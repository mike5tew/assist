#!/usr/bin/env python3
"""
Migrate skill_links → semantic_links for the CHISG graph demo.

Reads the local skill_links collection and writes an NDJSON file in the
SemanticLink schema expected by ChisgGraphHandler.

Usage:
    python3 migrate_skill_links_to_semantic_links.py
    # Produces: /tmp/semantic_links_seed.ndjson

Then on production:
    docker cp /tmp/semantic_links_seed.ndjson mongodb:/tmp/
    docker exec mongodb mongoimport \
        --authenticationDatabase admin \
        -u esp_admin -p 78NTGg8QgHEq5n4CLDBP \
        --db esp_organizer \
        --collection semantic_links \
        --file /tmp/semantic_links_seed.ndjson \
        --mode upsert \
        --upsertFields composite_key
"""

import json
from datetime import datetime, timezone
from pymongo import MongoClient

LOCAL_URI = "mongodb://admin:password123@localhost:27018/esp_organizer?authSource=admin"
OUTPUT_PATH = "/tmp/semantic_links_seed.ndjson"


def transform(doc):
    src = doc.get("source_name", "")
    tgt = doc.get("target_name", "")
    rel = doc.get("relation", "enables")

    # Derive a sensible inverse
    inverse_map = {
        "enables": "is_enabled_by",
        "requires": "is_required_by",
        "supports": "is_supported_by",
        "includes": "is_included_in",
        "precedes": "follows",
        "causes": "is_caused_by",
    }
    inverse = inverse_map.get(rel, f"is_{rel.replace(' ', '_')}_by")

    composite_key = f"{src}::{tgt}::{rel}"

    return {
        "composite_key": composite_key,
        "statement": f"{src} {rel} {tgt}",
        "source_term": src,
        "target_term": tgt,
        "forward_relation": rel,
        "inverse_relation": inverse,
        "domain": "education",
        "subject": "chisg",
        "confidence": 1.0,
        "link_type": "skill_link",
        "source_chisg_id": doc.get("source_chisg_id", ""),
        "target_chisg_id": doc.get("target_chisg_id", ""),
        "created_at": datetime.now(timezone.utc).isoformat(),
    }


def main():
    client = MongoClient(LOCAL_URI)
    db = client["esp_organizer"]
    docs = list(db["skill_links"].find({}))
    print(f"Found {len(docs)} skill_links documents")

    count = 0
    with open(OUTPUT_PATH, "w") as f:
        for doc in docs:
            doc.pop("_id", None)
            transformed = transform(doc)
            f.write(json.dumps(transformed) + "\n")
            count += 1

    print(f"Wrote {count} records to {OUTPUT_PATH}")
    print()
    print("Next steps — run on production:")
    print(f"  scp {OUTPUT_PATH} root@192.248.151.185:/tmp/")
    print("  ssh root@192.248.151.185 \\")
    print("    'docker cp /tmp/semantic_links_seed.ndjson mongodb:/tmp/ && \\")
    print("     docker exec mongodb mongoimport \\")
    print("       --authenticationDatabase admin \\")
    print("       -u esp_admin -p 78NTGg8QgHEq5n4CLDBP \\")
    print("       --db esp_organizer \\")
    print("       --collection semantic_links \\")
    print("       --file /tmp/semantic_links_seed.ndjson \\")
    print("       --mode upsert \\")
    print("       --upsertFields composite_key'")


if __name__ == "__main__":
    main()
