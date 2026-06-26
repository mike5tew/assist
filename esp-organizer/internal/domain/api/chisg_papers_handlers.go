package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"

	dbstore "esp-organizer/internal/store/db"
)

// ── /api/chisg/papers ────────────────────────────────────────────────────────

// CHISGPaperItem is the response shape for each paper.
type CHISGPaperItem struct {
	PaperID     string `json:"paper_id"`
	Title       string `json:"title"`
	Filename    string `json:"filename"`
	Pages       int    `json:"pages"`
	TotalChars  int    `json:"total_chars"`
	LinkCount   int    `json:"link_count"`
	EntityCount int    `json:"entity_count"`
}

// CHISGListPapersHandler handles GET /api/chisg/papers.
// Returns all AcademicPaper objects from Weaviate.
func CHISGListPapersHandler(w http.ResponseWriter, r *http.Request) {
	wc := dbstore.GetWeaviateClient()
	if wc == nil {
		if err := dbstore.InitializeWeaviateFromEnv(); err != nil {
			http.Error(w, "weaviate unavailable", http.StatusServiceUnavailable)
			return
		}
		wc = dbstore.GetWeaviateClient()
	}

	fields := []graphql.Field{
		{Name: "paper_id"},
		{Name: "title"},
		{Name: "filename"},
		{Name: "pages"},
		{Name: "total_chars"},
		{Name: "link_count"},
		{Name: "entity_count"},
	}

	resp, err := wc.GraphQL().Get().
		WithClassName("AcademicPaper").
		WithFields(fields...).
		WithLimit(500).
		Do(context.Background())

	if err != nil {
		log.Printf("CHISGListPapersHandler weaviate error: %v", err)
		http.Error(w, "weaviate query failed", http.StatusInternalServerError)
		return
	}
	if resp.Errors != nil {
		log.Printf("CHISGListPapersHandler weaviate errors: %v", resp.Errors)
	}

	data, _ := resp.Data["Get"].(map[string]interface{})
	rawItems, _ := data["AcademicPaper"].([]interface{})

	papers := make([]CHISGPaperItem, 0, len(rawItems))
	for _, raw := range rawItems {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		papers = append(papers, CHISGPaperItem{
			PaperID:     strVal(m, "paper_id"),
			Title:       strVal(m, "title"),
			Filename:    strVal(m, "filename"),
			Pages:       intVal(m, "pages"),
			TotalChars:  intVal(m, "total_chars"),
			LinkCount:   intVal(m, "link_count"),
			EntityCount: intVal(m, "entity_count"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"papers": papers,
		"count":  len(papers),
	})
}

// ── /api/chisg/links ─────────────────────────────────────────────────────────

// CHISGLinkItem is the response shape for each academic link.
type CHISGLinkItem struct {
	EntityA          string `json:"entity_a"`
	Relation         string `json:"relation"`
	EntityB          string `json:"entity_b"`
	BackwardRelation string `json:"backward_relation"`
	Statement        string `json:"statement"`
	Context          string `json:"context"`
	PaperID          string `json:"paper_id"`
	ChunkID          string `json:"chunk_id"`
	PageNum          int    `json:"page_num"`
}

// CHISGListLinksHandler handles GET /api/chisg/links?paper_id=xxx&limit=500.
// Returns AcademicLink objects from Weaviate, filtered by paper_id when provided.
func CHISGListLinksHandler(w http.ResponseWriter, r *http.Request) {
	wc := dbstore.GetWeaviateClient()
	if wc == nil {
		if err := dbstore.InitializeWeaviateFromEnv(); err != nil {
			http.Error(w, "weaviate unavailable", http.StatusServiceUnavailable)
			return
		}
		wc = dbstore.GetWeaviateClient()
	}

	paperID := r.URL.Query().Get("paper_id")

	limit := 500
	if s := r.URL.Query().Get("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 && v <= 2000 {
			limit = v
		}
	}

	fields := []graphql.Field{
		{Name: "entity_a"},
		{Name: "relation"},
		{Name: "entity_b"},
		{Name: "backward_relation"},
		{Name: "statement"},
		{Name: "context"},
		{Name: "paper_id"},
		{Name: "chunk_id"},
		{Name: "page_num"},
	}

	builder := wc.GraphQL().Get().
		WithClassName("AcademicLink").
		WithFields(fields...).
		WithLimit(limit)

	if paperID != "" {
		where := filters.Where().
			WithPath([]string{"paper_id"}).
			WithOperator(filters.Equal).
			WithValueString(paperID)
		builder = builder.WithWhere(where)
	}

	resp, err := builder.Do(context.Background())
	if err != nil {
		log.Printf("CHISGListLinksHandler weaviate error: %v", err)
		http.Error(w, "weaviate query failed", http.StatusInternalServerError)
		return
	}
	if resp.Errors != nil {
		log.Printf("CHISGListLinksHandler weaviate errors: %v", resp.Errors)
	}

	data, _ := resp.Data["Get"].(map[string]interface{})
	rawItems, _ := data["AcademicLink"].([]interface{})

	links := make([]CHISGLinkItem, 0, len(rawItems))
	for _, raw := range rawItems {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		links = append(links, CHISGLinkItem{
			EntityA:          strVal(m, "entity_a"),
			Relation:         strVal(m, "relation"),
			EntityB:          strVal(m, "entity_b"),
			BackwardRelation: strVal(m, "backward_relation"),
			Statement:        strVal(m, "statement"),
			Context:          strVal(m, "context"),
			PaperID:          strVal(m, "paper_id"),
			ChunkID:          strVal(m, "chunk_id"),
			PageNum:          intVal(m, "page_num"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"links":    links,
		"count":    len(links),
		"paper_id": paperID,
	})
}

// intVal extracts an integer value from a Weaviate response map.
func intVal(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		case int64:
			return int(n)
		}
	}
	return 0
}
