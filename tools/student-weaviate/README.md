# Student Weaviate Schema & Migration Tools

Tools for setting up and populating the student-facing Weaviate database with humanOS behavioral elements for intelligent action plan generation.

## Overview

This toolkit creates and populates 9 Weaviate collections that enable the assist platform to generate personalized learning interventions based on:

- **ETP Profiles**: 8-spectrum emotional trigger point profiles + 2 global moderators for each student
- **CHISG Skills**: 459 skills from the Contextualised Hierarchical Iterative Semantic Groupings
- **Barriers & Levers**: Student learning barriers and intervention levers
- **Teaching Patterns**: Proven classroom strategies adapted for AI tutoring

### Architecture

The system has two parallel skill layers:

1. **ETP (Emotional Trigger Points)** - 8 core biological spectra that define WHERE a student naturally sits + 44 moderator skills that enable expanding comfort zone on each spectrum
2. **CHISG** - 459 curriculum skills that students actually learn

The ETP profile personalizes HOW a student learns the CHISG skills. The `etp_modulation` on each CHISG skill indicates which spectrum gets stressed when learning that skill.

## Collections

| Collection | Purpose | Data Source |
|------------|---------|-------------|
| `ETPProfile` | Student emotional trigger profiles (8 spectra + 2 moderators) | User input |
| `LearningIntervention` | Planned interventions with actions | Generated |
| `ActionLog` | History of actions and outcomes | Runtime |
| `StudentProgress` | Aggregated progress metrics | Computed |
| `InterventionLibrary` | Template interventions | Manual curation |
| `CHISGSkill` | 459 skills with ETP modulation | humanOS |
| `StudentBarrier` | Learning barriers with detection patterns | humanOS |
| `InterventionLever` | Actions that shift student state | humanOS |
| `TeachingPattern` | Classroom strategies (voltage reduction, etc.) | humanOS |

## Quick Start

```bash
# Set environment variables (optional - defaults shown)
export STUDENT_WEAVIATE_URL=http://localhost:8088
export HUMANOS_PATH=../../humanOS

# Run all migrations
make migrate-all
```

## Individual Commands

```bash
# Just create schema (no data)
make setup-schema

# Import CHISG skills (~459 skills)
make migrate-chisg

# Import barriers and levers only
make migrate-barriers

# Import teaching patterns only
make migrate-patterns

# Verify connection and schema
make verify
```

## Prerequisites

1. **Weaviate running** at `STUDENT_WEAVIATE_URL` (default: `http://localhost:8088`)
2. **humanOS repository** accessible at `HUMANOS_PATH` (default: `../../humanOS`)
3. **Go 1.21+** installed

## Action Plan Generation Flow

These collections enable the following action plan generation pipeline:

```
Student ETP Profile
        ↓
Match relevant CHISGSkills (by ETP modulation)
        ↓
Predict likely StudentBarriers (from patterns)
        ↓
Select effective InterventionLevers
        ↓
Apply appropriate TeachingPatterns
        ↓
Generate personalized LearningIntervention
        ↓
Log actions to ActionLog
        ↓
Update StudentProgress
```

## Tool Structure

```
student-weaviate/
├── Makefile                    # Orchestration
├── README.md                   # This file
├── setup-student-schema/       # Creates 9 collections
│   └── main.go
├── migrate-chisg-skills/       # Imports ~459 skills
│   └── main.go
├── import-barriers/            # Imports barriers & levers
│   └── main.go
└── import-teaching-patterns/   # Imports teaching patterns
    └── main.go
```

## Schema Details

### ETPProfile (8 Core Spectra + 2 Global Moderators)

**8 Core Spectra** (biological, innate, orthogonal) - each ranges from -100 to +100:

1. **social_gravity**: Independent ↔ Cohesive (social interaction drains vs charges)
2. **energy_directionality**: Inward ↔ Outward (external stimulation overloads vs energizes)
3. **voltage_sensitivity**: Insulated ↔ Conductive (emotional current tolerance)
4. **threat_response**: Passive ↔ Aggressive (avoidance vs confrontation voltage)
5. **care_response**: Detached ↔ Nurturing (distance vs care voltage to vulnerability)
6. **risk_tolerance**: Averse ↔ Seeking (risk creates anxiety vs excitement)
7. **integrity_logic**: Relativistic ↔ Absolutist (moral flexibility vs rigidity)
8. **mirror_neuron_tuning**: Selective ↔ Absorbent (others' emotions distinct vs shared)

**2 Global Moderators** (affect ALL spectra) - each 0-100:

- **pilot_strength**: Executive function capacity - hand on all sliders
- **current_load**: Stress/depletion level - narrows range on all spectra

**44 Moderator Skills**: Trainable skills (4-6 per spectrum) that expand comfort zone range

### CHISGSkill (ETP Modulation)

Each skill includes:
- `etpSpectrumId`: Which of the 8 spectra gets stressed when learning this skill
- `etpVector`: Direction (Toward Positive/Negative)
- `etpMechanism`: How it modulates the spectrum
- `etpEnergyCost`: Energy cost (Low/Medium/High)
- `enablesSkills` / `requiresSkills`: Skill graph relationships

### StudentBarrier

Barrier types:
- `initiation_barrier`: Won't start
- `chronic_barrier`: Persistent avoidance
- `structural_barrier`: Environmental/systemic
- `enrichment_need`: Needs more challenge

### TeachingPattern

Pattern types:
- `engagement`: Get student engaged
- `scaffolding`: Build capability gradually
- `motivation`: Provide motivation
- `avoidance_prevention`: Block avoidance tactics
- `proximity_support`: Intensive presence
- `relationship_building`: Build trust
- `metacognition`: Pattern awareness

## Related Documentation

- [STUDENT_WEAVIATE_SCHEMA.md](../../docs/STUDENT_WEAVIATE_SCHEMA.md) - Full schema specification
- [humanOS CHISG Elements](../../humanOS/data/chisg_elements.json) - Source skills data
- [humanOS Barriers](../../humanOS/shared/schemas/barriers.json) - Source barrier definitions
- [Classroom Interaction Patterns](../../humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md) - Teaching pattern source
