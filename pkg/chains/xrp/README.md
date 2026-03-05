# pkg/chains/xrp

XRP (Ripple) cryptocurrency utilities for address encoding, hashing, key management, transaction serialization, and signing.

## Package Structure

```
pkg/chains/xrp/
├── address.go      # Address validation (ValidateAddress)
├── hash.go         # XRP hash types and address generation
├── interface.go    # Key and Hash interfaces
├── keygen.go       # Offline key generation (Ed25519 / secp256k1)
├── serializer.go   # XRP transaction serialization (canonical binary format)
├── sha.go          # XRP-specific SHA-512 half utilities
├── sign.go         # Native transaction signing (Signer)
├── types.go        # Amount/type utilities (ToFloat64, etc.)
│
├── protogen/       # Protocol Buffers generated types (DO NOT EDIT)
│   ├── account.pb.go / account_grpc.pb.go
│   ├── address.pb.go / address_grpc.pb.go
│   └── transaction.pb.go / transaction_grpc.pb.go
│
├── rpc/            # JSON-RPC client and types for direct XRP node communication
│   ├── server.go   # Server info RPC (server_info)
│   ├── account.go  # Account info RPC
│   ├── transaction.go # Transaction RPC
│   └── admin.go    # Admin RPC methods
│
├── xrplclient/     # gRPC client for apps/xrpl-grpc-server
│   └── client.go   # XRPLClient (Deprecated — gRPC server not currently in use)
│
└── xrplgo/         # WebSocket client for XRP Ledger node communication
    ├── client.go   # NewClient, implements ports/api/xrp interfaces
    ├── account.go  # AccountInfoProvider implementation
    ├── transaction.go # TransactionSubmitter implementation
    └── ledger.go   # Ledger utilities
```

## Features

- **Offline Key Generation**: Generate XRP keys without any network dependencies
- **Ed25519 Support**: Recommended algorithm for new accounts (2026+)
- **secp256k1 Support**: Legacy algorithm, compatible with existing accounts
- **Deterministic Generation**: Generate keys from entropy, passphrase, or HD wallet keys
- **Transaction Serialization**: Canonical binary format encoding for signing and submission
- **Native Signing**: Sign transactions offline using `Signer` without a gRPC connection
- **Address Validation**: Validate classic XRP addresses
- **WebSocket Client** (`xrplgo`): Connect to XRPL nodes and submit transactions

## Key Generation

### Using KeyGenerator

```go
import xrp "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"

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

## Transaction Signing

```go
signer := xrp.NewSigner()

txBlob, txHash, err := signer.Sign(tx, seed)
```

`Signer` performs all signing offline — no network connection or gRPC server required.

## WebSocket Client (xrplgo)

```go
import "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/xrplgo"

cfg := xrplgo.DefaultConfig("wss://s.altnet.rippletest.net:51233")
client, err := xrplgo.NewClient(cfg)
if err != nil {
    // handle error
}
defer client.Close()
```

The `xrplgo` client implements the interfaces defined in `internal/application/ports/api/xrp`:
- `AccountInfoProvider` — fetch account info and XRP balances
- `TransactionSubmitter` — submit signed transactions and wait for ledger validation

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

## Deprecated

- `xrplclient.XRPLClient` — gRPC client for `apps/xrpl-grpc-server`. The gRPC server is not currently in use. Do not use this type in new code; it is retained for future re-enablement via dependency injection.

## References

- [XRPL Cryptographic Keys](https://xrpl.org/docs/concepts/accounts/cryptographic-keys)
- [Ed25519 vs secp256k1](https://xrpl.org/docs/concepts/accounts/cryptographic-keys#signing-algorithms)
- [XRPL Transaction Format](https://xrpl.org/docs/references/protocol/transactions)
