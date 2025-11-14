# Semantic Links Index

**Purpose**: This directory tracks individual semantic link requirements and implementation details.

---

## Active Semantic Links

### Tier 1 (Root Requirements)
- [SL001: HSG Requires Hierarchy](./SL001_hsg_requires_hierarchy.md) - ✅ Schema updated, ⏳ Implementation pending

### Tier 2 (Implementation Details)
- **SL002**: Weaviate schema update - ✅ DONE (2025-01-05)
- **SL003**: Semantic link classification logic - ⏳ NEXT TASK
- **SL004**: Parent/child relationship building - ❌ NOT STARTED
- **SL005**: Multi-hop graph traversal - ❌ NOT STARTED

---

## How to Use This Index

1. Each semantic link is a **single, atomic requirement or implementation task**
2. Links are organized by **hierarchy level** (Tier 1 = root concepts, Tier 2 = implementation)
3. Each link has its own file with:
   - Statement (what needs to be done)
   - Current status
   - Dependencies (parent/child links)
   - Implementation steps
   - Acceptance criteria

---

## Creating New Semantic Links

Use this template:

```markdown
# SLXXX: [Title]

**Type**: [Requirement|Implementation|Bug Fix]  
**Source**: [Where this came from]  
**Target**: [What it affects]  
**Relation**: [requires|implements|fixes]  
**Hierarchy Level**: [1|2|3]  

## Statement
[Clear description of what needs to be done]

## Current Status
[✅ DONE | ⏳ IN PROGRESS | ❌ NOT STARTED]

## Dependencies
- **Parent Links**: [List of prerequisite SLs]
- **Child Links**: [List of dependent SLs]

## Implementation Steps
1. [Step 1]
2. [Step 2]
...

## Acceptance Criteria
- [ ] [Criterion 1]
- [ ] [Criterion 2]
...

## References
- Code: [File paths]
- Docs: [Related documentation]
```

---

**Last Updated**: 2025-01-05  
**Next Link ID**: SL006
