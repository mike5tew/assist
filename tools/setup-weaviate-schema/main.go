package main

import (
	"context"
	"log"
	"net/url"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	log.Println("Starting Weaviate schema setup...")

	// Get Weaviate URL from environment or use default
	weaviateURL := os.Getenv("WEAVIATE_URL")
	if weaviateURL == "" {
		weaviateURL = "http://localhost:8081" // Default for local development
	}
	log.Printf("Connecting to Weaviate at %s", weaviateURL)

	// Parse the URL to extract host and scheme
	parsedURL, err := url.Parse(weaviateURL)
	if err != nil {
		log.Fatalf("Invalid Weaviate URL: %v", err)
	}

	// Initialize the client
	cfg := weaviate.Config{
		Host:   parsedURL.Host,
		Scheme: parsedURL.Scheme,
		// Remove any authentication headers
		// Headers: map[string]string{
		//     "X-OpenAI-Api-Key": os.Getenv("OPENAI_API_KEY"), // Remove this
		// },
	}

	// cfg := weaviate.Config{
	// 	Host:   parsedURL.Host,   // Only the host:port part
	// 	Scheme: parsedURL.Scheme, // Only the scheme (http/https)
	// 	// No Auth field needed for anonymous access to Weaviate
	// 	Headers: map[string]string{
	// 		"X-OpenAI-Api-Key": os.Getenv("OPENAI_API_KEY"), // Add if needed
	// 	},
	// }
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Weaviate client: %v", err)
	}

	// Check connection
	_, err = client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		log.Fatalf("Weaviate is not ready: %v", err)
	}
	log.Println("✅ Successfully connected to Weaviate")

	// Define the schema classes
	setupSubjectAreaContentClass(client)

	// Verify that the classes now exist
	verifyClasses(client)

	log.Println("🎉 Weaviate schema setup completed successfully!")
}

func setupSubjectAreaContentClass(client *weaviate.Client) {
	className := "SubjectAreaContent"

	// Check if class already exists
	exists, err := classExists(client, className)
	if err != nil {
		log.Fatalf("Error checking if class exists: %v", err)
	}

	if exists {
		log.Printf("Class %s already exists in Weaviate", className)
		return
	}
	vectorizerName := os.Getenv("WEAVIATE_VECTORIZER")
	modelName := os.Getenv("AWS_BEDROCK_EMBEDDING_MODEL_ID") // Get model from ENV
	serviceName := "bedrock"                                 // Assuming you are using Bedrock as per your Docker setup

	if modelName == "" {
		modelName = "amazon.titan-embed-text-v2:0" // Default model
		log.Printf("Warning: AWS_BEDROCK_EMBEDDING_MODEL_ID not set, using default: %s", modelName)
	}

	// Define class schema
	class := &models.Class{
		Class: className,
		Properties: []*models.Property{
			{
				Name:        "title",
				DataType:    []string{"text"},
				Description: "The title of the content",
			},
			{
				Name:        "content",
				DataType:    []string{"text"},
				Description: "The textual content of the subject area",
			},
			{
				Name:        "type",
				DataType:    []string{"text"},
				Description: "The type of content (chapter, case study, term, etc.)",
			},
			{
				Name:        "source_id",
				DataType:    []string{"text"},
				Description: "Reference to the source in MongoDB",
			},
			{
				Name:        "batch_id",
				DataType:    []string{"text"},
				Description: "Batch ID for grouping related content",
			},
			{
				Name:        "created_at",
				DataType:    []string{"date"},
				Description: "When this content was created",
			},
		},
		Description: "Subject area content with semantic vectors for search and retrieval",
		Vectorizer:  vectorizerName,
		ModuleConfig: map[string]interface{}{
			vectorizerName: map[string]interface{}{
				"model":   modelName,
				"service": serviceName,
			},
		},
	}

	// Create the class
	err = client.Schema().ClassCreator().WithClass(class).Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to create %s class: %v", className, err)
	}
	log.Printf("Successfully created %s class in Weaviate", className)
}

func classExists(client *weaviate.Client, className string) (bool, error) {
	// Get the schema
	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		return false, err
	}

	// Check if the class exists
	for _, class := range schema.Classes {
		if class.Class == className {
			return true, nil
		}
	}
	return false, nil
}

func verifyClasses(client *weaviate.Client) {
	// List all classes to verify
	classesToVerify := []string{"SubjectAreaContent", "SemanticLinks", "MedicalExcerpt"}

	// Get the schema
	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to get schema: %v", err)
	}

	// Check each class
	existingClasses := make(map[string]bool)
	for _, class := range schema.Classes {
		existingClasses[class.Class] = true
	}

	log.Println("=== Weaviate Schema Verification ===")
	for _, className := range classesToVerify {
		if existingClasses[className] {
			log.Printf("✅ Class %s exists", className)
		} else {
			log.Printf("❌ Class %s does not exist", className)
		}
	}
}
