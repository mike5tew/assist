package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaperReference is a single entry from a paper's reference list.
// ResolvedID is populated when the cited paper is also present in the corpus,
// enabling a citation graph across ingested documents.
type PaperReference struct {
	FullCitation string              `bson:"full_citation" json:"full_citation"`
	DOI          string              `bson:"doi,omitempty" json:"doi,omitempty"`
	ResolvedID   *primitive.ObjectID `bson:"resolved_id,omitempty" json:"resolved_id,omitempty"`
}

// PaperFigure captures a figure or table from the paper.
type PaperFigure struct {
	Label   string `bson:"label" json:"label"`
	Caption string `bson:"caption" json:"caption"`
	Page    int    `bson:"page,omitempty" json:"page,omitempty"`
}

// SourceDocument represents a scientific paper ingested into CHISG.
// Semantic links store source_document_id + denormalised source_title and doi
// so Weaviate can filter by paper without a join, while this collection
// remains the single source of truth for full citation detail.
//
// MongoDB collection: chisg_knowledge_base.source_documents
type SourceDocument struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Core citation fields
	Title    string   `bson:"title" json:"title"`
	Authors  []string `bson:"authors" json:"authors"`
	DOI      string   `bson:"doi,omitempty" json:"doi,omitempty"`
	Year     int      `bson:"year,omitempty" json:"year,omitempty"`
	Journal  string   `bson:"journal,omitempty" json:"journal,omitempty"`
	Abstract string   `bson:"abstract,omitempty" json:"abstract,omitempty"`

	// Extracted structure (populated during / after Textract pipeline)
	Figures        []PaperFigure    `bson:"figures,omitempty" json:"figures,omitempty"`
	MethodsSummary string           `bson:"methods_summary,omitempty" json:"methods_summary,omitempty"`
	References     []PaperReference `bson:"references,omitempty" json:"references,omitempty"`

	// InCorpusCitationCount is incremented each time another ingested paper
	// is found to cite this document (matched by DOI).  Low count ≠ unimportant,
	// but high count does signal a foundational paper for the knowledge graph.
	InCorpusCitationCount int `bson:"in_corpus_citation_count" json:"in_corpus_citation_count"`

	// OriginalFilename is the uploaded file name, used as a fallback display title.
	OriginalFilename string `bson:"original_filename" json:"original_filename"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
