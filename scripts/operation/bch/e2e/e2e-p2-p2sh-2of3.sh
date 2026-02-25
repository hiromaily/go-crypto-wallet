#!/usr/bin/env bash

# Bitcoin Cash E2E Workflow Script - Pattern 2: P2SH 2-of-3 Multisig
# This script automates the complete Bitcoin Cash workflow for 2-of-3 multisig P2SH transactions
# Usage: ./scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh [OPTIONS]
#
# Transaction Pattern:
#   Pattern 2: BCH P2SH 2-of-3 Multisig
#   - Address Type: P2SH (BIP44 + BIP11)
#   - Address Format: CashAddr `bitcoincash:p...` (Mainnet), `2...` (Regtest)
#   - Signature Requirement: 2-of-3 Multisig (any 2 of Keygen, Sign1, Sign2)
#   - Key Derivation: m/44'/145'/account'/change/index (Mainnet), m/44'/1'/... (Testnet/Regtest)
#
# Note: BCH does NOT support SegWit or descriptor wallets

set -eu

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# E2E Pattern identifier for database file naming
# Format: watch-e2e-p2-{timestamp}.db, keygen-e2e-p2-{timestamp}.db, etc.
export E2E_PATTERN="p2"

# Source BCH common utilities (includes common.sh automatically)
# shellcheck source=../bch_common.sh
source "${SCRIPT_DIR}/../bch_common.sh"

# Initialize BCH config paths
bch_get_config_paths

# Pattern-specific configuration
SIGN_WALLET_NUM=2 # 2-of-3: need sign1 and sign2 wallets
VERBOSE=false
CLEANUP_ONLY=false
NON_INTERACTIVE=false
RESET_STATE=false

# Use 2-of-3 multisig account configuration for Pattern 2
CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account/account_2of3.yaml"
export BCH_ACCOUNT_CONF="${CONFIG_ACCOUNT}"

###############################################################################
# Key Generation Phase
###############################################################################

key_generation_phase() {
	log_step "Key Generation Phase"

	# Keygen wallet - create seed
	log_substep "Creating seed for keygen wallet"
	bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" create seed || {
		log_warn "Seed already exists or error occurred, continuing..."
	}

	# Keygen wallet - create hdkeys
	log_substep "Creating HD keys for keygen wallet (client, deposit, payment, stored)"
	for account in client deposit payment stored; do
		log_info "Creating HD keys for account: $account"
		bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" create hdkey --account "$account" --keynum 10
	done

	# Keygen wallet - import private keys
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

	# Sign wallets - create seed
	log_substep "Creating seeds for sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		config_var="BCH_CONFIG_SIGN${i}"
		bch_sign_cmd "$i" -c "${!config_var}" --coin "${BCH_COIN}" create seed || {
			log_warn "Sign${i} seed already exists, continuing..."
		}
	done

	# Sign wallets - create hdkeys
	log_substep "Creating HD keys for sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		log_info "Creating HD keys for sign${i}"
		config_var="BCH_CONFIG_SIGN${i}"
		bch_sign_cmd "$i" -c "${!config_var}" --coin "${BCH_COIN}" create hdkey
	done

	# Sign wallets - import private keys
	log_substep "Importing private keys into sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		log_info "Importing private keys for sign${i}"
		config_var="BCH_CONFIG_SIGN${i}"
		if [ "${BCH_ENCRYPTED}" = "true" ]; then
			bch_sign_cmd "$i" -c "${!config_var}" --coin "${BCH_COIN}" api walletpassphrase --passphrase "${BCH_WALLET_PASSPHRASE}"
		fi
		bch_sign_cmd "$i" -c "${!config_var}" --coin "${BCH_COIN}" import privkey
		if [ "${BCH_ENCRYPTED}" = "true" ]; then
			bch_sign_cmd "$i" -c "${!config_var}" --coin "${BCH_COIN}" api walletlock
		fi
	done

	# Sign wallets - export fullpubkey
	log_substep "Exporting full public keys from sign wallets"
	file_fullpubkey_auth1=$(bch_sign1_cmd -c "${BCH_CONFIG_SIGN1}" --coin "${BCH_COIN}" export fullpubkey)
	file_fullpubkey_auth2=$(bch_sign2_cmd -c "${BCH_CONFIG_SIGN2}" --coin "${BCH_COIN}" export fullpubkey)

	# Extract file paths
	FULLPUBKEY_FILE1="${file_fullpubkey_auth1##*\[fileName\]: }"
	FULLPUBKEY_FILE2="${file_fullpubkey_auth2##*\[fileName\]: }"

	log_info "Exported fullpubkey files:"
	log_info "  sign1: $FULLPUBKEY_FILE1"
	log_info "  sign2: $FULLPUBKEY_FILE2"

	# Store for next phase
	export FULLPUBKEY_FILE1
	export FULLPUBKEY_FILE2
}

###############################################################################
# Multisig Setup Phase (2-of-3)
###############################################################################

multisig_setup_phase() {
	log_step "Multisig Setup Phase (2-of-3)"

	# Import fullpubkeys
	log_substep "Importing full public keys into keygen wallet"
	log_info "Importing fullpubkey from sign1: $FULLPUBKEY_FILE1"
	bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" import fullpubkey --file "${FULLPUBKEY_FILE1}"

	log_info "Importing fullpubkey from sign2: $FULLPUBKEY_FILE2"
	bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" import fullpubkey --file "${FULLPUBKEY_FILE2}"

	# Create multisig addresses (2-of-3)
	log_substep "Creating 2-of-3 multisig addresses"
	for account in deposit payment stored; do
		log_info "Creating 2-of-3 multisig address for account: $account"
		bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" create multisig --account "$account"
	done

	# Export addresses
	log_substep "Exporting addresses from keygen wallet"
	file_address_client=$(bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" export address --account client)
	file_address_deposit=$(bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" export address --account deposit)
	file_address_payment=$(bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" export address --account payment)
	file_address_stored=$(bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" export address --account stored)

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
	bch_watch_cmd -c "${BCH_CONFIG_WATCH}" --coin "${BCH_COIN}" import address --file "${address_client}"

	log_info "Importing deposit addresses"
	bch_watch_cmd -c "${BCH_CONFIG_WATCH}" --coin "${BCH_COIN}" import address --file "${address_deposit}"

	log_info "Importing payment addresses"
	bch_watch_cmd -c "${BCH_CONFIG_WATCH}" --coin "${BCH_COIN}" import address --file "${address_payment}"

	log_info "Importing stored addresses"
	bch_watch_cmd -c "${BCH_CONFIG_WATCH}" --coin "${BCH_COIN}" import address --file "${address_stored}"

	# Store payment address file for UTXO generation
	export ADDRESS_FILE_PAYMENT="${address_payment}"
}

###############################################################################
# UTXO Generation Phase
###############################################################################

generate_test_utxos() {
	log_step "Generating Test UTXOs for Transaction Phase"

	log_info "Reading payment address from exported file..."
	# Get first payment address from the exported CSV file
	# CSV format: coin,account,address,,,,,pubkey,,,index
	# For BCH 2-of-3 multisig, the address is in field 3
	payment_address=$(grep -v '^#' "${ADDRESS_FILE_PAYMENT}" 2>/dev/null | head -n1 | cut -d',' -f3)

	if [ -z "$payment_address" ]; then
		log_error "Failed to extract payment address from ${ADDRESS_FILE_PAYMENT}"
		return 1
	fi
	export payment_address

	log_info "Using payment address: $payment_address"

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
# Transaction Flow Phase (2-of-3 Multisig)
###############################################################################

transaction_flow_phase() {
	log_step "Transaction Flow Phase (2-of-3 Multisig)"

	# Create unsigned transaction
	log_substep "Creating unsigned payment transaction"
	tx_file=$(bch_watch_cmd -c "${BCH_CONFIG_WATCH}" --coin "${BCH_COIN}" create payment 2>&1) || {
		log_error "Failed to create payment transaction"
		log_error "Error details: $tx_file"

		if echo "$tx_file" | grep -q "No utxo"; then
			bch_log_no_utxo_error
		fi

		return 1
	}

	# Extract file path
	tx_unsigned=$(bch_extract_file_path "$tx_file")
	log_info "Created unsigned transaction: $tx_unsigned"

	# Sign with keygen wallet (1st signature)
	log_substep "Signing with keygen wallet (1st signature)"
	if [ "${BCH_ENCRYPTED}" = "true" ]; then
		bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" api walletpassphrase --passphrase "${BCH_WALLET_PASSPHRASE}"
	fi
	tx_file_signed=$(bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" sign signature --file "${tx_unsigned}")
	if [ "${BCH_ENCRYPTED}" = "true" ]; then
		bch_keygen_cmd -c "${BCH_CONFIG_KEYGEN}" --coin "${BCH_COIN}" api walletlock
	fi

	tx_signed1=$(bch_extract_file_path "$tx_file_signed")
	log_info "Signed transaction (1st): $tx_signed1"

	# Sign with sign1 wallet (2nd signature - completing 2-of-3)
	log_substep "Signing with sign1 wallet (2nd signature - completing 2-of-3)"
	tx_file_signed2=$(bch_sign1_cmd -c "${BCH_CONFIG_SIGN1}" --coin "${BCH_COIN}" sign signature --file "${tx_signed1}")
	tx_signed2=$(bch_extract_file_path "$tx_file_signed2")
	log_info "Signed transaction (2nd): $tx_signed2"

	# Note: sign2 is NOT needed for 2-of-3 (we only need 2 signatures)
	log_info "2-of-3 multisig complete - sign2 not required"

	# Send transaction
	log_substep "Sending fully signed transaction"
	tx_result=$(bch_watch_cmd -c "${BCH_CONFIG_WATCH}" --coin "${BCH_COIN}" send tx --file "${tx_signed2}")
	tx_id="${tx_result##*txID: }"

	log_info "Transaction sent successfully!"
	log_info "Transaction ID: $tx_id"
}

###############################################################################
# Help & Main
###############################################################################

show_help() {
	cat <<EOF
Bitcoin Cash E2E Workflow Script - Pattern 2: P2SH 2-of-3 Multisig

This script automates the complete Bitcoin Cash workflow for 2-of-3 multisig
P2SH transactions. It serves as a regression test tool to verify that the BCH
2-of-3 multisig workflow functions correctly after code changes.

Usage: $0 [OPTIONS]

Options:
  --reset             Full reset: cleanup all state for fresh start
  --cleanup           Stop containers and cleanup state, then exit
  --verbose           Enable verbose output
  --non-interactive   Run without interactive prompts (for CI/CD)
  -h, --help          Display this help message

Transaction Pattern:
  Pattern 2: BCH P2SH 2-of-3 Multisig
  - Address Type: P2SH (BIP44 + BIP11)
  - Address Format: CashAddr \`bitcoincash:p...\` (Mainnet), \`2...\` (Regtest)
  - Signature Requirement: 2-of-3 (any 2 of Keygen, Sign1, Sign2)
  - Key Derivation: m/44'/145'/account'/change/index (Mainnet)

The script performs the following steps:
  1. Check prerequisites (Docker, CLI commands)
  2. Start infrastructure (database and Bitcoin Cash nodes)
  3. Create wallets in Bitcoin Cash nodes
  4. Generate keys for keygen and sign wallets (sign1, sign2)
  5. Create 2-of-3 multisig addresses and export to watch wallet
  6. Generate test UTXOs (automatically generates 101 blocks)
  7. Create, sign (2 signatures), and send a test transaction

Note: Unlike 3-of-3 multisig, this pattern only requires 2 signatures
(keygen + sign1), providing redundancy if one key is lost.

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
		bch_cleanup
		exit 0
	fi

	# Full reset if requested
	if [ "$RESET_STATE" = "true" ]; then
		bch_full_reset "${BCH_DOCKER_VOLUME_NAME}" "watch keygen sign1 sign2"
	fi

	log_info "Starting Bitcoin Cash E2E Workflow - Pattern 2: P2SH 2-of-3 Multisig"
	log_info "Coin: ${BCH_COIN}"
	log_info "Encrypted: ${BCH_ENCRYPTED}"
	log_info "Sign wallet count: $SIGN_WALLET_NUM (2-of-3 Multisig)"
	echo ""

	# Execute workflow phases
	bch_check_prerequisites "watch keygen sign1 sign2"
	bch_setup_infrastructure "bch-watch bch-keygen bch-sign1 bch-sign2"
	bch_setup_wallets "watch keygen sign1 sign2"
	key_generation_phase
	multisig_setup_phase
	generate_test_utxos
	create_payment_requests_phase
	transaction_flow_phase

	log_step "Bitcoin Cash E2E Workflow Completed Successfully!"
	log_info "Pattern 2: P2SH 2-of-3 multisig transaction completed"
	log_info "Summary:"
	log_info "  - Infrastructure setup complete"
	log_info "  - Wallets created and configured"
	log_info "  - Keys generated and imported"
	log_info "  - 2-of-3 multisig addresses created"
	log_info "  - Addresses imported into watch wallet"
	log_info "  - Test UTXOs generated"
	log_info "  - Transaction created, signed (2 of 3), and sent"
}

# Trap errors and cleanup
trap 'log_error "Script failed at line $LINENO"' ERR

# Run main
main "$@"
