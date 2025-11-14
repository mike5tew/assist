package main

import (
	"context"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"esp-organizer/internal/models"
	"flag"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Parse command line arguments
	sourceID := flag.String("source-id", "", "MongoDB ID of the source to process")
	testMode := flag.Bool("test-mode", false, "Run in test mode (simplified processing)")
	flag.Parse()

	if *sourceID == "" {
		log.Fatal("Error: --source-id is required")
	}

	// Connect to MongoDB - use your existing MongoDB connection code
	mongodb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Client.Disconnect(context.Background())
	// Define the collection name for immunology sources
	const sourcesCollectionName = "immunology_sources"
	sourcesCollection := mongodb.Database.Collection(sourcesCollectionName)

	sourceObjID, err := primitive.ObjectIDFromHex(*sourceID)
	if err != nil {
		log.Fatalf("Invalid source ID format: %v", err)
	}

	//Check the mongodb for the domain of immunology and the collection of immunology_sources
	// If not found, log an error and exit
	// If found, proceed with processing

	var source models.Source
	err = sourcesCollection.FindOne(context.Background(), bson.M{"_id": sourceObjID}).Decode(&source)
	if err != nil {
		log.Fatalf("Failed to find source with ID %s: %v", *sourceID, err)
	}

	log.Printf("📚 Processing source: %s", source.Title)

	if *testMode {
		log.Println("🧪 Running in test mode - bypassing normal processing")
		// Pass the mongodb object directly, not mongodb.Database
		createTestExcerpts(mongodb, source, *sourceID)
	} else {
		// Try to use the actual PDF processing code from your codebase
		log.Println("🔍 Starting real PDF processing using production code")

		// This would normally call your actual PDF processing implementation
		// For example: PDFProcessing.ProcessSourcePDF(source, mongodb)

		// Fallback if the actual implementation isn't available
		log.Println("⚠️ Real PDF processing not implemented or not accessible")
		log.Println("💡 Try running with --source-id for testing")
	}
}

// Create test excerpts for testing purposes using unified models
func createTestExcerpts(mongodb *db.MongoDB, source models.Source, sourceID string) {
	ctx := context.Background()

	// Use GetCollection method if it exists
	var excerptCollection *mongo.Collection
	if reflect.ValueOf(mongodb).MethodByName("GetCollection").IsValid() {
		// If GetCollection method exists on mongodb
		excerptCollection = mongodb.GetCollection("immunology_content")
	} else {
		// Fallback to direct collection access
		excerptCollection = mongodb.Database.Collection("immunology_content")
	}

	// Create a test excerpt using the unified MedicalExcerpt model
	excerpt := models.MedicalExcerpt{
		SourceID:  sourceID, // Updated to use SourceID instead of BookID
		Content:   "X-Linked Agammaglobulinemia is a primary immunodeficiency disease characterized by a lack of B-cell development.",
		Page:      42,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := excerptCollection.InsertOne(ctx, excerpt)
	if err != nil {
		log.Printf("⚠️ Failed to insert test excerpt: %v", err)
	} else {
		log.Println("✅ Created test excerpt")
	}

	// Add a second excerpt with different content
	excerpt2 := models.MedicalExcerpt{
		SourceID:  sourceID,
		Content:   "Immunodeficiency disorders impair the immune system's ability to defend the body against foreign cells.",
		Page:      43,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = excerptCollection.InsertOne(ctx, excerpt2)
	if err != nil {
		log.Printf("⚠️ Failed to insert second test excerpt: %v", err)
	} else {
		log.Println("✅ Created second test excerpt")
	}

	// Mark source as processed - use the same collection access pattern
	var sourcesCollection *mongo.Collection
	if reflect.ValueOf(mongodb).MethodByName("GetCollection").IsValid() {
		sourcesCollection = mongodb.GetCollection("medical_books")
	} else {
		sourcesCollection = mongodb.Database.Collection("medical_books")
	}

	update := bson.M{"$set": bson.M{"processed": true, "updatedAt": time.Now()}}
	objectID, _ := primitive.ObjectIDFromHex(source.ID)
	_, err = sourcesCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Printf("⚠️ Failed to update source as processed: %v", err)
	} else {
		log.Println("✅ Source marked as processed")
	}

	log.Println("✅ Test processing completed successfully")
}
