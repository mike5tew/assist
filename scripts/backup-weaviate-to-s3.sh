#!/bin/bash
#
# Weaviate Backup to S3
# Exports critical Weaviate classes to JSON and uploads to S3 with versioning
#
# Usage: ./backup-weaviate-to-s3.sh [weaviate_host] [s3_bucket]
#
# Requirements:
#   - AWS CLI configured with credentials
#   - jq installed
#   - curl installed
#
# Recommended: Run daily via cron
#   0 3 * * * /path/to/backup-weaviate-to-s3.sh >> /var/log/weaviate-backup.log 2>&1

set -euo pipefail

# Configuration
WEAVIATE_HOST="${1:-http://localhost:8081}"
S3_BUCKET="${2:-esp-weaviate-backups}"
BACKUP_DIR="/tmp/weaviate-backup-$(date +%Y%m%d-%H%M%S)"
DATE_STAMP=$(date +%Y-%m-%d)
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

# Classes to backup (critical data)
CLASSES=(
    "CHISGElement"
    "SkillLink"
    "Documentation"
    "SemanticLinks"
    "MedicalExcerpt"
    "CourseSkillSuggestions"
)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

log_success() {
    log "${GREEN}✓ $1${NC}"
}

log_warning() {
    log "${YELLOW}⚠ $1${NC}"
}

log_error() {
    log "${RED}✗ $1${NC}"
}

# Check dependencies
check_dependencies() {
    local missing=()
    command -v curl >/dev/null 2>&1 || missing+=("curl")
    command -v jq >/dev/null 2>&1 || missing+=("jq")
    command -v aws >/dev/null 2>&1 || missing+=("aws")
    
    if [ ${#missing[@]} -ne 0 ]; then
        log_error "Missing dependencies: ${missing[*]}"
        exit 1
    fi
}

# Check Weaviate connectivity
check_weaviate() {
    log "Checking Weaviate connectivity at ${WEAVIATE_HOST}..."
    if ! curl -sS "${WEAVIATE_HOST}/v1/.well-known/ready" >/dev/null 2>&1; then
        log_error "Cannot connect to Weaviate at ${WEAVIATE_HOST}"
        exit 1
    fi
    log_success "Weaviate is reachable"
}

# Get count of objects in a class
get_class_count() {
    local class_name=$1
    local count
    count=$(curl -sS -X POST "${WEAVIATE_HOST}/v1/graphql" \
        -H 'Content-Type: application/json' \
        -d "{\"query\":\"{ Aggregate { ${class_name} { meta { count } } } }\"}" \
        | jq -r ".data.Aggregate.${class_name}[0].meta.count // 0")
    echo "$count"
}

# Export a class to JSON (paginated)
export_class() {
    local class_name=$1
    local output_file="${BACKUP_DIR}/${class_name}.ndjson"
    local batch_size=100
    local offset=0
    local total_exported=0
    
    log "Exporting ${class_name}..."
    
    # Get schema for this class to know which properties to export
    local properties
    properties=$(curl -sS "${WEAVIATE_HOST}/v1/schema/${class_name}" \
        | jq -r '[.properties[].name] | join(" ")')
    
    if [ -z "$properties" ]; then
        log_warning "No properties found for ${class_name}, skipping"
        return 0
    fi
    
    # Clear output file
    > "$output_file"
    
    while true; do
        local batch
        batch=$(curl -sS -X POST "${WEAVIATE_HOST}/v1/graphql" \
            -H 'Content-Type: application/json' \
            -d "{\"query\":\"{ Get { ${class_name}(limit: ${batch_size}, offset: ${offset}) { _additional { id } ${properties} } } }\"}" \
            | jq -c ".data.Get.${class_name}[]?" 2>/dev/null)
        
        if [ -z "$batch" ]; then
            break
        fi
        
        # Count items in this batch
        local batch_count
        batch_count=$(echo "$batch" | wc -l | tr -d ' ')
        
        if [ "$batch_count" -eq 0 ]; then
            break
        fi
        
        # Append to output file
        echo "$batch" >> "$output_file"
        
        total_exported=$((total_exported + batch_count))
        offset=$((offset + batch_size))
        
        # Safety limit
        if [ $offset -gt 50000 ]; then
            log_warning "Safety limit reached for ${class_name}"
            break
        fi
    done
    
    if [ $total_exported -gt 0 ]; then
        log_success "Exported ${total_exported} objects from ${class_name}"
    else
        log_warning "No objects found in ${class_name}"
        rm -f "$output_file"
    fi
    
    echo "$total_exported"
}

# Create backup manifest
create_manifest() {
    local manifest_file="${BACKUP_DIR}/manifest.json"
    
    cat > "$manifest_file" << EOF
{
    "backup_timestamp": "${TIMESTAMP}",
    "backup_date": "${DATE_STAMP}",
    "weaviate_host": "${WEAVIATE_HOST}",
    "classes_exported": $(find "$BACKUP_DIR" -name "*.ndjson" -exec basename {} .ndjson \; | jq -R -s -c 'split("\n") | map(select(length > 0))'),
    "file_sizes": {
$(for f in "$BACKUP_DIR"/*.ndjson; do
    if [ -f "$f" ]; then
        echo "        \"$(basename "$f")\": $(stat -f%z "$f" 2>/dev/null || stat -c%s "$f" 2>/dev/null || echo 0),"
    fi
done | sed '$ s/,$//')
    }
}
EOF
    
    log_success "Created backup manifest"
}

# Compress backup
compress_backup() {
    local archive_name="weaviate-backup-${TIMESTAMP}.tar.gz"
    local archive_path="/tmp/${archive_name}"
    
    log "Compressing backup..."
    tar -czf "$archive_path" -C "$(dirname "$BACKUP_DIR")" "$(basename "$BACKUP_DIR")"
    
    log_success "Created archive: ${archive_path} ($(du -h "$archive_path" | cut -f1))"
    echo "$archive_path"
}

# Upload to S3
upload_to_s3() {
    local archive_path=$1
    local s3_key="backups/${DATE_STAMP}/$(basename "$archive_path")"
    
    log "Uploading to S3: s3://${S3_BUCKET}/${s3_key}"
    
    if aws s3 cp "$archive_path" "s3://${S3_BUCKET}/${s3_key}" \
        --storage-class STANDARD_IA \
        --metadata "backup-date=${DATE_STAMP},source=${WEAVIATE_HOST}"; then
        log_success "Uploaded to S3"
        
        # Also upload to latest/
        aws s3 cp "$archive_path" "s3://${S3_BUCKET}/latest/weaviate-backup-latest.tar.gz" \
            --storage-class STANDARD_IA
        log_success "Updated latest backup pointer"
    else
        log_error "Failed to upload to S3"
        return 1
    fi
}

# Cleanup old local files
cleanup() {
    log "Cleaning up temporary files..."
    rm -rf "$BACKUP_DIR"
    rm -f "/tmp/weaviate-backup-*.tar.gz"
    log_success "Cleanup complete"
}

# Main execution
main() {
    log "=========================================="
    log "Weaviate Backup to S3"
    log "=========================================="
    log "Weaviate Host: ${WEAVIATE_HOST}"
    log "S3 Bucket: ${S3_BUCKET}"
    log "Timestamp: ${TIMESTAMP}"
    log ""
    
    check_dependencies
    check_weaviate
    
    # Create backup directory
    mkdir -p "$BACKUP_DIR"
    
    # Export each class
    local total_objects=0
    for class in "${CLASSES[@]}"; do
        local count
        count=$(get_class_count "$class")
        if [ "$count" -gt 0 ]; then
            exported=$(export_class "$class")
            total_objects=$((total_objects + exported))
        else
            log_warning "Skipping empty class: ${class}"
        fi
    done
    
    if [ $total_objects -eq 0 ]; then
        log_warning "No objects exported, nothing to backup"
        cleanup
        exit 0
    fi
    
    # Create manifest
    create_manifest
    
    # Compress
    archive_path=$(compress_backup)
    
    # Upload to S3
    upload_to_s3 "$archive_path"
    
    # Cleanup
    cleanup
    
    log ""
    log "=========================================="
    log_success "Backup complete! ${total_objects} total objects backed up"
    log "=========================================="
}

main "$@"
