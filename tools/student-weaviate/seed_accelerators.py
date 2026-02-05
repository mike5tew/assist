#!/usr/bin/env python3
"""Seed Primary School accelerator skills into Weaviate CHISGElement class.

Usage:
  ./seed_accelerators.py [--weaviate <url>] [--file <path>] [--dry-run]

Defaults:
  WEAVIATE_URL = http://localhost:8081
  FILE = ../../docs/training_data/PRIMARY_SCHOOL_ACCELERATORS.json

This script is intentionally small: it reads the JSON file produced earlier and creates CHISGElement objects.
It maps the suggested ETP fields into existing CHISG properties where possible (etp_spectrum_id, etp_vector, etp_energy_cost).
"""

import argparse
import json
import os
import sys
from urllib import request, error

DEFAULT_FILE = os.path.join(os.path.dirname(__file__), "..", "..", "docs", "training_data", "PRIMARY_SCHOOL_ACCELERATORS.json")
DEFAULT_WEAVIATE = os.environ.get("WEAVIATE_URL", "http://localhost:8081")


def load_json(path):
    with open(path, "r") as f:
        return json.load(f)


import re


def normalize_chisg_id(raw_id, name_fallback=None):
    """Normalize an accelerator id/name into a safe CHISG chisg_id.

    - Removes leading 'accel_' if present
    - Replaces non-alphanumeric characters with underscores
    - Uppercases and prefixes with 'ACCEL_'
    """
    s = ''
    if raw_id:
        s = str(raw_id).strip()
    elif name_fallback:
        s = str(name_fallback).strip()

    # remove leading accel_ if present (case-insensitive)
    s = re.sub(r'(?i)^accel_', '', s)
    # replace whitespace and non-alphanum with underscore
    s = re.sub(r'[^A-Za-z0-9]+', '_', s)
    s = s.strip('_')
    if not s:
        s = 'UNKNOWN'
    return f"ACCEL_{s.upper()}"


def build_properties(accel, normalized_chisg_id=None):
    # Map accelerator entry into CHISGElement properties
    props = {
        "name": accel.get("name"),
        "description": accel.get("short_description", ""),
        "vector_content": f"{accel.get('name')}: {accel.get('short_description','')}",
        "domain": accel.get("domain"),
        "layer": "Accelerator",
        # deterministic chisg_id for traceability
        "chisg_id": normalized_chisg_id or f"ACCEL_{accel.get('id','').upper()}",
        "suggested_years": accel.get("suggested_years", ""),
        # ETP suggestions mapped to existing ETP properties where possible
        "etp_spectrum_id": accel.get("suggested_etp", {}).get("spectrum", ""),
        "etp_vector": accel.get("suggested_etp", {}).get("vector", ""),
        "etp_energy_cost": accel.get("suggested_etp", {}).get("energy_cost", ""),
        "source": "PRIMARY_SCHOOL_ACCELERATORS",
    }
    return props


def exists_object(weaviate_url, chisg_id):
    """Check whether a CHISGElement with the given chisg_id already exists.

    Returns (exists: bool, existing_record: dict|None)
    """
    query_body = json.dumps({
        "query": f"{{ Get {{ CHISGElement(where:{{path:[\"chisg_id\"], operator:Equal, valueString:\"{chisg_id}\"}}) {{ _additional {{ id }} name chisg_id }} }} }}"
    }).encode('utf-8')
    url = weaviate_url.rstrip('/') + '/v1/graphql'
    req = request.Request(url, data=query_body, headers={"Content-Type": "application/json"}, method='POST')
    try:
        with request.urlopen(req) as resp:
            j = json.loads(resp.read().decode('utf-8'))
            items = j.get('data', {}).get('Get', {}).get('CHISGElement', [])
            if items:
                return True, items[0]
            return False, None
    except Exception as e:
        # don't hard-fail on existence checks; warn and treat as not exists
        print(f"Warning: existence check failed for {chisg_id}: {e}")
        return False, None


def post_object(weaviate_url, obj):
    url = weaviate_url.rstrip("/") + "/v1/objects"
    data = json.dumps(obj).encode("utf-8")
    req = request.Request(url, data=data, headers={"Content-Type": "application/json"}, method="POST")
    try:
        with request.urlopen(req) as resp:
            body = resp.read().decode("utf-8")
            return True, json.loads(body)
    except error.HTTPError as e:
        return False, e.read().decode("utf-8")
    except Exception as e:
        return False, str(e)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--weaviate", default=DEFAULT_WEAVIATE, help="Weaviate base URL (default from WEAVIATE_URL or http://localhost:8081)")
    parser.add_argument("--file", default=DEFAULT_FILE, help="Path to accelerator JSON")
    parser.add_argument("--dry-run", action="store_true", help="Do not POST, just print what would be sent")
    args = parser.parse_args()

    if not os.path.exists(args.file):
        print(f"File not found: {args.file}")
        sys.exit(1)

    data = load_json(args.file)
    accels = data.get("accelerators", [])

    print(f"Found {len(accels)} accelerators in {args.file}")

    failures = []
    successes = []

    for a in accels:
        # Normalise chisg_id (prevents ACCEL_ACCEL_ duplication) and build properties
        normalized = normalize_chisg_id(a.get('id'), a.get('name'))
        props = build_properties(a, normalized_chisg_id=normalized)
        obj = {
            "class": "CHISGElement",
            "properties": props
        }

        # Check if object already exists by chisg_id to make seeding idempotent
        exists, existing = exists_object(args.weaviate, normalized)
        if exists:
            existing_id = existing.get('_additional', {}).get('id') if existing else None
            print(f"⏭️ SKIP (exists): {a.get('name')} -> chisg_id={normalized} (weaviate id={existing_id})")
            successes.append((a.get('id'), 'skipped', existing_id))
            continue

        if args.dry_run:
            print("--- DRY RUN ---")
            print(json.dumps(obj, indent=2))
            continue

        ok, res = post_object(args.weaviate, obj)
        if ok:
            created_id = res.get('id')
            print(f"✅ Created: {a.get('name')} -> chisg_id={normalized} id={created_id}")
            successes.append((a.get('id'), 'created', created_id))
        else:
            print(f"❌ Failed: {a.get('name')} -> {res}")
            failures.append((a.get('id'), res))

    print("\nSUMMARY")
    print(f"  successes: {len(successes)}")
    print(f"  failures:  {len(failures)}")

    if failures:
        print("Failures details:")
        for f in failures:
            print(f)


if __name__ == '__main__':
    main()
