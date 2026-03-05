# Requirements Document

## Project Description (Input)

@internal/infrastructure/api/xrp/xrp.go

```go
// current
type XRP struct {
 *WSClient                           // WebSocket operations (public + admin)
 API          *xrplclient.XRPLClient // gRPC operations (legacy, being phased out)
 chainConf    *chaincfg.Params
 coinTypeCode domainCoin.CoinTypeCode
}

// after (struct divided)
type publicClient struct {
 public       *websocket.WS
 chainConf    *chaincfg.Params // if needed
 coinTypeCode domainCoin.CoinTypeCode // if needed
}

type adminClient struct {
 admin        *websocket.WS
 chainConf    *chaincfg.Params // if needed
 coinTypeCode domainCoin.CoinTypeCode // if needed
}

type txClient struct {
 txClient      protogen.XRPTransactionAPIClient
 chainConf     *chaincfg.Params // if needed
 coinTypeCode  domainCoin.CoinTypeCode // if needed
}

type accountClient struct {
 client protogen.XRPAccountAPIClient
 chainConf     *chaincfg.Params // if needed
 coinTypeCode  domainCoin.CoinTypeCode // if needed
}

type addressClient struct {
 client protogen.XRPAddressAPIClient
 chainConf     *chaincfg.Params // if needed
 coinTypeCode  domainCoin.CoinTypeCode // if needed
}
```

The `XRP` struct has too many responsibilities and a name that does not communicate its role:

- `API *xrplclient.XRPLClient` bundles three gRPC concerns (transaction, account, address) into a struct whose primary role is WebSocket communication.
- The name `XRP` identifies only the chain, not the entity's purpose.
- `WSClient` bundles public and admin connections even when a use case needs only one.

The proposal is to **fully decompose** the `XRP` struct into five single-responsibility structs, each
named after its role. Each struct implements exactly one port interface and is injected directly into the
use cases that need it. The monolithic `XRP` struct is eliminated.

**Prerequisites:**

- `apps/xrpl-grpc-server` will remain and will not be modified.
- `pkg/chains/xrp/xrplclient/client.go` will not be modified.
- Library selection for new implementations must follow `.claude/rules/chains/xrp/3rd-library.md`.

## Requirements

### Introduction

The current `XRP` struct carries five distinct responsibilities:

1. **Public WebSocket communication** — querying the public XRP node
2. **Admin WebSocket communication** — sending admin commands to the node
3. **Transaction operations** — via `API.TxClient` (gRPC, currently non-functional)
4. **Address operations** — via `API.AddressClient` (gRPC, currently non-functional)
5. **Account operations** — via `API.AccountClient` (gRPC, currently non-functional)

Each responsibility shall become its own struct in `internal/infrastructure/api/xrp/`:

| Struct | Responsibility | Implements |
|--------|---------------|-----------|
| `publicClient` | Public WebSocket operations | public-facing port interfaces |
| `adminClient` | Admin WebSocket operations | admin port interfaces |
| `txClient` | Transaction preparation, signing, combining | `protogen.XRPTransactionAPIClient` |
| `accountClient` | Account queries | `protogen.XRPAccountAPIClient` |
| `addressClient` | Address generation and validation | `protogen.XRPAddressAPIClient` |

The `XRP` struct and `WSClient` struct are eliminated. The DI container wires each struct
independently to the use cases that need it.

The gRPC path remains available: `*xrplclient.XRPLClient` sub-fields already satisfy the
`protogen.*` interfaces and can be substituted directly in the DI container.

Library selection for non-gRPC implementations must follow `.claude/rules/chains/xrp/3rd-library.md`;
`github.com/XRPLF/xrpl-go` is **not** an approved dependency in this project.

---

### Requirement 1: Eliminate the `XRP` Struct and Replace with Five Single-Responsibility Structs

**Objective:** As a developer, I want the monolithic `XRP` struct replaced by five focused structs,
each with a name and responsibility that are immediately understandable, so that each use case depends
only on what it actually needs.

#### Acceptance Criteria

1. The XRP Infrastructure module shall delete the `XRP` struct and the `NewXRP` constructor from `internal/infrastructure/api/xrp/`.
2. The XRP Infrastructure module shall create five unexported structs in `internal/infrastructure/api/xrp/`: `publicClient`, `adminClient`, `txClient`, `accountClient`, and `addressClient`.
3. Each struct shall carry only the fields it genuinely requires; fields that a struct does not use shall not be included (e.g., a struct that never inspects `chainConf` shall not hold it).
4. Each struct shall have an exported constructor (e.g., `NewPublicClient`, `NewTxClient`) usable by the DI container.
5. The `WSClient` struct shall be eliminated or reduced to a pure connection-lifecycle helper (open/close) that is not exposed to use cases.

---

### Requirement 2: Each Struct Implements Exactly One Port Interface

**Objective:** As a developer, I want each new struct to implement a clearly defined port interface,
so that use cases can depend on the minimum surface they need.

#### Acceptance Criteria

1. `publicClient` shall implement the port interface(s) in `internal/application/ports/api/xrp/` that cover public-node WebSocket operations.
2. `adminClient` shall implement the port interface(s) in `internal/application/ports/api/xrp/` that cover admin-node WebSocket operations.
3. `txClient` shall implement `protogen.XRPTransactionAPIClient`, covering at minimum `PrepareTransaction`, `SignTransaction`, and `CombineTransaction`.
4. `accountClient` shall implement `protogen.XRPAccountAPIClient` for any account methods consumed by use cases.
5. `addressClient` shall implement `protogen.XRPAddressAPIClient`, covering at minimum `GenerateAddress`, `GenerateXAddress`, and `IsValidAddress`.
6. No struct shall implement methods belonging to another struct's responsibility.

---

### Requirement 3: Provide Non-gRPC Implementations for `txClient`, `accountClient`, and `addressClient`

**Objective:** As a developer, I want the three protogen-interface structs to work without the gRPC
server, so that all XRP wallet operations are functional today.

#### Acceptance Criteria

1. `txClient` shall execute transaction preparation, signing, and combining locally without making gRPC calls.
2. `accountClient` shall execute account queries via the approved WebSocket path (`pkg/chains/xrp/xrplgo`) or locally, without making gRPC calls.
3. `addressClient` shall execute address generation and validation locally using existing `pkg/chains/xrp/` utilities, without making gRPC calls.

4. Library selection must comply with `.claude/rules/chains/xrp/3rd-library.md`:
   - Offline cryptographic operations → `github.com/Peersyst/xrpl-go`
   - Live node queries → `pkg/chains/xrp/xrplgo` (wraps `github.com/xrpscan/xrpl-go`)
   - `github.com/XRPLF/xrpl-go` shall **not** be added as a new dependency.

---

### Requirement 4: Wire Each Struct Independently from the DI Container

**Objective:** As a developer, I want each use case to receive only the struct(s) it needs via the
DI container, so that dependencies are explicit and minimal.

#### Acceptance Criteria

1. The DI container (`internal/di/container.go`) shall instantiate each of the five structs independently.
2. Each use case constructor shall declare only the port interfaces it uses as parameters; it shall not receive the former `XRP` or `XRPer` monolith.
3. Use cases that previously depended on `XRP` (or `XRPer`) to reach gRPC operations shall be updated to accept the relevant focused interface as a separate parameter.
4. All references to `XRP`, `NewXRP`, and the old `XRPer` monolith shall be removed from `internal/` after the decomposition.

---

### Requirement 5: Preserve Backward Compatibility with the gRPC Path

**Objective:** As a developer, I want the gRPC-backed `*xrplclient.XRPLClient` to remain injectable
so that the gRPC server can be re-enabled in the future by changing only DI wiring.

#### Acceptance Criteria

1. The XRP Infrastructure module shall not modify any file under `apps/xrpl-grpc-server/`.
2. The XRP Infrastructure module shall not modify `pkg/chains/xrp/xrplclient/client.go`.
3. When `xrplclient.XRPLClient.TxClient`, `.AccountClient`, and `.AddressClient` are passed to the respective struct constructors instead of the non-gRPC implementations, the system shall compile and behave identically to the current gRPC-backed behavior.
4. Switching between gRPC-backed and non-gRPC implementations shall require changes only to DI wiring.

---

### Requirement 6: No Regression in Build and Tests

**Objective:** As a developer, I want the full build and test suite to continue passing after this
refactoring, so that no existing functionality is broken.

#### Acceptance Criteria

1. When `make check-build` is run after the refactoring, the XRP module shall compile without errors.
2. When `make go-lint` is run after the refactoring, the XRP module shall produce no new lint errors.
3. When `make go-test` is run, all unit tests in `internal/infrastructure/api/xrp/` and `pkg/chains/xrp/` shall pass.
4. The refactoring shall not break any existing port interface definitions in `internal/application/ports/api/xrp/`.
