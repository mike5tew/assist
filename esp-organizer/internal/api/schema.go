package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"esp-organizer/internal/InfoFlow/InfoStore/db"

	"github.com/gin-gonic/gin"
	"github.com/weaviate/weaviate/entities/models"
)

// CreateSemanticLinksClass creates a models.Class for semantic links
func CreateSemanticLinksClass() *models.Class {
	return &models.Class{
		Class: "SemanticLinks",
		Properties: []*models.Property{
			{Name: "sourceTerm", DataType: []string{"string"}},
			{Name: "targetTerm", DataType: []string{"string"}},
			{Name: "relationship", DataType: []string{"string"}},
			{Name: "confidence", DataType: []string{"number"}},
			{Name: "domain", DataType: []string{"string"}},
			{Name: "sourceDocId", DataType: []string{"string"}}, // Reference to MongoDB document
			{Name: "context", DataType: []string{"text"}},
		},
	}
}

// CreateMedicalExcerptClass creates a models.Class for medical excerpts
func CreateMedicalExcerptClass() *models.Class {
	return &models.Class{
		Class: "MedicalExcerpt",
		Properties: []*models.Property{
			{Name: "content", DataType: []string{"text"}},
			{Name: "source", DataType: []string{"string"}},
			{Name: "documentId", DataType: []string{"string"}}, // Reference to MongoDB
			{Name: "domain", DataType: []string{"string"}},
			{Name: "chapterNumber", DataType: []string{"string"}},
			{Name: "chapterTitle", DataType: []string{"string"}},
		},
	}
}

// InitializeSchema creates required Weaviate classes if they don't exist
func InitializeSchema(c *gin.Context) {
	client := db.GetWeaviateClient()
	ctx := context.Background()

	// Create classes
	classes := []*models.Class{
		CreateSemanticLinksClass(),
		CreateMedicalExcerptClass(),
	}

	for _, class := range classes {
		// Fix: Use ClassExistenceChecker instead of ClassExistence
		exists, err := client.Schema().ClassExistenceChecker().WithClassName(class.Class).Do(ctx)
		if err != nil {
			log.Printf("Error checking if class %s exists: %v", class.Class, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !exists {
			if err := client.Schema().ClassCreator().WithClass(class).Do(ctx); err != nil {
				log.Printf("Error creating class %s: %v", class.Class, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			log.Printf("Created Weaviate class: %s", class.Class)
		} else {
			log.Printf("Weaviate class already exists: %s", class.Class)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Schema initialized successfully",
		"schema": map[string]interface{}{
			"classes": []string{"SemanticLinks", "MedicalExcerpt"},
		},
	})
}

// InitializeSemanticLinksSchema initializes only the semantic links schema
// Can be called during application startup
func InitializeSemanticLinksSchema() error {
	ctx := context.Background()
	client := db.GetWeaviateClient()

	class := CreateSemanticLinksClass()
	className := class.Class

	// Fix: Use ClassExistenceChecker instead of ClassExistence
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to check class %s: %w", className, err)
	}

	if exists {
		log.Printf("Weaviate class already exists: %s", className)
		return nil
	}

	// Create class if it doesn't exist
	err = client.Schema().ClassCreator().WithClass(class).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class %s: %w", className, err)
	}

	log.Printf("Successfully created Weaviate class: %s", className)
	return nil
}

// GetCurrentConfig shows how to retrieve the current class configuration
func GetCurrentConfig(c *gin.Context) {
	ctx := context.Background()
	client := db.GetWeaviateClient()

	className := c.Param("className")
	if className == "" {
		className = "SemanticLinks"
	}

	// Get the class schema to see current configuration
	class, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if class == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Class not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"class":             class.Class,
		"vectorIndexType":   class.VectorIndexType,
		"vectorIndexConfig": class.VectorIndexConfig,
		"vectorizer":        class.Vectorizer,
		"moduleConfig":      class.ModuleConfig,
	})
}

// UpdateClassConfig updates the entire class configuration including vector index settings
func UpdateClassConfig(c *gin.Context) {
	ctx := context.Background()
	client := db.GetWeaviateClient()

	className := "SemanticLinks"

	// First get the current class to preserve existing properties
	currentClass, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current class: " + err.Error()})
		return
	}

	// Create updated class configuration
	updatedClass := &models.Class{
		Class:      className,
		Vectorizer: currentClass.Vectorizer,
		ModuleConfig: map[string]interface{}{
			"text2vec-aws": map[string]interface{}{
				"region":  "eu-west-2",
				"service": "bedrock",
				"model":   "amazon.titan-embed-text-v1",
			},
		},
		VectorIndexType: "hnsw",
		VectorIndexConfig: map[string]interface{}{
			"distance":         "cosine",
			"efConstruction":   128,
			"maxConnections":   64,
			"ef":               200, // Updated value
			"dynamicEfFactor":  8,
			"dynamicEfMin":     100,
			"dynamicEfMax":     500,
			"flatSearchCutoff": 40000,
		},
		Properties: currentClass.Properties, // Preserve existing properties
	}

	// Update the class configuration
	err = client.Schema().ClassUpdater().WithClass(updatedClass).Do(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update class: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Class configuration updated successfully",
		"class":   className,
	})
}

// IMPORTANT: Note that updating vector index config may require reindexing
// and could impact performance during the update process

// Alternative: Create a new class with the desired configuration and migrate data
// func RecreateClassWithNewConfig(c *gin.Context) {
// 	ctx := context.Background()
// 	client := db.GetWeaviateClient()

// 	className := "SemanticLinks"
// 	tempClassName := "SemanticLinks_temp"

// 	// This is a complex operation that would involve:
// 	// 1. Creating a temporary class with new configuration
// 	// 2. Copying data from old class to new class
// 	// 3. Deleting the old class
// 	// 4. Renaming the temporary class

// 	// This approach is recommended for major configuration changes
// 	// but requires careful implementation

// 	c.JSON(http.StatusNotImplemented, gin.H{
// 		"message": "This operation requires careful implementation",
// 		"warning": "Recreating classes with data migration is complex",
// 	})
// }
