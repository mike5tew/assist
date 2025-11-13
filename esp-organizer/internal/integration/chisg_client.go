package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// KnowledgeQuery represents a query to the CHISG service
type KnowledgeQuery struct {
	Topic          string   `json:"topic"`
	Concepts       []string `json:"concepts"`
	TargetAudience string   `json:"target_audience"`
}

// KnowledgeResponse represents the response from CHISG
type KnowledgeResponse struct {
	Summary         string            `json:"summary"`
	KeyConcepts     []string          `json:"key_concepts"`
	Prerequisites   []string          `json:"prerequisites"`
	Analogies       []string          `json:"analogies"`
	ConfidenceScore float64           `json:"confidence_score"`
	RelatedTopics   map[string]string `json:"related_topics"`
}

// CHISGClient handles communication with the CHISG service
type CHISGClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewCHISGClient creates a new CHISG client
func NewCHISGClient(baseURL string) *CHISGClient {
	if baseURL == "" {
		baseURL = "http://localhost:8080" // Default fallback
	}

	return &CHISGClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Query sends a knowledge query to CHISG and returns the response
func (c *CHISGClient) Query(ctx context.Context, query *KnowledgeQuery) (*KnowledgeResponse, error) {
	// Marshal the query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Create the HTTP request
	endpoint := fmt.Sprintf("%s/v1/knowledge/query", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(queryJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to CHISG: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CHISG returned status %d: %s", resp.StatusCode, string(body))
	}

	// Decode the response
	var response KnowledgeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode CHISG response: %w", err)
	}

	return &response, nil
}

// HealthCheck verifies that CHISG service is available
func (c *CHISGClient) HealthCheck(ctx context.Context) error {
	endpoint := fmt.Sprintf("%s/health", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("CHISG service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CHISG health check failed with status %d", resp.StatusCode)
	}

	return nil
}
