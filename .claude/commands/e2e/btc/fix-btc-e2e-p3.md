# Fix BTC E2E Pattern 3 #{issue_number}

Fix errors in BTC E2E test (Pattern 3: P2SH-P2WPKH Single-sig).

## Prerequisites

**MUST read:** @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)

## Pattern 3 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 3 - P2SH-P2WPKH Single-sig |
| **address_type** | `p2sh-segwit` |
| **Address Format** | `2...` (regtest P2SH) |
| **Descriptor** | `sh(wpkh([fp/49'/0'/0']xpub/0/*))` |
| **Signature** | Single-sig (1) |
| **Wallets** | watch, keygen |

## Execution

```bash
make btc-e2e-reset P=3
make btc-e2e-verbose P=3
make btc-e2e-cleanup P=3
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh` |
| Account Config | `config/wallet/account/account.yaml` |

## Pattern 3 vs Pattern 1

| Item | Pattern 1 (P2PKH) | Pattern 3 (P2SH-P2WPKH) |
|------|-------------------|-------------------------|
| address_type | `legacy` | `p2sh-segwit` |
| key_type | BIP44 | BIP49 |
| Witness | None | SegWit wrapped |

## Pattern-Specific Errors

### Descriptor Format Mismatch

```bash
# Check descriptor (should be sh(wpkh(...)))
jq '.[0].desc' data/descriptor/btc/payment_descriptors.json
```

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/btc-e2e-p3-{issue_number}`
- **Commit**: `fix(btc): ...`
