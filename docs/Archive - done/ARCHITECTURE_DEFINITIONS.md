# ESP Organizer Architecture Definitions

## 1. Semantic Links (SL) - Tier 1: Atomic Knowledge Units

**Purpose**: Clean, high-quality vector embeddings for precise knowledge retrieval.

**Definition**: A Semantic Link (SL) is a **data object** that captures a **single, atomic relationship** between two knowledge elements (keywords, concepts, definitions, or entities).

**Structure** (Updated - Removed Fixed Hierarchy):
```json
{
  "source_term": "X-Linked Agammaglobulinemia",
  "target_term": "B-cell deficiency",
  "relation_type": "causes",
  
  // 🆕 RELATIONSHIP METADATA (not hierarchy level!)
  "relationship_strength": 0.95,
  "relationship_direction": "source_causes_target",
  "semantic_distance": 0.15, // How "far apart" concepts are semantically
  
  // 🆕 TERM PROPERTIES (for dynamic hierarchy computation)
  "source_term_generality": 0.42, // 0=very specific, 1=very general
  "target_term_generality": 0.58,
  
  // Graph connectivity (still needed)
  "is_parent_of": [], // Broader concepts this links to
  "is_child_of": [],  // More specific concepts this derives from
  
  "context": "XLA results in absence of B cells due to BTK gene mutation",
  "confidence": 0.95,
  "source_mongo_id": "68e2ed2ffcc06550ef8703f8",
  "target_mongo_id": "68e2ed2ffcc06550ef8703f9"
}
```

**Key Characteristics**:
- ✅ **Atomic**: One relationship per link
- ✅ **Vectorized**: Each link gets its own embedding (from AWS Titan)
- ✅ **Bidirectional**: Links can be traversed source→target or target→source
- ✅ **Typed**: Relation types include: `causes`, `treats`, `is_a`, `related_to`, `component_of`
- ✅ **Hierarchical**: Links have `hierarchy_level` (1=root concept, 2=specific condition, 3=clinical finding)

**Hierarchy Levels**:
- **Level 1 (Root)**: Broad domain concepts (e.g., "Immunodeficiency Disorders")
- **Level 2 (Concept)**: Specific conditions/subtypes (e.g., "X-Linked Agammaglobulinemia")
- **Level 3 (Clinical)**: Specific findings/mechanisms (e.g., "BTK gene mutation", "Recurrent bacterial infections")

**Storage**:
- **MongoDB**: Raw text content (`immunology_content`, `immunology_terms`)
- **Weaviate**: Vectorized semantic links (`SemanticLinks` class)

---

## 2. Hierarchical Semantic Graph (HSG) - Tier 2: Narrative Threads

**Purpose**: Provide **contextual understanding** by creating **narrative threads** of connected Semantic Links with metadata.

**Definition**: An HSG is a **multi-hop traversal path** through a network of Semantic Links that tells a **coherent story** about a topic.

**Structure** (Implemented in MongoDB):
```json
{
  "summary_chunk_id": "68e2ed2ffcc06550ef8703fa",
  "summary_text": "XLA is a primary immunodeficiency caused by BTK gene mutations leading to B-cell absence and recurrent bacterial infections, treated with IVIG replacement therapy",
  "domain": "immunology",
  "semantic_link_ids": ["sl_001", "sl_045", "sl_102", "sl_203", "sl_304"],
  "derived_semantic_link_ids": ["sl_abstract_001", "sl_abstract_002"],
  "derived_link_count": 2,
  "confidence": 0.87,
  "metadata": {
    "topic": "causes",
    "link_count": 5,
    "source_document_id": "68e2ed2ffcc06550ef8703f8",
    "vector_source": "aws_bedrock_titan"
  }
}
```

**Key Characteristics**:
- ✅ **Multi-hop**: Chains together multiple SLs (Tier 1 granular links)
- ✅ **Contextual**: Metadata provides narrative understanding
- ✅ **Hierarchical**: Traverses from abstract concepts → specific findings
- ✅ **Summarized**: Includes LLM-generated summaries of the thread
- ✅ **🆕 RECURSIVE**: Extracts NEW higher-level semantic links from summaries

**Query Flow (Complete HSG)**:
1. **Query**: User asks "What causes X-Linked Agammaglobulinemia?"
2. **Tier 2 Search**: Find relevant SummaryChunks (abstract level)
3. **Retrieve Derived Links**: Get higher-level SemanticLinks extracted FROM summaries
4. **Thread Construction**: Traverse to Tier 1 granular SemanticLinks
5. **Content Retrieval**: Get raw DocChunk content from MongoDB
6. **Context Assembly**: Create HSG with full causal chain and metadata
7. **Answer Generation**: LLM uses hierarchical context for comprehensive answer

---

## 🆕 2.1. Tier 2 Recursive Processing

**The Complete Tier 2 Pipeline**:
```json
{
  "pipeline_id": "tier_2_recursive_processing",
  "description": "Complete processing pipeline for Tier 2 Recursive HSG",
  "steps": [
    {
      "step_id": "query_reception",
      "action": "Receive user query",
      "output": "Raw query text"
    },
    {
      "step_id": "initial_sl_search",
      "action": "Search Semantic Links (Tier 1) using vector similarity",
      "input": "Raw query text",
      "output": "Initial set of Semantic Links"
    },
    {
      "step_id": "hsg_thread_construction",
      "action": "Construct initial HSG thread from Semantic Links",
      "input": "Initial set of Semantic Links",
      "output": "Partial HSG thread"
    },
    {
      "step_id": "summary_generation",
      "action": "Generate summary for HSG thread using LLM",
      "input": "Partial HSG thread",
      "output": "Thread summary text"
    },
    {
      "step_id": "derived_link_extraction",
      "action": "Extract higher-level Semantic Links from summary",
      "input": "Thread summary text",
      "output": "Derived Semantic Links"
    },
    {
      "step_id": "final_hsg_construction",
      "action": "Construct final HSG with derived links and metadata",
      "input": "Derived Semantic Links",
      "output": "Complete HSG thread"
    },
    {
      "step_id": "answer_generation",
      "action": "Generate answer using complete HSG thread as context",
      "input": "Complete HSG thread",
      "output": "Comprehensive answer text"
    }
  ],
  "metadata": {
    "version": "1.0",
    "last_updated": "2023-10-10",
    "maintainer": "ESP Organizer Team"
  }
}
```

**Pipeline Steps**:
1. **Query Reception**: The pipeline receives and logs the user query.
2. **Initial SL Search**: Perform a vector similarity search to find relevant Semantic Links (Tier 1).
3. **HSG Thread Construction**: Build an initial HSG thread using the found Semantic Links.
4. **Summary Generation**: Use an LLM to generate a summary of the HSG thread.
5. **Derived Link Extraction**: Extract new higher-level Semantic Links from the summary.
6. **Final HSG Construction**: Construct the final HSG with the derived links and additional metadata.
7. **Answer Generation**: Generate a comprehensive answer using the complete HSG thread as context.

---

## 3. Skills Map - Specialized HSG

**Purpose**: Model skill development hierarchies with **component skills** (prerequisites) and **offspring skills** (next steps).

**Definition**: A Skill is a node in an HSG where:
- **Parent Skills** (components): Skills required to develop this skill
- **Offspring Skills**: Skills this skill enables you to develop
- **Assessment Criteria**: How to evaluate mastery
- **Practice Tasks**: Activities that develop the skill

**Structure**:
```json
{
  "skill_id": "handwriting",
  "name": "Handwriting",
  "description": "Ability to write legibly by hand",
  "parent_skills": [
    {"id": "fine_motor_skills", "name": "Fine Motor Skills"},
    {"id": "letter_recognition", "name": "Letter Recognition"}
  ],
  "offspring_skills": [
    {"id": "cursive_writing", "name": "Cursive Writing"},
    {"id": "note_taking", "name": "Note Taking"}
  ],
  "assessment_criteria": [
    "Letters are formed correctly",
    "Writing is legible to others",
    "Maintains consistent spacing"
  ],
  "practice_tasks": [
    {"id": "task_001", "title": "Trace letters"},
    {"id": "task_002", "title": "Copy sentences"}
  ]
}
```

**Key Difference from General HSG**:
- Skills have **explicit parent/child relationships** (not just semantic links)
- Skills have **assessment criteria** (measurable outcomes)
- Skills have **practice tasks** (actionable learning activities)
- The graph represents **learning pathways**, not just knowledge relationships

---

## 4. Current Implementation Status

### ✅ What Works (Tier 1):
1. **Semantic Link Extraction**: LLaMA extracts relationships from text
2. **Vectorization**: AWS Titan creates embeddings for links
3. **Storage**: Links stored in Weaviate `SemanticLinks` class
4. **Search**: Vector search finds relevant links

### ❌ What's Missing (Tier 2):
1. **`hierarchy_level` field**: Not being set during link creation
2. **Parent/child arrays**: `is_parent_of` and `is_child_of` are empty
3. **Thread construction**: No code to build multi-hop narrative chains
4. **Summary generation**: No LLM-powered thread summarization
5. **HSG query service**: Exists but doesn't do hierarchical traversal

### 🎯 Immediate Next Steps:
1. **Update semantic link creation** to set `hierarchy_level`
2. **Implement link classification** (root vs. concept vs. clinical)
3. **Add parent/child relationship tracking**
4. **Build multi-hop graph traversal** in `HSGQueryService`
5. **Generate thread summaries** using Claude

---

## 5. Data Flow

### Current (Flat Semantic Links):
````
