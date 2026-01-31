package etp

// IdentityHijackPattern represents patterns of cognitive identity hijacking
type IdentityHijackPattern struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	HijackedSpectra   []int    `json:"hijacked_spectra"`   // Which ETP spectra are being exploited
	TribalSignals     []string `json:"tribal_signals"`     // Language/behavior signals
	NeurochemicalHook string   `json:"neurochemical_hook"` // What biological reward it offers
	ExitDifficulty    string   `json:"exit_difficulty"`    // How hard to escape
}

// IdentityHijackPatterns defines common hijack patterns
var IdentityHijackPatterns = []IdentityHijackPattern{
	{
		Name:              "us_vs_them",
		Description:       "Creating artificial in-group/out-group division",
		HijackedSpectra:   []int{1, 6, 16}, // social_gravity, resource_allocation, status_sensitivity
		TribalSignals:     []string{"We/they language", "Loyalty tests", "Enemy naming", "Purity requirements"},
		NeurochemicalHook: "Belonging dopamine + enemy adrenaline",
		ExitDifficulty:    "High - threatens identity and belonging",
	},
	{
		Name:              "certainty_addiction",
		Description:       "Converting complex issues into simple binaries",
		HijackedSpectra:   []int{15, 17}, // ambiguity_tolerance, integrity_logic
		TribalSignals:     []string{"Absolute statements", "No nuance allowed", "Doubt = betrayal"},
		NeurochemicalHook: "Cognitive ease dopamine + anxiety reduction",
		ExitDifficulty:    "Medium - complexity feels threatening",
	},
	{
		Name:              "moral_superiority",
		Description:       "Identity based on being morally better than out-group",
		HijackedSpectra:   []int{17, 2}, // integrity_logic, guilt_response
		TribalSignals:     []string{"Virtue signaling", "Shame weaponization", "Moral purity tests"},
		NeurochemicalHook: "Status serotonin + righteousness dopamine",
		ExitDifficulty:    "High - admitting complexity feels like moral failure",
	},
	{
		Name:              "threat_amplification",
		Description:       "Exaggerating threat to justify extreme responses",
		HijackedSpectra:   []int{10, 11}, // risk_tolerance, anticipation_bias
		TribalSignals:     []string{"Catastrophizing", "Imminent danger claims", "Preemptive action demands"},
		NeurochemicalHook: "Adrenaline activation + protective dopamine",
		ExitDifficulty:    "Medium - requires felt safety to lower",
	},
	{
		Name:              "cruelty_as_strength",
		Description:       "Reframing cruelty as power/loyalty demonstration",
		HijackedSpectra:   []int{13, 14, 16}, // agency_threshold, authority_response, status_sensitivity
		TribalSignals:     []string{"Cruelty applause", "Weakness accusations", "Mercy = betrayal"},
		NeurochemicalHook: "Dominance testosterone + tribal dopamine + adrenaline",
		ExitDifficulty:    "Very high - backed by status rewards",
	},
}

// HijackDetection represents detection of identity hijacking in a context
type HijackDetection struct {
	DetectedPatterns []IdentityHijackPattern `json:"detected_patterns"`
	SeverityLevel    float64                 `json:"severity_level"`  // 0-1
	TribalThinking   float64                 `json:"tribal_thinking"` // 0-1
	SignalsDetected  []string                `json:"signals_detected"`
	Recommendation   string                  `json:"recommendation"`
	DeescalationPath string                  `json:"deescalation_path"`
}

// DetectIdentityHijacking analyzes language/behavior for hijack patterns
func DetectIdentityHijacking(signals []string) *HijackDetection {
	detected := []IdentityHijackPattern{}
	detectedSignals := []string{}

	for _, pattern := range IdentityHijackPatterns {
		for _, patternSignal := range pattern.TribalSignals {
			for _, inputSignal := range signals {
				if containsPattern(inputSignal, patternSignal) {
					detected = append(detected, pattern)
					detectedSignals = append(detectedSignals, inputSignal)
					break
				}
			}
		}
	}

	severity := float64(len(detected)) / float64(len(IdentityHijackPatterns))

	var recommendation string
	var deescalation string

	if severity > 0.6 {
		recommendation = "High identity hijacking detected. Switch to BOUNDARY mode. Do not engage dialectically - you'll be exploited."
		deescalation = "Remove yourself from the environment. Do not provide audience."
	} else if severity > 0.3 {
		recommendation = "Moderate tribal signals detected. Proceed with caution. Name patterns without attacking people."
		deescalation = "Redirect to shared concerns. 'We both want X, we disagree on how.'"
	} else {
		recommendation = "Low tribal signaling. Dialectical engagement may be productive."
		deescalation = "Continue with curiosity and boundary awareness."
	}

	return &HijackDetection{
		DetectedPatterns: detected,
		SeverityLevel:    severity,
		TribalThinking:   severity,
		SignalsDetected:  detectedSignals,
		Recommendation:   recommendation,
		DeescalationPath: deescalation,
	}
}

func containsPattern(input, pattern string) bool {
	// Simplified pattern matching - would be more sophisticated in production
	return len(input) > 0 && len(pattern) > 0
}

// CognitiveSelfDefenseSkill represents skills for defending against identity hijacking
type CognitiveSelfDefenseSkill struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Counters    []string `json:"counters"` // Which hijack patterns it counters
	Practice    string   `json:"practice"` // How to practice this skill
	Signs       string   `json:"signs"`    // Signs you need this skill
}

// CognitiveSelfDefenseSkills defines the cognitive self-defense curriculum
var CognitiveSelfDefenseSkills = []CognitiveSelfDefenseSkill{
	{
		Name:        "hijack_detection",
		Description: "Recognizing when thinking becomes tribal",
		Counters:    []string{"us_vs_them", "certainty_addiction"},
		Practice:    "Notice us/them language in yourself. Pause when you feel certainty.",
		Signs:       "Finding yourself using 'always', 'never', 'they all'",
	},
	{
		Name:        "identity_separation",
		Description: "Separating strategies from identity",
		Counters:    []string{"moral_superiority", "us_vs_them"},
		Practice:    "Say 'I'm using cooperation tools' not 'I am a cooperator'",
		Signs:       "Feeling personally attacked when your strategy is questioned",
	},
	{
		Name:        "complexity_tolerance",
		Description: "Holding contradictory truths simultaneously",
		Counters:    []string{"certainty_addiction"},
		Practice:    "Find one true thing in the opposing view. Hold both truths.",
		Signs:       "Needing to resolve ambiguity immediately",
	},
	{
		Name:        "voltage_awareness",
		Description: "Recognizing when emotional state affects thinking",
		Counters:    []string{"threat_amplification"},
		Practice:    "Check voltage before responding. High voltage = pause.",
		Signs:       "Making decisions you regret when calm",
	},
	{
		Name:        "status_source_audit",
		Description: "Examining where you get status/belonging",
		Counters:    []string{"cruelty_as_strength"},
		Practice:    "Ask: Am I doing this for status? Would I do it without audience?",
		Signs:       "Behavior changes based on who's watching",
	},
	{
		Name:        "dialectical_questioning",
		Description: "Actively seeking opposing perspectives",
		Counters:    []string{"us_vs_them", "moral_superiority"},
		Practice:    "Steel-man the opposing view. Make their best argument.",
		Signs:       "Only consuming media that agrees with you",
	},
}

// TribalThinkingIndicator represents measurable indicators of tribal thinking
type TribalThinkingIndicator struct {
	Indicator   string  `json:"indicator"`
	Weight      float64 `json:"weight"`       // How strongly this indicates tribal thinking
	CounterSign string  `json:"counter_sign"` // What dialectical thinking looks like instead
}

// TribalThinkingIndicators for assessment
var TribalThinkingIndicators = []TribalThinkingIndicator{
	{
		Indicator:   "Uses us/them language",
		Weight:      0.2,
		CounterSign: "Refers to shared humanity, varied perspectives within groups",
	},
	{
		Indicator:   "Claims certainty on complex issues",
		Weight:      0.15,
		CounterSign: "Acknowledges uncertainty, multiple valid views",
	},
	{
		Indicator:   "Attributes worst motives to out-group",
		Weight:      0.2,
		CounterSign: "Assumes others have reasons that make sense to them",
	},
	{
		Indicator:   "Dismisses complexity as 'both sides-ism'",
		Weight:      0.15,
		CounterSign: "Engages with nuance while maintaining values",
	},
	{
		Indicator:   "Treats disagreement as betrayal",
		Weight:      0.15,
		CounterSign: "Values loyal disagreement over compliant agreement",
	},
	{
		Indicator:   "Celebrates harm to out-group",
		Weight:      0.15,
		CounterSign: "Seeks solutions that address all parties' needs",
	},
}
