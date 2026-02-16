# ESP Thinking Ecosystem — Project Status

**Last Updated**: 2026-02-16
**Current Phase**: Infrastructure Complete — Assessing Commercial Viability
**Primary Focus**: CHISG Commercial Assessment + ETP Response Matrix Integration + LAO Testing

---

## Architecture Overview

This ecosystem has three layers:

1. **Foundational Tools** — reusable services applied across products
2. **Products** — user-facing applications built on those tools
3. **Concept Demonstrations** — prototypes used for stakeholder engagement

### Foundational Tools (not products in their own right)

| Tool | Purpose | Status | Applied In |
|------|---------|--------|------------|
| **CHISG** | Knowledge graph for skills/concepts | Weaviate schema + semantic search working | Skills Map, LAO, PrimaryOS, CareerOS |
| **Skills Map Platform** | Skills tree, course maps, student tracking | Production-ready (full CRUD, auth, admin) | Portfolio tools, PrimaryOS, classroom tools |
| **humanOS** | AI coaching orchestrator with safeguarding | Partial (orchestrator + LLM + chat UI, in-memory storage) | ETP Profile, coaching tools |
| **ETP Framework** | 9 biological spectra + voltage calculations | Domain logic complete (Go), Weaviate schema defined | All products |

### Products

| Product | Purpose | Status | Notes |
|---------|---------|--------|-------|
| **LAO** | GCSE Science revision (mobile) | Testing in TestFlight | 17 screens, full Go API + SQLite, EAS pipeline |
| **PrimaryOS** | Primary-age ETP + Neuron Navigators | Landing page only | Concept defined; no backend. Components would feed other products |
| **ParentOS** | Parent-facing spectrum explorer | Landing page only | Explorer/Guide pages are "Coming Soon" shells |
| **CareerOS** | Skills passport + role matching | Landing page only | Rich concept page, zero backend |
| **ETP Profile** | Interactive spectrum sliders | Working demo | 8 sliders + voltage calc + action plan generation |

### Concept Demonstrations

| Demo | Purpose | Status | Notes |
|------|---------|--------|-------|
| **DRB** | MAT strategic dashboard concept | Demo-functional | Shows MATs what's possible with finance/estates/governance. Seeded demo data, no real data ingestion |

---

## Current Sprint Status

**Goal**: Assess CHISG commercial viability; continue LAO field testing.

| Priority | Task | Status | Notes |
|----------|------|--------|-------|
| 1 | ETP Response Matrix | ✅ COMPLETE | 18 response matrix entries (9 spectra × 2 settings) coded in Go, deployed to Vultr, MongoDB persistence |
| 2 | MongoDB Migration (ETP) | ✅ COMPLETE | ETP profiles migrated from Weaviate → MongoDB. Assist Weaviate confirmed local-only. |
| 3 | CHISG Technical Spec | 🔄 IN PROGRESS | Technical outline for stakeholders covering hallucination reduction, analogy detection, gap analysis |
| 4 | CHISG Commercial Assessment | 🔄 IN PROGRESS | Preparing tech outline for Toby Fotherby to assess uniqueness/viability |
| 5 | LAO Field Testing | 🔄 IN PROGRESS | TestFlight build live; collecting usage feedback |
| 6 | PrimaryOS or ParentOS MVP | ⬜ CONSIDERING | PrimaryOS components reusable across ecosystem; ParentOS may be more commercial |

**Key Decision Pending**: Whether CHISG (as a knowledge-graph methodology) has standalone commercial value, which would change the development priority order significantly.

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
- [x] 8 landing pages (Portfolio, ETP, PrimaryOS, LAO, ESP World, CHISG, ParentOS, CareerOS)
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

These are described on landing pages but have **no implementation behind them**:

| Claim | Where Promised | Reality |
|-------|---------------|---------|
| CareerOS Role Architect | `/careeros` landing | Static example data only; no backend |
| CareerOS Skills Passport | `/careeros` landing | No backend, no data model |
| CareerOS Intelligent Matching | `/careeros` landing | No backend |
| PrimaryOS Teacher Dashboard | `/primary-os` landing | No backend; hero CTAs correctly commented out |
| PrimaryOS Sticker Book | `/primary-os` landing | Framework specified in docs but not implemented |
| PrimaryOS Neuron Navigators Journal | `/primary-os` landing | Concept only |
| ParentOS Spectrum Explorer | `/parent-os/explorer` | Coming Soon page with ContactReveal |
| ParentOS Parent Guide | `/parent-os/guide` | Coming Soon page with ContactReveal |
| humanOS persistent profiles | humanOS standalone | In-memory only (`map[string]*StudentProfile`) |

---

## Next Steps (Decision Tree)

```
If CHISG has commercial value (Fotherby assessment):
  → Build full CHISG demo with GCSE Science curriculum mapping
  → This becomes the core IP to protect/license

If PrimaryOS next:
  → Components (spectrum sliders, sticker tracking) reuse across ecosystem
  → Needs: backend API, teacher dashboard, student-facing sticker UI
  → Estimated: weeks of work

If ParentOS next:
  → Potentially more commercial (B2C vs B2B)
  → Simpler scope (spectrum explorer + guides)
  → Could be built as a layer on top of ETP Profile + CHISG
```

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
