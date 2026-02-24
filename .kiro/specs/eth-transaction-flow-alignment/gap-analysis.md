# Gap Analysis: eth-transaction-flow-alignment

## 1. Requirement-to-Asset Map

### Requirement 1: Offline Transaction Signing (Sign Wallet)

| Need | Existing Asset | Gap |
|------|---------------|-----|
| Functional Sign wallet | `internal/interface-adapters/wallet/eth/sign.go` — all 6 methods are stubs (log "no functionality", return nil) | **Missing** |
| DI wiring for Sign | `internal/di/container.go:288-289` — returns `"not implemented yet"` for ETH | **Missing** |
| Sign use case | `internal/application/usecase/sign/eth/sign_transaction.go` — exists but uses deprecated `ethereum.Ethereumer` from infrastructure + hardcoded `Password` | **Constraint**: architecture violation, needs rewrite |
| BIP-44 child key derivation | `internal/infrastructure/wallet/key/strategy/eth_strategy.go` — derives ETH addresses from HD keys | Partial: generates keys but Sign wallet doesn't use `accountXpriv` for derivation |
| Offline signing (no RPC) | Current `SignOnRawTransaction()` calls `e.GetPrivKey()` which reads from Geth keystore on-disk | **Missing**: signing requires local keystore, which is Geth-coupled |
| File-based tx exchange | `domainEthereum.RawTx` serialized via `serializer.EncodeToString()` | Exists but uses serialized binary format, not structured JSON |

### Requirement 2: EIP-1559 Transaction Support

| Need | Existing Asset | Gap |
|------|---------------|-----|
| EIP-1559 tx creation | `internal/infrastructure/api/eth/eth/transaction.go:CreateRawTransactionEIP1559()` — fully implemented | **Exists** (not wired to use case) |
| EIP-1559 detection | `SupportsEIP1559()` — checks block `baseFeePerGas` + Anvil auto-detection | **Exists** |
| Use case integration | `internal/application/usecase/watch/eth/create_transaction.go` — calls `CreateRawTransaction()` (legacy only) | **Missing**: use case doesn't call `CreateRawTransactionEIP1559()` |
| Fee config | `ethereum.conf.MaxPriorityFeePerGas` — exists in infrastructure | **Exists** (needs config YAML wiring) |
| Transaction type in file | `domainEthereum.RawTx` — no tx type field | **Missing**: file format doesn't distinguish legacy vs EIP-1559 |

### Requirement 3: Foundry/Anvil Development Environment

| Need | Existing Asset | Gap |
|------|---------------|-----|
| Anvil endpoint support | ETH client connects via `ethclient.Dial()` to any JSON-RPC endpoint | **Exists** (Anvil is compatible) |
| Anvil-compatible key import | `internal/application/usecase/keygen/eth/import_private_key.go` — uses `keystore.ImportECDSA()` (not `personal_importRawKey`) | **Exists** (already Anvil-compatible) |
| Network config update | Config files reference `goerli` (deprecated since 2023) | **Missing**: needs Sepolia/Holesky/local |
| Anvil startup scripts | None | **Missing**: no dev scripts/docs for Anvil |
| Client type detection | `ClientVersionAnvil` constant exists in infrastructure | **Exists** |

### Requirement 4: Robust Key Generation

| Need | Existing Asset | Gap |
|------|---------------|-----|
| BIP-39/44 HD wallet | `internal/infrastructure/wallet/key/hd_wallet.go` — shared HD key generation with BTC | **Exists** |
| ETH key strategy | `internal/infrastructure/wallet/key/strategy/eth_strategy.go` — Keccak-256 address derivation | **Exists** |
| `accountXpriv` storage | Shared keygen use case stores `accountXpriv` in DB | **Exists** (shared with BTC) |
| Full pubkey export/import | BTC has `ExportFullPubkey`/`ImportFullPubKey`; ETH Sign wallet stubs return nil | **Missing**: ETH Sign wallet needs implementation |
| Keygen adapter | `internal/interface-adapters/wallet/eth/keygen.go` | Exists but needs review for completeness |
| Password management | `apiethimpl.Password = "password"` hardcoded constant | **Constraint**: security issue |

### Requirement 5: Multisig Support (Safe Smart Contract Wallet)

| Need | Existing Asset | Gap |
|------|---------------|-----|
| Safe contract interaction | None | **Missing**: entirely new capability |
| EIP-712 typed data signing | None | **Missing**: no typed data signing in codebase |
| Safe ABI/contract types | None | **Missing**: Safe contract bindings needed |
| Threshold/owner config | None | **Missing**: config structure doesn't support Safe params |
| Partially-signed tx storage | `ETHDetailTx` domain type — has `SignedHexTx` (single field) | **Constraint**: needs extension for multiple signatures |
| Multisig tx file format | None | **Missing**: file format doesn't support signature aggregation |

### Requirement 6: Clean Architecture Compliance

| Need | Existing Asset | Gap |
|------|---------------|-----|
| Port-only use case imports | `internal/application/ports/api/eth/interface.go` — port interfaces exist | **Constraint**: 4 use cases import from infrastructure |
| ISP-compliant interfaces | BTC has `interfaces_small.go` with 15+ small interfaces | **Missing**: ETH has monolithic `Ethereumer` (70+ methods) |
| Remove deprecated interfaces | `internal/infrastructure/api/eth/api-interface.go` — marked `Deprecated:` | **Missing**: still imported by use cases |
| Configurable password | `apiethimpl.Password = "password"` hardcoded | **Missing**: no config injection |
| Port-level DTOs | `apieth.TxCreateParams` exists | Partial: some DTOs exist, others leak infrastructure types |

### Requirement 7: Transaction File Format and Flow Alignment

| Need | Existing Asset | Gap |
|------|---------------|-----|
| Structured unsigned tx file | Current: binary-serialized `RawTx` via `serializer.EncodeToString()`, line-separated | **Constraint**: not JSON, not human-readable |
| Partially-signed file | Current: signed file is `UUID,signedTxHex` CSV per line | **Constraint**: no metadata, no multi-sig tracking |
| Watch → Keygen → Sign flow | Current: Watch → Keygen signs (single sig, done) → Watch broadcasts. Sign wallet unused. | **Missing**: Sign wallet step entirely absent |
| `eth_sendRawTransaction` broadcast | `SendSignedRawTransaction()` exists | **Exists** |
| File repo | `file.TransactionFileRepositorier` — shared file I/O interface | **Exists** |

### Requirement 8: Configuration and Network Modernization

| Need | Existing Asset | Gap |
|------|---------------|-----|
| PostgreSQL support | BTC configs support postgres; ETH configs have mysql+sqlite only | **Missing**: no postgres in ETH config |
| Modern network types | Config: `network_type: "goerli"` with options `mainnet, goerli, rinkeby, ropsten` | **Missing**: all testnets deprecated except mainnet |
| EIP-1559 fee params in config | `MaxPriorityFeePerGas` exists in infrastructure struct but not in YAML config | **Missing**: not in config file |
| Chain ID config | Uses `netID` from `net_version` RPC at runtime | Partial: runtime detection exists but no explicit config |
| Network validation | None | **Missing**: no validation of deprecated network names |

### Requirement 9: Transaction Monitoring and Status Tracking

| Need | Existing Asset | Gap |
|------|---------------|-----|
| Confirmation tracking | `updateStatusTxTypeSent()` — checks confirmations and updates to `TxTypeDone` | **Exists** |
| `TxTypeDone` update | `updateStatusTxTypeDone()` — commented out in `monitor_transaction.go:44-51` | **Missing**: incomplete implementation |
| Failed/reverted detection | None | **Missing**: no revert reason checking |
| `is_allocated` update | `send_transaction.go:109` — `TODO: update is_allocated` | **Missing**: explicitly marked TODO |
| Retry with backoff | None | **Missing**: hard failure on node unreachable |
| Architecture compliance | Uses `ethereum.Ethereumer` (infrastructure import) | **Constraint**: violates Clean Architecture |

---

## 2. Implementation Approach Options

### Option A: Extend Existing Components

**Approach**: Modify existing ETH files in-place to add missing functionality.

**What to extend**:
- `internal/interface-adapters/wallet/eth/sign.go` — replace stubs with real implementations
- `internal/application/usecase/sign/eth/sign_transaction.go` — fix imports, add offline signing
- `internal/application/usecase/watch/eth/create_transaction.go` — switch to EIP-1559
- `internal/application/ports/api/eth/interface.go` — decompose `Ethereumer` into small interfaces
- `internal/di/container.go` — wire ETH Sign wallet
- Config YAML files — add postgres, modern networks, fee params

**Trade-offs**:
- (+) Minimal new files, leverages existing structure
- (+) Changes are traceable in existing file history
- (-) Risk of bloating existing files (especially `interface.go`)
- (-) Multisig (Safe) is fundamentally new — can't be "extended" from existing code

### Option B: Create New Components

**Approach**: Build new packages for multisig, offline signing, and ISP interfaces alongside existing code.

**New components**:
- `internal/application/ports/api/eth/interfaces_small.go` — ISP interfaces (following BTC pattern)
- `internal/infrastructure/api/eth/safe/` — Safe contract interaction package
- `internal/application/usecase/keygen/eth/create_multisig_address.go` — Safe address setup
- `internal/domain/ethereum/eth_multisig_tx.go` — Multisig transaction domain type
- `internal/infrastructure/api/eth/eth/eip712.go` — EIP-712 typed data signing

**Trade-offs**:
- (+) Clean separation for new capabilities (Safe, EIP-712)
- (+) Easier to test multisig in isolation
- (+) BTC provides clear template for ISP pattern
- (-) More files to navigate
- (-) Requires careful interface design for Safe interactions

### Option C: Hybrid Approach (Recommended)

**Approach**: Extend existing components for incremental improvements (EIP-1559, architecture fixes, config), create new components for fundamentally new capabilities (Safe multisig, offline signing).

**Phase 1 — Foundation (extend existing)**:
- Fix all architecture violations (use case imports)
- Decompose `Ethereumer` into ISP interfaces
- Switch `create_transaction.go` to EIP-1559
- Update configs (postgres, networks, fees)
- Replace hardcoded password with config injection

**Phase 2 — Sign Wallet (extend + new)**:
- Implement `ETHSign` adapter methods (replace stubs)
- Wire DI container for ETH Sign
- Create offline signing logic (no Geth keystore dependency)
- Design new JSON transaction file format

**Phase 3 — Multisig (new components)**:
- Create Safe contract interaction package
- Add EIP-712 typed data signing
- Extend domain types for multi-signature tracking
- Create multisig use cases
- Add Safe config parameters

**Phase 4 — Monitoring (extend existing)**:
- Uncomment and implement `updateStatusTxTypeDone()`
- Add `is_allocated` update after send
- Add failed/reverted transaction detection
- Add retry with backoff

**Trade-offs**:
- (+) Incremental delivery — each phase independently valuable
- (+) Architecture fixes in Phase 1 enable cleaner Phase 2-3
- (+) Risk mitigation through phased rollout
- (-) More complex planning required
- (-) Phase 3 (multisig) is the highest complexity and risk

---

## 3. Effort and Risk Assessment

| Requirement | Effort | Risk | Justification |
|---|---|---|---|
| R1: Sign Wallet | **L** (1-2 weeks) | **High** | Requires offline signing without Geth keystore, new key derivation flow, fundamental redesign of signing approach |
| R2: EIP-1559 | **S** (1-3 days) | **Low** | Infrastructure code already exists, just needs use case wiring and fallback logic |
| R3: Foundry/Anvil | **S** (1-3 days) | **Low** | Anvil is JSON-RPC compatible, key import already uses keystore (not `personal_importRawKey`) |
| R4: Key Generation | **M** (3-7 days) | **Medium** | HD wallet exists but Sign wallet integration and pubkey export need new flows |
| R5: Multisig (Safe) | **XL** (2+ weeks) | **High** | Entirely new capability — Safe contract ABI, EIP-712, signature aggregation, new domain types |
| R6: Clean Architecture | **M** (3-7 days) | **Medium** | Well-understood refactoring but touches many files, risk of breaking existing flows |
| R7: Tx File Format | **M** (3-7 days) | **Medium** | New JSON format design, migration from binary serialization, affects all wallet stages |
| R8: Config Modernization | **S** (1-3 days) | **Low** | Straightforward config additions following BTC pattern |
| R9: Monitoring | **S** (1-3 days) | **Low** | Mostly uncommenting and completing existing code, plus minor additions |

**Overall Effort**: **XL** (6-8 weeks total)
**Overall Risk**: **High** — primarily driven by R1 (offline signing redesign) and R5 (Safe multisig)

---

## 4. Research Items for Design Phase

### Critical Research

1. **Offline ECDSA Signing Without Geth Keystore**
   - Current `SignOnRawTransaction()` depends on `e.GetPrivKey()` which reads from Geth keystore files
   - Sign wallet must operate without any Geth node → need to sign using raw `*ecdsa.PrivateKey` from `crypto.ToECDSA()`
   - Research: Can `types.SignTx()` work with a raw private key derived from `accountXpriv` HD derivation without keystore?
   - Answer: Yes — `types.SignTx(tx, signer, privateKey)` accepts `*ecdsa.PrivateKey` directly

2. **Safe Contract Go Bindings**
   - Research: Are there existing Go bindings for Safe (Gnosis Safe) contracts?
   - Options: `abigen` from go-ethereum to generate from ABI, or manual `ethclient.CallContract()` calls
   - Need: Safe v1.4.1 ABI (latest audited version), `GnosisSafe.sol` and `GnosisSafeProxy.sol`

3. **EIP-712 Typed Data Signing in Go**
   - Research: How to implement EIP-712 structured data hashing in Go for Safe `execTransaction` signing?
   - `go-ethereum` has `signer.TypedData` and `core.SignTypedData` — evaluate for Safe compatibility
   - Alternative: Manual implementation of domain separator + message hash

4. **HD Wallet Library Stability**
   - `tyler-smith/go-bip39` (used transitively) was deleted by author
   - Research: Is current `go.mod` dependency resolution stable? Does `go-ethereum-hdwallet` still build?
   - Alternative: Migrate to `kslamph/bip39-hdwallet` if needed

### Non-Critical Research

5. **Safe Deployment on Anvil for Testing**
   - How to deploy Safe contracts to local Anvil for integration testing
   - Safe has a deployment script — verify Anvil compatibility

6. **ERC-4337 vs Safe-Only Multisig**
   - Current requirements specify Safe — evaluate if ERC-4337 account abstraction should be a future extension
   - Decision: Safe is sufficient for initial implementation; ERC-4337 is additive

---

## 5. Key Observations and Constraints

### Architecture Violations (Files to Fix)

| File | Violation | Fix |
|------|-----------|-----|
| `usecase/keygen/eth/sign_transaction.go` | Imports `ethereum` (infra) + `apiethimpl` (infra) | Use port interfaces |
| `usecase/sign/eth/sign_transaction.go` | Same | Use port interfaces |
| `usecase/watch/eth/send_transaction.go` | Imports `ethereum` (infra) | Use port interfaces |
| `usecase/watch/eth/monitor_transaction.go` | Imports `ethereum` (infra) | Use port interfaces |

### Security Constraints

- `Password = "password"` hardcoded in `internal/infrastructure/api/eth/eth/types.go` — used in keygen and sign use cases
- Sign wallet design must ensure private keys never leave the offline environment
- Safe multisig signature collection must validate signer addresses against configured owner list

### BTC Reference Patterns to Follow

- ISP interfaces: `interfaces_small.go` with granular interfaces + composed interfaces for common use cases
- Sign wallet adapter: `btcwallet.NewBTCSign()` with proper DI wiring
- DI container: `newBTCSigner()` pattern at `container.go:313`
- Transaction file: PSBT-based structured format (ETH equivalent: JSON)
