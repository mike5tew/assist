package etp

import "time"

// =============================================================================
// Assessment Feedback Loop — "The Control Loop"
// =============================================================================
//
// Engineering frame:
//
//   INPUT              TRANSFER FUNCTION           OUTPUT
//   (Activity)  →  (Child's Spectra + Drivers)  →  (Consequence)
//        ↑                                              │
//        │              FEEDBACK SENSOR                  │
//        └──────────── (Assessment) ←───────────────────┘
//
// Assessment in ToddlerOS is NOT a test. It is Differential Observation —
// measuring the DELTA between expected reaction and actual reaction.
//
// The system uses three sensors, each measuring a different dimension:
//
//   1. RECOVERY SENSOR   — "Latency to Baseline"
//      How long to get back to normal after a big feeling?
//      Measures: Pilot Strength (inhibitory control, regulation)
//
//   2. FRICTION SENSOR   — "Voltage Resistance"
//      Was the transition to the activity a grind or a flow?
//      Measures: Spectrum Alignment (is this too far from natural pole?)
//
//   3. AGENCY SENSOR     — "Internal Locus"
//      Did they look at ME for the win, or look at the TASK?
//      Measures: Driver Moderation (is the ventral driver sole power source?)
//
// Each sensor produces a simple signal that the parent can observe without
// technical knowledge. The app converts these into spectrum adjustments
// and driver confidence updates.
//
// NOISE FILTERING: Before any signal is processed, the system checks
// Current Load (tired/hungry/ill). Hardware noise is discarded.
// Software signals are kept.
// =============================================================================

// =============================================================================
// Core Types
// =============================================================================

// ObservationEvent is a single post-activity observation logged by the parent.
// This is the atomic unit of the feedback loop — one activity, one child,
// one set of sensor readings.
type ObservationEvent struct {
	ID        string    `json:"observation_id"`
	Timestamp time.Time `json:"timestamp"`

	// Which activity was performed?
	ActivityID string `json:"activity_id"` // → Activity.ID

	// Which child? (for multi-child households)
	ChildID string `json:"child_id"`

	// === SKIP FLAG ===
	// "Bad day? Skip it." — if true, this observation is not processed.
	SkippedDay bool `json:"skipped_day"`

	// === THE THREE SENSORS ===
	Recovery SensorReading `json:"recovery_sensor"`
	Friction SensorReading `json:"friction_sensor"`
	Agency   SensorReading `json:"agency_sensor"`

	// === FREEFORM ===
	// Optional parent note — captured for future pattern mining
	ParentNote string `json:"parent_note,omitempty"`

	// === SYSTEM-COMPUTED FIELDS ===
	// These are filled by ProcessObservation(), not by the parent.
	IsNoise             bool               `json:"is_noise"`           // true if SkippedDay was set
	SpectrumInsights    []SpectrumInsight  `json:"spectrum_insights"`  // derived spectrum adjustments
	DriverAdjustments   []DriverAdjustment `json:"driver_adjustments"` // driver confidence shifts
	PatchRecommendation *ActivityPatch     `json:"patch,omitempty"`    // suggested next-activity adjustment
}

// SensorReading is a single sensor observation. The parent answers ONE
// simple question; the system maps it to a numeric value.
type SensorReading struct {
	// The human-friendly question shown to the parent
	Question string `json:"question"`

	// The parent's response — always a simple choice
	Response string `json:"response"` // varies per sensor (see SensorQuestions)

	// System-computed: normalised value from the response
	// -1.0 (concern) to +1.0 (thriving), 0.0 = neutral
	Value float64 `json:"value"`
}

// SpectrumInsight is a derived adjustment to a spectrum reading based on
// observation. This feeds back into the ETP profile.
type SpectrumInsight struct {
	SpectrumID   int     `json:"spectrum_id"`
	SpectrumName string  `json:"spectrum_name"`
	Delta        float64 `json:"delta"`  // change to apply: -0.5 to +0.5
	Reason       string  `json:"reason"` // human-readable explanation
	Source       string  `json:"source"` // which sensor produced this: "recovery", "friction", "agency"
}

// DriverAdjustment updates the confidence score for a ventral driver based
// on observed behaviour. Feeds into DetectLikelyDrivers().
type DriverAdjustment struct {
	DriverID        string  `json:"driver_id"`
	ConfidenceDelta float64 `json:"confidence_delta"` // -0.2 to +0.2
	Reason          string  `json:"reason"`
}

// ActivityPatch is the system's recommendation for what to do NEXT based on
// the observation. This is how the feedback loop closes — the assessment
// modifies the next input.
type ActivityPatch struct {
	// What to adjust in the next session
	PatchType string `json:"patch_type"` // "increase_predictability", "reduce_audience", "extend_duration", "simplify", "escalate"

	// Suggested next activity (if specific)
	SuggestedActivityID string `json:"suggested_activity_id,omitempty"`

	// Human-readable recommendation for the parent
	ParentMessage string `json:"parent_message"`
	// e.g., "Next time, try giving a 10-second countdown before the music stops.
	//        Their brain needs a bridge between 'going' and 'stopped'."

	// How many sessions this patch should persist before re-evaluating
	PatchDuration int `json:"patch_duration"` // number of activity sessions

	// Which spectrum/driver this patch targets
	TargetSpectrumID int    `json:"target_spectrum_id,omitempty"`
	TargetDriverID   string `json:"target_driver_id,omitempty"`
}

// =============================================================================
// The Three Sensors — Question Banks
// =============================================================================
//
// Each sensor has 9 variants — one per spectrum theme. The parent sees
// ONE question per sensor after each activity, phrased in the language
// of the current theme week.
// =============================================================================

// RecoverySensorQuestion measures Latency to Baseline for a given spectrum theme.
// "How long did it take to get back to normal after the big feeling?"
type RecoverySensorQuestion struct {
	SpectrumID int    `json:"spectrum_id"`
	ThemeName  string `json:"theme_name"`
	Question   string `json:"question"`
	// Response options (parent picks one)
	Options []SensorOption `json:"options"`
}

// FrictionSensorQuestion measures Voltage Resistance — alignment between
// activity and natural pole.
// "Was the transition to the activity a grind or a flow?"
type FrictionSensorQuestion struct {
	SpectrumID int            `json:"spectrum_id"`
	ThemeName  string         `json:"theme_name"`
	Question   string         `json:"question"`
	Options    []SensorOption `json:"options"`
}

// AgencySensorQuestion measures Internal Locus — where the "win" signal
// came from.
// "Did they look at me for the win, or look at the task?"
type AgencySensorQuestion struct {
	SpectrumID int            `json:"spectrum_id"`
	ThemeName  string         `json:"theme_name"`
	Question   string         `json:"question"`
	Options    []SensorOption `json:"options"`
}

// SensorOption is a single response option with a display label and
// a numeric value for the system.
type SensorOption struct {
	Label string  `json:"label"` // What the parent sees
	Value float64 `json:"value"` // System value: -1.0 to +1.0
}

// =============================================================================
// Sensor Question Banks — 9 questions per sensor (one per spectrum theme)
// =============================================================================

// RecoverySensorQuestions — one per spectrum theme.
// WHAT IT MEASURES: Pilot Strength (inhibitory control, regulation).
// The parent observes how quickly the child returns to baseline after the
// activity triggered a big feeling.
var RecoverySensorQuestions = []RecoverySensorQuestion{
	{
		SpectrumID: 1, ThemeName: "Near and Far",
		Question: "When playing together got too much (or not enough), how quickly did they bounce back?",
		Options: []SensorOption{
			{Label: "Didn't need to — they were fine throughout", Value: 1.0},
			{Label: "A wobbly moment, then back to playing within a minute", Value: 0.5},
			{Label: "Took a few minutes of quiet before they were ready again", Value: 0.0},
			{Label: "Struggled for a long time — needed a lot of help to settle", Value: -0.5},
			{Label: "Couldn't recover — the session had to stop", Value: -1.0},
		},
	},
	{
		SpectrumID: 2, ThemeName: "Loud and Quiet",
		Question: "When the noise level changed (louder or quieter than they wanted), how fast did they adjust?",
		Options: []SensorOption{
			{Label: "Adjusted easily — barely noticed", Value: 1.0},
			{Label: "Brief grumble, then settled", Value: 0.5},
			{Label: "Needed help to find their comfortable level", Value: 0.0},
			{Label: "Got quite upset — took real effort to come back", Value: -0.5},
			{Label: "Couldn't handle it — hands over ears / ran away", Value: -1.0},
		},
	},
	{
		SpectrumID: 3, ThemeName: "Big Feelings, Small Feelings",
		Question: "After the biggest feeling moment, how long until they felt 'normal' again?",
		Options: []SensorOption{
			{Label: "Seconds — they surfed the wave", Value: 1.0},
			{Label: "A minute or two — needed a pause", Value: 0.5},
			{Label: "Several minutes — needed comfort", Value: 0.0},
			{Label: "A long time — the feeling stuck around", Value: -0.5},
			{Label: "Didn't recover during the activity", Value: -1.0},
		},
	},
	{
		SpectrumID: 4, ThemeName: "Yes and No",
		Question: "When they hit a 'no' moment (theirs or yours), how quickly did they find their feet again?",
		Options: []SensorOption{
			{Label: "Took it in stride — barely paused", Value: 1.0},
			{Label: "Brief freeze or flare, then moved on", Value: 0.5},
			{Label: "Needed a moment — but came back with help", Value: 0.0},
			{Label: "Got stuck in the 'no' — hard to shift", Value: -0.5},
			{Label: "Full shutdown or meltdown — session ended", Value: -1.0},
		},
	},
	{
		SpectrumID: 5, ThemeName: "Mine and Yours",
		Question: "When sharing was hard, how long before they could re-engage?",
		Options: []SensorOption{
			{Label: "Sharing came naturally — no issue", Value: 1.0},
			{Label: "One tricky moment, then generous again", Value: 0.5},
			{Label: "Needed reminding but got there", Value: 0.0},
			{Label: "Kept pulling things back — couldn't let go", Value: -0.5},
			{Label: "Sharing triggered a meltdown", Value: -1.0},
		},
	},
	{
		SpectrumID: 6, ThemeName: "Try and Wait",
		Question: "When the activity felt risky or uncertain, how quickly did they settle into it?",
		Options: []SensorOption{
			{Label: "Jumped straight in — no hesitation", Value: 1.0},
			{Label: "Watched for a moment, then joined", Value: 0.5},
			{Label: "Needed encouragement but eventually tried", Value: 0.0},
			{Label: "Took a lot of coaxing — very reluctant", Value: -0.5},
			{Label: "Refused to try", Value: -1.0},
		},
	},
	{
		SpectrumID: 7, ThemeName: "Same and Different",
		Question: "When the rules or routine changed, how quickly did they adapt?",
		Options: []SensorOption{
			{Label: "Rolled with it — flexible", Value: 1.0},
			{Label: "Noticed the change, brief wobble, then OK", Value: 0.5},
			{Label: "Needed the change explained before accepting it", Value: 0.0},
			{Label: "Got really upset — took a long time", Value: -0.5},
			{Label: "Could not accept the change at all", Value: -1.0},
		},
	},
	{
		SpectrumID: 8, ThemeName: "Your Feelings, My Feelings",
		Question: "If they picked up on someone else's big feeling, how quickly did they let it go?",
		Options: []SensorOption{
			{Label: "Noticed but didn't absorb it", Value: 1.0},
			{Label: "Got a bit wobbly but shook it off", Value: 0.5},
			{Label: "Carried it for a while — needed reassurance", Value: 0.0},
			{Label: "Really took it on — hard to separate their feeling from the other person's", Value: -0.5},
			{Label: "Completely overwhelmed by the other person's emotion", Value: -1.0},
		},
	},
	{
		SpectrumID: 9, ThemeName: "Tidy and Messy",
		Question: "When things got messy (or were forced to be tidy), how long before they felt OK about it?",
		Options: []SensorOption{
			{Label: "Didn't mind either way", Value: 1.0},
			{Label: "Brief resistance, then accepted", Value: 0.5},
			{Label: "Needed a transition activity to shift", Value: 0.0},
			{Label: "Really struggled — got quite stressed", Value: -0.5},
			{Label: "Could not tolerate it — distress", Value: -1.0},
		},
	},
	{
		SpectrumID: 10, ThemeName: "My Fault, Your Fault",
		Question: "When something went wrong, how quickly did they find a fair view of what happened?",
		Options: []SensorOption{
			{Label: "Saw their part clearly — fair and calm", Value: 1.0},
			{Label: "Wobbled but found a fair view with a nudge", Value: 0.5},
			{Label: "Hard to tell — avoided the topic", Value: 0.0},
			{Label: "Got stuck blaming others or blaming themselves", Value: -0.5},
			{Label: "Complete shutdown — refused to talk about it", Value: -1.0},
		},
	},
	{
		SpectrumID: 11, ThemeName: "Keeping and Letting Go",
		Question: "When something was taken away or lost, how long before they settled?",
		Options: []SensorOption{
			{Label: "Barely noticed — moved on instantly", Value: 1.0},
			{Label: "Brief upset, then found something else", Value: 0.5},
			{Label: "Needed comfort but eventually let go", Value: 0.0},
			{Label: "Really struggled — kept coming back to it", Value: -0.5},
			{Label: "Inconsolable — the loss felt catastrophic", Value: -1.0},
		},
	},
	{
		SpectrumID: 12, ThemeName: "Wanting and Waiting",
		Question: "When they wanted something, how well did they manage the gap between wanting and getting?",
		Options: []SensorOption{
			{Label: "Patient and measured — waited without fuss", Value: 1.0},
			{Label: "A bit restless but managed", Value: 0.5},
			{Label: "Needed help to wait", Value: 0.0},
			{Label: "Really struggled with waiting — got upset", Value: -0.5},
			{Label: "Could not wait at all — demanded it now", Value: -1.0},
		},
	},
}

// FrictionSensorQuestions — one per spectrum theme.
// WHAT IT MEASURES: Spectrum Alignment (is the activity too far from natural pole?).
// The parent observes whether the activity felt natural or forced.
var FrictionSensorQuestions = []FrictionSensorQuestion{
	{
		SpectrumID: 1, ThemeName: "Near and Far",
		Question: "Did the togetherness (or alone time) feel natural or forced?",
		Options: []SensorOption{
			{Label: "Totally natural — they leaned into it", Value: 1.0},
			{Label: "Mostly easy — a little stretching", Value: 0.5},
			{Label: "Neutral — neither easy nor hard", Value: 0.0},
			{Label: "Clearly pushed against their preference", Value: -0.5},
			{Label: "Active resistance — felt like swimming upstream", Value: -1.0},
		},
	},
	{
		SpectrumID: 2, ThemeName: "Loud and Quiet",
		Question: "Was the energy level of the activity a good fit for them today?",
		Options: []SensorOption{
			{Label: "Perfect match — they came alive", Value: 1.0},
			{Label: "Close enough — adjusted easily", Value: 0.5},
			{Label: "OK but could have gone either way", Value: 0.0},
			{Label: "Mismatch — too loud or too quiet for them", Value: -0.5},
			{Label: "Total mismatch — they checked out", Value: -1.0},
		},
	},
	{
		SpectrumID: 3, ThemeName: "Big Feelings, Small Feelings",
		Question: "Was the emotional intensity of the activity about right?",
		Options: []SensorOption{
			{Label: "Just right — engaged without overwhelm", Value: 1.0},
			{Label: "Mostly fine — one moment that was a bit much", Value: 0.5},
			{Label: "Hard to tell", Value: 0.0},
			{Label: "Too intense (or too flat) for them", Value: -0.5},
			{Label: "Way off — either overwhelmed or bored out", Value: -1.0},
		},
	},
	{
		SpectrumID: 4, ThemeName: "Yes and No",
		Question: "Did the boundary moments (saying yes or no) feel age-appropriate?",
		Options: []SensorOption{
			{Label: "They handled it confidently", Value: 1.0},
			{Label: "Managed with a bit of support", Value: 0.5},
			{Label: "Some moments OK, some not", Value: 0.0},
			{Label: "Most boundary moments were too hard", Value: -0.5},
			{Label: "Every 'no' was a battle", Value: -1.0},
		},
	},
	{
		SpectrumID: 5, ThemeName: "Mine and Yours",
		Question: "Was sharing and turn-taking the right level of challenge?",
		Options: []SensorOption{
			{Label: "Easy for them — sharing came naturally", Value: 1.0},
			{Label: "Manageable stretch — they grew a bit", Value: 0.5},
			{Label: "About 50/50 — some sharing, some friction", Value: 0.0},
			{Label: "Mostly too hard — not ready for this level", Value: -0.5},
			{Label: "Way beyond them — every share was a crisis", Value: -1.0},
		},
	},
	{
		SpectrumID: 6, ThemeName: "Try and Wait",
		Question: "Was the risk/uncertainty level right for them today?",
		Options: []SensorOption{
			{Label: "Thrived on it — wanted more challenge", Value: 1.0},
			{Label: "Good stretch without tipping over", Value: 0.5},
			{Label: "Neither easy nor too much", Value: 0.0},
			{Label: "A bit too much — needed a lot of scaffolding", Value: -0.5},
			{Label: "Way too risky / uncertain — paralysed", Value: -1.0},
		},
	},
	{
		SpectrumID: 7, ThemeName: "Same and Different",
		Question: "Was the amount of rule-changing/flexibility the right challenge?",
		Options: []SensorOption{
			{Label: "Loved it — adapted easily", Value: 1.0},
			{Label: "Slight discomfort but managed", Value: 0.5},
			{Label: "Could go either way", Value: 0.0},
			{Label: "The change was clearly too much", Value: -0.5},
			{Label: "Any change at all was unbearable", Value: -1.0},
		},
	},
	{
		SpectrumID: 8, ThemeName: "Your Feelings, My Feelings",
		Question: "Was the emotional exposure (other people's feelings) manageable?",
		Options: []SensorOption{
			{Label: "Noticed & moved on — good boundaries", Value: 1.0},
			{Label: "Got pulled in a bit but recovered", Value: 0.5},
			{Label: "Hard to tell how it affected them", Value: 0.0},
			{Label: "Absorbed too much — got weighed down", Value: -0.5},
			{Label: "Completely flooded by others' emotions", Value: -1.0},
		},
	},
	{
		SpectrumID: 9, ThemeName: "Tidy and Messy",
		Question: "Was the structure level (tidy/messy balance) comfortable?",
		Options: []SensorOption{
			{Label: "Perfect fit — neither too rigid nor too loose", Value: 1.0},
			{Label: "Slight preference for more/less structure, but fine", Value: 0.5},
			{Label: "Neither comfortable nor uncomfortable", Value: 0.0},
			{Label: "Clearly wanted more (or less) structure", Value: -0.5},
			{Label: "The structure level caused real distress", Value: -1.0},
		},
	},
	{
		SpectrumID: 10, ThemeName: "My Fault, Your Fault",
		Question: "Was the level of responsibility-taking expected in the activity right for them?",
		Options: []SensorOption{
			{Label: "Just right — owned their part naturally", Value: 1.0},
			{Label: "Manageable stretch — needed a small nudge", Value: 0.5},
			{Label: "Hard to tell", Value: 0.0},
			{Label: "Too much — deflected everything or absorbed everything", Value: -0.5},
			{Label: "Completely mismatched — blame-storm or guilt-spiral", Value: -1.0},
		},
	},
	{
		SpectrumID: 11, ThemeName: "Keeping and Letting Go",
		Question: "Was the letting-go / sharing-out demand of the activity the right level?",
		Options: []SensorOption{
			{Label: "Easy — they didn't mind giving things up", Value: 1.0},
			{Label: "A small stretch but they managed", Value: 0.5},
			{Label: "Neither easy nor hard", Value: 0.0},
			{Label: "Too much — they clung on or checked out", Value: -0.5},
			{Label: "Way too much — meltdown when anything was removed", Value: -1.0},
		},
	},
	{
		SpectrumID: 12, ThemeName: "Wanting and Waiting",
		Question: "Was the impulse-management demand of the activity right for them?",
		Options: []SensorOption{
			{Label: "Good fit — they channelled their energy well", Value: 1.0},
			{Label: "Manageable — a couple of moments but OK", Value: 0.5},
			{Label: "Hard to tell", Value: 0.0},
			{Label: "Too much waiting or restraint required", Value: -0.5},
			{Label: "Completely mismatched — couldn't contain themselves", Value: -1.0},
		},
	},
}

// AgencySensorQuestions — one per spectrum theme.
// WHAT IT MEASURES: Driver Moderation (is the ventral driver the sole power source,
// or is the child developing internal reward?).
// The parent observes WHERE the child looks for the "win."
var AgencySensorQuestions = []AgencySensorQuestion{
	{
		SpectrumID: 1, ThemeName: "Near and Far",
		Question: "When they did something good, who did they look at?",
		Options: []SensorOption{
			{Label: "Looked at what they'd made/done — satisfied with themselves", Value: 1.0},
			{Label: "Glanced at me then back to the task", Value: 0.5},
			{Label: "Hard to tell", Value: 0.0},
			{Label: "Looked at me first — needed my reaction", Value: -0.5},
			{Label: "Wouldn't continue until I praised them", Value: -1.0},
		},
	},
	{
		SpectrumID: 2, ThemeName: "Loud and Quiet",
		Question: "Was their engagement coming from inside (their own interest) or outside (your involvement)?",
		Options: []SensorOption{
			{Label: "Completely self-driven — forgot I was there", Value: 1.0},
			{Label: "Mostly internal — occasional check-in with me", Value: 0.5},
			{Label: "About half and half", Value: 0.0},
			{Label: "Mostly needed me to keep them going", Value: -0.5},
			{Label: "Stopped completely when I stepped back", Value: -1.0},
		},
	},
	{
		SpectrumID: 3, ThemeName: "Big Feelings, Small Feelings",
		Question: "After a big feeling, did they self-soothe or need you to fix it?",
		Options: []SensorOption{
			{Label: "Found their own way through — impressive", Value: 1.0},
			{Label: "Started self-soothing, then came to me for a top-up", Value: 0.5},
			{Label: "Needed me but calmed with minimal help", Value: 0.0},
			{Label: "Needed significant comfort from me", Value: -0.5},
			{Label: "Could not calm without me physically holding them", Value: -1.0},
		},
	},
	{
		SpectrumID: 4, ThemeName: "Yes and No",
		Question: "When they said 'no' (or 'yes'), did it feel like THEIR choice or a reaction to you?",
		Options: []SensorOption{
			{Label: "Clearly their own decision — considered and firm", Value: 1.0},
			{Label: "Mostly their own — some influence from me", Value: 0.5},
			{Label: "Unclear — could have gone either way", Value: 0.0},
			{Label: "Mostly reacting to what I wanted", Value: -0.5},
			{Label: "Just doing whatever I seemed to want (or the exact opposite)", Value: -1.0},
		},
	},
	{
		SpectrumID: 5, ThemeName: "Mine and Yours",
		Question: "When they shared, was it genuine or performed for you?",
		Options: []SensorOption{
			{Label: "Genuine — they wanted to share", Value: 1.0},
			{Label: "Mostly genuine — a bit of showing off", Value: 0.5},
			{Label: "Hard to tell", Value: 0.0},
			{Label: "Mostly for my benefit — checking I saw them share", Value: -0.5},
			{Label: "Only shared to get praise — took it back when I looked away", Value: -1.0},
		},
	},
	{
		SpectrumID: 6, ThemeName: "Try and Wait",
		Question: "When they tried something new, were they doing it for themselves or for your reaction?",
		Options: []SensorOption{
			{Label: "Purely for the thrill of trying — didn't check my face", Value: 1.0},
			{Label: "Mostly self-motivated — quick glance at me", Value: 0.5},
			{Label: "About equal — self-interest and wanting my approval", Value: 0.0},
			{Label: "Mainly to impress me", Value: -0.5},
			{Label: "Would not try unless I was watching and encouraging", Value: -1.0},
		},
	},
	{
		SpectrumID: 7, ThemeName: "Same and Different",
		Question: "When they solved a problem (new rule, different way), did they celebrate internally or need you to validate?",
		Options: []SensorOption{
			{Label: "Little fist-pump or smile to themselves — self-celebratory", Value: 1.0},
			{Label: "Pleased with themselves first, then showed me", Value: 0.5},
			{Label: "Mixed — some internal, some external", Value: 0.0},
			{Label: "Needed me to confirm it was right before feeling good", Value: -0.5},
			{Label: "Wouldn't accept it was 'right' unless I said so", Value: -1.0},
		},
	},
	{
		SpectrumID: 8, ThemeName: "Your Feelings, My Feelings",
		Question: "When they helped someone (or held back from helping), was it their own compass or echoing yours?",
		Options: []SensorOption{
			{Label: "Their own compass — genuine instinct", Value: 1.0},
			{Label: "Mostly their own — some mimicry of my behaviour", Value: 0.5},
			{Label: "Hard to separate their instinct from my influence", Value: 0.0},
			{Label: "Mainly copying what they thought I wanted", Value: -0.5},
			{Label: "Only helped (or held back) because I told them to", Value: -1.0},
		},
	},
	{
		SpectrumID: 9, ThemeName: "Tidy and Messy",
		Question: "When they tidied (or embraced mess), was it THEIR preference or compliance with yours?",
		Options: []SensorOption{
			{Label: "Clearly their own preference — did it without being asked", Value: 1.0},
			{Label: "Mostly their own — I nudged a little", Value: 0.5},
			{Label: "Half and half", Value: 0.0},
			{Label: "Mainly complying with what I wanted", Value: -0.5},
			{Label: "Only did it because I insisted — no internal drive", Value: -1.0},
		},
	},
	{
		SpectrumID: 10, ThemeName: "My Fault, Your Fault",
		Question: "When they took responsibility (or didn't), was it THEIR conscience or performing for you?",
		Options: []SensorOption{
			{Label: "Genuine ownership — said sorry/fixed it without prompting", Value: 1.0},
			{Label: "Mostly genuine — a bit of performance", Value: 0.5},
			{Label: "Hard to tell", Value: 0.0},
			{Label: "Mainly performing — said sorry because I was watching", Value: -0.5},
			{Label: "Only apologised under direct pressure — no internal compass", Value: -1.0},
		},
	},
	{
		SpectrumID: 11, ThemeName: "Keeping and Letting Go",
		Question: "When they let go of something (or held on), was it THEIR choice or yours?",
		Options: []SensorOption{
			{Label: "Their choice — genuine comfort with letting go", Value: 1.0},
			{Label: "Mostly theirs — I helped a little", Value: 0.5},
			{Label: "Hard to separate their choice from my influence", Value: 0.0},
			{Label: "Mainly doing what I wanted — no real choice", Value: -0.5},
			{Label: "Forced — they had no agency in the letting go", Value: -1.0},
		},
	},
	{
		SpectrumID: 12, ThemeName: "Wanting and Waiting",
		Question: "When they managed their impulse (or didn't), was it self-regulation or your control?",
		Options: []SensorOption{
			{Label: "Self-regulated — chose to wait on their own", Value: 1.0},
			{Label: "Mostly self-regulated — I reminded once", Value: 0.5},
			{Label: "About half and half", Value: 0.0},
			{Label: "Mainly my control — stopped only because I said so", Value: -0.5},
			{Label: "No self-regulation — only external limits worked", Value: -1.0},
		},
	},
}

// =============================================================================
// Range of Motion — The Zen Warrior Metric
// =============================================================================
//
// "Smooth running" is NOT quiet or obedient. It is RANGE OF MOTION.
//
// Low Range:  Child can ONLY function at their preferred pole.
//             (Only "Tidy," only "Alone," only "With praise")
// High Range: Child PREFERS one pole but maintains regulation at the other.
//             (Prefers "Tidy" but can handle "Messy" without meltdown)
//
// The Rosette visualisation shows 9 axes (one per spectrum). Each axis
// extends outward as Range of Motion increases. A Zen Warrior's rosette
// is a full circle. Most toddlers will have an asymmetric shape — and
// that's fine. Growth is measured by area increase, not symmetry.
// =============================================================================

// RangeOfMotion represents a child's regulatory capacity across all 9 spectra.
// Each spectrum has a score from 0.0 (locked to one pole) to 1.0 (full range).
type RangeOfMotion struct {
	ChildID   string    `json:"child_id"`
	UpdatedAt time.Time `json:"updated_at"`

	// Per-spectrum range scores. Key is spectrum_id (1-9).
	// 0.0 = can only function at preferred pole
	// 0.5 = can tolerate opposite pole with support
	// 1.0 = comfortable at both poles independently
	SpectrumRange map[int]float64 `json:"spectrum_range"`

	// Aggregate score — area of the rosette normalised to 0.0–1.0
	OverallRange float64 `json:"overall_range"`

	// Trend — is range growing, shrinking, or stable?
	// Computed from last 3 observations per spectrum.
	Trend string `json:"trend"` // "growing", "stable", "shrinking"
}

// =============================================================================
// Processing Logic
// =============================================================================

// GetSensorQuestionsForTheme returns the three sensor questions for a given
// spectrum theme. The parent app calls this after an activity to build
// the observation form.
func GetSensorQuestionsForTheme(spectrumID int) (recovery *RecoverySensorQuestion, friction *FrictionSensorQuestion, agency *AgencySensorQuestion) {
	for i := range RecoverySensorQuestions {
		if RecoverySensorQuestions[i].SpectrumID == spectrumID {
			recovery = &RecoverySensorQuestions[i]
			break
		}
	}
	for i := range FrictionSensorQuestions {
		if FrictionSensorQuestions[i].SpectrumID == spectrumID {
			friction = &FrictionSensorQuestions[i]
			break
		}
	}
	for i := range AgencySensorQuestions {
		if AgencySensorQuestions[i].SpectrumID == spectrumID {
			agency = &AgencySensorQuestions[i]
			break
		}
	}
	return
}

// ProcessObservation takes a raw ObservationEvent with parent responses
// and computes the derived fields: IsNoise, SpectrumInsights, DriverAdjustments,
// and PatchRecommendation.
func ProcessObservation(obs *ObservationEvent) {
	// === STEP 1: SKIP CHECK ===
	// "Bad day? Skip it." No load levels, no noise weighting.
	// Energy directionality already captures stress-reaction direction.
	if obs.SkippedDay {
		obs.IsNoise = true
		obs.SpectrumInsights = nil
		obs.DriverAdjustments = nil
		obs.PatchRecommendation = &ActivityPatch{
			PatchType:     "skip",
			ParentMessage: "You skipped today. That's fine — bad days are data about sleep, not about your child. Try again tomorrow.",
			PatchDuration: 0,
		}
		return
	}

	obs.IsNoise = false
	weight := 1.0

	// === STEP 2: DERIVE SPECTRUM INSIGHTS ===
	insights := make([]SpectrumInsight, 0, 3)

	// Recovery sensor → executive function insight (accelerator skill)
	if obs.Recovery.Value != 0 {
		insights = append(insights, SpectrumInsight{
			SpectrumID:   0, // executive_function is an accelerator skill, not a spectrum
			SpectrumName: "executive_function",
			Delta:        obs.Recovery.Value * 0.2 * weight,
			Reason:       recoverySummary(obs.Recovery.Value),
			Source:       "recovery",
		})
	}

	// Friction sensor → spectrum alignment insight
	// Negative friction means the activity was too far from natural pole.
	// This tells us WHERE the child sits on this spectrum.
	if obs.Friction.Value != 0 {
		// Look up which spectrum this activity primarily loads on
		activity := GetActivityByID(obs.ActivityID)
		if activity != nil {
			for _, sa := range activity.SpectrumActivations {
				if sa.Intensity == "primary" {
					insights = append(insights, SpectrumInsight{
						SpectrumID:   sa.SpectrumID,
						SpectrumName: sa.SpectrumName,
						Delta:        obs.Friction.Value * 0.15 * weight,
						Reason:       frictionSummary(obs.Friction.Value, sa.SpectrumName),
						Source:       "friction",
					})
					break
				}
			}
		}
	}

	// Agency sensor → driver moderation insight
	if obs.Agency.Value != 0 {
		insights = append(insights, SpectrumInsight{
			SpectrumID:   0,
			SpectrumName: "internal_locus",
			Delta:        obs.Agency.Value * 0.15 * weight,
			Reason:       agencySummary(obs.Agency.Value),
			Source:       "agency",
		})
	}

	obs.SpectrumInsights = insights

	// === STEP 3: DERIVE DRIVER ADJUSTMENTS ===
	adjustments := make([]DriverAdjustment, 0)

	// Low agency (looking to parent for validation) increases Impressor confidence
	if obs.Agency.Value < -0.3 {
		adjustments = append(adjustments, DriverAdjustment{
			DriverID:        "impressor",
			ConfidenceDelta: 0.1 * weight,
			Reason:          "Child sought external validation during activity",
		})
	}
	// High agency (self-driven) increases Explorer/Observer confidence
	if obs.Agency.Value > 0.5 {
		adjustments = append(adjustments, DriverAdjustment{
			DriverID:        "explorer",
			ConfidenceDelta: 0.05 * weight,
			Reason:          "Child was self-driven during activity",
		})
		adjustments = append(adjustments, DriverAdjustment{
			DriverID:        "observer",
			ConfidenceDelta: 0.05 * weight,
			Reason:          "Child was internally motivated",
		})
	}

	// High friction on orderliness → Sentinel signal
	if obs.Friction.Value < -0.5 {
		activity := GetActivityByID(obs.ActivityID)
		if activity != nil {
			for _, sa := range activity.SpectrumActivations {
				if sa.SpectrumName == "orderliness" && sa.Intensity == "primary" {
					adjustments = append(adjustments, DriverAdjustment{
						DriverID:        "sentinel",
						ConfidenceDelta: 0.1 * weight,
						Reason:          "High friction on structure change — safety-maintenance pattern",
					})
				}
			}
		}
	}

	// Very low recovery → may indicate Performer (high voltage, slow regulation)
	if obs.Recovery.Value < -0.5 {
		adjustments = append(adjustments, DriverAdjustment{
			DriverID:        "performer",
			ConfidenceDelta: 0.05 * weight,
			Reason:          "Slow recovery from high emotional voltage",
		})
	}

	obs.DriverAdjustments = adjustments

	// === STEP 4: GENERATE PATCH ===
	obs.PatchRecommendation = generatePatch(obs)
}

// generatePatch produces a next-session recommendation based on the sensor readings.
func generatePatch(obs *ObservationEvent) *ActivityPatch {
	// If all sensors are positive, no patch needed
	if obs.Recovery.Value >= 0.5 && obs.Friction.Value >= 0.0 && obs.Agency.Value >= 0.0 {
		return &ActivityPatch{
			PatchType:     "continue",
			ParentMessage: "That went well. Keep going with this level — they're in their growth zone.",
			PatchDuration: 0,
		}
	}

	// Priority: recovery problems first (safety), then friction (alignment), then agency (growth)
	if obs.Recovery.Value <= -0.5 {
		return &ActivityPatch{
			PatchType:     "increase_predictability",
			ParentMessage: "Recovery was hard today. Next time, try adding a countdown before transitions — '10 seconds until we stop.' Their brain needs a bridge between 'going' and 'stopped.'",
			PatchDuration: 3,
		}
	}

	if obs.Friction.Value <= -0.5 {
		return &ActivityPatch{
			PatchType:     "simplify",
			ParentMessage: "This was a big stretch from their comfort zone. Next time, try the same activity but dial it back — shorter duration, more scaffolding, or closer to their natural preference. Growth happens at the edge, not in the deep end.",
			PatchDuration: 2,
		}
	}

	if obs.Agency.Value <= -0.5 {
		return &ActivityPatch{
			PatchType:     "reduce_audience",
			ParentMessage: "They were doing it for you, not for themselves. Next time, try stepping back physically — sit further away, look at your phone for 10 seconds while they play. See if they keep going without your eyes on them. That's where internal reward grows.",
			PatchDuration: 3,
		}
	}

	// Mild signals — gentle suggestions
	if obs.Recovery.Value < 0 {
		return &ActivityPatch{
			PatchType:     "extend_duration",
			ParentMessage: "A bit of wobble in recovering — totally normal. Try the same activity again soon. Repetition builds the recovery pathway.",
			PatchDuration: 1,
		}
	}

	return nil
}

// =============================================================================
// Summary Helpers
// =============================================================================

func recoverySummary(value float64) string {
	switch {
	case value >= 0.5:
		return "Quick recovery — strong pilot strength"
	case value > 0:
		return "Adequate recovery — pilot strength developing"
	case value > -0.5:
		return "Slow recovery — pilot strength needs support"
	default:
		return "Recovery failure — activity may have exceeded current capacity"
	}
}

func frictionSummary(value float64, spectrumName string) string {
	switch {
	case value >= 0.5:
		return "Activity aligned well with natural " + spectrumName + " setting"
	case value > 0:
		return "Slight stretch on " + spectrumName + " — good growth zone"
	case value > -0.5:
		return "Moderate friction on " + spectrumName + " — at the edge of comfort"
	default:
		return "High friction on " + spectrumName + " — activity too far from natural pole"
	}
}

func agencySummary(value float64) string {
	switch {
	case value >= 0.5:
		return "Strong internal locus — self-motivated engagement"
	case value > 0:
		return "Developing internal locus — some self-direction"
	case value > -0.5:
		return "External locus dominant — primarily seeking parent validation"
	default:
		return "Fully dependent on external validation — ventral driver consuming executive function"
	}
}

// CalculateRangeOfMotion computes the rosette from a history of observations.
// It takes all non-noise observations for a child and produces per-spectrum
// range scores.
func CalculateRangeOfMotion(childID string, observations []ObservationEvent) RangeOfMotion {
	rom := RangeOfMotion{
		ChildID:       childID,
		UpdatedAt:     time.Now(),
		SpectrumRange: make(map[int]float64),
	}

	// Collect friction values per spectrum (friction is the best indicator
	// of range — high friction = low range, low friction = high range)
	spectrumFrictionSums := make(map[int]float64)
	spectrumFrictionCounts := make(map[int]int)

	for _, obs := range observations {
		if obs.IsNoise {
			continue
		}
		// Find which spectrum the activity primarily loads on
		activity := GetActivityByID(obs.ActivityID)
		if activity == nil {
			continue
		}
		for _, sa := range activity.SpectrumActivations {
			if sa.Intensity == "primary" {
				spectrumFrictionSums[sa.SpectrumID] += obs.Friction.Value
				spectrumFrictionCounts[sa.SpectrumID]++
				break
			}
		}
	}

	// Convert friction averages to range scores
	// Friction -1.0 → range 0.0 (locked)
	// Friction  0.0 → range 0.5 (developing)
	// Friction +1.0 → range 1.0 (full range)
	totalRange := 0.0
	counted := 0
	for specID := 1; specID <= 12; specID++ {
		if count, ok := spectrumFrictionCounts[specID]; ok && count > 0 {
			avg := spectrumFrictionSums[specID] / float64(count)
			rangeScore := (avg + 1.0) / 2.0 // map [-1,1] → [0,1]
			if rangeScore < 0 {
				rangeScore = 0
			}
			if rangeScore > 1 {
				rangeScore = 1
			}
			rom.SpectrumRange[specID] = rangeScore
			totalRange += rangeScore
			counted++
		}
	}

	if counted > 0 {
		rom.OverallRange = totalRange / float64(counted)
	}

	// Trend: compare recent observations (last 3) vs earlier
	// Simplified for now — will use proper windowing later
	rom.Trend = "stable"

	return rom
}
