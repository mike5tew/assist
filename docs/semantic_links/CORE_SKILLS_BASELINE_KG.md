# CHISG: HumanOS Skill Taxonomy (The Meta-Circuits)

This document organizes the **115+ Skills** from the living database (Local + Vultr) into four Meta-Buckets: **Leadership**, **The 7 Intelligences**, **Practical Skills**, and **Primary Precursors**. This layout ensures "Hallucination Protection" and "Thorough Coverage" for a student's entire developmental journey.

---

## 1. The Leadership Set
These high-level skills govern system direction and strategic outcome.

```json
[
  ["Strategic Foresight (1)", "requires", "Critical Thinking (14/4)", "identifies paths for", 
   ["domain:leadership", "circuit:strategic"]],
  ["Risk Intelligence (2)", "guides", "Decision Making", "calculates friction for", 
   ["domain:leadership", "circuit:strategic"]],
  ["Integrative Leadership (8/17)", "combines", "Interpersonal Intelligence (5)", "aligns goals for", 
   ["domain:leadership", "circuit:governance"]],
  ["Adaptive Planning (7)", "requires", "Contextualization (69/59)", "adjusts course for", 
   ["domain:leadership", "circuit:strategic"]],
  ["Goal Setting (Personal) (111)", "utilizes", "Self-Motivation (43/33)", "focuses energy on", 
   ["domain:leadership", "circuit:agency"]],
  ["Self-Advocacy (105)", "powered by", "Social Confidence (64/54)", "enables output of", 
   ["domain:leadership", "circuit:agency"]]
]
```

---

## 2. The 7 Intelligences Circuit
The core processing modules of the HumanOS, mapped to Gardner's framework.

### I. Linguistic (Communication High-Output)
- **Core Skills**: Vocabulary (18/8), Speaking (22/12), Reading (23/13), Phonics (24/14), Grammar (25/15), Spelling (27/17), Handwriting (28/18).
- **Advanced Skills**: Persuading (29/19), Questioning (30/20), Debating (31/21), Poetry (32), Storytelling (35), Public Speaking (62/52).

### II. Logical-Mathematical (Relational Logic)
- **Core Skills**: Critical Thinking (14/4), Problem Solving (15/5), Application of Ideas (48/38), Categorization (49/39).
- **Data Skills**: Data Handling (50/40), Data Contextualization (51/41), Understanding Variables (65/55), Graph - Reading (45/35).

### III. Spatial (Visual Orientation)
- **Core Skills**: Graph - Drawing (44/34), Graph - Reading (45/35), Visual Deconstruction (70/60).

### IV. Bodily-Kinesthetic (Physical Agency)
- **Core Skills**: Fine Motor Skills (56/46), Equipment Handling (53/43), Lab Safety (55/45).

### V. Musical / Audio (Pattern Recognition)
- **Core Skills**: Audio Processing (13/3), Rhythm Recognition (Pending), Phonetic Resonance (Linked to Phonics 24).

### VI. Interpersonal (Social Gravity)
- **Core Skills**: Empathy (42/32), Social Confidence (64/54), Managing Personal Conflict (104).
- **Global Skills**: Active Citizenship (107), Recognising Diversity (106).

### VII. Intrapersonal (Self-Governor)
- **Core Skills**: Self-awareness (30), Self-regulation (41/31), Self-Motivation (43/33).
- **Living Layer**: Zen Warrior Mindset (113), Coping Strategies (110).

---

## 3. The Technical & Practical Set
The "Operational Equipment" for interacting with the modern world. This bucket covers technical proficiency and resource management.

### I. IT & Digital Proficiency
```json
[
  ["Digital Literacy (126)", "requires", "Symbol Recognition (122)", "interfaces with", 
   ["domain:technical", "circuit:it"]],
  ["Basic Digital Safety (102)", "requires", "Critical Thinking (14/4)", "protects against", 
   ["domain:technical", "circuit:cyber-security"]],
  ["Keyboarding/Input Mastery (127)", "utilizes", "Fine Motor Skills (56/46)", "externalizes thought via", 
   ["domain:technical", "circuit:it"]],
  ["Information Searching (128)", "utilizes", "Categorization (49/39)", "retrieves data from", 
   ["domain:technical", "circuit:it"]]
]
```

### II. Resource & Money Handling
```json
[
  ["Basic Money Handling (129)", "requires", "Number Confidence (Precursor)", "manages physical", 
   ["domain:practical", "circuit:financial"]],
  ["Financial Awareness (115/101)", "requires", "Logical Reasoning", "calculates value for", 
   ["domain:practical", "circuit:financial-strategy"]],
  ["Resource Estimation (130)", "utilizes", "Critical Thinking (14/4)", "allocates supply for", 
   ["domain:practical", "circuit:operational"]]
]
```

### III. Operational Life Skills
```json
[
  ["Time Management (114)", "requires", "Working Memory (11)", "externalizes clock for", 
   ["domain:practical", "circuit:operational"]],
  ["Healthy Lifestyle Choices (103)", "leverages", "Self-regulation (41/31)", "enables long-term", 
   ["domain:practical", "circuit:biological-integrity"]],
  ["Following Protocols (57/47)", "requires", "Attention (12/2)", "ensures accuracy in", 
   ["domain:practical", "circuit:operational"]]
]
```

---

## 4. Primary Developmental Precursors (The DRB Foundation)
These skills (IDs 116-125) are the precursors for younger students (Primary/EYFS) before they can access the Core 62.

```json
[
  ["Gross Motor Skills (116)", "underpins", "Fine Motor Skills (56/46)", "provides stability for", 
   ["domain:developmental", "circuit:physical"]],
  ["Phonemic Awareness (117)", "precursor to", "Phonics (24/14)", "identifies sounds for", 
   ["domain:developmental", "circuit:literacy"]],
  ["Instruction Retention (2-Step) (118)", "precursor to", "Attention (12/2)", "buffers command for", 
   ["domain:developmental", "circuit:executive-function"]],
  ["Turn-Taking & Sharing (119)", "precursor to", "Interpersonal Intelligence (VI)", "governs rhythm of", 
   ["domain:developmental", "circuit:social"]],
  ["Emotional Naming (120)", "precursor to", "Self-regulation (41/31)", "labels voltage for", 
   ["domain:developmental", "circuit:emotional"]],
  ["Transition Management (121)", "precursor to", "Adaptive Planning (7)", "manages state-change for", 
   ["domain:developmental", "circuit:executive-function"]],
  ["Symbol Recognition (122)", "underpins", "Reading (23/13)", "decodes abstraction for", 
   ["domain:developmental", "circuit:literacy"]],
  ["Narrative Sequencing (123)", "precursor to", "Problem Solving (15/5)", "orders events for", 
   ["domain:developmental", "circuit:logical"]],
  ["Personal Belonging Care (124)", "precursor to", "Organizing Information (61/51)", "reduces entropy of", 
   ["domain:developmental", "circuit:practical"]],
  ["Self-Toileting & Hygiene (125)", "underpins", "Healthy Lifestyle (103)", "ensures biological", 
   ["domain:developmental", "circuit:practical"]],
  ["Number Correspondence (131)", "precursor to", "Basic Money Handling (129)", "counts units for", 
   ["domain:developmental", "circuit:logical"]]
]
```

---

## 5. Hallucination Protection & Audit Layer
To prevent the expansion of the skill map into redundant or synonymous definitions (e.g., adding "Resilience" when "Self-regulation" exists), the following rules apply:

1. **Circuit Saturation**: A new skill must either define a unique circuit path or be identified as a leaf node (alias) of an existing node.
2. **Gardner Alignment**: Any skill under "7 Intelligences" must be auditing one of the 7 core units.
3. **Leadership Convergence**: Leadership skills are composite skills—the orchestration of multiple intelligences towards a strategic goal.
