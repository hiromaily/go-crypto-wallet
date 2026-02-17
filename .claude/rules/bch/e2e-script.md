---
paths: "scripts/operation/bch/e2e/**"
---

# BCH E2E Script Development Rules

Rules applied when creating or modifying Bitcoin Cash E2E scripts.

## Critical: Read BCH Task Context First

**MANDATORY**: Before working on any BCH E2E task, read the BCH task context document:

- `docs/chains/bch/README.md` - **BCH vs BTC feature differences, workflow comparison, prohibited features**

This document explains what BCH can and cannot do compared to BTC.

## BCH Protocol Limitations (CRITICAL)

BCH does NOT support the following BTC features:

| Feature     | BTC | BCH    | Impact                     |
| ----------- | --- | ------ | -------------------------- |
| SegWit      | Yes | **NO** | No P2WPKH, P2WSH addresses |
| Taproot     | Yes | **NO** | No P2TR addresses          |
| Descriptor  | Yes | **NO** | Cannot use descriptor APIs |
| PSBT        | Yes | **NO** | Use raw transaction hex    |
| Schnorr     | Yes | **NO** | ECDSA only                 |
| MuSig2      | Yes | **NO** | Traditional multisig only  |
| BIP49/84/86 | Yes | **NO** | BIP44 only                 |

## Required Documentation

Read the following documents before creating or modifying scripts:

| Document                              | Contents                                                  |
| ------------------------------------- | --------------------------------------------------------- |
| `docs/chains/bch/README.md`           | **CRITICAL**: BCH vs BTC differences, prohibited features |
| `docs/chains/bch/README.md`           | Detailed BCH technical reference                          |
| `scripts/operation/bch/bch_common.sh` | BCH common utility functions                              |
| `scripts/operation/common.sh`         | Shared utility functions                                  |

## BCH Available Patterns

BCH supports **only 3 patterns** (compared to BTC's 11):

| Pattern | Description      | Address Type  | Address Format        | Signature |
| ------- | ---------------- | ------------- | --------------------- | --------- |
| 1       | P2PKH Single-sig | P2PKH (BIP44) | `m.../n...` (regtest) | Single    |
| 2       | P2SH 2-of-3      | P2SH          | `2...` (regtest)      | 2-of-3    |
| 3       | P2SH 3-of-3      | P2SH          | `2...` (regtest)      | 3-of-3    |

## Script Structure Conventions

### Header Comments

Each script must include header comments in the following format:

```bash
#!/usr/bin/env bash

# Bitcoin Cash E2E Workflow Script - Pattern N: [Pattern Name]
# This script automates the complete Bitcoin Cash workflow for [description]
# Usage: ./scripts/operation/bch/e2e/e2e-pN-{description}.sh [OPTIONS]
#
# Transaction Pattern:
#   Pattern N: BCH [Address Type] [Signature Requirement]
#   - Address Type: [P2PKH/P2SH]
#   - Address Format: CashAddr `bitcoincash:q...` (Mainnet), `m.../n...` (Regtest)
#   - Signature Requirement: [Single-sig/2-of-3/3-of-3]
#   - Key Derivation: m/44'/145'/account'/change/index (Mainnet)
#
# Note: BCH does NOT support SegWit or descriptor wallets
```

## Configuration File Policy

### Use BCH Config Files

BCH uses separate configuration files from BTC:

| File                            | Purpose                     |
| ------------------------------- | --------------------------- |
| `config/wallet/bch/watch.yaml`  | Watch wallet configuration  |
| `config/wallet/bch/keygen.yaml` | Keygen wallet configuration |
| `config/wallet/bch/sign1.yaml`  | Sign1 wallet configuration  |
| `config/wallet/bch/sign2.yaml`  | Sign2 wallet configuration  |

### Account Configuration Files

| Pattern         | Account Config                            |
| --------------- | ----------------------------------------- |
| Single-sig      | `config/wallet/account/account.yaml`      |
| 2-of-3 Multisig | `config/wallet/account/account_2of3.yaml` |
| 3-of-3 Multisig | `config/wallet/account/account_3of3.yaml` |

## Database Configuration

E2E scripts support two database backends via the `DB_TYPE` environment variable:

| DB_TYPE                | Description            | Docker MySQL | Use Case              |
| ---------------------- | ---------------------- | ------------ | --------------------- |
| `sqlite` (**default**) | Local SQLite file      | Not required | Fast testing, CI/CD   |
| `mysql`                | Docker MySQL container | Required     | Full integration test |

### Usage

```bash
# SQLite (default) - faster startup, no Docker MySQL needed
./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh --reset

# MySQL - traditional Docker-based testing
DB_TYPE=mysql ./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh --reset
```

### Database Debug Commands

Use the database abstraction functions from `bch_common.sh` for queries:

```bash
# Query addresses (works with both SQLite and MySQL)
db_query "watch" "SELECT wallet_address, account FROM address WHERE coin='bch' LIMIT 10"

# Query payment requests
db_query "watch" "SELECT * FROM payment_request WHERE coin='bch'"

# Query account keys
db_query "keygen" "SELECT * FROM account_key LIMIT 5"
```

**Manual queries** (when abstraction functions are not available):

| DB_TYPE  | Command                                                                          |
| -------- | -------------------------------------------------------------------------------- |
| `sqlite` | `sqlite3 ./data/sqlite/bch/e2e.db "SELECT ..."`                                  |
| `mysql`  | `docker compose exec -T wallet-mysql mysql -u root -proot watch -e "SELECT ..."` |

## BCH-Specific Implementation Rules

### DO NOT (Prohibited)

- ❌ DO NOT use Descriptor APIs
- ❌ DO NOT use PSBT format
- ❌ DO NOT use MuSig2
- ❌ DO NOT use Bech32/Bech32m addresses
- ❌ DO NOT use BIP49/84/86 derivation paths
- ❌ DO NOT reference BTC descriptor/psbt/musig2 files

### DO (Required)

- ✅ DO use CashAddr format for mainnet addresses
- ✅ DO use Raw Transaction Hex (not PSBT)
- ✅ DO use ECDSA signatures only
- ✅ DO use P2SH for multisig (not P2WSH)
- ✅ DO include SIGHASH_FORKID in signatures
- ✅ DO use BIP44 derivation path with coin type 145 (mainnet) or 1 (testnet)

## E2E Execution Rules

### ⚠️ MANDATORY: Always Use Makefile Targets

**AI Agents and developers MUST use Makefile targets to run E2E tests.**
Do NOT execute E2E scripts directly.

```bash
# ✅ CORRECT: Use Makefile target
make bch-e2e-reset P=1

# ❌ WRONG: Do not run scripts directly
./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh --reset
```

### Why Makefile Targets?

1. **Automatic Build**: `make bch-e2e-reset` includes `build-all` as a dependency
   - Incremental build: only rebuilds when Go sources change
   - No need to run `make build-all` separately
2. **Consistent Environment**: Properly validates pattern before execution
3. **Validated Patterns**: Validates pattern number before execution

### Available Makefile Targets

| Target                     | Description                          |
| -------------------------- | ------------------------------------ |
| `make bch-e2e-reset P=N`   | Fresh start with reset (recommended) |
| `make bch-e2e P=N`         | Run E2E test                         |
| `make bch-e2e-verbose P=N` | Run with verbose output              |
| `make bch-e2e-ci P=N`      | Run in non-interactive mode          |
| `make bch-e2e-cleanup P=N` | Cleanup only                         |

### Verification After Go Code Changes

```bash
# 1. Lint, build, and test
make go-lint && make check-build && make gotest

# 2. Run E2E test (build is automatic via dependency)
make bch-e2e-reset P=N
```

### Verification After Shell Script Changes

```bash
make shfmt
```

## Retry Limit

**If the fix-test cycle exceeds 5 iterations, organize progress and report.**

### Escalation Conditions

- Same error occurs repeatedly
- Deep understanding of BCH specifications required
- Large-scale Go code changes needed

### Progress Report Format

```markdown
## Progress Report

### Error Details

[Error message that occurred]

### Attempted Fixes

1. [Fix attempt 1]
2. [Fix attempt 2]

### Current State

[Description of current state]

### BCH-Specific Notes

- Address format used: [CashAddr/Legacy]
- Transaction format: Raw TX Hex

### Next Steps

[Required next actions]
```

## Security Rules

- ❌ Do NOT log private keys
- ❌ Do NOT use test passphrases/RPC credentials in production
- Reference: `docs/guidelines/security.md`

## Common BCH Errors

### "No utxo" Error

1. Verify addresses are correctly imported
2. Confirm block generation (101+) is complete
3. Check address format matches imported format

```bash
# Debug - Check balance
docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch getbalance "*" 1 true

# Debug - List UTXOs
docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch listunspent
```

### Address Format Issues

BCH uses different address formats:

| Network | P2PKH Format       | P2SH Format        |
| ------- | ------------------ | ------------------ |
| Mainnet | `bitcoincash:q...` | `bitcoincash:p...` |
| Testnet | `bchtest:q...`     | `bchtest:p...`     |
| Regtest | `m.../n...`        | `2...`             |

### Balance Detection Issues

BCH Node may not support all `getbalances` features. Use:

```bash
# For watch-only wallets
bitcoin-cli -regtest -rpcwallet=watch getbalance "*" 1 true
```

## Docker Container Names

BCH uses separate containers from BTC:

| Container    | Port  | Purpose            |
| ------------ | ----- | ------------------ |
| `bch-watch`  | 28332 | Watch wallet node  |
| `bch-keygen` | 29332 | Keygen wallet node |
| `bch-sign1`  | 30332 | Sign1 wallet node  |
| `bch-sign2`  | 31332 | Sign2 wallet node  |

## Makefile Targets

BCH E2E targets follow the pattern `bch-e2e-pN`:

| Pattern | Make Target  | Script                      |
| ------- | ------------ | --------------------------- |
| 1       | `bch-e2e-p1` | `e2e-p1-p2pkh-singlesig.sh` |
| 2       | `bch-e2e-p2` | `e2e-p2-p2sh-2of3.sh`       |
| 3       | `bch-e2e-p3` | `e2e-p3-p2sh-3of3.sh`       |

## Files NOT to Reference for BCH

These BTC files should **NEVER** be used for BCH:

```
internal/infrastructure/api/btc/btc/descriptor*.go
internal/infrastructure/api/btc/btc/psbt.go
internal/infrastructure/api/btc/btc/musig2.go
internal/application/usecase/*/btc/*musig2*.go
internal/application/usecase/*/btc/*descriptor*.go
```

## Related Documentation

- [BCH Documentation](../../../docs/chains/bch/README.md) - **CRITICAL**: Read first
- [BCH Technical Reference](../../../docs/chains/bch/README.md)
- [BCH E2E README](../../../scripts/operation/bch/README.md)
