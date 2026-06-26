# ESP Thinking Portfolio Site — Current Overview (assist/frontend)

_Last updated: 2026-02-27_

This document describes the current **ESP Thinking portfolio site** implemented in `assist/frontend`. It is written so that Weaviate-assist can index an accurate, up-to-date view of the portfolio, its routes, and how it connects to the wider ecosystem.

---

## 1. High-Level Purpose

The portfolio frontend is a **marketing and demo hub** for the ESP ecosystem. It:

- Presents product-specific landing pages (PrimaryOS, LAO, ToddlerOS, CareerOS, ETP, CHISG, ESP World)
- Provides deep dives into core ideas (CHISG, ETP, HumanOS framing)
- Exposes live demos and tools backed by the assist API and Weaviate
- Bridges out to other apps (Skills Map / skillstree, DRB, etc.) via Nginx routing

Code location:

- React app: `assist/frontend/`
- Main entry: `assist/frontend/src/App.tsx`
- Shared navigation: `assist/frontend/src/components/NavHeader.tsx`

---

## 2. Top-Level Routing & Layout

Routing is defined in `frontend/src/App.tsx` using `react-router-dom`.

- The router `basename` is:
  - `REACT_APP_BASENAME` if set, otherwise
  - `"/esp-organizer"` when the pathname starts with `/esp-organizer`, else `/`.
- Most pages share a common layout with `NavHeader`, but some landing pages own their own hero and load **without** the header wrapper.

### 2.1 Routes Without Shared Header

These routes render full-bleed marketing/landing experiences:

- `/` → PortfolioLanding
  - Master portfolio overview tying together all ESP products.
- `/etp-landing` → ETPLanding
  - Explains ETP (Emergent Tendency Profiles) and the 9 spectra.
- `/primary-os` → PrimaryOSLanding
  - "Operating system for growing learners"; Neuron Navigators Guide; school-facing CTA.
- `/lao` → LAOLanding
  - LAO GCSE revision experience; focuses on low-friction, sticker-album style revision.
- `/esp-world` → ESPWorldLanding
  - Umbrella view of ESP World tools and pathways.
- `/chisg` → CHISGLanding
  - Information quality / CHISG knowledge-graph framing.
- `/toddler-os` → ToddlerOSLanding
  - ToddlerOS framing — ETP spectra for parents of young children. Tagline: "The Coding Lesson You Never Knew You Needed".
  - Legacy `/parent-os` route redirects here.
- `/careeros/*` → CareerOSLanding (internal nested routes)
  - CareerOS/skills passport experience; CTA and sub-routes managed internally.
- `/skillstree` → SkillsTreeRedirect
  - Redirects to the separate Skills Map app hosted under `/skillstree/` via the main proxy.

### 2.2 Routes With Shared Header (AppLayout)

All other routes are nested under `AppLayout`, which provides the `NavHeader` and a scrolling content area.

Key demo/tool routes:

- `/etp-profile` → ETPProfilePage
  - Interactive ETP profile demo (sliders, text explanations).
- `/esp-world/demo` → ESPWorldDemo
  - Demonstrates a cross-product ESP World flow.
- `/esp-world/lesson-demo` → LessonPlanningDemo
  - Lesson-planning demo connected to CHISG/ETP ideas.
- `/toddler-os/explorer` → ToddlerOSExplorer
- `/toddler-os/guide` → ToddlerOSGuide
- `/coach-mvp-demo` → CoachMVPDemo

Core tools backed by assist API + Weaviate:

- `/tools` → Home (tool index)
- `/semantic-query` → SemanticQuery
  - Weaviate-backed semantic search over docs/skills.
- `/ai-chat` → AIChat
  - Chat assistant over assist/CHISG context.
- `/contentLoader` → ContentLoader
  - Internal content ingestion and testing.
- `/semantic-links/extract` → SemanticLinkExtractor
  - Frontend UI for generating SemanticLink relationships from PDFs.

Admin & diagnostics:

- `/admin/*` → AdminPanel
- `/diagnostics` → DiagnosticViewer
- `/upload/immunology` → ImmunologyUpload
- `*` → ErrorPage (fallback for unknown routes).

---

## 3. Skills Map / skillstree Integration

The portfolio does **not** host the skills tree directly; instead it hands off to the separate `skills-map-platform` frontend.

- Route: `/skillstree`
- Component: `SkillsTreeRedirect` in `App.tsx`
- Behaviour:
  - On mount, computes `target = protocol + "//" + hostname + "/skillstree/"`.
  - If `window.location.href !== target`, calls `window.location.replace(target)`.
  - This ensures that when running locally on `http://localhost:3000`, the user is redirected to `http://localhost/skillstree/` (served by the main Nginx proxy, which fronts `skills-frontend`).

This design avoids redirect loops and keeps `/skillstree/` owned by the skills-map application while still exposing it as a first-class item in the portfolio navigation.

---

## 4. Product Landing Pages & CTAs

Each landing page is opinionated about its audience and call-to-action, but pattern-aligned.

### 4.1 PrimaryOSLanding (`/primary-os`)

- Focus: classroom behaviour → strategy; Neuron Navigators Guide; sticker-book journey.
- Content sections:
  - Hero with PrimaryOS branding and summary.
  - Explanation of moving from shame to agency.
  - Two kinds of skills: spectrum skills vs foundational skills.
  - 9 ETP spectra cards (with classroom examples).
  - Neuron Navigators Guide Book stages (Reception → Year 6+).
  - Teacher role shift: judge → guide.
  - ToddlerOS bridge section linking to `/toddler-os`.
  - Final CTA section.
- CTAs:
  - Bottom CTA uses `ContactReveal` to open a protected email form to `world@espthinking.co.uk` with label **"Request School Demo"**.
  - ToddlerOS bridge CTA links to `/toddler-os`.
  - Hero CTAs (Teacher Dashboard, Neuron Navigators) are correctly commented out — routes don't exist yet.
- **Implementation status**: Landing page only. No backend, no teacher dashboard, no sticker book. Components described here would be reusable across the ecosystem if built.

### 4.2 LAOLanding (`/lao`)

- Focus: LAO GCSE revision experience.
- Hero: marketing copy and visuals tailored to revision resistance and low-friction engagement.
- CTA: primary buttons link directly to the public App Store listing for Little and Often LAO.

### 4.3 ETPLanding (`/etp-landing`)

- Focus: ETP spectra, sensitivity profiles, underlying neuroscience.
- Explains core spectra and moderators.
- CTA: uses `ContactReveal` with a label appropriate for ETP-focused enquiries.

### 4.4 ToddlerOSLanding (`/toddler-os`)

- Focus: parent-facing ETP framework and activity book for early years. "The Coding Lesson You Never Knew You Needed."
- Rebranded from ParentOS (Feb 2026). Legacy `/parent-os` routes redirect.
- **Product vision (Feb 2026)**: Weekly activity book (physical + PDF) with two modes per theme:
  - "With You" (3-5 min parent interaction — builds foundational capacities)
  - "While You..." (5-15 min independent activity — colouring, stickers, audio via QR)
  - Observation prompts build ETP profile as byproduct of play (not questionnaire)
  - Free PDF download → email capture. Printed subscription → revenue.
  - See `PROJECT_STATUS.md` for full product design.
- Links:
  - `/toddler-os/explorer` — "Coming Soon" page with description and ContactReveal ("Register Interest").
  - `/toddler-os/guide` — "Coming Soon" page with description and ContactReveal ("Register Interest").
- CTA: includes `ContactReveal` for parent/partner contact.
- **Implementation status**: Landing page only. Product design complete (activity book model, foundational skills, microdosing pedagogy). No activity book generator, no parent app, no PDF pipeline yet.

### 4.5 CareerOSLanding (`/careeros/*`)

- Focus: verified skills passport / CareerOS view.
- Implemented as a nested router inside `CareerOSLanding`.
- Exposes internal routes under `/careeros/...` (e.g. demo, enterprise, passport), all fronted by the main portfolio router at `/careeros/*`.
- **Product vision (Feb 2026)**: Gap analysis powered by pathfinder engine. "Distance to competency" metric: not "do you qualify?" but "how far, and what's the fastest path?" Decision matrix for skills direction using Passion × Aptitude × Demand × Longevity × ETP Fit. See `PROJECT_STATUS.md` for details.
- **Implementation status**: Landing page only (948 lines of rich concept content with static example data). No backend. Requires pathfinder engine + job-requirement dataset (both designed, neither built).

### 4.6 ESPWorldLanding (`/esp-world`)

- Focus: integrated ESP World experience spanning LAO, PrimaryOS, ToddlerOS, and Skills Map.
- Provides entry points into:
  - `ETPProfilePage` (`/etp-profile`)
  - lesson/ESP World demos
  - ToddlerOS explorer/guide

### 4.7 CHISGLanding (`/chisg`)

- Focus: CHISG as information-quality backbone.
- Explains how skills, links, and semantic enrichment combine into a robust knowledge graph for education.

All major landings share the same pattern: storytelling + visuals + concrete classroom/parent examples + clear CTA, with protected email via `ContactReveal` to avoid direct email harvesting.

---

## 5. Navigation & Header

`NavHeader` defines the primary navigation menu items exposed in the top bar when `AppLayout` is active, including (but not limited to):

- Portfolio home
- Product landings (PrimaryOS, LAO, ToddlerOS, CareerOS, ESP World, CHISG)
- Tools (semantic query, AI chat, content loader, semantic link extractor)
- Skills Map (via `/skillstree` → redirect)

Some landings (e.g. `/`, `/primary-os`, `/lao`, `/etp-landing`) intentionally **do not** use `AppLayout` so they can control the hero and header experience directly.

---

## 6. Backends & Data Dependencies

The portfolio frontend is primarily a **read-only + demo** layer. Implementation status varies significantly by product.

### Working backends:
- Semantic tools (`/semantic-query`, `/ai-chat`, `/contentLoader`, `/semantic-links/extract`):
  - Use the assist API (`esp-organizer`) and the `weaviate-assist` instance to search and manipulate docs/semantic links.
- ETP Profile (`/etp-profile`):
  - Full domain logic in Go (spectra, voltage calculation, compatibility). Weaviate ETPProfile schema defined.
- Skills Map hand-off (`/skillstree` → redirect):
  - Full production stack (Go API + MySQL + Weaviate CHISG instance) running behind the main proxy.
- DRB (`/drb/` via nginx):
  - Go + MongoDB backend with seeded demo data. Concept demonstration for MATs, not a standalone product.
- Analytics (`/api/analytics/*`):
  - Page view + event tracking, MongoDB storage, bearer token auth.

### Concept-only (landing pages with no backend):
- PrimaryOS (`/primary-os`): No backend, no teacher dashboard, no sticker book implementation.
- ToddlerOS (`/toddler-os`, `/toddler-os/explorer`, `/toddler-os/guide`): Coming Soon pages with ContactReveal.
- CareerOS (`/careeros/*`): Static concept content with hardcoded example data. No Skills Passport, Role Architect, or Matching backend.

### Foundational tools (applied across products, not standalone):
- CHISG knowledge graph: Weaviate-backed, used by Skills Map, Semantic Query, and AI Chat.
- humanOS coaching: Standalone project with orchestrator + LLM integration but in-memory student profiles.

---

## 7. Current Deployment

### Production (Live — espthinking.co.uk)

- **Domain**: `espthinking.co.uk` (DNS on IONOS → `192.248.151.185`)
- **SSL**: Let's Encrypt certificate via certbot container (expires 2026-05-15)
- **Server**: Vultr VPS running `docker-compose.prod.yml` with 10 containers
- **Docker Hub**: All images under `mike5tew/` namespace (linux/amd64)

Container routing via `main-proxy` Nginx:
- `/` → `assist-frontend` (portfolio React app)
- `/api/` → `assist-api` (esp-organizer Go API, port 8080)
- `/skillstree/` → `skills-frontend` (Skills Map React app)
- `/skillstree/api/` → `skills-api` (Skills Map Go API)
- `/drb/` → `drb-frontend` (DRB concept demo)
- `/drb/api/` → `drb-api` (DRB Go API, port 8082)

Data services:
- MongoDB (`esp_organizer`, `esp_analytics`, `drb_monitor`)
- MySQL `skills_db` (Skills Map operational data)
- Weaviate v1.24.10 (CHISG knowledge graph)

### Visitor Analytics

- Tracking: page views + CTA events (contact_unlocked, contact_sent)
- Privacy: daily hashed fingerprint (SHA256 of IP+UA), auto-deleted after 365 days
- Stats endpoint: `/api/analytics/stats` (bearer token auth)
- Dashboard: `/analytics` (period selector, summary cards, charts)

### Local / Dev

- Run via Docker Compose: `make up` or `docker compose up`
- Portfolio at `http://localhost:3000` (dev) or `http://localhost/` (via proxy)
- Skills Map at `http://localhost/skillstree/`

---

## 8. How This Document Is Intended to Be Used

- **For Weaviate-assist**: as a single, up-to-date source describing the portfolio site structure, so semantic search over project docs returns accurate answers about routes, CTAs, and integrations.
- **For deployment planning**: as the frontend reference when mapping routes and subpaths in Nginx/Traefik on Vultr.
- **For collaborators**: as an orientation guide before editing landing pages or wiring new demos/tools into the portfolio shell.
