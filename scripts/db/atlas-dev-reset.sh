#!/usr/bin/env bash
# Atlas migration reset script
# Regenerates migrations from HCL schemas and restores originals if no content changes detected
#
# Usage: ./scripts/db/atlas-dev-reset.sh

set -euo pipefail

# Configuration
ATLAS_ENV_WATCH="${ATLAS_ENV_WATCH:-local_watch}"
ATLAS_ENV_KEYGEN="${ATLAS_ENV_KEYGEN:-local_keygen}"
ATLAS_ENV_SIGN="${ATLAS_ENV_SIGN:-local_sign}"
SCHEMAS=("watch" "keygen" "sign")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
ATLAS_DIR="${PROJECT_ROOT}/tools/atlas"
MIGRATIONS_DIR="${ATLAS_DIR}/migrations"

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Cleanup function for trap
cleanup() {
    if [[ -n "${BACKUP_DIR:-}" ]] && [[ -d "${BACKUP_DIR}" ]]; then
        rm -rf "${BACKUP_DIR}"
    fi
}

trap cleanup EXIT

log_warning "Regenerating all migrations from HCL schemas..."
log_info "If no content changes are detected, original files will be restored."
echo ""

# Create backup directory
BACKUP_DIR=$(mktemp -d)
HAS_CHANGES=0

# Backup existing migrations
for schema in "${SCHEMAS[@]}"; do
    mkdir -p "${BACKUP_DIR}/${schema}"
    if ls "${MIGRATIONS_DIR}/${schema}"/*.sql 1>/dev/null 2>&1; then
        cp "${MIGRATIONS_DIR}/${schema}"/*.sql "${BACKUP_DIR}/${schema}/"
    fi
    if [[ -f "${MIGRATIONS_DIR}/${schema}/atlas.sum" ]]; then
        cp "${MIGRATIONS_DIR}/${schema}/atlas.sum" "${BACKUP_DIR}/${schema}/"
    fi
done

# Remove existing migrations
for schema in "${SCHEMAS[@]}"; do
    rm -f "${MIGRATIONS_DIR}/${schema}"/*.sql "${MIGRATIONS_DIR}/${schema}/atlas.sum"
done

# Generate new migrations
echo "=== Generating watch migration ==="
(cd "${ATLAS_DIR}" && atlas migrate diff initial_schema --config file://atlas.hcl --env "${ATLAS_ENV_WATCH}") || exit 1

echo "=== Generating keygen migration ==="
(cd "${ATLAS_DIR}" && atlas migrate diff initial_schema --config file://atlas.hcl --env "${ATLAS_ENV_KEYGEN}") || exit 1

echo "=== Generating sign migration ==="
(cd "${ATLAS_DIR}" && atlas migrate diff initial_schema --config file://atlas.hcl --env "${ATLAS_ENV_SIGN}") || exit 1

# Compare generated migrations with originals
echo "=== Comparing generated migrations with originals ==="
for schema in "${SCHEMAS[@]}"; do
    OLD_SQL=("${BACKUP_DIR}/${schema}"/*.sql)
    NEW_SQL=("${MIGRATIONS_DIR}/${schema}"/*.sql)
    
    # Check if both old and new SQL files exist
    if [[ -f "${OLD_SQL[0]}" ]] && [[ -f "${NEW_SQL[0]}" ]]; then
        if ! diff -q "${OLD_SQL[0]}" "${NEW_SQL[0]}" >/dev/null 2>&1; then
            echo "  ${schema}: Content changed"
            HAS_CHANGES=1
        else
            echo "  ${schema}: No content changes"
        fi
    else
        echo "  ${schema}: New migration created"
        HAS_CHANGES=1
    fi
done

# Restore or keep based on changes
if [[ "${HAS_CHANGES}" -eq 0 ]]; then
    echo ""
    log_info "No actual content changes detected."
    echo "   Restoring original migration files..."
    
    for schema in "${SCHEMAS[@]}"; do
        rm -f "${MIGRATIONS_DIR}/${schema}"/*.sql "${MIGRATIONS_DIR}/${schema}/atlas.sum"
        if ls "${BACKUP_DIR}/${schema}"/*.sql 1>/dev/null 2>&1; then
            cp "${BACKUP_DIR}/${schema}"/*.sql "${MIGRATIONS_DIR}/${schema}/"
        fi
        if [[ -f "${BACKUP_DIR}/${schema}/atlas.sum" ]]; then
            cp "${BACKUP_DIR}/${schema}/atlas.sum" "${MIGRATIONS_DIR}/${schema}/"
        fi
    done
    
    log_success "Original files restored. No changes made."
else
    echo ""
    log_success "Migrations regenerated with content changes."
    echo "  Run 'make reset-docker-db' to apply."
fi

