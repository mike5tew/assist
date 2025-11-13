# HSG Threads Index

**Purpose**: This directory tracks implementation narratives - multi-step workflows that tell the story of building a feature.

---

## Active Threads

### Phase 2: Hierarchy Implementation
- [THREAD001: Implementing HSG Hierarchy](./THREAD001_implementing_hierarchy.md) - ⏳ IN PROGRESS

### Future Threads
- **THREAD002**: Graph Traversal Algorithm - ❌ NOT STARTED
- **THREAD003**: Tier 2 Summarization Pipeline - ❌ NOT STARTED
- **THREAD004**: Two-Tier Query System - ❌ NOT STARTED

---

## What is an HSG Thread?

An HSG Thread is a **narrative chain** of related semantic links that tells the story of implementing a feature. It includes:

1. **Root Concept** - The high-level goal
2. **Semantic Link Chain** - The ordered steps to achieve the goal
3. **Progress Timeline** - Dates and status updates
4. **Current Blocker** - What's preventing progress
5. **Next Actions** - Concrete, atomic tasks to unblock

---

## Creating New Threads

Use this template:

```markdown
# THREADXXX: [Title]

**Thread Type**: [Implementation|Bug Fix|Refactor]  
**Root Concept**: [High-level goal]  
**Status**: [🟢 COMPLETE | 🟡 IN PROGRESS | 🔴 BLOCKED | ⚪ NOT STARTED]  
**Priority**: [🔴 CRITICAL | 🟡 HIGH | 🟢 MEDIUM | ⚪ LOW]  

## Narrative Summary
[Tell the story of this feature in 2-3 sentences]

## Semantic Link Chain
1. [SL001]: [Description] - [Status]
2. [SL002]: [Description] - [Status]
3. [SL003]: [Description] - [Status]
...

## Progress Timeline
- **YYYY-MM-DD HH:MM** - [Event description]
- **YYYY-MM-DD HH:MM** - [Event description]

## Current Blocker
[Describe what's preventing progress]

## Next Actions (Atomic Tasks)
- [ ] [Specific, actionable task]
- [ ] [Specific, actionable task]

## Success Metrics
- [Metric 1]: [Current] / [Target]
- [Metric 2]: [Current] / [Target]

## References
- Roadmap: [Link to phase in HSG_ROADMAP.md]
- Decision: [Link to DECISION_LOG.md entry]
- Code: [File paths]
```

---

**Last Updated**: 2025-01-05  
**Next Thread ID**: THREAD005
