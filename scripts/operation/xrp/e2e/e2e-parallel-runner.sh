#!/usr/bin/env bash

# XRP E2E Parallel Test Runner
# This script runs all XRP E2E test patterns in parallel using a shared rippled node.
# Usage: ./scripts/operation/xrp/e2e/e2e-parallel-runner.sh [OPTIONS]
#
# Options:
#   --patterns <list>      Run specific patterns (e.g., "1" or "1")
#   --max-parallel <N>     Limit concurrent processes (default: 1 for all patterns)
#   --verbose              Show real-time output from all processes
#   --ci                   Non-interactive mode with structured output for CI
#   -h, --help             Display help message
#
# Description:
#   This script is designed for CI/CD environments to run XRP E2E tests in parallel.
#   A single rippled standalone node is shared across all patterns; each pattern uses
#   an isolated SQLite database to prevent conflicts. The script collects exit codes
#   from all parallel processes and reports a summary.
#
# Examples:
#   ./scripts/operation/xrp/e2e/e2e-parallel-runner.sh --ci
#   ./scripts/operation/xrp/e2e/e2e-parallel-runner.sh --patterns 1 --verbose

set -euo pipefail

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# Default configuration
PATTERNS="1"
MAX_PARALLEL=1
VERBOSE=false
CI_MODE=false
LOG_DIR="${PROJECT_ROOT}/data/logs/e2e-parallel-xrp"
XRP_WS_HOST="${XRP_WS_HOST:-127.0.0.1}"
XRP_WS_PORT="${XRP_WS_PORT:-6006}"

# Pattern script mapping
declare -A PATTERN_SCRIPTS=(
	[1]="e2e-p1.sh"
)

# Pattern descriptions
declare -A PATTERN_DESCRIPTIONS=(
	[1]="Single-sig Payment Transfer"
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

show_help() {
	cat <<EOF
XRP E2E Parallel Test Runner

Usage: $0 [OPTIONS]

Options:
  --patterns <list>      Run specific patterns (e.g., "1")
                         Default: "1" (all patterns)
  --max-parallel <N>     Limit concurrent processes
                         Default: 1 (run all patterns in parallel)
  --verbose              Show real-time output from all processes
  --ci                   Non-interactive mode with structured output for CI
  -h, --help             Display this help message

Examples:
  # Run all patterns in parallel (default)
  $0 --ci

  # Run specific pattern with verbose output
  $0 --patterns 1 --verbose

Available Patterns:
  1:  Single-sig Payment Transfer
EOF
}

# Parse pattern range (e.g., "1" or "1,2")
parse_patterns() {
	local pattern_spec="$1"
	local -a patterns

	if [[ "$pattern_spec" =~ ^([0-9]+)-([0-9]+)$ ]]; then
		# Range format: 1-2
		local start="${BASH_REMATCH[1]}"
		local end="${BASH_REMATCH[2]}"
		for ((i = start; i <= end; i++)); do
			patterns+=("$i")
		done
	else
		# Comma-separated format: 1,2
		IFS=',' read -ra patterns <<<"$pattern_spec"
	fi

	echo "${patterns[@]}"
}

# Validate pattern number
validate_pattern() {
	local pattern="$1"
	if [[ ! "${PATTERN_SCRIPTS[$pattern]+_}" ]]; then
		local valid_patterns
		valid_patterns=$(printf "%s\n" "${!PATTERN_SCRIPTS[@]}" | sort -n | tr '\n' ',' | sed 's/,$//')
		log_error "Invalid pattern: P=$pattern. Valid patterns are ${valid_patterns}."
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

	# Full cleanup first (stop any existing containers)
	log_info "Cleaning up any existing infrastructure..."
	docker compose -f "${PROJECT_ROOT}/compose.xrp.yaml" down rippled 2>/dev/null || true

	# Clean old SQLite databases
	log_info "Cleaning old SQLite databases..."
	rm -f "${PROJECT_ROOT}/data/sqlite/xrp/"*-e2e-*.db 2>/dev/null || true

	# Start shared rippled node
	log_info "Starting shared rippled standalone node..."
	docker compose -f "${PROJECT_ROOT}/compose.xrp.yaml" up -d rippled

	# Wait for rippled to be ready
	log_info "Waiting for rippled on ${XRP_WS_HOST}:${XRP_WS_PORT}..."
	local max_wait=60
	local count=0
	while [ "${count}" -lt "${max_wait}" ]; do
		if docker compose -f "${PROJECT_ROOT}/compose.xrp.yaml" \
			exec -T rippled rippled --silent server_info >/dev/null 2>&1; then
			log_info "rippled is ready"
			echo ""
			return 0
		fi
		sleep 2
		count=$((count + 2))
	done

	log_error "rippled failed to start within ${max_wait} seconds"
	return 1
}

# shellcheck disable=SC2329  # Used in trap
cleanup_shared_infrastructure() {
	log_info ""
	log_info "=========================================="
	log_info "Cleaning up shared infrastructure"
	log_info "=========================================="

	# Stop rippled container
	log_info "Stopping rippled node..."
	docker compose -f "${PROJECT_ROOT}/compose.xrp.yaml" down rippled 2>/dev/null || true

	# Clean pattern-specific SQLite databases
	log_info "Cleaning SQLite databases..."
	rm -f "${PROJECT_ROOT}/data/sqlite/xrp/"*-e2e-*.db 2>/dev/null || true

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
	export XRP_WS_HOST="${XRP_WS_HOST}"
	export XRP_WS_PORT="${XRP_WS_PORT}"

	log_info "Starting P${pattern}: ${PATTERN_DESCRIPTIONS[$pattern]}"

	if [[ "$VERBOSE" == "true" ]]; then
		# Show real-time output
		"${script_path}" --non-interactive 2>&1 | tee "${log_file}"
		echo "${PIPESTATUS[0]}" >"${exit_code_file}"
	else
		# Redirect to log file
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
			local new_pids=()
			local finished_one=false
			for pid in "${pids[@]}"; do
				if kill -0 "$pid" 2>/dev/null; then
					new_pids+=("$pid")
				else
					wait "$pid" 2>/dev/null || true
					finished_one=true
				fi
			done
			pids=("${new_pids[@]}")
			running=${#pids[@]}
			if ! $finished_one; then
				sleep 0.5
			fi
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
	log_info "  Patterns:     ${patterns[*]}"
	log_info "  Max Parallel: ${MAX_PARALLEL}"
	log_info "  Verbose:      ${VERBOSE}"
	log_info "  CI Mode:      ${CI_MODE}"
	log_info "  Database:     SQLite (isolated per pattern)"
	log_info "  rippled WS:   ws://${XRP_WS_HOST}:${XRP_WS_PORT}"
	echo ""

	# Setup trap to ensure cleanup happens
	trap cleanup_shared_infrastructure EXIT

	# Setup shared infrastructure (rippled node)
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
