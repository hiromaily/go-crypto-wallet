# Fix BTC E2E Pattern 1 #{issue_number}

Fix errors in BTC E2E test (Pattern 1: P2PKH Single-sig).

## Prerequisites

**MUST read:** @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)

## Pattern 1 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 1 - P2PKH Single-sig |
| **address_type** | `legacy` |
| **Address Format** | `m.../n...` (regtest P2PKH) |
| **Descriptor** | `pkh([fp/44'/0'/0']xpub/0/*)` |
| **Signature** | Single-sig (1) |
| **Wallets** | watch, keygen |

## Execution

```bash
# Run E2E test (build is automatic)
make btc-e2e-reset P=1

# Verbose mode
make btc-e2e-verbose P=1

# Cleanup
make btc-e2e-cleanup P=1
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh` |
| Account Config | `config/wallet/account/account.yaml` |

## Pattern 1 vs Other Patterns

| Item | Pattern 1 | Pattern 2 (2-of-3) |
|------|-----------|---------------------|
| Signature | Single-sig | 2-of-3 Multisig |
| Wallets | keygen only | keygen + sign1 + sign2 |
| Address | `m.../n...` (P2PKH) | `2...` (P2SH) |

## Pattern-Specific Errors

### Descriptor Import Failed

```bash
# Check descriptor format (should be pkh(...))
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
```

**Expected**: `pkh([fingerprint/44'/0'/0']xpub.../0/*)`

### address_type Mismatch

```bash
# Verify environment variable
echo $WALLET_ADDRESS_TYPE  # Should be "legacy"
```

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/issue-{issue_number}-btc-e2e-p1`
- **Commit**: `fix(btc): ...`
- **Skill**: @.claude/skills/git-workflow/SKILL.md
