package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// List of collections the app *should* be using (for focused check)
var focusedMongoCollections = []string{
	"immunology_content",
	"immunology_terms",
	"immunology_sources",
	"unified_content",
	"semantic_links",
}

// List of Weaviate classes the app *should* be using
var weaviateClasses = []string{
	"SubjectAreaContent",
	"MedicalTerm",
	"SemanticLink",
}

func main() {
	// 1. Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_NAME")

	if mongoURI == "" {
		mongoURI = "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin"
	}
	if dbName == "" {
		dbName = "esp_organizer"
	}

	log.Printf("Connecting to MongoDB at %s", mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}
	log.Printf("✅ Successfully connected to MongoDB")
	log.Printf("Using MongoDB database: %s", dbName)

	// 2. LIST ALL EXISTING COLLECTIONS
	fmt.Println("\n=== MongoDB ALL Existing Collections ===")
	allNames, err := listAllMongoCollections(ctx, client, dbName)
	if err != nil {
		log.Printf("❌ WARNING: Failed to list all collection names: %v", err)
	} else {
		// This map tracks which collections were already checked
		checkedCollections := make(map[string]bool)

		// 3. CHECK FOCUSED COLLECTIONS
		fmt.Println("\n=== MongoDB Focused Medical Collections ===")
		for _, collName := range focusedMongoCollections {
			checkCollection(ctx, client, dbName, collName)
			checkedCollections[collName] = true
		}

		// 4. CHECK OTHER COLLECTIONS (found but not in the focused list)
		fmt.Println("\n=== MongoDB Other Existing Collections ===")
		foundOther := false
		for _, name := range allNames {
			if !checkedCollections[name] {
				checkCollection(ctx, client, dbName, name)
				foundOther = true
			}
		}
		if !foundOther {
			fmt.Println("ℹ️  No other collections found outside the focused list.")
		}
	}

	// 5. Check Weaviate classes
	fmt.Println("\n=== Weaviate Current Medical Classes ===")
	for _, className := range weaviateClasses {
		checkWeaviateClass(className)
	}
}

// listAllMongoCollections retrieves the names of all non-system collections in the database.
func listAllMongoCollections(ctx context.Context, client *mongo.Client, dbName string) ([]string, error) {
	// The filter excludes system collections like 'system.indexes' and 'fs.chunks'
	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: "^(?!system\\.)(?!fs\\.)(?!semantic_links_temp$)"}}}}
	names, err := client.Database(dbName).ListCollectionNames(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(names) == 0 {
		fmt.Println("ℹ️  No non-system collections found in the database.")
	} else {
		fmt.Printf("✅ Found %d collections: %v\n", len(names), names)
	}
	return names, nil
}

func checkCollection(ctx context.Context, client *mongo.Client, dbName, collectionName string) {
	// ... (The rest of the checkCollection and checkWeaviateClass functions remain the same) ...

	// The checkCollection function remains as defined in the previous response,
	// including the logic for counting documents and printing content snippets.
	// I'm omitting the body here for brevity, but it should be included in your file.

	// Count logic:
	count, err := client.Database(dbName).Collection(collectionName).CountDocuments(ctx, bson.M{})
	if err != nil {
		fmt.Printf("❌ Error checking collection %s: %v\n", collectionName, err)
		return
	}

	fmt.Printf("📊 %s: %d documents\n", collectionName, count)

	// Log Retrieval Activity: If data exists, show what is being retrieved.
	if count > 0 {
		fmt.Printf("   📝 RETRIEVAL LOG - SAMPLE:\n")
		cursor, err := client.Database(dbName).Collection(collectionName).Find(
			ctx,
			bson.M{},
			options.Find().SetLimit(2),
		)
		if err != nil {
			fmt.Printf("   ❌ Error retrieving samples: %v\n", err)
			return
		}
		defer cursor.Close(ctx)

		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			fmt.Printf("   ❌ Error decoding samples: %v\n", err)
			return
		}

		for i, doc := range results {
			fmt.Printf("   📄 Document %d: ", i+1)

			// --- LOGGING KEY IDENTIFIERS ---
			identifiers := ""
			if id, ok := doc["_id"]; ok {
				identifiers += fmt.Sprintf("ID: %v, ", id)
			}
			if title, ok := doc["title"]; ok {
				identifiers += fmt.Sprintf("Title: %v, ", title)
			}
			if bookId, ok := doc["bookId"]; ok {
				identifiers += fmt.Sprintf("BookID: %v, ", bookId)
			}
			if isbn, ok := doc["isbn13"]; ok {
				identifiers += fmt.Sprintf("ISBN: %v, ", isbn)
			}
			fmt.Println(identifiers)

			// --- CRITICAL DEBUG: LOGGING CONTENT TEXT ---
			// Check common content fields. Adjust 'content_text' to match your schema.
			if content, ok := doc["content"].(string); ok && len(content) > 0 {
				snippet := content
				if len(content) > 120 {
					snippet = content[:120] + "..."
				}
				fmt.Printf("      CONTENT SNIPPET: \"%s\"\n", snippet)
			} else if content, ok := doc["excerpt_text"].(string); ok && len(content) > 0 {
				snippet := content
				if len(content) > 120 {
					snippet = content[:120] + "..."
				}
				fmt.Printf("      CONTENT SNIPPET: \"%s\"\n", snippet)
			} else if content, ok := doc["definition"].(string); ok && len(content) > 0 {
				snippet := content
				if len(content) > 120 {
					snippet = content[:120] + "..."
				}
				fmt.Printf("      CONTENT SNIPPET (Definition): \"%s\"\n", snippet)
			} else {
				fmt.Printf("      CONTENT SNIPPET: [WARNING: Content field is missing or empty]\n")
			}
		}
	}
}

func checkWeaviateClass(className string) {
	// This function body is identical to the previous response.
	cfg := weaviate.Config{
		Host:   "localhost:8081",
		Scheme: "http",
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		fmt.Printf("❌ Error connecting to Weaviate: %v\n", err)
		return
	}

	ctx := context.Background()
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil {
		fmt.Printf("❌ Error checking Weaviate class %s: %v\n", className, err)
		return
	}

	if exists {
		// Use the correct syntax for aggregate query
		result, err := client.GraphQL().Aggregate().
			WithClassName(className).
			WithFields(graphql.Field{Name: "meta", Fields: []graphql.Field{{Name: "count"}}}).
			Do(ctx)

		if err != nil {
			fmt.Printf("❌ Error getting count for Weaviate class %s: %v\n", className, err)
			return
		}

		var count int
		if result != nil && result.Data != nil {
			if aggregations, ok := result.Data["Aggregate"].(map[string]interface{}); ok {
				if classAggs, ok := aggregations[className].([]interface{}); ok && len(classAggs) > 0 {
					if meta, ok := classAggs[0].(map[string]interface{})["meta"].(map[string]interface{}); ok {
						if countVal, ok := meta["count"].(float64); ok {
							count = int(countVal)
						}
					}
				}
			}
		}

		fmt.Printf("📊 %s: %d objects\n", className, count)
	} else {
		fmt.Printf("ℹ️  %s: class does not exist\n", className)
	}
}
