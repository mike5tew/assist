package db

import (
	"context"
	"log"
)

// FilterTimestampFields removes timestamp fields that might cause Weaviate schema validation errors
// This function can be used by any code that needs to store data in Weaviate
func FilterTimestampFields(properties map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// List of fields that commonly cause schema validation errors
	fieldsToFilter := []string{
		"created_at",
		"updated_at",
		"createdAt",
		"updatedAt",
		"timestamp",
	}

	// Copy all fields except timestamp fields that might cause problems
	for key, value := range properties {
		// Skip fields in our filter list
		shouldSkip := false
		for _, field := range fieldsToFilter {
			if key == field {
				shouldSkip = true
				log.Printf("Filtering out timestamp field '%s' for Weaviate storage", key)
				break
			}
		}

		if shouldSkip {
			continue
		}

		// Handle nested maps recursively
		if nestedMap, ok := value.(map[string]interface{}); ok {
			result[key] = FilterTimestampFields(nestedMap)
		} else {
			result[key] = value
		}
	}

	return result
}

// StoreDocumentWithVectorSafe is a wrapper around StoreDocumentWithVector that filters timestamp fields
func StoreDocumentWithVectorSafe(className string, properties map[string]interface{}, vector []float32) error {
	// Filter out timestamp fields that might cause schema validation errors
	filteredProperties := FilterTimestampFields(properties)

	// Use the existing function with filtered properties
	return StoreDocumentWithVector(context.TODO(), className, filteredProperties, vector)
}
