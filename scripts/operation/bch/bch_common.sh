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
BCH_WATCH_WALLET_RPC_HOST="${BCH_WATCH_WALLET_RPC_HOST:-127.0.0.1:28332/wallet/watch}"
BCH_KEYGEN_WALLET_RPC_HOST="${BCH_KEYGEN_WALLET_RPC_HOST:-127.0.0.1:29332/wallet/keygen}"

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

	BCH_CONFIG_WATCH="${PROJECT_ROOT}/config/wallet/bch_watch.yaml"
	BCH_CONFIG_KEYGEN="${PROJECT_ROOT}/config/wallet/bch_keygen.yaml"
	BCH_CONFIG_SIGN1="${PROJECT_ROOT}/config/wallet/bch_sign1.yaml"
	BCH_CONFIG_SIGN2="${PROJECT_ROOT}/config/wallet/bch_sign2.yaml"
	BCH_CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account.yaml"
	BCH_CONFIG_ACCOUNT_2OF3="${PROJECT_ROOT}/config/wallet/account_2of3.yaml"
	BCH_CONFIG_ACCOUNT_3OF3="${PROJECT_ROOT}/config/wallet/account_3of3.yaml"
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
bch_full_reset() {
	local volume_name="${1:-${BCH_DOCKER_VOLUME_NAME}}"
	local nodes="${2:-watch keygen}"
	log_step "Performing full reset for fresh state"

	# Stop and remove containers WITH VOLUMES
	log_info "Stopping Bitcoin Cash containers (with volume removal)..."
	docker compose -f compose.bch.yaml down -v 2>/dev/null || true

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

	# Clean data files
	bch_clean_data_files

	# Clean Bitcoin Cash wallet data
	bch_clean_wallet_data "$nodes"

	log_info "Full reset complete - system is in fresh state"
	log_info "Note: Database volumes were removed for complete cleanup"
}

# Basic cleanup: just stop containers
# Usage: bch_cleanup
bch_cleanup() {
	log_step "Cleaning up containers and state"

	log_info "Stopping Bitcoin Cash containers..."
	docker compose -f compose.bch.yaml down -v 2>/dev/null || true

	log_info "Stopping database container..."
	docker compose -f compose.yaml down -v 2>/dev/null || true

	log_info "Cleanup complete"
}

###############################################################################
# BCH Wallet Command Wrappers
###############################################################################

# Wrapper for watch wallet commands with host override
# Usage: bch_watch_cmd [args...]
bch_watch_cmd() {
	WALLET_BITCOIN_HOST="${BCH_WATCH_WALLET_RPC_HOST}" watch "$@"
}

# Wrapper for keygen wallet commands with host override
# Usage: bch_keygen_cmd [args...]
bch_keygen_cmd() {
	WALLET_BITCOIN_HOST="${BCH_KEYGEN_WALLET_RPC_HOST}" keygen "$@"
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
bch_setup_infrastructure() {
	local containers="${1:-bch-watch bch-keygen}"
	log_step "Setting up infrastructure"

	# Start database
	log_substep "Starting database container"
	docker compose -f compose.yaml up -d
	log_info "Database container started"

	# Wait for database to be healthy
	wait_for_healthy "wallet-db" 90

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
bch_setup_wallets() {
	local wallets="${1:-watch keygen}"
	log_step "Setting up Bitcoin Cash wallets"

	for wallet in $wallets; do
		local container="bch-${wallet}"
		bch_create_wallet_if_needed "$container" "$wallet"
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
		# Check balance using Bitcoin Cash RPC directly
		local balance_json
		balance_json=$(bch_cli "bch-watch" -rpcwallet=watch getbalances 2>&1 || true)
		local trusted_balance
		trusted_balance=$(echo "$balance_json" | jq -r '.mine.trusted // 0' 2>/dev/null || echo "0")

		# Check if we have any trusted (mature) balance
		if [ -n "$trusted_balance" ] && [ "$(echo "$trusted_balance > 0" | bc -l 2>/dev/null || echo 0)" -eq 1 ]; then
			log_info "Payment account balance verified: ${trusted_balance} BCH (took ${elapsed}s)"
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
	log_info "Generating $block_count blocks to create mature coinbase for testing..."
	bch_cli "bch-watch" generatetoaddress "$block_count" "$payment_address" >/dev/null

	log_info "Test UTXOs generated successfully"
}

# Extract file path from command output
# Expects format: "[fileName]: /path/to/file"
# Usage: file_path=$(bch_extract_file_path "$output")
bch_extract_file_path() {
	local output="$1"
	echo "${output##*\[fileName\]: }"
}

###############################################################################
# BCH Payment Request Functions
###############################################################################

# Get sender address from database
# Usage: sender_address=$(bch_get_sender_address "payment")
bch_get_sender_address() {
	local account="${1:-payment}"

	docker compose exec -e MYSQL_PWD="${BCH_MYSQL_ROOT_PASSWORD}" -T wallet-db mysql -u root watch -N -e \
		"SELECT wallet_address FROM address WHERE coin='bch' AND account='${account}' LIMIT 1" 2>/dev/null
}

# Generate receiver addresses (using legacy format for BCH regtest)
# Usage: receivers=$(bch_generate_receiver_addresses 3)
bch_generate_receiver_addresses() {
	local count="${1:-3}"
	local addresses=""

	log_substep "Generating receiver addresses for payment requests" >&2
	log_info "Creating $count receiver addresses in watch wallet..." >&2

	for i in $(seq 1 "$count"); do
		local addr
		# BCH uses legacy address format (no SegWit)
		addr=$(bch_cli "bch-watch" -rpcwallet=watch getnewaddress "" legacy)
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
		values="$values
	('bch', NULL, '${sender_address}', '${account}', '${receiver_array[$i]}', ${amount_array[$i]}, false)"
	done

	if ! docker compose exec -e MYSQL_PWD="${BCH_MYSQL_ROOT_PASSWORD}" -T wallet-db mysql -u root watch <<EOF; then
DELETE FROM payment_request;
INSERT INTO payment_request (coin, payment_id, sender_address, sender_account, receiver_address, amount, is_done)
VALUES${values};
EOF
		log_error "Failed to insert payment requests"
		return 1
	fi

	# Verify payment requests were created
	local count
	count=$(docker compose exec -e MYSQL_PWD="${BCH_MYSQL_ROOT_PASSWORD}" -T wallet-db mysql -u root watch -N -e \
		"SELECT COUNT(*) FROM payment_request WHERE coin='bch' AND is_done=false" 2>/dev/null)

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
