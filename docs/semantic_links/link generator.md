I don't have the exact description from your ESP project, but based on our conversation, I can reconstruct what the **Semantic Link Extraction Tool** should do and provide a complete specification.

---

# **ESP Semantic Link Extraction Tool**
## **Purpose**
A 3-click interface that allows users (teachers/researchers/curators) to quickly extract meaningful semantic relationships from educational texts and store them in a structured knowledge graph with quality metadata.

## **Core Philosophy**
- **3 clicks max** per semantic link
- **Quality over quantity** - reject vague relationships
- **Automatic provenance** - capture source context
- **Hierarchy-aware** - infer positioning in knowledge graph

---

## **Technical Specification**

### **Input Sources**
1. **PDF documents** (uploaded via OCR)
2. **Existing lesson_files** (from your legacy database)
3. **Plain text** (copy-paste interface)
4. **Web articles** (URL import)

### **Output Format**
```javascript
{
  "id": "sl_abc123def456",
  "source_term": "ribosomes",
  "relation_type": "location_of",
  "target_term": "protein synthesis",
  
  // Provenance (auto-captured)
  "document_id": "gcse_biology_2023",
  "page_number": 45,
  "sentence_id": "sentence_456",
  "text_context": "Ribosomes, where protein synthesis takes place...",
  "position_context": {
    "source_start": 0,
    "source_end": 9,
    "target_start": 18,
    "target_end": 35
  },
  
  // Quality metadata
  "confidence": 0.95,
  "hierarchy_level": 2,  // 1=root, 2=intermediate, 3=specific
  "quality_flags": [],  // e.g., "overgeneralization", "speculative"
  "extraction_method": "manual",  // or "llm_auto", "llm_suggested"
  "reviewed_by": "user_123",
  "created_at": "2025-01-06T10:30:00Z"
}
```

---

## **User Interface Flow**

### **Step 1: Document Selection & Viewing**
```
[PDF Viewer with text overlay]
_______________________________________
| 📄 Chapter 3: Protein Synthesis     |
|-------------------------------------|
| Ribosomes, where protein synthesis  |
| takes place. All the proteins needed|
| in the cell are made here.          |
|                                     |
| [Previous]       [Next]             |
---------------------------------------
```

### **Step 2: Term Selection (Click 1 & 2)**
User clicks on two terms in the text:
1. **Click 1**: "ribosomes" → becomes source term
2. **Click 2**: "protein synthesis" → becomes target term

**Visual feedback**: Terms highlighted in different colors.

### **Step 3: Relationship Selection (Click 3)**
```
Relationship type dropdown appears near selection:

[location_of] ← Default based on context
▼
  ── Biological Processes ──
  [catalyzes]              [produces]
  [regulates]              [inhibits]
  [transports]             [synthesizes]
  
  ── Spatial/Structural ──
  [contains]               [part_of]
  [attached_to]            [surrounds]
  
  ── Causal/Temporal ──
  [causes]                 [precedes]
  [enables]                [requires]
  
  ── Abstract ──
  [is_example_of]          [is_type_of]
  [contrasts_with]         [similar_to]
  
  ❌ [too_vague] ← Reject relationship
```

### **Step 4: Quality Quick-Set (Optional 4th click)**
If user notices issues:
```
[✅ High confidence]  [⚠️ Needs verification]  [❌ Problematic]
```

---

## **Backend Processing Pipeline**

### **1. Pre-processing (Automatic)**
```go
func preprocessDocument(doc Document) PreprocessedDoc {
    return PreprocessedDoc{
        Text:           extractText(doc),
        Sentences:      segmentSentences(extractText(doc)),
        TermPositions:  findTermPositions(extractText(doc)),
        DocumentMeta:   extractMetadata(doc),
        OCRConfidence:  calculateOCRConfidence(doc),
    }
}
```

### **2. Semantic Link Creation**
```go
func createSemanticLink(selection Selection, user User) SemanticLink {
    link := SemanticLink{
        ID:           generateUUID(),
        SourceTerm:   cleanTerm(selection.SourceTerm),
        RelationType: validateRelation(selection.RelationType),
        TargetTerm:   cleanTerm(selection.TargetTerm),
        
        // Auto-captured provenance
        DocumentID:   selection.DocumentID,
        PageNumber:  selection.PageNumber,
        SentenceID:  generateSentenceID(selection.TextContext),
        TextContext: selection.TextContext,
        Position:    selection.Position,
        
        // Quality scoring
        Confidence:    calculateInitialConfidence(selection, user),
        HierarchyLevel: inferHierarchyLevel(selection),
        ExtractionMethod: "manual",
        
        // User info
        CreatedBy:    user.ID,
        CreatedAt:    time.Now(),
        QualityFlags: []string{},
    }
    
    // Apply any quality overrides
    if selection.QualityFlag != "" {
        link.QualityFlags = append(link.QualityFlags, selection.QualityFlag)
        link.Confidence = adjustConfidence(link.Confidence, selection.QualityFlag)
    }
    
    return link
}
```

### **3. Hierarchy Inference**
```go
func inferHierarchyLevel(selection Selection) int {
    // Rule-based hierarchy inference
    terms := []string{selection.SourceTerm, selection.TargetTerm}
    
    // Level 1: Broad concepts (systems, categories)
    broadPatterns := []string{"system", "category", "class", "type", "field"}
    if containsAny(terms, broadPatterns) {
        return 1
    }
    
    // Level 2: Specific entities (diseases, processes, organs)
    specificPatterns := []string{"syndrome", "disease", "process", "organ", "protein"}
    if containsAny(terms, specificPatterns) {
        return 2
    }
    
    // Level 3: Detailed components (genes, molecules, specific parts)
    detailedPatterns := []string{"gene", "molecule", "enzyme", "receptor", "subunit"}
    if containsAny(terms, detailedPatterns) {
        return 3
    }
    
    // Default based on term specificity
    return calculateTermSpecificity(terms)
}
```

### **4. Storage**
```go
func storeSemanticLink(link SemanticLink) error {
    // Store in Weaviate (vector + metadata)
    weaviateClient.DataObject().Creator().
        WithClassName("SemanticLink").
        WithProperties(map[string]interface{}{
            "source":      link.SourceTerm,
            "relation":    link.RelationType,
            "target":      link.TargetTerm,
            "confidence":  link.Confidence,
            "hierarchy":   link.HierarchyLevel,
            "context":     link.TextContext,
            "provenance": map[string]interface{}{
                "document": link.DocumentID,
                "page":     link.PageNumber,
                "sentence": link.SentenceID,
            },
        }).
        Do(context.Background())
    
    // Store in relational DB for queries
    db.Exec(`
        INSERT INTO semantic_links 
        (id, source, relation, target, confidence, hierarchy_level, ...)
        VALUES (?, ?, ?, ?, ?, ?, ...)`,
        link.ID, link.SourceTerm, link.RelationType, 
        link.TargetTerm, link.Confidence, link.HierarchyLevel, ...)
    
    return nil
}
```

---

## **Key Features**

### **1. Smart Defaults**
- **Relation type suggestion**: Based on domain (biology, physics, chemistry)
- **Hierarchy auto-detection**: "cell → contains → nucleus" = level 2
- **Confidence scoring**: Manual extraction = 0.95, LLM-suggested = 0.65

### **2. Quality Controls**
```javascript
// Vague relationship rejection
const VAGUE_RELATIONS = [
    "related_to", "associated_with", "involved_in",
    "plays_role_in", "connected_to", "affects"  // too vague
];

// Overgeneralization detection
if (text.includes("all ") || text.includes("every ") || text.includes("always ")) {
    suggestQualityFlag("overgeneralization");
}

// Speculative language detection
if (text.includes("may ") || text.includes("could ") || text.includes("possibly ")) {
    adjustConfidence(0.7);  // Lower confidence for speculative statements
}
```

### **3. Batch Processing Mode**
```javascript
// For textbook chapters with repeated patterns
const batchTemplate = {
    pattern: "[PROTEIN] [ACTION] [MOLECULE]",
    examples: [
        "Insulin regulates blood glucose",
        "Hemoglobin transports oxygen"
    ],
    autoMap: {
        "[PROTEIN]": "source",
        "[ACTION]": "relation", 
        "[MOLECULE]": "target"
    }
};

// Process 50 similar sentences in 2 minutes vs 25 minutes manually
```

### **4. Duplicate & Conflict Detection**
```go
func checkForDuplicates(newLink SemanticLink) ([]SimilarLink, bool) {
    // Check for identical triples
    identical := findIdenticalTriple(newLink)
    if identical != nil {
        return []SimilarLink{identical}, true
    }
    
    // Check for similar meaning
    similar := findSimilarMeaning(newLink, threshold=0.8)
    
    // Check for contradictions
    contradictions := findContradictions(newLink)
    
    return append(similar, contradictions...), len(similar) > 0
}
```

---

## **Integration Points**

### **With Legacy System**
```go
// Link to existing lesson_files
func linkToLessonFile(semanticLinkID string, lessonFileID int64) {
    db.Exec(`
        INSERT INTO lesson_file_semantic_links 
        (semantic_link_id, lesson_files_id)
        VALUES (?, ?)`,
        semanticLinkID, lessonFileID)
}

// Connect to skills
func connectToSkill(semanticLinkID string, skillID int64) {
    db.Exec(`
        INSERT INTO skill_semantic_links
        (semantic_link_id, skills_key_id)
        VALUES (?, ?)`,
        semanticLinkID, skillID)
}
```

### **With ESP Knowledge Graph**
```go
// Build hierarchical relationships
func buildHierarchicalRelationships(link SemanticLink) {
    // Find parent concepts
    parents := findParentConcepts(link.SourceTerm)
    for _, parent := range parents {
        createHierarchyLink(parent, link.SourceTerm, "generalizes")
    }
    
    // Find child concepts  
    children := findChildConcepts(link.TargetTerm)
    for _, child := range children {
        createHierarchyLink(link.TargetTerm, child, "specializes")
    }
}
```

---

## **API Endpoints**

### **POST /api/semantic-links/extract**
```json
{
    "document_id": "textbook_chapter_3",
    "selections": [
        {
            "source_term": "ribosomes",
            "target_term": "protein synthesis", 
            "relation_type": "location_of",
            "text_context": "Ribosomes, where protein synthesis takes place...",
            "position": {
                "source_start": 0,
                "source_end": 9,
                "target_start": 18,
                "target_end": 35
            },
            "quality_flag": "high_confidence"
        }
    ]
}
```

### **GET /api/semantic-links/search?term=ribosomes**
Returns related semantic links with hierarchy context.

### **POST /api/semantic-links/batch**
For processing multiple similar extractions.

---

## **Implementation Priority (Week 1)**

### **Day 1-2: Basic Interface**
- PDF text selection handler
- 3-click workflow
- Simple relationship dropdown
- Local storage (IndexedDB)

### **Day 3-4: Backend Integration**
- Go API endpoints
- Weaviate schema setup
- Provenance capture
- Basic quality scoring

### **Day 5: Testing & Polish**
- Test with GCSE Biology chapter
- Performance optimization
- Error handling
- Export functionality

---

## **Success Metrics**
- **Extraction speed**: < 10 seconds per semantic link
- **Accuracy**: 95%+ of extractions are meaningful (manual review)
- **User satisfaction**: Teachers can process 50+ links/hour
- **Data quality**: < 5% vague relationships in final dataset

---

## **The Big Picture**

This tool isn't just extracting relationships—it's **building the foundational layer of your ESP knowledge graph**. Every click creates:
1. A **semantic link** for querying
2. **Quality metadata** for confidence scoring  
3. **Provenance data** for traceability
4. **Hierarchical positioning** for context
5. **Training data** for future LLM improvement

**Start simple**: Get the 3-click workflow working with GCSE Biology. Everything else builds from there.