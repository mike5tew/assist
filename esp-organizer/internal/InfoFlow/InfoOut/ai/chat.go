package ai

import (
	"esp-organizer/internal/config"
	"esp-organizer/internal/models"
)

// Chat represents a chat session with an AI models.

// formatSkillsForLLM formats skills data for LLMs
func FormatSkillsForLLM(skills []models.Skill) []map[string]interface{} {
	formatted := make([]map[string]interface{}, 0, len(skills))
	for _, skill := range skills {
		criteriaTexts := make([]string, len(skill.SkillCriteria))
		for i, criteria := range skill.SkillCriteria {
			criteriaTexts[i] = criteria.Description
		}
		skillData := map[string]interface{}{
			"skill_name":      skill.Name,
			"description":     skill.Description,
			"development_age": skill.DevelopmentAge,
			"criteria_levels": len(skill.SkillCriteria),
			"criteria":        criteriaTexts,
			"parent_count":    len(skill.ParentSkillIDs),
			"source_count":    len(skill.SourceRefs),
		}
		formatted = append(formatted, skillData)
	}
	return formatted
}

// generateSemanticInsights creates insights about the semantic relationships found
func GenerateSemanticInsights(response *models.QueryResponse) map[string]interface{} {
	insights := map[string]interface{}{
		"query_complexity":  "simple",
		"semantic_clusters": []string{},
		"confidence_levels": map[string]int{
			"high_confidence":   0,
			"medium_confidence": 0,
			"low_confidence":    0,
		},
	}
	for _, result := range response.WeaviateResults {
		if properties, ok := result["properties"].(map[string]interface{}); ok {
			certainty := config.GetCertaintyValue(properties)
			if certainty > 0.8 {
				insights["confidence_levels"].(map[string]int)["high_confidence"]++
			} else if certainty > 0.6 {
				insights["confidence_levels"].(map[string]int)["medium_confidence"]++
			} else {
				insights["confidence_levels"].(map[string]int)["low_confidence"]++
			}
		}
	}
	return insights
}

// formatWeaviateResultsForLLM formats Weaviate results for LLM consumption
func FormatWeaviateResultsForLLM(results []map[string]interface{}) []map[string]interface{} {
	formatted := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		if properties, ok := result["properties"].(map[string]interface{}); ok {
			skillData := map[string]interface{}{
				"skill_name":         config.GetStringValue(properties, "name"),
				"description":        config.GetStringValue(properties, "description"),
				"skill_type":         config.GetStringValue(properties, "skill_type"),
				"development_age":    config.GetIntValue(properties, "development_age"),
				"criteria":           config.GetStringValue(properties, "criteria_text"),
				"semantic_relevance": GetCertaintyValue(properties),
				"source_info": map[string]interface{}{
					"title":             config.GetStringValue(properties, "source_title"),
					"isbn":              config.GetStringValue(properties, "source_isbn"),
					"type":              config.GetStringValue(properties, "source_type"),
					"chapter_number":    config.GetStringValue(properties, "chapter_number"),
					"chapter_title":     config.GetStringValue(properties, "chapter_title"),
					"extraction_method": config.GetStringValue(properties, "extraction_method"),
					"extracted_at":      config.GetStringValue(properties, "extracted_at"),
				},
			}
			formatted = append(formatted, skillData)
		}
	}
	return formatted
}

// getIntValue safely extracts an int from a map
func GetIntValue(properties map[string]interface{}, key string) int {
	if val, ok := properties[key].(float64); ok {
		return int(val)
	}
	if val, ok := properties[key].(int); ok {
		return val
	}
	return 0
}

// getCertaintyValue extracts certainty from Weaviate _additional
func GetCertaintyValue(properties map[string]interface{}) float64 {
	if additional, ok := properties["_additional"].(map[string]interface{}); ok {
		if certainty, ok := additional["certainty"].(float64); ok {
			return certainty
		}
	}
	return 0.0
}
