package utils

import (
	"esp-organizer/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// convertRawSkillToSkill converts a raw JSON skill object to a models.Skill
func convertRawSkillToSkill(raw map[string]interface{}) (models.Skill, error) {
	skill := models.Skill{}

	// Handle ID field
	if idField, ok := raw["_id"]; ok {
		if idMap, ok := idField.(map[string]interface{}); ok {
			if oidStr, ok := idMap["$oid"].(string); ok {
				if oid, err := primitive.ObjectIDFromHex(oidStr); err == nil {
					skill.ID = oid
				}
			}
		}
	}

	// Handle basic string fields
	if name, ok := raw["name"].(string); ok {
		skill.Name = name
	}

	if description, ok := raw["description"].(string); ok {
		skill.Description = description
	}

	// Handle numeric fields
	if devAge, ok := raw["developmentAge"].(float64); ok {
		skill.DevelopmentAge = int(devAge)
	}

	// Handle criteria array
	if criteriaField, ok := raw["criteria"]; ok {
		if criteriaArray, ok := criteriaField.([]interface{}); ok {
			for _, criteriaItem := range criteriaArray {
				if criteriaMap, ok := criteriaItem.(map[string]interface{}); ok {
					criteria := models.SkillCriteria{}
					if level, ok := criteriaMap["level"].(float64); ok {
						criteria.Level = int(level)
					}
					if desc, ok := criteriaMap["description"].(string); ok {
						criteria.Description = desc
					}
					skill.SkillCriteria = append(skill.SkillCriteria, criteria)
				}
			}
		}
	}

	// Handle parent skills (multiple possible field names)
	parentFields := []string{"parentSkills", "parentSkillIDs", "parent_mongo_ids"}
	for _, fieldName := range parentFields {
		if parentField, ok := raw[fieldName]; ok {
			skill.ParentSkillIDs = convertToObjectIDArray(parentField)
			break
		}
	}

	// Handle source references
	sourceFields := []string{"sourceRefs", "sourceMongoIDs", "source_mongo_ids"}
	for _, fieldName := range sourceFields {
		if sourceField, ok := raw[fieldName]; ok {
			skill.SourceRefs = convertToSourceReferenceArray(sourceField)
			break
		}
	}

	return skill, nil
}

// convertToSourceReferenceArray converts various array formats into []models.SourceReference
func convertToSourceReferenceArray(field interface{}) []models.SourceReference {
	var result []models.SourceReference

	switch v := field.(type) {
	case []interface{}:
		for _, item := range v {
			ref := models.SourceReference{}
			switch itemType := item.(type) {
			case string:
				// Assumes the string is a source ID
				ref.SourceID = itemType
			case map[string]interface{}:
				// Handles structured source references
				if id, ok := itemType["source_id"].(string); ok {
					ref.SourceID = id
				}
				if title, ok := itemType["title"].(string); ok {
					ref.Title = title
				}
				// Also handle BSON ObjectID format
				if idMap, ok := itemType["$oid"].(string); ok {
					ref.SourceID = idMap
				}
			}
			if ref.SourceID != "" {
				result = append(result, ref)
			}
		}
	case []string:
		for _, str := range v {
			result = append(result, models.SourceReference{SourceID: str})
		}
	}

	return result
}

// convertToObjectIDArray converts various array formats to []primitive.ObjectID
func convertToObjectIDArray(field interface{}) []primitive.ObjectID {
	var result []primitive.ObjectID

	switch v := field.(type) {
	case []interface{}:
		for _, item := range v {
			switch itemType := item.(type) {
			case string:
				if oid, err := primitive.ObjectIDFromHex(itemType); err == nil {
					result = append(result, oid)
				}
			case map[string]interface{}:
				if oidStr, ok := itemType["$oid"].(string); ok {
					if oid, err := primitive.ObjectIDFromHex(oidStr); err == nil {
						result = append(result, oid)
					}
				}
			}
		}
	case []string:
		for _, str := range v {
			if oid, err := primitive.ObjectIDFromHex(str); err == nil {
				result = append(result, oid)
			}
		}
	}

	return result
}
