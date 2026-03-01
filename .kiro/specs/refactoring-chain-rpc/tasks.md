# Implementation Plan

- [ ] 0. (P) Move ERC20 contract ABI bindings to pkg/
- [x] 0.1 (P) Move internal/infrastructure/contract/ to pkg/chains/eth/contract/
  - Move `contract.go` (ERC20 contract factory) and `token-abi.go` (auto-generated DO NOT EDIT ABI bindings) into the new `pkg/chains/eth/contract/` directory
  - Neither file imports from `internal/` so no import cleanup is required inside the moved files
  - Update all import paths in files across `internal/` that reference the old `internal/infrastructure/contract` package
  - Confirm the project builds with zero compilation errors after the move
  - _Requirements: 1.4, 1.5, 8.1_

- [ ] 1. (P) Create the Ethereum RPC package in pkg/
- [x] 1.1 (P) Set up the ETH RPC package foundation with client interfaces
  - Define a minimal interface that abstracts the raw JSON-RPC transport so that the concrete Ethereum client satisfies it without modification
  - Define a separate interface for the higher-level typed Ethereum client (used by gas tip cap and raw transaction sending)
  - Establish the package directory and ensure it has no imports from `internal/`
  - _Requirements: 1.1, 1.2, 1.4, 1.5_

- [x] 1.2 (P) Move gas estimation and core ETH blockchain query functions
  - Relocate gas price, gas estimation, block number, balance query, transaction count, syncing status, protocol version, coinbase, accounts, and block/uncle query functions into the new package
  - Move the `ResponseSyncing` wire type (currently defined inside the infrastructure layer) into the package so it is no longer an infrastructure-only type
  - Preserve all existing error wrapping patterns exactly
  - _Requirements: 1.2, 2.1, 2.2, 2.3, 7.2, 8.1_

- [x] 1.3 (P) Move ETH transaction and node management RPC functions
  - Relocate sign, send transaction, send raw transaction, get transaction by hash, and get transaction receipt functions
  - Relocate admin peer management, node info, and data directory functions
  - Relocate personal account management functions (import raw key, new account, list accounts, lock/unlock account)
  - Relocate miner control, network version, net listening, peer count, and web3 utility functions
  - _Requirements: 1.2, 2.1, 2.2, 7.2, 8.1, 8.5_

- [x]* 1.4 Add unit tests for the ETH RPC package
  - Write tests using a mock implementation of the RPC client interface so no live Ethereum node is needed
  - Cover happy-path JSON parsing, hex decoding, and error propagation for each function group
  - _Requirements: 8.4_

- [ ] 2. (P) Create the XRP RPC package in pkg/
- [x] 2.0 Move XRP proto-generated code to pkg/chains/xrp/protogen/
  - Move all six generated files (`account.pb.go`, `account_grpc.pb.go`, `address.pb.go`, `address_grpc.pb.go`, `transaction.pb.go`, `transaction_grpc.pb.go`) from `internal/infrastructure/api/xrp/protogen/` to `pkg/chains/xrp/protogen/`
  - Update `PROTO_GO_OUT_DIR` in `make/codegen_proto.mk` from `internal/infrastructure/api/xrp/protogen` to `pkg/chains/xrp/protogen`
  - Update import paths in the five consumer files (`xrpapi.go`, `xrpapi_tx.go`, `xrpapi_account.go`, `xrpapi_address.go`, `converter.go`) to reference the new package path
  - Confirm the project builds with zero compilation errors after the move
  - _Requirements: 1.5, 7.3, 8.1_

- [x] 2.1 (P) Set up the XRP RPC package foundation with WebSocket client interface
  - Define a minimal interface that abstracts the WebSocket call transport so the existing WebSocket client satisfies it without modification
  - Establish the package directory with no imports from `internal/`
  - _Requirements: 1.1, 1.4, 1.5_

- [x] 2.2 (P) Move XRP public node query functions and their wire types
  - Relocate account info and account channels query functions along with their request/response wire types into the new package
  - Relocate server info query function and its response type
  - Define all wire types using Go-idiomatic field names with JSON struct tags matching the XRP WebSocket protocol
  - _Requirements: 1.2, 2.1, 2.2, 2.3, 7.3, 8.1_

- [x] 2.3 (P) Move XRP admin key generation functions and their wire types
  - Relocate wallet proposal and validation key creation functions and their response wire types into the new package
  - Ensure the wire types do not expose any third-party peersyst library types; define self-contained structs matching the rippled WebSocket protocol
  - _Requirements: 1.2, 2.1, 2.2, 7.3, 8.1, 8.5_

- [x]* 2.4 Add unit tests for the XRP RPC package
  - Write tests using a mock implementation of the WebSocket client interface so no live XRP node is needed
  - Cover JSON marshaling of request payloads and unmarshaling of response types for each function
  - _Requirements: 8.4_

- [ ] 3. Create the Bitcoin RPC package in pkg/
- [x] 3.1 Set up the BTC RPC package foundation with client interface
  - Define a minimal interface that abstracts the raw JSON-RPC call method so the existing Bitcoin RPC client satisfies it without modification
  - Establish the package directory with no imports from `internal/`
  - Move the `FlexibleLabels` custom unmarshaler (handles BTC string-array and BCH object-array label formats) and the `WarningsField` custom unmarshaler into the new package as they are wire-format concerns
  - _Requirements: 1.1, 1.4, 1.5, 7.1, 8.2_

- [x] 3.2 (P) Implement address and validation RPC functions with wire types
  - Move the address info query function and its wire response type into the package, using Go-idiomatic field names and JSON struct tags matching the Bitcoin RPC protocol
  - Move the address validation function and its wire response type
  - Move the get addresses by label function and its intermediate response type
  - Confirm that the `FlexibleLabels` unmarshaler from 3.1 is used by the address info response type
  - _Requirements: 1.2, 2.1, 2.2, 2.3, 7.1, 8.1, 8.2_

- [x] 3.3 (P) Implement network and blockchain info RPC functions with wire types
  - Move the get network info function and its complete wire response type (including nested network, local address, and warnings types) into the package
  - Move the get blockchain info function and its response type (including soft fork types)
  - Move the get block count function
  - _Requirements: 1.2, 2.1, 2.2, 7.1, 8.1_

- [x] 3.4 (P) Implement transaction RPC functions with wire types
  - Move transaction retrieval by ID and its wire response type
  - Move raw transaction decode, create, fund, and sign functions along with their request and response wire types
  - Move UTXO query (get tx out) function
  - Move send transaction by hex and by byte functions
  - _Requirements: 1.2, 2.1, 2.2, 7.1, 8.1_

- [x] 3.5 (P) Implement import, descriptor, wallet, and utility RPC functions with wire types
  - Move address import, private key import, descriptor import, and multi-import functions with their request and response wire types
  - Move descriptor info query and list descriptors functions with their wire types
  - Move wallet lifecycle functions (create, load, unload)
  - Move fee estimation, balance query, label setting, multisig address creation, and logging configuration functions with wire types
  - _Requirements: 1.2, 2.1, 2.2, 7.1, 8.1_

- [x]* 3.6 Add unit tests for the BTC RPC package
  - Write tests using a mock implementation of the RPC caller interface so no live Bitcoin node is needed
  - Cover JSON unmarshaling for the BTC and BCH label format variations (FlexibleLabels), the warnings field format variations (WarningsField), and all wire response types
  - _Requirements: 8.4_

- [ ] 4. Update infrastructure adapters to delegate to the new pkg/ RPC packages
- [x] 4.1 (P) Thin out the ETH infrastructure adapter
  - Replace each ETH adapter method body that contains a raw RPC call with a one-line delegation to the corresponding function in the ETH RPC package
  - Keep the fallback gas tip cap logic (which depends on configuration) in the infrastructure adapter; it calls the pkg function internally
  - Keep client lifecycle methods (close, get chain config, coin type code) unchanged in the adapter
  - Phase 1: adapter methods that return DTOs continue to convert inline for now; no port interface changes yet
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 7.2, 8.3_

- [x] 4.2 (P) Thin out the XRP infrastructure adapter
  - Replace the body of each WebSocket-based public account, server info, and admin key generation method with delegation to the corresponding function in the XRP RPC package
  - Keep the `xrpapi_tx.go` orchestration layer entirely in infrastructure (not moved in Phase 1)
  - Phase 1: adapter methods that return DTOs continue to convert inline; no port interface changes yet
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 7.3, 8.3_

- [x] 4.3 Update the BTC infrastructure adapter to delegate to the BTC RPC package
  - Replace each BTC adapter method body that makes a direct RPC call with delegation to the corresponding function in the BTC RPC package
  - Retain the private inline conversion to DTOs within each adapter method for this phase (port interfaces are unchanged until Phase 2)
  - Confirm that the BCH adapter (which embeds the BTC adapter) continues to work correctly through embedding; no BCH-specific overrides should be needed since `FlexibleLabels` handles the wire format difference
  - Retain all non-RPC logic in the adapter (PSBT operations, UTXO helpers, descriptor parsing, BIP32 utilities)
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 7.1, 7.4, 8.3_

- [x] 4.4 Verify Phase 1 correctness across all chains
  - Build the entire project and confirm zero compilation errors
  - Run the full unit test suite and confirm all tests that passed before Phase 1 still pass
  - Confirm that no `internal/` imports have been introduced in any `pkg/chains/*/rpc/` file
  - _Requirements: 8.1, 8.4, 8.5_

- [ ] 5. Phase 2 - Update ETH port interface signatures and eliminate infrastructure types
- [x] 5.1 Update ETH port interface method signatures to return pkg types
  - Change the `Syncing` method return type from the infrastructure-defined `ResponseSyncing` to the type now living in the ETH RPC package
  - For any other ETH port methods whose return types reference infrastructure-defined structs rather than domain or standard library types, update those signatures to reference the pkg types
  - Verify that the focused ETH port interfaces (`EthNodeAPIClient`, `TxMonitor`, etc.) remain syntactically valid and that all their method names are preserved
  - _Requirements: 3.1, 5.1, 5.2, 5.4_

- [x] 5.2 Update ETH adapter and use cases after interface signature change
  - Update the ETH infrastructure adapter so that the affected methods return the pkg type directly without any intermediate conversion
  - Update any ETH use case files that reference the infrastructure-defined `ResponseSyncing` type to import from the ETH RPC package instead
  - _Requirements: 3.2, 3.3, 4.1, 4.4_

- [x] 5.3 Regenerate ETH mocks and verify ETH tests
  - Run mock generation for all ETH port interfaces affected by signature changes
  - Run the ETH-specific test suite and confirm all tests pass
  - _Requirements: 8.4, 8.5_

- [x] 6. Phase 2 - Eliminate pure-wire BTC DTOs and update port interface signatures
- [x] 6.1 Identify and update BTC port interface method signatures for pure-wire DTOs
  - For each BTC port interface method that currently returns a pure-wire DTO (address info, validate address, network info, blockchain info, transaction result, raw transaction, fund raw transaction, descriptor info, list descriptors, import descriptor response, multisig address, logging result), update the return type to reference the corresponding type now in the BTC RPC package
  - Preserve all interface names (`Bitcoiner`, `AddressOperator`, `NetworkInformer`, etc.) and all method names without change
  - Leave PSBT-related and UTXO-with-domain-fields methods unchanged (those DTOs are domain-enriched and are not being eliminated)
  - _Requirements: 3.1, 5.1, 5.2, 5.4_

- [x] 6.2 Update BTC use cases to import from pkg instead of application DTO package
  - For each BTC use case file that imports address info, network info, transaction result, or other eliminated DTOs, replace the import with the corresponding type from the BTC RPC package
  - Where a use case needs to access a field of what was previously a DTO, confirm the pkg type exposes the same field under the same idiomatic name
  - For the `UnspentOutput` type (domain-enriched, contains account type), create an application-layer converter function that maps the RPC unspent list result plus domain account type into the retained domain type
  - _Requirements: 4.1, 4.2, 4.3, 4.4_

- [x] 6.3 Remove obsolete BTC mapper functions and DTO type definitions
  - Delete all mapper functions in the BTC infrastructure layer that previously converted RPC response structs to DTOs for the pure-wire cases now handled by the pkg types
  - Delete the corresponding DTO type definitions from the application DTO package for types that have been fully replaced by pkg types
  - Confirm that the remaining DTO types (PSBT-related, `UnspentOutput`, `PreviousTx`) are intact and still referenced correctly
  - _Requirements: 3.2, 3.3_

- [x] 6.4 Regenerate BTC mocks and verify BTC tests
  - Run mock generation for all BTC port interfaces affected by signature changes
  - Run the BTC and BCH test suites and confirm all tests pass
  - _Requirements: 8.4, 8.5_

- [x] 7. Phase 2 - Eliminate pure-wire XRP DTOs and update port interface signatures
- [x] 7.1 Update XRP port interface method signatures for pure-wire DTOs
  - For each XRP port interface method that returns a pure-wire DTO (account info, account channels, server info, validation create response, wallet propose response), update the return type to reference the corresponding type now in the XRP RPC package
  - Preserve all XRP interface names (`XRPer`, `XRPAdminer`, `XRPPublicer`, etc.) and all method names
  - Leave methods returning XRP transaction input types unchanged (those are orchestration-level types, not pure wire)
  - _Requirements: 3.1, 5.1, 5.2, 5.4_

- [x] 7.2 Update XRP adapter and use cases, remove obsolete XRP converter functions
  - Update the XRP infrastructure adapter so that the affected methods return the pkg type directly without conversion
  - Update XRP use case files that import the eliminated XRP DTO types to import from the XRP RPC package instead
  - Remove the converter functions in the XRP infrastructure layer that converted WebSocket response types into eliminated application DTOs
  - _Requirements: 3.2, 3.3, 4.1, 4.2, 4.4_

- [x] 7.3 Regenerate XRP mocks and verify XRP tests
  - Run mock generation for all XRP port interfaces affected by signature changes
  - Run the XRP test suite and confirm all tests pass
  - _Requirements: 8.4, 8.5_

- [x] 7.4 Move XRP xrplgo client wrapper to pkg/chains/xrp/client/
  - Move `internal/infrastructure/api/xrp/xrplgo/` (`client.go`, `account.go`, `transaction.go`, `ledger.go`) to `pkg/chains/xrp/client/`
  - At this point `account.go` and `transaction.go` no longer import `internal/application/dto/xrp` (those DTO types were eliminated in task 7.2), so the move is clean
  - Update the package declaration from `xrplgo` to `client`
  - Update import paths in all XRP infrastructure files that currently reference the `xrplgo` package
  - Move `client_test.go` alongside the package; update its imports to reflect the new location
  - Confirm the project builds with zero compilation errors after the move
  - _Requirements: 1.4, 1.5, 7.3, 8.1_

- [x] 8. Final integration validation and cleanup
- [x] 8.1 Remove any remaining obsolete infrastructure type definitions and empty DTO packages
  - Remove RPC response struct definitions that were previously defined inline in infrastructure files and have now been fully replaced by types in the pkg packages
  - Remove or consolidate any DTO package files that are now empty or contain only domain-enriched types that have been moved elsewhere
  - Confirm no infrastructure file imports from `pkg/chains/*/rpc/` is duplicating a type that already lives in pkg
  - _Requirements: 2.4, 3.3_

- [x] 8.2 Full project build and test verification
  - Build the complete project and confirm zero compilation errors
  - Run the full test suite across all chains and confirm all tests that passed before the refactoring still pass
  - Confirm that use case files contain no imports from `internal/infrastructure/api/` (the architectural boundary rule is now structurally enforced)
  - Run the linter and resolve any issues introduced during the refactoring
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_
