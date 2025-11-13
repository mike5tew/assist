package routes

import (
	"context"
	"esp-organizer/internal/InfoFlow/InfoIn"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RegisterSampleManagementRoutes registers all sample management endpoints
func RegisterSampleManagementRoutes(router *gin.Engine) {
	sampleRoutes := router.Group("/api/samples")
	{
		sampleRoutes.GET("", ListSamples)
		sampleRoutes.GET("/:id", GetSample)
		sampleRoutes.POST("", CreateSample)
		sampleRoutes.PUT("/:id", UpdateSample)
		sampleRoutes.DELETE("/:id", DeleteSample)
		sampleRoutes.POST("/curate", TriggerCuration)
		sampleRoutes.GET("/quality/:domain", GetQualityMetrics)
	}
}

// ListSamples retrieves extraction samples with filtering
func ListSamples(c *gin.Context) {
	ctx := context.Background()
	mongoDB, err := db.NewFromEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer mongoDB.Client.Disconnect(ctx)

	collection := mongoDB.Database.Collection("extraction_samples")

	// Build filter from query params
	filter := bson.M{}
	if domain := c.Query("domain"); domain != "" {
		filter["domain"] = domain
	}
	if minQuality := c.Query("min_quality"); minQuality != "" {
		filter["quality_score"] = bson.M{"$gte": minQuality}
	}

	// Pagination
	page := c.DefaultQuery("page", "1")
	//pageSize := c.DefaultQuery("page_size", "20")

	opts := options.Find().
		SetSort(bson.D{{Key: "quality_score", Value: -1}}).
		SetLimit(20)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch samples"})
		return
	}
	defer cursor.Close(ctx)

	var samples []InfoIn.ExtractionSample
	if err := cursor.All(ctx, &samples); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode samples"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"samples": samples,
		"count":   len(samples),
		"page":    page,
	})
}

// GetSample retrieves a single sample by ID
func GetSample(c *gin.Context) {
	ctx := context.Background()
	mongoDB, err := db.NewFromEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer mongoDB.Client.Disconnect(ctx)

	sampleID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(sampleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sample ID"})
		return
	}

	collection := mongoDB.Database.Collection("extraction_samples")
	var sample InfoIn.ExtractionSample

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&sample)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sample not found"})
		return
	}

	c.JSON(http.StatusOK, sample)
}

// CreateSample creates a new extraction sample
func CreateSample(c *gin.Context) {
	ctx := context.Background()
	mongoDB, err := db.NewFromEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer mongoDB.Client.Disconnect(ctx)

	var sample InfoIn.ExtractionSample
	if err := c.BindJSON(&sample); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set metadata
	sample.CreatedAt = time.Now()

	collection := mongoDB.Database.Collection("extraction_samples")
	result, err := collection.InsertOne(ctx, sample)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sample"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      result.InsertedID,
		"message": "Sample created successfully",
	})
}

// UpdateSample updates an existing sample
func UpdateSample(c *gin.Context) {
	ctx := context.Background()
	mongoDB, err := db.NewFromEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer mongoDB.Client.Disconnect(ctx)

	sampleID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(sampleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sample ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.BindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	collection := mongoDB.Database.Collection("extraction_samples")
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": updates},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sample"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sample not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sample updated successfully"})
}

// DeleteSample deletes a sample by ID
func DeleteSample(c *gin.Context) {
	ctx := context.Background()
	mongoDB, err := db.NewFromEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer mongoDB.Client.Disconnect(ctx)

	sampleID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(sampleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sample ID"})
		return
	}

	collection := mongoDB.Database.Collection("extraction_samples")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sample"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sample not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sample deleted successfully"})
}

// TriggerCuration manually triggers sample curation from high-quality extractions
func TriggerCuration(c *gin.Context) {
	ctx := context.Background()

	semanticService, err := InfoIn.NewSemanticLinkService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize service"})
		return
	}

	// Run curation in background
	go func() {
		if err := semanticService.CurateHighQualityExtractions(ctx); err != nil {
			// Log error but don't block response
			// log.Printf("Curation error: %v", err)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Curation job started",
		"status":  "processing",
	})
}

// GetQualityMetrics returns quality statistics for a domain
func GetQualityMetrics(c *gin.Context) {
	ctx := context.Background()
	mongoDB, err := db.NewFromEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer mongoDB.Client.Disconnect(ctx)

	domain := c.Param("domain")
	collection := mongoDB.Database.Collection("extraction_samples")

	// Aggregate quality metrics
	pipeline := []bson.M{
		{"$match": bson.M{"domain": domain}},
		{"$group": bson.M{
			"_id":         nil,
			"avg_quality": bson.M{"$avg": "$quality_score"},
			"min_quality": bson.M{"$min": "$quality_score"},
			"max_quality": bson.M{"$max": "$quality_score"},
			"count":       bson.M{"$sum": 1},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compute metrics"})
		return
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode metrics"})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"domain":  domain,
			"metrics": "No samples found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domain":  domain,
		"metrics": results[0],
	})
}
