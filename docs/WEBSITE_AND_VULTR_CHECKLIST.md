# Website & Vultr Deployment — Actionable Checklist

## Phase 1: Website Planning & Setup (This Week)

### 1a. Finalize Website Architecture
- [ ] Review [WEBSITE_AND_PANEL_ARCHITECTURE.md](WEBSITE_AND_PANEL_ARCHITECTURE.md)
- [ ] Decide: Next.js or Remix for public site?
- [ ] Decide: Contentful CMS or markdown-based content?
- [ ] Decide: Netlify for static sites + Vultr for backend, or everything on Vultr?

### 1b. Create Website Directory Structure
```bash
cd /Users/michaelstewart/Coding/assist

# Create websites directory
mkdir -p websites/{smartminds-marketing,smartminds-admin,smartminds-parent}

# Initialize each as separate project
cd websites/smartminds-marketing
npm create vite@latest . -- --template react-ts  # or Next.js
```

### 1c. Design Wireframes
- [ ] Public homepage (hero, features, pricing, CTA)
- [ ] Features page (CHISG explorer, ETP spectra, parent app overview)
- [ ] Pricing page (free tier, teacher, school, researcher, enterprise)
- [ ] Researcher portal page (CHISG explorer, API docs, download)
- [ ] Admin panel dashboard (school overview, teacher management, analytics)
- [ ] Parent portal child profile (skill progress, ETP sliders, milestones)

### 1d. Determine Domain & Email
- [ ] Register domain: `smartminds.education` (or preferred)
- [ ] Set up email: support@smartminds.education (Proton Mail, Gmail)
- [ ] Configure DNS MX records

---

## Phase 2: Vultr Server Setup (Next Week)

### 2a. Create Vultr Account & Server
- [ ] Create Vultr account
- [ ] Create 2GB+ RAM instance (Ubuntu 22.04 LTS)
- [ ] Note server IP: `__________________`
- [ ] Configure firewall (SSH 22, HTTP 80, HTTPS 443)
- [ ] Create non-root deploy user

### 2b. Initialize Server
```bash
# Run the commands from VULTR_MIGRATION_GUIDE.md Step 1
ssh root@<VULTR_IP>

# These will be copied to your server:
apt update && apt upgrade -y
apt install -y docker.io docker-compose
mkdir -p /opt/smart-minds
```

### 2c. Prepare .env Files
- [ ] Create `/Users/michaelstewart/Coding/assist/env.production` (NEVER commit)
- [ ] Fill in all secrets:
  - `MONGODB_PASSWORD=________`
  - `MYSQL_ROOT_PASSWORD=________`
  - `JWT_SECRET=________`
  - `DOCKER_HUB_USERNAME=________`
  - `DOCKER_HUB_PASSWORD=________`

---

## Phase 3: Build & Push Docker Images (Next Week)

### 3a. Build All Services Locally (Mac)
```bash
cd /Users/michaelstewart/Coding/assist

# Build using docker buildx (cross-compile for linux/amd64)
# Commands from VULTR_MIGRATION_GUIDE.md Step 2b

# Verify all images are on Docker Hub
# https://hub.docker.com/repositories
```

- [ ] smartminds-proxy
- [ ] assist-frontend
- [ ] assist-api
- [ ] smartminds-admin
- [ ] smartminds-parent
- [ ] smartminds-marketing

### 3b. Test Locally (Optional but Recommended)
```bash
docker-compose -f docker-compose.prod.yml up -d
# Visit http://localhost and verify all services load
```

---

## Phase 4: Deploy to Vultr (Next Week)

### 4a. Sync Configuration Files
```bash
# Run from Mac in /Users/michaelstewart/Coding/assist/
# Commands from VULTR_MIGRATION_GUIDE.md Step 3

REMOTE_HOST="root@<VULTR_IP>"
rsync -avz docker-compose.prod.yml "${REMOTE_HOST}:/opt/smart-minds/docker-compose.yml"
rsync -avz env.production "${REMOTE_HOST}:/opt/smart-minds/.env"
rsync -avz main-proxy/nginx.conf "${REMOTE_HOST}:/opt/smart-minds/main-proxy/nginx.conf"
```

### 4b. Start Services on Vultr
```bash
ssh root@<VULTR_IP>
cd /opt/smart-minds
docker-compose pull
docker-compose up -d
docker-compose ps  # Verify all services running
```

- [ ] main-proxy is running
- [ ] assist-api is running and responding to `/health`
- [ ] mongodb is running
- [ ] skills-mysql is initialized
- [ ] weaviate instances are responding

### 4c. Test via IP Address
```bash
# Test from Mac
curl http://<VULTR_IP>/api/health
curl http://<VULTR_IP>/

# If using Postman or Insomnia, test:
# GET http://<VULTR_IP>/api/chisg/skills
# GET http://<VULTR_IP>/api/parent/children/{child_id}/profile
```

- [ ] Homepage loads
- [ ] API endpoints respond
- [ ] Admin panel loads
- [ ] Parent portal loads

---

## Phase 5: SSL & DNS Setup (Next Week)

### 5a. Get SSL Certificates
```bash
ssh root@<VULTR_IP>

apt install -y certbot python3-certbot-nginx

certbot certonly --standalone \
  -d smartminds.education \
  -d www.smartminds.education \
  -d api.smartminds.education \
  -d app.smartminds.education \
  -d researcher.smartminds.education
```

### 5b. Update DNS Records
In your domain registrar, set:

```
smartminds.education       A  <VULTR_IP>
www.smartminds.education   CNAME smartminds.education
api.smartminds.education   CNAME smartminds.education
app.smartminds.education   CNAME smartminds.education
researcher.smartminds.education CNAME smartminds.education
```

### 5c. Update Nginx Config
- [ ] Add SSL certificates to `main-proxy/nginx.conf`
- [ ] Redirect HTTP → HTTPS
- [ ] Upload updated config: `rsync -avz main-proxy/nginx.conf "${REMOTE_HOST}:/opt/smart-minds/main-proxy/nginx.conf"`
- [ ] Restart Nginx: `docker-compose restart main-proxy`

---

## Phase 6: Database Initialization (Next Week)

### 6a. Initialize MySQL
```bash
ssh root@<VULTR_IP>
cd /opt/smart-minds

# Create skillsmap database and tables
docker-compose exec skills-mysql mysql -u root -p${MYSQL_ROOT_PASSWORD} << 'EOF'
CREATE DATABASE skillsmap;
CREATE TABLE skillsmap.skills (...);
CREATE TABLE skillsmap.users (...);
-- See skills-map-platform/init.sql for full schema
EOF
```

### 6b. Initialize MongoDB
```bash
docker-compose exec mongodb mongosh -u admin -p${MONGODB_PASSWORD} --authenticationDatabase admin << 'EOF'
use esp_organizer;
db.semantic_links.createIndex({ "source_skill_id": 1, "target_skill_id": 1 });
db.semantic_links.createIndex({ "relationship_type": 1 });
EOF
```

### 6c. Verify Databases
- [ ] Can connect to MySQL and see databases
- [ ] Can connect to MongoDB and see collections
- [ ] Weaviate-CHISG has skills loaded (curl `http://localhost:8081/v1/objects`)

---

## Phase 7: Backups & Monitoring (End of Week)

### 7a. Set Up Automated Backups
```bash
ssh root@<VULTR_IP>

# Install s3cmd for Vultr Object Storage
apt install -y s3cmd
s3cmd --configure

# Create backup script (from VULTR_MIGRATION_GUIDE.md)
cat > /opt/smart-minds/backup.sh << 'EOF'
#!/bin/bash
# ... backup script contents ...
EOF

chmod +x /opt/smart-minds/backup.sh

# Schedule daily backup at 2 AM
crontab -e
# 0 2 * * * /opt/smart-minds/backup.sh >> /var/log/smart-minds/backup.log 2>&1
```

- [ ] Test backup script runs: `./backup.sh`
- [ ] Verify backups uploaded to Vultr Object Storage
- [ ] Test restoration procedure (restore from backup)

### 7b. Set Up Monitoring
- [ ] Enable container health checks in docker-compose.prod.yml
- [ ] Configure log aggregation (optional: Sentry, DataDog, or just local logs)
- [ ] Set up cron job to check disk space: `df -h | mail -s "Disk Space Alert" admin@smartminds.education`

### 7c. Document Runbooks
- [ ] Create `.md` file: "What to do if assist-api crashes"
- [ ] Create `.md` file: "How to update a single service"
- [ ] Create `.md` file: "How to restore from backup"

---

## Phase 8: Content & Testing (Ongoing)

### 8a. Create Website Content
- [ ] Write homepage copy
- [ ] Create case studies (at least 2)
- [ ] Write FAQ
- [ ] Create API documentation
- [ ] Screenshot SkillsMarkbookMobile for marketing
- [ ] Record demo video (optional)

### 8b. Build Admin Panel Features
- [ ] Teacher management CRUD
- [ ] Student management CRUD
- [ ] Analytics dashboard (skill adoption rates, ETP distribution)
- [ ] Export to CSV/PDF
- [ ] School billing management

### 8c. Build Parent Portal Features
- [ ] Child profile view
- [ ] Skill explorer (with semantic link explanations)
- [ ] ETP spectrum explorer
- [ ] Milestones view
- [ ] Messages from teacher (optional Phase 2)

### 8d. User Testing
- [ ] Test with 2-3 real teachers (via SkillsMarkbookMobile)
- [ ] Test with 2-3 real parents (via parent portal)
- [ ] Iterate on feedback

---

## Phase 9: Launch Preparation (Final Week)

### 9a. Security Audit
- [ ] Enable UFW firewall on Vultr server
- [ ] Disable root SSH login
- [ ] Set up SSH key authentication
- [ ] Review all environment variables (no secrets in code)
- [ ] Audit database passwords (strong + unique)
- [ ] Test rate limiting on API endpoints
- [ ] Verify SSL/TLS configuration (use SSL Labs test)

### 9b. Performance Testing
```bash
# Load test the API
ab -n 1000 -c 10 https://api.smartminds.education/health

# Test database query performance
# Check slow query logs in MySQL
```

- [ ] API responds < 200ms for typical requests
- [ ] Database queries optimized (indexes created)
- [ ] Frontend loads < 3 seconds

### 9c. Documentation Audit
- [ ] All documentation is up-to-date
- [ ] API documentation is complete (OpenAPI/Swagger)
- [ ] User guides written (teacher, parent, admin)
- [ ] Troubleshooting guide created
- [ ] README updated with Vultr info

### 9d. Final Testing Checklist
- [ ] Homepage loads on all devices (mobile, tablet, desktop)
- [ ] Admin panel works (create teacher, add student, view analytics)
- [ ] Parent portal works (view child profile, explore skills)
- [ ] API endpoints all respond correctly
- [ ] Backups run successfully
- [ ] SSL certificate is valid
- [ ] All 4 data stores are accessible

---

## Phase 10: Launch! 🚀

### 10a. Update DNS
- [ ] Point domain to Vultr server
- [ ] Wait for DNS propagation (5-30 minutes)
- [ ] Test that smartminds.education resolves

### 10b. Monitor First 24 Hours
- [ ] Watch error logs in real-time
- [ ] Monitor disk space usage
- [ ] Check API response times
- [ ] Be ready to rollback if issues

### 10c. Send Announcements
- [ ] Email teachers: "New SkillsMarkbookMobile release + web access"
- [ ] Email parents: "You can now view your child's progress"
- [ ] Post on social media (if applicable)

### 10d. Ongoing Support
- [ ] Monitor logs daily for first week
- [ ] Update health status page with any incidents
- [ ] Respond to user feedback

---

## 📋 Summary: Timeline

| Phase | Duration | What Gets Done |
|-------|----------|---|
| Phase 1 | 2-3 days | Website architecture finalized, wireframes created |
| Phase 2 | 1-2 days | Vultr account created, server initialized |
| Phase 3 | 2-3 days | All Docker images built and pushed |
| Phase 4 | 1 day | Services deployed and tested on Vultr IP |
| Phase 5 | 1 day | SSL certificates obtained, DNS configured |
| Phase 6 | 1 day | Databases initialized and verified |
| Phase 7 | 1-2 days | Backups set up, monitoring enabled |
| Phase 8 | 1-2 weeks | Website content, admin panel, parent portal built |
| Phase 9 | 2-3 days | Security audit, performance testing, documentation |
| Phase 10 | Ongoing | Launch and monitor |

**Total: ~3-4 weeks from start to production launch**

---

## 🔗 Reference Documents

- [WEBSITE_AND_PANEL_ARCHITECTURE.md](WEBSITE_AND_PANEL_ARCHITECTURE.md) — Full architecture for all three websites
- [VULTR_MIGRATION_GUIDE.md](VULTR_MIGRATION_GUIDE.md) — Complete deployment walkthrough
- [DEPLOYMENT_PLAYBOOK.md](DEPLOYMENT_PLAYBOOK.md) — Legacy reference (still valid)
- [PROJECT_INDEX.md](../PROJECT_INDEX.md) — Project overview and status

---

## ❓ Questions to Answer Before Starting

1. What's your budget for Vultr hosting? (2GB RAM ≈ $12/month, scalable to larger)
2. Do you want the websites on Vultr or on separate hosting (Netlify, Vercel)?
3. Should the parent portal be separate domain (parents.smartminds.education) or subpath (/parent)?
4. Do you need email authentication or OAuth (Google/Microsoft)?
5. Will you use a CMS for website content or manage it in Git?
6. Do you need analytics? (Posthog, Plausible, Google Analytics)
7. Do you need uptime monitoring / alerting? (Uptime Robot, PagerDuty)
8. Will parents and teachers have account signup, or will admins provision accounts?

---

*Next step: Review Phase 1 items and let me know when you're ready to proceed!*
