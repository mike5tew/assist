package api

import (
	"bytes"
	"context"
	"encoding/json"
	"esp-organizer/internal/InfoFlow/InfoIn"
	"esp-organizer/internal/InfoFlow/InfoIn/extraction"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"esp-organizer/internal/models"

	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository interface for extraction job operations
type Repository interface {
	CreateExtractionJob(job *models.ExtractionJob) error
	UpdateExtractionJob(job *models.ExtractionJob) error
	GetExtractionJob(jobID string) (*models.ExtractionJob, error)
}

// Service handles extraction job operations
type Service struct {
	repo Repository
}

// NewService creates a new extraction service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Add this before starting the actual processing
func (s *Service) InitializeExtractionJob(jobID string) error {
	job := &models.ExtractionJob{
		ID:        jobID,
		Status:    "initializing",
		CreatedAt: time.Now(),
	}
	return s.repo.CreateExtractionJob(job)
}

// UpdateExtractionJobStatus updates the status of an extraction job
func (s *Service) UpdateExtractionJobStatus(jobID, status, message string, progress float64) error {
	job := &models.ExtractionJob{
		ID:        jobID,
		Status:    status,
		Message:   message,
		Progress:  progress,
		UpdatedAt: time.Now(),
	}
	return s.repo.UpdateExtractionJob(job)
}

// GetExtractionJob retrieves an extraction job by ID
func (s *Service) GetExtractionJob(jobID string) (*models.ExtractionJob, error) {
	return s.repo.GetExtractionJob(jobID)
}

// ExtractionStatusHandler checks the status of a Textract job
func ExtractionStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	jobID := vars["jobId"]

	if jobID == "" {
		// Try alternate parameter name
		jobID = vars["jobID"]
	}

	if jobID == "" {
		// If still empty, try to get from query string as fallback
		jobID = r.URL.Query().Get("jobId")
	}

	if jobID == "" {
		http.Error(w, "Missing job ID parameter", http.StatusBadRequest)
		return
	}

	// Check if this is a valid batch ID format
	if !strings.HasPrefix(jobID, "immunology-chapter-") && !strings.HasPrefix(jobID, "chapter-") && !strings.HasPrefix(jobID, "textract-") {
		log.Printf("Invalid job ID format: %s", jobID)
		http.Error(w, "Invalid job ID format", http.StatusBadRequest)
		return
	}

	// Connect to MongoDB
	mongOb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("Failed to connect to MongoDB for status check: %v", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mongOb.Client.Disconnect(r.Context())

	// Initialize the extraction repository and service
	repo := NewRepository(mongOb.Client)
	extractionService := NewService(repo)

	// Get job status from the extraction_jobs collection
	job, err := extractionService.GetExtractionJob(jobID)
	if err != nil {
		// If the job doesn't exist in our tracking system yet, check if documents exist
		filter := bson.M{"source.upload_id": jobID}
		count, countErr := mongOb.Database.Collection("subject_content").CountDocuments(r.Context(), filter)

		if countErr != nil {
			log.Printf("Error checking document count: %v", countErr)
			http.Error(w, "Failed to check document status", http.StatusInternalServerError)
			return
		}

		// If we found documents but no job record, the job completed but wasn't properly tracked
		if count > 0 {
			response := map[string]interface{}{
				"job_id":         jobID,
				"status":         "COMPLETE",
				"message":        "Document processing complete",
				"document_count": count,
				"timestamp":      time.Now(),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		// If we have neither a job record nor documents, return a "pending" status
		response := map[string]interface{}{
			"job_id":    jobID,
			"status":    "pending",
			"message":   "Job is being initialized",
			"timestamp": time.Now(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// We found the job record, return its status
	response := map[string]interface{}{
		"job_id":    job.ID,
		"status":    job.Status,
		"message":   job.Message,
		"progress":  job.Progress,
		"timestamp": time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// ImmunologyChapterUploadHandler handles uploads for immunology domain
func ImmunologyChapterUploadHandler(w http.ResponseWriter, r *http.Request) {
	// --- Synchronous Part: Quick Validation and Job Creation ---
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("chapter")
	if err != nil {
		http.Error(w, "Failed to get 'chapter' file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	// Copy form values
	formValues := make(map[string]string)
	for key, values := range r.PostForm {
		if len(values) > 0 {
			formValues[key] = values[0]
		}
	}

	batchID := fmt.Sprintf("immunology-chapter-%d", time.Now().Unix())

	// --- Asynchronous Part: Launch Background Processing ---
	go processImmunologyChapterInBackground(context.Background(), batchID, fileBytes, header.Filename, formValues)

	// --- Immediate Response ---
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Upload accepted. Processing has started in the background.",
		"batch_id": batchID,
		"domain":   "immunology",
		"status":   "PROCESSING",
	})

	// Check if this is a searchable PDF (from Adobe Scan)
	if isSearchablePDF(fileBytes) {
		log.Printf("[Batch %s] PDF is searchable (Adobe Scan OCR detected), extracting text locally", batchID)

		// Extract text locally (FREE, FAST)
		extractedText, err := extractTextFromSearchablePDF(fileBytes)
		if err == nil && len(extractedText) > 100 {
			log.Printf("[Batch %s] Local extraction successful: %d characters", batchID, len(extractedText))

			// Skip AWS Textract entirely!
			go processImmunologyChapterWithText(context.Background(), batchID, extractedText, header.Filename, formValues)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message":  "Searchable PDF detected. Processing locally (no AWS Textract needed).",
				"batch_id": batchID,
				"savings":  "$1.50 per 1000 pages",
			})
			return
		}
	}

	// Fallback to AWS Textract only for scanned (non-searchable) PDFs
	log.Printf("[Batch %s] PDF is not searchable, using AWS Textract", batchID)
	go processImmunologyChapterInBackground(context.Background(), batchID, fileBytes, header.Filename, formValues)

	// --- Immediate Response ---
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Upload accepted. Processing has started in the background.",
		"batch_id": batchID,
		"domain":   "immunology",
		"status":   "PROCESSING",
	})
}

// process pdf using pdfcpu
func processImmunologyChapterWithText(ctx context.Context, batchID, extractedText, filename string, formValues map[string]string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("FATAL PANIC in batch %s: %v\n%s", batchID, r, debug.Stack())
		}
	}()
	log.Printf("[Batch %s] Starting immunology processing for file: %s", batchID, filename)
	// --- 1. Initialize Services ---
	mongOb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[Batch %s] ERROR: Failed to connect to MongoDB: %v", batchID, err)
		return
	}
	defer mongOb.Client.Disconnect(ctx)
	// Initialize extraction repository and service for status tracking
	repo := NewRepository(mongOb.Client)
	extractionService := NewService(repo)
	// Initialize extraction job in database
	if err := extractionService.InitializeExtractionJob(batchID); err != nil {
		log.Printf("[Batch %s] ERROR: Failed to initialize extraction job: %v", batchID, err)
		// Continue anyway - we'll try to process without status tracking
	}
	// Domain-specific collections
	contentCollection := mongOb.Database.Collection("immunology_content")
	termsCollection := mongOb.Database.Collection("immunology_terms")
	// Initialize semantic link service for Weaviate
	semanticLinkService, err := InfoIn.NewSemanticLinkService()
	if err != nil {
		log.Printf("[Batch %s] ERROR: Failed to initialize SemanticLinkService: %v", batchID, err)
		extractionService.UpdateExtractionJobStatus(batchID, "failed", "Failed to initialize semantic link service", 0.0)
		return
	}
	// --- 2. Get Source Info ---
	var bookSource models.Source
	sourceID := formValues["source_id"]
	if sourceID != "" {
		objID, _ := primitive.ObjectIDFromHex(sourceID)
		_ = mongOb.Database.Collection("sources").FindOne(ctx, bson.M{"_id": objID}).Decode(&bookSource)
	}
	if bookSource.ID == "" {
		bookSource = models.Source{
			ID:        primitive.NewObjectID().Hex(),
			Title:     formValues["book_title"],
			ISBN:      formValues["book_isbn"],
			Type:      "medical_textbook",
			Domain:    "immunology",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if authors, ok := formValues["book_authors"]; ok {
			bookSource.Authors = strings.Split(authors, ",")
		}
		if publisher, ok := formValues["book_publisher"]; ok {
			bookSource.Publisher = publisher
		}
		if year, ok := formValues["book_year"]; ok {
			if _, err := strconv.Atoi(year); err == nil {
				bookSource.Year = year
			}
		}
	}
	chapterInfo := models.ChapterInfo{
		ChapterNumber: formValues["chapter_number"],
		ChapterTitle:  formValues["chapter_title"],
	}
	// --- 3. Process Extracted Text ---
	extractionService.UpdateExtractionJobStatus(batchID, "processing", "Processing extracted text", 0.2)

	job, err := extraction.ProcessExtractedText(ctx, extractedText, chapterInfo, bookSource, batchID)
	if err != nil {
		log.Printf("[Batch %s] ERROR: Text processing failed: %v", batchID, err)
		extractionService.UpdateExtractionJobStatus(batchID, "failed", fmt.Sprintf("Text processing failed: %v", err), 0.0)
		return
	}

	// Validate that we have extracted data
	if job.ExtractedData == nil {
		log.Printf("[Batch %s] ERROR: No extracted data returned from processing", batchID)
		extractionService.UpdateExtractionJobStatus(batchID, "failed", "No data extracted from text", 0.0)
		return
	}

	extractionService.UpdateExtractionJobStatus(batchID, "processing", "Text processing complete, storing content", 0.5)
	// --- 4. Store in Domain-Specific Collections ---
	// Enhanced validation
	hasContent := len(job.ExtractedData.RawText) > 0 ||

		len(job.ExtractedData.CaseStudies) > 0 ||
		len(job.ExtractedData.MedicalTerms) > 0
	if !hasContent {
		log.Printf("[Batch %s] ⚠️ No meaningful content extracted", batchID)
		return
	}
	// Store main chapter content in immunology_content
	chapterDoc := bson.M{
		"domain":         "immunology",
		"content_type":   "chapter",
		"title":          chapterInfo.ChapterTitle,
		"content":        job.ExtractedData.ChapterContent,
		"chapter_number": chapterInfo.ChapterNumber,
		"source": bson.M{
			"source_id":     bookSource.ID,
			"title":         bookSource.Title,
			"authors":       bookSource.Authors,
			"chapter_title": chapterInfo.ChapterTitle,
			"upload_id":     batchID,
		},
		"metadata": bson.M{
			"case_studies_count":  len(job.ExtractedData.CaseStudies),
			"medical_terms_count": len(job.ExtractedData.MedicalTerms),
			"raw_text_length":     len(job.ExtractedData.RawText),
		},
		"tags":       []string{"immunology", "chapter"},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	chapterResult, err := contentCollection.InsertOne(ctx, chapterDoc)
	if err != nil {
		log.Printf("[Batch %s] ERROR: Failed to store chapter content: %v", batchID, err)
		return
	}
	log.Printf("[Batch %s] Stored chapter in immunology_content with ID: %v", batchID, chapterResult.InsertedID)
	// Store case studies in immunology_content
	var storedDocuments []models.SubjectContent
	// Add chapter to documents for semantic linking
	if insertedID, ok := chapterResult.InsertedID.(primitive.ObjectID); ok {
		chapterContent := models.SubjectContent{
			ID:          insertedID,
			Domain:      "immunology",
			ContentType: "chapter",
			Title:       chapterInfo.ChapterTitle,
			Content:     job.ExtractedData.ChapterContent,
			Source: models.SourceReference{
				SourceID:     bookSource.ID,
				Title:        bookSource.Title,
				Authors:      bookSource.Authors,
				ChapterTitle: chapterInfo.ChapterTitle,
				UploadID:     batchID,
			},
			Tags:      []string{"immunology", "chapter"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		storedDocuments = append(storedDocuments, chapterContent)
	}
	for _, caseStudy := range job.ExtractedData.CaseStudies {
		caseStudyDoc := bson.M{
			"domain":       "immunology",
			"content_type": "case_study",
			"title":        caseStudy.Title,
			"content":      caseStudy.Content,
			"source": bson.M{
				"source_id":     bookSource.ID,
				"title":         bookSource.Title,
				"authors":       bookSource.Authors,
				"chapter_title": chapterInfo.ChapterTitle,
				"upload_id":     batchID,
			},
			"tags":       []string{"immunology", "case_study"},
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}
		caseStudyResult, err := contentCollection.InsertOne(ctx, caseStudyDoc)
		if err != nil {
			log.Printf("[Batch %s] ERROR: Failed to store case study: %v", batchID, err)
			continue
		}
		log.Printf("[Batch %s] Stored case study in immunology_content with ID: %v", batchID, caseStudyResult.InsertedID)
		// Add to documents for semantic linking
		if insertedID, ok := caseStudyResult.InsertedID.(primitive.ObjectID); ok {
			caseStudyContent := models.SubjectContent{
				ID:          insertedID,
				Domain:      "immunology",
				ContentType: "case_study",
				Title:       caseStudy.Title,
				Content:     caseStudy.Content,
				Source: models.SourceReference{
					SourceID:     bookSource.ID,
					Title:        bookSource.Title,
					Authors:      bookSource.Authors,
					ChapterTitle: chapterInfo.ChapterTitle,
					UploadID:     batchID,
				},
				Tags:      []string{"immunology", "case_study"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			storedDocuments = append(storedDocuments, caseStudyContent)
		}
	}
	// Store medical terms in immunology_terms
	for _, term := range job.ExtractedData.MedicalTerms {
		termDoc := bson.M{
			"domain":     "immunology",
			"term":       term.Term,
			"definition": term.Definition,
			"source": bson.M{
				"source_id":     bookSource.ID,
				"title":         bookSource.Title,
				"authors":       bookSource.Authors,
				"chapter_title": chapterInfo.ChapterTitle,
				"upload_id":     batchID,
			},
			"tags":       []string{"immunology", "medical_term"},
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}
		termResult, err := termsCollection.InsertOne(ctx, termDoc)
		if err != nil {
			log.Printf("[Batch %s] ERROR: Failed to store medical term: %v", batchID, err)
			continue
		}
		log.Printf("[Batch %s] Stored medical term in immunology_terms with ID: %v", batchID, termResult.InsertedID)
	}
	// --- 5. Create Semantic Links in Weaviate ---
	if len(storedDocuments) > 0 {
		err = semanticLinkService.CreateSemanticLinks(ctx, storedDocuments, "immunology", batchID)
		if err != nil {
			log.Printf("[Batch %s] ERROR: Failed to create semantic links: %v", batchID, err)
			extractionService.UpdateExtractionJobStatus(batchID, "failed", "Failed to create semantic links", 0.0)
			return
		}
		log.Printf("[Batch %s] Successfully created semantic links in Weaviate", batchID)
	}
	// --- 6. Finalize Extraction Job ---
	extractionService.UpdateExtractionJobStatus(batchID, "completed", "Processing complete", 1.0)
	log.Printf("[Batch %s] Immunology chapter processing completed successfully", batchID)
}

// processImmunologyChapterInBackground processes immunology content with Weaviate integration
func processImmunologyChapterInBackground(ctx context.Context, batchID string, fileBytes []byte, filename string, formValues map[string]string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("FATAL PANIC in batch %s: %v\n%s", batchID, r, debug.Stack())
		}
	}()

	log.Printf("[Batch %s] Starting immunology processing for file: %s", batchID, filename)

	// --- 1. Initialize Services ---
	mongOb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[Batch %s] ERROR: Failed to connect to MongoDB: %v", batchID, err)
		return
	}

	// Initialize extraction repository and service for status tracking
	repo := NewRepository(mongOb.Client)
	extractionService := NewService(repo)

	// Initialize extraction job in database
	if err := extractionService.InitializeExtractionJob(batchID); err != nil {
		log.Printf("[Batch %s] ERROR: Failed to initialize extraction job: %v", batchID, err)
		// Continue anyway - we'll try to process without status tracking
	}

	// Domain-specific collections
	contentCollection := mongOb.Database.Collection("immunology_content")
	termsCollection := mongOb.Database.Collection("immunology_terms")

	// Initialize semantic link service for Weaviate
	semanticLinkService, err := InfoIn.NewSemanticLinkService()
	if err != nil {
		log.Printf("[Batch %s] ERROR: Failed to initialize SemanticLinkService: %v", batchID, err)
		extractionService.UpdateExtractionJobStatus(batchID, "failed", "Failed to initialize semantic link service", 0.0)
		return
	}

	// --- 2. Get Source Info ---
	var bookSource models.Source
	sourceID := formValues["source_id"]
	if sourceID != "" {
		objID, _ := primitive.ObjectIDFromHex(sourceID)
		_ = mongOb.Database.Collection("sources").FindOne(ctx, bson.M{"_id": objID}).Decode(&bookSource)
	}
	if bookSource.ID == "" {
		bookSource = models.Source{
			ID:        primitive.NewObjectID().Hex(),
			Title:     formValues["book_title"],
			ISBN:      formValues["book_isbn"],
			Type:      "medical_textbook",
			Domain:    "immunology",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if authors, ok := formValues["book_authors"]; ok {
			bookSource.Authors = strings.Split(authors, ",")
		}
		if publisher, ok := formValues["book_publisher"]; ok {
			bookSource.Publisher = publisher
		}
		if year, ok := formValues["book_year"]; ok {
			if _, err := strconv.Atoi(year); err == nil {
				bookSource.Year = year
			}
		}
	}

	chapterInfo := models.ChapterInfo{
		ChapterNumber: formValues["chapter_number"],
		ChapterTitle:  formValues["chapter_title"],
	}

	// --- 3. Initialize and use the TextractProcessor directly ---
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "eu-west-2"
	}
	bucketName := os.Getenv("AWS_S3_BUCKET")
	if bucketName == "" {
		bucketName = "esp-new-organizer-immunology"
	}

	// Update status to processing
	extractionService.UpdateExtractionJobStatus(batchID, "processing", "Processing with AWS Textract", 0.1)

	processor, err := extraction.NewTextractProcessor(region, bucketName)
	if err != nil {
		log.Printf("[Batch %s] ERROR: Failed to initialize Textract processor: %v", batchID, err)
		extractionService.UpdateExtractionJobStatus(batchID, "failed", "Failed to initialize Textract processor", 0.0)
		return
	}

	// Process with Textract - pass the batchID directly instead of generating a new job ID
	fileReader := bytes.NewReader(fileBytes)
	extractionService.UpdateExtractionJobStatus(batchID, "processing", "Extracting text with AWS Textract", 0.2)

	job, err := processor.ProcessChapterFile(ctx, fileReader, chapterInfo, bookSource, batchID)
	if err != nil {
		log.Printf("[Batch %s] ERROR: Textract processing failed: %v", batchID, err)
		extractionService.UpdateExtractionJobStatus(batchID, "failed", fmt.Sprintf("Textract processing failed: %v", err), 0.0)
		return
	}

	if job.Status != "completed" || job.ExtractedData == nil {
		log.Printf("[Batch %s] WARNING: Textract job did not complete successfully", batchID)
		extractionService.UpdateExtractionJobStatus(batchID, "failed", "Textract extraction produced no usable data", 0.0)
		return
	}

	extractionService.UpdateExtractionJobStatus(batchID, "processing", "Text extraction complete, processing content", 0.5)

	// --- 4. Store in Domain-Specific Collections ---

	// Enhanced validation
	hasContent := len(job.ExtractedData.RawText) > 0 ||
		len(job.ExtractedData.CaseStudies) > 0 ||
		len(job.ExtractedData.MedicalTerms) > 0

	if !hasContent {
		log.Printf("[Batch %s] ⚠️ No meaningful content extracted", batchID)
		return
	}

	// Store main chapter content in immunology_content
	chapterDoc := bson.M{
		"domain":         "immunology",
		"content_type":   "chapter",
		"title":          chapterInfo.ChapterTitle,
		"content":        job.ExtractedData.ChapterContent,
		"chapter_number": chapterInfo.ChapterNumber,
		"source": bson.M{
			"source_id":     bookSource.ID,
			"title":         bookSource.Title,
			"authors":       bookSource.Authors,
			"chapter_title": chapterInfo.ChapterTitle,
			"upload_id":     batchID,
		},
		"metadata": bson.M{
			"case_studies_count":  len(job.ExtractedData.CaseStudies),
			"medical_terms_count": len(job.ExtractedData.MedicalTerms),
			"raw_text_length":     len(job.ExtractedData.RawText),
		},
		"tags":       []string{"immunology", "chapter"},
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	chapterResult, err := contentCollection.InsertOne(ctx, chapterDoc)
	if err != nil {
		log.Printf("[Batch %s] ERROR: Failed to store chapter content: %v", batchID, err)
		return
	}
	log.Printf("[Batch %s] Stored chapter in immunology_content with ID: %v", batchID, chapterResult.InsertedID)

	// Store case studies in immunology_content
	var storedDocuments []models.SubjectContent

	// Add chapter to documents for semantic linking
	if insertedID, ok := chapterResult.InsertedID.(primitive.ObjectID); ok {
		chapterContent := models.SubjectContent{
			ID:          insertedID,
			Domain:      "immunology",
			ContentType: "chapter",
			Title:       chapterInfo.ChapterTitle,
			Content:     job.ExtractedData.ChapterContent,
			Source: models.SourceReference{
				SourceID:      bookSource.ID,
				Title:         bookSource.Title,
				ChapterNumber: chapterInfo.ChapterNumber,
				ChapterTitle:  chapterInfo.ChapterTitle,
				UploadID:      batchID,
			},
			Tags: []string{"immunology", "chapter"},
			Metadata: map[string]interface{}{
				"case_studies_count":  len(job.ExtractedData.CaseStudies),
				"medical_terms_count": len(job.ExtractedData.MedicalTerms),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		storedDocuments = append(storedDocuments, chapterContent)
	}

	for _, caseStudy := range job.ExtractedData.CaseStudies {
		caseStudyDoc := bson.M{
			"domain":       "immunology",
			"content_type": "case_study",
			"title":        fmt.Sprintf("Case %s: %s", caseStudy.CaseNumber, caseStudy.Title),
			"content":      caseStudy.Content,
			"source": bson.M{
				"source_id":      bookSource.ID,
				"title":          bookSource.Title,
				"chapter_number": chapterInfo.ChapterNumber,
				"chapter_title":  chapterInfo.ChapterTitle,
				"upload_id":      batchID,
			},
			"metadata": bson.M{
				"case_number":       caseStudy.CaseNumber,
				"clinical_findings": caseStudy.ClinicalFindings,
				"parent_chapter_id": chapterResult.InsertedID,
			},
			"tags":       append([]string{"immunology", "case_study"}, caseStudy.Tags...),
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}

		caseResult, err := contentCollection.InsertOne(ctx, caseStudyDoc)
		if err != nil {
			log.Printf("[Batch %s] Warning: Failed to store case study: %v", batchID, err)
			continue
		}

		// Add case study to documents for semantic linking
		if caseID, ok := caseResult.InsertedID.(primitive.ObjectID); ok {
			caseStudyContent := models.SubjectContent{
				ID:          caseID,
				Domain:      "immunology",
				ContentType: "case_study",
				Title:       fmt.Sprintf("Case %s: %s", caseStudy.CaseNumber, caseStudy.Title),
				Content:     caseStudy.Content,
				Source: models.SourceReference{
					SourceID:      bookSource.ID,
					Title:         bookSource.Title,
					ChapterNumber: chapterInfo.ChapterNumber,
					ChapterTitle:  chapterInfo.ChapterTitle,
					UploadID:      batchID,
				},
				Tags: append([]string{"immunology", "case_study"}, caseStudy.Tags...),
				Metadata: map[string]interface{}{
					"case_number":       caseStudy.CaseNumber,
					"clinical_findings": caseStudy.ClinicalFindings,
					"parent_chapter_id": chapterResult.InsertedID,
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			storedDocuments = append(storedDocuments, caseStudyContent)
		}
	}

	// Store medical terms in immunology_terms
	for _, term := range job.ExtractedData.MedicalTerms {
		termDoc := bson.M{
			"domain":     "immunology",
			"term":       term.Term,
			"definition": term.Definition,
			"category":   term.Category,
			"source": bson.M{
				"source_id":      bookSource.ID,
				"title":          bookSource.Title,
				"chapter_number": chapterInfo.ChapterNumber,
				"upload_id":      batchID,
			},
			"tags":       append([]string{"immunology", "medical_term", term.Category}, term.Tags...),
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}

		termResult, err := termsCollection.InsertOne(ctx, termDoc)
		if err != nil {
			log.Printf("[Batch %s] Warning: Failed to store medical term: %v", batchID, err)
			continue
		}

		// Add medical term to documents for semantic linking
		if termID, ok := termResult.InsertedID.(primitive.ObjectID); ok {
			termContent := models.SubjectContent{
				ID:          termID,
				Domain:      "immunology",
				ContentType: "medical_term",
				Title:       term.Term,
				Content:     term.Definition,
				Source: models.SourceReference{
					SourceID:      bookSource.ID,
					Title:         bookSource.Title,
					ChapterNumber: chapterInfo.ChapterNumber,
					UploadID:      batchID,
				},
				Tags: append([]string{"immunology", "medical_term", term.Category}, term.Tags...),
				Metadata: map[string]interface{}{
					"term":     term.Term,
					"category": term.Category,
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			storedDocuments = append(storedDocuments, termContent)
		}
	}

	log.Printf("[Batch %s] ✅ Stored %d case studies and %d medical terms",
		batchID, len(job.ExtractedData.CaseStudies), len(job.ExtractedData.MedicalTerms))

	// --- 5. Semantic Linking to Weaviate ---
	log.Printf("[Batch %s] Starting semantic link extraction for %d documents...", batchID, len(storedDocuments))

	// Update status to completed BEFORE semantic linking (which may fail)
	// The core extraction and storage is done, so we can mark as completed
	extractionService.UpdateExtractionJobStatus(batchID, "completed", "Content extracted and stored successfully", 1.0)

	for _, document := range storedDocuments {
		if err := semanticLinkService.ProcessDocument(ctx, document); err != nil {
			log.Printf("[Batch %s] ERROR: Semantic link processing failed for %s: %v", batchID, document.Title, err)
		} else {
			log.Printf("[Batch %s] ✅ Semantic links created for: %s", batchID, document.Title)
		}
	}

	log.Printf("[Batch %s] Semantic link extraction complete.", batchID)

	// --- 6. Vectorization (if needed) ---
	// If you also need vector embeddings, add that logic here
	// vectorizer := InfoIn.NewVectorizationService()
	// ... vectorization logic ...

	log.Printf("✅ [Batch %s] Immunology processing finished successfully. %d documents processed for semantic links.", batchID, len(storedDocuments))
}

// getScannedPDFHelp provides user guidance for scanned PDF upload errors
func getScannedPDFHelp(details string) string {
	return "If you are uploading a scanned PDF and encounter this error, please check:\n" +
		"• The PDF is not encrypted or password-protected\n" +
		"• The file is not corrupted and opens in a PDF viewer\n" +
		"• The scan quality is high enough for OCR (text is readable)\n" +
		"• Try rescanning with higher contrast or resolution\n" +
		"• If the problem persists, contact support with this error: " + details
}

// DocumentUploadHandler handles document upload for Textract processing
func DocumentUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Handle file upload
	file, header, err := r.FormFile("document")
	if err != nil {
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	allowedTypes := map[string]bool{
		"application/pdf": true,
		"image/jpeg":      true,
		"image/png":       true,
		"image/tiff":      true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	// TODO: Upload to S3 and start Textract job
	// For now, return a placeholder response
	response := map[string]interface{}{
		"message":  "Document uploaded for processing",
		"filename": header.Filename,
		"size":     header.Size,
		"job_id":   "placeholder-job-id",
		"status":   "processing",
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func SearchByTagsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	tags := r.URL.Query()["tag"]
	if len(tags) == 0 {
		http.Error(w, "At least one 'tag' query parameter is required", http.StatusBadRequest)
		return
	}

	results, err := db.FindContentByTags(context.Background(), tags)
	if err != nil {
		http.Error(w, "Failed to retrieve content by tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// ImmunologySearchHandler handles POST /api/immunology/search requests
func ImmunologySearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Replace with real search logic
	results := []map[string]interface{}{
		{"title": "XLA Case Summary", "summary": "X-linked agammaglobulinemia (XLA) is a primary immunodeficiency..."},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
	})
}

// HSGSearchHandler handles POST /api/immunology/hsg-search requests for hierarchical semantic graph queries
func HSGSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Query      string                 `json:"query"`
		Domain     string                 `json:"domain"`  // Primary domain field
		Subject    string                 `json:"subject"` // Alternative name for domain
		MaxResults int                    `json:"max_results"`
		Metadata   map[string]interface{} `json:"metadata"` // Additional metadata
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set defaults and handle alternative field names
	if req.MaxResults == 0 {
		req.MaxResults = 5
	}

	// Determine the domain/subject from multiple possible sources
	domainValue := req.Domain
	if domainValue == "" {
		// Try alternative field name
		domainValue = req.Subject
	}
	if domainValue == "" && req.Metadata != nil {
		// Try to get from metadata
		if subj, ok := req.Metadata["subject"].(string); ok && subj != "" {
			domainValue = subj
		} else if dom, ok := req.Metadata["domain"].(string); ok && dom != "" {
			domainValue = dom
		}
	}
	if domainValue == "" {
		// Final fallback
		domainValue = "immunology"
	}

	log.Printf("HSG Search: Query=%s, Domain/Subject=%s", req.Query, domainValue)

	// Initialize HSG query service
	hsgService, err := InfoIn.NewHSGQueryService()
	if err != nil {
		log.Printf("Failed to initialize HSG query service: %v", err)
		http.Error(w, "Service initialization failed", http.StatusInternalServerError)
		return
	}

	// Perform HSG query with the determined domain/subject
	ragContext, answer, err := hsgService.QueryHSG(r.Context(), req.Query, domainValue, req.MaxResults)
	if err != nil {
		log.Printf("HSG query failed: %v", err)
		http.Error(w, "Query processing failed", http.StatusInternalServerError)
		return
	}

	// 🆕 Return BOTH normalized and original text
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"answer":      answer,
		"context":     ragContext,
		"query":       req.Query,
		"domain":      domainValue,
		"max_results": req.MaxResults,
		"metadata": map[string]interface{}{
			"processed_domain": domainValue,
			"timestamp":        time.Now(),
			"display_mode":     "preserve_symbols", // 🆕 Frontend can render original symbols
		},
	})
}

func AIChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "AI chat endpoint - not yet implemented",
		"status":  "placeholder",
	})
}

// Helper to detect if PDF has embedded text (OCR already done)
func isSearchablePDF(pdfBytes []byte) bool {
	// Check if PDF contains text layer
	// Simple heuristic: searchable PDFs contain "/Font" entries
	return bytes.Contains(pdfBytes, []byte("/Font"))
}

// Extract text from searchable PDF using local library
func extractTextFromSearchablePDF(pdfBytes []byte) (string, error) {
	// Use a Go PDF library (e.g., pdfcpu, unipdf)
	// For now, return placeholder
	return extractTextWithPDFCPU(pdfBytes)
}

// extractTextWithPDFCPU extracts text from PDF bytes
// TODO: Implement actual PDF text extraction using pdfcpu or similar library
func extractTextWithPDFCPU(pdfBytes []byte) (string, error) {
	// Placeholder implementation
	// You'll need to add a PDF library like:
	// - github.com/pdfcpu/pdfcpu
	// - github.com/ledongthuc/pdf
	// - github.com/unidoc/unipdf/v3

	// For now, return an error to indicate it's not implemented
	return "", fmt.Errorf("PDF text extraction not yet implemented - please add a PDF library like pdfcpu")

	// Example implementation with pdfcpu would look like:
	// import "github.com/pdfcpu/pdfcpu/pkg/api"
	// return api.ExtractText(bytes.NewReader(pdfBytes), nil)
}
