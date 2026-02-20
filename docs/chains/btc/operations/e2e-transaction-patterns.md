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
| **5** | **P2WPKH (BIP84)** | **Single-sig** | **`bc1q...`** | **✅ e2e/e2e-p5-p2wpkh-singlesig.sh** |
| **6** | **P2WSH (BIP84)** | **2-of-3 Multisig** | **`bc1q...`** | **✅ e2e/e2e-p6-p2wsh-2of3.sh** |
| **7** | **P2WSH (BIP84)** | **3-of-3 Multisig** | **`bc1q...`** | **✅ e2e/e2e-p7-p2wsh-3of3.sh** |
| **8** | **P2SH-P2WSH** | **3-of-3 Multisig** | **`3...`** | **✅ e2e/e2e-p8-p2sh-p2wsh-3of3.sh** |
| **9** | **P2TR (BIP86)** | **Single-sig** | **`bc1p...`** | **✅ e2e/e2e-p9-p2tr-singlesig.sh** |
| **10** | **P2TR (BIP86)** | **MuSig2 (N-of-N)** | **`bc1p...`** | **✅ e2e/e2e-p10-p2tr-musig2.sh** |
| **11** | **P2TR (BIP86)** | **Tapscript (M-of-N)** | **`bc1p...`** | **✅ e2e/e2e-p11-p2tr-tapscript.sh** |

### BCH Pattern Matrix

BCH supports **fewer patterns** than BTC due to lack of SegWit, Taproot, and Schnorr signatures.

| Pattern | Key Type | Signature Pattern | Address Format | E2E Script Support |
|---------|----------|-------------------|----------------|-------------------|
| **1** | **CashAddr P2PKH** | **Single-sig** | **`bitcoincash:q...`** | **✅ e2e/e2e-p1-p2pkh-singlesig.sh** |
| **2** | **CashAddr P2SH** | **2-of-3 Multisig** | **`bitcoincash:p...`** | **✅ e2e/e2e-p2-p2sh-2of3.sh** |
| **3** | **CashAddr P2SH** | **3-of-3 Multisig** | **`bitcoincash:p...`** | **✅ e2e/e2e-p3-p2sh-3of3.sh** |

**BCH Limitations:**

- ❌ No SegWit (no P2WPKH, P2WSH, P2SH-P2WPKH patterns)
- ❌ No Taproot (no P2TR patterns)
- ❌ No Schnorr signatures (ECDSA only)
- ❌ No MuSig2 (no signature aggregation)
- ❌ No Descriptor support (use address export/import instead)
- ❌ No PSBT format (use raw transaction hex)
- ❌ No Bech32/Bech32m encoding (use CashAddr)
- ❌ No BIP49/84/86 derivation paths (BIP44 only)
- ⚠️ Fee unit: sat/Byte (not sat/vByte - no witness discount)

For detailed BCH patterns, see [BCH Technical Reference](../../bch/README.md#e2e-transaction-patterns).

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

### Pattern 6: BTC P2WSH Native SegWit 2-of-3 Multisig

**✅ Fully implemented in `scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh`**

```
Address Type: P2WSH (BIP84 Native SegWit Multisig)
Signing Requirements: 2-of-3 Multisig (any 2 signatures out of 3)
Descriptor: wsh(sortedmulti(2,[fp1/84'/1'/1']xpub1/0/*,[fp2/84'/1'/1']xpub2/0/*,[fp3/84'/1'/1']xpub3/0/*))
Address Format: bcrt1q... (62 chars for P2WSH in regtest), bc1q... (mainnet)
```

#### What is P2WSH 2-of-3?

P2WSH (Pay-to-Witness-Script-Hash) is **native SegWit multisig without P2SH wrapper**. The 2-of-3 configuration requires **any 2 signatures out of 3 possible keys** to authorize a transaction. This is the most efficient multisig format for Bitcoin.

| BIP | Role |
|-----|------|
| **BIP84** | Defines native SegWit derivation path (`m/84'/...`) |
| **BIP141** | Defines SegWit witness program (P2WSH) |
| **BIP143** | Defines SegWit transaction signing |
| **BIP67** | Defines sorted multisig keys (lexicographic ordering) |

#### Address Structure

```
P2WSH 2-of-3 Address (Native SegWit):
┌─────────────────────────────────────────────────────────────┐
│  P2WSH (Native SegWit - no P2SH wrapper)                     │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  sortedmulti(2, pubkey1, pubkey2, pubkey3)               │ │
│  │  → 2-of-3 multisig witness script                        │ │
│  │  → SHA256 hash becomes witness program (32 bytes)        │ │
│  │  → Bech32-encoded as bc1q... address                     │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

#### Comparison with Related Patterns

| Item | Pattern 4 (P2SH-P2WSH 2-of-3) | **Pattern 6 (P2WSH 2-of-3)** | Pattern 7 (P2WSH 3-of-3) |
|------|-------------------------------|-----------------------------|--------------------------|
| BIP | BIP49 | **BIP84** | BIP84 |
| Address Prefix | `2...` (P2SH) | **`bc1q...` (62 chars)** | `bc1q...` (62 chars) |
| Descriptor | `sh(wsh(sortedmulti(2,...)))` | **`wsh(sortedmulti(2,...))`** | `wsh(sortedmulti(3,...))` |
| SegWit | Yes (wrapped) | **Yes (native)** | Yes (native) |
| Signature Requirement | 2-of-3 | **2-of-3** | 3-of-3 (all required) |
| Transaction Size | Medium | **Smallest** | Slightly larger than P6 |
| Legacy Compatible | Yes | **No** | No |

#### Why Use P2WSH 2-of-3?

| Aspect | Description |
|--------|-------------|
| ✅ Threshold security | Requires only 2 out of 3 keys - provides redundancy and operational flexibility |
| ✅ Maximum efficiency | No P2SH wrapper - smallest transaction size for multisig |
| ✅ Lowest fees | Native SegWit provides the best fee rates |
| ✅ Sorted keys | Uses BIP67 sorted multisig for deterministic key ordering |
| ✅ Descriptor-based | Full Bitcoin Core descriptor wallet support |
| ❌ No legacy compatibility | Cannot receive from very old wallets (requires Bech32 support) |
| ❌ Complex setup | Requires 3 separate wallets (keygen + sign1 + sign2) |

**Workflow:**

1. Generate Seeds in Keygen, Sign1, Sign2 wallets
2. Generate BIP84 HD Keys in all wallets (deposit, payment, stored accounts)
3. Export fullpubkeys from Sign1 and Sign2 wallets
4. Import fullpubkeys into Keygen wallet
5. Export 2-of-3 multisig descriptors from Keygen (generates `wsh(sortedmulti(2,...))`)
6. Import descriptors into Watch wallet
7. Generate Test UTXOs (regtest)
8. Create unsigned transaction → Sign with Keygen (1st) → Sign with Sign1 (2nd) → Broadcast

**Characteristics:**

- **2-of-3 threshold**: Any 2 signatures out of 3 are sufficient (Sign2 not needed if Keygen + Sign1 sign)
- Uses BIP84 key derivation path (`m/84'/1'/1'/change/index` for multisig)
- Native SegWit Bech32 address format (`bcrt1q...` 62 chars in regtest, `bc1q...` in mainnet)
- Most efficient multisig format (no P2SH overhead)
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

**Implementation Notes:**

- Requires native SegWit descriptor support in PSBT creation
- Address derivation uses `NewAddressWitnessScriptHash()` with SHA256 hash
- Witness script format: `OP_2 <pubkey1> <pubkey2> <pubkey3> OP_3 OP_CHECKMULTISIG`
- Public keys are sorted lexicographically per BIP67

---

### Pattern 8: BTC P2SH-P2WSH 3-of-3 Multisig

**✅ Fully implemented in `scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh`**

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

**✅ Fully implemented in `scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh`**

```
Address Type: P2TR (BIP86 Taproot)
Signing Requirements: Single-sig (Keygen only)
Descriptor: tr([fingerprint/86'/0'/0']xpub.../0/*)
Address Format: bc1p... (Mainnet), tb1p... (Testnet), bcrt1p... (regtest)
Bitcoin Core: v22.0+ required
```

#### What is P2TR (Taproot)?

P2TR (Pay-to-Taproot) is Bitcoin's latest major upgrade (activated November 2021) that introduces a new address format and signature scheme. It combines **Schnorr signatures** and **Merkelized Alternative Script Trees (MAST)** into a single output type.

| BIP | Role |
|-----|------|
| **BIP86** | Defines key derivation path (`m/86'/...`) for Taproot key path spend |
| **BIP340** | Defines Schnorr signature scheme (64 bytes vs 71-72 for ECDSA) |
| **BIP341** | Defines Taproot output structure and spending rules |
| **BIP342** | Defines Tapscript (script semantics for Taproot) |

#### Address Structure

```
P2TR Address (SegWit v1 - Taproot):
┌─────────────────────────────────────────────────────────────┐
│  Bech32m Encoding (NOT Bech32!)                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Witness Version: 1 (different from v0 in P2WPKH)       │ │
│  │  ┌─────────────────────────────────────────────────────┐ │ │
│  │  │  Tweaked Public Key (32 bytes, x-only format)        │ │ │
│  │  │  output_key = internal_key + H_TapTweak(P) * G       │ │ │
│  │  └─────────────────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘

Format: bc1p + 58 characters (mainnet, 62 chars total)
        tb1p + 58 characters (testnet/signet)
        bcrt1p + 58 characters (regtest)
```

#### Comparison with Other Single-sig Patterns

| Item | Pattern 1 (P2PKH) | Pattern 3 (P2SH-P2WPKH) | Pattern 5 (P2WPKH) | **Pattern 9 (P2TR)** |
|------|-------------------|-------------------------|--------------------|--------------------|
| BIP | BIP44 | BIP49 | BIP84 | **BIP86** |
| Address Prefix | `1...`/`m...` | `3...`/`2...` | `bc1q...` | **`bc1p...`** |
| Descriptor | `pkh(...)` | `sh(wpkh(...))` | `wpkh(...)` | **`tr(...)`** |
| SegWit Version | N/A | v0 (wrapped) | v0 (native) | **v1 (Taproot)** |
| Signature | ECDSA | ECDSA | ECDSA | **Schnorr (BIP340)** |
| Signature Size | 71-72 bytes | 71-72 bytes | 71-72 bytes | **64 bytes** |
| Transaction Size | Largest | Medium | Smaller | **Smallest** |
| Encoding | Base58Check | Base58Check | Bech32 | **Bech32m** |
| address_type | `legacy` | `p2sh-segwit` | `bech32` | **`taproot`** |

#### Key Path vs Script Path Spending

Taproot supports two spending paths:

| Spend Type | Description | Use Case | Pattern |
|------------|-------------|----------|---------|
| **Key Path** | Direct signature with tweaked key | Single-sig, efficient multisig | **Pattern 9** |
| Script Path | Reveal Merkle proof + satisfy script | Complex conditions, timelocks | Pattern 11 |

Pattern 9 uses **Key Path Spend** - the most efficient method when only a single key controls the funds.

#### Why Use P2TR (Taproot)?

| Aspect | Description |
|--------|-------------|
| ✅ Smallest transactions | ~30-40% smaller than P2PKH, ~15% smaller than P2WPKH |
| ✅ Lowest fees | Most efficient single-sig format |
| ✅ Enhanced privacy | All Taproot spends (single-sig and multisig) look identical |
| ✅ Schnorr signatures | 64 bytes vs 71-72 for ECDSA, batch verification possible |
| ✅ Future-proof | Latest Bitcoin address standard |
| ✅ Faster validation | Schnorr signatures verify faster than ECDSA |
| ❌ Requires Bitcoin Core v22.0+ | Older nodes cannot create/validate Taproot |
| ❌ Newer standard | Some services may not support `bc1p...` addresses yet |

#### Key Derivation Path (BIP86)

```
m / 86' / coin_type' / account' / change / address_index
         └─ 0' for mainnet, 1' for testnet/regtest
```

#### Schnorr Signature (BIP340)

```
Schnorr Signature: 64 bytes total
  ┌──────────────────────────────────────────────────────────┐
  │  r (32 bytes) │ s (32 bytes)                              │
  │  x-coord of R │ scalar value                              │
  └──────────────────────────────────────────────────────────┘

ECDSA Signature (DER encoded): 71-72 bytes
  ┌──────────────────────────────────────────────────────────┐
  │ 30 │ len │ 02 │ r_len │ r │ 02 │ s_len │ s              │
  └──────────────────────────────────────────────────────────┘
```

**Workflow:**

1. Generate Seed in Keygen
2. Generate BIP86 HD Key in Keygen (10 accounts each)
3. Export Taproot Descriptor from Keygen (generates `tr(...)`)
4. Import Descriptor to Watch
5. Generate Test UTXO (regtest)
6. Create unsigned transaction → Sign once (Schnorr) → Broadcast

**Signing Flow (Single-sig):**

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (sign with Schnorr signature - 64 bytes)
    ↓
Watch Wallet (broadcast)
```

**Characteristics:**

- Simple and fast (completed with single Schnorr signature)
- Sign1/Sign2 wallets not required
- Uses BIP86 key derivation path (`m/86'/0'/account'/change/index`)
- Taproot address format (`bcrt1p...` prefix in regtest)
- Most efficient single-sig transaction format
- **Requires Bitcoin Core v22.0+** (Taproot/Schnorr support)

**Implementation Notes:**

- Environment variable: `WALLET_ADDRESS_TYPE="taproot"`
- Descriptor format: `tr([fingerprint/86'/coin'/account']xpub.../0/*)`
- Uses Bech32m encoding (NOT Bech32 - different checksum constant)
- Witness version 1 (vs version 0 for P2WPKH)
- Address derivation uses x-only public keys (32 bytes, parity byte removed)

### Pattern 10: BTC P2TR MuSig2

**Implementation Status:** Framework implemented - CLI commands have TODOs for full MuSig2 protocol
**E2E Test Status:** ✅ Verified (2026-01-16) - Infrastructure, wallets, Taproot addresses, and workflow framework working correctly

```
Address Type: P2TR (BIP86)
Signing Requirements: N-of-N MuSig2 (all signatures required)
Descriptor: tr(musig([fingerprint1/86'/1'/1']xpub1,[fingerprint2]xpub2,[fingerprint3]xpub3)/0/*)
```

**2-Round Protocol:**

1. Round 1: Generate nonces in each wallet (parallelizable)
2. Round 2: Create partial signatures in each wallet (sequential)
3. Aggregate signatures in Watch and broadcast

**Implementation Notes:**

- Environment variable: `WALLET_ADDRESS_TYPE="taproot"`
- Uses BIP86 key derivation (m/86'/coin'/account'/change/index)
- Address format: `bcrt1p...` (regtest), `bc1p...` (mainnet) - bech32m encoding
- Requires Bitcoin Core v22.0+ for descriptor-based wallets
- 2-round signing protocol (BIP327)
- Framework script available at `scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh`

**E2E Test Verified Components:**

- ✅ Infrastructure setup (Docker, Bitcoin Core regtest, MySQL)
- ✅ Wallet creation and configuration (watch, keygen, sign1, sign2)
- ✅ HD key generation for all accounts (deposit, payment, stored)
- ✅ Taproot address generation with correct bech32m encoding
- ✅ Descriptor export and import (2000 addresses generated)
- ✅ UTXO generation and balance verification
- ✅ Payment workflow and transaction creation
- ⚠️ MuSig2 protocol (Round 1/2, aggregation) - Placeholder implementation pending full CLI support

## Details of Each Pattern for BCH

### BCH Pattern 3: BCH CashAddr P2SH 3-of-3 Multisig

**✅ Fully implemented and verified in `scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh`**

```
Address Type: CashAddr P2SH (BIP44 + BIP11)
Signing Requirements: 3-of-3 (Keygen + Sign1 + Sign2)
Address Format: bchreg:p... (P2SH multisig in regtest)
Key Derivation: m/44'/1'/account'/change/index (testnet/regtest)
```

**Implementation Status:**

- ✅ Infrastructure and wallet setup
- ✅ HD key generation and fullpubkey export/import
- ✅ 3-of-3 multisig address creation
- ✅ Payment request creation and database storage
- ✅ **UTXO retrieval** (Fixed in PR #426 - CashAddr format normalization)
- ✅ Transaction creation and signing workflow
- ⚠️ PSBT generation (Known issue with BCH label format)

**Workflow:**

1. Generate Seed in Keygen/Sign1/Sign2
2. Generate HD Key in Keygen (10 keys per account)
3. Generate HD Key in Sign1/Sign2
4. Export fullpubkey from Sign1/Sign2
5. Import fullpubkey to Keygen
6. Create 3-of-3 Multisig addresses in Keygen
7. Export addresses from Keygen (CashAddr format)
8. Import addresses to Watch wallet
9. Generate Test UTXO (regtest)
10. Create unsigned transaction → Sign 3 times (ECDSA) → Broadcast

**Signing Flow:**

```
Watch Wallet (create unsigned tx)
    ↓
Keygen Wallet (1st ECDSA signature)
    ↓
Sign1 Wallet (2nd ECDSA signature)
    ↓
Sign2 Wallet (3rd ECDSA signature)
    ↓
Watch Wallet (broadcast)
```

**Note:** Unlike BTC's SegWit patterns, BCH transactions don't benefit from witness data separation, resulting in larger transaction sizes. BCH compensates with very low fees (~1 sat/byte).

For more BCH patterns, see [BCH Technical Reference](../../bch/README.md#e2e-transaction-patterns)

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
| `scripts/operation/btc/e2e/e2e-p5-p2wpkh-singlesig.sh` | BTC | P2WPKH Native SegWit Single-sig (Pattern 5) | Single-sig |
| `scripts/operation/btc/e2e/e2e-p6-p2wsh-2of3.sh` | BTC | P2WSH Native SegWit 2-of-3 Multisig (Pattern 6) | 2-of-3 |
| `scripts/operation/btc/e2e/e2e-p7-p2wsh-3of3.sh` | BTC | P2WSH Native SegWit 3-of-3 Multisig (Pattern 7) | 3-of-3 |
| `scripts/operation/btc/e2e/e2e-p8-p2sh-p2wsh-3of3.sh` | BTC | P2SH-P2WSH 3-of-3 Multisig (Pattern 8) | 3-of-3 |
| `scripts/operation/btc/e2e/e2e-p9-p2tr-singlesig.sh` | BTC | P2TR Taproot Single-sig (Pattern 9) | Single-sig |
| `scripts/operation/btc/e2e/e2e-p10-p2tr-musig2.sh` | BTC | P2TR MuSig2 (Pattern 10) | N-of-N (framework) |
| `scripts/operation/btc/e2e/e2e-p11-p2tr-tapscript.sh` | BTC | P2TR Tapscript (Pattern 11) | 2-of-3 (framework) |
| `scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh` | BCH | CashAddr P2PKH Single-sig (Pattern 1) | Single-sig |
| `scripts/operation/bch/e2e/e2e-p2-p2sh-2of3.sh` | BCH | CashAddr P2SH 2-of-3 Multisig (Pattern 2) | 2-of-3 |
| `scripts/operation/bch/e2e/e2e-p3-p2sh-3of3.sh` | BCH | CashAddr P2SH 3-of-3 Multisig (Pattern 3) | 3-of-3 |

### Planned E2E Scripts

No additional E2E scripts planned at this time. All 11 Bitcoin patterns are implemented.

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

**Document Version:** 1.10
**Last Updated:** 2026-01-17
**Maintainer:** go-crypto-wallet team

---

## Changelog

### Version 1.10 (2026-01-17)

- ✅ Fixed BCH Pattern 3 (3-of-3 Multisig) UTXO retrieval issue (Closes #423, PR #426)
- Resolved CashAddr format mismatch in `ListUnspentByAccount` address comparison
- Updated BCH Pattern Matrix to show all 3 patterns with correct E2E script paths
- Updated E2E Script Reference with proper BCH script locations
- Added implementation status details for BCH Pattern 3
- All BCH patterns now have dedicated E2E scripts in `scripts/operation/bch/e2e/`

### Version 1.9 (2026-01-17)

- Updated BCH Pattern Matrix with accurate pattern numbering (1, 2, 3)
- Added BCH limitations section (no SegWit, Taproot, Schnorr, MuSig2)
- Updated BCH Pattern 3 (3-of-3 Multisig) with detailed workflow and signing flow
- Added cross-reference link to BCH Technical Reference
- Renamed "BCH Pattern 2" to "BCH Pattern 3" to reflect correct pattern numbering

### Version 1.8 (2026-01-16)

- ✅ Pattern 11 (P2TR Tapscript M-of-N) framework implemented
- E2E script `e2e-p11-p2tr-tapscript.sh` created (Closes #381)
- Implements Tapscript Script Path spending framework with 2-of-3 threshold
- Uses BIP86 key derivation + BIP342 Tapscript semantics
- Script tree with Merkle proof and control block structure
- M × Schnorr signatures for Script Path spend
- Address format: `bcrt1p...` (62 chars, Bech32m encoding)
- ~50% smaller than P2WSH 2-of-3 multisig
- Enhanced privacy: unused script paths hidden in Merkle tree
- Note: Full Tapscript implementation pending (currently uses placeholder)

### Version 1.7 (2026-01-16)

- ✅ Pattern 9 (P2TR Taproot Single-sig) is now fully working
- E2E script `e2e-p9-p2tr-singlesig.sh` completed and verified (Closes #377)
- Fixed Taproot address derivation to use x-only public keys (32 bytes)
- Implements BIP86 key derivation for Taproot key path spending
- Uses Schnorr signatures (BIP340, 64 bytes)
- Most efficient single-sig transaction format
- Address format: `bcrt1p...` (62 chars, Bech32m encoding)

### Version 1.6 (2026-01-16)

- ✅ Pattern 8 (P2SH-P2WSH 3-of-3 Multisig) is now fully working
- E2E script `e2e-p8-p2sh-p2wsh-3of3.sh` completed and verified
- Fixed receiver address generation to use P2SH-SegWit format (Closes #374)
- P2SH-wrapped SegWit multisig with legacy compatibility (`2...` addresses in regtest)
- Implements BIP49 key derivation for 3-of-3 multisig
- All 3 signatures required (Keygen + Sign1 + Sign2)

### Version 1.5 (2026-01-16)

- ✅ Pattern 6 (P2WSH Native SegWit 2-of-3 Multisig) is now fully working
- E2E script `e2e-p6-p2wsh-2of3.sh` completed and verified
- Native SegWit multisig with Bech32 encoding (`bcrt1q...` 62-char addresses)
- Added native SegWit descriptor support (`wsh`, `wpkh`) in PSBT infrastructure
- Most efficient multisig format - no P2SH wrapper overhead
- Implements BIP84 key derivation and BIP67 sorted multisig keys

### Version 1.4 (2026-01-15)

- ✅ Pattern 5 (P2WPKH Native SegWit Single-sig) is now fully working
- E2E script `e2e-p5-p2wpkh-singlesig.sh` completed and verified
- Native SegWit with Bech32 encoding (`bcrt1q...` addresses)

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
