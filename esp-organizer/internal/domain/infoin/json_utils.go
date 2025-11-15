package infoin

import (
	"strings"
)

// CleanJSONResponse removes markdown code fences and extra whitespace from LLM responses
func CleanJSONResponse(response string) string {
	// Remove markdown code fences (```json, ```, etc.)
	response = strings.TrimSpace(response)

	// Remove opening code fence
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
	}

	// Remove closing code fence
	if strings.HasSuffix(response, "```") {
		response = strings.TrimSuffix(response, "```")
	}

	// Trim remaining whitespace
	response = strings.TrimSpace(response)

	return response
}

// TrimJSONArray removes array brackets and trims whitespace
func TrimJSONArray(json string) string {
	json = strings.TrimSpace(json)
	json = strings.TrimPrefix(json, "[")
	json = strings.TrimSuffix(json, "]")
	return strings.TrimSpace(json)
}
