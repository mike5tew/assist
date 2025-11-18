package main

import (
	"context"
	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	log.Println("🔧 Medical Content Vector Repair Tool")
	log.Println("===================================")

	// 1. Initialize connections
	ctx := context.Background()

	log.Println("Connecting to Weaviate...")
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}
	defer db.CloseWeaviate()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not set in .env file")
	}
	mongoDbj, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDbj.Client.Disconnect(ctx)
	skillsCollection := mongoDbj.Client.Database(os.Getenv("MONGO_DB_NAME)")).Collection(os.Getenv("SKILLS_COLLECTION"))
	log.Println("✅ Successfully connected to MongoDB.")

	llamaClient := llm.NewLlamaClient()
	log.Println("✅ AWS Bedrock configured")

	// 2. Delete ALL existing medical content from Weaviate to ensure a clean slate.
	log.Println("Deleting all existing medical content from Weaviate...")
	client := db.GetWeaviateClient()

	medicalTypes := []string{"medical_term", "medical_chapter", "medical_case_study"}
	whereFilter := filters.Where().
		WithPath([]string{"skill_type"}).
		WithOperator(filters.ContainsAny).
		WithValueText(medicalTypes...)

	_, err = client.Batch().ObjectsBatchDeleter().
		WithClassName("Skills").
		WithWhere(whereFilter).
		Do(ctx)

	if err != nil {
		log.Fatalf("Failed to batch delete existing medical content: %v", err)
	}
	log.Println("✅ All previous medical content deleted from Weaviate.")

	// 3. Fetch all medical content from MongoDB to re-ingest
	log.Println("Fetching fresh medical content from MongoDB...")
	mongoFilter := bson.M{"skill_type": bson.M{"$in": medicalTypes}}
	cursor, err := skillsCollection.Find(ctx, mongoFilter)
	if err != nil {
		log.Fatalf("Failed to query MongoDB for medical skills: %v", err)
	}
	defer cursor.Close(ctx)

	var mongoSkills []models.Skill
	if err = cursor.All(ctx, &mongoSkills); err != nil {
		log.Fatalf("Failed to decode medical skills from MongoDB: %v", err)
	}

	if len(mongoSkills) == 0 {
		log.Println("No medical skills found in MongoDB to ingest. Exiting.")
		return
	}

	if len(mongoSkills) == 0 {
		log.Println("\n🔍 Diagnostic Information:")
		log.Println("- Checking for all available skill types in MongoDB...")

		// Query to get all distinct skill types
		skillTypes, err := skillsCollection.Distinct(ctx, "skill_type", bson.M{})
		if err != nil {
			log.Printf("⚠️ Error querying skill types: %v", err)
		} else {
			log.Printf("- Found skill types: %v", skillTypes)
			log.Println("- To add medical content, you need to create skills with skill_type: 'medical_term', 'medical_chapter', or 'medical_case_study'")
		}

		// Count total skills
		count, err := skillsCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			log.Printf("⚠️ Error counting skills: %v", err)
		} else {
			log.Printf("- Total skills in database: %d", count)
		}

		log.Println("\n📋 Next Steps:")
		log.Println("1. Use the seed-medical-data command to add sample medical content")
		log.Println("2. Check your MongoDB collection schema")
		log.Println("3. Verify the ESP Organizer is configured to use medical skill types correctly")
	}

	log.Printf("Found %d medical skills in MongoDB to re-ingest.", len(mongoSkills))

	// 4. Re-ingest the content into Weaviate with fresh embeddings
	var fixed, failed int
	for _, skill := range mongoSkills {
		content := fmt.Sprintf("Name: %s\nDescription: %s", skill.Name, skill.Description)
		log.Printf("Embedding and ingesting: %s", skill.Name)
		embedding, err := llamaClient.GenerateEmbedding(content)
		if err != nil {
			log.Printf("⚠️ Failed to generate embedding for %s: %v", skill.Name, err)
			failed++
			continue
		}

		properties := map[string]interface{}{
			"name":              skill.Name,
			"description":       skill.Description,
			"source_title":      skill.SourceTitle,
			"source_type":       skill.SourceType,
			"extraction_method": skill.ExtractionMethod,

			"mongo_id": skill.ID.Hex(),
		}
		if skill.ChapterTitle != "" {
			properties["chapter_title"] = skill.ChapterTitle
		}

		_, err = client.Data().Creator().
			WithClassName("Skills").
			WithProperties(properties).
			WithVector(embedding).
			Do(ctx)

		if err != nil {
			log.Printf("⚠️ Failed to re-ingest skill %s: %v", skill.Name, err)
			failed++
		} else {
			fixed++
		}
		time.Sleep(100 * time.Millisecond) // Small delay to avoid overwhelming any service
	}

	log.Printf("\n=== Vector Repair Complete ===")
	log.Printf("Total items to process: %d", len(mongoSkills))
	log.Printf("✅ Successfully ingested: %d", fixed)
	log.Printf("❌ Failed to ingest: %d", failed)

	if failed > 0 {
		log.Printf("⚠️ Some items could not be re-ingested. Check logs for details.")
	} else {
		log.Printf("🎉 All medical content vectors successfully repaired!")
	}

	log.Println("\nRun 'make debug-vector-search' to verify the fix worked.")
}
