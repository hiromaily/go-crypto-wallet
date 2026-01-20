# Fix BTC E2E Pattern 4 #{issue_number}

Fix errors in BTC E2E test (Pattern 4: P2SH-P2WSH 2-of-3 Multisig).

## Prerequisites

**MUST read:** @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)

## Pattern 4 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 4 - P2SH-P2WSH 2-of-3 Multisig |
| **address_type** | `p2sh-segwit` |
| **Address Format** | `2...` (regtest P2SH) |
| **Descriptor** | `sh(wsh(sortedmulti(2,[fp/49'/0'/0']xpub1/0/*,...)))` |
| **Signature** | 2-of-3 |
| **Wallets** | watch, keygen, sign1, sign2 |

## Execution

```bash
make btc-e2e-reset P=4
make btc-e2e-verbose P=4
make btc-e2e-cleanup P=4
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh` |
| Account Config | `config/wallet/account/account_2of3.yaml` |

## Signing Flow (2-of-3)

```
Watch → Keygen (1st) → Sign1 (2nd) → Watch (broadcast)
```

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/btc-e2e-p4`
- **Commit**: `fix(btc): ...`
