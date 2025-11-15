package query

import (
	"context"
	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/domain/skills"
	"esp-organizer/internal/models"

	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// QueryHandler implements the models.QueryResponse interface
type QueryHandler struct {
	skillService *skills.SkillService
	llamaClient  *llm.LlamaClient
}

// NewQueryHandler creates a new query handler
func NewQueryHandler(skillService *skills.SkillService, llamaCl *llm.LlamaClient) models.QueryProcessor {
	return &QueryHandler{
		skillService: skillService,
		llamaClient:  llamaCl,
	}
}

// searchMongoDB performs text search on MongoDB skills collection
func (qh *QueryHandler) searchMongoDB(ctx context.Context, query string) ([]models.Skill, error) {
	// Use the service method instead of direct collection access
	skills, err := qh.skillService.FindSkillsByName(ctx, query)
	if err != nil {
		log.Printf("MongoDB search error: %v", err)
		return nil, err
	}

	log.Printf("Found %d skills in MongoDB", len(skills))
	return skills, nil
}

func (qh *QueryHandler) searchWeaviate(ctx context.Context, query string, filters *models.Filter, weaviateClass string) ([]map[string]interface{}, error) {
	log.Printf("Performing Weaviate search for: %s in class: %s", query, weaviateClass)

	var results []map[string]interface{}
	var err error

	// If a specific class is requested, use it
	if weaviateClass != "" {
		results, err = qh.skillService.SearchSkillsByVectorInClass(ctx, query, filters, weaviateClass)
	} else {
		// Default to searching both classes
		results, err = qh.searchBothWeaviateClasses(ctx, query, filters)
	}

	if err != nil {
		return nil, err
	}

	// Post-process results to handle any data formatting
	for _, result := range results {
		if properties, ok := result["properties"].(map[string]interface{}); ok {
			qh.processWeaviateResults(properties)
		}
	}

	log.Printf("Found %d results in Weaviate", len(results))
	return results, nil
}

// searchBothWeaviateClasses searches both EducationalSkills and SemanticLinks
func (qh *QueryHandler) searchBothWeaviateClasses(ctx context.Context, query string, filters *models.Filter) ([]map[string]interface{}, error) {
	var allResults []map[string]interface{}

	// Search EducationalSkills
	skillsResults, err := qh.skillService.SearchSkillsByVectorInClass(ctx, query, filters, "EducationalSkills")
	if err != nil {
		log.Printf("Error searching EducationalSkills: %v", err)
	} else {
		allResults = append(allResults, skillsResults...)
	}

	// Search SemanticLinks
	linksResults, err := qh.skillService.SearchSkillsByVectorInClass(ctx, query, filters, "SemanticLinks")
	if err != nil {
		log.Printf("Error searching SemanticLinks: %v", err)
	} else {
		allResults = append(allResults, linksResults...)
	}

	return allResults, nil
}

// processWeaviateResults processes and formats Weaviate results
func (qh *QueryHandler) processWeaviateResults(properties map[string]interface{}) {
	// Handle different field names based on class type or data structure
	// This is a generic processor that can handle various Weaviate result formats

	// Convert string arrays if they exist
	arrayFields := []string{"parent_mongo_ids", "source_mongo_ids", "source_refs"}

	for _, field := range arrayFields {
		if value, ok := properties[field].(string); ok && value != "" {
			properties[field] = strings.Split(value, ",")
		}
	}

	// Handle criteria_text if it exists (from your Weaviate schema)
	if criteriaText, ok := properties["criteria_text"].(string); ok {
		// You might want to parse this into structured criteria if needed
		properties["criteria"] = qh.parseCriteriaText(criteriaText)
	}
}

// parseCriteriaText parses the criteria_text field into structured criteria
func (qh *QueryHandler) parseCriteriaText(criteriaText string) []map[string]interface{} {
	// This is a simple parser - adjust based on your actual criteria_text format
	var criteria []map[string]interface{}

	// Example: "Level 1: Description 1; Level 2: Description 2"
	levels := strings.Split(criteriaText, ";")
	for _, level := range levels {
		parts := strings.SplitN(level, ":", 2)
		if len(parts) == 2 {
			criteria = append(criteria, map[string]interface{}{
				"level":       strings.TrimSpace(parts[0]),
				"description": strings.TrimSpace(parts[1]),
			})
		}
	}

	return criteria
}

// findRelatedSkills finds skills related to the main results
func (qh *QueryHandler) findRelatedSkills(ctx context.Context, mainSkills []models.Skill) ([]models.Skill, error) {
	if len(mainSkills) == 0 {
		return nil, nil
	}

	var parentIDs []primitive.ObjectID
	for _, skill := range mainSkills {
		parentIDs = append(parentIDs, skill.ParentSkillIDs...)
	}

	if len(parentIDs) == 0 {
		return nil, nil
	}

	// Use the service method to find related skills
	relatedSkills, err := qh.skillService.FindSkillsByIDs(ctx, parentIDs)
	if err != nil {
		log.Printf("Error finding related skills: %v", err)
		return nil, err
	}

	return relatedSkills, nil
}

// ProcessQuery handles a natural language query about skills
func (qh *QueryHandler) ProcessQuery(ctx context.Context, query string, filters *models.Filter, weaviateClass string) (*models.QueryResponse, error) {
	log.Printf("Processing query: %s in class: %s", query, weaviateClass)

	response := &models.QueryResponse{
		Query:     query,
		Timestamp: time.Now(),
	}

	// 1. Search MongoDB for matching skills
	mongoResults, err := qh.searchMongoDB(ctx, query)
	if err != nil {
		log.Printf("MongoDB search error: %v", err)
	}
	response.MongoResults = mongoResults

	// 2. Search Weaviate for semantic matches
	weaviateResults, err := qh.searchWeaviate(ctx, query, filters, weaviateClass)
	if err != nil {
		log.Printf("Weaviate search error: %v", err)
	}
	response.WeaviateResults = weaviateResults

	// 3. Find related skills based on MongoDB results
	if len(mongoResults) > 0 {
		relatedSkills, err := qh.findRelatedSkills(ctx, mongoResults)
		if err != nil {
			log.Printf("Error finding related skills: %v", err)
		} else {
			response.RelatedSkills = relatedSkills
		}
	}

	// 4. Generate synthesis if LLM client is available
	if qh.llamaClient != nil {
		synthesis, err := qh.generateSynthesis(ctx, query, response)
		if err != nil {
			log.Printf("Error generating synthesis: %v", err)
		} else {
			response.Synthesis = synthesis
		}
	}

	return response, nil
}

// generateSynthesis creates a summary of the query results using LLM
func (qh *QueryHandler) generateSynthesis(ctx context.Context, query string, response *models.QueryResponse) (string, error) {
	if len(response.MongoResults) == 0 && len(response.WeaviateResults) == 0 {
		return "No skills found matching your query.", nil
	}

	// Build a prompt for the LLM
	prompt := qh.buildSynthesisPrompt(query, response)

	// Use the LLM client to generate a response
	synthesis, err := qh.llamaClient.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate synthesis: %w", err)
	}

	return synthesis, nil
}

// buildSynthesisPrompt creates a prompt for the LLM based on query results
func (qh *QueryHandler) buildSynthesisPrompt(query string, response *models.QueryResponse) string {
	prompt := fmt.Sprintf("Based on the following educational skills search results, provide a comprehensive answer to the query: '%s'\n\n", query)

	// Add MongoDB results (exact matches)
	if len(response.MongoResults) > 0 {
		prompt += "Exact Matches Found:\n"
		for i, skill := range response.MongoResults {
			if i >= 3 { // Limit to top 3 for prompt length
				break
			}
			prompt += fmt.Sprintf("- %s: %s\n", skill.Name, skill.Description)
		}
		prompt += "\n"
	}

	// Add Weaviate results (semantic matches)
	if len(response.WeaviateResults) > 0 {
		prompt += "Semantically Related Skills:\n"
		for i, result := range response.WeaviateResults {
			if i >= 3 { // Limit to top 3 for prompt length
				break
			}
			if properties, ok := result["properties"].(map[string]interface{}); ok {
				name := properties["name"].(string)
				desc := properties["description"].(string)
				prompt += fmt.Sprintf("- %s: %s\n", name, desc)
			}
		}
		prompt += "\n"
	}

	// Add related skills
	if len(response.RelatedSkills) > 0 {
		prompt += "Related Skills:\n"
		for i, skill := range response.RelatedSkills {
			if i >= 3 { // Limit to top 3 for prompt length
				break
			}
			prompt += fmt.Sprintf("- %s: %s\n", skill.Name, skill.Description)
		}
	}

	prompt += "\nPlease provide a comprehensive summary that connects these skills and answers the user's query."

	return prompt
}
