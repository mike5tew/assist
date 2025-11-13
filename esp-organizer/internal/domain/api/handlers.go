// Functioning - CORS issue resolved. Do not edit unless necessary.

package api

import (
	"net/http"
	// Add this import
	// Import the ChapterInfo type for chapter metadata
)

// Handler for the root endpoint
func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the ESP Organizer API"))
}

// Handler for getting all materials
func GetMaterialsHandler(w http.ResponseWriter, r *http.Request) {
	// Logic to retrieve materials from the database
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List of all materials"))
}

// Handler for creating a new material
func CreateMaterialHandler(w http.ResponseWriter, r *http.Request) {
	// Logic to create a new material in the database
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Material created"))
}

// SemanticSkillQueryHandler processes semantic queries against the skills database
// func SemanticSkillQueryHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	mongDB, err := mongoconnect.MongoDBConnect()
// 	if err != nil {
// 		http.Error(w, "Failed to connect to MongoDB: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Handle preflight requests
// 	if r.Method == "OPTIONS" {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	var req models.SemanticQueryRequest

// 	if r.Method == "GET" {
// 		// Support GET with query parameter for easy testing
// 		req.Query = r.URL.Query().Get("q")
// 		if req.Query == "" {
// 			http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
// 			return
// 		}
// 	} else {
// 		// Handle POST with JSON body
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 	}

// 	if req.Query == "" {
// 		http.Error(w, "Query text cannot be empty", http.StatusBadRequest)
// 		return
// 	}

// 	// Create services
// 	skillService := skills.NewSkillService(mongDB.Database, os.Getenv("SKILLS_COLLECTION"))
// 	llmClient := llm.NewLlamaClient()
// 	queryHandler := query.NewQueryHandler(skillService, llmClient)

// 	// Process the semantic query, now passing filters
// 	response, err := queryHandler.ProcessQuery(r.Context(), req.Query, req.Filters, "")
// 	if err != nil {
// 		http.Error(w, "Query processing failed: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Enhanced response for LLM understanding
// 	enhancedResponse := map[string]interface{}{
// 		"query":     req.Query,
// 		"timestamp": response.Timestamp,
// 		"results_summary": map[string]interface{}{
// 			"total_exact_matches":    len(response.MongoResults),
// 			"total_semantic_matches": len(response.WeaviateResults),
// 			"total_related_skills":   len(response.RelatedSkills),
// 		},
// 		"exact_matches":     ai.FormatSkillsForLLM(response.MongoResults),
// 		"semantic_matches":  ai.FormatWeaviateResultsForLLM(response.WeaviateResults),
// 		"related_skills":    ai.FormatSkillsForLLM(response.RelatedSkills),
// 		"llm_synthesis":     response.Synthesis,
// 		"semantic_insights": ai.GenerateSemanticInsights(response),
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(enhancedResponse)
// }
