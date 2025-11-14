package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/check-json/main.go path/to/skills.json")
		os.Exit(1)
	}

	filePath := os.Args[1]
	fmt.Printf("Checking JSON structure of: %s\n", filePath)

	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Parse the JSON
	var jsonData []map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	fmt.Printf("File contains %d skill objects\n", len(jsonData))

	// Check for parent_skill_ids
	skillsWithParents := 0

	for i, skill := range jsonData {
		name, _ := skill["name"].(string)

		// Check different possible formats of parent skill IDs
		parentIDs, hasParentIDs := skill["parent_skill_ids"]
		parentSkillIDs, hasParentSkillIDs := skill["parentSkillIDs"]

		if hasParentIDs || hasParentSkillIDs {
			skillsWithParents++

			// Print the first few with parents for inspection
			if skillsWithParents <= 5 {
				fmt.Printf("\nSkill #%d: %s\n", i+1, name)

				if hasParentIDs {
					fmt.Printf("  parent_skill_ids: %v\n", parentIDs)
				}

				if hasParentSkillIDs {
					fmt.Printf("  parentSkillIDs: %v\n", parentSkillIDs)
				}
			}
		}
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Total skills: %d\n", len(jsonData))
	fmt.Printf("- Skills with parent relationships: %d (%.1f%%)\n",
		skillsWithParents, float64(skillsWithParents)/float64(len(jsonData))*100)

	if skillsWithParents == 0 {
		fmt.Printf("\n❌ NO PARENT RELATIONSHIPS FOUND!\n")
		fmt.Printf("This is likely why no semantic links are being created.\n")
		fmt.Printf("Ensure your skills.json contains properly formatted parent_skill_ids arrays.\n")
	} else {
		fmt.Printf("\n✅ Found %d skills with parent relationships\n", skillsWithParents)
	}
}
