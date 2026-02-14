package api

import (
	"context"
	"encoding/json"
	"esp-organizer/internal/domain/skills"
	"esp-organizer/internal/integration"
	"esp-organizer/internal/models"
	"esp-organizer/internal/store/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CoachRespondRequest now includes optional answer for diagnostic checking
type CoachRespondRequest struct {
	Message        string `json:"message"`
	Age            int    `json:"age"`
	Domain         string `json:"domain"`
	StudentAnswer  string `json:"student_answer,omitempty"`  // Student's answer to check
	ExpectedAnswer string `json:"expected_answer,omitempty"` // Correct answer for comparison
	QuestionText   string `json:"question_text,omitempty"`   // The question asked
}

// CoachRespondResponse now includes diagnostic information
type CoachRespondResponse struct {
	// Original user message
	UserMessage string `json:"user_message"`
	Age         int    `json:"age"`
	Domain      string `json:"domain"`

	// HumanOS Barrier Detection
	BarriersDetected []string `json:"barriers_detected"`
	TriggerWords     []string `json:"trigger_words"`
	DevelopmentStage string   `json:"development_stage"`

	// CHISG Knowledge
	CHISGResponse *integration.KnowledgeResponse `json:"chisg_knowledge"`

	// MongoDB Enrichment
	MongoDBEnrichment map[string]interface{} `json:"mongodb_enrichment"`

	// Age Adjustment
	AgeAdjustedSummary string `json:"age_adjusted_summary"`

	// Combined Response
	RecommendedResponse string `json:"recommended_response"`
	InterventionNeeded  bool   `json:"intervention_needed"`
	InterventionScript  string `json:"intervention_script"`

	// Metadata
	ConfidenceScore float64   `json:"confidence_score"`
	Timestamp       time.Time `json:"timestamp"`
	ProcessingMs    int64     `json:"processing_ms"`

	// NEW: Diagnostic learning fields
	IsAnswerCorrect     bool                `json:"is_answer_correct,omitempty"`
	AnswerDiagnostics   *AnswerDiagnostics  `json:"answer_diagnostics,omitempty"`
	MisconceptionLinks  []MisconceptionLink `json:"misconception_links,omitempty"`
	CorrectionStrategy  string              `json:"correction_strategy,omitempty"`
	EncouragingResponse string              `json:"encouraging_response,omitempty"`
}

// AnswerDiagnostics explains why an answer is incorrect
type AnswerDiagnostics struct {
	StudentAnswer    string   `json:"student_answer"`
	CorrectAnswer    string   `json:"correct_answer"`
	ConfidenceLevel  float64  `json:"confidence_level"`
	ErrorType        string   `json:"error_type"` // conceptual, computational, misreading, etc.
	ErrorDescription string   `json:"error_description"`
	RelatedConcepts  []string `json:"related_concepts"`
	CommonMistakes   []string `json:"common_mistakes"`
}

// MisconceptionLink traces where the wrong answer came from
type MisconceptionLink struct {
	MisconceptTopic   string  `json:"misconception_topic"`
	CorrectConcept    string  `json:"correct_concept"`
	Connection        string  `json:"connection"`         // How they're related
	WhyConfusion      string  `json:"why_confusion"`      // Why student confused them
	CorrectiveInsight string  `json:"corrective_insight"` // How to distinguish them
	MongoDBReference  string  `json:"mongodb_reference"`  // Which doc explains this
	ConfidenceScore   float64 `json:"confidence_score"`
}

// FullDiagnosticCoachHandler handles full diagnostic checking. RENAMED to avoid collision.
func FullDiagnosticCoachHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	startTime := time.Now()
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// 1. Parse request
	var req CoachRespondRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set defaults
	if req.Domain == "" {
		req.Domain = "gcse"
	}
	if req.Age == 0 {
		req.Age = 14 // Default GCSE age
	}

	log.Printf("[Coach] Processing message for age %d, domain %s: %s", req.Age, req.Domain, req.Message)

	// Initialize response
	response := &CoachRespondResponse{
		UserMessage:      req.Message,
		Age:              req.Age,
		Domain:           req.Domain,
		Timestamp:        time.Now(),
		BarriersDetected: []string{},
		TriggerWords:     []string{},
	}

	// 2. HumanOS Barrier Detection
	barriers, triggerWords := detectBarriersForMessage(req.Message, req.Age)
	response.BarriersDetected = barriers
	response.TriggerWords = triggerWords
	response.DevelopmentStage = getDevelopmentStage(req.Age)

	log.Printf("[Coach] Barriers detected: %v, Triggers: %v", barriers, triggerWords)

	// 3. NEW: Diagnostic checking if student answer provided
	if req.StudentAnswer != "" && req.ExpectedAnswer != "" {
		diagnostics := performAnswerDiagnostics(ctx, req)
		response.AnswerDiagnostics = diagnostics
		response.IsAnswerCorrect = (strings.TrimSpace(strings.ToLower(req.StudentAnswer)) ==
			strings.TrimSpace(strings.ToLower(req.ExpectedAnswer)))

		// Find misconception links
		if !response.IsAnswerCorrect {
			misconceptionLinks := findMisconceptionLinks(ctx, req, diagnostics)
			response.MisconceptionLinks = misconceptionLinks
			response.CorrectionStrategy = generateCorrectionStrategy(diagnostics, misconceptionLinks, req.Age)
			response.EncouragingResponse = generateEncouragingResponse(diagnostics, req.Age)
		}
	}

	// 4. Extract topic and create CHISG query
	topic := extractTopicFromMessage(req.Message)
	concepts := extractConceptsFromMessage(req.Message)
	audience := ageToAudience(req.Age)

	query := &integration.KnowledgeQuery{
		Topic:          topic,
		Concepts:       concepts,
		TargetAudience: audience,
	}

	log.Printf("[Coach] CHISG Query: topic=%s, concepts=%v, audience=%s", topic, concepts, audience)

	// 5. Query CHISG via client
	chisgClient, err := integration.NewCHISGClient()
	chisgResponse, err := chisgClient.Query(ctx, query)
	if err != nil {
		log.Printf("[Coach] CHISG query failed: %v", err)
		response.CHISGResponse = &integration.KnowledgeResponse{
			Summary:         "I'm having difficulty accessing knowledge right now. Please try again.",
			ConfidenceScore: 0.0,
		}
	} else {
		response.CHISGResponse = chisgResponse
		log.Printf("[Coach] CHISG Response: confidence=%.2f", chisgResponse.ConfidenceScore)
	}

	// 6. MongoDB Enrichment - fetch related documents
	mongoEnrichment := enrichWithMongoDB(ctx, req.Domain, topic, concepts)
	response.MongoDBEnrichment = mongoEnrichment

	log.Printf("[Coach] MongoDB enrichment: %d results found", len(mongoEnrichment))

	// 7. Age Adjustment
	response.AgeAdjustedSummary = adjustForAge(response.CHISGResponse.Summary, req.Age)

	// 8. Barrier-based adjustments
	if len(barriers) > 0 {
		response.InterventionNeeded = true
		response.InterventionScript = generateIntervention(barriers, req.Age)
		response.ConfidenceScore = response.CHISGResponse.ConfidenceScore * 0.9 // Reduce confidence due to barriers
	} else {
		response.ConfidenceScore = response.CHISGResponse.ConfidenceScore
	}

	// 9. Generate combined response
	response.RecommendedResponse = generateCombinedResponse(response)

	// Calculate processing time
	response.ProcessingMs = time.Since(startTime).Milliseconds()

	log.Printf("[Coach] Response generated in %dms", response.ProcessingMs)

	// Return response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[Coach] Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// performAnswerDiagnostics analyzes a student's answer for errors
func performAnswerDiagnostics(ctx context.Context, req CoachRespondRequest) *AnswerDiagnostics {
	diagnostics := &AnswerDiagnostics{
		StudentAnswer: req.StudentAnswer,
		CorrectAnswer: req.ExpectedAnswer,
	}

	studentAnswerLower := strings.ToLower(strings.TrimSpace(req.StudentAnswer))
	correctAnswerLower := strings.ToLower(strings.TrimSpace(req.ExpectedAnswer))

	// Exact match
	if studentAnswerLower == correctAnswerLower {
		diagnostics.ConfidenceLevel = 1.0
		diagnostics.ErrorType = "none"
		return diagnostics
	}

	// Detect error type
	diagnostics.ErrorType = detectErrorType(studentAnswerLower, correctAnswerLower, req.QuestionText)
	diagnostics.ErrorDescription = describeError(diagnostics.ErrorType, req.StudentAnswer, req.ExpectedAnswer)
	diagnostics.ConfidenceLevel = calculateErrorConfidence(diagnostics.ErrorType)

	// Find related concepts
	diagnostics.RelatedConcepts = findRelatedConcepts(ctx, req.Domain, req.StudentAnswer, req.ExpectedAnswer)

	// Find common mistakes
	diagnostics.CommonMistakes = findCommonMistakes(ctx, req.Domain, req.QuestionText)

	return diagnostics
}

// detectErrorType categorizes what kind of mistake was made
func detectErrorType(studentAnswer, correctAnswer, questionText string) string {
	questionLower := strings.ToLower(questionText)
	studentLower := strings.ToLower(studentAnswer)
	correctLower := strings.ToLower(correctAnswer)

	// NEW: Phonetic similarity check (sounds like but spelled wrong)
	if isPhoneticallySimilar(studentLower, correctLower) {
		return "phonetic_error"
	}

	// Conceptual error: Right concept, wrong application
	if strings.Contains(studentLower, "photosynthesis") && strings.Contains(correctLower, "respiration") {
		return "conceptual_confusion"
	}

	// Computational error: Right method, wrong calculation
	if isNumeric(studentAnswer) && isNumeric(correctAnswer) {
		return "computational_error"
	}

	// Misreading error: Misread the question
	if strings.Contains(questionLower, "not") && !strings.Contains(studentLower, "not") {
		return "misreading_question"
	}

	// Partial understanding: Right direction, incomplete answer
	if strings.Contains(correctLower, studentLower) {
		return "incomplete_answer"
	}

	// Opposite/inverse: Got it backwards
	if isOpposite(studentAnswer, correctAnswer) {
		return "inverse_error"
	}

	// Generic wrong answer
	return "conceptual_error"
}

// isPhoneticallySimilar checks if two words sound similar (common mishearing)
func isPhoneticallySimilar(word1, word2 string) bool {
	// Common phonetic confusions
	phoneticPairs := map[string][]string{
		"argon":          {"argan", "organ"},     // Argon vs Argan oil
		"argan":          {"argon", "organ"},     // Argan oil (correct) vs Argon gas
		"their":          {"there", "they're"},   // Homophone confusion
		"affect":         {"effect"},             // Commonly confused
		"accept":         {"except"},             // Sound similar
		"complement":     {"compliment"},         // Sound identical
		"principle":      {"principal"},          // Sound similar
		"discrete":       {"discreet"},           // Sound similar
		"photosynthesis": {"photo", "synthesis"}, // Partial words confused
		"osmosis":        {"osmotic", "osmium"},  // Similar sounding
		"mitosis":        {"meiosis"},            // Sound very similar
		"bacteria":       {"bacterium"},          // Singular/plural confusion
		"phenomena":      {"phenomenon"},         // Singular/plural confusion
	}

	// Check both directions
	if similar, ok := phoneticPairs[word1]; ok {
		for _, s := range similar {
			if word2 == s {
				return true
			}
		}
	}

	if similar, ok := phoneticPairs[word2]; ok {
		for _, s := range similar {
			if word1 == s {
				return true
			}
		}
	}

	// Check for common letter swaps (transposition errors)
	if levenshteinDistance(word1, word2) <= 2 && len(word1) > 3 {
		return true
	}

	return false
}

// levenshteinDistance calculates the edit distance between two strings
// (used to detect similar spellings/pronunciations)
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(nums ...int) int {
	result := nums[0]
	for _, n := range nums {
		if n < result {
			result = n
		}
	}
	return result
}

// describeError explains what went wrong in plain language
func describeError(errorType, studentAnswer, correctAnswer string) string {
	descriptions := map[string]string{
		// NEW: Phonetic error explanation
		"phonetic_error": fmt.Sprintf(
			"Great ear! You said '%s' which sounds very similar to '%s' - an easy mix-up! "+
				"These are actually different things: '%s' is what you're looking for. "+
				"The similarity in pronunciation is why many people confuse these.",
			studentAnswer, correctAnswer, correctAnswer),

		"conceptual_confusion": fmt.Sprintf(
			"You answered '%s' but the correct answer is '%s'. These are related concepts, but they describe different processes. Let's explore what makes them different.",
			studentAnswer, correctAnswer),
		"computational_error": fmt.Sprintf(
			"Your calculation approach looks right, but the final answer '%s' should be '%s'. Let's walk through the steps together.",
			studentAnswer, correctAnswer),
		"misreading_question": fmt.Sprintf(
			"It looks like the question might have asked for something specific. You said '%s' but we're looking for '%s'. Let's read the question carefully together.",
			studentAnswer, correctAnswer),
		"incomplete_answer": fmt.Sprintf(
			"You're on the right track with '%s' - that's part of the answer! But the complete answer is '%s'. Can you see what's missing?",
			studentAnswer, correctAnswer),
		"inverse_error": fmt.Sprintf(
			"Interesting! You said '%s' but it's actually the opposite - '%s'. This is a common mix-up.",
			studentAnswer, correctAnswer),
		"conceptual_error": fmt.Sprintf(
			"You answered '%s', but the concept we're looking for is '%s'. These might seem similar, but there's an important difference.",
			studentAnswer, correctAnswer),
	}

	if desc, ok := descriptions[errorType]; ok {
		return desc
	}

	return fmt.Sprintf("You answered '%s' but the correct answer is '%s'. Let's explore why.", studentAnswer, correctAnswer)
}

// findMisconceptionLinks traces where the wrong answer originated from
func findMisconceptionLinks(ctx context.Context, req CoachRespondRequest, diagnostics *AnswerDiagnostics) []MisconceptionLink {
	links := []MisconceptionLink{}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[Coach] MongoDB connection failed for misconception linking: %v", err)
		return links
	}
	defer mongoDb.Client.Disconnect(ctx)

	// Search for documents that might explain the misconception
	collection := mongoDb.Database.Collection("misconception_links")

	filter := bson.M{
		"domain": req.Domain,
		"$or": []bson.M{
			{"wrong_concept": bson.M{"$regex": req.StudentAnswer, "$options": "i"}},
			{"common_mistakes": bson.M{"$in": []string{req.StudentAnswer}}},
			{"related_to_correct": bson.M{"$regex": req.ExpectedAnswer, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("[Coach] Error searching misconception links: %v", err)
		return links
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	cursor.All(ctx, &docs)

	for _, doc := range docs {
		link := MisconceptionLink{
			MisconceptTopic:   getStringField(doc, "wrong_concept"),
			CorrectConcept:    getStringField(doc, "correct_concept"),
			Connection:        getStringField(doc, "connection"),
			WhyConfusion:      getStringField(doc, "why_confusion"),
			CorrectiveInsight: getStringField(doc, "corrective_insight"),
			MongoDBReference:  getStringField(doc, "_id"),
			ConfidenceScore:   getFloatField(doc, "confidence"),
		}
		links = append(links, link)
	}

	// If no links found, create one based on the error type
	if len(links) == 0 {
		links = append(links, MisconceptionLink{
			MisconceptTopic:   req.StudentAnswer,
			CorrectConcept:    req.ExpectedAnswer,
			Connection:        fmt.Sprintf("Related to %s concept", diagnostics.ErrorType),
			WhyConfusion:      describeCommonConfusion(diagnostics.ErrorType, req.StudentAnswer, req.ExpectedAnswer),
			CorrectiveInsight: generateCorrectiveInsight(diagnostics.ErrorType, req.StudentAnswer, req.ExpectedAnswer),
			ConfidenceScore:   0.7,
		})
	}

	return links
}

// generateCorrectionStrategy creates a pedagogically sound response
func generateCorrectionStrategy(diagnostics *AnswerDiagnostics, links []MisconceptionLink, age int) string {
	// Base strategy
	var strategy string

	switch diagnostics.ErrorType {
	case "phonetic_error":
		strategy = "These words sound almost identical, so it's an understandable mix-up. Let's look at the spelling and pronunciation to see the difference."
	case "conceptual_confusion":
		if len(links) > 0 {
			strategy = fmt.Sprintf(
				"This is a common point of confusion. Let's explore the difference between '%s' and '%s'. The key insight is: %s",
				links[0].MisconceptTopic, links[0].CorrectConcept, links[0].CorrectiveInsight)
		} else {
			strategy = "These concepts are related, which can be tricky. Let's compare them side-by-side to clarify the distinction."
		}
	case "incomplete_answer":
		strategy = fmt.Sprintf(
			"That's a great start! You've correctly identified '%s', which is a key part of the answer. What else is needed to make the answer complete?",
			diagnostics.StudentAnswer)
	case "computational_error":
		strategy = "Your method seems correct, which is the most important part. Let's re-check the calculation together. Can you walk me through your steps?"
	case "misreading_question":
		strategy = "It's easy to misread a question, especially with words like 'not' or 'except'. Let's look at the question again carefully. Do you notice anything you might have missed?"
	case "inverse_error":
		strategy = fmt.Sprintf(
			"You're thinking about the right concept but in the opposite direction. You said '%s', but it's actually '%s'. This is a very common mix-up. Let's clarify why.",
			diagnostics.StudentAnswer, diagnostics.CorrectAnswer)
	default:
		strategy = "I can see the logic in your answer. Let's review the underlying concept, and the correct answer will become clear."
	}

	// Age-based adjustments
	if age < 12 {
		// Younger learners: More scaffolding and encouragement
		switch diagnostics.ErrorType {
		case "phonetic_error":
			return "Those words sound the same, it's tricky! Let's say them out loud and see how they are spelled differently."
		case "incomplete_answer":
			return fmt.Sprintf("You're on the right track with '%s'! That's a big piece of the puzzle. What's the next piece?", diagnostics.StudentAnswer)
		}
	} else if age > 16 {
		// Older learners: More direct and analytical
		switch diagnostics.ErrorType {
		case "conceptual_confusion":
			return "This highlights a critical distinction between two related concepts. Let's break down the specific differences to solidify your understanding."
		case "computational_error":
			return "The methodology is sound. Let's audit the calculation to pinpoint the arithmetic error."
		}
	}

	return strategy
}

// generateEncouragingResponse creates supportive language for wrong answers based on the specific error type.
func generateEncouragingResponse(diagnostics *AnswerDiagnostics, age int) string {
	// Select a base message based on the error type.
	errorSpecificMessages := map[string]string{
		"phonetic_error":       "That's a very easy mistake to make, those words sound so similar!",
		"incomplete_answer":    "You're definitely on the right track! That's a great start.",
		"conceptual_confusion": "I can see why you'd connect those two ideas. It's a common point of confusion.",
		"computational_error":  "Getting the method right is the hardest part, and you've done that! Let's just double-check the numbers.",
		"misreading_question":  "That's a sharp observation, even if it's for a slightly different question. It's an easy detail to miss!",
		"inverse_error":        "You've got the right concept but in reverse. That's a super common mix-up!",
	}

	baseMessage, ok := errorSpecificMessages[diagnostics.ErrorType]
	if !ok {
		baseMessage = "That's a thoughtful answer! You're clearly thinking about this topic."
	}

	// Add an age-appropriate follow-up.
	ageAppropriateFollowUp := ""
	if age < 12 {
		ageAppropriateFollowUp = "Mistakes are how our brains get stronger. Let's figure this out together!"
	} else if age < 16 {
		ageAppropriateFollowUp = "This is actually a really common misconception. You're not alone in thinking this!"
	} else {
		ageAppropriateFollowUp = "This highlights an important nuance in the concept. Let's dig deeper."
	}

	return baseMessage + " " + ageAppropriateFollowUp
}

// enrichWithMongoDB fetches related documents from MongoDB for enrichment
func enrichWithMongoDB(ctx context.Context, domain string, topic string, concepts []string) map[string]interface{} {
	result := make(map[string]interface{})

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[Coach] MongoDB connection failed: %v", err)
		return result
	}
	defer mongoDb.Client.Disconnect(ctx)

	// Determine which collections to search based on domain
	collectionsToSearch := getCollectionsForDomain(domain)

	result["domain"] = domain
	result["topic"] = topic
	result["collections_searched"] = collectionsToSearch

	// Search each collection
	enrichmentByCollection := make(map[string]interface{})

	for _, collName := range collectionsToSearch {
		collection := mongoDb.Database.Collection(collName)

		// Build search filter
		filter := bson.M{
			"$or": []bson.M{
				{"title": bson.M{"$regex": topic, "$options": "i"}},
				{"content": bson.M{"$regex": topic, "$options": "i"}},
				{"term": bson.M{"$regex": topic, "$options": "i"}},
				{"tags": bson.M{"$in": concepts}},
			},
		}

		// Search for documents
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			log.Printf("[Coach] Error searching collection %s: %v", collName, err)
			continue
		}
		defer cursor.Close(ctx)

		var documents []bson.M
		if err = cursor.All(ctx, &documents); err != nil {
			log.Printf("[Coach] Error reading results from %s: %v", collName, err)
			continue
		}

		if len(documents) > 0 {
			// Process documents for LLM consumption
			processedDocs := make([]map[string]interface{}, 0)

			for i, doc := range documents {
				if i >= 5 { // Limit to top 5 per collection
					break
				}

				processedDoc := make(map[string]interface{})

				// Extract key fields for LLM
				if id, ok := doc["_id"].(primitive.ObjectID); ok {
					processedDoc["id"] = id.Hex()
				}
				if title, ok := doc["title"].(string); ok {
					processedDoc["title"] = title
				}
				if content, ok := doc["content"].(string); ok {
					// Truncate content for efficiency
					if len(content) > 500 {
						processedDoc["content_preview"] = content[:500] + "..."
					} else {
						processedDoc["content"] = content
					}
				}
				if contentType, ok := doc["content_type"].(string); ok {
					processedDoc["type"] = contentType
				}
				if term, ok := doc["term"].(string); ok {
					processedDoc["term"] = term
				}
				if definition, ok := doc["definition"].(string); ok {
					processedDoc["definition"] = definition
				}
				if source, ok := doc["source"].(bson.M); ok {
					processedDoc["source"] = source
				}

				processedDocs = append(processedDocs, processedDoc)
			}

			enrichmentByCollection[collName] = map[string]interface{}{
				"count":       len(documents),
				"documents":   processedDocs,
				"total_found": len(documents),
			}

			log.Printf("[Coach] Found %d documents in %s", len(documents), collName)
		}
	}

	result["enrichment"] = enrichmentByCollection
	return result
}

// getCollectionsForDomain returns the MongoDB collections to search for a domain
func getCollectionsForDomain(domain string) []string {
	collections := map[string][]string{
		"immunology": {"immunology_content", "immunology_terms", "case_studies"},
		"gcse":       {"gcse_content", "subject_content"},
		"medical":    {"immunology_content", "medical_excerpts", "medical_terms"},
		"general":    {"subject_content", "skills"},
	}

	if colls, ok := collections[domain]; ok {
		return colls
	}

	return []string{"subject_content", "skills"} // Default collections
}

// detectBarriersForMessage detects barriers from the message
func detectBarriersForMessage(message string, age int) ([]string, []string) {
	barriers := []string{}
	triggers := []string{}

	messageLower := strings.ToLower(message)

	// Detect confusion/uncertainty
	confusionWords := []string{"don't understand", "confused", "what is", "how do", "why", "help me", "explain"}
	for _, word := range confusionWords {
		if strings.Contains(messageLower, word) {
			barriers = append(barriers, "confusion")
			triggers = append(triggers, word)
			break
		}
	}

	// Detect low confidence
	lowConfidenceWords := []string{"not sure", "maybe", "probably", "might", "possibly", "think"}
	for _, word := range lowConfidenceWords {
		if strings.Contains(messageLower, word) {
			barriers = append(barriers, "low_confidence")
			triggers = append(triggers, word)
			break
		}
	}

	// Detect frustration
	frustrationWords := []string{"frustrated", "annoyed", "hate", "can't", "stuck", "terrible"}
	for _, word := range frustrationWords {
		if strings.Contains(messageLower, word) {
			barriers = append(barriers, "frustration")
			triggers = append(triggers, word)
			break
		}
	}

	// Age-specific barriers
	if age < 8 {
		barriers = append(barriers, "age_too_young")
	}

	return removeDuplicates(barriers), removeDuplicates(triggers)
}

// getDevelopmentStage returns the developmental stage for an age
func getDevelopmentStage(age int) string {
	if age < 5 {
		return "early_childhood"
	} else if age < 8 {
		return "middle_childhood"
	} else if age < 12 {
		return "late_childhood"
	} else if age < 18 {
		return "adolescence"
	}
	return "adult"
}

// extractTopicFromMessage extracts the main topic from a message
func extractTopicFromMessage(message string) string {
	// Simple extraction - in production, use NLP
	messageLower := strings.ToLower(message)

	// Common question starters
	words := strings.Fields(messageLower)
	if len(words) > 3 {
		// Skip common words and return next meaningful word
		skipWords := map[string]bool{"what": true, "how": true, "why": true, "is": true, "do": true, "can": true, "the": true, "a": true, "an": true}
		for _, word := range words {
			if !skipWords[word] && len(word) > 2 {
				return strings.TrimSuffix(word, "?")
			}
		}
	}

	return message // Fallback to full message
}

// extractConceptsFromMessage extracts related concepts from a message
func extractConceptsFromMessage(message string) []string {
	// Simple concept extraction - in production, use NLP
	concepts := []string{}

	// Common concept keywords
	conceptKeywords := map[string]string{
		"photosynthesis": "chlorophyll,glucose,oxygen,light",
		"immunology":     "antibody,antigen,immunity,infection",
		"cell":           "nucleus,mitochondria,membrane",
		"energy":         "atp,glucose,metabolism",
	}

	messageLower := strings.ToLower(message)

	for topic, relatedConcepts := range conceptKeywords {
		if strings.Contains(messageLower, topic) {
			concepts = append(concepts, strings.Split(relatedConcepts, ",")...)
			break
		}
	}

	return removeDuplicates(concepts)
}

// ageToAudience converts age to CHISG audience string
func ageToAudience(age int) string {
	if age < 11 {
		return "primary_student"
	} else if age < 16 {
		return "gcse_student"
	} else if age < 19 {
		return "alevel_student"
	}
	return "university_student"
}

// adjustForAge adjusts text for age appropriateness
func adjustForAge(text string, age int) string {
	if age < 8 {
		// Use simpler words and shorter sentences
		text = strings.ReplaceAll(text, "photosynthesis", "how plants make food")
		text = strings.ReplaceAll(text, "molecule", "tiny piece")
		text = strings.ReplaceAll(text, "cellular", "cell")
	} else if age < 12 {
		// Slightly simplified but more accurate
		text = strings.ReplaceAll(text, "photosynthesis", "the process plants use to make food from sunlight")
	}

	return text
}

// generateIntervention creates supportive language based on barriers
// generateIntervention creates supportive language based on barriers and age.
func generateIntervention(barriers []string, age int) string {
	// Default scripts for a general audience (e.g., teens)
	scripts := map[string]string{
		"confusion":      "That's a complex topic, which can be confusing. Let me break it down into smaller, easier parts for you.",
		"low_confidence": "It's completely normal to feel unsure about this. You're asking great questions, which is the first step to understanding. Let's build that confidence together.",
		"frustration":    "I understand this can feel tough and frustrating. It's okay to feel that way. Take a deep breath, and we'll work through it step-by-step.",
		"age_too_young":  "This topic is a bit advanced for your age. Let's start with some of the basic ideas first to build a strong foundation.",
	}

	// Age-specific overrides for younger learners
	if age < 12 {
		scripts["confusion"] = "I know this seems tricky! Let's look at it in a simpler way."
		scripts["low_confidence"] = "It's okay to not be sure! Asking for help is smart. We can figure this out together."
		scripts["frustration"] = "It's frustrating when things are hard. Don't worry, we'll get through this. Every expert was once a beginner!"
	}

	// Age-specific overrides for older learners
	if age > 16 {
		scripts["confusion"] = "This is a nuanced topic. Let's deconstruct it to clarify the core components."
		scripts["low_confidence"] = "I see you're being cautious with your answer, which is a good analytical trait. Let's review the evidence to solidify your position."
		scripts["frustration"] = "I can see this part is a bottleneck. Let's isolate the issue and resolve it. Frustration is often a sign you're close to a breakthrough."
	}

	messages := []string{}
	for _, barrier := range barriers {
		if script, ok := scripts[barrier]; ok {
			messages = append(messages, script)
		}
	}

	if len(messages) > 0 {
		return strings.Join(messages, " ")
	}

	return "I'm here to help you learn!"
}

// generateCombinedResponse creates the final response for the student
func generateCombinedResponse(resp *CoachRespondResponse) string {
	// Start with intervention if needed
	response := ""

	if resp.InterventionNeeded {
		response = resp.InterventionScript + "\n\n"
	}

	// Add age-adjusted knowledge
	response += resp.AgeAdjustedSummary

	// Add MongoDB enrichment if available
	if enrichment, ok := resp.MongoDBEnrichment["enrichment"].(map[string]interface{}); ok && len(enrichment) > 0 {
		response += "\n\nHere's what I found in our resources:\n"

		for collName, collData := range enrichment {
			if data, ok := collData.(map[string]interface{}); ok {
				if count, ok := data["count"].(int); ok && count > 0 {
					response += fmt.Sprintf("- From %s (%d items found)\n", collName, count)
				}
			}
		}
	}

	// Add related topics from CHISG
	if resp.CHISGResponse != nil && len(resp.CHISGResponse.RelatedTopics) > 0 {
		response += "\nYou might also want to explore:\n"
		for topic, description := range resp.CHISGResponse.RelatedTopics {
			response += fmt.Sprintf("- %s: %s\n", topic, description)
		}
	}

	return response
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return err == nil
}

func isOpposite(s1, s2 string) bool {
	opposites := map[string]string{
		"increase": "decrease",
		"up":       "down",
		"more":     "less",
		"true":     "false",
		"yes":      "no",
	}

	for key, opposite := range opposites {
		if strings.Contains(strings.ToLower(s1), key) && strings.Contains(strings.ToLower(s2), opposite) {
			return true
		}
	}
	return false
}

func getStringField(doc bson.M, field string) string {
	if val, ok := doc[field].(string); ok {
		return val
	}
	return ""
}

func getFloatField(doc bson.M, field string) float64 {
	if val, ok := doc[field].(float64); ok {
		return val
	}
	return 0.0
}

func describeCommonConfusion(errorType, wrong, correct string) string {
	confusions := map[string]string{
		"conceptual_confusion": fmt.Sprintf(
			"'%s' and '%s' are often confused because they both relate to similar biological/chemical processes, but they operate in opposite directions.",
			wrong, correct),
		"inverse_error": fmt.Sprintf(
			"'%s' is the inverse of '%s'. It's easy to get these backwards when you're first learning.",
			wrong, correct),
	}

	if desc, ok := confusions[errorType]; ok {
		return desc
	}

	return fmt.Sprintf("Students sometimes confuse '%s' with '%s'.", wrong, correct)
}

func generateCorrectiveInsight(errorType, wrong, correct string) string {
	insights := map[string]string{
		"conceptual_confusion": fmt.Sprintf(
			"Remember: '%s' is when..., while '%s' is when...",
			wrong, correct),
		"incomplete_answer": fmt.Sprintf(
			"You've got the first part right ('%s'), but the complete answer also includes '%s'.",
			wrong, correct),
	}

	if insight, ok := insights[errorType]; ok {
		return insight
	}

	return fmt.Sprintf("The key difference between '%s' and '%s' is...", wrong, correct)
}

// calculateErrorConfidence determines confidence level based on error type
func calculateErrorConfidence(errorType string) float64 {
	// Phonetic errors are high-confidence - we're sure it's just a mix-up
	confidenceMap := map[string]float64{
		"phonetic_error":       0.95, // Very confident it's just a sound-alike issue
		"computational_error":  0.90, // Method is right, calculation is wrong
		"inverse_error":        0.85, // Clear reversal detected
		"misreading_question":  0.80, // Likely misread
		"incomplete_answer":    0.75, // Partial understanding
		"conceptual_confusion": 0.70, // Confused concepts
		"conceptual_error":     0.60, // Generic wrong answer
	}

	if confidence, ok := confidenceMap[errorType]; ok {
		return confidence
	}

	return 0.5
}

// findRelatedConcepts finds concepts related to both the wrong and correct answers
func findRelatedConcepts(ctx context.Context, domain string, wrongAnswer string, correctAnswer string) []string {
	concepts := []string{}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[Coach] MongoDB connection failed: %v", err)
		return concepts
	}
	defer mongoDb.Client.Disconnect(ctx)

	collection := mongoDb.Database.Collection("semantic_links")

	// Search for concepts related to both answers
	filter := bson.M{
		"domain": domain,
		"$or": []bson.M{
			{"source_term": bson.M{"$regex": wrongAnswer, "$options": "i"}},
			{"target_term": bson.M{"$regex": wrongAnswer, "$options": "i"}},
			{"source_term": bson.M{"$regex": correctAnswer, "$options": "i"}},
			{"target_term": bson.M{"$regex": correctAnswer, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("[Coach] Error searching semantic links: %v", err)
		return concepts
	}
	defer cursor.Close(ctx)

	var links []bson.M
	cursor.All(ctx, &links)

	// Extract unique concepts
	conceptMap := make(map[string]bool)
	for _, link := range links {
		if source, ok := link["source_term"].(string); ok && source != "" {
			conceptMap[source] = true
		}
		if target, ok := link["target_term"].(string); ok && target != "" {
			conceptMap[target] = true
		}
	}

	for concept := range conceptMap {
		concepts = append(concepts, concept)
	}

	return concepts
}

// findCommonMistakes searches for documented patterns of misunderstandings
func findCommonMistakes(ctx context.Context, domain string, questionText string) []string {
	mistakes := []string{}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[Coach] MongoDB connection failed: %v", err)
		return mistakes
	}
	defer mongoDb.Client.Disconnect(ctx)

	collection := mongoDb.Database.Collection("misunderstanding_map")

	// Search for documented misunderstandings in this domain
	filter := bson.M{
		"domain": domain,
		"$or": []bson.M{
			{"symptom_pattern": bson.M{"$regex": questionText, "$options": "i"}},
			{"question_keywords": bson.M{"$in": extractKeywords(questionText)}},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("[Coach] Error searching misunderstanding map: %v", err)
		return mistakes
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	cursor.All(ctx, &docs)

	// Extract common mistakes
	for _, doc := range docs {
		if wrongAnswer, ok := doc["wrong_answer"].(string); ok {
			mistakes = append(mistakes, wrongAnswer)
		}
	}

	return mistakes
}

// CreateMisunderstandingRecord creates a record of a student's misunderstanding
// This should be called when a misconception is detected
func CreateMisunderstandingRecord(ctx context.Context, req CoachRespondRequest, diagnostics *AnswerDiagnostics) error {
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[Coach] MongoDB connection failed: %v", err)
		return err
	}
	defer mongoDb.Client.Disconnect(ctx)

	collection := mongoDb.Database.Collection("misunderstanding_map")

	record := bson.M{
		"wrong_answer":              req.StudentAnswer,
		"correct_answer":            req.ExpectedAnswer,
		"domain":                    req.Domain,
		"age_group":                 getAgeGroup(req.Age),
		"error_type":                diagnostics.ErrorType,
		"symptom_pattern":           diagnostics.ErrorDescription,
		"root_causes":               identifyRootCauses(diagnostics),
		"correct_understanding":     req.ExpectedAnswer,
		"repair_strategy":           generateRepairStrategy(diagnostics.ErrorType),
		"related_misunderstandings": findRelatedMisunderstandings(ctx, req.Domain, req.StudentAnswer, req.ExpectedAnswer),
		"evidence_patterns":         []string{},
		"question_keywords":         extractKeywords(req.QuestionText),
		"frequency":                 1, // Will be incremented on repeated mistakes
		"first_observed":            time.Now(),
		"last_observed":             time.Now(),
	}

	_, err = collection.InsertOne(ctx, record)
	if err != nil {
		log.Printf("[Coach] Error creating misunderstanding record: %v", err)
		return err
	}

	log.Printf("[Coach] Created misunderstanding record for: %s -> %s", req.StudentAnswer, req.ExpectedAnswer)
	return nil
}

// UpdateMisunderstandingFrequency increments the frequency count for a known misunderstanding
func UpdateMisunderstandingFrequency(ctx context.Context, wrongAnswer string, domain string) error {
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[Coach] MongoDB connection failed: %v", err)
		return err
	}
	defer mongoDb.Client.Disconnect(ctx)

	collection := mongoDb.Database.Collection("misunderstanding_map")

	update := bson.M{
		"$inc": bson.M{"frequency": 1},
		"$set": bson.M{"last_observed": time.Now()},
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{
			"wrong_answer": wrongAnswer,
			"domain":       domain,
		},
		update,
	)

	if err != nil {
		log.Printf("[Coach] Error updating misunderstanding frequency: %v", err)
		return err
	}

	return nil
}

// identifyRootCauses analyzes the error to determine root causes
func identifyRootCauses(diagnostics *AnswerDiagnostics) []string {
	causes := []string{}

	rootCausesMap := map[string][]string{
		"phonetic_error": {
			"homophone_confusion",
			"auditory_processing",
			"similar_pronunciation",
		},
		"conceptual_confusion": {
			"incomplete_concept_understanding",
			"related_concepts_merged",
			"process_direction_confusion",
		},
		"computational_error": {
			"calculation_mistake",
			"arithmetic_error",
			"formula_misapplication",
		},
		"misreading_question": {
			"question_comprehension_failure",
			"instruction_not_followed",
			"negation_missed",
		},
		"incomplete_answer": {
			"partial_understanding",
			"missing_prerequisite_knowledge",
			"incomplete_recall",
		},
		"inverse_error": {
			"direction_reversal",
			"opposite_concept_selection",
			"bidirectional_confusion",
		},
	}

	if rootCauses, ok := rootCausesMap[diagnostics.ErrorType]; ok {
		causes = append(causes, rootCauses...)
	}

	return causes
}

// generateRepairStrategy creates a targeted repair strategy for the misunderstanding
func generateRepairStrategy(errorType string) string {
	strategies := map[string]string{
		"phonetic_error":       "Teach pronunciation and spelling differences through auditory/visual comparison. Use minimal pairs to highlight distinctions.",
		"conceptual_confusion": "Create explicit contrasts between the two concepts. Use analogies and visual representations to show differences.",
		"computational_error":  "Walk through the calculation step-by-step. Verify each intermediate result before proceeding.",
		"misreading_question":  "Teach active reading strategies. Emphasize keywords like 'not', 'except', 'all', 'some'. Practice question decomposition.",
		"incomplete_answer":    "Build on the correct part. Ask guiding questions to elicit missing elements. Provide prerequisite review if needed.",
		"inverse_error":        "Explicitly teach the directionality of the concept. Use arrows, timelines, or directional language to reinforce order.",
	}

	if strategy, ok := strategies[errorType]; ok {
		return strategy
	}

	return "Provide targeted instruction addressing the specific misconception. Use multiple modalities and examples."
}

// findRelatedMisunderstandings finds other common mistakes similar to this one
func findRelatedMisunderstandings(ctx context.Context, domain string, wrongAnswer string, correctAnswer string) []map[string]interface{} {
	related := []map[string]interface{}{}

	mongoDb, err := db.NewFromEnv()
	if err != nil {
		log.Printf("[Coach] MongoDB connection failed: %v", err)
		return related
	}
	defer mongoDb.Client.Disconnect(ctx)

	collection := mongoDb.Database.Collection("misunderstanding_map")

	// Find similar mistakes in the same domain
	filter := bson.M{
		"domain": domain,
		"$or": []bson.M{
			{"root_causes": bson.M{"$elemMatch": bson.M{}}}, // Same root causes
			{"error_type": bson.M{"$regex": getErrorCategory(wrongAnswer, correctAnswer), "$options": "i"}},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("[Coach] Error searching related misunderstandings: %v", err)
		return related
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	cursor.All(ctx, &docs)

	// Return top 5 related misunderstandings
	for i, doc := range docs {
		if i >= 5 {
			break
		}

		related = append(related, map[string]interface{}{
			"wrong_answer":       getStringField(doc, "wrong_answer"),
			"correct_answer":     getStringField(doc, "correct_answer"),
			"symptom_pattern":    getStringField(doc, "symptom_pattern"),
			"repair_strategy":    getStringField(doc, "repair_strategy"),
			"frequency_observed": doc["frequency"],
		})
	}

	return related
}

// getErrorCategory determines the general category of an error
func getErrorCategory(wrongAnswer, correctAnswer string) string {
	if strings.HasPrefix(strings.ToLower(correctAnswer), strings.ToLower(wrongAnswer)) {
		return "incomplete"
	}
	if isOpposite(wrongAnswer, correctAnswer) {
		return "inverse"
	}
	if isPhoneticallySimilar(strings.ToLower(wrongAnswer), strings.ToLower(correctAnswer)) {
		return "phonetic"
	}
	return "conceptual"
}

// getAgeGroup categorizes age into learning stages
func getAgeGroup(age int) string {
	if age < 6 {
		return "preschool"
	} else if age < 12 {
		return "elementary"
	} else if age < 16 {
		return "middle_school"
	} else if age < 19 {
		return "high_school"
	}
	return "higher_education"
}

// extractKeywords extracts searchable keywords from text
func extractKeywords(text string) []string {
	keywords := []string{}

	// Split into words
	words := strings.Fields(strings.ToLower(text))

	// Filter out common words
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "is": true, "are": true, "was": true, "were": true, "be": true, "been": true,
		"do": true, "does": true, "did": true, "will": true, "would": true, "could": true, "should": true, "may": true, "might": true,
		"what": true, "when": true, "where": true, "why": true, "how": true, "which": true, "who": true, "whom": true,
		"and": true, "or": true, "but": true, "if": true, "of": true, "in": true, "on": true, "at": true, "to": true, "for": true,
		"that": true, "this": true, "these": true, "those": true, "i": true, "you": true, "he": true, "she": true, "it": true, "we": true, "they": true,
	}

	for _, word := range words {
		// Clean punctuation
		word = strings.Trim(word, ".,!?;:")

		// Keep words longer than 3 chars and not stop words
		if len(word) > 3 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}

	return removeDuplicates(keywords)
}

// SemanticQueryHandler performs a lightweight skill search:
// - MongoDB exact-ish matches (name search)
// - Weaviate semantic matches (nearVector)
func SemanticQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Request contract kept minimal for the current frontend.
	// Supports optional fields for future extension.
	type semanticQueryRequest struct {
		Query string `json:"query"`
		Limit int    `json:"limit,omitempty"`
		Class string `json:"class,omitempty"`
	}

	var req semanticQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	queryText := strings.TrimSpace(req.Query)
	if queryText == "" {
		http.Error(w, "Query text cannot be empty", http.StatusBadRequest)
		return
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	ctx := r.Context()

	// --- Mongo exact matches (name regex) ---
	mongoDb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mongoDb.Client.Disconnect(ctx)

	skillsCollectionName := os.Getenv("SKILLS_COLLECTION")
	if skillsCollectionName == "" {
		skillsCollectionName = "skills"
	}
	skillSvc := skills.NewSkillService(mongoDb.Database.Collection(skillsCollectionName), nil, nil)

	exactSkills, err := skillSvc.FindSkillsByName(ctx, queryText)
	if err != nil {
		log.Printf("SemanticQueryHandler: Mongo exact search failed: %v", err)
		exactSkills = []models.Skill{}
	}

	// --- Related skills (parents of exact matches) ---
	parentIDs := make([]primitive.ObjectID, 0, 32)
	seen := make(map[primitive.ObjectID]struct{})
	for _, s := range exactSkills {
		for _, pid := range s.ParentSkillIDs {
			if _, ok := seen[pid]; ok {
				continue
			}
			seen[pid] = struct{}{}
			parentIDs = append(parentIDs, pid)
		}
	}

	relatedSkills, err := skillSvc.FindSkillsByIDs(ctx, parentIDs)
	if err != nil {
		log.Printf("SemanticQueryHandler: related skills lookup failed: %v", err)
		relatedSkills = []models.Skill{}
	}

	// --- Weaviate semantic matches (NearVector) ---
	weaviateClient := db.GetWeaviateClient()
	semanticMatches := make([]map[string]interface{}, 0)

	classToUse := strings.TrimSpace(req.Class)
	if classToUse == "" {
		// Choose the most “skill-like” class available.
		// NOTE: In the humanOS Weaviate instance, the skills class is typically CHISGElement.
		candidates := []string{"CHISGElement", "CHISGSkill", "EducationalSkills", "SemanticLinks", "Documentation"}
		for _, candidate := range candidates {
			exists, existsErr := db.WeaviateCollectionExists(ctx, candidate)
			if existsErr != nil {
				log.Printf("SemanticQueryHandler: failed to check Weaviate schema for %s: %v", candidate, existsErr)
				continue
			}
			if exists {
				classToUse = candidate
				break
			}
		}
	}

	if weaviateClient == nil {
		log.Printf("SemanticQueryHandler: Weaviate client not available")
	} else if classToUse != "" {
		// Some Weaviate instances (like the humanOS one) do not support nearText, so we use nearVector.
		embedClient := getEmbeddingClient()
		queryVector, err := embedClient.GenerateEmbedding(queryText)
		if err != nil {
			log.Printf("SemanticQueryHandler: embedding generation failed: %v", err)
			queryVector = nil
		}

		// Best-effort: infer vector dimension from Weaviate and resize the query vector if needed.
		if queryVector != nil {
			expectedDim := 0
			dimResp, dimErr := weaviateClient.GraphQL().Get().
				WithClassName(classToUse).
				WithFields(graphql.Field{Name: "_additional", Fields: []graphql.Field{{Name: "vector"}}}).
				WithLimit(1).
				Do(ctx)
			if dimErr == nil && dimResp != nil && dimResp.Errors == nil {
				if getBlock, ok := dimResp.Data["Get"].(map[string]interface{}); ok {
					if itemsAny, ok := getBlock[classToUse].([]interface{}); ok && len(itemsAny) > 0 {
						if item, ok := itemsAny[0].(map[string]interface{}); ok {
							if add, ok := item["_additional"].(map[string]interface{}); ok {
								switch v := add["vector"].(type) {
								case []interface{}:
									expectedDim = len(v)
								case []float64:
									expectedDim = len(v)
								}
							}
						}
					}
				}
			}
			if expectedDim > 0 && len(queryVector) != expectedDim {
				resized := make([]float32, expectedDim)
				copy(resized, queryVector)
				queryVector = resized
			}
		}

		if queryVector == nil {
			log.Printf("SemanticQueryHandler: no query vector available; skipping Weaviate semantic search")
		} else {
			nearVector := weaviateClient.GraphQL().NearVectorArgBuilder().WithVector(queryVector)

			additional := graphql.Field{Name: "_additional", Fields: []graphql.Field{{Name: "certainty"}, {Name: "distance"}, {Name: "id"}}}
			var fields []graphql.Field
			switch classToUse {
			case "CHISGElement":
				fields = []graphql.Field{{Name: "name"}, {Name: "description"}, {Name: "domain"}, {Name: "layer"}, {Name: "chisg_id"}, {Name: "source"}, {Name: "suggested_years"}, additional}
			case "CHISGSkill":
				fields = []graphql.Field{{Name: "name"}, {Name: "description"}, {Name: "domain"}, {Name: "layer"}, {Name: "chisgId"}, additional}
			case "EducationalSkills":
				fields = []graphql.Field{{Name: "name"}, {Name: "description"}, {Name: "skill_type"}, {Name: "development_age"}, {Name: "mongo_id"}, additional}
			case "SemanticLinks":
				fields = []graphql.Field{{Name: "statement"}, {Name: "source_term"}, {Name: "target_term"}, {Name: "forward_relation"}, {Name: "relation_type"}, {Name: "context"}, {Name: "confidence"}, {Name: "domain"}, additional}
			default:
				fields = []graphql.Field{{Name: "name"}, {Name: "description"}, additional}
			}

			resp, werr := weaviateClient.GraphQL().Get().
				WithClassName(classToUse).
				WithFields(fields...).
				WithNearVector(nearVector).
				WithLimit(limit).
				Do(ctx)
			if werr != nil {
				log.Printf("SemanticQueryHandler: Weaviate query failed: %v", werr)
			} else if resp != nil && resp.Errors != nil {
				log.Printf("SemanticQueryHandler: Weaviate returned errors: %v", resp.Errors)
			} else {
				getBlock, ok := resp.Data["Get"].(map[string]interface{})
				if ok {
					itemsAny, ok := getBlock[classToUse].([]interface{})
					if ok {
						for _, raw := range itemsAny {
							item, ok := raw.(map[string]interface{})
							if !ok {
								continue
							}

							// Map Weaviate object -> frontend semantic match shape
							match := map[string]interface{}{}
							var name, desc, skillType string

							switch classToUse {
							case "CHISGElement":
								name, _ = item["name"].(string)
								desc, _ = item["description"].(string)
								domain, _ := item["domain"].(string)
								layer, _ := item["layer"].(string)
								skillType = strings.TrimSpace(strings.Trim(strings.Join([]string{domain, layer}, " • "), " • "))
							case "CHISGSkill":
								name, _ = item["name"].(string)
								desc, _ = item["description"].(string)
								domain, _ := item["domain"].(string)
								layer, _ := item["layer"].(string)
								skillType = strings.TrimSpace(strings.Trim(strings.Join([]string{domain, layer}, " • "), " • "))
							case "EducationalSkills":
								name, _ = item["name"].(string)
								desc, _ = item["description"].(string)
								skillType, _ = item["skill_type"].(string)
							case "SemanticLinks":
								stmt, _ := item["statement"].(string)
								sourceTerm, _ := item["source_term"].(string)
								targetTerm, _ := item["target_term"].(string)
								rel, _ := item["forward_relation"].(string)
								if rel == "" {
									rel, _ = item["relation_type"].(string)
								}
								if stmt != "" {
									name = stmt
								} else {
									name = strings.TrimSpace(strings.Join([]string{sourceTerm, rel, targetTerm}, " "))
								}
								desc, _ = item["context"].(string)
								skillType, _ = item["domain"].(string)
							default:
								name, _ = item["name"].(string)
								desc, _ = item["description"].(string)
							}

							match["skill_name"] = name
							match["description"] = desc
							match["skill_type"] = skillType

							// certainty preferred; fall back to 1-distance.
							semanticRelevance := 0.0
							if add, ok := item["_additional"].(map[string]interface{}); ok {
								if c, ok := add["certainty"].(float64); ok {
									semanticRelevance = c
								} else if d, ok := add["distance"].(float64); ok {
									semanticRelevance = 1.0 - d
								}
							}
							if semanticRelevance < 0 {
								semanticRelevance = 0
							}
							if semanticRelevance > 1 {
								semanticRelevance = 1
							}
							match["semantic_relevance"] = semanticRelevance

							semanticMatches = append(semanticMatches, match)
						}
					}
				}
			}
		}
	} else {
		log.Printf("SemanticQueryHandler: no suitable Weaviate class found for semantic search")
	}

	// --- Response mapping to frontend contract ---
	mapSkill := func(s models.Skill) map[string]interface{} {
		criteria := make([]string, 0, len(s.SkillCriteria))
		for _, c := range s.SkillCriteria {
			criteria = append(criteria, fmt.Sprintf("Level %d: %s", c.Level, c.Description))
		}

		out := map[string]interface{}{
			"skill_name":      s.Name,
			"description":     s.Description,
			"skill_type":      s.Category,
			"development_age": s.DevelopmentAge,
			"criteria_levels": len(criteria),
			"criteria":        criteria,
		}

		if s.SourceTitle != "" || s.SourceType != "" || s.ExtractionMethod != "" || s.ChapterTitle != "" {
			out["source_info"] = map[string]interface{}{
				"title":             s.SourceTitle,
				"type":              s.SourceType,
				"chapter_title":     s.ChapterTitle,
				"extraction_method": s.ExtractionMethod,
			}
		}

		return out
	}

	exactOut := make([]map[string]interface{}, 0, len(exactSkills))
	for _, s := range exactSkills {
		exactOut = append(exactOut, mapSkill(s))
	}

	relatedOut := make([]map[string]interface{}, 0, len(relatedSkills))
	for _, s := range relatedSkills {
		relatedOut = append(relatedOut, mapSkill(s))
	}

	// --- Basic insights (no hand-wavy LLM claims) ---
	wordCount := len(strings.Fields(queryText))
	complexity := "medium"
	if wordCount <= 2 {
		complexity = "low"
	} else if wordCount >= 7 {
		complexity = "high"
	}

	high, medium, low := 0, 0, 0
	for _, m := range semanticMatches {
		v, _ := m["semantic_relevance"].(float64)
		switch {
		case v >= 0.75:
			high++
		case v >= 0.5:
			medium++
		default:
			low++
		}
	}

	synthesis := fmt.Sprintf(
		"Found %d exact match(es) in MongoDB and %d semantic match(es) in Weaviate%s.",
		len(exactOut),
		len(semanticMatches),
		func() string {
			if classToUse == "" {
				return ""
			}
			return fmt.Sprintf(" (class: %s)", classToUse)
		}(),
	)

	if len(semanticMatches) == 0 && len(exactOut) == 0 {
		synthesis = "No matches found. Try a broader skill phrase (e.g., 'working memory' or 'executive function')."
	}

	response := map[string]interface{}{
		"query":     queryText,
		"timestamp": time.Now(),
		"results_summary": map[string]int{
			"total_exact_matches":    len(exactOut),
			"total_semantic_matches": len(semanticMatches),
			"total_related_skills":   len(relatedOut),
		},
		"exact_matches":    exactOut,
		"semantic_matches": semanticMatches,
		"related_skills":   relatedOut,
		"llm_synthesis":    synthesis,
		"semantic_insights": map[string]interface{}{
			"query_complexity": complexity,
			"confidence_levels": map[string]int{
				"high_confidence":   high,
				"medium_confidence": medium,
				"low_confidence":    low,
			},
		},
	}

	// If the handler is being called but Weaviate isn't reachable/configured, surface it clearly.
	if weaviateClient == nil {
		response["semantic_error"] = "Weaviate client not available (check WEAVIATE_URL and container health)"
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// At this point headers are likely written; just log.
		log.Printf("SemanticQueryHandler: failed to encode response: %v", err)
		return
	}
}
