SemanticLink schema — canonical source

Location
- canonical Go struct: `assist/esp-organizer/internal/models/semantic_link.go` (source of truth)
- machine schema (this folder):
  - `semantic_link.schema.json` — JSON Schema (draft-07)
  - `semantic_link.openapi.json` — OpenAPI component wrapper

Purpose
- Use `semantic_link.schema.json` for validation, CI checks, migrations and OpenAPI component generation.

Important validation rules (recommended)
- `statement`, `source_term`, `target_term`, `forward_relation` should be present for any Tier-1 SL used in inference.
- Numeric metrics (`confidence`, `semantic_distance`, `*_generality`, `relationship_strength`, `quality_score`) must be in [0.0, 1.0].

Proposed extensions (NOT in Go struct)
- `code` (string): short stable code for UI/print (eg `PS-F-001`).
- `alsoKnownAs` (string[]): alias / synonym list to prevent duplicates.
- `normalizedLabel` (string): lowercased, punctuation-stripped label for fast matching.
- `deprecated` (boolean), `supersededBy` (string[]): deprecation metadata for code lifecycle.

Next steps
1. Review the schema and confirm any additional required fields (e.g., `student_metadata` shape, `vector` length constraints).
2. I can add CI validation (JSON Schema check) and generate TypeScript types from this schema.
3. If you want aliases/codes in-schema, confirm and I will add them as optional properties and implement DB/index recommendations.

How to use
- Validate incoming SL JSON with the schema before ingesting into Mongo/Weaviate.
- Use `semantic_link.openapi.json` as an OpenAPI `components.schemas` reference for CHISG API endpoints.

If you want, I can now:
- add `alsoKnownAs` + `code` to the canonical schema and open a PR draft, or
- run a repo-wide scan to find candidate skill labels to populate `alsoKnownAs` suggestions.
