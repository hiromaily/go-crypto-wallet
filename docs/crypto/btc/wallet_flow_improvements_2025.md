# Bitcoin Wallet Flow Improvements (2025)

This document outlines improvements for the Bitcoin wallet flow, including the roles of Keygen Wallet, Sign Wallet, and Watch Wallet, with a focus on modern key technologies (Taproot, MuSig2, etc.).

## Table of Contents

1. [Current Flow Overview](#1-current-flow-overview)
2. [Wallet Roles and Responsibilities](#2-wallet-roles-and-responsibilities)
3. [Improvements for Modern Key Technologies](#3-improvements-for-modern-key-technologies)
4. [Flow Improvements](#4-flow-improvements)
5. [Implementation Priority](#5-implementation-priority)

---

## 1. Current Flow Overview

### Current Architecture

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│ Watch Wallet│      │Keygen Wallet│      │ Sign Wallet │
│  (Online)   │      │  (Offline)  │      │  (Offline)  │
└─────────────┘      └─────────────┘      └─────────────┘
      │                    │                    │
      │ 1. Create TX       │                    │
      │───────────────────>│                    │
      │                    │                    │
      │                    │ 2. Sign (1st)      │
      │                    │<───────────────────│
      │                    │                    │
      │ 3. Partially Signed│                    │
      │<───────────────────│                    │
      │                    │                    │
      │                    │                    │ 4. Sign (2nd+)
      │                    │───────────────────>│
      │                    │                    │
      │ 5. Fully Signed     │                    │
      │<───────────────────│                    │
      │                    │                    │
      │ 6. Broadcast        │                    │
      │───────────────────>│                    │
```

### Current Transaction Flow

#### Single-Signature Address Flow

1. **Watch Wallet**: Create unsigned transaction
2. **Keygen Wallet**: Sign transaction (single signature)
3. **Watch Wallet**: Broadcast signed transaction

#### Multisig Address Flow (Current: P2SH/P2WSH)

1. **Watch Wallet**: Create unsigned transaction
2. **Keygen Wallet**: First signature
3. **Sign Wallet #1**: Second signature
4. **Sign Wallet #2**: Third signature (if required)
5. **Watch Wallet**: Broadcast fully signed transaction

---

## 2. Wallet Roles and Responsibilities

### 2.1 Watch Wallet (Online)

**Current Role:**
- Monitor addresses and UTXOs
- Create unsigned transactions
- Broadcast signed transactions
- Track transaction status

**Responsibilities:**
- ✅ Blockchain connectivity
- ✅ UTXO management
- ✅ Transaction construction
- ✅ Transaction broadcasting
- ✅ Balance monitoring

**Improvements Needed:**
- ⚠️ Taproot address support
- ⚠️ Descriptor wallet integration
- ⚠️ PSBT (Partially Signed Bitcoin Transaction) support
- ⚠️ Transaction fee estimation improvements

### 2.2 Keygen Wallet (Offline)

**Current Role:**
- Generate seeds and HD keys
- Create multisig addresses
- First signature for multisig transactions
- Single signature for non-multisig transactions

**Responsibilities:**
- ✅ Seed generation
- ✅ HD key derivation (BIP32/BIP44)
- ✅ Multisig address creation
- ✅ First signature in multisig flow
- ✅ Private key management (offline)

**Improvements Needed:**
- ⚠️ Taproot key derivation (BIP86)
- ⚠️ MuSig2 support
- ⚠️ PSBT handling
- ⚠️ Schnorr signature support
- ⚠️ Descriptor generation

### 2.3 Sign Wallet (Offline)

**Current Role:**
- Generate authorization account keys
- Export full public keys
- Second and subsequent signatures for multisig

**Responsibilities:**
- ✅ Authorization key generation
- ✅ Public key export
- ✅ Additional signatures for multisig
- ✅ Private key management (offline)

**Improvements Needed:**
- ⚠️ MuSig2 participation
- ⚠️ PSBT handling
- ⚠️ Schnorr signature support
- ⚠️ Parallel signing capability

---

## 3. Improvements for Modern Key Technologies

### 3.1 Taproot Support

#### Current State

- **Address Types Supported**: P2PKH, P2SH-SegWit, Bech32 (P2WPKH)
- **Taproot (P2TR)**: ❌ Not supported

#### Improvements

**1. Taproot Address Generation**

```go
// Keygen Wallet: Generate Taproot addresses
// Derivation path: m/86'/0'/account'/0/index
func (k *HDKey) CreateTaprootKey(
    seed []byte,
    accountType domainAccount.AccountType,
    idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
    // BIP86 derivation
    // Generate Taproot addresses (bc1p...)
}
```

**2. Taproot Transaction Support**

- Watch Wallet: Create Taproot transactions
- Keygen Wallet: Sign Taproot transactions (Schnorr signatures)
- Sign Wallet: Sign Taproot multisig transactions

**3. Benefits**

- ✅ Smaller transaction size
- ✅ Lower fees
- ✅ Better privacy (indistinguishable from single-sig)
- ✅ Script path flexibility

### 3.2 MuSig2 for Multisig

#### Current State

- **Multisig Type**: Traditional P2SH/P2WSH multisig
- **MuSig2**: ❌ Not supported

#### Improvements

**1. MuSig2 Protocol Implementation**

MuSig2 enables aggregated signatures that look like single signatures on-chain.

**Flow Changes:**

```
Current Flow (P2WSH Multisig):
1. Watch Wallet: Create unsigned TX
2. Keygen Wallet: Sign (ECDSA signature #1)
3. Sign Wallet #1: Sign (ECDSA signature #2)
4. Sign Wallet #2: Sign (ECDSA signature #3)
5. Watch Wallet: Broadcast (large transaction with 3 signatures)

Improved Flow (MuSig2):
1. Watch Wallet: Create unsigned TX
2. Keygen Wallet: Generate nonce (Round 1)
3. Sign Wallet #1: Generate nonce (Round 1)
4. Sign Wallet #2: Generate nonce (Round 1)
5. Keygen Wallet: Sign (Round 2)
6. Sign Wallet #1: Sign (Round 2)
7. Sign Wallet #2: Sign (Round 2)
8. Watch Wallet: Aggregate signatures and broadcast (looks like single sig)
```

**2. Implementation Requirements**

- **Keygen Wallet**: MuSig2 coordinator, nonce generation, signature aggregation
- **Sign Wallets**: Nonce generation, partial signature contribution
- **Watch Wallet**: Signature aggregation, transaction broadcasting

**3. Benefits**

- ✅ Transaction size reduction (single signature size vs. multiple)
- ✅ Privacy improvement (indistinguishable from single-sig)
- ✅ Lower fees
- ✅ Better scalability

### 3.3 PSBT (Partially Signed Bitcoin Transaction)

#### Current State

- **Transaction Format**: Custom JSON format
- **PSBT**: ❌ Not supported

#### Improvements

**1. PSBT Standard Adoption**

PSBT (BIP174) is a standard format for partially signed transactions.

**Benefits:**

- ✅ Standardized format (compatible with other wallets)
- ✅ Better security (structured data)
- ✅ Easier integration with hardware wallets
- ✅ Better error handling

**2. Flow Changes**

```
Current Flow:
Watch Wallet → JSON file → Keygen Wallet → JSON file → Sign Wallet → JSON file → Watch Wallet

Improved Flow (PSBT):
Watch Wallet → PSBT → Keygen Wallet → PSBT → Sign Wallet → PSBT → Watch Wallet
```

**3. Implementation**

- **Watch Wallet**: Create PSBT, finalize PSBT, extract transaction
- **Keygen Wallet**: Sign PSBT
- **Sign Wallet**: Sign PSBT

### 3.4 Descriptor Wallets

#### Current State

- **Wallet Type**: Legacy wallet format
- **Descriptor Wallets**: ❌ Not supported

#### Improvements

**1. Descriptor Generation**

Generate descriptors for all address types:

```go
// Taproot descriptor
tr([fingerprint/h/d]xpub.../0/*)

// Multisig descriptor (MuSig2)
tr([fingerprint/h/d]xpub1...,[fingerprint/h/d]xpub2.../0/*)

// Traditional multisig descriptor
wsh(sortedmulti(2,xpub1...,xpub2...))
```

**2. Benefits**

- ✅ Clear wallet functionality description
- ✅ Better compatibility with Bitcoin Core
- ✅ Easier backup and recovery
- ✅ More flexible script support

---

## 4. Flow Improvements

### 4.1 Improved Single-Signature Flow (Taproot)

```
1. Watch Wallet: Create unsigned Taproot transaction
   └─> watch create payment --address-type taproot

2. Keygen Wallet: Sign Taproot transaction (Schnorr signature)
   └─> keygen sign --file tx.psbt

3. Watch Wallet: Broadcast signed transaction
   └─> watch send --file tx_signed.psbt
```

**Benefits:**
- Smaller transaction size
- Lower fees
- Better privacy

### 4.2 Improved Multisig Flow (MuSig2)

```
1. Watch Wallet: Create unsigned transaction (PSBT)
   └─> watch create payment --multisig-type musig2

2. Keygen Wallet: Round 1 - Generate nonce
   └─> keygen musig2 nonce --file tx.psbt

3. Sign Wallet #1: Round 1 - Generate nonce
   └─> sign musig2 nonce --file tx.psbt

4. Sign Wallet #2: Round 1 - Generate nonce
   └─> sign musig2 nonce --file tx.psbt

5. Keygen Wallet: Round 2 - Sign
   └─> keygen musig2 sign --file tx.psbt

6. Sign Wallet #1: Round 2 - Sign
   └─> sign musig2 sign --file tx.psbt

7. Sign Wallet #2: Round 2 - Sign
   └─> sign musig2 sign --file tx.psbt

8. Watch Wallet: Aggregate signatures and broadcast
   └─> watch send --file tx_final.psbt
```

**Benefits:**
- Single signature on-chain (privacy)
- Smaller transaction size
- Lower fees

### 4.3 Parallel Signing Support

#### Current Limitation

Signatures must be collected sequentially (Keygen → Sign #1 → Sign #2).

#### Improvement

With MuSig2 Round 1, nonces can be generated in parallel:

```
Parallel Nonce Generation:
┌─────────────┐
│Keygen Wallet│──┐
└─────────────┘  │
                 ├─> All nonces generated in parallel
┌─────────────┐  │
│Sign Wallet #1│─┤
└─────────────┘  │
                 │
┌─────────────┐  │
│Sign Wallet #2│─┘
└─────────────┘
```

**Benefits:**
- Faster signature collection
- Better user experience
- Reduced waiting time

### 4.4 Transaction Type Improvements

#### Deposit Transaction

**Current:** P2SH-SegWit or Bech32 addresses

**Improved:** Taproot addresses

```bash
# Generate Taproot addresses for deposit account
keygen --coin btc create hdkey --account deposit --address-type taproot --keynum 10

# Create deposit transaction using Taproot addresses
watch --coin btc create deposit --address-type taproot --fee 0.0001
```

#### Payment Transaction

**Current:** Multisig P2WSH (large transaction size)

**Improved:** MuSig2 Taproot (single signature size)

```bash
# Create MuSig2 Taproot multisig address
keygen --coin btc create multisig --account payment --multisig-type musig2

# Create payment transaction
watch --coin btc create payment --multisig-type musig2 --fee 0.0001
```

#### Transfer Transaction

**Current:** Same as payment (multisig)

**Improved:** MuSig2 Taproot for better efficiency

---

## 5. Implementation Priority

### High Priority (Immediate Implementation)

1. **Taproot Address Support**
   - Keygen Wallet: BIP86 key derivation
   - Watch Wallet: Taproot transaction creation
   - Keygen/Sign Wallets: Schnorr signature support

2. **PSBT Support**
   - Replace custom JSON format with PSBT
   - Better compatibility and security

3. **Error Handling Improvements**
   - Better error messages
   - Transaction validation

### Medium Priority (Near Future)

4. **MuSig2 Implementation**
   - Two-round signature protocol
   - Signature aggregation
   - Parallel nonce generation

5. **Descriptor Wallets**
   - Descriptor generation
   - Bitcoin Core integration

6. **Transaction Fee Optimization**
   - Better fee estimation
   - RBF (Replace-By-Fee) support
   - CPFP (Child-Pays-For-Parent) support

### Low Priority (Long-term)

7. **Hardware Wallet Integration**
   - Ledger/Trezor support via PSBT
   - Better security for key management

8. **Advanced Script Support**
   - Time-locked transactions
   - Conditional payments
   - Complex multisig scenarios

---

## 6. Migration Strategy

### Phase 1: Taproot Support (3-6 months)

1. Implement BIP86 key derivation
2. Add Taproot address generation
3. Implement Schnorr signatures
4. Update transaction creation to support Taproot

### Phase 2: PSBT Adoption (3-6 months)

1. Replace custom JSON with PSBT
2. Update all wallet types to handle PSBT
3. Add PSBT validation
4. Update documentation

### Phase 3: MuSig2 Implementation (6-12 months)

1. Implement MuSig2 protocol
2. Add nonce generation and management
3. Implement signature aggregation
4. Update multisig flow to use MuSig2

### Phase 4: Descriptor Wallets (6-12 months)

1. Generate descriptors for all address types
2. Integrate with Bitcoin Core
3. Update backup/recovery procedures

---

## 7. Security Considerations

### 7.1 MuSig2 Security

- **Nonce Management**: Nonces must be unique and never reused
- **Key Aggregation**: Verify aggregated public keys
- **Signature Verification**: Validate all partial signatures

### 7.2 Taproot Security

- **Key Path vs Script Path**: Ensure proper key path spending
- **Script Path Security**: Verify script conditions
- **Address Validation**: Validate Taproot addresses before use

### 7.3 PSBT Security

- **Input Validation**: Verify all PSBT inputs
- **Signature Verification**: Validate all signatures
- **Finalization**: Ensure proper PSBT finalization

---

## 8. Testing Strategy

### 8.1 Unit Tests

- Taproot key derivation
- Schnorr signature generation
- MuSig2 protocol steps
- PSBT creation and parsing

### 8.2 Integration Tests

- End-to-end Taproot transaction flow
- MuSig2 multisig flow
- PSBT signing flow
- Descriptor generation and import

### 8.3 Testnet Testing

- Test all improvements on Bitcoin testnet/signet
- Verify transaction sizes and fees
- Test with multiple signers
- Validate privacy properties

---

## 9. Summary

### Key Improvements

1. **Taproot Support**: Modern address format with better privacy and lower fees
2. **MuSig2**: Efficient multisig that looks like single-sig on-chain
3. **PSBT**: Standardized transaction format for better compatibility
4. **Descriptor Wallets**: Clear wallet functionality description
5. **Parallel Signing**: Faster signature collection with MuSig2

### Expected Benefits

- **Transaction Size**: 30-50% reduction for multisig transactions
- **Fees**: 30-50% reduction due to smaller transaction size
- **Privacy**: Improved (Taproot and MuSig2 transactions look like single-sig)
- **Compatibility**: Better integration with Bitcoin Core and other wallets
- **User Experience**: Faster transaction processing with parallel signing

### Implementation Timeline

- **Phase 1 (Taproot)**: 3-6 months
- **Phase 2 (PSBT)**: 3-6 months
- **Phase 3 (MuSig2)**: 6-12 months
- **Phase 4 (Descriptors)**: 6-12 months

**Total**: 18-36 months for complete implementation

---

## References

- [BIP 174: Partially Signed Bitcoin Transaction Format](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)
- [BIP 340: Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [BIP 341: Taproot](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki)
- [BIP 86: Key Derivation for Single Key Taproot Outputs](https://github.com/bitcoin/bips/blob/master/bip-0086.mediawiki)
- [MuSig2: Simple Two-Round Schnorr Multisignatures](https://eprint.iacr.org/2020/1261)
- [Bitcoin Core: Descriptors](https://github.com/bitcoin/bitcoin/blob/master/doc/descriptors.md)

