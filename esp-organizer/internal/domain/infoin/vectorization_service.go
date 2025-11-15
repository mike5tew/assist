package infoin

import (
	"context"
	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// VectorizationService handles the process of creating and storing vector embeddings.
type VectorizationService struct {
	LlmClient *llm.LlamaClient
}

// NewVectorizationService creates a new instance of the vectorization service.
func NewVectorizationService() *VectorizationService {
	return &VectorizationService{
		LlmClient: llm.NewLlamaClient(),
	}
}

// Example re-vectorization function
func ReVectorizeAllContent(ctx context.Context) error {
	// Connect to MongoDB
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		return fmt.Errorf("mongo connect: %w", err)
	}
	defer mongoDb.Client.Disconnect(ctx)

	// Initialize vectorization service
	vectorizer := NewVectorizationService()

	// Process each collection
	collections := []string{"skills", "case_studies", "medical_terms"}
	for _, collName := range collections {
		log.Printf("Re-vectorizing collection: %s", collName)

		// Get all documents from collection
		cursor, err := mongoDb.Database.Collection(collName).Find(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("finding documents in %s: %w", collName, err)
		}

		// Process each document
		for cursor.Next(ctx) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				log.Printf("Warning: Could not decode document: %v", err)
				continue
			}

			// Get text to vectorize (adjust field names as needed)
			var textToVectorize string
			if name, ok := doc["name"].(string); ok {
				textToVectorize += name + " "
			}
			if desc, ok := doc["description"].(string); ok {
				textToVectorize += desc
			}

			// Generate vector embedding
			vector, err := vectorizer.GenerateEmbedding(ctx, textToVectorize)
			if err != nil {
				log.Printf("Warning: Failed to generate embedding: %v", err)
				continue
			}

			// Store in Weaviate with proper property names
			weaviateData := map[string]interface{}{
				"mongo_id":     doc["_id"].(primitive.ObjectID).Hex(),
				"content_type": getContentTypeForCollection(collName),
				"subject":      getSubjectForDocument(doc),
			}

			// Store in Weaviate
			if err := db.StoreDocumentWithVector(ctx, getClassNameForCollection(collName), weaviateData, vector); err != nil {
				log.Printf("Warning: Failed to store in Weaviate: %v", err)
			}
		}
	}

	return nil
}

// Add these helper functions to your file:

// getContentTypeForCollection maps MongoDB collection names to content types
func getContentTypeForCollection(collName string) string {
	switch collName {
	case "skills":
		return "skill"
	case "case_studies":
		return "case_study"
	case "medical_terms":
		return "medical_term"
	case "chapters":
		return "chapter"
	case "subject_content":
		return "section"
	default:
		return "document"
	}
}

// getSubjectForDocument extracts the subject field from a document
func getSubjectForDocument(doc bson.M) string {
	// Try various field names that might contain subject information
	for _, field := range []string{"subject", "domain", "category"} {
		if subject, ok := doc[field].(string); ok && subject != "" {
			return subject
		}
	}

	// Default subject if none found
	return "general"
}

// getClassNameForCollection maps collection names to Weaviate class names
func getClassNameForCollection(collName string) string {
	switch collName {
	// All these are deprecated. All semantic searches should go to SemanticLinks.
	case "skills":
		log.Println("DEPRECATION: Attempting to get Weaviate class for 'skills'. This should be 'SemanticLinks'.")
		return "SemanticLinks" // Was "EducationalSkills"
	case "case_studies", "medical_terms", "subject_content", "chapters":
		log.Printf("DEPRECATION: Attempting to get Weaviate class for '%s'. This should be 'SemanticLinks'.", collName)
		return "SemanticLinks" // Was "SubjectAreaContent" or others
	default:
		return "SemanticLinks"
	}
}

// GenerateEmbedding creates a vector embedding for the provided text.
func (s *VectorizationService) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Delegate to the underlying LlamaClient
	return s.LlmClient.GenerateEmbedding(text)
}

// // VectorizeAndStoreSubjectContent processes extracted chapter data and stores it in Weaviate.
// func (s *VectorizationService) VectorizeAndStoreSubjectContent(ctx context.Context, data *models.ExtractedChapterData, metadata map[string]interface{}) error {
// 	className := "SubjectAreaContent"

// 	// Process Sections
// 	for _, section := range data.Sections {
// 		if len(section.Content) < 50 { // Skip very short sections
// 			continue
// 		}
// 		embeddingText := fmt.Sprintf("Section Title: %s\nContent: %s", section.Title, section.Content)
// 		properties := s.createProperties(metadata, "section", section.Content)
// 		if err := s.storeWithVector(ctx, className, properties, embeddingText); err != nil {
// 			log.Printf("Warning: Failed to vectorize section '%s': %v", section.Title, err)
// 		}
// 	}

// 	// Process Case Studies
// 	for _, cs := range data.CaseStudies {
// 		embeddingText := fmt.Sprintf("Case Study: %s\nPatient: %s, %s\nHistory: %s\nContent: %s",
// 			cs.CaseNumber, cs.PatientInfo.Age, cs.PatientInfo.Gender, cs.PatientInfo.MedicalHistory, cs.Content)
// 		properties := s.createProperties(metadata, "case_study", cs.Content)
// 		if err := s.storeWithVector(ctx, className, properties, embeddingText); err != nil {
// 			log.Printf("Warning: Failed to vectorize case study '%s': %v", cs.CaseNumber, err)
// 		}
// 	}

// 	// Process Medical Terms
// 	for _, term := range data.MedicalTerms {
// 		contextStr := ""
// 		if len(term.Context) > 0 {
// 			contextStr = term.Context[0]
// 		}
// 		embeddingText := fmt.Sprintf("Medical Term: %s\nCategory: %s\nContext: %s", term.Term, term.Category, contextStr)
// 		properties := s.createProperties(metadata, "medical_term", embeddingText)
// 		if err := s.storeWithVector(ctx, className, properties, embeddingText); err != nil {
// 			log.Printf("Warning: Failed to vectorize medical term '%s': %v", term.Term, err)
// 		}
// 	}

// 	log.Println("✅ Successfully processed and vectorized subject content for Weaviate.")
// 	return nil
// }

// createProperties builds the map of properties for a Weaviate object.
func (s *VectorizationService) createProperties(metadata map[string]interface{}, contentType, contentChunk string) map[string]interface{} {
	properties := make(map[string]interface{})
	for k, v := range metadata {
		// Only add property if the value is not an empty string
		if s, ok := v.(string); !ok || s != "" {
			properties[k] = v
		}
	}
	properties["content_type"] = contentType
	properties["content_chunk"] = contentChunk
	return properties
}

// storeWithVector generates an embedding and stores the document in Weaviate.
func (s *VectorizationService) storeWithVector(ctx context.Context, className string, properties map[string]interface{}, textToEmbed string) error {
	vector, err := s.LlmClient.GenerateEmbedding(textToEmbed)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	if err := db.StoreDocumentWithVector(ctx, className, properties, vector); err != nil {
		return fmt.Errorf("failed to store document with vector: %w", err)
	}
	return nil
}

/*
DEPRECATED: This function represents an old pipeline that directly vectorizes raw text,
which conflicts with the modern semantic relationship vectorization pipeline.
It is the likely source of vector dimension mismatches.
All content processing should now go through the SemanticLinkService.

func (v *VectorizationService) VectorizeAndStoreSubjectContent(ctx context.Context, data *models.ExtractedChapterData, metadata map[string]interface{}) error {
	// Get the subject from metadata or default to "Skills"
	subject, ok := metadata["subject"].(string)
	if !ok || subject == "" {
		subject = "Skills"
	}

	// Convert to proper class name (capitalize first letter, remove spaces)
	className := formatAsClassName(subject)

	// Store in Weaviate with the appropriate class
	return v.storeInWeaviate(ctx, data, className, metadata)
}
*/

// formatAsClassName converts a subject string to a valid Weaviate class name
func formatAsClassName(subject string) string {
	// Remove spaces and capitalize first letter of each word
	words := strings.Fields(subject)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

// storeInWeaviate processes and stores extracted data in Weaviate with the specified class name
func (s *VectorizationService) storeInWeaviate(ctx context.Context, data *models.ExtractedChapterData, className string, metadata map[string]interface{}) error {
	// Process raw text content if available
	if len(data.RawText) > 0 {
		// Split content into manageable chunks if it's large
		chunks := splitTextIntoChunks(data.RawText, 3000)

		for i, chunk := range chunks {
			if len(chunk) < 50 { // Skip very short chunks
				continue
			}

			// Copy metadata and add chunk-specific info
			chunkMeta := copyMetadata(metadata)
			chunkMeta["chunk_index"] = i
			chunkMeta["content_type"] = "full_text"

			// Create embedding text with context
			embeddingText := fmt.Sprintf("Document Content Chunk %d: %s", i, chunk)

			// Store in Weaviate using the existing function
			if err := s.storeWithVector(ctx, className, chunkMeta, embeddingText); err != nil {
				return fmt.Errorf("failed to vectorize text chunk %d: %w", i, err)
			}
		}
	}

	// Process Sections
	for _, section := range data.Sections {
		if len(section.Content) < 50 { // Skip very short sections
			continue
		}
		embeddingText := fmt.Sprintf("Section Title: %s\nContent: %s", section.Title, section.Content)
		properties := s.createProperties(metadata, "section", section.Content)
		if err := s.storeWithVector(ctx, className, properties, embeddingText); err != nil {
			return fmt.Errorf("failed to vectorize section '%s': %w", section.Title, err)
		}
	}

	// Process Case Studies
	for _, cs := range data.CaseStudies {
		embeddingText := fmt.Sprintf("Case Study: %s\nContent: %s", cs.CaseNumber, cs.Content)
		properties := s.createProperties(metadata, "case_study", cs.Content)
		if err := s.storeWithVector(ctx, className, properties, embeddingText); err != nil {
			return fmt.Errorf("failed to vectorize case study '%s': %w", cs.CaseNumber, err)
		}
	}

	// Process Medical Terms
	for _, term := range data.MedicalTerms {
		contextStr := ""
		if len(term.Context) > 0 {
			contextStr = term.Context[0]
		}
		embeddingText := fmt.Sprintf("Medical Term: %s\nCategory: %s\nContext: %s", term.Term, term.Category, contextStr)
		properties := s.createProperties(metadata, "medical_term", embeddingText)
		if err := s.storeWithVector(ctx, className, properties, embeddingText); err != nil {
			return fmt.Errorf("failed to vectorize medical term '%s': %w", term.Term, err)
		}
	}

	return nil
}

// Helper functions for text processing
func splitTextIntoChunks(text string, maxChunkSize int) []string {
	var chunks []string

	// If text is smaller than max size, return it as a single chunk
	if len(text) <= maxChunkSize {
		return []string{text}
	}

	// Split by paragraphs first
	paragraphs := strings.Split(text, "\n\n")
	currentChunk := ""

	for _, paragraph := range paragraphs {
		// If adding this paragraph would exceed max size and we already have content,
		// finish the current chunk and start a new one
		if len(currentChunk)+len(paragraph) > maxChunkSize && len(currentChunk) > 0 {
			chunks = append(chunks, currentChunk)
			currentChunk = paragraph
		} else {
			// Add separator if not the first content in chunk
			if len(currentChunk) > 0 {
				currentChunk += "\n\n"
			}
			currentChunk += paragraph
		}
	}

	// Add the last chunk if it has content
	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

func copyMetadata(original map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for k, v := range original {
		copy[k] = v
	}
	return copy
}
