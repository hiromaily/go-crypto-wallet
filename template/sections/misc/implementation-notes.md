## Implementation Notes

### Library Integration: btcd/btcec/v2/schnorr/musig2

**Version**: v2.3.6
**Documentation**: <https://pkg.go.dev/github.com/btcsuite/btcd/btcec/v2/schnorr/musig2>
**License**: ISC (permissive)

#### Key Features Used

**1. MuSig2 Context Creation**:

```go
ctx, err := musig2.NewContext(
    privateKey,           // Signer's private key
    true,                 // Sort keys (MUST be same for all signers)
    musig2.WithKnownSigners(allPublicKeys),  // All participants
    musig2.WithBip86TweakCtx(),              // Taproot tweak for BIP86
)
```

**2. Session Management**:

```go
session, err := ctx.NewSession()
pubNonce := session.PublicNonce()  // [66]byte public nonce
```

**3. Nonce Registration**:

```go
haveAllNonces, err := session.RegisterPubNonce(otherNonce)
```

**4. Signing**:

```go
partialSig, err := session.Sign(messageHash)
```

**5. Aggregation**:

```go
haveAllSigs, err := session.CombineSig(partialSig)
finalSig := session.FinalSig()
```

**6. Verification**:

```go
isValid := finalSig.Verify(messageHash[:], aggregatedPubKey)
```

### PSBT Handling

**PSBT Library**: `github.com/btcsuite/btcd/psbt`

**MuSig2 Extension Fields**:

MuSig2 data is stored in PSBT proprietary fields (BIP174 allows custom fields):

```go
// Proprietary field format
type ProprietaryData struct {
    Identifier []byte  // "musig2"
    Subtype    []byte  // "nonce" or "psig"
    KeyData    []byte  // signer_id
    ValueData  []byte  // nonce or partial signature
}

// Example: Store nonce
func AddMuSig2NonceToPSBT(p *psbt.Packet, signerID string, nonce [66]byte) error {
    prop := psbt.ProprietaryData{
        Identifier: []byte("musig2"),
        Subtype:    []byte("nonce"),
        KeyData:    []byte(signerID),
        ValueData:  nonce[:],
    }
    p.Inputs[0].Unknowns = append(p.Inputs[0].Unknowns, &prop)
    return nil
}

// Example: Extract nonces
func ExtractMuSig2NoncesFromPSBT(p *psbt.Packet) ([][66]byte, error) {
    var nonces [][66]byte
    for _, unknown := range p.Inputs[0].Unknowns {
        if bytes.Equal(unknown.Identifier, []byte("musig2")) &&
           bytes.Equal(unknown.Subtype, []byte("nonce")) {
            var nonce [66]byte
            copy(nonce[:], unknown.ValueData)
            nonces = append(nonces, nonce)
        }
    }
    return nonces, nil
}
```

### Error Handling

**Error Wrapping**:

```go
if err != nil {
    return fmt.Errorf("failed to create MuSig2 context: %w", err)
}
```

**Domain Errors**:

```go
var (
    ErrDuplicateNonce       = errors.New("duplicate nonce detected")
    ErrInvalidSignerCount   = errors.New("invalid number of signers")
    ErrMissingPartialSig    = errors.New("missing partial signature")
    ErrSignatureVerifyFailed = errors.New("signature verification failed")
)
```

### Logging

**Structured Logging**:

```go
logger.Debug("create musig2 taproot address",
    "account_type", input.AccountType.String(),
)

logger.Info("MuSig2 signatures aggregated",
    "signer_count", len(partialSigs),
    "transaction_id", txID,
)

logger.Error("failed to aggregate signatures",
    "error", err,
    "partial_sig_count", len(partialSigs),
)
```

**Security**: Never log private keys, nonces (secret part), or sensitive data.

### Testing Strategy

**Unit Tests**:

- Test each use case in isolation
- Mock external dependencies
- Test error conditions

**Integration Tests**:

- Test complete signing workflows
- Test with real PSBT data
- Test nonce uniqueness enforcement

**Test Files**:

```
internal/application/usecase/keygen/btc/
├── create_musig2_address_test.go
├── musig2_nonce_test.go
└── musig2_sign_test.go

internal/application/usecase/sign/btc/
├── musig2_nonce_test.go
└── musig2_sign_test.go

internal/application/usecase/watch/btc/
└── musig2_aggregate_test.go

internal/infrastructure/api/btc/btc/
└── musig2_test.go
```

---
