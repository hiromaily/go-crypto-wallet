# Fix BCH E2E Pattern 3 #{issue_number}

Fix errors in BCH E2E test (Pattern 3: P2SH 3-of-3 Multisig).

## Prerequisites

**MUST read:**

- @docs/task-contexts/chains/bch.md (**CRITICAL**: BCH vs BTC differences)
- @.claude/rules/bch/e2e-script.md (common rules, errors, debug commands)

> **WARNING**: BCH does NOT support SegWit, Taproot, Descriptor, PSBT, MuSig2.

## Pattern 3 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 3 - P2SH 3-of-3 Multisig |
| **Key Type** | BIP44 |
| **Address Format** | `2...` (regtest P2SH) |
| **Transaction** | Raw TX Hex (**NOT PSBT**) |
| **Signature** | 3-of-3 (ALL required) |
| **Wallets** | watch, keygen, sign1, sign2 |

## Execution

```bash
make bch-e2e-reset P=3
make bch-e2e-verbose P=3
make bch-e2e-cleanup P=3
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh` |
| BCH Common | `scripts/operation/bch/bch_common.sh` |
| Account Config | `config/wallet/account/account_3of3.yaml` |

## Signing Flow (3-of-3)

```
Watch → Keygen (1st) → Sign1 (2nd) → Sign2 (3rd) → Watch (broadcast)

※ ALL 3 signatures required
```

## Pattern 3 vs Pattern 2

| Item | Pattern 2 | Pattern 3 |
|------|-----------|-----------|
| Signature | 2-of-3 | 3-of-3 |
| Sign2 | Not required | Required |

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/bch-e2e-p3`
- **Commit**: `fix(bch): ...`
