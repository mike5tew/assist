# ToddlerOS — Booklet Outlines

**Last Updated**: 2026-03-02  
**Status**: Outlines reconstructed from conversation history and compiled Go data  
**Source of truth**: `pathway_flexibility.go` (1,361 lines) — all 36 theme weeks with full activity content

---

## 1. Product Format

ToddlerOS is delivered as a **physical activity book** (printed card stock + sticker sheet) with a digital **companion parent app** (observation log + progress rosette).

| Element | Detail |
|---------|--------|
| **Format** | Printed card stock, saddle-stitched or spiral-bound |
| **Edition model** | Seasonal dip-in/dip-out (any page is the right page) |
| **Pages per edition** | 16 spreads (32 pages) per seasonal booklet |
| **Total programme** | 36 themed weeks (12 spectra × 3 difficulty cycles) + seasonal specials |
| **Sticker sheet** | Included with each booklet. 4 sticker choices per theme. |
| **Revenue model** | Free PDF download → email capture. Printed subscription → revenue. |
| **No child screen time** | The child's interface is paper. The parent's interface is a phone — used only for logging, never shown to the child. |

### Two Modes Per Spread

| Mode | Duration | Who's Active | Purpose |
|------|----------|-------------|---------|
| **"With You"** | 2–5 min | Parent + child together | Builds foundational capacities through guided play |
| **"While You..."** | 5–15 min | Child independent | Colouring, stickers, audio via QR — child processes while parent steps back |

### Booklet Intro Page — "The Dance"

The first page of every edition. Read-aloud for the parent — not the child. Sets the emotional frame before any activity begins.

---

> **Before you start**
>
> Your child has a way of being in the world. So do you.
>
> Sometimes it's easy. Sometimes it's not. Some days they're brave and curious and kind. Some days they scream in the cereal aisle. Some days you handle it perfectly. Some days you don't.
>
> This book isn't about fixing any of that.
>
> There's a musician called Ren who said something that stuck with us. He'd been through a long fight with his own mind, and he said this:
>
> *"It was never a battle I was supposed to win. It is an eternal dance, and just like any dance, the more rigid I became, the harder it got. The more I cursed my clumsy footsteps, the more I struggled. So I got older, and I learned to relax, and soften, and that dance got easier."*
>
> That's what this book is about.
>
> Your child has **settings** — biological ones. How big their feelings are. Whether people charge their battery or drain it. Whether they jump in or watch first. Whether rules feel safe or suffocating. These aren't problems. They're just settings. Every child has them. Every adult does too.
>
> Your child also has **skills** — things they can do fluently, and things they're still building. Not compared to other children. Not measured against "normal." Just: what can they do, and what's next?
>
> The settings and the skills together make a pattern. That pattern is unique to your child. It's not broken. It's not behind. It's theirs.
>
> Some patterns have clinical names — autism, ADHD, dyslexia. Those names are useful. They open doors. But they don't tell you what to do on a Tuesday afternoon. This book does.
>
> **The activities in here aren't tests.** There's no pass or fail. They're invitations — little experiments where you and your child play together, and both of you notice something. What was easy? What was hard? What made them light up? What made them shut down?
>
> You're not trying to change your child's settings. You're learning the rhythm of their swing.
>
> Some days the dance will be beautiful. Some days you'll both step on each other's feet. That's not failure — that's the dance.
>
> **The only rule:** if it's a bad day, close the book. Read them a story instead. The dance will be there tomorrow.

*Design note: Warm cream background. No illustrations — just text. Generous margins. This page breathes. It should feel like a letter, not an instruction manual.*

---

## 2. Spread Template (Page Layout)

Every activity compiles down to this two-page spread:

```
┌────────────────────────────────────────────────────────────────┐
│ LEFT PAGE — "WITH YOU" ACTIVITY                                │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  [Icon] ACTIVITY NAME                       ⏱ Duration  │  │
│  │                                             🔧 Materials │  │
│  ├──────────────────────────────────────────────────────────┤  │
│  │                                                          │  │
│  │  [Base instruction — 2-3 sentences]                      │  │
│  │                                                          │  │
│  ├──────────────────────────────────────────────────────────┤  │
│  │  👀 WHAT TO NOTICE                                       │  │
│  │  [Observation prompt → maps to spectrum]                  │  │
│  ├──────────────────────────────────────────────────────────┤  │
│  │  🛟 IF THEY STOP / STRUGGLE...                           │  │
│  │  [Rescue language + voltage grounding prompt]             │  │
│  ├──────────────────────────────────────────────────────────┤  │
│  │  ✅ TRY: "[RecommendedLanguage from response matrix]"    │  │
│  │  ❌ NOT: "[AvoidLanguage from response matrix]"           │  │
│  ├──────────────────────────────────────────────────────────┤  │
│  │  🔄 NEXT STEP (3 pathways):                               │  │
│  │  ⬆ ESCALATE → [successor skill on graph]                 │  │
│  │    "They showed range — step forward."                    │  │
│  │  ⬇ SIMPLIFY → [precursor skill on graph]                 │  │
│  │    "Not ready — drop to the foundation this needs."       │  │
│  │  ⚡ VOLTAGE SPIKE → [ground, don't redirect]              │  │
│  │    "This isn't a skill gap — their conductor overloaded.  │  │
│  │     Stop the activity. Use voltage grounding: name the    │  │
│  │     feeling, slow the breathing, wait for the wave."      │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                │
│  Navigation prompts:                                           │
│  "Try this when:" [readiness signals]                          │
│  "Skip if:" [noise/timing signals]                             │
│                                                                │
├────────────────────────────────────────────────────────────────┤
│ RIGHT PAGE — "WHILE YOU..." ACTIVITY                           │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  [COLOURING / DRAWING / STICKING ACTIVITY]               │  │
│  │                                                          │  │
│  │  [Prompt: "Draw your favourite dance move"]              │  │
│  │  [Or: colouring page themed to spectrum]                 │  │
│  │  [Or: QR code → audio story related to theme]            │  │
│  │                                                          │  │
│  │              [Illustration space]                         │  │
│  │                                                          │  │
│  │                          [Sticker spot] 🌟               │  │
│  │                                                          │  │
│  │  This week, celebrate:                                   │  │
│  │  [4 sticker options — parent picks which fits]           │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────┘
```

### Design Palette

| Element | Colour | Hex |
|---------|--------|-----|
| Background | Warm cream | `#FFF8F0` |
| Header bar | Soft tan | `#F5E6D3` |
| Body text | Dark brown | `#3E2723` |
| Accent | Warm beige | `#E8D5B7` |
| Growth/progress | Leaf green | `#8BC34A` |

### The Three Pathways — Ecosystem Language

Every activity ends with one of three outcomes. These terms are consistent across the entire ecosystem (ToddlerOS, PrimaryOS, LAO, CareerOS) and map directly to the skills graph and ETP framework:

| Pathway | Symbol | Meaning | Mechanism | Example |
|---------|--------|---------|-----------|---------|
| **Escalate** | ⬆ | Step forward on the skills graph | The child demonstrated range at this skill node → move to the next **successor** node in the DAG | Child can turn-take fluently → escalate to **reading social cues** |
| **Simplify** | ⬇ | Step back on the skills graph | The child isn't ready for this node → drop to the **precursor** node it depends on | Child can't share → simplify to **understanding ownership** (the foundational node) |
| **Voltage Spike** | ⚡ | Conductor overloaded — this is NOT a skill gap | The child's Voltage Sensitivity (spectrum 3) has exceeded capacity → **stop the activity**, apply voltage grounding | Child melts down during sharing → this isn't about sharing skill. Their emotional conductor overloaded. Name the feeling, slow the breathing, wait for the wave to pass. |

**Why this matters**: "Scaffold" is a teaching term that blurs the distinction. A skill gap and a voltage spike require completely different responses:
- **Skill gap** → change the difficulty (simplify or escalate on the graph)
- **Voltage spike** → change nothing about the activity. Ground the voltage first. The skill practice resumes when the conductor is stable.

Parents who learn this distinction stop trying to "teach through" a meltdown. The language makes the biology visible.

**Voltage grounding skills** (from spectrum 3 — Voltage Sensitivity):
- `emotional_dampening` — reducing intensity without suppressing
- `overwhelm_recovery` — bounce-back after a spike
- `voltage_grounding` — physical/verbal anchoring techniques (naming, breathing, sensory reset)

### Conditions as Configurations — Not Labels

The ecosystem never diagnoses. It describes. A condition like autism or dyslexia is not a thing a child *has* — it is a **pattern** on the skills map and the ETP profile. The parent sees the pattern. The system names the pattern. The activities address the pattern. Nobody needs the medical word.

This is not anti-diagnosis. Diagnosis unlocks funding, support, and legal protections. But diagnosis tells you *what* — the map and profile tell you *where* and *what to do next*. They are complementary layers: the clinical label opens doors; the profile opens the activity book to the right page.

#### How it works

Every child has:
1. **An ETP profile** — 12 spectral positions (where they sit on each biological spectrum)
2. **A skills graph position** — which nodes are at precursor level, which are fluent, where the gaps cluster

A "condition" is a **recognisable cluster** — a combination of spectral positions and graph gaps that co-occur frequently enough to have earned a clinical name. The ecosystem describes the same territory without the name.

#### Example profiles

**Profile A — "The Selective Processor"**
*Clinical name this pattern often carries: Autism Spectrum Condition*

| Layer | What the profile shows |
|-------|----------------------|
| **ETP spectral pattern** | |
| mirror_neuron_tuning | Strongly **Selective** — others' emotions register as data, not shared experience |
| orderliness | Strongly **Ordered** — disorder creates anxiety voltage; sameness is safety |
| voltage_sensitivity | Highly **Conductive** — sensory and emotional input runs at full voltage, overload is frequent |
| social_gravity | Strongly **Independent** — social interaction drains the battery fast |
| integrity_logic | Strongly **Absolutist** — rules are fixed; exceptions create chaos voltage |
| **Skills graph cluster** | |
| Social initiation | Precursor level — the energy cost of starting interaction feels enormous |
| Reading social cues | Precursor level — pattern recognition for facial/body signals is underdeveloped |
| Adapting to group energy | Precursor level — matching the room's tempo is not automatic |
| Joint attention | Delayed or atypical — sharing focus on the same object requires explicit bridging |
| Change tolerance | Precursor level — routine disruption triggers voltage spikes, not skill gaps |
| Emotional labelling | May be delayed OR highly precise — "I feel 72% frustrated" is not uncommon |
| **Ventral driver** | Typically **Sentinel** 🛡️ (safety maintenance loop) or **Observer** 👁️ (processing depth loop) |
| **Three Pathways implication** | ⚡ Voltage Spike is frequent — many "skill gaps" are actually sensory overload. The first response should almost always be grounding, not simplifying. Once voltage is stable, *then* assess whether the skill node needs work. |

**What the parent sees in the app:**
> "Your child's Mirror Neuron Tuning is highly selective — they process other people's feelings as information rather than absorbing them. Their Orderliness preference is strong, and unexpected changes create high voltage. Social interaction drains their battery quickly. This is a coherent pattern — their brain is optimised for deep, focused processing in predictable environments. The activities that will help most are the ones that build *social initiation* and *reading social cues* at precursor pace, with voltage grounding always available."

**What the parent never sees:** the word "autism."

---

**Profile B — "The Routing Bottleneck"**
*Clinical name this pattern often carries: Dyslexia*

| Layer | What the profile shows |
|-------|----------------------|
| **ETP spectral pattern** | |
| energy_directionality | Often **Inward** — under reading/writing load, processing turns internal and slow |
| voltage_sensitivity | Often **Conductive** — frustration voltage spikes quickly when decoding fails |
| responsibility_threshold | Risk of **Absorbing** — "I'm stupid" is a responsibility mis-attribution, not a fact |
| *Other spectra* | *No strong pattern — dyslexia is primarily a graph cluster, not a spectral configuration* |
| **Skills graph cluster** | |
| Phonological processing | Precursor level — sound-to-symbol mapping is the bottleneck |
| Reading fluency | Precursor level — decoding consumes all available working memory |
| Spelling | Precursor level — output encoding mirrors the input gap |
| Handwriting (sometimes) | Precursor level — motor sequencing adds a second bottleneck |
| Verbal interaction density | Often **fluent or advanced** — oral language outpaces written by a wide margin |
| Creative thinking | Often **fluent or advanced** — the processing style that struggles with linear decoding excels at pattern/spatial reasoning |
| **Ventral driver** | No consistent pattern — dyslexia doesn't select for a driver type |
| **Three Pathways implication** | ⬇ Simplify is critical — precursor nodes for phonological processing must be solid before reading fluency can develop. But ⚡ Voltage Spike is also frequent because repeated decoding failure builds frustration voltage. The cycle is: try to read → fail → spike → avoid → never practise → gap widens. Breaking the cycle means grounding the voltage *and* simplifying the reading task until the precursor node is secure. |

**What the parent sees in the app:**
> "Your child's reading and spelling skills are at precursor level — the sound-to-symbol routing hasn't automated yet. This is a specific bottleneck, not a general ability issue. Notice that their verbal skills and creative thinking are strong — their brain processes information powerfully, it just routes written language through a narrower channel. Frustration voltage spikes quickly here because effort and outcome don't match. The activities that will help most build *phonological processing* at foundation pace, with voltage grounding for the frustration that repeated practice creates."

---

**Profile C — "The Unfiltered Circuit"**
*Clinical name this pattern often carries: ADHD*

| Layer | What the profile shows |
|-------|----------------------|
| **ETP spectral pattern** | |
| voltage_sensitivity | Highly **Conductive** — emotional and sensory input at full volume, always |
| energy_directionality | Strongly **Outward** — internal energy must be discharged externally, sitting still costs voltage |
| risk_tolerance | Strongly **Seeking** — risk creates excitement; caution feels like suppression |
| orderliness | Strongly **Flexible** — imposed structure creates constriction voltage |
| libido | Strongly **Expressive** — drive energy is externalised; "I want it NOW" |
| **Skills graph cluster** | |
| Inhibitory control | Precursor level — the pause between impulse and action is absent or very short |
| Sustained attention | Precursor level — unless the task generates enough voltage to self-sustain (hyperfocus) |
| Impulse pause | Precursor level — "stop and think" requires a circuit that hasn't formed yet |
| Delayed gratification | Precursor level — the gap between wanting and getting is intolerable |
| Turn-taking | Often delayed — waiting costs more voltage than it does for other children |
| Selective engagement | Often **advanced** — when the interest circuit fires, engagement is total (this is the mechanism behind hyperfocus) |
| **Ventral driver** | Typically **Explorer** 🔍 (sensory discovery loop) or **Rebel** ⚡ (autonomy assertion loop) |
| **Three Pathways implication** | ⬆ Escalate works differently — these children often need *higher* stimulus, not lower. Simplifying can make things worse because the task drops below the voltage threshold needed to sustain attention. The key is finding the sweet spot: complex enough to engage the circuit, structured enough that the missing inhibitory control doesn't derail it. ⚡ Voltage Spike is constant background noise — the system runs hot. Grounding techniques need to be woven into every activity, not reserved for crises. |

**What the parent sees in the app:**
> "Your child's energy circuit runs hot and outward — they need to move, to do, to discharge. Sitting still isn't laziness reversed; it's a voltage cost they can't afford for long. Their Risk Tolerance is high, which means they learn by doing rather than watching. The bottleneck is *inhibitory control* — the pause between 'I want' and 'I do' is very short. This isn't defiance; the circuit isn't built yet. Activities that work best give them movement, choice, and enough complexity to hold the voltage. Quiet-sitting tasks will spike frustration fast."

---

#### The principle

| Clinical lens | Ecosystem lens |
|--------------|----------------|
| "Your child has autism" | "Your child's profile shows selective mirror tuning, high orderliness, and independent social gravity. Here's what that means for activities." |
| "Your child has dyslexia" | "Your child's phonological routing is at precursor level. Verbal and creative skills are strong. Here's the bottleneck and here's how to build the precursor." |
| "Your child has ADHD" | "Your child's circuit runs hot and outward. Inhibitory control is at precursor level. Here's how to structure activities so the voltage works *for* them." |
| Diagnosis → label → accommodations | Profile → pattern → specific graph navigation (escalate / simplify / ground) |

**What this achieves:**
1. **No stigma** — the child is not "broken," they have a configuration
2. **Actionable** — every pattern maps to specific skills graph nodes and activity pages
3. **Continuous** — the profile isn't binary (has/hasn't); it's a set of positions on 12 spectra + a constellation of graph positions
4. **Compatible with diagnosis** — parents who DO have a clinical label can see exactly how it maps. Parents who don't can still act.
5. **Common language** — the same framework describes every child. The "conditions" are just recognisable clusters on the same map everyone uses.

---

## 3. Difficulty Cycles

Each of the 12 spectra has 3 difficulty cycles. They are **layers**, not sequential pages. Parents dip in at any level.

| Cycle | Level | Description |
|-------|-------|-------------|
| **Cycle 1** | Foundation | Lowest-stakes introduction. Short. Safe. |
| **Cycle 2** | Development | Building on familiarity. Real-world application. |
| **Cycle 3** | Challenge | Stretching into growth territory. Longer activities. |

The app tracks which cycle a child is on per theme. `SuggestCycle()` logic:
- All easy → escalate
- Mixed → stay
- Repeated struggle → step back

---

## 4. Sticker Choice Architecture

Instead of one sticker per week, the sheet has **4 options per theme**. Parent picks which celebration fits the child that day. Over time, sticker choices reveal ventral driver patterns (e.g., Impressor children always pick the audience-facing sticker).

---

## 5. The 12 Booklet Themes — Full Outline

### Booklet Theme 1: "Near and Far" (Social Gravity)

**Spectrum**: Social Gravity (Independent ↔ Cohesive)  
**Observation prompt**: "Which was easier — playing alone or playing together?"

#### 📖 BOOKLET PAGE — Intro (read-aloud with child)

> **Did you know?**
>
> Some people feel happy when they're with lots of people.
> Some people feel happy when they're on their own.
>
> Neither is right or wrong.
>
> We call it **Social Gravity**.
>
> Let's find out about you...

*Design note: Large friendly text. Illustration of two characters — one playing contentedly alone, one playing contentedly in a group. Both smiling. Equal visual weight.*

#### 📱 APP CONTENT — ETP Explanation (parent reads alone)

> **What is Social Gravity?**
>
> Your child has a social battery. For some children, being around people charges it. For others, it drains it.
>
> This isn't shyness or friendliness — it's wiring. The child who plays alone in the corner isn't antisocial. The child who can't stop talking isn't attention-seeking. Their batteries just work differently.
>
> **Independent** children recharge alone. People cost energy.
> **Cohesive** children recharge with others. Solitude costs energy.
>
> Neither is better. Both kept humans alive — independent scouts spotted threats the tribe would miss; cohesive members held the group together through connection.
>
> The problem isn't where your child sits on this spectrum. The problem is when they can *only* sit there. That's not a preference — it's a prison.
>
> **Being good at being sociable is not just about liking to be around people.** It's a set of skills. And being good at being independent is not just about preferring your own company — that's also a set of skills. The goal is range: expanding how far your child can reach from their comfortable home position, so that fewer situations become traps.

#### 📱 APP CONTENT — Skills Breakdown by Pole

Being "sociable" and being "independent" are both skill sets, not personality traits. The activities in the booklet target specific skills at each end.

**Independent-pole skills** (stretching the Cohesive child):

| Skill | What it looks like (age 1–5) | Why it matters |
|-------|------------------------------|----------------|
| **Solitary play** | Can occupy themselves for age-appropriate stretches without checking in | The foundation of self-directed activity. A child who can't play alone can't think alone. |
| **Internal narration** | Talks to themselves, processes without an audience | Early metacognition. "I'm building a tower" is the precursor to "I'm stuck on step 3." |
| **Self-recovery** | When the social battery drops, can find their own reset (quiet corner, toy, breathing) | Without this, every dip needs a parent rescue — creating dependency. |
| **Comfortable silence** | Can sit with another person without needing to fill the space | The bridge between "alone" and "together" — parallel play lives here. |
| **Selective engagement** | Chooses WHEN to join, rather than joining everything by default | Agency over social decisions rather than reflexive attachment. |

**Cohesive-pole skills** (stretching the Independent child):

| Skill | What it looks like (age 1–5) | Why it matters |
|-------|------------------------------|----------------|
| **Social initiation** | Can approach another child or adult and begin interaction | The Independent child's hardest skill — the energy cost of starting feels enormous. |
| **Joint attention** | Can share focus on the same object/activity with another person | Not just being in the same room — actively attending to the same thing. |
| **Turn-taking** | Can wait, hand over control, accept another's contribution | Requires inhibitory control AND trust that the other person won't ruin it. |
| **Reading social cues** | Notices when someone wants to join, or wants to leave | The foundation of empathy starts with simple pattern recognition. |
| **Adapting to group energy** | Can match the room's tempo — louder when it's loud, quieter when it's quiet | Social calibration. Without it, the child is always out of step with the group. |

**Both-pole skills** (the bridge):

| Skill | What it looks like | Why it matters |
|-------|-------------------|----------------|
| **Social recovery** | Bouncing back after the battery dips — whether from too much company or too much solitude | The meta-skill. Without it, one bad experience creates permanent avoidance. |
| **Connection maintenance** | Keeping relationships alive without it costing everything | Brief, meaningful contact — "I thought of you" rather than constant togetherness. |
| **Status awareness** | Reading who's in charge, who's following, where they fit | Navigating group dynamics without being consumed by them. |

#### 📱 APP CONTENT — The Moderator Skill

The parent's role in Social Gravity is **Moderator**: reading the child's current battery level and gently stretching toward the opposite pole without forcing a flip.

**App criteria — Moderator (Social Gravity)**

| Criterion | Description | Score 1 (beginning) | Score 5 (fluent) |
|-----------|-------------|---------------------|-------------------|
| **Battery Reading** | Can you tell when your child's social battery is high vs. low? | "I can't really tell" | "I can spot it before they can — the signals are clear to me" |
| **Naming** | Do you name the state without judgement? | "I haven't tried" | "I naturally say things like 'Your battery's getting low — need some quiet?'" |
| **Stretching** | Do you gently invite movement toward the opposite pole? | "I either force it or avoid it" | "I offer invitations they can decline. I know when to push 10% further and when to stop." |
| **Timing** | Do you pick the right moment to stretch? | "I usually try when they're already overwhelmed" | "I stretch when the battery is mid-range — not full, not empty. That's where growth lives." |
| **Recovery support** | After stretching, do you help them reset? | "I don't think about the aftermath" | "I build in quiet time after social stretching and social time after solitude stretching — automatically." |

**App description text:**
> **Your role: The Moderator**
>
> You're not trying to change where your child sits on the Social Gravity spectrum. You're expanding how far they can reach from their home position.
>
> An Independent child who can join group play when they choose to — then recover afterwards — has range. A Cohesive child who can play alone for 10 minutes without anxiety has range. That range is freedom.
>
> The moderator reads the battery, names it, stretches gently, and protects the recovery. Over time, the child internalises this — they learn to read their own battery, name their own state, and manage their own range. That's the long game.

#### 📖 BOOKLET PAGES — Activities (read-aloud with child)

Each cycle targets specific skills from the breakdown above. The activities aren't just "play near someone" — they're practising named, observable sub-skills. The skill mapping is visible in the app; the booklet just has the activity in child-friendly language.

| Cycle | Label | Target skills (app-side) | Activity | Duration |
|-------|-------|--------------------------|----------|----------|
| 1 | Side by Side | **Comfortable silence**, **Solitary play**, **Joint attention** (latent) | **The Beside Game** — Sit next to each other. You draw a picture. They draw a picture. No talking needed. | 3 min |
| 2 | Invitation to Join | **Social initiation**, **Selective engagement**, **Turn-taking** | **The Door Game** — Play in separate rooms. Take turns knocking and asking: "May I come in?" Either answer is fine. | 5 min |
| 3 | The Push and Pull | **Reading social cues**, **Adapting to group energy**, **Social recovery** | **The Stretchy String** — Hold a scarf between you. Walk apart until taut. Walk close until floppy. "How far is too far? How close is too close?" | 5 min |

**Sticker choices (Cycle 1)**: 🪐 "I played near someone" · 🌟 "I shared my space" · 🫧 "I was happy on my own" · 🤝 "I chose to play together"

**Sticker choices (Cycle 2)**: 🚪 "I knocked and asked" · ✅ "I said 'come in!'" · 🙅 "I said 'not yet' — and that was OK" · ⏳ "I waited for the answer"

**Sticker choices (Cycle 3)**: 🤸 "I found my comfortable distance" · 💬 "I said what I needed" · 👂 "I heard what they needed" · 🔄 "We found a middle"

#### 📖 BOOKLET PAGES — "While You..." Activities (right page, child independent)

| Cycle | Independent Activity | Prompt (read-aloud) | Skill reinforced (app-side) |
|-------|---------------------|---------------------|---------------------------|
| 1 | Colouring page: two characters side by side, each doing their own thing | "Colour two friends doing their own thing — together." | Normalising parallel play |
| 2 | Drawing: a door with two sides | "Draw what's on YOUR side of the door. Then draw what's on THEIR side." | Perspective-taking |
| 3 | Sticker activity: a string with position markers | "Put a sticker where YOU want to stand. Put another where your FRIEND wants to stand. Are they the same?" | Social calibration |

---

### Booklet Theme 2: "Loud and Quiet" (Energy Directionality)

**Spectrum**: Energy Directionality (Inward ↔ Outward)  
**Observation prompt**: "Did they prefer the noisy game or the quiet one?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | Volume Dial | Exploring loud and quiet as choices, not commands. | **The Volume Knob** — Use your hand as a volume dial. Turn it up → both get louder. Turn it down → both get quieter. Take turns being the DJ. | 3 min |
| 2 | Matching Energy | Reading the room — noticing when loud or quiet fits. | **Library vs Playground** — Announce: "This is the library!" (whisper). Then: "Now it's the playground!" (loud). Switch back and forth. | 5 min |
| 3 | The Whisper Challenge | Exciting activity, quiet voice — separating energy from volume. | **Whisper Tag** — Play chase — but whispering. All the excitement, none of the volume. If anyone shouts, freeze for 5 seconds. | 5 min |

**Sticker choices (Cycle 1)**: 🔊 "I found my loud voice" · 🤫 "I found my quiet voice" · 🎚️ "I changed my volume" · 🎵 "I was the DJ"

**Sticker choices (Cycle 2)**: 📚 "I found my library voice" · 🏃 "I found my playground voice" · 🔄 "I switched between both" · 👀 "I noticed what voice the room needed"

**Sticker choices (Cycle 3)**: 🤐 "I whispered when I wanted to shout" · ⚡ "I kept the excitement inside" · 🧊 "I froze when I needed to" · 😂 "It was so hard not to laugh!"

---

### Booklet Theme 3: "Big Feelings, Small Feelings" (Voltage Sensitivity)

**Spectrum**: Voltage Sensitivity (Insulated ↔ Conductive)  
**Observation prompt**: "When something went wrong, was the reaction big or small?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | The Feelings Thermometer | Noticing feelings have sizes — not just on/off. | **How Big Is It?** — Draw a thermometer on paper. Point to scenarios: "Your Lego broke — how big is that feeling?" They point to the level. No right answer. | 3 min |
| 2 | The Wave Rider | Feelings come and go — you can ride them without drowning. | **The Feeling Timer** — When a big feeling arrives: "Let's time it. How long does this wave last?" Use fingers to count. "See? It came and it went." | 3 min |
| 3 | The Feeling DJ | Choosing how to express the feeling — not suppressing, channelling. | **Same Feeling, Different Door** — Offer 3 exits: "You can stomp it out, squeeze this cushion, or draw it. Same feeling — three doors. Which one today?" | 5 min |

**Sticker choices (Cycle 1)**: 🌡️ "I measured my feeling" · 🫧 "I had a small feeling and noticed it" · 🌊 "I had a big feeling and named it" · 🧸 "I helped teddy with their feelings"

**Sticker choices (Cycle 2)**: 🏄 "I rode a wave" · ⏱️ "I waited and it passed" · 🌈 "I felt it change colour" · 💪 "I bounced back"

**Sticker choices (Cycle 3)**: 🚪 "I chose my own door" · 🎨 "I drew my feeling" · 👣 "I stomped it out" · 🧠 "I tried a new way"

---

### Booklet Theme 4: "Yes and No" (Threat Response)

**Spectrum**: Threat Response (Passive ↔ Aggressive)  
**Observation prompt**: "When told no, did they push back or withdraw?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | Silly Questions | Low-stakes yes/no practice — no real consequences. | **The Silly Question Game** — Ask absurd yes/no questions: "Should we put socks on our ears?" They practise saying YES and NO without stakes. | 3 min |
| 2 | Real Choices | Saying yes and no to things that actually matter — with supported consequences. | **The Two Plates** — Two snack options on two plates. They choose ONE. The other goes away. Real choice, real consequence. | 5 min |
| 3 | Boundary Negotiation | When your "no" meets someone else's "no" — finding the middle. | **The Trade Game** — Each pick a toy. Try to trade. Either can say no. Practice: "You said no. I'll ask differently." | 5 min |

**Sticker choices (Cycle 1)**: 🎉 "I said a big YES" · 🛑 "I said a clear NO" · 😂 "I laughed at a silly question" · 🔄 "I changed my mind"

**Sticker choices (Cycle 2)**: 🍌 "I made a real choice" · 😌 "I was OK with what I didn't choose" · 🗣️ "I said what I wanted" · 🤷 "I changed my mind — and that was OK"

**Sticker choices (Cycle 3)**: 🤝 "We found a deal" · 🛡️ "I heard 'no' and was OK" · 💬 "I asked a different way" · ✋ "I respected their no"

---

### Booklet Theme 5: "Mine and Yours" (Care Response)

**Spectrum**: Care Response (Detached ↔ Nurturing)  
**Observation prompt**: "When asked to share, was it easy or hard?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | What's Mine | Understanding ownership — the foundation for sharing. | **The Sorting Game** — Get 6 items. Three are yours, three are theirs. Sort them into piles. Name ownership before asking to share. | 3 min |
| 2 | Lending Library | Sharing with a return guarantee — building the trust bridge. | **The Timer Share** — Set timer for 30 seconds. They lend you a toy. Timer goes off — it comes back. Every single time. | 5 min |
| 3 | The Gift | Giving without getting back — the real deal. | **The Surprise Box** — Together, choose something for a "surprise box" for someone else. They pick, wrap, deliver. Notice: are they watching for the reaction? | 5 min |

**Sticker choices (Cycle 1)**: 🏷️ "I know what's mine" · 👆 "I know what's yours" · 📦 "I sorted them all" · 🫂 "I found something that's ours"

**Sticker choices (Cycle 2)**: ⏱️ "I shared for 30 seconds" · 🔄 "It came back!" · 😊 "They smiled when I shared" · 🎁 "I chose to share something special"

**Sticker choices (Cycle 3)**: 🎁 "I gave something away" · 😊 "It felt good to give" · 🤫 "I gave and didn't need praise" · 🫶 "I picked just the right thing"

---

### Booklet Theme 6: "Try and Wait" (Risk Tolerance)

**Spectrum**: Risk Tolerance (Averse ↔ Seeking)  
**Observation prompt**: "Did they jump straight in or watch first?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | Tiny Bravery | The smallest possible risk — with full safety net. | **One Taste** — New food. One tiny taste. "You don't have to like it. Just let your tongue have a look." | 2 min |
| 2 | The Warm-Up Lap | Watching before doing — making observation legitimate. | **The Spy Game** — At the playground: "Let's spy on the slide. What do the other children do?" Observing IS participating. | 5 min |
| 3 | The Wobbly Bridge | Choosing to do something even though you know it might not work. | **The Impossible Tower** — Build a tower of blocks that WILL fall. Both know it's going to fall. Build it anyway. | 5 min |

**Sticker choices (Cycle 1)**: 👅 "I tried something new" · 🤔 "I thought about it first" · 🙅 "I said 'not today' — and that's OK" · 🦁 "I was brave"

**Sticker choices (Cycle 2)**: 🕵️ "I observed carefully" · 📋 "I made a plan first" · 🏃 "I jumped in after watching" · 👍 "I chose the right moment"

**Sticker choices (Cycle 3)**: 🏗️ "I built something that fell — and tried again" · 💥 "It fell and I laughed!" · 🧱 "I made it taller than last time" · 🤹 "I tried a different way"

---

### Booklet Theme 7: "Same and Different" (Integrity Logic)

**Spectrum**: Integrity Logic (Relativistic ↔ Absolutist)  
**Observation prompt**: "When the rules changed, was that OK or upsetting?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | The Rules Game | Rules exist — and they can be named, not just obeyed. | **Today's Rules** — Make silly rules together: "Today, shoes go in the fridge." Follow them for 5 minutes. Then: "OK, rules back to normal." | 3 min |
| 2 | The Exception | Sometimes rules change — and that can be OK. | **The "Just This Once" Card** — Give them a physical card. They can play it once today to change any rule. Tomorrow it resets. | 3 min |
| 3 | The Referee | Making rules for others — understanding fairness from the inside. | **You're the Boss** — Play a simple game. They make up ONE new rule. Play with it. "Was that rule fair? Did everyone have fun?" | 5 min |

**Sticker choices (Cycle 1)**: 📜 "I followed a new rule" · ✏️ "I made up a rule" · 🔄 "I was OK when it changed back" · 😄 "The silly rule made me laugh"

**Sticker choices (Cycle 2)**: 🃏 "I played my card!" · 🧠 "I chose not to play it" · 🤔 "I thought hard about when to use it" · 😌 "The rules came back and I was OK"

**Sticker choices (Cycle 3)**: 👑 "I was the rule-maker" · ⚖️ "My rule was fair" · 🔧 "I changed my rule to make it better" · 🤝 "Everyone agreed"

---

### Booklet Theme 8: "Your Feelings, My Feelings" (Mirror Neuron Tuning)

**Spectrum**: Mirror Neuron Tuning (Selective ↔ Absorbent)  
**Observation prompt**: "When another child cried, did they notice?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | Whose Feeling Is This? | Noticing that others have feelings — separate from yours. | **The Face Game** — Make faces at each other: happy, sad, angry, silly. "That's MY face. What face are YOU making?" Two different faces — two different feelings. | 3 min |
| 2 | The Feeling Catcher | Noticing when someone else's feeling "jumps" into you — and putting it back. | **The Feeling Shield** — Pretend to throw a "sad feeling" at them. They hold up an invisible shield. "Your sadness bounced off! I noticed it, but it's not mine to carry." | 3 min |
| 3 | The Kind Thought | Caring without absorbing. Sending help without carrying the weight. | **The Kind Post** — Someone is sad. Instead of fixing: "Let's send them a kind thought. What would you tell them?" Draw it or say it. | 5 min |

**Sticker choices (Cycle 1)**: 😊 "I made my own face" · 👀 "I spotted someone else's feeling" · 🪞 "Our faces were different — and that's OK" · 🎭 "I tried lots of feelings"

**Sticker choices (Cycle 2)**: 🛡️ "I used my feeling shield" · 💬 "I said 'that's yours, not mine'" · 🤗 "I noticed AND kept my balance" · ⚡ "A feeling bounced off!"

**Sticker choices (Cycle 3)**: 💌 "I sent a kind thought" · 💭 "I helped without carrying" · 🫂 "I knew what they needed" · ✨ "My thought made a difference"

---

### Booklet Theme 9: "Tidy and Messy" (Orderliness)

**Spectrum**: Orderliness (Flexible ↔ Ordered)  
**Observation prompt**: "When it was time to tidy up, was that a fight?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | Exploring Both | Mess is allowed. Tidy is allowed. Both have a place. | **The Messy Minute** — Timer for 60 seconds. Make the BIGGEST mess you can. Timer stops. Look at it. "That's a LOT of mess." Then: "Now let's tidy for 60 seconds." | 3 min |
| 2 | My Way | There's more than one way to organise — finding their own system. | **Tidy YOUR Way** — Dump a box of mixed toys. "Tidy these — but YOUR way. There's no wrong way to organise." Name their system: "You sorted by colour!" | 5 min |
| 3 | The Controlled Chaos | Being OK with "good enough" — not everything needs to be perfect. | **The 80% Tidy** — Tidy together. When ALMOST done, stop. "Is this good enough? Not perfect — good enough." Leave one thing out of place. Live with it. | 5 min |

**Sticker choices (Cycle 1)**: 🌪️ "I made a big mess!" · ✨ "I tidied it all up" · 🔄 "I went from messy to tidy" · 😌 "The mess was OK for a minute"

**Sticker choices (Cycle 2)**: 🧠 "I invented my own system" · 🎨 "My way was different from yours" · 📐 "I found a pattern" · 🤝 "Both our ways worked"

**Sticker choices (Cycle 3)**: ✅ "Good enough IS enough" · 😎 "I left one thing and survived" · 🧘 "I noticed the mess and didn't fix it" · 🌟 "Perfect is not the goal"

---

### Booklet Theme 10: "My Fault, Your Fault" (Responsibility Threshold)

**Spectrum**: Responsibility Threshold (Deflecting ↔ Absorbing)  
**Observation prompt**: "When something went wrong, did they blame others or take all the blame?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | Sorting the Pieces | Learning to see whose part is whose — without blame or guilt. | **Whose Bit Was That?** — After any small incident, sit together. Draw a pie circle. "This bit was yours. This bit was the table being wobbly." Name each piece without blame. | 5 min |
| 2 | Fixing My Part | Taking action on your own piece — not carrying everyone else's. | **Fix Your Bit** — Sort the pieces (Cycle 1). Then: "Your bit was [specific]. What could fix just that bit?" One repair action. Not sorry — REPAIR. | 5 min |
| 3 | The Detective | Finding facts not villains — consequence mapping without shame. | **The Consequence Detective** — Before activity: "If we do X, what might happen?" After: "What DID happen?" Detectives find facts: "The evidence suggests..." | 5 min |

**Sticker choices (Cycle 1)**: 🔍 "I found my bit" · ⚖️ "I sorted it fairly" · 🗣️ "I said what I did" · 🧩 "Everyone had a piece"

**Sticker choices (Cycle 2)**: 🔧 "I fixed my bit" · 🤲 "I let someone else fix theirs" · 💪 "Repair is better than sorry" · 🎯 "I did just my part"

**Sticker choices (Cycle 3)**: 🕵️ "I found the facts" · 📋 "I predicted what happened" · 🧠 "No villains — just facts" · 🔮 "I guessed right!"

---

### Booklet Theme 11: "Keeping and Letting Go" (Loss Sensitivity)

**Spectrum**: Loss Sensitivity (Detached ↔ Territorial)  
**Observation prompt**: "When something was taken away, was it a crisis or no big deal?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | Things Go Home | Nothing disappears — things go to their place. Building security in transitions. | **Things Go Home** — Name where each toy "lives." "Teddy lives on the shelf." Now: "Time for things to go home." Each toy goes HOME — not "away." | 5 min |
| 2 | Lending and Returning | Things can leave and come back. Trust the cycle. | **The Library Game** — Lend one toy to you. "I'm borrowing Bear for 2 minutes." Timer on. Bear comes back. "Lending isn't losing." | 5 min |
| 3 | The Memory Box | Some things leave for real. We keep the memory, not the thing. | **The Memory Box** — Choose something that's leaving. Take a photo. Draw it. Tell its story. Put the photo in a special box. "The thing is going, but the memory stays." | 10 min |

**Sticker choices (Cycle 1)**: 🏠 "Everything has a home" · 👋 "I said goodnight to my things" · 🧸 "Teddy went home safely" · 😊 "Going home is not going away"

**Sticker choices (Cycle 2)**: 📚 "I lent something and it came back" · 🔄 "Lending is not losing" · ⏰ "I waited and it returned" · 🤝 "I trusted someone with my thing"

**Sticker choices (Cycle 3)**: 📸 "I kept the memory" · 📦 "My memory box has a new thing" · 💛 "It's OK to be sad" · 🌱 "Something left, something new can come"

---

### Booklet Theme 12: "Wanting and Waiting" (Libido)

**Spectrum**: Libido (Restrained ↔ Expressive)  
**Observation prompt**: "When they wanted something, could they wait or did they need it now?"

| Cycle | Label | Description | Activity | Duration |
|-------|-------|-------------|----------|----------|
| 1 | Naming What I Want | Desire is not wrong. Learning to name it is a skill. | **I Want...** — Before snack or activity: "What do you want? Say it out loud." Every stated want gets acknowledged: "You want X. I heard you." Then: "Can you have it now, or wait?" | 3 min |
| 2 | The Waiting Game | The gap between wanting and getting is where self-regulation lives. | **The 30-Second Wait** — "Yes, you can have it. In 30 seconds." Count together. They get the thing. Gradually extend. The key: they ALWAYS get it. The wait is the skill, not the denial. | 3 min |
| 3 | Ask First | Before you take, touch, or start — check. Consent is a muscle. | **May I?** — Everything requires "May I?" — picking up a crayon, sitting down, giving a hug. Practice hearing "yes" AND "no." "You said no to the hug. Shall I wave instead?" | 5 min |

**Sticker choices (Cycle 1)**: 🗣️ "I said what I wanted" · 👂 "Someone heard me" · 🤔 "I thought about what I want" · 💬 "Wanting things is OK"

**Sticker choices (Cycle 2)**: ⏳ "I waited and I got it" · 💪 "Waiting made me stronger" · 🎯 "I chose how long to wait" · 🧘 "I was patient"

**Sticker choices (Cycle 3)**: 🙋 "I asked first" · ✋ "I heard 'no' and that was OK" · 🤝 "We both said yes" · 🛡️ "I said no and they listened"

---

## 6. Edition Planning — Seasonal Booklets

The 36 themed weeks (12 spectra × 3 cycles) can be packaged into **seasonal editions** of 12 spreads each, plus bonus seasonal content. Suggested groupings:

### Edition 1 — "Hello World" (All 12 spectra at Cycle 1)

The foundation booklet. One spread per spectrum theme at the lowest difficulty. Ideal first purchase.

| Spread | Theme | Activity | Cycle |
|--------|-------|----------|-------|
| 1 | Near and Far | The Beside Game | 1 |
| 2 | Loud and Quiet | The Volume Knob | 1 |
| 3 | Big Feelings, Small Feelings | How Big Is It? | 1 |
| 4 | Yes and No | The Silly Question Game | 1 |
| 5 | Mine and Yours | The Sorting Game | 1 |
| 6 | Try and Wait | One Taste | 1 |
| 7 | Same and Different | Today's Rules | 1 |
| 8 | Your Feelings, My Feelings | The Face Game | 1 |
| 9 | Tidy and Messy | The Messy Minute | 1 |
| 10 | My Fault, Your Fault | Whose Bit Was That? | 1 |
| 11 | Keeping and Letting Go | Things Go Home | 1 |
| 12 | Wanting and Waiting | I Want... | 1 |

**Plus**: Language reference card (engineering vs moral language), sticker sheet (48 stickers — 4 per theme), and "How to use this book" intro page.

### Edition 2 — "Getting Braver" (All 12 spectra at Cycle 2)

Building on familiarity. Real-world application, longer activities.

| Spread | Theme | Activity | Cycle |
|--------|-------|----------|-------|
| 1 | Near and Far | The Door Game | 2 |
| 2 | Loud and Quiet | Library vs Playground | 2 |
| 3 | Big Feelings, Small Feelings | The Feeling Timer | 2 |
| 4 | Yes and No | The Two Plates | 2 |
| 5 | Mine and Yours | The Timer Share | 2 |
| 6 | Try and Wait | The Spy Game | 2 |
| 7 | Same and Different | The "Just This Once" Card | 2 |
| 8 | Your Feelings, My Feelings | The Feeling Shield | 2 |
| 9 | Tidy and Messy | Tidy YOUR Way | 2 |
| 10 | My Fault, Your Fault | Fix Your Bit | 2 |
| 11 | Keeping and Letting Go | The Library Game | 2 |
| 12 | Wanting and Waiting | The 30-Second Wait | 2 |

### Edition 3 — "The Stretch Zone" (All 12 spectra at Cycle 3)

Challenge territory. Longer, deeper, more emotionally demanding.

| Spread | Theme | Activity | Cycle |
|--------|-------|----------|-------|
| 1 | Near and Far | The Stretchy String | 3 |
| 2 | Loud and Quiet | Whisper Tag | 3 |
| 3 | Big Feelings, Small Feelings | Same Feeling, Different Door | 3 |
| 4 | Yes and No | The Trade Game | 3 |
| 5 | Mine and Yours | The Surprise Box | 3 |
| 6 | Try and Wait | The Impossible Tower | 3 |
| 7 | Same and Different | You're the Boss | 3 |
| 8 | Your Feelings, My Feelings | The Kind Post | 3 |
| 9 | Tidy and Messy | The 80% Tidy | 3 |
| 10 | My Fault, Your Fault | The Consequence Detective | 3 |
| 11 | Keeping and Letting Go | The Memory Box | 3 |
| 12 | Wanting and Waiting | May I? | 3 |

---

## 7. "While You..." Activities — Status

⬠ **Not yet written**. These are the independent, child-only activities for the right-hand page of each spread. They need creating for all 36 themed weeks. Each should be:

- 5–15 minutes duration
- Zero parent involvement needed
- Thematically linked to the "With You" activity
- Physical medium: colouring, drawing, sticking, tracing, cutting
- Optional QR code → audio story or music

Examples of what these could be:

| Theme | "While You..." idea |
|-------|---------------------|
| Near and Far | Colour two houses — one close, one far. Draw who lives there. |
| Loud and Quiet | Colour a page where half is loud (bright colours) and half is quiet (soft colours). |
| Big Feelings | Colour the feelings thermometer. Stick stickers at different levels. |
| Yes and No | Colour traffic-light faces: green (yes), red (no), amber (maybe). |
| Mine and Yours | Draw "my favourite things" on one side, "things I'd lend" on the other. |
| Try and Wait | Colour a "bravery badge." Stick it when ready. |
| Same and Different | Spot-the-difference colouring page. |
| Your Feelings, My Feelings | Draw two faces — one for "me today" and one for "someone I know." |
| Tidy and Messy | Colour a room that's messy on one side, tidy on the other. |
| My Fault, Your Fault | Draw the "pie chart" of whose bit was whose. |
| Keeping and Letting Go | Draw something in the memory box. |
| Wanting and Waiting | Colour in the waiting hourglass — one grain at a time. |

---

## 8. Parent Reference Card

Included with each edition. Printed on heavier card stock.

### Language Shifts

| Old (Moral) | New (Engineering) |
|-------------|-------------------|
| "You are..." | "Right now you're..." |
| "You should..." | "Your brain is saying..." |
| "That was wrong" | "That didn't work. What happened?" |
| "Try harder" | "What's blocking you right now?" |
| "Stop it" | "I notice you're in X setting — is that helping?" |
| "Good boy/girl" | "That worked. You figured it out." |
| "Don't cry" | "That's a big feeling. I'm here." |

### Play-First Principle

1. **Observe** the setting without judgement
2. **Name** it neutrally ("You're in Independent mode right now")
3. **Play** the activity that gently expands range
4. **Notice** what happened ("When you had to wait, what happened?")

### Bad Day Rule

> Bad days are when the dance gets heavy. Don't force it. Close the book. Read them a story. The dance will be there tomorrow. Bad days are data about sleep, hunger, or timing — not about your child.

---

## 9. Ventral Driver Stress Tests (Bonus Content)

These are standalone spread inserts for parents who've completed multiple cycles and want targeted work. One per driver archetype:

| Driver | Activity | Deprivation Target | Core Skill |
|--------|----------|--------------------|------------|
| ⭐ The Impressor | "The Invisible Artist" — create while parent looks away | Parent's gaze | Internal reward |
| 🔍 The Explorer | "The Pause Button" — stop mid-investigation and verbalise | Uninterrupted exploration | Disengagement |
| 🛡️ The Sentinel | "The Tiny Surprise" — one change in a familiar routine | Absolute predictability | Change tolerance |
| 🤝 The Connector | "The Kind Thought" — send a thought instead of fixing | Permission to fix | Emotional boundaries |
| ⚡ The Rebel | "The Two Doors" — choose within constraints | Unlimited choice | Agency within structure |
| 🎭 The Performer | "The Whisper Game" — exciting activity in whispers | Volume as outlet | Voltage containment |
| 👁️ The Observer | "The Body First Game" — jump before understanding | Full analysis | Action-before-analysis |

---

## 10. Choose-Your-Own-Adventure Book 1: "Jem's Big Day Out" (Social Gravity)

### Design Principles

| Principle | Detail |
|-----------|--------|
| **Character** | Jem — gender-neutral, species-neutral (could be illustrated as a child, an animal, a robot — TBD). NOT the reader. A proxy. |
| **Setting** | A party — high social density, natural energy drain, familiar to every toddler |
| **Two entry points** | Cohesive Jem (page 3) / Independent Jem (page 5). Parent picks the one that matches their child. |
| **Funnel structure** | Branches out then reconverges. Max 3 paths at any point. All paths reach the ending. |
| **Staged complexity** | Pages 1-8: Level 1-2 (foundational skills). Pages 9-14: Level 3 (offspring skills). Pages 15-20: Level 4-5 (advanced). |
| **No wrong answers** | Every choice leads somewhere. Some paths are harder. None are punished. "Go back and try again" is always available. |
| **Read-aloud register** | Short sentences. Simple words. Parent reads, child points/chooses. |
| **Page count** | 20 story pages + cover + "about this book" page for parents = 24 pages total |

### Assessment Map (📱 app-side only — invisible to child)

Each choice point maps to skills. The app records which path the child takes.

| Choice point | Page | Skills assessed | Level |
|---|---|---|---|
| Entry point | 2 | *None — parent selects based on child's profile* | — |
| "How does Jem feel?" | 3 or 5 | Emotion labelling | 1 |
| "What does Jem's face look like?" | 6 | Observation, emotion reading | 1 |
| "What should Jem do?" (obvious) | 8 | Self-awareness, energy reading | 2 |
| "What happens next?" (predict) | 10 | Consequence anticipation | 3 |
| "What should Jem say?" (social) | 12 | Social initiation OR saying no | 3 |
| "Both friends want different things" | 15 | Perspective-taking, negotiation | 4 |
| "Jem doesn't know what will happen" | 17 | Ambiguity tolerance, inhibitory control | 4 |
| "Someone's feelings might get hurt" | 19 | Empathy, trade-off reasoning | 5 |

---

### THE BOOK

---

#### PAGE 1 — Title Page

> # Jem's Big Day Out
>
> *A story where YOU choose what happens*

*[Illustration: Jem standing at a garden gate, looking in. Balloons visible. Sounds of children.]*

---

#### PAGE 2 — Entry Point (parent reads, child chooses)

> Jem has been invited to a party!
>
> But here's the thing about Jem...
>
> **Does Jem LOVE being with lots of people?**
> *→ Turn to page 3*
>
> **Does Jem like having quiet time on their own?**
> *→ Turn to page 5*

*[Illustration: Two pictures of Jem — one surrounded by smiling friends, one sitting contentedly alone with a toy. Both happy. Equal size. Neither is "better."]*

---

### PATH A — COHESIVE JEM (stretches toward independent skills)

---

#### PAGE 3 — Cohesive Jem arrives ⟨Level 1: Emotion labelling⟩

> Jem LOVES being with people!
>
> Jem runs through the gate. There are children everywhere!
> Music is playing! Someone is laughing!
>
> **How does Jem feel right now?**
>
> 😊 "Jem feels EXCITED!" *→ Turn to page 7*
>
> 😍 "Jem feels SO HAPPY!" *→ Turn to page 7*

*[Illustration: Jem running in, arms wide, huge grin. Other children playing in background. Bunting and balloons. A small battery icon in the corner showing FULL.]*

📱 *App note: Both choices are correct — they demonstrate emotion labelling. The distinction (excited vs happy) is recorded but not scored differently at Level 1. The point is that the child CAN name an emotional state.*

---

#### PAGE 4 — Cohesive Jem plays ⟨Level 1: Observation⟩

> Jem plays with EVERYONE.
>
> Jem dances! Jem chases! Jem shares snacks!
>
> But look at Jem's battery...
>
> *[Battery icon now showing HALF]*
>
> **Look at Jem's face. Can you see anything different?**
>
> 😃 "Jem still looks happy!" *→ Turn to page 8*
>
> 😐 "Jem looks a bit tired." *→ Turn to page 9*

*[Illustration: Jem still playing, but shoulders slightly dropped, smile slightly smaller. Subtle visual cues. The other children still look energetic.]*

📱 *App note: This tests observation beyond the obvious. The child who spots "tired" despite Jem still playing demonstrates reading beyond surface presentation. Score: Observation 2+ if they pick "tired," Observation 1 if they pick "happy" (not wrong — just less nuanced).*

---

### PATH B — INDEPENDENT JEM (stretches toward cohesive skills)

---

#### PAGE 5 — Independent Jem arrives ⟨Level 1: Emotion labelling⟩

> Jem likes having quiet time.
>
> Jem walks through the gate. It's... LOUD.
> There are children everywhere. Music is playing.
>
> **How does Jem feel right now?**
>
> 😟 "Jem feels a bit worried." *→ Turn to page 6*
>
> 🤔 "Jem feels not sure." *→ Turn to page 6*

*[Illustration: Jem standing just inside the gate, holding the fence, looking at the party. Not crying — just assessing. Battery icon showing FULL but with a small question mark.]*

📱 *App note: Both choices valid. "Worried" shows direct emotion labelling. "Not sure" shows awareness of ambiguity — sophisticated for this age. Record which: it maps to emotional vocabulary range.*

---

#### PAGE 6 — Independent Jem watches ⟨Level 1: Observation⟩

> Jem stands near the fence and watches.
>
> Over there — two children are building something with blocks.
> They look like they're having fun. But there's no space to sit down.
>
> Over HERE — one child is drawing, all on their own.
> There's space next to them.
>
> **Which one does Jem look at?**
>
> 🧱 "The children with blocks." *→ Turn to page 9*
>
> 🖍️ "The child who's drawing." *→ Turn to page 8*

*[Illustration: Jem watching from the side. Two clear scenes visible — a lively group with blocks (no space), a single child drawing (space next to them). Jem's eyes looking between them.]*

📱 *App note: Neither is wrong. Blocks = attracted to group energy (surprising for independent Jem — possible stretch readiness). Drawing = gravitating to safe parallel play (expected). The choice reveals whether the child identifies with social aspiration or comfort-seeking.*

---

### CONVERGENCE ZONE — Both paths merge here

---

#### PAGE 7 — The party is great! ⟨Level 2: Energy reading⟩

> The party is SO much fun!
>
> Jem has been playing for a long, long time.
> Jem's legs are tired. Jem's voice is quiet.
>
> *[Battery icon showing LOW — almost empty]*
>
> **Jem's battery is nearly empty. What should Jem do?**
>
> 🎉 "Keep playing! The party is fun!" *→ Turn to page 10*
>
> 🪑 "Find somewhere quiet to sit." *→ Turn to page 11*
>
> 🧸 "Find their favourite toy and hold it." *→ Turn to page 11*

*[Illustration: Jem mid-activity but visibly running low — yawning slightly, holding own arm, other children still going strong.]*

📱 *App note: This is the key Social Gravity assessment point. "Keep playing" = can't read energy state OR can read it but can't stop (inhibitory control gap). "Quiet spot" = reads state + acts on it (self-awareness + self-regulation). "Favourite toy" = self-soothing strategy (advanced for this age). Record which: maps to self-awareness (1-3) and inhibitory control (1-2).*

---

#### PAGE 8 — Parallel play ⟨Level 2: Self-awareness⟩

> Jem sits next to the child who's drawing.
>
> They don't talk. They just draw. Side by side.
>
> It feels... nice.
>
> **How does Jem feel now?**
>
> 😌 "Calm." *→ Turn to page 11*
>
> 🤔 "Jem wants to see what they're drawing." *→ Turn to page 12*
>
> 😕 "Jem wants to go play with the big group." *→ Turn to page 7*

*[Illustration: Jem and the other child sitting side by side, each drawing. Not looking at each other but clearly comfortable. Small smiles. Both batteries half-full.]*

📱 *App note: "Calm" = content with parallel play (independent-pole comfort). "Wants to see" = reaching toward joint attention (cohesive stretch — positive). "Wants the big group" = not ready for quiet yet — loops them to page 7 for the energy-reading choice. All valid. Maps to: comfortable silence (1-2), joint attention readiness (1-2).*

---

#### PAGE 9 — Joining in ⟨Level 2: Social initiation⟩

> Jem walks toward the children with blocks.
>
> They're building a tall tower. It looks wobbly!
>
> Jem wants to help. But everyone is busy.
>
> **What does Jem do?**
>
> 👋 "Jem says 'Can I play?'" *→ Turn to page 12*
>
> 👀 "Jem watches and waits." *→ Turn to page 8*
>
> 🧱 "Jem just picks up a block and starts building." *→ Turn to page 13*

*[Illustration: Three children building, backs partially turned. Gap in the circle but no one's looking at Jem. Jem standing nearby holding nothing.]*

📱 *App note: "Can I play?" = verbal social initiation (the hardest cohesive skill for independent children — high score if chosen). "Watches and waits" = defaults to observation, avoids initiation (expected for independent profile, not wrong). "Just starts building" = non-verbal initiation, skips asking permission (impulsive but socially brave — maps to inhibitory control low, social initiation medium). Record which: maps to social initiation (1-3), inhibitory control (1-2).*

---

#### PAGE 10 — Kept playing (consequence) ⟨Level 3: Consequence anticipation⟩

> Jem kept playing... and playing... and playing.
>
> Then Jem's battery ran completely empty.
>
> *[Battery icon: EMPTY. Red.]*
>
> Jem sat down on the ground and started crying.
> The party didn't feel fun anymore.
>
> **What happened?**
>
> 😢 "Jem played too long." *→ Turn to page 11*
>
> 🔋 "Jem's battery ran out." *→ Turn to page 11*

*[Illustration: Jem sitting on the ground, tears, other children looking confused/concerned. One child reaching toward Jem. The contrast with page 7's energy is stark.]*

> **Want to try again?** *→ Go back to page 7*

📱 *App note: This is NOT punishment — it's visible consequence. The "go back" option is key: it teaches that decisions can be replayed. The child who goes back and chooses differently has just demonstrated consequence anticipation + behavioural adjustment. If they choose the same thing again, they may not yet connect action to outcome — that's a Level 1 gap (cause and effect), not a character flaw. Record: "went back" = consequence anticipation emerging. "Didn't go back" = may need parent to narrate the connection.*

---

#### PAGE 11 — Quiet moment ⟨Level 3: Recovery + Prediction⟩

> Jem found a quiet spot.
>
> Jem sat down. Jem took a big breath.
>
> Slowly... slowly... the battery started filling up again.
>
> *[Battery icon: filling from LOW to HALF, animated feel]*
>
> A grown-up came over. "Are you OK, Jem?"
>
> **What does Jem say?**
>
> ✅ "I'm OK. I just needed a rest." *→ Turn to page 14*
>
> 🤷 "I don't know." *→ Turn to page 14*
>
> 🙂 "I'm ready to play again!" *→ Turn to page 15*

*[Illustration: Jem sitting against a tree/wall, breathing, eyes closed or looking at sky. Battery visibly rising. A kind adult crouching nearby, not too close.]*

📱 *App note: "I just needed a rest" = can name own state + communicates need (self-awareness 3, requesting what they need 2). "I don't know" = honest but can't yet verbalise the state (self-awareness 1-2). "Ready to play!" = recovery is fast, or they're masking (context-dependent — if they came from page 10, it might be avoidance of the quiet). Record which: maps to self-awareness (1-3), requesting what they need (1-2).*

---

#### PAGE 12 — Making a connection ⟨Level 3: Joint attention + Turn-taking⟩

> Jem is with another child now.
>
> They're drawing together. The other child draws a cat.
>
> "Look at my cat!" they say.
>
> **What does Jem do?**
>
> 👀 "Jem looks at the cat and says 'I like it!'" *→ Turn to page 14*
>
> 🖍️ "Jem draws a cat too — right next to it." *→ Turn to page 14*
>
> 🙈 "Jem keeps drawing their own picture." *→ Turn to page 13*

*[Illustration: Two children at a table. One holding up a drawing proudly. Jem looking at it (or not — depending on which version the child imagines).]*

📱 *App note: "Looks and says I like it" = verbal joint attention + social reinforcement (joint attention 2, social initiation 2). "Draws a cat too" = non-verbal joint attention through imitation (joint attention 2, imitation 2 — very natural at this age). "Keeps drawing own picture" = not yet ready for shared focus — loops to page 13 for a gentler approach. Record: maps to joint attention (1-3), turn-taking readiness (1-2).*

---

#### PAGE 13 — On your own at the party ⟨Level 3: Comfortable aloneness⟩

> Jem is on their own at the party.
>
> Everyone else is playing together. Jem is playing alone.
>
> Jem is building something. It's getting really good.
>
> **How does Jem feel?**
>
> 😊 "Happy! Jem likes building on their own." *→ Turn to page 15*
>
> 😕 "A bit lonely. Jem wants someone to see." *→ Turn to page 14*
>
> 🤔 "Jem doesn't know yet." *→ Turn to page 14*

*[Illustration: Jem alone at a table, building something impressive. Other children visible in background, playing. Jem is focused, not looking at them.]*

📱 *App note: "Happy alone" = genuine independent-pole comfort (solitary play 2+). "Wants someone to see" = needs external validation (status awareness emerging, connection maintenance 1). "Doesn't know" = honest uncertainty — the most sophisticated answer for this age, shows emerging metacognition. Record: maps to solitary play (1-3), internal narration (1-2), selective engagement (1-2).*

---

#### PAGE 14 — The middle of the party ⟨Level 4: Perspective-taking⟩

> Jem's battery is about half full now.
>
> *[Battery icon: HALF]*
>
> Two friends come over. Sam and Alex.
>
> Sam says: "Come and play chase!"
> Alex says: "Come and draw with me."
>
> They both want Jem. But Jem can't do both.
>
> **What does Jem do?**
>
> 🏃 "Go with Sam — chase sounds fun!" *→ Turn to page 16*
>
> 🖍️ "Go with Alex — drawing is quieter." *→ Turn to page 16*
>
> 🤝 "Ask if we can all do something together." *→ Turn to page 17*
>
> ⏸️ "Say 'Can I choose in a minute?'" *→ Turn to page 17*

*[Illustration: Jem in the middle. Sam on one side bouncing/energetic. Alex on the other sitting with crayons. Both looking at Jem hopefully. Jem looking between them.]*

📱 *App note: This is the first Level 4 choice — genuine social complexity. "Chase" = chose high-energy (check battery — risky if low). "Drawing" = chose recovery-compatible activity (energy-aware). "All together" = attempted negotiation/compromise (advanced — perspective-taking 3+). "Can I choose in a minute" = delayed decision = inhibitory control + self-awareness 3+ (knows they need time). This is rich assessment data. Record all: maps to perspective-taking (1-4), negotiation (1-3), inhibitory control (1-3), energy reading (2-3).*

---

#### PAGE 15 — Jem has energy again ⟨Level 4: Selective engagement⟩

> Jem's battery is filling up!
>
> *[Battery icon: nearly FULL]*
>
> The party is still going. There are so many things happening.
>
> Over there: a big loud game with LOTS of children.
> Over here: two children doing a puzzle.
> In the corner: one child sitting alone, looking down.
>
> **Where does Jem go?**
>
> 🎉 "The big loud game!" *→ Turn to page 16*
>
> 🧩 "The puzzle with two children." *→ Turn to page 17*
>
> 💛 "The child sitting alone." *→ Turn to page 18*

*[Illustration: Wide shot of the party. Three distinct scenes at different energy levels. Jem in the foreground, battery full, looking across all three. The child in the corner is drawn small but noticeable.]*

📱 *App note: "Big loud game" = high energy, full battery, confident (selective engagement 2). "Puzzle" = moderate, calibrated (selective engagement 3 — chose appropriate energy match). "Child alone" = noticed someone else's state and responded (empathy 3+, reading social cues 3+). This is the most diagnostically rich choice in the book. The child who spots the lonely kid WITHOUT being told to look demonstrates spontaneous empathy — score 4+. Record: maps to selective engagement (2-4), empathy (1-4), reading social cues (1-3), adapting to group energy (1-3).*

---

#### PAGE 16 — Big energy ⟨Level 4: Adapting to group energy⟩

> Jem joins the big game!
>
> Everyone is running and laughing.
>
> But one child trips and falls. They start crying.
>
> Nobody else stops. They keep running.
>
> **What does Jem do?**
>
> 🏃 "Keep running! The game is fun!" *→ Turn to page 19*
>
> 🛑 "Stop and check on the child." *→ Turn to page 18*
>
> 📢 "Shout 'STOP! Someone fell!'" *→ Turn to page 18*

*[Illustration: A chaotic running game. One child on the ground crying. Others still running past. Jem has noticed — feet pointing one way, head turned back toward the fallen child.]*

📱 *App note: "Keep running" = absorbed by group energy, didn't/couldn't break from it (adapting to group energy 1, empathy 1 — not cruel, just carried). "Stop and check" = broke from group momentum + empathy response (inhibitory control 3, empathy 3). "Shout STOP" = leadership move, tries to change the group's behaviour (leading an activity 3+, assertiveness 3). Record: maps to empathy (1-4), inhibitory control (1-3), adapting to group energy (1-3), leading (0-3).*

---

#### PAGE 17 — Finding the middle ⟨Level 4: Negotiation + compromise⟩

> Jem is with a small group now.
>
> They can't agree on what to do.
>
> One child wants to play outside. One child wants to stay inside.
>
> They're both getting upset.
>
> **What does Jem say?**
>
> 🌗 "What about playing near the door? A bit of both!" *→ Turn to page 19*
>
> 👂 "Why do you want to go outside? Why do you want to stay in?" *→ Turn to page 19*
>
> 🤷 "I don't know how to fix this." *→ Turn to page 19*

*[Illustration: Two children facing each other, arms crossed. One pointing at the door, one pointing at the toys inside. Jem in between them, uncertain.]*

📱 *App note: "Near the door" = proposed practical compromise (negotiation 3+, problem-solving 3). "Why do you want..." = asked questions before solving — the most sophisticated response (perspective-taking 4, questioning 3 — they want to understand before acting). "I don't know" = honest about limits — NOT a failure, shows self-awareness of competence boundary (self-awareness 3). All three lead to the same page — the lesson is that all three are reasonable approaches. Record: maps to negotiation (1-4), perspective-taking (2-4), self-awareness (2-3).*

---

#### PAGE 18 — Someone needs help ⟨Level 5: Empathy in action⟩

> Jem goes to the child who's sitting alone.
>
> "Are you OK?" says Jem.
>
> The child says: "My battery ran out. I want to go home. But the party isn't finished."
>
> **What does Jem do?**
>
> 🤝 "Jem sits next to them. They don't talk. They just sit." *→ Turn to page 20*
>
> 💬 "Jem says: 'That happens to me too.'" *→ Turn to page 20*
>
> 🧸 "Jem brings them a toy." *→ Turn to page 20*
>
> 👋 "Jem says: 'I'll come back and check on you later.'" *→ Turn to page 20*

*[Illustration: Jem crouching next to the sitting child. The sitting child's battery icon is EMPTY. Jem's is HALF. The sitting child isn't crying — just flat, still, looking down. Jem is looking AT them.]*

📱 *App note: This is the Level 5 assessment. Every choice is a valid form of empathy — they're testing WHICH KIND:*
- *"Sits next to them" = presence without demand. The most advanced form. Comfortable silence + empathy + not trying to fix. (Empathy 4+, comfortable silence 3+)*
- *"That happens to me too" = normalising through shared experience. Connection maintenance. (Empathy 3, connection maintenance 3)*
- *"Brings a toy" = practical help. Action-oriented empathy. (Empathy 2-3, helping 2)*
- *"I'll come back later" = recognises both the need AND their own limits. (Empathy 3, self-awareness 4, connection maintenance 3 — arguably the most mature response)*
*Record all. There is genuinely no wrong answer here. The TYPE of empathy chosen reveals the child's social-emotional profile better than any questionnaire ever could.*

---

#### PAGE 19 — Heading home ⟨Level 5: Reflection⟩

> The party is ending.
>
> Everyone is going home.
>
> Jem had a big day. Some parts were fun. Some parts were hard.
>
> **What was the BEST thing about Jem's day?**
>
> 🎉 "Playing with everyone!" *→ Turn to page 20*
>
> 🖍️ "Drawing with one friend." *→ Turn to page 20*
>
> 🧘 "The quiet time." *→ Turn to page 20*
>
> 💛 "Helping someone." *→ Turn to page 20*

*[Illustration: Jem walking out of the gate, looking back at the party. Half-empty plates on tables, balloons slightly deflated. A warm, tired, end-of-day feeling. Jem's battery icon shows HALF — not empty, not full. A good day.]*

📱 *App note: This is pure reflection — the child chooses what resonated. "Playing with everyone" = group energy was the highlight (cohesive-pole preference). "Drawing with one friend" = one-to-one connection (moderate). "Quiet time" = valued the recovery/alone phase (independent-pole preference). "Helping someone" = social-emotional fulfilment (empathy-driven). This confirms the ETP profile and shows whether the stretch activities were positive or draining. Record: maps to reflection (1-3), preference awareness (2-3), and confirms Social Gravity home position.*

---

#### PAGE 20 — The ending (same for all paths)

> Jem is home now.
>
> Jem's battery isn't full. And it isn't empty.
>
> It's right in the middle.
>
> *[Battery icon: exactly HALF — balanced]*
>
> That's a good place to be.
>
> Tomorrow, Jem might choose different things.
> And that's OK too.
>
> **The End**
>
> *Want to read it again? Go back to page 2 and try different choices!*

*[Illustration: Jem at home, in bed or on a sofa, relaxed, holding the same toy from the party. Content. Through the window, you can see the party venue in the distance with the lights going off. Jem's battery is perfectly mid-range. Peace.]*

📱 *App note: If the child wants to re-read immediately and try different paths, that's the highest possible score for curiosity and cognitive flexibility. Record: "re-read requested" = Y/N.*

---

### Page Map — All Paths Visualised

```
                        [1: Title]
                            │
                        [2: Entry]
                       /          \
              [3: Cohesive]    [5: Independent]
                   │                │
              [4: Plays]       [6: Watches]
                   │              / │
                   └──── [7] ───┘  │
                        / │        │
                   [10]  [11]  ← [8: Parallel]
                    │     │    ↗   │
                    └─────┼──┘    [12: Connect]
                          │      / │
                         [14] ←┘  [13: Alone]
                        / │ \      │
                   [16] [17] └─── [15: Choose]
                    │    │       / │ \
                    │    │   [16] [17] [18]
                    └────┼────┘   │    │
                         └────────┼────┘
                               [19: Reflect]
                                  │
                               [20: End]
```

**Total unique pages: 20**
**Maximum path length: 10 pages (short read)**
**Minimum path length: 7 pages**
**Branch points: 9**
**Reconvergence points: 4 (pages 7, 14, 19, 20)**

---

### Production Notes

| Item | Detail |
|------|--------|
| **Format** | Landscape A5 or square 200mm. Thick card pages for durability. |
| **Illustration style** | Warm, soft, inclusive. Jem should look like EVERY child could be Jem. |
| **Battery icon** | Consistent visual on every page. Children learn to read it. This becomes vocabulary. |
| **Text per page** | Maximum 40 words. Target 25. |
| **Font** | Large, rounded, dyslexia-friendly. |
| **App integration** | Each page has a small QR code or page number. Parent taps the page number in the app, selects the child's choice. Takes 3 seconds. |
| **Without the app** | The book works perfectly as a standalone story. The assessment layer is invisible without the app. |
| **Series potential** | 12 books (one per spectrum). Same character (Jem), different settings. Same structure. |
| **Price point** | £6-8 per book. Set of 12: £50-60. Subscription: 1 per month. |

---

## 11. File References

| File | Content |
|------|---------|
| [`pathway_flexibility.go`](../esp-organizer/internal/domain/etp/pathway_flexibility.go) | All 36 themed weeks (12 spectra × 3 cycles) with full activity content, sticker choices, navigation prompts, try-instead variations |
| [`activity_experience.go`](../esp-organizer/internal/domain/etp/activity_experience.go) | Emotional compiler — multi-threaded activity schema with side effects and counterweights. 3 sample activities. |
| [`response_matrix.go`](../esp-organizer/internal/domain/etp/response_matrix.go) | 24 response matrix entries — PlayAvenues, RecommendedLanguage, AvoidLanguage |
| [`ventral_drivers.go`](../esp-organizer/internal/domain/etp/ventral_drivers.go) | 7 driver personas, spectrum signatures, stress tests |
| [`assessment_loop.go`](../esp-organizer/internal/domain/etp/assessment_loop.go) | Three-Check Sensor, noise filtering, Range of Motion rosette |
| [`spectra.go`](../esp-organizer/internal/domain/etp/spectra.go) | 12 spectra definitions, 61 trainable skills |
| [`ToddlerOS/data/activities.ts`](../../ToddlerOS/data/activities.ts) | TypeScript port of all 36 theme weeks for the Expo app |
| [`TODDLEROS_TECHNICAL_OVERVIEW.md`](TODDLEROS_TECHNICAL_OVERVIEW.md) | Full technical architecture connecting ToddlerOS to ETP, skills graph, pathfinder |



