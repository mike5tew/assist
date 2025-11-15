package infoin

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"

	"esp-organizer/internal/models"
)

// CuratedSemanticLink represents a human-verified semantic link with metrics
type CuratedSemanticLink struct {
	SourceTerm           string  `json:"source_term"`
	TargetTerm           string  `json:"target_term"`
	RelationType         string  `json:"relation_type"`
	SourceTermGenerality float64 `json:"source_term_generality"`
	TargetTermGenerality float64 `json:"target_term_generality"`
	SemanticDistance     float64 `json:"semantic_distance"`
	RelationshipStrength float64 `json:"relationship_strength"`
	Explanation          string  `json:"explanation"`
	HierarchyRationale   string  `json:"hierarchy_rationale"`
}

// HierarchyValidator validates LLM-generated metrics against curated examples
type HierarchyValidator struct {
	CuratedExamples []CuratedSemanticLink
}

// NewHierarchyValidator loads curated examples from JSON
func NewHierarchyValidator(curatedExamplesPath string) (*HierarchyValidator, error) {
	data, err := os.ReadFile(curatedExamplesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load curated examples: %w", err)
	}

	var examples []CuratedSemanticLink
	if err := json.Unmarshal(data, &examples); err != nil {
		return nil, fmt.Errorf("failed to parse curated examples: %w", err)
	}

	return &HierarchyValidator{
		CuratedExamples: examples,
	}, nil
}

// ValidateMetrics compares LLM output against curated examples
func (hv *HierarchyValidator) ValidateMetrics(llmMetrics HierarchyMetrics) ValidationReport {
	// Find the closest matching curated example
	closestExample := hv.findClosestExample(llmMetrics)

	// Calculate error metrics
	sourceGenError := math.Abs(llmMetrics.SourceTermGenerality - closestExample.SourceTermGenerality)
	targetGenError := math.Abs(llmMetrics.TargetTermGenerality - closestExample.TargetTermGenerality)
	distanceError := math.Abs(llmMetrics.SemanticDistance - closestExample.SemanticDistance)
	strengthError := math.Abs(llmMetrics.RelationshipStrength - closestExample.RelationshipStrength)

	avgError := (sourceGenError + targetGenError + distanceError + strengthError) / 4.0

	return ValidationReport{
		IsValid:        avgError < 0.15, // Accept if average error < 15%
		AverageError:   avgError,
		SourceGenError: sourceGenError,
		TargetGenError: targetGenError,
		DistanceError:  distanceError,
		StrengthError:  strengthError,
		ClosestExample: closestExample,
		Suggestion:     hv.generateSuggestion(llmMetrics, closestExample, avgError),
	}
}

// findClosestExample finds the curated example most similar to the LLM output
func (hv *HierarchyValidator) findClosestExample(llmMetrics HierarchyMetrics) CuratedSemanticLink {
	var closest CuratedSemanticLink
	minDistance := math.MaxFloat64

	for _, example := range hv.CuratedExamples {
		// Calculate Euclidean distance in metric space
		distance := math.Sqrt(
			math.Pow(llmMetrics.SourceTermGenerality-example.SourceTermGenerality, 2) +
				math.Pow(llmMetrics.TargetTermGenerality-example.TargetTermGenerality, 2) +
				math.Pow(llmMetrics.SemanticDistance-example.SemanticDistance, 2) +
				math.Pow(llmMetrics.RelationshipStrength-example.RelationshipStrength, 2),
		)

		if distance < minDistance {
			minDistance = distance
			closest = example
		}
	}

	return closest
}

// generateSuggestion creates feedback for the LLM
func (hv *HierarchyValidator) generateSuggestion(llmMetrics HierarchyMetrics, example CuratedSemanticLink, avgError float64) string {
	if avgError < 0.10 {
		return "Excellent match! Metrics align well with curated examples."
	}

	return fmt.Sprintf(`Your metrics differ from the closest curated example by %.1f%%.

Your Output:
  source_term_generality: %.2f
  target_term_generality: %.2f
  semantic_distance: %.2f
  relationship_strength: %.2f

Closest Curated Example (%s → %s):
  source_term_generality: %.2f
  target_term_generality: %.2f
  semantic_distance: %.2f
  relationship_strength: %.2f

Adjustment Hint: %s`,
		avgError*100,
		llmMetrics.SourceTermGenerality,
		llmMetrics.TargetTermGenerality,
		llmMetrics.SemanticDistance,
		llmMetrics.RelationshipStrength,
		example.SourceTerm,
		example.TargetTerm,
		example.SourceTermGenerality,
		example.TargetTermGenerality,
		example.SemanticDistance,
		example.RelationshipStrength,
		example.HierarchyRationale,
	)
}

// ValidationReport contains the result of comparing LLM output to curated data
type ValidationReport struct {
	IsValid        bool
	AverageError   float64
	SourceGenError float64
	TargetGenError float64
	DistanceError  float64
	StrengthError  float64
	ClosestExample CuratedSemanticLink
	Suggestion     string
}

// CorrectMetricsWithFeedback sends validation feedback back to LLM for correction
func (s *SemanticLinkService) CorrectMetricsWithFeedback(
	ctx context.Context,
	link *models.SemanticLink,
	initialMetrics HierarchyMetrics,
	validationReport ValidationReport,
) (*HierarchyMetrics, error) {

	if validationReport.IsValid {
		// Metrics are good, no correction needed
		return &initialMetrics, nil
	}

	// Send feedback to LLM for correction
	correctionPrompt := fmt.Sprintf(`Your initial analysis had an average error of %.1f%%.

%s

Please revise your metrics based on this feedback.

Original Analysis:
Source: "%s"
Target: "%s"
Relation: "%s"

Your Initial Metrics:
  source_term_generality: %.2f
  target_term_generality: %.2f
  semantic_distance: %.2f
  relationship_strength: %.2f

Provide corrected metrics as JSON:
{
  "source_term_generality": <0.0-1.0>,
  "target_term_generality": <0.0-1.0>,
  "semantic_distance": <0.0-1.0>,
  "relationship_strength": <0.0-1.0>,
  "explanation": "<why you adjusted>"
}`,
		validationReport.AverageError*100,
		validationReport.Suggestion,
		link.SourceTerm,
		link.TargetTerm,
		link.RelationType,
		initialMetrics.SourceTermGenerality,
		initialMetrics.TargetTermGenerality,
		initialMetrics.SemanticDistance,
		initialMetrics.RelationshipStrength,
	)

	response, err := s.LlmClient.Generate(ctx, correctionPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get corrected metrics: %w", err)
	}

	var correctedMetrics HierarchyMetrics
	if err := json.Unmarshal([]byte(cleanJSONResponse(response)), &correctedMetrics); err != nil {
		return nil, fmt.Errorf("failed to parse corrected metrics: %w", err)
	}

	return &correctedMetrics, nil
}
