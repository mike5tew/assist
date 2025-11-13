package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// BedrockLlamaClient is a client for interacting with the AWS Bedrock service
type BedrockLlamaClient struct {
	bedrockClient    *bedrockruntime.Client
	embeddingModelID string
}

// NewBedrockLlamaClient creates a new BedrockLlamaClient
func NewBedrockLlamaClient() (*BedrockLlamaClient, error) {
	// Load AWS configuration from environment variables or credentials file
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create the Bedrock runtime client
	bedrockClient := bedrockruntime.NewFromConfig(cfg)

	return &BedrockLlamaClient{
		bedrockClient:    bedrockClient,
		embeddingModelID: "amazon.titan-embed-text-v1", // Default model ID
	}, nil
}

// GenerateEmbedding creates vector embeddings with error handling for empty content
func (c *BedrockLlamaClient) GenerateEmbedding(text string) ([]float32, error) {
	// Validate and clean input text
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("cannot generate embedding for empty text")
	}

	// Ensure text meets minimum length requirement for Bedrock (at least 10 characters)
	if len(text) < 10 {
		// Pad text if too short by repeating it
		repeats := (10 + len(text) - 1) / len(text) // Calculate needed repeats to exceed 10 chars
		var paddedText strings.Builder
		for i := 0; i < repeats; i++ {
			paddedText.WriteString(text)
			if i < repeats-1 {
				paddedText.WriteString(" ")
			}
		}
		text = paddedText.String()
		log.Printf("Text was too short, padded to: %s", text)
	}

	// Set up Bedrock request with logging
	log.Printf("Generating embedding for text of length %d", len(text))
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.embeddingModelID),
		ContentType: aws.String("application/json"),
		Body:        []byte(fmt.Sprintf(`{"inputText": %q}`, text)),
	}

	// Call Bedrock API with error handling
	output, err := c.bedrockClient.InvokeModel(context.Background(), input)
	if err != nil {
		log.Printf("AWS Bedrock API error: %v", err)
		return nil, fmt.Errorf("AWS Bedrock API error: %w", err)
	}

	// Parse response with validation
	var response struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.Unmarshal(output.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse embedding response: %w", err)
	}

	// Validate response
	if len(response.Embedding) == 0 {
		return nil, fmt.Errorf("received empty embedding from AWS Bedrock")
	}

	log.Printf("Successfully generated embedding with %d dimensions", len(response.Embedding))
	return response.Embedding, nil
}
