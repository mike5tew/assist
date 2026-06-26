#!/usr/bin/env python3
"""
Log session changes from 2026-05-17 to Weaviate Documentation class.
Changes:
  - CHISG graph cleanup (deletions + new links + wrong roots fixed)
  - LAO Mobile version bump + EAS appVersionSource fix
  - IAP App Store Connect workflow diagnosis

Usage: python3 log_changes_20260517.py
"""

import weaviate
from weaviate.classes.init import AdditionalConfig, Timeout
import hashlib
from datetime import datetime, timezone

WEAVIATE_HOST = "localhost"
WEAVIATE_PORT = 8088
WEAVIATE_GRPC_PORT = 50052
DOC_CLASS = "Documentation"

CHANGES = [
    {
        "title": "CHISG Skills Graph Cleanup — Deletions and Link Fixes (2026-05-17)",
        "project": "assist",
        "section_path": "CHISG > Skills Graph > Maintenance",
        "file_path": "chisg/maintenance/2026-05-17-graph-cleanup.md",
        "content": """## CHISG Skills Graph Cleanup (2026-05-17)

### Summary
Cleaned up the CHISG skills graph on Vultr Weaviate (port 8088).
Result: 594 skills, 1716 links, 9 legitimate roots.

### Deletions
- **"Time Management Basics"** — not a real skill (children: Lifelong Learning Strategies,
  Project Management — both had 4+ other parents, safe to orphan)
- **"Basic Arithmetic"** — duplicate of "Basic Arithmetic Operations" (child "Operations"
  had 3 other parents, safe to orphan)
- All associated SkillLinks deleted for both skills

### New Links Added (parent → child)
| Link | Reason |
|------|--------|
| Fine Motor Skills → Fine Motor Control | Collection → specific skill |
| Reading Comprehension → Reading Comprehension Strategies | Sub-skill of comprehension |
| Number Recognition → Number Recognition & Writing | NR is precursor to NR&W |
| Abstract Thinking → Philosophical Reasoning | Fixed wrong root |
| Confidence → Influencing | Fixed wrong root |
| Conflict Engagement → Resistance Management | Fixed wrong root |
| Analysis → Stakeholder Analysis | Fixed wrong root |

### Script
`/tmp/chisg_cleanup.py` executed on Vultr via sshpass SSH.
Weaviate port on Vultr: 8088 (not 8081 which is skills-api).
Script used `http://localhost:8088` with GraphQL + REST DELETE endpoints.
Note: Basic Arithmetic UUID mismatch — first delete attempt found wrong UUID via
GraphQL; fixed by direct UUID deletion using the ID from the public API response.

### Verification
Live API at https://espthinking.co.uk/skillstree/api/chisg/skillsandlinks confirmed:
- 594 skills, 1716 links, 9 roots
- All 9 roots verified legitimate
- All 7 new links present
- Both deleted skills absent
""",
    },
    {
        "title": "CHISG Semantic Twin Audit Decisions (2026-05-17)",
        "project": "assist",
        "section_path": "CHISG > Skills Graph > Audit",
        "file_path": "chisg/audit/2026-05-17-twin-decisions.md",
        "content": """## CHISG Semantic Twin Pair Audit Decisions (2026-05-17)

### Already Correctly Linked — No Action Needed
- Market Analysis → Market Gap Analysis ✓
- Goal Setting → Aspirational Goal Setting ✓
- Sentence Structure → Sentence Structure & Variety ✓
- Problem Solving → Collaborative Problem Solving ✓
- Strategic Foresight Training → Strategic Foresight ✓

### By Design — Keep Both (ETP Modulation Leaves)
- Risk Tolerance Modulation vs Risk Tolerance (intentional ETP Self-Awareness leaf)
- Resource Allocation Modulation vs Resource Allocation (intentional ETP leaf)

### Legitimate Separate — Different Concepts/Domains
- Pattern Recognition (FOUNDATIONAL/cognitive) vs Pattern Recognition & Extension (NUMBERS/mathematical)
- Decision Making vs Holistic Decision Making
- Ethical Decision Making vs Decision Making
- Goal Setting (STRATEGIC) vs Goal Setting (Personal) (SELF)

### Undecided / Deferred
- Modelling [PROBLEM SOLVING] vs Modeling [SELF] — both parent Simulation, spelling variant
- Place Value [NUMBERS] vs Place Value Understanding [MATH] — different domains, need decision
- Scientific Method vs Scientific Method Application — already partially chained

### Actions Taken (links added)
See companion doc: CHISG Skills Graph Cleanup 2026-05-17
""",
    },
    {
        "title": "LAO Mobile — Version Bump to 2.1.2 and EAS Config Fix (2026-05-17)",
        "project": "assist",
        "section_path": "LAO Mobile > Publishing > App Store",
        "file_path": "lao-mobile/publishing/2026-05-17-version-bump.md",
        "content": """## LAO Mobile Version Bump + EAS Config Fix (2026-05-17)

### Problem
- app.json had version: "2.1.1"
- eas.json had appVersionSource: "remote" — meaning EAS ignored app.json version
  and pulled version from EAS Cloud (also 2.1.1)
- App Store Connect had a 2.1.2 version page created earlier
- Result: submitted binary was 2.1.1, didn't match ASC version page 2.1.2
- The binary appeared in TestFlight but couldn't be selected on the 2.1.2 version page

### Fix Applied
File: /Users/michaelstewart/Coding/LAOMobile/frontend/app.json
  - version: "2.1.1" → "2.1.2"

File: /Users/michaelstewart/Coding/LAOMobile/frontend/eas.json
  - appVersionSource: "remote" → "local"
  - autoIncrement: true still in production profile (handles build number)

### EAS Build Run
Command: npx eas-cli build --platform ios --profile production
Result: Build queued / completed (Exit Code 0)
Build will have version 2.1.2 with auto-incremented build number beyond 20260412.3

### Next Steps After Build Completes
1. Run: npx eas-cli submit --platform ios --latest
2. In App Store Connect, go to the 2.1.2 version page
3. Add build to the version (+ Build section)
4. Scroll to In-App Purchases section, add both IAP products
5. Submit for review with sandbox reviewer notes

### App Config Reference
- Bundle ID: com.laomobile.frontend
- ASC App ID: 6758043113
- RevenueCat production key: appl_uEENpJPVdDURKGHVTWAWtIgntOc
- IAP products: lao_gcse_science_trilogy, lao_gcse_science_separate
""",
    },
    {
        "title": "LAO Mobile — IAP App Store Connect Workflow Diagnosis (2026-05-17)",
        "project": "assist",
        "section_path": "LAO Mobile > Publishing > IAP",
        "file_path": "lao-mobile/publishing/2026-05-17-iap-workflow-diagnosis.md",
        "content": """## LAO Mobile IAP App Store Connect Workflow — Diagnosis (2026-05-17)

### Recurring Loop Identified
The IAP products have been stuck in a rejection loop since early 2026.
History from transcript:
- Earlier rejection (Guideline 3.1.1): IAP products submitted without a binary attached
- Rejection (May 7-8 2026, Guideline 2.1b): Build 2.1.1 (20260119.25) reviewed on
  iPad Air 11-inch M3, iPadOS 26.4.2 — purchase fails with error in sandbox

### Root Cause of Current Blockage
The user was trying to submit IAPs from the STANDALONE IAP PRODUCT PAGE.
On that page:
- Editing localization text → status changes to "Prepared for Submission" (correct)
- Editing Review Notes → Submit button turns blue
- But "binary attachment" field is empty even though builds exist in TestFlight

WHY: TestFlight builds do NOT appear in the binary selector on the standalone IAP page.
Builds only appear on the APP VERSION PAGE (App Store tab → version → Build section).

### Correct Workflow
1. Go to App Store Connect → App Store tab → version 2.1.2 (not In-App Purchases tab)
2. On the version page: Add Build (+ button in Build section)
   → TestFlight/EAS builds WILL appear here
3. Scroll down to "In-App Purchases and Subscriptions" section
4. Add both IAP products there
5. Submit for review from the VERSION page (not the IAP product page)
This bundles the binary + IAPs in a single review request.

### Paid Apps Agreement Check
If Apple mentioned Paid Apps Agreement in any rejection:
- App Store Connect → Business → Agreements, Tax, and Banking
- "Paid Applications" row must show "Active" (not Pending)
- If not Active: Account Holder must sign it before sandbox purchases work

### Before Submitting Checklist
- [ ] Privacy policy URL set in App Information
- [ ] IAP screenshot (paywall) attached to each product
- [ ] Review notes include sandbox Apple ID + password + how to trigger purchase
- [ ] Support URL set
- [ ] Paid Applications agreement Active

### IAP Product IDs
- lao_gcse_science_trilogy (non-consumable, £2.99)
- lao_gcse_science_separate (non-consumable, £2.99)

### RevenueCat
SDK: react-native-purchases v9.15.2
Production key: appl_uEENpJPVdDURKGHVTWAWtIgntOc
Sandbox detection: automatic (RevenueCat handles via StoreKit environment)
PurchaseContext.js has NativeModules.RNPurchases guard for Expo Go safety.
""",
    },
]


def make_id(title: str) -> str:
    return hashlib.md5(title.encode()).hexdigest()


def log_to_weaviate(client, change: dict):
    collection = client.collections.get(DOC_CLASS)
    obj_id = make_id(change["title"])
    props = {
        "title": change["title"],
        "content": change["content"],
        "file_path": change.get("file_path", ""),
        "section_path": change.get("section_path", ""),
        "project": change.get("project", "assist"),
        "source_type": "session_log",
    }
    try:
        collection.data.insert(properties=props, uuid=obj_id)
        print(f"  ✓ Inserted: {change['title'][:60]}")
    except Exception as e:
        if "already exists" in str(e).lower() or "422" in str(e):
            collection.data.replace(uuid=obj_id, properties=props)
            print(f"  ~ Updated:  {change['title'][:60]}")
        else:
            print(f"  ✗ Failed:   {change['title'][:60]}: {e}")


def main():
    print("Connecting to Weaviate at localhost:8088 ...")
    client = weaviate.connect_to_local(
        host=WEAVIATE_HOST,
        port=WEAVIATE_PORT,
        grpc_port=WEAVIATE_GRPC_PORT,
        additional_config=AdditionalConfig(timeout=Timeout(init=10, query=30, insert=30)),
    )
    try:
        print(f"Connected: {client.is_ready()}")
        print(f"\nLogging {len(CHANGES)} entries to '{DOC_CLASS}'...\n")
        for change in CHANGES:
            log_to_weaviate(client, change)
        print("\nDone.")
    finally:
        client.close()


if __name__ == "__main__":
    main()
