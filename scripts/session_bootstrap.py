#!/usr/bin/env python3
"""
Session Bootstrap — ESP Thinking Project
Run this at the start of any coding session to dump current project context
from Weaviate into /memories/session/esp-context.md for the LLM to read.

Usage:
  cd /Users/michaelstewart/Coding/humanOS/scripts
  source .venv/bin/activate
  python3 /Users/michaelstewart/Coding/assist/scripts/session_bootstrap.py
"""

import sys
import os
import datetime

sys.path.insert(0, os.path.dirname(__file__))

# Reuse the search function from the humanOS scripts directory
SCRIPTS_DIR = "/Users/michaelstewart/Coding/humanOS/scripts"
sys.path.insert(0, SCRIPTS_DIR)

from search_docs import search_docs

SESSION_MEMORY_PATH = os.path.expanduser("~/.config/github-copilot/memories/session/esp-context.md")
# VS Code Copilot memory path
COPILOT_MEMORY_PATH = "/memories/session/esp-context.md"

QUERIES = [
    ("vultr infrastructure docker deploy nginx ssl", "Infrastructure & Deployment"),
    ("outage incident fix docker network", "Recent Incidents"),
    ("esp pilot frontend deployment", "Frontend & Features"),
    ("project priority deadline status", "Project Status"),
    ("NTM extraction pipeline CHISG context bedrock", "NTM / CHISG Pipeline"),
]

def build_topic_tree(client) -> str:
    """Return a compact topic tree from all section_path values in Weaviate."""
    from collections import defaultdict
    try:
        col = client.collections.get("Documentation")
        r = col.query.fetch_objects(limit=2000, return_properties=["section_path", "project"])
        # Collect unique (project, path) pairs
        tree_lines = ["## Knowledge Base — Topic Tree\n",
                      "_Browse this to identify relevant areas before searching._\n"]
        by_project = defaultdict(set)
        for obj in r.objects:
            sp = obj.properties.get("section_path", "").strip()
            proj = obj.properties.get("project", "?")
            if sp and ">" in sp:
                # Take only the first two levels for the summary tree
                parts = [p.strip() for p in sp.split(">")]
                top = parts[0]
                sub = parts[1] if len(parts) > 1 else None
                by_project[proj].add((top, sub))
        
        for proj in sorted(by_project.keys()):
            tree_lines.append(f"\n### [{proj}]")
            # Group by top-level
            tops = defaultdict(set)
            for top, sub in by_project[proj]:
                if sub:
                    tops[top].add(sub)
                else:
                    tops[top]  # ensure key exists
            for top in sorted(tops.keys()):
                subs = tops[top]
                tree_lines.append(f"  {top}")
                for sub in sorted(subs)[:6]:  # limit breadth
                    tree_lines.append(f"    └─ {sub}")
                if len(subs) > 6:
                    tree_lines.append(f"    └─ ... ({len(subs)-6} more)")
        return "\n".join(tree_lines)
    except Exception as e:
        return f"## Knowledge Base — Topic Tree\n_Could not build tree: {e}_"

def run_bootstrap():
    print("🚀 ESP Thinking Session Bootstrap")
    print("=" * 60)
    print(f"Querying Weaviate at localhost:8088...")

    import weaviate as _weaviate
    try:
        _client = _weaviate.connect_to_local(host="localhost", port=8088, grpc_port=50052)
        tree_section = build_topic_tree(_client)
        _client.close()
    except Exception as e:
        tree_section = f"## Knowledge Base — Topic Tree\n_Weaviate unavailable: {e}_"

    sections = []
    sections.append(f"# ESP Thinking — Session Context\n_Generated: {datetime.datetime.now().strftime('%Y-%m-%d %H:%M')}_\n")
    sections.append(tree_section)
    sections.append("""## Critical Infrastructure Facts
- **Deploy dir**: `/opt/esp-thinking/docker-compose.prod.yml` ← ALWAYS use this
- **DO NOT USE**: `/opt/esp/` — old stack, wrong networks, empty SSL volume
- **main-proxy networks**: esp-thinking_shared-network, esp_shared-network, chisg-classifier_chisg-net
- **SSL cert**: esp-thinking_certbot_certs volume, cert2.pem, expires ~2026-08-14, NO auto-renew
- **SSH**: `ssh vultr` (192.248.151.185), password: E5P_Th!nk!ng?
- **Weaviate**: localhost:8088 | **MongoDB**: localhost:27018 admin:password123 db:esp_organizer
""")

    for query, label in QUERIES:
        print(f"\n🔍 {label}...")
        results = search_docs(query, project="assist", limit=3)

        if not results:
            print(f"   No results.")
            continue

        section_lines = [f"\n## {label}\n"]
        for r in results:
            section_lines.append(f"### {r['title']}")
            section_lines.append(f"_Path: {r['section_path']}_\n")
            content = r.get('content', '')
            if content:
                if len(content) > 600:
                    content = content[:600] + "..."
                section_lines.append(content)
            section_lines.append("")
            print(f"   ✓ {r['title']}")

        sections.append("\n".join(section_lines))

    full_content = "\n".join(sections)

    # Write to the workspace context file — read by copilot-instructions.md at session start
    output_path = "/Users/michaelstewart/Coding/assist/docs/CURRENT_CONTEXT.md"
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    with open(output_path, "w") as f:
        f.write(full_content)

    print(f"\n✅ Context written to: {output_path}")
    print(f"   Characters: {len(full_content)}")
    print("\n💡 Copilot will read this file automatically via copilot-instructions.md")
    print("\n--- CONTEXT PREVIEW ---")
    print(full_content[:1000])
    print("...")

if __name__ == "__main__":
    run_bootstrap()
