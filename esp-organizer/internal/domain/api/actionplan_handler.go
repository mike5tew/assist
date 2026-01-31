package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

// ETPProfileRequest matches the JSON from the frontend
type ETPProfileRequest struct {
	Values map[string]int `json:"values"`
}

var assistClient *weaviate.Client

func getAssistClient() *weaviate.Client {
	if assistClient == nil {
		// Use environment variable or default to Docker service name
		weaviateHost := os.Getenv("WEAVIATE_URL")
		if weaviateHost == "" {
			weaviateHost = "weaviate:8080"
		} else {
			// Strip scheme if present (WEAVIATE_URL might include http://)
			if len(weaviateHost) > 7 && weaviateHost[:7] == "http://" {
				weaviateHost = weaviateHost[7:]
			} else if len(weaviateHost) > 8 && weaviateHost[:8] == "https://" {
				weaviateHost = weaviateHost[8:]
			}
		}
		log.Printf("Connecting to Weaviate at: http://%s", weaviateHost)

		cfg := weaviate.Config{
			Host:    weaviateHost,
			Scheme:  "http",
			Headers: nil,
		}
		assistClient = weaviate.New(cfg)
	}
	return assistClient
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

	client := getAssistClient()

	// 1. Create ETPProfile object
	// Convert map to separate fields (if needed) or just store as JSON string or vector
	// The schema has "spectrumVector" as number[]
	// I need to order the values correctly to make a vector.
	// I'll define a canonical order.
	order := []string{
		"social_gravity", "guilt_response", "emotional_transparency", "energy_directionality", "mirror_neuron_tuning", "resource_allocation",
		"voltage_sensitivity", "impulse_gap", "self_righting_speed", "risk_tolerance", "anticipation_bias", "presence_sensitivity",
		"agency_threshold", "authority_response", "ambiguity_tolerance", "status_sensitivity", "integrity_logic",
	}

	vector := make([]float32, len(order))
	for i, key := range order {
		if val, ok := req.Values[key]; ok {
			vector[i] = float32(val)
		} else {
			vector[i] = 0 // Default to 0
		}
	}

	// Create the object
	dataObj := map[string]interface{}{
		"userId":          "test-user-001",               // TODO: Get from auth
		"spectrumProfile": fmt.Sprintf("%v", req.Values), // Simple string rep
		"spectrumVector":  vector,                        // The raw vector
		"timestamp":       time.Now().Format(time.RFC3339),
	}

	response, err := client.Data().Creator().
		WithClassName("ETPProfile").
		WithProperties(dataObj).
		Do(context.Background())

	if err != nil {
		log.Printf("Error creating ETPProfile in Weaviate: %v", err)
		http.Error(w, fmt.Sprintf("Failed to save profile: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Created ETPProfile with ID: %v", response.Object.ID)

	// 2. Mock Action Plan Generation
	// (In future: Query PersonalityLearningStrategy using the vector)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"profileId": response.Object.ID,
		"actionPlan": []string{
			"Review your Social Gravity settings.",
			"Practice Impulse Gap expansion in low-stress environments.",
			"Monitor Risk Tolerance thresholds.",
		},
	})
}
