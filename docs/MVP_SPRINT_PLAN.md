# AI Tutor - MVP Sprint Plan

**Objective**: Build a demonstrable AI Tutor prototype within 7 days that proves the core value proposition of the integrated learning OS.

**Guiding Principle**: "Mock what you can, build what you must." We will prioritize the backend logic that shows intelligence and use a minimal UI.

---

## The Demo Story

This is what we will be able to show at the end of the sprint:

1.  **User Input**: A user types a question like "What causes XLA?" into a simple web interface.
2.  **Intelligent Processing**: The Go backend receives the query. The `CoachRespondHandler` orchestrates the call.
3.  **Knowledge Retrieval**: The system queries the Weaviate knowledge graph to find relevant, interconnected `SemanticLinks`.
4.  **Psychological Adjustment**: The `HumanOS` logic layer adjusts the response for the target audience (e.g., simplifies it for a GCSE student).
5.  **Structured Output**: The UI displays a clear, context-aware answer that is factually correct and psychologically attuned.

This demo proves the end-to-end flow and the "magic" of your integrated system.

---

## 7-Day Critical Path

### Day 1: Seed the Brain (Data Foundation)
**Goal**: Have a small, perfect, queryable knowledge graph in Weaviate.

-   [ ] **Task**: Create a new script `tools/seed-demo-data/main.go`.
-   [ ] **Logic**: This script will manually create and store 10-15 high-quality, interconnected `SemanticLinks` for a single topic (e.g., "XLA").
-   [ ] **Details**: Hardcode the links, including `source_term`, `target_term`, `relation_type`, and `context`. This bypasses the entire complex extraction pipeline for the demo.
-   **Outcome**: A reliable, predictable data source for the backend to query.

### Day 2: Build the Query Service (The Librarian)
**Goal**: A service that can traverse the seeded knowledge graph.

-   [ ] **Task**: Implement the core logic in `esp-organizer/internal/domain/hsg_query_service.go`.
-   [ ] **Function**: `TraverseHierarchy(rootLinkID string)`.
-   [ ] **Logic**: Given a starting link, follow the `is_child_of` references to build a small, contextual graph of related information. Keep it simple: a 2-level traversal is enough for the demo.
-   **Outcome**: A function that can turn a single query result into a rich, contextual story.

### Day 3: The Conductor (Coach Handler)
**Goal**: An API endpoint that orchestrates the entire response.

-   [ ] **Task**: Implement the `CoachRespondHandler` in `esp-organizer/internal/domain/api/coach_handler.go`.
-   [ ] **Endpoint**: `POST /api/coach/respond`.
-   [ ] **Logic**:
    1.  Receive a user query.
    2.  Perform a simple vector search in Weaviate to find the best starting `SemanticLink`.
    3.  Call `hsg_query_service.TraverseHierarchy()` with the result.
    4.  (Mocked) Apply a simple age-appropriateness filter (e.g., `if age < 16, simplify language`).
    5.  Return the structured response as JSON.
-   **Outcome**: A working API endpoint that demonstrates the core intelligence.

### Day 4: The Face (Minimal UI)
**Goal**: A simple interface to interact with the backend.

-   [ ] **Task**: Modify `frontend/src/components/AIChat.tsx`.
-   [ ] **UI**: A single text input box and a "Submit" button.
-   [ ] **Display**: A pre-formatted `<div>` to render the JSON response from `/api/coach/respond`.
-   [ ] **Focus**: Make it functional, not beautiful. No need for chat bubbles or complex styling.
-   **Outcome**: A way to visually demonstrate the end-to-end flow.

### Day 5: Integration & Testing
**Goal**: Ensure all pieces work together smoothly.

-   [ ] **Task**: Run the entire system locally.
-   [ ] **Flow**: Use the UI to send a query, watch the backend logs, and see the response rendered.
-   [ ] **Debug**: Fix CORS issues, API endpoint mismatches, and data formatting errors.
-   **Outcome**: A stable, repeatable demo.

### Day 6: Record the Demo & Refine
**Goal**: Create a perfect, recorded video of the demo.

-   [ ] **Task**: Use screen recording software (like QuickTime or Loom) to capture the demo story in action.
-   [ ] **Narrate**: Explain what's happening at each step ("Now the system is querying the knowledge graph... now it's applying the psychological framework...").
-   [ ] **Iterate**: Record it multiple times until you have a smooth, compelling 2-minute video.
-   **Outcome**: A polished asset you can show anyone, anytime, without live-demo risks.

### Day 7: Rest & Prepare
**Goal**: Recharge and prepare to talk about the vision, not the code.

-   [ ] **Task**: Step away from the keyboard. Review your `INTERVIEW_STRATEGY.md` and `STRATEGIC_PLAYBOOK.md`.
-   [ ] **Focus**: Your value is the vision and the architecture. The demo is just proof you can execute.
-   **Outcome**: You enter your conversations confident and rested, ready to discuss strategy.

---

## What We Are NOT Building (De-Scoped for MVP)

-   The full PDF extraction pipeline.
-   User authentication.
-   Real-time updates or websockets.
-   A beautiful chat UI.
-   The `Misunderstanding Map` or `Federated Learning`.
-   Payment infrastructure.

This sprint is about building the **illusion of the complete system** by focusing on the most critical, intelligence-demonstrating components. You've got this. Let's build.
