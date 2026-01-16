# HumanOS: Spectrum-to-Skill Mapping (The Meta-Buckets)

This document defines the relationship between the **15 Sensitivity Spectrums** (the "Passive Pilot" UI) and the **115+ Skills** (the "Safe Hands" Markbook), organized by the **Meta-Bucket Taxonomy**.

## Meta-Bucket Architecture
The system organizes all skills into three core domains to ensure clarity and coverage:
1.  **Leadership Skills**: Composite skills for orchestration and agency.
2.  **The 7 Intelligences**: Core biological/cognitive processing units (Gardner model).
3.  **Practical Skills**: Operational equipment for daily success.

---

## Mapping Table (Spectrum to Meta-Bucket Anchor)

| HumanOS Spectrum | Primary Skill Anchor | Meta-Bucket | Logic / Circuit |
| :--- | :--- | :--- | :--- |
| **Shy — Outgoing** | Social Confidence (64) | 7-Int (VI) | Social Gravity calibration |
| **Guilt-Aware — Remorseless** | Emotional Naming (120) | Precursor | Labeling arousal (Pre-Reg) |
| **Private — Expressive** | Speaking (22) | 7-Int (I) | Linguistic output volume |
| **Aggressive — Passive** | Turn-Taking (119) | Precursor | Social force mechanics |
| **Empathic — Detached** | Empathy (42) | 7-Int (VI) | Relational logic |
| **Generous — Self-Interested** | Leadership (17) | Leadership | Altruistic resource logic |
| **Bored — Enthusiastic** | Self-Motivation (43) | 7-Int (VII) | Internal Voltage generation |
| **Patient — Impatient** | Time Management (114) | Practical | Externalizing the clock |
| **Fragile — Resilient** | Zen Warrior (113) | 7-Int (VII) | Stress dampening |
| **Risk-Averse — Risk-Seeking** | Risk Intelligence (2) | Leadership | Path friction calculation |
| **Independent — Dependent** | Self-Advocacy (105) | Leadership | Agency activation |
| **Conformist — Rebellious** | Debating (31) | 7-Int (I) | Constructive pushback |
| **Literal — Adaptable** | Transition Mgmt (121) | Precursor | State-change mechanicals |
| **Modest — Status-Driven** | Goal Setting (111) | Leadership | Kinetic intent redirection |
| **Truthful — Strategic** | Strategic Foresight (1) | Leadership | Information path simulation |

---

## Technical Dampeners (Independence Amplification)
While behavioral spectrums are primarily moderated by "7 Intelligences" or "Leadership," high mastery in **Technical & Operational skills** acts as a force-multiplier for the **Independent — Dependent** spectrum.

| Technical Skill | Target Spectrum | Independence Logic |
| :--- | :--- | :--- |
| **Digital Literacy (126)** | Independent — Dependent | Ability to resolve information gaps via technology reduces hovering. |
| **Money Handling (129)** | Independent — Dependent | Mechanical ability to participate in the "System Economy" (Lunch, fees). |
| **Time Management (114)** | Patient — Impatient | Mastering the external clock reduces anxiety and wait-time friction. |

---

## Skill Coverage & Integrity (The Hallucination Guard)

To prevent semantic bloat and AI hallucinations when adding new skills, the **CHISG Integrity Layer** performs a "Semantic Audit":

| Audit Rule | Purpose | CHISG Action |
| :--- | :--- | :--- |
| **Circuit Saturation** | Prevent Redundancy | Flag new Skill if it maps to the exact same [Meta-Bucket + Intelligence Type] as an existing Skill. |
| **Intelligence Mapping** | Ensure Grounding | Any skill in the "7-Int" bucket must audit one of the 7 core units. If it cannot, it is moved to "Practical." |
| **Leadership Composition** | Prevent "Ghost" skills | Leadership skills must be defined as a **composite** of 2+ Intelligence units. |

### The "62 Core" Anchor
The 62 original skills represent the "Ground Truth." Any new skill (63+) from the Vultr living layer must prove its uniqueness within this 3-bucket taxonomy.

---

## Vultr Extended Coverage (Living Skills)
The community-added skills from the Vultr platform have started to fill critical gaps in the HumanOS logic:

| Sensitivity Spectrum | Extended Skill from Vultr | Impact Logic |
| :--- | :--- | :--- |
| **Fragile — Resilient** | 110 (Coping Strategies) | Provides the "Active Recovery" equipment for the Fragile end of the spectrum. |
| **Bored — Enthusiastic** | 113 (Zen Warrior Mindset) | Teaches "Arousal Control" (Voltage management) for high-amplitude students. |
| **Aggressive — Passive** | 104 (Managing Personal Conflict) | Provides practical alternatives to externalized energy (Aggression). |
| **Independent — Dependent** | 112 (Identifying Support Needs) | Dependency is re-framed as a skill of "Targeted Support Identification." |

---

## Data Loop (Recursive Mastery)
1.  **Observation**: Teacher logs a "Tell" on the **Passive Pilot UI** (e.g., student was "Highly Strategic" in a group task).
2.  **Profile Update**: The **Strategic** end of the Spectrum (15) is updated in the student's Profile.
3.  **Skill Suggestion**: The system flags **Critical Thinking (14)** and **Problem Solving (15)** for the next "Safe Hands" assessment.
4.  **Verification**: If the student scores high in **Critical Thinking**, the "Strategic" behavior is re-categorized as "Intelligent Resource Management" rather than "Deception."

## Technical Implementation
*   **MySQL**: `skillsmarkbook` stores the raw scores.
*   **Profile Engine**: A JSON object in the `students` table (or a separate `etp_profiles` table) stores the spectrum values (-5 to +5).
*   **Logical Integrity Layer (CHISG)**: All mappings are validated against [HUMAN_OS_INTEGRITY_LINKS.md](semantic_links/HUMAN_OS_INTEGRITY_LINKS.md). This prevents "hallucinated skills" by requiring each link to have a triangulation density (Expert + Science sources).
*   **Mapping Layer**: A lookup table or middleware logic maps the 15 spectrum IDs to 62 skill IDs for recursive weightings.
