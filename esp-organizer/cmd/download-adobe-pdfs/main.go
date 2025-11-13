package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"esp-organizer/internal/adobe"
)

func main() {
	linksFile := flag.String("links", "resources/immunology/adobe_links.txt", "Path to file containing Adobe share links")
	outputDir := flag.String("output", "resources/immunology/case_studies", "Output directory for downloaded PDFs")
	flag.Parse()

	// Read links from file
	file, err := os.Open(*linksFile)
	if err != nil {
		log.Fatalf("Failed to open links file: %v", err)
	}
	defer file.Close()

	var links []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		links = append(links, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading links file: %v", err)
	}

	fmt.Printf("📚 Found %d Adobe share links to download\n", len(links))

	// Create downloader
	downloader := adobe.NewAdobeDownloader("", "") // ClientID/Secret if needed

	// Download all files
	if err := downloader.DownloadBatch(links, *outputDir); err != nil {
		log.Fatalf("Batch download failed: %v", err)
	}

	fmt.Println("✅ All downloads complete!")
}
