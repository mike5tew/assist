# Action Plan: HumanOS Ecosystem Development

**Last Updated**: 2025-03-14
**Current Focus**: GCSE Revision Tool (First Revenue) → CHISG Enhancement → Skills Tree Rising

## Setup Complete ✅

- [x] GitHub repository created
- [x] Initial commit pushed
- [x] Go module initialized
- [x] Project structure organized
- [x] **Backend: Go (decision finalized)** ✅
- [x] **Complete classroom framework preserved** ✅ (moved to `/docs/reference/`)
- [x] **Individual concepts being extracted** ✅ (for AI tutor MVP)
- [x] **Deprecated code cleaned up** ✅ (removed `/project/backend/`)
- [x] **Documentation consolidated** ✅ (moved to `/docs/`)

## Strategic Overview

### The 80/20 Principle
**Focus on the 20% that generates 80% of value:**
1. **HumanOS Core** (barrier detection + intervention) - foundation for everything
2. **CHISG Integration** - makes AI responses smart
3. **Payment Infrastructure** - enables all revenue
4. **Three Target Products** - AI Tutor, GCSE Tool, Skills Tree

**Deliberately NOT doing (for now):**
- Full ESP suite completion (50-70% → 100% is massive work)
- Every service at 100% (diminishing returns)
- Features beyond MVP (perfectionism trap)

## Phase 1: Foundation & First Revenue (Commercial Pilot) ✅
**Goal**: Launch GCSE Revision Tool MVP and acquire first pilot users.

### Priority 1: LittleAndOften (GCSE Revision Tool) ✅
**Current**: 95% complete (Pilot Ready)
**Target**: 100% complete (Global Release)

**Completed**:
- [x] **Core Content Engine**: Restored relational SQLite architecture with full curriculum support.
- [x] **Dynamic Summaries**: File-based Markdown rendering for all lesson content.
- [x] **Study Reward Logic**: Implemented Daily Effort and Activity Completion trackers.
- [x] **Mobile UX**: Standardized headers, navigation, and activity modes (Audio, Speed Reader, Card).
- [x] **Deployment**: Successful EAS build submission to TestFlight.

**Next Work Needed**:
- [ ] **Content Expansion**: Populate remaining summaries for all 25 modules.
- [ ] **Multi-Subject Tuning**: Polish the subject selection flow for broader appeal.
- [ ] **AI Question Gen**: Prototype auto-generated quiz questions from Markdown summaries.

### Priority 2: Skills Tree Rising (Primary Schools)
**Goal**: Launch the Sticker Album extension and high-fidelity dashboard.

**Next Work Needed**:
- [ ] **Physical-Digital Loop**: Implement the [Sticker Album Extension](../../assist/docs/STICKER_ALBUM_EXTENSION.md) logic (rewarding physical stamps for digital effort).
- [ ] **Strategic Portfolio**: Refine the "Non-School Functions" dashboard for Finance, HR, and Estates.

---

## Phase 2: CHISG Enhancement & Sticker Album Integration (Current Focus)
  - [ ] Implement "Question Practice" interface.
- [ ] **Week 3: UI, Auth & Progress Tracking**
  - [ ] Student dashboard with progress visualization.
  - [ ] User authentication (can leverage existing JWT system).
- [ ] **Week 4: Commercialization & Launch**
  - [ ] Set up Stripe for subscription payments.
  - [ ] Create landing page.
  - [ ] Beta launch and feedback collection.

### Priority 2: Semantic Link Extraction Tool (Enabler)
**Current**: 30% complete (spec & architecture designed)
**Target**: 100% complete (live, integrated tool)

**Purpose**: Integrated tool to rapidly populate the knowledge graph via the PDF extraction pipeline. Bridges PDF text extraction → semantic relationship curation → Weaviate storage.

**Architecture**: React frontend + Go API (integrated with `esp-organizer` backend)

**Implementation Plan (4-week sprint)**:

**Week 1 (Backend foundation — 3 days)**
- [ ] Add Go API endpoints to `esp-organizer`:
  - `POST /api/semantic-links/extract` — accept selections and persist links
  - `GET /api/semantic-links/search?term=X` — query by term
  - `POST /api/semantic-links/validate` — trigger quality scoring
- [ ] Implement `SemanticLink` model with provenance fields
  - `source_term, target_term, relation_type`
  - `document_id, page_number, sentence_id, text_context`
  - `position_context` (char offsets), `confidence`, `hierarchy_level`, `quality_flags`
- [ ] Wire Weaviate storage & index for full-text search
- [ ] Wire PostgreSQL for audit and provenance logging

**Week 2 (Frontend & UX — 4 days)**
- [ ] Create React component `<SemanticLinkExtractor />`
  - PDF viewer with selectable text overlay
  - 3-click workflow: select source → select target → choose relation
  - Visual highlighting and inline relation suggestions
- [ ] Relation dropdown with domain-aware presets
- [ ] Quick quality flags (High / Needs verification / Problematic)

**Week 3 (Pipeline integration & inference — 3 days)**
- [ ] Hook into PDF extraction pipeline: add "Extract Semantic Links" action to OCR results
- [ ] Implement hierarchy inference (rule-based + thresholds)
- [ ] Duplicate & contradiction detection (identical triples, similarity threshold)

**Week 4 (Validation, batch, launch — 4 days)**
- [ ] Vague relation rejection and speculative language detection
- [ ] Batch processing mode for textbook patterns
- [ ] E2E & performance testing (<10s/link), UAT with teachers/curators
- [ ] Deploy to staging and production; document curator workflow

**Acceptance Criteria / Success Metrics**
- Extraction speed: < 10 seconds per semantic link
- Accuracy: ≥ 95% of manual review calls extraction meaningful
- Curator throughput: 50+ links/hour for experienced users
- Data quality: < 5% vague relationships in final dataset

### Priority 2b: Academic Paper CHISG Pipeline (McGrath Review Workflow)
**Context**: Michael Stewart (PhD Biophysics) is working through microbiology papers to build CHISG training data and knowledge graph content. Volume constraint means a fully manual approach is not sustainable across 50+ papers.

**Agreed Workflow**:
1. **Few-shot prompt development** — Work through a small number of papers collaboratively (MS + Copilot) to produce high-quality human-curated reduced semantic cores (as JSON). These serve as few-shot examples.
2. **Pipeline extraction** — `extract_chisg_papers.py` processes PDFs via Bedrock, producing CHISG elements using the few-shot examples in the prompt. Output includes `evidence_context` (conditions under which the supporting evidence was gathered — not a limiter on the claim, but provenance of the evidence) and conformant controlled vocabulary.
3. **Human-readable review artefact** — After pipeline runs on a paper, auto-generate a clean human-readable summary (Markdown or tabular) of extracted links, organised by figure/section, suitable for a non-programmer reviewer.
4. **McGrath review** — Michael McGrath reads the paper and marks up the review artefact: flagging missing links, incorrect relation types, wrong evidence_context, or vocabulary violations.
5. **Correction loop** — Corrections feed back as additional few-shot examples, improving subsequent pipeline runs.

**Key design decisions recorded here**:
- `evidence_context` field stores the experimental conditions of the supporting evidence (organism, method, timepoint, treatment). This is provenance — *where we got the evidence* — not a claim that the relationship only holds in those conditions.
- Human JSON extractions are the reduced semantic core only (entity_a, relation, entity_b, evidence_context). Pipeline expands to full CHISG element (backward relation, full context array, trust metadata stubs).
- Controlled vocabulary reference: `extract_chisg_papers.py` lines 60–70 (20 relations, microbiology domain).

**Immediate next actions**:
- [ ] Update `extract_chisg_papers.py` system prompt: add `evidence_context` to output schema, rename `context` to `text_quote`, add `inverse_relation` field
- [ ] Add le_chen_eLife_2022 human JSON as first few-shot example in prompt
- [ ] Work through 2–3 more papers with MS to build the few-shot set to ~5 examples
- [ ] Build human-readable review artefact generator (script or pipeline step)
- [ ] Define McGrath review workflow and handoff format

### Priority 3: Integrated Insight Markbook (Skills Map + ETP Profile Builder)
**Strategic Context**: Sprung from the drb Ignite business case (Jan 2026). A tool that combines "Safe Hands" (attainment data) with "Smart Minds" (pedagogical logic) using the AISA methodology.

**Core Vision**:
- **Skills Markbook (The Compass)**: Real-time "Sticker Album" tracking of competencies. 
- **HumanOS Engine (The Map)**: Using "Recursive Mastery" logic (Stimulus → ETP → Skill → Behaviour) to explain *why* students succeed or fail in specific contexts.
- **Goal**: Workload displacement (5x faster than marking) while providing headteacher-level strategic insights through CHISG Knowledge Graph integrity.

**Implementation Plan (4-week sprint)**:
- [ ] **Week 1: Unified Logic Architecture**
  - [ ] Implement "Recursive Mastery" schema: Link Skills to the ETPs they master.
  - [ ] Map the "Behavioral Pipeline" (Stimulus -> Mechanism -> ETP -> Response) into the MySQL backend.
- [ ] **Week 2: Profile Builder & Passive Pilot UI**
  - [ ] Implement the "2 Positive / 2 Negative" Temporal Marker system for objective behavior logging.
  - [ ] Build the 15-spectrum Sensitivity Profile sliders with mandatory justification notes for downward shifts.
  - [ ] Implement spatial memory (seating plan) interface for high-speed logging.
- [ ] **Week 3: CHISG Integrity Layer**
  - [ ] Create the CHISG validation layer to prevent AI diagnostic "hallucinations."
  - [ ] Cross-link teacher observations with the CHISG Knowledge Graph.
- [ ] **Week 4: Strategic Analytics (MAT Dashboard)**
  - [ ] Develop longitudinal profile visualization (Heatmaps of Social Gravity & Voltage).
  - [ ] Design "Safe Hands" automated reporting engine.

## Phase 2: Product Development & Initial Revenue (Months 3-4)
**Goal**: Launch GCSE Tool + Skills Tree Rising beta  
**Revenue Target**: $500-2000/month (combined)

### Month 3: GCSE Tool Launch
**Current**: 0% complete (initial setup)  
**Target**: 100% complete (live product)

**Leverages Existing Work**:
- HumanOS Core (50% → 100% in Phase 1)
- CHISG Integration (80% complete)
- Payment Infrastructure (basic version)

**New Work Needed**:
- [ ] **Week 1: Core Content Engine (Science Focus)**
  - [ ] Load GCSE Science curriculum (lesson headings).
  - [ ] Implement "Lesson Overview" generation based on headings.
  - [ ] Implement "Question Generation" based on overviews.
- [ ] **Week 2: Content Consumption Features**
  - [ ] Implement "Spreeder" speed-reading UI for overviews.
  - [ ] Implement "Downloadable Audio" feature (Text-to-Speech API integration).
  - [ ] Implement "Keywords & Definitions" extraction and display.
- [ ] **Week 3: UI & Progress Tracking**
  - [ ] Design and implement the "Curriculum-Aware Study Planner" with progress tracking.
  - [ ] Design a unified student dashboard to display all features.
- [ ] **Week 4: Commercialization & Launch**
  - [ ] Set up Stripe/PayPal for subscription payments.
  - [ ] Create landing page + SEO optimization.
  - [ ] Reach out to schools/tutors for partnerships.

**Success Metrics**:
- 100+ free tier signups
- 20+ paying subscribers (£100-200/month revenue)
- >80% accuracy on practice questions
- Avg 20 minutes/day engagement per active user

### Month 4: Skills Tree Rising
**Current**: 80% base complete (via `skills-map-platform` project)
**Target**: 95% complete (integrated and refined for B2B launch)

**Leverages Existing Work**:
- **`skills-map-platform`**: Provides a 95% complete foundation, including backend, frontend, auth, and skills visualization.

**New Work Needed (Integration & Refinement)**:
- [ ] **Week 1: Architectural Integration Plan**
  - [x] Decide on a unified database strategy (MySQL vs. MongoDB/Weaviate). **DECISION: Hybrid Model.**
  - [ ] **Implement MySQL connector** in the `esp-organizer` backend.
  - [ ] Plan the merge of the `skills-map-platform` Go backend features into the `esp-organizer` backend.
  - [ ] Plan the integration of the `SkillsTree.tsx` component into the `assist` frontend.
- [ ] **Week 2: Leadership Skills Taxonomy & Content**
  - [ ] Work with collaborators to define skills and assessment criteria.
- [ ] **Week 3: B2B Feature Polish**
  - [ ] Implement PDF report generation for trainers.
  - [ ] Refine cohort progress visualization.
  - [ ] Implement certificate generation.
- [ ] **Week 4: Payment & Trainer Tools**
  - [ ] Implement revenue share model.
  - [ ] Polish the trainer dashboard for cohort management.

**Launch Strategy**:
- [ ] Pilot with collaborators' existing clients
- [ ] Revenue share: 70% collaborators, 30% you
- [ ] Target: 20-50 learners in first cohort
- [ ] Pricing: £200-500/learner (handled by collaborators)

**Success Metrics**:
- 20+ learners enrolled
- £1000-2500 revenue share (£300-750 to you)
- >90% completion rate for pilot cohort
- Testimonials for future marketing

## Phase 3: Scaling & Product Expansion (Months 5-6)
**Goal**: Scale existing products + launch Solo Skills Map  
**Revenue Target**: $2000-5000/month combined

### Month 5: Scale AI Tutor & GCSE Tool
**Tasks**:
- [ ] AI Tutor improvements
  - [ ] **Implement Curriculum Intelligence**: Use CHISG + LLM to analyze curriculum structure and provide insights.
  - [ ] Add voice interface (leveraging existing OpenAI APIs)
  - [ ] Multi-subject expansion (math, English, languages)
  - [ ] Parent/teacher dashboards
  - [ ] Marketing: $500/month ad spend (Google, Facebook)
  - [ ] Target: 50-100 paid subscribers ($1000-2000/month)
- [ ] GCSE Tool expansion
  - [ ] Add more subjects (English, Math, Geography)
  - [ ] Partnership with tutoring centers
  - [ ] School pilot program (freemium for schools)
  - [ ] Target: 100+ free, 30+ paid subscribers (£300-1000/month)

### Month 6: Solo Skills Map Launch
**Current**: 0% career content, 80% base complete  
**Target**: 80% complete (public beta)

**New Work Needed**:
- [ ] Week 1: Career skills taxonomy
  - [ ] Research job postings for common skill requirements
  - [ ] Map skills to career paths (tech, business, creative, etc.)
  - [ ] Create skill assessment questionnaires
- [ ] Week 2: CV generation engine
  - [ ] Templates for different industries
  - [ ] AI-assisted bullet point writing
  - [ ] ATS optimization (keyword matching)
- [ ] Week 3: Job matching algorithm
  - [ ] Skill requirements from job postings
  - [ ] Match user skills to opportunities
  - [ ] Gap analysis + learning recommendations
- [ ] Week 4: LinkedIn integration
  - [ ] Import skills from LinkedIn
  - [ ] Export updated profile sections
  - [ ] Share achievements + certificates

**Monetization**:
- Free tier: Basic skill tracking
- £10/month: CV generation + job matching
- £20/month: + AI career coaching + LinkedIn optimization

**Marketing**:
- [ ] Product Hunt launch
- [ ] LinkedIn articles about skill-based hiring
- [ ] Partnership with career coaching services
- [ ] Target: 100+ free signups, 20+ paid (£200-400/month)

**Phase 3 Success Metrics**:
- **Total MRR**: £1500-3000 ($2000-4000)
- **Active users**: 200+ across all products
- **Churn rate**: <10% monthly
- **NPS score**: >40

## Phase 4: Immunology Assistant & Research Tools (Months 7-9)
**Goal**: Fulfill personal commitment + build academic credibility  
**Revenue Target**: $0-500/month (not primary goal)

### Month 7-8: Immunology Assistant Development
**Current**: 20% complete (content loading)  
**Target**: 80% complete (functional for exam prep)

**Leverages Existing Work**:
- CHISG (60% → 80% in Phase 1)
- PDF RAG (90% complete)
- Skills Map (80% complete)
- ESP Assist (60% → 80% in Phase 1-2)

**New Work Needed**:
- [ ] Medical terminology ontology
  - [ ] Import MeSH (Medical Subject Headings)
  - [ ] Map immunology concepts to clinical pathology
  - [ ] Build prerequisite chains (undergrad → clinical)
- [ ] Research paper integration
  - [ ] Load key immunology papers (Nature Immunology, etc.)
  - [ ] Extract figures + explanations
  - [ ] Link to exam board specifications
- [ ] Clinical case studies
  - [ ] Integrate case presentations
  - [ ] Link symptoms → pathology → diagnosis pathway
  - [ ] Practice question generation based on cases
- [ ] Study schedule generation
  - [ ] Based on exam date + current knowledge
  - [ ] Spaced repetition algorithm
  - [ ] Progress tracking toward exam readiness

**Launch Strategy**:
- [ ] Beta test with medical student friends
- [ ] Partner with university immunology departments
- [ ] Academic paper: "AI-Assisted Medical Education Using Semantic Knowledge Graphs"
- [ ] Pricing: Free for students, £50-100/month for clinicians (CPD)

### Month 9: Research Tools & Academic Credibility
**Tasks**:
- [ ] Write academic paper on HumanOS framework
  - [ ] "Emotional Trigger Points: A Framework for Adaptive Educational AI"
  - [ ] Submit to conferences (EDM, AIED, LAK)
- [ ] Open-source core components
  - [ ] HumanOS ETP framework (GitHub)
  - [ ] CHISG knowledge graph builder (GitHub)
  - [ ] Build developer community
- [ ] Conference presentations
  - [ ] Demo at AI in Education conferences
  - [ ] Network with EdTech researchers
  - [ ] Explore research partnerships

### Future: University-Level Knowledge Visualization
**Context**: Unlike GCSE's hierarchical curriculum (Subject → Module → Topic → Lesson), university knowledge is web-structured with cross-cutting dependencies. Treemaps break down because concepts don't fit neatly in boxes - "linear algebra" appears in physics, CS, economics, pure maths simultaneously.

**Capability**: Force-directed semantic map of student understanding

**Visualization Approach**:
- **Position**: t-SNE/UMAP projection of concept embeddings (nearby = semantically related)
- **Color**: Student mastery level (red → yellow → green)
- **Edges**: `requires_understanding_of` prerequisite links from CHISG
- **Insight**: Reveals blocking concepts - "You can't understand thermodynamics because you're missing partial derivatives"

**Technical Requirements**:
- [ ] Student concept mastery tracking (per CHISG node, not per curriculum item)
- [ ] Concept embeddings from CHISG definitions
- [ ] t-SNE/UMAP projection to 2D
- [ ] Force-directed graph layout with prerequisite edges
- [ ] Interactive visualization (zoom, pan, click for details)

**Data Model**:
```sql
student_concept_mastery:
  - student_id
  - concept_id (CHISG node)
  - mastery_level (0.0-1.0)
  - last_assessed (timestamp)
  - evidence[] (quiz results, assignments, self-assessment)
```

**CHISG Integration**:
- Uses existing `requires_understanding_of` / `enables_understanding_of` links
- Cross-domain pattern detection highlights "same concept, different course" opportunities
- Gap detection: find unmastered concepts blocking downstream understanding

**Use Cases**:
- Student self-assessment: "Where are my knowledge gaps?"
- Advisor tool: "This student should take X before Y"
- Course recommendation: "Based on your graph, these modules will fill structural holes"

**Prerequisite**: CHISG graph with sufficient university-level content and prerequisite links.

## Phase 5: Federated Learning & Long-term Vision (Months 10-12)
**Goal**: Build collective intelligence system  
**Revenue Target**: Improved product value (indirect revenue)

### Month 10-11: Federated Coordinator Development
**Current**: 0% complete  
**Target**: 50% complete (functional aggregation)

**Architecture**:
```go
type FederatedCoordinator struct {
    patternAggregator  *PatternAggregator
    privacyEngine      *DifferentialPrivacy
    modelDistributor   *ModelDistributor
}

// Deployments submit anonymized patterns
func (fc *FederatedCoordinator) SubmitPattern(pattern AnonymousPattern) error {
    // Validate no PII present
    if fc.privacyEngine.ContainsPII(pattern) {
        return errors.New("PII detected - pattern rejected")
    }
    
    // Aggregate with existing data
    fc.patternAggregator.Add(pattern)
    
    // Check if enough data for model update
    if fc.patternAggregator.ReadyForUpdate() {
        improvedModel := fc.trainModel()
        fc.modelDistributor.Distribute(improvedModel)
    }
    
    return nil
}
```

**Tasks**:
- [ ] Week 1-2: Privacy infrastructure
  - [ ] Differential privacy implementation
  - [ ] K-anonymity checks
  - [ ] PII detection + rejection
  - [ ] Opt-in/opt-out management
- [ ] Week 3-4: Pattern aggregation
  - [ ] Time-series aggregation (weekly/monthly insights)
  - [ ] Cross-deployment pattern matching
  - [ ] Intervention effectiveness scoring
- [ ] Week 5-6: Model distribution
  - [ ] Improved barrier detection models
  - [ ] Better intervention selection logic
  - [ ] Age-appropriate language improvements
- [ ] Week 7-8: Monitoring + compliance
  - [ ] Privacy audit logs
  - [ ] GDPR compliance checks
  - [ ] Deployment health monitoring

**Success Metrics**:
- 10+ deployments contributing patterns
- >95% privacy guarantee compliance
- Measurable improvement in intervention success rates
- Zero PII leaks or breaches

### Month 12: Full ESP Suite Consideration
**Decision Point**: Do we complete the full ESP suite?

**Factors to consider**:
- Are AI Tutor + GCSE Tool + Skills Tree generating $5000+/month?
- Do we have school customers asking for full suite?
- Is there funding available (grants, investors)?
- Is completion worth 6-12 months of work?

**If YES, proceed with full ESP completion**:
- Allocate 6-12 months for remaining 30-50% of ESP components
- Target: School licenses at $10,000-30,000/year
- Need: Sales team, customer support, implementation consultants

**If NO, continue focusing on high-value products**:
- Double down on AI Tutor (expand subjects, add features)
- Grow GCSE Tool (more exam boards, international)
- Scale Solo Skills Map (enterprise version for companies)
- Explore new product ideas (language learning, test prep)

## Contingency Plans

### If Phase 1 Takes Longer Than Expected
**Adjust**:
- Move GCSE Tool to Month 4 (delay but don't skip)
- Simplify AI Tutor MVP (fewer features, faster launch)
- Reduce scope of Skills Tree Rising (manual workarounds)

### If No Job Offers by Month 2
**Adjust**:
- Double down on AI Tutor revenue (need to replace employment income)
- Consider contract/freelance work (AI consulting)
- Accelerate GCSE Tool launch (family need is real)

### If Products Don't Gain Traction
**Pivot Options**:
- B2B focus: Sell to schools/tutoring centers instead of direct-to-consumer
- White-label: License technology to existing EdTech companies
- Consulting: Offer AI implementation services using your tech stack
- Open-source + support model: Free product, charge for hosting/support

## Success Metrics Summary

### By Month 2 (End of Phase 1)
- [ ] HumanOS: 50% complete, production-ready
- [ ] AI Tutor: Launched with 10+ beta users
- [ ] Job offers: 2-5 interviews, 1+ offer
- [ ] Revenue: $100-500/month from beta users

### By Month 4 (End of Phase 2)
- [ ] GCSE Tool: Launched with 30+ paid subscribers
- [ ] Skills Tree Rising: Delivered to collaborators, 20+ learners
- [ ] Revenue: $500-1500/month combined

### By Month 6 (End of Phase 3)
- [ ] AI Tutor: 50-100 paid subscribers
- [ ] GCSE Tool: 50+ paid subscribers
- [ ] Solo Skills Map: Launched, 20+ paid users
- [ ] Revenue: $2000-4000/month combined

### By Month 9 (End of Phase 4)
- [ ] Immunology Assistant: Functional for personal use
- [ ] Academic paper: Submitted to conference
- [ ] Research credibility: Established in AI + Education

### By Month 12 (End of Phase 5)
- [ ] Federated learning: 10+ deployments contributing
- [ ] Strategic decision: Full ESP suite or continue current products
- [ ] Revenue: $5000+/month or full-time employment secured

## Risk Mitigation

**Technical Risks**:
- HumanOS complexity underestimated → Build incrementally, test continuously
- CHISG knowledge graph accuracy issues → Start with narrow domains, expand gradually
- Federated learning privacy concerns → Over-engineer privacy, get legal review

**Business Risks**:
- Low user adoption → Extensive beta testing, iterate based on feedback
- High customer acquisition cost → Focus on organic growth, word-of-mouth
- Competition from established players → Emphasize unique ETP approach

**Personal Risks**:
- Burnout from overwork → Set sustainable pace, don't skip Phase 1 employment
- Scope creep → Ruthlessly prioritize 80/20 features
- Perfectionism blocking launches → Ship MVPs, improve based on real usage

## Weekly Review Cadence

**Every Friday**:
- [ ] Review week's progress against action plan
- [ ] Update completion percentages
- [ ] Identify blockers + create unblock plan
- [ ] Plan next week's priorities
- [ ] Celebrate wins (even small ones!)

**Every Month**:
- [ ] Review phase progress
- [ ] Adjust timeline if needed
- [ ] Update revenue projections
- [ ] Reassess priorities based on results

---

## 🎯 SkillsMarkBook Extension: Sticker Album Feature

**Strategic Context**: Build a printable, physical component of the Skills Markbook that makes skill mastery tangible and rewarding. Teachers award physical stickers corresponding to skills demonstrated; students collect them in an A5 album, creating a motivational artifact.

**Core Vision**:
Transform abstract skill tracking into a *collectible experience* - where the Skills Map becomes a physical album filled with custom stickers, each representing a demonstrated competency area. This combines the motivation of gamification with the permanence of a physical portfolio.

### Design Principles

1. **A5 Page Format**: Each skill area fits cleanly on a single A5 page (~148×210mm)
2. **Paired Explainer Page**: Opposite side explains the skill area in simple, student-friendly language
3. **One Sticker Per Award**: Teacher/assessor prints a sticker when skillsmark is awarded
4. **Progressive Narrative**: As students fill the album, they see growth and mastery emerge organically
5. **Physical + Digital Bridge**: Album has QR codes linking to detailed digital profile

### Phase 1 Implementation Plan (2-3 weeks)

#### Week 1: Design & Content Structure

**Deliverables**:
- [ ] Define 8-12 core skill areas for MVP (e.g., "Communication", "Problem Solving", "Creativity", "Resilience", "Collaboration", "Leadership", "Research", "Critical Thinking")
- [ ] Create A5 page templates:
  - **Front side**: Space for 5-6 stickers + skill area title + visual theme (icon/color coding)
  - **Back side**: 100-150 word explanation of the skill area, what it means, and how to demonstrate it
- [ ] Design custom sticker art (5cm × 5cm) for each skill area with:
  - Skill area name
  - Achievement tier (if multi-level, e.g., Bronze/Silver/Gold)
  - Unique visual identifier (animal, icon, color scheme)
- [ ] Create album cover template (A4 folded to A5) with:
  - Student name + class + year
  - QR code to digital profile
  - Simple instruction: "Collect stickers as you demonstrate skills"

**Database Schema Updates** (MySQL):
```sql
-- New tables
CREATE TABLE skill_areas (
  SkillAreaID INT PRIMARY KEY AUTO_INCREMENT,
  AreaName VARCHAR(50) NOT NULL,
  Description TEXT NOT NULL,
  IconID INT,
  ColorCode VARCHAR(7),
  Order INT,
  Active BOOLEAN DEFAULT TRUE
);

CREATE TABLE skill_area_stickers (
  StickerID INT PRIMARY KEY AUTO_INCREMENT,
  SkillAreaID INT NOT NULL,
  TierLevel INT (1-3 for Bronze/Silver/Gold),
  StickerDesignURL VARCHAR(255),
  DateCreated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (SkillAreaID) REFERENCES skill_areas(SkillAreaID)
);

CREATE TABLE student_sticker_awards (
  AwardID INT PRIMARY KEY AUTO_INCREMENT,
  StudentID INT NOT NULL,
  StickerID INT NOT NULL,
  AwardedByTeacherID INT NOT NULL,
  AwardedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  ObservationNotes TEXT,
  FOREIGN KEY (StudentID) REFERENCES students(StudentID),
  FOREIGN KEY (StickerID) REFERENCES skill_area_stickers(StickerID),
  FOREIGN KEY (AwardedByTeacherID) REFERENCES users(UserID)
);

CREATE TABLE student_album_progress (
  ProgressID INT PRIMARY KEY AUTO_INCREMENT,
  StudentID INT NOT NULL,
  TotalStickersAwarded INT DEFAULT 0,
  SkillAreaCompletionStatus JSON,
  LastUpdated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (StudentID) REFERENCES students(StudentID)
);
```

#### Week 2: Teacher & Student Interfaces

**Teacher Interface (Skills Markbook Mobile + Web)**:
- [ ] "Award Sticker" button in the Skills Markbook marking interface
- [ ] Quick action: Select skill area → select student(s) → confirm award
- [ ] Bulk print interface: Generate PDF with selected stickers for printing
- [ ] Dashboard showing:
  - Total stickers awarded by skill area
  - Student engagement metrics (who has more stickers)
  - "Ready to Print" sticker queue

**Student Facing** (Skills Markbook Mobile + Web):
- [ ] "My Album" tab in Skills Markbook showing:
  - Visual preview of album pages (A5 grid layout)
  - Number of stickers collected in each area
  - "Next sticker" motivation display
  - QR code to digital profile details
- [ ] Celebrate notifications when new sticker earned
- [ ] Historical view: timeline of earned stickers with dates

**Print & Export Features**:
- [ ] Generate PDF: 
  - Single album (all pages + stickers earned so far)
  - Batch print job (selected students' albums)
  - Sticker sheet (printable 5cm×5cm stickers for manual insertion)
- [ ] QR code generation linking to:
  - Student's digital profile
  - Detailed skill descriptions
  - Evidence/observations from teachers
- [ ] Export as image: album preview for sharing with parents

#### Week 3: Integration & MVP Launch

**Backend API Endpoints**:
- [ ] `POST /api/student-awards/sticker` - Award sticker to student
- [ ] `GET /api/student-album/{studentID}` - Get album progress + stickers
- [ ] `GET /api/skill-areas` - List all skill areas and their designs
- [ ] `POST /api/album/export-pdf` - Generate printable PDF album
- [ ] `POST /api/album/export-stickers` - Generate sticker sheet PDF
- [ ] `GET /api/statistics/sticker-awards` - Aggregate stats by skill area, teacher, class

**Mobile UI Components**:
- [ ] `<StickerAlbumView />` - Grid display of A5 pages with stickers
- [ ] `<AwardStickerModal />` - Quick-action modal for awarding stickers
- [ ] `<SkillAreaCard />` - Shows skill explanation + stickers earned + next milestone
- [ ] `<AlbumExportMenu />` - Export options (PDF, print, share)

**QA & Testing**:
- [ ] Sticker design renders correctly at 5cm×5cm
- [ ] PDF generation is clean and printable
- [ ] QR codes resolve correctly to student profiles
- [ ] Mobile responsiveness for A5 preview
- [ ] Data integrity: one sticker award = one print-ready sticker

### Phase 2 (Later): Advanced Features

**Tier System** (Bronze/Silver/Gold):
- [ ] Students earn multiple stickers in the same skill area as they progress
- [ ] Visual ranking: upgrade sticker color/design with each tier
- [ ] Teacher can award "Bronze" for first demonstration, "Silver" for consistent application, "Gold" for mastery
- [ ] Album automatically fills in visual tiers

**Badges & Milestones**:
- [ ] "Skill Master" badge when student collects all stickers in an area (optional achievement)
- [ ] "Album Keeper" badge for consistency (weekly sticker awards)
- [ ] "Growth Mindset" badge for progressing through tier levels

**Parent Integration**:
- [ ] Share student album with parents via secure link
- [ ] Parents can print a copy to display at home
- [ ] Optional: parents can leave encouraging comments on student profiles

**Customization**:
- [ ] School can create custom skill areas and sticker designs
- [ ] Teachers can set "house rules" (e.g., max 1 sticker/skill/week to prevent saturation)
- [ ] Different designs for different year groups (primary vs. secondary)

**Analytics & Reporting**:
- [ ] Heatmap: which skills are being demonstrated most frequently?
- [ ] Equity check: are all students receiving stickers fairly?
- [ ] Trend analysis: which students are accelerating vs. plateauing?
- [ ] Export for Ofsted/assessment: evidence of skill development over time

### Physical Production Workflow

1. **Teacher Awards in App** → System generates print job queue
2. **Print Stickers** → Batch print at school (color sticker labels)
3. **Manual or Automatic Distribution**:
   - Manual: Teacher gives sticker to student to place in album
   - Automatic: School prints full albums monthly/termly
4. **Student Keeps Album** → Portable portfolio of demonstrated skills
5. **End of Year** → Album becomes keepsake; export digital copy for records

### Success Metrics

- **Engagement**: 90%+ of students have at least 1 sticker within first 2 weeks
- **Teacher Adoption**: 80%+ of teachers use award function weekly
- **Printing**: <30 seconds per album PDF generation
- **Motivation**: Student survey: "I feel proud of my sticker album" (target: 85% agree)
- **Data Quality**: Sticker awards track to actual skill improvements in digital profile

### Budget & Resources

- **Design**: 40 hours (sticker art + album templates)
- **Development**: 80 hours (backend API + mobile UI + PDF generation)
- **QA & Integration**: 20 hours
- **Physical Production**: School provides color printer + sticker label sheets (~£0.50/student/album)

### Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Sticker printing becomes tedious | Teachers abandon feature | Pre-print batches; provide bulk sticker sheets |
| Sticker inflation (too easy to earn) | Loses motivational value | Define clear rubrics; teacher training on consistency |
| Students lose/damage album | Demotivates | Digital backup; print replacement albums on demand |
| Sticker design prints poorly | Quality perception | Test designs at actual print size; provide high-res files |
| Doesn't improve behavior | ROI question | Track correlation between stickers and progress in other metrics |

---

## Final Notes

**Remember**:
- Perfect is the enemy of shipped
- 80% complete and launched beats 100% complete and never released
- User feedback is more valuable than your assumptions
- Revenue validates product-market fit
- This is a marathon, not a sprint

**When in doubt**:
1. Does this help students learn better? (Mission)
2. Will this generate revenue? (Sustainability)
3. Can this be done in 20% of the time for 80% of the value? (Efficiency)

If yes to all three → Do it.  
If no to any → Deprioritize or skip.

---

**Next Action**: Start Week 1 of Phase 1 → HumanOS Core to 50%

**Let's build this! 🚀**
