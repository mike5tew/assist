package api

import (
	"esp-organizer/internal/models"
	"fmt"
	"math"
	"strings"
	"time"
)

// isMedicalQuery detects if a query is related to medical content
func isMedicalQuery(query string) bool {
	queryLower := strings.ToLower(query)

	medicalKeywords := []string{
		"immunology", "case study", "medical", "clinical", "patient",
		"processed", "textract", "real", "summary", "summarize",
		"immunoglobulin", "agammaglobulinemia", "antibody", "antigen",
		"lymphocyte", "cytokine", "complement", "macrophage", "neutrophil",
		"chapter", "extracted", "deficiency", "syndrome", "disease",
		"disorder", "infection", "autoimmune", "hypersensitivity", "allergy",
	}

	for _, keyword := range medicalKeywords {
		if strings.Contains(queryLower, keyword) {
			return true
		}
	}
	return false
}

// GenerateEnhancedAIResponse creates AI responses for skills queries
func GenerateEnhancedAIResponse(question string, data *models.QueryResponse) map[string]interface{} {
	if data == nil {
		return map[string]interface{}{
			"type":       "error",
			"content":    "No data received from query processor",
			"confidence": "low",
		}
	}

	totalResults := len(data.MongoResults) + len(data.WeaviateResults)

	if totalResults == 0 {
		return map[string]interface{}{
			"type": "no_results",
			"content": fmt.Sprintf(
				"I couldn't find specific information about '%s' in the current knowledge base. "+
					"Try using different keywords or check if the topic is covered.",
				question),
			"confidence": "low",
		}
	}

	return map[string]interface{}{
		"type":         "comprehensive_answer",
		"content":      fmt.Sprintf("Found %d relevant results for your query about '%s'.", totalResults, question),
		"confidence":   "high",
		"result_count": totalResults,
		"details": map[string]interface{}{
			"mongo_results":    len(data.MongoResults),
			"weaviate_results": len(data.WeaviateResults),
		},
	}
}

// createMedicalResponseFromEmbeddings creates response based on embedding search results
func createMedicalResponseFromEmbeddings(query string, results []map[string]interface{}) map[string]interface{} {
	if len(results) == 0 {
		return createFallbackMedicalResponse(query)
	}

	// Since we're using embeddings (not text generation), return structured results
	return map[string]interface{}{
		"type":        "embedding_based_results",
		"query":       query,
		"content":     fmt.Sprintf("Found %d semantically relevant medical content matches.", len(results)),
		"confidence":  calculateEmbeddingConfidence(results),
		"total_count": len(results),
		"results":     formatMedicalResults(results),
		"search_type": "semantic_embedding",
		"timestamp":   time.Now(),
	}
}

// formatMedicalResults formats the search results for presentation
func formatMedicalResults(results []map[string]interface{}) []map[string]interface{} {
	var formatted []map[string]interface{}

	for i, result := range results {
		if i >= 5 { // Limit results
			break
		}

		formattedResult := map[string]interface{}{
			"rank": i + 1,
		}

		// Extract relevant fields
		if title, ok := result["title"].(string); ok {
			formattedResult["title"] = title
		}

		if content, ok := result["content"].(string); ok {
			// Shorten content for display
			if len(content) > 200 {
				formattedResult["content_preview"] = content[:200] + "..."
			} else {
				formattedResult["content_preview"] = content
			}
		}

		if contentType, ok := result["contentType"].(string); ok {
			formattedResult["type"] = contentType
		}

		if score, ok := result["score"].(float64); ok {
			formattedResult["similarity_score"] = math.Round(score*1000) / 1000
		}

		formatted = append(formatted, formattedResult)
	}

	return formatted
}

// calculateEmbeddingConfidence determines confidence based on search results
func calculateEmbeddingConfidence(results []map[string]interface{}) string {
	if len(results) == 0 {
		return "low"
	}

	// Check if we have high similarity scores
	for _, result := range results {
		if score, ok := result["score"].(float64); ok && score > 0.8 {
			return "high"
		}
	}

	if len(results) >= 3 {
		return "medium"
	}

	return "low"
}

// createFallbackMedicalResponse when no results found
func createFallbackMedicalResponse(query string) map[string]interface{} {
	medicalKnowledge := map[string]string{
		"agammaglobulinemia": "Primary immunodeficiency with absent B cells. Requires immunoglobulin therapy.",
		"immunodeficiency":   "Impaired immune function. Can be primary (genetic) or secondary (acquired).",
		"antibody":           "Proteins produced by B cells to neutralize pathogens.",
		"immunology":         "Study of the immune system and its functions.",
	}

	queryLower := strings.ToLower(query)
	for term, info := range medicalKnowledge {
		if strings.Contains(queryLower, term) {
			return map[string]interface{}{
				"type":       "general_medical_knowledge",
				"content":    info,
				"confidence": "medium",
				"note":       "Based on general knowledge. Upload content for specific information.",
			}
		}
	}

	return map[string]interface{}{
		"type":       "no_specific_content",
		"content":    fmt.Sprintf("No specific medical content found for '%s'. Upload immunology materials for better results.", query),
		"confidence": "low",
		"suggestions": []string{
			"Upload immunology textbook chapters",
			"Try different medical terminology",
			"Check spelling or use broader terms",
		},
	}
}
