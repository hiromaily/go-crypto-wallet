#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern 1: P2PKH Single-sig
# This script automates the complete Bitcoin workflow for single-sig P2PKH transactions
# Usage: ./scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh [OPTIONS]
# Options:
#   --cleanup  Stop containers and cleanup state
#   --reset    Full reset and run from scratch
#   --verbose  Enable verbose output
#   --non-interactive  Run without prompts (for CI/CD)
#   -h, --help Display help message
#
# Reference Documentation:
#   docs/crypto/btc/e2e_transaction_patterns.md - E2E transaction patterns
#
# Transaction Pattern:
#   Pattern 1: BTC P2PKH Single-sig
#   - Address Type: P2PKH (BIP44 Legacy)
#   - Address Format: `1...` (Mainnet), `m.../n...` (Testnet/Regtest)
#   - Signature Requirement: Single-sig (Keygen only)
#   - Descriptor: pkh([fingerprint/44'/0'/0']xpub.../0/*)
#
# Required Config Settings:
#   - config/wallet/btc_watch.yaml:  address_type: "legacy"
#   - config/wallet/btc_keygen.yaml: address_type: "legacy"

set -eu

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"

# Source common utilities
# shellcheck source=../../common.sh
source "${SCRIPT_DIR}/../../common.sh"

# Configuration
COIN="btc"
ENCRYPTED="false"
SIGN_WALLET_NUM=0  # Single-sig: no additional sign wallets needed
VERBOSE=false
CLEANUP_ONLY=false
NON_INTERACTIVE=false
RESET_STATE=false

# RPC credentials (can be overridden via environment variables)
# Note: Default values are for regtest/development only
RPC_USER="${RPC_USER:-xyz}"
RPC_PASSWORD="${RPC_PASSWORD:-xyz}"

# Wallet passphrase (only used if ENCRYPTED=true)
# Note: Default value is for testing only - use strong passphrase in production
WALLET_PASSPHRASE="${WALLET_PASSPHRASE:-test}"

# Config file paths (absolute)
CONFIG_WATCH="${PROJECT_ROOT}/config/wallet/btc_watch.yaml"
CONFIG_KEYGEN="${PROJECT_ROOT}/config/wallet/btc_keygen.yaml"
# Use single-sig account configuration for Pattern 1
CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account_singlesig.yaml"

# Export account config for keygen wallet (required for configuration)
export BTC_ACCOUNT_CONF="${CONFIG_ACCOUNT}"

###############################################################################
# Environment Variable Overrides for Configuration
###############################################################################
# These environment variables override config file values.
# Priority: Environment Variables > Config File > Default Values
#
# Pattern 1 (P2PKH Single-sig) requires:
#   - address_type: "legacy" (derives key_type: bip44 automatically)
# Note: key_type is automatically derived from address_type in Go code
#       (see internal/domain/address/types.go AddrType.ToKeyType())
# Note: Bitcoin RPC wallet names are set via sed in setup_wallets() function
export WALLET_ADDRESS_TYPE="legacy"

###############################################################################
# Cleanup Functions
###############################################################################

# Clean generated data files
clean_data_files() {
	log_substep "Cleaning generated data files"

	# Clean address files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/address/btc"

	# Clean fullpubkey files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/fullpubkey/btc"

	# Clean descriptor files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/descriptor/btc"

	# Clean transaction files (keeping .gitkeep)
	clean_dir_except_gitkeep "data/tx/btc"
}

# Clean Bitcoin node wallet data
clean_bitcoin_wallet_data() {
	log_substep "Cleaning Bitcoin node wallet data"

	# Clean regtest data directories (keeping bitcoin.conf)
	for node in watch keygen; do
		wallet_dir="docker/nodes/btc/${node}/regtest"
		if [ -d "$wallet_dir" ]; then
			# Remove all files/dirs except bitcoin.conf
			find "$wallet_dir" -mindepth 1 ! -name 'bitcoin.conf' -exec rm -rf {} + 2>/dev/null || true
			log_info "Cleaned ${node} wallet data"
		fi
	done
}

# Full reset: cleanup everything for fresh state
full_reset() {
	log_step "Performing full reset for fresh state"

	# Stop and remove containers WITH VOLUMES
	# The -v flag removes volumes, which clears the database state
	log_info "Stopping Bitcoin containers (with volume removal)..."
	docker compose -f compose.btc.yaml down -v 2>/dev/null || true

	log_info "Stopping database container (with volume removal)..."
	# This removes all database data, ensuring a clean slate
	docker compose -f compose.yaml down -v 2>/dev/null || true

	# Wait for containers to fully stop before attempting volume removal
	log_info "Waiting for containers to stop completely..."
	sleep 3

	# Explicitly remove database volume (in case -v flag didn't work)
	log_info "Forcefully removing database volume..."
	local volume_name="go-crypto-wallet_wallet-db"

	# Try multiple times in case volume is still being used
	local removal_attempts=0
	local max_removal_attempts=5

	while [ $removal_attempts -lt $max_removal_attempts ]; do
		if docker volume rm "$volume_name" 2>/dev/null; then
			log_info "Volume removed successfully on attempt $((removal_attempts + 1))"
			break
		fi
		removal_attempts=$((removal_attempts + 1))
		if [ $removal_attempts -lt $max_removal_attempts ]; then
			log_warn "Volume removal failed, retrying in 2 seconds... (attempt $removal_attempts/$max_removal_attempts)"
			sleep 2
		fi
	done

	# Verify volume is truly deleted (wait up to 10 seconds)
	log_info "Verifying volume deletion..."
	local max_wait=10
	local counter=0
	local volume_deleted=false

	while [ $counter -lt $max_wait ]; do
		if ! docker volume inspect "$volume_name" >/dev/null 2>&1; then
			log_info "Volume successfully deleted"
			volume_deleted=true
			break
		fi
		counter=$((counter + 1))
		if [ $counter -lt $max_wait ]; then
			log_warn "Volume still exists, waiting... (${counter}s/${max_wait}s)"
			sleep 1
		fi
	done

	if [ "$volume_deleted" = "false" ]; then
		log_error "Volume still exists after ${max_wait}s - this may cause duplicate key errors"
		log_error "Manual cleanup required: docker volume rm -f $volume_name"
		return 1
	fi

	# Clean data files
	clean_data_files

	# Clean Bitcoin wallet data
	clean_bitcoin_wallet_data

	log_info "Full reset complete - system is in fresh state"
	log_info "Note: Database volumes were removed for complete cleanup"
}

# Basic cleanup: just stop containers
cleanup() {
	log_step "Cleaning up containers and state"

	# Stop and remove containers
	log_info "Stopping Bitcoin containers..."
	docker compose -f compose.btc.yaml down -v 2>/dev/null || true

	log_info "Stopping database container..."
	docker compose -f compose.yaml down -v 2>/dev/null || true

	log_info "Cleanup complete"
}

###############################################################################
# Prerequisites Check
###############################################################################

check_prerequisites() {
	log_step "Checking prerequisites"

	# Check Docker and Docker Compose
	check_docker || exit 1

	# Check CLI commands (only watch and keygen for single-sig)
	for cmd in watch keygen; do
		if ! command_exists "$cmd"; then
			log_error "$cmd command is not available"
			log_error "Please build the project first: make build"
			exit 1
		fi
	done

	log_info "All prerequisites satisfied"
}

###############################################################################
# Infrastructure Setup
###############################################################################

setup_infrastructure() {
	log_step "Setting up infrastructure"

	# Start database
	log_substep "Starting database container"
	docker compose -f compose.yaml up -d
	log_info "Database container started"

	# Wait for database to be healthy
	wait_for_healthy "wallet-db" 90

	# Start Bitcoin nodes (only watch and keygen for single-sig)
	log_substep "Starting Bitcoin node containers"
	docker compose -f compose.btc.yaml up -d btc-watch btc-keygen
	log_info "Bitcoin node containers started"

	# Wait for containers to be healthy
	log_substep "Waiting for containers to be healthy"
	wait_for_healthy "btc-watch" 90
	wait_for_healthy "btc-keygen" 90

	log_info "All containers are healthy"
}

###############################################################################
# Wallet Setup
###############################################################################

setup_wallets() {
	log_step "Setting up Bitcoin wallets"

	# Create wallets in Bitcoin nodes (only watch and keygen for single-sig)
	btc_create_wallet_if_needed "btc-watch" "watch"
	btc_create_wallet_if_needed "btc-keygen" "keygen"

	log_info "All wallets are ready"

	# Note: Wallet-specific RPC endpoints are configured via environment variables
	# for each command invocation (see wrapper functions below)
	# This avoids modifying config files and creating backups
	# Environment variables use WALLET_ prefix to override config file settings
}

###############################################################################
# Wallet Command Wrappers with Environment Variable Overrides
###############################################################################

# Wrapper for watch wallet commands with host override
watch_with_wallet() {
	WALLET_BITCOIN_HOST="127.0.0.1:18332/wallet/watch" watch "$@"
}

# Wrapper for keygen wallet commands with host override
keygen_with_wallet() {
	WALLET_BITCOIN_HOST="127.0.0.1:19332/wallet/keygen" keygen "$@"
}

###############################################################################
# Key Generation Phase
###############################################################################

key_generation_phase() {
	log_step "Key Generation Phase"

	# Keygen wallet - create seed
	log_substep "Creating seed for keygen wallet"
	keygen -c "${CONFIG_KEYGEN}" create seed || {
		log_warn "Seed already exists or error occurred, continuing..."
	}

	# Keygen wallet - create hdkeys
	log_substep "Creating HD keys for keygen wallet (client, deposit, payment, stored)"
	for account in client deposit payment stored; do
		log_info "Creating HD keys for account: $account"
		keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" create hdkey --account "$account" --keynum 10
	done

	# Keygen wallet - import private keys
	# NOTE: Private key import is NOT needed for descriptor-based wallets (Pattern 1)
	# The keygen wallet will sign PSBTs using its internal descriptor wallet
	# Skipping this step to avoid "Cannot import private keys to a wallet with private keys disabled" error
	# log_substep "Importing private keys into keygen wallet"
	# if [ "$ENCRYPTED" = "true" ]; then
	# 	keygen -c "${CONFIG_KEYGEN}" api walletpassphrase --passphrase "${WALLET_PASSPHRASE}"
	# fi
	# for account in client deposit payment stored; do
	# 	log_info "Importing private keys for account: $account"
	# 	keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" import privkey --account "$account"
	# done
	# if [ "$ENCRYPTED" = "true" ]; then
	# 	keygen -c "${CONFIG_KEYGEN}" api walletlock
	# fi
}

###############################################################################
# Single-sig Address Setup Phase
###############################################################################

singlesig_setup_phase() {
	log_step "Single-sig Address Setup Phase"

	# For Pattern 1 (P2PKH Single-sig), use account_singlesig.yaml
	# which configures all accounts as single-sig
	# We export descriptors for both client and payment accounts
	# to test the complete deposit → payment flow
	log_substep "Exporting descriptors from keygen wallet"

	# Process client and payment accounts for single-sig testing
	local accounts=(client payment)

	# Export descriptors for accounts
	# Note: Using wrapper function to set WALLET_BITCOIN_HOST for Bitcoin Core RPC
	declare -A descriptor_files
	for account in "${accounts[@]}"; do
		log_info "Exporting ${account} descriptors"
		file_output=$(keygen_with_wallet -c "${CONFIG_KEYGEN}" --coin "${COIN}" descriptor export \
			--account "${account}" \
			--output "data/descriptor/btc/${account}_descriptors.json" \
			--format bitcoin-core \
			--include-change)
		# Extract file path from output
		descriptor_files[$account]="${file_output##*exported to }"
		log_info "  ${account}: ${descriptor_files[$account]}"
	done

	# Import descriptors into watch wallet
	# Note: Using wrapper function to set WALLET_BITCOIN_HOST for Bitcoin Core RPC
	log_substep "Importing descriptors into watch wallet"
	for account in "${accounts[@]}"; do
		log_info "Importing ${account} descriptors"
		watch_with_wallet -c "${CONFIG_WATCH}" --coin "${COIN}" import descriptor \
			--file "${descriptor_files[$account]}" \
			--account "${account}"
	done

	log_info "All descriptors imported successfully"
	log_info "Note: Pattern 1 uses account_singlesig.yaml (all accounts are single-sig)"

	# Derive payment address from descriptor for UTXO generation
	# We generate UTXOs to the payment account address for testing payment transactions
	log_substep "Deriving payment address from descriptor for UTXO generation"

	# Extract first descriptor from payment_descriptors.json
	# For P2PKH (legacy), we use the first descriptor (index 0)
	first_descriptor=$(jq -r '.[0].desc // empty' "${descriptor_files[payment]}" 2>/dev/null)

	if [ -z "$first_descriptor" ]; then
		log_error "Failed to extract descriptor from ${descriptor_files[payment]}"
		return 1
	fi

	log_info "  Using descriptor: ${first_descriptor:0:50}..."

	# Export for use in generate_test_utxos
	export first_descriptor
}

###############################################################################
# UTXO Generation Phase (for regtest)
###############################################################################

generate_test_utxos() {
	log_step "Generating Test UTXOs for Transaction Phase"

	log_info "Deriving payment address from descriptor..."

	# Add checksum to descriptor if not present
	if ! echo "$first_descriptor" | grep -q '#'; then
		log_info "Adding checksum to descriptor..."
		descriptor_with_checksum=$(btc_cli "btc-watch" getdescriptorinfo "$first_descriptor" | jq -r '.descriptor')
	else
		descriptor_with_checksum="$first_descriptor"
	fi

	# Derive first address (index 0) from descriptor
	payment_address=$(btc_cli "btc-watch" deriveaddresses "$descriptor_with_checksum" "[0,0]" 2>/dev/null | jq -r '.[0]')

	if [ -z "$payment_address" ] || [ "$payment_address" = "null" ]; then
		log_error "Failed to derive payment address from descriptor"
		log_error "Descriptor: $descriptor_with_checksum"
		return 1
	fi

	log_info "Using payment address: $payment_address"

	# Generate blocks with coinbase reward to payment address
	log_info "Generating 101 blocks to create mature coinbase for testing..."
	btc_cli "btc-watch" generatetoaddress 101 "$payment_address" >/dev/null

	log_info "Test UTXOs generated successfully"
	log_info "Waiting for blockchain sync and balance update..."

	# Poll for balance update with timeout
	max_wait=60
	wait_interval=3
	elapsed=0
	balance_found=false

	while [ $elapsed -lt $max_wait ]; do
		# Check balance using Bitcoin Core RPC directly
		balance_json=$(btc_cli "btc-watch" -rpcwallet=watch getbalances 2>&1 || true)
		trusted_balance=$(echo "$balance_json" | jq -r '.mine.trusted // 0' 2>/dev/null || echo "0")

		# Check if we have any trusted (mature) balance
		if [ -n "$trusted_balance" ] && (($(echo "$trusted_balance > 0" | bc -l))); then
			log_info "Payment account balance verified: ${trusted_balance} BTC (took ${elapsed}s)"
			balance_found=true
			break
		fi

		sleep $wait_interval
		elapsed=$((elapsed + wait_interval))
		if [ $elapsed -lt $max_wait ]; then
			log_info "Still waiting for balance update... (${elapsed}s/${max_wait}s)"
		fi
	done

	if [ "$balance_found" = false ]; then
		log_error "Balance not detected within ${max_wait}s"
		log_error "This indicates a failure in UTXO generation or blockchain sync"
		log_error "Please check:"
		log_error "  - Bitcoin node logs: docker compose -f compose.btc.yaml logs btc-watch"
		log_error "  - Block generation succeeded"
		log_error "  - Address import into watch wallet succeeded"
		log_error "Debug: Last balance check output:"
		echo "$balance_json"
		return 1
	fi
}

###############################################################################
# Payment Request Creation Phase
###############################################################################

create_payment_requests_phase() {
	log_step "Payment Request Creation Phase"

	# Get a payment sender address from database
	# For P2PKH testing, query for legacy addresses (starting with 'm' or 'n' in regtest)
	log_substep "Retrieving payment sender address from database"
	sender_address=$(docker compose exec -T wallet-db mysql -u root -proot watch -N -e \
		"SELECT wallet_address FROM address WHERE coin='btc' AND account='payment' LIMIT 1" 2>/dev/null)

	if [ -z "$sender_address" ]; then
		log_error "No payment addresses found in database"
		log_error "Please check:"
		log_error "  - Descriptor import succeeded"
		log_error "  - Addresses were derived and stored in database"
		return 1
	fi

	log_info "Using sender address: $sender_address"

	# Generate anonymous receiver addresses for testing
	log_substep "Generating receiver addresses for payment requests"
	log_info "Creating 3 receiver addresses in watch wallet..."
	receiver1=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" legacy)
	receiver2=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" legacy)
	receiver3=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" legacy)

	log_info "Generated receiver addresses:"
	log_info "  1. $receiver1"
	log_info "  2. $receiver2"
	log_info "  3. $receiver3"

	# Create payment requests using payment account
	log_substep "Inserting payment requests into database"
	docker compose exec -T wallet-db mysql -u root -proot watch <<EOF
DELETE FROM payment_request;
INSERT INTO payment_request (coin, payment_id, sender_address, sender_account, receiver_address, amount, is_done)
VALUES
	('btc', NULL, '${sender_address}', 'payment', '${receiver1}', 0.001, false),
	('btc', NULL, '${sender_address}', 'payment', '${receiver2}', 0.002, false),
	('btc', NULL, '${sender_address}', 'payment', '${receiver3}', 0.0015, false);
EOF

	if [ $? -ne 0 ]; then
		log_error "Failed to insert payment requests"
		return 1
	fi

	# Verify payment requests were created
	count=$(docker compose exec -T wallet-db mysql -u root -proot watch -N -e \
		"SELECT COUNT(*) FROM payment_request WHERE coin='btc' AND is_done=false" 2>/dev/null)

	log_info "Created $count payment requests"

	if [ "$count" -eq 0 ]; then
		log_error "No payment requests were created"
		return 1
	fi

	log_info "Payment requests ready for transaction creation"
}

###############################################################################
# Transaction Flow Phase (Single-sig)
###############################################################################

transaction_flow_phase() {
	log_step "Transaction Flow Phase (Single-sig)"

	# Create unsigned transaction
	log_substep "Creating unsigned payment transaction"
	tx_file=$(watch_with_wallet -c "${CONFIG_WATCH}" create payment 2>&1) || {
		log_error "Failed to create payment transaction"
		log_error "Output: $tx_file"

		if echo "$tx_file" | grep -q "No utxo"; then
			log_error "Transaction creation failed"
			log_error "This could indicate:"
			log_error "  - No payment requests in database"
			log_error "  - No UTXOs available for payment account"
			log_error "  - UTXOs not mature enough (need 100+ confirmations)"
			return 1
		fi

		return 1
	}

	if echo "$tx_file" | grep -q "No utxo"; then
		log_error "Transaction creation failed"
		log_error "This could indicate:"
		log_error "  - No payment requests in database"
		log_error "  - No UTXOs available for payment account"
		log_error "  - UTXOs not mature enough (need 100+ confirmations)"
		return 1
	fi

	# Extract file path
	tx_unsigned=$(echo "${tx_file}" | sed -n 's/.*\[fileName\]: //p')
	log_info "Created unsigned transaction: $tx_unsigned"

	# Sign with keygen wallet (single signature)
	log_substep "Signing with keygen wallet (single signature)"
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletpassphrase --passphrase "${WALLET_PASSPHRASE}"
	fi
	tx_file_signed=$(keygen -c "${CONFIG_KEYGEN}" sign signature --file "${tx_unsigned}")
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletlock
	fi

	tx_signed=$(echo "${tx_file_signed}" | sed -n 's/.*\[fileName\]: //p')
	log_info "Signed transaction: $tx_signed"

	# Send transaction
	log_substep "Sending fully signed transaction"
	tx_result=$(watch_with_wallet -c "${CONFIG_WATCH}" send --file "${tx_signed}")
	tx_id="${tx_result##*txID: }"

	log_info "Transaction sent successfully!"
	log_info "Transaction ID: $tx_id"
}

###############################################################################
# Help Message
###############################################################################

show_help() {
	cat <<EOF
Bitcoin E2E Workflow Script - Pattern 1: P2PKH Single-sig

This script automates the complete Bitcoin workflow for single-signature P2PKH
transactions. It serves as a regression test tool to verify that the Bitcoin
single-sig workflow functions correctly after code changes.

Usage: $0 [OPTIONS]

Options:
  --reset             Full reset: cleanup all state for fresh start
  --cleanup           Stop containers and cleanup state, then exit
  --verbose           Enable verbose output
  --non-interactive   Run without interactive prompts (for CI/CD)
  -h, --help          Display this help message

Examples:
  # Run from completely fresh state
  $0 --reset

  # Run complete E2E workflow
  $0

  # Run with verbose output
  $0 --verbose

  # Cleanup containers and state
  $0 --cleanup

The script performs the following steps:
  1. Check prerequisites (Docker, CLI commands)
  2. Start infrastructure (database and Bitcoin nodes)
  3. Create wallets in Bitcoin nodes
  4. Generate keys for keygen wallet
  5. Export descriptors from keygen wallet
  6. Import descriptors into watch wallet
  7. Generate test UTXOs (automatically generates 101 blocks)
  8. Create payment requests
  9. Create unsigned transaction
 10. Sign with keygen wallet (single signature)
 11. Broadcast transaction

The script uses descriptor-based import for P2PKH single-sig addresses,
ensuring compatibility with Bitcoin Core's modern wallet infrastructure.
Test UTXOs are automatically generated for the transaction phase, making it
fully automated and suitable for CI/CD pipelines.

Transaction Pattern Details:
  - Address Type: P2PKH (BIP44 Legacy)
  - Address Format: m.../n... (Regtest)
  - Signature Requirement: Single-sig (Keygen only)
  - No sign wallets needed

Environment Variables:
  RPC_USER          Bitcoin RPC username (default: xyz for regtest)
  RPC_PASSWORD      Bitcoin RPC password (default: xyz for regtest)
  WALLET_PASSPHRASE Wallet passphrase for encrypted wallets (default: test)

EOF
}

###############################################################################
# Main Execution
###############################################################################

main() {
	# Parse arguments
	while [ $# -gt 0 ]; do
		case "$1" in
		--cleanup)
			CLEANUP_ONLY=true
			shift
			;;
		--reset)
			RESET_STATE=true
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
			exit 0
			;;
		*)
			log_error "Unknown option: $1"
			show_help
			exit 1
			;;
		esac
	done

	# Cleanup and exit if requested
	if [ "$CLEANUP_ONLY" = "true" ]; then
		cleanup
		exit 0
	fi

	# Full reset if requested
	if [ "$RESET_STATE" = "true" ]; then
		full_reset
	fi

	log_info "Starting Bitcoin E2E Workflow - Pattern 1: P2PKH Single-sig"
	log_info "Coin: $COIN"
	log_info "Encrypted: $ENCRYPTED"
	log_info "Sign wallet count: $SIGN_WALLET_NUM (Single-sig)"
	echo ""

	# Execute workflow phases
	check_prerequisites
	setup_infrastructure
	setup_wallets
	key_generation_phase
	singlesig_setup_phase
	generate_test_utxos
	create_payment_requests_phase
	transaction_flow_phase

	log_step "Bitcoin E2E Workflow Completed Successfully!"
	log_info "Summary:"
	log_info "  ✓ Infrastructure setup complete"
	log_info "  ✓ Wallets created and configured"
	log_info "  ✓ HD keys generated for keygen wallet"
	log_info "  ✓ Descriptors exported and imported (client and payment accounts)"
	log_info "  ✓ P2PKH single-sig addresses created"
	log_info "  ✓ Test UTXOs generated"
	log_info "  ✓ Payment requests created (using payment account)"
	log_info "  ✓ Transaction created, signed (1 signature), and sent"
	echo ""
	log_info "Transaction Pattern Used:"
	log_info "  • P2PKH (BIP44 Legacy) single-signature"
	log_info "  • Descriptor-based address management"
	log_info "  • Simple single-key signing workflow"
	log_info "  • Uses account_singlesig.yaml (all accounts configured as single-sig)"
	echo ""
	log_info "You can now use the wallet system for Bitcoin single-sig operations"
	log_info "To cleanup, run: $0 --cleanup"
	log_info "To full reset for fresh state, run: $0 --reset"
}

# Trap errors and cleanup
trap 'log_error "Script failed at line $LINENO"' ERR

# Run main
main "$@"
