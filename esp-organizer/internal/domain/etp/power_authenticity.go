package etp

// ProtectionMode represents stages in the protection response
type ProtectionMode string

const (
	ProtectionHealthy   ProtectionMode = "healthy"   // Appropriate boundaries, assertiveness
	ProtectionWarped    ProtectionMode = "warped"    // Aggression as deterrent
	ProtectionCorrupted ProtectionMode = "corrupted" // Cruelty as "strength signaling"
	ProtectionPerverted ProtectionMode = "perverted" // Cruelty as status/tribal membership
)

// ProtectionStage represents a stage in the protection→cruelty pipeline
type ProtectionStage struct {
	Mode         ProtectionMode `json:"mode"`
	Description  string         `json:"description"`
	Behavior     string         `json:"behavior"`
	Underlying   string         `json:"underlying"` // What's really happening
	Indicators   []string       `json:"indicators"`
	Intervention string         `json:"intervention"`
}

// ProtectionPipeline shows the progression from healthy to corrupted
var ProtectionPipeline = []ProtectionStage{
	{
		Mode:         ProtectionHealthy,
		Description:  "Appropriate protective response",
		Behavior:     "Setting boundaries, assertive communication",
		Underlying:   "Genuine protection need met proportionally",
		Indicators:   []string{"Proportional response", "De-escalates when threat gone", "No audience needed"},
		Intervention: "Support and reinforce - this is healthy",
	},
	{
		Mode:         ProtectionWarped,
		Description:  "Aggression as deterrent",
		Behavior:     "Excessive force, preemptive aggression",
		Underlying:   "Fear driving over-response",
		Indicators:   []string{"Disproportionate response", "Stays aggressive after threat gone", "Hypervigilance"},
		Intervention: "Address underlying fear. 'What are you really afraid of?'",
	},
	{
		Mode:         ProtectionCorrupted,
		Description:  "Cruelty as 'strength signaling'",
		Behavior:     "Cruelty to demonstrate power",
		Underlying:   "Deep insecurity needing external validation",
		Indicators:   []string{"Needs witnesses", "Cruelty beyond protection need", "Enjoys target's fear"},
		Intervention: "Name the insecurity gently. 'You don't need them afraid to be safe.'",
	},
	{
		Mode:         ProtectionPerverted,
		Description:  "Cruelty as status symbol",
		Behavior:     "Cruelty for tribal membership and status",
		Underlying:   "Belonging need hijacked by toxic group",
		Indicators:   []string{"Performed for group", "Applauded by peers", "Escalating cruelty for status"},
		Intervention: "Offer alternative belonging. Remove status rewards for cruelty.",
	},
}

// PowerAuthenticity represents the distinction between real and disguised power
type PowerAuthenticity string

const (
	PowerAuthentic PowerAuthenticity = "authentic" // Real power through competence
	PowerDisguised PowerAuthenticity = "disguised" // Insecurity masquerading as strength
)

// PowerAssessment represents an assessment of power authenticity
type PowerAssessment struct {
	Authenticity    PowerAuthenticity `json:"authenticity"`
	Indicators      []string          `json:"indicators"`
	UnderlyingState string            `json:"underlying_state"`
	RealPowerLevel  float64           `json:"real_power_level"` // 0-1
	InsecurityLevel float64           `json:"insecurity_level"` // 0-1
	Recommendation  string            `json:"recommendation"`
}

// AuthenticPowerIndicators - signs of real power
var AuthenticPowerIndicators = []string{
	"Doesn't need to prove it",
	"Can be kind without feeling weak",
	"Secure enough to admit mistakes",
	"Power comes from competence, not fear",
	"Can protect without cruelty",
	"Others follow voluntarily, not from fear",
	"Can hold complexity without collapsing",
	"Doesn't require others' submission to feel safe",
}

// DisguisedPowerIndicators - signs of insecurity masquerading as strength
var DisguisedPowerIndicators = []string{
	"Needs constant proof/performance",
	"Kindness feels threatening",
	"Mistakes are intolerable",
	"Power comes from others' fear",
	"Cruelty is 'necessary' for protection",
	"Others comply from fear, not respect",
	"Collapses complexity into us/them",
	"Needs others diminished to feel elevated",
}

// AssessPowerAuthenticity evaluates whether displayed power is authentic
func AssessPowerAuthenticity(behaviors []string) *PowerAssessment {
	authenticScore := 0.0
	disguisedScore := 0.0

	authenticMatches := []string{}
	disguisedMatches := []string{}

	for _, behavior := range behaviors {
		for _, indicator := range AuthenticPowerIndicators {
			if containsPattern(behavior, indicator) {
				authenticScore += 1
				authenticMatches = append(authenticMatches, behavior)
			}
		}
		for _, indicator := range DisguisedPowerIndicators {
			if containsPattern(behavior, indicator) {
				disguisedScore += 1
				disguisedMatches = append(disguisedMatches, behavior)
			}
		}
	}

	total := authenticScore + disguisedScore
	if total == 0 {
		total = 1 // Avoid division by zero
	}

	var authenticity PowerAuthenticity
	var underlying string
	var recommendation string

	if authenticScore > disguisedScore {
		authenticity = PowerAuthentic
		underlying = "Genuine competence and security"
		recommendation = "This is real power. Reinforce and develop."
	} else {
		authenticity = PowerDisguised
		underlying = "Insecurity requiring external validation"
		recommendation = "Address underlying insecurity. Offer real competence building."
	}

	return &PowerAssessment{
		Authenticity:    authenticity,
		Indicators:      append(authenticMatches, disguisedMatches...),
		UnderlyingState: underlying,
		RealPowerLevel:  authenticScore / total,
		InsecurityLevel: disguisedScore / total,
		Recommendation:  recommendation,
	}
}

// RealPowerSkill represents skills that constitute authentic power
type RealPowerSkill struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	CountersCruelty bool     `json:"counters_cruelty"` // Does this replace cruelty-based "power"?
	Practice        []string `json:"practice"`
	ETPSpectra      []int    `json:"etp_spectra"` // Which spectra this develops
}

// RealPowerSkills defines the authentic power curriculum
var RealPowerSkills = []RealPowerSkill{
	{
		Name:            "boundary_without_cruelty",
		Description:     "Setting clear limits without aggression",
		CountersCruelty: true,
		Practice:        []string{"'I' statements", "Calm repetition", "Walking away as power"},
		ETPSpectra:      []int{14, 13}, // authority_response, agency_threshold
	},
	{
		Name:            "de_escalation",
		Description:     "Reducing others' fear without capitulation",
		CountersCruelty: true,
		Practice:        []string{"Lowering voice/energy", "Finding shared ground", "Naming dynamics"},
		ETPSpectra:      []int{7, 8}, // voltage_sensitivity, impulse_gap
	},
	{
		Name:            "complex_problem_solving",
		Description:     "Addressing underlying needs rather than symptoms",
		CountersCruelty: true,
		Practice:        []string{"Root cause analysis", "Win-win design", "Systems thinking"},
		ETPSpectra:      []int{15, 17}, // ambiguity_tolerance, integrity_logic
	},
	{
		Name:            "earned_respect",
		Description:     "Being followed because effective, not feared",
		CountersCruelty: true,
		Practice:        []string{"Competence demonstration", "Keeping commitments", "Helping others succeed"},
		ETPSpectra:      []int{16, 6}, // status_sensitivity, resource_allocation
	},
	{
		Name:            "vulnerability_as_strength",
		Description:     "Admitting mistakes and uncertainty without losing status",
		CountersCruelty: true,
		Practice:        []string{"Public learning", "Asking for help", "Crediting others"},
		ETPSpectra:      []int{3, 2}, // emotional_transparency, guilt_response
	},
	{
		Name:            "protective_presence",
		Description:     "Creating safety for others without aggression",
		CountersCruelty: true,
		Practice:        []string{"Standing with vulnerable", "Calm confidence", "Predictable reliability"},
		ETPSpectra:      []int{1, 12}, // social_gravity, presence_sensitivity
	},
}

// ProtectionReframe helps reframe cruelty back to healthy protection
type ProtectionReframe struct {
	CrueltyBehavior    string `json:"cruelty_behavior"`
	UnderlyingNeed     string `json:"underlying_need"`
	HealthyAlternative string `json:"healthy_alternative"`
	Script             string `json:"script"` // What to say
}

// ProtectionReframes provides reframes for common cruelty patterns
var ProtectionReframes = []ProtectionReframe{
	{
		CrueltyBehavior:    "Mocking someone's mistake",
		UnderlyingNeed:     "Status/competence validation",
		HealthyAlternative: "Demonstrate competence through helping",
		Script:             "I see you noticed the mistake. Want to help them fix it? That's harder but more impressive.",
	},
	{
		CrueltyBehavior:    "Excluding someone from group",
		UnderlyingNeed:     "In-group security",
		HealthyAlternative: "Build in-group through shared challenge, not shared enemy",
		Script:             "You're protecting your group. What if you got stronger by including different perspectives?",
	},
	{
		CrueltyBehavior:    "Physical intimidation",
		UnderlyingNeed:     "Physical safety/territory",
		HealthyAlternative: "Boundaries through presence, not threat",
		Script:             "You're showing you're strong. Real strength doesn't need them scared - it just IS strong.",
	},
	{
		CrueltyBehavior:    "Humiliating someone publicly",
		UnderlyingNeed:     "Status elevation",
		HealthyAlternative: "Status through demonstrated skill",
		Script:             "That got a reaction. But respect from fear doesn't last. What could earn real respect?",
	},
	{
		CrueltyBehavior:    "Spreading rumors/gossip",
		UnderlyingNeed:     "Social power/belonging",
		HealthyAlternative: "Build alliances through shared value, not shared target",
		Script:             "You're building connections. What if you connected over something you're building, not destroying?",
	},
}

// CrueltyToStrengthTransition represents the path from cruelty to authentic strength
type CrueltyToStrengthTransition struct {
	Stage       int    `json:"stage"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Signs       string `json:"signs"`
	Support     string `json:"support"`
}

// TransitionPath from cruelty to authentic power
var CrueltyToStrengthPath = []CrueltyToStrengthTransition{
	{
		Stage:       1,
		Name:        "awareness",
		Description: "Recognizing cruelty as weakness disguised",
		Signs:       "Moments of shame after cruelty, wondering 'why did I do that?'",
		Support:     "Name the pattern without shaming. 'That's insecurity talking, not strength.'",
	},
	{
		Stage:       2,
		Name:        "pause",
		Description: "Creating space between trigger and cruelty",
		Signs:       "Sometimes stopping before being cruel, catching self",
		Support:     "Celebrate pauses. 'You stopped yourself. That's the hard part.'",
	},
	{
		Stage:       3,
		Name:        "alternative",
		Description: "Trying different responses to same triggers",
		Signs:       "Experimenting with boundary-without-cruelty",
		Support:     "Provide specific alternatives. 'Instead of that, try this.'",
	},
	{
		Stage:       4,
		Name:        "competence",
		Description: "Building real skills that provide authentic power",
		Signs:       "Finding status through competence, not cruelty",
		Support:     "Create opportunities to demonstrate real capability.",
	},
	{
		Stage:       5,
		Name:        "integration",
		Description: "Cruelty no longer needed because security is internal",
		Signs:       "Can be challenged without cruelty response",
		Support:     "Reinforce identity as 'someone who handles things well'",
	},
}
