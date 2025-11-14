package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"esp-organizer/internal/InfoFlow/InfoIn"
	"esp-organizer/internal/InfoFlow/InfoOut/llm"
	"esp-organizer/internal/models"
)

func main() {
	// Parse command-line flags
	chapterFile := flag.String("chapter", "", "Path to chapter text file")
	outputFile := flag.String("output", "extracted_links.json", "Output JSON file")
	domain := flag.String("domain", "immunology", "Subject domain")
	flag.Parse()

	if *chapterFile == "" {
		log.Fatal("Usage: extract-semantic-links -chapter <file> [-output <file>] [-domain <name>]")
	}

	// Read chapter content
	content, err := os.ReadFile(*chapterFile)
	if err != nil {
		log.Fatalf("Failed to read chapter file: %v", err)
	}

	log.Printf("📚 Extracting semantic links from: %s", *chapterFile)
	log.Printf("📊 Chapter size: %d characters", len(content))

	// Initialize LLM client
	llmClient := llm.NewLlamaClient()

	// **SIMPLE APPROACH: Direct LLM extraction, no pattern matching**
	ctx := context.Background()
	extractedLinks, err := extractWithLLM(ctx, llmClient, string(content), *domain)
	if err != nil {
		log.Fatalf("LLM extraction failed: %v", err)
	}

	log.Printf("✅ Extracted %d semantic links", len(extractedLinks))

	// Save to JSON for manual review
	output, err := json.MarshalIndent(extractedLinks, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(*outputFile, output, 0644); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	log.Printf("💾 Saved to: %s", *outputFile)
	log.Println("")
	log.Println("📝 Next steps:")
	log.Println("   1. Review the extracted links in the output file")
	log.Println("   2. Correct any mistakes (this is your revision exercise!)")
	log.Println("   3. Run: make validate-curated-links (when you have corrections)")
	log.Println("   4. Add corrected links to: docs/training_data/human_extracted_links.md")
}

// extractWithLLM performs simple LLM-based extraction without pattern matching
func extractWithLLM(ctx context.Context, llmClient *llm.LlamaClient, content string, domain string) ([]models.SemanticLink, error) {
	prompt := fmt.Sprintf(`Extract atomic semantic relationships from this %s text.

IMPORTANT RULES:
1. Each relationship should be ONE clear fact (e.g., "A causes B")
2. Use minimal context (just enough to be clear)
3. NO explanatory narrative
4. Mark case-specific facts with * at the start of source_term

For each relationship, provide:
- source_term: The subject (with * if case-specific)
- target_term: The object
- relation_type: [causes, part_of, develops_into, has_measurement, normal_range, etc.]
- context: Minimal (e.g., "during B-cell development")

Example:
{
  "source_term": "Large pre-B cell",
  "target_term": "mu heavy chain",
  "relation_type": "expresses",
  "context": "Transiently during development",
  "source_term_generality": 0.40,
  "target_term_generality": 0.35,
  "semantic_distance": 0.08,
  "relationship_strength": 0.95
}

Case-specific example:
{
  "source_term": "*XLA patient 9 years old",
  "target_term": "5100 white blood cells per microliter",
  "relation_type": "has_measurement",
  "context": "Laboratory test",
  "is_case_specific": true,
  "case_study_id": "XLA_Case_1"
}

Text to analyze:
%s

Return ONLY a JSON array of relationships.`, domain, content)

	response, err := llmClient.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Clean response (remove markdown code fences if present)
	response = InfoIn.CleanJSONResponse(response)

	// Parse JSON
	var links []models.SemanticLink
	if err := json.Unmarshal([]byte(response), &links); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	// **POST-PROCESSING: Add timestamps, normalize, generate composite keys**
	for i := range links {
		links[i].Domain = domain
		links[i].CreatedAt = time.Now()
		links[i].CompositeKey = generateCompositeKey(links[i].SourceTerm, links[i].TargetTerm, links[i].RelationType)

		// Set case-specific metadata if * marker present
		if strings.HasPrefix(links[i].SourceTerm, "*") {
			links[i].IsCaseSpecific = true
			links[i].SourceTerm = strings.TrimPrefix(links[i].SourceTerm, "*")
			links[i].SourceTerm = strings.TrimSpace(links[i].SourceTerm)

			// Extract disease category from source term if possible
			if strings.Contains(strings.ToLower(links[i].SourceTerm), "xla") {
				links[i].CaseStudyID = "XLA_Case_1"
				links[i].PatientMetadata = &models.PatientMetadata{
					DiseaseCategory: "XLA",
					CaseNumber:      1,
				}
			}
		}
	}

	return links, nil
}

func generateCompositeKey(source, target, relationType string) string {
	return fmt.Sprintf("%s::%s::%s",
		strings.ToLower(strings.TrimSpace(source)),
		strings.ToLower(strings.TrimSpace(target)),
		strings.ToLower(strings.TrimSpace(relationType)),
	)
}
