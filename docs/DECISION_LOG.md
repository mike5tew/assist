# Decision Log

## 2025-01-05 - Avoid Rewriting Upload Handlers
**Temptation**: Consolidate `immunology_routes.go` and `upload_handlers.go`
**Decision**: DEFER - this is refactoring, not HSG implementation
**Rationale**: Current upload system works. Focus on adding hierarchy to semantic links first.

## 2025-01-05 - Don't Build New SemanticLinks Schema Yet
**Temptation**: Start fresh with a new Weaviate class
**Decision**: MODIFY existing `SemanticLinks` class
**Rationale**: We have working semantic links. Add fields, don't rebuild.

## Future Template
**Temptation**: [What I want to do]
**Decision**: [DO IT / DEFER / DON'T DO IT]
**Rationale**: [Does this directly implement Tier 2 or multi-hop traversal?]
