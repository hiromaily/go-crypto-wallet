#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern 1: P2PKH Single-sig
# This script automates the complete Bitcoin workflow for single-sig P2PKH transactions
# Usage: ./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh [OPTIONS]
#
# Transaction Pattern:
#   Pattern 1: BTC P2PKH Single-sig
#   - Address Type: P2PKH (BIP44 Legacy)
#   - Address Format: `1...` (Mainnet), `m.../n...` (Testnet/Regtest)
#   - Signature Requirement: Single-sig (Keygen only)
#   - Descriptor: pkh([fingerprint/44'/0'/0']xpub.../0/*)

set -eu

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# E2E Pattern identifier for database file naming
# Format: watch-e2e-p1-{timestamp}.db, keygen-e2e-p1-{timestamp}.db, etc.
export E2E_PATTERN="p1"

# Source BTC common utilities (includes common.sh automatically)
# shellcheck source=../btc_common.sh
source "${SCRIPT_DIR}/../btc_common.sh"

# Initialize BTC config paths
btc_get_config_paths

# Pattern-specific configuration
SIGN_WALLET_NUM=0 # Single-sig: no additional sign wallets needed
VERBOSE=false
CLEANUP_ONLY=false
NON_INTERACTIVE=false
RESET_STATE=false

# Use single-sig account configuration for Pattern 1
CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account/account.yaml"
export BTC_ACCOUNT_CONF="${CONFIG_ACCOUNT}"

# Pattern 1 requires: address_type: "legacy"
export WALLET_ADDRESS_TYPE="legacy"

###############################################################################
# Key Generation Phase
###############################################################################

key_generation_phase() {
	log_step "Key Generation Phase"

	log_substep "Creating seed for keygen wallet"
	btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" create seed || {
		log_warn "Seed already exists or error occurred, continuing..."
	}

	log_substep "Creating HD keys for keygen wallet (client, deposit, payment, stored)"
	for account in client deposit payment stored; do
		log_info "Creating HD keys for account: $account"
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" create hdkey --account "$account" --keynum 10
	done
}

###############################################################################
# Single-sig Address Setup Phase
###############################################################################

singlesig_setup_phase() {
	log_step "Single-sig Address Setup Phase"

	log_substep "Exporting descriptors from keygen wallet"
	local accounts=(client payment)
	declare -A descriptor_files

	for account in "${accounts[@]}"; do
		log_info "Exporting ${account} descriptors"
		file_output=$(btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" descriptor export \
			--account "${account}" \
			--output "data/descriptor/btc/${account}_descriptors.json" \
			--format bitcoin-core \
			--include-change)
		descriptor_files[$account]="${file_output##*exported to }"
		log_info "  ${account}: ${descriptor_files[$account]}"
	done

	log_substep "Importing descriptors into watch wallet"
	for account in "${accounts[@]}"; do
		log_info "Importing ${account} descriptors"
		btc_watch_cmd -c "${BTC_CONFIG_WATCH}" --coin "${BTC_COIN}" import descriptor \
			--file "${descriptor_files[$account]}" \
			--account "${account}"
	done

	log_info "All descriptors imported successfully"

	log_substep "Deriving payment address from descriptor for UTXO generation"
	first_descriptor=$(jq -r '.[0].desc // empty' "${descriptor_files[payment]}" 2>/dev/null)

	if [ -z "$first_descriptor" ]; then
		log_error "Failed to extract descriptor from ${descriptor_files[payment]}"
		return 1
	fi

	log_info "  Using descriptor: ${first_descriptor:0:50}..."
	export first_descriptor
}

###############################################################################
# UTXO Generation Phase
###############################################################################

generate_test_utxos() {
	log_step "Generating Test UTXOs for Transaction Phase"

	payment_address=$(btc_derive_address_from_descriptor "$first_descriptor")
	export payment_address

	btc_generate_test_utxos "$payment_address" 101
	btc_wait_for_balance 60
}

###############################################################################
# Payment Request Creation Phase
###############################################################################

create_payment_requests_phase() {
	log_step "Payment Request Creation Phase"

	log_substep "Retrieving payment sender address from database"
	sender_address=$(btc_get_sender_address "payment")

	if [ -z "$sender_address" ]; then
		log_error "No payment addresses found in database"
		return 1
	fi

	log_info "Using sender address: $sender_address"

	receiver1=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" legacy)
	receiver2=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" legacy)
	receiver3=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" legacy)

	btc_insert_payment_requests "$sender_address" "$receiver1 $receiver2 $receiver3" "0.001 0.002 0.0015"
}

###############################################################################
# Transaction Flow Phase (Single-sig)
###############################################################################

transaction_flow_phase() {
	log_step "Transaction Flow Phase (Single-sig)"

	log_substep "Creating unsigned payment transaction"
	tx_file=$(btc_watch_cmd -c "${BTC_CONFIG_WATCH}" --coin "${BTC_COIN}" create payment 2>&1) || {
		log_error "Failed to create payment transaction"
		if echo "$tx_file" | grep -q "No utxo"; then
			btc_log_no_utxo_error
		fi
		return 1
	}

	tx_unsigned=$(btc_extract_file_path "$tx_file")
	log_info "Created unsigned transaction: $tx_unsigned"

	log_substep "Signing with keygen wallet (single signature)"
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" api walletpassphrase --passphrase "${BTC_WALLET_PASSPHRASE}"
	fi
	tx_file_signed=$(btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" sign signature --file "${tx_unsigned}")
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" api walletlock
	fi

	tx_signed=$(btc_extract_file_path "$tx_file_signed")
	log_info "Signed transaction: $tx_signed"

	log_substep "Sending fully signed transaction"
	tx_result=$(btc_watch_cmd -c "${BTC_CONFIG_WATCH}" --coin "${BTC_COIN}" send --file "${tx_signed}")
	tx_id="${tx_result##*txID: }"

	log_info "Transaction sent successfully!"
	log_info "Transaction ID: $tx_id"
}

###############################################################################
# Help & Main
###############################################################################

show_help() {
	cat <<EOF
Bitcoin E2E Workflow Script - Pattern 1: P2PKH Single-sig

Usage: $0 [OPTIONS]

Options:
  --reset, --cleanup, --verbose, --non-interactive, -h/--help

Transaction Pattern: P2PKH Single-sig (Address: m.../n... in Regtest)
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
		btc_cleanup
		exit 0
	fi
	if [ "$RESET_STATE" = "true" ]; then btc_full_reset "${BTC_DOCKER_VOLUME_NAME}" "watch keygen"; fi

	log_info "Starting Bitcoin E2E Workflow - Pattern 1: P2PKH Single-sig"
	log_info "Coin: ${BTC_COIN}"
	log_info "Encrypted: ${BTC_ENCRYPTED}"
	log_info "Sign wallet count: $SIGN_WALLET_NUM (Single-sig)"
	echo ""

	btc_check_prerequisites "watch keygen"
	btc_setup_infrastructure "btc-watch btc-keygen"
	btc_setup_wallets "watch keygen"
	key_generation_phase
	singlesig_setup_phase
	generate_test_utxos
	create_payment_requests_phase
	transaction_flow_phase

	log_step "Bitcoin E2E Workflow Completed Successfully!"
	log_info "Pattern 1: P2PKH single-sig transaction completed"
}

trap 'log_error "Script failed at line $LINENO"' ERR
main "$@"
