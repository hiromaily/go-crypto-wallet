## Security Architecture

### Nonce Security

**Critical Requirement**: Each nonce must be used exactly once. Nonce reuse will leak the private key!

#### Multi-Layer Protection

**1. Application Layer**:

```go
// Validate nonce uniqueness before use
func ValidateNonceUniqueness(nonces []NonceCommitment) error {
    seen := make(map[string]bool)
    for _, nonce := range nonces {
        key := string(nonce.Nonce[:])
        if seen[key] {
            return ErrDuplicateNonce
        }
        seen[key] = true
    }
    return nil
}
```

**2. Database Layer** (Future Enhancement):

```sql
CREATE TABLE musig2_nonces (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    signer_id VARCHAR(255) NOT NULL,
    nonce BINARY(66) NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_nonce (nonce),
    UNIQUE KEY unique_tx_signer (transaction_id, signer_id)
);
```

**3. Cryptographic Layer**:

- Nonces generated using secure random number generator
- Each session creates fresh nonces
- Nonces are automatically deleted after use in btcd library

#### Nonce Lifecycle

```
1. Generate → Fresh random nonce
2. Store    → PSBT proprietary field + DB tracking
3. Exchange → Via PSBT file transfers
4. Use      → Sign message once
5. Delete   → Remove after signing complete
```

### Key Aggregation Security

**Public Key Validation**:

```go
// Validate all public keys before aggregation
func ValidatePublicKeysForMuSig2(pubKeys [][]byte) error {
    // Check count (2-15 signers)
    if len(pubKeys) < 2 || len(pubKeys) > 15 {
        return ErrInvalidSignerCount
    }

    // Check each key
    for i, pubKey := range pubKeys {
        if len(pubKey) == 0 {
            return fmt.Errorf("empty public key at index %d", i)
        }
        // Validate key format (33 or 65 bytes)
        if len(pubKey) != 33 && len(pubKey) != 65 {
            return fmt.Errorf("invalid key length at index %d", i)
        }
    }

    // Check for duplicates
    seen := make(map[string]bool)
    for _, pubKey := range pubKeys {
        key := string(pubKey)
        if seen[key] {
            return ErrDuplicatePublicKey
        }
        seen[key] = true
    }

    return nil
}
```

**Taproot Tweak Application**:

```go
// Apply BIP86 Taproot tweak for key-path spending
aggregatedKey, err := musig2Service.AggregatePublicKeys(
    pubKeys,
    true, // applyTaprootTweak
)
```

### Signature Validation

**Partial Signature Verification**:

```go
// Each partial signature is validated before aggregation
func ValidatePartialSignatures(sigs []PartialSignature, expected int) error {
    if len(sigs) != expected {
        return fmt.Errorf("expected %d signatures, got %d", expected, len(sigs))
    }

    for i, sig := range sigs {
        if len(sig.Signature) != 32 {
            return fmt.Errorf("invalid signature length at index %d", i)
        }
        if sig.SignerID == "" {
            return fmt.Errorf("missing signer ID at index %d", i)
        }
    }

    return nil
}
```

**Final Signature Verification**:

```go
// Verify aggregated signature before broadcasting
isValid := musig2Service.VerifyAggregatedSignature(
    aggregatedPubKey,
    messageHash,
    finalSignature,
)
if !isValid {
    return errors.New("aggregated signature verification failed")
}
```

### Attack Vectors and Mitigations

#### 1. Nonce Reuse Attack

**Attack**: Reusing the same nonce for different messages allows an attacker to compute the private key.

**Mitigation**:

- Database unique constraint on nonces
- Application-level validation
- Automatic nonce deletion after use
- Session-based nonce management

#### 2. Rogue Key Attack

**Attack**: Attacker provides a crafted public key that allows them to control the aggregated key.

**Mitigation**:

- MuSig2 protocol includes built-in rogue key protection
- Key sorting ensures deterministic aggregation
- All signers must participate in signing

#### 3. Signature Forgery

**Attack**: Attempt to create valid signature without all required partial signatures.

**Mitigation**:

- Cryptographic proof prevents forgery without all partial signatures
- Verification step before broadcasting
- PSBT validation at each stage

#### 4. PSBT Tampering

**Attack**: Modify PSBT fields to change transaction or signatures.

**Mitigation**:

- PSBT format includes checksums
- Transaction hash validation at each step
- Offline signing prevents network attacks

---
