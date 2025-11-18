package main

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/domain/infoin"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Parse command line flags
	query := flag.String("query", "", "The query to search for")
	domain := flag.String("domain", "immunology", "The domain to search in")
	maxResults := flag.Int("max", 5, "Maximum number of results to return")
	verbose := flag.Bool("verbose", false, "Show verbose output including context chunks")
	setupSchema := flag.Bool("setup", false, "Set up the HSG schema if it doesn't exist")
	flag.Parse()

	if *query == "" {
		log.Fatal("Query cannot be empty. Use --query to specify a search term.")
	}

	// Initialize HSG query service
	fmt.Println("🔍 Initializing HSG Query Service...")

	// Check if schema setup is requested - use the main setup command instead of our custom function
	if *setupSchema {
		fmt.Println("🏗️ Setting up HSG Schema first...")
		cmd := exec.Command("go", "run", "cmd/setup-hsg-schema/main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to set up HSG schema: %v", err)
		}
		fmt.Println("✅ HSG Schema setup complete")
	}

	hsgService, err := infoin.NewHSGQueryService()
	if err != nil {
		log.Fatalf("Failed to initialize HSG query service: %v", err)
	}
	fmt.Println("✅ HSG Query Service initialized")

	fmt.Printf("\n🔎 Executing HSG search:\n")
	fmt.Printf("  • Query: %s\n", *query)
	fmt.Printf("  • Domain: %s\n", *domain)
	fmt.Printf("  • Max Results: %d\n\n", *maxResults)

	startTime := time.Now()

	// Perform the search with proper error handling
	ragContext, answer, err := hsgService.QueryHSG(context.Background(), *query, *domain, *maxResults)
	if err != nil {
		log.Printf("HSG query failed: %v", err)
		fmt.Println("\n⚠️ SEARCH FAILED")
		fmt.Println("==================")
		fmt.Printf("Error: %v\n\n", err)

		if strings.Contains(err.Error(), "class does not exist") {
			fmt.Println("💡 The HSG schema needs to be set up. Run this command:")
			fmt.Println("    make setup-hsg-schema")
			fmt.Println("\nThen try your search again.")
		}

		os.Exit(1)
	}

	duration := time.Since(startTime)

	// Print the results
	fmt.Println("📋 SEARCH RESULTS")
	fmt.Println("==================")
	fmt.Printf("⏱️  Query completed in: %.2f seconds\n\n", duration.Seconds())

	fmt.Println("📝 ANSWER:")
	fmt.Println(answer)
	fmt.Println()

	// Handle the case where context may be nil
	if ragContext == nil {
		fmt.Println("⚠️ No context data was returned")
	} else {
		if *verbose {
			fmt.Println("📚 CONTEXT DETAILS:")

			// Safely check and display summary context
			summaryCount := 0
			if ragContext.SummaryContext != nil {
				summaryCount = len(ragContext.SummaryContext)
			}
			fmt.Printf("  • Found %d summary chunks\n", summaryCount)

			// Safely check and display document context
			docCount := 0
			if ragContext.DocumentContext != nil {
				docCount = len(ragContext.DocumentContext)
			}
			fmt.Printf("  • Found %d document chunks\n", docCount)

			// Safely check and display semantic context
			semanticCount := 0
			if ragContext.SemanticContext != nil {
				semanticCount = len(ragContext.SemanticContext)
			}
			fmt.Printf("  • Found %d semantic links\n", semanticCount)

			// Display summary chunks if available
			if summaryCount > 0 {
				fmt.Println("\n📄 SUMMARY CHUNKS:")
				for i, chunk := range ragContext.SummaryContext {
					// Just print the chunk directly as it's a string, not a struct with Content field
					firstLine := strings.Split(chunk, "\n")[0]
					fmt.Printf("  [%d] %s\n", i+1, firstLine)
				}
			}

			// Display semantic links if available
			if semanticCount > 0 {
				fmt.Println("\n🔗 SEMANTIC LINKS:")
				for i, link := range ragContext.SemanticContext {
					// Just print the link as a string, since it's not a struct with fields
					fmt.Printf("  [%d] %s\n", i+1, link)
				}
			}
		} else {
			summaryCount := len(ragContext.SummaryContext)
			docCount := len(ragContext.DocumentContext)
			semanticCount := len(ragContext.SemanticContext)

			fmt.Printf("📊 STATS: %d summary chunks, %d document chunks, %d semantic links\n",
				summaryCount, docCount, semanticCount)
			fmt.Println("(Use --verbose flag to see full context details)")
		}
	}

	// Save results to file if desired
	if os.Getenv("SAVE_RESULTS") == "true" {
		resultData := map[string]interface{}{
			"query":      *query,
			"domain":     *domain,
			"answer":     answer,
			"context":    ragContext,
			"timestamp":  time.Now(),
			"duration_s": duration.Seconds(),
		}

		resultJSON, _ := json.MarshalIndent(resultData, "", "  ")
		filename := fmt.Sprintf("hsg-search-%d.json", time.Now().Unix())
		if err := os.WriteFile(filename, resultJSON, 0644); err != nil {
			log.Printf("Warning: Could not save results to file: %v", err)
		} else {
			fmt.Printf("\n💾 Results saved to: %s\n", filename)
		}
	}
}

// Add this new function to handle schema setup
func setupHSGSchema() error {
	// Get Weaviate client
	weaviateClient, err := getWeaviateClient()
	if err != nil {
		return fmt.Errorf("failed to get Weaviate client: %w", err)
	}

	ctx := context.Background()

	// Define and create the necessary schema classes
	// These match what's in cmd/setup-hsg-schema/main.go

	// SummaryChunk class
	summaryChunkClass := &models.Class{
		Class:       "SummaryChunk",
		Description: "Summary chunks for hierarchical semantic graph",
		Properties: []*models.Property{
			{
				Name:        "content",
				DataType:    []string{"text"},
				Description: "The summary content",
			},
			{
				Name:        "domain",
				DataType:    []string{"string"},
				Description: "The knowledge domain",
			},
			// ...other properties...
		},
	}

	// SemanticLink class
	semanticLinkClass := &models.Class{
		Class:       "SemanticLink",
		Description: "Semantic links between concepts",
		Properties: []*models.Property{
			{
				Name:            "sourceTerm",
				DataType:        []string{"string"},
				Description:     "Source term for the link",
				IndexFilterable: ptrBool(true),
				IndexSearchable: ptrBool(true),
			},
			// ...other properties...
		},
	}

	// SubjectAreaContent class
	subjectContentClass := &models.Class{
		Class:       "SubjectAreaContent",
		Description: "Subject area content for domain-specific knowledge",
		Properties: []*models.Property{
			{
				Name:            "title",
				DataType:        []string{"string"},
				Description:     "Content title",
				IndexFilterable: ptrBool(true),
				IndexSearchable: ptrBool(true),
			},
			// ...other properties...
		},
	}

	// Create or update the classes
	for _, class := range []*models.Class{summaryChunkClass, semanticLinkClass, subjectContentClass} {
		// Check if class exists
		exists, err := weaviateClient.Schema().ClassExistenceChecker().WithClassName(class.Class).Do(ctx)
		if err != nil {
			return fmt.Errorf("error checking if class %s exists: %w", class.Class, err)
		}

		if !exists {
			if err := weaviateClient.Schema().ClassCreator().WithClass(class).Do(ctx); err != nil {
				return fmt.Errorf("failed to create class %s: %w", class.Class, err)
			}
			fmt.Printf("✅ Created Weaviate class: %s\n", class.Class)
		} else {
			fmt.Printf("ℹ️ Class %s already exists\n", class.Class)
		}
	}

	return nil
}

func ptrBool(b bool) *bool {
	return &b
}

func getWeaviateClient() (*weaviate.Client, error) {
	// This is a simplified version that should be adjusted to match your actual Weaviate client initialization
	weaviateURL := os.Getenv("WEAVIATE_URL_HOST")
	if weaviateURL == "" {
		weaviateURL = "http://localhost:8081" // Default URL
	}

	cfg := weaviate.Config{
		Host:   weaviateURL,
		Scheme: "http",
	}

	return weaviate.New(cfg), nil
}
