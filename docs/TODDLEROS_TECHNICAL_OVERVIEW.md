# ToddlerOS — Technical Overview

**Last Updated**: 2025-06-28  
**Status**: 12-spectrum model implemented. Global moderators removed. Transient states removed.  
**Document Purpose**: Technical architecture connecting ToddlerOS to the skills map, ETP framework, and pathfinder engine.

---

## 1. What ToddlerOS Is

ToddlerOS is a parent-facing product that replaces generational memory with structured observation prompts. It is delivered as a **physical activity book** (printed card stock + sticker sheet) with a **companion parent app** (observation log + progress rosette). The child's interface is paper. The parent's interface is a phone — used only for logging, never shown to the child.

**Tagline**: "The Coding Lesson You Never Knew You Needed"

ToddlerOS is not a standalone system. It is a **lens** into the shared ESP Thinking infrastructure — the same three-layer architecture that powers LAO, PrimaryOS, CareerOS, and humanOS. This document describes the technical connections.

---

## 2. Three-Layer Architecture

Every product in the ecosystem uses the same three layers. ToddlerOS filters each one for ages 0–5:

| Layer | System | What ToddlerOS Uses |
|-------|--------|---------------------|
| **1. Emotional/Biological** | ETP Framework (`spectra.go`, `response_matrix.go`) | 12 spectra → 12 weekly themes. Parent observations build an ETP profile as a byproduct of play. |
| **2. Skills/Capacities** | Universal Skills Graph (`skills_key` + `skill_link` DAG) | 30–50 foundational skill nodes (ages 0–5) — the base layer that academic skills assume but never name. |
| **3. Academic/Domain** | Not used directly | ToddlerOS operates *below* the academic layer. Its purpose is to build the capacities that make Layer 3 possible. |

### How the Layers Connect in ToddlerOS

```
┌─────────────────────────────────────────────────────────────┐
│                    ACTIVITY BOOK (Output)                     │
│  ┌──────────────────┐  ┌──────────────────────────────────┐  │
│  │   "With You"     │  │        "While You..."            │  │
│  │   3–5 min        │  │        5–15 min                  │  │
│  │   Parent active  │  │        Child independent         │  │
│  └────────┬─────────┘  └──────────────┬───────────────────┘  │
│           │                            │                      │
│           ▼                            ▼                      │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              RESPONSE MATRIX (Layer 1)                  │  │
│  │  PlayAvenues + RecommendedLanguage + AvoidLanguage      │  │
│  │  24 entries: 12 spectra × 2 settings                     │  │
│  │  Source: response_matrix.go (576 lines, compiled Go)    │  │
│  └────────┬────────────────────────────────────────────────┘  │
│           │                                                   │
│           ▼                                                   │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │            FOUNDATIONAL SKILLS (Layer 2)                │  │
│  │  ~30–50 nodes in skills_key, DevelopmentAge 0.0–5.0     │  │
│  │  7 high-leverage capacities + sub-skills                │  │
│  │  Connected via skill_link DAG to higher-age nodes       │  │
│  └────────┬────────────────────────────────────────────────┘  │
│           │                                                   │
│           ▼                                                   │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │           PATHFINDER ENGINE (Shared)                     │  │
│  │  Given: current scores + ETP profile                    │  │
│  │  Output: which theme/skill to practise next             │  │
│  │  ETP modulation: edge weights vary by spectrum profile  │  │
│  └─────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Connection to the ETP Framework

### 3.1 The 12 Spectra

The ETP (Emotional Tuning Profile) framework defines 12 biological spectra, each a morally-neutral axis with two poles. These are defined in [`spectra.go`](../esp-organizer/internal/domain/etp/spectra.go):

| ID | Spectrum | Negative Pole | Positive Pole | ToddlerOS Theme |
|----|----------|---------------|---------------|------------------|
| 1 | `social_gravity` | Independent | Cohesive | **"Near and Far"** (alone vs. together play) |
| 2 | `energy_directionality` | Inward | Outward | **"Loud and Quiet"** (internal vs. external processing) |
| 3 | `voltage_sensitivity` | Insulated | Conductive | **"Big Feelings, Small Feelings"** (emotional intensity) |
| 4 | `threat_response` | Passive | Aggressive | **"Yes and No"** (standing ground vs. stepping back) |
| 5 | `care_response` | Detached | Nurturing | **"Mine and Yours"** (sharing vs. keeping) |
| 6 | `risk_tolerance` | Averse | Seeking | **"Try and Wait"** (jumping in vs. holding back) |
| 7 | `integrity_logic` | Relativistic | Absolutist | **"Same and Different"** (rules are fixed vs. flexible) |
| 8 | `mirror_neuron_tuning` | Selective | Absorbent | **"Your Feelings, My Feelings"** (emotional boundaries) |
| 9 | `orderliness` | Flexible | Ordered | **"Tidy and Messy"** (structure vs. spontaneity) |
| 10 | `responsibility_threshold` | Deflecting | Absorbing | **"My Fault, Your Fault"** (blame vs. ownership) |
| 11 | `loss_sensitivity` | Detached | Territorial | **"Keeping and Letting Go"** (attachment to things) |
| 12 | `libido` | Restrained | Expressive | **"Wanting and Waiting"** (drive energy expression) |

Each spectrum has 4–6 trainable skills (61 total across all 12), defined in `SpectrumSkills` in the same file. These are the adult/older-child skills — ToddlerOS addresses the same spectra through play rather than explicit skill training.

### 3.2 The Response Matrix

The response matrix ([`response_matrix.go`](../esp-organizer/internal/domain/etp/response_matrix.go)) contains 24 entries — one per pole of each spectrum. Each entry defines:

```go
type ResponseMatrixEntry struct {
    SpectrumID          int      // Which of the 12 spectra
    SpectrumName        string   // e.g., "social_gravity"
    Setting             string   // e.g., "independent" or "cohesive"
    SettingEnd          string   // "negative" or "positive"
    BarrierType         string   // "verbalising", "starting", or "mistakes"
    BarrierDescription  string   // What the barrier looks like in practice
    AvoidLanguage       []string // What NOT to say (moral framework)
    RecommendedLanguage []string // What TO say (engineering framework)
    PlayAvenues         []string // Activity strategies for this setting
}
```

**ToddlerOS uses these fields directly**:

| Field | Book Usage |
|-------|-----------|
| `PlayAvenues` | Source content for "With You" and "While You" activities |
| `RecommendedLanguage` | Parent prompt text printed on each spread |
| `AvoidLanguage` | "Common traps" sidebar — what to avoid saying |
| `BarrierDescription` | "What you might notice" observation prompt |
| `BarrierType` | Classifies which kind of difficulty the child may show |

Example — Social Gravity, Independent setting:
```
BarrierDescription: "Won't ask for help when stuck"
AvoidLanguage:      ["You should have asked for help earlier"]
RecommendedLanguage: ["Your battery was running low. Next time, you could 
                       signal me and I'll come to you — no need to come find me."]
PlayAvenues:        ["Parallel play with gradual invitation",
                     "Let them work alone in the same space",
                     "Occasionally wonder aloud near them: 'I'm wondering if...'
                      — invitation without demand"]
```

### 3.3 Core Pedagogy: Microdosing the Poles

Toddlers' ETP defaults haven't locked in yet — ages 0–5 represent the neuroplasticity sweet spot. Instead of teaching *about* spectra, ToddlerOS lets children **experience both poles** in tiny, safe doses:

- Each spectrum = one theme week
- 12 spectra × 3 difficulty cycles = 36 themed weeks + seasonal specials = expandable programme
- Delivered as seasonal editions (dip-in/dip-out model)
- Non-sequential: any page is the right page

### 3.4 The Play-First Principle

The adult's job is codified as four steps in `PlayFirstPrinciple` (exported constant):

```go
var PlayFirstPrinciple = []string{
    "Observe the setting without judgment",
    "Name it neutrally (\"You're in Independent mode right now\")",
    "Design the play that gently expands range",
    "Debrief in engineering language (\"When you had to wait, what happened in your brain?\")",
}
```

For ToddlerOS these are simplified to age-appropriate versions: Observe → Name → Play → Notice.

### 3.5 Observation as Profile Building

The parent never fills in an ETP questionnaire. Instead, after each theme week, a single observation prompt maps to a spectrum data point:

| Theme | Observation Prompt | Maps To |
|-------|-------------------|--------|
| Near and Far | "Which was easier — playing alone or playing together?" | `social_gravity` value |
| Loud and Quiet | "Did they prefer the noisy game or the quiet one?" | `energy_directionality` value |
| Big Feelings, Small Feelings | "When something went wrong, was the reaction big or small?" | `voltage_sensitivity` value |
| Yes and No | "When told no, did they push back or withdraw?" | `threat_response` value |
| Mine and Yours | "When asked to share, was it easy or hard?" | `care_response` value |
| Try and Wait | "Did they jump straight in or watch first?" | `risk_tolerance` value |
| Same and Different | "When the rules changed, was that OK or upsetting?" | `integrity_logic` value |
| Your Feelings, My Feelings | "When another child cried, did they notice?" | `mirror_neuron_tuning` value |
| Tidy and Messy | "When it was time to tidy up, was that a fight?" | `orderliness` value |
| My Fault, Your Fault | "When something went wrong, did they blame others or take all the blame?" | `responsibility_threshold` value |
| Keeping and Letting Go | "When something was taken away, was it a crisis or no big deal?" | `loss_sensitivity` value |
| Wanting and Waiting | "When they wanted something, could they wait or did they need it now?" | `libido` value |

Over the full programme (3 cycles), each spectrum gets 3 observations. That's enough to build a meaningful ETP profile without the parent ever knowing they're creating one.

### 3.6 Accelerator Skills

Executive function capacities (formerly "global moderators") are now properly classified as **accelerator skills** in the skills DAG. They include:

- `accel_executive_function` — the ability to deliberately move sliders rather than having triggers move them
- `accel_self_regulation` — the ability to tolerate discomfort without dysregulation
- `accel_working_memory`, `accel_cognitive_flexibility`, `accel_inhibitory_control`, etc.

These are defined in `PRIMARY_SCHOOL_ACCELERATORS.json` (12 skills) and sit in the skills DAG as cross-cutting enablers. They are NOT personality spectra — they are trainable capacity muscles.

The previous `pilot_strength` and `current_load` moderators have been removed:
- **pilot_strength** → `accel_executive_function` + `accel_self_regulation` in the skills DAG
- **current_load** → removed entirely. Transient states (tired, hungry, ill) change hourly and are impractical for teachers/parents to maintain. Energy directionality already captures stress-reaction direction. The book says: "Bad day? Skip it."

### 3.7 Language Shifts

The unified language shifts (`UnifiedLanguageShifts` in `response_matrix.go`) apply across all spectra:

| Old (Moral) | New (Engineering) |
|-------------|-------------------|
| "You are..." | "Your setting is..." |
| "You should..." | "Your brain is telling you..." |
| "That was wrong" | "That didn't work — let's look at why" |
| "Try harder" | "What's blocking you right now?" |
| "Stop it" | "I notice you're in X setting — is that helping?" |

These are printed in each activity book as a reference card. For ToddlerOS, the language is simplified:

| Old | ToddlerOS Version |
|-----|-------------------|
| "You are..." | "Right now you're..." |
| "Your brain is telling you..." | "Your brain is saying..." |
| "That didn't work — let's look at why" | "That didn't work. What happened?" |

---

## 3.8 Activity Experience Model (Emotional Compiler)

**New**: [`activity_experience.go`](../esp-organizer/internal/domain/etp/activity_experience.go)

The response matrix gives per-spectrum guidance. But real activities are **multi-threaded events** — a single activity like "Dance Party" simultaneously activates `social_gravity` (cohesive), `voltage_sensitivity` (conductive), and `mirror_neuron_tuning` (absorbent), while also developing `gross_motor`, `rhythm`, and `midline_crossing` skills. The Activity Experience model captures this full causal chain.

### Core Insight: Activities Have Side Effects

An activity designed to build one capacity can accidentally create a **negative emotional imprint** (a "bug" in the child's OS). Example:

1. **Dance Party** targets `social_gravity` (cohesive) and `mirror_neuron_tuning` (absorbent)
2. A peer observes → triggers `threat_response` (passive pole)
3. Child freezes → associates dancing with shame (**`shame_gate` imprint**)
4. Left unmanaged: child withdraws from movement activities entirely

The schema models these side effects as first-class citizens with a **rescue route**:

| Side Effect Field | Purpose |
|-------------------|---------|
| `trigger` | What causes it (peer observation, spill, failure) |
| `triggered_spectrum` | Which spectrum spikes |
| `imprint` | What lasting association forms if unmanaged |
| `signals` | Observable behaviours the parent watches for |
| `rescue_language` | Engineering language to name and validate |
| `counterweight` | The spectrum shift that discharges the voltage |

### The Counterweight Strategy (Zen Warriors)

Every side effect has a three-step counterweight — this is the core ToddlerOS mechanism:

```
1. NAME IT    → "Your brain noticed people watching. That made your chest feel tight."
2. VALIDATE   → "That prickly feeling is real. Your body is protecting you."
3. INVITE     → "Let's find a spot where it's just us. Your voltage will come back down."
```

This is the "Zen Warrior" approach: expose the trigger, acknowledge the pain, practise the counterweight on that spectrum. Over time, toddlers build **spectrum range** — the ability to stay regulated across wider voltage swings.

### Spread Template (Book Page Layout)

Each activity compiles down to a `BookSpreadContent` struct — the "view" layer that hides the complexity behind calm, simple slots:

```
┌─────────────────────────────────────────────────────────────────┐
│ LEFT PAGE                                                       │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  💃 DANCE PARTY                            ⏱ 3 min      │   │
│  │                                           🎵 Just music  │   │
│  ├──────────────────────────────────────────────────────────┤   │
│  │                                                          │   │
│  │  Put on a song you love and dance together.              │   │
│  │  Copy each other's moves.                                │   │
│  │                                                          │   │
│  ├──────────────────────────────────────────────────────────┤   │
│  │  👀 WHAT TO NOTICE                                       │   │
│  │  Do they copy your moves or invent their own?            │   │
│  ├──────────────────────────────────────────────────────────┤   │
│  │  🛟 IF THEY STOP OR HIDE...                              │   │
│  │  That's OK. Sit down and sway gently.                    │   │
│  │  The dance can be just fingers.                          │   │
│  ├──────────────────────────────────────────────────────────┤   │
│  │  ✅ TRY: "Your body doesn't want to move. That's fine." │   │
│  │  ❌ NOT: "Come on, it's fun! Nobody's watching!"         │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│ RIGHT PAGE                                                      │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  [COLOURING / ACTIVITY]                                  │   │
│  │  "Draw your favourite dance move"                        │   │
│  │                                                          │   │
│  │  What does your dance look like on paper?                │   │
│  │                                                          │   │
│  │              [Illustration space]                         │   │
│  │                                                          │   │
│  │                                    [Sticker spot] 🌟     │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### Sample Activities in Schema

Three activities are seeded in `activity_experience.go` as proof-of-concept:

| Activity | Mode | Primary Skills | Spectra Activated | Side Effects |
|----------|------|---------------|-------------------|--------------|
| Dance Party | with_you | gross_motor, rhythm, joint_attention | social_gravity (primary), voltage_sensitivity, mirror_neuron_tuning, threat_response (latent) | shame_gate (peer observation), sensory_flooding (overstimulation) |
| The Pouring Game | with_you | fine_motor, sustained_attention, inhibitory_control | orderliness (primary), risk_tolerance, threat_response (latent) | failure_avoidance (spill→freeze), frustration_explosion (spill→aggression) |
| Find My Hands | with_you | joint_attention, turn_taking, inhibitory_control | social_gravity (primary), energy_directionality, care_response (latent) | failure_avoidance (wrong guess), turn_taking_refusal (won't swap roles) |

---

## 3.9 Ventral Emotional Drivers — "The Ghost in the Machine"

**New**: [`ventral_drivers.go`](../esp-organizer/internal/domain/etp/ventral_drivers.go)

The spectra tell us **what** is happening (slider positions). Ventral Drivers tell us **why** — the deep-brain reward loop that pulls all the sliders into a characteristic configuration. A Ventral Driver is not a single spectrum reading. It is a **pattern** — a combination of 2–3 dominant spectra that form a persistent motivational force.

### The Key Difference

| Concept | Measures | Changes With | Example |
|---------|----------|-------------|---------|
| Spectrum reading | Current slider position | Context (tired, excited, safe, threatened) | `social_gravity = +1.2 (cohesive)` |
| Ventral Driver | Why the slider keeps going there | Very slowly — deep-brain reward circuitry | "I seek connection because being seen = being real" |

### The 7 Archetypal Drivers

Each driver maps to a recognisable **persona** — a warm, non-clinical label parents can identify without any technical knowledge:

| Persona | Icon | Spectrum Signature | Reward Condition | Threat Condition | Child's Voice |
|---------|------|-------------------|------------------|------------------|---------------|
| **The Impressor** | ⭐ | social_gravity↑ + mirror_neuron↑ + energy↑ | Positive reflection in the adult's eyes | Being unobserved or receiving neutral feedback | "I exist because you see me" |
| **The Explorer** | 🔍 | risk_tolerance↑ + energy↓ + social_gravity↓ | Sensory feedback from environment | Being pulled away from investigation | "I understand by touching" |
| **The Sentinel** | 🛡️ | threat_response↓ + orderliness↑ + risk_tolerance↓ | Absence of change. Predictability confirmed | Unexpected change, broken routine | "I'm safe when nothing changes" |
| **The Connector** | 🤝 | care_response↑ + mirror_neuron↑ + social_gravity↑ | Making someone feel better | Witnessing distress they can't fix | "Are you OK? Can I help?" |
| **The Rebel** | ⚡ | threat_response↑ + integrity_logic↓ + social_gravity↓ | Agency confirmed. "I chose this" | Loss of control. Being overridden | "I decide. Not you." |
| **The Performer** | 🎭 | voltage_sensitivity↑ + energy↑ + risk_tolerance↑ | Full emotional expression received and matched | Being told to calm down | "If you can't feel it, it doesn't count" |
| **The Observer** | 👁️ | energy↓ + mirror_neuron↓ + voltage↓ | Understanding achieved. Pattern is clear | Being forced to act before processing completes | "I need to understand before I move" |

### The Heat Shield Mechanism

ToddlerOS does NOT "fix" ventral drivers — they are **power sources**. The goal is to build a **Heat Shield** (metacognitive distance) between the driver and the executive function:

```
1. NAME the driver    → "You're using your People Power right now."
2. VALIDATE the energy → "It feels big when I look away, doesn't it?"
3. STRESS-TEST safely  → Activity that loads the threat channel in micro-doses
4. DEBRIEF             → "How did that feel in your tummy?"
```

Over time, the child learns: **"I HAVE this feeling, but I am NOT this feeling."**

### Stress Tests (One Per Driver)

Each driver has a specific stress-test activity that safely loads its anxiety channel:

| Driver | Stress Test | Deprivation Target | Core Skill Built |
|--------|------------|-------------------|-----------------|
| Impressor | "The Invisible Artist" — create while parent looks away | Parent's gaze | Internal reward |
| Explorer | "The Pause Button" — stop mid-investigation and verbalise | Uninterrupted exploration | Disengagement |
| Sentinel | "The Tiny Surprise" — one change in a familiar routine | Absolute predictability | Change tolerance |
| Connector | "The Kind Thought" — send a thought instead of fixing | Permission to physically help | Emotional boundaries |
| Rebel | "The Two Doors" — choose within constraints | Unlimited choice | Agency within structure |
| Performer | "The Whisper Game" — exciting activity in whispers | Volume as outlet | Voltage containment |
| Observer | "The Body First Game" — jump before understanding | Full analysis first | Action-before-analysis |

### Driver × Activity Interactions

Activities are now tagged with `DriverInteraction` entries that describe how each relevant driver experiences the activity differently. The same Dance Party is:

- **Reward** for an Impressor (audience is right there)
- **Threat** for an Observer (must act before computing)
- **Reward** for a Performer (full-voltage expression)

This allows the book to include driver-specific parent notes, and the future parent app to personalise activity recommendations based on detected driver patterns.

### Driver Detection (App Layer)

A `DetectLikelyDrivers()` function takes spectrum observations (accumulated over 3+ theme cycles) and returns ranked driver matches. The algorithm weights `primary_fuel` spectra 3×, `amplifier` 2×, `filter` 1×, then normalises to a 0–1 confidence score. The parent never sees "driver detection" — they see warm framing: *"Your child seems to light up when..."*

### 3.10 Assessment Feedback Loop — "The Control Loop"

**File:** `assessment_loop.go` (~550 lines)

The assessment layer closes the engineering control loop. Without it, activities are *open-loop* — fire and forget. With it, every observation patches the next recommendation.

#### Control Theory Frame

```
INPUT              TRANSFER FUNCTION           OUTPUT
(Activity)  →  (Child's Spectra + Drivers)  →  (Consequence)
     ↑                                              │
     │              FEEDBACK SENSOR                  │
     └──────────── (Assessment) ←───────────────────┘
```

#### The Three-Check Sensor

Every post-activity observation reads three dimensions:

| # | Sensor | Measures | Question Frame | Maps To |
|---|--------|----------|----------------|---------|
| 1 | **Recovery** | Latency to Baseline | "How quickly did they bounce back?" | `pilot_strength` (global moderator) |
| 2 | **Friction** | Voltage Resistance | "Was it a grind or a flow?" | Spectrum alignment (per-spectrum) |
| 3 | **Agency** | Internal Locus | "Who did they look at for the win?" | Driver moderation (ventral drivers) |

Each sensor has **9 spectrum-themed variants** — one per theme week. The parent sees 3 questions per activity (one per sensor), phrased in everyday language. All responses map to a -1.0 to +1.0 scale internally.

#### 9 Parent App Observation Questions (per sensor)

Each spectrum theme generates three observation prompts:

| Spectrum | Recovery Question | Friction Question | Agency Question |
|----------|-------------------|-------------------|-----------------|
| 1. Near and Far | "When togetherness got too much, how quickly did they bounce back?" | "Did the togetherness feel natural or forced?" | "When they did something good, who did they look at?" |
| 2. Loud and Quiet | "When noise changed, how fast did they adjust?" | "Was the energy level a good fit?" | "Was engagement from inside or outside?" |
| 3. Big Feelings | "After the biggest feeling, how long until normal?" | "Was the emotional intensity about right?" | "Did they self-soothe or need you to fix it?" |
| 4. Yes and No | "When they hit a 'no' moment, how quickly did they find their feet?" | "Did boundary moments feel age-appropriate?" | "Was their yes/no THEIR choice or a reaction?" |
| 5. Mine and Yours | "When sharing was hard, how long before re-engaging?" | "Was sharing the right level of challenge?" | "Was sharing genuine or performed for you?" |
| 6. Try and Wait | "When it felt risky, how quickly did they settle in?" | "Was the risk level right today?" | "Were they trying for themselves or your reaction?" |
| 7. Same and Different | "When rules changed, how quickly did they adapt?" | "Was the flexibility the right challenge?" | "Did they celebrate internally or need validation?" |
| 8. Your/My Feelings | "If they picked up someone's feeling, how quickly did they let go?" | "Was the emotional exposure manageable?" | "Was helping their compass or echoing yours?" |
| 9. Tidy and Messy | "When things got messy/tidy, how long before they felt OK?" | "Was the structure level comfortable?" | "Was tidying their preference or compliance?" |

#### Noise Filtering — Hardware vs Software

Before processing any sensor signal, the system checks **Current Load** (the existing `current_load` global moderator in `spectra.go`):

| Load Level | Weight | System Response |
|------------|--------|-----------------|
| **Low** (rested, fed, well) | 100% | Full signal — process normally |
| **Medium** (bit tired, recovering) | 50% | Reduced weight — half signal |
| **High** (exhausted, hungry, ill) | 0% | **Noise** — discard entirely. "Bad days are data about sleep, not about your child." |

This prevents the system from confusing a hungry toddler's meltdown with a genuine spectrum-alignment problem.

#### Processing Pipeline

`ProcessObservation()` runs four steps:

1. **Noise filter** — if `CurrentLoad.Level == "high"`, mark as noise, suggest rest, return
2. **Spectrum insights** — Recovery → pilot_strength delta; Friction → primary spectrum delta; Agency → internal_locus delta
3. **Driver adjustments** — Low agency → Impressor +0.1; High agency → Explorer/Observer +0.05; High friction on orderliness → Sentinel +0.1; Slow recovery → Performer +0.05
4. **Patch generation** — Priority order: recovery failure → "increase_predictability"; high friction → "simplify"; low agency → "reduce_audience"

#### Patch Types

| Patch | When | Parent Message (excerpt) |
|-------|------|--------------------------|
| `rest` | High current load | "Battery was low today. This one doesn't count." |
| `increase_predictability` | Recovery ≤ -0.5 | "Try adding a countdown before transitions." |
| `simplify` | Friction ≤ -0.5 | "Dial it back — shorter, more scaffolding." |
| `reduce_audience` | Agency ≤ -0.5 | "Step back physically. See if they keep going." |
| `extend_duration` | Mild recovery wobble | "Try the same activity again soon. Repetition builds the pathway." |
| `continue` | All sensors positive | "That went well. Keep going." |

#### Range of Motion — The Zen Warrior Metric

"Smooth running" is NOT quiet or obedient. It is **Range of Motion**.

```
LOW RANGE:  Child can ONLY function at preferred pole.
            ┌──●────────────────────┐  ← locked to one end
            
HIGH RANGE: Child PREFERS one pole but maintains regulation at the other.
            ┌────────●──────────────┐  ← comfortable across the range
```

`RangeOfMotion` is computed from historical friction values per spectrum. Each spectrum gets a 0.0–1.0 score:
- **0.0** = locked to one pole (only "tidy" OR only "messy")
- **0.5** = can tolerate opposite pole with support
- **1.0** = comfortable at both poles independently

The **Progress Rosette** is a radial visualisation with 9 axes (one per spectrum). The frontier expands outward as Range of Motion increases. A full circle = Zen Warrior. Growth is measured by **area increase**, not symmetry — an asymmetric rosette that's bigger than last month is progress.

#### Key Structs

| Struct | Purpose |
|--------|---------|
| `ObservationEvent` | Atomic unit: one activity, one child, three sensor readings + noise check |
| `LoadState` | Hardware noise filter: tired/hungry/ill flags |
| `SensorReading` | Single sensor response: question + response + normalised value |
| `SpectrumInsight` | Derived spectrum adjustment (delta + source sensor) |
| `DriverAdjustment` | Driver confidence shift from observation |
| `ActivityPatch` | Next-session recommendation: type + parent message + duration |
| `RangeOfMotion` | Per-child rosette: 9 spectrum range scores + overall + trend |

### 3.11 Pathway Flexibility — "The Differentiable Curriculum"

**File:** [`pathway_flexibility.go`](../esp-organizer/internal/domain/etp/pathway_flexibility.go) (~900 lines)

The book doesn't change. The path through it does. This is the **differentiable curriculum** — one physical object that adapts to each child through **choice architecture**, not code.

| Element | Fixed | Differentiated |
|---------|-------|----------------|
| Pages | Same 32 activities | Order varies by child |
| Activities | Same instructions | Parent language shifts by ETP |
| Stickers | Same 32 stickers | Which ones you celebrate varies |
| Difficulty | Same 3 cycles | Which cycle you're on varies |

**The book is the hardware. The parent is the operating system. The app is the compiler.**

#### Non-Linear Navigation

Each spread has "Try this when" and "Skip if" prompts printed on the page. The parent scans these and chooses — no prescribed order:

```
Try this when:
  • Your child is learning to set boundaries
  • You're noticing power struggles
  • They need practice saying "no"

Skip if:
  • They're tired or hungry
  • You're short on time
  • They just had a big "no" fight
```

The app later asks: "Did that match?" — feeding the suggestion engine.

#### Difficulty Cycles as Layers

Each of the 9 themes has 3 difficulty cycles. They're **layers**, not sequential pages:

| Spectrum | Cycle 1 | Cycle 2 | Cycle 3 |
|----------|---------|---------|---------|
| Near and Far | Side by Side (parallel play) | Invitation to Join | The Push and Pull |
| Loud and Quiet | Volume Dial | Matching Energy | The Whisper Challenge |
| Big Feelings | Feelings Thermometer | The Wave Rider | The Feeling DJ |
| Yes and No | Silly Questions | Real Choices | Boundary Negotiation |
| Mine and Yours | What's Mine | Lending Library | The Gift |
| Try and Wait | Tiny Bravery | The Warm-Up Lap | The Wobbly Bridge |
| Same and Different | The Rules Game | The Exception | The Referee |
| Your/My Feelings | Whose Feeling Is This? | The Feeling Catcher | The Kind Thought |
| Tidy and Messy | Exploring Both | My Way | The Controlled Chaos |

The app tracks which cycle the child is on per theme. `SuggestCycle()` looks at the last 3 attempts: all easy → escalate, repeated struggle → step back, mixed → stay.

#### Sticker Choice Architecture

Instead of one sticker per week, the sticker sheet has **4 options per theme**:

```
This week, celebrate:
🎉 "I said a big YES"
🛑 "I said a clear NO"
😂 "I laughed at a silly question"
🔄 "I changed my mind"
```

Parent chooses which fits. The app learns which celebrations resonate — different sticker choices reveal driver patterns (Impressor children always pick the audience-facing sticker).

#### "Try Instead" Variations

Every activity has 3 variations printed directly on the page:

| Condition | Adjustment | Example |
|-----------|------------|---------|
| "If they loved this" | Escalate | "Next time, try it with stuffed animals" |
| "If they struggled" | Simplify | "Next time, just do the first step and stop" |
| "If they were overwhelmed" | Scaffold | "Next time, just watch you do it first" |

No app needed — the adaptations are in the ink. The app notes which variation was chosen.

#### Suggestion Engine

`SuggestNextThemes()` returns 3 recommendations by scoring all 9 themes on:

1. **Coverage** — least practised themes score higher (+3.0 if untried, +2.0 if done once)
2. **Recency** — themes not attempted in 4+ weeks get a +1.0 boost
3. **Driver alignment** — themes matching detected ventral drivers get +0.5 × confidence
4. **Load adjustment** — on high-load days, prefer aligned themes at cycle 1; de-prioritise stretch themes

Output is a warm suggestion card, not a prescription:

```
THIS WEEK'S SUGGESTIONS
┌────────────────────────────────┐
│ 🛑 YES AND NO (Cycle 1)        │
│ They're learning boundaries.   │
│ Try the silly questions.       │
├────────────────────────────────┤
│ 🌪️ TIDY AND MESSY (Cycle 2)    │
│ Haven't done this in 4 weeks.  │
│ Try the "tidy YOUR way" one.   │
├────────────────────────────────┤
│ 🌊 DANCE PARTY (Repeat)         │
│ They loved this. Do it again!  │
│ Repetition builds pathways.    │
└────────────────────────────────┘
```

Parent picks one. App logs. Algorithm updates. Nobody is told what to do.

#### Choice Logging

Every parent decision is captured in `PathwayChoice`:

| Field | Captures |
|-------|----------|
| `ThemeWeekID` | Which theme + cycle they chose |
| `ActivityID` | Which activity they did |
| `ChosenVariation` | Which "try instead" they picked |
| `ChosenCelebration` | Which sticker they used |
| `DifficultyRating` | 1=easy, 2=about right, 3=hard |
| `ObservationID` | Link to Three-Check Sensor reading (if logged) |

Over time, the choice log reveals: which themes the parent avoids (gap signal), which variations they default to (difficulty calibration), which stickers resonate (driver signal). The book doesn't change — the experience always does.

#### Theme Week Data

All 36 themed weeks (12 spectra × 3 cycles) are seeded in `ThemeWeeks` with full content:

- Navigation prompts (try this when / skip if)
- Activities with 3 printed variations each
- 4 sticker choices per theme, each tagged with `SkillFocus`

Every activity, variation, and sticker is linked back to the skill graph and ETP framework — the same data that feeds the pathfinder, the rosette, and the driver detector.

#### Key Structs

| Struct | Purpose |
|--------|---------|
| `ThemeWeek` | One spectrum theme at one difficulty cycle: activities + stickers + navigation prompts |
| `ThemeActivity` | Activity within a theme week: instructions + 3 "try instead" variations |
| `ActivityVariation` | Printed adaptation: condition + suggestion + adjustment type |
| `StickerChoice` | One celebration option: icon + label + skill focus |
| `PathwayChoice` | Parent decision log: what they chose + how it went |
| `ThemeSuggestion` | App recommendation: theme + cycle + reason + priority |

---

## 4. Connection to the Skills Map

### 4.1 The Universal Skills Graph

The skills map is a directed acyclic graph (DAG) stored across two systems:

| System | Database | Schema | Purpose |
|--------|----------|--------|---------|
| **Skills Map Platform** | MySQL (`dare2lead`) | `skills_key` + `skill_link` | Operational store — CRUD, scoring (1–5), graph traversal, visualisation |
| **CHISG (Weaviate 8081)** | Weaviate | `CHISGElement` + `SkillLink` | Semantic search, ETP modulation tags, relationship inference |

#### MySQL Schema (`skills_key`)

```sql
CREATE TABLE skills_key (
    Skills_keyID     bigint PRIMARY KEY AUTO_INCREMENT,
    Skill_name       varchar(50) NOT NULL,
    Criteria1        varchar(500) DEFAULT '',  -- Level 1: Not yet observed
    Criteria2        varchar(500) DEFAULT '',  -- Level 2: Glimmers
    Criteria3        varchar(500) DEFAULT '',  -- Level 3: With support
    Criteria4        varchar(500) DEFAULT '',  -- Level 4: Independent
    Criteria5        varchar(500) DEFAULT '',  -- Level 5: Shows others
    Hidden           tinyint DEFAULT 0,
    DevelopmentAge   double DEFAULT 0,         -- Age in years (0.0–18.0+)
    SkillDescription varchar(500) NOT NULL
);
```

#### MySQL Schema (`skill_link`)

```sql
CREATE TABLE skill_link (
    Skill_linkID      bigint PRIMARY KEY AUTO_INCREMENT,
    ParentSkillID     bigint DEFAULT 0,    -- The prerequisite skill
    OffspringSkillID  bigint DEFAULT 0,    -- The skill that depends on it
    CourseID          bigint DEFAULT 0     -- Context (which lens)
);
```

#### Weaviate Schema (`CHISGElement`)

```
CHISGElement:
    name              TEXT
    description       TEXT
    domain            TEXT       (physical, social-emotional, cognitive, linguistic)
    layer             TEXT       (foundational, intermediate, academic)
    etp_spectrum_id   TEXT       (which of the 9 spectra this skill loads on)
    etp_mechanism     TEXT       (how ETP modulates this skill)
    etp_energy_cost   TEXT       (voltage cost for different settings)
    enables           TEXT[]     (forward skill links)
    requires          TEXT[]     (backward skill links)
```

### 4.2 ToddlerOS as a Lens

ToddlerOS queries the same graph as every other product, filtered by:

```
WHERE DevelopmentAge BETWEEN 0.0 AND 5.0
AND   domain IN ('physical', 'social-emotional', 'cognitive', 'linguistic')
AND   layer = 'foundational'
```

This returns the ~30–50 skill nodes that form the base of the DAG. Higher-age products (PrimaryOS, LAO, CareerOS) query upward from these same nodes.

### 4.3 The 7 High-Leverage Foundational Capacities

These are the **root nodes** of the skill graph — the capacities that predict downstream outcomes across all the developmental research (Hart & Risley, Heckman, marshmallow replication studies):

| Capacity | Domain | DevelopmentAge | ETP Spectrum Link | Practise Time |
|----------|--------|----------------|-------------------|---------------|
| **Inhibitory control** | cognitive | 1.5–4.0 | `pilot_strength` (moderator) | 2 min (stop-go games) |
| **Joint attention** | social-emotional | 0.5–2.0 | `social_gravity`, `mirror_neuron_tuning` | 1 min (pointing, shared book) |
| **Emotional labelling** | social-emotional | 1.0–3.0 | `voltage_sensitivity`, `mirror_neuron_tuning` | 30 sec ("You look frustrated") |
| **Turn-taking** | social-emotional | 1.0–3.0 | `social_gravity`, `care_response` | 2 min (rolling ball back and forth) |
| **Sustained attention** | cognitive | 0.5–5.0 | `energy_directionality`, `orderliness` | 0 min (protect by not interrupting) |
| **Verbal interaction density** | linguistic | 0.0–5.0 | `energy_directionality`, `social_gravity` | 5 min woven into any activity |
| **Frustration tolerance** | cognitive | 1.5–4.0 | `threat_response`, `risk_tolerance` | 1 min (wait before helping) |

Each capacity will be ingested as a `skills_key` row with `DevelopmentAge` in the 0–5 range, `Criteria1-5` filled with age-appropriate observation descriptors, and `skill_link` edges connecting upward to the academic skills they enable.

### 4.4 Criteria 1–5 for Early Years (Steiner-Compatible)

The 1–5 scoring scale was originally developed with Kevin Avison of the Steiner Waldorf Advisory Service for Ofsted observation reporting. The levels translate naturally to early years:

| Level | Generic | ToddlerOS (Parent Language) | Steiner Translation |
|-------|---------|----------------------------|---------------------|
| 1 | Not yet observed | "Haven't seen this yet" | Not yet observed |
| 2 | Glimmers | "Starting to notice glimpses" | Glimmers |
| 3 | With support | "Can do it with help" | With support |
| 4 | Independent | "Does it on their own" | Independent |
| 5 | Shows others | "Helps other children do it" | Shows others |

In the parent app, these aren't presented as numbers. The observation prompts ("Which was easier — near or far?") implicitly place the child at a point on the scale. The system infers the score.

### 4.5 How Activities Connect to Skill Nodes

The `skillsenable` table maps activities to skills:

```sql
CREATE TABLE skillsenable (
    SkillsEnableID  bigint PRIMARY KEY AUTO_INCREMENT,
    Activity        varchar(50) NOT NULL,  -- e.g., "stop-go game"
    Skills_keyID    bigint DEFAULT 0       -- → skills_key.Skills_keyID
);
```

For ToddlerOS, each `PlayAvenue` from the response matrix is an `Activity` row linked to the foundational skill nodes it develops. This creates a queryable chain:

```
ETP Spectrum → ResponseMatrixEntry.PlayAvenues → skillsenable.Activity → skills_key node
```

Example:
```
social_gravity (independent setting)
  → PlayAvenue: "Parallel play with gradual invitation"
    → skillsenable: Activity="Parallel play" → Skills_keyID=<joint_attention>
    → skillsenable: Activity="Parallel play" → Skills_keyID=<social_gravity_tolerance>
```

---

## 5. Connection to the Extraction Pipeline

### 5.1 CHISG Extraction (Existing)

The extraction pipeline (`extract_chisg_links.py`) uses Bedrock/Sonnet to extract semantic relationships from source definitions. It has already processed the LAO GCSE science corpus:

- **Input**: 2,746 keyword definitions from `LAOMobile/backend/lao.db`
- **Output**: 6,415 semantic links in [`extraction_full_sonnet.json`](../data/chisg/extraction_full_sonnet.json)
- **Controlled vocabulary**: 13 relation types (causes, enables, inhibits, is composed of, is a type of, has property, has value, is used for, is found in, is an example of, determines, represents, contradicts)
- **Ingested into**: Weaviate 8081 as `CHISGElement` (579 nodes) + `SkillLink` (1,701 edges)

### 5.2 Early Years Extraction (Planned)

The same pipeline will process early years source documents:

| Source | Content | Expected Yield |
|--------|---------|----------------|
| EYFS Statutory Framework | 7 areas of learning, 17 early learning goals | ~50 skill definitions |
| Development Matters | Age-band descriptors (birth–5) | ~150 observation statements |
| Montessori scope-and-sequence | Practical life, sensorial, language, maths, cultural | ~80 skill progressions |
| Steiner curriculum (Avison) | Kindergarten capacities, will-based learning | ~40 foundational nodes |

The extraction will produce:
1. New `skills_key` rows with `DevelopmentAge` 0.0–5.0
2. New `skill_link` edges connecting early years nodes to each other and upward to existing academic nodes
3. New `CHISGElement` entries in Weaviate 8081 with `etp_spectrum_id` tags
4. Activity-to-skill mappings in `skillsenable`

### 5.3 ETP Modulation Tags

Each skill node in CHISG carries ETP modulation metadata:

| Field | Purpose | Example |
|-------|---------|---------|
| `etp_spectrum_id` | Which spectrum most affects this skill | `4` (threat_response) for "frustration tolerance" |
| `etp_mechanism` | How the spectrum interacts | "High threat-response passive setting increases freeze probability when frustrated" |
| `etp_energy_cost` | Voltage cost for different settings | "High cost for passive (avoidance voltage), low cost for aggressive (confrontation voltage)" |

These tags allow the pathfinder engine to weight graph edges differently per child's ETP profile.

---

## 6. The Pathfinder Engine

### 6.1 Shared Infrastructure

The pathfinder is not a ToddlerOS component — it is a **shared engine** consumed by all products:

| Product | Pathfinder Query | Output |
|---------|-----------------|--------|
| **ToddlerOS** | "Given this child's current skill scores and ETP profile, which theme week should we suggest?" | Weekly theme recommendation |
| **PrimaryOS** | "Which skill should this student work on next?" | Next skill in progression |
| **LAO** | "Which revision topic gives the most marks-per-hour?" | Revision order |
| **CareerOS** | "How far is this person from this career's skill requirements?" | Distance-to-competency metric |

### 6.2 Algorithm (Planned)

```
pathfinder(
    current_scores:  map[skill_id → 1–5 score],
    destination:     set[skill_id],      // target skills (or "all" for ToddlerOS)
    etp_profile:     map[spectrum → -2..+2 value],
    constraints:     { max_time, max_energy, ... }
) → ordered list of (skill_id, priority_score, suggested_activity)
```

**Edge weight calculation**:
```
base_weight = difficulty(skill) × (5 - current_score) / 5
etp_modifier = energy_cost(skill.etp_spectrum_id, etp_profile[skill.etp_spectrum_id])
adjusted_weight = base_weight × etp_modifier
```

For a child with strong `threat_response = passive` (-1.5), skills that require confrontation have a higher `etp_modifier` — they're not avoided, but scheduled when pilot_strength is high (rested, fed, good day). For a child with `orderliness = ordered` (+1.8), spontaneity-heavy activities are weighted to follow structured ones.

### 6.3 ToddlerOS-Specific Pathfinding

In the dip-in/dip-out model, the pathfinder doesn't enforce a sequence. Instead it:

1. Ranks the 9 spectrum themes by current need (least-practised + highest developmental impact)
2. Adjusts for ETP profile (themes aligned with natural settings first to build confidence, cross-pole themes scheduled strategically)
3. Suggests a theme without mandating it
4. Presents as "You might try..." not "You must do..."

The parent can ignore the suggestion. The book works without the app. The app enhances without being required.

---

## 7. Data Flow Summary

```
┌────────────────────────────────────────────────────────────────────────┐
│                        SOURCE MATERIAL                                 │
│  EYFS · Development Matters · Montessori · Steiner/Avison             │
└──────────────────────────────┬─────────────────────────────────────────┘
                               │ extract_chisg_links.py (Bedrock/Sonnet)
                               ▼
┌────────────────────────────────────────────────────────────────────────┐
│                     UNIVERSAL SKILLS GRAPH                             │
│  MySQL: skills_key (nodes) + skill_link (DAG edges)                   │
│  Weaviate 8081: CHISGElement + SkillLink + ETP modulation tags        │
│  MongoDB: Skill model (name, description, category, criteria[])       │
│  ~30–50 foundational nodes (age 0–5) at base of graph                 │
└──────────────┬───────────────────────────────────────┬────────────────┘
               │                                       │
               ▼                                       ▼
┌──────────────────────────────┐   ┌────────────────────────────────────┐
│      ETP FRAMEWORK           │   │     SCORING & OBSERVATION          │
│  9 spectra (spectra.go)      │   │  skillsmarkbook (1–5 per skill)   │
│  47 trainable skills         │   │  Parent observations (app)        │
│  18 response matrix entries  │   │  Implicit Criteria1–5 placement   │
│  2 global moderators         │   │                                    │
│  Voltage calculator          │   │                                    │
└──────────────┬───────────────┘   └──────────────┬─────────────────────┘
               │                                   │
               ▼                                   ▼
┌────────────────────────────────────────────────────────────────────────┐
│                     PATHFINDER ENGINE (Go)                              │
│  Input: current scores + destination + ETP profile                     │
│  Output: ordered (skill, priority, activity) tuples                    │
│  ETP modulation: edge weights adjusted per spectrum alignment          │
└──────────────────────────────┬─────────────────────────────────────────┘
                               │
                               ▼
┌────────────────────────────────────────────────────────────────────────┐
│                     TODDLEROS OUTPUT                                   │
│                                                                        │
│  ┌──────────────────────┐  ┌────────────────────────────────────────┐  │
│  │  ACTIVITY BOOK (PDF) │  │  PARENT APP (Companion)               │  │
│  │  Seasonal editions   │  │  Observation log per theme            │  │
│  │  16 spreads each     │  │  Progress rosette (not checklist)     │  │
│  │  "With You" mode     │  │  ETP profile builds over 32 weeks    │  │
│  │  "While You" mode    │  │  Theme suggestions (pathfinder)       │  │
│  │  Paper, not screen   │  │  Parent-only — no child interface     │  │
│  └──────────────────────┘  └────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 8. File Reference

### ETP Domain Logic

| File | Lines | Content |
|------|-------|---------|
| [`spectra.go`](../esp-organizer/internal/domain/etp/spectra.go) | 347 | 9 spectra definitions, 47 trainable skills, 2 global moderators, lookup functions |
| [`response_matrix.go`](../esp-organizer/internal/domain/etp/response_matrix.go) | 576 | 18 response matrix entries, 3 core barriers, language shifts, play-first principle, `GeneratePersonalisedPlan()`, `GetAdultLanguageForProfile()` |
| [`activity_experience.go`](../esp-organizer/internal/domain/etp/activity_experience.go) | ~380 | Activity Experience model (Emotional Compiler) — multi-threaded activity schema with spectrum activations, side effects, counterweight strategies, driver interactions, book spread content. 3 sample activities. |
| [`ventral_drivers.go`](../esp-organizer/internal/domain/etp/ventral_drivers.go) | ~460 | Ventral Emotional Drivers — 7 archetypal driver personas with spectrum signatures, heat shield mechanisms, stress tests, detection algorithm. |
| [`assessment_loop.go`](../esp-organizer/internal/domain/etp/assessment_loop.go) | ~550 | Assessment Feedback Loop — Three-Check Sensor (Recovery/Friction/Agency), 27 observation questions (9 spectra × 3 sensors), noise filtering, Range of Motion metric, `ProcessObservation()`, `CalculateRangeOfMotion()` |
| [`pathway_flexibility.go`](../esp-organizer/internal/domain/etp/pathway_flexibility.go) | ~900 | Differentiable Curriculum — 27 themed weeks (9 spectra × 3 cycles), non-linear navigation, sticker choice architecture, "try instead" variations, `SuggestNextThemes()`, `SuggestCycle()`, choice logging |
| [`voltage_calculator.go`](../esp-organizer/internal/domain/etp/voltage_calculator.go) | 318 | Voltage impact calculations, energy cost functions |
| [`compatibility.go`](../esp-organizer/internal/domain/etp/compatibility.go) | 192 | Compatibility solutions for all spectra |

### Skills Map Platform

| File | Content |
|------|---------|
| [`init.sql`](../../skills-map-platform/init.sql) | Full MySQL schema — `skills_key`, `skill_link`, `skillsmarkbook`, `skillsenable`, `skillsmapelements` |
| [`structures.go`](../../skills-map-platform/api/structures/structures.go) | Go structs: `Skill` (with Criteria1–5, DevelopmentAge, ParentSkills, OffspringSkills), `SkillNode`, `SkillLink` |
| [`skill.go`](../esp-organizer/internal/models/skill.go) | MongoDB Skill model: Name, Description, Category, DevelopmentAge, ParentSkillIDs, ChildSkillIDs, SkillCriteria[] |

### CHISG / Extraction

| File | Content |
|------|---------|
| [`extraction_full_sonnet.json`](../data/chisg/extraction_full_sonnet.json) | 6,415 semantic links from 2,746 GCSE definitions (Bedrock/Sonnet extraction) |
| [`ingest_chisg_complete.py`](../../humanOS/scripts/ingest_chisg_complete.py) | Weaviate 8081 ingestion script — CHISGElement + SkillLink schemas with ETP modulation fields |
| [`extract_chisg_links.py`](../scripts/extract_chisg_links.py) | LLM extraction pipeline (controlled vocabulary of 13 relations) |

### Product & Strategy

| File | Content |
|------|---------|
| [`PROJECT_STATUS.md`](PROJECT_STATUS.md) | Full product design, strategic decisions, build priorities |
| [`PORTFOLIO_SITE_OVERVIEW.md`](PORTFOLIO_SITE_OVERVIEW.md) | Live site routes including ToddlerOS landing page |
| [`ToddlerOSLanding.tsx`](../frontend/src/pages/ESPWorld/ToddlerOSLanding.tsx) | 1,010-line landing page with 9 parent-friendly spectrum questions |

---

## 9. What Exists vs. What's Planned

| Component | Status | Location |
|-----------|--------|----------|
| ETP spectra definitions | ✅ Complete | `spectra.go` (compiled Go) |
| Response matrix (18 entries) | ✅ Complete | `response_matrix.go` (compiled Go) |
| Voltage calculator | ✅ Complete | `voltage_calculator.go` |
| `GeneratePersonalisedPlan()` | ✅ Complete | `response_matrix.go` |
| Activity Experience schema | ✅ Complete | `activity_experience.go` — 9 structs, 3 sample activities, side effect model, counterweight strategies, driver interactions, book spread template |
| Ventral Emotional Drivers | ✅ Complete | `ventral_drivers.go` — 7 archetypal personas (Impressor, Explorer, Sentinel, Connector, Rebel, Performer, Observer), spectrum signatures, stress tests, `DetectLikelyDrivers()` |
| Assessment Feedback Loop | ✅ Complete | `assessment_loop.go` — Three-Check Sensor (Recovery/Friction/Agency), 27 observation questions, noise filtering (current_load), `ProcessObservation()`, `CalculateRangeOfMotion()`, Progress Rosette, activity patches |
| Pathway Flexibility | ✅ Complete | `pathway_flexibility.go` — 27 themed weeks (9 spectra × 3 difficulty cycles), non-linear navigation prompts, 4 sticker choices per theme, 3 activity variations per spread, `SuggestNextThemes()`, `SuggestCycle()`, `PathwayChoice` logging |
| Skills graph schema | ✅ Complete | `init.sql` (MySQL), Weaviate 8081 |
| CHISG extraction pipeline | ✅ Complete | `extract_chisg_links.py` |
| CHISG GCSE data | ✅ Complete | 579 elements, 1,701 links in Weaviate |
| Landing page (ToddlerOS) | ✅ Complete | `ToddlerOSLanding.tsx` |
| Foundational skill nodes (0–5) | ⬜ Not built | Need ~30–50 nodes from developmental sources |
| Early years extraction | ⬜ Not built | Same pipeline, different source documents |
| Pathfinder engine | ⬜ Not built | Go service — shared by all products |
| Activity book PDF generator | ⬜ Not built | Template → skill query → PDF |
| Parent observation app | ⬜ Not built | Observation log → implicit ETP profile |
| ETP-weighted edge costs | ⬜ Not built | `etp_spectrum_id` fields exist, weighting logic doesn't |

---

## 10. Build Sequence

1. **Ingest foundational skill nodes** — Add 30–50 early years nodes to `skills_key` with `DevelopmentAge` 0.0–5.0, Criteria1–5, and `skill_link` edges. Tag each with `etp_spectrum_id` in Weaviate.

2. **Build pathfinder engine** — Go service: shortest weighted path through DAG with ETP-adjusted edge costs. Query interface serves all products.

3. **Map PlayAvenues to skill nodes** — Populate `skillsenable` with response_matrix.PlayAvenues → foundational skill nodes. Creates the chain: spectrum → activity → skill.

4. **Build activity book generator** — Template engine: theme → query pathfinder → slot activities into spread templates → produce PDF. 16 spreads per seasonal edition.

5. **Build parent observation app** — Log observations per theme. Infer Criteria1–5 placement from natural-language prompts. Build ETP profile over 32 weeks. Progress rosette visualisation.
