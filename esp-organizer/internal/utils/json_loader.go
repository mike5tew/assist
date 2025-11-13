package utils

import (
	"encoding/json"
	"esp-organizer/internal/models"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// jsonSkill is a temporary struct to robustly unmarshal skills from JSON.
// It uses flexible types to avoid immediate unmarshalling errors.
type jsonSkill struct {
	ID             json.RawMessage   `json:"_id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Category       string            `json:"category"`
	DevelopmentAge int               `json:"developmentAge"`
	ParentSkillIDs []json.RawMessage `json:"parent_skill_ids"` // Corrected tag
	SourceRefs     json.RawMessage   `json:"sourceRefs"`
	LastUpdated    json.RawMessage   `json:"lastUpdated"`
	Criteria       json.RawMessage   `json:"criteria"`
}

// LoadSkillsFromJSON loads skills from a JSON file, handling MongoDB's Extended JSON format
func LoadSkillsFromJSON(filePath string) ([]models.Skill, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills JSON file: %w", err)
	}

	// Unmarshal into the flexible jsonSkill struct to prevent immediate failure
	var tempSkills []jsonSkill
	if err := json.Unmarshal(fileData, &tempSkills); err != nil {
		return nil, fmt.Errorf("failed to perform initial unmarshal of skills JSON: %w", err)
	}

	var skills []models.Skill
	for _, tempSkill := range tempSkills {
		skill := models.Skill{
			Name:           tempSkill.Name,
			Description:    tempSkill.Description,
			Category:       tempSkill.Category,
			DevelopmentAge: tempSkill.DevelopmentAge,
		}

		// Manually parse the _id field
		var idStr string
		var idHolder struct {
			Oid string `json:"$oid"`
		}
		// Try parsing as extended JSON object
		if err := json.Unmarshal(tempSkill.ID, &idHolder); err == nil && idHolder.Oid != "" {
			idStr = idHolder.Oid
		} else {
			// Fallback to parsing as a simple string
			_ = json.Unmarshal(tempSkill.ID, &idStr)
		}
		skill.IDString = idStr

		// Manually parse the parentSkillIDs
		var parentIDs []primitive.ObjectID
		for _, rawParentID := range tempSkill.ParentSkillIDs {
			var parentIDStr string
			var parentIDHolder struct {
				Oid string `json:"$oid"`
			}
			// Try parsing as extended JSON object
			if err := json.Unmarshal(rawParentID, &parentIDHolder); err == nil && parentIDHolder.Oid != "" {
				parentIDStr = parentIDHolder.Oid
			} else {
				// Fallback to parsing as a simple string
				_ = json.Unmarshal(rawParentID, &parentIDStr)
			}

			if objID, err := primitive.ObjectIDFromHex(parentIDStr); err == nil {
				parentIDs = append(parentIDs, objID)
			} else {
				log.Printf("Warning: Invalid parent ObjectID format '%s' for skill '%s'", parentIDStr, skill.Name)
			}
		}
		skill.ParentSkillIDs = parentIDs

		// Unmarshal remaining fields that are less likely to cause issues
		_ = json.Unmarshal(tempSkill.Criteria, &skill.SkillCriteria)
		_ = json.Unmarshal(tempSkill.SourceRefs, &skill.SourceRefs)
		_ = json.Unmarshal(tempSkill.LastUpdated, &skill.LastUpdated)

		skills = append(skills, skill)
	}

	if len(skills) == 0 && len(tempSkills) > 0 {
		return nil, fmt.Errorf("all skills failed to process, please check JSON format and log warnings")
	}

	return skills, nil
}
