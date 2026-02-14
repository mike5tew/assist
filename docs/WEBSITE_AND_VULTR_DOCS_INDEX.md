# 📚 Website & Vultr Deployment Documentation Index

## Start Here 👇

**New to this deployment?** Start with [QUICK_START.md](QUICK_START.md) (5 min read)

**Have questions?** Check the relevant guide below based on what you need to know.

---

## 🎯 Documentation by Purpose

### I want to understand the overall architecture
→ **[WEBSITE_AND_PANEL_ARCHITECTURE.md](WEBSITE_AND_PANEL_ARCHITECTURE.md)**
- What gets deployed (public website, admin panel, parent portal)
- How the 3 websites are built and deployed
- Nginx routing strategy
- Database architecture
- Frontend build process

**Key sections**:
- Layer 1: Public Website
- Layer 2: Admin Control Panel
- Layer 3: Parent Portal
- Unified Frontend Build Process
- Domain & Routing Strategy

---

### I want the step-by-step deployment instructions
→ **[VULTR_MIGRATION_GUIDE.md](VULTR_MIGRATION_GUIDE.md)**
- Pre-migration checklist
- Server initialization
- Building Docker images on Mac
- Pushing to Docker Hub
- Syncing configs to Vultr
- Starting services
- SSL/HTTPS setup
- Database initialization
- Backup strategy
- Monitoring & health checks
- Troubleshooting

**Key sections**:
- Step 1: Server Initialization
- Step 2: Build & Push Docker Images
- Step 3: Sync Configuration Files
- Step 4: Deploy to Vultr
- Step 5: SSL Setup
- Step 6: Database Initialization
- Step 7: Backup Strategy
- Step 8: Monitoring
- Step 9: DNS Finalization

---

### I want a week-by-week plan with actionable tasks
→ **[WEBSITE_AND_VULTR_CHECKLIST.md](WEBSITE_AND_VULTR_CHECKLIST.md)**
- Phase-by-phase breakdown (10 phases total)
- Checkboxes for each task
- Specific commands to run (copy-paste ready)
- Timeline (3-4 weeks to launch)
- Decision points
- Testing procedures
- Launch checklist

**Key sections**:
- Phase 1: Website Planning & Setup
- Phase 2: Vultr Server Setup
- Phase 3: Build & Push Docker Images
- Phase 4: Deploy to Vultr
- Phase 5: SSL & DNS Setup
- Phase 6: Database Initialization
- Phase 7: Backups & Monitoring
- Phase 8: Content & Testing
- Phase 9: Launch Preparation
- Phase 10: Launch!

---

### I want a quick reference
→ **[QUICK_START.md](QUICK_START.md)**
- Summary of what you're getting
- 3-4 week launch plan
- Key decisions to make
- Quick command reference
- Next immediate steps
- Troubleshooting quick fixes

---

## 📋 File Structure

```
assist/docs/
├── QUICK_START.md ⭐ START HERE
├── WEBSITE_AND_PANEL_ARCHITECTURE.md (Read 2nd — understand architecture)
├── VULTR_MIGRATION_GUIDE.md (Read 3rd — understand deployment)
├── WEBSITE_AND_VULTR_CHECKLIST.md (Use during implementation — actionable tasks)
├── WEBSITE_AND_VULTR_DOCS_INDEX.md (This file)
└── DEPLOYMENT_PLAYBOOK.md (Legacy reference)
```

---

## 🚀 Getting Started in 4 Steps

### 1️⃣ Understand What You're Building (20 min)
Read: [WEBSITE_AND_PANEL_ARCHITECTURE.md](WEBSITE_AND_PANEL_ARCHITECTURE.md)

**You'll learn**:
- What's the public website for?
- What's the admin panel for?
- What's the parent portal for?
- How are they all deployed together?

**Key decision**: Next.js or Vite for public site? Netlify or Vultr for static hosting?

---

### 2️⃣ Understand How to Deploy (30 min)
Read: [VULTR_MIGRATION_GUIDE.md](VULTR_MIGRATION_GUIDE.md) Steps 1-5

**You'll learn**:
- How to set up a Vultr server
- How to build Docker images on Mac
- How to push images to Docker Hub
- How to deploy everything to Vultr
- How to set up SSL

**Key decision**: Vultr instance size? (2GB RAM ≈ $12/month)

---

### 3️⃣ Get Your Week-by-Week Plan (15 min)
Read: [WEBSITE_AND_VULTR_CHECKLIST.md](WEBSITE_AND_VULTR_CHECKLIST.md) Phases 1-2

**You'll learn**:
- What to do this week (planning + setup)
- What to do next week (build + deploy)
- Which decisions need to be made first
- Exact commands to run

**Key decision**: Timeline — are you starting now?

---

### 4️⃣ Pick a Phase and Start! 🎯
Use: [WEBSITE_AND_VULTR_CHECKLIST.md](WEBSITE_AND_VULTR_CHECKLIST.md) as your daily todo

**Options**:
- **Phase 1a**: Start with website architecture
- **Phase 2a**: Start with Vultr server setup
- **Phase 3a**: Start with Docker image building
- **Phase 4a**: Start with deployment

---

## 💡 Key Concepts

### Build Local, Deploy to Vultr
```
Mac (local development)
  ↓
Build Docker images (cross-compile for linux/amd64)
  ↓
Push to Docker Hub
  ↓
Vultr Server
  ↓
docker-compose pull
docker-compose up -d
  ↓
Services running (Nginx routes traffic)
```

### Three Websites, One Backend
```
smartminds.education (Nginx reverse proxy)
├── / → smartminds-marketing (Next.js public site)
├── /app → smartminds-admin (React admin panel)
├── /parent → smartminds-parent (React parent portal)
└── /api → assist-api (Go backend)
```

### Databases Are Containerized
```
Vultr Server
├── mongodb (semantic links)
├── skills-mysql (users, progress, courses)
├── weaviate (assist docs/ideas)
└── weaviate-chisg (579 skills + 1120 links + ETPs)
```

---

## ❓ FAQ

**Q: What if I don't want to use Vultr?**
A: The same docker-compose.prod.yml works on any Linux server. Just update the IP in Step 4 of the checklist.

**Q: Can I deploy just the backend first?**
A: Yes! Just skip the website phases (1, 8) and focus on Phases 2-7 (server, docker, deployment, databases).

**Q: How long does deployment actually take?**
A: From fresh server to running services: ~30 minutes. Full launch with testing: 3-4 weeks.

**Q: What if something breaks in production?**
A: Check the Troubleshooting section in VULTR_MIGRATION_GUIDE.md or see the Quick Reference in QUICK_START.md.

**Q: Can I update services without downtime?**
A: Yes! See "Scenario 2: Updating a Single Service" in DEPLOYMENT_PLAYBOOK.md.

**Q: How much will hosting cost?**
A: Vultr 2GB RAM = $12/month. Add Object Storage for backups = ~$5/month. Total ≈ $20/month.

---

## 🔗 Related Documentation

- [PROJECT_INDEX.md](../PROJECT_INDEX.md) — Overall project overview
- [PROJECT_STRUCTURE.md](../PROJECT_STRUCTURE.md) — Code organization
- [DEPLOYMENT_PLAYBOOK.md](DEPLOYMENT_PLAYBOOK.md) — Legacy (original deployment guide)
- [PRODUCT_VISION.md](PRODUCT_VISION.md) — Why this system exists

---

## 📞 Quick Links

- **Vultr Docs**: https://www.vultr.com/docs/
- **Docker Docs**: https://docs.docker.com/
- **Docker Compose Reference**: https://docs.docker.com/compose/compose-file/
- **Let's Encrypt**: https://letsencrypt.org/
- **Nginx Docs**: https://nginx.org/en/docs/

---

## ✅ Checklist: You're Ready to Deploy When You've

- [ ] Read QUICK_START.md
- [ ] Read WEBSITE_AND_PANEL_ARCHITECTURE.md
- [ ] Read VULTR_MIGRATION_GUIDE.md Steps 1-5
- [ ] Answered the 8 decisions in WEBSITE_AND_VULTR_CHECKLIST.md
- [ ] Created .env.production with all secrets
- [ ] Set up Docker Hub account
- [ ] Registered your domain (or decided on domain name)
- [ ] Created Vultr account

**Then**: Pick a phase from WEBSITE_AND_VULTR_CHECKLIST.md and start! 🚀

---

*Last Updated: 6 February 2026*
