package etp

// OperatingMode represents the three cognitive operating modes
type OperatingMode string

const (
	ModeDialectical OperatingMode = "dialectical" // Engage complexity, seek understanding
	ModeBoundary    OperatingMode = "boundary"    // Protect self/others, maintain integrity
	ModeDisengage   OperatingMode = "disengage"   // Withdraw, conserve energy, fight later
)

// ModeCharacteristics describes an operating mode
type ModeCharacteristics struct {
	Mode          OperatingMode `json:"mode"`
	When          string        `json:"when"`
	Goal          string        `json:"goal"`
	EnergyCost    string        `json:"energy_cost"`
	VoltageDemand string        `json:"voltage_demand"`
	Strengths     []string      `json:"strengths"`
	Risks         []string      `json:"risks"`
}

// OperatingModes defines all three modes
var OperatingModes = map[OperatingMode]ModeCharacteristics{
	ModeDialectical: {
		Mode:          ModeDialectical,
		When:          "Both sides operate in good faith",
		Goal:          "Synthesis, understanding, integration",
		EnergyCost:    "High voltage cost, high potential gain",
		VoltageDemand: "Requires cognitive surplus",
		Strengths:     []string{"Creates lasting solutions", "Builds bridges", "Expands understanding", "Develops wisdom"},
		Risks:         []string{"Exhausting", "Can be exploited by bad-faith actors", "May appear weak"},
	},
	ModeBoundary: {
		Mode:          ModeBoundary,
		When:          "One side crosses ethical lines or operates in bad faith",
		Goal:          "Protection, integrity maintenance",
		EnergyCost:    "Defensive, necessary preservation",
		VoltageDemand: "Moderate - focused energy",
		Strengths:     []string{"Protects vulnerable", "Maintains integrity", "Clear limits", "Stops exploitation"},
		Risks:         []string{"Can become rigid", "May miss synthesis opportunities", "Can escalate conflict"},
	},
	ModeDisengage: {
		Mode:          ModeDisengage,
		When:          "Engagement fuels harmful systems or energy is depleted",
		Goal:          "System change through non-participation, recovery",
		EnergyCost:    "Withdrawal to fight another day",
		VoltageDemand: "Low - conservation mode",
		Strengths:     []string{"Preserves energy", "Denies fuel to toxic systems", "Strategic patience", "Self-care"},
		Risks:         []string{"Can become avoidance", "May miss opportunities", "Can be misread as surrender"},
	},
}

// CircuitBreaker represents conditions that should trigger mode switching
type CircuitBreaker struct {
	Condition   string        `json:"condition"`
	Indicator   string        `json:"indicator"`
	TriggerMode OperatingMode `json:"trigger_mode"`
	Response    string        `json:"response"`
}

// CircuitBreakers define when to switch OUT of dialectical mode
var CircuitBreakers = []CircuitBreaker{
	{
		Condition:   "power_imbalance",
		Indicator:   "One side has significantly more power and is using it coercively",
		TriggerMode: ModeBoundary,
		Response:    "Protect vulnerable first, dialectics later",
	},
	{
		Condition:   "dehumanization",
		Indicator:   "One side denies others' humanity or basic dignity",
		TriggerMode: ModeBoundary,
		Response:    "Boundary defense, not dialectical engagement",
	},
	{
		Condition:   "bad_faith",
		Indicator:   "One side is using dialectics to manipulate, not understand",
		TriggerMode: ModeBoundary,
		Response:    "Disengage from false dialogue",
	},
	{
		Condition:   "exploitation",
		Indicator:   "Dialectics being used to maintain oppression or extract value",
		TriggerMode: ModeDisengage,
		Response:    "Resistance through non-participation",
	},
	{
		Condition:   "voltage_depletion",
		Indicator:   "Personal energy too low to maintain dialectical complexity",
		TriggerMode: ModeDisengage,
		Response:    "Withdraw, recover, return stronger",
	},
	{
		Condition:   "cruelty_as_status",
		Indicator:   "Cruelty is being rewarded with social status in the group",
		TriggerMode: ModeDisengage,
		Response:    "Remove yourself as audience; cruelty needs witnesses",
	},
}

// ModeRecommendation represents a recommendation for which mode to use
type ModeRecommendation struct {
	RecommendedMode   OperatingMode    `json:"recommended_mode"`
	TriggeredBreakers []CircuitBreaker `json:"triggered_breakers"`
	EnergyLevel       float64          `json:"energy_level"`     // 0-1
	GoodFaithLevel    float64          `json:"good_faith_level"` // 0-1
	PowerBalance      float64          `json:"power_balance"`    // -1 to +1 (0 = balanced)
	Rationale         string           `json:"rationale"`
	Actions           []string         `json:"actions"`
}

// ContextSignals represents signals about the current context
type ContextSignals struct {
	GoodFaithDetected   bool    `json:"good_faith_detected"`
	PowerImbalance      float64 `json:"power_imbalance"`       // 0-1
	DehumanizationRisk  float64 `json:"dehumanization_risk"`   // 0-1
	ExploitationRisk    float64 `json:"exploitation_risk"`     // 0-1
	CrueltyStatusSignal float64 `json:"cruelty_status_signal"` // 0-1
	PersonalVoltage     float64 `json:"personal_voltage"`      // 0-100
}

// RecommendMode analyzes context and recommends operating mode
func RecommendMode(signals ContextSignals) *ModeRecommendation {
	triggeredBreakers := []CircuitBreaker{}

	// Check each circuit breaker
	if signals.PowerImbalance > 0.7 {
		triggeredBreakers = append(triggeredBreakers, CircuitBreakers[0])
	}
	if signals.DehumanizationRisk > 0.5 {
		triggeredBreakers = append(triggeredBreakers, CircuitBreakers[1])
	}
	if !signals.GoodFaithDetected {
		triggeredBreakers = append(triggeredBreakers, CircuitBreakers[2])
	}
	if signals.ExploitationRisk > 0.6 {
		triggeredBreakers = append(triggeredBreakers, CircuitBreakers[3])
	}
	if signals.PersonalVoltage > 85 {
		triggeredBreakers = append(triggeredBreakers, CircuitBreakers[4])
	}
	if signals.CrueltyStatusSignal > 0.5 {
		triggeredBreakers = append(triggeredBreakers, CircuitBreakers[5])
	}

	// Determine mode
	var mode OperatingMode
	var rationale string
	var actions []string

	if len(triggeredBreakers) == 0 && signals.GoodFaithDetected && signals.PersonalVoltage < 70 {
		mode = ModeDialectical
		rationale = "Context supports dialectical engagement: good faith detected, energy available, no breakers triggered"
		actions = []string{
			"Seek to understand all perspectives",
			"Look for synthesis opportunities",
			"Hold complexity without collapsing to sides",
		}
	} else if signals.PersonalVoltage > 85 || signals.ExploitationRisk > 0.7 {
		mode = ModeDisengage
		rationale = "Energy depleted or engagement fuels harm: strategic withdrawal recommended"
		actions = []string{
			"Withdraw from active engagement",
			"Conserve energy for better opportunity",
			"Do not provide audience for harmful dynamics",
		}
	} else {
		mode = ModeBoundary
		rationale = "Protection needed: maintain boundaries while staying engaged at safe distance"
		actions = []string{
			"Set clear limits on acceptable behavior",
			"Protect yourself and others from harm",
			"Name problematic dynamics without engaging them",
		}
	}

	return &ModeRecommendation{
		RecommendedMode:   mode,
		TriggeredBreakers: triggeredBreakers,
		EnergyLevel:       1 - (signals.PersonalVoltage / 100),
		GoodFaithLevel:    boolToFloat(signals.GoodFaithDetected),
		PowerBalance:      signals.PowerImbalance,
		Rationale:         rationale,
		Actions:           actions,
	}
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// DialecticalThinkingStage represents developmental stages of dialectical thinking
type DialecticalThinkingStage string

const (
	StageDualistic   DialecticalThinkingStage = "dualistic"   // One right answer
	StageMultiplicit DialecticalThinkingStage = "multiplicit" // Multiple perspectives, one right in context
	StageDialectical DialecticalThinkingStage = "dialectical" // Multiple truths coexist
)

// ThinkingStages describes the developmental progression
var ThinkingStages = map[DialecticalThinkingStage]struct {
	Description     string
	Characteristics []string
	Limitation      string
	Growth          string
}{
	StageDualistic: {
		Description:     "One right answer exists",
		Characteristics: []string{"Good vs bad", "My side vs their side", "Authority has answers"},
		Limitation:      "Can't handle complexity or ambiguity",
		Growth:          "Exposure to multiple valid perspectives",
	},
	StageMultiplicit: {
		Description:     "Multiple perspectives exist, context matters",
		Characteristics: []string{"Different views for different situations", "Relativism", "Contextual judgment"},
		Limitation:      "Still seeks 'one right answer' for each context",
		Growth:          "Recognizing that contradictions contain truth",
	},
	StageDialectical: {
		Description:     "Multiple truths can coexist in tension",
		Characteristics: []string{"Contradictions reveal deeper truth", "Synthesis creates new understanding", "Comfortable with paradox"},
		Limitation:      "High cognitive cost, can lead to paralysis",
		Growth:          "Practical wisdom about when to act despite uncertainty",
	},
}
