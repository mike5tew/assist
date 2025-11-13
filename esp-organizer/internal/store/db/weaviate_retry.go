package db

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// WeaviateIsReady checks if Weaviate is fully ready for schema and data operations
func WeaviateIsReady(ctx context.Context) bool {
	if weaviateClient != nil {
		ready, err := weaviateClient.Misc().ReadyChecker().Do(ctx)
		if err == nil && ready {
			// Double-check schema access
			if _, err := weaviateClient.Schema().Getter().Do(ctx); err == nil {
				return true
			}
		}
	}

	// HTTP fallback to /v1/.well-known/ready
	host := GetEnvWithDefault("WEAVIATE_HOST", "localhost:8081")
	scheme := GetEnvWithDefault("WEAVIATE_SCHEME", "http")
	url := fmt.Sprintf("%s://%s/v1/.well-known/ready", scheme, host)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// WaitForWeaviateReady waits until Weaviate is ready with timeout (exponential backoff)
func WaitForWeaviateReady(ctx context.Context, maxAttempts int, delay time.Duration) error {
	next := delay
	for i := 0; i < maxAttempts; i++ {
		if WeaviateIsReady(ctx) {
			log.Printf("✅ Weaviate cluster is ready (attempt %d/%d)", i+1, maxAttempts)
			return nil
		}
		log.Printf("⏳ Waiting for Weaviate to be ready (attempt %d/%d)...", i+1, maxAttempts)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(next):
			// backoff but cap to 15s
			if next < 15*time.Second {
				next *= 2
				if next > 15*time.Second {
					next = 15 * time.Second
				}
			}
		}
	}
	return fmt.Errorf("weaviate not ready after %d attempts", maxAttempts)
}

// EnsureWeaviateSchema makes sure the schema exists and is properly initialized
func EnsureWeaviateSchema(ctx context.Context, className string, maxAttempts int) error {
	if weaviateClient == nil {
		return fmt.Errorf("weaviate client not initialized")
	}

	// First check if class exists
	var exists bool
	for i := 0; i < maxAttempts; i++ {
		schema, err := weaviateClient.Schema().Getter().Do(ctx)
		if err != nil {
			log.Printf("⚠️ Failed to get schema (attempt %d/%d): %v", i+1, maxAttempts, err)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		exists = false
		for _, class := range schema.Classes {
			if class.Class == className {
				exists = true
				break
			}
		}

		if exists {
			log.Printf("✅ Weaviate class %s exists", className)
			return nil
		}

		// Create it if it doesn't exist
		log.Printf("Class '%s' does not exist. Attempting to create it.", className)
		switch className {
		case "SemanticLinks":
			err = CreateSemanticLinksClass(ctx)
		case "MedicalExcerpt":
			err = CreateMedicalExcerptClass(ctx)
		// "EducationalSkills" is deprecated and has no creation function.
		default:
			err = fmt.Errorf("no creation logic defined for class '%s'", className)
		}

		if err == nil {
			log.Printf("✅ Created Weaviate class %s", className)
			return nil
		}

		log.Printf("⚠️ Failed to create class %s (attempt %d/%d): %v",
			className, i+1, maxAttempts, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return fmt.Errorf("failed to ensure schema for class %s after %d attempts",
		className, maxAttempts)
}

// WithRetry runs a function with retry logic specifically for Weaviate operations
func WithRetry(op string, fn func() error, maxRetries int, initialDelay time.Duration) error {
	var err error
	delay := initialDelay

	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		errStr := err.Error()
		// Check if it's a "leader not found" or consistency-related error
		if i < maxRetries-1 && (strings.Contains(errStr, "leader not found") ||
			strings.Contains(errStr, "consistency") ||
			strings.Contains(errStr, "timeout") ||
			strings.Contains(errStr, "unavailable")) {

			log.Printf("⚠️ Weaviate %s failed (attempt %d/%d): %v. Retrying in %v...",
				op, i+1, maxRetries, err, delay)
			time.Sleep(delay)
			delay = delay * 2 // Exponential backoff
			continue
		}

		// If it's a different error, return immediately
		return err
	}

	return fmt.Errorf("operation '%s' failed after %d attempts: %w", op, maxRetries, err)
}
