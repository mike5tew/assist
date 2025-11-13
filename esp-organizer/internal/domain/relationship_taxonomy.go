package InfoIn

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/InfoFlow/InfoOut/llm"
	"fmt"
	"os"
	"strings"
)

// RelationshipTaxonomy defines the standardized relationship types
type RelationshipTaxonomy struct {
	// Core categories (always normalized to these)
	CoreTypes map[string]RelationshipType

	// Synonym map (user phrases → core type)
	Synonyms map[string]string
}

type RelationshipType struct {
	Name        string   // Canonical name (e.g., "causes")
	Category    string   // Broad category (causal, structural, functional, etc.)
	Direction   string   // directional | bidirectional
	Description string   // Human-readable explanation
	Examples    []string // Example usage
}

// NewRelationshipTaxonomy creates the canonical relationship taxonomy
func NewRelationshipTaxonomy() *RelationshipTaxonomy {
	taxonomy := &RelationshipTaxonomy{
		CoreTypes: make(map[string]RelationshipType),
		Synonyms:  make(map[string]string),
	}

	// Define core relationship types
	coreTypes := []RelationshipType{
		// Causal relationships
		{
			Name:        "causes",
			Category:    "causal",
			Direction:   "directional",
			Description: "A directly causes B",
			Examples:    []string{"BTK mutation causes B-cell deficiency"},
		},
		{
			Name:        "leads_to",
			Category:    "causal",
			Direction:   "directional",
			Description: "A leads to B (indirect causation)",
			Examples:    []string{"B-cell deficiency leads to recurrent infections"},
		},

		// Structural relationships
		{
			Name:        "part_of",
			Category:    "structural",
			Direction:   "directional",
			Description: "A is a component of B",
			Examples:    []string{"B cell part_of immune system"},
		},
		{
			Name:        "contains",
			Category:    "structural",
			Direction:   "directional",
			Description: "A contains B",
			Examples:    []string{"Immune system contains B cells"},
		},

		// Developmental relationships
		{
			Name:        "develops_into",
			Category:    "developmental",
			Direction:   "directional",
			Description: "A matures/transitions into B",
			Examples:    []string{"Pro-B cell develops_into Pre-B cell"},
		},
		{
			Name:        "during_this",
			Category:    "temporal",
			Direction:   "directional",
			Description: "Event/property occurs during A",
			Examples:    []string{"Large pre-B cell during_this mu chain expression"},
		},

		// Functional relationships
		{
			Name:        "regulates",
			Category:    "functional",
			Direction:   "directional",
			Description: "A controls/modulates B",
			Examples:    []string{"BTK regulates B-cell receptor signaling"},
		},
		{
			Name:        "enables",
			Category:    "functional",
			Direction:   "directional",
			Description: "A makes B possible",
			Examples:    []string{"BTK enables B-cell survival"},
		},
		{
			Name:        "inhibits",
			Category:    "functional",
			Direction:   "directional",
			Description: "A prevents/blocks B",
			Examples:    []string{"Regulatory T cell inhibits immune response"},
		},

		// Clinical relationships
		{
			Name:        "treats",
			Category:    "clinical",
			Direction:   "directional",
			Description: "A is a treatment for B",
			Examples:    []string{"IVIG treats antibody deficiency"},
		},
		{
			Name:        "diagnoses",
			Category:    "clinical",
			Direction:   "directional",
			Description: "A is used to diagnose B",
			Examples:    []string{"Flow cytometry diagnoses XLA"},
		},
		{
			Name:        "manifests_as",
			Category:    "clinical",
			Direction:   "directional",
			Description: "Disease A presents as symptom B",
			Examples:    []string{"XLA manifests_as recurrent infections"},
		},

		// Measurement relationships
		{
			Name:        "has_measurement",
			Category:    "measurement",
			Direction:   "directional",
			Description: "Entity A has measured value B",
			Examples:    []string{"*XLA patient has_measurement WBC 5100/µl"},
		},
		{
			Name:        "normal_range",
			Category:    "reference",
			Direction:   "directional",
			Description: "Entity A has normal range B",
			Examples:    []string{"9-year-old child normal_range 4500-13500 WBC/µl"},
		},

		// Equivalence/Classification
		{
			Name:        "is_a",
			Category:    "classification",
			Direction:   "directional",
			Description: "A is a type of B (hierarchical)",
			Examples:    []string{"B cell is_a lymphocyte"},
		},
		{
			Name:        "equivalent_to",
			Category:    "classification",
			Direction:   "bidirectional",
			Description: "A and B are the same thing (synonyms)",
			Examples:    []string{"XLA equivalent_to X-linked agammaglobulinemia"},
		},

		// Location
		{
			Name:        "located_at",
			Category:    "spatial",
			Direction:   "directional",
			Description: "A is found at location B",
			Examples:    []string{"BTK gene located_at Xq22"},
		},
	}

	// Store core types
	for _, coreType := range coreTypes {
		taxonomy.CoreTypes[coreType.Name] = coreType
	}

	// Define synonyms (natural language → canonical)
	synonymMappings := map[string]string{
		// Causal synonyms
		"results_in":    "causes",
		"produces":      "causes",
		"induces":       "causes",
		"triggers":      "causes",
		"gives_rise_to": "causes",
		"brings_about":  "causes",
		"results_from":  "caused_by", // Reverse direction
		"caused_by":     "caused_by",
		"due_to":        "caused_by",

		// Developmental synonyms
		"matures_into":        "develops_into",
		"transitions_to":      "develops_into",
		"becomes":             "develops_into",
		"differentiates_into": "develops_into",

		// Structural synonyms
		"component_of":    "part_of",
		"subcomponent_of": "part_of",
		"element_of":      "part_of",
		"comprises":       "contains",
		"made_up_of":      "contains",
		"consists_of":     "contains",

		// Functional synonyms
		"controls":    "regulates",
		"modulates":   "regulates",
		"influences":  "regulates",
		"activates":   "enables",
		"facilitates": "enables",
		"promotes":    "enables",
		"suppresses":  "inhibits",
		"blocks":      "inhibits",
		"prevents":    "inhibits",

		// Clinical synonyms
		"therapy_for":      "treats",
		"used_to_treat":    "treats",
		"presents_as":      "manifests_as",
		"characterized_by": "manifests_as",
		"shows":            "manifests_as",

		// Measurement synonyms
		"measured_at":     "has_measurement",
		"value_of":        "has_measurement",
		"reference_range": "normal_range",
		"normal_value":    "normal_range",

		// Classification synonyms
		"type_of":       "is_a",
		"subtype_of":    "is_a",
		"instance_of":   "is_a",
		"same_as":       "equivalent_to",
		"also_known_as": "equivalent_to",
		"synonym_of":    "equivalent_to",

		// Location synonyms
		"found_at":      "located_at",
		"positioned_at": "located_at",
		"situated_at":   "located_at",
	}

	taxonomy.Synonyms = synonymMappings

	return taxonomy
}

// NormalizeRelationType converts any relationship phrase to canonical form
func (rt *RelationshipTaxonomy) NormalizeRelationType(input string) string {
	// Clean and normalize input
	normalized := strings.ToLower(strings.TrimSpace(input))
	normalized = strings.ReplaceAll(normalized, " ", "_")

	// Check if it's already a core type
	if _, exists := rt.CoreTypes[normalized]; exists {
		return normalized
	}

	// Check if it's a known synonym
	if canonical, exists := rt.Synonyms[normalized]; exists {
		return canonical
	}

	// If unknown, return as-is (will be sent to LLM for classification)
	return normalized
}

// GetRelationshipInfo returns detailed information about a relationship type
func (rt *RelationshipTaxonomy) GetRelationshipInfo(relationType string) (*RelationshipType, bool) {
	if info, exists := rt.CoreTypes[relationType]; exists {
		return &info, true
	}
	return nil, false
}

// SuggestRelationshipType uses LLM to classify unknown relationship types
func (rt *RelationshipTaxonomy) SuggestRelationshipType(llmClient *llm.LlamaClient, ctx context.Context, input string, sourceTerm string, targetTerm string) (string, error) {
	// Build list of core types for LLM
	coreTypesList := make([]string, 0, len(rt.CoreTypes))
	for name := range rt.CoreTypes {
		coreTypesList = append(coreTypesList, name)
	}

	prompt := fmt.Sprintf(`You are a relationship classifier for a medical knowledge graph.

Given the relationship phrase "%s" connecting:
  Source: "%s"
  Target: "%s"

Which of these canonical relationship types best describes it?
%s

Respond with ONLY the canonical relationship type name, nothing else.

If none match closely, respond with "unknown".`,
		input, sourceTerm, targetTerm, strings.Join(coreTypesList, ", "))

	response, err := llmClient.Generate(ctx, prompt)
	if err != nil {
		return "unknown", err
	}

	suggested := strings.ToLower(strings.TrimSpace(response))

	// Validate LLM suggestion
	if _, exists := rt.CoreTypes[suggested]; exists {
		return suggested, nil
	}

	return "unknown", fmt.Errorf("LLM suggested invalid type: %s", suggested)
}

// IsKnownType checks if a relationship type is in the taxonomy
func (rt *RelationshipTaxonomy) IsKnownType(relationType string) bool {
	normalized := strings.ToLower(strings.TrimSpace(relationType))

	// Check core types
	if _, exists := rt.CoreTypes[normalized]; exists {
		return true
	}

	// Check synonyms
	if _, exists := rt.Synonyms[normalized]; exists {
		return true
	}

	return false
}

// GetAllCoreTypes returns a list of all canonical relationship types
func (rt *RelationshipTaxonomy) GetAllCoreTypes() []string {
	types := make([]string, 0, len(rt.CoreTypes))
	for name := range rt.CoreTypes {
		types = append(types, name)
	}
	return types
}

// ExportTaxonomy exports the taxonomy as JSON for documentation
func (rt *RelationshipTaxonomy) ExportTaxonomy(filepath string) error {
	output := map[string]interface{}{
		"core_types": rt.CoreTypes,
		"synonyms":   rt.Synonyms,
		"version":    "1.0.0",
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}
