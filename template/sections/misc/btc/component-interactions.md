## Component Interactions

### Address Creation Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CLI Command                                              │
│    keygen create musig2-address --account payment           │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. CreateMuSig2AddressUseCase                               │
│    - Validate account is multisig                           │
│    - Get signer public keys                                 │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. AuthFullPubkeyRepository                                 │
│    - GetOne(auth1) → pubKey1                                │
│    - GetOne(auth2) → pubKey2                                │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. AccountKeyRepository                                     │
│    - GetAllAddrStatus(payment, PrivKeyImported)             │
│    → Returns account keys needing multisig addresses        │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 5. MuSig2Service.AggregatePublicKeys()                      │
│    Input: [pubKey1, pubKey2, accountPubKey]                 │
│    - Sort public keys                                       │
│    - Aggregate using MuSig2 protocol                        │
│    - Apply Taproot tweak (BIP86)                            │
│    Output: aggregatedPubKey                                 │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 6. Create Taproot Address                                   │
│    - Create P2TR address from aggregatedPubKey              │
│    - Format: bc1p... (mainnet) or tb1p... (testnet)        │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 7. AccountKeyRepository.UpdateMultisigAddr()                │
│    - Store P2TR address                                     │
│    - Update addr_status                                     │
└─────────────────────────────────────────────────────────────┘
```

### Transaction Signing Flow (Two-Round Protocol)

#### Round 1: Nonce Generation (Parallel)

```
┌─────────────────────────────────────────────────────────────┐
│ KEYGEN WALLET                                               │
│ 1. keygen musig2 nonce --file payment_15_unsigned_0.psbt   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. GenerateMuSig2NonceUseCase (Keygen)                     │
│    - Parse PSBT                                             │
│    - Validate transaction                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. MuSig2Service.GenerateNonce()                            │
│    - Generate secure random nonce                           │
│    - Create public nonce (66 bytes)                         │
│    Output: [66]byte nonce                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. Store in PSBT                                            │
│    - Add nonce to PSBT proprietary field                    │
│    - Keygen nonce: field_id = "musig2_nonce_keygen"        │
│    Output: payment_15_unsigned_0_...1.psbt                  │
└─────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────┐
│ SIGN WALLET 1 (Parallel)                                    │
│ 1. sign musig2 nonce --file payment_15_unsigned_0_...1.psbt│
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. GenerateMuSig2NonceUseCase (Sign)                       │
│    - Parse PSBT                                             │
│    - Validate transaction                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. MuSig2Service.GenerateNonce()                            │
│    - Generate secure random nonce                           │
│    Output: [66]byte nonce                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. Store in PSBT                                            │
│    - Add nonce to PSBT proprietary field                    │
│    - Sign nonce: field_id = "musig2_nonce_sign1"           │
│    Output: payment_15_unsigned_0_...2.psbt                  │
└─────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────┐
│ SIGN WALLET 2 (Parallel)                                    │
│ 1. sign musig2 nonce --file payment_15_unsigned_0_...2.psbt│
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. GenerateMuSig2NonceUseCase (Sign)                       │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. MuSig2Service.GenerateNonce()                            │
│    Output: [66]byte nonce                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. Store in PSBT                                            │
│    - Sign nonce: field_id = "musig2_nonce_sign2"           │
│    Output: payment_15_nonce_0_...3.psbt                     │
│    Status: All nonces collected (3/3)                       │
└─────────────────────────────────────────────────────────────┘
```

#### Round 2: Signing (Sequential)

```
┌─────────────────────────────────────────────────────────────┐
│ KEYGEN WALLET                                               │
│ 1. keygen musig2 sign --file payment_15_nonce_0.psbt       │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. MuSig2SignUseCase (Keygen)                              │
│    - Validate all nonces present                            │
│    - Extract nonces from PSBT                               │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. AccountKeyRepository                                     │
│    - Get private key for signing                            │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. MuSig2Service.CreatePartialSignature()                  │
│    Input:                                                    │
│    - privateKey (keygen)                                    │
│    - allPublicKeys                                          │
│    - allNonces [nonce1, nonce2, nonce3]                     │
│    - messageHash (transaction hash)                         │
│    Process:                                                  │
│    - Create MuSig2 context with all signers                 │
│    - Register all nonces                                    │
│    - Sign message hash                                      │
│    Output: partialSignature (32 bytes)                      │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 5. Store in PSBT                                            │
│    - Add partial signature to PSBT proprietary field        │
│    - Field ID: "musig2_partialsig_keygen"                  │
│    Output: payment_15_unsigned_1.psbt                       │
└─────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────┐
│ SIGN WALLET 1 (Sequential after Keygen)                    │
│ 1. sign musig2 sign --file payment_15_unsigned_1.psbt      │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. MuSig2SignUseCase (Sign)                                │
│    - Validate all nonces present                            │
│    - Extract nonces from PSBT                               │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. AuthAccountKeyRepository                                 │
│    - Get auth private key for signing                       │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. MuSig2Service.CreatePartialSignature()                  │
│    Input:                                                    │
│    - privateKey (sign1/auth1)                               │
│    - allPublicKeys                                          │
│    - allNonces                                              │
│    - messageHash                                            │
│    Output: partialSignature (32 bytes)                      │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 5. Store in PSBT                                            │
│    - Field ID: "musig2_partialsig_sign1"                   │
│    Output: payment_15_unsigned_2.psbt                       │
└─────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────┐
│ SIGN WALLET 2 (Sequential after Sign1)                     │
│ 1. sign musig2 sign --file payment_15_unsigned_2.psbt      │
└────────────────────┬────────────────────────────────────────┘
                     │
                 [Same process as Sign1]
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 5. Store in PSBT                                            │
│    - Field ID: "musig2_partialsig_sign2"                   │
│    Output: payment_15_unsigned_3.psbt                       │
│    Status: All partial signatures collected (3/3)           │
└─────────────────────────────────────────────────────────────┘
```

#### Aggregation (Watch Wallet)

```
┌─────────────────────────────────────────────────────────────┐
│ WATCH WALLET                                                │
│ 1. watch musig2 aggregate --file payment_15_unsigned_3.psbt│
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. AggregateMuSig2SignaturesUseCase                        │
│    - Validate all partial signatures present                │
│    - Extract partial signatures from PSBT                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. MuSig2Service.AggregateSignatures()                     │
│    Input:                                                    │
│    - allPublicKeys                                          │
│    - allNonces                                              │
│    - partialSignatures [sig1, sig2, sig3]                  │
│    - messageHash                                            │
│    Process:                                                  │
│    - Create MuSig2 context                                  │
│    - Register all nonces                                    │
│    - Combine all partial signatures                         │
│    - Generate final signature                               │
│    Output: finalSignature (64 bytes Schnorr)               │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. MuSig2Service.VerifyAggregatedSignature()               │
│    - Verify signature against aggregated public key        │
│    - Ensure signature is valid before broadcasting          │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 5. Finalize PSBT                                            │
│    - Add final signature to PSBT witness                    │
│    - Remove temporary MuSig2 fields                         │
│    Output: payment_15_signed_3.psbt                         │
│    Status: Ready for broadcasting                           │
└─────────────────────────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 6. Broadcast Transaction                                    │
│    watch send --file payment_15_signed_3.psbt              │
│    → Transaction broadcast to Bitcoin network               │
└─────────────────────────────────────────────────────────────┘
```

---
