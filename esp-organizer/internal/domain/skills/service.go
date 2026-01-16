package skills

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/domain/filters" // Add an import alias to avoid namespace collision
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	weaviatefilters "github.com/weaviate/weaviate-go-client/v4/weaviate/filters" // Add an alias for weaviate filters
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SkillService struct {
	mongoCollection *mongo.Collection
	weaviateClient  *weaviate.Client
	llmClient       *llm.LlamaClient
}

func NewSkillService(
	mongoCollection *mongo.Collection,
	weaviateClient *weaviate.Client,
	llmClient *llm.LlamaClient,
) *SkillService {
	return &SkillService{
		mongoCollection: mongoCollection,
		weaviateClient:  weaviateClient,
		llmClient:       llmClient,
	}
}

// CreateSemanticLink creates a single semantic relationship between skills
func (s *SkillService) CreateSemanticLink(ctx context.Context, link models.SemanticLink) error {
	if s.weaviateClient == nil {
		return fmt.Errorf("weaviate client not available")
	}
	if s.llmClient == nil {
		return fmt.Errorf("llm client not available for vectorization")
	}

	// Vectorize the relationship
	embeddingText := fmt.Sprintf("Relationship: %s %s %s. Context: %s",
		link.SourceTerm, link.RelationType, link.TargetTerm, link.Context)
	vector, err := s.llmClient.GenerateEmbedding(embeddingText)
	if err != nil {
		return fmt.Errorf("failed to generate embedding for semantic link: %w", err)
	}

	// Prepare properties for Weaviate
	properties := map[string]interface{}{
		"source_term":     link.SourceTerm,
		"target_term":     link.TargetTerm,
		"source_mongo_id": link.SourceID.Hex(),
		"target_mongo_id": link.TargetID.Hex(),
		"relation_type":   link.RelationType,
		"context":         link.Context,
		"confidence":      link.Confidence,
		"domain":          link.Domain,
	}

	// Store in Weaviate
	if err := db.StoreDocumentWithVector(ctx, "SemanticLinks", properties, vector); err != nil {
		return fmt.Errorf("failed to store semantic link in Weaviate: %w", err)
	}

	log.Printf("Successfully created semantic link: %s -> %s (%s)",
		link.SourceTerm, link.TargetTerm, link.RelationType)
	return nil
}

// StoreSemanticLinks vectorizes and stores a batch of semantic links.
// Uses the statement field as the primary vectorization target.
func (s *SkillService) StoreSemanticLinks(ctx context.Context, links []models.SemanticLink) error {
	if s.weaviateClient == nil {
		return fmt.Errorf("weaviate client not available")
	}
	if s.llmClient == nil {
		return fmt.Errorf("llm client not available for vectorization")
	}

	for _, link := range links {
		// Build embedding text from statement (the full citable fact)
		// If no statement provided, fall back to legacy format
		embeddingText := link.Statement
		if embeddingText == "" {
			embeddingText = fmt.Sprintf("Relationship: %s %s %s. Context: %s",
				link.SourceTerm, link.ForwardRelation, link.TargetTerm, link.Context)
		}
		vector, err := s.llmClient.GenerateEmbedding(embeddingText)
		if err != nil {
			log.Printf("Warning: failed to generate embedding for link '%s -> %s': %v",
				link.SourceTerm, link.TargetTerm, err)
			continue
		}

		// Serialize conditions to JSON for storage
		conditionsJSON := "[]"
		if len(link.Conditions) > 0 {
			if jsonBytes, err := json.Marshal(link.Conditions); err == nil {
				conditionsJSON = string(jsonBytes)
			}
		}

		// Prepare properties for Weaviate - include new statement-centric fields
		properties := map[string]interface{}{
			"statement":        link.Statement,
			"source_term":      link.SourceTerm,
			"target_term":      link.TargetTerm,
			"forward_relation": link.ForwardRelation,
			"inverse_relation": link.InverseRelation,
			"relation_type":    link.RelationType, // Legacy compatibility
			"context":          link.Context,
			"conditions_json":  conditionsJSON,
			"confidence":       link.Confidence,
			"domain":           link.Domain,
		}

		// Store in Weaviate
		if err := db.StoreDocumentWithVector(ctx, "SemanticLinks", properties, vector); err != nil {
			log.Printf("Warning: failed to store semantic link in Weaviate for '%s -> %s': %v",
				link.SourceTerm, link.TargetTerm, err)
		}
	}
	log.Printf("Successfully processed and stored %d semantic links in Weaviate.", len(links))
	return nil
}

// FindRelatedSemanticLinks finds semantic relationships for a given skill
func (s *SkillService) FindRelatedSemanticLinks(ctx context.Context, skillID primitive.ObjectID, relationType string) ([]models.SemanticLink, error) {
	if s.weaviateClient == nil {
		return nil, fmt.Errorf("weaviate client not available")
	}

	skillIDStr := skillID.Hex()

	// Build filter to find links where this skill is either source or target
	filter := weaviatefilters.Where().
		WithOperator(weaviatefilters.Or).
		WithOperands([]*weaviatefilters.WhereBuilder{
			weaviatefilters.Where().
				WithPath([]string{"source_mongo_id"}).
				WithOperator(weaviatefilters.Equal).
				WithValueString(skillIDStr),
			weaviatefilters.Where().
				WithPath([]string{"target_mongo_id"}).
				WithOperator(weaviatefilters.Equal).
				WithValueString(skillIDStr),
		})

	// Add relation type filter if specified
	if relationType != "" {
		relationFilter := weaviatefilters.Where().
			WithPath([]string{"relation_type"}).
			WithOperator(weaviatefilters.Equal).
			WithValueString(relationType)
		filter = weaviatefilters.Where().
			WithOperator(weaviatefilters.And).
			WithOperands([]*weaviatefilters.WhereBuilder{filter, relationFilter})
	}

	fields := []graphql.Field{
		{Name: "source_term"},
		{Name: "target_term"},
		{Name: "source_mongo_id"},
		{Name: "target_mongo_id"},
		{Name: "relation_type"},
		{Name: "context"},
		{Name: "confidence"},
		{Name: "domain"},
	}

	queryBuilder := s.weaviateClient.GraphQL().Get().
		WithClassName("SemanticLinks").
		WithFields(fields...).
		WithWhere(filter).
		WithLimit(50)

	response, err := queryBuilder.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("weaviate query failed: %w", err)
	}

	if response.Errors != nil {
		return nil, fmt.Errorf("weaviate query returned errors: %v", response.Errors)
	}

	data, ok := response.Data["Get"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from Weaviate")
	}

	linksInterface, ok := data["SemanticLinks"].([]interface{})
	if !ok {
		return []models.SemanticLink{}, nil
	}

	var semanticLinks []models.SemanticLink
	for _, linkInterface := range linksInterface {
		if linkMap, ok := linkInterface.(map[string]interface{}); ok {
			link, err := s.parseSemanticLinkFromWeaviate(linkMap)
			if err != nil {
				log.Printf("Warning: Failed to parse semantic link: %v", err)
				continue
			}
			semanticLinks = append(semanticLinks, *link)
		}
	}

	return semanticLinks, nil
}

// parseSemanticLinkFromWeaviate converts Weaviate response to SemanticLink model
func (s *SkillService) parseSemanticLinkFromWeaviate(linkMap map[string]interface{}) (*models.SemanticLink, error) {
	sourceID, err := primitive.ObjectIDFromHex(linkMap["source_mongo_id"].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid source ID: %w", err)
	}

	targetID, err := primitive.ObjectIDFromHex(linkMap["target_mongo_id"].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid target ID: %w", err)
	}

	return &models.SemanticLink{
		SourceTerm:   linkMap["source_term"].(string),
		TargetTerm:   linkMap["target_term"].(string),
		SourceID:     sourceID,
		TargetID:     targetID,
		RelationType: linkMap["relation_type"].(string),
		Context:      linkMap["context"].(string),
		Confidence:   linkMap["confidence"].(float64),
		Domain:       linkMap["domain"].(string),
	}, nil
}

// FindSkillsBySemanticSimilarity finds skills semantically similar to a query
func (s *SkillService) FindSkillsBySemanticSimilarity(ctx context.Context, queryText string, limit int) ([]models.Skill, error) {
	if s.weaviateClient == nil {
		return nil, fmt.Errorf("weaviate client not available")
	}

	vector := s.generateQueryVector(queryText)
	if vector == nil {
		return nil, fmt.Errorf("failed to generate query vector")
	}

	// Search in EducationalSkills class using semantic similarity
	nearVector := s.weaviateClient.GraphQL().
		NearVectorArgBuilder().
		WithVector(vector)

	fields := []graphql.Field{
		{Name: "name"},
		{Name: "description"},
		{Name: "skill_type"},
		{Name: "development_age"},
		{Name: "mongo_id"},
		{Name: "_additional", Fields: []graphql.Field{{Name: "certainty"}}},
	}

	queryBuilder := s.weaviateClient.GraphQL().Get().
		WithClassName("EducationalSkills").
		WithFields(fields...).
		WithNearVector(nearVector).
		WithLimit(limit)

	response, err := queryBuilder.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("weaviate query failed: %w", err)
	}

	if response.Errors != nil {
		return nil, fmt.Errorf("weaviate query returned errors: %v", response.Errors)
	}

	data, ok := response.Data["Get"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from Weaviate")
	}

	skillsInterface, ok := data["EducationalSkills"].([]interface{})
	if !ok {
		return []models.Skill{}, nil
	}

	var skills []models.Skill
	for _, skillInterface := range skillsInterface {
		if skillMap, ok := skillInterface.(map[string]interface{}); ok {
			skill, err := s.parseSkillFromWeaviate(skillMap)
			if err != nil {
				log.Printf("Warning: Failed to parse skill: %v", err)
				continue
			}
			skills = append(skills, *skill)
		}
	}

	return skills, nil
}

// parseSkillFromWeaviate converts Weaviate response to Skill model
func (s *SkillService) parseSkillFromWeaviate(skillMap map[string]interface{}) (*models.Skill, error) {
	mongoID, err := primitive.ObjectIDFromHex(skillMap["mongo_id"].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid mongo ID: %w", err)
	}

	// Get the actual skill from MongoDB for complete data
	var skill models.Skill
	err = s.mongoCollection.FindOne(context.Background(), bson.M{"_id": mongoID}).Decode(&skill)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch complete skill data from MongoDB: %w", err)
	}

	return &skill, nil
}

// GetSkillGraph builds a complete graph of skills and their semantic relationships
func (s *SkillService) GetSkillGraph(ctx context.Context, centerSkillID primitive.ObjectID, depth int) (map[string]interface{}, error) {
	if depth < 1 {
		depth = 1
	}
	if depth > 3 {
		depth = 3 // Limit depth for performance
	}

	// Get the center skill
	centerSkill, err := s.GetSkillByID(ctx, centerSkillID)
	if err != nil {
		return nil, fmt.Errorf("failed to get center skill: %w", err)
	}

	result := map[string]interface{}{
		"center_skill": centerSkill,
		"relationships": map[string]interface{}{
			"outgoing": []models.SemanticLink{},
			"incoming": []models.SemanticLink{},
		},
		"related_skills": []models.Skill{},
	}

	// Get direct relationships
	outgoingLinks, err := s.FindRelatedSemanticLinks(ctx, centerSkillID, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get semantic relationships: %w", err)
	}

	// Separate outgoing and incoming relationships
	var outgoing, incoming []models.SemanticLink
	relatedSkillIDs := make(map[primitive.ObjectID]bool)

	for _, link := range outgoingLinks {
		if link.SourceID == centerSkillID {
			outgoing = append(outgoing, link)
			relatedSkillIDs[link.TargetID] = true
		} else {
			incoming = append(incoming, link)
			relatedSkillIDs[link.SourceID] = true
		}
	}

	result["relationships"].(map[string]interface{})["outgoing"] = outgoing
	result["relationships"].(map[string]interface{})["incoming"] = incoming

	// Get related skills details
	var relatedSkills []models.Skill
	for skillID := range relatedSkillIDs {
		skill, err := s.GetSkillByID(ctx, skillID)
		if err != nil {
			log.Printf("Warning: Failed to get related skill %s: %v", skillID.Hex(), err)
			continue
		}
		relatedSkills = append(relatedSkills, *skill)
	}
	result["related_skills"] = relatedSkills

	return result, nil
}

// GetSkillByID gets a skill by its MongoDB ID
func (s *SkillService) GetSkillByID(ctx context.Context, id primitive.ObjectID) (*models.Skill, error) {
	var skill models.Skill
	err := s.mongoCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&skill)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}
	return &skill, nil
}

// GenerateQueryVector generates embedding vector for a query
func (s *SkillService) GenerateQueryVector(queryText string) ([]float32, error) {
	if s.llmClient == nil {
		return nil, fmt.Errorf("llm client not available")
	}
	return s.llmClient.GenerateEmbedding(queryText)
}

// CreateSkill creates a new skill in MongoDB and optionally creates semantic links for parent relationships
func (s *SkillService) CreateSkill(ctx context.Context, skill *models.Skill) (*models.Skill, error) {
	// First create the skill in MongoDB
	result, err := s.mongoCollection.InsertOne(ctx, skill)
	if err != nil {
		return nil, fmt.Errorf("failed to insert skill into MongoDB: %w", err)
	}

	skill.ID = result.InsertedID.(primitive.ObjectID)

	// Create semantic links for parent relationships
	for _, parentID := range skill.ParentSkillIDs {
		parentSkill, err := s.GetSkillByID(ctx, parentID)
		if err != nil {
			log.Printf("Warning: Failed to get parent skill %s: %v", parentID.Hex(), err)
			continue
		}

		semanticLink := models.SemanticLink{
			SourceID:     parentID,
			TargetID:     skill.ID,
			SourceTerm:   parentSkill.Name,
			TargetTerm:   skill.Name,
			RelationType: "parent_of",
			Context:      "Educational skill hierarchy",
			Confidence:   0.9,
			Domain:       "education",
		}

		if err := s.CreateSemanticLink(ctx, semanticLink); err != nil {
			log.Printf("Warning: Failed to create semantic link for parent relationship: %v", err)
		}
	}

	return skill, nil
}

// CreateSkills creates multiple skills in MongoDB
func (s *SkillService) CreateSkills(ctx context.Context, skills []models.Skill) ([]models.Skill, error) {
	var documents []interface{}
	for i := range skills {
		documents = append(documents, skills[i])
	}

	if len(documents) == 0 {
		return skills, nil
	}

	result, err := s.mongoCollection.InsertMany(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("failed to insert skills into MongoDB: %w", err)
	}

	// Update the skills with their generated IDs
	for i, insertedID := range result.InsertedIDs {
		if i < len(skills) {
			skills[i].ID = insertedID.(primitive.ObjectID)
		}
	}

	return skills, nil
}

// GetAllSkills retrieves all skills from the MongoDB collection
func (s *SkillService) GetAllSkills(ctx context.Context) ([]models.Skill, error) {
	cursor, err := s.mongoCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to query skills from MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	var skills []models.Skill
	if err := cursor.All(ctx, &skills); err != nil {
		return nil, fmt.Errorf("failed to decode skills from MongoDB: %w", err)
	}
	return skills, nil
}

// FindSkillsByName finds skills by name using regex search
func (s *SkillService) FindSkillsByName(ctx context.Context, name string) ([]models.Skill, error) {
	cursor, err := s.mongoCollection.Find(ctx, bson.M{"name": bson.M{"$regex": name, "$options": "i"}})
	if err != nil {
		return nil, fmt.Errorf("failed to find skills by name: %w", err)
	}
	defer cursor.Close(ctx)

	var skills []models.Skill
	if err = cursor.All(ctx, &skills); err != nil {
		return nil, fmt.Errorf("failed to decode skills: %w", err)
	}
	return skills, nil
}

// FindSkillsByIDs finds multiple skills by their MongoDB IDs
func (s *SkillService) FindSkillsByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Skill, error) {
	if len(ids) == 0 {
		return []models.Skill{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := s.mongoCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find skills by IDs: %w", err)
	}
	defer cursor.Close(ctx)

	var skills []models.Skill
	if err = cursor.All(ctx, &skills); err != nil {
		return nil, fmt.Errorf("failed to decode skills: %w", err)
	}

	return skills, nil
}

// Helper method to generate query vector
func (s *SkillService) generateQueryVector(queryText string) []float32 {
	if s.llmClient == nil {
		return nil
	}
	vector, err := s.llmClient.GenerateEmbedding(queryText)
	if err != nil {
		log.Printf("Warning: Failed to generate query vector: %v", err)
		return nil
	}
	return vector
}

// SearchSkillsByVectorInClass performs a vector-based semantic search in Weaviate
func (s *SkillService) SearchSkillsByVectorInClass(ctx context.Context, queryText string, filter *models.Filter, className string) ([]map[string]interface{}, error) {
	if s.weaviateClient == nil {
		return nil, fmt.Errorf("weaviate client not available")
	}

	// Generate vector using AWS through your LLM client
	vector, err := s.llmClient.GenerateEmbedding(queryText)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query vector using AWS: %w", err)
	}

	if len(vector) == 0 {
		return nil, fmt.Errorf("empty vector generated for query")
	}

	log.Printf("Performing vector search on class: %s with vector length: %d", className, len(vector))

	// Define fields based on the class
	var fields []graphql.Field
	switch className {
	case "SemanticLinks":
		fields = []graphql.Field{
			{Name: "source_term"},
			{Name: "target_term"},
			{Name: "relation_type"},
			{Name: "context"},
			{Name: "confidence"},
			{Name: "domain"},
			{Name: "source_mongo_id"},
			{Name: "target_mongo_id"},
			{Name: "_additional", Fields: []graphql.Field{{Name: "certainty"}}},
		}
	case "EducationalSkills":
		fields = []graphql.Field{
			{Name: "name"},
			{Name: "description"},
			{Name: "skill_type"},
			{Name: "development_age"},
			{Name: "mongo_id"},
			{Name: "source_refs"},
			{Name: "_additional", Fields: []graphql.Field{{Name: "certainty"}}},
		}
	default:
		return nil, fmt.Errorf("unsupported class for vector search: %s", className)
	}

	nearVector := s.weaviateClient.GraphQL().
		NearVectorArgBuilder().
		WithVector(vector)

	queryBuilder := s.weaviateClient.GraphQL().Get().
		WithClassName(className).
		WithFields(fields...).
		WithNearVector(nearVector).
		WithLimit(10)

	if filter != nil {
		where := filters.BuildWhereFilterFromAPI(filter)
		queryBuilder = queryBuilder.WithWhere(where)
		log.Printf("Applying filter for class %s", className)
	}

	response, err := queryBuilder.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("weaviate query failed: %w", err)
	}

	if response.Errors != nil {
		return nil, fmt.Errorf("weaviate query errors: %v", response.Errors)
	}

	data, ok := response.Data["Get"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from Weaviate")
	}

	items, ok := data[className].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil // No results is not an error
	}

	results := make([]map[string]interface{}, len(items))
	for i, item := range items {
		results[i] = map[string]interface{}{
			"properties": item,
		}
	}

	log.Printf("Found %d results in class %s", len(results), className)
	return results, nil
}

// SearchSkillsByVector performs a semantic search in Weaviate against the SemanticLinks class.
// This is kept for backward compatibility, but the SearchSkillsByVectorInClass method should be used.
func (s *SkillService) SearchSkillsByVector(ctx context.Context, queryText string, filter *models.Filter) ([]map[string]interface{}, error) {
	// By default, we search the SemanticLinks class
	return s.SearchSkillsByVectorInClass(ctx, queryText, filter, "SemanticLinks")
}
