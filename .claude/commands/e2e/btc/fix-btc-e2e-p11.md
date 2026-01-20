# Fix BTC E2E Pattern 11 #{issue_number}

Fix errors in BTC E2E test (Pattern 11: P2TR Tapscript M-of-N).

## Prerequisites

**MUST read:**

- @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)
- @.claude/skills/btc-terminology/SKILL.md (**CRITICAL**: taproot vs bech32m)
- @docs/crypto/btc/taproot/user-guide.md (Tapscript)

## Pattern 11 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 11 - P2TR Tapscript M-of-N |
| **address_type** | `taproot` ⚠️ NOT `bech32m` |
| **Address Format** | `bcrt1p...` (62 chars, regtest) |
| **Script** | Tapscript with script path spending |
| **Signature** | M-of-N via Tapscript |
| **Wallets** | watch, keygen, sign1, sign2 |

## Execution

```bash
make btc-e2e-reset P=11
make btc-e2e-verbose P=11
make btc-e2e-cleanup P=11
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh` |
| Account Config | `config/wallet/account/account_2of3.yaml` |

## Pattern 11 vs Pattern 10

| Item | Pattern 10 (MuSig2) | Pattern 11 (Tapscript) |
|------|---------------------|------------------------|
| Spending Path | Key path | Script path |
| Signature | Single aggregated | Multiple individual |
| Privacy | Higher (looks like single-sig) | Lower (reveals script) |
| Flexibility | N-of-N only | M-of-N possible |

## Tapscript-Specific Errors

### Script Path Verification Failed

Check that the Tapscript tree is correctly constructed.

### Control Block Invalid

Verify the control block includes correct internal key and merkle proof.

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/btc-e2e-p11`
- **Commit**: `fix(btc): ...`
