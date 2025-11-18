package main

import (
	"context"
	"log"

	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	ctx := context.Background()
	mongoDB, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	collection := mongoDB.Database.Collection("semantic_links")

	// Example: Migrate old "personal" domain to split into "personal" and "students"
	cursor, err := collection.Find(ctx, bson.M{
		"domain":      "personal",
		"source_term": bson.M{"$regex": ".*student.*", "$options": "i"},
	})
	if err != nil {
		log.Fatalf("Failed to find links: %v", err)
	}
	defer cursor.Close(ctx)

	updateCount := 0
	for cursor.Next(ctx) {
		var link models.SemanticLink
		if err := cursor.Decode(&link); err != nil {
			log.Printf("Warning: Failed to decode link: %v", err)
			continue
		}

		// Update domain to "students"
		_, err := collection.UpdateOne(ctx,
			bson.M{"_id": link.ID},
			bson.M{"$set": bson.M{"domain": "students"}},
		)
		if err != nil {
			log.Printf("Warning: Failed to update link %s: %v", link.ID.Hex(), err)
			continue
		}
		updateCount++
	}

	log.Printf("✅ Migrated %d links from 'personal' to 'students' domain", updateCount)
}
