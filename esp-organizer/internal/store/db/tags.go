package db

import (
	"context"
	"esp-organizer/internal/models"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

// FindContentByTags searches for SubjectContent documents that contain all of the specified tags.
func FindContentByTags(ctx context.Context, tags []string) ([]models.SubjectContent, error) {
	mngoDB, err := NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mngoDB.Client.Disconnect(ctx); err != nil {
			log.Printf("Warning: Error disconnecting from MongoDB: %v", err)
		}
	}()

	collection := mngoDB.Client.Database("esp-organizer").Collection("subject_content")

	// Use the $all operator to find documents where the tags array field contains all specified values.
	// Use $in if you want to find documents that contain any of the specified tags.
	filter := bson.M{"tags": bson.M{"$all": tags}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find query: %w", err)
	}
	defer cursor.Close(ctx)

	var results []models.SubjectContent
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}

	return results, nil
}
