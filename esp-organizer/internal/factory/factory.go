package factory

import (
	"esp-organizer/internal/InfoFlow/InfoIn/rawText/query"
	"esp-organizer/internal/InfoFlow/InfoOut/llm"
	"esp-organizer/internal/InfoFlow/InfoStore/skills"
	"esp-organizer/internal/models"
)

// NewQueryProcessor creates a new query processor
func NewQueryProcessor(skillService *skills.SkillService, llamaClient *llm.LlamaClient) models.QueryProcessor {
	return query.NewQueryHandler(skillService, llamaClient)
}
