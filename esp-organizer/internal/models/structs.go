package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RelatedSkill struct {
	ID   primitive.ObjectID `bson:"id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

// CreateSkillRequest defines the payload for creating a new skill
type CreateSkillRequest struct {
	Name           string               `json:"name" binding:"required"`
	Type           string               `json:"type" binding:"required"`
	Description    string               `json:"description" binding:"required"`
	DevelopmentAge int                  `json:"developmentAge" binding:"required"`
	Hidden         bool                 `json:"hidden"`
	Criteria       []SkillCriteria      `json:"criteria" binding:"required,dive"`
	ParentSkillIDs []primitive.ObjectID `json:"parentSkillIDs,omitempty"`
	SourceRefs     []primitive.ObjectID `json:"sourceRefs,omitempty"`
}

// UpdateSkillRequest defines the payload for updating an existing skill
type UpdateSkillRequest struct {
	Name           *string              `json:"name,omitempty"`
	Type           *string              `json:"type,omitempty"`
	Description    *string              `json:"description,omitempty"`
	DevelopmentAge *int                 `json:"developmentAge,omitempty"`
	Hidden         *bool                `json:"hidden,omitempty"`
	Criteria       *[]SkillCriteria     `json:"criteria,omitempty"`
	ParentSkillIDs []primitive.ObjectID `json:"parentSkillIDs,omitempty"`
	SourceRefs     []primitive.ObjectID `json:"sourceRefs,omitempty"`
}

// --- Request/Response DTOs for Source CRUD ---

// CreateSourceRequest defines the payload for creating a new source
type CreateSourceRequest struct {
	Name            string `json:"name" binding:"required"`
	Type            string `json:"type" binding:"required"`
	Author          string `json:"author,omitempty"`
	PublicationYear int    `json:"publicationYear,omitempty"`
	URL             string `json:"url,omitempty"`
	Description     string `json:"description,omitempty"`
	ISBN            string `json:"isbn,omitempty"`
}

// UpdateSourceRequest defines the payload for updating an existing source
type UpdateSourceRequest struct {
	Name            *string `json:"name,omitempty"`
	Type            *string `json:"type,omitempty"`
	Author          *string `json:"author,omitempty"`
	PublicationYear *int    `json:"publicationYear,omitempty"`
	URL             *string `json:"url,omitempty"`
	Description     *string `json:"description,omitempty"`
	ISBN            *string `json:"isbn,omitempty"`
}

// --- Weaviate Data Structures ---

// WeaviateSkill represents the data structure for a skill in Weaviate
type WeaviateSkill struct {
	MongoID        string   `json:"mongo_id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	SkillType      string   `json:"skill_type"`
	DevelopmentAge int      `json:"development_age"`
	CriteriaText   string   `json:"criteria_text"`
	ParentMongoIDs []string `json:"parent_mongo_ids,omitempty"`
	SourceMongoIDs []string `json:"source_mongo_ids,omitempty"`
}

// WeaviateSource represents the data structure for a source in Weaviate
type WeaviateSource struct {
	MongoID         string `json:"mongo_id"`
	Name            string `json:"name"`
	SourceType      string `json:"source_type"`
	Author          string `json:"author,omitempty"`
	PublicationYear int    `json:"publication_year,omitempty"`
	Description     string `json:"description,omitempty"`
	ContentSummary  string `json:"content_summary,omitempty"`
}
