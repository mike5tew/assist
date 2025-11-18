package main

import (
	"context"
	"log"
	"time"

	"esp-organizer/internal/store/db"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	log.Println("=== ESP Organizer - Working Features Demo ===")

	ctx := context.Background()

	// Test 1: MongoDB Connection
	log.Println("\n1️⃣ Testing MongoDB Connection...")
	mongoDB, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}
	log.Println("✅ MongoDB connected")

	// Test 2: Check Immunology Content
	log.Println("\n2️⃣ Checking Immunology Content...")
	collection := mongoDB.Database.Collection("immunology_content")
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalf("❌ Failed to count documents: %v", err)
	}
	log.Printf("✅ Found %d immunology documents", count)

	// Test 3: Weaviate Connection
	log.Println("\n3️⃣ Testing Weaviate Connection...")
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("❌ Weaviate connection failed: %v", err)
	}
	log.Println("✅ Weaviate connected")

	// Test 4: Check SemanticLinks Class
	log.Println("\n4️⃣ Checking Weaviate SemanticLinks Class...")
	exists, err := db.WeaviateCollectionExists(ctx, "SemanticLinks")
	if err != nil {
		log.Fatalf("❌ Failed to check Weaviate class: %v", err)
	}
	if !exists {
		log.Println("⚠️  SemanticLinks class doesn't exist, creating...")
		if err := db.CreateSemanticLinksClass(ctx); err != nil {
			log.Fatalf("❌ Failed to create class: %v", err)
		}
	}
	log.Println("✅ SemanticLinks class exists")

	// Summary
	log.Println("\n=== Demo Complete ===")
	log.Println("✅ All core systems are operational")
	log.Println("")
	log.Println("📋 Next Steps:")
	log.Println("   1. Upload a PDF: make run-full-test PDF=./resources/test.pdf")
	log.Println("   2. Query content: make test-hsg-search QUERY=\"immunology\"")
	log.Println("   3. Check API health: curl http://localhost:8080/health")

	time.Sleep(1 * time.Second)
}
