package main // 🛑 CORRECTED: Remove duplicate "package validatellmextraction"

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"esp-organizer/internal/models" // 🆕 Import canonical model
)

func main() {
	// Load human-curated links (now using canonical SemanticLink)
	humanData, _ := os.ReadFile("docs/training_data/curated_semantic_links_human_extracted.json")
	var humanLinks []models.SemanticLink // 🆕 Use canonical type
	json.Unmarshal(humanData, &humanLinks)

	// Load LLM-extracted links
	llmData, _ := os.ReadFile("docs/training_data/llm_extracted_links.json")
	var llmLinks []models.SemanticLink // 🆕 Use canonical type
	json.Unmarshal(llmData, &llmLinks)

	// Build index of human-curated composite keys
	humanKeys := make(map[string]bool)
	for _, link := range humanLinks {
		key := generateCompositeKey(link.SourceTerm, link.TargetTerm, link.RelationType)
		humanKeys[key] = true
	}

	// Count matches
	matches := 0
	for _, link := range llmLinks {
		key := generateCompositeKey(link.SourceTerm, link.TargetTerm, link.RelationType)
		if humanKeys[key] {
			matches++
			fmt.Printf("✅ MATCH: %s\n", key)
		} else {
			fmt.Printf("❌ MISS: %s\n", key)
		}
	}

	accuracy := float64(matches) / float64(len(llmLinks)) * 100
	fmt.Printf("\n📊 Accuracy: %.1f%% (%d matches / %d total LLM links)\n", accuracy, matches, len(llmLinks))
}

func generateCompositeKey(source, target, relationType string) string {
	return fmt.Sprintf("%s::%s::%s",
		strings.ToLower(strings.TrimSpace(source)),
		strings.ToLower(strings.TrimSpace(target)),
		strings.ToLower(strings.TrimSpace(relationType)),
	)
}
