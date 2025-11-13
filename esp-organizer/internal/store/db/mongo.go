package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewFromEnv creates a new MongoDB client from environment variables
func NewFromEnv() (*MongoDB, error) {
	// Get MongoDB URI from environment with fallback
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		// Try different environment variable names
		uri = os.Getenv("MONGO_URI")
	}

	if uri == "" {
		// Final fallback to localhost
		uri = "mongodb://localhost:27017"
		log.Printf("MONGODB_URI not set, using default: %s", uri)
	}

	// Set database name with fallback
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "esp_organizer"
	}

	return New(uri, dbName)
}

// New creates a new MongoDB client
func New(uri, dbName string) (*MongoDB, error) {
	log.Printf("Connecting to MongoDB using URI: %s (database: %s)",
		maskConnectionString(uri), dbName)

	// Configure client with more resilient settings
	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second).
		SetRetryWrites(true).
		SetRetryReads(true)

	// Set up connection context with longer timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect with retry logic
	var client *mongo.Client
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			log.Printf("Retrying MongoDB connection (attempt %d/%d)...", i+1, maxRetries)
			time.Sleep(2 * time.Second * time.Duration(i)) // Exponential backoff
		}

		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			// Ping to verify connection is successful
			pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = client.Ping(pingCtx, readpref.Primary())
			pingCancel()

			if err == nil {
				// Successfully connected
				log.Printf("Successfully connected to MongoDB")
				break
			}
		}

		log.Printf("MongoDB connection attempt %d failed: %v", i+1, err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

// maskConnectionString hides credentials in the connection string for logging
func maskConnectionString(uri string) string {
	// Simple implementation - just return with credentials masked
	return "[masked for security]"
}

// GetCollection is a helper method to get a collection
func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

// Other methods...
