#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern 7: P2WSH Native SegWit 3-of-3 Multisig
# This script automates the complete Bitcoin workflow for 3-of-3 multisig P2WSH transactions
# Usage: ./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh [OPTIONS]
#
# Transaction Pattern:
#   Pattern 7: BTC P2WSH Native SegWit 3-of-3 Multisig
#   - Address Type: P2WSH (BIP84 Native SegWit)
#   - Address Format: `bc1q...` (Mainnet), `bcrt1q...` (Regtest)
#   - Signature Requirement: 3-of-3 (all 3 signatures required)

set -euo pipefail

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# E2E Pattern identifier for database file naming
# Format: watch-e2e-p7-{timestamp}.db, keygen-e2e-p7-{timestamp}.db, etc.
export E2E_PATTERN="p7"

# Source BTC common utilities
# shellcheck source=../btc_common.sh
source "${SCRIPT_DIR}/../btc_common.sh"

# Initialize BTC config paths
btc_get_config_paths

# Pattern-specific configuration
SIGN_WALLET_NUM=2 # keygen + sign1 + sign2 = 3 signers
VERBOSE=false
CLEANUP_ONLY=false
NON_INTERACTIVE=false
RESET_STATE=false

# Use 3-of-3 multisig account configuration
CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account/account_3of3.yaml"
export BTC_ACCOUNT_CONF="${CONFIG_ACCOUNT}"

# Pattern 7 requires: address_type: "bech32"
export WALLET_ADDRESS_TYPE="bech32"

###############################################################################
# Key Generation Phase
###############################################################################

key_generation_phase() {
	log_step "Key Generation Phase"

	log_substep "Creating seed for keygen wallet"
	keygen -c "${BTC_CONFIG_KEYGEN}" create seed || {
		log_warn "Seed already exists, continuing..."
	}

	log_substep "Creating HD keys for keygen wallet"
	for account in client deposit payment stored; do
		keygen -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" create hdkey --account "$account" --keynum 10
	done

	log_substep "Importing private keys into keygen wallet"
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		keygen -c "${BTC_CONFIG_KEYGEN}" api walletpassphrase --passphrase "${BTC_WALLET_PASSPHRASE}"
	fi
	for account in client deposit payment stored; do
		keygen -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" import privkey --account "$account"
	done
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		keygen -c "${BTC_CONFIG_KEYGEN}" api walletlock
	fi

	log_substep "Creating seeds for sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		config_var="BTC_CONFIG_SIGN${i}"
		"sign${i}" --conf "${!config_var}" --coin "${BTC_COIN}" create seed || {
			log_warn "Sign${i} seed already exists, continuing..."
		}
	done

	log_substep "Creating HD keys for sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		config_var="BTC_CONFIG_SIGN${i}"
		"sign${i}" --conf "${!config_var}" --coin "${BTC_COIN}" --wallet "sign${i}" create hdkey
	done

	log_substep "Importing private keys into sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		config_var="BTC_CONFIG_SIGN${i}"
		if [ "${BTC_ENCRYPTED}" = "true" ]; then
			"sign${i}" --conf "${!config_var}" --coin "${BTC_COIN}" --wallet "sign${i}" api walletpassphrase --passphrase "${BTC_WALLET_PASSPHRASE}"
		fi
		"sign${i}" --conf "${!config_var}" --coin "${BTC_COIN}" --wallet "sign${i}" import privkey
		if [ "${BTC_ENCRYPTED}" = "true" ]; then
			"sign${i}" --conf "${!config_var}" --coin "${BTC_COIN}" --wallet "sign${i}" api walletlock
		fi
	done

	log_substep "Exporting full public keys from sign wallets"
	file_fullpubkey_auth1=$(sign1 --conf "${BTC_CONFIG_SIGN1}" --coin "${BTC_COIN}" --wallet sign1 export fullpubkey)
	file_fullpubkey_auth2=$(sign2 --conf "${BTC_CONFIG_SIGN2}" --coin "${BTC_COIN}" --wallet sign2 export fullpubkey)

	export FULLPUBKEY_FILE1="${file_fullpubkey_auth1##*\[fileName\]: }"
	export FULLPUBKEY_FILE2="${file_fullpubkey_auth2##*\[fileName\]: }"
}

###############################################################################
# Multisig Setup Phase (3-of-3)
###############################################################################

multisig_setup_phase() {
	log_step "Multisig Setup Phase (3-of-3 Native SegWit)"

	log_substep "Importing full public keys into keygen wallet"
	keygen -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" import fullpubkey --file "${FULLPUBKEY_FILE1}"
	keygen -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" import fullpubkey --file "${FULLPUBKEY_FILE2}"

	log_substep "Exporting descriptors from keygen wallet"
	file_descriptor_deposit=$(btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" descriptor export --account deposit --output data/descriptor/btc/deposit_descriptors.json --format bitcoin-core --include-change)
	file_descriptor_payment=$(btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" descriptor export --account payment --output data/descriptor/btc/payment_descriptors.json --format bitcoin-core --include-change)
	file_descriptor_stored=$(btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" descriptor export --account stored --output data/descriptor/btc/stored_descriptors.json --format bitcoin-core --include-change)

	descriptor_deposit="${file_descriptor_deposit##*exported to }"
	descriptor_payment="${file_descriptor_payment##*exported to }"
	descriptor_stored="${file_descriptor_stored##*exported to }"

	log_substep "Importing descriptors into watch wallet"
	btc_watch_cmd -c "${BTC_CONFIG_WATCH}" --coin "${BTC_COIN}" import descriptor --file "${descriptor_deposit}" --account deposit
	btc_watch_cmd -c "${BTC_CONFIG_WATCH}" --coin "${BTC_COIN}" import descriptor --file "${descriptor_payment}" --account payment
	btc_watch_cmd -c "${BTC_CONFIG_WATCH}" --coin "${BTC_COIN}" import descriptor --file "${descriptor_stored}" --account stored

	first_descriptor=$(jq -r '.[0].desc // empty' "${descriptor_payment}" 2>/dev/null)
	if [ -z "$first_descriptor" ]; then
		log_error "Failed to extract descriptor"
		return 1
	fi
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

	sender_address="$payment_address"

	receiver1=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" bech32)
	receiver2=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" bech32)
	receiver3=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" bech32)

	btc_insert_payment_requests "$sender_address" "$receiver1 $receiver2 $receiver3" "0.001 0.002 0.0015"
}

###############################################################################
# Transaction Flow Phase (3-of-3 Multisig)
###############################################################################

transaction_flow_phase() {
	log_step "Transaction Flow Phase (3-of-3 P2WSH Native SegWit)"

	log_substep "Creating unsigned payment transaction"
	tx_file=$(btc_watch_cmd -c "${BTC_CONFIG_WATCH}" create payment 2>&1) || {
		log_error "Failed to create payment transaction"
		if echo "$tx_file" | grep -q "No utxo"; then
			btc_log_no_utxo_error
		fi
		return 1
	}

	tx_unsigned=$(btc_extract_file_path "$tx_file")

	log_substep "Signing with keygen wallet (1st signature)"
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		keygen -c "${BTC_CONFIG_KEYGEN}" api walletpassphrase --passphrase "${BTC_WALLET_PASSPHRASE}"
	fi
	tx_file_signed=$(keygen -c "${BTC_CONFIG_KEYGEN}" sign signature --file "${tx_unsigned}")
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		keygen -c "${BTC_CONFIG_KEYGEN}" api walletlock
	fi

	tx_signed1=$(btc_extract_file_path "$tx_file_signed")

	log_substep "Signing with sign1 wallet (2nd signature)"
	tx_file_signed2=$(sign1 --conf "${BTC_CONFIG_SIGN1}" --wallet sign1 sign signature --file "${tx_signed1}")
	tx_signed2=$(btc_extract_file_path "$tx_file_signed2")

	log_substep "Signing with sign2 wallet (3rd signature - completing 3-of-3)"
	tx_file_signed3=$(sign2 --conf "${BTC_CONFIG_SIGN2}" --wallet sign2 sign signature --file "${tx_signed2}")
	tx_signed3=$(btc_extract_file_path "$tx_file_signed3")

	log_substep "Sending fully signed transaction"
	tx_result=$(btc_watch_cmd -c "${BTC_CONFIG_WATCH}" send --file "${tx_signed3}")
	tx_id="${tx_result##*txID: }"

	log_info "Transaction sent successfully!"
	log_info "Transaction ID: $tx_id"
}

###############################################################################
# Help & Main
###############################################################################

show_help() {
	cat <<EOF
Bitcoin E2E Workflow Script - Pattern 7: P2WSH Native SegWit 3-of-3 Multisig

Usage: $0 [OPTIONS]

Options:
  --reset, --cleanup, --verbose, --non-interactive, -h/--help

Transaction Pattern: P2WSH 3-of-3 Multisig (Address: bcrt1q... in Regtest)
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
	if [ "$RESET_STATE" = "true" ]; then btc_full_reset "${BTC_DOCKER_VOLUME_NAME}" "watch keygen sign1 sign2"; fi

	log_info "Starting Bitcoin E2E Workflow - Pattern 7: P2WSH Native SegWit 3-of-3 Multisig"
	echo ""

	btc_check_prerequisites "watch keygen sign1 sign2"
	btc_setup_infrastructure "btc-watch btc-keygen btc-sign1 btc-sign2"
	btc_setup_wallets "watch keygen sign1 sign2"
	key_generation_phase
	multisig_setup_phase
	generate_test_utxos
	create_payment_requests_phase
	transaction_flow_phase

	log_step "Bitcoin E2E Workflow Completed Successfully!"
}

trap 'log_error "Script failed at line $LINENO"' ERR
main "$@"
