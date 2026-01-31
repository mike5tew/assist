package etp

// CompatibilitySolution defines how to bridge different voltage settings on a spectrum
type CompatibilitySolution struct {
	SpectrumID      int      `json:"spectrum_id"`
	Name            string   `json:"name"`
	NegativeSupport string   `json:"negative_support"` // Support for -2 end
	PositiveSupport string   `json:"positive_support"` // Support for +2 end
	BridgeActivity  string   `json:"bridge_activity"`  // Activity for mixed groups
	Language        []string `json:"language"`         // Voltage-aware phrases
}

// CompatibilitySolutions maps each of the 8 core ETP spectra to its compatibility solution
var CompatibilitySolutions = map[int]CompatibilitySolution{
	1: {
		SpectrumID:      1,
		Name:            "Social Voltage Zones",
		NegativeSupport: "Recharge Zone: Quiet, alone space for Independents",
		PositiveSupport: "Connection Zone: Group interaction for Cohesives",
		BridgeActivity:  "Paired work with clear recharge breaks",
		Language: []string{
			"I need to recharge",
			"I need connection",
			"I'm going to recharge so I can be fully present later",
		},
	},
	2: {
		SpectrumID:      2,
		Name:            "Energy Circuit Design",
		NegativeSupport: "Parallel play: Same activity, lower engagement level",
		PositiveSupport: "Energy outlet: External expression of internal energy",
		BridgeActivity:  "Energy trading: 'You talk while I listen, then we switch'",
		Language: []string{
			"I need parallel play right now",
			"Can we trade energy? You lead, I'll follow, then switch",
			"I need buffer time between social activities",
		},
	},
	3: {
		SpectrumID:      3,
		Name:            "Voltage Step-Down Transformers",
		NegativeSupport: "Gradual exposure: Low-voltage → medium → high",
		PositiveSupport: "Full participation with processing time after",
		BridgeActivity:  "Voltage mixing: Insulated observes, Conductive participates, then discuss",
		Language: []string{
			"Too much voltage",
			"More voltage please",
			"I need to observe first, then participate",
		},
	},
	4: {
		SpectrumID:      4,
		Name:            "Threat Voltage Calibration",
		NegativeSupport: "Passive strategies: Validate withdrawal as legitimate tactic",
		PositiveSupport: "Aggressive channels: Safe outlets (debate, advocacy, sports)",
		BridgeActivity:  "Threat assessment: 'When to stand ground, when to withdraw - both are tactics'",
		Language: []string{
			"Walking away IS a power move",
			"I choose my battles strategically",
			"Confrontation is a tool, not a personality",
			"Passive here, aggressive there - that's fluency",
		},
	},
	5: {
		SpectrumID:      5,
		Name:            "Care Voltage Boundaries",
		NegativeSupport: "Healthy detachment: Boundaries ARE care - for yourself",
		PositiveSupport: "Sustainable giving: Care without depletion",
		BridgeActivity:  "Care calibration: 'Put your own oxygen mask on first, then help others'",
		Language: []string{
			"I can care without drowning",
			"Boundaries ARE care - for myself",
			"Detached is not cold, it's sustainable",
			"I choose who receives my care voltage",
		},
	},
	6: {
		SpectrumID:      6,
		Name:            "Risk Voltage Gradients",
		NegativeSupport: "Scaffolded risk: 1% risk → 5% → 10%",
		PositiveSupport: "Risk roles: Seeker scouts ahead",
		BridgeActivity:  "Risk pairing: Seeker scouts, Averse maps safety",
		Language: []string{
			"This feels X% risky to me",
			"I'll scout, you map safety",
			"Let's start with 1% risk",
		},
	},
	7: {
		SpectrumID:      7,
		Name:            "Integrity Voltage Framing",
		NegativeSupport: "Context markers: 'In this context, the principle is...'",
		PositiveSupport: "Core principles: 'These 3 things are absolute'",
		BridgeActivity:  "Principle evolution: 'This used to be absolute, now contextual because...'",
		Language: []string{
			"These 3 things are absolute, the rest is relative",
			"In this context, the principle is...",
			"This used to be absolute, now it's contextual because...",
		},
	},
	8: {
		SpectrumID:      8,
		Name:            "Empathy Voltage Filters",
		NegativeSupport: "Cueing: 'They're experiencing this strongly' (signal, not share)",
		PositiveSupport: "Boundary reminder: 'This is their feeling, not yours'",
		BridgeActivity:  "Shared language: 'I'm resonating with your X' / 'I see you're feeling X'",
		Language: []string{
			"I'm resonating with your feeling",
			"I see you're feeling X",
			"This is their feeling, not mine",
			"They're experiencing this strongly",
		},
	},
}

// GetCompatibilitySolution returns the solution for a given spectrum ID
func GetCompatibilitySolution(spectrumID int) *CompatibilitySolution {
	if sol, ok := CompatibilitySolutions[spectrumID]; ok {
		return &sol
	}
	return nil
}

// VoltageCompatibilityCheck assesses compatibility between two ETP settings on a spectrum
type VoltageCompatibilityCheck struct {
	SpectrumID     int     `json:"spectrum_id"`
	SpectrumName   string  `json:"spectrum_name"`
	Setting1       float64 `json:"setting_1"`       // -2 to +2
	Setting2       float64 `json:"setting_2"`       // -2 to +2
	Difference     float64 `json:"difference"`      // Absolute difference
	ConflictRisk   string  `json:"conflict_risk"`   // low, medium, high
	SolutionNeeded bool    `json:"solution_needed"` // True if difference > 2
	Solution       string  `json:"solution"`        // Recommended solution name
}

// AssessCompatibility checks how compatible two ETP settings are
func AssessCompatibility(spectrumID int, setting1, setting2 float64) *VoltageCompatibilityCheck {
	spectrum := GetSpectrumByID(spectrumID)
	if spectrum == nil {
		return nil
	}

	diff := abs(setting1 - setting2)

	risk := "low"
	if diff > 2 {
		risk = "medium"
	}
	if diff > 3 {
		risk = "high"
	}

	solution := ""
	needsSolution := diff > 2
	if needsSolution {
		if sol := GetCompatibilitySolution(spectrumID); sol != nil {
			solution = sol.Name
		}
	}

	return &VoltageCompatibilityCheck{
		SpectrumID:     spectrumID,
		SpectrumName:   spectrum.Name,
		Setting1:       setting1,
		Setting2:       setting2,
		Difference:     diff,
		ConflictRisk:   risk,
		SolutionNeeded: needsSolution,
		Solution:       solution,
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
