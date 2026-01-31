package intervention

import (
	"context"
	"fmt"
	"sort"
	"time"
)

// InterventionSelector selects optimal interventions based on student state
type InterventionSelector struct {
	weaviateClient interface{} // *weaviate.Client
}

// FrictionSignal represents real-time friction data from the mobile app
type FrictionSignal struct {
	Timestamp            time.Time `json:"timestamp"`
	HesitationTime       float64   `json:"hesitation_time"`        // Seconds of hesitation
	ErrorType            string    `json:"error_type"`             // Type of error made
	ErrorCount           int       `json:"error_count"`            // Number of errors
	SelfReportAnxiety    int       `json:"self_report_anxiety"`    // 0-10
	SelfReportConfidence int       `json:"self_report_confidence"` // 0-10
	ProgressionSpeed     float64   `json:"progression_speed"`      // 0-1, how fast through material
	AttemptDuration      float64   `json:"attempt_duration"`       // Seconds spent on current attempt
}

// BarrierMatch represents a detected barrier with confidence
type BarrierMatch struct {
	BarrierID       string  `json:"barrier_id"`
	BarrierName     string  `json:"barrier_name"`
	Confidence      float64 `json:"confidence"` // 0-1
	DetectionReason string  `json:"detection_reason"`
}

// LeverRecommendation represents a recommended intervention lever
type LeverRecommendation struct {
	LeverID              string  `json:"lever_id"`
	LeverName            string  `json:"lever_name"`
	LeverType            string  `json:"lever_type"`
	Score                float64 `json:"score"`                  // Combined effectiveness score
	HistoricalSuccess    float64 `json:"historical_success"`     // Past success rate for this ETP profile
	EnergyCost           string  `json:"energy_cost"`            // low, medium, high
	ExpectedVoltageShift int     `json:"expected_voltage_shift"` // How much it should reduce voltage
}

// SelectedIntervention is the final recommended intervention
type SelectedIntervention struct {
	InterventionID   string              `json:"intervention_id"`
	StudentID        string              `json:"student_id"`
	SkillID          string              `json:"skill_id"`
	SelectedLever    LeverRecommendation `json:"selected_lever"`
	TeachingPattern  string              `json:"teaching_pattern_id"`
	Actions          []ActionStep        `json:"actions"`
	Priority         int                 `json:"priority"`
	Urgency          string              `json:"urgency"`
	DetectedBarriers []BarrierMatch      `json:"detected_barriers"`
	VoltageState     int                 `json:"voltage_state"`
	Rationale        string              `json:"rationale"`
	CreatedAt        time.Time           `json:"created_at"`
}

// ActionStep is a specific action to take
type ActionStep struct {
	Order      int    `json:"order"`
	Action     string `json:"action"`
	Duration   string `json:"duration"`
	Checkpoint string `json:"checkpoint"` // How to know this step worked
}

// ETPProfile represents the student's ETP profile (simplified)
type ETPProfile struct {
	StudentID             string `json:"student_id"`
	ResponseToScarcity    int    `json:"response_to_scarcity"`
	ResponseToChallenge   int    `json:"response_to_challenge"`
	ResponseToFeedback    int    `json:"response_to_feedback"`
	ResponseToAuthority   int    `json:"response_to_authority"`
	ResponseToUncertainty int    `json:"response_to_uncertainty"`
	SocialEnergyBalance   int    `json:"social_energy_balance"`
	// ... other spectra
}

// SelectIntervention is the main entry point - selects optimal intervention
func (s *InterventionSelector) SelectIntervention(
	ctx context.Context,
	studentID string,
	skillID string,
	frictionSignals []FrictionSignal,
) (*SelectedIntervention, error) {

	// 1. Get student's current ETP profile
	etpProfile, err := s.getETPProfile(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ETP profile: %w", err)
	}

	// 2. Calculate current voltage state from friction signals
	voltageState := s.calculateVoltageState(frictionSignals)

	// 3. Detect active barriers based on friction patterns
	activeBarriers := s.detectBarriers(ctx, frictionSignals, etpProfile)

	// 4. Get skill-specific data
	skillData, err := s.getSkillData(ctx, skillID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill data: %w", err)
	}

	// 5. Match barriers to candidate intervention levers
	candidateLevers := s.findLeversForBarriers(ctx, activeBarriers, etpProfile)

	// 6. Filter by voltage state (can't use high-energy interventions if voltage is low)
	filteredLevers := s.filterByVoltage(candidateLevers, voltageState)

	// 7. Rank by historical effectiveness for this ETP profile + skill
	rankedLevers := s.rankByEffectiveness(ctx, filteredLevers, etpProfile, skillID)

	// 8. Select top lever
	if len(rankedLevers) == 0 {
		return nil, fmt.Errorf("no suitable interventions found")
	}
	selectedLever := rankedLevers[0]

	// 9. Get appropriate teaching pattern
	teachingPattern := s.selectTeachingPattern(ctx, skillID, selectedLever, etpProfile)

	// 10. Generate action steps
	actions := s.generateActionSteps(selectedLever, teachingPattern, skillData)

	// 11. Determine urgency
	urgency := s.determineUrgency(voltageState, activeBarriers)

	// 12. Build the intervention
	intervention := &SelectedIntervention{
		InterventionID:   fmt.Sprintf("int_%s_%d", studentID, time.Now().Unix()),
		StudentID:        studentID,
		SkillID:          skillID,
		SelectedLever:    selectedLever,
		TeachingPattern:  teachingPattern,
		Actions:          actions,
		Priority:         s.calculatePriority(voltageState, activeBarriers),
		Urgency:          urgency,
		DetectedBarriers: activeBarriers,
		VoltageState:     voltageState,
		Rationale:        s.buildRationale(selectedLever, activeBarriers, etpProfile),
		CreatedAt:        time.Now(),
	}

	// 13. Log selection for learning
	s.logInterventionChoice(ctx, intervention, frictionSignals)

	return intervention, nil
}

// ============================================================================
// STEP 1: Get ETP Profile
// ============================================================================

func (s *InterventionSelector) getETPProfile(ctx context.Context, studentID string) (*ETPProfile, error) {
	// TODO: Query Weaviate ETPProfile collection
	// For now, return placeholder
	return &ETPProfile{
		StudentID:             studentID,
		ResponseToScarcity:    0,
		ResponseToChallenge:   -20,
		ResponseToFeedback:    10,
		ResponseToAuthority:   -10,
		ResponseToUncertainty: -30,
		SocialEnergyBalance:   -20,
	}, nil
}

// ============================================================================
// STEP 2: Calculate Voltage State
// ============================================================================

func (s *InterventionSelector) calculateVoltageState(signals []FrictionSignal) int {
	if len(signals) == 0 {
		return 50 // Default mid-range
	}

	// Calculate voltage from multiple signals
	var totalVoltage float64

	for _, signal := range signals {
		signalVoltage := 50.0 // Baseline

		// Anxiety increases voltage
		signalVoltage += float64(signal.SelfReportAnxiety) * 5

		// Low confidence increases voltage
		signalVoltage += float64(10-signal.SelfReportConfidence) * 3

		// Hesitation indicates building voltage
		if signal.HesitationTime > 5 {
			signalVoltage += 10
		}

		// Errors increase voltage
		signalVoltage += float64(signal.ErrorCount) * 5

		// Slow progression suggests struggle
		if signal.ProgressionSpeed < 0.3 {
			signalVoltage += 15
		}

		totalVoltage += signalVoltage
	}

	avgVoltage := int(totalVoltage / float64(len(signals)))

	// Clamp to 0-100
	if avgVoltage < 0 {
		return 0
	}
	if avgVoltage > 100 {
		return 100
	}
	return avgVoltage
}

// ============================================================================
// STEP 3: Detect Barriers
// ============================================================================

func (s *InterventionSelector) detectBarriers(
	ctx context.Context,
	signals []FrictionSignal,
	profile *ETPProfile,
) []BarrierMatch {
	barriers := []BarrierMatch{}

	// Analyze signals for barrier patterns
	avgAnxiety := 0.0
	avgConfidence := 0.0
	totalHesitation := 0.0
	errorCount := 0

	for _, sig := range signals {
		avgAnxiety += float64(sig.SelfReportAnxiety)
		avgConfidence += float64(sig.SelfReportConfidence)
		totalHesitation += sig.HesitationTime
		errorCount += sig.ErrorCount
	}

	if len(signals) > 0 {
		avgAnxiety /= float64(len(signals))
		avgConfidence /= float64(len(signals))
	}

	// Check for "I don't know" / avoidance pattern
	if avgConfidence < 3 && totalHesitation > 10 {
		barriers = append(barriers, BarrierMatch{
			BarrierID:       "lack_of_motivation",
			BarrierName:     "Lack of Motivation / Won't Start",
			Confidence:      0.8,
			DetectionReason: "Low confidence + high hesitation indicates avoidance",
		})
	}

	// Check for anxiety barrier
	if avgAnxiety > 7 {
		barriers = append(barriers, BarrierMatch{
			BarrierID:       "anxiety_block",
			BarrierName:     "Anxiety Block",
			Confidence:      float64(avgAnxiety) / 10.0,
			DetectionReason: "High self-reported anxiety",
		})
	}

	// Check for silent avoider pattern
	if avgConfidence < 4 && errorCount == 0 && totalHesitation > 15 {
		barriers = append(barriers, BarrierMatch{
			BarrierID:       "silent_avoider",
			BarrierName:     "Silent Avoider / Hider",
			Confidence:      0.7,
			DetectionReason: "No attempts made, extended hesitation",
		})
	}

	// Sort by confidence
	sort.Slice(barriers, func(i, j int) bool {
		return barriers[i].Confidence > barriers[j].Confidence
	})

	return barriers
}

// ============================================================================
// STEP 4: Get Skill Data
// ============================================================================

func (s *InterventionSelector) getSkillData(ctx context.Context, skillID string) (map[string]interface{}, error) {
	// TODO: Query Weaviate CHISGSkill collection
	return map[string]interface{}{
		"skill_id":     skillID,
		"name":         "Placeholder Skill",
		"domain":       "FOCUS & TOOLS",
		"etp_spectrum": "Response to Challenge",
		"etp_vector":   "Toward Positive",
	}, nil
}

// ============================================================================
// STEP 5: Find Levers for Barriers
// ============================================================================

func (s *InterventionSelector) findLeversForBarriers(
	ctx context.Context,
	barriers []BarrierMatch,
	profile *ETPProfile,
) []LeverRecommendation {
	levers := []LeverRecommendation{}

	// Map barriers to effective levers
	leverMap := map[string][]LeverRecommendation{
		"lack_of_motivation": {
			{
				LeverID:    "game_access_incentive",
				LeverName:  "Game-Access Incentive",
				LeverType:  "extrinsic_motivation",
				EnergyCost: "low",
			},
			{
				LeverID:    "ban_idk_response",
				LeverName:  "Ban 'I Don't Know' Responses",
				LeverType:  "avoidance_prevention",
				EnergyCost: "medium",
			},
		},
		"anxiety_block": {
			{
				LeverID:    "voltage_reduction",
				LeverName:  "Voltage Reduction Through Familiarity",
				LeverType:  "voltage_reduction",
				EnergyCost: "low",
			},
			{
				LeverID:    "micro_success",
				LeverName:  "Micro-Success Generation",
				LeverType:  "scaffolding",
				EnergyCost: "low",
			},
		},
		"silent_avoider": {
			{
				LeverID:    "shoulder_sitting",
				LeverName:  "Shoulder Sitting (Intensive Proximity)",
				LeverType:  "proximity_support",
				EnergyCost: "high",
			},
			{
				LeverID:    "micro_success",
				LeverName:  "Micro-Success Generation",
				LeverType:  "scaffolding",
				EnergyCost: "low",
			},
		},
	}

	// Collect levers for all detected barriers
	for _, barrier := range barriers {
		if barrierLevers, ok := leverMap[barrier.BarrierID]; ok {
			for _, lever := range barrierLevers {
				lever.Score = barrier.Confidence * 0.5 // Initial score based on barrier confidence
				levers = append(levers, lever)
			}
		}
	}

	return levers
}

// ============================================================================
// STEP 6: Filter by Voltage
// ============================================================================

func (s *InterventionSelector) filterByVoltage(
	levers []LeverRecommendation,
	voltageState int,
) []LeverRecommendation {
	filtered := []LeverRecommendation{}

	for _, lever := range levers {
		// High voltage = can't use high-energy interventions
		if voltageState > 80 && lever.EnergyCost == "high" {
			continue
		}

		// Very high voltage = only low-energy interventions
		if voltageState > 90 && lever.EnergyCost != "low" {
			continue
		}

		filtered = append(filtered, lever)
	}

	// If nothing passes filter, return the lowest energy options
	if len(filtered) == 0 {
		for _, lever := range levers {
			if lever.EnergyCost == "low" {
				filtered = append(filtered, lever)
			}
		}
	}

	return filtered
}

// ============================================================================
// STEP 7: Rank by Effectiveness
// ============================================================================

func (s *InterventionSelector) rankByEffectiveness(
	ctx context.Context,
	levers []LeverRecommendation,
	profile *ETPProfile,
	skillID string,
) []LeverRecommendation {
	// TODO: Query historical effectiveness from ActionLog
	// For now, apply ETP-based adjustments

	for i := range levers {
		lever := &levers[i]

		// Adjust score based on ETP profile match
		switch lever.LeverType {
		case "extrinsic_motivation":
			// Works better for students with low internal drive
			if profile.ResponseToScarcity < 0 {
				lever.Score += 0.2
			}
		case "voltage_reduction":
			// Works better for anxious students
			if profile.ResponseToUncertainty < 0 {
				lever.Score += 0.3
			}
		case "proximity_support":
			// May not work for students who value independence
			if profile.SocialEnergyBalance < -30 {
				lever.Score -= 0.1
			}
		case "scaffolding":
			// Universal applicability
			lever.Score += 0.15
		}

		lever.HistoricalSuccess = 0.5 + lever.Score // Placeholder
	}

	// Sort by score
	sort.Slice(levers, func(i, j int) bool {
		return levers[i].Score > levers[j].Score
	})

	return levers
}

// ============================================================================
// STEP 8-12: Supporting Functions
// ============================================================================

func (s *InterventionSelector) selectTeachingPattern(
	ctx context.Context,
	skillID string,
	lever LeverRecommendation,
	profile *ETPProfile,
) string {
	// Map lever types to teaching patterns
	patternMap := map[string]string{
		"voltage_reduction":     "voltage_reduction_familiarity",
		"extrinsic_motivation":  "game_access_incentive",
		"avoidance_prevention":  "ban_idk_response",
		"scaffolding":           "micro_success_generation",
		"proximity_support":     "shoulder_sitting",
		"relationship_building": "relationship_building",
	}

	if pattern, ok := patternMap[lever.LeverType]; ok {
		return pattern
	}
	return "voltage_reduction_familiarity" // Default
}

func (s *InterventionSelector) generateActionSteps(
	lever LeverRecommendation,
	patternID string,
	skillData map[string]interface{},
) []ActionStep {
	// Generate action steps based on lever and pattern
	steps := []ActionStep{}

	switch lever.LeverID {
	case "voltage_reduction":
		steps = []ActionStep{
			{Order: 1, Action: "Start with 4-7-8 breathing exercise", Duration: "2 min", Checkpoint: "Student reports feeling calmer"},
			{Order: 2, Action: "Recap what we did last session", Duration: "3 min", Checkpoint: "Student recalls previous content"},
			{Order: 3, Action: "Connect to a topic they enjoy", Duration: "2 min", Checkpoint: "Student shows interest"},
			{Order: 4, Action: "Start with an easy win question", Duration: "1 min", Checkpoint: "Student answers correctly"},
		}
	case "game_access_incentive":
		steps = []ActionStep{
			{Order: 1, Action: "Offer 5-min game access for completing task", Duration: "30 sec", Checkpoint: "Student acknowledges deal"},
			{Order: 2, Action: "Define exactly what 'complete' means", Duration: "1 min", Checkpoint: "Criteria understood"},
			{Order: 3, Action: "Student works on task", Duration: "Variable", Checkpoint: "Task completed"},
			{Order: 4, Action: "Provide game access code", Duration: "30 sec", Checkpoint: "Reward delivered"},
		}
	case "micro_success":
		steps = []ActionStep{
			{Order: 1, Action: "Find smallest possible positive action", Duration: "1 min", Checkpoint: "Micro-goal identified"},
			{Order: 2, Action: "Student attempts micro-goal", Duration: "1-2 min", Checkpoint: "Any attempt made"},
			{Order: 3, Action: "Celebrate completion (not perfection)", Duration: "30 sec", Checkpoint: "Student acknowledged"},
			{Order: 4, Action: "Build to next micro-goal", Duration: "1 min", Checkpoint: "Momentum building"},
		}
	default:
		steps = []ActionStep{
			{Order: 1, Action: "Apply selected intervention", Duration: "5 min", Checkpoint: "Intervention started"},
			{Order: 2, Action: "Observe student response", Duration: "2 min", Checkpoint: "Response recorded"},
			{Order: 3, Action: "Adjust if needed", Duration: "Variable", Checkpoint: "Progress made"},
		}
	}

	return steps
}

func (s *InterventionSelector) determineUrgency(voltageState int, barriers []BarrierMatch) string {
	if voltageState > 90 {
		return "critical"
	}
	if voltageState > 75 || len(barriers) > 2 {
		return "high"
	}
	if voltageState > 50 || len(barriers) > 0 {
		return "medium"
	}
	return "low"
}

func (s *InterventionSelector) calculatePriority(voltageState int, barriers []BarrierMatch) int {
	priority := 5 // Default mid-priority

	// Higher voltage = higher priority
	priority += (voltageState - 50) / 20

	// More barriers = higher priority
	priority += len(barriers)

	// Clamp to 1-10
	if priority < 1 {
		return 1
	}
	if priority > 10 {
		return 10
	}
	return priority
}

func (s *InterventionSelector) buildRationale(
	lever LeverRecommendation,
	barriers []BarrierMatch,
	profile *ETPProfile,
) string {
	if len(barriers) == 0 {
		return fmt.Sprintf("Applying %s as general support", lever.LeverName)
	}

	barrierNames := ""
	for i, b := range barriers {
		if i > 0 {
			barrierNames += ", "
		}
		barrierNames += b.BarrierName
	}

	return fmt.Sprintf(
		"Detected barriers: %s. Selected '%s' (%s) based on ETP profile compatibility and historical effectiveness.",
		barrierNames,
		lever.LeverName,
		lever.LeverType,
	)
}

func (s *InterventionSelector) logInterventionChoice(
	ctx context.Context,
	intervention *SelectedIntervention,
	signals []FrictionSignal,
) {
	// TODO: Write to ActionLog collection for learning
	// This enables the system to learn which interventions work for which profiles
}

// ============================================================================
// LEARNING: Update effectiveness based on outcomes
// ============================================================================

// LearnFromOutcome updates the system based on intervention results
func (s *InterventionSelector) LearnFromOutcome(
	ctx context.Context,
	interventionID string,
	outcome string, // success, partial, blocked, abandoned
	skillImprovement float64, // 0-1, how much skill improved
	voltageReduction int, // How much voltage decreased
) error {
	// TODO:
	// 1. Update InterventionLibrary effectiveness scores
	// 2. Update SkillTeachingPattern mappings
	// 3. Refine ETP-intervention correlations
	// 4. Store in ActionLog for future queries
	return nil
}
