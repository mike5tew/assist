package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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

// InvokeClaude sends a prompt to Claude via Bedrock and returns the text response.
// Uses the Anthropic Messages API format.
func (c *BedrockLlamaClient) InvokeClaude(systemPrompt, userMessage string) (string, error) {
	modelID := os.Getenv("AWS_BEDROCK_CLAUDE_MODEL")
	if modelID == "" {
		modelID = "anthropic.claude-haiku-4-5-20251001-v1:0"
	}

	body := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        2048,
		"system":            systemPrompt,
		"messages": []map[string]interface{}{
			{"role": "user", "content": userMessage},
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Claude request: %w", err)
	}

	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        bodyBytes,
	}

	output, err := c.bedrockClient.InvokeModel(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("Claude invocation error: %w", err)
	}

	var resp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(output.Body, &resp); err != nil {
		return "", fmt.Errorf("failed to parse Claude response: %w", err)
	}
	for _, block := range resp.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}
	return "", fmt.Errorf("no text content in Claude response")
}
