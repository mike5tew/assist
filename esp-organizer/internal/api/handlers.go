package api

import (
	"context"
	"esp-organizer/internal/InfoFlow/InfoOut/llm"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"esp-organizer/internal/InfoFlow/InfoStore/skills"
	"fmt"
	"log"
	"math"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HandleQuery is the HTTP handler for processing queries.
// NOTE: This function appears to be for debugging or a legacy CLI tool, not a standard HTTP handler.
// It has been updated to use the modern SkillService for searching.
func HandleQuery(ctx context.Context, query string, limit int, whereFilter *filters.WhereBuilder) {
	// Log the incoming query
	log.Printf("Received query: %s", query)

	// Initialize services
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}
	client := llm.NewLlamaClient()
	skillService := skills.NewSkillService(mongoDb.Database.Collection("skills"), db.GetWeaviateClient(), client)

	// The modern search function takes a models.Filter, not a weaviate.WhereBuilder.
	// For this example, we'll perform the search without the filter.
	// A proper implementation would convert the WhereBuilder or pass the models.Filter.
	weaviateResults, err := skillService.SearchSkillsByVector(ctx, query, nil)
	if err != nil {
		log.Printf("Weaviate query failed for '%s': %v", query, err)
		return
	}

	// Prepare to collect enriched results
	enrichedResults := make([]map[string]interface{}, 0, len(weaviateResults))

	fmt.Printf("Top %d Weaviate matches (searching SemanticLinks):\n", limit)
	for i, result := range weaviateResults {
		if i >= limit {
			break
		}

		if properties, ok := result["properties"].(map[string]interface{}); ok {
			// Extract basic information from the semantic link
			sourceTerm := "Unknown"
			targetTerm := "Unknown"
			certainty := 0.0

			if st, ok := properties["source_term"].(string); ok {
				sourceTerm = st
			}
			if tt, ok := properties["target_term"].(string); ok {
				targetTerm = tt
			}

			if additional, ok := properties["_additional"].(map[string]interface{}); ok {
				if c, ok := additional["certainty"].(float64); ok {
					certainty = c
				}
			}

			// Create enriched result with basic info
			enrichedResult := map[string]interface{}{
				"source_term": sourceTerm,
				"target_term": targetTerm,
				"certainty":   certainty,
				"weaviate":    properties,
			}

			fmt.Printf("  %d. Link: %s -> %s (Certainty: %.4f)\n", i+1, sourceTerm, targetTerm, certainty)

			// Enrich with MongoDB data if mongo_id is available
			if mongoIDStr, ok := properties["source_mongo_id"].(string); ok && mongoIDStr != "" {
				mongoID, err := primitive.ObjectIDFromHex(mongoIDStr)
				if err == nil {
					skill, err := db.GetSkillByID(ctx, mongoID)
					if err == nil {
						enrichedResult["mongodb_source_skill"] = skill
						fmt.Printf("     Source Skill: %s (%v)\n", skill.Name, skill)
					}
				}
			}

			// Add to collection of enriched results
			enrichedResults = append(enrichedResults, enrichedResult)
		}
	}
	fmt.Println()

	// Return enrichedResults for API responses, etc.
	// You could modify this function signature to return the enriched data
	// Example: return enrichedResults, nil
}

func calculateMagnitude(vec []float32) float32 {
	var sum float32
	for _, v := range vec {
		sum += v * v
	}
	return float32(math.Sqrt(float64(sum)))
}

func sumVector(vec []float32) float32 {
	var sum float32
	for _, v := range vec {
		sum += v
	}
	return sum
}

func cosineSimilarity(a, b []float32) float32 {
	var dotProduct, normA, normB float32
	for i := 0; i < len(a) && i < len(b); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
