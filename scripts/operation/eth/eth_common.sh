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

# Ensure GOPATH is set — CI runners may not have it as an environment variable
GOPATH="${GOPATH:-$(go env GOPATH)}"
export GOPATH

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
# SQLite Schema Initialization
###############################################################################

# Initialize SQLite databases for E2E testing using the composite schemas.
# Must be called after eth_init_database() sets SQLITE_*_DB_PATH.
eth_init_sqlite_db() {
	if [ "${DB_TYPE}" != "sqlite" ]; then
		return 0
	fi

	local project_root
	project_root="$(cd "${_ETH_COMMON_DIR}/../../.." && pwd)"

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

	# Skip container management when using shared infrastructure (parallel E2E)
	if [ "${E2E_SHARED_INFRASTRUCTURE:-false}" = "true" ]; then
		log_info "Using shared infrastructure (skipping Docker container management)"
		log_substep "Initializing SQLite databases for pattern: ${E2E_PATTERN}"
		eth_init_sqlite_db
		return 0
	fi

	case "${NODE_TYPE}" in
	geth)
		# Use geth-dev profile: local PoA dev chain (geth --dev), no sync required.
		# The testnet geth profile requires Sepolia sync (hours) and is not for E2E testing.
		log_substep "Starting Geth dev node (geth --dev, chain ID 1337)"
		docker compose -f compose.eth.yaml --profile geth-dev up -d
		;;
	anvil | *)
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

	# Skip container management when using shared infrastructure (parallel E2E)
	if [ "${E2E_SHARED_INFRASTRUCTURE:-false}" = "true" ]; then
		log_info "Using shared infrastructure (skipping cleanup)"
		return 0
	fi

	case "${NODE_TYPE}" in
	geth)
		log_substep "Stopping Geth dev node"
		docker compose -f compose.eth.yaml --profile geth-dev down || true
		;;
	anvil | *)
		log_substep "Stopping Anvil"
		docker compose -f compose.eth.yaml --profile anvil down anvil || true
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

	# Skip container management when using shared infrastructure (parallel E2E)
	if [ "${E2E_SHARED_INFRASTRUCTURE:-false}" = "true" ]; then
		log_info "Using shared infrastructure (skipping full reset)"
		if [ "${DB_TYPE}" = "sqlite" ]; then
			log_info "Cleaning SQLite databases for pattern: ${E2E_PATTERN}"
			rm -f "${PROJECT_ROOT:-./}/data/sqlite/eth/watch-e2e-${E2E_PATTERN}.db" \
				"${PROJECT_ROOT:-./}/data/sqlite/eth/keygen-e2e-${E2E_PATTERN}.db" 2>/dev/null || true
			eth_init_sqlite_db
		fi
		return 0
	fi

	eth_cleanup

	log_substep "Removing keystore files"
	rm -f ./data/keystore/*

	log_info "Full reset completed"
}

###############################################################################
# Database Query Functions (SQLite)
###############################################################################

# Get the first unallocated address for an account from the watch wallet DB
# Usage: addr=$(eth_get_payment_address "payment")
eth_get_payment_address() {
	local account="${1:-payment}"

	if [ "${DB_TYPE}" = "sqlite" ]; then
		sqlite3 "${SQLITE_WATCH_DB_PATH}" \
			"SELECT wallet_address FROM address WHERE coin='eth' AND account='${account}' AND is_allocated=0 LIMIT 1" 2>/dev/null || true
	else
		log_warn "eth_get_payment_address: MySQL support not yet implemented"
		echo ""
	fi
}

# Export ETH addresses from keygen DB as a CSV file compatible with watch import address.
# The watch's import address use case uses ConvertAddressLine which expects 8 fields:
#   coinTypeCode,accountType,P2PKHAddress,P2SHSegwitAddress,Bech32Address,FullPublicKey,MultisigAddress,Idx
# For ETH, the ETH address goes in the P2PKHAddress slot (field 2), others are empty.
#
# Usage: csv_file=$(eth_export_watch_address_csv "payment")
eth_export_watch_address_csv() {
	local account="$1"
	local project_root
	project_root="$(cd "${_ETH_COMMON_DIR}/../../.." && pwd)"
	local addr_dir="${project_root}/data/address/eth"
	mkdir -p "${addr_dir}"

	local tmp_file
	tmp_file="${addr_dir}/${account}_$(date +%s%N 2>/dev/null || date +%s).csv"

	if [ "${DB_TYPE}" = "sqlite" ]; then
		# Query keygen DB for ETH addresses, emit 8-field CSV (P2PKHAddress=ETH address, rest empty)
		sqlite3 "${SQLITE_KEYGEN_DB_PATH}" \
			"SELECT 'eth'||','||account||','||address||',,,,,'||idx FROM eth_account_key WHERE account='${account}'" \
			>"${tmp_file}" 2>/dev/null || true
	else
		log_warn "eth_export_watch_address_csv: MySQL support not yet implemented"
	fi

	echo "${tmp_file}"
}

###############################################################################
# Utility Functions
###############################################################################

# Fund an address with ETH.
# - Anvil: uses anvil_setBalance (instant, no gas needed)
# - Geth dev: sends ETH from the pre-funded coinbase account via eth_sendTransaction
eth_fund_address() {
	local address="$1"
	local amount_eth="${2:-100}" # Default 100 ETH

	# Convert ETH to Wei using bc for precision (1 ETH = 10^18 Wei)
	# Use bc for both multiplication and hex conversion: values exceed 2^63 (e.g. 100 ETH = 10^20 Wei)
	local amount_wei_dec
	amount_wei_dec=$(echo "${amount_eth} * 1000000000000000000" | bc)

	local amount_wei_hex
	amount_wei_hex="0x$(echo "obase=16; ${amount_wei_dec}" | bc)"

	log_substep "Funding ${address} with ${amount_eth} ETH (${amount_wei_hex} Wei)"

	if [ "${NODE_TYPE}" = "geth" ]; then
		# geth --dev auto-unlocks the first account (developer account).
		# Use eth_accounts (eth_coinbase is not available in recent geth versions).
		local dev_account
		dev_account=$(curl -s -X POST -H "Content-Type: application/json" \
			--data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' \
			"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" |
			grep -oE '"0x[a-fA-F0-9]{40}"' | head -1 | tr -d '"' || true)

		if [ -z "${dev_account}" ]; then
			log_error "eth_fund_address: failed to get dev account from geth dev node"
			return 1
		fi

		curl -s -X POST -H "Content-Type: application/json" \
			--data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendTransaction\",\"params\":[{\"from\":\"${dev_account}\",\"to\":\"${address}\",\"value\":\"${amount_wei_hex}\",\"gas\":\"0x5208\"}],\"id\":1}" \
			"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" >/dev/null

		# Wait for the transaction to be mined (geth --dev mines on each tx submission,
		# but there may be a brief delay before the balance is reflected).
		local wait_count=0
		local max_wait_balance=10
		while [ "${wait_count}" -lt "${max_wait_balance}" ]; do
			local balance
			balance=$(curl -s -X POST -H "Content-Type: application/json" \
				--data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"${address}\",\"latest\"],\"id\":1}" \
				"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" |
				grep -oE '"result":"0x[a-fA-F0-9]+"' | sed 's/"result":"//;s/"//' || true)
			if [ -n "${balance}" ] && [ "${balance}" != "0x0" ]; then
				break
			fi
			sleep 1
			wait_count=$((wait_count + 1))
		done
	else
		curl -s -X POST -H "Content-Type: application/json" \
			--data "{\"jsonrpc\":\"2.0\",\"method\":\"anvil_setBalance\",\"params\":[\"${address}\",\"${amount_wei_hex}\"],\"id\":1}" \
			"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" >/dev/null
	fi

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
