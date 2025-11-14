package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	className := flag.String("class", "", "Name of the class to reset")
	flag.Parse()

	if *className == "" {
		log.Fatal("Please specify a class name with -class flag")
	}

	log.Printf("Starting Weaviate class reset for: %s", *className)

	// Get Weaviate URL from environment or use default
	weaviateURL := os.Getenv("WEAVIATE_URL")
	if weaviateURL == "" {
		weaviateURL = "http://localhost:8081" // Default for local development
	}

	// Parse the URL to extract host and scheme
	parsedURL, err := url.Parse(weaviateURL)
	if err != nil {
		log.Fatalf("Invalid Weaviate URL: %v", err)
	}

	// Initialize the client
	cfg := weaviate.Config{
		Host:   parsedURL.Host,   // Only the host:port part
		Scheme: parsedURL.Scheme, // Only the scheme (http/https)
		// No Auth field needed for anonymous access to Weaviate
		Headers: map[string]string{
			"X-OpenAI-Api-Key": os.Getenv("OPENAI_API_KEY"), // Add if needed
		},
	}
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

	// Delete the class if it exists
	err = client.Schema().ClassDeleter().WithClassName(*className).Do(context.Background())
	if err != nil {
		log.Printf("Warning: Could not delete class %s: %v", *className, err)
	} else {
		log.Printf("Successfully deleted class: %s", *className)
	}

	// Recreate the class with proper configuration
	if *className == "SubjectAreaContent" {
		setupSubjectAreaContentClass(client)
	} else if *className == "SemanticLinks" {
		setupSemanticLinksClass(client)
	}

	log.Printf("🎉 Weaviate class %s has been reset successfully!", *className)
}

func setupSubjectAreaContentClass(client *weaviate.Client) {
	className := "SubjectAreaContent"
	vectorizerName := os.Getenv("WEAVIATE_VECTORIZER")
	modelName := os.Getenv("AWS_BEDROCK_EMBEDDING_MODEL_ID") // Get model from ENV
	serviceName := "bedrock"
	// Define class schema with proper vectorizer configuration
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
	err := client.Schema().ClassCreator().WithClass(class).Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to create %s class: %v", className, err)
	}
	log.Printf("Successfully created %s class in Weaviate", className)
}

func setupSemanticLinksClass(client *weaviate.Client) {
	className := "SemanticLinks"
	vectorizerName := os.Getenv("WEAVIATE_VECTORIZER")
	modelName := os.Getenv("AWS_BEDROCK_EMBEDDING_MODEL_ID") // Get model from ENV
	serviceName := "bedrock"
	// Define class schema with proper vectorizer configuration
	class := &models.Class{
		Class: className,
		Properties: []*models.Property{
			{
				Name:        "source_term",
				DataType:    []string{"text"},
				Description: "The source term in the semantic link",
			},
			{
				Name:        "target_term",
				DataType:    []string{"text"},
				Description: "The target term in the semantic link",
			},
			{
				Name:        "relationship_type",
				DataType:    []string{"text"},
				Description: "The type of relationship between terms",
			},
			{
				Name:        "confidence_score",
				DataType:    []string{"number"},
				Description: "Confidence score of the semantic link",
			},
			{
				Name:        "context",
				DataType:    []string{"text"},
				Description: "Context where this relationship was found",
			},
			{
				Name:        "batch_id",
				DataType:    []string{"text"},
				Description: "Batch ID for grouping related links",
			},
			{
				Name:        "created_at",
				DataType:    []string{"date"},
				Description: "When this link was created",
			},
		},
		Description: "Semantic links between medical terms and concepts",
		Vectorizer:  vectorizerName,
		ModuleConfig: map[string]interface{}{
			vectorizerName: map[string]interface{}{
				"model":   modelName,
				"service": serviceName,
			},
		}}

	// Create the class
	err := client.Schema().ClassCreator().WithClass(class).Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to create %s class: %v", className, err)
	}
	log.Printf("Successfully created %s class in Weaviate", className)
}
