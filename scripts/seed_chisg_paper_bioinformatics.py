#!/usr/bin/env python3
"""
Seed the CHISG semantic_links collection with bioinformatics relationships
extracted from the 1998 Protein Science paper:
  "Domain assignment for protein structures using a consensus approach"
  Jones, Stewart, Michie, Swindells, Orengo, Thornton. Protein Science 7(2):233-242, 1998.

Usage (local):
  python3 seed_chisg_paper_bioinformatics.py --env local

Usage (production via SSH tunnel — run from Mac):
  python3 seed_chisg_paper_bioinformatics.py --env prod

The script connects directly to MongoDB and upserts all records.
"""

import argparse
import sys
from datetime import datetime, timezone
from pymongo import MongoClient, UpdateOne
from pymongo.errors import BulkWriteError

PAPER_TITLE = "Domain assignment for protein structures using a consensus approach"
PAPER_YEAR = 1998
DOMAIN = "bioinformatics"
SUBJECT = "protein_structure"
LINK_TYPE = "semantic"
STATUS = "validated"
BATCH_ID = "chisg_paper_jones1998_v1"

# All semantic links extracted from Jones et al. 1998
# Format: (source_term, target_term, forward_relation, inverse_relation, confidence,
#          source_generality, target_generality, statement, context, page_number)
LINKS = [
    # ── Consensus approach and accuracy ──────────────────────────────────────
    (
        "consensus approach",
        "domain assignment accuracy",
        "increases",
        "is_increased_by",
        0.98,
        0.7, 0.6,
        "A consensus approach to domain assignment increases the accuracy of protein structural domain identification.",
        "Four methods—PUU, DOMAK, DETECTIVE, and Islam—were combined in a consensus scheme.",
        1,
    ),
    (
        "consensus of four algorithms",
        "100% accuracy",
        "achieves",
        "is_achieved_by",
        0.99,
        0.6, 0.5,
        "Agreement among all four domain-assignment algorithms achieves 100% accuracy on the 55-chain validation set.",
        "When all four methods agree, the assignment is always correct.",
        4,
    ),
    (
        "consensus of three algorithms",
        "97.5% accuracy",
        "achieves",
        "is_achieved_by",
        0.99,
        0.6, 0.5,
        "Agreement among any three of the four algorithms achieves 97.5% accuracy.",
        "Three-way agreement still dramatically outperforms individual methods.",
        4,
    ),
    (
        "individual domain-assignment algorithms",
        "67–76% accuracy",
        "achieves",
        "is_achieved_by",
        0.99,
        0.5, 0.5,
        "Individual domain-assignment algorithms achieve only 67–76% accuracy when applied independently.",
        "Accuracy range across PUU, DOMAK, DETECTIVE, and Islam on the validation set.",
        4,
    ),
    (
        "consensus approach",
        "787 PDB protein chains",
        "validated_on",
        "used_for_validation_of",
        0.98,
        0.7, 0.5,
        "The consensus approach was validated on 787 non-redundant protein chains from the PDB.",
        "Large-scale application to PDB chains confirmed the method's generalisability.",
        6,
    ),

    # ── Individual algorithm characteristics ─────────────────────────────────
    (
        "PUU algorithm",
        "single-domain proteins",
        "identifies",
        "is_identified_by",
        0.92,
        0.6, 0.6,
        "The PUU algorithm (Holm & Sander) identifies single-domain proteins using normal mode analysis.",
        "PUU performs best on single-chain, single-domain proteins.",
        2,
    ),
    (
        "PUU algorithm",
        "normal mode analysis",
        "based_on",
        "underpins",
        0.95,
        0.6, 0.5,
        "The PUU algorithm is based on normal mode analysis to detect structurally rigid units.",
        "Holm & Sander's approach models domain rigidity via normal modes.",
        2,
    ),
    (
        "DOMAK algorithm",
        "protein structures",
        "tends_to_over_divide",
        "is_over_divided_by",
        0.93,
        0.6, 0.5,
        "The DOMAK algorithm (Siddiqui & Barton) tends to over-divide protein structures into too many domains.",
        "DOMAK fragments some genuine single-domain proteins into two or more parts.",
        2,
    ),
    (
        "DOMAK algorithm",
        "inter-domain contact minimization",
        "based_on",
        "underpins",
        0.90,
        0.6, 0.5,
        "DOMAK assigns domain boundaries by minimizing contacts between proposed domains.",
        "Contact-ratio criterion used to partition chains into domains.",
        2,
    ),
    (
        "DETECTIVE algorithm",
        "hydrophobic cores",
        "identifies",
        "is_identified_by",
        0.94,
        0.6, 0.6,
        "The DETECTIVE algorithm (Swindells) identifies domains by locating hydrophobic cores within protein structures.",
        "DETECTIVE locates compact hydrophobic regions as proxies for domain cores.",
        2,
    ),
    (
        "Islam method",
        "intramolecular contact areas",
        "uses",
        "is_used_by",
        0.93,
        0.6, 0.5,
        "The Islam domain-assignment method uses intramolecular contact areas to delineate domain boundaries.",
        "Contact-area analysis distinguishes tightly packed intra-domain residues from inter-domain contacts.",
        2,
    ),

    # ── Domain structural properties ──────────────────────────────────────────
    (
        "structural domain",
        "independently folding unit",
        "is",
        "exemplified_by",
        0.99,
        0.4, 0.5,
        "A structural domain is an independently folding unit within a protein chain.",
        "Domains fold and function as modular units; they are the fundamental unit of CATH classification.",
        1,
    ),
    (
        "intra-domain contacts",
        "structural domains",
        "characterize",
        "is_characterized_by",
        0.96,
        0.5, 0.5,
        "Intra-domain contacts are denser than inter-domain contacts, characterizing compact structural domain cores.",
        "The density of contacts within a domain is the principal criterion for most assignment algorithms.",
        2,
    ),
    (
        "domain size",
        "100 residues",
        "peaks_at",
        "is_peak_size_of",
        0.95,
        0.6, 0.5,
        "Domain size in protein structures peaks at approximately 100 residues.",
        "Histogram of domain sizes from PDB application shows modal value near 100 residues.",
        7,
    ),
    (
        "protein domains",
        "less than 200 residues",
        "are_predominantly",
        "is_size_of_majority_of",
        0.97,
        0.5, 0.5,
        "80.3% of protein domains identified in the PDB are less than 200 residues in length.",
        "Size distribution of all automatically assigned domains from the 787-chain set.",
        7,
    ),
    (
        "continuous domains",
        "protein domains",
        "constitute_majority_of",
        "includes",
        0.96,
        0.5, 0.4,
        "72% of protein domains are continuous (sequential), meaning domain residues form an uninterrupted chain segment.",
        "Discontinuous domains involve residue segments from different parts of the sequence.",
        7,
    ),
    (
        "alpha/beta domains",
        "protein domain classes",
        "are_most_common_in",
        "dominated_by",
        0.95,
        0.5, 0.5,
        "Alpha/beta domains are the most common structural class among PDB protein domains.",
        "Based on CATH secondary-structure class assignments from the PDB application.",
        7,
    ),
    (
        "single-domain proteins",
        "PDB protein chains",
        "constitute_majority_of",
        "includes",
        0.97,
        0.5, 0.4,
        "66.8% of assignable PDB protein chains are single-domain proteins.",
        "Majority of chains in the PDB that can be automatically assigned consist of one domain.",
        7,
    ),

    # ── PDB coverage and automated assignment ────────────────────────────────
    (
        "automated consensus method",
        "PDB protein chains",
        "assigns_55.7_percent_of",
        "is_auto_assigned_by",
        0.98,
        0.6, 0.5,
        "The automated consensus method is able to assign 55.7% of PDB protein chains without human intervention.",
        "Chains where three or more methods agree are assigned automatically; the rest require manual curation.",
        6,
    ),
    (
        "protein chain",
        "structural domains",
        "composed_of",
        "are_components_of",
        0.99,
        0.5, 0.5,
        "A protein chain is composed of one or more structural domains.",
        "Fundamental relationship in structural biology; the number of domains per chain ranges from 1 to >10.",
        1,
    ),

    # ── Overlap score and comparability ──────────────────────────────────────
    (
        "overlap score",
        "comparable domain assignments",
        "defines",
        "is_defined_by",
        0.97,
        0.5, 0.6,
        "An overlap score of ≥85% between two domain assignments defines them as comparable.",
        "Two assignments agree if ≥85% of residues are in the same domain in both.",
        3,
    ),
    (
        "overlap score",
        "domain assignment accuracy",
        "measures",
        "is_measured_by",
        0.97,
        0.5, 0.5,
        "Overlap score is the primary measure of domain assignment accuracy in this study.",
        "Gold-standard domain boundaries from CATH curators used as reference.",
        3,
    ),

    # ── Databases and resources ───────────────────────────────────────────────
    (
        "CATH classification",
        "consensus domain assignments",
        "uses",
        "is_used_by",
        0.98,
        0.6, 0.6,
        "The CATH structural classification database uses the consensus domain-assignment approach to classify protein chains.",
        "Jones et al. consensus method was adopted into the CATH pipeline.",
        8,
    ),
    (
        "CATH",
        "protein domain hierarchy",
        "classifies_by",
        "is_classified_by",
        0.99,
        0.5, 0.5,
        "CATH classifies protein domains into a hierarchy of Class, Architecture, Topology, and Homologous superfamily.",
        "Each domain is placed in the CATH hierarchy based on structural and sequence similarity.",
        1,
    ),
    (
        "PDB",
        "protein structure data",
        "contains",
        "is_contained_in",
        0.99,
        0.4, 0.5,
        "The Protein Data Bank (PDB) contains experimentally determined three-dimensional protein structure data.",
        "The PDB is the primary source of protein chains analysed in this study.",
        1,
    ),
    (
        "FSSP",
        "protein fold families",
        "classifies",
        "is_classified_by",
        0.95,
        0.5, 0.5,
        "The FSSP database classifies protein fold families using structure-based pairwise alignment.",
        "FSSP used alongside CATH and SCOP as a comparative classification resource.",
        1,
    ),
    (
        "SCOP",
        "protein structural classes",
        "classifies",
        "is_classified_by",
        0.96,
        0.5, 0.5,
        "SCOP (Structural Classification of Proteins) classifies proteins by structural and evolutionary relationships.",
        "SCOP provides an alternative hierarchical classification to CATH.",
        1,
    ),
    (
        "3Dee",
        "protein domain definitions",
        "contains",
        "is_contained_in",
        0.94,
        0.5, 0.5,
        "The 3Dee database contains a comprehensive set of protein domain definitions from the PDB.",
        "3Dee provides an independent reference set of domain boundaries for comparison.",
        1,
    ),

    # ── Structural determinants of domain boundaries ─────────────────────────
    (
        "protein secondary structure",
        "domain boundary detection",
        "influences",
        "is_influenced_by",
        0.90,
        0.6, 0.6,
        "Protein secondary structure elements (helices, strands) influence where domain boundaries are detected.",
        "Boundaries rarely split secondary structure elements; they typically fall in loop regions.",
        2,
    ),
    (
        "domain boundary",
        "inter-domain contact minimization",
        "determined_by",
        "determines",
        0.92,
        0.5, 0.5,
        "Domain boundaries are determined by minimizing inter-domain contacts and maximizing intra-domain contacts.",
        "The contact-area criterion is used by DOMAK and Islam methods.",
        2,
    ),
    (
        "multi-domain proteins",
        "domain-swapping artefacts",
        "prone_to",
        "is_caused_by",
        0.85,
        0.5, 0.5,
        "Multi-domain proteins are prone to domain-swapping artefacts that complicate automated assignment.",
        "Artificially concatenated or swapped domains can confound assignment algorithms.",
        5,
    ),

    # ── Validation set ────────────────────────────────────────────────────────
    (
        "55-chain validation set",
        "CATH-curated domain assignments",
        "validated_against",
        "serves_as_reference_for",
        0.99,
        0.5, 0.6,
        "The 55-chain validation set was validated against manually curated CATH domain assignments as gold standard.",
        "Expert CATH curators provided the reference domain boundaries for the 55-chain set.",
        3,
    ),
    (
        "consensus approach",
        "55-chain validation set",
        "evaluated_on",
        "used_for_evaluation_of",
        0.99,
        0.7, 0.5,
        "The consensus approach was evaluated on a 55-chain validation set before large-scale PDB application.",
        "55 chains were used as a held-out test set to measure accuracy before PDB-wide deployment.",
        3,
    ),

    # ── Algorithm interplay ───────────────────────────────────────────────────
    (
        "PUU algorithm",
        "DOMAK algorithm",
        "complements",
        "is_complemented_by",
        0.88,
        0.6, 0.6,
        "PUU and DOMAK use complementary physical principles—normal mode analysis vs. contact minimization.",
        "Combining PUU and DOMAK in consensus reduces the biases inherent in either method alone.",
        2,
    ),
    (
        "DETECTIVE algorithm",
        "Islam method",
        "complements",
        "is_complemented_by",
        0.87,
        0.6, 0.6,
        "DETECTIVE (hydrophobic-core detection) and the Islam contact-area method provide complementary domain evidence.",
        "Together, the four methods cover different aspects of domain organisation.",
        2,
    ),
]


def make_composite_key(source: str, target: str, relation: str) -> str:
    return f"{source}::{target}::{relation}"


def build_doc(row) -> dict:
    (
        source, target, fwd_rel, inv_rel, confidence,
        src_gen, tgt_gen, statement, context, page_num,
    ) = row

    composite_key = make_composite_key(source, target, fwd_rel)
    now = datetime.now(timezone.utc)

    return {
        "composite_key": composite_key,
        "source_term": source,
        "target_term": target,
        "forward_relation": fwd_rel,
        "inverse_relation": inv_rel,
        "relation_type": fwd_rel,  # legacy field
        "statement": statement,
        "context": context,
        "confidence": confidence,
        "source_term_generality": src_gen,
        "target_term_generality": tgt_gen,
        "semantic_distance": round(1.0 - confidence, 2),
        "relationship_strength": confidence,
        "domain": DOMAIN,
        "subject": SUBJECT,
        "link_type": LINK_TYPE,
        "status": STATUS,
        "is_manual": True,
        "source_title": PAPER_TITLE,
        "page_number": page_num,
        "batch_id": BATCH_ID,
        "created_at": now,
        "is_case_specific": False,
        "metadata": {
            "paper_year": PAPER_YEAR,
            "paper_authors": "Jones DT, Stewart M, Michie A, Swindells MB, Orengo C, Thornton JM",
            "journal": "Protein Science",
            "volume": "7",
            "issue": "2",
            "pages": "233-242",
        },
    }


def seed(mongo_uri: str, db_name: str, dry_run: bool = False):
    client = MongoClient(mongo_uri, serverSelectionTimeoutMS=10_000)
    db = client[db_name]
    col = db["semantic_links"]

    docs = [build_doc(row) for row in LINKS]
    print(f"[seed] Prepared {len(docs)} semantic link documents.")

    if dry_run:
        import json
        for d in docs:
            print(json.dumps({k: v for k, v in d.items() if k not in ("created_at",)}, indent=2, default=str))
        print("[dry-run] No writes performed.")
        return

    ops = [
        UpdateOne(
            {"composite_key": d["composite_key"]},
            {"$setOnInsert": d},
            upsert=True,
        )
        for d in docs
    ]

    try:
        result = col.bulk_write(ops, ordered=False)
        print(
            f"[seed] Done — inserted: {result.upserted_count}, "
            f"matched (already existed): {result.matched_count}"
        )
    except BulkWriteError as e:
        print(f"[seed] BulkWriteError: {e.details}")
        raise

    # Print current count
    total = col.count_documents({})
    domain_count = col.count_documents({"domain": DOMAIN})
    print(f"[seed] semantic_links total: {total}, domain='{DOMAIN}': {domain_count}")
    client.close()


def main():
    parser = argparse.ArgumentParser(description="Seed CHISG semantic_links from Jones et al. 1998")
    parser.add_argument(
        "--env",
        choices=["local", "prod"],
        default="local",
        help="Target MongoDB environment (default: local)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Print docs without writing",
    )
    parser.add_argument(
        "--ndjson",
        action="store_true",
        help="Print NDJSON to stdout (one doc per line) for use with mongoimport",
    )
    args = parser.parse_args()

    if args.ndjson:
        import json
        docs = [build_doc(row) for row in LINKS]
        for d in docs:
            # Convert datetime to ISO string for JSON serialisation
            d["created_at"] = d["created_at"].isoformat()
            print(json.dumps(d))
        return

    if args.env == "local":
        uri = "mongodb://admin:password123@localhost:27018"
        db_name = "esp_organizer"
        print(f"[seed] Connecting to LOCAL MongoDB at localhost:27018 / {db_name}")
    else:
        # Production — expects an SSH tunnel or direct connection
        # To use: open a tunnel first:
        #   sshpass -p 'E5P_Th!nk!ng?' ssh -o StrictHostKeyChecking=no \
        #     -L 27017:localhost:27017 root@192.248.151.185 -N &
        # Then run:  python3 seed_chisg_paper_bioinformatics.py --env prod
        uri = "mongodb://esp_admin:78NTGg8QgHEq5n4CLDBP@localhost:27017"
        db_name = "esp_organizer"
        print(f"[seed] Connecting to PRODUCTION MongoDB via SSH tunnel at localhost:27017 / {db_name}")
        print("[seed] NOTE: Ensure SSH tunnel is open: sshpass -p 'E5P_Th!nk!ng?' ssh -o StrictHostKeyChecking=no -L 27017:localhost:27017 root@192.248.151.185 -N &")

    seed(uri, db_name, dry_run=args.dry_run)


if __name__ == "__main__":
    main()
