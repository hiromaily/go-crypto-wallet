# XRP Library Selection: Dual xrpl-go Architecture

This document explains why go-crypto-wallet uses **two** different `xrpl-go` libraries and the rationale behind this architectural decision.

## Executive Summary

| Library | Version | Primary Use Case | Key Strength |
|---------|---------|------------------|--------------|
| **xrpscan/xrpl-go** | v0.2.11 | Transaction submission, ledger queries | WebSocket communication, reliable submission |
| **Peersyst/xrpl-go** | v0.1.15 | Transaction signing (offline) | Comprehensive wallet/signing API, native Go |

**Why two libraries?** Each library excels in different areas. Rather than compromise on either signing or submission quality, we use both libraries in specialized roles.

---

## Problem Statement

### Original Architecture (Pre-2026)

```
Watch Wallet (Online)
    ↓
gRPC Server (ripple-lib-server)
    ↓
ripple-lib (JavaScript/Node.js)
    ↓
Sign Transaction
```

**Problems:**

1. ❌ **Network dependency for signing**: gRPC server must be running
2. ❌ **Offline signing impossible**: Cannot sign on air-gapped keygen/sign wallets
3. ❌ **JavaScript dependency**: ripple-lib requires Node.js runtime
4. ❌ **Complex deployment**: Additional gRPC server infrastructure
5. ❌ **Single point of failure**: gRPC server outage blocks all signing

### Requirements for 2026 Architecture

1. ✅ **Native Go signing**: No JavaScript/Node.js dependencies
2. ✅ **Offline-capable**: Sign on air-gapped systems without network
3. ✅ **Multi-signature support**: Sequential signing workflow
4. ✅ **Reliable submission**: Robust transaction delivery to XRP Ledger
5. ✅ **Clean Architecture**: Separation of concerns (signing vs submission)

---

## Library Comparison

### Evaluated Libraries

We evaluated three Go libraries for XRP Ledger integration:

1. **xrpscan/xrpl-go** (v0.2.11) - Currently integrated
2. **Peersyst/xrpl-go** (v0.1.15) - New for native signing
3. **XRPLF/xrpl-go** (official) - Future alternative

### Feature Matrix

| Feature | xrpscan/xrpl-go | Peersyst/xrpl-go | XRPLF/xrpl-go |
|---------|-----------------|------------------|---------------|
| **WebSocket Client** | ✅ Excellent | ❌ Not included | ⚠️ Basic |
| **Transaction Submission** | ✅ `SubmitTransaction` | ❌ Not included | ⚠️ Basic |
| **Wallet Derivation** | ❌ Limited | ✅ `wallet.FromSeed()` | ✅ Serialization |
| **Transaction Signing** | ❌ Limited | ✅ `wallet.Sign()` | ✅ Via serialization |
| **Multi-Signature** | ❌ Manual | ✅ `wallet.Multisign()` | ⚠️ Manual combination |
| **Offline Capability** | ⚠️ Partial | ✅ Full (zero network) | ✅ Full |
| **Documentation** | ⚠️ Minimal | ⚠️ Moderate | ✅ Official (XRPLF) |
| **Maintenance** | ⚠️ Community | ✅ Active | ✅ Active (official) |
| **API Completeness** | ✅ WebSocket focus | ✅ Wallet focus | ✅ Comprehensive |

### Decision Matrix

| Criterion | Winner | Rationale |
|-----------|--------|-----------|
| **Transaction Submission** | xrpscan | Already integrated, proven WebSocket implementation |
| **Offline Signing** | Peersyst | Best wallet/signing API, `wallet.Sign()` + `wallet.Multisign()` |
| **Multi-Signature** | Peersyst | Native `wallet.Multisign()` method (vs manual combination) |
| **Deployment Simplicity** | Both | Pure Go, no JavaScript/Node.js dependencies |
| **Long-term** | XRPLF | Official library, but requires more investigation |

---

## Architectural Decision

### Selected Approach: **Specialized Dual-Library Architecture**

```
┌─────────────────────────────────────────────────────────────┐
│                     Watch Wallet (Online)                    │
├─────────────────────────────────────────────────────────────┤
│  CreateTransactionUseCase                                    │
│    ↓                                                         │
│  xrpscan/xrpl-go (AccountInfoProvider)                      │
│    • GetAccountInfo() - query account details               │
│    • GetBalance() - check available balance                 │
└─────────────────────────────────────────────────────────────┘
                          ↓ (JSON file transfer)
┌─────────────────────────────────────────────────────────────┐
│                 Keygen/Sign Wallet (Offline)                 │
├─────────────────────────────────────────────────────────────┤
│  SignTransactionUseCase                                      │
│    ↓                                                         │
│  Peersyst/xrpl-go (TransactionSigner)                       │
│    • wallet.FromSeed(secret) - derive wallet                │
│    • wallet.Sign(tx) - single-sig signing                   │
│    • wallet.Multisign(tx) - multi-sig signing               │
│                                                              │
│  ✅ ZERO network calls (offline guarantee)                  │
└─────────────────────────────────────────────────────────────┘
                          ↓ (JSON file transfer)
┌─────────────────────────────────────────────────────────────┐
│                     Watch Wallet (Online)                    │
├─────────────────────────────────────────────────────────────┤
│  SendTransactionUseCase                                      │
│    ↓                                                         │
│  xrpscan/xrpl-go (TransactionSubmitter)                     │
│    • SubmitTransaction() - broadcast to ledger              │
│    • WaitValidation() - confirm validation                  │
│    • GetTransaction() - query tx status                     │
└─────────────────────────────────────────────────────────────┘
```

### Rationale

**Why Peersyst for Signing?**

1. ✅ **Mature Wallet API**: `wallet.FromSeed()`, `wallet.Sign()`, `wallet.Multisign()`
2. ✅ **Offline-First Design**: Zero network dependencies in signing code
3. ✅ **Multi-Signature Support**: Native `wallet.Multisign()` method (not manual blob combination)
4. ✅ **Developer Experience**: Clean, intuitive API for signing operations
5. ✅ **Proven in Production**: Used by Peersyst's XRPL projects

**Why Keep xrpscan for Submission?**

1. ✅ **Already Integrated**: Working WebSocket communication (v0.2.11)
2. ✅ **Reliable Submission**: Proven implementation of `SubmitTransaction()`
3. ✅ **No Migration Risk**: Avoid breaking existing submission logic
4. ✅ **Specialization**: Purpose-built for WebSocket/RPC communication
5. ✅ **Lower Risk**: Incremental migration (signing first, submission later)

**Why Not XRPLF (Official Library)?**

- ⚠️ **Investigation Required**: Needs deeper evaluation for production use
- ⚠️ **API Differences**: May require significant refactoring
- ⚠️ **Migration Effort**: Would need to replace both signing AND submission
- 📅 **Future Optimization**: Deferred to post-Phase 5 (after transaction flow stable)

---

## Code Examples

### Using Peersyst for Signing (Offline)

```go
import (
    "github.com/Peersyst/xrpl-go/xrpl/wallet"
    "github.com/Peersyst/xrpl-go/xrpl/transaction"
)

// Offline signing on keygen/sign wallet
func signTransaction(secret string, txInput *TxInput) (string, string, error) {
    // 1. Derive wallet from secret (offline)
    w, err := wallet.FromSeed(secret, "")
    if err != nil {
        return "", "", fmt.Errorf("failed to derive wallet: %w", err)
    }

    // 2. Build transaction
    tx := transaction.Payment{
        Account:            txInput.Account,
        Destination:        txInput.Destination,
        Amount:             txInput.Amount,
        Fee:                txInput.Fee,
        Sequence:           txInput.Sequence,
        LastLedgerSequence: txInput.LastLedgerSequence,
    }

    // 3. Sign (offline, no network calls)
    signed, hash, err := w.Sign(&tx)
    if err != nil {
        return "", "", fmt.Errorf("failed to sign: %w", err)
    }

    return hash, signed, nil
}
```

### Using xrpscan for Submission (Online)

```go
import (
    xrplgo "github.com/xrpscan/xrpl-go"
)

// Online submission on watch wallet
func submitTransaction(signedBlob string) (*SentTx, error) {
    // 1. Connect to XRP Ledger (online)
    client := xrplgo.NewWebSocketClient("wss://xrplcluster.com")
    defer client.Close()

    // 2. Submit signed transaction
    result, err := client.SubmitTransaction(context.Background(), signedBlob)
    if err != nil {
        return nil, fmt.Errorf("failed to submit: %w", err)
    }

    // 3. Wait for validation
    validated, err := client.WaitValidation(context.Background(), result.LedgerVersion)
    if err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    return &SentTx{
        Hash:          result.Hash,
        LedgerVersion: validated,
    }, nil
}
```

---

## Trade-offs and Limitations

### Advantages

1. ✅ **Best-of-Breed**: Each library used for its strongest capability
2. ✅ **Offline Security**: Peersyst enables true air-gapped signing
3. ✅ **gRPC Eliminated**: No more xrpl-grpc-server dependency; pure Go stack
4. ✅ **Zero JavaScript**: Pure Go stack (no Node.js runtime)
5. ✅ **Lower Risk**: Keep working submission code untouched

### Disadvantages

1. ⚠️ **Two Dependencies**: Increases dependency count (manageable)
2. ⚠️ **Type Conversion**: May need to convert between library types
3. ⚠️ **Binary Size**: Both libraries included in build (~500KB total)
4. ⚠️ **Maintenance**: Two libraries to track for updates

### Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Type incompatibility | Use shared DTO layer (`internal/application/dto/xrp`) |
| Version conflicts | Both libraries have no overlapping dependencies |
| Binary bloat | Acceptable (~500KB for secure offline signing) |
| Future migration | XRPLF evaluation planned for future optimization |

---

## Migration Path

### Phase 1: Native Go Signing (Completed)

```
✅ Task 1.1: Add Peersyst/xrpl-go dependency
✅ Task 1.2: Create segregated interfaces (TransactionSigner)
✅ Task 3.1: Implement PeersystSigner
✅ Task 5.1: Migrate SignTransactionUseCase to native Go signing
```

### Phase 2: Retire gRPC Server (In Progress)

> **Status**: xrpl-grpc-server is **deprecated**. The wallet now connects directly to rippled via
> WebSocket (port 6006). New code must NOT use the gRPC adapter.
> Remaining gRPC code in `xrpapi.go` / `xrpapi_tx.go` is legacy and will be removed.

```
✅ Decision: xrpl-grpc-server deprecated (direct WebSocket preferred)
✅ E2E scripts updated to use WALLET_RIPPLE_WEBSOCKET_PUBLIC_URL (ws://host:6006)
⏳ Remove remaining gRPC calls in xrpapi_tx.go (PrepareTransaction)
⏳ Remove grpc.ClientConn and protogen imports from xrpapi.go
⏳ Remove gRPC protobuf definitions (internal/infrastructure/api/xrp/xrp/*.pb.go)
⏳ Remove xrpl-grpc-server Docker service from compose files
```

### Phase 3: Evaluate Unified Library (Future)

```
1. Investigate XRPLF/xrpl-go maturity
2. POC: Replace both Peersyst + xrpscan with XRPLF
3. Compare: API quality, performance, maintenance burden
4. Decision: Migrate or maintain dual-library approach
```

---

## Performance Considerations

### Signing Performance (Peersyst)

- **Single transaction**: ~10-50ms (CPU-bound, no network)
- **Batch (100 txs)**: ~1-5 seconds
- **Multi-sig overhead**: ~2x single-sig time per signature

### Submission Performance (xrpscan)

- **WebSocket latency**: ~100-500ms (network-dependent)
- **Validation wait**: ~3-5 seconds (XRP Ledger consensus)
- **Batch submission**: Sequential (avoid sequence conflicts)

### Total Workflow (Create → Sign → Send)

```
Create:   ~100ms  (account queries via xrpscan)
Transfer: manual  (JSON file to offline system)
Sign:     ~50ms   (native Go via Peersyst, offline)
Transfer: manual  (signed file back to watch)
Submit:   ~4s     (broadcast + validation via xrpscan)
```

**Total**: ~4.2 seconds (excluding manual file transfers)

---

## Future Considerations

### When to Revisit This Decision?

1. **XRPLF Maturity**: When XRPLF/xrpl-go reaches feature parity
2. **Unified Need**: If maintaining two libraries becomes burdensome
3. **Performance Issues**: If type conversion overhead becomes significant
4. **Security Audit**: If audit recommends single-source dependency

### Alternative Approaches (Not Chosen)

| Approach | Why Not Chosen |
|----------|----------------|
| **xrpscan only** | Poor signing API, lacks `wallet.Multisign()` |
| **Peersyst only** | No WebSocket client, submission would need rewrite |
| **XRPLF only** | Needs investigation, higher migration risk now |
| **Custom implementation** | Reinventing wheel, security risk, high effort |

---

## References

### Libraries

- [Peersyst/xrpl-go](https://github.com/Peersyst/xrpl-go) - Wallet and signing
- [xrpscan/xrpl-go](https://github.com/xrpscan/xrpl-go) - WebSocket and submission
- [XRPLF/xrpl-go](https://github.com/XRPLF/xrpl-go) - Official (future evaluation)

### XRP Ledger Documentation

- [Multi-Signing](https://xrpl.org/docs/concepts/accounts/multi-signing)
- [Reliable Transaction Submission](https://xrpl.org/docs/concepts/transactions/reliable-transaction-submission)
- [Secure Signing](https://xrpl.org/docs/concepts/transactions/secure-signing)

### Internal Documentation

- [xrpl-go.md](../../../docs/chains/xrp/xrpl-go.md) - xrpscan/xrpl-go capabilities
- [architecture-2026.md](../../../docs/chains/xrp/architecture-2026.md) - Overall XRP architecture
- `.kiro/specs/xrp-transaction-flow-alignment/design.md` - Transaction flow design

---

## Decision Record

| Attribute | Value |
|-----------|-------|
| **Decision Date** | 2026-02-14 |
| **Decision** | Use Peersyst/xrpl-go for signing, xrpscan/xrpl-go for submission |
| **Status** | ✅ Approved (Task 1.1) |
| **Review Date** | Post-Phase 5 (after transaction flow stable) |
| **Owned By** | @hiromaily |

**Approval**: This dual-library approach is the optimal solution for native Go signing while maintaining stable transaction submission. Future unification with XRPLF is possible but deferred to reduce current migration risk.
