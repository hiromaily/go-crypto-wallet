# Working History: ETH MPC-TSS Pattern 4 E2E Fix

## Status: Blocked on `serve_mpc.go` — TSS signing not implemented

---

## Goal

`make eth-e2e-p4-reset` must pass end-to-end.

The script tests a 2-of-3 MPC threshold signing flow:
1. Generate Paillier pre-params per node
2. Run DKG ceremony (3 nodes, P2P mode)
3. Fund the joint ETH address
4. Create an unsigned MPC transaction
5. Start 2 `serve mpc` node daemons
6. `watch send mpc` → coordinator signs and broadcasts

---

## What Was Already Working (Before This Session)

- DKG ceremony (P2P mode) ✅ — produces a joint ETH address
- Fund phase ✅
- Create MPC tx phase ✅
- Node daemon startup ✅ (binds gRPC ports 9001/9002)
- Coordinator connects to nodes ✅

---

## Root Cause of Current Hang

**`make eth-e2e-p4-reset` hangs at the "Send Phase"** because `serve_mpc.go` does not implement TSS signing. After verifying the session request, it blocks forever with `<-ctx.Done()`.

The coordinator (`MPCCoordinator.relayLoop`) waits for `MPCWireMessage{IsSignature: true}` from nodes, but nodes never produce it.

---

## Protocol Flow Analysis

```
Watch (coordinator)                      Node (serve-mpc)
─────────────────────────────────────────────────────────
1. InitSigning(sessionId, hash,   ──►   GRPCInboundTransport.InitSigning():
              partyIds, threshold)         - stores sessionID only
                                           - discards hash, partyIds, threshold ← BUG
2. RelaySession (bidi stream)     ◄──►   bidirectional gRPC stream
3. relayLoop: wait for messages   ◄────  ??? nodes never send anything ← BUG
```

**Problem 1**: `GRPCInboundTransport.InitSigning()` only stores `sessionID`. It discards `hash`, `partyIds`, `threshold` — the data the node needs to run TSS signing.

**Problem 2**: `serve_mpc.go` reads the first message from `recvCh` expecting JSON `MPCSessionRequest` (with `raw_tx_hex`), but nothing sends this. The coordinator's `InitSigning` RPC doesn't put anything into `recvCh`.

**Problem 3**: After the first message check, `serve_mpc.go` just blocks. No TSS signing party is created, no messages are processed, no signature is produced.

---

## Files Involved

### Files Needing Changes
| File | Problem | Fix Needed |
|------|---------|------------|
| `internal/infrastructure/api/eth/mpc/grpc_inbound.go` | `InitSigning` discards hash/partyIds/threshold | Expose session data (see options below) |
| `internal/application/usecase/keygen/eth/serve_mpc.go` | No TSS signing implemented | Full rewrite of Serve() |
| (possibly) `internal/application/ports/api/eth/interfaces_mpc.go` | May need new port for signing node | Add `MPCSigningNodePort` |

### Files Already Fixed (This Branch)
| File | Change |
|------|--------|
| `internal/application/ports/api/eth/interfaces_mpc.go` | Added `ListenAddr` to `DKGParams` |
| `internal/application/usecase/keygen/interfaces.go` | Added `ListenAddr` to `RunDKGInput` |
| `internal/application/usecase/keygen/eth/run_dkg.go` | Pass `ListenAddr` to `DKGParams` |
| `internal/infrastructure/api/eth/mpc/node_server.go` | P2P DKG mode (runDKGP2P, sendToPeers, waitForTCPAddr, etc.) |
| `internal/interface-adapters/cli/keygen/dkg/dkg.go` | Added `--listen-addr` flag |
| `scripts/operation/eth/e2e/e2e-p4.sh` | Per-node DBs, DKG P2P flags, port cleanup, 2-of-3 signing |

---

## What Needs to Be Implemented

### Architecture Decision: Where Does TSS Signing Logic Live?

`serve_mpc.go` is in `application/usecase/` — it **cannot** import `tss-lib/v2` directly (clean architecture violation).

**Recommended approach** (matches existing DKG pattern):
1. Define `MPCSigningNodePort` interface in `application/ports/api/eth/interfaces_mpc.go`
2. Implement it in `internal/infrastructure/api/eth/mpc/signing_node.go` (or extend `node_server.go`)
3. Inject into `serveMPCUseCase` via `NewServeMPCUseCase()`
4. `serve_mpc.go` calls `port.RunSigning(ctx, params)` — no TSS imports needed

This mirrors exactly how DKG works:
- `MPCKeyGeneratorPort.RunDKG()` ← defined in ports, implemented in `node_server.go`
- `MPCSigningNodePort.RunSigning()` ← new, same pattern

### Step 1: Define `MPCSigningNodePort`

**File:** `internal/application/ports/api/eth/interfaces_mpc.go`

```go
// MPCSigningNodeParams carries configuration for a signing session on this node.
type MPCSigningNodeParams struct {
    SessionID   string
    Hash        []byte   // 32-byte Keccak256 hash of the unsigned transaction
    PartyID     string   // this node's party ID
    AllPartyIDs []string // all T signing parties (from coordinator's InitSigning)
    Threshold   int
    ShardJSON   []byte   // decrypted shard (LocalPartySaveData JSON)
}

// MPCSigningNodePort is implemented by MPC node servers to run the TSS signing state machine.
type MPCSigningNodePort interface {
    RunSigning(ctx context.Context, params MPCSigningNodeParams) error
}
```

### Step 2: Expose Session Data from `InitSigning`

**File:** `internal/infrastructure/api/eth/mpc/grpc_inbound.go`

The simplest approach that doesn't break `MPCInboundTransport` interface:
- Add a `sessionInfoCh chan signingSessionInfo` field to `GRPCInboundTransport`
- In `InitSigning()`, send `{hash, partyIDs, threshold}` to `sessionInfoCh`
- Export via a new method `AwaitSessionInfo(ctx) (signingSessionInfo, error)`

Since `serve_mpc.go` gets the transport as `apieth.MPCInboundTransport` (interface), we need to either:
- A. Add `AwaitSessionInfo` to the `MPCInboundTransport` interface (requires mock update)
- B. Type-assert to `*GRPCInboundTransport` in the use case (violates clean arch)
- C. Create a separate `MPCSigningSessionProvider` interface with just `AwaitSessionInfo`

**Recommended: Option A** — extend `MPCInboundTransport`:
```go
// Add to MPCInboundTransport interface:
// AwaitSessionInfo blocks until InitSigning is received, then returns the session parameters.
AwaitSessionInfo(ctx context.Context) (MPCSigningSessionInfo, error)

// MPCSigningSessionInfo carries the session parameters delivered by InitSigning.
type MPCSigningSessionInfo struct {
    SessionID string
    Hash      []byte
    PartyIDs  []string
    Threshold int
}
```

### Step 3: Implement `MPCSigningNode` in Infrastructure

**File:** `internal/infrastructure/api/eth/mpc/signing_node.go` (new file)

```go
type MPCSigningNode struct {
    transport    apieth.MPCInboundTransport  // already has the running gRPC server
    logger       *slog.Logger
}

func (s *MPCSigningNode) RunSigning(ctx context.Context, params apieth.MPCSigningNodeParams) error {
    // 1. Unmarshal shard JSON into keygen.LocalPartySaveData
    // 2. Build sorted party IDs from params.AllPartyIDs
    // 3. Find this node's *tss.PartyID
    // 4. Create tss.Parameters (threshold-1 for GG18)
    // 5. Create outCh, endCh, signing.NewLocalParty(msgBig, tssParams, saveData, outCh, endCh)
    // 6. party.Start()
    // 7. Get recvCh from transport.Receive()
    // 8. Run message loop:
    //    - outCh: marshal MPCWireMessage, EnqueueOutbound
    //    - recvCh: unmarshal MPCWireMessage, party.Update()
    //    - endCh: assemble 65-byte sig, EnqueueOutbound with IsSignature=true
}
```

Key tss-lib/v2 imports (all in infrastructure layer):
```go
"math/big"
"github.com/bnb-chain/tss-lib/v2/common"
"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
"github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
"github.com/bnb-chain/tss-lib/v2/tss"
```

Signature from `*common.SignatureData`:
```go
sig := make([]byte, 65)
copy(sig[:32], sigData.R)
copy(sig[32:64], sigData.S)
sig[64] = sigData.SignatureRecovery[0]  // 0 or 1
```

### Step 4: Rewrite `serve_mpc.go`

```go
func (u *serveMPCUseCase) Serve(ctx context.Context, input keygenusecase.ServeMPCInput) error {
    // 1. Load key shard
    shardJSON, err := u.shardStorage.LoadShard(ctx, input.ShardPath, input.Passphrase)

    // 2. Start gRPC listener
    u.transport.Listen(ctx, input.ListenAddr)
    defer u.transport.Close()

    // 3. Wait for session info from coordinator's InitSigning
    sessionInfo, err := u.transport.AwaitSessionInfo(ctx)
    // sessionInfo: {SessionID, Hash, PartyIDs, Threshold}

    // 4. Run TSS signing
    params := apieth.MPCSigningNodeParams{
        SessionID:   sessionInfo.SessionID,
        Hash:        sessionInfo.Hash,
        PartyID:     input.PartyID,
        AllPartyIDs: sessionInfo.PartyIDs,  // coordinator sends only T parties
        Threshold:   sessionInfo.Threshold,
        ShardJSON:   shardJSON,
    }
    return u.signingNode.RunSigning(ctx, params)
}
```

`serveMPCUseCase` struct gains a new field `signingNode apieth.MPCSigningNodePort`.

### Step 5: Update DI Wiring

**File:** `internal/di/container.go`

- Add `MPCSigningNode` factory
- Pass it to `NewServeMPCUseCase`

### Step 6: Update Mock

Run `make mockery` to regenerate mock for updated `MPCInboundTransport` interface.

---

## Signature Verification Note

`applySignature()` in `send_mpc_transaction.go` normalizes v:
```go
if sigCopy[64] >= 27 {
    sigCopy[64] -= 27
}
```

So the node should send `v ∈ {0, 1}` or `v ∈ {27, 28}` — both work.

---

## Testing Strategy

After implementation:
1. `make go-lint && make tidy && make check-build` — no compile errors
2. `make build-keygen && make build-watch` — build binaries
3. `make eth-e2e-p4-reset` — full E2E test
4. Check logs in `data/mpc/e2e-p4/` for TSS round messages being exchanged

---

## Known Issues / Risks

1. **Concurrency**: Two nodes run the same TSS signing session simultaneously. They exchange messages via the coordinator relay. The coordinator routes `from=node1` messages to `node2` and vice versa. Ensure `handleInbound` correctly identifies sender from `wm.From`.

2. **Session ID**: In the signing phase, the coordinator's `relayLoop` discards messages without `IsSignature`. But the node sends TSS round messages wrapped as `MPCWireMessage{From, To, IsBroadcast, Data}`. The coordinator routes them. The `GRPCInboundTransport` session ID filter uses the session ID set by `InitSigning`. This should work since `InitSigning` sets the session ID before `RelaySession` opens.

3. **`recvCh` ordering**: After `AwaitSessionInfo()` consumes the session init from `InitSigning`, subsequent messages on `recvCh` are all TSS round messages. The signing loop should work correctly.

4. **Party IDs**: The coordinator sends `party_ids` = the T signing parties (e.g., `[node1, node2]`). This matches `ServeMPCInput.AllPartyIDs` from the CLI (`--all-party-ids "node1,node2"`). Verify these are consistent.
