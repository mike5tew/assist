package infoin

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	wvmodels "github.com/weaviate/weaviate/entities/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Defintion of the SemanticLinkService
// SL are the relationships between key terms in documents
// HSG are the higher-level summary chunks derived from groups of SL

// SemanticLinkService handles extraction and management of semantic links
type SemanticLinkService struct {
	LlmClient      *llm.LlamaClient
	MongoClient    *db.MongoDB
	WeaviateClient *weaviate.Client
	normalizer     *SymbolNormalizer // 🆕 ADD THIS
}

// NewSemanticLinkService creates a new instance of the semantic link service
func NewSemanticLinkService() (*SemanticLinkService, error) {
	mongoDB, err := db.NewFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	return &SemanticLinkService{
		LlmClient:      llm.NewLlamaClient(),
		MongoClient:    mongoDB,
		WeaviateClient: db.GetWeaviateClient(),
		normalizer:     NewSymbolNormalizer(), // 🆕 ADD THIS
	}, nil
}

// ProcessDocument extracts semantic links from a single document.
func (s *SemanticLinkService) ProcessDocument(ctx context.Context, doc models.SubjectContent) error {
	log.Printf("🔍 Extracting semantic links from: %s", doc.Title)

	// 🆕 ADD THIS: LLM prompt for extraction
	prompt := fmt.Sprintf(`Extract semantic relationships from this medical text. Return JSON array.

Text: %s

Format each relationship as:
{
  "source_term": "term A",
  "target_term": "term B", 
  "relation_type": "causes|treats|is_a|part_of|etc",
  "context": "brief explanation",
  "confidence": 0.0-1.0
}

Example:
[
  {
    "source_term": "X-linked agammaglobulinemia",
    "target_term": "B-cell deficiency",
    "relation_type": "causes",
    "context": "BTK gene mutation prevents B-cell maturation",
    "confidence": 0.95
  }
]`, doc.Content)

	// Call AWS Bedrock Claude
	response, err := s.LlmClient.Generate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("LLM extraction failed: %w", err)
	}

	// Parse JSON response
	var extractedLinks []models.SemanticLink
	if err := json.Unmarshal([]byte(response), &extractedLinks); err != nil {
		return fmt.Errorf("failed to parse LLM response: %w", err)
	}

	// Store each link
	for _, link := range extractedLinks {
		link.Domain = doc.Domain
		link.SourceID = doc.ID
		// Generate embedding and store in Weaviate
		if err := s.StoreSemanticLink(ctx, link); err != nil {
			log.Printf("Warning: Failed to store link: %v", err)
		}
	}

	return nil
}

// ProcessDocumentBatch extracts semantic links from a batch of documents
func (s *SemanticLinkService) ProcessDocumentBatch(ctx context.Context, batchID string, collectionName string) error {
	collection := s.MongoClient.Database.Collection(collectionName)

	// Get all documents from this batch
	filter := bson.M{"source.upload_id": batchID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("error finding documents: %w", err)
	}
	defer cursor.Close(ctx)

	// Create semantic links collection if it doesn't exist
	semanticLinksCollection := s.MongoClient.Database.Collection("semantic_links")

	// Collect all documents to process
	var documents []models.SubjectContent
	for cursor.Next(ctx) {
		var doc models.SubjectContent
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Warning: Failed to decode document: %v", err)
			continue
		}
		documents = append(documents, doc)
	}

	log.Printf("Processing %d documents for semantic link extraction", len(documents))

	// Process each document to extract terms and relationships
	for _, doc := range documents {
		// 1. Extract key terms from the document
		terms, err := s.extractKeyTerms(ctx, doc.Content)
		if err != nil {
			log.Printf("Warning: Failed to extract terms from document %s: %v", doc.ID.Hex(), err)
			continue
		}

		// 2. Verify terms against MongoDB skills collection
		verifiedTerms, err := s.verifyTermsInMongoDB(ctx, terms)
		if err != nil {
			log.Printf("Warning: Failed to verify terms: %v", err)
			// Continue with unverified terms
			verifiedTerms = make([]models.Term, len(terms))
			for i, termText := range terms {
				verifiedTerms[i] = models.Term{
					Text:     termText,
					Verified: false,
					MongoID:  primitive.NilObjectID,
				}
			}
		}

		// 3. Extract relationships between terms
		links, err := s.extractRelationships(ctx, doc.Content, verifiedTerms)
		if err != nil {
			log.Printf("Warning: Failed to extract relationships: %v", err)
			continue
		}

		// 4. Store links in MongoDB with vectors
		insertedCount, err := s.StoreSemanticLinks(ctx, semanticLinksCollection, links, doc.ID)
		if err != nil {
			log.Printf("Warning: Failed to store semantic links: %v", err)
			continue
		}

		log.Printf("✅ Processed document %s: Found %d terms, created %d semantic links",
			doc.ID.Hex(), len(verifiedTerms), insertedCount)

		// 5. Update the original document with link information
		s.updateDocumentWithLinkInfo(ctx, collection, doc.ID, insertedCount)
	}

	return nil
}

// extractKeyTerms uses LLaMA to identify key domain-specific terms in content
func (s *SemanticLinkService) extractKeyTerms(ctx context.Context, content string) ([]string, error) {
	// Trim content if it's too long
	if len(content) > 8000 {
		content = content[:8000] + "..."
	}

	prompt := fmt.Sprintf(`Extract the key domain-specific terms from the following text.
Return only a JSON array of strings with the important terms.

Text:
%s

Output Format:
["term1", "term2", "term3", ...]`, content)

	response, err := s.LlmClient.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLaMA generation failed: %w", err)
	}

	// Validate response
	if response == "" {
		return nil, fmt.Errorf("received empty response from LLaMA")
	}

	// Extract JSON array from response
	startIdx := strings.Index(response, "[")
	endIdx := strings.LastIndex(response, "]")
	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return nil, fmt.Errorf("invalid JSON format in LLaMA response: %s", response)
	}

	jsonArray := response[startIdx : endIdx+1]
	var terms []string
	if err := json.Unmarshal([]byte(jsonArray), &terms); err != nil {
		return nil, fmt.Errorf("error unmarshaling terms from response '%s': %w", jsonArray, err)
	}

	return terms, nil
}

func (s *SemanticLinkService) verifyTermsInMongoDB(ctx context.Context, terms []string) ([]models.Term, error) {
	if s.MongoClient == nil {
		return nil, fmt.Errorf("mongodb client not initialized")
	}

	skillsCollection := s.MongoClient.Database.Collection("skills")
	var verifiedTerms []models.Term

	for _, termText := range terms {
		filter := bson.M{"name": bson.M{"$regex": fmt.Sprintf("^%s$", termText), "$options": "i"}}
		var skill models.EducationalSkill

		err := skillsCollection.FindOne(ctx, filter).Decode(&skill)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// Term not found
				verifiedTerms = append(verifiedTerms, models.Term{Text: termText, Verified: false})
			} else {
				// Other database error
				log.Printf("Warning: Error verifying term '%s' in MongoDB: %v", termText, err)
				verifiedTerms = append(verifiedTerms, models.Term{Text: termText, Verified: false})
			}
			continue
		}

		// Term found and verified
		verifiedTerms = append(verifiedTerms, models.Term{
			Text:     skill.Name,
			Verified: true,
			MongoID:  skill.ID,
			Metadata: map[string]interface{}{"description": skill.Description},
		})
	}

	return verifiedTerms, nil
}

func (s *SemanticLinkService) verifyTermsInWeaviate(ctx context.Context, terms []string) ([]models.Term, error) {
	if s.WeaviateClient == nil {
		return nil, fmt.Errorf("weaviate client not initialized")
	}

	var verifiedTerms []models.Term
	for _, term := range terms {
		whereFilter := filters.Where().
			WithPath([]string{"name"}).
			WithOperator(filters.Equal).
			WithValueString(term)

		result, err := s.WeaviateClient.GraphQL().
			Get().
			WithClassName("EducationalSkills"). // was "Skills"
			WithFields(
				graphql.Field{Name: "mongo_id"},
				graphql.Field{Name: "name"},
				graphql.Field{Name: "description"},
			).
			WithWhere(whereFilter).
			WithLimit(1).
			Do(ctx)

		if err != nil || (result != nil && len(result.Errors) > 0) {
			verifiedTerms = append(verifiedTerms, models.Term{
				Text:     term,
				Verified: false,
				MongoID:  primitive.NilObjectID,
			})
			continue
		}

		exists, mongoID, metadata := s.processWeaviateSearchResult(result.Data, term)

		verifiedTerms = append(verifiedTerms, models.Term{
			Text:     term,
			Verified: exists,
			MongoID:  mongoID,
			Metadata: metadata,
		})
	}
	return verifiedTerms, nil
}

// processWeaviateSearchResult extracts information from Weaviate search results
func (s *SemanticLinkService) processWeaviateSearchResult(
	data map[string]wvmodels.JSONObject,
	_ string,
) (bool, primitive.ObjectID, map[string]interface{}) {
	if data == nil {
		return false, primitive.NilObjectID, nil
	}

	// Convert the top-level data to a standard map for easier processing
	dataMap := make(map[string]interface{})
	for k, v := range data {
		dataMap[k] = v
	}

	getSection, ok := dataMap["Get"].(map[string]interface{})
	if !ok {
		return false, primitive.NilObjectID, nil
	}

	items, ok := getSection["EducationalSkills"].([]interface{}) // was getSection["Skills"]
	if !ok || len(items) == 0 {
		return false, primitive.NilObjectID, nil
	}

	item, ok := items[0].(map[string]interface{})
	if !ok {
		return false, primitive.NilObjectID, nil
	}

	mongoIDStr, _ := item["mongo_id"].(string)
	var mongoID primitive.ObjectID
	if mongoIDStr != "" {
		if id, err := primitive.ObjectIDFromHex(mongoIDStr); err == nil {
			mongoID = id
		}
	}

	metadata := make(map[string]interface{})
	if desc, ok := item["description"].(string); ok && desc != "" {
		metadata["description"] = desc
	}

	return true, mongoID, metadata
}

// extractRelationships uses LLaMA to identify relationships between verified terms
func (s *SemanticLinkService) extractRelationships(ctx context.Context, content string, verifiedTerms []models.Term) ([]models.SemanticLink, error) {
	// Convert verifiedTerms to string slice for prompt
	termStrings := make([]string, len(verifiedTerms))
	for _, term := range verifiedTerms {
		termStrings = append(termStrings, term.Text)
	}

	// Use sample-based prompt builder
	domain := "immunology" // TODO: Extract from context
	prompt := s.buildExtractionPromptWithSamples(content, domain)

	response, err := s.LlmClient.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLaMA generation failed: %w", err)
	}

	// Validate response
	if response == "" {
		return nil, fmt.Errorf("received empty response from LLaMA")
	}

	// Extract JSON array from response
	startIdx := strings.Index(response, "[")
	endIdx := strings.LastIndex(response, "]")
	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return nil, fmt.Errorf("invalid JSON format in LLaMA response: %s", response)
	}

	jsonArray := response[startIdx : endIdx+1]
	var rawRelationships []struct {
		SourceTerm     string  `json:"source_term"`
		TargetTerm     string  `json:"target_term"`
		RelationType   string  `json:"relation_type"`
		Context        string  `json:"context"`
		Strength       float64 `json:"strength"`
		IsCaseSpecific bool    `json:"is_case_specific"`
		CaseStudyID    string  `json:"case_study_id"`
	}

	if err := json.Unmarshal([]byte(jsonArray), &rawRelationships); err != nil {
		return nil, fmt.Errorf("error unmarshaling relationships from response '%s': %w", jsonArray, err)
	}

	// Convert to SemanticLink objects
	var links []models.SemanticLink
	for _, rel := range rawRelationships {
		// Filter out low-confidence relationships
		if rel.Strength < 0.5 {
			continue
		}

		// Find the source and target terms in our verified terms
		var sourceMongoID, targetMongoID primitive.ObjectID

		var sourceVerified, targetVerified bool

		for _, term := range verifiedTerms {
			if strings.EqualFold(term.Text, rel.SourceTerm) {
				sourceMongoID = term.MongoID
				sourceVerified = term.Verified
			}
			if strings.EqualFold(term.Text, rel.TargetTerm) {
				targetMongoID = term.MongoID
				targetVerified = term.Verified
			}
		}

		// Create the semantic link
		link := models.SemanticLink{
			ID:           primitive.NewObjectID(),
			SourceID:     sourceMongoID,
			TargetID:     targetMongoID,
			SourceTerm:   rel.SourceTerm,
			TargetTerm:   rel.TargetTerm,
			RelationType: rel.RelationType,
			Context:      rel.Context,
			Confidence:   rel.Strength,
			Domain:       "extracted_relationship",
			CreatedAt:    time.Now(),
			Metadata: map[string]interface{}{

				"source_verified": sourceVerified,
				"target_verified": targetVerified,
			},
		}

		links = append(links, link)
	}

	// 🆕 NEW: Compute relativistic hierarchy metrics for each link
	for i := range links {
		metrics, err := s.computeHierarchyMetrics(ctx, &links[i], content)
		if err != nil {
			log.Printf("Warning: Failed to compute hierarchy metrics for link %d: %v", i, err)
			// Set defaults if computation fails
			links[i].SourceTermGenerality = 0.5
			links[i].TargetTermGenerality = 0.5
			links[i].SemanticDistance = 0.5
			links[i].RelationshipStrength = links[i].Confidence
			continue
		}

		// Update link with computed metrics
		links[i].SourceTermGenerality = metrics.SourceTermGenerality
		links[i].TargetTermGenerality = metrics.TargetTermGenerality
		links[i].SemanticDistance = metrics.SemanticDistance
		links[i].RelationshipStrength = metrics.RelationshipStrength
	}

	return links, nil
}

// 🆕 NEW: Compute hierarchy metrics using pattern-based LLM prompting
func (s *SemanticLinkService) computeHierarchyMetrics(
	ctx context.Context,
	link *models.SemanticLink,
	fullContext string,
) (*HierarchyMetrics, error) {
	// Extract relevant context
	contextSnippet := extractRelevantContext(fullContext, link.SourceTerm, link.TargetTerm, 500)

	// Generate pattern-based prompt
	prompt := HSGAnalysisPrompt(
		link.SourceTerm,
		link.TargetTerm,
		link.RelationType,
		contextSnippet,
	)

	// Call LLM
	response, err := s.LlmClient.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLM hierarchy analysis failed: %w", err)
	}

	// Parse JSON response
	var initialMetrics HierarchyMetrics
	if err := json.Unmarshal([]byte(cleanJSONResponse(response)), &initialMetrics); err != nil {
		return nil, fmt.Errorf("failed to parse hierarchy metrics: %w", err)
	}

	// 🆕 VALIDATION STEP: Compare against curated examples
	validator, err := NewHierarchyValidator("docs/training_data/curated_semantic_links.json")
	if err != nil {
		log.Printf("Warning: Failed to load validator, skipping validation: %v", err)
		return &initialMetrics, nil // Proceed without validation
	}

	validationReport := validator.ValidateMetrics(initialMetrics)

	if !validationReport.IsValid {
		log.Printf("⚠️ Metrics validation failed (error: %.1f%%). Requesting correction from LLM...",
			validationReport.AverageError*100)

		// 🆕 CORRECTION STEP: Send feedback to LLM
		correctedMetrics, err := s.CorrectMetricsWithFeedback(ctx, link, initialMetrics, validationReport)
		if err != nil {
			log.Printf("Warning: Failed to correct metrics: %v. Using initial metrics.", err)
			return &initialMetrics, nil
		}

		// Validate corrected metrics
		correctedValidation := validator.ValidateMetrics(*correctedMetrics)
		if correctedValidation.IsValid {
			log.Printf("✅ Corrected metrics are now valid (error: %.1f%%)",
				correctedValidation.AverageError*100)
			return correctedMetrics, nil
		} else {
			log.Printf("⚠️ Correction didn't improve metrics significantly. Using corrected version anyway.")
			return correctedMetrics, nil
		}
	}

	log.Printf("✅ Metrics validated successfully (error: %.1f%%)",
		validationReport.AverageError*100)
	return &initialMetrics, nil
}

// HierarchyMetrics holds the computed relativistic hierarchy values
type HierarchyMetrics struct {
	SourceTermGenerality float64 `json:"source_term_generality"`
	TargetTermGenerality float64 `json:"target_term_generality"`
	SemanticDistance     float64 `json:"semantic_distance"`
	RelationshipStrength float64 `json:"relationship_strength"`
	Explanation          string  `json:"explanation"`
}

// extractRelevantContext finds text surrounding both terms
func extractRelevantContext(fullText, term1, term2 string, maxChars int) string {
	lowerText := strings.ToLower(fullText)
	lowerTerm1 := strings.ToLower(term1)
	lowerTerm2 := strings.ToLower(term2)

	idx1 := strings.Index(lowerText, lowerTerm1)
	idx2 := strings.Index(lowerText, lowerTerm2)

	if idx1 == -1 && idx2 == -1 {
		// Neither term found, return first maxChars
		if len(fullText) > maxChars {
			return fullText[:maxChars]
		}
		return fullText
	}

	// Find start and end indices to capture both terms
	start := idx1
	if idx2 != -1 && (idx1 == -1 || idx2 < idx1) {
		start = idx2
	}

	end := idx1 + len(lowerTerm1)
	if idx2 != -1 && idx2+len(lowerTerm2) > end {
		end = idx2 + len(lowerTerm2)
	}

	// Expand context window
	start = max(0, start-maxChars/2)
	end = min(len(fullText), end+maxChars/2)

	return fullText[start:end]
}

// cleanJSONResponse strips markdown code fences and extracts JSON
func cleanJSONResponse(response string) string {
	// Remove markdown code fences
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")

	// Trim whitespace
	response = strings.TrimSpace(response)

	return response
}

// New helper method to ensure schema is correct
func (s *SemanticLinkService) ensureCleanSchema(ctx context.Context) error {
	if s.WeaviateClient == nil {
		return fmt.Errorf("weaviate client not initialized")
	}

	// The only class we care about for semantic processing is SemanticLinks.
	// SubjectAreaContent is deprecated.

	// Also ensure SemanticLinks class is ready
	exists, err := s.WeaviateClient.Schema().ClassExistenceChecker().WithClassName("SemanticLinks").Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to check SemanticLinks class existence: %w", err)
	}

	if !exists {
		if err := db.CreateSemanticLinksClass(ctx); err != nil {
			return fmt.Errorf("failed to create SemanticLinks class: %w", err)
		}
		log.Printf("✅ Successfully created SemanticLinks class with vectorizer=none")
	}

	return nil
}

// ProcessDocumentWithHSG extends ProcessDocument to include Tier 2 HSG processing
func (s *SemanticLinkService) ProcessDocumentWithHSG(ctx context.Context, doc models.SubjectContent) error {
	// TIER 1: Existing semantic link processing
	if err := s.ProcessDocument(ctx, doc); err != nil {
		return fmt.Errorf("tier 1 processing failed: %w", err)
	}

	// TIER 2: HSG Summary Chunk Creation
	log.Printf("TIER 2 HSG ⏳: Beginning summary chunk creation for document %s", doc.ID.Hex())

	if err := s.processTier2SummaryChunks(ctx, doc); err != nil {
		log.Printf("Warning: Tier 2 processing failed for doc %s: %v", doc.ID.Hex(), err)
		// Don't fail the entire process if Tier 2 fails
	} else {
		log.Printf("TIER 2 HSG ✅: Successfully created summary chunks for document %s", doc.ID.Hex())
	}

	return nil
}

// processTier2SummaryChunks creates Tier 2 summary chunks from related semantic links
// AND extracts new higher-level semantic links from the summaries
func (s *SemanticLinkService) processTier2SummaryChunks(ctx context.Context, doc models.SubjectContent) error {
	// 1. Query related semantic links from this document
	semanticLinksCollection := s.MongoClient.Database.Collection("semantic_links")

	filter := bson.M{"metadata.document_id": doc.ID.Hex()}
	cursor, err := semanticLinksCollection.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to find semantic links: %w", err)
	}
	defer cursor.Close(ctx)

	var semanticLinks []models.SemanticLink
	if err := cursor.All(ctx, &semanticLinks); err != nil {
		return fmt.Errorf("failed to decode semantic links: %w", err)
	}

	if len(semanticLinks) < 2 {
		log.Printf("Insufficient semantic links (%d) for summary chunk creation", len(semanticLinks))
		return nil
	}

	// 2. Group related semantic links by domain/topic
	linkGroups := s.groupSemanticLinksByTopic(semanticLinks)

	// 3. Create summary chunks for each group
	summaryChunksCollection := s.MongoClient.Database.Collection("summary_chunks")

	for topic, links := range linkGroups {
		if len(links) < 2 {
			continue // Skip groups with insufficient links
		}

		// 4. Generate summary text using LLaMA
		summaryText, confidence, err := s.synthesizeSummaryFromLinks(ctx, links, topic)
		if err != nil {
			log.Printf("Warning: Failed to synthesize summary for topic %s: %v", topic, err)
			continue
		}

		// 5. Generate vector for summary text
		vector, err := s.LlmClient.GenerateEmbedding(summaryText)
		if err != nil {
			log.Printf("Warning: Failed to vectorize summary for topic %s: %v", topic, err)
			continue
		}

		// 6. Create and store SummaryChunk
		linkIDs := make([]primitive.ObjectID, len(links))
		for i, link := range links {
			linkIDs[i] = link.ID
		}

		summaryChunk := models.SummaryChunk{
			ID:              primitive.NewObjectID(),
			SummaryText:     summaryText,
			Vector:          vector,
			Domain:          doc.Domain,
			SemanticLinkIDs: linkIDs,
			Confidence:      confidence,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			Metadata: map[string]interface{}{
				"source_document_id": doc.ID.Hex(),
				"topic":              topic,
				"link_count":         len(links),
				"vector_source":      "aws_bedrock_titan",
			},
		}

		// Store in MongoDB
		insertResult, err := summaryChunksCollection.InsertOne(ctx, summaryChunk)
		if err != nil {
			log.Printf("Warning: Failed to store summary chunk in MongoDB: %v", err)
			continue
		}

		// Store in Weaviate
		properties := map[string]interface{}{
			"summary_text":      summaryChunk.SummaryText,
			"domain":            summaryChunk.Domain,
			"confidence":        summaryChunk.Confidence,
			"semantic_link_ids": linkIDs,
			"metadata":          summaryChunk.Metadata,
		}

		err = db.StoreDocumentWithVector(ctx, "SummaryChunks", properties, vector)
		if err != nil {
			log.Printf("Warning: Failed to store summary chunk in Weaviate: %v", err)
		} else {
			log.Printf("✅ Created Tier 2 summary chunk for topic: %s", topic)
		}

		// 🆕 TIER 2 RECURSIVE STEP: Extract NEW higher-level semantic links from the summary
		log.Printf("🔄 TIER 2 RECURSIVE: Extracting higher-level semantic links from summary...")

		summaryChunkID, ok := insertResult.InsertedID.(primitive.ObjectID)
		if !ok {
			log.Printf("Warning: Could not extract SummaryChunk ID for recursive extraction")
			continue
		}

		newHighLevelLinks, err := s.extractHigherLevelLinksFromSummary(ctx, summaryText, summaryChunkID, topic)
		if err != nil {
			log.Printf("Warning: Failed to extract higher-level links from summary: %v", err)
			continue
		}

		if len(newHighLevelLinks) > 0 {
			// Store the new higher-level semantic links
			insertedCount, err := s.StoreSemanticLinks(ctx, semanticLinksCollection, newHighLevelLinks, doc.ID)
			if err != nil {
				log.Printf("Warning: Failed to store higher-level semantic links: %v", err)
			} else {
				log.Printf("✅ TIER 2 RECURSIVE: Extracted %d new higher-level semantic links from summary", insertedCount)
			}

			// Link these new semantic links back to the summary chunk
			s.linkHigherLevelLinksToSummary(ctx, summaryChunksCollection, summaryChunkID, newHighLevelLinks)
		}
	}

	return nil
}

// 🆕 extractHigherLevelLinksFromSummary extracts NEW semantic links from summary text
// These are higher-level relationships that may not be visible in granular content
func (s *SemanticLinkService) extractHigherLevelLinksFromSummary(
	ctx context.Context,
	summaryText string,
	summaryChunkID primitive.ObjectID,
	topic string,
) ([]models.SemanticLink, error) {
	// First, extract key terms from the summary
	terms, err := s.extractKeyTerms(ctx, summaryText)
	if err != nil {
		return nil, fmt.Errorf("failed to extract terms from summary: %w", err)
	}

	if len(terms) == 0 {
		log.Printf("No key terms extracted from summary")
		return []models.SemanticLink{}, nil
	}

	// Create Term objects (unverified for now, as these are higher-level concepts)
	var summaryTerms []models.Term
	for _, termText := range terms {
		summaryTerms = append(summaryTerms, models.Term{
			Text:     termText,
			Verified: false, // Higher-level terms may not exist in our base ontology
			MongoID:  primitive.NilObjectID,
		})
	}

	// Extract relationships with a special prompt for higher-level abstraction
	prompt := fmt.Sprintf(`Analyze this high-level summary and identify ABSTRACT semantic relationships.
Focus on:
- Overarching themes and patterns
- Causal chains that span multiple concepts
- Hierarchical relationships (broader categories)
- Meta-level connections not obvious in granular details

Summary (Topic: %s):
%s

Key terms identified: %s

Return ONLY higher-level, abstract relationships as JSON:
[
  {
    "source_term": "broader_concept",
    "target_term": "related_concept",
    "relation_type": "causes|enables|encompasses|relates_to",
    "context": "high-level description of the abstract relationship",
    "strength": 0.0-1.0
  }
]

Focus on relationships that provide NEW insights not present in the original content.`,
		topic, summaryText, strings.Join(terms, ", "))

	response, err := s.LlmClient.Generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("LLaMA generation failed for higher-level links: %w", err)
	}

	// Parse the response
	startIdx := strings.Index(response, "[")
	endIdx := strings.LastIndex(response, "]")
	if startIdx == -1 || endIdx == -1 {
		return []models.SemanticLink{}, nil // No relationships found
	}

	jsonArray := response[startIdx : endIdx+1]

	var rawRelationships []struct {
		SourceTerm   string  `json:"source_term"`
		TargetTerm   string  `json:"target_term"`
		RelationType string  `json:"relation_type"`
		Context      string  `json:"context"`
		Strength     float64 `json:"strength"`
	}

	if err := json.Unmarshal([]byte(jsonArray), &rawRelationships); err != nil {
		log.Printf("Warning: Failed to parse higher-level relationships: %v", err)
		return []models.SemanticLink{}, nil
	}

	// Convert to SemanticLink objects with special metadata
	var higherLevelLinks []models.SemanticLink
	for _, rel := range rawRelationships {
		// Only keep high-confidence abstract relationships
		if rel.Strength < 0.6 {
			continue
		}

		link := models.SemanticLink{
			ID:           primitive.NewObjectID(),
			SourceID:     primitive.NilObjectID, // Higher-level concepts may not map to existing entities
			TargetID:     primitive.NilObjectID,
			SourceTerm:   rel.SourceTerm,
			TargetTerm:   rel.TargetTerm,
			RelationType: rel.RelationType,
			Context:      rel.Context,
			Confidence:   rel.Strength,
			Domain:       topic,
			CreatedAt:    time.Now(),
			Metadata: map[string]interface{}{
				"abstraction_level": "tier2_summary",
				"derived_from":      "summary_chunk",
				"summary_chunk_id":  summaryChunkID.Hex(),
				"higher_level_link": true,
				"original_topic":    topic,
			},
		}

		higherLevelLinks = append(higherLevelLinks, link)
	}

	log.Printf("🔍 Extracted %d higher-level semantic links from summary (topic: %s)", len(higherLevelLinks), topic)
	return higherLevelLinks, nil
}

// 🆕 linkHigherLevelLinksToSummary creates bidirectional references
// between summary chunks and their derived higher-level semantic links
func (s *SemanticLinkService) linkHigherLevelLinksToSummary(
	ctx context.Context,
	summaryCollection *mongo.Collection,
	summaryChunkID primitive.ObjectID,
	higherLevelLinks []models.SemanticLink,
) {
	// Extract IDs of the higher-level links
	linkIDs := make([]primitive.ObjectID, len(higherLevelLinks))
	for i, link := range higherLevelLinks {
		linkIDs[i] = link.ID
	}

	// Update the summary chunk with references to derived links
	filter := bson.M{"_id": summaryChunkID}
	update := bson.M{
		"$set": bson.M{
			"derived_semantic_link_ids": linkIDs,
			"updated_at":                time.Now(),
		},
		"$inc": bson.M{
			"derived_link_count": len(linkIDs),
		},
	}

	_, err := summaryCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Warning: Failed to link higher-level links to summary %s: %v", summaryChunkID.Hex(), err)
	} else {
		log.Printf("✅ Linked %d higher-level semantic links to summary chunk %s", len(linkIDs), summaryChunkID.Hex())
	}
}

// StoreSemanticLinks stores semantic links in MongoDB with vectorization
func (s *SemanticLinkService) StoreSemanticLinks(
	ctx context.Context,
	collection *mongo.Collection,
	links []models.SemanticLink,
	documentID primitive.ObjectID,
) (int, error) {
	if len(links) == 0 {
		return 0, nil
	}

	// Prepare bulk write operations
	var ops []mongo.WriteModel

	for _, link := range links {
		// Normalize terms for consistent storage
		link.SourceTerm = s.normalizer.NormalizeText(link.SourceTerm)
		link.TargetTerm = s.normalizer.NormalizeText(link.TargetTerm)
		link.Context = s.normalizer.NormalizeForVectorization(link.Context)

		// Generate composite key AFTER normalization
		compositeKey := GenerateCompositeKey(link.SourceTerm, link.TargetTerm, link.RelationType)
		link.CompositeKey = compositeKey

		// Check if this link already exists
		existingLink := &models.SemanticLink{}
		err := collection.FindOne(ctx, bson.M{"composite_key": compositeKey}).Decode(existingLink)

		if err == nil {
			// Link already exists, skip or update
			log.Printf("⚠️ Semantic link already exists: %s → %s (%s). Skipping.",
				link.SourceTerm, link.TargetTerm, link.RelationType)
			continue
		} else if err != mongo.ErrNoDocuments {
			// Real error, not just "not found"
			log.Printf("❌ Error checking for duplicate link: %v", err)
			continue
		}

		// Link doesn't exist, proceed with insertion
		ops = append(ops, mongo.NewInsertOneModel().SetDocument(link))
	}

	// Execute bulk write
	if len(ops) > 0 {
		result, err := collection.BulkWrite(ctx, ops)
		if err != nil {
			return 0, fmt.Errorf("failed to store semantic links: %w", err)
		}
		return int(result.InsertedCount), nil
	}

	return 0, nil
}

// StoreSemanticLink stores a single semantic link (used by ProcessDocument)
func (s *SemanticLinkService) StoreSemanticLink(ctx context.Context, link models.SemanticLink) error {
	collection := s.MongoClient.Database.Collection("semantic_links")

	// Normalize terms
	link.SourceTerm = s.normalizer.NormalizeText(link.SourceTerm)
	link.TargetTerm = s.normalizer.NormalizeText(link.TargetTerm)
	link.Context = s.normalizer.NormalizeForVectorization(link.Context)

	// Generate composite key
	link.CompositeKey = GenerateCompositeKey(link.SourceTerm, link.TargetTerm, link.RelationType)

	// Check for duplicates
	existingLink := &models.SemanticLink{}
	err := collection.FindOne(ctx, bson.M{"composite_key": link.CompositeKey}).Decode(existingLink)

	if err == nil {
		log.Printf("⚠️ Semantic link already exists: %s → %s (%s). Skipping.",
			link.SourceTerm, link.TargetTerm, link.RelationType)
		return nil
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking for duplicate: %w", err)
	}

	// Insert new link
	_, err = collection.InsertOne(ctx, link)
	return err
}

// updateDocumentWithLinkInfo updates the original document with semantic link metadata
func (s *SemanticLinkService) updateDocumentWithLinkInfo(ctx context.Context, collection *mongo.Collection, docID primitive.ObjectID, linkCount int) {
	update := bson.M{
		"$set": bson.M{
			"semantic_links_count": linkCount,
			"updated_at":           time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": docID}, update)
	if err != nil {
		log.Printf("Warning: Failed to update document with link info: %v", err)
	}
}

// groupSemanticLinksByTopic groups semantic links by topic/domain for summary creation
func (s *SemanticLinkService) groupSemanticLinksByTopic(links []models.SemanticLink) map[string][]models.SemanticLink {
	groups := make(map[string][]models.SemanticLink)

	for _, link := range links {
		topic := link.Domain
		if topic == "" {
			topic = "general"
		}
		groups[topic] = append(groups[topic], link)
	}

	return groups
}

// synthesizeSummaryFromLinks generates a summary from a group of semantic links
func (s *SemanticLinkService) synthesizeSummaryFromLinks(ctx context.Context, links []models.SemanticLink, topic string) (string, float64, error) {
	// Build context from links
	var linkDescriptions []string
	for _, link := range links {
		desc := fmt.Sprintf("%s %s %s", link.SourceTerm, link.RelationType, link.TargetTerm)
		linkDescriptions = append(linkDescriptions, desc)
	}

	prompt := fmt.Sprintf(`Synthesize a coherent summary from these semantic relationships about %s:

%s

Create a brief, informative summary (2-3 sentences) that captures the key insights.
Return ONLY the summary text, no JSON.`, topic, strings.Join(linkDescriptions, "\n"))

	summary, err := s.LlmClient.Generate(ctx, prompt)
	if err != nil {
		return "", 0.0, err
	}

	// Calculate average confidence
	var totalConfidence float64
	for _, link := range links {
		totalConfidence += link.Confidence
	}
	avgConfidence := totalConfidence / float64(len(links))

	return strings.TrimSpace(summary), avgConfidence, nil
}

// CorrectMetricsWithFeedback sends validation feedback to LLM for correction
// func (s *SemanticLinkService) CorrectMetricsWithFeedback(
// 	ctx context.Context,
// 	link *models.SemanticLink,
// 	initialMetrics HierarchyMetrics,
// 	validationReport ValidationReport,
// ) (*HierarchyMetrics, error) {
// 	prompt := fmt.Sprintf(`Your initial hierarchy metrics were incorrect. Please correct them based on this feedback:

// Relationship: %s %s %s
// Context: %s

// Your initial metrics:
// - Source Term Generality: %.2f
// - Target Term Generality: %.2f
// - Semantic Distance: %.2f
// - Relationship Strength: %.2f

// Validation errors:
// %s

// Please provide corrected metrics following the same JSON format:
// {
//   "source_term_generality": 0.0-1.0,
//   "target_term_generality": 0.0-1.0,
//   "semantic_distance": 0.0-1.0,
//   "relationship_strength": 0.0-1.0,
//   "explanation": "brief explanation of corrections"
// }`,
// 		link.SourceTerm, link.RelationType, link.TargetTerm, link.Context,
// 		initialMetrics.SourceTermGenerality,
// 		initialMetrics.TargetTermGenerality,
// 		initialMetrics.SemanticDistance,
// 		initialMetrics.RelationshipStrength,
// 		validationReport.ErrorMessage)

// 	response, err := s.LlmClient.Generate(ctx, prompt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var correctedMetrics HierarchyMetrics
// 	if err := json.Unmarshal([]byte(cleanJSONResponse(response)), &correctedMetrics); err != nil {
// 		return nil, fmt.Errorf("failed to parse corrected metrics: %w", err)
// 	}

// 	return &correctedMetrics, nil
// }

// GenerateCompositeKey creates a normalized composite key (fixed duplicate definitions)
func GenerateCompositeKey(source, target, relationType string) string {
	return fmt.Sprintf("%s::%s::%s",
		strings.ToLower(strings.TrimSpace(source)),
		strings.ToLower(strings.TrimSpace(target)),
		strings.ToLower(strings.TrimSpace(relationType)),
	)
}

// 🆕 NEW: Sample-driven extraction prompt builder
func (s *SemanticLinkService) buildExtractionPromptWithSamples(content string, domain string) string {
	// Load domain-specific samples
	samples := s.loadSamplesByDomain(domain)

	samplesText := ""
	for i, sample := range samples {
		samplesText += fmt.Sprintf(`
EXAMPLE %d:
Input: %s
Output: %s
Quality Score: %.2f
`, i+1, sample.InputText, sample.ExpectedOutput, sample.QualityScore)
	}

	return fmt.Sprintf(`Extract ATOMIC semantic relationships following these CURATED EXAMPLES:

%s

Now extract relationships from this NEW text using the SAME pattern and quality standards:

Text: %s

Return ONLY a JSON array matching the example format.`, samplesText, content)
}

// 🆕 NEW: Sample management
type ExtractionSample struct {
	InputText      string    `json:"input_text"`
	ExpectedOutput string    `json:"expected_output"`
	QualityScore   float64   `json:"quality_score"`
	Domain         string    `json:"domain"`
	CreatedAt      time.Time `json:"created_at"`
}

func (s *SemanticLinkService) loadSamplesByDomain(domain string) []ExtractionSample {
	// Try loading from MongoDB first
	ctx := context.Background()
	samplesCollection := s.MongoClient.Database.Collection("extraction_samples")

	filter := bson.M{
		"domain":        domain,
		"quality_score": bson.M{"$gte": 0.8}, // Only high-quality samples
	}

	cursor, err := samplesCollection.Find(ctx, filter, nil)
	if err != nil {
		log.Printf("Warning: Failed to load samples from MongoDB: %v", err)
		return s.getDefaultSamples(domain)
	}
	defer cursor.Close(ctx)

	var samples []ExtractionSample
	if err := cursor.All(ctx, &samples); err != nil {
		log.Printf("Warning: Failed to decode samples: %v", err)
		return s.getDefaultSamples(domain)
	}

	// Limit to top 3-5 samples for context window efficiency
	if len(samples) > 5 {
		samples = samples[:5]
	}

	return samples
}

func (s *SemanticLinkService) getDefaultSamples(domain string) []ExtractionSample {
	// Fallback hardcoded samples by domain
	immunologySamples := []ExtractionSample{
		{
			InputText: "*XLA patient 9 years old with absent B cells and low IgG",
			ExpectedOutput: `{
  "source_term": "XLA patient 9 years old",
  "target_term": "absent B cells",
  "relation_type": "has_clinical_finding",
  "context": "Primary immunodeficiency presentation",
  "is_case_specific": true,
  "case_study_id": "XLA_Case_1",
  "confidence": 0.95
}`,
			QualityScore: 0.95,
			Domain:       "immunology",
		},
	}

	switch domain {
	case "immunology":
		return immunologySamples
	default:
		return []ExtractionSample{}
	}
}

// 🆕 NEW: Quality assessment and sample curation
func (s *SemanticLinkService) assessExtractionQuality(
	ctx context.Context,
	extractedLink models.SemanticLink,
	originalText string,
) (float64, error) {
	prompt := fmt.Sprintf(`Assess the quality of this semantic link extraction:

Original Text: %s

Extracted Link:
- Source: %s
- Target: %s
- Relation: %s
- Context: %s

Rate from 0.0-1.0 based on:
1. Accuracy (is the relationship correct?)
2. Completeness (captures all nuance?)
3. Atomicity (single clear relationship?)
4. Context quality (minimal but sufficient?)

Return ONLY a JSON object:
{
  "quality_score": 0.0-1.0,
  "reasoning": "brief explanation",
  "suggestions": ["improvement 1", "improvement 2"]
}`, originalText, extractedLink.SourceTerm, extractedLink.TargetTerm,
		extractedLink.RelationType, extractedLink.Context)

	response, err := s.LlmClient.Generate(ctx, prompt)
	if err != nil {
		return 0.0, err
	}

	var assessment struct {
		QualityScore float64  `json:"quality_score"`
		Reasoning    string   `json:"reasoning"`
		Suggestions  []string `json:"suggestions"`
	}

	if err := json.Unmarshal([]byte(cleanJSONResponse(response)), &assessment); err != nil {
		return 0.0, err
	}

	return assessment.QualityScore, nil
}

// 🆕 NEW: Automatic sample curation from high-quality extractions
func (s *SemanticLinkService) CurateHighQualityExtractions(ctx context.Context) error {
	semanticLinksCollection := s.MongoClient.Database.Collection("semantic_links")
	samplesCollection := s.MongoClient.Database.Collection("extraction_samples")

	// Find recent high-quality extractions
	filter := bson.M{
		"confidence": bson.M{"$gte": 0.9},
		"created_at": bson.M{"$gte": time.Now().AddDate(0, 0, -7)}, // Last 7 days
	}

	cursor, err := semanticLinksCollection.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to find high-quality links: %w", err)
	}
	defer cursor.Close(ctx)

	var highQualityLinks []models.SemanticLink
	if err := cursor.All(ctx, &highQualityLinks); err != nil {
		return err
	}

	// Convert to samples and store
	for _, link := range highQualityLinks {
		sample := ExtractionSample{
			InputText: link.ContextOriginal, // Use original text with symbols
			ExpectedOutput: fmt.Sprintf(`{
  "source_term": "%s",
  "target_term": "%s",
  "relation_type": "%s",
  "confidence": %.2f
}`, link.SourceTermOriginal, link.TargetTermOriginal, link.RelationType, link.Confidence),
			QualityScore: link.Confidence,
			Domain:       link.Domain,
			CreatedAt:    time.Now(),
		}

		// Store with deduplication
		_, err := samplesCollection.UpdateOne(
			ctx,
			bson.M{
				"input_text": sample.InputText,
				"domain":     sample.Domain,
			},
			bson.M{"$set": sample},
			// Use upsert to avoid duplicates
			// options.Update().SetUpsert(true), // Uncomment if you have mongo options imported
		)

		if err != nil {
			log.Printf("Warning: Failed to store sample: %v", err)
		}
	}

	log.Printf("✅ Curated %d new extraction samples", len(highQualityLinks))
	return nil
}

// CreateSemanticLinks processes multiple documents and creates semantic links for them
func (s *SemanticLinkService) CreateSemanticLinks(
	ctx context.Context,
	documents []models.SubjectContent,
	domain string,
	batchID string,
) error {
	if len(documents) == 0 {
		log.Printf("[Batch %s] No documents to process for semantic links", batchID)
		return nil
	}

	log.Printf("[Batch %s] Starting semantic link creation for %d documents in domain: %s", batchID, len(documents), domain)

	// Process each document
	var processedCount int
	var failedCount int

	for _, doc := range documents {
		// Process the document to extract semantic links
		if err := s.ProcessDocument(ctx, doc); err != nil {
			log.Printf("[Batch %s] ⚠️ Failed to process document %s (%s): %v", batchID, doc.ID.Hex(), doc.Title, err)
			failedCount++
			continue
		}

		processedCount++
		log.Printf("[Batch %s] ✅ Processed document: %s", batchID, doc.Title)
	}

	log.Printf("[Batch %s] Semantic link creation complete. Processed: %d, Failed: %d", batchID, processedCount, failedCount)

	// If all documents failed, return an error
	if processedCount == 0 && failedCount > 0 {
		return fmt.Errorf("failed to process any documents for semantic links")
	}

	return nil
}
