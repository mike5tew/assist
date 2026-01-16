# Decision Log

## 2025-12-27 - CHISG Metadata Architecture: Simple vs Research Grade

**Temptation**: Implement full research-grade knowledge graph with logical validation, consensus scoring, and comprehensive provenance tracking for GCSE revision tool.

**Decision**: **Two-tier approach. GCSE uses simplified metadata. Research-grade features are designed but deferred.**

**Rationale**:
1. **GCSE Needs**: Source attribution, topic/level filtering, basic quality indicators. Textbook content is already curated - heavy validation is redundant.
2. **Research Needs**: Provenance chains, claim type classification, logical consistency checks, computed consensus scores.
3. **Architecture Compatibility**: Design the data model to support both. GCSE ignores extended fields; research uses them.
4. **Build Order**: Ship LAO with simple model. Extend for research use cases later.

**GCSE Fields** (implement now):
- Context array (source, domain, topic, level, exam_board)
- knowledge_type (declarative, procedural, conditional, misconception, confusable)
- pedagogical_delta (1-10 complexity)

**Research Fields** (design now, implement later):
- source_hierarchy, predictive_power, author_expertise
- claim_type, validation_status, temporal_status
- Computed: coherence_score, consensus_score

## 2025-12-27 - Context as Array, Not String

**Temptation**: Store context as single text field for simplicity.

**Decision**: **Context must be an array of searchable elements.**

**Rationale**:
1. **Searchability**: Can filter by any dimension (level:gcse AND topic:energy)
2. **Flexibility**: Add new context types without schema changes
3. **Structure**: Prefix convention (`type:value`) enables parsing while remaining grep-friendly

**Format**: `["source:pearson-physics-2024", "domain:physics", "level:gcse", "topic:energy.conservation"]`

## 2025-12-27 - Knowledge Types for Pedagogical Precision

**Temptation**: Treat all claims as equivalent facts.

**Decision**: **Classify knowledge_type to enable misconception handling, procedure teaching, and conditional application.**

**Types**:
- `declarative` - facts (electrons are negative)
- `procedural` - how-to (rearranging equations)
- `conditional` - when to apply (use F=ma when...)
- `misconception` - common errors (heavier falls faster)
- `confusable` - easily confused pairs (speed vs velocity)
- `model` - simplified representations (particle model)

## 2026-01-16 - The Bootstrap Paradox: Skills Map Sourcing

**The Issue**: The current Skills Map is "unsourced" (based on personal expert heuristics rather than a literature review of 1,000+ sources).

**Decision**: **Accept "Expert Bootstrap" as the initial validation, with a pivot to "Empirical Truth" through intervention results.**

**Rationale**:
1. **Speed to Value**: Waiting for PhD-level sourcing is a project-killer.
2. **Empirical Validation as a Source**: In the CHISG model, "Observable Success" is a primary source. If an intervention based on a "broken-down skill definition" improves a student's outcome, that success *is* the validation of the definition.
3. **Backfilling Plan**: Use the diagnostic success data of the project to justify and fund the retroactive academic "backfilling" of sources later.

**Result**: We stop asking "is this sourced?" and start asking "does this definition lead to effective interventions?"

## 2026-01-22 - Shift to Objective Temporal Markers (Behaviour Lite)

**The Issue**: Subjective "Physical Tells" (e.g., "looked at friend") introduce teacher bias and lack evidentiary weight for profiling.

**Decision**: **Adopt a 4-Marker system (2 Positive / 2 Negative) based on frequency and duration, mapped to 15 Personality/Sensitivity sliders.**

**Rationale**:
1. **Evidence-Based Calibration**: Profiles (Sliders) are calibrated based on occurrences of specific, objective markers.
2. **Mandatory Justification**: Any downward shift on a slider requires a manual note referencing the markers, ensuring accountability.
3. **Cognitive Load Reduction**: Limiting to precisely 4 markers (Behaviour Lite) keeps the "eyes-up" teaching goal viable while removing the ambiguity of open-ended physical tells.
4. **Data Integrity**: Moves HumanOS from a "description" tool to a "diagnostic" tool that furnishes proof for skill attainment.

**Result**: UI updated to use toggle markers (Temporal Data) and a slider-based modal (Calibration Data).

## 2026-01-16 - Deployment Strategy: Sampling & Outliers

**The Issue**: Real-time logging of 30 students is cognitively impossible for a teacher while teaching.

**Decision**: **Move from "Census Logging" to "Sampling & Extreme Capture."**

**Rationale**:
1. **The Pareto of Insight**: Most useful data comes from the extremes (the very quick to finish, the very slow to engage).
2. **Rotational Sampling**: Teachers focus on a deep-dive of ~3 students per hour. Over a week, a full class profile emerges without cognitive overload.
3. **The Reassurance Effect**: Being "watched" is re-framed from "Surveillance" to "Reassurance" (Culture of Care). Students perform differently when they know their effort is being seen and logged in a meaningful way.
4. **Future-Proofing**: The manual marker interface prepares the ground for future automated inputs (e.g., Computer Vision for safeguarding/engagement detection).

**Result**: UI will prioritize "Extreme" shortcuts and "Focus Student" rotations.

**Link Types Added**: `is_misconception_of`, `corrected_by`, `confused_with`, `distinguished_by`, `applies_when`, `limited_by`

## 2025-12-27 - LLM + KG Hybrid Architecture

**Temptation**: Use KG as replacement for LLM knowledge.

**Decision**: **Complementary layers. LLM for fluency/generation, KG for grounding/verification.**

**Rationale**:
1. **LLM Weakness**: No provenance, can't update surgically, hallucinates
2. **KG Weakness**: Expensive to build, requires extraction quality
3. **Combined**: LLM generates from KG context, cites sources, detects gaps

**Pattern**: Enhanced RAG with structured semantic links rather than flat document chunks.

## 2025-12-22 - Multi-Store Architecture for Semantic Links
**Temptation**: Store semantic links only in MongoDB/Weaviate and query from there for all use cases.
**Decision**: **Design semantic links to be exportable to multiple stores: MongoDB (primary), Weaviate (vector search), and SQLite/Realm (mobile offline-first).**
**Rationale**:
1. **Primary Store (MongoDB)**: All semantic links are first stored in MongoDB `semantic_links` collection. This is the source of truth.
2. **Vector Store (Weaviate)**: Links are asynchronously indexed in Weaviate via `text2vec-aws` for semantic similarity search. The `statement` field is vectorized.
3. **Mobile Store (SQLite/Realm)**: For the GCSE Revision Tool mobile app, semantic links need to be exportable to SQLite or MongoDB Realm for offline-first local storage. This enables the app to work without internet.
4. **Export Format**: Semantic links should be exportable as JSON that can be imported into SQLite tables with matching schema.
**Action**: ✅ COMPLETE
- MongoDB schema already supports all fields
- Weaviate indexing is in place (async via `StoreSemanticLinks`)
- Export endpoint implemented: `GET /api/semantic-links/export?format=json|sqlite-sql&domain=...&limit=1000`

## 2025-12-22 - Statement-Centric Semantic Link Model (Bidirectional Relations)
**Temptation**: Store term-pairs with a single relationship type, requiring two separate records for bidirectional traversal.
**Decision**: **Store statement + bidirectional relations in a single record.**
**Rationale**:
1. **Statement is the Unit of Truth**: The vectorized field should be a complete, citable fact ("XLA causes reduction in Ig"), not a term-pair. This supports the Retrieval-Constrained Generation (RCG) architecture where LLMs compose retrieved facts rather than generating from parametric knowledge.
2. **Single Record, Two Directions**: One record with `forward_relation` and `inverse_relation` fields. When querying from either direction, the appropriate relation is used.
3. **Conditions as Embedded Array**: Conditions (e.g., "occurs when maternal antibodies wane") are embedded as an array in the record, each with their own bidirectional relations. This keeps related data together.
4. **Query Pattern**: Find by `source_term`, `target_term`, or `conditions.term` in MongoDB. Vector search on `statement` in Weaviate.
**Action**: Updated `SemanticLink` struct, API handlers, and frontend to use `statement`, `forward_relation`, `inverse_relation`, and `conditions[]`.

## 2025-12-22 - Embedding Strategy: Weaviate text2vec-aws + Bedrock Titan
**Temptation**: Use free/open-source embedding models (e.g., sentence-transformers) or generate embeddings manually in Go code.
**Decision**: **Use AWS Bedrock Titan Embed v2 via Weaviate's text2vec-aws module.**
**Rationale**:
1. **Embedding Quality**: Titan Embed v2 produces higher-quality embeddings for semantic similarity tasks compared to many open-source alternatives. This is critical for the RCG architecture where retrieval precision directly impacts output quality.
2. **Automatic Vectorization**: Weaviate's `text2vec-aws` module calls AWS Bedrock Titan Embed v2 automatically when documents are inserted. No manual embedding code required in the application.
3. **Configuration**: Weaviate is configured with `DEFAULT_VECTORIZER_MODULE: text2vec-aws` and `ENABLE_MODULES: text2vec-aws`. AWS credentials are passed via environment variables.
4. **Model**: Using `amazon.titan-embed-text-v2:0` with 1024-dimension vectors (v2 has better multilingual and domain-specific performance than v1).
5. **What Gets Vectorized**: The `statement` field (the full fact) is the primary vectorization target. Weaviate will embed this automatically.
**Action**: Semantic links are inserted into Weaviate asynchronously after MongoDB storage. The `statement` field is indexed for semantic search.

## 2025-01-05 - Avoid Rewriting Upload Handlers
**Temptation**: Consolidate `immunology_routes.go` and `upload_handlers.go`
**Decision**: DEFER - this is refactoring, not HSG implementation
**Rationale**: Current upload system works. Focus on adding hierarchy to semantic links first.

## 2025-01-05 - Don't Build New SemanticLinks Schema Yet
**Temptation**: Start fresh with a new Weaviate class
**Decision**: MODIFY existing `SemanticLinks` class
**Rationale**: We have working semantic links. Add fields, don't rebuild.

## 2025-01-13 - Database Strategy for Integration
**Temptation**: Migrate the `skills-map-platform`'s MySQL database to MongoDB/Weaviate to create a single, unified data store.
**Decision**: **DO NOT MIGRATE.** Adopt a hybrid data store model. The `esp-organizer` backend will be enhanced to connect to both the existing MongoDB/Weaviate stack AND the `skills-map-platform`'s MySQL database.
**Rationale**:
1.  **Preserve Investment**: The `init.sql` for the MySQL database represents months of work modeling complex, relational school data (like CTF files). Migrating this to NoSQL would be a massive undertaking and would lose the relational integrity.
2.  **Right Tool for the Job**: Relational data (student records, course assignments) is best handled by MySQL. Semantic/vector data (knowledge graphs, embeddings) is best handled by Weaviate/MongoDB. A hybrid approach uses each tool for its strengths.
3.  **Long-Term Vision**: The final goal of a full school MIS will require a robust relational database. Keeping MySQL aligns with this long-term requirement.
4.  **De-risks Integration**: Connecting to an existing, working database is far less risky than attempting a complex data migration.

**Action**: The `esp-organizer` backend will add a MySQL connector. Handlers from the `skills-map-platform` API will be merged into `esp-organizer` and will use this new connector, while the AI Tutor features will continue to use the existing Mongo/Weaviate connectors.

## 2025-01-22 - Dataset Sourcing for AI Training (Golden Set) - REVISED
**Temptation**: Only use public domain texts to create the "Golden Set" of semantic links, avoiding all copyright risk.
**Decision**: **Use curriculum-aligned textbooks (even if copyrighted) as a reference for manually creating the Golden Set.** This is acceptable under a "transformative use" rationale with strict conditions.
**Rationale**:
1.  **Curriculum Alignment is Paramount**: The primary challenge is modeling a specific curriculum's structure. Generic public domain texts do not meet this requirement. Textbooks are the only reliable source for this structure.
2.  **Use is Transformative**: The copyrighted text is used only as a reference to create a new, distinct work: a dataset of semantic relationships. The model learns the *process* of creating these links, not the book's content.
3.  **No Reproduction**: The critical mitigation is that **the original source text is never stored, retrieved, or reproduced by the system in any form.** The final product contains only the derived semantic links and the trained model. This makes it impossible to reconstruct the original work.
**Action**: The process of creating the Golden Set must be documented to show that the source text is used for reference only and is not ingested into any production data store.

## 2025-01-24 - Defer Microservices Migration
**Temptation**: Break the backend into microservices for authorization, skills, knowledge, etc., to prepare for future scale.
**Decision**: **DEFER.** Continue with a "Modular Monolith" architecture for the `esp-organizer` backend.
**Rationale**:
1.  **Speed**: A monolithic architecture is significantly faster for development and iteration at this early stage. Refactoring and adding features does not require cross-service coordination.
2.  **Simplicity**: The operational overhead of deploying, managing, and monitoring a distributed system is immense and would distract from the primary goal of delivering product features.
3.  **Flexibility**: It is much easier to refactor and define service boundaries within a monolith than it is to change them once they are distributed across a network.
**Action**: We will continue to build within the `esp-organizer` service, but with a strong emphasis on internal modularity (e.g., distinct packages for `auth`, `knowledge`, `skills`). This will allow us to easily extract these modules into true microservices in the future, if and when the scale of the system justifies the added complexity.

## 2025-01-25 - Adopt a Single Monorepo Strategy
**Temptation**: Create a new folder/repository for each product (e.g., a new "Revision Assistant" folder).
**Decision**: **DO NOT CREATE NEW FOLDERS.** Consolidate all active development for HumanOS, CHISG, the Coach/Revision Assistant, and Skills Tree Rising into the existing `/assist` monorepo.
**Rationale**:
1.  **Eliminates Fragmentation**: The user's key pain point is multiple, out-of-date `docs` folders. A single monorepo enforces a single source of truth for all documentation, code, and configuration.
2.  **Enables Code Sharing**: All products are built on the same HumanOS/CHISG foundation. A monorepo makes sharing code between the "Coach" and "Skills Tree" products trivial.
3.  **Supports Modular Monolith**: This decision is the physical implementation of our "Modular Monolith" strategy. The products are logically separate (different UI routes, different API endpoints) but live in one codebase for development speed.
**Action**: The `/assist` directory is now the single source of truth. Other project folders (`HumanOS`, `skillsdocker`, `espdata`) should be considered for archival or have their relevant code merged into `/assist`. The "Revision Assistant" will be built within the existing `/assist/frontend` and `/assist/esp-organizer` structure.

## 2026-01-09 - Narrative Spines and Fundamental Concepts (Exploratory)

**Observation**: While teaching, created compressed narrative chains like:
> DNA (code) → Bases (letters) → Genes (instruction) → Protein (machine) → Mutation (change) → Variation (differences) → Adaptation (advantage) → Survival → Natural Selection

This is a **narrative spine**: minimal sequence with logical dependency that tells a complete causal story and provides hooks for elaboration.

**Insight 1 - Narrative Spines**:
- Different from "important concepts" - it's a compressed *explanatory structure*
- Each term requires the previous (logical dependency)
- Expert teachers do this intuitively
- CHISG could help *discover* or *construct* these from the graph

**Insight 2 - Defining "Fundamental"**:
> Fundamentality = cross-domain frequency of the same relational pattern

A concept is fundamental not because it's popular, but because its **relational pattern** appears across many different domains. The wider the scope of domains, the more fundamental.

| Metric | What It Finds |
|--------|---------------|
| PageRank centrality | Popularity (what gets linked to) |
| Betweenness centrality | Bridging (connects clusters) |
| Cross-domain pattern frequency | **Fundamentality** (same pattern, different domains) |

**Example**: "Differential causes change" appears in:
- Physics (thermal expansion)
- Biology (auxin gradients, turgor)
- Economics (price differentials)
- Chemistry (concentration gradients)

This is more fundamental than a concept appearing frequently within *one* domain.

**Connection to Structural Analogy Detection**: This is computable from existing CHISG architecture. Query for relation patterns that appear across multiple `domain:` tags in the context array.

**LLM Probing Won't Work**: Querying LLM completions would find co-occurrence (popularity), not logical dependency or cross-domain structure. The graph approach is superior for this.

**Status**: Exploratory. Potential features:
1. Surface existing narrative spines from graph
2. Compute fundamentality scores from cross-domain pattern frequency
3. Suggest minimal covering narratives for a topic

## Future Template
**Temptation**: [What I want to do]
**Decision**: [DO IT / DEFER / DON'T DO IT]
**Rationale**: [Does this directly implement Tier 2 or multi-hop traversal?]
