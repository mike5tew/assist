package insight

import (
	"time"
)

// Tell represents a single observable behavior (e.g., "Micro-withdrawal", "Peer-check")
type Tell struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"` // e.g., "Engagement", "Social", "Cognitive"
	ScoreChange float64 `json:"score_change"`
}

// Observation is the record of a Tell being witnessed in the classroom
type Observation struct {
	ID            string    `json:"id"`
	StudentID     string    `json:"student_id"`
	TeacherID     string    `json:"teacher_id"`
	TellID        string    `json:"tell_id"`
	Timestamp     time.Time `json:"timestamp"`
	Context       string    `json:"context"`        // e.g., "Maths Lesson", "Transition"
	Justification string    `json:"justification"`  // Objective evidence for the observation (crucial for downward shifts)
	Voltage       float64   `json:"voltage"`        // The observed/set arousal level
	SocialGravity float64   `json:"social_gravity"` // Impact of peer group
}

// ProfileAdjustment tracks an explicit change to the student's sensitivity baseline
type ProfileAdjustment struct {
	ID            string    `json:"id"`
	StudentID     string    `json:"student_id"`
	SpectrumID    string    `json:"spectrum_id"` // One of the 15 markers
	OldValue      float64   `json:"old_value"`
	NewValue      float64   `json:"new_value"`
	Justification string    `json:"justification"` // Why was this calibrated downwards?
	Timestamp     time.Time `json:"timestamp"`
	TeacherID     string    `json:"teacher_id"`
}

// ETPProfile is the longitudinal map of a student's Emotional Trigger Points
type ETPProfile struct {
	StudentID     string             `json:"student_id"`
	TriggerPoints map[string]float64 `json:"trigger_points"` // e.g., "RiskTolerance": 0.4
	LastUpdated   time.Time          `json:"last_updated"`
}

// RecursiveMapping links Skills to the ETPs they master (The "Recursive Mastery" logic)
type RecursiveMapping struct {
	SkillID  string  `json:"skill_id"`
	ETPID    string  `json:"etp_id"`
	Strength float64 `json:"strength"` // How much this skill modulates this ETP
}
