# Fix BTC E2E Pattern 6 #{issue_number}

Fix errors in BTC E2E test (Pattern 6: P2WSH 2-of-3 Multisig).

## Prerequisites

**MUST read:** @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)

## Pattern 6 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 6 - P2WSH 2-of-3 Multisig |
| **address_type** | `bech32` |
| **Address Format** | `bcrt1q...` (62 chars, regtest) |
| **Descriptor** | `wsh(sortedmulti(2,[fp/84'/0'/0']xpub1/0/*,...))` |
| **Signature** | 2-of-3 |
| **Wallets** | watch, keygen, sign1, sign2 |

## Execution

```bash
make btc-e2e-reset P=6
make btc-e2e-verbose P=6
make btc-e2e-cleanup P=6
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh` |
| Account Config | `config/wallet/account/account_2of3.yaml` |

## Signing Flow (2-of-3)

```
Watch → Keygen (1st) → Sign1 (2nd) → Watch (broadcast)
```

## P2WSH Address Length

- P2WSH addresses are 62 characters (longer than P2WPKH's 42)
- Contains witness script hash

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/btc-e2e-p6`
- **Commit**: `fix(btc): ...`
