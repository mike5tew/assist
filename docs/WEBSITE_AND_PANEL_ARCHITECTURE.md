# Website & Admin Panel Architecture

## Overview

This document outlines the website and admin panel infrastructure for the Smart Minds ecosystem. The solution consists of three layers:

1. **Public Website** (marketing/info)
2. **Admin Control Panel** (teacher/school management)
3. **Parent Portal** (SkillsProfileParent parent interface)

---

## Layer 1: Public Website (smartminds.education)

### Purpose
- Explain the system to potential schools and families
- Showcase CHISG framework, ETP spectra, skill hierarchies
- Provide onboarding flows for new schools
- Host documentation, case studies, research findings

### Technology Stack
- **Framework**: Next.js 14 (React + SSR) or Remix (if you prefer Vite)
- **Hosting**: Vultr (static + serverless functions)
- **CMS**: Optional (Contentful, Sanity) or Git-based (markdown in repo)
- **Styling**: Tailwind CSS (consistent with existing projects)
- **Analytics**: Posthog or Plausible (privacy-first)

### Key Pages

#### Homepage
```
Header:
  - Logo + Navigation (Features, Pricing, Docs, Sign In)
  
Hero Section:
  - Headline: "Understand Your Child's Unique Learning Profile"
  - Subheadline: "CHISG + ETP for parents, teachers, and researchers"
  - CTA: "Start Free Trial" → SkillsMarkbookMobile onboarding

Social Proof:
  - 3-4 quotes from teachers/schools using the system
  - Stats: "X schools, Y students, Z skills mapped"

Features Carousel:
  - Card 1: CHISG (579 skills, semantic relationships)
  - Card 2: ETP (8 spectra understanding, personalized learning)
  - Card 3: Home + School sync (parents see child's progress)
  - Card 4: Evidence-based (semantic links explain WHY)
```

#### Features Page
```
1. CHISG Skills Framework
   - Interactive skill tree preview (readonly)
   - "Explore a sample skill" → Show Joint Attention with semantic context
   - Comparison to other taxonomies (ICD-11, DSM-5, traditional curricula)

2. ETP Spectra Profiles
   - "Discover your learning style" → Interactive quiz (optional)
   - Show 8 spectra with descriptions
   - Case studies: "High energy_directionality child + how to support"

3. Parent App Features
   - Screenshot carousel of SkillsProfileParent
   - "See what parents say" testimonials

4. Teacher Dashboard
   - Screenshot of SkillsMarkbookMobile
   - Show graph visualization
   - "Real-time skill observation"

5. Researcher Access
   - "Download datasets"
   - "Query CHISG via API"
   - "Contribute semantic links"
```

#### Pricing Page
```
Tiers:
  - Individual Parent: Free or $5/month (read-only access)
  - Teacher: Free (up to 5 students) or $15/month (unlimited)
  - School: $500/month (site license, teacher + admin tools)
  - Researcher: Free (with attribution requirement)
  - Enterprise: Custom pricing (white-label, API access, data export)
```

#### Docs/Blog
```
- Getting started guides (teacher, parent, school admin)
- Case studies (XLA immunology, specific classroom outcomes)
- API documentation (for integrations)
- Research papers (user studies, CHISG validation)
- FAQ
```

#### Researcher Portal
```
- Interactive CHISG explorer (full graph)
- API documentation (GET /api/chisg/skills, /api/chisg/links)
- Download options (CSV, JSON, RDF)
- Attribution tools (show CHISG citations in your work)
- Contribution workflow (propose new skills, relationships)
```

### Architecture

```
smartminds.education (Next.js on Vultr)
├── pages/
│   ├── index.tsx (homepage)
│   ├── features.tsx
│   ├── pricing.tsx
│   ├── docs/[...slug].tsx (markdown-based)
│   ├── researcher/
│   │   ├── explorer.tsx (CHISG interactive viz)
│   │   ├── api.tsx (API docs)
│   │   └── download.tsx
│   └── auth/
│       ├── signup.tsx
│       └── login.tsx (redirects to admin panel)
├── components/
│   ├── SkillTree.tsx (react-flow based)
│   ├── ETPSpectrum.tsx (slider visualization)
│   ├── SkillCard.tsx (reusable skill preview)
│   └── ...
├── api/
│   ├── /auth/signup → skillstree-mysql (create user account)
│   ├── /chisg/skills → weaviate-ETPs-HumanOS-skillsmapinCHISG (read-only)
│   ├── /chisg/links → Same Weaviate
│   └── /newsletter/subscribe → Mailchimp/SendGrid
├── lib/
│   ├── weaviate.ts (GraphQL queries)
│   ├── mysql.ts (user accounts)
│   └── markdown-loader.ts
└── public/
    ├── case-studies/
    ├── screenshots/
    └── api-swagger.json
```

---

## Layer 2: Admin Control Panel (app.smartminds.education)

### Purpose
- School administrators manage teachers, students, courses
- Teachers record observations, view progress, export reports
- Admins manage billing, integrate with school systems, view analytics

### Technology Stack
- **Frontend**: React SPA (TypeScript, Vite)
- **Auth**: JWT + OAuth (Google/Microsoft for single sign-on)
- **Backend**: Go API (esp-organizer, extended)
- **Database**: skillstree-mysql (existing)
- **Hosting**: Vultr (containerized)

### Key Modules

#### School Administration
```
Dashboard:
  - Overview stats (teachers, students, active observations)
  - Recent activity log
  - Health indicators (data quality, system uptime)

Teachers Management:
  - List, add, remove teachers
  - Assign students/courses to teachers
  - View teacher activity (observations logged, time spent)
  - Set permissions (edit access, export rights, etc.)

Students Management:
  - Bulk import (CSV from school system)
  - Assign to teachers and courses
  - View enrollment status
  - Link parent accounts (if parent has opted in)

Courses Management:
  - Define custom courses (or use CHISG standard courses)
  - Set skill sequences (prerequisites, milestones)
  - View completion rates

Billing:
  - Current plan, usage, upcoming renewals
  - Invoice history
  - Upgrade/downgrade
  - Per-seat pricing or site license
```

#### Analytics & Reporting
```
School-Level Reports:
  - Skill adoption rates (which skills are most observed)
  - ETP spectrum distribution (how diverse are learners)
  - Progress timelines (skill mastery curves per cohort)
  - Teacher workload (observations/week per teacher)

Export Options:
  - CSV (skill progress for individual students)
  - PDF (individual progress reports for parent communication)
  - API access (pull data for integration with school SIS)
```

#### Integrations
```
- Google Classroom sync (pull student roster)
- Microsoft Teams (SkillsMarkbookMobile notifications)
- School SIS (Powerschool, Infinite Campus, etc.)
- LMS (Canvas, Blackboard student data)
```

### Architecture

```
app.smartminds.education (React SPA on Vultr)
├── src/
│   ├── pages/
│   │   ├── dashboard.tsx
│   │   ├── teachers/
│   │   │   ├── list.tsx
│   │   │   ├── [id].tsx (edit)
│   │   │   └── new.tsx
│   │   ├── students/
│   │   ├── courses/
│   │   ├── analytics/
│   │   ├── billing/
│   │   └── auth/login.tsx
│   ├── components/
│   │   ├── TeacherTable.tsx
│   │   ├── StudentProgress.tsx
│   │   ├── SkillAdoptionChart.tsx
│   │   └── ...
│   ├── api/
│   │   └── (no backend here, all calls to Go API)
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   └── useSchool.ts
│   └── lib/
│       └── api-client.ts
└── .env.production
```

**Backend API** (existing assist-api, extended):
```
POST /api/admin/schools/{school_id}/teachers → Create teacher
POST /api/admin/schools/{school_id}/students → Bulk import students
GET /api/admin/schools/{school_id}/analytics/skills → Skill adoption stats
GET /api/admin/schools/{school_id}/reports/{report_type} → Generate export
```

---

## Layer 3: Parent Portal (smartminds.education/parent)

### Purpose
- Read-only access to child's progress
- Understand CHISG skill definitions via semantic links
- View ETP profile with guidance
- See developmental milestones and next steps
- Optional: two-way messaging with teachers

### Technology Stack
- Shared frontend codebase with SkillsProfileParent (React Native → web version)
- Auth: Email/password or SSO (Google, Apple)
- Backend: New /api/parent/* endpoints (Go)
- Database: skillstree-mysql (existing) + parent-database (new)

### Key Features

#### Child Profile Page
```
Header:
  - Child's name, school, teacher(s)
  - Age + ETP profile (8 spectra as visual sliders)
  - Overall progress indicator

Skill Progress:
  - Pie chart: Mastered vs. Emerging vs. Not Started
  - Timeline graph: Skill mastery rate over time
  - Comparison to age group (optional, privacy-aware)

Quick Stats:
  - Total skills being tracked
  - Skills mastered this term
  - Top growth areas
  - Areas to focus on
```

#### Skill Explorer
```
When parent taps a skill (e.g., "Joint Attention"):

Part 1: Simple Explanation
  "Your child is learning to look where you point. This is the foundation
   of learning from others and understanding what others care about."

Part 2: Semantic Context (generated from CHISG links)
  Enables: Understanding gestures, Following instructions, Learning names
  Requires: Eye contact, Interest in social interaction
  Similar to: Shared attention, Turn-taking
  (Each link is clickable to explore related skills)

Part 3: What You Can Do
  - Point at things and see if child follows
  - Narrate what you're looking at
  - Play turn-taking games
  (Personalized based on child's ETP profile)

Part 4: Evidence
  "Research shows that children who develop joint attention earlier
   typically have better language outcomes. CHISG links show how this
   skill connects to social communication and learning capacity."
```

#### ETP Explorer
```
For each spectrum (e.g., social_gravity):
  - Description: "How drawn is your child to social interaction?"
  - Scale: Explanation of low vs. high
  - Your child's profile: Visual placement on spectrum
  - What it means: Plain language explanation
  - Skills alignment: Which skills are natural for this child, which need extra support
  - Sample activities: "Since your child has low social_gravity, they may prefer 1-1 or 
                       small group settings. Try these activities..."
```

#### Milestones & Growth
```
Age-appropriate skill sequencing:
  - "Skills your child should be developing now"
  - Visual progress: 3 skills mastered, 2 emerging, 5 not started
  - Prerequisite view: "These 3 skills will unlock these 5 new skills"
  - Comparison (opt-in): "Typical children this age are working on..."
```

#### Messages from School
```
Teacher can send structured updates:
  "Your child showed improvement in Joint Attention this week!
   We've noticed they're now looking at what others point to.
   This is a key step in learning together.
   
   Try these activities at home:
   - Point at interesting things
   - Play peek-a-boo
   - Take turns with toys"
```

### Architecture

```
smartminds.education/parent (React SPA on Vultr)
├── [Shared codebase with SkillsProfileParent]
├── src/
│   ├── pages/
│   │   ├── dashboard.tsx
│   │   ├── child/[child_id]/profile.tsx
│   │   ├── child/[child_id]/skill/[skill_id].tsx
│   │   ├── child/[child_id]/etp.tsx
│   │   ├── child/[child_id]/milestones.tsx
│   │   ├── messages.tsx
│   │   └── auth/login.tsx
│   ├── components/
│   │   ├── ChildProfile.tsx
│   │   ├── SkillExplainer.tsx
│   │   ├── ETPExplorer.tsx
│   │   └── ...
│   └── lib/
│       └── parent-api-client.ts
└── .env.production
```

**Backend API** (new Go endpoints):
```
POST /api/parent/auth/login → Email + password
GET /api/parent/children/{child_id} → Child profile
GET /api/parent/children/{child_id}/skill/{skill_id} → Skill details (queries CHISG)
GET /api/parent/children/{child_id}/etp → ETP profile + guidance
GET /api/parent/children/{child_id}/milestones → Age-based skill sequencing
GET /api/parent/children/{child_id}/messages → Messages from teacher
POST /api/parent/children/{child_id}/messages → Parent reply (optional)
```

---

## Unified Frontend Build Process

### Single Repository, Multiple Deployments

```
assist/ (monorepo)
├── websites/
│   ├── smartminds-marketing/     (Next.js public site)
│   ├── smartminds-admin/         (Vite React admin panel)
│   └── smartminds-parent/        (Vite React parent portal)
├── esp-organizer/                (Go backend, extended)
├── docker-compose.prod.yml
└── scripts/
    └── deploy-all.sh             (builds & pushes all 3 sites + backend)
```

### Build & Deployment Strategy

**On Local Mac**:
```bash
# 1. Build all three websites (optimize for production)
yarn build:marketing
yarn build:admin
yarn build:parent

# 2. Build backend
docker buildx build --platform linux/amd64 -t mike5tew/assist-api:latest ...

# 3. Push to Docker Hub + static hosting (Vultr Spaces or Netlify)
docker push mike5tew/assist-api:latest
netlify deploy --prod --dir=websites/smartminds-marketing/dist

# Or: Upload to Vultr Spaces, point Nginx at the bucket
```

**On Vultr Server**:
```bash
# docker-compose.prod.yml pulls all images
# Nginx routes based on domain:
#   smartminds.education → static marketing site (or Next.js container)
#   app.smartminds.education → admin panel (static SPA)
#   smartminds.education/parent → parent portal (static SPA)
#   api.smartminds.education → Go backend
```

---

## Domain & Routing Strategy

### DNS Setup
```
smartminds.education          → A record → Vultr server IP
api.smartminds.education      → CNAME → smartminds.education
app.smartminds.education      → CNAME → smartminds.education
www.smartminds.education      → CNAME → smartminds.education
researcher.smartminds.education → CNAME → smartminds.education
```

### Nginx Routing (in main-proxy)
```nginx
# Public website
server {
  server_name smartminds.education www.smartminds.education;
  location / {
    proxy_pass http://smartminds-marketing:3000;  # or static bucket
  }
}

# Admin panel
server {
  server_name app.smartminds.education;
  location / {
    proxy_pass http://smartminds-admin:3001;
  }
}

# Parent portal
server {
  server_name smartminds.education;
  location /parent {
    proxy_pass http://smartminds-parent:3002;
  }
}

# API backend
server {
  server_name api.smartminds.education smartminds.education;
  location /api/ {
    proxy_pass http://assist-api:8080;
  }
}

# Researcher portal
server {
  server_name researcher.smartminds.education smartminds.education;
  location /researcher {
    proxy_pass http://smartminds-marketing:3000;  # shared with marketing site
  }
}
```

---

## Summary: What Gets Built

| Layer | Purpose | Tech | Where |
|-------|---------|------|-------|
| **Public Website** | Marketing, docs, researcher explorer | Next.js + Tailwind | `smartminds.education` |
| **Admin Panel** | School/teacher management, analytics | React + Vite | `app.smartminds.education` |
| **Parent Portal** | Child progress + skill exploration | React + Vite | `smartminds.education/parent` |
| **Backend API** | Data layer for all three | Go (esp-organizer) | `api.smartminds.education` |
| **Databases** | Skills, users, progress, parents | Weaviate + MySQL + MongoDB | Vultr (containerized) |

---

## Next Steps

1. **Create `/websites` directory** with three Next.js/Vite projects
2. **Design the public website** (homepage, pricing, researcher explorer)
3. **Build admin panel** MVP (teacher/student management)
4. **Extend Go backend** with `/api/parent/*` endpoints
5. **Deploy to Vultr** (see VULTR_MIGRATION_GUIDE.md)
