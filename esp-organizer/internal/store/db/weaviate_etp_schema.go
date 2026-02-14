package db

import (
	"context"
	"fmt"
	"log"

	wvmodels "github.com/weaviate/weaviate/entities/models"
)

// ETP Weaviate Schema
// This defines the collections for storing ETP slider/dashboard data
// for the HumanOS moral neutrality framework.
//
// Core Insight: ETPs are not personality flaws - they are positioning on biological spectra.
// The goal is not to keep everyone at 0 (middle), but to ensure agency to move sliders
// rather than having triggers move them involuntarily.

// CreateETPProfileClass creates the ETPProfile collection
// Stores a student's 9-slider dashboard configuration + 2 global moderators
func CreateETPProfileClass(ctx context.Context) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	className := "ETPProfile"
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Class %s already exists", className)
		return nil
	}

	class := &wvmodels.Class{
		Class:       className,
		Description: "Student ETP profile with 9 core spectrum slider positions + global moderators",
		Vectorizer:  "none", // No vectorization needed for structured data
		Properties: []*wvmodels.Property{
			{
				Name:        "student_id",
				DataType:    []string{"text"},
				Description: "Unique student identifier",
			},
			{
				Name:        "mongo_id",
				DataType:    []string{"text"},
				Description: "MongoDB ObjectID reference",
			},
			// 8 CORE ETP Slider Positions (-2 to +2)
			{
				Name:        "social_gravity",
				DataType:    []string{"number"},
				Description: "Independent (-2) to Cohesive (+2) - social interaction voltage",
			},
			{
				Name:        "energy_directionality",
				DataType:    []string{"number"},
				Description: "Inward (-2) to Outward (+2) - energy flow direction",
			},
			{
				Name:        "voltage_sensitivity",
				DataType:    []string{"number"},
				Description: "Insulated (-2) to Conductive (+2) - emotional current handling",
			},
			{
				Name:        "threat_response",
				DataType:    []string{"number"},
				Description: "Passive (-2) to Aggressive (+2) - threat confrontation style",
			},
			{
				Name:        "care_response",
				DataType:    []string{"number"},
				Description: "Detached (-2) to Nurturing (+2) - caregiving tendency",
			},
			{
				Name:        "risk_tolerance",
				DataType:    []string{"number"},
				Description: "Averse (-2) to Seeking (+2) - risk appetite",
			},
			{
				Name:        "integrity_logic",
				DataType:    []string{"number"},
				Description: "Relativistic (-2) to Absolutist (+2) - moral reasoning style",
			},
			{
				Name:        "mirror_neuron_tuning",
				DataType:    []string{"number"},
				Description: "Selective (-2) to Absorbent (+2) - empathy sensitivity",
			},
			{
				Name:        "orderliness",
				DataType:    []string{"number"},
				Description: "Flexible (-2) to Ordered (+2) - preference for structure vs spontaneity",
			},
			// Global Moderators
			{
				Name:        "pilot_strength",
				DataType:    []string{"number"},
				Description: "0-1 score of executive function capacity to move sliders deliberately",
			},
			{
				Name:        "current_load",
				DataType:    []string{"number"},
				Description: "0-1 score of stress/depletion (higher = less range of motion)",
			},
			// Meta fields
			{
				Name:        "overall_fluidity",
				DataType:    []string{"number"},
				Description: "0-1 score of how fluid all sliders are (vs stuck)",
			},
			{
				Name:        "assessment_date",
				DataType:    []string{"text"},
				Description: "ISO date of last assessment",
			},
		},
	}

	err = weaviateClient.Schema().ClassCreator().WithClass(class).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class %s: %w", className, err)
	}

	log.Printf("✅ Created Weaviate class: %s", className)
	return nil
}

// CreateSliderStateClass creates the SliderState collection
// Tracks the state of individual sliders (fluid, sticky, stuck, hijacked)
func CreateSliderStateClass(ctx context.Context) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	className := "SliderState"
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Class %s already exists", className)
		return nil
	}

	class := &wvmodels.Class{
		Class:       className,
		Description: "Individual slider state assessment - tracks if slider is fluid, sticky, stuck, or hijacked",
		Vectorizer:  "none",
		Properties: []*wvmodels.Property{
			{
				Name:        "student_id",
				DataType:    []string{"text"},
				Description: "Student identifier",
			},
			{
				Name:        "spectrum_id",
				DataType:    []string{"int"},
				Description: "ETP spectrum ID (1-17)",
			},
			{
				Name:        "spectrum_name",
				DataType:    []string{"text"},
				Description: "ETP spectrum name",
			},
			{
				Name:        "current_position",
				DataType:    []string{"number"},
				Description: "Current slider position (-2 to +2)",
			},
			{
				Name:        "resting_position",
				DataType:    []string{"number"},
				Description: "Natural resting position (comfort zone default)",
			},
			{
				Name:        "state",
				DataType:    []string{"text"},
				Description: "fluid, sticky, stuck, or hijacked",
			},
			{
				Name:        "range_of_motion",
				DataType:    []string{"number"},
				Description: "0-4 (full range is 4, meaning can access -2 to +2)",
			},
			{
				Name:        "trigger_points",
				DataType:    []string{"number[]"},
				Description: "Positions that cause involuntary slider movement",
			},
			{
				Name:        "can_access_neg2",
				DataType:    []string{"boolean"},
				Description: "Can reach -2 when appropriate",
			},
			{
				Name:        "can_access_neg1",
				DataType:    []string{"boolean"},
				Description: "Can reach -1 when appropriate",
			},
			{
				Name:        "can_access_center",
				DataType:    []string{"boolean"},
				Description: "Can reach 0 (center) when appropriate",
			},
			{
				Name:        "can_access_pos1",
				DataType:    []string{"boolean"},
				Description: "Can reach +1 when appropriate",
			},
			{
				Name:        "can_access_pos2",
				DataType:    []string{"boolean"},
				Description: "Can reach +2 when appropriate",
			},
			{
				Name:        "assessment_date",
				DataType:    []string{"text"},
				Description: "ISO date of assessment",
			},
		},
	}

	err = weaviateClient.Schema().ClassCreator().WithClass(class).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class %s: %w", className, err)
	}

	log.Printf("✅ Created Weaviate class: %s", className)
	return nil
}

// CreateSystemShiftClass creates the SystemShift collection
// Records observed multi-slider movements (e.g., defensive activation, deep focus)
func CreateSystemShiftClass(ctx context.Context) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	className := "SystemShift"
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Class %s already exists", className)
		return nil
	}

	class := &wvmodels.Class{
		Class:       className,
		Description: "Records system-level ETP shifts - deliberate multi-slider movements",
		Vectorizer:  "none",
		Properties: []*wvmodels.Property{
			{
				Name:        "student_id",
				DataType:    []string{"text"},
				Description: "Student identifier",
			},
			{
				Name:        "shift_name",
				DataType:    []string{"text"},
				Description: "Name of the shift pattern (e.g., defensive_activation, deep_focus)",
			},
			{
				Name:        "trigger",
				DataType:    []string{"text"},
				Description: "What initiated the shift",
			},
			{
				Name:        "from_state",
				DataType:    []string{"text"},
				Description: "JSON of slider positions before shift",
			},
			{
				Name:        "to_state",
				DataType:    []string{"text"},
				Description: "JSON of slider positions after shift",
			},
			{
				Name:        "was_deliberate",
				DataType:    []string{"boolean"},
				Description: "Whether the Pilot consciously initiated the shift",
			},
			{
				Name:        "recovery_time",
				DataType:    []string{"number"},
				Description: "Seconds to return to baseline",
			},
			{
				Name:        "success_rating",
				DataType:    []string{"number"},
				Description: "0-1 rating of how well the shift served its purpose",
			},
			{
				Name:        "timestamp",
				DataType:    []string{"text"},
				Description: "ISO timestamp of the shift",
			},
		},
	}

	err = weaviateClient.Schema().ClassCreator().WithClass(class).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class %s: %w", className, err)
	}

	log.Printf("✅ Created Weaviate class: %s", className)
	return nil
}

// CreateRangeOfMotionAssessmentClass creates the RangeOfMotionAssessment collection
// Tracks student's ability to access different parts of each spectrum
func CreateRangeOfMotionAssessmentClass(ctx context.Context) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	className := "RangeOfMotionAssessment"
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Class %s already exists", className)
		return nil
	}

	class := &wvmodels.Class{
		Class:       className,
		Description: "Assesses student ability to access all positions on each ETP spectrum",
		Vectorizer:  "none",
		Properties: []*wvmodels.Property{
			{
				Name:        "student_id",
				DataType:    []string{"text"},
				Description: "Student identifier",
			},
			{
				Name:        "spectrum_id",
				DataType:    []string{"int"},
				Description: "ETP spectrum ID (1-17)",
			},
			{
				Name:        "spectrum_name",
				DataType:    []string{"text"},
				Description: "ETP spectrum name",
			},
			{
				Name:        "exercise_name",
				DataType:    []string{"text"},
				Description: "Name of the range of motion exercise",
			},
			{
				Name:        "observed_positions",
				DataType:    []string{"number[]"},
				Description: "Array of positions observed during exercise",
			},
			{
				Name:        "range_score",
				DataType:    []string{"number"},
				Description: "0-1 score (1 = full range accessed)",
			},
			{
				Name:        "fluidity_score",
				DataType:    []string{"number"},
				Description: "0-1 score of how smoothly student moves between positions",
			},
			{
				Name:        "stuck_positions",
				DataType:    []string{"number[]"},
				Description: "Positions the student could not access",
			},
			{
				Name:        "assessment_date",
				DataType:    []string{"text"},
				Description: "ISO date of assessment",
			},
			{
				Name:        "assessor_notes",
				DataType:    []string{"text"},
				Description: "Qualitative notes from observer",
			},
		},
	}

	err = weaviateClient.Schema().ClassCreator().WithClass(class).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class %s: %w", className, err)
	}

	log.Printf("✅ Created Weaviate class: %s", className)
	return nil
}

// CreateOperatingModeLogClass creates the OperatingModeLog collection
// Tracks when students use Dialectical, Boundary, or Disengage modes
func CreateOperatingModeLogClass(ctx context.Context) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	className := "OperatingModeLog"
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Class %s already exists", className)
		return nil
	}

	class := &wvmodels.Class{
		Class:       className,
		Description: "Logs operating mode selections and outcomes (Dialectical/Boundary/Disengage)",
		Vectorizer:  "none",
		Properties: []*wvmodels.Property{
			{
				Name:        "student_id",
				DataType:    []string{"text"},
				Description: "Student identifier",
			},
			{
				Name:        "mode_selected",
				DataType:    []string{"text"},
				Description: "dialectical, boundary, or disengage",
			},
			{
				Name:        "situation_description",
				DataType:    []string{"text"},
				Description: "Description of the situation requiring mode selection",
			},
			{
				Name:        "circuit_breaker_triggered",
				DataType:    []string{"text"},
				Description: "Which circuit breaker triggered (if any): power_imbalance, dehumanization, bad_faith, cruelty_as_status",
			},
			{
				Name:        "mode_appropriate",
				DataType:    []string{"boolean"},
				Description: "Whether the mode selection was appropriate for the situation",
			},
			{
				Name:        "outcome_rating",
				DataType:    []string{"number"},
				Description: "0-1 rating of how well the mode served the student",
			},
			{
				Name:        "reflection_notes",
				DataType:    []string{"text"},
				Description: "Student reflection on the experience",
			},
			{
				Name:        "timestamp",
				DataType:    []string{"text"},
				Description: "ISO timestamp",
			},
		},
	}

	err = weaviateClient.Schema().ClassCreator().WithClass(class).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class %s: %w", className, err)
	}

	log.Printf("✅ Created Weaviate class: %s", className)
	return nil
}

// CreateStrategyFluencyAssessmentClass creates the StrategyFluencyAssessment collection
// Tracks student ability to use both Cooperation and Competition as tools
func CreateStrategyFluencyAssessmentClass(ctx context.Context) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	className := "StrategyFluencyAssessment"
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Class %s already exists", className)
		return nil
	}

	class := &wvmodels.Class{
		Class:       className,
		Description: "Assesses fluency in using Cooperation and Competition as strategic tools",
		Vectorizer:  "none",
		Properties: []*wvmodels.Property{
			{
				Name:        "student_id",
				DataType:    []string{"text"},
				Description: "Student identifier",
			},
			{
				Name:        "natural_bias",
				DataType:    []string{"text"},
				Description: "cooperation or competition - the student's default strategy",
			},
			{
				Name:        "cooperation_skill_level",
				DataType:    []string{"number"},
				Description: "0-1 skill in using cooperation tools",
			},
			{
				Name:        "competition_skill_level",
				DataType:    []string{"number"},
				Description: "0-1 skill in using competition tools",
			},
			{
				Name:        "switching_ability",
				DataType:    []string{"number"},
				Description: "0-1 ability to switch between strategies",
			},
			{
				Name:        "situation_reading",
				DataType:    []string{"number"},
				Description: "0-1 accuracy in diagnosing which strategy fits situation",
			},
			{
				Name:        "integration_skill",
				DataType:    []string{"number"},
				Description: "0-1 ability to design systems using both strategies",
			},
			{
				Name:        "fluency_level",
				DataType:    []string{"text"},
				Description: "identity_locked, learning, situational, or fluent",
			},
			{
				Name:        "assessment_date",
				DataType:    []string{"text"},
				Description: "ISO date of assessment",
			},
		},
	}

	err = weaviateClient.Schema().ClassCreator().WithClass(class).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class %s: %w", className, err)
	}

	log.Printf("✅ Created Weaviate class: %s", className)
	return nil
}

// CreateStickerBookProgressClass creates the StickerBookProgress collection
// Tracks primary school self-calibration curriculum progress
func CreateStickerBookProgressClass(ctx context.Context) error {
	if err := EnsureWeaviateClient(ctx); err != nil {
		return fmt.Errorf("failed to ensure weaviate client: %w", err)
	}

	className := "StickerBookProgress"
	exists, err := WeaviateCollectionExists(ctx, className)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Class %s already exists", className)
		return nil
	}

	class := &wvmodels.Class{
		Class:       className,
		Description: "Primary school sticker book progress for self-calibration curriculum",
		Vectorizer:  "none",
		Properties: []*wvmodels.Property{
			{
				Name:        "student_id",
				DataType:    []string{"text"},
				Description: "Student identifier",
			},
			{
				Name:        "year_group",
				DataType:    []string{"int"},
				Description: "Year group (3-6)",
			},
			{
				Name:        "skill_name",
				DataType:    []string{"text"},
				Description: "Name of the skill being developed",
			},
			{
				Name:        "sticker_name",
				DataType:    []string{"text"},
				Description: "Name of the sticker earned (e.g., Master of the Dial)",
			},
			{
				Name:        "sticker_earned",
				DataType:    []string{"boolean"},
				Description: "Whether the sticker has been earned",
			},
			{
				Name:        "exercise_completed",
				DataType:    []string{"text"},
				Description: "Description of exercise that earned the sticker",
			},
			{
				Name:        "evidence_notes",
				DataType:    []string{"text"},
				Description: "Teacher notes on evidence of skill demonstration",
			},
			{
				Name:        "earned_date",
				DataType:    []string{"text"},
				Description: "ISO date sticker was earned",
			},
		},
	}

	err = weaviateClient.Schema().ClassCreator().WithClass(class).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class %s: %w", className, err)
	}

	log.Printf("✅ Created Weaviate class: %s", className)
	return nil
}

// EnsureETPSchema creates all ETP-related Weaviate classes
func EnsureETPSchema(ctx context.Context) error {
	log.Println("Ensuring ETP Weaviate schema...")

	if err := CreateETPProfileClass(ctx); err != nil {
		return fmt.Errorf("failed to create ETPProfile class: %w", err)
	}

	if err := CreateSliderStateClass(ctx); err != nil {
		return fmt.Errorf("failed to create SliderState class: %w", err)
	}

	if err := CreateSystemShiftClass(ctx); err != nil {
		return fmt.Errorf("failed to create SystemShift class: %w", err)
	}

	if err := CreateRangeOfMotionAssessmentClass(ctx); err != nil {
		return fmt.Errorf("failed to create RangeOfMotionAssessment class: %w", err)
	}

	if err := CreateOperatingModeLogClass(ctx); err != nil {
		return fmt.Errorf("failed to create OperatingModeLog class: %w", err)
	}

	if err := CreateStrategyFluencyAssessmentClass(ctx); err != nil {
		return fmt.Errorf("failed to create StrategyFluencyAssessment class: %w", err)
	}

	if err := CreateStickerBookProgressClass(ctx); err != nil {
		return fmt.Errorf("failed to create StickerBookProgress class: %w", err)
	}

	log.Println("✅ All ETP Weaviate classes created successfully")
	return nil
}
