package main

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/InfoFlow/InfoIn/rawText/query"
	"esp-organizer/internal/InfoFlow/InfoOut/llm"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"esp-organizer/internal/InfoFlow/InfoStore/skills"
	"esp-organizer/internal/config"
	"esp-organizer/internal/models"
	"esp-organizer/internal/utils"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	queryTimeout      = 30 * time.Second
	maxDisplayResults = 3
)

func main() {
	// Centralized config loading at the start of the application.
	projectRoot, err := utils.FindProjectRoot()
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}
	config.LoadConfig(projectRoot, ".env")

	// Initialize context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	// 1. Initialize MongoDB
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer func() {
		if err := mongoDb.Client.Disconnect(ctx); err != nil {
			log.Printf("Warning: MongoDB disconnect error: %v", err)
		}
	}()

	// 2. Initialize Weaviate
	log.Println("Connecting to Weaviate...")
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}
	defer db.CloseWeaviate()

	// 3. Initialize LLAMA Client
	llamaClient := llm.NewLlamaClient()
	if llamaClient == nil {
		log.Fatal("Failed to create LlamaClient")
	}
	log.Println("LLAMA Client initialized successfully")

	// 4. Create Services
	skillsCollectionName := os.Getenv("SKILLS_COLLECTION")
	if skillsCollectionName == "" {
		skillsCollectionName = "skills"
	}
	skillService := skills.NewSkillService(
		mongoDb.Database.Collection(skillsCollectionName),
		db.GetWeaviateClient(),
		llamaClient,
	)

	queryHandler := query.NewQueryHandler(skillService, llamaClient)
	if queryHandler == nil {
		log.Fatal("Failed to create QueryHandler")
	}

	// 5. Run Test Queries
	runTestQueries(ctx, queryHandler)

	// 6. Test JSON Serialization
	testJSONResponse(ctx, queryHandler)

	log.Println("=== Query System Test Complete ===")
}

func runTestQueries(ctx context.Context, handler models.QueryProcessor) {
	testQueries := []string{
		"working memory",
		"problem solving skills",
		"communication and speaking",
		"reading comprehension",
		"fine motor skills",
		"emotional intelligence",
	}

	filter := createTestFilter()

	fmt.Println("=== Testing Query System ===")
	for i, testQuery := range testQueries {
		testSingleQuery(ctx, handler, i+1, testQuery, filter)
	}
}

func createTestFilter() *models.Filter {
	return &models.Filter{
		Operator: "And",
		Operands: []models.Operand{
			{
				Path:        []string{"development_age"},
				Operator:    "GreaterThanEqual",
				ValueNumber: 5,
			},
		},
	}
}

func testSingleQuery(ctx context.Context, handler models.QueryProcessor, idx int, queryStr string, filter *models.Filter) {
	fmt.Printf("--- Test Query %d: %s ---\n", idx, queryStr)

	// All semantic queries now go to the SemanticLinks class
	response, err := handler.ProcessQuery(ctx, queryStr, filter, "SemanticLinks")
	if err != nil {
		log.Printf("Query failed: %v", err)
		return
	}

	printQueryResults(response)
}

func printQueryResults(response *models.QueryResponse) {
	// Print MongoDB results
	printMongoResults(response.MongoResults)

	// Print Weaviate results
	printWeaviateResults(response.WeaviateResults)

	// Print Related Skills
	printRelatedSkills(response.RelatedSkills)

	// Print Synthesis
	fmt.Printf("Synthesis: %s\n\n", config.TruncateString(response.Synthesis, 120))
}

func printMongoResults(results []models.Skill) {
	fmt.Printf("MongoDB Results: %d skills found\n", len(results))
	for i, skill := range results {
		if i >= maxDisplayResults {
			fmt.Printf("  ... and %d more\n", len(results)-maxDisplayResults)
			break
		}
		fmt.Printf("  - %s (Age: %d): %s\n",
			skill.Name,
			skill.DevelopmentAge,
			config.TruncateString(skill.Description, 100))
	}
}

func printWeaviateResults(results []map[string]interface{}) {
	fmt.Printf("Weaviate Results: %d semantic links found\n", len(results))
	for i, result := range results {
		if i >= maxDisplayResults {
			fmt.Printf("  ... and %d more semantic links\n", len(results)-maxDisplayResults)
			break
		}
		if properties, ok := result["properties"].(map[string]interface{}); ok {
			fmt.Printf("  - Link: %s -> %s (Certainty: %.3f)\n    Context: %s\n",
				config.GetStringValue(properties, "source_term"),
				config.GetStringValue(properties, "target_term"),
				config.GetCertaintyValue(properties),
				config.TruncateString(config.GetStringValue(properties, "context"), 80))
		}
	}
}

func printRelatedSkills(skills []models.Skill) {
	fmt.Printf("Related Skills: %d parent skills found\n", len(skills))
	for i, skill := range skills {
		if i >= maxDisplayResults-1 { // Show one less for related skills
			fmt.Printf("  ... and %d more\n", len(skills)-(maxDisplayResults-1))
			break
		}
		fmt.Printf("  - %s: %s\n",
			skill.Name,
			config.TruncateString(skill.Description, 80))
	}
}

func testJSONResponse(ctx context.Context, handler models.QueryProcessor) {
	fmt.Println("\n=== Testing JSON Response ===")
	response, err := handler.ProcessQuery(ctx, "reading skills", nil, "SemanticLinks")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %v", err)
	}

	fmt.Printf("JSON Response length: %d characters\n", len(jsonData))
	fmt.Printf("Sample JSON (first 500 chars):\n%.500s...\n", string(jsonData))
}
