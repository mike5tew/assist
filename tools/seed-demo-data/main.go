package main

import (
	"context"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	log.Println("🚀 Starting MVP Demo Data Seeder...")

	// Find and load .env file from project root
	envPath := findEnvFile()
	log.Printf("📍 Looking for .env at: %s", envPath)

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️ Warning: Could not load .env file at %s: %v", envPath, err)
		log.Println("Using environment variables instead...")
	} else {
		log.Printf("✅ Loaded .env from: %s", envPath)
	}

	// CRITICAL FIX: Use localhost URI for host-machine execution
	// The .env has mongodb://admin:password123@mongodb:27017 (Docker DNS)
	// But this script runs on the HOST, so we need localhost
	mongoURILocal := "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin"
	os.Setenv("MONGODB_URI", mongoURILocal)
	log.Printf("✅ Override MONGODB_URI for host execution: localhost:27017")

	// Debug: Log what we're connecting with
	debugEnvVars()

	// CRITICAL FIX #2: Use localhost URI for Weaviate host execution
	// The .env has http://weaviate:8080 (Docker DNS)
	// But this script runs on the HOST, so we need localhost:8081
	weaviateURLLocal := "http://localhost:8081"
	os.Setenv("WEAVIATE_URL", weaviateURLLocal)
	log.Printf("✅ Override WEAVIATE_URL for host execution: %s", weaviateURLLocal)

	ctx := context.Background()

	// 1. Connect to Databases
	mongoClient, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Client.Disconnect(ctx)

	weaviateClient := db.GetWeaviateClient()
	if weaviateClient == nil {
		log.Fatalf("❌ Failed to get Weaviate client")
	}

	log.Println("✅ Connected to databases.")

	// 2. Define the demo data
	// This is a small, interconnected graph about X-Linked Agammaglobulinemia (XLA).
	demoLinks := []models.SemanticLink{
		{
			SourceTerm:   "X-Linked Agammaglobulinemia",
			TargetTerm:   "Primary Immunodeficiency",
			RelationType: "is_a",
			Context:      "X-linked agammaglobulinemia (XLA) is a primary immunodeficiency disease that affects the body's ability to fight infection.",
			Confidence:   0.99,
			Domain:       "immunology",
		},
		{
			SourceTerm:   "X-Linked Agammaglobulinemia",
			TargetTerm:   "BTK gene mutation",
			RelationType: "caused_by",
			Context:      "XLA is caused by mutations in the Bruton's tyrosine kinase (BTK) gene, which is critical for B-cell development.",
			Confidence:   0.98,
			Domain:       "immunology",
		},
		{
			SourceTerm:   "BTK gene mutation",
			TargetTerm:   "B-cell development failure",
			RelationType: "leads_to",
			Context:      "Mutations in the BTK gene lead to a failure of B-cell development, resulting in the absence of mature B cells.",
			Confidence:   0.97,
			Domain:       "immunology",
		},
		{
			SourceTerm:   "B-cell development failure",
			TargetTerm:   "Absent B-cells",
			RelationType: "results_in",
			Context:      "The failure of B-cell development results in profoundly low numbers or complete absence of B cells in the circulation.",
			Confidence:   0.99,
			Domain:       "immunology",
		},
		{
			SourceTerm:   "Absent B-cells",
			TargetTerm:   "Recurrent bacterial infections",
			RelationType: "causes",
			Context:      "The absence of B cells means the body cannot produce antibodies, leading to recurrent bacterial infections, especially from encapsulated bacteria.",
			Confidence:   0.96,
			Domain:       "immunology",
		},
		{
			SourceTerm:   "X-Linked Agammaglobulinemia",
			TargetTerm:   "Immunoglobulin replacement therapy",
			RelationType: "treated_by",
			Context:      "The standard treatment for XLA is lifelong immunoglobulin replacement therapy (IVIG) to provide passive immunity.",
			Confidence:   0.98,
			Domain:       "immunology",
		},
	}

	// 3. Get MongoDB collection
	semanticLinksCollection := mongoClient.Database.Collection("semantic_links")

	// 4. Insert data
	var insertedCount int
	for _, link := range demoLinks {
		link.ID = primitive.NewObjectID()
		link.CreatedAt = time.Now()
		link.CompositeKey = fmt.Sprintf("%s::%s::%s", link.SourceTerm, link.TargetTerm, link.RelationType)

		// Check if link already exists to make this script idempotent
		var existing models.SemanticLink
		err := semanticLinksCollection.FindOne(ctx, bson.M{"composite_key": link.CompositeKey}).Decode(&existing)
		if err == nil {
			log.Printf("⚠️ Link already exists, skipping: %s", link.CompositeKey)
			continue
		}

		// Insert into MongoDB
		_, err = semanticLinksCollection.InsertOne(ctx, link)
		if err != nil {
			log.Printf("❌ Failed to insert link into MongoDB: %v", err)
			continue
		}

		// Insert into Weaviate (without vector, as we'll query by text for the demo)
		properties := map[string]interface{}{
			"sourceTerm":   link.SourceTerm,
			"targetTerm":   link.TargetTerm,
			"relationType": link.RelationType,
			"context":      link.Context,
			"confidence":   link.Confidence,
			"domain":       link.Domain,
			"mongoId":      link.ID.Hex(), // Store the Mongo ID for cross-reference
		}

		// Create the object
		_, err = weaviateClient.Data().Creator().
			WithClassName("SemanticLinks").
			WithProperties(properties).
			Do(ctx)

		if err != nil {
			log.Printf("❌ Failed to insert link into Weaviate: %v", err)
			// For MVP, we can skip Weaviate if it fails - MongoDB has the data
			log.Printf("⚠️ Continuing with MongoDB insertion only (Weaviate optional for MVP)")
			continue
		}

		insertedCount++
		log.Printf("✅ Seeded link: %s", link.CompositeKey)
	}

	log.Printf("🎉 Seeding complete. Inserted %d new links into MongoDB.", insertedCount)
}

// findEnvFile searches for .env file starting from current directory and going up
func findEnvFile() string {
	// Try common paths relative to where the binary is run
	searchPaths := []string{
		".env",          // Current directory
		"../.env",       // Parent
		"../../.env",    // Grandparent
		"../../../.env", // Great-grandparent
		"../../../.env", // From esp-organizer root
	}

	for _, path := range searchPaths {
		absPath, err := filepath.Abs(path)
		if err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath
			}
		}
	}

	// Fallback to current working directory
	return ".env"
}

// debugEnvVars logs the environment variables being used
func debugEnvVars() {
	vars := map[string]string{
		"MONGO_INITDB_ROOT_USERNAME": "",
		"MONGO_INITDB_ROOT_PASSWORD": "",
		"MONGO_DB_NAME":              "",
		"MONGODB_URI":                "",
		"WEAVIATE_URL":               "",
	}

	log.Println("📊 Environment Variables:")
	for key := range vars {
		val := os.Getenv(key)
		if val == "" {
			log.Printf("  ❌ %s: NOT SET", key)
		} else if key == "MONGO_INITDB_ROOT_PASSWORD" || key == "MONGODB_URI" {
			log.Printf("  ✅ %s: [MASKED]", key)
		} else {
			log.Printf("  ✅ %s: %s", key, val)
		}
	}
}
