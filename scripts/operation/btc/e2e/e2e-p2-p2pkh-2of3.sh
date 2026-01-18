#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern 2: P2PKH 2-of-3 Multisig
# This script automates the complete Bitcoin workflow for 2-of-3 multisig P2PKH transactions
# Usage: ./scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh [OPTIONS]
#
# Transaction Pattern:
#   Pattern 2: BTC P2PKH 2-of-3 Multisig
#   - Address Type: P2PKH (BIP44 Legacy) wrapped in P2SH
#   - Address Format: `3...` (Mainnet), `2...` (Testnet/Regtest)
#   - Signature Requirement: 2-of-3 (any 2 signatures out of 3)
#   - Descriptor: sh(multi(2,[fingerprint/44'/0'/0']xpub1.../0/*,xpub2.../0/*,xpub3.../0/*))

set -euo pipefail

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# E2E Pattern identifier for database file naming
# Format: watch-e2e-p2-{timestamp}.db, keygen-e2e-p2-{timestamp}.db, etc.
export E2E_PATTERN="p2"

# Source BTC common utilities (includes common.sh automatically)
# shellcheck source=../btc_common.sh
source "${SCRIPT_DIR}/../btc_common.sh"

# Initialize BTC config paths
btc_get_config_paths

# Pattern-specific configuration
SIGN_WALLET_NUM=2 # 2-of-3: need sign1 and sign2 wallets
VERBOSE=false
CLEANUP_ONLY=false
NON_INTERACTIVE=false
RESET_STATE=false

# Use 2-of-3 multisig account configuration for Pattern 2
CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account/account_2of3.yaml"
export BTC_ACCOUNT_CONF="${CONFIG_ACCOUNT}"

# Pattern 2 requires: address_type: "legacy"
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

	log_substep "Creating HD keys for keygen wallet"
	for account in client deposit payment stored; do
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" create hdkey --account "$account" --keynum 10
	done

	log_substep "Importing private keys into keygen wallet"
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" api walletpassphrase --passphrase "${BTC_WALLET_PASSPHRASE}"
	fi
	for account in client deposit payment stored; do
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" import privkey --account "$account"
	done
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" api walletlock
	fi

	log_substep "Creating seeds for sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		config_var="BTC_CONFIG_SIGN${i}"
		btc_sign_cmd "$i" -c "${!config_var}" --coin "${BTC_COIN}" create seed || {
			log_warn "Sign${i} seed already exists, continuing..."
		}
	done

	log_substep "Creating HD keys for sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		config_var="BTC_CONFIG_SIGN${i}"
		btc_sign_cmd "$i" -c "${!config_var}" --coin "${BTC_COIN}" --wallet "sign${i}" create hdkey
	done

	log_substep "Importing private keys into sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		config_var="BTC_CONFIG_SIGN${i}"
		if [ "${BTC_ENCRYPTED}" = "true" ]; then
			btc_sign_cmd "$i" -c "${!config_var}" --coin "${BTC_COIN}" --wallet "sign${i}" api walletpassphrase --passphrase "${BTC_WALLET_PASSPHRASE}"
		fi
		btc_sign_cmd "$i" -c "${!config_var}" --coin "${BTC_COIN}" --wallet "sign${i}" import privkey
		if [ "${BTC_ENCRYPTED}" = "true" ]; then
			btc_sign_cmd "$i" -c "${!config_var}" --coin "${BTC_COIN}" --wallet "sign${i}" api walletlock
		fi
	done

	log_substep "Exporting full public keys from sign wallets"
	file_fullpubkey_auth1=$(btc_sign1_cmd -c "${BTC_CONFIG_SIGN1}" --coin "${BTC_COIN}" --wallet sign1 export fullpubkey)
	file_fullpubkey_auth2=$(btc_sign2_cmd -c "${BTC_CONFIG_SIGN2}" --coin "${BTC_COIN}" --wallet sign2 export fullpubkey)

	export FULLPUBKEY_FILE1="${file_fullpubkey_auth1##*\[fileName\]: }"
	export FULLPUBKEY_FILE2="${file_fullpubkey_auth2##*\[fileName\]: }"

	log_info "Exported fullpubkey files:"
	log_info "  sign1: $FULLPUBKEY_FILE1"
	log_info "  sign2: $FULLPUBKEY_FILE2"
}

###############################################################################
# Multisig Setup Phase (2-of-3)
###############################################################################

multisig_setup_phase() {
	log_step "Multisig Setup Phase (2-of-3)"

	log_substep "Importing full public keys into keygen wallet"
	btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" import fullpubkey --file "${FULLPUBKEY_FILE1}"
	btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" import fullpubkey --file "${FULLPUBKEY_FILE2}"

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

	log_info "All descriptors imported successfully"

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
	log_info "Using sender address: $sender_address"

	receiver1=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" legacy)
	receiver2=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" legacy)
	receiver3=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" legacy)

	btc_insert_payment_requests "$sender_address" "$receiver1 $receiver2 $receiver3" "0.001 0.002 0.0015"
}

###############################################################################
# Transaction Flow Phase (2-of-3 Multisig)
###############################################################################

transaction_flow_phase() {
	log_step "Transaction Flow Phase (2-of-3 Multisig)"

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

	log_substep "Signing with keygen wallet (1st signature)"
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" api walletpassphrase --passphrase "${BTC_WALLET_PASSPHRASE}"
	fi
	tx_file_signed=$(btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" sign signature --file "${tx_unsigned}")
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" api walletlock
	fi

	tx_signed1=$(btc_extract_file_path "$tx_file_signed")
	log_info "Signed transaction (1st): $tx_signed1"

	log_substep "Signing with sign1 wallet (2nd signature - completing 2-of-3)"
	tx_file_signed2=$(btc_sign1_cmd -c "${BTC_CONFIG_SIGN1}" --coin "${BTC_COIN}" --wallet sign1 sign signature --file "${tx_signed1}")
	tx_signed2=$(btc_extract_file_path "$tx_file_signed2")
	log_info "Signed transaction (2nd): $tx_signed2"

	log_substep "Sending fully signed transaction"
	tx_result=$(btc_watch_cmd -c "${BTC_CONFIG_WATCH}" --coin "${BTC_COIN}" send --file "${tx_signed2}")
	tx_id="${tx_result##*txID: }"

	log_info "Transaction sent successfully!"
	log_info "Transaction ID: $tx_id"
}

###############################################################################
# Help & Main
###############################################################################

show_help() {
	cat <<EOF
Bitcoin E2E Workflow Script - Pattern 2: P2PKH 2-of-3 Multisig

Usage: $0 [OPTIONS]

Options:
  --reset, --cleanup, --verbose, --non-interactive, -h/--help

Transaction Pattern: P2PKH 2-of-3 Multisig (Address: 2... in Regtest)
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

	log_info "Starting Bitcoin E2E Workflow - Pattern 2: P2PKH 2-of-3 Multisig"
	log_info "Coin: ${BTC_COIN}"
	log_info "Encrypted: ${BTC_ENCRYPTED}"
	log_info "Sign wallet count: $SIGN_WALLET_NUM (2-of-3 Multisig)"
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
	log_info "Pattern 2: P2PKH 2-of-3 multisig transaction completed"
}

trap 'log_error "Script failed at line $LINENO"' ERR
main "$@"
