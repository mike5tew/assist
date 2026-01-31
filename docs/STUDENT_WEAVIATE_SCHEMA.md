# Student Weaviate DB Schema

This document defines the collections and structure for the student-facing Weaviate database deployed on Vultr. This DB stores actual student data, learning interventions, and educational pathways.

## Architecture Overview

The system has two parallel skill systems:

1. **ETP (Emotional Trigger Points)** - 8 core biological spectra that define WHERE a student naturally sits + 44 moderator skills that enable expanding comfort zone on each spectrum
2. **CHISG (Contextualised Hierarchical Iterative Semantic Groupings)** - 459 curriculum skills that students actually learn

The ETP profile personalizes HOW a student learns the CHISG skills. The `etp_modulation` on each CHISG skill indicates which spectrum gets stressed when learning that skill.

## Collections

### 1. ETPProfile
Stores individual student's Emotional Trigger Points profile - 8 core spectra + 2 global moderators.

```json
{
  "class": "ETPProfile",
  "description": "Student's ETP profile across 8 core spectra + 2 global moderators",
  "properties": [
    {
      "name": "studentId",
      "dataType": ["string"],
      "description": "Unique identifier for the student"
    },
    {
      "name": "profileTimestamp",
      "dataType": ["date"],
      "description": "When this profile was created/assessed"
    },
    {
      "name": "social_gravity",
      "dataType": ["int"],
      "description": "Independent (-100) to Cohesive (+100) - social interaction drains vs charges voltage"
    },
    {
      "name": "energy_directionality",
      "dataType": ["int"],
      "description": "Inward (-100) to Outward (+100) - external stimulation overloads vs energizes"
    },
    {
      "name": "voltage_sensitivity",
      "dataType": ["int"],
      "description": "Insulated (-100) to Conductive (+100) - emotional current tolerance"
    },
    {
      "name": "threat_response",
      "dataType": ["int"],
      "description": "Passive (-100) to Aggressive (+100) - avoidance vs confrontation voltage"
    },
    {
      "name": "care_response",
      "dataType": ["int"],
      "description": "Detached (-100) to Nurturing (+100) - distance vs care voltage to vulnerability"
    },
    {
      "name": "risk_tolerance",
      "dataType": ["int"],
      "description": "Averse (-100) to Seeking (+100) - risk creates anxiety vs excitement"
    },
    {
      "name": "integrity_logic",
      "dataType": ["int"],
      "description": "Relativistic (-100) to Absolutist (+100) - moral flexibility vs rigidity"
    },
    {
      "name": "mirror_neuron_tuning",
      "dataType": ["int"],
      "description": "Selective (-100) to Absorbent (+100) - others' emotions distinct vs shared"
    },
    {
      "name": "pilot_strength",
      "dataType": ["int"],
      "description": "Global moderator (0-100): Executive function capacity - hand on all sliders"
    },
    {
      "name": "current_load",
      "dataType": ["int"],
      "description": "Global moderator (0-100): Stress/depletion level - narrows range on all spectra"
    },
    {
    {
      "name": "spectrumVector",
      "dataType": ["number[]"],
      "description": "Vector representation of all 8 spectra for semantic search"
    },
    {
      "name": "profileSummary",
      "dataType": ["text"],
      "description": "Human-readable summary of the student's ETP profile"
    }
  ]
}
```

### 2. LearningIntervention
Personalized learning recommendations generated from a student's ETP profile.

```json
{
  "class": "LearningIntervention",
  "description": "Personalized learning pathways and interventions based on ETP profile",
  "properties": [
    {
      "name": "studentId",
      "dataType": ["string"],
      "description": "Reference to student"
    },
    {
      "name": "etpProfileId",
      "dataType": ["string"],
      "description": "Reference to the ETPProfile this intervention was based on"
    },
    {
      "name": "generatedTimestamp",
      "dataType": ["date"],
      "description": "When the intervention was generated"
    },
    {
      "name": "interventionType",
      "dataType": ["string"],
      "description": "e.g., 'social_approach', 'emotional_development', 'authority_navigation'"
    },
    {
      "name": "targetSpectrum",
      "dataType": ["string"],
      "description": "Which ETP spectrum this intervention targets (e.g., 'social_gravity')"
    },
    {
      "name": "currentValue",
      "dataType": ["int"],
      "description": "Student's current value on the targeted spectrum"
    },
    {
      "name": "recommendedApproach",
      "dataType": ["text"],
      "description": "The recommended learning/behavioral approach"
    },
    {
      "name": "rationale",
      "dataType": ["text"],
      "description": "Why this intervention is recommended given the student's profile"
    },
    {
      "name": "practicalExamples",
      "dataType": ["text[]"],
      "description": "Concrete examples of how to apply this approach"
    },
    {
      "name": "resources",
      "dataType": ["string[]"],
      "description": "Links/references to learning materials, articles, etc."
    },
    {
      "name": "difficulty",
      "dataType": ["string"],
      "description": "e.g., 'beginner', 'intermediate', 'advanced'"
    },
    {
      "name": "estimatedTimeToMastery",
      "dataType": ["string"],
      "description": "e.g., '2-3 weeks', '1-2 months'"
    },
    {
      "name": "status",
      "dataType": ["string"],
      "description": "e.g., 'recommended', 'in_progress', 'completed'"
    },
    {
      "name": "vector",
      "dataType": ["number[]"],
      "description": "Embedding of the intervention content for semantic search"
    }
  ]
}
```

### 3. ActionLog
Records of actions/assignments given to students based on their interventions.

```json
{
  "class": "ActionLog",
  "description": "Action items and assignments tracking student progress",
  "properties": [
    {
      "name": "studentId",
      "dataType": ["string"],
      "description": "Reference to student"
    },
    {
      "name": "interventionId",
      "dataType": ["string"],
      "description": "Reference to the LearningIntervention this action is tied to"
    },
    {
      "name": "actionTitle",
      "dataType": ["string"],
      "description": "Brief title of the action/assignment"
    },
    {
      "name": "actionDescription",
      "dataType": ["text"],
      "description": "Detailed description of what the student should do"
    },
    {
      "name": "createdDate",
      "dataType": ["date"],
      "description": "When the action was assigned"
    },
    {
      "name": "dueDate",
      "dataType": ["date"],
      "description": "Target completion date"
    },
    {
      "name": "completedDate",
      "dataType": ["date"],
      "description": "When student marked it complete (if applicable)"
    },
    {
      "name": "status",
      "dataType": ["string"],
      "description": "e.g., 'assigned', 'in_progress', 'completed', 'overdue'"
    },
    {
      "name": "priority",
      "dataType": ["string"],
      "description": "e.g., 'low', 'medium', 'high'"
    },
    {
      "name": "reflectionNotes",
      "dataType": ["text"],
      "description": "Student's reflection on their learning/progress"
    },
    {
      "name": "impactOnSpectrum",
      "dataType": ["string"],
      "description": "Which ETP spectrum(s) this action is expected to influence"
    }
  ]
}
```

### 4. StudentProgress
Aggregated view of student's journey and profile evolution over time.

```json
{
  "class": "StudentProgress",
  "description": "Student's overall progress, milestones, and profile evolution",
  "properties": [
    {
      "name": "studentId",
      "dataType": ["string"],
      "description": "Unique student identifier"
    },
    {
      "name": "enrollmentDate",
      "dataType": ["date"],
      "description": "When student joined the learning system"
    },
    {
      "name": "latestProfileId",
      "dataType": ["string"],
      "description": "Most recent ETPProfile assessment"
    },
    {
      "name": "profileEvaluationCount",
      "dataType": ["int"],
      "description": "Number of times student has been assessed"
    },
    {
      "name": "interventionsCreated",
      "dataType": ["int"],
      "description": "Total interventions generated for this student"
    },
    {
      "name": "actionsCompleted",
      "dataType": ["int"],
      "description": "Number of actions successfully completed"
    },
    {
      "name": "averageCompletionTime",
      "dataType": ["string"],
      "description": "e.g., '5.2 days' average time to complete actions"
    },
    {
      "name": "spectrumGrowthAreas",
      "dataType": ["string[]"],
      "description": "Spectra where student has shown most improvement"
    },
    {
      "name": "focusAreas",
      "dataType": ["string[]"],
      "description": "Current priority areas for development"
    },
    {
      "name": "learningPathway",
      "dataType": ["text"],
      "description": "Summary of the student's learning journey"
    },
    {
      "name": "lastActivityDate",
      "dataType": ["date"],
      "description": "When student last interacted with the system"
    }
  ]
}
```

### 5. InterventionLibrary
Pre-built, reusable intervention templates that can be customized for students.

```json
{
  "class": "InterventionLibrary",
  "description": "Repository of intervention templates for different ETP spectrum combinations",
  "properties": [
    {
      "name": "templateId",
      "dataType": ["string"],
      "description": "Unique identifier for this template"
    },
    {
      "name": "name",
      "dataType": ["string"],
      "description": "Name of the intervention template"
    },
    {
      "name": "targetSpectra",
      "dataType": ["string[]"],
      "description": "Which ETP spectra this template addresses"
    },
    {
      "name": "description",
      "dataType": ["text"],
      "description": "What this intervention is designed to help with"
    },
    {
      "name": "approaches",
      "dataType": ["text[]"],
      "description": "List of recommended approaches/strategies"
    },
    {
      "name": "examples",
      "dataType": ["text[]"],
      "description": "Practical examples"
    },
    {
      "name": "resources",
      "dataType": ["string[]"],
      "description": "URLs and references to supporting materials"
    },
    {
      "name": "category",
      "dataType": ["string"],
      "description": "e.g., 'behavioral', 'cognitive', 'social', 'emotional'"
    },
    {
      "name": "vector",
      "dataType": ["number[]"],
      "description": "Semantic embedding for retrieval"
    }
  ]
}
```

## Relationships

```
ETPProfile
  ├── 1:N → LearningIntervention (one profile can have many interventions)
  └── 1:N → StudentProgress (historical profile records)

LearningIntervention
  ├── N:1 → ETPProfile
  ├── 1:N → ActionLog (one intervention can have many action items)
  └── M:N → InterventionLibrary (based on library templates)

ActionLog
  ├── N:1 → LearningIntervention
  └── N:1 → StudentProgress

StudentProgress
  └── N:1 → Student (external reference, could be in MongoDB)
```

## Deployment Considerations

- **Vectorization**: Use semantic embeddings (AWS Bedrock + text2vec-aws) for:
  - `spectrumVector` in ETPProfile
  - Intervention descriptions in LearningIntervention and InterventionLibrary
  - Enables semantic similarity search for finding relevant interventions

- **Indexing**: Create indexes on:
  - `studentId` (across all collections)
  - `status` (for queries on active/completed items)
  - `targetSpectrum` (for targeted recommendations)

- **Scaling**: As student count grows, consider:
  - Sharding by `studentId`
  - Archiving old ActionLogs
  - Caching popular InterventionLibrary templates

## Next Steps

1. Review and refine this schema
2. Create Weaviate schema initialization script
3. Build action plan generation algorithm
4. Implement student profile assessment flow

---

## Integration: humanOS Behavioral Elements

The humanOS project contains extensive behavioral modeling that can be integrated into this schema to generate deeper, more contextual action plans.

### Additional Collections (from humanOS)

### 6. CHISGSkill
Skills graph with ETP modulation data (from `chisg_elements.json` - 23,680 lines).

```json
{
  "class": "CHISGSkill",
  "description": "Comprehensive Human Intelligence Skills Graph - skills with ETP relationships",
  "properties": [
    {
      "name": "chisgId",
      "dataType": ["string"],
      "description": "Unique CHISG identifier (e.g., 'CHISG_1')"
    },
    {
      "name": "name",
      "dataType": ["string"],
      "description": "Skill name (e.g., 'Strategic Foresight')"
    },
    {
      "name": "description",
      "dataType": ["text"],
      "description": "What this skill enables"
    },
    {
      "name": "domain",
      "dataType": ["string"],
      "description": "e.g., 'FOCUS & TOOLS', 'EMOTIONAL', 'SOCIAL'"
    },
    {
      "name": "layer",
      "dataType": ["string"],
      "description": "e.g., 'Intervention', 'Foundation', 'Applied'"
    },
    {
      "name": "etpSpectrumId",
      "dataType": ["string"],
      "description": "Which ETP spectrum this skill modulates"
    },
    {
      "name": "etpVector",
      "dataType": ["string"],
      "description": "Direction: 'Toward Positive (Growth)' or 'Toward Negative'"
    },
    {
      "name": "etpMechanism",
      "dataType": ["string"],
      "description": "How it affects the spectrum (e.g., 'Creating internal drivers')"
    },
    {
      "name": "energyCost",
      "dataType": ["string"],
      "description": "Low/Medium/High energy required"
    },
    {
      "name": "enables",
      "dataType": ["string[]"],
      "description": "Skills this one enables"
    },
    {
      "name": "requires",
      "dataType": ["string[]"],
      "description": "Prerequisites for this skill"
    },
    {
      "name": "vector",
      "dataType": ["number[]"],
      "description": "Semantic embedding for skill matching"
    }
  ]
}
```

### 7. StudentBarrier
Barrier profiles with detection patterns and effective interventions.

```json
{
  "class": "StudentBarrier",
  "description": "Learning barriers with avoidance tactics and intervention strategies",
  "properties": [
    {
      "name": "barrierId",
      "dataType": ["string"],
      "description": "e.g., 'confrontational_showoff', 'silent_avoider'"
    },
    {
      "name": "name",
      "dataType": ["string"],
      "description": "Human-readable barrier name"
    },
    {
      "name": "category",
      "dataType": ["string"],
      "description": "acute | chronic | structural | enrichment"
    },
    {
      "name": "description",
      "dataType": ["text"],
      "description": "What this barrier looks like"
    },
    {
      "name": "activatedETPs",
      "dataType": ["string[]"],
      "description": "Which ETPs are triggered by this barrier"
    },
    {
      "name": "avoidanceTactics",
      "dataType": ["string[]"],
      "description": "How students with this barrier avoid engagement"
    },
    {
      "name": "effectiveLeverIds",
      "dataType": ["string[]"],
      "description": "References to InterventionLever documents"
    },
    {
      "name": "underlyingCause",
      "dataType": ["text"],
      "description": "Root cause of the barrier"
    },
    {
      "name": "vector",
      "dataType": ["number[]"],
      "description": "Embedding for semantic matching"
    }
  ]
}
```

### 8. InterventionLever
Specific intervention strategies from 12 years of classroom experience.

```json
{
  "class": "InterventionLever",
  "description": "Specific teaching/coaching intervention strategies",
  "properties": [
    {
      "name": "leverId",
      "dataType": ["string"],
      "description": "Unique lever identifier"
    },
    {
      "name": "name",
      "dataType": ["string"],
      "description": "e.g., 'Voltage Reduction Through Familiarity'"
    },
    {
      "name": "description",
      "dataType": ["text"],
      "description": "What this lever does"
    },
    {
      "name": "steps",
      "dataType": ["text[]"],
      "description": "Step-by-step implementation"
    },
    {
      "name": "prerequisites",
      "dataType": ["string[]"],
      "description": "What needs to be in place first"
    },
    {
      "name": "benefits",
      "dataType": ["string[]"],
      "description": "Expected outcomes"
    },
    {
      "name": "etpReduction",
      "dataType": ["string[]"],
      "description": "Which ETPs this lever reduces/modulates"
    },
    {
      "name": "brainStateTarget",
      "dataType": ["string"],
      "description": "primal | emotional | rational"
    },
    {
      "name": "whenToUse",
      "dataType": ["text[]"],
      "description": "Situational triggers for using this lever"
    },
    {
      "name": "timeInvestment",
      "dataType": ["string"],
      "description": "e.g., '5 minutes', '1 week'"
    },
    {
      "name": "expectedTimeline",
      "dataType": ["string"],
      "description": "When to expect results"
    },
    {
      "name": "vector",
      "dataType": ["number[]"],
      "description": "Embedding for semantic matching"
    }
  ]
}
```

### 9. TeachingPattern
Classroom interaction patterns adapted for individual coaching.

```json
{
  "class": "TeachingPattern",
  "description": "Proven classroom interaction patterns for coaching",
  "properties": [
    {
      "name": "patternId",
      "dataType": ["string"],
      "description": "Unique pattern identifier"
    },
    {
      "name": "name",
      "dataType": ["string"],
      "description": "e.g., 'Question Killer Game', 'Micro-Victory Celebration'"
    },
    {
      "name": "classroomPattern",
      "dataType": ["text"],
      "description": "How it works in a classroom setting"
    },
    {
      "name": "aiTutorAdaptation",
      "dataType": ["text"],
      "description": "How to adapt for 1-on-1 AI coaching"
    },
    {
      "name": "whyItWorks",
      "dataType": ["text"],
      "description": "Psychology behind the pattern"
    },
    {
      "name": "targetBarriers",
      "dataType": ["string[]"],
      "description": "Which barriers this pattern addresses"
    },
    {
      "name": "etpContext",
      "dataType": ["string[]"],
      "description": "Which ETP spectra are relevant"
    },
    {
      "name": "implementationExamples",
      "dataType": ["text[]"],
      "description": "Concrete examples"
    },
    {
      "name": "vector",
      "dataType": ["number[]"],
      "description": "Embedding for retrieval"
    }
  ]
}
```

---

## Enhanced Action Plan Generation Flow

With humanOS integration, the action plan generation becomes:

```
1. Student submits ETP Profile (17 spectra values)
                ↓
2. Query CHISGSkill collection for skills that:
   - Match the student's ETP vector direction
   - Have appropriate energy cost for their profile
   - Build on existing strengths
                ↓
3. Query StudentBarrier collection using:
   - ETP profile to predict likely barriers
   - Historical interaction patterns (if available)
                ↓
4. Query InterventionLever collection for:
   - Levers that address predicted barriers
   - Levers appropriate for current brain state
                ↓
5. Query TeachingPattern collection for:
   - Patterns that work with student's ETP profile
   - 1-on-1 coaching adaptations
                ↓
6. Generate LearningIntervention documents:
   - Personalized to ETP profile
   - Informed by barrier prediction
   - Using proven teaching strategies
   - With specific action items
                ↓
7. Store in ActionLog for tracking
```

## Data Migration Plan

1. **CHISG Skills**: Parse `humanOS/data/chisg_elements.json` → bulk import to Weaviate
2. **Barriers**: Extract from `humanOS/backend/internal/etp/barriers.json` → import
3. **Intervention Levers**: Extract from ETP types and barrier levers → import
4. **Teaching Patterns**: Parse `humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md` → import

## Benefits of Integration

- **Semantic Matching**: Find skills that modulate specific ETP spectra
- **Barrier Prediction**: Use ETP profile to predict likely learning obstacles
- **Evidence-Based Interventions**: 12 years of classroom-tested strategies
- **Personalized Learning Paths**: Connect skills → barriers → interventions → actions
