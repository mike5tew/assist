package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SemanticLink struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// 🆕 Composite key for deduplication
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

	// 🆕 Store normalized context for vectorization, original for display
	Context         string               `bson:"context" json:"context"`                                       // Normalized
	ContextOriginal string               `bson:"context_original,omitempty" json:"context_original,omitempty"` // 🆕 With symbols
	IsParentOf      []primitive.ObjectID `bson:"is_parent_of,omitempty" json:"is_parent_of,omitempty"`
	IsChildOf       []primitive.ObjectID `bson:"is_child_of,omitempty" json:"is_child_of,omitempty"`

	Confidence float64                `bson:"confidence" json:"confidence"`
	Domain     string                 `bson:"domain" json:"domain"`                       // e.g., "immunology", "personal", "neurology"
	Subject    string                 `bson:"subject,omitempty" json:"subject,omitempty"` // e.g., "esp_organizer_development", "go_programming"
	Vector     []float32              `bson:"vector,omitempty" json:"vector,omitempty"`
	CreatedAt  time.Time              `bson:"created_at" json:"created_at"`
	Metadata   map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`

	// 🆕 Store BOTH original and normalized versions
	SourceTermOriginal string `bson:"source_term_original,omitempty" json:"source_term_original,omitempty"` // 🆕 With symbols
	TargetTermOriginal string `bson:"target_term_original,omitempty" json:"target_term_original,omitempty"` // 🆕 With symbols

	// 🆕 CONTEXTUAL QUALIFIERS
	TemporalQualifier string `bson:"temporal_qualifier,omitempty" json:"temporal_qualifier,omitempty"` // permanent | transient | conditional | developmental
	SpatialQualifier  string `bson:"spatial_qualifier,omitempty" json:"spatial_qualifier,omitempty"`   // intracellular | surface | secreted | circulating | nuclear | cytoplasmic
	FunctionalState   string `bson:"functional_state,omitempty" json:"functional_state,omitempty"`     // immature | mature | activated | resting | apoptotic | proliferating
	LifecycleStage    string `bson:"lifecycle_stage,omitempty" json:"lifecycle_stage,omitempty"`       // Specific developmental or disease stage

	// 🆕 CONDITIONAL LOGIC
	Conditions []Condition `bson:"conditions,omitempty" json:"conditions,omitempty"`

	// 🆕 CONTRASTS (what distinguishes this from similar relationships)
	ContrastsWith []Contrast `bson:"contrasts_with,omitempty" json:"contrasts_with,omitempty"`

	// 🆕 CONTEXT CLASSIFICATION
	LinkType string `bson:"link_type" json:"link_type"` // "case_specific" | "reference_value" | "clinical_interpretation" | "general_fact"

	// 🆕 CROSS-REFERENCES (for connecting case-specific to reference links)
	RelatedLinkIDs []string `bson:"related_link_ids,omitempty" json:"related_link_ids,omitempty"` // Links this connects to

	// 🆕 CASE-SPECIFIC FLAG (derived from * marker)
	IsCaseSpecific bool   `bson:"is_case_specific" json:"is_case_specific"`               // True if extracted from *-marked entity
	CaseStudyID    string `bson:"case_study_id,omitempty" json:"case_study_id,omitempty"` // e.g., "XLA_Case_1" (NOT "Bill_Grignard_XLA")

	// 🆕 PATIENT METADATA (stored separately for privacy/flexibility)
	PatientMetadata *PatientMetadata `bson:"patient_metadata,omitempty" json:"patient_metadata,omitempty"`

	// 🆕 Optional: Add student-specific metadata
	StudentMetadata *StudentMetadata `bson:"student_metadata,omitempty" json:"student_metadata,omitempty"`

	// 🆕 ADD THIS: Allows storing additional fields without schema changes
	ExtendedProperties map[string]interface{} `bson:"extended_properties,omitempty" json:"extended_properties,omitempty"`

	// 🆕 ADDED (optional field, old links work without it)
	PsychologicalProfile *PsychologicalProfile `bson:"psychological_profile,omitempty" json:"psychological_profile,omitempty"`
	BatchID              string                `bson:"batch_id,omitempty" json:"batch_id,omitempty"` // For tracking ingestion batches
}

// 🆕 Condition represents a prerequisite for this relationship to be valid
type Condition struct {
	ConditionType  string `bson:"condition_type" json:"condition_type"`               // e.g., "developmental_stage", "gene_rearrangement", "activation_state"
	ConditionValue string `bson:"condition_value" json:"condition_value"`             // e.g., "large_pre-B", "VDJ_rearranged"
	Required       bool   `bson:"required" json:"required"`                           // Must be true for relationship to apply
	Description    string `bson:"description,omitempty" json:"description,omitempty"` // Human-readable explanation
}

// 🆕 Contrast explicitly states how this relationship differs from similar ones
type Contrast struct {
	Term       string `bson:"term" json:"term"`                               // The contrasting term (e.g., "Mature B cell")
	Difference string `bson:"difference" json:"difference"`                   // What's different
	Rationale  string `bson:"rationale,omitempty" json:"rationale,omitempty"` // Why this matters
}

// 🆕 PatientMetadata stores patient-specific information separately
type PatientMetadata struct {
	PatientName      string `bson:"patient_name,omitempty" json:"patient_name,omitempty"`           // "Bill Grignard"
	PatientInitials  string `bson:"patient_initials,omitempty" json:"patient_initials,omitempty"`   // "B.G."
	DiseaseCategory  string `bson:"disease_category" json:"disease_category"`                       // "XLA", "SCID", etc.
	CaseNumber       int    `bson:"case_number" json:"case_number"`                                 // 1, 2, 3...
	ChapterReference string `bson:"chapter_reference,omitempty" json:"chapter_reference,omitempty"` // "Chapter 1"
}

// 🆕 StudentMetadata stores student-specific information
type StudentMetadata struct {
	StudentID            string                `bson:"student_id,omitempty" json:"student_id,omitempty"` // e.g., "STUDENT_OW_2015"
	StudentName          string                `bson:"student_name,omitempty" json:"student_name,omitempty"`
	Age                  int                   `bson:"age,omitempty" json:"age,omitempty"`
	YearGroup            string                `bson:"year_group,omitempty" json:"year_group,omitempty"`
	SpecialNeeds         []string              `bson:"special_needs,omitempty" json:"special_needs,omitempty"` // ["autism", "high_IQ"]
	PsychologicalProfile *PsychologicalProfile `bson:"psychological_profile,omitempty" json:"psychological_profile,omitempty"`
}

// 🆕 PsychologicalProfile for detailed student profiling
type PsychologicalProfile struct {
	LearningStyle     string   `bson:"learning_style,omitempty" json:"learning_style,omitempty"`
	AttentionSpan     int      `bson:"attention_span,omitempty" json:"attention_span,omitempty"` // minutes
	MotivationFactors []string `bson:"motivation_factors,omitempty" json:"motivation_factors,omitempty"`
	AnxietyTriggers   []string `bson:"anxiety_triggers,omitempty" json:"anxiety_triggers,omitempty"`
	Strengths         []string `bson:"strengths,omitempty" json:"strengths,omitempty"`
	Challenges        []string `bson:"challenges,omitempty" json:"challenges,omitempty"`
}
