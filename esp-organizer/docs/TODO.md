# ESP Organizer - Implementation TODO

---

## 🔴 CRITICAL PATH: HSG Hierarchy (Phase 2)

### ✅ DONE (2025-01-05)
- [x] Update Weaviate `SemanticLinks` schema with hierarchy fields
- [x] Create architecture definitions document
- [x] Create roadmap and decision log
- [x] Identify root cause of missing hierarchy data

### ⏳ IN PROGRESS
- [ ] **CURRENT TASK**: Modify semantic link creation to set `hierarchy_level`
  - **File**: `/internal/InfoFlow/InfoIn/semantic_link_service.go`
  - **Function**: `ProcessDocument()`
  - **What to add**: `classifyHierarchyLevel()` logic
  - **Acceptance**: Links in Weaviate have `hierarchy_level` 1/2/3
  - **Time estimate**: 2-3 hours
  - **Started**: Not yet
  - **Blocked by**: None

### ❌ NOT STARTED (Priority Order)

#### Phase 2.1: Hierarchy Classification (2-3 hours total)
1. [ ] Create `classifyHierarchyLevel()` function
   - Input: `sourceTerm, targetTerm, relationType`
   - Output: `int` (1, 2, or 3)
   - Logic: Pattern matching on term semantics
   - Duration: 1 hour

2. [ ] Modify `ProcessDocument()` to call classification
   - Set `properties["hierarchy_level"]` before storing
   - Duration: 30 min

3. [ ] Test hierarchy classification
   - Upload test chapter
   - Query Weaviate for Level 1 concepts
   - Manually review 20 links for accuracy
   - Duration: 1 hour

#### Phase 2.2: Parent/Child Relationships (2-3 hours total)
4. [ ] Create `/internal/InfoFlow/InfoIn/link_relationship_builder.go`
   - Function: `BuildRelationships(links []SemanticLink) error`
   - Logic: Infer parent/child from `relation_type` and `hierarchy_level`
   - Duration: 1.5 hours

5. [ ] Implement `UpdateLinkRelationships()`
   - Update Weaviate with `is_parent_of`, `is_child_of` arrays
   - Duration: 1 hour

6. [ ] Test parent/child traversal
   - Query root concept
   - Traverse to children
   - Verify "Immunodeficiency" → "XLA" → "BTK mutation" path
   - Duration: 30 min

#### Phase 2.3: Multi-Hop Traversal (3-4 hours total)
7. [ ] Design traversal algorithm
   - Document in `/docs/hsg_threads/THREAD002_graph_traversal.md`
   - Breadth-first search, max depth 3
   - Duration: 1 hour

8. [ ] Implement `TraverseHierarchy()` in `hsg_query_service.go`
   - Input: `rootLinkID, maxDepth`
   - Output: `*HSGContext` (hierarchical structure)
   - Duration: 2 hours

9. [ ] Update `QueryHSG()` to use hierarchical traversal
   - Replace flat vector search
   - Return full context chain
   - Duration: 1 hour

---

## 🟡 BACKLOG (Defer Until Phase 2 Complete)

### Code Quality & Refactoring
- [ ] Consolidate upload handlers (`immunology_routes.go` + `upload_handlers.go`)
- [ ] Add health endpoint to API router
- [ ] Fix Claude API Converse call format
- [ ] Add unit tests for semantic link service
- [ ] Add integration tests for HSG query

### Tier 2 Implementation (Phase 3)
- [ ] Create `SummaryChunk` Weaviate class
- [ ] Implement background summarization job
- [ ] Implement two-tier query system (Summary → Links → Content)

### Production Hardening (Phase 4)
- [ ] Add monitoring and logging
- [ ] Performance optimization (caching)
- [ ] Error handling and retry logic
- [ ] Load testing
- [ ] API documentation

---

## 📊 METRICS

### Phase 2 Progress:
- **Semantic links with hierarchy**: 0/100 (target: 90+)
- **Links with parent/child data**: 0/100 (target: 70+)
- **Multi-hop queries working**: NO (target: YES)
- **Time invested**: 4 hours / ~10 hours estimated

### Code Changes:
- **Files modified**: 2 (`weaviate.go`, documentation files)
- **Files to modify**: 3 (`semantic_link_service.go`, `hsg_query_service.go`, new files)
- **Lines of code added**: ~50 (documentation)
- **Lines of code to add**: ~300-400 (implementation)

---

## 🚨 BLOCKERS & RISKS

### Current Blockers:
- None (path is clear)

### Potential Risks:
1. **Classification accuracy** - hierarchy classification may be inaccurate
   - Mitigation: Manual review + iterative refinement
   - Impact: Medium

2. **Graph traversal performance** - may be slow for large graphs
   - Mitigation: Implement caching + limit depth
   - Impact: Low (can optimize later)

3. **Claude API still broken** - affects summarization (Phase 3)
   - Mitigation: Fix Converse API call before Phase 3
   - Impact: High (blocks Tier 2)

---

## 📅 SCHEDULE

| Date | Phase | Deliverable | Status |
|------|-------|-------------|--------|
| 2025-01-05 | Setup | Documentation + Schema | ✅ Done |
| 2025-01-06 AM | 2.1 | Hierarchy classification | ⏳ In Progress |
| 2025-01-06 PM | 2.2 | Parent/child relationships | ⏳ Waiting |
| 2025-01-07 | 2.3 | Multi-hop traversal | ⏳ Waiting |
| 2025-01-08 | 3.1 | SummaryChunk class | ❌ Not Started |
| 2025-01-09 | 3.2 | Background summarization | ❌ Not Started |
| 2025-01-10 | 3.3 | Two-tier queries | ❌ Not Started |

---

**Last Updated**: 2025-01-05 20:30 UTC  
**Next Review**: Tomorrow morning before starting work
