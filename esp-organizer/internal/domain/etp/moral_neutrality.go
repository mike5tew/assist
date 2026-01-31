package etp

// MoralNeutralityPrinciple represents the core philosophical stance
// ETPs are not moral categories - they are positioning on biological spectra
// Like a physicist studying forces, not judging them

// SliderState represents the state of an ETP slider
type SliderState string

const (
	SliderFluid    SliderState = "fluid"    // Can move freely (healthy)
	SliderSticky   SliderState = "sticky"   // Moves with difficulty
	SliderStuck    SliderState = "stuck"    // Cannot move (dysfunction)
	SliderHijacked SliderState = "hijacked" // Moves involuntarily via triggers
)

// ETPSlider represents a single spectrum as a movable slider
type ETPSlider struct {
	SpectrumID      int         `json:"spectrum_id"`
	SpectrumName    string      `json:"spectrum_name"`
	CurrentPosition float64     `json:"current_position"` // -2 to +2
	RestingPosition float64     `json:"resting_position"` // Where it naturally returns
	State           SliderState `json:"state"`            // fluid, sticky, stuck, hijacked
	RangeOfMotion   float64     `json:"range_of_motion"`  // 0-4 (full range is 4)
	TriggerPoints   []float64   `json:"trigger_points"`   // Positions that cause involuntary movement
}

// AnimalPilotModel represents the core HumanOS architecture
type AnimalPilotModel struct {
	Description string `json:"description"`
	Animal      string `json:"animal"` // The hardware - the sliders themselves
	Pilot       string `json:"pilot"`  // Executive function - the hand on the sliders
	Dysfunction string `json:"dysfunction"`
	Goal        string `json:"goal"`
}

// TheAnimalPilotModel defines the core insight
var TheAnimalPilotModel = AnimalPilotModel{
	Description: "The 'Wild Animal' isn't a glitch in the OS; it IS the OS",
	Animal:      "The hardware - the 17 ETP sliders themselves. Biological, amoral, functional.",
	Pilot:       "Executive Function - the hand on the slider. The capacity to CHOOSE position.",
	Dysfunction: "Not having a slider at -2 (aggression). Dysfunction is having it STUCK there.",
	Goal:        "Agency to move sliders yourself, rather than having triggers move them for you.",
}

// KineticAggressionProfile represents the ETP state during violence
type KineticAggressionProfile struct {
	Description string             `json:"description"`
	Settings    map[string]float64 `json:"settings"` // Which sliders, what positions
	Trigger     string             `json:"trigger"`
	Duration    string             `json:"duration"`
	Recovery    string             `json:"recovery"`
}

// KineticAggression defines the ETP profile during violent state
var KineticAggression = KineticAggressionProfile{
	Description: "Violence is what happens when multiple ETP dials are turned to one extreme simultaneously",
	Settings: map[string]float64{
		"impulse_gap":          -2.0, // Maximum Impulsivity / Zero Pause
		"social_gravity":       -2.0, // Total Independence / Disregard for Social Contract
		"authority_response":   +2.0, // Absolute Challenge to Authority/Rules
		"risk_tolerance":       +2.0, // Maximum Risk Seeking
		"presence_sensitivity": +2.0, // Fully Present / Attuned to threat
	},
	Trigger:  "Perceived threat to resource, status, or safety",
	Duration: "Until threat resolved or energy depleted",
	Recovery: "Sliders return toward resting position as threat recedes",
}

// SliderDashboard represents a person's full ETP dashboard
type SliderDashboard struct {
	StudentID       string             `json:"student_id"`
	Sliders         map[int]*ETPSlider `json:"sliders"`
	OverallFluidity float64            `json:"overall_fluidity"` // 0-1, how fluid all sliders are
	StuckSliders    []int              `json:"stuck_sliders"`    // Which sliders are stuck
	HijackedSliders []int              `json:"hijacked_sliders"` // Which respond to triggers
	PilotStrength   float64            `json:"pilot_strength"`   // Executive function capacity
}

// CreateDashboard initializes a slider dashboard from ETP profile
func CreateDashboard(profile *ETPProfile) *SliderDashboard {
	sliders := make(map[int]*ETPSlider)

	for spectrumID, setting := range profile.Settings {
		spectrum := GetSpectrumByID(spectrumID)
		if spectrum == nil {
			continue
		}

		sliders[spectrumID] = &ETPSlider{
			SpectrumID:      spectrumID,
			SpectrumName:    spectrum.Name,
			CurrentPosition: setting,
			RestingPosition: setting,
			State:           SliderFluid, // Assume fluid until assessed
			RangeOfMotion:   4.0,         // Assume full range until assessed
			TriggerPoints:   []float64{}, // No known triggers until assessed
		}
	}

	return &SliderDashboard{
		StudentID:       profile.StudentID,
		Sliders:         sliders,
		OverallFluidity: 1.0, // Assume fluid until assessed
		StuckSliders:    []int{},
		HijackedSliders: []int{},
		PilotStrength:   0.5, // Default middle
	}
}

// RangeOfMotionSkill represents the ability to access different parts of a spectrum
type RangeOfMotionSkill struct {
	SpectrumID      int     `json:"spectrum_id"`
	SpectrumName    string  `json:"spectrum_name"`
	CanAccessNeg2   bool    `json:"can_access_neg_2"` // Can reach -2 when appropriate
	CanAccessNeg1   bool    `json:"can_access_neg_1"`
	CanAccessCenter bool    `json:"can_access_center"`
	CanAccessPos1   bool    `json:"can_access_pos_1"`
	CanAccessPos2   bool    `json:"can_access_pos_2"` // Can reach +2 when appropriate
	RangeScore      float64 `json:"range_score"`      // 0-1 (1 = full range)
	FluidityScore   float64 `json:"fluidity_score"`   // How smoothly they can move
}

// AssessRangeOfMotion evaluates range of motion for a spectrum
func AssessRangeOfMotion(slider *ETPSlider, observedPositions []float64) *RangeOfMotionSkill {
	canNeg2, canNeg1, canCenter, canPos1, canPos2 := false, false, false, false, false

	for _, pos := range observedPositions {
		if pos <= -1.5 {
			canNeg2 = true
		}
		if pos <= -0.5 && pos > -1.5 {
			canNeg1 = true
		}
		if pos > -0.5 && pos < 0.5 {
			canCenter = true
		}
		if pos >= 0.5 && pos < 1.5 {
			canPos1 = true
		}
		if pos >= 1.5 {
			canPos2 = true
		}
	}

	accessCount := 0
	if canNeg2 {
		accessCount++
	}
	if canNeg1 {
		accessCount++
	}
	if canCenter {
		accessCount++
	}
	if canPos1 {
		accessCount++
	}
	if canPos2 {
		accessCount++
	}

	return &RangeOfMotionSkill{
		SpectrumID:      slider.SpectrumID,
		SpectrumName:    slider.SpectrumName,
		CanAccessNeg2:   canNeg2,
		CanAccessNeg1:   canNeg1,
		CanAccessCenter: canCenter,
		CanAccessPos1:   canPos1,
		CanAccessPos2:   canPos2,
		RangeScore:      float64(accessCount) / 5.0,
		FluidityScore:   0.5, // Would need more data
	}
}

// SystemLevelShift represents a deliberate multi-slider movement
type SystemLevelShift struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	FromState   map[int]float64 `json:"from_state"` // Starting positions
	ToState     map[int]float64 `json:"to_state"`   // Target positions
	Trigger     string          `json:"trigger"`    // What initiates the shift
	Purpose     string          `json:"purpose"`    // Why this shift is useful
	Recovery    string          `json:"recovery"`   // How to return to baseline
}

// SystemLevelShifts defines common multi-slider patterns
var SystemLevelShifts = []SystemLevelShift{
	{
		Name:        "defensive_activation",
		Description: "Shifting to protective mode when threat detected",
		FromState: map[int]float64{
			8:  0.0, // impulse_gap: balanced
			10: 0.0, // risk_tolerance: balanced
			12: 0.0, // presence_sensitivity: balanced
		},
		ToState: map[int]float64{
			8:  -1.5, // impulse_gap: more reactive
			10: +1.0, // risk_tolerance: more risk-accepting (to fight)
			12: +2.0, // presence_sensitivity: fully attuned
		},
		Trigger:  "Genuine threat to self or others",
		Purpose:  "Survival - matching energy to threat level",
		Recovery: "Controlled de-escalation as threat reduces",
	},
	{
		Name:        "deep_focus",
		Description: "Shifting to concentrated work mode",
		FromState: map[int]float64{
			1:  0.0, // social_gravity: balanced
			4:  0.0, // energy_directionality: balanced
			12: 0.0, // presence_sensitivity: balanced
		},
		ToState: map[int]float64{
			1:  -1.5, // social_gravity: more independent
			4:  -1.5, // energy_directionality: inward
			12: +1.5, // presence_sensitivity: highly attuned
		},
		Trigger:  "Complex cognitive task requiring concentration",
		Purpose:  "Blocking distraction, maximizing cognitive resources",
		Recovery: "Natural return to baseline as task completes",
	},
	{
		Name:        "social_engagement",
		Description: "Shifting to connection mode for group interaction",
		FromState: map[int]float64{
			1: 0.0, // social_gravity: balanced
			3: 0.0, // emotional_transparency: balanced
			5: 0.0, // mirror_neuron_tuning: balanced
		},
		ToState: map[int]float64{
			1: +1.5, // social_gravity: more cohesive
			3: +1.0, // emotional_transparency: more transparent
			5: +1.0, // mirror_neuron_tuning: more absorbent
		},
		Trigger:  "Social situation requiring connection",
		Purpose:  "Building relationship, reading group dynamics",
		Recovery: "Recharge time afterward (especially for introverts)",
	},
	{
		Name:        "principled_stand",
		Description: "Shifting to hold position against pressure",
		FromState: map[int]float64{
			14: 0.0, // authority_response: balanced
			17: 0.0, // integrity_logic: balanced
			13: 0.0, // agency_threshold: balanced
		},
		ToState: map[int]float64{
			14: +1.5, // authority_response: more challenging
			17: +1.5, // integrity_logic: more absolutist
			13: +1.5, // agency_threshold: more active
		},
		Trigger:  "Core values threatened, boundary violated",
		Purpose:  "Maintaining integrity despite social pressure",
		Recovery: "Return to flexibility once stand is acknowledged",
	},
}

// TurningOtherCheekPower represents the ultimate power demonstration
type TurningOtherCheekPower struct {
	Description   string `json:"description"`
	Message       string `json:"message"`
	Requirement   string `json:"requirement"`
	Demonstration string `json:"demonstration"`
}

// TheTurningCheekPower captures the power dynamic of restraint
var TheTurningCheekPower = TurningOtherCheekPower{
	Description:   "Showing the aggressor you could match their energy but CHOOSE not to",
	Message:       "I have the same sliders you do. I could move mine to match yours in a heartbeat. But I am so much more powerful than you that I am choosing to keep mine at 'Cooperative' while you scream.",
	Requirement:   "Must actually HAVE the capacity to match aggression - not weakness disguised as virtue",
	Demonstration: "The visible pause where you could escalate but don't. The calm that comes from strength, not fear.",
}

// PeriodicTableAnalogy represents ETPs as elements of human nature
type PeriodicTableAnalogy struct {
	Principle   string `json:"principle"`
	Example     string `json:"example"`
	Lesson      string `json:"lesson"`
	Application string `json:"application"`
}

// ThePeriodicTableOfHumanNature captures the chemistry analogy
var ThePeriodicTableOfHumanNature = PeriodicTableAnalogy{
	Principle:   "The 17 ETPs are the periodic table of human nature",
	Example:     "You don't get angry at Hydrogen for being explosive; you just learn how to handle it so it doesn't blow up the lab",
	Lesson:      "By treating the 'Wild Animal' as spectrum settings, you remove shame and judgment",
	Application: "Give the child the Schematic to their own brain - owner's manual, not moral rulebook",
}

// StuckSliderIndicator represents signs that a slider is stuck
type StuckSliderIndicator struct {
	Indicator    string `json:"indicator"`
	Implication  string `json:"implication"`
	Intervention string `json:"intervention"`
}

// StuckSliderIndicators help identify when sliders aren't moving
var StuckSliderIndicators = []StuckSliderIndicator{
	{
		Indicator:    "Same response regardless of context",
		Implication:  "Slider not adjusting to situation appropriately",
		Intervention: "Gentle exposure to contexts requiring different positions",
	},
	{
		Indicator:    "High distress when situation requires different position",
		Implication:  "Slider movement causes voltage spike",
		Intervention: "Gradual desensitization to other positions",
	},
	{
		Indicator:    "Can't imagine self at different position",
		Implication:  "Identity fused with current position",
		Intervention: "Role-play exercises accessing other positions safely",
	},
	{
		Indicator:    "Judges others at different positions as 'wrong'",
		Implication:  "Position moralized rather than seen as tool",
		Intervention: "Teach functional value of all positions",
	},
	{
		Indicator:    "Physical/emotional shutdown when pushed to move",
		Implication:  "Trauma or deep conditioning at current position",
		Intervention: "Professional support, very gradual, trauma-informed approach",
	},
}

// RangeOfMotionExercise represents a practice for developing slider fluidity
type RangeOfMotionExercise struct {
	Name         string `json:"name"`
	SpectrumID   int    `json:"spectrum_id"`
	SpectrumName string `json:"spectrum_name"`
	Exercise     string `json:"exercise"`
	AtNeg2       string `json:"at_neg_2"` // What it looks like at -2
	AtPos2       string `json:"at_pos_2"` // What it looks like at +2
	Debrief      string `json:"debrief"`
}

// RangeOfMotionExercises for developing fluidity
var RangeOfMotionExercises = []RangeOfMotionExercise{
	{
		Name:         "challenging_compliant_switch",
		SpectrumID:   14,
		SpectrumName: "authority_response",
		Exercise:     "Game where you must be Challenging (+2) then circle time where you must be Compliant (-2)",
		AtNeg2:       "Following all instructions precisely, supporting the leader, not questioning",
		AtPos2:       "Questioning everything, proposing alternatives, testing boundaries",
		Debrief:      "What did each feel like? When is each useful?",
	},
	{
		Name:         "independent_cohesive_switch",
		SpectrumID:   1,
		SpectrumName: "social_gravity",
		Exercise:     "Solo work time (Independent) then group collaboration (Cohesive)",
		AtNeg2:       "Working entirely alone, not seeking input, self-reliant",
		AtPos2:       "Constantly checking in, sharing everything, group decisions",
		Debrief:      "When did you feel more productive? More connected?",
	},
	{
		Name:         "reactive_reflective_switch",
		SpectrumID:   8,
		SpectrumName: "impulse_gap",
		Exercise:     "Fast-paced game requiring instant reactions, then slow discussion requiring pauses",
		AtNeg2:       "Immediate responses, no filtering, first thought out loud",
		AtPos2:       "Long pauses, careful consideration, edited responses",
		Debrief:      "What mistakes did you make at each? What did you catch?",
	},
	{
		Name:         "risk_averse_seeking_switch",
		SpectrumID:   10,
		SpectrumName: "risk_tolerance",
		Exercise:     "Safe, known activity then something with managed uncertainty",
		AtNeg2:       "Choosing the safest option, avoiding uncertainty, planning extensively",
		AtPos2:       "Jumping into unknown, excited by uncertainty, minimal planning",
		Debrief:      "What opportunities did each create? What did each cost?",
	},
}

// StickerBookIntegration represents how this maps to primary school curriculum
type StickerBookIntegration struct {
	YearGroup   int    `json:"year_group"`
	SkillName   string `json:"skill_name"`
	Concept     string `json:"concept"`
	Exercise    string `json:"exercise"`
	StickerName string `json:"sticker_name"`
}

// StickerBookSkills for primary school range of motion curriculum
var StickerBookSkills = []StickerBookIntegration{
	{
		YearGroup:   3,
		SkillName:   "Slider Awareness",
		Concept:     "I have different settings for different situations",
		Exercise:    "Name which 'setting' you're on right now",
		StickerName: "Setting Spotter",
	},
	{
		YearGroup:   4,
		SkillName:   "Slider Movement",
		Concept:     "I can change my setting when I need to",
		Exercise:    "Deliberately change one slider for an activity",
		StickerName: "Dial Turner",
	},
	{
		YearGroup:   5,
		SkillName:   "Situation Matching",
		Concept:     "Different situations need different settings",
		Exercise:    "Match the situation to the helpful setting",
		StickerName: "Setting Matcher",
	},
	{
		YearGroup:   6,
		SkillName:   "Range of Motion",
		Concept:     "Can I access all parts of my spectra?",
		Exercise:    "Practice being at +2 and -2 on the same spectrum",
		StickerName: "Master of the Dial",
	},
}
