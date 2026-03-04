#!/usr/bin/env bash

# XRP E2E Workflow Script - Pattern 1: Single-sig Payment Transfer
# This script automates the complete XRP single-sig transaction flow:
#   keygen seed/key → address export → watch address import →
#   fund addresses (genesis) → create unsigned transfer tx →
#   keygen offline sign → watch send → ledger_accept → monitor confirmation
#
# Usage: ./scripts/operation/xrp/e2e/e2e-p1.sh [OPTIONS]
#
# Options:
#   --cleanup           Stop containers and cleanup state
#   --reset             Full reset and run from scratch
#   --verbose           Enable verbose output (set -x)
#   --non-interactive   Run without prompts (for CI/CD)
#   -h, --help          Display this help message
#
# Transaction Pattern:
#   Pattern 1: XRP single-sig transfer (payment → deposit)
#   - Address Type: XRP classic address (ed25519 or secp256k1)
#   - Address Format: r... (base58)
#   - Signing: Keygen wallet (offline, no Sign wallet needed)
#
# Environment Variables:
#   DB_TYPE       Database type: sqlite (default) or mysql
#   XRP_WS_HOST   rippled WebSocket host (default: 127.0.0.1)
#   XRP_WS_PORT   rippled WebSocket port (default: 6006)
#
# Required Config:
#   config/wallet/xrp/watch.yaml:  ripple.websocket_public_url overridden by XRP_WS_HOST:XRP_WS_PORT
#   config/wallet/xrp/keygen.yaml: ripple.offline_keygen: true
#
# Infrastructure:
#   rippled v3.1.0 standalone mode (compose.xrp.yaml)
#   WebSocket admin port 6006 (ws://) — direct connection, no xrpl-grpc-server
#   Ledger must be manually advanced via 'rippled ledger_accept' after each tx

set -euo pipefail

###############################################################################
# Script Configuration
###############################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the XRP common utilities (which also sources common.sh)
# shellcheck source=../xrp_common.sh
source "${SCRIPT_DIR}/../xrp_common.sh"

# Pattern identifier for DB isolation
export E2E_PATTERN="p1"

# Re-initialize database with correct pattern now that E2E_PATTERN is set
xrp_init_database

###############################################################################
# Configuration
###############################################################################

xrp_get_config_paths

# Number of HD keys to generate per account
KEY_NUM=5

# Accounts used in this E2E test
PAYMENT_ACCOUNT="payment"
CLIENT_ACCOUNT="deposit"

# XRP amounts
FUNDING_AMOUNT_PAYMENT_XRP=100 # Fund payment account (sender)
FUNDING_AMOUNT_CLIENT_XRP=10   # Fund deposit account (activate with base reserve)
TRANSFER_AMOUNT_XRP=50.0       # Transfer amount

###############################################################################
# Argument Parsing
###############################################################################

MODE="run"
VERBOSE=false
NON_INTERACTIVE=false

show_help() {
	cat <<EOF
XRP E2E Workflow Script - Pattern 1: Single-sig Payment Transfer

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
  DB_TYPE       sqlite (default) | mysql
  XRP_WS_HOST   rippled WebSocket host (default: 127.0.0.1)
  XRP_WS_PORT   rippled WebSocket port (default: 6006)

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
	xrp_keygen_cmd -c "${XRP_CONFIG_KEYGEN}" --coin "${XRP_COIN}" create seed || {
		log_warn "Seed already exists or error, continuing..."
	}

	log_substep "Creating ${KEY_NUM} HD keys for '${PAYMENT_ACCOUNT}' account"
	xrp_keygen_cmd -c "${XRP_CONFIG_KEYGEN}" --coin "${XRP_COIN}" create hdkey \
		--account "${PAYMENT_ACCOUNT}" \
		--keynum "${KEY_NUM}"

	log_substep "Creating ${KEY_NUM} HD keys for '${CLIENT_ACCOUNT}' account"
	xrp_keygen_cmd -c "${XRP_CONFIG_KEYGEN}" --coin "${XRP_COIN}" create hdkey \
		--account "${CLIENT_ACCOUNT}" \
		--keynum "${KEY_NUM}"

	log_info "Key generation completed"
}

###############################################################################
# Phase 2: Address Export → Watch Address Import
###############################################################################

address_setup_phase() {
	log_step "Address Export (Keygen DB) → Watch Address Import Phase"

	for account in "${PAYMENT_ACCOUNT}" "${CLIENT_ACCOUNT}"; do
		log_substep "Exporting '${account}' addresses from keygen"
		local export_output
		export_output=$(xrp_keygen_cmd -c "${XRP_CONFIG_KEYGEN}" --coin "${XRP_COIN}" export address \
			--account "${account}" 2>&1) || {
			log_error "Failed to export addresses for '${account}'"
			log_error "Output: ${export_output}"
			return 1
		}

		local csv_file
		csv_file=$(xrp_extract_file_path "${export_output}")

		if [ -z "${csv_file}" ] || [ ! -f "${csv_file}" ]; then
			log_error "Failed to get address CSV for '${account}'"
			log_error "Output: ${export_output}"
			return 1
		fi
		log_info "  ${account} address CSV: ${csv_file}"

		log_substep "Importing '${account}' addresses into watch wallet"
		xrp_watch_cmd -c "${XRP_CONFIG_WATCH}" --coin "${XRP_COIN}" import address --file "${csv_file}"
	done

	log_info "Address setup completed"
}

###############################################################################
# Phase 3: Fund Addresses (from genesis account)
###############################################################################

funding_phase() {
	log_step "Funding Phase (from genesis account)"

	# Get payment address (sender)
	local payment_addr
	payment_addr=$(xrp_get_address "${PAYMENT_ACCOUNT}")

	if [ -z "${payment_addr}" ]; then
		log_error "No '${PAYMENT_ACCOUNT}' address found in watch DB"
		log_error "Ensure Phase 2 (address import) ran successfully"
		return 1
	fi

	# Get deposit address (receiver — needs activation with base reserve)
	local client_addr
	client_addr=$(xrp_get_address "${CLIENT_ACCOUNT}")

	if [ -z "${client_addr}" ]; then
		log_error "No '${CLIENT_ACCOUNT}' address found in watch DB"
		return 1
	fi

	log_info "Payment address (sender):   ${payment_addr}"
	log_info "Deposit address (receiver): ${client_addr}"

	# Fund payment account
	xrp_fund_address "${payment_addr}" "${FUNDING_AMOUNT_PAYMENT_XRP}"

	# Fund deposit account with base reserve to activate it
	xrp_fund_address "${client_addr}" "${FUNDING_AMOUNT_CLIENT_XRP}"

	log_info "Funding completed"
}

###############################################################################
# Phase 4: Transaction Creation (Watch wallet — online)
###############################################################################

# Globals set by create_tx_phase / sign_tx_phase
UNSIGNED_TX_FILE=""
SIGNED_TX_FILE=""

create_tx_phase() {
	log_step "Transaction Creation Phase (Watch wallet)"

	log_substep "Creating unsigned transfer tx: ${PAYMENT_ACCOUNT} → ${CLIENT_ACCOUNT} (${TRANSFER_AMOUNT_XRP} XRP)"
	local create_output
	create_output=$(xrp_watch_cmd -c "${XRP_CONFIG_WATCH}" --coin "${XRP_COIN}" create transfer \
		--account1 "${PAYMENT_ACCOUNT}" \
		--account2 "${CLIENT_ACCOUNT}" \
		--amount "${TRANSFER_AMOUNT_XRP}" 2>&1) || {
		log_error "Failed to create transfer transaction"
		log_error "Output: ${create_output}"
		return 1
	}

	UNSIGNED_TX_FILE=$(xrp_extract_file_path "${create_output}")

	if [ -z "${UNSIGNED_TX_FILE}" ]; then
		log_error "Failed to extract unsigned tx file path"
		log_error "Output: ${create_output}"
		return 1
	fi

	log_info "Unsigned transaction file: ${UNSIGNED_TX_FILE}"
}

###############################################################################
# Phase 5: Offline Signing (Keygen wallet — offline, no network calls)
###############################################################################

sign_tx_phase() {
	log_step "Offline Signing Phase (Keygen wallet)"

	log_substep "Signing transaction: ${UNSIGNED_TX_FILE}"
	local sign_output
	sign_output=$(xrp_keygen_cmd -c "${XRP_CONFIG_KEYGEN}" --coin "${XRP_COIN}" sign signature \
		--file "${UNSIGNED_TX_FILE}" 2>&1) || {
		log_error "Failed to sign transaction"
		log_error "Output: ${sign_output}"
		return 1
	}

	SIGNED_TX_FILE=$(xrp_extract_file_path "${sign_output}")

	if [ -z "${SIGNED_TX_FILE}" ]; then
		log_error "Failed to extract signed tx file path"
		log_error "Output: ${sign_output}"
		return 1
	fi

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
	send_output=$(xrp_watch_cmd -c "${XRP_CONFIG_WATCH}" --coin "${XRP_COIN}" send tx \
		--file "${SIGNED_TX_FILE}" 2>&1) || {
		log_error "Failed to send transaction"
		log_error "Output: ${send_output}"
		return 1
	}

	local tx_id
	tx_id=$(xrp_extract_tx_id "${send_output}")

	if [ -z "${tx_id}" ]; then
		log_warn "Could not extract txID from output; the tx may still have been sent"
		log_info "Raw output: ${send_output}"
	else
		log_info "Transaction sent! txID: ${tx_id}"
	fi

	# Advance ledger to validate the transaction (required in standalone mode)
	log_substep "Advancing ledger to validate transaction"
	xrp_ledger_accept
}

###############################################################################
# Phase 7: Confirmation Monitoring (Watch wallet — online)
###############################################################################

monitor_phase() {
	log_step "Monitoring Phase (Watch wallet)"

	log_substep "Monitoring sent transactions for confirmation"
	xrp_watch_cmd -c "${XRP_CONFIG_WATCH}" --coin "${XRP_COIN}" monitor senttx || {
		log_warn "Monitor returned non-zero (may be no transactions to monitor yet)"
	}

	log_info "Monitoring completed"
}

###############################################################################
# Full E2E Workflow
###############################################################################

e2e_full_workflow() {
	log_step "XRP E2E Workflow - Pattern 1: Single-sig Payment Transfer"
	log_info "DB type: ${DB_TYPE}"
	log_info "WS:      ws://${XRP_WS_HOST}:${XRP_WS_PORT}"
	echo ""

	xrp_check_prerequisites
	xrp_setup_infrastructure
	xrp_init_sqlite_db

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
		xrp_cleanup
		;;
	reset)
		xrp_full_reset
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
