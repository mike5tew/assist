# Vultr Migration & Deployment Guide

## Overview

This guide covers migrating the entire Smart Minds ecosystem to Vultr Cloud Compute. The deployment strategy is:

- **Build locally on Mac** (cross-compile for `linux/amd64`)
- **Push to Docker Hub** (central image registry)
- **Deploy to Vultr** (pull images, run via docker-compose)
- **Configuration management** (via .env files and rsync)

---

## Architecture: What's Running on Vultr

```
Vultr Server (1 instance)
├── Docker Engine
├── docker-compose.prod.yml (orchestrates all services)
├── Services:
│   ├── main-proxy (Nginx reverse proxy) → ports 80/443
│   ├── assist-frontend (React website) → port 3000
│   ├── assist-api (Go backend) → port 8080
│   ├── smartminds-admin (React admin panel) → port 3001
│   ├── smartminds-parent (React parent portal) → port 3002
│   ├── smartminds-marketing (Next.js public site) → port 3003
│   ├── mongodb (semantic links + medical content)
│   ├── skills-mysql (operational relational DB)
│   ├── weaviate (vector DB on port 8080 internal)
│   └── weaviate-chisg (master knowledge graph on port 8081 internal)
├── Volumes:
│   ├── mongodb-data
│   ├── mysql-data
│   ├── weaviate-data
│   └── weaviate-chisg-data
└── Networks:
    └── smart-minds-network (all containers interconnected)
```

---

## Pre-Migration Checklist

### 1. Vultr Account & Server Setup

- [ ] Create Vultr account (or use existing)
- [ ] Create API token for automated operations
- [ ] Create 2GB+ RAM instance (Ubuntu 22.04 LTS recommended)
- [ ] Note server IP (example: `192.248.151.185`)
- [ ] Create firewall rules: Allow SSH (22), HTTP (80), HTTPS (443)

### 2. Local Mac Preparation

- [ ] Ensure all code is committed to Git
- [ ] Create `.env.production` file with all secrets (do NOT commit)
- [ ] Build all Docker images locally
- [ ] Push images to Docker Hub
- [ ] Test docker-compose.prod.yml locally (if possible)

### 3. Domain & DNS

- [ ] Register or transfer domain (e.g., `smartminds.education`)
- [ ] Create DNS A records pointing to Vultr IP
- [ ] Set up SSL certificate (Let's Encrypt via Certbot)

---

## Step 1: Server Initialization (One-Time Setup)

Run these commands on the Vultr server via SSH:

```bash
# SSH into server
ssh root@<VULTR_IP>

# Update system packages
apt update && apt upgrade -y

# Install Docker
apt install -y docker.io docker-compose

# Add current user to docker group (optional, for non-root access)
usermod -aG docker root

# Create deployment directory
mkdir -p /opt/smart-minds
mkdir -p /opt/smart-minds/main-proxy
mkdir -p /opt/smart-minds/volumes/{mongodb-data,mysql-data,weaviate-data,weaviate-chisg-data}

# Create log directory
mkdir -p /var/log/smart-minds

# Set permissions
chmod 755 /opt/smart-minds

# Create a non-root deploy user (optional but recommended)
useradd -m -s /bin/bash deploy
usermod -aG docker deploy

# Enable Docker to start on boot
systemctl enable docker
systemctl start docker

# Verify Docker is running
docker --version
docker-compose --version
```

---

## Step 2: Build & Push Docker Images (On Mac)

### 2a. Create Docker Hub Account & Login

```bash
# On Mac
docker login

# Enter Docker Hub credentials when prompted
```

### 2b. Build All Production Images

Run these commands from `/Users/michaelstewart/Coding/assist/`:

```bash
# Main Proxy (Nginx)
docker buildx build \
  --platform linux/amd64 \
  -t <DOCKER_HUB_USERNAME>/smart-minds-proxy:latest \
  -f main-proxy/Dockerfile \
  --push .

# Assist Frontend
docker buildx build \
  --platform linux/amd64 \
  -t <DOCKER_HUB_USERNAME>/assist-frontend:latest \
  -f frontend/Dockerfile \
  --push ./frontend

# Assist API (Go backend)
docker buildx build \
  --platform linux/amd64 \
  -t <DOCKER_HUB_USERNAME>/assist-api:latest \
  -f esp-organizer/Dockerfile \
  --push ./esp-organizer

# Admin Panel
docker buildx build \
  --platform linux/amd64 \
  -t <DOCKER_HUB_USERNAME>/smartminds-admin:latest \
  -f websites/smartminds-admin/Dockerfile \
  --push ./websites/smartminds-admin

# Parent Portal
docker buildx build \
  --platform linux/amd64 \
  -t <DOCKER_HUB_USERNAME>/smartminds-parent:latest \
  -f websites/smartminds-parent/Dockerfile \
  --push ./websites/smartminds-parent

# Marketing Website (Next.js)
docker buildx build \
  --platform linux/amd64 \
  -t <DOCKER_HUB_USERNAME>/smartminds-marketing:latest \
  -f websites/smartminds-marketing/Dockerfile \
  --push ./websites/smartminds-marketing
```

**Note**: Replace `<DOCKER_HUB_USERNAME>` with your actual Docker Hub username.

### 2c: Verify All Images Are Pushed

```bash
# View Docker Hub images
# Visit: https://hub.docker.com/repositories
# Or: docker search <your-username>/smart-minds
```

---

## Step 3: Sync Configuration Files to Vultr

Run these commands from Mac (in `/Users/michaelstewart/Coding/assist/`):

```bash
# Variables
REMOTE_HOST="root@<VULTR_IP>"
REMOTE_DIR="/opt/smart-minds"
DOCKER_HUB_USER="<your-docker-hub-username>"

# Step 1: Upload docker-compose.prod.yml
rsync -avz --delete docker-compose.prod.yml \
  "${REMOTE_HOST}:${REMOTE_DIR}/docker-compose.yml"

# Step 2: Upload .env.production as .env on server
# (Make sure .env.production is in .gitignore!)
rsync -avz env.production \
  "${REMOTE_HOST}:${REMOTE_DIR}/.env"

# Step 3: Upload Nginx config
rsync -avz main-proxy/nginx.conf \
  "${REMOTE_HOST}:${REMOTE_DIR}/main-proxy/nginx.conf"

# Step 4: Upload database init scripts
rsync -avz skills-map-platform/init.sql \
  "${REMOTE_HOST}:${REMOTE_DIR}/skills-init.sql"

# Step 5: Create a .env file with Docker Hub credentials (optional, for private registries)
# ssh ${REMOTE_HOST} "cat > ${REMOTE_DIR}/.dockercfg << 'EOF'
# {\"auths\": {\"https://index.docker.io/v1/\": {\"auth\": \"<base64_encoded_credentials>\"}}}
# EOF"

echo "Configuration files synced to ${REMOTE_HOST}:${REMOTE_DIR}"
```

---

## Step 4: Update docker-compose.prod.yml

Before deploying, ensure `docker-compose.prod.yml` references correct image names:

```yaml
version: '3.8'

services:
  main-proxy:
    image: <DOCKER_HUB_USERNAME>/smart-minds-proxy:latest
    container_name: main-proxy
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./main-proxy/nginx.conf:/etc/nginx/nginx.conf:ro
      - /etc/letsencrypt:/etc/letsencrypt:ro  # SSL certificates
    networks:
      - smart-minds-network
    depends_on:
      - assist-frontend
      - assist-api
      - smartminds-admin
      - smartminds-parent

  assist-frontend:
    image: <DOCKER_HUB_USERNAME>/assist-frontend:latest
    container_name: assist-frontend
    restart: always
    networks:
      - smart-minds-network

  assist-api:
    image: <DOCKER_HUB_USERNAME>/assist-api:latest
    container_name: assist-api
    restart: always
    env_file: .env
    ports:
      - "8080:8080"
    depends_on:
      - mongodb
      - skills-mysql
      - weaviate
      - weaviate-chisg
    networks:
      - smart-minds-network

  smartminds-admin:
    image: <DOCKER_HUB_USERNAME>/smartminds-admin:latest
    container_name: smartminds-admin
    restart: always
    networks:
      - smart-minds-network

  smartminds-parent:
    image: <DOCKER_HUB_USERNAME>/smartminds-parent:latest
    container_name: smartminds-parent
    restart: always
    networks:
      - smart-minds-network

  smartminds-marketing:
    image: <DOCKER_HUB_USERNAME>/smartminds-marketing:latest
    container_name: smartminds-marketing
    restart: always
    networks:
      - smart-minds-network

  # Databases
  mongodb:
    image: mongo:6.0
    container_name: mongodb
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: ${MONGODB_PASSWORD}
    volumes:
      - mongodb-data:/data/db
    networks:
      - smart-minds-network
    ports:
      - "27017:27017"  # Internal only

  skills-mysql:
    image: mysql:8.0
    container_name: skills-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: skillsmap
    volumes:
      - mysql-data:/var/lib/mysql
      - ./skills-init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    networks:
      - smart-minds-network
    ports:
      - "3306:3306"  # Internal only

  weaviate:
    image: semitechnologies/weaviate:latest
    container_name: weaviate
    restart: always
    environment:
      QUERY_DEFAULTS_LIMIT: 20
      AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED: 'true'
      PERSISTENCE_DATA_PATH: '/var/lib/weaviate'
    volumes:
      - weaviate-data:/var/lib/weaviate
    networks:
      - smart-minds-network
    ports:
      - "8080:8080"  # Internal only

  weaviate-chisg:
    image: semitechnologies/weaviate:latest
    container_name: weaviate-chisg
    restart: always
    environment:
      QUERY_DEFAULTS_LIMIT: 20
      AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED: 'true'
      PERSISTENCE_DATA_PATH: '/var/lib/weaviate'
    volumes:
      - weaviate-chisg-data:/var/lib/weaviate
    networks:
      - smart-minds-network
    ports:
      - "8081:8080"  # Internal only

networks:
  smart-minds-network:
    driver: bridge

volumes:
  mongodb-data:
  mysql-data:
  weaviate-data:
  weaviate-chisg-data:
```

---

## Step 5: Deploy to Vultr (Full Deployment)

SSH into the Vultr server and run:

```bash
# SSH into server
ssh root@<VULTR_IP>

# Navigate to deployment directory
cd /opt/smart-minds

# Pull all images from Docker Hub
docker-compose pull

# Start all services in the background
docker-compose up -d

# Verify all services are running
docker-compose ps

# Check logs (if there are issues)
docker-compose logs -f main-proxy
docker-compose logs -f assist-api
docker-compose logs -f mongodb
```

---

## Step 6: SSL Setup with Let's Encrypt

```bash
# SSH into server
ssh root@<VULTR_IP>

# Install Certbot
apt install -y certbot python3-certbot-nginx

# Generate SSL certificate for your domain
certbot certonly --standalone \
  -d smartminds.education \
  -d www.smartminds.education \
  -d api.smartminds.education \
  -d app.smartminds.education \
  -d researcher.smartminds.education

# Verify certificates were created
ls -la /etc/letsencrypt/live/smartminds.education/

# Update Nginx config to use SSL (see main-proxy/nginx.conf)
# Then reload Nginx:
docker-compose restart main-proxy

# Set up auto-renewal
certbot renew --dry-run  # Test renewal
systemctl enable certbot.timer
systemctl start certbot.timer
```

---

## Step 7: Database Initialization

```bash
# SSH into server
ssh root@<VULTR_IP>

cd /opt/smart-minds

# Initialize MySQL database
docker-compose exec skills-mysql mysql -u root -p${MYSQL_ROOT_PASSWORD} < skills-init.sql

# Verify MySQL is ready
docker-compose exec skills-mysql mysql -u root -p${MYSQL_ROOT_PASSWORD} -e "SHOW DATABASES;"

# Initialize MongoDB
docker-compose exec mongodb mongosh -u admin -p${MONGODB_PASSWORD} --authenticationDatabase admin << 'EOF'
use esp_organizer;
db.semantic_links.createIndex({ "source_skill_id": 1, "target_skill_id": 1 });
db.semantic_links.createIndex({ "relationship_type": 1 });
EOF

# Check Weaviate is responding
curl -s http://localhost:8080/v1/meta | jq .

# Check Weaviate-CHISG is responding
curl -s http://localhost:8081/v1/meta | jq .
```

---

## Step 8: Backup Strategy

### Automated Daily Backups

Create a cron job to back up databases to Vultr Object Storage (S3-compatible):

```bash
# SSH into server
ssh root@<VULTR_IP>

# Install s3cmd (S3 client)
apt install -y s3cmd

# Configure s3cmd
s3cmd --configure
# Enter Vultr Object Storage credentials (see Vultr dashboard)

# Create backup script
cat > /opt/smart-minds/backup.sh << 'EOF'
#!/bin/bash
set -e

BACKUP_DIR="/tmp/smart-minds-backups"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BUCKET="smart-minds-backups"

mkdir -p $BACKUP_DIR

# Backup MySQL
docker-compose exec -T skills-mysql mysqldump \
  -u root -p${MYSQL_ROOT_PASSWORD} \
  --all-databases \
  > $BACKUP_DIR/mysql-$TIMESTAMP.sql

# Backup MongoDB
docker-compose exec -T mongodb mongodump \
  -u admin -p${MONGODB_PASSWORD} \
  --authenticationDatabase admin \
  --out $BACKUP_DIR/mongodb-$TIMESTAMP

# Backup Weaviate-CHISG data
tar -czf $BACKUP_DIR/weaviate-chisg-$TIMESTAMP.tar.gz \
  /opt/smart-minds/volumes/weaviate-chisg-data

# Upload to Vultr Object Storage
s3cmd sync $BACKUP_DIR/ s3://$BUCKET/

# Clean local backups older than 7 days
find $BACKUP_DIR -mtime +7 -delete

echo "Backup completed at $(date)"
EOF

chmod +x /opt/smart-minds/backup.sh

# Schedule daily backup at 2 AM
crontab -e
# Add line: 0 2 * * * /opt/smart-minds/backup.sh >> /var/log/smart-minds/backup.log 2>&1
```

---

## Step 9: Monitoring & Health Checks

### Enable Container Health Checks

In `docker-compose.prod.yml`, add healthchecks:

```yaml
assist-api:
  healthcheck:
    test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    start_period: 40s

mongodb:
  healthcheck:
    test: ["CMD", "mongosh", "-u", "admin", "-p", "${MONGODB_PASSWORD}", "--authenticationDatabase", "admin", "--eval", "db.adminCommand('ping')"]
    interval: 30s
    timeout: 10s
    retries: 3

skills-mysql:
  healthcheck:
    test: ["CMD", "mysqladmin", "ping", "-u", "root", "-p${MYSQL_ROOT_PASSWORD}"]
    interval: 30s
    timeout: 10s
    retries: 3
```

### Monitor Logs

```bash
# View logs for all services
docker-compose logs -f

# View logs for specific service
docker-compose logs -f assist-api

# View error logs only
docker-compose logs -f | grep ERROR

# Save logs to file
docker-compose logs > /var/log/smart-minds/docker.log
```

### Set Up Alerts (Optional)

Use a service like Sentry, DataDog, or New Relic:

```bash
# Example: Sentry (error tracking)
# Add to .env:
SENTRY_DSN=https://your-sentry-key@sentry.io/project-id

# Add to assist-api container:
environment:
  SENTRY_DSN: ${SENTRY_DSN}
```

---

## Step 10: DNS Finalization

After verifying everything is working on the Vultr IP directly, update DNS:

```bash
# Current (testing):
# Access via IP: http://<VULTR_IP>

# Final (production):
# Update DNS A records to point domain to Vultr IP
smartminds.education       A  <VULTR_IP>
api.smartminds.education   CNAME smartminds.education
app.smartminds.education   CNAME smartminds.education
www.smartminds.education   CNAME smartminds.education
researcher.smartminds.education CNAME smartminds.education

# Wait for DNS propagation (10 minutes to 24 hours)
nslookup smartminds.education
```

---

## Ongoing Maintenance

### Weekly Tasks

- [ ] Check disk space: `df -h`
- [ ] Check container logs for errors
- [ ] Verify backups completed successfully

### Monthly Tasks

- [ ] Review database sizes: `docker-compose exec skills-mysql du -sh /var/lib/mysql`
- [ ] Update Docker images: `docker-compose pull && docker-compose up -d`
- [ ] Test backup restoration procedure

### Quarterly Tasks

- [ ] Review and update SSL certificates
- [ ] Audit user access logs
- [ ] Test disaster recovery plan

---

## Troubleshooting

### Containers Won't Start

```bash
# Check container logs
docker-compose logs <container-name>

# Rebuild and restart
docker-compose down
docker-compose pull
docker-compose up -d
```

### Database Connection Issues

```bash
# Test connectivity from assist-api to mongodb
docker-compose exec assist-api nc -zv mongodb 27017

# Test connectivity from assist-api to weaviate-chisg
docker-compose exec assist-api curl -s http://weaviate-chisg:8080/v1/meta

# Check .env variables are loaded
docker-compose config | grep MONGODB_URI
```

### Out of Disk Space

```bash
# Check disk usage
docker system df

# Clean up unused images/volumes
docker system prune -a
```

### SSL Certificate Issues

```bash
# Check certificate expiration
certbot certificates

# Renew manually
certbot renew --force-renewal
```

---

## Quick Reference: Useful Commands

```bash
# SSH into server
ssh root@<VULTR_IP>

# Navigate to deployment
cd /opt/smart-minds

# View running containers
docker-compose ps

# View all logs (realtime)
docker-compose logs -f

# View logs for specific service
docker-compose logs -f assist-api

# Execute command in container
docker-compose exec assist-api curl http://weaviate-chisg:8080/v1/meta

# Stop all services
docker-compose stop

# Start all services
docker-compose start

# Restart specific service
docker-compose restart assist-api

# Pull latest images and restart
docker-compose pull && docker-compose up -d

# View resource usage
docker stats

# Check container status
docker-compose ps

# Remove all containers (destructive!)
docker-compose down
```

---

## Security Checklist

- [ ] Enable UFW firewall: `ufw enable`
- [ ] Allow only SSH, HTTP, HTTPS: `ufw allow 22/tcp` `ufw allow 80/tcp` `ufw allow 443/tcp`
- [ ] Disable root SSH login
- [ ] Set up SSH key authentication (no passwords)
- [ ] Keep Docker and OS updated regularly
- [ ] Use strong database passwords (stored in .env, not in code)
- [ ] Enable HTTPS/SSL for all domains
- [ ] Implement rate limiting in Nginx
- [ ] Set up DDoS protection (Vultr offers this)
- [ ] Regular backup tests

---

## Next Steps

1. **Create Vultr account** and configure server
2. **Prepare images** (build on Mac, push to Docker Hub)
3. **Execute deployment** (sync config, start containers)
4. **Test services** (verify all endpoints work)
5. **Configure DNS** (point domain to Vultr)
6. **Set up monitoring** (logging, health checks, backups)
7. **Document runbooks** (what to do when things break)

---

*This guide assumes familiarity with Docker, Linux, and basic DevOps practices. For questions, refer to Vultr documentation or Docker documentation.*
