package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}

	fmt.Println("🔍 Verifying AWS Bedrock configuration for Weaviate...")

	// Check AWS credentials
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	embeddingModelID := os.Getenv("AWS_BEDROCK_EMBEDDING_MODEL_ID")

	fmt.Printf("AWS Region: %s\n", awsRegion)
	fmt.Printf("AWS Access Key ID: %s...\n", maskString(awsAccessKey))
	fmt.Printf("AWS Secret Key: %s...\n", maskString(awsSecretKey))
	fmt.Printf("AWS Bedrock Embedding Model ID: %s\n", embeddingModelID)

	// Validate AWS credentials
	if awsAccessKey == "" || awsSecretKey == "" {
		fmt.Println("❌ AWS credentials are missing!")
		fmt.Println("Please set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY in your .env file.")
		os.Exit(1)
	}

	if awsRegion == "" {
		fmt.Println("❌ AWS_REGION is missing!")
		fmt.Println("Please set AWS_REGION in your .env file.")
		os.Exit(1)
	}

	if embeddingModelID == "" {
		fmt.Println("⚠️ AWS_BEDROCK_EMBEDDING_MODEL_ID is not set.")
		fmt.Println("Defaulting to 'amazon.titan-embed-text-v2:0'...")
		embeddingModelID = "amazon.titan-embed-text-v2:0"
	}

	// Test AWS connectivity
	fmt.Println("\nTesting AWS Bedrock connectivity...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(awsRegion),
	)

	if err != nil {
		fmt.Printf("❌ Failed to load AWS config: %v\n", err)
		os.Exit(1)
	}

	// Create Bedrock client
	bedrockClient := bedrock.NewFromConfig(cfg)

	// List Bedrock foundation models to verify access
	resp, err := bedrockClient.ListFoundationModels(ctx, &bedrock.ListFoundationModelsInput{})
	if err != nil {
		fmt.Printf("❌ Failed to connect to AWS Bedrock: %v\n", err)
		fmt.Println("Please check your AWS credentials and region.")
		os.Exit(1)
	}

	fmt.Println("✅ Successfully connected to AWS Bedrock!")
	fmt.Printf("Found %d foundation models\n", len(resp.ModelSummaries))

	// Verify that the embedding model exists
	modelExists := false
	for _, model := range resp.ModelSummaries {
		if model.ModelId != nil && *model.ModelId == embeddingModelID {
			modelExists = true
			break
		}
	}

	if modelExists {
		fmt.Printf("✅ Model '%s' exists and is accessible\n", embeddingModelID)
	} else {
		fmt.Printf("⚠️ Model '%s' was not found in available models\n", embeddingModelID)
		fmt.Println("Available models:")
		for _, model := range resp.ModelSummaries {
			if model.ModelId != nil {
				fmt.Printf("- %s\n", *model.ModelId)
			}
		}
	}

	fmt.Println("\n✅ AWS Bedrock verification complete")
}

// maskString masks all but the first few characters of a string
func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:4] + "****"
}
