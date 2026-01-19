#!/usr/bin/env bash

# Bitcoin E2E Parallel Test Runner
# This script runs all BTC E2E test patterns (P1-P11) in parallel using SQLite backend
# Usage: ./scripts/operation/btc/e2e/e2e-parallel-runner.sh [OPTIONS]
#
# Options:
#   --patterns <list>      Run specific patterns (e.g., "1,2,3" or "1-11")
#   --max-parallel <N>     Limit concurrent processes (default: 11 for all patterns)
#   --verbose              Show real-time output from all processes
#   --ci                   Non-interactive mode with structured output for CI
#   -h, --help             Display help message
#
# Description:
#   This script is designed for CI/CD environments to run BTC E2E tests in parallel.
#   Each pattern uses an isolated SQLite database to prevent conflicts.
#   The script collects exit codes from all parallel processes and reports a summary.
#
# Examples:
#   ./scripts/operation/btc/e2e/e2e-parallel-runner.sh --ci
#   ./scripts/operation/btc/e2e/e2e-parallel-runner.sh --patterns 1,2,3 --verbose
#   ./scripts/operation/btc/e2e/e2e-parallel-runner.sh --patterns 1-5 --max-parallel 3

set -euo pipefail

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# Default configuration
PATTERNS="1-11"
MAX_PARALLEL=11
VERBOSE=false
CI_MODE=false
LOG_DIR="${PROJECT_ROOT}/data/logs/e2e-parallel"

# Pattern script mapping
declare -A PATTERN_SCRIPTS=(
	[1]="e2e-p1-p2pkh-singlesig.sh"
	[2]="e2e-p2-p2pkh-2of3.sh"
	[3]="e2e-p3-p2sh-p2wpkh-singlesig.sh"
	[4]="e2e-p4-p2sh-p2wsh-2of3.sh"
	[5]="e2e-p5-p2wpkh-singlesig.sh"
	[6]="e2e-p6-p2wsh-2of3.sh"
	[7]="e2e-p7-p2wsh-3of3.sh"
	[8]="e2e-p8-p2sh-p2wsh-3of3.sh"
	[9]="e2e-p9-p2tr-singlesig.sh"
	[10]="e2e-p10-p2tr-musig2.sh"
	[11]="e2e-p11-p2tr-tapscript.sh"
)

# Pattern descriptions
declare -A PATTERN_DESCRIPTIONS=(
	[1]="P2PKH Single-sig"
	[2]="P2PKH 2-of-3 Multisig"
	[3]="P2SH-P2WPKH Single-sig"
	[4]="P2SH-P2WSH 2-of-3 Multisig"
	[5]="P2WPKH Native SegWit Single-sig"
	[6]="P2WSH Native SegWit 2-of-3 Multisig"
	[7]="P2WSH Native SegWit 3-of-3 Multisig"
	[8]="P2SH-P2WSH 3-of-3 Multisig"
	[9]="P2TR Taproot Single-sig"
	[10]="P2TR MuSig2 N-of-N"
	[11]="P2TR Tapscript M-of-N"
)

###############################################################################
# Utility Functions
###############################################################################

log_info() {
	if [[ "$CI_MODE" == "true" ]]; then
		echo "[INFO] $*"
	else
		echo -e "\033[0;36m[INFO]\033[0m $*"
	fi
}

log_success() {
	if [[ "$CI_MODE" == "true" ]]; then
		echo "[SUCCESS] $*"
	else
		echo -e "\033[0;32m[SUCCESS]\033[0m $*"
	fi
}

log_error() {
	if [[ "$CI_MODE" == "true" ]]; then
		echo "[ERROR] $*" >&2
	else
		echo -e "\033[0;31m[ERROR]\033[0m $*" >&2
	fi
}

log_warn() {
	if [[ "$CI_MODE" == "true" ]]; then
		echo "[WARN] $*"
	else
		echo -e "\033[0;33m[WARN]\033[0m $*"
	fi
}

show_help() {
	cat <<EOF
Bitcoin E2E Parallel Test Runner

Usage: $0 [OPTIONS]

Options:
  --patterns <list>      Run specific patterns (e.g., "1,2,3" or "1-11")
                         Default: "1-11" (all patterns)
  --max-parallel <N>     Limit concurrent processes
                         Default: 11 (run all patterns in parallel)
  --verbose              Show real-time output from all processes
  --ci                   Non-interactive mode with structured output for CI
  -h, --help             Display this help message

Examples:
  # Run all patterns in parallel (default)
  $0 --ci

  # Run specific patterns
  $0 --patterns 1,2,3 --verbose

  # Run patterns with limited parallelism
  $0 --patterns 1-5 --max-parallel 3 --ci

  # Run single pattern for debugging
  $0 --patterns 1 --verbose

Available Patterns:
  1:  P2PKH Single-sig
  2:  P2PKH 2-of-3 Multisig
  3:  P2SH-P2WPKH Single-sig
  4:  P2SH-P2WSH 2-of-3 Multisig
  5:  P2WPKH Native SegWit Single-sig
  6:  P2WSH Native SegWit 2-of-3 Multisig
  7:  P2WSH Native SegWit 3-of-3 Multisig
  8:  P2SH-P2WSH 3-of-3 Multisig
  9:  P2TR Taproot Single-sig
  10: P2TR MuSig2 N-of-N
  11: P2TR Tapscript M-of-N
EOF
}

# Parse pattern range (e.g., "1-11" or "1,2,3")
parse_patterns() {
	local pattern_spec="$1"
	local -a patterns

	if [[ "$pattern_spec" =~ ^([0-9]+)-([0-9]+)$ ]]; then
		# Range format: 1-11
		local start="${BASH_REMATCH[1]}"
		local end="${BASH_REMATCH[2]}"
		for ((i = start; i <= end; i++)); do
			patterns+=("$i")
		done
	else
		# Comma-separated format: 1,2,3
		IFS=',' read -ra patterns <<<"$pattern_spec"
	fi

	echo "${patterns[@]}"
}

# Validate pattern number
validate_pattern() {
	local pattern="$1"
	if [[ ! "${PATTERN_SCRIPTS[$pattern]+_}" ]]; then
		log_error "Invalid pattern: P=$pattern. Valid patterns are 1-11."
		return 1
	fi
	return 0
}

###############################################################################
# Infrastructure Management
###############################################################################

setup_shared_infrastructure() {
	log_info "=========================================="
	log_info "Setting up shared infrastructure"
	log_info "=========================================="

	# Source BTC common utilities
	# shellcheck source=../btc_common.sh
	source "${SCRIPT_DIR}/../btc_common.sh"

	# Set SQLite mode
	export DB_TYPE="sqlite"

	# Full cleanup first (stop any existing containers)
	log_info "Cleaning up any existing infrastructure..."
	log_info "Stopping Bitcoin containers..."
	docker compose -f "${PROJECT_ROOT}/compose.btc.yaml" down -v 2>/dev/null || true
	docker compose -f "${PROJECT_ROOT}/compose.yaml" down -v 2>/dev/null || true

	# Clean old SQLite databases
	log_info "Cleaning old SQLite databases..."
	rm -rf "${PROJECT_ROOT}/data/sqlite/btc"/*-e2e-p*.db 2>/dev/null || true

	# Start Bitcoin node containers (shared across all patterns)
	log_info "Starting shared Bitcoin node containers..."
	log_substep "Starting Bitcoin node containers"
	docker compose -f "${PROJECT_ROOT}/compose.btc.yaml" up -d btc-watch btc-keygen btc-sign1 btc-sign2

	# Wait for containers to be healthy
	log_substep "Waiting for containers to be healthy"
	for container in btc-watch btc-keygen btc-sign1 btc-sign2; do
		wait_for_healthy "$container" 90
	done

	log_info "Shared infrastructure is ready"
	log_info "Note: Each pattern will initialize its own SQLite database"
	echo ""
}

cleanup_shared_infrastructure() {
	log_info ""
	log_info "=========================================="
	log_info "Cleaning up shared infrastructure"
	log_info "=========================================="

	# Stop Bitcoin containers
	log_info "Stopping Bitcoin containers..."
	docker compose -f "${PROJECT_ROOT}/compose.btc.yaml" down -v 2>/dev/null || true

	# Clean up pattern-specific SQLite databases
	log_info "Cleaning pattern-specific SQLite databases..."
	rm -rf "${PROJECT_ROOT}/data/sqlite/btc"/*-e2e-p*.db 2>/dev/null || true

	log_info "Cleanup complete"
}

###############################################################################
# Main Test Execution
###############################################################################

run_pattern() {
	local pattern="$1"
	local script="${PATTERN_SCRIPTS[$pattern]}"
	local script_path="${SCRIPT_DIR}/${script}"
	local log_file="${LOG_DIR}/p${pattern}.log"
	local exit_code_file="${LOG_DIR}/p${pattern}.exitcode"

	# Ensure log directory exists
	mkdir -p "${LOG_DIR}"

	# Set up environment for this pattern
	export DB_TYPE="sqlite"
	export E2E_PATTERN="p${pattern}"
	export E2E_SHARED_INFRASTRUCTURE="true" # Tell script to skip infrastructure setup

	log_info "Starting P${pattern}: ${PATTERN_DESCRIPTIONS[$pattern]}"

	if [[ "$VERBOSE" == "true" ]]; then
		# Show real-time output (skip --reset since infrastructure is shared)
		"${script_path}" --non-interactive 2>&1 | tee "${log_file}"
		echo "$?" >"${exit_code_file}"
	else
		# Redirect to log file (skip --reset since infrastructure is shared)
		"${script_path}" --non-interactive >"${log_file}" 2>&1
		echo "$?" >"${exit_code_file}"
	fi
}

run_patterns_parallel() {
	local -a patterns=("$@")
	local -a pids=()
	local running=0

	log_info "Running ${#patterns[@]} patterns with max parallelism: ${MAX_PARALLEL}"
	log_info "Log directory: ${LOG_DIR}"

	# Start patterns with parallelism control
	for pattern in "${patterns[@]}"; do
		# Wait if we've reached max parallel processes
		while [[ $running -ge $MAX_PARALLEL ]]; do
			# Check for completed processes
			for pid in "${pids[@]}"; do
				if ! kill -0 "$pid" 2>/dev/null; then
					wait "$pid" 2>/dev/null || true
					running=$((running - 1))
					# Remove completed PID
					pids=("${pids[@]/$pid/}")
				fi
			done
			sleep 0.5
		done

		# Start new pattern
		run_pattern "$pattern" &
		local pid=$!
		pids+=("$pid")
		running=$((running + 1))

		# Small delay to avoid overwhelming the system
		sleep 0.1
	done

	# Wait for all remaining processes
	log_info "Waiting for all patterns to complete..."
	for pid in "${pids[@]}"; do
		[[ -n "$pid" ]] && wait "$pid" 2>/dev/null || true
	done
}

generate_summary() {
	local -a patterns=("$@")
	local -a passed=()
	local -a failed=()
	local total=${#patterns[@]}

	echo ""
	log_info "=========================================="
	log_info "E2E Test Summary"
	log_info "=========================================="

	for pattern in "${patterns[@]}"; do
		local exit_code_file="${LOG_DIR}/p${pattern}.exitcode"
		local log_file="${LOG_DIR}/p${pattern}.log"
		local exit_code=1

		if [[ -f "$exit_code_file" ]]; then
			exit_code=$(cat "$exit_code_file")
		fi

		if [[ "$exit_code" -eq 0 ]]; then
			passed+=("$pattern")
			log_success "P${pattern}: ${PATTERN_DESCRIPTIONS[$pattern]} - PASSED"
		else
			failed+=("$pattern")
			log_error "P${pattern}: ${PATTERN_DESCRIPTIONS[$pattern]} - FAILED (exit code: ${exit_code})"
			if [[ "$CI_MODE" == "true" && -f "$log_file" ]]; then
				log_error "  Log: ${log_file}"
				log_error "  Last 20 lines:"
				tail -20 "$log_file" | sed 's/^/    /'
			fi
		fi
	done

	echo ""
	log_info "Total: ${total}, Passed: ${#passed[@]}, Failed: ${#failed[@]}"

	if [[ ${#failed[@]} -gt 0 ]]; then
		log_error "Failed patterns: ${failed[*]}"
		log_info "Check logs in: ${LOG_DIR}"
		return 1
	else
		log_success "All patterns passed!"
		return 0
	fi
}

###############################################################################
# Main Script
###############################################################################

main() {
	# Parse arguments
	while [[ $# -gt 0 ]]; do
		case "$1" in
		--patterns)
			PATTERNS="$2"
			shift 2
			;;
		--max-parallel)
			MAX_PARALLEL="$2"
			shift 2
			;;
		--verbose)
			VERBOSE=true
			shift
			;;
		--ci)
			CI_MODE=true
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

	# Parse and validate patterns
	local -a patterns
	IFS=' ' read -ra patterns <<<"$(parse_patterns "$PATTERNS")"

	if [[ ${#patterns[@]} -eq 0 ]]; then
		log_error "No valid patterns specified"
		exit 1
	fi

	for pattern in "${patterns[@]}"; do
		if ! validate_pattern "$pattern"; then
			exit 1
		fi
	done

	# Clean up old log directory
	rm -rf "${LOG_DIR}"
	mkdir -p "${LOG_DIR}"

	# Show configuration
	log_info "Configuration:"
	log_info "  Patterns: ${patterns[*]}"
	log_info "  Max Parallel: ${MAX_PARALLEL}"
	log_info "  Verbose: ${VERBOSE}"
	log_info "  CI Mode: ${CI_MODE}"
	log_info "  Database: SQLite (isolated per pattern)"
	log_info "  Infrastructure: Shared Bitcoin nodes"
	echo ""

	# Setup trap to ensure cleanup happens
	trap cleanup_shared_infrastructure EXIT

	# Setup shared infrastructure (Bitcoin containers)
	setup_shared_infrastructure

	# Run patterns in parallel
	run_patterns_parallel "${patterns[@]}"

	# Generate and display summary
	local summary_exit=0
	if ! generate_summary "${patterns[@]}"; then
		summary_exit=1
	fi

	# Cleanup will be called by trap
	exit $summary_exit
}

# Execute main function
main "$@"
