package api

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/domain/etp"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ETPProfileRequest matches the JSON from the frontend
type ETPProfileRequest struct {
	Values map[string]int `json:"values"`
}

// ETPProfileDoc is the MongoDB document for an ETP profile
type ETPProfileDoc struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"userId" json:"userId"`
	SpectrumProfile map[string]int     `bson:"spectrumProfile" json:"spectrumProfile"`
	SpectrumVector  []float32          `bson:"spectrumVector" json:"spectrumVector"`
	Timestamp       time.Time          `bson:"timestamp" json:"timestamp"`
}

// ---------- MongoDB connection for ETP data ----------

var etpDB *mongo.Database

func getETPDB() (*mongo.Database, error) {
	if etpDB != nil {
		return etpDB, nil
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	etpDB = client.Database("esp_organizer")

	// Ensure indexes on etp_profiles collection
	coll := etpDB.Collection("etp_profiles")
	coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "userId", Value: 1}, {Key: "timestamp", Value: -1}},
	})

	log.Println("✅ ETP MongoDB connection established (database: esp_organizer)")
	return etpDB, nil
}

func ActionPlanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	var req ETPProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received ETP Profile: %+v", req.Values)

	// 1. Store ETP profile in MongoDB
	// Build the spectrum vector using canonical order
	order := []string{
		"social_gravity", "guilt_response", "emotional_transparency", "energy_directionality", "mirror_neuron_tuning", "resource_allocation",
		"voltage_sensitivity", "impulse_gap", "self_righting_speed", "risk_tolerance", "anticipation_bias", "presence_sensitivity",
		"agency_threshold", "authority_response", "ambiguity_tolerance", "status_sensitivity", "integrity_logic",
		"orderliness",
	}

	vector := make([]float32, len(order))
	for i, key := range order {
		if val, ok := req.Values[key]; ok {
			vector[i] = float32(val)
		} else {
			vector[i] = 0
		}
	}

	profileDoc := ETPProfileDoc{
		UserID:          "test-user-001", // TODO: Get from auth
		SpectrumProfile: req.Values,
		SpectrumVector:  vector,
		Timestamp:       time.Now(),
	}

	var profileId string

	db, err := getETPDB()
	if err != nil {
		log.Printf("[ActionPlan] MongoDB not available, continuing without persistence: %v", err)
		profileId = "ephemeral"
	} else {
		result, err := db.Collection("etp_profiles").InsertOne(r.Context(), profileDoc)
		if err != nil {
			log.Printf("[ActionPlan] Failed to save ETP profile to MongoDB: %v", err)
			profileId = "save-failed"
		} else {
			profileId = result.InsertedID.(primitive.ObjectID).Hex()
			log.Printf("[ActionPlan] Saved ETP profile to MongoDB: %s", profileId)
		}
	}

	// 2. Generate personalised action plan from the in-memory response matrix
	profileValues := make(map[string]float64)
	for k, v := range req.Values {
		profileValues[k] = float64(v)
	}

	actions := etp.GeneratePersonalisedPlan(profileValues, 0.5)

	// Build the action plan response
	actionPlan := buildActionPlanResponse(actions)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "success",
		"profileId":       profileId,
		"actionPlan":      actionPlan.Summary,
		"detailedActions": actionPlan.DetailedActions,
		"barrierProfile":  actionPlan.BarrierProfile,
		"languageGuide":   actionPlan.LanguageGuide,
		"playStrategies":  actionPlan.PlayStrategies,
		"unifiedLanguage": etp.UnifiedLanguageShifts,
		"playFirstSteps":  etp.PlayFirstPrinciple,
	})
}

// ActionPlanResponse holds the structured plan returned to the frontend
type ActionPlanResponse struct {
	Summary         []string                 `json:"summary"`
	DetailedActions []etp.PersonalisedAction `json:"detailed_actions"`
	BarrierProfile  map[string][]string      `json:"barrier_profile"`
	LanguageGuide   map[string]LanguageEntry `json:"language_guide"`
	PlayStrategies  []string                 `json:"play_strategies"`
}

// LanguageEntry pairs avoid/use language for a spectrum
type LanguageEntry struct {
	SpectrumName string   `json:"spectrum_name"`
	Setting      string   `json:"setting"`
	Avoid        []string `json:"avoid"`
	Use          []string `json:"use"`
}

// buildActionPlanResponse creates the action plan from in-memory response matrix data
func buildActionPlanResponse(actions []etp.PersonalisedAction) ActionPlanResponse {
	resp := ActionPlanResponse{
		BarrierProfile: make(map[string][]string),
		LanguageGuide:  make(map[string]LanguageEntry),
	}

	playSet := make(map[string]bool)

	for _, a := range actions {
		spectrum := etp.GetSpectrumByName(a.SpectrumName)
		spectrumLabel := a.SpectrumName
		if spectrum != nil {
			spectrumLabel = spectrum.SolutionName
		}

		summaryLine := fmt.Sprintf("[%s] %s setting (%s): %s",
			strings.ToUpper(a.BarrierType), a.Setting, a.SettingStrength, a.BarrierDesc)
		resp.Summary = append(resp.Summary, summaryLine)

		resp.BarrierProfile[a.BarrierType] = append(resp.BarrierProfile[a.BarrierType], a.BarrierDesc)

		resp.LanguageGuide[a.SpectrumName] = LanguageEntry{
			SpectrumName: spectrumLabel,
			Setting:      a.Setting,
			Avoid:        a.AvoidSaying,
			Use:          a.InsteadSay,
		}

		for _, p := range a.PlayStrategies {
			if !playSet[p] {
				playSet[p] = true
				resp.PlayStrategies = append(resp.PlayStrategies, p)
			}
		}
	}

	resp.DetailedActions = actions

	if len(resp.Summary) == 0 {
		resp.Summary = []string{
			"Your ETP profile is relatively balanced across all spectra.",
			"Focus on the Play-First Principle: observe settings, name them neutrally, design play that expands range.",
			"Use engineering language: 'Your setting is...' instead of 'You are...'",
		}
	}

	return resp
}
