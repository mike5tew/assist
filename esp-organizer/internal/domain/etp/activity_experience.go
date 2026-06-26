package etp

// =============================================================================
// Activity Experience Model
// =============================================================================
//
// An Activity is NOT a single skill exercise. It is a multi-threaded event
// that activates multiple ETP spectra simultaneously, develops several skills,
// and has potential side effects (negative emotional imprints) that must be
// anticipated and mitigated.
//
// This model is the "Emotional Compiler" — it takes a human interaction and
// decomposes it into:
//   1. What skills it develops (Layer 2)
//   2. Which spectra it activates and at which pole (Layer 1)
//   3. What can go wrong (side effects / "bugs")
//   4. How to handle it when it does (mitigation strategies)
//   5. What the book page should show (the "view" of this data)
//
// The response_matrix.go provides per-spectrum guidance. This model provides
// per-ACTIVITY guidance — because real activities cross spectrum boundaries.
// =============================================================================

// Activity represents a single ToddlerOS activity — a multi-threaded event
// that crosses spectrum boundaries and has emotional side effects.
type Activity struct {
	ID          string `json:"activity_id"`
	Name        string `json:"name"`        // Parent-friendly name: "Dance Party"
	Description string `json:"description"` // What actually happens: "Put on music and dance together"
	Mode        string `json:"mode"`        // "with_you" or "while_you"
	Duration    string `json:"duration"`    // "3 min", "5 min", "10 min"
	Materials   string `json:"materials"`   // "None", "Music", "Paper and crayons"
	AgeRange    string `json:"age_range"`   // "1.0-3.0", "2.0-5.0"

	// === LAYER 2: SKILLS DEVELOPED ===
	// What foundational capacities does this activity build?
	PrimarySkills   []string `json:"primary_skills"`   // Main skills targeted: ["gross_motor", "rhythm"]
	SecondarySkills []string `json:"secondary_skills"` // Incidental skills: ["midline_crossing", "joint_attention"]

	// === LAYER 1: SPECTRUM ACTIVATION ===
	// Which ETP spectra does this activity load on, and at which pole?
	// An activity can activate multiple spectra simultaneously.
	SpectrumActivations []SpectrumActivation `json:"spectrum_activations"`

	// === SIDE EFFECTS ===
	// What can go wrong? Each side effect is a potential "bug" —
	// an unintended negative emotional imprint.
	SideEffects []SideEffect `json:"side_effects"`

	// === VENTRAL DRIVER INTERACTIONS ===
	// Which deep-brain drivers does this activity load on?
	// A child driven by "The Impressor" will experience Dance Party differently
	// from a child driven by "The Observer". This field maps the activity to
	// each driver's reward/threat conditions.
	DriverInteractions []DriverInteraction `json:"driver_interactions,omitempty"`

	// === BOOK OUTPUT SLOTS ===
	// Pre-computed content for the physical page spread.
	// Generated from the activations + side effects + response matrix.
	BookContent BookSpreadContent `json:"book_content"`
}

// SpectrumActivation describes how an activity loads on a single spectrum.
// An activity like "Dance Party" might activate 3-4 spectra at once.
type SpectrumActivation struct {
	SpectrumID   int    `json:"spectrum_id"`
	SpectrumName string `json:"spectrum_name"`

	// Which pole does this activity push toward?
	// "positive", "negative", or "both" (the activity exercises range)
	ActivationPole string `json:"activation_pole"`

	// How strongly does this activity load on this spectrum?
	// "primary" = this is the main emotional thread
	// "secondary" = activated but not the focus
	// "latent" = could be triggered by circumstances (side effect territory)
	Intensity string `json:"intensity"` // "primary", "secondary", "latent"

	// What the parent should observe on this spectrum during this activity
	ObservationPrompt string `json:"observation_prompt"`
}

// SideEffect models a potential negative emotional imprint.
// This is the "bug" that can occur when an activity triggers an unintended
// spectrum response — e.g., dance → peer observation → shame gate.
type SideEffect struct {
	ID string `json:"side_effect_id"`

	// What causes the side effect? An external event or internal state.
	Trigger string `json:"trigger"` // "peer_observation", "failure", "overstimulation"

	// Which spectrum does the trigger activate?
	TriggeredSpectrum string `json:"triggered_spectrum"` // "threat_response"
	TriggeredPole     string `json:"triggered_pole"`     // "passive" or "aggressive"

	// What is the emotional imprint if unmanaged?
	// This is the "wedge" — the lasting association the child forms.
	Imprint     string `json:"imprint"`      // "shame_gate", "avoidance_pattern", "parent_rejection"
	ImprintDesc string `json:"imprint_desc"` // Human description: "Child associates dancing with embarrassment"

	// === THE RESCUE ROUTE ===

	// What does the parent WATCH for? Observable signals that this
	// side effect is occurring.
	Signals []string `json:"signals"` // ["sudden stillness", "hiding face", "saying 'no'"]

	// What should the parent SAY? (Engineering language, not moral)
	RescueLanguage []string `json:"rescue_language"`

	// What should the parent DO? The physical action to take.
	RescueAction string `json:"rescue_action"`

	// What should the parent NOT do?
	Antipattern string `json:"antipattern"` // "Don't force them to continue"

	// Which spectrum dial to turn to resolve this?
	// This is the "counterweight" — the spectrum adjustment that
	// discharges the voltage without creating the imprint.
	Counterweight CounterweightStrategy `json:"counterweight"`
}

// CounterweightStrategy describes how to shift the spectrum dial
// to discharge the voltage from a side effect without creating an imprint.
//
// This is the "Zen Warrior" mechanism:
//  1. Expose the trigger (name it)
//  2. Acknowledge the pain (validate it)
//  3. Practice the counterweight (shift the dial)
type CounterweightStrategy struct {
	// Which spectrum to shift on
	SpectrumName string `json:"spectrum_name"`

	// Which direction to shift: move toward this pole to discharge
	ShiftDirection string `json:"shift_direction"` // "toward_independent", "toward_cohesive"

	// The specific micro-activity that performs the shift
	MicroActivity string `json:"micro_activity"` // "Move to a low-observation zone and dance alone first"

	// The naming phrase (step 1: expose the trigger)
	NamingPhrase string `json:"naming_phrase"` // "Your brain noticed people watching. That made your chest feel tight."

	// The validation phrase (step 2: acknowledge the pain)
	ValidationPhrase string `json:"validation_phrase"` // "That's your threat voltage. It's trying to protect you."

	// The invitation phrase (step 3: practice the counterweight)
	InvitationPhrase string `json:"invitation_phrase"` // "Let's find a spot where it's just us. Your voltage will come back down."
}

// BookSpreadContent holds the pre-computed content for a physical book page.
// This is the "View" layer — what the parent actually sees.
// All the complexity above compiles down to these simple slots.
type BookSpreadContent struct {
	// === LEFT PAGE: "WITH YOU" or "WHILE YOU" ===

	// The theme header — warm, parent-friendly
	ThemeTitle string `json:"theme_title"` // "Dance Party"
	ThemeIcon  string `json:"theme_icon"`  // Emoji or icon reference

	// The activity instruction — one sentence
	Instruction string `json:"instruction"` // "Put on a song you love and dance together."

	// Time + materials badge
	TimeBadge     string `json:"time_badge"`     // "3 min"
	MaterialBadge string `json:"material_badge"` // "Just music"

	// === OBSERVATION ZONE ===

	// "What to notice" — the primary observation prompt
	// Derived from the primary SpectrumActivation.ObservationPrompt
	ObservationPrompt string `json:"observation_prompt"` // "Do they copy your moves or invent their own?"

	// === RESCUE SIDEBAR ===

	// "If you notice..." — the most likely side effect, simplified
	RescueSignal   string `json:"rescue_signal"`   // "If they suddenly stop or hide..."
	RescueResponse string `json:"rescue_response"` // "That's OK. Sit down and sway gently instead."

	// === LANGUAGE CARD ===

	// One phrase to use. One phrase to avoid.
	TrySaying   string `json:"try_saying"`   // "Your body doesn't want to move right now. That's fine."
	AvoidSaying string `json:"avoid_saying"` // "Come on, it's fun!"

	// === RIGHT PAGE: INDEPENDENT / COLOURING ===

	// The "While You" companion activity
	IndependentActivity string `json:"independent_activity"` // "Draw your favourite dance move"
	IndependentPrompt   string `json:"independent_prompt"`   // A simple creative prompt

	// QR link for audio (optional)
	AudioQR string `json:"audio_qr,omitempty"` // URL for audio story/music
}

// DriverInteraction describes how a specific ventral driver experiences
// this activity. The same activity can be rewarding for one driver and
// threatening for another.
type DriverInteraction struct {
	DriverID    string `json:"driver_id"`    // References VentralDriver.ID
	PersonaName string `json:"persona_name"` // "The Impressor", "The Observer"

	// Does this activity feed the driver's reward loop, load its threat channel, or both?
	Interaction string `json:"interaction"` // "reward", "threat", "both"

	// What the parent should know about this driver during THIS activity
	DriverNote string `json:"driver_note"`
	// e.g., "An Impressor will love the audience aspect of dance. Watch for
	//        collapse if you look away."

	// Adjusted observation prompt for this driver (overrides the generic one)
	AdjustedObservation string `json:"adjusted_observation,omitempty"`
}

// =============================================================================
// Seed Data: Sample Activities
// =============================================================================

// SampleActivities provides initial activity definitions to validate the schema.
// In production, these will be generated from the response matrix PlayAvenues
// crossed with foundational skill nodes.
var SampleActivities = []Activity{
	{
		ID:          "dance_001",
		Name:        "Dance Party",
		Description: "Put on a song and dance together. Copy each other's moves, then invent new ones.",
		Mode:        "with_you",
		Duration:    "3 min",
		Materials:   "Music (phone speaker, face down)",
		AgeRange:    "1.5-5.0",

		PrimarySkills:   []string{"gross_motor", "rhythm", "joint_attention"},
		SecondarySkills: []string{"midline_crossing", "turn_taking", "emotional_labelling"},

		SpectrumActivations: []SpectrumActivation{
			{
				SpectrumID:        1,
				SpectrumName:      "social_gravity",
				ActivationPole:    "positive", // cohesive — together activity
				Intensity:         "primary",
				ObservationPrompt: "Do they want to dance with you or dance alone?",
			},
			{
				SpectrumID:        3,
				SpectrumName:      "voltage_sensitivity",
				ActivationPole:    "positive", // conductive — high emotional charge
				Intensity:         "secondary",
				ObservationPrompt: "How big is their reaction — giggly and wild, or calm and swaying?",
			},
			{
				SpectrumID:        8,
				SpectrumName:      "mirror_neuron_tuning",
				ActivationPole:    "positive", // absorbent — copying movements
				Intensity:         "secondary",
				ObservationPrompt: "Do they copy your moves or ignore them?",
			},
			{
				SpectrumID:        4,
				SpectrumName:      "threat_response",
				ActivationPole:    "both",
				Intensity:         "latent", // only activates if peer observation occurs
				ObservationPrompt: "",       // Not observed unless side effect triggers
			},
		},

		SideEffects: []SideEffect{
			{
				ID:                "dance_001_se_01",
				Trigger:           "peer_observation",
				TriggeredSpectrum: "threat_response",
				TriggeredPole:     "passive",
				Imprint:           "shame_gate",
				ImprintDesc:       "Child associates dancing with embarrassment. Withdraws from movement activities. May extend to rejecting parent-led play.",
				Signals:           []string{"sudden stillness", "hiding face", "turning away", "saying 'no' or 'stop'"},
				RescueLanguage:    []string{"Your brain noticed someone watching. That made your body want to stop."},
				RescueAction:      "Sit down. Reduce to finger dancing or hand clapping. Remove the audience.",
				Antipattern:       "Don't say 'Come on, it's fun!' or 'Nobody's watching!' — both deny the child's reality.",
				Counterweight: CounterweightStrategy{
					SpectrumName:     "social_gravity",
					ShiftDirection:   "toward_independent",
					MicroActivity:    "Move to a private space. Dance sitting down, just hands. Build back to standing when voltage drops.",
					NamingPhrase:     "Your brain saw people and put up a shield. That's your threat voltage.",
					ValidationPhrase: "That prickly feeling is real. Your body is protecting you.",
					InvitationPhrase: "Let's find a spot where it's just us. We can wiggle our fingers instead.",
				},
			},
			{
				ID:                "dance_001_se_02",
				Trigger:           "overstimulation",
				TriggeredSpectrum: "voltage_sensitivity",
				TriggeredPole:     "positive", // conductive child gets overwhelmed
				Imprint:           "sensory_flooding",
				ImprintDesc:       "Music + movement + social contact overloads a high-conductive child. Meltdown risk.",
				Signals:           []string{"hands over ears", "screaming", "running away", "aggression spike"},
				RescueLanguage:    []string{"There's a lot happening. Your voltage is very high right now."},
				RescueAction:      "Turn music off. Sit on floor. Slow breathing. Physical contact only if child initiates.",
				Antipattern:       "Don't say 'Calm down' or 'Stop being silly' — the child literally cannot regulate yet.",
				Counterweight: CounterweightStrategy{
					SpectrumName:     "voltage_sensitivity",
					ShiftDirection:   "toward_insulated",
					MicroActivity:    "Silence. Dim lights if possible. Sit together without talking. Wait.",
					NamingPhrase:     "Everything got very loud and fast. Your body had too much voltage.",
					ValidationPhrase: "That overwhelm feeling is real. Your brain needs quiet now.",
					InvitationPhrase: "Let's sit here until your body feels ready. No rush.",
				},
			},
		},

		DriverInteractions: []DriverInteraction{
			{
				DriverID:            "impressor",
				PersonaName:         "The Impressor",
				Interaction:         "reward",
				DriverNote:          "An Impressor will LOVE this. The audience (you) is right there. Watch for collapse if you look away — that's Impressor voltage, not bad behaviour.",
				AdjustedObservation: "Do they keep checking your face for a reaction? That's the external validation loop.",
			},
			{
				DriverID:            "observer",
				PersonaName:         "The Observer",
				Interaction:         "threat",
				DriverNote:          "An Observer may refuse to start. They're not shy — they're computing. Let them watch you dance for a full minute before inviting. Their body may join before their brain gives permission.",
				AdjustedObservation: "How long do they watch before moving? A long pause is processing, not resistance.",
			},
			{
				DriverID:            "performer",
				PersonaName:         "The Performer",
				Interaction:         "reward",
				DriverNote:          "A Performer will go BIG. This is perfect fuel. The risk is overstimulation (see side effect). Use the whisper game variant if voltage spikes.",
				AdjustedObservation: "Is the volume escalating? That's Performer voltage — exciting but watch the ceiling.",
			},
		},

		BookContent: BookSpreadContent{
			ThemeTitle:    "Dance Party",
			ThemeIcon:     "💃",
			Instruction:   "Put on a song you love and dance together. Copy each other's moves.",
			TimeBadge:     "3 min",
			MaterialBadge: "Just music",

			ObservationPrompt: "Do they copy your moves or invent their own?",

			RescueSignal:   "If they suddenly stop or hide...",
			RescueResponse: "That's OK. Sit down and sway gently. The dance can be just fingers.",

			TrySaying:   "Your body doesn't want to move right now. That's fine. Let's wiggle our fingers instead.",
			AvoidSaying: "Come on, it's fun! Nobody's watching!",

			IndependentActivity: "Draw your favourite dance move",
			IndependentPrompt:   "What does your dance look like on paper?",
		},
	},

	{
		ID:          "pour_001",
		Name:        "The Pouring Game",
		Description: "Two cups, one with water. Pour back and forth without spilling.",
		Mode:        "with_you",
		Duration:    "5 min",
		Materials:   "Two cups, water, towel",
		AgeRange:    "1.5-4.0",

		PrimarySkills:   []string{"fine_motor", "sustained_attention", "inhibitory_control"},
		SecondarySkills: []string{"frustration_tolerance", "task_sequencing"},

		SpectrumActivations: []SpectrumActivation{
			{
				SpectrumID:        9,
				SpectrumName:      "orderliness",
				ActivationPole:    "positive", // ordered — precise, sequential
				Intensity:         "primary",
				ObservationPrompt: "Do they pour carefully or splash happily?",
			},
			{
				SpectrumID:        6,
				SpectrumName:      "risk_tolerance",
				ActivationPole:    "both", // pouring fuller = more risk
				Intensity:         "secondary",
				ObservationPrompt: "Do they fill the cup to the top (risk seeking) or keep it safe and low (risk averse)?",
			},
			{
				SpectrumID:        4,
				SpectrumName:      "threat_response",
				ActivationPole:    "both",
				Intensity:         "latent", // only if spill occurs
				ObservationPrompt: "",
			},
		},

		SideEffects: []SideEffect{
			{
				ID:                "pour_001_se_01",
				Trigger:           "spill",
				TriggeredSpectrum: "threat_response",
				TriggeredPole:     "passive", // freeze after mistake
				Imprint:           "failure_avoidance",
				ImprintDesc:       "Child associates spilling with failure. Refuses to try pouring tasks. Fine motor development stalls.",
				Signals:           []string{"freezing", "pushing cups away", "crying", "saying 'I can't'"},
				RescueLanguage:    []string{"The water went on the table. That's what towels are for."},
				RescueAction:      "Hand them the towel. Make wiping up the game. 'Can you catch it before it reaches the edge?'",
				Antipattern:       "Don't rush to clean it up yourself — that confirms 'spilling = bad'.",
				Counterweight: CounterweightStrategy{
					SpectrumName:     "risk_tolerance",
					ShiftDirection:   "toward_seeking",
					MicroActivity:    "Pour deliberately badly. Laugh together. 'My turn to spill! Your turn!'",
					NamingPhrase:     "Your brain said 'uh oh!' when the water went. That's your caution voltage.",
					ValidationPhrase: "Spilling feels surprising. That's OK. Water dries.",
					InvitationPhrase: "Let's try again. This time, see how full you can make it before it wobbles.",
				},
			},
			{
				ID:                "pour_001_se_02",
				Trigger:           "spill",
				TriggeredSpectrum: "threat_response",
				TriggeredPole:     "aggressive", // lash out after mistake
				Imprint:           "frustration_explosion",
				ImprintDesc:       "Child throws cup or hits table. Associates frustration with aggression as default response.",
				Signals:           []string{"throwing cups", "hitting table", "screaming", "pushing adult away"},
				RescueLanguage:    []string{"Something felt like an attack. The water surprised you and your body fought back."},
				RescueAction:      "Remove cups calmly. 'We'll pause. The cups will be here when your body is ready.'",
				Antipattern:       "Don't say 'We don't throw things' — the child already knows. They couldn't stop.",
				Counterweight: CounterweightStrategy{
					SpectrumName:     "threat_response",
					ShiftDirection:   "toward_passive",
					MicroActivity:    "Pause. Count to 5 together. Then offer the cup back. If they take it, continue. If not, finished for today.",
					NamingPhrase:     "Your body went BOOM. That's your fight voltage. It fires fast.",
					ValidationPhrase: "That frustration is real. Your hands moved before your brain could catch up.",
					InvitationPhrase: "Let's count to 5. Then you decide — try again, or finished for today. Both are fine.",
				},
			},
		},

		BookContent: BookSpreadContent{
			ThemeTitle:    "The Pouring Game",
			ThemeIcon:     "🫗",
			Instruction:   "Two cups. One with water. Pour it back and forth without spilling.",
			TimeBadge:     "5 min",
			MaterialBadge: "2 cups, water, towel",

			ObservationPrompt: "Do they pour carefully or fill it to the brim?",

			RescueSignal:   "If they spill and freeze (or throw the cup)...",
			RescueResponse: "Hand them the towel. 'Can you catch the puddle before it reaches the edge?'",

			TrySaying:   "The water went on the table. That's what towels are for.",
			AvoidSaying: "Be careful! / I told you to be careful.",

			IndependentActivity: "Colour the cups — which one is full and which is empty?",
			IndependentPrompt:   "Can you draw the water moving from one cup to the other?",
		},
	},

	{
		ID:          "hide_001",
		Name:        "Find My Hands",
		Description: "Hide your hands behind your back. Child guesses which hand holds the object. Swap roles.",
		Mode:        "with_you",
		Duration:    "2 min",
		Materials:   "One small object (toy, stone, spoon)",
		AgeRange:    "1.0-3.0",

		PrimarySkills:   []string{"joint_attention", "turn_taking", "inhibitory_control"},
		SecondarySkills: []string{"verbal_interaction_density", "sustained_attention"},

		SpectrumActivations: []SpectrumActivation{
			{
				SpectrumID:        1,
				SpectrumName:      "social_gravity",
				ActivationPole:    "positive", // requires face-to-face engagement
				Intensity:         "primary",
				ObservationPrompt: "How long do they stay engaged before looking away?",
			},
			{
				SpectrumID:        2,
				SpectrumName:      "energy_directionality",
				ActivationPole:    "positive", // outward — external focus on hands
				Intensity:         "secondary",
				ObservationPrompt: "Do they point and shout, or think quietly before choosing?",
			},
			{
				SpectrumID:        5,
				SpectrumName:      "care_response",
				ActivationPole:    "both", // sharing / turn-taking
				Intensity:         "latent",
				ObservationPrompt: "",
			},
		},

		SideEffects: []SideEffect{
			{
				ID:                "hide_001_se_01",
				Trigger:           "wrong_guess",
				TriggeredSpectrum: "threat_response",
				TriggeredPole:     "passive",
				Imprint:           "failure_avoidance",
				ImprintDesc:       "Child stops guessing. Associates 'wrong' with 'bad'. Withdraws from guessing games.",
				Signals:           []string{"refusing to guess", "always pointing to same hand", "looking away"},
				RescueLanguage:    []string{"Nearly! It was hiding in this one. Your brain made a prediction — that's the skill."},
				RescueAction:      "Show both hands open. Celebrate the guess, not the result.",
				Antipattern:       "Don't say 'Wrong!' or 'Try again' — both frame guessing as pass/fail.",
				Counterweight: CounterweightStrategy{
					SpectrumName:     "risk_tolerance",
					ShiftDirection:   "toward_seeking",
					MicroActivity:    "Make it obviously easy — leave fingers slightly open. Let them 'win' to rebuild confidence.",
					NamingPhrase:     "Your brain made a guess. Sometimes guesses land, sometimes they don't.",
					ValidationPhrase: "Guessing feels risky. That's brave.",
					InvitationPhrase: "I'll make it easier. Watch my fingers...",
				},
			},
			{
				ID:                "hide_001_se_02",
				Trigger:           "won't_take_turn",
				TriggeredSpectrum: "care_response",
				TriggeredPole:     "negative", // detached — won't share the role
				Imprint:           "turn_taking_refusal",
				ImprintDesc:       "Child always wants to guess but won't hide. Avoids the vulnerability of being 'the one who hides'.",
				Signals:           []string{"'Your turn!' refused", "grabbing both hands", "tantrums when roles swap"},
				RescueLanguage:    []string{"You love choosing. Hiding is a different kind of power — you get to trick me."},
				RescueAction:      "Model it first. 'Watch — I'll hide it, and you try to trick me. Then we swap.'",
				Antipattern:       "Don't force the swap. Offer it. If refused, play their way and try again next time.",
				Counterweight: CounterweightStrategy{
					SpectrumName:     "care_response",
					ShiftDirection:   "toward_nurturing",
					MicroActivity:    "Frame hiding as 'giving' — 'You get to give me a surprise.' Reframes turn-taking as generosity.",
					NamingPhrase:     "You want to keep choosing. That feels safe.",
					ValidationPhrase: "Swapping feels like losing control. That's hard.",
					InvitationPhrase: "What if you hide it AND choose which hand? You're still in charge.",
				},
			},
		},

		BookContent: BookSpreadContent{
			ThemeTitle:    "Find My Hands",
			ThemeIcon:     "🤲",
			Instruction:   "Hide a small object behind your back. Which hand is it in?",
			TimeBadge:     "2 min",
			MaterialBadge: "Any small object",

			ObservationPrompt: "Do they guess quickly or think about it first?",

			RescueSignal:   "If they stop guessing or refuse to swap roles...",
			RescueResponse: "Make it easy. Leave your fingers slightly open. Celebrate the guess, not the answer.",

			TrySaying:   "Your brain made a prediction. That's the skill — not getting it right.",
			AvoidSaying: "Wrong! Try again.",

			IndependentActivity: "Which hand is the teddy hiding in? Circle the right one.",
			IndependentPrompt:   "Colour the hands and hide a sticker under one.",
		},
	},
}

// GetActivityByID returns an activity by its ID
func GetActivityByID(id string) *Activity {
	for i := range SampleActivities {
		if SampleActivities[i].ID == id {
			return &SampleActivities[i]
		}
	}
	return nil
}

// GetActivitiesForSpectrum returns all activities that activate a given spectrum
func GetActivitiesForSpectrum(spectrumID int) []Activity {
	var result []Activity
	for _, a := range SampleActivities {
		for _, sa := range a.SpectrumActivations {
			if sa.SpectrumID == spectrumID {
				result = append(result, a)
				break
			}
		}
	}
	return result
}

// GetActivitiesForSkill returns all activities that develop a given skill
func GetActivitiesForSkill(skill string) []Activity {
	var result []Activity
	for _, a := range SampleActivities {
		for _, s := range a.PrimarySkills {
			if s == skill {
				result = append(result, a)
				break
			}
		}
		for _, s := range a.SecondarySkills {
			if s == skill {
				result = append(result, a)
				break
			}
		}
	}
	return result
}

// GetActivitiesByMode returns all activities matching a mode ("with_you" or "while_you")
func GetActivitiesByMode(mode string) []Activity {
	var result []Activity
	for _, a := range SampleActivities {
		if a.Mode == mode {
			result = append(result, a)
		}
	}
	return result
}
