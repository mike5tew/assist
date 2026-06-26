#!/usr/bin/env python3
"""
CHISG Full Pipeline: PDF -> Text -> Chunks -> LLM (OpenAI Structured) -> JSON
=============================================================================
This script ties together the PDF extraction (via pdfplumber) and the 
structured LLM extraction (via extract_llm_pipeline_v2.py).
"""

import os
import sys
import json
import argparse
from pathlib import Path
from tqdm import tqdm

try:
    import pdfplumber
except ImportError:
    print("pdfplumber is not installed. Please run: pip install pdfplumber", file=sys.stderr)
    sys.exit(1)

# Import the new extraction logic
try:
    from extract_llm_pipeline_v2 import extract_links
except ImportError:
    print("Could not import extract_llm_pipeline_v2. Make sure it is in the same directory.", file=sys.stderr)
    sys.exit(1)

# Import chunking logic from the old script
try:
    import re
    from extract_chisg_papers import extract_text_from_pdf, chunk_text
except ImportError:
    print("Could not import from extract_chisg_papers.", file=sys.stderr)
    sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description="End-to-end PDF to Semantic Links pipeline.")
    parser.add_argument("--pdf", required=True, help="Path to the source PDF file.")
    parser.add_argument("--doc-id", required=True, help="Document ID (e.g. le_chen_2022).")
    parser.add_argument("--title", required=True, help="Paper title for metadata.")
    parser.add_argument("--domain", default="microbiology", help="Domain of the text.")
    parser.add_argument("--output", help="Output JSON file. Defaults to <doc_id>_extracted.json")
    args = parser.parse_args()

    if not os.getenv("OPENAI_API_KEY"):
        print("Error: OPENAI_API_KEY not found in environment.", file=sys.stderr)
        sys.exit(1)

    pdf_path = Path(args.pdf)
    if not pdf_path.exists():
        print(f"Error: Target PDF {pdf_path} does not exist.", file=sys.stderr)
        sys.exit(1)

    print(f"Extracting text from: {pdf_path.name}...")
    pages = extract_text_from_pdf(str(pdf_path))
    print(f"Extracted {len(pages)} pages.")

    print("Chunking text...")
    chunks = chunk_text(pages, args.doc_id)
    print(f"Generated {len(chunks)} chunks.")

    all_links = []
    
    print("Sending chunks to OpenAI via Structured Outputs...")
    # Process each chunk sequentially (could be parallelized later)
    for i, chunk in enumerate(tqdm(chunks, desc="Processing chunks")):
        chunk_text_data = chunk['text']
        
        # Add a bit of page context to the chunk for the LLM
        contextual_text = f"PAGE {chunk['page_num']}:\n{chunk_text_data}"
        
        result = extract_links(
            text=contextual_text,
            document_id=args.doc_id,
            paper_title=args.title,
            domain=args.domain
        )
        
        if result and "links" in result:
            # Inject page metadata into the extracted links just to be thorough
            for link in result["links"]:
                link["source_metadata"]["page_reference"] = f"Page {chunk['page_num']}"
            all_links.extend(result["links"])

    output_path = args.output if args.output else f"{args.doc_id}_extracted.json"
    
    with open(output_path, "w", encoding="utf-8") as f:
        json.dump(all_links, f, indent=2)

    print(f"\nExtraction complete! Extracted {len(all_links)} semantic links.")
    print(f"Saved to: {output_path}")

if __name__ == "__main__":
    main()
