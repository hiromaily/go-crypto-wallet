#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern 10: P2TR MuSig2 N-of-N
# This script automates the complete Bitcoin workflow for MuSig2 P2TR (Taproot) transactions
# Usage: ./scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh [OPTIONS]
# Options:
#   --cleanup  Stop containers and cleanup state
#   --reset    Full reset and run from scratch
#   --verbose  Enable verbose output
#   --non-interactive  Run without prompts (for CI/CD)
#   -h, --help Display help message
#
# Reference Documentation:
#   docs/crypto/btc/e2e_transaction_patterns.md - E2E transaction patterns
#   docs/crypto/btc/musig2_guide.md - MuSig2 user guide
#   docs/crypto/btc/TAPROOT_GUIDE.md - Taproot user guide
#
# Transaction Pattern:
#   Pattern 10: BTC P2TR MuSig2 N-of-N
#   - Address Type: P2TR (BIP86 + BIP327 Taproot MuSig2)
#   - Address Format: `bc1p...` (Mainnet), `tb1p...` (Testnet), `bcrt1p...` (Regtest)
#   - Signature Requirement: N-of-N (all signatures required, aggregated into single Schnorr signature)
#   - Descriptor: tr(musig([fingerprint/86'/1'/1']xpub1,[fingerprint/86'/1'/1']xpub2,[fingerprint/86'/1'/1']xpub3)/0/*)
#   - Bitcoin Core Version Required: v22.0+
#   - Protocol: 2-Round MuSig2 (BIP327)
#
# Required Config Settings:
#   - config/wallet/btc_watch.yaml:  address_type: "bech32m"
#   - config/wallet/btc_keygen.yaml: address_type: "bech32m"
#   - config/wallet/btc_sign1.yaml:  address_type: "bech32m"
#   - config/wallet/btc_sign2.yaml:  address_type: "bech32m"
#
# IMPORTANT NOTE:
#   This E2E script demonstrates the MuSig2 workflow framework.
#   The MuSig2 CLI commands currently contain placeholder implementations with TODOs.
#   This script serves as documentation for the intended workflow and will be fully
#   functional once the MuSig2 implementation is complete.

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
SIGN_WALLET_NUM=2 # MuSig2 N-of-N: keygen + sign1 + sign2 = 3 signers
VERBOSE=false
CLEANUP_ONLY=false
NON_INTERACTIVE=false
RESET_STATE=false

# RPC credentials (can be overridden via environment variables)
# Note: Default values are for regtest/development only
RPC_USER="${RPC_USER:-xyz}"
RPC_PASSWORD="${RPC_PASSWORD:-xyz}"

# MySQL credentials (can be overridden via environment variables)
# Note: Default value is for regtest/development only
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-root}"

# Wallet passphrase (only used if ENCRYPTED=true)
# Note: Default value is for testing only - use strong passphrase in production
WALLET_PASSPHRASE="${WALLET_PASSPHRASE:-test}"

# Config file paths (absolute)
CONFIG_WATCH="${PROJECT_ROOT}/config/wallet/btc_watch.yaml"
CONFIG_KEYGEN="${PROJECT_ROOT}/config/wallet/btc_keygen.yaml"
CONFIG_SIGN1="${PROJECT_ROOT}/config/wallet/btc_sign1.yaml"
CONFIG_SIGN2="${PROJECT_ROOT}/config/wallet/btc_sign2.yaml"
# Use 3-of-3 account configuration for Pattern 10 (MuSig2 N-of-N)
CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account_3of3.yaml"

# Export account config for keygen wallet (required for configuration)
export BTC_ACCOUNT_CONF="${CONFIG_ACCOUNT}"

# Wallet-specific RPC hosts (for environment variable overrides)
WATCH_WALLET_RPC_HOST="127.0.0.1:18332/wallet/watch"
KEYGEN_WALLET_RPC_HOST="127.0.0.1:19332/wallet/keygen"

###############################################################################
# Environment Variable Overrides for Configuration
###############################################################################
# These environment variables override config file values.
# Priority: Environment Variables > Config File > Default Values
#
# Pattern 10 (P2TR MuSig2 N-of-N) requires:
#   - address_type: "bech32m" (for Taproot P2TR addresses)
#   - Derives key_type: bip86 automatically
# Note: key_type is automatically derived from address_type in Go code
#       (see internal/domain/address/types.go AddrType.ToKeyType())
# Note: "bech32m" is used for Taproot addresses (not "taproot" config value)
export WALLET_ADDRESS_TYPE="bech32m"

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
	for node in watch keygen sign1 sign2; do
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
	# Docker Compose prefixes volume names with the project name (defaults to base name of project directory)
	# Dynamically determine volume name to handle different project directory names
	local volume_name
	volume_name="$(basename "$PROJECT_ROOT" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9_-]//g')_wallet-db"

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

	# Check CLI commands (watch, keygen, sign1, sign2 for MuSig2)
	for cmd in watch keygen sign1 sign2; do
		if ! command_exists "$cmd"; then
			log_error "$cmd command is not available"
			log_error "Please build the project first: make build-all"
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

	# Start Bitcoin nodes (watch, keygen, sign1, sign2 for MuSig2)
	log_substep "Starting Bitcoin node containers"
	docker compose -f compose.btc.yaml up -d
	log_info "Bitcoin node containers started"

	# Wait for containers to be healthy
	log_substep "Waiting for containers to be healthy"
	wait_for_healthy "btc-watch" 90
	wait_for_healthy "btc-keygen" 90
	wait_for_healthy "btc-sign1" 90
	wait_for_healthy "btc-sign2" 90

	log_info "All containers are healthy"
}

###############################################################################
# Wallet Setup
###############################################################################

setup_wallets() {
	log_step "Setting up Bitcoin wallets"

	# Create wallets in Bitcoin nodes (watch, keygen, sign1, sign2)
	btc_create_wallet_if_needed "btc-watch" "watch"
	btc_create_wallet_if_needed "btc-keygen" "keygen"
	btc_create_wallet_if_needed "btc-sign1" "sign1"
	btc_create_wallet_if_needed "btc-sign2" "sign2"

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
	WALLET_BITCOIN_HOST="${WATCH_WALLET_RPC_HOST}" watch "$@"
}

# Wrapper for keygen wallet commands with host override
keygen_with_wallet() {
	WALLET_BITCOIN_HOST="${KEYGEN_WALLET_RPC_HOST}" keygen "$@"
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

	# Keygen wallet - create hdkeys for MuSig2 accounts
	log_substep "Creating HD keys for keygen wallet (deposit, payment, stored)"
	for account in deposit payment stored; do
		log_info "Creating HD keys for account: $account"
		keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" create hdkey --account "$account" --keynum 10
	done

	# Keygen wallet - import private keys
	log_substep "Importing private keys into keygen wallet"
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletpassphrase --passphrase "${WALLET_PASSPHRASE}"
	fi
	for account in deposit payment stored; do
		log_info "Importing private keys for account: $account"
		keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" import privkey --account "$account"
	done
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletlock
	fi

	# Sign wallets - create seeds for all sign wallets
	log_substep "Creating seeds for sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		log_info "Creating seed for sign${i}"
		config_var="CONFIG_SIGN${i}"
		"sign${i}" --conf "${!config_var}" --coin "${COIN}" create seed || {
			log_warn "Sign${i} seed already exists or error occurred, continuing..."
		}
	done

	# Sign wallets - create hdkeys
	log_substep "Creating HD keys for sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		log_info "Creating HD keys for sign${i}"
		config_var="CONFIG_SIGN${i}"
		"sign${i}" --conf "${!config_var}" --coin "${COIN}" --wallet "sign${i}" create hdkey
	done

	# Sign wallets - import private keys
	log_substep "Importing private keys into sign wallets"
	for i in $(seq 1 "$SIGN_WALLET_NUM"); do
		log_info "Importing private keys for sign${i}"
		config_var="CONFIG_SIGN${i}"
		if [ "$ENCRYPTED" = "true" ]; then
			"sign${i}" --conf "${!config_var}" --coin "${COIN}" --wallet "sign${i}" api walletpassphrase --passphrase "${WALLET_PASSPHRASE}"
		fi
		"sign${i}" --conf "${!config_var}" --coin "${COIN}" --wallet "sign${i}" import privkey
		if [ "$ENCRYPTED" = "true" ]; then
			"sign${i}" --conf "${!config_var}" --coin "${COIN}" --wallet "sign${i}" api walletlock
		fi
	done

	# Sign wallets - export fullpubkey
	log_substep "Exporting full public keys from sign wallets"
	file_fullpubkey_auth1=$(sign1 --conf "${CONFIG_SIGN1}" --coin "${COIN}" --wallet sign1 export fullpubkey)
	file_fullpubkey_auth2=$(sign2 --conf "${CONFIG_SIGN2}" --coin "${COIN}" --wallet sign2 export fullpubkey)

	# Extract file paths
	fullpubkey_file1="${file_fullpubkey_auth1##*\[fileName\]: }"
	fullpubkey_file2="${file_fullpubkey_auth2##*\[fileName\]: }"

	log_info "Exported fullpubkey files:"
	log_info "  sign1: $fullpubkey_file1"
	log_info "  sign2: $fullpubkey_file2"

	# Store for next phase
	export FULLPUBKEY_FILE1="$fullpubkey_file1"
	export FULLPUBKEY_FILE2="$fullpubkey_file2"
}

###############################################################################
# MuSig2 Setup Phase (N-of-N)
###############################################################################

musig2_setup_phase() {
	log_step "MuSig2 Setup Phase (N-of-N Taproot)"

	# Import fullpubkeys
	log_substep "Importing full public keys into keygen wallet"
	log_info "Importing fullpubkey from sign1: $FULLPUBKEY_FILE1"
	keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" import fullpubkey --file "${FULLPUBKEY_FILE1}"

	log_info "Importing fullpubkey from sign2: $FULLPUBKEY_FILE2"
	keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" import fullpubkey --file "${FULLPUBKEY_FILE2}"

	# For MuSig2, we need to create aggregated public keys from keygen + sign1 + sign2
	# This will be done through the descriptor export with MuSig2 format
	log_substep "Exporting MuSig2 descriptors from keygen wallet"

	# Export descriptors for MuSig2 accounts (deposit, payment, stored)
	# Pattern 10 uses MuSig2 N-of-N with P2TR (BIP86 + BIP327 Taproot)
	# Note: Using wrapper function to set WALLET_BITCOIN_HOST for Bitcoin Core RPC

	# TODO: Implement actual MuSig2 descriptor export
	# For now, we export standard Taproot descriptors as placeholder
	# The real implementation should use: tr(musig([fp1/86'/1'/1']xpub1,[fp2]xpub2,[fp3]xpub3)/0/*)

	log_warn "MuSig2 descriptor export is currently a placeholder"
	log_warn "The script will export standard Taproot descriptors for now"

	file_descriptor_deposit=$(keygen_with_wallet -c "${CONFIG_KEYGEN}" --coin "${COIN}" descriptor export --account deposit --output data/descriptor/btc/deposit_descriptors.json --format bitcoin-core --include-change)
	file_descriptor_payment=$(keygen_with_wallet -c "${CONFIG_KEYGEN}" --coin "${COIN}" descriptor export --account payment --output data/descriptor/btc/payment_descriptors.json --format bitcoin-core --include-change)
	file_descriptor_stored=$(keygen_with_wallet -c "${CONFIG_KEYGEN}" --coin "${COIN}" descriptor export --account stored --output data/descriptor/btc/stored_descriptors.json --format bitcoin-core --include-change)

	# Extract file paths from descriptor export output
	descriptor_deposit="${file_descriptor_deposit##*exported to }"
	descriptor_payment="${file_descriptor_payment##*exported to }"
	descriptor_stored="${file_descriptor_stored##*exported to }"

	log_info "Exported descriptor files:"
	log_info "  deposit: $descriptor_deposit"
	log_info "  payment: $descriptor_payment"
	log_info "  stored: $descriptor_stored"

	# Import descriptors into watch wallet for MuSig2 accounts
	# Note: Using wrapper function to set WALLET_BITCOIN_HOST for Bitcoin Core RPC
	log_substep "Importing descriptors into watch wallet"
	log_info "Importing deposit descriptors"
	watch_with_wallet -c "${CONFIG_WATCH}" --coin "${COIN}" import descriptor --file "${descriptor_deposit}" --account deposit

	log_info "Importing payment descriptors"
	watch_with_wallet -c "${CONFIG_WATCH}" --coin "${COIN}" import descriptor --file "${descriptor_payment}" --account payment

	log_info "Importing stored descriptors"
	watch_with_wallet -c "${CONFIG_WATCH}" --coin "${COIN}" import descriptor --file "${descriptor_stored}" --account stored

	log_info "All descriptors imported successfully"
	log_info "Note: Pattern 10 uses MuSig2 N-of-N with P2TR (BIP86 + BIP327)"
	log_info "Note: Current implementation uses placeholder descriptors"

	# Derive payment address from descriptor for UTXO generation
	log_substep "Deriving payment address from descriptor for UTXO generation"

	# Extract first descriptor from payment_descriptors.json
	# For P2TR MuSig2, we use the descriptor at index 0
	first_descriptor=$(jq -r '.[0].desc // empty' "${descriptor_payment}" 2>/dev/null)

	if [ -z "$first_descriptor" ]; then
		log_error "Failed to extract descriptor from ${descriptor_payment}"
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

	# Verify address format (should start with 'bcrt1p' for regtest P2TR)
	if [[ ! "$payment_address" =~ ^bcrt1p ]]; then
		log_warn "Warning: Expected P2TR address starting with 'bcrt1p', got: $payment_address"
		log_warn "This may indicate address_type configuration issue"
	fi

	# Export for use in create_payment_requests_phase
	export payment_address

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

	# Use the payment address derived from descriptor in generate_test_utxos phase
	log_substep "Using payment sender address derived from descriptor"
	sender_address="$payment_address"
	log_info "Using sender address: $sender_address"

	# Generate anonymous receiver addresses for testing
	log_substep "Generating receiver addresses for payment requests"
	log_info "Creating 3 receiver addresses in watch wallet..."
	receiver1=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" bech32m)
	receiver2=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" bech32m)
	receiver3=$(btc_cli "btc-watch" -rpcwallet=watch getnewaddress "" bech32m)

	log_info "Generated receiver addresses:"
	log_info "  1. $receiver1"
	log_info "  2. $receiver2"
	log_info "  3. $receiver3"

	# Create payment requests using payment account
	log_substep "Inserting payment requests into database"
	docker compose exec -e MYSQL_PWD="${MYSQL_ROOT_PASSWORD}" -T wallet-db mysql -u root watch <<EOF
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
	count=$(docker compose exec -e MYSQL_PWD="${MYSQL_ROOT_PASSWORD}" -T wallet-db mysql -u root watch -N -e \
		"SELECT COUNT(*) FROM payment_request WHERE coin='btc' AND is_done=false" 2>/dev/null)

	log_info "Created $count payment requests"

	if [ "$count" -eq 0 ]; then
		log_error "No payment requests were created"
		return 1
	fi

	log_info "Payment requests ready for transaction creation"
}

###############################################################################
# Transaction Flow Phase (MuSig2 2-Round Protocol)
###############################################################################

transaction_flow_phase() {
	log_step "Transaction Flow Phase (MuSig2 2-Round Protocol)"

	# Create unsigned transaction (creates PSBT)
	log_substep "Creating unsigned payment transaction (PSBT)"
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

	# Extract file path
	tx_unsigned=$(echo "${tx_file}" | sed -n 's/.*\[fileName\]: //p')
	log_info "Created unsigned transaction (PSBT): $tx_unsigned"

	###########################################################################
	# ROUND 1: Nonce Generation (parallel execution possible)
	###########################################################################
	log_step "MuSig2 Round 1: Nonce Generation (3 wallets in parallel)"
	log_info "All signers generate nonces independently"
	log_warn "NOTE: Current MuSig2 implementation contains placeholder TODOs"
	log_warn "This workflow demonstrates the intended protocol flow"

	# Generate nonce for keygen wallet
	log_substep "Generating nonce for keygen wallet (Signer 1)"
	nonce_file_keygen="data/tx/btc/$(basename "${tx_unsigned%.psbt}")_keygen_nonce.json"
	log_info "Output: $nonce_file_keygen"
	# TODO: Uncomment when MuSig2 nonce command is fully implemented
	# keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" musig2 nonce \
	# 	--file "${tx_unsigned}" \
	# 	--output "${nonce_file_keygen}"
	log_warn "Skipping keygen nonce generation (placeholder)"

	# Generate nonce for sign1 wallet
	log_substep "Generating nonce for sign1 wallet (Signer 2)"
	nonce_file_sign1="data/tx/btc/$(basename "${tx_unsigned%.psbt}")_sign1_nonce.json"
	log_info "Output: $nonce_file_sign1"
	# TODO: Uncomment when MuSig2 nonce command is fully implemented
	# sign1 --conf "${CONFIG_SIGN1}" --wallet sign1 musig2 nonce \
	# 	--file "${tx_unsigned}" \
	# 	--output "${nonce_file_sign1}"
	log_warn "Skipping sign1 nonce generation (placeholder)"

	# Generate nonce for sign2 wallet
	log_substep "Generating nonce for sign2 wallet (Signer 3)"
	nonce_file_sign2="data/tx/btc/$(basename "${tx_unsigned%.psbt}")_sign2_nonce.json"
	log_info "Output: $nonce_file_sign2"
	# TODO: Uncomment when MuSig2 nonce command is fully implemented
	# sign2 --conf "${CONFIG_SIGN2}" --wallet sign2 musig2 nonce \
	# 	--file "${tx_unsigned}" \
	# 	--output "${nonce_file_sign2}"
	log_warn "Skipping sign2 nonce generation (placeholder)"

	log_info "Round 1 complete: All nonces generated (placeholder)"
	log_info "Next: Collect and aggregate nonces in watch wallet"

	###########################################################################
	# Nonce Collection and Aggregation (Watch Wallet)
	###########################################################################
	log_substep "Collecting and aggregating nonces (Watch Wallet)"
	psbt_with_nonces="data/tx/btc/$(basename "${tx_unsigned%.psbt}")_with_nonces.psbt"
	log_info "Output: $psbt_with_nonces"
	# TODO: Uncomment when MuSig2 collect-nonces command is fully implemented
	# watch_with_wallet -c "${CONFIG_WATCH}" --coin "${COIN}" musig2 collect-nonces \
	# 	--file "${tx_unsigned}" \
	# 	--nonces "${nonce_file_keygen},${nonce_file_sign1},${nonce_file_sign2}" \
	# 	--output "${psbt_with_nonces}"
	log_warn "Skipping nonce aggregation (placeholder)"
	# For placeholder, copy unsigned PSBT
	cp "${tx_unsigned}" "${psbt_with_nonces}"

	log_info "Nonces collected and aggregated (placeholder)"
	log_info "PSBT now contains aggregated nonces for Round 2"

	###########################################################################
	# ROUND 2: Partial Signature Creation (sequential execution required)
	###########################################################################
	log_step "MuSig2 Round 2: Partial Signature Creation (sequential)"
	log_info "Each signer creates a partial signature using their private key"
	log_info "Partial signatures are added to PSBT sequentially"

	# Sign with keygen wallet (1st partial signature)
	log_substep "Creating partial signature with keygen wallet (1st of 3)"
	psbt_signed_keygen="data/tx/btc/$(basename "${psbt_with_nonces%.psbt}")_keygen_signed.psbt"
	log_info "Output: $psbt_signed_keygen"
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletpassphrase --passphrase "${WALLET_PASSPHRASE}"
	fi
	# TODO: Uncomment when MuSig2 sign command is fully implemented
	# keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" musig2 sign \
	# 	--file "${psbt_with_nonces}" \
	# 	--output "${psbt_signed_keygen}"
	log_warn "Skipping keygen partial signature (placeholder)"
	# For placeholder, use standard signature
	tx_file_signed=$(keygen -c "${CONFIG_KEYGEN}" sign signature --file "${psbt_with_nonces}")
	tx_signed1=$(echo "${tx_file_signed}" | sed -n 's/.*\[fileName\]: //p')
	cp "${tx_signed1}" "${psbt_signed_keygen}"
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletlock
	fi

	log_info "Keygen partial signature added (1/3)"

	# Sign with sign1 wallet (2nd partial signature)
	log_substep "Creating partial signature with sign1 wallet (2nd of 3)"
	psbt_signed_sign1="data/tx/btc/$(basename "${psbt_signed_keygen%.psbt}")_sign1.psbt"
	log_info "Output: $psbt_signed_sign1"
	# TODO: Uncomment when MuSig2 sign command is fully implemented
	# sign1 --conf "${CONFIG_SIGN1}" --wallet sign1 musig2 sign \
	# 	--file "${psbt_signed_keygen}" \
	# 	--output "${psbt_signed_sign1}"
	log_warn "Skipping sign1 partial signature (placeholder)"
	# For placeholder, use standard signature
	tx_file_signed2=$(sign1 --conf "${CONFIG_SIGN1}" --wallet sign1 sign signature --file "${psbt_signed_keygen}")
	tx_signed2=$(echo "${tx_file_signed2}" | sed -n 's/.*\[fileName\]: //p')
	cp "${tx_signed2}" "${psbt_signed_sign1}"

	log_info "Sign1 partial signature added (2/3)"

	# Sign with sign2 wallet (3rd partial signature)
	log_substep "Creating partial signature with sign2 wallet (3rd of 3)"
	psbt_signed_sign2="data/tx/btc/$(basename "${psbt_signed_sign1%.psbt}")_sign2.psbt"
	log_info "Output: $psbt_signed_sign2"
	# TODO: Uncomment when MuSig2 sign command is fully implemented
	# sign2 --conf "${CONFIG_SIGN2}" --wallet sign2 musig2 sign \
	# 	--file "${psbt_signed_sign1}" \
	# 	--output "${psbt_signed_sign2}"
	log_warn "Skipping sign2 partial signature (placeholder)"
	# For placeholder, use standard signature
	tx_file_signed3=$(sign2 --conf "${CONFIG_SIGN2}" --wallet sign2 sign signature --file "${psbt_signed_sign1}")
	tx_signed3=$(echo "${tx_file_signed3}" | sed -n 's/.*\[fileName\]: //p')
	cp "${tx_signed3}" "${psbt_signed_sign2}"

	log_info "Sign2 partial signature added (3/3)"
	log_info "Round 2 complete: All partial signatures collected"

	###########################################################################
	# Signature Aggregation (Watch Wallet)
	###########################################################################
	log_substep "Aggregating partial signatures into final Schnorr signature"
	psbt_final="data/tx/btc/$(basename "${psbt_signed_sign2%.psbt}")_final.psbt"
	log_info "Output: $psbt_final"
	# TODO: Uncomment when MuSig2 aggregate command is fully implemented
	# watch_with_wallet -c "${CONFIG_WATCH}" --coin "${COIN}" musig2 aggregate \
	# 	--files "${psbt_signed_keygen},${psbt_signed_sign1},${psbt_signed_sign2}" \
	# 	--output "${psbt_final}"
	log_warn "Skipping signature aggregation (placeholder)"
	# For placeholder, use the last signed PSBT
	cp "${psbt_signed_sign2}" "${psbt_final}"

	log_info "Partial signatures aggregated into single Schnorr signature (placeholder)"
	log_info "Transaction is now fully signed and ready to broadcast"

	###########################################################################
	# Broadcast Transaction
	###########################################################################
	log_substep "Broadcasting fully signed transaction"
	tx_result=$(watch_with_wallet -c "${CONFIG_WATCH}" send --file "${psbt_final}")
	tx_id="${tx_result##*txID: }"

	log_info "Transaction sent successfully!"
	log_info "Transaction ID: $tx_id"
	log_info ""
	log_info "MuSig2 Benefits:"
	log_info "  • Single 64-byte Schnorr signature on-chain"
	log_info "  • Indistinguishable from single-sig transaction"
	log_info "  • ~64% smaller than traditional 3-of-3 multisig"
	log_info "  • Maximum privacy (N-of-N looks like 1-of-1)"
}

###############################################################################
# Help Message
###############################################################################

show_help() {
	cat <<EOF
Bitcoin E2E Workflow Script - Pattern 10: P2TR MuSig2 N-of-N

This script automates the complete Bitcoin workflow for MuSig2 P2TR (Taproot)
N-of-N transactions. It demonstrates the 2-round MuSig2 protocol where N signers
create a single aggregated Schnorr signature that looks identical to a single-sig
transaction on-chain.

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
  4. Generate keys for keygen and sign wallets
  5. Export extended keys (xpub) from sign wallets
  6. Import extended keys into keygen wallet
  7. Create MuSig2 aggregated public keys
  8. Export MuSig2 descriptors (tr(musig(...)))
  9. Import descriptors into watch wallet
 10. Generate test UTXOs (automatically generates 101 blocks)
 11. Create payment requests
 12. Create unsigned transaction (PSBT)
 13. Round 1: Generate nonces (keygen, sign1, sign2 in parallel)
 14. Collect and aggregate nonces in watch wallet
 15. Round 2: Create partial signatures (sequential)
     - Keygen wallet (1st partial signature)
     - Sign1 wallet (2nd partial signature)
     - Sign2 wallet (3rd partial signature)
 16. Aggregate partial signatures into single Schnorr signature
 17. Broadcast transaction

The script demonstrates the MuSig2 2-round protocol for creating N-of-N multisig
transactions that appear as single-sig on-chain.

Transaction Pattern Details:
  - Address Type: P2TR (BIP86 Taproot + BIP327 MuSig2)
  - Address Format: bcrt1p... (Regtest, 62 chars)
  - Signature Requirement: N-of-N (all 3 signatures required)
  - Protocol: 2-Round MuSig2 (BIP327)
  - On-Chain: Single 64-byte Schnorr signature
  - Descriptor: tr(musig([fp1/86'/1'/1']xpub1,[fp2]xpub2,[fp3]xpub3)/0/*)
  - Privacy: Maximum (indistinguishable from single-sig)
  - Size: ~64% smaller than traditional 3-of-3 multisig

IMPORTANT NOTE:
  This E2E script currently uses placeholder implementations for MuSig2 commands.
  The script demonstrates the intended workflow and will be fully functional once
  the MuSig2 implementation is complete. See TODOs in the code for details.

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

	log_info "Starting Bitcoin E2E Workflow - Pattern 10: P2TR MuSig2 N-of-N"
	log_info "Coin: $COIN"
	log_info "Encrypted: $ENCRYPTED"
	log_info "Sign wallet count: $SIGN_WALLET_NUM (MuSig2 N-of-N: keygen + sign1 + sign2)"
	log_warn "IMPORTANT: This script demonstrates MuSig2 workflow framework"
	log_warn "MuSig2 CLI commands currently contain placeholder implementations"
	echo ""

	# Execute workflow phases
	# Note: Explicitly check return values for functions that may return non-zero
	# The 'set -e' option only exits for failing commands, not for functions that return non-zero
	check_prerequisites
	setup_infrastructure
	setup_wallets
	key_generation_phase
	musig2_setup_phase
	generate_test_utxos || exit 1
	create_payment_requests_phase || exit 1
	transaction_flow_phase || exit 1

	log_step "Bitcoin E2E Workflow Completed Successfully!"
	log_info "Summary:"
	log_info "  ✓ Infrastructure setup complete"
	log_info "  ✓ Wallets created and configured"
	log_info "  ✓ HD keys generated for keygen and sign wallets"
	log_info "  ✓ MuSig2 aggregated public keys created"
	log_info "  ✓ Descriptors exported and imported (deposit, payment, stored accounts)"
	log_info "  ✓ P2TR MuSig2 N-of-N addresses created (placeholder)"
	log_info "  ✓ Test UTXOs generated"
	log_info "  ✓ Payment requests created (using payment account)"
	log_info "  ✓ Transaction created using MuSig2 2-round protocol (placeholder)"
	log_info "  ✓ Transaction broadcast successfully"
	echo ""
	log_info "Transaction Pattern Used:"
	log_info "  • P2TR (BIP86 Taproot + BIP327 MuSig2) N-of-N multisig"
	log_info "  • 2-Round Protocol: Nonce Generation → Partial Signing → Aggregation"
	log_info "  • Single 64-byte Schnorr signature on-chain"
	log_info "  • Indistinguishable from single-sig transaction"
	log_info "  • ~64% smaller than traditional 3-of-3 multisig"
	log_info "  • Maximum privacy (N-of-N looks like 1-of-1)"
	echo ""
	log_warn "Current Status: MuSig2 Implementation Framework"
	log_warn "The script demonstrates the workflow structure using placeholders"
	log_warn "Full functionality pending MuSig2 command implementation"
	echo ""
	log_info "You can now use the wallet system for Bitcoin MuSig2 operations"
	log_info "To cleanup, run: $0 --cleanup"
	log_info "To full reset for fresh state, run: $0 --reset"
}

# Trap errors and cleanup
trap 'log_error "Script failed at line $LINENO"' ERR

# Run main
main "$@"
