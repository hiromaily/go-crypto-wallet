# Research & Design Decisions

## Summary

- **Feature**: `eth-transaction-flow-alignment`
- **Discovery Scope**: Extension (existing system with significant gaps)
- **Key Findings**:
  - EIP-1559 infrastructure already exists in codebase but uses static fee defaults; needs dynamic estimation via `SuggestGasTipCap()`
  - Sign wallet is entirely non-functional (all stubs); BTC reference provides clear ISP-compliant port interface pattern to follow
  - RLP encoding in `ethtx.go` does not correctly handle EIP-2718 typed transactions; must switch to `MarshalBinary`/`UnmarshalBinary`
  - Two use cases (`sign/eth/sign_transaction.go`, `keygen/eth/import_private_key.go`) import from infrastructure layer — Clean Architecture violations
  - Anvil fully supports standard `eth_*` JSON-RPC including `eth_maxPriorityFeePerGas`; no custom adapter needed

## Research Log

### ETH Sign Wallet Current State

- **Context**: Requirement 1 — assess current sign wallet functionality
- **Sources**: `internal/interface-adapters/wallet/eth/sign.go`
- **Findings**:
  - All methods are stubs returning `nil`: `GenerateSeed`, `StoreSeed`, `GenerateAuthKey`, `ImportPrivKey`, `ExportFullPubkey`, `SignTx`
  - Log messages are misleading (reference `CreateMultisigAddress()` in ETH)
  - Only `Done()` method is functional (closes DB and client connections)
  - DI container returns "not implemented yet" for ETH sign wallet
- **Implications**: Full implementation required — no existing logic to extend

### EIP-1559 Implementation Status

- **Context**: Requirement 2 — assess existing EIP-1559 support
- **Sources**: `internal/infrastructure/api/eth/eth/transaction.go` lines 319-493
- **Findings**:
  - `CreateRawTransactionEIP1559()` fully implemented with `types.DynamicFeeTx`
  - `SupportsEIP1559()` detection works via `baseFeePerGas` in block header
  - Fee calculation: `maxFeePerGas = baseFee * 2 + maxPriorityFeePerGas` (correct formula)
  - `maxPriorityFeePerGas` uses static config default (2 Gwei) instead of dynamic `eth_maxPriorityFeePerGas` RPC
  - `NewLondonSigner(chainID)` used at line 221; should be `LatestSignerForChainID` for forward compatibility
  - ERC-20 token transfers only support legacy transactions
- **Implications**: Infrastructure exists but needs wiring to use case layer and dynamic fee estimation

### Transaction Serialization for File Exchange

- **Context**: Requirement 6 — file-based transaction exchange
- **Sources**: `internal/infrastructure/api/eth/ethtx/ethtx.go`
- **Findings**:
  - Current `EncodeTx`/`DecodeTx` use `rlp.EncodeToBytes`/`rlp.Decode`
  - RLP encoding does not produce correct EIP-2718 typed transaction envelope (`type || payload`)
  - `types.Transaction.MarshalBinary()` handles all transaction types correctly (Legacy, EIP-2930, EIP-1559, EIP-4844)
  - JSON marshaling also available via `tx.MarshalJSON()` for human-readable format
- **Implications**: Must migrate to `MarshalBinary`/`UnmarshalBinary` before enabling EIP-1559 file exchange

### BTC Reference: ISP Port Interfaces

- **Context**: Requirement 5 — Clean Architecture compliance
- **Sources**: `internal/application/ports/api/btc/interfaces_small.go`
- **Findings**:
  - BTC defines small, focused interfaces: `ChainConfigProvider`, `AmountConverter`, `UTXOProvider`, `RawTransactionCreator`, `AddressOperator`, `BalanceChecker`, `TransactionSender`, `TransactionMonitor`, `PSBTCreator`, `PSBTSigner`, `PSBTFinalizer`
  - Use cases define local interfaces embedding only needed ports (e.g., `signTxBTCClient`)
  - Composed interfaces combine small interfaces for specific use cases
  - Go's implicit interface satisfaction bridges full implementation to small ports
- **Implications**: ETH should define equivalent small interfaces in `internal/application/ports/api/eth/interfaces_small.go`

### BTC Reference: File-Based Exchange Flow

- **Context**: Requirement 6 — transaction file format
- **Sources**: `internal/infrastructure/storage/file/transaction/`
- **Findings**:
  - File naming: `{actionType}_{txID}_{txType}_{signedCount}_{timestamp}.{ext}`
  - BTC uses PSBT (base64) with `.psbt` extension
  - `ValidateFilePath()` parses filename components for metadata
  - `WritePSBTFile`/`ReadPSBTFile` handle serialization
  - Flow: Watch creates unsigned → Keygen signs (count 1) → Sign signs (count 2, fully signed) → Watch broadcasts
- **Implications**: ETH should use JSON with `.json` extension for human readability; same naming convention and flow pattern

### Anvil/Foundry Compatibility

- **Context**: Requirement 3 — development environment
- **Sources**: Foundry docs, GitHub issues
- **Findings**:
  - Anvil supports all standard `eth_*` JSON-RPC methods
  - `eth_maxPriorityFeePerGas` supported (returns static value, acceptable for dev)
  - No `personal_importRawKey` needed; uses pre-funded accounts or `anvil_impersonateAccount`
  - Default chain ID: 31337; configurable via `--chain-id`
  - 10 pre-funded accounts with 10,000 ETH each
  - `ethclient.Dial("http://localhost:8545")` works without changes
- **Implications**: No special adapter needed; standard ethclient works with Anvil

### Clean Architecture Violations

- **Context**: Requirement 5 — identify violations
- **Sources**: Codebase grep for infrastructure imports in use case layer
- **Findings**:
  - `internal/application/usecase/sign/eth/sign_transaction.go`: imports `internal/infrastructure/api/eth` and `internal/infrastructure/api/eth/eth`
  - `internal/application/usecase/keygen/eth/import_private_key.go`: same violations
  - Root cause: `ethereum.Ethereumer` interface from infrastructure, `apiethimpl.Password` constant
  - The `Ethereumer` interface in ports already exists but use cases reference the deprecated one
- **Implications**: Fix imports to use port interfaces; move password constant to config

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| Mirror BTC port pattern | Small ISP interfaces in `ports/api/eth/` | Proven in codebase, consistent | Requires interface decomposition work | Selected |
| Keep monolithic Ethereumer | Use existing large interface | Less refactoring | Perpetuates ISP violation | Rejected |
| Adapter pattern | Wrap infrastructure with adapters | Good isolation | Over-engineering for this case | Rejected |

## Design Decisions

### Decision: Transaction File Format

- **Context**: ETH needs file-based exchange like BTC's PSBT
- **Alternatives Considered**:
  1. Binary (RLP hex) — compact but not human-readable
  2. JSON — human-readable, debuggable, larger
  3. Protobuf — efficient but requires schema management
- **Selected Approach**: JSON with `.json` extension
- **Rationale**: Aligns with Requirement 6.4 (human readability); JSON debugging is valuable for offline signing workflow; size overhead is negligible for single transactions
- **Trade-offs**: Slightly larger files vs much better debuggability
- **Follow-up**: Define JSON schema for transaction file

### Decision: Transaction Encoding (Internal)

- **Context**: Current RLP encoding breaks with EIP-1559 typed transactions
- **Alternatives Considered**:
  1. Keep RLP (`rlp.EncodeToBytes`) — current approach, breaks with typed txs
  2. Use `MarshalBinary`/`UnmarshalBinary` — handles all tx types
  3. Use `MarshalJSON`/`UnmarshalJSON` — human-readable but less efficient
- **Selected Approach**: `MarshalBinary`/`UnmarshalBinary` for internal hex encoding; JSON wrapper for file exchange
- **Rationale**: `MarshalBinary` is the canonical go-ethereum method for transaction serialization and correctly handles the EIP-2718 envelope format
- **Trade-offs**: Requires updating `ethtx.go` but gains correctness for all transaction types

### Decision: Signer Selection

- **Context**: Current code uses `types.NewLondonSigner(chainID)`
- **Selected Approach**: Switch to `types.LatestSignerForChainID(chainID)`
- **Rationale**: Forward-compatible with EIP-4844 (blob txs) and EIP-7702 (set-code txs); subsumes all signers that `NewLondonSigner` supports
- **Follow-up**: Update `transaction.go:221`

### Decision: Fee Estimation Strategy

- **Context**: Static 2 Gwei default vs dynamic RPC estimation
- **Selected Approach**: Dynamic `SuggestGasTipCap()` with config fallback
- **Rationale**: Requirement 2.3 specifies `eth_maxPriorityFeePerGas` RPC; static defaults may over/under-pay on production
- **Follow-up**: Verify `SuggestGasTipCap` availability on target nodes

### Decision: Single-Signature Flow (No Multisig)

- **Context**: Multisig excluded from this spec
- **Selected Approach**: ETH uses Watch → Sign → Watch flow (no Keygen intermediate signing)
- **Rationale**: Without multisig, only one signature is needed; Keygen wallet generates keys and exports pubkeys but does not sign transactions
- **Trade-offs**: Simpler flow but different from BTC's three-step signing
- **Follow-up**: Multisig spec will add Keygen signing step later

## Risks & Mitigations

- **Risk**: RLP → MarshalBinary migration may break existing stored transaction hex — **Mitigation**: Version the encoding; add format detection in `DecodeTx`
- **Risk**: `SuggestGasTipCap` returns unreasonable values on some networks — **Mitigation**: Cap at configured `maxPriorityFeePerGas` ceiling; use config value as fallback
- **Risk**: Sign wallet implementation complexity (HD key derivation + offline signing) — **Mitigation**: Reuse existing `key/` infrastructure; BIP-44 derivation already implemented in `hd_wallet.go`
- **Risk**: Interface decomposition may cause regressions in DI wiring — **Mitigation**: Go's implicit interface satisfaction means existing implementations automatically satisfy new small interfaces

## References

- [go-ethereum types package](https://pkg.go.dev/github.com/ethereum/go-ethereum/core/types) — Transaction types and signers
- [go-ethereum ethclient](https://pkg.go.dev/github.com/ethereum/go-ethereum/ethclient) — `SuggestGasTipCap`, `HeaderByNumber`
- [Anvil Reference](https://getfoundry.sh/anvil/reference/) — JSON-RPC methods and custom methods
- [EIP-1559](https://eips.ethereum.org/EIPS/eip-1559) — Fee market change specification
- [EIP-2718](https://eips.ethereum.org/EIPS/eip-2718) — Typed Transaction Envelope
- [BIP-44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki) — Multi-Account Hierarchy for Deterministic Wallets
