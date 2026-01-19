#!/usr/bin/env bash

# Bitcoin Cash E2E Workflow Script - Pattern 1: P2PKH Single-sig
# This script automates the complete Bitcoin Cash workflow for single-sig P2PKH transactions
# Usage: ./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh [OPTIONS]
#
# Transaction Pattern:
#   Pattern 1: BCH P2PKH Single-sig
#   - Address Type: P2PKH (BIP44 Legacy)
#   - Address Format: CashAddr `bitcoincash:q...` (Mainnet), `m.../n...` (Regtest)
#   - Signature Requirement: Single-sig (Keygen only)
#   - Key Derivation: m/44'/145'/account'/change/index (Mainnet), m/44'/1'/... (Testnet/Regtest)
#
# Note: BCH does NOT support SegWit or descriptor wallets

set -eu

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# E2E Pattern identifier for database file naming
# Format: watch-e2e-p1-{timestamp}.db, keygen-e2e-p1-{timestamp}.db, etc.
export E2E_PATTERN="p1"

# Source BCH common utilities (includes common.sh automatically)
# shellcheck source=../bch_common.sh
source "${SCRIPT_DIR}/../bch_common.sh"

# Initialize BCH config paths
bch_get_config_paths

# Pattern-specific configuration
SIGN_WALLET_NUM=0 # Single-sig: no additional sign wallets needed
VERBOSE=false
CLEANUP_ONLY=false
NON_INTERACTIVE=false
RESET_STATE=false

# Use single-sig account configuration for Pattern 1
CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account/account.yaml"
export BCH_ACCOUNT_CONF="${CONFIG_ACCOUNT}"

###############################################################################
# Key Generation Phase
###############################################################################

key_generation_phase() {
	log_step "Key Generation Phase"

	log_substep "Creating seed for keygen wallet"
	bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" create seed || {
		log_warn "Seed already exists or error occurred, continuing..."
	}

	log_substep "Creating HD keys for keygen wallet (client, deposit, payment, stored)"
	for account in client deposit payment stored; do
		log_info "Creating HD keys for account: $account"
		bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" create hdkey --account "$account" --keynum 10
	done
}

###############################################################################
# Single-sig Address Setup Phase
###############################################################################

singlesig_setup_phase() {
	log_step "Single-sig Address Setup Phase"

	# BCH doesn't support descriptor wallets, use address export/import workflow
	log_substep "Importing private keys into keygen wallet"
	if [ "${BCH_ENCRYPTED}" = "true" ]; then
		bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" api walletpassphrase --passphrase "${BCH_WALLET_PASSPHRASE}"
	fi
	for account in client deposit payment stored; do
		log_info "Importing private keys for account: $account"
		bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" import privkey --account "$account"
	done
	if [ "${BCH_ENCRYPTED}" = "true" ]; then
		bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" api walletlock
	fi

	# Export addresses from keygen
	log_substep "Exporting addresses from keygen wallet"
	local accounts=("client" "payment")
	declare -A address_files

	for account in "${accounts[@]}"; do
		log_info "Exporting ${account} addresses"
		file_output=$(bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" export address --account "${account}")
		address_files[$account]=$(bch_extract_file_path "$file_output")
		log_info "  ${account}: ${address_files[$account]}"
	done

	# Import addresses into watch wallet
	log_substep "Importing addresses into watch wallet"
	for account in "${accounts[@]}"; do
		log_info "Importing ${account} addresses"
		bch_watch_cmd -c "${BCH_CONFIG_WATCH}" --coin "${BCH_COIN}" import address --file "${address_files[$account]}"
	done

	log_info "All addresses imported successfully"

	# Store payment address file for UTXO generation
	export ADDRESS_FILE_PAYMENT="${address_files[payment]}"
}

###############################################################################
# UTXO Generation Phase
###############################################################################

generate_test_utxos() {
	log_step "Generating Test UTXOs for Transaction Phase"

	log_info "Reading payment address from exported file..."
	# Get first payment address from the exported CSV file
	# CSV format: coin,account,P2PKH,P2SH-segwit,bech32,taproot,pubkey,multisig,redeem_script,index
	# For BCH Pattern 1 (single-sig), use field 3 (P2PKH address)
	payment_address=$(grep -v '^#' "${ADDRESS_FILE_PAYMENT}" 2>/dev/null | head -n1 | cut -d',' -f3)

	if [ -z "$payment_address" ]; then
		log_error "Failed to extract payment address from ${ADDRESS_FILE_PAYMENT}"
		return 1
	fi
	export payment_address

	bch_generate_test_utxos "$payment_address" 101
	bch_wait_for_balance 60
}

###############################################################################
# Payment Request Creation Phase
###############################################################################

create_payment_requests_phase() {
	log_step "Payment Request Creation Phase"

	log_substep "Retrieving payment sender address from database"
	sender_address=$(bch_get_sender_address "payment")

	if [ -z "$sender_address" ]; then
		log_error "No payment addresses found in database"
		return 1
	fi

	log_info "Using sender address: $sender_address"

	# Generate receiver addresses (BCH uses legacy format)
	# Use pattern-specific wallet name
	local wallet_name="${BCH_WATCH_WALLET_NAME:-watch}"
	receiver1=$(bch_cli "bch-watch" -rpcwallet="${wallet_name}" getnewaddress "" legacy)
	receiver2=$(bch_cli "bch-watch" -rpcwallet="${wallet_name}" getnewaddress "" legacy)
	receiver3=$(bch_cli "bch-watch" -rpcwallet="${wallet_name}" getnewaddress "" legacy)

	bch_insert_payment_requests "$sender_address" "$receiver1 $receiver2 $receiver3" "0.001 0.002 0.0015"
}

###############################################################################
# Transaction Flow Phase (Single-sig)
###############################################################################

transaction_flow_phase() {
	log_step "Transaction Flow Phase (Single-sig)"

	log_substep "Creating unsigned payment transaction"
	tx_file=$(bch_watch_cmd -c "${BCH_CONFIG_WATCH}" --coin "${BCH_COIN}" create payment 2>&1) || {
		log_error "Failed to create payment transaction"
		log_error "Error details: $tx_file"
		if echo "$tx_file" | grep -q "No utxo"; then
			bch_log_no_utxo_error
		fi
		return 1
	}

	tx_unsigned=$(bch_extract_file_path "$tx_file")
	log_info "Created unsigned transaction: $tx_unsigned"

	log_substep "Signing with keygen wallet (single signature)"
	if [ "${BCH_ENCRYPTED}" = "true" ]; then
		bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" api walletpassphrase --passphrase "${BCH_WALLET_PASSPHRASE}"
	fi
	tx_file_signed=$(bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" sign signature --file "${tx_unsigned}")
	if [ "${BCH_ENCRYPTED}" = "true" ]; then
		bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" api walletlock
	fi

	tx_signed=$(bch_extract_file_path "$tx_file_signed")
	log_info "Signed transaction: $tx_signed"

	log_substep "Sending fully signed transaction"
	tx_result=$(bch_watch_cmd -c "${BCH_CONFIG_WATCH}" --coin "${BCH_COIN}" send --file "${tx_signed}")
	tx_id="${tx_result##*txID: }"

	log_info "Transaction sent successfully!"
	log_info "Transaction ID: $tx_id"
}

###############################################################################
# Help & Main
###############################################################################

show_help() {
	cat <<EOF
Bitcoin Cash E2E Workflow Script - Pattern 1: P2PKH Single-sig

This script automates the complete Bitcoin Cash workflow for single-sig P2PKH
transactions. It serves as a regression test tool to verify that the BCH
single-sig workflow functions correctly after code changes.

Usage: $0 [OPTIONS]

Options:
  --reset             Full reset: cleanup all state for fresh start
  --cleanup           Stop containers and cleanup state, then exit
  --verbose           Enable verbose output
  --non-interactive   Run without interactive prompts (for CI/CD)
  -h, --help          Display this help message

Transaction Pattern:
  Pattern 1: BCH P2PKH Single-sig
  - Address Type: P2PKH (BIP44 Legacy)
  - Address Format: CashAddr \`bitcoincash:q...\` (Mainnet), \`m.../n...\` (Regtest)
  - Signature Requirement: Single-sig (Keygen only)
  - Key Derivation: m/44'/145'/account'/change/index (Mainnet)

The script performs the following steps:
  1. Check prerequisites (Docker, CLI commands)
  2. Start infrastructure (database and Bitcoin Cash nodes)
  3. Create wallets in Bitcoin Cash nodes
  4. Generate HD keys for keygen wallet
  5. Export addresses and import to watch wallet
  6. Generate test UTXOs (automatically generates 101 blocks)
  7. Create, sign (single-sig), and send a test transaction

Environment Variables:
  RPC_USER                Bitcoin Cash RPC username (default: xyz for regtest)
  RPC_PASSWORD            Bitcoin Cash RPC password (default: xyz for regtest)
  BCH_WALLET_PASSPHRASE   Wallet passphrase for encrypted wallets (default: test)

Examples:
  # Run from completely fresh state
  $0 --reset

  # Run complete E2E workflow
  $0

  # Run with verbose output
  $0 --verbose

  # Cleanup containers and state
  $0 --cleanup

EOF
}

main() {
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

	if [ "$CLEANUP_ONLY" = "true" ]; then
		bch_cleanup
		exit 0
	fi
	if [ "$RESET_STATE" = "true" ]; then
		bch_full_reset "${BCH_DOCKER_VOLUME_NAME}" "watch keygen"
	fi

	log_info "Starting Bitcoin Cash E2E Workflow - Pattern 1: P2PKH Single-sig"
	log_info "Coin: ${BCH_COIN}"
	log_info "Encrypted: ${BCH_ENCRYPTED}"
	log_info "Sign wallet count: $SIGN_WALLET_NUM (Single-sig)"
	echo ""

	bch_check_prerequisites "watch keygen"
	bch_setup_infrastructure "bch-watch bch-keygen"
	bch_setup_wallets "watch keygen"
	key_generation_phase
	singlesig_setup_phase
	generate_test_utxos
	create_payment_requests_phase
	transaction_flow_phase

	log_step "Bitcoin Cash E2E Workflow Completed Successfully!"
	log_info "Pattern 1: P2PKH single-sig transaction completed"
}

trap 'log_error "Script failed at line $LINENO"' ERR
main "$@"
