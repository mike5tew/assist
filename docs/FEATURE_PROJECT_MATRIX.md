# Feature-Project Matrix

**Purpose**: Maps features to projects to identify overlaps, gaps, and consolidation opportunities.

**Last Updated**: 2026-02-02

---

## Feature Categories

### 1. Student-Facing Learning Tools

| Feature | LAOMobile | ESP Platform | SkillsMarkbook | skills-map-platform |
|---------|:---------:|:------------:|:--------------:|:-------------------:|
| Lesson Overviews/Narratives | ✅ | ⬜ | ⬜ | ⬜ |
| Spreeder Mode (speed reading) | ✅ | ⬜ | ⬜ | ⬜ |
| Downloadable Audio Summaries | ✅ | ⬜ | ⬜ | ⬜ |
| Automated Glossaries | ✅ | ⬜ | ⬜ | ⬜ |
| Keyword-Definition Drill | ✅ | ⬜ | ⬜ | ⬜ |
| Multiple Choice Questions | ✅ | 🔲 | ⬜ | ⬜ |
| Subject Chat Mode (LLM) | 🔲 | 🔲 | ⬜ | ⬜ |
| Curriculum Heatmap | ✅ | 🔲 | ⬜ | ⬜ |
| Spaced Repetition | 🔲 | ⬜ | ⬜ | ⬜ |

### 2. Skills & Competency Tracking

| Feature | LAOMobile | ESP Platform | SkillsMarkbook | skills-map-platform |
|---------|:---------:|:------------:|:--------------:|:-------------------:|
| Skills Tree Visualization | ⬜ | ✅ | 🔲 | ✅ |
| Skills Map Rendering | ⬜ | ✅ | ⬜ | ✅ |
| Competency-Based Progression | ⬜ | ✅ | ✅ | ✅ |
| Skill Relationship Links | ⬜ | ✅ | ⬜ | ✅ |
| CHISG Integration | ✅ | ✅ | 🔲 | ⬜ |
| ETP Profiles (8 spectra) | ⬜ | ✅ | ✅ | ⬜ |
| Sticker Album UI | ⬜ | 🔲 | ✅ | ⬜ |

### 3. Teacher/Admin Tools

| Feature | LAOMobile | ESP Platform | SkillsMarkbook | ESP Seating* | ESP Behaviour* |
|---------|:---------:|:------------:|:--------------:|:------------:|:--------------:|
| Class Management | ⬜ | ✅ | ⬜ | ✅ | ✅ |
| Seating Plan Interface | ⬜ | 🔲 | ⬜ | ✅ | ⬜ |
| Behaviour Logging | ⬜ | 🔲 | ✅ | ⬜ | ✅ |
| Action Recorder (4 markers) | ⬜ | ⬜ | ✅ | ⬜ | ⬜ |
| Profile Sliders (15 spectra→8) | ⬜ | ✅ | ✅ | ⬜ | ⬜ |
| Real-time Progress View | ⬜ | ✅ | ✅ | ⬜ | ✅ |
| Work Assignment | ⬜ | ✅ | ⬜ | ⬜ | ⬜ |

*\* ESP Seating and ESP Behaviour Lite were WinDev Mobile apps — original code lost, need rewrite*

### 4. Knowledge Graph & AI

| Feature | assist | humanOS | CHISG Engine |
|---------|:------:|:-------:|:------------:|
| Semantic Link Extraction | ✅ | ⬜ | ✅ |
| CHISG Parsing/Validation | ✅ | ⬜ | ✅ |
| Gap Analysis | 🔲 | ⬜ | ✅ |
| Hallucination Reduction | 🔲 | ⬜ | ✅ |
| LLM Self-Evaluation | 🔲 | ⬜ | 🔲 |
| Emotional Motivation Model | ⬜ | ✅ | ⬜ |
| Teaching Pattern Codification | ⬜ | ✅ | ⬜ |
| Classroom Interaction Patterns | ⬜ | ✅ | ⬜ |

### 5. Infrastructure & Data

| Feature | assist (main) | Weaviate | MongoDB | MySQL |
|---------|:-------------:|:--------:|:-------:|:-----:|
| Documentation Index | ✅ | ✅ | ⬜ | ⬜ |
| Skills/CHISG Elements | ✅ | ✅ | ⬜ | ⬜ |
| Semantic Links | ✅ | ✅ | ⬜ | ⬜ |
| User/Student Data | ⬜ | ⬜ | ⬜ | ✅ |
| Lesson Content | ⬜ | ⬜ | ✅ | ✅ |
| Audio Files | ⬜ | ⬜ | ⬜ | S3 |
| ETP Profiles | ✅ | ✅ | ⬜ | ⬜ |

---

## Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Implemented and active |
| 🔲 | Planned/Partial |
| ⬜ | Not applicable/Not present |

---

## Project Status Summary

| Project | Status | Key Value | Consolidation Notes |
|---------|--------|-----------|---------------------|
| **assist** | Active monorepo | Documentation, orchestration | Primary development target |
| **LAOMobile** | App Store approved | GCSE revision, heatmaps | Keep separate; next step is live release/public distribution alignment |
| **ESPProj** | Legacy migration | Full MIS features | Migrate to assist |
| **SkillsMarkbookMobile** | Partial | Primary skills + stickers | Consider merge with LAOMobile or standalone |
| **humanOS** | Research/Docs | Emotional models, teaching | Integrate into Weaviate + assist docs |
| **skills-map-platform** | Maintenance | Skills tree rendering | Check features before archiving |
| **ESP Seating** | Lost (WinDev) | Seating plan UI | REWRITE NEEDED |
| **ESP Behaviour Lite** | Lost (WinDev) | Behaviour logging | REWRITE NEEDED |

---

## Recommended Consolidation Path

### Phase 1: Immediate (Backups & Stability)
1. ✅ Backup Weaviate to S3 (script created)
2. Export humanOS docs → Weaviate
3. Archive skills-map-platform after feature audit

### Phase 2: Feature Migration
1. Migrate ESP legacy features → assist/esp-organizer
2. Rewrite Seating Plan interface (React/TypeScript)
3. Rewrite Behaviour logging (can share UI with SkillsMarkbook)

### Phase 3: Portfolio Site Merge
1. Combine portfolio + ESP showcase site
2. Structure:
   - `/` — Big ideas showcase (ETP, CHISG, humanOS concepts)
   - `/apps` — App Store links (LAOMobile, ESP Seating, Behaviour)
   - `/demo` — Live demos (Skills Tree, Curriculum Heatmap)
   - `/platform` — ESP Platform login (school-facing)

---

## ESP Seating / Behaviour Lite: Rewrite Scope

Since original WinDev code is lost, these need fresh implementations:

### ESP Seating (Rewrite)
- **Purpose**: Visual seating plan editor + student placement
- **Core Features**:
  - Drag-drop desk/student placement
  - Classroom template library
  - Student profile quick-view on hover
  - Export to PDF/print
- **Tech**: React + Canvas/SVG, Go API
- **Estimate**: 2-3 weeks for MVP

### ESP Behaviour Lite (Rewrite)
- **Purpose**: Quick behaviour logging during lessons
- **Core Features**:
  - 4-button quick log (2 positive, 2 negative)
  - Student picker (class view or search)
  - Session timeline view
  - Sync to ESP Platform
- **Tech**: React Native (share code with SkillsMarkbook)
- **Estimate**: 1-2 weeks for MVP

---

## Weaviate Classes for Backup

| Class | Count | Priority |
|-------|-------|----------|
| CHISGElement | 534 | 🔴 Critical |
| SkillLink | 1,120 | 🔴 Critical |
| Documentation | 3,764 | 🟡 High |
| SemanticLinks | 0 | 🟢 When populated |
| MedicalExcerpt | ? | 🟡 High (LAO content) |
| CourseSkillSuggestions | ? | 🟢 Medium |
