# Requirements Document

## Introduction

This specification defines the requirements for implementing the **HYC ERC-20 token** smart contract and integrating it with the existing go-crypto-wallet CLI infrastructure.

The scope covers three areas:
1. Solidity smart contract implementation (HYC token)
2. Foundry-based compilation and deployment to a local Ethereum node
3. CLI support for ERC-20 token transfers using EIP-1559 transactions

Target network: Ethereum (local development via `anvil` / `geth` Docker nodes at `http://localhost:8545`).

---

## Requirements

### Requirement 1: HYC ERC-20 Smart Contract

**Objective:** As a developer, I want an ERC-20 compliant Solidity contract named HYC, so that the token can be deployed to Ethereum nodes and used by the wallet CLI.

#### Acceptance Criteria

1. The HYC contract shall implement the ERC-20 standard by inheriting from OpenZeppelin's `ERC20` base contract.
2. The HYC contract shall set the token name to `"hiromaily Coin"` and the symbol to `"HYC"`.
3. The HYC contract shall use 18 decimal places (OpenZeppelin default).
4. When the HYC contract is deployed, the HYC contract shall mint the entire initial supply to the deployer's address.
5. The HYC contract shall accept `initialSupply` as a constructor parameter to allow configurable supply at deploy time.
6. The HYC contract shall use Solidity version `^0.8.34`.

---

### Requirement 2: Foundry Build Environment

**Objective:** As a developer, I want a Foundry-based build environment configured for the HYC project, so that contracts can be compiled, tested, and deployed consistently.

#### Acceptance Criteria

1. The HYC project shall include a `foundry.toml` configuration file at the project root.
2. The foundry configuration shall set `src = "contracts"`, `out = "out"`, and `libs = ["node_modules"]`.
3. The foundry configuration shall set `solc_version = "0.8.34"`, `optimizer = true`, and `optimizer_runs = 200`.
4. When `forge build` is executed, the ERC-20 compilation system shall produce compiled contract artifacts without errors.
5. The project shall use `bun` as the primary package manager; if bun is unavailable, `pnpm` shall be used as fallback.
6. When dependencies are installed, the package manager shall install `@openzeppelin/contracts` as a runtime dependency.

---

### Requirement 3: Contract Deployment to Local Ethereum Node

**Objective:** As a developer, I want a Foundry deployment script for the HYC token, so that the contract can be deployed to local Ethereum development nodes.

#### Acceptance Criteria

1. The project shall include a Foundry deployment script at `script/DeployHYC.s.sol`.
2. The deployment script shall read the deployer's private key from the `PRIVATE_KEY` environment variable.
3. When the deployment script is executed against a local Ethereum node, the deployment system shall broadcast the contract deployment transaction.
4. When deployment succeeds, the deployment system shall output the deployed contract address, transaction hash, and gas usage.
5. The deployment shall use an initial supply of `1_000_000 ether` (i.e., `1_000_000 × 10^18` token units).
6. The deployment script shall be compatible with both `anvil` and `geth` nodes reachable at `http://localhost:8545`.

---

### Requirement 4: ERC-20 Transfer Encoding

> **Out of scope for this spec.** This requirement is fully satisfied by the existing Go infrastructure at `internal/infrastructure/api/eth/erc20/erc20.go` (`createTransferData()` method), which already encodes `transfer(address,uint256)` calldata using method selector `0xa9059cbb` with correct 32-byte ABI encoding for address and amount.

---

### Requirement 5: EIP-1559 Transaction Construction for ERC-20 Transfers

> **Out of scope for this spec.** The existing Go ERC-20 infrastructure (`erc20.go`) uses legacy (type 0) transactions. EIP-1559 support for ERC-20 transfers is a separate enhancement to the Go wallet and is not part of this Foundry/Solidity spec.

---

### Requirement 6: CLI Command for HYC Token Transfer

> **Out of scope for this spec.** The existing Go wallet CLI (`internal/interface-adapters/cli/`) and use-case layer already support ERC-20 token transfers. Registering the HYC token in the wallet configuration (adding `HYC` as an `ERC20Token` constant in `internal/domain/coin/types.go` and adding a config entry) is a separate Go-side task.

---

### Requirement 7: Development Tooling — Formatter and Linter

**Objective:** As a developer, I want Solidity and TypeScript/JavaScript code quality tools configured, so that code style is enforced consistently across the project.

#### Acceptance Criteria

1. The project shall include `solhint` as a development dependency for Solidity linting.
2. The project shall include a `.solhint.json` configuration file extending `solhint:recommended`.
3. When `solhint contracts/**/*.sol` is executed, the Solidity linter shall report zero errors on the HYC contract source.
4. The project shall include `dprint` as a development dependency for TypeScript/JavaScript formatting.
5. The project shall include a `dprint.json` configuration file specifying the TypeScript plugin (latest available version).
6. When `dprint fmt` is executed, the formatter shall format all TypeScript/JavaScript files without errors.
7. The project shall install development dependencies (`solhint`, `dprint`) via `bun add -d` (or `pnpm add -D` as fallback).
