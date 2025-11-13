package utils

import (
	"context"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"esp-organizer/internal/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FindBookByID retrieves a book/source by ID from MongoDB
// Tries multiple ID formats (string, ObjectID, ISBN)
func FindBookByID(sourceID string) *models.Source {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("⚠️ Warning: Could not connect to MongoDB: %v", err)
		return nil
	}
	defer mongoDb.Client.Disconnect(ctx)

	collection := mongoDb.Database.Collection("sources")

	// Try both string ID and ObjectID formats
	var objectID primitive.ObjectID
	if len(sourceID) == 24 {
		objectID, err = primitive.ObjectIDFromHex(sourceID)
		if err != nil {
			objectID = primitive.NilObjectID
		}
	}

	filter := bson.M{"$or": []bson.M{
		{"_id": sourceID},
		{"_id": objectID},
	}}

	var source models.Source
	err = collection.FindOne(ctx, filter).Decode(&source)
	if err != nil {
		return nil
	}

	return &source
}

// FindBookByISBN retrieves a book by ISBN
func FindBookByISBN(isbn string) *models.Source {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("⚠️ Warning: Could not connect to MongoDB: %v", err)
		return nil
	}
	defer mongoDb.Client.Disconnect(ctx)

	collection := mongoDb.Database.Collection("sources")
	var source models.Source
	err = collection.FindOne(ctx, bson.M{"isbn": isbn}).Decode(&source)
	if err != nil {
		log.Printf("Error finding book by ISBN %s: %v", isbn, err)
		return nil
	}

	return &source
}

// MaskConnectionString hides sensitive parts of connection strings
func MaskConnectionString(uri string) string {
	if uri == "" {
		return ""
	}

	// Hide passwords in connection strings
	if idx := len(uri) - 1; idx > 0 {
		parts := make([]rune, len(uri))
		copy(parts, []rune(uri))

		inPassword := false
		for i, ch := range parts {
			if ch == ':' && i > 0 && parts[i-1] != '@' {
				inPassword = true
			}
			if ch == '@' {
				inPassword = false
			}
			if inPassword && ch != ':' && ch != '@' {
				parts[i] = '*'
			}
		}
		return string(parts)
	}

	return uri
}
