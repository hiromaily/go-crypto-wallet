# Fix BCH E2E Pattern 2 Errors #{issue_number}

Fix errors in BCH E2E test (Pattern 2: P2SH 2-of-3 Multisig).

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

This command diagnoses and fixes errors when running BCH E2E Pattern 2 test.
The script already exists, so focus on identifying and fixing the root cause.

### Pattern 2 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 2 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2SH (BIP44 + BIP11) |
| **Script Type** | 2-of-3 Multisig |
| **Address Format** | `2...` (regtest P2SH) |
| **Signature Requirement** | 2-of-3 Multisig (any 2 of Keygen, Sign1, Sign2) |
| **Key Derivation** | `m/44'/1'/account'/change/index` (regtest) |
| **Required Wallets** | watch, keygen, sign1, sign2 |
| **Transaction Format** | Raw TX Hex (NOT PSBT) |
| **Account Config** | `account_2of3.yaml` |

### BCH vs BTC Pattern 2 Differences

| Item | BTC Pattern 2 | BCH Pattern 2 |
|------|---------------|---------------|
| Descriptor | `sh(multi(2,...))` | **NOT supported** |
| Transaction Format | PSBT | **Raw TX Hex** |
| Address Import | Descriptor import | **Address export/import** |
| Coin Type | 0 (mainnet), 1 (testnet) | 145 (mainnet), **1 (regtest)** |
| Docker Containers | `btc-*` | **`bch-*`** |
| SegWit Equivalent | P2SH-P2WSH (Pattern 4) | **N/A - No SegWit** |

### Differences from BCH Pattern 1 (Single-sig)

| Item | Pattern 1 | Pattern 2 |
|------|-----------|-----------|
| Signature Requirement | Single-sig (1) | **2-of-3 Multisig** |
| Required Wallets | keygen only | **keygen + sign1 + sign2** |
| Address Format | `m.../n...` (P2PKH) | **`2...` (P2SH)** |
| fullpubkey Exchange | Not required | **Required** |
| Account Config | `account.yaml` | **`account_2of3.yaml`** |
| Multisig Creation | Not needed | **Required** |

### Differences from BCH Pattern 3 (3-of-3 Multisig)

| Item | Pattern 2 | Pattern 3 |
|------|-----------|-----------|
| Signature Requirement | **2-of-3** (any 2) | 3-of-3 (all required) |
| Signatures Needed | 2 (keygen + sign1) | 3 (keygen + sign1 + sign2) |
| Redundancy | **Yes** (1 key loss OK) | No |
| Account Config | `account_2of3.yaml` | `account_3of3.yaml` |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-bch-e2e-p2`
- **Commit type**: `fix(bch)`
- **Scope**: BCH E2E Pattern 2

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Fix locally without creating branch or PR.

## Pattern 2 Specific Documentation

In addition to Required Documentation in prerequisites, refer to:

- @scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh - **Target script**
- @scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh - Pattern 3 script (3-of-3 reference)
- @scripts/operation/bch/bch_common.sh - BCH common functions
- @scripts/operation/bch/README.md - BCH E2E workflow documentation

## BCH-Specific Reminders

Before debugging, remember these BCH limitations:

| Feature | Available? | Alternative |
|---------|------------|-------------|
| Descriptor APIs | **NO** | Use address export/import |
| PSBT format | **NO** | Use Raw TX Hex |
| P2WSH (SegWit Multisig) | **NO** | Use P2SH |
| Schnorr/MuSig2 | **NO** | Use ECDSA |
| BIP49/84/86 | **NO** | Use BIP44 only |

## 2-of-3 Multisig Workflow

```
┌─────────────────────────────────────────────────────────┐
│              BCH 2-of-3 MULTISIG FLOW                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Keygen: Create seed, HD keys, import privkey        │
│  2. Sign1/Sign2: Create seed, HD keys, import privkey   │
│  3. Sign1/Sign2: Export fullpubkey                      │
│  4. Keygen: Import fullpubkeys from Sign1/Sign2         │
│  5. Keygen: Create 2-of-3 multisig addresses            │
│  6. Keygen: Export addresses                            │
│  7. Watch: Import addresses                             │
│          ↓                                              │
│  8. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  9. Keygen Wallet: Sign (1st signature - ECDSA)        │
│          ↓                                              │
│  10. Sign1 Wallet: Sign (2nd signature - ECDSA)        │
│          ↓                                              │
│  11. Watch Wallet: Broadcast transaction                │
│                                                         │
│  (Sign2 not required - 2-of-3 already satisfied)       │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Error Diagnosis Steps

### Step 1: Reproduce Error

```bash
# Full reset and run E2E test
./scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh --reset

# Or with verbose output
./scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh --reset --verbose
```

Check error message and categorize below.

### Step 2: Identify Error Category

| Error Message | Category | Reference Section |
|---------------|----------|-------------------|
| `No utxo` | UTXO-related | [UTXO Errors](#utxo-errors) |
| `connection refused` | Infrastructure | [Infrastructure Errors](#infrastructure-errors) |
| `wallet not found` | Wallet | [Wallet Errors](#wallet-errors) |
| `signing failed` | Signing | [Signing Errors](#signing-errors) |
| `fullpubkey` | Fullpubkey | [Fullpubkey Errors](#fullpubkey-errors) |
| `multisig` | Multisig | [Multisig Errors](#multisig-errors) |
| `address format` | Address | [Address Errors](#address-errors) |
| `duplicate key` | DB | [Database Errors](#database-errors) |

## Pattern 2 Specific Errors and Solutions

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

1. **Multisig addresses not created correctly**

   ```bash
   # Debug: Check address in database
   docker compose exec -T wallet-db mysql -u root -proot watch -e \
     "SELECT wallet_address, account FROM address WHERE coin='bch' AND account='payment' LIMIT 5"
   ```

   Address should start with `2...` (P2SH format in regtest).

2. **Address format mismatch in CSV extraction**

   ```bash
   # Check field extraction (Pattern 2 uses field 4 for P2SH)
   grep "cut -d',' -f" scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh
   
   # Verify CSV structure
   cat data/address/bch/address_payment_*.csv | head -5
   ```

3. **Insufficient block generation**

   ```bash
   # Check block count
   docker exec bch-watch bitcoin-cli -regtest getblockcount
   # Should be 101 or more
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
# Check container status (all 4 BCH containers)
docker compose ps | grep bch

# Restart BCH containers
docker compose -f compose.yaml -f compose.bch.yaml up -d \
  bch-watch bch-keygen bch-sign1 bch-sign2

# Check container health
docker exec bch-watch bitcoin-cli -regtest getblockchaininfo
```

### Wallet Errors

#### "wallet not found" Error

**Solution**:

```bash
# List wallets (all 4)
docker exec bch-watch bitcoin-cli -regtest listwallets
docker exec bch-keygen bitcoin-cli -regtest listwallets
docker exec bch-sign1 bitcoin-cli -regtest listwallets
docker exec bch-sign2 bitcoin-cli -regtest listwallets

# Create wallets if missing
docker exec bch-watch bitcoin-cli -regtest createwallet "watch" true true
docker exec bch-keygen bitcoin-cli -regtest createwallet "keygen" false true
docker exec bch-sign1 bitcoin-cli -regtest createwallet "sign1" false true
docker exec bch-sign2 bitcoin-cli -regtest createwallet "sign2" false true
```

### Signing Errors

#### Signing Failed Error (2-of-3)

**Cause**: Private key not imported, incorrect multisig setup, or only 1 signature

**Steps**:

```bash
# 1. Check keygen wallet has keys
docker exec bch-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  getwalletinfo

# 2. Check sign1 wallet has keys
docker exec bch-sign1 bitcoin-cli -regtest -rpcwallet=sign1 \
  getwalletinfo

# 3. Verify 2-of-3 only needs 2 signatures
# Pattern 2 should use keygen + sign1 (sign2 is NOT required)
```

### Fullpubkey Errors

#### Fullpubkey Export/Import Failed

**Steps**:

```bash
# Check fullpubkey files exist
ls -la data/fullpubkey/bch/

# Verify file contents
cat data/fullpubkey/bch/fullpubkey_auth1_*.csv | head -5
cat data/fullpubkey/bch/fullpubkey_auth2_*.csv | head -5

# Re-export if needed
sign1 -c config/wallet/bch/sign1.yaml --coin bch --wallet sign1 export fullpubkey
sign2 -c config/wallet/bch/sign2.yaml --coin bch --wallet sign2 export fullpubkey
```

### Multisig Errors

#### Multisig Address Creation Failed

**Steps**:

```bash
# 1. Verify fullpubkeys were imported into keygen
# Check database
docker compose exec -T wallet-db mysql -u root -proot keygen -e \
  "SELECT * FROM auth_fullpubkey LIMIT 5"

# 2. Re-create multisig addresses
keygen -c config/wallet/bch/keygen.yaml --coin bch create multisig --account payment
```

### Address Errors

#### Address Format Mismatch

**BCH P2SH Address Formats by Network**:

| Network | P2SH Format |
|---------|-------------|
| Mainnet | `bitcoincash:p...` |
| Testnet | `bchtest:p...` |
| Regtest | `2...` |

**Steps**:

```bash
# Check exported address format (should be P2SH `2...` for multisig)
cat data/address/bch/address_payment_*.csv | head -5

# Verify field extraction
# Pattern 2 CSV field 4 should contain P2SH address
```

### Database Errors

#### "duplicate key" Error

**Solution**:

```bash
# Full reset
./scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh --reset
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

# Check DB addresses (should be P2SH `2...` format)
docker compose exec -T wallet-db mysql -u root -proot watch -e \
  "SELECT wallet_address, account FROM address WHERE coin='bch' LIMIT 10"

# Check payment requests
docker compose exec -T wallet-db mysql -u root -proot watch -e \
  "SELECT * FROM payment_request WHERE coin='bch'"

# Check fullpubkeys in keygen DB
docker compose exec -T wallet-db mysql -u root -proot keygen -e \
  "SELECT * FROM auth_fullpubkey LIMIT 5"
```

### Log Check

```bash
# Run in verbose mode
./scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh --verbose --reset

# Check container logs
docker logs bch-watch --tail 50
docker logs bch-keygen --tail 50
docker logs bch-sign1 --tail 50
docker logs bch-sign2 --tail 50
```

## Identifying Files to Fix

| Error Type | Target File |
|------------|-------------|
| Script logic | `scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh` |
| BCH common functions | `scripts/operation/bch/bch_common.sh` |
| Common functions | `scripts/operation/common.sh` |
| BCH API (Go) | `internal/infrastructure/api/btc/bch/` |
| Multisig creation | `internal/application/usecase/keygen/btc/create_multisig*.go` |
| Watch wallet | `internal/application/usecase/watch/btc/` |
| Config loading | `pkg/config/` |

## Related Code (Go)

| Path | Role |
|------|------|
| `internal/infrastructure/api/btc/bch/` | **BCH API implementation (overrides BTC)** |
| `internal/infrastructure/api/btc/btc/` | BTC API (inherited by BCH) |
| `internal/application/usecase/keygen/btc/create_multisig*.go` | **Multisig address creation** |
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

- Do not break BCH Pattern 1 (`e2e-p1-p2pkh-singlesig.sh`)
- Do not break BCH Pattern 3 (`e2e-p3-p2sh-3of3.sh`)
- When modifying `bch_common.sh`, verify impact on other BCH patterns
- When modifying `common.sh`, verify impact on both BTC and BCH patterns

### BCH-Specific Considerations

- BCH uses **Raw TX Hex**, not PSBT
- BCH requires **SIGHASH_FORKID** for replay protection
- BCH regtest uses **P2SH address format** (`2...`) for multisig
- BCH 2-of-3 needs only **2 signatures** (keygen + sign1), sign2 is optional
- BCH node API responses may differ from BTC (override in BCH layer)

> **Note**: For build rules, verification commands, security, see @.claude/rules/bch/e2e-script.md

## Cleanup

```bash
# Stop containers only
./scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh --cleanup

# Full reset (including data)
./scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh --reset
```

## Known Issues

As of 2026-01-17, BCH E2E Pattern 2 has the following known issues:

- **Not Verified**: Script exists but has not been fully tested
- **Likely has address format issues** similar to Pattern 1 (field extraction)

See `scripts/operation/bch/README.md` for current verification status.
