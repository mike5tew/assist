package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	// Parse command line flags
	text := flag.String("text", "Test embedding generation for troubleshooting", "Text to embed")
	model := flag.String("model", os.Getenv("AWS_BEDROCK_MODEL"), "Bedrock model ID")
	region := flag.String("region", os.Getenv("AWS_BEDROCK_REGION"), "AWS region")
	verbose := flag.Bool("verbose", true, "Enable verbose logging")
	flag.Parse()

	fmt.Println("🔍 AWS Bedrock Diagnostic Tool")
	fmt.Println("===============================")

	// Print environment information
	fmt.Println("\n📊 Environment Information:")
	fmt.Printf("AWS Region: %s\n", *region)
	fmt.Printf("AWS Bedrock Model: %s\n", *model)
	fmt.Printf("AWS_ACCESS_KEY_ID: %s\n", maskString(os.Getenv("AWS_ACCESS_KEY_ID")))
	fmt.Printf("AWS_SECRET_ACCESS_KEY: %s\n", maskString(os.Getenv("AWS_SECRET_ACCESS_KEY")))
	sess := os.Getenv("AWS_SESSION_TOKEN")
	if sess != "" {
		fmt.Printf("AWS_SESSION_TOKEN: %s\n", maskString(sess))
	} else {
		fmt.Println("AWS_SESSION_TOKEN: (not set)")
	}
	fmt.Printf("Text length: %d characters\n", len(*text))

	// Normalize region preference
	if os.Getenv("AWS_BEDROCK_REGION") == "" && os.Getenv("AWS_REGION") != "" {
		_ = os.Setenv("AWS_BEDROCK_REGION", os.Getenv("AWS_REGION"))
	}

	// Guard: ignore placeholder/invalid session tokens
	if tok := os.Getenv("AWS_SESSION_TOKEN"); tok != "" {
		if strings.Contains(strings.ToLower(tok), "your") || strings.Contains(strings.ToLower(tok), "here") {
			fmt.Println("⚠️ Detected placeholder AWS_SESSION_TOKEN; ignoring it for this run")
			_ = os.Unsetenv("AWS_SESSION_TOKEN")
		}
	}

	// Create AWS config with debug logging
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(*region),
	)
	if err != nil {
		log.Fatalf("❌ Failed to load AWS config: %v", err)
	}

	// Enable verbose logging
	if *verbose {
		cfg.ClientLogMode = aws.LogRequestWithBody | aws.LogResponseWithBody
	}

	// Print AWS configuration
	fmt.Println("\n📊 AWS SDK Configuration:")
	fmt.Printf("Client Log Mode: %v\n", cfg.ClientLogMode)
	fmt.Printf("Credentials Provider: %T\n", cfg.Credentials)

	// Test credentials
	fmt.Println("\n🔑 Testing AWS credentials...")
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		log.Fatalf("❌ Failed to retrieve credentials: %v", err)
	}
	fmt.Printf("✅ Retrieved credentials for: %s\n", creds.AccessKeyID)
	fmt.Printf("Expiration: %v\n", creds.Expires)

	// Create Bedrock client
	client := bedrockruntime.NewFromConfig(cfg)

	// Test Bedrock connection
	fmt.Println("\n📡 Testing Bedrock connection...")

	// Create input for the model
	inputText := *text
	inputStruct := struct {
		InputText string `json:"inputText"`
	}{
		InputText: inputText,
	}

	inputBytes, err := json.Marshal(inputStruct)
	if err != nil {
		log.Fatalf("❌ Failed to marshal input: %v", err)
	}

	// Create the request
	fmt.Println("\n📤 Sending embedding request...")
	startTime := time.Now()

	resp, err := client.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     model,
		Body:        inputBytes,
		ContentType: aws.String("application/json"),
	})

	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Bedrock request failed in %v: %v\n", duration, err)
		fmt.Println("\n🔍 Troubleshooting Tips:")
		fmt.Println("1. Check AWS credentials are correct and have Bedrock permissions")
		fmt.Println("2. Verify the model ID is correct and available in your region")
		fmt.Println("3. Confirm your AWS account has access to the model")
		fmt.Println("4. Check for network connectivity issues")
		os.Exit(1)
	}

	// Parse the response
	fmt.Printf("✅ Bedrock request successful in %v\n", duration)
	fmt.Printf("Response Content Type: %s\n", *resp.ContentType)
	fmt.Printf("Response Body Length: %d bytes\n", len(resp.Body))

	// Parse the embedding
	var response struct {
		Embedding []float32 `json:"embedding"`
	}

	if err := json.Unmarshal(resp.Body, &response); err != nil {
		fmt.Printf("❌ Failed to parse response: %v\n", err)
		fmt.Printf("Raw response: %s\n", string(resp.Body))
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully extracted embedding vector of dimension %d\n", len(response.Embedding))
	fmt.Println("\n🔍 Vector Sample (first 5 dimensions):")
	for i := 0; i < min(5, len(response.Embedding)); i++ {
		fmt.Printf("[%d]: %f\n", i, response.Embedding[i])
	}

	fmt.Println("\n✅ Bedrock diagnostics completed successfully")
}

// Helper function to mask credentials
func maskString(s string) string {
	if len(s) <= 8 {
		return "********"
	}
	return s[:4] + "..." + s[len(s)-4:]
}

// Helper for min of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
