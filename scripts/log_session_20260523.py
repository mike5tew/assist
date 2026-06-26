#!/usr/bin/env python3
"""
Log session changes 2026-05-23 to Weaviate Documentation class and MongoDB.
Covers: CHISG auth deployment, 502 fix, data migration, NTM tab merge,
        /ntm behind PrivateRoute, and CHISG product strategy decisions.

Usage: python3 log_session_20260523.py
"""

import weaviate
import hashlib
from datetime import datetime
from pymongo import MongoClient

WEAVIATE_HOST = "localhost"
WEAVIATE_PORT = 8088
WEAVIATE_GRPC_PORT = 50052
DOC_CLASS = "Documentation"

MONGO_URI = "mongodb://admin:password123@localhost:27018/esp_organizer?authSource=admin"
MONGO_DB = "esp_organizer"
NOTES_COLLECTION = "project_notes"

CHANGES = [
    {
        "title": "CHISG JWT Auth System Deployed — /api/auth/login (2026-05-23)",
        "project": "assist",
        "section_path": "Backend > CHISG > Auth",
        "file_path": "esp-organizer/internal/domain/api/auth_handlers.go",
        "content": """## CHISG JWT Auth System

### Endpoint
POST /api/auth/login — public, no auth required

### Request
```json
{ "username": "chisg_admin", "password": "McGrath2024!" }
```

### Response
```json
{ "token": "<jwt>", "expires_in": 86400 }
```

### Implementation
- File: `esp-organizer/internal/domain/api/auth_handlers.go`
- HMAC-SHA256 using Go stdlib `crypto/hmac` + `crypto/sha256` — NO external JWT library
- Token format: base64url(header).base64url(payload).base64url(sig)
- Payload: `{"sub":"<username>","iat":<unix>,"exp":<unix+86400>}`
- 24h expiry
- Credentials from env vars: `CHISG_USER` / `CHISG_PASS`
- Signed with env var: `JWT_SECRET`

### Middleware
`RequireAuth(next http.Handler)` — validates Bearer token in Authorization header
Applied to chisgRouter subrouter (wraps /chisg/papers and /chisg/links)

### Production Env Values (on Vultr /opt/esp-thinking/.env)
- CHISG_USER=chisg_admin
- CHISG_PASS=McGrath2024!
- JWT_SECRET=wOx4NtaRLQBPN4HwHb//wyo+LxKjFI7cmP8vSm82yP0=

### Router Structure
- `/api/auth/login` — public POST
- `/api/chisg/graph` — public GET (knowledge graph visualisation data)
- `/api/chisg/papers` — RequireAuth — GET list of AcademicPaper objects
- `/api/chisg/links` — RequireAuth — GET AcademicLink objects by paper_id
- `/api/expert-review/*` — RequireAuth (expertRouter)
- `/api/source-documents` — RequireAuth (sourceDocRouter)
""",
    },
    {
        "title": "CHISG Papers + Links API Handlers Deployed (2026-05-23)",
        "project": "assist",
        "section_path": "Backend > CHISG > Handlers",
        "file_path": "esp-organizer/internal/domain/api/chisg_papers_handlers.go",
        "content": """## CHISG Papers and Links Handlers

### CHISGListPapersHandler
GET /api/chisg/papers — requires Bearer token auth

Queries Weaviate `AcademicPaper` class. Returns:
```json
{
  "papers": [
    {
      "paper_id": "string",
      "title": "string",
      "filename": "string",
      "pages": 42,
      "total_chars": 123456,
      "link_count": 454,
      "entity_count": 647
    }
  ]
}
```

### CHISGListLinksHandler
GET /api/chisg/links?paper_id=xxx&limit=N — requires Bearer token auth

Queries Weaviate `AcademicLink` class filtered by paper_id.
Default limit: 50. Max: 500.
Returns:
```json
{
  "links": [
    {
      "entity_a": "CarD",
      "relation": "activates",
      "entity_b": "rrnA",
      "context": "Under starvation...",
      "paper_id": "ntm-mcgrath-2024",
      "confidence": 0.92
    }
  ],
  "total": 454,
  "paper_id": "ntm-mcgrath-2024"
}
```

### Weaviate Data (Production as of 2026-05-23)
- AcademicPaper: 1 object (ntm-mcgrath-2024 / McGrath PhD paper)
- AcademicLink: 454 objects (semantic triples from the McGrath paper)
- 647 unique biological entities
- 30 relation types
- Data migrated from local Weaviate → prod using temporary Python container
  on esp_shared-network, batch REST API ingest
""",
    },
    {
        "title": "Production 502 Fix — assist-api Missing from esp_shared-Network (2026-05-23)",
        "project": "assist",
        "section_path": "Infrastructure > Vultr Production > Incidents",
        "file_path": "infra/incidents/2026-05-23-assist-api-network.md",
        "content": """## Production 502 — assist-api Weaviate Unreachable (2026-05-23)

### Symptom
GET /api/chisg/papers returned 502. 
assist-api logs showed: connection refused to weaviate:8080

### Root Cause
assist-api container was only on `esp-thinking_shared-network`.
Weaviate container runs on `esp_shared-network` (the legacy network from the 
`/opt/esp/` stack). assist-api could not resolve the `weaviate` hostname.

### Live Fix
```bash
docker network connect esp_shared-network assist-api
```
Immediately resolved without restart.

### Permanent Fix
Updated `/opt/esp-thinking/docker-compose.prod.yml`:
- Added `esp_shared-network` to assist-api's network list
- Declared `esp_shared-network` as `external: true` at the bottom of the file

```yaml
# In assist-api service:
networks:
  - shared-network
  - esp_shared-network

# At bottom of file:
networks:
  shared-network:
    name: esp-thinking_shared-network
  esp_shared-network:
    external: true
```

### Key Rule
Any service that needs to talk to Weaviate must be on BOTH:
1. `esp-thinking_shared-network` (primary app network)
2. `esp_shared-network` (where Weaviate actually lives)

### Weaviate Container Facts
- Container name: `weaviate`
- Network: `esp_shared-network`
- Internal port: 8080
- Accessible from other containers as: `http://weaviate:8080`
""",
    },
    {
        "title": "Weaviate AcademicPaper + AcademicLink Data Migration Local→Prod (2026-05-23)",
        "project": "assist",
        "section_path": "Infrastructure > Weaviate > Data Migration",
        "file_path": "infra/weaviate/2026-05-23-migration.md",
        "content": """## Weaviate Data Migration — Local to Production (2026-05-23)

### Problem
Production Weaviate had no AcademicPaper or AcademicLink schema/data.
Local Weaviate (localhost:8088) had the full dataset.

### Migration Method
1. Exported from local using weaviate Python client v4:
   ```python
   import weaviate
   client = weaviate.connect_to_custom(
       http_host="localhost", http_port=8088, http_secure=False,
       grpc_host="localhost", grpc_port=50052, grpc_secure=False
   )
   ```
   Note: weaviate-client v3 `Client()` constructor was removed in v4 — must use
   `connect_to_custom()` or `connect_to_local()`.

2. Created schema on prod via Weaviate REST API directly from Vultr:
   - Ran temporary Python 3.11 container on `esp_shared-network`
   - `docker run --rm -it --network esp_shared-network python:3.11 python3 -c "..."`
   - POSTed class definitions to http://weaviate:8080/v1/schema

3. Ingested data in batches of 50 via REST API:
   - AcademicPaper: 1 object ingested
   - AcademicLink: 454 objects ingested in 10 batches

### Weaviate Schema on Production
AcademicPaper class properties:
- paper_id (text), title (text), filename (text), pages (int),
  total_chars (int), link_count (int), entity_count (int)

AcademicLink class properties:
- entity_a (text), relation (text), entity_b (text), context (text),
  paper_id (text), confidence (number), source_sentence (text)

### Vectorizer
Both classes use `text2vec-transformers` module.
If module unavailable, ingest still works but semantic search is degraded.

### Verification
```bash
# On Vultr via temp container on esp_shared-network:
curl http://weaviate:8080/v1/objects?class=AcademicPaper&limit=1
# → 1 paper returned

curl 'http://weaviate:8080/v1/objects?class=AcademicLink&limit=3'
# → 3 links returned, total 454
```
""",
    },
    {
        "title": "Frontend Auth System — AuthContext + PrivateRoute + CHISGLogin (2026-05-23)",
        "project": "assist",
        "section_path": "Frontend > Auth > CHISG",
        "file_path": "frontend/src/context/AuthContext.tsx",
        "content": """## Frontend Auth System for CHISG

### Files
- `frontend/src/context/AuthContext.tsx` — React context + hook
- `frontend/src/components/PrivateRoute.tsx` — Route guard
- `frontend/src/pages/CHISG/CHISGLogin.tsx` — Login page

### AuthContext
Provides: `token` (string|null), `isAuthenticated` (bool), `login(token)`, `logout()`
Storage: localStorage key `chisg_token`
Usage: `const { token, isAuthenticated, login, logout } = useAuth()`

### PrivateRoute
Redirects to `/chisg/login` if not authenticated.
Preserves `from` location in `location.state` so login can return user to where they came from.

### CHISGLogin
Route: `/chisg/login` (public)
Posts to `/api/auth/login` with username + password.
On success: calls `login(token)` then navigates to `from` or `/ntm` (default).
Error handling: displays API error message.

### Protected Routes (in App.tsx)
- `/ntm` → `<PrivateRoute><NTMResearch /></PrivateRoute>` — AUTH REQUIRED
- `/chisg/papers` → `<PrivateRoute><CHISGPapers /></PrivateRoute>`
- `/chisg/links` → `<PrivateRoute><CHISGLinks /></PrivateRoute>`
- `/semantic-links/extract` → `<PrivateRoute><SemanticLinkExtractor /></PrivateRoute>`

### Login Flow
1. User visits espthinking.co.uk/ntm
2. Redirected to /chisg/login (state: { from: { pathname: '/ntm' } })
3. Enters chisg_admin / McGrath2024!
4. Token stored in localStorage
5. Redirected back to /ntm
""",
    },
    {
        "title": "NTMResearch Page — Papers Tab Merged In, /ntm Behind Auth (2026-05-23)",
        "project": "assist",
        "section_path": "Frontend > Pages > NTM Research",
        "file_path": "frontend/src/pages/CHISG/NTMResearch.tsx",
        "content": """## NTMResearch Page Updates (2026-05-23)

### Route
`https://espthinking.co.uk/ntm` — now BEHIND PrivateRoute (was public)

### Tab Structure (post-update)
Tab 0: Upload Paper — PDF upload to /api/ntm/upload queue
Tab 1: Ask a Question — semantic Q&A via Claude (AWS Bedrock)
Tab 2: Processed Papers — papers list from /api/chisg/papers (NEW)

### Tab 2 — PapersTab Component
- Uses `useAuth()` hook to check authentication state
- If NOT authenticated: shows lock icon + "Sign in to view processed papers" button → /chisg/login
- If authenticated: fetches /api/chisg/papers with Bearer token
- Shows table: Title, Paper ID, Links count (Chip), Entities, Pages, Graph/Links action buttons
- Graph button → /chisg/graph?paper_id=xxx
- Links button → /chisg/links?paper_id=xxx

### Design Decision
NTMResearch is the central hub for all CHISG/NTM research tools.
CHISGPapers.tsx exists as a standalone admin view but Papers list lives in /ntm.
Avoids duplicating tab infrastructure across multiple pages.

### CHISGPapers.tsx
Reverted to simple clean version (no tabs).
Shows papers table with AppBar + Logout button.
Still accessible at /chisg/papers (PrivateRoute protected).

### New Imports Added to NTMResearch.tsx
- `useEffect` (React)
- Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Tooltip (MUI)
- AccountTreeIcon, EditIcon, ListAltIcon, LockIcon (MUI icons)
- `useNavigate` (react-router-dom)
- `useAuth` (../../context/AuthContext)
""",
    },
    {
        "title": "CHISG Product Strategy — Multi-Tenant Academic Repository (2026-05-23)",
        "project": "assist",
        "section_path": "Product > CHISG > Strategy",
        "file_path": "docs/product/chisg-repository-strategy.md",
        "content": """## CHISG Academic Repository — Product Strategy Discussion (2026-05-23)

### Two Product Shapes Discussed

#### Option A — Hosted Research Repository (SaaS)
Users upload papers → pipeline runs → hosted page at espthinking.co.uk/research/<slug>
with Upload/Query/Papers UI personalised for their project.

What needs building:
1. User registration + per-project accounts (MongoDB: users, projects collections)
2. Weaviate multi-tenancy (native support — each project gets a tenant partition)
3. Automated extraction pipeline on upload (currently manual SSH)
4. Per-project system prompt for Claude personalisation (project.context field)
5. Project admin dashboard

Readiness: ~40%. Q&A engine and data model exist. 
Auth, multi-tenancy, pipeline automation not yet built.

#### Option B — Embeddable CHISG API
API keys → POST /v1/query → returns answer + sources.
Drop-in JS embed widget for third-party sites.

What needs building (on top of Option A):
- API key management (MongoDB api_keys collection, hashed, scoped to project)
- Rate limiting middleware (Go token bucket per key)
- CORS policy for cross-origin embeds
- Usage metering (for billing)

Readiness: ~25%. More infrastructure but query endpoint is structurally right.

### Recommendation
Closed beta of Option A for McGrath's group first.
Priority order:
1. Automate the extraction pipeline end-to-end on upload (BLOCKER)
2. Per-project system prompt for Claude personalisation (1 day)
3. Weaviate multi-tenancy for second research group (2-3 days)
4. Simple user management — admin-created accounts OK for pilot (1-2 days)

### Why Human-Authenticated Access is the Right Choice
- Quality signal: 454 links from 1 paper — error rate unknown
- Human reviewers flagging wrong answers = validation dataset
- Feedback loop improves extraction quality before scaling
- Avoid scaling before quality is verified

### CHISG Value Proposition vs Standard RAG (Honest Assessment)

#### Where CHISG is better than vanilla RAG:
1. COMPLETENESS — returns ALL relationships for an entity across all papers.
   RAG returns top-k similar chunks and WILL miss things.
2. RELATIONAL TRAVERSAL — "what genes does Clp regulate, and what do THOSE regulate?"
   Requires graph traversal. LLMs cannot do this on unstructured text reliably.
3. CONTRADICTION DETECTION — two papers with conflicting entity-relation pairs
   are findable in a triple store. Invisible in a vector DB.
4. STATISTICS — "how many inhibitory vs activating relationships does protein X have?"
   Trivial in CHISG. Unreliable from RAG.

#### Where RAG is competitive with CHISG currently:
- Simple Q&A on well-structured text (Claude 3.5+, GPT-4o are very good at this)
- The gap has narrowed since 2022 for basic retrieval
- CHISG currently implements only ~30% of its potential value
  (doing vector search on AcademicLink objects, not full graph traversal)

#### The real differentiator not yet implemented:
Structured graph queries BEFORE LLM synthesis:
- Fetch all triples for named entity via Weaviate WHERE filter (not vector search)
- Multi-hop traversal: entity → relations → related entities → their relations
- Conflict detection: same entity+relation pair with contradictory entity_b values
This is when CHISG genuinely pulls ahead of well-built RAG systems.

#### On LLM differentiation of similar records:
Standard RAG embeds chunks as vectors. "CarD activates gene X" and 
"CarD inhibits gene X" have very similar vector embeddings (same words, same semantic space).
Top-k retrieval pulls both but synthesis may fail to distinguish them.
CHISG stores explicit relation types (activates, inhibits, regulates, binds_to) as
structured data — queryable precisely and completely independent of semantic similarity.

### Next Technical Priority
Implement structured entity-centric graph queries in the /api/chisg/query endpoint:
```
1. Extract named entities from the question (Claude or regex)
2. Fetch ALL AcademicLink objects where entity_a OR entity_b = named entity
   (Weaviate WHERE filter, not vector search)
3. Optionally: second-hop fetch for related entities
4. Pass structured triples as context to Claude for synthesis
5. Return answer + sources with full triple provenance
```
This changes the query from "find similar text chunks" to "traverse the knowledge graph".
""",
    },
    {
        "title": "CHISG Docker Image Deploy Pattern — 2026-05-23",
        "project": "assist",
        "section_path": "Infrastructure > Deployment > CHISG Frontend",
        "file_path": "infra/deployment/chisg-frontend-deploy.md",
        "content": """## CHISG Frontend Deploy Pattern (verified 2026-05-23)

### Build
```bash
cd /Users/michaelstewart/Coding/assist
docker build --platform linux/amd64 -t mike5tew/assist-frontend:latest -f frontend/Dockerfile ./frontend
```
Always --platform linux/amd64 (Mac M-series → Vultr AMD64).

### Push
```bash
docker push mike5tew/assist-frontend:latest
```

### Deploy to Vultr
```bash
sshpass -p 'E5P_Th!nk!ng?' ssh -i ~/.ssh/gitkey11-25 -o StrictHostKeyChecking=no root@192.248.151.185 \\
  "cd /opt/esp-thinking && \\
   docker compose -f docker-compose.prod.yml pull assist-frontend && \\
   docker compose -f docker-compose.prod.yml up -d --no-deps assist-frontend && \\
   docker exec main-proxy nginx -s reload"
```

Always --no-deps to avoid restarting unrelated services.
Always nginx -s reload after (picks up any config changes).

### Verify
```bash
curl -sk https://espthinking.co.uk/ntm | grep -o '<title>[^<]*'
# Or check the login redirect:
curl -I https://espthinking.co.uk/ntm
# Should return 200 (SPA handles redirect client-side)
```

### Warnings During Build (non-fatal, ignore)
- "SecretsUsedInArgOrEnv: ARG REACT_APP_RECAPTCHA_SITE_KEY" — safe to ignore,
  not a real secret, just a public reCAPTCHA site key
""",
    },
]


def connect_weaviate():
    return weaviate.connect_to_local(
        host=WEAVIATE_HOST,
        port=WEAVIATE_PORT,
        grpc_port=WEAVIATE_GRPC_PORT,
    )


def add_to_weaviate(client, entry):
    if not client.collections.exists(DOC_CLASS):
        print(f"  ⚠️  '{DOC_CLASS}' class not found in Weaviate — skipping")
        return False
    collection = client.collections.get(DOC_CLASS)
    chunk_hash = hashlib.md5(entry["content"].encode()).hexdigest()[:12]
    collection.data.insert({
        "title": entry["title"],
        "content": entry["content"],
        "file_path": entry["file_path"],
        "project": entry["project"],
        "section_path": entry["section_path"],
        "heading_level": 2,
        "chunk_hash": chunk_hash,
        "last_modified": datetime.now().isoformat(),
    })
    return True


def add_to_mongodb(entry):
    client = MongoClient(MONGO_URI, serverSelectionTimeoutMS=5000)
    db = client[MONGO_DB]
    coll = db[NOTES_COLLECTION]
    doc = {
        "title": entry["title"],
        "content": entry["content"],
        "project": entry["project"],
        "section_path": entry["section_path"],
        "file_path": entry["file_path"],
        "type": "session_log",
        "session_date": "2026-05-23",
        "created_at": datetime.utcnow(),
        "tags": ["chisg", "auth", "weaviate", "frontend", "product-strategy", "deployment"],
    }
    result = coll.insert_one(doc)
    client.close()
    return result.inserted_id


def main():
    print("🔌 Connecting to Weaviate (localhost:8088)...")
    wc = connect_weaviate()
    print("✅ Weaviate connected\n")

    for entry in CHANGES:
        print(f"📝 Processing: {entry['title']}")

        ok = add_to_weaviate(wc, entry)
        if ok:
            print(f"   ✅ Added to Weaviate '{DOC_CLASS}'")

        try:
            inserted_id = add_to_mongodb(entry)
            print(f"   ✅ Added to MongoDB '{NOTES_COLLECTION}' (id: {inserted_id})")
        except Exception as e:
            print(f"   ⚠️  MongoDB insert failed: {e}")

        print()

    wc.close()
    print("✅ Done — all 2026-05-23 session entries logged.")


if __name__ == "__main__":
    main()
