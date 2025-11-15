package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcessingJob struct {
	ID            string                `json:"id"`
	Status        string                `json:"status"` // "processing", "completed", "failed"
	ChapterInfo   ChapterInfo           `json:"chapter_info"`
	BookSource    Source                `json:"book_source"`
	CreatedAt     time.Time             `json:"created_at"`
	CompletedAt   *time.Time            `json:"completed_at,omitempty"`
	ErrorMessage  string                `json:"error_message,omitempty"`
	ExtractedData *ExtractedChapterData `json:"extracted_data,omitempty"`
}

// ExtractedChapterData holds the data extracted from a chapter
type ExtractedChapterData struct {
	CaseStudies    []CaseStudy   `json:"case_studies"`
	MedicalTerms   []MedicalTerm `json:"medical_terms"`
	Sections       []Section     `json:"sections"`
	ChapterContent string        `json:"chapter_content"` // Full chapter content as a single string
	RawText        string        `json:"raw_text"`
	TextractBlocks int           `json:"textract_blocks"`
}

// PatientInfo represents patient information in a case study
type PatientInfo struct {
	Age            string `json:"age"`
	Gender         string `json:"gender"`
	MedicalHistory string `json:"medical_history"`
}

type SubjectContent struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Content     string                 `bson:"content" json:"content"`                   // The actual text content
	ContentType string                 `bson:"content_type" json:"content_type"`         // "chapter", "case_study", "definition", etc.
	Domain      string                 `bson:"domain" json:"domain"`                     // "immunology", "gcse", "civil-engineering"
	Title       string                 `bson:"title" json:"title"`                       // Title/name of the content
	Source      SourceReference        `bson:"source" json:"source"`                     // Updated to use SourceReference
	Tags        []string               `bson:"tags" json:"tags"`                         // Domain-specific tags
	Metadata    map[string]interface{} `bson:"metadata" json:"metadata"`                 // Domain-specific extended metadata
	Vector      []float32              `bson:"vector,omitempty" json:"vector,omitempty"` // Embedding vector
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updated_at"`
}

// Source represents any content source (book, document, web page, etc.)
type Source struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Type      string                 `json:"type" bson:"type"` // "book", "article", "web", "upload"
	Title     string                 `json:"title" bson:"title"`
	Authors   []string               `json:"authors" bson:"authors,omitempty"`
	Publisher string                 `json:"publisher" bson:"publisher,omitempty"`
	Year      string                 `json:"year" bson:"year,omitempty"`
	ISBN      string                 `json:"isbn" bson:"isbn,omitempty"`
	URL       string                 `json:"url" bson:"url,omitempty"`
	Domain    string                 `json:"domain" bson:"domain,omitempty"` // e.g., "immunology", "gcse"
	FilePath  string                 `json:"file_path" bson:"file_path,omitempty"`
	Processed bool                   `json:"processed" bson:"processed"`
	Metadata  map[string]interface{} `json:"metadata" bson:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" bson:"updated_at"`
}

// Book is now an alias for Source with book-specific methods
type Book = Source

// MedicalExcerpt now references the unified Source model
type MedicalExcerpt struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	SourceID  string    `json:"source_id" bson:"source_id"` // References Source.ID
	Content   string    `json:"content" bson:"content"`
	Page      int       `json:"page" bson:"page"`
	Chapter   int       `json:"chapter" bson:"chapter,omitempty"`
	Keywords  []string  `json:"keywords" bson:"keywords,omitempty"`
	Processed bool      `json:"processed" bson:"processed"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// SourceReference is a lightweight reference to a source with specific location info
type SourceReference struct {
	SourceID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `json:"title" bson:"title"` // Cache of source title
	ChapterNumber string             `json:"chapter_number,omitempty" bson:"chapter_number,omitempty"`
	ChapterTitle  string             `json:"chapter_title,omitempty" bson:"chapter_title,omitempty"`
	PageNumbers   string             `json:"page_numbers,omitempty" bson:"page_numbers,omitempty"`
	UploadID      string             `json:"upload_id,omitempty" bson:"upload_id,omitempty"`
	Authors       []string           `json:"authors,omitempty" bson:"authors,omitempty"`
}

// Term represents a verified term from Weaviate search results
type Term struct {
	Text     string                 `json:"text"`
	Verified bool                   `json:"verified"`
	MongoID  primitive.ObjectID     `json:"mongo_id"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type MedicalTerm struct {
	Term       string   `json:"term"`
	Definition string   `json:"definition,omitempty"`
	Context    []string `json:"context,omitempty"`
	Category   string   `json:"category,omitempty"` // e.g., "Cell Type", "Cytokine", "Disorder"
	Tags       []string `json:"tags,omitempty"`
}

type CaseStudy struct {
	CaseNumber       string   `json:"case_number"`
	Title            string   `json:"title"`
	Content          string   `json:"content"`
	ClinicalFindings string   `json:"clinical_findings,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

type Section struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Level   int    `json:"level"` // e.g., 1 for chapter title, 2 for section, 3 for subsection
}

type ChapterInfo struct {
	ChapterNumber string `json:"chapter_number"`
	ChapterTitle  string `json:"chapter_title"`
}
