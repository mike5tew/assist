#!/usr/bin/env python3
"""
CHISG Academic Paper Extraction Pipeline
=========================================
Reads PDF academic papers, extracts text in sections, runs CHISG semantic
link extraction via Claude (AWS Bedrock), writes to JSON.

Designed for: Michael McGrath's PhD repository (Brighton).
Project: "Understanding the Global Control of Mechanisms for Starvation
         Survival in Non-Tuberculous Mycobacteria (NTM)"

Pipeline:
  1. PDF → text (pdfplumber)
  2. Text → sections/chunks (~300 tokens)
  3. Chunks → CHISG links (Claude via Bedrock)
  4. Validate against controlled vocabulary
  5. Write resumable JSON output

Prerequisites:
    pip install boto3 pdfplumber
    AWS credentials configured (.env or ~/.aws/credentials)

Usage:
    # Process all PDFs in data/chisg/mcgrath/papers/
    python3 scripts/extract_chisg_papers.py

    # Dry run (cost estimate only)
    python3 scripts/extract_chisg_papers.py --dry-run

    # Use Sonnet for higher quality
    python3 scripts/extract_chisg_papers.py --model anthropic.claude-3-7-sonnet-20250219-v1:0

    # Process a single PDF
    python3 scripts/extract_chisg_papers.py --input path/to/paper.pdf

    # Use local Ollama
    python3 scripts/extract_chisg_papers.py --provider ollama --ollama-model llama3.2
"""

import json
import os
import re
import sys
import time
import argparse
import hashlib
from pathlib import Path
from datetime import datetime

# ── Configuration ──────────────────────────────────────────────────────────────

SCRIPT_DIR = Path(__file__).parent
PROJECT_ROOT = SCRIPT_DIR.parent
DEFAULT_PAPERS_DIR = PROJECT_ROOT / "data" / "chisg" / "mcgrath" / "papers"
DEFAULT_OUTPUT = PROJECT_ROOT / "data" / "chisg" / "mcgrath" / "extraction.json"

# Target chunk size in tokens (1 token ≈ 4 chars)
CHUNK_TARGET_TOKENS = 300
CHUNK_MAX_TOKENS = 500

CONTROLLED_RELATIONS = [
    "causes", "enables", "inhibits", "is composed of", "is a type of",
    "has property", "has value", "is used for", "is found in",
    "is an example of", "determines", "represents", "contradicts",
    # Extended for academic microbiology
    "correlates with", "supports", "refutes", "regulates",
    "is a mechanism of", "is a marker for", "interacts with",
]

# ── System Prompt for Academic Papers ────────────────────────────────────────

SYSTEM_PROMPT = """You are a semantic link extractor for the CHISG knowledge graph.
You receive sections of academic papers (microbiology/infectious disease domain)
and extract structured semantic links.

EXTRACTION RULES (mandatory):

R1: Named substances, organisms, genes, proteins, and mechanisms are entities.
    If the text names a specific molecule (e.g., "TNF-α"), organism (e.g.,
    "Mycobacterium tuberculosis"), gene, protein, drug, or cell type, extract
    it as a separate entity with its functional role.

R2: Replace sequence markers with causal relations.
    Convert 'leads to', 'results in', 'subsequently', 'triggers' into explicit
    causal relations.

R3: Extract what the source states — nothing more.
    Do NOT inject domain knowledge beyond what is written. If the paper says
    A causes C with no mechanism, extract A→C. Do NOT insert intermediate steps.

R4: Granularity is set by the source, not the extractor.
    Preserve the level of detail in the original text. Do not over-simplify
    or over-elaborate.

R5: Two layers: faithful extraction vs graph-level analysis.
    You do Layer 1 only — faithful mapping. Do NOT detect gaps or add
    missing steps. Do not speculate.

R6: Preserve quantitative claims.
    If the text states specific values, concentrations, percentages, or
    statistical findings, capture these using "has value" relations.

R7: Distinguish claims from evidence.
    Statements that report experimental results should include the method
    context in the entity name where relevant (e.g., "growth inhibition (MIC assay)"
    rather than just "growth inhibition").

R8: Use ONLY these controlled relations:
    causes, enables, inhibits, is composed of, is a type of, has property,
    has value, is used for, is found in, is an example of, determines,
    represents, contradicts, correlates with, supports, refutes, regulates,
    is a mechanism of, is a marker for, interacts with

RELATION SYNONYM MAPPINGS:
    If you want to express a relation not in the controlled vocabulary, use this
    mapping table:

    produces / generates    → A enables B
    activates / upregulates → A enables B
    suppresses / downregulates / blocks → A inhibits B
    binds to / attaches to  → A interacts with B
    encodes / expresses     → A determines B
    attenuates / reduces    → A inhibits B
    is associated with      → A correlates with B
    is implicated in        → A is a mechanism of B
    is indicative of        → A is a marker for B
    promotes / facilitates  → A enables B
    prevents / blocks       → A inhibits B
    modulates               → A regulates B
    is a component of       → A is composed of B (swap: B is composed of A)
    is a subtype of / is a strain of → A is a type of B
    is located in / is expressed in  → A is found in B
    contradicts / conflicts with     → A contradicts B
    confirms / validates    → A supports B
    challenges / disproves  → A refutes B

R9: Capture the structural hierarchy of experimental evidence.
    A paper is composed of experiments. An experiment is composed of figures
    and/or methods. A figure is composed of panels. Each panel represents
    specific data. Extract these structural links using "is composed of" and
    "represents" — they are essential for preserving the contextual meaning
    of every biological claim. Without them, a finding cannot be traced back
    to the conditions that produced it.

    Structural node naming convention:
      paper level:     {paper_id}
      experiment:      {paper_id}_experiment_{n}     (n = 1, 2, ...)
      figure:          {paper_id}_figure_{n}          (n = 1, 2, ...)
      panel:           {paper_id}_figure_{n}_{letter} (e.g. figure_1_a)
      method:          {paper_id}_method_{name}       (e.g. method_western_blot)

    Always extract:
      - paper → experiment(s):  paper_id "is composed of" experiment node
      - experiment → figures:   experiment "is composed of" figure node
      - figure → panels:        figure "is composed of" panel node
      - panel → data content:   panel "represents" what the panel shows
      - method → analyte:       method node "determines" what it measures
      - normalisation control:  control "is used for" normalisation node
      - condition label → protocol: condition label "represents" exact procedure

    Expand all abbreviations on first use. If the text defines CK, CL, NL, PL
    or similar shorthand, extract a "represents" link from the abbreviation to
    its full meaning before using the abbreviation as an entity elsewhere.

R10: Populate evidence_context for every biological claim.
    evidence_context records the experimental conditions under which the
    supporting evidence was gathered. This is PROVENANCE — it records where
    the evidence came from, NOT a statement that the claim is only true under
    those conditions.

    Include in evidence_context (all that are stated in the text):
      - organism and strain (e.g. "Mycobacterium smegmatis mc2155")
      - treatment or condition (e.g. "PBS starvation", "10 mM H2O2")
      - timepoint (e.g. "0.5 hr", "1 hr")
      - detection method (e.g. "Western blot, normalised to SigA")
      - figure/panel reference (e.g. "Figure 1B")

    For general background claims (not tied to this paper's experiments),
    set evidence_context to "general — not this paper's experiment" and
    include the cited reference if given.

DOMAIN: Microbiology / Infectious Disease / Mycobacterium

OUTPUT FORMAT:
    For each text section, return a JSON object:
    {
        "chunk_id": "...",          // the chunk_id provided
        "paper_id": "...",          // the paper_id provided
        "links": [
            {
                "entity_a": "...",
                "relation": "...",          // MUST be from controlled vocabulary
                "inverse_relation": "...",  // reverse direction in plain English
                "entity_b": "...",
                "text_quote": "...",        // direct quote or close paraphrase from source
                "evidence_context": "..."   // experimental conditions of supporting evidence
            }
        ]
    }

    inverse_relation examples:
      "is composed of"  → "is a component of"
      "causes"          → "is caused by"
      "determines"      → "is determined by"
      "inhibits"        → "is inhibited by"
      "enables"         → "is enabled by"
      "represents"      → "is represented by"
      "is used for"     → "uses"
      "is a type of"    → "is a category that includes"
      "has property"    → "is a property of"

    If a section contains NO extractable links (e.g., acknowledgements,
    author affiliations), return the object with an empty links array.

    Return ONLY a JSON array of these objects — no markdown, no explanation."""

USER_PROMPT_TEMPLATE = """Extract semantic links from these {count} sections of academic papers.
Apply rules R1-R10 strictly. Use ONLY the controlled relations.

For each section, first extract the structural hierarchy (R9: paper → experiment → figure → panel → data), then extract biological claims (R1-R8), then populate evidence_context for every claim (R10).

Sections:
{sections}

Return a JSON array of extraction objects."""


# ── PDF Processing ───────────────────────────────────────────────────────────

def extract_text_from_pdf(pdf_path):
    """Extract text from a PDF, preserving page structure."""
    import pdfplumber

    pages = []
    with pdfplumber.open(pdf_path) as pdf:
        for i, page in enumerate(pdf.pages):
            text = page.extract_text()
            if text and text.strip():
                pages.append({
                    "page_num": i + 1,
                    "text": text.strip()
                })
    return pages


def extract_metadata_from_text(full_text):
    """Try to extract title and authors from the first page text."""
    lines = full_text.split("\n")
    title = lines[0].strip() if lines else "Unknown"
    # Heuristic: title is usually the first non-empty line that isn't a journal name
    for line in lines[:5]:
        line = line.strip()
        if len(line) > 20 and not any(w in line.lower() for w in ["journal", "vol.", "doi:", "issn"]):
            title = line
            break
    return {"title": title}


def chunk_text(pages, paper_id, target_tokens=CHUNK_TARGET_TOKENS, max_tokens=CHUNK_MAX_TOKENS):
    """Split pages into chunks of approximately target_tokens size.

    Splits on paragraph boundaries (double newline), then sentence boundaries
    if paragraphs are too long.
    """
    chunks = []
    chunk_idx = 0

    for page in pages:
        # Split into paragraphs
        paragraphs = re.split(r'\n\s*\n', page["text"])

        current_chunk = ""
        for para in paragraphs:
            para = para.strip()
            if not para:
                continue

            # Skip very short fragments (headers, page numbers)
            if len(para) < 20:
                continue

            # Skip references section entries (heuristic)
            if re.match(r'^\[\d+\]', para) or re.match(r'^\d+\.\s+[A-Z][a-z]+,?\s+[A-Z]', para):
                continue

            est_tokens = len(para) / 4  # rough token estimate

            if est_tokens > max_tokens:
                # Split long paragraphs on sentences
                if current_chunk:
                    chunks.append(_make_chunk(current_chunk, chunk_idx, paper_id, page["page_num"]))
                    chunk_idx += 1
                    current_chunk = ""

                sentences = re.split(r'(?<=[.!?])\s+', para)
                sent_chunk = ""
                for sent in sentences:
                    if len(sent_chunk + " " + sent) / 4 > target_tokens and sent_chunk:
                        chunks.append(_make_chunk(sent_chunk, chunk_idx, paper_id, page["page_num"]))
                        chunk_idx += 1
                        sent_chunk = sent
                    else:
                        sent_chunk = (sent_chunk + " " + sent).strip()
                if sent_chunk:
                    current_chunk = sent_chunk
            elif len(current_chunk + "\n\n" + para) / 4 > target_tokens and current_chunk:
                chunks.append(_make_chunk(current_chunk, chunk_idx, paper_id, page["page_num"]))
                chunk_idx += 1
                current_chunk = para
            else:
                current_chunk = (current_chunk + "\n\n" + para).strip()

        # Flush remaining text from this page
        if current_chunk and len(current_chunk) / 4 > 30:  # skip tiny remnants
            chunks.append(_make_chunk(current_chunk, chunk_idx, paper_id, page["page_num"]))
            chunk_idx += 1
            current_chunk = ""

    return chunks


def _make_chunk(text, idx, paper_id, page_num):
    """Create a chunk dict."""
    return {
        "chunk_id": f"{paper_id}-c{idx:03d}",
        "paper_id": paper_id,
        "page_num": page_num,
        "text": text.strip(),
        "est_tokens": len(text.strip()) / 4
    }


def make_paper_id(pdf_path):
    """Generate a stable paper_id from filename."""
    stem = Path(pdf_path).stem
    # Clean up filename to make a readable ID
    clean = re.sub(r'[^a-zA-Z0-9_-]', '_', stem)
    clean = re.sub(r'_+', '_', clean).strip('_')
    return clean[:60]


# ── API (AWS Bedrock) ────────────────────────────────────────────────────────

_bedrock_client = None

def get_bedrock_client(region="eu-west-2"):
    global _bedrock_client
    if _bedrock_client is None:
        import boto3
        _bedrock_client = boto3.client('bedrock-runtime', region_name=region)
    return _bedrock_client


def format_batch_for_prompt(batch):
    """Format a batch of chunks as numbered entries for the prompt."""
    lines = []
    for i, chunk in enumerate(batch, 1):
        lines.append(f"--- Section {i} [chunk_id: {chunk['chunk_id']}, paper: {chunk['paper_id']}, page: {chunk['page_num']}] ---")
        lines.append(chunk["text"])
        lines.append("")
    return "\n".join(lines)


def call_claude(batch, model, max_retries=3, region="eu-west-2"):
    """Send a batch of chunks to Claude via AWS Bedrock."""
    client = get_bedrock_client(region)
    prompt_text = USER_PROMPT_TEMPLATE.format(
        count=len(batch),
        sections=format_batch_for_prompt(batch)
    )

    # More output per chunk than GCSE defs — academic text is denser
    max_tokens = min(len(batch) * 800 + 500, 16384)

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

            stop_reason = result.get('stop_reason', '')
            if stop_reason == 'max_tokens':
                print(f"\n  ⚠ Response truncated (max_tokens)")

            # Extract usage for cost tracking
            usage = result.get('usage', {})

            text = result['content'][0]['text'].strip()

            # Strip markdown code fences if present
            if text.startswith("```"):
                text = text.split("\n", 1)[1]
                if text.endswith("```"):
                    text = text.rsplit("\n", 1)[0]

            extractions = json.loads(text)
            return extractions, usage

        except json.JSONDecodeError as e:
            print(f"  ⚠ JSON parse error (attempt {attempt+1}/{max_retries}): {e}")
            if attempt < max_retries - 1:
                time.sleep(2 ** attempt)
            else:
                print(f"  ✗ Failed to parse response after {max_retries} attempts")
                return None, {}

        except Exception as e:
            err_str = str(e)
            print(f"  ⚠ Bedrock API error (attempt {attempt+1}/{max_retries}): {err_str[:150]}")
            if "ThrottlingException" in err_str:
                wait = 5 * (attempt + 1)
            else:
                wait = 2 ** (attempt + 1)
            if attempt < max_retries - 1:
                print(f"    Retrying in {wait}s...")
                time.sleep(wait)
            else:
                print(f"  ✗ Failed after {max_retries} attempts")
                return None, {}


def call_ollama(batch, model="llama3.2", max_retries=3, base_url="http://localhost:11434"):
    """Send a batch of chunks to a local Ollama model."""
    import urllib.request

    prompt_text = USER_PROMPT_TEMPLATE.format(
        count=len(batch),
        sections=format_batch_for_prompt(batch)
    )

    request_body = json.dumps({
        "model": model,
        "messages": [
            {"role": "system", "content": SYSTEM_PROMPT},
            {"role": "user", "content": prompt_text}
        ],
        "stream": False,
        "options": {
            "num_predict": min(len(batch) * 800 + 500, 16384),
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

            fence_match = re.search(r'```(?:json)?\s*\n(.*?)```', text, re.DOTALL)
            if fence_match:
                text = fence_match.group(1).strip()
            else:
                first_bracket = text.find("[")
                last_bracket = text.rfind("]")
                if first_bracket != -1 and last_bracket > first_bracket:
                    text = text[first_bracket:last_bracket+1]

            extractions = json.loads(text)
            return extractions, {}

        except json.JSONDecodeError as e:
            print(f"  ⚠ JSON parse error (attempt {attempt+1}/{max_retries}): {e}")
            if attempt < max_retries - 1:
                time.sleep(2)
            else:
                try:
                    text_trimmed = text[:text.rfind("]")+1]
                    extractions = json.loads(text_trimmed)
                    return extractions, {}
                except:
                    print(f"  ✗ Failed after {max_retries} attempts")
                    return None, {}

        except Exception as e:
            print(f"  ⚠ Ollama error (attempt {attempt+1}/{max_retries}): {str(e)[:150]}")
            if attempt < max_retries - 1:
                time.sleep(2)
            else:
                return None, {}


# ── Validation ───────────────────────────────────────────────────────────────

def validate_extractions(extractions):
    """Validate extractions against controlled vocabulary."""
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
                    f"  ⚠ Invalid relation '{relation}': "
                    f"{link.get('entity_a', '?')} → {link.get('entity_b', '?')}"
                )

    return valid_count, invalid_count, warnings


# ── Output ───────────────────────────────────────────────────────────────────

def load_existing_output(output_path):
    """Load existing output for resume functionality."""
    if not output_path.exists():
        return {"extraction_meta": {}, "papers": {}, "extractions": [], "stats": {}}, set()

    with open(output_path, 'r') as f:
        data = json.load(f)

    processed_ids = {ext.get("chunk_id") for ext in data.get("extractions", []) if ext.get("chunk_id")}
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
    parser = argparse.ArgumentParser(description="CHISG academic paper extraction pipeline")
    parser.add_argument("--input", type=Path, help="Single PDF file to process")
    parser.add_argument("--papers-dir", type=Path, default=DEFAULT_PAPERS_DIR,
                        help="Directory containing PDF files")
    parser.add_argument("--output", type=Path, default=DEFAULT_OUTPUT, help="Output JSON path")
    parser.add_argument("--batch-size", type=int, default=5,
                        help="Chunks per API call (default: 5)")
    parser.add_argument("--model", default="anthropic.claude-haiku-4-5-20251001-v1:0",
                        help="Bedrock model ID")
    parser.add_argument("--provider", choices=["bedrock", "ollama"], default="bedrock")
    parser.add_argument("--ollama-model", default="llama3.2")
    parser.add_argument("--ollama-url", default="http://localhost:11434")
    parser.add_argument("--region", default="eu-west-2", help="AWS region for Bedrock")
    parser.add_argument("--dry-run", action="store_true", help="Show plan and cost estimate only")
    parser.add_argument("--no-resume", action="store_true", help="Start fresh")
    parser.add_argument("--text-only", action="store_true",
                        help="Extract text and chunks only (no LLM calls)")
    args = parser.parse_args()

    # Load .env
    env_file = PROJECT_ROOT / ".env"
    if env_file.exists():
        for line in env_file.read_text().splitlines():
            line = line.strip()
            if line and not line.startswith("#") and "=" in line:
                key, val = line.split("=", 1)
                os.environ.setdefault(key.strip(), val.strip())

    # Find PDFs
    if args.input:
        pdf_files = [args.input]
    else:
        if not args.papers_dir.exists():
            print(f"Error: Papers directory not found: {args.papers_dir}")
            print(f"  Create it and add PDF files: mkdir -p {args.papers_dir}")
            sys.exit(1)
        pdf_files = sorted(args.papers_dir.glob("*.pdf"))

    if not pdf_files:
        print(f"No PDF files found in {args.papers_dir}")
        print(f"  Add PDF papers to: {args.papers_dir}/")
        sys.exit(0)

    print(f"Found {len(pdf_files)} PDF files")

    # ── Stage 1: PDF → Text → Chunks ────────────────────────────────────────

    print("\n═══ Stage 1: PDF Extraction & Chunking ═══\n")

    all_chunks = []
    paper_metadata = {}

    for pdf_path in pdf_files:
        paper_id = make_paper_id(pdf_path)
        print(f"  📄 {pdf_path.name} → {paper_id}")

        try:
            pages = extract_text_from_pdf(pdf_path)
            if not pages:
                print(f"    ⚠ No text extracted (scanned PDF? Try OCR)")
                continue

            full_text = "\n\n".join(p["text"] for p in pages)
            meta = extract_metadata_from_text(full_text)
            meta["filename"] = pdf_path.name
            meta["pages"] = len(pages)
            meta["total_chars"] = len(full_text)
            meta["est_tokens"] = len(full_text) / 4
            paper_metadata[paper_id] = meta

            chunks = chunk_text(pages, paper_id)
            all_chunks.extend(chunks)

            print(f"    {len(pages)} pages → {len(chunks)} chunks "
                  f"(~{meta['est_tokens']:.0f} tokens)")

        except Exception as e:
            print(f"    ✗ Error processing: {e}")

    if not all_chunks:
        print("\nNo chunks extracted from any PDF. Nothing to process.")
        sys.exit(0)

    total_tokens = sum(c["est_tokens"] for c in all_chunks)
    print(f"\n  Total: {len(pdf_files)} papers → {len(all_chunks)} chunks (~{total_tokens:.0f} tokens)")

    if args.text_only:
        # Save just the chunks for inspection
        text_output = args.output.with_name("chunks_preview.json")
        with open(text_output, 'w') as f:
            json.dump({
                "papers": paper_metadata,
                "chunks": all_chunks
            }, f, indent=2)
        print(f"\n  Chunks saved to: {text_output}")
        return

    # ── Cost Estimate ────────────────────────────────────────────────────────

    n_batches = (len(all_chunks) + args.batch_size - 1) // args.batch_size
    # Input: system prompt (~3000 tokens) + batch chunks
    est_input_per_batch = 3000 + args.batch_size * CHUNK_TARGET_TOKENS
    est_output_per_batch = args.batch_size * 250  # ~250 output tokens per chunk
    total_input = n_batches * est_input_per_batch
    total_output = n_batches * est_output_per_batch

    if args.provider == "ollama":
        cost_str = "$0.00 (local)"
    else:
        model_name = args.model.lower()
        if "haiku" in model_name:
            # Claude 3.5 Haiku: $0.80/MTok in, $4/MTok out
            cost = (total_input * 0.80 + total_output * 4.0) / 1_000_000
            pricing_str = "Haiku ($0.80/MTok in, $4.00/MTok out)"
        elif "sonnet" in model_name:
            # Claude 3.5 Sonnet: $3/MTok in, $15/MTok out
            cost = (total_input * 3.0 + total_output * 15.0) / 1_000_000
            pricing_str = "Sonnet ($3.00/MTok in, $15.00/MTok out)"
        else:
            cost = (total_input * 3.0 + total_output * 15.0) / 1_000_000
            pricing_str = "Unknown model (using Sonnet pricing)"
        cost_str = f"~${cost:.2f} ({pricing_str})"

    print(f"\n═══ Extraction Plan ═══\n")
    print(f"  Papers:      {len(pdf_files)}")
    print(f"  Chunks:      {len(all_chunks)}")
    print(f"  Batches:     {n_batches} (batch size: {args.batch_size})")
    print(f"  Provider:    {args.provider}")
    print(f"  Model:       {args.model if args.provider == 'bedrock' else args.ollama_model}")
    print(f"  Est tokens:  ~{total_input:,} in / ~{total_output:,} out")
    print(f"  Est cost:    {cost_str}")
    print(f"  Output:      {args.output}")

    if args.dry_run:
        print("\n  [DRY RUN — no API calls made]")
        return

    # ── Stage 2: CHISG Link Extraction ───────────────────────────────────────

    if args.provider == "bedrock":
        try:
            import boto3
        except ImportError:
            print("Error: pip install boto3")
            sys.exit(1)
        try:
            get_bedrock_client(args.region)
            print(f"\n✓ AWS Bedrock connected (region: {args.region})")
        except Exception as e:
            print(f"Error connecting to AWS Bedrock: {e}")
            sys.exit(1)

    # Resume handling
    processed_chunk_ids = set()
    if not args.no_resume:
        output_data, processed_chunk_ids = load_existing_output(args.output)
        if processed_chunk_ids:
            print(f"  Resuming: {len(processed_chunk_ids)} chunks already processed")
    else:
        output_data = {"extraction_meta": {}, "papers": {}, "extractions": [], "stats": {}}

    remaining = [c for c in all_chunks if c["chunk_id"] not in processed_chunk_ids]
    print(f"  Remaining: {len(remaining)} chunks to process")

    if not remaining:
        print("  Nothing new to process.")
        return

    # Update metadata
    output_data["extraction_meta"] = {
        "project": "Understanding the Global Control of Mechanisms for Starvation Survival in Non-Tuberculous Mycobacteria (NTM)",
        "researcher": "Michael McGrath",
        "institution": "University of Brighton",
        "provider": args.provider,
        "model": args.model if args.provider == "bedrock" else args.ollama_model,
        "controlled_relations": CONTROLLED_RELATIONS,
        "date_started": output_data.get("extraction_meta", {}).get(
            "date_started", datetime.now().isoformat()),
        "date_updated": datetime.now().isoformat(),
    }
    output_data["papers"] = paper_metadata

    print(f"\n═══ Stage 2: CHISG Link Extraction ═══\n")

    total_links = 0
    total_valid = 0
    total_invalid = 0
    total_input_tokens = 0
    total_output_tokens = 0
    batch_count = 0

    for i in range(0, len(remaining), args.batch_size):
        batch = remaining[i:i + args.batch_size]
        batch_count += 1
        paper_ids_in_batch = set(c["paper_id"] for c in batch)

        print(f"  Batch {batch_count}/{n_batches} "
              f"({len(batch)} chunks from {', '.join(paper_ids_in_batch)})...",
              end="", flush=True)

        if args.provider == "bedrock":
            extractions, usage = call_claude(batch, args.model, region=args.region)
        else:
            extractions, usage = call_ollama(batch, args.ollama_model, base_url=args.ollama_url)

        if extractions is None:
            print(" ✗ (skipped)")
            continue

        # Track token usage
        if usage:
            total_input_tokens += usage.get("input_tokens", 0)
            total_output_tokens += usage.get("output_tokens", 0)

        # Validate
        valid, invalid, warnings = validate_extractions(extractions)
        total_valid += valid
        total_invalid += invalid

        batch_links = sum(len(e.get("links", [])) for e in extractions)
        total_links += batch_links

        # Add to output
        output_data["extractions"].extend(extractions)

        print(f" ✓ {batch_links} links ({valid} valid, {invalid} invalid)")

        if warnings and invalid > 0:
            for w in warnings[:3]:
                print(w)
            if len(warnings) > 3:
                print(f"    ... and {len(warnings)-3} more")

        # Save after each batch (resumable)
        output_data["stats"] = {
            "total_papers": len(pdf_files),
            "total_chunks": len(all_chunks),
            "chunks_processed": len(all_chunks) - len(remaining) + i + len(batch),
            "total_links": total_links,
            "valid_relations": total_valid,
            "invalid_relations": total_invalid,
            "input_tokens": total_input_tokens,
            "output_tokens": total_output_tokens,
        }
        save_output(args.output, output_data)

        # Rate limiting
        time.sleep(1)

    # ── Summary ──────────────────────────────────────────────────────────────

    print(f"\n═══ Extraction Complete ═══\n")
    print(f"  Papers processed: {len(pdf_files)}")
    print(f"  Chunks processed: {len(all_chunks)}")
    print(f"  Total links:      {total_links}")
    print(f"  Valid relations:   {total_valid} ({total_valid/(total_valid+total_invalid)*100:.1f}%)" if (total_valid+total_invalid) > 0 else "")
    print(f"  Invalid relations: {total_invalid}")
    if total_input_tokens:
        print(f"  API tokens used:   {total_input_tokens:,} in / {total_output_tokens:,} out")
        if "haiku" in args.model.lower():
            actual_cost = (total_input_tokens * 0.80 + total_output_tokens * 4.0) / 1_000_000
        else:
            actual_cost = (total_input_tokens * 3.0 + total_output_tokens * 15.0) / 1_000_000
        print(f"  Actual cost:       ~${actual_cost:.2f}")
    print(f"  Output:            {args.output}")


if __name__ == "__main__":
    main()
