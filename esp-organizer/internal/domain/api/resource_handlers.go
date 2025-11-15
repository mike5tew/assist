package api

import (
	"encoding/json"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetSourcesHandler returns a list of all sources (e.g., books).
func GetSourcesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(r.Context())

	q := r.URL.Query()
	srcType := strings.TrimSpace(q.Get("type"))
	limitStr := strings.TrimSpace(q.Get("limit"))
	raw := strings.EqualFold(q.Get("raw"), "true")

	filter := bson.M{}
	if srcType != "" {
		filter["type"] = srcType
	}

	findOpts := options.Find()
	var limit int64
	if limitStr != "" {
		if v, perr := strconv.ParseInt(limitStr, 10, 64); perr == nil && v > 0 {
			limit = v
			findOpts.SetLimit(limit)
		}
	}

	sourcesCollection := mongoDb.Database.Collection("sources")
	cursor, err := sourcesCollection.Find(r.Context(), filter, findOpts)
	if err != nil {
		http.Error(w, "Failed to query sources", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var rawDocs []bson.M
	if err := cursor.All(r.Context(), &rawDocs); err != nil {
		http.Error(w, "Failed to decode sources", http.StatusInternalServerError)
		return
	}

	sources := make([]map[string]interface{}, 0, len(rawDocs))
	for _, d := range rawDocs {
		out := make(map[string]interface{}, len(d)+1)
		for k, v := range d {
			if k == "_id" {
				switch idv := v.(type) {
				case primitive.ObjectID:
					out["id"] = idv.Hex()
				default:
					out["id"] = fmt.Sprintf("%v", v)
				}
				continue
			}
			out[k] = v
		}
		sources = append(sources, out)
	}

	if sources == nil {
		sources = []map[string]interface{}{}
	}

	if raw {
		json.NewEncoder(w).Encode(sources)
		return
	}

	response := map[string]interface{}{
		"sources": sources,
		"count":   len(sources),
		"filter": map[string]interface{}{
			"type":  srcType,
			"limit": limit,
		},
		"timestamp": time.Now(),
	}
	json.NewEncoder(w).Encode(response)
}

// CreateSourceHandler adds a new book or other source to the database
func CreateSourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var source models.Source

	if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if source.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if source.Type == "" {
		source.Type = "book" // Default type
	}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(r.Context())

	now := time.Now()
	doc := bson.M{
		"type":       source.Type,
		"isbn":       source.ISBN,
		"title":      source.Title,
		"authors":    source.Authors,
		"publisher":  source.Publisher,
		"year":       source.Year,
		"processed":  false,
		"created_at": now,
		"updated_at": now,
	}

	sourcesCollection := mongoDb.Database.Collection("sources")
	result, err := sourcesCollection.InsertOne(r.Context(), doc)
	if err != nil {
		http.Error(w, "Failed to insert source: "+err.Error(), http.StatusInternalServerError)
		return
	}

	doc["id"] = result.InsertedID
	delete(doc, "_id")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
}

// ListSourcesHandler returns a list of sources (books) for selection
func ListSourcesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(r.Context())

	filter := bson.M{}
	if t := r.URL.Query().Get("type"); t != "" {
		filter["type"] = t
	}
	sourcesCollection := mongoDb.Database.Collection("sources")
	cursor, err := sourcesCollection.Find(r.Context(), filter)
	if err != nil {
		http.Error(w, "Failed to query sources", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var sources []map[string]interface{}
	for cursor.Next(r.Context()) {
		var source bson.M
		if err := cursor.Decode(&source); err != nil {
			continue
		}
		sources = append(sources, source)
	}
	response := map[string]interface{}{
		"sources":     sources,
		"total_count": len(sources),
	}
	json.NewEncoder(w).Encode(response)
}

// GetProcessingBatchesHandler returns processing batches for rollback management
func GetProcessingBatchesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(r.Context())

	batchesCollection := mongoDb.Database.Collection("processing_batches")
	cursor, err := batchesCollection.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to query batches", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var batches []map[string]interface{}
	if err = cursor.All(r.Context(), &batches); err != nil {
		http.Error(w, "Failed to decode batches", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(batches)
}
