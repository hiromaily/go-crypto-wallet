# XRP Transaction Flow (Current Implementation)

This document describes the **as-built** transaction flow for the XRP integration as of PR #635.
It maps each step to the concrete Go files and interfaces that implement it.

For the common multi-wallet architecture shared across all chains, see
[docs/transaction-flow.md](../../transaction-flow.md).
For the future architecture proposal, see [architecture-2026.md](./architecture-2026.md).

---

## Overview

All network communication uses **WebSocket** directly to a `rippled` node.
The former gRPC adapter (`apps/xrpl-grpc-server`) has been fully removed.

The flow spans three wallet processes communicating via **files** (transferred offline):

```
Watch Wallet (online)          Sign/Keygen Wallet (offline)       Watch Wallet (online)
        │                               │                                  │
  Step 1: Create                  Step 2: Sign                      Step 3: Submit
  unsigned tx                     (offline, no network)             signed tx
        │                               │                                  │
  ──► file ──────────────────────────► file ────────────────────────────► XRP network
```

---

## Transport Layer

| Connection | Type | Port | Used by |
|---|---|---|---|
| Public WebSocket | WebSocket | 6006 | All read/submit operations (`PublicXRP`) |
| Admin WebSocket | WebSocket | 6005 | `ledger_accept` standalone mode only (`AdminXRP`) |

`PublicXRP` and `AdminXRP` are separate structs in `internal/infrastructure/api/xrp/public/` and
`admin/` respectively. Both embed the corresponding low-level RPC client from `pkg/chains/xrp/rpc/`.
Construction is via `NewPublicXRPFromCoinType()` and `NewAdminXRPFromCoinType()` in `connection.go`.

---

## Step 1 — Watch Wallet: Create Unsigned Transaction

**Entry point**: `internal/interface-adapters/cli/watch/api/xrp/`
**Use case**: `internal/application/usecase/watch/xrp/create_transaction.go`

### Interfaces used

| Interface | File | Method |
|---|---|---|
| `apixrp.AccountInfoProvider` | `ports/api/xrp/account_info.go` | `GetBalance()` |
| `apixrp.TransactionPreparer` | `ports/api/xrp/interface.go` | `CreateRawTransaction()` |

### Implementation path

```
createTransactionUseCase.Execute()
    │
    ├─ accountInfo.GetBalance(senderAddr)
    │       └─ WSClient: account_info → reads XRP balance
    │
    └─ txPreparer.CreateRawTransaction(sender, receiver, amount, instructions)
            └─ WSClient.PrepareTransaction()
                    │  1. xrprpc.AccountInfo() → WebSocket account_info
                    │     → fetches Sequence + LedgerCurrentIndex
                    │  2. Builds dtoxrp.TxInput:
                    │       TransactionType: "Payment"
                    │       Account:            sender address
                    │       Destination:        receiver address
                    │       Amount:             drops (amount × 1_000_000)
                    │       Fee:                "12" (minimum, overridable via Instructions)
                    │       Sequence:           from account_info
                    │       LastLedgerSequence: LedgerCurrentIndex + MaxLedgerVersionOffset
                    └─ Returns (*dtoxrp.TxInput, rawTxJSON string, error)
```

### Output file format

The unsigned transaction is written to a JSON file via `txFileRepo.WriteXRPJSONFile()`:

```
storage/xrp/tx/deposit_unsigned_<txID>_0.json
```

File content (XRPTransactionFile schema in `ports/file/`):
```json
{
  "transactions": [
    {
      "uuid": "<uuid-v7>",
      "sender_account_type": "client",
      "sender_account": "r...",
      "unsigned_data": { "TransactionType": "Payment", "Account": "r...", ... },
      "required_signatures": 1,
      "signature_count": 0,
      "is_complete": false
    }
  ]
}
```

---

## Step 2 — Sign/Keygen Wallet: Sign Transaction (Offline)

**Entry point**: `internal/interface-adapters/cli/sign/` or `keygen/`
**Use case**: `internal/application/usecase/sign/xrp/sign_transaction.go`
(keygen uses the equivalent at `internal/application/usecase/keygen/xrp/sign_transaction.go`)

### Interface used

| Interface | File | Methods |
|---|---|---|
| `apixrp.TransactionSigner` | `ports/api/xrp/transaction_signer.go` | `SignTransaction()`, `SignTransactionNative()` |

### Implementation path

```
signTransactionUseCase.Sign()
    │
    ├─ txFileRepo.ReadXRPJSONFile(filePath)
    │       └─ reads unsigned JSON file from disk
    │
    ├─ txFile.Validate()
    │       └─ checks invariants (non-empty, valid fields)
    │
    └─ for each tx entry:
            │
            ├─ xrpAccountKeyRepo.GetSecret(senderAccountType, senderAccount)
            │       └─ reads seed/secret from local DB (never logged)
            │
            └─ xrp.SignTransactionNative(ctx, &tx.UnsignedData, secret, isMultiSig, existingBlob)
                    └─ (*XRP).SignTransactionNative()  [xrpapi_tx.go]
                            └─ xrpsigner.NewPeersystSigner().SignTransactionNative()
                                    │  Uses Peersyst/xrpl-go library
                                    │  ZERO network calls — fully offline
                                    │  Single-sig:  wallet.Sign()
                                    │  Multi-sig:   wallet.Multisign() + blob accumulation
                                    └─ Returns (txHash string, txBlob string, error)
```

### Multi-signature accumulation

For M-of-N signing, each signer calls `SignTransactionNative` with the previous signer's
`txBlob` passed as `existingSignedBlob`. The library appends the new `Signer` entry to the
existing `Signers` array rather than replacing it. Signing is fully serial and offline.

### Output file format

The signed result is written back as JSON:

```
storage/xrp/tx/deposit_signed_<txID>_<signedCount>.json
```

The `signed_blob` field of each completed transaction entry is populated with the hex-encoded
signed transaction blob ready for submission.

---

## Step 3 — Watch Wallet: Submit and Confirm

**Entry point**: `internal/interface-adapters/cli/watch/`
**Use case**: `internal/application/usecase/watch/xrp/send_transaction.go`

### Interfaces used

| Interface | File | Methods |
|---|---|---|
| `apixrp.TransactionSubmitter` | `ports/api/xrp/transaction_submitter.go` | `SubmitTransaction()`, `GetTransaction()` |
| `apixrp.LedgerPoller` | `ports/api/xrp/transaction_submitter.go` | `LedgerCurrent()` |
| `apixrp.LedgerAdvancer` (optional) | `ports/api/xrp/transaction_submitter.go` | `LedgerAccept()` |

### Ledger Validation — Plan C Architecture

`WaitValidation` was **removed** from `TransactionSubmitter`. The use case manages the
ledger polling loop itself (Plan C):

- `LedgerPoller` → `PublicXRP.LedgerCurrent()` — polls current ledger index (always available)
- `LedgerAdvancer` → `AdminXRP.LedgerAccept()` — advances ledger in standalone mode (optional, may be nil)

The local `xrpSendTxClient` interface in `send_transaction.go` composes both:
```go
type xrpSendTxClient interface {
    apixrp.TransactionSubmitter
    apixrp.LedgerPoller
}
```
`ledgerAdv apixrp.LedgerAdvancer` is a separate optional field (nil in production without admin).

### Implementation path

```
sendTransactionUseCase.Execute()
    │
    ├─ txFileRepo.ReadFileSlice(filePath)
    │       └─ reads signed tx file lines: "uuid,txHash,txBlob"
    │
    └─ for each entry (concurrent goroutines):
            │
            ├─ xrper.SubmitTransaction(ctx, txBlob)
            │       └─ PublicXRP.SubmitTransaction()
            │               └─ xrprpc.Submit() → WebSocket submit
            │                  Returns SentTx{ResultCode, TxJSON{Hash, LastLedgerSequence}, ...}
            │                  Checks ResultCode contains "tesSUCCESS"
            │
            ├─ use case: waitValidation(ctx, sentTx.TxJSON.LastLedgerSequence)
            │       │  Polls xrper.LedgerCurrent() up to 30 times × 1s
            │       │  If ledgerAdv != nil: calls LedgerAccept() between polls (standalone)
            │       └─ Returns when currentLedger >= LastLedgerSequence
            │
            └─ xrper.GetTransaction(ctx, txHash, earliestLedgerVersion)
                    └─ PublicXRP.GetTransaction()
                            └─ xrprpc.GetTx() → WebSocket tx
                               Confirms Meta.TransactionResult and Validated flag
```

---

## Port Interface Summary

```
apixrp.XRPPublicClient (focused — used everywhere except DI admin setup)
    embeds:
        XRPPublicer         → public queries (AccountChannels, AccountInfo, ServerInfo)
        TransactionSubmitter→ SubmitTransaction, GetTransaction
        TransactionSigner   → SignTransaction, SignTransactionNative
        LedgerPoller        → LedgerCurrent
        AccountInfoProvider → GetAccountInfo, GetBalance
        TransactionPreparer → CreateRawTransaction
        CoinTypeProvider, Closer, ChainConfigProvider

apixrp.XRPAdminClient (focused — used only for keygen + standalone ledger)
    embeds:
        XRPAdminer          → admin keygen ops (ValidationCreate, WalletPropose)
        LedgerAdvancer      → LedgerAccept
        CoinTypeProvider, Closer, ChainConfigProvider

apixrp.XRPer (deprecated — kept for backward compatibility only)
    embeds: XRPAdminClient + XRPPublicClient (effectively the union)
```

Use cases depend only on the **focused interfaces**:

| Use case | Interface dependency |
|---|---|
| `createTransactionUseCase` | `AccountInfoProvider`, `TransactionPreparer` |
| `signTransactionUseCase` (sign/keygen) | `TransactionSigner` |
| `sendTransactionUseCase` | `TransactionSubmitter` + `LedgerPoller` (via local `xrpSendTxClient`), optional `LedgerAdvancer` |
| `monitorTransactionUseCase` | `AccountInfoProvider` |
| `generateKeyUseCase` (keygen) | `XRPAdminClient` (needs `WalletPropose`) |

---

## Removed Functionality (gRPC era)

The following operations were implemented via `apps/xrpl-grpc-server` and have been **removed**
as of PR #632. They have no WebSocket-based replacement in the current codebase:

| Operation | Former file | Status |
|---|---|---|
| `GenerateAddress` / `GenerateXAddress` | `xrpapi_address.go` | Removed |
| `SetRegularKey` transaction | `xrpapi_tx_account.go` | Removed (DI panics) |
| `SignerListSet` transaction | `xrpapi_tx_account.go` | Removed (DI panics) |
| Escrow transactions | `xrpapi_tx_escrow.go` | Removed |
| NFToken transactions | `xrpapi_tx_nftoken.go` | Removed |
| Payment channel transactions | `xrpapi_tx_payment_channel.go` | Removed |
| `CombineTransaction` (multisig aggregation via gRPC) | `xrpapi_tx.go` | Removed |

Multi-signature signing is still supported via `SignTransactionNative` (signature accumulation
in the sign wallet), but the gRPC-based `CombineTransaction` aggregation path is gone.

---

## Key Files Reference

| File | Role |
|---|---|
| `internal/infrastructure/api/xrp/public/client.go` | `PublicXRP` struct (implements `XRPPublicClient`) |
| `internal/infrastructure/api/xrp/public/transaction.go` | `PrepareTransaction`, `SubmitTransaction`, `GetTransaction`, `LedgerCurrent` |
| `internal/infrastructure/api/xrp/public/sign.go` | `SignTransaction`, `SignTransactionNative` |
| `internal/infrastructure/api/xrp/public/account.go` | `GetAccountInfo`, `GetBalance` |
| `internal/infrastructure/api/xrp/admin/client.go` | `AdminXRP` struct (implements `XRPAdminClient`) |
| `internal/infrastructure/api/xrp/admin/keygen.go` | `WalletPropose`, `ValidationCreate` |
| `internal/infrastructure/api/xrp/connection.go` | Factory functions `NewPublicXRPFromCoinType`, `NewAdminXRPFromCoinType` |
| `internal/infrastructure/api/xrp/signer/` | `PeersystSigner` — offline signing via Peersyst/xrpl-go |
| `internal/application/ports/api/xrp/interface.go` | `XRPPublicClient`, `XRPAdminClient`, deprecated `XRPer` |
| `internal/application/ports/api/xrp/transaction_signer.go` | `TransactionSigner` interface |
| `internal/application/ports/api/xrp/transaction_submitter.go` | `TransactionSubmitter`, `LedgerPoller`, `LedgerAdvancer` interfaces |
| `internal/application/usecase/watch/xrp/create_transaction.go` | Step 1: build unsigned tx |
| `internal/application/usecase/sign/xrp/sign_transaction.go` | Step 2: sign (sign wallet) |
| `internal/application/usecase/keygen/xrp/sign_transaction.go` | Step 2: sign (keygen wallet) |
| `internal/application/usecase/watch/xrp/send_transaction.go` | Step 3: submit + confirm (Plan C ledger polling) |
