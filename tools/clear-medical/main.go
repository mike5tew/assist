package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	"esp-organizer/internal/store/db"
	"esp-organizer/internal/utils"

	weaviate "github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"go.mongodb.org/mongo-driver/bson"
)

// --- UPDATED GLOBALS REFLECTING NEW STRUCTURE ---

// MongoDB collections that are part of the Immunology domain's new structure.
// NOTE: We assume 'immunology_content' and 'immunology_terms' are the primary collections.
// var medicalCollections = []string{
// 	"immunology_content", // Chapters and Case Studies now live here
// 	"immunology_terms",   // Definitions and Keywords now live here
// 	"immunology_sources", // Administrative/book metadata (as established previously)
// 	"unified_content",    // Kept for backward compatibility, filtered by subject
// 	"subject_content",    // Kept for backward compatibility, filtered by subject
// }

// // Weaviate classes that likely contain vectorized Immunology content.
// // NOTE: Classes are typically singular (e.g., MedicalExcerpt).
// var medicalClasses = []string{
// 	"SubjectAreaContent", // The primary class for domain content (chapters, case studies)
// 	"MedicalTerm",        // Class for definitions/terms
// 	"SemanticLink",       // Class for relationships/keywords
// }

// A map of collections that require a filter to only delete "immunology" content.
var medicalContentFilters = map[string]bson.M{
	"unified_content": {"domain": "immunology"},
	"subject_content": {"domain": "immunology"},
	// Note: immunology_content and immunology_terms are cleared completely as they are domain-specific
}

// --- END UPDATED GLOBALS ---

func main() {
	// Ensure logs directory exists before initializing logging
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			log.Fatalf("Failed to create logs directory: %v", err)
		}
	}
	utils.InitDefaultLogging()

	// Load environment variables from .env file if present
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// Connect to MongoDB
	log.Println("Connecting to MongoDB at", os.Getenv("MONGODB_URI"))
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDb.Client.Disconnect(context.Background())
	log.Println("✅ Successfully connected to MongoDB")
	log.Println("Using MongoDB database:", mongoDb.Database.Name())

	// Add flag for HSG-specific clearing
	hsgFlag := flag.Bool("hsg", false, "Clear HSG-related collections (summary_chunks, semantic_links, etc.)")
	flag.Parse()

	log.Printf("🧹 Clearing all medical content from MongoDB database: %s", mongoDb.Database.Name())

	// Basic collections always cleared
	collections := []string{
		"immunology_content",
		"medical_terms",
		"immunology_terms", // Explicitly add this collection
		"extraction_jobs",
	}

	// Add HSG-specific collections when flag is set
	if *hsgFlag {
		hsgCollections := []string{
			"subject_content",
			"summary_chunks",
			"semantic_links",
		}
		collections = append(collections, hsgCollections...)
		log.Printf("HSG flag enabled - will clear additional collections: %v", hsgCollections)
	}

	// Clear MongoDB collections
	for _, collName := range collections {
		log.Printf("   Attempting to clear collection: %s", collName)

		collection := mongoDb.Database.Collection(collName)

		// Determine filter: use specific filter for backward-compatible collections, otherwise delete all
		filter := bson.M{}
		if specificFilter, exists := medicalContentFilters[collName]; exists {
			filter = specificFilter
			log.Printf("   -> Using FILTERED deletion (domain: 'immunology')...")
		} else {
			// For domain-specific collections (immunology_content, immunology_terms, etc.)
			log.Printf("   -> Using COMPLETE deletion (domain-specific collection)...")
		}

		// Get document count before deletion for reporting
		count, err := collection.CountDocuments(context.Background(), filter)
		if err != nil {
			log.Printf("   Error counting documents in collection %s: %v", collName, err)
			continue
		}

		if count == 0 {
			log.Printf("   No matching documents in collection '%s'. Nothing to clear.", collName)
			continue
		}

		// Delete documents
		deleteResult, err := collection.DeleteMany(context.Background(), filter)
		if err != nil {
			log.Printf("   ❌ Error deleting from collection %s: %v", collName, err)
			continue
		}

		log.Printf("✅ Cleared %d documents from MongoDB collection: %s", deleteResult.DeletedCount, collName)
	}

	// Clear Weaviate classes
	weaviateClient := db.GetWeaviateClient()
	if weaviateClient == nil {
		log.Println("❌ Failed to initialize Weaviate client")
		return
	}

	weaviateClasses := []string{"MedicalExcerpt"}
	if *hsgFlag {
		weaviateClasses = append(weaviateClasses, []string{
			"SubjectAreaContent",
			"SemanticLinks",
		}...)
	}

	// Clear entire classes (schema + objects) for clean slate
	for _, className := range weaviateClasses {
		deleteWeaviateClass(weaviateClient, className)
	}

	// Special handling: Clear any leftover content in SemanticLinks or SubjectAreaContent
	// that might have been indexed under the 'immunology' subject
	clearWeaviateContentBySubject(weaviateClient, "immunology")

	log.Println("✅ All medical data has been cleared.")
}

// clearWeaviateClass clears all objects of a specific class (kept for structural reference)
func deleteWeaviateClass(client *weaviate.Client, className string) {
	ctx := context.Background()

	// Check if class exists
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil || !exists {
		log.Printf("ℹ️ Weaviate class %s does not exist", className)
		return
	}

	// Delete the entire class (objects + schema)
	err = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
	if err != nil {
		log.Printf("❌ Error deleting Weaviate class %s: %v", className, err)
		return
	}

	log.Printf("✅ Deleted Weaviate class: %s", className)
}

// Advanced version with options (kept for structural reference)
// NOTE: Removed unused ClearOptions struct
func clearWeaviateContentBySubject(client *weaviate.Client, subjectName string) {
	//ctx := context.Background()

	// You can customize which class to search in, or search all classes
	classesToCheck := []string{"SemanticLink", "MedicalTerm", "SubjectAreaContent"} // Updated names

	totalCleared := 0

	for _, className := range classesToCheck {
		cleared := clearSubjectFromClass(client, className, subjectName)
		totalCleared += cleared
	}

	log.Printf("✅ Cleared %d objects with subject '%s' from Weaviate", totalCleared, subjectName)
}

// Clear specific subject from a particular class (kept for structural reference)
func clearSubjectFromClass(client *weaviate.Client, className, subjectName string) int {
	ctx := context.Background()

	// Check if class exists
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil || !exists {
		log.Printf("ℹ️ Weaviate class %s does not exist, skipping.", className)
		return 0
	}

	// Count objects with this subject
	count, err := countObjectsWithField(client, className, "domain", subjectName) // Using 'domain' field for subject check
	if err != nil {
		log.Printf("❌ Error counting objects in %s with subject '%s': %v", className, subjectName, err)
		return 0
	}

	if count == 0 {
		log.Printf("ℹ️ No objects with subject '%s' in class %s", subjectName, className)
		return 0
	}

	log.Printf("🧹 Found %d objects with subject '%s' in Weaviate class %s. Clearing...", count, subjectName, className)

	// Delete objects with this subject
	subjectFilter := filters.Where().
		WithPath([]string{"domain"}). // Using 'domain' field for subject check
		WithOperator(filters.Equal).
		WithValueString(subjectName)

	result, err := client.Batch().ObjectsBatchDeleter().
		WithClassName(className).
		WithWhere(subjectFilter).
		Do(ctx)

	if err != nil {
		log.Printf("❌ Error clearing subject '%s' from class %s: %v", subjectName, className, err)
		return 0
	}
	if result == nil || result.Results == nil {
		log.Printf("❌ No result returned when clearing subject '%s' from class %s", subjectName, className)
		return 0
	}

	log.Printf("✅ Cleared %d objects with subject '%s' from class %s", count, subjectName, className)
	return count
}

// Helper function to count objects with specific field value (kept for structural reference)
func countObjectsWithField(client *weaviate.Client, className, fieldName, fieldValue string) (int, error) {
	ctx := context.Background()

	response, err := client.GraphQL().Aggregate().
		WithClassName(className).
		WithFields(
			graphql.Field{
				Name: "meta",
				Fields: []graphql.Field{
					{Name: "count"},
				},
			},
		).
		WithWhere(filters.Where().
			WithPath([]string{fieldName}).
			WithOperator(filters.Equal).
			WithValueString(fieldValue)).
		Do(ctx)

	if err != nil {
		return 0, err
	}

	return parseCountFromResponse(response, className), nil
}

// Helper function to parse count from GraphQL response (kept for structural reference)
func parseCountFromResponse(response *models.GraphQLResponse, className string) int {
	if response == nil || response.Data == nil {
		return 0
	}

	if aggregate, ok := response.Data["Aggregate"].(map[string]interface{}); ok {
		if classData, ok := aggregate[className].([]interface{}); ok && len(classData) > 0 {
			if firstItem, ok := classData[0].(map[string]interface{}); ok {
				if meta, ok := firstItem["meta"].(map[string]interface{}); ok {
					if countVal, ok := meta["count"].(float64); ok {
						return int(countVal)
					}
				}
			}
		}
	}

	return 0
}
