---
paths: ["scripts/operation/btc/e2e/**"]
---

# BTC E2E Script Development Rules

Rules applied when creating or modifying Bitcoin E2E scripts.
**Also load**: `.claude/rules/chains/e2e-script.md` for universal E2E rules (DB config, Makefile policy, security, etc.).

## Required Documentation

| Document                                                 | Contents                                                    |
| -------------------------------------------------------- | ----------------------------------------------------------- |
| `docs/chains/btc/operations/e2e-transaction-patterns.md` | Detailed specifications for all 11 patterns                 |
| `docs/chains/btc/overview/address-types.md`              | **CRITICAL**: Address types vs formats (taproot vs bech32m) |
| `scripts/operation/common.sh`                            | Common utility functions                                    |
| `pkg/config/README.md`                                   | Configuration override via environment variables            |
| `config/wallet/README.md`                                | Wallet configuration file policies                          |

## Script Header Template

```bash
#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern N: [Pattern Name]
# This script automates the complete Bitcoin workflow for [description]
# Usage: ./scripts/operation/btc/e2e/e2e-pN-{description}.sh [OPTIONS]
# Options:
#   --cleanup  Stop containers and cleanup state
#   --reset    Full reset and run from scratch
#   --verbose  Enable verbose output
#   --non-interactive  Run without prompts (for CI/CD)
#   -h, --help Display help message
#
# Reference Documentation:
#   docs/chains/btc/operations/e2e-transaction-patterns.md - E2E transaction patterns
#
# Transaction Pattern:
#   Pattern N: BTC [Address Type] [Signature Requirement]
#   - Address Type: [P2PKH/P2SH-P2WPKH/etc.]
#   - Address Format: `...` (Mainnet), `...` (Testnet/Regtest)
#   - Signature Requirement: [Single-sig/2-of-3/3-of-3]
#   - Descriptor: [descriptor format]
#
# Required Config Settings:
#   - config/wallet/btc/watch.yaml:  address_type: "[type]"
#   - config/wallet/btc/keygen.yaml: address_type: "[type]"
```

## Environment Variable Override Section

```bash
###############################################################################
# Environment Variable Overrides for Configuration
###############################################################################
# These environment variables override config file values.
# Priority: Environment Variables > Config File > Default Values
#
# Pattern N requires:
#   - address_type: "[type]" (derives key_type automatically)
export WALLET_ADDRESS_TYPE="[type]"
```

## Account Configuration Files

| Pattern         | Account Config                            |
| --------------- | ----------------------------------------- |
| Single-sig      | `config/wallet/account/account.yaml`      |
| 2-of-3 Multisig | `config/wallet/account/account_2of3.yaml` |
| 3-of-3 Multisig | `config/wallet/account/account_3of3.yaml` |

## Database Configuration (BTC-Specific)

| Variable         | Default                    | Description          |
| ---------------- | -------------------------- | -------------------- |
| `DB_TYPE`        | `sqlite`                   | Database type        |
| `SQLITE_DB_PATH` | `./data/sqlite/btc/e2e.db` | SQLite database file |

When `DB_TYPE=sqlite`, `btc_setup_infrastructure` initializes SQLite with all schemas. No Docker MySQL container is started.

### Database Debug Commands

```bash
# Via btc_common.sh abstraction (works for both SQLite and MySQL)
db_query "watch"  "SELECT wallet_address, account FROM address WHERE coin='btc' LIMIT 10"
db_query "watch"  "SELECT * FROM payment_request WHERE coin='btc'"
db_query "keygen" "SELECT * FROM account_key LIMIT 5"
```

| DB_TYPE  | Manual query command                                                                     |
| -------- | ---------------------------------------------------------------------------------------- |
| `sqlite` | `sqlite3 ./data/sqlite/btc/e2e.db "SELECT ..."`                                          |
| `mysql`  | `docker compose exec -T wallet-mysql mysql -u root -proot watch -e "SELECT ..."`        |

## Address Type Configuration

⚠️ **CRITICAL**: For P2TR patterns (9, 10, 11), use `address_type="taproot"` NOT `"bech32m"`. See `btc-terminology` skill for details.

`key_type` is **automatically derived** from `address_type` — do not set `WALLET_KEY_TYPE` explicitly.

| address_type          | Derived key_type | Use Case                                  |
| --------------------- | ---------------- | ----------------------------------------- |
| `legacy`              | `bip44`          | P2PKH (Pattern 1, 2)                      |
| `p2sh-segwit`         | `bip49`          | P2SH-P2WPKH/P2SH-P2WSH (Pattern 3, 4, 8) |
| `bech32`              | `bip84`          | Native SegWit (Pattern 5, 6, 7)           |
| `taproot` / `bech32m` | `bip86`          | Taproot (Pattern 9, 10, 11)               |

Reference: `AddrType.ToKeyType()` in `internal/domain/address/types.go`

## Transaction Patterns

| Pattern | Description                     | address_type  | Address Format         | Signature |
| ------- | ------------------------------- | ------------- | ---------------------- | --------- |
| 1       | P2PKH Single-sig                | `legacy`      | `m.../n...`            | Single    |
| 2       | P2PKH 2-of-3                    | `legacy`      | `2...` (P2SH)          | 2-of-3    |
| 3       | P2SH-P2WPKH Single-sig          | `p2sh-segwit` | `2...`                 | Single    |
| 4       | P2SH-P2WSH 2-of-3               | `p2sh-segwit` | `2...`                 | 2-of-3    |
| 5       | P2WPKH Native SegWit Single-sig | `bech32`      | `bcrt1q...`            | Single    |
| 6       | P2WSH 2-of-3                    | `bech32`      | `bcrt1q...` (62 chars) | 2-of-3    |
| 7       | P2WSH 3-of-3                    | `bech32`      | `bcrt1q...` (62 chars) | 3-of-3    |
| 8       | P2SH-P2WSH 3-of-3               | `p2sh-segwit` | `2...`                 | 3-of-3    |
| 9       | P2TR Taproot Single-sig         | `taproot`     | `bcrt1p...`            | Single    |
| 10      | P2TR MuSig2 N-of-N              | `taproot`     | `bcrt1p...`            | N-of-N    |
| 11      | P2TR Tapscript M-of-N           | `taproot`     | `bcrt1p...`            | M-of-N    |

## Makefile Targets

Add targets to `make/btc_e2e.mk`. Naming convention: `btc-e2e-pN`.

| Target                     | Description                          |
| -------------------------- | ------------------------------------ |
| `make btc-e2e-reset P=N`   | Fresh start with reset (recommended) |
| `make btc-e2e P=N`         | Run E2E test                         |
| `make btc-e2e-verbose P=N` | Run with verbose output              |
| `make btc-e2e-ci P=N`      | Run in non-interactive mode          |
| `make btc-e2e-cleanup P=N` | Cleanup only                         |

**Current Scripts**:

| Pattern | Script                            | Make Target   | Status |
| ------- | --------------------------------- | ------------- | ------ |
| 1       | `e2e-p1-p2pkh-singlesig.sh`       | `btc-e2e-p1`  | ✅     |
| 2       | `e2e-p2-p2pkh-2of3.sh`            | `btc-e2e-p2`  | ✅     |
| 3       | `e2e-p3-p2sh-p2wpkh-singlesig.sh` | `btc-e2e-p3`  | ✅     |
| 4       | `e2e-p4-p2sh-p2wsh-2of3.sh`       | `btc-e2e-p4`  | ✅     |
| 5       | `e2e-p5-p2wpkh-singlesig.sh`      | `btc-e2e-p5`  | ✅     |
| 6       | `e2e-p6-p2wsh-2of3.sh`            | `btc-e2e-p6`  | ✅     |
| 7       | `e2e-p7-p2wsh-3of3.sh`            | `btc-e2e-p7`  | ✅     |
| 8       | `e2e-p8-p2sh-p2wsh-3of3.sh`       | `btc-e2e-p8`  | ✅     |
| 9       | `e2e-p9-p2tr-singlesig.sh`        | `btc-e2e-p9`  | ✅     |
| 10      | `e2e-p10-p2tr-musig2.sh`          | `btc-e2e-p10` | ✅     |
| 11      | `e2e-p11-p2tr-tapscript.sh`       | `btc-e2e-p11` | ✅     |

## Common Errors

### "No utxo" Error

1. Verify Descriptor is correctly imported
2. Confirm block generation (101+) is complete
3. Check `address_type` is correct

```bash
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch getaddressinfo "<address>"
```

### address_type Mismatch

```bash
echo $WALLET_ADDRESS_TYPE
```

### duplicate key Error (Data from Previous Run)

```bash
make btc-e2e-reset P=N
```
