package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateMisunderstandingIndexes creates indexes for efficient querying of misconceptions
func CreateMisunderstandingIndexes(db *mongo.Database) error {
	collection := db.Collection("misunderstanding_map")

	indexes := []mongo.IndexModel{
		// Index for domain + error_type queries
		{
			Keys: bson.D{
				{Key: "domain", Value: 1},
				{Key: "error_type", Value: 1},
			},
		},
		// Index for finding by wrong_answer (for similarity matching)
		{
			Keys: bson.D{
				{Key: "wrong_answer", Value: "text"},
				{Key: "symptom_pattern", Value: "text"},
			},
		},
		// Index for age group lookups
		{
			Keys: bson.D{
				{Key: "domain", Value: 1},
				{Key: "age_group", Value: 1},
				{Key: "frequency", Value: -1}, // Most frequent first
			},
		},
		// Index for question_keywords search
		{
			Keys: bson.D{
				{Key: "question_keywords", Value: 1},
			},
		},
		// TTL index: Auto-delete entries older than 365 days
		{
			Keys: bson.D{
				{Key: "last_observed", Value: 1},
			},
			Options: options.Index().SetExpireAfterSeconds(365 * 24 * 60 * 60), // Convert to seconds
		},
	}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := collection.Indexes().CreateMany(context.Background(), indexes, opts)
	if err != nil {
		log.Printf("Error creating misunderstanding indexes: %v", err)
		return err
	}

	log.Println("✅ Misunderstanding map indexes created successfully")
	return nil
}
