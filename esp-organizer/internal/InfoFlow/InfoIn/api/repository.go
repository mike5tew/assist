package api

import (
	"context"
	"esp-organizer/internal/models"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRepository implements Repository interface using MongoDB
type MongoRepository struct {
	client *mongo.Client
}

// NewRepository creates a new MongoDB-backed repository
func NewRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{
		client: client,
	}
}

const defaultTimeout = 10 * time.Second

// extractionJobsCollection returns the MongoDB collection for extraction jobs
func (r *MongoRepository) extractionJobsCollection() *mongo.Collection {
	return r.client.Database("esp_organizer").Collection("extraction_jobs")
}

// CreateExtractionJob creates a new extraction job record
func (r *MongoRepository) CreateExtractionJob(job *models.ExtractionJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	_, err := r.extractionJobsCollection().InsertOne(ctx, job)
	if err != nil {
		return fmt.Errorf("failed to create extraction job: %w", err)
	}
	return nil
}

// UpdateExtractionJob updates an existing extraction job
func (r *MongoRepository) UpdateExtractionJob(job *models.ExtractionJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     job.Status,
			"message":    job.Message,
			"progress":   job.Progress,
			"updated_at": job.UpdatedAt,
		},
	}

	_, err := r.extractionJobsCollection().UpdateOne(
		ctx,
		bson.M{"id": job.ID},
		update,
	)

	if err != nil {
		return fmt.Errorf("failed to update extraction job: %w", err)
	}
	return nil
}

// GetExtractionJob retrieves an extraction job by ID
func (r *MongoRepository) GetExtractionJob(jobID string) (*models.ExtractionJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var job models.ExtractionJob
	err := r.extractionJobsCollection().FindOne(
		ctx,
		bson.M{"id": jobID},
	).Decode(&job)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("extraction job not found: %s", jobID)
		}
		return nil, fmt.Errorf("failed to get extraction job: %w", err)
	}

	return &job, nil
}
