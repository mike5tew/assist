package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// BedrockClient is a client for interacting with AWS Bedrock
type BedrockClient struct {
	client         *bedrockruntime.Client
	config         aws.Config
	modelID        string
	maxInputLength int
}

// EmbeddingResponse represents the structure of Bedrock embedding responses
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"` // Bedrock typically returns float64
	InputText string    `json:"inputText,omitempty"`
}

// NewBedrockClient creates a new BedrockClient
func NewBedrockClient(modelID string) (*BedrockClient, error) {
	// Resolve region from env with sensible defaults
	region := os.Getenv("AWS_BEDROCK_REGION")
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}
	if region == "" {
		region = "eu-west-2"
	}

	// Ignore placeholder session token if present
	if tok := os.Getenv("AWS_SESSION_TOKEN"); tok != "" {
		if strings.Contains(strings.ToLower(tok), "your") || strings.Contains(strings.ToLower(tok), "here") {
			log.Println("⚠️ Ignoring placeholder AWS_SESSION_TOKEN in Bedrock client")
			_ = os.Unsetenv("AWS_SESSION_TOKEN")
		}
	}

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	// Create Bedrock client
	svc := bedrockruntime.NewFromConfig(cfg)

	return &BedrockClient{
		client:         svc,
		config:         cfg,
		modelID:        modelID,
		maxInputLength: 500, // Set a sensible default max input length
	}, nil
}

// EnableDebugMode enables verbose logging for AWS Bedrock interactions
func (c *BedrockClient) EnableDebugMode() {
	// AWS SDK v2 logging configuration
	c.config.ClientLogMode = aws.LogRequestWithBody | aws.LogResponseWithBody | aws.LogRetries
	c.client = bedrockruntime.NewFromConfig(c.config)
	log.Println("🔍 AWS Bedrock DEBUG mode enabled - full request/response logging activated")
}

// GenerateEmbedding generates an embedding for the given text using AWS Bedrock
func (c *BedrockClient) GenerateEmbedding(text string) ([]float32, error) {
	// Log the start of embedding generation with text length
	log.Printf("📤 Bedrock embedding request: %d chars, model: %s", len(text), c.modelID)

	// Validate input
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("empty text provided for embedding")
	}

	// Truncate text if needed
	if len(text) > c.maxInputLength {
		log.Printf("⚠️ Text exceeds maximum length (%d > %d), truncating", len(text), c.maxInputLength)
		text = text[:c.maxInputLength]
	}

	// Create properly formatted JSON request
	requestBody, err := json.Marshal(map[string]string{
		"inputText": text,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create embedding request
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.modelID),
		Body:        requestBody,
		ContentType: aws.String("application/json"),
	}

	// Invoke the model
	resp, err := c.client.InvokeModel(context.Background(), input)
	if err != nil {
		// AWS SDK v2 error handling - check for common error patterns
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "ThrottlingException"):
			log.Printf("❌ Bedrock throttling error: %v", err)
			return nil, fmt.Errorf("bedrock throttling error: %w", err)
		case strings.Contains(errMsg, "AccessDeniedException"):
			log.Printf("❌ Access denied to Bedrock: %v", err)
			return nil, fmt.Errorf("access denied to bedrock: %w", err)
		case strings.Contains(errMsg, "ModelNotFound"):
			log.Printf("❌ Bedrock model not found: %v", err)
			return nil, fmt.Errorf("bedrock model not found: %w", err)
		case strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "Timeout"):
			log.Printf("❌ Bedrock request timeout: %v", err)
			return nil, fmt.Errorf("bedrock request timeout: %w", err)
		default:
			log.Printf("❌ Bedrock request failed: %v", err)
			return nil, fmt.Errorf("bedrock embedding request failed: %w", err)
		}
	}

	// Parse the response
	var embeddingResp EmbeddingResponse
	if err := json.Unmarshal(resp.Body, &embeddingResp); err != nil {
		log.Printf("❌ Failed to parse embedding response: %v", err)
		return nil, fmt.Errorf("failed to parse embedding response: %w", err)
	}

	if len(embeddingResp.Embedding) == 0 {
		return nil, errors.New("empty embedding returned from Bedrock")
	}

	// Convert float64 to float32 if needed
	embedding := make([]float32, len(embeddingResp.Embedding))
	for i, v := range embeddingResp.Embedding {
		embedding[i] = float32(v)
	}

	log.Printf("✅ Bedrock embedding generated: %d dimensions", len(embedding))
	return embedding, nil
}

// SetMaxInputLength allows configuring the maximum input length
func (c *BedrockClient) SetMaxInputLength(length int) {
	if length > 0 {
		c.maxInputLength = length
		log.Printf("📏 Max input length set to: %d", length)
	}
}

// GetModelInfo returns information about the current model
func (c *BedrockClient) GetModelInfo() string {
	return fmt.Sprintf("Model ID: %s, Max Input Length: %d", c.modelID, c.maxInputLength)
}
