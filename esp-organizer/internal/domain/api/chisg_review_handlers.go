package api

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/aws/llm"
	"esp-organizer/internal/domain/extraction"
	"esp-organizer/internal/store/db"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReviewTask maps to the MongoDB collection chisg_knowledge_base.review_tasks
type ReviewTask struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	PaperID          string              `bson:"paper_id" json:"paper_id"`
	SourceDocumentID *primitive.ObjectID `bson:"source_document_id,omitempty" json:"source_document_id,omitempty"`
	SourceTitle      string              `bson:"source_title,omitempty" json:"source_title,omitempty"`
	DOI              string              `bson:"doi,omitempty" json:"doi,omitempty"`
	ChunkID          string              `bson:"chunk_id" json:"chunk_id"`
	SourceText       string              `bson:"source_text" json:"source_text"`
	ProposedLinks    []SemanticLinkData  `bson:"proposed_links" json:"proposed_links"`
	Status           string              `bson:"status" json:"status"`
	AssignedTo       *string             `bson:"assigned_to" json:"assigned_to"`
	ReviewedBy       *string             `bson:"reviewed_by" json:"reviewed_by"`
	ReviewedAt       *time.Time          `bson:"reviewed_at" json:"reviewed_at"`
	CreatedAt        time.Time           `bson:"created_at" json:"created_at"`
}

// SemanticLinkData represents a proposed semantic link
type SemanticLinkData struct {
	EntityA     string `bson:"entity_a" json:"entity_a"`
	Relation    string `bson:"relation" json:"relation"`
	EntityB     string `bson:"entity_b" json:"entity_b"`
	Context     string `bson:"context" json:"context"`
	SourceQuote string `bson:"source_quote" json:"source_quote"`
	Attribution string `bson:"attribution,omitempty" json:"attribution,omitempty"`
}

// UploadCHISGDocumentHandler receives a text or PDF file, returns immediately, and
// processes Textract + Claude extraction in the background so the reviewer can start
// working on the first chunks while the rest of the paper is still being read.
func UploadCHISGDocumentHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, `{"error": "File too large"}`, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("document")
	if err != nil {
		http.Error(w, `{"error": "No document provided"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. Read file bytes synchronously (fast — just memory copy)
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, `{"error": "Failed to read file"}`, http.StatusInternalServerError)
		return
	}

	isPDF := strings.HasSuffix(strings.ToLower(handler.Filename), ".pdf")
	filename := handler.Filename
	jobID := primitive.NewObjectID().Hex()

	// 3. Parse paper metadata from form fields
	title := strings.TrimSpace(r.FormValue("title"))
	authorsRaw := strings.TrimSpace(r.FormValue("authors"))
	doi := strings.TrimSpace(r.FormValue("doi"))
	year := 0
	if y := strings.TrimSpace(r.FormValue("year")); y != "" {
		year, _ = strconv.Atoi(y)
	}
	journal := strings.TrimSpace(r.FormValue("journal"))

	// Split comma-separated authors
	var authors []string
	for _, a := range strings.Split(authorsRaw, ",") {
		a = strings.TrimSpace(a)
		if a != "" {
			authors = append(authors, a)
		}
	}

	// 4. Upsert SourceDocument (synchronous — fast DB write before returning to client)
	sourceDocID, sourceTitle, err := upsertSourceDocument(title, authors, doi, year, journal, filename)
	if err != nil {
		fmt.Printf("[%s] Warning: could not upsert source document: %v\n", jobID, err)
		// Non-fatal — continue without provenance
		sourceTitle = filename
	}
	paperDOI := doi

	// 3. Return immediately — the client can start polling right away
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"job_id":  jobID,
		"message": "Processing started. Tasks will appear in the review queue as they are extracted.",
	})

	// 4. Everything else runs in the background
	go func() {
		var sourceText string

		if isPDF {
			region := os.Getenv("AWS_REGION")
			if region == "" {
				region = "eu-west-2"
			}
			bucket := os.Getenv("AWS_S3_BUCKET")
			if bucket == "" {
				bucket = "esp-new-organizer-immunology"
			}
			tp, err := extraction.NewTextractProcessor(region, bucket)
			if err != nil {
				fmt.Printf("[%s] Failed to init Textract: %v\n", jobID, err)
				return
			}
			extracted, err := tp.ExtractRawText(context.Background(), buf, jobID)
			if err != nil || strings.TrimSpace(extracted) == "" {
				fmt.Printf("[%s] Textract failed: %v\n", jobID, err)
				return
			}
			sourceText = extracted
		} else {
			sourceText = string(buf)
		}

		// 5. Clean and chunk
		sourceText = cleanPDFText(sourceText)
		paragraphs := splitIntoParagraphs(sourceText)
		fmt.Printf("[%s] Split into %d chunks\n", jobID, len(paragraphs))

		mongoDb, err := db.NewFromEnv()
		if err != nil {
			fmt.Printf("[%s] Background DB connection failed: %v\n", jobID, err)
			return
		}
		defer mongoDb.Client.Disconnect(context.Background())
		coll := mongoDb.Client.Database("chisg_knowledge_base").Collection("review_tasks")

		chunkIndex := 0
		for _, para := range paragraphs {
			para = strings.TrimSpace(para)
			if len(para) < 100 {
				continue
			}

			// 6. Extract semantic links via Claude — insert into MongoDB immediately
			//    so the reviewer sees this chunk the next time they poll
			llmClient := llm.NewLlamaClient()
			prompt := fmt.Sprintf(`Extract up to 5 of the most crucial core semantic relationships STRICTLY from this text segment.
CRITICAL RULES:
1. Do NOT use external knowledge.
2. If the text is a document title, list of authors, university affiliation, email address, or general intro, YOU MUST RETURN [].
3. Only extract from full, declarative scientific sentences describing actual biological mechanisms.
4. Only output a raw JSON array. Do not include markdown codeblocks.
5. If the claim is attributed to a cited source (e.g. "Li et al., 2022", "(Gourse et al., 1996)", "according to Smith 2020"), record the citation string in "attribution". If the claim is made directly by this paper's authors, omit "attribution" or set it to "".

Text: %s

Format each relationship exactly as:
{
  "entity_a": "term A",
  "entity_b": "term B",
  "relation": "causes|leads_to|part_of|contains|develops_into|regulates|enables|inhibits|treats|diagnoses|manifests_as|is_a|located_at",
  "source_quote": "the exact verbatim sentence or clause from the text above that contains this specific link",
  "context": "short biological/taxonomic/conditional context for this link (e.g. 'in Mycobacterium tuberculosis', 'under starvation conditions', 'in Gram-negative bacteria', 'during heat stress', 'in cancer cells') — NOT a verbatim quote",
  "attribution": "cited source if claim is from another paper, e.g. 'Li et al., 2022' — omit if this paper's own claim"
}`, para)

			var links []SemanticLinkData
			response, err := llmClient.Generate(context.Background(), prompt)
			if err == nil && response != "" {
				response = strings.TrimPrefix(response, "```json")
				response = strings.TrimPrefix(response, "```")
				response = strings.TrimSuffix(response, "```")
				response = strings.TrimSpace(response)
				json.Unmarshal([]byte(response), &links)
			}

			if len(links) == 0 {
				continue // Claude found no biological relationships — skip
			}

			// Insert immediately so it's reviewable right now
			newTask := ReviewTask{
				ID:            primitive.NewObjectID(),
				PaperID:       filename,
				SourceTitle:   sourceTitle,
				DOI:           paperDOI,
				ChunkID:       fmt.Sprintf("chunk_%d", chunkIndex),
				SourceText:    para,
				ProposedLinks: links,
				Status:        "pending_review",
				CreatedAt:     time.Now(),
			}
			if sourceDocID != primitive.NilObjectID {
				newTask.SourceDocumentID = &sourceDocID
			}
			coll.InsertOne(context.Background(), newTask)
			fmt.Printf("[%s] Inserted chunk_%d (%d links)\n", jobID, chunkIndex, len(links))
			chunkIndex++
		}
		fmt.Printf("[%s] Finished — %d chunks with links inserted\n", jobID, chunkIndex)
	}()
}

// splitIntoParagraphs

func cleanPDFText(text string) string {
	// 1. Remove copyright phrases
	re版权 := regexp.MustCompile(`(?i)copyright.*?all rights reserved`)
	text = re版权.ReplaceAllString(text, "")
	re版权2 := regexp.MustCompile(`(?i)all rights reserved`)
	text = re版权2.ReplaceAllString(text, "")

	re版权3 := regexp.MustCompile(`(?i)Copyright ©.*?by Annual Reviews[\.]?`)
	text = re版权3.ReplaceAllString(text, "")

	// 2. Remove DOIs or URLs
	reDOI := regexp.MustCompile(`(?i)https?://[^\s]+`)
	text = reDOI.ReplaceAllString(text, "")

	reDOI2 := regexp.MustCompile(`(?i)doi\.org/[^\s]+`)
	text = reDOI2.ReplaceAllString(text, "")

	// Removes numerical journal garbage like 122343 iology
	reNumbers := regexp.MustCompile(`\b\d+\s+iology\b`)
	text = reNumbers.ReplaceAllString(text, "iology")

	// 3. Fix hyphenated words across spaces or line breaks
	reHyphen := regexp.MustCompile(`([a-zA-Z]+)-\s+([a-zA-Z]+)`)
	text = reHyphen.ReplaceAllString(text, "${1}${2}")

	return text
}

// isReferenceSectionChunk returns true when a chunk looks like a bibliography /
// references section rather than body text. Such chunks produce citation noise
// (author names, journal names, DOIs) that pollute the semantic link graph.
func isReferenceSectionChunk(chunk string) bool {
	lines := strings.Split(chunk, "\n")
	total := len(lines)
	if total == 0 {
		return false
	}

	doiCount := 0
	numberedRefCount := 0
	authorPatternCount := 0
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		// DOI pattern
		if strings.Contains(strings.ToLower(l), "doi:") || strings.Contains(l, "https://doi.org") {
			doiCount++
		}
		// Numbered reference: "1. ", "12. ", "[1]", "[12]"
		if len(l) > 3 {
			if (l[0] >= '1' && l[0] <= '9' && (l[1] == '.' || (l[1] >= '0' && l[1] <= '9' && l[2] == '.'))) ||
				(l[0] == '[' && l[1] >= '1' && l[1] <= '9') {
				numberedRefCount++
			}
		}
		// Author-year pattern: "Smith J, Jones A. (2020)" or "Smith et al."
		if strings.Contains(l, "et al.") || strings.Contains(l, "et al,") {
			authorPatternCount++
		}
	}

	// Heuristic: skip if >30% of lines have numbered ref patterns, or chunk has 3+ DOIs
	if doiCount >= 3 {
		return true
	}
	nonEmpty := 0
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			nonEmpty++
		}
	}
	if nonEmpty > 0 && (float64(numberedRefCount)/float64(nonEmpty) > 0.3 || float64(authorPatternCount)/float64(nonEmpty) > 0.3) {
		return true
	}
	return false
}

func splitIntoParagraphs(text string) []string {
	// First try double-newline splits (works for plain TXT files).
	parts := strings.Split(text, "\n\n")
	if len(parts) >= 3 {
		return parts
	}

	// Textract / OCR output uses single newlines — collapse whitespace into
	// one flat string, then chunk every ~1200 chars at a sentence boundary.
	reSpace := regexp.MustCompile(`\s+`)
	flat := strings.TrimSpace(reSpace.ReplaceAllString(text, " "))

	const target = 1200
	var chunks []string
	start := 0
	for start < len(flat) {
		end := start + target
		if end >= len(flat) {
			chunks = append(chunks, flat[start:])
			break
		}
		// Walk forward up to 400 extra chars to find ". A" (period → space → uppercase)
		boundary := -1
		for i := end; i < len(flat) && i < end+400; i++ {
			if flat[i] == '.' && i+2 < len(flat) && flat[i+1] == ' ' && flat[i+2] >= 'A' && flat[i+2] <= 'Z' {
				boundary = i + 1
				break
			}
		}
		if boundary == -1 {
			boundary = end
		}
		chunks = append(chunks, strings.TrimSpace(flat[start:boundary]))
		start = boundary
		for start < len(flat) && flat[start] == ' ' {
			start++
		}
	}
	return chunks
}
func GetPendingReviewTasksHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Add strict checking of User Claims from Context here based on JWT auth

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(context.Background())

	coll := mongoDb.Client.Database("chisg_knowledge_base").Collection("review_tasks")

	// Fetch matching pending tasks limit 10
	importOpts := options.Find().SetLimit(10)
	cursor, err := coll.Find(context.Background(), bson.M{"status": "pending_review"}, importOpts)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch tasks"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var tasks []ReviewTask
	if err := cursor.All(context.Background(), &tasks); err != nil {
		http.Error(w, `{"error": "Failed to decode tasks"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"tasks":   tasks,
	})
}

// ApproveReviewTaskHandler accepts edits, marks task approved, and (eventually) sends to Weaviate
type ReviewApprovalPayload struct {
	ApprovedLinks []SemanticLinkData `json:"approved_links"`
}

func ApproveReviewTaskHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Add strict checking of User Claims from Context here based on JWT auth

	vars := mux.Vars(r)
	taskIDHex := vars["id"]
	taskID, err := primitive.ObjectIDFromHex(taskIDHex)
	if err != nil {
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	var payload ReviewApprovalPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"error": "Invalid payload"}`, http.StatusBadRequest)
		return
	}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(context.Background())

	coll := mongoDb.Client.Database("chisg_knowledge_base").Collection("review_tasks")

	// Read the task first to capture provenance fields for Weaviate
	var task ReviewTask
	if err := coll.FindOne(context.Background(), bson.M{"_id": taskID}).Decode(&task); err != nil {
		http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
		return
	}

	// Update MongoDB with the approved status and the final edited links
	update := bson.M{
		"$set": bson.M{
			"status":         "approved",
			"proposed_links": payload.ApprovedLinks, // Overwriting array with the heavily vetted ones
			"reviewed_at":    time.Now(),
			// "reviewed_by": userFromJWT.ID // To be implemented
		},
	}

	res, err := coll.UpdateByID(context.Background(), taskID, update)
	if err != nil || res.ModifiedCount == 0 {
		http.Error(w, `{"error": "Failed to approve task"}`, http.StatusInternalServerError)
		return
	}

	// Write approved links to Weaviate AcademicLink class
	weaviateClient := db.GetWeaviateClient()
	sourceDocIDStr := ""
	if task.SourceDocumentID != nil {
		sourceDocIDStr = task.SourceDocumentID.Hex()
	}
	for _, link := range payload.ApprovedLinks {
		props := map[string]interface{}{
			"entity_a":           link.EntityA,
			"relation":           link.Relation,
			"entity_b":           link.EntityB,
			"context":            link.Context,
			"statement":          link.SourceQuote,
			"paper_id":           task.PaperID,
			"chunk_id":           task.ChunkID,
			"source_document_id": sourceDocIDStr,
		}
		if link.Attribution != "" {
			props["attribution"] = link.Attribution
		}
		_, wErr := weaviateClient.Data().Creator().
			WithClassName("AcademicLink").
			WithProperties(props).
			Do(context.Background())
		if wErr != nil {
			fmt.Printf("[ApproveTask] Weaviate insert error (%s → %s): %v\n", link.EntityA, link.EntityB, wErr)
			// Non-fatal: continue storing remaining links
		}
	}

	// Log diffs between original AI proposals and human-approved links
	if len(task.ProposedLinks) > 0 {
		corrColl := mongoDb.Client.Database("chisg_knowledge_base").Collection("extraction_corrections")
		for i, approved := range payload.ApprovedLinks {
			if i >= len(task.ProposedLinks) {
				break
			}
			original := task.ProposedLinks[i]
			if original.EntityA != approved.EntityA ||
				original.EntityB != approved.EntityB ||
				original.Relation != approved.Relation ||
				original.Context != approved.Context {
				_, _ = corrColl.InsertOne(context.Background(), bson.M{
					"task_id":     taskID,
					"paper_id":    task.PaperID,
					"chunk_text":  task.SourceText,
					"original":    original,
					"approved":    approved,
					"reviewed_at": time.Now(),
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Task approved.",
		"provenance": map[string]interface{}{
			"source_document_id": func() string {
				if task.SourceDocumentID != nil {
					return task.SourceDocumentID.Hex()
				}
				return ""
			}(),
			"source_title": task.SourceTitle,
			"doi":          task.DOI,
		},
	})
}

func ClearReviewTasksHandler(w http.ResponseWriter, r *http.Request) {
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(context.Background())

	coll := mongoDb.Client.Database("chisg_knowledge_base").Collection("review_tasks")
	_, err = coll.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, `{"error": "Failed to clear queue"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
