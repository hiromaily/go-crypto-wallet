# Bitcoin Cash E2E Workflow

This document describes the complete Bitcoin Cash end-to-end workflow implemented in `e2e-workflow.sh`.

## Overview

The E2E workflow automates the complete Bitcoin Cash wallet operation from infrastructure setup to transaction execution. It serves as both a regression test tool and documentation of the standard Bitcoin Cash operation flow.

## Prerequisites

Before running the workflow, ensure the following are available:

- **Docker & Docker Compose**: For running Bitcoin Cash nodes and database
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
│   wallet-mysql     │     │   bch-watch     │
│   (MySQL)       │     │   (BCH Node)    │
└─────────────────┘     └─────────────────┘
                        ┌─────────────────┐
                        │   bch-keygen    │
                        │   (BCH Node)    │
                        └─────────────────┘
                        ┌─────────────────┐
                        │   bch-sign1     │
                        │   (BCH Node)    │
                        └─────────────────┘
                        ┌─────────────────┐
                        │   bch-sign2     │
                        │   (BCH Node)    │
                        └─────────────────┘
```

1. Start database container (`compose.yaml`)
2. Wait for database to be healthy
3. Start Bitcoin Cash node containers (`compose.bch.yaml`)
4. Wait for all Bitcoin Cash nodes to be healthy

### Phase 3: Wallet Setup

Creates wallets in each Bitcoin Cash node:

| Node       | Wallet Name |
|------------|-------------|
| bch-watch  | watch       |
| bch-keygen | keygen      |
| bch-sign1  | sign1       |
| bch-sign2  | sign2       |

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
│ 3. Import private keys into Bitcoin Cash node              │
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
│ 3. Import private keys into Bitcoin Cash nodes             │
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

1. Extract payment address from exported address file (CashAddr format)
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
./scripts/operation/bch/e2e-workflow.sh --reset

# Run workflow (assumes clean state or continuing from previous run)
./scripts/operation/bch/e2e-workflow.sh

# Run with verbose output
./scripts/operation/bch/e2e-workflow.sh --verbose
```

### Cleanup

```bash
# Stop containers and cleanup state
./scripts/operation/bch/e2e-workflow.sh --cleanup
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
| `RPC_USER`          | `xyz`   | Bitcoin Cash RPC username (for regtest)          |
| `RPC_PASSWORD`      | `xyz`   | Bitcoin Cash RPC password (for regtest)          |
| `WALLET_PASSPHRASE` | `test`  | Wallet passphrase for encrypted wallets          |

## Configuration Files

| File                              | Purpose                    |
|-----------------------------------|----------------------------|
| `config/wallet/bch/watch.yaml`    | Watch wallet configuration |
| `config/wallet/bch/keygen.yaml`   | Keygen wallet configuration|
| `config/wallet/bch/sign1.yaml`    | Sign1 wallet configuration |
| `config/wallet/bch/sign2.yaml`    | Sign2 wallet configuration |

## Generated Files

### Address Files (`data/address/bch/`)

- `address_client_*.csv` - Client addresses (non-multisig)
- `address_deposit_*.csv` - Deposit multisig addresses
- `address_payment_*.csv` - Payment multisig addresses
- `address_stored_*.csv` - Stored multisig addresses

### Public Key Files (`data/fullpubkey/bch/`)

- `fullpubkey_auth1_*.csv` - Full public keys from sign1
- `fullpubkey_auth2_*.csv` - Full public keys from sign2

### Transaction Files (`data/tx/bch/`)

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

## Bitcoin Cash Specific Notes

### Address Format

Bitcoin Cash uses the CashAddr format (`bitcoincash:qp...`) for addresses, which is different from Bitcoin's legacy, P2SH-SegWit, and Bech32 formats.

### Node Software

This workflow uses Bitcoin Cash Node (BCHN) version 28.0.1, which is the most widely used full node implementation for Bitcoin Cash. See [bitcoincashnode.org](https://bitcoincashnode.org/) for more information.

### Port Configuration

Default ports (to avoid conflicts with Bitcoin):

| Node       | Host Port |
|------------|-----------|
| bch-watch  | 28332     |
| bch-keygen | 29332     |
| bch-sign1  | 30332     |
| bch-sign2  | 31332     |

## Troubleshooting

### "No UTXOs available"

If the transaction phase reports no UTXOs:

1. Verify blocks were generated successfully
2. Check balance with: `watch -c config/wallet/bch/watch.yaml --coin bch monitor balance`
3. The script automatically generates 101 blocks; if this fails, manually run:

   ```bash
   docker exec bch-watch bitcoin-cli -regtest generatetoaddress 101 <address>
   ```

### Container Health Issues

If containers fail health checks:

1. Check container logs: `docker logs bch-watch`
2. Verify Docker resources are available
3. Try full reset: `./scripts/operation/bch/e2e-workflow.sh --reset`

### Key Import Failures

If key import fails:

1. Ensure wallets were created successfully in Bitcoin Cash nodes
2. Check if encrypted mode is properly configured
3. Verify configuration file paths are correct

### Address Format Issues

Bitcoin Cash uses CashAddr format. If you see address format errors:

1. Ensure `address_type = "bch-cashaddr"` is set in configuration files
2. Verify the wallet software supports CashAddr format

## E2E Script Verification Status

### Pattern 1: P2PKH Single-sig (`e2e-p1-p2pkh-singlesig.sh`)

**Status**: ⚠️ Partially Verified (4 bugs fixed, 1 blocking issue remains)

**Last Verified**: 2026-01-17
**Issue**: #404

#### Test Results

| Acceptance Criteria | Status | Notes |
|---------------------|--------|-------|
| Script executes with `--reset` | ⚠️ Partial | Progresses through most phases |
| All phases complete | ❌ Failed | Stops at transaction creation |
| Transaction broadcast | ❌ Not reached | Blocked by UTXO query issue |
| `--cleanup` works | ✅ Pass | Containers stop properly |
| `--verbose` works | ✅ Pass | Debug output displayed |
| `--help` works | ✅ Pass | Help message displayed |

#### Bugs Fixed

1. **Balance Detection** - Fixed to use `getbalance` with watch-only flag (BCH doesn't support `getbalances`)
2. **Missing Rescan** - Added blockchain rescan after block generation to detect UTXOs
3. **Address Format** - Fixed to use field 4 (legacy format) matching imported addresses
4. **Sign Command** - Fixed syntax to use `keygen sign signature --file`

#### Known Issues

- **UTXO Query Issue**: Transaction creation fails with "No utxo" error despite UTXOs existing in wallet
- Requires investigation of watch wallet UTXO query logic for BCH watch-only addresses

#### Execution Progress

- ✅ Prerequisites check
- ✅ Infrastructure setup
- ✅ Wallet creation  
- ✅ HD key generation
- ✅ Address import
- ✅ UTXO generation (50 BCH)
- ✅ Balance verification
- ✅ Payment request creation
- ❌ Transaction creation (blocked)
- ⏹️ Transaction signing (not reached)
- ⏹️ Transaction broadcast (not reached)

### Pattern 2: P2SH 2-of-3 (`e2e-p2-p2sh-2of3.sh`)

**Status**: ⚠️ Not Verified

**Notes**: Likely has similar address format issues as Pattern 1 (uses field 3 extraction on line 192)

### Pattern 3: P2SH 3-of-3 (`e2e-p3-p2sh-3of3.sh`)

**Status**: ⚠️ Not Verified

**Notes**: Likely has similar address format issues as Pattern 1 (uses field 3 extraction on line 192)
