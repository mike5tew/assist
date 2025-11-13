package api

// import (
// 	"context"
// 	"encoding/json"
// 	"esp-organizer/internal/InfoFlow/InfoIn"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"runtime/debug"
// 	"time"
// )

// // This is a placeholder for your actual document processing service
// var semanticLinkService *InfoIn.SemanticLinkService

// func init() {
// 	// Initialize your semantic service here. This is a critical step.
// 	// You may need to pass database clients or other dependencies.
// 	// Example: semanticLinkService = InfoIn.NewSemanticLinkService(db.GetMongoClient(), db.GetWeaviateClient())
// 	// For now, we'll assume it can be created. If this fails, the service will be nil.
// 	var err error
// 	semanticLinkService, err = InfoIn.NewSemanticLinkService()
// 	if err != nil {
// 		log.Printf("FATAL: Failed to initialize SemanticLinkService: %v", err)
// 		semanticLinkService = nil // Ensure it's nil on failure
// 	}
// }

// // HandleImmunologyUpload is the handler for the chapter upload endpoint.
// // It now uses standard http interfaces to be compatible with both gin and gorilla/mux.
// func HandleImmunologyUpload(w http.ResponseWriter, r *http.Request) {
// 	if semanticLinkService == nil {
// 		http.Error(w, "Internal Server Error: Semantic processing service is not available.", http.StatusInternalServerError)
// 		return
// 	}

// 	// 1. Parse multipart form data from the request
// 	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
// 		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
// 		return
// 	}

// 	file, header, err := r.FormFile("chapter")
// 	if err != nil {
// 		http.Error(w, "Failed to get chapter file from form", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	fileBytes, err := io.ReadAll(file)
// 	if err != nil {
// 		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
// 		return
// 	}

// 	// Extract other form data
// 	chapterNumber := r.FormValue("chapter_number")
// 	chapterTitle := r.FormValue("chapter_title")
// 	sourceID := r.FormValue("source_id")
// 	// ... and so on

// 	jobID := fmt.Sprintf("chapter-%d", time.Now().Unix())

// 	// 2. Launch the background processing in a robust goroutine
// 	go func(ctx context.Context, jobID string, data []byte, filename string) {
// 		// IMPORTANT: Use recover() to catch panics within the goroutine.
// 		defer func() {
// 			if rec := recover(); rec != nil {
// 				log.Printf("FATAL: Goroutine panicked while processing job %s: %v\n%s", jobID, rec, debug.Stack())
// 			}
// 		}()

// 		log.Printf("✅ [Job %s] Starting background processing for chapter '%s'", jobID, filename)

// 		// --- THIS IS THE REAL WORK ---
// 		// Replace the placeholder with the actual call to your service.
// 		// This assumes your service has a method like ProcessDocument.
// 		// You will need to create the RawTextSubmission struct with the form data.
// 		submission := InfoIn.RawTextSubmission{
// 			Content:  string(data), // Or however you need to pass it
// 			Source:   "book",
// 			SourceID: sourceID,
// 			Metadata: map[string]interface{}{
// 				"filename":       filename,
// 				"chapter_number": chapterNumber,
// 				"chapter_title":  chapterTitle,
// 			},
// 		}

// 		// Call the actual processing pipeline
// 		err := semanticLinkService.ProcessAndStoreDocument(ctx, submission)
// 		if err != nil {
// 			log.Printf("❌ [Job %s] Error during semantic processing: %v", jobID, err)
// 			// Here you could update the job status in the DB to "FAILED"
// 			return
// 		}
// 		// --- END OF REAL WORK ---

// 		log.Printf("✅ [Job %s] Finished background processing successfully.", jobID)

// 	}(context.Background(), jobID, fileBytes, header.Filename) // Pass necessary data into the goroutine

// 	// 3. Immediately respond with 202 Accepted
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusAccepted)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message":  "Chapter upload accepted and is being processed in the background.",
// 		"batch_id": jobID,
// 		"status":   "PROCESSING",
// 	})
// }

// // You would need to register this handler in your routes setup.
// // For example, in a routes.go file:
// // router.POST("/api/immunology/upload-chapter", HandleImmunologyUpload)
