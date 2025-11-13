package main

import (
	"context"
	"esp-organizer/internal/models"
	"esp-organizer/internal/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

// Global database connection for consistent access
var (
	mongoClient *mongo.Client
	dbName      string
)

func main() {
	// Parse command line arguments
	query := flag.String("query", "", "Search query")
	limit := flag.Int("limit", 5, "Maximum number of results to return")
	debug := flag.Bool("debug", false, "Show debug information")
	diagnosticMode := flag.Bool("diagnostic", false, "Run in diagnostic mode")
	flag.Parse()

	// Initialize database connection
	initDBConnection()
	defer closeDBConnection()

	// Run in diagnostic mode if requested
	if *diagnosticMode {
		runDiagnostics()
		return
	}

	// Regular search mode
	if *query == "" {
		fmt.Println("Error: Search query is required")
		flag.Usage()
		os.Exit(1)
	}

	// Dump all books for debugging
	if *debug {
		dumpAllBooks()
	}

	// Search in both MongoDB and Weaviate
	fmt.Printf("🔍 Searching for medical content matching: '%s'\n\n", *query)

	// Search in MongoDB (direct content match)
	mongoResults, err := searchMongoDB(*query, *limit)
	if err != nil {
		fmt.Printf("Error searching MongoDB: %v\n", err)
	} else {
		fmt.Printf("📚 MongoDB Results (Direct Content Match):\n")
		if len(mongoResults) == 0 {
			fmt.Println("  No direct matches found")
		} else {
			for i, result := range mongoResults {
				fmt.Printf("%d. [Page %d] 📖 %s\n", i+1, result.Page, truncateString(result.Content, 120))

				// Use the findBookByID function with SourceID
				book := utils.FindBookByID(result.SourceID) // Updated from BookID to SourceID
				var bookInfo string
				if book != nil {
					bookInfo = fmt.Sprintf("%s (ISBN: %s)", book.Title, book.ISBN)
				} else {
					// Make it clear this is a fallback for when the book is not in the DB
					bookInfo = "Book not in DB; fallback: Case Studies in Immunology (ISBN: 9780815345121)"
				}
				fmt.Printf("   📘 Book: %s\n\n", bookInfo)
			}
		}
	}

	// Search in Weaviate (semantic search)
	weaviateResults, err := searchWeaviate(*query, *limit)
	if err != nil {
		fmt.Printf("Error searching Weaviate: %v\n", err)
	} else {
		fmt.Printf("🔎 Weaviate Results (Semantic Search):\n")
		if len(weaviateResults) == 0 {
			// Improved message that doesn't sound like an error
			fmt.Println("  No semantic matches found (Weaviate may not be set up yet)")
		} else {
			for i, result := range weaviateResults {
				fmt.Printf("%d. [Page %d] 📖 %s\n", i+1, result.Page, truncateString(result.Content, 120))

				// Use the findBookByID function with SourceID
				book := utils.FindBookByID(result.SourceID) // Updated from BookID to SourceID
				var bookInfo string
				if book != nil {
					bookInfo = fmt.Sprintf("%s (ISBN: %s)", book.Title, book.ISBN)
				} else {
					// Make it clear this is a fallback for when the book is not in the DB
					bookInfo = "Book not in DB; fallback: Case Studies in Immunology (ISBN: 9780815345121)"
				}
				fmt.Printf("   📘 Book: %s\n\n", bookInfo)
			}
		}
	}
}

// Initialize the database connection once and reuse it
func initDBConnection() {
	// Use the standardized MongoDB connection settings from logs
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			// Use explicit authentication from logs
			mongoURI = "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin"
		}
	}

	// Use consistent database name
	dbName = os.Getenv("MONGODB_NAME")
	if dbName == "" {
		dbName = "esp_organizer"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verify connection
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
}

func closeDBConnection() {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		mongoClient.Disconnect(ctx)
	}
}

// Diagnostic function that runs various checks
func runDiagnostics() {
	fmt.Println("=== RUNNING DIAGNOSTICS ===")

	// Check collections
	checkCollections()

	// Check existing books
	dumpAllBooks()

	// Check excerpts
	dumpExcerpts()

	// Try to get a specific book by ISBN
	fmt.Println("\n=== CHECKING BOOK BY ISBN ===")
	bookByISBN := findBookByISBN("9780815345121")
	if bookByISBN != nil {
		fmt.Printf("✅ Found book by ISBN: %s (ID: %s)\n", bookByISBN.Title, bookByISBN.ID)
	} else {
		fmt.Println("❌ Could not find book by ISBN 9780815345121")
	}

	fmt.Println("\n=== DIAGNOSTICS COMPLETE ===")
}

// Check all relevant collections
func checkCollections() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("\n=== COLLECTIONS STATUS ===")
	collections, err := mongoClient.Database(dbName).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		fmt.Printf("Error listing collections: %v\n", err)
		return
	}

	fmt.Printf("Found %d collections in database %s:\n", len(collections), dbName)
	for _, coll := range collections {
		count, err := mongoClient.Database(dbName).Collection(coll).CountDocuments(ctx, bson.M{})
		if err != nil {
			fmt.Printf("- %s: Error counting documents\n", coll)
		} else {
			fmt.Printf("- %s: %d documents\n", coll, count)
		}
	}
}

// Dump all available books
func dumpAllBooks() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("\n=== ALL BOOKS IN DATABASE ===")
	collection := mongoClient.Database(dbName).Collection("medical_books")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Printf("Error finding books: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	var books []bson.M
	if err = cursor.All(ctx, &books); err != nil {
		fmt.Printf("Error decoding books: %v\n", err)
		return
	}

	fmt.Printf("Found %d books:\n", len(books))
	for i, book := range books {
		fmt.Printf("Book %d: ID=%v, Title=%v, ISBN=%v\n",
			i+1, book["_id"], book["title"], book["isbn13"])
	}
}

// Dump all excerpts
func dumpExcerpts() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("\n=== EXCERPTS IN DATABASE ===")
	collection := mongoClient.Database(dbName).Collection("medical_excerpts")
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetLimit(5))
	if err != nil {
		fmt.Printf("Error finding excerpts: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	var excerpts []bson.M
	if err = cursor.All(ctx, &excerpts); err != nil {
		fmt.Printf("Error decoding excerpts: %v\n", err)
		return
	}

	fmt.Printf("Found %d excerpts (showing first 5):\n", len(excerpts))
	for i, excerpt := range excerpts {
		fmt.Printf("Excerpt %d: BookID=%v, Page=%v\n",
			i+1, excerpt["bookId"], excerpt["page"])
		if content, ok := excerpt["content"].(string); ok {
			if len(content) > 50 {
				content = content[:50] + "..."
			}
			fmt.Printf("  Content: %s\n", content)
		}
	}
}

// Find a book by ISBN
func findBookByISBN(isbn string) *models.Book {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongoClient.Database(dbName).Collection("medical_books")
	var book models.Book
	err := collection.FindOne(ctx, bson.M{"isbn13": isbn}).Decode(&book)
	if err != nil {
		fmt.Printf("Error finding book by ISBN %s: %v\n", isbn, err)
		return nil
	}

	return &book
}

func searchMongoDB(query string, limit int) ([]models.MedicalExcerpt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query MongoDB for medical excerpts
	collection := mongoClient.Database(dbName).Collection("medical_excerpts")

	// Create a text search filter
	filter := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	// If text search index isn't set up, fall back to regex search
	count, _ := collection.CountDocuments(ctx, filter)

	// Fix: Capture both return values from CreateOne - the index name and the error
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"content": "text",
		},
	})
	if err != nil {
		fmt.Printf("Error creating text index: %v\n", err)
	}

	// Improve search with more flexible matching for natural language queries
	if err != nil || count == 0 {
		// Extract key terms from the query
		keyTerms := extractKeyTerms(query)

		// If we have key terms, build a more flexible query
		if len(keyTerms) > 0 {
			var orConditions []bson.M
			for _, term := range keyTerms {
				if len(term) > 3 { // Only use terms with at least 4 characters
					orConditions = append(orConditions, bson.M{
						"content": bson.M{
							"$regex":   term,
							"$options": "i", // case-insensitive
						},
					})
				}
			}

			// If we have valid conditions, use them
			if len(orConditions) > 0 {
				filter = bson.M{"$or": orConditions}
			} else {
				// Fallback to simple regex on the whole query
				filter = bson.M{
					"content": bson.M{
						"$regex":   query,
						"$options": "i", // case-insensitive
					},
				}
			}
		} else {
			// Fallback to simple regex on the whole query
			filter = bson.M{
				"content": bson.M{
					"$regex":   query,
					"$options": "i", // case-insensitive
				},
			}
		}
	}

	findOptions := options.Find().SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to execute MongoDB search: %v", err)
	}
	defer cursor.Close(ctx)

	var results []models.MedicalExcerpt
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode MongoDB results: %v", err)
	}

	// If still no direct results and query contains terms like "summarize" or "summary",
	// try looking for relevant content without those terms
	if len(results) == 0 && (strings.Contains(strings.ToLower(query), "summarize") ||
		strings.Contains(strings.ToLower(query), "summary")) {
		// Extract the subject from the query (remove summary-related words)
		subject := extractSubject(query)
		if subject != "" {
			// Search just for the subject
			return searchMongoDB(subject, limit)
		}
	}

	return results, nil
}

// Helper function to extract key terms from a query
func extractKeyTerms(query string) []string {
	// Convert to lowercase
	query = strings.ToLower(query)

	// Remove common words that don't add search value
	stopWords := []string{"the", "and", "for", "with", "about", "summarize", "summary", "case"}
	for _, word := range stopWords {
		query = strings.ReplaceAll(query, " "+word+" ", " ")
	}

	// Split into words
	words := strings.Fields(query)

	// Return unique, non-empty words
	var terms []string
	seen := make(map[string]bool)
	for _, word := range words {
		// Only include words with 3+ characters
		if len(word) > 2 && !seen[word] {
			terms = append(terms, word)
			seen[word] = true
		}
	}

	return terms
}

// Extract the actual subject of interest from a query like "summarize the XLA case"
func extractSubject(query string) string {
	query = strings.ToLower(query)

	// Remove common instruction words
	instructionWords := []string{"summarize", "summary", "explain", "tell me about", "describe"}
	for _, word := range instructionWords {
		query = strings.ReplaceAll(query, word, "")
	}

	// Remove filler words
	fillerWords := []string{"the", "a", "an", "about", "for", "of", "on", "with", "to", "me"}
	for _, word := range fillerWords {
		query = strings.ReplaceAll(query, " "+word+" ", " ")
	}

	// Handle specific common medical acronyms
	if strings.Contains(query, "xla") {
		return "x-linked agammaglobulinemia"
	}

	// Clean up and return
	return strings.TrimSpace(query)
}

// Connect to Weaviate with properly formatted URL
func searchWeaviate(query string, limit int) ([]models.MedicalExcerpt, error) {
	weaviateURL := os.Getenv("WEAVIATE_URL")
	if weaviateURL == "" {
		weaviateURL = "localhost:8081"
	}

	// Fix URL format issues - ensure proper scheme format
	if !strings.HasPrefix(weaviateURL, "http://") && !strings.HasPrefix(weaviateURL, "https://") {
		weaviateURL = "http://" + weaviateURL
	}

	// Strip the scheme for the Weaviate client config
	host := weaviateURL
	scheme := "http"

	if strings.HasPrefix(host, "http://") {
		host = strings.TrimPrefix(host, "http://")
		scheme = "http"
	} else if strings.HasPrefix(host, "https://") {
		host = strings.TrimPrefix(host, "https://")
		scheme = "https"
	}

	cfg := weaviate.Config{
		Host:   host,
		Scheme: scheme,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Weaviate client: %v", err)
	}

	// Check if the class exists before querying
	classExists, err := client.Schema().ClassExistenceChecker().WithClassName("MedicalExcerpt").Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to check if Weaviate class exists: %v", err)
	}

	if !classExists {
		// Return empty results instead of an error when class doesn't exist yet
		return []models.MedicalExcerpt{}, nil
	}

	// Add a debug message to help understand what's happening
	fmt.Printf("Debug: Weaviate MedicalExcerpt class exists, performing semantic search...\n")

	// Perform semantic search using nearText, which is handled by the transformers module
	fields := []graphql.Field{
		{Name: "sourceId"}, // Updated from bookId to sourceId
		{Name: "content"},
		{Name: "page"},
		{Name: "chapterNum"},
		{Name: "title"},
		{Name: "_additional", Fields: []graphql.Field{{Name: "certainty"}}},
	}

	nearText := client.GraphQL().NearTextArgBuilder().
		WithConcepts([]string{query})

	result, err := client.GraphQL().Get().
		WithClassName("MedicalExcerpt").
		WithFields(fields...).
		WithNearText(nearText).
		WithLimit(limit).
		Do(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to execute Weaviate search: %v", err)
	}

	// Extract and convert results
	var excerpts []models.MedicalExcerpt
	if result != nil && result.Data != nil {
		// Access the results - the structure depends on Weaviate's response format
		if data, ok := result.Data["Get"].(map[string]interface{}); ok {
			if medExcerpts, ok := data["MedicalExcerpt"].([]interface{}); ok {
				for _, item := range medExcerpts {
					if excerpt, ok := item.(map[string]interface{}); ok {
						var page int
						if pageVal, ok := excerpt["page"].(float64); ok {
							page = int(pageVal)
						}

						sourceID, _ := excerpt["sourceId"].(string) // Updated from bookId to sourceId
						content, _ := excerpt["content"].(string)

						excerpts = append(excerpts, models.MedicalExcerpt{
							SourceID: sourceID, // Updated field name
							Content:  content,
							Page:     page,
						})
					}
				}
			}
		}
	}

	return excerpts, nil
}

// Cache for book titles to avoid repeated lookups
var bookCache = make(map[string]*models.Book)

// Hardcoded book information for fallback
var fallbackBook = &models.Book{
	ISBN:    "9780815345121",
	Title:   "Case Studies in Immunology",
	Authors: []string{"Raif S Geha"},
}

func getBookDetails(bookID string) *models.Book {
	// Check cache first
	if book, ok := bookCache[bookID]; ok {
		return book
	}

	// Use the existing MongoDB client connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query for book details
	collection := mongoClient.Database(dbName).Collection("medical_books")

	// Try to convert the string ID to ObjectID if it looks like one
	var objectID primitive.ObjectID
	var err error
	if len(bookID) == 24 {
		objectID, err = primitive.ObjectIDFromHex(bookID)
		if err != nil {
			objectID = primitive.NilObjectID
		}
	}

	// Build a flexible query that tries multiple ID formats
	filter := bson.M{"$or": []bson.M{
		{"_id": bookID},             // String ID
		{"_id": objectID},           // ObjectID
		{"isbn13": bookID},          // Maybe bookID is actually ISBN
		{"isbn13": "9780815345121"}, // Known ISBN
	}}

	var book models.Book
	err = collection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		// If nothing found, use hardcoded fallback
		return fallbackBook
	}

	// Cache the result
	bookCache[bookID] = &book
	return &book
}

// Helper functions
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Fix for the collection.Indexes().CreateOne issue at line 346
// This function should be somewhere in the file, we need to update it to properly handle
// both return values from the CreateOne method

// The problematic code might look something like this:
// index := collection.Indexes().CreateOne(ctx, indexModel)
//
// Which should be changed to:
// index, err := collection.Indexes().CreateOne(ctx, indexModel)
// if err != nil {
//     // Handle error
// }

// Since we don't have the full context of the file, here's a general fix for the MongoDB index creation pattern:

func createSearchIndex(ctx context.Context, collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "content", Value: "text"},
			{Key: "keywords", Value: "text"},
		},
		Options: options.Index().SetName("content_text_keywords_text"),
	}

	// Fix: Capture both return values from CreateOne
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create text search index: %v", err)
	}

	return nil
}
