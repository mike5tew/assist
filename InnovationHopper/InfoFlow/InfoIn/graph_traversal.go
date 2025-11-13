package InfoIn

type TraversalConfig struct {
	// Depth & Scope
	MaxDepth                  int     `json:"max_depth"`                   // 1-5 hops
	SemanticDistanceThreshold float64 `json:"semantic_distance_threshold"` // 0.0-1.0
	GeneralityGap             float64 `json:"generality_gap"`              // Min diff between source/target

	// Quality Filtering
	MinRelationshipStrength float64 `json:"min_relationship_strength"` // 0.0-1.0
	ConfidenceThreshold     float64 `json:"confidence_threshold"`

	// Result Ranking
	WeightByDistance   bool `json:"weight_by_distance"`   // Prefer close relationships?
	WeightByGenerality bool `json:"weight_by_generality"` // Prefer hierarchical jumps?
	WeightByStrength   bool `json:"weight_by_strength"`   // Prefer confident links?
}

// QualityProfile presets for different use cases
var QualityProfiles = map[string]TraversalConfig{
	"high_precision": {
		MaxDepth:                  2,
		SemanticDistanceThreshold: 0.15,
		MinRelationshipStrength:   0.85,
		ConfidenceThreshold:       0.90,
		WeightByStrength:          true,
	},
	"balanced": {
		MaxDepth:                  3,
		SemanticDistanceThreshold: 0.30,
		MinRelationshipStrength:   0.70,
		ConfidenceThreshold:       0.75,
		WeightByDistance:          true,
	},
	"discovery": {
		MaxDepth:                  4,
		SemanticDistanceThreshold: 0.50,
		MinRelationshipStrength:   0.60,
		ConfidenceThreshold:       0.60,
		WeightByGenerality:        true, // Find surprising jumps
	},
}
