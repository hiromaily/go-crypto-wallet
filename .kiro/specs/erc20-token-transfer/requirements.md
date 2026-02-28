# Requirements Document

## Introduction

This specification defines the requirements for implementing **ERC-20 token transfer** support
in the go-crypto-wallet Go application, using the HYC token (deployed via the `erc20-token` spec)
as the target token.

The scope covers:

1. Registering the HYC token in the Go domain layer and wallet configuration
2. Adding EIP-1559 transaction support to the ERC-20 infrastructure
3. Implementing CLI commands for the ERC-20 token transfer workflow (keygen/sign/watch pattern)
4. End-to-end transfer flow: create unsigned transaction → offline sign → broadcast

This spec is a **Go-only** change. Solidity contracts are out of scope (handled by the `erc20-token` spec).

---

## Domain Type System Context

The domain layer (`internal/domain/coin/types.go`) uses two separate type systems for ERC-20 tokens:

| Type | Purpose | Examples |
|------|---------|---------|
| `CoinType` (uint32) | BIP-44 HD derivation path index per coin/token | `CoinTypeERC20 = 9000`, `CoinTypeERC20HYT = 9001` |
| `CoinTypeCode` (string) | Human-readable coin identifier used in config and CLI | `"erc20"`, `"hyt"` |
| `ERC20Token` (string) | Token identifier for contract dispatch and validation | `TokenHYT = "hyt"` |

**Key insight**: `CoinTypeCode` and `ERC20Token` share the same string value for a given token (e.g., both `"hyt"`),
but serve different roles. `CoinTypeCode` drives BIP-44 key derivation; `ERC20Token` identifies which
contract address and metadata to load from config. Both registrations are required for a token to be
fully operational.

---

## Requirements

### Requirement 1: HYC Token Domain Registration

**Objective:** As a developer, I want the HYC token registered in all three domain type systems,
so that key derivation, config loading, and token validation all work consistently.

#### Acceptance Criteria

1. The ERC-20 Token Transfer system shall add `CoinTypeERC20HYC CoinType = 9002` to the `CoinType` constants in `internal/domain/coin/types.go`.
2. The ERC-20 Token Transfer system shall add `HYC CoinTypeCode = "hyc"` to the `CoinTypeCode` constants and register `HYC: CoinTypeERC20HYC` in `CoinTypeCodeValue`.
3. The ERC-20 Token Transfer system shall add `TokenHYC ERC20Token = "hyc"` to the `ERC20Token` constants and register `TokenHYC: true` in `ERC20Map`.
4. When `IsCoinTypeCode("hyc")` is called, the domain system shall return `true`.
5. When `IsERC20Token("hyc")` is called, the domain system shall return `true`.
6. When `IsETHGroup("hyc")` is called, the domain system shall return `true` (because `IsERC20Token` is called internally).
7. The ERC-20 Token Transfer system shall maintain backward compatibility with existing `HYT` and `BAT` registrations.

---

### Requirement 2: HYC Token Wallet Configuration

**Objective:** As an operator, I want the HYC token configured in all relevant wallet config files,
so that keygen, sign, and watch wallets can operate with HYC tokens.

#### Acceptance Criteria

1. The ERC-20 Token Transfer system shall add an `hyc` entry under `ethereum.erc20s` in `config/wallet/eth/keygen.yaml`, including `symbol`, `name`, `contract_address`, `master_address`, and `decimals` fields.
2. The ERC-20 Token Transfer system shall add the same `hyc` entry to `config/wallet/eth/sign.yaml`.
3. The ERC-20 Token Transfer system shall add the same `hyc` entry to `config/wallet/eth/watch.yaml`.
4. When `ethereum.erc20_token` is set to `"hyc"` in a config file, the wallet shall load HYC token parameters correctly.
5. The `decimals` field for HYC shall be set to `18` to match the deployed contract.

---

### Requirement 3: EIP-1559 Transaction Support for ERC-20

**Objective:** As a developer, I want ERC-20 token transfers to support EIP-1559 (Type 2) transactions,
so that gas fees are handled efficiently on modern Ethereum networks.

#### Acceptance Criteria

1. When the connected Ethereum node supports EIP-1559, the ERC-20 infrastructure shall construct Type 2 (EIP-1559) transactions for token transfers.
2. The ERC-20 infrastructure shall implement `SupportsEIP1559()` to query the node's fee history and return `true` when EIP-1559 is available, replacing the current hardcoded `false` return.
3. When EIP-1559 is supported, the ERC-20 infrastructure shall populate `MaxFeePerGas` and `MaxPriorityFeePerGas` from the wallet configuration.
4. When EIP-1559 is not supported by the node, the ERC-20 infrastructure shall fall back to legacy (Type 0) transactions.
5. The ERC-20 infrastructure shall reuse the existing `Ethereum` struct's EIP-1559 helper methods (via embedding or composition) to eliminate the code duplication noted in the existing FIXME comments.
6. If `MaxFeePerGas` or `MaxPriorityFeePerGas` config values are zero when EIP-1559 is active, the ERC-20 infrastructure shall return an error rather than broadcast a zero-fee transaction.

---

### Requirement 4: Create Unsigned ERC-20 Transfer Transaction (Watch Wallet)

**Objective:** As a watch-wallet operator, I want to create unsigned ERC-20 token transfer transactions,
so that they can be passed to the offline sign wallet for signing.

#### Acceptance Criteria

1. When the `create-tx` watch-wallet CLI command is invoked with coin type `hyc`, the Watch Wallet CLI shall invoke `CreateTransactionUseCase` with the ERC-20 implementation.
2. When a transfer transaction is created, the ERC-20 Token Transfer system shall write an unsigned transaction JSON file to the configured transaction directory.
3. The unsigned transaction JSON shall include `from_address`, `to_address`, `amount`, `nonce`, `chain_id`, and the ABI-encoded `data` field (using method selector `0xa9059cbb`).
4. When multiple transfer targets exist in the database, the ERC-20 Token Transfer system shall batch them into a single JSON file with incrementing nonces.
5. If the sender account balance is insufficient to cover the gas fees, the ERC-20 Token Transfer system shall return an error and write no file.

---

### Requirement 5: Sign ERC-20 Transfer Transaction (Keygen/Sign Wallet)

**Objective:** As an offline sign-wallet operator, I want to sign ERC-20 token transfer transactions
without a network connection, so that private keys remain isolated.

#### Acceptance Criteria

1. When the `sign-tx` sign-wallet CLI command is invoked with coin type `hyc`, the Sign Wallet CLI shall invoke `SignTransactionUseCase` with the unsigned transaction file.
2. While signing, the ERC-20 Token Transfer system shall not make any network calls.
3. When the transaction is a Type 2 (EIP-1559) transaction, the Sign Wallet CLI shall sign it using the EIP-1559 signer with the correct chain ID.
4. When the transaction is a Type 0 (legacy) transaction, the Sign Wallet CLI shall sign it using the legacy signer.
5. When signing completes, the Sign Wallet CLI shall write the signed transaction JSON to the configured output directory.
6. If the private key is not found for the `from_address`, the Sign Wallet CLI shall return an error with the missing address.

---

### Requirement 6: Send Signed ERC-20 Transfer Transaction (Watch Wallet)

**Objective:** As a watch-wallet operator, I want to broadcast signed ERC-20 token transfer transactions,
so that token transfers are submitted to the Ethereum network.

#### Acceptance Criteria

1. When the `send-tx` watch-wallet CLI command is invoked with coin type `hyc`, the Watch Wallet CLI shall read the signed transaction file and invoke `SendTransactionUseCase`.
2. When `SendTransactionUseCase` succeeds, the ERC-20 Token Transfer system shall record the transaction hash in the database and mark the transaction as `sent`.
3. If the Ethereum node rejects the transaction (e.g., nonce too low, insufficient gas), the ERC-20 Token Transfer system shall return the node error without crashing.
4. When sending completes, the Watch Wallet CLI shall display the transaction hash to the operator.

---

### Requirement 7: Monitor ERC-20 Transfer Transaction (Watch Wallet)

**Objective:** As a watch-wallet operator, I want to monitor the status of sent ERC-20 token transfer
transactions, so that I know when they are confirmed or failed.

#### Acceptance Criteria

1. When the `monitor-tx` watch-wallet CLI command is invoked with coin type `hyc`, the Watch Wallet CLI shall invoke `MonitorTransactionUseCase`.
2. When a transaction reaches the required confirmation count, the ERC-20 Token Transfer system shall update the transaction status to `done` in the database.
3. When a monitored transaction is reverted on-chain, the ERC-20 Token Transfer system shall update the transaction status to `cancelled` and log the revert reason.
4. While monitoring, the ERC-20 Token Transfer system shall use exponential backoff retry logic to avoid hammering the Ethereum node.
5. If a transaction remains unconfirmed beyond the configured timeout, the ERC-20 Token Transfer system shall log a warning and continue monitoring.

---

### Requirement 8: Unit and Integration Tests

**Objective:** As a developer, I want comprehensive tests covering the ERC-20 token transfer flow,
so that regressions are caught automatically.

#### Acceptance Criteria

1. The ERC-20 Token Transfer system shall include unit tests for `CreateRawTransactionEIP1559()` verifying correct Type 2 transaction field population.
2. The ERC-20 Token Transfer system shall include unit tests for `SupportsEIP1559()` covering both true and false node responses.
3. The ERC-20 Token Transfer system shall include unit tests for the HYC domain registration verifying `IsCoinTypeCode("hyc")`, `IsERC20Token("hyc")`, and `IsETHGroup("hyc")`.
4. When `make go-test` is executed, all new tests shall pass with no failures.
5. When `make go-lint` is executed, no new lint errors shall be introduced.
