package extraction

import (
	"context"
	"crypto/rand"
	"esp-organizer/internal/domain/codeanalysis"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type CodeExtractionJob struct {
	ID            string             `json:"id"`
	Status        string             `json:"status"`
	FilePath      string             `json:"file_path"`
	Language      string             `json:"language"`
	ExtractedData *ExtractedCodeData `json:"extracted_data,omitempty"`
	Error         string             `json:"error,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
}

type ExtractedCodeData struct {
	Functions    []codeanalysis.CodeFunction `json:"functions"`
	Components   []codeanalysis.CodeFunction `json:"components"`
	Classes      []codeanalysis.CodeFunction `json:"classes"`
	Summary      string                      `json:"summary"`
	Dependencies []string                    `json:"dependencies"`
}

// Remove duplicate generateJobID function - it's already in textract.go
// func generateJobID() string {
// 	bytes := make([]byte, 16)
// 	rand.Read(bytes)
// 	return fmt.Sprintf("%x", bytes)
// }

// Use a different function name for code extraction
func generateCodeJobID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("code-%x", bytes)
}

func ProcessCodeFile(ctx context.Context, filePath string) (*CodeExtractionJob, error) {
	job := &CodeExtractionJob{
		ID:        generateCodeJobID(), // Use the renamed function
		Status:    "processing",
		FilePath:  filePath,
		Language:  detectLanguage(filePath),
		CreatedAt: time.Now(),
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		job.Status = "failed"
		job.Error = err.Error()
		return job, err
	}

	// Analyze code
	analyzer := codeanalysis.NewCodeAnalyzer()
	var functions []codeanalysis.CodeFunction

	switch job.Language {
	case "react", "javascript", "typescript":
		functions, err = analyzer.AnalyzeReactFile(filePath, string(content))
	case "php":
		functions, err = analyzer.AnalyzePHPFile(filePath, string(content))
	default:
		err = fmt.Errorf("unsupported language: %s", job.Language)
	}

	if err != nil {
		job.Status = "failed"
		job.Error = err.Error()
		return job, err
	}

	// Generate summary
	summary := generateCodeSummary(functions, job.Language)

	job.ExtractedData = &ExtractedCodeData{
		Functions:    functions,
		Summary:      summary,
		Dependencies: extractDependencies(string(content), job.Language),
	}

	job.Status = "completed"

	// Database storage removed - will be handled by handlers
	return job, nil
}

// Remove these functions entirely:
// - storeCodeInDatabases
// - storeCodeInWeaviate
// - buildSemanticDescription

// Helper functions
func detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".js":
		return "javascript"
	case ".jsx":
		return "react"
	case ".ts":
		return "typescript"
	case ".tsx":
		return "typescript-react"
	case ".php":
		return "php"
	case ".py":
		return "python"
	case ".go":
		return "go"
	default:
		return "unknown"
	}
}

func generateCodeSummary(functions []codeanalysis.CodeFunction, language string) string {
	if len(functions) == 0 {
		return fmt.Sprintf("No %s functions or components found", language)
	}

	summary := fmt.Sprintf("Found %d %s ", len(functions), language)

	switch language {
	case "react", "javascript", "typescript":
		summary += "components and functions"
	case "php":
		summary += "functions"
	default:
		summary += "code elements"
	}

	// Add function names
	if len(functions) <= 3 {
		names := make([]string, len(functions))
		for i, fn := range functions {
			names[i] = fn.Name
		}
		summary += ": " + strings.Join(names, ", ")
	} else {
		names := make([]string, 3)
		for i := 0; i < 3; i++ {
			names[i] = functions[i].Name
		}
		summary += fmt.Sprintf(": %s and %d more", strings.Join(names, ", "), len(functions)-3)
	}

	return summary
}

func extractDependencies(content string, language string) []string {
	dependencies := []string{}

	switch language {
	case "react", "javascript", "typescript":
		// Extract import statements
		importPattern := regexp.MustCompile(`import\s+.*?\s+from\s+['"]([^'"]+)['"]`)
		matches := importPattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				dependencies = append(dependencies, match[1])
			}
		}

	case "php":
		// Extract require/include statements
		requirePattern := regexp.MustCompile(`(?:require|include)(?:_once)?\s*\(?['"]([^'"]+)['"]`)
		matches := requirePattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				dependencies = append(dependencies, match[1])
			}
		}

		// Extract use statements
		usePattern := regexp.MustCompile(`use\s+([^;]+);`)
		useMatches := usePattern.FindAllStringSubmatch(content, -1)
		for _, match := range useMatches {
			if len(match) > 1 {
				dependencies = append(dependencies, match[1])
			}
		}

	case "python":
		// Extract import statements
		importPattern := regexp.MustCompile(`(?:from\s+(\S+)\s+)?import\s+([^#\n]+)`)
		matches := importPattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 && match[1] != "" {
				dependencies = append(dependencies, match[1])
			} else if len(match) > 2 {
				deps := strings.Split(match[2], ",")
				for _, dep := range deps {
					dependencies = append(dependencies, strings.TrimSpace(dep))
				}
			}
		}
	}

	return dependencies
}
