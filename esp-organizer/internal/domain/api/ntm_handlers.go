package api

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/aws/llm"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	dbstore "esp-organizer/internal/store/db"
)

// ── Types ─────────────────────────────────────────────────────────────────────

type NTMQueryRequest struct {
	Question string `json:"question"`
}

type NTMQueryResponse struct {
	Answer    string          `json:"answer"`
	Sources   []NTMSourceSnip `json:"sources"`
	LinksUsed int             `json:"links_used"`
}

type NTMSourceSnip struct {
	EntityA  string `json:"entity_a"`
	Relation string `json:"relation"`
	EntityB  string `json:"entity_b"`
	Context  string `json:"context,omitempty"`
}

type NTMUploadResponse struct {
	JobID    string `json:"job_id"`
	Status   string `json:"status"`
	Filename string `json:"filename"`
	Message  string `json:"message"`
}

type NTMPaper struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	JobID      string             `bson:"job_id" json:"job_id"`
	Filename   string             `bson:"filename" json:"filename"`
	Title      string             `bson:"title" json:"title"`
	Status     string             `bson:"status" json:"status"`
	LinkCount  int                `bson:"link_count" json:"link_count"`
	UploadedAt time.Time          `bson:"uploaded_at" json:"uploaded_at"`
	FilePath   string             `bson:"file_path" json:"file_path"`
}

// ── Upload handler ────────────────────────────────────────────────────────────

// NTMUploadHandler saves an uploaded PDF to the papers staging directory.
func NTMUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, `{"error":"file too large (max 50MB)"}`, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("paper")
	if err != nil {
		http.Error(w, `{"error":"no paper field in form"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(handler.Filename), ".pdf") {
		http.Error(w, `{"error":"only PDF files are accepted"}`, http.StatusBadRequest)
		return
	}

	stagingDir := os.Getenv("NTM_PAPERS_DIR")
	if stagingDir == "" {
		stagingDir = "/data/ntm-papers"
	}
	if err := os.MkdirAll(stagingDir, 0755); err != nil {
		log.Printf("NTMUpload: failed to create staging dir: %v", err)
		http.Error(w, `{"error":"server storage error"}`, http.StatusInternalServerError)
		return
	}

	jobID := primitive.NewObjectID().Hex()
	safeFilename := fmt.Sprintf("%s_%s", jobID, filepath.Base(handler.Filename))
	destPath := filepath.Join(stagingDir, safeFilename)

	dest, err := os.Create(destPath)
	if err != nil {
		log.Printf("NTMUpload: create file error: %v", err)
		http.Error(w, `{"error":"failed to save file"}`, http.StatusInternalServerError)
		return
	}
	defer dest.Close()
	if _, err := io.Copy(dest, file); err != nil {
		log.Printf("NTMUpload: copy error: %v", err)
		http.Error(w, `{"error":"failed to write file"}`, http.StatusInternalServerError)
		return
	}

	// Store metadata in MongoDB
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		title = strings.TrimSuffix(handler.Filename, ".pdf")
	}
	paperDoc := NTMPaper{
		ID:         primitive.NewObjectID(),
		JobID:      jobID,
		Filename:   handler.Filename,
		Title:      title,
		Status:     "queued",
		UploadedAt: time.Now().UTC(),
		FilePath:   destPath,
	}
	if mdb := getNTMMongoDB(); mdb != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		coll := mdb.Database.Collection("ntm_papers")
		if _, err := coll.InsertOne(ctx, paperDoc); err != nil {
			log.Printf("NTMUpload: MongoDB insert warning: %v", err)
		}
	}

	json.NewEncoder(w).Encode(NTMUploadResponse{
		JobID:    jobID,
		Status:   "queued",
		Filename: handler.Filename,
		Message:  "Paper queued for extraction. Run the extraction pipeline to process it.",
	})
}

// NTMListPapersHandler returns a list of uploaded papers from MongoDB.
func NTMListPapersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mdb := getNTMMongoDB()
	if mdb == nil {
		json.NewEncoder(w).Encode([]NTMPaper{})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	coll := mdb.Database.Collection("ntm_papers")
	cur, err := coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "uploaded_at", Value: -1}}).SetLimit(50))
	if err != nil {
		http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)
	var papers []NTMPaper
	if err := cur.All(ctx, &papers); err != nil || papers == nil {
		papers = []NTMPaper{}
	}
	json.NewEncoder(w).Encode(papers)
}

// ── Query handler ─────────────────────────────────────────────────────────────

// NTMQueryHandler searches AcademicLink in Weaviate for relevant triples,
// builds context, and calls Claude to answer the research question.
func NTMQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req NTMQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Question) == "" {
		http.Error(w, `{"error":"question field required"}`, http.StatusBadRequest)
		return
	}
	question := strings.TrimSpace(req.Question)

	// Search Weaviate AcademicLink with BM25
	links, err := searchAcademicLinks(question, 30)
	if err != nil {
		log.Printf("NTMQuery: Weaviate search error: %v", err)
		// Fall through — we can still try Claude with empty context
		links = []NTMSourceSnip{}
	}

	// Build context block for Claude
	var contextLines []string
	for i, l := range links {
		line := fmt.Sprintf("%d. %s %s %s", i+1, l.EntityA, l.Relation, l.EntityB)
		if l.Context != "" {
			line += fmt.Sprintf(" [context: %s]", l.Context)
		}
		contextLines = append(contextLines, line)
	}
	contextBlock := strings.Join(contextLines, "\n")

	systemPrompt := `You are a research assistant specialising in mycobacterial biology, 
specifically the starvation survival mechanisms of Non-Tuberculous Mycobacteria (NTM).
You have access to a knowledge graph extracted from academic papers by the CHISG 
(Contextualised Hierarchical Iterative Semantic Grouping) pipeline.
Answer questions using the provided semantic links as your evidence base.
Be precise, cite entity relationships where relevant, and note if the context 
does not contain enough information to fully answer the question.`

	userMessage := fmt.Sprintf("Research question: %s\n\nRelevant semantic links from the knowledge graph:\n%s\n\nPlease answer the research question using these links as evidence.", question, contextBlock)
	if len(contextLines) == 0 {
		userMessage = fmt.Sprintf("Research question: %s\n\n(No specific links found in the knowledge graph for this query. Please answer from general mycobacteriology knowledge.)", question)
	}

	bedrockClient, err := llm.NewBedrockLlamaClient()
	if err != nil {
		log.Printf("NTMQuery: Bedrock init error: %v", err)
		http.Error(w, `{"error":"LLM service unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	answer, err := bedrockClient.InvokeClaude(systemPrompt, userMessage)
	if err != nil {
		log.Printf("NTMQuery: Claude error: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"LLM error: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(NTMQueryResponse{
		Answer:    answer,
		Sources:   links,
		LinksUsed: len(links),
	})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// searchAcademicLinks performs BM25 keyword search on the AcademicLink Weaviate class.
func searchAcademicLinks(query string, limit int) ([]NTMSourceSnip, error) {
	wc := dbstore.GetWeaviateClient()
	if wc == nil {
		return nil, fmt.Errorf("weaviate client not initialised")
	}

	fields := []graphql.Field{
		{Name: "entity_a"},
		{Name: "relation"},
		{Name: "entity_b"},
		{Name: "context"},
		{Name: "statement"},
	}

	response, err := wc.GraphQL().Get().
		WithClassName("AcademicLink").
		WithFields(fields...).
		WithBM25(wc.GraphQL().Bm25ArgBuilder().
			WithQuery(query).
			WithProperties("statement", "entity_a", "entity_b", "context")).
		WithLimit(limit).
		Do(context.Background())

	if err != nil {
		return nil, fmt.Errorf("weaviate BM25 search: %w", err)
	}
	if response.Errors != nil {
		return nil, fmt.Errorf("weaviate errors: %v", response.Errors)
	}

	data, ok := response.Data["Get"].(map[string]interface{})
	if !ok {
		return []NTMSourceSnip{}, nil
	}
	items, ok := data["AcademicLink"].([]interface{})
	if !ok {
		return []NTMSourceSnip{}, nil
	}

	var results []NTMSourceSnip
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		snip := NTMSourceSnip{
			EntityA:  strVal(m, "entity_a"),
			Relation: strVal(m, "relation"),
			EntityB:  strVal(m, "entity_b"),
			Context:  strVal(m, "context"),
		}
		results = append(results, snip)
	}
	return results, nil
}

func strVal(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getNTMMongoDB() *dbstore.MongoDB {
	mdb, err := dbstore.NewMongoDBFromEnv()
	if err != nil {
		log.Printf("NTM: MongoDB not available: %v", err)
		return nil
	}
	return mdb
}

// ── Extract handler ───────────────────────────────────────────────────────────

var ntmBackwardRelations = map[string]string{
	"regulates":         "is regulated by",
	"causes":            "is caused by",
	"enables":           "is enabled by",
	"inhibits":          "is inhibited by",
	"is a type of":      "has subtype",
	"interacts with":    "interacts with",
	"is a mechanism of": "has mechanism",
	"is found in":       "contains",
	"correlates with":   "correlates with",
	"is composed of":    "is a component of",
	"determines":        "is determined by",
	"is a marker for":   "has marker",
}

type ntmExtractRequest struct {
	JobID string `json:"job_id"`
}

// NTMExtractHandler triggers background CHISG extraction for an uploaded NTM paper.
// POST /api/ntm/extract   body: {"job_id": "..."}
func NTMExtractHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ntmExtractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.JobID) == "" {
		http.Error(w, `{"error":"job_id required"}`, http.StatusBadRequest)
		return
	}

	mdb := getNTMMongoDB()
	if mdb == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var paper NTMPaper
	if err := mdb.Database.Collection("ntm_papers").FindOne(ctx, bson.M{"job_id": req.JobID}).Decode(&paper); err != nil {
		http.Error(w, `{"error":"paper not found — upload the paper first"}`, http.StatusNotFound)
		return
	}

	if paper.Status == "processing" {
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "processing", "job_id": req.JobID, "message": "Extraction already running."})
		return
	}

	go ntmRunExtraction(paper)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "processing",
		"job_id":  req.JobID,
		"message": "Extraction started. Poll /api/ntm/status/" + req.JobID + " for progress.",
	})
}

// NTMStatusHandler returns extraction status for a paper.
// GET /api/ntm/status/{jobID}
func NTMStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jobID := mux.Vars(r)["jobID"]
	if jobID == "" {
		http.Error(w, `{"error":"jobID required"}`, http.StatusBadRequest)
		return
	}
	mdb := getNTMMongoDB()
	if mdb == nil {
		http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var paper NTMPaper
	if err := mdb.Database.Collection("ntm_papers").FindOne(ctx, bson.M{"job_id": jobID}).Decode(&paper); err != nil {
		http.Error(w, `{"error":"paper not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"job_id":     paper.JobID,
		"status":     paper.Status,
		"filename":   paper.Filename,
		"title":      paper.Title,
		"link_count": paper.LinkCount,
	})
}

// ntmRunExtraction runs in a background goroutine: extracts text, calls Claude per chunk,
// inserts AcademicLink objects into Weaviate, then updates MongoDB status.
func ntmRunExtraction(paper NTMPaper) {
	log.Printf("NTMExtract[%s]: starting extraction for %s", paper.JobID, paper.Filename)
	ntmUpdateStatus(paper.JobID, "processing", 0)

	pdfBytes, err := os.ReadFile(paper.FilePath)
	if err != nil {
		log.Printf("NTMExtract[%s]: failed to read file: %v", paper.JobID, err)
		ntmUpdateStatus(paper.JobID, "error: file not found", 0)
		return
	}

	rawText, err := extractTextWithPDFCPU(pdfBytes)
	if err != nil {
		log.Printf("NTMExtract[%s]: text extraction failed: %v", paper.JobID, err)
		ntmUpdateStatus(paper.JobID, "error: could not extract text — ensure the PDF is searchable (not a scanned image)", 0)
		return
	}

	rawText = cleanPDFText(rawText)
	chunks := splitIntoParagraphs(rawText)
	log.Printf("NTMExtract[%s]: %d chunks to process", paper.JobID, len(chunks))

	bedrockClient, err := llm.NewBedrockLlamaClient()
	if err != nil {
		log.Printf("NTMExtract[%s]: Bedrock init failed: %v", paper.JobID, err)
		ntmUpdateStatus(paper.JobID, "error: LLM unavailable", 0)
		return
	}

	wc := dbstore.GetWeaviateClient()
	if wc == nil {
		if initErr := dbstore.InitializeWeaviateFromEnv(); initErr == nil {
			wc = dbstore.GetWeaviateClient()
		}
	}

	if wc != nil {
		_, _ = wc.Data().Creator().
			WithClassName("AcademicPaper").
			WithProperties(map[string]interface{}{
				"paper_id":     paper.JobID,
				"title":        paper.Title,
				"filename":     paper.Filename,
				"pages":        0,
				"total_chars":  len(rawText),
				"link_count":   0,
				"entity_count": 0,
			}).
			Do(context.Background())
	}

	const systemPrompt = `You are a CHISG semantic extraction engine for Non-Tuberculous Mycobacteria (NTM) research.
Extract biological, clinical, and epidemiological relationships from the text.
Return ONLY a raw JSON array — no markdown fences, no explanation, just the array starting with [ and ending with ].
Each object must have exactly these four keys: "entity_a", "relation", "entity_b", "context".

CRITICAL — the "context" field is NOT the source sentence. It encodes two things:
1. The EXPERIMENT that produced this finding (for reproducibility):
   - assay / method (e.g. "qRT-PCR", "western blot", "RNA-seq", "ChIP-seq", "survival assay", "growth curve")
   - figure / panel reference if stated (e.g. "Fig. 3B", "Figure 2A", "Supplementary Fig. S4")
   - e.g. "survival assay Fig. 3B", "qRT-PCR analysis Fig. 4C", "RNA-seq experiment"
2. The CONDITIONS under which the relationship holds:
   - organism / strain (e.g. "in M. smegmatis mc2155", "in M. tuberculosis H37Rv")
   - experimental state (e.g. "under starvation conditions", "during exponential growth")
   - scope (e.g. "in vitro", "in vivo", "under oxidative stress")
   - key qualifiers (e.g. "SigF-dependent", "Clp-mediated")
Combine both into one concise phrase (max 150 characters), experiment first then conditions.
Example: "survival assay Fig. 3B; M. smegmatis mc2155 under starvation"
Example: "qRT-PCR Fig. 4C; M. tuberculosis H37Rv, exponential phase"
Example: "western blot; M. abscessus in vitro, CarD overexpression"
If no experiment type is inferable, state the organism and conditions only.
Never copy the source sentence as the context.

Use ONLY these relation types (exact spelling, words separated by spaces — never underscores):
  regulates, causes, enables, inhibits, is a type of, interacts with, is a mechanism of,
  is found in, correlates with, is composed of, determines, is a marker for, associated with,
  treated with, expressed in, produced by.
Be inclusive — extract all meaningful entity relationships, not just the most obvious ones.
Only return [] if the text is purely references, figure numbers, page headers, or contains no entity relationships at all.`

	totalLinks := 0
	for i, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if len(chunk) < 100 {
			continue
		}

		// Skip reference/bibliography sections — they produce citation noise
		if isReferenceSectionChunk(chunk) {
			log.Printf("NTMExtract[%s]: chunk %d skipped (reference section)", paper.JobID, i)
			continue
		}

		userMsg := fmt.Sprintf("Extract CHISG semantic links from this NTM research text:\n\n%s", chunk)
		if i < 2 {
			preview := chunk
			if len(preview) > 300 {
				preview = preview[:300]
			}
			log.Printf("NTMExtract[%s]: chunk %d text preview: %q", paper.JobID, i, preview)
		}
		response, err := bedrockClient.InvokeClaude(systemPrompt, userMsg)
		if err != nil {
			log.Printf("NTMExtract[%s]: Claude error on chunk %d: %v", paper.JobID, i, err)
			continue
		}
		if i < 2 {
			preview := response
			if len(preview) > 300 {
				preview = preview[:300]
			}
			log.Printf("NTMExtract[%s]: chunk %d Claude response: %q", paper.JobID, i, preview)
		}

		// Strip markdown code fences (```json ... ``` or ``` ... ```)
		response = strings.TrimSpace(response)
		for _, fence := range []string{"```json", "```"} {
			if strings.HasPrefix(response, fence) {
				response = strings.TrimPrefix(response, fence)
				break
			}
		}
		if idx := strings.LastIndex(response, "```"); idx >= 0 {
			response = response[:idx]
		}
		response = strings.TrimSpace(response)
		// Find the JSON array bounds
		if idx := strings.Index(response, "["); idx >= 0 {
			response = response[idx:]
		}
		if idx := strings.LastIndex(response, "]"); idx >= 0 {
			response = response[:idx+1]
		}

		var rawLinks []struct {
			EntityA  string `json:"entity_a"`
			Relation string `json:"relation"`
			EntityB  string `json:"entity_b"`
			Context  string `json:"context"`
		}
		if err := json.Unmarshal([]byte(response), &rawLinks); err != nil {
			log.Printf("NTMExtract[%s]: JSON parse error on chunk %d: %v", paper.JobID, i, err)
			continue
		}

		for _, lk := range rawLinks {
			if lk.EntityA == "" || lk.EntityB == "" || lk.Relation == "" {
				continue
			}
			// Normalize relation: replace underscores with spaces (e.g. "is_a_type_of" → "is a type of")
			lk.Relation = strings.ReplaceAll(lk.Relation, "_", " ")
			backward := ntmBackwardRelations[lk.Relation]
			if backward == "" {
				backward = "related to"
			}
			props := map[string]interface{}{
				"entity_a":          lk.EntityA,
				"relation":          lk.Relation,
				"entity_b":          lk.EntityB,
				"backward_relation": backward,
				"statement":         fmt.Sprintf("%s %s %s", lk.EntityA, lk.Relation, lk.EntityB),
				"context":           lk.Context,
				"paper_id":          paper.JobID,
				"chunk_id":          fmt.Sprintf("chunk_%d", i),
				"page_num":          0,
			}
			if wc != nil {
				if _, wErr := wc.Data().Creator().
					WithClassName("AcademicLink").
					WithProperties(props).
					Do(context.Background()); wErr != nil {
					log.Printf("NTMExtract[%s]: Weaviate insert error: %v", paper.JobID, wErr)
				}
			}
			totalLinks++
		}
		log.Printf("NTMExtract[%s]: chunk %d → %d links (total: %d)", paper.JobID, i, len(rawLinks), totalLinks)
	}

	log.Printf("NTMExtract[%s]: complete — %d links inserted", paper.JobID, totalLinks)
	ntmUpdateStatus(paper.JobID, "complete", totalLinks)
}

// ntmUpdateStatus writes status and link count back to MongoDB.
func ntmUpdateStatus(jobID, status string, linkCount int) {
	mdb := getNTMMongoDB()
	if mdb == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = mdb.Database.Collection("ntm_papers").UpdateOne(
		ctx,
		bson.M{"job_id": jobID},
		bson.M{"$set": bson.M{"status": status, "link_count": linkCount}},
	)
}
