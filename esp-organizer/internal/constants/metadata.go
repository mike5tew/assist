package constants

// RevisionMetadata is the standard metadata structure for all revision content
type RevisionMetadata struct {
	Difficulty      string                 `json:"difficulty"`
	AcademicLevel   string                 `json:"academic_level"`
	EstimatedTimeMs int                    `json:"estimated_time_ms"`
	ContentType     string                 `json:"content_type"`
	RelatedConcepts []string               `json:"related_concepts"`
	Prerequisites   []string               `json:"prerequisites"`
	Keywords        []string               `json:"keywords"`
	LastReviewedAt  string                 `json:"last_reviewed_at,omitempty"`
	ReviewCount     int                    `json:"review_count,omitempty"`
	StudentNotes    string                 `json:"student_notes,omitempty"`
	IsVerified      bool                   `json:"is_verified"`
	VerifiedBy      string                 `json:"verified_by,omitempty"`
	QualityScore    float64                `json:"quality_score"`
	Custom          map[string]interface{} `json:"custom"`
}

// MetadataTemplate provides a base template for creating metadata
func MetadataTemplate(difficulty, level, contentType string, estimatedMs int) RevisionMetadata {
	return RevisionMetadata{
		Difficulty:      difficulty,
		AcademicLevel:   level,
		EstimatedTimeMs: estimatedMs,
		ContentType:     contentType,
		RelatedConcepts: []string{},
		Prerequisites:   []string{},
		Keywords:        []string{},
		ReviewCount:     0,
		IsVerified:      false,
		QualityScore:    0.5,
		Custom:          make(map[string]interface{}),
	}
}

// SubjectMetadataExtensions define additional metadata fields per subject
var SubjectMetadataExtensions = map[string][]string{
	"immunology": {
		"case_number",
		"disease_category",
		"clinical_findings",
		"diagnosis",
		"treatment",
		"patient_initials",
		"age_at_diagnosis",
		"prognosis",
	},
	"civil_engineering": {
		"calculation_type",
		"unit_system",
		"material_type",
		"design_code",
		"failure_modes",
		"safety_factor",
		"construction_method",
		"cost_estimate",
	},
	"gcse": {
		"exam_board",
		"spec_section",
		"command_word",
		"mark_allocation",
		"suggested_time",
		"past_paper_ref",
		"grade_range",
	},
}

// MetadataValidation rules for ensuring consistency
type MetadataValidation struct {
	AllowEmptyPrerequisites  bool
	AllowEmptyRelatedContent bool
	RequireVerification      bool
	MinQualityScore          float64
	MaxCustomFields          int
}

// DefaultMetadataValidation returns standard validation rules
func DefaultMetadataValidation() MetadataValidation {
	return MetadataValidation{
		AllowEmptyPrerequisites:  true,
		AllowEmptyRelatedContent: true,
		RequireVerification:      false,
		MinQualityScore:          0.3,
		MaxCustomFields:          20,
	}
}
