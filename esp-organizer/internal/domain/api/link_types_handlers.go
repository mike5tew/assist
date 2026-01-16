package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// getMongoDatabase returns a connected MongoDB database
func getMongoDatabase() (*db.MongoDB, error) {
	return db.NewFromEnv()
}

// LinkTypesListHandler returns all saved link types
// GET /api/link-types
func LinkTypesListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	mongoDB, err := getMongoDatabase()
	if err != nil {
		log.Printf("Error getting mongo client: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	collection := mongoDB.Database.Collection("link_types")

	// Sort by usage count (most used first), then alphabetically
	opts := options.Find().SetSort(bson.D{{Key: "usage_count", Value: -1}, {Key: "link_a_to_b", Value: 1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Printf("Error querying link types: %v", err)
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var linkTypes []models.LinkType
	if err := cursor.All(ctx, &linkTypes); err != nil {
		log.Printf("Error decoding link types: %v", err)
		http.Error(w, "decode failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"link_types": linkTypes,
		"count":      len(linkTypes),
	})
}

// LinkTypeCreateRequest represents the request body for creating a link type
type LinkTypeCreateRequest struct {
	LinkAToB string `json:"link_a_to_b"`
	LinkBToA string `json:"link_b_to_a"`
	Category string `json:"category,omitempty"`
}

// LinkTypesCreateHandler creates a new link type
// POST /api/link-types
func LinkTypesCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LinkTypeCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.LinkAToB == "" {
		http.Error(w, "link_a_to_b is required", http.StatusBadRequest)
		return
	}

	// If reverse not provided, try to infer it
	if req.LinkBToA == "" {
		req.LinkBToA = inferReverseLink(req.LinkAToB)
	}

	mongoDB, err := getMongoDatabase()
	if err != nil {
		log.Printf("Error getting mongo client: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	collection := mongoDB.Database.Collection("link_types")

	// Check if already exists
	var existing models.LinkType
	err = collection.FindOne(ctx, bson.M{"link_a_to_b": req.LinkAToB}).Decode(&existing)
	if err == nil {
		// Already exists, return it
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"link_type": existing,
			"created":   false,
			"message":   "link type already exists",
		})
		return
	}

	linkType := models.LinkType{
		ID:         primitive.NewObjectID(),
		LinkAToB:   req.LinkAToB,
		LinkBToA:   req.LinkBToA,
		Category:   req.Category,
		UsageCount: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = collection.InsertOne(ctx, linkType)
	if err != nil {
		log.Printf("Error inserting link type: %v", err)
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Created new link type: %s ↔ %s", linkType.LinkAToB, linkType.LinkBToA)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"link_type": linkType,
		"created":   true,
	})
}

// LinkTypesSeedHandler seeds the database with common link types
// POST /api/link-types/seed
func LinkTypesSeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	mongoDB, err := getMongoDatabase()
	if err != nil {
		log.Printf("Error getting mongo database: %v", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	collection := mongoDB.Database.Collection("link_types")

	inserted := 0
	skipped := 0

	for _, lt := range models.CommonLinkTypes {
		// Check if exists
		var existing models.LinkType
		err := collection.FindOne(ctx, bson.M{"link_a_to_b": lt.LinkAToB}).Decode(&existing)
		if err == nil {
			skipped++
			continue
		}

		linkType := models.LinkType{
			ID:         primitive.NewObjectID(),
			LinkAToB:   lt.LinkAToB,
			LinkBToA:   lt.LinkBToA,
			Category:   lt.Category,
			UsageCount: 0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		_, err = collection.InsertOne(ctx, linkType)
		if err != nil {
			log.Printf("Error inserting seed link type '%s': %v", lt.LinkAToB, err)
			continue
		}
		inserted++
	}

	log.Printf("✅ Seeded link types: %d inserted, %d skipped (already exist)", inserted, skipped)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"inserted": inserted,
		"skipped":  skipped,
		"total":    len(models.CommonLinkTypes),
	})
}

// inferReverseLink attempts to generate a reverse link label
func inferReverseLink(forward string) string {
	// Check common patterns
	reverseMap := map[string]string{
		"causes":             "is caused by",
		"produces":           "is produced by",
		"inhibits":           "is inhibited by",
		"activates":          "is activated by",
		"requires":           "is required by",
		"contains":           "is contained in",
		"includes":           "is included in",
		"precedes":           "follows",
		"follows":            "precedes",
		"leads to":           "results from",
		"results from":       "leads to",
		"is a type of":       "includes",
		"is part of":         "contains",
		"treats":             "is treated by",
		"prevents":           "is prevented by",
		"is associated with": "is associated with",
		"is similar to":      "is similar to",
		"is different from":  "is different from",
	}

	if reverse, ok := reverseMap[forward]; ok {
		return reverse
	}

	// Default: prepend "is ... by"
	return "is " + forward + " by"
}
