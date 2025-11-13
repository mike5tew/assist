package api

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// DiagnosticsConfigHandler returns the current server configuration
func DiagnosticsConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	config := map[string]interface{}{
		"endpoints": map[string]interface{}{
			"immunology": map[string]string{
				"search":        "/api/immunology/search",
				"upload":        "/api/immunology/upload-chapter",
				"chapters":      "/api/immunology/chapters",
				"case_studies":  "/api/immunology/case-studies",
				"medical_terms": "/api/immunology/medical-terms",
			},
		},
		"domain_configs": domainConfigs,
		"timestamp":      time.Now(),
	}

	json.NewEncoder(w).Encode(config)
}

// DiagnosticsDocumentHandler retrieves a document by ID from a specific collection
func DiagnosticsDocumentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	collection := r.URL.Query().Get("collection")

	if id == "" || collection == "" {
		http.Error(w, "Missing required parameters: id and collection", http.StatusBadRequest)
		return
	}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(context.Background())

	var result bson.M
	err = mongoDb.Database.Collection(collection).FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)
	if err != nil {
		log.Printf("Document not found with exact ID, trying alternative formats: %v", err)
		err = mongoDb.Database.Collection(collection).FindOne(context.Background(), bson.M{"id": id}).Decode(&result)
		if err != nil {
			log.Printf("Document not found with ID %s in collection %s: %v", id, collection, err)
			http.Error(w, fmt.Sprintf("Document not found: %v", err), http.StatusNotFound)
			return
		}
	}

	json.NewEncoder(w).Encode(result)
}

// DiagnosticsCollectionsHandler returns a list of collections and document counts
func DiagnosticsCollectionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(context.Background())

	collections, err := mongoDb.Database.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Failed to list collections: %v", err)
		http.Error(w, "Failed to list collections", http.StatusInternalServerError)
		return
	}

	result := make(map[string]interface{})
	collectionStats := make([]map[string]interface{}, 0)

	for _, collectionName := range collections {
		collection := mongoDb.Database.Collection(collectionName)
		count, err := collection.CountDocuments(context.Background(), bson.M{})
		if err != nil {
			log.Printf("Failed to count documents in collection %s: %v", collectionName, err)
			continue
		}

		collectionStats = append(collectionStats, map[string]interface{}{
			"name":           collectionName,
			"document_count": count,
		})
	}

	result["collections"] = collectionStats
	result["timestamp"] = time.Now()

	json.NewEncoder(w).Encode(result)
}
