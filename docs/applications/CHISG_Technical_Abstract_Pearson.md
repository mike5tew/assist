# CHISG: Contextualised Hierarchical Iterative Semantic Grouping Pipeline

**Technical Abstract**

**Author:** Michael Stewart, PhD (Biophysics)  
**Date:** December 2024

---

## Executive Summary

The CHISG (Contextualised Hierarchical Iterative Semantic Grouping) pipeline addresses a critical challenge in AI-enhanced education: how do we ensure that large language models remain grounded in curriculum-accurate knowledge when supporting learners? This system extracts and structures the semantic relationships embedded in educational content, creating a "ground truth" layer that can anchor AI tutoring systems to verified curriculum knowledge.

Initial application targets UK GCSE science content, with architecture designed for extension across subjects and educational levels.

## The Problem: Semantic Drift in Educational AI

Large language models exhibit semantic drift and hallucination when applied to educational contexts—generating plausible but incorrect explanations, conflating related concepts, or introducing knowledge inappropriate for the learner's level. This is particularly problematic in high-stakes educational settings where accuracy and curriculum alignment are non-negotiable.

The root cause: **LLMs lack structured access to the specific semantic relationships that constitute curriculum knowledge.**

Educational content contains rich conceptual structures critical for learning:
- **Prerequisite dependencies** (concept A must be understood before concept B)
- **Causal mechanisms** (process X causes outcome Y)
- **Taxonomic relationships** (entity A is a type of entity B)
- **Comparative structures** (approach A contrasts with approach B)

These relationships remain implicit in textbook prose. Current approaches rely on expensive expert annotation or generic NLP techniques that lack educational grounding—neither scales to the breadth of content required for comprehensive AI tutoring systems.

## The CHISG Approach

CHISG creates a deterministic grounding layer for educational AI through three integrated components:

### 1. Human-in-the-Loop Training Data Generation

The pipeline centers on a purpose-built annotation tool that enables rapid extraction of **bidirectional semantic links** from curriculum materials. Design decisions are grounded in cognitive science research:

- **Bidirectional relationships**: Every link captures both forward (A→B) and inverse (B→A) directions. This reflects how semantic memory encodes relationships—"photosynthesis *requires* light" activates different retrieval pathways than "light *is required for* photosynthesis." Capturing both supports diverse query patterns in downstream applications.

- **Quality calibration**: Each extracted link receives a quality score (0-100%) with explicit reasoning, based on the Human-OS ETP (Emotional Trigger Points) framework. This enables:
  - Filtering for high-value training examples
  - Research into what makes relationships "valuable" for learning
  - Training signal for models to recognize low-utility extractions

- **Context linking**: Relationships link to related extractions, preserving the interconnected structure of scientific knowledge that research shows is critical for transfer.

- **Source provenance**: Full traceability to source document, page, and excerpt—essential for verification and curriculum alignment auditing.

### 2. Educationally-Grounded Relationship Taxonomy

The system employs a controlled vocabulary of relationship types derived from analysis of exam board specification language and grounded in learning science research on conceptual change:

| Category | Example Relationships | Learning Science Basis |
|----------|----------------------|------------------------|
| **Hierarchical** | is-a-type-of, includes, is-part-of | Taxonomic knowledge structures |
| **Causal** | causes, enables, inhibits, requires | Mechanistic reasoning |
| **Comparative** | contrasts-with, is-similar-to | Discrimination learning |
| **Functional** | is-used-for, performs, transfers | Procedural knowledge |
| **Definitional** | is-defined-as, has-property | Declarative foundations |

This taxonomy enables curriculum-appropriate relationship extraction that aligns with how exam boards frame required knowledge.

### 3. Scalable Technical Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  PDF Ingestion  │────▶│  Human Annotation │────▶│  Training Data  │
│  + Text Extract │     │  Interface        │     │  Export         │
└─────────────────┘     └──────────────────┘     └─────────────────┘
                                                          │
                                                          ▼
                              ┌──────────────────────────────────────┐
                              │  Fine-tuned Extraction Model         │
                              │  (Automated relationship detection)  │
                              └──────────────────────────────────────┘
                                                          │
                                                          ▼
                              ┌──────────────────────────────────────┐
                              │  Curriculum Knowledge Graph          │
                              │  (Vector-indexed semantic layer)     │
                              └──────────────────────────────────────┘
                                                          │
                                                          ▼
                              ┌──────────────────────────────────────┐
                              │  AI Grounding Layer                  │
                              │  (Deterministic retrieval for LLMs)  │
                              └──────────────────────────────────────┘
```

The architecture separates concerns to enable institutional deployment:
- **Backend:** Go with MongoDB (document storage) and Weaviate (vector database for semantic search)
- **Frontend:** React/TypeScript with integrated PDF viewer for efficient annotation
- **Export:** JSON-LD, CSV, and custom formats for model training and integration

### 4. Quality Calibration Framework

A key contribution is the quality scoring rubric, designed to operationalize: *What makes a semantic relationship valuable for learning?*

**High-value indicators:**
- Captures mechanism or causal structure (not mere association)
- Non-obvious to target learner demographic
- Connects concepts across topic boundaries (supports transfer)
- Uses curriculum specification language

**Low-value indicators:**
- Trivial enumeration without conceptual insight
- Common knowledge for target audience
- Redundant with other captured relationships

This rubric represents a testable hypothesis about educational value, enabling empirical research on curriculum quality.

## Applications

### Product Development Support

1. **AI tutoring grounding**: Provide deterministic knowledge retrieval for LLM-based tutoring systems, preventing hallucination and ensuring curriculum alignment
2. **Adaptive learning pathways**: Prerequisite relationships enable intelligent content sequencing based on learner knowledge state
3. **Diagnostic assessment**: Explicit relationship structures support targeted assessment item generation

### Research & Thought Leadership

1. **Curriculum coherence analysis**: Quantify how well curricula build conceptual connections; identify gaps where relationships are assumed but never taught
2. **Cross-curriculum comparison**: Compare semantic structures across publishers, exam boards, or educational systems
3. **Misconception research**: Analyze patterns in how students misrepresent relationships compared to curriculum ground truth

### Business Applications

1. **Content quality metrics**: Objective measurement of conceptual density and relationship quality in educational materials
2. **Competitive analysis**: Map semantic coverage differences between Pearson and competitor content
3. **Localisation validation**: Ensure translated content preserves critical semantic relationships

## Current Status

- **Operational**: Annotation interface deployed with GCSE Physics content (AQA specification)
- **In progress**: Extraction of energy transfer and storage domain (~100+ semantic relationships)
- **Functional**: Export pipeline producing model-ready training data with quality scoring

**Validation roadmap:**
- Inter-annotator reliability study with domain experts
- Alignment verification against exam board specifications
- Extraction model performance benchmarking on held-out content

## Alignment with Pearson R&D Priorities

| Pearson Priority | CHISG Contribution |
|------------------|-------------------|
| **AI in education research** | Solves grounding problem for LLM deployment in high-stakes learning |
| **Learning science application** | Operationalizes conceptual change theory in production systems |
| **Evidence-based insights** | Generates empirical data on curriculum semantic structures |
| **Product development support** | Provides foundation for next-generation adaptive systems |
| **Thought leadership** | Novel framework for defining "educational value" in AI context |

The pipeline translates theoretical learning science—particularly research on knowledge integration, conceptual change, and mechanistic reasoning—into operational infrastructure for AI-enhanced learning products.

---

**Michael Stewart, PhD**  
michael.stewart@espthinking.co.uk | 07525 359661  
LinkedIn: [Profile](https://linkedin.com/in/michaelstewart)

*Technical demonstration and repository access available upon request.*

