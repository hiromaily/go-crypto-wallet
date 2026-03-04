#!/usr/bin/env bash

# XRP-Specific Common Utility Functions for E2E Operations
# This file provides shared utility functions specifically for XRP E2E scripts.
#
# Usage:
#   source "$(dirname "${BASH_SOURCE[0]}")/xrp_common.sh"
#
# Note: This file sources common.sh automatically, so you don't need to source both.
#
# Infrastructure:
#   - rippled v3.1.0 in standalone mode (via compose.xrp.yaml)
#   - WebSocket admin port: 6006 (ws://)
#   - No xrpl-grpc-server (wallet connects directly via WebSocket)
#
# Genesis Account (standalone mode only):
#   Address: rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh
#   Secret:  snoPBrXtMeMyMHUVTgbuqAfg1SUTb

# Script directory for relative paths
_XRP_COMMON_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the general common utilities
# shellcheck source=../common.sh
source "${_XRP_COMMON_DIR}/../common.sh"

###############################################################################
# XRP Default Configuration
###############################################################################

# Ensure GOPATH is set — CI runners may not have it as an environment variable
GOPATH="${GOPATH:-$(go env GOPATH)}"
export GOPATH

# Coin identifier
XRP_COIN="xrp"

# rippled WebSocket host/port (standalone mode)
# Port 6006 is the WebSocket admin port in rippled.cfg
XRP_WS_HOST="${XRP_WS_HOST:-127.0.0.1}"
XRP_WS_PORT="${XRP_WS_PORT:-6006}"

# MySQL credentials (can be overridden via environment variables)
XRP_MYSQL_ROOT_PASSWORD="${XRP_MYSQL_ROOT_PASSWORD:-${MYSQL_ROOT_PASSWORD:-root}}"

# rippled genesis account for standalone funding
# These credentials are well-known for standalone mode — never use in production.
XRP_GENESIS_ACCOUNT="${XRP_GENESIS_ACCOUNT:-rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh}"
XRP_GENESIS_SECRET="${XRP_GENESIS_SECRET:-snoPBrXtMeMyMHUVTgbuqAfg1SUTb}"

###############################################################################
# Database Configuration
###############################################################################

# Database type: sqlite (default) or mysql
DB_TYPE="${DB_TYPE:-sqlite}"

# E2E Pattern identifier (set by each E2E script)
E2E_PATTERN="${E2E_PATTERN:-}"

# Get database file path with pattern suffix for a given wallet type
xrp_get_db_path() {
	local wallet_type="${1:-watch}"
	local db_suffix=""
	if [ -n "${E2E_PATTERN}" ]; then
		db_suffix="-e2e-${E2E_PATTERN}"
	else
		db_suffix="-e2e"
	fi
	echo "./data/sqlite/xrp/${wallet_type}${db_suffix}.db"
}

# Separate DB paths for each wallet type
SQLITE_WATCH_DB_PATH=""
SQLITE_KEYGEN_DB_PATH=""

# Initialize database configuration — sets SQLITE_*_DB_PATH and exports env vars
xrp_init_database() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		export WALLET_DATABASE_TYPE="sqlite"

		SQLITE_WATCH_DB_PATH="$(xrp_get_db_path "watch")"
		SQLITE_KEYGEN_DB_PATH="$(xrp_get_db_path "keygen")"

		# Create directories if they don't exist
		mkdir -p "$(dirname "${SQLITE_WATCH_DB_PATH}")"
		mkdir -p "$(dirname "${SQLITE_KEYGEN_DB_PATH}")"

		log_info "Using SQLite databases:"
		log_info "  watch:  ${SQLITE_WATCH_DB_PATH}"
		log_info "  keygen: ${SQLITE_KEYGEN_DB_PATH}"
	else
		export WALLET_DATABASE_TYPE="mysql"
		log_info "Using MySQL database"
	fi
}

###############################################################################
# Configuration Path Management
###############################################################################

xrp_get_config_paths() {
	local project_root
	project_root="$(cd "${_XRP_COMMON_DIR}/../../.." && pwd)"

	XRP_CONFIG_WATCH="${project_root}/config/wallet/xrp/watch.yaml"
	XRP_CONFIG_KEYGEN="${project_root}/config/wallet/xrp/keygen.yaml"

	export XRP_CONFIG_WATCH XRP_CONFIG_KEYGEN
}

###############################################################################
# Wallet Command Wrappers
# Each wrapper injects the wallet-specific SQLite DB path and WebSocket URL
# so that watch and keygen operate on isolated databases.
###############################################################################

xrp_watch_cmd() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		WALLET_DATABASE_SQLITE_PATH="${SQLITE_WATCH_DB_PATH}" \
			WALLET_RIPPLE_WEBSOCKET_PUBLIC_URL="ws://${XRP_WS_HOST}:${XRP_WS_PORT}" \
			WALLET_RIPPLE_WEBSOCKET_ADMIN_URL="ws://${XRP_WS_HOST}:${XRP_WS_PORT}" \
			"${GOPATH}/bin/watch" "$@"
	else
		WALLET_RIPPLE_WEBSOCKET_PUBLIC_URL="ws://${XRP_WS_HOST}:${XRP_WS_PORT}" \
			WALLET_RIPPLE_WEBSOCKET_ADMIN_URL="ws://${XRP_WS_HOST}:${XRP_WS_PORT}" \
			"${GOPATH}/bin/watch" "$@"
	fi
}

xrp_keygen_cmd() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		WALLET_DATABASE_SQLITE_PATH="${SQLITE_KEYGEN_DB_PATH}" \
			"${GOPATH}/bin/keygen" "$@"
	else
		"${GOPATH}/bin/keygen" "$@"
	fi
}

###############################################################################
# SQLite Schema Initialization
###############################################################################

# Initialize SQLite databases for E2E testing using the composite schemas.
# Must be called after xrp_init_database() sets SQLITE_*_DB_PATH.
xrp_init_sqlite_db() {
	if [ "${DB_TYPE}" != "sqlite" ]; then
		return 0
	fi

	local project_root
	project_root="$(cd "${_XRP_COMMON_DIR}/../../.." && pwd)"

	local watch_schema="${project_root}/tools/sqlc/schemas/sqlite/e2e/01_watch.sql"
	local keygen_schema="${project_root}/tools/sqlc/schemas/sqlite/e2e/02_keygen.sql"

	if [ -f "${watch_schema}" ]; then
		rm -f "${SQLITE_WATCH_DB_PATH}"
		sqlite3 "${SQLITE_WATCH_DB_PATH}" <"${watch_schema}"
		log_info "Watch SQLite DB initialized"
	else
		log_warn "Watch schema not found: ${watch_schema}"
	fi

	if [ -f "${keygen_schema}" ]; then
		rm -f "${SQLITE_KEYGEN_DB_PATH}"
		sqlite3 "${SQLITE_KEYGEN_DB_PATH}" <"${keygen_schema}"
		log_info "Keygen SQLite DB initialized"
	else
		log_warn "Keygen schema not found: ${keygen_schema}"
	fi
}

###############################################################################
# Infrastructure Setup
###############################################################################

xrp_check_prerequisites() {
	log_step "Checking Prerequisites"

	# Check if binaries exist
	if [ ! -f "${GOPATH}/bin/watch" ]; then
		log_error "watch binary not found. Run 'make build-all' first."
		return 1
	fi

	if [ ! -f "${GOPATH}/bin/keygen" ]; then
		log_error "keygen binary not found. Run 'make build-all' first."
		return 1
	fi

	if ! command_exists docker; then
		log_error "docker is not installed or not in PATH."
		return 1
	fi

	if ! command_exists sqlite3; then
		log_error "sqlite3 is not installed or not in PATH."
		return 1
	fi

	log_info "All prerequisites met"
}

xrp_setup_infrastructure() {
	log_step "Setting Up Infrastructure (rippled standalone)"

	# Skip container management when using shared infrastructure (parallel E2E)
	if [ "${E2E_SHARED_INFRASTRUCTURE:-false}" = "true" ]; then
		log_info "Using shared infrastructure (skipping Docker container management)"
		log_substep "Initializing SQLite databases for pattern: ${E2E_PATTERN}"
		xrp_init_sqlite_db
		return 0
	fi

	local project_root
	project_root="$(cd "${_XRP_COMMON_DIR}/../../.." && pwd)"

	log_substep "Starting rippled standalone node"
	docker compose -f "${project_root}/compose.xrp.yaml" up -d rippled

	# Wait for rippled to be healthy
	log_substep "Waiting for rippled on ${XRP_WS_HOST}:${XRP_WS_PORT}..."
	local max_wait=60
	local count=0
	while [ "${count}" -lt "${max_wait}" ]; do
		if docker compose -f "${project_root}/compose.xrp.yaml" \
			exec -T rippled rippled --silent server_info >/dev/null 2>&1; then
			log_info "rippled is ready"
			return 0
		fi
		sleep 2
		count=$((count + 2))
	done

	log_error "rippled failed to start within ${max_wait} seconds"
	return 1
}

###############################################################################
# Ledger Management
###############################################################################

# Advance the rippled ledger in standalone mode.
# Must be called after each transaction submission for it to become validated.
xrp_ledger_accept() {
	local project_root
	project_root="$(cd "${_XRP_COMMON_DIR}/../../.." && pwd)"

	log_info "Advancing rippled ledger (ledger_accept)..."
	docker compose -f "${project_root}/compose.xrp.yaml" \
		exec -T rippled rippled --silent ledger_accept >/dev/null 2>&1 || {
		log_warn "ledger_accept returned non-zero (ledger may already be closed)"
	}
}

###############################################################################
# Cleanup Functions
###############################################################################

xrp_cleanup() {
	log_step "Cleaning Up"

	# Skip container management when using shared infrastructure (parallel E2E)
	if [ "${E2E_SHARED_INFRASTRUCTURE:-false}" = "true" ]; then
		log_info "Using shared infrastructure (skipping cleanup)"
		return 0
	fi

	local project_root
	project_root="$(cd "${_XRP_COMMON_DIR}/../../.." && pwd)"

	log_substep "Stopping rippled node"
	docker compose -f "${project_root}/compose.xrp.yaml" down rippled 2>/dev/null || true
	# Remove the named volume so rippled starts with a clean ledger on next run
	docker volume rm xrp-node-db 2>/dev/null || true

	if [ "${DB_TYPE}" = "sqlite" ]; then
		log_substep "Removing SQLite databases"
		rm -f "${project_root}/data/sqlite/xrp/"*-e2e*.db 2>/dev/null || true
	fi

	log_info "Cleanup completed"
}

xrp_full_reset() {
	log_step "Performing Full Reset"

	# Skip container management when using shared infrastructure (parallel E2E)
	if [ "${E2E_SHARED_INFRASTRUCTURE:-false}" = "true" ]; then
		log_info "Using shared infrastructure (skipping full reset)"
		if [ "${DB_TYPE}" = "sqlite" ]; then
			log_info "Cleaning SQLite databases for pattern: ${E2E_PATTERN}"
			local project_root
			project_root="$(cd "${_XRP_COMMON_DIR}/../../.." && pwd)"
			rm -f "${project_root}/data/sqlite/xrp/watch-e2e-${E2E_PATTERN}.db" \
				"${project_root}/data/sqlite/xrp/keygen-e2e-${E2E_PATTERN}.db" 2>/dev/null || true
			xrp_init_sqlite_db
		fi
		return 0
	fi

	xrp_cleanup

	log_substep "Removing keystore files"
	local project_root
	project_root="$(cd "${_XRP_COMMON_DIR}/../../.." && pwd)"
	rm -f "${project_root}/data/keystore/"*xrp* 2>/dev/null || true

	log_info "Full reset completed"
}

###############################################################################
# Database Query Functions (SQLite)
###############################################################################

# Get the first unallocated address for an account from the watch wallet DB
# Usage: addr=$(xrp_get_address "payment")
xrp_get_address() {
	local account="${1:-payment}"

	if [ "${DB_TYPE}" = "sqlite" ]; then
		sqlite3 "${SQLITE_WATCH_DB_PATH}" \
			"SELECT wallet_address FROM address WHERE coin='xrp' AND account='${account}' AND is_allocated=0 LIMIT 1" \
			2>/dev/null || true
	else
		log_warn "xrp_get_address: MySQL support not yet implemented"
		echo ""
	fi
}

###############################################################################
# Funding Functions
###############################################################################

# Fund an XRP address from the genesis account (standalone mode only).
# Uses rippled's HTTP admin API (port 5005) inside the Docker container via docker compose exec.
# The HTTP admin port is bound to localhost inside the container, so we must use exec to reach it.
#
# Usage: xrp_fund_address "rXXX..." 100
xrp_fund_address() {
	local address="$1"
	local amount_xrp="${2:-100}"
	# Convert XRP to drops (1 XRP = 1,000,000 drops)
	local amount_drops
	amount_drops=$(echo "${amount_xrp} * 1000000" | bc | cut -d'.' -f1)

	log_substep "Funding ${address} with ${amount_xrp} XRP from genesis account"

	local project_root
	project_root="$(cd "${_XRP_COMMON_DIR}/../../.." && pwd)"

	# Use rippled HTTP admin API inside the container (port 5005, localhost only)
	# The submit method with secret+tx_json auto-fills Sequence, Fee, etc.
	local response
	response=$(docker compose -f "${project_root}/compose.xrp.yaml" \
		exec -T rippled \
		curl -s -X POST \
		-H "Content-Type: application/json" \
		-d "{\"method\":\"submit\",\"params\":[{\"secret\":\"${XRP_GENESIS_SECRET}\",\"tx_json\":{\"TransactionType\":\"Payment\",\"Account\":\"${XRP_GENESIS_ACCOUNT}\",\"Destination\":\"${address}\",\"Amount\":\"${amount_drops}\"}}]}" \
		http://127.0.0.1:5005) || {
		log_error "Failed to fund ${address}"
		return 1
	}

	# Check for error in response
	if echo "${response}" | grep -q '"error"'; then
		log_error "rippled API error while funding ${address}: ${response}"
		return 1
	fi

	log_info "Funding transaction submitted for ${address} (${amount_xrp} XRP / ${amount_drops} drops)"

	# Advance ledger to finalize the funding transaction
	xrp_ledger_accept

	log_info "Address funded successfully: ${address} (${amount_xrp} XRP)"
}

###############################################################################
# Utility Functions
###############################################################################

# Extract a file path tagged with [fileName]: from command output
# Usage: file=$(xrp_extract_file_path "$output")
xrp_extract_file_path() {
	local output="$1"
	echo "$output" | grep '^\[fileName\]:' | sed 's/^\[fileName\]: //'
}

# Extract txID from watch send tx output
# Usage: txid=$(xrp_extract_tx_id "$output")
# Returns empty string (and exit 0) when no txID is found.
xrp_extract_tx_id() {
	local output="$1"
	echo "$output" | grep -oE '[0-9A-F]{64}' | head -1 || true
}

###############################################################################
# Initialization
###############################################################################

# Initialize database configuration (runs on source)
xrp_init_database

log_info "XRP common utilities loaded (coin=${XRP_COIN}, ws=${XRP_WS_HOST}:${XRP_WS_PORT}, db_type=${DB_TYPE})"
