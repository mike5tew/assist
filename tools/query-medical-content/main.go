package main

import (
	"context"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {

	ctx := context.Background()
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer mongoDb.Client.Disconnect(context.Background())

	// Step 2: Initialize Weaviate vector database connection
	// Retrieve connection parameters from environment variables or use defaults
	//	weaviateHost := db.GetEnvWithDefault("WEAVIATE_HOST", db.DefaultWeaviateHost)
	//	weaviateScheme := db.GetEnvWithDefault("WEAVIATE_SCHEME", db.DefaultWeaviateScheme)

	// Query immunology chapters
	fmt.Println("=== Processed Immunology Chapters ===")
	chaptersCollection := mongoDb.Database.Collection("immunology_chapters")
	cursor, err := chaptersCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to query chapters: %v", err)
	} else {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var chapter bson.M
			if err := cursor.Decode(&chapter); err != nil {
				log.Printf("Failed to decode chapter: %v", err)
				continue
			}
			fmt.Printf("Chapter: %v\n", chapter["chapter_title"])
			fmt.Printf("Book: %v\n", chapter["book_title"])
			fmt.Printf("Case Studies: %v, Medical Terms: %v\n",
				chapter["case_studies_count"], chapter["medical_terms_count"])
			fmt.Printf("Job ID: %v\n", chapter["job_id"])
			fmt.Printf("Extracted: %v\n\n", chapter["extracted_at"])
		}
	}

	// Query case studies
	fmt.Println("=== Extracted Case Studies ===")
	caseStudiesCollection := mongoDb.Database.Collection("case_studies")
	cursor, err = caseStudiesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to query case studies: %v", err)
	} else {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var caseStudy bson.M
			if err := cursor.Decode(&caseStudy); err != nil {
				log.Printf("Failed to decode case study: %v", err)
				continue
			}
			fmt.Printf("Case Number: %v\n", caseStudy["case_number"])
			fmt.Printf("Content Preview: %.200s...\n", caseStudy["content"])
			if patientInfo, ok := caseStudy["patient_info"].(bson.M); ok {
				fmt.Printf("Patient: %v, %v\n", patientInfo["age"], patientInfo["gender"])
			}
			fmt.Println()
		}
	}

	// Query medical terms
	fmt.Println("=== Extracted Medical Terms ===")
	medicalTermsCollection := mongoDb.Database.Collection("medical_terms")
	cursor, err = medicalTermsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to query medical terms: %v", err)
	} else {
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var term bson.M
			if err := cursor.Decode(&term); err != nil {
				log.Printf("Failed to decode medical term: %v", err)
				continue
			}
			fmt.Printf("Term: %v (Category: %v, Frequency: %v)\n",
				term["term"], term["category"], term["frequency"])
			if context, ok := term["context"].(bson.A); ok && len(context) > 0 {
				fmt.Printf("Context: %v\n", context[0])
			}
			fmt.Println()
		}
	}
}
