# ESP Thinking Portfolio Site — Current Overview (assist/frontend)

_Last updated: 2026-02-14_

This document describes the current **ESP Thinking portfolio site** implemented in `assist/frontend`. It is written so that Weaviate-assist can index an accurate, up-to-date view of the portfolio, its routes, and how it connects to the wider ecosystem.

---

## 1. High-Level Purpose

The portfolio frontend is a **marketing and demo hub** for the ESP ecosystem. It:

- Presents product-specific landing pages (PrimaryOS, LAO, ParentOS, CareerOS, ETP, CHISG, ESP World)
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
- `/parent-os` → ParentOSLanding
  - ParentOS framing and connection to PrimaryOS.
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
- `/parent-os/explorer` → ParentOSExplorer
- `/parent-os/guide` → ParentOSGuide
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
  - ParentOS bridge section linking to `/parent-os`.
  - Final CTA section.
- CTAs:
  - Bottom CTA uses `ContactReveal` to open a protected email form to `world@espthinking.co.uk` with label **"Request School Demo"**.
  - ParentOS bridge CTA links to `/parent-os`.

### 4.2 LAOLanding (`/lao`)

- Focus: LAO GCSE revision experience.
- Hero: marketing copy and visuals tailored to revision resistance and low-friction engagement.
- CTA: bottom-of-page button wired through `ContactReveal` for protected contact with `world@espthinking.co.uk`.

### 4.3 ETPLanding (`/etp-landing`)

- Focus: ETP spectra, sensitivity profiles, underlying neuroscience.
- Explains core spectra and moderators.
- CTA: uses `ContactReveal` with a label appropriate for ETP-focused enquiries.

### 4.4 ParentOSLanding (`/parent-os`)

- Focus: parent-facing framing of the same CHISG/ETP language used in school.
- Links:
  - `/parent-os/explorer` for interactive exploration.
  - `/parent-os/guide` for deeper written guidance.
- CTA: includes `ContactReveal` for parent/partner contact.

### 4.5 CareerOSLanding (`/careeros/*`)

- Focus: verified skills passport / CareerOS view.
- Implemented as a nested router inside `CareerOSLanding`.
- Exposes internal routes under `/careeros/...` (e.g. demo, enterprise, passport), all fronted by the main portfolio router at `/careeros/*`.

### 4.6 ESPWorldLanding (`/esp-world`)

- Focus: integrated ESP World experience spanning LAO, PrimaryOS, ParentOS, and Skills Map.
- Provides entry points into:
  - `ETPProfilePage` (`/etp-profile`)
  - lesson/ESP World demos
  - ParentOS explorer/guide

### 4.7 CHISGLanding (`/chisg`)

- Focus: CHISG as information-quality backbone.
- Explains how skills, links, and semantic enrichment combine into a robust knowledge graph for education.

All major landings share the same pattern: storytelling + visuals + concrete classroom/parent examples + clear CTA, with protected email via `ContactReveal` to avoid direct email harvesting.

---

## 5. Navigation & Header

`NavHeader` defines the primary navigation menu items exposed in the top bar when `AppLayout` is active, including (but not limited to):

- Portfolio home
- Product landings (PrimaryOS, LAO, ParentOS, CareerOS, ESP World, CHISG)
- Tools (semantic query, AI chat, content loader, semantic link extractor)
- Skills Map (via `/skillstree` → redirect)

Some landings (e.g. `/`, `/primary-os`, `/lao`, `/etp-landing`) intentionally **do not** use `AppLayout` so they can control the hero and header experience directly.

---

## 6. Backends & Data Dependencies

The portfolio frontend is primarily a **read-only + demo** layer, but several routes depend on backend services:

- Semantic tools (`/semantic-query`, `/ai-chat`, `/contentLoader`, `/semantic-links/extract`):
  - Use the assist API (`esp-organizer`) and the `weaviate-assist` instance to search and manipulate docs/semantic links.
- ParentOS, ESP World, and CHISG content:
  - Read from existing docs and CHISG/ETP definitions; some demos are stubbed for now but wired for future data integration.
- Skills Map hand-off:
  - Depends on `skills-map-platform` stack (Go API + MySQL + Weaviate CHISG instance) running behind the main proxy.

---

## 7. Current Deployment Assumptions

Local / dev:

- Run via Docker Compose in the `assist` repo (`make up` or equivalent commands).
- Portfolio React app is served by `assist-frontend` container.
- For dev convenience, `assist-frontend` can be exposed on `http://localhost:3000` while the main Nginx proxy serves `http://localhost/`.
- The `/skillstree/` path is served by `skills-frontend` (skills-map-platform) behind the `main-proxy` Nginx container.

Production (planned via Vultr docs):

- `assist-frontend` is fronted by a main Nginx proxy container that also forwards:
  - `/skillstree/` → `skills-frontend`
  - `/skillstree/api/` → `skills-api`
  - Additional apps (e.g. DRB) under subpaths.
- Domain mapping is handled at Nginx + DNS level (see Vultr migration docs for smartminds.education; ESPThinking.co.uk will follow the same pattern but with a different primary domain).

---

## 8. How This Document Is Intended to Be Used

- **For Weaviate-assist**: as a single, up-to-date source describing the portfolio site structure, so semantic search over project docs returns accurate answers about routes, CTAs, and integrations.
- **For deployment planning**: as the frontend reference when mapping routes and subpaths in Nginx/Traefik on Vultr.
- **For collaborators**: as an orientation guide before editing landing pages or wiring new demos/tools into the portfolio shell.
