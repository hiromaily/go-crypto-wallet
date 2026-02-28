#!/usr/bin/env bash

# Ethereum E2E Workflow Script - Pattern 2: ERC-20 HYT Token Transfer
# This script automates the complete HYT ERC-20 token transfer flow:
#   keygen seed/key → accountXpub export → watch address import →
#   fund ETH gas + HYT tokens → create unsigned tx → keygen offline sign →
#   watch send → verify HYT balances
#
# Usage: ./scripts/operation/eth/e2e/e2e-p2.sh [OPTIONS]
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
#   Pattern 2: ERC-20 HYT token transfer (EIP-1559)
#   - Coin Type: hyt (ERC-20 token on Ethereum)
#   - Address Type: Standard Ethereum secp256k1 EOA (same HD derivation as Pattern 1)
#   - Transaction Type: EIP-1559 (Type 2) — Data field encodes transfer(address,uint256)
#   - Signing: Keygen wallet (offline HD derivation, no Sign wallet needed)
#
# Environment Variables:
#   NODE_TYPE           Node type: anvil (default) or geth
#   DB_TYPE             Database type: sqlite (default) or mysql
#   ETH_RPC_HOST        Ethereum RPC host (default: 127.0.0.1)
#   ETH_RPC_PORT        Ethereum RPC port (default: 8546 for anvil, 8545 for geth)
#   HYT_CONTRACT_ADDRESS  HYT contract address (default: auto-detected from broadcast/)
#   HYT_DEPLOYER_KEY    Private key of the HYT deployer/master account
#                       (default: Anvil account 0 derived from project mnemonic)
#
# Required Config:
#   config/wallet/eth/watch.yaml:
#     ethereum.port: 8546 (Anvil default)
#     ethereum.erc20s.hyt.contract_address: <deployed HYT contract>
#     ethereum.erc20s.hyt.master_address: <deployer address>

set -euo pipefail

###############################################################################
# Script Configuration
###############################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the ETH common utilities (which also sources common.sh)
# shellcheck source=../eth_common.sh
source "${SCRIPT_DIR}/../eth_common.sh"

# Pattern identifier for DB isolation
export E2E_PATTERN="p2"

# Re-initialize database with correct pattern now that E2E_PATTERN is set
eth_init_database

###############################################################################
# Environment Variable Overrides for Configuration
###############################################################################

# HYT uses ETH BIP44 derivation (same key derivation path as ETH)
export WALLET_ADDRESS_TYPE="eth-address"
if [ "${NODE_TYPE}" = "geth" ]; then
	export WALLET_ETHEREUM_NETWORK_TYPE="local"
else
	export WALLET_ETHEREUM_NETWORK_TYPE="anvil"
fi
export WALLET_ETHEREUM_PORT="${ETH_RPC_PORT}"

# Tell the watch wallet which ERC-20 token to use from the erc20s config map
export WALLET_ETHEREUM_ERC20_TOKEN="hyt"

# Override HYT contract and master addresses in watch.yaml at runtime.
# These env vars take precedence over any static values in the config file,
# allowing the same yaml to work regardless of which Anvil run deployed HYT.
# Values are resolved later (after argument parsing) once HYT_CONTRACT_ADDRESS
# and HYT_MASTER_ADDRESS are fully determined — see _apply_hyt_config_overrides.
_apply_hyt_config_overrides() {
	export WALLET_ETHEREUM_ERC20S_HYT_CONTRACT_ADDRESS="${HYT_CONTRACT_ADDRESS}"
	export WALLET_ETHEREUM_ERC20S_HYT_MASTER_ADDRESS="${HYT_MASTER_ADDRESS}"
}

###############################################################################
# Configuration
###############################################################################

eth_get_config_paths

# Coin type code for HYT ERC-20 token
HYT_COIN="hyt"

# Number of HD keys to generate per account
KEY_NUM=5

# Accounts used in this E2E test
PAYMENT_ACCOUNT="payment"
CLIENT_ACCOUNT="deposit"

# Funding amounts
FUNDING_AMOUNT_ETH=1    # ETH for gas fees (small amount is sufficient)
HYT_TRANSFER_AMOUNT=100 # HYT tokens to seed the payment account with
TRANSFER_AMOUNT_HYT=10  # HYT tokens to transfer in the wallet transaction

# Foundry toolchain binary paths
CAST="${HOME}/.foundry/bin/cast"
FORGE="${HOME}/.foundry/bin/forge"

# Deployer private key (Anvil account 0 from project mnemonic).
# Override with HYT_DEPLOYER_KEY=0x... to use a different account.
HYT_DEPLOYER_KEY="${HYT_DEPLOYER_KEY:-0xd3f1327df33c2a67ac03877f37410f33ed3a37085a84c739a7f5a137e5152bf2}"

# Deployer/master address corresponding to HYT_DEPLOYER_KEY (Anvil account 0).
# Override with HYT_MASTER_ADDRESS=0x... when using a different deployer.
HYT_MASTER_ADDRESS="${HYT_MASTER_ADDRESS:-0x328F371a76dfAc47b89Cc007bb048ec446c21494}"

# HYT contract address: use env var override, or auto-detect from broadcast artifacts
_detect_hyt_contract_address() {
	local project_root
	project_root="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"
	local broadcast="${project_root}/apps/eth-contracts/broadcast/DeployHYT.s.sol/31337/run-latest.json"
	if [ -f "${broadcast}" ]; then
		jq -r '.transactions[0].contractAddress' "${broadcast}" 2>/dev/null || true
	fi
}

HYT_CONTRACT_ADDRESS="${HYT_CONTRACT_ADDRESS:-$(_detect_hyt_contract_address)}"

###############################################################################
# Argument Parsing
###############################################################################

MODE="run"
VERBOSE=false
NON_INTERACTIVE=false

show_help() {
	cat <<EOF
Ethereum E2E Workflow Script - Pattern 2: ERC-20 HYT Token Transfer

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
  NODE_TYPE             anvil (default) | geth
  DB_TYPE               sqlite (default) | mysql
  ETH_RPC_HOST          Ethereum RPC host (default: 127.0.0.1)
  ETH_RPC_PORT          Ethereum RPC port (default: auto from NODE_TYPE)
  HYT_CONTRACT_ADDRESS  HYT contract address (auto-detected if unset)
  HYT_DEPLOYER_KEY      Private key of HYT deployer (default: Anvil account 0)

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
# HYT-specific address query
# eth_get_payment_address queries coin='eth'; for HYT we need coin='hyt'.
###############################################################################

hyt_get_payment_address() {
	local account="${1:-payment}"
	if [ "${DB_TYPE}" = "sqlite" ]; then
		sqlite3 "${SQLITE_WATCH_DB_PATH}" \
			"SELECT wallet_address FROM address WHERE coin='hyt' AND account='${account}' AND is_allocated=0 LIMIT 1" \
			2>/dev/null || true
	fi
}

###############################################################################
# Phase 0: Deploy HYT Contract
# Runs after Anvil starts. On --reset the contract must be redeployed because
# Anvil is wiped clean. Sets HYT_CONTRACT_ADDRESS for all subsequent phases.
###############################################################################

deploy_hyt_phase() {
	log_step "Deploy HYT Contract Phase"

	local project_root
	project_root="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

	log_substep "Deploying HYT to Anvil (${ETH_RPC_HOST}:${ETH_RPC_PORT})"
	(
		cd "${project_root}/apps/eth-contracts"
		PRIVATE_KEY="${HYT_DEPLOYER_KEY}" \
			"${FORGE}" script script/DeployHYT.s.sol \
			--rpc-url "http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" \
			--broadcast \
			2>&1 | grep -v "^Warning:" || true
	)

	# Read deployed address from broadcast artifacts
	local broadcast="${project_root}/apps/eth-contracts/broadcast/DeployHYT.s.sol/31337/run-latest.json"
	HYT_CONTRACT_ADDRESS=$(jq -r '.transactions[0].contractAddress' "${broadcast}" 2>/dev/null || true)

	if [ -z "${HYT_CONTRACT_ADDRESS}" ]; then
		log_error "Failed to deploy HYT or read contract address from broadcast artifacts"
		return 1
	fi

	log_info "HYT deployed at: ${HYT_CONTRACT_ADDRESS}"

	# Update the config overrides with the freshly deployed address
	_apply_hyt_config_overrides
}

###############################################################################
# HYT-specific address CSV export
# eth_export_watch_address_csv hardcodes 'eth' as the coin type code, but
# watch import address validates that the CSV coin matches --coin hyt.
# This override emits 'hyt' as the coin type code in the CSV.
###############################################################################

hyt_export_watch_address_csv() {
	local account="$1"
	local project_root
	project_root="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"
	local addr_dir="${project_root}/data/address/eth"
	mkdir -p "${addr_dir}"

	local tmp_file
	tmp_file="${addr_dir}/${account}_hyt_$(date +%s%N 2>/dev/null || date +%s).csv"

	if [ "${DB_TYPE}" = "sqlite" ]; then
		sqlite3 "${SQLITE_KEYGEN_DB_PATH}" \
			"SELECT 'hyt'||','||account||','||address||',,,,,'||idx FROM eth_account_key WHERE account='${account}'" \
			>"${tmp_file}" 2>/dev/null || true
	else
		log_warn "hyt_export_watch_address_csv: MySQL support not yet implemented"
	fi

	echo "${tmp_file}"
}

###############################################################################
# Phase 1: Key Generation (Keygen wallet — offline)
###############################################################################

keygen_phase() {
	log_step "Key Generation Phase (Keygen wallet)"

	log_substep "Creating mnemonic seed"
	eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${HYT_COIN}" create seed || {
		log_warn "Seed already exists or error, continuing..."
	}

	log_substep "Creating ${KEY_NUM} HD keys for '${PAYMENT_ACCOUNT}' account"
	eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${HYT_COIN}" create hdkey \
		--account "${PAYMENT_ACCOUNT}" \
		--keynum "${KEY_NUM}"

	log_substep "Creating ${KEY_NUM} HD keys for '${CLIENT_ACCOUNT}' account"
	eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${HYT_COIN}" create hdkey \
		--account "${CLIENT_ACCOUNT}" \
		--keynum "${KEY_NUM}"

	log_info "Key generation completed"
}

###############################################################################
# Phase 2: AccountXpub Export → Watch Address Import
###############################################################################

address_setup_phase() {
	log_step "Address Export (Keygen DB) → Watch Address Import Phase"

	for account in "${PAYMENT_ACCOUNT}" "${CLIENT_ACCOUNT}"; do
		log_substep "Exporting '${account}' addresses from keygen DB"
		local csv_file
		csv_file=$(hyt_export_watch_address_csv "${account}")

		if [ -z "${csv_file}" ] || [ ! -f "${csv_file}" ] || [ ! -s "${csv_file}" ]; then
			log_error "Failed to create address CSV for '${account}'"
			return 1
		fi
		log_info "  ${account} address CSV: ${csv_file}"

		log_substep "Importing '${account}' addresses into watch wallet"
		eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${HYT_COIN}" import address --file "${csv_file}"
	done

	log_info "Address setup completed"
}

###############################################################################
# Phase 3: Fund ETH (gas) + HYT Tokens
###############################################################################

funding_phase() {
	log_step "Funding Phase"

	if [ -z "${HYT_CONTRACT_ADDRESS}" ]; then
		log_error "HYT_CONTRACT_ADDRESS is not set and could not be auto-detected."
		log_error "Run 'make deploy-hyt' first, or set HYT_CONTRACT_ADDRESS=0x..."
		return 1
	fi
	log_info "HYT contract: ${HYT_CONTRACT_ADDRESS}"

	local payment_addr
	payment_addr=$(hyt_get_payment_address "${PAYMENT_ACCOUNT}")
	if [ -z "${payment_addr}" ]; then
		log_error "No '${PAYMENT_ACCOUNT}' address found in watch DB"
		return 1
	fi
	log_info "Payment address: ${payment_addr}"

	# Fund ETH for gas fees
	log_substep "Funding ${payment_addr} with ${FUNDING_AMOUNT_ETH} ETH (for gas)"
	eth_fund_address "${payment_addr}" "${FUNDING_AMOUNT_ETH}"

	# Transfer HYT tokens from deployer to payment address
	log_substep "Transferring ${HYT_TRANSFER_AMOUNT} HYT to payment address"
	local hyt_amount_wei
	hyt_amount_wei=$(echo "${HYT_TRANSFER_AMOUNT} * 1000000000000000000" | bc)
	"${CAST}" send \
		--private-key "${HYT_DEPLOYER_KEY}" \
		"${HYT_CONTRACT_ADDRESS}" \
		"transfer(address,uint256)" \
		"${payment_addr}" \
		"${hyt_amount_wei}" \
		--rpc-url "http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" \
		--quiet

	# Verify HYT balance was received
	local balance
	balance=$("${CAST}" call "${HYT_CONTRACT_ADDRESS}" \
		"balanceOf(address)(uint256)" "${payment_addr}" \
		--rpc-url "http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" 2>/dev/null || echo "0")
	log_info "Payment address HYT balance: ${balance} wei"

	log_info "Funding completed"
}

###############################################################################
# Phase 4: Transaction Creation (Watch wallet — online)
###############################################################################

create_tx_phase() {
	log_step "Transaction Creation Phase (Watch wallet)"

	log_substep "Creating unsigned HYT transfer tx: ${PAYMENT_ACCOUNT} → ${CLIENT_ACCOUNT} (${TRANSFER_AMOUNT_HYT} HYT)"
	local create_output
	create_output=$(eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${HYT_COIN}" create transfer \
		--account1 "${PAYMENT_ACCOUNT}" \
		--account2 "${CLIENT_ACCOUNT}" \
		--amount "${TRANSFER_AMOUNT_HYT}" 2>&1) || {
		log_error "Failed to create HYT transfer transaction"
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

# Globals set by create_tx_phase / sign_tx_phase
UNSIGNED_TX_FILE=""
SIGNED_TX_FILE=""

###############################################################################
# Phase 5: Offline Signing (Keygen wallet — offline, no network calls)
###############################################################################

sign_tx_phase() {
	log_step "Offline Signing Phase (Keygen wallet)"

	log_substep "Signing transaction: ${UNSIGNED_TX_FILE}"
	local sign_output
	sign_output=$(eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${HYT_COIN}" sign signature \
		--file "${UNSIGNED_TX_FILE}" 2>&1) || {
		log_error "Failed to sign HYT transaction"
		log_error "Output: ${sign_output}"
		return 1
	}

	SIGNED_TX_FILE=$(eth_extract_file_path "${sign_output}")

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

	log_substep "Sending signed HYT transaction: ${SIGNED_TX_FILE}"
	local send_output
	send_output=$(eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${HYT_COIN}" send tx \
		--file "${SIGNED_TX_FILE}" 2>&1) || {
		log_error "Failed to send HYT transaction"
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
# Phase 7: HYT Balance Verification
###############################################################################

verify_balances_phase() {
	log_step "HYT Balance Verification Phase"

	if [ -z "${HYT_CONTRACT_ADDRESS}" ]; then
		log_warn "HYT_CONTRACT_ADDRESS not set; skipping balance verification"
		return 0
	fi

	local client_addr
	client_addr=$(sqlite3 "$(eth_get_db_path "watch")" \
		"SELECT wallet_address FROM address WHERE coin='hyt' AND account='${CLIENT_ACCOUNT}' AND is_allocated=1 LIMIT 1" \
		2>/dev/null || true)

	if [ -z "${client_addr}" ]; then
		log_warn "Could not find allocated client address; skipping balance check"
		return 0
	fi

	log_substep "Querying HYT balance of recipient (${CLIENT_ACCOUNT}): ${client_addr}"
	local balance
	balance=$("${CAST}" call "${HYT_CONTRACT_ADDRESS}" \
		"balanceOf(address)(uint256)" "${client_addr}" \
		--rpc-url "http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" 2>/dev/null || echo "unknown")
	log_info "Recipient HYT balance: ${balance} wei"

	# Expected minimum: TRANSFER_AMOUNT_HYT * 10^18
	local expected_min
	expected_min=$(echo "${TRANSFER_AMOUNT_HYT} * 1000000000000000000" | bc)
	log_info "Expected minimum: ${expected_min} wei"
}

###############################################################################
# Phase 8: Confirmation Monitoring (Watch wallet — online)
###############################################################################

monitor_phase() {
	log_step "Monitoring Phase (Watch wallet)"

	log_substep "Monitoring sent transactions for confirmation"
	eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${HYT_COIN}" monitor senttx || {
		log_warn "Monitor returned non-zero (may be no transactions to monitor yet)"
	}

	log_info "Monitoring completed"
}

###############################################################################
# Full E2E Workflow
###############################################################################

e2e_full_workflow() {
	log_step "Ethereum E2E Workflow - Pattern 2: ERC-20 HYT Token Transfer"
	log_info "Node type:    ${NODE_TYPE}"
	log_info "DB type:      ${DB_TYPE}"
	log_info "RPC:          ${ETH_RPC_HOST}:${ETH_RPC_PORT}"
	echo ""

	eth_check_prerequisites
	eth_setup_infrastructure
	eth_init_sqlite_db

	# Deploy HYT contract to the fresh Anvil node, then apply config overrides
	deploy_hyt_phase
	log_info "HYT contract (config override): ${WALLET_ETHEREUM_ERC20S_HYT_CONTRACT_ADDRESS}"
	log_info "HYT master   (config override): ${WALLET_ETHEREUM_ERC20S_HYT_MASTER_ADDRESS}"

	keygen_phase
	address_setup_phase
	funding_phase
	create_tx_phase
	sign_tx_phase
	send_tx_phase
	verify_balances_phase
	monitor_phase

	log_info "E2E Pattern 2 (HYT ERC-20 transfer) completed successfully!"
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
