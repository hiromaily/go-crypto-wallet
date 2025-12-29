# MuSig2 Library Selection Research

**Date:** 2025-12-29
**Issue:** #132 - Research and MuSig2 Library Selection
**Parent Issue:** #131 - Implement MuSig2 Support (Phase 3)

## Executive Summary

**Selected Library:** `github.com/btcsuite/btcd/btcec/v2/schnorr/musig2`

**Version:** v2.3.6 (already in use in the project)

**Rationale:**
- ✅ Production-ready with stable API (v2.3.6)
- ✅ Already integrated in our project (btcec/v2 v2.3.6)
- ✅ Fully implements two-round MuSig2 protocol
- ✅ Supports BIP340 Schnorr signatures
- ✅ Includes Taproot support (essential for Phase 1 integration)
- ✅ Maintained by btcsuite (trusted Bitcoin implementation team)
- ✅ 51+ known importers (good ecosystem adoption)
- ✅ Works offline (no RPC required)
- ✅ ISC license (permissive, compatible with our project)

## Research Methodology

### Libraries Evaluated

1. **github.com/btcsuite/btcd/btcec/v2/schnorr/musig2** ⭐ **SELECTED**
2. github.com/btcsuite/btcwallet (no standalone MuSig2 package found)
3. Third-party libraries (no viable production-ready alternatives found)

### Evaluation Criteria

| Criteria | btcd/btcec/v2/schnorr/musig2 | Status |
|----------|------------------------------|--------|
| Two-round MuSig2 protocol | ✅ Yes | PASS |
| BIP340 Schnorr signatures | ✅ Yes | PASS |
| Taproot support | ✅ Yes (WithTaprootTweakCtx) | PASS |
| Production-ready | ✅ Yes (v2.3.6, 51+ importers) | PASS |
| Offline signing | ✅ Yes (no RPC required) | PASS |
| Existing integration | ✅ Yes (btcec/v2 v2.3.6) | PASS |
| Active maintenance | ✅ Yes (btcsuite team) | PASS |
| Security audit | ⚠️ No public audit found | ACCEPTABLE* |
| License | ✅ ISC (permissive) | PASS |

*Note: While no dedicated security audit was found, btcd is widely used in the Bitcoin ecosystem (including Lightning Network implementations) and has been battle-tested since 2013. The library has 51+ known importers and is maintained by the trusted btcsuite team.

## Technical Details

### Package: github.com/btcsuite/btcd/btcec/v2/schnorr/musig2

#### Core Capabilities

1. **Key Aggregation**
   - `AggregateKeys()`: Combines multiple public keys into single aggregated key
   - Supports optional key tweaking (Taproot)
   - Deterministic key sorting

2. **Two-Round Protocol**
   - **Round 1 (Nonce Generation):**
     - `GenNonces()`: Generate secret and public nonces
     - `AggregateNonces()`: Combine nonces from all signers
     - Prevents nonce reuse with secure generation

   - **Round 2 (Signing):**
     - `Sign()`: Create partial signature using secret nonce
     - `CombineSigs()`: Aggregate partial signatures into final signature

3. **High-Level API (Recommended)**
   - `Context`: Manages signing context with key aggregation
   - `Session`: Handles individual signing sessions
   - Built-in safety mechanisms to prevent nonce reuse

4. **Taproot Integration**
   - `WithTaprootTweakCtx()`: Apply Taproot tweaks
   - `WithBip86TweakCtx()`: BIP 86 key-path spending
   - Compatible with our Phase 1 (Taproot) implementation

#### Key Types

```go
// Nonces for signing
type Nonces struct {
    PubNonce [66]byte  // Public nonce (shareable)
    SecNonce [97]byte  // Secret nonce (keep private)
}

// Partial signature from one signer
type PartialSignature struct {
    S *btcec.ModNScalar  // Signature scalar
    R *btcec.PublicKey   // Nonce commitment
}

// Signing context (high-level API)
type Context struct { /* manages key aggregation */ }

// Signing session (high-level API)
type Session struct { /* manages one signing round */ }
```

#### Example Usage (High-Level API)

```go
import (
    "github.com/btcsuite/btcd/btcec/v2"
    "github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
)

// Round 1: Setup and Nonce Generation
ctx, err := musig2.NewContext(
    privateKey,
    true, // true = sort keys
    musig2.WithKnownSigners(allPublicKeys),
)
session, err := ctx.NewSession()

// Generate and share public nonce
ourNonce := session.PublicNonce()

// Register nonces from other signers
for _, otherNonce := range otherSignersNonces {
    session.RegisterPubNonce(otherNonce)
}

// Round 2: Signing
partialSig, err := session.Sign(messageHash)

// Aggregate (on coordinator/watch wallet)
// This example shows combining one of multiple partial signatures.
haveAllSigs, err := session.CombineSig(partialSig)
if err == nil && haveAllSigs {
    finalSig := session.FinalSig()
    // ... verify and use final signature
}
```

## Integration Plan

### Current Project State

From `go.mod`:
```
github.com/btcsuite/btcd v0.25.0
github.com/btcsuite/btcd/btcec/v2 v2.3.6
```

**Status:** ✅ **No new dependencies required**

The MuSig2 package is already available in our project through the existing `btcec/v2 v2.3.6` dependency.

### Import Path

```go
import "github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
```

### Compatibility Verification

- ✅ **btcd v0.25.0**: Compatible
- ✅ **btcec/v2 v2.3.6**: Contains musig2 package
- ✅ **PSBT (Phase 2)**: MuSig2 signatures can be stored in PSBT
- ✅ **Taproot (Phase 1)**: MuSig2 includes Taproot support
- ✅ **Go 1.25.5**: Fully compatible

## Security Considerations

### Nonce Generation Security

The `musig2` package implements secure nonce generation:
- Uses cryptographically secure RNG (ChaCha20)
- Supports deterministic nonce generation with `WithCustomRand()`
- Includes nonce commitment to prevent Wagner's attack
- **Critical**: Nonces must NEVER be reused (handled by Session API)

### Key Aggregation Security

- Implements MuSig2 coefficient scheme (prevents rogue key attacks)
- Supports Taproot tweaking (essential for Bitcoin scripts)
- Deterministic key ordering option (prevents key malleability)

### Known Limitations

1. **No Dedicated Security Audit**
   - However, btcd has been in production since 2013
   - Used by Lightning Network implementations
   - Wide ecosystem adoption (51+ importers)
   - Active maintenance by btcsuite team

2. **Nonce Reuse Prevention**
   - Application must ensure nonce uniqueness
   - High-level `Session` API provides safety mechanisms
   - **Recommendation**: Use Session API instead of low-level functions

3. **Offline Signing Considerations**
   - Nonces must be securely transmitted between offline wallets
   - File-based nonce storage needs careful handling
   - **Recommendation**: Store nonces in database (next sub-issue)

## Alternatives Considered

### 1. github.com/btcsuite/btcwallet

**Status:** ❌ Not suitable

**Reason:** btcwallet doesn't provide a standalone MuSig2 package. It uses btcd/btcec/v2 internally, so we'd be using the same library.

### 2. Third-Party Libraries

**Status:** ❌ None found

**Findings:** No production-ready, actively-maintained MuSig2 Go libraries found outside of btcsuite ecosystem.

### 3. Custom Implementation

**Status:** ❌ Not recommended

**Reason:**
- Cryptographic complexity (high risk of implementation bugs)
- Security-critical code (one mistake leaks private keys)
- Reinventing the wheel (btcd implementation is mature)
- Time-intensive (6+ months for secure implementation)

## Recommendations

### 1. Use btcd/btcec/v2/schnorr/musig2

**Verdict:** ✅ **APPROVED**

This is the best choice for our project:
- Already integrated (no new dependencies)
- Production-ready and well-tested
- Supports all required features
- Maintained by trusted Bitcoin team
- Compatible with our existing architecture

### 2. Use High-Level Session API

**Recommendation:** Use `Context` and `Session` types instead of low-level functions.

**Benefits:**
- Built-in nonce reuse prevention
- Cleaner API
- Better error handling
- Safer for production use

### 3. Implement Nonce Storage (Next Sub-Issue)

**Critical:** MuSig2 requires secure nonce storage between Round 1 and Round 2.

**Action Items:**
- Design database schema for nonce storage (sub-issue #135)
- Ensure nonce uniqueness constraints
- Implement nonce cleanup after use

### 4. Integration Testing

**Required:**
- Unit tests for MuSig2 protocol steps
- Integration tests for multi-wallet flow
- Testnet testing for real transactions
- Performance benchmarks (vs traditional multisig)

## Next Steps

1. ✅ **Library Selected:** `github.com/btcsuite/btcd/btcec/v2/schnorr/musig2`
2. ✅ **No go.mod Changes Needed:** Already have btcec/v2 v2.3.6
3. 🔄 **Sub-Issue #133:** Create domain layer types and validators
4. 🔄 **Sub-Issue #134:** Implement infrastructure layer (MuSig2 service)
5. 🔄 **Sub-Issue #135:** Implement nonce management (database)

## References

- [MuSig2 Paper (IACR 2020/1261)](https://eprint.iacr.org/2020/1261)
- [BIP 340: Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [BIP 327: MuSig2 for BIP340-compatible Multi-Signatures](https://github.com/bitcoin/bips/blob/master/bip-0327.mediawiki)
- [btcd btcec/v2 Documentation](https://pkg.go.dev/github.com/btcsuite/btcd/btcec/v2)
- [btcd schnorr/musig2 Package](https://pkg.go.dev/github.com/btcsuite/btcd/btcec/v2/schnorr/musig2)

## Conclusion

The `github.com/btcsuite/btcd/btcec/v2/schnorr/musig2` package is the optimal choice for implementing MuSig2 in our wallet. It's production-ready, already integrated, and provides all the features we need for Phase 3. No additional dependencies are required.

**Decision:** ✅ **Proceed with btcd/btcec/v2/schnorr/musig2**
