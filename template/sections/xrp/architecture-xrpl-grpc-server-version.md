# XRPL (XRP Ledger) Architecture for go-crypto-wallet (OBSOLETE)

> **⚠️ OBSOLETE: This document is no longer valid and contradicts the current project direction.**
>
> **Current Approach**: The project has adopted native Go signing using the `xrpl-go` library to eliminate gRPC dependencies and simplify the architecture. This document represents a previous architectural exploration that is no longer being pursued.
>
> **For current XRP implementation details, see**:
>
> - `.kiro/specs/xrp-transaction-flow-alignment/` - Current specification
> - `docs/chains/xrp/` - Current XRP documentation (excluding this file)
>
> **This file is retained for historical reference only.**

---

## gRPC Bridge via TypeScript (xrpl.js) + Go CLI Integration (Historical Proposal)

This document proposed an XRPL module architecture for **go-crypto-wallet** under the constraint that **Go has limited XRPL library maturity**, so XRPL functionality would be implemented in a **TypeScript gRPC server** built on **xrpl.js**, and the **Go CLI would call that gRPC endpoint**.

- TypeScript SDK: `xrpl.js` (XRPLF) — <https://github.com/XRPLF/xrpl.js>
- gRPC server location: `apps/xrpl-grpc-server/`
- Go CLI uses XRP-specific operations by calling the gRPC API.

---

## 1. Why this architecture (2026 context)

### Goals

- Keep **go-crypto-wallet’s wallet separation model** (Watch / Keygen / Sign).
- Use **xrpl.js** as the canonical XRPL protocol implementation (transaction building, serialization, submission logic).
- Provide a **stable, language-agnostic contract** (gRPC) so the Go CLI stays consistent with other coins.
- Preserve **offline signing** and **auditability** even with a remote protocol service.

### Key principle

> The gRPC server must never require access to private keys.
> Signing remains in the offline components (Keygen/Sign), and the server focuses on:
>
> - building canonical unsigned tx
> - network queries
> - reliable submission and confirmation

---

## 2. High-level runtime topology

```

+----------------------------+        gRPC        +------------------------------+
| Go CLI (watch/keygen/sign) | <----------------> | apps/xrpl-grpc-server (TS)  |
| - orchestrates workflow    |                   | - xrpl.js client            |
| - persists artifacts       |                   | - network access (rippled)  |
+--------------+-------------+                   +--------------+---------------+
|                                                    |
| (signed tx_blob)                                   | (JSON-RPC/WebSocket)
v                                                    v
Local artifact store                               XRPL Nodes / Providers

```

### Responsibilities split

- **Go CLI**: user-facing commands, artifact persistence, offline signing workflows, policy enforcement.
- **TS gRPC server**: XRPL network IO + canonical tx preparation using xrpl.js.

---

## 3. Core transaction lifecycle (single-sig)

### 3.1 Create unsigned tx (Watch/online)

1. Go CLI calls gRPC: `GetAccountState(Account)` to obtain `Sequence`.
2. Go CLI calls gRPC: `SuggestFee()` and applies fee policy.
3. Go CLI calls gRPC: `BuildUnsignedTx(intent)` → returns canonical `tx_json` (unsigned).
4. Go CLI persists an **UnsignedTx artifact**.

### 3.2 Sign (offline Sign wallet)

1. Go CLI loads unsigned artifact.
2. Offline signer derives keypair from `(seed, algorithm)` and signs.
3. Output:
   - `tx_blob`
   - `tx_hash`
   - `tx_json`
4. Go CLI persists **SignedTx artifact**.

### 3.3 Submit & confirm (Watch/online)

1. Go CLI calls gRPC: `SubmitBlob(tx_blob)`
2. Go CLI calls gRPC: `WaitForValidated(tx_hash)`
3. Go CLI persists **SubmissionReceipt**.

> Avoid "sign-and-submit" patterns. Always submit a pre-signed blob.

---

## 4. Multisig transaction lifecycle (serial/offline)

```

Watch (Go) -> BuildUnsignedTx (gRPC)
Signer #1 offline -> AddSignature -> artifact
Signer #2 offline -> AddSignature -> artifact
...
Watch (Go) -> SubmitBlob (gRPC) -> WaitForValidated (gRPC)

```

Suggested artifact states:

- `CREATED_UNSIGNED`
- `SIGNED_PARTIAL`
- `SIGNED_FINAL`
- `SUBMITTED`
- `VALIDATED_SUCCESS`
- `FINAL_FAILURE`

---

## 5. gRPC API design (recommended)

### 5.1 Service boundaries

#### Network / state

- `GetServerInfo()`
- `GetLedgerIndex()`
- `GetAccountState(address)`
- `GetFee()`
- `SuggestFee()`

#### Transaction building

- `BuildUnsignedPaymentTx(request)`
- `BuildUnsignedTrustSetTx(request)`
- `AutofillTx(tx_json)`

#### Submission / confirmation

- `SubmitBlob(tx_blob)`
- `TxStatus(tx_hash)`
- `WaitForValidated(tx_hash, timeout)`

---

## 6. Repository layout (recommended)

### TypeScript server

```
apps/
└── xrpl-grpc-server/
    ├── src/
    │   ├── server.ts
    │   ├── xrpl/
    │   │   └── client.ts
    │   ├── builders/
    │   │   └── payment.ts
    │   └── submit/
    │       └── reliable.ts
    ├── package.json
    └── README.md

proto/
└── xrpapi/
    ├── account.proto
    ├── address.proto
    └── transaction.proto

```

### Go CLI integration

```
├── internal
│   ├── application
│   │   ├── dto
│   │   │   ├── btc
│   │   │   └── xrp
│   │   ├── ports
│   │   │   ├── api
│   │   │   │   ├── btc
│   │   │   │   ├── eth
│   │   │   │   └── xrp
│   │   │   ├── file
│   │   │   ├── persistence
│   │   │   ├── repository
│   │   │   │   ├── cold
│   │   │   │   └── watch
│   │   │   └── wallet
│   │   └── usecase
│   │       ├── bchutil
│   │       ├── keygen
│   │       │   ├── shared
│   │       │   └── xrp
│   │       ├── shared
│   │       ├── sign
│   │       │   ├── shared
│   │       │   └── xrp
│   │       └── watch
│   │           ├── shared
│   │           └── xrp
│   ├── di
│   ├── domain
│   │   └── xrp
│   ├── infrastructure
│   │   ├── api
│   │   │   └── xrp
│   │   │       ├── protogen
│   │   │       ├── testutil
│   │   │       └── xrplgo
│   │   └── wallet
│   │       ├── key
│   │       │   ├── generator
│   │       │   └── strategy
│   │       └── sign
│   │           └── xrp
│   ├── integration_test
│   └── interface-adapters
│       ├── cli
│       │   ├── app
│       │   ├── keygen
│       │   │   ├── api
│       │   │   │   ├── btc
│       │   │   │   └── eth
│       │   │   ├── btc
│       │   │   ├── create
│       │   │   ├── export
│       │   │   ├── imports
│       │   │   └── sign
│       │   ├── sign
│       │   │   ├── create
│       │   │   ├── export
│       │   │   ├── imports
│       │   │   └── sign
│       │   └── watch
│       │       ├── api
│       │       │   ├── btc
│       │       │   ├── eth
│       │       │   └── xrp
│       │       ├── btc
│       │       ├── create
│       │       ├── imports
│       │       ├── monitor
│       │       └── send
│       ├── http
│       └── wallet
│           ├── btc
│           ├── eth
│           └── xrp
```

---

## 7. Security model implications

### Must-have controls

- No private keys in gRPC server memory.
- mTLS between CLI and server.
- Do not log full tx blobs.
- Store artifact hashes.
- Deterministic serialization before signing.

---

## 8. Operational reliability

- Endpoint pool (primary + fallback)
- Health checks
- Backoff + circuit breaking
- Reliable submission loop

---

## 9. Docs to add

1. `docs/chains/xrp/architecture.md` (Done)
2. `docs/chains/xrp/grpc-api.md` (Not yet)
3. `docs/chains/xrp/operations/payment.md` (Not yet)

---

## Appendix A: ADR

**Context**: Go XRPL libraries insufficient.
**Decision**: Use xrpl.js via TypeScript gRPC server.
**Consequences**:

- - protocol correctness
- - stable proto boundary
- - extra service to deploy
- - version sync required
