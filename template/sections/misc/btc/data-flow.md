## Data Flow

### Nonce Data Flow

```
Round 1: Nonce Generation
=========================

Keygen Wallet:
  Private Key → MuSig2Service.GenerateNonce()
                       ↓
                [66]byte Public Nonce
                       ↓
              PSBT Proprietary Field
                       ↓
         payment_15_unsigned_0_...1.psbt

Sign Wallet 1:
  Private Key → MuSig2Service.GenerateNonce()
                       ↓
                [66]byte Public Nonce
                       ↓
              PSBT Proprietary Field
                       ↓
         payment_15_unsigned_0_...2.psbt

Sign Wallet 2:
  Private Key → MuSig2Service.GenerateNonce()
                       ↓
                [66]byte Public Nonce
                       ↓
              PSBT Proprietary Field
                       ↓
          payment_15_nonce_0_...3.psbt
          (All nonces collected)
```

### Signature Data Flow

```
Round 2: Partial Signature Creation
====================================

Each Signer:
  Private Key + All Nonces + Message Hash
              ↓
  MuSig2Service.CreatePartialSignature()
              ↓
  [32]byte Partial Signature
              ↓
  PSBT Proprietary Field
              ↓
  Updated PSBT file

Aggregation:
  All Partial Signatures + All Nonces + Message Hash
              ↓
  MuSig2Service.AggregateSignatures()
              ↓
  [64]byte Schnorr Signature
              ↓
  PSBT Witness Field
              ↓
  Finalized PSBT
              ↓
  Broadcast to Bitcoin Network
```

### PSBT Field Structure

MuSig2 uses PSBT proprietary fields (BIP174 extension):

```
Proprietary Field Format:
Key:   <identifier> <subtype> <keydata>
Value: <valuedata>

MuSig2 Nonces:
Key:   FC <musig2> <nonce> <signer_id>
Value: [66 bytes] public nonce

MuSig2 Partial Signatures:
Key:   FC <musig2> <psig> <signer_id>
Value: [32 bytes] partial signature

Example PSBT with MuSig2 data:
┌─────────────────────────────────────┐
│ PSBT Global Fields                  │
│ - Version                           │
│ - Transaction                       │
└─────────────────────────────────────┘
┌─────────────────────────────────────┐
│ PSBT Input #0                       │
│ - Non-witness UTXO                  │
│ - Witness UTXO                      │
│ - Derivation paths                  │
│ - Proprietary Fields:               │
│   * musig2_nonce_keygen: [66]byte  │
│   * musig2_nonce_sign1:  [66]byte  │
│   * musig2_nonce_sign2:  [66]byte  │
│   * musig2_psig_keygen:  [32]byte  │
│   * musig2_psig_sign1:   [32]byte  │
│   * musig2_psig_sign2:   [32]byte  │
└─────────────────────────────────────┘
┌─────────────────────────────────────┐
│ PSBT Output #0 (Recipient)          │
│ - Script                            │
│ - Amount                            │
└─────────────────────────────────────┘
┌─────────────────────────────────────┐
│ PSBT Output #1 (Change)             │
│ - Script                            │
│ - Amount                            │
└─────────────────────────────────────┘
```

---
