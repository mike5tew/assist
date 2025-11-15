package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// KnowledgeQuery matches the CHISG API contract request format
type KnowledgeQuery struct {
	Topic          string   `json:"topic"`
	Concepts       []string `json:"concepts"`
	TargetAudience string   `json:"target_audience"` // e.g., "gcse_student", "medical_student"
}

// HumanOSResponse wraps CHISG knowledge with age-appropriate adjustments
type HumanOSResponse struct {
	OriginalKnowledge *KnowledgeResponse `json:"knowledge"`
	AdjustedForAge    string             `json:"adjusted_summary"`
	DevelopmentStage  string             `json:"development_stage"`
	BarriersDetected  []string           `json:"barriers_detected"`
	Interventions     []string           `json:"recommended_interventions"`
}

// KnowledgeResponse matches the CHISG API contract response format
type KnowledgeResponse struct {
	Summary         string            `json:"summary"`
	KeyConcepts     []string          `json:"key_concepts"`
	Prerequisites   []string          `json:"prerequisites"`
	Analogies       []string          `json:"analogies"`
	ConfidenceScore float64           `json:"confidence_score"`
	RelatedTopics   map[string]string `json:"related_topics"`
	Timestamp       time.Time         `json:"timestamp"`
}

// cacheEntry stores a response with its creation time for TTL management
type cacheEntry struct {
	response  *KnowledgeResponse
	createdAt time.Time
}

// CHISGClient provides a Go interface to the CHISG microservice
// Features:
// - Connection pooling with configurable timeouts
// - Response caching (in-memory) for performance
// - Retry logic with exponential backoff
// - Request/response logging for debugging
// - Thread-safe operations with RWMutex
type CHISGClient struct {
	baseURL        string
	apiKey         string
	httpClient     *http.Client
	cache          map[string]*cacheEntry
	cacheMutex     sync.RWMutex
	cacheTTL       time.Duration
	maxRetries     int
	initialBackoff time.Duration
}

// NewCHISGClient creates a new CHISG client with API key from environment
func NewCHISGClient() (*CHISGClient, error) {
	baseURL := os.Getenv("CHISG_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080" // Default for local dev
	}

	apiKey := os.Getenv("CHISG_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("CHISG_API_KEY environment variable not set")
	}

	client := &CHISGClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		cache:          make(map[string]*cacheEntry),
		cacheTTL:       5 * time.Minute, // Cache responses for 5 minutes
		maxRetries:     3,
		initialBackoff: 100 * time.Millisecond,
	}

	return client, nil
}

// Query sends a knowledge query to CHISG and returns a response
// Returns error if the query fails after retries
func (c *CHISGClient) Query(ctx context.Context, query *KnowledgeQuery) (*KnowledgeResponse, error) {
	if query == nil {
		return nil, fmt.Errorf("query cannot be nil")
	}

	// Generate cache key from query
	cacheKey := c.generateCacheKey(query)

	// Check cache first
	if cached := c.getFromCache(cacheKey); cached != nil {
		log.Printf("[CHISG] Cache hit for query: %s", query.Topic)
		return cached, nil
	}

	log.Printf("[CHISG] Querying CHISG for topic: %s, audience: %s", query.Topic, query.TargetAudience)

	// Retry logic with exponential backoff
	var lastErr error
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			backoff := c.initialBackoff * time.Duration(1<<uint(attempt-1)) // exponential: 100ms, 200ms, 400ms
			log.Printf("[CHISG] Retry attempt %d/%d, waiting %v", attempt+1, c.maxRetries, backoff)

			select {
			case <-time.After(backoff):
				// Continue to retry
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		response, err := c.queryWithTimeout(ctx, query)
		if err == nil {
			// Cache the successful response
			c.putInCache(cacheKey, response)
			log.Printf("[CHISG] Query successful, cached result for %s", query.Topic)
			return response, nil
		}

		lastErr = err
		log.Printf("[CHISG] Query attempt %d failed: %v", attempt+1, err)
	}

	return nil, fmt.Errorf("failed to query CHISG after %d attempts: %w", c.maxRetries, lastErr)
}

// queryWithTimeout executes a single query attempt with context-aware timeout
func (c *CHISGClient) queryWithTimeout(ctx context.Context, query *KnowledgeQuery) (*KnowledgeResponse, error) {
	// Create request body
	requestBody, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Create HTTP request
	endpoint := fmt.Sprintf("%s/v1/knowledge/query", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "HumanOS-CHISGClient/1.0")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Execute request
	start := time.Now()
	resp, err := c.httpClient.Do(req)
	latency := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Log latency for performance monitoring
	if latency > 200*time.Millisecond {
		log.Printf("[CHISG] ⚠️  High latency detected: %v (target: <200ms)", latency)
	} else {
		log.Printf("[CHISG] ✅ Request completed in %v", latency)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CHISG returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response KnowledgeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	response.Timestamp = time.Now()
	return &response, nil
}

// generateCacheKey creates a unique cache key for a query
func (c *CHISGClient) generateCacheKey(query *KnowledgeQuery) string {
	// Simple key generation - combines topic, audience, and concepts
	return fmt.Sprintf("%s:%s:%s", query.Topic, query.TargetAudience, joinStrings(query.Concepts))
}

// joinStrings joins slice elements with comma (helper for cache key)
func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += ","
		}
		result += s
	}
	return result
}

// getFromCache retrieves a cached response if it exists and hasn't expired
func (c *CHISGClient) getFromCache(key string) *KnowledgeResponse {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil
	}

	// Check if cache entry has expired
	if time.Since(entry.createdAt) > c.cacheTTL {
		log.Printf("[CHISG] Cache expired for key: %s", key)
		return nil
	}

	return entry.response
}

// putInCache stores a response in the cache with timestamp
func (c *CHISGClient) putInCache(key string, response *KnowledgeResponse) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	c.cache[key] = &cacheEntry{
		response:  response,
		createdAt: time.Now(),
	}
	log.Printf("[CHISG] Cached response for key: %s (TTL: %v)", key, c.cacheTTL)
}

// ClearCache removes all cached entries
func (c *CHISGClient) ClearCache() {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	c.cache = make(map[string]*cacheEntry)
	log.Printf("[CHISG] Cache cleared")
}

// SetCacheTTL updates the cache time-to-live duration
func (c *CHISGClient) SetCacheTTL(ttl time.Duration) {
	c.cacheTTL = ttl
	log.Printf("[CHISG] Cache TTL updated to %v", ttl)
}

// GetStats returns cache statistics for monitoring
func (c *CHISGClient) GetStats() map[string]interface{} {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	// Count expired entries
	expiredCount := 0
	for _, entry := range c.cache {
		if time.Since(entry.createdAt) > c.cacheTTL {
			expiredCount++
		}
	}

	return map[string]interface{}{
		"cached_items":  len(c.cache),
		"expired_items": expiredCount,
		"cache_ttl":     c.cacheTTL.String(),
		"max_retries":   c.maxRetries,
		"base_url":      c.baseURL,
		"http_timeout":  "10s",
		"backoff_start": c.initialBackoff.String(),
	}
}

// SetBaseURL allows runtime configuration of the CHISG service URL
func (c *CHISGClient) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
	log.Printf("[CHISG] Base URL updated to: %s", baseURL)
}

// SetMaxRetries configures the number of retry attempts
func (c *CHISGClient) SetMaxRetries(retries int) {
	if retries > 0 {
		c.maxRetries = retries
		log.Printf("[CHISG] Max retries updated to: %d", retries)
	}
}

// SetInitialBackoff configures the initial backoff duration
func (c *CHISGClient) SetInitialBackoff(backoff time.Duration) {
	if backoff > 0 {
		c.initialBackoff = backoff
		log.Printf("[CHISG] Initial backoff updated to: %v", backoff)
	}
}

// Health checks if the CHISG service is reachable
func (c *CHISGClient) Health(ctx context.Context) error {
	healthURL := fmt.Sprintf("%s/health", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	log.Printf("[CHISG] ✅ Health check passed")
	return nil
}
