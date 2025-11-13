package models

// HSGTraversalResult contains the result of hierarchical graph traversal
type HSGTraversalResult struct {
	Tier2Chunks     []SummaryChunk   `json:"tier2_chunks"`
	Tier1Links      []SemanticLink   `json:"tier1_links"`
	Tier0Documents  []SubjectContent `json:"tier0_documents"`
	TotalConfidence float64          `json:"total_confidence"`
}

// RAGQueryContext contains context assembled for final generation
type RAGQueryContext struct {
	Query           string          `json:"query"`
	Domain          string          `json:"domain"`
	SummaryContext  []string        `json:"summary_context"`
	SemanticContext []string        `json:"semantic_context"`
	DocumentContext []string        `json:"document_context"`
	SourceChain     []TraversalStep `json:"source_chain"`
}

// TraversalStep tracks the path through the HSG for explainability
type TraversalStep struct {
	Tier       int     `json:"tier"`
	ObjectID   string  `json:"object_id"`
	ObjectType string  `json:"object_type"`
	Content    string  `json:"content"`
	Confidence float64 `json:"confidence"`
}
