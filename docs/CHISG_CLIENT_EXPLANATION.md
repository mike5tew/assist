# CHISG Go Client - Architecture & Implementation Guide

## What is the CHISG Go Client?

The CHISG Go Client (`/backend/internal/integration/chisg_client.go`) is a **bridge** between the HumanOS AI system and the CHISG semantic knowledge graph. It's essentially a Go wrapper that makes it easy and safe to query CHISG for structured knowledge.

### Simple Analogy
Think of it like this:
- **CHISG** = A library (the knowledge source)
- **Go Client** = A librarian (helps you find books efficiently)
- **HumanOS** = A student (asks questions)

The client handles all the complexity of talking to CHISG, so HumanOS just needs to ask simple questions.

---

## What Does It Actually Do?

### 1. **Handles Connection & Communication**
```
User Request
    ↓
Go Client (formats request)
    ↓
HTTP Connection Pool
    ↓
CHISG Service
    ↓
Go Client (parses response)
    ↓
Structured Response
```

The client:
- Opens HTTP connections efficiently (connection pooling)
- Formats your question into CHISG's expected format
- Sends the request with proper headers
- Receives the response and parses it back into Go objects

### 2. **Makes Requests Fast (Caching)**
```
Query: "What is photosynthesis?"
    ↓
Check Cache (is this in memory?)
    ├─ YES → Return immediately (no network call!)
    └─ NO → Query CHISG → Store in cache → Return
```

**Why cache?** Because the same questions get asked frequently. Caching saves:
- Time (no network round trip)
- Money (fewer API calls to CHISG)
- Bandwidth

Cache expires after 5 minutes, so you get fresh data regularly.

### 3. **Handles Failures Gracefully (Retry Logic)**
```
Attempt 1: "CHISG is down" 
    ↓ Wait 100ms
Attempt 2: "CHISG is down"
    ↓ Wait 200ms  
Attempt 3: "CHISG is down"
    ↓ Wait 400ms
Attempt 4: FAIL (after 3 retries)
```

**Why retry?** Network hiccups happen. The client automatically retries with exponential backoff:
- 1st retry: Wait 100 milliseconds
- 2nd retry: Wait 200 milliseconds
- 3rd retry: Wait 400 milliseconds
- If all 3 fail: Tell user there's an error

This makes the system more reliable without slowing things down.

### 4. **Monitors Performance (Latency Tracking)**
```
Query Time: 150ms ✅ Good!
Query Time: 250ms ⚠️ Slow (target: <200ms)
```

The client logs how long each query takes. If it's slow, you know there's a problem.

---

## Core Components

### The Query Struct (What You Send)
```go
type KnowledgeQuery struct {
	Topic          string   // e.g., "photosynthesis"
	Concepts       []string // e.g., ["chlorophyll", "glucose", "light"]
	TargetAudience string   // e.g., "gcse_student" or "medical_student"
}
```

**What this means:**
- **Topic**: The main subject you're asking about
- **Concepts**: Related ideas you want included in the response
- **TargetAudience**: Who is asking (so CHISG can adjust complexity)

### The Response Struct (What You Get Back)
```go
type KnowledgeResponse struct {
	Summary         string            // Plain English explanation
	KeyConcepts     []string          // Main ideas to learn
	Prerequisites   []string          // What you need to know first
	Analogies       []string          // Simple comparisons
	ConfidenceScore float64           // How sure is CHISG? (0.0 to 1.0)
	RelatedTopics   map[string]string // Links to related topics
	Timestamp       time.Time         // When this was generated
}
```

**What this means:**
- **Summary**: "Photosynthesis is how plants turn sunlight into food"
- **KeyConcepts**: ["chlorophyll", "glucose", "oxygen"]
- **Prerequisites**: ["what is light", "what are atoms"]
- **Analogies**: ["like a solar panel converting sun to electricity"]
- **ConfidenceScore**: 0.95 (95% sure this is right)
- **RelatedTopics**: {"cellular_respiration": "the reverse process", ...}

---

## How It Works: The Complete Flow

### Simple Example: A Student Asking About Photosynthesis

**Step 1: Student Asks HumanOS**
```
"What is photosynthesis?"
```

**Step 2: HumanOS Creates a Query**
```go
query := &KnowledgeQuery{
	Topic:          "photosynthesis",
	Concepts:       []string{"chlorophyll", "light", "glucose"},
	TargetAudience: "gcse_student",  // Age 14-16
}
```

**Step 3: HumanOS Calls the Client**
```go
response, err := chisgClient.Query(ctx, query)
```

**Step 4: Client Does Everything**
```
a) Check cache
   "Is photosynthesis already in memory?"
   → NO, so continue

b) Format the request
   Convert query to JSON

c) Connect to CHISG
   Use connection pool (fast!)

d) Send request with retries
   Try up to 3 times if network fails

e) Parse response
   Convert JSON back to Go structs

f) Cache the result
   Store in memory for 5 minutes

g) Return to HumanOS
```

**Step 5: HumanOS Gets Back Structured Knowledge**
```json
{
  "summary": "Photosynthesis is the process where plants convert sunlight, water, and CO2 into glucose and oxygen",
  "key_concepts": ["chlorophyll", "light-dependent reaction", "Calvin cycle", "ATP", "NADPH"],
  "prerequisites": ["What is a plant cell", "What is ATP", "What is the chloroplast"],
  "analogies": ["Like a solar panel turning sunlight into electricity"],
  "confidence_score": 0.98,
  "related_topics": {
    "cellular_respiration": "The reverse process - using glucose to make energy",
    "light_absorption": "How chlorophyll captures light energy"
  }
}
```

**Step 6: HumanOS Adjusts for Age & Barriers**
```
Apply HumanOS barrier detection:
- Is this student struggling with basic concepts?
- Is this age-appropriate for a 12-year-old?
- Are there triggers (trauma, learning disabilities)?

Then adjust the response accordingly:
- Simplify language if needed
- Add more analogies if confused
- Skip prerequisites already known
- Add supportive intervention if traumatized
```

---

## Configuration (Settings You Control)

### Via Environment Variables (.env)
```bash
CHISG_BASE_URL=http://localhost:8080    # Where CHISG lives
CHISG_API_KEY=dev-key-12345             # Authentication key
```

### Via Code Configuration
```go
client := NewCHISGClient("http://localhost:8080")

// Customize behavior
client.SetCacheTTL(10 * time.Minute)        // Cache for 10 mins instead of 5
client.SetMaxRetries(5)                      // Retry 5 times instead of 3
client.SetInitialBackoff(200 * time.Millisecond) // Start waiting 200ms

// Check health
if err := client.Health(context.Background()); err != nil {
    log.Fatal("CHISG is not responding!")
}

// Get statistics
stats := client.GetStats()
log.Printf("Cached items: %v", stats["cached_items"])
```

---

## What Needs to Be Built Next

### 1. **Shared Types File** (`/backend/internal/integration/chisg_types.go`)
**Purpose**: Define bridge types between HumanOS and CHISG

**What it should contain:**
```go
// Request type that HumanOS barrier detection can create
type HumanOSKnowledgeRequest struct {
	Query             string   // User's actual question
	Age               int      // Student's age
	DetectedBarriers  []string // ["trauma", "confusion", "low_confidence"]
	TriggerWords      []string // Words that might upset student
	PrerequisitesKnown []string // Topics already mastered
}

// Response wrapper with HumanOS metadata
type HumanOSEnrichedResponse struct {
	CHISGKnowledge    *KnowledgeResponse // Original CHISG response
	AgeAdjusted       string             // Simplified for age
	InterventionNeeded bool              // Does student need support?
	SupportScript     string             // What to say to help
	ConfidenceAfter   float64            // Confidence score after adjustments
}
```

### 2. **Coach Respond Handler** (`/backend/internal/integration/coach_handler.go`)
**Purpose**: Orchestrate the full flow: HumanOS → CHISG → Response

**What it should do:**
```go
func CoachRespondHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse user message
	var req struct {
		Message string `json:"message"`
		Age     int    `json:"age"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// 2. Detect barriers (HumanOS)
	barriers := detectBarriers(req.Message, req.Age)

	// 3. Create CHISG query
	query := &KnowledgeQuery{
		Topic:          extractTopic(req.Message),
		Concepts:       extractConcepts(req.Message),
		TargetAudience: ageToAudience(req.Age),
	}

	// 4. Query CHISG via client
	chisgResponse, err := chisgClient.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Knowledge lookup failed", 500)
		return
	}

	// 5. Adjust for age & barriers
	adjusted := adjustForAge(chisgResponse, req.Age)
	adjusted = adjustForBarriers(adjusted, barriers)

	// 6. Return combined response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adjusted)
}
```

**Flow diagram:**
```
User Message
    ↓
[Barrier Detection] ← HumanOS
    ↓
[Topic Extraction]
    ↓
[CHISG Query] → CHISG Client → CHISG Service
    ↓
[Parse Response]
    ↓
[Age Adjustment] ← HumanOS
    ↓
[Barrier Adjustment] ← HumanOS
    ↓
[Return to User]
```

### 3. **Integration Tests** (`/backend/tests/integration_chisg_test.go`)
**Purpose**: Verify the whole system works together

**What tests should cover:**
```go
func TestCHISGClientConnection(t *testing.T)
	// Does client connect to CHISG?

func TestCHISGQueryLatency(t *testing.T)
	// Is response < 200ms?

func TestCoachRespondFlow(t *testing.T)
	// Does full pipeline work?
	// User → HumanOS → CHISG → Response

func TestAgeAppropriateResponse(t *testing.T)
	// Does age 8 get simpler response than age 16?

func TestCacheEffectiveness(t *testing.T)
	// Does cache actually speed things up?

func TestErrorHandling(t *testing.T)
	// What happens if CHISG is down?
	// Do retries work?
```

---

## Request/Response Examples

### Example 1: GCSE Student Learning Photosynthesis

**Request:**
```json
{
  "message": "I don't understand photosynthesis",
  "age": 14
}
```

**Internal Processing:**
1. Detect barriers: ["confusion", "low_confidence"]
2. Extract topic: "photosynthesis"
3. Create CHISG query with audience "gcse_student"

**CHISG Response:**
```json
{
  "summary": "Plants use sunlight to convert water and carbon dioxide into glucose (sugar) and oxygen",
  "key_concepts": ["chlorophyll", "light reaction", "dark reaction", "glucose"],
  "prerequisites": ["photosynthesis location (chloroplast)", "what is ATP"],
  "analogies": ["Like a solar panel converting sun to electricity"],
  "confidence_score": 0.95,
  "related_topics": {
    "cellular_respiration": "The reverse - plants use glucose for energy",
    "light_wavelengths": "Different colors of light absorbed"
  }
}
```

**After HumanOS Adjustment:**
```json
{
  "original_response": { ... CHISG response ... },
  "adjusted_summary": "Plants are like little solar panels. They use sunlight to turn water and air into food (sugar) and release oxygen.",
  "barriers_detected": ["confusion"],
  "support_offered": "That's a complex topic! Let me break it into steps...",
  "next_steps": ["Learn what a chloroplast is", "Learn what ATP is"],
  "confidence_after_adjustment": 0.92
}
```

### Example 2: Medical Student Learning Immunology

**Request:**
```json
{
  "message": "Explain X-linked agammaglobulinemia",
  "age": 22
}
```

**Internal Processing:**
1. Detect barriers: [] (none - advanced student)
2. Extract topic: "X-linked agammaglobulinemia"
3. Create CHISG query with audience "medical_student"

**CHISG Response:**
```json
{
  "summary": "X-linked agammaglobulinemia is a primary immunodeficiency caused by mutations in the BTK gene, resulting in absent B cells and severe hypogammaglobulinemia",
  "key_concepts": ["B cell development", "BTK protein", "immunoglobulin deficiency", "Bruton's tyrosine kinase"],
  "prerequisites": ["B cell biology", "immunoglobulin structure", "genetic inheritance patterns"],
  "analogies": ["B cell production line breaks down - no antibodies manufactured"],
  "confidence_score": 0.97,
  "related_topics": {
    "other primary immunodeficiencies": "Similar genetic conditions",
    "immunoglobulin replacement therapy": "Standard treatment approach"
  }
}
```

**After HumanOS Adjustment:**
```json
{
  "original_response": { ... CHISG response ... },
  "barriers_detected": [],
  "confidence_after_adjustment": 0.97,
  "notes": "No adjustments needed - response appropriate for audience"
}
```

---

## Error Scenarios & Handling

### Scenario 1: CHISG is Down

**Request:**
```
Query: "What is photosynthesis?"
```

**Client Behavior:**
```
Attempt 1 → Connection refused → Wait 100ms
Attempt 2 → Connection refused → Wait 200ms
Attempt 3 → Connection refused → Wait 400ms
All retries failed → Return error to HumanOS
```

**Response to User:**
```json
{
  "error": "Knowledge service temporarily unavailable",
  "fallback": "I'm having trouble accessing my knowledge base right now. Try again in a moment.",
  "recovery_time": "~1 minute"
}
```

### Scenario 2: Topic Not Found

**Request:**
```
Query: "Explain quantum superposition with medieval tapestries"
```

**CHISG Response:**
```json
{
  "summary": "I don't have information linking quantum superposition to medieval tapestries",
  "confidence_score": 0.1,
  "suggestions": [
    "Try asking about quantum superposition alone",
    "Try asking about historical parallels separately"
  ]
}
```

---

## Performance Expectations

| Operation | Expected Time | Status |
|-----------|---------------|--------|
| Cache hit (already asked) | 1-5ms | ✅ Excellent |
| CHISG query (network + processing) | 50-150ms | ✅ Good |
| Total request (barrier detection + CHISG + adjustment) | <200ms | ✅ Target |
| Retry (with backoff) | 500-1500ms | ✅ Acceptable |

---

## Security Considerations

### 1. **API Key Protection**
```bash
# GOOD ✅
export CHISG_API_KEY="prod-key-xyz"  # Environment variable
# Key never appears in code

# BAD ❌
const CHISG_API_KEY = "prod-key-xyz"  # Hardcoded
# Key visible in source code
```

### 2. **Request Validation**
```go
// Always validate input before sending to CHISG
if len(query.Topic) == 0 {
    return nil, fmt.Errorf("topic cannot be empty")
}

if query.ConfidenceScore < 0 || query.ConfidenceScore > 1 {
    return nil, fmt.Errorf("confidence score must be 0.0-1.0")
}
```

### 3. **Response Sanitization**
```go
// Don't just return CHISG response as-is
// Filter for appropriate content before showing student
if containsInappropriateContent(response) {
    return nil, fmt.Errorf("response filtered for age appropriateness")
}
```

---

## When Everything Works

**The Happy Path:**
```
✅ Student asks question
✅ HumanOS detects barriers (if any)
✅ Client queries CHISG (< 100ms if cached, < 200ms if fresh)
✅ Response adjusted for age & barriers
✅ Student gets clear, helpful answer
✅ Response cached for next similar question
```

**Result:** Students get intelligent, personalized, trauma-aware tutoring at scale! 🎉

