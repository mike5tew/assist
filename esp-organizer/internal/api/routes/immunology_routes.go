package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"esp-organizer/internal/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

// handleImmunologyUpload delegates to the existing upload handler logic
func handleImmunologyUpload(c *gin.Context) {
	immunologyService := services.NewImmunologyService()

	file, _ := c.FormFile("file")
	uploadJob, err := immunologyService.ProcessChapterUpload(
		c.Request.Context(),
		file,
		c.PostForm("chapter_number"),
		c.PostForm("chapter_title"),
		c.PostForm("source_id"),
		c.PostForm("book_title"),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_job_id": uploadJob.ID.Hex(),
		"message":       "Processing started",
	})
}

// Use this instead of defining locally:
// book := utils.FindBookByID(sourceID)

// Remove the local findBookByID function entirely - it's now in utils

// Helper function to update upload job status
func updateUploadJobStatus(ctx context.Context, jobID primitive.ObjectID, status string, progress float64, detail string) {

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[Coach] MongoDB connection failed: %v", err)
		return // Early return on DB connection failure
	}
	defer mongoDb.Client.Disconnect(ctx)
	uploadJobsCollection := mongoDb.Database.Collection("upload_jobs")

	update := bson.M{
		"$set": bson.M{
			"status":        status,
			"progress":      progress,
			"status_detail": detail,
			"updated_at":    time.Now(),
		},
	}

	uploadJobsCollection.UpdateOne(ctx, bson.M{"_id": jobID}, update)
}

// // Helper function to create or find a source
// func createOrFindSource(ctx context.Context, bookTitle string) (primitive.ObjectID, error) {
// 	mongoDb, err := db.NewFromEnv()
// 	if err != nil {
// 		log.Printf("[Coach] MongoDB connection failed: %v", err)
// 		return primitive.NilObjectID, err
// 	}
// 	defer mongoDb.Client.Disconnect(ctx)
// 	sourcesCollection := mongoDb.Database.Collection("sources")

// 	// Try to find existing source
// 	var existingSource models.Source
// 	err = sourcesCollection.FindOne(ctx, bson.M{"title": bookTitle}).Decode(&existingSource)
// 	if err == nil {
// 		return existingSource.ID, nil
// 	}

// 	// Create new source
// 	newSource := models.Source{
// 		ID:        primitive.NewObjectID(),
// 		Title:     bookTitle,
// 		Type:      "book",
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	_, err = sourcesCollection.InsertOne(ctx, newSource)
// 	if err != nil {
// 		return primitive.NilObjectID, err
// 	}

// 	return newSource.ID, nil
// }
