# BTC E2E Scripts

This directory contains Bitcoin E2E (End-to-End) test scripts.
Each script automates the complete workflow from infrastructure setup to transaction execution.

## Documentation

For detailed transaction pattern explanations, technical references, and implementation status, see:

- **[E2E Transaction Patterns Guide](../../../../docs/chains/btc/operations/e2e-transaction-patterns.md)** - Key types, signature patterns, and workflow details

## Script List

| Script | Pattern | Signature Requirement | Address Format |
|--------|---------|----------------------|----------------|
| `e2e-p1-p2pkh-singlesig.sh` | P2PKH Single-sig (Pattern 1) | Single-sig | `1...` / `m...` |
| `e2e-p2-p2pkh-2of3.sh` | P2PKH 2-of-3 Multisig (Pattern 2) | 2-of-3 | `3...` / `2...` |
| `e2e-p3-p2sh-p2wpkh-singlesig.sh` | P2SH-P2WPKH Single-sig (Pattern 3) | Single-sig | `3...` / `2...` |
| `e2e-p4-p2sh-p2wsh-2of3.sh` | P2SH-P2WSH 2-of-3 Multisig (Pattern 4) | 2-of-3 | `3...` / `2...` |
| `e2e-p5-p2wpkh-singlesig.sh` | P2WPKH Native SegWit Single-sig (Pattern 5) | Single-sig | `bc1q...` / `bcrt1q...` |
| `e2e-p6-p2wsh-2of3.sh` | P2WSH Native SegWit 2-of-3 Multisig (Pattern 6) | 2-of-3 | `bc1q...` / `bcrt1q...` |
| `e2e-p7-p2wsh-3of3.sh` | P2WSH Native SegWit 3-of-3 Multisig (Pattern 7) | 3-of-3 | `bc1q...` / `bcrt1q...` |
| `e2e-p8-p2sh-p2wsh-3of3.sh` | P2SH-P2WSH 3-of-3 Multisig (Pattern 8) | 3-of-3 | `3...` / `2...` |
| `e2e-p9-p2tr-singlesig.sh` | P2TR Taproot Single-sig (Pattern 9) | Single-sig | `bc1p...` / `bcrt1p...` |
| `e2e-p10-p2tr-musig2.sh` | P2TR MuSig2 (Pattern 10) | N-of-N (framework) | `bc1p...` / `bcrt1p...` |
| `e2e-p11-p2tr-tapscript.sh` | P2TR Tapscript M-of-N (Pattern 11) | 2-of-3 (framework) | `bc1p...` / `bcrt1p...` |

## Verification Status

| Pattern | E2E Test Status | Last Verified | Notes |
|---------|----------------|---------------|-------|
| 1-9 | ✅ Fully operational | - | Complete transaction signing |
| 10 | ✅ Framework verified | 2026-01-16 | Infrastructure and workflow verified; MuSig2 protocol pending CLI implementation |
| 11 | ✅ Framework verified | 2026-01-16 | Infrastructure and workflow verified; Tapscript Script Path pending CLI implementation |

**Pattern 10 Verified Components:**

- Infrastructure setup (Docker, Bitcoin Core, MySQL)
- Wallet creation and configuration
- Taproot address generation (bech32m encoding)
- Descriptor import (2000 addresses)
- Payment workflow

**Pattern 10 Pending:**

- Full MuSig2 2-round protocol (nonce generation, partial signatures, aggregation)
- Requires completion of MuSig2 CLI commands

**Pattern 11 Verified Components:**

- Infrastructure setup (Docker, Bitcoin Core, MySQL)
- Wallet creation and configuration
- Taproot address generation (bech32m encoding)
- Descriptor import
- Payment workflow

**Pattern 11 Pending:**

- Script tree construction (Merkle tree)
- Internal key tweaking with Merkle root
- Control block construction
- Tapscript signing with sortedmulti_a
- Witness construction (sigs + script + control block)
- Requires completion of Tapscript CLI commands

## Usage

### Basic Execution

```bash
# Pattern 1: Single-sig E2E test
./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh

# Pattern 2: 2-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh

# Pattern 3: P2SH-P2WPKH Single-sig E2E test
./scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh

# Pattern 4: P2SH-P2WSH 2-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh

# Pattern 5: P2WPKH Native SegWit Single-sig E2E test
./scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh

# Pattern 6: P2WSH Native SegWit 2-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh

# Pattern 7: P2WSH Native SegWit 3-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh

# Pattern 8: P2SH-P2WSH 3-of-3 Multisig E2E test
./scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh

# Pattern 9: P2TR Taproot Single-sig E2E test
./scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh

# Pattern 10: P2TR MuSig2 N-of-N E2E test
./scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh

# Pattern 11: P2TR Tapscript M-of-N E2E test
./scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh
```

### Make Targets

#### Single Pattern Execution

Use `P=<pattern>` parameter to specify the test pattern (1-11):

```bash
# Run E2E test (default P=1)
make btc-e2e P=<pattern>

# Run with fresh state (recommended)
make btc-e2e-reset P=<pattern>

# Run with verbose output
make btc-e2e-verbose P=<pattern>

# Run in non-interactive mode (for CI/CD)
make btc-e2e-ci P=<pattern>

# Cleanup test environment
make btc-e2e-cleanup P=<pattern>

# Show help and pattern list
make btc-e2e-help
```

#### Parallel Pattern Execution

Run multiple patterns in parallel for faster CI/CD execution:

```bash
# Run all patterns in parallel (recommended for CI)
make btc-e2e-parallel

# Run all patterns in CI mode
make btc-e2e-ci-all

# Run specific patterns in parallel
make btc-e2e-parallel PATTERNS=1,2,3

# Run pattern range with limited parallelism
make btc-e2e-parallel PATTERNS=1-5 MAX_PARALLEL=3

# Run with verbose output
make btc-e2e-parallel PATTERNS=1-11 VERBOSE=true
```

**Parallel Execution Parameters:**

| Parameter | Default | Description |
|-----------|---------|-------------|
| `PATTERNS` | `1-11` | Comma-separated list or range (e.g., "1,2,3" or "1-11") |
| `MAX_PARALLEL` | `11` | Maximum number of concurrent processes |
| `VERBOSE` | `false` | Show real-time output from all processes |

**Benefits:**

- Uses SQLite backend for isolated database per pattern
- Each pattern runs independently without database conflicts
- Significantly reduces total CI execution time
- Logs are saved to `data/logs/e2e-parallel/`

**Examples:**

```bash
# Pattern 1: P2PKH Single-sig (fresh start)
make btc-e2e-reset P=1

# Pattern 3: P2SH-P2WPKH Single-sig
make btc-e2e P=3

# Pattern 9: P2TR Taproot Single-sig (verbose)
make btc-e2e-verbose P=9

# Pattern 10: P2TR MuSig2 N-of-N (CI mode)
make btc-e2e-ci P=10
```

**Available Patterns:**

| P | Pattern |
|---|---------|
| 1 | P2PKH Single-sig |
| 2 | P2PKH 2-of-3 Multisig |
| 3 | P2SH-P2WPKH Single-sig |
| 4 | P2SH-P2WSH 2-of-3 Multisig |
| 5 | P2WPKH Native SegWit Single-sig |
| 6 | P2WSH Native SegWit 2-of-3 Multisig |
| 7 | P2WSH Native SegWit 3-of-3 Multisig |
| 8 | P2SH-P2WSH 3-of-3 Multisig |
| 9 | P2TR Taproot Single-sig |
| 10 | P2TR MuSig2 N-of-N |
| 11 | P2TR Tapscript M-of-N |

### Common Options

| Option | Description |
|--------|-------------|
| `--cleanup` | Stop containers and cleanup state |
| `--reset` | Full reset and run from scratch |
| `--verbose` | Enable verbose output |
| `--non-interactive` | Run without prompts (for CI/CD) |
| `-h, --help` | Display help message |

## Required Configuration

Each script requires matching `address_type` in the corresponding config files:

### Pattern 1: P2PKH Single-sig

```yaml
# config/wallet/btc/watch.yaml, btc/keygen.yaml
address_type: "legacy"
```

### Pattern 2: P2PKH 2-of-3 Multisig

```yaml
# config/wallet/btc/watch.yaml, btc/keygen.yaml, btc/sign1.yaml, btc/sign2.yaml
address_type: "legacy"
```

### Pattern 3: P2SH-P2WPKH Single-sig

```yaml
# config/wallet/btc/watch.yaml, btc/keygen.yaml
address_type: "p2sh-segwit"
```

### Pattern 4: P2SH-P2WSH 2-of-3 Multisig

```yaml
# config/wallet/btc/watch.yaml, btc/keygen.yaml, btc/sign1.yaml, btc/sign2.yaml
address_type: "p2sh-segwit"
```

### Pattern 5: P2WPKH Native SegWit Single-sig

```yaml
# config/wallet/btc/watch.yaml, btc/keygen.yaml
address_type: "bech32"
```

### Pattern 6: P2WSH Native SegWit 2-of-3 Multisig

```yaml
# config/wallet/btc/watch.yaml, btc/keygen.yaml, btc/sign1.yaml, btc/sign2.yaml
address_type: "bech32"
```

### Pattern 7: P2WSH Native SegWit 3-of-3 Multisig

```yaml
# config/wallet/btc/watch.yaml, btc/keygen.yaml, btc/sign1.yaml, btc/sign2.yaml
address_type: "bech32"
```

### Pattern 8: P2SH-P2WSH 3-of-3 Multisig

```yaml
# config/wallet/btc/watch.yaml, btc/keygen.yaml, btc/sign1.yaml, btc/sign2.yaml
address_type: "p2sh-segwit"
```

### Pattern 9: P2TR Taproot Single-sig

```yaml
# config/wallet/btc/watch.yaml, btc/keygen.yaml
address_type: "taproot"
```

### Pattern 10: P2TR MuSig2 N-of-N

**Note:** Pattern 10 uses environment variable override instead of editing config files.

```bash
# Script automatically sets:
export WALLET_ADDRESS_TYPE="taproot"
```

Required:

- Uses `config/wallet/account/account_3of3.yaml` for N-of-N configuration
- Bitcoin Core v22.0+ with descriptor-based wallet support
- Address encoding: bech32m (`bcrt1p...` for regtest)

### Pattern 11: P2TR Tapscript M-of-N

**Note:** Pattern 11 uses environment variable override instead of editing config files.

```bash
# Script automatically sets:
export WALLET_ADDRESS_TYPE="taproot"
```

Required:

- Uses `config/wallet/account/account_2of3.yaml` for 2-of-3 configuration
- Bitcoin Core v22.0+ with Taproot/Tapscript support
- Address encoding: bech32m (`bcrt1p...` for regtest)

## Environment Variables

```bash
# RPC credentials (defaults are for regtest/development only)
RPC_USER=xyz
RPC_PASSWORD=xyz

# MySQL credentials (defaults are for regtest/development only)
MYSQL_ROOT_PASSWORD=root

# Database type selection (mysql or sqlite)
DB_TYPE=mysql  # Default: mysql
```

## Database Configuration

E2E scripts support two database backends:

### MySQL (Default)

Uses Docker MySQL container. This is the default and requires Docker running.

```bash
# Run with MySQL (default)
./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh
```

### SQLite (Lightweight Testing)

Uses local SQLite files. No Docker MySQL container required, enabling:

- Faster test startup (no Docker database container)
- Parallel test execution (each test can use separate DB files)
- Lighter CI/CD environments
- Testing without Docker MySQL overhead

```bash
# Run with SQLite
DB_TYPE=sqlite ./scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh
```

### SQLite Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_TYPE` | `mysql` | Database type: `mysql` or `sqlite` |
| `SQLITE_WATCH_DB_PATH` | `./data/sqlite/btc/watch.db` | SQLite database file for watch wallet |
| `SQLITE_KEYGEN_DB_PATH` | `./data/sqlite/btc/keygen.db` | SQLite database file for keygen wallet |
| `SQLITE_SIGN_DB_PATH` | `./data/sqlite/btc/sign.db` | SQLite database file for sign wallet |

### Wallet Configuration for SQLite

When using SQLite, update wallet config files to use SQLite:

```yaml
# config/wallet/btc/watch.yaml
database:
  type: "sqlite"
  sqlite:
    path: "./data/sqlite/btc/watch.db"
    debug: true
```

Or use environment variable override:

```bash
WALLET_DATABASE_TYPE=sqlite WALLET_DATABASE_SQLITE_PATH="./data/sqlite/btc/watch.db" watch ...
```

## Common Utilities

All E2E scripts use shared utility functions from:

- **`../btc_common.sh`** - BTC-specific common functions
- **`../../common.sh`** - General common functions (auto-sourced by btc_common.sh)

### Usage in E2E Scripts

```bash
# Source BTC common utilities (includes common.sh automatically)
source "${SCRIPT_DIR}/../btc_common.sh"

# Initialize config paths
btc_get_config_paths

# Use common functions
btc_check_prerequisites "watch keygen"
btc_setup_infrastructure "btc-watch btc-keygen"
btc_setup_wallets "btc-watch:watch btc-keygen:keygen"
btc_full_reset "watch keygen"
```

See `../README.md` for a complete list of available functions.

## Related Documentation

- [E2E Transaction Patterns Guide](../../../../docs/chains/btc/operations/e2e-transaction-patterns.md) - Pattern details
- [BTC Technical Reference](../../../../docs/chains/btc/README.md) - Bitcoin technical reference
- [Descriptor Examples](../../../../docs/chains/btc/descriptor/examples.md) - Descriptor examples
- [PSBT Developer Guide](../../../../docs/chains/btc/psbt/developer-guide.md) - PSBT developer guide
