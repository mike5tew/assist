package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MisunderstandingMap represents a documented pattern of student misconceptions
type MisunderstandingMap struct {
	ID                       primitive.ObjectID        `bson:"_id,omitempty" json:"id,omitempty"`
	WrongAnswer              string                    `bson:"wrong_answer" json:"wrong_answer"`
	CorrectAnswer            string                    `bson:"correct_answer" json:"correct_answer"`
	Domain                   string                    `bson:"domain" json:"domain"` // e.g., "immunology", "gcse"
	AgeGroup                 string                    `bson:"age_group" json:"age_group"`
	ErrorType                string                    `bson:"error_type" json:"error_type"` // phonetic, conceptual, computational, etc.
	SymptomPattern           string                    `bson:"symptom_pattern" json:"symptom_pattern"`
	RootCauses               []string                  `bson:"root_causes" json:"root_causes"`
	CorrectUnderstanding     string                    `bson:"correct_understanding" json:"correct_understanding"`
	RepairStrategy           string                    `bson:"repair_strategy" json:"repair_strategy"`
	RelatedMisunderstandings []RelatedMisunderstanding `bson:"related_misunderstandings" json:"related_misunderstandings"`
	EvidencePatterns         []string                  `bson:"evidence_patterns" json:"evidence_patterns"`
	QuestionKeywords         []string                  `bson:"question_keywords" json:"question_keywords"`
	Frequency                int64                     `bson:"frequency" json:"frequency"` // How many times this error occurred
	FirstObserved            time.Time                 `bson:"first_observed" json:"first_observed"`
	LastObserved             time.Time                 `bson:"last_observed" json:"last_observed"`
	ConfidenceScore          float64                   `bson:"confidence_score" json:"confidence_score"` // How certain we are this is a real pattern
}

// RelatedMisunderstanding represents a similar misconception
type RelatedMisunderstanding struct {
	WrongAnswer       string `bson:"wrong_answer" json:"wrong_answer"`
	CorrectAnswer     string `bson:"correct_answer" json:"correct_answer"`
	SymptomPattern    string `bson:"symptom_pattern" json:"symptom_pattern"`
	RepairStrategy    string `bson:"repair_strategy" json:"repair_strategy"`
	FrequencyObserved int64  `bson:"frequency_observed" json:"frequency_observed"`
}

// MisunderstandingStats aggregates stats about misconceptions
type MisunderstandingStats struct {
	Domain                 string
	TotalMisunderstandings int64
	MostCommon             *MisunderstandingMap
	ErrorTypeBreakdown     map[string]int64
	AgeGroupBreakdown      map[string]int64
}
