package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HSGThread struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ThreadID string             `bson:"thread_id" json:"thread_id"` // e.g., "THREAD_BCELL_DEVELOPMENT"

	Title  string `bson:"title" json:"title"`
	Domain string `bson:"domain" json:"domain"`

	// Chain of semantic links
	SemanticLinkChain []string `bson:"semantic_link_chain" json:"semantic_link_chain"` // Array of link_ids

	// 🆕 THIS IS THE KEY: Narrative provides the wider context
	NarrativeSummary string `bson:"narrative_summary" json:"narrative_summary"`

	// 🆕 Contextual metadata that applies to multiple links
	ContextualNotes map[string]ContextualDetail `bson:"contextual_notes,omitempty" json:"contextual_notes,omitempty"`

	// Graph metadata
	RootConcept      string   `bson:"root_concept" json:"root_concept"`
	TerminalConcepts []string `bson:"terminal_concepts" json:"terminal_concepts"`
	ThreadType       string   `bson:"thread_type" json:"thread_type"` // e.g., "developmental_pathway", "causal_chain"

	Confidence float64   `bson:"confidence" json:"confidence"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`

	Metadata map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// 🆕 ContextualDetail stores rich context for specific concepts within a thread
type ContextualDetail struct {
	Temporal               string   `bson:"temporal,omitempty" json:"temporal,omitempty"`                               // When this applies
	Spatial                string   `bson:"spatial,omitempty" json:"spatial,omitempty"`                                 // Where this occurs
	Functional             string   `bson:"functional,omitempty" json:"functional,omitempty"`                           // What role it plays
	DistinguishingFeatures []string `bson:"distinguishing_features,omitempty" json:"distinguishing_features,omitempty"` // What makes it unique

	// Relationships to other concepts in the thread
	Preconditions []string `bson:"preconditions,omitempty" json:"preconditions,omitempty"` // What must happen first
	Consequences  []string `bson:"consequences,omitempty" json:"consequences,omitempty"`   // What happens next

	Explanation string `bson:"explanation,omitempty" json:"explanation,omitempty"`
}
