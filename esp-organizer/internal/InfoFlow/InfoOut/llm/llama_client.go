package llm

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/go-resty/resty/v2"
)

// LlamaClient is a client for interacting with various LLM providers
type LlamaClient struct {
	// Common fields
	mutex              sync.Mutex
	fallbackVectorSize int  // Size for fallback vectors
	useDummy           bool // Whether to use dummy embeddings

	// AWS Bedrock specific
	bedrockClient     *bedrockruntime.Client
	embeddingModelID  string
	completionModelID string

	// OpenAI/Ollama specific
	httpClient *HTTPClient
	apiKey     string
	baseURL    string
	model      string

	// Config flags
	provider string // "bedrock", "openai", "ollama", "llama" or "dummy"
}

// NewLlamaClient creates a new client for interacting with LLM providers
func NewLlamaClient() *LlamaClient {
	// Start with a base client with reasonable defaults
	client := &LlamaClient{

		fallbackVectorSize: 384,
		useDummy:           false,
		httpClient:         NewHTTPClient(),
	}

	// First try AWS Bedrock
	if bedrockClient := initializeBedrockClient(); bedrockClient != nil {
		client.bedrockClient = bedrockClient
		client.provider = "bedrock"

		// Get model IDs from environment variables with defaults
		client.embeddingModelID = os.Getenv("AWS_BEDROCK_EMBEDDING_MODEL_ID")
		if client.embeddingModelID == "" {
			// Update to a valid model ID - titan-embed-text-v1 is the correct format
			client.embeddingModelID = "amazon.titan-embed-text-v2:0" // Default embedding model
			log.Printf("AWS_BEDROCK_EMBEDDING_MODEL_ID not set, using default: %s", client.embeddingModelID)
		}

		client.completionModelID = os.Getenv("AWS_BEDROCK_COMPLETION_MODEL_ID")
		if client.completionModelID == "" {
			// Update to a valid model ID - claude-v2 is the correct format
			client.completionModelID = os.Getenv("AWS_BEDROCK_GENERATIVE_MODEL_ID") // Default completion model
			log.Printf("AWS_BEDROCK_COMPLETION_MODEL_ID not set, using default: %s", client.completionModelID)
		}

		// Add validation to check model IDs
		if !isValidBedrockModelID(client.embeddingModelID) {
			log.Printf("Warning: Invalid embedding model ID: %s - using fallback embeddings", client.embeddingModelID)
			client.useDummy = true
		}

		if !isValidBedrockModelID(client.completionModelID) {
			log.Printf("Warning: Invalid completion model ID: %s - using fallback text generation", client.completionModelID)
		}

		log.Printf("Using AWS Bedrock for LLM with models: %s (embeddings), %s (completions)",
			client.embeddingModelID, client.completionModelID)
		return client
	}

	// Then try OpenAI
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey != "" {
		client.apiKey = openaiKey
		client.baseURL = "https://api.openai.com/v1"
		client.provider = "openai"
		log.Println("✅ OpenAI API key found - using OpenAI for embeddings and completions")
		return client
	}

	// Check for Ollama/Llama configuration
	llamaURL := os.Getenv("LLAMA_API_URL")
	llamaModel := os.Getenv("LLAMA_MODEL")
	ollamaURL := os.Getenv("OLLAMA_BASE_URL")

	if ollamaURL != "" {
		client.baseURL = ollamaURL
		client.model = os.Getenv("OLLAMA_MODEL")
		if client.model == "" {
			client.model = "nomic-embed-text"
		}
		client.provider = "ollama"
		log.Printf("✅ Ollama configured - using %s model for local embeddings", client.model)
		return client
	} else if llamaURL != "" {
		client.baseURL = llamaURL
		client.model = "nomic-embed-text"
		if llamaModel != "" {
			client.model = llamaModel
		}
		client.provider = "llama"
		log.Printf("✅ Llama API configured - using %s model", client.model)
		return client
	}

	// Fallback to dummy mode if no providers are available
	client.useDummy = true
	client.provider = "dummy"
	log.Println("⚠️ No LLM providers configured - using dummy embeddings")
	log.Println("🔧 For real embeddings, configure one of these in .env:")
	log.Println("   • AWS Bedrock: AWS_BEDROCK_REGION + AWS_BEDROCK_MODEL (recommended)")
	log.Println("   • OpenAI: OPENAI_API_KEY")
	log.Println("   • Local Ollama: OLLAMA_BASE_URL + OLLAMA_MODEL")

	return client
}

// Initialize AWS Bedrock client
func initializeBedrockClient() *bedrockruntime.Client {
	// Load AWS region from environment variable
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1" // Default region
		log.Printf("AWS_REGION not set, using default: %s", region)
	}

	// Initialize AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Printf("Error loading AWS configuration: %v", err)
		return nil
	}

	// Create the Bedrock runtime client
	return bedrockruntime.NewFromConfig(cfg)
}

// Generate generates text using the configured LLM provider
func (c *LlamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	switch c.provider {
	case "bedrock":
		return c.generateBedrockCompletion(ctx, prompt)
	case "openai":
		return c.generateOpenAICompletion(ctx, prompt)
	case "ollama":
		return c.generateOllamaCompletion(ctx, prompt)
	case "llama":
		return c.generateLlamaCompletion(ctx, prompt)
	default:
		return c.fallbackGenerate(prompt), nil
	}
}

// GenerateEmbedding generates an embedding vector for the given text
func (c *LlamaClient) GenerateEmbedding(text string) ([]float32, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	switch c.provider {
	case "bedrock":
		return c.generateBedrockEmbedding(text)
	case "openai":
		return c.generateOpenAIEmbedding(text)
	case "ollama":
		return c.generateOllamaEmbedding(text)
	case "llama":
		return c.generateLlamaEmbedding(text)
	default:
		return c.generateDummyEmbedding(text), nil
	}
}

// BatchGenerateEmbeddings generates embeddings for multiple texts
func (c *LlamaClient) BatchGenerateEmbeddings(texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))

	for i, text := range texts {
		embedding, err := c.GenerateEmbedding(text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text %d: %w", i, err)
		}
		embeddings[i] = embedding
	}

	return embeddings, nil
}

// AWS Bedrock implementations

// generateBedrockCompletion calls AWS Bedrock's InvokeModel API for text generation
// NOTE: Claude 3+ models require the Converse API, but this is a simplified version
// that falls back to the legacy InvokeModel format for compatibility
func (c *LlamaClient) generateBedrockCompletion(ctx context.Context, prompt string) (string, error) {
	if c.bedrockClient == nil {
		log.Printf("Error: Bedrock client not initialized")
		return c.fallbackGenerate(prompt), nil
	}

	if c.completionModelID == "" {
		log.Printf("Error: Completion model ID not set")
		return c.fallbackGenerate(prompt), nil
	}

	// Prepare the request based on the model
	var requestBody []byte
	var err error

	if strings.HasPrefix(c.completionModelID, "anthropic") {
		// Claude format - use the Messages API format for Claude 3+
		if strings.Contains(c.completionModelID, "claude-3") {
			// For Claude 3+, use the Messages API format
			requestBody, err = json.Marshal(map[string]interface{}{
				"anthropic_version": "bedrock-2023-05-31",
				"max_tokens":        2000,
				"messages": []map[string]interface{}{
					{
						"role": "user",
						"content": []map[string]interface{}{
							{
								"type": "text",
								"text": prompt,
							},
						},
					},
				},
			})
		} else {
			// Claude 2.x format (legacy)
			requestBody, err = json.Marshal(map[string]interface{}{
				"prompt":               fmt.Sprintf("\n\nHuman: %s\n\nAssistant:", prompt),
				"max_tokens_to_sample": 2000,
				"temperature":          0.7,
			})
		}
	} else if strings.HasPrefix(c.completionModelID, "amazon.titan") {
		// Titan model format
		requestBody, err = json.Marshal(map[string]interface{}{
			"inputText": prompt,
			"textGenerationConfig": map[string]interface{}{
				"maxTokenCount": 2000,
				"temperature":   0.7,
			},
		})
	} else {
		log.Printf("Unsupported model: %s, using default format", c.completionModelID)
		requestBody, err = json.Marshal(map[string]interface{}{
			"prompt":      prompt,
			"max_tokens":  2000,
			"temperature": 0.7,
		})
	}

	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return c.fallbackGenerate(prompt), nil
	}

	// Invoke the model
	output, err := c.bedrockClient.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.completionModelID),
		Body:        requestBody,
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		log.Printf("AWS Bedrock API error: %v", err)
		return c.fallbackGenerate(prompt), nil
	}

	var response map[string]interface{}
	if err := json.Unmarshal(output.Body, &response); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return c.fallbackGenerate(prompt), nil
	}

	// Extract the result based on the model
	var result string
	if strings.HasPrefix(c.completionModelID, "anthropic") {
		// Try Claude 3+ format first
		if content, ok := response["content"].([]interface{}); ok && len(content) > 0 {
			if textBlock, ok := content[0].(map[string]interface{}); ok {
				result, _ = textBlock["text"].(string)
			}
		} else {
			// Fallback to Claude 2.x format
			result, _ = response["completion"].(string)
		}
	} else if strings.HasPrefix(c.completionModelID, "amazon.titan") {
		if results, ok := response["results"].([]interface{}); ok && len(results) > 0 {
			if resultMap, ok := results[0].(map[string]interface{}); ok {
				result, _ = resultMap["outputText"].(string)
			}
		}
	} else {
		// Generic extraction attempt
		if text, ok := response["generated_text"].(string); ok {
			result = text
		} else if text, ok := response["text"].(string); ok {
			result = text
		} else {
			log.Printf("Could not extract text from response: %+v", response)
			result = c.fallbackGenerate(prompt)
		}
	}

	return result, nil
}

// generateBedrockEmbedding calls AWS Bedrock embedding API
func (c *LlamaClient) generateBedrockEmbedding(text string) ([]float32, error) {
	// Check if client is initialized
	if c.bedrockClient == nil {
		log.Printf("Warning: AWS Bedrock client not initialized, using fallback embedding")
		return c.generateDummyEmbedding(text), nil
	}

	// Check if model ID is set
	if c.embeddingModelID == "" {
		log.Printf("Error: Embedding model ID not set")
		return c.generateDummyEmbedding(text), nil
	}

	// Prepare the request based on the model
	var requestBody []byte
	var err error

	// Trim text to max length (8000 characters seems safe for most models)
	if len(text) > 8000 {
		text = text[:8000]
	}

	if strings.HasPrefix(c.embeddingModelID, "amazon.titan") {
		// Titan embedding format
		requestBody, err = json.Marshal(map[string]interface{}{
			"inputText": text,
		})
	} else if strings.HasPrefix(c.embeddingModelID, "cohere") {
		// Cohere embedding format
		requestBody, err = json.Marshal(map[string]interface{}{
			"texts":      []string{text},
			"input_type": "search_document",
		})
	} else {
		// Generic embedding format
		requestBody, err = json.Marshal(map[string]interface{}{
			"text": text,
		})
	}

	if err != nil {
		log.Printf("Error marshalling embedding request: %v", err)
		return c.generateDummyEmbedding(text), nil
	}

	// Call the Bedrock API with the specific model ID
	output, err := c.bedrockClient.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.embeddingModelID),
		Body:        requestBody,
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		log.Printf("AWS Bedrock API error: %v", err)
		return c.generateDummyEmbedding(text), nil
	}

	var response map[string]interface{}
	if err := json.Unmarshal(output.Body, &response); err != nil {
		log.Printf("Error unmarshalling embedding response: %v", err)
		return c.generateDummyEmbedding(text), nil
	}

	// Extract embeddings based on the model
	var embeddings []float32
	if strings.HasPrefix(c.embeddingModelID, "amazon.titan") {
		// Titan embedding format
		if embedding, ok := response["embedding"].([]interface{}); ok {
			embeddings = make([]float32, len(embedding))
			for i, val := range embedding {
				if floatVal, ok := val.(float64); ok {
					embeddings[i] = float32(floatVal)
				}
			}
		}
	} else if strings.HasPrefix(c.embeddingModelID, "cohere") {
		// Cohere embedding format
		if embeddingArr, ok := response["embeddings"].([]interface{}); ok && len(embeddingArr) > 0 {
			if embedding, ok := embeddingArr[0].([]interface{}); ok {
				embeddings = make([]float32, len(embedding))
				for i, val := range embedding {
					if floatVal, ok := val.(float64); ok {
						embeddings[i] = float32(floatVal)
					}
				}
			}
		}
	} else {
		// Generic embedding extraction attempt
		if embedding, ok := response["embedding"].([]interface{}); ok {
			embeddings = make([]float32, len(embedding))
			for i, val := range embedding {
				if floatVal, ok := val.(float64); ok {
					embeddings[i] = float32(floatVal)
				}
			}
		} else if embeddingVec, ok := response["vector"].([]interface{}); ok {
			embeddings = make([]float32, len(embeddingVec))
			for i, val := range embeddingVec {
				if floatVal, ok := val.(float64); ok {
					embeddings[i] = float32(floatVal)
				}
			}
		}
	}

	if len(embeddings) == 0 {
		log.Printf("Could not extract embeddings from response: %+v", response)
		return c.generateDummyEmbedding(text), nil
	}

	return embeddings, nil
}

// OpenAI implementations

// generateOpenAICompletion calls OpenAI's completion API
func (c *LlamaClient) generateOpenAICompletion(ctx context.Context, prompt string) (string, error) {
	type reqBody struct {
		Model       string  `json:"model"`
		Prompt      string  `json:"prompt"`
		MaxTokens   int     `json:"max_tokens"`
		Temperature float32 `json:"temperature"`
	}
	type choice struct {
		Text string `json:"text"`
	}
	type respBody struct {
		Choices []choice `json:"choices"`
	}

	body := reqBody{
		Model:       "text-davinci-003",
		Prompt:      prompt,
		MaxTokens:   150,
		Temperature: 0.7,
	}
	var resp respBody
	_, err := c.httpClient.client.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&resp).
		Post(c.baseURL + "/completions")
	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		return c.fallbackGenerate(prompt), nil
	}
	if len(resp.Choices) == 0 {
		log.Printf("OpenAI returned empty response")
		return c.fallbackGenerate(prompt), nil
	}
	return resp.Choices[0].Text, nil
}

// generateOpenAIEmbedding calls OpenAI's embedding API
func (c *LlamaClient) generateOpenAIEmbedding(text string) ([]float32, error) {
	type OpenAIRequest struct {
		Input string `json:"input"`
		Model string `json:"model"`
	}

	type OpenAIResponse struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	request := OpenAIRequest{
		Input: text,
		Model: c.model,
	}

	var response OpenAIResponse
	_, err := c.httpClient.client.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&response).
		Post(c.baseURL + "/embeddings")

	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		return c.generateDummyEmbedding(text), nil
	}

	if len(response.Data) == 0 {
		log.Printf("OpenAI returned no embeddings")
		return c.generateDummyEmbedding(text), nil
	}

	return response.Data[0].Embedding, nil
}

// Ollama implementations

// generateOllamaCompletion calls Ollama's generation API
func (c *LlamaClient) generateOllamaCompletion(ctx context.Context, prompt string) (string, error) {
	type OllamaRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}

	type OllamaResponse struct {
		Response string `json:"response"`
	}

	request := OllamaRequest{
		Model:  c.model,
		Prompt: prompt,
	}

	var response OllamaResponse
	_, err := c.httpClient.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&response).
		Post(c.baseURL + "/api/generate")

	if err != nil {
		log.Printf("Ollama API error: %v", err)
		return c.fallbackGenerate(prompt), nil
	}

	return response.Response, nil
}

// generateOllamaEmbedding calls Ollama API endpoint for embeddings
func (c *LlamaClient) generateOllamaEmbedding(text string) ([]float32, error) {
	type OllamaRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}

	type OllamaResponse struct {
		Embedding []float32 `json:"embedding"`
	}

	request := OllamaRequest{
		Model:  c.model,
		Prompt: text,
	}

	var response OllamaResponse
	_, err := c.httpClient.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&response).
		Post(c.baseURL + "/api/embeddings")

	if err != nil {
		log.Printf("Ollama API error: %v", err)
		return c.generateDummyEmbedding(text), nil
	}

	return response.Embedding, nil
}

// Llama implementations

// generateLlamaCompletion calls a custom Llama API
func (c *LlamaClient) generateLlamaCompletion(ctx context.Context, prompt string) (string, error) {
	type LlamaRequest struct {
		Text string `json:"text"`
	}

	type LlamaResponse struct {
		Completion string `json:"completion"`
	}

	request := LlamaRequest{Text: prompt}

	var response LlamaResponse
	_, err := c.httpClient.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&response).
		Post(c.baseURL + "/completion")

	if err != nil {
		log.Printf("Llama API error: %v", err)
		return c.fallbackGenerate(prompt), nil
	}

	return response.Completion, nil
}

// generateLlamaEmbedding calls custom Llama API for embeddings
func (c *LlamaClient) generateLlamaEmbedding(text string) ([]float32, error) {
	type LlamaRequest struct {
		Text string `json:"text"`
	}

	type LlamaResponse struct {
		Embedding []float32 `json:"embedding"`
	}

	request := LlamaRequest{Text: text}

	var response LlamaResponse
	_, err := c.httpClient.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&response).
		Post(c.baseURL + "/embeddings")

	if err != nil {
		log.Printf("Llama API error: %v", err)
		return c.generateDummyEmbedding(text), nil
	}

	return response.Embedding, nil
}

// Fallback implementations

// fallbackGenerate provides a simple response when no LLM provider is available
func (c *LlamaClient) fallbackGenerate(prompt string) string {
	log.Printf("Using fallback text generation for prompt: %.50s...", prompt)

	// Return a simple fallback response
	return fmt.Sprintf("I processed your request about '%s', but I'm currently in fallback mode with limited capabilities. Please ensure an LLM provider is properly configured.", truncateString(prompt, 50))
}

// generateDummyEmbedding creates a deterministic dummy embedding for development
func (c *LlamaClient) generateDummyEmbedding(text string) []float32 {
	if !c.useDummy {
		// Only log when we're not explicitly in dummy mode
		log.Printf("Using fallback embedding generation for: %.50s...", text)
	}

	// First try the hash-based approach for better semantic matching
	if c.useDummy {
		// Create a deterministic hash-based embedding
		hash := md5.Sum([]byte(strings.ToLower(text)))

		// Convert hash to vector (use consistent size)
		embedding := make([]float32, c.fallbackVectorSize)

		// Use hash bytes to seed the embedding
		for i := 0; i < c.fallbackVectorSize; i++ {
			// Use hash bytes cyclically and add some variation
			byteIndex := i % 16
			seed := int(hash[byteIndex]) + i

			// Generate deterministic but varied values
			value := float32(math.Sin(float64(seed)*0.1)) * 0.5

			// Add text-specific patterns for better similarity
			if strings.Contains(strings.ToLower(text), "memory") && i < 50 {
				value += 0.3
			}
			if strings.Contains(strings.ToLower(text), "reading") && i >= 50 && i < 100 {
				value += 0.3
			}
			if strings.Contains(strings.ToLower(text), "motor") && i >= 100 && i < 150 {
				value += 0.3
			}
			if strings.Contains(strings.ToLower(text), "emotional") && i >= 150 && i < 200 {
				value += 0.3
			}

			embedding[i] = value
		}

		// Normalize the vector
		magnitude := float32(0)
		for _, v := range embedding {
			magnitude += v * v
		}
		magnitude = float32(math.Sqrt(float64(magnitude)))

		if magnitude > 0 {
			for i := range embedding {
				embedding[i] /= magnitude
			}
		}

		return embedding
	}

	// Alternative approach using FNV hash for more randomness
	vector := make([]float32, c.fallbackVectorSize)

	// Create a deterministic hash from the text
	h := fnv.New32a()
	h.Write([]byte(text))
	seed := h.Sum32()

	// Generate vector values with this seed
	r := rand.New(rand.NewSource(int64(seed)))
	for i := range vector {
		vector[i] = (r.Float32() * 2) - 1 // Values between -1 and 1
	}

	// Normalize vector to unit length
	var sum float32
	for _, v := range vector {
		sum += v * v
	}
	length := float32(math.Sqrt(float64(sum)))
	for i := range vector {
		vector[i] = vector[i] / length
	}

	return vector
}

// Helper function to truncate strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// HTTPClient type with needed methods

// HTTPClient wraps an HTTP client
type HTTPClient struct {
	client *resty.Client
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: resty.New(),
	}
}

// Add this helper function to validate model IDs
func isValidBedrockModelID(modelID string) bool {
	// Basic validation - ensure model ID has a namespace (contains a dot)
	// and isn't empty
	return modelID != "" && strings.Contains(modelID, ".")
}
