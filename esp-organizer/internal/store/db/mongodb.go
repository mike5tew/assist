package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionTimeout = 10 * time.Second
	operationTimeout  = 5 * time.Second
)

// MongoDBConfig holds configuration from environment variables
type MongoDBConfig struct {
	URI      string
	Database string
	TLS      bool
}

// NewFromEnv creates a new MongoDB connection from environment variables.
// This is now the single source of truth for creating a MongoDB connection.

func LoadMongoConfig() MongoDBConfig {
	// The Makefile now sets MONGODB_URI to MONGODB_URI_EXTERNAL for host commands.
	// The docker-compose.yml sets MONGODB_URI to MONGODB_URI_INTERNAL for the API service.
	// This function simply reads the final MONGODB_URI value.
	uri := os.Getenv("MONGODB_URI")

	// Extract database name from URI
	dbName := GetEnvWithDefault("MONGO_DB_NAME", "esp_organizer") // Default
	if parsedURI, err := url.Parse(uri); err == nil {
		pathDB := strings.TrimPrefix(parsedURI.Path, "/")
		if pathDB != "" {
			dbName = pathDB
		}
	}

	return MongoDBConfig{
		URI:      uri,
		Database: dbName,
		TLS:      getEnvBool("MONGO_TLS_ENABLED"),
	}
}

// Helper functions
func GetEnvWithDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getEnvBool(key string) bool {
	return os.Getenv(key) == "true"
}

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	ctx      context.Context
}

// NewMongoDB creates a new MongoDB connection wrapper
// func NewMongoDB(uri string) (*MongoDB, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
// 	defer cancel()

// 	clientOptions := options.Client().ApplyURI(uri)
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
// 	}

// 	// Verify connection
// 	if err := client.Ping(ctx, nil); err != nil {
// 		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
// 	}

// 	// Get database name from environment or extract from URI
// 	dbName := os.Getenv("MONGO_DB_NAME")
// 	if dbName == "" {
// 		// Extract database name from URI as fallback
// 		parts := strings.Split(uri, "/")
// 		if len(parts) > 3 {
// 			dbParts := strings.Split(parts[3], "?")
// 			dbName = dbParts[0]
// 		}
// 	}

// 	// Final fallback to default database name
// 	if dbName == "" {
// 		dbName = "esp_project"
// 		log.Printf("Warning: Using default database name: %s", dbName)
// 	}

// 	db := client.Database(dbName)
// 	log.Printf("✅ Successfully connected to MongoDB database: %s", dbName)

// 	log.Printf((".env variables loaded (in main.go): MONGO_DB_NAME=%s, MONGO_HOST=%s, MONGO_TLS_ENABLED=%t"), dbName, os.Getenv("MONGO_HOST"), getEnvBool("MONGO_TLS_ENABLED"))
// 	return &MongoDB{
// 		Client:   client,
// 		Database: db,
// 	}, nil
// }

// CollectionExists checks if a collection exists
func (m *MongoDB) CollectionExists(ctx context.Context, name string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	collections, err := m.Database.ListCollectionNames(ctx, bson.M{"name": name})
	if err != nil {
		return false, fmt.Errorf("failed to list collections: %w", err)
	}

	for _, coll := range collections {
		if coll == name {
			return true, nil
		}
	}
	return false, nil
}

// DeleteDocument deletes a document by ID
func (m *MongoDB) DeleteDocument(ctx context.Context, collectionName, id string) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	collection := m.GetCollection(collectionName)
	res, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("delete operation failed: %w", err)
	}

	if res.DeletedCount == 0 {
		return errors.New("document not found")
	}

	log.Printf("Deleted document %s from %s", id, collectionName)
	return nil
}

// Close safely closes the MongoDB connection
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	if err := m.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	return nil
}

func MongoDeleteDocument(ctx context.Context, db *MongoDB, collectionName, id string) error {
	if db == nil {
		return errors.New("MongoDB client is nil")
	}

	return db.DeleteDocument(ctx, collectionName, id)
}

func MongoCollectionExists(ctx context.Context, db *MongoDB, collectionName string) (bool, error) {
	if db == nil {
		return false, errors.New("MongoDB client is nil")
	}

	return db.CollectionExists(ctx, collectionName)
}

// InitializeWeaviateWithRetry attempts to connect to Weaviate with a retry mechanism
func InitializeWeaviateWithRetry(maxRetries int, delay time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = InitializeWeaviateFromEnv()
		if err == nil {
			log.Printf("✅ Successfully connected to Weaviate on attempt %d/%d", i+1, maxRetries)
			return nil
		}
		log.Printf("⚠️ Weaviate connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		time.Sleep(delay)
	}
	return fmt.Errorf("failed to initialize Weaviate after %d attempts: %w", maxRetries, err)
}

// StoreDocumentWithVectorRetry wraps the original function with retry logic for leader not found errors
func StoreDocumentWithVectorRetry(ctx context.Context, className string, properties map[string]interface{}, vector []float32, maxRetries int, delay time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = StoreDocumentWithVector(ctx, className, properties, vector)
		if err == nil {
			return nil
		}

		// Check if it's a "leader not found" error
		if i < maxRetries-1 && (err.Error() == "leader not found" || err.Error() == "failed to execute query: leader not found" || err.Error() == "invalid object: failed to execute query: leader not found") {
			log.Printf("⚠️ Weaviate leader not found on attempt %d/%d, retrying in %v...", i+1, maxRetries, delay)
			time.Sleep(delay * time.Duration(i+1)) // Exponential backoff
			continue
		}

		// If it's a different error, return immediately
		return err
	}
	return fmt.Errorf("failed to store document after %d attempts: %w", maxRetries, err)
}

// EnsureSemanticLinkIndexes creates optimized indexes
func EnsureSemanticLinkIndexes(ctx context.Context, collection *mongo.Collection) error {
	indexes := []mongo.IndexModel{
		// Unique composite key
		{
			Keys:    bson.D{{Key: "composite_key", Value: 1}},
			Options: options.Index().SetUnique(true),
		},

		// 🆕 Case study queries
		{
			Keys:    bson.D{{Key: "case_study_id", Value: 1}},
			Options: options.Index(),
		},

		// 🆕 Disease category queries
		{
			Keys:    bson.D{{Key: "patient_metadata.disease_category", Value: 1}},
			Options: options.Index(),
		},

		// 🆕 Patient name queries
		{
			Keys:    bson.D{{Key: "patient_metadata.patient_name", Value: 1}},
			Options: options.Index(),
		},

		// 🆕 Combined queries (e.g., "XLA patient measurements")
		{
			Keys: bson.D{
				{Key: "patient_metadata.disease_category", Value: 1},
				{Key: "relation_type", Value: 1},
			},
			Options: options.Index(),
		},
	}

	for _, index := range indexes {
		_, err := collection.Indexes().CreateOne(ctx, index)
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	log.Println("✅ Created optimized indexes for semantic_links collection")
	return nil
}
