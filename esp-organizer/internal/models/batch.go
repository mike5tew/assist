package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProcessingBatch represents a batch of related operations that can be rolled back together
type ProcessingBatch struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	BatchID        string             `bson:"batch_id" json:"batch_id"`
	ProcessingType string             `bson:"processing_type" json:"processing_type"` // "textract_chapter", "manual_upload", "bulk_import"
	Status         string             `bson:"status" json:"status"`                   // "completed", "rolled_back", "partial"
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	RolledBackAt   *time.Time         `bson:"rolled_back_at,omitempty" json:"rolled_back_at,omitempty"`

	// Track what was created in this batch
	CreatedRecords BatchRecordSummary `bson:"created_records" json:"created_records"`

	// Metadata about the batch
	SourceInfo  map[string]interface{} `bson:"source_info,omitempty" json:"source_info,omitempty"`
	UserID      string                 `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Description string                 `bson:"description,omitempty" json:"description,omitempty"`
}

// BatchRecordSummary tracks what records were created in a batch
type BatchRecordSummary struct {
	ChapterIDs     []primitive.ObjectID `bson:"chapter_ids,omitempty" json:"chapter_ids,omitempty"`
	CaseStudyIDs   []primitive.ObjectID `bson:"case_study_ids,omitempty" json:"case_study_ids,omitempty"`
	MedicalTermIDs []primitive.ObjectID `bson:"medical_term_ids,omitempty" json:"medical_term_ids,omitempty"`
	SkillIDs       []primitive.ObjectID `bson:"skill_ids,omitempty" json:"skill_ids,omitempty"`
	WeaviateIDs    []string             `bson:"weaviate_ids,omitempty" json:"weaviate_ids,omitempty"`

	// Summary counts
	TotalRecords int `bson:"total_records" json:"total_records"`
}
