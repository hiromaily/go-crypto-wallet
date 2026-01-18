#!/usr/bin/env bash

# Bitcoin-Specific Common Utility Functions for E2E Operations
# This file provides shared utility functions specifically for Bitcoin E2E scripts.
#
# Usage:
#   source "$(dirname "${BASH_SOURCE[0]}")/btc_common.sh"
#
# Note: This file sources common.sh automatically, so you don't need to source both.
#
# Database Support:
#   This file supports both MySQL and SQLite databases for E2E testing.
#   Set DB_TYPE environment variable to choose:
#     DB_TYPE=mysql (default) - Uses Docker MySQL container
#     DB_TYPE=sqlite          - Uses local SQLite file (no Docker required)

# Script directory for relative paths
_BTC_COMMON_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the general common utilities
# shellcheck source=../common.sh
source "${_BTC_COMMON_DIR}/../common.sh"

###############################################################################
# BTC Default Configuration
###############################################################################

# Coin identifier
BTC_COIN="btc"

# Encrypted wallet flag (default: false for regtest)
BTC_ENCRYPTED="${BTC_ENCRYPTED:-false}"

# RPC credentials (can be overridden via environment variables)
# Note: Default values are for regtest/development only
BTC_RPC_USER="${BTC_RPC_USER:-${RPC_USER:-xyz}}"
BTC_RPC_PASSWORD="${BTC_RPC_PASSWORD:-${RPC_PASSWORD:-xyz}}"

# MySQL credentials (can be overridden via environment variables)
# Note: Default value is for regtest/development only
BTC_MYSQL_ROOT_PASSWORD="${BTC_MYSQL_ROOT_PASSWORD:-${MYSQL_ROOT_PASSWORD:-root}}"

# Wallet passphrase (only used if BTC_ENCRYPTED=true)
# Note: Default value is for testing only - use strong passphrase in production
BTC_WALLET_PASSPHRASE="${BTC_WALLET_PASSPHRASE:-${WALLET_PASSPHRASE:-test}}"

# Docker volume name (can be overridden via environment variable)
# Note: Docker Compose prepends project name to volume names
BTC_DOCKER_VOLUME_NAME="${BTC_DOCKER_VOLUME_NAME:-go-crypto-wallet_wallet-db}"

# Wallet-specific RPC hosts (for environment variable overrides)
BTC_WATCH_WALLET_RPC_HOST="${BTC_WATCH_WALLET_RPC_HOST:-127.0.0.1:18332/wallet/watch}"
BTC_KEYGEN_WALLET_RPC_HOST="${BTC_KEYGEN_WALLET_RPC_HOST:-127.0.0.1:19332/wallet/keygen}"

###############################################################################
# Database Configuration
###############################################################################

# Database type: sqlite (default) or mysql
DB_TYPE="${DB_TYPE:-sqlite}"

# SQLite configuration
# For E2E testing, we use a single database file containing all schemas
# This simplifies testing as all wallet commands share the same database
SQLITE_DB_PATH="${SQLITE_DB_PATH:-./data/sqlite/btc/e2e.db}"

# Individual wallet DB paths (kept for backward compatibility, default to shared DB)
SQLITE_WATCH_DB_PATH="${SQLITE_WATCH_DB_PATH:-${SQLITE_DB_PATH}}"
SQLITE_KEYGEN_DB_PATH="${SQLITE_KEYGEN_DB_PATH:-${SQLITE_DB_PATH}}"
SQLITE_SIGN_DB_PATH="${SQLITE_SIGN_DB_PATH:-${SQLITE_DB_PATH}}"

# SQLite schema file location
SQLITE_WATCH_SCHEMA="${SQLITE_WATCH_SCHEMA:-tools/sqlc/schemas_sqlite/01_watch.sql}"
SQLITE_KEYGEN_SCHEMA="${SQLITE_KEYGEN_SCHEMA:-tools/sqlc/schemas_sqlite/02_keygen.sql}"
SQLITE_SIGN_SCHEMA="${SQLITE_SIGN_SCHEMA:-tools/sqlc/schemas_sqlite/03_sign.sql}"

# Export database environment variables for wallet CLI commands
# When DB_TYPE=sqlite, set WALLET_DATABASE_TYPE and WALLET_DATABASE_SQLITE_PATH
# so that wallet commands (watch, keygen, sign) use SQLite instead of MySQL
if [ "${DB_TYPE}" = "sqlite" ]; then
	export WALLET_DATABASE_TYPE="sqlite"
	export WALLET_DATABASE_SQLITE_PATH="${SQLITE_DB_PATH}"
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

# Initialize SQLite database with schema
# Usage: sqlite_init_db "watch"
#        sqlite_init_db "all"  # Initialize with all schemas
# Creates the database file and applies the schema
sqlite_init_db() {
	local db_name="$1"
	local db_path="${SQLITE_DB_PATH}"

	# Create directory if it doesn't exist
	mkdir -p "$(dirname "$db_path")"

	if [ "$db_name" = "all" ]; then
		# Initialize with all schemas for E2E testing
		log_info "Initializing SQLite database with all schemas: $db_path"

		# Remove existing database to start fresh
		rm -f "$db_path"

		# Apply all schemas
		for schema_name in watch keygen sign; do
			local schema_file
			case "$schema_name" in
			watch) schema_file="${PROJECT_ROOT}/${SQLITE_WATCH_SCHEMA}" ;;
			keygen) schema_file="${PROJECT_ROOT}/${SQLITE_KEYGEN_SCHEMA}" ;;
			sign) schema_file="${PROJECT_ROOT}/${SQLITE_SIGN_SCHEMA}" ;;
			esac

			if [ -f "$schema_file" ]; then
				log_info "  Applying ${schema_name} schema..."
				sqlite3 "$db_path" <"$schema_file"
			else
				log_warn "  Schema file not found: $schema_file"
			fi
		done

		log_info "SQLite database initialized with all schemas"
		return 0
	fi

	# Initialize with specific schema
	local schema_file
	case "$db_name" in
	watch) schema_file="${PROJECT_ROOT}/${SQLITE_WATCH_SCHEMA}" ;;
	keygen) schema_file="${PROJECT_ROOT}/${SQLITE_KEYGEN_SCHEMA}" ;;
	sign | sign1 | sign2) schema_file="${PROJECT_ROOT}/${SQLITE_SIGN_SCHEMA}" ;;
	*)
		log_error "Unknown database name: $db_name"
		return 1
		;;
	esac

	# Check if schema file exists
	if [ ! -f "$schema_file" ]; then
		log_error "Schema file not found: $schema_file"
		return 1
	fi

	log_info "Initializing SQLite database: $db_path (${db_name} schema)"

	# Create database and apply schema
	sqlite3 "$db_path" <"$schema_file"

	log_info "SQLite database initialized: $db_path"
}

# Clean SQLite database
# Usage: sqlite_clean_db "watch"
#        sqlite_clean_db "all"  # Remove the shared E2E database
sqlite_clean_db() {
	local db_name="$1"
	local db_path="${SQLITE_DB_PATH}"

	if [ "$db_name" = "all" ]; then
		if [ -f "$db_path" ]; then
			rm -f "$db_path"
			log_info "Removed SQLite database: $db_path"
		fi
		return 0
	fi

	# For backward compatibility, individual wallet cleanup
	case "$db_name" in
	watch) db_path="${SQLITE_WATCH_DB_PATH}" ;;
	keygen) db_path="${SQLITE_KEYGEN_DB_PATH}" ;;
	sign | sign1 | sign2) db_path="${SQLITE_SIGN_DB_PATH}" ;;
	*)
		log_error "Unknown database name: $db_name"
		return 1
		;;
	esac

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

	case "$db_name" in
	watch) db_path="${SQLITE_WATCH_DB_PATH}" ;;
	keygen) db_path="${SQLITE_KEYGEN_DB_PATH}" ;;
	sign | sign1 | sign2) db_path="${SQLITE_SIGN_DB_PATH}" ;;
	*)
		log_error "Unknown database name: $db_name"
		return 1
		;;
	esac

	sqlite3 -batch "$db_path" "$query"
}

# Execute MySQL query via Docker
# Usage: mysql_query "watch" "SELECT * FROM address"
mysql_query() {
	local db_name="$1"
	local query="$2"

	docker compose exec -e MYSQL_PWD="${BTC_MYSQL_ROOT_PASSWORD}" -T wallet-db mysql -u root "$db_name" -N -e "$query"
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
		docker compose exec -e MYSQL_PWD="${BTC_MYSQL_ROOT_PASSWORD}" -T wallet-db mysql -u root "$db_name" <<EOF
$query
EOF
	fi
}

###############################################################################
# BTC Config Path Functions
###############################################################################

# Get BTC config file paths (requires PROJECT_ROOT to be set)
# Usage: btc_get_config_paths
#        echo "$BTC_CONFIG_WATCH"
btc_get_config_paths() {
	if [ -z "${PROJECT_ROOT:-}" ]; then
		log_error "PROJECT_ROOT is not set"
		return 1
	fi

	BTC_CONFIG_WATCH="${PROJECT_ROOT}/config/wallet/btc/watch.yaml"
	BTC_CONFIG_KEYGEN="${PROJECT_ROOT}/config/wallet/btc/keygen.yaml"
	BTC_CONFIG_SIGN1="${PROJECT_ROOT}/config/wallet/btc/sign1.yaml"
	BTC_CONFIG_SIGN2="${PROJECT_ROOT}/config/wallet/btc/sign2.yaml"
	BTC_CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account/account.yaml"
	BTC_CONFIG_ACCOUNT_2OF3="${PROJECT_ROOT}/config/wallet/account/account_2of3.yaml"
	BTC_CONFIG_ACCOUNT_3OF3="${PROJECT_ROOT}/config/wallet/account/account_3of3.yaml"
}

###############################################################################
# BTC Cleanup Functions
###############################################################################

# Clean BTC generated data files
# Usage: btc_clean_data_files
btc_clean_data_files() {
	log_substep "Cleaning BTC generated data files"

	# Clean address files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/address/btc"

	# Clean fullpubkey files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/fullpubkey/btc"

	# Clean descriptor files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/descriptor/btc"

	# Clean transaction files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/tx/btc"
}

# Clean Bitcoin node wallet data
# Usage: btc_clean_wallet_data [node_list]
# Example: btc_clean_wallet_data "watch keygen"
#          btc_clean_wallet_data "watch keygen sign1 sign2"
btc_clean_wallet_data() {
	local nodes="${1:-watch keygen}"
	log_substep "Cleaning Bitcoin node wallet data"

	# Clean regtest data directories (keeping bitcoin.conf)
	for node in $nodes; do
		wallet_dir="docker/nodes/btc/${node}/regtest"
		if [ -d "$wallet_dir" ]; then
			# Remove all files/dirs except bitcoin.conf
			find "$wallet_dir" -mindepth 1 ! -name 'bitcoin.conf' -exec rm -rf {} + 2>/dev/null || true
			log_info "Cleaned ${node} wallet data"
		fi
	done
}

# Get Docker volume name for database
# Usage: volume_name=$(btc_get_volume_name)
btc_get_volume_name() {
	local project_root="${PROJECT_ROOT:-$(pwd)}"
	# Docker Compose prefixes volume names with the project name (defaults to base name of project directory)
	# Dynamically determine volume name to handle different project directory names
	echo "$(basename "$project_root" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9_-]//g')_wallet-db"
}

# Remove Docker volume with retry
# Usage: btc_remove_volume "volume_name"
btc_remove_volume() {
	local volume_name="$1"
	local max_attempts="${2:-5}"
	local removal_attempts=0

	log_info "Forcefully removing database volume: $volume_name"

	while [ $removal_attempts -lt $max_attempts ]; do
		if docker volume rm "$volume_name" 2>/dev/null; then
			log_info "Volume removed successfully on attempt $((removal_attempts + 1))"
			return 0
		fi
		removal_attempts=$((removal_attempts + 1))
		if [ $removal_attempts -lt $max_attempts ]; then
			log_warn "Volume removal failed, retrying in 2 seconds... (attempt $removal_attempts/$max_attempts)"
			sleep 2
		fi
	done

	log_warn "Volume removal failed after $max_attempts attempts"
	return 1
}

# Verify Docker volume is deleted
# Usage: btc_verify_volume_deleted "volume_name" [max_wait_seconds]
btc_verify_volume_deleted() {
	local volume_name="$1"
	local max_wait="${2:-10}"
	local counter=0

	log_info "Verifying volume deletion..."

	while [ $counter -lt $max_wait ]; do
		if ! docker volume inspect "$volume_name" >/dev/null 2>&1; then
			log_info "Volume successfully deleted"
			return 0
		fi
		counter=$((counter + 1))
		if [ $counter -lt $max_wait ]; then
			log_warn "Volume still exists, waiting... (${counter}s/${max_wait}s)"
			sleep 1
		fi
	done

	log_error "Volume still exists after ${max_wait}s - this may cause duplicate key errors"
	log_error "Manual cleanup required: docker volume rm -f $volume_name"
	return 1
}

# Full reset: cleanup everything for fresh state
# Usage: btc_full_reset [volume_name] [node_list]
# Example: btc_full_reset
#          btc_full_reset "go-crypto-wallet_wallet-db" "watch keygen sign1 sign2"
btc_full_reset() {
	local volume_name="${1:-${BTC_DOCKER_VOLUME_NAME}}"
	local nodes="${2:-watch keygen}"
	log_step "Performing full reset for fresh state"

	# Stop and remove containers WITH VOLUMES
	log_info "Stopping Bitcoin containers (with volume removal)..."
	docker compose -f compose.btc.yaml down -v 2>/dev/null || true

	if db_is_sqlite; then
		# SQLite mode: Remove SQLite database file
		log_info "Removing SQLite database file..."
		sqlite_clean_db "all"
	else
		# MySQL mode: Stop database container and remove volumes
		log_info "Stopping database container (with volume removal)..."
		docker compose -f compose.yaml down -v 2>/dev/null || true

		# Wait for containers to fully stop before attempting volume removal
		log_info "Waiting for containers to stop completely..."
		sleep 3

		# Remove volume
		btc_remove_volume "$volume_name" || true

		# Verify volume is deleted
		if ! btc_verify_volume_deleted "$volume_name"; then
			return 1
		fi
	fi

	# Clean data files
	btc_clean_data_files

	# Clean Bitcoin wallet data
	btc_clean_wallet_data "$nodes"

	log_info "Full reset complete - system is in fresh state"
	if db_is_sqlite; then
		log_info "Note: SQLite database files were removed for complete cleanup"
	else
		log_info "Note: Database volumes were removed for complete cleanup"
	fi
}

# Basic cleanup: just stop containers
# Usage: btc_cleanup
btc_cleanup() {
	log_step "Cleaning up containers and state"

	log_info "Stopping Bitcoin containers..."
	docker compose -f compose.btc.yaml down -v 2>/dev/null || true

	if db_is_sqlite; then
		log_info "SQLite mode: No database container to stop"
	else
		log_info "Stopping database container..."
		docker compose -f compose.yaml down -v 2>/dev/null || true
	fi

	log_info "Cleanup complete"
}

###############################################################################
# BTC Wallet Command Wrappers
###############################################################################

# Wrapper for watch wallet commands with host override
# Database environment is already exported via WALLET_DATABASE_TYPE and
# WALLET_DATABASE_SQLITE_PATH when DB_TYPE=sqlite
# Usage: btc_watch_cmd [args...]
btc_watch_cmd() {
	WALLET_BITCOIN_HOST="${BTC_WATCH_WALLET_RPC_HOST}" watch "$@"
}

# Wrapper for keygen wallet commands with host override
# Database environment is already exported via WALLET_DATABASE_TYPE and
# WALLET_DATABASE_SQLITE_PATH when DB_TYPE=sqlite
# Usage: btc_keygen_cmd [args...]
btc_keygen_cmd() {
	WALLET_BITCOIN_HOST="${BTC_KEYGEN_WALLET_RPC_HOST}" keygen "$@"
}

###############################################################################
# BTC Infrastructure Setup Functions
###############################################################################

# Check BTC prerequisites
# Usage: btc_check_prerequisites [command_list]
# Example: btc_check_prerequisites "watch keygen"
#          btc_check_prerequisites "watch keygen sign1 sign2"
btc_check_prerequisites() {
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

# Setup BTC infrastructure
# Usage: btc_setup_infrastructure [container_list]
# Example: btc_setup_infrastructure "btc-watch btc-keygen"
#          btc_setup_infrastructure "btc-watch btc-keygen btc-sign1 btc-sign2"
btc_setup_infrastructure() {
	local containers="${1:-btc-watch btc-keygen}"
	log_step "Setting up infrastructure"

	if db_is_sqlite; then
		# SQLite mode: Initialize SQLite database with all schemas
		log_substep "Initializing SQLite database"
		sqlite_init_db "all"
		log_info "SQLite database initialized"
	else
		# MySQL mode: Start database container
		log_substep "Starting database container"
		docker compose -f compose.yaml up -d
		log_info "Database container started"

		# Wait for database to be healthy
		wait_for_healthy "wallet-db" 90
	fi

	# Start Bitcoin nodes
	log_substep "Starting Bitcoin node containers"
	# shellcheck disable=SC2086
	docker compose -f compose.btc.yaml up -d $containers
	log_info "Bitcoin node containers started"

	# Wait for containers to be healthy
	log_substep "Waiting for containers to be healthy"
	for container in $containers; do
		wait_for_healthy "$container" 90
	done

	log_info "All containers are healthy"
}

# Setup BTC wallets
# Usage: btc_setup_wallets [wallet_names]
# Example: btc_setup_wallets "watch keygen"
#          btc_setup_wallets "watch keygen sign1 sign2"
# Note: Container name is derived as "btc-${wallet_name}"
btc_setup_wallets() {
	local wallets="${1:-watch keygen}"
	log_step "Setting up Bitcoin wallets"

	for wallet in $wallets; do
		local container="btc-${wallet}"
		btc_create_wallet_if_needed "$container" "$wallet"
	done

	log_info "All wallets are ready"
}

###############################################################################
# BTC Transaction Helper Functions
###############################################################################

# Log detailed error message for "No utxo" errors
# Usage: btc_log_no_utxo_error
btc_log_no_utxo_error() {
	log_error "Transaction creation failed"
	log_error "This could indicate:"
	log_error "  - No payment requests in database"
	log_error "  - No UTXOs available for payment account"
	log_error "  - UTXOs not mature enough (need 100+ confirmations)"
}

# Wait for balance to update after UTXO generation
# Usage: btc_wait_for_balance [max_wait] [wait_interval]
btc_wait_for_balance() {
	local max_wait="${1:-60}"
	local wait_interval="${2:-3}"
	local elapsed=0
	local balance_found=false

	log_info "Waiting for blockchain sync and balance update..."

	while [ $elapsed -lt $max_wait ]; do
		# Check balance using Bitcoin Core RPC directly
		local balance_json
		balance_json=$(btc_cli "btc-watch" -rpcwallet=watch getbalances 2>&1 || true)
		local trusted_balance
		trusted_balance=$(echo "$balance_json" | jq -r '.mine.trusted // 0' 2>/dev/null || echo "0")

		# Check if we have any trusted (mature) balance
		if [ -n "$trusted_balance" ] && [ "$(echo "$trusted_balance > 0" | bc -l 2>/dev/null || echo 0)" -eq 1 ]; then
			log_info "Payment account balance verified: ${trusted_balance} BTC (took ${elapsed}s)"
			balance_found=true
			break
		fi

		sleep "$wait_interval"
		elapsed=$((elapsed + wait_interval))
		if [ $elapsed -lt $max_wait ]; then
			log_info "Still waiting for balance update... (${elapsed}s/${max_wait}s)"
		fi
	done

	if [ "$balance_found" = false ]; then
		log_error "Balance not detected within ${max_wait}s"
		log_error "This indicates a failure in UTXO generation or blockchain sync"
		log_error "Please check:"
		log_error "  - Bitcoin node logs: docker compose -f compose.btc.yaml logs btc-watch"
		log_error "  - Block generation succeeded"
		log_error "  - Address import into watch wallet succeeded"
		return 1
	fi

	return 0
}

# Generate test UTXOs for a payment address
# Usage: btc_generate_test_utxos "payment_address" [block_count]
btc_generate_test_utxos() {
	local payment_address="$1"
	local block_count="${2:-101}"

	log_info "Using payment address: $payment_address"

	# Generate blocks with coinbase reward to payment address
	log_info "Generating $block_count blocks to create mature coinbase for testing..."
	btc_cli "btc-watch" generatetoaddress "$block_count" "$payment_address" >/dev/null

	log_info "Test UTXOs generated successfully"
}

# Derive payment address from descriptor
# Usage: payment_address=$(btc_derive_address_from_descriptor "$descriptor")
btc_derive_address_from_descriptor() {
	local descriptor="$1"
	local descriptor_with_checksum

	log_info "Deriving payment address from descriptor..." >&2

	# Add checksum to descriptor if not present
	if ! echo "$descriptor" | grep -q '#'; then
		log_info "Adding checksum to descriptor..." >&2
		descriptor_with_checksum=$(btc_cli "btc-watch" getdescriptorinfo "$descriptor" | jq -r '.descriptor')
	else
		descriptor_with_checksum="$descriptor"
	fi

	# Derive first address (index 0) from descriptor
	local payment_address
	payment_address=$(btc_cli "btc-watch" deriveaddresses "$descriptor_with_checksum" "[0,0]" 2>/dev/null | jq -r '.[0]')

	if [ -z "$payment_address" ] || [ "$payment_address" = "null" ]; then
		log_error "Failed to derive payment address from descriptor" >&2
		log_error "Descriptor: $descriptor_with_checksum" >&2
		return 1
	fi

	echo "$payment_address"
}

# Extract file path from command output
# Expects format: "[fileName]: /path/to/file"
# Usage: file_path=$(btc_extract_file_path "$output")
btc_extract_file_path() {
	local output="$1"
	echo "$output" | sed -n 's/.*\[fileName\]: //p'
}

# Extract descriptor export path from command output
# Expects format: "exported to /path/to/file"
# Usage: file_path=$(btc_extract_descriptor_path "$output")
btc_extract_descriptor_path() {
	local output="$1"
	echo "${output##*exported to }"
}

###############################################################################
# BTC Payment Request Functions
###############################################################################

# Get sender address from database
# Usage: sender_address=$(btc_get_sender_address "payment")
btc_get_sender_address() {
	local account="${1:-payment}"

	db_query "watch" "SELECT wallet_address FROM address WHERE coin='btc' AND account='${account}' LIMIT 1"
}

# Generate receiver addresses
# Usage: receivers=$(btc_generate_receiver_addresses 3 "legacy")
#        receivers=$(btc_generate_receiver_addresses 3 "bech32m")
btc_generate_receiver_addresses() {
	local count="${1:-3}"
	local address_type="${2:-legacy}"
	local addresses=""

	log_substep "Generating receiver addresses for payment requests" >&2
	log_info "Creating $count receiver addresses in watch wallet..." >&2

	for i in $(seq 1 "$count"); do
		local addr
		addr=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" "$address_type")
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
# Usage: btc_insert_payment_requests "sender_address" "receiver1 receiver2 receiver3" "0.001 0.002 0.0015"
btc_insert_payment_requests() {
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
	('btc', NULL, '${sender_address}', '${account}', '${receiver_array[$i]}', '${amount_array[$i]}', 0)"
		else
			values="$values
	('btc', NULL, '${sender_address}', '${account}', '${receiver_array[$i]}', ${amount_array[$i]}, false)"
		fi
	done

	local sql_statements
	sql_statements="DELETE FROM payment_request WHERE coin='btc';
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
		count=$(db_query "watch" "SELECT COUNT(*) FROM payment_request WHERE coin='btc' AND is_done=0")
	else
		count=$(db_query "watch" "SELECT COUNT(*) FROM payment_request WHERE coin='btc' AND is_done=false")
	fi

	log_info "Created $count payment requests"

	if [ "$count" -eq 0 ]; then
		log_error "No payment requests were created"
		return 1
	fi

	log_info "Payment requests ready for transaction creation"
}

###############################################################################
# BTC Argument Parsing Helper
###############################################################################

# Parse common e2e script arguments
# Usage: btc_parse_args "$@"
#        This sets: CLEANUP_ONLY, RESET_STATE, VERBOSE, NON_INTERACTIVE
btc_parse_args() {
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
