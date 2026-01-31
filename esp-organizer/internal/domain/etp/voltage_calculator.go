package etp

import (
	"math"
)

// VoltageCalculator calculates voltage changes based on ETP profile and activities
type VoltageCalculator struct{}

// ActivityType categorizes activities by their voltage impact pattern
type ActivityType string

const (
	ActivitySocial        ActivityType = "social"        // Group interaction
	ActivitySolitary      ActivityType = "solitary"      // Alone work
	ActivityRisky         ActivityType = "risky"         // Uncertainty involved
	ActivityStructured    ActivityType = "structured"    // Clear expectations
	ActivityEmotional     ActivityType = "emotional"     // Emotional content
	ActivityCognitive     ActivityType = "cognitive"     // Mental processing
	ActivityPhysical      ActivityType = "physical"      // Physical activity
	ActivityCreative      ActivityType = "creative"      // Open-ended creation
	ActivityCompetitive   ActivityType = "competitive"   // Status/ranking involved
	ActivityCollaborative ActivityType = "collaborative" // Shared work
)

// VoltageImpact describes how an activity affects voltage for an ETP setting
type VoltageImpact struct {
	SpectrumID     int          `json:"spectrum_id"`
	SpectrumName   string       `json:"spectrum_name"`
	Setting        float64      `json:"setting"`        // Student's setting on this spectrum
	ActivityType   ActivityType `json:"activity_type"`  // Type of activity
	VoltageChange  float64      `json:"voltage_change"` // Positive = voltage increase, Negative = decrease
	IsAligned      bool         `json:"is_aligned"`     // Does activity align with their preference?
	Recommendation string       `json:"recommendation"` // What to do about it
}

// spectrumActivityAlignment maps each spectrum to which activities align with negative/positive settings
var spectrumActivityAlignment = map[int]struct {
	negativeActivities []ActivityType // Activities that reduce voltage for negative settings
	positiveActivities []ActivityType // Activities that reduce voltage for positive settings
}{
	1:  {[]ActivityType{ActivitySolitary}, []ActivityType{ActivitySocial, ActivityCollaborative}},                   // social_gravity
	2:  {[]ActivityType{ActivityCognitive}, []ActivityType{ActivitySocial}},                                         // guilt_response
	3:  {[]ActivityType{ActivityCognitive, ActivitySolitary}, []ActivityType{ActivityEmotional, ActivitySocial}},    // emotional_transparency
	4:  {[]ActivityType{ActivitySolitary, ActivityCognitive}, []ActivityType{ActivitySocial, ActivityPhysical}},     // energy_directionality
	5:  {[]ActivityType{ActivityCognitive, ActivitySolitary}, []ActivityType{ActivityEmotional, ActivitySocial}},    // mirror_neuron_tuning
	6:  {[]ActivityType{ActivitySolitary, ActivityStructured}, []ActivityType{ActivityCollaborative}},               // resource_allocation
	7:  {[]ActivityType{ActivityStructured, ActivitySolitary}, []ActivityType{ActivityEmotional, ActivityCreative}}, // voltage_sensitivity
	8:  {[]ActivityType{ActivityPhysical, ActivityCreative}, []ActivityType{ActivityCognitive, ActivityStructured}}, // impulse_gap
	9:  {[]ActivityType{ActivityStructured, ActivitySolitary}, []ActivityType{ActivityPhysical, ActivitySocial}},    // self_righting_speed
	10: {[]ActivityType{ActivityStructured}, []ActivityType{ActivityRisky, ActivityCreative}},                       // risk_tolerance
	11: {[]ActivityType{ActivityStructured}, []ActivityType{ActivityCreative, ActivityRisky}},                       // anticipation_bias
	12: {[]ActivityType{ActivityCognitive, ActivitySolitary}, []ActivityType{ActivityPhysical, ActivityEmotional}},  // presence_sensitivity
	13: {[]ActivityType{ActivityStructured}, []ActivityType{ActivityCreative, ActivityCompetitive}},                 // agency_threshold
	14: {[]ActivityType{ActivityStructured}, []ActivityType{ActivityCreative, ActivityCompetitive}},                 // authority_response
	15: {[]ActivityType{ActivityStructured}, []ActivityType{ActivityCreative, ActivityRisky}},                       // ambiguity_tolerance
	16: {[]ActivityType{ActivitySolitary, ActivityCognitive}, []ActivityType{ActivityCompetitive, ActivitySocial}},  // status_sensitivity
	17: {[]ActivityType{ActivityStructured}, []ActivityType{ActivityCreative}},                                      // integrity_logic
}

// CalculateVoltageImpact calculates how an activity impacts voltage for a given ETP setting
func (vc *VoltageCalculator) CalculateVoltageImpact(spectrumID int, setting float64, activity ActivityType) *VoltageImpact {
	spectrum := GetSpectrumByID(spectrumID)
	if spectrum == nil {
		return nil
	}

	alignment, exists := spectrumActivityAlignment[spectrumID]
	if !exists {
		return nil
	}

	// Check if activity aligns with the setting
	isAligned := false
	alignmentStrength := 0.0

	if setting < 0 {
		// Negative setting - check if activity is in negative alignment list
		for _, act := range alignment.negativeActivities {
			if act == activity {
				isAligned = true
				alignmentStrength = 1.0
				break
			}
		}
		// Check if it's in the opposite list (misaligned)
		if !isAligned {
			for _, act := range alignment.positiveActivities {
				if act == activity {
					alignmentStrength = -1.0
					break
				}
			}
		}
	} else {
		// Positive setting - check if activity is in positive alignment list
		for _, act := range alignment.positiveActivities {
			if act == activity {
				isAligned = true
				alignmentStrength = 1.0
				break
			}
		}
		// Check if it's in the opposite list (misaligned)
		if !isAligned {
			for _, act := range alignment.negativeActivities {
				if act == activity {
					alignmentStrength = -1.0
					break
				}
			}
		}
	}

	// Calculate voltage change
	// Aligned activities reduce voltage, misaligned increase it
	// Stronger settings = stronger effect
	voltageChange := -alignmentStrength * math.Abs(setting) * 10

	// Build recommendation
	var recommendation string
	if isAligned {
		recommendation = "Activity aligns with voltage preference - continue"
	} else if alignmentStrength < 0 {
		solution := GetCompatibilitySolution(spectrumID)
		if solution != nil {
			recommendation = "Activity misaligned with voltage preference. Consider: " + solution.BridgeActivity
		} else {
			recommendation = "Activity misaligned with voltage preference - provide support"
		}
	} else {
		recommendation = "Neutral activity - monitor voltage"
	}

	return &VoltageImpact{
		SpectrumID:     spectrumID,
		SpectrumName:   spectrum.Name,
		Setting:        setting,
		ActivityType:   activity,
		VoltageChange:  voltageChange,
		IsAligned:      isAligned,
		Recommendation: recommendation,
	}
}

// ETPProfile represents a student's full ETP profile (17 spectra)
type ETPProfile struct {
	StudentID string          `json:"student_id"`
	Settings  map[int]float64 `json:"settings"` // spectrumID -> setting (-2 to +2)
	Labels    map[int]string  `json:"labels"`   // spectrumID -> label (e.g., "Moderately Independent")
}

// CalculateTotalVoltageImpact calculates the combined voltage impact of an activity across all spectra
func (vc *VoltageCalculator) CalculateTotalVoltageImpact(profile *ETPProfile, activities []ActivityType) float64 {
	totalVoltage := 0.0

	for spectrumID, setting := range profile.Settings {
		for _, activity := range activities {
			impact := vc.CalculateVoltageImpact(spectrumID, setting, activity)
			if impact != nil {
				totalVoltage += impact.VoltageChange
			}
		}
	}

	return totalVoltage
}

// GetLabel returns a human-readable label for a spectrum setting
func GetLabel(spectrumID int, setting float64) string {
	spectrum := GetSpectrumByID(spectrumID)
	if spectrum == nil {
		return "Unknown"
	}

	intensity := ""
	if math.Abs(setting) < 0.5 {
		intensity = "Balanced"
	} else if math.Abs(setting) < 1.0 {
		intensity = "Slightly"
	} else if math.Abs(setting) < 1.5 {
		intensity = "Moderately"
	} else {
		intensity = "Strongly"
	}

	if intensity == "Balanced" {
		return "Balanced"
	}

	if setting < 0 {
		return intensity + " " + spectrum.NegativeLabel
	}
	return intensity + " " + spectrum.PositiveLabel
}

// AnalyzeClassroomCompatibility analyzes voltage compatibility across a group of students
type ClassroomCompatibilityAnalysis struct {
	SpectrumID       int              `json:"spectrum_id"`
	SpectrumName     string           `json:"spectrum_name"`
	Spread           float64          `json:"spread"`             // Max - Min settings
	HighConflictRisk bool             `json:"high_conflict_risk"` // True if spread > 3
	Clusters         []SettingCluster `json:"clusters"`           // Groups of similar settings
	RecommendedSetup string           `json:"recommended_setup"`  // How to structure activity
}

// SettingCluster represents a group of students with similar settings
type SettingCluster struct {
	Label      string  `json:"label"` // e.g., "Independent group"
	MinSetting float64 `json:"min_setting"`
	MaxSetting float64 `json:"max_setting"`
	Count      int     `json:"count"` // Number of students
}

// AnalyzeClassroomForSpectrum analyzes classroom compatibility for a single spectrum
func (vc *VoltageCalculator) AnalyzeClassroomForSpectrum(spectrumID int, settings []float64) *ClassroomCompatibilityAnalysis {
	spectrum := GetSpectrumByID(spectrumID)
	if spectrum == nil || len(settings) == 0 {
		return nil
	}

	// Find min/max
	minSetting := settings[0]
	maxSetting := settings[0]
	for _, s := range settings {
		if s < minSetting {
			minSetting = s
		}
		if s > maxSetting {
			maxSetting = s
		}
	}
	spread := maxSetting - minSetting

	// Cluster students
	negativeCount := 0
	neutralCount := 0
	positiveCount := 0

	for _, s := range settings {
		if s < -0.5 {
			negativeCount++
		} else if s > 0.5 {
			positiveCount++
		} else {
			neutralCount++
		}
	}

	clusters := []SettingCluster{}
	if negativeCount > 0 {
		clusters = append(clusters, SettingCluster{
			Label:      spectrum.NegativeLabel + " group",
			MinSetting: -2,
			MaxSetting: -0.5,
			Count:      negativeCount,
		})
	}
	if neutralCount > 0 {
		clusters = append(clusters, SettingCluster{
			Label:      "Balanced group",
			MinSetting: -0.5,
			MaxSetting: 0.5,
			Count:      neutralCount,
		})
	}
	if positiveCount > 0 {
		clusters = append(clusters, SettingCluster{
			Label:      spectrum.PositiveLabel + " group",
			MinSetting: 0.5,
			MaxSetting: 2,
			Count:      positiveCount,
		})
	}

	// Determine recommended setup
	var recommendedSetup string
	if spread < 2 {
		recommendedSetup = "Low voltage spread - whole class activities work well"
	} else if len(clusters) == 2 {
		solution := GetCompatibilitySolution(spectrumID)
		if solution != nil {
			recommendedSetup = "Two distinct groups - use " + solution.Name + ". Bridge activity: " + solution.BridgeActivity
		}
	} else {
		recommendedSetup = "High spread - consider choice-based activities with different voltage zones"
	}

	return &ClassroomCompatibilityAnalysis{
		SpectrumID:       spectrumID,
		SpectrumName:     spectrum.Name,
		Spread:           spread,
		HighConflictRisk: spread > 3,
		Clusters:         clusters,
		RecommendedSetup: recommendedSetup,
	}
}
