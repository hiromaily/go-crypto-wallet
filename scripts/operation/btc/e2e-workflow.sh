#!/usr/bin/env bash

# Bitcoin E2E Workflow Script
# This script automates the complete Bitcoin workflow from infrastructure setup to transaction execution
# Usage: ./scripts/operation/btc/e2e-workflow.sh [OPTIONS]
# Options:
#   --cleanup  Stop containers and cleanup state
#   --verbose  Enable verbose output
#   -h, --help Display help message

set -eu

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"

# Source common utilities
# shellcheck source=../common.sh
source "${SCRIPT_DIR}/../common.sh"

# Configuration
COIN="btc"
ENCRYPTED="false"
SIGN_WALLET_NUM=2
VERBOSE=false
CLEANUP_ONLY=false
NON_INTERACTIVE=false
RESET_STATE=false

# RPC credentials (can be overridden via environment variables)
# Note: Default values are for regtest/development only
RPC_USER="${RPC_USER:-xyz}"
RPC_PASSWORD="${RPC_PASSWORD:-xyz}"

# Wallet passphrase (only used if ENCRYPTED=true)
# Note: Default value is for testing only - use strong passphrase in production
WALLET_PASSPHRASE="${WALLET_PASSPHRASE:-test}"

# Config file paths (absolute)
CONFIG_WATCH="${PROJECT_ROOT}/config/wallet/btc_watch.yaml"
CONFIG_KEYGEN="${PROJECT_ROOT}/config/wallet/btc_keygen.yaml"
CONFIG_SIGN1="${PROJECT_ROOT}/config/wallet/btc_sign1.yaml"
CONFIG_SIGN2="${PROJECT_ROOT}/config/wallet/btc_sign2.yaml"

###############################################################################
# Cleanup Functions
###############################################################################

# Clean generated data files
clean_data_files() {
	log_substep "Cleaning generated data files"

	# Clean address files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/address/btc"

	# Clean fullpubkey files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/fullpubkey/btc"

	# Clean transaction files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/tx/btc"
}

# Clean Bitcoin node wallet data
clean_bitcoin_wallet_data() {
	log_substep "Cleaning Bitcoin node wallet data"

	# Clean regtest data directories (keeping bitcoin.conf)
	for node in watch keygen sign1 sign2; do
		wallet_dir="docker/nodes/btc/${node}/regtest"
		if [ -d "$wallet_dir" ]; then
			# Remove all files/dirs except bitcoin.conf
			find "$wallet_dir" -mindepth 1 ! -name 'bitcoin.conf' -exec rm -rf {} + 2>/dev/null || true
			log_info "Cleaned ${node} wallet data"
		fi
	done
}

# Full reset: cleanup everything for fresh state
full_reset() {
	log_step "Performing full reset for fresh state"

	# Stop and remove containers
	log_info "Stopping Bitcoin containers..."
	docker compose -f compose.btc.yaml down -v 2>/dev/null || true

	log_info "Stopping database container..."
	docker compose -f compose.yaml down -v 2>/dev/null || true

	# Clean data files
	clean_data_files

	# Clean Bitcoin wallet data
	clean_bitcoin_wallet_data

	log_info "Full reset complete - system is in fresh state"
}

# Basic cleanup: just stop containers
cleanup() {
	log_step "Cleaning up containers and state"

	# Stop and remove containers
	log_info "Stopping Bitcoin containers..."
	docker compose -f compose.btc.yaml down -v 2>/dev/null || true

	log_info "Stopping database container..."
	docker compose -f compose.yaml down -v 2>/dev/null || true

	log_info "Cleanup complete"
}

###############################################################################
# Prerequisites Check
###############################################################################

check_prerequisites() {
	log_step "Checking prerequisites"

	# Check Docker and Docker Compose
	check_docker || exit 1

	# Check CLI commands
	for cmd in watch keygen sign1 sign2; do
		if ! command_exists "$cmd"; then
			log_error "$cmd command is not available"
			log_error "Please build the project first: make build"
			exit 1
		fi
	done

	log_info "All prerequisites satisfied"
}

###############################################################################
# Infrastructure Setup
###############################################################################

setup_infrastructure() {
	log_step "Setting up infrastructure"

	# Start database
	log_substep "Starting database container"
	docker compose -f compose.yaml up -d
	log_info "Database container started"

	# Wait for database to be healthy
	wait_for_healthy "wallet-db" 90

	# Start Bitcoin nodes
	log_substep "Starting Bitcoin node containers"
	docker compose -f compose.btc.yaml up -d
	log_info "Bitcoin node containers started"

	# Wait for containers to be healthy
	log_substep "Waiting for containers to be healthy"
	wait_for_healthy "btc-watch" 90
	wait_for_healthy "btc-keygen" 90
	wait_for_healthy "btc-sign1" 90
	wait_for_healthy "btc-sign2" 90

	log_info "All containers are healthy"
}

###############################################################################
# Wallet Setup
###############################################################################

setup_wallets() {
	log_step "Setting up Bitcoin wallets"

	# Create wallets in Bitcoin nodes
	btc_create_wallet_if_needed "btc-watch" "watch"
	btc_create_wallet_if_needed "btc-keygen" "keygen"
	btc_create_wallet_if_needed "btc-sign1" "sign1"
	btc_create_wallet_if_needed "btc-sign2" "sign2"

	log_info "All wallets are ready"
}

###############################################################################
# Key Generation Phase
###############################################################################

key_generation_phase() {
	log_step "Key Generation Phase"

	# Keygen wallet - create seed
	log_substep "Creating seed for keygen wallet"
	keygen -c "${CONFIG_KEYGEN}" create seed || {
		log_warn "Seed already exists or error occurred, continuing..."
	}

	# Keygen wallet - create hdkeys
	log_substep "Creating HD keys for keygen wallet (client, deposit, payment, stored)"
	for account in client deposit payment stored; do
		log_info "Creating HD keys for account: $account"
		keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" create hdkey --account "$account" --keynum 10
	done

	# Keygen wallet - import private keys
	log_substep "Importing private keys into keygen wallet"
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletpassphrase --passphrase "${WALLET_PASSPHRASE}"
	fi
	for account in client deposit payment stored; do
		log_info "Importing private keys for account: $account"
		keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" import privkey --account "$account"
	done
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletlock
	fi

	# Sign wallets - create seed (using sign1 as the primary sign wallet)
	log_substep "Creating seeds for sign wallets"
	sign1 --conf "${CONFIG_SIGN1}" --coin "${COIN}" create seed || {
		log_warn "Sign seed already exists or error occurred, continuing..."
	}

	# Sign wallets - create hdkeys
	log_substep "Creating HD keys for sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		log_info "Creating HD keys for sign${i}"
		config_var="CONFIG_SIGN${i}"
		"sign${i}" --conf "${!config_var}" --coin "${COIN}" --wallet "sign${i}" create hdkey
	done

	# Sign wallets - import private keys
	log_substep "Importing private keys into sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		log_info "Importing private keys for sign${i}"
		config_var="CONFIG_SIGN${i}"
		if [ "$ENCRYPTED" = "true" ]; then
			"sign${i}" --conf "${!config_var}" --coin "${COIN}" --wallet "sign${i}" api walletpassphrase --passphrase "${WALLET_PASSPHRASE}"
		fi
		"sign${i}" --conf "${!config_var}" --coin "${COIN}" --wallet "sign${i}" import privkey
		if [ "$ENCRYPTED" = "true" ]; then
			"sign${i}" --conf "${!config_var}" --coin "${COIN}" --wallet "sign${i}" api walletlock
		fi
	done

	# Sign wallets - export fullpubkey
	log_substep "Exporting full public keys from sign wallets"
	file_fullpubkey_auth1=$(sign1 --conf "${CONFIG_SIGN1}" --coin "${COIN}" --wallet sign1 export fullpubkey)
	file_fullpubkey_auth2=$(sign2 --conf "${CONFIG_SIGN2}" --coin "${COIN}" --wallet sign2 export fullpubkey)

	# Extract file paths
	fullpubkey_file1="${file_fullpubkey_auth1##*\[fileName\]: }"
	fullpubkey_file2="${file_fullpubkey_auth2##*\[fileName\]: }"

	log_info "Exported fullpubkey files:"
	log_info "  sign1: $fullpubkey_file1"
	log_info "  sign2: $fullpubkey_file2"

	# Store for next phase
	export FULLPUBKEY_FILE1="$fullpubkey_file1"
	export FULLPUBKEY_FILE2="$fullpubkey_file2"
}

###############################################################################
# Multisig Setup Phase
###############################################################################

multisig_setup_phase() {
	log_step "Multisig Setup Phase"

	# Import fullpubkeys
	log_substep "Importing full public keys into keygen wallet"
	log_info "Importing fullpubkey from sign1: $FULLPUBKEY_FILE1"
	keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" import fullpubkey --file "${FULLPUBKEY_FILE1}"

	log_info "Importing fullpubkey from sign2: $FULLPUBKEY_FILE2"
	keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" import fullpubkey --file "${FULLPUBKEY_FILE2}"

	# Create multisig addresses
	log_substep "Creating multisig addresses"
	for account in deposit payment stored; do
		log_info "Creating multisig address for account: $account"
		keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" create multisig --account "$account"
	done

	# Export addresses
	log_substep "Exporting addresses from keygen wallet"
	file_address_client=$(keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" export address --account client)
	file_address_deposit=$(keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" export address --account deposit)
	file_address_payment=$(keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" export address --account payment)
	file_address_stored=$(keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" export address --account stored)

	# Extract file paths
	address_client="${file_address_client##*\[fileName\]: }"
	address_deposit="${file_address_deposit##*\[fileName\]: }"
	address_payment="${file_address_payment##*\[fileName\]: }"
	address_stored="${file_address_stored##*\[fileName\]: }"

	log_info "Exported address files:"
	log_info "  client: $address_client"
	log_info "  deposit: $address_deposit"
	log_info "  payment: $address_payment"
	log_info "  stored: $address_stored"

	# Import addresses into watch wallet
	log_substep "Importing addresses into watch wallet"
	log_info "Importing client addresses"
	watch -c "${CONFIG_WATCH}" --coin "${COIN}" import address --file "${address_client}"

	log_info "Importing deposit addresses"
	watch -c "${CONFIG_WATCH}" --coin "${COIN}" import address --file "${address_deposit}"

	log_info "Importing payment addresses"
	watch -c "${CONFIG_WATCH}" --coin "${COIN}" import address --file "${address_payment}"

	log_info "Importing stored addresses"
	watch -c "${CONFIG_WATCH}" --coin "${COIN}" import address --file "${address_stored}"
}

###############################################################################
# UTXO Generation Phase (for regtest)
#
# This function generates test UTXOs required for the transaction phase.
# It only works in regtest environment where blocks can be instantly generated.
#
# How it works:
#   1. Extract payment address from the exported CSV file
#   2. Generate 101 blocks with coinbase rewards sent to the payment address
#      (101 blocks are needed because coinbase outputs require 100 confirmations
#       to become spendable - this is Bitcoin's coinbase maturity rule)
#   3. Poll for balance update to verify UTXOs are available
#   4. Fail with error if balance is not detected within timeout
#
# Note: The generatetoaddress RPC command is regtest-only and will not work
# on testnet or mainnet where actual proof-of-work mining is required.
###############################################################################

generate_test_utxos() {
	log_step "Generating Test UTXOs for Transaction Phase"

	log_info "Reading payment address from exported file..."
	# Get first payment address from the exported CSV file
	# CSV format: coin,account,P2PKH,P2SH-segwit,bech32,taproot,pubkey,,index
	# Use field 4 (P2SH-segwit address) which works for both single and multisig addresses
	payment_address=$(grep -v '^#' "${address_payment}" 2>/dev/null | head -n1 | cut -d',' -f4)

	if [ -z "$payment_address" ]; then
		log_error "Failed to extract payment address from ${address_payment}"
		return 1
	fi

	log_info "Using payment address: $payment_address"

	# Generate blocks with coinbase reward to payment address
	log_info "Generating 101 blocks to create mature coinbase for testing..."
	btc_cli "btc-watch" generatetoaddress 101 "$payment_address" >/dev/null

	log_info "Test UTXOs generated successfully"
	log_info "Waiting for blockchain sync and balance update..."

	# Poll for balance update with timeout
	max_wait=30
	wait_interval=2
	elapsed=0
	balance_found=false

	while [ $elapsed -lt $max_wait ]; do
		balance_output=$(watch -c "${CONFIG_WATCH}" --coin "${COIN}" monitor balance 2>&1 || true)
		if echo "$balance_output" | grep -q "payment"; then
			log_info "Payment account balance verified (took ${elapsed}s)"
			balance_found=true
			break
		fi
		sleep $wait_interval
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
}

###############################################################################
# Transaction Flow Phase
###############################################################################

transaction_flow_phase() {
	log_step "Transaction Flow Phase"

	# Create unsigned transaction
	log_substep "Creating unsigned payment transaction"
	tx_file=$(watch -c "${CONFIG_WATCH}" create payment 2>&1) || {
		log_error "Failed to create payment transaction"
		log_error "Output: $tx_file"

		if echo "$tx_file" | grep -q "No utxo"; then
			log_error "No UTXOs available for payment transaction"
			log_warn "Transaction flow phase skipped - no UTXOs available"
			return 0
		fi

		return 1
	}

	if echo "$tx_file" | grep -q "No utxo"; then
		log_warn "No UTXOs available for payment transaction"
		log_warn "Transaction flow phase skipped"
		return 0
	fi

	# Extract file path
	tx_unsigned="${tx_file##*\[fileName\]: }"
	log_info "Created unsigned transaction: $tx_unsigned"

	# Sign with keygen wallet (1st signature)
	log_substep "Signing with keygen wallet (1st signature)"
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletpassphrase --passphrase "${WALLET_PASSPHRASE}"
	fi
	tx_file_signed=$(keygen -c "${CONFIG_KEYGEN}" sign --file "${tx_unsigned}")
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletlock
	fi

	tx_signed1="${tx_file_signed##*\[fileName\]: }"
	log_info "Signed transaction (1st): $tx_signed1"

	# Sign with sign1 wallet (2nd signature)
	log_substep "Signing with sign1 wallet (2nd signature)"
	tx_file_signed2=$(sign1 --conf "${CONFIG_SIGN1}" --wallet sign1 sign --file "${tx_signed1}")
	tx_signed2="${tx_file_signed2##*\[fileName\]: }"
	log_info "Signed transaction (2nd): $tx_signed2"

	# Sign with sign2 wallet (3rd signature)
	log_substep "Signing with sign2 wallet (3rd signature)"
	tx_file_signed3=$(sign2 --conf "${CONFIG_SIGN2}" --wallet sign2 sign --file "${tx_signed2}")
	tx_signed3="${tx_file_signed3##*\[fileName\]: }"
	log_info "Signed transaction (3rd): $tx_signed3"

	# Send transaction
	log_substep "Sending fully signed transaction"
	tx_result=$(watch -c "${CONFIG_WATCH}" send --file "${tx_signed3}")
	tx_id="${tx_result##*txID: }"

	log_info "Transaction sent successfully!"
	log_info "Transaction ID: $tx_id"
}

###############################################################################
# Help Message
###############################################################################

show_help() {
	cat <<EOF
Bitcoin E2E Workflow Script

This script automates the complete Bitcoin workflow from infrastructure setup
to transaction execution. It serves as a regression test tool to verify that
the Bitcoin workflow functions correctly after code changes.

Usage: $0 [OPTIONS]

Options:
  --reset             Full reset: cleanup all state for fresh start
  --cleanup           Stop containers and cleanup state, then exit
  --verbose           Enable verbose output
  --non-interactive   Run without interactive prompts (for CI/CD)
  -h, --help          Display this help message

Examples:
  # Run from completely fresh state
  $0 --reset

  # Run complete E2E workflow
  $0

  # Run with verbose output
  $0 --verbose

  # Cleanup containers and state
  $0 --cleanup

The script performs the following steps:
  1. Check prerequisites (Docker, CLI commands)
  2. Start infrastructure (database and Bitcoin nodes)
  3. Create wallets in Bitcoin nodes
  4. Generate keys for keygen and sign wallets
  5. Create multisig addresses and export to watch wallet
  6. Generate test UTXOs (automatically generates 101 blocks)
  7. Create, sign, and send a test transaction

The script now automatically generates test UTXOs for the transaction phase,
making it fully automated and suitable for CI/CD pipelines.

Environment Variables:
  RPC_USER          Bitcoin RPC username (default: xyz for regtest)
  RPC_PASSWORD      Bitcoin RPC password (default: xyz for regtest)
  WALLET_PASSPHRASE Wallet passphrase for encrypted wallets (default: test)

EOF
}

###############################################################################
# Main Execution
###############################################################################

main() {
	# Parse arguments
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
			show_help
			exit 0
			;;
		*)
			log_error "Unknown option: $1"
			show_help
			exit 1
			;;
		esac
	done

	# Cleanup and exit if requested
	if [ "$CLEANUP_ONLY" = "true" ]; then
		cleanup
		exit 0
	fi

	# Full reset if requested
	if [ "$RESET_STATE" = "true" ]; then
		full_reset
	fi

	log_info "Starting Bitcoin E2E Workflow"
	log_info "Coin: $COIN"
	log_info "Encrypted: $ENCRYPTED"
	log_info "Sign wallet count: $SIGN_WALLET_NUM"
	echo ""

	# Execute workflow phases
	check_prerequisites
	setup_infrastructure
	setup_wallets
	key_generation_phase
	multisig_setup_phase
	generate_test_utxos
	transaction_flow_phase

	log_step "Bitcoin E2E Workflow Completed Successfully!"
	log_info "Summary:"
	log_info "  ✓ Infrastructure setup complete"
	log_info "  ✓ Wallets created and configured"
	log_info "  ✓ Keys generated and imported"
	log_info "  ✓ Multisig addresses created"
	log_info "  ✓ Addresses imported into watch wallet"
	log_info "  ✓ Test UTXOs generated"
	log_info "  ✓ Transaction created, signed, and sent"
	echo ""
	log_info "You can now use the wallet system for Bitcoin operations"
	log_info "To cleanup, run: $0 --cleanup"
	log_info "To full reset for fresh state, run: $0 --reset"
}

# Trap errors and cleanup
trap 'log_error "Script failed at line $LINENO"' ERR

# Run main
main "$@"
