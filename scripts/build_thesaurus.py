#!/usr/bin/env python3
"""
Build the Scoped Semantic Thesaurus from CHISG extraction data.

Implements "Continuous Relational Abstraction":
- Every relation (controlled + invalid) is placed in a hierarchical group
- Each has a similarity score (0.0–1.0) relative to its group canonical
- Queries at threshold T include all relations with score >= T
- Direction flag enables reverse-relation recovery

Hierarchy:
  Level 0: Root group (e.g. "causal")
  Level 1: Family (e.g. "causal.direct", "causal.enabling")
  Level 2: Individual relation entries

Storage: MongoDB → esp_organizer.chisg_thesaurus_groups + chisg_thesaurus_entries
"""

import json
import sys
from datetime import datetime
from pymongo import MongoClient

# ============================================================
# THESAURUS HIERARCHY DEFINITION
# ============================================================
# Derived from relation_analysis.json co-occurrence + subject profiles

HIERARCHY = [
    # Root groups (level 0)
    {
        "group_id": "causal",
        "level": 0,
        "parent_id": None,
        "label": "Causal Relations",
        "description": "Relations about cause, effect, production, and transformation. "
                       "Physics: 8.9% causes, 8.6% enables. Biology: 8.3% causes, 14.9% enables. "
                       "This is the largest root group by variety of invalid synonyms."
    },
    {
        "group_id": "structural",
        "level": 0,
        "parent_id": None,
        "label": "Structural Relations",
        "description": "Relations about composition, location, properties, and values. "
                       "Dominates all subjects: has property is #1 everywhere (16-25%). "
                       "Chemistry especially strong on composition (12.9%)."
    },
    {
        "group_id": "taxonomic",
        "level": 0,
        "parent_id": None,
        "label": "Taxonomic Relations",
        "description": "Relations about classification, identity, exemplification, and representation. "
                       "is a type of is #2 or #3 across all subjects. "
                       "Co-occurs strongly with has property (241 times)."
    },
    {
        "group_id": "functional",
        "level": 0,
        "parent_id": None,
        "label": "Functional Relations",
        "description": "Relations about purpose, use, requirements, and interactions. "
                       "is used for ranges from 6.6% (physics) to 11% (chemistry). "
                       "Strong co-occurrence with is a type of (90) and has property (106)."
    },
    {
        "group_id": "oppositional",
        "level": 0,
        "parent_id": None,
        "label": "Oppositional Relations",
        "description": "Relations about contradiction, conflict, and negation. "
                       "Smallest group: contradicts at 48 total (mostly biology 31). "
                       "Low co-occurrence with everything — these are semantic outliers."
    },

    # Families (level 1)
    # -- Causal --
    {
        "group_id": "causal.direct",
        "level": 1,
        "parent_id": "causal",
        "label": "Direct Causation",
        "description": "One thing directly causing, producing, or transforming another. "
                       "The core causal relation. 559 controlled 'causes' links."
    },
    {
        "group_id": "causal.enabling",
        "level": 1,
        "parent_id": "causal",
        "label": "Enabling / Facilitating",
        "description": "One thing enabling, transmitting, carrying, or facilitating another. "
                       "Second most common relation overall (782 links). Biology-heavy (433/782)."
    },
    {
        "group_id": "causal.inhibiting",
        "level": 1,
        "parent_id": "causal",
        "label": "Inhibiting / Preventing",
        "description": "One thing blocking, preventing, or reducing another. "
                       "Biology-dominated (58/70). Semantic opposite of enabling."
    },
    {
        "group_id": "causal.determining",
        "level": 1,
        "parent_id": "causal",
        "label": "Controlling / Determining",
        "description": "One thing controlling, regulating, or determining another. "
                       "Relatively rare (80 links). Distributed across all subjects."
    },

    # -- Structural --
    {
        "group_id": "structural.composition",
        "level": 1,
        "parent_id": "structural",
        "label": "Composition / Containment",
        "description": "What things are made of, contain, or are formed from. "
                       "630 controlled 'is composed of' links. Chemistry-heavy (197/630)."
    },
    {
        "group_id": "structural.location",
        "level": 1,
        "parent_id": "structural",
        "label": "Location / Context",
        "description": "Where things are found, stored, occur, or originate. "
                       "349 links. Biology-dominant (219/349)."
    },
    {
        "group_id": "structural.property",
        "level": 1,
        "parent_id": "structural",
        "label": "Properties / Attributes",
        "description": "Attributes, features, symptoms, or characteristics. "
                       "LARGEST single relation: 1293 'has property' links. #1 everywhere."
    },
    {
        "group_id": "structural.value",
        "level": 1,
        "parent_id": "structural",
        "label": "Values / Measurements",
        "description": "Quantitative properties, measurements, and units. "
                       "237 links. Physics-heavy (156/237)."
    },

    # -- Taxonomic --
    {
        "group_id": "taxonomic.classification",
        "level": 1,
        "parent_id": "taxonomic",
        "label": "Classification / Typing",
        "description": "Category, type, or class membership. "
                       "718 links — third most common relation overall."
    },
    {
        "group_id": "taxonomic.exemplification",
        "level": 1,
        "parent_id": "taxonomic",
        "label": "Exemplification / Evidence",
        "description": "Examples, instances, evidence, and predictions. "
                       "338 links. Biology: 167, Physics: 111, Chemistry: 60."
    },
    {
        "group_id": "taxonomic.representation",
        "level": 1,
        "parent_id": "taxonomic",
        "label": "Representation / Symbolism",
        "description": "Symbolic, representational, or descriptive relations. "
                       "358 links. Physics-heavy (189/358) — formulas, units, models."
    },

    # -- Functional --
    {
        "group_id": "functional.application",
        "level": 1,
        "parent_id": "functional",
        "label": "Application / Usage",
        "description": "How things are used, what they require, what they achieve. "
                       "575 controlled 'is used for' links."
    },
    {
        "group_id": "functional.interaction",
        "level": 1,
        "parent_id": "functional",
        "label": "Interaction / Connection",
        "description": "How things interact, bind, connect, or respond to each other. "
                       "Mostly biology (cellular/molecular interactions)."
    },

    # -- Oppositional --
    {
        "group_id": "oppositional.contradiction",
        "level": 1,
        "parent_id": "oppositional",
        "label": "Contradiction / Conflict",
        "description": "Direct semantic contradiction or opposition. "
                       "48 links total. Biology: 31, Physics: 9, Chemistry: 8."
    },
]

# ============================================================
# RELATION → GROUP MAPPING WITH SCORES
# ============================================================
# Score semantics:
#   1.00       = canonical controlled relation
#   0.85-0.95  = direct synonym, same direction
#   0.75-0.84  = very close synonym or reverse of canonical
#   0.60-0.74  = same family, weaker synonym
#   0.45-0.59  = same family, loose or domain-specific
#   0.30-0.44  = very loose, placed for recovery

# Format: (relation, group_id, score, direction, notes)
# direction: "forward" = same direction as canonical
#            "reverse" = swapped source↔target from canonical
#            "neutral" = direction doesn't meaningfully change

RELATION_MAP = [
    # ============ CAUSAL.DIRECT (canonical: causes) ============
    ("causes",              "causal.direct",      1.00, "forward", "Canonical controlled relation"),
    ("produces",            "causal.direct",      0.90, "forward", "Production is direct causation"),
    ("results in",          "causal.direct",      0.90, "forward", "Direct outcome"),
    ("converts",            "causal.direct",      0.85, "forward", "Transformation is specific causation"),
    ("releases",            "causal.direct",      0.85, "forward", "Release is a form of production"),
    ("produced",            "causal.direct",      0.85, "forward", "Past tense of produces"),
    ("secretes",            "causal.direct",      0.85, "forward", "Biological production/release"),
    ("affects",             "causal.direct",      0.80, "forward", "Weaker/unspecified causation"),
    ("is produced",         "causal.direct",      0.75, "reverse", "Passive of produces"),
    ("is caused by",        "causal.direct",      0.75, "reverse", "Reverse of causes"),
    ("is produced by",      "causal.direct",      0.75, "reverse", "Reverse of produces"),
    ("is formed by",        "causal.direct",      0.75, "reverse", "Reverse of forms (≈ causes)"),
    ("is generated by",     "causal.direct",      0.75, "reverse", "Reverse of generates (≈ causes)"),
    ("is converted to",     "causal.direct",      0.75, "forward", "Transformation result"),
    ("is emitted by",       "causal.direct",      0.75, "reverse", "Reverse of emits (≈ produces)"),
    ("is emitted from",     "causal.direct",      0.70, "reverse", "Emission source"),
    ("is secreted by",      "causal.direct",      0.75, "reverse", "Reverse of secretes"),
    ("is given off by",     "causal.direct",      0.70, "reverse", "Reverse of gives off (≈ produces)"),
    ("is a result of",      "causal.direct",      0.70, "reverse", "Reverse of results in"),
    ("is a consequence of", "causal.direct",      0.70, "reverse", "Reverse of causes"),
    ("is produced in",      "causal.direct",      0.70, "reverse", "Location of production"),
    ("is released by",      "causal.direct",      0.70, "reverse", "Reverse of releases"),
    ("is released from",    "causal.direct",      0.70, "reverse", "Source of release"),
    ("is released during",  "causal.direct",      0.65, "neutral", "Temporal context of release"),
    ("is formed from",      "causal.direct",      0.75, "reverse", "Material cause"),
    ("is replaced by",      "causal.direct",      0.60, "reverse", "Succession/replacement"),
    ("can reproduce",       "causal.direct",      0.55, "forward", "Reproduction is biological causation"),
    ("can reproduce with",  "causal.direct",      0.55, "forward", "Reproduction with partner"),
    ("become resistant to", "causal.direct",      0.50, "forward", "Resistance is a causal outcome"),
    ("increases",           "causal.direct",      0.70, "forward", "Quantitative causation"),
    ("increases speed",     "causal.direct",      0.65, "forward", "Specific quantitative causation"),

    # ============ CAUSAL.ENABLING (canonical: enables) ============
    ("enables",             "causal.enabling",    1.00, "forward", "Canonical controlled relation"),
    ("is transmitted by",   "causal.enabling",    0.80, "reverse", "Transmission is enabling movement"),
    ("is enabled by",       "causal.enabling",    0.75, "reverse", "Reverse of enables"),
    ("carries",             "causal.enabling",    0.80, "forward", "Carrying enables delivery"),
    ("supplies",            "causal.enabling",    0.80, "forward", "Supply enables use"),
    ("transports",          "causal.enabling",    0.80, "forward", "Transport enables access"),
    ("transports to",       "causal.enabling",    0.75, "forward", "Directed transport"),
    ("conducts",            "causal.enabling",    0.75, "forward", "Conduction enables flow"),
    ("travels to",          "causal.enabling",    0.70, "forward", "Movement enables reaching"),
    ("powers",              "causal.enabling",    0.80, "forward", "Powering enables function"),
    ("is mediated by",      "causal.enabling",    0.70, "reverse", "Mediation is enabling mechanism"),
    ("is inherited",        "causal.enabling",    0.65, "neutral", "Genetic transmission"),
    ("is inherited from",   "causal.enabling",    0.65, "reverse", "Source of genetic inheritance"),
    ("receives",            "causal.enabling",    0.70, "reverse", "Receiving what was transported"),
    ("carries blood away from", "causal.enabling", 0.70, "forward", "Specific biological transport"),
    ("is absorbed by",      "causal.enabling",    0.65, "reverse", "Absorption as receiving/enabling"),

    # ============ CAUSAL.INHIBITING (canonical: inhibits) ============
    ("inhibits",            "causal.inhibiting",  1.00, "forward", "Canonical controlled relation"),
    ("neutralizes",         "causal.inhibiting",  0.80, "forward", "Neutralization prevents action"),
    ("kill",                "causal.inhibiting",  0.75, "forward", "Killing is ultimate inhibition"),
    ("overcomes",           "causal.inhibiting",  0.70, "forward", "Overcoming resistance/pathogen"),
    ("cannot use",          "causal.inhibiting",  0.65, "forward", "Inability — functional inhibition"),
    ("fails to secrete sufficient", "causal.inhibiting", 0.60, "forward", "Failure of production"),

    # ============ CAUSAL.DETERMINING (canonical: determines) ============
    ("determines",          "causal.determining", 1.00, "forward", "Canonical controlled relation"),
    ("controls",            "causal.determining", 0.90, "forward", "Control is determining"),
    ("regulates",           "causal.determining", 0.90, "forward", "Regulation is determining within bounds"),
    ("is determined by",    "causal.determining", 0.75, "reverse", "Reverse of determines"),
    ("monitors",            "causal.determining", 0.70, "forward", "Monitoring supports determining"),
    ("coordinates",         "causal.determining", 0.70, "forward", "Coordination is multi-actor determining"),
    ("is done when",        "causal.determining", 0.55, "neutral", "Conditional trigger — loose determining"),

    # ============ STRUCTURAL.COMPOSITION (canonical: is composed of) ============
    ("is composed of",      "structural.composition", 1.00, "forward", "Canonical controlled relation"),
    ("contains",            "structural.composition", 0.90, "forward", "Containment ≈ composition"),
    ("is made from",        "structural.composition", 0.85, "reverse", "Material origin"),
    ("combined",            "structural.composition", 0.75, "forward", "Combination into composite"),
    ("are held together by","structural.composition", 0.75, "reverse", "Binding force for composition"),
    ("is shared between",   "structural.composition", 0.65, "neutral", "Shared component"),
    ("is derived from",     "structural.composition", 0.70, "reverse", "Derivation from material"),
    ("are made by",         "structural.composition", 0.70, "reverse", "Production of components"),
    ("is bound by",         "structural.composition", 0.70, "reverse", "Bound into structure"),
    ("absorbs",             "structural.composition", 0.60, "forward", "Absorption into structure"),
    ("is coded for by",     "structural.composition", 0.65, "reverse", "Genetic coding → composition"),
    ("is connected by",     "structural.composition", 0.65, "reverse", "Connection between parts"),
    ("is added to",         "structural.composition", 0.65, "forward", "Addition to composition"),

    # ============ STRUCTURAL.LOCATION (canonical: is found in) ============
    ("is found in",         "structural.location",    1.00, "forward", "Canonical controlled relation"),
    ("is stored in",        "structural.location",    0.90, "forward", "Storage is persistent location"),
    ("occurs in",           "structural.location",    0.85, "forward", "Occurrence location"),
    ("occurs during",       "structural.location",    0.75, "forward", "Temporal location"),
    ("is found near",       "structural.location",    0.80, "forward", "Proximity"),
    ("is between",          "structural.location",    0.70, "forward", "Positional relationship"),
    ("extends throughout",  "structural.location",    0.70, "forward", "Spatial extent"),
    ("are concentrated in", "structural.location",    0.80, "forward", "High-density location"),
    ("covers",              "structural.location",    0.70, "forward", "Surface location/coverage"),
    ("is formed in",        "structural.location",    0.70, "forward", "Location of formation"),
    ("landed on",           "structural.location",    0.55, "forward", "Specific locational event"),
    ("has been in orbit around", "structural.location", 0.50, "forward", "Orbital location"),
    ("is moved out of",     "structural.location",    0.65, "reverse", "Movement from location"),

    # ============ STRUCTURAL.PROPERTY (canonical: has property) ============
    ("has property",        "structural.property",    1.00, "forward", "Canonical controlled relation"),
    ("has symptom",         "structural.property",    0.85, "forward", "Symptom is a medical property"),
    ("is a property of",    "structural.property",    0.75, "reverse", "Reverse of has property"),
    ("is defined as",       "structural.property",    0.55, "neutral", "Definition captures properties (R1 violation)"),
    ("is defined by",       "structural.property",    0.55, "reverse", "Definitional property"),
    ("describes",           "structural.property",    0.60, "forward", "Description captures properties"),
    ("is specialized for",  "structural.property",    0.60, "forward", "Specialization is a property"),
    ("is often associated with", "structural.property", 0.55, "forward", "Association is weak property link"),
    ("is based on",         "structural.property",    0.55, "neutral", "Foundation is a property"),

    # ============ STRUCTURAL.VALUE (canonical: has value) ============
    ("has value",           "structural.value",       1.00, "forward", "Canonical controlled relation"),
    ("is measured in",      "structural.value",       0.90, "forward", "Measurement unit"),
    ("is a measure of",     "structural.value",       0.85, "reverse", "What a measurement represents"),
    ("is a value",          "structural.value",       0.90, "forward", "Direct value statement"),
    ("measures",            "structural.value",       0.80, "forward", "Measuring function"),
    ("is more than",        "structural.value",       0.55, "forward", "Comparative value"),

    # ============ TAXONOMIC.CLASSIFICATION (canonical: is a type of) ============
    ("is a type of",        "taxonomic.classification", 1.00, "forward", "Canonical controlled relation"),
    ("are arranged by",     "taxonomic.classification", 0.65, "neutral", "Arrangement by category"),
    ("distinguishes from",  "taxonomic.classification", 0.60, "neutral", "Differentiation within taxonomy"),

    # ============ TAXONOMIC.EXEMPLIFICATION (canonical: is an example of) ============
    ("is an example of",    "taxonomic.exemplification", 1.00, "forward", "Canonical controlled relation"),
    ("is evidence for",     "taxonomic.exemplification", 0.65, "forward", "Evidence is empirical exemplification"),
    ("correctly predicted", "taxonomic.exemplification", 0.50, "forward", "Prediction confirmed by example"),

    # ============ TAXONOMIC.REPRESENTATION (canonical: represents) ============
    ("represents",          "taxonomic.representation", 1.00, "forward", "Canonical controlled relation"),
    ("is displayed on",     "taxonomic.representation", 0.70, "forward", "Visual representation"),
    ("is conserved in",     "taxonomic.representation", 0.55, "neutral", "Conservation represents invariance"),
    ("proposed",            "taxonomic.representation", 0.50, "neutral", "Proposing a representation/model"),
    ("discovered",          "taxonomic.representation", 0.50, "neutral", "Discovery reveals representation"),
    ("performed experiment","taxonomic.representation", 0.40, "neutral", "Experiment represents method"),
    ("proved",              "taxonomic.representation", 0.45, "neutral", "Proof represents verification"),

    # ============ FUNCTIONAL.APPLICATION (canonical: is used for) ============
    ("is used for",         "functional.application",   1.00, "forward", "Canonical controlled relation"),
    ("uses",                "functional.application",   0.90, "reverse", "Reverse perspective of 'is used for'"),
    ("requires",            "functional.application",   0.85, "forward", "Requirement for function"),
    ("is required for",     "functional.application",   0.85, "forward", "Requirement relationship"),
    ("depends on",          "functional.application",   0.80, "forward", "Dependency for function"),
    ("is needed for",       "functional.application",   0.80, "forward", "Necessity for function"),
    ("is needed by",        "functional.application",   0.75, "reverse", "Reverse of needed"),
    ("involves",            "functional.application",   0.70, "forward", "Involvement in process"),
    ("is involved in",      "functional.application",   0.70, "forward", "Participation"),
    ("does work by",        "functional.application",   0.75, "forward", "Mechanism of work"),
    ("is achieved by",      "functional.application",   0.70, "reverse", "Achievement through method"),
    ("has function",        "functional.application",   0.85, "forward", "Direct function statement"),
    ("is obtained by",      "functional.application",   0.65, "reverse", "Method of obtaining"),
    ("is shone through",    "functional.application",   0.50, "neutral", "Light application method"),
    ("extracts nutrients from", "functional.application", 0.65, "forward", "Specific biological function"),
    ("digests",             "functional.application",   0.65, "forward", "Digestive function"),
    ("translates",          "functional.application",   0.65, "forward", "Translation function"),
    ("is treated with",     "functional.application",   0.70, "reverse", "Medical application"),
    ("reacts with",         "functional.application",   0.70, "forward", "Chemical interaction/use"),
    ("involved",            "functional.application",   0.60, "neutral", "Past tense involvement"),
    ("looking for evidence of", "functional.application", 0.40, "forward", "Search function"),
    ("scan the stars",      "functional.application",   0.40, "forward", "Observation function"),
    ("are messages to",     "functional.application",   0.45, "forward", "Communication function"),
    ("funds",               "functional.application",   0.45, "forward", "Financial support function"),
    ("will not reach another star for", "functional.application", 0.35, "neutral", "Temporal constraint"),
    ("left gaps for undiscovered elements", "functional.application", 0.35, "neutral", "Historical gap"),
    ("was found",           "functional.application",   0.40, "neutral", "Discovery as outcome"),

    # ============ FUNCTIONAL.INTERACTION (canonical: interacts with — new effective canonical) ============
    ("connects",            "functional.interaction",   0.85, "forward", "Physical connection"),
    ("binds to",            "functional.interaction",   0.90, "forward", "Molecular binding"),
    ("binds",               "functional.interaction",   0.90, "forward", "Binding interaction"),
    ("interacts with",      "functional.interaction",   0.90, "forward", "General interaction"),
    ("responds to",         "functional.interaction",   0.75, "forward", "Response interaction"),
    ("is a response to",    "functional.interaction",   0.75, "reverse", "Reverse of responds to"),
    ("engulfs",             "functional.interaction",   0.70, "forward", "Phagocytic interaction"),
    ("targets",             "functional.interaction",   0.75, "forward", "Targeting interaction"),
    ("infects",             "functional.interaction",   0.70, "forward", "Pathogenic interaction"),
    ("is connected by",     "functional.interaction",   0.70, "reverse", "Reverse of connects"),

    # ============ OPPOSITIONAL.CONTRADICTION (canonical: contradicts) ============
    ("contradicts",         "oppositional.contradiction", 1.00, "forward", "Canonical controlled relation"),

    # ============ UNMAPPED RECOVERY ============
    ("is excreted by",      "causal.direct",      0.65, "reverse", "Reverse of excretes (≈ produces/releases)"),
]


# ============================================================
# BUILD AND LOAD
# ============================================================

def build_entries(extraction_data):
    """Enrich each relation entry with frequency data from extraction."""
    meta = extraction_data["extraction_meta"]
    controlled = set(meta["controlled_relations"])
    
    # Count actual frequencies
    freq = {}
    subj_dist = {}
    examples = {}
    for ext in extraction_data["extractions"]:
        subject = ext.get("subject", "unknown")
        for link in ext.get("links", []):
            rel = link.get("relation", "UNKNOWN")
            target = link.get("target", "UNKNOWN")
            freq[rel] = freq.get(rel, 0) + 1
            if rel not in subj_dist:
                subj_dist[rel] = {}
            subj_dist[rel][subject] = subj_dist[rel].get(subject, 0) + 1
            if rel not in examples or len(examples[rel]) < 3:
                if rel not in examples:
                    examples[rel] = []
                examples[rel].append({
                    "source": ext["source_term"],
                    "target": target,
                    "subject": subject
                })
    
    entries = []
    mapped_relations = set()
    
    for (relation, group_id, score, direction, notes) in RELATION_MAP:
        root_group = group_id.split(".")[0]
        family = group_id
        is_controlled = relation in controlled
        
        entry = {
            "relation": relation,
            "group_id": group_id,
            "root_group": root_group,
            "score": score,
            "direction": direction,
            "is_controlled": is_controlled,
            "is_canonical": score == 1.0 and is_controlled,
            "frequency": freq.get(relation, 0),
            "subjects": subj_dist.get(relation, {}),
            "examples": examples.get(relation, [])[:3],
            "notes": notes,
            "created_at": datetime.utcnow().isoformat()
        }
        entries.append(entry)
        mapped_relations.add(relation)
    
    # Check for any relations in extraction that we missed
    all_extraction_relations = set(freq.keys())
    unmapped = all_extraction_relations - mapped_relations
    if unmapped:
        print(f"\n⚠️  {len(unmapped)} relations in extraction data not mapped in thesaurus:")
        for rel in sorted(unmapped):
            print(f"    '{rel}' (frequency: {freq[rel]})")
    
    return entries, unmapped


def print_thesaurus_summary(entries):
    """Print a human-readable summary of the thesaurus."""
    print("\n" + "=" * 70)
    print("THESAURUS SUMMARY")
    print("=" * 70)
    
    # Group entries by group_id
    by_group = {}
    for e in entries:
        gid = e["group_id"]
        if gid not in by_group:
            by_group[gid] = []
        by_group[gid].append(e)
    
    total_controlled = sum(1 for e in entries if e["is_controlled"])
    total_invalid = sum(1 for e in entries if not e["is_controlled"])
    total_freq = sum(e["frequency"] for e in entries)
    recovered_freq = sum(e["frequency"] for e in entries if not e["is_controlled"])
    
    print(f"\nTotal entries: {len(entries)}")
    print(f"  Controlled: {total_controlled}")
    print(f"  Recovered invalids: {total_invalid}")
    print(f"  Total links covered: {total_freq}")
    print(f"  Recovered invalid links: {recovered_freq}")
    
    # By root group
    root_groups = {}
    for e in entries:
        rg = e["root_group"]
        if rg not in root_groups:
            root_groups[rg] = {"count": 0, "freq": 0, "families": set()}
        root_groups[rg]["count"] += 1
        root_groups[rg]["freq"] += e["frequency"]
        root_groups[rg]["families"].add(e["group_id"])
    
    print(f"\n{'Root Group':<20} {'Entries':>8} {'Families':>8} {'Links':>8}")
    print("-" * 50)
    for rg in sorted(root_groups.keys()):
        info = root_groups[rg]
        print(f"{rg:<20} {info['count']:>8} {len(info['families']):>8} {info['freq']:>8}")
    
    # Detailed by family
    print(f"\n{'Group':<30} {'#Rel':>5} {'Links':>6} {'Score Range':<12} {'Canonical'}")
    print("-" * 80)
    for hier in HIERARCHY:
        if hier["level"] != 1:
            continue
        gid = hier["group_id"]
        members = by_group.get(gid, [])
        if not members:
            continue
        scores = [m["score"] for m in members]
        links = sum(m["frequency"] for m in members)
        canon = [m["relation"] for m in members if m.get("is_canonical")]
        canon_str = canon[0] if canon else members[0]["relation"]
        print(f"{gid:<30} {len(members):>5} {links:>6} {min(scores):.2f}–{max(scores):.2f}    {canon_str}")
        
        # Show members sorted by score
        for m in sorted(members, key=lambda x: -x["score"]):
            marker = "★" if m["is_canonical"] else ("●" if m["is_controlled"] else "○")
            dir_str = {"forward": "→", "reverse": "←", "neutral": "↔"}[m["direction"]]
            print(f"  {marker} {m['relation']:<40} {m['score']:.2f} {dir_str} freq={m['frequency']}")


def load_to_mongodb(hierarchy, entries):
    """Load thesaurus into MongoDB."""
    client = MongoClient("mongodb://admin:password123@localhost:27018/?authSource=admin")
    db = client["esp_organizer"]
    
    # Drop existing collections
    db.drop_collection("chisg_thesaurus_groups")
    db.drop_collection("chisg_thesaurus_entries")
    
    # Insert hierarchy
    for h in hierarchy:
        h["created_at"] = datetime.utcnow().isoformat()
    result = db.chisg_thesaurus_groups.insert_many(hierarchy)
    print(f"\n✅ Inserted {len(result.inserted_ids)} hierarchy groups")
    
    # Insert entries
    result = db.chisg_thesaurus_entries.insert_many(entries)
    print(f"✅ Inserted {len(result.inserted_ids)} thesaurus entries")
    
    # Create indexes
    db.chisg_thesaurus_entries.create_index("relation")
    db.chisg_thesaurus_entries.create_index("group_id")
    db.chisg_thesaurus_entries.create_index("root_group")
    db.chisg_thesaurus_entries.create_index("score")
    db.chisg_thesaurus_entries.create_index([("group_id", 1), ("score", -1)])
    db.chisg_thesaurus_entries.create_index([("root_group", 1), ("score", -1)])
    print("✅ Created indexes on relation, group_id, root_group, score")
    
    # Create metadata document
    db.chisg_thesaurus_meta.drop()
    db.chisg_thesaurus_meta.insert_one({
        "version": "1.0",
        "created_at": datetime.utcnow().isoformat(),
        "source": "extraction_full_sonnet.json",
        "total_groups": len(hierarchy),
        "total_entries": len(entries),
        "root_groups": list(set(e["root_group"] for e in entries)),
        "score_semantics": {
            "1.00": "Canonical controlled relation",
            "0.85-0.95": "Direct synonym, same direction",
            "0.75-0.84": "Very close synonym or reverse of canonical",
            "0.60-0.74": "Same family, weaker synonym",
            "0.45-0.59": "Same family, loose or domain-specific",
            "0.30-0.44": "Very loose, placed for recovery"
        },
        "threshold_guide": {
            "1.0": "Exact controlled relation only",
            "0.8": "Close synonyms + reverse forms",
            "0.6": "Whole family including weak synonyms",
            "0.3": "Everything in group (maximum recall)"
        }
    })
    print("✅ Created metadata document")
    
    return client


def test_threshold_queries(db):
    """Demonstrate threshold-based querying."""
    print("\n" + "=" * 70)
    print("THRESHOLD QUERY DEMONSTRATIONS")
    print("=" * 70)
    
    coll = db.chisg_thesaurus_entries
    
    test_cases = [
        ("causes", [1.0, 0.8, 0.6, 0.3]),
        ("enables", [1.0, 0.8, 0.6]),
        ("is composed of", [1.0, 0.8, 0.6]),
        ("has property", [1.0, 0.8, 0.6]),
    ]
    
    for canon_rel, thresholds in test_cases:
        # Find the group for this canonical relation
        entry = coll.find_one({"relation": canon_rel, "is_canonical": True})
        if not entry:
            entry = coll.find_one({"relation": canon_rel})
        if not entry:
            print(f"\n⚠️ '{canon_rel}' not found")
            continue
        
        group_id = entry["group_id"]
        
        print(f"\n--- Query: '{canon_rel}' (group: {group_id}) ---")
        
        for thresh in thresholds:
            results = list(coll.find(
                {"group_id": group_id, "score": {"$gte": thresh}},
                {"relation": 1, "score": 1, "direction": 1, "frequency": 1, "_id": 0}
            ).sort("score", -1))
            
            total_links = sum(r.get("frequency", 0) for r in results)
            rels = [f"{r['relation']}({r['score']:.2f})" for r in results]
            print(f"  threshold ≥ {thresh}: {len(results)} relations, {total_links} links")
            print(f"    {', '.join(rels)}")
    
    # Cross-group query: "All causal relations"
    print(f"\n--- Cross-group: All 'causal' root relations at threshold ≥ 0.6 ---")
    results = list(coll.find(
        {"root_group": "causal", "score": {"$gte": 0.6}},
        {"relation": 1, "score": 1, "group_id": 1, "frequency": 1, "_id": 0}
    ).sort([("group_id", 1), ("score", -1)]))
    
    current_group = None
    total_links = 0
    for r in results:
        if r["group_id"] != current_group:
            current_group = r["group_id"]
            print(f"  [{current_group}]")
        total_links += r.get("frequency", 0)
        print(f"    {r['relation']:<40} {r['score']:.2f}  freq={r.get('frequency', 0)}")
    print(f"  TOTAL: {len(results)} relations covering {total_links} links")


DATA_FILE = "/Users/michaelstewart/Coding/assist/data/chisg/extraction_full_sonnet.json"

def main():
    # Load extraction data
    print("Loading extraction data...")
    with open(DATA_FILE) as f:
        extraction_data = json.load(f)
    
    # Build entries
    entries, unmapped = build_entries(extraction_data)
    
    # Print summary
    print_thesaurus_summary(entries)
    
    # Load to MongoDB
    print("\n" + "=" * 70)
    print("LOADING TO MONGODB")
    print("=" * 70)
    
    client = load_to_mongodb(HIERARCHY, entries)
    db = client["esp_organizer"]
    
    # Run test queries
    test_threshold_queries(db)
    
    client.close()
    print("\n✅ Thesaurus build complete!")


if __name__ == "__main__":
    main()
