# MuSig2 Usage Guide

**Date:** 2025-12-29
**Issue:** #132 - Research and MuSig2 Library Selection
**Library:** `github.com/btcsuite/btcd/btcec/v2/schnorr/musig2`
**Version:** v2.3.6

## Quick Start

### Import

```go
import (
    "github.com/btcsuite/btcd/btcec/v2"
    "github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
)
```

### Basic MuSig2 Flow (High-Level API)

#### Step 1: Create Context (Each Signer)

```go
// Each signer creates their own context
ctx, err := musig2.NewContext(
    privateKey,           // Signer's private key
    true,                 // Sort keys (must be same across all signers)
    musig2.WithKnownSigners(allPublicKeys),  // All participants' public keys
)
if err != nil {
    return fmt.Errorf("failed to create context: %w", err)
}
```

#### Step 2: Create Session (Each Signer)

```go
session, err := ctx.NewSession()
if err != nil {
    return fmt.Errorf("failed to create session: %w", err)
}
```

#### Step 3: Round 1 - Nonce Generation (Parallel)

```go
// Each signer generates their public nonce
pubNonce := session.PublicNonce()  // Returns [66]byte

// Share pubNonce with all other signers
// (via file, PSBT, or network in production)
```

#### Step 4: Nonce Exchange (Each Signer)

```go
// Register public nonces from all other signers
for _, otherNonce := range otherSignersNonces {
    haveAllNonces, err := session.RegisterPubNonce(otherNonce)
    if err != nil {
        return fmt.Errorf("failed to register nonce: %w", err)
    }
}
```

#### Step 5: Round 2 - Signing (Each Signer)

```go
// Sign the message hash (32 bytes)
messageHash := sha256.Sum256(transaction)

partialSig, err := session.Sign(messageHash)
if err != nil {
    return fmt.Errorf("failed to sign: %w", err)
}

// Share partialSig with coordinator (watch wallet)
```

#### Step 6: Signature Aggregation (Coordinator/Watch Wallet)

```go
// Coordinator collects all partial signatures
for _, partialSig := range otherPartialSigs {
    haveAllSigs, err := session.CombineSig(partialSig)
    if err != nil {
        return fmt.Errorf("failed to combine signature: %w", err)
    }

    if haveAllSigs {
        break  // All signatures received
    }
}

// Get final signature
finalSig := session.FinalSig()
if finalSig == nil {
    return fmt.Errorf("final signature not available")
}
```

#### Step 7: Verification

```go
// Get aggregated public key
aggregatedKey, err := ctx.CombinedKey()
if err != nil {
    return fmt.Errorf("failed to get combined key: %w", err)
}

// Verify signature
if !finalSig.Verify(messageHash[:], aggregatedKey) {
    return fmt.Errorf("invalid signature")
}
```

## Integration with Our Wallet Architecture

### Keygen Wallet (Offline)

```go
// Round 1: Generate nonce
ctx, err := musig2.NewContext(keygenPrivKey, true,
    musig2.WithKnownSigners(allPubKeys))
if err != nil {
    return fmt.Errorf("failed to create context: %w", err)
}

session, err := ctx.NewSession()
if err != nil {
    return fmt.Errorf("failed to create session: %w", err)
}

keygenNonce := session.PublicNonce()

// Save nonce to file (transfer to Watch wallet)
SaveNonceToFile("keygen_nonce.json", keygenNonce)

// Round 2: Sign after receiving all individual nonces
allNonces := LoadAllNoncesFromFile("all_nonces.json")
// Register each other signer's nonce
for _, otherNonce := range allNonces {
    if otherNonce != keygenNonce { // Skip own nonce
        _, err := session.RegisterPubNonce(otherNonce)
        if err != nil {
            // Handle error
        }
    }
}

partialSig, err := session.Sign(txHash)
if err != nil {
    return fmt.Errorf("failed to sign: %w", err)
}

// Save partial signature to file (transfer to Watch wallet)
SavePartialSigToFile("keygen_partialsig.json", partialSig)
```

### Sign Wallet (Offline)

```go
// Round 1: Generate nonce
ctx, err := musig2.NewContext(signPrivKey, true,
    musig2.WithKnownSigners(allPubKeys))
if err != nil {
    return fmt.Errorf("failed to create context: %w", err)
}

session, err := ctx.NewSession()
if err != nil {
    return fmt.Errorf("failed to create session: %w", err)
}

signNonce := session.PublicNonce()

// Save nonce to file
SaveNonceToFile("sign_nonce.json", signNonce)

// Round 2: Sign after receiving all individual nonces
allNonces := LoadAllNoncesFromFile("all_nonces.json")
// Register each other signer's nonce
for _, otherNonce := range allNonces {
    if otherNonce != signNonce { // Skip own nonce
        _, err := session.RegisterPubNonce(otherNonce)
        if err != nil {
            // Handle error
        }
    }
}

partialSig, err := session.Sign(txHash)
if err != nil {
    return fmt.Errorf("failed to sign: %w", err)
}

// Save partial signature to file
SavePartialSigToFile("sign_partialsig.json", partialSig)
```

### Watch Wallet (Online)

```go
// Round 1: Collect nonces from offline wallets
keygenNonce := LoadNonceFromFile("keygen_nonce.json")
signNonce := LoadNonceFromFile("sign_nonce.json")

// Bundle all individual nonces (Watch wallet acts as coordinator)
// DO NOT aggregate nonces - each signer needs to register individual nonces
allNonces := [][66]byte{keygenNonce, signNonce}

// Save the bundle of individual nonces for transfer to offline wallets
// Each offline wallet will register all other signers' individual nonces
SaveAllNoncesToFile("all_nonces.json", allNonces)

// Round 2: Aggregate partial signatures
keygenPartialSig := LoadPartialSigFromFile("keygen_partialsig.json")
signPartialSig := LoadPartialSigFromFile("sign_partialsig.json")

// Use low-level API to combine signatures
finalSig := musig2.CombineSigs(combinedNonce,
    []*musig2.PartialSignature{keygenPartialSig, signPartialSig})

// Add signature to transaction and broadcast
AddSignatureToTx(tx, finalSig)
BroadcastTx(tx)
```

## Taproot Integration

For Taproot addresses (compatible with Phase 1):

```go
// Create context with Taproot tweak
ctx, err := musig2.NewContext(
    privateKey,
    true,
    musig2.WithKnownSigners(allPublicKeys),
    musig2.WithTaprootTweakCtx(scriptRoot),  // Add Taproot tweak
)
```

Or for BIP 86 (key-path only):

```go
ctx, err := musig2.NewContext(
    privateKey,
    true,
    musig2.WithKnownSigners(allPublicKeys),
    musig2.WithBip86TweakCtx(),  // BIP 86 tweak
)
```

## PSBT Integration (Phase 2 Compatible)

MuSig2 data can be stored in PSBT extension fields:

### Store Public Nonces in PSBT

```go
// PSBT proprietary field for MuSig2 nonces
// Key: 0xFC (proprietary) + identifier + subtype
// Value: Public nonce [66]byte
func StoreMuSig2NonceInPSBT(psbt *psbt.Packet, signerID string, nonce [66]byte) {
    // Store in PSBT proprietary field
    // Implementation in sub-issue #136-#138
}
```

### Store Partial Signatures in PSBT

```go
// PSBT proprietary field for MuSig2 partial signatures
func StorePartialSigInPSBT(psbt *psbt.Packet, signerID string, sig *musig2.PartialSignature) {
    // Store S value and R point
    // Implementation in sub-issue #136-#138
}
```

## Important Types

### Nonces

```go
type Nonces struct {
    PubNonce [66]byte  // 2x 33-byte compressed points (shareable)
    SecNonce [97]byte  // 2x 32-byte scalars + metadata (KEEP SECRET)
}
```

**Security:**
- `PubNonce` is safe to share with all signers
- `SecNonce` must be kept SECRET and NEVER shared
- Nonces must be generated fresh for each transaction
- **NEVER reuse nonces** (will leak private key)

### Partial Signature

```go
type PartialSignature struct {
    S *btcec.ModNScalar  // Signature scalar
    R *btcec.PublicKey   // Nonce commitment
}
```

Access signature data:
```go
sBytes := partialSig.S.Bytes()  // [32]byte
```

### Final Signature

```go
*schnorr.Signature  // 64-byte Schnorr signature
```

Serialize:
```go
sigBytes := finalSig.Serialize()  // [64]byte
```

## Security Best Practices

### 1. Nonce Management

✅ **DO:**
- Generate nonces with secure RNG (default)
- Store secret nonces securely (encrypted database)
- Delete secret nonces after use
- Use unique nonces per transaction

❌ **DO NOT:**
- Reuse nonces (CRITICAL: leaks private key)
- Share secret nonces with anyone
- Log secret nonces
- Store nonces in plain text files

### 2. Key Aggregation

✅ **DO:**
- Use consistent key sorting across all signers
- Verify all public keys before aggregation
- Use Taproot tweaking for Taproot addresses

❌ **DO NOT:**
- Mix sorted and unsorted keys
- Skip key validation
- Ignore Taproot tweaks when required

### 3. Signature Verification

✅ **DO:**
- Always verify final signature before broadcasting
- Check that all partial signatures are received
- Validate aggregated public key

❌ **DO NOT:**
- Skip signature verification
- Broadcast unverified transactions
- Assume partial signatures are valid

## Error Handling Patterns

```go
// Wrap errors with context
ctx, err := musig2.NewContext(privKey, true, opts...)
if err != nil {
    return fmt.Errorf("failed to create MuSig2 context: %w", err)
}

// Check boolean returns
haveAllNonces, err := session.RegisterPubNonce(nonce)
if err != nil {
    return fmt.Errorf("failed to register nonce: %w", err)
}
if !haveAllNonces {
    // Still need more nonces
    return ErrInsufficientNonces
}

// Verify nil returns
finalSig := session.FinalSig()
if finalSig == nil {
    return fmt.Errorf("final signature not available")
}
```

## Testing

### Unit Test Example

```go
func TestMuSig2Signing(t *testing.T) {
    // Generate test keys
    privKey1, _ := btcec.NewPrivateKey()
    privKey2, _ := btcec.NewPrivateKey()
    pubKeys := []*btcec.PublicKey{privKey1.PubKey(), privKey2.PubKey()}

    // Create contexts
    ctx1, err := musig2.NewContext(privKey1, true,
        musig2.WithKnownSigners(pubKeys))
    require.NoError(t, err)

    ctx2, err := musig2.NewContext(privKey2, true,
        musig2.WithKnownSigners(pubKeys))
    require.NoError(t, err)

    // Create sessions
    session1, _ := ctx1.NewSession()
    session2, _ := ctx2.NewSession()

    // Exchange nonces
    nonce1 := session1.PublicNonce()
    nonce2 := session2.PublicNonce()

    _, err = session1.RegisterPubNonce(nonce2)
    require.NoError(t, err)

    _, err = session2.RegisterPubNonce(nonce1)
    require.NoError(t, err)

    // Sign
    msg := [32]byte{1, 2, 3}
    partialSig1, _ := session1.Sign(msg)
    partialSig2, _ := session2.Sign(msg)

    // Aggregate
    haveAll, err := session1.CombineSig(partialSig2)
    require.NoError(t, err)
    require.True(t, haveAll)

    // Verify
    finalSig := session1.FinalSig()
    require.NotNil(t, finalSig)

    aggregatedKey, _ := ctx1.CombinedKey()
    require.True(t, finalSig.Verify(msg[:], aggregatedKey))
}
```

## Performance Characteristics

From POC testing:

- **MuSig2 Signature Size:** 64 bytes
- **Traditional 2-of-2 P2WSH:** ~249 bytes
- **Size Reduction:** 74.3%
- **Fee Reduction:** ~74% (proportional to size)
- **Privacy:** Indistinguishable from single-sig on-chain

## Common Pitfalls

### 1. Nonce Reuse

```go
// ❌ WRONG: Reusing session/nonce
session, _ := ctx.NewSession()
nonce := session.PublicNonce()
sig1, _ := session.Sign(msg1)  // First sign
sig2, _ := session.Sign(msg2)  // DANGER: Nonce reuse!

// ✅ CORRECT: New session per transaction
session1, _ := ctx.NewSession()
sig1, _ := session1.Sign(msg1)

session2, _ := ctx.NewSession()  // New session = new nonce
sig2, _ := session2.Sign(msg2)
```

### 2. Key Sorting Mismatch

```go
// ❌ WRONG: Different sort settings
ctx1, _ := musig2.NewContext(key1, true, ...)   // sorted
ctx2, _ := musig2.NewContext(key2, false, ...)  // unsorted

// ✅ CORRECT: Same sort setting
ctx1, _ := musig2.NewContext(key1, true, ...)
ctx2, _ := musig2.NewContext(key2, true, ...)
```

### 3. Missing Nonce Registration

```go
// ❌ WRONG: Signing without registering all nonces
session, _ := ctx.NewSession()
// Missing: session.RegisterPubNonce(otherNonce)
sig, _ := session.Sign(msg)  // ERROR: not all nonces registered

// ✅ CORRECT: Register all nonces first
session, _ := ctx.NewSession()
session.RegisterPubNonce(nonce1)
session.RegisterPubNonce(nonce2)
sig, _ := session.Sign(msg)
```

## Next Steps for Production

1. **Sub-issue #133:** Create domain layer types
   - Wrap musig2 types in domain layer
   - Add validators for nonces and signatures

2. **Sub-issue #134:** Infrastructure layer implementation
   - Create MuSig2Service interface
   - Implement nonce generation, signing, aggregation

3. **Sub-issue #135:** Nonce management
   - Database schema for nonce storage
   - Uniqueness constraints
   - Cleanup mechanisms

4. **Sub-issues #136-#138:** Use case implementation
   - Keygen wallet use cases
   - Sign wallet use cases
   - Watch wallet use cases

5. **Sub-issue #139:** CLI commands
   - `musig2 nonce` command
   - `musig2 sign` command
   - `musig2 aggregate` command

## References

- [MuSig2 Paper](https://eprint.iacr.org/2020/1261)
- [BIP 327: MuSig2](https://github.com/bitcoin/bips/blob/master/bip-0327.mediawiki)
- [BIP 340: Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [btcd musig2 Package](https://pkg.go.dev/github.com/btcsuite/btcd/btcec/v2/schnorr/musig2)
- [POC Implementation](./musig2_poc.go)
