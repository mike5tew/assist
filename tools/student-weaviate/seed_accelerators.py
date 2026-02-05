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


def build_properties(accel):
    # Map accelerator entry into CHISGElement properties
    props = {
        "name": accel.get("name"),
        "description": accel.get("short_description", ""),
        "vector_content": f"{accel.get('name')}: {accel.get('short_description','')}",
        "domain": accel.get("domain"),
        "layer": "Accelerator",
        # deterministic chisg_id for traceability
        "chisg_id": f"ACCEL_{accel.get('id','').upper()}",
        "suggested_years": accel.get("suggested_years", ""),
        # ETP suggestions mapped to existing ETP properties where possible
        "etp_spectrum_id": accel.get("suggested_etp", {}).get("spectrum", ""),
        "etp_vector": accel.get("suggested_etp", {}).get("vector", ""),
        "etp_energy_cost": accel.get("suggested_etp", {}).get("energy_cost", ""),
        "source": "PRIMARY_SCHOOL_ACCELERATORS",
    }
    return props


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
        props = build_properties(a)
        obj = {
            "class": "CHISGElement",
            "properties": props
        }

        if args.dry_run:
            print("--- DRY RUN ---")
            print(json.dumps(obj, indent=2))
            continue

        ok, res = post_object(args.weaviate, obj)
        if ok:
            print(f"✅ Created: {a.get('name')} -> id {res.get('id')}")
            successes.append((a.get('id'), res.get('id')))
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
