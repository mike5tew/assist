package api

import (
	"encoding/json"
	"net/http"
	"time"

	"esp-organizer/internal/InfoFlow/InfoIn"
	wfilters "esp-organizer/internal/InfoFlow/InfoIn/filters"
	"esp-organizer/internal/InfoFlow/InfoStore/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DomainConfig holds the collections for a given domain.

// HandleKnowledgeQuery runs a Mongo text‐search + Weaviate semantic search.
// HandleKnowledgeQuery runs a Mongo text-search + Weaviate semantic search with MongoDB enrichment.
func HandleKnowledgeQuery(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain string `json:"domain"`
		Query  string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	// 1) Connect to MongoDB
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "mongo connect: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(ctx)

	// 2) Load domain→collections mapping
	var cfg DomainConfig
	if err := mongoDb.Database.Collection("domains").
		FindOne(ctx, bson.M{"domain": req.Domain}).
		Decode(&cfg); err != nil {
		http.Error(w, "unsupported domain or missing config", http.StatusBadRequest)
		return
	}

	// 3) Mongo text searches
	mongoResults := make(map[string][]bson.M, len(cfg.MongoCollections))
	for _, coll := range cfg.MongoCollections {
		cur, err := mongoDb.Database.Collection(coll).
			Find(ctx, bson.M{"$text": bson.M{"$search": req.Query}},
				options.Find().
					SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}}).
					SetSort(bson.M{"score": bson.M{"$meta": "textScore"}}),
			)
		if err != nil {
			http.Error(w, "mongo find: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var docs []bson.M
		_ = cur.All(ctx, &docs)
		mongoResults[coll] = docs
	}

	if db.GetWeaviateClient() == nil {
		if err := db.InitializeWeaviateFromEnv(); err != nil {
			http.Error(w, "weaviate init: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// turn the user query into a vector via your existing service
	vectorizer := InfoIn.NewVectorizationService()
	vec, err := vectorizer.GenerateEmbedding(ctx, req.Query)
	if err != nil {
		http.Error(w, "embed error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// build a simple domain-filter
	where := wfilters.BuildWhereFilter(map[string]interface{}{
		"path":  []string{"domain"},
		"op":    "Equal",
		"value": req.Domain,
	})

	// execute the nearVector search
	semMatches, err := db.QueryWeaviateData(ctx, "Skills", vec, req.Query, 10, where)
	if err != nil {
		http.Error(w, "weaviate search: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// NEW CODE: Enrich semantic matches with MongoDB data
	enrichedMatches := make([]map[string]interface{}, 0, len(semMatches))
	for _, match := range semMatches {
		properties, ok := match["properties"].(map[string]interface{})
		if !ok {
			continue
		}

		// Create base enriched result
		enrichedMatch := map[string]interface{}{
			"weaviate_data": properties,
			"certainty":     0.0,
		}

		// Extract certainty score if available
		if additional, ok := properties["_additional"].(map[string]interface{}); ok {
			if certainty, ok := additional["certainty"].(float64); ok {
				enrichedMatch["certainty"] = certainty
			}
		}

		// Check for MongoDB ID reference
		if mongoIDStr, ok := properties["mongo_id"].(string); ok {
			// Determine which collection to query based on content_type
			var collectionName string
			if contentType, ok := properties["content_type"].(string); ok {
				// Map content_type to collection name - adjust this mapping as needed
				switch contentType {
				case "skill":
					collectionName = "skills"
				case "case_study":
					collectionName = "case_studies"
				case "medical_term":
					collectionName = "medical_terms"
				// Add other mappings as needed
				default:
					collectionName = "skills" // Default fallback
				}
			} else {
				collectionName = "skills" // Default fallback
			}

			// Convert string ID to ObjectID
			mongoID, err := primitive.ObjectIDFromHex(mongoIDStr)
			if err == nil {
				// Fetch complete document from the appropriate MongoDB collection
				var fullDoc bson.M
				err = mongoDb.Database.Collection(collectionName).
					FindOne(ctx, bson.M{"_id": mongoID}).
					Decode(&fullDoc)

				if err == nil {
					// Add MongoDB data to enriched result
					enrichedMatch["mongodb_data"] = fullDoc

					// Add key fields directly to top level for convenience
					if name, ok := fullDoc["name"].(string); ok {
						enrichedMatch["name"] = name
					}
					if description, ok := fullDoc["description"].(string); ok {
						enrichedMatch["description"] = description
					}
				}
			}
		}

		enrichedMatches = append(enrichedMatches, enrichedMatch)
	}

	// 5) Write combined response with enriched semantic matches
	resp := map[string]interface{}{
		"domain":                    req.Domain,
		"query":                     req.Query,
		"mongo_results":             mongoResults,
		"enriched_semantic_matches": enrichedMatches,
		"timestamp":                 time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "json encode: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
