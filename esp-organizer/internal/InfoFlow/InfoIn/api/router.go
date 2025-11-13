package api

import (
	"encoding/json"
	"net/http"
	"time"

	"log"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(r *mux.Router) {
	// Add a debug log to confirm routes are registered
	log.Println("Registering API routes...")

	// --- All Application Routes Defined Here ---
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ESP Organizer API is running"))
	}).Methods("GET")

	apiRouter := r.PathPrefix("/api").Subrouter()

	// Source management
	apiRouter.HandleFunc("/sources", GetSourcesHandler).Methods(http.MethodGet)
	apiRouter.HandleFunc("/sources", CreateSourceHandler).Methods(http.MethodPost)

	// Domain-aware search
	//apiRouter.HandleFunc("/search/domain", DomainSearchHandler).Methods(http.MethodGet, http.MethodPost)
	apiRouter.HandleFunc("/search/tags", SearchByTagsHandler).Methods(http.MethodGet)
	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"service":   "esp-organizer-api",
			"timestamp": time.Now(),
		})
	}).Methods("GET")

	// Immunology endpoints
	apiRouter.HandleFunc("/immunology/upload-chapter", ImmunologyChapterUploadHandler).Methods(http.MethodPost)
	apiRouter.HandleFunc("/immunology/search", SearchImmunologyContent).Methods(http.MethodGet)
	apiRouter.HandleFunc("/immunology/search", ImmunologySearchHandler).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/immunology/hsg-search", HSGSearchHandler).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.HandleFunc("/immunology/chapters", GetImmunologyChapters).Methods(http.MethodGet)
	apiRouter.HandleFunc("/immunology/terms", GetImmunologyTerms).Methods(http.MethodGet)
	// CHISG integration
	// Add to RegisterRoutes()
	apiRouter.HandleFunc("/coach/respond", CoachRespondHandler).Methods(http.MethodPost)
	// Extraction status
	apiRouter.HandleFunc("/extraction/status/{jobId}", ExtractionStatusHandler).Methods(http.MethodGet)
	apiRouter.HandleFunc("/upload/immunology/chapter", ImmunologyChapterUploadHandler).Methods(http.MethodPost)
	apiRouter.HandleFunc("/upload/document", DocumentUploadHandler).Methods(http.MethodPost)

	// Study area specific routes
	apiRouter.HandleFunc("/gcse/search", GCSESearchHandler).Methods(http.MethodGet, http.MethodPost)
	apiRouter.HandleFunc("/civil-engineering/search", CivilEngineeringSearchHandler).Methods(http.MethodGet, http.MethodPost)

	// AI Chat
	apiRouter.HandleFunc("/ai/chat", AIChatHandler).Methods(http.MethodPost)
	apiRouter.HandleFunc("/batches", GetProcessingBatchesHandler).Methods(http.MethodGet)

	apiRouter.HandleFunc("/diagnostics/config", DiagnosticsConfigHandler).Methods(http.MethodGet)
	apiRouter.HandleFunc("/diagnostics/document", DiagnosticsDocumentHandler).Methods(http.MethodGet)
	apiRouter.HandleFunc("/diagnostics/collections", DiagnosticsCollectionsHandler).Methods(http.MethodGet)

	// Add a simple test endpoint that's guaranteed to work
	r.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "API is running"})
	}).Methods("GET")

	// Health check
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"service":   "esp-organizer-api",
			"timestamp": time.Now(),
		})
	}).Methods("GET")

	log.Println("API routes registered successfully. 2")
}

// RegisterAdditionalRoutes adds potentially missing routes
func RegisterAdditionalRoutes(r *mux.Router) {
	log.Println("Registering additional diagnostic routes...")

	r.HandleFunc("/api/diagnostics/config", DiagnosticsConfigHandler).Methods("GET")
	r.HandleFunc("/api/diagnostics/document", DiagnosticsDocumentHandler).Methods("GET")
	r.HandleFunc("/api/diagnostics/collections", DiagnosticsCollectionsHandler).Methods("GET")

	log.Println("Additional diagnostic routes registered")
}
