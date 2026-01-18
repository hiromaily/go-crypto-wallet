# Fix BTC E2E Pattern 7 #{issue_number}

Fix errors in BTC E2E test (Pattern 7: P2WSH 3-of-3 Multisig).

## Prerequisites

**MUST read:** @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)

## Pattern 7 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 7 - P2WSH 3-of-3 Multisig |
| **address_type** | `bech32` |
| **Address Format** | `bcrt1q...` (62 chars, regtest) |
| **Descriptor** | `wsh(sortedmulti(3,[fp/84'/0'/0']xpub1/0/*,...))` |
| **Signature** | 3-of-3 (all required) |
| **Wallets** | watch, keygen, sign1, sign2 |

## Execution

```bash
make btc-e2e-reset P=7
make btc-e2e-verbose P=7
make btc-e2e-cleanup P=7
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh` |
| Account Config | `config/wallet/account/account_3of3.yaml` |

## Signing Flow (3-of-3)

```
Watch → Keygen (1st) → Sign1 (2nd) → Sign2 (3rd) → Watch (broadcast)

※ ALL 3 signatures required
```

## Pattern 7 vs Pattern 6

| Item | Pattern 6 | Pattern 7 |
|------|-----------|-----------|
| Signature | 2-of-3 | 3-of-3 |
| Sign2 | Not required | Required |

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/issue-{issue_number}-btc-e2e-p7`
- **Commit**: `fix(btc): ...`
