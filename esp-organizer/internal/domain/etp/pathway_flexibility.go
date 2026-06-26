package etp

import "time"

// =============================================================================
// Pathway Flexibility — "The Differentiable Curriculum"
// =============================================================================
//
// The book doesn't change. The path through it does.
//
// Fixed elements:        32 activities, 9 themes, 3 cycles, sticker sheet
// Differentiated elements: order, language, celebration, difficulty, variation
//
// The book is the hardware. The parent is the operating system. The app is
// the compiler.
//
// This file models the choice architecture that makes one physical book
// produce infinite pathways:
//
//   1. NON-LINEAR NAVIGATION — parents choose themes by readiness, not order
//   2. DIFFICULTY CYCLES      — 3 layers per theme (not sequential pages)
//   3. STICKER CHOICES        — multiple celebrations per theme
//   4. "TRY INSTEAD" VARIANTS — printed adaptations on each spread
//   5. SUGGESTION ENGINE       — algorithm recommends next theme + cycle
//   6. CHOICE LOGGING          — every decision feeds back into the model
//
// The parent never feels prescribed to. The app never demands. The child
// leads, the parent follows, the system learns.
// =============================================================================

// =============================================================================
// Theme Weeks & Difficulty Cycles
// =============================================================================

// ThemeWeek represents one of the 9 spectrum themes at a specific difficulty
// cycle. 9 themes × 3 cycles = 27 themed weeks + 5 seasonal specials = 32.
type ThemeWeek struct {
	ID         string `json:"theme_week_id"` // e.g., "yes_and_no_cycle2"
	SpectrumID int    `json:"spectrum_id"`   // 1-9
	ThemeName  string `json:"theme_name"`    // e.g., "Yes and No"
	Cycle      int    `json:"cycle"`         // 1, 2, or 3

	// Cycle descriptions — what changes at each level
	CycleLabel       string `json:"cycle_label"`       // "Silly Questions", "Real Choices", "Boundary Negotiation"
	CycleDescription string `json:"cycle_description"` // What this difficulty level looks like

	// Navigation prompts — printed on the spread
	TryThisWhen []string `json:"try_this_when"` // readiness signals
	SkipIf      []string `json:"skip_if"`       // noise/timing signals

	// Activities at this cycle level
	Activities []ThemeActivity `json:"activities"`

	// Sticker options — multiple celebrations per theme
	StickerChoices []StickerChoice `json:"sticker_choices"`
}

// ThemeActivity is a single activity within a theme week, with "try instead"
// variations printed directly on the page.
type ThemeActivity struct {
	ActivityID  string `json:"activity_id"` // → Activity.ID in activity_experience.go
	Name        string `json:"name"`
	Mode        string `json:"mode"` // "with_you" or "while_you"
	Duration    string `json:"duration"`
	Instruction string `json:"instruction"` // Base instructions

	// "Try Instead" variations — printed on the spread so no app needed
	Variations []ActivityVariation `json:"variations"`
}

// ActivityVariation is a printed adaptation on the page. The parent reads
// all three and picks the one that fits today.
type ActivityVariation struct {
	Condition      string `json:"condition"`       // "If they loved this", "If they struggled", "If they were overwhelmed"
	Suggestion     string `json:"suggestion"`      // "Next time, try it with stuffed animals"
	AdjustmentType string `json:"adjustment_type"` // "escalate", "simplify", "scaffold"
}

// StickerChoice represents one celebration option on the sticker sheet.
// Multiple options per theme — the parent picks which fits.
type StickerChoice struct {
	Icon       string `json:"icon"`        // emoji for the sticker
	Label      string `json:"label"`       // "I said a clear NO"
	SkillFocus string `json:"skill_focus"` // which skill this celebration reinforces
}

// =============================================================================
// Choice Logging — What the App Captures
// =============================================================================

// PathwayChoice records every decision the parent makes. This is the
// raw signal that feeds into the suggestion engine.
type PathwayChoice struct {
	ID        string    `json:"choice_id"`
	ChildID   string    `json:"child_id"`
	Timestamp time.Time `json:"timestamp"`

	// What they chose
	ThemeWeekID       string `json:"theme_week_id"`                // which theme + cycle
	ActivityID        string `json:"activity_id"`                  // which activity
	ChosenVariation   string `json:"chosen_variation,omitempty"`   // which "try instead" they picked
	ChosenCelebration string `json:"chosen_celebration,omitempty"` // which sticker they used

	// How it went
	DifficultyRating int    `json:"difficulty_rating"` // 1=easy, 2=about right, 3=hard
	ParentNote       string `json:"parent_note,omitempty"`

	// System-linked observation (if they also logged a Three-Check Sensor reading)
	ObservationID string `json:"observation_id,omitempty"` // → ObservationEvent.ID
}

// =============================================================================
// Suggestion Engine — The Compiler
// =============================================================================

// ThemeSuggestion is a single recommendation from the suggestion engine.
// The app shows 3 of these; the parent picks one (or ignores all).
type ThemeSuggestion struct {
	ThemeWeekID string `json:"theme_week_id"`
	ThemeName   string `json:"theme_name"`
	Cycle       int    `json:"cycle"`
	Reason      string `json:"reason"` // human-readable explanation

	// Priority score — internal, not shown to parent
	// Higher = more recommended
	Priority float64 `json:"priority"`

	// Suggestion framing
	Icon    string `json:"icon"`    // emoji for the card
	Tagline string `json:"tagline"` // one-line summary shown on the card
}

// SuggestNextThemes returns up to 3 theme recommendations for a child.
// It considers:
//   - Which themes are least practised (coverage)
//   - Which themes have low cycle progression (growth opportunity)
//   - Which themes align with current ventral driver needs (driver support)
//   - Which themes the child has been avoiding (gentle nudge)
func SuggestNextThemes(
	childID string,
	profile map[int]float64, // spectrum_id → current value (-2 to +2)
	history []PathwayChoice,
	drivers []DriverMatch, // from DetectLikelyDrivers()
) []ThemeSuggestion {

	suggestions := make([]ThemeSuggestion, 0, 3)

	// === STEP 1: Count theme coverage ===
	themeCounts := make(map[int]int)       // spectrum_id → times attempted
	themeCycles := make(map[int]int)       // spectrum_id → highest cycle completed
	lastAttempt := make(map[int]time.Time) // spectrum_id → most recent attempt

	for _, choice := range history {
		tw := GetThemeWeekByID(choice.ThemeWeekID)
		if tw == nil {
			continue
		}
		themeCounts[tw.SpectrumID]++
		if tw.Cycle > themeCycles[tw.SpectrumID] {
			themeCycles[tw.SpectrumID] = tw.Cycle
		}
		if choice.Timestamp.After(lastAttempt[tw.SpectrumID]) {
			lastAttempt[tw.SpectrumID] = choice.Timestamp
		}
	}

	// === STEP 2: Score each theme ===
	type scoredTheme struct {
		spectrumID int
		cycle      int
		score      float64
		reason     string
	}

	scored := make([]scoredTheme, 0, 12)

	for specID := 1; specID <= 12; specID++ {
		spec := GetSpectrumByID(specID)
		if spec == nil {
			continue
		}

		score := 0.0
		reasons := make([]string, 0, 3)

		// Coverage factor — less practised themes get higher priority
		count := themeCounts[specID]
		if count == 0 {
			score += 3.0
			reasons = append(reasons, "haven't tried this yet")
		} else if count < 3 {
			score += 2.0 - float64(count)*0.5
			reasons = append(reasons, "only tried this "+pluralise(count, "time"))
		}

		// Recency factor — themes not attempted recently get a boost
		if last, ok := lastAttempt[specID]; ok {
			weeksSince := time.Since(last).Hours() / (24 * 7)
			if weeksSince > 4 {
				score += 1.0
				reasons = append(reasons, "haven't done this in "+pluralise(int(weeksSince), "week"))
			}
		}

		// Driver alignment — themes that support detected drivers get a boost
		for _, dm := range drivers {
			driver := GetDriverByID(dm.DriverID)
			if driver == nil {
				continue
			}
			for _, ds := range driver.SpectrumSignature {
				if ds.SpectrumID == specID {
					score += dm.Confidence * 0.5
					reasons = append(reasons, "aligns with their "+driver.PersonaName+" energy")
					break
				}
			}
		}

		// Determine cycle
		currentCycle := themeCycles[specID]
		nextCycle := currentCycle + 1
		if nextCycle > 3 {
			nextCycle = 3 // cap at cycle 3, but can repeat
		}
		if nextCycle == 0 {
			nextCycle = 1 // first attempt
		}

		// No load adjustment — transient states are not part of this system.
		// Energy directionality already captures stress-reaction direction.
		// "Bad day? Skip it." is handled at the observation level, not the suggestion level.

		reason := reasons[0]
		if len(reasons) > 1 {
			reason = reasons[0]
			for i := 1; i < len(reasons); i++ {
				reason += "; " + reasons[i]
			}
		}

		scored = append(scored, scoredTheme{
			spectrumID: specID,
			cycle:      nextCycle,
			score:      score,
			reason:     reason,
		})
	}

	// === STEP 3: Sort by score (descending) and pick top 3 ===
	// Simple selection sort — only 9 items
	for i := 0; i < len(scored); i++ {
		maxIdx := i
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[maxIdx].score {
				maxIdx = j
			}
		}
		scored[i], scored[maxIdx] = scored[maxIdx], scored[i]
	}

	for i := 0; i < 3 && i < len(scored); i++ {
		st := scored[i]
		spec := GetSpectrumByID(st.spectrumID)
		if spec == nil {
			continue
		}

		tw := findThemeWeek(st.spectrumID, st.cycle)
		twID := ""
		tagline := ""
		icon := ""
		if tw != nil {
			twID = tw.ID
			tagline = tw.CycleLabel
			icon = themeIcon(st.spectrumID)
		}

		suggestions = append(suggestions, ThemeSuggestion{
			ThemeWeekID: twID,
			ThemeName:   spec.Name,
			Cycle:       st.cycle,
			Reason:      capitalise(st.reason),
			Priority:    st.score,
			Icon:        icon,
			Tagline:     tagline,
		})
	}

	return suggestions
}

// SuggestCycle determines the appropriate difficulty cycle for a theme
// based on the child's history. Looks at the last 3 attempts:
//   - All success (difficulty ≤ 2) → suggest next cycle
//   - Mixed → suggest same cycle
//   - Repeated struggle (difficulty = 3) → suggest previous cycle
func SuggestCycle(spectrumID int, history []PathwayChoice) int {
	// Collect recent attempts at this theme
	var recent []PathwayChoice
	for i := len(history) - 1; i >= 0 && len(recent) < 3; i-- {
		tw := GetThemeWeekByID(history[i].ThemeWeekID)
		if tw != nil && tw.SpectrumID == spectrumID {
			recent = append(recent, history[i])
		}
	}

	if len(recent) == 0 {
		return 1 // first attempt → cycle 1
	}

	// Find current cycle and assess difficulty
	currentCycle := 1
	struggles := 0
	successes := 0

	for _, r := range recent {
		tw := GetThemeWeekByID(r.ThemeWeekID)
		if tw != nil && tw.Cycle > currentCycle {
			currentCycle = tw.Cycle
		}
		if r.DifficultyRating >= 3 {
			struggles++
		} else if r.DifficultyRating <= 1 {
			successes++
		}
	}

	// Decision logic
	if successes == len(recent) && currentCycle < 3 {
		return currentCycle + 1 // all easy → escalate
	}
	if struggles >= 2 && currentCycle > 1 {
		return currentCycle - 1 // repeated struggle → step back
	}
	return currentCycle // mixed → stay
}

// =============================================================================
// Theme Week Data — 9 themes × 3 cycles = 27 themed weeks
// =============================================================================

// ThemeWeeks contains all 27 themed weeks (9 spectra × 3 cycles).
// Seasonal specials (5 extra) will be added separately.
var ThemeWeeks = []ThemeWeek{
	// =========================================================================
	// SPECTRUM 1: SOCIAL GRAVITY — "Near and Far"
	// =========================================================================
	{
		ID: "near_and_far_cycle1", SpectrumID: 1, ThemeName: "Near and Far", Cycle: 1,
		CycleLabel: "Side by Side", CycleDescription: "Parallel play — same room, separate tasks. No pressure to interact.",
		TryThisWhen: []string{
			"Your child plays alone even when others are around",
			"They cling to you when other children appear",
			"You want to understand their social comfort zone",
		},
		SkipIf: []string{"They're overtired", "You're rushed", "They just had a social meltdown"},
		Activities: []ThemeActivity{
			{
				ActivityID: "near_far_c1_parallel", Name: "The Beside Game", Mode: "with_you", Duration: "3 min",
				Instruction: "Sit next to each other. You draw a picture. They draw a picture. No talking needed.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Next time, swap pictures halfway through", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Next time, just sit together — no drawing needed", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Try it in different rooms with the door open", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🪐", Label: "I played near someone", SkillFocus: "social_proximity"},
			{Icon: "🌟", Label: "I shared my space", SkillFocus: "joint_attention"},
			{Icon: "🫧", Label: "I was happy on my own", SkillFocus: "independent_play"},
			{Icon: "🤝", Label: "I chose to play together", SkillFocus: "social_initiation"},
		},
	},
	{
		ID: "near_and_far_cycle2", SpectrumID: 1, ThemeName: "Near and Far", Cycle: 2,
		CycleLabel: "Invitation to Join", CycleDescription: "One person invites, the other decides. Practice saying yes AND no to togetherness.",
		TryThisWhen: []string{
			"They can play near others comfortably",
			"You want to practise social invitations",
			"They need help joining group play",
		},
		SkipIf: []string{"They're clingy today", "There's been recent separation anxiety", "They've had a rough drop-off"},
		Activities: []ThemeActivity{
			{
				ActivityID: "near_far_c2_invite", Name: "The Door Game", Mode: "with_you", Duration: "5 min",
				Instruction: "Play in separate rooms. Take turns knocking on each other's 'door' and asking: 'May I come in?' Either answer is fine.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Add stuffed animals who also knock", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "You always knock first — let them just answer", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Use a tent or blanket fort instead of rooms", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🚪", Label: "I knocked and asked", SkillFocus: "social_initiation"},
			{Icon: "✅", Label: "I said 'come in!'", SkillFocus: "social_acceptance"},
			{Icon: "🙅", Label: "I said 'not yet' — and that was OK", SkillFocus: "boundary_setting"},
			{Icon: "⏳", Label: "I waited for the answer", SkillFocus: "inhibitory_control"},
		},
	},
	{
		ID: "near_and_far_cycle3", SpectrumID: 1, ThemeName: "Near and Far", Cycle: 3,
		CycleLabel: "The Push and Pull", CycleDescription: "Negotiating how close or far — in real time, with real feelings.",
		TryThisWhen: []string{
			"They can invite and accept invitations",
			"You want to practise mid-play renegotiation",
			"They're ready for bigger social challenges",
		},
		SkipIf: []string{"They're fragile today", "There's been a friendship conflict", "You don't have 10+ minutes"},
		Activities: []ThemeActivity{
			{
				ActivityID: "near_far_c3_negotiate", Name: "The Stretchy String", Mode: "with_you", Duration: "5 min",
				Instruction: "Hold a scarf between you. Walk apart until it's taut. Walk close until it's floppy. Talk about what feels right. 'How far is too far? How close is too close?'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Try it with a friend instead of parent", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Use a shorter scarf — less range to negotiate", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Just do the 'close' end — skip the walking apart", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🤸", Label: "I found my comfortable distance", SkillFocus: "social_calibration"},
			{Icon: "💬", Label: "I said what I needed", SkillFocus: "assertive_communication"},
			{Icon: "👂", Label: "I heard what they needed", SkillFocus: "perspective_taking"},
			{Icon: "🔄", Label: "We found a middle", SkillFocus: "negotiation"},
		},
	},

	// =========================================================================
	// SPECTRUM 2: ENERGY DIRECTIONALITY — "Loud and Quiet"
	// =========================================================================
	{
		ID: "loud_and_quiet_cycle1", SpectrumID: 2, ThemeName: "Loud and Quiet", Cycle: 1,
		CycleLabel: "Volume Dial", CycleDescription: "Exploring loud and quiet as choices, not commands.",
		TryThisWhen: []string{
			"Your child is always at one volume",
			"They can't seem to find 'quiet' (or 'loud')",
			"You want to explore their energy default",
		},
		SkipIf: []string{"Nap time is close", "Siblings are sleeping", "They're already overstimulated"},
		Activities: []ThemeActivity{
			{
				ActivityID: "loud_quiet_c1_dial", Name: "The Volume Knob", Mode: "with_you", Duration: "3 min",
				Instruction: "Use your hand as a 'volume dial'. Turn it up → both get louder. Turn it down → both get quieter. Take turns being the DJ.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Add instruments — pots and wooden spoons", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Just do whisper-to-talking — skip the shouting end", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Let them control the dial the whole time", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🔊", Label: "I found my loud voice", SkillFocus: "energy_expression"},
			{Icon: "🤫", Label: "I found my quiet voice", SkillFocus: "energy_containment"},
			{Icon: "🎚️", Label: "I changed my volume", SkillFocus: "self_regulation"},
			{Icon: "🎵", Label: "I was the DJ", SkillFocus: "social_leadership"},
		},
	},
	{
		ID: "loud_and_quiet_cycle2", SpectrumID: 2, ThemeName: "Loud and Quiet", Cycle: 2,
		CycleLabel: "Matching Energy", CycleDescription: "Reading the room — noticing when loud or quiet fits.",
		TryThisWhen: []string{
			"They can control their volume when asked",
			"They need help reading social energy cues",
			"You want to practise context-switching",
		},
		SkipIf: []string{"They're very wound up", "You need quiet right now", "They're shutting down"},
		Activities: []ThemeActivity{
			{
				ActivityID: "loud_quiet_c2_match", Name: "Library vs Playground", Mode: "with_you", Duration: "5 min",
				Instruction: "Pick a room. Announce: 'This is the library!' (whisper). Then: 'Now it's the playground!' (loud). Switch back and forth. Notice the gear change.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Let them invent the places — 'spaceship' or 'underwater'", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Just do two switches, not rapid-fire", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Use teddy bears to demonstrate first", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "📚", Label: "I found my library voice", SkillFocus: "contextual_regulation"},
			{Icon: "🏃", Label: "I found my playground voice", SkillFocus: "contextual_expression"},
			{Icon: "🔄", Label: "I switched between both", SkillFocus: "cognitive_flexibility"},
			{Icon: "👀", Label: "I noticed what voice the room needed", SkillFocus: "social_awareness"},
		},
	},
	{
		ID: "loud_and_quiet_cycle3", SpectrumID: 2, ThemeName: "Loud and Quiet", Cycle: 3,
		CycleLabel: "The Whisper Challenge", CycleDescription: "Exciting activity, quiet voice — separating energy from volume.",
		TryThisWhen: []string{
			"They can match energy to context",
			"You want to build voltage containment",
			"They're ready for the Performer stress-test",
		},
		SkipIf: []string{"They need to be loud right now", "They're emotionally flat", "It's been a hard day"},
		Activities: []ThemeActivity{
			{
				ActivityID: "loud_quiet_c3_whisper", Name: "Whisper Tag", Mode: "with_you", Duration: "5 min",
				Instruction: "Play chase — but whispering. All the excitement, none of the volume. If anyone shouts, freeze for 5 seconds.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Add a 'silent' level — communicate with gestures only", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Switch to whisper hide-and-seek — less running", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Just do whisper counting — no chase element", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🤐", Label: "I whispered when I wanted to shout", SkillFocus: "impulse_containment"},
			{Icon: "⚡", Label: "I kept the excitement inside", SkillFocus: "voltage_regulation"},
			{Icon: "🧊", Label: "I froze when I needed to", SkillFocus: "inhibitory_control"},
			{Icon: "😂", Label: "It was so hard not to laugh!", SkillFocus: "emotional_awareness"},
		},
	},

	// =========================================================================
	// SPECTRUM 3: VOLTAGE SENSITIVITY — "Big Feelings, Small Feelings"
	// =========================================================================
	{
		ID: "big_small_feelings_cycle1", SpectrumID: 3, ThemeName: "Big Feelings, Small Feelings", Cycle: 1,
		CycleLabel: "The Feelings Thermometer", CycleDescription: "Noticing feelings have sizes — not just on/off.",
		TryThisWhen: []string{
			"Everything is either FINE or a DISASTER",
			"They can't tell the difference between annoyed and furious",
			"You want to introduce emotional vocabulary",
		},
		SkipIf: []string{"They're mid-meltdown", "They're already emotionally drained", "You need them to just comply right now"},
		Activities: []ThemeActivity{
			{
				ActivityID: "feelings_c1_thermo", Name: "How Big Is It?", Mode: "with_you", Duration: "3 min",
				Instruction: "Draw a thermometer on paper. Point to scenarios: 'Your Lego broke — how big is that feeling?' They point to the level. No right answer.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Use real scenarios from today", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Just do 'big' or 'small' — skip the middle", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Use teddy's feelings instead of theirs", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🌡️", Label: "I measured my feeling", SkillFocus: "emotional_granularity"},
			{Icon: "🫧", Label: "I had a small feeling and noticed it", SkillFocus: "interoception"},
			{Icon: "🌊", Label: "I had a big feeling and named it", SkillFocus: "emotional_labelling"},
			{Icon: "🧸", Label: "I helped teddy with their feelings", SkillFocus: "empathy_practice"},
		},
	},
	{
		ID: "big_small_feelings_cycle2", SpectrumID: 3, ThemeName: "Big Feelings, Small Feelings", Cycle: 2,
		CycleLabel: "The Wave Rider", CycleDescription: "Feelings come and go — you can ride them without drowning.",
		TryThisWhen: []string{
			"They can name big and small feelings",
			"They get stuck in feelings (can't let go)",
			"You want to build the 'this will pass' muscle",
		},
		SkipIf: []string{"They're in a difficult emotional period", "You don't have patience today", "They're ill"},
		Activities: []ThemeActivity{
			{
				ActivityID: "feelings_c2_wave", Name: "The Feeling Timer", Mode: "with_you", Duration: "3 min",
				Instruction: "When a big feeling arrives, say: 'Let's time it. How long does this wave last?' Use fingers to count. Afterwards: 'See? It came and it went.'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Let them time YOUR feelings too", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Just notice: 'Your face looks different now than 2 minutes ago'", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Skip the timing — just name the feeling and sit with it", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🏄", Label: "I rode a wave", SkillFocus: "distress_tolerance"},
			{Icon: "⏱️", Label: "I waited and it passed", SkillFocus: "emotional_patience"},
			{Icon: "🌈", Label: "I felt it change colour", SkillFocus: "emotional_awareness"},
			{Icon: "💪", Label: "I bounced back", SkillFocus: "resilience"},
		},
	},
	{
		ID: "big_small_feelings_cycle3", SpectrumID: 3, ThemeName: "Big Feelings, Small Feelings", Cycle: 3,
		CycleLabel: "The Feeling DJ", CycleDescription: "Choosing how to express the feeling — not suppressing, channelling.",
		TryThisWhen: []string{
			"They can ride a wave",
			"They need to learn feelings have outlets",
			"They're ready to channel rather than contain",
		},
		SkipIf: []string{"This is not about suppression", "They're already regulating well", "Today is not the day for big work"},
		Activities: []ThemeActivity{
			{
				ActivityID: "feelings_c3_channel", Name: "Same Feeling, Different Door", Mode: "with_you", Duration: "5 min",
				Instruction: "When a big feeling arrives, offer 3 exits: 'You can stomp it out, squeeze this cushion, or draw it. Same feeling — three doors. Which one today?'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Let them invent a 4th door", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Offer just 2 doors at first", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "You demonstrate all 3 doors first — they watch", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🚪", Label: "I chose my own door", SkillFocus: "emotional_agency"},
			{Icon: "🎨", Label: "I drew my feeling", SkillFocus: "emotional_expression"},
			{Icon: "👣", Label: "I stomped it out", SkillFocus: "physical_regulation"},
			{Icon: "🧠", Label: "I tried a new way", SkillFocus: "cognitive_flexibility"},
		},
	},

	// =========================================================================
	// SPECTRUM 4: THREAT RESPONSE — "Yes and No"
	// =========================================================================
	{
		ID: "yes_and_no_cycle1", SpectrumID: 4, ThemeName: "Yes and No", Cycle: 1,
		CycleLabel: "Silly Questions", CycleDescription: "Low-stakes yes/no practice — no real consequences.",
		TryThisWhen: []string{
			"Your child is learning to set boundaries",
			"You're noticing power struggles",
			"They need practice saying 'no' safely",
		},
		SkipIf: []string{"They're tired or hungry", "You're short on time", "They just had a big 'no' fight"},
		Activities: []ThemeActivity{
			{
				ActivityID: "yes_no_c1_silly", Name: "The Silly Question Game", Mode: "with_you", Duration: "3 min",
				Instruction: "Ask absurd yes/no questions: 'Should we put socks on our ears?' 'Should we eat dinner on the ceiling?' They practise saying YES and NO without stakes.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Let them ask YOU the silly questions", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Use hand puppets to ask the questions", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Just do 3 questions and stop. Less is fine.", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🎉", Label: "I said a big YES", SkillFocus: "assertive_agreement"},
			{Icon: "🛑", Label: "I said a clear NO", SkillFocus: "boundary_setting"},
			{Icon: "😂", Label: "I laughed at a silly question", SkillFocus: "playful_engagement"},
			{Icon: "🔄", Label: "I changed my mind", SkillFocus: "cognitive_flexibility"},
		},
	},
	{
		ID: "yes_and_no_cycle2", SpectrumID: 4, ThemeName: "Yes and No", Cycle: 2,
		CycleLabel: "Real Choices", CycleDescription: "Saying yes and no to things that actually matter — with supported consequences.",
		TryThisWhen: []string{
			"Silly questions are easy now",
			"They need practice with real-world boundaries",
			"You want to build genuine assertiveness",
		},
		SkipIf: []string{"They're fragile today", "You can't follow through on the consequences", "Bedtime is soon"},
		Activities: []ThemeActivity{
			{
				ActivityID: "yes_no_c2_real", Name: "The Two Plates", Mode: "with_you", Duration: "5 min",
				Instruction: "At snack time: two options on two plates. They choose ONE. The other goes away. Real choice, real consequence. 'You chose the banana. The apple goes back. That's a clear choice.'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Extend to non-food: 'Which park? Which book?'", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Make both options things they like equally", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Don't remove the unchosen item yet — build up to it", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🍌", Label: "I made a real choice", SkillFocus: "decision_making"},
			{Icon: "😌", Label: "I was OK with what I didn't choose", SkillFocus: "loss_tolerance"},
			{Icon: "🗣️", Label: "I said what I wanted", SkillFocus: "assertive_communication"},
			{Icon: "🤷", Label: "I changed my mind — and that was OK", SkillFocus: "flexibility"},
		},
	},
	{
		ID: "yes_and_no_cycle3", SpectrumID: 4, ThemeName: "Yes and No", Cycle: 3,
		CycleLabel: "Boundary Negotiation", CycleDescription: "When your 'no' meets someone else's 'no' — finding the middle.",
		TryThisWhen: []string{
			"They can make real choices confidently",
			"They need practice when their choice conflicts with yours",
			"You're ready for honest negotiation",
		},
		SkipIf: []string{"You're not in a position to negotiate today", "They're dysregulated", "Power dynamics are already tense"},
		Activities: []ThemeActivity{
			{
				ActivityID: "yes_no_c3_negotiate", Name: "The Trade Game", Mode: "with_you", Duration: "5 min",
				Instruction: "Each pick a toy. Try to trade. Either can say no. Practice: 'You said no. I'll ask differently.' or 'I said no. That's OK too.'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Try with a friend — real social negotiation", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "You always say yes first — model acceptance", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Use stuffed animals to negotiate first", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🤝", Label: "We found a deal", SkillFocus: "negotiation"},
			{Icon: "🛡️", Label: "I heard 'no' and was OK", SkillFocus: "rejection_tolerance"},
			{Icon: "💬", Label: "I asked a different way", SkillFocus: "persuasion"},
			{Icon: "✋", Label: "I respected their no", SkillFocus: "respect_for_boundaries"},
		},
	},

	// =========================================================================
	// SPECTRUM 5: CARE RESPONSE — "Mine and Yours"
	// =========================================================================
	{
		ID: "mine_and_yours_cycle1", SpectrumID: 5, ThemeName: "Mine and Yours", Cycle: 1,
		CycleLabel: "What's Mine", CycleDescription: "Understanding ownership — the foundation for sharing.",
		TryThisWhen: []string{
			"Everything is 'MINE!'",
			"They struggle with the concept of borrowing",
			"You want to build a foundation for sharing",
		},
		SkipIf: []string{"A sibling just took something", "They're already generous today", "Possessiveness is very high right now"},
		Activities: []ThemeActivity{
			{
				ActivityID: "mine_yours_c1_sort", Name: "The Sorting Game", Mode: "with_you", Duration: "3 min",
				Instruction: "Get 6 items. Three are yours, three are theirs. Sort them into piles. 'This is mine. This is yours.' Name ownership before asking them to share it.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Add a 'nobody's' pile for shared items", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Use items they don't care about (spoons, socks)", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Just do yours — don't sort theirs yet", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🏷️", Label: "I know what's mine", SkillFocus: "ownership_awareness"},
			{Icon: "👆", Label: "I know what's yours", SkillFocus: "other_awareness"},
			{Icon: "📦", Label: "I sorted them all", SkillFocus: "categorisation"},
			{Icon: "🫂", Label: "I found something that's ours", SkillFocus: "shared_ownership"},
		},
	},
	{
		ID: "mine_and_yours_cycle2", SpectrumID: 5, ThemeName: "Mine and Yours", Cycle: 2,
		CycleLabel: "Lending Library", CycleDescription: "Sharing with a return guarantee — building the trust bridge.",
		TryThisWhen: []string{
			"They know what's theirs but won't let go",
			"You want to introduce temporary sharing",
			"They need the 'it comes back' experience",
		},
		SkipIf: []string{"Something was recently broken or lost", "Trust is low today", "A new sibling just arrived"},
		Activities: []ThemeActivity{
			{
				ActivityID: "mine_yours_c2_lend", Name: "The Timer Share", Mode: "with_you", Duration: "5 min",
				Instruction: "Set a timer for 30 seconds. They lend you a toy. Timer goes off — it comes back. Every single time. Build the 'return' reflex.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Extend the timer to 1 minute", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Use a toy they don't care about first", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "YOU lend them your toy first — model the trust", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "⏱️", Label: "I shared for 30 seconds", SkillFocus: "temporary_sharing"},
			{Icon: "🔄", Label: "It came back!", SkillFocus: "trust_building"},
			{Icon: "😊", Label: "They smiled when I shared", SkillFocus: "social_reward"},
			{Icon: "🎁", Label: "I chose to share something special", SkillFocus: "generosity"},
		},
	},
	{
		ID: "mine_and_yours_cycle3", SpectrumID: 5, ThemeName: "Mine and Yours", Cycle: 3,
		CycleLabel: "The Gift", CycleDescription: "Giving without getting back — the real deal.",
		TryThisWhen: []string{
			"Timer-sharing is easy now",
			"They're ready for genuine giving",
			"You want to notice if sharing is genuine or performed",
		},
		SkipIf: []string{"They've been forced to share too much today", "Generosity should not be forced", "They're feeling deprived"},
		Activities: []ThemeActivity{
			{
				ActivityID: "mine_yours_c3_gift", Name: "The Surprise Box", Mode: "with_you", Duration: "5 min",
				Instruction: "Together, choose something to put in a 'surprise box' for someone else (sibling, friend, teddy). They pick, wrap, deliver. Notice: are they watching for the reaction?",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Let them choose the recipient too", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Make the gift something small — a drawing or sticker", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Make the recipient a teddy — lower stakes", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🎁", Label: "I gave something away", SkillFocus: "generosity"},
			{Icon: "😊", Label: "It felt good to give", SkillFocus: "intrinsic_reward"},
			{Icon: "🤫", Label: "I gave and didn't need praise", SkillFocus: "internal_locus"},
			{Icon: "🫶", Label: "I picked just the right thing", SkillFocus: "perspective_taking"},
		},
	},

	// =========================================================================
	// SPECTRUM 6: RISK TOLERANCE — "Try and Wait"
	// =========================================================================
	{
		ID: "try_and_wait_cycle1", SpectrumID: 6, ThemeName: "Try and Wait", Cycle: 1,
		CycleLabel: "Tiny Bravery", CycleDescription: "The smallest possible risk — with full safety net.",
		TryThisWhen: []string{"They won't try new things", "They watch but won't join", "You want to understand their approach to risk"},
		SkipIf:      []string{"They've been forced to try something today", "They're in cautious mode for good reason", "They're already being brave"},
		Activities: []ThemeActivity{
			{
				ActivityID: "try_wait_c1_tiny", Name: "One Taste", Mode: "with_you", Duration: "2 min",
				Instruction: "New food. One tiny taste. That's it. 'You don't have to like it. Just let your tongue have a look.' Same principle applies to any new experience.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Try 3 new tastes in a row — a tasting plate", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Just smell it — don't even taste", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "You taste it first and describe it", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "👅", Label: "I tried something new", SkillFocus: "novelty_tolerance"},
			{Icon: "🤔", Label: "I thought about it first", SkillFocus: "risk_assessment"},
			{Icon: "🙅", Label: "I said 'not today' — and that's OK", SkillFocus: "self_advocacy"},
			{Icon: "🦁", Label: "I was brave", SkillFocus: "courage"},
		},
	},
	{
		ID: "try_and_wait_cycle2", SpectrumID: 6, ThemeName: "Try and Wait", Cycle: 2,
		CycleLabel: "The Warm-Up Lap", CycleDescription: "Watching before doing — making observation legitimate.",
		TryThisWhen: []string{"One Taste is easy", "They need permission to watch before joining", "You want to validate cautious approaches"},
		SkipIf:      []string{"They're already jumping in", "Pressure to 'just try it' is high", "Social comparison is active"},
		Activities: []ThemeActivity{
			{
				ActivityID: "try_wait_c2_warmup", Name: "The Spy Game", Mode: "with_you", Duration: "5 min",
				Instruction: "At the playground: 'Let's spy on the slide. What do the other children do? How fast? Which way?' Observing IS participating. When ready, they try.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Spy on a new activity they've never tried", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Spy from inside the car if needed — window watching", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Spy at home on a video — no live pressure", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🕵️", Label: "I observed carefully", SkillFocus: "observational_learning"},
			{Icon: "📋", Label: "I made a plan first", SkillFocus: "planning"},
			{Icon: "🏃", Label: "I jumped in after watching", SkillFocus: "calculated_risk"},
			{Icon: "👍", Label: "I chose the right moment", SkillFocus: "timing"},
		},
	},
	{
		ID: "try_and_wait_cycle3", SpectrumID: 6, ThemeName: "Try and Wait", Cycle: 3,
		CycleLabel: "The Wobbly Bridge", CycleDescription: "Choosing to do something even though you know it might not work.",
		TryThisWhen: []string{"Watching is comfortable", "They need practice with acceptable failure", "You want to build frustration tolerance"},
		SkipIf:      []string{"They've had a failure recently", "Confidence is low", "They need a win today"},
		Activities: []ThemeActivity{
			{
				ActivityID: "try_wait_c3_wobbly", Name: "The Impossible Tower", Mode: "with_you", Duration: "5 min",
				Instruction: "Build a tower of blocks that WILL fall. Both know it's going to fall. The point is: build it anyway. When it falls: 'There it goes! Shall we try again or try something different?'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Time how long before it falls — beat the record", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Build a short tower that won't fall — then add one risky block", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "You build, they knock down — make the falling part the game", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🏗️", Label: "I built something that fell — and tried again", SkillFocus: "persistence"},
			{Icon: "💥", Label: "It fell and I laughed!", SkillFocus: "failure_tolerance"},
			{Icon: "🧱", Label: "I made it taller than last time", SkillFocus: "incremental_progress"},
			{Icon: "🤹", Label: "I tried a different way", SkillFocus: "creative_problem_solving"},
		},
	},

	// =========================================================================
	// SPECTRUM 7: INTEGRITY LOGIC — "Same and Different"
	// =========================================================================
	{
		ID: "same_and_different_cycle1", SpectrumID: 7, ThemeName: "Same and Different", Cycle: 1,
		CycleLabel: "The Rules Game", CycleDescription: "Rules exist — and they can be named, not just obeyed.",
		TryThisWhen: []string{"They can't cope with rule changes", "They need rules to feel safe", "You want to explore rigidity vs flexibility"},
		SkipIf:      []string{"They're already coping with too many changes", "A routine just changed", "Predictability is needed right now"},
		Activities: []ThemeActivity{
			{
				ActivityID: "same_diff_c1_rules", Name: "Today's Rules", Mode: "with_you", Duration: "3 min",
				Instruction: "Make silly rules together: 'Today, shoes go in the fridge. Cups live on the floor.' Follow them for 5 minutes. Then: 'OK, rules back to normal.' Notice: relief or disappointment?",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Let them make ALL the rules tomorrow", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Only change ONE rule — keep everything else the same", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Frame it as 'pretend' — 'let's PRETEND the rule is...'", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "📜", Label: "I followed a new rule", SkillFocus: "rule_following"},
			{Icon: "✏️", Label: "I made up a rule", SkillFocus: "rule_creation"},
			{Icon: "🔄", Label: "I was OK when it changed back", SkillFocus: "transition_tolerance"},
			{Icon: "😄", Label: "The silly rule made me laugh", SkillFocus: "playful_flexibility"},
		},
	},
	{
		ID: "same_and_different_cycle2", SpectrumID: 7, ThemeName: "Same and Different", Cycle: 2,
		CycleLabel: "The Exception", CycleDescription: "Sometimes rules change — and that can be OK.",
		TryThisWhen: []string{"They understand rules but can't flex them", "You want to introduce exceptions", "They need to see that changed rules aren't broken rules"},
		SkipIf:      []string{"They need structure right now", "Too many exceptions lately", "Authority trust is low"},
		Activities: []ThemeActivity{
			{
				ActivityID: "same_diff_c2_exception", Name: "The 'Just This Once' Card", Mode: "with_you", Duration: "3 min",
				Instruction: "Give them a physical card: 'Just This Once.' They can play it once today to change any rule. Pudding before dinner? Card played. Tomorrow it resets.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Give them TWO cards next time", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "You play the card first — show it's safe", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Just look at the card together — don't play it yet", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🃏", Label: "I played my card!", SkillFocus: "agency"},
			{Icon: "🧠", Label: "I chose not to play it", SkillFocus: "delayed_gratification"},
			{Icon: "🤔", Label: "I thought hard about when to use it", SkillFocus: "strategic_thinking"},
			{Icon: "😌", Label: "The rules came back and I was OK", SkillFocus: "routine_return"},
		},
	},
	{
		ID: "same_and_different_cycle3", SpectrumID: 7, ThemeName: "Same and Different", Cycle: 3,
		CycleLabel: "The Referee", CycleDescription: "Making rules for others — understanding fairness from the inside.",
		TryThisWhen: []string{"Exceptions are manageable", "They need to see rules from the maker's perspective", "You want to build fairness reasoning"},
		SkipIf:      []string{"They'd use rule-making to control others", "Power dynamics are tricky right now", "They're not yet comfortable with exceptions"},
		Activities: []ThemeActivity{
			{
				ActivityID: "same_diff_c3_referee", Name: "You're the Boss", Mode: "with_you", Duration: "5 min",
				Instruction: "Play a simple game (snakes and ladders, catch). They make up ONE new rule. Play with it. Then discuss: 'Was that rule fair? Did everyone have fun?'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Let them referee a game between two teddy bears", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Give them 2 rules to choose from instead of inventing", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Co-create the rule together — shared authorship", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "👑", Label: "I was the rule-maker", SkillFocus: "leadership"},
			{Icon: "⚖️", Label: "My rule was fair", SkillFocus: "fairness_reasoning"},
			{Icon: "🔧", Label: "I changed my rule to make it better", SkillFocus: "iterative_thinking"},
			{Icon: "🤝", Label: "Everyone agreed", SkillFocus: "consensus_building"},
		},
	},

	// =========================================================================
	// SPECTRUM 8: MIRROR NEURON TUNING — "Your Feelings, My Feelings"
	// =========================================================================
	{
		ID: "your_my_feelings_cycle1", SpectrumID: 8, ThemeName: "Your Feelings, My Feelings", Cycle: 1,
		CycleLabel: "Whose Feeling Is This?", CycleDescription: "Noticing that others have feelings — separate from yours.",
		TryThisWhen: []string{"They don't notice when others are upset", "They absorb everyone's feelings as their own", "You want to build emotional boundaries"},
		SkipIf:      []string{"Someone close is very upset right now", "They're emotionally saturated", "They need their own feelings attended to first"},
		Activities: []ThemeActivity{
			{
				ActivityID: "feelings_c1_whose", Name: "The Face Game", Mode: "with_you", Duration: "3 min",
				Instruction: "Make faces at each other: happy, sad, angry, silly. 'That's MY face. What face are YOU making?' Two different faces — two different feelings. That's OK.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Add a mirror — 3 faces at once (you, them, reflection)", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Start with just happy and sad", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Use emoji cards instead of real faces", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "😊", Label: "I made my own face", SkillFocus: "emotional_differentiation"},
			{Icon: "👀", Label: "I spotted someone else's feeling", SkillFocus: "emotion_recognition"},
			{Icon: "🪞", Label: "Our faces were different — and that's OK", SkillFocus: "emotional_boundary"},
			{Icon: "🎭", Label: "I tried lots of feelings", SkillFocus: "emotional_range"},
		},
	},
	{
		ID: "your_my_feelings_cycle2", SpectrumID: 8, ThemeName: "Your Feelings, My Feelings", Cycle: 2,
		CycleLabel: "The Feeling Catcher", CycleDescription: "Noticing when someone else's feeling 'jumps' into you — and putting it back.",
		TryThisWhen: []string{"They can name others' feelings", "They absorb emotions without knowing it", "You want to build emotional separation"},
		SkipIf:      []string{"They're already carrying someone's feelings", "Empathy is very low today", "They need connection, not boundaries"},
		Activities: []ThemeActivity{
			{
				ActivityID: "feelings_c2_catcher", Name: "The Feeling Shield", Mode: "with_you", Duration: "3 min",
				Instruction: "Pretend to throw a 'sad feeling' at them. They hold up an invisible shield. 'Your sadness bounced off! I noticed it, but it's not mine to carry.'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Let them throw feelings at you — you demonstrate the shield", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Use a real cushion as the shield", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Just practise: 'That's your feeling, not mine' — no throwing", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🛡️", Label: "I used my feeling shield", SkillFocus: "emotional_boundary"},
			{Icon: "💬", Label: "I said 'that's yours, not mine'", SkillFocus: "emotional_differentiation"},
			{Icon: "🤗", Label: "I noticed AND kept my balance", SkillFocus: "empathic_regulation"},
			{Icon: "⚡", Label: "A feeling bounced off!", SkillFocus: "emotional_resilience"},
		},
	},
	{
		ID: "your_my_feelings_cycle3", SpectrumID: 8, ThemeName: "Your Feelings, My Feelings", Cycle: 3,
		CycleLabel: "The Kind Thought", CycleDescription: "Caring without absorbing. Sending help without carrying the weight.",
		TryThisWhen: []string{"Shield is established", "They want to help but get overwhelmed", "You want to build sustainable empathy"},
		SkipIf:      []string{"They need to learn to notice first", "Empathy fatigue is present", "They need to receive care today"},
		Activities: []ThemeActivity{
			{
				ActivityID: "feelings_c3_kindthought", Name: "The Kind Post", Mode: "with_you", Duration: "5 min",
				Instruction: "Someone is sad (friend, teddy, character in a book). Instead of fixing: 'Let's send them a kind thought. What would you tell them?' Draw it or say it. Help without carrying.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Make a 'kind post box' for the week — collect kind thoughts", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "You say the kind thought first — they repeat", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Send the kind thought to a teddy, not a real person", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "💌", Label: "I sent a kind thought", SkillFocus: "compassion"},
			{Icon: "💭", Label: "I helped without carrying", SkillFocus: "sustainable_empathy"},
			{Icon: "🫂", Label: "I knew what they needed", SkillFocus: "perspective_taking"},
			{Icon: "✨", Label: "My thought made a difference", SkillFocus: "prosocial_efficacy"},
		},
	},

	// =========================================================================
	// SPECTRUM 9: ORDERLINESS — "Tidy and Messy"
	// =========================================================================
	{
		ID: "tidy_and_messy_cycle1", SpectrumID: 9, ThemeName: "Tidy and Messy", Cycle: 1,
		CycleLabel: "Exploring Both", CycleDescription: "Mess is allowed. Tidy is allowed. Both have a place.",
		TryThisWhen: []string{"Everything must be 'just so'", "Mess causes real distress", "Tidying is a constant battle"},
		SkipIf:      []string{"You need a clean house right now", "They've just been told off for mess", "Structure is needed for safety"},
		Activities: []ThemeActivity{
			{
				ActivityID: "tidy_messy_c1_explore", Name: "The Messy Minute", Mode: "with_you", Duration: "3 min",
				Instruction: "Set a timer for 60 seconds. Make the BIGGEST mess you can (cushions, toys, clothes). Timer stops. Look at it together. 'That's a LOT of mess.' Then: 'Now let's tidy for 60 seconds.'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Extend to 2 minutes each way", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Just messy one thing — dump one box", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "You make the mess — they just watch", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🌪️", Label: "I made a big mess!", SkillFocus: "structure_tolerance"},
			{Icon: "✨", Label: "I tidied it all up", SkillFocus: "task_completion"},
			{Icon: "🔄", Label: "I went from messy to tidy", SkillFocus: "transition_management"},
			{Icon: "😌", Label: "The mess was OK for a minute", SkillFocus: "discomfort_tolerance"},
		},
	},
	{
		ID: "tidy_and_messy_cycle2", SpectrumID: 9, ThemeName: "Tidy and Messy", Cycle: 2,
		CycleLabel: "My Way", CycleDescription: "There's more than one way to organise — finding their own system.",
		TryThisWhen: []string{"They can tolerate mess", "They only accept YOUR version of tidy", "You want to build their own ordering system"},
		SkipIf:      []string{"Mess tolerance isn't there yet", "You need things done a specific way", "They're feeling controlled"},
		Activities: []ThemeActivity{
			{
				ActivityID: "tidy_messy_c2_myway", Name: "Tidy YOUR Way", Mode: "with_you", Duration: "5 min",
				Instruction: "Dump a box of mixed toys. 'Tidy these — but YOUR way. There's no wrong way to organise.' Watch what system they invent. Name it: 'You sorted by colour! That's your system.'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Compare systems — you tidy half one way, they tidy half another", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Give them 2 categories: 'big and small'", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Start with just 5 items — not a whole box", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🧠", Label: "I invented my own system", SkillFocus: "classification"},
			{Icon: "🎨", Label: "My way was different from yours", SkillFocus: "independent_thinking"},
			{Icon: "📐", Label: "I found a pattern", SkillFocus: "pattern_recognition"},
			{Icon: "🤝", Label: "Both our ways worked", SkillFocus: "perspective_flexibility"},
		},
	},
	{
		ID: "tidy_and_messy_cycle3", SpectrumID: 9, ThemeName: "Tidy and Messy", Cycle: 3,
		CycleLabel: "The Controlled Chaos", CycleDescription: "Being OK with 'good enough' — not everything needs to be perfect.",
		TryThisWhen: []string{"They have their own system", "Perfectionism is emerging", "You want to build 'good enough' tolerance"},
		SkipIf:      []string{"They need a win today", "OCD-like patterns are present (consult professional)", "Mess is genuinely unsafe"},
		Activities: []ThemeActivity{
			{
				ActivityID: "tidy_messy_c3_goodenough", Name: "The 80% Tidy", Mode: "with_you", Duration: "5 min",
				Instruction: "Tidy the room together. When it's ALMOST done, stop. 'Is this good enough? Not perfect — good enough.' Leave one thing out of place. Live with it.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Leave TWO things out of place — raise the tolerance", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Leave something very small out of place — a single crayon", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Tell them YOU are leaving something out — it's your mess, not theirs", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "✅", Label: "Good enough IS enough", SkillFocus: "perfectionism_tolerance"},
			{Icon: "😎", Label: "I left one thing and survived", SkillFocus: "discomfort_tolerance"},
			{Icon: "🧘", Label: "I noticed the mess and didn't fix it", SkillFocus: "impulse_control"},
			{Icon: "🌟", Label: "Perfect is not the goal", SkillFocus: "growth_mindset"},
		},
	},

	// =========================================================================
	// SPECTRUM 10: RESPONSIBILITY THRESHOLD — "My Fault, Your Fault"
	// =========================================================================
	{
		ID: "my_fault_your_fault_cycle1", SpectrumID: 10, ThemeName: "My Fault, Your Fault", Cycle: 1,
		CycleLabel: "Sorting the Pieces", CycleDescription: "Learning to see whose part is whose — without blame or guilt.",
		TryThisWhen: []string{"'It wasn't me' is the default", "They take on blame for everything", "Consequences are always someone else's fault"},
		SkipIf:      []string{"They've just been told off", "A genuine injustice happened to them", "Shame is already high"},
		Activities: []ThemeActivity{
			{
				ActivityID: "responsibility_c1_sorting", Name: "Whose Bit Was That?", Mode: "with_you", Duration: "5 min",
				Instruction: "After any small incident (spilled drink, knocked tower), sit together. Draw a circle. 'What happened? Who did what?' Split it into pieces like a pie. 'This bit was yours. This bit was the table being wobbly. This bit was the cup being too full.' Name each piece without blame.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Try it with a story — 'whose bit was it in this book?'", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Use toys to act it out — teddy did this, bunny did that", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "YOU do the sorting — they just watch and nod", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🔍", Label: "I found my bit", SkillFocus: "ownership_activation"},
			{Icon: "⚖️", Label: "I sorted it fairly", SkillFocus: "consequence_prediction"},
			{Icon: "🗣️", Label: "I said what I did", SkillFocus: "accountability_language"},
			{Icon: "🧩", Label: "Everyone had a piece", SkillFocus: "delegation_comfort"},
		},
	},
	{
		ID: "my_fault_your_fault_cycle2", SpectrumID: 10, ThemeName: "My Fault, Your Fault", Cycle: 2,
		CycleLabel: "Fixing My Part", CycleDescription: "Taking action on your own piece — not carrying everyone else's.",
		TryThisWhen: []string{"They can identify their part", "They try to fix everyone's problems", "They need practice with repair not guilt"},
		SkipIf:      []string{"Blame is still reactive", "They're carrying too much already", "A real grievance needs adult resolution"},
		Activities: []ThemeActivity{
			{
				ActivityID: "responsibility_c2_repair", Name: "Fix Your Bit", Mode: "with_you", Duration: "5 min",
				Instruction: "When something goes wrong, sort the pieces (Cycle 1). Then: 'Your bit was [specific thing]. What could fix just that bit?' They do ONE repair action. Not sorry — REPAIR. 'You knocked it over. You rebuilt that part.'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Let them suggest repairs for fictional characters too", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Offer two repair options — they pick one", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Model it — 'Watch me fix MY bit first'", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🔧", Label: "I fixed my bit", SkillFocus: "ownership_activation"},
			{Icon: "🤲", Label: "I let someone else fix theirs", SkillFocus: "delegation_comfort"},
			{Icon: "💪", Label: "Repair is better than sorry", SkillFocus: "accountability_language"},
			{Icon: "🎯", Label: "I did just my part", SkillFocus: "responsibility_release"},
		},
	},
	{
		ID: "my_fault_your_fault_cycle3", SpectrumID: 10, ThemeName: "My Fault, Your Fault", Cycle: 3,
		CycleLabel: "The Detective", CycleDescription: "Finding facts not villains — consequence mapping without shame.",
		TryThisWhen: []string{"They can repair their bit", "Blame patterns are softening", "You want to build prediction skills"},
		SkipIf:      []string{"Emotional regulation is shaky today", "They need comfort not challenge", "The situation is genuinely unfair"},
		Activities: []ThemeActivity{
			{
				ActivityID: "responsibility_c3_detective", Name: "The Consequence Detective", Mode: "with_you", Duration: "5 min",
				Instruction: "Before an activity: 'If we do X, what might happen?' After: 'What DID happen? Was it what we predicted?' No blame — detectives find facts. 'The evidence suggests... the tower fell because it was too tall AND someone bumped the table.'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Keep a 'detective notebook' of predictions vs results", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Give them two possible outcomes to pick from", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Be the detective together — shared investigation", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🕵️", Label: "I found the facts", SkillFocus: "consequence_prediction"},
			{Icon: "📋", Label: "I predicted what happened", SkillFocus: "consequence_prediction"},
			{Icon: "🧠", Label: "No villains — just facts", SkillFocus: "accountability_language"},
			{Icon: "🔮", Label: "I guessed right!", SkillFocus: "consequence_prediction"},
		},
	},

	// =========================================================================
	// SPECTRUM 11: LOSS SENSITIVITY — "Keeping and Letting Go"
	// =========================================================================
	{
		ID: "keeping_letting_go_cycle1", SpectrumID: 11, ThemeName: "Keeping and Letting Go", Cycle: 1,
		CycleLabel: "Things Go Home", CycleDescription: "Nothing disappears — things go to their place. Building security in transitions.",
		TryThisWhen: []string{"Tidy-up triggers meltdowns", "They hoard or cling to objects", "Transitions feel like theft"},
		SkipIf:      []string{"They've just lost something real", "A pet or person has died recently", "Change is already overwhelming"},
		Activities: []ThemeActivity{
			{
				ActivityID: "loss_c1_gohome", Name: "Things Go Home", Mode: "with_you", Duration: "5 min",
				Instruction: "Pick 3 toys. Name where each one 'lives.' 'Teddy lives on the shelf. Lego lives in the box. Book lives on the table.' Now: 'Time for things to go home.' Each toy goes to its HOME — not 'away.' 'Teddy is going home now. Wave goodnight.'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "They choose where new things live — naming the home", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Just ONE thing goes home — Teddy gets a named shelf", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "YOU send your things home — they watch", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🏠", Label: "Everything has a home", SkillFocus: "secure_attachment"},
			{Icon: "👋", Label: "I said goodnight to my things", SkillFocus: "transition_rituals"},
			{Icon: "🧸", Label: "Teddy went home safely", SkillFocus: "secure_attachment"},
			{Icon: "😊", Label: "Going home is not going away", SkillFocus: "abundance_awareness"},
		},
	},
	{
		ID: "keeping_letting_go_cycle2", SpectrumID: 11, ThemeName: "Keeping and Letting Go", Cycle: 2,
		CycleLabel: "Lending and Returning", CycleDescription: "Things can leave and come back. Trust the cycle.",
		TryThisWhen: []string{"They understand 'things go home'", "Sharing is hard because it feels like losing", "They need to learn that lending isn't giving away"},
		SkipIf:      []string{"Trust is broken", "Someone took something and didn't return it", "They're in a possessive phase for safety reasons"},
		Activities: []ThemeActivity{
			{
				ActivityID: "loss_c2_lending", Name: "The Library Game", Mode: "with_you", Duration: "5 min",
				Instruction: "Choose one toy to 'lend' to you. 'I'm borrowing Bear for 2 minutes. He's visiting my house.' Timer on. Bear comes back. 'Bear came back! Lending isn't losing.' Gradually extend the time.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "They lend to a sibling or friend — same return guarantee", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Keep the timer visible — they can SEE when it comes back", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "YOU lend YOUR thing to them first — model the trust", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "📚", Label: "I lent something and it came back", SkillFocus: "secure_attachment"},
			{Icon: "🔄", Label: "Lending is not losing", SkillFocus: "abundance_awareness"},
			{Icon: "⏰", Label: "I waited and it returned", SkillFocus: "transition_rituals"},
			{Icon: "🤝", Label: "I trusted someone with my thing", SkillFocus: "secure_attachment"},
		},
	},
	{
		ID: "keeping_letting_go_cycle3", SpectrumID: 11, ThemeName: "Keeping and Letting Go", Cycle: 3,
		CycleLabel: "The Memory Box", CycleDescription: "Some things leave for real. We keep the memory, not the thing.",
		TryThisWhen: []string{"They can lend and receive back", "A real goodbye is coming (end of term, moving)", "They need to process genuine loss at their level"},
		SkipIf:      []string{"Grief is raw and unprocessed", "Professional support is needed", "The loss is too big for a game"},
		Activities: []ThemeActivity{
			{
				ActivityID: "loss_c3_memorybox", Name: "The Memory Box", Mode: "with_you", Duration: "10 min",
				Instruction: "Choose something that's leaving (outgrown clothes, a broken toy, end-of-term project). Take a photo. Draw it. Tell its story. Put the photo/drawing in a special box. 'The thing is going, but the memory stays. You can open this box whenever you want.'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Make it a regular ritual — seasonal memory boxing", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Just the photo — no drawing needed", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Do it with YOUR thing first — show that adults feel loss too", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "📸", Label: "I kept the memory", SkillFocus: "grief_processing"},
			{Icon: "📦", Label: "My memory box has a new thing", SkillFocus: "transition_rituals"},
			{Icon: "💛", Label: "It's OK to be sad", SkillFocus: "grief_processing"},
			{Icon: "🌱", Label: "Something left, something new can come", SkillFocus: "abundance_awareness"},
		},
	},

	// =========================================================================
	// SPECTRUM 12: LIBIDO — "Wanting and Waiting"
	// =========================================================================
	{
		ID: "wanting_and_waiting_cycle1", SpectrumID: 12, ThemeName: "Wanting and Waiting", Cycle: 1,
		CycleLabel: "Naming What I Want", CycleDescription: "Desire is not wrong. Learning to name it is a skill.",
		TryThisWhen: []string{"They grab without asking", "They can't articulate what they want", "Desire comes out as frustration"},
		SkipIf:      []string{"They're overwhelmed", "Basic needs aren't met", "Emotional regulation is fragile today"},
		Activities: []ThemeActivity{
			{
				ActivityID: "libido_c1_naming", Name: "I Want...", Mode: "with_you", Duration: "3 min",
				Instruction: "Before snack or activity, practise: 'What do you want? Say it out loud. I want the red cup. I want to play outside.' No judgement. Every stated want gets acknowledged: 'You want X. I heard you.' Then: 'Can you have it now, or do you need to wait?'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Make a 'want list' — things I want today, this week, someday", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Offer two choices — 'Do you want A or B?'", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "YOU say what you want first — model the language", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🗣️", Label: "I said what I wanted", SkillFocus: "desire_articulation"},
			{Icon: "👂", Label: "Someone heard me", SkillFocus: "desire_articulation"},
			{Icon: "🤔", Label: "I thought about what I want", SkillFocus: "desire_articulation"},
			{Icon: "💬", Label: "Wanting things is OK", SkillFocus: "desire_articulation"},
		},
	},
	{
		ID: "wanting_and_waiting_cycle2", SpectrumID: 12, ThemeName: "Wanting and Waiting", Cycle: 2,
		CycleLabel: "The Waiting Game", CycleDescription: "The gap between wanting and getting is where self-regulation lives.",
		TryThisWhen: []string{"They can name what they want", "Instant gratification is the default", "They need practice with delay"},
		SkipIf:      []string{"They're hungry or tired", "The wait would be cruel not educational", "They need a win not a challenge"},
		Activities: []ThemeActivity{
			{
				ActivityID: "libido_c2_waiting", Name: "The 30-Second Wait", Mode: "with_you", Duration: "3 min",
				Instruction: "When they want something (snack, toy, screen): 'Yes, you can have it. In 30 seconds.' Count together. They get the thing. Gradually extend: 1 minute, 2 minutes. The key: they ALWAYS get it. The wait is the skill, not the denial.",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "They choose the wait time — 'How long can you wait?'", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "10-second wait with a visual countdown", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "Wait together — 'We'll both wait. I want the biscuit too.'", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "⏳", Label: "I waited and I got it", SkillFocus: "delayed_gratification"},
			{Icon: "💪", Label: "Waiting made me stronger", SkillFocus: "impulse_channeling"},
			{Icon: "🎯", Label: "I chose how long to wait", SkillFocus: "delayed_gratification"},
			{Icon: "🧘", Label: "I was patient", SkillFocus: "delayed_gratification"},
		},
	},
	{
		ID: "wanting_and_waiting_cycle3", SpectrumID: 12, ThemeName: "Wanting and Waiting", Cycle: 3,
		CycleLabel: "Ask First", CycleDescription: "Before you take, touch, or start — check. Consent is a muscle.",
		TryThisWhen: []string{"They can wait", "Physical boundaries need work", "They take/touch without asking"},
		SkipIf:      []string{"Safeguarding concerns need professional input", "They're too young for this abstraction", "Trust is broken"},
		Activities: []ThemeActivity{
			{
				ActivityID: "libido_c3_askfirst", Name: "May I?", Mode: "with_you", Duration: "5 min",
				Instruction: "Play a game where EVERYTHING requires 'May I?' — 'May I pick up the crayon? May I sit here? May I give you a hug?' Practice hearing 'yes' AND 'no.' Both are OK. 'You said no to the hug. That's your right. Shall I wave instead?'",
				Variations: []ActivityVariation{
					{Condition: "If they loved this", Suggestion: "Extend to real life — 'May I?' before touching someone's things all day", AdjustmentType: "escalate"},
					{Condition: "If they struggled", Suggestion: "Just 3 rounds of 'May I?' with clear yes/no outcomes", AdjustmentType: "simplify"},
					{Condition: "If they were overwhelmed", Suggestion: "THEY are the gatekeeper — they say yes or no to YOUR requests", AdjustmentType: "scaffold"},
				},
			},
		},
		StickerChoices: []StickerChoice{
			{Icon: "🙋", Label: "I asked first", SkillFocus: "consent_navigation"},
			{Icon: "✋", Label: "I heard 'no' and that was OK", SkillFocus: "boundary_awareness"},
			{Icon: "🤝", Label: "We both said yes", SkillFocus: "consent_navigation"},
			{Icon: "🛡️", Label: "I said no and they listened", SkillFocus: "boundary_awareness"},
		},
	},
}

// =============================================================================
// Lookup Functions
// =============================================================================

// GetThemeWeekByID returns a ThemeWeek by its ID string.
func GetThemeWeekByID(id string) *ThemeWeek {
	for i := range ThemeWeeks {
		if ThemeWeeks[i].ID == id {
			return &ThemeWeeks[i]
		}
	}
	return nil
}

// GetThemeWeeksForSpectrum returns all ThemeWeeks for a given spectrum.
func GetThemeWeeksForSpectrum(spectrumID int) []ThemeWeek {
	var result []ThemeWeek
	for _, tw := range ThemeWeeks {
		if tw.SpectrumID == spectrumID {
			result = append(result, tw)
		}
	}
	return result
}

// GetThemeWeekForCycle returns the specific ThemeWeek for a spectrum at a cycle.
func GetThemeWeekForCycle(spectrumID, cycle int) *ThemeWeek {
	for i := range ThemeWeeks {
		if ThemeWeeks[i].SpectrumID == spectrumID && ThemeWeeks[i].Cycle == cycle {
			return &ThemeWeeks[i]
		}
	}
	return nil
}

// findThemeWeek is an internal helper used by the suggestion engine.
func findThemeWeek(spectrumID, cycle int) *ThemeWeek {
	return GetThemeWeekForCycle(spectrumID, cycle)
}

// =============================================================================
// Helpers
// =============================================================================

func pluralise(n int, word string) string {
	if n == 1 {
		return "1 " + word
	}
	return string(rune('0'+n)) + " " + word + "s"
}

func capitalise(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:] // ASCII uppercase — safe for our English strings
}

func themeIcon(spectrumID int) string {
	icons := map[int]string{
		1:  "🪐",  // Near and Far
		2:  "🔊",  // Loud and Quiet
		3:  "🌊",  // Big Feelings, Small Feelings
		4:  "🛑",  // Yes and No
		5:  "🎁",  // Mine and Yours
		6:  "🦁",  // Try and Wait
		7:  "📜",  // Same and Different
		8:  "🪞",  // Your Feelings, My Feelings
		9:  "🌪️", // Tidy and Messy
		10: "⚖️", // My Fault, Your Fault
		11: "🔒",  // Keeping and Letting Go
		12: "🔥",  // Wanting and Waiting
	}
	if icon, ok := icons[spectrumID]; ok {
		return icon
	}
	return "⭐"
}
