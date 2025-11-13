package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChapterInfo represents metadata about a chapter being uploaded or processed
type ChapterInfo struct {
	ChapterNumber string `bson:"chapter_number,omitempty" json:"chapter_number,omitempty"`
	ChapterTitle  string `bson:"chapter_title,omitempty" json:"chapter_title,omitempty"`
	Domain        string `bson:"domain,omitempty" json:"domain,omitempty"`
}

// Section represents a section in the chapter
type Section struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// DataSource represents the original source of content (files, uploads, web scrapes)
type DataSource struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Type        string                 `bson:"type" json:"type"`
	Source      string                 `bson:"source" json:"source"`
	ProcessedAt int64                  `bson:"processed_at" json:"processed_at"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// PageRange represents the page range for content
type PageRange struct {
	StartPage int `bson:"start_page" json:"start_page"`
	EndPage   int `bson:"end_page" json:"end_page"`
}

// VectorSource represents the source attribution for vector embeddings
type VectorSource struct {
	SourceID         primitive.ObjectID `bson:"source_id" json:"source_id"`
	DataSourceID     primitive.ObjectID `bson:"data_source_id" json:"data_source_id"`
	SourceType       string             `bson:"source_type" json:"source_type"`
	ChapterInfo      *ChapterInfo       `bson:"chapter_info,omitempty" json:"chapter_info,omitempty"`
	PageRange        *PageRange         `bson:"page_range,omitempty" json:"page_range,omitempty"`
	ExtractedAt      time.Time          `bson:"extracted_at" json:"extracted_at"`
	ExtractionMethod string             `bson:"extraction_method" json:"extraction_method"`
	BatchID          string             `bson:"batch_id,omitempty" json:"batch_id,omitempty"`
}

// Chunk structure with proper source tracking
type Chunk struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Content      string                 `bson:"content" json:"content"`
	VectorSource VectorSource           `bson:"vector_source" json:"vector_source"`
	ContentType  string                 `bson:"content_type" json:"content_type"`
	CreatedAt    int64                  `bson:"created_at" json:"created_at"`
	Metadata     map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// Material represents processed content
type Material struct {
	ID              primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Title           string               `json:"title" bson:"title"`
	Description     string               `json:"description,omitempty" bson:"description,omitempty"`
	OriginalContent string               `json:"original_content,omitempty" bson:"original_content,omitempty"`
	DataSourceID    primitive.ObjectID   `json:"data_source_id" bson:"data_source_id"`
	Chunks          []primitive.ObjectID `json:"chunks,omitempty" bson:"chunks,omitempty"`
	SemanticLinks   []primitive.ObjectID `json:"semantic_links,omitempty" bson:"semantic_links,omitempty"`
	CreatedAt       int64                `json:"created_at" bson:"created_at"`
	UpdatedAt       int64                `json:"updated_at" bson:"updated_at"`
}
