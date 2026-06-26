#!/usr/bin/env python3
"""
Log A-Level Physics /aphy changes from 2026-05-20 to Weaviate and MongoDB.
Usage: python3 log_aphy_changes_20260520.py
"""

import weaviate
import hashlib
import sys
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
        "title": "A-Level Physics /aphy — 4 Revision Modes Built (2026-05-20)",
        "project": "assist",
        "section_path": "Frontend > Pages > APhysicsRevision",
        "file_path": "frontend/src/pages/APhysicsRevision.tsx",
        "content": """## A-Level Physics Revision Page — /aphy

### Route
`https://espthinking.co.uk/aphy`

### Files
- `frontend/src/pages/APhysicsRevision.tsx` — main component
- `frontend/src/pages/aPhysicsContent.ts` — 44 lessons, 4 modules

### 4 Revision Modes
1. **read** — Full markdown render (MUI tables, KaTeX equations via ReactMarkdown)
2. **speedread** — RSVP word-by-word, 100-600 WPM, configurable delay
3. **audio** — Web Speech API TTS with pause/resume/stop
4. **flashcards** — Flip cards: h2/h3 headings as fronts, body text as ReactMarkdown backs

### Content (aPhysicsContent.ts)
- 44 lessons across 4 modules (IDs 4001–4004)
- Module 4001: Thermal Physics
- Module 4002: Nuclear Physics
- Module 4003: Particles & Radiation
- Module 4004: Exam Skills — 6-Mark Questions

### Dependencies added
- `remark-gfm` installed via `npm install remark-gfm` — enables GFM table parsing
- `react-markdown`, `remark-math`, `rehype-katex`, `katex` already present

### ReactMarkdown config (card back + read mode)
```jsx
<ReactMarkdown
  remarkPlugins={[remarkMath, remarkGfm]}
  rehypePlugins={[rehypeKatex]}
  components={mdComponents}
>
```

### MUI Table mapping (mdComponents)
- `table` → `TableContainer + Table`
- `thead` → `TableHead` (accent background)
- `tbody` → `TableBody`
- `tr` → `TableRow`
- `th` → `TableCell` (white text)
- `td` → `TableCell`
""",
    },
    {
        "title": "A-Level Physics /aphy — Flashcard Table & Equation Rendering Fixed (2026-05-20)",
        "project": "assist",
        "section_path": "Frontend > Pages > APhysicsRevision > Flashcards",
        "file_path": "frontend/src/pages/APhysicsRevision.tsx",
        "content": """## Flashcard Rendering Bugs Fixed

### Bug 1: [table row] shown instead of table
- **Root cause**: `extractFlashcards()` had:
  `.replace(/\\|[^\\n]+\\|/g, '[table row]')` — stripping all markdown table cells
- **Fix**: Removed that replace entirely.
- **Result**: Card backs now render full GFM tables via ReactMarkdown + remark-gfm.

### Bug 2: [equation] / [eq] shown instead of LaTeX
- **Root cause**: `extractFlashcards()` had three replacements stripping LaTeX:
  1. `$$...$$` → `[equation]`
  2. `$...$` → `[eq]`
  3. `` ```...``` `` → `[code block]`
- **Fix**: Removed all three replacements from `extractFlashcards()`.
- **Result**: Block and inline LaTeX equations render correctly via rehype-katex.

### Card Back Design
- Light background: `#f0f4ff` with dark text `#1a1a2e` (avoids white-on-dark table issue)
- `alignItems: 'stretch'`, `textAlign: 'left'` — prevents table/equation squeezing
- Flip animation: `rotateY(180deg)`, `backfaceVisibility: hidden`

### extractFlashcards() logic (post-fix)
- Parses markdown; h2/h3 headings become card fronts
- Body text (bullet-normalised only) becomes card back markdown — no stripping
- `markdownToPlainText()` still used for audio/speedread modes (unchanged)

### Build & Deploy
- `docker build --no-cache --platform linux/amd64 -t mike5tew/assist-frontend:latest ./frontend`
- Digest: `sha256:9d3bbc8e6dc519f233790f8583384291ae8e323e69d30b09e28fd90775ea407d`
- Deployed: `cd /opt/esp-thinking && docker compose -f docker-compose.prod.yml pull assist-frontend && docker compose -f docker-compose.prod.yml up -d --force-recreate assist-frontend`
- Live: `https://espthinking.co.uk/aphy`
""",
        "tags": ["frontend", "aphy", "flashcards", "reactmarkdown", "katex"],
    },
    {
        "title": "A-Level Physics /aphy — Narrative-First Definitions Note (2026-05-30)",
        "project": "assist",
        "section_path": "Docs > LinkedIn Posts > Working Note > Narrative-First Definitions in A-Level Physics",
        "file_path": "docs/linkedin/LINKEDIN_POSTS.md",
        "content": """## Narrative-First Definitions in A-Level Physics

### Core insight
Definitions become stronger when they are taught as compressed physical narratives rather than static labels.

### Captured examples from electromagnetism
- Magnetic flux density: how tightly the field is packed through an area.
- Magnetic flux: how much field is caught by a chosen area.
- Flux linkage: how many field-lines-worth of change a coil intercepts across all turns.
- Induced emf: the rate at which those line cuts are changing.
- Lenz's law: the system pushes back against the change that created the induced current.

### Teaching implication
The missing step for many learners is between naming the quantity and using the equation. The definition should already carry the mechanism, so the formula reads as shorthand for a causal story rather than a disconnected symbol rule.

### Pedagogic ladder
1. Name the object.
2. Give the definition.
3. Translate the definition into a physical story.
4. Show the equation as shorthand for that story.

### Summary
A strong science definition is a frozen mechanism. If the learner can re-expand it into the underlying process, they understand the idea rather than merely recalling the term.
""",
        "tags": ["aphy", "pedagogy", "definitions", "narrative", "electromagnetism"],
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
        print(f"  ⚠️  '{DOC_CLASS}' class not found in Weaviate — skipping Weaviate insert")
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
    tags = entry.get("tags", ["frontend", "aphy", "flashcards", "reactmarkdown", "katex"])
    doc = {
        "title": entry["title"],
        "content": entry["content"],
        "project": entry["project"],
        "section_path": entry["section_path"],
        "file_path": entry["file_path"],
        "type": "feature_change",
        "created_at": datetime.utcnow(),
        "tags": tags,
    }
    result = coll.insert_one(doc)
    client.close()
    return result.inserted_id


def main():
    title_filter = " ".join(sys.argv[1:]).strip().lower()
    entries = CHANGES
    if title_filter:
        entries = [entry for entry in CHANGES if title_filter in entry["title"].lower()]
        if not entries:
            print(f"⚠️  No CHANGES entries matched filter: {title_filter}")
            return

    print("🔌 Connecting to Weaviate (localhost:8088)...")
    wc = connect_weaviate()
    print("✅ Weaviate connected\n")

    for entry in entries:
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
    print("✅ Done — all entries logged.")


if __name__ == "__main__":
    main()
