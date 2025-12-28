# Bitcoin Node Setup Script

This directory contains the automated setup script for Bitcoin nodes used by go-crypto-wallet.

## Overview

The `setup_btc.sh` script initializes Bitcoin nodes when Docker containers start. It configures the Bitcoin nodes (btc-watch, btc-keygen, btc-sign1, btc-sign2) to a state where they can be operated by the go-crypto-wallet CLI.

## Features

- Waits for all Bitcoin nodes to be ready and responsive
- Creates necessary wallets for each node type:
  - `btc-watch` container: `watch` wallet
  - `btc-keygen` container: `keygen` wallet
  - `btc-sign1` container: `sign1` wallet
  - `btc-sign2` container: `sign2` wallet
- Generates initial blocks (101 blocks) in regtest mode for testing
- Handles errors gracefully and is idempotent (safe to run multiple times)
- Provides colored logging output for easy monitoring

## Usage

### Manual Execution

After starting the Bitcoin containers:

```bash
# Start Bitcoin nodes
docker compose -f compose.btc.yaml up -d

# Run the setup script
./scripts/setup/btc/setup_btc.sh
```

### First-Time Setup

For a clean setup from scratch:

```bash
# Clean up existing data (optional)
./scripts/operation/reset-data.sh

# Start containers
docker compose -f compose.btc.yaml up -d

# Run setup script
./scripts/setup/btc/setup_btc.sh
```

### Verify Setup

After running the setup script, verify the wallets are created:

```bash
# List wallets on btc-watch
docker compose -f compose.btc.yaml exec btc-watch bitcoin-cli -regtest -rpcuser=xyz -rpcpassword=xyz listwallets

# Check blockchain info
docker compose -f compose.btc.yaml exec btc-watch bitcoin-cli -regtest -rpcuser=xyz -rpcpassword=xyz getblockchaininfo
```

## Configuration

The script uses the following default configuration:

- **RPC User**: `xyz` (default, can be overridden with `BTC_RPC_USER` environment variable)
- **RPC Password**: `xyz` (default, can be overridden with `BTC_RPC_PASSWORD` environment variable)
- **Max Retries**: 30 attempts
- **Retry Interval**: 2 seconds
- **Initial Blocks**: 101 blocks (for regtest mode)

### Using Environment Variables

You can override the default RPC credentials using environment variables:

```bash
# Set custom credentials
export BTC_RPC_USER="myuser"
export BTC_RPC_PASSWORD="mypassword"

# Run setup script with custom credentials
./scripts/setup/btc/setup_btc.sh
```

This is more secure than hardcoding credentials and allows for different configurations in different environments.

## Error Handling

The script handles common errors gracefully:

- **Wallet already exists**: Skips wallet creation and continues
- **Node not ready**: Retries up to 30 times with 2-second intervals
- **Block generation failure**: Logs warning but continues (optional feature)

## Troubleshooting

### Nodes not becoming ready

If nodes fail to become ready after 30 attempts:

1. Check if containers are running: `docker compose -f compose.btc.yaml ps`
2. Check container logs: `docker compose -f compose.btc.yaml logs btc-watch`
3. Verify RPC credentials in `docker/nodes/btc/data*/bitcoin.conf`

### Wallet creation fails

If wallet creation fails unexpectedly:

1. Check if wallet already exists: `docker compose -f compose.btc.yaml exec btc-watch bitcoin-cli -regtest -rpcuser=xyz -rpcpassword=xyz listwallets`
2. Try creating wallet manually to see detailed error
3. Check container logs for Bitcoin Core errors

### Block generation fails

Block generation is optional and failure doesn't stop the setup. If you need blocks:

1. Verify the watch wallet is created and loaded
2. Generate blocks manually: `docker compose -f compose.btc.yaml exec btc-watch bitcoin-cli -regtest -rpcuser=xyz -rpcpassword=xyz -rpcwallet=watch -generate 101`

## Architecture Notes

### Multi-Sign Wallet Support

This setup supports multiple sign wallets (btc-sign1 and btc-sign2) for multisig operations:

- **btc-sign1**: Uses `sign1` directory, RPC port 20332
- **btc-sign2**: Uses `sign2` directory, RPC port 21332

Both sign nodes run in offline mode (`-maxconnections=0`) for security.

### Wallet Types

- **Watch Wallet** (btc-watch): Online node with network connectivity, used for monitoring and broadcasting transactions
- **Keygen Wallet** (btc-keygen): Offline node for key generation, provides first signature for multisig
- **Sign Wallets** (btc-sign1, btc-sign2): Offline nodes for signing, provide additional signatures for multisig

## Integration

### Automatic Execution

To run the setup script automatically when containers start, you can:

1. Add it to your Docker Compose startup sequence
2. Use a post-start hook
3. Include it in your CI/CD pipeline
4. Run it as part of your deployment scripts

Example using Docker Compose command override:

```yaml
# Not implemented yet - for future consideration
btc-watch:
  command: >
    sh -c "bitcoind -printtoconsole & 
           sleep 10 && 
           /scripts/setup_btc.sh && 
           wait"
```

## Related Files

- `compose.btc.yaml`: Docker Compose configuration for Bitcoin nodes
- `scripts/operation/btc/create-bitcoind-wallet.sh`: Manual wallet creation script (legacy)
- `scripts/operation/btc/load-bitcoind-wallet.sh`: Wallet loading script
- `scripts/operation/reset-data.sh`: Data cleanup script
- `docker/nodes/btc/data*/bitcoin.conf`: Bitcoin Core configuration files
