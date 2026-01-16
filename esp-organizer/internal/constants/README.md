# Constants Package - Universal Revision Tool Framework

This package defines all hardcoded, subject-independent constants for the revision tool. It's designed to work across multiple academic subjects and levels.

## Overview

The revision tool serves **three different subjects** as POC:
- **GCSE** - Secondary education (age ~14-16)
- **Civil Engineering** - University undergraduate
- **Clinical Immunology** - Professional exam/postgraduate

All three use the **same underlying structure** with subject-specific extensions.

## Files

### 1. `academic.go` - Learning Framework
Defines educational levels and difficulty scales that apply universally.

**Key Types:**
- `AcademicLevel` - GCSE, A-Level, Undergraduate, Postgraduate, Professional
- `Difficulty` - Easy, Medium, Hard
- `EstimatedTimeMinutes` - Standard time allocations for study

**Usage:**
```go
import "esp-organizer/internal/constants"

// Validate academic level
if constants.IsValidAcademicLevel("undergraduate") {
    // Valid level
}

// Get all levels
for _, level := range constants.AllAcademicLevels() {
    fmt.Println(level)
}
```

### 2. `content.go` - Content Structure
Defines what types of learning materials exist and how content is sourced.

**Key Types:**
- `ContentType` - Chapter, Case Study, Definition, Practice Question, etc.
- `SourceType` - Book, PDF, Upload, Exam Board, etc.

**Content Type Examples:**
- `chapter` - "Chapter 1: X-linked Agammaglobulinemia" (immunology)
- `chapter` - "Chapter 3: Structural Design" (engineering)
- `practice_question` - GCSE exam question
- `case_study` - Patient case (immunology) or design example (engineering)
- `definition` - "B-cell deficiency" or "Young's modulus"

**Usage:**
```go
// Validate content type
if constants.IsValidContentType("case_study") {
    // Store with this type
}

// Validate source
if constants.IsValidSourceType("book") {
    // This is a published book
}
```

### 3. `tags.go` - Categorization System
Universal tags that work across all subjects, plus subject-specific examples.

**Universal Tags:**
- Study: `revision`, `exam_prep`, `practice_question`, `key_concept`
- Type: `case_study`, `worked_example`, `formula`, `diagram`
- Level: `foundation`, `core_knowledge`, `advanced`
- Assessment: `past_exam_question`, `mock_exam`

**Subject-Specific Examples:**
The function `SubjectSpecificTagExamples()` shows tags typically used per subject, but these should come from database configuration.

**Usage:**
```go
// Check if tag is valid
if constants.IsValidUniversalTag("exam_prep") {
    // OK to use this tag
}

// See subject-specific examples
examples := constants.SubjectSpecificTagExamples()
immunologyTags := examples["immunology"]
```

### 4. `metadata.go` - Data Structure
Standardized metadata that all content should have, with optional subject-specific extensions.

**Core RevisionMetadata Fields:**
```go
type RevisionMetadata struct {
    // Required universal fields
    Difficulty      string   // easy, medium, hard
    AcademicLevel   string   // gcse, undergraduate, etc.
    EstimatedTimeMs int      // How long to study (milliseconds)
    ContentType     string   // chapter, case_study, etc.
    
    // Learning connections
    RelatedConcepts []string // Linked content
    Prerequisites   []string // Must learn first
    Keywords        []string // For searching
    
    // Revision tracking
    LastReviewedAt  string   // When last studied
    ReviewCount     int      // How many times reviewed
    StudentNotes    string   // Student's annotations
    
    // Quality
    IsVerified      bool     // Checked by educator?
    QualityScore    float64  // 0.0-1.0 rating
    
    // Subject-specific
    Custom          map[string]interface{} // Per-subject extensions
}
```

**Subject-Specific Extensions:**

Clinical Immunology additional fields:
- `case_number`, `disease_category`, `clinical_findings`, `diagnosis`, `treatment`

Civil Engineering additional fields:
- `calculation_type`, `material_type`, `design_code`, `failure_modes`, `safety_factor`

GCSE additional fields:
- `exam_board`, `spec_section`, `command_word`, `mark_allocation`, `past_paper_ref`

**Usage:**
```go
// Create metadata for a new content item
metadata := constants.MetadataTemplate("medium", "undergraduate", "case_study", 1800000)

// Add subject-specific extension
metadata.Custom["material_type"] = "steel"
metadata.Custom["safety_factor"] = 1.5

// Validate metadata
validation := constants.DefaultMetadataValidation()
if metadata.QualityScore < validation.MinQualityScore {
    // Quality too low
}
```

## Design Principles

### ✅ Subject-Independent (Hardcoded)
- Academic levels (GCSE, A-Level, etc.)
- Difficulty scales (Easy, Medium, Hard)
- Content types (Chapter, Case Study, etc.)
- Source types (Book, PDF, etc.)
- Universal tags (revision, exam_prep, etc.)
- Core metadata structure

### ❌ Subject-Dependent (From Database/Config)
- Subject/domain names themselves
- Subject-specific tags
- Subject-specific books and sources
- Subject-specific metadata extension values

## Adding a New Subject

To add a new subject (e.g., "History"):

1. **No changes needed to constants** - The framework already supports it
2. **Add to database/config:**
   - Create subject entry with domain name, academic levels
   - Define subject-specific tags
   - Add known sources (books, exam boards, etc.)
3. **Extend metadata (optional):**
   - Add custom fields to `SubjectMetadataExtensions` if needed
   - Example for History: `"historical_period"`, `"primary_source"`, `"era"`

## Validation Functions

All constant types have validation helpers:

```go
constants.IsValidAcademicLevel(level string) bool
constants.IsValidDifficulty(difficulty string) bool
constants.IsValidContentType(contentType string) bool
constants.IsValidSourceType(sourceType string) bool
constants.IsValidUniversalTag(tag string) bool
```

Use these before storing data to prevent invalid values.

## Integration Points

These constants are used in:
- **Upload handlers** - Validate source type, content type
- **Content model** - Default values for metadata fields
- **Search/filter** - Valid values for queries
- **API validation** - Accept only valid enum values
- **Frontend** - Dropdown lists, validation

## Future Extensions

The design supports:
- Adding new academic levels without code change (add to database)
- Subject-specific validation rules
- Localization (translate constants)
- Difficulty scaling per subject (easier GCSE vs harder postgrad)
- Estimated time calibration per subject

---

**Last Updated:** 2025-12-22  
**Applies to:** All three POC subjects (GCSE, Civil Engineering, Clinical Immunology)
