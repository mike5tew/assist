package main

import (
	"context"
	"esp-organizer/internal/InfoFlow/InfoOut/llm"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Command-line flags
	fixFlag := flag.Bool("fix", false, "Automatically fix dimension mismatch by resetting problematic classes")
	dimensionFlag := flag.Int("dimension", 0, "Expected vector dimension (0 = auto-detect)")
	classFlag := flag.String("class", "", "Specific class to check/fix (default: all classes)")
	flag.Parse()

	// Connect to Weaviate
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to connect to Weaviate: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Determine which classes to check
	var classesToCheck []string
	if *classFlag != "" {
		classesToCheck = strings.Split(*classFlag, ",")
	} else {
		classesToCheck = []string{
			"EducationalSkills",
			"SubjectAreaContent",
			"SemanticLinks",
			"MedicalExcerpt",
		}
	}

	fmt.Println("=== SEMANTIC PROCESSING PIPELINE MONITOR ===")

	// Check each class
	problemFound := false
	for _, className := range classesToCheck {
		className = strings.TrimSpace(className)
		if className == "" {
			continue
		}

		// Check if class exists
		exists, err := db.WeaviateCollectionExists(ctx, className)
		if err != nil {
			log.Printf("Error checking class %s: %v", className, err)
			continue
		}

		if !exists {
			fmt.Printf("✅ Class %s does not exist yet (will be created when needed)\n", className)
			continue
		}

		// Attempt a simple query to check for dimension issues
		expectedDimension := *dimensionFlag
		if expectedDimension == 0 {
			// Try to detect from AWS Bedrock client
			expectedDimension = detectBedrockDimension()
			if expectedDimension == 0 {
				// Fall back to current popular Titan model dimension
				expectedDimension = 384
			}
		}

		vectorDimensions, err := checkVectorDimensions(ctx, className)
		if err != nil {
			// Handle specific dimension mismatch error
			if strings.Contains(err.Error(), "vector with length") &&
				strings.Contains(err.Error(), "Existing nodes have vectors with length") {
				// Extract dimensions from the error message
				existingDim := extractDimensionFromError(err.Error())
				if existingDim > 0 && existingDim != expectedDimension {
					problemFound = true
					fmt.Printf("❌ Class %s has vector dimension mismatch: existing=%d, expected=%d\n",
						className, existingDim, expectedDimension)

					if *fixFlag {
						fmt.Printf("🔄 Resetting class %s to fix dimension mismatch...\n", className)
						if err := db.ResetWeaviateClass(ctx, className); err != nil {
							fmt.Printf("❌ Failed to reset class %s: %v\n", className, err)
						} else {
							fmt.Printf("✅ Successfully reset class %s\n", className)
						}
					} else {
						fmt.Printf("ℹ️ Run with -fix flag to automatically reset this class\n")
					}
				}
			} else {
				fmt.Printf("⚠️ Class %s error: %v\n", className, err)
			}
		} else if vectorDimensions > 0 {
			if vectorDimensions != expectedDimension {
				problemFound = true
				fmt.Printf("❌ Class %s has vector dimension mismatch: current=%d, expected=%d\n",
					className, vectorDimensions, expectedDimension)

				if *fixFlag {
					fmt.Printf("🔄 Resetting class %s to fix dimension mismatch...\n", className)
					if err := db.ResetWeaviateClass(ctx, className); err != nil {
						fmt.Printf("❌ Failed to reset class %s: %v\n", className, err)
					} else {
						fmt.Printf("✅ Successfully reset class %s\n", className)
					}
				} else {
					fmt.Printf("ℹ️ Run with -fix flag to automatically reset this class\n")
				}
			} else {
				fmt.Printf("✅ Class %s has correct vector dimensions: %d\n", className, vectorDimensions)
			}
		} else {
			fmt.Printf("✅ Class %s is empty, unable to check vector dimensions\n", className)
		}

		// Get additional class info
		printClassInfo(ctx, className)
	}

	if !problemFound {
		fmt.Println("\n✅ No vector dimension mismatches detected")
	} else if !*fixFlag {
		fmt.Println("\nTo fix the detected issues, run:")
		fmt.Printf("  %s -fix\n", os.Args[0])
	}
}

// extractDimensionFromError parses error message to extract dimension value
func extractDimensionFromError(errMsg string) int {
	// Look for pattern like: "Existing nodes have vectors with length 1024"
	var dimension int
	_, err := fmt.Sscanf(errMsg, "%*s %*s %*s %*s %*s %*s %*s %d", &dimension)
	if err != nil || dimension == 0 {
		// Try alternative pattern
		_, err = fmt.Sscanf(errMsg, "%*s %*s %*s %*s %*s %*s %d", &dimension)
	}
	return dimension
}

// detectBedrockDimension attempts to detect the current AWS Bedrock Titan model's dimension
func detectBedrockDimension() int {
	// Create a temporary Llama client
	client := llm.NewLlamaClient()
	if client == nil {
		return 0
	}

	// Generate a test embedding to check dimensions
	vector, err := client.GenerateEmbedding("Test embedding for dimension detection")
	if err != nil {
		log.Printf("Warning: Could not detect AWS Bedrock dimensions: %v", err)
		return 0
	}

	return len(vector)
}

// checkVectorDimensions attempts a simple operation on a class to detect vector dimension issues
func checkVectorDimensions(ctx context.Context, className string) (int, error) {
	// Get the client
	client := db.GetWeaviateClient()
	if client == nil {
		return 0, fmt.Errorf("Weaviate client not initialized")
	}

	// Try to get an object to check its vector dimensions
	result, err := client.Data().ObjectsGetter().
		WithClassName(className).
		WithLimit(1).
		WithVector().
		Do(ctx)
	if err != nil {
		return 0, err
	}

	if len(result) == 0 {
		// It's not an error if a class is empty, we just can't check dimensions.
		return 0, nil
	}

	// Get the vector from the first object
	vector := result[0].Vector
	if vector == nil {
		return 0, fmt.Errorf("object has no vector")
	}

	return len(vector), nil
}

// printClassInfo gets and displays detailed information about a Weaviate class
func printClassInfo(ctx context.Context, className string) {
	client := db.GetWeaviateClient()
	if client == nil {
		fmt.Printf("❌ %s: Weaviate client not initialized\n", className)
		return
	}

	// Get class info
	schema, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
	if err != nil {
		fmt.Printf("❌ %s: Error getting schema - %v\n", className, err)
		return
	}

	// Get count
	result, err := client.GraphQL().Aggregate().
		WithClassName(className).
		WithFields(graphql.Field{Name: "meta", Fields: []graphql.Field{{Name: "count"}}}).
		Do(ctx)

	var count int
	if err == nil {
		data, ok := result.Data["Aggregate"].(map[string]interface{})
		if ok {
			classData, ok := data[className].([]interface{})
			if ok && len(classData) > 0 {
				metaData, ok := classData[0].(map[string]interface{})
				if ok {
					countVal, ok := metaData["meta"].(map[string]interface{})
					if ok {
						count, _ = countVal["count"].(int)
					}
				}
			}
		}
	}

	fmt.Printf("✅ %s: %d objects, vectorizer=%s\n",
		className, count, schema.Vectorizer)
}
