package api

import (
	"encoding/json"
	"esp-organizer/internal/domain/coach"
	"net/http"
	"time"

	"log"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all API routes and accepts the coach service
func RegisterRoutes(r *mux.Router, coachService coach.CoachServiceMVP) {
	log.Println("Registering API routes with MVP Coach service...")

	// Add CORS middleware
	r.Use(corsMiddleware)

	// Basic routes
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ESP Organizer API is running"))
	}).Methods("GET")

	apiRouter := r.PathPrefix("/api").Subrouter()

	// Health checks
	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"service":   "esp-organizer-api",
			"timestamp": time.Now(),
		})
	}).Methods("GET")

	apiRouter.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "API is running"})
	}).Methods("GET")

	// MVP Coach endpoints
	coachHandler := NewCoachHandlerMVP(coachService)
	apiRouter.HandleFunc("/coach/mvp-demo", coachHandler.CoachMVPDemoHandler).Methods("POST", "OPTIONS")
	log.Println("✅ Coach endpoints registered")

	// Full Diagnostic Coach endpoint
	apiRouter.HandleFunc("/coach/respond", FullDiagnosticCoachHandler).Methods("POST", "OPTIONS")
	log.Println("✅ Full diagnostic coach endpoint registered")

	// Skills/Semantic endpoints (MVP stubs)
	apiRouter.HandleFunc("/skills/semantic-query", SemanticQueryHandler).Methods("POST", "OPTIONS")
	log.Println("✅ Semantic query endpoint registered")

	// Semantic Links API (extract/search/validate/export)
	apiRouter.HandleFunc("/semantic-links/extract", SemanticLinksExtractHandler).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/semantic-links/search", SemanticLinksSearchHandler).Methods("GET", "POST", "OPTIONS")
	apiRouter.HandleFunc("/semantic-links/validate", SemanticLinksValidateHandler).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/semantic-links/export", SemanticLinksExportHandler).Methods("GET", "OPTIONS")
	log.Println("✅ Semantic-links endpoints registered (extract/search/validate/export)")

	// Link Types API (vocabulary for relationships)
	apiRouter.HandleFunc("/link-types", LinkTypesListHandler).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/link-types", LinkTypesCreateHandler).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/link-types/seed", LinkTypesSeedHandler).Methods("POST", "OPTIONS")
	log.Println("✅ Link-types endpoints registered (list/create/seed)")

	// Existing routes (keep all your current routes)
	apiRouter.HandleFunc("/sources", GetSourcesHandler).Methods("GET")
	apiRouter.HandleFunc("/sources", CreateSourceHandler).Methods("POST")
	apiRouter.HandleFunc("/search/tags", SearchByTagsHandler).Methods("GET")
	apiRouter.HandleFunc("/immunology/upload-chapter", ImmunologyChapterUploadHandler).Methods("POST")
	apiRouter.HandleFunc("/immunology/search", SearchImmunologyContent).Methods("GET")
	apiRouter.HandleFunc("/immunology/search", ImmunologySearchHandler).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/immunology/hsg-search", HSGSearchHandler).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/immunology/chapters", GetImmunologyChapters).Methods("GET")
	apiRouter.HandleFunc("/immunology/terms", GetImmunologyTerms).Methods("GET")
	apiRouter.HandleFunc("/extraction/status/{jobId}", ExtractionStatusHandler).Methods("GET")
	apiRouter.HandleFunc("/upload/immunology/chapter", ImmunologyChapterUploadHandler).Methods("POST")
	apiRouter.HandleFunc("/upload/document", DocumentUploadHandler).Methods("POST")
	apiRouter.HandleFunc("/gcse/search", GCSESearchHandler).Methods("GET", "POST")
	apiRouter.HandleFunc("/civil-engineering/search", CivilEngineeringSearchHandler).Methods("GET", "POST")
	apiRouter.HandleFunc("/ai/chat", AIChatHandler).Methods("POST")
	apiRouter.HandleFunc("/batches", GetProcessingBatchesHandler).Methods("GET")
	apiRouter.HandleFunc("/diagnostics/config", DiagnosticsConfigHandler).Methods("GET")
	apiRouter.HandleFunc("/diagnostics/document", DiagnosticsDocumentHandler).Methods("GET")
	apiRouter.HandleFunc("/diagnostics/collections", DiagnosticsCollectionsHandler).Methods("GET")

	log.Println("API routes registered successfully with MVP Coach integration")
}

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
