# Bitcoin E2E Workflow

This document describes the complete Bitcoin end-to-end workflow implemented in `e2e-workflow.sh`.

## Overview

The E2E workflow automates the complete Bitcoin wallet operation from infrastructure setup to transaction execution. It serves as both a regression test tool and documentation of the standard Bitcoin operation flow.

## Prerequisites

Before running the workflow, ensure the following are available:

- **Docker & Docker Compose**: For running Bitcoin nodes and database
- **CLI Commands**: `watch`, `keygen`, `sign1`, `sign2` (build with `make build`)

## Workflow Phases

### Phase 1: Prerequisites Check

Verifies that all required tools and dependencies are available:

- Docker and Docker Compose are installed and running
- All CLI commands (`watch`, `keygen`, `sign1`, `sign2`) are built and accessible

### Phase 2: Infrastructure Setup

Starts the required infrastructure containers:

```
┌─────────────────┐     ┌─────────────────┐
│   wallet-db     │     │   btc-watch     │
│   (MySQL)       │     │   (Bitcoin Node)│
└─────────────────┘     └─────────────────┘
                        ┌─────────────────┐
                        │   btc-keygen    │
                        │   (Bitcoin Node)│
                        └─────────────────┘
                        ┌─────────────────┐
                        │   btc-sign1     │
                        │   (Bitcoin Node)│
                        └─────────────────┘
                        ┌─────────────────┐
                        │   btc-sign2     │
                        │   (Bitcoin Node)│
                        └─────────────────┘
```

1. Start database container (`compose.yaml`)
2. Wait for database to be healthy
3. Start Bitcoin node containers (`compose.btc.yaml`)
4. Wait for all Bitcoin nodes to be healthy

### Phase 3: Wallet Setup

Creates wallets in each Bitcoin node:

| Node       | Wallet Name |
|------------|-------------|
| btc-watch  | watch       |
| btc-keygen | keygen      |
| btc-sign1  | sign1       |
| btc-sign2  | sign2       |

### Phase 4: Key Generation

This phase generates all required keys for the wallet system.

#### 4.1 Keygen Wallet Operations

```
┌─────────────────────────────────────────────────────────────┐
│                    KEYGEN WALLET                            │
├─────────────────────────────────────────────────────────────┤
│ 1. Create seed                                              │
│    └─ keygen create seed                                    │
│                                                             │
│ 2. Create HD keys for each account (10 keys each)          │
│    ├─ client account                                        │
│    ├─ deposit account                                       │
│    ├─ payment account                                       │
│    └─ stored account                                        │
│                                                             │
│ 3. Import private keys into Bitcoin node                   │
│    └─ For all accounts (client, deposit, payment, stored)  │
└─────────────────────────────────────────────────────────────┘
```

#### 4.2 Sign Wallet Operations

```
┌─────────────────────────────────────────────────────────────┐
│                    SIGN WALLETS (sign1, sign2)              │
├─────────────────────────────────────────────────────────────┤
│ 1. Create seed                                              │
│    └─ sign1 create seed                                     │
│                                                             │
│ 2. Create HD keys                                           │
│    ├─ sign1 create hdkey                                    │
│    └─ sign2 create hdkey                                    │
│                                                             │
│ 3. Import private keys into Bitcoin nodes                  │
│    ├─ sign1 import privkey                                  │
│    └─ sign2 import privkey                                  │
│                                                             │
│ 4. Export full public keys                                 │
│    ├─ sign1 export fullpubkey → fullpubkey_auth1.csv       │
│    └─ sign2 export fullpubkey → fullpubkey_auth2.csv       │
└─────────────────────────────────────────────────────────────┘
```

### Phase 5: Multisig Setup

Creates multisig addresses and exports them to the watch wallet.

```
┌─────────────────────────────────────────────────────────────┐
│                    MULTISIG SETUP                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  KEYGEN WALLET                                              │
│  ├─ Import fullpubkey from sign1                           │
│  ├─ Import fullpubkey from sign2                           │
│  │                                                          │
│  ├─ Create multisig addresses                              │
│  │   ├─ deposit account                                     │
│  │   ├─ payment account                                     │
│  │   └─ stored account                                      │
│  │                                                          │
│  └─ Export addresses                                       │
│      ├─ client addresses  → address_client.csv             │
│      ├─ deposit addresses → address_deposit.csv            │
│      ├─ payment addresses → address_payment.csv            │
│      └─ stored addresses  → address_stored.csv             │
│                                                             │
│  WATCH WALLET                                               │
│  └─ Import all address files                               │
│      ├─ client addresses                                    │
│      ├─ deposit addresses                                   │
│      ├─ payment addresses                                   │
│      └─ stored addresses                                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Phase 6: Test UTXO Generation

For regtest environment, generates test UTXOs for transaction testing:

1. Extract payment address from exported address file
2. Generate 101 blocks to the payment address (creates mature coinbase)
3. Wait for balance update in watch wallet

### Phase 7: Transaction Flow

Creates, signs, and sends a payment transaction.

```
┌─────────────────────────────────────────────────────────────┐
│                    TRANSACTION FLOW                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. CREATE UNSIGNED TRANSACTION                             │
│     └─ watch create payment → tx_unsigned.hex              │
│                                                             │
│  2. SIGN WITH KEYGEN WALLET (1st signature)                │
│     └─ keygen sign --file tx_unsigned.hex                  │
│        → tx_signed_1.hex                                    │
│                                                             │
│  3. SIGN WITH SIGN1 WALLET (2nd signature)                 │
│     └─ sign1 sign --file tx_signed_1.hex                   │
│        → tx_signed_2.hex                                    │
│                                                             │
│  4. SIGN WITH SIGN2 WALLET (3rd signature)                 │
│     └─ sign2 sign --file tx_signed_2.hex                   │
│        → tx_signed_3.hex                                    │
│                                                             │
│  5. SEND TRANSACTION                                        │
│     └─ watch send --file tx_signed_3.hex                   │
│        → Transaction ID                                     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Usage

### Run Complete E2E Workflow

```bash
# From fresh state (full reset first)
./scripts/operation/btc/e2e-workflow.sh --reset

# Run workflow (assumes clean state or continuing from previous run)
./scripts/operation/btc/e2e-workflow.sh

# Run with verbose output
./scripts/operation/btc/e2e-workflow.sh --verbose
```

### Cleanup

```bash
# Stop containers and cleanup state
./scripts/operation/btc/e2e-workflow.sh --cleanup
```

### Options

| Option              | Description                                        |
|---------------------|----------------------------------------------------|
| `--reset`           | Full reset: cleanup all state for fresh start      |
| `--cleanup`         | Stop containers and cleanup state, then exit       |
| `--verbose`         | Enable verbose output                              |
| `--non-interactive` | Run without interactive prompts (for CI/CD)        |
| `-h, --help`        | Display help message                               |

## Environment Variables

| Variable            | Default | Description                                      |
|---------------------|---------|--------------------------------------------------|
| `RPC_USER`          | `xyz`   | Bitcoin RPC username (for regtest)               |
| `RPC_PASSWORD`      | `xyz`   | Bitcoin RPC password (for regtest)               |
| `WALLET_PASSPHRASE` | `test`  | Wallet passphrase for encrypted wallets          |

## Configuration Files

| File                              | Purpose                    |
|-----------------------------------|----------------------------|
| `config/wallet/btc_watch.toml`    | Watch wallet configuration |
| `config/wallet/btc_keygen.toml`   | Keygen wallet configuration|
| `config/wallet/btc_sign1.toml`    | Sign1 wallet configuration |
| `config/wallet/btc_sign2.toml`    | Sign2 wallet configuration |

## Generated Files

### Address Files (`data/address/btc/`)

- `address_client_*.csv` - Client addresses (non-multisig)
- `address_deposit_*.csv` - Deposit multisig addresses
- `address_payment_*.csv` - Payment multisig addresses
- `address_stored_*.csv` - Stored multisig addresses

### Public Key Files (`data/fullpubkey/btc/`)

- `fullpubkey_auth1_*.csv` - Full public keys from sign1
- `fullpubkey_auth2_*.csv` - Full public keys from sign2

### Transaction Files (`data/tx/btc/`)

- `tx_*.hex` - Unsigned and signed transaction files

## Account Types

| Account  | Purpose                                           | Multisig |
|----------|---------------------------------------------------|----------|
| client   | Client-facing addresses for receiving funds       | No       |
| deposit  | Deposit addresses for initial fund receipt        | Yes      |
| payment  | Payment addresses for outgoing transactions       | Yes      |
| stored   | Cold storage addresses for long-term holding      | Yes      |

## Signature Flow (3-of-3 Multisig)

The transaction requires signatures from three wallets:

1. **Keygen Wallet** - Primary key holder
2. **Sign1 Wallet** - First authorization signer
3. **Sign2 Wallet** - Second authorization signer

This implements a 3-of-3 multisig scheme where all three parties must sign to authorize a transaction.

## Troubleshooting

### "No UTXOs available"

If the transaction phase reports no UTXOs:

1. Verify blocks were generated successfully
2. Check balance with: `watch -c config/wallet/btc_watch.toml --coin btc monitor balance`
3. The script automatically generates 101 blocks; if this fails, manually run:

   ```bash
   docker exec btc-watch bitcoin-cli -regtest generatetoaddress 101 <address>
   ```

### Container Health Issues

If containers fail health checks:

1. Check container logs: `docker logs btc-watch`
2. Verify Docker resources are available
3. Try full reset: `./scripts/operation/btc/e2e-workflow.sh --reset`

### Key Import Failures

If key import fails:

1. Ensure wallets were created successfully in Bitcoin nodes
2. Check if encrypted mode is properly configured
3. Verify configuration file paths are correct
