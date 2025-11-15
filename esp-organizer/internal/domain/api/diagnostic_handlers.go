package api

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/store/db"
	"esp-organizer/internal/utils"

	"fmt"
	"net/http"
	"time"
)

// DiagnosticHandler handles requests for system diagnostics
func DiagnosticHandler(w http.ResponseWriter, r *http.Request) {
	requestID := fmt.Sprintf("diag-%d", time.Now().UnixNano())
	utils.Log(requestID, "Starting system diagnostics")

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	diagnostics := map[string]interface{}{
		"timestamp": time.Now(),
		"requestID": requestID,
	}

	// Run MongoDB diagnostics
	dbDiagnostics, err := performDatabaseDiagnostics(ctx, requestID)
	if err != nil {
		diagnostics["database_error"] = err.Error()
	} else {
		diagnostics["database"] = dbDiagnostics
	}

	// Add environment info
	diagnostics["env"] = map[string]string{
		"log_path": utils.GetCurrentLogPath(),
		// Add more environment info as needed
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(diagnostics)
}

func performDatabaseDiagnostics(ctx context.Context, requestID string) (map[string]interface{}, error) {
	utils.Log(requestID, "Performing database diagnostics")

	mongoDb, err := db.NewFromEnv()

	dbInfo := make(map[string]interface{})
	if mongoDb.Database == nil {
		return nil, fmt.Errorf("MongoDB client not initialized")
	}

	// List collections
	collections, err := mongoDb.Database.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}
	dbInfo["collections"] = collections

	// Count documents in key collections
	collectionCounts := make(map[string]int64)
	for _, collName := range collections {
		count, err := mongoDb.Database.Collection(collName).CountDocuments(ctx, map[string]interface{}{})
		if err != nil {
			collectionCounts[collName] = -1 // Indicate error
		} else {
			collectionCounts[collName] = count
		}
	}
	dbInfo["collection_counts"] = collectionCounts

	utils.Log(requestID, "Database diagnostics completed")
	return dbInfo, nil
}

// Add this to your server routes:
// mux.HandleFunc("/api/diagnostics", api.DiagnosticHandler)
