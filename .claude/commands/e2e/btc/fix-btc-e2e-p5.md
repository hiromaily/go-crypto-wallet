# Fix BTC E2E Pattern 5 #{issue_number}

Fix errors in BTC E2E test (Pattern 5: P2WPKH Native SegWit Single-sig).

## Prerequisites

**MUST read:** @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)

## Pattern 5 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 5 - P2WPKH Native SegWit Single-sig |
| **address_type** | `bech32` |
| **Address Format** | `bcrt1q...` (42 chars, regtest) |
| **Descriptor** | `wpkh([fp/84'/0'/0']xpub/0/*)` |
| **Signature** | Single-sig (1) |
| **Wallets** | watch, keygen |

## Execution

```bash
make btc-e2e-reset P=5
make btc-e2e-verbose P=5
make btc-e2e-cleanup P=5
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh` |
| Account Config | `config/wallet/account/account.yaml` |

## Pattern 5 Address Characteristics

- Native SegWit (no P2SH wrapper)
- 42 characters total
- Lower fees than P2SH-wrapped

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/btc-e2e-p5-{issue_number}`
- **Commit**: `fix(btc): ...`
