package main

import (
	"context"
	"esp-organizer/internal/store/db"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer mongoDb.Client.Disconnect(context.Background())

	// Step 2: Initialize Weaviate vector database connection
	// Retrieve connection parameters from environment variables or use defaults
	//weaviateHost := db.GetEnvWithDefault("WEAVIATE_HOST", db.DefaultWeaviateHost)
	//weaviateScheme := db.GetEnvWithDefault("WEAVIATE_SCHEME", db.DefaultWeaviateScheme)

	ctx := context.Background()

	log.Println("🔍 Searching for mock/test medical data...")

	// Clear chapters that look like mock data
	chaptersCollection := mongoDb.Database.Collection("immunology_chapters")

	// Find chapters with suspicious patterns
	mockChapterFilter := bson.M{
		"$or": []bson.M{
			{"case_studies_count": 1, "medical_terms_count": 1}, // Exact 1-1 pattern is suspicious
			{"job_id": bson.M{"$regex": "^[a-f0-9]{32}$"}},      // Old hex job IDs (not textract- prefixed)
		},
	}

	chapters, err := chaptersCollection.Find(ctx, mockChapterFilter)
	if err != nil {
		log.Printf("Error finding chapters: %v", err)
		return
	}
	defer chapters.Close(ctx)

	var chapterIDs []interface{}
	chapterCount := 0
	for chapters.Next(ctx) {
		var chapter bson.M
		if err := chapters.Decode(&chapter); err != nil {
			continue
		}
		chapterIDs = append(chapterIDs, chapter["_id"])
		chapterCount++
		log.Printf("Found suspicious chapter: %s (job: %v, cases: %v, terms: %v)",
			chapter["chapter_title"], chapter["job_id"], chapter["case_studies_count"], chapter["medical_terms_count"])
	}

	if chapterCount > 0 {
		// Delete chapters
		result, err := chaptersCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": chapterIDs}})
		if err != nil {
			log.Printf("Error deleting chapters: %v", err)
		} else {
			log.Printf("✅ Deleted %d suspicious chapters", result.DeletedCount)
		}

		// Clear case studies that look like mock data
		caseStudiesCollection := mongoDb.Database.Collection("case_studies")
		mockCaseFilter := bson.M{
			"$or": []bson.M{
				{"content": "Sample case study content"},
				{"clinical_findings": "Finding 1"},
				{"patient_info.age": "25", "patient_info.gender": "Female"}, // Exact match pattern
			},
		}

		result, err = caseStudiesCollection.DeleteMany(ctx, mockCaseFilter)
		if err != nil {
			log.Printf("Error deleting case studies: %v", err)
		} else {
			log.Printf("✅ Deleted %d mock case studies", result.DeletedCount)
		}

		// Clear medical terms that look like mock data
		medicalTermsCollection := mongoDb.Database.Collection("medical_terms")
		mockTermFilter := bson.M{
			"$or": []bson.M{
				{"term": "Immunoglobulin", "category": "Protein", "frequency": 5},
				{"context": []interface{}{"antibody", "immune response"}}, // Exact array match
			},
		}

		result, err = medicalTermsCollection.DeleteMany(ctx, mockTermFilter)
		if err != nil {
			log.Printf("Error deleting medical terms: %v", err)
		} else {
			log.Printf("✅ Deleted %d mock medical terms", result.DeletedCount)
		}
	} else {
		log.Println("No suspicious mock data found")
	}

	log.Println("✅ Mock data cleanup complete!")
	log.Println("💡 Now try uploading a fresh PDF to test real AWS Textract processing")
}
