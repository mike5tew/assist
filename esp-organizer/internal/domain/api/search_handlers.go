package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"esp-organizer/internal/InfoFlow/InfoIn"
	"esp-organizer/internal/InfoFlow/InfoOut/llm"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"esp-organizer/internal/InfoFlow/InfoStore/skills"
	"esp-organizer/internal/factory"
	"esp-organizer/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

// DomainConfig defines the configuration for each search domain
type DomainConfig struct {
	Name              string
	MongoCollections  []string
	WeaviateClass     string
	RequiresAI        bool
	Implemented       bool
	UseSubjectContent bool
}

// Domain configurations
var domainConfigs = map[string]DomainConfig{
	"immunology": {
		Name:              "immunology",
		MongoCollections:  []string{"immunology_content", "medical_terms"},
		WeaviateClass:     "MedicalExcerpt",
		RequiresAI:        true,
		Implemented:       true, // <--- MUST BE TRUE
		UseSubjectContent: true,
	},
	"gcse": {
		Name:              "gcse",
		MongoCollections:  []string{"gcse_content"},
		WeaviateClass:     "SemanticLinks",
		RequiresAI:        false,
		Implemented:       false,
		UseSubjectContent: true,
	},
	"civil-engineering": {
		Name:              "civil-engineering",
		MongoCollections:  []string{"engineering_content"},
		WeaviateClass:     "SemanticLinks",
		RequiresAI:        false,
		Implemented:       false,
		UseSubjectContent: true,
	},
	"general": {
		Name:              "general",
		MongoCollections:  nil, // Explicitly nil if no collections are used
		WeaviateClass:     "SemanticLinks",
		RequiresAI:        false,
		Implemented:       true,
		UseSubjectContent: false, // Assuming general search doesn't filter by subject content
	},
}

// SearchImmunologyContent searches immunology content and terms
func SearchImmunologyContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(r.Context())

	// Search in both immunology collections
	contentCollection := mongoDb.Database.Collection("immunology_content")
	termsCollection := mongoDb.Database.Collection("immunology_terms")

	var results map[string]interface{} = make(map[string]interface{})

	// Search content (chapters, case studies)
	contentFilter := bson.M{
		"domain": "immunology",
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"content": bson.M{"$regex": query, "$options": "i"}},
			{"tags": bson.M{"$in": []string{query}}},
		},
	}

	contentCursor, err := contentCollection.Find(r.Context(), contentFilter)
	if err == nil {
		var contentResults []bson.M
		contentCursor.All(r.Context(), &contentResults)
		results["content"] = contentResults
		contentCursor.Close(r.Context())
	}

	// Search terms
	termsFilter := bson.M{
		"domain": "immunology",
		"$or": []bson.M{
			{"term": bson.M{"$regex": query, "$options": "i"}},
			{"definition": bson.M{"$regex": query, "$options": "i"}},
			{"category": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	termsCursor, err := termsCollection.Find(r.Context(), termsFilter)
	if err == nil {
		var termsResults []bson.M
		termsCursor.All(r.Context(), &termsResults)
		results["terms"] = termsResults
		termsCursor.Close(r.Context())
	}

	results["domain"] = "immunology"
	results["query"] = query
	results["timestamp"] = time.Now()

	json.NewEncoder(w).Encode(results)
}

// GetImmunologyChapters returns all immunology chapters
func GetImmunologyChapters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(r.Context())

	collection := mongoDb.Database.Collection("immunology_content")

	cursor, err := collection.Find(r.Context(), bson.M{
		"domain":       "immunology",
		"content_type": "chapter",
	})
	if err != nil {
		http.Error(w, "Failed to query chapters", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var chapters []bson.M
	cursor.All(r.Context(), &chapters)

	response := map[string]interface{}{
		"domain":    "immunology",
		"chapters":  chapters,
		"count":     len(chapters),
		"timestamp": time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

// GetImmunologyTerms returns all immunology medical terms
func GetImmunologyTerms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(r.Context())

	collection := mongoDb.Database.Collection("immunology_terms")

	cursor, err := collection.Find(r.Context(), bson.M{"domain": "immunology"})
	if err != nil {
		http.Error(w, "Failed to query terms", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var terms []bson.M
	cursor.All(r.Context(), &terms)

	response := map[string]interface{}{
		"domain":    "immunology",
		"terms":     terms,
		"count":     len(terms),
		"timestamp": time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

// GCSESearchHandler is a placeholder for the GCSE-specific search endpoint.
func GCSESearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "Not Implemented",
		"message": "GCSE-specific search is not yet implemented.",
	})
}

// CivilEngineeringSearchHandler is a placeholder for the civil engineering-specific search endpoint.
func CivilEngineeringSearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "Not Implemented",
		"message": "Civil Engineering-specific search is not yet implemented.",
	})
}

// Extracts the subject value from the top-level operands of the Filter.
func extractSubjectFilter(filters *models.Filter) string {
	if filters == nil || filters.Operands == nil {
		return ""
	}

	// Assuming filters.Operands is a slice of your Operand struct.
	for _, op := range filters.Operands {

		// 1. Check if the Path is specifically ["subject"]
		if len(op.Path) == 1 && op.Path[0] == "subject" {

			// 2. Check for the "Equal" operator (optional, but good practice for filtering)
			// You can remove this check if you assume any path match is the subject filter.
			if op.Operator == "Equal" || op.Operator == "" {

				// 3. Return the string value (assuming subject is always a string)
				return op.ValueString
			}
		}
	}
	return ""
}

// handleDomainSearch processes searches for any domain based on configuration
func handleDomainSearch(w http.ResponseWriter, r *http.Request, queryText string, domain string, filters *models.Filter) {
	log.Printf("2 handleDomainSearch called for domain: %s, query: %s", domain, queryText)

	config, exists := domainConfigs[domain]
	if !exists {
		// Fall back to general search if domain not found
		log.Printf("Domain '%s' not found, falling back to general domain", domain)
		config = domainConfigs["general"]
	} else {
		log.Printf("Found configuration for domain: %s, implemented: %v", domain, config.Implemented)
	}

	if !config.Implemented {
		// Return a "coming soon" response for unimplemented domains
		log.Printf("Domain '%s' is not yet implemented", domain)
		response := map[string]interface{}{
			"domain":  config.Name,
			"query":   queryText,
			"message": fmt.Sprintf("%s content search will be implemented soon", config.Name),
			"status":  "coming_soon",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Load environment variables
	log.Printf("Initializing MongoDB connection for domain: %s", domain)
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("Database initialization failed: %v", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(r.Context()) // Ensure connection is closed

	collections, _ := mongoDb.Database.ListCollectionNames(context.Background(), bson.M{})
	log.Printf("Available collections: %v", collections)
	log.Println("MongoDB connected successfully")

	// gather the skills collection name from the environment
	skillsCollectionName := os.Getenv("SKILLS_COLLECTION")
	if skillsCollectionName == "" {
		skillsCollectionName = "skills" // Default collection name
	}
	log.Printf("Using skills collection: %s", skillsCollectionName)

	// Initialize the skills service with the specified collection
	log.Printf("Initializing services for domain: %s", domain)
	llamaClient := llm.NewLlamaClient()
	weaviateClient := db.GetWeaviateClient() // Assuming weaviate is initialized
	if weaviateClient == nil {
		log.Printf("WARNING: Weaviate client is nil for domain: %s", domain)
	}

	skillService := skills.NewSkillService(
		mongoDb.Database.Collection(skillsCollectionName),
		weaviateClient,
		llamaClient,
	)

	// Prepare response object with common fields
	response := map[string]interface{}{
		"domain":    config.Name,
		"query":     queryText,
		"timestamp": time.Now(),
	}

	// Add AI-generated content if required by the domain
	if config.RequiresAI {
		log.Printf("Domain '%s' requires AI content generation", domain)
		if domain == "immunology" {
			log.Printf("Using HSG query service for domain '%s' query: %s", domain, queryText)

			// Initialize HSG query service instead of using semanticSearchMedicalContent
			hsgService, err := InfoIn.NewHSGQueryService()
			if err != nil {
				log.Printf("Error initializing HSG query service: %v", err)
				response["medical_content_error"] = "Failed to initialize HSG service: " + err.Error()
			} else {
				// Use HSG query service with domain and query
				ragContext, answer, err := hsgService.QueryHSG(r.Context(), queryText, domain, 5)
				if err != nil {
					log.Printf("HSG query failed: %v", err)
					response["medical_content_error"] = "HSG query failed: " + err.Error()
				} else {
					log.Printf("HSG query successful, retrieved %d summary chunks",
						len(ragContext.SummaryContext))

					// Format response with the HSG results
					response["medical_content"] = map[string]interface{}{
						"answer":      answer,
						"context":     ragContext,
						"query":       queryText,
						"search_type": "hierarchical_semantic_graph",
						"timestamp":   time.Now(),
					}
				}
			}
		}
		// Add other domain-specific AI handlers as needed
	}

	// Search in Weaviate for semantic matches
	log.Printf("Starting semantic search with class '%s' for query: %s", config.WeaviateClass, queryText)
	queryHandler := factory.NewQueryProcessor(skillService, llamaClient)
	// Use the domain-specific Weaviate class for search
	semanticResponse, err := queryHandler.ProcessQuery(r.Context(), queryText, filters, config.WeaviateClass)
	if err != nil {
		log.Printf("Error in semantic search for %s: %v", config.WeaviateClass, err)
		response["semantic_error"] = err.Error()
	} else {
		resultCount := 0
		if semanticResponse != nil && semanticResponse.WeaviateResults != nil {
			resultCount = len(semanticResponse.WeaviateResults)
		}
		log.Printf("Semantic search completed successfully with %d results", resultCount)
		if config.UseSubjectContent {
			response["semantic_matches"] = semanticResponse
		} else {
			response["results"] = semanticResponse
		}
	}

	// If domain has MongoDB collections configured, search them
	if len(config.MongoCollections) > 0 && mongoDb.Database != nil {
		log.Printf("Searching MongoDB collections for domain '%s': %v", domain, config.MongoCollections)
		mongoResults := make(map[string]interface{})

		// --- CORRECTED LOGIC START ---
		// Extract the subject filter value from the provided models.Filter object
		subjectFilterValue := extractSubjectFilter(filters)
		log.Printf("Extracted subject filter value: %s", subjectFilterValue)

		for _, collectionName := range config.MongoCollections {
			log.Printf("Searching collection: %s", collectionName)
			collection := mongoDb.Database.Collection(collectionName)

			// Base filter for text search across multiple fields
			filter := bson.M{
				"$or": []bson.M{
					{"content": bson.M{"$regex": queryText, "$options": "i"}},
					{"chapter_title": bson.M{"$regex": queryText, "$options": "i"}},
					{"term": bson.M{"$regex": queryText, "$options": "i"}},
				},
			}

			// Add subject filter if the value was successfully extracted
			if subjectFilterValue != "" {
				filter["subject"] = subjectFilterValue
			}
			// --- CORRECTED LOGIC END ---

			cursor, err := collection.Find(r.Context(), filter)
			if err != nil {
				log.Printf("Error searching collection %s: %v", collectionName, err)
				continue
			}

			var results []bson.M
			if err := cursor.All(r.Context(), &results); err != nil {
				log.Printf("Error retrieving results from cursor for collection %s: %v", collectionName, err)
				cursor.Close(r.Context())
				continue
			}

			log.Printf("Found %d results in collection %s", len(results), collectionName)
			if len(results) > 0 {
				mongoResults[collectionName] = results
			}
			cursor.Close(r.Context())
		}

		if len(mongoResults) > 0 {
			log.Printf("Adding MongoDB results to response")
			response["collection_results"] = mongoResults
		} else {
			log.Printf("No MongoDB results found")
		}
	}

	log.Printf("Sending response for domain '%s' search", domain)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// SemanticSkillQueryHandler handles semantic searches against the skills database
// This is maintained for backwards compatibility with existing clients
func SemanticSkillQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var queryText string
	var filters *models.Filter

	if r.Method == "GET" {
		queryText = r.URL.Query().Get("q")
	} else {
		var request struct {
			Query   string         `json:"query"`
			Filters *models.Filter `json:"filters"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		queryText = request.Query
		filters = request.Filters
	}

	if queryText == "" {
		http.Error(w, "Query text cannot be empty", http.StatusBadRequest)
		return
	}

	// Redirect to the general domain search handler
	handleDomainSearch(w, r, queryText, "general", filters)
}

// GET /api/immunology/case-studies
func GetCaseStudiesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: Replace with real DB query
	json.NewEncoder(w).Encode(map[string]interface{}{
		"case_studies": []interface{}{},
	})
}

// GET /api/immunology/medical-terms
func GetMedicalTermsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: Replace with real DB query
	json.NewEncoder(w).Encode(map[string]interface{}{
		"medical_terms": []interface{}{},
	})
}
