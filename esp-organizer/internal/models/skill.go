package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectIDHolder handles MongoDB Extended JSON format for ObjectIDs
type ObjectIDHolder struct {
	Oid string `json:"$oid"`
}

// SkillCriteria represents the criteria for assessing a skill
type SkillCriteria struct {
	Level       int    `bson:"level" json:"level"`
	Description string `bson:"description" json:"description"`
}

// Skill represents a unified model for educational skills.
type Skill struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	IDString       string             `bson:"-" json:"_id"` // Temporary field for JSON unmarshalling
	Name           string             `bson:"name" json:"name"`
	Description    string             `bson:"description" json:"description"`
	Category       string             `bson:"category" json:"category"`
	DevelopmentAge int                `bson:"developmentAge" json:"developmentAge"`

	// Unified field for parent skill IDs
	ParentSkillIDs []primitive.ObjectID `bson:"parentSkillIDs,omitempty" json:"parentSkillIDs,omitempty"`

	// Populated fields for relationships
	ParentSkills  []*Skill             `bson:"-" json:"parentSkills,omitempty"`
	ChildSkillIDs []primitive.ObjectID `bson:"childSkillIDs,omitempty" json:"childSkillIDs,omitempty"`
	SourceRefs    []SourceReference    `bson:"sourceRefs,omitempty" json:"sourceRefs,omitempty"`
	LastUpdated   time.Time            `bson:"lastUpdated,omitempty" json:"lastUpdated,omitempty"`
	SkillCriteria []SkillCriteria      `bson:"criteria" json:"criteria"`

	// Fields for medical content and source tracking
	SourceTitle      string `bson:"source_title,omitempty" json:"source_title,omitempty"`
	SourceType       string `bson:"source_type,omitempty" json:"source_type,omitempty"`
	ExtractionMethod string `bson:"extraction_method,omitempty" json:"extraction_method,omitempty"`
	ChapterTitle     string `bson:"chapter_title,omitempty" json:"chapter_title,omitempty"`

	// Populated fields for relationships (not stored in MongoDB)
	Sources []*Source `bson:"-" json:"sources,omitempty"`
}

// EducationalSkill represents the structure of a skill in the skills.json file.
// This is now DEPRECATED in favor of the unified `Skill` model.
type EducationalSkill struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name           string               `bson:"name" json:"name"`
	Description    string               `bson:"description" json:"description"`
	SkillType      string               `bson:"skill_type" json:"skill_type"` // Corrected from 'Type' to match DB
	DevelopmentAge int                  `bson:"development_age" json:"development_age"`
	Criteria       []SkillCriteria      `bson:"criteria" json:"criteria"`
	ParentSkillIDs []primitive.ObjectID `bson:"parent_skill_ids,omitempty" json:"parent_skill_ids,omitempty"`
	SourceRefs     []primitive.ObjectID `bson:"source_refs,omitempty" json:"source_refs,omitempty"`
	Hidden         bool                 `bson:"hidden" json:"hidden"`

	// Fields for medical content and source tracking
	SourceTitle      string `bson:"source_title,omitempty" json:"source_title,omitempty"`
	SourceType       string `bson:"source_type,omitempty" json:"source_type,omitempty"`
	ExtractionMethod string `bson:"extraction_method,omitempty" json:"extraction_method,omitempty"`
	ChapterTitle     string `bson:"chapter_title,omitempty" json:"chapter_title,omitempty"`

	// Populated fields for relationships (not stored in MongoDB)
	ParentSkills []*Skill  `bson:"-" json:"parent_skills,omitempty"`
	Sources      []*Source `bson:"-" json:"sources,omitempty"`
}
