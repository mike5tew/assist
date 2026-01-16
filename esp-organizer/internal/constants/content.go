package constants

// ContentType represents the category/structure of a piece of content
type ContentType string

const (
	ContentTypeChapter          ContentType = "chapter"
	ContentTypeCaseStudy        ContentType = "case_study"
	ContentTypeDefinition       ContentType = "definition"
	ContentTypeSection          ContentType = "section"
	ContentTypeDesignExample    ContentType = "design_example"
	ContentTypePracticeQuestion ContentType = "practice_question"
	ContentTypeConceptMap       ContentType = "concept_map"
	ContentTypeKeyConcept       ContentType = "key_concept"
	ContentTypeSummary          ContentType = "summary"
	ContentTypeFormula          ContentType = "formula"
	ContentTypeDiagram          ContentType = "diagram"
)

// AllContentTypes returns all valid content types
func AllContentTypes() []ContentType {
	return []ContentType{
		ContentTypeChapter,
		ContentTypeCaseStudy,
		ContentTypeDefinition,
		ContentTypeSection,
		ContentTypeDesignExample,
		ContentTypePracticeQuestion,
		ContentTypeConceptMap,
		ContentTypeKeyConcept,
		ContentTypeSummary,
		ContentTypeFormula,
		ContentTypeDiagram,
	}
}

// IsValidContentType checks if the provided type is valid
func IsValidContentType(ct string) bool {
	for _, valid := range AllContentTypes() {
		if string(valid) == ct {
			return true
		}
	}
	return false
}

// SourceType represents the original format/source of content
type SourceType string

const (
	SourceTypeBook      SourceType = "book"
	SourceTypeTextbook  SourceType = "textbook"
	SourceTypeArticle   SourceType = "article"
	SourceTypePaper     SourceType = "paper"
	SourceTypeWebpage   SourceType = "web"
	SourceTypePDF       SourceType = "pdf"
	SourceTypeUpload    SourceType = "upload"
	SourceTypeCourse    SourceType = "course"
	SourceTypeExamBoard SourceType = "exam_board"
)

// AllSourceTypes returns all valid source types
func AllSourceTypes() []SourceType {
	return []SourceType{
		SourceTypeBook,
		SourceTypeTextbook,
		SourceTypeArticle,
		SourceTypePaper,
		SourceTypeWebpage,
		SourceTypePDF,
		SourceTypeUpload,
		SourceTypeCourse,
		SourceTypeExamBoard,
	}
}

// IsValidSourceType checks if the provided type is valid
func IsValidSourceType(st string) bool {
	for _, valid := range AllSourceTypes() {
		if string(valid) == st {
			return true
		}
	}
	return false
}
