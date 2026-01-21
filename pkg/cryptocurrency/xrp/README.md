# cryptocurrency/xrp

XRP (Ripple) cryptocurrency utilities for address encoding, hashing, and key management.

## Features

- **Offline Key Generation**: Generate XRP keys without any network dependencies
- **Ed25519 Support**: Recommended algorithm for new accounts (2026+)
- **secp256k1 Support**: Legacy algorithm, compatible with existing accounts
- **Deterministic Generation**: Generate keys from entropy, passphrase, or HD wallet keys

## Key Generation

### Using KeyGenerator

```go
import xrp "github.com/hiromaily/go-crypto-wallet/pkg/cryptocurrency/xrp"

// Create a key generator with Ed25519 (recommended)
gen := xrp.NewKeyGenerator(xrp.AlgorithmEd25519)

// Generate a random key pair
keyPair, err := gen.GenerateRandom()
if err != nil {
    // handle error
}

// Access key pair fields
fmt.Println("Address:", keyPair.ClassicAddress)  // r...
fmt.Println("Seed:", keyPair.Seed)               // s... (NEVER log in production!)
fmt.Println("Algorithm:", keyPair.Algorithm)     // ed25519
```

### From HD Wallet Key

```go
// Derive XRP key from BTC HD wallet private key
hdPrivateKey := // ... 32-byte private key from HD wallet
keyPair, err := gen.DeriveKeyFromHDKey(hdPrivateKey)
```

### From Passphrase (for testing/compatibility)

```go
// Compatible with rippled's wallet_propose passphrase
keyPair, err := gen.GenerateFromPassphrase("my passphrase")
```

## Algorithm Comparison

| Feature | Ed25519 | secp256k1 |
|---------|---------|-----------|
| Signature Size | 64 bytes | 70-72 bytes |
| Performance | Faster | Slower |
| Signature Type | Deterministic | Requires nonce |
| Recommendation | **New accounts** | Legacy compatibility |

## Security Notes

- **NEVER log** seed, seedHex, or private key values
- Store keys encrypted at rest
- Use Ed25519 for new accounts (recommended by XRPL)
- Generate keys offline for maximum security

## Configuration

In `config/wallet/xrp/*.yaml`:

```yaml
ripple:
  key_algorithm: "ed25519"  # or "secp256k1"
  offline_keygen: true      # Use native Go implementation
```

## References

- [XRPL Cryptographic Keys](https://xrpl.org/docs/concepts/accounts/cryptographic-keys)
- [Ed25519 vs secp256k1](https://xrpl.org/docs/concepts/accounts/cryptographic-keys#signing-algorithms)

## Original Source

Original address encoding files are from [github.com/rubblelabs/ripple/tree/master/crypto](https://github.com/rubblelabs/ripple/tree/master/crypto)

