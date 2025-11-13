package main

import (
	"context"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"esp-organizer/internal/models"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/rollback-batch/main.go <batch_id>")
	}
	ctx := context.Background()
	batchID := os.Args[1]

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer mongoDb.Client.Disconnect(context.Background())

	// Step 2: Initialize Weaviate vector database connection
	// Retrieve connection parameters from environment variables or use defaults
	//weaviateHost := db.GetEnvWithDefault("WEAVIATE_HOST", db.DefaultWeaviateHost)
	//weaviateScheme := db.GetEnvWithDefault("WEAVIATE_SCHEME", db.DefaultWeaviateScheme)

	log.Println("Connecting to Weaviate...")
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}
	defer db.CloseWeaviate()

	// Find the batch to rollback
	batchesCollection := mongoDb.Database.Collection("processing_batches")
	var batch models.ProcessingBatch

	err = batchesCollection.FindOne(ctx, bson.M{"batch_id": batchID}).Decode(&batch)
	if err != nil {
		log.Fatalf("Batch not found: %v", err)
	}

	if batch.Status == "rolled_back" {
		log.Printf("Batch %s is already rolled back", batchID)
		return
	}

	log.Printf("Rolling back batch: %s (%s)", batch.BatchID, batch.Description)
	log.Printf("Created: %s", batch.CreatedAt.Format(time.RFC3339))
	log.Printf("Records to remove: %d total", batch.CreatedRecords.TotalRecords)

	// Rollback MongoDB records
	rollbackCount := 0

	// Remove chapters
	if len(batch.CreatedRecords.ChapterIDs) > 0 {
		result, err := mongoDb.Database.Collection("immunology_chapters").DeleteMany(ctx,
			bson.M{"_id": bson.M{"$in": batch.CreatedRecords.ChapterIDs}})
		if err != nil {
			log.Printf("Warning: Failed to delete chapters: %v", err)
		} else {
			log.Printf("✅ Removed %d chapters", result.DeletedCount)
			rollbackCount += int(result.DeletedCount)
		}
	}

	// Remove case studies
	if len(batch.CreatedRecords.CaseStudyIDs) > 0 {
		result, err := mongoDb.Database.Collection("case_studies").DeleteMany(ctx,
			bson.M{"_id": bson.M{"$in": batch.CreatedRecords.CaseStudyIDs}})
		if err != nil {
			log.Printf("Warning: Failed to delete case studies: %v", err)
		} else {
			log.Printf("✅ Removed %d case studies", result.DeletedCount)
			rollbackCount += int(result.DeletedCount)
		}
	}

	// Remove medical terms
	if len(batch.CreatedRecords.MedicalTermIDs) > 0 {
		result, err := mongoDb.Database.Collection("medical_terms").DeleteMany(ctx,
			bson.M{"_id": bson.M{"$in": batch.CreatedRecords.MedicalTermIDs}})
		if err != nil {
			log.Printf("Warning: Failed to delete medical terms: %v", err)
		} else {
			log.Printf("✅ Removed %d medical terms", result.DeletedCount)
			rollbackCount += int(result.DeletedCount)
		}
	}

	// Remove skills
	if len(batch.CreatedRecords.SkillIDs) > 0 {
		result, err := mongoDb.Database.Collection("skills").DeleteMany(ctx,
			bson.M{"_id": bson.M{"$in": batch.CreatedRecords.SkillIDs}})
		if err != nil {
			log.Printf("Warning: Failed to delete skills: %v", err)
		} else {
			log.Printf("✅ Removed %d skills", result.DeletedCount)
			rollbackCount += int(result.DeletedCount)
		}
	}

	// Remove Weaviate vectors
	if len(batch.CreatedRecords.WeaviateIDs) > 0 {
		weaviateCount := 0
		for _, vectorID := range batch.CreatedRecords.WeaviateIDs {
			if err := db.WeaviateDeleteDocument(ctx, "Skills", vectorID); err != nil {
				log.Printf("Warning: Failed to delete Weaviate vector %s: %v", vectorID, err)
			} else {
				weaviateCount++
			}
		}
		log.Printf("✅ Removed %d Weaviate vectors", weaviateCount)
	}

	// Mark batch as rolled back
	now := time.Now()
	_, err = batchesCollection.UpdateOne(ctx,
		bson.M{"_id": batch.ID},
		bson.M{
			"$set": bson.M{
				"status":         "rolled_back",
				"rolled_back_at": now,
			},
		})
	if err != nil {
		log.Printf("Warning: Failed to update batch status: %v", err)
	}

	log.Printf("✅ Rollback completed for batch %s", batchID)
	log.Printf("📊 Total records removed: %d", rollbackCount)
}
