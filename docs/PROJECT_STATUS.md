# ESP Thinking Ecosystem — Project Status

**Last Updated**: 2026-05-30
**Current Phase**: ToddlerOS Book 1 production + LAO App Store Review
**Primary Focus**: ToddlerOS booklet markup/content pipeline | LAO 2.2.1 submitted for review

---

## Architecture Overview

### The Three-Layer Model

Every application in this ecosystem uses the same three layers. No product reinvents them — each one is a **lens** into shared infrastructure:

1. **Emotional/Biological** — humanOS + ETP (who is this person, how do they process)
2. **Skills/Capacities** — the universal skills graph (what can humans do, in what order)
3. **Academic/Domain** — subject knowledge (what does this context require)

**Core thesis**: Layers 1 and 2 *facilitate* layer 3. A child who has strong sequencing capacity and whose ETP profile is understood will learn more effectively than one whose foundational skills are assumed and whose neurological wiring is ignored.

### Foundational Tools (not products in their own right)

| Tool | Purpose | Status | Applied In |
|------|---------|--------|------------|
| **CHISG** | Knowledge graph for skills/concepts | Weaviate schema + semantic search working | Skills Map, LAO, PrimaryOS, CareerOS |
| **Skills Map Platform** | Universal skills graph — DAG with 1-5 scoring | Production-ready (full CRUD, auth, admin, visualisation) | All products view subsets via lens queries |
| **humanOS** | AI coaching orchestrator with safeguarding | Partial (orchestrator + LLM + chat UI, in-memory storage) | ETP Profile, coaching tools |
| **ETP Framework** | 12 biological spectra + voltage calculations + response matrix | Domain logic complete (Go), 24 response matrix entries | All products |
| **Pathfinder Engine** | Graph pathfinding with ETP-weighted edges | ⬜ NOT BUILT — later priority | ToddlerOS (theme selection), CareerOS (gap analysis), LAO (revision order), PrimaryOS (next skill) |

### The Skills Map Is Universal

The skills map is not a product — it is the **territory**. Every application is a different **lens** looking at a subset:

| Application | Lens Filter | What It Shows |
|---|---|---|
| **ToddlerOS** | `DevelopmentAge` 0-5, domains: physical, social-emotional, cognitive | Parent activity book themes; "what to practise this week" |
| **PrimaryOS** | `DevelopmentAge` 5-11, all domains | Teacher sticker book; classroom skill tracking |
| **LAO** | Academic domain skills, GCSE-level | Revision pathway; which topics to study next |
| **CareerOS** | All ages, filtered by job-requirement destination nodes | Gap analysis; "distance to competency" for career paths |
| **humanOS** | ETP-modulated edge weights on any subgraph | Coaching plan; "how does this person best develop this skill?" |

### Products

| Product | Purpose | Status | Notes |
|---------|---------|--------|-------|
| **LAO** | GCSE Science revision (mobile) | Approved for App Store distribution | 17 screens, full Go API + SQLite, EAS pipeline. Version 2.2.1 approved on 31 May 2026 with both IAPs (`lao_gcse_science_separate_2`, `lao_science_trilogy_1`) approved. |
| **ToddlerOS** | Parent-facing activity book + ETP for early years | Product design complete, build starting | Activity book (physical+digital), parent observation app, weekly themes. B2C — shortest route to market |
| **PrimaryOS** | Primary-age ETP + sticker book for teachers | Landing page only | Concept defined; no backend. B2B (schools) — longer sales cycle. Shares pathfinder with ToddlerOS |
| **CareerOS** | Skills passport + gap analysis + role matching | Landing page only | Requires pathfinder engine + job-requirement dataset. Rich concept page. |
| **ETP Profile** | Interactive spectrum sliders | Working demo | 8 sliders + voltage calc + action plan generation |

### Concept Demonstrations

| Demo | Purpose | Status | Notes |
|------|---------|--------|-------|
| **DRB** | MAT strategic dashboard concept | Archived concept demo | Partnership route closed. Go+MongoDB backend, React dashboard. 13-school demo remains deployed at espthinking.co.uk/drb as a portfolio piece. |

---

## Current Sprint Status

**Goal**: Clear LAO App Review and move back onto ToddlerOS Book 1 production work.

| Priority | Task | Status | Notes |
|----------|------|--------|-------|
| 1 | ETP Response Matrix | ✅ COMPLETE | 24 response matrix entries (12 spectra × 2 settings) coded in Go |
| 2 | MongoDB Migration (ETP) | ✅ COMPLETE | ETP profiles in MongoDB. Assist Weaviate confirmed local-only. |
| 3 | CHISG Technical Spec | ✅ COMPLETE | Formal spec sent to Fotherby (18 Feb 2026) |
| 4 | CHISG Commercial Assessment | ❌ NO RESPONSE | CHISG classified as academic (internal grounding layer) |
| 5 | ParentOS → ToddlerOS Rebrand | ✅ COMPLETE | All files, routes, cross-references, docs (25 Feb 2026) |
| 6 | Website Deployment | ✅ COMPLETE | ToddlerOS rebrand live on espthinking.co.uk (25 Feb 2026) |
| 7 | ToddlerOS Product Design | ✅ COMPLETE | Activity book model, foundational skills, microdosing pedagogy (27 Feb 2026) |
| 8 | ToddlerOS Book 1 PDF | ✅ COMPLETE | 53-page "Hello, World!" book generated (12 Apr 2026) |
| 9 | CHISG Master Taxonomy | ✅ COMPLETE | 508 skills ingested into Weaviate SkillsMap (12 Apr 2026) |
| 10 | LAO GoLang Course Removal | ✅ COMPLETE | Removed from courseCatalogue, all screen refs cleaned (12 Apr 2026) |
| 11 | LAO RevenueCat Integration | ✅ COMPLETE | SDK wired, offerings configured, production key in place, review flow debugged. |
| 12 | LAO App Store Approval | ✅ APPROVED | Version 2.2.1 approved for App Store distribution on 31 May 2026 after resolving first-IAP attachment workflow. Separate and Trilogy IAPs approved alongside the binary. |
| 13 | DRB Meeting Preparation | ❌ DROPPED | DRB route no longer active. ToddlerOS is the primary path. |
| 14 | ToddlerOS Book 1 markup + content | ⬜ NEXT | Finish booklet production slice: right-page "While You..." content, spread copy, cluster/IDML population. |
| 15 | CHISG Graph Cleanup | ✅ COMPLETE | 142 cycles → 0 cycles. 1828 → 1753 links. All 597 skills with descriptions. DAG confirmed (26 Apr 2026) |
| 16 | CHISG Production Seed | ✅ COMPLETE | Weaviate 8088 seeded: 597 CHISGElement, 1753 SkillLink. /skillstree/explore 500 errors fixed (26 Apr 2026) |
| 17 | skills-api SkillLink limit | ⬜ NEEDS DEPLOY | limit: 1200 → 2000 fix coded. Docker rebuild + push + deploy needed. 23 isolated nodes in explorer until done |
| 18 | Pathfinder Engine | ⬜ LATER | Shared graph pathfinding with ETP-weighted edges — still required, but no longer the immediate slice. |
| 19 | DRB Stakeholder Response | ❌ CLOSED | Partnership route ended. |

### Strategic Decisions

| Date | Decision | Outcome | Rationale |
|------|----------|---------|----------|
| 25 Feb 2026 | CHISG commercial viability | **Academic only** | No response from Fotherby. Internal grounding layer. |
| 25 Feb 2026 | DRB next steps | **Reactivated (Apr 2026)** | Alvin Walters contacted re: Trust Data Manager role. Classifier added as portfolio piece. |
| 01 Apr 2026 | CHISG Classifier as portfolio | **Added** | Load-balanced Go microservice positioned for DRB/MAT context. Live at /classifier-demo |
| 25 Feb 2026 | Next product priority | **ToddlerOS** | B2C, shortest route to market |
| 25 Feb 2026 | PrimaryOS timing | **After ToddlerOS** | B2B to schools = longer sales cycle. Shares components. |
| 27 Feb 2026 | Skills map is universal | **One graph, many lenses** | All products query the same skill graph with different filters (age, domain, ETP). No separate content per product. |
| 27 Feb 2026 | Steiner/Montessori content | **Ingest, don't invent** | Published developmental frameworks are source material. Extract and ingest same as GCSE definitions. Fills the 0-7 base layer. |
| 27 Feb 2026 | ToddlerOS format | **Weekly activity book** | Physical (printed) + digital (PDF). Dip-in/dip-out, not sequential. Free PDF download = acquisition. Printed subscription = revenue. |
| 27 Feb 2026 | Pathfinder is shared engine | **Build once, all products use** | Same algorithm: given current scores + destination + ETP profile → compute optimal path. |
| 27 Feb 2026 | Meta-skillsmap / summits | **Deferred to ESP Thinking / CareerOS** | Apex skills, specialisation timing, decision matrix — important but not ToddlerOS priority. |
| 30 May 2026 | LAO review status | **Submitted and parked** | App version 2.2.1 submitted. Next execution focus returns to ToddlerOS production assets. |
| 31 May 2026 | LAO App Store approval | **Approved for distribution** | iOS app version 2.2.1 approved by App Review. Both IAPs approved. LAO moves from review-risk to live-distribution/marketing task set. |

---

## ToddlerOS Product Design (27 Feb 2026)

### Vision
ToddlerOS replaces generational memory with structured observation prompts. Parents who didn't have engaged role models get a framework for interacting with their children that naturally builds the foundational capacities that predict school readiness and life outcomes.

**Tagline**: "The Coding Lesson You Never Knew You Needed"

### Core Insight: Microdosing the Poles
Toddlers' ETP defaults haven't locked in yet (0-5 is the neuroplasticity sweet spot). Instead of teaching *about* spectra, let children **experience both poles** in tiny, safe doses:
- Loud and quiet (Energy Directionality + Voltage Sensitivity)
- Near and far (Social Gravity + Risk Tolerance)
- Mine and yours (Care Response + Integrity Logic)
- Yes and no (Threat Response + Integrity Logic)

Each spectrum = one theme. 12 spectra × 3 difficulty cycles = 36 themed weeks.

### The Activity Book Format

Two modes per theme spread:

| Mode | Duration | Parent Role | Purpose |
|---|---|---|---|
| **"With You"** | 3-5 min | Active participant | High-leverage interaction (joint attention, turn-taking, emotional labelling) |
| **"While You..."** | 5-15 min | Nearby but freed | Child engages independently (colouring, stickers, dot-to-dot, audio story via QR) |

**Design principles**:
- No prep required — use whatever's to hand
- No sequence — any page is the right page
- No screen for the child — QR links to audio only, phone goes face-down
- Progress rosette, not checklist — stickers fill a wheel, no "correct" order
- Explicit permission to skip — "Do one page this week. Or five. Or none."
- Seasonal editions (Spring/Summer/Autumn/Winter, 8 themes each = 16 spreads per book)

### Distribution Model

| Channel | Cost | Price | Purpose |
|---|---|---|---|
| **PDF download** (free) | Hosting only | Free | Email capture, app signup funnel |
| **Printed book** (card stock + sticker sheet) | ~£3-4 POD | £7.99 + postage | Physical product parents want |
| **Annual subscription** (4 seasonal books) | ~£12-16/year | £24.99/year | Recurring revenue |
| **App companion** (observation log + progress) | Dev time | Free with book / £1.99/mo | Parent-only — child's interface is paper |

### The 7 High-Leverage Foundational Capacities

These predict downstream outcomes across all the research (Hart & Risley, Heckman, marshmallow replication studies). They are the **skill nodes at the base of the graph**:

| Capacity | Why It Matters | Time to Practise |
|---|---|---|
| **Inhibitory control** | Predicts academic readiness, social competence, health outcomes | 2 min (stop-go games) |
| **Joint attention** | Gateway to language, social referencing, instruction-based learning | 1 min (pointing, shared book) |
| **Emotional labelling** | Naming feelings builds prefrontal regulation circuit | 30 sec ("You look frustrated") |
| **Turn-taking** | Foundation of conversation, negotiation, cooperation | 2 min (rolling ball back and forth) |
| **Sustained attention** | Predicts literacy readiness | 0 min (protect it by not interrupting) |
| **Verbal interaction density** | Responsive conversation matters more than vocabulary drills | 5 min woven into any activity |
| **Frustration tolerance** | Real mechanism behind persistence | 1 min (wait before helping) |

### Observation as Byproduct

The parent doesn't score 1-5 on "Social Gravity." They answer after each week's theme:
> "Which was easier for them — near or far?"

That maps to an ETP data point. Over 32 weeks, a profile builds without the parent ever knowing they're creating one. The assessment happens through play, not questionnaire.

### Social Mobility Mechanism

The demographic achievement gap is largely a proxy for responsive parent interaction time. ToddlerOS addresses this by:
1. Making each interaction 5 minutes or less (compliance research: 30 min gets abandoned, 5 min gets done)
2. Zero prep, zero materials, fits into existing routines
3. One observation prompt per activity — parent starts noticing without the prompt after ~15 uses
4. Free PDF download removes financial barrier

**Goal**: ToddlerOS has done its job when the parent no longer needs it.

### Kevin Avison / Steiner Connection

The skills map was originally developed during discussions with Kevin Avison of the Steiner Waldorf Advisory Service as a mechanism for showing student progress for Ofsted inspections. The schools themselves rejected it — the framing was wrong ("measure progress for Ofsted" vs "document what you already see").

The `Criteria1-5` fields on `skills_key` are the exact 1-5 observation scale discussed. Steiner-compatible translations: 1=Not yet observed / 2=Glimmers / 3=With support / 4=Independent / 5=Shows others.

Steiner and Montessori are not target markets — they are **source material** for identifying foundational capacity nodes that academic curricula assume but never name (sensorimotor, executive function, social-emotional, metacognitive, physical coordination as cognitive prerequisite).

---

## Completed Milestones

### Production Deployment ✅ (Feb 2026)
- [x] All 7 Docker images built for linux/amd64 and pushed to Docker Hub (`mike5tew/`)
- [x] Vultr server (`192.248.151.185`) running 10 containers via `docker-compose.prod.yml`
- [x] Domain `espthinking.co.uk` live with Let's Encrypt SSL (expires 2026-05-15)
- [x] Nginx reverse proxy serving: portfolio, skills map, DRB, all APIs
- [x] Visitor analytics system with MongoDB + bearer token auth
- [x] All landing pages verified returning 200 over HTTPS

### Portfolio Site ✅ (Feb 2026)
- [x] 9 landing pages (Portfolio, ETP, PrimaryOS, LAO, ESP World, CHISG, ToddlerOS, CareerOS, Classifier Demo)
- [x] ContactReveal email forms on all landing CTAs (protected from harvesting)
- [x] Live interactive tools: ETP Profile, Semantic Query, AI Chat
- [x] Analytics dashboard at `/analytics` with period selector and event tracking

### LAO Mobile ✅ (Jan 2026)
- [x] Expo/React Native app with 17 screens
- [x] SQLite-backed Go API with 11 handler modules
- [x] Audio summaries, speed reader, flashcards, progress tracking
- [x] EAS build pipeline → TestFlight distribution

### Skills Map Platform ✅
- [x] Full CRUD API (Go) with JWT auth and admin roles
- [x] Graph explorer, skills tree, course management
- [x] Weaviate CHISG integration (579 elements, 1,120 links, 27 courses)
- [x] MySQL operational database for assignments and markbook scores
- [x] Production deployed behind main proxy

### ETP Domain Logic ✅
- [x] 9 spectra definitions with trainable skills (`spectra.go`, 347 lines)
- [x] Voltage calculator with real impact calculations (`voltage_calculator.go`, 318 lines)
- [x] Compatibility solutions for all 8 spectra (`compatibility.go`, 192 lines)
- [x] 3 cognitive modes, educational enigmas, identity hijacking frameworks
- [x] Weaviate ETPProfile schema (660 lines)
- [x] Interactive frontend with sliders, undo/redo, help dialog

### ETP Response Matrix ✅ (Feb 2026)
- [x] 18 response matrix entries (9 spectra × 2 settings) encoded in `response_matrix.go`
- [x] Three barrier types mapped: verbalising, starting, mistakes
- [x] Adult language shifts: avoid/use pairs for each spectrum setting
- [x] Play-first strategies for each setting
- [x] In-memory response matrix (zero database dependency for plan generation)
- [x] Action plan endpoint (`POST /api/actionplan`) returns personalised interventions
- [x] Coach handler personalisation (`personaliseWithETP()`) uses response matrix for live coaching
- [x] Unified language shifts and Play-First Principle constants exported

### MongoDB Migration (ETP Profiles) ✅ (Feb 2026)
- [x] ETP profile storage migrated from Weaviate to MongoDB (`etp_profiles` collection in `esp_organizer` database)
- [x] `getETPDB()` lazy singleton follows same pattern as `getAnalyticsDB()`
- [x] Profiles stored as proper BSON documents: `userId`, `spectrumProfile` (map), `spectrumVector`, `timestamp`
- [x] Non-blocking save: if MongoDB unavailable, plan still generates from in-memory response matrix
- [x] Coach handler queries MongoDB for most recent profile by userId
- [x] `docker-compose.prod.yml` updated: removed `WEAVIATE_URL` from assist-api, removed Weaviate from depends_on
- [x] Deployed and verified on Vultr — MongoDB profile ID returned in API response

### Semantic Link Extraction Tool (70% complete)
- [x] Backend API endpoints (extract, search, validate) + Weaviate integration
- [x] React `<SemanticLinkExtractor />` with PDF viewer and 3-click workflow
- [ ] LLM-as-judge validation handler
- [ ] Hierarchy inference logic
- [ ] Batch processing pipeline

### DRB Concept Demo ✅
- [x] Go + MongoDB backend with CRUD routes (schools, finance, estates, contractors, keys)
- [x] React dashboard + estates monitor
- [x] Seeded with representative demo data for 13 schools

---

## What's NOT Built Yet (Honest Assessment)

These are described on landing pages or in product design but have **no implementation behind them**:

| Claim | Where Promised | Reality |
|-------|---------------|---------|
| ToddlerOS Activity Book | Product design (this doc) | Design complete. No PDF generator, no content pipeline, no observation app. |
| ToddlerOS Parent App | Product design (this doc) | No app. No observation logging. No profile building. |
| Foundational Skill Nodes (0-7 age) | Product design (this doc) | Zero nodes in graph. Need ~30-50 from developmental science sources. |
| Pathfinder Engine | Architecture design (this doc) | Not built. Required by all products for "what comes next" queries. |
| CareerOS Role Architect | `/careeros` landing | Static example data only; no backend. Requires pathfinder + job-requirement dataset. |
| CareerOS Skills Passport | `/careeros` landing | No backend, no data model |
| CareerOS Decision Matrix | Product design (this doc) | Variables defined (Passion, Aptitude, Demand, Longevity, Fit, Stage, Cost). No implementation. |
| PrimaryOS Teacher Dashboard | `/primary-os` landing | No backend; hero CTAs correctly commented out |
| PrimaryOS Sticker Book | `/primary-os` landing | Framework specified in docs but not implemented |
| PrimaryOS Neuron Navigators Journal | `/primary-os` landing | Concept only |
| ToddlerOS Spectrum Explorer | `/toddler-os/explorer` | Coming Soon page with ContactReveal |
| ToddlerOS Parent Guide | `/toddler-os/guide` | Coming Soon page with ContactReveal |
| humanOS persistent profiles | humanOS standalone | In-memory only (`map[string]*StudentProfile`) |
| Meta-skills map / Summits | Strategic discussion (this doc) | Concept defined. Apex skills, specialisation timing, temporal weighting. Deferred to CareerOS/ESP Thinking. |

---

## Next Steps

### Immediate Priority: Finish ToddlerOS Book 1 Production Assets

**Step 1: Write the missing "While You..." pages**
- `docs/TODDLEROS_BOOKLET_OUTLINES.md` explicitly marks the right-hand child-independent pages as not yet written for all themed weeks
- Start with **Edition 1 / Hello World / Cycle 1** only: 12 themes, one right-page concept per spread
- Keep each one 5-15 minutes, zero parent input, thematically paired to the left-page "With You" activity
- Output format should be final spread-ready copy, not abstract notes

**Step 2: Convert spread copy into the print pipeline inputs**
- Use the existing cluster/IDML tooling in `ToddlerOS/scripts/populate-cluster.py`
- Confirm the Cycle 1 spread data for Book 1 is complete enough to populate all 12 themes cleanly
- Keep the work bounded to Book 1 rather than reopening all 36 themed weeks

**Step 3: Populate and review the printable file**
- Run the IDML population workflow for Book 1 and inspect the generated spreads for overflow, weak copy, and pairing issues
- Fix copy/layout mismatches at the markup/content level before adding new engine work

**Step 4: Resume engine work after Book 1 is materially complete**
- Foundational skill nodes and the shared pathfinder still matter
- They are now second-order work behind having a finished ToddlerOS saleable booklet asset

### Parallel: LAO Review Watch
- Wait for App Review outcome on 2.2.1
- Do not reopen LAO engineering work unless review feedback requires it

### Later: PrimaryOS MVP
- Reuses foundational skill nodes + pathfinder engine
- Adds teacher dashboard + sticker book
- B2B school partnerships
- Teacher scores skills 1-5 (Criteria1-5 with Steiner-compatible language)

### Later: CareerOS MVP
- Requires pathfinder engine + job-requirement destination node sets
- Map 10-20 careers as skill clusters in the graph
- Decision matrix interface (Passion × Aptitude × Demand × Longevity × ETP Fit)
- "Distance to competency" metric: not "do you qualify?" but "how far, and what's the fastest path?"

### Deferred: Meta-Skills Map / Summits
- Apex skills as emergent properties of skill combinations (Clarity, Mastery, Creativity, Connection, Contribution, Contentment, Integration)
- Temporal weighting function: breadth (0-12) → depth (12-25) → integration (25-40) → wisdom (40+)
- Specialisation timing model
- Relevant to CareerOS and ESP Thinking high school products

---

## Infrastructure Reference

| Service | Container | Port | Database |
|---------|-----------|------|----------|
| Portfolio + assist API | assist-frontend / assist-api | 80 / 8080 | MongoDB (`esp_organizer`, `esp_analytics`) |
| Skills Map | skills-frontend / skills-api | 80 / 8080 | MySQL (`skills_db`), Weaviate |
| DRB | drb-frontend / drb-api | 80 / 8082 | MongoDB (`drb_monitor`) |
| Weaviate | weaviate | 8080 (internal) | Used by skills-api only (CHISG semantic search) |
| Main Proxy | main-proxy | 80 / 443 | — |

**Note (Feb 2026)**: assist-api no longer depends on Weaviate. ETP profiles stored in MongoDB. The in-memory response matrix (`response_matrix.go`) has zero database dependency — all 18 entries live in compiled Go code.

**Docker Hub**: `mike5tew/` — main-proxy, assist-frontend, assist-api, skills-frontend, skills-api, drb-frontend, drb-api
**Domain**: espthinking.co.uk → 192.248.151.185 (Vultr)
**SSL**: Let's Encrypt via certbot container (ssl profile)

---

## Documentation Map

| Document | Purpose |
|----------|---------|
| `PROJECT_STATUS.md` | This file — current state and priorities |
| `PROJECT_INDEX.md` | Architecture, data systems, file locations |
| `PORTFOLIO_SITE_OVERVIEW.md` | Live site routes, CTAs, integrations |
| `STRATEGIC_PLAYBOOK.md` | Pitch/vision strategy |
| `SEMANTIC_SEARCH_APPROACH.md` | Weaviate integration design |
| `ETP_VOLTAGE_COMPATIBILITY_FRAMEWORK.md` | ETP domain logic reference |
| `FINANCIAL_ESTATES_SYSTEM.md` | DRB concept definition |
| `STICKER_ALBUM_IMPLEMENTATION_SUMMARY.md` | SkillsMarkBook sticker album spec |
