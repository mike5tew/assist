#!/usr/bin/env python3
"""Quick analysis of extraction output."""
import json, sys

path = sys.argv[1] if len(sys.argv) > 1 else "data/chisg/extraction_test_s10.json"
d = json.load(open(path))
s = d["stats"]
n = s.get("definitions_processed", 0)

print(f"Processed: {n} defs")
print(f"Links: {s['total_links']}")
print(f"Avg: {s['total_links']/max(n,1):.1f} links/def")
print(f"Invalid: {s['invalid_relations']}")
print(f"Zero-link: {s['zero_link_definitions']}")
print()

CV = ["causes","enables","inhibits","is composed of","is a type of","has property",
      "has value","is used for","is found in","is an example of","determines",
      "represents","contradicts"]

valid_total = sum(v for k,v in s["relation_counts"].items() if k in CV)
invalid_total = sum(v for k,v in s["relation_counts"].items() if k not in CV)
total = valid_total + invalid_total
print(f"Valid relations: {valid_total} ({valid_total*100//max(total,1)}%)")
print(f"Invalid relations: {invalid_total} ({invalid_total*100//max(total,1)}%)")

print("\nControlled vocabulary usage:")
for k,v in sorted(s["relation_counts"].items(), key=lambda x:-x[1]):
    if k in CV:
        print(f"  ✓ {k}: {v}")

print("\nInvalid relation types (candidates for vocabulary expansion):")
for k,v in sorted(s["relation_counts"].items(), key=lambda x:-x[1]):
    if k not in CV:
        print(f"  ✗ {k}: {v}")
