package api

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// extractTextWithPDFCPU extracts text from a PDF using pdftotext (poppler-utils).
// The container has poppler-utils installed, making this the reliable extraction path.
func extractTextWithPDFCPU(pdfBytes []byte) (string, error) {
	if len(pdfBytes) == 0 {
		return "", fmt.Errorf("PDF is empty")
	}

	// Write to a temp file — pdftotext requires a file path
	tmpFile, err := os.CreateTemp("", "ntm-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(pdfBytes); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write PDF to temp file: %w", err)
	}
	tmpFile.Close()

	// -layout preserves whitespace structure; "-" outputs to stdout
	out, err := exec.Command("pdftotext", "-layout", tmpFile.Name(), "-").Output()
	if err != nil {
		return "", fmt.Errorf("pdftotext failed: %w", err)
	}

	text := strings.TrimSpace(string(out))
	if len(text) == 0 {
		return "", fmt.Errorf("pdftotext produced no text — PDF may be scanned/image-only")
	}

	log.Printf("✅ Successfully extracted %d characters from PDF via pdftotext", len(text))
	return text, nil
}
