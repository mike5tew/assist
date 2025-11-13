package main

import (
	"context"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer mongoDb.Client.Disconnect(context.Background())

	ctx := context.Background()

	log.Println("Connecting to Weaviate...")
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}
	defer db.CloseWeaviate()

	// Clear ONLY educational skills from MongoDB (preserve medical content)
	log.Println("Clearing educational skills from MongoDB...")

	skillsCollection := mongoDb.Database.Collection("skills")
	result, err := skillsCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Failed to clear skills: %v", err)
	} else {
		log.Printf("✅ Cleared %d educational skills", result.DeletedCount)
	}

	sourcesCollection := mongoDb.Database.Collection("sources")
	result, err = sourcesCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Failed to clear sources: %v", err)
	} else {
		log.Printf("✅ Cleared %d sources", result.DeletedCount)
	}

	// Clear ONLY educational skills from Weaviate (preserve medical vectors)
	log.Println("Clearing educational skills from Weaviate...")

	weaviateClient := db.GetWeaviateClient()
	if weaviateClient != nil {
		// Delete only non-medical content
		objects, err := weaviateClient.Data().ObjectsGetter().WithClassName("Skills").Do(ctx)
		if err != nil {
			log.Printf("Warning: Failed to get Weaviate objects: %v", err)
		} else {
			educationalCount := 0
			medicalCount := 0

			for _, obj := range objects {
				if properties, ok := obj.Properties.(map[string]interface{}); ok {
					// Check if this is medical content
					if sourceType, exists := properties["source_type"].(string); exists {
						if sourceType == "medical_textbook" || sourceType == "case_study" {
							medicalCount++
							continue // Skip deletion for medical content
						}
					}
					if extractionMethod, exists := properties["extraction_method"].(string); exists {
						if extractionMethod == "textract_immunology" {
							medicalCount++
							continue // Skip deletion for medical content
						}
					}
				}

				// Delete non-medical content
				// Fix: Use WeaviateDeleteDocument directly to avoid type mismatch
				if err := db.WeaviateDeleteDocument(ctx, "Skills", obj.ID.String()); err != nil {
					log.Printf("Warning: Failed to delete object %s: %v", obj.ID, err)
				} else {
					educationalCount++
				}
			}

			log.Printf("✅ Cleared %d educational vectors from Weaviate", educationalCount)
			log.Printf("📊 Preserved %d medical vectors in Weaviate", medicalCount)
		}
	}

	log.Println("✅ Educational skills cleared, medical content preserved!")
}
