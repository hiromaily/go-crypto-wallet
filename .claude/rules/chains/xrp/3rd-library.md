---
paths:
  - internal/application/usecase/keygen/xrp/*.go
  - internal/application/usecase/sign/xrp/*.go
  - internal/application/usecase/watch/xrp/*.go
  - internal/application/dto/xrp/*.go
  - internal/application/ports/api/xrp/*.go
  - internal/application/ports/repository/watch/xrp_transaction.go
  - internal/infrastructure/api/xrp/connection.go
  - internal/infrastructure/api/xrp/**/*.go
  - internal/interface-adapters/cli/keygen/api/xrp/*.go
  - internal/interface-adapters/cli/watch/api/xrp/*.go
  - internal/interface-adapters/wallet/xrp/*.go
---

# XRP Third-Party Library Rules

## Overview

This project uses **two separate XRP Go libraries** for distinct purposes.
The official XRPL documentation recommends `github.com/XRPLF/xrpl-go`, but this project does **not** use that library.
Do not confuse the three libraries — their module paths are different and they are not interchangeable.

---

## Library Summary

| Library | Module Path | Version | Purpose |
|---------|-------------|---------|---------|
| **Peersyst xrpl-go** | `github.com/Peersyst/xrpl-go` | v0.1.15 | Offline signing, address codec, binary codec |
| **XRPScan xrpl-go** | `github.com/xrpscan/xrpl-go` | v0.2.11 | WebSocket client for XRPL node communication |
| ~~XRPLF xrpl-go~~ | ~~`github.com/XRPLF/xrpl-go`~~ | — | **NOT USED** in this project |

---

## Library 1: `github.com/Peersyst/xrpl-go` — Offline Cryptographic Operations

### What it does

Provides offline (no-network) XRP Ledger cryptographic operations:
- **Binary codec** (`binary-codec`): canonical serialization of XRP transaction objects
- **Address codec** (`address-codec`): Base58Check encoding/decoding of XRP addresses and seeds
- **Wallet** (`xrpl/wallet`): key pair derivation from seed, transaction signing

### Where it is used in this project

| File | Package Used | Purpose |
|------|-------------|---------|
| `pkg/chains/xrp/hash.go` | `address-codec` | `Base58CheckEncode/Decode` for seed encoding |
| `pkg/chains/xrp/sign.go` | `address-codec` | Seed decoding in `decodeSeed()` and `DetectAlgorithmFromSeed()` |
| `internal/infrastructure/api/xrp/signer/peersyst_signer.go` | `binary-codec`, `xrpl/wallet` | Transaction serialization and signing in `PeersystSigner` |

### When to use it

- Any offline cryptographic operation: signing, serialization, seed/address encoding
- Implementing or extending `TransactionSigner` (`internal/application/ports/api/xrp/transaction_signer.go`)
- Extending `pkg/chains/xrp/` utilities (hash, sign, keygen)

### Key types and packages

```go
import addresscodec "github.com/Peersyst/xrpl-go/address-codec"
import binarycodec  "github.com/Peersyst/xrpl-go/binary-codec"
import xrplwallet   "github.com/Peersyst/xrpl-go/xrpl/wallet"

// Address encoding
encoded, err := addresscodec.Base58CheckEncode(bytes, prefix)
decoded, err  := addresscodec.Base58CheckDecode(str)

// Transaction signing
wallet, err   := xrplwallet.FromSeed(seed, "")
txBlob, err   := binarycodec.Encode(txMap)
```

---

## Library 2: `github.com/xrpscan/xrpl-go` — WebSocket Node Client

### What it does

Provides a WebSocket client for communicating with a live XRPL node:
- Connect to public/private XRPL WebSocket endpoints
- Send JSON-RPC requests (`account_info`, `submit`, `tx`, `server_info`, etc.)
- Subscribe to ledger and transaction streams

### Where it is used in this project

| File | Purpose |
|------|---------|
| `pkg/chains/xrp/xrplgo/client.go` | Wraps `xrpl.Client`; implements `AccountInfoProvider` and `TransactionSubmitter` ports |
| `pkg/chains/xrp/xrplgo/account.go` | `GetAccountInfo`, `GetBalance`, `GetTotalBalance` |
| `pkg/chains/xrp/xrplgo/transaction.go` | `SubmitTransaction`, `GetTransaction` |
| `pkg/chains/xrp/xrplgo/ledger.go` | `WaitValidation`, `SubscribeLedger` |

The `xrplgo` package is the **only** location that should import `xrpscan/xrpl-go`.
All other code should depend on the `xrplgo` wrapper or the port interfaces.

### When to use it

- Any online operation requiring a live XRPL node: balance queries, transaction submission, ledger subscription
- Extending or fixing `pkg/chains/xrp/xrplgo/`

See [`pkg/chains/xrp/xrplgo/README.md`](../../../../../pkg/chains/xrp/xrplgo/README.md) for full usage, configuration, implemented interfaces, and known limitations.

### Key types

```go
import xrpl "github.com/xrpscan/xrpl-go"

client, err := xrpl.NewClient(xrpl.ClientConfig{URL: "wss://s.altnet.rippletest.net:51233"})
req := xrpl.BaseRequest{"command": "account_info", "account": addr}
res, err := client.Request(req)
```

---

## Library NOT Used: `github.com/XRPLF/xrpl-go`

The [official XRPL Go tutorial](https://xrpl.org/docs/tutorials/get-started/get-started-go) references `github.com/XRPLF/xrpl-go`.
This library is **not imported** in this project.

### Why it is not used

`XRPLF/xrpl-go` and `Peersyst/xrpl-go` overlap significantly — both provide `addresscodec`, `binarycodec`, and `keypairs` packages. This project adopted `Peersyst/xrpl-go` for offline cryptographic operations because:

1. **Adopted earlier**: `Peersyst/xrpl-go` was integrated first; switching would require replacing working, tested code with no functional benefit.
2. **Same version maturity**: Both libraries are pre-1.0 (`Peersyst` is also at v0.1.x), so `XRPLF` offers no stability advantage.
3. **No WebSocket client needed from it**: `XRPLF/xrpl-go` does not provide a clear WebSocket client for live node communication; `xrpscan/xrpl-go` was chosen specifically for that role.
4. **Duplication risk**: Adding a third xrpl library for the same functionality would create confusion about which codec or signing path to use.

**Do not add it.** The two libraries already in use cover all required functionality:
- Offline operations → `Peersyst/xrpl-go`
- Online/WebSocket operations → `xrpscan/xrpl-go`

---

## Decision Guide

```
Need to sign or serialize a transaction offline?
  → github.com/Peersyst/xrpl-go  (binary-codec, xrpl/wallet)

Need to encode/decode XRP addresses or seeds?
  → github.com/Peersyst/xrpl-go  (address-codec)

Need to query a live XRPL node or submit a transaction?
  → pkg/chains/xrp/xrplgo  (wraps xrpscan/xrpl-go)
  → Never import xrpscan/xrpl-go directly outside xrplgo/

Saw github.com/XRPLF/xrpl-go in documentation or examples?
  → Ignore it. This project does not use it.
```
