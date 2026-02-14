#!/bin/bash
#
# IONOS → Vultr Migration Script
# Migrates espthinking.co.uk from WordPress on IONOS to React portfolio on Vultr
#
# Prerequisites:
# - SSH access to Vultr (vultr host in ~/.ssh/config)
# - DNS access to IONOS
# - WordPress backup completed via UpdraftPlus
#
# Usage: ./migrate-to-vultr.sh [step]
#   Steps: backup, deploy, dns-check, ssl

set -euo pipefail

VULTR_HOST="vultr"
VULTR_IP="192.248.151.185"
DOMAIN="espthinking.co.uk"
LOCAL_PROJECT="/Users/michaelstewart/Coding/assist"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $1"; }
success() { echo -e "${GREEN}✓${NC} $1"; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1"; }

# Step 1: Backup WordPress (manual step reminder)
step_backup() {
    log "Step 1: WordPress Backup"
    echo ""
    echo "Before proceeding, ensure you have:"
    echo "  1. Logged into WordPress admin at https://www.espthinking.co.uk/wp-admin"
    echo "  2. Gone to Settings → UpdraftPlus Backups"
    echo "  3. Clicked 'Backup Now' with Database + Files"
    echo "  4. Downloaded the backup to your local machine"
    echo ""
    read -p "Have you completed the WordPress backup? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        success "WordPress backup confirmed"
    else
        warn "Please complete the backup before continuing"
        exit 1
    fi
}

# Step 2: Deploy to Vultr
step_deploy() {
    log "Step 2: Deploy to Vultr"
    
    # Build frontend
    log "Building frontend..."
    cd "$LOCAL_PROJECT/frontend"
    npm run build
    success "Frontend built"
    
    # Copy nginx config
    log "Updating nginx configuration..."
    scp "$LOCAL_PROJECT/main-proxy/nginx-vultr-portfolio.conf" \
        "$VULTR_HOST:/root/assist/main-proxy/default.conf"
    success "Nginx config uploaded"
    
    # Rebuild and restart on Vultr
    log "Rebuilding containers on Vultr..."
    ssh "$VULTR_HOST" "cd /root/assist && docker-compose build assist-frontend && docker-compose up -d assist-frontend main-proxy"
    success "Containers rebuilt"
    
    # Test
    log "Testing deployment..."
    sleep 5
    if curl -sS "http://$VULTR_IP" | grep -q "Educational Data Systems"; then
        success "Portfolio is live at http://$VULTR_IP"
    else
        warn "Portfolio may not be serving correctly - check manually"
    fi
}

# Step 3: DNS Check
step_dns() {
    log "Step 3: DNS Configuration"
    echo ""
    echo "Current DNS for $DOMAIN:"
    dig +short "$DOMAIN"
    echo ""
    echo "To migrate, update DNS at IONOS:"
    echo "  1. Log into IONOS control panel"
    echo "  2. Go to Domains & SSL → $DOMAIN → DNS"
    echo "  3. Update A record for @ to: $VULTR_IP"
    echo "  4. Update A record for www to: $VULTR_IP"
    echo "  5. TTL: 300 (5 minutes) for quick propagation"
    echo ""
    echo "After updating, wait for propagation (5-60 minutes)"
    echo "Check with: dig $DOMAIN +short"
    echo ""
    echo "Expected result: $VULTR_IP"
}

# Step 4: SSL Setup
step_ssl() {
    log "Step 4: SSL Setup with Let's Encrypt"
    
    # Check DNS first
    current_ip=$(dig +short "$DOMAIN" | head -1)
    if [ "$current_ip" != "$VULTR_IP" ]; then
        error "DNS not yet pointing to Vultr ($current_ip != $VULTR_IP)"
        echo "Wait for DNS propagation before setting up SSL"
        exit 1
    fi
    
    success "DNS verified: $DOMAIN → $VULTR_IP"
    
    # Install certbot and get certificate
    log "Setting up SSL on Vultr..."
    ssh "$VULTR_HOST" << 'REMOTESCRIPT'
        # Install certbot if not present
        if ! command -v certbot &> /dev/null; then
            apt-get update && apt-get install -y certbot python3-certbot-nginx
        fi
        
        # Get certificate (standalone mode - stop nginx briefly)
        docker stop main-proxy
        certbot certonly --standalone -d espthinking.co.uk -d www.espthinking.co.uk --non-interactive --agree-tos --email admin@espthinking.co.uk
        docker start main-proxy
        
        echo "SSL certificate obtained"
REMOTESCRIPT
    
    success "SSL certificate obtained"
    echo ""
    echo "Next steps:"
    echo "  1. Update nginx config to use SSL (see nginx-vultr-portfolio.conf comments)"
    echo "  2. Mount /etc/letsencrypt into main-proxy container"
    echo "  3. Restart main-proxy"
}

# Main
case "${1:-all}" in
    backup)
        step_backup
        ;;
    deploy)
        step_deploy
        ;;
    dns|dns-check)
        step_dns
        ;;
    ssl)
        step_ssl
        ;;
    all)
        echo "=========================================="
        echo "IONOS → Vultr Migration"
        echo "=========================================="
        echo ""
        echo "This script will guide you through migrating"
        echo "$DOMAIN from WordPress on IONOS to the new"
        echo "React portfolio on Vultr."
        echo ""
        echo "Steps:"
        echo "  1. backup  - Confirm WordPress backup"
        echo "  2. deploy  - Build and deploy to Vultr"
        echo "  3. dns     - Update DNS at IONOS"
        echo "  4. ssl     - Set up HTTPS with Let's Encrypt"
        echo ""
        echo "Run each step individually:"
        echo "  ./migrate-to-vultr.sh backup"
        echo "  ./migrate-to-vultr.sh deploy"
        echo "  ./migrate-to-vultr.sh dns"
        echo "  ./migrate-to-vultr.sh ssl"
        ;;
    *)
        echo "Unknown step: $1"
        echo "Valid steps: backup, deploy, dns, ssl, all"
        exit 1
        ;;
esac
