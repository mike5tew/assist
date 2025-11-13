package skills

import (
	"context"
	"esp-organizer/internal/models"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdvancedSkillService provides comprehensive skill management with relationships
type AdvancedSkillService struct {
	skillsCollection  *mongo.Collection
	sourcesCollection *mongo.Collection
	// ... other collections
}

// NewAdvancedSkillService creates a new AdvancedSkillService
func NewAdvancedSkillService(db *mongo.Database) *AdvancedSkillService {
	return &AdvancedSkillService{
		skillsCollection:  db.Collection("skills"),
		sourcesCollection: db.Collection("sources"),
	}
}

// GetSkillWithDetails fetches a skill and its associated parent skills and sources
func (s *AdvancedSkillService) GetSkillWithDetails(ctx context.Context, skillID primitive.ObjectID) (*models.Skill, error) {
	var skill models.Skill
	err := s.skillsCollection.FindOne(ctx, bson.M{"_id": skillID}).Decode(&skill)
	if err != nil {
		return nil, fmt.Errorf("failed to find skill: %w", err)
	}

	// Hydrate Parent Skills
	if len(skill.ParentSkillIDs) > 0 {
		cursor, err := s.skillsCollection.Find(ctx, bson.M{"_id": bson.M{"$in": skill.ParentSkillIDs}})
		if err != nil {
			return nil, fmt.Errorf("failed to find parent skills: %w", err)
		}
		defer cursor.Close(ctx)
		var parentSkills []*models.Skill
		if err = cursor.All(ctx, &parentSkills); err != nil {
			return nil, fmt.Errorf("failed to decode parent skills: %w", err)
		}
		skill.ParentSkills = parentSkills
	}

	// Hydrate Sources
	if len(skill.SourceRefs) > 0 {
		cursor, err := s.sourcesCollection.Find(ctx, bson.M{"_id": bson.M{"$in": skill.SourceRefs}})
		if err != nil {
			return nil, fmt.Errorf("failed to find sources: %w", err)
		}
		defer cursor.Close(ctx)
		var sources []*models.Source
		if err = cursor.All(ctx, &sources); err != nil {
			return nil, fmt.Errorf("failed to decode sources: %w", err)
		}
		// Assign sources to the skill
		skill.Sources = make([]*models.Source, len(sources))
		for i, src := range sources {
			skill.Sources[i] = &models.Source{
				ID:      src.ID,
				Title:   src.Title,
				Authors: src.Authors,
				Year:    src.Year,
				Type:    src.Type,
			}
		}
	} else {
		// If no sources, initialize to avoid nil slice
		skill.Sources = []*models.Source{}
	}

	return &skill, nil
}

// ToRAGChunk generates the text string for embedding
func ToRAGChunk(s *models.Skill) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Skill Name: %s. Description: %s. ", s.Name, s.Description))

	if len(s.SkillCriteria) > 0 {
		sb.WriteString("Criteria: ")
		for _, c := range s.SkillCriteria {
			sb.WriteString(fmt.Sprintf("%d. %s. ", c.Level, c.Description))
		}
	}

	if len(s.ParentSkills) > 0 {
		sb.WriteString("Parent Skills: ")
		for i, ps := range s.ParentSkills {
			sb.WriteString(ps.Name)
			if i < len(s.ParentSkills)-1 {
				sb.WriteString(", ")
			} else {
				sb.WriteString(". ")
			}
		}
	}

	if len(s.Sources) > 0 {
		sb.WriteString("Sources: ")
		for _, src := range s.Sources {
			// Use Authors field (which is now the primary field)
			author := strings.Join(src.Authors, ", ")
			if author == "" {
				author = "Unknown Author"
			}

			// Use Year field
			year := src.Year

			// Use Title field
			title := src.Title

			sb.WriteString(fmt.Sprintf("'%s' by %s (%s, %s). ", title, author, year, src.Type))
		}
	}

	return sb.String()
}

// CreateSkill inserts a single skill into MongoDB
func (s *AdvancedSkillService) CreateSkill(skill models.Skill) error {
	collection := s.skillsCollection
	ctx := context.TODO()

	_, err := collection.InsertOne(ctx, skill)
	return err
}

// CreateSkills inserts multiple skills into MongoDB
func (s *AdvancedSkillService) CreateSkills(skills []models.Skill) error {
	collection := s.skillsCollection
	ctx := context.TODO()

	for _, skill := range skills {
		_, err := collection.InsertOne(ctx, skill)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAllSkills retrieves all skills from MongoDB
func (s *AdvancedSkillService) GetAllSkills() ([]models.Skill, error) {
	collection := s.skillsCollection
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var skills []models.Skill
	if err = cursor.All(ctx, &skills); err != nil {
		return nil, err
	}

	return skills, nil
}

// GetSkillByID retrieves a skill by its ID
func (s *AdvancedSkillService) GetSkillByID(id primitive.ObjectID) (*models.Skill, error) {
	collection := s.skillsCollection
	ctx := context.TODO()

	var skill models.Skill
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&skill)
	if err != nil {
		return nil, err
	}

	return &skill, nil
}

// UpdateSkill updates an existing skill
func (s *AdvancedSkillService) UpdateSkill(id primitive.ObjectID, skill models.Skill) error {
	collection := s.skillsCollection
	ctx := context.TODO()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": skill},
	)
	return err
}

// DeleteSkill removes a skill from MongoDB
func (s *AdvancedSkillService) DeleteSkill(id primitive.ObjectID) error {
	collection := s.skillsCollection
	ctx := context.TODO()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteAllSkills removes all skills from MongoDB
func (s *AdvancedSkillService) DeleteAllSkills() error {
	collection := s.skillsCollection
	ctx := context.TODO()

	_, err := collection.DeleteMany(ctx, bson.M{})
	return err
}

// GetCollection returns the skills collection for advanced queries
func (s *AdvancedSkillService) GetCollection() *mongo.Collection {
	return s.skillsCollection
}
