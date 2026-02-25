#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern 11: P2TR Tapscript M-of-N
# This script automates the complete Bitcoin workflow for M-of-N multisig P2TR (Taproot) Tapscript transactions
# Usage: ./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh [OPTIONS]
#
# Transaction Pattern:
#   Pattern 11: BTC P2TR Tapscript M-of-N
#   - Address Type: P2TR (BIP86 Taproot + BIP342 Tapscript)
#   - Address Format: `bc1p...` (Mainnet), `bcrt1p...` (Regtest)
#   - Signature Requirement: M-of-N threshold (2-of-3)
#
# IMPORTANT NOTE:
#   This script demonstrates the Tapscript M-of-N workflow framework.
#   Tapscript implementation is currently under development.

set -eu

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# E2E Pattern identifier for database file naming
# Format: watch-e2e-p11-{timestamp}.db, keygen-e2e-p11-{timestamp}.db, etc.
export E2E_PATTERN="p11"

# Source BTC common utilities
# shellcheck source=../btc_common.sh
source "${SCRIPT_DIR}/../btc_common.sh"

# Initialize BTC config paths
btc_get_config_paths

# Pattern-specific configuration
SIGN_WALLET_NUM=2 # 2-of-3: keygen + sign1 + sign2 = 3 signers
VERBOSE=false
CLEANUP_ONLY=false
NON_INTERACTIVE=false
RESET_STATE=false

# Use 2-of-3 account configuration for Tapscript 2-of-3
CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account/account_2of3.yaml"
export BTC_ACCOUNT_CONF="${CONFIG_ACCOUNT}"

# Pattern 11 requires: address_type: "taproot"
export WALLET_ADDRESS_TYPE="taproot"

###############################################################################
# Key Generation Phase
###############################################################################

key_generation_phase() {
	log_step "Key Generation Phase"

	log_substep "Creating seed for keygen wallet"
	btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" create seed || {
		log_warn "Seed already exists, continuing..."
	}

	log_substep "Creating HD keys for keygen wallet"
	for account in deposit payment stored; do
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" create hdkey --account "$account" --keynum 10
	done

	log_substep "Importing private keys into keygen wallet"
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" api walletpassphrase --passphrase "${BTC_WALLET_PASSPHRASE}"
	fi
	for account in deposit payment stored; do
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
		btc_sign_cmd "$i" -c "${!config_var}" --coin "${BTC_COIN}" create hdkey
	done

	log_substep "Importing private keys into sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		config_var="BTC_CONFIG_SIGN${i}"
		if [ "${BTC_ENCRYPTED}" = "true" ]; then
			btc_sign_cmd "$i" -c "${!config_var}" --coin "${BTC_COIN}" api walletpassphrase --passphrase "${BTC_WALLET_PASSPHRASE}"
		fi
		btc_sign_cmd "$i" -c "${!config_var}" --coin "${BTC_COIN}" import privkey
		if [ "${BTC_ENCRYPTED}" = "true" ]; then
			btc_sign_cmd "$i" -c "${!config_var}" --coin "${BTC_COIN}" api walletlock
		fi
	done

	log_substep "Exporting full public keys from sign wallets"
	file_fullpubkey_auth1=$(btc_sign1_cmd -c "${BTC_CONFIG_SIGN1}" --coin "${BTC_COIN}" export fullpubkey)
	file_fullpubkey_auth2=$(btc_sign2_cmd -c "${BTC_CONFIG_SIGN2}" --coin "${BTC_COIN}" export fullpubkey)

	export FULLPUBKEY_FILE1="${file_fullpubkey_auth1##*\[fileName\]: }"
	export FULLPUBKEY_FILE2="${file_fullpubkey_auth2##*\[fileName\]: }"
}

###############################################################################
# Tapscript Setup Phase (2-of-3)
###############################################################################

tapscript_setup_phase() {
	log_step "Tapscript Setup Phase (2-of-3 P2TR Script Path)"

	log_substep "Importing full public keys into keygen wallet"
	btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" import fullpubkey --file "${FULLPUBKEY_FILE1}"
	btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" import fullpubkey --file "${FULLPUBKEY_FILE2}"

	log_warn "Tapscript descriptor export is currently a placeholder"

	log_substep "Exporting descriptors from keygen wallet"
	declare -A descriptor_paths
	# Use pattern-specific descriptor file path for parallel execution
	local descriptor_suffix=""
	if [ -n "${E2E_PATTERN}" ]; then
		descriptor_suffix="-${E2E_PATTERN}"
	fi
	for account in deposit payment stored; do
		output_file="data/descriptor/btc/${account}_descriptors${descriptor_suffix}.json"
		cmd_output=$(btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --account-config "${BTC_ACCOUNT_CONF}" --coin "${BTC_COIN}" create descriptor export --account "$account" --output "$output_file" --format bitcoin-core --include-change)
		descriptor_paths[$account]="${cmd_output##*exported to }"
	done

	log_substep "Importing descriptors into watch wallet"
	for account in deposit payment stored; do
		btc_watch_cmd -c "${BTC_CONFIG_WATCH}" --coin "${BTC_COIN}" import descriptor import --file "${descriptor_paths[$account]}" --account "$account"
	done

	first_descriptor=$(jq -r '.[0].desc // empty' "${descriptor_paths[payment]}" 2>/dev/null)
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

	receiver1=$(btc_cli "btc-watch" -rpcwallet="${BTC_WATCH_WALLET_NAME:-watch}" getnewaddress "" bech32m)
	receiver2=$(btc_cli "btc-watch" -rpcwallet="${BTC_WATCH_WALLET_NAME:-watch}" getnewaddress "" bech32m)
	receiver3=$(btc_cli "btc-watch" -rpcwallet="${BTC_WATCH_WALLET_NAME:-watch}" getnewaddress "" bech32m)

	btc_insert_payment_requests "$sender_address" "$receiver1 $receiver2 $receiver3" "0.001 0.002 0.0015"
}

###############################################################################
# Transaction Flow Phase (Tapscript 2-of-3)
###############################################################################

transaction_flow_phase() {
	log_step "Transaction Flow Phase (Tapscript 2-of-3 Script Path)"
	log_warn "NOTE: Current Tapscript implementation contains placeholder TODOs"

	log_substep "Creating unsigned payment transaction (PSBT)"
	tx_file=$(btc_watch_cmd -c "${BTC_CONFIG_WATCH}" --coin "${BTC_COIN}" create payment 2>&1) || {
		log_error "Failed to create payment transaction"
		log_error "Error details: $tx_file"
		if echo "$tx_file" | grep -q "No utxo"; then
			btc_log_no_utxo_error
		fi
		return 1
	}

	tx_unsigned=$(btc_extract_file_path "$tx_file")

	# Using standard signature as placeholder for Tapscript
	log_substep "Signing with keygen wallet (1st Schnorr signature - placeholder)"
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" api walletpassphrase --passphrase "${BTC_WALLET_PASSPHRASE}"
	fi
	tx_file_signed=$(btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" sign signature --file "${tx_unsigned}")
	if [ "${BTC_ENCRYPTED}" = "true" ]; then
		btc_keygen_cmd -c "${BTC_CONFIG_KEYGEN}" --coin "${BTC_COIN}" api walletlock
	fi

	tx_signed1=$(btc_extract_file_path "$tx_file_signed")

	log_substep "Signing with sign1 wallet (2nd signature - completing 2-of-3)"
	tx_file_signed2=$(btc_sign1_cmd -c "${BTC_CONFIG_SIGN1}" --coin "${BTC_COIN}" sign signature --file "${tx_signed1}")
	tx_signed2=$(btc_extract_file_path "$tx_file_signed2")

	log_substep "Sending transaction"
	tx_result=$(btc_watch_cmd -c "${BTC_CONFIG_WATCH}" --coin "${BTC_COIN}" send tx --file "${tx_signed2}")
	tx_id="${tx_result##*txID: }"

	log_info "Transaction sent successfully!"
	log_info "Transaction ID: $tx_id"
}

###############################################################################
# Help & Main
###############################################################################

show_help() {
	cat <<EOF
Bitcoin E2E Workflow Script - Pattern 11: P2TR Tapscript M-of-N

Usage: $0 [OPTIONS]

Options:
  --reset, --cleanup, --verbose, --non-interactive, -h/--help

Transaction Pattern: P2TR Tapscript 2-of-3 (Address: bcrt1p... in Regtest)
Note: Tapscript implementation is currently a placeholder.
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

	log_info "Starting Bitcoin E2E Workflow - Pattern 11: P2TR Tapscript M-of-N"
	log_warn "IMPORTANT: Tapscript CLI commands contain placeholder implementations"
	echo ""

	btc_check_prerequisites "watch keygen sign1 sign2"
	btc_setup_infrastructure "btc-watch btc-keygen btc-sign1 btc-sign2"
	btc_setup_wallets "watch keygen sign1 sign2"
	key_generation_phase
	tapscript_setup_phase
	generate_test_utxos || exit 1
	create_payment_requests_phase || exit 1
	transaction_flow_phase || exit 1

	log_step "Bitcoin E2E Workflow Completed Successfully!"
}

trap 'log_error "Script failed at line $LINENO"' ERR
main "$@"
