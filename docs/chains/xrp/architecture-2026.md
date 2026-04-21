<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/xrp/architecture-2026.tpl.md · Run `make docs` to regenerate.
-->

# XRPL (XRP Ledger) Architecture (2026)

This document proposes an XRP Ledger (XRPL) module architecture aligned with **go-crypto-wallet’s Clean Architecture + 3-wallet security model** (Watch / Keygen / Sign).

> **As-built reference**: For the current implementation details (post PR #632), see
> [transaction-flow.md](./transaction-flow.md).

It focuses on:

* **Transaction flow**: build → sign (single or multisig) → submit → confirm (validated ledger)
* **Key generation & custody**: offline-first, algorithm-explicit (secp256k1 / Ed25519), recovery-ready
* **2026-ready operational architecture**: resilient node access, reliable submission, observability, and auditability

> XRPL references:
>
> * Transactions are only final when included in a **validated ledger**.
> * Use **Reliable Transaction Submission** patterns for robust delivery/confirmation.
> * Handle **Fee** and **Sequence** correctly (core requirements for submission).
> * Be explicit about key algorithm defaults (Ed25519 vs secp256k1 differ by tool/library).
> * Avoid “sign-and-submit” for production; prefer signing locally and submitting a `tx_blob`.

---

# 1. Fit with go-crypto-wallet’s Overall Architecture

go-crypto-wallet already defines:

* **Domain / Application / Infrastructure / Interface Adapters**
* **Three wallet roles**:

  * Watch (online)
  * Keygen (offline)
  * Sign (offline)

XRPL should follow the same pattern:

```
internal/domain/xrp/...
internal/application/usecase/.../xrp/...
internal/application/ports/xrp/...
internal/infrastructure/api/xrp/...
internal/interface-adapters/wallet/xrp/...
```

---

# 2. XRPL Module Goals

## 2.1 Security Goals

* Private keys **never** touch online machines (Watch wallet stores only public data).
* Deterministic recovery (seed-based) with **explicit algorithm selection** (secp256k1 or Ed25519).
* Support **multisig** flows (serial signing across offline signers).

## 2.2 Reliability Goals

* Submission uses **reliable transaction submission**:

  * submit
  * verify inclusion in validated ledger
  * treat expiration as final failure
* Transactions must include correct:

  * `Fee`
  * `Sequence`
  * bounded ledger window (`LastLedgerSequence`)

## 2.3 Operability Goals (2026)

* Multiple node endpoints (primary + fallbacks)
* Health checks + metrics
* Full audit trail:

  * unsigned tx
  * partially signed tx
  * fully signed tx
  * submission result
  * validated proof

---

# 3. XRPL Transaction Model Essentials (Domain Encoding)

XRPL transactions require:

* `TransactionType`
* `Account` (sender address)
* `Fee`
* `Sequence`
* Per-type fields (e.g., `Payment` requires `Destination`, `Amount`)

### Finality

A transaction is only final once accepted into a **validated ledger**.

### Recommended Submission Pattern

* Construct + sign locally
* Persist signed artifact
* Submit `tx_blob`
* Confirm validated inclusion

---

# 4. Key Generation & Management

## 4.1 Key Algorithms

XRPL supports:

* `secp256k1`
* `Ed25519`

⚠ Different tools default differently.
If you do not specify the algorithm, the **same seed can generate different addresses**.

### Recommended Rule

Always store the tuple:

* `seed`
* `algorithm` (secp256k1 | ed25519)
* `address`
* `public_key`

This ensures deterministic, tool-agnostic recovery.

---

## 4.2 Offline Keygen Wallet Responsibilities

Keygen wallet creates:

* Seed material (OS CSPRNG)
* Derived keypair
* XRPL classic address
* Public export artifact for Watch wallet

For multisig:

* Export signer public keys
* Export signer list metadata

Keys must be generated only on trusted devices.

---

## 4.3 Suggested Storage Format

### Offline (Sensitive)

`xrp_key_bundle.json`

* seed (encrypted-at-rest)
* algorithm
* public_key
* address
* created_at

### Public (Watch Safe)

`xrp_pub_bundle.json`

* algorithm
* public_key
* address
* signer_role (optional)

---

# 5. Wallet-Level Flows

> The common 3-wallet setup, signing, and monitoring flows are defined in the chain-agnostic reference:
> [docs/transaction-flow.md](../../transaction-flow.md).
> This section describes XRPL-specific concerns on top of that common flow.

## 5.1 Single-Sig Flow (XRPL)

Follows the [common single-sig flow](../../transaction-flow.md#single-sig-flow).

XRPL-specific steps:
- Watch Wallet fetches `Sequence` and decides `Fee` / `LastLedgerSequence` before building the unsigned tx
- Keygen Wallet signs the serialized transaction as a `tx_blob`
- Watch Wallet submits `tx_blob` and polls for validated ledger inclusion

## 5.2 Multisig Flow (XRPL)

Follows the [common multisig flow](../../transaction-flow.md#multisig-flow-m-of-n).

XRPL-specific steps:
- Each signer produces a `Signer` entry appended to the transaction (rather than replacing it)
- Signing is fully offline; signers do not need to coordinate in real time
- Watch Wallet aggregates all `Signer` entries and submits the final `tx_blob`

---

# 6. Transaction Lifecycle (Reliable Submission)

Recommended lifecycle:

1. Create & sign locally
2. Submit to network
3. Confirm validated ledger inclusion
   OR treat as final failure after ledger window expiry

---

## 6.1 Watch Wallet: Build Unsigned Transaction

### Inputs

* intent (Payment / TrustSet / OfferCreate / etc.)
* sender address
* destination
* amount
* memos / destination tag (optional)

### Steps

1. Query account info → get `Sequence`
2. Decide `Fee` strategy
3. Set `LastLedgerSequence`
4. Persist canonical unsigned JSON + hash

---

## 6.2 Offline Wallets: Sign

### Single-Sig

* Derive keypair (seed + algorithm explicit)
* Sign serialized transaction
* Output:

  * `tx_blob`
  * `tx_json`

### Multisig

* Each signer produces `Signer` entry
* Signing is fully offline
* Output updated transaction with appended signer entries

---

## 6.3 Watch Wallet: Submit & Confirm

Rules:

* Prefer submitting `tx_blob`
* Never use sign-and-submit in production
* Persist before/after each network call

Reliable loop:

```
submit →
  check result →
    poll tx status →
      if validated → SUCCESS
      if expired/rejected → FINAL FAILURE
      else → continue polling
```

---

# 7. 2026-Ready Runtime Architecture

## 7.1 Node Access Strategy

* Primary rippled node
* Secondary / tertiary fallback nodes
* Different operators / regions

Add:

* Health checks (ledger freshness, latency)
* Circuit breaker
* Rate limit + backoff

---

## 7.2 Submission Service (Infrastructure)

`XRPLClient` (behind application port):

* `GetAccountInfo(address)`
* `GetSuggestedFee()`
* `SubmitBlob(tx_blob)`
* `GetTxStatus(tx_hash)`
* `GetValidatedLedgerIndex()`

Ensure:

* timeouts
* safe retries
* deterministic audit logging

---

## 7.3 Observability

Emit structured logs/metrics:

* unsigned_tx_id
* tx_hash
* submission attempts
* ledger index
* final outcome:

  * validated_success
  * expired
  * rejected

---

# 8. Clean Architecture Mapping

## 8.1 Domain Layer

`internal/domain/xrp/`

### Value Objects

* `XRPAddress`
* `XRPPublicKey`
* `XRPLSeed`
* `XRPLAlgorithm`
* `XRPAmount` (drops)
* `LedgerIndex`
* `Sequence`
* `FeeDrops`

### Entities

* `UnsignedTransaction`
* `SignedTransaction`
* `SubmissionReceipt`

### Domain Services

* `TxSerializer`
* `TxHashCalculator`

---

## 8.2 Application Layer

`internal/application/usecase/.../xrp/`

Suggested Use Cases:

* `CreateUnsignedPaymentTx`
* `CreateUnsignedTrustSetTx`
* `SignTxSingle`
* `SignTxMulti`
* `SubmitSignedTxReliable`
* `QueryTxStatus`

---

## 8.3 Ports

`internal/application/ports/xrp/`

* `XRPLNetworkPort`
* `XRPLStoragePort`
* `ClockPort`
* `AuditPort` (optional)

---

## 8.4 Infrastructure

`internal/infrastructure/api/xrp/`

* WebSocket client (`xrpscan/xrpl-go`) — online queries and transaction submission
* Offline signing (`Peersyst/xrpl-go`) — zero network, used by keygen/sign wallets

> **Deprecated (do not use)**: The xrpl-grpc-server gRPC adapter has been retired.
> The wallet connects directly to rippled via WebSocket on port 6006.
> See [library-selection.md](./library-selection.md) for the dual-library architecture.

---

## 8.5 Interface Adapters

`internal/interface-adapters/wallet/xrp/`

CLI examples:

```
watch  -coin xrp tx create-payment ...
keygen -coin xrp keygen ...
signX  -coin xrp tx sign ...
watch  -coin xrp tx submit ...
watch  -coin xrp tx status ...
```

---

# 9. Security Checklist

* Algorithm explicitness (store algorithm with seed)
* Offline signing only
* No sign-and-submit in production
* Reliable submission validation
* Deterministic audit artifacts
* Key generation only on trusted machines

---

# 10. Suggested Repository Placement

```
docs/chains/xrp/architecture.md
```

Optional:

```
docs/chains/xrp/operations/
docs/chains/xrp/tx-lifecycle.md
```

---

# Appendix A: Minimal Reliable Submission State Machine

### States

* `CREATED_UNSIGNED`
* `SIGNED_PARTIAL`
* `SIGNED_FINAL`
* `SUBMITTED`
* `VALIDATED_SUCCESS`
* `FINAL_FAILURE`

### Core Principle

Never report “success” until the transaction is included in a **validated ledger**.
