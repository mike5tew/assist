package infoin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	weaviatefilters "github.com/weaviate/weaviate-go-client/v4/weaviate/filters" // 🆕 ADD THIS
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// HSGQueryService implements the tiered RAG retrieval system
type HSGQueryService struct {
	LlmClient      *llm.LlamaClient
	MongoClient    *db.MongoDB
	WeaviateClient *weaviate.Client
	normalizer     *SymbolNormalizer // 🆕 ADD THIS
}

// NewHSGQueryService creates a new HSG query service
func NewHSGQueryService() (*HSGQueryService, error) {
	mongoDB, err := db.NewFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Initialize Weaviate client directly to ensure it works
	weaviateURL := os.Getenv("WEAVIATE_URL")
	if weaviateURL == "" {
		weaviateURL = "http://localhost:8081"
	}

	// Parse the URL to validate and get components
	parsedURL, err := url.Parse(weaviateURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Weaviate URL: %w", err)
	}

	// Create weaviate client config
	cfg := weaviate.Config{
		Host:   parsedURL.Host,
		Scheme: parsedURL.Scheme,
	}

	// Add API key authentication if provided
	apiKey := os.Getenv("WEAVIATE_API_KEY")
	if apiKey != "" {
		cfg.Headers = map[string]string{
			"Authorization": "Bearer " + apiKey,
		}

		// Add API key auth
		authConfig := auth.ApiKey{Value: apiKey}
		cfg.AuthConfig = &authConfig
	}

	// Initialize the client
	weaviateClient := weaviate.New(cfg)

	// Test the connection
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ready, err := weaviateClient.Misc().ReadyChecker().Do(ctx)
	if err != nil {
		log.Printf("Warning: Weaviate connectivity test failed: %v", err)
		// Continue anyway, but log the warning
	} else if !ready {
		log.Printf("Warning: Weaviate is not ready")
		// Continue anyway, but log the warning
	} else {
		log.Printf("✅ Successfully connected to Weaviate at %s", weaviateURL)
	}

	return &HSGQueryService{
		LlmClient:      llm.NewLlamaClient(),
		MongoClient:    mongoDB,
		WeaviateClient: weaviateClient,
		normalizer:     NewSymbolNormalizer(), // 🆕 ADD THIS
	}, nil
}

// QueryHSG performs the complete tiered RAG retrieval and generation
func (s *HSGQueryService) QueryHSG(ctx context.Context, query string, domain string, maxResults int) (*models.RAGQueryContext, string, error) {
	log.Printf("🔍 HSG Query: Starting tiered retrieval for query: %s", query)

	// Initialize an empty RAG context for error cases
	ragContext := &models.RAGQueryContext{
		Query:           query,
		Domain:          domain,
		SummaryContext:  []string{},
		DocumentContext: []string{},
		SemanticContext: []string{},
		SourceChain:     []models.TraversalStep{},
	}

	// Check if Weaviate client is initialized
	if s.WeaviateClient == nil {
		return ragContext, "Sorry, the search system is not properly initialized.",
			fmt.Errorf("Weaviate client is not initialized")
	}

	// 🆕 Support "personal" domain and subject filtering
	if domain == "" {
		domain = "immunology" // Default fallback
	}

	log.Printf("HSG Query: '%s' in domain '%s'", query, domain)

	// 🆕 Create multiple search variants to handle different symbol representations
	searchVariants := s.normalizer.CreateSearchVariants(query)

	log.Printf("Query: '%s' → Search variants: %v", query, searchVariants)

	// Try each variant until we get results
	var semanticLinks []models.SemanticLink
	var err error

	for _, variant := range searchVariants {
		normalizedQuery := s.normalizer.NormalizeForVectorization(variant)

		// Generate embedding for normalized query
		embedding, err := s.LlmClient.GenerateEmbedding(normalizedQuery)
		if err != nil {
			log.Printf("Failed to generate embedding for variant '%s': %v", variant, err)
			continue
		}

		// Search Weaviate
		semanticLinks, err = s.searchSemanticLinks(ctx, embedding, domain, maxResults)
		if err != nil {
			log.Printf("Search failed for variant '%s': %v", variant, err)
			continue
		}

		if len(semanticLinks) > 0 {
			log.Printf("✅ Found %d results using search variant: '%s'", len(semanticLinks), variant)
			break
		}
	}

	// Step 1: Generate query vector
	queryVector, err := s.LlmClient.GenerateEmbedding(query)
	if err != nil {
		return ragContext, "I couldn't generate an embedding for your query.",
			fmt.Errorf("failed to generate query vector: %w", err)
	}

	// Step 2: Tier 2 - Search SummaryChunks first
	log.Printf("🔍 Tier 2: Searching SummaryChunks...")

	// Use the more reliable schema getter approach
	log.Printf("🔍 Getting full schema to verify SummaryChunk class...")
	schema, err := s.WeaviateClient.Schema().Getter().Do(ctx)
	if err != nil {
		log.Printf("❌ Error getting schema: %v", err)
		return ragContext, "I'm having trouble accessing the knowledge graph.",
			fmt.Errorf("failed to get schema: %v", err)
	}

	// Check for SummaryChunk class in retrieved schema
	classExists := false
	for _, class := range schema.Classes {
		log.Printf("Found schema class: %s", class.Class)
		if class.Class == "SummaryChunk" {
			classExists = true
			log.Printf("✅ SummaryChunk class verified in schema")
			break
		}
	}

	if !classExists {
		log.Printf("❌ SummaryChunk class not found in schema")
		return ragContext, "The search system needs to be set up. Please run 'make setup-hsg-schema' first.",
			fmt.Errorf("SummaryChunk class does not exist in Weaviate")
	}

	// Continue with the search now that we've confirmed the class exists
	summaryChunks, err := s.searchSummaryChunks(ctx, queryVector, domain, maxResults)
	if err != nil {
		log.Printf("Warning: Tier 2 search failed: %v", err)
		summaryChunks = []models.SummaryChunk{} // Continue with empty results
	}

	// Step 3: Tier 1 - Traverse to related SemanticLinks
	log.Printf("🔍 Tier 1: Traversing to SemanticLinks...")
	semanticLinks, err = s.traverseToSemanticLinks(ctx, summaryChunks, maxResults*2)
	if err != nil {
		log.Printf("Warning: Tier 1 traversal failed: %v", err)
		semanticLinks = []models.SemanticLink{}
	}

	// Step 4: Tier 0 - Traverse to source documents
	log.Printf("🔍 Tier 0: Traversing to source documents...")
	documents, err := s.traverseToDocuments(ctx, semanticLinks, maxResults*3)
	if err != nil {
		log.Printf("Warning: Tier 0 traversal failed: %v", err)
		documents = []models.SubjectContent{}
	}

	// Step 5: Build comprehensive context
	ragContext = s.buildRAGContext(query, summaryChunks, semanticLinks, documents)

	// Step 6: Generate final answer
	log.Printf("🔍 Final Generation: Synthesizing answer...")
	answer, err := s.generateFinalAnswer(ctx, ragContext)
	if err != nil {
		return ragContext, "I'm sorry, I couldn't generate a final answer based on the retrieved information.",
			fmt.Errorf("failed to generate final answer: %w", err)
	}

	log.Printf("✅ HSG Query Complete: Retrieved %d summary chunks, %d semantic links, %d documents",
		len(summaryChunks), len(semanticLinks), len(documents))

	return ragContext, answer, nil
}

// searchSummaryChunks performs vector search on Tier 2 SummaryChunks
func (h *HSGQueryService) searchSummaryChunks(ctx context.Context, queryVector []float32, domain string, limit int) ([]models.SummaryChunk, error) {
	// Check if Weaviate client is initialized
	if h.WeaviateClient == nil {
		return nil, fmt.Errorf("Weaviate client is not initialized")
	}

	// First try getting the schema - this is more reliable than class existence checker
	log.Printf("🔍 Getting full schema to verify SummaryChunk class...")
	schema, err := h.WeaviateClient.Schema().Getter().Do(ctx)
	if err != nil {
		log.Printf("❌ Error getting schema: %v", err)
		return nil, fmt.Errorf("failed to get schema: %v", err)
	}

	// Check for SummaryChunk class in retrieved schema
	classExists := false
	for _, class := range schema.Classes {
		if class.Class == "SummaryChunk" {
			classExists = true
			log.Printf("✅ SummaryChunk class verified in schema")
			break
		}
	}

	if !classExists {
		log.Printf("❌ SummaryChunk class not found in schema")
		return nil, fmt.Errorf("SummaryChunk class does not exist in Weaviate")
	}

	// Use Weaviate vector search on SummaryChunk class
	className := "SummaryChunk"

	// Print query info for debugging
	log.Printf("🔍 Executing vector search on class '%s' with %d dimensions", className, len(queryVector))

	nearVector := h.WeaviateClient.GraphQL().Get().
		WithClassName(className).
		WithFields(
			graphql.Field{Name: "content"}, // Using "content" property instead of "summary_text"
			graphql.Field{Name: "domain"},
			graphql.Field{Name: "chunkId"},
			graphql.Field{Name: "parentId"},
			graphql.Field{Name: "_additional", Fields: []graphql.Field{
				{Name: "id"},
				{Name: "distance"},
			}},
		).
		WithNearVector(h.WeaviateClient.GraphQL().NearVectorArgBuilder().
			WithVector(queryVector).
			WithCertainty(0.7)).
		WithLimit(limit)

	// Add domain filter if specified
	if domain != "" {
		whereFilter := filters.Where().
			WithPath([]string{"domain"}).
			WithOperator(filters.Equal).
			WithValueString(domain)
		nearVector = nearVector.WithWhere(whereFilter)
	}

	// Execute the search with extra error handling
	result, err := nearVector.Do(ctx)
	if err != nil {
		// Print detailed error information
		log.Printf("❌ Weaviate search failed: %v", err)

		// Check for common error patterns
		if strings.Contains(strings.ToLower(err.Error()), "class not found") ||
			strings.Contains(strings.ToLower(err.Error()), "not exist") {
			log.Printf("⚠️ Class existence error detected in search response")
			return nil, fmt.Errorf("SummaryChunk class does not exist in Weaviate")
		}

		return nil, fmt.Errorf("weaviate search failed: %w", err)
	}

	// Fix: Convert result.Data to the correct type
	if result.Data == nil {
		log.Printf("⚠️ Weaviate returned nil data")
		return []models.SummaryChunk{}, nil
	}

	// Convert weaviate models.JSONObject to map[string]interface{}
	dataMap := make(map[string]interface{})
	for key, value := range result.Data {
		dataMap[key] = value
	}

	// Parse results and fetch full objects from MongoDB
	return h.parseSummaryChunkResults(ctx, dataMap)
}

// searchSemanticLinks performs vector search on SemanticLinks class
func (s *HSGQueryService) searchSemanticLinks(ctx context.Context, queryVector []float32, domain string, limit int) ([]models.SemanticLink, error) {
	// 🆕 ADD THIS: Actual Weaviate GraphQL query
	fields := []graphql.Field{
		{Name: "source_term"},
		{Name: "target_term"},
		{Name: "relation_type"},
		{Name: "context"},
		{Name: "confidence"},
		{Name: "_additional", Fields: []graphql.Field{{Name: "distance"}}},
	}

	nearVector := s.WeaviateClient.GraphQL().
		NearVectorArgBuilder().
		WithVector(queryVector)

	queryBuilder := s.WeaviateClient.GraphQL().Get().
		WithClassName("SemanticLinks").
		WithFields(fields...).
		WithNearVector(nearVector).
		WithLimit(limit)

	if domain != "" {
		whereFilter := weaviatefilters.Where().
			WithPath([]string{"domain"}).
			WithOperator(weaviatefilters.Equal).
			WithValueString(domain)
		queryBuilder = queryBuilder.WithWhere(whereFilter)
	}

	response, err := queryBuilder.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("weaviate search failed: %w", err)
	}

	if response.Errors != nil {
		return nil, fmt.Errorf("weaviate query returned errors: %v", response.Errors)
	}

	data, ok := response.Data["Get"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from Weaviate")
	}

	items, ok := data["SemanticLinks"].([]interface{})
	if !ok {
		return []models.SemanticLink{}, nil
	}

	var semanticLinks []models.SemanticLink
	for _, item := range items {
		if itemMap, ok := item.(map[string]interface{}); ok {
			link, err := s.parseSemanticLinkFromWeaviate(itemMap)
			if err != nil {
				log.Printf("Warning: Failed to parse semantic link: %v", err)
				continue
			}
			semanticLinks = append(semanticLinks, *link)
		}
	}

	return semanticLinks, nil
}

// parseSemanticLinkFromWeaviate converts Weaviate response to SemanticLink model
func (s *HSGQueryService) parseSemanticLinkFromWeaviate(linkMap map[string]interface{}) (*models.SemanticLink, error) {
	sourceMongoID, _ := linkMap["source_mongo_id"].(string)
	targetMongoID, _ := linkMap["target_mongo_id"].(string)

	sourceID, _ := primitive.ObjectIDFromHex(sourceMongoID)
	targetID, _ := primitive.ObjectIDFromHex(targetMongoID)

	return &models.SemanticLink{
		SourceTerm:   getString(linkMap, "source_term"),
		TargetTerm:   getString(linkMap, "target_term"),
		SourceID:     sourceID,
		TargetID:     targetID,
		RelationType: getString(linkMap, "relation_type"),
		Context:      getString(linkMap, "context"),
		Confidence:   getFloat(linkMap, "confidence"),
	}, nil
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0.0
}

// traverseToSemanticLinks retrieves SemanticLinks connected to SummaryChunks
func (h *HSGQueryService) traverseToSemanticLinks(ctx context.Context, summaryChunks []models.SummaryChunk, limit int) ([]models.SemanticLink, error) {
	if len(summaryChunks) == 0 {
		return []models.SemanticLink{}, nil
	}

	// Collect all semantic link IDs from summary chunks
	var linkIDs []primitive.ObjectID
	for _, chunk := range summaryChunks {
		linkIDs = append(linkIDs, chunk.SemanticLinkIDs...)
	}

	if len(linkIDs) == 0 {
		return []models.SemanticLink{}, nil
	}

	// Query MongoDB for semantic links
	collection := h.MongoClient.Database.Collection("semantic_links")
	filter := bson.M{"_id": bson.M{"$in": linkIDs}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query semantic links: %w", err)
	}
	defer cursor.Close(ctx)

	var semanticLinks []models.SemanticLink
	if err := cursor.All(ctx, &semanticLinks); err != nil {
		return nil, fmt.Errorf("failed to decode semantic links: %w", err)
	}

	// Limit results if necessary
	if len(semanticLinks) > limit {
		semanticLinks = semanticLinks[:limit]
	}

	return semanticLinks, nil
}

// traverseToDocuments retrieves source documents connected to SemanticLinks
func (h *HSGQueryService) traverseToDocuments(ctx context.Context, semanticLinks []models.SemanticLink, limit int) ([]models.SubjectContent, error) {
	if len(semanticLinks) == 0 {
		return []models.SubjectContent{}, nil
	}

	// Extract document IDs from semantic links metadata
	documentIDMap := make(map[string]bool)
	for _, link := range semanticLinks {
		if docID, ok := link.Metadata["document_id"].(string); ok {
			documentIDMap[docID] = true
		}
	}

	var documentObjectIDs []primitive.ObjectID
	for docIDStr := range documentIDMap {
		if objID, err := primitive.ObjectIDFromHex(docIDStr); err == nil {
			documentObjectIDs = append(documentObjectIDs, objID)
		}
	}

	if len(documentObjectIDs) == 0 {
		return []models.SubjectContent{}, nil
	}

	// Query MongoDB for documents (try multiple collections)
	documents := []models.SubjectContent{}

	collections := []string{"immunology_content", "subject_content", "unified_content"}
	for _, collName := range collections {
		collection := h.MongoClient.Database.Collection(collName)
		filter := bson.M{"_id": bson.M{"$in": documentObjectIDs}}

		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			continue // Try next collection
		}

		var collDocs []models.SubjectContent
		if err := cursor.All(ctx, &collDocs); err == nil {
			documents = append(documents, collDocs...)
		}
		cursor.Close(ctx)
	}

	// Limit results if necessary
	if len(documents) > limit {
		documents = documents[:limit]
	}

	return documents, nil
}

// buildRAGContext assembles the comprehensive context for final generation
func (h *HSGQueryService) buildRAGContext(query string, summaryChunks []models.SummaryChunk, semanticLinks []models.SemanticLink, documents []models.SubjectContent) *models.RAGQueryContext {
	ragContext := &models.RAGQueryContext{
		Query:           query,
		SummaryContext:  []string{},
		SemanticContext: []string{},
		DocumentContext: []string{},
		SourceChain:     []models.TraversalStep{},
	}

	// Add summary context
	for i, chunk := range summaryChunks {
		ragContext.SummaryContext = append(ragContext.SummaryContext, chunk.SummaryText)
		ragContext.SourceChain = append(ragContext.SourceChain, models.TraversalStep{
			Tier:       2,
			ObjectID:   chunk.ID.Hex(),
			ObjectType: "SummaryChunk",
			Content:    chunk.SummaryText[:min(100, len(chunk.SummaryText))] + "...",
			Confidence: chunk.Confidence,
		})

		if i >= 3 { // Limit to prevent context overflow
			break
		}
	}

	// Add semantic context
	for i, link := range semanticLinks {
		contextStr := fmt.Sprintf("%s → %s (%s): %s",
			link.SourceTerm, link.TargetTerm, link.RelationType, link.Context)
		ragContext.SemanticContext = append(ragContext.SemanticContext, contextStr)
		ragContext.SourceChain = append(ragContext.SourceChain, models.TraversalStep{
			Tier:       1,
			ObjectID:   link.ID.Hex(),
			ObjectType: "SemanticLink",
			Content:    contextStr,
			Confidence: link.Confidence,
		})

		if i >= 5 { // Limit to prevent context overflow
			break
		}
	}

	// Add document context
	for i, doc := range documents {
		content := doc.Content
		if len(content) > 500 {
			content = content[:500] + "..."
		}
		ragContext.DocumentContext = append(ragContext.DocumentContext, content)
		ragContext.SourceChain = append(ragContext.SourceChain, models.TraversalStep{
			Tier:       0,
			ObjectID:   doc.ID.Hex(),
			ObjectType: "Document",
			Content:    content,
			Confidence: 1.0, // Documents have inherent confidence
		})

		if i >= 3 { // Limit to prevent context overflow
			break
		}
	}

	return ragContext
}

// generateFinalAnswer synthesizes the final answer from the assembled context
func (h *HSGQueryService) generateFinalAnswer(ctx context.Context, ragContext *models.RAGQueryContext) (string, error) {
	// Build comprehensive prompt
	var promptBuilder strings.Builder

	promptBuilder.WriteString(fmt.Sprintf("User Query: %s\n\n", ragContext.Query))

	if len(ragContext.SummaryContext) > 0 {
		promptBuilder.WriteString("High-Level Summaries:\n")
		for i, summary := range ragContext.SummaryContext {
			promptBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, summary))
		}
		promptBuilder.WriteString("\n")
	}

	if len(ragContext.SemanticContext) > 0 {
		promptBuilder.WriteString("Semantic Relationships:\n")
		for i, semantic := range ragContext.SemanticContext {
			promptBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, semantic))
		}
		promptBuilder.WriteString("\n")
	}

	if len(ragContext.DocumentContext) > 0 {
		promptBuilder.WriteString("Source Documents:\n")
		for i, doc := range ragContext.DocumentContext {
			promptBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, doc))
		}
		promptBuilder.WriteString("\n")
	}

	promptBuilder.WriteString(`Based on the hierarchical context provided above, please provide a comprehensive, accurate answer to the user's query. 
Use the high-level summaries for conceptual understanding, the semantic relationships for detailed connections, and the source documents for specific evidence.

Answer:`)

	return h.LlmClient.Generate(ctx, promptBuilder.String())
}

// Helper function parseSummaryChunkResults - Parse Weaviate search results and fetch full objects from MongoDB
func (h *HSGQueryService) parseSummaryChunkResults(ctx context.Context, data map[string]interface{}) ([]models.SummaryChunk, error) {
	log.Printf("🔍 Parsing Weaviate SummaryChunk search results...")

	// Add data dump for debugging - helps understand the structure
	bytes, _ := json.Marshal(data)
	log.Printf("📊 Weaviate response data structure: %s", string(bytes[:min(500, len(bytes))]))

	// Parse the Weaviate GraphQL response structure
	// Expected structure: data["Get"]["SummaryChunk"] = []interface{}

	getSection, ok := data["Get"].(map[string]interface{})
	if !ok {
		log.Printf("⚠️ No 'Get' section found in Weaviate response")
		return []models.SummaryChunk{}, nil
	}

	// Fix: Use singular "SummaryChunk" to match the class name we created
	summaryChunksData, ok := getSection["SummaryChunk"].([]interface{})
	if !ok {
		log.Printf("⚠️ No 'SummaryChunk' array found in Weaviate Get section")

		// Add more detailed response structure logging
		for k := range getSection {
			log.Printf("📋 Found section: %s", k)
		}

		// Try both singular and plural forms
		pluralData, pluralOk := getSection["SummaryChunks"].([]interface{})
		if pluralOk {
			log.Printf("ℹ️ Found results under plural name 'SummaryChunks' instead")
			summaryChunksData = pluralData
		} else {
			log.Printf("❌ Neither 'SummaryChunk' nor 'SummaryChunks' found in response")
			return []models.SummaryChunk{}, nil
		}
	}

	if len(summaryChunksData) == 0 {
		log.Printf("ℹ️ No SummaryChunk results found in Weaviate")
		return []models.SummaryChunk{}, nil
	}

	log.Printf("📊 Found %d SummaryChunk results from Weaviate", len(summaryChunksData))

	// Extract Weaviate IDs and collect basic properties
	var weaviateIDs []string
	var weaviateResults []map[string]interface{}

	for i, item := range summaryChunksData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			log.Printf("⚠️ Skipping malformed result #%d", i)
			continue
		}

		// Extract Weaviate ID from _additional section
		var weaviateID string
		if additional, ok := itemMap["_additional"].(map[string]interface{}); ok {
			if id, ok := additional["id"].(string); ok {
				weaviateID = id
			}
		}

		if weaviateID == "" {
			log.Printf("⚠️ No Weaviate ID found for result #%d", i)
			continue
		}

		weaviateIDs = append(weaviateIDs, weaviateID)
		weaviateResults = append(weaviateResults, itemMap)

		// Log what we extracted from Weaviate
		summaryText, _ := itemMap["summary_text"].(string)
		domain, _ := itemMap["domain"].(string)
		confidence, _ := itemMap["confidence"].(float64)

		log.Printf("📄 Weaviate Result #%d: ID=%s, Domain=%s, Confidence=%.2f, Preview=%s...",
			i+1, weaviateID[:8], domain, confidence,
			func() string {
				if len(summaryText) > 50 {
					return summaryText[:50]
				}
				return summaryText
			}())
	}

	if len(weaviateIDs) == 0 {
		log.Printf("⚠️ No valid Weaviate IDs extracted from results")
		return []models.SummaryChunk{}, nil
	}

	// Now fetch the complete SummaryChunk objects from MongoDB
	// We need to find them by some correlation - either store Weaviate ID in MongoDB metadata
	// or use another approach to match them

	log.Printf("🔍 Fetching full SummaryChunk objects from MongoDB...")
	summaryChunksCollection := h.MongoClient.Database.Collection("summary_chunks")

	// Strategy 1: If we store Weaviate IDs in MongoDB metadata
	filter := bson.M{
		"metadata.weaviate_id": bson.M{"$in": weaviateIDs},
	}

	cursor, err := summaryChunksCollection.Find(ctx, filter)
	if err != nil {
		// Strategy 2: Fallback - get recent summary chunks and match by content similarity
		log.Printf("⚠️ Failed to find by Weaviate ID, trying fallback approach: %v", err)
		return h.fallbackSummaryChunkRetrieval(ctx, weaviateResults)
	}
	defer cursor.Close(ctx)

	var summaryChunks []models.SummaryChunk
	if err := cursor.All(ctx, &summaryChunks); err != nil {
		return nil, fmt.Errorf("failed to decode SummaryChunk objects: %w", err)
	}

	log.Printf("✅ Retrieved %d complete SummaryChunk objects from MongoDB", len(summaryChunks))

	// Log what we're returning for debugging
	for i, chunk := range summaryChunks {
		log.Printf("📄 MongoDB SummaryChunk #%d: ID=%s, Domain=%s, LinkCount=%d, Confidence=%.2f",
			i+1, chunk.ID.Hex(), chunk.Domain, len(chunk.SemanticLinkIDs), chunk.Confidence)
	}

	return summaryChunks, nil
}

// fallbackSummaryChunkRetrieval - Alternative strategy when Weaviate ID matching fails
func (h *HSGQueryService) fallbackSummaryChunkRetrieval(ctx context.Context, weaviateResults []map[string]interface{}) ([]models.SummaryChunk, error) {
	log.Printf("🔄 Using fallback retrieval strategy for SummaryChunks...")

	summaryChunksCollection := h.MongoClient.Database.Collection("summary_chunks")

	// Get recent summary chunks (last 100) and try to match by content similarity
	limit := int64(100)
	cursor, err := summaryChunksCollection.Find(ctx, bson.M{},
		&options.FindOptions{
			Sort:  bson.M{"created_at": -1},
			Limit: &limit,
		})
	if err != nil {
		return nil, fmt.Errorf("fallback query failed: %w", err)
	}
	defer cursor.Close(ctx)

	var allChunks []models.SummaryChunk
	if err := cursor.All(ctx, &allChunks); err != nil {
		return nil, fmt.Errorf("failed to decode fallback chunks: %w", err)
	}

	// Match by partial content similarity (simple string matching)
	var matchedChunks []models.SummaryChunk

	for _, weaviateResult := range weaviateResults {
		weaviateSummary, ok := weaviateResult["summary_text"].(string)
		if !ok || len(weaviateSummary) < 20 {
			continue
		}

		// Look for MongoDB chunks with similar content
		for _, mongoChunk := range allChunks {
			if len(mongoChunk.SummaryText) < 20 {
				continue
			}

			// Simple similarity check - first 50 characters
			weaviatePrefix := weaviateSummary
			if len(weaviatePrefix) > 50 {
				weaviatePrefix = weaviatePrefix[:50]
			}

			mongoPrefix := mongoChunk.SummaryText
			if len(mongoPrefix) > 50 {
				mongoPrefix = mongoPrefix[:50]
			}

			if strings.Contains(weaviatePrefix, mongoPrefix) || strings.Contains(mongoPrefix, weaviatePrefix) {
				// Check for duplicates before adding
				var alreadyAdded bool
				for _, added := range matchedChunks {
					if added.ID == mongoChunk.ID {
						alreadyAdded = true
						break
					}
				}
				if alreadyAdded {
					continue
				}

				// Add to results
				if len(matchedChunks) >= 5 { // Limit fallback results
					break
				}
				log.Printf("📄 Matched SummaryChunk by content similarity: %s", mongoChunk.ID.Hex())
				matchedChunks = append(matchedChunks, mongoChunk)
			}
		}

		if len(matchedChunks) >= 5 {
			break // Limit total fallback results
		}
	}

	if len(matchedChunks) == 0 {
		log.Printf("⚠️ Fallback retrieval found no matched SummaryChunks")
	} else {
		log.Printf("🔄 Fallback retrieval found %d matched SummaryChunks", len(matchedChunks))
	}

	return matchedChunks, nil

}

// Helper to infer reverse relationship types
func inferReverseRelationType(relationType string) string {
	reverseRelations := map[string]string{
		"expresses":     "expressed_by",
		"causes":        "caused_by",
		"develops_into": "develops_from",
		"part_of":       "contains",
		"enables":       "enabled_by",
		"inhibits":      "inhibited_by",
	}

	if reverse, exists := reverseRelations[relationType]; exists {
		return reverse
	}

	// Default: add "_inverse" suffix
	return relationType + "_inverse"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
