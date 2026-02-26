#!/usr/bin/env bash

# Ethereum E2E Workflow Script - Pattern 1: Single-sig EIP-1559
# This script automates the complete Ethereum single-sig transaction flow:
#   keygen seed/key → accountXpub export → watch address import →
#   fund addresses → create unsigned tx → keygen offline sign →
#   watch send → monitor confirmation
#
# Usage: ./scripts/operation/eth/e2e/e2e-p1.sh [OPTIONS]
#
# Options:
#   --cleanup           Stop containers and cleanup state
#   --reset             Full reset and run from scratch
#   --verbose           Enable verbose output (set -x)
#   --non-interactive   Run without prompts (for CI/CD)
#   -h, --help          Display this help message
#
# Reference Documentation:
#   docs/chains/eth/operations/e2e-transaction-patterns.md
#
# Transaction Pattern:
#   Pattern 1: ETH Single-sig (EOA)
#   - Address Type: Standard Ethereum secp256k1 EOA
#   - Address Format: 0x... (checksummed)
#   - Transaction Type: EIP-1559 (Type 2) with Anvil auto-detection
#   - Signing: Keygen wallet (offline HD derivation, no Sign wallet needed)
#
# Environment Variables:
#   NODE_TYPE    Node type: anvil (default) or geth
#   DB_TYPE      Database type: sqlite (default) or mysql
#   ETH_RPC_HOST Ethereum RPC host (default: 127.0.0.1)
#   ETH_RPC_PORT Ethereum RPC port (default: 8546 for anvil, 8545 for geth)
#
# Required Config:
#   config/wallet/eth/watch.yaml:  ethereum.port: 8546 (Anvil default)
#   config/wallet/eth/keygen.yaml: default Ethereum settings

set -euo pipefail

###############################################################################
# Script Configuration
###############################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the ETH common utilities (which also sources common.sh)
# shellcheck source=../eth_common.sh
source "${SCRIPT_DIR}/../eth_common.sh"

# Pattern identifier for DB isolation (also used by eth_get_db_path)
export E2E_PATTERN="p1"

# Re-initialize database with correct pattern now that E2E_PATTERN is set
eth_init_database

###############################################################################
# Environment Variable Overrides for Configuration
###############################################################################
# Priority: Environment Variables > Config File > Default Values
#
# Pattern 1 uses Anvil by default; override with NODE_TYPE=geth for Geth.
# ETH uses BIP44 derivation (eth-address -> bip44), required to initialize key generator.
export WALLET_ADDRESS_TYPE="eth-address"
# NODE_TYPE is the client type (anvil/geth); WALLET_ETHEREUM_NETWORK_TYPE is the network name.
# For geth --dev (E2E), use "local" (chain ID 1337, same as anvil); for anvil, use "anvil".
if [ "${NODE_TYPE}" = "geth" ]; then
	export WALLET_ETHEREUM_NETWORK_TYPE="local"
else
	export WALLET_ETHEREUM_NETWORK_TYPE="anvil"
fi
export WALLET_ETHEREUM_PORT="${ETH_RPC_PORT}"

###############################################################################
# Configuration
###############################################################################

eth_get_config_paths

# Number of HD keys to generate per account
KEY_NUM=5

# Accounts used in this E2E test
# createTransferTx restricts receiver to internal accounts (not client, not authorization)
PAYMENT_ACCOUNT="payment"
CLIENT_ACCOUNT="deposit"

# Amount (ETH) to fund sender address and transfer
FUNDING_AMOUNT_ETH=100
TRANSFER_AMOUNT_ETH=1.0

###############################################################################
# Argument Parsing
###############################################################################

MODE="run"
VERBOSE=false
NON_INTERACTIVE=false

show_help() {
	cat <<EOF
Ethereum E2E Workflow Script - Pattern 1: Single-sig EIP-1559

Usage: $0 [OPTIONS]

Options:
  --cleanup           Stop containers and cleanup state
  --reset             Full reset and run from scratch
  --verbose           Enable verbose output
  --non-interactive   Run without prompts (for CI/CD)
  -h, --help          Display this help message

Examples:
  $0               # Run E2E workflow
  $0 --reset       # Fresh start with full reset
  $0 --verbose     # Run with detailed logging
  $0 --cleanup     # Stop containers and cleanup

Environment Variables:
  NODE_TYPE    anvil (default) | geth
  DB_TYPE      sqlite (default) | mysql
  ETH_RPC_HOST Ethereum RPC host (default: 127.0.0.1)
  ETH_RPC_PORT Ethereum RPC port (default: auto from NODE_TYPE)

EOF
	exit 0
}

while [[ $# -gt 0 ]]; do
	case $1 in
	--cleanup)
		MODE="cleanup"
		shift
		;;
	--reset)
		MODE="reset"
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
		;;
	*)
		log_error "Unknown option: $1"
		show_help
		;;
	esac
done

###############################################################################
# Phase 1: Key Generation (Keygen wallet — offline)
###############################################################################

keygen_phase() {
	log_step "Key Generation Phase (Keygen wallet)"

	log_substep "Creating mnemonic seed"
	eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" create seed || {
		log_warn "Seed already exists or error, continuing..."
	}

	log_substep "Creating ${KEY_NUM} HD keys for '${PAYMENT_ACCOUNT}' account"
	eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" create hdkey \
		--account "${PAYMENT_ACCOUNT}" \
		--keynum "${KEY_NUM}"

	log_substep "Creating ${KEY_NUM} HD keys for '${CLIENT_ACCOUNT}' account"
	eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" create hdkey \
		--account "${CLIENT_ACCOUNT}" \
		--keynum "${KEY_NUM}"

	log_info "Key generation completed"
}

###############################################################################
# Phase 2: AccountXpub Export → Watch Address Import
###############################################################################

address_setup_phase() {
	log_step "Address Export (Keygen DB) → Watch Address Import Phase"

	# Export ETH addresses from keygen DB as a watch-compatible CSV and import into watch wallet.
	# watch import address uses ConvertAddressLine which expects BTC-style 8-field CSV;
	# eth_export_watch_address_csv creates that format using the keygen SQLite DB directly.
	for account in "${PAYMENT_ACCOUNT}" "${CLIENT_ACCOUNT}"; do
		log_substep "Exporting '${account}' addresses from keygen DB"
		local csv_file
		csv_file=$(eth_export_watch_address_csv "${account}")

		if [ -z "${csv_file}" ] || [ ! -f "${csv_file}" ] || [ ! -s "${csv_file}" ]; then
			log_error "Failed to create address CSV for '${account}'"
			return 1
		fi
		log_info "  ${account} address CSV: ${csv_file}"

		log_substep "Importing '${account}' addresses into watch wallet"
		eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${ETH_COIN}" import address --file "${csv_file}"
	done

	log_info "Address setup completed"
}

###############################################################################
# Phase 3: Fund Payment Addresses (Anvil only)
###############################################################################

funding_phase() {
	log_step "Funding Phase"

	# Retrieve first payment address from watch DB
	local payment_addr
	payment_addr=$(eth_get_payment_address "${PAYMENT_ACCOUNT}")

	if [ -z "${payment_addr}" ]; then
		log_error "No '${PAYMENT_ACCOUNT}' address found in watch DB"
		log_error "Ensure Phase 2 (address import) ran successfully"
		return 1
	fi

	log_info "Funding payment address: ${payment_addr}"
	eth_fund_address "${payment_addr}" "${FUNDING_AMOUNT_ETH}"
}

###############################################################################
# Phase 4: Transaction Creation (Watch wallet — online)
###############################################################################

create_tx_phase() {
	log_step "Transaction Creation Phase (Watch wallet)"

	log_substep "Creating unsigned transfer tx: ${PAYMENT_ACCOUNT} → ${CLIENT_ACCOUNT} (${TRANSFER_AMOUNT_ETH} ETH)"
	local create_output
	create_output=$(eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${ETH_COIN}" create transfer \
		--account1 "${PAYMENT_ACCOUNT}" \
		--account2 "${CLIENT_ACCOUNT}" \
		--amount "${TRANSFER_AMOUNT_ETH}" 2>&1) || {
		log_error "Failed to create transfer transaction"
		log_error "Output: ${create_output}"
		return 1
	}

	UNSIGNED_TX_FILE=$(eth_extract_file_path "${create_output}")

	if [ -z "${UNSIGNED_TX_FILE}" ]; then
		log_error "Failed to extract unsigned tx file path"
		log_error "Output: ${create_output}"
		return 1
	fi

	log_info "Unsigned transaction file: ${UNSIGNED_TX_FILE}"
}

# Declare global for signed tx path (set in create_tx_phase, read in sign_phase)
UNSIGNED_TX_FILE=""
SIGNED_TX_FILE=""

###############################################################################
# Phase 5: Offline Signing (Keygen wallet — offline, no network calls)
###############################################################################

sign_tx_phase() {
	log_step "Offline Signing Phase (Keygen wallet)"

	log_substep "Signing transaction: ${UNSIGNED_TX_FILE}"
	local sign_output
	sign_output=$(eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" sign signature \
		--file "${UNSIGNED_TX_FILE}" 2>&1) || {
		log_error "Failed to sign transaction"
		log_error "Output: ${sign_output}"
		return 1
	}

	SIGNED_TX_FILE=$(eth_extract_file_path "${sign_output}")

	if [ -z "${SIGNED_TX_FILE}" ]; then
		log_error "Failed to extract signed tx file path"
		log_error "Output: ${sign_output}"
		return 1
	fi

	# Verify signing is complete
	local is_completed
	is_completed=$(echo "${sign_output}" | grep '^\[isCompleted\]:' | sed 's/^\[isCompleted\]: //')
	if [ "${is_completed}" != "true" ]; then
		log_warn "Signing reported isCompleted=${is_completed}"
	fi

	log_info "Signed transaction file: ${SIGNED_TX_FILE}"
}

###############################################################################
# Phase 6: Broadcast (Watch wallet — online)
###############################################################################

send_tx_phase() {
	log_step "Broadcast Phase (Watch wallet)"

	log_substep "Sending signed transaction: ${SIGNED_TX_FILE}"
	local send_output
	send_output=$(eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${ETH_COIN}" send tx \
		--file "${SIGNED_TX_FILE}" 2>&1) || {
		log_error "Failed to send transaction"
		log_error "Output: ${send_output}"
		return 1
	}

	local tx_id
	tx_id=$(eth_extract_tx_id "${send_output}")

	if [ -z "${tx_id}" ]; then
		log_warn "Could not extract txID from output; the tx may still have been sent"
		log_info "Raw output: ${send_output}"
	else
		log_info "Transaction sent! txID: ${tx_id}"
	fi
}

###############################################################################
# Phase 7: Confirmation Monitoring (Watch wallet — online)
###############################################################################

monitor_phase() {
	log_step "Monitoring Phase (Watch wallet)"

	log_substep "Monitoring sent transactions for confirmation"
	eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${ETH_COIN}" monitor senttx || {
		log_warn "Monitor returned non-zero (may be no transactions to monitor yet)"
	}

	log_info "Monitoring completed"
}

###############################################################################
# Full E2E Workflow
###############################################################################

e2e_full_workflow() {
	log_step "Ethereum E2E Workflow - Pattern 1: Single-sig EIP-1559"
	log_info "Node type: ${NODE_TYPE}"
	log_info "DB type:   ${DB_TYPE}"
	log_info "RPC:       ${ETH_RPC_HOST}:${ETH_RPC_PORT}"
	echo ""

	eth_check_prerequisites
	eth_setup_infrastructure
	eth_init_sqlite_db

	keygen_phase
	address_setup_phase
	funding_phase
	create_tx_phase
	sign_tx_phase
	send_tx_phase
	monitor_phase

	log_info "E2E Pattern 1 completed successfully!"
}

###############################################################################
# Main
###############################################################################

trap 'log_error "Script failed at line $LINENO"' ERR

main() {
	case "${MODE}" in
	cleanup)
		eth_cleanup
		;;
	reset)
		eth_full_reset
		e2e_full_workflow
		;;
	run)
		e2e_full_workflow
		;;
	*)
		log_error "Invalid mode: ${MODE}"
		exit 1
		;;
	esac
}

main
