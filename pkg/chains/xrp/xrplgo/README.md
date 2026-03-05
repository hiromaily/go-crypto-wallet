# pkg/chains/xrp/xrplgo

A WebSocket client for XRP Ledger (XRPL) node communication.

This package is a project-specific extension of [`github.com/xrpscan/xrpl-go`](https://github.com/xrpscan/xrpl-go).
It wraps the low-level `xrpl.Client` and adds typed request/response handling, drop-to-XRP conversion,
parallel balance queries, and ledger stream subscription — while implementing the port interfaces
defined in `internal/application/ports/api/xrp`.

## Relationship to xrpscan/xrpl-go

| Layer | Role |
|-------|------|
| `github.com/xrpscan/xrpl-go` | Low-level WebSocket transport: `Client`, `BaseRequest`, `BaseResponse`, stream channels |
| `pkg/chains/xrp/xrplgo` (this package) | Typed wrapper: implements port interfaces, decodes responses, converts amounts |

All imports of `xrpscan/xrpl-go` are confined to this package. Other packages must use this wrapper, never the underlying library directly.

## Files

| File | Responsibility |
|------|---------------|
| `doc.go` | Package documentation |
| `client.go` | `Client`, `ClientConfig`, `NewClient` — wraps `xrpl.Client`, manages connection lifetime |
| `account.go` | `GetAccountInfo`, `GetBalance`, `GetTotalBalance` — `account_info` RPC, parallel balance fetch |
| `transaction.go` | `SubmitTransaction`, `GetTransaction` — `submit` and `tx` RPC, amount parsing |
| `ledger.go` | `WaitValidation`, `SubscribeLedger` — ledger stream subscription, validation polling |
| `types.go` | `AccountInfo`, `TxInput`, `SentTx`, `TxInfo`, `TxOutcome`, and supporting types |

## Implemented Port Interfaces

| Interface (internal/application/ports/api/xrp) | Methods |
|------------------------------------------------|---------|
| `AccountInfoProvider` | `GetAccountInfo` |
| `BalanceChecker` | `GetBalance`, `GetTotalBalance` |
| `TransactionSubmitter` | `SubmitTransaction`, `GetTransaction`, `WaitValidation` |

## Usage

```go
import "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/xrplgo"

// Create client with defaults (recommended)
cfg := xrplgo.DefaultConfig("wss://s.altnet.rippletest.net:51233")
client, err := xrplgo.NewClient(cfg)
if err != nil {
    // handle error
}
defer client.Close()

// Get account info
info, err := client.GetAccountInfo(ctx, "rXXX...")

// Get XRP balance (returns XRP, not drops)
balance, err := client.GetBalance(ctx, "rXXX...")

// Get total balance across multiple addresses (parallel)
total := client.GetTotalBalance(ctx, []string{"rAAA...", "rBBB..."})

// Submit a signed transaction
sentTx, lastLedgerSeq, err := client.SubmitTransaction(ctx, signedTxBlob)

// Wait for ledger validation
validatedLedger, err := client.WaitValidation(ctx, lastLedgerSeq)

// Subscribe to ledger stream
ledgerCh, err := client.SubscribeLedger(ctx)
for ledgerIndex := range ledgerCh {
    // called on each validated ledger close
}
```

## Configuration

```go
cfg := xrplgo.ClientConfig{
    WebSocketURL:      "wss://s.altnet.rippletest.net:51233",
    ReadTimeout:       60 * time.Second,   // default
    WriteTimeout:      60 * time.Second,   // default
    HeartbeatInterval: 5 * time.Second,    // default
    ValidationTimeout: 5 * time.Minute,   // default; used by WaitValidation
}
```

Use `DefaultConfig(url)` to get these defaults with only the URL specified.

## Amount Representation

All amounts returned by this package are in **XRP** (not drops).
The XRPL protocol uses drops internally (1 XRP = 1,000,000 drops); this conversion is handled transparently.

```
raw XRPL response: "1000000"  →  this package returns: "1"  (XRP)
```

## Known Limitations

- `xrpscan/xrpl-go` does not expose context cancellation on the underlying WebSocket calls;
  the `ctx` parameter on request methods is accepted for interface compatibility but does not cancel in-flight network I/O.
- The ledger stream subscription (`WaitValidation`, `SubscribeLedger`) respects `ctx.Done()` between messages but not mid-read.

## Public XRPL WebSocket Endpoints

| Network | URL |
|---------|-----|
| Mainnet | `wss://xrplcluster.com` |
| Testnet | `wss://s.altnet.rippletest.net:51233` |
| Devnet  | `wss://s.devnet.rippletest.net:51233` |
