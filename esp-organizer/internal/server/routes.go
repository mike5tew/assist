package server

import (
	"esp-organizer/internal/InfoFlow/InfoIn/api"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *mux.Router) {
	// API routes
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Extraction routes
	extractionRouter := apiRouter.PathPrefix("/extraction").Subrouter()
	extractionRouter.HandleFunc("/status/{jobId}", api.ExtractionStatusHandler).Methods("GET")

	// Make sure we also handle the route with a different parameter format for compatibility
	extractionRouter.HandleFunc("/status/{jobID}", api.ExtractionStatusHandler).Methods("GET")

	// Also register a version that accepts the job ID as a query parameter
	extractionRouter.HandleFunc("/status", api.ExtractionStatusHandler).Methods("GET")

	// Log the routes we're registering
	log.Printf("Registered route: /api/extraction/status/{jobId}")
	log.Printf("Registered route: /api/extraction/status/{jobID}")
	log.Printf("Registered route: /api/extraction/status?jobId=...")

	// Search routes
	apiRouter.HandleFunc("/immunology/search", api.ImmunologySearchHandler).Methods(http.MethodPost, http.MethodOptions)

}
