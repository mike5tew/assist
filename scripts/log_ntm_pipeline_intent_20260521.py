"""
Log the NTM full-pipeline intent and CHISG review architecture to Weaviate and MongoDB.
Run once to capture the design decisions made on 2026-05-21.
"""
import weaviate
import datetime
import sys
sys.path.insert(0, '/Users/michaelstewart/Coding/humanOS/scripts')

from pymongo import MongoClient

WEAVIATE_HOST = "localhost"
WEAVIATE_HTTP_PORT = 8088
WEAVIATE_GRPC_PORT = 50052

MONGO_URI = "mongodb://admin:password123@localhost:27018"
MONGO_DB = "esp_organizer"

ENTRIES = [
    {
        "title": "CHISG Review Pipeline — Full Architecture (chisg_review_handlers.go)",
        "content": """The CHISG review pipeline is the backbone of the research-team dataset-building workflow.
It is implemented in esp-organizer/internal/domain/api/chisg_review_handlers.go.

PIPELINE FLOW:
1. UploadCHISGDocumentHandler — accepts multipart PDF or text file.
   - If PDF: calls AWS Textract (eu-west-2, bucket esp-new-organizer-immunology) to extract raw text.
   - Cleans and chunks the text.
   - For each chunk of ≥100 chars: calls LlamaClient (Claude via Bedrock) with a strict extraction prompt.
   - Prompt enforces 13-relation controlled vocabulary:
     causes|leads_to|part_of|contains|develops_into|regulates|enables|inhibits|treats|diagnoses|manifests_as|is_a|located_at
   - Each extracted triple has: entity_a, entity_b, relation, source_quote (verbatim), context, attribution (cited source if applicable).
   - Results stored immediately to MongoDB chisg_knowledge_base.review_tasks (status: pending_review).
   - Returns job_id to client instantly so reviewers can poll while extraction continues.

2. GetPendingReviewTasksHandler — GET, returns up to 10 pending tasks from review_tasks collection.

3. ApproveReviewTaskHandler — POST /{id}/approve
   - Accepts approved_links array (edited by human reviewer).
   - Updates task status to "approved", stores final links.
   - Captures provenance: source_document_id, source_title, DOI.
   - *** TODO (NOT YET IMPLEMENTED): Push approved links to Weaviate SemanticLinks class ***
   - This Weaviate write-back is the critical loop that feeds the learning algorithm.

4. ClearReviewTasksHandler — clears the review queue.

MONGODB SCHEMA (chisg_knowledge_base.review_tasks):
- paper_id, source_document_id, source_title, doi, chunk_id, source_text
- proposed_links: [{entity_a, entity_b, relation, source_quote, context, attribution}]
- status: pending_review | approved | rejected
- assigned_to, reviewed_by, reviewed_at, created_at

MISSING PIECES (as of 2026-05-21):
1. ApproveReviewTaskHandler Weaviate write-back (TODO in code)
2. Frontend review UI — no page exists yet for researchers to review tasks
3. NTM upload not yet wired to this pipeline (uses simpler NTMUploadHandler instead)
""",
        "file_path": "esp-organizer/internal/domain/api/chisg_review_handlers.go",
        "section_path": "CHISG Review Pipeline > Architecture",
        "source_type": "infrastructure",
    },
    {
        "title": "NTM Research Page — Full Pipeline Intent (2026-05-21)",
        "content": """The NTM Research page (/ntm) was created 2026-05-20 as the research interface for the
Michael McGrath / University of Brighton NTM mycobacteria collaboration.

INTENT: This page is for research teams to:
1. Upload academic papers (PDF) about NTM/mycobacteria
2. The system extracts semantic triples (entity_a, relation, entity_b) with source_quote provenance
3. A researcher reviews, edits, and approves the extracted links
4. Approved links are fed back into Weaviate (AcademicLink class) and the CHISG knowledge graph
5. The graph grows iteratively — each approved paper enriches the knowledge base used for AI queries
6. Researcher corrections (edits to relations, contexts, entities) feed back to improve extraction quality over time

CURRENT STATUS (as of 2026-05-21):
- Frontend page exists at /ntm with Upload tab and Query tab (NTMResearch.tsx)
- NTMUploadHandler (ntm_handlers.go) only does file storage — NOT the full pipeline
- NTMQueryHandler uses Weaviate BM25 on AcademicLink class + Claude for synthesis
- The full CHISG pipeline (chisg_review_handlers.go) is NOT yet wired to the NTM page

WHAT NEEDS TO BE BUILT:
1. Wire NTM upload to UploadCHISGDocumentHandler (or a variant scoped to AcademicLink)
2. Add a "Review" tab to NTMResearch.tsx:
   - Poll GetPendingReviewTasksHandler
   - Show each chunk's source_text alongside proposed triples
   - Researcher can edit entity names, relations, context inline
   - Approve sends to ApproveReviewTaskHandler
3. Implement the Weaviate write-back in ApproveReviewTaskHandler:
   - On approval, insert link into Weaviate AcademicLink class
   - This makes it immediately queryable by NTMQueryHandler
4. Log corrections for future model fine-tuning (store original vs edited triple)
""",
        "file_path": "frontend/src/pages/CHISG/NTMResearch.tsx",
        "section_path": "NTM Research > Pipeline Intent",
        "source_type": "feature",
    },
    {
        "title": "Weaviate Write-back — ApproveReviewTaskHandler TODO",
        "content": """The Weaviate write-back step is the critical missing piece in the CHISG review pipeline.

LOCATION: esp-organizer/internal/domain/api/chisg_review_handlers.go — ApproveReviewTaskHandler

WHAT NEEDS TO HAPPEN after human approves a task:
1. For each approved SemanticLink, insert into Weaviate with:
   - entity_a, relation, entity_b, context, source_quote, attribution
   - paper_id (from task.PaperID), chunk_id (from task.ChunkID)
   - source_document_id (from task.SourceDocumentID.Hex()), source_title, doi
2. Use weaviate-go-client v4 (module: github.com/weaviate/weaviate-go-client/v4 v4.11.0 via replace)
3. Client: use dbstore.GetWeaviateClient() — same pattern as searchAcademicLinks in ntm_handlers.go
4. Class: AcademicLink (already has the correct schema from ingestion)

CORRECTION LOGGING (for ML feedback loop):
- Store original_link (Claude's first extraction) vs approved_link (human-edited) in a separate
  MongoDB collection: chisg_knowledge_base.extraction_corrections
- This dataset can be used to fine-tune the extraction prompt or train a correction model
- Fields: task_id, chunk_text, original_triple, approved_triple, correction_type (relation_changed|entity_renamed|context_added|rejected), reviewed_by, reviewed_at

TARGET: When fully implemented, every approved paper adds to the shared knowledge graph
immediately, making it available to all NTM research queries without redeployment.
""",
        "file_path": "esp-organizer/internal/domain/api/chisg_review_handlers.go",
        "section_path": "CHISG Review Pipeline > Weaviate Write-back",
        "source_type": "infrastructure",
    },
]


def log_to_weaviate(client, entry):
    collection = client.collections.get("Documentation")
    uuid = collection.data.insert({
        "title": entry["title"],
        "content": entry["content"],
        "file_path": entry["file_path"],
        "section_path": entry["section_path"],
        "project": "assist",
        "source_type": entry["source_type"],
    })
    return uuid


def log_to_mongo(entry):
    client = MongoClient(MONGO_URI)
    db = client[MONGO_DB]
    coll = db["project_notes"]
    result = coll.insert_one({
        "title": entry["title"],
        "content": entry["content"],
        "file_path": entry["file_path"],
        "section_path": entry["section_path"],
        "project": "assist",
        "source_type": entry["source_type"],
        "created_at": datetime.datetime.utcnow(),
    })
    client.close()
    return result.inserted_id


def main():
    print("Connecting to Weaviate...")
    weaviate_client = weaviate.connect_to_local(
        host=WEAVIATE_HOST,
        port=WEAVIATE_HTTP_PORT,
        grpc_port=WEAVIATE_GRPC_PORT,
    )

    for entry in ENTRIES:
        print(f"\nLogging: {entry['title']}")
        wuid = log_to_weaviate(weaviate_client, entry)
        print(f"  Weaviate UUID: {wuid}")
        mid = log_to_mongo(entry)
        print(f"  MongoDB ID:    {mid}")

    weaviate_client.close()
    print("\nDone — 3 entries logged.")


if __name__ == "__main__":
    main()
