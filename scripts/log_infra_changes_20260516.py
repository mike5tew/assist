#!/usr/bin/env python3
"""
Log infrastructure changes from 2026-05-16 to Weaviate Documentation class
and MongoDB project_notes collection.

Usage: python3 log_infra_changes_20260516.py
"""

import weaviate
import hashlib
from datetime import datetime
from pymongo import MongoClient

# Weaviate — assist instance (port 8088)
WEAVIATE_HOST = "localhost"
WEAVIATE_PORT = 8088
WEAVIATE_GRPC_PORT = 50052
DOC_CLASS = "Documentation"

# MongoDB — local assist instance (port 27018 on host, admin:password123)
MONGO_URI = "mongodb://admin:password123@localhost:27018/esp_organizer?authSource=admin"
MONGO_DB = "esp_organizer"
NOTES_COLLECTION = "project_notes"

CHANGES = [
    {
        "title": "Production Outage Fix — espthinking.co.uk (2026-05-16)",
        "project": "assist",
        "section_path": "Infrastructure > Vultr Production > Incidents",
        "file_path": "infra/incidents/2026-05-16-outage.md",
        "content": """## Production Outage — espthinking.co.uk (2026-05-16)

### Summary
Site was returning ERR_CONNECTION_TIMED_OUT. Root cause: main-proxy nginx container 
was crash-looping on Vultr. Total downtime ~7 days (main-proxy was last restarted 
and failed, other services remained up).

### Root Cause
Two separate docker-compose stacks existed on Vultr:
- `/opt/esp/` — old stack with empty certbot volumes and `esp_shared-network`
- `/opt/esp-thinking/` — current stack with real certs and `esp-thinking_shared-network`

main-proxy was running from `/opt/esp/` (wrong stack), so:
1. nginx upstream `assist-frontend:80` resolved to `esp-thinking_shared-network` but 
   main-proxy was only on `esp_shared-network` → DNS resolution failed → crash loop
2. Certbot volume `esp_certbot_certs` was empty → no SSL certs → nginx couldn't start
3. SSL certificate had also expired (last issued 2026-02-14)

### Fix Applied
1. Diagnosed via SSH docker logs: `host not found in upstream "assist-frontend:80"`
2. Stopped and removed main-proxy from old `/opt/esp/` stack
3. Started main-proxy from `/opt/esp-thinking/docker-compose.prod.yml` (correct stack)
   → now uses `esp-thinking_certbot_certs` volume which has real certs
4. Connected main-proxy to `esp_shared-network` (for drb-frontend) and 
   `chisg-classifier_chisg-net` (for chisg-lb) since those services run on different networks
5. Renewed expired SSL certificate via certbot:
   `docker compose -f docker-compose.prod.yml --profile ssl run --rm certbot renew`
   → cert renewed to cert2.pem, valid until ~2026-08-14
6. Reloaded nginx: `docker exec main-proxy nginx -s reload`
7. Verified: `curl -I https://espthinking.co.uk` → HTTP/2 200

### Docker Network Map (Post-Fix)
- `esp-thinking_shared-network`: assist-frontend, assist-api, skills-frontend, skills-api, 
  drb-api, sampletrack-web, sampletrack-api, mongodb, weaviate, skills-db, main-proxy
- `esp_shared-network`: drb-frontend, skills-frontend, skills-api, drb-api, main-proxy
- `chisg-classifier_chisg-net`: chisg-lb, main-proxy

### Prevention
All services should be managed exclusively from `/opt/esp-thinking/docker-compose.prod.yml`. 
The old `/opt/esp/` stack should not be used to start any containers.
""",
    },
    {
        "title": "SSL Certificate Renewal — espthinking.co.uk (2026-05-16)",
        "project": "assist",
        "section_path": "Infrastructure > Vultr Production > SSL",
        "file_path": "infra/ssl/2026-05-16-cert-renewal.md",
        "content": """## SSL Certificate Renewal — espthinking.co.uk (2026-05-16)

### Context
Certificate expired (original issue date: 2026-02-14). Was not auto-renewed because 
main-proxy was not running, so the webroot challenge could not be served.

### Renewal Method
- Certbot webroot challenge via docker-compose profile:
  ```bash
  cd /opt/esp-thinking
  docker compose -f docker-compose.prod.yml --profile ssl run --rm certbot renew
  ```
- Certbot uses volume `esp-thinking_certbot_www` for challenge files
- Certs stored in volume `esp-thinking_certbot_certs`

### Result
- New cert: `archive/espthinking.co.uk/cert2.pem`
- Symlinks updated: `live/espthinking.co.uk/fullchain.pem → ../../archive/espthinking.co.uk/fullchain2.pem`
- nginx reloaded to pick up new cert

### Auto-Renewal Note
Certbot is not running as a daemon — renewal must be triggered manually or via cron.
To add auto-renewal on Vultr:
```bash
crontab -e
# Add: 0 3 * * * cd /opt/esp-thinking && docker compose -f docker-compose.prod.yml --profile ssl run --rm certbot renew && docker exec main-proxy nginx -s reload
```
""",
    },
    {
        "title": "ESP Pilot Study Landing Page Added — /esp-pilot (2026-05-16)",
        "project": "assist",
        "section_path": "Frontend > Pages > ESP Research",
        "file_path": "frontend/src/pages/ESPPilotLanding.tsx",
        "content": """## ESP Pilot Study Landing Page

### Route
`https://espthinking.co.uk/esp-pilot`

### File
`/Users/michaelstewart/Coding/assist/frontend/src/pages/ESPPilotLanding.tsx`

### Purpose
Landing page for recruiting participants to the ESP personal knowledge oracle pilot study.
Explains the research concept, shows a 4-step 'how it works', lists Stage 2 features, 
and provides a contact reveal for expressing interest.

### Key Components
- `ESPPilotNav`: Logo + CHISG + Back to Portfolio nav (matches CHISGLanding pattern)
- Hero: Dark gradient (`#020617 → #1e1b4b`), two CTAs (Join pilot / How it works)
- Problem section: Why personal scientific knowledge is hard to manage
- How it Works: 4-step numbered process (Upload → Link → Query → Discover)
- Stage 2 Features: 6-card grid (Citation Graph, Collaboration, etc.)
- Pilot Study CTA: ContactReveal for `michael.stewart@espthinking.co.uk`
- Cross-link: Card linking back to CHISG at `/chisg`
- SEO: title="ESP Research Pilot — Personal Scientific Knowledge Oracle"

### Cross-Links
- CHISGLanding.tsx: Added a dark call-out card before existing CTA section, 
  linking to `/esp-pilot` with title "CHISG applied to academic research"
- App.tsx: Route added after `/chisg` route

### Status
Built and validated (no TypeScript errors). Not yet deployed to Vultr production 
(needs docker build + push + pull on Vultr).
""",
    },
    {
        "title": "Vultr Docker Stack Architecture — Two-Stack Problem (2026-05-16)",
        "project": "assist",
        "section_path": "Infrastructure > Vultr Production > Architecture",
        "file_path": "infra/vultr/docker-stack-architecture.md",
        "content": """## Vultr Docker Stack Architecture

### Server
IP: 192.248.151.185, SSH alias: `vultr`, key: `~/.ssh/gitkey11-25`

### Active Stack (use this)
Location: `/opt/esp-thinking/`
Compose file: `docker-compose.prod.yml`
Networks: `esp-thinking_shared-network` (primary), `esp_shared-network` (legacy)
Cert volumes: `esp-thinking_certbot_certs`, `esp-thinking_certbot_www`

### Legacy Stack (DO NOT USE)
Location: `/opt/esp/`
Compose file: `docker-compose.yml`
Networks: `esp_shared-network`
Cert volumes: `esp_certbot_certs` (EMPTY — no real certs)

### Services and Their Networks
| Container | Networks |
|-----------|----------|
| main-proxy | esp-thinking_shared-network, esp_shared-network, chisg-classifier_chisg-net |
| assist-frontend | esp-thinking_shared-network |
| assist-api | esp-thinking_shared-network |
| skills-frontend | esp-thinking_shared-network, esp_shared-network |
| skills-api | esp-thinking_shared-network, esp_shared-network |
| drb-frontend | esp_shared-network |
| drb-api | esp-thinking_shared-network, esp_shared-network |
| sampletrack-web | esp-thinking_shared-network |
| sampletrack-api | esp-thinking_shared-network |
| chisg-lb | chisg-classifier_chisg-net, shared-network |
| mongodb | esp-thinking_shared-network |
| weaviate | esp-thinking_shared-network |
| skills-db | esp-thinking_shared-network |

### Nginx Upstreams (main-proxy)
Config file: `/opt/esp-thinking/main-proxy/nginx-vultr-portfolio.conf`
Upstreams: assist-frontend:80, assist-api:8080, skills-frontend:80, skills-api:8080,
           drb-frontend:80, drb-api:8082, sampletrack-web:80, sampletrack-api:8085,
           chisg-lb:80

### Deploy Workflow
1. Build locally: `docker build -t mike5tew/<service>:latest ./<service>`
2. Push: `docker push mike5tew/<service>:latest`
3. On Vultr: `cd /opt/esp-thinking && docker compose -f docker-compose.prod.yml pull <service> && docker compose -f docker-compose.prod.yml up -d <service>`
""",
    },
]


def connect_weaviate():
    return weaviate.connect_to_local(
        host=WEAVIATE_HOST,
        port=WEAVIATE_PORT,
        grpc_port=WEAVIATE_GRPC_PORT,
    )


def add_to_weaviate(client, entry):
    if not client.collections.exists(DOC_CLASS):
        print(f"  ⚠️  '{DOC_CLASS}' class not found in Weaviate — skipping Weaviate insert")
        return False
    collection = client.collections.get(DOC_CLASS)
    chunk_hash = hashlib.md5(entry["content"].encode()).hexdigest()[:12]
    collection.data.insert({
        "title": entry["title"],
        "content": entry["content"],
        "file_path": entry["file_path"],
        "project": entry["project"],
        "section_path": entry["section_path"],
        "heading_level": 2,
        "chunk_hash": chunk_hash,
        "last_modified": datetime.now().isoformat(),
    })
    return True


def add_to_mongodb(entry):
    client = MongoClient(MONGO_URI, serverSelectionTimeoutMS=5000)
    db = client[MONGO_DB]
    coll = db[NOTES_COLLECTION]
    doc = {
        "title": entry["title"],
        "content": entry["content"],
        "project": entry["project"],
        "section_path": entry["section_path"],
        "file_path": entry["file_path"],
        "type": "infrastructure_change",
        "created_at": datetime.utcnow(),
        "tags": ["infrastructure", "vultr", "docker", "production"],
    }
    result = coll.insert_one(doc)
    client.close()
    return result.inserted_id


def main():
    print("🔌 Connecting to Weaviate (localhost:8088)...")
    wc = connect_weaviate()
    print("✅ Weaviate connected\n")

    for entry in CHANGES:
        print(f"📝 Processing: {entry['title']}")

        # Weaviate
        ok = add_to_weaviate(wc, entry)
        if ok:
            print(f"   ✅ Added to Weaviate '{DOC_CLASS}'")

        # MongoDB
        try:
            inserted_id = add_to_mongodb(entry)
            print(f"   ✅ Added to MongoDB '{NOTES_COLLECTION}' (id: {inserted_id})")
        except Exception as e:
            print(f"   ⚠️  MongoDB insert failed: {e}")

        print()

    wc.close()
    print("✅ Done — all entries logged.")


if __name__ == "__main__":
    main()
