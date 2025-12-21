package api

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/store/db"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type CoachHandlerMVP struct {
	service interface{}
}

func NewCoachHandlerMVP(service interface{}) *CoachHandlerMVP {
	return &CoachHandlerMVP{service: service}
}

// CoachMVPDemoHandler - Dead simple MVP handler that queries MongoDB
func (h *CoachHandlerMVP) CoachMVPDemoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[MVP] Decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	log.Printf("[MVP] Query: %s", req.Message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Connect to MongoDB
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[MVP] MongoDB error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database connection failed"})
		return
	}
	defer mongoDb.Client.Disconnect(ctx)

	// Query for semantic links
	collection := mongoDb.Database.Collection("semantic_links")
	queryLower := strings.ToLower(req.Message)

	filter := bson.M{
		"$or": []bson.M{
			{"source_term": bson.M{"$regex": queryLower, "$options": "i"}},
			{"target_term": bson.M{"$regex": queryLower, "$options": "i"}},
			{"context": bson.M{"$regex": queryLower, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("[MVP] Query error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Query failed"})
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Printf("[MVP] Cursor error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read results"})
		return
	}

	// Convert to response format
	routes := make([]map[string]string, 0)
	for _, result := range results {
		route := map[string]string{
			"source_term":   getString(result, "source_term"),
			"target_term":   getString(result, "target_term"),
			"relation_type": getString(result, "relation_type"),
			"context":       getString(result, "context"),
		}
		routes = append(routes, route)
	}

	// If no results, return default
	if len(routes) == 0 {
		// Return a rich subgraph for the demo visualization (The "Map")
		routes = []map[string]string{
			{
				"source_term":   "X-Linked Agammaglobulinemia",
				"target_term":   "BTK Gene Mutation",
				"relation_type": "caused_by",
				"context":       "XLA is caused by mutations in the Bruton Tyrosine Kinase gene.",
			},
			{
				"source_term":   "X-Linked Agammaglobulinemia",
				"target_term":   "B-cell Deficiency",
				"relation_type": "causes",
				"context":       "The mutation prevents B-cells from maturing.",
			},
			{
				"source_term":   "B-cell Deficiency",
				"target_term":   "Low Antibodies",
				"relation_type": "results_in",
				"context":       "Without B-cells, the body cannot produce immunoglobulins.",
			},
			{
				"source_term":   "Low Antibodies",
				"target_term":   "Recurrent Infections",
				"relation_type": "leads_to",
				"context":       "Lack of antibodies makes the patient susceptible to bacterial infections.",
			},
			{
				"source_term":   "Recurrent Infections",
				"target_term":   "IVIG Therapy",
				"relation_type": "treated_with",
				"context":       "Patients require lifelong immunoglobulin replacement therapy.",
			},
		}
	}

	response := map[string]interface{}{
		"message":                "Found relevant knowledge for: " + req.Message,
		"response_id":            "resp_" + time.Now().Format("20060102150405"),
		"timestamp":              time.Now().Format(time.RFC3339),
		"knowledge_graph_routes": routes,
	}

	log.Printf("[MVP] Returning %d routes", len(routes))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getString(m bson.M, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
