package rawText

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/InfoFlow/InfoIn"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"esp-organizer/internal/config"
	"esp-organizer/internal/models"
	"esp-organizer/internal/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This package handles raw text submissions and concept submissions.
// It develops chunks by scanning for potenial keywords, checking these against the vector database,

// RawTextSubmission represents the expected JSON payload for raw text submission
type RawTextSubmission struct {
	Content  string            `json:"content"`
	Source   string            `json:"source"`    // e.g., "user_upload", "web_scrape_XYZ"
	SourceID string            `json:"source_id"` // e.g., filename, URL, or a user-defined ID
	Metadata map[string]string `json:"metadata"`  // Additional metadata from the frontend
}

// ConceptSubmission represents the payload for submitting a concept with influences
type ConceptSubmission struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Influences  []models.Influence `json:"influences"`
	Source      string             `json:"source"`             // e.g., "manual_entry_caffeine_study"
	SourceID    string             `json:"source_id"`          // e.g., a specific document ID or user input session
	Strength    string             `json:"strength,omitempty"` // Optional strength of the concept
}

// RawTextPOSTHandler now correctly uses the modern SemanticLinkService.
// This ensures that uploads from the frontend follow the same, correct data processing pipeline.
func RawTextPOSTHandler(w http.ResponseWriter, r *http.Request) {
	// Centralized config loading at the start of the application.
	projectRoot, err := utils.FindProjectRoot()
	if err != nil {
		log.Printf("Error finding project root: %v", err)
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}
	config.LoadConfig(projectRoot, ".env") // Ensure config is loaded for this handler.

	var submission RawTextSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure required Weaviate classes exist before processing.
	// This is critical now that auto-schema is disabled.
	ctx := r.Context()
	if err := db.CreateSemanticLinksClass(ctx); err != nil {
		http.Error(w, "Failed to prepare SemanticLinks class: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize the modern semantic link service
	semanticService, err := InfoIn.NewSemanticLinkService()
	if err != nil {
		http.Error(w, "Failed to create semantic service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a document model from the submission
	doc := models.SubjectContent{
		ID:          primitive.NewObjectID(),
		Content:     submission.Content,
		ContentType: "raw_text_upload",
		Domain:      submission.Metadata["domain"], // Assuming domain is passed in metadata
		Title:       submission.SourceID,
		Source: models.SourceReference{
			SourceID: submission.SourceID,
			Title:    submission.Source,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Default the domain to "immunology" if not provided (prevents empty filters downstream)
	if strings.TrimSpace(doc.Domain) == "" {
		doc.Domain = "immunology"
	}

	// Persist the document so background processors can update it later
	mdb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "Failed to init database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := mdb.Database.Collection("subject_content").InsertOne(r.Context(), doc); err != nil {
		http.Error(w, "Failed to save document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Process the document using the correct, modern pipeline in a goroutine
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		log.Printf("Starting semantic processing for uploaded document: %s", doc.Title)
		if err := semanticService.ProcessDocument(ctx, doc); err != nil {
			log.Printf("Error processing uploaded document %s: %v", doc.Title, err)
		} else {
			log.Printf("✅ Successfully processed uploaded document: %s", doc.Title)
		}
		// Best-effort backfill for legacy UI endpoints
		backfillLegacyImmunologyCollections(ctx, mdb, doc)
	}()

	// Respond immediately that the task has been accepted for processing
	response := map[string]string{
		"message":    "Document accepted for semantic processing.",
		"documentId": doc.ID.Hex(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

// backfillLegacyImmunologyCollections provides a compatibility layer for older UI components
// that query legacy MongoDB collections directly. It performs a best-effort update
// to collections like 'immunology_chapters'.
func backfillLegacyImmunologyCollections(ctx context.Context, mdb *db.MongoDB, doc models.SubjectContent) {
	// Only backfill for the immunology domain
	if doc.Domain != "immunology" {
		return
	}

	log.Printf("Performing legacy backfill for document: %s", doc.Title)

	// Create a placeholder entry in the 'immunology_chapters' collection.
	// This ensures that legacy endpoints that query this collection will see the new content.
	chaptersCollection := mdb.Database.Collection("immunology_chapters")
	chapterEntry := bson.M{
		"book_title":      doc.Source.Title,
		"chapter_title":   doc.Title,
		"source_document": doc.ID,
		"content_preview": doc.Content[:200] + "...", // Provide a short preview
		"created_at":      doc.CreatedAt,
		"updated_at":      doc.UpdatedAt,
	}

	_, err := chaptersCollection.InsertOne(ctx, chapterEntry)
	if err != nil {
		log.Printf("Warning: Failed to backfill legacy 'immunology_chapters' collection for doc %s: %v", doc.ID.Hex(), err)
	} else {
		log.Printf("✅ Successfully backfilled legacy chapter for doc %s", doc.ID.Hex())
	}

	// NOTE: A more advanced implementation could be added here to extract and backfill
	// 'case_studies' and 'medical_terms' if required by legacy endpoints.
	// For now, this placeholder approach resolves the immediate issue.
}

// SubmitConceptHandler handles the creation of a Concept and its associated Chunks and SemanticLinks
func SubmitConceptHandler(w http.ResponseWriter, r *http.Request) {
	var submission ConceptSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	//ctx := r.Context() // Use request context

	// 1. Create DataSource for this concept submission
	dataSource := models.DataSource{
		ID:          primitive.NewObjectID(),
		Type:        "concept_submission", // Or derive from submission.SourceType if added
		Source:      submission.SourceID,  // e.g., "manual_entry_caffeine_study"
		ProcessedAt: time.Now().Unix(),
	}
	// _, err := db.GetCollection("datasources").InsertOne(ctx, dataSource)
	// if err != nil {
	// 	http.Error(w, "Failed to create data source: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// For now, we simulate success

	// 2. Create the main Chunk for the Concept's description
	// This chunk will be vectorized and stored in Pinecone.
	vectSource := models.VectorSource{
		SourceID: dataSource.ID,
		// VectorID: "vector_id_from_pinecone_for_concept", // Placeholder, would be set after vectorization
	}
	conceptDescriptionChunk := models.Chunk{
		ID:           primitive.NewObjectID(),
		VectorSource: vectSource,             // Link to the DataSource
		Content:      submission.Description, // The concept's description becomes a chunk
		// VectorID: "vector_id_from_pinecone_for_description", // Placeholder
		Metadata:  map[string]any{"name": submission.Name, "type": "concept_description"},
		CreatedAt: time.Now().Unix(),
	}
	// _, err = db.GetCollection("chunks").InsertOne(ctx, conceptDescriptionChunk)
	// if err != nil {
	// 	http.Error(w, "Failed to create concept chunk: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// TODO: Vectorize conceptDescriptionChunk.Content and upsert to Pinecone, then update VectorID

	// 3. Create the Concept document in MongoDB

	concept := models.Concept{
		ID:             primitive.NewObjectID(), // This will be the source for semantic links
		Name:           submission.Name,
		Description:    submission.Description,
		DataSourceID:   dataSource.ID,
		PrimaryChunkID: conceptDescriptionChunk.ID,
		Influences:     submission.Influences, // Store raw influences for now, or process into SemanticLinks directly
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	}
	// _, err = db.GetCollection("concepts").InsertOne(ctx, concept)
	// if err != nil {
	// 	http.Error(w, "Failed to create concept: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// 4. Create SemanticLink documents based on Influences
	// Each Influence translates into a SemanticLink.
	// The TargetConceptID in the Influence refers to an *existing* Concept's ID.
	// The SourceChunkID for these links will be the conceptDescriptionChunk.ID created above.
	var createdSemanticLinks []models.SemanticLink
	for _, influence := range submission.Influences {
		// We need to find the PrimaryChunkID of the target concept.
		// This requires fetching the target concept from the DB.
		// targetConcept, err := db.GetConceptByID(ctx, influence.TargetConceptID)
		// if err != nil {
		//    log.Printf("Warning: Could not find target concept %s for influence: %v", influence.TargetConceptID.Hex(), err)
		//    continue // or handle error more robustly
		// }
		// if targetConcept.PrimaryChunkID.IsZero() {
		//    log.Printf("Warning: Target concept %s has no primary chunk ID", influence.TargetConceptID.Hex())
		//    continue
		// }

		semanticLink := models.SemanticLink{
			ID:           primitive.NewObjectID(),
			SourceID:     concept.ID,                      // The concept's ID is the source
			TargetID:     influence.TargetConceptID,       // This is the ID of the target
			SourceTerm:   concept.Name,                    // The term from the source concept
			TargetTerm:   influence.TargetConceptID.Hex(), // The term from the target concept
			RelationType: string(influence.RelationType),  // e.g., "influences", "related_to"
			Context:      influence.Context,               // Optional context for the link
			Confidence:   influence.Strength,              // Convert Strength to Confidence
			Vector:       nil,                             // Placeholder for vector, would be set after vectorization
			BatchID:      "",                              // Optional batch ID if this is part of a batch process
			Domain:       "concept_influence",             // Domain for this semantic link
			CreatedAt:    time.Now(),                      // Timestamp for when this link was created
			// Additional fields can be added as needed
		}

		// _, err = db.GetCollection("semantic_links").InsertOne(ctx, semanticLink)
		// if err != nil {
		// 	log.Printf("Failed to create semantic link for influence %+v: %v", influence, err)
		// 	// Decide on error handling: continue, or fail the whole request?
		// } else {
		createdSemanticLinks = append(createdSemanticLinks, semanticLink)
		// }
	}

	// Placeholder for actual database operations
	// Simulate saving these to MongoDB
	// db.SaveConcept(concept)
	// for _, sl := range createdSemanticLinks {
	// 	db.SaveSemanticLink(sl)
	// }
	// db.SaveChunk(conceptDescriptionChunk) // (and upsert to Pinecone)

	response := map[string]any{
		"message":           "Concept and influences submitted for processing.",
		"conceptId":         concept.ID.Hex(),
		"primaryChunkId":    conceptDescriptionChunk.ID.Hex(),
		"semanticLinkCount": len(createdSemanticLinks),
		// "createdSemanticLinkIds": getIDs(createdSemanticLinks), // Helper to get just IDs
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

// strengthToConfidence converts a string representation of strength to a float64 confidence score.
// func strengthToConfidence(strength string) float64 {
// 	switch strings.ToLower(strength) {
// 	case "high", "strong":
// 		return 0.9
// 	case "medium", "moderate":
// 		return 0.7
// 	case "low", "weak":
// 		return 0.4
// 	default:
// 		return 0.5 // Default confidence for unknown strength
// 	}
// }

// Helper to get IDs from a slice of structs having an ID field (if needed)
// func getIDs(links []models.SemanticLink) []string {
// 	ids := make([]string, len(links))
// 	for i, link := range links {
// 		ids[i] = link.ID.Hex()
// 	}
// 	return ids
// }

// TODO:
// 1. Implement actual database calls (db.CreateDataSource, db.CreateChunk, db.CreateConcept, db.CreateSemanticLink).
// 2. Implement vectorization logic (e.g., in an internal/vector package) and Pinecone upsert.
// 3. Implement robust error handling and logging.
// 4. In SubmitConceptHandler, implement the lookup for TargetConcept's PrimaryChunkID.
// 5. Add a new route in internal/api/routes.go (or main.go if routes are there) for "/concept" POST -> SubmitConceptHandler.
// 2. Implement vectorization logic (e.g., in an internal/vector package) and Pinecone upsert.
// 3. Implement robust error handling and logging.
// 4. In SubmitConceptHandler, implement the lookup for TargetConcept's PrimaryChunkID.
// 5. Add a new route in internal/api/routes.go (or main.go if routes are there) for "/concept" POST -> SubmitConceptHandler.
