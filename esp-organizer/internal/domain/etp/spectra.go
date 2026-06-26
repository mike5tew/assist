package etp

// ETPSpectrum defines a single ETP spectrum with voltage semantics
type ETPSpectrum struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	NegativeLabel   string   `json:"negative_label"`
	PositiveLabel   string   `json:"positive_label"`
	NegativeEffect  string   `json:"negative_effect"`
	PositiveEffect  string   `json:"positive_effect"`
	ConflictPattern string   `json:"conflict_pattern"`
	SolutionName    string   `json:"solution_name"`
	Skills          []string `json:"skills"` // Trainable skills to expand range on this spectrum
}

// SpectrumSkill defines a trainable skill associated with a spectrum
type SpectrumSkill struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SpectrumID  int    `json:"spectrum_id"`
	Direction   string `json:"direction"` // "negative", "positive", or "both"
	Description string `json:"description"`
	Exercises   string `json:"exercises"`
}

// ETPSpectra defines the 12 CORE biological spectra
// These are innate, morally neutral, and orthogonal to each other
var ETPSpectra = []ETPSpectrum{
	{
		ID:              1,
		Name:            "social_gravity",
		NegativeLabel:   "Independent",
		PositiveLabel:   "Cohesive",
		NegativeEffect:  "Social interaction DRAINS voltage",
		PositiveEffect:  "Social interaction CHARGES voltage",
		ConflictPattern: "Independent feels overwhelmed, Cohesive feels lonely",
		SolutionName:    "Social Voltage Zones",
		Skills: []string{
			"solitude_tolerance",     // Skill for cohesive to handle alone time
			"group_navigation",       // Skill for independent to handle group settings
			"social_recovery",        // Skill for managing post-social depletion
			"status_awareness",       // From old status_sensitivity - reading social hierarchies
			"connection_maintenance", // Keeping relationships alive without draining
		},
	},
	{
		ID:              2,
		Name:            "energy_directionality",
		NegativeLabel:   "Inward",
		PositiveLabel:   "Outward",
		NegativeEffect:  "External stimulation OVERLOADS system",
		PositiveEffect:  "Internal energy NEEDS external outlet",
		ConflictPattern: "Introverts drain, Extroverts starve",
		SolutionName:    "Energy Circuit Design",
		Skills: []string{
			"quiet_focus",            // Skill for outward to focus internally
			"public_activation",      // Skill for inward to engage externally
			"energy_boundaries",      // Setting limits on energy expenditure
			"stimulation_modulation", // Adjusting input levels
		},
	},
	{
		ID:              3,
		Name:            "voltage_sensitivity",
		NegativeLabel:   "Insulated",
		PositiveLabel:   "Conductive",
		NegativeEffect:  "Needs thick emotional insulation",
		PositiveEffect:  "Thrives on emotional current",
		ConflictPattern: "Insulated seems numb, Conductive seems dramatic",
		SolutionName:    "Voltage Step-Down Transformers",
		Skills: []string{
			"emotional_dampening",    // Skill for conductive to reduce intensity
			"emotional_opening",      // Skill for insulated to feel more
			"overwhelm_recovery",     // From old self_righting_speed - recovering after overload
			"emotional_transparency", // From old spectrum - choosing what to show
			"voltage_grounding",      // Discharging excess without explosion
		},
	},
	{
		ID:              4,
		Name:            "threat_response",
		NegativeLabel:   "Passive",
		PositiveLabel:   "Aggressive",
		NegativeEffect:  "Threats trigger AVOIDANCE voltage (freeze/flee)",
		PositiveEffect:  "Threats trigger CONFRONTATION voltage (fight)",
		ConflictPattern: "Passive seen as cowardly, Aggressive seen as dangerous",
		SolutionName:    "Threat Voltage Calibration",
		Skills: []string{
			"strategic_withdrawal",    // Skill for aggressive to disengage
			"assertive_confrontation", // Skill for passive to stand ground
			"de_escalation",           // Reducing threat voltage in situations
			"agency_activation",       // From old agency_threshold - taking action
			"authority_navigation",    // From old authority_response - dealing with power
			"boundary_setting",        // Clear limits without aggression
		},
	},
	{
		ID:              5,
		Name:            "care_response",
		NegativeLabel:   "Detached",
		PositiveLabel:   "Nurturing",
		NegativeEffect:  "Vulnerability triggers DISTANCE voltage (self-protection)",
		PositiveEffect:  "Vulnerability triggers CARE voltage (other-protection)",
		ConflictPattern: "Detached seen as cold, Nurturing seen as smothering",
		SolutionName:    "Care Voltage Boundaries",
		Skills: []string{
			"healthy_detachment",    // Skill for nurturing to maintain boundaries
			"sustainable_giving",    // Skill for nurturing to not deplete
			"compassion_activation", // Skill for detached to engage care
			"resource_sharing",      // From old resource_allocation - sharing without depleting
			"self_care_priority",    // Caring for self first (oxygen mask)
			"empathic_action",       // Translating feeling into helping
		},
	},
	{
		ID:              6,
		Name:            "risk_tolerance",
		NegativeLabel:   "Averse",
		PositiveLabel:   "Seeking",
		NegativeEffect:  "Risk creates ANXIETY voltage",
		PositiveEffect:  "Risk creates EXCITEMENT voltage",
		ConflictPattern: "Averse holds back, Seeking pushes forward",
		SolutionName:    "Risk Voltage Gradients",
		Skills: []string{
			"calculated_risk_taking",  // Skill for averse to take measured risks
			"prudent_restraint",       // Skill for seeking to hold back
			"risk_assessment",         // Accurate evaluation of actual risk
			"anticipation_management", // From old anticipation_bias - managing expectations
			"ambiguity_tolerance",     // From old spectrum - sitting with uncertainty
			"failure_recovery",        // Bouncing back from risks gone wrong
		},
	},
	{
		ID:              7,
		Name:            "integrity_logic",
		NegativeLabel:   "Relativistic",
		PositiveLabel:   "Absolutist",
		NegativeEffect:  "Absolutes create CONSTRICTION voltage",
		PositiveEffect:  "Relativity creates CHAOS voltage",
		ConflictPattern: "Relativistic seems unprincipled, Absolutist seems rigid",
		SolutionName:    "Integrity Voltage Framing",
		Skills: []string{
			"contextual_reasoning",  // Skill for absolutist to see shades
			"principle_standing",    // Skill for relativist to hold firm
			"dialectical_synthesis", // Finding truth in opposing views
			"guilt_processing",      // From old guilt_response - healthy guilt handling
			"moral_flexibility",     // Adapting principles without abandoning them
			"value_articulation",    // Clearly stating what you stand for
		},
	},
	{
		ID:              8,
		Name:            "mirror_neuron_tuning",
		NegativeLabel:   "Selective",
		PositiveLabel:   "Absorbent",
		NegativeEffect:  "Others' emotions are DISTINCT signals",
		PositiveEffect:  "Others' emotions are SHARED experience",
		ConflictPattern: "Selective seems uncaring, Absorbent seems overwhelmed",
		SolutionName:    "Empathy Voltage Filters",
		Skills: []string{
			"empathy_filtering",           // Skill for absorbent to not take on everything
			"emotional_attunement",        // Skill for selective to tune in more
			"self_other_boundary",         // Distinguishing your feelings from theirs
			"presence_modulation",         // From old presence_sensitivity - attention control
			"emotional_contagion_control", // Not catching unwanted emotions
			"compassionate_witnessing",    // Feeling without fixing
		},
	},
	{
		ID:              9,
		Name:            "orderliness",
		NegativeLabel:   "Flexible",
		PositiveLabel:   "Ordered",
		NegativeEffect:  "Structure creates CONSTRICTION voltage (feels caged)",
		PositiveEffect:  "Disorder creates ANXIETY voltage (feels chaotic)",
		ConflictPattern: "Flexible seems chaotic, Ordered seems controlling",
		SolutionName:    "Structure Voltage Gradients",
		Skills: []string{
			"routine_building",      // Skill for flexible to establish useful structure
			"spontaneity_tolerance", // Skill for ordered to handle disruption
			"task_sequencing",       // Breaking work into manageable ordered steps
			"workspace_management",  // Physical environment organisation
			"transition_handling",   // Moving between activities without dysregulation
			"flexible_planning",     // Creating structure that bends without breaking
		},
	},
	{
		ID:              10,
		Name:            "responsibility_threshold",
		NegativeLabel:   "Deflecting",
		PositiveLabel:   "Absorbing",
		NegativeEffect:  "Consequence feels EXTERNAL — not my problem",
		PositiveEffect:  "Consequence feels INTERNAL — everything is my fault",
		ConflictPattern: "Deflecting seems irresponsible, Absorbing seems codependent",
		SolutionName:    "Responsibility Voltage Calibration",
		Skills: []string{
			"ownership_activation",    // Skill for deflecting to take responsibility
			"responsibility_release",  // Skill for absorbing to let go of others' consequences
			"consequence_prediction",  // Foreseeing outcomes of actions
			"accountability_language", // Expressing ownership without shame
			"delegation_comfort",      // Letting others own their part
		},
	},
	{
		ID:              11,
		Name:            "loss_sensitivity",
		NegativeLabel:   "Detached",
		PositiveLabel:   "Territorial",
		NegativeEffect:  "Removal creates MINIMAL voltage — lets things go easily",
		PositiveEffect:  "Removal creates EXISTENTIAL voltage — feels like annihilation",
		ConflictPattern: "Detached seems uncaring about possessions, Territorial seems obsessive",
		SolutionName:    "Loss Voltage Gradients",
		Skills: []string{
			"secure_attachment",   // Skill for territorial to feel safe without possessions
			"value_recognition",   // Skill for detached to appreciate what they have
			"transition_rituals",  // Managing the process of things leaving
			"abundance_awareness", // Recognising that loss is not permanent
			"grief_processing",    // Healthy processing of genuine loss
		},
	},
	{
		ID:              12,
		Name:            "libido",
		NegativeLabel:   "Restrained",
		PositiveLabel:   "Expressive",
		NegativeEffect:  "Drive energy is CONTAINED — low outward expression of desire",
		PositiveEffect:  "Drive energy is EXTERNALISED — high outward expression of desire",
		ConflictPattern: "Restrained seems disengaged, Expressive seems overwhelming",
		SolutionName:    "Drive Voltage Modulation",
		Skills: []string{
			"desire_articulation",   // Naming what you want without shame
			"impulse_channeling",    // Directing drive energy into productive outlets
			"boundary_awareness",    // Recognising others' comfort thresholds
			"delayed_gratification", // Tolerating the gap between wanting and having
			"consent_navigation",    // Understanding mutual agreement
		},
	},
}

// SpectrumSkills defines all trainable skills organized by spectrum
var SpectrumSkills = map[int][]SpectrumSkill{
	// Social Gravity Skills
	1: {
		{ID: "solitude_tolerance", Name: "Solitude Tolerance", SpectrumID: 1, Direction: "negative", Description: "Ability to be comfortable alone", Exercises: "Solo meals, solo walks, journaling"},
		{ID: "group_navigation", Name: "Group Navigation", SpectrumID: 1, Direction: "positive", Description: "Ability to function in group settings", Exercises: "Structured group activities, defined roles"},
		{ID: "social_recovery", Name: "Social Recovery", SpectrumID: 1, Direction: "both", Description: "Restoring energy after social interaction", Exercises: "Post-social quiet time, decompression rituals"},
		{ID: "status_awareness", Name: "Status Awareness", SpectrumID: 1, Direction: "positive", Description: "Reading and navigating social hierarchies", Exercises: "Observation practice, role analysis"},
		{ID: "connection_maintenance", Name: "Connection Maintenance", SpectrumID: 1, Direction: "positive", Description: "Keeping relationships alive efficiently", Exercises: "Scheduled check-ins, meaningful brief contacts"},
	},
	// Energy Directionality Skills
	2: {
		{ID: "quiet_focus", Name: "Quiet Focus", SpectrumID: 2, Direction: "negative", Description: "Sustaining internal attention", Exercises: "Meditation, deep work blocks, reading"},
		{ID: "public_activation", Name: "Public Activation", SpectrumID: 2, Direction: "positive", Description: "Engaging with external world when needed", Exercises: "Planned outings, structured social time"},
		{ID: "energy_boundaries", Name: "Energy Boundaries", SpectrumID: 2, Direction: "both", Description: "Setting limits on energy expenditure", Exercises: "Time limits, exit strategies, saying no"},
		{ID: "stimulation_modulation", Name: "Stimulation Modulation", SpectrumID: 2, Direction: "both", Description: "Adjusting environmental input", Exercises: "Noise control, lighting, space design"},
	},
	// Voltage Sensitivity Skills
	3: {
		{ID: "emotional_dampening", Name: "Emotional Dampening", SpectrumID: 3, Direction: "negative", Description: "Reducing emotional intensity", Exercises: "Grounding techniques, breathing, distraction"},
		{ID: "emotional_opening", Name: "Emotional Opening", SpectrumID: 3, Direction: "positive", Description: "Allowing more emotional experience", Exercises: "Art, music, vulnerability practice"},
		{ID: "overwhelm_recovery", Name: "Overwhelm Recovery", SpectrumID: 3, Direction: "both", Description: "Returning to baseline after emotional flooding", Exercises: "Recovery protocols, safe spaces, time allowance"},
		{ID: "emotional_transparency", Name: "Emotional Transparency", SpectrumID: 3, Direction: "both", Description: "Choosing what emotions to show", Exercises: "Graduated disclosure, trusted confidants"},
		{ID: "voltage_grounding", Name: "Voltage Grounding", SpectrumID: 3, Direction: "both", Description: "Discharging excess emotion safely", Exercises: "Physical exercise, journaling, talking it out"},
	},
	// Threat Response Skills
	4: {
		{ID: "strategic_withdrawal", Name: "Strategic Withdrawal", SpectrumID: 4, Direction: "negative", Description: "Disengaging when confrontation is counterproductive", Exercises: "Exit planning, cooling off periods, choosing battles"},
		{ID: "assertive_confrontation", Name: "Assertive Confrontation", SpectrumID: 4, Direction: "positive", Description: "Standing ground when necessary", Exercises: "Assertiveness training, role play, scripts"},
		{ID: "de_escalation", Name: "De-escalation", SpectrumID: 4, Direction: "both", Description: "Reducing threat voltage in situations", Exercises: "Tone control, validation, space creation"},
		{ID: "agency_activation", Name: "Agency Activation", SpectrumID: 4, Direction: "positive", Description: "Taking action rather than waiting", Exercises: "Small initiatives, decision practice, ownership"},
		{ID: "authority_navigation", Name: "Authority Navigation", SpectrumID: 4, Direction: "both", Description: "Effective interaction with power structures", Exercises: "Choosing compliance vs challenge strategically"},
		{ID: "boundary_setting", Name: "Boundary Setting", SpectrumID: 4, Direction: "both", Description: "Clear limits without aggression", Exercises: "I-statements, consequence clarity, consistency"},
	},
	// Care Response Skills
	5: {
		{ID: "healthy_detachment", Name: "Healthy Detachment", SpectrumID: 5, Direction: "negative", Description: "Maintaining boundaries while caring", Exercises: "Professional distance, not-my-problem practice"},
		{ID: "sustainable_giving", Name: "Sustainable Giving", SpectrumID: 5, Direction: "both", Description: "Caring without depleting self", Exercises: "Giving budgets, reciprocity tracking, rest"},
		{ID: "compassion_activation", Name: "Compassion Activation", SpectrumID: 5, Direction: "positive", Description: "Engaging care when appropriate", Exercises: "Volunteer work, kindness practice, empathy building"},
		{ID: "resource_sharing", Name: "Resource Sharing", SpectrumID: 5, Direction: "positive", Description: "Sharing without security anxiety", Exercises: "Generosity experiments, abundance mindset"},
		{ID: "self_care_priority", Name: "Self-Care Priority", SpectrumID: 5, Direction: "negative", Description: "Putting on your own oxygen mask first", Exercises: "Non-negotiable self-care, guilt-free rest"},
		{ID: "empathic_action", Name: "Empathic Action", SpectrumID: 5, Direction: "positive", Description: "Translating feeling into practical help", Exercises: "Identifying needs, offering specific help"},
	},
	// Risk Tolerance Skills
	6: {
		{ID: "calculated_risk_taking", Name: "Calculated Risk Taking", SpectrumID: 6, Direction: "positive", Description: "Taking measured risks", Exercises: "Small bets, graduated exposure, risk ladders"},
		{ID: "prudent_restraint", Name: "Prudent Restraint", SpectrumID: 6, Direction: "negative", Description: "Holding back when wise", Exercises: "Pause protocols, second opinion seeking, sleep on it"},
		{ID: "risk_assessment", Name: "Risk Assessment", SpectrumID: 6, Direction: "both", Description: "Accurate evaluation of actual risk", Exercises: "Probability training, base rate awareness"},
		{ID: "anticipation_management", Name: "Anticipation Management", SpectrumID: 6, Direction: "both", Description: "Managing expectations realistically", Exercises: "Pre-mortem, best/worst/likely scenarios"},
		{ID: "ambiguity_tolerance", Name: "Ambiguity Tolerance", SpectrumID: 6, Direction: "both", Description: "Sitting with uncertainty", Exercises: "Open-ended projects, delayed decisions, comfort with not-knowing"},
		{ID: "failure_recovery", Name: "Failure Recovery", SpectrumID: 6, Direction: "both", Description: "Bouncing back from risks gone wrong", Exercises: "Failure analysis, self-compassion, lesson extraction"},
	},
	// Integrity Logic Skills
	7: {
		{ID: "contextual_reasoning", Name: "Contextual Reasoning", SpectrumID: 7, Direction: "negative", Description: "Seeing shades of gray", Exercises: "Devil's advocate, steelmanning, perspective taking"},
		{ID: "principle_standing", Name: "Principle Standing", SpectrumID: 7, Direction: "positive", Description: "Holding firm on core values", Exercises: "Value clarification, line-drawing, non-negotiables"},
		{ID: "dialectical_synthesis", Name: "Dialectical Synthesis", SpectrumID: 7, Direction: "both", Description: "Finding truth in opposing views", Exercises: "Thesis-antithesis practice, both/and thinking"},
		{ID: "guilt_processing", Name: "Guilt Processing", SpectrumID: 7, Direction: "both", Description: "Healthy handling of guilt", Exercises: "Responsibility sorting, repair actions, self-forgiveness"},
		{ID: "moral_flexibility", Name: "Moral Flexibility", SpectrumID: 7, Direction: "negative", Description: "Adapting without abandoning principles", Exercises: "Edge case analysis, value hierarchy"},
		{ID: "value_articulation", Name: "Value Articulation", SpectrumID: 7, Direction: "both", Description: "Clearly stating what you stand for", Exercises: "Values writing, explaining to others, living examples"},
	},
	// Mirror Neuron Tuning Skills
	8: {
		{ID: "empathy_filtering", Name: "Empathy Filtering", SpectrumID: 8, Direction: "negative", Description: "Not absorbing everything", Exercises: "Shielding visualizations, selective attention, limits"},
		{ID: "emotional_attunement", Name: "Emotional Attunement", SpectrumID: 8, Direction: "positive", Description: "Tuning in to others' states", Exercises: "Active listening, body language reading, checking in"},
		{ID: "self_other_boundary", Name: "Self-Other Boundary", SpectrumID: 8, Direction: "both", Description: "Distinguishing your feelings from theirs", Exercises: "Whose feeling is this?, body scanning, grounding"},
		{ID: "presence_modulation", Name: "Presence Modulation", SpectrumID: 8, Direction: "both", Description: "Controlling attention to present moment", Exercises: "Mindfulness, attention training, focus practice"},
		{ID: "emotional_contagion_control", Name: "Emotional Contagion Control", SpectrumID: 8, Direction: "negative", Description: "Not catching unwanted emotions", Exercises: "Pre-exposure preparation, mid-exposure check, post-exposure clearing"},
		{ID: "compassionate_witnessing", Name: "Compassionate Witnessing", SpectrumID: 8, Direction: "both", Description: "Feeling with without fixing", Exercises: "Presence without advice, holding space, just listening"},
	},
	// Orderliness Skills
	9: {
		{ID: "routine_building", Name: "Routine Building", SpectrumID: 9, Direction: "positive", Description: "Establishing useful daily structure", Exercises: "Morning routine design, checklist creation, habit stacking"},
		{ID: "spontaneity_tolerance", Name: "Spontaneity Tolerance", SpectrumID: 9, Direction: "negative", Description: "Handling disruption to plans without dysregulation", Exercises: "Planned surprises, flexible scheduling, improvisation games"},
		{ID: "task_sequencing", Name: "Task Sequencing", SpectrumID: 9, Direction: "positive", Description: "Breaking work into manageable ordered steps", Exercises: "Task lists, workflow mapping, priority matrices"},
		{ID: "workspace_management", Name: "Workspace Management", SpectrumID: 9, Direction: "positive", Description: "Organising physical and digital environments", Exercises: "Desk resets, folder organisation, tool placement rituals"},
		{ID: "transition_handling", Name: "Transition Handling", SpectrumID: 9, Direction: "both", Description: "Moving between activities without dysregulation", Exercises: "Transition warnings, bridging rituals, closure practices"},
		{ID: "flexible_planning", Name: "Flexible Planning", SpectrumID: 9, Direction: "negative", Description: "Creating structure that bends without breaking", Exercises: "Plan B thinking, loose scheduling, outcome vs process focus"},
	},
	// Responsibility Threshold Skills
	10: {
		{ID: "ownership_activation", Name: "Ownership Activation", SpectrumID: 10, Direction: "positive", Description: "Taking responsibility for outcomes", Exercises: "Post-mortem reflections, 'what was my part?' practice, ownership journaling"},
		{ID: "responsibility_release", Name: "Responsibility Release", SpectrumID: 10, Direction: "negative", Description: "Letting go of others' consequences", Exercises: "Sorting cards: my problem vs not my problem, boundary statements"},
		{ID: "consequence_prediction", Name: "Consequence Prediction", SpectrumID: 10, Direction: "both", Description: "Foreseeing outcomes of actions", Exercises: "If-then scenarios, consequence mapping, pre-mortem practice"},
		{ID: "accountability_language", Name: "Accountability Language", SpectrumID: 10, Direction: "both", Description: "Expressing ownership without shame", Exercises: "'I did X and Y happened' scripts, repair conversations, no-blame reviews"},
		{ID: "delegation_comfort", Name: "Delegation Comfort", SpectrumID: 10, Direction: "negative", Description: "Letting others own their part", Exercises: "Shared project roles, explicit ownership splitting, trust exercises"},
	},
	// Loss Sensitivity Skills
	11: {
		{ID: "secure_attachment", Name: "Secure Attachment", SpectrumID: 11, Direction: "negative", Description: "Feeling safe without possessions as anchors", Exercises: "Object rotation, gradual lending, 'things go home' rituals"},
		{ID: "value_recognition", Name: "Value Recognition", SpectrumID: 11, Direction: "positive", Description: "Appreciating what you have while you have it", Exercises: "Gratitude naming, favourite-thing journal, appreciation rituals"},
		{ID: "transition_rituals", Name: "Transition Rituals", SpectrumID: 11, Direction: "both", Description: "Managing the process of things leaving", Exercises: "Goodbye ceremonies, photo memories, 'things go home' language"},
		{ID: "abundance_awareness", Name: "Abundance Awareness", SpectrumID: 11, Direction: "negative", Description: "Recognising that loss is not permanent", Exercises: "Seasonal cycles, sharing-and-return games, library visits"},
		{ID: "grief_processing", Name: "Grief Processing", SpectrumID: 11, Direction: "both", Description: "Healthy processing of genuine loss", Exercises: "Memory boxes, feeling naming, 'it's OK to be sad' validation"},
	},
	// Libido Skills
	12: {
		{ID: "desire_articulation", Name: "Desire Articulation", SpectrumID: 12, Direction: "both", Description: "Naming what you want without shame", Exercises: "Want-lists, choice boards, 'I want' practice in safe contexts"},
		{ID: "impulse_channeling", Name: "Impulse Channeling", SpectrumID: 12, Direction: "negative", Description: "Directing drive energy into productive outlets", Exercises: "Physical outlets, creative expression, structured challenge activities"},
		{ID: "boundary_awareness", Name: "Boundary Awareness", SpectrumID: 12, Direction: "both", Description: "Recognising others' comfort thresholds", Exercises: "Personal space games, consent practice, stop-signal activities"},
		{ID: "delayed_gratification", Name: "Delayed Gratification", SpectrumID: 12, Direction: "negative", Description: "Tolerating the gap between wanting and having", Exercises: "Waiting games, earn-then-receive sequences, marshmallow-style practice"},
		{ID: "consent_navigation", Name: "Consent Navigation", SpectrumID: 12, Direction: "both", Description: "Understanding mutual agreement", Exercises: "Ask-before-touching games, 'may I?' practice, two-yes-one-no activities"},
	},
}

// GetSpectrumByID returns the spectrum for a given ID (1-12)
func GetSpectrumByID(id int) *ETPSpectrum {
	if id < 1 || id > len(ETPSpectra) {
		return nil
	}
	return &ETPSpectra[id-1]
}

// GetSpectrumByName returns the spectrum for a given name
func GetSpectrumByName(name string) *ETPSpectrum {
	for i := range ETPSpectra {
		if ETPSpectra[i].Name == name {
			return &ETPSpectra[i]
		}
	}
	return nil
}

// SpectrumNames returns all spectrum names in order
func SpectrumNames() []string {
	names := make([]string, len(ETPSpectra))
	for i, s := range ETPSpectra {
		names[i] = s.Name
	}
	return names
}

// GetSkillsForSpectrum returns all skills for a given spectrum ID
func GetSkillsForSpectrum(spectrumID int) []SpectrumSkill {
	return SpectrumSkills[spectrumID]
}

// GetSkillByID returns a specific skill by its ID
func GetSkillByID(skillID string) *SpectrumSkill {
	for _, skills := range SpectrumSkills {
		for i := range skills {
			if skills[i].ID == skillID {
				return &skills[i]
			}
		}
	}
	return nil
}
