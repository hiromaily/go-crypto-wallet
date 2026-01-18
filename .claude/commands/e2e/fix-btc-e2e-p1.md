# Fix BTC E2E Pattern 1 Errors #{issue_number}

Fix errors in BTC E2E test (Pattern 1: P2PKH Single-sig).

## Prerequisites

**Read the following common rules first:**

- @.claude/rules/btc/e2e-script.md - BTC E2E common rules (build, verification, escalation, security)

## Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `{issue_number}` | Optional | GitHub issue number. Follow git-workflow when specified |

## Overview

This command diagnoses and fixes errors when running `make btc-e2e-reset P=1`.
The script already exists, so focus on identifying and fixing the root cause.

### Pattern 1 Technical Specifications

| Item | Value |
|------|-------|
| **Pattern Number** | 1 |
| **Network** | **regtest** (local environment) |
| **Key Type** | P2PKH (BIP44 Legacy) |
| **Script Type** | Single-sig |
| **Address Format** | `m.../n...` (regtest/testnet P2PKH) |
| **Signature Requirement** | Single-sig (1 signature) |
| **Descriptor** | `pkh([fingerprint/44'/0'/0']xpub.../0/*)` |
| **Required Wallets** | watch, keygen |
| **Environment Variable** | `WALLET_ADDRESS_TYPE="legacy"` |

### Differences from Pattern 2 (2-of-3 Multisig)

| Item | Pattern 1 | Pattern 2 |
|------|-----------|-----------|
| Signature Requirement | Single-sig (1) | 2-of-3 Multisig |
| Required Wallets | keygen only | keygen + sign1 + sign2 |
| Address Format | `m.../n...` (P2PKH) | `2...` (P2SH) |
| fullpubkey Exchange | Not required | Required |
| Account Config | `account_singlesig.yaml` | `account_2of3.yaml` |

### When issue number is specified

Load `git-workflow` skill and work with these settings:

- **Branch name**: `fix/issue-{issue_number}-btc-e2e-p1`
- **Commit type**: `fix(btc)`
- **Scope**: BTC E2E Pattern 1

→ See @.claude/skills/git-workflow/SKILL.md for details

### When issue number is not specified

Fix locally without creating branch or PR.

## Pattern 1 Specific Documentation

In addition to Required Documentation in common rules, refer to:

- @scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh - **Target script**
- @scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh - Pattern 2 script (Multisig reference)

## Pre-check: Environment Variables

**Pattern 1 requires `WALLET_ADDRESS_TYPE="legacy"`.**

Auto-configured in script, but for verification:

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "legacy"
```

> **Note**: Do not edit config files directly. Override with environment variables.
> See "Configuration File Policy" in common rules for details.

## Error Diagnosis Steps

### Step 1: Reproduce Error

```bash
# Full reset and run E2E test
make btc-e2e-reset P=1
```

Check error message and categorize below.

### Step 2: Identify Error Category

| Error Message | Category | Reference Section |
|---------------|----------|-------------------|
| `No utxo` | UTXO-related | [UTXO Errors](#utxo-errors) |
| `connection refused` | Infrastructure | See common rules |
| `wallet not found` | Wallet | [Wallet Errors](#wallet-errors) |
| `signing failed` | Signing | [Signing Errors](#signing-errors) |
| `descriptor` | Descriptor | [Descriptor Errors](#descriptor-errors) |
| `address_type` mismatch | Config | [Config Errors](#config-errors) |
| `duplicate key` | DB | See common rules |

## Pattern 1 Specific Errors and Solutions

For common errors (connection refused, duplicate key, etc.), see common rules.

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

1. **Descriptor not imported correctly**

   ```bash
   # Debug: Check address info
   docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch \
     getaddressinfo "<payment_address>"
   ```

   Check:
   - `solvable: true`
   - `ismine: true` (false is OK for watch-only wallet)

2. **Insufficient block generation**

   ```bash
   # Check block count
   docker exec btc-watch bitcoin-cli -regtest getblockcount
   # Should be 101 or more
   ```

3. **address_type mismatch**

   ```bash
   # Check environment variable
   echo $WALLET_ADDRESS_TYPE  # Should be "legacy"
   ```

### Wallet Errors

#### "wallet not found" Error

**Solution**:

```bash
# List wallets
docker exec btc-watch bitcoin-cli -regtest listwallets

# Create wallets
docker exec btc-watch bitcoin-cli -regtest createwallet "watch" true true
docker exec btc-keygen bitcoin-cli -regtest createwallet "keygen" false true
```

### Signing Errors

#### Signing Failed Error

**Cause**: Private key not imported or address_type mismatch

**Steps**:

```bash
# 1. Check environment variable
echo "WALLET_ADDRESS_TYPE: $WALLET_ADDRESS_TYPE"  # Should be "legacy"

# 2. Check keygen wallet private keys
docker exec btc-keygen bitcoin-cli -regtest -rpcwallet=keygen \
  listdescriptors true
```

### Descriptor Errors

#### Descriptor Import Failed

**Steps**:

```bash
# Check descriptor file
cat data/descriptor/btc/payment_descriptors.json

# Check format (P2PKH should be pkh(...) format)
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
```

**Expected format** (Pattern 1):

```
pkh([fingerprint/44'/0'/0']xpub.../0/*)
```

### Config Errors

#### address_type Mismatch Error

**Symptoms**: Generated address differs from expected

**Steps**:

```bash
# Check environment variable export in script
grep -A2 "Environment Variable Overrides" \
  scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh
```

**Fix**: Ensure the following is set in script

```bash
export WALLET_ADDRESS_TYPE="legacy"
```

## Debug Commands

### Status Check

```bash
# Bitcoin node status
docker exec btc-watch bitcoin-cli -regtest getblockchaininfo

# Wallet balance
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch getbalances

# UTXO list
docker exec btc-watch bitcoin-cli -regtest -rpcwallet=watch listunspent
```

### Database Queries

> **Note**: Database queries depend on `DB_TYPE` setting. See "Database Debug Commands" in common rules.

```bash
# Check DB addresses (using abstraction function)
db_query "watch" "SELECT wallet_address, account FROM address WHERE coin='btc' LIMIT 10"

# Check payment requests
db_query "watch" "SELECT * FROM payment_request WHERE coin='btc'"
```

### Log Check

```bash
# Run in verbose mode (build is automatic)
make btc-e2e-verbose P=1
```

## Identifying Files to Fix

| Error Type | Target File |
|------------|-------------|
| Script logic | `scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh` |
| Common functions | `scripts/operation/common.sh` |
| Descriptor generation | `internal/application/usecase/keygen/btc/` |
| Watch wallet | `internal/application/usecase/watch/btc/` |
| Bitcoin RPC | `internal/infrastructure/wallet/api/btc/` |
| Config loading | `pkg/config/` |

## Related Code (Go)

| Path | Role |
|------|------|
| `internal/application/usecase/keygen/btc/` | Key generation use case |
| `internal/application/usecase/watch/btc/` | Watch wallet use case |
| `internal/infrastructure/wallet/api/btc/` | Bitcoin RPC implementation |
| `internal/domain/wallet/key/` | Key domain model |
| `pkg/config/loader.go` | Config loader |

## Cautions

### Avoid Impact on Existing Scripts

- Do not break Pattern 2 (`e2e-p2-p2pkh-2of3.sh`)
- Do not break Pattern 8 (`e2e-p8-p2sh-p2wsh-3of3.sh`)
- When modifying `common.sh`, verify impact on other patterns
- Set environment variables locally within each script

> **Note**: For build rules, verification commands, security, see common rules.

## Cleanup

```bash
# Stop containers only
make btc-e2e-cleanup P=1

# Full reset (including data)
make btc-e2e-reset P=1
```
