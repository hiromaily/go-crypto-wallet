<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/eth/multisig.tpl.md · Run `make docs` to regenerate.
-->

# Ethereum Safe Multisig Reference

This document covers the Safe (Gnosis Safe v1.4.1) multisig implementation for Ethereum in go-crypto-wallet. It explains the architecture, EIP-712 signing flow, file exchange format, CLI commands, and E2E Pattern 3 workflow.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Safe Contract Background](#safe-contract-background)
4. [EIP-712 Signing Flow](#eip-712-signing-flow)
5. [Transaction File Format](#transaction-file-format)
6. [CLI Commands](#cli-commands)
7. [Wallet Roles in Multisig](#wallet-roles-in-multisig)
8. [E2E Pattern 3: 2-of-2 Safe Payment](#e2e-pattern-3-2-of-2-safe-payment)
9. [Configuration](#configuration)
10. [Key Implementation Files](#key-implementation-files)
11. [Official References](#official-references)

---

## Overview

Ethereum EOA (Externally Owned Account) does not support multisig at the protocol level. This system implements multisig via **Gnosis Safe v1.4.1** — an audited smart contract that enforces an m-of-n threshold before executing any transaction.

The implementation uses a **file-based, offline-signing workflow** consistent with the existing EOA single-sig flow:

1. Watch Wallet proposes a transaction → writes JSON file
2. Each signer (Keygen or Sign wallet) reads the file, verifies the hash offline, appends a signature → writes a new file
3. Watch Wallet reads the fully-signed file → calls `execTransaction` on the Safe contract

No private keys touch the online Watch Wallet at any point.

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│  Watch Wallet (Online)                                            │
│  ├─ watch create multisig    → writes unsigned JSON file         │
│  ├─ watch safe info          → reads Safe contract state         │
│  └─ watch send multisig send-eth → submits execTransaction       │
└──────────────────────────────────────────────────────────────────┘
                            │ JSON files
                   ┌────────┴────────┐
                   ↓                 ↓
┌─────────────────────┐   ┌─────────────────────┐
│ Keygen Wallet       │   │ Sign Wallet         │
│ (Offline/Air-gapped)│   │ (Offline/Air-gapped)│
│ keygen sign         │   │ sign sign           │
│ signature           │   │ signature           │
│ --signer-address    │   │ --signer-address    │
└─────────────────────┘   └─────────────────────┘
         │                          │
         └──────────┬───────────────┘
                    ↓
          Signed JSON file (signature count reaches threshold)
```

### Use Case Layer

| Use Case | Package | Description |
|----------|---------|-------------|
| `CreateETHMultisigTransactionUseCase` | `usecase/watch/eth` | Proposes an unsigned Safe transaction |
| `SignMultisigTransactionUseCase` | `usecase/keygen/eth` | Verifies hash offline, appends signature |
| `SendETHMultisigTransactionUseCase` | `usecase/watch/eth` | Calls `execTransaction` on the Safe |
| `ETHSafeInfoUseCase` | `usecase/watch/eth` | Reads current Safe contract state |

---

## Safe Contract Background

### What is Safe?

[Safe (formerly Gnosis Safe)](https://safe.global/) is the most widely-used smart contract multisig wallet on Ethereum. Key properties:

- **m-of-n threshold**: Requires at least `m` out of `n` owner signatures
- **`execTransaction`**: The on-chain entry point; accepts packed signatures from all signers
- **Nonce**: The Safe contract maintains its own nonce (separate from EOA nonce) to prevent replay attacks across Safe instances
- **EIP-712 `safeTxHash`**: A typed structured data hash that uniquely identifies a proposed transaction; signers sign this hash offline

### Safe Contract ABI (relevant methods)

```solidity
// Read the current Safe nonce
function nonce() external view returns (uint256);

// Compute the EIP-712 typed hash for a proposed transaction
function getTransactionHash(
    address to,
    uint256 value,
    bytes calldata data,
    uint8 operation,
    uint256 safeTxGas,
    uint256 baseGas,
    uint256 gasPrice,
    address gasToken,
    address refundReceiver,
    uint256 _nonce
) external view returns (bytes32);

// Execute a multisig transaction
function execTransaction(
    address to,
    uint256 value,
    bytes calldata data,
    uint8 operation,
    uint256 safeTxGas,
    uint256 baseGas,
    uint256 gasPrice,
    address gasToken,
    address payable refundReceiver,
    bytes memory signatures
) external payable returns (bool success);
```

---

## EIP-712 Signing Flow

### Why EIP-712?

EIP-712 defines typed structured data hashing. The Safe contract uses it to produce a canonical hash (`safeTxHash`) from the transaction parameters that:

- Binds the hash to a specific Safe address and chain ID (prevents cross-chain replay)
- Binds the hash to the Safe's current nonce (prevents replay within the same Safe)
- Is verifiable offline without any network calls

### Online: `safeTxHash` Computation

The Watch Wallet calls `getTransactionHash` on the Safe contract to obtain the canonical `safeTxHash`:

```go
// internal/infrastructure/api/eth/safe_client.go
func (s *SafeClient) GetSafeTxHash(
    ctx context.Context,
    safeAddress, to string,
    value *big.Int,
    data []byte,
    operation uint8,
    nonce *big.Int,
) (string, error) {
    // Calls Safe.getTransactionHash() via eth_call
}
```

### Offline: Hash Verification and Signing

Each signer wallet **recomputes** the `safeTxHash` locally from the file fields using EIP-712 encoding — no network call required. This guards against a tampered file where the `safe_tx_hash` field has been replaced with a different hash.

```go
// internal/application/usecase/keygen/eth/safe_tx_hash.go
func ComputeSafeTxHash(f *dtoeth.ETHMultisigTransactionFile) ([32]byte, error) {
    // 1. Build SafeTx struct type hash
    // 2. Encode all transaction fields with ABI encoding
    // 3. Compute domainSeparator from chainID + verifyingContract (Safe address)
    // 4. Return keccak256("\x19\x01" || domainSeparator || safeTxStructHash)
}
```

If the recomputed hash does not match `safe_tx_hash` in the file, signing is aborted with `ErrSafeTxHashMismatch`.

### Signature Format

Safe requires signatures in the format `[R || S || V]` where `V ∈ {27, 28}` (Ethereum legacy convention). The standard `crypto.Sign` output has `V ∈ {0, 1}`; the signing use case adjusts this:

```go
sig, _ := crypto.Sign(computedHash[:], privKey)
sig[64] += 27  // adjust V to Safe's expected range
```

Signatures are packed in signer address order when submitting to `execTransaction`.

---

## Transaction File Format

The `ETHMultisigTransactionFile` struct (`internal/application/dto/eth/multisig_transaction_file.go`) is the canonical JSON file exchanged between wallets.

### File Naming Convention

```
{actionType}_multisig_{uuid}_{signedCount}.json

Examples:
  payment_multisig_550e8400-e29b-41d4-a716-446655440000_0.json  (unsigned)
  payment_multisig_550e8400-e29b-41d4-a716-446655440000_1.json  (1 signature)
  payment_multisig_550e8400-e29b-41d4-a716-446655440000_2.json  (threshold met → "signed")
```

### JSON Fields

| Field | Type | Description |
|-------|------|-------------|
| `version` | `int` | File format version (currently `1`) |
| `tx_type` | `string` | `"unsigned"` or `"signed"` |
| `uuid` | `string` | UUIDv4 identifying this proposal |
| `action_type` | `string` | `"deposit"`, `"payment"`, or `"transfer"` |
| `safe_address` | `string` | EIP-55 checksummed Safe proxy address |
| `to` | `string` | EIP-55 checksummed recipient address |
| `value` | `string` | Wei amount as decimal string |
| `data` | `string` | `"0x"` for plain ETH transfer |
| `operation` | `uint8` | `0` = Call, `1` = DelegateCall |
| `safe_tx_gas` | `string` | `"0"` for plain ETH transfer |
| `base_gas` | `string` | `"0"` for plain ETH transfer |
| `gas_price` | `string` | `"0"` for plain ETH transfer |
| `gas_token` | `string` | Zero address for ETH gas |
| `refund_receiver` | `string` | Zero address |
| `nonce` | `string` | Safe contract nonce as decimal string |
| `chain_id` | `uint64` | EIP-155 chain identifier |
| `safe_tx_hash` | `string` | `0x`-prefixed EIP-712 hash (authoritative) |
| `threshold` | `int` | Required signer count (m in m-of-n) |
| `signatures` | `[]SignEntry` | Accumulated signer entries |

### SignEntry

```json
{
  "signer_address": "0xAbCd...",
  "signature_hex": "0x{R 32 bytes}{S 32 bytes}{V 1 byte}"
}
```

### Example File

```json
{
  "version": 1,
  "tx_type": "signed",
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "action_type": "payment",
  "safe_address": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
  "to": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
  "value": "100000000000000000",
  "data": "0x",
  "operation": 0,
  "safe_tx_gas": "0",
  "base_gas": "0",
  "gas_price": "0",
  "gas_token": "0x0000000000000000000000000000000000000000",
  "refund_receiver": "0x0000000000000000000000000000000000000000",
  "nonce": "0",
  "chain_id": 31337,
  "safe_tx_hash": "0xabcdef...",
  "threshold": 2,
  "signatures": [
    {
      "signer_address": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
      "signature_hex": "0x..."
    },
    {
      "signer_address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
      "signature_hex": "0x..."
    }
  ]
}
```

### State Machine

```
tx_type: "unsigned"  (len(signatures) < threshold)
    → (signer appends signature, threshold reached)
tx_type: "signed"    (len(signatures) >= threshold)
```

---

## CLI Commands

### Watch Wallet

#### Propose a multisig transaction

```bash
watch create multisig \
  --safe-address 0xSafeAddr \
  --to 0xRecipient \
  --amount 0.1 \
  --threshold 2 \
  --action-type payment
```

Writes: `payment_multisig_{uuid}_0.json`

#### Submit a fully-signed transaction

```bash
watch send multisig send-eth --file /path/to/payment_multisig_{uuid}_2.json
```

Calls `execTransaction` on the Safe. Gas is paid by the configured `safe_gas_payer_hex_key` account.

#### Inspect Safe contract state

```bash
watch safe info --safe-address 0xSafeAddr
```

Prints: owners, threshold, nonce, version.

### Keygen Wallet (first signer)

```bash
keygen sign signature \
  --file /path/to/payment_multisig_{uuid}_0.json \
  --signer-address 0xSignerAddr
```

Writes: `payment_multisig_{uuid}_1.json`

### Sign Wallet (additional signers)

```bash
sign sign signature \
  --file /path/to/payment_multisig_{uuid}_1.json \
  --signer-address 0xSignerAddr
```

Writes: `payment_multisig_{uuid}_2.json` (if threshold = 2)

---

## Wallet Roles in Multisig

| Wallet | Role | Network |
|--------|------|---------|
| **Watch** | Propose transaction, submit `execTransaction` | Online |
| **Keygen** | Sign as owner 1 (offline, air-gapped) | Offline |
| **Sign** | Sign as owner 2…n (offline, air-gapped) | Offline |

The Sign Wallet **is used for ETH Safe multisig** (unlike the single-sig EOA flow where only Keygen is needed).

---

## E2E Pattern 3: 2-of-2 Safe Payment

Pattern 3 demonstrates a complete Safe 2-of-2 multisig payment flow on a local Anvil node.

### Setup

| Component | Account | Purpose |
|-----------|---------|---------|
| Deployer | Anvil account 2 | Deploys Safe contract via Foundry `forge script` |
| Safe owner 1 | Keygen wallet address | Signs the proposal |
| Safe owner 2 | Sign wallet address | Cosigns the proposal |
| Gas payer | Anvil account 2 | Pays ETH gas for `execTransaction` |

### Sequence

```
1. Deploy Safe proxy (Foundry forge script)
   - factory: SafeProxyFactory
   - singleton: Safe v1.4.1
   - owners: [keygen_addr, sign_addr], threshold: 2

2. Fund Safe with ETH (anvil_setBalance)

3. Watch: create multisig → payment_multisig_{uuid}_0.json

4. Keygen: sign signature → payment_multisig_{uuid}_1.json

5. Sign: sign signature → payment_multisig_{uuid}_2.json

6. Watch: send multisig send-eth → execTransaction on-chain
```

### Make Targets

| Target | Description |
|--------|-------------|
| `make eth-e2e-p3` | Run full P3 workflow |
| `make eth-e2e-p3-reset` | Reset Anvil state and re-run |
| `make eth-e2e-p3-ci` | CI-mode run |

### Parallel Execution

P3 uses **Anvil account 2** as its deployer/gas-payer, ensuring it does not conflict with P1 (no deployer) or P2 (Anvil account 1) in parallel runs:

```bash
make eth-e2e-parallel PATTERNS=1-3 MAX_PARALLEL=3
```

---

## Configuration

### Config struct (`pkg/config/wallet.go`)

```go
type Ethereum struct {
    // ...existing fields...
    SafeGasPayerHexKey string `toml:"safe_gas_payer_hex_key" env:"WALLET_ETHEREUM_SAFE_GAS_PAYER_HEX_KEY"`
}
```

`SafeGasPayerHexKey` is the raw hex private key of the EOA that submits `execTransaction` and pays gas. It is intentionally separate from the Safe owners (an owner may also be the gas payer in test environments).

### Environment Variable

```bash
export WALLET_ETHEREUM_SAFE_GAS_PAYER_HEX_KEY=0x{64-hex-chars}
```

---

## Key Implementation Files

| File | Description |
|------|-------------|
| `internal/application/dto/eth/multisig_transaction_file.go` | `ETHMultisigTransactionFile` struct, validation, sentinel errors |
| `internal/application/ports/api/eth/interface.go` | `SafeNonceReader`, `SafeTxHashComputer`, `SafeExecutor`, `SafeInfoReader` port interfaces |
| `internal/application/ports/file/multisig_file.go` | `MultisigFileRepositorier` port interface |
| `internal/application/usecase/watch/eth/create_multisig_transaction.go` | Propose use case (online, reads Safe nonce) |
| `internal/application/usecase/watch/eth/send_multisig_transaction.go` | Submit use case (online, calls execTransaction) |
| `internal/application/usecase/watch/eth/safe_info.go` | Safe info use case (online, reads owners/threshold) |
| `internal/application/usecase/keygen/eth/sign_multisig_transaction.go` | Sign use case (offline, EIP-712 recomputation) |
| `internal/application/usecase/keygen/eth/safe_tx_hash.go` | EIP-712 `safeTxHash` offline computation |
| `internal/infrastructure/api/eth/safe_client.go` | `SafeClient` — wraps Safe contract ABI via `eth_call` |
| `internal/infrastructure/storage/file/transaction/multisig.go` | JSON file read/write implementation |
| `internal/di/container.go` | DI wiring for `SafeClient` and multisig use cases |
| `apps/eth-contracts/script/DeploySafe.s.sol` | Foundry deployment script for Safe proxy |
| `scripts/operation/eth/e2e/e2e-p3.sh` | E2E Pattern 3 shell script |
| `make/wallet/eth_e2e.mk` | Make targets for P3 and parallel E2E |

---

## Official References

| Resource | URL |
|----------|-----|
| Safe documentation | https://docs.safe.global/ |
| Safe contracts v1.4.1 | https://github.com/safe-global/safe-smart-account |
| EIP-712: Typed structured data hashing | https://eips.ethereum.org/EIPS/eip-712 |
| EIP-155: Replay attack protection | https://eips.ethereum.org/EIPS/eip-155 |

---

**Document Version:** 1.0
**Last Updated:** 2026-03-08
**Maintainer:** go-crypto-wallet team
