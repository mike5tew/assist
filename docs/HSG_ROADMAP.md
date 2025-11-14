# HSG Implementation Roadmap

**Project**: ESP Organizer - Hierarchical Semantic Graph (HSG) RAG System  
**Start Date**: 2025-01-05  
**Target Completion**: Phase 2 by 2025-01-12

---

## Phase 1: Tier 1 Semantic Links ✅ COMPLETE

### Deliverables (All Complete)
- [x] PDF upload and OCR processing
- [x] Content storage in MongoDB
- [x] Semantic link extraction with LLM
- [x] External vectorization with AWS Titan
- [x] Semantic links stored in Weaviate
- [x] Basic vector search working

**Completion Date**: 2025-01-05  
**Status**: ✅ Production Ready

---

## Phase 2: Add Hierarchy to Tier 1 🔴 CURRENT PRIORITY

**Goal**: Enable hierarchical classification and traversal of semantic links

### 2.1: Hierarchy Classification (2-3 hours) ⏳ IN PROGRESS
**Target**: 2025-01-06 Morning

#### Tasks:
- [ ] **2.1.1**: Update `SemanticLink` model
  - File: `/internal/models/semantic_link.go`
  - Add fields: `HierarchyLevel int`, `IsParentOf []string`, `IsChildOf []string`
  - Duration: 15 min

- [ ] **2.1.2**: Create hierarchy classification logic
  - File: `/internal/InfoFlow/InfoIn/semantic_link_service.go`
  - Function: `classifyHierarchyLevel(sourceTerm, targetTerm, relationType) int`
  - Logic:
    ```
    Level 1: Broad concepts (contains "disorder", "disease", "system")
    Level 2: Specific conditions (contains "syndrome", specific disease names)
    Level 3: Clinical details (contains "symptom", "finding", "mutation")
    ```
  - Duration: 1 hour

- [ ] **2.1.3**: Modify `ProcessDocument()` to set hierarchy
  - File: `/internal/InfoFlow/InfoIn/semantic_link_service.go`
  - Change: Set `hierarchy_level` field when creating Weaviate properties
  - Duration: 30 min

- [ ] **2.1.4**: Test hierarchy classification
  - Action: Upload test chapter
  - Validation: Query Weaviate `{ where: { hierarchy_level: 1 } }`
  - Expected: Returns only root concepts
  - Duration: 30 min

#### Acceptance Criteria:
- [ ] Every new semantic link has `hierarchy_level` set (1, 2, or 3)
- [ ] Can filter semantic links by hierarchy level in Weaviate
- [ ] Root concepts (Level 1) are clearly identifiable
- [ ] 80%+ classification accuracy (manual review of 20 links)

---

### 2.2: Parent/Child Relationship Building (2-3 hours)
**Target**: 2025-01-06 Afternoon

#### Tasks:
- [ ] **2.2.1**: Create relationship builder service
  - File: `/internal/InfoFlow/InfoIn/link_relationship_builder.go` (NEW)
  - Function: `BuildRelationships(links []SemanticLink) error`
  - Logic: Analyze `source_term`, `target_term`, `relation_type` to infer parent/child
  - Duration: 1.5 hours

- [ ] **2.2.2**: Update existing links with relationships
  - File: `/internal/InfoFlow/InfoIn/link_relationship_builder.go`
  - Function: `UpdateLinkRelationships(ctx context.Context, linkID string, parentIDs, childIDs []string) error`
  - Duration: 1 hour

- [ ] **2.2.3**: Test hierarchical traversal
  - Action: Query root concept, traverse to children
  - Validation: Can navigate "Immunodeficiency" → "XLA" → "BTK mutation"
  - Duration: 30 min

#### Acceptance Criteria:
- [ ] `is_parent_of` and `is_child_of` arrays populated for 70%+ of links
- [ ] Can traverse from any link to its parents and children
- [ ] No circular references in parent/child relationships

---

### 2.3: Multi-Hop Graph Traversal (3-4 hours)
**Target**: 2025-01-07

#### Tasks:
- [ ] **2.3.1**: Design traversal algorithm
  - Document: `/docs/hsg_threads/THREAD002_graph_traversal.md` (NEW)
  - Algorithm: Breadth-first search from query match
  - Max depth: 3 levels
  - Duration: 1 hour

- [ ] **2.3.2**: Implement `TraverseHierarchy()` function
  - File: `/internal/InfoFlow/InfoIn/hsg_query_service.go`
  - Function signature:
    ```go
    func (h *HSGQueryService) TraverseHierarchy(
        ctx context.Context,
        rootLinkID string,
        maxDepth int,
    ) (*HSGContext, error)
    ```
  - Duration: 2 hours

- [ ] **2.3.3**: Update `QueryHSG()` to use traversal
  - File: `/internal/InfoFlow/InfoIn/hsg_query_service.go`
  - Change: Replace flat search with hierarchical traversal
  - Duration: 1 hour

#### Acceptance Criteria:
- [ ] Can traverse graph 3 levels deep
- [ ] Returns hierarchical context (parent → child → grandchild)
- [ ] Includes provenance (which links led to which content)
- [ ] Performance: <2s for traversal of 100 links

---

## Phase 3: Tier 2 SummaryChunk Implementation 🟡 NOT STARTED

**Goal**: Enable abstract, high-level semantic search

### 3.1: SummaryChunk Weaviate Class (1 hour)
**Target**: 2025-01-08 Morning

#### Tasks:
- [ ] **3.1.1**: Define SummaryChunk schema
  - File: `/internal/InfoFlow/InfoStore/db/weaviate.go`
  - Properties:
    ```
    - summary_text: text (searchable)
    - abstraction_level: int (1=most abstract, 3=most specific)
    - linked_semantic_links: text[] (array of SemanticLink IDs)
    - domain: text ("immunology", "neurology", etc.)
    - created_at: date
    ```
  - Duration: 30 min

- [ ] **3.1.2**: Implement `CreateSummaryChunkClass()`
  - Function: Create Weaviate class with AWS Titan vectorization
  - Duration: 30 min

#### Acceptance Criteria:
- [ ] `SummaryChunk` class exists in Weaviate
- [ ] Can store and retrieve summary chunks
- [ ] AWS Titan generates embeddings for summaries

---

### 3.2: Background Summarization Job (4-5 hours)
**Target**: 2025-01-08 Afternoon + 2025-01-09

#### Tasks:
- [ ] **3.2.1**: Create summarization service
  - File: `/internal/InfoFlow/InfoIn/tier2_summarization.go` (NEW)
  - Function: `SummarizeSemanticLinkGroup(links []SemanticLink) (*SummaryChunk, error)`
  - Duration: 2 hours

- [ ] **3.2.2**: Implement grouping logic
  - Strategy: Group links by topic similarity (vector clustering)
  - Group size: 10-20 links per summary
  - Duration: 1.5 hours

- [ ] **3.2.3**: LLM summarization with Claude
  - Prompt: "Summarize the relationships between these medical concepts..."
  - Model: Claude 3.5 Sonnet
  - Duration: 1 hour

- [ ] **3.2.4**: Store summaries in Weaviate
  - Back-reference: Store IDs of constituent links
  - Metadata: abstraction_level, domain, created_at
  - Duration: 30 min

#### Acceptance Criteria:
- [ ] Background job runs every 1 hour
- [ ] Groups related semantic links automatically
- [ ] Generates coherent summaries using Claude
- [ ] Stores summaries with back-references to links
- [ ] At least 10 summary chunks created from test data

---

### 3.3: Two-Tier Query System (2-3 hours)
**Target**: 2025-01-10

#### Tasks:
- [ ] **3.3.1**: Implement Tier 2 search
  - File: `/internal/InfoFlow/InfoIn/hsg_query_service.go`
  - Function: `SearchSummaries(query string) ([]SummaryChunk, error)`
  - Duration: 1 hour

- [ ] **3.3.2**: Implement Tier 1 refinement
  - Logic: From matched summary, retrieve constituent links
  - Duration: 1 hour

- [ ] **3.3.3**: Integrate with RAG answer generation
  - Context: Hierarchical (Summary → Links → Raw Content)
  - LLM: Claude generates answer with full provenance
  - Duration: 1 hour

#### Acceptance Criteria:
- [ ] Query first searches SummaryChunk (Tier 2)
- [ ] Then traverses to SemanticLinks (Tier 1)
- [ ] Then retrieves raw content from MongoDB
- [ ] Final answer includes hierarchical context
- [ ] Query "What causes XLA?" returns full causal chain

---

## Phase 4: Production Hardening 🔵 FUTURE

**Goal**: Make HSG system production-ready

### Tasks (High-Level):
- [ ] Add monitoring and logging (Prometheus/Grafana)
- [ ] Performance optimization (caching, batch processing)
- [ ] Error handling and retry logic
- [ ] Unit and integration tests (80% coverage)
- [ ] Load testing (100 concurrent queries)
- [ ] Documentation for API endpoints
- [ ] Deployment automation (CI/CD)

**Target**: 2025-01-15 - 2025-01-20

---

## Milestones & Checkpoints

| Milestone | Target Date | Status | Deliverable |
|-----------|-------------|--------|-------------|
| Phase 1 Complete | 2025-01-05 | ✅ Done | Tier 1 semantic links working |
| Phase 2.1 Complete | 2025-01-06 AM | ⏳ In Progress | Hierarchy classification |
| Phase 2.2 Complete | 2025-01-06 PM | ⏳ Waiting | Parent/child relationships |
| Phase 2.3 Complete | 2025-01-07 | ⏳ Waiting | Multi-hop traversal |
| Phase 3.1 Complete | 2025-01-08 AM | ❌ Not Started | SummaryChunk class |
| Phase 3.2 Complete | 2025-01-09 | ❌ Not Started | Background summarization |
| Phase 3.3 Complete | 2025-01-10 | ❌ Not Started | Two-tier query system |
| Phase 4 Complete | 2025-01-20 | ❌ Not Started | Production ready |

---

## Success Metrics

### Phase 2 (Hierarchy) Success:
- [ ] 90%+ of semantic links have `hierarchy_level` set
- [ ] Can query by hierarchy level in Weaviate
- [ ] Multi-hop traversal works for 95% of test queries
- [ ] Traversal completes in <2 seconds

### Phase 3 (Tier 2) Success:
- [ ] At least 50 summary chunks created
- [ ] Query latency <3 seconds (end-to-end)
- [ ] Answer quality: 4/5+ rating (manual review)
- [ ] Provenance accuracy: 95%+ (correct source attribution)

---

## Risk Management

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Claude API fails | Medium | High | Implement retry logic + fallback to Titan |
| Hierarchy classification inaccurate | Medium | Medium | Manual review + iterative refinement |
| Vector DB performance issues | Low | High | Implement caching + index optimization |
| Summarization quality poor | Medium | High | Refine prompts + use longer context window |

---

## Dependencies

### External Services:
- ✅ AWS Textract (OCR)
- ✅ AWS S3 (file storage)
- ✅ AWS Titan (embeddings)
- ⚠️ AWS Bedrock Claude (generative - needs API fix)
- ✅ MongoDB (content storage)
- ✅ Weaviate (vector DB)

### Internal Services:
- ✅ Semantic link extraction
- ⏳ Hierarchy classification (in progress)
- ❌ Relationship building (not started)
- ❌ Tier 2 summarization (not started)
- ❌ Multi-hop traversal (not started)

---

**Next Review**: After Phase 2.1 completion (2025-01-06)  
**Updated**: 2025-01-05 20:30 UTC
