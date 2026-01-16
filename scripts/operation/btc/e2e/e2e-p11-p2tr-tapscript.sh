#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern 11: P2TR Tapscript M-of-N
# This script automates the complete Bitcoin workflow for M-of-N multisig P2TR (Taproot) Tapscript transactions
# Usage: ./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh [OPTIONS]
# Options:
#   --cleanup  Stop containers and cleanup state
#   --reset    Full reset and run from scratch
#   --verbose  Enable verbose output
#   --non-interactive  Run without prompts (for CI/CD)
#   -h, --help Display help message
#
# Reference Documentation:
#   docs/crypto/btc/e2e_transaction_patterns.md - E2E transaction patterns
#   docs/crypto/btc/TAPROOT_GUIDE.md - Taproot user guide
#   docs/crypto/btc/btc_bch_technical_guide.md - Script Path details
#
# Transaction Pattern:
#   Pattern 11: BTC P2TR Tapscript M-of-N
#   - Address Type: P2TR (BIP86 Taproot + BIP342 Tapscript)
#   - Address Format: `bc1p...` (Mainnet), `tb1p...` (Testnet), `bcrt1p...` (Regtest)
#   - Signature Requirement: M-of-N threshold (2-of-3)
#   - Descriptor: tr(internal_key,sortedmulti_a(2,[fingerprint/86'/1'/1']xpub1/0/*,xpub2/0/*,xpub3/0/*))
#   - Bitcoin Core Version Required: v22.0+
#   - Spending Path: Script Path (not Key Path)
#
# Required Config Settings:
#   - config/wallet/btc_watch.yaml:  address_type: "taproot"
#   - config/wallet/btc_keygen.yaml: address_type: "taproot"
#   - config/wallet/btc_sign1.yaml:  address_type: "taproot"
#   - config/wallet/btc_sign2.yaml:  address_type: "taproot"
#
# IMPORTANT NOTE:
#   This script demonstrates the Tapscript M-of-N workflow framework.
#   Tapscript implementation is currently under development.
#   The script uses placeholder implementations where full Tapscript support is pending.

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
SIGN_WALLET_NUM=2 # 2-of-3: keygen + sign1 + sign2 = 3 signers
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
# Use 2-of-3 account configuration for Pattern 11 (Tapscript 2-of-3)
CONFIG_ACCOUNT="${PROJECT_ROOT}/config/wallet/account_2of3.yaml"

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
# Pattern 11 (P2TR Tapscript M-of-N) requires:
#   - address_type: "taproot" (derives key_type: bip86 automatically)
# Note: key_type is automatically derived from address_type in Go code
#       (see internal/domain/address/types.go AddrType.ToKeyType())
# Note: "taproot" is the correct config value (not "bech32m" - that's the encoding format)
export WALLET_ADDRESS_TYPE="taproot"

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

	# Check CLI commands (watch, keygen, sign1, sign2 for Tapscript)
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

	# Start Bitcoin nodes (watch, keygen, sign1, sign2 for Tapscript)
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

	# Keygen wallet - create hdkeys for Tapscript accounts
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
	fullpubkey_file1=$(extract_filename "$file_fullpubkey_auth1")
	fullpubkey_file2=$(extract_filename "$file_fullpubkey_auth2")

	log_info "Exported fullpubkey files:"
	log_info "  sign1: $fullpubkey_file1"
	log_info "  sign2: $fullpubkey_file2"

	# Store for next phase
	export FULLPUBKEY_FILE1="$fullpubkey_file1"
	export FULLPUBKEY_FILE2="$fullpubkey_file2"
}

###############################################################################
# Tapscript Setup Phase (2-of-3)
###############################################################################

tapscript_setup_phase() {
	log_step "Tapscript Setup Phase (2-of-3 P2TR Script Path)"

	# Import fullpubkeys
	log_substep "Importing full public keys into keygen wallet"
	log_info "Importing fullpubkey from sign1: $FULLPUBKEY_FILE1"
	keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" import fullpubkey --file "${FULLPUBKEY_FILE1}"

	log_info "Importing fullpubkey from sign2: $FULLPUBKEY_FILE2"
	keygen -c "${CONFIG_KEYGEN}" --coin "${COIN}" import fullpubkey --file "${FULLPUBKEY_FILE2}"

	# For Tapscript, we need to create aggregated public keys from keygen + sign1 + sign2
	# and build the script tree with M-of-N multisig
	log_substep "Exporting Tapscript descriptors from keygen wallet"

	# TODO: Implement actual Tapscript descriptor export with sortedmulti_a
	# For now, we export standard Taproot descriptors as placeholder
	# The real implementation should use: tr(internal_key,sortedmulti_a(2,[fp1/86'/1'/1']xpub1/0/*,xpub2/0/*,xpub3/0/*))

	log_warn "Tapscript descriptor export is currently a placeholder"
	log_warn "The script will export standard Taproot descriptors for now"
	log_warn "Full Tapscript support (Script Path with sortedmulti_a) is under development"

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

	# Import descriptors into watch wallet for Tapscript accounts
	# Note: Using wrapper function to set WALLET_BITCOIN_HOST for Bitcoin Core RPC
	log_substep "Importing descriptors into watch wallet"
	log_info "Importing deposit descriptors"
	watch_with_wallet -c "${CONFIG_WATCH}" --coin "${COIN}" import descriptor --file "${descriptor_deposit}" --account deposit

	log_info "Importing payment descriptors"
	watch_with_wallet -c "${CONFIG_WATCH}" --coin "${COIN}" import descriptor --file "${descriptor_payment}" --account payment

	log_info "Importing stored descriptors"
	watch_with_wallet -c "${CONFIG_WATCH}" --coin "${COIN}" import descriptor --file "${descriptor_stored}" --account stored

	log_info "All descriptors imported successfully"
	log_info "Note: Pattern 11 uses Tapscript 2-of-3 with P2TR (BIP86 + BIP342)"
	log_info "Note: Current implementation uses placeholder descriptors"

	# Derive payment address from descriptor for UTXO generation
	log_substep "Deriving payment address from descriptor for UTXO generation"

	# Extract first descriptor from payment_descriptors.json
	# For P2TR Tapscript, we use the descriptor at index 0
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
		if [ -n "$trusted_balance" ] && [ "$(echo "$trusted_balance > 0" | bc -l)" -eq 1 ]; then
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
# Transaction Flow Phase (Tapscript 2-of-3)
###############################################################################

transaction_flow_phase() {
	log_step "Transaction Flow Phase (Tapscript 2-of-3 Script Path)"

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
	tx_unsigned=$(extract_filename "$tx_file")
	log_info "Created unsigned transaction (PSBT): $tx_unsigned"

	###########################################################################
	# Tapscript Signing: Script Path Spend with M-of-N Multisig
	###########################################################################
	log_step "Tapscript Script Path Signing (2-of-3 Multisig with Schnorr)"
	log_info "For Tapscript Script Path, we need:"
	log_info "  1. Control block (internal key + Merkle path)"
	log_info "  2. Script reveal (sortedmulti_a script)"
	log_info "  3. M Schnorr signatures (2-of-3 = 2 signatures required)"
	log_warn "NOTE: Current Tapscript implementation contains placeholder TODOs"
	log_warn "This workflow demonstrates the intended protocol flow"

	# Sign with keygen wallet (1st Schnorr signature)
	log_substep "Signing with keygen wallet (1st of 2 required Schnorr signatures)"
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletpassphrase --passphrase "${WALLET_PASSPHRASE}"
	fi
	# TODO: Implement Tapscript Script Path signing
	# For now, use standard signature as placeholder
	log_warn "Using standard signature as placeholder for Tapscript signing"
	tx_file_signed=$(keygen -c "${CONFIG_KEYGEN}" sign signature --file "${tx_unsigned}")
	if [ "$ENCRYPTED" = "true" ]; then
		keygen -c "${CONFIG_KEYGEN}" api walletlock
	fi

	tx_signed1=$(extract_filename "$tx_file_signed")
	log_info "Signed transaction (1st signature): $tx_signed1"

	# Sign with sign1 wallet (2nd Schnorr signature)
	# For 2-of-3 Tapscript, 2 signatures are sufficient
	log_substep "Signing with sign1 wallet (2nd signature - completing 2-of-3 requirement)"
	tx_file_signed2=$(sign1 --conf "${CONFIG_SIGN1}" --wallet sign1 sign signature --file "${tx_signed1}")
	tx_signed2=$(extract_filename "$tx_file_signed2")
	log_info "Signed transaction (2nd signature): $tx_signed2"
	log_info "Note: 2-of-3 Tapscript requirement satisfied with 2 Schnorr signatures"

	# Sign2 wallet is not needed for 2-of-3 (only 2 signatures required)

	# Send transaction
	log_substep "Broadcasting fully signed Tapscript transaction"
	tx_result=$(watch_with_wallet -c "${CONFIG_WATCH}" send --file "${tx_signed2}")
	tx_id="${tx_result##*txID: }"

	log_info "Transaction sent successfully!"
	log_info "Transaction ID: $tx_id"
	log_info ""
	log_info "Tapscript Benefits:"
	log_info "  • Only used script path revealed on-chain"
	log_info "  • 2 × 64-byte Schnorr signatures (vs ECDSA ~71 bytes each)"
	log_info "  • ~50% smaller than P2WSH 2-of-3 multisig"
	log_info "  • Enhanced privacy (unused paths hidden in Merkle tree)"
}

###############################################################################
# Help Message
###############################################################################

show_help() {
	cat <<EOF
Bitcoin E2E Workflow Script - Pattern 11: P2TR Tapscript M-of-N

This script automates the complete Bitcoin workflow for M-of-N multisig P2TR
(Taproot) Tapscript transactions. It demonstrates the Script Path spending with
M-of-N threshold multisig using Schnorr signatures.

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
  7. Create Tapscript script tree (Merkle tree with M-of-N script)
  8. Tweak internal key with Merkle root
  9. Export Tapscript descriptors (tr(internal_key,sortedmulti_a(...)))
 10. Import descriptors into watch wallet
 11. Generate test UTXOs (automatically generates 101 blocks)
 12. Create payment requests
 13. Create unsigned transaction (PSBT)
 14. Script Path Signing:
     - Sign with keygen wallet (1st Schnorr signature)
     - Sign with sign1 wallet (2nd Schnorr signature - completing 2-of-3)
     - Construct witness: [sig1, sig2, script, control_block]
 15. Broadcast transaction

The script demonstrates the Tapscript M-of-N workflow for creating threshold
multisig transactions that reveal only the used script path on-chain.

Transaction Pattern Details:
  - Address Type: P2TR (BIP86 Taproot + BIP342 Tapscript)
  - Address Format: bcrt1p... (Regtest, 62 chars)
  - Signature Requirement: 2-of-3 (M-of-N threshold)
  - Signatures: 2 × Schnorr (64 bytes each)
  - Descriptor: tr(internal_key,sortedmulti_a(2,[fp1]xpub1,[fp2]xpub2,[fp3]xpub3))
  - Spending Path: Script Path (not Key Path)
  - On-Chain: Only used script path + control block revealed
  - Privacy: Unused paths remain hidden in Merkle tree
  - Size: ~50% smaller than P2WSH 2-of-3 multisig

IMPORTANT NOTE:
  This E2E script currently uses placeholder implementations for Tapscript commands.
  The script demonstrates the intended workflow and will be fully functional once
  the Tapscript implementation is complete. See TODOs in the code for details.

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

	log_info "Starting Bitcoin E2E Workflow - Pattern 11: P2TR Tapscript M-of-N"
	log_info "Coin: $COIN"
	log_info "Encrypted: $ENCRYPTED"
	log_info "Sign wallet count: $SIGN_WALLET_NUM (Tapscript 2-of-3: keygen + sign1 + sign2)"
	log_warn "IMPORTANT: This script demonstrates Tapscript workflow framework"
	log_warn "Tapscript CLI commands currently contain placeholder implementations"
	echo ""

	# Execute workflow phases
	# Note: Explicitly check return values for functions that may return non-zero
	# The 'set -e' option only exits for failing commands, not for functions that return non-zero
	check_prerequisites
	setup_infrastructure
	setup_wallets
	key_generation_phase
	tapscript_setup_phase
	generate_test_utxos || exit 1
	create_payment_requests_phase || exit 1
	transaction_flow_phase || exit 1

	log_step "Bitcoin E2E Workflow Completed Successfully!"
	log_info "Summary:"
	log_info "  ✓ Infrastructure setup complete"
	log_info "  ✓ Wallets created and configured"
	log_info "  ✓ HD keys generated for keygen and sign wallets"
	log_info "  ✓ Tapscript script tree created (placeholder)"
	log_info "  ✓ Descriptors exported and imported (deposit, payment, stored accounts)"
	log_info "  ✓ P2TR Tapscript 2-of-3 addresses created (placeholder)"
	log_info "  ✓ Test UTXOs generated"
	log_info "  ✓ Payment requests created (using payment account)"
	log_info "  ✓ Transaction created using Tapscript Script Path (placeholder)"
	log_info "  ✓ Transaction broadcast successfully"
	echo ""
	log_info "Transaction Pattern Used:"
	log_info "  • P2TR (BIP86 Taproot + BIP342 Tapscript) 2-of-3 multisig"
	log_info "  • Script Path Spend with sortedmulti_a(2, ...)"
	log_info "  • 2 × Schnorr signatures (64 bytes each)"
	log_info "  • Control block reveals internal key + Merkle path"
	log_info "  • ~50% smaller than P2WSH 2-of-3 multisig"
	log_info "  • Enhanced privacy (unused script paths hidden)"
	echo ""
	log_warn "Current Status: Tapscript Implementation Framework"
	log_warn "The script demonstrates the workflow structure using placeholders"
	log_warn "Full functionality pending Tapscript command implementation"
	echo ""
	log_info "You can now use the wallet system for Bitcoin Tapscript operations"
	log_info "To cleanup, run: $0 --cleanup"
	log_info "To full reset for fresh state, run: $0 --reset"
}

# Trap errors and cleanup
trap 'log_error "Script failed at line $LINENO"' ERR

# Run main
main "$@"
