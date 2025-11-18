package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/domain/infoin"
	"esp-organizer/internal/models"
)

func main() {
	// Initialize taxonomy and LLM client
	taxonomy := infoin.NewRelationshipTaxonomy()
	llmClient := llm.NewLlamaClient()
	ctx := context.Background()

	// Read markdown file
	content, err := os.ReadFile("docs/training_data/human_extracted_links.md")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")

	var semanticLinks []models.SemanticLink
	unknownRelationships := make(map[string]int)

	// Parse three-line groups
	for i := 0; i < len(lines)-2; i += 3 {
		sourceLine := strings.TrimSpace(lines[i])
		relationLine := strings.TrimSpace(lines[i+1])
		targetLine := strings.TrimSpace(lines[i+2])

		if sourceLine == "" || relationLine == "" || targetLine == "" {
			continue
		}

		// Extract case-specific marker
		isCaseSpecific := strings.HasPrefix(sourceLine, "*")
		if isCaseSpecific {
			sourceLine = strings.TrimPrefix(sourceLine, "*")
			sourceLine = strings.TrimSpace(sourceLine)
		}

		// Normalize relationship type
		normalizedRelation := taxonomy.NormalizeRelationType(relationLine)

		// If normalization failed, use LLM
		if normalizedRelation == relationLine && !taxonomy.IsKnownType(normalizedRelation) {
			fmt.Printf("⚠️ Unknown relationship: '%s' (using LLM to classify)\n", relationLine)

			suggested, err := taxonomy.SuggestRelationshipType(llmClient, ctx, relationLine, sourceLine, targetLine)
			if err != nil || suggested == "unknown" {
				fmt.Printf("❌ Could not classify '%s', keeping as-is\n", relationLine)
				unknownRelationships[relationLine]++
				normalizedRelation = relationLine // Keep original
			} else {
				fmt.Printf("✅ LLM classified '%s' → '%s'\n", relationLine, suggested)
				normalizedRelation = suggested
			}
		}

		// Create semantic link
		link := models.SemanticLink{
			SourceTerm:     sourceLine,
			TargetTerm:     targetLine,
			RelationType:   normalizedRelation,
			IsCaseSpecific: isCaseSpecific,
			// ... other fields
		}

		semanticLinks = append(semanticLinks, link)
	}

	// Report statistics
	fmt.Printf("\n📊 Extraction Statistics:\n")
	fmt.Printf("Total links extracted: %d\n", len(semanticLinks))
	fmt.Printf("Unknown relationships: %d\n", len(unknownRelationships))

	if len(unknownRelationships) > 0 {
		fmt.Printf("\n⚠️ Unknown Relationships (consider adding to taxonomy):\n")
		for rel, count := range unknownRelationships {
			fmt.Printf("  • '%s' (%d occurrences)\n", rel, count)
		}
	}

	// Save to JSON
	output, _ := json.MarshalIndent(semanticLinks, "", "  ")
	os.WriteFile("docs/training_data/validated_links.json", output, 0644)

	fmt.Println("\n✅ Validation complete. Output saved to validated_links.json")
}
