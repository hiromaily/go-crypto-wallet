#!/usr/bin/env bash

# Ethereum-Specific Common Utility Functions for E2E Operations
# This file provides shared utility functions specifically for Ethereum E2E scripts.
#
# Usage:
#   source "$(dirname "${BASH_SOURCE[0]}")/eth_common.sh"
#
# Note: This file sources common.sh automatically, so you don't need to source both.

# Script directory for relative paths
_ETH_COMMON_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the general common utilities
# shellcheck source=../common.sh
source "${_ETH_COMMON_DIR}/../common.sh"

###############################################################################
# ETH Default Configuration
###############################################################################

# Coin identifier
ETH_COIN="eth"

# Node type: anvil (default) or geth
# Override via NODE_TYPE=geth environment variable
NODE_TYPE="${NODE_TYPE:-anvil}"

# RPC configuration (Anvil defaults to port 8546 to avoid conflict with Geth)
ETH_RPC_HOST="${ETH_RPC_HOST:-127.0.0.1}"
if [ "${NODE_TYPE}" = "geth" ]; then
	ETH_RPC_PORT="${ETH_RPC_PORT:-8545}"
else
	ETH_RPC_PORT="${ETH_RPC_PORT:-8546}"
fi

# MySQL credentials (can be overridden via environment variables)
ETH_MYSQL_ROOT_PASSWORD="${ETH_MYSQL_ROOT_PASSWORD:-${MYSQL_ROOT_PASSWORD:-root}}"

# Docker volume name
ETH_DOCKER_VOLUME_NAME="${ETH_DOCKER_VOLUME_NAME:-go-crypto-wallet_eth-mysql}"

###############################################################################
# Database Configuration
###############################################################################

# Database type: sqlite (default) or mysql
DB_TYPE="${DB_TYPE:-sqlite}"

# E2E Pattern identifier (set by each E2E script)
E2E_PATTERN="${E2E_PATTERN:-}"

# Get database file path with pattern suffix for a given wallet type
eth_get_db_path() {
	local wallet_type="${1:-watch}"
	local db_suffix=""
	if [ -n "${E2E_PATTERN}" ]; then
		db_suffix="-e2e-${E2E_PATTERN}"
	else
		db_suffix="-e2e"
	fi
	echo "./data/sqlite/eth/${wallet_type}${db_suffix}.db"
}

# Separate DB paths for each wallet type (mirrors BTC pattern)
SQLITE_WATCH_DB_PATH=""
SQLITE_KEYGEN_DB_PATH=""

# Initialize database configuration — sets SQLITE_*_DB_PATH and exports env vars
eth_init_database() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		export WALLET_DATABASE_TYPE="sqlite"

		SQLITE_WATCH_DB_PATH="$(eth_get_db_path "watch")"
		SQLITE_KEYGEN_DB_PATH="$(eth_get_db_path "keygen")"

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

eth_get_config_paths() {
	local project_root
	project_root="$(cd "${_ETH_COMMON_DIR}/../../.." && pwd)"

	ETH_CONFIG_WATCH="${project_root}/config/wallet/eth/watch.yaml"
	ETH_CONFIG_KEYGEN="${project_root}/config/wallet/eth/keygen.yaml"
	ETH_CONFIG_SIGN="${project_root}/config/wallet/eth/sign.yaml"

	export ETH_CONFIG_WATCH ETH_CONFIG_KEYGEN ETH_CONFIG_SIGN
}

###############################################################################
# Wallet Command Wrappers
# Each wrapper injects the wallet-specific SQLite DB path so that watch and
# keygen operate on isolated databases (mirrors BTC btc_watch_cmd pattern).
###############################################################################

eth_watch_cmd() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		WALLET_DATABASE_SQLITE_PATH="${SQLITE_WATCH_DB_PATH}" \
			"${GOPATH}/bin/watch" "$@"
	else
		"${GOPATH}/bin/watch" "$@"
	fi
}

eth_keygen_cmd() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		WALLET_DATABASE_SQLITE_PATH="${SQLITE_KEYGEN_DB_PATH}" \
			"${GOPATH}/bin/keygen" "$@"
	else
		"${GOPATH}/bin/keygen" "$@"
	fi
}

eth_sign_cmd() {
	"${GOPATH}/bin/sign" "$@"
}

###############################################################################
# Infrastructure Setup
###############################################################################

eth_check_prerequisites() {
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

	log_info "All prerequisites met"
}

eth_setup_infrastructure() {
	log_step "Setting Up Infrastructure (NODE_TYPE=${NODE_TYPE})"

	case "${NODE_TYPE}" in
	geth)
		compose_profile="geth"
		log_substep "Starting Geth + Lodestar nodes"
		docker compose -f compose.eth.yaml --profile geth up -d
		;;
	anvil | *)
		compose_profile="anvil"
		log_substep "Starting Anvil node"
		docker compose -f compose.eth.yaml --profile anvil up -d anvil
		;;
	esac

	# Wait for RPC endpoint to be ready
	log_substep "Waiting for ETH node (${NODE_TYPE}) on ${ETH_RPC_HOST}:${ETH_RPC_PORT}..."
	local max_wait=60
	local count=0
	while [ $count -lt $max_wait ]; do
		if curl -s -X POST -H "Content-Type: application/json" \
			--data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
			"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" >/dev/null 2>&1; then
			log_info "ETH node is ready"
			return 0
		fi
		sleep 1
		count=$((count + 1))
	done

	log_error "ETH node (${NODE_TYPE}) failed to start within ${max_wait} seconds"
	return 1
}

###############################################################################
# Cleanup Functions
###############################################################################

eth_cleanup() {
	log_step "Cleaning Up"

	case "${NODE_TYPE}" in
	geth)
		log_substep "Stopping Geth + Lodestar"
		docker compose -f compose.eth.yaml --profile geth stop || true
		;;
	anvil | *)
		log_substep "Stopping Anvil"
		docker compose -f compose.eth.yaml --profile anvil stop anvil || true
		;;
	esac

	if [ "${DB_TYPE}" = "sqlite" ]; then
		log_substep "Removing SQLite databases"
		rm -f ./data/sqlite/eth/*-e2e*.db
	fi

	log_info "Cleanup completed"
}

eth_full_reset() {
	log_step "Performing Full Reset"

	eth_cleanup

	log_substep "Removing keystore files"
	rm -f ./data/keystore/*

	log_info "Full reset completed"
}

###############################################################################
# Database Query Functions (SQLite)
###############################################################################

# Get the first payment address from the watch wallet DB
# Usage: addr=$(eth_get_payment_address "payment")
eth_get_payment_address() {
	local account="${1:-payment}"

	if [ "${DB_TYPE}" = "sqlite" ]; then
		sqlite3 "${SQLITE_WATCH_DB_PATH}" \
			"SELECT wallet_address FROM account_pubkey_table WHERE account='${account}' AND is_allocated=0 LIMIT 1" 2>/dev/null || true
	else
		log_warn "eth_get_payment_address: MySQL support not yet implemented"
		echo ""
	fi
}

###############################################################################
# Utility Functions
###############################################################################

# Fund an address using anvil_setBalance RPC (Anvil only)
eth_fund_address() {
	local address="$1"
	local amount_eth="${2:-100}" # Default 100 ETH

	if [ "${NODE_TYPE}" = "geth" ]; then
		log_warn "eth_fund_address: auto-funding not supported on Geth; fund address manually"
		return 0
	fi

	# Convert ETH to Wei using bc for precision (1 ETH = 10^18 Wei)
	# Use bc for both multiplication and hex conversion: values exceed 2^63 (e.g. 100 ETH = 10^20 Wei)
	local amount_wei_dec
	amount_wei_dec=$(echo "${amount_eth} * 1000000000000000000" | bc)

	local amount_wei_hex
	amount_wei_hex="0x$(echo "obase=16; ${amount_wei_dec}" | bc)"

	log_substep "Funding ${address} with ${amount_eth} ETH (${amount_wei_hex} Wei)"

	curl -s -X POST -H "Content-Type: application/json" \
		--data "{\"jsonrpc\":\"2.0\",\"method\":\"anvil_setBalance\",\"params\":[\"${address}\",\"${amount_wei_hex}\"],\"id\":1}" \
		"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" >/dev/null

	log_info "Address funded successfully"
}

# Extract a file path tagged with [fileName]: from command output
# Usage: file=$(eth_extract_file_path "$output")
eth_extract_file_path() {
	local output="$1"
	echo "$output" | grep '^\[fileName\]:' | sed 's/^\[fileName\]: //'
}

# Extract txID from watch send tx output
# Usage: txid=$(eth_extract_tx_id "$output")
eth_extract_tx_id() {
	local output="$1"
	echo "$output" | grep -oE '0x[a-fA-F0-9]{64}' | head -1
}

###############################################################################
# Initialization
###############################################################################

# Initialize database configuration (runs on source)
eth_init_database

log_info "ETH common utilities loaded (coin=${ETH_COIN}, node_type=${NODE_TYPE}, db_type=${DB_TYPE})"
