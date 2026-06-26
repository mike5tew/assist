# Operations Log

Record of production deployment fixes and configuration changes on the Vultr VPS (192.248.151.185).

---

## 26 April 2026 — skills-api/frontend Deploy + Network Fix + Graph Layout Fix

### What Changed

**Docker network — permanent fix** (`/opt/esp-thinking/docker-compose.yml`):
- `skills-api` and `skills-frontend` were being created only on `esp-thinking_shared-network` by compose, but `main-proxy` and `skills-db` live on `esp_shared-network` (the `esp` compose project)
- Every deploy caused a 502 (proxy can't reach frontend) and skills-db DNS failure until `docker network connect esp_shared-network <container>` was run manually
- Fix: added `esp_shared-network` as an external network in `/opt/esp-thinking/docker-compose.yml` and added it to `skills-frontend` and `skills-api` services
- **Future deploys of skills-api and skills-frontend no longer need manual `docker network connect`**

```yaml
# Add to both skills-frontend and skills-api services:
networks:
  - shared-network
  - esp_shared-network

# Add to bottom-level networks block:
networks:
  shared-network:
    driver: bridge
  esp_shared-network:
    external: true
    name: esp_shared-network
```

**Also fixed this session — esp-thinking `.env` DB credentials** (`/opt/esp-thinking/.env`):
- `DB_NAME` was `skills_db` → corrected to `dare2lead` (matches `skills-db` container's actual init)
- `DB_PASSWORD` was `7oEW1UZGUJJnRFCnK0ak` → corrected to `xiMqun-sezgob-virjo9`
- Root cause: `skills-db` is managed by the `esp` compose project (uses `/opt/esp/.env`), not `esp-thinking`

**skills-api deployed** (limit: 1200 → 2000 fix):
- `api/weaviate/client.go` `GetAllCHISGSkillsAndLinks()` `SkillLink(limit: 2000)` deployed
- Old standalone container removed, compose-managed container started
- Resolves 23 isolated nodes on `/skillstree/explore`

**skills-frontend deployed** (BFS layout fix):
- `frontend/src/components/GraphExplorer.tsx` — replaced Longest Path depth algorithm with BFS minimum-depth
- Old algorithm: 74 levels × 280px = 20,440px wide (caused extreme map stretch after all 597 nodes connected)
- New algorithm: 7 levels × 280px = 1,680px wide
- Old standalone container removed, compose-managed container started

### Key Infrastructure Facts (critical for future sessions)

- **Two compose projects on Vultr**:
  - `/opt/esp/docker-compose.yml` → project `esp` → manages: `skills-db`, `weaviate`, `mongodb`, `main-proxy`, `drb-api`, `drb-frontend`, `chisg-*`
  - `/opt/esp-thinking/docker-compose.yml` → project `esp-thinking` → manages: `skills-api`, `skills-frontend`, `assist-api`, `assist-frontend`, `sampletrack-*`
- `skills-db` credentials come from `/opt/esp/.env` (NOT `/opt/esp-thinking/.env`)  
  - DB_NAME=dare2lead, DB_USER=skills_user, DB_PASSWORD=xiMqun-sezgob-virjo9
- **Deploy pattern for skills-api or skills-frontend**:
  ```bash
  # Build locally:
  cd /Users/michaelstewart/Coding && docker buildx build --platform linux/amd64 \
    -t mike5tew/skills-api:latest -f skills-map-platform/api/Dockerfile --push .
  # Or for frontend:
  docker buildx build --platform linux/amd64 \
    -t mike5tew/skills-frontend:latest -f skills-map-platform/frontend/Dockerfile \
    --push skills-map-platform/frontend
  # On Vultr (no manual network connect needed after compose fix):
  docker stop skills-api && docker rm skills-api  # if orphan exists
  cd /opt/esp-thinking && docker compose pull skills-api && docker compose up -d --no-deps skills-api
  ```

### Container State (end of 26 Apr session)

| Container | Project | Networks | Status |
|-----------|---------|----------|--------|
| skills-api | esp-thinking | esp-thinking_shared-network + esp_shared-network | Up, connected |
| skills-frontend | esp-thinking | esp-thinking_shared-network + esp_shared-network | Up, connected |
| skills-db | esp | esp_shared-network | Up |
| weaviate (8088) | esp | esp_shared-network | Up, 597 CHISGElement + 1753 SkillLink |
| main-proxy | standalone | esp_shared-network + esp-thinking_shared-network | Up |
| drb-api | standalone | esp_shared-network | Up |

---

## 26 April 2026 — CHISG Graph Cleanup + Production Weaviate Seeding

### What Changed

**CHISG Knowledge Graph** (skills-map-platform):
- Graph had 1828 links with 142 cycles. 75 reversed edges removed across 4 passes → **1753 links, 0 cycles (DAG)**
- All 597 skills confirmed to have descriptions
- `data/chisg_skill_tree.json` rebuilt from live Weaviate 8081: 597 skills, descriptions, domains, parents[], offspring[]
- CHISG triple extraction run via AWS Bedrock (`eu.anthropic.claude-haiku-4-5-20251001-v1:0`): 540/597 skills processed, 1705 triples → `data/chisg/skill_triples.json`

**Production Weaviate 8088 seeded** (was empty — causing 500 errors on /skillstree/explore):
- Schema recreated: `CHISGElement`, `SkillLink` (with `forward_link`/`backward_link`), `CourseSkillSuggestions`
- 597 CHISGElement seeded, 1753 SkillLink seeded
- Script: `skills-map-platform/scripts/seed_production_weaviate.py --no-vectorizer`
- Source: local Weaviate 8081 (humanos-weaviate-1 container)

**Bug found but NOT yet deployed — skills-api SkillLink limit**:
- `api/weaviate/client.go` `GetAllCHISGSkillsAndLinks()` queried `SkillLink(limit: 1200)` but graph has 1753 links
- Result: 23 skills appeared isolated in the /explore graph (their only connections were in the missing 553 links)
- **Fix**: changed to `limit: 2000` in source code
- **Still needed**: rebuild Docker image `mike5tew/skills-api`, push, deploy on Vultr
- **Deploy command** (safe — only restarts skills-api, not DRB or proxy):
  ```bash
  cd /opt/esp-thinking
  docker compose pull skills-api
  docker compose up -d --no-deps skills-api
  ```

### New Scripts Created

| Script | Purpose |
|--------|--------|
| `skills-map-platform/scripts/seed_production_weaviate.py` | Full seed pipeline: reads local Weaviate 8081, builds vector_content, seeds target Weaviate. Flags: `--dry-run`, `--no-vectorizer`, `--no-drop`, `--local` |
| `skills-map-platform/scripts/fix_long_cycles.py` | Cycle fix pass 1 (removed 13 edges) |
| `skills-map-platform/scripts/fix_cycles_pass2.py` | Cycle fix pass 2 (removed 21 edges) |
| `skills-map-platform/scripts/extract_skill_triples.py` | Extract CHISG triples from skill descriptions via Bedrock |

### Container State (post-session)

| Container | Status | Notes |
|-----------|--------|-------|
| weaviate (8088) | Up | Seeded: 597 CHISGElement, 1753 SkillLink |
| skills-api | Up (OLD IMAGE) | Needs rebuild for limit: 1200 → 2000 fix |
| skills-frontend | Up | No changes |
| drb-api | Up (standalone) | Not compose-managed |
| drb-frontend | Up (compose: esp) | Not compose-managed in esp-thinking |
| main-proxy | Up (standalone) | Not compose-managed |

---

## 19 April 2026 — SampleTrack Deployment (Web + API + TestFlight)

### What Was Deployed

**SampleTrack** — NHS Blood Science Sample Tracker (Go API + React Web + React Native Mobile).

- **Web dashboard**: https://espthinking.co.uk/sampletrack/
- **API**: https://espthinking.co.uk/sampletrack/api/
- **Mobile**: Submitted to TestFlight (ASC App ID: 6762559687)

### Docker Images Pushed

| Image | Platform | Builder |
|-------|----------|---------|
| `mike5tew/sampletrack-api:latest` | linux/amd64 | docker buildx (multiplatform) |
| `mike5tew/sampletrack-web:latest` | linux/amd64 | docker buildx (multiplatform) |

### Changes to docker-compose.prod.yml

Added `sampletrack-web` and `sampletrack-api` services:

```yaml
sampletrack-web:
  image: mike5tew/sampletrack-web:latest
  container_name: sampletrack-web
  restart: always
  networks:
    - shared-network

sampletrack-api:
  image: mike5tew/sampletrack-api:latest
  container_name: sampletrack-api
  restart: always
  expose:
    - "8085"
  environment:
    - MONGO_URI=mongodb://mongodb:27017/sampletrack
    - DB_NAME=sampletrack
    - JWT_SECRET=${JWT_SECRET}
    - PORT=8085
  depends_on:
    - mongodb
  networks:
    - shared-network
```

**Key**: `MONGO_URI` uses NO authentication (consistent with existing mongodb container which has no auth — see 16 April entry).

### Changes to nginx-vultr-portfolio.conf

Added upstream definitions and location blocks for SampleTrack in both HTTP and HTTPS server blocks:

```nginx
# Upstreams
upstream sampletrack_frontend {
    server sampletrack-web:80;
}
upstream sampletrack_api {
    server sampletrack-api:8085;
}

# Location blocks (in both server blocks)
location /sampletrack/api/ {
    rewrite ^/sampletrack/api/(.*)$ /api/$1 break;
    proxy_pass http://sampletrack_api;
    # standard proxy headers...
}
location /sampletrack/ {
    rewrite ^/sampletrack/(.*)$ /$1 break;
    proxy_pass http://sampletrack_frontend;
    # standard proxy headers...
}
```

### Deployment Steps on Vultr

```bash
# SSH to Vultr (password auth — key ~/.ssh/gitkey11-25 NOT working)
ssh root@192.248.151.185

# Configs were copied from Mac via scp:
# - docker-compose.prod.yml → /opt/esp-thinking/docker-compose.yml
# - nginx-vultr-portfolio.conf → /opt/esp/main-proxy/nginx-vultr-portfolio.conf
#   ⚠️ CRITICAL: nginx config must go to /opt/esp/main-proxy/ (NOT /opt/esp-thinking/main-proxy/)
#   because main-proxy container volume mounts from /opt/esp/

# Started containers using --no-deps (mongodb owned by different compose project)
cd /opt/esp-thinking
docker compose pull sampletrack-api sampletrack-web
docker compose up -d --no-deps sampletrack-api sampletrack-web

# Reloaded nginx to pick up new config
docker exec main-proxy nginx -s reload
```

### Production Database Seeded

```bash
# Admin user created
curl -X POST https://espthinking.co.uk/sampletrack/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123","role":"admin"}'

# 15 NHS locations seeded (Phlebotomy, Pathology Lab, A&E, etc.)
```

### iOS App — EAS Build & TestFlight Submission

- **EAS Project ID**: `7c6bc804-ca9b-4456-ae14-1249d4e7f2ee`
- **Apple Team**: Michael Stewart (5EKL6W683T), Individual account
- **Apple ID**: mike5tew@hotmail.com
- **Distribution Cert**: serial `62EA063879FF35276EF2FA8407B460F9`
- **ASC App ID**: `6762559687`
- **App Name on ASC**: "SampleTrack (07979e)" (name "SampleTrack" was taken)
- **TestFlight Group**: "Team (Expo)" — mike5tew@hotmail.com, expertanswerz@gmail.com

Build issues resolved during session:
1. npm peer dependency conflicts → added `.npmrc` with `legacy-peer-deps=true`
2. Missing `./assets/icon.png` → generated placeholder blue (0d47a1) PNGs
3. Empty `ascAppId` in `eas.json` → removed, replaced with `appleId` field

### Local File Changes (SampleTrack repo)

| File | Change |
|------|--------|
| `GoLangAPI/` | Docker image built and pushed (no code changes) |
| `ReactWebFrontEnd/` | Docker image built and pushed (no code changes) |
| `ReactNativeMobileFrontEnd/eas.json` | Removed empty ascAppId, set appleId |
| `ReactNativeMobileFrontEnd/.npmrc` | Created — `legacy-peer-deps=true` |
| `ReactNativeMobileFrontEnd/.gitignore` | Added `ios/`, `android/`, `.expo/` |
| `ReactNativeMobileFrontEnd/assets/` | Created placeholder icon.png, adaptive-icon.png, splash-icon.png |
| `ReactNativeMobileFrontEnd/app.json` | Updated by `eas init` with projectId |

### Current Container State (post-deployment)

| Container | Image | Status | Notes |
|-----------|-------|--------|-------|
| sampletrack-web | mike5tew/sampletrack-web:latest | Up | Compose-managed (esp-thinking) |
| sampletrack-api | mike5tew/sampletrack-api:latest | Up | Compose-managed (esp-thinking), port 8085 |
| main-proxy | nginx:alpine | Up | Standalone (NOT compose-managed — see 16 April) |
| drb-api | mike5tew/drb-api:latest | Up | Standalone (NOT compose-managed — see 16 April) |
| drb-frontend | mike5tew/drb-frontend:latest | Up | Compose-managed (esp) |
| assist-frontend | mike5tew/assist-frontend:latest | Up | Compose-managed (esp-thinking) |
| assist-api | mike5tew/assist-api:latest | Up | Compose-managed (esp-thinking) |
| mongodb | mongo:latest | Up | NO AUTH — shared by sampletrack, drb, assist |
| skills-frontend | mike5tew/skills-frontend:latest | Up | Compose-managed (esp) |
| skills-api | mike5tew/skills-api:latest | Up | Compose-managed (esp) |
| skills-db | mysql:8.0 | Up | Compose-managed (esp) |
| weaviate | semitechnologies/weaviate | Up | Compose-managed (esp) |

### Known Issues / Follow-up

1. **Placeholder app icons**: iOS app uses solid blue placeholder PNGs — replace with proper branding before wider distribution
2. **Admin password**: Production admin user is `admin` / `admin123` — should be changed
3. **SSH key auth**: `~/.ssh/gitkey11-25` is not accepted by Vultr — using password auth as workaround
4. **MongoDB still has no auth** — sampletrack joins drb and assist in using unauthenticated connections
5. **ASC app name**: Registered as "SampleTrack (07979e)" — can be renamed at https://appstoreconnect.apple.com

---

## 16 April 2026 — DRB Deployment & Network Fix

### Problem
- DRB dashboard at `espthinking.co.uk/drb/` returning 404
- The production `docker-compose.yml` on the server was stale (no DRB services)
- After syncing updated compose, the `main-proxy` nginx container got stuck in a crash loop due to Docker overlay2 corruption on the volume mount

### Root Causes Found

1. **Two Docker networks exist on the server:**
   - `esp-thinking_shared-network` — used by `assist-frontend`, `assist-api`, `mongodb`
   - `esp_shared-network` — used by `skills-frontend`, `skills-api`, `skills-db`, `drb-frontend`, `chisg-*`, `weaviate`, `chisg-redis`
   - The `main-proxy` must be connected to **both** networks to resolve all upstream hostnames

2. **Nginx config filename mismatch:**
   - Compose expects: `./main-proxy/nginx-vultr-portfolio.conf` mounted to `/etc/nginx/conf.d/default.conf`
   - The DEPLOYMENT_PLAYBOOK.md references `nginx.conf` — this is outdated
   - Correct file on server: `/opt/esp/main-proxy/nginx-vultr-portfolio.conf`

3. **MongoDB has no authentication enabled:**
   - `MONGO_INITDB_ROOT_USERNAME` and `MONGO_INITDB_ROOT_PASSWORD` are both empty on the running `mongodb` container
   - The compose file references `${MONGO_ROOT_USER}` and `${MONGO_ROOT_PASSWORD}` which don't exist in `.env`
   - Connection URIs with auth credentials cause `MongoServerError: Authentication failed`
   - DRB API must use a plain URI: `mongodb://mongodb:27017/drb_monitor`

4. **DRB API panic on nil cursor:**
   - The Go code does `cursor, _ := col.Find(...)` then `cursor.All(...)` — if `Find` returns an error (e.g. auth failure), cursor is nil and causes a panic
   - This is a code quality issue to fix later, but the immediate fix was correcting the MongoDB URI

### Resolution Steps

Since `docker compose` had corrupted overlay2 state for `main-proxy`, we bypassed compose and ran containers manually:

```bash
# 1. Sync updated configs from Mac
cd ~/Coding/assist
rsync -avz docker-compose.prod.yml root@192.248.151.185:/opt/esp/docker-compose.yml
rsync -avz main-proxy/nginx-vultr-portfolio.conf root@192.248.151.185:/opt/esp/main-proxy/nginx-vultr-portfolio.conf

# 2. Create main-proxy on BOTH networks
ssh root@192.248.151.185
docker stop main-proxy && docker rm main-proxy
docker create --name main-proxy \
  --network esp-thinking_shared-network \
  -p 80:80 -p 443:443 \
  -v /opt/esp/main-proxy/nginx-vultr-portfolio.conf:/etc/nginx/conf.d/default.conf:ro \
  -v esp-thinking_certbot_www:/var/www/certbot:ro \
  -v esp-thinking_certbot_certs:/etc/letsencrypt:ro \
  --restart always nginx:alpine
docker network connect esp_shared-network main-proxy

# 3. Create drb-api on BOTH networks with NO-AUTH MongoDB URI
docker create --name drb-api \
  --network esp-thinking_shared-network \
  -e 'MONGO_URI=mongodb://mongodb:27017/drb_monitor' \
  -e 'PORT=8082' \
  --expose 8082 \
  --restart always mike5tew/drb-api:latest
docker network connect esp_shared-network drb-api

# 4. Start both
docker start drb-api
docker start main-proxy
```

### Current Container State (post-fix)

| Container | Image | Network(s) | Status |
|-----------|-------|------------|--------|
| main-proxy | nginx:alpine | esp-thinking_shared-network + esp_shared-network | Up (standalone, NOT compose-managed) |
| drb-api | mike5tew/drb-api:latest | esp-thinking_shared-network + esp_shared-network | Up (standalone, NOT compose-managed) |
| drb-frontend | mike5tew/drb-frontend:latest | esp_shared-network | Up (compose-managed) |
| assist-frontend | mike5tew/assist-frontend:latest | esp-thinking_shared-network | Up |
| assist-api | mike5tew/assist-api:latest | esp-thinking_shared-network | Up |
| mongodb | mongo:latest | esp-thinking_shared-network | Up (NO AUTH) |
| skills-frontend | mike5tew/skills-frontend:latest | esp_shared-network | Up |
| skills-api | mike5tew/skills-api:latest | esp_shared-network | Up |
| skills-db | mysql:8.0 | esp_shared-network | Up |
| chisg-1/2/3 | — | esp_shared-network | Up |
| chisg-lb | — | esp_shared-network | Up |
| chisg-redis | — | esp_shared-network | Up |
| weaviate | — | esp_shared-network | Up |

### Known Issues / Follow-up Required

1. **main-proxy and drb-api are NOT compose-managed** — they were created with `docker run`/`docker create`. Running `docker compose up -d` may conflict. Before running compose again, stop and remove these two containers first.

2. **Two Docker networks need consolidation** — All services should ideally be on one network. This split happened because compose was run with different project names at different times (`esp-thinking` vs `esp`).

3. **MongoDB has no authentication** — The `.env` file has placeholder credentials but they were never applied to the running MongoDB instance. The existing data volume was created without auth. Enabling auth now would require re-initializing or creating users manually.

4. **Nginx config filename in DEPLOYMENT_PLAYBOOK.md is outdated** — References `nginx.conf` but the actual file is `nginx-vultr-portfolio.conf`. The rsync command in Scenario 3 needs updating.

5. **DRB API error handling** — `main.go` ignores errors from MongoDB `Find()` calls (`cursor, _ := ...`), causing panics on nil cursors. Should be fixed in code.

### Verified Working Endpoints
- `https://espthinking.co.uk/` — Portfolio site ✅
- `https://espthinking.co.uk/drb/` — DRB Dashboard (React SPA) ✅
- `https://espthinking.co.uk/drb/api/health/ping` — "DRB Monitor API: Active (Mongo Enabled)" ✅
- `https://espthinking.co.uk/drb/api/schools` — Returns 13 schools with data ✅
- `https://espthinking.co.uk/drb/api/estates/contractors` — Returns contractor data ✅
