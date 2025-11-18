package main

import (
	"context"
	"esp-organizer/internal/domain/infoin/extraction"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// Parse command line flags
	pdfPath := flag.String("pdf", "", "Path to the PDF file to process")
	isbn := flag.String("isbn", "", "ISBN of the book")
	chapterNum := flag.Int("chapter", 1, "Chapter number to process")
	domain := flag.String("domain", "immunology", "Domain for the content (e.g. immunology)")
	flag.Parse()

	if *pdfPath == "" || *isbn == "" {
		log.Fatal("PDF path and ISBN are required")
	}

	// Connect to MongoDB
	log.Println("Connecting to MongoDB...")
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDb.Client.Disconnect(context.Background())

	// Find the book by ISBN
	var book models.Source
	err = mongoDb.Database.Collection("medical_books").FindOne(
		context.Background(),
		bson.M{"isbn": *isbn},
	).Decode(&book)
	if err != nil {
		log.Fatalf("Failed to find book with ISBN %s: %v", *isbn, err)
	}

	// Open the PDF file
	file, err := os.Open(*pdfPath)
	if err != nil {
		log.Fatalf("Failed to open PDF file: %v", err)
	}
	defer file.Close()

	// Set up the chapter info
	chapterInfo := models.ChapterInfo{
		ChapterNumber: fmt.Sprintf("%d", *chapterNum),
		ChapterTitle:  fmt.Sprintf("Chapter %d", *chapterNum),
		Domain:        *domain,
	}

	// Initialize text extraction processor
	processor, err := extraction.NewTextractProcessor(
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_S3_BUCKET"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Textract processor: %v", err)
	}

	jobID := fmt.Sprintf("textract-%d", time.Now().Unix())
	log.Printf("Starting extraction job %s for chapter %d", jobID, *chapterNum)

	// Process the PDF file
	job, err := processor.ProcessChapterFile(
		context.Background(),
		file,
		chapterInfo,
		book,
		jobID,
	)
	if err != nil {
		log.Fatalf("Failed to process chapter file: %v", err)
	}

	// Store the extracted content in MongoDB
	log.Printf("Storing extracted content in MongoDB...")
	if job.ExtractedData == nil {
		log.Fatalf("No data was extracted from the PDF")
	}

	// Create a document for the extracted chapter
	chapterDoc := bson.M{
		"chapter_number": chapterInfo.ChapterNumber,
		"chapter_title":  chapterInfo.ChapterTitle,
		"domain":         *domain,
		"book_id":        book.ID,
		"content":        job.ExtractedData.ChapterContent,
		"raw_text":       job.ExtractedData.RawText,
		"created_at":     primitive.NewDateTimeFromTime(time.Now()),
		"content_type":   "chapter",
		"source": bson.M{
			"id":    book.ID,
			"title": book.Title,
			"isbn":  book.ISBN,
		},
	}

	// Insert into immunology_content collection
	result, err := mongoDb.Database.Collection("immunology_content").InsertOne(
		context.Background(),
		chapterDoc,
	)
	if err != nil {
		log.Fatalf("Failed to insert chapter content: %v", err)
	}
	log.Printf("Inserted chapter document with ID %v", result.InsertedID)

	// Update extraction job status
	_, err = mongoDb.Database.Collection("extraction_jobs").InsertOne(
		context.Background(),
		job,
	)
	if err != nil {
		log.Printf("Warning: Failed to update job status: %v", err)
	}

	log.Printf("✅ Successfully processed chapter %d from %s", *chapterNum, book.Title)
}
