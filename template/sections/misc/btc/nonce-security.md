## Nonce Security

### Why Nonce Uniqueness is Critical

#### The Mathematics

In Schnorr signatures, the signature equation is:

```
s = k + Hash(R || P || m) * x

Where:
- s = signature scalar (the signature itself)
- k = secret nonce (random value, must be unique)
- R = public nonce (R = k * G, where G is base point)
- P = public key
- m = message (transaction hash)
- x = private key
```

If you sign two different messages (m1 and m2) with the same nonce k:

```
s1 = k + Hash(R || P || m1) * x
s2 = k + Hash(R || P || m2) * x

Subtract equations:
s1 - s2 = (Hash(R || P || m1) - Hash(R || P || m2)) * x

Solve for private key:
x = (s1 - s2) / (Hash(R || P || m1) - Hash(R || P || m2))
```

**Result**: Anyone observing both signatures can calculate your private key.

#### Real-World Analogy

```
Nonce Reuse = Revealing Your Private Key in Public

It's like using the same encryption key twice for two different messages
with a one-time pad. The security breaks down completely.
```

### Nonce Generation

#### Secure Nonce Generation

The `btcd` library generates nonces securely:

```go
// From github.com/btcsuite/btcd/btcec/v2/schnorr/musig2

// Session manages nonce generation and partial signing
session, err := context.NewSession()
if err != nil {
    return fmt.Errorf("failed to create session: %w", err)
}

// Public nonce (safe to share)
publicNonce := session.PublicNonce()  // [66]byte

// Secret nonce (NEVER shared, managed internally)
// The session object maintains the secret nonce
// It is used when Sign() is called
```

**Security Properties**:

- Uses `crypto/rand` for randomness (cryptographically secure)
- Nonce is unique with overwhelming probability (2^256 keyspace)
- Secret nonce never leaves the session object

#### Nonce Format

A MuSig2 nonce consists of **two nonce points**:

```
Public Nonce Structure:
├─ R1 (33 bytes) - First nonce point (compressed)
└─ R2 (33 bytes) - Second nonce point (compressed)
Total: 66 bytes

Why two nonces?
- MuSig2 protocol improvement over MuSig1
- Provides protection against certain attack scenarios
- See BIP327 for detailed explanation
```

### Nonce Storage and Tracking

#### Database Schema for Nonce Tracking

**Proposed Table** (for future implementation):

```sql
CREATE TABLE musig2_nonces (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,                  -- 'btc', 'bch'
    transaction_id BIGINT NOT NULL,             -- Links to transaction
    signer_id VARCHAR(255) NOT NULL,            -- 'keygen', 'auth1', 'auth2', etc.
    nonce BINARY(66) NOT NULL,                  -- Public nonce (66 bytes)
    used BOOLEAN DEFAULT FALSE,                 -- Has been used for signing?
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    used_at TIMESTAMP NULL,

    -- Constraints
    UNIQUE KEY unique_nonce (nonce),            -- CRITICAL: Prevent nonce reuse
    UNIQUE KEY unique_tx_signer (transaction_id, signer_id),  -- One nonce per signer per tx
    KEY idx_coin_tx (coin, transaction_id),     -- Query by transaction
    KEY idx_used (used)                         -- Query unused nonces
) ENGINE=InnoDB;
```

**Key Constraints**:

1. **`UNIQUE KEY unique_nonce (nonce)`**
   - Most critical constraint
   - Prevents same nonce from being stored twice
   - Database will reject duplicate nonce with error
   - **DO NOT OVERRIDE THIS CONSTRAINT**

2. **`UNIQUE KEY unique_tx_signer (transaction_id, signer_id)`**
   - Ensures each signer generates exactly one nonce per transaction
   - Prevents confusion about which nonce to use

#### Current Implementation (PSBT-Based)

Currently, nonces are stored in PSBT proprietary fields:

```
PSBT Proprietary Field Structure:
├─ Identifier: "musig2"
├─ Subtype: "nonce"
├─ KeyData: signer_id (e.g., "keygen", "auth1")
└─ ValueData: nonce (66 bytes)

Advantages:
- Self-contained (no external state)
- Travels with transaction
- All signers can see all nonces

Disadvantages:
- No database-level enforcement of uniqueness
- Harder to track nonce usage across transactions
```

### Nonce Usage Validation

#### Before Signing (Round 2)

```go
// Pseudo-code for nonce validation before signing

func validateNonceBeforeSigning(txID int64, signerID string, nonce [66]byte) error {
    // 1. Check if nonce already used
    used, err := nonceRepo.IsNonceUsed(txID, signerID)
    if err != nil {
        return fmt.Errorf("failed to check nonce usage: %w", err)
    }
    if used {
        return errors.New("CRITICAL: Nonce already used - DO NOT SIGN")
    }

    // 2. Check if nonce exists in database (if tracking)
    exists, err := nonceRepo.NonceExists(nonce)
    if err != nil {
        return fmt.Errorf("failed to check nonce existence: %w", err)
    }
    if exists {
        return errors.New("CRITICAL: Nonce collision detected - DO NOT SIGN")
    }

    // 3. Validate nonce format
    if len(nonce) != 66 {
        return fmt.Errorf("invalid nonce length: got %d, expected 66", len(nonce))
    }

    // All checks passed
    return nil
}
```

#### After Signing

```go
// Mark nonce as used immediately after creating partial signature

func createPartialSignature(session *musig2.Session, messageHash [32]byte, txID int64, signerID string) error {
    // 1. Create partial signature
    partialSig, err := session.Sign(messageHash)
    if err != nil {
        return fmt.Errorf("failed to sign: %w", err)
    }

    // 2. IMMEDIATELY mark nonce as used
    err = nonceRepo.MarkNonceUsed(txID, signerID)
    if err != nil {
        // CRITICAL: If we can't mark as used, log error but don't block
        // (the signature is already created, can't undo it)
        logger.Error("CRITICAL: Failed to mark nonce as used",
            "tx_id", txID,
            "signer_id", signerID,
            "error", err,
        )
        // Alert operators
        alerting.SendCriticalAlert("Nonce tracking failure", ...)
    }

    // 3. Store partial signature
    // ...

    return nil
}
```

### Nonce Reuse Detection

#### Monitoring

```bash
#!/bin/bash
# Monitor for nonce reuse attempts

# Check for duplicate nonces in database
mysql -u wallet_user -p wallet_db <<EOF
SELECT
    nonce,
    COUNT(*) as usage_count,
    GROUP_CONCAT(signer_id) as signers,
    GROUP_CONCAT(transaction_id) as transactions
FROM musig2_nonces
GROUP BY nonce
HAVING COUNT(*) > 1;
EOF

# If any rows returned, ALERT IMMEDIATELY
# Nonce has been reused - private key may be compromised
```

#### Alerting

```go
// Alert on nonce reuse detection

type NonceReuseAlert struct {
    Nonce         [66]byte
    SignerIDs     []string
    TransactionIDs []int64
    Timestamp     time.Time
}

func checkForNonceReuse() {
    duplicates, err := nonceRepo.FindDuplicateNonces()
    if err != nil {
        logger.Error("Failed to check for nonce reuse", "error", err)
        return
    }

    if len(duplicates) > 0 {
        // CRITICAL ALERT
        for _, dup := range duplicates {
            alert := NonceReuseAlert{
                Nonce:         dup.Nonce,
                SignerIDs:     dup.SignerIDs,
                TransactionIDs: dup.TransactionIDs,
                Timestamp:     time.Now(),
            }

            // Send immediate alert via multiple channels
            alerting.SendEmail("CRITICAL: Nonce Reuse Detected", alert)
            alerting.SendSMS("Nonce reuse - Check email immediately", ...)
            alerting.SendSlack("@channel CRITICAL: Nonce reuse detected", ...)

            // Log for audit trail
            logger.Error("CRITICAL: Nonce reuse detected",
                "nonce", hex.EncodeToString(dup.Nonce[:]),
                "signers", dup.SignerIDs,
                "transactions", dup.TransactionIDs,
            )
        }

        // STOP ALL OPERATIONS
        // Investigate immediately
        // Assume private key may be compromised
    }
}

// Run check every 5 minutes
func startNonceReuseMonitoring() {
    ticker := time.NewTicker(5 * time.Minute)
    go func() {
        for range ticker.C {
            checkForNonceReuse()
        }
    }()
}
```

#### Incident Response for Nonce Reuse

If nonce reuse is detected:

1. **STOP ALL OPERATIONS IMMEDIATELY**
   - Stop generating new nonces
   - Stop signing any transactions
   - Do not broadcast any pending transactions

2. **ASSESS IMPACT**
   - Which private key was affected? (which signer)
   - Was the nonce actually used to sign two different messages?
   - Have any transactions been broadcast?

3. **ASSUME KEY COMPROMISE**
   - Treat the private key as potentially compromised
   - Plan emergency fund sweep (see [Incident Response](#incident-response))

4. **INVESTIGATE ROOT CAUSE**
   - How did nonce reuse occur?
   - Software bug?
   - Database failure?
   - Operator error?

5. **TAKE CORRECTIVE ACTION**
   - Fix root cause
   - Generate new keys if necessary
   - Update procedures to prevent recurrence

---
