# ESP/HumanOS Ecosystem — Project Index

> **Purpose**: Single source of truth for understanding this interconnected educational technology ecosystem.  
> **Audience**: Collaborators, employers, AI assistants, and future-me.

---

## Executive Summary

This ecosystem comprises **foundational tools** for education data management and AI-assisted learning, and **products** built on top of them. The core innovation is **CHISG** (Contextualised Hierarchical Iterative Semantic Groupings) — a knowledge-graph methodology for structuring educational content — combined with **ETP** (Emergent Tendency Profiles) for personalised learning.

### Architecture Layers

1. **Foundational Tools** (reusable services, not standalone products):
   - CHISG knowledge graph (Weaviate + semantic links)
   - Skills Map Platform (skills CRUD, course maps, student tracking)
   - humanOS (AI coaching orchestrator with safeguarding)
   - ETP framework (9 biological spectra + voltage calculations + response matrix)

2. **Products** (user-facing applications):
   - LAO (GCSE revision mobile app — in TestFlight)
   - PrimaryOS (primary-age ETP — concept stage)
   - ParentOS (parent-facing spectrum tools — concept stage)
   - CareerOS (skills passport — concept stage)
   - ETP Profile (interactive spectrum demo — working)

3. **Concept Demonstrations** (stakeholder engagement):
   - DRB (MAT strategic dashboard — demo with seeded data, shows MATs what's possible)

### Primary Value Propositions

1. **For Schools**: Real-time competency tracking, behaviour logging, and strategic oversight
2. **For Students**: Adaptive revision tools with engagement-first design
3. **For Parents**: Understanding and supporting their child's development across the spectra
4. **For Trusts/MATs**: Consolidated financial and estates monitoring with data-driven governance

---

## Generated: 2026-02-15

---

## Project Overview

### Foundational Tools

| Project | Purpose | Status | Tech Stack |
|---------|---------|--------|------------|
| **CHISG** | Knowledge graph methodology for skills/concepts | Weaviate schema working, semantic search functional | Weaviate, Go, MongoDB |
| **skills-map-platform** | Skills tree, course maps, student proficiency | Production-ready | Go, React, MySQL, Weaviate |
| **humanOS** | AI coaching orchestrator with safeguarding | Partial (in-memory storage) | Go, React, Weaviate, AWS Bedrock |
| **ETP Framework** | 9 biological spectra, voltage calculations, response matrix | Domain logic + response matrix complete, MongoDB persistence | Go, MongoDB |

### Products

| Project | Purpose | Status | Tech Stack |
|---------|---------|--------|------------|
| **LAOMobile** | GCSE revision app | TestFlight (field testing) | React Native, Go, SQLite |
| **PrimaryOS** | Primary-age ETP + Neuron Navigators | Landing page only (concept) | React (landing) |
| **ParentOS** | Parent-facing spectrum explorer | Landing page only (concept) | React (landing) |
| **CareerOS** | Skills passport + role matching | Landing page only (concept) | React (landing) |
| **ETP Profile** | Interactive spectrum demo | Working | React, Go |
| **SkillsMarkbookMobile** | Primary teacher competency logging | Development (early) | React Native, Expo |

### Concept Demonstrations

| Project | Purpose | Status | Tech Stack |
|---------|---------|--------|------------|
| **DRB** | MAT financial & estates dashboard concept | Demo-functional (seeded data) | Go, React, MongoDB |

### Infrastructure

| Project | Purpose | Status | Tech Stack |
|---------|---------|--------|------------|
| **assist** | Monorepo (docs, API, frontend, deployment) | Active | Go, React, Weaviate, Docker |
| **ESPProj** | Legacy MIS features | Migration | Go, React, MySQL |

---

## 🎯 Priority Projects (Job-Relevant)

### 1. SkillsMarkbookMobile
**Location**: `/Users/michaelstewart/Coding/SkillsMarkbookMobile`

A mobile app for primary school teachers to log student competency and character in real-time.

**Core Modes**:
- **Action Recorder**: Spatial seating plan with 4-button behaviour logging (2+/2-)
- **Profile Builder**: 8 ETP sensitivity spectrums as sliders with justification notes
- **Skills Markbook**: Proficiency tracking (1-5) against 131 CHISG Core Skills

**Key Files**:
```
SkillsMarkbookMobile/
├── app/
│   └── index.tsx              # Main entry point
├── prototypes/
│   └── PassivePilotUI.web.tsx # Sticker album prototype
├── package.json               # Dependencies
└── README.md                  # Project overview
```

**Integration Points**:
- Pushes data to ESP Platform API
- Uses CHISG Knowledge Graph for skill definitions
- ETP profiles sync to Weaviate

---

### 2. DRB (MAT Strategic Dashboard — Concept Demo)
**Location**: `/Users/michaelstewart/Coding/DRB`

A concept demonstration for Multi-Academy Trusts, showing the type of strategic oversight tooling that could be built. Not a standalone product — used for stakeholder engagement to illustrate what's possible with consolidated financial and estates data.

**Key Modules**:
- **Bank Aggregation**: Open Banking (AISP) integration
- **Estates Health**: Utility spend correlation with physical assets
- **Governance Radar**: Traffic-light monitoring for governors

**Key Files**:
```
DRB/
├── main.go                    # API entry point
├── Data/                      # Data processing logic
├── frontend/                  # React dashboard
├── docs/                      # Architecture docs
├── docker-compose.yml         # Local dev environment
└── README.md                  # Project overview
```

**Value for Academy Chains**:
- Consolidated cash position across schools
- Proactive infrastructure intervention
- Non-technical stakeholder dashboards

---

## Core Concepts

### CHISG (Contextualised Hierarchical Iterative Semantic Groupings)
A method of parsing and validating information for vector database ingestion that:
- Maps semantic relationships between skills and concepts
- Enables gap analysis in learning pathways
- Reduces AI hallucinations through grounded knowledge

**Location**: `assist/docs/chisg/CHISG_KNOWLEDGE_GRAPH_ARCHITECTURE.md`

### ETP (Emergent Tendency Profiles)
9 biological spectrums that describe learner tendencies:

| Spectrum | Description |
|----------|-------------|
| social_gravity | Introvert ↔ Extrovert tendency |
| energy_directionality | Task focus ↔ Social focus |
| voltage_sensitivity | Calm baseline ↔ High reactivity |
| threat_response | Approach ↔ Avoid threats |
| care_response | Independent ↔ Nurturing tendency |
| risk_tolerance | Risk-averse ↔ Risk-seeking |
| integrity_logic | Flexible principles ↔ Rigid principles |
| mirror_neuron_tuning | Low empathy ↔ High empathy |
| orderliness | Flexible/Spontaneous ↔ Ordered/Structured |

Plus 2 global moderators: `pilot_strength`, `current_load`

**Location**: `assist/esp-organizer/internal/domain/etp/spectra.go`

---

## Data Architecture

### ⚠️ Dual Data System — Critical for AI Assistants

The skills-map-platform uses **two separate data stores** that serve different purposes.
AI assistants and collaborators MUST understand which store holds which data to avoid
querying the wrong source.

#### Weaviate (Vector Database) — The CHISG Knowledge Graph
**Port**: 8081 (humanos-weaviate), 8088 (skills-map weaviate admin)
**Container**: `humanos-weaviate-1`

This is the **source of truth for skills, skill links, and course definitions**.
The explore page (`/explore`) reads subjects and courses from here.

| Class | ~Count | Purpose |
|-------|--------|---------|  
| CHISGElement | 579 | Skills with domain, description, ETP mappings |
| SkillLink | 1,120 | Parent→offspring skill relationships |
| CourseSkillSuggestions | 27 | **Subjects, courses, and their skill lists** |
| Documentation | 3,764 | Indexed project docs |
| SemanticLinks | — | Concept relationships |
| MedicalExcerpt | — | LAO revision content |

**CourseSkillSuggestions fields**: `course`, `course_code`, `subject`, `key_stage`,
`year_group`, `description`, `core_skills[]`, `practice_skills[]`, `implicit_skills[]`

#### MySQL — Operational/Relational Data
**Port**: 3307 (host) → 3306 (container)
**Container**: `skills-mysql`
**Database**: `dare2lead`
**Credentials**: `root` / `skills_password123`

Used for user accounts, skill assignments to courses, map positions, and markbook scores.
The `skills_key` table here may be **empty** — skills are defined in Weaviate.

| Table | Purpose |
|-------|---------|  
| skills_key | Skill definitions (may be empty — CHISG is in Weaviate) |
| skill_link | Skill parent/offspring links (operational) |
| skill_attached | Which skills are attached to which course |
| skillsmapelements | Node positions for the tree visualisation |
| skillsmarkbook | Student proficiency scores |
| subject | Subject records |
| courses | Course records |

### Other Databases
- **MongoDB**: Raw content storage (lessons, overviews)
---

## Key File Locations

### Documentation (Source of Truth)
```
assist/docs/
├── PRODUCT_VISION.md          # Strategic vision
├── TECHNICAL_ARCHITECTURE.md  # System architecture
├── FEATURE_PROJECT_MATRIX.md  # Feature-project mapping
├── chisg/
│   ├── CHISG_KNOWLEDGE_GRAPH_ARCHITECTURE.md
│   └── ETP_DEFINITIONS.md
└── humanos/
    └── CLASSROOM_DYNAMICS_HARMONICS.md
```

### Backend (Go)
```
assist/esp-organizer/
├── internal/
│   ├── domain/etp/spectra.go  # ETP spectrum definitions
│   └── domain/chisg/          # CHISG logic
└── cmd/server/main.go         # API entry

DRB/
├── main.go                    # Financial API
└── Data/                      # Data handlers
```

### Frontend (React/TypeScript)
```
assist/frontend/src/
├── pages/
│   └── ETPProfilePage.tsx     # ETP profile UI
├── components/
│   └── ETP/                   # ETP components

SkillsMarkbookMobile/
├── app/index.tsx              # React Native entry
└── prototypes/                # UI experiments
```

### skills-map-platform (SEPARATE REPO — NOT inside assist/)
**Location**: `/Users/michaelstewart/Coding/skills-map-platform/`

⚠️ `assist/skills-map-platform/` contains only old cached frontend artifacts.
The real project is a **separate git repo** outside the assist workspace.

```
skills-map-platform/           # /Users/michaelstewart/Coding/skills-map-platform/
├── api/
│   ├── skills/skills.go       # Skill CRUD + map save/load handlers
│   ├── subjects/SubjectsAndTopics.go  # Subject/course CRUD
│   ├── routes/routes.go       # All route definitions (public vs protected)
│   ├── handlers/course_suggestions.go # Weaviate-backed course/subject API
│   ├── weaviate/client.go     # Weaviate GraphQL client
│   ├── config/config.go       # DB config (env vars)
│   ├── auth/                  # JWT auth
│   └── structures/structures.go # Go struct definitions
├── frontend/
│   └── src/
│       ├── components/
│       │   ├── GraphExplorerWithLayout.tsx  # /explore page
│       │   ├── SkillsTree.tsx              # /skills-tree page
│       │   └── layout/NavigationDrawer.tsx  # Sidebar (fetches from Weaviate)
│       ├── config/api.ts      # API endpoint definitions
│       └── routes.tsx         # Route → component mapping
├── init.sql                   # Full MySQL schema (dare2lead)
└── docker-compose.yml
```

**Docker containers**:
| Container | Internal Port | External Port | Purpose |
|-----------|--------------|---------------|---------|  
| skills-api | 8080 | (Docker network only) | Go API |
| skills-frontend | 80 | (Docker network only) | React app |
| skills-mysql | 3306 | 3307 | MySQL (dare2lead) |
| weaviate | 8080 | 8088 | skills-map Weaviate |
| humanos-weaviate-1 | 8080 | 8081 | CHISG Weaviate (primary) |

**Key API routes (public — no auth)**:
- `GET /api/courses/subjects` → Weaviate CourseSkillSuggestions (powers explore menu)
- `GET /api/courses/by-subject/{subject}` → Weaviate courses for a subject
- `GET /api/public/skillsandlinks` → MySQL skills_key + skill_link
- `GET /api/chisg/skillsandlinks` → Weaviate CHISGElement + SkillLink

**Key API routes (protected — JWT required)**:
- `POST /api/SkillsMapPOST` → Save map node positions
- `GET /api/SkillsAttachedCourseGET/{id}` → Skills attached to a course
- `POST /api/SkillAttachedPOST` → Attach skill to course
- `POST /api/subjectPOST` / `POST /api/coursePOST` → Create subject/course (MySQL)

### Tools & Scripts
```
assist/tools/
├── skills/
│   ├── check_weaviate_skills.py    # Check if skills exist in Weaviate CHISG
│   └── add_reception_course_weaviate.py  # Add courses to Weaviate
├── student-weaviate/          # Weaviate schema setup
└── scripts/
    ├── backup-weaviate-to-s3.sh
    └── generate-project-index.sh (this script)

humanOS/scripts/
├── search_docs.py             # Weaviate doc search
└── ingest_docs.py             # Doc ingestion
```

---

## Quick Start Commands

### Search Documentation (Weaviate)
```bash
cd /Users/michaelstewart/Coding/humanOS/scripts
source .venv/bin/activate
python search_docs.py "your query" -c
```

### Run Local Development
```bash
cd /Users/michaelstewart/Coding/assist
make up                        # Start all containers
make logs                      # View logs
```

### Backup Weaviate to S3
```bash
./scripts/backup-weaviate-to-s3.sh http://localhost:8081 esp-weaviate-backups
```

### Regenerate This Index
```bash
make index
# or directly:
./scripts/generate-project-index.sh
```

---

## Context for Collaborators

### What This Solves
1. **Fragmented Student Data**: Schools have competency data in spreadsheets, behaviour data in separate systems, and no unified view.
2. **Teacher Cognitive Load**: Real-time logging is impossible with current tools; teachers resort to memory and end-of-day notes.
3. **Trust Blindspots**: Multi-academy trusts lack consolidated financial visibility until it's too late.
4. **Revision Resistance**: Students avoid homework; current tools are intimidating and high-effort.

### The Approach
- **Little and Often**: 10-minute sessions, lowest possible barrier to entry
- **Sticker Album Metaphor**: Celebrate collection, not completion
- **Grounded AI**: CHISG prevents hallucinations by validating against structured knowledge
- **Eyes-Up Teaching**: Log behaviours without looking down at a device

### Why These Projects Matter for Data Management Roles
- **DRB**: Concept demo showing Open Banking integration, financial data consolidation, and governance dashboards for MATs
- **Skills Map / CHISG**: Foundational tool — vector database management, semantic search, and knowledge graph design applied across all products
- **SkillsMarkbookMobile**: Real-time data capture, competency tracking, and cross-system sync
- **humanOS**: Foundational tool — AI coaching orchestrator with safeguarding framework, applied across products

---

