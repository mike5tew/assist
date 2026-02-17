# CHISG: Technical Specification

**Contextualised Hierarchical Iterative Semantic Groupings**

**Author**: Michael Stewart, PhD (Biophysics)  
**Date**: February 2026  
**Contact**: michael.stewart@espthinking.co.uk  
**Live System**: https://espthinking.co.uk

---

## What CHISG Is

CHISG is a methodology for structuring knowledge into a machine-readable graph that solves three problems no current system addresses together:

1. **Hallucination Reduction** — Grounding AI responses in verified, sourced knowledge claims
2. **Cross-Domain Connection Discovery** — Automatically detecting structural analogies between disparate fields
3. **Gap Analysis** — Identifying where our collective knowledge has holes, contradictions, or untested assumptions

It is not a product. It is infrastructure — a knowledge-graph methodology that can be applied to any domain where structured, trustworthy knowledge matters.

---

## The Three Core Aims

### 1. Reducing AI Hallucination

**The Problem**: Large language models generate plausible but incorrect information. In education, law, medicine, and policy, this is dangerous. Current mitigation strategies (prompt engineering, RLHF, retrieval-augmented generation) reduce hallucination frequency but cannot eliminate it because the underlying knowledge remains unstructured text.

**The CHISG Approach**: Build a deterministic grounding layer.

Every knowledge claim in CHISG is stored as a **semantic unit**:

```
[Entity A] —[relation]→ [Entity B]
    + backward relation
    + context[] (source, domain, level, topic, exam board, page)
    + trust score (computed, not assigned)
```

When an AI system generates a response, CHISG provides:

- **Verification**: Does this claim exist in the graph? With what trust score?
- **Contradiction Detection**: Does the graph contain claims that contradict this?
- **Scope Boundaries**: Is this claim supported at the learner's level, or is it a simplification that breaks at higher levels?
- **Provenance**: Which sources support this claim? Are they independent?

The key insight is that **truth is computed from structure, not declared by authority**. A claim supported by multiple independent sources via multiple logical routes has high structural integrity. A claim supported by one source repeated many times has high popularity but low structural integrity. CHISG distinguishes between these.

**Current Implementation**: 579 skills, 1,120 semantic links, 27 course definitions stored in Weaviate (vector database). The Go backend generates embeddings and queries via `nearVector` for semantic matching. GCSE Science content (AQA specification) is the initial domain.

---

### 2. Finding Connections Between Disparate Elements

**The Problem**: Knowledge is taught and stored in silos. Physics concepts are separated from biology concepts, which are separated from psychology concepts. But the underlying structural patterns are often identical. A teacher who sees the connection can use it as an analogy. An AI system without structured knowledge cannot.

**The CHISG Approach**: Structural analogy detection via relational pattern matching.

Every semantic link has the same structure: `[A] —[relation]→ [B]`. When the same relational pattern appears across different domains, CHISG detects it automatically.

**Worked Example**:

| Domain | Entity A | Relation | Entity B |
|--------|----------|----------|----------|
| Physics | Bimetallic strip | differential expansion causes | curvature |
| Biology | Auxin gradient | differential elongation causes | phototropic bending |
| Biology | Guard cells | differential turgor causes | stomatal curvature |

All three share the identical relational pattern: **differential X causes curvature**. A student who deeply understands the bimetallic strip has already learned the logic pattern for auxin response — they need the domain vocabulary mapped onto the same structural relationship.

**Why This Matters**:

- **For teaching**: Analogies are the most powerful pedagogical tool. Currently, every analogy must be manually authored by an expert teacher. CHISG makes analogy detection automatic and exhaustive.
- **For research**: Cross-domain structural patterns reveal deep principles. When the same pattern appears in physics, biology, and economics, it points to a fundamental mechanism worth investigating.
- **For AI systems**: Structured analogy retrieval replaces vague "you might also be interested in" recommendations with precise structural matches: "This is like the bimetallic strip you learned in physics — same pattern, different molecules."

**Implementation**: The `context[]` array on each semantic link includes domain tags. Pattern matching queries identify identical relational structures across different domains:

```
MATCH (a1)-[r1:causes]->(b1), (a2)-[r2:causes]->(b2)
WHERE r1.type = r2.type
AND a1.domain ≠ a2.domain
RETURN a1, b1, a2, b2
```

Links identified as cross-domain structural analogies are tagged with `bridge:structural-analogy` for targeted retrieval.

---

### 3. Gap Analysis: Finding the Cutting Edge of Knowledge

**The Problem**: We don't know what we don't know. Curriculum content is presented as complete, but every field has gaps — relationships that should exist based on structural patterns but have never been tested, or areas where claims contradict each other without resolution.

**The CHISG Approach**: Compute what's missing from the graph.

Three types of gap analysis:

#### 3a. Structural Gaps (Missing Links)

If the graph contains `A → B` and `B → C` but not `A → C`, and the relation types are transitive, this is a structural gap. Either:
- The relationship exists but hasn't been captured yet (data gap)
- The relationship doesn't exist, which is itself interesting (research gap)

#### 3b. Analogy Gaps (Incomplete Patterns)

If a relational pattern appears in three domains but is absent from a fourth where it would be expected, this is an analogy gap. Example: "Differential X causes curvature" appears in physics, botany, and cell biology — does it appear in developmental biology? If not, is that because it doesn't apply, or because no one has looked?

#### 3c. Trust Gaps (Low-Confidence Claims)

The trust scoring system identifies claims that are:
- **High popularity, low structural support** — repeated everywhere but with only one logical path (possible echo chamber)
- **High structural support, low popularity** — independently derivable through multiple routes but rarely stated (possible overlooked insight)
- **Contradicted** — one source claims `A causes B`, another claims `A inhibits B` (unresolved contradiction requiring investigation)

**The Trust Score Formula**:

$$\mathcal{T}(A \xrightarrow{r} B) = \phi(M_a) \cdot \left[ 1 - e^{-( \alpha \ln(1+N_s) + \beta N_p )} \right]$$

Where:
- $\mathcal{T}$ = computed trust score for the claim $A \xrightarrow{r} B$, bounded $(0, 1]$
- $N_s$ = number of independent sources claiming this relationship
- $N_p$ = number of distinct logical routes leading to the same conclusion (structural redundancy)
- $M_a$ = number of times the same relational pattern appears in other domains (analogy matches)
- $\phi(M_a)$ = analogy scaling factor — amplifies trust when the same structural pattern is independently confirmed across domains
- $\alpha, \beta$ = weighting coefficients for source vs path contributions
- The saturating term $1 - e^{-x}$ ensures diminishing returns: the first few independent sources matter most; additional repetitions add progressively less confidence
- The logarithmic $\ln(1 + N_s)$ further compresses source count, reflecting that 10 sources citing one paper are not 10× more trustworthy than one — the echo chamber filter handles source independence before $N_s$ is computed

**Interpretation**: A claim with high $N_s$ but low $N_p$ saturates slowly (consensus by correlation — **opinion**). A claim with high $N_p$ across domains drives $\phi(M_a)$ upward and $\beta N_p$ dominates (consensus by derivation — **structural truth**).

**The Echo Chamber Filter**: If ten sources all cite one original paper, CHISG clusters them as a single opinion node rather than counting them as ten independent confirmations. Source independence is tracked through provenance metadata.

---

## Technical Architecture

### Data Model

```
┌─────────────────────────────────────────────────────┐
│                   Semantic Unit                      │
│                                                      │
│  entity_a:        "kinetic energy"                  │
│  forward_link:    "calculated using"                │
│  entity_b:        "½mv²"                            │
│  backward_link:   "used to calculate"               │
│  context[]:       [source:pearson-physics,           │
│                    domain:physics,                   │
│                    level:gcse,                       │
│                    topic:energy,                     │
│                    exam_board:aqa,                   │
│                    page:142]                         │
│  knowledge_type:  procedural                        │
│  trust_score:     0.87 (computed)                   │
│  vector:          [1536d embedding]                  │
└─────────────────────────────────────────────────────┘
```

### Technology Stack

| Component | Technology | Role |
|-----------|-----------|------|
| Semantic Store | Weaviate | Vector-indexed knowledge graph with `nearVector` semantic search |
| Operational Store | MongoDB | User profiles, extraction metadata, analytics |
| Relational Store | MySQL | Skills assignments, course structures, markbook scores |
| API | Go | Backend orchestration, embedding generation, trust computation |
| Embeddings | AWS Titan | 1536d vectors for semantic similarity |
| LLM | AWS Bedrock (Claude) | Response generation grounded in graph context |
| Frontend | React/TypeScript | Annotation tool, explorer, coaching interface |

### Weaviate Schema (Current)

| Class | Count | Purpose |
|-------|-------|---------|
| CHISGElement | 579 | Skills with domain, description, ETP mappings, vectors |
| SkillLink | 1,120 | Parent→offspring skill relationships with provenance |
| CourseSkillSuggestions | 27 | Curriculum course definitions with skill lists |
| SemanticLink | — | Extracted knowledge claims with full provenance |
| Documentation | 3,764 | Indexed project documentation |

### Knowledge Types

| Type | Example | Pedagogical Handling |
|------|---------|---------------------|
| `declarative` | "Electrons have negative charge" | State fact, verify recall |
| `procedural` | "Rearranging equations" | Sequence steps, practise |
| `conditional` | "Use F=ma when calculating resultant force" | Teach application context |
| `misconception` | "Heavier objects fall faster" | Surface error, correct explicitly |
| `confusable` | "Speed vs Velocity" | Contrast, distinguish features |
| `model` | "Particle model of matter" | Explain scope and limitations |

---

## What Makes CHISG Different

### vs. Traditional Knowledge Graphs (Wikidata, ConceptNet)

Traditional knowledge graphs store facts. CHISG stores **claims with provenance and computed trust**. The difference:

- Wikidata: `(Earth, instance-of, planet)` — a fact, presented as ground truth
- CHISG: `(Earth, instance-of, planet)` — a claim, supported by N sources, via M independent paths, with trust score T, at knowledge level L

CHISG treats knowledge as a hypothesis to be strengthened or weakened by evidence, not a binary true/false flag.

### vs. RAG (Retrieval-Augmented Generation)

RAG retrieves text chunks and feeds them to an LLM. CHISG retrieves **structured relationships** and feeds verified claims to an LLM. The difference:

- RAG: "Here's a paragraph about kinetic energy. Generate a response."
- CHISG: "Here are 12 verified relationships involving kinetic energy, their trust scores, and two structural analogies from other domains. Generate a response that stays within these verified claims and flags anything you're uncertain about."

The LLM knows what it knows and what it doesn't, because the graph explicitly represents both.

### vs. Fine-Tuning

Fine-tuning bakes knowledge into model weights, making it opaque and non-updatable without retraining. CHISG keeps knowledge external, queryable, and updatable. When a curriculum changes, update the graph — don't retrain the model.

---

## The Bootstrap Problem and Solution

### The Problem

Building a comprehensive knowledge graph requires massive annotation effort. Academic sourcing for every claim is impractical at scale.

### The Three-Phase Solution

**Phase 1: Expert Bootstrap (Current)**  
Domain experts capture knowledge claims based on professional expertise. These are tagged as `source:expert-heuristic` and treated as foundational hypotheses.

**Phase 2: Empirical Validation**  
As the system is used in classrooms, observable outcomes validate or challenge claims:
- If teaching skill X (as defined) leads to measurable improvement in outcome Y, the definition is empirically validated
- Trust scores are automatically up-ranked for claims that survive contact with reality

**Phase 3: Academic Grounding**  
Retroactively map published research onto validated nodes, completing the provenance chain. The graph now has expert origin, empirical validation, and academic grounding — three independent legs of trust.

---

## Current State

| Component | Status |
|-----------|--------|
| Weaviate schema with CHISG classes | ✅ Production (Vultr) |
| 579 skills + 1,120 links + 27 courses | ✅ Populated and queryable |
| Semantic search via nearVector | ✅ Working (Go API + AWS Titan embeddings) |
| Annotation tool (semantic link extraction) | 🔄 70% complete (React + Go) |
| Trust scoring computation | ⬜ Designed, not implemented |
| Structural analogy detection | ⬜ Designed, not implemented |
| Gap analysis queries | ⬜ Designed, not implemented |
| Echo chamber filter | ⬜ Designed, not implemented |

---

## Applications

### Education (Primary Focus)
- AI tutoring grounded in curriculum-verified knowledge
- Adaptive learning pathways using prerequisite relationships
- Cross-topic revision using structural analogies
- Misconception detection and targeted correction

### Research
- Curriculum coherence analysis (quantify how well curricula build conceptual connections)
- Cross-curriculum comparison (compare semantic structures across publishers or exam boards)
- Cross-domain pattern discovery (surface structural parallels between fields)

### Content Publishing
- Quality metrics for educational materials (conceptual density, relationship coverage)
- Gap identification in course content
- Localisation validation (ensure translated content preserves critical semantic relationships)

---

## Why This Matters

The dominant approach to educational AI is to fine-tune language models on curriculum content and hope they don't hallucinate. CHISG takes the opposite approach: **extract the knowledge structure first, then use it to constrain and ground AI responses**.

This is not an incremental improvement to RAG. It is a different architecture:

- Knowledge is structured, not chunked text
- Truth is computed from graph structure, not assumed from sources
- Analogies are detected automatically, not authored manually
- Gaps are visible — the system knows what it doesn't know

The methodology is domain-agnostic. The initial application is UK GCSE Science. The architecture extends to any field where structured, trustworthy, interconnected knowledge matters.

---

**Michael Stewart, PhD**  
michael.stewart@espthinking.co.uk | 07525 359661  
https://espthinking.co.uk

*Live system demonstration, repository access, and technical discussion available on request.*
