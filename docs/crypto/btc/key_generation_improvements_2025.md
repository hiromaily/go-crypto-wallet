# Bitcoin Key Generation Modernization Improvements (End of 2025)

This document summarizes improvements for modernizing Bitcoin key generation as of the end of 2025.

## Table of Contents

1. [Taproot (BIP341/BIP86) Support](#1-taproot-bip341bip86-support)
2. [Complete BIP49 (P2WPKH-P2SH) Implementation](#2-complete-bip49-p2wpkh-p2sh-implementation)
3. [BIP85 (Deterministic Entropy) Consideration](#3-bip85-deterministic-entropy-consideration)
4. [Descriptor Wallets Support](#4-descriptor-wallets-support)
5. [Multisig Improvements with MuSig2](#5-multisig-improvements-with-musig2)
6. [Random Number Generation Enhancement](#6-random-number-generation-enhancement)
7. [Extended BIP32/BIP44 Support](#7-extended-bip32bip44-support)
8. [Security Enhancements](#8-security-enhancements)
9. [Implementation Priority](#9-implementation-priority)

---

## 1. Taproot (BIP341/BIP86) Support

### Current Status

The current implementation only supports the following address formats:

- P2PKH (Legacy)
- P2SH-SegWit (P2WPKH-P2SH)
- Bech32 (Native SegWit, P2WPKH)

**Taproot addresses (P2TR) are not supported**.

### Improvements

Taproot was activated on the Bitcoin network in November 2021 and has become a standard address format as of 2025.

**What to implement:**

1. **BIP86 (Taproot Key Path Spending) Support**
   - Generate Taproot addresses (`bc1p...`)
   - BIP32 derivation path: `m/86'/0'/0'/0/0` (BIP86 purpose)
   - Or generate Taproot addresses from existing BIP44 paths

2. **Taproot Signature Support**
   - Implement Schnorr signatures (BIP340)
   - Create and sign Taproot transactions

3. **Integration with Existing Code**

   ```go
   // Add to internal/infrastructure/wallet/key/hd_wallet.go
   // Generate Taproot address
   func (k *HDKey) getTaprootAddr(privKey *btcec.PrivateKey) (*btcutil.AddressTaproot, error) {
       // BIP340 Schnorr public key generation
       // BIP341 Taproot output creation
   }
   ```

4. **Domain Model Extension**

   ```go
   // Add to internal/domain/key/valueobject.go
   type WalletKey struct {
       // ... existing fields
       TaprootAddr string // Taproot address (bc1p...)
   }
   ```

### References

- [BIP 340: Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [BIP 341: Taproot](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki)
- [BIP 86: Key Derivation for Single Key Taproot Outputs](https://github.com/bitcoin/bips/blob/master/bip-0086.mediawiki)

---

## 2. Complete BIP49 (P2WPKH-P2SH) Implementation

### Current Status

`PurposeTypeBIP49` is defined in the code, but its actual usage has not been confirmed.

```go
// internal/infrastructure/wallet/key/hd_wallet.go
const (
    PurposeTypeBIP44 PurposeType = 44 // BIP44
    PurposeTypeBIP49 PurposeType = 49 // BIP49
)
```

### Improvements

BIP49 is a format that wraps P2WPKH in P2SH, providing SegWit benefits while maintaining compatibility with legacy wallets.

**What to implement:**

1. **BIP49 Derivation Path Support**
   - Path: `m/49'/0'/0'/0/0`
   - Generate P2SH-SegWit addresses (already implemented, but should be explicitly supported as BIP49 path)

2. **Purpose Type Selection Feature**
   - Allow users to select BIP44/BIP49/BIP86
   - Make it configurable in configuration files

---

## 3. BIP85 (Deterministic Entropy) Consideration

### Current Status

Currently, seeds are generated directly from BIP39 mnemonics.

```go
// internal/infrastructure/wallet/key/seed.go
func GenerateMnemonic(passphrase string) ([]byte, string, error) {
    entropy, _ := bip39.NewEntropy(256)
    mnemonic, err := bip39.NewMnemonic(entropy)
    // ...
}
```

### Improvements

BIP85 provides a method to deterministically derive entropy from an existing BIP32 seed. This enables:

- Generate independent entropy for multiple applications from a single master seed
- More secure key management
- Simplified backup process

**What to implement:**

1. **BIP85 Entropy Derivation Implementation**

   ```go
   // BIP85: Deterministic Entropy From BIP32 Seed
   func DeriveBIP85Entropy(masterSeed []byte, applicationIndex uint32, entropyBits uint32) ([]byte, error) {
       // BIP85 derivation logic
   }
   ```

2. **Application-Specific Entropy Generation**
   - Generate different entropy for each application
   - More secure key management

### References

- [BIP 85: Deterministic Entropy From BIP32 Seed](https://github.com/bitcoin/bips/blob/master/bip-0085.mediawiki)

---

## 4. Descriptor Wallets Support

### Current Status

Bitcoin Core has recommended Descriptor Wallets since 2020, but the current implementation uses the traditional wallet format.

### Improvements

Descriptor Wallets are a new format that expresses wallet functionality using descriptors.

**Benefits:**

- More flexible script support
- Clear description of wallet functionality
- Easier multisig management

**What to implement:**

1. **Descriptor Generation**

   ```go
   // Taproot descriptor example
   // tr([fingerprint/h/d]xpub.../0/*)
   
   // Multisig descriptor example
   // wsh(sortedmulti(2,xpub1...,xpub2...))
   ```

2. **Integration with Bitcoin Core**
   - Use `importdescriptors` RPC
   - Generate descriptors when creating wallets

### References

- [Bitcoin Core: Descriptors](https://github.com/bitcoin/bitcoin/blob/master/doc/descriptors.md)

---

## 5. Multisig Improvements with MuSig2

### Current Status

Currently using traditional multisig (P2SH/P2WSH).

### Improvements

MuSig2 is a Schnorr signature-based aggregate signature protocol that significantly improves multisig efficiency.

**Benefits:**

- Reduced transaction size
- Improved privacy (indistinguishable from regular single signatures)
- Efficiency through signature aggregation

**What to implement:**

1. **MuSig2 Protocol Implementation**
   - Two-round signature protocol
   - Signature aggregation

2. **Integration with Taproot Multisig**
   - Use MuSig2 with Taproot script paths
   - More efficient multisig transactions

### References

- [MuSig2: Simple Two-Round Schnorr Multisignatures](https://eprint.iacr.org/2020/1261)

---

## 6. Random Number Generation Enhancement

### Current Status

`hdkeychain.GenerateSeed()` and `bip39.NewEntropy()` are used, but internal implementation verification is needed.

### Improvements

1. **Explicit Verification of crypto/rand Usage**
   - Verify that `crypto/rand` is being used
   - Verify that the system's random number generator is properly initialized

2. **Entropy Source Verification**
   - Quality checks for entropy
   - Verification in tests

3. **Enhanced Error Handling**

   ```go
   // Current code
   entropy, _ := bip39.NewEntropy(256) // Error is ignored
   
   // Improved version
   entropy, err := bip39.NewEntropy(256)
   if err != nil {
       return nil, "", fmt.Errorf("failed to generate entropy: %w", err)
   }
   ```

---

## 7. Extended BIP32/BIP44 Support

### Current Status

Only BIP44 is implemented, and support for BIP49, BIP84, and BIP86 is incomplete.

### Improvements

1. **Complete Purpose Type Support**
   - BIP44 (Legacy): `m/44'/0'/0'/0/0`
   - BIP49 (P2SH-SegWit): `m/49'/0'/0'/0/0`
   - BIP84 (Native SegWit): `m/84'/0'/0'/0/0` (already implemented as Bech32)
   - BIP86 (Taproot): `m/86'/0'/0'/0/0`

2. **Selection via Configuration**
   - Allow users to select Purpose Type according to their needs
   - Recommend Taproot (BIP86) as default

---

## 8. Security Enhancements

### Improvements

1. **Memory Clear Implementation**
   - Explicitly clear private keys from memory
   - Implement `memset` equivalent functionality

2. **Key Derivation Path Validation**
   - Detect invalid derivation paths
   - Verify hardening

3. **Entropy Verification**
   - Quality checks for generated entropy
   - Detect weak entropy

4. **Exclusion of Secret Information from Logs**
   - Verify that private keys, seeds, and mnemonics are not logged
   - Likely already implemented, but should be reconfirmed

---

## 9. Implementation Priority

### High Priority (Should be implemented immediately)

1. **Taproot (BIP341/BIP86) Support**
   - Standard address format as of 2025
   - Already supported by the existing btcd library

2. **Complete BIP49 Implementation**
   - Defined in code but unused
   - Can leverage existing P2SH-SegWit implementation

3. **Error Handling Improvements**
   - Fix locations where `bip39.NewEntropy()` errors are ignored

### Medium Priority (To be implemented in the near future)

1. **Descriptor Wallets Support**
   - Improved compatibility with Bitcoin Core
   - More flexible script support

2. **BIP85 Consideration**
   - More secure key management
   - Consider implementation complexity

### Low Priority (Long-term improvements)

1. **MuSig2 Implementation**
   - Multisig efficiency improvements
   - High implementation complexity

2. **Quantum Resistance Consideration**
   - Not practical at this time
   - Long-term research topic

---

## Implementation Examples

### Taproot Address Generation Example

```go
// Add to internal/infrastructure/wallet/key/hd_wallet.go

import (
    "github.com/btcsuite/btcd/btcec/v2"
    "github.com/btcsuite/btcd/btcutil"
    "github.com/btcsuite/btcd/btcutil/hdkeychain"
    "github.com/btcsuite/btcd/chaincfg"
    "github.com/btcsuite/btcd/txscript"
)

// getTaprootAddr returns Taproot address (BIP86)
func (k *HDKey) getTaprootAddr(privKey *btcec.PrivateKey) (*btcutil.AddressTaproot, error) {
    // BIP340: Schnorr public key generation
    pubKey := privKey.PubKey()
    
    // BIP341: Taproot output creation
    // Taproot uses 32-byte public keys
    taprootKey := txscript.ComputeTaprootKeyNoScript(pubKey)
    
    // Generate Taproot address
    taprootAddr, err := btcutil.NewAddressTaproot(
        schnorr.SerializePubKey(taprootKey),
        k.conf,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create taproot address: %w", err)
    }
    
    return taprootAddr, nil
}
```

### BIP86 Derivation Path Implementation Example

```go
// Add BIP86 to PurposeType
const (
    PurposeTypeBIP44 PurposeType = 44 // BIP44
    PurposeTypeBIP49 PurposeType = 49 // BIP49
    PurposeTypeBIP84 PurposeType = 84 // BIP84 (Native SegWit)
    PurposeTypeBIP86 PurposeType = 86 // BIP86 (Taproot)
)
```

---

## References

### BIPs

- [BIP 32: Hierarchical Deterministic Wallets](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)
- [BIP 39: Mnemonic Code for generating deterministic keys](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)
- [BIP 44: Multi-Account Hierarchy for Deterministic Wallets](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)
- [BIP 49: Derivation scheme for P2WPKH-nested-in-P2SH](https://github.com/bitcoin/bips/blob/master/bip-0049.mediawiki)
- [BIP 84: Derivation scheme for P2WPKH based accounts](https://github.com/bitcoin/bips/blob/master/bip-0084.mediawiki)
- [BIP 85: Deterministic Entropy From BIP32 Seed](https://github.com/bitcoin/bips/blob/master/bip-0085.mediawiki)
- [BIP 86: Key Derivation for Single Key Taproot Outputs](https://github.com/bitcoin/bips/blob/master/bip-0086.mediawiki)

### Libraries

- [btcd/btcutil](https://pkg.go.dev/github.com/btcsuite/btcd/btcutil) - Verify Taproot support
- [btcd/btcec/v2](https://pkg.go.dev/github.com/btcsuite/btcd/btcec/v2) - Schnorr signature support

---

## Summary

The most important improvements for modernizing Bitcoin key generation as of the end of 2025 are:

1. **Taproot (BIP86) Support** - Has become a standard address format
2. **Complete BIP49 Implementation** - Already defined in code but unused
3. **Error Handling Improvements** - Enhanced security and robustness

These improvements will enable compliance with the latest Bitcoin standards and more secure and efficient key management.

