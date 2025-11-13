package main

import (
	"esp-organizer/internal/InfoFlow/InfoIn/extraction"
	"esp-organizer/internal/models"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	fmt.Println("=== AWS Textract Setup Test ===")

	// Check AWS credentials
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	awsBucket := os.Getenv("AWS_S3_BUCKET")

	fmt.Printf("AWS Access Key: %s\n", maskCredential(awsAccessKey))
	fmt.Printf("AWS Secret Key: %s\n", maskCredential(awsSecretKey))
	fmt.Printf("AWS Region: %s\n", awsRegion)
	fmt.Printf("AWS S3 Bucket: %s\n", awsBucket)

	if awsAccessKey == "" || awsSecretKey == "" {
		fmt.Println("❌ AWS credentials not configured")
		fmt.Println("   Add AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY to .env")
		return
	}

	if awsRegion == "" {
		awsRegion = "eu-west-2"
		fmt.Printf("⚠️  Using default region: %s\n", awsRegion)
	}

	if awsBucket == "" {
		awsBucket = "esp-new-organizer-immunology"
		fmt.Printf("⚠️  Using default bucket: %s\n", awsBucket)
	}

	// Test Textract processor initialization
	fmt.Println("\n=== Testing Textract Processor ===")
	_, err := extraction.NewTextractProcessor(awsRegion, awsBucket)
	if err != nil {
		fmt.Printf("❌ Failed to initialize Textract processor: %v\n", err)
		return
	}

	fmt.Println("✅ Textract processor initialized successfully")

	// Test with sample chapter info
	chapterInfo := models.ChapterInfo{
		ChapterNumber: "1",
		ChapterTitle:  "Case 1: X-Linked Agammaglobulinemia",
	}

	bookSource := models.Source{
		Title:     "Case Studies in Immunology: A Clinical Companion",
		ISBN:      "978-0-8153-4532-1",
		Authors:   []string{"Charles A. Janeway Jr.", "Paul Travers"},
		Publisher: "Garland Science",
		Year:      "2019",
		Type:      "medical_textbook",
	}

	fmt.Printf("\n=== Sample Data for Processing ===\n")
	fmt.Printf("Book Title: %s\n", bookSource.Title)
	fmt.Printf("Book ISBN: %s\n", bookSource.ISBN)
	fmt.Printf("Publisher: %s\n", bookSource.Publisher)
	fmt.Printf("Year: %s\n", bookSource.Year)
	fmt.Printf("Book Type: %s\n", bookSource.Type)
	fmt.Printf("Book Authors: %s\n", strings.Join(bookSource.Authors, ", "))
	fmt.Printf("Chapter to Process: Chapter %s - %s\n", chapterInfo.ChapterNumber, chapterInfo.ChapterTitle)

	fmt.Println("\n=== Setup Complete ===")
	fmt.Println("Ready to process immunology chapters!")
}

func maskCredential(credential string) string {
	if credential == "" {
		return "<not set>"
	}
	if len(credential) <= 8 {
		return strings.Repeat("*", len(credential))
	}
	return credential[:4] + strings.Repeat("*", len(credential)-8) + credential[len(credential)-4:]
}
