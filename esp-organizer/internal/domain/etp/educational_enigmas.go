package etp

// EducationalEnigma represents a dialectical tension in education
type EducationalEnigma struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Pole1           string   `json:"pole_1"`           // One extreme
	Pole2           string   `json:"pole_2"`           // Other extreme
	RelatedSpectra  []int    `json:"related_spectra"`  // ETP spectra that map to this enigma
	TraditionalView string   `json:"traditional_view"` // How it's usually framed (choosing a side)
	DialecticalView string   `json:"dialectical_view"` // HumanOS reframe
	IntegrationPath string   `json:"integration_path"` // How to use both
	SkillsNeeded    []string `json:"skills_needed"`    // Skills for navigating this enigma
}

// EducationalEnigmas defines core dialectical tensions in education
var EducationalEnigmas = []EducationalEnigma{
	{
		ID:              1,
		Name:            "discipline_vs_relationships",
		Pole1:           "Discipline",
		Pole2:           "Relationships",
		RelatedSpectra:  []int{14}, // authority_response
		TraditionalView: "Choose: Be strict OR be caring",
		DialecticalView: "Discipline IS caring when done through relationship",
		IntegrationPath: "Build relationship first, discipline within that trust",
		SkillsNeeded:    []string{"boundary_without_cruelty", "connection_before_correction", "warm_demanding"},
	},
	{
		ID:              2,
		Name:            "data_access_vs_security",
		Pole1:           "Data Access",
		Pole2:           "Security",
		RelatedSpectra:  []int{6, 10}, // resource_allocation, risk_tolerance
		TraditionalView: "Choose: Open access OR locked down",
		DialecticalView: "Security ENABLES access by making it safe",
		IntegrationPath: "Tiered access, earned trust, transparent rules",
		SkillsNeeded:    []string{"risk_assessment", "graduated_permissions", "trust_building"},
	},
	{
		ID:              3,
		Name:            "consistency_vs_individualism",
		Pole1:           "Consistency",
		Pole2:           "Individualism",
		RelatedSpectra:  []int{15, 1}, // ambiguity_tolerance, social_gravity
		TraditionalView: "Choose: Same for everyone OR personalized",
		DialecticalView: "Consistent PRINCIPLES, individualized APPLICATION",
		IntegrationPath: "Core standards with flexible pathways",
		SkillsNeeded:    []string{"differentiation", "principled_flexibility", "fair_vs_equal"},
	},
	{
		ID:              4,
		Name:            "challenge_vs_support",
		Pole1:           "Challenge",
		Pole2:           "Support",
		RelatedSpectra:  []int{10, 9}, // risk_tolerance, self_righting_speed
		TraditionalView: "Choose: Push hard OR protect from failure",
		DialecticalView: "Support ENABLES challenge by providing safety net",
		IntegrationPath: "Productive struggle zone: challenge WITH support",
		SkillsNeeded:    []string{"zone_of_proximal_development", "scaffolding", "gradual_release"},
	},
	{
		ID:              5,
		Name:            "content_vs_process",
		Pole1:           "Content Knowledge",
		Pole2:           "Process Skills",
		RelatedSpectra:  []int{15, 8}, // ambiguity_tolerance, impulse_gap
		TraditionalView: "Choose: Facts OR thinking skills",
		DialecticalView: "Process operates ON content; content develops process",
		IntegrationPath: "Teach process through content, assess both",
		SkillsNeeded:    []string{"metacognition", "transfer", "deep_understanding"},
	},
	{
		ID:              6,
		Name:            "competition_vs_cooperation",
		Pole1:           "Competition",
		Pole2:           "Cooperation",
		RelatedSpectra:  []int{1, 6, 14}, // social_gravity, resource_allocation, authority_response
		TraditionalView: "Choose: Win/lose OR everyone gets along",
		DialecticalView: "Compete to DISCOVER, cooperate to BUILD",
		IntegrationPath: "Phased approach: diverge competitively, converge cooperatively",
		SkillsNeeded:    []string{"strategy_fluency", "mode_switching", "integration_design"},
	},
	{
		ID:              7,
		Name:            "structure_vs_freedom",
		Pole1:           "Structure",
		Pole2:           "Freedom",
		RelatedSpectra:  []int{15, 13}, // ambiguity_tolerance, agency_threshold
		TraditionalView: "Choose: Tight control OR complete autonomy",
		DialecticalView: "Structure CREATES freedom by defining safe space",
		IntegrationPath: "Clear boundaries enable exploration within them",
		SkillsNeeded:    []string{"constraint_design", "gradual_autonomy", "scaffold_removal"},
	},
	{
		ID:              8,
		Name:            "efficiency_vs_exploration",
		Pole1:           "Efficiency",
		Pole2:           "Exploration",
		RelatedSpectra:  []int{10, 15}, // risk_tolerance, ambiguity_tolerance
		TraditionalView: "Choose: Get it done fast OR learn through discovery",
		DialecticalView: "Efficient at WHAT matters; explore to find what matters",
		IntegrationPath: "Explore for direction, efficient for execution",
		SkillsNeeded:    []string{"time_allocation", "curiosity_protection", "productive_inefficiency"},
	},
	{
		ID:              9,
		Name:            "depth_vs_breadth",
		Pole1:           "Depth",
		Pole2:           "Breadth",
		RelatedSpectra:  []int{12, 15}, // presence_sensitivity, ambiguity_tolerance
		TraditionalView: "Choose: Deep expertise OR broad knowledge",
		DialecticalView: "Depth in SOME enables connection across MANY",
		IntegrationPath: "T-shaped learning: deep in areas, connected across",
		SkillsNeeded:    []string{"transfer", "pattern_recognition", "expertise_development"},
	},
	{
		ID:              10,
		Name:            "tradition_vs_innovation",
		Pole1:           "Tradition",
		Pole2:           "Innovation",
		RelatedSpectra:  []int{10, 17}, // risk_tolerance, integrity_logic
		TraditionalView: "Choose: Preserve the past OR embrace the new",
		DialecticalView: "Innovation is INFORMED by tradition; tradition is RENEWED by innovation",
		IntegrationPath: "Know tradition deeply to innovate wisely",
		SkillsNeeded:    []string{"historical_understanding", "creative_destruction", "adaptive_preservation"},
	},
}

// GetEnigmaByName returns an enigma by its name
func GetEnigmaByName(name string) *EducationalEnigma {
	for i := range EducationalEnigmas {
		if EducationalEnigmas[i].Name == name {
			return &EducationalEnigmas[i]
		}
	}
	return nil
}

// EnigmaNavigationLevel represents a person's ability to navigate a specific enigma
type EnigmaNavigationLevel struct {
	EnigmaID         int     `json:"enigma_id"`
	EnigmaName       string  `json:"enigma_name"`
	Pole1Skill       float64 `json:"pole_1_skill"`      // 0-1
	Pole2Skill       float64 `json:"pole_2_skill"`      // 0-1
	IntegrationSkill float64 `json:"integration_skill"` // 0-1
	NaturalBias      string  `json:"natural_bias"`      // Which pole they default to
	DevelopmentFocus string  `json:"development_focus"` // What to work on
}

// AssessEnigmaNavigation evaluates ability to navigate a specific enigma
func AssessEnigmaNavigation(enigmaID int, profile *ETPProfile) *EnigmaNavigationLevel {
	enigma := &EducationalEnigmas[0]
	for i := range EducationalEnigmas {
		if EducationalEnigmas[i].ID == enigmaID {
			enigma = &EducationalEnigmas[i]
			break
		}
	}

	// Calculate bias from related spectra
	spectrumSum := 0.0
	for _, spectrumID := range enigma.RelatedSpectra {
		if setting, ok := profile.Settings[spectrumID]; ok {
			spectrumSum += setting
		}
	}
	avgSetting := spectrumSum / float64(len(enigma.RelatedSpectra))

	var naturalBias string
	var pole1Skill, pole2Skill float64

	if avgSetting < 0 {
		naturalBias = enigma.Pole1
		pole1Skill = (abs(avgSetting) + 2) / 4 // Stronger in pole 1
		pole2Skill = 0.3                       // Weaker in pole 2
	} else {
		naturalBias = enigma.Pole2
		pole2Skill = (avgSetting + 2) / 4 // Stronger in pole 2
		pole1Skill = 0.3                  // Weaker in pole 1
	}

	var developmentFocus string
	if pole1Skill < 0.4 {
		developmentFocus = "Develop " + enigma.Pole1 + " skills"
	} else if pole2Skill < 0.4 {
		developmentFocus = "Develop " + enigma.Pole2 + " skills"
	} else {
		developmentFocus = "Practice integration and switching"
	}

	return &EnigmaNavigationLevel{
		EnigmaID:         enigmaID,
		EnigmaName:       enigma.Name,
		Pole1Skill:       pole1Skill,
		Pole2Skill:       pole2Skill,
		IntegrationSkill: (pole1Skill + pole2Skill) / 2 * 0.8, // Integration is harder
		NaturalBias:      naturalBias,
		DevelopmentFocus: developmentFocus,
	}
}

// HegelianDialectic represents the thesis-antithesis-synthesis pattern
type HegelianDialectic struct {
	Thesis     string `json:"thesis"`
	Antithesis string `json:"antithesis"`
	Synthesis  string `json:"synthesis"`
	Example    string `json:"example"`
	Lesson     string `json:"lesson"`
}

// HistoricalDialectics provides examples for teaching
var HistoricalDialectics = []HegelianDialectic{
	{
		Thesis:     "Pure Competition (19th century capitalism)",
		Antithesis: "Pure Cooperation (20th century communism)",
		Synthesis:  "Social market economy (mixed systems)",
		Example:    "Nordic countries combining market competition with social cooperation",
		Lesson:     "Neither extreme worked alone; synthesis required",
	},
	{
		Thesis:     "Authority-based education (traditional)",
		Antithesis: "Child-centered education (progressive)",
		Synthesis:  "Structured exploration (balanced approach)",
		Example:    "Montessori: Structure enables freedom",
		Lesson:     "Children need both guidance AND agency",
	},
	{
		Thesis:     "Individualism (self-reliance)",
		Antithesis: "Collectivism (group-reliance)",
		Synthesis:  "Interdependence (mutual support)",
		Example:    "Open source: Individual contributions building collective resources",
		Lesson:     "Strong individuals make strong communities make strong individuals",
	},
}

// DialecticalExercise represents a practice for developing dialectical thinking
type DialecticalExercise struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Steps        []string `json:"steps"`
	ForEnigma    []int    `json:"for_enigma"` // Which enigmas this helps with
	Difficulty   string   `json:"difficulty"` // beginner, intermediate, advanced
	TimeRequired string   `json:"time_required"`
}

// DialecticalExercises for curriculum
var DialecticalExercises = []DialecticalExercise{
	{
		Name:        "steel_man",
		Description: "Make the strongest possible case for the opposing view",
		Steps: []string{
			"Identify a position you disagree with",
			"Research it thoroughly from its advocates",
			"Write the best possible argument FOR it",
			"Check with someone who holds that view - did you get it right?",
			"Reflect: What did you learn?",
		},
		ForEnigma:    []int{1, 3, 6},
		Difficulty:   "intermediate",
		TimeRequired: "30-60 minutes",
	},
	{
		Name:        "both_and_design",
		Description: "Design a system that meets both poles' needs",
		Steps: []string{
			"Identify the two poles of an enigma",
			"List what each pole genuinely needs",
			"Brainstorm: How could BOTH needs be met?",
			"Design a system that integrates both",
			"Test: Does each pole feel their need is addressed?",
		},
		ForEnigma:    []int{2, 4, 7},
		Difficulty:   "advanced",
		TimeRequired: "60-90 minutes",
	},
	{
		Name:        "strategy_switch",
		Description: "Practice switching between cooperation and competition",
		Steps: []string{
			"Start a task competitively (5 minutes)",
			"Switch to cooperative mode (5 minutes)",
			"Reflect: What changed? What worked better when?",
			"Try the opposite order",
			"Design: When would you use each?",
		},
		ForEnigma:    []int{6},
		Difficulty:   "beginner",
		TimeRequired: "20-30 minutes",
	},
	{
		Name:        "dissonance_sitting",
		Description: "Practice holding contradictory truths without resolving them",
		Steps: []string{
			"Find a genuine contradiction (e.g., 'I need structure AND freedom')",
			"Don't try to resolve it",
			"Sit with the tension for 5 minutes",
			"Notice the discomfort",
			"Ask: What does the tension reveal?",
		},
		ForEnigma:    []int{3, 7, 8},
		Difficulty:   "advanced",
		TimeRequired: "15-20 minutes",
	},
	{
		Name:        "pole_cross_training",
		Description: "Deliberately practice your weaker pole",
		Steps: []string{
			"Identify your natural bias (which pole you default to)",
			"Pick one skill from the opposite pole",
			"Practice that skill in a low-stakes situation",
			"Reflect: What was hard? What was valuable?",
			"Repeat until it feels less foreign",
		},
		ForEnigma:    []int{1, 4, 5, 6},
		Difficulty:   "beginner",
		TimeRequired: "15-30 minutes daily",
	},
}
