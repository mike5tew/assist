package InfoIn

import (
	"fmt"
	"strings"
)

type MarkdownLink struct {
	SourceTerm     string
	TargetTerm     string
	RelationType   string
	IsCaseSpecific bool
	OriginalSource string // Preserves the * for logging/debugging
}

func ParseMarkdownLink(lines []string) (*MarkdownLink, error) {
	if len(lines) < 3 {
		return nil, fmt.Errorf("insufficient lines for a semantic link")
	}

	sourceLine := strings.TrimSpace(lines[0])
	relationType := strings.TrimSpace(lines[1])
	targetLine := strings.TrimSpace(lines[2])

	// Check for case-specific marker
	isCaseSpecific := strings.HasPrefix(sourceLine, "*")
	originalSource := sourceLine

	// Remove the * marker for storage
	if isCaseSpecific {
		sourceLine = strings.TrimPrefix(sourceLine, "*")
		sourceLine = strings.TrimSpace(sourceLine)
	}

	return &MarkdownLink{
		SourceTerm:     sourceLine,
		TargetTerm:     targetLine,
		RelationType:   relationType,
		IsCaseSpecific: isCaseSpecific,
		OriginalSource: originalSource, // Keep for debugging
	}, nil
}
