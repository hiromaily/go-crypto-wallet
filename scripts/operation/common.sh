#!/usr/bin/env bash

# Common Utility Functions for Wallet Operations
# This file provides shared utility functions for all operation scripts.
#
# Usage:
#   source "$(dirname "${BASH_SOURCE[0]}")/../common.sh"
#   or
#   source "${SCRIPT_DIR}/../common.sh"

###############################################################################
# Color Definitions
###############################################################################

# Only set colors if not already defined and if terminal supports colors
if [ -z "${RED:-}" ]; then
	if [ -t 1 ] && command -v tput >/dev/null 2>&1; then
		RED='\033[0;31m'
		GREEN='\033[0;32m'
		YELLOW='\033[1;33m'
		BLUE='\033[0;34m'
		NC='\033[0m' # No Color
	else
		RED=''
		GREEN=''
		YELLOW=''
		BLUE=''
		NC=''
	fi
fi

###############################################################################
# Logging Functions
###############################################################################

# Log info message (green)
log_info() {
	printf "${GREEN}[INFO]${NC} %s\n" "$1"
}

# Log warning message (yellow)
log_warn() {
	printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

# Log error message (red)
log_error() {
	printf "${RED}[ERROR]${NC} %s\n" "$1"
}

# Log debug message (blue) - only shown when VERBOSE=true
log_debug() {
	if [ "${VERBOSE:-false}" = "true" ]; then
		printf "${BLUE}[DEBUG]${NC} %s\n" "$1"
	fi
}

# Log step header (major section)
log_step() {
	echo ""
	echo "=================================================="
	echo "$1"
	echo "=================================================="
}

# Log substep header (minor section)
log_substep() {
	echo ""
	echo "--------------------------------------------------"
	echo "$1"
	echo "--------------------------------------------------"
}

###############################################################################
# General Utility Functions
###############################################################################

# Check if command exists
# Usage: command_exists docker && echo "docker is installed"
command_exists() {
	command -v "$1" >/dev/null 2>&1
}

# Check if Docker is available
# Usage: check_docker || exit 1
check_docker() {
	if ! command_exists docker; then
		log_error "docker is not installed"
		return 1
	fi

	if ! docker compose version >/dev/null 2>&1; then
		log_error "docker compose is not available"
		return 1
	fi

	return 0
}

# Wait for Docker container to be healthy
# Features:
#   - Detects and restarts exited/crashed containers
#   - Shows progress with current status and elapsed time
#   - Shows container logs on failure for debugging
# Usage: wait_for_healthy "container-name" [max_wait_seconds]
wait_for_healthy() {
	local container_name=$1
	local max_wait=${2:-120}
	local counter=0
	local last_status=""
	local restart_count=0
	local max_restarts=2

	log_info "Waiting for $container_name to be healthy (timeout: ${max_wait}s)..."

	while [ $counter -lt $max_wait ]; do
		# First check if container is running
		local running_state
		running_state=$(docker inspect --format='{{.State.Status}}' "$container_name" 2>/dev/null || echo "not_found")

		# Handle exited/crashed container
		if [ "$running_state" = "exited" ] || [ "$running_state" = "dead" ]; then
			if [ $restart_count -lt $max_restarts ]; then
				restart_count=$((restart_count + 1))
				log_warn "$container_name has exited, restarting (attempt $restart_count/$max_restarts)..."
				docker start "$container_name" >/dev/null 2>&1 || true
				sleep 3
				counter=$((counter + 3))
				continue
			else
				log_error "$container_name keeps crashing after $max_restarts restart attempts"
				log_error "Last 20 lines of container logs:"
				docker logs --tail 20 "$container_name" 2>&1 | while IFS= read -r line; do
					echo "  $line"
				done
				return 1
			fi
		fi

		# Container not found
		if [ "$running_state" = "not_found" ]; then
			log_error "Container $container_name not found"
			return 1
		fi

		# Get health status (only if container is running)
		local health_status
		health_status=$(docker inspect --format='{{.State.Health.Status}}' "$container_name" 2>/dev/null || echo "none")

		# Check if healthy
		if [ "$health_status" = "healthy" ]; then
			log_info "$container_name is healthy (took ${counter}s)"
			return 0
		fi

		# No healthcheck defined - check if container is running
		if [ "$health_status" = "none" ] || [ "$health_status" = "" ]; then
			if [ "$running_state" = "running" ]; then
				log_warn "$container_name has no healthcheck, assuming ready"
				return 0
			fi
		fi

		# Show progress every 10 seconds
		if [ $((counter % 10)) -eq 0 ]; then
			log_info "  [${counter}s] $container_name: $health_status"
		fi

		counter=$((counter + 1))
		sleep 1
	done

	log_error "$container_name did not become healthy within ${max_wait}s"
	log_error "Last 10 lines of container logs:"
	docker logs --tail 10 "$container_name" 2>&1 | while IFS= read -r line; do
		echo "  $line"
	done

	return 1
}

# Wait for a port to be available
# Usage: wait_for_port "localhost" 8080 [max_wait_seconds]
wait_for_port() {
	local host=$1
	local port=$2
	local max_wait=${3:-30}
	local counter=0

	log_info "Waiting for ${host}:${port} to be available..."

	while [ $counter -lt $max_wait ]; do
		if nc -z "$host" "$port" 2>/dev/null; then
			log_info "${host}:${port} is available"
			return 0
		fi

		counter=$((counter + 1))
		sleep 1
	done

	log_error "${host}:${port} did not become available within ${max_wait}s"
	return 1
}

# Wait for migration containers to complete
# Usage: wait_for_migrations [max_wait_seconds]
wait_for_migrations() {
	local max_wait=${1:-120}
	local counter=0
	local migration_containers=("wallet-db-migrate-watch" "wallet-db-migrate-keygen" "wallet-db-migrate-sign" "wallet-db-migrate-sign2")

	log_info "Waiting for database migrations to complete..."

	while [ $counter -lt $max_wait ]; do
		local all_completed=true

		for container in "${migration_containers[@]}"; do
			# Check if container exists first
			if ! docker ps -a --format '{{.Names}}' | grep -q "^${container}$"; then
				continue
			fi

			# Get container state
			local state
			state=$(docker inspect --format='{{.State.Status}}' "$container" 2>/dev/null || echo "not_found")

			if [ "$state" = "exited" ]; then
				# Check exit code
				local exit_code
				exit_code=$(docker inspect --format='{{.State.ExitCode}}' "$container" 2>/dev/null || echo "1")
				if [ "$exit_code" != "0" ]; then
					log_error "Migration container $container failed with exit code $exit_code"
					log_error "Container logs:"
					docker logs "$container" 2>&1 | tail -20
					return 1
				fi
			elif [ "$state" = "running" ] || [ "$state" = "created" ]; then
				all_completed=false
			fi
		done

		if [ "$all_completed" = true ]; then
			log_info "All database migrations completed successfully"
			return 0
		fi

		counter=$((counter + 1))
		sleep 1
	done

	log_error "Database migrations did not complete within ${max_wait}s"
	return 1
}

# Extract filename from command output
# Usage: filename=$(extract_filename "$output")
# Expects format: "[fileName]: /path/to/file"
extract_filename() {
	local output="$1"
	echo "${output##*\[fileName\]: }"
}

###############################################################################
# Bitcoin/Bitcoin Cash Common Utility Functions
###############################################################################

# Check if wallet exists in Bitcoin/BCH node
# Usage: bitcoin_wallet_exists "btc-watch" "watch" && echo "exists"
# Note: Works for both BTC and BCH as they share the same RPC interface
bitcoin_wallet_exists() {
	local container=$1
	local wallet_name=$2
	local rpc_user="${RPC_USER:-xyz}"
	local rpc_password="${RPC_PASSWORD:-xyz}"

	docker exec "$container" bitcoin-cli -regtest \
		-rpcuser="${rpc_user}" \
		-rpcpassword="${rpc_password}" \
		listwallets 2>/dev/null | grep -q "\"$wallet_name\"" && return 0 || return 1
}

# Create Bitcoin/BCH wallet if not exists
# Usage: bitcoin_create_wallet_if_needed "btc-watch" "watch"
# Note: Works for both BTC and BCH as they share the same RPC interface
bitcoin_create_wallet_if_needed() {
	local container=$1
	local wallet_name=$2
	local rpc_user="${RPC_USER:-xyz}"
	local rpc_password="${RPC_PASSWORD:-xyz}"

	if bitcoin_wallet_exists "$container" "$wallet_name"; then
		log_info "Wallet '$wallet_name' already exists in $container"
		docker exec "$container" bitcoin-cli -regtest \
			-rpcuser="${rpc_user}" \
			-rpcpassword="${rpc_password}" \
			loadwallet "$wallet_name" >/dev/null 2>&1 || true
	else
		log_info "Creating wallet '$wallet_name' in $container"
		# Bitcoin Cash Node (BCH) has different createwallet API than Bitcoin Core
		# BCH: createwallet "wallet_name" (disable_private_keys blank)
		# BTC: createwallet "wallet_name" (disable_private_keys blank passphrase avoid_reuse descriptors load_on_startup external_signer)
		if [[ "$container" == bch-* ]]; then
			# BCH: Use 3-parameter format (wallet_name, disable_private_keys, blank)
			docker exec "$container" bitcoin-cli -regtest \
				-rpcuser="${rpc_user}" \
				-rpcpassword="${rpc_password}" \
				createwallet "$wallet_name" false false >/dev/null
		else
			# BTC: Create wallet with appropriate descriptor setting
			# Watch wallet needs descriptors=true for importdescriptors RPC
			# Other wallets use descriptors=false for compatibility with importprivkey command
			# Parameters: wallet_name, disable_private_keys, blank, passphrase, avoid_reuse, descriptors
			local disable_private_keys=false
			local use_descriptors=false

			# Watch wallet: watch-only with descriptors
			if [[ "$wallet_name" == "watch" ]]; then
				disable_private_keys=true
				use_descriptors=true
			fi

			docker exec "$container" bitcoin-cli -regtest \
				-rpcuser="${rpc_user}" \
				-rpcpassword="${rpc_password}" \
				createwallet "$wallet_name" "$disable_private_keys" false "" false "$use_descriptors" >/dev/null
		fi
	fi
}

# Execute bitcoin-cli command in a container
# Usage: bitcoin_cli "btc-watch" "getblockcount"
# Note: Works for both BTC and BCH as they share the same RPC interface
bitcoin_cli() {
	local container=$1
	shift
	local rpc_user="${RPC_USER:-xyz}"
	local rpc_password="${RPC_PASSWORD:-xyz}"

	docker exec "$container" bitcoin-cli -regtest \
		-rpcuser="${rpc_user}" \
		-rpcpassword="${rpc_password}" \
		"$@"
}

###############################################################################
# Bitcoin-Specific Utility Functions (Aliases for backward compatibility)
###############################################################################

# Check if wallet exists in Bitcoin node (alias for backward compatibility)
# Usage: btc_wallet_exists "btc-watch" "watch" && echo "exists"
btc_wallet_exists() {
	bitcoin_wallet_exists "$@"
}

# Create Bitcoin wallet if not exists (alias for backward compatibility)
# Usage: btc_create_wallet_if_needed "btc-watch" "watch"
btc_create_wallet_if_needed() {
	bitcoin_create_wallet_if_needed "$@"
}

# Execute bitcoin-cli command in a container (alias for backward compatibility)
# Usage: btc_cli "btc-watch" "getblockcount"
btc_cli() {
	bitcoin_cli "$@"
}

###############################################################################
# Bitcoin Cash Specific Utility Functions (Aliases)
###############################################################################

# Check if wallet exists in BCH node
# Usage: bch_wallet_exists "bch-watch" "watch" && echo "exists"
bch_wallet_exists() {
	bitcoin_wallet_exists "$@"
}

# Create BCH wallet if not exists
# Usage: bch_create_wallet_if_needed "bch-watch" "watch"
bch_create_wallet_if_needed() {
	bitcoin_create_wallet_if_needed "$@"
}

# Execute bitcoin-cli command in BCH container
# Usage: bch_cli "bch-watch" "getblockcount"
bch_cli() {
	bitcoin_cli "$@"
}

###############################################################################
# File/Directory Utility Functions
###############################################################################

# Get project root directory
# Usage: PROJECT_ROOT=$(get_project_root)
get_project_root() {
	local script_dir
	script_dir="$(cd "$(dirname "${BASH_SOURCE[1]}")" && pwd)"
	# Navigate up from scripts/operation/{coin}/ to project root
	cd "$script_dir" && pwd | sed 's|/scripts/operation.*||'
}

# Clean files in directory except .gitkeep
# Usage: clean_dir_except_gitkeep "data/address/btc"
clean_dir_except_gitkeep() {
	local dir=$1
	if [ -d "$dir" ]; then
		find "$dir" -type f ! -name '.gitkeep' -delete 2>/dev/null || true
		log_info "Cleaned files in $dir"
	fi
}

###############################################################################
# Validation Functions
###############################################################################

# Validate required environment variable
# Usage: require_env "RPC_USER" || exit 1
require_env() {
	local var_name=$1
	if [ -z "${!var_name:-}" ]; then
		log_error "Required environment variable $var_name is not set"
		return 1
	fi
	return 0
}

# Validate file exists
# Usage: require_file "/path/to/file" || exit 1
require_file() {
	local file_path=$1
	if [ ! -f "$file_path" ]; then
		log_error "Required file not found: $file_path"
		return 1
	fi
	return 0
}

# Validate directory exists
# Usage: require_dir "/path/to/dir" || exit 1
require_dir() {
	local dir_path=$1
	if [ ! -d "$dir_path" ]; then
		log_error "Required directory not found: $dir_path"
		return 1
	fi
	return 0
}
