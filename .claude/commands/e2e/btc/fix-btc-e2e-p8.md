# Fix BTC E2E Pattern 8 #{issue_number}

Fix errors in BTC E2E test (Pattern 8: P2SH-P2WSH 3-of-3 Multisig).

## Prerequisites

**MUST read:** @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)

## Pattern 8 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 8 - P2SH-P2WSH 3-of-3 Multisig |
| **address_type** | `p2sh-segwit` |
| **Address Format** | `2...` (regtest P2SH) |
| **Descriptor** | `sh(wsh(sortedmulti(3,[fp/49'/0'/0']xpub1/0/*,...)))` |
| **Signature** | 3-of-3 (all required) |
| **Wallets** | watch, keygen, sign1, sign2 |

## Execution

```bash
make btc-e2e-reset P=8
make btc-e2e-verbose P=8
make btc-e2e-cleanup P=8
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh` |
| Account Config | `config/wallet/account/account_3of3.yaml` |

## Signing Flow (3-of-3)

```
Watch → Keygen (1st) → Sign1 (2nd) → Sign2 (3rd) → Watch (broadcast)
```

## Pattern 8 vs Pattern 2

| Item | Pattern 2 | Pattern 8 |
|------|-----------|-----------|
| Key Type | BIP44 (Legacy) | BIP49 (P2SH-SegWit) |
| Signature | 2-of-3 | 3-of-3 |
| Descriptor | `sh(multi(2,...))` | `sh(wsh(sortedmulti(3,...)))` |

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/issue-{issue_number}-btc-e2e-p8`
- **Commit**: `fix(btc): ...`
