package extraction

import (
	"bytes"
	"context"
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
	ID            string                `json:"id"`
	Status        string                `json:"status"`
	CreatedAt     string                `json:"created_at"`
	CompletedAt   string                `json:"completed_at,omitempty"`
	ExtractedData *models.ExtractedData `json:"extracted_data,omitempty"`
	Error         string                `json:"error,omitempty"`
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

func (tp *TextractProcessor) ProcessChapterFile(ctx context.Context, file io.Reader, chapterInfo models.ChapterInfo, bookSource models.Source, jobID string) (*models.ProcessingJob, error) {
	// Step 1: Upload to S3
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

	// Step 2: Start Textract LAYOUT analysis
	startInput := &textract.StartDocumentAnalysisInput{
		DocumentLocation: &types.DocumentLocation{
			S3Object: &types.S3Object{
				Bucket: aws.String(tp.bucketName),
				Name:   aws.String(s3Key),
			},
		},
		FeatureTypes: []types.FeatureType{
			types.FeatureTypeLayout,
			types.FeatureTypeTables,
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

	// Step 3: Wait for completion
	maxWaitTime := 8 * time.Minute
	pollInterval := 15 * time.Second
	startTime := time.Now()

	var allBlocks []types.Block

	for {
		if time.Since(startTime) > maxWaitTime {
			return nil, fmt.Errorf("LAYOUT analysis timed out after %v", maxWaitTime)
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

			// Handle pagination
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

			// Step 4: Extract and process text
			extractedText := extractOCRText(allBlocks)

			if len(extractedText) == 0 {
				log.Printf("⚠️ NO OCR TEXT EXTRACTED")
				return nil, fmt.Errorf("no text could be extracted from the document")
			}

			log.Printf("OCR text extracted: %d characters", len(extractedText))

			// Step 5: Process the extracted text using the provided source
			// Step 5: Process the extracted text using the provided source
			extractedData, err := ProcessExtractedText(ctx, extractedText, chapterInfo, bookSource, jobID)
			if err != nil {
				log.Printf("❌ Failed to process extracted text: %v", err)
				return nil, fmt.Errorf("failed to process extracted text: %v", err)
			}
			// Step 6: Clean up S3
			_, err = tp.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(tp.bucketName),
				Key:    aws.String(s3Key),
			})
			if err != nil {
				log.Printf("Warning: Failed to clean up S3 object: %v", err)
			}

			return &models.ProcessingJob{
				ID:            jobID,
				Status:        "completed",
				ChapterInfo:   chapterInfo,
				BookSource:    bookSource,
				CreatedAt:     time.Now(),
				CompletedAt:   aws.Time(time.Now()),
				ExtractedData: extractedData.ExtractedData,
			}, nil

		} else if status == types.JobStatusFailed {
			errorMsg := "unknown error"
			if statusResult.StatusMessage != nil {
				errorMsg = *statusResult.StatusMessage
			}
			log.Printf("❌ LAYOUT analysis failed: %s", errorMsg)
			return nil, fmt.Errorf("LAYOUT analysis failed: %s", errorMsg)
		}

		// Still processing, wait and check again
		time.Sleep(pollInterval)
	}
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

func abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
