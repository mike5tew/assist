#!/usr/bin/env python3
"""
CHISG Semantic Link Extraction Pipeline
========================================
Reads all definitions from LAO SQLite, extracts semantic links using Claude
via AWS Bedrock, validates against controlled vocabulary, writes to JSON.

Resumable — picks up where it left off if interrupted.

Prerequisites:
    pip install boto3
    AWS credentials configured (.env or ~/.aws/credentials)

Usage:
    python3 scripts/extract_chisg_links.py
    python3 scripts/extract_chisg_links.py --dry-run
    python3 scripts/extract_chisg_links.py --batch-size 30 --model anthropic.claude-3-haiku-20240307-v1:0

Options:
    --batch-size N     Definitions per API call (default: 25)
    --model MODEL      Bedrock model ID (default: anthropic.claude-3-haiku-20240307-v1:0)
    --output PATH      Output file (default: data/chisg/extraction_full.json)
    --dry-run          Show what would be processed without calling API
    --no-resume        Start fresh, ignore existing output
    --subject SUBJ     Filter by subject: physics, chemistry, all (default: all)
    --region REGION    AWS region (default: eu-west-2)
"""

import json
import os
import sys
import time
import sqlite3
import argparse
from pathlib import Path
from datetime import datetime

# ── Configuration ──────────────────────────────────────────────────────────────

SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DB_PATH = PROJECT_ROOT.parent / "LAOMobile" / "backend" / "lao.db"
DEFAULT_OUTPUT = PROJECT_ROOT / "data" / "chisg" / "extraction_full.json"

SUBJECT_MAP = {1: "physics", 3: "chemistry"}
# Physics (SubjectID=1) covers biology, physics, and mixed GCSE science
# Chemistry (SubjectID=3) is chemistry-specific

CONTROLLED_RELATIONS = [
    "causes", "enables", "inhibits", "is composed of", "is a type of",
    "has property", "has value", "is used for", "is found in",
    "is an example of", "determines", "represents", "contradicts"
]

# ── Extraction Rules (embedded in system prompt) ─────────────────────────────

SYSTEM_PROMPT = """You are a semantic link extractor for the CHISG knowledge graph.
You receive science definitions and extract structured semantic links.

EXTRACTION RULES (mandatory):

R1: Named substances in conditions are entities.
    If a condition names a specific substance, organism, or mechanism, extract it as
    a separate entity with its functional role.

R2: Replace sequence markers with causal relations.
    Convert 'first step', 'leads to', 'then' into explicit causal relations.

R3: Extract what the source states — nothing more.
    Do NOT inject domain knowledge. If the source says A causes C with no mechanism,
    extract A→C. Do NOT insert B.

R4: Granularity is set by the source, not the extractor.
    A GCSE definition that omits a mechanism is complete at GCSE level.

R5: Two layers: faithful extraction vs graph-level analysis.
    You do Layer 1 only — faithful mapping. Do NOT detect gaps or add missing steps.

R6: Definitions are input, not knowledge.
    The definition is consumed to produce links. Do not preserve "is defined as".

R7: No 'is defined as' relation.
    Definitions generate structural claims (causes, is composed of, has value, etc.),
    never "is defined as" or "are defined as" links.

R8: Use ONLY these 13 controlled relations:
    causes, enables, inhibits, is composed of, is a type of, has property,
    has value, is used for, is found in, is an example of, determines,
    represents, contradicts

RELATION SYNONYM MAPPINGS:
    If you want to express a relation not in the controlled vocabulary, use this
    mapping table. The left column is what you might naturally write; the right
    column is the controlled relation to use instead, with direction guidance.

    converts A to B        → A enables B           (the process enables the product)
    produces / generates   → A enables B           (the process enables the output)
    uses / utilises        → B is used for A       (swap direction: the thing used → its purpose)
    is derived from        → A is composed of B    (swap: the product is composed of its source)
    contains               → A is composed of B    (what contains is composed of what it holds)
    absorbs / emits        → A has property B      (absorbing/emitting is a property)
    expands / contracts    → A has property B      (physical behaviour is a property)
    is measured in / units → A has value B         (measurement unit is a value relationship)
    is attached to         → A is composed of B    (structural connection = composition)
    releases               → A causes B            (releasing energy = causing energy output)
    is limited to          → A is found in B       (geographic/context limitation = location)
    leads to / results in  → A causes B            (causal synonym)
    prevents / blocks      → A inhibits B          (prevention synonym)
    requires / needs       → A enables B           (swap: B enables A — the requirement enables the process)
    is classified as       → A is a type of B      (classification synonym)
    includes               → B is a type of A      (swap: the member is a type of the group)
    is connected to        → A is composed of B    (structural connection = composition)
    is specialized for     → A is used for B       (specialisation = purpose)
    is described by        → A has property B      (descriptor = property)
    extracts X from Y      → Y is composed of X    (extraction implies composition)

    IMPORTANT: When a mapping says "swap direction", reverse entity_a and entity_b
    so the controlled relation reads naturally. For example:
        ✗ "wind turbine" uses "kinetic energy"
        ✓ "kinetic energy" is used for "wind power generation"

DOMAIN ASSIGNMENT:
    SubjectID 1 often contains biology and physics mixed together.
    Assign the correct domain based on content: physics, biology, or chemistry.

OUTPUT FORMAT:
    For each definition, return a JSON object:
    {
        "source_id": "lao-kw-XXXX",     // KeywordID
        "source_term": "...",
        "source_definition": "...",
        "subject": "physics|biology|chemistry",
        "links": [
            {
                "entity_a": "...",
                "relation": "...",           // MUST be from controlled vocabulary
                "entity_b": "...",
                "source_id": "lao-kw-XXXX"   // provenance on every link
            }
        ]
    }

    If a definition contains NO extractable links (e.g., it's just a label or is
    too vague), return the object with an empty links array.

    Return ONLY a JSON array of these objects — no markdown, no explanation."""

USER_PROMPT_TEMPLATE = """Extract semantic links from these {count} definitions.
Apply rules R1-R8 strictly. Use ONLY the 13 controlled relations.

Definitions:
{definitions}

Return a JSON array of extraction objects."""

# ── Database ─────────────────────────────────────────────────────────────────

def load_definitions(db_path, subject_filter="all"):
    """Load all definitions from LAO SQLite."""
    conn = sqlite3.connect(db_path)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()

    if subject_filter == "all":
        cursor.execute("SELECT KeywordID, Term, Definition_full, SubjectID FROM keywords ORDER BY KeywordID")
    elif subject_filter == "physics":
        cursor.execute("SELECT KeywordID, Term, Definition_full, SubjectID FROM keywords WHERE SubjectID = 1 ORDER BY KeywordID")
    elif subject_filter == "chemistry":
        cursor.execute("SELECT KeywordID, Term, Definition_full, SubjectID FROM keywords WHERE SubjectID = 3 ORDER BY KeywordID")
    else:
        raise ValueError(f"Unknown subject filter: {subject_filter}")

    definitions = []
    for row in cursor:
        definitions.append({
            "keyword_id": row["KeywordID"],
            "term": row["Term"],
            "definition": row["Definition_full"],
            "subject_id": row["SubjectID"]
        })

    conn.close()
    return definitions


def format_batch_for_prompt(batch):
    """Format a batch of definitions as numbered entries for the prompt."""
    lines = []
    for i, d in enumerate(batch, 1):
        lines.append(f"{i}. [ID: lao-kw-{d['keyword_id']:04d}] {d['term']}")
        lines.append(f"   Definition: {d['definition']}")
        lines.append(f"   SubjectID: {d['subject_id']}")
        lines.append("")
    return "\n".join(lines)


# ── API (AWS Bedrock) ────────────────────────────────────────────────────────

_bedrock_client = None

def get_bedrock_client(region="eu-west-2"):
    """Lazy-init Bedrock client, reuse across calls."""
    global _bedrock_client
    if _bedrock_client is None:
        import boto3
        _bedrock_client = boto3.client('bedrock-runtime', region_name=region)
    return _bedrock_client


def call_claude(batch, model="anthropic.claude-3-haiku-20240307-v1:0", max_retries=3, region="eu-west-2"):
    """Send a batch of definitions to Claude via AWS Bedrock."""
    client = get_bedrock_client(region)
    prompt_text = USER_PROMPT_TEMPLATE.format(
        count=len(batch),
        definitions=format_batch_for_prompt(batch)
    )

    # ~500 output tokens per definition, plus safety margin
    max_tokens = min(len(batch) * 600 + 500, 16384)

    request_body = json.dumps({
        "anthropic_version": "bedrock-2023-05-31",
        "max_tokens": max_tokens,
        "system": SYSTEM_PROMPT,
        "messages": [{"role": "user", "content": [{"type": "text", "text": prompt_text}]}]
    })

    for attempt in range(max_retries):
        try:
            response = client.invoke_model(
                modelId=model,
                body=request_body,
                contentType="application/json"
            )

            result = json.loads(response['body'].read())

            # Check for truncation
            stop_reason = result.get('stop_reason', '')
            if stop_reason == 'max_tokens':
                print(f"\n  ⚠ Response truncated (max_tokens) — reducing batch would help")

            text = result['content'][0]['text'].strip()

            # Strip markdown code fences if present
            if text.startswith("```"):
                text = text.split("\n", 1)[1]  # Remove first line
                if text.endswith("```"):
                    text = text.rsplit("\n", 1)[0]  # Remove last line

            extractions = json.loads(text)
            return extractions

        except json.JSONDecodeError as e:
            print(f"  ⚠ JSON parse error (attempt {attempt+1}/{max_retries}): {e}")
            if attempt < max_retries - 1:
                time.sleep(2 ** attempt)
            else:
                print(f"  ✗ Failed to parse response after {max_retries} attempts")
                print(f"    Raw response: {text[:200]}...")
                return None

        except Exception as e:
            err_str = str(e)
            print(f"  ⚠ Bedrock API error (attempt {attempt+1}/{max_retries}): {err_str[:150]}")
            if "ThrottlingException" in err_str:
                wait = 5 * (attempt + 1)  # Back off harder on throttle
            else:
                wait = 2 ** (attempt + 1)
            if attempt < max_retries - 1:
                print(f"    Retrying in {wait}s...")
                time.sleep(wait)
            else:
                print(f"  ✗ Failed after {max_retries} attempts")
                return None


# ── API (Ollama / Local Llama) ───────────────────────────────────────────────

def call_ollama(batch, model="llama3.2", max_retries=3, base_url="http://localhost:11434"):
    """Send a batch of definitions to a local Ollama model."""
    import urllib.request

    prompt_text = USER_PROMPT_TEMPLATE.format(
        count=len(batch),
        definitions=format_batch_for_prompt(batch)
    )

    # Combine system prompt and user prompt for Ollama's chat API
    request_body = json.dumps({
        "model": model,
        "messages": [
            {"role": "system", "content": SYSTEM_PROMPT},
            {"role": "user", "content": prompt_text}
        ],
        "stream": False,
        "options": {
            "num_predict": min(len(batch) * 600 + 500, 16384),
            "temperature": 0.3
        }
    })

    for attempt in range(max_retries):
        try:
            req = urllib.request.Request(
                f"{base_url}/api/chat",
                data=request_body.encode("utf-8"),
                headers={"Content-Type": "application/json"}
            )
            with urllib.request.urlopen(req, timeout=300) as resp:
                result = json.loads(resp.read())

            text = result["message"]["content"].strip()

            # Robust JSON extraction — models often wrap in prose + code fences
            # 1) Strip markdown code fences anywhere in text
            import re
            fence_match = re.search(r'```(?:json)?\s*\n(.*?)```', text, re.DOTALL)
            if fence_match:
                text = fence_match.group(1).strip()
            else:
                # 2) Find first [ ... last ] (the JSON array)
                first_bracket = text.find("[")
                last_bracket = text.rfind("]")
                if first_bracket != -1 and last_bracket > first_bracket:
                    text = text[first_bracket:last_bracket+1]

            extractions = json.loads(text)
            return extractions

        except json.JSONDecodeError as e:
            print(f"  ⚠ JSON parse error (attempt {attempt+1}/{max_retries}): {e}")
            if attempt < max_retries - 1:
                time.sleep(2)
            else:
                # Try to salvage partial JSON
                try:
                    text_trimmed = text[:text.rfind("]")+1]
                    extractions = json.loads(text_trimmed)
                    print(f"  ⚠ Salvaged partial response ({len(extractions)} items)")
                    return extractions
                except:
                    print(f"  ✗ Failed to parse response after {max_retries} attempts")
                    print(f"    Raw response: {text[:200]}...")
                    return None

        except Exception as e:
            print(f"  ⚠ Ollama error (attempt {attempt+1}/{max_retries}): {str(e)[:150]}")
            if attempt < max_retries - 1:
                time.sleep(2)
            else:
                print(f"  ✗ Failed after {max_retries} attempts")
                return None


# ── Validation ───────────────────────────────────────────────────────────────

def validate_extractions(extractions):
    """Validate extractions against controlled vocabulary. Returns (valid, warnings)."""
    warnings = []
    valid_count = 0
    invalid_count = 0

    for ext in extractions:
        for link in ext.get("links", []):
            relation = link.get("relation", "")
            if relation in CONTROLLED_RELATIONS:
                valid_count += 1
            else:
                invalid_count += 1
                warnings.append(
                    f"  ⚠ Invalid relation '{relation}' in {ext.get('source_id', '?')}: "
                    f"{link.get('entity_a', '?')} → {link.get('entity_b', '?')}"
                )

    return valid_count, invalid_count, warnings


# ── Output ───────────────────────────────────────────────────────────────────

def load_existing_output(output_path):
    """Load existing output for resume functionality."""
    if not output_path.exists():
        return {"extraction_meta": {}, "extractions": []}, set()

    with open(output_path, 'r') as f:
        data = json.load(f)

    processed_ids = set()
    for ext in data.get("extractions", []):
        # Extract keyword ID from source_id like "lao-kw-0001"
        sid = ext.get("source_id", "")
        if sid.startswith("lao-kw-"):
            try:
                processed_ids.add(int(sid.replace("lao-kw-", "")))
            except ValueError:
                pass

    return data, processed_ids


def save_output(output_path, data):
    """Write output JSON atomically."""
    output_path.parent.mkdir(parents=True, exist_ok=True)
    tmp_path = output_path.with_suffix('.tmp')
    with open(tmp_path, 'w') as f:
        json.dump(data, f, indent=2)
    tmp_path.rename(output_path)


# ── Main ─────────────────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(description="CHISG semantic link extraction pipeline")
    parser.add_argument("--batch-size", type=int, default=10, help="Definitions per API call")
    parser.add_argument("--model", default="anthropic.claude-3-haiku-20240307-v1:0",
                        help="Bedrock model ID (ignored if --provider ollama)")
    parser.add_argument("--provider", choices=["bedrock", "ollama"], default="bedrock",
                        help="LLM provider: bedrock (AWS) or ollama (local)")
    parser.add_argument("--ollama-model", default="llama3.2",
                        help="Ollama model name (default: llama3.2)")
    parser.add_argument("--ollama-url", default="http://localhost:11434",
                        help="Ollama base URL")
    parser.add_argument("--output", type=Path, default=DEFAULT_OUTPUT, help="Output JSON path")
    parser.add_argument("--dry-run", action="store_true", help="Show plan without calling API")
    parser.add_argument("--no-resume", action="store_true", help="Start fresh, ignore existing output")
    parser.add_argument("--subject", choices=["all", "physics", "chemistry"], default="all")
    parser.add_argument("--db", type=Path, default=DB_PATH, help="Path to LAO SQLite database")
    parser.add_argument("--region", default="eu-west-2", help="AWS region for Bedrock")
    parser.add_argument("--limit", type=int, default=0, help="Max definitions to process (0=all)")
    args = parser.parse_args()

    # Load .env file for AWS credentials
    env_file = PROJECT_ROOT / ".env"
    if env_file.exists():
        for line in env_file.read_text().splitlines():
            line = line.strip()
            if line and not line.startswith("#") and "=" in line:
                key, val = line.split("=", 1)
                os.environ.setdefault(key.strip(), val.strip())

    # Check dependencies
    if not args.dry_run:
        if args.provider == "bedrock":
            try:
                import boto3
            except ImportError:
                print("Error: pip install boto3")
                sys.exit(1)
            # Quick connectivity test
            try:
                get_bedrock_client(args.region)
                print(f"✓ AWS Bedrock connected (region: {args.region}, model: {args.model})")
            except Exception as e:
                print(f"Error connecting to AWS Bedrock: {e}")
                sys.exit(1)
        elif args.provider == "ollama":
            # Check Ollama server is reachable
            try:
                import urllib.request
                req = urllib.request.Request(f"{args.ollama_url}/api/tags")
                with urllib.request.urlopen(req, timeout=5) as resp:
                    data = json.loads(resp.read())
                    models = [m["name"] for m in data.get("models", [])]
                    if not any(args.ollama_model in m for m in models):
                        print(f"Warning: model '{args.ollama_model}' not found in Ollama. Available: {models}")
                        print(f"  Run: ollama pull {args.ollama_model}")
                        sys.exit(1)
                    print(f"✓ Ollama connected ({args.ollama_url}, model: {args.ollama_model})")
            except Exception as e:
                print(f"Error connecting to Ollama at {args.ollama_url}: {e}")
                print("  Is 'ollama serve' running?")
                sys.exit(1)

    # Load definitions
    print(f"Loading definitions from {args.db}...")
    if not args.db.exists():
        print(f"Error: Database not found at {args.db}")
        sys.exit(1)

    definitions = load_definitions(args.db, args.subject)
    if args.limit > 0:
        definitions = definitions[:args.limit]
    print(f"  Found {len(definitions)} definitions (filter: {args.subject}, limit: {args.limit or 'none'})")

    # Resume handling
    processed_ids = set()
    if not args.no_resume:
        existing_data, processed_ids = load_existing_output(args.output)
        if processed_ids:
            print(f"  Resuming: {len(processed_ids)} already processed")

    remaining = [d for d in definitions if d["keyword_id"] not in processed_ids]
    print(f"  Remaining: {len(remaining)} to process")

    if args.dry_run:
        n_batches = (len(remaining) + args.batch_size - 1) // args.batch_size
        model_display = args.ollama_model if args.provider == "ollama" else args.model
        print(f"\nDry run summary:")
        print(f"  Provider: {args.provider}")
        print(f"  Batches: {n_batches}")
        print(f"  Batch size: {args.batch_size}")
        print(f"  Model: {model_display}")
        if args.provider == "bedrock":
            print(f"  Region: {args.region}")
        print(f"  Output: {args.output}")
        if args.provider == "ollama":
            print(f"  Estimated cost: $0.00 (local)")
        else:
            # Estimate cost — ~500 input tokens per def (inc. system prompt amortised), ~200 output tokens per def
            input_tokens = len(remaining) * 500
            output_tokens = len(remaining) * 200
            # Haiku pricing: $0.25/MTok input, $1.25/MTok output
            # Sonnet 3 pricing: $3/MTok input, $15/MTok output
            if "haiku" in args.model:
                cost_estimate = (input_tokens * 0.25 + output_tokens * 1.25) / 1_000_000
                print(f"  Pricing: Haiku ($0.25/MTok in, $1.25/MTok out)")
            else:
                cost_estimate = (input_tokens * 3 + output_tokens * 15) / 1_000_000
                print(f"  Pricing: Sonnet ($3/MTok in, $15/MTok out)")
            print(f"  Estimated tokens: ~{input_tokens:,} in / ~{output_tokens:,} out")
            print(f"  Estimated cost: ~${cost_estimate:.2f}")
        return

    # Prepare output structure    
    if processed_ids and not args.no_resume:
        output_data = existing_data
    else:
        output_data = {
            "extraction_meta": {
                "source": str(args.db),
                "provider": args.provider,
                "model": args.ollama_model if args.provider == "ollama" else args.model,
                "method": f"LLM extraction ({args.provider}: {args.ollama_model if args.provider == 'ollama' else args.model})",
                "date_started": datetime.now().isoformat(),
                "rules_version": "0.2",
                "controlled_relations": CONTROLLED_RELATIONS,
                "extraction_rules": [
                    "R1: Named substances in conditions are entities",
                    "R2: Replace sequence markers with causal relations",
                    "R3: Extract what the source states — nothing more",
                    "R4: Granularity is set by the source, not the extractor",
                    "R5: Two layers: faithful extraction vs graph-level analysis",
                    "R6: Definitions are input, not knowledge",
                    "R7: No 'is defined as' relation",
                    "R8: Controlled vocabulary of 13 relations"
                ]
            },
            "extractions": [],
            "stats": {
                "total_definitions": len(definitions),
                "total_links": 0,
                "relation_counts": {},
                "zero_link_definitions": 0,
                "invalid_relations": 0,
                "errors": 0
            }
        }

    # Process in batches
    n_batches = (len(remaining) + args.batch_size - 1) // args.batch_size
    total_links = sum(len(e.get("links", [])) for e in output_data.get("extractions", []))
    total_warnings = []

    model_display = args.ollama_model if args.provider == "ollama" else args.model
    print(f"\nProcessing {len(remaining)} definitions in {n_batches} batches of {args.batch_size}...")
    print(f"Provider: {args.provider} | Model: {model_display}")
    print(f"Output: {args.output}\n")

    for batch_idx in range(n_batches):
        start = batch_idx * args.batch_size
        end = min(start + args.batch_size, len(remaining))
        batch = remaining[start:end]

        batch_ids = [d["keyword_id"] for d in batch]
        print(f"Batch {batch_idx+1}/{n_batches} (IDs {batch_ids[0]}-{batch_ids[-1]})...", end=" ", flush=True)

        # Call API — route to provider
        if args.provider == "ollama":
            extractions = call_ollama(batch, model=args.ollama_model, base_url=args.ollama_url)
        else:
            extractions = call_claude(batch, model=args.model, region=args.region)

        if extractions is None:
            print("✗ FAILED — skipping batch")
            output_data["stats"]["errors"] += len(batch)
            continue

        # Validate
        valid, invalid, warnings = validate_extractions(extractions)
        total_warnings.extend(warnings)

        batch_links = sum(len(e.get("links", [])) for e in extractions)
        zero_links = sum(1 for e in extractions if not e.get("links"))
        total_links += batch_links

        output_data["extractions"].extend(extractions)
        output_data["stats"]["total_links"] = total_links
        output_data["stats"]["invalid_relations"] += invalid
        output_data["stats"]["zero_link_definitions"] += zero_links

        # Update relation counts
        for ext in extractions:
            for link in ext.get("links", []):
                rel = link.get("relation", "UNKNOWN")
                output_data["stats"]["relation_counts"][rel] = \
                    output_data["stats"]["relation_counts"].get(rel, 0) + 1

        print(f"✓ {batch_links} links ({zero_links} empty, {invalid} invalid)")

        # Print warnings for this batch
        for w in warnings:
            print(w)

        # Save after each batch (resumable)
        output_data["extraction_meta"]["date_updated"] = datetime.now().isoformat()
        output_data["stats"]["definitions_processed"] = len(output_data["extractions"])
        save_output(args.output, output_data)

        # Rate limiting: 1s between batches
        if batch_idx < n_batches - 1:
            time.sleep(1)

    # Final summary
    stats = output_data["stats"]
    print(f"\n{'='*60}")
    print(f"EXTRACTION COMPLETE")
    print(f"{'='*60}")
    print(f"Definitions processed: {stats.get('definitions_processed', 0)}/{len(definitions)}")
    print(f"Total links extracted: {stats['total_links']}")
    print(f"Average links/def:     {stats['total_links']/max(stats.get('definitions_processed',1),1):.1f}")
    print(f"Zero-link definitions: {stats['zero_link_definitions']}")
    print(f"Invalid relations:     {stats['invalid_relations']}")
    print(f"Errors (skipped):      {stats['errors']}")
    print(f"\nRelation distribution:")
    for rel, count in sorted(stats["relation_counts"].items(), key=lambda x: -x[1]):
        marker = "✓" if rel in CONTROLLED_RELATIONS else "✗"
        print(f"  {marker} {rel}: {count}")

    if total_warnings:
        print(f"\n⚠ {len(total_warnings)} invalid relation warnings (see output for details)")

    print(f"\nOutput saved to: {args.output}")


if __name__ == "__main__":
    main()
