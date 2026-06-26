#!/usr/bin/env python3
"""
CHISG Academic Repository — Weaviate Schema & Ingestion
=========================================================
Creates Weaviate collections for academic paper knowledge graphs and
ingests extraction output from extract_chisg_papers.py.

Collections:
  1. AcademicEntity  — Unique entities (organisms, genes, proteins, mechanisms)
  2. AcademicLink    — Semantic links between entities (with provenance)
  3. AcademicPaper   — Paper metadata (title, authors, year)

Usage:
    # Create schema (drops existing)
    python3 scripts/ingest_chisg_academic.py --create-schema

    # Ingest extraction output
    python3 scripts/ingest_chisg_academic.py --input data/chisg/mcgrath/extraction.json

    # Both
    python3 scripts/ingest_chisg_academic.py --create-schema --input data/chisg/mcgrath/extraction.json

    # Use different Weaviate instance
    python3 scripts/ingest_chisg_academic.py --weaviate-port 8088 --input data/chisg/mcgrath/extraction.json

    # Stats only (don't ingest, just show what's there)
    python3 scripts/ingest_chisg_academic.py --stats
"""

import json
import os
import sys
import argparse
from pathlib import Path
from collections import defaultdict, Counter

SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DEFAULT_INPUT = PROJECT_ROOT / "data" / "chisg" / "mcgrath" / "extraction.json"

ENTITY_CLASS = "AcademicEntity"
LINK_CLASS = "AcademicLink"
PAPER_CLASS = "AcademicPaper"

BACKWARD_RELATIONS = {
    "causes": "is caused by",
    "enables": "is enabled by",
    "inhibits": "is inhibited by",
    "is composed of": "is a component of",
    "is a type of": "has subtype",
    "has property": "is a property of",
    "has value": "is a value of",
    "is used for": "uses",
    "is found in": "contains",
    "is an example of": "has example",
    "determines": "is determined by",
    "represents": "is represented by",
    "contradicts": "is contradicted by",
    "correlates with": "correlates with",
    "supports": "is supported by",
    "refutes": "is refuted by",
    "regulates": "is regulated by",
    "is a mechanism of": "has mechanism",
    "is a marker for": "has marker",
    "interacts with": "interacts with",
}


def get_weaviate_client(host="localhost", port=8081, grpc_port=50051):
    """Connect to Weaviate."""
    import weaviate
    client = weaviate.connect_to_local(host=host, port=port, grpc_port=grpc_port)
    return client


def create_schema(client):
    """Create the academic repository schema."""
    import weaviate.classes.config as wvc

    # --- AcademicPaper ---
    if client.collections.exists(PAPER_CLASS):
        print(f"  Deleting existing {PAPER_CLASS}...")
        client.collections.delete(PAPER_CLASS)

    print(f"  Creating {PAPER_CLASS}...")
    client.collections.create(
        name=PAPER_CLASS,
        description="An academic paper in the repository",
        vectorizer_config=wvc.Configure.Vectorizer.none(),
        properties=[
            wvc.Property(name="paper_id", data_type=wvc.DataType.TEXT),
            wvc.Property(name="title", data_type=wvc.DataType.TEXT),
            wvc.Property(name="filename", data_type=wvc.DataType.TEXT),
            wvc.Property(name="pages", data_type=wvc.DataType.INT),
            wvc.Property(name="total_chars", data_type=wvc.DataType.INT),
            wvc.Property(name="link_count", data_type=wvc.DataType.INT),
            wvc.Property(name="entity_count", data_type=wvc.DataType.INT),
        ]
    )

    # --- AcademicEntity ---
    if client.collections.exists(ENTITY_CLASS):
        print(f"  Deleting existing {ENTITY_CLASS}...")
        client.collections.delete(ENTITY_CLASS)

    print(f"  Creating {ENTITY_CLASS}...")
    client.collections.create(
        name=ENTITY_CLASS,
        description="A unique entity from academic papers (organism, gene, protein, mechanism, etc.)",
        vectorizer_config=wvc.Configure.Vectorizer.none(),
        properties=[
            wvc.Property(name="name", data_type=wvc.DataType.TEXT,
                         description="Canonical entity name"),
            wvc.Property(name="name_lower", data_type=wvc.DataType.TEXT,
                         description="Lowercase name for deduplication"),
            wvc.Property(name="entity_type", data_type=wvc.DataType.TEXT,
                         description="Inferred type: organism, gene, protein, drug, mechanism, structure, other"),
            wvc.Property(name="domain", data_type=wvc.DataType.TEXT,
                         description="Domain: microbiology, immunology, pharmacology, etc."),
            wvc.Property(name="paper_ids", data_type=wvc.DataType.TEXT_ARRAY,
                         description="Papers mentioning this entity"),
            wvc.Property(name="outgoing_relations", data_type=wvc.DataType.TEXT_ARRAY,
                         description="Relations where this entity is source"),
            wvc.Property(name="incoming_relations", data_type=wvc.DataType.TEXT_ARRAY,
                         description="Relations where this entity is target"),
            wvc.Property(name="outgoing_count", data_type=wvc.DataType.INT),
            wvc.Property(name="incoming_count", data_type=wvc.DataType.INT),
            wvc.Property(name="total_mentions", data_type=wvc.DataType.INT),
        ]
    )

    # --- AcademicLink ---
    if client.collections.exists(LINK_CLASS):
        print(f"  Deleting existing {LINK_CLASS}...")
        client.collections.delete(LINK_CLASS)

    print(f"  Creating {LINK_CLASS}...")
    client.collections.create(
        name=LINK_CLASS,
        description="A semantic link between two academic entities",
        vectorizer_config=wvc.Configure.Vectorizer.none(),
        properties=[
            wvc.Property(name="entity_a", data_type=wvc.DataType.TEXT),
            wvc.Property(name="relation", data_type=wvc.DataType.TEXT,
                         description="Forward relation (from controlled vocabulary)"),
            wvc.Property(name="entity_b", data_type=wvc.DataType.TEXT),
            wvc.Property(name="backward_relation", data_type=wvc.DataType.TEXT,
                         description="Reverse relation"),
            wvc.Property(name="statement", data_type=wvc.DataType.TEXT,
                         description="Natural language: 'entity_a relation entity_b'"),
            wvc.Property(name="context", data_type=wvc.DataType.TEXT,
                         description="Source quote or paraphrase"),
            wvc.Property(name="paper_id", data_type=wvc.DataType.TEXT),
            wvc.Property(name="chunk_id", data_type=wvc.DataType.TEXT),
            wvc.Property(name="page_num", data_type=wvc.DataType.INT),
        ]
    )

    print(f"\n  ✓ Schema created: {PAPER_CLASS}, {ENTITY_CLASS}, {LINK_CLASS}")


def infer_entity_type(name):
    """Heuristic entity type inference for microbiology domain."""
    lower = name.lower()

    # Organisms
    if any(w in lower for w in ["mycobacterium", "tuberculosis", "m. tb", "m.tb",
                                  "bacillus", "bacterium", "bacteria", "pathogen",
                                  "e. coli", "staphylococcus", "streptococcus"]):
        return "organism"

    # Genes/genetic
    if any(w in lower for w in ["gene", "operon", "promoter", "codon", "allele"]):
        return "gene"
    # Gene names are often short uppercase: katG, rpoB, inhA, etc.
    if len(name) <= 6 and any(c.isupper() for c in name) and any(c.islower() for c in name):
        return "gene"

    # Proteins
    if any(w in lower for w in ["protein", "enzyme", "kinase", "synthase", "reductase",
                                  "polymerase", "transferase", "receptor", "antibody",
                                  "cytokine", "interleukin", "interferon", "tnf",
                                  "catalase", "peroxidase"]):
        return "protein"

    # Drugs
    if any(w in lower for w in ["isoniazid", "rifampicin", "ethambutol", "pyrazinamide",
                                  "streptomycin", "fluoroquinolone", "antibiotic",
                                  "drug", "treatment", "therapy", "regimen"]):
        return "drug"

    # Cellular structures
    if any(w in lower for w in ["cell wall", "membrane", "ribosome", "dna", "rna",
                                  "lipid", "peptidoglycan", "mycolic acid", "capsule"]):
        return "structure"

    # Immune
    if any(w in lower for w in ["macrophage", "t cell", "b cell", "neutrophil",
                                  "immune", "immunity", "phagocyte", "granuloma",
                                  "inflammation", "cd4", "cd8", "nk cell"]):
        return "immune_component"

    # Processes/mechanisms
    if any(w in lower for w in ["resistance", "mutation", "infection", "replication",
                                  "transcription", "translation", "metabolism",
                                  "pathogenesis", "virulence", "latency", "dormancy"]):
        return "mechanism"

    return "other"


def build_entity_index(extractions):
    """Build a deduplicated entity index from all extraction links."""
    entities = defaultdict(lambda: {
        "name": "",
        "mentions": 0,
        "paper_ids": set(),
        "outgoing": [],
        "incoming": [],
    })

    for ext in extractions:
        paper_id = ext.get("paper_id", "unknown")
        for link in ext.get("links", []):
            a = link.get("entity_a", "").strip()
            b = link.get("entity_b", "").strip()
            rel = link.get("relation", "")
            if not a or not b:
                continue

            a_key = a.lower()
            b_key = b.lower()

            # Entity A
            entities[a_key]["name"] = a  # keep original casing
            entities[a_key]["mentions"] += 1
            entities[a_key]["paper_ids"].add(paper_id)
            entities[a_key]["outgoing"].append(f"{rel} → {b}")

            # Entity B
            entities[b_key]["name"] = b
            entities[b_key]["mentions"] += 1
            entities[b_key]["paper_ids"].add(paper_id)
            entities[b_key]["incoming"].append(f"{a} → {rel}")

    return entities


def ingest_data(client, input_path):
    """Ingest extraction output into Weaviate."""
    print(f"\n  Loading {input_path}...")
    with open(input_path) as f:
        data = json.load(f)

    papers = data.get("papers", {})
    extractions = data.get("extractions", [])

    if not extractions:
        print("  No extractions to ingest.")
        return

    # Build entity index
    entity_index = build_entity_index(extractions)
    print(f"  Found {len(entity_index)} unique entities from {len(extractions)} chunks")

    # Count links per paper
    paper_link_counts = Counter()
    paper_entity_sets = defaultdict(set)
    all_links = []

    for ext in extractions:
        paper_id = ext.get("paper_id", "unknown")
        chunk_id = ext.get("chunk_id", "")
        for link in ext.get("links", []):
            a = link.get("entity_a", "").strip()
            b = link.get("entity_b", "").strip()
            if a and b:
                paper_link_counts[paper_id] += 1
                paper_entity_sets[paper_id].add(a.lower())
                paper_entity_sets[paper_id].add(b.lower())
                all_links.append({**link, "paper_id": paper_id, "chunk_id": chunk_id})

    # --- Ingest Papers ---
    print(f"\n  Ingesting {len(papers)} papers into {PAPER_CLASS}...")
    paper_col = client.collections.get(PAPER_CLASS)
    with paper_col.batch.dynamic() as batch:
        for paper_id, meta in papers.items():
            batch.add_object(properties={
                "paper_id": paper_id,
                "title": meta.get("title", "Unknown"),
                "filename": meta.get("filename", ""),
                "pages": meta.get("pages", 0),
                "total_chars": meta.get("total_chars", 0),
                "link_count": paper_link_counts.get(paper_id, 0),
                "entity_count": len(paper_entity_sets.get(paper_id, set())),
            })

    # --- Ingest Entities ---
    print(f"  Ingesting {len(entity_index)} entities into {ENTITY_CLASS}...")
    entity_col = client.collections.get(ENTITY_CLASS)
    with entity_col.batch.dynamic() as batch:
        for key, ent in entity_index.items():
            batch.add_object(properties={
                "name": ent["name"],
                "name_lower": key,
                "entity_type": infer_entity_type(ent["name"]),
                "domain": "microbiology",
                "paper_ids": list(ent["paper_ids"]),
                "outgoing_relations": ent["outgoing"][:50],  # cap array size
                "incoming_relations": ent["incoming"][:50],
                "outgoing_count": len(ent["outgoing"]),
                "incoming_count": len(ent["incoming"]),
                "total_mentions": ent["mentions"],
            })

    # --- Ingest Links ---
    print(f"  Ingesting {len(all_links)} links into {LINK_CLASS}...")
    link_col = client.collections.get(LINK_CLASS)
    with link_col.batch.dynamic() as batch:
        for link in all_links:
            a = link.get("entity_a", "")
            rel = link.get("relation", "")
            b = link.get("entity_b", "")
            backward = BACKWARD_RELATIONS.get(rel, f"reverse of {rel}")
            statement = f"{a} {rel} {b}"

            batch.add_object(properties={
                "entity_a": a,
                "relation": rel,
                "entity_b": b,
                "backward_relation": backward,
                "statement": statement,
                "context": link.get("context", ""),
                "paper_id": link.get("paper_id", ""),
                "chunk_id": link.get("chunk_id", ""),
            })

    print(f"\n  ✓ Ingestion complete:")
    print(f"    Papers:   {len(papers)}")
    print(f"    Entities: {len(entity_index)}")
    print(f"    Links:    {len(all_links)}")


def show_stats(client):
    """Show current collection stats."""
    print(f"\n═══ Weaviate Academic Repository Stats ═══\n")
    for cls_name in [PAPER_CLASS, ENTITY_CLASS, LINK_CLASS]:
        if client.collections.exists(cls_name):
            col = client.collections.get(cls_name)
            count = col.aggregate.over_all(total_count=True).total_count
            print(f"  {cls_name}: {count} objects")
        else:
            print(f"  {cls_name}: not created")


def main():
    parser = argparse.ArgumentParser(description="CHISG academic repository — Weaviate schema & ingestion")
    parser.add_argument("--create-schema", action="store_true", help="Create/recreate Weaviate schema")
    parser.add_argument("--input", type=Path, default=DEFAULT_INPUT, help="Extraction JSON to ingest")
    parser.add_argument("--ingest", action="store_true", help="Ingest extraction data")
    parser.add_argument("--stats", action="store_true", help="Show collection stats")
    parser.add_argument("--weaviate-host", default="localhost")
    parser.add_argument("--weaviate-port", type=int, default=8081)
    parser.add_argument("--grpc-port", type=int, default=50051)
    args = parser.parse_args()

    # Default action: if no flags, show help
    if not any([args.create_schema, args.ingest, args.stats]) and not args.input.exists():
        parser.print_help()
        return

    try:
        client = get_weaviate_client(args.weaviate_host, args.weaviate_port, args.grpc_port)
        print(f"✓ Connected to Weaviate at {args.weaviate_host}:{args.weaviate_port}")
    except Exception as e:
        print(f"✗ Failed to connect to Weaviate: {e}")
        sys.exit(1)

    try:
        if args.create_schema:
            print(f"\n═══ Creating Schema ═══\n")
            create_schema(client)

        if args.ingest or (args.input.exists() and not args.stats and not args.create_schema):
            if not args.input.exists():
                print(f"Error: Input file not found: {args.input}")
                sys.exit(1)
            ingest_data(client, args.input)

        if args.stats or args.create_schema or args.ingest:
            show_stats(client)

    finally:
        client.close()


if __name__ == "__main__":
    main()
