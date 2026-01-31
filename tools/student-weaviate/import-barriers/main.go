package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

const version = "1.0.0"

type BarriersFile struct {
	Barriers []Barrier `json:"barriers"`
}

type Barrier struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`

	Manifestations struct {
		Observable []string `json:"observable"`
		Verbal     []string `json:"verbal"`
	} `json:"manifestations"`

	UnderlyingCauses struct {
		Primary       string   `json:"primary"`
		ETPActivation []string `json:"etpActivation"`
	} `json:"underlyingCauses"`

	InterventionStrategy map[string]interface{} `json:"interventionStrategy"`

	AICoachImplementation struct {
		DetectionSignals []string `json:"detectionSignals"`
	} `json:"aiCoachImplementation"`

	TimelineExpectation map[string]string `json:"timelineExpectation"`
}

func main() {
	log.Printf("🚧 Barriers & Levers Import Tool v%s", version)

	humanOSPath := os.Getenv("HUMANOS_PATH")
	if humanOSPath == "" {
		humanOSPath = "../../../humanOS"
	}

	barriersPath := filepath.Join(humanOSPath, "shared", "schemas", "barriers.json")
	if len(os.Args) > 1 {
		barriersPath = os.Args[1]
	}

	log.Printf("Loading barriers from: %s", barriersPath)

	barriersData, err := loadBarriersFile(barriersPath)
	if err != nil {
		log.Fatalf("Failed to load barriers: %v", err)
	}
	log.Printf("✅ Loaded %d barriers", len(barriersData.Barriers))

	client, err := connectToWeaviate()
	if err != nil {
		log.Fatalf("Failed to connect to Weaviate: %v", err)
	}
	log.Println("✅ Connected to Weaviate")

	log.Println("\n📋 Importing barriers...")
	barrierSuccess, barrierErrors := importBarriers(client, barriersData.Barriers)
	log.Printf("   ✅ %d barriers imported, ❌ %d errors", barrierSuccess, barrierErrors)

	log.Println("\n🔧 Extracting levers from intervention strategies...")
	levers := extractLeversFromBarriers(barriersData.Barriers)
	log.Printf("   Extracted %d unique levers", len(levers))

	leverSuccess, leverErrors := importLevers(client, levers)
	log.Printf("   ✅ %d levers imported, ❌ %d errors", leverSuccess, leverErrors)

	log.Println("\n🎉 Import complete!")
}

func loadBarriersFile(path string) (*BarriersFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var barriersFile BarriersFile
	if err := json.Unmarshal(data, &barriersFile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &barriersFile, nil
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

func importBarriers(client *weaviate.Client, barriers []Barrier) (int, int) {
	successCount := 0
	errorCount := 0

	for _, b := range barriers {
		behavioralSigns := append(b.Manifestations.Observable, b.Manifestations.Verbal...)
		triggerPatterns := b.AICoachImplementation.DetectionSignals
		severity := determineSeverity(b.Category)
		effectiveLevers := extractLeverIDsFromStrategy(b.InterventionStrategy, b.ID)

		data := map[string]interface{}{
			"barrierId":             b.ID,
			"name":                  b.Name,
			"description":           b.UnderlyingCauses.Primary,
			"category":              b.Category,
			"triggerPatterns":       triggerPatterns,
			"behavioralSigns":       behavioralSigns,
			"avoidanceTactics":      b.Manifestations.Verbal,
			"effectiveLevers":       effectiveLevers,
			"ineffectiveApproaches": []string{},
			"escalationPath":        "",
			"affectedEtpSpectra":    b.UnderlyingCauses.ETPActivation,
			"typicalDuration":       getTypicalDuration(b.TimelineExpectation),
			"severity":              severity,
		}

		err := insertObject(client, "StudentBarrier", data)
		if err != nil {
			log.Printf("   ❌ Error importing barrier %s: %v", b.ID, err)
			errorCount++
		} else {
			log.Printf("   ✅ Imported barrier: %s", b.Name)
			successCount++
		}
	}

	return successCount, errorCount
}

type Lever struct {
	ID             string
	Name           string
	Description    string
	LeverType      string
	TargetBarriers []string
}

func extractLeversFromBarriers(barriers []Barrier) []Lever {
	leverMap := make(map[string]Lever)

	for _, b := range barriers {
		for phaseKey, phaseValue := range b.InterventionStrategy {
			phase, ok := phaseValue.(map[string]interface{})
			if !ok {
				continue
			}

			name, _ := phase["name"].(string)
			if name == "" {
				continue
			}

			leverID := fmt.Sprintf("lever_%s_%s", b.ID, phaseKey)

			purpose, _ := phase["purpose"].(string)
			rationale, _ := phase["rationale"].(string)
			description := purpose
			if description == "" {
				description = rationale
			}

			leverType := determineLeverType(phaseKey, name)

			lever := Lever{
				ID:             leverID,
				Name:           name,
				Description:    description,
				LeverType:      leverType,
				TargetBarriers: []string{b.ID},
			}

			leverMap[leverID] = lever
		}
	}

	levers := make([]Lever, 0, len(leverMap))
	for _, lever := range leverMap {
		levers = append(levers, lever)
	}

	return levers
}

func importLevers(client *weaviate.Client, levers []Lever) (int, int) {
	successCount := 0
	errorCount := 0

	for _, lever := range levers {
		data := map[string]interface{}{
			"leverId":                lever.ID,
			"name":                   lever.Name,
			"description":            lever.Description,
			"leverType":              lever.LeverType,
			"applicationMethod":      "",
			"examplePhrases":         []string{},
			"exampleActions":         []string{},
			"targetBarriers":         lever.TargetBarriers,
			"targetEtpSpectra":       []string{},
			"expectedShiftDirection": "",
			"prerequisites":          []string{},
			"contraindications":      []string{},
			"energyCost":             "",
			"timeRequired":           "",
		}

		err := insertObject(client, "InterventionLever", data)
		if err != nil {
			log.Printf("   ❌ Error importing lever %s: %v", lever.ID, err)
			errorCount++
		} else {
			log.Printf("   ✅ Imported lever: %s", lever.Name)
			successCount++
		}
	}

	return successCount, errorCount
}

func insertObject(client *weaviate.Client, className string, data map[string]interface{}) error {
	_, err := client.Data().Creator().
		WithClassName(className).
		WithProperties(data).
		Do(context.Background())
	return err
}

func determineSeverity(category string) string {
	switch category {
	case "initiation_barrier":
		return "medium"
	case "chronic_barrier":
		return "high"
	case "structural_barrier":
		return "critical"
	case "enrichment_need":
		return "low"
	default:
		return "medium"
	}
}

func determineLeverType(phaseKey, name string) string {
	lowerPhase := strings.ToLower(phaseKey)
	lowerName := strings.ToLower(name)

	switch {
	case strings.Contains(lowerPhase, "victory") || strings.Contains(lowerName, "victory"):
		return "positive_reinforcement"
	case strings.Contains(lowerPhase, "avoidance") || strings.Contains(lowerName, "ban"):
		return "avoidance_prevention"
	case strings.Contains(lowerPhase, "relationship") || strings.Contains(lowerName, "chat"):
		return "relationship_building"
	case strings.Contains(lowerPhase, "awareness") || strings.Contains(lowerName, "pattern"):
		return "pattern_awareness"
	case strings.Contains(lowerPhase, "engagement") || strings.Contains(lowerName, "incentive"):
		return "extrinsic_motivation"
	case strings.Contains(lowerPhase, "proximity") || strings.Contains(lowerName, "shoulder"):
		return "proximity_support"
	case strings.Contains(lowerPhase, "success") || strings.Contains(lowerName, "success"):
		return "micro_success"
	default:
		return "general"
	}
}

func extractLeverIDsFromStrategy(strategy map[string]interface{}, barrierID string) []string {
	var leverIDs []string
	for phaseKey := range strategy {
		leverIDs = append(leverIDs, fmt.Sprintf("lever_%s_%s", barrierID, phaseKey))
	}
	return leverIDs
}

func getTypicalDuration(timeline map[string]string) string {
	if timeline == nil {
		return ""
	}
	for _, value := range timeline {
		if value != "" {
			return value
		}
	}
	return ""
}
