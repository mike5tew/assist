# Deployment Playbook

**Purpose**: This document provides the exact commands and procedures for deploying the HumanOS ecosystem to the Vultr production server. It covers both initial full deployments and routine updates to individual services.

---

## Core Principles

1.  **Build on Mac, Run on Server**: We build Docker images on the local Mac (cross-compiling for `linux/amd64`) and push them to Docker Hub. The Vultr server **only** pulls images from Docker Hub; it never builds them.
2.  **Docker Hub is the Intermediary**: Docker Hub acts as the central registry for all production-ready images.
3.  **Configuration is Synced**: All configuration files (`docker-compose.prod.yml`, `.env`, `nginx.conf`) are managed in the `assist` monorepo and synced to the server using `rsync`.

---

## Scenario 1: Initial Full Deployment (Wipe & Replace)

Use this scenario when setting up the server for the first time or when a major architectural change requires a clean slate.

### Step 0: Prerequisites (Run on Mac)
Before you can deploy, you must create a production environment file.

1.  Copy the example environment file: `cp env.example env.production`
2.  Open `env.production` in a text editor.
3.  Fill in all the required secrets (passwords, keys, etc.). **This file must not be committed to Git.**

### Step 1: Build and Push All Production Images (Run on Mac)
Ensure every service has a `linux/amd64` image pushed to Docker Hub.

- **Main Proxy**: (from `/assist`)
  `docker buildx build --platform linux/amd64 -t mike5tew/main-proxy:latest -f main-proxy/Dockerfile --push .`

- **Assist Frontend**: (from `/assist/frontend`)
  `docker buildx build --platform linux/amd64 -t mike5tew/assist-frontend:latest --push .`

- **Assist API**: (from `/assist/esp-organizer`)
  `docker buildx build --platform linux/amd64 -t mike5tew/assist-api:latest --push .`

- **Skills Frontend**: (from `/assist/skills-map-platform/frontend`)
  `docker buildx build --platform linux/amd64 -t mike5tew/skills-frontend:latest --push .`

- **Skills API**: (from `/assist/skills-map-platform/api`)
  `docker buildx build --platform linux/amd64 -t mike5tew/skills-api:latest --push .`

### Step 2: Full Server Cleanup (Run on Mac)
This script connects to the server and removes everything to ensure a clean start. **This is destructive and will remove all old deployments.**

```bash
ssh root@192.248.151.185 << 'EOF'
  # Stop and remove all existing containers
  docker stop $(docker ps -aq) && docker rm $(docker ps -aq)
  
  # Prune unused networks and volumes
  docker network prune -f
  docker volume prune -f
  
  # Remove all old deployment directories and create the new one
  # The /opt directory is a standard top-level directory on Linux for optional software packages.
  # It is located at the root of the filesystem (like /home or /etc), not inside another folder.
  rm -rf /opt/esp
  rm -rf /opt/skills-map-platform
  rm -rf /opt/main-proxy
  # Create the root directory AND the necessary subdirectory for nginx
  mkdir -p /opt/esp/main-proxy
EOF
```

### Step 3: Upload All Configurations (Run on Mac)
This script syncs all necessary configuration files to the server's new `/opt/esp` directory.

```bash
# (Run from /Users/michaelstewart/Coding/assist/)
REMOTE_HOST="root@192.248.151.185"
REMOTE_DIR="/opt/esp"

ssh ${REMOTE_HOST} "mkdir -p ${REMOTE_DIR}/main-proxy"

rsync -avz docker-compose.prod.yml "${REMOTE_HOST}:${REMOTE_DIR}/docker-compose.yml"
rsync -avz env.production "${REMOTE_HOST}:${REMOTE_DIR}/.env"
rsync -avz main-proxy/nginx.conf "${REMOTE_HOST}:${REMOTE_DIR}/main-proxy/nginx.conf"
rsync -avz skills-map-platform/init.sql "${REMOTE_HOST}:${REMOTE_DIR}/init.sql"
```

### Step 4: Run the Entire Stack (Run on Vultr Server)
Connect to the server and start everything.

```bash
# ssh root@192.248.151.185
cd /opt/esp
docker compose pull
docker compose up -d
```

---

## Scenario 2: Updating a Single Service (Partial Update)

Use this common scenario when you've made code changes to one service (e.g., `assist-frontend`) and want to deploy only that update.

### Step 1: Build and Push the Updated Image (Run on Mac)
Re-build and push the image for the specific service you changed.

```bash
# Example: Updating the Assist Frontend
cd /Users/michaelstewart/Coding/assist/frontend/
docker buildx build --platform linux/amd64 -t mike5tew/assist-frontend:latest --push .
```

### Step 2: Deploy the Service on the Server (Run on Vultr Server)
Connect to the server, pull the new image, and restart only that service.

```bash
# ssh root@192.248.151.185
cd /opt/esp

# Pull the new image for the specific service
docker compose pull assist-frontend

# Restart the service. --no-deps prevents restarting dependencies.
docker compose up -d --no-deps assist-frontend
```
You can replace `assist-frontend` with any other service name from your `docker-compose.prod.yml` file (e.g., `assist-api`, `skills-frontend`).

---

## Scenario 3: Updating Configuration Files

Use this scenario when you've changed a configuration file like `docker-compose.prod.yml` or `main-proxy/nginx.conf`.

### Step 1: Sync the Changed File (Run on Mac)
Use `rsync` to upload only the file that has changed.

```bash
# Example: Updating the main Nginx proxy configuration
rsync -avz main-proxy/nginx.conf root@192.248.151.185:/opt/esp/main-proxy/nginx.conf
```

### Step 2: Apply the Configuration on the Server (Run on Vultr Server)
Connect to the server and run `docker compose up`. Docker Compose is smart enough to see which services are affected by the configuration change and will recreate only those.

```bash
# ssh root@192.248.151.185
cd /opt/esp
docker compose up -d
```
This command will automatically restart the `main-proxy` container because its configuration has changed, while leaving other, unaffected containers running.

---

## Debugging

If you encounter issues, such as a container failing to start, follow these steps to diagnose the problem.

### Step 1: Check Container Logs

Use the `docker logs` command to view the logs for the failed container. This will provide insight into what went wrong.

```bash
# Example: Checking logs for the skills-db container
docker logs skills-db
```

### Step 2: Common Issues

1.  **Database Initialization Errors**: If the logs indicate an error with the `init.sql` script, double-check the script for mistakes.
2.  **Missing Environment Variables**: Ensure that all required environment variables are set in the `.env` file on the server.

### Step 3: Restarting Containers

After fixing any issues, you may need to restart the affected containers.

```bash
# Example: Restarting the skills-db container
docker compose up -d --no-deps skills-db
```

Refer to the specific service's documentation for any additional troubleshooting steps.

---

## SSH Connection Issues

If you encounter SSH connection errors, such as `Connection reset by peer`, follow these steps to diagnose and fix the problem.

### Step 1: Wait and Retry

Often, this error is temporary. Wait 10-15 minutes and then try to connect again.

```bash
# Example: Trying to SSH again after waiting
ssh root@192.248.151.185
```

### Step 2: Check Server Status

Ensure the server is running and accessible. You can check this from the Vultr control panel.

### Step 3: Restart SSH Service

If you have access to the server console (e.g., through the Vultr control panel), you can try restarting the SSH service.

```bash
# Example: Restarting the SSH service
sudo systemctl restart sshd
```

### Step 4: Check Fail2Ban

If the server uses `fail2ban` or similar software, your IP might be temporarily banned. Check the `fail2ban` logs and unban your IP if necessary.

```bash
# Example: Checking fail2ban logs
sudo cat /var/log/fail2ban.log
```

### Step 5: Reboot the Server

As a last resort, you can try rebooting the server from the Vultr control panel. This can resolve issues with unresponsive services.

---

## Verification

After deployment, it's crucial to verify that everything is working as expected. Follow these steps to confirm a successful deployment.

### Step 1: Check Environment File

Ensure that the `.env` file exists and is not empty.

```bash
# 1. Check that the .env file exists and has content
ls -la
# You should see the .env file listed.
```

### Step 2: Check Docker Compose Version

Verify that the version of the `docker-compose` file is compatible.

```bash
# 2. Check the version of the docker-compose file
cat docker-compose.yml | grep "Version:"
# This should output "# Version: 1.4" or higher.
```

### Step 3: Review Environment Variables

Check the content of the `.env` file to ensure all secrets are correctly set.

```bash
# 3. Check the content of the .env file
cat .env
# This should show the secrets you entered in env.production
```

### Step 4: Validate Service Status

Confirm that all services are up and running.

```bash
# 4. Validate that all services are running
docker ps
# You should see all your containers listed and running.
```

### Step 5: Check Logs for Errors

Review the logs for any errors or warnings.

```bash
# 5. Check the logs for any errors
docker-compose logs
# Look for any ERROR or WARNING messages.
```

## Log rotation (Recommended)

Logs can grow over time and may fill disks. Install `logrotate` on your server and create a rotation config. An example file is provided at `scripts/logrotate.conf` in this repo — deploy it to `/etc/logrotate.d/esp-organizer` and adjust paths if your logs live elsewhere (for example `/opt/esp/logs/*.log` or `/var/log/esp/*.log`). Example rotation policy:

## Log level (Useful for runtime control)

The application respects a `LOG_LEVEL` environment variable which controls the verbosity of the structured logger. Supported values: `debug`, `info`, `warn`, `error` (default: `info`). Set `LOG_LEVEL=debug` on staging or local environments for verbose logs.

- Rotate daily
- Keep 14 rotations
- Compress rotated files
- Use `copytruncate` if the process cannot be restarted during rotation

Example deployment steps (on the server):

```bash
# Install logrotate (Debian/Ubuntu)
sudo apt-get update && sudo apt-get install -y logrotate

# Copy the example config and verify
sudo cp /opt/esp/scripts/logrotate.conf /etc/logrotate.d/esp-organizer
sudo logrotate --debug /etc/logrotate.d/esp-organizer
```

If all these checks pass, your deployment was successful. If you encounter any issues, refer to the debugging section of this playbook.
