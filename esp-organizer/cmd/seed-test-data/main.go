package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	fmt.Println("🌱 Seeding test data for HSG search...")

	// Get Weaviate client
	client := getWeaviateClient()
	if client == nil {
		log.Fatal("Failed to initialize Weaviate client")
	}

	// Create test summary chunks for XLA
	createTestSummaryChunks(client)

	fmt.Println("✅ Test data seeding complete")
}

func getWeaviateClient() *weaviate.Client {
	// Initialize from same code as in hsg_query_service.go
	// ...implementation here...
	return nil // Replace with actual implementation
}

func createTestSummaryChunks(client *weaviate.Client) {
	ctx := context.Background()

	// Sample data about X-linked agammaglobulinemia (XLA)
	summaryData := []struct {
		Content string
		Level   int
	}{
		{
			Content: "X-linked agammaglobulinemia (XLA) is a primary immunodeficiency disorder characterized by a lack of mature B cells and severe antibody deficiency. It is caused by mutations in the BTK gene which is located on the X chromosome, resulting in defective B-cell development and maturation.",
			Level:   1,
		},
		{
			Content: "Clinical manifestations of XLA typically begin after 6 months of age when maternal antibodies are depleted. Patients present with recurrent bacterial infections, particularly of the respiratory tract, including sinusitis, pneumonia, and otitis media. Less commonly, they may develop gastrointestinal infections, conjunctivitis, and skin infections.",
			Level:   2,
		},
		{
			Content: "Diagnosis of XLA involves laboratory findings such as very low or absent serum immunoglobulins (IgG, IgA, IgM), absent or markedly reduced B cells (<2% of lymphocytes), and normal T-cell numbers and function. Genetic testing to identify BTK gene mutations confirms the diagnosis.",
			Level:   2,
		},
		{
			Content: "Treatment for XLA consists primarily of immunoglobulin replacement therapy (IgRT), either intravenously (IVIG) or subcutaneously (SCIG), to provide passive immunity. Prompt antibiotic therapy for infections is crucial, and some patients may require prophylactic antibiotics. Live vaccines should be avoided.",
			Level:   2,
		},
		{
			Content: "The prognosis for patients with XLA has significantly improved with early diagnosis and appropriate management. Most patients lead relatively normal lives with regular immunoglobulin replacement therapy, though they remain at increased risk for chronic lung disease and certain malignancies.",
			Level:   1,
		},
	}

	// Create objects in Weaviate
	for i, data := range summaryData {
		props := map[string]interface{}{
			"content":  data.Content,
			"domain":   "immunology",
			"chunkId":  fmt.Sprintf("xla-summary-%d", i+1),
			"parentId": "root",
			"level":    data.Level,
		}

		// Create object
		id, err := client.Data().Creator().
			WithClassName("SummaryChunk").
			WithProperties(props).
			WithID(fmt.Sprintf("SummaryChunk-%s", primitive.NewObjectID().Hex())).
			Do(ctx)

		if err != nil {
			log.Printf("❌ Error creating test data object: %v", err)
		} else {
			log.Printf("✅ Created test summary chunk with ID: %s", id)
		}

		// Add a small delay between operations
		time.Sleep(500 * time.Millisecond)
	}

	// 3. Add the SummaryChunk to Weaviate
	weaviateProps := map[string]interface{}{
		"chunkId": "xla-summary-1",
		"domain":  "immunology",
		"content": "XLA is a genetic disorder affecting B-cells due to BTK gene mutations.",
	}
	_, err := client.Data().Creator().
		WithClassName("SummaryChunk").
		WithID(fmt.Sprintf("SummaryChunk-%s", primitive.NewObjectID().Hex())).
		WithProperties(weaviateProps).
		Do(context.Background())
	if err != nil {
		log.Fatalf("❌ Failed to create SummaryChunk object in Weaviate: %v", err)
	}
	log.Printf("✅ Added SummaryChunk to Weaviate with ID: %s", weaviateProps["chunkId"])
}
