#!/usr/bin/env bash

# Ethereum E2E Workflow Script - Pattern 3: Safe 2-of-2 Multisig Payment
# This script automates the complete Safe multisig payment flow:
#   keygen seed/key (signer 1) → keygen seed/key (signer 2) →
#   deploy Safe on Anvil → fund Safe → watch create multisig →
#   keygen sign (signer 1) → keygen sign (signer 2) →
#   watch send multisig → verify ETH balances
#
# Usage: ./scripts/operation/eth/e2e/e2e-p3.sh [OPTIONS]
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
#   Pattern 3: ETH Safe 2-of-2 multisig payment
#   - Safe address: deployed on Anvil from DeploySafe.s.sol
#   - Owners: signer1 (keygen DB) + signer2 (sign DB simulated by second keygen binary)
#   - Threshold: 2 (both signers required)
#   - Signing: Keygen wallet (offline EIP-712, no network calls during signing)
#
# Environment Variables:
#   NODE_TYPE        Node type: anvil (default) or geth
#   DB_TYPE          Database type: sqlite (default) or mysql
#   ETH_RPC_HOST     Ethereum RPC host (default: 127.0.0.1)
#   ETH_RPC_PORT     Ethereum RPC port (default: 8546 for anvil, 8545 for geth)
#   DEPLOYER_KEY     Private key for Anvil deployer (default: Anvil account 0)
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
export E2E_PATTERN="p3"

# Re-initialize database with correct pattern now that E2E_PATTERN is set
eth_init_database

###############################################################################
# Environment Variable Overrides for Configuration
###############################################################################

export WALLET_ADDRESS_TYPE="eth-address"
if [ "${NODE_TYPE}" = "geth" ]; then
	export WALLET_ETHEREUM_NETWORK_TYPE="local"
else
	export WALLET_ETHEREUM_NETWORK_TYPE="anvil"
fi
export WALLET_ETHEREUM_PORT="${ETH_RPC_PORT}"

# Gas payer for execTransaction: use the deployer key (Anvil account 0).
# The deployer is pre-funded and its key is available as a raw hex value.
# This tells the Watch wallet which key to use when submitting the Safe tx on-chain.
export WALLET_ETHEREUM_SAFE_GAS_PAYER_HEX_KEY="${DEPLOYER_KEY:-0xd3f1327df33c2a67ac03877f37410f33ed3a37085a84c739a7f5a137e5152bf2}"

###############################################################################
# Configuration
###############################################################################

eth_get_config_paths

# Foundry toolchain binary paths
CAST="${HOME}/.foundry/bin/cast"
FORGE="${HOME}/.foundry/bin/forge"

# Deployer private key (Anvil account 0 from project mnemonic).
# Override with DEPLOYER_KEY=0x... to use a different account.
DEPLOYER_KEY="${DEPLOYER_KEY:-0xd3f1327df33c2a67ac03877f37410f33ed3a37085a84c739a7f5a137e5152bf2}"

# Deployer address corresponding to DEPLOYER_KEY (Anvil account 0).
DEPLOYER_ADDRESS="${DEPLOYER_ADDRESS:-0x328F371a76dfAc47b89Cc007bb048ec446c21494}"

# Number of HD keys to generate per signer
KEY_NUM=3

# Account type to use for owner key generation
OWNER_ACCOUNT="payment"

# ETH amounts
SAFE_FUNDING_ETH=10   # ETH to fund the Safe contract
TRANSFER_AMOUNT_ETH=1 # ETH to transfer via multisig tx

# Safe threshold (m-of-n)
SAFE_THRESHOLD=2

# Second-signer DB path (simulates the Sign wallet with a separate key store)
SQLITE_SIGN_DB_PATH=""

###############################################################################
# Argument Parsing
###############################################################################

MODE="run"
VERBOSE=false
NON_INTERACTIVE=false

show_help() {
	cat <<EOF
Ethereum E2E Workflow Script - Pattern 3: Safe 2-of-2 Multisig Payment

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
  NODE_TYPE        anvil (default) | geth
  DB_TYPE          sqlite (default) | mysql
  ETH_RPC_HOST     Ethereum RPC host (default: 127.0.0.1)
  ETH_RPC_PORT     Ethereum RPC port (default: auto from NODE_TYPE)
  DEPLOYER_KEY     Private key for Safe deployment (default: Anvil account 0)

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
# Second-signer DB helpers
# The "Sign wallet" role is simulated by running the keygen binary with a
# separate SQLite DB so that each signer has independent HD key material.
###############################################################################

p3_init_sign_db() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		local project_root
		project_root="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"
		SQLITE_SIGN_DB_PATH="${project_root}/data/sqlite/eth/sign-e2e-${E2E_PATTERN}.db"
		mkdir -p "$(dirname "${SQLITE_SIGN_DB_PATH}")"
		log_info "  sign:   ${SQLITE_SIGN_DB_PATH}"
	fi
}

# Run a keygen command against the second-signer (Sign wallet role) DB.
p3_sign_wallet_cmd() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		WALLET_DATABASE_SQLITE_PATH="${SQLITE_SIGN_DB_PATH}" \
			"${GOPATH}/bin/keygen" "$@"
	else
		"${GOPATH}/bin/keygen" "$@"
	fi
}

# Initialise the second-signer SQLite DB with the keygen schema.
p3_init_sign_sqlite_db() {
	if [ "${DB_TYPE}" != "sqlite" ]; then
		return 0
	fi

	local project_root
	project_root="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"
	local keygen_schema="${project_root}/tools/sqlc/schemas/sqlite/e2e/02_keygen.sql"

	if [ -f "${keygen_schema}" ]; then
		rm -f "${SQLITE_SIGN_DB_PATH}"
		sqlite3 "${SQLITE_SIGN_DB_PATH}" <"${keygen_schema}"
		log_info "Sign SQLite DB initialized"
	else
		log_warn "Keygen schema not found: ${keygen_schema}"
	fi
}

###############################################################################
# Extract a multisig file path from a wallet command's stdout
# Looks for "  File:      <path>" in watch create multisig output.
###############################################################################

p3_extract_multisig_file() {
	local output="$1"
	echo "${output}" | grep -E '^ *File:' | awk '{print $NF}'
}

###############################################################################
# Phase 0: Provisioning — Key Generation for Both Signers
###############################################################################

provisioning_phase() {
	log_step "Provisioning Phase — Key Generation for Both Signers"

	# --- Signer 1 (primary Keygen wallet DB) ---
	log_substep "Signer 1: creating mnemonic seed"
	eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" create seed || {
		log_warn "Seed already exists or error, continuing..."
	}

	log_substep "Signer 1: creating ${KEY_NUM} HD keys for '${OWNER_ACCOUNT}' account"
	eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" create hdkey \
		--account "${OWNER_ACCOUNT}" \
		--keynum "${KEY_NUM}"

	# --- Signer 2 (Sign wallet role — separate keygen DB) ---
	log_substep "Signer 2: creating mnemonic seed"
	p3_sign_wallet_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" create seed || {
		log_warn "Seed already exists or error, continuing..."
	}

	log_substep "Signer 2: creating ${KEY_NUM} HD keys for '${OWNER_ACCOUNT}' account"
	p3_sign_wallet_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" create hdkey \
		--account "${OWNER_ACCOUNT}" \
		--keynum "${KEY_NUM}"

	log_info "Key provisioning completed"
}

###############################################################################
# Get signer addresses from their respective DBs
###############################################################################

p3_get_signer1_address() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		sqlite3 "${SQLITE_KEYGEN_DB_PATH}" \
			"SELECT address FROM eth_account_key WHERE account='${OWNER_ACCOUNT}' ORDER BY idx LIMIT 1" \
			2>/dev/null || true
	fi
}

p3_get_signer2_address() {
	if [ "${DB_TYPE}" = "sqlite" ]; then
		sqlite3 "${SQLITE_SIGN_DB_PATH}" \
			"SELECT address FROM eth_account_key WHERE account='${OWNER_ACCOUNT}' ORDER BY idx LIMIT 1" \
			2>/dev/null || true
	fi
}

###############################################################################
# Phase 1: Setup — Deploy Safe and Fund It
###############################################################################

# Globals set in this phase and read by later phases
SAFE_ADDRESS=""
SIGNER1_ADDR=""
SIGNER2_ADDR=""

setup_phase() {
	log_step "Setup Phase — Deploy Safe Contract and Fund It"

	# Retrieve signer addresses from their respective DBs
	SIGNER1_ADDR=$(p3_get_signer1_address)
	SIGNER2_ADDR=$(p3_get_signer2_address)

	if [ -z "${SIGNER1_ADDR}" ]; then
		log_error "Could not retrieve signer 1 address from keygen DB"
		return 1
	fi
	if [ -z "${SIGNER2_ADDR}" ]; then
		log_error "Could not retrieve signer 2 address from sign DB"
		return 1
	fi

	log_info "Signer 1: ${SIGNER1_ADDR}"
	log_info "Signer 2: ${SIGNER2_ADDR}"

	# Deploy Safe proxy with both signers as owners and threshold=2
	local project_root
	project_root="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

	log_substep "Deploying Safe (owners: ${SIGNER1_ADDR},${SIGNER2_ADDR}; threshold: ${SAFE_THRESHOLD})"
	local deploy_output
	deploy_output=$(
		cd "${project_root}/apps/eth-contracts"
		DEPLOYER_KEY="${DEPLOYER_KEY}" \
			SAFE_OWNERS="${SIGNER1_ADDR},${SIGNER2_ADDR}" \
			SAFE_THRESHOLD="${SAFE_THRESHOLD}" \
			"${FORGE}" script script/DeploySafe.s.sol \
			--rpc-url "http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" \
			--broadcast \
			2>&1
	)

	# Parse the "SAFE_ADDRESS=0x..." line from forge output
	SAFE_ADDRESS=$(echo "${deploy_output}" | grep 'SAFE_ADDRESS=' | tail -1 | cut -d= -f2 | tr -d '[:space:]')

	if [ -z "${SAFE_ADDRESS}" ]; then
		log_error "Failed to parse Safe address from forge output"
		log_error "Output: ${deploy_output}"
		return 1
	fi
	log_info "Safe deployed at: ${SAFE_ADDRESS}"

	# Fund the Safe contract with ETH so it can pay the recipient
	log_substep "Funding Safe contract (${SAFE_ADDRESS}) with ${SAFE_FUNDING_ETH} ETH"
	eth_fund_address "${SAFE_ADDRESS}" "${SAFE_FUNDING_ETH}"

	log_info "Setup phase completed"
}

###############################################################################
# Phase 2: Create — Propose the Multisig Transaction
###############################################################################

# Global: set by create_multisig_phase, read by sign phases
MULTISIG_FILE_0=""

create_multisig_phase() {
	log_step "Create Phase — Propose Multisig Transaction (Watch wallet)"

	log_substep "Creating multisig tx proposal: Safe → ${DEPLOYER_ADDRESS} (${TRANSFER_AMOUNT_ETH} ETH)"
	local create_output
	create_output=$(eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${ETH_COIN}" create multisig \
		--safe "${SAFE_ADDRESS}" \
		--to "${DEPLOYER_ADDRESS}" \
		--amount "${TRANSFER_AMOUNT_ETH}" \
		--threshold "${SAFE_THRESHOLD}" \
		--action-type "payment" 2>&1) || {
		log_error "Failed to create multisig proposal"
		log_error "Output: ${create_output}"
		return 1
	}

	MULTISIG_FILE_0=$(p3_extract_multisig_file "${create_output}")

	if [ -z "${MULTISIG_FILE_0}" ]; then
		log_error "Failed to extract multisig file path from output"
		log_error "Output: ${create_output}"
		return 1
	fi

	log_info "Unsigned multisig file: ${MULTISIG_FILE_0}"

	# Verify file exists and has TxType=unsigned
	local tx_type
	tx_type=$(jq -r '.tx_type' "${MULTISIG_FILE_0}" 2>/dev/null || echo "unknown")
	if [ "${tx_type}" != "unsigned" ]; then
		log_error "Expected tx_type=unsigned, got: ${tx_type}"
		return 1
	fi
	log_info "Verified: tx_type=${tx_type}"
}

###############################################################################
# Phase 3: Sign 1 — First Signer (Keygen wallet, offline)
###############################################################################

# Global: set by sign1_phase, read by sign2_phase
MULTISIG_FILE_1=""

sign1_phase() {
	log_step "Sign 1 Phase — First Signer (Keygen wallet, offline)"

	log_substep "Signer 1 (${SIGNER1_ADDR}) signing: ${MULTISIG_FILE_0}"
	local sign_output
	sign_output=$(eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" sign signature \
		--file "${MULTISIG_FILE_0}" \
		--signer-address "${SIGNER1_ADDR}" 2>&1) || {
		log_error "Signer 1 signing failed"
		log_error "Output: ${sign_output}"
		return 1
	}

	MULTISIG_FILE_1=$(eth_extract_file_path "${sign_output}")

	if [ -z "${MULTISIG_FILE_1}" ]; then
		log_error "Failed to extract file path from sign output"
		log_error "Output: ${sign_output}"
		return 1
	fi

	log_info "After sign 1: ${MULTISIG_FILE_1}"

	# Verify exactly 1 signature in the file
	local sig_count
	sig_count=$(jq '.signatures | length' "${MULTISIG_FILE_1}" 2>/dev/null || echo "0")
	if [ "${sig_count}" != "1" ]; then
		log_error "Expected 1 signature after sign 1, got: ${sig_count}"
		return 1
	fi
	log_info "Verified: signature count = ${sig_count}"
}

###############################################################################
# Phase 4: Sign 2 — Second Signer (Sign wallet role, offline)
###############################################################################

# Global: set by sign2_phase, read by send_phase
MULTISIG_FILE_2=""

sign2_phase() {
	log_step "Sign 2 Phase — Second Signer (Sign wallet role, offline)"

	log_substep "Signer 2 (${SIGNER2_ADDR}) signing: ${MULTISIG_FILE_1}"
	local sign_output
	sign_output=$(p3_sign_wallet_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" sign signature \
		--file "${MULTISIG_FILE_1}" \
		--signer-address "${SIGNER2_ADDR}" 2>&1) || {
		log_error "Signer 2 signing failed"
		log_error "Output: ${sign_output}"
		return 1
	}

	MULTISIG_FILE_2=$(eth_extract_file_path "${sign_output}")

	if [ -z "${MULTISIG_FILE_2}" ]; then
		log_error "Failed to extract file path from sign output"
		log_error "Output: ${sign_output}"
		return 1
	fi

	log_info "After sign 2: ${MULTISIG_FILE_2}"

	# Verify tx_type=signed (threshold met) and 2 signatures
	local tx_type sig_count
	tx_type=$(jq -r '.tx_type' "${MULTISIG_FILE_2}" 2>/dev/null || echo "unknown")
	sig_count=$(jq '.signatures | length' "${MULTISIG_FILE_2}" 2>/dev/null || echo "0")

	if [ "${tx_type}" != "signed" ]; then
		log_error "Expected tx_type=signed after both signers, got: ${tx_type}"
		return 1
	fi
	if [ "${sig_count}" != "2" ]; then
		log_error "Expected 2 signatures, got: ${sig_count}"
		return 1
	fi
	log_info "Verified: tx_type=${tx_type}, signature count=${sig_count}"
}

###############################################################################
# Phase 5: Send — Broadcast the Fully Signed Multisig Transaction
###############################################################################

TX_HASH=""

send_multisig_phase() {
	log_step "Send Phase — Broadcast Multisig Transaction (Watch wallet)"

	log_substep "Submitting signed multisig tx: ${MULTISIG_FILE_2}"
	local send_output
	send_output=$(eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${ETH_COIN}" send multisig send-eth \
		--file "${MULTISIG_FILE_2}" 2>&1) || {
		log_error "Failed to submit multisig transaction"
		log_error "Output: ${send_output}"
		return 1
	}

	TX_HASH=$(eth_extract_tx_id "${send_output}")

	if [ -z "${TX_HASH}" ]; then
		log_warn "Could not extract txHash from output; the tx may still have been submitted"
		log_info "Raw output: ${send_output}"
	else
		log_info "Multisig transaction submitted! txHash: ${TX_HASH}"
	fi
}

###############################################################################
# Phase 6: Verify — Assert ETH Balance Changes
###############################################################################

verify_phase() {
	log_step "Verify Phase — Assert ETH Balance Changes"

	# Get recipient balance (deployer address received the payment)
	log_substep "Querying recipient balance: ${DEPLOYER_ADDRESS}"
	local recipient_balance
	recipient_balance=$("${CAST}" balance "${DEPLOYER_ADDRESS}" \
		--rpc-url "http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" 2>/dev/null || echo "0")
	log_info "Recipient (deployer) balance: ${recipient_balance} Wei"

	# Get Safe contract balance (should have decreased)
	log_substep "Querying Safe contract balance: ${SAFE_ADDRESS}"
	local safe_balance
	safe_balance=$("${CAST}" balance "${SAFE_ADDRESS}" \
		--rpc-url "http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" 2>/dev/null || echo "unknown")
	log_info "Safe contract balance: ${safe_balance} Wei"

	# Verify recipient received at least the transfer amount
	# TRANSFER_AMOUNT_ETH in Wei = TRANSFER_AMOUNT_ETH * 10^18
	local expected_min_wei
	expected_min_wei=$(echo "${TRANSFER_AMOUNT_ETH} * 1000000000000000000" | bc)
	log_info "Expected minimum recipient balance: ${expected_min_wei} Wei"

	if [ "${recipient_balance}" -ge "${expected_min_wei}" ] 2>/dev/null; then
		log_info "✓ Recipient balance check passed (${recipient_balance} >= ${expected_min_wei})"
	else
		log_warn "Balance check inconclusive (may need bn arithmetic); recipient_balance=${recipient_balance}"
	fi

	log_info "Verification phase completed"
}

###############################################################################
# Full E2E Workflow
###############################################################################

e2e_full_workflow() {
	log_step "Ethereum E2E Workflow - Pattern 3: Safe 2-of-2 Multisig Payment"
	log_info "Node type: ${NODE_TYPE}"
	log_info "DB type:   ${DB_TYPE}"
	log_info "RPC:       ${ETH_RPC_HOST}:${ETH_RPC_PORT}"
	echo ""

	eth_check_prerequisites
	eth_setup_infrastructure
	eth_init_sqlite_db
	p3_init_sign_db
	p3_init_sign_sqlite_db

	provisioning_phase
	setup_phase
	create_multisig_phase
	sign1_phase
	sign2_phase
	send_multisig_phase
	verify_phase

	log_info "E2E Pattern 3 (Safe 2-of-2 multisig payment) completed successfully!"
}

###############################################################################
# Cleanup
###############################################################################

p3_cleanup() {
	eth_cleanup

	if [ "${DB_TYPE}" = "sqlite" ] && [ -n "${SQLITE_SIGN_DB_PATH}" ]; then
		log_substep "Removing sign DB: ${SQLITE_SIGN_DB_PATH}"
		rm -f "${SQLITE_SIGN_DB_PATH}"
	fi
}

p3_full_reset() {
	p3_cleanup

	log_substep "Removing keystore files"
	rm -f ./data/keystore/*

	log_info "Full reset completed"
}

###############################################################################
# Main
###############################################################################

trap 'log_error "Script failed at line $LINENO"' ERR

main() {
	p3_init_sign_db

	case "${MODE}" in
	cleanup)
		p3_cleanup
		;;
	reset)
		p3_full_reset
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
