package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

const version = "1.0.0"

type TeachingPattern struct {
	PatternID               string   `json:"patternId"`
	Name                    string   `json:"name"`
	Description             string   `json:"description"`
	PatternType             string   `json:"patternType"`
	Steps                   []string `json:"steps"`
	Variations              []string `json:"variations"`
	Adaptations             string   `json:"adaptations"`
	BestUsedWhen            []string `json:"bestUsedWhen"`
	AvoidWhen               []string `json:"avoidWhen"`
	CombinesWellWith        []string `json:"combinesWellWith"`
	TargetBarrierCategories []string `json:"targetBarrierCategories"`
	SuccessIndicators       []string `json:"successIndicators"`
	SourceReference         string   `json:"sourceReference"`
}

func main() {
	log.Printf("📚 Teaching Patterns Import Tool v%s", version)

	var patterns []TeachingPattern
	var err error

	if len(os.Args) > 1 {
		patterns, err = loadPatternsFromFile(os.Args[1])
		if err != nil {
			log.Fatalf("Failed to load patterns: %v", err)
		}
	} else {
		log.Println("Using embedded patterns from humanOS documentation...")
		patterns = getEmbeddedPatterns()
	}

	log.Printf("✅ Loaded %d teaching patterns", len(patterns))

	client, err := connectToWeaviate()
	if err != nil {
		log.Fatalf("Failed to connect to Weaviate: %v", err)
	}
	log.Println("✅ Connected to Weaviate")

	successCount := 0
	errorCount := 0

	for _, pattern := range patterns {
		err := importPattern(client, pattern)
		if err != nil {
			log.Printf("   ❌ Error importing pattern %s: %v", pattern.PatternID, err)
			errorCount++
		} else {
			log.Printf("   ✅ Imported: %s", pattern.Name)
			successCount++
		}
	}

	log.Printf("\n🎉 Import complete! ✅ %d patterns imported, ❌ %d errors", successCount, errorCount)
}

func loadPatternsFromFile(path string) ([]TeachingPattern, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var patterns []TeachingPattern
	if err := json.Unmarshal(data, &patterns); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return patterns, nil
}

func getEmbeddedPatterns() []TeachingPattern {
	return []TeachingPattern{
		{
			PatternID:               "voltage_reduction_familiarity",
			Name:                    "Voltage Reduction Through Familiarity",
			Description:             "Start with familiar ground to lower emotional voltage before introducing new content.",
			PatternType:             "engagement",
			Steps:                   []string{"Acknowledge student with warmth", "Quick recap of last session", "Connect to student interests", "Lower voltage before new content", "Select a micro-success opportunity"},
			Variations:              []string{"Interest-based entry point", "Previous success reminder", "Comfort topic bridge"},
			Adaptations:             "Adjust familiarity depth based on student's anxiety level.",
			BestUsedWhen:            []string{"Starting a new session", "Introducing difficult content", "Student shows anxiety or resistance"},
			AvoidWhen:               []string{"Student is eager to learn new material", "Time is very limited"},
			CombinesWellWith:        []string{"engagement_hook", "micro_success_generation"},
			TargetBarrierCategories: []string{"initiation_barrier", "chronic_barrier"},
			SuccessIndicators:       []string{"Student relaxes visibly", "Engagement increases", "Student makes connections unprompted"},
			SourceReference:         "humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md",
		},
		{
			PatternID:               "engagement_hook",
			Name:                    "Engagement Hook (Welcome Pattern)",
			Description:             "Set positive tone immediately with warm, personalized greeting.",
			PatternType:             "engagement",
			Steps:                   []string{"Personalized greeting using interests", "Warm acknowledgment of presence", "Set positive emotional tone", "Reference previous successes"},
			Variations:              []string{"Interest-based greeting", "Achievement reminder opening", "Mood check-in start"},
			Adaptations:             "Match energy to student's current state.",
			BestUsedWhen:            []string{"Beginning of any session", "Student seems disengaged"},
			AvoidWhen:               []string{"Student is in crisis mode"},
			CombinesWellWith:        []string{"voltage_reduction_familiarity"},
			TargetBarrierCategories: []string{"initiation_barrier"},
			SuccessIndicators:       []string{"Student responds positively", "Engagement metrics increase"},
			SourceReference:         "humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md",
		},
		{
			PatternID:               "question_killer_game",
			Name:                    "Question Killer Game (Safe Structured Questioning)",
			Description:             "No-hands-up alternative creating safety through structure.",
			PatternType:             "scaffolding",
			Steps:                   []string{"Start with impossibly easy question", "Track success streak", "If struggling, guide toward guaranteed success", "If 3+ successes, slightly increase difficulty", "Wrong answers OK if genuine attempt"},
			Variations:              []string{"Difficulty progression", "Topic-specific adaptation", "Confidence-based adjustment"},
			Adaptations:             "Adjust starting difficulty based on ETP profile.",
			BestUsedWhen:            []string{"Student avoids answering questions", "Building participation confidence", "Student uses 'I don't know' as avoidance"},
			AvoidWhen:               []string{"Student is already highly engaged", "Deep discussion is needed"},
			CombinesWellWith:        []string{"micro_success_generation", "voltage_reduction_familiarity"},
			TargetBarrierCategories: []string{"initiation_barrier", "chronic_barrier"},
			SuccessIndicators:       []string{"Student attempts answers", "Reduced 'I don't know' responses", "Increasing answer quality"},
			SourceReference:         "humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md",
		},
		{
			PatternID:               "micro_success_generation",
			Name:                    "Micro-Success Generation",
			Description:             "Create small wins to build momentum.",
			PatternType:             "scaffolding",
			Steps:                   []string{"Find smallest possible positive action", "Acknowledge any engagement", "Celebrate completion (not perfection)", "Build on success to next micro-goal", "Gradually raise expectations"},
			Variations:              []string{"Task completion celebration", "Effort recognition", "Progress acknowledgment"},
			Adaptations:             "For confrontational students, celebrate even staying in conversation.",
			BestUsedWhen:            []string{"Student won't start tasks", "After failure or setback", "Building initial momentum"},
			AvoidWhen:               []string{"Student is already performing well", "Would feel patronizing"},
			CombinesWellWith:        []string{"question_killer_game", "game_access_incentive"},
			TargetBarrierCategories: []string{"initiation_barrier", "chronic_barrier"},
			SuccessIndicators:       []string{"Student completes small tasks", "Positive momentum builds", "Student initiates next step"},
			SourceReference:         "humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md",
		},
		{
			PatternID:               "game_access_incentive",
			Name:                    "Game-Access Incentive",
			Description:             "Use extrinsic motivation to get initial engagement.",
			PatternType:             "motivation",
			Steps:                   []string{"Offer game access reward for task completion", "Define clear completion criteria", "Generate time-limited unlock code", "Track reward-triggered vs self-initiated work", "Gradually increase work required per reward"},
			Variations:              []string{"Time-based rewards", "Task-quantity rewards", "Quality-threshold rewards"},
			Adaptations:             "Adjust reward frequency based on progress toward intrinsic motivation.",
			BestUsedWhen:            []string{"Student has no intrinsic motivation", "Need to break avoidance cycle", "Building initial work habits"},
			AvoidWhen:               []string{"Student is intrinsically motivated", "Would create dependency"},
			CombinesWellWith:        []string{"micro_success_generation", "ban_idk_response"},
			TargetBarrierCategories: []string{"initiation_barrier"},
			SuccessIndicators:       []string{"Student completes tasks for reward", "Work quality meets criteria", "Gradual reduction in reward-seeking"},
			SourceReference:         "humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md",
		},
		{
			PatternID:               "ban_idk_response",
			Name:                    "Ban 'I Don't Know' Responses",
			Description:             "Eliminate primary avoidance tactic by requiring attempts.",
			PatternType:             "avoidance_prevention",
			Steps:                   []string{"Establish rule: 'I don't know' is not acceptable", "Require attempt even if unsure", "Frame: 'Trying and being wrong teaches more'", "Wait in supportive silence for attempt", "Celebrate any genuine attempt"},
			Variations:              []string{"Gentle persistence", "Scaffold toward attempt", "Offer multiple choice if stuck"},
			Adaptations:             "For very anxious students, offer scaffolded options.",
			BestUsedWhen:            []string{"Student uses 'I don't know' repeatedly", "Avoidance pattern is established", "Need to reveal actual capability"},
			AvoidWhen:               []string{"Student is genuinely overwhelmed", "Would cause shame spiral"},
			CombinesWellWith:        []string{"question_killer_game", "micro_success_generation"},
			TargetBarrierCategories: []string{"initiation_barrier", "chronic_barrier"},
			SuccessIndicators:       []string{"Reduced 'I don't know' responses", "Student makes attempts", "Quality of attempts improves"},
			SourceReference:         "humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md",
		},
		{
			PatternID:               "shoulder_sitting",
			Name:                    "Shoulder Sitting (Intensive Proximity)",
			Description:             "Provide constant nearby presence to prevent hiding behavior.",
			PatternType:             "proximity_support",
			Steps:                   []string{"Maintain constant virtual presence", "Frequent check-ins", "Quick response to any disengagement", "Supportive monitoring without pressure"},
			Variations:              []string{"Active monitoring", "Gentle presence", "Responsive availability"},
			Adaptations:             "Adjust presence intensity based on student needs.",
			BestUsedWhen:            []string{"Student is a 'hider' or 'avoider'", "Silent avoidance pattern", "Student needs external structure"},
			AvoidWhen:               []string{"Student finds it intrusive", "Student is self-directed"},
			CombinesWellWith:        []string{"micro_success_generation", "question_killer_game"},
			TargetBarrierCategories: []string{"chronic_barrier"},
			SuccessIndicators:       []string{"Student engages when present", "Reduced hiding behavior", "Gradual increase in self-direction"},
			SourceReference:         "humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md",
		},
		{
			PatternID:               "relationship_building",
			Name:                    "Chat is Not All About Work",
			Description:             "Show interest in student as person, not just student.",
			PatternType:             "relationship_building",
			Steps:                   []string{"Ask about interests, hobbies, life", "Remember and reference personal details", "Balance non-work chat with work requests", "Show genuine curiosity about their world"},
			Variations:              []string{"Interest exploration", "Life check-ins", "Shared interest discovery"},
			Adaptations:             "Match conversation style to student preferences.",
			BestUsedWhen:            []string{"Building initial relationship", "Student is confrontational", "Need to build trust"},
			AvoidWhen:               []string{"Student wants to focus on work", "Time is very limited"},
			CombinesWellWith:        []string{"engagement_hook", "pattern_awareness_discussion"},
			TargetBarrierCategories: []string{"chronic_barrier"},
			SuccessIndicators:       []string{"Student shares personal information", "Relationship warmth increases", "Reduced confrontational behavior"},
			SourceReference:         "humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md",
		},
		{
			PatternID:               "pattern_awareness_discussion",
			Name:                    "Pattern Discussion (No Moral Judgment)",
			Description:             "Discuss behavioral patterns without good/bad framework.",
			PatternType:             "metacognition",
			Steps:                   []string{"Wait for relationship foundation", "Frame: 'Let's talk about patterns'", "Explain why pattern is harming them (not 'you're bad')", "Collaborate on how to fix it", "Focus on practical consequences"},
			Variations:              []string{"Observation sharing", "Impact discussion", "Solution brainstorming"},
			Adaptations:             "Time carefully based on relationship strength.",
			BestUsedWhen:            []string{"Strong relationship established", "Pattern is clearly harmful", "Student shows some self-awareness"},
			AvoidWhen:               []string{"Relationship not established", "Would feel like attack", "Student is in crisis"},
			CombinesWellWith:        []string{"relationship_building"},
			TargetBarrierCategories: []string{"chronic_barrier"},
			SuccessIndicators:       []string{"Student acknowledges pattern", "Collaborative problem-solving", "Self-correction attempts"},
			SourceReference:         "humanOS/docs/CLASSROOM_INTERACTION_PATTERNS.md",
		},
	}
}

func connectToWeaviate() (*weaviate.Client, error) {
	weaviateURL := os.Getenv("STUDENT_WEAVIATE_URL")
	if weaviateURL == "" {
		weaviateURL = "http://localhost:8088"
	}

	parsedURL, err := url.Parse(weaviateURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Weaviate URL: %w", err)
	}

	cfg := weaviate.Config{
		Host:   parsedURL.Host,
		Scheme: parsedURL.Scheme,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	_, err = client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Weaviate not ready: %w", err)
	}

	return client, nil
}

func importPattern(client *weaviate.Client, pattern TeachingPattern) error {
	stepsJSON, _ := json.Marshal(pattern.Steps)
	variationsJSON, _ := json.Marshal(pattern.Variations)

	data := map[string]interface{}{
		"patternId":               pattern.PatternID,
		"name":                    pattern.Name,
		"description":             pattern.Description,
		"patternType":             pattern.PatternType,
		"steps":                   string(stepsJSON),
		"variations":              string(variationsJSON),
		"adaptations":             pattern.Adaptations,
		"bestUsedWhen":            pattern.BestUsedWhen,
		"avoidWhen":               pattern.AvoidWhen,
		"combinesWellWith":        pattern.CombinesWellWith,
		"targetEtpProfiles":       "{}",
		"targetBarrierCategories": pattern.TargetBarrierCategories,
		"successIndicators":       pattern.SuccessIndicators,
		"sourceReference":         pattern.SourceReference,
	}

	_, err := client.Data().Creator().
		WithClassName("TeachingPattern").
		WithProperties(data).
		Do(context.Background())

	return err
}
