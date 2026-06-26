#!/usr/bin/env python3
"""
Analyse extraction_full_sonnet.json to discover natural thesaurus clusters.

Outputs:
1. Full relation inventory (controlled + invalid) with counts, subjects, examples
2. Co-occurrence matrix: which relations appear together on the same source term?
3. Subject distribution per relation
4. Suggested clusters based on semantic similarity patterns
"""

import json
import sys
from collections import Counter, defaultdict

DATA_FILE = "/Users/michaelstewart/Coding/assist/data/chisg/extraction_full_sonnet.json"

print("Loading extraction data...")
with open(DATA_FILE) as f:
    data = json.load(f)

meta = data["extraction_meta"]
controlled = set(meta["controlled_relations"])
extractions = data["extractions"]

print(f"Total definitions: {len(extractions)}")
print(f"Controlled relations: {len(controlled)}")
print()

# ============================================================
# 1. Full relation inventory
# ============================================================
print("=" * 70)
print("1. FULL RELATION INVENTORY")
print("=" * 70)

relation_counts = Counter()
relation_by_subject = defaultdict(Counter)  # relation -> {subject: count}
relation_examples = defaultdict(list)       # relation -> [(source, target, subject)]
relation_cooccurrence = defaultdict(Counter) # relation -> {other_relation: count}
term_relations = defaultdict(set)           # source_term -> set of relations used

for ext in extractions:
    source = ext["source_term"]
    subject = ext.get("subject", "unknown")
    rels_used = set()
    
    for link in ext.get("links", []):
        rel = link.get("relation", "UNKNOWN")
        target = link.get("target", "UNKNOWN")
        
        relation_counts[rel] += 1
        relation_by_subject[rel][subject] += 1
        rels_used.add(rel)
        
        if len(relation_examples[rel]) < 3:
            relation_examples[rel].append((source, target, subject))
    
    term_relations[source] = rels_used
    
    # co-occurrence: which relations appear together?
    for r1 in rels_used:
        for r2 in rels_used:
            if r1 != r2:
                relation_cooccurrence[r1][r2] += 1

# Sort by count
sorted_relations = relation_counts.most_common()

print(f"\nTotal unique relation types: {len(sorted_relations)}")
print(f"  Controlled: {sum(1 for r, _ in sorted_relations if r in controlled)}")
print(f"  Invalid: {sum(1 for r, _ in sorted_relations if r not in controlled)}")
print(f"  Total links: {sum(c for _, c in sorted_relations)}")
print()

# Print all relations with counts, validity, subject breakdown
print(f"{'Relation':<45} {'Count':>6} {'Valid':>5} {'Phys':>5} {'Bio':>5} {'Chem':>5}")
print("-" * 75)
for rel, count in sorted_relations:
    valid = "✅" if rel in controlled else "❌"
    phys = relation_by_subject[rel].get("physics", 0)
    bio = relation_by_subject[rel].get("biology", 0)
    chem = relation_by_subject[rel].get("chemistry", 0)
    print(f"{rel:<45} {count:>6} {valid:>5} {phys:>5} {bio:>5} {chem:>5}")

# ============================================================
# 2. Invalid relations — detailed breakdown
# ============================================================
print()
print("=" * 70)
print("2. INVALID RELATIONS — DETAILED")
print("=" * 70)

invalid_rels = [(r, c) for r, c in sorted_relations if r not in controlled]
print(f"\n{len(invalid_rels)} unique invalid relation types, {sum(c for _, c in invalid_rels)} total links")
print()

# Group invalid relations by likely parent controlled relation
# Use simple heuristic: keyword matching
def suggest_parent(rel):
    """Heuristic: suggest which controlled relation this invalid one maps to."""
    r = rel.lower()
    
    # Causal
    if any(w in r for w in ["cause", "produce", "result", "lead", "trigger", "induce", "generate"]):
        return "causes"
    if any(w in r for w in ["enable", "allow", "permit", "facilitate"]):
        return "enables"
    if any(w in r for w in ["inhibit", "prevent", "block", "suppress", "reduce", "limit", "restrict"]):
        return "inhibits"
    if any(w in r for w in ["determine", "control", "regulate", "govern", "dictate"]):
        return "determines"
    
    # Compositional
    if any(w in r for w in ["compose", "consist", "made of", "contain", "include", "comprise"]):
        return "is composed of"
    if any(w in r for w in ["found in", "located", "present in", "exist in", "occur in"]):
        return "is found in"
    if any(w in r for w in ["property", "characteristic", "attribute", "feature", "trait"]):
        return "has property"
    
    # Taxonomic
    if any(w in r for w in ["type of", "kind of", "form of", "class of", "category"]):
        return "is a type of"
    if any(w in r for w in ["example", "instance", "such as", "like"]):
        return "is an example of"
    if any(w in r for w in ["represent", "symbol", "denote", "signif", "stand for"]):
        return "represents"
    
    # Functional
    if any(w in r for w in ["used for", "function", "purpose", "role", "serve", "employ", "utiliz", "apply"]):
        return "is used for"
    if any(w in r for w in ["value", "measure", "quantit", "amount", "unit"]):
        return "has value"
    
    # Oppositional
    if any(w in r for w in ["contradict", "opposite", "counter", "conflict", "negate"]):
        return "contradicts"
    
    # Passive reversals (is X by)
    if r.startswith("is ") and r.endswith(" by"):
        core = r[3:-3].strip()
        if any(w in core for w in ["cause", "produce", "affect"]):
            return "causes"  # passive reversal
        if any(w in core for w in ["transmit", "transfer", "carry"]):
            return "enables"
        if any(w in core for w in ["form", "creat", "generat"]):
            return "causes"
        return "PASSIVE_REVERSAL"
    
    # Uses / involves
    if any(w in r for w in ["use", "involve", "require", "need", "depend"]):
        return "is used for"
    
    # Transformation
    if any(w in r for w in ["transform", "convert", "change", "turn into", "become"]):
        return "causes"
    
    return "UNMATCHED"

print(f"{'Invalid Relation':<45} {'Count':>5} {'Suggested Parent':<25} {'Examples'}")
print("-" * 120)
parent_mapping = {}
for rel, count in invalid_rels:
    parent = suggest_parent(rel)
    parent_mapping[rel] = parent
    examples = relation_examples[rel][:2]
    ex_str = "; ".join(f"{s}->{t}" for s, t, _ in examples)[:45]
    print(f"{rel:<45} {count:>5} {parent:<25} {ex_str}")

# Summary of mappings
print()
parent_counts = Counter(parent_mapping.values())
print("Suggested parent mapping distribution:")
for parent, count in parent_counts.most_common():
    print(f"  {parent}: {count} invalid types")

# ============================================================
# 3. Co-occurrence patterns
# ============================================================
print()
print("=" * 70)
print("3. CO-OCCURRENCE PATTERNS (which relations appear together?)")
print("=" * 70)

# Only show controlled relations co-occurrence
print(f"\n{'':>20}", end="")
controlled_list = sorted(controlled)
abbrevs = {
    "causes": "caus",
    "enables": "enab",
    "inhibits": "inhi",
    "is composed of": "comp",
    "is a type of": "type",
    "has property": "prop",
    "has value": "valu",
    "is used for": "used",
    "is found in": "fnd ",
    "is an example of": "exmp",
    "determines": "dtrm",
    "represents": "repr",
    "contradicts": "cont"
}
for r in controlled_list:
    print(f"{abbrevs.get(r, r[:4]):>5}", end="")
print()

for r1 in controlled_list:
    print(f"{abbrevs.get(r1, r1[:4]):>20}", end="")
    for r2 in controlled_list:
        if r1 == r2:
            print(f"{'---':>5}", end="")
        else:
            count = relation_cooccurrence[r1].get(r2, 0)
            print(f"{count:>5}", end="")
    print()

# ============================================================
# 4. Subject-specific relation profiles
# ============================================================
print()
print("=" * 70)
print("4. SUBJECT-SPECIFIC RELATION PROFILES")
print("=" * 70)

subjects = ["physics", "biology", "chemistry"]
for subj in subjects:
    print(f"\n--- {subj.upper()} ---")
    subj_total = sum(relation_by_subject[r].get(subj, 0) for r in relation_counts)
    for rel, total_count in sorted_relations[:20]:
        subj_count = relation_by_subject[rel].get(subj, 0)
        if subj_count > 0:
            pct = subj_count / subj_total * 100 if subj_total > 0 else 0
            bar = "█" * int(pct)
            print(f"  {rel:<35} {subj_count:>5} ({pct:>5.1f}%) {bar}")

# ============================================================
# 5. Definitions with many relation types (co-occurrence hubs)
# ============================================================
print()
print("=" * 70)
print("5. TOP CO-OCCURRENCE HUBS (definitions using most distinct relations)")
print("=" * 70)

hub_terms = sorted(term_relations.items(), key=lambda x: len(x[1]), reverse=True)[:20]
for term, rels in hub_terms:
    print(f"  {term:<40} {len(rels)} relations: {', '.join(sorted(rels))}")

# ============================================================
# 6. Output structured data for thesaurus builder
# ============================================================
output = {
    "analysis_date": "2026-02-22",
    "total_links": sum(c for _, c in sorted_relations),
    "total_unique_relations": len(sorted_relations),
    "controlled_count": sum(1 for r, _ in sorted_relations if r in controlled),
    "invalid_count": sum(1 for r, _ in sorted_relations if r not in controlled),
    "relation_inventory": [
        {
            "relation": rel,
            "count": count,
            "is_controlled": rel in controlled,
            "suggested_parent": parent_mapping.get(rel, rel if rel in controlled else "UNMATCHED"),
            "subjects": dict(relation_by_subject[rel]),
            "examples": [{"source": s, "target": t, "subject": subj} for s, t, subj in relation_examples[rel][:3]]
        }
        for rel, count in sorted_relations
    ],
    "cooccurrence": {
        r1: {r2: c for r2, c in relation_cooccurrence[r1].most_common(5)}
        for r1 in controlled
    }
}

output_path = "/Users/michaelstewart/Coding/assist/data/chisg/relation_analysis.json"
with open(output_path, "w") as f:
    json.dump(output, f, indent=2)
print(f"\n\nStructured analysis written to: {output_path}")
print("This file will be used by the thesaurus builder script.")
