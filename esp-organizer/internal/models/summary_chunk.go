package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SummaryChunk represents a Tier 2 high-level summary of grouped semantic links
type SummaryChunk struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	SummaryText     string               `bson:"summary_text" json:"summary_text"`
	Vector          []float32            `bson:"vector,omitempty" json:"vector,omitempty"`
	Domain          string               `bson:"domain" json:"domain"`
	SemanticLinkIDs []primitive.ObjectID `bson:"semantic_link_ids" json:"semantic_link_ids"` // Tier 1 links used to create this summary

	// 🆕 NEW FIELDS for recursive Tier 2
	DerivedSemanticLinkIDs []primitive.ObjectID `bson:"derived_semantic_link_ids,omitempty" json:"derived_semantic_link_ids,omitempty"` // Higher-level links extracted FROM this summary
	DerivedLinkCount       int                  `bson:"derived_link_count,omitempty" json:"derived_link_count,omitempty"`

	Confidence float64                `bson:"confidence" json:"confidence"`
	CreatedAt  time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time              `bson:"updated_at" json:"updated_at"`
	Metadata   map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}
