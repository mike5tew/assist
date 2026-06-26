#!/usr/bin/env python3
"""
CHISG LLM Extraction Pipeline V2 (Structured Outputs)
=====================================================
Uses OpenAI's Structured Outputs (JSON Schema) to reliably extract 
SemanticLinks from target texts, specifically matching the canonical schema 
and including nested contexts (e.g. parent_node, experimental conditions).
"""

import os
import sys
import json
import argparse
from typing import List, Optional
from pydantic import BaseModel, Field
from openai import OpenAI

class SourceMetadata(BaseModel):
    document_id: str
    paper_title: str
    page_reference: str
    domain: str
    organism: Optional[str] = None
    extraction_method: str = "human_curated"
    parent_node: Optional[str] = None

class SemanticLinkModel(BaseModel):
    thread_id: str = Field(description="Logical grouping of concepts, e.g. THREAD_CARD_STARVATION_RESPONSE")
    source_term: str = Field(description="Normalized source term. If an experimental context, use descriptive top-level node names.")
    target_term: str = Field(description="Normalized target term.")
    relation_type: str = Field(description="Ontology relationship (e.g., causes, represents, regulates, is used for)")
    source_term_generality: float = Field(description="0.0 to 1.0 scale (specific to general)")
    target_term_generality: float = Field(description="0.0 to 1.0 scale (specific to general)")
    semantic_distance: float = Field(description="0.0 representing definitional, 1.0 representing distant connection")
    relationship_strength: float = Field(description="0.0 to 1.0 scale")
    context: str = Field(description="The contextual explanation of the link, capturing experimental bounds.")
    explanation: str = Field(description="Why this specific relation was chosen.")
    hierarchy_rationale: str = Field(description="Rationale for the assigned generality and semantic distance metrics.")
    source_metadata: SourceMetadata

class ExtractionResponse(BaseModel):
    links: List[SemanticLinkModel]

SYSTEM_PROMPT = """You are a highly precise semantic extraction engine for the CHISG Knowledge Graph.
Your task is to take scientific papers (e.g., PhD microbiology papers) and extract atomic assertions as SemanticLinks.

CRITICAL EXTRACTION RULES:
1. ALWAYS map assertions back to their experimental context. In academic papers, claims are NOT universal truth; they exist under experimental conditions.
2. Use nested context (e.g., `parent_node`) to nest findings under the experiment node.
3. If the fact is specific to a very precise case, append an asterisk `*` to the term to flag it as a case-specific fact.
4. Provide numeric values for depth and generality (e.g. `source_term_generality`), keeping highly specific instances near 0.1 and broad concepts closer to 0.8-1.0.

Your output must precisely follow the provided JSON schema.
"""

def extract_links(text: str, document_id: str, paper_title: str, domain: str) -> dict:
    client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))
    
    # We load our 'few-shot' training data as an example
    few_shot_msg = (
        "Here is the style of output we expect. Observe the deep metadata and hierarchy_rationale:\n"
        '{\n'
        '  "links": [\n'
        '    {\n'
        '      "thread_id": "THREAD_CARD_STARVATION_RESPONSE",\n'
        '      "source_term": "experiment_le_chen_et_al_eLife_2022",\n'
        '      "target_term": "measuring the change in CarD and transcription in different starvation conditions",\n'
        '      "relation_type": "is used for",\n'
        '      "source_term_generality": 0.80,\n'
        '      "target_term_generality": 0.55,\n'
        '      "semantic_distance": 0.10,\n'
        '      "relationship_strength": 1.00,\n'
        '      "context": "Experiment by Le Chen et al., eLife 2022 — investigates CarD protein levels...",\n'
        '      "explanation": "Top-level experiment node is used for its stated purpose.",\n'
        '      "hierarchy_rationale": "Source is the experiment as a whole (0.80), target is the specific experimental aim (0.55).",\n'
        '      "source_metadata": {\n'
        '        "document_id": "le_chen_et_al_eLife_2022",\n'
        '        "paper_title": "Le Chen et al., eLife 2022",\n'
        '        "page_reference": "Page 3",\n'
        '        "domain": "microbiology",\n'
        '        "organism": "Mycobacterium smegmatis mc2155",\n'
        '        "extraction_method": "human_curated"\n'
        '      }\n'
        '    }\n'
        '  ]\n'
        '}'
    )

    try:
        response = client.beta.chat.completions.parse(
            model="gpt-4o",
            messages=[
                {"role": "system", "content": SYSTEM_PROMPT},
                {"role": "user", "content": few_shot_msg},
                {"role": "user", "content": f"Extract semantic links from the following text (Doc ID: {document_id}, Title: {paper_title}, Domain: {domain}):\n\n{text}"}
            ],
            response_format=ExtractionResponse
        )
        return response.choices[0].message.parsed.model_dump()
    except Exception as e:
        print(f"Extraction failed: {e}", file=sys.stderr)
        return None

def main():
    parser = argparse.ArgumentParser(description="Extract SemanticLinks from text via LLM.")
    parser.add_argument("--text", required=True, help="Text to extract from")
    parser.add_argument("--doc-id", default="doc_123", help="Source document ID")
    parser.add_argument("--title", default="Unknown Title", help="Paper title")
    parser.add_argument("--domain", default="microbiology", help="Domain of the text")
    args = parser.parse_args()

    if not os.getenv("OPENAI_API_KEY"):
        print("Error: OPENAI_API_KEY not found in environment.", file=sys.stderr)
        sys.exit(1)

    result = extract_links(args.text, args.doc_id, args.title, args.domain)
    if result:
        print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()
