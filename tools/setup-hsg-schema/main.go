package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	// Ensure logs directory exists before initializing logging
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			log.Fatalf("Failed to create logs directory: %v", err)
		}
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	fmt.Println("🏗️ Setting up HSG Schema in Weaviate...")

	// Get Weaviate URL - note that in docker-compose we're mapping from 8080 (container) to 8081 (host)
	weaviateURL := os.Getenv("WEAVIATE_URL")
	if weaviateURL == "" {
		weaviateURL = "http://localhost:8081" // Default URL for local dev
		log.Printf("⚠️ WEAVIATE_URL not set in environment, using default: %s", weaviateURL)
	} else {
		log.Printf("ℹ️ Using WEAVIATE_URL from environment: %s", weaviateURL)
		// If we're using localhost in the URL, ensure we're using port 8081 (host port)
		if strings.Contains(weaviateURL, "localhost:") && !strings.Contains(weaviateURL, "localhost:8081") {
			weaviateURL = "http://localhost:8081"
			log.Printf("⚠️ Overriding to use correct host port: %s", weaviateURL)
		}
	}

	// Fix potential URL formatting issues
	if !strings.HasPrefix(weaviateURL, "http://") && !strings.HasPrefix(weaviateURL, "https://") {
		weaviateURL = "http://" + weaviateURL
		log.Printf("ℹ️ Added http:// prefix to URL: %s", weaviateURL)
	}

	// Parse the URL to validate and get components
	parsedURL, err := url.Parse(weaviateURL)
	if err != nil {
		log.Fatalf("❌ Invalid Weaviate URL: %v", err)
	}

	// Try connecting without credentials first (to isolate connectivity issues)
	fmt.Printf("🔄 Testing basic connectivity to %s...\n", parsedURL.Host)
	checkURL := fmt.Sprintf("%s://%s/v1/.well-known/ready", parsedURL.Scheme, parsedURL.Host)
	resp, err := http.Get(checkURL)
	if err != nil {
		fmt.Printf("❌ Basic connectivity test failed: %v\n", err)
		fmt.Println("Checking if Weaviate is running in Docker...")

		// Check Docker container logs
		checkLogs := exec.Command("docker", "logs", "--tail", "20", "esp_weaviate")
		logs, _ := checkLogs.CombinedOutput()
		fmt.Println("--- Last 20 lines of Weaviate container logs ---")
		fmt.Println(string(logs))
		fmt.Println("-----------------------------------------------")

		log.Fatalf("❌ Cannot connect to Weaviate. Make sure the container is running and port 8081 is available.")
	} else {
		resp.Body.Close()
		fmt.Printf("✅ Basic connectivity test succeeded with status %d\n", resp.StatusCode)
	}

	// Initialize Weaviate client with authentication
	fmt.Println("🔄 Initializing Weaviate client...")

	cfg := weaviate.Config{
		Host:   parsedURL.Host,
		Scheme: parsedURL.Scheme,
	}

	// Remove authentication code since we're using anonymous access
	// Just create the client without auth config
	client := weaviate.New(cfg)

	if client == nil {
		log.Fatal("❌ Failed to initialize Weaviate client: client is nil")
	}

	// Test the connection with detailed diagnostics
	fmt.Println("🔄 Testing connection to Weaviate...")
	ctx := context.Background()

	// Add timeout to the context
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	ready, err := client.Misc().ReadyChecker().Do(ctx)
	if err != nil {
		log.Printf("❌ Weaviate connectivity test failed: %v", err)
		log.Printf("   Is Weaviate running and accepting connections?")
		log.Fatal("   Cannot continue without Weaviate connectivity")
	}

	if !ready {
		log.Fatal("❌ Weaviate is not ready")
	}

	log.Printf("✅ Successfully connected to Weaviate at %s", weaviateURL)

	// Create the necessary schema classes
	createHSGSchema(client)

	fmt.Println("✅ HSG Schema setup complete")
}

func createHSGSchema(client *weaviate.Client) {
	ctx := context.Background()

	// Get the AWS region from environment
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "eu-west-2" // Default region if not set
		log.Printf("⚠️ AWS_REGION not set, using default: %s", awsRegion)
	}

	// Get the embedding model ID from environment
	embeddingModelID := os.Getenv("AWS_BEDROCK_EMBEDDING_MODEL_ID")
	if embeddingModelID == "" {
		embeddingModelID = "amazon.titan-embed-text-v2:0" // Default model
		log.Printf("⚠️ AWS_BEDROCK_EMBEDDING_MODEL_ID not set, using default: %s", embeddingModelID)
	}

	// Verify model is in the allowed list
	validModels := map[string]bool{
		"amazon.titan-embed-text-v1":   true,
		"amazon.titan-embed-text-v2:0": true,
		"cohere.embed-english-v3":      true,
		"cohere.embed-multilingual-v3": true,
	}

	if !validModels[embeddingModelID] {
		log.Printf("⚠️ Selected embedding model '%s' not in allowed list. Using amazon.titan-embed-text-v2:0", embeddingModelID)
		embeddingModelID = "amazon.titan-embed-text-v2:0"
	}

	// Set up vectorizer module configuration
	vectorizerConfig := map[string]interface{}{
		"model":  embeddingModelID,
		"region": awsRegion,
	}

	// Create classes with explicit vectorizer configuration
	summaryChunkClass := &models.Class{
		Class:       "SummaryChunk",
		Description: "Summary chunks for hierarchical semantic graph",
		VectorIndexConfig: map[string]interface{}{
			"distance": "cosine",
		},
		Vectorizer: "text2vec-aws",
		ModuleConfig: map[string]interface{}{
			"text2vec-aws": vectorizerConfig,
		},
		Properties: []*models.Property{
			{
				Name:        "content",
				DataType:    []string{"text"},
				Description: "The summary content",
			},
			{
				Name:        "domain",
				DataType:    []string{"string"},
				Description: "The knowledge domain",
			},
			{
				Name:        "chunkId",
				DataType:    []string{"string"},
				Description: "Unique identifier for the chunk",
			},
			{
				Name:            "parentId",
				DataType:        []string{"string"},
				Description:     "ID of parent summary chunk",
				IndexFilterable: ptrBool(true),
				IndexSearchable: ptrBool(true),
			},
			{
				Name:            "level",
				DataType:        []string{"int"},
				Description:     "Hierarchy level in the graph",
				IndexFilterable: ptrBool(true),
			},
			{
				Name:        "createdAt",
				DataType:    []string{"date"},
				Description: "Creation timestamp",
			},
		},
	}

	// Class for semantic links
	semanticLinkClass := &models.Class{
		Class:       "SemanticLink",
		Description: "Semantic links between concepts",
		VectorIndexConfig: map[string]interface{}{
			"distance": "cosine", // Using cosine distance
		},
		// Explicitly disable vectorizer for this class
		Vectorizer: "none", // Set to "none" to disable vectorization during schema creation
		Properties: []*models.Property{
			{
				Name:            "sourceTerm",
				DataType:        []string{"string"},
				Description:     "Source term for the link",
				IndexFilterable: ptrBool(true),
				IndexSearchable: ptrBool(true),
			},
			{
				Name:            "targetTerm",
				DataType:        []string{"string"},
				Description:     "Target term for the link",
				IndexFilterable: ptrBool(true),
				IndexSearchable: ptrBool(true),
			},
			{
				Name:        "sourceId",
				DataType:    []string{"string"},
				Description: "Source document ID",
			},
			{
				Name:        "targetId",
				DataType:    []string{"string"},
				Description: "Target document ID",
			},
			{
				Name:            "relationType",
				DataType:        []string{"string"},
				Description:     "Type of relationship",
				IndexFilterable: ptrBool(true),
			},
			{
				Name:        "confidence",
				DataType:    []string{"number"},
				Description: "Confidence score of the relationship",
			},
			{
				Name:        "context",
				DataType:    []string{"text"},
				Description: "Context where the relationship was found",
			},
			{
				Name:            "domain",
				DataType:        []string{"string"},
				Description:     "Knowledge domain",
				IndexFilterable: ptrBool(true),
			},
			{
				Name:        "createdAt",
				DataType:    []string{"date"},
				Description: "Creation timestamp",
			},
		},
	}

	// Class for subject area content
	subjectContentClass := &models.Class{
		Class:       "SubjectAreaContent",
		Description: "Subject area content for domain-specific knowledge",
		VectorIndexConfig: map[string]interface{}{
			"distance": "cosine", // Using cosine distance
		},
		// Explicitly disable vectorizer for this class
		Vectorizer: "none", // Set to "none" to disable vectorization during schema creation
		Properties: []*models.Property{
			{
				Name:            "title",
				DataType:        []string{"string"},
				Description:     "Content title",
				IndexFilterable: ptrBool(true),
				IndexSearchable: ptrBool(true),
			},
			{
				Name:        "content",
				DataType:    []string{"text"},
				Description: "The main content",
			},
			{
				Name:            "domain",
				DataType:        []string{"string"},
				Description:     "Knowledge domain",
				IndexFilterable: ptrBool(true),
			},
			{
				Name:            "contentType",
				DataType:        []string{"string"},
				Description:     "Type of content (chapter, term, case study, etc.)",
				IndexFilterable: ptrBool(true),
			},
			{
				Name:        "sourceId",
				DataType:    []string{"string"},
				Description: "Source document ID",
			},
			{
				Name:        "createdAt",
				DataType:    []string{"date"},
				Description: "Creation timestamp",
			},
		},
	}

	// Create or update the classes
	createOrUpdateClass(client, ctx, summaryChunkClass)
	createOrUpdateClass(client, ctx, semanticLinkClass)
	createOrUpdateClass(client, ctx, subjectContentClass)

	// Give Weaviate a moment to update its schema
	time.Sleep(2 * time.Second)

	// Verify the classes were created
	verifyClasses(client, ctx, []string{"SummaryChunk", "SemanticLink", "SubjectAreaContent"})
}

func createOrUpdateClass(client *weaviate.Client, ctx context.Context, class *models.Class) {
	// Check if class exists
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(class.Class).Do(ctx)
	if err != nil {
		log.Printf("Error checking if class %s exists: %v", class.Class, err)
		return
	}

	if exists {
		// Update existing class
		log.Printf("Class %s already exists. Skipping creation.", class.Class)
		return
	}

	// Create new class
	err = client.Schema().ClassCreator().WithClass(class).Do(ctx)
	if err != nil {
		log.Printf("❌ Failed to create class %s: %v", class.Class, err)
		return
	}

	log.Printf("✅ Successfully created class: %s", class.Class)
}

func verifyClasses(client *weaviate.Client, ctx context.Context, classNames []string) {
	// Get schema to verify classes exist
	schema, err := client.Schema().Getter().Do(ctx)
	if err != nil {
		log.Printf("❌ Failed to get schema: %v", err)
		return
	}

	fmt.Println("\n🔍 Verifying created classes:")

	foundClasses := make(map[string]bool)
	for _, class := range schema.Classes {
		foundClasses[class.Class] = true
	}

	for _, className := range classNames {
		if foundClasses[className] {
			fmt.Printf("✅ Class %s exists\n", className)
		} else {
			fmt.Printf("❌ Class %s is missing\n", className)
		}
	}

	// Check count of objects in each class (if any)
	for _, className := range classNames {
		if !foundClasses[className] {
			continue
		}

		response, err := client.GraphQL().Aggregate().
			WithClassName(className).
			WithFields(
				graphql.Field{
					Name: "meta",
					Fields: []graphql.Field{
						{Name: "count"},
					},
				},
			).
			Do(ctx)

		if err != nil {
			log.Printf("Error counting objects in %s: %v", className, err)
			continue
		}

		count := 0
		if response != nil && response.Data != nil {
			if aggregate, ok := response.Data["Aggregate"].(map[string]interface{}); ok {
				if classData, ok := aggregate[className].([]interface{}); ok && len(classData) > 0 {
					if firstItem, ok := classData[0].(map[string]interface{}); ok {
						if meta, ok := firstItem["meta"].(map[string]interface{}); ok {
							if countVal, ok := meta["count"].(float64); ok {
								count = int(countVal)
							}
						}
					}
				}
			}
		}

		fmt.Printf("   Class %s contains %d objects\n", className, count)
	}
}

// Helper function for boolean pointers (for IndexFilterable/IndexSearchable flags)
func ptrBool(b bool) *bool {
	return &b
}
