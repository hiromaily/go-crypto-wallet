# Bitcoin Operation Scripts

This directory contains scripts for Bitcoin wallet operations.

## E2E Workflow Script

### Overview

The `e2e-workflow.sh` script provides a unified, automated workflow for complete Bitcoin end-to-end testing. It serves as a regression test tool to verify that program modifications don't break the Bitcoin workflow.

### Usage

```bash
# Run complete E2E workflow
./scripts/operation/btc/e2e-workflow.sh

# Or use Makefile target
make btc-e2e-test

# Run with verbose output
make btc-e2e-test-verbose

# Cleanup containers and state
make btc-e2e-cleanup
```

### What the Script Does

The script automates the following workflow:

1. **Prerequisites Check**
   - Verifies Docker and Docker Compose are installed
   - Checks that CLI commands (watch, keygen, sign1, sign2) are available

2. **Infrastructure Setup**
   - Starts database container
   - Starts Bitcoin node containers (btc-watch, btc-keygen, btc-sign1, btc-sign2)
   - Waits for all services to be healthy

3. **Wallet Setup**
   - Creates wallets in Bitcoin nodes (if not already created)
   - Loads wallets automatically

4. **Key Generation Phase**
   - **Keygen wallet**: Creates seed → Generates HD keys → Imports private keys
   - **Sign wallets**: Creates seed → Generates HD keys → Imports private keys → Exports full public keys

5. **Multisig Setup Phase**
   - Imports full public keys into keygen wallet
   - Creates multisig addresses for deposit, payment, and stored accounts
   - Exports addresses from keygen wallet
   - Imports addresses into watch wallet

6. **Transaction Flow Phase** (requires UTXOs)
   - Creates unsigned payment transaction
   - Signs with keygen wallet (1st signature)
   - Signs with sign1 wallet (2nd signature)
   - Signs with sign2 wallet (3rd signature)
   - Sends fully signed transaction
   - Outputs transaction ID

### Options

- `--cleanup`: Stop containers and cleanup state, then exit
- `--verbose`: Enable verbose output (shows all commands)
- `-h, --help`: Display help message

### Prerequisites

1. **Docker and Docker Compose**: Required for running Bitcoin nodes and database
2. **Built CLI commands**: Run `make build` to build the wallet CLI commands
3. **Docker images**: Bitcoin Core 29.2 image will be pulled automatically

### Testing in Regtest Mode

The transaction phase requires UTXOs to be available. For testing in regtest mode:

```bash
# 1. Get a payment address from the exported address CSV files
# Look in: data/address/btc/address_payment_*.csv

# 2. Generate test coins to a payment address
docker exec btc-watch bitcoin-cli -regtest -rpcuser=xyz -rpcpassword=xyz generatetoaddress 101 <payment_address>

# 3. Check balance
watch -coin btc monitor balance

# 4. Re-run the E2E script
make btc-e2e-test
```

### Directory Structure

After running the script, the following directories will contain generated files:

```
data/
├── address/btc/          # Exported address CSV files
├── fullpubkey/btc/       # Exported full public key CSV files
└── tx/btc/               # Transaction files (unsigned and signed)
```

### Use Cases

1. **Regression Testing**: Run after code changes to verify Bitcoin workflow still works
2. **Development Setup**: Quickly set up a complete Bitcoin test environment
3. **CI/CD Integration**: Can be integrated into CI/CD pipelines (returns appropriate exit codes)
4. **Learning**: Demonstrates the complete Bitcoin wallet workflow

### Troubleshooting

**"No utxo" Error**

If you see "No utxo" when creating a transaction:
- The script will skip the transaction flow phase
- You need to generate test coins (see "Testing in Regtest Mode" above)
- Re-run the script after generating coins

**Container Health Check Failures**

If containers don't become healthy:
- Check Docker is running
- Check container logs: `docker compose -f compose.btc.yaml logs`
- Try cleanup and restart: `make btc-e2e-cleanup && make btc-e2e-test`

**CLI Commands Not Found**

If CLI commands are not available:
- Build the project: `make build`
- Ensure binaries are in PATH or install them: `make install`

### Cleanup

To stop containers and cleanup state:

```bash
make btc-e2e-cleanup
```

This will:
- Stop and remove all Bitcoin node containers
- Stop and remove database container
- Remove Docker volumes

Note: Wallet data directories (`docker/nodes/btc/`) are not automatically deleted for safety. To completely reset, manually delete these directories.

## Other Scripts

### generate-btc-key.sh

Generates keys for Bitcoin wallets (legacy script, now superseded by e2e-workflow.sh for testing).

### create-btc-tx-payment.sh

Creates and signs a payment transaction (legacy script, now superseded by e2e-workflow.sh for testing).

### create-bitcoind-wallet.sh

Creates wallets in Bitcoin nodes (now integrated into e2e-workflow.sh).

### load-bitcoind-wallet.sh

Loads existing wallets in Bitcoin nodes (now integrated into e2e-workflow.sh).
