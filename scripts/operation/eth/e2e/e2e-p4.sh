#!/usr/bin/env bash

# Ethereum E2E Workflow Script - Pattern 4: MPC-TSS 2-of-3 Threshold Signing
# This script automates the complete MPC-TSS payment flow:
#   keygen pre-params (x3) → keygen dkg (x3 concurrently) →
#   fund joint address → watch create mpc →
#   keygen serve mpc (2 nodes background) → watch send mpc →
#   verify ETH balance and signed file
#
# Usage: ./scripts/operation/eth/e2e/e2e-p4.sh [OPTIONS]
#
# Options:
#   --cleanup           Stop containers and cleanup state
#   --reset             Full reset and run from scratch
#   --verbose           Enable verbose output (set -x)
#   --non-interactive   Run without prompts (for CI/CD)
#   -h, --help          Display this help message
#
# Reference Documentation:
#   docs/chains/eth/transaction-patterns.md
#
# Transaction Pattern:
#   Pattern 4: ETH MPC-TSS 2-of-3 threshold signing
#   - Joint address: derived from DKG ceremony across 3 keygen nodes
#   - Parties: node1, node2, node3 (threshold: 2-of-3)
#   - Signing: Distributed threshold signatures via tss-lib/v2
#   - Transport: gRPC relay coordinated by Watch wallet
#
# Environment Variables:
#   NODE_TYPE        Node type: anvil (default) or geth
#   DB_TYPE          Database type: sqlite (default) or mysql
#   ETH_RPC_HOST     Ethereum RPC host (default: 127.0.0.1)
#   ETH_RPC_PORT     Ethereum RPC port (default: 8546 for anvil, 8545 for geth)
#   MPC_PASSPHRASE   Passphrase for shard encryption (default: e2e-mpc-test-passphrase)
#   DEPLOYER_ADDRESS Recipient address for the MPC transfer (default: Anvil account 3)

set -euo pipefail

###############################################################################
# Script Configuration
###############################################################################

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Source the ETH common utilities (which also sources common.sh)
# shellcheck source=../eth_common.sh
source "${SCRIPT_DIR}/../eth_common.sh"

# Pattern identifier for DB isolation (also used by eth_get_db_path)
export E2E_PATTERN="p4"

# Re-initialize database with correct pattern now that E2E_PATTERN is set
eth_init_database

###############################################################################
# Environment Variable Overrides for Configuration
###############################################################################

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

# Foundry toolchain binary
CAST="${HOME}/.foundry/bin/cast"

# MPC ceremony parameters
MPC_THRESHOLD=2
MPC_ALL_PARTY_IDS="node1,node2,node3"
MPC_NODE1_ADDR="localhost:9001"
MPC_NODE2_ADDR="localhost:9002"
MPC_NODE3_ADDR="localhost:9003"
MPC_PASSPHRASE="${MPC_PASSPHRASE:-e2e-mpc-test-passphrase}"

# ETH transfer amounts
FUNDING_AMOUNT_ETH=10
TRANSFER_AMOUNT_ETH=1

# Recipient for the MPC payment (Anvil account 3 — isolated for P4 in parallel runner)
RECIPIENT_ADDRESS="${DEPLOYER_ADDRESS:-0x0505bCf2c03Af6F5D36b1591BC9175295986fD2B}"

# Data directory for MPC artefacts (pre-params, shards, logs)
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"
MPC_DATA_DIR="${PROJECT_ROOT}/data/mpc/e2e-p4"

# Pre-params paths (one per node)
PREPARAMS_NODE1="${MPC_DATA_DIR}/pre_params_node1.json"
PREPARAMS_NODE2="${MPC_DATA_DIR}/pre_params_node2.json"
PREPARAMS_NODE3="${MPC_DATA_DIR}/pre_params_node3.json"

# Encrypted shard paths (one per node)
SHARD_NODE1="${MPC_DATA_DIR}/shard_node1.json"
SHARD_NODE2="${MPC_DATA_DIR}/shard_node2.json"
SHARD_NODE3="${MPC_DATA_DIR}/shard_node3.json"

# DKG stdout/stderr capture files
DKG_LOG_NODE1="${MPC_DATA_DIR}/dkg_node1.log"
DKG_LOG_NODE2="${MPC_DATA_DIR}/dkg_node2.log"
DKG_LOG_NODE3="${MPC_DATA_DIR}/dkg_node3.log"

# DKG wall-clock timeout in seconds (CPU-intensive, allow generous budget)
DKG_TIMEOUT=300

# How many seconds to wait for serve-mpc daemons to open their gRPC ports
SERVE_READY_TIMEOUT=15

###############################################################################
# Argument Parsing
###############################################################################

MODE="run"
VERBOSE=false
NON_INTERACTIVE=false

show_help() {
	cat <<EOF
Ethereum E2E Workflow Script - Pattern 4: MPC-TSS 2-of-3 Threshold Signing

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
  MPC_PASSPHRASE   Shard encryption passphrase
  DEPLOYER_ADDRESS Recipient address (default: Anvil account 3)

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
# MPC Process Management
###############################################################################

# PIDs for background serve-mpc daemons
NODE1_PID=0
NODE2_PID=0

# Stop all background MPC node daemon processes
p4_stop_node_daemons() {
	if [ "${NODE1_PID}" -ne 0 ]; then
		log_substep "Stopping MPC node1 daemon (PID: ${NODE1_PID})"
		kill "${NODE1_PID}" 2>/dev/null || true
		wait "${NODE1_PID}" 2>/dev/null || true
		NODE1_PID=0
	fi
	if [ "${NODE2_PID}" -ne 0 ]; then
		log_substep "Stopping MPC node2 daemon (PID: ${NODE2_PID})"
		kill "${NODE2_PID}" 2>/dev/null || true
		wait "${NODE2_PID}" 2>/dev/null || true
		NODE2_PID=0
	fi
}

# Wait until a TCP port accepts connections (daemon ready check)
p4_wait_for_port() {
	local host_port="$1"
	local max_wait="${2:-${SERVE_READY_TIMEOUT}}"
	local host="${host_port%%:*}"
	local port="${host_port##*:}"
	local count=0

	while [ "${count}" -lt "${max_wait}" ]; do
		if nc -z "${host}" "${port}" 2>/dev/null; then
			return 0
		fi
		sleep 1
		count=$((count + 1))
	done

	log_error "Port ${host_port} not ready after ${max_wait} seconds"
	return 1
}

###############################################################################
# Phase 0: Pre-Params Generation
###############################################################################

preparams_phase() {
	log_step "Pre-Params Phase — Generate Paillier parameters for 3 nodes"

	mkdir -p "${MPC_DATA_DIR}"

	for node_num in 1 2 3; do
		local preparams_file="${MPC_DATA_DIR}/pre_params_node${node_num}.json"
		if [ -f "${preparams_file}" ]; then
			log_info "  node${node_num}: pre-params already exist, skipping"
			continue
		fi
		log_substep "Generating pre-params for node${node_num} (CPU-intensive, may take minutes)..."
		eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" pre-params \
			--output "${preparams_file}"
		log_info "  node${node_num}: pre-params written to ${preparams_file}"
	done

	log_info "Pre-params generation completed"
}

###############################################################################
# Phase 1: DKG Ceremony — All 3 Nodes Concurrently
###############################################################################

# Global: set by dkg_phase, read by all subsequent phases
JOINT_ADDRESS=""

dkg_phase() {
	log_step "DKG Phase — Distributed Key Generation (3 nodes concurrently)"

	mkdir -p "${MPC_DATA_DIR}"

	# Remove stale shard and log files from a previous run
	rm -f "${SHARD_NODE1}" "${SHARD_NODE2}" "${SHARD_NODE3}"
	rm -f "${DKG_LOG_NODE1}" "${DKG_LOG_NODE2}" "${DKG_LOG_NODE3}"

	log_substep "Starting DKG ceremony (all 3 nodes, timeout: ${DKG_TIMEOUT}s)"

	# Run all three DKG nodes concurrently in the background.
	# Each node:
	#  - Identifies itself with --party-id
	#  - Knows all participants via --all-party-ids
	#  - Lists the other two nodes' gRPC addresses via --peers
	#  - Uses locally generated pre-params
	#  - Writes its encrypted key shard to --shard-output
	(eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" dkg \
		--party-id "node1" \
		--all-party-ids "${MPC_ALL_PARTY_IDS}" \
		--threshold "${MPC_THRESHOLD}" \
		--peers "${MPC_NODE2_ADDR},${MPC_NODE3_ADDR}" \
		--pre-params-path "${PREPARAMS_NODE1}" \
		--shard-output "${SHARD_NODE1}" \
		--passphrase "${MPC_PASSPHRASE}" >"${DKG_LOG_NODE1}" 2>&1) &
	local dkg_pid1=$!

	(eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" dkg \
		--party-id "node2" \
		--all-party-ids "${MPC_ALL_PARTY_IDS}" \
		--threshold "${MPC_THRESHOLD}" \
		--peers "${MPC_NODE1_ADDR},${MPC_NODE3_ADDR}" \
		--pre-params-path "${PREPARAMS_NODE2}" \
		--shard-output "${SHARD_NODE2}" \
		--passphrase "${MPC_PASSPHRASE}" >"${DKG_LOG_NODE2}" 2>&1) &
	local dkg_pid2=$!

	(eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" dkg \
		--party-id "node3" \
		--all-party-ids "${MPC_ALL_PARTY_IDS}" \
		--threshold "${MPC_THRESHOLD}" \
		--peers "${MPC_NODE1_ADDR},${MPC_NODE2_ADDR}" \
		--pre-params-path "${PREPARAMS_NODE3}" \
		--shard-output "${SHARD_NODE3}" \
		--passphrase "${MPC_PASSPHRASE}" >"${DKG_LOG_NODE3}" 2>&1) &
	local dkg_pid3=$!

	# Enforce a wall-clock timeout: kill all DKG nodes if they exceed it
	(sleep "${DKG_TIMEOUT}" &&
		kill "${dkg_pid1}" "${dkg_pid2}" "${dkg_pid3}" 2>/dev/null || true) &
	local timeout_pid=$!

	# Wait for all three DKG processes
	local dkg_exit=0
	wait "${dkg_pid1}" || {
		log_error "DKG node1 failed (exit $?)"
		dkg_exit=1
	}
	wait "${dkg_pid2}" || {
		log_error "DKG node2 failed (exit $?)"
		dkg_exit=1
	}
	wait "${dkg_pid3}" || {
		log_error "DKG node3 failed (exit $?)"
		dkg_exit=1
	}

	# Cancel the timeout watchdog
	kill "${timeout_pid}" 2>/dev/null || true
	wait "${timeout_pid}" 2>/dev/null || true

	if [ "${dkg_exit}" -ne 0 ]; then
		log_error "DKG ceremony failed. Printing last lines from node logs:"
		for log_file in "${DKG_LOG_NODE1}" "${DKG_LOG_NODE2}" "${DKG_LOG_NODE3}"; do
			if [ -f "${log_file}" ]; then
				log_error "--- $(basename "${log_file}") ---"
				tail -5 "${log_file}" | sed 's/^/  /' >&2
			fi
		done
		return 1
	fi

	# Verify that all three shard files were created
	for shard_file in "${SHARD_NODE1}" "${SHARD_NODE2}" "${SHARD_NODE3}"; do
		if [ ! -f "${shard_file}" ]; then
			log_error "Expected shard file not found: ${shard_file}"
			return 1
		fi
	done

	# Extract the joint ETH address from node1 output.
	# Expected format: "  Joint ETH Address: 0x..."
	JOINT_ADDRESS=$(grep 'Joint ETH Address:' "${DKG_LOG_NODE1}" | awk '{print $NF}' | head -1 || true)
	if [ -z "${JOINT_ADDRESS}" ]; then
		log_error "Could not extract joint ETH address from DKG node1 output"
		log_error "--- dkg_node1.log ---"
		cat "${DKG_LOG_NODE1}" >&2
		return 1
	fi

	log_info "DKG ceremony completed. Joint ETH address: ${JOINT_ADDRESS}"
	log_info "Shard files written to ${MPC_DATA_DIR}/"
}

###############################################################################
# Phase 2: Fund Joint MPC Address
###############################################################################

fund_phase() {
	log_step "Fund Phase — Fund joint address ${JOINT_ADDRESS} with ${FUNDING_AMOUNT_ETH} ETH"

	eth_fund_address "${JOINT_ADDRESS}" "${FUNDING_AMOUNT_ETH}"

	log_info "Funding completed"
}

###############################################################################
# Phase 3: Create Unsigned MPC Transaction File
###############################################################################

# Global: set by create_mpc_phase, read by send_mpc_phase and verify_phase
MPC_TX_FILE=""

create_mpc_phase() {
	log_step "Create Phase — Create unsigned MPC transaction file (Watch wallet)"

	log_substep "Creating MPC tx: ${JOINT_ADDRESS} → ${RECIPIENT_ADDRESS} (${TRANSFER_AMOUNT_ETH} ETH)"
	local create_output
	create_output=$(eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${ETH_COIN}" create mpc \
		--from "${JOINT_ADDRESS}" \
		--to "${RECIPIENT_ADDRESS}" \
		--amount "${TRANSFER_AMOUNT_ETH}" \
		--threshold "${MPC_THRESHOLD}" \
		--party-ids "${MPC_ALL_PARTY_IDS}" \
		--action-type "payment" 2>&1) || {
		log_error "Failed to create MPC transaction"
		log_error "Output: ${create_output}"
		return 1
	}

	# Extract file path from "  File:      <path>" line
	MPC_TX_FILE=$(echo "${create_output}" | grep -E '^ *File:' | awk '{print $NF}' | head -1)
	if [ -z "${MPC_TX_FILE}" ]; then
		log_error "Failed to extract MPC transaction file path from output"
		log_error "Output: ${create_output}"
		return 1
	fi

	log_info "Unsigned MPC transaction file: ${MPC_TX_FILE}"

	# Verify file exists and contains tx_type: unsigned
	if [ ! -f "${MPC_TX_FILE}" ]; then
		log_error "MPC transaction file not found: ${MPC_TX_FILE}"
		return 1
	fi
	local tx_type
	tx_type=$(jq -r '.tx_type' "${MPC_TX_FILE}" 2>/dev/null || echo "unknown")
	if [ "${tx_type}" != "unsigned" ]; then
		log_error "Expected tx_type=unsigned, got: ${tx_type}"
		return 1
	fi
	log_info "Verified: tx_type=${tx_type}"
}

###############################################################################
# Phase 4: Start 2-of-3 MPC Signing Node Daemons
###############################################################################

start_nodes_phase() {
	log_step "Node Daemon Phase — Start 2 of 3 MPC signing nodes in background"

	# Start node1
	log_substep "Starting MPC node1 daemon on ${MPC_NODE1_ADDR}"
	eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" serve mpc \
		--listen-addr "${MPC_NODE1_ADDR}" \
		--shard-path "${SHARD_NODE1}" \
		--passphrase "${MPC_PASSPHRASE}" \
		--party-id "node1" \
		--all-party-ids "${MPC_ALL_PARTY_IDS}" >"${MPC_DATA_DIR}/serve_node1.log" 2>&1 &
	NODE1_PID=$!
	log_info "  node1 started (PID: ${NODE1_PID})"

	# Start node2
	log_substep "Starting MPC node2 daemon on ${MPC_NODE2_ADDR}"
	eth_keygen_cmd -c "${ETH_CONFIG_KEYGEN}" --coin "${ETH_COIN}" serve mpc \
		--listen-addr "${MPC_NODE2_ADDR}" \
		--shard-path "${SHARD_NODE2}" \
		--passphrase "${MPC_PASSPHRASE}" \
		--party-id "node2" \
		--all-party-ids "${MPC_ALL_PARTY_IDS}" >"${MPC_DATA_DIR}/serve_node2.log" 2>&1 &
	NODE2_PID=$!
	log_info "  node2 started (PID: ${NODE2_PID})"

	# Wait for both daemons to bind their gRPC ports before proceeding
	log_substep "Waiting for node daemons to be ready..."
	p4_wait_for_port "${MPC_NODE1_ADDR}" || {
		log_error "node1 daemon did not bind ${MPC_NODE1_ADDR} in time"
		if [ -f "${MPC_DATA_DIR}/serve_node1.log" ]; then
			cat "${MPC_DATA_DIR}/serve_node1.log" >&2
		fi
		return 1
	}
	p4_wait_for_port "${MPC_NODE2_ADDR}" || {
		log_error "node2 daemon did not bind ${MPC_NODE2_ADDR} in time"
		if [ -f "${MPC_DATA_DIR}/serve_node2.log" ]; then
			cat "${MPC_DATA_DIR}/serve_node2.log" >&2
		fi
		return 1
	}
	log_info "Both MPC node daemons are ready"
}

###############################################################################
# Phase 5: Send — Initiate TSS Signing Session and Broadcast
###############################################################################

TX_HASH=""

send_mpc_phase() {
	log_step "Send Phase — Initiate MPC signing session and broadcast (Watch wallet)"

	log_substep "Signing and broadcasting: ${MPC_TX_FILE}"
	log_info "  Signing nodes: ${MPC_NODE1_ADDR}, ${MPC_NODE2_ADDR}"
	local send_output
	send_output=$(eth_watch_cmd -c "${ETH_CONFIG_WATCH}" --coin "${ETH_COIN}" send mpc \
		--file "${MPC_TX_FILE}" \
		--peer-addrs "${MPC_NODE1_ADDR},${MPC_NODE2_ADDR}" 2>&1) || {
		log_error "Failed to send MPC transaction"
		log_error "Output: ${send_output}"
		return 1
	}

	# Extract TX Hash from "  TX Hash: 0x..." line
	TX_HASH=$(echo "${send_output}" | grep -E '^ *TX Hash:' | awk '{print $NF}' | head -1 || true)
	if [ -z "${TX_HASH}" ]; then
		log_warn "Could not extract TX Hash from output; tx may still have been submitted"
		log_info "Raw output: ${send_output}"
	else
		log_info "MPC transaction broadcast! TX Hash: ${TX_HASH}"
	fi
}

###############################################################################
# Phase 6: Verify — ETH Balance and Signed File
###############################################################################

verify_phase() {
	log_step "Verify Phase — Assert ETH balance change and signed file"

	# Check recipient balance
	log_substep "Querying recipient balance: ${RECIPIENT_ADDRESS}"
	local recipient_balance
	recipient_balance=$("${CAST}" balance "${RECIPIENT_ADDRESS}" \
		--rpc-url "http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" 2>/dev/null || echo "0")
	log_info "Recipient balance: ${recipient_balance} Wei"

	local expected_min_wei
	expected_min_wei=$(echo "${TRANSFER_AMOUNT_ETH} * 1000000000000000000" | bc)
	log_info "Expected minimum recipient balance: ${expected_min_wei} Wei"

	if [ "${recipient_balance}" -ge "${expected_min_wei}" ] 2>/dev/null; then
		log_info "✓ Recipient balance check passed (${recipient_balance} >= ${expected_min_wei})"
	else
		log_warn "Balance check inconclusive; recipient_balance=${recipient_balance}"
	fi

	# Check joint address balance (should have decreased)
	log_substep "Querying joint address balance: ${JOINT_ADDRESS}"
	local joint_balance
	joint_balance=$("${CAST}" balance "${JOINT_ADDRESS}" \
		--rpc-url "http://${ETH_RPC_HOST}:${ETH_RPC_PORT}" 2>/dev/null || echo "unknown")
	log_info "Joint address balance after transfer: ${joint_balance} Wei"

	# Verify the signed MPC file exists with tx_type=signed and non-empty signed_tx_hex.
	# Naming: {base}_signed.json (e.g., payment_mpc_<uuid>_signed.json)
	local signed_file="${MPC_TX_FILE%.json}_signed.json"
	log_substep "Verifying signed MPC transaction file: ${signed_file}"

	if [ ! -f "${signed_file}" ]; then
		log_error "Signed MPC transaction file not found: ${signed_file}"
		return 1
	fi

	local signed_tx_type
	signed_tx_type=$(jq -r '.tx_type' "${signed_file}" 2>/dev/null || echo "unknown")
	if [ "${signed_tx_type}" = "signed" ]; then
		log_info "✓ tx_type=${signed_tx_type}"
	else
		log_error "Expected tx_type=signed, got: ${signed_tx_type}"
		return 1
	fi

	local signed_tx_hex
	signed_tx_hex=$(jq -r '.signed_tx_hex // ""' "${signed_file}" 2>/dev/null || echo "")
	if [ -n "${signed_tx_hex}" ]; then
		log_info "✓ signed_tx_hex is non-empty"
	else
		log_error "signed_tx_hex is empty in signed file"
		return 1
	fi

	log_info "Signed file verified: ${signed_file}"
	log_info "Verification phase completed"
}

###############################################################################
# Full E2E Workflow
###############################################################################

e2e_full_workflow() {
	log_step "Ethereum E2E Workflow - Pattern 4: MPC-TSS 2-of-3 Threshold Signing"
	log_info "Node type: ${NODE_TYPE}"
	log_info "DB type:   ${DB_TYPE}"
	log_info "RPC:       ${ETH_RPC_HOST}:${ETH_RPC_PORT}"
	log_info "MPC nodes: ${MPC_NODE1_ADDR}, ${MPC_NODE2_ADDR}, ${MPC_NODE3_ADDR}"
	log_info "Threshold: ${MPC_THRESHOLD}-of-3"
	log_info "Recipient: ${RECIPIENT_ADDRESS}"
	echo ""

	eth_check_prerequisites
	eth_setup_infrastructure
	eth_init_sqlite_db

	preparams_phase
	dkg_phase
	fund_phase
	create_mpc_phase
	start_nodes_phase
	send_mpc_phase
	verify_phase

	log_info "E2E Pattern 4 (MPC-TSS 2-of-3 threshold signing) completed successfully!"
}

###############################################################################
# Cleanup
###############################################################################

p4_cleanup_data() {
	log_substep "Removing MPC data directory: ${MPC_DATA_DIR}"
	rm -rf "${MPC_DATA_DIR}"
}

p4_cleanup() {
	p4_stop_node_daemons
	eth_cleanup
}

p4_full_reset() {
	p4_stop_node_daemons
	p4_cleanup_data
	eth_cleanup
	log_info "Full reset completed"
}

###############################################################################
# Main
###############################################################################

# Always stop background daemons and report the failed line on any error
trap 'p4_stop_node_daemons; log_error "Script failed at line $LINENO"' ERR
trap 'p4_stop_node_daemons' EXIT

main() {
	case "${MODE}" in
	cleanup)
		p4_cleanup
		;;
	reset)
		p4_full_reset
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
