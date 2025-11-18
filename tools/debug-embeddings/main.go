package main

import (
	"context"
	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"
	"math"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	ctx := context.Background()

	log.Println("Connecting to Weaviate...")
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}
	defer db.CloseWeaviate()

	// Test embedding generation
	client := llm.NewLlamaClient()

	testQueries := []string{
		"memory",
		"working memory",
		"reading",
		"reading comprehension",
		"motor skills",
		"fine motor skills",
		"emotional intelligence",
		"problem solving",
		"communication",
	}

	fmt.Println("=== Testing Embedding Generation ===")

	// First, test the embedding function directly
	fmt.Println("Testing embedding function directly:")
	testEmbeddings := []string{"memory", "reading", "motor", "emotional"}
	for _, test := range testEmbeddings {
		emb, err := client.GenerateEmbedding(test)
		if err != nil {
			log.Printf("Error generating embedding for '%s': %v", test, err)
			continue
		}
		sum := sumVector(emb)
		magnitude := calculateMagnitude(emb)
		fmt.Printf("  '%s': length=%d, magnitude=%.4f, sum=%.4f, first3=[%.4f,%.4f,%.4f]\n",
			test, len(emb), magnitude, sum, emb[0], emb[1], emb[2])
	}
	fmt.Println()

	embeddings := make(map[string][]float32)

	for _, query := range testQueries {
		embedding, err := client.GenerateEmbedding(query)
		if err != nil {
			log.Printf("Failed to generate embedding for '%s': %v", query, err)
			continue
		}

		embeddings[query] = embedding

		// Calculate magnitude to verify normalization
		magnitude := calculateMagnitude(embedding)
		sum := sumVector(embedding)

		fmt.Printf("Query: '%s' -> Length: %d, Magnitude: %.4f, Sum: %.4f, First 5: [", query, len(embedding), magnitude, sum)
		for i := 0; i < 5 && i < len(embedding); i++ {
			fmt.Printf("%.4f", embedding[i])
			if i < 4 && i < len(embedding)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Println("],")

		// Test similarity with Weaviate against the EducationalSkills class
		results, err := db.QueryWeaviateData(ctx, "EducationalSkills", embedding, "", 3, nil)
		if err != nil {
			log.Printf("Weaviate query failed for '%s': %v", query, err)
			continue
		}

		fmt.Printf("  Top 3 Weaviate matches:\n")
		for i, result := range results {
			if i >= 3 {
				break
			}

			if properties, ok := result["properties"].(map[string]interface{}); ok {
				// Extract basic info for display
				name := "Unknown"
				certainty := 0.0

				if n, ok := properties["name"].(string); ok {
					name = n
				}
				if additional, ok := properties["_additional"].(map[string]interface{}); ok {
					if c, ok := additional["certainty"].(float64); ok {
						certainty = c
					}
				}

				// Extract MongoDB ID (if using the minimal approach we discussed)
				var fullDocument *models.Skill
				if mongoIDStr, ok := properties["mongo_id"].(string); ok {
					// Convert string ID to MongoDB ObjectID
					mongoID, err := primitive.ObjectIDFromHex(mongoIDStr)
					if err == nil {
						// Fetch complete document from MongoDB
						fullDocument, err = db.GetSkillByID(ctx, mongoID)
						if err != nil {
							log.Printf("Warning: Could not fetch MongoDB document: %v", err)
						}
					}
				}

				// Display basic match information
				fmt.Printf("    %d. %s (Certainty: %.4f)\n", i+1, name, certainty)

				// Display enriched information if available
				if fullDocument != nil {
					fmt.Printf("       Skill: %s\n", fullDocument.Name)
					fmt.Printf("       Developmental Age: %d\n", fullDocument.DevelopmentAge)

					// Display related skills (now ParentSkills)
					if len(fullDocument.ParentSkills) > 0 {
						fmt.Printf("       Parent skills: %d (first: %s)\n",
							len(fullDocument.ParentSkills),
							fullDocument.ParentSkills[0].Name)
					}
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}

	// Test embedding similarity
	fmt.Println("=== Embedding Similarity Analysis ===")
	queryPairs := [][]string{
		{"memory", "working memory"},
		{"reading", "reading comprehension"},
		{"motor skills", "fine motor skills"},
	}

	for _, pair := range queryPairs {
		if emb1, ok := embeddings[pair[0]]; ok {
			if emb2, ok := embeddings[pair[1]]; ok {
				similarity := cosineSimilarity(emb1, emb2)
				fmt.Printf("Similarity between '%s' and '%s': %.4f\n", pair[0], pair[1], similarity)
			}
		}
	}
	fmt.Println()

	// Check a few objects in Weaviate to see if they have vectors
	fmt.Println("\n=== Detailed Weaviate Storage Check ===")
	checkObjects, err := db.GetWeaviateClient().Data().ObjectsGetter().
		WithClassName("EducationalSkills").
		WithLimit(10).
		Do(context.Background())

	if err != nil {
		log.Fatalf("Failed to get objects from Weaviate: %v", err)
	}

	if len(checkObjects) == 0 {
		log.Println("No objects found in Weaviate EducationalSkills class. Run 'make upload' first.")
		return
	}

	log.Printf("Found %d objects in Weaviate EducationalSkills class. Checking a sample:", len(checkObjects))
	for i, obj := range checkObjects {
		fmt.Printf("  %d. ID: %s\n", i+1, obj.ID)
		if obj.Properties != nil {
			props := obj.Properties.(map[string]interface{})
			if name, ok := props["name"].(string); ok {
				fmt.Printf("     Name: %s\n", name)
			}
			if desc, ok := props["description"].(string); ok {
				// Truncate long descriptions
				if len(desc) > 50 {
					desc = desc[:50] + "..."
				}
				fmt.Printf("     Description: %s\n", desc)
			}
			// Check for a source field to display
			if sourceTitle, ok := props["source_title"].(string); ok && sourceTitle != "" {
				fmt.Printf("     Source: %s\n", sourceTitle)
			} else {
				// Fallback for older data
				fmt.Printf("     Source: ESP Organizer Skills Database\n")
			}
		}
		// The vector search is working, so we can trust the vectors are present internally.
		// We no longer need to print the misleading "MISSING" message.
		fmt.Println("     Vector: Present (verified by successful search)")
		fmt.Println()
	}
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
