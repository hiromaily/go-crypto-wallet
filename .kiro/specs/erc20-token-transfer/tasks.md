# Implementation Plan

## Overview

5 major tasks covering domain registration (add HYC + remove generic ERC20 type), config, EIP-1559 infrastructure upgrade, DI wiring, and tests.
Tasks 1, 2, and 3 can be executed in parallel (different files, no data dependency).
Task 4 depends on Tasks 1 and 3. Task 5 is sequential (depends on Task 4).
Task 5.3 deploys the HYC contract from `apps/eth-contracts` to a local Anvil node.
Task 5.4 runs a Go wallet integration test against the live deployed contract.

---

- [x] 1. Register HYC token and clean up generic ERC-20 type in the domain type system

- [x] 1.1 Add HYC constants to the domain type system
  - Add `CoinTypeERC20HYC CoinType = 9002` to the `CoinType` constant block (follows the existing HYT = 9001 pattern)
  - Add `HYC CoinTypeCode = "hyc"` to the `CoinTypeCode` constant block and register `HYC: CoinTypeERC20HYC` in `CoinTypeCodeValue` map
  - Add `TokenHYC ERC20Token = "hyc"` to the `ERC20Token` constant block and register `TokenHYC: true` in `ERC20Map`
  - Confirm existing HYT and BAT registrations are unchanged
  - Verify that `IsCoinTypeCode("hyc")`, `IsERC20Token("hyc")`, and `IsETHGroup("hyc")` all return `true`
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7_

- [x] 1.2 Remove the generic `CoinTypeERC20` and `ERC20 CoinTypeCode` constants
  - Remove `CoinTypeERC20 CoinType = 9000` from the `CoinType` constant block and remove its `ERC20: CoinTypeERC20` entry from `CoinTypeCodeValue`
  - Remove `ERC20 CoinTypeCode = "erc20"` from the `CoinTypeCode` constant block
  - Update `IsETHGroup` to remove the `val == ERC20` arm — the function becomes `val == ETH || IsERC20Token(val.String())`
  - In `internal/infrastructure/wallet/key/strategy/factory.go`, change `case domainCoin.ETH, domainCoin.ERC20:` to `case domainCoin.ETH:` and add a new arm using `IsETHGroup` (or list specific token codes) to return `NewETHKeyStrategy()` for ERC-20 tokens
  - Remove `domainCoin.ERC20` from all "not implemented" guard clauses in BTC-specific switch statements — affected files: `internal/di/container.go` (2 switch arms), `internal/infrastructure/api/btc/connection.go`, `internal/infrastructure/api/btc/btc/multisig.go`, `internal/application/usecase/watch/btc/import_address.go`, `internal/application/usecase/keygen/btc/import_private_key.go`, `internal/application/usecase/sign/btc/import_private_key.go`
  - Run `make check-build` to verify no compilation errors after removal
  - _Requirements: 1.1, 1.2, 1.3_

- [ ] 2. Add HYC token entries to wallet configuration files

- [ ] 2.1 (P) Add HYC entry to watch wallet configuration
  - Add `hyc` block under `ethereum.erc20s` in `config/wallet/eth/watch.yaml`
  - Include `symbol: "hyc"`, `name: "HY Coin"`, `contract_address: "0x..."` (placeholder for deployed HYC contract), `master_address: "0x..."` (placeholder for fee-paying master account), `decimals: 18`
  - Note: placeholder addresses are replaced in Task 5.3 when `forge script` deploys HYC from `apps/eth-contracts/script/DeployHYC.s.sol` to Anvil; production operators replace them after a real deployment
  - Does not conflict with existing `hyt` entry; both coexist in the same map
  - _Requirements: 2.1, 2.4, 2.5_

- [ ] 2.2 (P) Add HYC entries to keygen and sign wallet configurations
  - Add `erc20s.hyc` block to `config/wallet/eth/keygen.yaml` (currently has no `erc20s` section)
  - Add identical `erc20s.hyc` block to `config/wallet/eth/sign.yaml` (currently has no `erc20s` section)
  - Both entries use the same fields and placeholder values as Task 2.1
  - The `erc20_token` key can remain unset in keygen/sign configs (watch wallet drives token selection)
  - _Requirements: 2.2, 2.3, 2.4, 2.5_

- [ ] 3. (P) Upgrade ERC-20 infrastructure to support EIP-1559 transactions

- [ ] 3.1 (P) Add Ethereum named field to ERC20 struct and update constructor
  - Add `eth *eth.Ethereum` as a private named field in the `ERC20` struct (not Go anonymous embedding — avoids promoting all Ethereum methods)
  - Add `eth *eth.Ethereum` as the first parameter of `NewERC20`
  - Remove the now-redundant direct `*ethclient.Client` field from `ERC20` if the client can be accessed via `e.eth` — or retain it for the existing balance and token contract calls that use it directly (assess during implementation; retain `client *ethclient.Client` if removing it causes too many cascading changes)
  - Confirm the compile-time check `var _ apieth.ERC20er = (*ERC20)(nil)` still passes
  - _Requirements: 3.5_

- [ ] 3.2 Replace hardcoded SupportsEIP1559 with real detection
  - Replace the `return false` body of `ERC20.SupportsEIP1559` with a delegation to `e.eth.SupportsEIP1559(ctx)`
  - This reuses the existing Anvil detection and `baseFeePerGas` block-header check already implemented in `Ethereum.SupportsEIP1559`
  - Depends on Task 3.1 (the `eth` field must exist)
  - _Requirements: 3.2, 3.4_

- [ ] 3.3 Implement EIP-1559 transaction creation for ERC-20 transfers
  - Replace the `CreateRawTransactionEIP1559` delegation to `CreateRawTransaction` with a real EIP-1559 implementation
  - Build a `types.DynamicFeeTx` with the ABI-encoded `transfer(address,uint256)` calldata in the `Data` field (method selector `0xa9059cbb`) — same calldata as the existing legacy path
  - Set `GasTipCap` from `e.eth.SuggestGasTipCap(ctx)` and `GasFeeCap` from `(baseFee × 2) + tip` — same fee formula as `Ethereum.CreateRawTransactionEIP1559`
  - Set `To` to the ERC-20 contract address (not the token recipient address — same pattern as the legacy path)
  - Populate `TxCreateParams` with `EthTxType = 2`, `ChainID`, `MaxFeePerGas`, `MaxPriorityFeePerGas`
  - Fall back to `CreateRawTransaction` (legacy Type 0) when `SupportsEIP1559` returns `false`
  - Return an error if `SupportsEIP1559` returns `true` but EIP-1559 fee fields cannot be obtained
  - Depends on Task 3.1 and 3.2
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.6, 4.2, 4.3, 4.4, 4.5_

- [ ] 3.4 Resolve nonce duplication by delegating to Ethereum field
  - Replace the body of `ERC20.getNonce` with a call to `e.eth`'s nonce retrieval (the existing FIXME on line 279 of `erc20.go`)
  - This removes the duplicate `client.PendingNonceAt` logic
  - Depends on Task 3.1
  - _Requirements: 3.5_

- [ ] 4. Update DI container to wire Ethereum instance into ERC-20 constructor

- [ ] 4.1 Pass cached Ethereum instance to NewERC20 in the DI container
  - In `container.go`, update `newERC20()` to call `c.newETH()` and pass the returned `*eth.Ethereum` as the first argument to `NewERC20`
  - `c.newETH()` is already lazy-cached, so no duplicate connections are created
  - Confirm no circular initialization: `newETH()` must not call `newERC20()` in any code path
  - Depends on Tasks 1 and 3 (HYC registered and new constructor signature ready)
  - _Requirements: 3.5, 4.1_

- [ ] 4.2 Verify build and routing for HYC and existing ERC-20 tokens
  - Run `make check-build` to confirm the updated constructor compiles cleanly across all wallet binaries (watch, keygen, sign)
  - Confirm that existing `hyt` token routing is unaffected by tracing `IsERC20Token("hyt")` through the DI dispatch
  - Confirm that `IsETHGroup("hyc")` now routes through `newETHWatchCreateTransactionUseCase` → `newERC20()` path when `erc20_token: "hyc"` is set in config
  - Run `make go-lint` to confirm no new lint issues
  - Depends on Task 4.1
  - _Requirements: 4.1, 5.1, 6.1, 7.1_

- [ ] 5. Write unit tests for domain registration and EIP-1559 ERC-20 infrastructure

- [ ] 5.1 Write domain registration unit tests
  - Add tests to `internal/domain/coin/` verifying that `IsCoinTypeCode("hyc")` returns `true`
  - Add test verifying `IsERC20Token("hyc")` returns `true`
  - Add test verifying `IsETHGroup("hyc")` returns `true`
  - Add test verifying `BIP44AccountPath(HYC, 0)` produces the path `"m/44'/9002'/0'"`
  - Add regression tests confirming existing `HYT`, `BAT`, `ETH` behavior is unchanged
  - Depends on Task 1
  - _Requirements: 8.3, 8.4, 8.5_

- [ ] 5.2 Write ERC-20 EIP-1559 infrastructure unit tests
  - Add tests to `internal/infrastructure/api/eth/erc20/` using a mock `*eth.Ethereum` or a test double
  - Test `SupportsEIP1559`: when the Ethereum field returns `true`, `ERC20.SupportsEIP1559` returns `true`
  - Test `SupportsEIP1559`: when the Ethereum field returns `false`, `ERC20.SupportsEIP1559` returns `false`
  - Test `CreateRawTransactionEIP1559`: when EIP-1559 is supported, the returned `TxCreateParams.EthTxType` equals `2` and `ChainID` is non-zero
  - Test `CreateRawTransactionEIP1559`: when EIP-1559 is not supported, the method falls back to `CreateRawTransaction` and `TxCreateParams.EthTxType` equals `0`
  - Test `CreateRawTransactionEIP1559`: the encoded transaction's `Data` field starts with the `0xa9059cbb` method selector
  - Depends on Task 3
  - _Requirements: 8.1, 8.2, 8.4, 8.5_

- [ ] 5.3 Deploy HYC contract to local Anvil node for integration testing
  - Ensure an Anvil node is running locally (`anvil` from the Foundry toolchain)
  - From `apps/eth-contracts/`, run `forge script script/DeployHYC.s.sol --rpc-url http://localhost:8545 --broadcast` with a funded test private key set in `PRIVATE_KEY`
  - Capture the deployed HYC contract address from the script output or `broadcast/` directory artifacts
  - Update `config/wallet/eth/watch.yaml` (and test config) with the real deployed `contract_address` and a funded `master_address` from Anvil's pre-funded accounts
  - Verify the deployed contract by calling `balanceOf` on the deployer address; expect `1_000_000 ether` total supply (as defined in `DeployHYC.s.sol`)
  - Depends on Task 2 (config structure must exist to update)
  - _Requirements: 2.1, 2.4, 2.5_

- [ ] 5.4 Integration test: Go wallet ERC-20 transfer against deployed HYC contract
  - Using the Anvil node and deployed HYC contract from Task 5.3, add an integration test scenario in `internal/application/usecase/watch/eth/` (alongside `eth_file_exchange_test.go`)
  - Configure the test with `erc20_token: "hyc"`, the deployed contract address, and Anvil pre-funded account addresses
  - Run the full create-unsigned-tx → sign-tx → send-tx flow for a HYC token transfer:
    - `CreateTransactionUseCase` with coin type `hyc` — verify unsigned JSON contains `EthTxType = 2` (EIP-1559) and `Data` starting with `0xa9059cbb`
    - `SignTransactionUseCase` — verify signed JSON is produced without network calls
    - `SendTransactionUseCase` — broadcast the signed transaction to Anvil
  - After send, call `GetBalance` on the recipient address to confirm the token balance increased by the transferred amount
  - Verify token balance decreased on the sender address by the same amount
  - Depends on Tasks 1, 2, 3, 4, 5.1, 5.2, 5.3
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 5.1, 5.2, 5.3, 6.1, 6.2, 7.1, 8.4_
