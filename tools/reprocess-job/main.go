package main

import (
	"context"
	"flag"
	"log"
	"time"

	"esp-organizer/internal/config"
	"esp-organizer/internal/store/db"
	"esp-organizer/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	jobID := flag.String("job", "", "Job ID to reprocess (e.g., immunology-chapter-1759699625)")
	flag.Parse()

	if *jobID == "" {
		log.Fatal("Usage: go run cmd/reprocess-job/main.go -job=<job_id>")
	}

	projectRoot, err := utils.FindProjectRoot()
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}
	config.LoadConfig(projectRoot, ".env")

	ctx := context.Background()

	// Connect to MongoDB
	mongoDB, err := db.NewMongoDBFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Client.Disconnect(ctx)

	log.Printf("🔄 Reprocessing job: %s", *jobID)

	// Find the job
	jobsCollection := mongoDB.Database.Collection("extraction_jobs")
	var job bson.M
	err = jobsCollection.FindOne(ctx, bson.M{"_id": *jobID}).Decode(&job)
	if err != nil {
		log.Fatalf("Job not found: %v", err)
	}

	log.Printf("✅ Found job with status: %v", job["status"])

	// Find the associated content by job creation time
	contentCollection := mongoDB.Database.Collection("immunology_content")
	var content bson.M

	err = contentCollection.FindOne(ctx, bson.M{
		"created_at": bson.M{"$gte": job["created_at"]},
	}).Decode(&content)

	if err != nil {
		log.Fatalf("Content not found for job: %v", err)
	}

	// Extract data from BSON map with nil checks
	contentID := content["_id"].(primitive.ObjectID)
	chapterContent := ""
	if cc, ok := content["chaptercontent"].(string); ok {
		chapterContent = cc
	}

	caseStudies := []interface{}{}
	medicalTerms := []interface{}{}

	// Check if extracteddata exists and is not nil
	if extractedData, ok := content["extracteddata"].(bson.M); ok && extractedData != nil {
		if cs, ok := extractedData["casestudies"].([]interface{}); ok {
			caseStudies = cs
		}
		if mt, ok := extractedData["medicalterms"].([]interface{}); ok {
			medicalTerms = mt
		}
	} else {
		log.Printf("⚠️ No extracted data found - job may not have completed successfully")
	}

	log.Printf("✅ Found content: %d case studies, %d medical terms",
		len(caseStudies), len(medicalTerms))
	log.Printf("📝 Chapter content length: %d characters", len(chapterContent))

	// Update job status to completed
	_, err = jobsCollection.UpdateOne(ctx,
		bson.M{"_id": *jobID},
		bson.M{
			"$set": bson.M{
				"status":     "completed",
				"updated_at": time.Now(),
			},
		},
	)

	if err != nil {
		log.Fatalf("Failed to update job status: %v", err)
	}

	log.Printf("✅ Job reprocessing complete!")
	log.Printf("📊 Summary:")
	log.Printf("   • Job ID: %s", *jobID)
	log.Printf("   • Content ID: %s", contentID.Hex())
	log.Printf("   • Status: completed")
	log.Printf("   • Case Studies: %d", len(caseStudies))
	log.Printf("   • Medical Terms: %d", len(medicalTerms))
}
