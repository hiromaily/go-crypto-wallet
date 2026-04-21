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
