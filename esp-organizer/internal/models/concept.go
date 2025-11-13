package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Influence represents the effect one concept has on another
type Influence struct {
	TargetConceptID primitive.ObjectID `json:"target_concept_id" bson:"target_concept_id"`
	RelationType    RelationType       `json:"relation_type" bson:"relation_type"`
	Strength        float64            `json:"strength" bson:"strength"` // 0.0 to 1.0 confidence score
	Bidirectional   bool               `json:"bidirectional" bson:"bidirectional,omitempty"`
	Evidence        string             `json:"evidence,omitempty" bson:"evidence,omitempty"` // Text evidence for this relationship
	Context         string             `json:"context,omitempty" bson:"context,omitempty"`   // Surrounding context from source
	SourceText      string             `json:"source_text,omitempty" bson:"source_text,omitempty"`
	ExtractedBy     string             `json:"extracted_by,omitempty" bson:"extracted_by,omitempty"` // "manual", "llama", etc.
	CreatedAt       int64              `json:"created_at" bson:"created_at"`
	UpdatedAt       int64              `json:"updated_at" bson:"updated_at"`
}

type RelationType string

const (
	Causes          RelationType = "causes"
	IsA             RelationType = "is_a"
	PartOf          RelationType = "part_of"
	AssociatedWith  RelationType = "associated_with"
	TreatedWith     RelationType = "treated_with"
	HasSymptom      RelationType = "has_symptom"
	Contraindicates RelationType = "contraindicates"
)

// Concept represents a central idea or entity and its direct relationships
type Concept struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Description    string             `json:"description" bson:"description"`
	DataSourceID   primitive.ObjectID `json:"data_source_id" bson:"data_source_id"`
	PrimaryChunkID primitive.ObjectID `json:"primary_chunk_id,omitempty" bson:"primary_chunk_id,omitempty"`
	Influences     []Influence        `json:"influences,omitempty" bson:"influences,omitempty"`
	CreatedAt      int64              `json:"created_at" bson:"created_at"`
	UpdatedAt      int64              `json:"updated_at" bson:"updated_at"`
}

// GraphQLResponse represents a response from a Weaviate GraphQL query.
type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []*GraphQLError        `json:"errors"`
}

type GraphQLError struct {
	Message   string                   `json:"message"`
	Locations []map[string]interface{} `json:"locations"`
	Path      []interface{}            `json:"path"`
}
