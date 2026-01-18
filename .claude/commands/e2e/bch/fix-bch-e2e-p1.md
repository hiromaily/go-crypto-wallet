# Fix BCH E2E Pattern 1 #{issue_number}

Fix errors in BCH E2E test (Pattern 1: P2PKH Single-sig).

## Prerequisites

**MUST read:**

- @docs/task-contexts/chains/bch.md (**CRITICAL**: BCH vs BTC differences)
- @.claude/rules/bch/e2e-script.md (common rules, errors, debug commands)

> **WARNING**: BCH does NOT support SegWit, Taproot, Descriptor, PSBT, MuSig2.

## Pattern 1 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 1 - P2PKH Single-sig |
| **Key Type** | BIP44 |
| **Address Format** | `m.../n...` (regtest P2PKH) |
| **Transaction** | Raw TX Hex (**NOT PSBT**) |
| **Signature** | Single-sig (ECDSA + SIGHASH_FORKID) |
| **Wallets** | watch, keygen |

## Execution

```bash
make bch-e2e-reset P=1
make bch-e2e-verbose P=1
make bch-e2e-cleanup P=1
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh` |
| BCH Common | `scripts/operation/bch/bch_common.sh` |
| Account Config | `config/wallet/account/account.yaml` |

## BCH vs BTC Pattern 1

| Item | BTC Pattern 1 | BCH Pattern 1 |
|------|---------------|---------------|
| Descriptor | `pkh(...)` | **NOT supported** |
| Transaction | PSBT | **Raw TX Hex** |
| Address Import | Descriptor | **Address export/import** |
| Container | `btc-watch` | **`bch-watch`** |

## BCH-Specific Errors

### UTXO Query Issue

```bash
# Check balance (watch-only compatible)
docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch \
  getbalance "*" 1 true
```

### Address Format

BCH regtest uses legacy format (`m.../n...`), not CashAddr.

```bash
# Check exported addresses
cat data/address/bch/address_payment_*.csv | head -5
```

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/issue-{issue_number}-bch-e2e-p1`
- **Commit**: `fix(bch): ...`
- **Skill**: @.claude/skills/git-workflow/SKILL.md
