package constants

// AcademicLevel represents the educational level/tier of content
type AcademicLevel string

const (
	LevelGCSE          AcademicLevel = "gcse"
	LevelALevel        AcademicLevel = "a-level"
	LevelUndergraduate AcademicLevel = "undergraduate"
	LevelPostgraduate  AcademicLevel = "postgraduate"
	LevelProfessional  AcademicLevel = "professional"
)

// AllAcademicLevels returns all valid academic levels
func AllAcademicLevels() []AcademicLevel {
	return []AcademicLevel{
		LevelGCSE,
		LevelALevel,
		LevelUndergraduate,
		LevelPostgraduate,
		LevelProfessional,
	}
}

// IsValidAcademicLevel checks if the provided level is valid
func IsValidAcademicLevel(level string) bool {
	for _, valid := range AllAcademicLevels() {
		if string(valid) == level {
			return true
		}
	}
	return false
}

// Difficulty represents the relative difficulty of content
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

// AllDifficulties returns all valid difficulty levels
func AllDifficulties() []Difficulty {
	return []Difficulty{
		DifficultyEasy,
		DifficultyMedium,
		DifficultyHard,
	}
}

// IsValidDifficulty checks if the provided difficulty is valid
func IsValidDifficulty(d string) bool {
	for _, valid := range AllDifficulties() {
		if string(valid) == d {
			return true
		}
	}
	return false
}

// EstimatedTimeMinutes maps description to minutes for study planning
var EstimatedTimeMinutes = map[string]int{
	"quick_read":   5,   // Quick definition, formula
	"quick_review": 10,  // Short concept explanation
	"standard":     20,  // Typical study material
	"deep_dive":    45,  // Case study, complex topic
	"full_chapter": 90,  // Complete chapter
	"exam_prep":    120, // Comprehensive exam preparation
}
