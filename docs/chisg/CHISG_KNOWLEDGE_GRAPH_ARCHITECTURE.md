# CHISG Knowledge Graph Architecture

## Summary

A structured approach to capturing knowledge claims with provenance metadata that enables downstream reasoning, filtering, and logical validation.

---

## Core Principle

**Extract what sources claim. Tag where they came from. Let coherence emerge from the graph.**

CHISG does not judge truth during extraction. It faithfully captures claims with rich metadata. Truth is a computed property - consistency across sources, not correspondence to reality.

---

## Core Capability: Structural Analogy Detection

CHISG's primary teaching value is **surfacing structural analogies across domains**. When the same relation pattern appears in different contexts, CHISG can automatically detect and surface these connections.

### How It Works

The semantic link structure `[Entity A, Forward Link, Entity B]` enables pattern matching across domains:

| Domain | Entity A | Relation | Entity B |
|--------|----------|----------|----------|
| Physics | Bimetallic strip | differential expansion causes | curvature |
| Biology | Auxin gradient | differential elongation causes | phototropic bending |
| Biology | Guard cells | differential turgor causes | stomatal curvature |

All three share the **identical relational pattern**: differential X causes curvature.

### Why This Matters

1. **Learn once, apply everywhere**: A student who deeply understands the bimetallic strip has already learned the logic pattern for auxin response - they just need the domain vocabulary mapped onto the same structural relationship.

2. **Teaching using analogies**: This is not an optional enhancement - it is the core requirement for analogy-based teaching. Without automatic detection of structural matches, every analogy must be manually authored.

3. **Cross-domain bridges**: The `context[]` array can tag links with `bridge:structural-analogy` when the same relation pattern appears across different `domain:` tags, enabling targeted retrieval of analogous concepts.

### Implementation Pattern

Query for structural analogies:
```
MATCH (a1)-[r1:causes]->(b1), (a2)-[r2:causes]->(b2)
WHERE r1.type = r2.type 
AND a1.context CONTAINS 'domain:physics'
AND a2.context CONTAINS 'domain:biology'
RETURN a1, b1, a2, b2
```

This enables the system to automatically generate: *"This is like the bimetallic strip you learned in physics - same pattern, different molecules."*

---

## Two-Layer Architecture

| Layer | Purpose | Content Type |
|-------|---------|--------------|
| **Source Layer (CHISG)** | Faithful capture of what sources claim | Extracted claims with provenance |
| **Teaching Layer (LAO)** | Pedagogical enrichment | Analogies, narratives, mnemonics |

The teaching layer is generated from (and cites) the source layer, but adds interpretive value.

---

## The Semantic Unit

### Core Structure
```
[Entity A, Forward Link, Entity B, Backward Link, Context[]]
```

**Example**:
```
["kinetic energy", "calculated using", "½mv²", "used to calculate", 
 ["source:pearson-physics", "domain:physics", "level:gcse", "topic:energy"]]
```

### Context Array

Searchable, filterable elements:
- `source:` - publication identifier
- `domain:` - subject area
- `topic:` - hierarchical topic path
- `level:` - educational level (gcse, alevel, undergraduate)
- `exam_board:` - awarding body
- `chapter:`, `page:` - location in source

---

## Claim Validation & Trust Scoring (Triangulation)

To eliminate AI hallucinations and ensure structural integrity, CHISG ranks claims based on **Triangulation Density** rather than binary truth.

### 1. Source Density (Breadth)
Ranking is determined by the number of **independent contemporaneous sources** claiming the same semantic link.
- `High Density`: A link appearing in multiple textbooks and official specifications.
- `Low Density`: A link appearing in only one source or a transient observation.

### 2. Path Density (Depth/Routes)
Ranking is determined by the number of independent **logical routes** that lead to the same conclusion.
- If `A -> B` is a direct claim, its score is $S$.
- If `A -> C -> B` is also a valid path, the "Trust Score" for the relationship between A and B increases.
- This represents **Structural Robustness**. Multiple routes to the same idea imply that the relationship is a fundamental pillar of the domain's logic, not a peripheral detail.

### 3. Computed Trust Value
Trust is not a static flag but a computed score:

$$\mathcal{T}(A \xrightarrow{r} B) = \phi(M_a) \cdot \left[ 1 - e^{-( \alpha \ln(1+N_s) + \beta N_p )} \right]$$

Where $N_s$ = independent source count, $N_p$ = independent path count, $M_a$ = cross-domain analogy matches, $\phi(M_a)$ = analogy scaling factor, and $\alpha, \beta$ = weighting coefficients. The saturating exponential ensures diminishing returns from additional sources, while $\phi(M_a)$ amplifies trust when the same structural pattern is confirmed across domains.

Relationships with low $\mathcal{T}$ scores are treated as "hypotheses" and are flagged for human review or excluded from automated coaching responses.

### 4. Opinion vs. Structural Truth (The Echo Chamber Filter)
A high `SourceCount` can be misleading if it represents **Consensus by Correlation** (everyone repeating the same person) rather than **Consensus by Derivation** (everyone finding the same truth independently).

- **Independence Weighting**: The system analyzes the "Provenance Pedigree." If ten sources link back to a single original publication, they are clustered as a single "Opinion Node." 
- **The Opinion Buffer**: Claims that have high `SourceCount` but low `PathDensity` (only one way to describe it) are categorized as **"Subjective Narratives"** or **"Social Norms."**
- **Validation Rule**: Truth in CHISG requires **Structural Redundancy**. If a "popular opinion" cannot be cross-verified by a different logical route (e.g., a scientific model or a complementary domain), it remains an *opinion attribute* of the source, not a *fact node* in the graph.

### 5. The Bootstrap Paradox (Expert-Lead to Empirical-Lead)
The current Skills Map is undergoing a transition from **Expert-Heuristic** to **Empirical-Truth**.

- **Phase 1: Expert Bootstrap (Current)**: The skills map is "unsourced" in the academic sense, relying on the condensed expertise of the lead consultant. These nodes are treated as **Foundational Hypotheses**.
- **Phase 2: Empirical Backfilling**: As the system becomes diagnostic, the "Primary Source" becomes **Observable Intervention Results**. 
    - *The Logic*: "If exercising Skill X (as defined here) leads to a measurable improvement in Outcome Y, then the definition of Skill X is validated."
- **Phase 3: Academic Grounding**: Retroactively mapping historical and pedagogical research (backfilling PhD-level sourcing) to the validated nodes to complete the provenance chain.

---

## Universal Motif: Skills & Competencies (HumanOS Integration)

The CHISG structure is universal. The Skills Map and HumanOS are not distinct systems; they are specialized clusters within the graph that follow the same relational logic.

### 1. Skills as Graph Motifs
Skills and their criteria are represented as nodes with specific semantic links:

- `is_a`: Hierarchical classification (e.g., `[Working Memory] -[:is_a]-> [Executive Function]`).
- `component_of`: Part-whole relationships (e.g., `[Decoding] -[:component_of]-> [Reading Fluency]`).
- `is_required_for`: Logical prerequisites (e.g., `[Phonemic Awareness] -[:is_required_for]-> [Blending]`).
- `evidences`: Linking behavioral "Tells" to skill attainment (e.g., `[Micro-withdrawal] -[:evidences]-> [Low Resilience]`).
- `mitigates`: How a skill alters a psychological state (e.g., `[Metacognition] -[:mitigates]-> [Impulsivity]`).

### 2. HumanOS: The Behavioral Engine
HumanOS characterizes the dynamic state of the student. While Skills represent *capacity*, HumanOS represents *disposition* and *real-time response*.

**The Pipeline:**
`[Stimuli] -> [Processing Mechanisms] -> [Emotional Trigger Point (ETP)] -> [Behaviour]`

Within CHISG, these are nodes:
- `:Stimulus`: Environmental or academic input.
- `:Mechanism`: Evolutionary instincts, social norms, peer pressure, personal experience.
- `:ETP`: Specific emotional trigger states (e.g., "Risk Aversion", "Pain Tolerance").
- `:Behaviour`: The observed output (e.g., "Right vs Left choice", "Withdrawal").

### 3. The Philosophy of Dissolution
*Question: Does the HumanOS just dissolve into the Skills Map?*

**Resolution:** In a Knowledge Graph, they are unified. However, they serve different functional roles:
- **Skills (Attainment)**: Static tracking of what a student *can* do.
- **HumanOS (Insight)**: Dynamic tracking of *why* a student does what they do.

By linking them (e.g., `[:Skill]-[:regulates]->[:ETP]`), we avoid "hallucinations" by providing a structural map of human behavior. The graph ensures that an AI doesn't just guess a student's motive; it traverses the logic from Stimulus to ETP to Behaviour.

---

## Metadata Fields

### GCSE Implementation (Phase 1)

| Field | Values | Purpose |
|-------|--------|---------|
| `knowledge_type` | declarative, procedural, conditional, misconception, confusable, model | Enables targeted pedagogy |
| `pedagogical_delta` | 1-10 | Complexity level for sequencing |
| `context[]` | Array of tagged strings | Filtering and search |

### Research Extension (Phase 2)

| Field | Values | Purpose |
|-------|--------|---------|
| `source_hierarchy` | primary, secondary, tertiary | Source authority |
| `predictive_power` | high, medium, low, none | Empirical backing |
| `publication_type` | journal, book, news, opinion, blog, framework | Source type |
| `peer_reviewed` | boolean | Academic validation |
| `author_expertise` | domain-expert, adjacent-expert, informed-amateur, unknown | Author credibility |
| `claim_type` | fact, prediction, interpretation, model, opinion | Nature of claim |
| `validation_status` | empirically-tested, peer-reviewed, informal, untested | Verification level |
| `temporal_status` | established, current, predicted | Time dimension |

---

## Knowledge Types

| Type | Example | Pedagogical Handling |
|------|---------|---------------------|
| `declarative` | "Electrons have negative charge" | State fact, check recall |
| `procedural` | "Rearranging equations" | Sequence steps, practice |
| `conditional` | "Use F=ma when calculating resultant force" | Teach application context |
| `misconception` | "Heavier objects fall faster" | Surface error, correct explicitly |
| `confusable` | "Speed vs Velocity" | Contrast, distinguish features |
| `model` | "Particle model of matter" | Explain scope and limitations |

### Pedagogical Link Types

- `is_misconception_of` / `corrected_by`
- `confused_with` / `distinguished_by`
- `applies_when` / `limited_by`
- `requires_understanding_of` / `enables_understanding_of`

---

## Logical Validation (Computed, Not Stored)

### Checks
1. **Internal Consistency**: Contradictions within same source
2. **Logical Validity**: Valid inference structure
3. **Hidden Premises**: Unstated assumptions
4. **Circular Reasoning**: Self-referential loops
5. **Transitivity Violations**: Broken logical chains

### Cache Strategy
- **Dirty flagging** with K-hop propagation (2-3 hops)
- No dependency lists stored (avoids bloat)
- Recompute on query if dirty

---

## LLM + KG Hybrid

| Component | Role |
|-----------|------|
| **Knowledge Graph** | Grounding, provenance, contradiction detection, structured retrieval |
| **LLM** | Fluency, generation, synthesis, pedagogical framing |

**Pattern**: Enhanced RAG where LLM generates from KG context, can cite sources, and detects gaps ("I don't know" is possible when graph has no relevant links).

---

## Implementation Plan

### Phase 1: GCSE/LAO (Now)
1. Update SemanticLink model with `context[]` array
2. Add `knowledge_type` enum
3. Add `pedagogical_delta` integer
4. Update UI to capture new fields
5. Build LAO MVP with simple model

### Phase 2: Research Extension (Later)
1. Add extended provenance fields
2. Implement logical validation functions
3. Add computed consensus scoring
4. Build contradiction detection UI

---

## Use Cases Validated

| Scenario | Required Fields | Works? |
|----------|-----------------|--------|
| Student revision query | context[], pedagogical_delta | ✅ |
| Exam prep filtering | context[], knowledge_type | ✅ |
| Misconception surfacing | knowledge_type, corrected_by link | ✅ |
| Concept dependencies | forward/backward links | ✅ |
| Contradiction detection | computed coherence | ✅ (Phase 2) |
| Source reliability check | source_hierarchy, predictive_power | ✅ (Phase 2) |

---

## Key Insight

For GCSE textbook content, heavy validation is unnecessary - the content is already curated. The rich architecture exists for:
1. Handling diverse sources (journals, opinions, frameworks)
2. Detecting contradictions across sources
3. Enabling AI to make nuanced judgments about claim reliability

Build simple, design extensible.

---

## Future Capability: Adaptive Personalisation

### Concept

Integrate personality profiling (Human-OS ETP framework) with KG traversal weights to deliver personalised pedagogy.

### Mechanism

```
Student Profile (Traits) → Tuned Presentation Weights → Content Selection/Ordering
         ↓                                                        ↓
   Behaviour Signals  ←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←  Outcome Metrics
         ↓                                                        ↓
                    Correlation Analysis → Refined Weights
```

### Personality → Pedagogy Mapping

| Trait | Learning Weight Adjustment |
|-------|---------------------------|
| High Certainty | Prioritise procedural, reduce ambiguity |
| High Curiosity | Show edge cases, model limitations |
| High Social | Peer comparisons, collaborative examples |
| High Openness | Alternative explanations, "why" emphasis |
| High Neuroticism | Confidence building, error normalisation |

### Data Flow

1. **Initial**: Cluster-based defaults from personality assessment
2. **Behavioural**: Refine from session data (time on content, skip patterns, revisits)
3. **Outcome**: Calibrate from assessment results
4. **Meta-learning**: Patterns discovered across students inform new defaults

### Ethical Constraints

- Student controls their data
- Transparency in why content is ordered
- Override capability for preferences
- Anonymised aggregation only for insights

### Implementation Phase

**Deferred until after LAO MVP ships.** This is Phase 3+ capability.

---

## Temporal Update Patterns

Different domains require different recalculation strategies:

| Domain | Update Trigger | Latency Tolerance | Strategy |
|--------|---------------|-------------------|----------|
| GCSE | Annual spec release | Days/weeks | Batch recalculation, versioned snapshots |
| Research | Paper ingestion | Hours/days | Event-driven, priority zones |
| Clinical | Data stream + guidelines | Seconds (critical) to months (stable) | Tiered by criticality class |

### Node Update Classes

```
update_class:
  - "realtime"  → streaming, no cache
  - "daily"     → overnight batch
  - "periodic"  → trigger-based (spec release)
  - "stable"    → indefinite cache
```

### Meta-Metadata for Self-Tuning

The system can learn optimal weights from outcomes:

```
meta_config: {
  weights: { current tuned values },
  success_metrics: [ tracked outcomes ],
  insights: [ discovered correlations ],
  adjustments: [ tuning history ]
}
```

Tuning parameters become first-class knowledge claims with provenance.

---

## Scope Decision: Codebase KG

**Status: DEFERRED**

Building a knowledge graph of project code/plans was considered and rejected as a distraction. The current tooling (AI assistant + logs + conversation summaries) provides sufficient project memory. Revisit only after LAO ships.

