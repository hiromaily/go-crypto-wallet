# Fix BTC E2E Pattern 9 #{issue_number}

Fix errors in BTC E2E test (Pattern 9: P2TR Taproot Single-sig).

## Prerequisites

**MUST read:**

- @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)
- @.claude/skills/btc-terminology/SKILL.md (**CRITICAL**: taproot vs bech32m)

## Pattern 9 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 9 - P2TR Taproot Single-sig |
| **address_type** | `taproot` ⚠️ NOT `bech32m` |
| **Address Format** | `bcrt1p...` (62 chars, regtest) |
| **Descriptor** | `tr([fp/86'/0'/0']xpub/0/*)` |
| **Signature** | Single-sig (Schnorr) |
| **Wallets** | watch, keygen |

## Execution

```bash
make btc-e2e-reset P=9
make btc-e2e-verbose P=9
make btc-e2e-cleanup P=9
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh` |
| Account Config | `config/wallet/account/account.yaml` |

## ⚠️ Critical: address_type Configuration

```bash
# CORRECT
export WALLET_ADDRESS_TYPE="taproot"

# WRONG - bech32m is encoding format, not address type
export WALLET_ADDRESS_TYPE="bech32m"  # ❌
```

## Taproot-Specific Errors

### Invalid address_type

```bash
echo $WALLET_ADDRESS_TYPE  # Must be "taproot"
```

### Descriptor Format

```bash
# Should be tr(...) format
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
```

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/issue-{issue_number}-btc-e2e-p9`
- **Commit**: `fix(btc): ...`
