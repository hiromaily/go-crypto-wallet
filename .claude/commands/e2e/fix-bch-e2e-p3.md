# Fix BCH E2E Pattern 3 Errors #{issue_number}

Fix errors in BCH E2E test (Pattern 3: P2SH 3-of-3 Multisig).

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

This command diagnoses and fixes errors when running BCH E2E Pattern 3 test.
The script already exists, so focus on identifying and fixing the root cause.

### Pattern 3 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 3 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2SH (BIP44 + BIP11) |
| **Script Type** | 3-of-3 Multisig |
| **Address Format** | `2...` (regtest P2SH) |
| **Signature Requirement** | 3-of-3 Multisig (**ALL** of Keygen, Sign1, Sign2) |
| **Key Derivation** | `m/44'/1'/account'/change/index` (regtest) |
| **Required Wallets** | watch, keygen, sign1, sign2 |
| **Transaction Format** | Raw TX Hex (NOT PSBT) |
| **Account Config** | `account_3of3.yaml` |
| **Multisig Accounts** | deposit, payment, stored (**NOT client**) |

### BCH vs BTC Pattern 3 Differences

| Item | BTC Pattern 3 | BCH Pattern 3 |
|------|---------------|---------------|
| Descriptor | `sh(multi(3,...))` | **NOT supported** |
| Transaction Format | PSBT | **Raw TX Hex** |
| Address Import | Descriptor import | **Address export/import** |
| Coin Type | 0 (mainnet), 1 (testnet) | 145 (mainnet), **1 (regtest)** |
| Docker Containers | `btc-*` | **`bch-*`** |
| SegWit Equivalent | P2SH-P2WSH | **N/A - No SegWit** |

### Differences from BCH Pattern 1 (Single-sig)

| Item | Pattern 1 | Pattern 3 |
|------|-----------|-----------|
| Signature Requirement | Single-sig (1) | **3-of-3 Multisig (ALL 3)** |
| Required Wallets | keygen only | **keygen + sign1 + sign2** |
| Address Format | `m.../n...` (P2PKH) | **`2...` (P2SH)** |
| fullpubkey Exchange | Not required | **Required** |
| Account Config | `account.yaml` | **`account_3of3.yaml`** |
| Multisig Creation | Not needed | **Required** |
| Use Case | Hot wallet operations | **Cold storage** |

### Differences from BCH Pattern 2 (2-of-3 Multisig)

| Item | Pattern 2 | Pattern 3 |
|------|-----------|-----------|
| Signature Requirement | 2-of-3 (any 2) | **3-of-3 (ALL required)** |
| Signatures Needed | 2 (keygen + sign1) | **3 (keygen + sign1 + sign2)** |
| Redundancy | Yes (1 key loss OK) | **No (all keys required)** |
| Account Config | `account_2of3.yaml` | `account_3of3.yaml` |
| Security Level | High | **Maximum** |
| Use Case | Payment operations | **Cold storage (stored)** |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-bch-e2e-p3`
- **Commit type**: `fix(bch)`
- **Scope**: BCH E2E Pattern 3

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Fix locally without creating branch or PR.

## Pattern 3 Specific Documentation

In addition to Required Documentation in prerequisites, refer to:

- @scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh - **Target script**
- @scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh - Pattern 2 script (2-of-3 reference)
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

## 3-of-3 Multisig Workflow

```
┌─────────────────────────────────────────────────────────┐
│              BCH 3-of-3 MULTISIG FLOW                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Keygen: Create seed, HD keys, import privkey        │
│  2. Sign1/Sign2: Create seed, HD keys, import privkey   │
│  3. Sign1/Sign2: Export fullpubkey                      │
│  4. Keygen: Import fullpubkeys from Sign1/Sign2         │
│  5. Keygen: Create 3-of-3 multisig addresses            │
│     (deposit, payment, stored - NOT client)             │
│  6. Keygen: Export addresses                            │
│  7. Watch: Import addresses                             │
│          ↓                                              │
│  8. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  9. Keygen Wallet: Sign (1st signature - ECDSA)        │
│          ↓                                              │
│  10. Sign1 Wallet: Sign (2nd signature - ECDSA)        │
│          ↓                                              │
│  11. Sign2 Wallet: Sign (3rd signature - ECDSA)        │
│          ↓                                              │
│  12. Watch Wallet: Broadcast transaction                │
│                                                         │
│  (ALL 3 signatures required - no redundancy)            │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## Error Diagnosis Steps

### Step 1: Reproduce Error

```bash
# Full reset and run E2E test
./scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh --reset

# Or with verbose output
./scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh --reset --verbose
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
| `incomplete signatures` | Signing | [Incomplete Signature Errors](#incomplete-signature-errors) |

## Pattern 3 Specific Errors and Solutions

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
   # Debug: Check address in database (using abstraction function)
   db_query "watch" "SELECT wallet_address, account FROM address WHERE coin='bch' AND account='payment' LIMIT 5"
   ```

   Address should start with `2...` (P2SH format in regtest).

2. **Client account not exported** (Pattern 3 difference)

   ```bash
   # Pattern 3 uses account_3of3.yaml which does NOT configure client as multisig
   # Only deposit, payment, stored are multisig accounts
   # Verify exported files (client should NOT be exported as multisig)
   ls -la data/address/bch/
   ```

3. **Address format mismatch in CSV extraction**

   ```bash
   # Check field extraction (Pattern 3 uses field 4 for P2SH via awk)
   grep "awk" scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh
   
   # Verify CSV structure
   cat data/address/bch/address_payment_*.csv | head -5
   ```

4. **Insufficient block generation**

   ```bash
   # Check block count
   docker exec bch-watch bitcoin-cli -regtest getblockcount
   # Should be 101 or more
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

#### Signing Failed Error (3-of-3)

**Cause**: Private key not imported, incorrect multisig setup, or missing signature

**Steps**:

```bash
# 1. Check keygen wallet has keys
docker exec bch-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  getwalletinfo

# 2. Check sign1 wallet has keys
docker exec bch-sign1 bitcoin-cli -regtest -rpcwallet=sign1 \
  getwalletinfo

# 3. Check sign2 wallet has keys
docker exec bch-sign2 bitcoin-cli -regtest -rpcwallet=sign2 \
  getwalletinfo

# Pattern 3 requires ALL 3 signatures - verify all wallets are ready
```

### Incomplete Signature Errors

#### "incomplete signatures" or Transaction Rejected

**Cause**: Pattern 3 requires ALL 3 signatures (keygen + sign1 + sign2)

**Steps**:

```bash
# Verify all 3 signing steps were completed in order:
# 1. keygen → tx_signed1
# 2. sign1 → tx_signed2
# 3. sign2 → tx_signed3 (final)

# Check the final signed transaction file
ls -la data/tx/bch/

# The tx_signed3 file should be used for broadcast
```

**Important**: Unlike Pattern 2, Pattern 3 MUST have sign2's signature.

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
# 1. Verify fullpubkeys were imported into keygen (using abstraction function)
db_query "keygen" "SELECT * FROM auth_fullpubkey LIMIT 5"

# 2. Re-create multisig addresses (only for deposit, payment, stored)
keygen -c config/wallet/bch/keygen.yaml --coin bch create multisig --account deposit
keygen -c config/wallet/bch/keygen.yaml --coin bch create multisig --account payment
keygen -c config/wallet/bch/keygen.yaml --coin bch create multisig --account stored
```

**Note**: `client` account is NOT configured as multisig in `account_3of3.yaml`.

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

# Verify field extraction (Pattern 3 uses awk instead of cut)
# awk -F, '!/^#/ {print $4; exit}'
```

### Database Errors

#### "duplicate key" Error

**Solution**:

```bash
# Full reset
./scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh --reset
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
# Check DB addresses (should be P2SH `2...` format)
db_query "watch" "SELECT wallet_address, account FROM address WHERE coin='bch' LIMIT 10"

# Check payment requests
db_query "watch" "SELECT * FROM payment_request WHERE coin='bch'"

# Check fullpubkeys in keygen DB
db_query "keygen" "SELECT * FROM auth_fullpubkey LIMIT 5"
```

### Log Check

```bash
# Run in verbose mode
./scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh --verbose --reset

# Check container logs
docker logs bch-watch --tail 50
docker logs bch-keygen --tail 50
docker logs bch-sign1 --tail 50
docker logs bch-sign2 --tail 50
```

## Identifying Files to Fix

| Error Type | Target File |
|------------|-------------|
| Script logic | `scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh` |
| BCH common functions | `scripts/operation/bch/bch_common.sh` |
| Common functions | `scripts/operation/common.sh` |
| BCH API (Go) | `internal/infrastructure/api/btc/bch/` |
| Multisig creation | `internal/application/usecase/keygen/btc/create_multisig*.go` |
| Watch wallet | `internal/application/usecase/watch/btc/` |
| Config loading | `pkg/config/` |
| Account config | `config/wallet/account/account_3of3.yaml` |

## Related Code (Go)

| Path | Role |
|------|------|
| `internal/infrastructure/api/btc/bch/` | **BCH API implementation (overrides BTC)** |
| `internal/infrastructure/api/btc/btc/` | BTC API (inherited by BCH) |
| `internal/application/usecase/keygen/btc/create_multisig*.go` | **Multisig address creation** |
| `internal/application/usecase/keygen/btc/` | Key generation use case |
| `internal/application/usecase/watch/btc/` | Watch wallet use case |
| `internal/application/usecase/sign/btc/` | **Sign use case (3 signatures)** |

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
- Do not break BCH Pattern 2 (`e2e-p2-p2sh-2of3.sh`)
- When modifying `bch_common.sh`, verify impact on other BCH patterns
- When modifying `common.sh`, verify impact on both BTC and BCH patterns

### BCH-Specific Considerations

- BCH uses **Raw TX Hex**, not PSBT
- BCH requires **SIGHASH_FORKID** for replay protection
- BCH regtest uses **P2SH address format** (`2...`) for multisig
- BCH 3-of-3 needs **ALL 3 signatures** (keygen + sign1 + sign2)
- **No redundancy** - if any key is lost, funds are unrecoverable
- BCH node API responses may differ from BTC (override in BCH layer)

### Pattern 3 Account Configuration

```yaml
# account_3of3.yaml - only these accounts are configured as 3-of-3 multisig:
# - deposit: 3-of-3
# - payment: 3-of-3
# - stored: 3-of-3
# 
# client account is NOT multisig in Pattern 3
```

> **Note**: For build rules, verification commands, security, see @.claude/rules/bch/e2e-script.md

## Cleanup

```bash
# Stop containers only
./scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh --cleanup

# Full reset (including data)
./scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh --reset
```

## Known Issues

As of 2026-01-17, BCH E2E Pattern 3 has the following known issues:

- **Not Verified**: Script exists but has not been fully tested
- **Likely has address format issues** similar to Pattern 1 and 2 (field extraction)
- **Account config difference**: Uses `account_3of3.yaml` which excludes client from multisig

See `scripts/operation/bch/README.md` for current verification status.
