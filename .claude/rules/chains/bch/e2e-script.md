---
paths: ["scripts/operation/bch/e2e/**"]
---

# BCH E2E Script Development Rules

Rules applied when creating or modifying Bitcoin Cash E2E scripts.
**Also load**: `.claude/rules/chains/e2e-script.md` for universal E2E rules (DB config, Makefile policy, security, etc.).

## Critical: Read BCH Task Context First

**MANDATORY**: Before working on any BCH E2E task, read:

- `docs/chains/bch/README.md` — **BCH vs BTC feature differences, workflow comparison, prohibited features**

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

| Document                              | Contents                                                  |
| ------------------------------------- | --------------------------------------------------------- |
| `docs/chains/bch/README.md`           | **CRITICAL**: BCH vs BTC differences, prohibited features |
| `scripts/operation/bch/bch_common.sh` | BCH common utility functions                              |
| `scripts/operation/common.sh`         | Shared utility functions                                  |

## Script Header Template

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

## Transaction Patterns

BCH supports **only 3 patterns** (compared to BTC's 11):

| Pattern | Description      | Address Type  | Address Format        | Signature |
| ------- | ---------------- | ------------- | --------------------- | --------- |
| 1       | P2PKH Single-sig | P2PKH (BIP44) | `m.../n...` (regtest) | Single    |
| 2       | P2SH 2-of-3      | P2SH          | `2...` (regtest)      | 2-of-3    |
| 3       | P2SH 3-of-3      | P2SH          | `2...` (regtest)      | 3-of-3    |

## Configuration Files

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

## Database Configuration (BCH-Specific)

### Database Debug Commands

```bash
# Via bch_common.sh abstraction (works for both SQLite and MySQL)
db_query "watch"  "SELECT wallet_address, account FROM address WHERE coin='bch' LIMIT 10"
db_query "watch"  "SELECT * FROM payment_request WHERE coin='bch'"
db_query "keygen" "SELECT * FROM account_key LIMIT 5"
```

| DB_TYPE  | Manual query command                                                                     |
| -------- | ---------------------------------------------------------------------------------------- |
| `sqlite` | `sqlite3 ./data/sqlite/bch/e2e.db "SELECT ..."`                                          |
| `mysql`  | `docker compose exec -T wallet-mysql mysql -u root -proot watch -e "SELECT ..."`        |

## BCH Implementation Rules

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

## Makefile Targets

BCH E2E targets use the `P=N` parameter style. Naming convention: `bch-e2e-pN`.

| Target                     | Description                          |
| -------------------------- | ------------------------------------ |
| `make bch-e2e-reset P=N`   | Fresh start with reset (recommended) |
| `make bch-e2e P=N`         | Run E2E test                         |
| `make bch-e2e-verbose P=N` | Run with verbose output              |
| `make bch-e2e-ci P=N`      | Run in non-interactive mode          |
| `make bch-e2e-cleanup P=N` | Cleanup only                         |

**Current Scripts**:

| Pattern | Script                      | Make Target  | Status |
| ------- | --------------------------- | ------------ | ------ |
| 1       | `e2e-p1-p2pkh-singlesig.sh` | `bch-e2e-p1` | ✅     |
| 2       | `e2e-p2-p2sh-2of3.sh`       | `bch-e2e-p2` | ✅     |
| 3       | `e2e-p3-p2sh-3of3.sh`       | `bch-e2e-p3` | ✅     |

## Docker Container Names

| Container    | Port  | Purpose            |
| ------------ | ----- | ------------------ |
| `bch-watch`  | 28332 | Watch wallet node  |
| `bch-keygen` | 29332 | Keygen wallet node |
| `bch-sign1`  | 30332 | Sign1 wallet node  |
| `bch-sign2`  | 31332 | Sign2 wallet node  |

## Common Errors

### "No utxo" Error

1. Verify addresses are correctly imported
2. Confirm block generation (101+) is complete
3. Check address format matches imported format

```bash
docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch getbalance "*" 1 true
docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch listunspent
```

### Address Format Issues

| Network | P2PKH Format       | P2SH Format        |
| ------- | ------------------ | ------------------ |
| Mainnet | `bitcoincash:q...` | `bitcoincash:p...` |
| Testnet | `bchtest:q...`     | `bchtest:p...`     |
| Regtest | `m.../n...`        | `2...`             |

### Balance Detection Issues

```bash
# For watch-only wallets
bitcoin-cli -regtest -rpcwallet=watch getbalance "*" 1 true
```

## Files NOT to Reference for BCH

```
internal/infrastructure/api/btc/btc/descriptor*.go
internal/infrastructure/api/btc/btc/psbt.go
internal/infrastructure/api/btc/btc/musig2.go
internal/application/usecase/*/btc/*musig2*.go
internal/application/usecase/*/btc/*descriptor*.go
```
