package api

import (
	"context"
	"esp-organizer/internal/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository is the concrete implementation for MongoDB operations.
type MongoRepository struct {
	collection *mongo.Collection
}

// NewRepository creates a new repository for extraction jobs.
func NewRepository(client *mongo.Client) *MongoRepository {
	collection := client.Database("esp_organizer").Collection("extraction_jobs")
	return &MongoRepository{
		collection: collection,
	}
}

// CreateExtractionJob creates a new job record.
func (r *MongoRepository) CreateExtractionJob(job *models.ExtractionJob) error {
	_, err := r.collection.InsertOne(context.Background(), job)
	return err
}

// UpdateExtractionJob updates an existing job record.
func (r *MongoRepository) UpdateExtractionJob(job *models.ExtractionJob) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": job.ID}
	update := bson.M{
		"$set": bson.M{
			"status":     job.Status,
			"message":    job.Message,
			"progress":   job.Progress,
			"updated_at": job.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": job.CreatedAt,
		},
	}

	_, err := r.collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

// GetExtractionJob retrieves a job by its ID.
func (r *MongoRepository) GetExtractionJob(jobID string) (*models.ExtractionJob, error) {
	var job models.ExtractionJob
	filter := bson.M{"_id": jobID}
	err := r.collection.FindOne(context.Background(), filter).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("job with ID '%s' not found", jobID)
		}
		return nil, err
	}
	return &job, nil
}
