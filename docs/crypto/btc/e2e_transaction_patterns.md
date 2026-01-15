# E2E Transaction Patterns Guide

This document explains transaction combination patterns for Bitcoin/Bitcoin Cash. Various E2E workflow patterns exist depending on key types and whether multisig is used.

## Table of Contents

1. [Overview](#overview)
2. [Supported Key Types](#supported-key-types)
3. [Signature Patterns](#signature-patterns)
4. [E2E Workflow Matrix](#e2e-workflow-matrix)
5. [Details of Each Pattern](#details-of-each-pattern)
6. [Account Types and Signing Requirements](#account-types-and-signing-requirements)
7. [Implementation Status](#implementation-status)
8. [E2E Script Reference](#e2e-script-reference)

---

## Overview

Bitcoin transactions can be classified along two main axes:

1. **Key Type (Address Type)** - Which BIP standard is used to generate addresses
2. **Signature Pattern** - Single-sig or multisig

The combination of these creates various E2E workflows.

---

## Supported Key Types

### Bitcoin (BTC)

| Address Type | BIP | Prefix (Mainnet) | Prefix (Testnet) | Description |
|--------------|-----|------------------|------------------|-------------|
| **P2PKH** (Legacy) | BIP44 | `1...` | `m.../n...` | Traditional Pay-to-Public-Key-Hash |
| **P2SH-P2WPKH** | BIP49 | `3...` | `2...` | SegWit wrapped in P2SH |
| **P2WPKH** (Native SegWit) | BIP84 | `bc1q...` | `tb1q...` | Native SegWit |
| **P2TR** (Taproot) | BIP86 | `bc1p...` | `tb1p...` | Taproot (recommended) |

### Bitcoin Cash (BCH)

| Address Type | Prefix | Description |
|--------------|--------|-------------|
| **CashAddr** | `bitcoincash:q...` | Bitcoin Cash dedicated format |
| **Legacy** | `1...` | Legacy format (for compatibility) |

### Key Derivation Paths

| Standard | Path | Usage |
|----------|------|-------|
| BIP44 | `m/44'/0'/account'/change/index` | P2PKH (Legacy) |
| BIP49 | `m/49'/0'/account'/change/index` | P2SH-P2WPKH |
| BIP84 | `m/84'/0'/account'/change/index` | P2WPKH (Native SegWit) |
| BIP86 | `m/86'/0'/account'/change/index` | P2TR (Taproot) |

---

## Signature Patterns

### Single-Sig

Pattern where a single private key is used for signing.

```
┌─────────────────────────────────────────────────────────┐
│                  SINGLE-SIG FLOW                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  2. Keygen Wallet: Sign with single key                │
│          ↓                                              │
│  3. Watch Wallet: Broadcast transaction                 │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**Characteristics:**

- Simple and fast
- Completed with a single signature
- Risk concentrated since there's only one private key

### Multi-Sig

Pattern where multiple private keys are used for signing. M-of-N (M signatures required out of N keys).

#### 3-of-3 Multisig

```
┌─────────────────────────────────────────────────────────┐
│                  3-of-3 MULTISIG FLOW                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  2. Keygen Wallet: Sign (1st signature)                │
│          ↓                                              │
│  3. Sign1 Wallet: Sign (2nd signature)                 │
│          ↓                                              │
│  4. Sign2 Wallet: Sign (3rd signature)                 │
│          ↓                                              │
│  5. Watch Wallet: Broadcast transaction                 │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

#### 2-of-3 Multisig

```
┌─────────────────────────────────────────────────────────┐
│                  2-of-3 MULTISIG FLOW                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Watch Wallet: Create unsigned transaction           │
│          ↓                                              │
│  2. Keygen Wallet: Sign (1st signature)                │
│          ↓                                              │
│  3. Sign1 Wallet: Sign (2nd signature)                 │
│          ↓                                              │
│  4. Watch Wallet: Broadcast transaction                 │
│                                                         │
│  (Sign2 Wallet not required - completed with 2 sigs)   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### MuSig2 (Signature Aggregation)

Aggregate signature protocol based on Schnorr signatures. N-of-N multisig becomes the same size as single-sig.

```
┌─────────────────────────────────────────────────────────┐
│                    MUSIG2 FLOW                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Round 1: Nonce Generation (can be parallelized)       │
│  ├─ Keygen Wallet: Generate nonce                       │
│  ├─ Sign1 Wallet: Generate nonce                        │
│  └─ Sign2 Wallet: Generate nonce                        │
│          ↓                                              │
│  Round 2: Signing (sequential)                         │
│  ├─ Keygen Wallet: Create partial signature             │
│  ├─ Sign1 Wallet: Create partial signature              │
│  └─ Sign2 Wallet: Create partial signature              │
│          ↓                                              │
│  Aggregation:                                           │
│  └─ Watch Wallet: Aggregate & broadcast                 │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**MuSig2 Benefits:**

- Transaction size reduced by 30-50%
- Improved privacy (indistinguishable from single-sig)
- Reduced fees

---

## E2E Workflow Matrix

### BTC Pattern Matrix

| Pattern | Key Type | Signature Pattern | Address Format | E2E Script Support |
|---------|----------|-------------------|----------------|-------------------|
| **1** | **P2PKH (BIP44)** | **Single-sig** | **`1...`** | **✅ e2e/e2e-p2pkh-singlesig.sh** |
| **2** | **P2PKH (BIP44)** | **2-of-3 Multisig** | **`3...` (P2SH wrapped)** | **✅ e2e/e2e-p2pkh-2of3.sh (Fixed in #357)** |
| 3 | P2SH-P2WPKH (BIP49) | Single-sig | `3...` | 🔶 Manual testing |
| 4 | P2SH-P2WPKH (BIP49) | 2-of-3 Multisig | `3...` | ❌ Not supported |
| 5 | P2WPKH (BIP84) | Single-sig | `bc1q...` | 🔶 Manual testing |
| 6 | P2WSH (BIP84) | 2-of-3 Multisig | `bc1q...` | ❌ Not supported |
| 7 | P2WSH (BIP84) | 3-of-3 Multisig | `bc1q...` | ❌ Not supported |
| **8** | **P2SH-P2WSH** | **3-of-3 Multisig** | **`3...`** | **✅ e2e/e2e-p2sh-p2wsh-3of3.sh** |
| 9 | P2TR (BIP86) | Single-sig | `bc1p...` | 🔶 Manual testing |
| 10 | P2TR (BIP86) | MuSig2 (N-of-N) | `bc1p...` | 🔜 In development |
| 11 | P2TR (BIP86) | Tapscript (M-of-N) | `bc1p...` | 🔜 In development |

### BCH Pattern Matrix

| Pattern | Key Type | Signature Pattern | Address Format | E2E Script Support |
|---------|----------|-------------------|----------------|-------------------|
| 1 | CashAddr | Single-sig | `bitcoincash:q...` | 🔶 Manual testing |
| **2** | **CashAddr** | **3-of-3 Multisig** | **`bitcoincash:p...`** | **✅ e2e-workflow.sh** |
| 3 | CashAddr | 2-of-3 Multisig | `bitcoincash:p...` | ❌ Not supported |

---

## Details of Each Pattern

### Pattern 1: BTC P2PKH Single-sig

**Currently implemented in `scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh`**

```
Address Type: P2PKH (BIP44 Legacy)
Signing Requirements: Single-sig (Keygen only)
Descriptor: pkh([fingerprint/44'/0'/0']xpub.../0/*)
```

**Workflow:**

1. Generate Seed in Keygen
2. Generate HD Key in Keygen (10 accounts each)
3. Export Descriptor from Keygen
4. Import Descriptor to Watch
5. Generate Test UTXO (regtest)
6. Create unsigned transaction → Sign once → Broadcast

**Characteristics:**

- Simple and fast (completed with single signature)
- Sign1/Sign2 wallets not required
- Uses BIP44 key derivation path
- Legacy address format (`m...`/`n...` in regtest)

### Pattern 8: BTC P2SH-P2WSH 3-of-3 Multisig (Current E2E)

**Currently implemented in `scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh`**

```
Address Type: P2SH-P2WSH (BIP49 wrapped SegWit)
Signing Requirements: 3-of-3 (Keygen + Sign1 + Sign2)
Descriptor: sh(wsh(sortedmulti(3, xpub1, xpub2, xpub3)))
```

**Workflow:**

1. Generate Seed in Keygen/Sign1/Sign2
2. Generate HD Key in Keygen (10 accounts each)
3. Generate HD Key in Sign1/Sign2
4. Export fullpubkey from Sign1/Sign2
5. Import fullpubkey to Keygen
6. Export Descriptor from Keygen
7. Import Descriptor to Watch
8. Generate Test UTXO (regtest)
9. Create unsigned transaction → Sign 3 times → Broadcast

### Pattern 2: BTC P2PKH 2-of-3 Multisig

**✅ Fully implemented in `scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh` (Fixed in #357)**

```
Address Type: P2PKH (BIP44 Legacy) with 2-of-3 Multisig
Signing Requirements: 2-of-3 (Keygen + Sign1, Sign2 is optional)
Descriptor: sh(multi(2, [fingerprint/44'/1'/1]xpub1/0/*, [fingerprint/44'/1'/1]xpub2/0/*, [fingerprint/44'/1'/1]xpub3/0/*))
Address Format: 2... (P2SH in regtest), 3... (P2SH in mainnet)
```

**Fixed Issues (PR #357):**

- ✅ P2SH multisig descriptor generation now working correctly
- ✅ Key derivation path mismatch resolved (non-hardened account paths for multisig)
- ✅ P2SH address generation (not P2WSH) working as expected
- ✅ Transaction signing and broadcasting verified end-to-end

**Root Cause of Previous Failure:**

Sign/keygen wallets were using hardened account derivation (`m/44'/1'/[account]'`) while descriptors specified non-hardened paths (`m/44'/1'/[account]`). For multisig with xpub derivation, non-hardened account paths are required since xpubs cannot derive hardened children.

**Workflow:**

1. Generate Seed in Keygen/Sign1/Sign2
2. Generate HD Key in Keygen (10 accounts each)
3. Generate HD Key in Sign1/Sign2 (with non-hardened account derivation for multisig)
4. Export fullpubkey from Sign1/Sign2
5. Import fullpubkey to Keygen
6. Export Descriptor from Keygen (generates `sh(multi(2, ...))`)
7. Import Descriptor to Watch
8. Generate Test UTXO (regtest)
9. Create unsigned transaction → Sign 2 times → Broadcast

**Characteristics:**

- 2-of-3 multisig (completed with Keygen + Sign1 signatures, Sign2 not required)
- Uses BIP44 key derivation path with non-hardened account index for multisig
- Generates proper P2SH addresses (`2...` prefix in regtest)
- Fully compatible with Bitcoin Core descriptor import

### Pattern 2: BCH CashAddr 3-of-3 Multisig (Current E2E)

**Currently implemented in `scripts/operation/bch/e2e-workflow.sh`**

```
Address Type: CashAddr P2SH
Signing Requirements: 3-of-3 (Keygen + Sign1 + Sign2)
Address Format: bitcoincash:p... (P2SH multisig)
```

**Workflow:**

1. Generate Seed in Keygen/Sign1/Sign2
2. Generate HD Key in Keygen
3. Generate HD Key in Sign1/Sign2
4. Export fullpubkey from Sign1/Sign2
5. Import fullpubkey to Keygen
6. Create Multisig address in Keygen
7. Export address from Keygen
8. Import address to Watch
9. Generate Test UTXO (regtest)
10. Create unsigned transaction → Sign 3 times → Broadcast

### Pattern 9: BTC P2TR Single-sig (Taproot)

```
Address Type: P2TR (BIP86)
Signing Requirements: Single-sig (Keygen only)
Descriptor: tr([fingerprint/86'/0'/0']xpub.../0/*)
```

**Simple Workflow:**

1. Generate Seed in Keygen
2. Generate BIP86 HD Key in Keygen
3. Export Taproot address from Keygen
4. Import Taproot address to Watch
5. Create unsigned transaction → Sign once (Schnorr) → Broadcast

### Pattern 10: BTC P2TR MuSig2 (In Development)

```
Address Type: P2TR (BIP86)
Signing Requirements: N-of-N MuSig2 (all signatures required)
Descriptor: tr(musig(xpub1, xpub2, xpub3))
```

**2-Round Protocol:**

1. Round 1: Generate nonces in each wallet
2. Round 2: Create partial signatures in each wallet
3. Aggregate signatures in Watch and broadcast

---

## Account Types and Signing Requirements

| Account | Purpose | Recommended Signature Pattern | Reason |
|---------|---------|------------------------------|--------|
| **client** | Customer deposit addresses | Single-sig | Required for customer-side operations |
| **deposit** | Deposit aggregation | Multisig (2-of-3 or 3-of-3) | Enhanced security |
| **payment** | Payments | Multisig (2-of-3 or 3-of-3) | Approval workflow |
| **stored** | Long-term storage | Multisig (3-of-3) | Highest level of security |

---

## Implementation Status

### Key Type Implementation Status

| Key Type | BTC | BCH |
|----------|-----|-----|
| P2PKH (Legacy) | ✅ Implemented | N/A |
| P2SH-P2WPKH (BIP49) | ✅ Implemented | N/A |
| P2WPKH (BIP84) | ✅ Implemented | N/A |
| P2TR (BIP86) | ✅ Implemented | N/A |
| CashAddr | N/A | ✅ Implemented |

### Signature Pattern Implementation Status

| Signature Pattern | BTC | BCH |
|-------------------|-----|-----|
| Single-sig | ✅ Implemented | ✅ Implemented |
| 2-of-3 Multisig | ✅ Implemented (Fixed in #357) | ⚠️ Partial |
| 3-of-3 Multisig | ✅ Implemented | ✅ Implemented |
| MuSig2 | 🔜 In development | N/A |

### Descriptor Support

| Feature | BTC | BCH |
|---------|-----|-----|
| Descriptor Export | ✅ Implemented | ❌ Not supported |
| Descriptor Import | ✅ Implemented | ❌ Not supported |
| Bitcoin Core Integration | ✅ Implemented | N/A |

---

## E2E Script Reference

### Currently Available E2E Scripts

| Script | Coin | Pattern | Signing Requirements |
|--------|------|---------|---------------------|
| `scripts/operation/btc/e2e/e2e-p2pkh-singlesig.sh` | BTC | P2PKH Single-sig | Single-sig |
| `scripts/operation/btc/e2e/e2e-p2pkh-2of3.sh` | BTC | P2PKH 2-of-3 Multisig | 2-of-3 |
| `scripts/operation/btc/e2e/e2e-p2sh-p2wsh-3of3.sh` | BTC | P2SH-P2WSH Multisig | 3-of-3 |
| `scripts/operation/bch/e2e-workflow.sh` | BCH | CashAddr Multisig | 3-of-3 |

### Planned E2E Scripts

| Script (Planned) | Coin | Pattern | Signing Requirements | Priority |
|------------------|------|---------|---------------------|----------|
| `e2e-singlesig.sh` | BTC | P2WPKH/P2TR Single-sig | 1 | High |
| `e2e-musig2.sh` | BTC | P2TR MuSig2 | N-of-N | Medium |
| `e2e-2of3.sh` | BTC | P2WSH 2-of-3 | 2-of-3 | Low |
| `e2e-tapscript.sh` | BTC | P2TR Script Path | M-of-N | Low |

---

## Quick Reference

### Identifying BTC Address Types

| Prefix | Type | BIP | SegWit |
|--------|------|-----|--------|
| `1...` | P2PKH | BIP44 | ❌ |
| `3...` | P2SH or P2SH-P2WPKH | BIP16/BIP49 | △ |
| `bc1q...` | P2WPKH or P2WSH | BIP84 | ✅ |
| `bc1p...` | P2TR (Taproot) | BIP86 | ✅ |

### Identifying BCH Address Types

| Prefix | Type | Multisig |
|--------|------|----------|
| `bitcoincash:q...` | P2PKH | ❌ |
| `bitcoincash:p...` | P2SH | ✅ |

### Transaction Size Comparison

| Pattern | Weight | vBytes | Notes |
|---------|--------|--------|-------|
| P2PKH Single-sig (1-in, 2-out) | ~680 | ~170 | Legacy |
| P2WPKH Single-sig (1-in, 2-out) | ~440 | ~110 | Native SegWit |
| P2TR Single-sig (1-in, 2-out) | ~396 | ~99 | Taproot |
| 2-of-3 P2WSH Multisig | ~1,100 | ~275 | Traditional Multisig |
| 2-of-3 MuSig2 (P2TR) | ~560 | ~140 | Signature Aggregation |

---

## Related Documents

- [BTC Technical Reference](./README.md) - Bitcoin technical reference
- [Taproot User Guide](./TAPROOT_GUIDE.md) - How to use Taproot
- [MuSig2 User Guide](./musig2_guide.md) - How to use MuSig2
- [Descriptor Examples](./descriptor_examples.md) - Descriptor examples
- [PSBT Developer Guide](./psbt_developer_guide.md) - PSBT development guide
- [BCH E2E Workflow](../../../scripts/operation/bch/README.md) - BCH E2E workflow

---

**Document Version:** 1.2
**Last Updated:** 2026-01-15
**Maintainer:** go-crypto-wallet team

---

## Changelog

### Version 1.2 (2026-01-15)

- ✅ Pattern 2 (P2PKH 2-of-3 Multisig) is now fully working
- Fixed key derivation path mismatch for multisig accounts (PR #357)
- Added detailed explanation of the fix and root cause
- Updated implementation status for 2-of-3 multisig

### Version 1.1 (2026-01-14)

- Initial comprehensive documentation of all patterns
