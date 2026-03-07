# Requirements Document

## Project Description (Input)

Implement Ethereum Multi-signature Transaction support using Safe (Gnosis Safe) smart contract accounts,
following the file-based offline signing workflow used by BTC and XRP multisig flows.
The implementation uses the Safe contract's `execTransaction` function with EIP-712 off-chain signatures
collected from each authorized signer via the existing wallet infrastructure (Watch, Keygen, Sign).

---

## Background

### Ethereum Multisig via Safe (Smart Account)

Ethereum does not have native protocol-level multisig. The standard approach in 2026 is to use
**Safe (formerly Gnosis Safe)** — a battle-tested smart contract wallet deployed as an on-chain account.

Key characteristics:

- **Safe contract**: An Ethereum smart contract deployed at a specific address that acts as the multisig wallet
- **Owners**: A set of Ethereum EOA addresses authorized to approve transactions
- **Threshold**: Minimum number of owner signatures required to execute a transaction
- **safeTxHash**: An EIP-712 typed data hash computed deterministically from the transaction parameters
- **execTransaction**: The Safe contract function that verifies collected signatures and executes the transaction

### Transaction Flow (4-Phase)

| Phase | Description | Wallet | Online/Offline |
|-------|-------------|--------|----------------|
| 1. Proposal | Watch wallet calls `getTransactionHash` on Safe contract to compute `safeTxHash`; writes transaction file | Watch | Online |
| 2. Off-chain Signing | Each owner signs `safeTxHash` using EIP-712 with their `ecdsa.PrivateKey`; signature appended to file | Keygen / Sign | **Offline** |
| 3. Aggregation | Watch wallet collects all signatures and sorts them by signer address (ascending) | Watch | Offline |
| 4. Execution | Watch wallet calls `execTransaction` on the Safe contract with aggregated signatures | Watch | Online |

### Signature Rules (EIP-712)

- Signatures must conform to `eth_signTypedData` (EIP-712)
- **Critical**: When concatenating signatures for `execTransaction`, signer addresses MUST be sorted in ascending order (e.g., `0x111... < 0x222...`). Failure to sort will cause a validation error on the contract
- Each signature is 65 bytes: `r (32) || s (32) || v (1)`
- The final `signatures` parameter is the concatenation of all sorted individual signatures

### Safe Contract Interface

The core function for execution:

```solidity
function execTransaction(
    address to,
    uint256 value,
    bytes calldata data,
    Enum.Operation operation,
    uint256 safeTxGas,
    uint256 baseGas,
    uint256 gasPrice,
    address gasToken,
    address payable refundReceiver,
    bytes calldata signatures
) external payable returns (bool success);
```

The hash computation function (used offline):

```solidity
function getTransactionHash(
    address to, uint256 value, bytes calldata data,
    Enum.Operation operation, uint256 safeTxGas, uint256 baseGas,
    uint256 gasPrice, address gasToken, address payable refundReceiver,
    uint256 _nonce
) public view returns (bytes32);
```

### Current Implementation State

The following components are **already implemented** and must not be re-implemented:

| Component | Location | Status |
|-----------|----------|--------|
| ETH single-sig EOA transaction create/sign/send | `usecase/watch/eth/`, `usecase/keygen/eth/` | ✅ Done |
| ETH EIP-1559 transaction creation | `infrastructure/api/eth/eth/transaction.go` | ✅ Done |
| `ETHTransactionFile` DTO (single-sig) | `application/dto/eth/transaction_file.go` | ✅ Done |
| ETH key derivation (BIP-44, `deriveChildPrivKey`) | `usecase/keygen/eth/sign_transaction.go` | ✅ Done |
| ETH wallet adapters (keygen, sign, watch) | `interface-adapters/wallet/eth/` | ✅ Done |
| ETH E2E P1 (single-sig) | `scripts/operation/eth/e2e/e2e-p1.sh` | ✅ Done |
| Foundry-based Solidity test infrastructure | `apps/eth-contracts/` | ✅ Done |
| `go-ethereum` library | `go.mod` | ✅ In use |

The following components are **missing** and must be implemented:

| Component | Gap | Required Action |
|-----------|-----|-----------------|
| Safe contract ABI binding | No Safe Go bindings in codebase | Generate with `abigen` from Safe ABI |
| Safe contract deployment script | No Foundry script for Safe | Add Foundry deploy script for local E2E |
| `ETHMultisigTransactionFile` DTO | No multisig file format for ETH | Define new DTO with Safe-specific fields |
| `SafeTxHashComputer` port interface | No port for `getTransactionHash` call | Define in `application/ports/api/eth/` |
| `SafeExecuter` port interface | No port for `execTransaction` call | Define in `application/ports/api/eth/` |
| Safe infrastructure implementation | No Safe client in `infrastructure/api/eth/` | Implement Safe contract client |
| Multisig create-transaction use case (Watch) | Existing create-tx is single-sig only | New use case: `CreateMultisigTransactionUseCase` |
| Multisig sign use case (Keygen/Sign) | `ETHSign.SignTx()` logs "no functionality" | Implement EIP-712 offline signing |
| Multisig send use case (Watch) | Existing send-tx is single-sig only | New use case: `SendMultisigTransactionUseCase` |
| CLI commands for multisig (Watch) | No multisig CLI commands for ETH | Add `create multisig`, `send multisig` commands |
| DI wiring for multisig use cases | No DI wiring for Safe use cases | Wire in `internal/di/container.go` |
| Safe contract management CLI | No CLI for Safe deployment/setup | Add `watch safe deploy`, `watch safe add-owner` commands |
| ETH E2E P2 (multisig) | Only P1 (single-sig) exists | Add P2 E2E test script |

---

## Requirements

### FR-1: Safe Contract Integration (Infrastructure)

**Goal**: Provide Go bindings and a client implementation for the Safe smart contract so that the
application layer can interact with Safe accounts via port interfaces.

**Acceptance Criteria**:

1. Safe contract ABI bindings are generated using `abigen` from the official Safe v1.4.1 ABI
2. The generated bindings are placed in `internal/infrastructure/contract/safe/` as an auto-generated file (with `DO NOT EDIT` header)
3. A `SafeClient` struct is implemented in `internal/infrastructure/api/eth/safe/` that wraps the generated bindings
4. `SafeClient` implements the `SafeTxHashComputer` port interface (calls `getTransactionHash` on-chain)
5. `SafeClient` implements the `SafeExecuter` port interface (calls `execTransaction` on-chain)
6. `SafeClient` dynamically obtains the chain ID from the connected Ethereum node (`client.ChainID(ctx)`)
7. Gas estimation for `execTransaction` is performed before submission
8. All Safe contract parameters not set by the operator default to zero/`Address(0)` (safeTxGas=0, baseGas=0, gasPrice=0, gasToken=Address(0), refundReceiver=Address(0), data=`[]`)
9. `SafeClient` accepts `operation` as a parameter; for simple ETH transfers, `operation=0` (Call)

### FR-2: Multisig Transaction File Format

**Goal**: Define a new DTO for the Safe multisig transaction file format that captures all data
needed by each signer and the final executor.

**Acceptance Criteria**:

1. A new `ETHMultisigTransactionFile` struct is defined in `internal/application/dto/eth/`
2. The file contains all fields required to reconstruct the `safeTxHash` offline:
   - `safe_address` — the Safe contract address
   - `to` — target address
   - `value` — ETH amount in Wei (decimal string)
   - `data` — hex-encoded call data (empty for simple ETH transfers)
   - `operation` — 0=Call, 1=DelegateCall
   - `safe_tx_gas`, `base_gas`, `gas_price`, `gas_token`, `refund_receiver` — Safe gas parameters
   - `nonce` — Safe's internal nonce (from `safe.nonce()`)
   - `chain_id` — for EIP-712 domain separator reconstruction
   - `safe_tx_hash` — the precomputed hash (for verification by signers)
   - `threshold` — number of signatures required
   - `signatures` — ordered list of collected `{signer_address, signature_hex}` entries
   - `tx_type` — `"unsigned"` or `"signed"` (becomes `"signed"` when `len(signatures) >= threshold`)
3. A `Validate()` method enforces all invariants
4. File naming convention uses a UUID instead of a DB-assigned txID: `{action}_multisig_{uuid}_{signedCount}.json`
   - A UUID is generated by the Watch wallet at proposal time and embedded in the file
   - This avoids requiring a database record for the multisig flow
5. The existing `ETHTransactionFile` is NOT modified — single-sig and multisig use separate file types

### FR-3: Port Interface Definitions for Safe Operations

**Goal**: Define focused port interfaces for Safe contract operations, following ISP.

**Acceptance Criteria**:

1. `SafeTxHashComputer` interface is defined in `internal/application/ports/api/eth/interface.go`:
   - Method: `GetSafeTxHash(ctx, safeAddr, to string, value *big.Int, data []byte, operation uint8, nonce *big.Int) (string, error)`
2. `SafeExecuter` interface is defined in `internal/application/ports/api/eth/interface.go`:
   - Method: `ExecuteSafeTransaction(ctx context.Context, params SafeExecParams) (string, error)`
   - `SafeExecParams` is a plain struct defined in `internal/application/ports/api/eth/interface.go` containing all execution parameters (safe address, to, value, data, operation, gas fields, nonce, aggregated signatures bytes). This decouples the port from the file DTO.
3. `SafeNonceReader` interface is defined in `internal/application/ports/api/eth/interface.go`:
   - Method: `GetSafeNonce(ctx context.Context, safeAddr string) (*big.Int, error)`
4. All interfaces use only primitive types or structs defined within the ports package — no `dto/eth/` types appear in port interface signatures
5. The `SafeClient` infrastructure implementation satisfies all three interfaces
6. The use case layer is responsible for mapping from `ETHMultisigTransactionFile` fields to `SafeExecParams` before calling `SafeExecuter`

### FR-4: Multisig Transaction Creation (Watch Wallet)

**Goal**: The Watch wallet creates a Safe transaction proposal file by computing the `safeTxHash`
on-chain and writing it to a file for offline signing.

**Acceptance Criteria**:

1. A new `CreateMultisigTransactionUseCase` interface is defined in `internal/application/usecase/watch/interfaces.go`
2. A new `CreateMultisigTransactionInput` struct includes:
   - `SafeAddress` — the Safe contract address
   - `ToAddress` — recipient
   - `Amount` — ETH amount in Ether (float64, converted to Wei internally)
   - `Threshold` — number of required signatures
   - `ActionType` — deposit/payment/transfer
3. The implementation (`usecase/watch/eth/create_multisig_transaction.go`):
   - Calls `SafeNonceReader.GetSafeNonce()` to fetch the current Safe nonce
   - Calls `SafeTxHashComputer.GetSafeTxHash()` to compute the hash on-chain
   - Writes an `ETHMultisigTransactionFile` with `tx_type: "unsigned"` and empty `signatures`
4. The generated file is passed to the Keygen wallet for the first signature
5. No database records are created for multisig transactions — the file is the sole state carrier; the UUID in the file name uniquely identifies the proposal across wallet hops
6. The command output prints the generated file path in format: `[fileName]: {path}`

### FR-5: Offline Multisig Signing (Keygen / Sign Wallet)

**Goal**: Each authorized signer computes the `safeTxHash` locally (from file data), signs it
with EIP-712 using their private key, and appends the signature to the file.

**Acceptance Criteria**:

1. A new `SignMultisigTransactionUseCase` is implemented for both Keygen and Sign wallets
2. The use case reads an `ETHMultisigTransactionFile` from disk
3. The signer **verifies** the `safeTxHash` by recomputing it locally from the file's EIP-712 fields:
   - Domain separator: `keccak256(domainSeparatorType || chainId || safeAddress)`
   - Struct hash: `keccak256(SAFE_TX_TYPEHASH || to || value || data || ...)`
   - Final hash: `keccak256("\x19\x01" || domainSeparator || structHash)`
   - If the recomputed hash does not match `safe_tx_hash` in the file, the signing is aborted with an error
4. Signing is performed offline using the derived `ecdsa.PrivateKey` — no Ethereum node connection
5. The signer's address and 65-byte hex signature are appended to the file's `signatures` array
6. When `len(signatures) >= threshold`, the file's `tx_type` is set to `"signed"`
7. The signed file is saved with an incremented counter in the file name (existing naming pattern)
8. The use case returns `IsComplete: true` when all required signatures are collected
9. The Keygen wallet's `SignTx` adapter method is wired to the new use case (used for signer 1)
10. The Sign wallet's `SignTx` adapter method (currently a no-op) is wired to the new use case (used for signers 2..N)
11. The Sign wallet is used **for multisig only** — single-sig ETH transactions continue to be signed exclusively by the Keygen wallet (existing behaviour is unchanged)

**Implementation Notes**:
- EIP-712 hashing must be implemented in pure Go using `go-ethereum/crypto` — no external Safe SDK dependency
- The private key derivation reuses the existing `deriveChildPrivKey` function from the Keygen signing use case
- The signer's address is derived from the public key (standard `crypto.PubkeyToAddress`)
- Both Keygen and Sign wallets share the same `SignMultisigTransactionUseCase` implementation (same package, different DI wiring)

### FR-6: Multisig Transaction Submission (Watch Wallet)

**Goal**: The Watch wallet reads a fully-signed multisig file, sorts signatures by signer address,
and submits the transaction to the Ethereum network via Safe's `execTransaction`.

**Acceptance Criteria**:

1. A new `SendMultisigTransactionUseCase` is implemented in `usecase/watch/eth/`
2. The use case reads an `ETHMultisigTransactionFile` and validates `tx_type == "signed"`
3. Signatures are sorted by signer address in ascending order before forming the `signatures` bytes
4. The use case calls `SafeExecuter.ExecuteSafeTransaction()` with the aggregated signatures
5. The use case polls for transaction receipt using the existing monitoring pattern (`GetTxReceipt`)
6. Transaction hash is logged; success/failure is returned
7. The submitted transaction hash is printed to stdout on success

### FR-7: Safe Account Inspection (Watch Wallet CLI)

**Goal**: The Watch wallet exposes a CLI command to inspect the current state of a Safe contract.
Safe deployment is handled exclusively via the Foundry script used in the E2E test (FR-10) —
the Watch wallet does not deploy Safe contracts, as deployment requires a funded deployer EOA
with a private key that the Watch wallet does not hold.

**Acceptance Criteria**:

1. `watch safe info` command displays current Safe state:
   - Accepts `--safe` (Safe contract address) flag
   - Prints owners list, threshold, current nonce, and ETH balance
   - Follows the existing Cobra CLI pattern in `internal/interface-adapters/cli/watch/`
2. The `ETHWatch` wallet adapter is updated to include the `SafeInfoUseCase`
3. Safe deployment for local E2E testing is handled by the Foundry script in FR-10 AC-3 (`DeploySafe.s.sol`) — not by a Watch wallet CLI command

### FR-8: Multisig Transaction CLI Commands (Watch Wallet)

**Goal**: Expose CLI commands for the complete Safe multisig transaction lifecycle on the Watch wallet.

**Acceptance Criteria**:

1. `watch create multisig` command creates a Safe transaction proposal file:
   - Accepts `--safe` (Safe address), `--to`, `--amount`, `--threshold` flags
   - Prints the generated file path
2. `watch send multisig` command submits a fully-signed multisig file:
   - Accepts `--file` flag pointing to the signed `ETHMultisigTransactionFile`
   - Validates the file before submission
   - Prints the transaction hash on success
3. The Sign wallet CLI wires the existing `sign tx` command (currently a no-op for ETH) to `SignMultisigTransactionUseCase` — it handles multisig only; single-sig ETH signing remains Keygen-only
4. All commands validate inputs at the CLI layer before passing to use cases

### FR-9: DI Container Wiring

**Goal**: Wire all new Safe multisig use cases and infrastructure components in the DI container.

**Acceptance Criteria**:

1. `SafeClient` is instantiated in the DI container and injected into:
   - `CreateMultisigTransactionUseCase` (Watch wallet)
   - `SendMultisigTransactionUseCase` (Watch wallet)
2. `SignMultisigTransactionUseCase` is wired in both Keygen and Sign wallet DI containers (shared implementation, separate factory calls)
3. `SafeInfoUseCase` is wired in the Watch wallet DI container (no `DeployUseCase` — deployment is handled by Foundry script)
4. All factory functions follow the existing `newXxx` naming convention in `internal/di/container.go`
5. No panics occur during DI initialization for ETH wallets

### FR-10: E2E Test — Safe Multisig Payment (P2)

**Goal**: A complete end-to-end test validates the ETH Safe multisig flow against a local Anvil node.

**Acceptance Criteria**:

1. A new E2E test script `scripts/operation/eth/e2e/e2e-p2.sh` covers the 2-of-2 Safe multisig flow:
   - Setup: Two signer accounts generated; Safe contract deployed with both as owners, threshold=2
   - Fund: ETH sent to the Safe contract address
   - Create: Watch wallet creates a multisig transaction file (`create multisig`)
   - Sign 1: Keygen wallet signs (signer 1), producing `signed_1` file
   - Sign 2: Sign wallet signs (signer 2), producing `signed_2` file with `tx_type: "signed"`
   - Send: Watch wallet submits via `execTransaction` (`send multisig`)
   - Verify: Transaction confirmed; recipient balance updated correctly
2. `make eth-e2e-p2` and `make eth-e2e-p2-ci` targets are added to the Makefile
3. A Foundry deployment script (`apps/eth-contracts/script/DeploySafe.s.sol`) deploys a Safe to the Anvil node
4. The test runs against Anvil (local) and passes reliably in CI

---

## Out of Scope

The following are explicitly **out of scope** for this specification:

| Item | Rationale |
|------|-----------|
| ERC-4337 / Account Abstraction (UserOperation / Bundler) | Future extension; direct `execTransaction` is simpler to start |
| Safe L2 adapter contracts | Out of scope; target local/mainnet Safe singleton factory |
| `SetRegularKey` / account recovery | Not applicable to ETH Safe |
| Web UI or HTTP API for multisig coordination | No HTTP/REST layer in this project |
| ERC-20 token multisig transfers | Phase 2 extension; initial scope is native ETH only |
| Safe modules (spending limits, social recovery) | Advanced features; not required for basic multisig |
| Meta-multisig (multisig to change Safe owners/threshold) | Out of scope; single-sig deployment of Safe is sufficient |
| Safe v1.3.x compatibility | Target Safe v1.4.1 only |

---

## Constraints

1. **Architecture**: Must follow Clean Architecture. Domain layer has zero infrastructure dependencies.
2. **Offline signing**: Keygen and Sign wallets must never make network calls during EIP-712 signing (FR-5).
3. **File transfer**: Transaction files move between wallets via physical transfer (USB). No networked key sharing.
4. **Library**: Use `go-ethereum` (`github.com/ethereum/go-ethereum`) for all cryptographic operations. Do not add a Safe SDK dependency.
5. **ABI binding**: Use `abigen` to generate Safe contract bindings from the official Safe v1.4.1 ABI. The generated file is auto-generated and must not be hand-edited.
6. **Signature ordering**: Signer address ordering (ascending) must be implemented at the Watch wallet layer (not at signing time), as required by the Safe contract.
7. **Chain ID**: Must be dynamically obtained from the node (`client.ChainID(ctx)`) — never hardcoded.
8. **No gRPC**: Do not introduce any gRPC dependencies. All Safe contract interactions use the existing `go-ethereum` JSON-RPC client.
9. **Backward compatibility**: All existing single-sig ETH flows (P1) must continue to function without modification.
