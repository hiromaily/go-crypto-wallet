# Requirements Document

## Project Description (Input)

Functions specifically designed to call Node RPCs should be moved to an RPC directory created under each chain directory under pkg/.

For example, GetAddressInfo in internal/infrastructure/api/btc/btc/address.go calls RPC as follows:

```go
rawResult, err := b.Client.RawRequest("getaddressinfo", []json.RawMessage{input})
```

This function should be moved to pkg/chains/btc/rpc.

If this process involves conversion using DTOs, they need to be removed. Define the RPC request and response types under pkg as they are, and only define a domain type if a use case requires a different type, and convert them using DTOs in the application layer.

While much of the implementation will be moved under pkg, the abstraction layer will remain in the application/ports layer.

## Introduction

This specification defines requirements for relocating chain-specific Node RPC call functions from the infrastructure layer (`internal/infrastructure/api/`) into dedicated `pkg/chains/*/rpc/` packages. The goal is to separate pure RPC communication logic from infrastructure wiring, eliminate premature DTO conversion at the infrastructure layer, and ensure that type conversion between RPC types and domain types happens only in the application layer where domain context exists.

The refactoring applies to all three supported blockchain chains: Bitcoin (BTC/BCH), Ethereum (ETH/ERC-20), and Ripple (XRP).

## Requirements

### Requirement 1: RPC Package Structure

**Objective:** As a developer, I want RPC call functions organized in dedicated `pkg/chains/*/rpc/` packages, so that RPC communication logic is decoupled from infrastructure wiring and reusable across the codebase.

#### Acceptance Criteria

1. The RPC Package shall create the following directories: `pkg/chains/btc/rpc/`, `pkg/chains/eth/rpc/`, and `pkg/chains/xrp/rpc/`.
2. The RPC Package shall place all functions that directly invoke blockchain Node RPC methods (via `RawRequest`, `CallContext`, or WebSocket `Call`) into the appropriate chain's `rpc/` directory.
3. When a chain (e.g., BCH) shares the same RPC protocol as another chain (e.g., BTC), the RPC Package shall reuse the BTC rpc package rather than duplicating code.
4. The RPC Package shall expose only exported, standalone functions that accept the RPC client as an explicit parameter (not as a struct method on an infrastructure type).
5. If an RPC function requires chain-specific parameters (e.g., config), the RPC Package shall accept them as explicit function parameters.

---

### Requirement 2: RPC Type Definitions in pkg

**Objective:** As a developer, I want RPC request and response types defined in the `pkg/chains/*/rpc/` packages, so that type definitions are co-located with the RPC calls that use them and are not duplicated across layers.

#### Acceptance Criteria

1. The RPC Type system shall define all raw RPC request struct types (e.g., `GetAddressInfoParams`) in the `pkg/chains/*/rpc/` package where the corresponding RPC call function resides.
2. The RPC Type system shall define all raw RPC response struct types (e.g., `GetAddressInfoResult`) in the `pkg/chains/*/rpc/` package where the corresponding RPC call function resides.
3. The RPC Type system shall use JSON struct tags that match the exact field names expected by the Node RPC protocol.
4. If a type is already defined under `pkg/` (e.g., from a library), the RPC Package shall use or alias that type rather than redefining it.
5. The RPC Type system shall NOT define application-layer DTO types (types shaped for use cases) in `pkg/chains/*/rpc/`.

---

### Requirement 3: Elimination of Infrastructure-Layer DTO Conversion

**Objective:** As a developer, I want the infrastructure layer to return raw RPC types rather than application DTOs, so that the infrastructure layer does not carry application-layer concerns and DTOs are only created where domain context justifies them.

#### Acceptance Criteria

1. When an infrastructure API method (e.g., `Bitcoin.GetAddressInfo`) is refactored, the RPC Package shall return the raw RPC response type defined in `pkg/chains/*/rpc/`, not an application DTO.
2. If infrastructure-layer mapper functions (e.g., `ToAddressInfo()`) exist solely to convert from an RPC response type to an application DTO, the Refactoring shall remove those mapper functions from the infrastructure layer.
3. The Refactoring shall remove application DTO types from the infrastructure package when they exist only to transport RPC response data to use cases.
4. While the infrastructure-layer DTO conversion is removed, the application/ports interface method signatures may temporarily use RPC types until the application layer conversion is put in place.

---

### Requirement 4: Application-Layer Type Conversion

**Objective:** As a developer, I want type conversion between RPC types and domain/use-case types to happen in the application layer, so that domain rules and conversions are maintained in the correct architectural layer.

#### Acceptance Criteria

1. When a use case requires a type that differs structurally from the raw RPC response, the Application Layer shall define the needed type in the domain or application layer (not in `pkg/chains/*/rpc/`).
2. The Application Layer shall perform conversion from `pkg/chains/*/rpc/` types to domain/use-case types using DTO functions defined in the application layer.
3. If a raw RPC response type is structurally identical to what a use case needs, the Application Layer shall use the `pkg/chains/*/rpc/` type directly without wrapping it in a redundant DTO.
4. The Application Layer shall not import types from `internal/infrastructure/api/` directly; use cases shall only depend on port interfaces and `pkg/` types.

---

### Requirement 5: Port Interface Preservation

**Objective:** As a developer, I want the abstraction interfaces in `application/ports/api/` to remain stable, so that use cases are not affected by the infrastructure reorganization.

#### Acceptance Criteria

1. The Refactoring shall not remove or rename any existing interface types in `internal/application/ports/api/{btc,eth,xrp}/`.
2. When port interface method return types reference application DTOs that are being eliminated, the Port Interface shall be updated to reference the appropriate `pkg/chains/*/rpc/` type or a new application-layer type.
3. The Refactoring shall keep the DI layer (`internal/di/`) as the only location that depends on the concrete infrastructure struct types.
4. The Refactoring shall keep all focused port interfaces (e.g., `AddressOperator`, `BalanceChecker`) intact; only method signatures may change when DTO types are replaced.

---

### Requirement 6: Infrastructure Layer Becomes a Thin Adapter

**Objective:** As a developer, I want the infrastructure API structs to act as thin adapters that delegate to `pkg/chains/*/rpc/` functions, so that infrastructure concerns are minimal and the logic lives in the portable `pkg/` layer.

#### Acceptance Criteria

1. After refactoring, each method on an infrastructure struct (e.g., `Bitcoin`, `Ethereum`, `XRP`) that previously contained RPC call logic shall delegate to a corresponding function in `pkg/chains/*/rpc/` by passing its client field.
2. The infrastructure struct shall retain only the fields and methods needed for client lifecycle management and DI wiring; pure RPC logic shall reside in `pkg/chains/*/rpc/`.
3. If an infrastructure method contains logic beyond RPC call and return (e.g., retry, caching, fallback), that logic shall remain in the infrastructure layer, calling the `pkg/chains/*/rpc/` function internally.
4. The infrastructure layer shall not duplicate RPC type definitions; it shall import them from `pkg/chains/*/rpc/`.

---

### Requirement 7: Multi-Chain Coverage

**Objective:** As a developer, I want the relocation applied consistently to all supported chains, so that the codebase follows a uniform structure regardless of chain.

#### Acceptance Criteria

1. The Refactoring shall apply the RPC relocation to all Bitcoin RPC methods currently in `internal/infrastructure/api/btc/btc/` (including `getaddressinfo`, `getaddressesbylabel`, `validateaddress`, `getbalance`, `estimatesmartfee`, `getnetworkinfo`, `getblockchaininfo`, `gettransaction`, `gettxout`, `decoderawtransaction`, `createrawtransaction`, `fundrawtransaction`, `signrawtransactionwithwallet`, `importaddress`, `importdescriptors`, `importmulti`, `getdescriptorinfo`, `createwallet`, `listdescriptors`, `setlabel`, `logging`, `addmultisigaddress`).
2. The Refactoring shall apply the RPC relocation to all Ethereum RPC methods currently in `internal/infrastructure/api/eth/eth/rpc_*.go` (including all `eth_*`, `net_*`, `web3_*`, `personal_*`, `admin_*`, `miner_*` methods).
3. The Refactoring shall apply the RPC relocation to all XRP WebSocket-based RPC methods currently in `internal/infrastructure/api/xrp/` (including `account_info`, `server_info`, `ledger`, `sign`, `submit`, `validation_create`, `wallet_propose` and related methods).
4. The Refactoring shall ensure BCH chain uses BTC's `pkg/chains/btc/rpc/` package functions since BCH shares the same Bitcoin RPC protocol.

---

### Requirement 8: Behavioral Equivalence

**Objective:** As a developer, I want the refactored code to produce identical observable behavior, so that existing functionality is not broken by the structural changes.

#### Acceptance Criteria

1. When an RPC function is relocated, the Refactoring shall preserve all existing error wrapping patterns (e.g., `fmt.Errorf("fail to call %s: %w", ...)`) so that error messages remain consistent.
2. The Refactoring shall preserve any special unmarshaling logic (e.g., BTC's `FlexibleLabels` for BTC/BCH compatibility) in the `pkg/chains/*/rpc/` package.
3. The Refactoring shall preserve any fallback or retry logic that exists in the infrastructure layer.
4. When the existing tests pass before the refactoring, the Refactoring shall ensure they continue to pass after.
5. The Refactoring shall not change any public RPC method name or behavior observable by use cases.
