#!/usr/bin/env bash

# Bitcoin Cash-Specific Common Utility Functions for E2E Operations
# This file provides shared utility functions specifically for Bitcoin Cash E2E scripts.
#
# Usage:
#   source "$(dirname "${BASH_SOURCE[0]}")/bch_common.sh"
#
# Note: This file sources common.sh automatically, so you don't need to source both.

# Script directory for relative paths
_BCH_COMMON_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the general common utilities
# shellcheck source=../common.sh
source "${_BCH_COMMON_DIR}/../common.sh"

###############################################################################
# BCH Default Configuration
###############################################################################

# Coin identifier
BCH_COIN="bch"

# Encrypted wallet flag (default: false for regtest)
BCH_ENCRYPTED="${BCH_ENCRYPTED:-false}"

# RPC credentials (can be overridden via environment variables)
# Note: Default values are for regtest/development only
BCH_RPC_USER="${BCH_RPC_USER:-${RPC_USER:-xyz}}"
BCH_RPC_PASSWORD="${BCH_RPC_PASSWORD:-${RPC_PASSWORD:-xyz}}"

# MySQL credentials (can be overridden via environment variables)
# Note: Default value is for regtest/development only
BCH_MYSQL_ROOT_PASSWORD="${BCH_MYSQL_ROOT_PASSWORD:-${MYSQL_ROOT_PASSWORD:-root}}"

# Wallet passphrase (only used if BCH_ENCRYPTED=true)
# Note: Default value is for testing only - use strong passphrase in production
BCH_WALLET_PASSPHRASE="${BCH_WALLET_PASSPHRASE:-${WALLET_PASSPHRASE:-test}}"

# Docker volume name (can be overridden via environment variable)
# Note: Docker Compose prepends project name to volume names
BCH_DOCKER_VOLUME_NAME="${BCH_DOCKER_VOLUME_NAME:-go-crypto-wallet_wallet-db}"

# Wallet-specific RPC hosts (for environment variable overrides)
# Note: These are initialized with default values and updated by bch_init_rpc_hosts()
# when E2E_PATTERN is set for parallel execution
BCH_WATCH_WALLET_RPC_HOST="${BCH_WATCH_WALLET_RPC_HOST:-}"
BCH_KEYGEN_WALLET_RPC_HOST="${BCH_KEYGEN_WALLET_RPC_HOST:-}"

# Initialize RPC hosts based on E2E_PATTERN
# Usage: bch_init_rpc_hosts
# Sets BCH_WATCH_WALLET_RPC_HOST and BCH_KEYGEN_WALLET_RPC_HOST
# If E2E_PATTERN is set (e.g., "p1"), wallet names will be suffixed (e.g., "watch-p1")
bch_init_rpc_hosts() {
	local watch_wallet="watch"
	local keygen_wallet="keygen"
	local sign1_wallet="sign1"
	local sign2_wallet="sign2"

	# Use pattern-specific wallet names for parallel E2E execution
	if [ -n "${E2E_PATTERN:-}" ]; then
		watch_wallet="watch-${E2E_PATTERN}"
		keygen_wallet="keygen-${E2E_PATTERN}"
		sign1_wallet="sign1-${E2E_PATTERN}"
		sign2_wallet="sign2-${E2E_PATTERN}"
	fi

	# Set RPC hosts with wallet-specific paths (BCH uses different ports than BTC)
	BCH_WATCH_WALLET_RPC_HOST="127.0.0.1:28332/wallet/${watch_wallet}"
	BCH_KEYGEN_WALLET_RPC_HOST="127.0.0.1:29332/wallet/${keygen_wallet}"
	BCH_SIGN1_WALLET_RPC_HOST="127.0.0.1:30332/wallet/${sign1_wallet}"
	BCH_SIGN2_WALLET_RPC_HOST="127.0.0.1:31332/wallet/${sign2_wallet}"

	export BCH_WATCH_WALLET_RPC_HOST BCH_KEYGEN_WALLET_RPC_HOST
	export BCH_SIGN1_WALLET_RPC_HOST BCH_SIGN2_WALLET_RPC_HOST

	# Export wallet names for use in bitcoin-cli commands
	export BCH_WATCH_WALLET_NAME="${watch_wallet}"
	export BCH_KEYGEN_WALLET_NAME="${keygen_wallet}"
	export BCH_SIGN1_WALLET_NAME="${sign1_wallet}"
	export BCH_SIGN2_WALLET_NAME="${sign2_wallet}"
}

###############################################################################
# Database Configuration
###############################################################################

# Database type: sqlite (default) or mysql
DB_TYPE="${DB_TYPE:-sqlite}"

# E2E Pattern identifier (set by each E2E script, e.g., "p1", "p2", etc.)
# This is used to generate unique SQLite database file names
E2E_PATTERN="${E2E_PATTERN:-}"

# E2E Timestamp for unique database files (generated once per run)
# Format: YYYYMMDD-HHMMSS
E2E_TIMESTAMP="${E2E_TIMESTAMP:-$(date +%Y%m%d-%H%M%S)}"

# SQLite base directory
SQLITE_DB_DIR="${SQLITE_DB_DIR:-./data/sqlite/bch}"

# SQLite schema file location
SQLITE_WATCH_SCHEMA="${SQLITE_WATCH_SCHEMA:-tools/sqlc/schemas_sqlite/01_watch.sql}"
SQLITE_KEYGEN_SCHEMA="${SQLITE_KEYGEN_SCHEMA:-tools/sqlc/schemas_sqlite/02_keygen.sql}"
SQLITE_SIGN_SCHEMA="${SQLITE_SIGN_SCHEMA:-tools/sqlc/schemas_sqlite/03_sign.sql}"

# Generate SQLite database file path for a specific wallet type
# Usage: db_path=$(sqlite_get_db_path "watch")
# Output format: ./data/sqlite/bch/watch-e2e-p1-20260118-123456.db
sqlite_get_db_path() {
	local wallet_type="$1"
	local pattern_suffix=""

	if [ -n "${E2E_PATTERN:-}" ]; then
		pattern_suffix="-e2e-${E2E_PATTERN}-${E2E_TIMESTAMP}"
	else
		pattern_suffix="-e2e-${E2E_TIMESTAMP}"
	fi

	echo "${SQLITE_DB_DIR}/${wallet_type}${pattern_suffix}.db"
}

# Individual wallet DB paths (generated dynamically based on pattern and timestamp)
# These are computed lazily when needed
sqlite_init_db_paths() {
	SQLITE_WATCH_DB_PATH="${SQLITE_WATCH_DB_PATH:-$(sqlite_get_db_path "watch")}"
	SQLITE_KEYGEN_DB_PATH="${SQLITE_KEYGEN_DB_PATH:-$(sqlite_get_db_path "keygen")}"
	SQLITE_SIGN_DB_PATH="${SQLITE_SIGN_DB_PATH:-$(sqlite_get_db_path "sign")}"
	SQLITE_SIGN2_DB_PATH="${SQLITE_SIGN2_DB_PATH:-$(sqlite_get_db_path "sign2")}"

	export SQLITE_WATCH_DB_PATH SQLITE_KEYGEN_DB_PATH SQLITE_SIGN_DB_PATH SQLITE_SIGN2_DB_PATH
}

# Export database environment variables for wallet CLI commands
# Override config file's database.type based on DB_TYPE environment variable
# Note: WALLET_DATABASE_SQLITE_PATH is set per-command in bch_watch_cmd, bch_keygen_cmd, etc.
if [ "${DB_TYPE}" = "sqlite" ]; then
	export WALLET_DATABASE_TYPE="sqlite"
	# Increase max_open_conns for E2E tests to handle concurrent operations
	export WALLET_DATABASE_SQLITE_MAX_OPEN_CONNS=20
elif [ "${DB_TYPE}" = "mysql" ]; then
	export WALLET_DATABASE_TYPE="mysql"
fi

###############################################################################
# Database Abstraction Functions
###############################################################################

# Check if using SQLite database
# Usage: if db_is_sqlite; then ...; fi
db_is_sqlite() {
	[ "${DB_TYPE}" = "sqlite" ]
}

# Check if using MySQL database
# Usage: if db_is_mysql; then ...; fi
db_is_mysql() {
	[ "${DB_TYPE}" = "mysql" ]
}

# Resolve database name to its SQLite path
# Usage: db_path=$(_sqlite_resolve_db_path "watch")
# Returns: The path to the SQLite database file, or empty string with error for unknown names
# Note: This is a helper function to avoid duplicating the case statement in multiple functions
_sqlite_resolve_db_path() {
	local db_name="$1"

	case "$db_name" in
	watch) echo "${SQLITE_WATCH_DB_PATH}" ;;
	keygen) echo "${SQLITE_KEYGEN_DB_PATH}" ;;
	sign | sign1) echo "${SQLITE_SIGN_DB_PATH}" ;;
	sign2) echo "${SQLITE_SIGN2_DB_PATH}" ;;
	*)
		log_error "Unknown database name: $db_name"
		return 1
		;;
	esac
}

# Resolve database name to its schema file path
# Usage: schema_file=$(_sqlite_resolve_schema_path "watch")
# Returns: The path to the schema file, or empty string with error for unknown names
_sqlite_resolve_schema_path() {
	local db_name="$1"

	case "$db_name" in
	watch) echo "${PROJECT_ROOT}/${SQLITE_WATCH_SCHEMA}" ;;
	keygen) echo "${PROJECT_ROOT}/${SQLITE_KEYGEN_SCHEMA}" ;;
	sign | sign1 | sign2) echo "${PROJECT_ROOT}/${SQLITE_SIGN_SCHEMA}" ;;
	*)
		log_error "Unknown database name: $db_name"
		return 1
		;;
	esac
}

# Initialize SQLite database with schema
# Usage: sqlite_init_db "watch"
#        sqlite_init_db "all"  # Initialize all wallet databases (watch, keygen, sign)
# Creates separate database files for each wallet type with pattern and timestamp
sqlite_init_db() {
	local db_name="$1"

	# Initialize database paths first
	sqlite_init_db_paths

	# Create directory if it doesn't exist
	mkdir -p "${SQLITE_DB_DIR}"

	if [ "$db_name" = "all" ]; then
		# Initialize separate databases for each wallet type
		log_info "Initializing SQLite databases for E2E testing"
		log_info "  Pattern: ${E2E_PATTERN:-unset}"
		log_info "  Timestamp: ${E2E_TIMESTAMP}"

		# Initialize each wallet database using helper functions
		for wallet_type in watch keygen sign sign2; do
			local db_path schema_file
			db_path=$(_sqlite_resolve_db_path "$wallet_type")
			schema_file=$(_sqlite_resolve_schema_path "$wallet_type")

			if [ -f "$schema_file" ]; then
				# Remove existing database to start fresh
				rm -f "$db_path"
				log_info "  Creating ${wallet_type} database: $(basename "$db_path")"
				sqlite3 "$db_path" <"$schema_file"
			else
				log_warn "  Schema file not found: $schema_file"
			fi
		done

		log_info "SQLite databases initialized:"
		log_info "  Watch:  ${SQLITE_WATCH_DB_PATH}"
		log_info "  Keygen: ${SQLITE_KEYGEN_DB_PATH}"
		log_info "  Sign:   ${SQLITE_SIGN_DB_PATH}"
		log_info "  Sign2:  ${SQLITE_SIGN2_DB_PATH}"
		return 0
	fi

	# Initialize with specific schema using helper functions
	local db_path schema_file
	db_path=$(_sqlite_resolve_db_path "$db_name") || return 1
	schema_file=$(_sqlite_resolve_schema_path "$db_name") || return 1

	# Check if schema file exists
	if [ ! -f "$schema_file" ]; then
		log_error "Schema file not found: $schema_file"
		return 1
	fi

	log_info "Initializing SQLite database: $db_path (${db_name} schema)"

	# Remove existing and create fresh database
	rm -f "$db_path"
	sqlite3 "$db_path" <"$schema_file"

	log_info "SQLite database initialized: $db_path"
}

# Clean SQLite database files
# Usage: sqlite_clean_db "watch"
#        sqlite_clean_db "all"  # Remove all E2E database files matching current pattern/timestamp
#        sqlite_clean_db "all" "pattern"  # Remove all database files matching the pattern (e.g., "*.db")
sqlite_clean_db() {
	local db_name="$1"
	local pattern="${2:-}"

	if [ "$db_name" = "all" ]; then
		if [ -n "$pattern" ]; then
			# Remove files matching pattern
			log_info "Removing SQLite database files matching: ${SQLITE_DB_DIR}/${pattern}"
			rm -f "${SQLITE_DB_DIR}"/${pattern} 2>/dev/null || true
		else
			# Remove current session's database files using helper function
			for wallet_type in watch keygen sign sign2; do
				local db_path
				db_path=$(_sqlite_resolve_db_path "$wallet_type" 2>/dev/null) || continue
				if [ -n "$db_path" ] && [ -f "$db_path" ]; then
					rm -f "$db_path"
					log_info "Removed: $db_path"
				fi
			done
			# Also remove any old-style e2e.db file
			if [ -f "${SQLITE_DB_DIR}/e2e.db" ]; then
				rm -f "${SQLITE_DB_DIR}/e2e.db"
				log_info "Removed legacy: ${SQLITE_DB_DIR}/e2e.db"
			fi
		fi
		return 0
	fi

	# Initialize paths if not already set
	sqlite_init_db_paths

	# Individual wallet cleanup using helper function
	local db_path
	db_path=$(_sqlite_resolve_db_path "$db_name") || return 1

	if [ -f "$db_path" ]; then
		rm -f "$db_path"
		log_info "Removed SQLite database: $db_path"
	fi
}

# Execute SQLite query
# Usage: sqlite_query "watch" "SELECT * FROM address"
sqlite_query() {
	local db_name="$1"
	local query="$2"
	local db_path

	db_path=$(_sqlite_resolve_db_path "$db_name") || return 1

	sqlite3 -batch "$db_path" "$query"
}

# Execute MySQL query via Docker
# Usage: mysql_query "watch" "SELECT * FROM address"
mysql_query() {
	local db_name="$1"
	local query="$2"

	docker compose exec -e MYSQL_PWD="${BCH_MYSQL_ROOT_PASSWORD}" -T wallet-db mysql -u root "$db_name" -N -e "$query"
}

# Execute database query (abstraction layer)
# Usage: db_query "watch" "SELECT * FROM address"
db_query() {
	local db_name="$1"
	local query="$2"

	if db_is_sqlite; then
		sqlite_query "$db_name" "$query"
	else
		mysql_query "$db_name" "$query"
	fi
}

# Execute database command (for INSERT, UPDATE, DELETE)
# Usage: db_execute "watch" "INSERT INTO ..."
db_execute() {
	local db_name="$1"
	local query="$2"

	if db_is_sqlite; then
		sqlite_query "$db_name" "$query"
	else
		docker compose exec -e MYSQL_PWD="${BCH_MYSQL_ROOT_PASSWORD}" -T wallet-db mysql -u root "$db_name" <<EOF
$query
EOF
	fi
}

###############################################################################
# BCH Config Path Functions
###############################################################################

# Get BCH config file paths (requires PROJECT_ROOT to be set)
# Usage: bch_get_config_paths
#        echo "$BCH_CONFIG_WATCH"
bch_get_config_paths() {
	if [ -z "${PROJECT_ROOT:-}" ]; then
		log_error "PROJECT_ROOT is not set"
		return 1
	fi

	# Initialize RPC hosts based on E2E_PATTERN
	bch_init_rpc_hosts

	BCH_CONFIG_WATCH="${PROJECT_ROOT}/config/wallet/bch/watch.yaml"
	BCH_CONFIG_KEYGEN="${PROJECT_ROOT}/config/wallet/bch/keygen.yaml"
	BCH_CONFIG_SIGN1="${PROJECT_ROOT}/config/wallet/bch/sign1.yaml"
	BCH_CONFIG_SIGN2="${PROJECT_ROOT}/config/wallet/bch/sign2.yaml"
	BCH_CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account/account.yaml"
	BCH_CONFIG_ACCOUNT_2OF3="${PROJECT_ROOT}/config/wallet/account/account_2of3.yaml"
	BCH_CONFIG_ACCOUNT_3OF3="${PROJECT_ROOT}/config/wallet/account/account_3of3.yaml"
}

###############################################################################
# BCH Cleanup Functions
###############################################################################

# Clean BCH generated data files
# Usage: bch_clean_data_files
bch_clean_data_files() {
	log_substep "Cleaning BCH generated data files"

	# Clean address files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/address/bch"

	# Clean fullpubkey files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/fullpubkey/bch"

	# Clean transaction files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/tx/bch"
}

# Clean Bitcoin Cash node wallet data
# Usage: bch_clean_wallet_data [node_list]
# Example: bch_clean_wallet_data "watch keygen"
#          bch_clean_wallet_data "watch keygen sign1 sign2"
bch_clean_wallet_data() {
	local nodes="${1:-watch keygen}"
	log_substep "Cleaning Bitcoin Cash node wallet data"

	# Clean regtest data directories (keeping bitcoin.conf)
	for node in $nodes; do
		wallet_dir="docker/nodes/bch/${node}/regtest"
		if [ -d "$wallet_dir" ]; then
			# Remove all files/dirs except bitcoin.conf
			find "$wallet_dir" -mindepth 1 ! -name 'bitcoin.conf' -exec rm -rf {} + 2>/dev/null || true
			log_info "Cleaned ${node} wallet data"
		fi
	done
}

# Get Docker volume name for database
# Usage: volume_name=$(bch_get_volume_name)
bch_get_volume_name() {
	local project_root="${PROJECT_ROOT:-$(pwd)}"
	# Docker Compose prefixes volume names with the project name (defaults to base name of project directory)
	# Dynamically determine volume name to handle different project directory names
	echo "$(basename "$project_root" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9_-]//g')_wallet-db"
}

# Remove Docker volume with retry
# Usage: bch_remove_volume "volume_name"
bch_remove_volume() {
	local volume_name="$1"
	local max_attempts="${2:-5}"
	local removal_attempts=0

	log_info "Forcefully removing database volume: $volume_name"

	while [ "$removal_attempts" -lt "$max_attempts" ]; do
		if docker volume rm "$volume_name" 2>/dev/null; then
			log_info "Volume removed successfully on attempt $((removal_attempts + 1))"
			return 0
		fi
		removal_attempts=$((removal_attempts + 1))
		if [ "$removal_attempts" -lt "$max_attempts" ]; then
			log_warn "Volume removal failed, retrying in 2 seconds... (attempt $removal_attempts/$max_attempts)"
			sleep 2
		fi
	done

	log_warn "Volume removal failed after $max_attempts attempts"
	return 1
}

# Verify Docker volume is deleted
# Usage: bch_verify_volume_deleted "volume_name" [max_wait_seconds]
bch_verify_volume_deleted() {
	local volume_name="$1"
	local max_wait="${2:-10}"
	local counter=0

	log_info "Verifying volume deletion..."

	while [ "$counter" -lt "$max_wait" ]; do
		if ! docker volume inspect "$volume_name" >/dev/null 2>&1; then
			log_info "Volume successfully deleted"
			return 0
		fi
		counter=$((counter + 1))
		if [ "$counter" -lt "$max_wait" ]; then
			log_warn "Volume still exists, waiting... (${counter}s/${max_wait}s)"
			sleep 1
		fi
	done

	log_error "Volume still exists after ${max_wait}s - this may cause duplicate key errors"
	log_error "Manual cleanup required: docker volume rm -f $volume_name"
	return 1
}

# Full reset: cleanup everything for fresh state
# Usage: bch_full_reset [volume_name] [node_list]
# Example: bch_full_reset
#          bch_full_reset "go-crypto-wallet_wallet-db" "watch keygen sign1 sign2"
# Note: When E2E_SHARED_INFRASTRUCTURE=true, skips Docker container management (for parallel E2E)
bch_full_reset() {
	local volume_name="${1:-${BCH_DOCKER_VOLUME_NAME}}"
	local nodes="${2:-watch keygen}"

	# Skip container management when using shared infrastructure (parallel E2E)
	if [ "${E2E_SHARED_INFRASTRUCTURE:-false}" = "true" ]; then
		log_info "Using shared infrastructure (skipping full reset)"
		# Only clean SQLite databases for this pattern
		if db_is_sqlite; then
			log_info "Cleaning SQLite databases for pattern: ${E2E_PATTERN}"
			sqlite_clean_db "all" "*-e2e-${E2E_PATTERN}-*.db"
		fi
		return 0
	fi

	log_step "Performing full reset for fresh state"

	# Stop and remove containers WITH VOLUMES
	log_info "Stopping Bitcoin Cash containers (with volume removal)..."
	docker compose -f compose.bch.yaml down -v 2>/dev/null || true

	if db_is_sqlite; then
		# SQLite mode: Remove SQLite database files
		log_info "Using SQLite mode (no MySQL container needed)"
		log_info "Removing all SQLite E2E database files..."
		# Remove all E2E database files (pattern: *-e2e-*.db)
		sqlite_clean_db "all" "*-e2e-*.db"
		# Also remove any old-style single e2e.db file
		sqlite_clean_db "all" "e2e.db"

		# Also stop MySQL containers if running from previous runs (for clean state)
		if docker compose -f compose.yaml ps -q 2>/dev/null | grep -q .; then
			log_info "Stopping leftover MySQL containers from previous runs..."
			docker compose -f compose.yaml down -v 2>/dev/null || true
		fi
	else
		# MySQL mode: Stop database container and remove volumes
		log_info "Using MySQL mode"
		log_info "Stopping database container (with volume removal)..."
		docker compose -f compose.yaml down -v 2>/dev/null || true

		# Wait for containers to fully stop before attempting volume removal
		log_info "Waiting for containers to stop completely..."
		sleep 3

		# Remove volume
		bch_remove_volume "$volume_name" || true

		# Verify volume is deleted
		if ! bch_verify_volume_deleted "$volume_name"; then
			return 1
		fi
	fi

	# Clean data files
	bch_clean_data_files

	# Clean Bitcoin Cash wallet data
	bch_clean_wallet_data "$nodes"

	log_info "Full reset complete - system is in fresh state"
	if db_is_sqlite; then
		log_info "Note: SQLite database files were removed for complete cleanup"
	else
		log_info "Note: MySQL database volumes were removed for complete cleanup"
	fi
}

# Basic cleanup: just stop containers
# Usage: bch_cleanup
# Note: When E2E_SHARED_INFRASTRUCTURE=true, skips Docker container management (for parallel E2E)
bch_cleanup() {
	# Skip container management when using shared infrastructure (parallel E2E)
	if [ "${E2E_SHARED_INFRASTRUCTURE:-false}" = "true" ]; then
		log_info "Using shared infrastructure (skipping cleanup)"
		return 0
	fi

	log_step "Cleaning up containers and state"

	log_info "Stopping Bitcoin Cash containers..."
	docker compose -f compose.bch.yaml down -v 2>/dev/null || true

	if db_is_sqlite; then
		log_info "SQLite mode: No database container needed"
		# Also stop MySQL containers if running from previous runs (for clean state)
		if docker compose -f compose.yaml ps -q 2>/dev/null | grep -q .; then
			log_info "Stopping leftover MySQL containers from previous runs..."
			docker compose -f compose.yaml down -v 2>/dev/null || true
		fi
	else
		log_info "Stopping database container..."
		docker compose -f compose.yaml down -v 2>/dev/null || true
	fi

	log_info "Cleanup complete"
}

###############################################################################
# BCH RPC Command Wrappers
###############################################################################

# Wrapper for bitcoin-cli commands with BCH RPC credentials
# Usage: bch_cli <container_name> [bitcoin-cli args...]
# Example: bch_cli "bch-watch" -rpcwallet=watch getbalance
bch_cli() {
	local container=$1
	shift
	local rpc_user="${BCH_RPC_USER}"
	local rpc_password="${BCH_RPC_PASSWORD}"

	docker exec "$container" bitcoin-cli -regtest \
		-rpcuser="${rpc_user}" \
		-rpcpassword="${rpc_password}" \
		"$@"
}

###############################################################################
# BCH Wallet Command Wrappers
###############################################################################

# Internal helper function for wallet command wrappers
# Handles environment variable setup for both SQLite and MySQL modes
# Usage: _bch_wallet_cmd <wallet_name> <db_path> <rpc_host> [args...]
_bch_wallet_cmd() {
	local wallet_name="$1"
	local db_path="$2"
	local rpc_host="$3"
	shift 3

	if db_is_sqlite; then
		WALLET_DATABASE_SQLITE_PATH="${db_path}" \
			WALLET_BITCOIN_HOST="${rpc_host}" "${wallet_name}" "$@"
	else
		WALLET_BITCOIN_HOST="${rpc_host}" "${wallet_name}" "$@"
	fi
}

# Wrapper for watch wallet commands with host and database path override
# When DB_TYPE=sqlite, sets WALLET_DATABASE_SQLITE_PATH to the watch-specific database
# Usage: bch_watch_cmd [args...]
bch_watch_cmd() {
	_bch_wallet_cmd "watch" "${SQLITE_WATCH_DB_PATH}" "${BCH_WATCH_WALLET_RPC_HOST}" "$@"
}

# Wrapper for keygen wallet commands with host and database path override
# When DB_TYPE=sqlite, sets WALLET_DATABASE_SQLITE_PATH to the keygen-specific database
# Usage: bch_keygen_cmd [args...]
bch_keygen_cmd() {
	_bch_wallet_cmd "keygen" "${SQLITE_KEYGEN_DB_PATH}" "${BCH_KEYGEN_WALLET_RPC_HOST}" "$@"
}

# Wrapper for sign1 wallet commands with host and database path override
# When DB_TYPE=sqlite, sets WALLET_DATABASE_SQLITE_PATH to the sign-specific database
# Usage: bch_sign1_cmd [args...]
bch_sign1_cmd() {
	bch_sign_cmd 1 "$@"
}

# Wrapper for sign2 wallet commands with host and database path override
# When DB_TYPE=sqlite, sets WALLET_DATABASE_SQLITE_PATH to the sign2-specific database
# Usage: bch_sign2_cmd [args...]
bch_sign2_cmd() {
	bch_sign_cmd 2 "$@"
}

# Generic wrapper for sign wallet commands with dynamic sign number
# When DB_TYPE=sqlite, sets WALLET_DATABASE_SQLITE_PATH to the appropriate sign database
# Usage: bch_sign_cmd <sign_number> [args...]
# Example: bch_sign_cmd 1 -c config.yaml --coin bch create seed
#          bch_sign_cmd 2 -c config.yaml --coin bch create hdkey
bch_sign_cmd() {
	local sign_num="$1"
	shift

	local db_path
	local rpc_host
	if [ "${sign_num}" = "1" ]; then
		db_path="${SQLITE_SIGN_DB_PATH}"
		rpc_host="${BCH_SIGN1_WALLET_RPC_HOST}"
	elif [ "${sign_num}" = "2" ]; then
		db_path="${SQLITE_SIGN2_DB_PATH}"
		rpc_host="${BCH_SIGN2_WALLET_RPC_HOST}"
	else
		log_error "Invalid sign number provided to bch_sign_cmd: ${sign_num}"
		return 1
	fi

	_bch_wallet_cmd "sign${sign_num}" "${db_path}" "${rpc_host}" "$@"
}

###############################################################################
# BCH Infrastructure Setup Functions
###############################################################################

# Check BCH prerequisites
# Usage: bch_check_prerequisites [command_list]
# Example: bch_check_prerequisites "watch keygen"
#          bch_check_prerequisites "watch keygen sign1 sign2"
bch_check_prerequisites() {
	local commands="${1:-watch keygen}"
	log_step "Checking prerequisites"

	# Check Docker and Docker Compose
	check_docker || return 1

	# Check CLI commands
	for cmd in $commands; do
		if ! command_exists "$cmd"; then
			log_error "$cmd command is not available"
			log_error "Please build the project first: make build"
			return 1
		fi
	done

	log_info "All prerequisites satisfied"
}

# Setup BCH infrastructure
# Usage: bch_setup_infrastructure [container_list]
# Example: bch_setup_infrastructure "bch-watch bch-keygen"
#          bch_setup_infrastructure "bch-watch bch-keygen bch-sign1 bch-sign2"
# Note: When E2E_SHARED_INFRASTRUCTURE=true, skips Docker container management (for parallel E2E)
bch_setup_infrastructure() {
	local containers="${1:-bch-watch bch-keygen}"

	# Skip infrastructure setup when using shared infrastructure (parallel E2E)
	if [ "${E2E_SHARED_INFRASTRUCTURE:-false}" = "true" ]; then
		log_info "Using shared infrastructure (skipping Docker container management)"
		log_substep "Initializing SQLite databases for pattern: ${E2E_PATTERN}"
		sqlite_init_db "all"
		return 0
	fi

	log_step "Setting up infrastructure (DB_TYPE=${DB_TYPE})"

	if db_is_sqlite; then
		# SQLite mode: Initialize separate SQLite databases for each wallet type
		# No Docker MySQL container is started - this is faster and simpler
		log_substep "Initializing SQLite databases (no MySQL container needed)"
		sqlite_init_db "all"
		# Database paths are logged by sqlite_init_db
	else
		# MySQL mode: Start database container
		log_substep "Starting MySQL database container"
		docker compose -f compose.yaml up -d
		log_info "MySQL database container started"

		# Wait for database to be healthy
		wait_for_healthy "wallet-db" 90

		# Wait for database migrations to complete
		wait_for_migrations 120
	fi

	# Start Bitcoin Cash nodes
	log_substep "Starting Bitcoin Cash node containers"
	# shellcheck disable=SC2086
	docker compose -f compose.bch.yaml up -d $containers
	log_info "Bitcoin Cash node containers started"

	# Wait for containers to be healthy
	log_substep "Waiting for containers to be healthy"
	for container in $containers; do
		wait_for_healthy "$container" 90
	done

	log_info "All containers are healthy"
}

# Setup BCH wallets
# Usage: bch_setup_wallets [wallet_names]
# Example: bch_setup_wallets "watch keygen"
#          bch_setup_wallets "watch keygen sign1 sign2"
# Note: Container name is derived as "bch-${wallet_name}"
#       If E2E_PATTERN is set (e.g., "p1"), wallet names will be suffixed (e.g., "watch-p1")
bch_setup_wallets() {
	local wallets="${1:-watch keygen}"
	log_step "Setting up Bitcoin Cash wallets"

	if [ -n "${E2E_PATTERN:-}" ]; then
		log_info "E2E Pattern: ${E2E_PATTERN} (wallets will have pattern suffix)"
	fi

	for wallet in $wallets; do
		local container="bch-${wallet}"
		local wallet_name="$wallet"

		# Use pattern-specific wallet name for parallel E2E execution
		if [ -n "${E2E_PATTERN:-}" ]; then
			wallet_name="${wallet}-${E2E_PATTERN}"
			log_info "Creating/loading wallet: ${wallet_name} in container ${container}"
		else
			log_info "Creating/loading wallet: ${wallet_name} in container ${container}"
		fi

		if ! bch_create_wallet_if_needed "$container" "$wallet_name"; then
			log_error "Failed to create/load wallet ${wallet_name} in ${container}"
			return 1
		fi

		# Verify wallet is loaded
		if bch_wallet_exists "$container" "$wallet_name"; then
			log_info "  ✓ Wallet ${wallet_name} is ready and loaded"
		else
			log_error "  ✗ Wallet ${wallet_name} creation reported success but wallet not found!"
			log_error "    Container: ${container}"
			log_error "    Wallet name: ${wallet_name}"
			# List available wallets for debugging
			log_error "    Available wallets:"
			bch_cli "$container" listwallets 2>/dev/null || log_error "    Failed to list wallets"
			return 1
		fi
	done

	log_info "All wallets are ready"
}

###############################################################################
# BCH Transaction Helper Functions
###############################################################################

# Log detailed error message for "No utxo" errors
# Usage: bch_log_no_utxo_error
bch_log_no_utxo_error() {
	log_error "Transaction creation failed"
	log_error "This could indicate:"
	log_error "  - No payment requests in database"
	log_error "  - No UTXOs available for payment account"
	log_error "  - UTXOs not mature enough (need 100+ confirmations)"
}

# Wait for balance to update after UTXO generation
# Usage: bch_wait_for_balance [max_wait] [wait_interval]
bch_wait_for_balance() {
	local max_wait="${1:-60}"
	local wait_interval="${2:-3}"
	local elapsed=0
	local balance_found=false

	log_info "Waiting for blockchain sync and balance update..."

	while [ "$elapsed" -lt "$max_wait" ]; do
		# Check balance using Bitcoin Cash RPC
		# BCH watch wallet uses watch-only addresses, so we need to include them
		# getbalance "*" minconf include_watchonly
		# Use pattern-specific wallet name if E2E_PATTERN is set
		local wallet_name="${BCH_WATCH_WALLET_NAME:-watch}"
		local balance
		balance=$(bch_cli "bch-watch" -rpcwallet="${wallet_name}" getbalance "*" 1 true 2>&1 || echo "0")

		# Check if we have any balance
		if [ -n "$balance" ] && [ "$(echo "$balance > 0" | bc -l 2>/dev/null || echo 0)" -eq 1 ]; then
			log_info "Payment account balance verified: ${balance} BCH (took ${elapsed}s)"
			balance_found=true
			break
		fi

		sleep "$wait_interval"
		elapsed=$((elapsed + wait_interval))
		if [ "$elapsed" -lt "$max_wait" ]; then
			log_info "Still waiting for balance update... (${elapsed}s/${max_wait}s)"
		fi
	done

	if [ "$balance_found" = false ]; then
		log_error "Balance not detected within ${max_wait}s"
		log_error "This indicates a failure in UTXO generation or blockchain sync"
		log_error "Please check:"
		log_error "  - Bitcoin Cash node logs: docker compose -f compose.bch.yaml logs bch-watch"
		log_error "  - Block generation succeeded"
		log_error "  - Address import into watch wallet succeeded"
		return 1
	fi

	return 0
}

# Generate test UTXOs for a payment address
# Usage: bch_generate_test_utxos "payment_address" [block_count]
bch_generate_test_utxos() {
	local payment_address="$1"
	local block_count="${2:-101}"

	log_info "Using payment address: $payment_address"

	# Generate blocks with coinbase reward to payment address
	# Use pattern-specific wallet name if E2E_PATTERN is set
	local wallet_name="${BCH_WATCH_WALLET_NAME:-watch}"
	log_info "Generating $block_count blocks to create mature coinbase for testing..."
	log_info "Using wallet: ${wallet_name}"
	bch_cli "bch-watch" -rpcwallet="${wallet_name}" generatetoaddress "$block_count" "$payment_address" >/dev/null

	log_info "Test UTXOs generated successfully"

	# Rescan blockchain to detect the newly generated UTXOs
	# This is necessary because addresses were imported before blocks were generated
	log_info "Rescanning blockchain to detect UTXOs..."
	if ! output=$(bch_cli "bch-watch" -rpcwallet="${wallet_name}" rescanblockchain 0 2>&1); then
		log_error "Blockchain rescan failed with output:"
		log_error "$output"
		return 1
	fi
	log_info "Blockchain rescan completed"
}

# Extract file path from command output
# Expects format: "[fileName]: /path/to/file"
# Usage: file_path=$(bch_extract_file_path "$output")
bch_extract_file_path() {
	local output="$1"
	# Extract everything after [fileName]: and take only the first line
	# Output format: [fileName]: /path/to/file.hex\n[signedCount]: 1\n...
	echo "${output##*\[fileName\]: }" | head -n1
}

###############################################################################
# BCH Payment Request Functions
###############################################################################

# Get sender address from database
# Usage: sender_address=$(bch_get_sender_address "payment")
bch_get_sender_address() {
	local account="${1:-payment}"

	db_query "watch" "SELECT wallet_address FROM address WHERE coin='bch' AND account='${account}' LIMIT 1"
}

# Generate receiver addresses (using legacy format for BCH regtest)
# Usage: receivers=$(bch_generate_receiver_addresses 3)
bch_generate_receiver_addresses() {
	local count="${1:-3}"
	local addresses=""

	log_substep "Generating receiver addresses for payment requests" >&2
	log_info "Creating $count receiver addresses in watch wallet..." >&2

	# Use pattern-specific wallet name if E2E_PATTERN is set
	local wallet_name="${BCH_WATCH_WALLET_NAME:-watch}"

	for i in $(seq 1 "$count"); do
		local addr
		# BCH uses legacy address format (no SegWit)
		addr=$(bch_cli "bch-watch" -rpcwallet="${wallet_name}" getnewaddress "" legacy)
		if [ -n "$addresses" ]; then
			addresses="$addresses $addr"
		else
			addresses="$addr"
		fi
		log_info "  $i. $addr" >&2
	done

	echo "$addresses"
}

# Insert payment requests into database
# Usage: bch_insert_payment_requests "sender_address" "receiver1 receiver2 receiver3" "0.001 0.002 0.0015"
bch_insert_payment_requests() {
	local sender_address="$1"
	local receivers="$2"
	local amounts="$3"
	local account="${4:-payment}"

	log_substep "Inserting payment requests into database"

	# Convert space-separated lists to arrays
	local receiver_array=()
	read -ra receiver_array <<<"$receivers"
	local amount_array=()
	read -ra amount_array <<<"$amounts"

	# Build INSERT statement
	local values=""
	for i in "${!receiver_array[@]}"; do
		if [ -n "$values" ]; then
			values="$values,"
		fi
		# SQLite uses 0/1 for boolean, MySQL uses false/true
		if db_is_sqlite; then
			values="$values
	('bch', NULL, '${sender_address}', '${account}', '${receiver_array[$i]}', '${amount_array[$i]}', 0)"
		else
			values="$values
	('bch', NULL, '${sender_address}', '${account}', '${receiver_array[$i]}', ${amount_array[$i]}, false)"
		fi
	done

	local sql_statements
	sql_statements="DELETE FROM payment_request WHERE coin='bch';
INSERT INTO payment_request (coin, payment_id, sender_address, sender_account, receiver_address, amount, is_done)
VALUES${values};"

	db_execute "watch" "$sql_statements"

	# shellcheck disable=SC2181
	if [ $? -ne 0 ]; then
		log_error "Failed to insert payment requests"
		return 1
	fi

	# Verify payment requests were created
	local count
	if db_is_sqlite; then
		count=$(db_query "watch" "SELECT COUNT(*) FROM payment_request WHERE coin='bch' AND is_done=0")
	else
		count=$(db_query "watch" "SELECT COUNT(*) FROM payment_request WHERE coin='bch' AND is_done=false")
	fi

	log_info "Created $count payment requests"

	if [ "$count" -eq 0 ]; then
		log_error "No payment requests were created"
		return 1
	fi

	log_info "Payment requests ready for transaction creation"
}

###############################################################################
# BCH Address Export Functions (for non-descriptor workflows)
###############################################################################

# Export addresses and import to watch wallet (legacy BCH workflow)
# Usage: bch_export_and_import_addresses "client deposit payment stored"
bch_export_and_import_addresses() {
	local accounts="${1:-client deposit payment stored}"
	log_substep "Exporting and importing addresses"

	for account in $accounts; do
		log_info "Exporting ${account} addresses"
		local file_output
		file_output=$(keygen -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" export address --account "$account")
		local address_file="${file_output##*\[fileName\]: }"

		log_info "Importing ${account} addresses"
		watch -c "${BCH_CONFIG_WATCH}" --coin "${BCH_COIN}" import address --file "${address_file}"
	done

	log_info "All addresses exported and imported"
}

###############################################################################
# BCH Argument Parsing Helper
###############################################################################

# Parse common e2e script arguments
# Usage: bch_parse_args "$@"
#        This sets: CLEANUP_ONLY, RESET_STATE, VERBOSE, NON_INTERACTIVE
bch_parse_args() {
	CLEANUP_ONLY="${CLEANUP_ONLY:-false}"
	RESET_STATE="${RESET_STATE:-false}"
	VERBOSE="${VERBOSE:-false}"
	NON_INTERACTIVE="${NON_INTERACTIVE:-false}"

	while [ $# -gt 0 ]; do
		case "$1" in
		--cleanup)
			CLEANUP_ONLY=true
			shift
			;;
		--reset)
			RESET_STATE=true
			shift
			;;
		--verbose)
			VERBOSE=true
			set -x
			shift
			;;
		--non-interactive)
			NON_INTERACTIVE=true
			shift
			;;
		-h | --help)
			# Return special value to indicate help was requested
			return 2
			;;
		*)
			log_error "Unknown option: $1"
			return 1
			;;
		esac
	done

	export CLEANUP_ONLY RESET_STATE VERBOSE NON_INTERACTIVE
	return 0
}
