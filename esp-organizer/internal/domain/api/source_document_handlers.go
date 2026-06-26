package api

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// upsertSourceDocument finds an existing SourceDocument by DOI or creates a new one.
// Returns (docID, displayTitle, error).
// displayTitle is the value suitable for denormalising onto SemanticLinks/ReviewTasks.
func upsertSourceDocument(title string, authors []string, doi string, year int, journal, originalFilename string) (primitive.ObjectID, string, error) {
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		return primitive.NilObjectID, "", err
	}
	defer mongoDb.Client.Disconnect(context.Background())

	coll := mongoDb.Client.Database("chisg_knowledge_base").Collection("source_documents")

	// --- Try to find by DOI first ---
	if doi != "" {
		var existing models.SourceDocument
		err := coll.FindOne(context.Background(), bson.M{"doi": doi}).Decode(&existing)
		if err == nil {
			// Already exists — return the canonical ID and stored title
			displayTitle := existing.Title
			if displayTitle == "" {
				displayTitle = existing.OriginalFilename
			}
			return existing.ID, displayTitle, nil
		}
		if err != mongo.ErrNoDocuments {
			return primitive.NilObjectID, "", err
		}
	}

	// --- Create new SourceDocument ---
	displayTitle := title
	if displayTitle == "" {
		displayTitle = originalFilename
	}

	doc := models.SourceDocument{
		ID:               primitive.NewObjectID(),
		Title:            title,
		Authors:          authors,
		DOI:              doi,
		Year:             year,
		Journal:          journal,
		OriginalFilename: originalFilename,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err = coll.InsertOne(context.Background(), doc)
	if err != nil {
		// Handle duplicate DOI race (two concurrent uploads of the same paper)
		if strings.Contains(err.Error(), "duplicate key") && doi != "" {
			var existing models.SourceDocument
			if findErr := coll.FindOne(context.Background(), bson.M{"doi": doi}).Decode(&existing); findErr == nil {
				dt := existing.Title
				if dt == "" {
					dt = existing.OriginalFilename
				}
				return existing.ID, dt, nil
			}
		}
		return primitive.NilObjectID, "", err
	}

	return doc.ID, displayTitle, nil
}

// IncrementCitationCounts scans all source_documents references and, for any
// reference whose DOI matches an existing document, increments that document's
// in_corpus_citation_count.  Should be called after references are extracted for
// a newly ingested document.
func IncrementCitationCounts(newDocID primitive.ObjectID, references []models.PaperReference) error {
	if len(references) == 0 {
		return nil
	}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		return err
	}
	defer mongoDb.Client.Disconnect(context.Background())

	coll := mongoDb.Client.Database("chisg_knowledge_base").Collection("source_documents")

	for i, ref := range references {
		if ref.DOI == "" {
			continue
		}
		// Find the referenced document in the corpus
		var cited models.SourceDocument
		err := coll.FindOne(context.Background(), bson.M{"doi": ref.DOI}).Decode(&cited)
		if err != nil {
			continue // Not in corpus yet — skip
		}

		// Increment its citation count
		coll.UpdateByID(context.Background(), cited.ID, bson.M{
			"$inc": bson.M{"in_corpus_citation_count": 1},
			"$set": bson.M{"updated_at": time.Now()},
		})

		// Back-fill resolved_id on the new document's reference entry
		coll.UpdateOne(
			context.Background(),
			bson.M{"_id": newDocID},
			bson.M{"$set": bson.M{
				fmt.Sprintf("references.%d.resolved_id", i): cited.ID,
				"updated_at": time.Now(),
			}},
		)
	}
	return nil
}

// ── HTTP handlers ───────────────────────────────────────────────────────────

// CreateSourceDocumentHandler allows manually creating a SourceDocument record
// without uploading a paper (e.g. to pre-register a paper before extraction).
func CreateSourceDocumentHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.SourceDocument
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error":"Invalid payload"}`, http.StatusBadRequest)
		return
	}

	payload.ID = primitive.NewObjectID()
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(context.Background())

	coll := mongoDb.Client.Database("chisg_knowledge_base").Collection("source_documents")
	if _, err := coll.InsertOne(context.Background(), payload); err != nil {
		http.Error(w, `{"error":"Failed to create document"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": payload.ID.Hex()})
}

// ListSourceDocumentsHandler returns all source documents, sorted by creation date desc.
func ListSourceDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(context.Background())

	coll := mongoDb.Client.Database("chisg_knowledge_base").Collection("source_documents")
	cursor, err := coll.Find(
		context.Background(),
		bson.M{},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(200),
	)
	if err != nil {
		http.Error(w, `{"error":"Failed to list documents"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var docs []models.SourceDocument
	if err := cursor.All(context.Background(), &docs); err != nil {
		http.Error(w, `{"error":"Failed to decode documents"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "documents": docs})
}

// GetSourceDocumentHandler returns a single SourceDocument by ID.
func GetSourceDocumentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, `{"error":"Invalid document ID"}`, http.StatusBadRequest)
		return
	}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(context.Background())

	coll := mongoDb.Client.Database("chisg_knowledge_base").Collection("source_documents")
	var doc models.SourceDocument
	if err := coll.FindOne(context.Background(), bson.M{"_id": docID}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"error":"Not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}
