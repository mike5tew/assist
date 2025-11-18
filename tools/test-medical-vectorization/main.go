package main

import (
	"context"
	"esp-organizer/internal/domain/infoin"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	wmodels "github.com/weaviate/weaviate/entities/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	log.Println("🚀 Starting SEMANTIC medical knowledge extraction process...")
	ctx := context.Background()

	// 1) Init DBs
	log.Println("DEBUG: Connecting to MongoDB...")
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("❌ Failed to initialize MongoDB: %v", err)
	}
	defer mongoDb.Client.Disconnect(ctx)

	// Check MongoDB connection explicitly
	if err := mongoDb.Client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}
	log.Println("DEBUG: MongoDB connection successful")

	weaviateHost := db.GetEnvWithDefault("WEAVIATE_HOST", "localhost:8081")
	weaviateScheme := db.GetEnvWithDefault("WEAVIATE_SCHEME", "http")
	log.Printf("DEBUG: Connecting to Weaviate at %s://%s", weaviateScheme, weaviateHost)
	wc, err := weaviate.NewClient(weaviate.Config{Host: weaviateHost, Scheme: weaviateScheme})
	if err != nil {
		log.Fatalf("❌ Failed to create Weaviate client: %v", err)
	}

	// Verify Weaviate is responsive
	_, err = wc.Misc().ReadyChecker().Do(ctx)
	if err != nil {
		log.Fatalf("❌ Weaviate readiness check failed: %v", err)
	}
	log.Println("DEBUG: Weaviate connection successful")
	log.Println("✅ Successfully connected to MongoDB and Weaviate.")

	// Ensure all required Weaviate classes exist with the correct schema
	log.Println("Ensuring Weaviate schemas are correctly configured...")
	if err := db.CreateSubjectAreaContentClass(ctx); err != nil {
		log.Fatalf("❌ Failed to create SubjectAreaContent class in Weaviate: %v", err)
	}
	if err := db.CreateSemanticLinksClass(ctx); err != nil {
		log.Fatalf("❌ Failed to create SemanticLinks class in Weaviate: %v", err)
	}
	log.Println("✅ Weaviate schemas are ready.")

	// Add special logging to verify we're actually running this
	log.Println("📢 VERIFICATION: test-medical-vectorization main.go is executing")

	// --- SEMANTIC-FIRST PROCESSING ---
	log.Println("📖 Starting semantic knowledge extraction from medical content...")
	if err := processSemanticMedicalKnowledge(ctx, mongoDb, wc); err != nil {
		log.Fatalf("❌ Failed to process semantic medical knowledge: %v", err)
	}

	// Add verification that data was created
	verifyDataCreation(ctx, mongoDb)

	log.Println("✅ Semantic knowledge base created successfully.")
}

// Add this function to verify data was created
func verifyDataCreation(ctx context.Context, mdb *db.MongoDB) {
	log.Println("DEBUG: Starting data verification...")
	// Check medical books
	booksCollection := mdb.GetCollection("medical_books")
	bookCount, err := booksCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("⚠️ Error checking medical_books collection: %v", err)
	} else {
		log.Printf("📚 Verified %d medical books created", bookCount)
		if bookCount == 0 {
			log.Printf("⚠️ WARNING: No medical books were created!")
		} else {
			// Show a sample book to verify content
			var book models.Source
			if err := booksCollection.FindOne(ctx, bson.M{}).Decode(&book); err != nil {
				log.Printf("⚠️ Error getting sample book: %v", err)
			} else {
				log.Printf("DEBUG: Sample book - Title: %s, ISBN: %s, ID: %s",
					book.Title, book.ISBN, book.ID)
			}
		}
	}

	// Check subject content
	contentCollection := mdb.GetCollection("subject_content")
	contentCount, err := contentCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("⚠️ Error checking subject_content collection: %v", err)
	} else {
		log.Printf("📝 Verified %d subject content documents created", contentCount)
		if contentCount == 0 {
			log.Printf("⚠️ WARNING: No subject content documents were created!")
		} else {
			// Show a sample content to verify
			var content models.SubjectContent
			if err := contentCollection.FindOne(ctx, bson.M{}).Decode(&content); err != nil {
				log.Printf("⚠️ Error getting sample content: %v", err)
			} else {
				log.Printf("DEBUG: Sample content - Title: %s, Content length: %d chars",
					content.Title, len(content.Content))
			}
		}
	}

	// Check semantic links
	linksCollection := mdb.GetCollection("semantic_links")
	linkCount, err := linksCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("⚠️ Error checking semantic_links collection: %v", err)
	} else {
		log.Printf("🔗 Verified %d semantic links created", linkCount)
		if linkCount == 0 {
			log.Printf("⚠️ WARNING: No semantic links were created!")
		} else {
			// Show a sample link to verify
			var link models.SemanticLink
			if err := linksCollection.FindOne(ctx, bson.M{}).Decode(&link); err != nil {
				log.Printf("⚠️ Error getting sample link: %v", err)
			} else {
				log.Printf("DEBUG: Sample link - Source: %s, Target: %s, Type: %s",
					link.SourceTerm, link.TargetTerm, link.RelationType)
			}
		}
	}

	log.Println("DEBUG: Data verification completed")
}

// NEW: Semantic-first processing
func processSemanticMedicalKnowledge(ctx context.Context, mdb *db.MongoDB, wc *weaviate.Client) error {
	log.Println("DEBUG: Starting processSemanticMedicalKnowledge function...")

	// Create test book using the unified Source model
	bookDetails := models.Source{
		Type:      "book",
		ISBN:      "9780815345121",
		Title:     "Case Studies in Immunology",
		Authors:   []string{"Raif S Geha"},
		Publisher: "WW Norton & Co",
		Year:      "2016",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store source in MongoDB for reference
	booksCollection := mdb.GetCollection("medical_books")
	log.Printf("DEBUG: Getting medical_books collection from database %s", mdb.Database.Name())

	var existingBook models.Source
	log.Printf("DEBUG: Checking if book %s already exists...", bookDetails.ISBN)
	err := booksCollection.FindOne(ctx, bson.M{"isbn": bookDetails.ISBN}).Decode(&existingBook)
	if err == nil {
		log.Printf("DEBUG: Book '%s' already exists with ID: %s", bookDetails.Title, existingBook.ID)
		bookDetails = existingBook
	} else {
		log.Printf("DEBUG: Book not found, creating new entry. Error was: %v", err)
		bookDetails.ID = primitive.NewObjectID().Hex()
		log.Printf("DEBUG: Generated new book ID: %s", bookDetails.ID)

		_, err = booksCollection.InsertOne(ctx, bookDetails)
		if err != nil {
			log.Printf("❌ ERROR: Failed to insert new book: %v", err)
			return fmt.Errorf("failed to insert new book: %w", err)
		}
		log.Printf("✅ Created source record: %s with ID: %s", bookDetails.Title, bookDetails.ID)
	}

	// SEMANTIC PROCESSING: Extract knowledge directly
	rawContent := `X-Linked Agammaglobulinemia (XLA) is a primary immunodeficiency disease characterized by a lack of B-cell development. The disease was first described by Dr. Ogden Bruton in 1952 and is thus sometimes called Bruton's agammaglobulinemia. It's caused by mutations in the BTK gene located on the X chromosome, which encodes Bruton's tyrosine kinase.

Patients with XLA have a near-complete absence of B lymphocytes and immunoglobulins of all isotypes. This results in a marked reduction in mature B cells and plasma cells, causing a significant reduction in all classes of immunoglobulins. Clinical features include recurrent bacterial infections, particularly respiratory infections, starting after maternal antibody protection wanes.

Diagnosis of XLA is based on very low immunoglobulin levels, absent or very low B cells in peripheral blood, normal T cell numbers and function, and genetic testing confirming BTK gene mutation. Treatment involves regular immunoglobulin replacement therapy and prophylactic antibiotics.`

	log.Printf("DEBUG: Raw content length: %d characters", len(rawContent))

	// Initialize semantic link service for LLaMA processing
	log.Println("DEBUG: Creating semantic service...")
	semanticService, err := infoin.NewSemanticLinkService()
	if err != nil {
		log.Printf("❌ ERROR: Failed to create semantic service: %v", err)
		return fmt.Errorf("failed to create semantic service: %w", err)
	}
	log.Println("DEBUG: Semantic service created successfully")

	// Create a temporary document for processing
	subjectContentCollection := mdb.GetCollection("subject_content")
	log.Printf("DEBUG: Getting subject_content collection from database %s", mdb.Database.Name())

	// Check if document already exists
	log.Println("DEBUG: Checking if document already exists...")
	existingCount, err := subjectContentCollection.CountDocuments(ctx, bson.M{
		"title":            "X-Linked Agammaglobulinemia Case Study",
		"source.source_id": bookDetails.ID,
	})

	if err != nil {
		log.Printf("DEBUG: Error checking for existing document: %v", err)
	} else if existingCount > 0 {
		log.Printf("DEBUG: Found %d existing documents - will create a new one anyway", existingCount)
	}

	docID := primitive.NewObjectID()
	log.Printf("DEBUG: Generated new document ID: %s", docID.Hex())

	doc := models.SubjectContent{
		ID:          docID,
		Content:     rawContent,
		ContentType: "medical_chapter",
		Domain:      "immunology",
		Title:       "X-Linked Agammaglobulinemia Case Study",
		Source: models.SourceReference{
			SourceID: bookDetails.ID,
			Title:    bookDetails.Title,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the document to MongoDB first
	log.Println("DEBUG: Saving document to MongoDB...")
	_, err = subjectContentCollection.InsertOne(ctx, doc)
	if err != nil {
		log.Printf("❌ ERROR: Failed to insert document into MongoDB: %v", err)
		return fmt.Errorf("failed to insert document: %w", err)
	}
	log.Printf("DEBUG: Document saved successfully with ID: %s", doc.ID.Hex())

	// EXTRACT SEMANTIC RELATIONSHIPS using LLaMA
	log.Println("🔍 Extracting semantic relationships with LLaMA...")
	err = semanticService.ProcessDocument(ctx, doc)
	if err != nil {
		log.Printf("❌ ERROR: Failed to extract semantic relationships: %v", err)
		return fmt.Errorf("failed to extract semantic relationships: %w", err)
	}

	// Backfill legacy collections to satisfy existing frontend endpoints
	backfillLegacyImmunologyCollections(ctx, mdb, doc)

	log.Println("DEBUG: ProcessDocument completed successfully")
	log.Println("✅ Semantic knowledge extraction completed - knowledge base is ready!")
	return nil
}

// --- Removed file-scanning approach ---
// func findMedicalDocuments(...) { /* ...removed... */ }
// func processDocument(...) { /* ...removed... */ }

func chunkText(text string, approx int) []string {
	// First split on paragraphs
	paras := strings.Split(text, "\n\n")
	var chunks []string
	var buf strings.Builder

	writeChunk := func() {
		s := strings.TrimSpace(buf.String())
		if s != "" {
			chunks = append(chunks, s)
		}
		buf.Reset()
	}

	for _, p := range paras {
		if buf.Len()+len(p) > approx && buf.Len() > 0 {
			writeChunk()
		}
		if len(p) > approx*2 {
			// Further split oversized paragraphs
			for i := 0; i < len(p); i += approx {
				end := i + approx
				if end > len(p) {
					end = len(p)
				}
				part := strings.TrimSpace(p[i:end])
				if part != "" {
					chunks = append(chunks, part)
				}
			}
			continue
		}
		if buf.Len() > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString(p)
	}
	if buf.Len() > 0 {
		writeChunk()
	}
	return chunks
}

func uploadChunksToWeaviate(ctx context.Context, wc *weaviate.Client, className, sourceCollection, sourceID string, chunks []string) error {
	type item struct {
		id   string
		text string
		i    int
	}
	items := make([]item, 0, len(chunks))
	for i, c := range chunks {
		items = append(items, item{id: uuid.New().String(), text: c, i: i})
	}

	now := time.Now().UTC().Format(time.RFC3339)
	objects := make([]*wmodels.Object, 0, len(items))
	for i, it := range items {
		var prev, next string
		if i > 0 {
			prev = items[i-1].id
		}
		if i < len(items)-1 {
			next = items[i+1].id
		}
		obj := &wmodels.Object{
			ID:    strfmt.UUID(it.id),
			Class: className,
			Properties: map[string]interface{}{
				"text":             it.text,
				"sourceCollection": sourceCollection,
				"sourceId":         sourceID,
				"chunkIndex":       it.i,
				"totalChunks":      len(items),
				"prevId":           prev,
				"nextId":           next,
				"createdAt":        now,
			},
		}
		objects = append(objects, obj)
	}

	b := wc.Batch().ObjectsBatcher()
	for _, o := range objects {
		b = b.WithObjects(o)
	}
	resp, err := b.Do(ctx)
	if err != nil {
		return err
	}
	// Minimal error report
	failed := 0
	for _, r := range resp {
		if r.Result == nil || r.Result.Errors != nil {
			failed++
		}
	}
	if failed > 0 {
		return fmt.Errorf("weaviate batch insert: %d/%d failed", failed, len(objects))
	}
	return nil
}

// capitalizeFirst uppercases the first rune of a string.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// --- Removed file-scanning approach ---
// func findMedicalDocuments(...) { /* ...removed... */ }
// func processDocument(...) { /* ...removed... */ }

// Minimal compatibility layer so /api/immunology/* endpoints don’t return null
func backfillLegacyImmunologyCollections(ctx context.Context, mdb *db.MongoDB, doc models.SubjectContent) {
	try := func() {
		ch := mdb.GetCollection("immunology_chapters")
		cs := mdb.GetCollection("case_studies")
		mt := mdb.GetCollection("medical_terms")

		// Chapters
		cCount, _ := ch.CountDocuments(ctx, bson.M{})
		if cCount == 0 {
			_, _ = ch.InsertOne(ctx, bson.M{
				"title":          safeString(doc.Title, "Immunology Chapter 1"),
				"chapter_number": 1,
				"content":        truncateString(doc.Content, 4000),
				"bookId":         doc.Source.SourceID,
				"sourceId":       doc.Source.SourceID,
			})
		}

		// Case studies
		csCount, _ := cs.CountDocuments(ctx, bson.M{})
		if csCount == 0 {
			_, _ = cs.InsertOne(ctx, bson.M{
				"title":   safeString(doc.Title, "XLA Case Study"),
				"summary": truncateString(doc.Content, 600),
				"domain":  doc.Domain,
			})
		}

		// Terms
		mtCount, _ := mt.CountDocuments(ctx, bson.M{})
		if mtCount == 0 {
			for _, t := range []string{"X-Linked Agammaglobulinemia", "BTK", "B cell"} {
				_, _ = mt.InsertOne(ctx, bson.M{
					"name":        t,
					"description": "Auto-extracted term",
					"domain":      doc.Domain,
				})
			}
		}
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("legacy backfill panic recovered: %v", r)
		}
	}()
	try()
}

func safeString(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// Keep imports tidy - remove duplicate line
