package extraction

import (
	"bytes"
	"context"
	"esp-organizer/internal/config"
	"esp-organizer/internal/models"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// TextractExtractionJob represents a Textract extraction job
type TextractExtractionJob struct {
	ID            string                       `json:"id"`
	Status        string                       `json:"status"`
	CreatedAt     string                       `json:"created_at"`
	CompletedAt   string                       `json:"completed_at,omitempty"`
	ExtractedData *models.ExtractedChapterData `json:"extracted_data,omitempty"`
	Error         string                       `json:"error,omitempty"`
}

// TextractProcessor handles AWS Textract document processing
type TextractProcessor struct {
	region         string
	bucketName     string
	s3Client       *s3.Client
	textractClient *textract.Client
}

// NewTextractProcessor creates a new Textract processor
func NewTextractProcessor(region, bucketName string) (*TextractProcessor, error) {
	// Create AWS config using v2 SDK
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &TextractProcessor{
		region:         region,
		bucketName:     bucketName,
		s3Client:       s3.NewFromConfig(cfg),
		textractClient: textract.NewFromConfig(cfg),
	}, nil
}

// Add function to generate unique job IDs with proper format
func generateJobID() string {
	return fmt.Sprintf("textract-%d", time.Now().UnixNano())
}

// ProcessChapterFile processes a chapter file and extracts content using real AWS Textract OCR
func (tp *TextractProcessor) ProcessChapterFile(ctx context.Context, file io.Reader, chapterInfo models.ChapterInfo, bookSource models.Source, jobID string) (*TextractExtractionJob, error) {
	// Step 1: Upload to S3 (automated)
	s3Key := fmt.Sprintf("textract-input/%s.pdf", jobID)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	_, err = tp.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(tp.bucketName),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %v", err)
	}

	// Step 2: Start Textract LAYOUT analysis (automated)
	startInput := &textract.StartDocumentAnalysisInput{
		DocumentLocation: &types.DocumentLocation{
			S3Object: &types.S3Object{
				Bucket: aws.String(tp.bucketName),
				Name:   aws.String(s3Key),
			},
		},
		FeatureTypes: []types.FeatureType{
			types.FeatureTypeLayout, // Essential for proper text extraction from scanned images
			types.FeatureTypeTables, // Additional structure detection
		},
		JobTag: aws.String(fmt.Sprintf("scanned-book-layout-%s", jobID)),
	}

	result, err := tp.textractClient.StartDocumentAnalysis(ctx, startInput)
	if err != nil {
		log.Printf("❌ Failed to start LAYOUT analysis: %v", err)
		return nil, fmt.Errorf("failed to start LAYOUT analysis: %v", err)
	}

	textractJobID := *result.JobId
	log.Printf("⏳ LAYOUT analysis job started with ID: %s", textractJobID)

	// Step 3: Wait for completion (automated)
	maxWaitTime := 8 * time.Minute // Longer timeout for layout analysis
	pollInterval := 15 * time.Second
	startTime := time.Now()

	var allBlocks []types.Block

	for {
		if time.Since(startTime) > maxWaitTime {
			return nil, textractJobID, fmt.Errorf("LAYOUT analysis timed out after %v", maxWaitTime)
		}

		statusInput := &textract.GetDocumentAnalysisInput{
			JobId: aws.String(textractJobID),
		}

		statusResult, err := tp.textractClient.GetDocumentAnalysis(ctx, statusInput)
		if err != nil {
			log.Printf("Error checking LAYOUT analysis status: %v", err)
			time.Sleep(pollInterval)
			continue
		}

		status := statusResult.JobStatus
		log.Printf("📊 LAYOUT analysis status: %s (elapsed: %v)", status, time.Since(startTime))

		if status == types.JobStatusSucceeded {
			log.Printf("🎉 LAYOUT analysis completed successfully!")

			// Handle pagination with detailed logging
			allBlocks = statusResult.Blocks
			log.Printf("📄 Retrieved %d blocks from first page", len(allBlocks))

			nextToken := statusResult.NextToken
			pageCount := 1

			for nextToken != nil {
				pageCount++
				log.Printf("📄 Retrieving page %d with NextToken...", pageCount)

				nextPageInput := &textract.GetDocumentAnalysisInput{
					JobId:     aws.String(textractJobID),
					NextToken: nextToken,
				}

				nextPageResult, err := tp.textractClient.GetDocumentAnalysis(ctx, nextPageInput)
				if err != nil {
					log.Printf("❌ Error retrieving page %d: %v", pageCount, err)
					break
				}

				allBlocks = append(allBlocks, nextPageResult.Blocks...)
				log.Printf("📄 Retrieved %d blocks from page %d (total: %d)",
					len(nextPageResult.Blocks), pageCount, len(allBlocks))

				nextToken = nextPageResult.NextToken
			}

			log.Printf("✅ LAYOUT analysis completed: %d total blocks across %d pages",
				len(allBlocks), pageCount)

			// Analyze what types of blocks we got
			blockTypes := make(map[string]int)
			for _, block := range allBlocks {
				blockType := string(block.BlockType)
				blockTypes[blockType]++
			}
			log.Printf("📋 LAYOUT analysis block breakdown: %v", blockTypes)

			// Step 4: Extract structured data (automated)
			// Enhanced OCR text extraction with detailed analysis
			log.Printf("Analyzing %d Textract blocks for OCR content...", len(allBlocks))

			// Detailed block analysis
			blockTypesCount := make(map[string]int)
			hasText := false

			for _, block := range allBlocks {
				blockType := string(block.BlockType)
				blockTypesCount[blockType]++

				if blockType == "LINE" || blockType == "WORD" {
					if block.Text != nil && strings.TrimSpace(*block.Text) != "" {
						hasText = true
					}
				}
			}

			log.Printf("📊 Detailed block analysis:")
			for blockType, count := range blockTypesCount {
				log.Printf("  • %s: %d blocks", blockType, count)
			}
			log.Printf("📝 Contains readable text: %t", hasText)

			// Extract text using optimized method
			extractedText := extractOCRText(allBlocks)

			log.Printf("OCR text extracted: %d characters", len(extractedText))
			if len(extractedText) > 0 {
				log.Printf("OCR text preview: %s", config.TruncateString(extractedText, 200))
			} else {
				log.Printf("⚠️ NO OCR TEXT EXTRACTED - Possible issues:")
				log.Printf("  • PDF may be corrupted or unreadable")
				log.Printf("  • Image quality too poor for OCR")
				log.Printf("  • PDF contains only graphics/diagrams")
				log.Printf("  • S3 upload/Textract processing failed")

				if blockTypesCount["PAGE"] == 0 {
					log.Printf("  • Critical: No PAGE blocks found - Textract didn't process the PDF")
				}
			}

			// Step 5: Clean up S3 (automated)
			_, err = tp.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(tp.bucketName),
				Key:    aws.String(s3Key),
			})
			if err != nil {
				log.Printf("Warning: Failed to clean up S3 object: %v", err)
			}

			return &TextractExtractionJob{
				ID:            jobID,
				Status:        "completed",
				CreatedAt:     time.Now().Format(time.RFC3339),
				CompletedAt:   time.Now().Format(time.RFC3339),
				ExtractedData: nil, // Set to nil for automated workflow
			}, nil

		} else if status == types.JobStatusFailed {
			errorMsg := "unknown error"
			if statusResult.StatusMessage != nil {
				errorMsg = *statusResult.StatusMessage
			}
			log.Printf("❌ LAYOUT analysis failed: %s", errorMsg)
			return nil, textractJobID, fmt.Errorf("LAYOUT analysis failed: %s", errorMsg)
		}

		// Still processing, wait and check again
		time.Sleep(pollInterval)
	}
}

// tryLayoutAnalysis attempts asynchronous document analysis with LAYOUT features
func (tp *TextractProcessor) tryLayoutAnalysis(ctx context.Context, s3Key, jobID string) ([]types.Block, string, error) {
	log.Printf("🚀 Starting LAYOUT analysis for S3 key: %s", s3Key)

	startInput := &textract.StartDocumentAnalysisInput{
		DocumentLocation: &types.DocumentLocation{
			S3Object: &types.S3Object{
				Bucket: aws.String(tp.bucketName),
				Name:   aws.String(s3Key),
			},
		},
		FeatureTypes: []types.FeatureType{
			types.FeatureTypeLayout, // Essential for proper text extraction from scanned images
			types.FeatureTypeTables, // Additional structure detection
		},
		JobTag: aws.String(fmt.Sprintf("scanned-book-layout-%s", jobID)),
	}

	result, err := tp.textractClient.StartDocumentAnalysis(ctx, startInput)
	if err != nil {
		log.Printf("❌ Failed to start LAYOUT analysis: %v", err)
		return nil, "", fmt.Errorf("failed to start LAYOUT analysis: %v", err)
	}

	textractJobID := *result.JobId
	log.Printf("⏳ LAYOUT analysis job started with ID: %s", textractJobID)

	// Wait for completion with timeout
	maxWaitTime := 8 * time.Minute // Longer timeout for layout analysis
	pollInterval := 15 * time.Second
	startTime := time.Now()

	for {
		if time.Since(startTime) > maxWaitTime {
			return nil, textractJobID, fmt.Errorf("LAYOUT analysis timed out after %v", maxWaitTime)
		}

		statusInput := &textract.GetDocumentAnalysisInput{
			JobId: aws.String(textractJobID),
		}

		statusResult, err := tp.textractClient.GetDocumentAnalysis(ctx, statusInput)
		if err != nil {
			log.Printf("Error checking LAYOUT analysis status: %v", err)
			time.Sleep(pollInterval)
			continue
		}

		status := statusResult.JobStatus
		log.Printf("📊 LAYOUT analysis status: %s (elapsed: %v)", status, time.Since(startTime))

		if status == types.JobStatusSucceeded {
			log.Printf("🎉 LAYOUT analysis completed successfully!")

			// Handle pagination with detailed logging
			allBlocks := statusResult.Blocks
			log.Printf("📄 Retrieved %d blocks from first page", len(allBlocks))

			nextToken := statusResult.NextToken
			pageCount := 1

			for nextToken != nil {
				pageCount++
				log.Printf("📄 Retrieving page %d with NextToken...", pageCount)

				nextPageInput := &textract.GetDocumentAnalysisInput{
					JobId:     aws.String(textractJobID),
					NextToken: nextToken,
				}

				nextPageResult, err := tp.textractClient.GetDocumentAnalysis(ctx, nextPageInput)
				if err != nil {
					log.Printf("❌ Error retrieving page %d: %v", pageCount, err)
					break
				}

				allBlocks = append(allBlocks, nextPageResult.Blocks...)
				log.Printf("📄 Retrieved %d blocks from page %d (total: %d)",
					len(nextPageResult.Blocks), pageCount, len(allBlocks))

				nextToken = nextPageResult.NextToken
			}

			log.Printf("✅ LAYOUT analysis completed: %d total blocks across %d pages",
				len(allBlocks), pageCount)

			// Analyze what types of blocks we got
			blockTypes := make(map[string]int)
			for _, block := range allBlocks {
				blockType := string(block.BlockType)
				blockTypes[blockType]++
			}
			log.Printf("📋 LAYOUT analysis block breakdown: %v", blockTypes)

			return allBlocks, textractJobID, nil

		} else if status == types.JobStatusFailed {
			errorMsg := "unknown error"
			if statusResult.StatusMessage != nil {
				errorMsg = *statusResult.StatusMessage
			}
			log.Printf("❌ LAYOUT analysis failed: %s", errorMsg)
			return nil, textractJobID, fmt.Errorf("LAYOUT analysis failed: %s", errorMsg)
		}

		// Still processing, wait and check again
		time.Sleep(pollInterval)
	}
}

// tryDirectTextDetection attempts synchronous text detection (simpler, more reliable)
func (tp *TextractProcessor) tryDirectTextDetection(ctx context.Context, s3Key string) ([]types.Block, error) {
	log.Printf("🔍 Attempting synchronous TEXT detection for S3 key: %s", s3Key)

	detectInput := &textract.DetectDocumentTextInput{
		Document: &types.Document{
			S3Object: &types.S3Object{
				Bucket: aws.String(tp.bucketName),
				Name:   aws.String(s3Key),
			},
		},
	}

	detectResult, err := tp.textractClient.DetectDocumentText(ctx, detectInput)
		log.Printf("   • PDF may be corrupted or in an unsupported format")
		log.Printf("   • S3 object may not have uploaded correctly")
	}

	return detectResult.Blocks, nil
}

// extractOCRText is optimized for OCR content from scanned book pages
func extractOCRText(blocks []types.Block) string {
	var textBuilder strings.Builder

	// First try LINE blocks (preferred for OCR content)
	lineBlocks := []types.Block{}
	wordBlocks := []types.Block{}

	for _, block := range blocks {
		blockType := string(block.BlockType)
		if blockType == "LINE" && block.Text != nil {
			lineBlocks = append(lineBlocks, block)
		} else if blockType == "WORD" && block.Text != nil {
			wordBlocks = append(wordBlocks, block)
		}
	}

	log.Printf("OCR extraction: Found %d LINE blocks, %d WORD blocks", len(lineBlocks), len(wordBlocks))

	// Use LINE blocks if available (better for reading order)
	if len(lineBlocks) > 0 {
		for _, block := range lineBlocks {
			text := strings.TrimSpace(*block.Text)
			if text != "" {
				textBuilder.WriteString(text)
				textBuilder.WriteString("\n")
			}
		}
		log.Printf("Extracted text using LINE blocks: %d lines", len(lineBlocks))
	} else if len(wordBlocks) > 0 {
		// Fallback to WORD blocks, but try to reconstruct lines
		log.Printf("No LINE blocks found, reconstructing from WORD blocks...")
		textBuilder.WriteString(reconstructTextFromWords(wordBlocks))
	}

	return textBuilder.String()
}

// reconstructTextFromWords attempts to reconstruct readable text from WORD blocks
func reconstructTextFromWords(wordBlocks []types.Block) string {
	if len(wordBlocks) == 0 {
		return ""
	}

	// Sort words by vertical position (top to bottom) then horizontal (left to right)
	type WordPosition struct {
		Text string
		Top  float32
		Left float32
	}

	var words []WordPosition
	for _, block := range wordBlocks {
		word := WordPosition{Text: *block.Text}

		if block.Geometry != nil && block.Geometry.BoundingBox != nil {
			word.Top = block.Geometry.BoundingBox.Top
			word.Left = block.Geometry.BoundingBox.Left
		}
		words = append(words, word)
	}

	// Group words into lines by similar vertical positions
	const lineThreshold = 0.02 // Adjust based on testing
	var lines [][]WordPosition

	for _, word := range words {
		placed := false
		for i, line := range lines {
			if len(line) > 0 && abs32(line[0].Top-word.Top) < lineThreshold {
				lines[i] = append(lines[i], word)
				placed = true
				break
			}
		}
		if !placed {
			lines = append(lines, []WordPosition{word})
		}
	}

	// Sort each line by horizontal position and combine
	var textBuilder strings.Builder
	for _, line := range lines {
		// Sort words in line by left position
		for i := 0; i < len(line)-1; i++ {
			for j := i + 1; j < len(line); j++ {
				if line[i].Left > line[j].Left {
					line[i], line[j] = line[j], line[i]
				}
			}
		}

		// Combine words in line
		for i, word := range line {
			if i > 0 {
				textBuilder.WriteString(" ")
			}
			textBuilder.WriteString(word.Text)
		}
		textBuilder.WriteString("\n")
	}

	log.Printf("Reconstructed %d lines from %d words", len(lines), len(words))
	return textBuilder.String()
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

// extractCaseStudies finds case study patterns in the extracted text
func extractCaseStudies(text string, chapterInfo models.ChapterInfo) []models.CaseStudy {
	var caseStudies []models.CaseStudy

	// If no text was extracted, don't create fallback case studies
	if len(strings.TrimSpace(text)) == 0 {
		log.Printf("⚠️ No text available - skipping case study extraction")
		return caseStudies
	}

	// Look for case study patterns (case-insensitive)
	lowerText := strings.ToLower(text)

	// Common case study indicators
	caseIndicators := []string{
		"case study",
		"clinical case",
		"patient:",
		"case:",
		"scenario:",
	}

	caseNumber := 1
	for _, indicator := range caseIndicators {
		if strings.Contains(lowerText, indicator) {
			// Found a potential case study
			// Extract surrounding text (simplified extraction)
			startIdx := strings.Index(lowerText, indicator)
			if startIdx != -1 {
				endIdx := startIdx + 500 // Extract 500 chars
				if endIdx > len(text) {
					endIdx = len(text)
				}

				content := text[startIdx:endIdx]

				caseStudy := models.CaseStudy{
					CaseNumber: fmt.Sprintf("%s-%d", chapterInfo.ChapterNumber, caseNumber),
					Content:    strings.TrimSpace(content),
					PatientInfo: models.PatientInfo{
						Age:            extractAge(content),
						Gender:         extractGender(content),
						MedicalHistory: fmt.Sprintf("From %s", chapterInfo.ChapterTitle),
					},
					ClinicalFindings: extractFindings(content),
					DiscussionPoints: extractDiscussionPoints(content),
				}

				caseStudies = append(caseStudies, caseStudy)
				caseNumber++

				// For now, limit to first few case studies to avoid duplicates
				if len(caseStudies) >= 3 {
					break
				}
			}
		}
	}

	// Remove the fallback case study creation when no text is found
	// Only create case studies when we actually have content to work with
	if len(caseStudies) == 0 {
		log.Printf("No case study patterns found in extracted text")
	}

	return caseStudies
}

// extractMedicalTerms finds medical terminology in the extracted text
func extractMedicalTerms(text string, chapterInfo models.ChapterInfo) []models.MedicalTerm {
	var medicalTerms []models.MedicalTerm

	// If no text was extracted, don't create fallback medical terms
	if len(strings.TrimSpace(text)) == 0 {
		log.Printf("⚠️ No text available - skipping medical term extraction")
		return medicalTerms
	}

	// Common immunology terms to look for
	immunologyTerms := []string{
		"immunoglobulin", "antibody", "antigen", "lymphocyte", "cytokine",
		"complement", "macrophage", "neutrophil", "eosinophil", "basophil",
		"plasma cell", "memory cell", "helper t cell", "cytotoxic t cell",
		"natural killer", "dendritic cell", "mhc", "hla", "immunodeficiency",
		"autoimmune", "hypersensitivity", "allergy", "inflammation",
		"interferon", "interleukin", "tumor necrosis factor", "chemokine",
	}

	lowerText := strings.ToLower(text)

	for _, term := range immunologyTerms {
		count := strings.Count(lowerText, term)
		if count > 0 {
			// Extract context around the term
			context := extractTermContext(text, term)

			medicalTerm := models.MedicalTerm{
				Term:      strings.Title(term),
				Category:  categorizeImmunologyTerm(term),
				Frequency: count,
				Context:   append(context, fmt.Sprintf("Found in %s", chapterInfo.ChapterTitle)),
			}

			medicalTerms = append(medicalTerms, medicalTerm)
		}
	}

	// Remove the fallback medical term creation when no text is found
	// Only create terms when we actually have content to work with
	if len(medicalTerms) == 0 {
		log.Printf("No medical terms found in extracted text")
	}

	return medicalTerms
}

// Helper functions for text extraction
func extractAge(text string) string {
	// Simple age extraction patterns
	agePatterns := []string{"age", "year", "yo", "y/o"}
	lowerText := strings.ToLower(text)

	for _, pattern := range agePatterns {
		if idx := strings.Index(lowerText, pattern); idx != -1 {
			// Look for numbers before or after the pattern
			start := max(0, idx-10)
			end := min(len(text), idx+20)
			segment := text[start:end]

			// Extract any numbers from this segment
			for i, char := range segment {
				if char >= '0' && char <= '9' {
					// Found a number, try to extract it
					numStart := i
					numEnd := i
					for numEnd < len(segment) && segment[numEnd] >= '0' && segment[numEnd] <= '9' {
						numEnd++
					}
					return segment[numStart:numEnd]
				}
			}
		}
	}
	return "Not specified"
}

func extractGender(text string) string {
	lowerText := strings.ToLower(text)
	if strings.Contains(lowerText, "female") || strings.Contains(lowerText, "woman") || strings.Contains(lowerText, "she") {
		return "Female"
	}
	if strings.Contains(lowerText, "male") || strings.Contains(lowerText, "man") || strings.Contains(lowerText, "he") {
		return "Male"
	}
	return "Not specified"
}

func extractFindings(text string) []string {
	// Simple extraction of clinical findings
	findings := []string{}

	findingKeywords := []string{"findings", "symptoms", "presents with", "shows", "demonstrates"}
	lowerText := strings.ToLower(text)

	for _, keyword := range findingKeywords {
		if strings.Contains(lowerText, keyword) {
			findings = append(findings, fmt.Sprintf("Clinical finding related to %s", keyword))
		}
	}

	if len(findings) == 0 {
		findings = []string{"Extracted from case study content"}
	}

	return findings
}

func extractDiscussionPoints(text string) []string {
	// Simple extraction of discussion points
	points := []string{}

	discussionKeywords := []string{"discussion", "analysis", "mechanism", "pathogenesis", "treatment"}
	lowerText := strings.ToLower(text)

	for _, keyword := range discussionKeywords {
		if strings.Contains(lowerText, keyword) {
			points = append(points, fmt.Sprintf("Discussion point regarding %s", keyword))
		}
	}

	if len(points) == 0 {
		points = []string{"Refer to case study for detailed discussion"}
	}

	return points
}

func extractTermContext(text, term string) []string {
	context := []string{}
	lowerText := strings.ToLower(text)
	lowerTerm := strings.ToLower(term)

	// Find occurrences and extract surrounding words
	words := strings.Fields(lowerText)
	for i, word := range words {
		if strings.Contains(word, lowerTerm) {
			// Add surrounding words as context
			start := max(0, i-2)
			end := min(len(words), i+3)

			contextWords := words[start:end]
			contextStr := strings.Join(contextWords, " ")
			context = append(context, contextStr)

			// Limit context entries
			if len(context) >= 3 {
				break
			}
		}
	}

	if len(context) == 0 {
		context = []string{term}
	}

	return context
}

func categorizeImmunologyTerm(term string) string {
	lowerTerm := strings.ToLower(term)

	if strings.Contains(lowerTerm, "cell") {
		return "Cell Type"
	}
	if strings.Contains(lowerTerm, "antibody") || strings.Contains(lowerTerm, "immunoglobulin") {
		return "Antibody"
	}
	if strings.Contains(lowerTerm, "cytokine") || strings.Contains(lowerTerm, "interleukin") || strings.Contains(lowerTerm, "interferon") {
		return "Cytokine"
	}
	if strings.Contains(lowerTerm, "antigen") {
		return "Antigen"
	}

	return "Immunology Concept"
}
