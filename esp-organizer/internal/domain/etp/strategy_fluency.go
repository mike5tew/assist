package etp

// StrategyType represents the two fundamental survival strategies
type StrategyType string

const (
	StrategyCooperation StrategyType = "cooperation"
	StrategyCompetition StrategyType = "competition"
)

// Strategy represents a survival strategy with its characteristics
type Strategy struct {
	Type         StrategyType `json:"type"`
	VoltageLogic string       `json:"voltage_logic"`
	Strengths    []string     `json:"strengths"`
	Weaknesses   []string     `json:"weaknesses"`
	ETPProfile   []string     `json:"etp_profile"` // ETP spectra that align with this strategy
}

// Strategies defines the two fundamental strategies
var Strategies = map[StrategyType]Strategy{
	StrategyCooperation: {
		Type:         StrategyCooperation,
		VoltageLogic: "We optimize together",
		Strengths:    []string{"Builds complex systems", "Creates safety nets", "Long-term stability", "Resource pooling"},
		Weaknesses:   []string{"Can stagnate", "Vulnerable to defectors", "Slow decision-making", "Free-rider problem"},
		ETPProfile:   []string{"social_gravity:+", "resource_allocation:+", "authority_response:-", "presence_sensitivity:+"},
	},
	StrategyCompetition: {
		Type:         StrategyCompetition,
		VoltageLogic: "I optimize first",
		Strengths:    []string{"Rapid innovation", "Strong individual agency", "Quick decisions", "Merit selection"},
		Weaknesses:   []string{"Destructive cycles", "Burns out systems", "Winner-take-all", "Trust erosion"},
		ETPProfile:   []string{"social_gravity:-", "resource_allocation:-", "authority_response:+", "presence_sensitivity:-"},
	},
}

// SituationFactor represents factors that influence strategy selection
type SituationFactor string

const (
	FactorProblemComplexity SituationFactor = "problem_complexity" // Simple vs Complex
	FactorTimeHorizon       SituationFactor = "time_horizon"       // Short-term vs Long-term
	FactorResourceState     SituationFactor = "resource_state"     // Scarce vs Abundant
	FactorTrustLevel        SituationFactor = "trust_level"        // Low vs High
	FactorPowerBalance      SituationFactor = "power_balance"      // Imbalanced vs Balanced
	FactorStakesLevel       SituationFactor = "stakes_level"       // Low vs High
)

// SituationAnalysis represents analysis of a situation for strategy selection
type SituationAnalysis struct {
	Factors             map[SituationFactor]float64 `json:"factors"` // -1 to +1 for each factor
	RecommendedStrategy StrategyType                `json:"recommended_strategy"`
	RecommendedPhases   []StrategyPhase             `json:"recommended_phases"`
	IntegrationPossible bool                        `json:"integration_possible"`
	Rationale           string                      `json:"rationale"`
}

// StrategyPhase represents a phase in an integrated approach
type StrategyPhase struct {
	Phase       int          `json:"phase"`
	Strategy    StrategyType `json:"strategy"`
	Description string       `json:"description"`
	Duration    string       `json:"duration"`
	Transition  string       `json:"transition"` // How to transition to next phase
}

// AnalyzeSituation determines optimal strategy mix for a situation
func AnalyzeSituation(factors map[SituationFactor]float64) *SituationAnalysis {
	// Calculate strategy scores
	cooperationScore := 0.0
	competitionScore := 0.0

	// Problem complexity: Complex favors cooperation
	if complexity, ok := factors[FactorProblemComplexity]; ok {
		cooperationScore += complexity * 0.3
		competitionScore -= complexity * 0.2
	}

	// Time horizon: Long-term favors cooperation
	if timeHorizon, ok := factors[FactorTimeHorizon]; ok {
		cooperationScore += timeHorizon * 0.25
		competitionScore -= timeHorizon * 0.15
	}

	// Resource state: Abundance enables cooperation
	if resources, ok := factors[FactorResourceState]; ok {
		cooperationScore += resources * 0.2
		competitionScore -= resources * 0.1
	}

	// Trust level: High trust enables cooperation
	if trust, ok := factors[FactorTrustLevel]; ok {
		cooperationScore += trust * 0.25
		competitionScore -= trust * 0.2
	}

	// Determine recommendation
	var recommended StrategyType
	var rationale string
	phases := []StrategyPhase{}

	diff := cooperationScore - competitionScore
	integrationPossible := abs(diff) < 0.3

	if integrationPossible {
		// Recommend phased approach
		phases = []StrategyPhase{
			{Phase: 1, Strategy: StrategyCompetition, Description: "Divergent ideation", Duration: "20%", Transition: "When options identified"},
			{Phase: 2, Strategy: StrategyCooperation, Description: "Evaluate and select", Duration: "30%", Transition: "When direction chosen"},
			{Phase: 3, Strategy: StrategyCompetition, Description: "Individual execution", Duration: "30%", Transition: "When components ready"},
			{Phase: 4, Strategy: StrategyCooperation, Description: "Integration and refinement", Duration: "20%", Transition: "Completion"},
		}
		recommended = StrategyCooperation // Default to cooperation for integration
		rationale = "Situation supports integrated approach: use competition for discovery, cooperation for building"
	} else if cooperationScore > competitionScore {
		recommended = StrategyCooperation
		rationale = "High complexity, trust, or long-term focus suggests cooperation-first approach"
	} else {
		recommended = StrategyCompetition
		rationale = "Low trust, short timeline, or need for rapid testing suggests competition-first approach"
	}

	return &SituationAnalysis{
		Factors:             factors,
		RecommendedStrategy: recommended,
		RecommendedPhases:   phases,
		IntegrationPossible: integrationPossible,
		Rationale:           rationale,
	}
}

// StrategyFluencyLevel represents a person's ability to use both strategies
type StrategyFluencyLevel struct {
	CooperationSkill   float64      `json:"cooperation_skill"`   // 0-1
	CompetitionSkill   float64      `json:"competition_skill"`   // 0-1
	SwitchingAbility   float64      `json:"switching_ability"`   // 0-1
	IntegrationAbility float64      `json:"integration_ability"` // 0-1
	NaturalBias        StrategyType `json:"natural_bias"`        // Which they default to
	DevelopmentAreas   []string     `json:"development_areas"`
}

// AssessStrategyFluency evaluates strategy fluency from ETP profile
func AssessStrategyFluency(profile *ETPProfile) *StrategyFluencyLevel {
	// Calculate natural bias from relevant spectra
	cooperationIndicators := []int{1, 5, 6, 12} // social_gravity, mirror_neuron, resource_allocation, presence
	competitionIndicators := []int{10, 13, 14}  // risk_tolerance, agency, authority

	coopSum := 0.0
	compSum := 0.0

	for _, id := range cooperationIndicators {
		if setting, ok := profile.Settings[id]; ok {
			coopSum += setting
		}
	}
	for _, id := range competitionIndicators {
		if setting, ok := profile.Settings[id]; ok {
			compSum += setting
		}
	}

	coopAvg := coopSum / float64(len(cooperationIndicators))
	compAvg := compSum / float64(len(competitionIndicators))

	var naturalBias StrategyType
	if coopAvg > compAvg {
		naturalBias = StrategyCooperation
	} else {
		naturalBias = StrategyCompetition
	}

	// Estimate skills (positive = skill in that direction)
	coopSkill := (coopAvg + 2) / 4 // Normalize to 0-1
	compSkill := (compAvg + 2) / 4

	// Development areas
	developmentAreas := []string{}
	if coopSkill < 0.4 {
		developmentAreas = append(developmentAreas, "Cooperation tools: trust-building, resource sharing, consensus")
	}
	if compSkill < 0.4 {
		developmentAreas = append(developmentAreas, "Competition tools: assertiveness, individual initiative, risk-taking")
	}

	return &StrategyFluencyLevel{
		CooperationSkill:   coopSkill,
		CompetitionSkill:   compSkill,
		SwitchingAbility:   0.5, // Requires assessment
		IntegrationAbility: 0.5, // Requires assessment
		NaturalBias:        naturalBias,
		DevelopmentAreas:   developmentAreas,
	}
}
