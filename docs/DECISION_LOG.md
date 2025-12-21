# Decision Log

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

## Future Template
**Temptation**: [What I want to do]
**Decision**: [DO IT / DEFER / DON'T DO IT]
**Rationale**: [Does this directly implement Tier 2 or multi-hop traversal?]
