package main

import (
	"context"
	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// check-data is a diagnostic utility that verifies the presence and integrity
// of data across different storage systems (MongoDB and Weaviate vector database).
// It checks for educational skills, academic sources, immunology chapters, case studies,
// and medical terms, displaying summary statistics about the data available in the system.
func main() {
	// Step 1: Initialize MongoDB client from environment variables
	mongoCl, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer mongoCl.Client.Disconnect(context.Background())

	// Step 2: Initialize Weaviate vector database connection
	// Retrieve connection parameters from environment variables or use defaults

	log.Println("Connecting to Weaviate...")
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}
	defer db.CloseWeaviate()

	// 3. Initialize LLAMA Client
	llamaClient := llm.NewLlamaClient()
	if llamaClient == nil {
		log.Fatal("Failed to create LlamaClient")
	}
	log.Println("LLAMA Client initialized successfully")

	// 4. Create Services
	// skillService := skills.NewSkillService(
	// 	mongDB,
	// 	db.GetEnvWithDefault("SKILLS_COLLECTION", db.DefaultSkillsColl),
	// )

	//queryHandler := query.NewQueryHandler(skillService, llamaClient)

	// Create a context for database operations
	ctx := context.Background()

	// Check Skills Collection - This section verifies educational skills data in MongoDB
	fmt.Println("=== EDUCATIONAL SKILLS DATA ===")
	skillsCollection := mongoCl.Database.Collection("skills")
	skillCount, err := skillsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error counting skills: %v", err)
	} else {
		fmt.Printf("📚 Skills in MongoDB: %d\n", skillCount)

		// If skills exist, display a sample skill name for verification
		if skillCount > 0 {
			var sampleSkill bson.M
			err := skillsCollection.FindOne(ctx, bson.M{}).Decode(&sampleSkill)
			if err == nil {
				fmt.Printf("   Sample: %v\n", sampleSkill["name"])
			}
		}
	}

	// Check Sources Collection - Verify academic source references
	fmt.Println("\n=== ACADEMIC SOURCES DATA ===")
	sourcesCollection := mongoCl.Database.Collection("sources")
	sourceCount, err := sourcesCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error counting sources: %v", err)
	} else {
		fmt.Printf("📖 Sources in MongoDB: %d\n", sourceCount)
	}

	// Check Immunology Collections - Verify medical/immunology content across different collections
	fmt.Println("\n=== IMMUNOLOGY CHAPTER DATA ===")

	// Check for immunology chapters in the database
	chaptersCollection := mongoCl.Database.Collection("immunology_chapters")
	chapterCount, err := chaptersCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error counting chapters: %v", err)
	} else {
		fmt.Printf("🔬 Immunology Chapters: %d\n", chapterCount)

		// Display a sample chapter title if any exist
		if chapterCount > 0 {
			var sampleChapter bson.M
			err := chaptersCollection.FindOne(ctx, bson.M{}).Decode(&sampleChapter)
			if err == nil {
				fmt.Printf("   Sample Chapter: %v\n", sampleChapter["chapter_title"])
			}
		}
	}

	// Check for medical case studies
	caseStudiesCollection := mongoCl.Database.Collection("case_studies")
	caseCount, err := caseStudiesCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error counting case studies: %v", err)
	} else {
		fmt.Printf("🏥 Case Studies: %d\n", caseCount)
	}

	// Check for medical terminology
	medicalTermsCollection := mongoCl.Database.Collection("medical_terms")
	termCount, err := medicalTermsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error counting medical terms: %v", err)
	} else {
		fmt.Printf("🧬 Medical Terms: %d\n", termCount)
	}

	// Check Weaviate Vector Database - Verify vector embeddings for semantic search
	fmt.Println("\n=== WEAVIATE VECTOR DATA ===")
	weaviateClient := db.GetWeaviateClient()
	if weaviateClient != nil {
		// Check Skills class in Weaviate
		objects, err := weaviateClient.Data().ObjectsGetter().WithClassName("Skills").Do(ctx)
		if err != nil {
			log.Printf("Error getting Weaviate objects: %v", err)
		} else {
			fmt.Printf("🔍 Vector embeddings in Weaviate: %d\n", len(objects))

			// Analyze vector types to check for medical content specifically
			medicalCount := 0
			for _, obj := range objects {
				if properties, ok := obj.Properties.(map[string]interface{}); ok {
					if sourceType, ok := properties["source_type"].(string); ok {
						if sourceType == "medical_textbook" || sourceType == "case_study" {
							medicalCount++
						}
					}
				}
			}
			fmt.Printf("   Medical content vectors: %d\n", medicalCount)
		}
	} else {
		fmt.Println("❌ Weaviate not connected")
	}

	// Provide overall system status and setup instructions if needed
	fmt.Println("\n=== SUMMARY ===")
	if skillCount > 0 {
		fmt.Println("✅ Educational skills database is loaded and ready")
	} else {
		fmt.Println("❌ Educational skills not loaded - run 'make setup'")
	}

	if chapterCount > 0 || caseCount > 0 || termCount > 0 {
		fmt.Println("✅ Some immunology content is loaded")
	} else {
		fmt.Println("❌ No immunology chapter content found")
		fmt.Println("   To load chapter content:")
		fmt.Println("   1. Place PDF in resources/ directory")
		fmt.Println("   2. Run: make upload-chapter")
		fmt.Println("   3. Or use frontend upload interface")
	}
}
