package main

import (
	"context"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"esp-organizer/internal/config"
	"esp-organizer/internal/models"
	"esp-organizer/internal/utils"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Initialize configuration
	projectRoot, err := utils.FindProjectRoot()
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}
	config.LoadConfig(projectRoot, ".env")

	ctx := context.Background()

	// Connect to MongoDB
	mngoDB, err := db.NewMongoDBFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mngoDB.Client.Disconnect(ctx); err != nil {
			log.Printf("Warning: Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Connect to Weaviate
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}
	weaviateClient := db.GetWeaviateClient()
	log.Println("✅ Connected to databases")

	// Parse command line arguments
	if len(os.Args) > 1 && os.Args[1] == "fix" {
		log.Println("🔧 Running in FIX mode - will attempt to repair relationships")
		fixRelationships(ctx, mngoDB.Database)
		return
	}

	// Run diagnostics
	log.Println("🔍 Running relationship diagnostics")
	runDiagnostics(ctx, mngoDB.Database, weaviateClient)
}

func runDiagnostics(ctx context.Context, db *mongo.Database, weaviateClient interface{}) {
	// 1. Count total skills
	skillsCollection := db.Collection("skills")
	totalSkills, err := skillsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to count skills: %v", err)
	}
	log.Printf("Total skills in MongoDB: %d", totalSkills)

	// 2. Count skills with parent relationships
	parentFilter := bson.M{"parentSkillIDs.0": bson.M{"$exists": true}}
	skillsWithParents, err := skillsCollection.CountDocuments(ctx, parentFilter)
	if err != nil {
		log.Fatalf("Failed to count skills with parents: %v", err)
	}
	log.Printf("Skills with parent relationships: %d (%.1f%%)",
		skillsWithParents, float64(skillsWithParents)/float64(totalSkills)*100)

	// 3. Count semantic links in MongoDB
	semanticLinksCollection := db.Collection("semantic_links")
	totalLinks, err := semanticLinksCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Warning: Failed to count semantic links: %v", err)
	} else {
		log.Printf("Total semantic links in MongoDB: %d", totalLinks)
	}

	// 4. Check skills with broken parent references
	cursor, err := skillsCollection.Find(ctx, parentFilter)
	if err != nil {
		log.Fatalf("Failed to query skills with parents: %v", err)
	}
	defer cursor.Close(ctx)

	var brokenRelationships int
	var skills []models.Skill
	if err := cursor.All(ctx, &skills); err != nil {
		log.Fatalf("Failed to decode skills: %v", err)
	}

	// Build a map of all skill IDs for quick lookup
	allSkillsMap := make(map[primitive.ObjectID]bool)
	allSkillsCursor, err := skillsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to query all skills: %v", err)
	}
	defer allSkillsCursor.Close(ctx)

	var allSkills []models.Skill
	if err := allSkillsCursor.All(ctx, &allSkills); err != nil {
		log.Fatalf("Failed to decode all skills: %v", err)
	}

	for _, skill := range allSkills {
		allSkillsMap[skill.ID] = true
	}

	// Check for broken parent references
	log.Println("Checking for broken parent references...")
	for _, skill := range skills {
		for _, parentID := range skill.ParentSkillIDs {
			if !allSkillsMap[parentID] {
				brokenRelationships++
				log.Printf("⚠️ Skill '%s' has invalid parent ID: %s", skill.Name, parentID.Hex())
			}
		}
	}

	if brokenRelationships > 0 {
		log.Printf("⚠️ Found %d broken parent references", brokenRelationships)
	} else {
		log.Printf("✅ All parent references are valid")
	}

	// 5. Analyze semantic links
	if totalLinks > 0 {
		analyzeLinkTypes(ctx, semanticLinksCollection)
		analyzeLinkDomains(ctx, semanticLinksCollection)
	}

	// 6. Check for disconnected skills (no parents or children)
	log.Println("Checking for disconnected skills...")
	disconnectedFilter := bson.M{
		"$and": []bson.M{
			{"parentSkillIDs": bson.M{"$size": 0}},
			{"childSkillIDs": bson.M{"$size": 0}},
		},
	}

	disconnectedCount, err := skillsCollection.CountDocuments(ctx, disconnectedFilter)
	if err != nil {
		log.Printf("Warning: Failed to count disconnected skills: %v", err)
	} else {
		log.Printf("Disconnected skills (no parents or children): %d (%.1f%%)",
			disconnectedCount, float64(disconnectedCount)/float64(totalSkills)*100)
	}

	// Print summary
	log.Println("\n=== Relationship Diagnostics Summary ===")
	log.Printf("Total Skills: %d", totalSkills)
	log.Printf("Skills with Parents: %d (%.1f%%)",
		skillsWithParents, float64(skillsWithParents)/float64(totalSkills)*100)
	log.Printf("Semantic Links: %d", totalLinks)
	log.Printf("Broken References: %d", brokenRelationships)
	log.Printf("Disconnected Skills: %d (%.1f%%)",
		disconnectedCount, float64(disconnectedCount)/float64(totalSkills)*100)

	if skillsWithParents == 0 {
		log.Println("\n❌ CRITICAL ISSUE: No parent-child relationships found")
		log.Println("This is likely why no semantic links are being created.")
		log.Println("Check the skills.json file to ensure parent_skill_ids are properly defined.")
		log.Println("Try running with the 'fix' parameter to attempt automatic repair.")
	} else if brokenRelationships > 0 {
		log.Println("\n⚠️ WARNING: Some parent references are broken")
		log.Println("This may cause semantic links to be incomplete.")
		log.Println("Try running with the 'fix' parameter to attempt repair.")
	} else if totalLinks == 0 {
		log.Println("\n⚠️ WARNING: No semantic links found despite having valid parent relationships")
		log.Println("This suggests an issue in the vectorization or link creation process.")
	} else {
		log.Println("\n✅ No major relationship issues detected")
	}
}

func analyzeLinkTypes(ctx context.Context, collection *mongo.Collection) {
	// Get counts of different relation types
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$RelationType"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Warning: Failed to aggregate link types: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Printf("Warning: Failed to decode link types: %v", err)
		return
	}

	log.Println("Semantic link types:")
	for _, result := range results {
		linkType := result["_id"]
		count := result["count"]
		log.Printf("  - %v: %v", linkType, count)
	}
}

func analyzeLinkDomains(ctx context.Context, collection *mongo.Collection) {
	// Get counts of different domains
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$Domain"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Warning: Failed to aggregate link domains: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Printf("Warning: Failed to decode link domains: %v", err)
		return
	}

	log.Println("Semantic link domains:")
	for _, result := range results {
		domain := result["_id"]
		count := result["count"]
		log.Printf("  - %v: %v", domain, count)
	}
}

func fixRelationships(ctx context.Context, db *mongo.Database) {
	log.Println("Starting relationship repair process...")
	skillsCollection := db.Collection("skills")

	// Load all skills
	cursor, err := skillsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to query skills: %v", err)
	}
	defer cursor.Close(ctx)

	var skills []models.Skill
	if err := cursor.All(ctx, &skills); err != nil {
		log.Fatalf("Failed to decode skills: %v", err)
	}

	// Build ID map and name map for lookups
	idMap := make(map[primitive.ObjectID]*models.Skill)
	nameMap := make(map[string]*models.Skill)

	for i := range skills {
		skill := &skills[i]
		idMap[skill.ID] = skill
		nameMap[strings.ToLower(skill.Name)] = skill
	}

	// Find and fix parent-child relationships
	updatedCount := 0

	for i := range skills {
		skill := &skills[i]

		// Skip if already has parents
		if len(skill.ParentSkillIDs) > 0 {
			continue
		}

		// Look for pattern matches in names to guess relationships
		for _, potentialParent := range skills {
			// Skip self
			if potentialParent.ID == skill.ID {
				continue
			}

			// Simple relationship detection by name substring
			// This is a heuristic approach - customize as needed for your data
			if strings.Contains(strings.ToLower(skill.Name), strings.ToLower(potentialParent.Name)) {
				log.Printf("Found potential parent relationship: '%s' could be parent of '%s'",
					potentialParent.Name, skill.Name)

				// Update the skill
				update := bson.M{
					"$addToSet": bson.M{"parentSkillIDs": potentialParent.ID},
				}

				_, err := skillsCollection.UpdateOne(ctx, bson.M{"_id": skill.ID}, update)
				if err != nil {
					log.Printf("Error updating skill relationship: %v", err)
				} else {
					updatedCount++
				}
			}
		}
	}

	log.Printf("Repair complete. Updated %d skills with new parent relationships.", updatedCount)

	if updatedCount > 0 {
		log.Println("You should now run the upload-skills tool again to recreate semantic links.")
	}
}
