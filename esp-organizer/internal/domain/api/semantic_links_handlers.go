package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"esp-organizer/internal/domain/infoin"
	"esp-organizer/internal/domain/skills"
	"esp-organizer/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Request structures for the semantic links extractor

// ConditionSelection represents a condition attached to a semantic link
type ConditionSelection struct {
	Term            string `json:"term"`
	ForwardRelation string `json:"forward_relation"`
	InverseRelation string `json:"inverse_relation"`
	Required        bool   `json:"required,omitempty"`
	Description     string `json:"description,omitempty"`
}

type semanticLinkSelection struct {
	// The full statement that gets vectorized (e.g., "XLA causes reduction in Ig")
	Statement string `json:"statement"`

	// Terms for graph traversal
	SourceTerm string `json:"source_term"`
	TargetTerm string `json:"target_term"`

	// Bidirectional relations
	ForwardRelation string `json:"forward_relation"`
	InverseRelation string `json:"inverse_relation"`

	// Legacy field - deprecated, use ForwardRelation
	RelationType string `json:"relation_type,omitempty"`

	// Conditions under which the relationship holds
	Conditions []ConditionSelection `json:"conditions,omitempty"`

	// Training data fields
	Chapter     string `json:"chapter,omitempty"`      // Chapter/section reference
	ExcerptText string `json:"excerpt_text,omitempty"` // Source passage this was extracted from
	PageNumber  int    `json:"page_number,omitempty"`  // Page reference

	// Quality scoring for training data
	QualityScore  *float64 `json:"quality_score,omitempty"`  // 0.0-1.0 quality rating
	QualityReason string   `json:"quality_reason,omitempty"` // Explanation for score

	// Context links - IDs of related semantic links
	ContextLinkIDs []string `json:"context_link_ids,omitempty"`

	// Optional metadata
	TextContext string         `json:"text_context,omitempty"` // Deprecated: use Statement
	Position    map[string]int `json:"position,omitempty"`
	QualityFlag string         `json:"quality_flag,omitempty"`
	Confidence  *float64       `json:"confidence,omitempty"`
}

type semanticLinksExtractRequest struct {
	DocumentID  string                  `json:"document_id"`
	SourceID    string                  `json:"source_id,omitempty"`    // MongoDB ObjectID of the source document
	SourceTitle string                  `json:"source_title,omitempty"` // Denormalized title for display
	Domain      string                  `json:"domain,omitempty"`       // e.g., immunology, biology
	Subject     string                  `json:"subject,omitempty"`      // e.g., GCSE Biology
	Selections  []semanticLinkSelection `json:"selections"`
}

// SemanticLinksExtractHandler handles manual extraction submissions from the UI
// POST /api/semantic-links/extract
func SemanticLinksExtractHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req semanticLinksExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	svc, err := infoin.NewSemanticLinkService()
	if err != nil {
		log.Printf("Error initializing semantic link service: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// Prepare counters and collectors
	inserted := 0
	errs := make([]string, 0)
	var storedLinks []models.SemanticLink

	collection := svc.MongoClient.Database.Collection("semantic_links")

	for _, sel := range req.Selections {
		// Basic validation
		if sel.SourceTerm == "" || sel.TargetTerm == "" {
			errs = append(errs, "selection missing required source_term or target_term")
			continue
		}

		// Determine the relation type - prefer ForwardRelation, fallback to RelationType
		forwardRel := sel.ForwardRelation
		if forwardRel == "" {
			forwardRel = sel.RelationType // Backward compatibility
		}
		if forwardRel == "" {
			errs = append(errs, "selection missing forward_relation or relation_type")
			continue
		}

		// Determine statement - prefer Statement, fallback to TextContext
		statement := sel.Statement
		if statement == "" {
			statement = sel.TextContext // Backward compatibility
		}

		// Build conditions from selection
		var conditions []models.Condition
		for _, c := range sel.Conditions {
			conditions = append(conditions, models.Condition{
				Term:            infoin.NormalizeText(c.Term),
				ForwardRelation: c.ForwardRelation,
				InverseRelation: c.InverseRelation,
				Required:        c.Required,
				Description:     c.Description,
			})
		}

		// Parse source ID if provided
		var sourceOID primitive.ObjectID
		if req.SourceID != "" {
			if oid, err := primitive.ObjectIDFromHex(req.SourceID); err == nil {
				sourceOID = oid
			}
		}

		// Determine domain - prefer request-level, fallback to "manual"
		domain := req.Domain
		if domain == "" {
			domain = "manual"
		}

		// Build the semantic link with normalized terms
		link := models.SemanticLink{
			ID:              primitive.NewObjectID(),
			Statement:       statement,
			SourceTerm:      infoin.NormalizeText(sel.SourceTerm),
			TargetTerm:      infoin.NormalizeText(sel.TargetTerm),
			ForwardRelation: forwardRel,
			InverseRelation: sel.InverseRelation,
			RelationType:    forwardRel, // Legacy field for backward compat
			Context:         infoin.NormalizeForVectorization(statement),
			Conditions:      conditions,
			Domain:          domain,
			Subject:         req.Subject,
			CreatedAt:       time.Now(),
			Metadata:        map[string]interface{}{"document_id": req.DocumentID},

			// Training data fields
			SourceID:       sourceOID,
			SourceTitle:    req.SourceTitle,
			Chapter:        sel.Chapter,
			ExcerptText:    sel.ExcerptText,
			PageNumber:     sel.PageNumber,
			Status:         "draft", // New links start as draft
			IsManual:       true,    // UI-created links are manual
			ContextLinkIDs: sel.ContextLinkIDs,
		}

		// Quality scoring
		if sel.QualityScore != nil {
			link.QualityScore = *sel.QualityScore
		} else {
			link.QualityScore = 0.8 // Default quality score
		}
		link.QualityReason = sel.QualityReason

		if sel.Position != nil {
			link.Metadata["position"] = sel.Position
		}
		if sel.QualityFlag != "" {
			link.Metadata["quality_flag"] = sel.QualityFlag
		}
		if sel.Confidence != nil {
			link.Confidence = *sel.Confidence
		} else {
			link.Confidence = 0.95
		}

		// Compute composite key to detect duplicates before attempting insert
		link.CompositeKey = infoin.GenerateCompositeKey(link.SourceTerm, link.TargetTerm, link.ForwardRelation)

		// Check for existing link
		existing := &models.SemanticLink{}
		if err := collection.FindOne(ctx, bson.M{"composite_key": link.CompositeKey}).Decode(existing); err == nil {
			// already exists, skip
			log.Printf("⚠️ Semantic link already exists (skipping): %s", link.CompositeKey)
			continue
		} else if err != mongo.ErrNoDocuments {
			log.Printf("Error checking for existing link: %v", err)
			errs = append(errs, err.Error())
			continue
		}

		// Insert directly into collection (avoid double normalization in StoreSemanticLink)
		if _, err := collection.InsertOne(ctx, link); err != nil {
			log.Printf("Failed to insert semantic link: %v", err)
			errs = append(errs, err.Error())
			continue
		}

		inserted++
		storedLinks = append(storedLinks, link)
	}

	// Asynchronously vectorize / index stored links in Weaviate using SkillService
	if len(storedLinks) > 0 {
		go func(links []models.SemanticLink) {
			skillSvc := skills.NewSkillService(svc.MongoClient.Database.Collection("skills"), svc.WeaviateClient, svc.LlmClient)
			if err := skillSvc.StoreSemanticLinks(context.Background(), links); err != nil {
				log.Printf("Failed to store semantic links in Weaviate: %v", err)
			} else {
				log.Printf("✅ Indexed %d semantic links in Weaviate", len(links))
			}
		}(storedLinks)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"inserted": inserted,
		"errors":   errs,
	})
}

// SemanticLinksSearchHandler provides a convenient search over semantic links
// GET /api/semantic-links/search?term=...&domain=...&max=10
func SemanticLinksSearchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query().Get("term")
	if q == "" {
		http.Error(w, "query param 'term' is required", http.StatusBadRequest)
		return
	}
	domain := r.URL.Query().Get("domain")
	max := 10
	if s := r.URL.Query().Get("max"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			max = v
		}
	}

	svc, err := infoin.NewHSGQueryService()
	if err != nil {
		log.Printf("Error creating HSGQueryService: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	ragContext, answer, err := svc.QueryHSG(ctx, q, domain, max)
	if err != nil {
		log.Printf("HSG query failed: %v", err)
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"semantic_context": ragContext.SemanticContext,
		"summary_context":  ragContext.SummaryContext,
		"documents":        ragContext.DocumentContext,
		"answer":           answer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SemanticLinksValidateHandler is a placeholder for validation and scoring
// POST /api/semantic-links/validate
func SemanticLinksValidateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "validation endpoint not implemented yet"})
}

// SemanticLinkExportItem represents a semantic link in a mobile-app friendly format
type SemanticLinkExportItem struct {
	ID              string                 `json:"id"`
	Statement       string                 `json:"statement"`
	SourceTerm      string                 `json:"source_term"`
	TargetTerm      string                 `json:"target_term"`
	ForwardRelation string                 `json:"forward_relation"`
	InverseRelation string                 `json:"inverse_relation"`
	Conditions      []ConditionExportItem  `json:"conditions,omitempty"`
	Domain          string                 `json:"domain"`
	Subject         string                 `json:"subject,omitempty"`
	Confidence      float64                `json:"confidence"`
	CreatedAt       string                 `json:"created_at"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`

	// Training data fields
	SourceID    string `json:"source_id,omitempty"`
	SourceTitle string `json:"source_title,omitempty"`
	Chapter     string `json:"chapter,omitempty"`
	ExcerptText string `json:"excerpt_text,omitempty"`
	PageNumber  int    `json:"page_number,omitempty"`
	Status      string `json:"status"`
	IsManual    bool   `json:"is_manual"`
	ValidatedBy string `json:"validated_by,omitempty"`
	ValidatedAt string `json:"validated_at,omitempty"`
}

type ConditionExportItem struct {
	Term            string `json:"term"`
	ForwardRelation string `json:"forward_relation"`
	InverseRelation string `json:"inverse_relation"`
}

// SemanticLinksExportHandler exports semantic links in SQLite-compatible JSON format
// GET /api/semantic-links/export?format=json&domain=...&limit=1000
func SemanticLinksExportHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	domain := r.URL.Query().Get("domain")
	limit := 1000
	if s := r.URL.Query().Get("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 && v <= 10000 {
			limit = v
		}
	}

	// Additional filters for training data
	status := r.URL.Query().Get("status")      // draft, validated, rejected
	isManual := r.URL.Query().Get("is_manual") // true, false
	sourceID := r.URL.Query().Get("source_id") // Filter by source document

	svc, err := infoin.NewSemanticLinkService()
	if err != nil {
		log.Printf("Error initializing semantic link service: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// Build query filter
	filter := bson.M{}
	if domain != "" {
		filter["domain"] = domain
	}
	if status != "" {
		filter["status"] = status
	}
	if isManual == "true" {
		filter["is_manual"] = true
	} else if isManual == "false" {
		filter["is_manual"] = false
	}
	if sourceID != "" {
		if oid, err := primitive.ObjectIDFromHex(sourceID); err == nil {
			filter["source_id"] = oid
		}
	}

	collection := svc.MongoClient.Database.Collection("semantic_links")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error querying semantic links: %v", err)
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var links []models.SemanticLink
	if err := cursor.All(ctx, &links); err != nil {
		log.Printf("Error decoding semantic links: %v", err)
		http.Error(w, "decode failed", http.StatusInternalServerError)
		return
	}

	// Apply limit
	if len(links) > limit {
		links = links[:limit]
	}

	// Convert to export format (simpler structure for SQLite import)
	exportItems := make([]SemanticLinkExportItem, 0, len(links))
	for _, link := range links {
		conditions := make([]ConditionExportItem, 0, len(link.Conditions))
		for _, c := range link.Conditions {
			conditions = append(conditions, ConditionExportItem{
				Term:            c.Term,
				ForwardRelation: c.ForwardRelation,
				InverseRelation: c.InverseRelation,
			})
		}

		// Format validated_at if set
		validatedAt := ""
		if !link.ValidatedAt.IsZero() {
			validatedAt = link.ValidatedAt.Format(time.RFC3339)
		}

		// Get source ID as string
		sourceID := ""
		if !link.SourceID.IsZero() {
			sourceID = link.SourceID.Hex()
		}

		exportItems = append(exportItems, SemanticLinkExportItem{
			ID:              link.ID.Hex(),
			Statement:       link.Statement,
			SourceTerm:      link.SourceTerm,
			TargetTerm:      link.TargetTerm,
			ForwardRelation: link.ForwardRelation,
			InverseRelation: link.InverseRelation,
			Conditions:      conditions,
			Domain:          link.Domain,
			Subject:         link.Subject,
			Confidence:      link.Confidence,
			CreatedAt:       link.CreatedAt.Format(time.RFC3339),
			Metadata:        link.Metadata,

			// Training data fields
			SourceID:    sourceID,
			SourceTitle: link.SourceTitle,
			Chapter:     link.Chapter,
			ExcerptText: link.ExcerptText,
			PageNumber:  link.PageNumber,
			Status:      link.Status,
			IsManual:    link.IsManual,
			ValidatedBy: link.ValidatedBy,
			ValidatedAt: validatedAt,
		})
	}

	// Return based on format
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"count":   len(exportItems),
			"links":   exportItems,
			"version": "1.0",
			"schema": map[string]string{
				"table":       "semantic_links",
				"primary_key": "id",
			},
		})
	case "sqlite-sql":
		// Return SQL statements for SQLite import
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=semantic_links.sql")
		// Generate CREATE TABLE and INSERT statements
		sqlStatements := generateSQLiteStatements(exportItems)
		w.Write([]byte(sqlStatements))
	default:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"count":   len(exportItems),
			"links":   exportItems,
			"version": "1.0",
		})
	}
}

// generateSQLiteStatements creates SQLite-compatible SQL for importing semantic links
func generateSQLiteStatements(items []SemanticLinkExportItem) string {
	var sb strings.Builder

	sb.WriteString("-- Semantic Links Export for SQLite/Revision App\n")
	sb.WriteString("-- Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")

	// Create table
	sb.WriteString(`CREATE TABLE IF NOT EXISTS semantic_links (
    id TEXT PRIMARY KEY,
    statement TEXT NOT NULL,
    source_term TEXT NOT NULL,
    target_term TEXT NOT NULL,
    forward_relation TEXT,
    inverse_relation TEXT,
    domain TEXT,
    confidence REAL,
    created_at TEXT,
    conditions_json TEXT
);

CREATE INDEX IF NOT EXISTS idx_source_term ON semantic_links(source_term);
CREATE INDEX IF NOT EXISTS idx_target_term ON semantic_links(target_term);
CREATE INDEX IF NOT EXISTS idx_domain ON semantic_links(domain);

`)

	// Insert statements
	for _, item := range items {
		conditionsJSON, _ := json.Marshal(item.Conditions)
		sb.WriteString("INSERT OR REPLACE INTO semantic_links (id, statement, source_term, target_term, forward_relation, inverse_relation, domain, confidence, created_at, conditions_json) VALUES (")
		sb.WriteString("'" + escapeSQLite(item.ID) + "', ")
		sb.WriteString("'" + escapeSQLite(item.Statement) + "', ")
		sb.WriteString("'" + escapeSQLite(item.SourceTerm) + "', ")
		sb.WriteString("'" + escapeSQLite(item.TargetTerm) + "', ")
		sb.WriteString("'" + escapeSQLite(item.ForwardRelation) + "', ")
		sb.WriteString("'" + escapeSQLite(item.InverseRelation) + "', ")
		sb.WriteString("'" + escapeSQLite(item.Domain) + "', ")
		sb.WriteString(strconv.FormatFloat(item.Confidence, 'f', 4, 64) + ", ")
		sb.WriteString("'" + escapeSQLite(item.CreatedAt) + "', ")
		sb.WriteString("'" + escapeSQLite(string(conditionsJSON)) + "');\n")
	}

	return sb.String()
}

// escapeSQLite escapes single quotes for SQLite
func escapeSQLite(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
