# CHISG Project - AI Prompt Guide

**Objective**: To provide clear instructions for an AI assistant working on the CHISG (Contextual Human-like Intelligence Scaffolding Graph) project.

---

## 1. Project Overview

**What is CHISG?**
CHISG is a Go-based microservice that acts as a semantic knowledge graph. Its primary function is to take a query about a topic and return structured, easy-to-understand knowledge. It is the "brain" that provides factual and conceptual information to other services like the HumanOS AI Tutor.

**Core Purpose**:
- Deconstruct complex topics into simple, related concepts.
- Provide analogies and prerequisites to aid learning.
- Serve this structured knowledge via a low-latency REST API.
- **🆕 Extract and manage semantic links** from unstructured content using sample-driven LLM transformation.

## 2. Core Architecture

- **Language**: Go (Golang)
- **Vector Database**: Weaviate (for semantic search and graph relationships)
- **Metadata Store**: MongoDB (for storing raw documents, user data, extraction samples)
- **API**: RESTful API built with Gin framework
- **🆕 LLM Integration**: AWS Bedrock (Claude/Titan) for semantic extraction and embedding generation
- **🆕 Document Processing**: AWS Textract for scanned PDFs, local extraction for searchable PDFs

### Document Processing Pipeline

1. **Upload** → `ImmunologyChapterUploadHandler`
   - Accepts PDF files
   - Detects if PDF is searchable (has embedded text)
   - Routes to appropriate processor

2. **Text Extraction**
   - **Searchable PDFs**: Local extraction (free, fast)
   - **Scanned PDFs**: AWS Textract (OCR, more expensive)

3. **Content Parsing** → `ProcessExtractedText`
   - Extracts case studies
   - Identifies medical terms
   - Structures chapter content

4. **Semantic Link Extraction** → `SemanticLinkService`
   - Uses sample-driven LLM prompts
   - Creates relationships between concepts
   - Stores in MongoDB and Weaviate

## 3. API Contract (The Rules)

### Knowledge Query Endpoint

The primary endpoint is `/v1/knowledge/query`.

**Request (Input):**
A JSON object matching the `KnowledgeQuery` struct:
```go
type KnowledgeQuery struct {
	Topic          string   `json:"topic"`
	Concepts       []string `json:"concepts"`
	TargetAudience string   `json:"target_audience"` // e.g., "gcse_student", "medical_student"
}
```

**Response (Output):**
A JSON object matching the `KnowledgeResponse` struct:
```go
type KnowledgeResponse struct {
	Summary         string            `json:"summary"`
	KeyConcepts     []string          `json:"key_concepts"`
	Prerequisites   []string          `json:"prerequisites"`
	Analogies       []string          `json:"analogies"`
	ConfidenceScore float64           `json:"confidence_score"`
	RelatedTopics   map[string]string `json:"related_topics"`
}
```

### 🆕 Semantic Link Extraction Endpoints

#### 1. Document Processing
`POST /api/immunology/upload`
- Uploads and processes documents
- Extracts semantic links using sample-driven prompts
- Returns batch ID for status tracking

#### 2. Sample Management
`GET /api/samples` - List extraction samples
`POST /api/samples` - Create new sample
`PUT /api/samples/:id` - Update sample
`DELETE /api/samples/:id` - Delete sample
`POST /api/samples/curate` - Trigger automatic curation

#### 3. HSG Query
`POST /api/immunology/hsg-search`
```json
{
  "query": "What causes XLA?",
  "domain": "immunology",
  "max_results": 5
}
```

## 4. Development Patterns & Expectations

### Core Principles

1.  **Stay within the API Contract**: All changes must respect the input and output structs defined above.
2.  **Weaviate is for "What"**: Use Weaviate for semantic queries, finding related concepts, and answering "what is X?" questions.
3.  **MongoDB is for "Who" and "Where"**: Use MongoDB to store document metadata, user profiles, extraction samples, and other structured data that doesn't belong in the vector graph.
4.  **Low Latency is Critical**: All operations should be optimized for speed. The target response time is **< 200ms**.
5.  **Think in Graphs**: When adding new knowledge, consider the prerequisites (what comes before?) and related topics (what comes next or sideways?).
6.  **Explain Your Work**: Clearly describe the changes you are making and why, referencing the architecture.

### 🆕 Semantic Link Extraction Pattern

#### Sample-Driven Transformation

The system uses a **sample-driven transformation agent** pattern for extracting semantic relationships:

1. **Sample Collection** (`extraction_samples` MongoDB collection)
   - Stores high-quality input/output examples
   - Domain-specific (immunology, personal, etc.)
   - Quality-scored (0.0-1.0)

2. **Extraction Process**
   ```go
   // Load domain-specific samples
   samples := loadSamplesByDomain("immunology")
   
   // Build prompt with examples
   prompt := buildExtractionPromptWithSamples(content, samples)
   
   // LLM extracts following the pattern
   response := llmClient.Generate(ctx, prompt)
   ```

3. **Quality Assessment**
   - Automatic validation against curated examples
   - Self-correction loop if quality < threshold
   - High-quality extractions (>0.9 confidence) become new samples

4. **Automatic Curation**
   - Scheduler runs every 24 hours
   - Harvests recent high-quality extractions
   - Adds to sample collection for continuous improvement

#### Key Data Structures

```go
// Document extraction tracking
type ExtractionJob struct {
    ID            string         `json:"id" bson:"_id"`
    Status        string         `json:"status" bson:"status"`
    Progress      float64        `json:"progress" bson:"progress"`
    ExtractedData *ExtractedData `json:"extracted_data,omitempty" bson:"extracted_data,omitempty"`
}

type ExtractedData struct {
    RawText        string        `json:"raw_text" bson:"raw_text"`
    ChapterContent string        `json:"chapter_content" bson:"chapter_content"`
    CaseStudies    []CaseStudy   `json:"case_studies" bson:"case_studies"`
    MedicalTerms   []MedicalTerm `json:"medical_terms" bson:"medical_terms"`
}

// Semantic relationship storage
type SemanticLink struct {
    ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    SourceTerm         string             `bson:"source_term" json:"source_term"`
    TargetTerm         string             `bson:"target_term" json:"target_term"`
    RelationType       string             `bson:"relation_type" json:"relation_type"`
    Context            string             `bson:"context" json:"context"`
    Confidence         float64            `bson:"confidence" json:"confidence"`
    SourceTermOriginal string             `bson:"source_term_original,omitempty" json:"source_term_original,omitempty"`
    TargetTermOriginal string             `bson:"target_term_original,omitempty" json:"target_term_original,omitempty"`
    SourceTermGenerality float64          `bson:"source_term_generality" json:"source_term_generality"`
    TargetTermGenerality float64          `bson:"target_term_generality" json:"target_term_generality"`
    SemanticDistance     float64          `bson:"semantic_distance" json:"semantic_distance"`
}

// Sample for few-shot LLM prompting
type ExtractionSample struct {
    InputText      string  `json:"input_text"`
    ExpectedOutput string  `json:"expected_output"`
    QualityScore   float64 `json:"quality_score"`
    Domain         string  `json:"domain"`
}
```

### Best Practices for Adding Features

1. **New Extraction Logic**: Always add corresponding samples to `extraction_samples` collection
2. **Domain Expansion**: Create domain-specific default samples in `getDefaultSamples()`
3. **Quality Control**: Use `assessExtractionQuality()` to validate new extraction logic
4. **Testing**: Compare against curated samples before deploying
5. **Monitoring**: Track quality metrics per domain via `/api/samples/quality/:domain`

### Error Handling

- **LLM Failures**: Log and continue with default samples
- **Validation Failures**: Trigger self-correction loop (max 2 attempts)
- **Database Errors**: Graceful degradation with in-memory fallbacks
- **Duplicate Links**: Skip silently using composite key deduplication

## 5. Scheduled Tasks

### Curation Scheduler

Runs every 24 hours to harvest high-quality extractions:

```go
// Start in main.go
scheduler, _ := scheduler.NewCurationScheduler(24)
scheduler.Start()
defer scheduler.Stop()
```

Manual trigger: `POST /api/samples/curate`

## How to Use This Guide

When you need help with the `chisg` project, copy and paste the contents of this file as the initial prompt to your AI assistant. This will ensure the assistant has the full context and follows the established patterns for this specific service.

---

## 🆕 Quick Reference: Adding a New Domain

1. Create sample in MongoDB:
   ```json
   {
     "domain": "neurology",
     "input_text": "Example input",
     "expected_output": "{...}",
     "quality_score": 0.95
   }
   ```

2. Add default samples in code:
   ```go
   func getDefaultSamples(domain string) []ExtractionSample {
       switch domain {
       case "neurology":
           return neurologySamples
       // ...
       }
   }
   ```

3. Test extraction quality:
   ```bash
   curl -X POST /api/samples/curate
   curl /api/samples/quality/neurology
   ```

4. Monitor and refine samples based on quality metrics.

## 🆕 Text Extraction Methods

### Method 1: Searchable PDF (Local)
- **Cost**: Free
- **Speed**: Fast (~seconds)
- **Use case**: PDFs with embedded text (e.g., from Adobe Scan)
- **Process**: Direct text extraction from PDF structure

### Method 2: AWS Textract (Cloud OCR)
- **Cost**: $1.50 per 1000 pages
- **Speed**: Slower (~minutes)
- **Use case**: Scanned documents without embedded text
- **Process**: OCR via AWS Textract → Text extraction

### Detection Logic
```go
func isSearchablePDF(pdfBytes []byte) bool {
    // Check for embedded fonts and text operators
    return bytes.Contains(pdfBytes, []byte("/Font"))
}
```
