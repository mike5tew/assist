package infoin

import (
	"regexp"
	"strings"
	"unicode"
)

// SymbolNormalizer handles conversion of special symbols to searchable text
type SymbolNormalizer struct {
	symbolMap    map[string]string
	greekLetters map[string]string
	unitMap      map[string]string
}

// NewSymbolNormalizer creates a new symbol normalizer with comprehensive mappings
func NewSymbolNormalizer() *SymbolNormalizer {
	return &SymbolNormalizer{
		// Greek letters (both uppercase and lowercase)
		greekLetters: map[string]string{
			"α": "alpha",
			"β": "beta",
			"γ": "gamma",
			"δ": "delta",
			"ε": "epsilon",
			"ζ": "zeta",
			"η": "eta",
			"θ": "theta",
			"ι": "iota",
			"κ": "kappa",
			"λ": "lambda",
			"μ": "mu",
			"ν": "nu",
			"ξ": "xi",
			"ο": "omicron",
			"π": "pi",
			"ρ": "rho",
			"σ": "sigma",
			"τ": "tau",
			"υ": "upsilon",
			"φ": "phi",
			"χ": "chi",
			"ψ": "psi",
			"ω": "omega",
			"Α": "Alpha",
			"Β": "Beta",
			"Γ": "Gamma",
			"Δ": "Delta",
			"Ε": "Epsilon",
			"Ζ": "Zeta",
			"Η": "Eta",
			"Θ": "Theta",
			"Ι": "Iota",
			"Κ": "Kappa",
			"Λ": "Lambda",
			"Μ": "Mu",
			"Ν": "Nu",
			"Ξ": "Xi",
			"Ο": "Omicron",
			"Π": "Pi",
			"Ρ": "Rho",
			"Σ": "Sigma",
			"Τ": "Tau",
			"Υ": "Upsilon",
			"Φ": "Phi",
			"Χ": "Chi",
			"Ψ": "Psi",
			"Ω": "Omega",
		},

		// Common scientific symbols
		symbolMap: map[string]string{
			"±": "plus-or-minus",
			"≤": "less-than-or-equal-to",
			"≥": "greater-than-or-equal-to",
			"≠": "not-equal-to",
			"≈": "approximately",
			"×": "times",
			"÷": "divided-by",
			"→": "to",
			"←": "from",
			"↔": "reversible",
			"°": "degrees",
			"'": "prime",
			"″": "double-prime",
			"⁰": "0", "¹": "1", "²": "2", "³": "3", "⁴": "4",
			"⁵": "5", "⁶": "6", "⁷": "7", "⁸": "8", "⁹": "9",
			"₀": "0", "₁": "1", "₂": "2", "₃": "3", "₄": "4",
			"₅": "5", "₆": "6", "₇": "7", "₈": "8", "₉": "9",
		},

		// Medical/scientific units with special symbols
		unitMap: map[string]string{
			"µl": "microliter",
			"μl": "microliter",
			"µL": "microliter",
			"μL": "microliter",
			"µg": "microgram",
			"μg": "microgram",
			"µm": "micrometer",
			"μm": "micrometer",
			"µM": "micromolar",
			"μM": "micromolar",
			"mg": "milligram",
			"ml": "milliliter",
			"mL": "milliliter",
			"dl": "deciliter",
			"dL": "deciliter",
			"ng": "nanogram",
			"pg": "picogram",
			"kg": "kilogram",
			"°C": "degrees-Celsius",
			"°F": "degrees-Fahrenheit",
		},
	}
}

// NormalizeText performs comprehensive text normalization
func (sn *SymbolNormalizer) NormalizeText(text string) string {
	// Step 1: Replace Greek letters
	for symbol, name := range sn.greekLetters {
		text = strings.ReplaceAll(text, symbol, name)
	}

	// Step 2: Replace units (order matters - do compound units first)
	for unit, name := range sn.unitMap {
		text = strings.ReplaceAll(text, unit, name)
	}

	// Step 3: Replace other symbols
	for symbol, replacement := range sn.symbolMap {
		text = strings.ReplaceAll(text, symbol, replacement)
	}

	// Step 4: Handle superscripts in scientific notation (e.g., "10⁹" → "10^9")
	text = sn.normalizeSuperscripts(text)

	// Step 5: Handle subscripts (e.g., "H₂O" → "H2O")
	text = sn.normalizeSubscripts(text)

	// Step 6: Remove multiple spaces
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return strings.TrimSpace(text)
}

// NormalizeForSearch creates a search-friendly version (more aggressive)
func (sn *SymbolNormalizer) NormalizeForSearch(text string) string {
	// Start with standard normalization
	normalized := sn.NormalizeText(text)

	// Additional search-specific transformations
	normalized = strings.ToLower(normalized)

	// Remove punctuation except hyphens (important for terms like "T-cell")
	normalized = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(normalized, " ")

	// Normalize whitespace
	normalized = regexp.MustCompile(`\s+`).ReplaceAllString(normalized, " ")

	return strings.TrimSpace(normalized)
}

// NormalizeForVectorization preserves more structure for embeddings
func (sn *SymbolNormalizer) NormalizeForVectorization(text string) string {
	// Light normalization - preserve semantic meaning
	normalized := sn.NormalizeText(text)

	// Keep punctuation that provides context (periods, commas)
	// Remove only problematic characters
	normalized = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return ' '
	}, normalized)

	return strings.TrimSpace(normalized)
}

// normalizeSuperscripts converts Unicode superscripts to ^notation
func (sn *SymbolNormalizer) normalizeSuperscripts(text string) string {
	superscriptPattern := regexp.MustCompile(`([0-9]+)([⁰¹²³⁴⁵⁶⁷⁸⁹]+)`)

	return superscriptPattern.ReplaceAllStringFunc(text, func(match string) string {
		parts := superscriptPattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		base := parts[1]
		exponent := parts[2]

		// Convert superscript digits to normal digits
		normalExponent := ""
		for _, char := range exponent {
			if replacement, ok := sn.symbolMap[string(char)]; ok {
				normalExponent += replacement
			}
		}

		return base + "^" + normalExponent
	})
}

// normalizeSubscripts converts Unicode subscripts to normal notation
func (sn *SymbolNormalizer) normalizeSubscripts(text string) string {
	subscriptPattern := regexp.MustCompile(`([A-Za-z]+)([₀₁₂₃₄₅₆₇₈₉]+)`)

	return subscriptPattern.ReplaceAllStringFunc(text, func(match string) string {
		parts := subscriptPattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		base := parts[1]
		subscript := parts[2]

		// Convert subscript digits to normal digits
		normalSubscript := ""
		for _, char := range subscript {
			if replacement, ok := sn.symbolMap[string(char)]; ok {
				normalSubscript += replacement
			}
		}

		return base + normalSubscript
	})
}

// CreateSearchVariants creates multiple search-friendly variants of a term
func (sn *SymbolNormalizer) CreateSearchVariants(term string) []string {
	variants := make(map[string]bool)

	// Original term
	variants[term] = true

	// Normalized version
	variants[sn.NormalizeForSearch(term)] = true

	// Handle common medical abbreviations
	if strings.Contains(term, "µ") || strings.Contains(term, "μ") {
		// Add "micro" variant
		microVariant := strings.ReplaceAll(term, "µ", "micro")
		microVariant = strings.ReplaceAll(microVariant, "μ", "micro")
		variants[sn.NormalizeForSearch(microVariant)] = true
	}

	// Convert map to slice
	result := make([]string, 0, len(variants))
	for variant := range variants {
		if variant != "" {
			result = append(result, variant)
		}
	}

	return result
}
