# Fix BTC E2E Pattern 2 #{issue_number}

Fix errors in BTC E2E test (Pattern 2: P2PKH 2-of-3 Multisig).

## Prerequisites

**MUST read:** @.claude/rules/btc/e2e-script.md (common rules, errors, debug commands)

## Pattern 2 Specifications

| Item | Value |
|------|-------|
| **Pattern** | 2 - P2PKH 2-of-3 Multisig |
| **address_type** | `legacy` |
| **Address Format** | `2...` (regtest P2SH) |
| **Descriptor** | `sh(multi(2,[fp/44'/0'/0']xpub1/0/*,xpub2/0/*,xpub3/0/*))` |
| **Signature** | 2-of-3 (any 2 signatures) |
| **Wallets** | watch, keygen, sign1, sign2 |

## Execution

```bash
# Run E2E test (build is automatic)
make btc-e2e-reset P=2

# Verbose mode
make btc-e2e-verbose P=2

# Cleanup
make btc-e2e-cleanup P=2
```

## Pattern-Specific Files

| Type | File |
|------|------|
| Script | `scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh` |
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

※ Sign2 not required - 2 signatures satisfy 2-of-3
```

## Pattern-Specific Errors

### Descriptor Format Error

Using BIP49 format instead of BIP44.

```bash
echo $WALLET_ADDRESS_TYPE  # Should be "legacy"
```

### fullpubkey Import Error

```bash
# Check fullpubkey files exist
ls data/fullpubkey/btc/
```

### Insufficient/Too Many Signatures

```bash
# Check PSBT signature status
btc_cli "btc-watch" analyzepsbt "${psbt_hex}"
```

## Git Workflow

When `{issue_number}` is specified:

- **Branch**: `fix/btc-e2e-p2-{issue_number}`
- **Commit**: `fix(btc): ...`
- **Skill**: @.claude/skills/git-workflow/SKILL.md
