package infoin

import "fmt"

// HSGPatternExamples provides curated pattern examples for hierarchy metric computation
const HSGPatternExamples = `
You are a medical knowledge graph builder. Your task is to analyze term relationships and compute relativistic hierarchy metrics.

CURATED PATTERN EXAMPLES (Human-Verified Ground Truth):

Example 1: ROOT CONCEPT → CATEGORY
Source: "Immunodeficiency Disorders"
Target: "Primary Immunodeficiency"
Relation: "encompasses"

Human-Curated Metrics:
  source_term_generality: 0.95 (very broad domain concept)
  target_term_generality: 0.75 (subcategory)
  semantic_distance: 0.15 (direct hierarchical relationship)
  relationship_strength: 0.90 (strong categorical link)

Rationale: Broad domain term encompassing a more specific category. 
Source is at the system level, target is a major subdivision.

---

Example 2: SPECIFIC DISEASE → MOLECULAR MECHANISM
Source: "X-Linked Agammaglobulinemia"
Target: "BTK gene mutation"
Relation: "caused_by"

Human-Curated Metrics:
  source_term_generality: 0.40 (specific disease entity)
  target_term_generality: 0.15 (molecular-level detail)
  semantic_distance: 0.10 (very close - direct causation)
  relationship_strength: 0.95 (definitive causal link)

Rationale: Named disease caused by known genetic defect. 
Source is clinical entity, target is genetic mechanism.

---

Example 3: PATHOLOGICAL STATE → CLINICAL SYMPTOM
Source: "B-cell deficiency"
Target: "Recurrent bacterial infections"
Relation: "causes"

Human-Curated Metrics:
  source_term_generality: 0.50 (intermediate - pathological state)
  target_term_generality: 0.30 (clinical manifestation)
  semantic_distance: 0.20 (moderate - indirect causation)
  relationship_strength: 0.85 (strong but not absolute)

Rationale: Immunological defect leading to clinical outcome.
Causation is strong but mediated by immune system function.
`

// HSGPatternExamplesWithFeedback includes self-assessment instructions
const HSGPatternExamplesWithFeedback = `
You are a medical knowledge graph builder. Your task is to analyze term relationships and compute relativistic hierarchy metrics.

LEARNING APPROACH:
1. Study the pattern examples below
2. Apply those patterns to new term pairs
3. Self-assess your output against the examples
4. Adjust your scoring if your output deviates from the pattern

CURATED PATTERN EXAMPLES (Human-Verified Ground Truth):

Example 1: ROOT CONCEPT → CATEGORY
Source: "Immunodeficiency Disorders"
Target: "Primary Immunodeficiency"
Relation: "encompasses"

Human-Curated Metrics:
  source_term_generality: 0.95 (very broad domain concept)
  target_term_generality: 0.75 (subcategory)
  semantic_distance: 0.15 (direct hierarchical relationship)
  relationship_strength: 0.90 (strong categorical link)

Rationale: Broad domain term encompassing a more specific category. 
Source is at the system level, target is a major subdivision.

---

Example 2: SPECIFIC DISEASE → MOLECULAR MECHANISM
Source: "X-Linked Agammaglobulinemia"
Target: "BTK gene mutation"
Relation: "caused_by"

Human-Curated Metrics:
  source_term_generality: 0.40 (specific disease entity)
  target_term_generality: 0.15 (molecular-level detail)
  semantic_distance: 0.10 (very close - direct causation)
  relationship_strength: 0.95 (definitive causal link)

Rationale: Named disease caused by known genetic defect. 
Source is clinical entity, target is genetic mechanism.

---

Example 3: PATHOLOGICAL STATE → CLINICAL SYMPTOM
Source: "B-cell deficiency"
Target: "Recurrent bacterial infections"
Relation: "causes"

Human-Curated Metrics:
  source_term_generality: 0.50 (intermediate - pathological state)
  target_term_generality: 0.30 (clinical manifestation)
  semantic_distance: 0.20 (moderate - indirect causation)
  relationship_strength: 0.85 (strong but not absolute)

Rationale: Immunological defect leading to clinical outcome.
Causation is strong but mediated by immune system function.

---

SELF-ASSESSMENT INSTRUCTIONS:
After generating metrics for a new term pair:
1. Compare your source_term_generality to Example 1 (0.95), Example 2 (0.40), Example 3 (0.50)
2. Ask: "Is my source term as broad as 'Immunodeficiency Disorders' (0.95)?"
3. Ask: "Is my source term as specific as 'BTK gene mutation' (0.15)?"
4. Adjust your score to match the pattern

NOW ANALYZE THIS NEW TERM PAIR:
Source: "%s"
Target: "%s"
Relation: "%s"
Context: "%s"

REQUIRED OUTPUT FORMAT (JSON):
{
  "source_term_generality": <0.0-1.0>,
  "target_term_generality": <0.0-1.0>,
  "semantic_distance": <0.0-1.0>,
  "relationship_strength": <0.0-1.0>,
  "explanation": "<brief justification>",
  "self_assessment": {
    "compared_to_example": <1|2|3>,
    "confidence": <0.0-1.0>,
    "reasoning": "<why I chose these metrics>"
  }
}
`

// HSGAnalysisPrompt generates a prompt for computing hierarchy metrics
func HSGAnalysisPrompt(sourceTerm, targetTerm, relationType, context string) string {
	return HSGPatternExamples + fmt.Sprintf(`

NOW ANALYZE THIS NEW TERM PAIR:
Source: "%s"
Target: "%s"
Relation Type: "%s"
Context: "%s"

Return ONLY valid JSON (no markdown):
{
  "source_term_generality": <float 0.0-1.0>,
  "target_term_generality": <float 0.0-1.0>,
  "semantic_distance": <float 0.0-1.0>,
  "relationship_strength": <float 0.0-1.0>,
  "explanation": "<brief justification>"
}
`, sourceTerm, targetTerm, relationType, context)
}
