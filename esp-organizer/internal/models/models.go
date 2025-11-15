package models

import (
	"context"
	"time"
)

// ExtractionJob tracks the status of a document extraction job.
type ExtractionJob struct {
	ID            string         `json:"id" bson:"_id"`
	Status        string         `json:"status" bson:"status"` // e.g., initializing, processing, completed, failed
	Message       string         `json:"message,omitempty" bson:"message,omitempty"`
	Progress      float64        `json:"progress" bson:"progress"` // 0.0 to 1.0
	CreatedAt     time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" bson:"updated_at"`
	ExtractedData *ExtractedData `json:"extracted_data,omitempty" bson:"extracted_data,omitempty"`
}

// ExtractedData holds the parsed content from a document.
type ExtractedData struct {
	RawText        string        `json:"raw_text" bson:"raw_text"`
	ChapterContent string        `json:"chapter_content" bson:"chapter_content"`
	CaseStudies    []CaseStudy   `json:"case_studies" bson:"case_studies"`
	MedicalTerms   []MedicalTerm `json:"medical_terms" bson:"medical_terms"`
}

// SourceInfo contains metadata about the source of the content.
type SourceInfo struct {
	Type      string `bson:"type" json:"type"`           // e.g., "book", "article", "pdf"
	UploadID  string `bson:"upload_id" json:"upload_id"` // Batch ID for the upload
	Author    string `bson:"author,omitempty" json:"author,omitempty"`
	Reference string `bson:"reference,omitempty" json:"reference,omitempty"`
}

// QueryResponse represents the response to a query
type QueryResponse struct {
	Query           string                   `json:"query"`
	Timestamp       time.Time                `json:"timestamp"`
	MongoResults    []Skill                  `json:"mongo_results,omitempty"`
	WeaviateResults []map[string]interface{} `json:"weaviate_results,omitempty"`
	Synthesis       string                   `json:"synthesis,omitempty"`
	RelatedSkills   []Skill                  `json:"related_skills,omitempty"`
}

type QueryProcessor interface {
	ProcessQuery(ctx context.Context, query string, filters *Filter, weaviateClass string) (*QueryResponse, error)
}

type Filter struct {
	Operator     string    `json:"operator,omitempty"`     // "And" or "Or"
	Operands     []Operand `json:"operands,omitempty"`     // List of operands in this filter
	NestedFilter *Filter   `json:"nestedFilter,omitempty"` // Support for nested filters
}
type Operand struct {
	Path         []string `json:"path,omitempty"`         // Property path to filter on
	Operator     string   `json:"operator,omitempty"`     // Operator like "Equal", "GreaterThan", etc.
	ValueString  string   `json:"valueString,omitempty"`  // String value for the filter
	ValueNumber  float64  `json:"valueNumber,omitempty"`  // Numeric value for the filter
	ValueBoolean bool     `json:"valueBoolean,omitempty"` // Boolean value for the filter
	NestedFilter *Filter  `json:"nestedFilter,omitempty"` // Support for nested filters in operands
}
