package main

import (
	"bufio"
	"context"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	queryTimeout = 30 // seconds
)

func main() {
	// Load environment variables from .env file. This is the crucial fix.
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on system environment variables.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	// 1. Initialize MongoDB
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoDb.Client.Disconnect(ctx); err != nil {
			log.Printf("Warning: Error disconnecting from MongoDB: %v", err)
		}
	}()
	// WARNING: This script can clear data - require explicit confirmation
	fmt.Println("⚠️  WARNING: Database Clearing Script")
	fmt.Println("This script will clear all data from MongoDB and Weaviate.")
	fmt.Println("Data volumes will be preserved but emptied.")
	fmt.Println("")
	fmt.Print("Are you sure you want to continue? (type 'yes' to confirm): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(input) != "yes" {
		log.Println("Operation cancelled by user.")
		return
	}

	log.Println("Connecting to Weaviate...")
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}
	defer db.CloseWeaviate()

	// Clear MongoDB collections (but preserve database structure)
	log.Println("Clearing MongoDB collections...")

	collectionsToClear := []string{
		"skills",
		"sources",
		"immunology_chapters",
		"case_studies",
		"medical_terms",
		"processing_batches",
	}

	for _, collectionName := range collectionsToClear {
		collection := mongoDb.Database.Collection(collectionName)
		result, err := collection.DeleteMany(ctx, map[string]interface{}{})
		if err != nil {
			log.Printf("Warning: Failed to clear %s collection: %v", collectionName, err)
		} else {
			log.Printf("Successfully cleared %s collection (%d documents)", collectionName, result.DeletedCount)
		}
	}

	// Clear Weaviate classes (but preserve class structure)
	log.Println("Clearing Weaviate classes...")
	weaviateClasses := []string{"EducationalSkills", "SubjectAreaContent", "Skills"} // Include old "Skills" for cleanup

	for _, className := range weaviateClasses {
		exists, err := db.WeaviateCollectionExists(ctx, className)
		if err != nil {
			log.Printf("Warning: Failed to check if %s class exists: %v", className, err)
			continue
		}

		if exists {
			weaviateClient := db.GetWeaviateClient()
			objects, err := weaviateClient.Data().ObjectsGetter().WithClassName(className).Do(ctx)
			if err != nil {
				log.Printf("Warning: Failed to get objects from %s: %v", className, err)
			} else {
				log.Printf("Found %d objects in %s class, deleting...", len(objects), className)
				if len(objects) > 0 {
					for _, obj := range objects {
						if err := db.WeaviateDeleteDocument(ctx, className, obj.ID.String()); err != nil {
							log.Printf("Warning: Failed to delete object %s from %s: %v", obj.ID, className, err)
						}
					}
				}
				log.Printf("Successfully cleared Weaviate %s class", className)
			}
		} else {
			log.Printf("%s class does not exist in Weaviate, skipping.", className)
		}
	}

	log.Println("✅ Database clearing completed successfully!")
	log.Println("📁 Data volumes preserved for future use")
	log.Println("🚀 Ready for fresh data upload")
}
