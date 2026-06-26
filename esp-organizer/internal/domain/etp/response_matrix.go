package etp

// ResponseMatrixEntry captures how a specific ETP setting manifests barriers
// and what language/strategies adults should use in response.
// This is the core of "The Architecture of Emotional Readiness."
type ResponseMatrixEntry struct {
	SpectrumID   int    `json:"spectrum_id"`
	SpectrumName string `json:"spectrum_name"`
	Setting      string `json:"setting"`      // e.g., "independent", "cohesive"
	SettingEnd   string `json:"setting_end"`  // "negative" or "positive"
	BarrierType  string `json:"barrier_type"` // "verbalising", "starting", "mistakes"

	// What the barrier looks like in practice
	BarrierDescription string `json:"barrier_description"`

	// Language the adult should AVOID (moral framework)
	AvoidLanguage []string `json:"avoid_language"`

	// Language the adult should USE (engineering framework)
	RecommendedLanguage []string `json:"recommended_language"`

	// Play/exploration strategies for developing resilience
	PlayAvenues []string `json:"play_avenues"`
}

// UnifiedLanguageShift maps old moral language to new engineering language.
type UnifiedLanguageShift struct {
	OldMoral       string `json:"old_moral"`
	NewEngineering string `json:"new_engineering"`
}

// The three core barriers that the response matrix addresses
var CoreBarriers = []struct {
	Name    string
	Problem string
	ETPLens string
}{
	{
		Name:    "verbalising",
		Problem: "\"I don't know what I think until I say it\" vs \"I know what I think but can't find the words\"",
		ETPLens: "Energy Directionality + Mirror Neuron Tuning",
	},
	{
		Name:    "starting",
		Problem: "Inertia, overwhelm, perfectionism, fear",
		ETPLens: "Risk Tolerance + Threat Response + Orderliness",
	},
	{
		Name:    "mistakes",
		Problem: "Shutdown, deflection, self-criticism, blame",
		ETPLens: "Integrity Logic + Care Response + Threat Response",
	},
}

// UnifiedLanguageShifts — the pattern across all spectra
var UnifiedLanguageShifts = []UnifiedLanguageShift{
	{OldMoral: "You are...", NewEngineering: "Your setting is..."},
	{OldMoral: "You should...", NewEngineering: "Your brain is telling you..."},
	{OldMoral: "That was wrong", NewEngineering: "That didn't work — let's look at why"},
	{OldMoral: "Try harder", NewEngineering: "What's blocking you right now?"},
	{OldMoral: "Stop it", NewEngineering: "I notice you're in X setting — is that helping?"},
}

// PlayFirstPrinciple — the adult's job for PrimaryOS/ToddlerOS
var PlayFirstPrinciple = []string{
	"Observe the setting without judgment",
	"Name it neutrally (\"You're in Independent mode right now\")",
	"Design the play that gently expands range",
	"Debrief in engineering language (\"When you had to wait, what happened in your brain?\")",
}

// ResponseMatrix contains all 18 entries (9 spectra × 2 settings)
var ResponseMatrix = []ResponseMatrixEntry{
	// ============================================================
	// 1. SOCIAL GRAVITY
	// ============================================================
	{
		SpectrumID:         1,
		SpectrumName:       "social_gravity",
		Setting:            "independent",
		SettingEnd:         "negative",
		BarrierType:        "starting",
		BarrierDescription: "Won't ask for help when stuck",
		AvoidLanguage: []string{
			"You should have asked for help earlier",
		},
		RecommendedLanguage: []string{
			"Your battery was running low. Next time, you could signal me and I'll come to you — no need to come find me.",
		},
		PlayAvenues: []string{
			"Parallel play with gradual invitation",
			"Let them work alone in the same space",
			"Occasionally wonder aloud near them: \"I'm wondering if...\" — invitation without demand",
		},
	},
	{
		SpectrumID:         1,
		SpectrumName:       "social_gravity",
		Setting:            "cohesive",
		SettingEnd:         "positive",
		BarrierType:        "starting",
		BarrierDescription: "Can't start without someone else",
		AvoidLanguage: []string{
			"You need to learn to work alone",
		},
		RecommendedLanguage: []string{
			"You work best with a thinking partner. Let's find you one, then I'll check back in 10 minutes when you're rolling.",
		},
		PlayAvenues: []string{
			"Paired exploration tasks where the goal is exploration, not the product",
			"Gradually increase the \"alone time\" within the paired context",
		},
	},

	// ============================================================
	// 2. ENERGY DIRECTIONALITY
	// ============================================================
	{
		SpectrumID:         2,
		SpectrumName:       "energy_directionality",
		Setting:            "inward",
		SettingEnd:         "negative",
		BarrierType:        "verbalising",
		BarrierDescription: "\"I don't know\" when asked directly",
		AvoidLanguage: []string{
			"You must know, you just did it",
		},
		RecommendedLanguage: []string{
			"You're processing. Take a minute, then come find me when you've got words.",
			"Draw me a picture of what happened.",
		},
		PlayAvenues: []string{
			"Build in \"processing pauses\" before any verbal response",
			"\"Let's all think for 30 seconds before anyone speaks\"",
			"Normalise the pause",
		},
	},
	{
		SpectrumID:         2,
		SpectrumName:       "energy_directionality",
		Setting:            "outward",
		SettingEnd:         "positive",
		BarrierType:        "verbalising",
		BarrierDescription: "Talks over others, can't focus until thoughts are spoken",
		AvoidLanguage: []string{
			"Stop interrupting",
		},
		RecommendedLanguage: []string{
			"I can see you're thinking. Jot that thought down so you don't lose it, then I'll come back to you.",
		},
		PlayAvenues: []string{
			"Use talking objects, thinking partners, \"parking lots\" for ideas",
			"Give them a legitimate channel for external processing",
		},
	},

	// ============================================================
	// 3. VOLTAGE SENSITIVITY
	// ============================================================
	{
		SpectrumID:         3,
		SpectrumName:       "voltage_sensitivity",
		Setting:            "insulated",
		SettingEnd:         "negative",
		BarrierType:        "mistakes",
		BarrierDescription: "Misses social cues, seems indifferent — doesn't realise impact",
		AvoidLanguage: []string{
			"You upset them and you don't even care",
		},
		RecommendedLanguage: []string{
			"You didn't feel that, but they felt it. Let me show you what I noticed.",
		},
		PlayAvenues: []string{
			"Use video playback of social situations",
			"\"Watch — what do you see on their face here?\"",
			"Build pattern recognition explicitly",
		},
	},
	{
		SpectrumID:         3,
		SpectrumName:       "voltage_sensitivity",
		Setting:            "conductive",
		SettingEnd:         "positive",
		BarrierType:        "mistakes",
		BarrierDescription: "Absorbs classroom tension, overwhelmed by others' mistakes — takes on others' emotions",
		AvoidLanguage: []string{
			"Stop being so sensitive",
		},
		RecommendedLanguage: []string{
			"You're feeling the room. That's a gift, but right now you need your own energy back. Let's find a calm corner.",
		},
		PlayAvenues: []string{
			"Create \"voltage shelters\" — quiet spaces, noise-cancelling headphones, predictable routines",
			"Teach them to recognise when they're absorbing",
		},
	},

	// ============================================================
	// 4. THREAT RESPONSE
	// ============================================================
	{
		SpectrumID:         4,
		SpectrumName:       "threat_response",
		Setting:            "passive",
		SettingEnd:         "negative",
		BarrierType:        "starting",
		BarrierDescription: "Freezes when stuck or wrong",
		AvoidLanguage: []string{
			"Just try, it doesn't matter if you're wrong",
		},
		RecommendedLanguage: []string{
			"Your brain just put up a stop sign. Let's look at the sign together. What's it saying?",
		},
		PlayAvenues: []string{
			"Low-stakes exploration with guaranteed safety",
			"\"Let's make the best wrong answer we can think of. What's the most interesting way to be wrong?\"",
		},
	},
	{
		SpectrumID:         4,
		SpectrumName:       "threat_response",
		Setting:            "aggressive",
		SettingEnd:         "positive",
		BarrierType:        "mistakes",
		BarrierDescription: "Lashes out when stuck or wrong",
		AvoidLanguage: []string{
			"Don't you dare speak to me like that",
		},
		RecommendedLanguage: []string{
			"Something just felt like an attack. I'm not attacking you — show me where it landed.",
		},
		PlayAvenues: []string{
			"Rough-and-tumble play with clear rules",
			"Practice stopping when someone says \"pause\"",
			"Build the experience of safe conflict",
		},
	},

	// ============================================================
	// 5. CARE RESPONSE
	// ============================================================
	{
		SpectrumID:         5,
		SpectrumName:       "care_response",
		Setting:            "detached",
		SettingEnd:         "negative",
		BarrierType:        "mistakes",
		BarrierDescription: "Seems not to care about impact",
		AvoidLanguage: []string{
			"You should feel bad about that",
		},
		RecommendedLanguage: []string{
			"I'm not telling you how to feel. I'm telling you how they felt. Two different things.",
		},
		PlayAvenues: []string{
			"Use third-person perspective: stories, puppets, role-play",
			"\"How do you think that character felt?\" No personal charge",
		},
	},
	{
		SpectrumID:         5,
		SpectrumName:       "care_response",
		Setting:            "nurturing",
		SettingEnd:         "positive",
		BarrierType:        "mistakes",
		BarrierDescription: "Takes responsibility for others' mistakes",
		AvoidLanguage: []string{
			"That's not your fault",
		},
		RecommendedLanguage: []string{
			"You're trying to carry something that isn't yours. Let's put it down together.",
		},
		PlayAvenues: []string{
			"Caretaking roles with clear boundaries",
			"Looking after a class plant vs. managing a peer's emotions",
			"Explicitly label the difference",
		},
	},

	// ============================================================
	// 6. RISK TOLERANCE
	// ============================================================
	{
		SpectrumID:         6,
		SpectrumName:       "risk_tolerance",
		Setting:            "averse",
		SettingEnd:         "negative",
		BarrierType:        "starting",
		BarrierDescription: "Won't start unless certain of success",
		AvoidLanguage: []string{
			"Just have a go",
		},
		RecommendedLanguage: []string{
			"Let's find the smallest possible first step. What's the least risky thing you could try?",
		},
		PlayAvenues: []string{
			"Scaffolded risk ladders",
			"\"Today we try step 1. Tomorrow step 2. We don't jump to step 5.\"",
		},
	},
	{
		SpectrumID:         6,
		SpectrumName:       "risk_tolerance",
		Setting:            "seeking",
		SettingEnd:         "positive",
		BarrierType:        "starting",
		BarrierDescription: "Starts but doesn't evaluate consequences — wrong kind of start",
		AvoidLanguage: []string{
			"Think before you act",
		},
		RecommendedLanguage: []string{
			"You're great at starting. Now let's practise the 10-second pause before you do. What might happen?",
		},
		PlayAvenues: []string{
			"Build consequence prediction into play",
			"\"Before you press that, what do you think will happen? Let's check.\"",
		},
	},

	// ============================================================
	// 7. INTEGRITY LOGIC
	// ============================================================
	{
		SpectrumID:         7,
		SpectrumName:       "integrity_logic",
		Setting:            "relativistic",
		SettingEnd:         "negative",
		BarrierType:        "mistakes",
		BarrierDescription: "\"It's not fair\" becomes excuse",
		AvoidLanguage: []string{
			"Life isn't fair",
		},
		RecommendedLanguage: []string{
			"You're right, this situation is unfair. And you still have a choice about what you do next.",
		},
		PlayAvenues: []string{
			"Explore different rule systems in games",
			"\"This game has different rules — let's try both and see how it feels.\"",
		},
	},
	{
		SpectrumID:         7,
		SpectrumName:       "integrity_logic",
		Setting:            "absolutist",
		SettingEnd:         "positive",
		BarrierType:        "mistakes",
		BarrierDescription: "Cannot recover from rule-breaking",
		AvoidLanguage: []string{
			"It's just a game",
		},
		RecommendedLanguage: []string{
			"The rule broke. That feels terrible. Let's fix the rule together so we can keep playing.",
		},
		PlayAvenues: []string{
			"Repair rituals",
			"When something breaks, the focus is on rebuilding the structure, not the mistake",
		},
	},

	// ============================================================
	// 8. MIRROR NEURON TUNING
	// ============================================================
	{
		SpectrumID:         8,
		SpectrumName:       "mirror_neuron_tuning",
		Setting:            "selective",
		SettingEnd:         "negative",
		BarrierType:        "verbalising",
		BarrierDescription: "Misses group emotional dynamics — doesn't read the room",
		AvoidLanguage: []string{
			"Can't you see everyone's upset?",
		},
		RecommendedLanguage: []string{
			"I'm going to point out what I'm noticing. Watch their faces now — what do you see?",
		},
		PlayAvenues: []string{
			"Emotion-identification games",
			"\"What emotion is this face? What might have caused it?\"",
		},
	},
	{
		SpectrumID:         8,
		SpectrumName:       "mirror_neuron_tuning",
		Setting:            "absorbent",
		SettingEnd:         "positive",
		BarrierType:        "verbalising",
		BarrierDescription: "Overwhelmed by group emotion — can't separate own feelings",
		AvoidLanguage: []string{
			"You're overreacting",
		},
		RecommendedLanguage: []string{
			"You're feeling the room. Which part is yours and which part is theirs? Let's sort them.",
		},
		PlayAvenues: []string{
			"Practice labelling: \"This feeling is mine. This feeling is from the group. This feeling is from the story.\"",
		},
	},

	// ============================================================
	// 9. ORDERLINESS
	// ============================================================
	{
		SpectrumID:         9,
		SpectrumName:       "orderliness",
		Setting:            "flexible",
		SettingEnd:         "negative",
		BarrierType:        "starting",
		BarrierDescription: "Difficulty with routines that require consistency — needs novelty to engage",
		AvoidLanguage: []string{
			"You never stick with anything",
		},
		RecommendedLanguage: []string{
			"You love new things. Let's find the novelty in this routine — what could we discover today?",
		},
		PlayAvenues: []string{
			"Build \"planned spontaneity\" — predictable times for unpredictable activities",
		},
	},
	{
		SpectrumID:         9,
		SpectrumName:       "orderliness",
		Setting:            "ordered",
		SettingEnd:         "positive",
		BarrierType:        "starting",
		BarrierDescription: "Cannot cope with unexpected changes — disruption derails them",
		AvoidLanguage: []string{
			"Just go with it",
		},
		RecommendedLanguage: []string{
			"The plan changed. That's hard. Let's make a new plan together right now.",
		},
		PlayAvenues: []string{
			"Practice \"change rituals\"",
			"When something changes, we stop, name it, and rebuild the structure",
		},
	},

	// ============================================================
	// 10. RESPONSIBILITY THRESHOLD
	// ============================================================
	{
		SpectrumID:         10,
		SpectrumName:       "responsibility_threshold",
		Setting:            "deflecting",
		SettingEnd:         "negative",
		BarrierType:        "persisting",
		BarrierDescription: "Avoids ownership of outcomes — 'it wasn't me' becomes default response",
		AvoidLanguage: []string{
			"It IS your fault",
			"Stop blaming everyone else",
		},
		RecommendedLanguage: []string{
			"What happened? What was the bit you did? That's your bit — let's fix just that part.",
			"Everyone has a part. Your part was small, but it's yours. Owning it makes it smaller.",
		},
		PlayAvenues: []string{
			"Sorting games: 'My part / Your part / Nobody's part'",
			"Post-mortem without blame — detectives find facts, not villains",
		},
	},
	{
		SpectrumID:         10,
		SpectrumName:       "responsibility_threshold",
		Setting:            "absorbing",
		SettingEnd:         "positive",
		BarrierType:        "verbalising",
		BarrierDescription: "Takes on blame for everything — 'it's all my fault' even when it isn't",
		AvoidLanguage: []string{
			"Don't worry about it",
			"It's fine, it wasn't you",
		},
		RecommendedLanguage: []string{
			"You're carrying more than your share. Let's divide the responsibility fairly — which part is actually yours?",
			"Your care is a strength. But if you carry everyone's load, you'll drop yours.",
		},
		PlayAvenues: []string{
			"Weight-sorting: physical objects representing 'my bit' vs 'not my bit'",
			"Responsibility pie chart — draw slices for each person's part",
		},
	},

	// ============================================================
	// 11. LOSS SENSITIVITY
	// ============================================================
	{
		SpectrumID:         11,
		SpectrumName:       "loss_sensitivity",
		Setting:            "detached",
		SettingEnd:         "negative",
		BarrierType:        "starting",
		BarrierDescription: "Doesn't engage with shared resources — lets go too easily, appears not to care",
		AvoidLanguage: []string{
			"Don't you care about anything?",
			"Look after your things",
		},
		RecommendedLanguage: []string{
			"This one matters to you even if you don't show it. Let's find one thing worth keeping safe today.",
			"You're good at letting go. Sometimes holding on is the brave thing.",
		},
		PlayAvenues: []string{
			"Favourite-thing journal — photograph and name one meaningful object per week",
			"Lending library — practice lending AND asking for return",
		},
	},
	{
		SpectrumID:         11,
		SpectrumName:       "loss_sensitivity",
		Setting:            "territorial",
		SettingEnd:         "positive",
		BarrierType:        "sharing",
		BarrierDescription: "Cannot tolerate removal of possessions or changes to 'their' space — tidy-up feels like theft",
		AvoidLanguage: []string{
			"It's just a toy",
			"Stop being selfish",
			"Tidy up now",
		},
		RecommendedLanguage: []string{
			"Your things are going home, not away. They'll be here tomorrow.",
			"You love this thing. That's not wrong. Let's find it a safe place for overnight.",
		},
		PlayAvenues: []string{
			"'Things go home' ritual — named places where objects rest, not disappear",
			"Security box — one item that never has to be shared, everything else rotates",
		},
	},

	// ============================================================
	// 12. LIBIDO
	// ============================================================
	{
		SpectrumID:         12,
		SpectrumName:       "libido",
		Setting:            "restrained",
		SettingEnd:         "negative",
		BarrierType:        "starting",
		BarrierDescription: "Low drive expression may mask needs — difficulty articulating desires or engaging in physical play",
		AvoidLanguage: []string{
			"Join in like everyone else",
			"Don't be shy",
		},
		RecommendedLanguage: []string{
			"You get to choose how much energy you put out. What feels right for you today?",
			"Not everyone shows excitement the same way. Your quiet way counts.",
		},
		PlayAvenues: []string{
			"Choice-board activities — 'I want' practice in low-pressure settings",
			"Graduated physical play — from solo sensory to paired activities at their pace",
		},
	},
	{
		SpectrumID:         12,
		SpectrumName:       "libido",
		Setting:            "expressive",
		SettingEnd:         "positive",
		BarrierType:        "regulating",
		BarrierDescription: "High drive expression may overwhelm peers — physical affection or desire expressed without calibration",
		AvoidLanguage: []string{
			"Stop touching people",
			"That's inappropriate",
			"Control yourself",
		},
		RecommendedLanguage: []string{
			"Your energy is big and that's OK. Let's find where it fits best.",
			"Before you touch, check: did they say yes? Asking first makes it better for everyone.",
		},
		PlayAvenues: []string{
			"Consent-checkpoint games — ask-before-touching becomes automatic",
			"High-energy channelling — structured physical outlets matched to drive level",
		},
	},
}

// GetResponseForSpectrum returns both matrix entries for a given spectrum ID
func GetResponseForSpectrum(spectrumID int) []ResponseMatrixEntry {
	var entries []ResponseMatrixEntry
	for _, e := range ResponseMatrix {
		if e.SpectrumID == spectrumID {
			entries = append(entries, e)
		}
	}
	return entries
}

// GetResponseForSetting returns the matrix entry for a specific spectrum + setting combination
func GetResponseForSetting(spectrumName, setting string) *ResponseMatrixEntry {
	for i := range ResponseMatrix {
		if ResponseMatrix[i].SpectrumName == spectrumName && ResponseMatrix[i].Setting == setting {
			return &ResponseMatrix[i]
		}
	}
	return nil
}

// GetResponsesForBarrier returns all matrix entries for a specific barrier type
func GetResponsesForBarrier(barrierType string) []ResponseMatrixEntry {
	var entries []ResponseMatrixEntry
	for _, e := range ResponseMatrix {
		if e.BarrierType == barrierType {
			entries = append(entries, e)
		}
	}
	return entries
}

// InterpretProfileValue converts a -2 to +2 slider value to the dominant setting name
func InterpretProfileValue(spectrumName string, value float64) string {
	for _, e := range ResponseMatrix {
		if e.SpectrumName == spectrumName {
			if value < 0 && e.SettingEnd == "negative" {
				return e.Setting
			}
			if value > 0 && e.SettingEnd == "positive" {
				return e.Setting
			}
			// Neutral — return whichever end is first (negative)
			if value == 0 && e.SettingEnd == "negative" {
				return e.Setting
			}
		}
	}
	return "neutral"
}

// GeneratePersonalisedPlan takes a map of spectrum_name -> value and returns
// actionable strategies tailored to the student's ETP profile.
type PersonalisedAction struct {
	SpectrumName    string   `json:"spectrum_name"`
	Setting         string   `json:"setting"`
	BarrierType     string   `json:"barrier_type"`
	BarrierDesc     string   `json:"barrier_description"`
	AvoidSaying     []string `json:"avoid_saying"`
	InsteadSay      []string `json:"instead_say"`
	PlayStrategies  []string `json:"play_strategies"`
	SettingStrength string   `json:"setting_strength"` // "mild", "moderate", "strong"
}

// GeneratePersonalisedPlan analyses an ETP profile and returns targeted strategies.
// profileValues maps spectrum_name -> float64 value (range -2 to +2).
// Only spectra with |value| > threshold are included (mild settings get lighter guidance).
func GeneratePersonalisedPlan(profileValues map[string]float64, threshold float64) []PersonalisedAction {
	if threshold == 0 {
		threshold = 0.5 // Default: ignore near-neutral settings
	}

	var actions []PersonalisedAction

	for spectrumName, value := range profileValues {
		absVal := value
		if absVal < 0 {
			absVal = -absVal
		}

		if absVal < threshold {
			continue // Near-neutral, no specific intervention needed
		}

		setting := InterpretProfileValue(spectrumName, value)
		entry := GetResponseForSetting(spectrumName, setting)
		if entry == nil {
			continue
		}

		strength := "mild"
		if absVal >= 1.5 {
			strength = "strong"
		} else if absVal >= 1.0 {
			strength = "moderate"
		}

		actions = append(actions, PersonalisedAction{
			SpectrumName:    spectrumName,
			Setting:         setting,
			BarrierType:     entry.BarrierType,
			BarrierDesc:     entry.BarrierDescription,
			AvoidSaying:     entry.AvoidLanguage,
			InsteadSay:      entry.RecommendedLanguage,
			PlayStrategies:  entry.PlayAvenues,
			SettingStrength: strength,
		})
	}

	return actions
}

// GetAdultLanguageForProfile returns the language shifts relevant to a student's
// dominant settings — used by the coach handler to personalise responses.
func GetAdultLanguageForProfile(profileValues map[string]float64) map[string][]string {
	result := make(map[string][]string)

	for spectrumName, value := range profileValues {
		absVal := value
		if absVal < 0 {
			absVal = -absVal
		}
		if absVal < 0.5 {
			continue
		}

		setting := InterpretProfileValue(spectrumName, value)
		entry := GetResponseForSetting(spectrumName, setting)
		if entry == nil {
			continue
		}

		result[spectrumName] = entry.RecommendedLanguage
	}

	return result
}
