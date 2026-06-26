# CHISG: Technical Specification

**Contextualised Hierarchical Iterative Semantic Groupings**

**Author**: Michael Stewart, PhD (Biophysics)  
**Date**: February 2026 (updated June 2026)  
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

**Current Implementation**: 579 skills, 1,120 semantic links, 27 course definitions stored in Weaviate (vector database). The Go backend generates embeddings and queries via `nearVector` for semantic matching. GCSE Science content (AQA specification) is the initial domain. The live extraction pipeline (see §5) has extended the graph into the immunology domain via the McGrath mycobacteriology corpus.

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
| PDF upload → Textract → Claude extraction pipeline | ✅ Live (June 2026) |
| Streaming review queue (MongoDB) | ✅ Live — chunks inserted immediately as extracted |
| Human review + approve/edit interface | ✅ Live — `/api/chisg/review/tasks` |
| Weaviate write-back on approval | ✅ Live — `AcademicLink` class populated on approve |
| Extraction corrections log (training signal) | ✅ Live — `extraction_corrections` collection |
| Papers & links query API | ✅ Live — `/api/chisg/papers`, `/api/chisg/links` |
| Annotation tool (semantic link extraction) | ✅ Superseded by live pipeline above |
| Trust scoring computation | ⬜ Designed, not implemented |
| Structural analogy detection | ⬜ Designed, not implemented |
| Gap analysis queries | ⬜ Designed, not implemented |
| Echo chamber filter | ⬜ Designed, not implemented |

---

## §5 Live Extraction Pipeline (Implemented June 2026)

The extraction pipeline described in the original spec as a "next phase" is now live in production. This section documents the actual implementation.

### 5.1 Architecture Overview

```
PDF / Text file
      ↓
  UploadCHISGDocumentHandler  (Go — esp-organizer)
      ↓ (async background goroutine)
  AWS Textract  (eu-west-2, bucket: esp-new-organizer-immunology)
      ↓
  cleanPDFText() + splitIntoParagraphs()  (~1200 chars/chunk at sentence boundary)
      ↓
  Reference section detector  (skips bibliography chunks)
      ↓
  Claude 3 (AWS Bedrock)  — structured extraction prompt
      ↓
  MongoDB  chisg_knowledge_base.review_tasks  (streaming — inserted per chunk)
      ↓  (reviewer approves/edits via UI)
  ApproveReviewTaskHandler
      ↓
  Weaviate  AcademicLink class  (approved links only)
      ↓
  MongoDB  extraction_corrections  (diff log: AI-proposed vs human-approved)
```

The upload endpoint returns immediately with a `job_id`. Chunks are inserted to MongoDB as they complete, so a reviewer can begin working on the first chunks of a long paper while the remainder is still being processed.

### 5.2 Controlled Relation Vocabulary

The 13-relation vocabulary enforced in the Claude extraction prompt:

| Relation | Meaning |
|----------|---------|
| `causes` | A directly produces B |
| `leads_to` | A sets conditions that result in B |
| `part_of` | A is a structural component of B |
| `contains` | A holds or includes B |
| `develops_into` | A matures or transforms into B |
| `regulates` | A controls the activity or expression of B |
| `enables` | A makes B possible |
| `inhibits` | A suppresses or prevents B |
| `treats` | A is a therapeutic intervention for B |
| `diagnoses` | A is used to identify B |
| `manifests_as` | A presents clinically or phenotypically as B |
| `is_a` | A is a subtype or instance of B |
| `located_at` | A exists at or within B |

### 5.3 Semantic Unit (Live Schema)

Each extracted and approved link stored in Weaviate `AcademicLink`:

```
{
  entity_a:           "CarD",
  relation:           "regulates",
  entity_b:           "rRNA transcription",
  context:            "under starvation conditions in Mycobacterium tuberculosis",
  statement:          "CarD directly contacts the β-subunit of RNAP...",  // verbatim source quote
  attribution:        "Stallings et al., 2009",                           // omitted if original paper's own claim
  paper_id:           "elife-73347-v2_1",
  chunk_id:           "chunk_12",
  source_document_id: "64a3f..."
}
```

The `attribution` field is the key provenance distinction: it separates claims this paper's authors make directly (high structural integrity — original data) from claims they cite from prior literature (lower weight — may be echo-chamber repetition).

### 5.4 Training Signal

Every time a reviewer edits an AI-proposed link before approving it, the diff is stored in `extraction_corrections`. This collection is the ground-truth training set for Level 2 (algorithmic extraction), capturing:
- Which entity names the model got wrong
- Which relation types the model misclassified
- Which context strings were too vague or too verbose

### 5.5 Papers & Links Query API

- `GET /api/chisg/papers` — lists all papers in Weaviate with link and entity counts
- `GET /api/chisg/links?paper_id=xxx&limit=N` — returns semantic links for a paper
- `GET /api/chisg/review/tasks` — returns up to 10 pending review tasks
- `POST /api/chisg/review/tasks/{id}/approve` — submit approved/edited links
- `DELETE /api/chisg/review/tasks` — clear the review queue

### 5.6 First Domain: McGrath Immunology Corpus

The live pipeline was validated on the McGrath mycobacteriology eLife paper corpus. The offline extraction (`extraction_elife_sonnet.json`, 6,415 links from 2,746 definitions) served as the initial test of extraction quality. The live pipeline supersedes this workflow — all future domain ingestion uses PDF upload → review → approve.

---

## §6 (Former §5): LLM-to-CHISG Extraction Pipeline — Design Notes

### The Goal

Build a learning model that enables an LLM to convert source material into CHISG semantic units — extracting structured `[Entity A] —[relation]→ [Entity B]` triples with full provenance metadata.

### Difficulty Gradient

The extraction problem has a clear difficulty gradient based on source material density:

#### Level 1: Definitions (Easy)

**Source**: The existing LAO keyword database (2,746 GCSE Science definitions from Quizlet + AI-generated, manually verified).

Definitions are the ideal training candidate because the information is already condensed. A definition like *"Generator: A device that converts kinetic energy into electrical energy"* maps almost directly to a semantic unit:

```
Generator —[converts]→ kinetic energy → electrical energy
  knowledge_type: declarative
  source: quizlet-expert-verified
  domain: physics
  level: gcse
```

The relation type, entities, and domain are largely explicit in the text. Extraction at this level is a parsing problem, not a reasoning problem.

#### Level 2: Textbook Passages (Medium)

Longer structured text where relationships are stated but spread across sentences and paragraphs. The LLM must identify which statements are knowledge claims (vs. pedagogical scaffolding, worked examples, or motivational filler) and assign appropriate relation types.

#### Level 3: Research Papers (Hard — The Real Test)

**This is where the methodology must prove itself.**

Research papers present three challenges that definitions do not:

1. **Context-dependent truth**: A relationship extracted from a paper may only hold under specific experimental conditions, in a particular organism, at a certain scale, or within stated assumptions. The metadata — the `context[]` array — becomes the critical component. A claim like *"Protein X inhibits pathway Y"* is meaningless without: species, cell type, concentration range, temperature, and whether this was in vivo or in vitro. **The context in which the link between items is actually true is far more difficult to extract than the link itself.**

2. **Implicit relationships**: Research papers assume domain knowledge. The relationship between two concepts may never be explicitly stated — it is implied by the experimental design or by the juxtaposition of results. The LLM must infer relationships that the author assumed the reader would recognise.

3. **Hedged and qualified claims**: Research language is deliberately cautious — "suggests", "may contribute to", "is consistent with". These qualifications map to trust metadata (confidence weighting), not to the relationship itself. The extraction model must separate the claim from its hedging.

### Training Strategy

The definitions serve as **supervised training data** for the extraction pipeline:

1. **Phase A**: Convert 2,746 definitions → semantic units (high accuracy, known-good output)
2. **Phase B**: Use these as ground truth to train/evaluate LLM extraction prompts
3. **Phase C**: Graduate to textbook passages, using definition-derived units to validate extracted relationships
4. **Phase D**: Attack research papers, where the context extraction problem becomes the primary challenge

Each phase produces training signal for the next. The definitions are not the goal — they are the **calibration set**.

### Key Metadata Fields for Research-Grade Extraction

| Field | Why It Matters | Easy at Definition Level | Hard at Paper Level |
|-------|---------------|------------------------|-------------------|
| `relation_type` | What connects A to B | Usually explicit ("is a", "converts") | Often implicit or hedged |
| `context.conditions` | When is this true | Usually universal at GCSE | Species, concentration, temperature, etc. |
| `context.scope` | At what scale/level | GCSE level, straightforward | May only apply in vitro, or at molecular scale |
| `confidence` | How sure are we | High (verified definition) | Varies — "suggests" vs "demonstrates" |
| `context.assumptions` | What must be true for this to hold | Rarely stated | Critical and often unstated |
| `contradicts` | What claims conflict | Rare at GCSE | Common in active research areas |

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
