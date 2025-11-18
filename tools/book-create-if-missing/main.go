package main

import (
	"context"
	"esp-organizer/internal/models"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"esp-organizer/internal/domain/infoin/vectorize"
	"esp-organizer/internal/store/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Parse command line flags for book details
	isbn := flag.String("isbn", "", "ISBN-13 of the book")
	title := flag.String("title", "", "Title of the book")
	author := flag.String("author", "", "Author of the book")
	publisher := flag.String("publisher", "", "Publisher of the book")
	year := flag.String("year", "", "Year published")
	pdfPath := flag.String("pdf", "", "Path to the PDF file")
	testMode := flag.Bool("test", false, "Run in test mode (force processing)")
	forceReset := flag.Bool("force-reset", false, "Force clear existing data and recreate")
	// Add new flags for full chapter loading
	fullChapter := flag.Bool("full-chapter", false, "Load full chapter instead of excerpts")
	chapterNum := flag.Int("chapter", 1, "Chapter number to load (when using --full-chapter)")

	flag.Parse()

	// Validate required fields
	if *isbn == "" || *title == "" || *author == "" {
		log.Fatalf("Required fields missing. ISBN13, title, and author are required.")
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			mongoURI = "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin"
		}
	}

	dbName := os.Getenv("MONGODB_NAME")
	if dbName == "" {
		dbName = "esp_organizer" // Match the database name from logs
	}

	log.Printf("Connecting to MongoDB at %s", mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Check if the database is available
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}
	log.Printf("✅ Successfully connected to MongoDB")

	collection := client.Database(dbName).Collection("medical_books")

	// If force reset is enabled, remove existing book first
	if *forceReset {
		fmt.Println("🧹 Force reset enabled - removing existing book data...")

		// Clear medical content first using the clear-medical command if available
		if *testMode {
			fmt.Println("🧪 Test mode: Clearing all medical content...")
			clearMedicalContent()
		} else {
			// Just delete this specific book
			filter := bson.M{"isbn13": *isbn}
			deleteResult, err := collection.DeleteOne(ctx, filter)
			if err != nil {
				fmt.Printf("⚠️ Warning: Failed to delete existing book: %v\n", err)
			} else if deleteResult.DeletedCount > 0 {
				fmt.Println("✅ Existing book entry deleted")
			}
		}
	}

	// Check if book already exists
	var existingSource models.Source
	filter := bson.M{"isbn": *isbn}
	err = collection.FindOne(ctx, filter).Decode(&existingSource)

	if err == nil {
		fmt.Printf("✅ Book already exists in database with ID: %s\n", existingSource.ID)

		// IMPORTANT: Always update the PDF path if provided
		if *pdfPath != "" && *pdfPath != existingSource.FilePath {
			update := bson.M{"$set": bson.M{"filePath": *pdfPath}}
			_, err = collection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Printf("⚠️ Warning: Failed to update PDF path: %v\n", err)
			} else {
				fmt.Println("✅ Updated PDF path for existing book")
				existingSource.FilePath = *pdfPath
			}
		}

		// Check if PDF needs to be processed
		if existingSource.Processed && !*testMode {
			fmt.Println("📚 Book has already been processed")
		} else {
			if *testMode {
				fmt.Println("🧪 Test mode: Forcing reprocessing of PDF...")
			} else {
				fmt.Println("📝 Book exists but hasn't been processed. Starting processing...")
			}
			// Here you would call your PDF processing logic
			triggerPDFProcessing(existingSource, *pdfPath, testMode, *fullChapter, *chapterNum)
		}
		return
	}

	if err != mongo.ErrNoDocuments {
		log.Fatalf("Error checking for existing book: %v", err)
	}

	// Book doesn't exist, create a new one
	fmt.Println("📚 Book not found in database. Creating new entry...")

	// Validate PDF path
	if *pdfPath != "" {
		if _, err := os.Stat(*pdfPath); os.IsNotExist(err) {
			fmt.Printf("⚠️ Warning: PDF file not found at %s\n", *pdfPath)

			// Create the default PDF directory if it doesn't exist
			pdfDir := filepath.Dir(*pdfPath)
			if _, err := os.Stat(pdfDir); os.IsNotExist(err) {
				os.MkdirAll(pdfDir, 0755)
				fmt.Printf("📁 Created directory: %s\n", pdfDir)
			}

			// Create an empty placeholder file or test PDF
			if *testMode {
				// For test mode, create a minimal valid PDF
				createTestPDF(*pdfPath)
			} else {
				// Just create an empty file as placeholder
				emptyFile, err := os.Create(*pdfPath)
				if err != nil {
					fmt.Printf("⚠️ Could not create placeholder PDF: %v\n", err)
				} else {
					emptyFile.Close()
					fmt.Printf("📄 Created empty placeholder PDF at %s\n", *pdfPath)
				}
			}
		}
	}

	// Create new book using unified Source model
	newSource := models.Source{
		Type:      "book",
		ISBN:      *isbn,
		Title:     *title,
		Authors:   []string{*author},
		Publisher: *publisher,
		Year:      *year,
		FilePath:  *pdfPath,
		Processed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := collection.InsertOne(ctx, newSource)
	if err != nil {
		log.Fatalf("Failed to insert book: %v", err)
	}

	// Get the inserted ID and update the source
	insertedID := result.InsertedID.(primitive.ObjectID)
	newSource.ID = insertedID.Hex()
	fmt.Printf("✅ Book created with ID: %v\n", insertedID.Hex())

	// Print the actual collection name for debugging
	fmt.Printf("📊 Book stored in collection: %s.%s\n", dbName, "medical_books")

	// Process the PDF if path was provided
	if *pdfPath != "" {
		fmt.Println("📝 Starting PDF processing...")
		// Pass the full chapter flags to the processing function
		triggerPDFProcessing(newSource, *pdfPath, testMode, *fullChapter, *chapterNum)
	}

	// If in test mode, run validation after processing
	if *testMode {
		fmt.Println("🧪 Test mode: Running validation checks...")
		validateProcessing(newSource)
	}
}

// Update function signature to use models.Source instead of Book
func triggerPDFProcessing(source models.Source, pdfPath string, testMode *bool, fullChapter bool, chapterNum int) {
	if *testMode {
		fmt.Println("🧪 Test mode: Processing PDF...")

		if fullChapter {
			fmt.Printf("📄 Loading full chapter %d instead of excerpts...\n", chapterNum)
			createFullChapter(source, chapterNum)
		} else {
			// Original excerpt creation
			createTestExcerpts(source)
		}

		// Also update the collection that check_medical looks for
		updateCheckMedicalCollections(source)
	} else {
		fmt.Println("🔄 PDF processing would be triggered here")
		fmt.Println("📋 To process the PDF manually, run: make process-medical-pdf")
	}
}

// Update function signature to use models.Source
func createFullChapter(source models.Source, chapterNum int) {
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			mongoURI = "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin"
		}
	}

	dbName := os.Getenv("MONGODB_NAME")
	if dbName == "" {
		dbName = "esp_organizer"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("⚠️ Error connecting to MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)

	// Create full chapter directly
	excerptCollection := client.Database(dbName).Collection("medical_excerpts")

	// In a real implementation, this would extract the full chapter text from the PDF
	// For this example, we'll create a larger excerpt that simulates a full chapter

	// Sample chapter content (in a real implementation, this would be extracted from the PDF)
	chapterContent := `X-Linked Agammaglobulinemia (XLA) is a primary immunodeficiency disease characterized by a lack of B-cell development.

The disease was first described by Dr. Ogden Bruton in 1952 and is thus sometimes called Bruton's agammaglobulinemia. It's caused by mutations in the BTK gene located on the X chromosome, which encodes Bruton's tyrosine kinase.

Patients with XLA have a near-complete absence of B lymphocytes and immunoglobulins of all isotypes. This results in a marked reduction in mature B cells and plasma cells, causing a significant reduction in all classes of immunoglobulins.

Clinical features include:
1. Recurrent bacterial infections, particularly respiratory infections, starting after maternal antibody protection wanes
2. Recurrent sinusitis and pneumonia
3. Increased susceptibility to certain pathogens like Streptococcus pneumoniae and Haemophilus influenzae
4. Absence of lymphoid tissue such as tonsils and peripheral lymph nodes

Diagnosis is based on:
- Very low immunoglobulin levels
- Absent or very low B cells in peripheral blood
- Normal T cell numbers and function
- Genetic testing confirming BTK gene mutation

Treatment involves regular immunoglobulin replacement therapy and prophylactic antibiotics. With proper treatment, most patients can lead relatively normal lives, although they require lifelong medical management.

Immunodeficiency disorders impair the immune system's ability to defend the body against foreign cells. These disorders can be primary (present at birth) or secondary (acquired during life). X-Linked Agammaglobulinemia is one example of a primary immunodeficiency.`

	// Create a single document for the full chapter
	fullChapterDoc := struct {
		ID         primitive.ObjectID `bson:"_id,omitempty"`
		SourceID   string             `bson:"source_id"` // Updated to use source_id
		Content    string             `bson:"content"`
		Page       int                `bson:"page"`
		ChapterNum int                `bson:"chapterNum"`
		Title      string             `bson:"title"`
		Keywords   []string           `bson:"keywords,omitempty"`
		Processed  bool               `bson:"processed"`
	}{
		SourceID:   source.ID, // Use source ID instead of book ID
		Content:    chapterContent,
		Page:       chapterNum * 10, // Simulate start page
		ChapterNum: chapterNum,
		Title:      fmt.Sprintf("Chapter %d: Immunodeficiency Disorders", chapterNum),
		Keywords:   []string{"immunodeficiency", "X-Linked", "agammaglobulinemia", "BTK gene", "Bruton"},
		Processed:  true,
	}

	_, err = excerptCollection.InsertOne(ctx, fullChapterDoc)
	if err != nil {
		fmt.Printf("⚠️ Failed to insert full chapter: %v\n", err)
	} else {
		fmt.Printf("✅ Created full chapter %d with %d characters\n", chapterNum, len(chapterContent))

		// Now store in Weaviate for vector search
		storeExcerptInWeaviate(fullChapterDoc.SourceID, fullChapterDoc.Content, fullChapterDoc.Page, fullChapterDoc.ChapterNum, fullChapterDoc.Title, fullChapterDoc.Keywords)
	}

	// Mark source as processed
	sourcesCollection := client.Database(dbName).Collection("medical_books")
	objectID, _ := primitive.ObjectIDFromHex(source.ID)
	update := bson.M{"$set": bson.M{"processed": true, "updatedAt": time.Now()}}
	_, err = sourcesCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		fmt.Printf("⚠️ Failed to update source processed status: %v\n", err)
	} else {
		fmt.Println("✅ Marked source as processed")
	}

	fmt.Println("✅ Full chapter processing completed successfully")
}

// Update function signature to use models.Source
func createTestExcerpts(source models.Source) {
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			mongoURI = "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin"
		}
	}

	dbName := os.Getenv("MONGODB_NAME")
	if dbName == "" {
		dbName = "esp_organizer"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("⚠️ Error connecting to MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)

	// Create test excerpts directly
	excerptCollection := client.Database(dbName).Collection("medical_excerpts")

	// Create a test excerpt
	excerpt := struct {
		ID        primitive.ObjectID `bson:"_id,omitempty"`
		SourceID  string             `bson:"source_id"` // Updated to use source_id
		Content   string             `bson:"content"`
		Page      int                `bson:"page"`
		Keywords  []string           `bson:"keywords,omitempty"`
		Processed bool               `bson:"processed"`
	}{
		SourceID:  source.ID, // Use source ID instead of book ID
		Content:   "X-Linked Agammaglobulinemia is a primary immunodeficiency disease characterized by a lack of B-cell development.",
		Page:      42,
		Keywords:  []string{"immunodeficiency", "X-Linked", "B-cell"},
		Processed: true,
	}

	_, err = excerptCollection.InsertOne(ctx, excerpt)
	if err != nil {
		fmt.Printf("⚠️ Failed to insert test excerpt: %v\n", err)
	} else {
		fmt.Println("✅ Created test excerpt")

		// Now store in Weaviate for vector search
		storeExcerptInWeaviate(excerpt.SourceID, excerpt.Content, excerpt.Page, 1, "", excerpt.Keywords)
	}

	// Add a second excerpt with different content
	excerpt2 := struct {
		ID        primitive.ObjectID `bson:"_id,omitempty"`
		SourceID  string             `bson:"source_id"` // Updated to use source_id
		Content   string             `bson:"content"`
		Page      int                `bson:"page"`
		Keywords  []string           `bson:"keywords,omitempty"`
		Processed bool               `bson:"processed"`
	}{
		SourceID:  source.ID, // Use source ID instead of book ID
		Content:   "Immunodeficiency disorders impair the immune system's ability to defend the body against foreign cells.",
		Page:      43,
		Keywords:  []string{"immunodeficiency", "immune system", "disorder"},
		Processed: true,
	}

	_, err = excerptCollection.InsertOne(ctx, excerpt2)
	if err != nil {
		fmt.Printf("⚠️ Failed to insert second test excerpt: %v\n", err)
	} else {
		fmt.Println("✅ Created second test excerpt")

		// Now store in Weaviate for vector search
		storeExcerptInWeaviate(excerpt2.SourceID, excerpt2.Content, excerpt2.Page, 1, "", excerpt2.Keywords)
	}

	// Mark source as processed
	sourcesCollection := client.Database(dbName).Collection("medical_books")
	objectID, _ := primitive.ObjectIDFromHex(source.ID)
	update := bson.M{"$set": bson.M{"processed": true, "updatedAt": time.Now()}}
	_, err = sourcesCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		fmt.Printf("⚠️ Failed to update source processed status: %v\n", err)
	} else {
		fmt.Println("✅ Marked source as processed")
	}

	fmt.Println("✅ Test processing completed successfully")
}

// Helper function to store medical excerpts in Weaviate
func storeExcerptInWeaviate(sourceID string, content string, page int, chapterNum int, title string, keywords []string) {
	// Initialize Weaviate connection
	err := db.InitializeWeaviateFromEnv()
	if err != nil {
		fmt.Printf("⚠️ Failed to connect to Weaviate: %v\n", err)
		return
	}

	// Create MedicalExcerpt class if it doesn't exist
	err = db.CreateMedicalExcerptClass(context.Background())
	if err != nil {
		fmt.Printf("⚠️ Failed to create MedicalExcerpt class: %v\n", err)
		return
	}

	// Prepare properties for Weaviate
	properties := map[string]interface{}{
		"sourceId":   sourceID, // Updated to use sourceId instead of bookId
		"content":    content,
		"page":       page,
		"chapterNum": chapterNum,
	}

	// Add optional fields if present
	if title != "" {
		properties["title"] = title
	}

	if len(keywords) > 0 {
		properties["keywords"] = keywords
	}

	// Check if we should use client-side vectorization (default is server-side)
	useClientVectorization := os.Getenv("USE_CLIENT_VECTORIZATION") == "true"

	if useClientVectorization {
		// CLIENT-SIDE VECTORIZATION
		// Generate vector embedding using our own code
		vector, err := vectorize.GetEmbedding(content)
		if err != nil {
			fmt.Printf("⚠️ Failed to generate vector embedding: %v\n", err)
			// Fall back to server-side vectorization
			useClientVectorization = false
		} else {
			// Store with our generated vector
			err = db.StoreDocumentWithVector(context.Background(), "MedicalExcerpt", properties, vector)
			if err != nil {
				fmt.Printf("⚠️ Failed to store in Weaviate with vector: %v\n", err)
				// Fall back to server-side vectorization
				useClientVectorization = false
			} else {
				fmt.Println("✅ Successfully stored medical excerpt in Weaviate with client-side vector")
			}
		}
	}

	// If client-side vectorization was not requested or failed, use server-side
	if !useClientVectorization {
		// SERVER-SIDE VECTORIZATION
		// Let Weaviate generate the vector using its configured vectorizer
		err = db.StoreDocumentWithoutVector(context.Background(), "MedicalExcerpt", properties)
		if err != nil {
			fmt.Printf("⚠️ Failed to store in Weaviate: %v\n", err)
			return
		}
		fmt.Println("✅ Successfully stored medical excerpt in Weaviate using server-side vectorization")
	}
}

// Update function signature to use models.Source
func updateCheckMedicalCollections(source models.Source) {
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			mongoURI = "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin"
		}
	}

	dbName := os.Getenv("MONGODB_NAME")
	if dbName == "" {
		dbName = "esp_organizer"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("⚠️ Error connecting to MongoDB for collection update: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)

	// Create a sample entry in immunology_chapters collection
	immunologyChapters := client.Database(dbName).Collection("immunology_chapters")
	chapterDoc := bson.M{
		"title":     "Case 1: X-Linked Agammaglobulinemia",
		"sourceId":  source.ID, // Updated to use sourceId
		"isbn":      source.ISBN,
		"content":   "X-Linked Agammaglobulinemia is a primary immunodeficiency disease.",
		"page":      1,
		"chapter":   1,
		"processed": true,
	}

	_, err = immunologyChapters.InsertOne(ctx, chapterDoc)
	if err != nil {
		fmt.Printf("⚠️ Failed to insert into immunology_chapters: %v\n", err)
	} else {
		fmt.Println("✅ Created test entry in immunology_chapters")
	}
}

// ClearMedicalContent function remains the same
func clearMedicalContent() {
	// Try to run the clear-medical command
	cmd := exec.Command("go", "run", "cmd/clear-medical/main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("⚠️ Error clearing medical content: %v\n", err)
		fmt.Println(string(output))
	} else {
		fmt.Println("✅ Medical content cleared successfully")
	}
}

// Create a minimal valid PDF for testing
func createTestPDF(path string) {
	// This creates a minimal test PDF with content about immunology
	minimalPDF := `%PDF-1.3
1 0 obj
<< /Type /Catalog /Pages 2 0 R >>
endobj
2 0 obj
<< /Type /Pages /Kids [3 0 R] /Count 1 >>
endobj
3 0 obj
<< /Type /Page /Parent 2 0 R /Resources << /Font << /F1 4 0 R >> >> /MediaBox [0 0 612 792] /Contents 5 0 R >>
endobj
4 0 obj
<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>
endobj
5 0 obj
<< /Length 161 >>
stream
BT
/F1 24 Tf
100 700 Td
(Case Studies in Immunology) Tj
/F1 12 Tf
0 -50 Td
(X-Linked Agammaglobulinemia is a primary immunodeficiency disease.) Tj
ET
endstream
endobj
xref
0 6
0000000000 65535 f
0000000009 00000 n
0000000058 00000 n
0000000115 00000 n
0000000234 00000 n
0000000302 00000 n
trailer << /Size 6 /Root 1 0 R >>
startxref
514
%%EOF`

	err := os.WriteFile(path, []byte(minimalPDF), 0644)
	if err != nil {
		fmt.Printf("⚠️ Could not create test PDF: %v\n", err)
	} else {
		fmt.Printf("📄 Created test PDF with immunology content at %s\n", path)
	}
}

func validateProcessing(source models.Source) {
	fmt.Println("🔍 Validating source processing...")

	// Connect to MongoDB to check if book was marked as processed
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			mongoURI = "mongodb://admin:password123@localhost:27017/esp_organizer?authSource=admin"
		}
	}

	dbName := os.Getenv("MONGODB_NAME")
	if dbName == "" {
		dbName = "esp_organizer"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		fmt.Printf("⚠️ Validation failed - couldn't connect to MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)

	// Check if the source was marked as processed
	var updatedSource models.Source
	collection := client.Database(dbName).Collection("medical_books")
	objectID, _ := primitive.ObjectIDFromHex(source.ID)
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedSource)
	if err != nil {
		fmt.Printf("⚠️ Validation failed - couldn't find source: %v\n", err)
		return
	}

	if updatedSource.Processed {
		fmt.Println("✅ Validation successful - source was marked as processed")
	} else {
		fmt.Println("❌ Validation failed - source was not marked as processed")
	}

	// Check if any excerpts were created
	excerptCollection := client.Database(dbName).Collection("medical_excerpts")
	count, err := excerptCollection.CountDocuments(ctx, bson.M{"source_id": source.ID}) // Updated field name
	if err != nil {
		fmt.Printf("⚠️ Validation failed - couldn't check excerpts: %v\n", err)
		return
	}

	if count > 0 {
		fmt.Printf("✅ Validation successful - %d excerpts were created\n", count)
	} else {
		fmt.Println("❌ Validation failed - no excerpts were created")
	}

	// Run a simple search test if excerpts exist
	if count > 0 {
		fmt.Println("🔍 Running test search query...")
		cmd := exec.Command("go", "run", "cmd/search-medical/main.go", "--query", "immunodeficiency", "--limit", "3")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("⚠️ Search test failed: %v\n", err)
		} else {
			fmt.Println("✅ Search test results:")
			fmt.Println(string(output))
		}
	}
}
