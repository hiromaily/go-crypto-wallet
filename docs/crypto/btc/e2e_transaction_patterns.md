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
| **1** | **P2PKH (BIP44)** | **Single-sig** | **`1...`** | **✅ e2e/e2e-p1-p2pkh-singlesig.sh** |
| **2** | **P2PKH (BIP44)** | **2-of-3 Multisig** | **`3...` (P2SH wrapped)** | **✅ e2e/e2e-p2-p2pkh-2of3.sh** |
| **3** | **P2SH-P2WPKH (BIP49)** | **Single-sig** | **`3...`** | **✅ e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh** |
| **4** | **P2SH-P2WSH (BIP49)** | **2-of-3 Multisig** | **`3...`/`2...`** | **✅ e2e/e2e-p4-p2sh-p2wsh-2of3.sh** |
| 5 | P2WPKH (BIP84) | Single-sig | `bc1q...` | 🔶 Manual testing |
| 6 | P2WSH (BIP84) | 2-of-3 Multisig | `bc1q...` | ❌ Not supported |
| 7 | P2WSH (BIP84) | 3-of-3 Multisig | `bc1q...` | ❌ Not supported |
| **8** | **P2SH-P2WSH** | **3-of-3 Multisig** | **`3...`** | **🔶 e2e/e2e-p8-p2sh-p2wsh-3of3.sh** (WIP) |
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

## Details of Each Pattern for BTC

### Pattern 1: BTC P2PKH Single-sig

**Currently implemented in `scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh`**

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

### Pattern 2: BTC P2PKH 2-of-3 Multisig

**✅ Fully implemented in `scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh`**

```
Address Type: P2PKH (BIP44 Legacy) with 2-of-3 Multisig
Signing Requirements: 2-of-3 (Keygen + Sign1, Sign2 is optional)
Descriptor: sh(multi(2, [fingerprint/44'/1'/1]xpub1/0/*, [fingerprint/44'/1'/1]xpub2/0/*, [fingerprint/44'/1'/1]xpub3/0/*))
Address Format: 2... (P2SH in regtest), 3... (P2SH in mainnet)
```

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

### Pattern 3: BTC P2SH-P2WPKH Single-sig

**Fully implemented and verified in `scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh`**

```
Address Type: P2SH-P2WPKH (BIP49 Nested SegWit)
Signing Requirements: Single-sig (Keygen only)
Descriptor: sh(wpkh([fingerprint/49'/0'/0']xpub.../0/*))
Address Format: 2... (P2SH in regtest), 3... (P2SH in mainnet)
```

#### What is P2SH-P2WPKH (BIP49)?

P2SH-P2WPKH is **SegWit (P2WPKH) wrapped in P2SH for backward compatibility**. It was introduced in **BIP49** as a transitional format during SegWit adoption.

| BIP | Role |
|-----|------|
| **BIP16** | Defines P2SH (Pay-to-Script-Hash) |
| **BIP49** | Defines P2SH-P2WPKH derivation path (`m/49'/...`) |
| **BIP141** | Defines SegWit (P2WPKH witness program) |
| **BIP143** | Defines SegWit transaction signing |

#### Address Structure

```
P2SH-P2WPKH Address:
┌─────────────────────────────────────────────────────┐
│  P2SH wrapper (for legacy wallet compatibility)     │
│  ┌─────────────────────────────────────────────────┐│
│  │  P2WPKH (Native SegWit)                         ││
│  │  ┌─────────────────────────────────────────────┐││
│  │  │  Public Key Hash (20 bytes)                 │││
│  │  └─────────────────────────────────────────────┘││
│  └─────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────┘
```

#### Comparison with Other Single-sig Patterns

| Item | Pattern 1 (P2PKH) | Pattern 3 (P2SH-P2WPKH) | Pattern 5 (P2WPKH) | Pattern 9 (P2TR) |
|------|-------------------|-------------------------|--------------------|--------------------|
| BIP | BIP44 | **BIP49** | BIP84 | BIP86 |
| Address Prefix | `1...`/`m...` | **`3...`/`2...`** | `bc1q...` | `bc1p...` |
| Descriptor | `pkh(...)` | **`sh(wpkh(...))`** | `wpkh(...)` | `tr(...)` |
| SegWit | No | **Yes (wrapped)** | Yes (native) | Yes (Taproot) |
| Transaction Size | Largest | Medium | Smaller | Smallest |
| Legacy Compatible | Yes | **Yes** | No | No |

#### Why Use P2SH-P2WPKH?

| Aspect | Description |
|--------|-------------|
| ✅ SegWit benefits | Reduced transaction size compared to P2PKH, lower fees |
| ✅ Legacy compatibility | `3...` addresses can receive from any wallet (including old wallets) |
| ✅ Widely supported | Supported by most exchanges and services |
| ❌ Not optimal | Native SegWit (P2WPKH) is more efficient |
| ❌ Transitional | Primarily for backward compatibility during SegWit migration |

**Workflow:**

1. Generate Seed in Keygen
2. Generate BIP49 HD Key in Keygen (10 accounts each)
3. Export Descriptor from Keygen (generates `sh(wpkh(...))`)
4. Import Descriptor to Watch
5. Generate Test UTXO (regtest)
6. Create unsigned transaction → Sign once → Broadcast

**Characteristics:**

- Simple and fast (completed with single signature)
- Sign1/Sign2 wallets not required
- Uses BIP49 key derivation path (`m/49'/0'/account'/change/index`)
- P2SH address format (`2...` prefix in regtest)
- SegWit transaction efficiency with legacy address compatibility

---

### Pattern 4: BTC P2SH-P2WSH 2-of-3 Multisig

**Fully implemented and verified in `scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh`**

```
Address Type: P2SH-P2WSH (BIP49 Nested SegWit Multisig)
Signing Requirements: 2-of-3 Multisig (any 2 signatures out of 3)
Descriptor: sh(wsh(sortedmulti(2,[fp1/49'/1'/1']xpub1/0/*,[fp2/49'/1'/1']xpub2/0/*,[fp3/49'/1'/1']xpub3/0/*)))
Address Format: 2... (P2SH in regtest), 3... (P2SH in mainnet)
```

#### What is P2SH-P2WSH 2-of-3?

P2SH-P2WSH is **SegWit multisig (P2WSH) wrapped in P2SH for backward compatibility**. The 2-of-3 configuration requires **any 2 signatures out of 3 possible keys** to authorize a transaction.

| BIP | Role |
|-----|------|
| **BIP16** | Defines P2SH (Pay-to-Script-Hash) |
| **BIP49** | Defines P2SH-SegWit derivation path (`m/49'/...`) |
| **BIP141** | Defines SegWit witness program |
| **BIP143** | Defines SegWit transaction signing |
| **BIP67** | Defines sorted multisig keys (lexicographic ordering) |

#### Address Structure

```
P2SH-P2WSH 2-of-3 Address:
┌─────────────────────────────────────────────────────────────┐
│  P2SH wrapper (for legacy wallet compatibility)              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  P2WSH (SegWit Script Hash)                              │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  sortedmulti(2, pubkey1, pubkey2, pubkey3)           │ │ │
│  │  │  → 2-of-3 multisig script                            │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

#### Comparison with Related Patterns

| Item | Pattern 2 (P2PKH 2-of-3) | **Pattern 4 (P2SH-P2WSH 2-of-3)** | Pattern 8 (P2SH-P2WSH 3-of-3) |
|------|-------------------------|------------------------------------|-------------------------------|
| BIP | BIP44 | **BIP49** | BIP49 |
| Address Prefix | `2...` (P2SH) | **`2...` (P2SH)** | `2...` (P2SH) |
| Descriptor | `sh(multi(2,...))` | **`sh(wsh(sortedmulti(2,...)))`** | `sh(wsh(sortedmulti(3,...)))` |
| SegWit | No | **Yes (wrapped)** | Yes (wrapped) |
| Signature Requirement | 2-of-3 | **2-of-3** | 3-of-3 (all required) |
| Transaction Size | Larger | **Smaller** | Similar to P4 |
| Threshold Flexibility | Medium | **Medium** | Low (all must sign) |

#### Why Use P2SH-P2WSH 2-of-3?

| Aspect | Description |
|--------|-------------|
| ✅ Threshold security | Requires only 2 out of 3 keys - provides redundancy and operational flexibility |
| ✅ SegWit efficiency | Reduced transaction size and fees compared to legacy multisig |
| ✅ Legacy compatibility | `2...` addresses can receive from any wallet |
| ✅ Sorted keys | Uses BIP67 sorted multisig for deterministic key ordering |
| ❌ Complex setup | Requires 3 separate wallets (keygen + sign1 + sign2) |
| ❌ Multiple signatures | Requires coordination between 2 signers |

**Workflow:**

1. Generate Seeds in Keygen, Sign1, Sign2 wallets
2. Generate BIP49 HD Keys in all wallets (deposit, payment, stored accounts)
3. Export fullpubkeys from Sign1 and Sign2 wallets
4. Import fullpubkeys into Keygen wallet
5. Export 2-of-3 multisig descriptors from Keygen (generates `sh(wsh(sortedmulti(2,...)))`)
6. Import descriptors into Watch wallet
7. Generate Test UTXOs (regtest)
8. Create unsigned transaction → Sign with Keygen (1st) → Sign with Sign1 (2nd) → Broadcast

**Characteristics:**

- **2-of-3 threshold**: Any 2 signatures out of 3 are sufficient (Sign2 not needed if Keygen + Sign1 sign)
- Uses BIP49 key derivation path (`m/49'/1'/1'/change/index` for multisig)
- P2SH address format (`2...` prefix in regtest, `3...` in mainnet)
- SegWit transaction efficiency with legacy address compatibility
- Configuration: Uses `account_2of3.yaml` for multisig account settings
- **Security model**: Provides redundancy - losing one key doesn't prevent access to funds

**Signing Flow:**

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (1st signature)
    ↓
Sign1 Wallet (2nd signature) ← Complete here! (2-of-3 satisfied)
    ↓
Watch Wallet (broadcast)

Note: Sign2 is optional - 2 signatures already satisfy the 2-of-3 requirement
```

---

### Pattern 8: BTC P2SH-P2WSH 3-of-3 Multisig (WIP)

**Currently WIP implemented in `scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh`**

```
Address Type: P2SH-P2WSH (SegWit multisig wrapped in P2SH)
Signing Requirements: 3-of-3 (Keygen + Sign1 + Sign2)
Descriptor: sh(wsh(sortedmulti(3, xpub1, xpub2, xpub3)))
```

#### What is P2SH-P2WSH?

P2SH-P2WSH is a **nested structure where P2WSH (Pay-to-Witness-Script-Hash) is wrapped inside P2SH (Pay-to-Script-Hash)**. This is achieved through a combination of multiple BIPs:

| BIP | Role |
|-----|------|
| **BIP16** | Defines P2SH (Pay-to-Script-Hash) |
| **BIP141** | Defines SegWit (including P2WSH) |
| **BIP143** | Defines SegWit transaction signing |

**Note:** P2SH-P2WSH is different from BIP49 (P2SH-P2WPKH):

- **BIP49 (P2SH-P2WPKH)**: Wraps **P2WPKH** (single public key) in P2SH → For Single-sig
- **P2SH-P2WSH**: Wraps **P2WSH** (witness script) in P2SH → **For Multisig**

#### Why No Single-sig or 2-of-3 Variants?

| Variant | Reason for Absence |
|---------|-------------------|
| **Single-sig** | P2SH-P2WPKH (BIP49, Pattern 3) is sufficient. P2SH-P2WSH is designed for complex scripts |
| **2-of-3** | Technically possible but not implemented in this project (prioritization) |

P2SH-P2WSH is primarily used for **complex multisig scripts** while maintaining SegWit efficiency and backward compatibility with legacy wallets.

#### Comparison with Other Patterns

| Pattern | Descriptor | Use Case |
|---------|------------|----------|
| Pattern 2 | `sh(multi(2, ...))` | Legacy P2SH multisig |
| Pattern 6-7 | `wsh(sortedmulti(...))` | Native SegWit multisig |
| **Pattern 8** | **`sh(wsh(sortedmulti(...)))`** | **SegWit multisig + Legacy compatibility** |

#### Pros and Cons

| Aspect | Description |
|--------|-------------|
| ✅ Legacy compatibility | `3...` addresses can receive from legacy wallets |
| ✅ SegWit efficiency | Witness data stored separately, reducing transaction size |
| ❌ Complexity | Nested structure makes implementation and debugging more complex |
| ❌ Size | Slightly larger than Native SegWit (P2WSH)

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

## Details of Each Pattern for BCH

### BCH Pattern 2: BCH CashAddr 3-of-3 Multisig (Current E2E)

**Currently WIP implemented in `scripts/operation/bch/e2e-workflow.sh`**

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
| `scripts/operation/btc/e2e/e2e-p1-p2pkh-singlesig.sh` | BTC | P2PKH Single-sig (Pattern 1) | Single-sig |
| `scripts/operation/btc/e2e/e2e-p2-p2pkh-2of3.sh` | BTC | P2PKH 2-of-3 Multisig (Pattern 2) | 2-of-3 |
| `scripts/operation/btc/e2e/e2e-p3-p2sh-p2wpkh-singlesig.sh` | BTC | P2SH-P2WPKH Single-sig (Pattern 3) | Single-sig |
| `scripts/operation/btc/e2e/e2e-p4-p2sh-p2wsh-2of3.sh` | BTC | P2SH-P2WSH 2-of-3 Multisig (Pattern 4) | 2-of-3 |
| `scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh` | BTC | P2SH-P2WSH 3-of-3 Multisig (Pattern 8) | 3-of-3 |
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

**Document Version:** 1.3
**Last Updated:** 2026-01-15
**Maintainer:** go-crypto-wallet team

---

## Changelog

### Version 1.3 (2026-01-15)

- ✅ Pattern 3 (P2SH-P2WPKH Single-sig) is now fully working
- E2E script `e2e-p3-p2sh-p2wpkh-singlesig.sh` completed and verified

### Version 1.2 (2026-01-15)

- ✅ Pattern 2 (P2PKH 2-of-3 Multisig) is now fully working
- Fixed key derivation path mismatch for multisig accounts (PR #357)
- Added detailed explanation of the fix and root cause
- Updated implementation status for 2-of-3 multisig

### Version 1.1 (2026-01-14)

- Initial comprehensive documentation of all patterns
