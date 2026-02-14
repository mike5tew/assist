package main

import (
	"context"
	"log"
	"net/url"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

const version = "1.0.0"

func main() {
	log.Printf("🎓 Student Weaviate Schema Setup v%s", version)
	log.Println("Setting up schema for student-facing learning intervention database...")

	weaviateURL := os.Getenv("STUDENT_WEAVIATE_URL")
	if weaviateURL == "" {
		weaviateURL = "http://localhost:8088"
	}
	log.Printf("Connecting to Weaviate at %s", weaviateURL)

	parsedURL, err := url.Parse(weaviateURL)
	if err != nil {
		log.Fatalf("Invalid Weaviate URL: %v", err)
	}

	cfg := weaviate.Config{
		Host:   parsedURL.Host,
		Scheme: parsedURL.Scheme,
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Weaviate client: %v", err)
	}

	_, err = client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		log.Fatalf("Weaviate is not ready: %v", err)
	}
	log.Println("✅ Successfully connected to Weaviate")

	classes := []struct {
		name  string
		setup func(*weaviate.Client)
	}{
		{"ETPProfile", setupETPProfileClass},
		{"LearningIntervention", setupLearningInterventionClass},
		{"ActionLog", setupActionLogClass},
		{"StudentProgress", setupStudentProgressClass},
		{"InterventionLibrary", setupInterventionLibraryClass},
		{"CHISGSkill", setupCHISGSkillClass},
		{"StudentBarrier", setupStudentBarrierClass},
		{"InterventionLever", setupInterventionLeverClass},
		{"TeachingPattern", setupTeachingPatternClass},
		// Integration layer collections
		{"FrictionMonitor", setupFrictionMonitorClass},
		{"InterventionSequence", setupInterventionSequenceClass},
		{"SkillTeachingPattern", setupSkillTeachingPatternClass},
	}

	for _, c := range classes {
		log.Printf("Setting up class: %s", c.name)
		c.setup(client)
	}

	verifyClasses(client, classes)
	log.Println("🎉 Student Weaviate schema setup completed successfully!")
}

func setupETPProfileClass(client *weaviate.Client) {
	className := "ETPProfile"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Student's Emotional Trigger Point profile across 9 core spectra + 2 global moderators",
		Properties: []*models.Property{
			{Name: "studentId", DataType: []string{"text"}, Description: "Unique student identifier"},
			{Name: "schoolId", DataType: []string{"text"}, Description: "School identifier for multi-tenancy"},
			{Name: "profileName", DataType: []string{"text"}, Description: "Human-readable profile name"},
			// 8 Core Spectra (biological, innate, orthogonal)
			{Name: "socialGravity", DataType: []string{"int"}, Description: "Independent (-100) to Cohesive (+100) - social interaction drains vs charges"},
			{Name: "energyDirectionality", DataType: []string{"int"}, Description: "Inward (-100) to Outward (+100) - external stimulation overloads vs energizes"},
			{Name: "voltageSensitivity", DataType: []string{"int"}, Description: "Insulated (-100) to Conductive (+100) - emotional current tolerance"},
			{Name: "threatResponse", DataType: []string{"int"}, Description: "Passive (-100) to Aggressive (+100) - avoidance vs confrontation voltage"},
			{Name: "careResponse", DataType: []string{"int"}, Description: "Detached (-100) to Nurturing (+100) - distance vs care voltage to vulnerability"},
			{Name: "riskTolerance", DataType: []string{"int"}, Description: "Averse (-100) to Seeking (+100) - risk creates anxiety vs excitement"},
			{Name: "integrityLogic", DataType: []string{"int"}, Description: "Relativistic (-100) to Absolutist (+100) - moral flexibility vs rigidity"},
			{Name: "mirrorNeuronTuning", DataType: []string{"int"}, Description: "Selective (-100) to Absorbent (+100) - others' emotions distinct vs shared"},
			{Name: "orderliness", DataType: []string{"int"}, Description: "Flexible (-100) to Ordered (+100) - preference for structure vs spontaneity"},
			// 2 Global Moderators (affect ALL spectra)
			{Name: "pilotStrength", DataType: []string{"int"}, Description: "Executive function capacity (0-100) - hand on all sliders"},
			{Name: "currentLoad", DataType: []string{"int"}, Description: "Stress/depletion level (0-100) - narrows range on all spectra"},
			// Metadata
			{Name: "createdAt", DataType: []string{"date"}, Description: "Profile creation timestamp"},
			{Name: "updatedAt", DataType: []string{"date"}, Description: "Last update timestamp"},
			{Name: "version", DataType: []string{"int"}, Description: "Profile version for history tracking"},
		},
		Vectorizer: "none",
	}
	createClass(client, class)
}

func setupLearningInterventionClass(client *weaviate.Client) {
	className := "LearningIntervention"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Planned learning interventions based on ETP profile analysis",
		Properties: []*models.Property{
			{Name: "studentId", DataType: []string{"text"}, Description: "Student identifier"},
			{Name: "schoolId", DataType: []string{"text"}, Description: "School identifier"},
			{Name: "etpProfileId", DataType: []string{"text"}, Description: "Reference to ETPProfile"},
			{Name: "interventionType", DataType: []string{"text"}, Description: "Type: skill_building, barrier_removal, enrichment"},
			{Name: "targetSkillIds", DataType: []string{"text[]"}, Description: "CHISG skill IDs being targeted"},
			{Name: "detectedBarrierIds", DataType: []string{"text[]"}, Description: "Barrier IDs being addressed"},
			{Name: "selectedLevers", DataType: []string{"text[]"}, Description: "Intervention lever IDs to apply"},
			{Name: "teachingPatternIds", DataType: []string{"text[]"}, Description: "Teaching pattern IDs to use"},
			{Name: "actions", DataType: []string{"text"}, Description: "JSON array of action steps"},
			{Name: "priority", DataType: []string{"int"}, Description: "Priority 1-10"},
			{Name: "status", DataType: []string{"text"}, Description: "Status: planned, active, completed, paused"},
			{Name: "estimatedDuration", DataType: []string{"text"}, Description: "Estimated time to complete"},
			{Name: "createdAt", DataType: []string{"date"}, Description: "Creation timestamp"},
			{Name: "scheduledFor", DataType: []string{"date"}, Description: "When to execute"},
		},
		Vectorizer: "none",
	}
	createClass(client, class)
}

func setupActionLogClass(client *weaviate.Client) {
	className := "ActionLog"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Log of actions taken during interventions with outcomes",
		Properties: []*models.Property{
			{Name: "studentId", DataType: []string{"text"}, Description: "Student identifier"},
			{Name: "interventionId", DataType: []string{"text"}, Description: "Reference to LearningIntervention"},
			{Name: "actionType", DataType: []string{"text"}, Description: "Type of action taken"},
			{Name: "actionDescription", DataType: []string{"text"}, Description: "Description of the action"},
			{Name: "skillsAddressed", DataType: []string{"text[]"}, Description: "CHISG skill IDs addressed"},
			{Name: "barriersEncountered", DataType: []string{"text[]"}, Description: "Barrier IDs encountered"},
			{Name: "leversApplied", DataType: []string{"text[]"}, Description: "Lever IDs that were applied"},
			{Name: "outcome", DataType: []string{"text"}, Description: "Outcome: success, partial, blocked, abandoned"},
			{Name: "confidenceScore", DataType: []string{"number"}, Description: "Confidence in outcome assessment 0-1"},
			{Name: "studentResponse", DataType: []string{"text"}, Description: "How student responded"},
			{Name: "teacherNotes", DataType: []string{"text"}, Description: "Additional teacher observations"},
			{Name: "etpShiftObserved", DataType: []string{"text"}, Description: "JSON of observed ETP changes"},
			{Name: "timestamp", DataType: []string{"date"}, Description: "When action occurred"},
			{Name: "duration", DataType: []string{"int"}, Description: "Duration in minutes"},
		},
		Vectorizer: "none",
	}
	createClass(client, class)
}

func setupStudentProgressClass(client *weaviate.Client) {
	className := "StudentProgress"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Aggregated progress metrics for students",
		Properties: []*models.Property{
			{Name: "studentId", DataType: []string{"text"}, Description: "Student identifier"},
			{Name: "schoolId", DataType: []string{"text"}, Description: "School identifier"},
			{Name: "skillMastery", DataType: []string{"text"}, Description: "JSON map of skill_id -> mastery_level (0-100)"},
			{Name: "barriersOvercome", DataType: []string{"text[]"}, Description: "Barrier IDs that have been resolved"},
			{Name: "activeBarriers", DataType: []string{"text[]"}, Description: "Barrier IDs still active"},
			{Name: "interventionCount", DataType: []string{"int"}, Description: "Total interventions attempted"},
			{Name: "successRate", DataType: []string{"number"}, Description: "Overall success rate 0-1"},
			{Name: "currentStreak", DataType: []string{"int"}, Description: "Current positive engagement streak"},
			{Name: "longestStreak", DataType: []string{"int"}, Description: "Longest engagement streak"},
			{Name: "lastActiveAt", DataType: []string{"date"}, Description: "Last activity timestamp"},
			{Name: "engagementTrend", DataType: []string{"text"}, Description: "Trend: improving, stable, declining"},
			{Name: "recommendedFocus", DataType: []string{"text[]"}, Description: "Skill IDs recommended for focus"},
		},
		Vectorizer: "none",
	}
	createClass(client, class)
}

func setupInterventionLibraryClass(client *weaviate.Client) {
	className := "InterventionLibrary"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Library of intervention templates organized by ETP profile patterns",
		Properties: []*models.Property{
			{Name: "templateId", DataType: []string{"text"}, Description: "Unique template identifier"},
			{Name: "name", DataType: []string{"text"}, Description: "Template name"},
			{Name: "description", DataType: []string{"text"}, Description: "What this intervention addresses"},
			{Name: "targetEtpPattern", DataType: []string{"text"}, Description: "JSON describing target ETP pattern"},
			{Name: "targetBarrierTypes", DataType: []string{"text[]"}, Description: "Barrier categories this addresses"},
			{Name: "requiredLevers", DataType: []string{"text[]"}, Description: "Lever IDs required"},
			{Name: "suggestedPatterns", DataType: []string{"text[]"}, Description: "Teaching pattern IDs to use"},
			{Name: "actionSteps", DataType: []string{"text"}, Description: "JSON array of action step templates"},
			{Name: "successIndicators", DataType: []string{"text[]"}, Description: "How to measure success"},
			{Name: "contraindications", DataType: []string{"text[]"}, Description: "When NOT to use this"},
			{Name: "difficultyLevel", DataType: []string{"text"}, Description: "Teacher skill required: novice, intermediate, expert"},
			{Name: "evidenceBase", DataType: []string{"text"}, Description: "Research/evidence supporting this"},
			{Name: "usageCount", DataType: []string{"int"}, Description: "How many times used"},
			{Name: "successRate", DataType: []string{"number"}, Description: "Aggregate success rate 0-1"},
		},
		Vectorizer: "text2vec-transformers",
	}
	createClass(client, class)
}

func setupCHISGSkillClass(client *weaviate.Client) {
	className := "CHISGSkill"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "CHISG (Coherent Human Intelligence Skills Graph) skills from humanOS",
		Properties: []*models.Property{
			{Name: "chisgId", DataType: []string{"text"}, Description: "Unique CHISG identifier (e.g., CHISG_1)"},
			{Name: "originalSkillId", DataType: []string{"int"}, Description: "Original skill ID from source"},
			{Name: "name", DataType: []string{"text"}, Description: "Skill name"},
			{Name: "description", DataType: []string{"text"}, Description: "Skill description"},
			{Name: "domain", DataType: []string{"text"}, Description: "Domain: FOCUS & TOOLS, INTERPERSONAL, etc."},
			{Name: "layer", DataType: []string{"text"}, Description: "Layer: Intervention, Foundation, etc."},
			{Name: "etpSpectrumId", DataType: []string{"text"}, Description: "Which ETP spectrum this skill modulates"},
			{Name: "etpVector", DataType: []string{"text"}, Description: "Direction: Toward Positive, Toward Negative, Neutral"},
			{Name: "etpMechanism", DataType: []string{"text"}, Description: "How it modulates the spectrum"},
			{Name: "etpEnergyCost", DataType: []string{"text"}, Description: "Energy cost: Low, Medium, High"},
			{Name: "enablesSkills", DataType: []string{"text[]"}, Description: "Skill names this enables"},
			{Name: "requiresSkills", DataType: []string{"text[]"}, Description: "Skill names required first"},
			{Name: "enablesSkillIds", DataType: []string{"text[]"}, Description: "CHISG IDs this enables"},
			{Name: "requiresSkillIds", DataType: []string{"text[]"}, Description: "CHISG IDs required first"},
		},
		Vectorizer: "text2vec-transformers",
	}
	createClass(client, class)
}

func setupStudentBarrierClass(client *weaviate.Client) {
	className := "StudentBarrier"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Barrier definitions for learning blockages from humanOS",
		Properties: []*models.Property{
			{Name: "barrierId", DataType: []string{"text"}, Description: "Unique barrier identifier"},
			{Name: "name", DataType: []string{"text"}, Description: "Barrier name"},
			{Name: "description", DataType: []string{"text"}, Description: "Barrier description"},
			{Name: "category", DataType: []string{"text"}, Description: "Category: acute, chronic, structural, enrichment"},
			{Name: "triggerPatterns", DataType: []string{"text[]"}, Description: "Regex patterns that indicate this barrier"},
			{Name: "behavioralSigns", DataType: []string{"text[]"}, Description: "Observable behaviors"},
			{Name: "avoidanceTactics", DataType: []string{"text[]"}, Description: "How students avoid (e.g., 'I don't know')"},
			{Name: "effectiveLevers", DataType: []string{"text[]"}, Description: "Lever IDs that work for this barrier"},
			{Name: "ineffectiveApproaches", DataType: []string{"text[]"}, Description: "What NOT to do"},
			{Name: "escalationPath", DataType: []string{"text"}, Description: "What to do if barrier persists"},
			{Name: "affectedEtpSpectra", DataType: []string{"text[]"}, Description: "ETP spectra this barrier impacts"},
			{Name: "typicalDuration", DataType: []string{"text"}, Description: "How long this typically lasts"},
			{Name: "severity", DataType: []string{"text"}, Description: "Severity: low, medium, high, critical"},
		},
		Vectorizer: "text2vec-transformers",
	}
	createClass(client, class)
}

func setupInterventionLeverClass(client *weaviate.Client) {
	className := "InterventionLever"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Intervention levers - specific actions that can shift student state",
		Properties: []*models.Property{
			{Name: "leverId", DataType: []string{"text"}, Description: "Unique lever identifier"},
			{Name: "name", DataType: []string{"text"}, Description: "Lever name"},
			{Name: "description", DataType: []string{"text"}, Description: "What this lever does"},
			{Name: "leverType", DataType: []string{"text"}, Description: "Type: voltage_reduction, reframing, bridging, etc."},
			{Name: "applicationMethod", DataType: []string{"text"}, Description: "How to apply this lever"},
			{Name: "examplePhrases", DataType: []string{"text[]"}, Description: "Example things to say"},
			{Name: "exampleActions", DataType: []string{"text[]"}, Description: "Example things to do"},
			{Name: "targetBarriers", DataType: []string{"text[]"}, Description: "Barrier IDs this lever addresses"},
			{Name: "targetEtpSpectra", DataType: []string{"text[]"}, Description: "ETP spectra this can shift"},
			{Name: "expectedShiftDirection", DataType: []string{"text"}, Description: "Which way it shifts the spectrum"},
			{Name: "prerequisites", DataType: []string{"text[]"}, Description: "What must be true first"},
			{Name: "contraindications", DataType: []string{"text[]"}, Description: "When NOT to use"},
			{Name: "energyCost", DataType: []string{"text"}, Description: "Energy cost: low, medium, high"},
			{Name: "timeRequired", DataType: []string{"text"}, Description: "Time to apply: seconds, minutes, session"},
		},
		Vectorizer: "text2vec-transformers",
	}
	createClass(client, class)
}

func setupTeachingPatternClass(client *weaviate.Client) {
	className := "TeachingPattern"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Teaching patterns from humanOS classroom strategies",
		Properties: []*models.Property{
			{Name: "patternId", DataType: []string{"text"}, Description: "Unique pattern identifier"},
			{Name: "name", DataType: []string{"text"}, Description: "Pattern name (e.g., 'Voltage Reduction Through Familiarity')"},
			{Name: "description", DataType: []string{"text"}, Description: "What this pattern achieves"},
			{Name: "patternType", DataType: []string{"text"}, Description: "Type: engagement, de-escalation, scaffolding, etc."},
			{Name: "steps", DataType: []string{"text"}, Description: "JSON array of pattern steps"},
			{Name: "variations", DataType: []string{"text"}, Description: "JSON array of variations"},
			{Name: "adaptations", DataType: []string{"text"}, Description: "How to adapt for different students"},
			{Name: "bestUsedWhen", DataType: []string{"text[]"}, Description: "Situations where this works best"},
			{Name: "avoidWhen", DataType: []string{"text[]"}, Description: "Situations to avoid this pattern"},
			{Name: "combinesWellWith", DataType: []string{"text[]"}, Description: "Other pattern IDs that work together"},
			{Name: "targetEtpProfiles", DataType: []string{"text"}, Description: "JSON describing ideal ETP profiles"},
			{Name: "targetBarrierCategories", DataType: []string{"text[]"}, Description: "Barrier categories this addresses"},
			{Name: "successIndicators", DataType: []string{"text[]"}, Description: "How to know it's working"},
			{Name: "sourceReference", DataType: []string{"text"}, Description: "Source documentation reference"},
		},
		Vectorizer: "text2vec-transformers",
	}
	createClass(client, class)
}

// ============================================================================
// INTEGRATION LAYER CLASSES
// ============================================================================

func setupFrictionMonitorClass(client *weaviate.Client) {
	className := "FrictionMonitor"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Real-time monitoring of student friction signals during skill practice",
		Properties: []*models.Property{
			{Name: "studentId", DataType: []string{"text"}, Description: "Student identifier"},
			{Name: "sessionId", DataType: []string{"text"}, Description: "Current session identifier"},
			{Name: "currentSkillId", DataType: []string{"text"}, Description: "CHISG skill being worked on"},
			{Name: "currentActivity", DataType: []string{"text"}, Description: "What the student is doing"},
			// Friction signals
			{Name: "hesitationPatterns", DataType: []string{"text"}, Description: "JSON array of hesitation events with timestamps"},
			{Name: "errorPatterns", DataType: []string{"text"}, Description: "JSON array of error types and frequencies"},
			{Name: "selfReportedState", DataType: []string{"text"}, Description: "JSON with anxiety, confidence, frustration levels"},
			{Name: "progressionSpeed", DataType: []string{"number"}, Description: "How fast they're moving through material 0-1"},
			// Barrier detection
			{Name: "barrierLikelihoods", DataType: []string{"text"}, Description: "JSON map of barrier_id -> probability (0-1)"},
			{Name: "dominantBarrier", DataType: []string{"text"}, Description: "Most likely barrier ID"},
			{Name: "barrierConfidence", DataType: []string{"number"}, Description: "Confidence in barrier detection 0-1"},
			// Intervention recommendation
			{Name: "recommendedLeverId", DataType: []string{"text"}, Description: "Which InterventionLever to pull now"},
			{Name: "recommendedPatternId", DataType: []string{"text"}, Description: "Which TeachingPattern to apply"},
			{Name: "urgency", DataType: []string{"text"}, Description: "Urgency level: low, medium, high, critical"},
			// State
			{Name: "voltageState", DataType: []string{"int"}, Description: "Current cognitive/emotional load 0-100"},
			{Name: "engagementLevel", DataType: []string{"int"}, Description: "Current engagement 0-100"},
			{Name: "timestamp", DataType: []string{"date"}, Description: "When this reading was taken"},
		},
		Vectorizer: "none",
	}
	createClass(client, class)
}

func setupInterventionSequenceClass(client *weaviate.Client) {
	className := "InterventionSequence"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Dynamic sequences of interventions that adapt to student response",
		Properties: []*models.Property{
			{Name: "studentId", DataType: []string{"text"}, Description: "Student identifier"},
			{Name: "skillId", DataType: []string{"text"}, Description: "Skill being worked on"},
			{Name: "sessionId", DataType: []string{"text"}, Description: "Session identifier"},
			// Current state
			{Name: "currentInterventionId", DataType: []string{"text"}, Description: "Active intervention ID"},
			{Name: "currentLeverId", DataType: []string{"text"}, Description: "Active lever being applied"},
			{Name: "currentPatternId", DataType: []string{"text"}, Description: "Active teaching pattern"},
			// Student response tracking
			{Name: "studentResponse", DataType: []string{"text"}, Description: "Response type: resistant, engaged, confused, progressing, stuck"},
			{Name: "responseIntensity", DataType: []string{"number"}, Description: "How strongly they responded 0-1"},
			{Name: "responseTimestamp", DataType: []string{"date"}, Description: "When response was observed"},
			// Next steps
			{Name: "nextInterventionOptions", DataType: []string{"text"}, Description: "JSON array of possible next interventions with scores"},
			{Name: "recommendedNextId", DataType: []string{"text"}, Description: "Recommended next intervention ID"},
			{Name: "escalationNeeded", DataType: []string{"boolean"}, Description: "Whether to escalate to different approach"},
			// History
			{Name: "sequenceHistory", DataType: []string{"text"}, Description: "JSON array of interventions tried with results"},
			{Name: "totalInterventions", DataType: []string{"int"}, Description: "Number of interventions in this sequence"},
			{Name: "successfulInterventions", DataType: []string{"int"}, Description: "Number that had positive response"},
			// Learning metrics
			{Name: "learningRate", DataType: []string{"number"}, Description: "How quickly student adapts to interventions 0-1"},
			{Name: "adaptationPattern", DataType: []string{"text"}, Description: "Pattern type: fast-responder, gradual, resistant, volatile"},
			{Name: "effectivenessScore", DataType: []string{"number"}, Description: "Overall sequence effectiveness 0-1"},
		},
		Vectorizer: "none",
	}
	createClass(client, class)
}

func setupSkillTeachingPatternClass(client *weaviate.Client) {
	className := "SkillTeachingPattern"
	if classExists(client, className) {
		log.Printf("  ⏭️  Class %s already exists", className)
		return
	}

	class := &models.Class{
		Class:       className,
		Description: "Maps which teaching patterns work best for which skills",
		Properties: []*models.Property{
			{Name: "skillId", DataType: []string{"text"}, Description: "CHISG skill ID"},
			{Name: "skillName", DataType: []string{"text"}, Description: "Skill name for readability"},
			{Name: "skillDomain", DataType: []string{"text"}, Description: "Skill domain"},
			// Optimal patterns
			{Name: "optimalPatterns", DataType: []string{"text"}, Description: "JSON array of {patternId, effectiveness, conditions}"},
			{Name: "primaryPatternId", DataType: []string{"text"}, Description: "Best pattern for most students"},
			{Name: "primaryPatternEffectiveness", DataType: []string{"number"}, Description: "Effectiveness score 0-1"},
			// ETP-specific mappings
			{Name: "patternETPAlignment", DataType: []string{"text"}, Description: "JSON map of patternId -> ETP profile ranges"},
			{Name: "etpSensitivePatterns", DataType: []string{"text[]"}, Description: "Patterns that vary significantly by ETP"},
			// Sequence recommendations
			{Name: "patternSequence", DataType: []string{"text[]"}, Description: "Optimal order to apply patterns"},
			{Name: "sequenceRationale", DataType: []string{"text"}, Description: "Why this sequence works"},
			// What to avoid
			{Name: "contraindicatedPatterns", DataType: []string{"text[]"}, Description: "Patterns that make this skill harder"},
			{Name: "contraindicationReasons", DataType: []string{"text"}, Description: "JSON map of patternId -> reason to avoid"},
			// Default levers
			{Name: "defaultLevers", DataType: []string{"text[]"}, Description: "Lever IDs to try first when stuck"},
			{Name: "expectedBarriers", DataType: []string{"text[]"}, Description: "Barrier IDs commonly seen with this skill"},
			// Learning data
			{Name: "totalAttempts", DataType: []string{"int"}, Description: "How many times this skill has been attempted"},
			{Name: "successRate", DataType: []string{"number"}, Description: "Overall success rate 0-1"},
			{Name: "lastUpdated", DataType: []string{"date"}, Description: "When effectiveness data was last updated"},
		},
		Vectorizer: "text2vec-transformers",
	}
	createClass(client, class)
}

func classExists(client *weaviate.Client, className string) bool {
	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		log.Printf("Warning: couldn't check schema: %v", err)
		return false
	}
	for _, class := range schema.Classes {
		if class.Class == className {
			return true
		}
	}
	return false
}

func createClass(client *weaviate.Client, class *models.Class) {
	err := client.Schema().ClassCreator().WithClass(class).Do(context.Background())
	if err != nil {
		log.Printf("  ❌ Error creating class %s: %v", class.Class, err)
		return
	}
	log.Printf("  ✅ Created class %s", class.Class)
}

func verifyClasses(client *weaviate.Client, classes []struct {
	name  string
	setup func(*weaviate.Client)
}) {
	log.Println("\nVerifying all classes...")
	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		log.Printf("❌ Error getting schema: %v", err)
		return
	}

	existingClasses := make(map[string]bool)
	for _, class := range schema.Classes {
		existingClasses[class.Class] = true
	}

	allPresent := true
	for _, c := range classes {
		if existingClasses[c.name] {
			log.Printf("  ✅ %s", c.name)
		} else {
			log.Printf("  ❌ %s NOT FOUND", c.name)
			allPresent = false
		}
	}

	if allPresent {
		log.Println("\n✅ All 12 classes verified successfully!")
	} else {
		log.Println("\n⚠️  Some classes are missing - check errors above")
	}
}
