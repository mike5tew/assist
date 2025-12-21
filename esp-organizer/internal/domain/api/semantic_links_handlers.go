package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"esp-organizer/internal/domain/infoin"
	"esp-organizer/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Request structures for the semantic links extractor
type semanticLinkSelection struct {
	SourceTerm   string             `json:"source_term"`
	TargetTerm   string             `json:"target_term"`
	RelationType string             `json:"relation_type"`
	TextContext  string             `json:"text_context"`
	Position     map[string]int     `json:"position,omitempty"`
	QualityFlag  string             `json:"quality_flag,omitempty"`
	Confidence   *float64           `json:"confidence,omitempty"`
}

type semanticLinksExtractRequest struct {
	DocumentID string                   `json:"document_id"`
	Selections []semanticLinkSelection `json:"selections"`
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

	inserted := 0
	errs := make([]string, 0)

	for _, sel := range req.Selections {
		link := models.SemanticLink{
			SourceTerm: sel.SourceTerm,
			TargetTerm: sel.TargetTerm,
			RelationType: sel.RelationType,
			Context: sel.TextContext,
			Domain: "manual",
			CreatedAt: time.Now(),
		}

		// Set confidence: explicit override or default high confidence for manual curation
		if sel.Confidence != nil {
			link.Confidence = *sel.Confidence
		} else {
			link.Confidence = 0.95
		}

		// Attach provenance in metadata
		link.Metadata = map[string]interface{}{"document_id": req.DocumentID}
		if sel.Position != nil {
			link.Metadata["position"] = sel.Position
		}
		if sel.QualityFlag != "" {
			link.Metadata["quality_flag"] = sel.QualityFlag
		}

		// Use the service to persist
		if err := svc.StoreSemanticLink(ctx, link); err != nil {
			log.Printf("Failed to store semantic link: %v", err)
			errs = append(errs, err.Error())
			continue
		}
		inserted++
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
