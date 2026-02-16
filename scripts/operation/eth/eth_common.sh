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

# RPC configuration for Anvil (default port 8546 to avoid conflict with Geth)
ETH_RPC_HOST="${ETH_RPC_HOST:-127.0.0.1}"
ETH_RPC_PORT="${ETH_RPC_PORT:-8546}"

# MySQL credentials (can be overridden via environment variables)
ETH_MYSQL_ROOT_PASSWORD="${ETH_MYSQL_ROOT_PASSWORD:-${MYSQL_ROOT_PASSWORD:-root}}"

# Docker volume name
ETH_DOCKER_VOLUME_NAME="${ETH_DOCKER_VOLUME_NAME:-go-crypto-wallet_wallet-mysql}"

###############################################################################
# Database Configuration
###############################################################################

# Database type: sqlite (default) or mysql
DB_TYPE="${DB_TYPE:-sqlite}"

# E2E Pattern identifier (set by each E2E script)
E2E_PATTERN="${E2E_PATTERN:-}"

# Get database file path with pattern suffix
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

# Initialize database configuration
eth_init_database() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		export WALLET_DATABASE_TYPE="sqlite"
		WALLET_DATABASE_SQLITE_PATH="$(eth_get_db_path "watch")"
		export WALLET_DATABASE_SQLITE_PATH
		log_info "Using SQLite database: ${WALLET_DATABASE_SQLITE_PATH}"

		# Create directory if it doesn't exist
		mkdir -p "$(dirname "${WALLET_DATABASE_SQLITE_PATH}")"
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
###############################################################################

eth_watch_cmd() {
	"${GOPATH}/bin/watch" "$@"
}

eth_keygen_cmd() {
	"${GOPATH}/bin/keygen" "$@"
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

	if [ ! -f "${GOPATH}/bin/sign" ]; then
		log_error "sign binary not found. Run 'make build-all' first."
		return 1
	fi

	log_info "All prerequisites met"
}

eth_setup_infrastructure() {
	log_step "Setting Up Infrastructure"

	# Start Anvil
	log_substep "Starting Anvil node"
	docker compose -f compose.eth.yaml up -d anvil

	# Wait for Anvil to be ready
	log_substep "Waiting for Anvil to be ready..."
	local max_wait=30
	local count=0
	while [ $count -lt $max_wait ]; do
		if curl -s -X POST -H "Content-Type: application/json" \
			--data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
			"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" >/dev/null 2>&1; then
			log_info "Anvil is ready"
			return 0
		fi
		sleep 1
		count=$((count + 1))
	done

	log_error "Anvil failed to start within ${max_wait} seconds"
	return 1
}

###############################################################################
# Cleanup Functions
###############################################################################

eth_cleanup() {
	log_step "Cleaning Up"

	log_substep "Stopping Anvil"
	docker compose -f compose.eth.yaml stop anvil

	log_substep "Removing SQLite databases"
	rm -f ./data/sqlite/eth/*-e2e*.db

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
# Utility Functions
###############################################################################

# Fund an address using anvil_setBalance RPC
eth_fund_address() {
	local address="$1"
	local amount_eth="${2:-100}" # Default 100 ETH

	# Convert ETH to Wei using bc for precision (1 ETH = 10^18 Wei)
	# Note: bash arithmetic can't handle values > 2^63-1, so we use bc
	local amount_wei_dec
	amount_wei_dec=$(echo "${amount_eth} * 1000000000000000000" | bc)

	# Convert decimal Wei to hex
	local amount_wei_hex
	amount_wei_hex=$(printf "0x%x" "${amount_wei_dec}")

	log_substep "Funding address ${address} with ${amount_eth} ETH"

	curl -s -X POST -H "Content-Type: application/json" \
		--data "{\"jsonrpc\":\"2.0\",\"method\":\"anvil_setBalance\",\"params\":[\"${address}\",\"${amount_wei_hex}\"],\"id\":1}" \
		"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" >/dev/null

	log_info "Address funded successfully"
}

# Extract file path from command output
eth_extract_file_path() {
	local output="$1"
	echo "$output" | grep -oE '(data/tx|\.)/[^ ]+\.(hex|json)' | tail -1
}

###############################################################################
# Initialization
###############################################################################

# Initialize database configuration
eth_init_database

log_info "ETH common utilities loaded (coin=${ETH_COIN}, db_type=${DB_TYPE})"
