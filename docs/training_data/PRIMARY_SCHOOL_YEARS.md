# Primary School Years — Subject Guide (Draft)

Status: DRAFT (2026-02-05)

## Purpose
This document collects definitions, linked skills, suggested learning progressions, and data pointers for the "Primary School Years" subject group (early years through Year 6). Its goal is to provide:

- A clear subject label and canonical skill list for primary-aged learners
- Pointers to CHISG skills to be associated with year levels
- Links to frontend components and backend endpoints that consume or manage these skills

## Suggested Domains to Include
- FOUNDATION / FOUNDATIONAL skills (listening, follow instructions)
- LITERACY (early reading, phonics)
- NUMBERS / MATH (counting, number recognition)
- LANGUAGE (speaking, vocabulary)
- SOCIAL-EMOTIONAL (self, people, empathy)
- PHYSICAL / SENSORIMOTOR (gross and fine motor skills)

## How this connects to the codebase
- Skills and links are stored in the CHISG graph (Weaviate). See: `skills-map-platform/api/skills/skills.go`
- Expose and manage subjects via the Skills API endpoint: `/api/courses/subjects` (skills-map-platform/frontend `NavigationDrawer` reads this list)
- Frontend pages that reference subjects:
  - `skills-map-platform/frontend/src/components/SubjectsAndTopics.tsx`
  - `skills-map-platform/frontend/src/components/SkillsTree.tsx`
  - `skills-map-platform/frontend/src/components/GraphExplorer.tsx`

## Adding or Updating Skills for Primary Years
1. Create or update skills in CHISG (Weaviate) with `domain` fields that match our subject naming (e.g., `LITERACY`, `NUMBERS`).
2. Use `tools/seed-demo-data` or dedicated loader scripts to batch ingest skills for year levels.
3. Update `skills-map-platform` frontend where subjects are listed if a new canonical subject name is introduced.

## Finding Information (search tips)
- Find primary-focused skill names: `grep -R "phonics\|count\|reading\|number" -n` at repo root
- Find CHISG domain usage: `grep -R "domain: '" -n skills-map-platform | sed -n '1,120p'`

## Next Steps
- Curate an initial canonical list of ~80–120 primary skills mapped to Years R–6.
- Run the ingestion script in `tools/seed-demo-data` and verify via `/api/chisg/skillsandlinks`.
- Add a short UI panel (optional) in the Subjects drawer to show "Primary School Years" grouped by year level.

---

Created by: GitHub Copilot on request of project maintainer.
