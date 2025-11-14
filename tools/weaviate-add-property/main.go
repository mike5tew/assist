package main

import (
	"context"
	"flag"
	"log"

	"esp-organizer/internal/InfoFlow/InfoStore/db"

	wvmodels "github.com/weaviate/weaviate/entities/models"
)

func main() {
	className := flag.String("class", "", "Weaviate class name")
	propertyName := flag.String("property", "", "Property name to add")
	flag.Parse()

	if *className == "" || *propertyName == "" {
		log.Fatal("Usage: weaviate-add-property -class <name> -property <name>")
	}

	ctx := context.Background()
	if err := db.InitializeWeaviateFromEnv(); err != nil {
		log.Fatalf("Failed to initialize Weaviate: %v", err)
	}

	client := db.GetWeaviateClient()

	// Define new property
	newProperty := &wvmodels.Property{
		Name:            *propertyName,
		DataType:        []string{"text"},
		IndexFilterable: boolPtr(true),
	}

	// Add property to class
	err := client.Schema().PropertyCreator().
		WithClassName(*className).
		WithProperty(newProperty).
		Do(ctx)

	if err != nil {
		log.Fatalf("Failed to add property: %v", err)
	}

	log.Printf("✅ Successfully added property '%s' to class '%s'", *propertyName, *className)
}

func boolPtr(b bool) *bool {
	return &b
}
