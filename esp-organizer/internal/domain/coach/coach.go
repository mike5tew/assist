package coach

import (
	"context"
	"esp-organizer/internal/domain/infoin"
	"esp-organizer/internal/models"
	"fmt"
)

// CoachServiceMVP defines the interface for the MVP coach's business logic.
type CoachServiceMVP interface {
	GetKnowledgeRoutes(ctx context.Context, rootID string) ([]models.SemanticLink, error)
}

// coachServiceMVP is the implementation of the CoachServiceMVP interface.
type coachServiceMVP struct {
	hsgService *infoin.HSGQueryService
}

// NewCoachServiceMVP creates a new instance of the coach service.
func NewCoachServiceMVP() (CoachServiceMVP, error) {
	hsgService, err := infoin.NewHSGQueryService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize HSGQueryService for coach: %w", err)
	}
	return &coachServiceMVP{
		hsgService: hsgService,
	}, nil
}

// GetKnowledgeRoutes encapsulates the logic to traverse the knowledge graph.
func (s *coachServiceMVP) GetKnowledgeRoutes(ctx context.Context, rootID string) ([]models.SemanticLink, error) {
	return s.hsgService.TraverseHierarchy(ctx, rootID)
}
