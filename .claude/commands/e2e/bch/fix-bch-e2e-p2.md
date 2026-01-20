# Fix BCH E2E Pattern 2 #{issue_number}

Fix errors in BCH E2E test (Pattern 2: P2SH 2-of-3 Multisig).

## Prerequisites

**MUST read:**

- @docs/task-contexts/chains/bch.md (**CRITICAL**: BCH vs BTC differences)
- @.claude/rules/bch/e2e-script.md (common rules, errors, debug commands)

> **WARNING**: BCH does NOT support SegWit, Taproot, Descriptor, PSBT, MuSig2.

## Pattern 2 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 2 - P2SH 2-of-3 Multisig |
| **Key Type** | BIP44 |
| **Address Format** | `2...` (regtest P2SH) |
| **Transaction** | Raw TX Hex (**NOT PSBT**) |
| **Signature** | 2-of-3 (ECDSA + SIGHASH_FORKID) |
| **Wallets** | watch, keygen, sign1, sign2 |

## Execution

```bash
make bch-e2e-reset P=2
make bch-e2e-verbose P=2
make bch-e2e-cleanup P=2
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh` |
| BCH Common | `scripts/operation/bch/bch_common.sh` |
| Account Config | `config/wallet/account/account_2of3.yaml` |

## Signing Flow (2-of-3)

```
Watch (create unsigned tx)
    ↓
Keygen (1st signature)
    ↓
Sign1 (2nd signature) ← Complete here
    ↓
Watch (broadcast)
```

## BCH-Specific Considerations

- SIGHASH_FORKID required for replay protection
- Raw TX Hex format (not PSBT)
- No descriptor-based address import

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/bch-e2e-p2`
- **Commit**: `fix(bch): ...`
