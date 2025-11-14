package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UploadResponse struct {
	BatchID string `json:"batch_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type StatusResponse struct {
	JobID         string  `json:"job_id"`
	Status        string  `json:"status"`
	Message       string  `json:"message"`
	CurrentStep   string  `json:"current_step,omitempty"`
	Detail        string  `json:"detail,omitempty"`
	Progress      float64 `json:"progress,omitempty"`
	TextractJobID string  `json:"textract_job_id,omitempty"`
}

func main() {
	pdfPath := flag.String("pdf", "", "Path to the PDF file to upload and process.")
	apiURL := flag.String("api-url", "http://localhost:8080", "URL of the API server.")
	timeout := flag.Duration("timeout", 10*time.Minute, "Timeout for waiting for processing to complete.")
	flag.Parse()

	if *pdfPath == "" {
		fmt.Println("❌ Error: -pdf flag is required.")
		os.Exit(1)
	}

	// 1. Upload the PDF
	fmt.Printf("🚀 Starting full test workflow for: %s\n", *pdfPath)
	uploadURL := *apiURL + "/api/immunology/upload-chapter"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(*pdfPath)
	if err != nil {
		fmt.Printf("❌ Error opening PDF file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("chapter", filepath.Base(*pdfPath))
	if err != nil {
		fmt.Printf("❌ Error creating form file: %v\n", err)
		os.Exit(1)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("❌ Error copying file data: %v\n", err)
		os.Exit(1)
	}

	// Add all required form fields
	_ = writer.WriteField("chapter_title", "Test Chapter")
	_ = writer.WriteField("chapter_number", "1")
	_ = writer.WriteField("source_id", "")
	_ = writer.WriteField("book_title", "Test Book")
	_ = writer.WriteField("book_isbn", "123456789")
	_ = writer.WriteField("book_authors", "Test Author")
	_ = writer.WriteField("book_publisher", "Test Publisher")
	_ = writer.WriteField("book_year", "2023")
	_ = writer.WriteField("extraction_method", "manual")
	_ = writer.WriteField("content_type", "textbook")
	_ = writer.WriteField("target_exam", "test")
	_ = writer.WriteField("domain", "immunology")

	err = writer.Close()
	if err != nil {
		fmt.Printf("❌ Error closing multipart writer: %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		fmt.Printf("❌ Error creating request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	fmt.Println("1️⃣  Uploading PDF to API...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error during upload request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Error reading upload response body: %v\n", err)
		os.Exit(1)
	}

	// Check for non-successful status codes (anything not in the 200-299 range)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("❌ Upload failed with status %d: %s\n", resp.StatusCode, string(respBody))
		os.Exit(1)
	}

	var uploadResp UploadResponse
	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		fmt.Printf("❌ Error decoding upload response: %v. Body: %s\n", err, string(respBody))
		os.Exit(1)
	}

	if uploadResp.BatchID == "" {
		fmt.Printf("❌ Failed to get batch ID from upload response. Body: %s\n", string(respBody))
		os.Exit(1)
	}

	fmt.Printf("✅ Upload successful. Batch ID: %s\n", uploadResp.BatchID)

	// 2. Poll for status
	fmt.Println("2️⃣  Monitoring processing status...")
	statusURL := *apiURL + "/api/extraction/status/" + uploadResp.BatchID
	startTime := time.Now()

	for {
		if time.Since(startTime) > *timeout {
			fmt.Println("❌ Error: Timed out waiting for processing to complete.")
			os.Exit(1)
		}

		statusResp, err := http.Get(statusURL)
		if err != nil {
			fmt.Printf("⚠️ Warning: Could not get status: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var status StatusResponse
		bodyBytes, _ := io.ReadAll(statusResp.Body)
		if err := json.Unmarshal(bodyBytes, &status); err != nil {
			statusResp.Body.Close()
			fmt.Printf("⚠️ Warning: Could not decode status response: %v. Body: %s\n", err, string(bodyBytes))
			time.Sleep(5 * time.Second)
			continue
		}
		statusResp.Body.Close()

		fmt.Printf("   [%s] Status: %s, Step: %s, Progress: %.1f%%, Detail: %s, TextractID: %s\n",
			time.Now().Format("15:04:05"), status.Status, status.CurrentStep, status.Progress, status.Detail, status.TextractJobID)

		if status.Status == "COMPLETED" {
			fmt.Println("✅ Processing completed successfully!")
			break
		}
		if status.Status == "FAILED" {
			fmt.Printf("❌ Processing failed. Final status: %+v\n", status)
			os.Exit(1)
		}

		time.Sleep(5 * time.Second)
	}

	fmt.Println("🎉 Full test workflow finished successfully.")
}
