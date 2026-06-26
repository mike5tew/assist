package etp

// =============================================================================
// Ventral Emotional Drivers — "The Ghost in the Machine"
// =============================================================================
//
// A Ventral Driver is the deep-brain "WHY" behind a toddler's behaviour.
// It is NOT a single spectrum reading. It is a PATTERN — a characteristic
// combination of 2–3 spectra that form a persistent reward loop.
//
// Where the spectra tell us WHAT is happening (slider positions), the
// Ventral Driver tells us WHY — the motivational force that pulls the
// sliders into that configuration in the first place.
//
// Example: The "Impressor" archetype
//   - social_gravity → high cohesive (NEEDS the audience)
//   - mirror_neuron_tuning → high absorbent (READS the audience's reaction)
//   - energy_directionality → high outward (PROJECTS toward the audience)
//   - Reward loop: "I exist because you see me"
//   - Risk: External validation becomes the only power source
//
// ToddlerOS does NOT "fix" ventral drivers. They are power sources.
// The goal is to build a HEAT SHIELD (metacognitive distance) between the
// driver and the executive function, so the child can use the energy
// without being consumed by it.
//
// The Zen Warrior formula:
//   1. NAME the driver ("You're using your People Power right now")
//   2. VALIDATE the energy ("It feels big when I look away, doesn't it?")
//   3. STRESS-TEST in safe doses ("Can you do one more turn just for YOUR eyes?")
//   4. DEBRIEF the discovery ("How did that feel in your tummy?")
//
// Over time, the child learns: "I HAVE this feeling, but I am NOT this feeling."
// =============================================================================

// VentralDriver represents a persistent motivational pattern — the deep-brain
// "reward loop" that pulls a toddler's spectrum sliders into a characteristic
// configuration. Each driver is a combination of 2–3 dominant spectra that
// together form a recognisable persona.
type VentralDriver struct {
	ID   string `json:"driver_id"`
	Name string `json:"name"` // Engineering name: "external_validation_loop"

	// === PARENT-FACING IDENTITY ===

	// The persona name — a warm, non-clinical label for the book
	PersonaName string `json:"persona_name"` // "The Impressor"
	PersonaIcon string `json:"persona_icon"` // Emoji for the book spread

	// One-sentence description the parent can recognise
	ParentDescription string `json:"parent_description"`
	// e.g., "Your child lights up when you watch. They perform, show, present.
	//        The audience IS the fuel."

	// What the child would say if they had the words
	ChildVoice string `json:"child_voice"`
	// e.g., "I exist because you see me."

	// === SPECTRUM SIGNATURE ===

	// The 2–3 spectra that combine to form this driver.
	// Each has a dominant pole — the direction the driver pulls.
	SpectrumSignature []DriverSpectrum `json:"spectrum_signature"`

	// The reward condition — what "success" feels like to this driver
	RewardCondition string `json:"reward_condition"`
	// e.g., "A positive reflection in the adult's eyes"

	// The threat condition — what triggers the driver's anxiety
	ThreatCondition string `json:"threat_condition"`
	// e.g., "Being unobserved, ignored, or receiving neutral/negative feedback"

	// === THE HEAT SHIELD ===

	// The Naming Phrase — how the parent identifies this driver in the moment
	NamingPhrase string `json:"naming_phrase"`
	// e.g., "You're using your People Power right now. You really want me to see this."

	// The Moderation Move — the Zen Warrior micro-intervention
	ModerationMove string `json:"moderation_move"`
	// e.g., "It feels big when I look away, doesn't it? That's your Social Gravity
	//        pulling hard. Let's see if we can do one more turn just for YOUR eyes,
	//        and see how that feels in your tummy."

	// The Stress Test — a specific activity that safely loads the driver's
	// anxiety channel, building tolerance
	StressTest DriverStressTest `json:"stress_test"`

	// The Long Game — what this driver becomes when well-regulated vs. unregulated
	RegulatedOutcome   string `json:"regulated_outcome"`   // "A charismatic leader who can perform without needing approval"
	UnregulatedOutcome string `json:"unregulated_outcome"` // "A people-pleaser who collapses without external validation"

	// === INTERACTION WITH OTHER DRIVERS ===

	// Drivers that amplify this one (positive feedback loop risk)
	Amplifiers []string `json:"amplifiers"` // driver IDs
	// Drivers that naturally moderate this one (counterbalance)
	Counterbalances []string `json:"counterbalances"` // driver IDs
}

// DriverSpectrum describes one spectrum's role within a ventral driver pattern.
type DriverSpectrum struct {
	SpectrumID   int    `json:"spectrum_id"`
	SpectrumName string `json:"spectrum_name"`
	DominantPole string `json:"dominant_pole"` // "positive" or "negative"
	Role         string `json:"role"`          // "primary_fuel", "amplifier", "filter"
	Mechanism    string `json:"mechanism"`     // How this spectrum serves the driver
}

// DriverStressTest is a specific activity designed to safely load the driver's
// anxiety channel. The child practises tolerating the "threat condition" in
// micro-doses, building metacognitive distance.
type DriverStressTest struct {
	ActivityName string `json:"activity_name"` // "The Blind Chef"
	Description  string `json:"description"`   // What happens
	Duration     string `json:"duration"`      // "30 sec", "2 min"
	// What the parent removes (the reward the driver craves)
	DeprivationTarget string `json:"deprivation_target"` // "parent's gaze"
	// What the child must do without it
	Challenge string `json:"challenge"` // "Keep playing for 30 seconds without looking at me"
	// The debrief question
	DebriefPrompt string `json:"debrief_prompt"` // "How did that feel in your tummy?"
}

// =============================================================================
// The Driver Map — 7 Archetypal Ventral Drivers
// =============================================================================
//
// These are NOT personality types. A child can express multiple drivers,
// and drivers shift with context and development. They are PATTERNS that
// help parents recognise the "why" behind recurring behaviours.
//
// The 7 archetypes cover the most common deep-brain reward loops observed
// in 0–5 development. Each maps to a recognisable "persona" that parents
// can spot without any technical knowledge.
// =============================================================================

var VentralDriverMap = []VentralDriver{
	// =========================================================================
	// 1. THE IMPRESSOR
	// =========================================================================
	{
		ID:                "impressor",
		Name:              "external_validation_loop",
		PersonaName:       "The Impressor",
		PersonaIcon:       "⭐",
		ParentDescription: "Your child lights up when you watch. They perform, show, present. The audience IS the fuel.",
		ChildVoice:        "I exist because you see me.",

		SpectrumSignature: []DriverSpectrum{
			{
				SpectrumID:   1,
				SpectrumName: "social_gravity",
				DominantPole: "positive", // cohesive
				Role:         "primary_fuel",
				Mechanism:    "Social connection is the reward. Being seen = being real.",
			},
			{
				SpectrumID:   8,
				SpectrumName: "mirror_neuron_tuning",
				DominantPole: "positive", // absorbent
				Role:         "amplifier",
				Mechanism:    "Reads the audience's emotional state with high fidelity. A positive reflection amplifies; a neutral one threatens.",
			},
			{
				SpectrumID:   2,
				SpectrumName: "energy_directionality",
				DominantPole: "positive", // outward
				Role:         "filter",
				Mechanism:    "Energy projects outward toward the audience. Internal experience is secondary to the reaction it generates.",
			},
		},

		RewardCondition: "A positive reflection in the adult's eyes. Applause, attention, delight.",
		ThreatCondition: "Being unobserved, ignored, or receiving neutral/negative feedback. The gaze turning away.",

		NamingPhrase:   "You're using your People Power right now. You really want me to see this.",
		ModerationMove: "It feels big when I look away, doesn't it? That's your Social Gravity pulling hard. Let's see if we can do one more turn just for YOUR eyes, and see how that feels in your tummy.",

		StressTest: DriverStressTest{
			ActivityName:      "The Invisible Artist",
			Description:       "Child draws/builds while the parent looks away for 30 seconds. Parent then looks and asks 'Tell me about it' (not 'That's great!').",
			Duration:          "30 sec deprivation, then 2 min shared",
			DeprivationTarget: "Parent's gaze and approval",
			Challenge:         "Keep creating without checking if I'm watching",
			DebriefPrompt:     "When I wasn't looking, did you keep going? How did that feel?",
		},

		RegulatedOutcome:   "A charismatic communicator who can perform brilliantly AND sit comfortably alone. Uses People Power as a tool, not a life-support system.",
		UnregulatedOutcome: "A people-pleaser who collapses without external validation. Performance anxiety. Identity dependent on audience reaction.",

		Amplifiers:      []string{"performer"},            // performer + impressor = runaway validation loop
		Counterbalances: []string{"explorer", "sentinel"}, // inward focus + safety-seeking moderate the need for audience
	},

	// =========================================================================
	// 2. THE EXPLORER
	// =========================================================================
	{
		ID:                "explorer",
		Name:              "sensory_discovery_loop",
		PersonaName:       "The Explorer",
		PersonaIcon:       "🔍",
		ParentDescription: "Your child wanders, touches, tests. They don't need you to watch — they need the world to respond. The environment IS the conversation.",
		ChildVoice:        "I understand by touching.",

		SpectrumSignature: []DriverSpectrum{
			{
				SpectrumID:   6,
				SpectrumName: "risk_tolerance",
				DominantPole: "positive", // seeking
				Role:         "primary_fuel",
				Mechanism:    "Risk creates excitement voltage. The unknown is rewarding, not threatening.",
			},
			{
				SpectrumID:   2,
				SpectrumName: "energy_directionality",
				DominantPole: "negative", // inward
				Role:         "amplifier",
				Mechanism:    "Processes the sensory feedback internally. Doesn't need to narrate the discovery — just absorb it.",
			},
			{
				SpectrumID:   1,
				SpectrumName: "social_gravity",
				DominantPole: "negative", // independent
				Role:         "filter",
				Mechanism:    "Other people are background. The object/environment is foreground.",
			},
		},

		RewardCondition: "Sensory feedback from the environment. Something moves, changes, or reveals a property.",
		ThreatCondition: "Being pulled away from investigation. Interruption. 'Come here!' during deep exploration.",

		NamingPhrase:   "Your Discovery Brain is switched on. You're figuring something out.",
		ModerationMove: "You really want to keep going. That's your Explorer energy. But right now we need to pause. Can you put the question in your pocket and come back to it?",

		StressTest: DriverStressTest{
			ActivityName:      "The Pause Button",
			Description:       "During deep exploration, parent says 'Pause! Tell me one thing about it.' Child practises stopping mid-investigation and verbalising.",
			Duration:          "10 sec pause, then resume",
			DeprivationTarget: "Uninterrupted exploration",
			Challenge:         "Stop exploring for 10 seconds and tell me what you found",
			DebriefPrompt:     "Was it hard to stop? What did your hands want to do?",
		},

		RegulatedOutcome:   "A deep thinker who can investigate intensely AND disengage when needed. Uses curiosity as a power source, not an escape from connection.",
		UnregulatedOutcome: "A child who disappears into objects and resists all social engagement. Exploration becomes avoidance. Relationships suffer.",

		Amplifiers:      []string{"sentinel"},               // both resist social pull — can reinforce isolation
		Counterbalances: []string{"impressor", "connector"}, // social reward + care response draw the explorer back to people
	},

	// =========================================================================
	// 3. THE SENTINEL
	// =========================================================================
	{
		ID:                "sentinel",
		Name:              "safety_maintenance_loop",
		PersonaName:       "The Sentinel",
		PersonaIcon:       "🛡️",
		ParentDescription: "Your child watches before acting. They check the room, the rules, the exits. Change is the enemy. Sameness IS safety.",
		ChildVoice:        "I'm safe when nothing changes.",

		SpectrumSignature: []DriverSpectrum{
			{
				SpectrumID:   4,
				SpectrumName: "threat_response",
				DominantPole: "negative", // passive — freeze/scan before acting
				Role:         "primary_fuel",
				Mechanism:    "Threat detection is always running. The sentinel scans for danger before engaging.",
			},
			{
				SpectrumID:   9,
				SpectrumName: "orderliness",
				DominantPole: "positive", // ordered — routine is the shield
				Role:         "amplifier",
				Mechanism:    "Structure reduces threat surface. If everything is in its place, nothing can surprise.",
			},
			{
				SpectrumID:   6,
				SpectrumName: "risk_tolerance",
				DominantPole: "negative", // averse — risk = voltage spike
				Role:         "filter",
				Mechanism:    "Risk creates anxiety voltage, not excitement. The safe path is always preferred.",
			},
		},

		RewardCondition: "Absence of change. Predictability confirmed. Same routine, same rules, same outcome.",
		ThreatCondition: "Unexpected change. New environment, new person, broken routine, violated rule.",

		NamingPhrase:   "Your Safety Brain is doing a big scan right now. It's checking everything.",
		ModerationMove: "Something changed and your brain said 'ALERT!' That's your Sentinel energy — it's trying to protect you. Let's look at what's different together. Your brain can stand down when it sees there's no real danger.",

		StressTest: DriverStressTest{
			ActivityName:      "The Tiny Surprise",
			Description:       "During a familiar routine, parent introduces ONE small change (different coloured cup, different route to the table). Child practises tolerating the deviation.",
			Duration:          "The whole routine (with one change)",
			DeprivationTarget: "Absolute predictability",
			Challenge:         "Can you spot what's different? Is it dangerous or just different?",
			DebriefPrompt:     "The thing that changed — was it scary-different or just new-different?",
		},

		RegulatedOutcome:   "A careful, strategic thinker who assesses situations accurately AND can adapt when plans change. Uses vigilance as a superpower, not a cage.",
		UnregulatedOutcome: "An anxious child who cannot tolerate any change. Rigidity. Meltdowns at transitions. Control-seeking behaviour that escalates with age.",

		Amplifiers:      []string{"explorer"},           // both prefer solitary processing — can reinforce withdrawal
		Counterbalances: []string{"performer", "rebel"}, // spontaneity + risk-seeking crack the rigidity gently
	},

	// =========================================================================
	// 4. THE CONNECTOR
	// =========================================================================
	{
		ID:                "connector",
		Name:              "care_attachment_loop",
		PersonaName:       "The Connector",
		PersonaIcon:       "🤝",
		ParentDescription: "Your child brings you things, checks on you, notices when you're sad. Other people's feelings ARE their business. Love is the operating system.",
		ChildVoice:        "Are you OK? Can I help?",

		SpectrumSignature: []DriverSpectrum{
			{
				SpectrumID:   5,
				SpectrumName: "care_response",
				DominantPole: "positive", // nurturing
				Role:         "primary_fuel",
				Mechanism:    "Vulnerability in others triggers care voltage. Helping IS the reward.",
			},
			{
				SpectrumID:   8,
				SpectrumName: "mirror_neuron_tuning",
				DominantPole: "positive", // absorbent
				Role:         "amplifier",
				Mechanism:    "Absorbs others' emotional states with high fidelity. Feels what they feel.",
			},
			{
				SpectrumID:   1,
				SpectrumName: "social_gravity",
				DominantPole: "positive", // cohesive
				Role:         "filter",
				Mechanism:    "Proximity to people is the default. Being alone = being without purpose.",
			},
		},

		RewardCondition: "Making someone feel better. Being needed. The emotional temperature of the room dropping because they helped.",
		ThreatCondition: "Being unable to help. Witnessing distress they can't fix. Being told 'I'm fine' when they can feel the person isn't.",

		NamingPhrase:   "Your Caring Brain noticed something. You want to help!",
		ModerationMove: "You can feel that something's wrong. That's your mirror power — very strong. But this feeling belongs to ME, not to you. You don't have to carry it. Can you send me a kind thought instead of trying to fix it?",

		StressTest: DriverStressTest{
			ActivityName:      "The Kind Thought",
			Description:       "Parent pretends to be mildly sad. Child practises sending a 'kind thought' instead of physically trying to fix it. Builds emotional boundary.",
			Duration:          "1 min",
			DeprivationTarget: "Permission to physically fix the problem",
			Challenge:         "Can you send me a kind thought with your eyes instead of your hands?",
			DebriefPrompt:     "Did you still feel like helping? What happened when you just sent the thought?",
		},

		RegulatedOutcome:   "A deeply empathic person who can sit with others' pain without drowning in it. Uses care as a compass, not a compulsion.",
		UnregulatedOutcome: "A child who takes on everyone's emotions. Burnout. Boundary collapse. Loses self in service to others. Cannot tolerate others' unhappiness.",

		Amplifiers:      []string{"impressor"},         // both need social connection — can become codependent
		Counterbalances: []string{"explorer", "rebel"}, // independence + self-assertion build healthy boundaries
	},

	// =========================================================================
	// 5. THE REBEL
	// =========================================================================
	{
		ID:                "rebel",
		Name:              "autonomy_assertion_loop",
		PersonaName:       "The Rebel",
		PersonaIcon:       "⚡",
		ParentDescription: "Your child pushes back. 'No' is their favourite word. They test every boundary, not to be difficult, but because THEY need to be the one who decides. Control IS identity.",
		ChildVoice:        "I decide. Not you.",

		SpectrumSignature: []DriverSpectrum{
			{
				SpectrumID:   4,
				SpectrumName: "threat_response",
				DominantPole: "positive", // aggressive — confrontation over compliance
				Role:         "primary_fuel",
				Mechanism:    "Being told 'no' triggers confrontation voltage. Compliance feels like submission.",
			},
			{
				SpectrumID:   7,
				SpectrumName: "integrity_logic",
				DominantPole: "negative", // relativistic — MY rules, not yours
				Role:         "amplifier",
				Mechanism:    "External rules feel arbitrary. The rebel's internal logic is the only valid authority.",
			},
			{
				SpectrumID:   1,
				SpectrumName: "social_gravity",
				DominantPole: "negative", // independent — don't need your approval
				Role:         "filter",
				Mechanism:    "Social approval is not the currency. Autonomy is.",
			},
		},

		RewardCondition: "Agency confirmed. 'I chose this.' Even choosing the wrong thing feels better than being told the right thing.",
		ThreatCondition: "Loss of control. Being overridden, forced, or ignored. Being physically moved.",

		NamingPhrase:   "Your Boss Brain just switched on. You REALLY want to be in charge right now.",
		ModerationMove: "I hear your Boss Brain. It's very loud right now. Here's the deal: you can choose THIS or THIS. Both are fine. But doing nothing isn't one of the choices.",

		StressTest: DriverStressTest{
			ActivityName:      "The Two Doors",
			Description:       "Parent gives two acceptable choices. Child practises choosing within constraints instead of demanding a third option. Builds agency WITHIN boundaries.",
			Duration:          "30 sec decision window",
			DeprivationTarget: "Unlimited choice / total control",
			Challenge:         "Pick one of these two. Both are real choices.",
			DebriefPrompt:     "You chose! How did it feel to pick? Was it hard when you couldn't do the other thing?",
		},

		RegulatedOutcome:   "A strong-willed leader who can assert boundaries AND accept constraints. Uses defiance as discernment, not destruction.",
		UnregulatedOutcome: "A child in permanent opposition. Every instruction is a battle. Oppositional behaviour escalates. Authority figures become enemies.",

		Amplifiers:      []string{"performer"},             // both resist external control — can become explosive
		Counterbalances: []string{"connector", "sentinel"}, // care response + safety-seeking provide reasons to cooperate
	},

	// =========================================================================
	// 6. THE PERFORMER
	// =========================================================================
	{
		ID:                "performer",
		Name:              "intensity_expression_loop",
		PersonaName:       "The Performer",
		PersonaIcon:       "🎭",
		ParentDescription: "Your child is LOUD. Big emotions, big movements, big reactions. They don't do anything at half-volume. The intensity IS the message.",
		ChildVoice:        "If you can't feel it, it doesn't count.",

		SpectrumSignature: []DriverSpectrum{
			{
				SpectrumID:   3,
				SpectrumName: "voltage_sensitivity",
				DominantPole: "positive", // conductive — high emotional current
				Role:         "primary_fuel",
				Mechanism:    "Emotions run at full voltage. The signal is always loud.",
			},
			{
				SpectrumID:   2,
				SpectrumName: "energy_directionality",
				DominantPole: "positive", // outward — expression is the outlet
				Role:         "amplifier",
				Mechanism:    "Internal emotional voltage MUST be discharged externally. Containment is painful.",
			},
			{
				SpectrumID:   6,
				SpectrumName: "risk_tolerance",
				DominantPole: "positive", // seeking — goes big
				Role:         "filter",
				Mechanism:    "Cautious expression feels like suppression. The big move is always preferred.",
			},
		},

		RewardCondition: "Full emotional expression received and matched. 'I showed you how I feel and you felt it too.'",
		ThreatCondition: "'Calm down.' 'Use your indoor voice.' 'Stop being so dramatic.' — any instruction to reduce voltage.",

		NamingPhrase:   "Your Volume Dial is turned right up. Your feelings are VERY big right now.",
		ModerationMove: "I can see how big this feels. Your voltage is really high. I'm not going to tell you to calm down — but let's find a way to let some of that energy out without it hurting anyone. Can we stamp our feet together?",

		StressTest: DriverStressTest{
			ActivityName:      "The Whisper Game",
			Description:       "An exciting activity conducted entirely in whispers. Child practises containing high voltage in a low-output channel.",
			Duration:          "2 min",
			DeprivationTarget: "Volume as an outlet",
			Challenge:         "Can you tell me the exciting thing... but whisper it?",
			DebriefPrompt:     "You had BIG feelings but tiny words. Did the feelings still fit?",
		},

		RegulatedOutcome:   "An expressive, passionate person who can communicate intensity AND modulate output. Uses emotional power as fuel, not as a weapon.",
		UnregulatedOutcome: "Emotional dysregulation. Every feeling at max volume. Meltdowns. Adults start managing BY avoidance. Child learns that feelings = danger.",

		Amplifiers:      []string{"rebel"},                // both run hot — can create explosive combinations
		Counterbalances: []string{"sentinel", "explorer"}, // structure + inward focus teach voltage management
	},

	// =========================================================================
	// 7. THE OBSERVER
	// =========================================================================
	{
		ID:                "observer",
		Name:              "processing_depth_loop",
		PersonaName:       "The Observer",
		PersonaIcon:       "👁️",
		ParentDescription: "Your child watches. They don't join in immediately. They study the room, the people, the game. They're not shy — they're COMPUTING. Depth IS the pace.",
		ChildVoice:        "I need to understand before I move.",

		SpectrumSignature: []DriverSpectrum{
			{
				SpectrumID:   2,
				SpectrumName: "energy_directionality",
				DominantPole: "negative", // inward — processes internally first
				Role:         "primary_fuel",
				Mechanism:    "External stimulation is raw data to be processed, not a call to action. Action comes after understanding.",
			},
			{
				SpectrumID:   8,
				SpectrumName: "mirror_neuron_tuning",
				DominantPole: "negative", // selective — watches but doesn't absorb
				Role:         "amplifier",
				Mechanism:    "Observes emotions as data points, not as shared experience. Analytical, not empathic.",
			},
			{
				SpectrumID:   3,
				SpectrumName: "voltage_sensitivity",
				DominantPole: "negative", // insulated — emotional noise filtered out
				Role:         "filter",
				Mechanism:    "Low emotional voltage means the analytical process isn't disrupted by feelings.",
			},
		},

		RewardCondition: "Understanding achieved. The pattern is clear. 'I see how this works now.'",
		ThreatCondition: "Being forced to act before processing is complete. 'Just try it!' 'Stop watching and join in!'",

		NamingPhrase:   "Your Thinking Brain is running. You're watching to understand.",
		ModerationMove: "You're watching. That's YOUR way of learning. But I wonder — could your body try it while your brain is still figuring it out? Sometimes your hands know before your head does.",

		StressTest: DriverStressTest{
			ActivityName:      "The Body First Game",
			Description:       "A physical activity where acting comes before understanding (jumping into a pile of leaves, rolling down a hill). Child practises action-before-analysis.",
			Duration:          "30 sec (the action), then unlimited analysis time",
			DeprivationTarget: "Full understanding before acting",
			Challenge:         "Jump first, think second. Ready? GO!",
			DebriefPrompt:     "Your body did it before your brain was ready. What did your body know that your brain didn't?",
		},

		RegulatedOutcome:   "A deep analyst who can process AND act. Uses observation as preparation, not perpetual delay. Comfortable with incomplete understanding.",
		UnregulatedOutcome: "Paralysis by analysis. Refuses to participate. Labelled 'shy' or 'slow' when actually processing deeply. Misses opportunities. Social isolation by default.",

		Amplifiers:      []string{"sentinel"},           // both prefer safety of non-action — can reinforce paralysis
		Counterbalances: []string{"performer", "rebel"}, // high voltage + agency push the observer into action
	},
}

// =============================================================================
// Driver Detection Helpers
// =============================================================================

// GetDriverByID returns a ventral driver by ID
func GetDriverByID(driverID string) *VentralDriver {
	for i := range VentralDriverMap {
		if VentralDriverMap[i].ID == driverID {
			return &VentralDriverMap[i]
		}
	}
	return nil
}

// GetDriverByPersona returns a ventral driver by persona name (case-insensitive match)
func GetDriverByPersona(name string) *VentralDriver {
	for i := range VentralDriverMap {
		if VentralDriverMap[i].PersonaName == name {
			return &VentralDriverMap[i]
		}
	}
	return nil
}

// GetDriversForSpectrum returns all drivers that use a given spectrum in their signature
func GetDriversForSpectrum(spectrumID int) []VentralDriver {
	var result []VentralDriver
	for _, d := range VentralDriverMap {
		for _, s := range d.SpectrumSignature {
			if s.SpectrumID == spectrumID {
				result = append(result, d)
				break
			}
		}
	}
	return result
}

// DetectLikelyDrivers takes a set of spectrum observations and returns
// drivers ranked by how closely the observations match the driver's signature.
// Each observation is a spectrum_id → value (-2.0 to +2.0) where negative
// means the negative pole and positive means the positive pole.
//
// This is used by the parent app: after 3+ observation cycles, the system
// can suggest which drivers are most active. The parent never sees "driver
// detection" — they see "Your child seems to light up when..." framing.
func DetectLikelyDrivers(observations map[int]float64) []DriverMatch {
	var matches []DriverMatch

	for _, driver := range VentralDriverMap {
		score := 0.0
		maxScore := 0.0

		for _, sig := range driver.SpectrumSignature {
			// Each spectrum in the signature contributes to the match
			weight := 1.0
			switch sig.Role {
			case "primary_fuel":
				weight = 3.0
			case "amplifier":
				weight = 2.0
			case "filter":
				weight = 1.0
			}
			maxScore += weight * 2.0 // max observation is ±2.0

			if obs, ok := observations[sig.SpectrumID]; ok {
				// Check if observation aligns with the driver's dominant pole
				if sig.DominantPole == "positive" && obs > 0 {
					score += weight * obs
				} else if sig.DominantPole == "negative" && obs < 0 {
					score += weight * (-obs) // flip sign
				}
				// Opposing observation = 0 contribution (not negative — absence isn't counter-evidence)
			}
		}

		if maxScore > 0 {
			normalised := score / maxScore // 0.0 to 1.0
			matches = append(matches, DriverMatch{
				DriverID:    driver.ID,
				PersonaName: driver.PersonaName,
				Confidence:  normalised,
			})
		}
	}

	// Sort by confidence descending
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].Confidence > matches[i].Confidence {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	return matches
}

// DriverMatch represents the confidence that a ventral driver is active,
// based on spectrum observations.
type DriverMatch struct {
	DriverID    string  `json:"driver_id"`
	PersonaName string  `json:"persona_name"`
	Confidence  float64 `json:"confidence"` // 0.0 to 1.0
}
