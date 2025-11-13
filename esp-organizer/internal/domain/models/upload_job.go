package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UploadJob represents a file upload and processing job
type UploadJob struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// File metadata
	OriginalFilename string `bson:"original_filename" json:"original_filename"`
	FileSize         int64  `bson:"file_size" json:"file_size"`
	ContentType      string `bson:"content_type" json:"content_type"`
	S3Key            string `bson:"s3_key,omitempty" json:"s3_key,omitempty"` // Where the file is stored

	// Source linkage
	SourceID      primitive.ObjectID `bson:"source_id" json:"source_id"`           // Links to Source collection
	ChapterNumber string             `bson:"chapter_number" json:"chapter_number"` // Which chapter
	ChapterTitle  string             `bson:"chapter_title" json:"chapter_title"`

	// Processing status
	Status       string  `bson:"status" json:"status"` // pending, processing, completed, failed
	Progress     float64 `bson:"progress" json:"progress"`
	StatusDetail string  `bson:"status_detail,omitempty" json:"status_detail,omitempty"`

	// Processing details
	TextractJobID string `bson:"textract_job_id,omitempty" json:"textract_job_id,omitempty"`

	// Results
	ExtractedContentIDs []primitive.ObjectID `bson:"extracted_content_ids,omitempty" json:"extracted_content_ids,omitempty"` // Links to immunology_content
	ErrorMessage        string               `bson:"error_message,omitempty" json:"error_message,omitempty"`

	// Timestamps
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	StartedAt   *time.Time `bson:"started_at,omitempty" json:"started_at,omitempty"`
	CompletedAt *time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
}
