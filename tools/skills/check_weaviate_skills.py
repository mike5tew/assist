#!/usr/bin/env python3
"""Check which reception sticker-book skills exist in the Weaviate CHISG graph."""

import json
import urllib.request

WEAVIATE_URL = "http://localhost:8081/v1/graphql"

RECEPTION_SKILLS = [
    "Attention",
    "Working Memory",
    "Fine Motor Skills",
    "Following Protocols",
    "Social Confidence",
    "Speaking",
    "Vocabulary",
    "Non-verbal communication",
    "Self-regulation",
    "Empathy",
]

def query_weaviate(query_str):
    payload = json.dumps({"query": query_str}).encode()
    req = urllib.request.Request(WEAVIATE_URL, data=payload, headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req) as resp:
        return json.loads(resp.read())

def main():
    # Fetch all CHISG skills
    result = query_weaviate("""{ Get { CHISGElement(limit: 700) { name domain description _additional { id } } } }""")
    skills = result["data"]["Get"]["CHISGElement"]
    print(f"Total CHISG skills in Weaviate: {len(skills)}")
    print()

    # Build lookup (case-insensitive)
    by_lower = {}
    for s in skills:
        by_lower[s["name"].lower()] = s

    print("=" * 60)
    print("RECEPTION STICKER BOOK — SKILL CHECK")
    print("=" * 60)

    found = []
    missing = []

    for target in RECEPTION_SKILLS:
        key = target.lower()
        if key in by_lower:
            s = by_lower[key]
            found.append((target, s))
            uid = s["_additional"]["id"][:12]
            domain = s.get("domain", "")
            desc = (s.get("description", "") or "")[:60]
            print(f"  ✅  {target:30s}  id={uid}..  domain={domain}")
            if desc:
                print(f"      desc: {desc}")
        else:
            # Try partial / substring matches
            partials = []
            for name_lower, s in by_lower.items():
                if key in name_lower or name_lower in key:
                    partials.append(s["name"])
            if partials:
                missing.append((target, partials))
                print(f"  ⚠️   {target:30s}  NOT EXACT — close matches: {partials[:5]}")
            else:
                missing.append((target, []))
                print(f"  ❌  {target:30s}  NOT FOUND")

    print()
    print(f"Found exact:  {len(found)}/10")
    print(f"Missing:      {len(missing)}/10")

    if missing:
        print()
        print("Skills that need integration:")
        for name, partials in missing:
            if partials:
                print(f"  - {name}  (possible matches: {partials})")
            else:
                print(f"  - {name}  (no similar skills found)")

if __name__ == "__main__":
    main()
