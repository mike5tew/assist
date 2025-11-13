package extraction

import (
	"context"
	"esp-organizer/internal/models"
	"strings"
	"time"
)

// ProcessExtractedText processes already-extracted text (e.g., from searchable PDFs)
// without using AWS Textract
func ProcessExtractedText(
	ctx context.Context,
	extractedText string,
	chapterInfo models.ChapterInfo,
	bookSource models.Source,
	batchID string,
) (*models.ExtractionJob, error) {

	job := &models.ExtractionJob{
		ID:        batchID,
		Status:    "processing",
		CreatedAt: time.Now(),
		ExtractedData: &models.ExtractedData{
			RawText:        extractedText,
			ChapterContent: extractedText, // Use raw text as chapter content
		},
	}

	// Parse the text to extract structured data
	// This is a simplified version - you may want to use LLM for better extraction
	job.ExtractedData.CaseStudies = extractCaseStudiesFromText(extractedText)
	job.ExtractedData.MedicalTerms = extractMedicalTermsFromText(extractedText)

	job.Status = "completed"
	job.UpdatedAt = time.Now()

	return job, nil
}

// extractCaseStudiesFromText attempts to identify case studies in the text
func extractCaseStudiesFromText(text string) []models.CaseStudy {
	var caseStudies []models.CaseStudy

	// Simple heuristic: look for sections that start with "Case" followed by a number
	lines := strings.Split(text, "\n")
	var currentCase *models.CaseStudy
	var caseContent strings.Builder

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if this line starts a new case study
		if strings.HasPrefix(strings.ToLower(trimmed), "case ") {
			// Save previous case if exists
			if currentCase != nil {
				currentCase.Content = caseContent.String()
				caseStudies = append(caseStudies, *currentCase)
			}

			// Start new case
			parts := strings.SplitN(trimmed, ":", 2)
			caseNumber := strings.TrimPrefix(strings.ToLower(parts[0]), "case ")
			caseNumber = strings.TrimSpace(caseNumber)

			title := ""
			if len(parts) > 1 {
				title = strings.TrimSpace(parts[1])
			}

			currentCase = &models.CaseStudy{
				CaseNumber: caseNumber,
				Title:      title,
				Tags:       []string{"extracted_case"},
			}
			caseContent.Reset()
		} else if currentCase != nil {
			// Accumulate content for current case
			caseContent.WriteString(line)
			caseContent.WriteString("\n")
		}
	}

	// Save last case
	if currentCase != nil {
		currentCase.Content = caseContent.String()
		caseStudies = append(caseStudies, *currentCase)
	}

	return caseStudies
}

// extractMedicalTermsFromText attempts to identify medical terms and their definitions
func extractMedicalTermsFromText(text string) []models.MedicalTerm {
	var terms []models.MedicalTerm

	// Simple heuristic: look for bolded terms followed by definitions
	// This is a placeholder - you'd want to use NLP or LLM for better extraction
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Look for lines that might be term definitions
		// (e.g., lines with colons, or short lines followed by longer explanations)
		if strings.Contains(trimmed, ":") && len(trimmed) < 100 {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				term := strings.TrimSpace(parts[0])
				definition := strings.TrimSpace(parts[1])

				// If definition is too short, check next line
				if len(definition) < 20 && i+1 < len(lines) {
					definition += " " + strings.TrimSpace(lines[i+1])
				}

				if len(term) > 0 && len(definition) > 0 {
					terms = append(terms, models.MedicalTerm{
						Term:       term,
						Definition: definition,
						Category:   "general", // Could be improved with classification
						Tags:       []string{"extracted_term"},
					})
				}
			}
		}
	}

	return terms
}
