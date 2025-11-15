package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SemanticLink captures a single, atomic relationship between two concepts.
type SemanticLink struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	CompositeKey string `bson:"composite_key" json:"composite_key"` // "source::target::relation"

	SourceID     primitive.ObjectID `bson:"source_id,omitempty" json:"source_id"`
	TargetID     primitive.ObjectID `bson:"target_id,omitempty" json:"target_id"`
	SourceTerm   string             `bson:"source_term" json:"source_term"` // Normalized
	TargetTerm   string             `bson:"target_term" json:"target_term"` // Normalized
	RelationType string             `bson:"relation_type" json:"relation_type"`

	// Relativistic Hierarchy Metrics
	SourceTermGenerality float64 `bson:"source_term_generality" json:"source_term_generality"`
	TargetTermGenerality float64 `bson:"target_term_generality" json:"target_term_generality"`
	SemanticDistance     float64 `bson:"semantic_distance" json:"semantic_distance"`
	RelationshipStrength float64 `bson:"relationship_strength" json:"relationship_strength"`

	Context         string               `bson:"context" json:"context"`
	ContextOriginal string               `bson:"context_original,omitempty" json:"context_original,omitempty"`
	IsParentOf      []primitive.ObjectID `bson:"is_parent_of,omitempty" json:"is_parent_of,omitempty"`
	IsChildOf       []primitive.ObjectID `bson:"is_child_of,omitempty" json:"is_child_of,omitempty"`

	Confidence float64                `bson:"confidence" json:"confidence"`
	Domain     string                 `bson:"domain" json:"domain"`
	Subject    string                 `bson:"subject,omitempty" json:"subject,omitempty"`
	Vector     []float32              `bson:"vector,omitempty" json:"vector,omitempty"`
	CreatedAt  time.Time              `bson:"created_at" json:"created_at"`
	Metadata   map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`

	SourceTermOriginal string `bson:"source_term_original,omitempty" json:"source_term_original,omitempty"`
	TargetTermOriginal string `bson:"target_term_original,omitempty" json:"target_term_original,omitempty"`

	// Contextual Qualifiers
	TemporalQualifier string `bson:"temporal_qualifier,omitempty" json:"temporal_qualifier,omitempty"`
	SpatialQualifier  string `bson:"spatial_qualifier,omitempty" json:"spatial_qualifier,omitempty"`
	FunctionalState   string `bson:"functional_state,omitempty" json:"functional_state,omitempty"`
	LifecycleStage    string `bson:"lifecycle_stage,omitempty" json:"lifecycle_stage,omitempty"`

	Conditions    []Condition `bson:"conditions,omitempty" json:"conditions,omitempty"`
	ContrastsWith []Contrast  `bson:"contrasts_with,omitempty" json:"contrasts_with,omitempty"`
	LinkType      string      `bson:"link_type" json:"link_type"`

	RelatedLinkIDs []string `bson:"related_link_ids,omitempty" json:"related_link_ids,omitempty"`
	IsCaseSpecific bool     `bson:"is_case_specific" json:"is_case_specific"`
	CaseStudyID    string   `bson:"case_study_id,omitempty" json:"case_study_id,omitempty"`

	PatientMetadata *PatientMetadata `bson:"patient_metadata,omitempty" json:"patient_metadata,omitempty"`
	StudentMetadata *StudentMetadata `bson:"student_metadata,omitempty" json:"student_metadata,omitempty"`

	ExtendedProperties   map[string]interface{} `bson:"extended_properties,omitempty" json:"extended_properties,omitempty"`
	PsychologicalProfile *PsychologicalProfile  `bson:"psychological_profile,omitempty" json:"psychological_profile,omitempty"`
	BatchID              string                 `bson:"batch_id,omitempty" json:"batch_id,omitempty"`
}

// Condition represents a prerequisite for a relationship to be valid.
type Condition struct {
	ConditionType  string `bson:"condition_type" json:"condition_type"`
	ConditionValue string `bson:"condition_value" json:"condition_value"`
	Required       bool   `bson:"required" json:"required"`
	Description    string `bson:"description,omitempty" json:"description,omitempty"`
}

// Contrast explicitly states how a relationship differs from similar ones.
type Contrast struct {
	Term       string `bson:"term" json:"term"`
	Difference string `bson:"difference" json:"difference"`
	Rationale  string `bson:"rationale,omitempty" json:"rationale,omitempty"`
}

// PatientMetadata stores patient-specific information.
type PatientMetadata struct {
	PatientName      string `bson:"patient_name,omitempty" json:"patient_name,omitempty"`
	PatientInitials  string `bson:"patient_initials,omitempty" json:"patient_initials,omitempty"`
	DiseaseCategory  string `bson:"disease_category" json:"disease_category"`
	CaseNumber       int    `bson:"case_number" json:"case_number"`
	ChapterReference string `bson:"chapter_reference,omitempty" json:"chapter_reference,omitempty"`
}

// StudentMetadata stores student-specific information.
type StudentMetadata struct {
	StudentID            string                `bson:"student_id,omitempty" json:"student_id,omitempty"`
	StudentName          string                `bson:"student_name,omitempty" json:"student_name,omitempty"`
	Age                  int                   `bson:"age,omitempty" json:"age,omitempty"`
	YearGroup            string                `bson:"year_group,omitempty" json:"year_group,omitempty"`
	SpecialNeeds         []string              `bson:"special_needs,omitempty" json:"special_needs,omitempty"`
	PsychologicalProfile *PsychologicalProfile `bson:"psychological_profile,omitempty" json:"psychological_profile,omitempty"`
}

// PsychologicalProfile contains detailed student profiling information.
type PsychologicalProfile struct {
	LearningStyle     string   `bson:"learning_style,omitempty" json:"learning_style,omitempty"`
	AttentionSpan     int      `bson:"attention_span,omitempty" json:"attention_span,omitempty"` // in minutes
	MotivationFactors []string `bson:"motivation_factors,omitempty" json:"motivation_factors,omitempty"`
	AnxietyTriggers   []string `bson:"anxiety_triggers,omitempty" json:"anxiety_triggers,omitempty"`
	Strengths         []string `bson:"strengths,omitempty" json:"strengths,omitempty"`
	Challenges        []string `bson:"challenges,omitempty" json:"challenges,omitempty"`
}
