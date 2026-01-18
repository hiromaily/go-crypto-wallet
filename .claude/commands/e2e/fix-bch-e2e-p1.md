# Fix BCH E2E Pattern 1 Errors #{issue_number}

Fix errors in BCH E2E test (Pattern 1: P2PKH Single-sig).

## Prerequisites (MANDATORY)

**Read the following documents BEFORE starting any work:**

1. **@docs/task-contexts/chains/bch.md** - **CRITICAL: BCH vs BTC feature differences, prohibited features, workflow comparison**
2. @.claude/rules/bch/e2e-script.md - BCH E2E common rules (build, verification, escalation, security)

> **WARNING**: BCH does NOT support SegWit, Taproot, Descriptor, PSBT, MuSig2. Read the task context document to understand what is available.

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command diagnoses and fixes errors when running BCH E2E Pattern 1 test.
The script already exists, so focus on identifying and fixing the root cause.

### Pattern 1 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 1 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2PKH (BIP44) |
| **Script Type** | Single-sig |
| **Address Format** | `m.../n...` (regtest P2PKH) |
| **Signature Requirement** | Single-sig (1 signature) |
| **Key Derivation** | `m/44'/1'/account'/change/index` (regtest) |
| **Required Wallets** | watch, keygen |
| **Transaction Format** | Raw TX Hex (NOT PSBT) |

### BCH vs BTC Pattern 1 Differences

| Item | BTC Pattern 1 | BCH Pattern 1 |
|------|---------------|---------------|
| Descriptor | `pkh([fp/44'/0'/0']xpub/*)` | **NOT supported** |
| Transaction Format | PSBT | **Raw TX Hex** |
| Address Import | Descriptor import | **Address export/import** |
| Coin Type | 0 (mainnet), 1 (testnet) | 145 (mainnet), **1 (regtest)** |
| Docker Container | `btc-watch`, `btc-keygen` | **`bch-watch`**, **`bch-keygen`** |

### Differences from BCH Pattern 3 (3-of-3 Multisig)

| Item | Pattern 1 | Pattern 3 |
|------|-----------|-----------|
| Signature Requirement | Single-sig (1) | 3-of-3 Multisig |
| Required Wallets | keygen only | keygen + sign1 + sign2 |
| Address Format | `m.../n...` (P2PKH) | `2...` (P2SH) |
| fullpubkey Exchange | Not required | Required |
| Account Config | `account.yaml` | `account_3of3.yaml` |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-bch-e2e-p1`
- **Commit type**: `fix(bch)`
- **Scope**: BCH E2E Pattern 1

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Fix locally without creating branch or PR.

## Pattern 1 Specific Documentation

In addition to Required Documentation in prerequisites, refer to:

- @scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh - **Target script**
- @scripts/operation/bch/bch_common.sh - BCH common functions
- @scripts/operation/bch/README.md - BCH E2E workflow documentation

## BCH-Specific Reminders

Before debugging, remember these BCH limitations:

| Feature | Available? | Alternative |
|---------|------------|-------------|
| Descriptor APIs | **NO** | Use address export/import |
| PSBT format | **NO** | Use Raw TX Hex |
| SegWit addresses | **NO** | Use P2PKH (`q...`) only |
| Schnorr/MuSig2 | **NO** | Use ECDSA |
| BIP49/84/86 | **NO** | Use BIP44 only |

## Error Diagnosis Steps

### Step 1: Reproduce Error

```bash
# Full reset and run E2E test
./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh --reset

# Or with verbose output
./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh --reset --verbose
```

Check error message and categorize below.

### Step 2: Identify Error Category

| Error Message | Category | Reference Section |
|---------------|----------|-------------------|
| `No utxo` | UTXO-related | [UTXO Errors](#utxo-errors) |
| `connection refused` | Infrastructure | [Infrastructure Errors](#infrastructure-errors) |
| `wallet not found` | Wallet | [Wallet Errors](#wallet-errors) |
| `signing failed` | Signing | [Signing Errors](#signing-errors) |
| `address format` | Address | [Address Errors](#address-errors) |
| `duplicate key` | DB | [Database Errors](#database-errors) |

## Pattern 1 Specific Errors and Solutions

### UTXO Errors

#### "No utxo" Error

**Symptoms**:

```
Transaction creation failed
This could indicate:
  - No payment requests in database
  - No UTXOs available for payment account
  - UTXOs not mature enough (need 100+ confirmations)
```

**Causes and Solutions**:

1. **Addresses not imported correctly**

   ```bash
   # Debug: Check wallet balance
   docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch \
     getbalance "*" 1 true

   # Debug: List UTXOs
   docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch \
     listunspent
   ```

2. **Insufficient block generation**

   ```bash
   # Check block count
   docker exec bch-watch bitcoin-cli -regtest getblockcount
   # Should be 101 or more
   ```

3. **Address format mismatch**

   BCH regtest uses legacy format (`m.../n...`), not CashAddr.

   ```bash
   # Check address in database (using abstraction function)
   db_query "watch" "SELECT wallet_address FROM address WHERE coin='bch' AND account='payment' LIMIT 5"
   ```

4. **Rescan not performed**

   ```bash
   # Force rescan
   docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch \
     rescanblockchain
   ```

### Infrastructure Errors

#### "connection refused" Error

**Solution**:

```bash
# Check container status
docker compose ps | grep bch

# Restart BCH containers
docker compose -f compose.yaml -f compose.bch.yaml up -d bch-watch bch-keygen

# Check container health
docker exec bch-watch bitcoin-cli -regtest getblockchaininfo
```

### Wallet Errors

#### "wallet not found" Error

**Solution**:

```bash
# List wallets
docker exec bch-watch bitcoin-cli -regtest listwallets

# Create wallets if missing
docker exec bch-watch bitcoin-cli -regtest createwallet "watch" true true
docker exec bch-keygen bitcoin-cli -regtest createwallet "keygen" false true
```

### Signing Errors

#### Signing Failed Error

**Cause**: Private key not imported or incorrect key format

**Steps**:

```bash
# 1. Check keygen wallet keys
docker exec bch-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  getwalletinfo

# 2. Check if wallet is encrypted
docker exec bch-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  walletpassphrase "test" 60
```

### Address Errors

#### Address Format Mismatch

**BCH Address Formats by Network**:

| Network | P2PKH | P2SH |
|---------|-------|------|
| Mainnet | `bitcoincash:q...` | `bitcoincash:p...` |
| Testnet | `bchtest:q...` | `bchtest:p...` |
| Regtest | `m.../n...` | `2...` |

**Steps**:

```bash
# Check exported address format
cat data/address/bch/address_payment_*.csv | head -5

# Verify field extraction in script
# BCH P1 script uses field 4 for legacy format
grep "cut -d',' -f" scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh
```

### Database Errors

#### "duplicate key" Error

**Solution**:

```bash
# Full reset
./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh --reset
```

## Debug Commands

### Status Check

```bash
# BCH node status
docker exec bch-watch bitcoin-cli -regtest getblockchaininfo

# Wallet balance (watch-only compatible)
docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch \
  getbalance "*" 1 true

# UTXO list
docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch listunspent
```

### Database Queries

> **Note**: Database queries depend on `DB_TYPE` setting. See "Database Debug Commands" in @.claude/rules/bch/e2e-script.md

```bash
# Check DB addresses (using abstraction function)
db_query "watch" "SELECT wallet_address, account FROM address WHERE coin='bch' LIMIT 10"

# Check payment requests
db_query "watch" "SELECT * FROM payment_request WHERE coin='bch'"
```

### Log Check

```bash
# Run in verbose mode
./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh --verbose --reset

# Check container logs
docker logs bch-watch --tail 50
docker logs bch-keygen --tail 50
```

## Identifying Files to Fix

| Error Type | Target File |
|------------|-------------|
| Script logic | `scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh` |
| BCH common functions | `scripts/operation/bch/bch_common.sh` |
| Common functions | `scripts/operation/common.sh` |
| BCH API (Go) | `internal/infrastructure/api/btc/bch/` |
| Watch wallet | `internal/application/usecase/watch/btc/` |
| Config loading | `pkg/config/` |

## Related Code (Go)

| Path | Role |
|------|------|
| `internal/infrastructure/api/btc/bch/` | **BCH API implementation (overrides BTC)** |
| `internal/infrastructure/api/btc/btc/` | BTC API (inherited by BCH) |
| `internal/application/usecase/keygen/btc/` | Key generation use case |
| `internal/application/usecase/watch/btc/` | Watch wallet use case |

### Files NOT to Use for BCH

> **WARNING**: These BTC files should NEVER be referenced for BCH fixes:

```
internal/infrastructure/api/btc/btc/descriptor*.go  # No Descriptor
internal/infrastructure/api/btc/btc/psbt.go         # No PSBT
internal/infrastructure/api/btc/btc/musig2.go       # No MuSig2
internal/application/usecase/*/btc/*musig2*.go
internal/application/usecase/*/btc/*descriptor*.go
```

## Cautions

### Avoid Impact on Existing Scripts

- Do not break BCH Pattern 2 (`e2e-p2-p2sh-2of3.sh`)
- Do not break BCH Pattern 3 (`e2e-p3-p2sh-3of3.sh`)
- When modifying `bch_common.sh`, verify impact on other BCH patterns
- When modifying `common.sh`, verify impact on both BTC and BCH patterns

### BCH-Specific Considerations

- BCH uses **Raw TX Hex**, not PSBT
- BCH requires **SIGHASH_FORKID** for replay protection
- BCH regtest uses **legacy address format** (`m.../n...`)
- BCH node API responses may differ from BTC (override in BCH layer)

> **Note**: For build rules, verification commands, security, see @.claude/rules/bch/e2e-script.md

## Cleanup

```bash
# Stop containers only
./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh --cleanup

# Full reset (including data)
./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh --reset
```

## Known Issues

As of 2026-01-17, BCH E2E Pattern 1 has the following known issues:

- **UTXO Query Issue**: Transaction creation may fail with "No utxo" error
- Requires investigation of watch wallet UTXO query logic for BCH watch-only addresses

See `scripts/operation/bch/README.md` for current verification status.
