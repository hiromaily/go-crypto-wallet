# MuSig2 Security Documentation

This document provides comprehensive security guidance for operating MuSig2 (Simple Two-Round Schnorr Multisignatures) in go-crypto-wallet. It covers threat models, security requirements, common pitfalls, and best practices.

## Table of Contents

1. [Security Overview](#security-overview)
2. [Threat Model](#threat-model)
3. [Nonce Security](#nonce-security)
4. [Key Security](#key-security)
5. [Operational Security](#operational-security)
6. [Security Checklist](#security-checklist)
7. [Incident Response](#incident-response)
8. [Appendix: Attack Scenarios](#appendix-attack-scenarios)

---

## Security Overview

### Critical Security Principle

**The #1 Rule of MuSig2 Security**:

```
NEVER REUSE A NONCE

Using the same nonce to sign two different messages
leaks your private key to anyone who observes both signatures.
```

This is not a theoretical vulnerability - it is a **mathematical certainty**. If you sign two different messages with the same nonce, an attacker can calculate your private key with simple algebra.

### Security Architecture

MuSig2 security relies on multiple defense layers:

```
┌─────────────────────────────────────────────────────────┐
│ Layer 1: Cryptographic Protocol (MuSig2)                │
│ - Schnorr signatures (BIP340)                           │
│ - Key aggregation with rogue key protection             │
│ - Two-round signing protocol                            │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│ Layer 2: Application-Level Protection                    │
│ - Nonce uniqueness validation                           │
│ - Signature verification before broadcast               │
│ - Transaction state tracking                            │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│ Layer 3: Database-Level Protection                       │
│ - Unique constraints on nonce columns                   │
│ - Transaction isolation                                 │
│ - Atomic operations                                     │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│ Layer 4: Operational Procedures                         │
│ - Trained operators                                     │
│ - File management protocols                             │
│ - Monitoring and alerting                               │
└─────────────────────────────────────────────────────────┘
```

### Security Properties

MuSig2 provides these security guarantees (when used correctly):

#### Unforgeability

Without knowledge of private keys, it is computationally infeasible to forge a valid MuSig2 signature.

```
Security Level: 128-bit
Equivalent to ECDSA security for 256-bit keys
```

#### Key Aggregation Security

The MuSig2 protocol protects against "rogue key attacks" where a malicious signer crafts their public key to gain full control over the aggregated key.

```
Protection: Deterministic key aggregation with all signer public keys
No single signer can manipulate the aggregated key
```

#### Non-Repudiation

Once a partial signature is created, that signer cannot deny having participated in creating the final signature (given proof of nonce commitment).

```
Property: Each signer commits to their nonce in Round 1
Partial signatures in Round 2 are tied to those commitments
```

---

## Threat Model

### Assets to Protect

1. **Private Keys**
   - **Keygen private key** (offline wallet)
   - **Sign wallet private keys** (offline wallets)
   - **Impact if compromised**: Complete loss of funds

2. **Nonces**
   - **Secret nonces** (generated in Round 1, never shared)
   - **Impact if reused**: Private key leakage → loss of funds

3. **Transaction Data**
   - **Unsigned PSBTs** (before signing)
   - **Partial signatures** (during signing)
   - **Impact if compromised**: Transaction details leaked, no fund loss

### Threat Actors

#### External Attackers

**Capabilities**:

- Observe blockchain transactions
- Observe network traffic (if not encrypted)
- Attempt to exploit software vulnerabilities
- Social engineering attacks

**Goals**:

- Steal funds by exploiting vulnerabilities
- Learn about transaction patterns (privacy attack)
- Disrupt operations (DoS)

**Mitigation**:

- Offline wallets (air-gapped) for key operations
- Encrypted file transport
- Regular security updates
- Employee security training

#### Insider Threats (Malicious)

**Capabilities**:

- Access to operational systems
- Knowledge of procedures
- Ability to manipulate files or processes

**Goals**:

- Steal funds
- Sabotage operations

**Mitigation**:

- Multi-signature requirement (no single person can steal funds)
- Audit logging
- Access controls
- Background checks

#### Insider Threats (Accidental)

**Capabilities**:

- Authorized access to systems
- Perform operations

**Goals**:

- None (accidental errors, not malicious)

**Mitigation**:

- Training and procedures
- Monitoring and alerts
- Error detection (multiple validation layers)

### Attack Surfaces

#### 1. Nonce Generation

**Vulnerability**: Weak random number generation
**Impact**: Predictable nonces → private key recovery
**Mitigation**:

- Use cryptographically secure RNG (`crypto/rand` in Go)
- btcd library handles nonce generation securely
- Never implement custom nonce generation

#### 2. Nonce Storage

**Vulnerability**: Nonce reuse due to storage issues
**Impact**: Private key leakage
**Mitigation**:

- Database unique constraints
- Application-level validation
- Monitoring for reuse attempts

#### 3. File Transport

**Vulnerability**: PSBT file interception/modification
**Impact**: Transaction manipulation, DoS
**Mitigation**:

- Encrypted transport (if over network)
- File integrity checks (checksums)
- Physical security for USB transport

#### 4. Signature Aggregation

**Vulnerability**: Accepting invalid partial signatures
**Impact**: Transaction broadcast failure, fund stuck
**Mitigation**:

- Verify each partial signature
- Verify final aggregated signature
- Test transaction validity before broadcast

#### 5. Software Bugs

**Vulnerability**: Implementation errors in MuSig2 code
**Impact**: Various (depends on bug)
**Mitigation**:

- Use well-tested library (`btcd/btcec/v2/schnorr/musig2`)
- Extensive testing
- Code review
- Security audits

---

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

## Key Security

### Key Generation

#### Secure Key Generation

Keys must be generated on offline systems using cryptographically secure randomness:

```bash
# Keygen wallet (offline system)
keygen create-key --account deposit

# Internally uses:
# - crypto/rand for secure randomness
# - BIP32 for deterministic key derivation
# - BIP39 for mnemonic seed (if applicable)
```

**Security Requirements**:

- Generate keys on air-gapped system (no network connection)
- Use high-quality entropy source
- Never generate keys on internet-connected systems
- Verify key generation process (test recovery)

#### Key Derivation

MuSig2 uses same key derivation as traditional multisig:

```
BIP32 Hierarchical Deterministic Keys:

m / purpose' / coin_type' / account' / change / address_index

For Bitcoin:
m / 44' / 0' / account' / 0 / index     (P2PKH, legacy)
m / 49' / 0' / account' / 0 / index     (P2SH-SegWit)
m / 84' / 0' / account' / 0 / index     (P2WPKH, native SegWit)
m / 86' / 0' / account' / 0 / index     (P2TR, Taproot/MuSig2)
```

For MuSig2:

- **Purpose**: `86'` (BIP86 - Taproot key derivation)
- **Coin Type**: `0'` (Bitcoin)
- **Account**: `0'`, `1'`, etc. (different accounts)
- **Change**: `0` (external/receiving), `1` (internal/change)
- **Index**: `0`, `1`, `2`, ... (address index)

### Key Aggregation

#### Rogue Key Attack Prevention

**Rogue Key Attack**: A malicious signer crafts their public key to gain full control over the aggregated key.

**Example**:

```
Honest signers have public keys: P1, P2
Attacker claims public key: P3' = P3 - P1 - P2

Aggregated key: P_agg = P1 + P2 + P3'
              = P1 + P2 + (P3 - P1 - P2)
              = P3

Result: Attacker has full control (only their private key needed)
```

**MuSig2 Protection**:

MuSig2 uses **deterministic key aggregation coefficients** to prevent this attack:

```
For each signer i, calculate coefficient:
a_i = Hash(L || P_i)

Where L = Hash(P_1 || P_2 || ... || P_n) (all public keys)

Aggregated key:
P_agg = a_1 * P_1 + a_2 * P_2 + ... + a_n * P_n

Now the attacker cannot manipulate P_agg because:
- Coefficients depend on all public keys
- Attacker cannot control coefficients
- Cannot cancel out other signers' keys
```

**Implementation in `btcd`**:

```go
// From github.com/btcsuite/btcd/btcec/v2/schnorr/musig2

// Create context with all signer public keys
ctx, err := musig2.NewContext(
    privateKey,
    true,  // Sort keys (all signers must use same order)
    musig2.WithKnownSigners(allPublicKeys),  // Provide all public keys
    musig2.WithBip86TweakCtx(),              // Apply Taproot tweak
)

// The library handles rogue key protection automatically:
// - Computes coefficients based on all public keys
// - Applies coefficients during aggregation
// - No single signer can manipulate the aggregated key
```

### Key Storage

#### Offline Storage

**Keygen and Sign wallets must be offline (air-gapped)**:

```
Offline System Requirements:
├─ No network interfaces (WiFi, Ethernet disabled)
├─ No Bluetooth
├─ Physical access controls (locked room)
├─ Tamper-evident seals on hardware
└─ Audit logging of all access
```

#### Encryption at Rest

Private keys must be encrypted when stored:

```go
// Example key encryption (simplified)

func encryptPrivateKey(privKey []byte, passphrase string) ([]byte, error) {
    // Derive encryption key from passphrase
    salt := make([]byte, 32)
    _, err := rand.Read(salt)
    if err != nil {
        return nil, err
    }

    // Use Argon2 for key derivation (strong KDF)
    encKey := argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, 32)

    // Encrypt with AES-256-GCM
    block, err := aes.NewCipher(encKey)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    _, err = rand.Read(nonce)
    if err != nil {
        return nil, err
    }

    ciphertext := gcm.Seal(nil, nonce, privKey, nil)

    // Return: salt + nonce + ciphertext
    result := append(salt, nonce...)
    result = append(result, ciphertext...)

    return result, nil
}
```

**Storage Format** (WIF - Wallet Import Format):

```
Wallet Import Format (WIF):
- Private key is base58-encoded with checksum
- Optionally compressed public key indicator
- Used by btcd and Bitcoin Core

Example WIF (testnet):
cT3G9vXZr7T8MhQQmJ9KJR3qL3cPt6nZ3eFt6pKyMwPnEJ7HvDxN
```

#### Backup Strategy

**Rule**: Always maintain encrypted backups of private keys.

```
Backup Locations:
1. Primary offline system (operational)
2. Encrypted USB drive (fireproof safe #1)
3. Encrypted USB drive (fireproof safe #2, different location)
4. Paper backup (BIP39 mnemonic, split with Shamir Secret Sharing)
```

**Recovery Testing**:

- Test key recovery from backup quarterly
- Document recovery procedure
- Train multiple team members on recovery

### Key Access Control

#### Multi-Person Control

**No single person should have access to all keys**:

```
Key Custody Model (2-of-3 multisig):
├─ Keygen Key: Person A has access to offline system
├─ Sign Key 1 (auth1): Person B has access to offline system
└─ Sign Key 2 (auth2): Person C has access to offline system

Transaction requires: Any 2 out of 3 people
- Prevents single point of failure
- Prevents single insider theft
- Maintains availability (1 person can be unavailable)
```

#### Role Separation

| Role | Responsibilities | Key Access |
|------|------------------|------------|
| **Keygen Operator** | Generate keys, First partial signature | Keygen private key only |
| **Sign Operator 1** | Second partial signature | Sign key 1 only |
| **Sign Operator 2** | Third partial signature | Sign key 2 only |
| **Watch Operator** | Create transactions, Aggregate, Broadcast | No private keys |
| **Security Officer** | Audit, Monitor, Review | No operational keys |

**Important**: Watch wallet operator has NO private keys (watch-only).

---

## Operational Security

### Common Pitfalls

#### Pitfall #1: Nonce Reuse

**Description**: Using the same nonce to sign two different transactions.

**Causes**:

- Database failure (nonce not recorded properly)
- Operator error (re-running Round 1 for same transaction)
- Software bug (nonce not generated uniquely)
- File management error (using old PSBT with same nonce)

**Consequences**:

```
CRITICAL: Private key leakage
→ Immediate fund loss risk
→ All funds accessible with that key are at risk
```

**Prevention**:

1. Database unique constraints (technical control)
2. Application validation (technical control)
3. Operator training (human control)
4. File naming conventions (process control)
5. Monitoring (detective control)

**Detection**:

```bash
# Check for nonce reuse
mysql> SELECT nonce, COUNT(*) FROM musig2_nonces GROUP BY nonce HAVING COUNT(*) > 1;
```

**Response**: See [Incident Response - Nonce Reuse](#nonce-reuse-incident)

#### Pitfall #2: Wrong Signer Order

**Description**: Signers use different ordering of public keys when aggregating.

**Causes**:

- Configuration mismatch between wallets
- Manual key entry (typos, wrong order)
- Software version mismatch

**Consequences**:

```
Result: Different aggregated public keys
→ Signature verification fails
→ Funds remain safe, but transaction unusable
→ Need to restart signing process
```

**Prevention**:

1. Use consistent key sorting (sort alphabetically by public key)
2. Configuration management (all wallets use same config)
3. Automated key import (no manual entry)

**Detection**:

```go
// Verify all signers compute same aggregated key

aggregatedKey1 := keygen.GetAggregatedKey()
aggregatedKey2 := sign1.GetAggregatedKey()

if !bytes.Equal(aggregatedKey1, aggregatedKey2) {
    return errors.New("aggregated key mismatch - check key order")
}
```

**Response**:

- Fix key ordering configuration
- Restart signing process with correct order

#### Pitfall #3: Missing Partial Signature

**Description**: Not all signers provide partial signatures.

**Causes**:

- Operator error (forgot to sign)
- File transfer failure (partial signature lost)
- Hardware failure (signing device offline)

**Consequences**:

```
Result: Incomplete signature set
→ Cannot aggregate final signature
→ Transaction cannot be broadcast
→ Funds remain safe, transaction pending
```

**Prevention**:

1. Tracking system (checklist for each transaction)
2. Automated validation (check signature count before aggregation)
3. Redundant file backups

**Detection**:

```go
// Before aggregation, verify all signatures present

func validatePartialSignatures(sigs []PartialSignature, expectedCount int) error {
    if len(sigs) < expectedCount {
        return fmt.Errorf("missing partial signatures: got %d, expected %d",
            len(sigs), expectedCount)
    }
    return nil
}
```

**Response**:

- Identify which signer is missing
- Have that signer complete their partial signature
- Continue aggregation

#### Pitfall #4: File Management Errors

**Description**: Using wrong PSBT file, overwriting files, losing files.

**Causes**:

- Confusing file names
- Manual file management (human error)
- Lack of version control for files

**Consequences**:

```
Result varies:
- Using old PSBT → may include wrong transaction data
- Overwriting file → loss of signatures, need to restart
- Losing file → may need to recreate transaction
```

**Prevention**:

1. Strict file naming convention:

   ```
   payment_{request_id}_{stage}_{step}_{timestamp}.psbt

   Example:
   payment_15_unsigned_0_{timestamp}.psbt   (initial PSBT)
   payment_15_nonce_1_{timestamp}.psbt      (after keygen nonce)
   payment_15_nonce_2_{timestamp}.psbt      (after sign1 nonce)
   payment_15_nonce_3_{timestamp}.psbt      (all nonces collected)
   payment_15_signed_1_{timestamp}.psbt     (after keygen signature)
   payment_15_signed_2_{timestamp}.psbt     (after sign1 signature)
   payment_15_signed_3_{timestamp}.psbt     (all partial signatures)
   payment_15_final_{timestamp}.psbt        (final, ready to broadcast)
   ```

2. File checksums (detect corruption)

   ```bash
   sha256sum payment_15_unsigned_0_1704067200.psbt > payment_15_unsigned_0_1704067200.psbt.sha256
   ```

3. Automated file management (scripts handle naming)

4. File backup (keep copies of all intermediate files)

**Detection**:

- Checksum verification before use
- PSBT analysis (check expected contents)

**Response**:

- If wrong file used: Restart from last correct checkpoint
- If file corrupted: Restore from backup or recreate

#### Pitfall #5: Network Synchronization Issues

**Description**: Offline wallets have stale data.

**Causes**:

- Offline wallets not updated with latest public keys
- Configuration out of sync between wallets
- Key rotation not communicated

**Consequences**:

```
Result: Signing failures
→ Different aggregated keys computed
→ Signature verification fails
→ Need to re-sync and restart
```

**Prevention**:

1. Regular sync schedule (weekly)
2. Manifest files tracking current keys/config
3. Version numbers in configuration files
4. Checksum verification after sync

**Sync Procedure**:

```bash
# 1. Export current public keys from watch wallet
watch export-keys --output keys_v5.json

# 2. Transfer to offline systems (encrypted USB)
# Physically transfer USB to offline systems

# 3. Import on keygen wallet
keygen import-keys --input keys_v5.json --verify

# 4. Import on sign wallets
sign import-keys --input keys_v5.json --verify

# 5. Verify all systems have same keys
# Check checksums of key files on all systems
```

**Detection**:

- Pre-flight checks before signing
- Aggregated key verification

**Response**:

- Re-sync offline wallets
- Verify sync with checksums
- Restart signing with synced config

### Best Practices

#### Operator Training

**Minimum Training Requirements**:

1. **MuSig2 Basics** (2 hours)
   - Read: `docs/chains/btc/musig2/user-guide.md`
   - Understand two-round protocol
   - Understand nonce security

2. **Security Requirements** (1 hour)
   - Read this document
   - Understand threat model
   - Understand nonce reuse consequences

3. **Operational Procedures** (3 hours)
   - Documented SOPs for each role
   - File management protocols
   - Error recovery procedures

4. **Hands-On Practice** (8 hours)
   - Testnet practice (10+ transactions)
   - Error scenario practice
   - Emergency procedures drill

5. **Assessment** (1 hour)
   - Written test on security concepts
   - Practical test on testnet
   - Must score 90%+ to qualify

**Ongoing Training**:

- Quarterly refresher (1 hour)
- Review of any incidents
- Updates for procedure changes

#### File Management

**Standard File Naming Convention**:

```
{tx_type}_{request_id}_{stage}_{step}_{timestamp}.psbt

Fields:
- tx_type: "payment", "deposit", "sweep", etc.
- request_id: Numeric ID from database
- stage: "unsigned", "nonce", "signed", "final"
- step: Numeric counter (0, 1, 2, 3)
- timestamp: Unix timestamp (for uniqueness)

Examples:
payment_42_unsigned_0_1704067200.psbt    # Initial unsigned PSBT
payment_42_nonce_1_1704067201.psbt       # After keygen adds nonce
payment_42_nonce_2_1704067202.psbt       # After sign1 adds nonce
payment_42_nonce_3_1704067203.psbt       # All nonces collected
payment_42_signed_1_1704067204.psbt      # After keygen signature
payment_42_signed_2_1704067205.psbt      # After sign1 signature
payment_42_signed_3_1704067206.psbt      # All partial signatures
payment_42_final_1704067230.psbt         # Final signed PSBT
```

**File Organization**:

```
data/tx/btc/
├── pending/              # Active transactions being signed
│   ├── payment_42_unsigned_0_1704067200.psbt
│   ├── payment_42_nonce_1_1704067201.psbt
│   ├── payment_42_nonce_2_1704067202.psbt
│   ├── payment_42_nonce_3_1704067203.psbt
│   ├── payment_42_signed_1_1704067204.psbt
│   ├── payment_42_signed_2_1704067205.psbt
│   └── payment_42_signed_3_1704067206.psbt
├── completed/            # Successfully broadcast transactions
│   └── payment_42_final_1704067230.psbt
└── failed/               # Failed transactions (for investigation)
    └── payment_15_failed_reason.txt
```

**File Lifecycle**:

```
1. Creation (watch wallet)
   └─> data/tx/btc/pending/payment_42_unsigned_0.psbt

2. Round 1: Nonce Generation
   ├─> Keygen adds nonce → payment_42_unsigned_0_...1.psbt
   ├─> Sign1 adds nonce  → payment_42_unsigned_0_...2.psbt
   └─> Sign2 adds nonce  → payment_42_nonce_0_...3.psbt

3. Round 2: Signing
   ├─> Keygen signs → payment_42_unsigned_0_...1.psbt
   ├─> Sign1 signs  → payment_42_unsigned_1_...2.psbt
   └─> Sign2 signs  → payment_42_unsigned_2_...3.psbt

4. Aggregation (watch wallet)
   └─> Watch aggregates → payment_42_signed_3.psbt

5. Broadcast (watch wallet)
   └─> Transaction broadcast, move to completed/
       └─> data/tx/btc/completed/payment_42_signed_3.psbt

6. Archive (monthly)
   └─> Move old completed transactions to archive/
       └─> data/tx/btc/archive/2024-01/payment_42_signed_3.psbt
```

#### Monitoring

**Key Metrics to Monitor**:

1. **Nonce Uniqueness**

   ```sql
   -- Run every 5 minutes
   SELECT COUNT(*) as duplicate_nonces
   FROM (
       SELECT nonce, COUNT(*) as cnt
       FROM musig2_nonces
       GROUP BY nonce
       HAVING COUNT(*) > 1
   ) duplicates;

   -- Alert if result > 0
   ```

2. **Signing Success Rate**

   ```sql
   -- Run daily
   SELECT
       DATE(created_at) as date,
       COUNT(*) as total_transactions,
       SUM(CASE WHEN status = 'broadcast' THEN 1 ELSE 0 END) as successful,
       SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed,
       (SUM(CASE WHEN status = 'broadcast' THEN 1 ELSE 0 END) * 100.0 / COUNT(*)) as success_rate
   FROM transactions
   WHERE coin = 'btc' AND tx_type = 'musig2'
   GROUP BY DATE(created_at)
   ORDER BY date DESC
   LIMIT 30;

   -- Alert if success_rate < 95%
   ```

3. **Pending Transaction Age**

   ```sql
   -- Run hourly
   SELECT
       id,
       request_id,
       TIMESTAMPDIFF(HOUR, created_at, NOW()) as age_hours
   FROM transactions
   WHERE status = 'pending' AND tx_type = 'musig2'
   HAVING age_hours > 24;

   -- Alert if any transactions older than 24 hours
   ```

4. **File System Health**

   ```bash
   # Check for stale PSBT files
   find data/tx/btc/pending/ -name "*.psbt" -mtime +7

   # Alert if any files older than 7 days in pending/
   ```

**Alerting Channels**:

- **Email**: All alerts (operators + management)
- **SMS**: Critical alerts only (nonce reuse, key compromise)
- **Slack**: All alerts + daily summaries
- **Dashboard**: Real-time metrics display

---

## Security Checklist

### Pre-Deployment Checklist

Before deploying MuSig2 to production:

- [ ] **Documentation Review**
  - [ ] All operators read MuSig2 user guide
  - [ ] All operators read this security documentation
  - [ ] All operators understand nonce reuse consequences

- [ ] **Infrastructure**
  - [ ] Bitcoin Core ≥ v22.0 installed (Taproot support)
  - [ ] Database schema updated (taproot_address column)
  - [ ] Backup systems tested and verified
  - [ ] Offline systems properly air-gapped

- [ ] **Security Controls**
  - [ ] Database unique constraints on nonce column (if using nonce table)
  - [ ] Application-level nonce validation implemented
  - [ ] Signature verification before broadcast
  - [ ] Access controls configured (least privilege)

- [ ] **Testing**
  - [ ] Testnet testing completed (minimum 50 transactions)
  - [ ] Error scenarios practiced (nonce issues, file errors, etc.)
  - [ ] Rollback procedures tested
  - [ ] Recovery procedures tested

- [ ] **Monitoring**
  - [ ] Nonce reuse monitoring configured
  - [ ] Transaction success rate monitoring configured
  - [ ] Alerting channels tested (email, SMS, Slack)
  - [ ] Dashboard created for real-time visibility

- [ ] **Procedures**
  - [ ] Standard Operating Procedures (SOPs) documented
  - [ ] File management protocols defined
  - [ ] Error recovery procedures documented
  - [ ] Incident response plan created

- [ ] **Team Readiness**
  - [ ] All operators completed training
  - [ ] All operators passed assessment (90%+)
  - [ ] On-call rotation schedule defined
  - [ ] Escalation procedures documented

### Operational Security Checklist

Before each MuSig2 transaction signing:

- [ ] **Pre-Flight Checks**
  - [ ] All wallets have same key configuration
  - [ ] PSBT file integrity verified (checksum)
  - [ ] Transaction details reviewed and approved
  - [ ] No pending alerts or warnings

- [ ] **Round 1: Nonce Generation**
  - [ ] Each signer generates fresh nonce
  - [ ] Nonces are unique (not reused)
  - [ ] All nonces collected before proceeding to Round 2
  - [ ] PSBT file contains all expected nonces

- [ ] **Round 2: Signing**
  - [ ] Verify all nonces present before signing
  - [ ] Each signer creates partial signature
  - [ ] Nonce marked as "used" after signing
  - [ ] Partial signatures collected

- [ ] **Aggregation**
  - [ ] All partial signatures present
  - [ ] Aggregated signature verified
  - [ ] Transaction validity tested (testmempoolaccept)
  - [ ] Transaction size and fee verified

- [ ] **Broadcast**
  - [ ] Final approval obtained
  - [ ] Transaction broadcast to network
  - [ ] Transaction ID recorded
  - [ ] Monitoring started for confirmation

### Monthly Security Review

- [ ] **Audit Log Review**
  - [ ] Review all operator actions
  - [ ] Check for suspicious activity
  - [ ] Verify access patterns

- [ ] **Nonce Analysis**
  - [ ] Verify no nonce reuse occurred
  - [ ] Check nonce generation randomness
  - [ ] Review nonce storage integrity

- [ ] **Transaction Analysis**
  - [ ] Review success rates
  - [ ] Analyze any failures
  - [ ] Identify process improvements

- [ ] **Security Posture**
  - [ ] Review access controls
  - [ ] Check backup integrity
  - [ ] Verify monitoring effectiveness
  - [ ] Test alerting (send test alerts)

- [ ] **Team Readiness**
  - [ ] Conduct refresher training if needed
  - [ ] Review and update procedures
  - [ ] Practice emergency scenarios

---

## Incident Response

### Incident Classification

#### Severity 1 (CRITICAL)

**Nonce Reuse or Private Key Compromise Suspected**

- **Impact**: Immediate risk of fund loss
- **Response Time**: Immediate (within 15 minutes)
- **Escalation**: CEO, CTO, Security Officer

**Actions**:

1. STOP all operations
2. Assess scope of compromise
3. Assume key is compromised
4. Plan emergency fund sweep

#### Severity 2 (HIGH)

**Repeated Transaction Failures, System Unavailability**

- **Impact**: Operations disrupted, funds safe
- **Response Time**: Within 1 hour
- **Escalation**: Technical Lead, Operations Manager

**Actions**:

1. Investigate root cause
2. Activate backup procedures
3. Communicate ETA to stakeholders

#### Severity 3 (MEDIUM)

**Single Transaction Failure, Minor Issues**

- **Impact**: Minimal, isolated issue
- **Response Time**: Within 4 hours
- **Escalation**: On-duty operator

**Actions**:

1. Debug and resolve
2. Document in incident log
3. Monitor for recurrence

### Nonce Reuse Incident

**Trigger**: Nonce reuse detected (database alert, manual discovery).

#### Immediate Actions (0-15 minutes)

```
1. STOP ALL OPERATIONS
   - No new nonce generation
   - No signing operations
   - No transaction broadcasts

2. ALERT TEAM
   - Page on-call security officer
   - Alert all MuSig2 operators
   - Notify management

3. PRESERVE EVIDENCE
   - Do NOT modify database
   - Snapshot current state
   - Save all log files

4. ASSESS SCOPE
   - Which nonce was reused?
   - Which signer is affected?
   - Which transactions?
   - Was nonce used to sign different messages?
```

#### Investigation (15 minutes - 2 hours)

```
1. DETERMINE IF KEY IS COMPROMISED
   Query:
   - Were two different messages signed with same nonce?
   - Have those signatures been broadcast?
   - Are they visible on the blockchain?

   If YES to all:
   └─> ASSUME PRIVATE KEY IS COMPROMISED

   If NO (nonce reuse caught before signing):
   └─> KEY LIKELY SAFE, but verify carefully

2. IDENTIFY ROOT CAUSE
   - Software bug?
   - Database failure?
   - Operator error?
   - File reuse?

3. CALCULATE EXPOSURE
   - How many BTC controlled by compromised key?
   - Which addresses?
   - Which UTXOs?
```

#### Mitigation (2-24 hours)

**If Key is Compromised**:

```
EMERGENCY FUND SWEEP

1. PREPARE NEW ADDRESSES
   - Generate NEW keys (not from same seed)
   - Use traditional P2WSH temporarily (faster setup)
   - Verify new keys are properly backed up

2. CREATE SWEEP TRANSACTION
   - Collect all UTXOs controlled by compromised key
   - Send to new safe addresses
   - Use HIGH fee (priority confirmation)

3. SIGN WITH REMAINING KEYS
   - 2-of-3 multisig: Use the 2 uncompromised keys
   - Sign immediately (race against attacker)

4. BROADCAST IMMEDIATELY
   - Submit to mempool with high priority
   - Submit to multiple nodes
   - Monitor for confirmation

5. MONITOR MEMPOOL
   - Watch for competing transactions (attacker trying to steal)
   - If attacker transaction seen, consider RBF (Replace-By-Fee)

6. CONFIRMATION
   - Wait for 1 confirmation (10 min average)
   - Funds are safe once confirmed
```

**If Key is Safe** (nonce reuse caught before signing):

```
1. FIX ROOT CAUSE
   - Patch software bug
   - Repair database
   - Update procedures

2. VERIFY FIX
   - Test in testnet
   - Verify nonce uniqueness enforced

3. GENERATE NEW NONCES
   - Discard old nonces
   - Regenerate fresh nonces
   - Continue with transaction

4. RESUME OPERATIONS
   - Gradual restart
   - Enhanced monitoring
```

#### Post-Incident (24-48 hours)

```
1. ROOT CAUSE ANALYSIS
   - Detailed investigation report
   - Timeline of events
   - Contributing factors
   - Why controls failed

2. CORRECTIVE ACTIONS
   - Fix identified weaknesses
   - Enhance controls
   - Update procedures

3. PREVENTIVE MEASURES
   - Additional safeguards
   - Enhanced monitoring
   - Operator retraining

4. LESSONS LEARNED
   - Team meeting
   - Update documentation
   - Share knowledge

5. FOLLOW-UP
   - Verify corrective actions effective
   - Monitor for recurrence
   - Schedule review in 30 days
```

### Communication Plan

#### Internal Communication

**Severity 1 (Critical)**:

- Immediately alert: CEO, CTO, Security Officer, All Operators
- Method: SMS + Phone Call (ensure acknowledgment)
- Updates: Every 15 minutes until resolved

**Severity 2 (High)**:

- Alert: Technical Lead, Operations Manager, On-call team
- Method: Slack + Email
- Updates: Hourly

**Severity 3 (Medium)**:

- Inform: Operations team
- Method: Slack
- Updates: As resolved

#### External Communication (if applicable)

**Customers/Partners**:

- Inform IF funds are at risk or operations disrupted
- Method: Email + status page update
- Frequency: As situation develops
- Transparency: Provide facts, no speculation

**Regulators** (if applicable):

- Report IF required by regulations
- Consult legal team first
- Follow compliance procedures

---

## Appendix: Attack Scenarios

### Scenario A1: Nonce Reuse Attack

**Attacker Goal**: Extract private key by observing two signatures with same nonce.

**Attack Steps**:

1. Attacker observes transaction 1 on blockchain:
   - Message: m1 (transaction hash)
   - Signature: s1
   - Public nonce: R

2. Attacker observes transaction 2 on blockchain:
   - Message: m2 (different transaction hash)
   - Signature: s2
   - Public nonce: R (same as transaction 1!)

3. Attacker computes private key:

   ```
   s1 = k + Hash(R || P || m1) * x
   s2 = k + Hash(R || P || m2) * x

   Subtract: s1 - s2 = (Hash(R || P || m1) - Hash(R || P || m2)) * x

   Solve for x:
   x = (s1 - s2) / (Hash(R || P || m1) - Hash(R || P || m2))
   ```

4. Attacker now has private key, can steal all funds

**Defense**:

- Nonce uniqueness enforced at multiple layers
- Monitoring detects nonce reuse before broadcast
- Emergency fund sweep if compromise detected

### Scenario A2: Rogue Key Attack (Prevented)

**Attacker Goal**: Gain full control over multisig by manipulating their public key.

**Attack Steps** (Attempted):

1. Honest parties have keys: P1, P2
2. Attacker claims key: P3' = P3 - P1 - P2
3. Naive aggregation: P_agg = P1 + P2 + P3' = P3
4. Attacker controls P_agg with only their private key

**Why This Fails in MuSig2**:

- MuSig2 uses deterministic coefficients
- Coefficients depend on all public keys
- Attacker cannot manipulate aggregated key
- Attack is mathematically prevented

### Scenario A3: File Substitution Attack

**Attacker Goal**: Substitute malicious PSBT file to redirect funds.

**Attack Steps**:

1. Attacker gains access to file transport (USB drive)
2. Attacker replaces PSBT file with malicious version
   - Changes output address to attacker's address
3. Signers unknowingly sign malicious transaction
4. Funds sent to attacker

**Defense**:

1. File integrity checks (checksums)
2. Review transaction details before signing
3. Physical security for file transport
4. Encrypted file containers

**Detection**:

- Checksum verification fails
- Manual review notices wrong address

### Scenario A4: Timing Attack on Nonce Generation

**Attacker Goal**: Predict nonces by observing timing.

**Attack Steps**:

1. Attacker measures time taken for nonce generation
2. Attempts to correlate timing with random values
3. Predicts nonce values
4. Uses predicted nonces to forge signatures

**Why This Fails**:

- `crypto/rand` provides cryptographically secure randomness
- Timing does not leak information about random values
- 2^256 keyspace makes prediction infeasible

### Scenario A5: Social Engineering

**Attacker Goal**: Trick operator into signing malicious transaction.

**Attack Steps**:

1. Attacker impersonates manager
2. Requests urgent transaction signing
3. Operator skips verification steps
4. Signs transaction sending funds to attacker

**Defense**:

1. Verification procedures mandatory (no exceptions)
2. Multi-person approval for transactions
3. Out-of-band confirmation for unusual requests
4. Security awareness training

---

## Conclusion

MuSig2 provides significant benefits (smaller transactions, lower fees, better privacy) but requires strict operational security. Key takeaways:

### Critical Security Rules

1. **NEVER reuse a nonce** - This leaks your private key
2. **Verify everything** - Nonces, signatures, transactions
3. **Offline keys** - Never expose private keys to network
4. **Monitor continuously** - Detect issues before they cause damage
5. **Train thoroughly** - Security requires knowledgeable operators

### Defense in Depth

MuSig2 security uses multiple protective layers:

- Cryptographic protocol (MuSig2 itself)
- Application validation
- Database constraints
- Operational procedures
- Monitoring and alerting

No single layer is perfect, but together they provide robust security.

### Incident Preparedness

- Have incident response procedures documented and tested
- Practice emergency scenarios (nonce reuse, key compromise)
- Know how to quickly sweep funds if needed
- Maintain clear escalation paths

### Continuous Improvement

- Review incidents and near-misses
- Update procedures based on lessons learned
- Stay informed about MuSig2 security research
- Regularly test backups and recovery procedures

---

**Document Version**: 1.0
**Last Updated**: 2025-01-30
**Author**: go-crypto-wallet Security Team
**Related Documents**:

- [MuSig2 User Guide](../crypto/btc/musig2_guide.md)
- [Migration Guide](../migration/traditional_to_musig2.md)
- [Architecture Documentation](../architecture/musig2_architecture.md)

**Related Issues**: #141, #171
