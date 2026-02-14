package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

// WeaviateBackfillVectorsHandler backfills missing vectors for an existing Weaviate class.
//
// This is intended for local/dev use where classes are configured with vectorizer "none" but
// you still want semantic search via nearVector.
//
// Guarded by ENABLE_WEAVIATE_VECTOR_BACKFILL=true.
func WeaviateBackfillVectorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if strings.ToLower(strings.TrimSpace(os.Getenv("ENABLE_WEAVIATE_VECTOR_BACKFILL"))) != "true" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	type reqBody struct {
		Class     string `json:"class,omitempty"`
		Offset    int    `json:"offset,omitempty"`
		BatchSize int    `json:"batch_size,omitempty"`
		DryRun    bool   `json:"dry_run,omitempty"`
	}

	var req reqBody
	_ = json.NewDecoder(r.Body).Decode(&req)

	className := strings.TrimSpace(req.Class)
	if className == "" {
		className = "CHISGElement"
	}

	batchSize := req.BatchSize
	if batchSize <= 0 {
		batchSize = 25
	}
	if batchSize > 100 {
		batchSize = 100
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	ctx := r.Context()
	if err := db.EnsureWeaviateClient(ctx); err != nil {
		http.Error(w, "Weaviate not available", http.StatusServiceUnavailable)
		return
	}
	weaviateClient := db.GetWeaviateClient()
	if weaviateClient == nil {
		http.Error(w, "Weaviate not available", http.StatusServiceUnavailable)
		return
	}

	exists, err := db.WeaviateCollectionExists(ctx, className)
	if err != nil {
		http.Error(w, "Weaviate schema check failed", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, fmt.Sprintf("Weaviate class %q not found", className), http.StatusBadRequest)
		return
	}

	start := time.Now()

	// Fetch total count (best-effort)
	totalCount := 0
	aggResp, aggErr := weaviateClient.GraphQL().Aggregate().
		WithClassName(className).
		WithFields(graphql.Field{Name: "meta", Fields: []graphql.Field{{Name: "count"}}}).
		Do(ctx)
	if aggErr == nil && aggResp != nil && aggResp.Errors == nil {
		if aggBlock, ok := aggResp.Data["Aggregate"].(map[string]interface{}); ok {
			if rowsAny, ok := aggBlock[className].([]interface{}); ok && len(rowsAny) > 0 {
				if row, ok := rowsAny[0].(map[string]interface{}); ok {
					if meta, ok := row["meta"].(map[string]interface{}); ok {
						if count, ok := meta["count"].(float64); ok {
							totalCount = int(count)
						}
					}
				}
			}
		}
	}

	// Get a page of objects
	additional := graphql.Field{Name: "_additional", Fields: []graphql.Field{{Name: "id"}, {Name: "vector"}}}
	var fields []graphql.Field
	switch className {
	case "Documentation":
		fields = []graphql.Field{{Name: "title"}, {Name: "content"}, additional}
	case "SemanticLinks":
		fields = []graphql.Field{{Name: "statement"}, {Name: "context"}, {Name: "source_term"}, {Name: "target_term"}, additional}
	default:
		fields = []graphql.Field{{Name: "name"}, {Name: "description"}, additional}
	}
	getResp, getErr := weaviateClient.GraphQL().Get().
		WithClassName(className).
		WithFields(fields...).
		WithLimit(batchSize).
		WithOffset(offset).
		Do(ctx)
	if getErr != nil {
		http.Error(w, fmt.Sprintf("Weaviate get failed: %v", getErr), http.StatusBadGateway)
		return
	}
	if getResp != nil && getResp.Errors != nil {
		http.Error(w, fmt.Sprintf("Weaviate get returned errors: %v", getResp.Errors), http.StatusBadGateway)
		return
	}

	getBlock, ok := getResp.Data["Get"].(map[string]interface{})
	if !ok {
		http.Error(w, "Unexpected Weaviate response format", http.StatusBadGateway)
		return
	}
	itemsAny, ok := getBlock[className].([]interface{})
	if !ok {
		itemsAny = []interface{}{}
	}

	embedClient := getEmbeddingClient()
	updated := 0
	skipped := 0
	failed := 0
	sampleVectorDim := 0
	failures := make([]map[string]string, 0, 5)

	for _, raw := range itemsAny {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		add, _ := item["_additional"].(map[string]interface{})
		id, _ := add["id"].(string)
		if id == "" {
			continue
		}

		// Skip already-vectorized objects (re-runs are common in dev).
		if existing := add["vector"]; existing != nil {
			switch v := existing.(type) {
			case []interface{}:
				if len(v) > 0 {
					skipped++
					continue
				}
			case []float64:
				if len(v) > 0 {
					skipped++
					continue
				}
			}
		}

		text := ""
		switch className {
		case "Documentation":
			title, _ := item["title"].(string)
			content, _ := item["content"].(string)
			text = strings.TrimSpace(strings.Join([]string{title, content}, "\n\n"))
		case "SemanticLinks":
			stmt, _ := item["statement"].(string)
			ctxText, _ := item["context"].(string)
			sourceTerm, _ := item["source_term"].(string)
			targetTerm, _ := item["target_term"].(string)
			text = strings.TrimSpace(strings.Join([]string{stmt, ctxText, sourceTerm, targetTerm}, "\n\n"))
		default:
			name, _ := item["name"].(string)
			desc, _ := item["description"].(string)
			text = strings.TrimSpace(strings.Join([]string{name, desc}, "\n\n"))
		}
		if text == "" {
			continue
		}

		vector, err := embedClient.GenerateEmbedding(text)
		if err != nil {
			failed++
			if len(failures) < 5 {
				failures = append(failures, map[string]string{"id": id, "error": err.Error()})
			}
			continue
		}
		if sampleVectorDim == 0 {
			sampleVectorDim = len(vector)
		}

		if req.DryRun {
			updated++
			continue
		}

		uerr := weaviateClient.Data().Updater().
			WithClassName(className).
			WithID(id).
			WithVector(vector).
			WithMerge().
			Do(ctx)
		if uerr != nil {
			failed++
			if len(failures) < 5 {
				failures = append(failures, map[string]string{"id": id, "error": uerr.Error()})
			}
			log.Printf("WeaviateBackfillVectorsHandler: failed updating %s/%s: %v", className, id, uerr)
			continue
		}
		updated++
	}

	nextOffset := offset + len(itemsAny)
	elapsed := time.Since(start)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"class":             className,
		"offset":            offset,
		"batch_size":        batchSize,
		"processed":         len(itemsAny),
		"updated":           updated,
		"skipped":           skipped,
		"failed":            failed,
		"next_offset":       nextOffset,
		"total_count":       totalCount,
		"sample_vector_dim": sampleVectorDim,
		"dry_run":           req.DryRun,
		"failures":          failures,
		"elapsed_ms":        elapsed.Milliseconds(),
		"timestamp":         time.Now(),
	})
}

// WeaviateAddDocumentationHandler inserts a single Documentation record into Weaviate.
// Guarded by ENABLE_WEAVIATE_VECTOR_BACKFILL=true (local/dev only).
//
// The handler generates an embedding for the content and stores it as the object's vector so
// the document becomes immediately searchable via nearVector.
func WeaviateAddDocumentationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if strings.ToLower(strings.TrimSpace(os.Getenv("ENABLE_WEAVIATE_VECTOR_BACKFILL"))) != "true" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	type reqBody struct {
		Title        string `json:"title"`
		Content      string `json:"content"`
		FilePath     string `json:"file_path,omitempty"`
		Project      string `json:"project,omitempty"`
		SectionPath  string `json:"section_path,omitempty"`
		HeadingLevel int    `json:"heading_level,omitempty"`
		LastModified string `json:"last_modified,omitempty"`
		ChunkHash    string `json:"chunk_hash,omitempty"`
	}

	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	if title == "" || content == "" {
		http.Error(w, "title and content are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := db.EnsureWeaviateClient(ctx); err != nil {
		http.Error(w, "Weaviate not available", http.StatusServiceUnavailable)
		return
	}
	weaviateClient := db.GetWeaviateClient()
	if weaviateClient == nil {
		http.Error(w, "Weaviate not available", http.StatusServiceUnavailable)
		return
	}

	// Compute chunk hash if not provided
	chunkHash := strings.TrimSpace(req.ChunkHash)
	if chunkHash == "" {
		h := sha256.Sum256([]byte(strings.Join([]string{title, content, strings.TrimSpace(req.FilePath)}, "\n")))
		chunkHash = hex.EncodeToString(h[:])
	}

	lastModified := strings.TrimSpace(req.LastModified)
	if lastModified == "" {
		lastModified = time.Now().Format(time.RFC3339)
	}

	props := map[string]interface{}{
		"title":         title,
		"content":       content,
		"file_path":     strings.TrimSpace(req.FilePath),
		"project":       strings.TrimSpace(req.Project),
		"section_path":  strings.TrimSpace(req.SectionPath),
		"heading_level": req.HeadingLevel,
		"chunk_hash":    chunkHash,
		"last_modified": lastModified,
	}

	vector, err := getEmbeddingClient().GenerateEmbedding(content)
	if err != nil {
		http.Error(w, fmt.Sprintf("embedding generation failed: %v", err), http.StatusBadGateway)
		return
	}

	// Create with explicit vector.
	obj, err := weaviateClient.Data().Creator().
		WithClassName("Documentation").
		WithProperties(props).
		WithVector(vector).
		Do(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("weaviate insert failed: %v", err), http.StatusBadGateway)
		return
	}

	createdID := ""
	if obj != nil && obj.Object != nil {
		createdID = fmt.Sprint(obj.Object.ID)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "ok",
		"class":      "Documentation",
		"id":         createdID,
		"chunk_hash": chunkHash,
		"timestamp":  time.Now(),
	})
}
