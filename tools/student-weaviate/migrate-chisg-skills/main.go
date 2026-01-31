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

type CHISGElement struct {
	CHISGID         string `json:"chisg_id"`
	OriginalSkillID int    `json:"original_skill_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Domain          string `json:"domain"`
	Layer           string `json:"layer"`
	ETPModulation   struct {
		SpectrumID string `json:"spectrum_id"`
		Vector     string `json:"vector"`
		Mechanism  string `json:"mechanism"`
		EnergyCost string `json:"energy_cost"`
	} `json:"etp_modulation"`
	Relationships struct {
		Enables  []string `json:"enables"`
		Requires []string `json:"requires"`
	} `json:"relationships"`
}

func main() {
	log.Printf("🧠 CHISG Skills Migration Tool v%s", version)

	humanOSPath := os.Getenv("HUMANOS_PATH")
	if humanOSPath == "" {
		humanOSPath = "../../../humanOS"
	}

	chisgPath := filepath.Join(humanOSPath, "data", "chisg_elements.json")
	if len(os.Args) > 1 {
		chisgPath = os.Args[1]
	}

	log.Printf("Loading CHISG elements from: %s", chisgPath)

	elements, err := loadCHISGElements(chisgPath)
	if err != nil {
		log.Fatalf("Failed to load CHISG elements: %v", err)
	}
	log.Printf("✅ Loaded %d CHISG elements", len(elements))

	nameToID := buildNameToIDMap(elements)
	log.Printf("✅ Built name->ID lookup with %d entries", len(nameToID))

	client, err := connectToWeaviate()
	if err != nil {
		log.Fatalf("Failed to connect to Weaviate: %v", err)
	}
	log.Println("✅ Connected to Weaviate")

	batchSize := 100
	successCount := 0
	errorCount := 0

	for i := 0; i < len(elements); i += batchSize {
		end := i + batchSize
		if end > len(elements) {
			end = len(elements)
		}

		batch := elements[i:end]
		for _, elem := range batch {
			enablesIDs := resolveNamesToIDs(elem.Relationships.Enables, nameToID)
			requiresIDs := resolveNamesToIDs(elem.Relationships.Requires, nameToID)

			data := map[string]interface{}{
				"chisgId":          elem.CHISGID,
				"originalSkillId":  elem.OriginalSkillID,
				"name":             elem.Name,
				"description":      elem.Description,
				"domain":           elem.Domain,
				"layer":            elem.Layer,
				"etpSpectrumId":    elem.ETPModulation.SpectrumID,
				"etpVector":        elem.ETPModulation.Vector,
				"etpMechanism":     elem.ETPModulation.Mechanism,
				"etpEnergyCost":    elem.ETPModulation.EnergyCost,
				"enablesSkills":    elem.Relationships.Enables,
				"requiresSkills":   elem.Relationships.Requires,
				"enablesSkillIds":  enablesIDs,
				"requiresSkillIds": requiresIDs,
			}

			err := insertObject(client, "CHISGSkill", data)
			if err != nil {
				errorCount++
			} else {
				successCount++
			}
		}
		log.Printf("  Batch %d-%d: processed", i+1, end)
	}

	log.Printf("\n🎉 Migration complete!")
	log.Printf("   ✅ Successfully inserted: %d", successCount)
	log.Printf("   ❌ Errors: %d", errorCount)
}

func loadCHISGElements(path string) ([]CHISGElement, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var elements []CHISGElement
	if err := json.Unmarshal(data, &elements); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return elements, nil
}

func buildNameToIDMap(elements []CHISGElement) map[string]string {
	lookup := make(map[string]string)
	for _, elem := range elements {
		normalizedName := strings.ToLower(strings.TrimSpace(elem.Name))
		lookup[normalizedName] = elem.CHISGID
	}
	return lookup
}

func resolveNamesToIDs(names []string, lookup map[string]string) []string {
	ids := make([]string, 0, len(names))
	for _, name := range names {
		normalizedName := strings.ToLower(strings.TrimSpace(name))
		if id, ok := lookup[normalizedName]; ok {
			ids = append(ids, id)
		}
	}
	return ids
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

func insertObject(client *weaviate.Client, className string, data map[string]interface{}) error {
	_, err := client.Data().Creator().
		WithClassName(className).
		WithProperties(data).
		Do(context.Background())
	return err
}
