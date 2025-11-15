package extraction

import (
	"context"
	"esp-organizer/internal/models"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// ProcessExtractedText takes raw text and extracts structured data like case studies and medical terms.
func ProcessExtractedText(ctx context.Context, text string, chapterInfo models.ChapterInfo, bookSource models.Source, jobID string) (*models.ProcessingJob, error) {
	log.Printf("[Job %s] Starting detailed text processing. Text length: %d", jobID, len(text))

	caseStudies := extractCaseStudies(text, chapterInfo)
	medicalTerms := extractMedicalTerms(text)

	extractedData := &models.ExtractedChapterData{
		RawText:        text,
		ChapterContent: text, // For now, use raw text as chapter content
		CaseStudies:    caseStudies,
		MedicalTerms:   medicalTerms,
	}

	job := &models.ProcessingJob{
		ID:            jobID,
		Status:        "processing_completed",
		ChapterInfo:   chapterInfo,
		BookSource:    bookSource,
		ExtractedData: extractedData,
		CreatedAt:     time.Now(),
		CompletedAt:   func() *time.Time { t := time.Now(); return &t }(),
	}

	log.Printf("[Job %s] Finished detailed text processing. Found %d case studies and %d medical terms.", jobID, len(caseStudies), len(medicalTerms))

	return job, nil
}

// extractCaseStudies finds case study patterns in the extracted text
func extractCaseStudies(text string, chapterInfo models.ChapterInfo) []models.CaseStudy {
	var caseStudies []models.CaseStudy

	if len(strings.TrimSpace(text)) == 0 {
		log.Printf("⚠️ No text available - skipping case study extraction")
		return caseStudies
	}

	patterns := []struct {
		pattern string
		title   string
	}{
		{`Case\s+[Ss]tudy\s*[#:]?\s*(\d+)[:.\-]*(.*?)(?=\n\n|Case\s+[Ss]tudy|$)`, "Case Study $1: $2"},
		{`Clinical\s+[Cc]ase\s*[#:]?\s*(\d+)[:.\-]*(.*?)(?=\n\n|Clinical\s+[Cc]ase|$)`, "Clinical Case $1: $2"},
		{`Patient\s+[Cc]ase\s*[#:]?\s*(\d+)[:.\-]*(.*?)(?=\n\n|Patient\s+[Cc]ase|$)`, "Patient Case $1: $2"},
	}

	caseNumber := 1

	for _, p := range patterns {
		re := regexp.MustCompile(p.pattern)
		matches := re.FindAllStringSubmatch(text, -1)

		for _, match := range matches {
			if len(match) >= 3 {
				contentStart := strings.Index(text, match[0])
				if contentStart == -1 {
					continue
				}

				remainingText := text[contentStart+len(match[0]):]
				endPatterns := []string{"\n\nCase Study", "\n\nClinical Case", "\n\nPatient Case", "\n\nDiscussion", "\n\nConclusion"}

				var contentEnd int
				for _, endPattern := range endPatterns {
					if idx := strings.Index(remainingText, endPattern); idx != -1 {
						if contentEnd == 0 || idx < contentEnd {
							contentEnd = idx
						}
					}
				}

				var content string
				if contentEnd > 0 {
					content = strings.TrimSpace(match[0] + remainingText[:contentEnd])
				} else {
					content = strings.TrimSpace(match[0] + remainingText)
				}

				if len(content) > 2000 {
					content = content[:2000] + "..."
				}

				title := fmt.Sprintf("Case %s-%d", chapterInfo.ChapterNumber, caseNumber)
				if match[2] != "" {
					cleanedTitle := strings.TrimSpace(match[2])
					if len(cleanedTitle) > 0 {
						title = fmt.Sprintf("Case %s-%d: %s", chapterInfo.ChapterNumber, caseNumber, cleanedTitle)
					}
				}

				caseStudy := models.CaseStudy{
					CaseNumber:       fmt.Sprintf("%s-%d", chapterInfo.ChapterNumber, caseNumber),
					Title:            title,
					Content:          content,
					ClinicalFindings: extractClinicalFindings(content),
					Tags:             []string{"immunology", "case_study", "chapter-" + chapterInfo.ChapterNumber},
				}

				caseStudies = append(caseStudies, caseStudy)
				caseNumber++

				if len(caseStudies) >= 5 {
					return caseStudies
				}
			}
		}
	}

	if len(caseStudies) == 0 {
		caseStudies = extractCaseStudyFallbacks(text, chapterInfo, caseNumber)
	}

	log.Printf("Extracted %d case studies from chapter %s", len(caseStudies), chapterInfo.ChapterNumber)
	return caseStudies
}

// extractClinicalFindings extracts clinical findings from case study content
func extractClinicalFindings(content string) string {
	patterns := []string{
		`presented with[^.]*\.`,
		`complaining of[^.]*\.`,
		`symptoms included[^.]*\.`,
		`physical exam[^.]*\.`,
		`laboratory findings[^.]*\.`,
	}

	lowerContent := strings.ToLower(content)
	var findings []string

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindString(lowerContent); match != "" {
			findings = append(findings, strings.TrimSpace(match))
		}
	}

	if len(findings) > 0 {
		return strings.Join(findings, "; ")
	}

	if len(content) > 200 {
		return content[:200] + "..."
	}
	return content
}

// extractCaseStudyFallbacks provides fallback extraction when regex patterns don't match
func extractCaseStudyFallbacks(text string, chapterInfo models.ChapterInfo, startNumber int) []models.CaseStudy {
	var caseStudies []models.CaseStudy
	sections := strings.Split(text, "\n\n")

	for _, section := range sections {
		section = strings.TrimSpace(section)
		if len(section) < 100 || len(section) > 2000 {
			continue
		}

		lowerSection := strings.ToLower(section)
		caseIndicators := []string{"patient", "case", "presented", "symptoms", "diagnosis", "treatment"}
		indicatorCount := 0

		for _, indicator := range caseIndicators {
			if strings.Contains(lowerSection, indicator) {
				indicatorCount++
			}
		}

		if indicatorCount >= 3 {
			caseStudy := models.CaseStudy{
				CaseNumber:       fmt.Sprintf("%s-%d", chapterInfo.ChapterNumber, startNumber),
				Title:            fmt.Sprintf("Clinical Scenario %s-%d", chapterInfo.ChapterNumber, startNumber),
				Content:          section,
				ClinicalFindings: extractClinicalFindings(section),
				Tags:             []string{"immunology", "case_study", "auto-detected"},
			}
			caseStudies = append(caseStudies, caseStudy)
			startNumber++
			if len(caseStudies) >= 3 {
				break
			}
		}
	}
	return caseStudies
}

// extractMedicalTerms finds medical terminology in the extracted text
func extractMedicalTerms(text string) []models.MedicalTerm {
	var medicalTerms []models.MedicalTerm
	if len(strings.TrimSpace(text)) == 0 {
		log.Printf("⚠️ No text available - skipping medical term extraction")
		return medicalTerms
	}

	re := regexp.MustCompile(`\b([A-Z][a-zA-Z-]{2,}(?:\s[a-zA-Z-]+)*|[A-Z]{2,})\b`)
	matches := re.FindAllString(text, -1)

	termCounts := make(map[string]int)
	for _, match := range matches {
		if len(match) > 2 && !isCommonWord(match) {
			termCounts[match]++
		}
	}

	for term, count := range termCounts {
		if count > 0 {
			context := extractTermContext(text, term)
			medicalTerm := models.MedicalTerm{
				Term:     term,
				Category: categorizeImmunologyTerm(term),
				Context:  context,
			}
			medicalTerms = append(medicalTerms, medicalTerm)
		}
	}

	if len(medicalTerms) == 0 {
		log.Printf("No specific medical terms identified via pattern matching in extracted text.")
	} else {
		log.Printf("Identified %d potential medical terms via pattern matching.", len(medicalTerms))
	}

	return medicalTerms
}

// isCommonWord checks if a word is a common English word that should be ignored.
func isCommonWord(word string) bool {
	commonWords := map[string]bool{
		"The": true, "And": true, "For": true, "With": true, "From": true,
		"This": true, "That": true, "But": true, "Not": true, "Are": true,
		"Was": true, "Were": true, "Has": true, "Have": true, "Case": true,
		"Study": true, "Figure": true, "Table": true,
	}
	_, found := commonWords[strings.Title(strings.ToLower(word))]
	return found
}

func extractTermContext(text, term string) []string {
	context := []string{}
	lowerText := strings.ToLower(text)
	lowerTerm := strings.ToLower(term)

	words := strings.Fields(lowerText)
	for i, word := range words {
		if strings.Contains(word, lowerTerm) {
			start := max(0, i-2)
			end := min(len(words), i+3)
			contextWords := words[start:end]
			contextStr := strings.Join(contextWords, " ")
			context = append(context, contextStr)
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
