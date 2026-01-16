package constants

// UniversalTags are topic-independent tags that apply to any subject
var UniversalTags = []string{
	// Study/Revision tags
	"revision",
	"exam_prep",
	"practice_question",
	"key_concept",
	"definition",
	"summary",
	"common_mistake",
	"tip",
	"formula",
	"diagram",

	// Content-specific tags
	"case_study",
	"clinical_example",
	"design_example",
	"worked_example",
	"concept_map",

	// Assessment tags
	"past_exam_question",
	"mock_exam",
	"homework",
	"problem_set",

	// Difficulty indicators
	"foundation",
	"core_knowledge",
	"advanced",

	// Learning style
	"visual",
	"textual",
	"interactive",
	"audio",
}

// SubjectSpecificTagExamples shows example tags for different subjects
func SubjectSpecificTagExamples() map[string][]string {
	return map[string][]string{
		"immunology": {
			"immune_response",
			"antibody",
			"antigen",
			"infection",
			"clinical_findings",
			"pathophysiology",
			"diagnosis",
			"treatment",
			"patient_management",
		},
		"civil_engineering": {
			"structural_analysis",
			"materials",
			"design",
			"construction",
			"calculations",
			"failure_modes",
			"load_bearing",
			"stress_strain",
			"design_process",
			"safety_factors",
		},
		"gcse": {
			"biology",
			"chemistry",
			"physics",
			"revision_summary",
			"equation",
			"calculation",
			"practical",
			"past_paper",
			"specification",
		},
	}
}

// IsValidUniversalTag checks if a tag is in the universal tags list
func IsValidUniversalTag(tag string) bool {
	for _, valid := range UniversalTags {
		if valid == tag {
			return true
		}
	}
	return false
}
