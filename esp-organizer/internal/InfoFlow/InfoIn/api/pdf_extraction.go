package api

import (
	"bytes"
	"fmt"
	"log"
)

// extractTextWithPDFCPU attempts to extract text from a searchable PDF
// This is a FREE, LOCAL alternative to AWS Textract for PDFs with embedded text
// Note: Full text extraction from PDFs is complex. This function validates PDF structure
// and indicates whether AWS Textract processing is needed.
func extractTextWithPDFCPU(pdfBytes []byte) (string, error) {
	// For now, we'll use a simple heuristic approach:
	// If the PDF has text markers, it's likely searchable.
	// Full extraction would require a more sophisticated library.

	if len(pdfBytes) == 0 {
		return "", fmt.Errorf("PDF is empty")
	}

	// Check if PDF has embedded text streams (indicates searchable PDF)
	hasTextStreams := bytes.Contains(pdfBytes, []byte("stream")) &&
		bytes.Contains(pdfBytes, []byte("BT")) && // BT = Begin Text operator
		bytes.Contains(pdfBytes, []byte("ET")) // ET = End Text operator

	if !hasTextStreams {
		return "", fmt.Errorf("PDF appears to be scanned (no text streams found)")
	}

	// Extract basic text markers from PDF
	extractedText := extractTextMarkers(pdfBytes)

	if len(extractedText) == 0 {
		return "", fmt.Errorf("no extractable text found in PDF")
	}

	log.Printf("✅ Successfully extracted %d characters from searchable PDF", len(extractedText))
	return extractedText, nil
}

// extractTextMarkers attempts to extract text strings from PDF byte content
// This is a simple heuristic approach that works for text-heavy PDFs
func extractTextMarkers(pdfBytes []byte) string {
	var result bytes.Buffer
	inTextMode := false
	currentString := bytes.NewBuffer([]byte{})

	i := 0
	for i < len(pdfBytes) {
		// Look for text mode start
		if i+1 < len(pdfBytes) && pdfBytes[i] == 'B' && pdfBytes[i+1] == 'T' {
			inTextMode = true
			i += 2
			continue
		}

		// Look for text mode end
		if i+1 < len(pdfBytes) && pdfBytes[i] == 'E' && pdfBytes[i+1] == 'T' {
			if currentString.Len() > 0 {
				result.Write(currentString.Bytes())
				result.WriteString("\n")
				currentString.Reset()
			}
			inTextMode = false
			i += 2
			continue
		}

		// In text mode, extract printable characters
		if inTextMode {
			b := pdfBytes[i]
			// Keep printable ASCII and common extended ASCII
			if (b >= 32 && b <= 126) || b == '\n' || b == '\r' || b == '\t' {
				currentString.WriteByte(b)
			}
		}

		i++
	}

	return result.String()
}

// For a production-quality solution, use one of these specialized libraries:
// - github.com/ledongthuc/pdf (pure Go, basic text extraction)
// - github.com/mandykoh/prism (more advanced extraction)
// - Or continue using AWS Textract for reliable OCR and structured data extraction
//
// The pdfcpu library is primarily designed for PDF manipulation (splitting, merging)
// rather than text extraction, so AWS Textract remains the best choice for
// reliable text extraction from both searchable and scanned PDFs.
