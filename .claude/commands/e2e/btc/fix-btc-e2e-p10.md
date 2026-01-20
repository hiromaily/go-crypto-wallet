# Fix BTC E2E Pattern 10 #{issue_number}

Fix errors in BTC E2E test (Pattern 10: P2TR MuSig2 N-of-N).

## Prerequisites

**MUST read:**

- @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)
- @.claude/skills/btc-terminology/SKILL.md (**CRITICAL**: taproot vs bech32m)
- @docs/crypto/btc/musig2/user-guide.md (MuSig2 protocol)

## Pattern 10 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 10 - P2TR MuSig2 N-of-N |
| **address_type** | `taproot` ⚠️ NOT `bech32m` |
| **Address Format** | `bcrt1p...` (62 chars, regtest) |
| **Descriptor** | `tr(musig(xpub1,xpub2,xpub3)/0/*)` |
| **Signature** | N-of-N MuSig2 (Schnorr aggregated) |
| **Wallets** | watch, keygen, sign1, sign2 |

## Execution

```bash
make btc-e2e-reset P=10
make btc-e2e-verbose P=10
make btc-e2e-cleanup P=10
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh` |
| Account Config | `config/wallet/account/account_3of3.yaml` |

## MuSig2 Signing Flow

```
Watch (create unsigned tx)
    ↓
All signers exchange nonces (Round 1)
    ↓
All signers create partial signatures (Round 2)
    ↓
Aggregate into single Schnorr signature
    ↓
Watch (broadcast)
```

## MuSig2-Specific Errors

### Nonce Exchange Failed

Check that all signers completed nonce generation.

### Partial Signature Aggregation Failed

Verify all partial signatures are valid before aggregation.

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/btc-e2e-p10`
- **Commit**: `fix(btc): ...`
