#!/bin/sh

# Bitcoin Node Setup Script
# This script initializes Bitcoin nodes when Docker containers start.
# It configures the Bitcoin nodes (btc-watch, btc-keygen, btc-sign1, btc-sign2)
# to a state where they can be operated by the go-crypto-wallet CLI.

set -e
set -u

# Configuration
COMPOSE_FILE="compose.btc.yaml"
RPC_USER="${BTC_RPC_USER:-xyz}"
RPC_PASSWORD="${BTC_RPC_PASSWORD:-xyz}"
MAX_RETRIES=30
RETRY_INTERVAL=2

# Node names
NODES="btc-watch btc-keygen btc-sign1 btc-sign2"

# Colors for output (only use if output is a TTY)
if [ -t 1 ]; then
	RED='\033[0;31m'
	GREEN='\033[0;32m'
	YELLOW='\033[1;33m'
	NC='\033[0m' # No Color
else
	RED=''
	GREEN=''
	YELLOW=''
	NC=''
fi

# Logging functions
log_info() {
	echo "${GREEN}[INFO]${NC} $1"
}

log_warn() {
	echo "${YELLOW}[WARN]${NC} $1"
}

log_error() {
	echo "${RED}[ERROR]${NC} $1"
}

# Wait for a Bitcoin node to be ready
# Args: $1 - node name (e.g., btc-watch)
wait_for_node() {
	local node_name=$1
	local retry_count=0

	log_info "Waiting for $node_name to be ready..."

	while [ $retry_count -lt $MAX_RETRIES ]; do
		if docker compose -f "$COMPOSE_FILE" exec -T "$node_name" \
			bitcoin-cli -regtest -rpcuser="$RPC_USER" -rpcpassword="$RPC_PASSWORD" \
			getblockchaininfo >/dev/null 2>&1; then
			log_info "$node_name is ready"
			return 0
		fi

		retry_count=$((retry_count + 1))
		log_warn "$node_name not ready yet (attempt $retry_count/$MAX_RETRIES)"
		sleep $RETRY_INTERVAL
	done

	log_error "$node_name failed to become ready after $MAX_RETRIES attempts"
	return 1
}

# Create a wallet on a Bitcoin node
# Args: $1 - node name, $2 - wallet name
create_wallet() {
	local node_name=$1
	local wallet_name=$2

	log_info "Creating wallet '$wallet_name' on $node_name..."

	# Try to create the wallet and capture output
	local output
	output=$(docker compose -f "$COMPOSE_FILE" exec -T "$node_name" \
		bitcoin-cli -regtest -rpcuser="$RPC_USER" -rpcpassword="$RPC_PASSWORD" \
		createwallet "$wallet_name" 2>&1)
	local exit_code=$?

	if [ $exit_code -eq 0 ]; then
		log_info "Wallet '$wallet_name' created successfully on $node_name"
		return 0
	# Check if the error is because the wallet already exists
	elif echo "$output" | grep -q "already exists"; then
		log_warn "Wallet '$wallet_name' already exists on $node_name (skipping)"
		return 0
	else
		log_error "Failed to create wallet '$wallet_name' on $node_name. Error: $output"
		return 1
	fi
}

# Generate initial blocks for regtest mode
# Args: $1 - node name, $2 - number of blocks
generate_blocks() {
	local node_name=$1
	local num_blocks=$2

	log_info "Generating $num_blocks blocks on $node_name..."

	# Check current block height
	local current_height
	current_height=$(docker compose -f "$COMPOSE_FILE" exec -T "$node_name" \
		bitcoin-cli -regtest -rpcuser="$RPC_USER" -rpcpassword="$RPC_PASSWORD" \
		getblockcount 2>/dev/null || echo "0")

	if [ "$current_height" -ge "$num_blocks" ]; then
		log_warn "Already have $current_height blocks on $node_name (skipping generation)"
		return 0
	fi

	# Generate blocks to a temporary address
	local address
	address=$(docker compose -f "$COMPOSE_FILE" exec -T "$node_name" \
		bitcoin-cli -regtest -rpcuser="$RPC_USER" -rpcpassword="$RPC_PASSWORD" \
		-rpcwallet=watch getnewaddress 2>/dev/null)

	if [ -z "$address" ]; then
		log_error "Failed to get new address for block generation"
		return 1
	fi

	if docker compose -f "$COMPOSE_FILE" exec -T "$node_name" \
		bitcoin-cli -regtest -rpcuser="$RPC_USER" -rpcpassword="$RPC_PASSWORD" \
		generatetoaddress "$num_blocks" "$address" >/dev/null 2>&1; then
		log_info "Generated $num_blocks blocks successfully on $node_name"
		return 0
	else
		log_error "Failed to generate blocks on $node_name"
		return 1
	fi
}

# Main setup logic
main() {
	log_info "Starting Bitcoin node setup..."

	# Step 1: Wait for all nodes to be ready
	log_info "Step 1: Waiting for all nodes to be ready..."
	for node in $NODES; do
		if ! wait_for_node "$node"; then
			log_error "Setup failed: $node did not become ready"
			exit 1
		fi
	done
	log_info "All nodes are ready"

	# Step 2: Create wallets
	log_info "Step 2: Creating wallets..."

	# Create wallets for each node
	local wallets="btc-watch:watch btc-keygen:keygen btc-sign1:sign1 btc-sign2:sign2"
	for item in $wallets; do
		local node_name=${item%%:*}
		local wallet_name=${item##*:}
		if ! create_wallet "$node_name" "$wallet_name"; then
			log_error "Setup failed: could not create '$wallet_name' wallet on $node_name"
			exit 1
		fi
	done

	log_info "All wallets created successfully"

	# Step 3: Generate initial blocks for regtest (optional)
	log_info "Step 3: Generating initial blocks for regtest mode..."
	if ! generate_blocks "btc-watch" 101; then
		log_warn "Block generation failed, but continuing (this is optional)"
	fi

	log_info "Bitcoin node setup completed successfully!"
	log_info "You can now use the go-crypto-wallet CLI to operate the nodes"
}

# Run main function
main "$@"
