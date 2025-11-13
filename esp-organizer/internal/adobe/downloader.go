package adobe

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type AdobeDownloader struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
}

func NewAdobeDownloader(clientID, clientSecret string) *AdobeDownloader {
	return &AdobeDownloader{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// GetAccessToken retrieves OAuth token from Adobe
func (ad *AdobeDownloader) GetAccessToken() error {
	// Implementation for OAuth token retrieval
	// See: https://developer.adobe.com/document-services/docs/overview/pdf-services-api/
	return nil
}

// DownloadDocument downloads a PDF from Adobe Document Cloud
func (ad *AdobeDownloader) DownloadDocument(shareLink, outputPath string) error {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	req, err := http.NewRequest("GET", shareLink, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers to mimic browser behavior
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")
	req.Header.Set("Accept", "application/pdf,*/*")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create output directory
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Copy response body to file
	bytesWritten, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✅ Downloaded %d bytes to %s\n", bytesWritten, outputPath)
	return nil
}

// DownloadBatch downloads multiple documents from a list of share links
func (ad *AdobeDownloader) DownloadBatch(links []string, outputDir string) error {
	for i, link := range links {
		outputPath := filepath.Join(outputDir, fmt.Sprintf("case_%02d.pdf", i+1))

		fmt.Printf("📥 Downloading %d/%d: %s\n", i+1, len(links), link)

		if err := ad.DownloadDocument(link, outputPath); err != nil {
			fmt.Printf("⚠️  Failed to download %s: %v\n", link, err)
			continue
		}

		// Be nice to Adobe's servers
		time.Sleep(2 * time.Second)
	}

	return nil
}
