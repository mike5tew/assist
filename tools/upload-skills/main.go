package main

import (
	"context"
	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/config"
	"esp-organizer/internal/domain/skills"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"esp-organizer/internal/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Centralized config loading at the start of the application.
	projectRoot, err := utils.FindProjectRoot()
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}
	config.LoadConfig(projectRoot, ".env")

	ctx := context.Background()

	// The Makefile now handles sourcing the .env file, so this Go program
	// will inherit the correct environment variables.

	mngoDB, err := db.NewMongoDBFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mngoDB.Client.Disconnect(ctx); err != nil {
			log.Printf("Warning: Error disconnecting from MongoDB: %v", err)
		}
	}()

	skillsCollection := mngoDB.Database.Collection("skills")
	log.Printf("Using collection: %s", "skills")

	// Ensure the debug book source exists for API/UI workflows
	ensureDebugBookSource(ctx, mngoDB)

	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}
	weaviateClient := db.GetWeaviateClient()
	log.Println("✅ Successfully connected to Weaviate.")

	// The "EducationalSkills" class is deprecated. We no longer upload raw skills to Weaviate.
	// Vector search will be performed on the "SemanticLinks" class.
	log.Println("Ensuring core Weaviate schema is present...")
	if err := db.EnsureCoreSchema(ctx); err != nil {
		log.Fatalf("Failed to ensure Weaviate schema: %v", err)
	}

	log.Println("Clearing existing educational skills data before upload...")
	if _, err := skillsCollection.DeleteMany(ctx, bson.M{}); err != nil {
		log.Fatalf("Failed to clear skills collection: %v", err)
	}
	log.Println("✅ Cleared MongoDB skills collection")

	// Initialize the LLM client for vectorization
	llmClient := llm.NewLlamaClient()
	log.Println("✅ Initialized LLM Client for vectorization.")

	skillsJSONPath := "./internal/temp/skills.json"
	log.Println("Loading skills from JSON...")
	skillsData, err := utils.LoadSkillsFromJSON(skillsJSONPath)
	if err != nil {
		log.Fatalf("Failed to load skills: %v", err)
	}
	log.Printf("Loaded %d skills from JSON", len(skillsData))

	// Pre-process skill IDs to convert strings to ObjectIDs before insertion
	preprocessSkillIDs(skillsData)

	skillService := skills.NewSkillService(
		mngoDB.Database.Collection("skills"),
		weaviateClient,
		llmClient,
	)
	log.Println("Uploading skills to MongoDB...")
	if _, err := skillService.CreateSkills(ctx, skillsData); err != nil {
		log.Fatalf("Failed to upload skills to MongoDB: %v", err)
	}
	log.Println("✅ Successfully uploaded skills to MongoDB")

	log.Println("Extracting and vectorizing semantic links from skill hierarchy...")
	if err := createSkillHierarchyLinks(ctx, skillService, skillsData); err != nil {
		log.Fatalf("Failed to create skill hierarchy links in Weaviate: %v", err)
	}

	log.Println("🎉 Skills upload and semantic link vectorization completed successfully!")
}

// createSkillHierarchyLinks iterates through skills, identifies parent-child relationships,
// and stores them as vectorized semantic links in Weaviate.
func createSkillHierarchyLinks(ctx context.Context, skillService *skills.SkillService, skillsData []models.Skill) error {
	var links []models.SemanticLink
	skillMap := make(map[primitive.ObjectID]models.Skill)
	skillNameMap := make(map[string]primitive.ObjectID) // Additional map for name lookups

	// First, get all skills from MongoDB to ensure we have their generated IDs
	log.Println("Fetching all skills from MongoDB...")
	allSkills, err := skillService.GetAllSkills(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve all skills from MongoDB: %w", err)
	}
	log.Printf("Retrieved %d skills from MongoDB", len(allSkills))

	// First, populate the map with all skills for reliable lookups
	for _, skill := range allSkills {
		skillMap[skill.ID] = skill
		skillNameMap[strings.ToLower(skill.Name)] = skill.ID
	}

	// Now, log all skill names and their resolved parent names
	log.Println("Available skills in MongoDB:")
	for _, skill := range allSkills {
		log.Printf("  - Skill: %s (ID: %s)", skill.Name, skill.ID.Hex())

		// Debug parent skill IDs
		if len(skill.ParentSkillIDs) > 0 {
			parentNames := make([]string, 0, len(skill.ParentSkillIDs))
			for _, pid := range skill.ParentSkillIDs {
				if parent, ok := skillMap[pid]; ok {
					parentNames = append(parentNames, parent.Name)
				} else {
					// This should no longer happen with the pre-populated map
					parentNames = append(parentNames, fmt.Sprintf("Unknown(%s)", pid.Hex()))
				}
			}
			log.Printf("    Parent Skills: %s", strings.Join(parentNames, ", "))
		} else {
			log.Printf("    No parent skills defined")
		}
	}

	// Now, create links based on parent-child relationships
	log.Println("Creating semantic links from parent-child relationships...")
	relationshipCount := 0

	for _, skill := range allSkills {
		if len(skill.ParentSkillIDs) > 0 {
			log.Printf("Processing relationships for skill '%s' with %d parent(s)",
				skill.Name, len(skill.ParentSkillIDs))

			for _, parentID := range skill.ParentSkillIDs {
				if parentSkill, ok := skillMap[parentID]; ok {
					log.Printf("  ✓ Found parent skill: %s", parentSkill.Name)

					link := models.SemanticLink{
						SourceTerm:   skill.Name,
						TargetTerm:   parentSkill.Name,
						SourceID:     skill.ID,
						TargetID:     parentSkill.ID,
						RelationType: "is_subskill_of",
						Context:      fmt.Sprintf("The skill '%s' is a sub-skill of the broader category '%s'.", skill.Name, parentSkill.Name),
						Confidence:   1.0, // This is a defined, not extracted, relationship
						Domain:       "educational_skills_taxonomy",
						CreatedAt:    time.Now(),
					}
					links = append(links, link)
					relationshipCount++
				} else {
					log.Printf("  ❌ ERROR: Parent skill ID %s not found for skill '%s'",
						parentID.Hex(), skill.Name)
				}
			}
		}
	}

	if len(links) == 0 {
		log.Println("❌ No parent-child skill relationships found to vectorize.")
		log.Println("This may be because:")
		log.Println("1. No skills have parentSkillIDs defined in the JSON file")
		log.Println("2. Skills with parentSkillIDs weren't properly loaded")
		log.Println("3. The parent skill IDs couldn't be resolved to existing skills")

		// Try to detect any skills that should have parent relationships
		detectPotentialRelationships(allSkills)

		return nil
	}

	log.Printf("✅ Found %d skill relationships to vectorize.", len(links))
	return skillService.StoreSemanticLinks(ctx, links)
}

// detectPotentialRelationships tries to identify skills that might have parent-child
// relationships based on naming patterns.
func detectPotentialRelationships(skills []models.Skill) {
	log.Println("\nAnalyzing skill names to detect potential relationships...")

	for i, skill := range skills {
		for j, potentialParent := range skills {
			// Skip self-comparison
			if i == j {
				continue
			}

			// Simple heuristic: if skill name contains another skill's name,
			// it might be a child of that skill
			if strings.Contains(strings.ToLower(skill.Name), strings.ToLower(potentialParent.Name)) {
				log.Printf("Potential relationship: '%s' might be a parent of '%s'",
					potentialParent.Name, skill.Name)
			}
		}
	}
}

// preprocessSkillIDs converts string representations of ObjectIDs in ParentSkillIDs
// into actual primitive.ObjectID types. This is necessary because the data is
// loaded from a JSON file where ObjectIDs could be strings or extended JSON format.
func preprocessSkillIDs(skillsData []models.Skill) {
	log.Println("Preprocessing ParentSkillIDs from strings to ObjectIDs...")
	skillIDMap := make(map[string]primitive.ObjectID)
	relationshipCount := 0

	// First, create a map of all skill string IDs to their ObjectID
	for i, skill := range skillsData {
		// The json unmarshaller creates a new ObjectID for `ID` but we need to preserve the original.
		// We'll parse the original string ID from the JSON and set it on the struct.
		objID := primitive.NewObjectID()
		if skill.IDString != "" {
			if parsedID, err := primitive.ObjectIDFromHex(skill.IDString); err == nil {
				objID = parsedID
			} else {
				log.Printf("Warning: could not convert _id string '%s' to ObjectID for skill '%s'. Using a new ID.", skill.IDString, skill.Name)
			}
		} else {
			log.Printf("Note: Skill '%s' has no IDString, generating new ObjectID", skill.Name)
		}
		skillsData[i].ID = objID
		skillIDMap[skill.IDString] = objID
		skillIDMap[objID.Hex()] = objID // Map both string formats

		// Also map by name as a fallback strategy
		skillIDMap[skill.Name] = objID

		// Ensure timestamps are set if not already present
		if skillsData[i].LastUpdated.IsZero() {
			skillsData[i].LastUpdated = time.Now()
		}
	}

	log.Printf("Created ID map for %d skills", len(skillIDMap))

	// Now, process parent-child relationships
	for i := range skillsData {
		if len(skillsData[i].ParentSkillIDs) > 0 {
			var resolvedIDs []primitive.ObjectID
			unresolvedCount := 0
			for _, parentID := range skillsData[i].ParentSkillIDs {
				// The parentID from JSON might be a string that needs resolving.
				// This logic is complex because the JSON unmarshalling is tricky.
				// We will rely on the string representation if the ObjectID is zero.
				if !parentID.IsZero() {
					resolvedIDs = append(resolvedIDs, parentID)
					relationshipCount++
				} else {
					unresolvedCount++
				}
			}
			skillsData[i].ParentSkillIDs = resolvedIDs
			if unresolvedCount > 0 {
				log.Printf("Skill '%s' has %d unresolved parent IDs.", skillsData[i].Name, unresolvedCount)
			}
		}
	}

	if relationshipCount == 0 {
		log.Println("⚠️ No parent-child relationships were established during preprocessing.")
		log.Println("   This could be due to the format of 'parentSkillIDs' in your JSON.")
	} else {
		log.Printf("✅ Finished preprocessing skill IDs with %d parent-child relationships",
			relationshipCount)
	}
}

func ensureEducationalSkillsClass(ctx context.Context, client *weaviate.Client, className string) error {
	// This function is deprecated as the EducationalSkills class is no longer used for vector search.
	// We ensure the core schema instead.
	return db.EnsureCoreSchema(ctx)
}

// ensureDebugBookSource seeds a known 'book' source for quick debugging across API/UI and CLI tools.
func ensureDebugBookSource(ctx context.Context, mdb *db.MongoDB) {
	now := time.Now()

	// 1) Seed into 'sources' (used by /api/sources)
	srcColl := mdb.Database.Collection("sources")
	srcFilter := bson.M{"type": "book", "isbn": "9780815345121"}
	if err := srcColl.FindOne(ctx, srcFilter).Err(); err == mongo.ErrNoDocuments {
		_, insErr := srcColl.InsertOne(ctx, bson.M{
			"type":       "book",
			"isbn":       "9780815345121",
			"title":      "Case Studies in Immunology",
			"authors":    []string{"Raif S Geha"},
			"publisher":  "WW Norton & Co",
			"year":       "2016",
			"url":        "",
			"file_path":  "",
			"processed":  false,
			"metadata":   nil,
			"created_at": now,
			"updated_at": now,
		})
		if insErr == nil {
			log.Println("Seeded debug 'book' source in 'sources' collection")
		}
	}

	// 2) Seed into 'medical_books' (used by CLI tools and legacy code)
	booksColl := mdb.Database.Collection("medical_books")
	bookFilter := bson.M{"$or": []bson.M{{"isbn13": "9780815345121"}, {"isbn": "9780815345121"}}}
	if err := booksColl.FindOne(ctx, bookFilter).Err(); err == mongo.ErrNoDocuments {
		_, insErr := booksColl.InsertOne(ctx, bson.M{
			"isbn13":     "9780815345121",
			"isbn":       "9780815345121",
			"title":      "Case Studies in Immunology",
			"authors":    []string{"Raif S Geha"},
			"publisher":  "WW Norton & Co",
			"year":       "2016",
			"created_at": now,
			"updated_at": now,
		})
		if insErr == nil {
			log.Println("Seeded debug 'book' in 'medical_books' collection")
		}
	}
}
