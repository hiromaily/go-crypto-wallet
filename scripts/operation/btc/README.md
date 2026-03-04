# Bitcoin Operation Scripts

This directory contains Bitcoin operation and E2E workflow scripts.

## Directory Structure

```
btc/
├── btc_common.sh         # BTC-specific common utilities (NEW)
├── README.md             # This file
├── create-bitcoind-wallet.sh
├── create-btc-tx-deposit.sh
├── create-btc-tx-payment.sh
├── create-btc-tx-transfer-all.sh
├── create-btc-tx-transfer.sh
├── generate-btc-key.sh
├── load-bitcoind-wallet.sh
└── e2e/                  # E2E test scripts
    ├── README.md
    ├── e2e-p1-p2pkh-singlesig.sh
    ├── e2e-p2-p2pkh-2of3.sh
    └── ... (other patterns)
```

## Common Utilities

### btc_common.sh

BTC-specific common functions for E2E scripts. This file automatically sources `../common.sh`, so you don't need to source both.

**Usage in E2E scripts:**

```bash
# Source BTC common utilities
source "${SCRIPT_DIR}/../btc_common.sh"

# Initialize config paths
btc_get_config_paths

# Use common functions
btc_check_prerequisites "watch keygen"
btc_setup_infrastructure "btc-watch btc-keygen"
btc_setup_wallets "btc-watch:watch btc-keygen:keygen"
```

**Available Functions:**

| Function | Description |
|----------|-------------|
| `btc_get_config_paths` | Set standard config file paths |
| `btc_clean_data_files` | Clean BTC data directories |
| `btc_clean_wallet_data` | Clean Bitcoin node wallet data |
| `btc_full_reset` | Full reset with volume deletion |
| `btc_cleanup` | Stop containers |
| `btc_check_prerequisites` | Check Docker and CLI commands |
| `btc_setup_infrastructure` | Start database and Bitcoin nodes |
| `btc_setup_wallets` | Create wallets in Bitcoin nodes |
| `btc_watch_cmd` | Wrapper for watch commands |
| `btc_keygen_cmd` | Wrapper for keygen commands |
| `btc_log_no_utxo_error` | Log UTXO error details |
| `btc_wait_for_balance` | Wait for balance update |
| `btc_generate_test_utxos` | Generate test UTXOs |
| `btc_derive_address_from_descriptor` | Derive address from descriptor |
| `btc_extract_file_path` | Extract file path from output |
| `btc_extract_descriptor_path` | Extract descriptor path from output |
| `btc_get_sender_address` | Get sender address from database |
| `btc_generate_receiver_addresses` | Generate receiver addresses |
| `btc_insert_payment_requests` | Insert payment requests |
| `btc_parse_args` | Parse common e2e script arguments |

**Database Abstraction Functions:**

| Function | Description |
|----------|-------------|
| `db_is_sqlite` | Check if using SQLite database |
| `db_is_mysql` | Check if using MySQL database |
| `db_query` | Execute database query (SELECT) |
| `db_execute` | Execute database command (INSERT, UPDATE, DELETE) |
| `sqlite_init_db` | Initialize SQLite database with schema |
| `sqlite_clean_db` | Remove SQLite database file |
| `sqlite_query` | Execute SQLite query directly |
| `mysql_query` | Execute MySQL query via Docker |

### Database Configuration

The scripts support three database backends:

#### SQLite (Default)

Uses local SQLite files. No Docker database container required.

```bash
# Run with SQLite (default)
./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh
# or explicitly:
DB_TYPE=sqlite ./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh
```

Benefits:
- Faster test startup (no Docker database container)
- Parallel test execution (each test uses separate DB files)
- Lighter CI/CD environments

#### PostgreSQL

Uses Docker PostgreSQL container. Set `DB_TYPE=postgres`.

```bash
# Run with PostgreSQL
DB_TYPE=postgres ./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh
```

Start the container with: `docker compose --profile postgres up`

#### MySQL

Uses Docker MySQL container. Set `DB_TYPE=mysql`.

```bash
# Run with MySQL
DB_TYPE=mysql ./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh
```

Start the container with: `docker compose --profile mysql up`

**Environment Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_TYPE` | `sqlite` | Database type: `sqlite`, `postgres`, or `mysql` |
| `SQLITE_DB_DIR` | `./data/sqlite/btc` | SQLite database directory |
| `SQLITE_WATCH_DB_PATH` | `${SQLITE_DB_DIR}/watch.db` | Watch wallet database |
| `SQLITE_KEYGEN_DB_PATH` | `${SQLITE_DB_DIR}/keygen.db` | Keygen wallet database |
| `SQLITE_SIGN_DB_PATH` | `${SQLITE_DB_DIR}/sign.db` | Sign wallet database |
| `SQLITE_SIGN2_DB_PATH` | `${SQLITE_DB_DIR}/sign2.db` | Sign2 wallet database |
| `BTC_POSTGRES_USER` | `postgres` | PostgreSQL username |
| `BTC_POSTGRES_PASSWORD` | `postgres` | PostgreSQL password |

---

# Bitcoin E2E Workflow

This document describes the complete Bitcoin end-to-end workflow implemented in `e2e/e2e-p2sh-p2wsh-3of3.sh`.

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

**With MySQL (default):**

```
┌─────────────────┐     ┌─────────────────┐
│   wallet-mysql     │     │   btc-watch     │
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

**With SQLite (`DB_TYPE=sqlite`):**

```
┌─────────────────┐     ┌─────────────────┐
│   SQLite Files  │     │   btc-watch     │
│   (Local)       │     │   (Bitcoin Node)│
└─────────────────┘     └─────────────────┘
 ├─ watch.db            ┌─────────────────┐
 ├─ keygen.db           │   btc-keygen    │
 └─ sign.db             │   (Bitcoin Node)│
                        └─────────────────┘
                        ... (sign nodes)
```

**MySQL Mode:**

1. Start database container (`compose.yaml`)
2. Wait for database to be healthy
3. Start Bitcoin node containers (`compose.btc.yaml`)
4. Wait for all Bitcoin nodes to be healthy

**SQLite Mode:**

1. Initialize SQLite database files with schema
2. Start Bitcoin node containers (`compose.btc.yaml`)
3. Wait for all Bitcoin nodes to be healthy

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
./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh --reset

# Run workflow (assumes clean state or continuing from previous run)
./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh

# Run with verbose output
./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh --verbose
```

### Cleanup

```bash
# Stop containers and cleanup state
./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh --cleanup
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
| `config/wallet/btc/watch.yaml`    | Watch wallet configuration |
| `config/wallet/btc/keygen.yaml`   | Keygen wallet configuration|
| `config/wallet/btc/sign1.yaml`    | Sign1 wallet configuration |
| `config/wallet/btc/sign2.yaml`    | Sign2 wallet configuration |

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
2. Check balance with: `watch -c config/wallet/btc/watch.yaml --coin btc monitor balance`
3. The script automatically generates 101 blocks; if this fails, manually run:

   ```bash
   docker exec btc-watch bitcoin-cli -regtest generatetoaddress 101 <address>
   ```

### Container Health Issues

If containers fail health checks:

1. Check container logs: `docker logs btc-watch`
2. Verify Docker resources are available
3. Try full reset: `./scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh --reset`

### Key Import Failures

If key import fails:

1. Ensure wallets were created successfully in Bitcoin nodes
2. Check if encrypted mode is properly configured
3. Verify configuration file paths are correct
