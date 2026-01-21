#!/usr/bin/env bash

# Ethereum E2E Workflow Script - Pattern 1: Anvil Basic
# This script automates the complete Ethereum workflow with Anvil (Foundry's local node)
# Usage: ./scripts/operation/eth/e2e/e2e-p1-anvil-basic.sh [OPTIONS]
# Options:
#   --cleanup  Stop containers and cleanup state
#   --reset    Full reset and run from scratch
#   --verbose  Enable verbose output
#   --non-interactive  Run without prompts (for CI/CD)
#   -h, --help Display help message
#
# Reference Documentation:
#   docs/crypto/eth/operations/e2e-transaction-patterns.md - E2E transaction patterns
#
# Transaction Pattern:
#   Pattern 1: ETH Basic Single-sig
#   - Address Type: Standard Ethereum (secp256k1)
#   - Address Format: `0x...` (checksummed)
#   - Transaction Type: Legacy or EIP-1559 (auto-detected)
#
# Required Config Settings:
#   - config/wallet/eth/watch.yaml:  Port 8546 for Anvil
#   - config/wallet/eth/keygen.yaml: Default Ethereum settings

set -euo pipefail

###############################################################################
# Script Configuration
###############################################################################

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the ETH common utilities (which also sources common.sh)
# shellcheck source=../eth_common.sh
source "${SCRIPT_DIR}/../eth_common.sh"

# Pattern identifier for database isolation
export E2E_PATTERN="p1"

###############################################################################
# Environment Variable Overrides for Configuration
###############################################################################
# These environment variables override config file values.
# Priority: Environment Variables > Config File > Default Values
#
# Pattern 1 requires:
#   - network_type: "anvil" (for Anvil compatibility)
#   - port: 8546 (Anvil default)
export WALLET_ETHEREUM_NETWORK_TYPE="anvil"
export WALLET_ETHEREUM_PORT="8546"

###############################################################################
# Configuration
###############################################################################

# Get configuration paths
eth_get_config_paths

# Account configuration (single-sig)
ACCOUNT_CONFIG="${SCRIPT_DIR}/../../../config/wallet/account/account.yaml"

# Test addresses and amounts
FUNDING_AMOUNT_ETH=100
TRANSFER_AMOUNT_ETH=1
DEPOSIT_ACCOUNT="deposit"
CLIENT_ACCOUNT="client"

###############################################################################
# Parse Arguments
###############################################################################

MODE="run"
VERBOSE=false
NON_INTERACTIVE=false

show_help() {
	cat <<EOF
Ethereum E2E Workflow Script - Pattern 1: Anvil Basic

Usage: $0 [OPTIONS]

Options:
  --cleanup           Stop containers and cleanup state
  --reset             Full reset and run from scratch
  --verbose           Enable verbose output
  --non-interactive   Run without prompts (for CI/CD)
  -h, --help          Display this help message

Examples:
  $0                  # Run E2E workflow
  $0 --reset          # Fresh start with full reset
  $0 --verbose        # Run with detailed logging
  $0 --cleanup        # Stop containers and cleanup

Environment Variables:
  DB_TYPE             Database type: sqlite (default) or mysql
  ETH_RPC_HOST        Anvil RPC host (default: 127.0.0.1)
  ETH_RPC_PORT        Anvil RPC port (default: 8546)

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
		set -x # Enable bash debug mode
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
# Utility Functions
###############################################################################

# Extract address from wallet command output
extract_address() {
	local output="$1"
	echo "$output" | grep -oE '0x[a-fA-F0-9]{40}'
}

# Extract single address (first match)
extract_single_address() {
	local output="$1"
	echo "$output" | grep -oE '0x[a-fA-F0-9]{40}' | head -1
}

# Wait for transaction confirmation
wait_for_confirmation() {
	local tx_hash="$1"
	local max_wait=30
	local count=0

	log_substep "Waiting for transaction confirmation: ${tx_hash}"

	while [ $count -lt $max_wait ]; do
		local receipt
		receipt=$(curl -s -X POST -H "Content-Type: application/json" \
			--data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionReceipt\",\"params\":[\"${tx_hash}\"],\"id\":1}" \
			"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}")

		if echo "$receipt" | grep -q '"status":"0x1"'; then
			log_info "Transaction confirmed successfully"
			return 0
		fi

		sleep 1
		count=$((count + 1))
	done

	log_error "Transaction not confirmed within ${max_wait} seconds"
	return 1
}

###############################################################################
# Workflow Functions
###############################################################################

e2e_keygen_workflow() {
	log_step "Key Generation Workflow"

	# 1. Create seed
	log_substep "Creating seed"
	eth_keygen_cmd -coin eth create seed

	# 2. Create HD wallet keys for deposit account
	log_substep "Creating HD keys for ${DEPOSIT_ACCOUNT} account"
	eth_keygen_cmd -coin eth create hdkey --account "${DEPOSIT_ACCOUNT}"

	# 3. Create HD wallet keys for client account
	log_substep "Creating HD keys for ${CLIENT_ACCOUNT} account"
	eth_keygen_cmd -coin eth create hdkey --account "${CLIENT_ACCOUNT}"

	# 4. Import private keys to local keystore
	log_substep "Importing private keys to keystore for ${DEPOSIT_ACCOUNT}"
	eth_keygen_cmd -coin eth import key --account "${DEPOSIT_ACCOUNT}"

	log_substep "Importing private keys to keystore for ${CLIENT_ACCOUNT}"
	eth_keygen_cmd -coin eth import key --account "${CLIENT_ACCOUNT}"

	log_info "Key generation completed"
}

e2e_address_workflow() {
	log_step "Address Creation Workflow"

	# 1. Create deposit addresses
	log_substep "Creating deposit addresses"
	local deposit_output
	deposit_output=$(eth_watch_cmd -coin eth api address create --account "${DEPOSIT_ACCOUNT}" --count 2)

	# Extract all addresses then select specific ones
	local addresses
	addresses=$(extract_address "$deposit_output")
	DEPOSIT_ADDR_1=$(echo "$addresses" | sed -n '1p')
	DEPOSIT_ADDR_2=$(echo "$addresses" | sed -n '2p')

	if [ -z "$DEPOSIT_ADDR_1" ] || [ -z "$DEPOSIT_ADDR_2" ]; then
		log_error "Failed to extract deposit addresses"
		return 1
	fi

	log_info "Deposit address 1: ${DEPOSIT_ADDR_1}"
	log_info "Deposit address 2: ${DEPOSIT_ADDR_2}"

	# 2. Create client addresses
	log_substep "Creating client addresses"
	local client_output
	client_output=$(eth_watch_cmd -coin eth api address create --account "${CLIENT_ACCOUNT}" --count 1)

	CLIENT_ADDR=$(extract_single_address "$client_output")

	if [ -z "$CLIENT_ADDR" ]; then
		log_error "Failed to extract client address"
		return 1
	fi

	log_info "Client address: ${CLIENT_ADDR}"

	log_info "Address creation completed"
}

e2e_funding_workflow() {
	log_step "Funding Workflow"

	# Fund deposit addresses using Anvil's anvil_setBalance
	eth_fund_address "${DEPOSIT_ADDR_1}" "${FUNDING_AMOUNT_ETH}"
	eth_fund_address "${DEPOSIT_ADDR_2}" "${FUNDING_AMOUNT_ETH}"

	# Verify balances
	log_substep "Verifying balances"

	local balance1
	balance1=$(curl -s -X POST -H "Content-Type: application/json" \
		--data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"${DEPOSIT_ADDR_1}\",\"latest\"],\"id\":1}" \
		"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" | grep -oE '"result":"0x[0-9a-f]+"' | cut -d'"' -f4)

	local balance2
	balance2=$(curl -s -X POST -H "Content-Type: application/json" \
		--data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"${DEPOSIT_ADDR_2}\",\"latest\"],\"id\":1}" \
		"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" | grep -oE '"result":"0x[0-9a-f]+"' | cut -d'"' -f4)

	log_info "Deposit address 1 balance: ${balance1}"
	log_info "Deposit address 2 balance: ${balance2}"

	log_info "Funding completed"
}

e2e_transaction_workflow() {
	log_step "Transaction Workflow"

	# 1. Create payment request
	log_substep "Creating payment request"
	local payment_output
	payment_output=$(eth_watch_cmd -coin eth api payment create \
		--receiver "${CLIENT_ADDR}" \
		--amount "${TRANSFER_AMOUNT_ETH}")

	log_info "Payment request created"

	# 2. Create unsigned transaction
	log_substep "Creating unsigned transaction"
	local tx_output
	tx_output=$(eth_sign_cmd -coin eth api tx create payment 2>&1) || {
		log_error "Failed to create payment transaction"
		log_error "Error details: $tx_output"
		return 1
	}

	# Extract transaction file path
	local tx_file
	tx_file=$(echo "$tx_output" | grep -oE '(data/tx|\.)/[^ ]+\.(json|hex)' | tail -1)

	if [ -z "$tx_file" ]; then
		log_error "Failed to extract transaction file path"
		log_error "Output was: $tx_output"
		return 1
	fi

	log_info "Transaction created: ${tx_file}"

	# 3. Sign transaction
	log_substep "Signing transaction"
	local signed_output
	signed_output=$(eth_sign_cmd -coin eth api tx sign --file "${tx_file}" 2>&1) || {
		log_error "Failed to sign transaction"
		log_error "Error details: $signed_output"
		return 1
	}

	# Extract signed transaction file
	local signed_file
	signed_file=$(echo "$signed_output" | grep -oE '(data/tx|\.)/[^ ]+\.(json|hex)' | tail -1)

	if [ -z "$signed_file" ]; then
		log_error "Failed to extract signed transaction file path"
		log_error "Output was: $signed_output"
		return 1
	fi

	log_info "Transaction signed: ${signed_file}"

	# 4. Send transaction
	log_substep "Sending transaction to Anvil"
	local send_output
	send_output=$(eth_watch_cmd -coin eth api tx send --file "${signed_file}" 2>&1) || {
		log_error "Failed to send transaction"
		log_error "Error details: $send_output"
		return 1
	}

	# Extract transaction hash
	local tx_hash
	tx_hash=$(echo "$send_output" | grep -oE '0x[a-fA-F0-9]{64}')

	if [ -z "$tx_hash" ]; then
		log_error "Failed to extract transaction hash"
		log_error "Output was: $send_output"
		return 1
	fi

	log_info "Transaction sent: ${tx_hash}"

	# 5. Wait for confirmation
	wait_for_confirmation "${tx_hash}"

	# 6. Verify client address balance
	log_substep "Verifying client address balance"
	local client_balance
	client_balance=$(curl -s -X POST -H "Content-Type: application/json" \
		--data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"${CLIENT_ADDR}\",\"latest\"],\"id\":1}" \
		"http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" | grep -oE '"result":"0x[0-9a-f]+"' | cut -d'"' -f4)

	log_info "Client address balance: ${client_balance}"

	log_info "Transaction workflow completed"
}

e2e_full_workflow() {
	log_header "Ethereum E2E Workflow - Pattern 1: Anvil Basic"

	# Prerequisites
	eth_check_prerequisites

	# Setup infrastructure
	eth_setup_infrastructure

	# Run workflows
	e2e_keygen_workflow
	e2e_address_workflow
	e2e_funding_workflow
	e2e_transaction_workflow

	log_success "E2E workflow completed successfully"
}

###############################################################################
# Main Execution
###############################################################################

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

# Execute main function
main
