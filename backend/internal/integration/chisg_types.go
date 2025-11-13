// /shared/types/chisg_integration.go
package types

import "time"

// KnowledgeQuery is the request sent to CHISG
type KnowledgeQuery struct {
	Topic          string   `json:"topic"`
	Concepts       []string `json:"concepts"`
	TargetAudience string   `json:"target_audience"`
}

// KnowledgeResponse is what CHISG returns
type KnowledgeResponse struct {
	Summary         string            `json:"summary"`
	KeyConcepts     []string          `json:"key_concepts"`
	Prerequisites   []string          `json:"prerequisites"`
	Analogies       []string          `json:"analogies"`
	ConfidenceScore float64           `json:"confidence_score"`
	RelatedTopics   map[string]string `json:"related_topics"`
	Timestamp       time.Time         `json:"timestamp"`
}

// HumanOSResponse wraps CHISG knowledge with age-appropriate adjustments
type HumanOSResponse struct {
	OriginalKnowledge *KnowledgeResponse `json:"knowledge"`
	AdjustedForAge    string             `json:"adjusted_summary"`
	DevelopmentStage  string             `json:"development_stage"`
	BarriersDetected  []string           `json:"barriers_detected"`
	Interventions     []string           `json:"recommended_interventions"`
}
