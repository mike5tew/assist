# CHISG System Manifest: Vultr Semantic Overlay

## Overview
This document describes the technical implementation of the CHISG (Knowledge Claim) logic as it sits atop the Vultr production stack. It bridges the gap between Relational Data (MySQL) and Semantic Logic (Weaviate).

## 1. The Multi-Model Data Stack

| Layer | Technology | Role | "The Persona" |
| :--- | :--- | :--- | :--- |
| **Relational Layer** | MySQL | Raw attainment, student logs, user accounts, skill definitions (62 core skills). | **Safe Hands** |
| **Semantic Layer** | Weaviate | Vectorized semantic links, analogy detection, and structural relationship mapping. | **Smart Minds** |
| **Reasoning Layer** | CHISG Logic | The rules engine that validates claims using Triangulation, Recursive Mastery, and the Echo Chamber filter. | **The Architect** |

## 2. CHISG Implementation on Weaviate

### Class Design: `CHISG_SemanticLink`
The core atom of knowledge in Weaviate.

*   **Properties**:
    *   `entityA` (String): The source node (e.g., "Working Memory").
    *   `relation` (String): The forward link (e.g., "masters", "component_of").
    *   `entityB` (String): The target node (e.g., "Impulsivity").
    *   `trustScore` (Number): Computed value based on Source and Path Density.
    *   `sourcePedigree` (String[]): Array of independent sources or personal heuristics (Expert Bootstrap).
    *   `vector`: 1536d embedding (AWS Titan) for analogy matching.

### Class Design: `HumanOS_Tell`
The capture class for student behavioral patterns.

*   **Properties**:
    *   `tellName` (String): e.g., "Micro-withdrawal".
    *   `contextId` (String): UUID mapping to a MySQL session.
    *   `voltageState` (Number): Real-time challenge level.
    *   `socialGravityWeight` (Number): Influence of peer group.

## 3. The Recursive Loop (The "Faking It" Engine)

The system operates on a "Bootstrap" logic to generate valid insights from unsourced data:

1.  **Observed Behavior**: Teacher logs a "Tell" via Passive Pilot UI.
2.  **MySQL Entry**: Transaction recorded in `skillsmarkbook`.
3.  **CHISG Inference**:
    *   Query Weaviate for: `MATCH (Tell) -[:evidences]-> (Skill) -[:masters]-> (ETP)`.
    *   If `Skill Improvement` is later recorded in MySQL, **Up-rank** the `trustScore` of the links in Weaviate.
    *   This converts "Expert Heuristics" into "Empirical Evidence" automatically.

## 4. Current State (Vultr)
*   **MySQL**: Active. Contains the 62 core skills and criteria.
*   **Weaviate**: Active. Ready to receive the vectorized `SemanticLink` schema.
*   **CHISG Logic**: Resides in the Go backend (`esp-organizer`), acting as the coordinator between the two databases.
