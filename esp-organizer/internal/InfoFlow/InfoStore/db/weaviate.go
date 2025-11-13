package db

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"esp-organizer/internal/models"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	wvmodels "github.com/weaviate/weaviate/entities/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB connection timeout
const mongoConnectionTimeout = 10 * time.Second

// NewMongoDBFromEnv creates a new MongoDB connection using environment variables
func NewMongoDBFromEnv() (*MongoDB, error) {
	// .env loading is now handled by config.LoadConfig() at application startup.
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable not set")
	}

	log.Printf("Connecting to MongoDB using URI from environment...")

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectionTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)

	// Optional TLS configuration
	if os.Getenv("MONGODB_TLS") == "true" {
		clientOptions.SetTLSConfig(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping with timeout
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Extract database name from URI or use default
	parsedURI, err := url.Parse(mongoURI)
	if err != nil {
		return nil, fmt.Errorf("invalid MongoDB URI: %w", err)
	}
	dbName := strings.TrimPrefix(parsedURI.Path, "/")
	if dbName == "" {
		dbName = "esp_organizer" // Fallback
	}

	log.Printf("✅ Successfully connected to MongoDB database: %s", dbName)

	return &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

// // Close closes the MongoDB connection
// func (m *MongoDB) Close() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	return m.Client.Disconnect(ctx)
// }

// Helper function to get a pointer to a boolean value.
func boolPtr(b bool) *bool {
	return &b
}

var (
	weaviateClient *weaviate.Client
	weaviateOnce   sync.Once
)

func CreateWeaviateMetadata(mongoID primitive.ObjectID, contentType, subject string) map[string]interface{} {
	return map[string]interface{}{
		"mongo_id":     mongoID.Hex(), // Convert ObjectID to string for Weaviate
		"content_type": contentType,   // e.g. "case_study", "section", "medical_term"
		"subject":      subject,       // e.g. "immunology", "neurology"
	}
}

func InitializeWeaviateFromEnv() error {
	// .env loading is now handled by config.LoadConfig() at application startup.
	weaviateURL := os.Getenv("WEAVIATE_URL")
	if weaviateURL == "" {
		log.Println("WEAVIATE_URL not set, using default http://localhost:8081")
		weaviateURL = "http://localhost:8081"
	}

	weaviateAPIKey := os.Getenv("WEAVIATE_API_KEY")
	// Optional API Key

	// Call the existing InitializeWeaviate function with the loaded config
	if err := InitializeWeaviate(weaviateURL, weaviateAPIKey); err != nil {
		return err
	}

	// Ensure core classes exist (auto-schema is disabled)
	if err := EnsureCoreSchema(context.Background()); err != nil {
		log.Printf("Warning: failed to ensure core Weaviate schema: %v", err)
	}
	return nil
}

// InitializeWeaviate sets up the connection to Weaviate and returns the client.
func InitializeWeaviate(fullURL string, apiKey string) error {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		log.Printf("Failed to parse Weaviate URL: %v", err)
		return err
	}

	cfg := weaviate.Config{
		Host:   parsedURL.Host,   // e.g., "localhost:8081"
		Scheme: parsedURL.Scheme, // e.g., "http"
	}

	// Add API key authentication if provided
	if apiKey != "" {
		cfg.AuthConfig = auth.ApiKey{Value: apiKey}
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return err
	}

	// Check connection
	_, err = client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		log.Printf("Failed to connect to Weaviate: %v", err)
		return err
	}

	weaviateClient = client
	log.Println("✅ Successfully connected to Weaviate.")
	return nil
}

// EnsureCoreSchema creates required classes if missing (idempotent).
func EnsureCoreSchema(ctx context.Context) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}
	// EducationalSkills and SubjectAreaContent classes are deprecated.
	// Raw content is stored in MongoDB. Vector search is performed on semantic links.

	// Ensure the core SemanticLinks class exists.
	if err := CreateSemanticLinksClass(ctx); err != nil {
		return err
	}
	// Medical excerpts are a specific type of content chunk, also to be deprecated
	// in favor of semantic link searching. For now, we keep it for legacy queries.
	if err := CreateMedicalExcerptClass(ctx); err != nil {
		return err
	}
	return nil
}

// GetWeaviateClient initializes and returns a singleton Weaviate client.
func GetWeaviateClient() *weaviate.Client {
	weaviateOnce.Do(func() {
		client, err := NewWeaviateClientFromEnv()
		if err != nil {
			log.Fatalf("FATAL: Failed to initialize Weaviate client: %v", err)
			return
		}
		weaviateClient = client
	})
	return weaviateClient
}

// NewWeaviateClientFromEnv creates a new Weaviate client from environment variables.
// Added retry logic to handle delayed service availability.
func NewWeaviateClientFromEnv() (*weaviate.Client, error) {
	weaviateURL := os.Getenv("WEAVIATE_URL")
	if weaviateURL == "" {
		weaviateURL = "http://localhost:8081" // Default for local development
	}
	parsedURL, err := url.Parse(weaviateURL)
	if err != nil {
		return nil, err
	}
	cfg := weaviate.Config{
		Host:   parsedURL.Host,
		Scheme: parsedURL.Scheme,
	}
	client := weaviate.New(cfg)

	// Retry readiness check up to 5 times with 5-second intervals.
	var checkErr error
	for i := 0; i < 5; i++ {
		_, checkErr = client.Misc().ReadyChecker().Do(context.Background())
		if checkErr == nil {
			return client, nil
		}
		log.Printf("Weaviate not ready (attempt %d): %v; retrying in 5 sec", i+1, checkErr)
		time.Sleep(5 * time.Second)
	}
	return nil, checkErr
}

// CloseWeaviate is a placeholder for any cleanup logic if needed in the future.
func CloseWeaviate() {
	// No explicit disconnect needed for the Weaviate client
}

// WeaviateCollectionExists checks if a class exists in Weaviate
func WeaviateCollectionExists(ctx context.Context, className string) (bool, error) {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return false, fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	// Check if the class exists in the schema
	schema, err := weaviateClient.Schema().Getter().Do(ctx)
	if err != nil {
		return false, err
	}

	// Look for the class in the schema
	for _, class := range schema.Classes {
		if class.Class == className {
			return true, nil
		}
	}

	return false, nil
}

// WeaviateDeleteDocument deletes an object from Weaviate by ID
func WeaviateDeleteDocument(ctx context.Context, className string, id string) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	err := weaviateClient.Data().Deleter().
		WithClassName(className).
		WithID(id).
		Do(ctx)

	if err != nil {
		log.Printf("Failed to delete object %s from class %s: %v", id, className, err)
		return err
	}

	log.Printf("Successfully deleted object %s from class %s", id, className)
	return nil
}

// ResetWeaviateClass deletes and recreates a class with the given definition
func ResetWeaviateClass(ctx context.Context, className string) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	// Check if class exists
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}

	// Delete the class if it exists
	if exists {
		log.Printf("Deleting existing class: %s", className)
		err = weaviateClient.Schema().ClassDeleter().WithClassName(className).Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete class %s: %w", className, err)
		}
		log.Printf("Successfully deleted class %s", className)
	}

	// Recreate the class based on className
	switch className {
	// "EducationalSkills" & "SubjectAreaContent" are deprecated and removed.
	case "SemanticLinks":
		return CreateSemanticLinksClass(ctx)
	case "MedicalExcerpt":
		return CreateMedicalExcerptClass(ctx)
	default:
		return fmt.Errorf("no class definition available for %s", className)
	}
}

// StoreDocumentWithVector stores a document in Weaviate with optional vector embedding
func StoreDocumentWithVector(ctx context.Context, className string, properties map[string]interface{}, vector []float32) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	// Add validation for empty vector
	if len(vector) == 0 {
		return fmt.Errorf("empty vector provided, cannot store document")
	}

	// Log vector size for debugging
	log.Printf("Storing document in Weaviate class %s with vector of length %d", className, len(vector))

	// Filter out application-specific timestamps that might conflict with Weaviate schema
	properties = FilterWeaviateProperties(properties)

	// Fail fast if class does not exist
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("weaviate class %q does not exist", className)
	}

	// Create an object builder
	objectBuilder := weaviateClient.Data().Creator().
		WithClassName(className).
		WithProperties(properties)

	objectBuilder = objectBuilder.WithVector(vector)

	_, err = objectBuilder.Do(ctx)
	if err != nil {
		// Check for vector dimension mismatch error
		if strings.Contains(err.Error(), "new node has a vector with length") &&
			strings.Contains(err.Error(), "Existing nodes have vectors with length") {
			log.Printf("❌ VECTOR DIMENSION MISMATCH: %v", err)
			log.Printf("Vector dimensions from AWS Bedrock: %d", len(vector))
			log.Printf("Run 'go run cmd/weaviate-reset/main.go -class %s' to reset this class", className)
			return fmt.Errorf("vector dimension mismatch: %w", err)
		}
		log.Printf("Failed to store document in Weaviate class %s: %v", className, err)
		return err
	}

	log.Printf("✅ Successfully stored document in Weaviate class %s with vector of length %d", className, len(vector))
	return nil
}

// FilterWeaviateProperties removes fields that shouldn't go to Weaviate
func FilterWeaviateProperties(properties map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range properties {
		// Skip application-specific metadata fields
		if key == "created_at" || key == "updated_at" || key == "app_metadata" {
			continue
		}

		// Handle nested maps
		if nestedMap, ok := value.(map[string]interface{}); ok {
			result[key] = FilterWeaviateProperties(nestedMap)
		} else {
			result[key] = value
		}
	}

	return result
}

// StoreDocumentWithoutVector stores a document in Weaviate without providing a vector
func StoreDocumentWithoutVector(ctx context.Context, className string, properties map[string]interface{}) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	// Add timestamps to properties
	properties = AddTimestampToProperties(properties)

	// Create an object builder
	objectBuilder := weaviateClient.Data().Creator().
		WithClassName(className).
		WithProperties(properties)

	// Create the object in Weaviate - Weaviate will generate the vector automatically
	_, err := objectBuilder.Do(ctx)
	if err != nil {
		log.Printf("Failed to store document in Weaviate class %s: %v", className, err)
		return err
	}

	return nil
}

// CreateEducationalSkillsClass is deprecated. Skills are stored in MongoDB.
// Vector search is performed on the SemanticLinks class.

// CreateSubjectAreaContentClass is deprecated. Raw content is stored in MongoDB and
// only the semantic links extracted from it are vectorized and stored in Weaviate.
func CreateSubjectAreaContentClass(ctx context.Context) error {
	className := "SubjectAreaContent"
	log.Printf("DEPRECATION WARNING: CreateSubjectAreaContentClass for class '%s' is called, but this class is deprecated. Vector search should be performed on the 'SemanticLinks' class.", className)
	return nil
}

// Add this function to check and reinitialize the client if needed
func EnsureWeaviateClient(ctx context.Context) error {
	// If client already exists, check connection
	if weaviateClient != nil {
		_, err := weaviateClient.Misc().ReadyChecker().Do(ctx)
		if err == nil {
			// Connection is working
			return nil
		}
		// Connection is broken, will attempt reinitialization
		log.Println("Weaviate client connection lost, attempting to reinitialize...")
	}

	// Reinitialize from environment variables
	return InitializeWeaviateFromEnv()
}

// CreateSemanticLinksClass creates the class for semantic links in Weaviate.
func CreateSemanticLinksClass(ctx context.Context) error {
	className := "SemanticLinks"
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("%s class already exists in Weaviate", className)
		return nil
	}

	// --- Configuration from Environment ---
	// Get WEAVIATE_VECTORIZER (e.g., "text2vec-aws") and model ID
	vectorizerName := os.Getenv("WEAVIATE_VECTORIZER")
	modelName := os.Getenv("AWS_BEDROCK_EMBEDDING_MODEL_ID")
	//serviceName := "bedrock" // This is the service name for the AWS module

	if vectorizerName == "" || modelName == "" {
		log.Println("WARNING: WEAVIATE_VECTORIZER or AWS_BEDROCK_EMBEDDING_MODEL_ID not set. Using 'none' vectorizer for SemanticLinks. Vector search will NOT work.")
		vectorizerName = "none"
	} else {
		log.Printf("Configuring %s class with vectorizer %s (Model: %s)", className, vectorizerName, modelName)
	}
	// --------------------------------------

	classObj := &wvmodels.Class{
		Class:      "SemanticLinks",
		Vectorizer: vectorizerName,
		Properties: []*wvmodels.Property{
			// 🔑 BOTH fields are indexed and searchable
			{Name: "source_term", DataType: []string{"text"}, IndexSearchable: boolPtr(true)},
			{Name: "target_term", DataType: []string{"text"}, IndexSearchable: boolPtr(true)},
			{Name: "relation_type", DataType: []string{"text"}, IndexFilterable: boolPtr(true)},
			{Name: "context", DataType: []string{"text"}, IndexSearchable: boolPtr(true)},

			// 🆕 RELATIVE ADDRESSING: Properties that enable dynamic hierarchy computation
			{Name: "source_term_generality", DataType: []string{"number"}, IndexFilterable: boolPtr(true)}, // 0.0-1.0
			{Name: "target_term_generality", DataType: []string{"number"}, IndexFilterable: boolPtr(true)}, // 0.0-1.0
			{Name: "semantic_distance", DataType: []string{"number"}, IndexFilterable: boolPtr(true)},      // 0.0-1.0
			{Name: "relationship_strength", DataType: []string{"number"}, IndexFilterable: boolPtr(true)},  // 0.0-1.0

			// Keep parent/child for graph connectivity
			{Name: "is_parent_of", DataType: []string{"text[]"}, IndexFilterable: boolPtr(true)},
			{Name: "is_child_of", DataType: []string{"text[]"}, IndexFilterable: boolPtr(true)},

			{Name: "context", DataType: []string{"text"}, IndexSearchable: boolPtr(true)},
			{Name: "confidence", DataType: []string{"number"}, IndexFilterable: boolPtr(true)},
			{Name: "created_at", DataType: []string{"date"}, IndexFilterable: boolPtr(true)},
			{Name: "updated_at", DataType: []string{"date"}, IndexFilterable: boolPtr(true)},
		},
	}

	err = weaviateClient.Schema().ClassCreator().WithClass(classObj).Do(ctx)
	if err != nil {
		log.Printf("Failed to create %s class: %v", className, err)
		return err
	}
	log.Printf("Successfully created %s class in Weaviate", className)
	return nil
}

// CreateMedicalExcerptClass creates the class for medical excerpts in Weaviate.
func CreateMedicalExcerptClass(ctx context.Context) error {
	className := "MedicalExcerpt"
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("%s class already exists in Weaviate", className)
		return nil
	}

	// --- Configuration from Environment ---
	// Get WEAVIATE_VECTORIZER (e.g., "text2vec-aws") and model ID
	vectorizerName := os.Getenv("WEAVIATE_VECTORIZER")
	modelName := os.Getenv("AWS_BEDROCK_EMBEDDING_MODEL_ID")
	serviceName := "bedrock"

	if vectorizerName == "" || modelName == "" {
		log.Println("WARNING: WEAVIATE_VECTORIZER or AWS_BEDROCK_EMBEDDING_MODEL_ID not set. Using 'none' vectorizer for MedicalExcerpt. Vector search will NOT work.")
		vectorizerName = "none"
	} else {
		log.Printf("Configuring %s class with vectorizer %s (Model: %s)", className, vectorizerName, modelName)
	}
	// --------------------------------------

	classObj := &wvmodels.Class{
		Class: className,

		// 🛑 CORRECTED: Use the configured vectorizer
		Vectorizer: vectorizerName,

		// 🛑 ADDED: Module config for Bedrock integration
		ModuleConfig: map[string]interface{}{
			vectorizerName: map[string]interface{}{
				"model":   modelName,
				"service": serviceName,
			},
		},
		Properties: []*wvmodels.Property{
			{Name: "mongo_id", DataType: []string{"text"}, IndexFilterable: boolPtr(true)},
			{Name: "excerpt_text", DataType: []string{"text"}, IndexSearchable: boolPtr(true)},
			{Name: "source_title", DataType: []string{"text"}, IndexFilterable: boolPtr(true)},
			{Name: "source_isbn", DataType: []string{"text"}, IndexFilterable: boolPtr(true)},
			{Name: "chapter_number", DataType: []string{"text"}, IndexFilterable: boolPtr(true)},
			{Name: "chapter_title", DataType: []string{"text"}, IndexFilterable: boolPtr(true)},
			{Name: "subject", DataType: []string{"text"}, IndexFilterable: boolPtr(true)},
			{Name: "target_exam", DataType: []string{"text"}, IndexFilterable: boolPtr(true)},
			{Name: "created_at", DataType: []string{"date"}, IndexFilterable: boolPtr(true)},
			{Name: "updated_at", DataType: []string{"date"}, IndexFilterable: boolPtr(true)},
		},
	}

	err = weaviateClient.Schema().ClassCreator().WithClass(classObj).Do(ctx)
	if err != nil {
		log.Printf("Failed to create %s class: %v", className, err)
		return err
	}
	log.Printf("Successfully created %s class in Weaviate", className)
	return nil
}

// VectorSearch performs a "near vector" search in Weaviate.
// This replaces the deprecated QueryWeaviateData function.
func VectorSearch(ctx context.Context, vector []float32, className string, limit int) ([]map[string]interface{}, error) {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return nil, fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	// Define the fields to retrieve. These are common across different classes.
	fields := []graphql.Field{
		{Name: "title"},
		{Name: "content"},
		{Name: "contentType"},
		{Name: "domain"},
		{Name: "source_term"},
		{Name: "target_term"},
		{Name: "relation_type"},
		{Name: "_additional", Fields: []graphql.Field{{Name: "distance"}, {Name: "score"}, {Name: "certainty"}}},
	}

	nearVector := weaviateClient.GraphQL().
		NearVectorArgBuilder().
		WithVector(vector)

	queryBuilder := weaviateClient.GraphQL().Get().
		WithClassName(className).
		WithFields(fields...).
		WithNearVector(nearVector).
		WithLimit(limit)

	response, err := queryBuilder.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("weaviate vector search failed for class %s: %w", className, err)
	}

	if response.Errors != nil {
		return nil, fmt.Errorf("weaviate query returned errors: %v", response.Errors)
	}

	data, ok := response.Data["Get"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from Weaviate: 'Get' field is missing")
	}

	items, ok := data[className].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil // No results is not an error
	}

	results := make([]map[string]interface{}, len(items))
	for i, item := range items {
		if itemMap, ok := item.(map[string]interface{}); ok {
			results[i] = itemMap
		}
	}

	return results, nil
}

// QueryWeaviateData performs a vector search on a Weaviate class.
// DEPRECATED: This is a low-level function. Use the SkillService.SearchSkillsByVector for application-level searches,
// as it correctly targets the 'SemanticLinks' class.
func QueryWeaviateData(ctx context.Context, className string, vector []float32, query string, limit int, whereFilter *filters.WhereBuilder) ([]map[string]interface{}, error) {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return nil, fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	fields := []graphql.Field{
		{Name: "name"},
		{Name: "description"},
		{Name: "mongo_id"},
		{Name: "_additional", Fields: []graphql.Field{{Name: "certainty"}, {Name: "id"}}},
	}

	nearVector := weaviateClient.GraphQL().
		NearVectorArgBuilder().
		WithVector(vector)

	queryBuilder := weaviateClient.GraphQL().Get().
		WithClassName(className).
		WithFields(fields...).
		WithNearVector(nearVector).
		WithLimit(limit)

	if whereFilter != nil {
		queryBuilder = queryBuilder.WithWhere(whereFilter)
	}

	response, err := queryBuilder.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("weaviate query failed: %w", err)
	}

	if response.Errors != nil {
		return nil, fmt.Errorf("weaviate query returned errors: %v", response.Errors)
	}

	data, ok := response.Data["Get"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from Weaviate")
	}

	items, ok := data[className].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil // No results is not an error
	}

	results := make([]map[string]interface{}, len(items))
	for i, item := range items {
		if itemMap, ok := item.(map[string]interface{}); ok {
			results[i] = map[string]interface{}{"properties": itemMap}
		}
	}

	return results, nil
}

// GetSkillByID fetches a single skill from MongoDB by its ObjectID.
func GetSkillByID(ctx context.Context, id primitive.ObjectID) (*models.Skill, error) {
	mongoDB, err := NewMongoDBFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer mongoDB.Client.Disconnect(ctx)

	skillsCollection := mongoDB.Database.Collection("skills")
	var skill models.Skill
	err = skillsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&skill)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("skill with ID %s not found", id.Hex())
		}
		return nil, fmt.Errorf("failed to find skill by ID: %w", err)
	}

	return &skill, nil
}

// WeaviateImportSchemaFromFile imports a schema from a JSON file into Weaviate
func WeaviateImportSchemaFromFile(ctx context.Context, filePath string) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open schema file: %w", err)
	}
	defer file.Close()

	var schema wvmodels.Schema
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&schema); err != nil {
		return fmt.Errorf("failed to decode schema JSON: %w", err)
	}

	for _, class := range schema.Classes {
		exists, err := WeaviateCollectionExists(ctx, class.Class)
		if err != nil {
			return err
		}
		if exists {
			log.Printf("Class %s already exists, skipping creation", class.Class)
			continue
		}

		err = weaviateClient.Schema().ClassCreator().WithClass(class).Do(ctx)
		if err != nil {
			log.Printf("Failed to create class %s: %v", class.Class, err)
			return err
		}
		log.Printf("Successfully created class %s", class.Class)
	}

	return nil
}

// HydrateSkill fetches and populates the parent skills for a given skill.
// A skill is defined by its parent components (e.g., "Handwriting" is composed of "Fine Motor Skills").
func HydrateSkill(ctx context.Context, skillsCollection *mongo.Collection, skillID primitive.ObjectID) (*models.Skill, error) {
	var skill models.Skill
	err := skillsCollection.FindOne(ctx, bson.M{"_id": skillID}).Decode(&skill)
	if err != nil {
		return nil, fmt.Errorf("failed to find skill: %w", err)
	}

	// A skill is defined by its parent components. Hydrate ParentSkills from ParentSkillIDs.
	if len(skill.ParentSkillIDs) > 0 {
		cursor, err := skillsCollection.Find(ctx, bson.M{"_id": bson.M{"$in": skill.ParentSkillIDs}})
		if err != nil {
			return nil, fmt.Errorf("failed to find parent skills: %w", err)
		}
		defer cursor.Close(ctx)

		var parentSkills []*models.Skill
		if err := cursor.All(ctx, &parentSkills); err != nil {
			return nil, fmt.Errorf("failed to decode parent skills: %w", err)
		}
		skill.ParentSkills = parentSkills
	}

	// Child skills are not hydrated by default, as they are not part of the skill's definition.
	// The ChildSkillIDs are sufficient for downward traversal if needed elsewhere.

	return &skill, nil
}

// Add the AddTimestampToProperties function to add timestamps to properties
func AddTimestampToProperties(properties map[string]interface{}) map[string]interface{} {
	now := time.Now().UTC().Format(time.RFC3339)

	// Add created_at if not present
	if _, exists := properties["created_at"]; !exists {
		properties["created_at"] = now
	}

	// Always update updated_at
	properties["updated_at"] = now

	return properties
}
