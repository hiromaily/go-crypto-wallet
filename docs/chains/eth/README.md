# Ethereum (ETH) Technical Reference

This document provides a comprehensive technical reference for Ethereum implementation in the go-crypto-wallet system. It covers specifications, protocol details, and implementation notes to help AI agents and developers understand Ethereum's architecture and implement features correctly.

## Architecture & Transaction Flow

| Document | Description |
|----------|-------------|
| [architecture.md](./architecture.md) | **ETH wallet architecture** — wallet roles, use case boundary map, Clean Architecture layers, KeygenSignTx offline signing detail |
| [docs/transaction-flow.md](../../transaction-flow.md) | Chain-agnostic 3-wallet setup, signing, and monitoring flows |

For Ethereum-specific concerns on top of the common flow, see [ETH-Specific Flow Details](#eth-specific-flow-details) below.

> **Key point:** ETH uses single-sig EOA only. Only Watch and Keygen wallets are required.
> The Keygen Wallet performs transaction signing. The Sign Wallet is not used for ETH.

---

## Table of Contents

1. [Overview](#overview)
2. [Core Specifications](#core-specifications)
3. [Address Types & Key Derivation](#address-types--key-derivation)
4. [Transaction Architecture](#transaction-architecture)
5. [Transaction Types](#transaction-types)
6. [Signing Mechanism](#signing-mechanism)
7. [ERC-20 Token Support](#erc-20-token-support)
8. [Network & Confirmation Guidelines](#network--confirmation-guidelines)
9. [Fee Management](#fee-management)
10. [Wallet Implementation](#wallet-implementation)
11. [RPC & API Reference](#rpc--api-reference)
12. [Security Considerations](#security-considerations)
13. [Testing Resources](#testing-resources)
14. [Official References](#official-references)
15. [ETH-Specific Flow Details](#eth-specific-flow-details)

---

## Overview

### What is Ethereum?

Ethereum is a decentralized platform with smart contract functionality. Unlike Bitcoin's UTXO model, Ethereum uses an **account-based model** where each account has a balance and a nonce. Transactions modify account state rather than consuming and creating outputs.

### Key Characteristics (2026)

| Property | Value |
|----------|-------|
| **Launch Date** | July 30, 2015 |
| **Block Time** | ~12 seconds (post-Merge) |
| **Consensus** | Proof of Stake (since September 2022 Merge) |
| **Native Currency** | Ether (ETH) |
| **Smallest Unit** | Wei (1 ETH = 10^18 Wei) |
| **Cryptographic Curve** | secp256k1 |
| **Signature Algorithm** | ECDSA |
| **Address Length** | 20 bytes (40 hex characters) |

### Protocol Upgrade Timeline

| Year | Upgrade | Key Features |
|------|---------|--------------|
| 2019 | **Istanbul** | EIP-1884, opcode repricing |
| 2021 | **London (EIP-1559)** | Base fee + priority fee model |
| 2022 | **The Merge** | Proof of Stake transition |
| 2023 | **Shanghai** | Staking withdrawals |
| 2024 | **Dencun** | Proto-danksharding (EIP-4844) |

---

## Core Specifications

### Cryptographic Primitives

#### Elliptic Curve (secp256k1)

Ethereum uses the same secp256k1 elliptic curve as Bitcoin for signing:

```
Private Key:  32 bytes (256-bit scalar)
Public Key:   64 bytes (uncompressed, without 04 prefix) or 65 bytes with prefix
Address:      20 bytes = Keccak256(pubkey)[12:]
```

#### Hash Functions

| Function | Usage |
|----------|-------|
| **Keccak256** | Address derivation, transaction hashing, function selectors |
| **RLP** | Transaction serialization encoding |
| **SHA3** | Generic hashing (note: Ethereum uses Keccak256, not NIST SHA3) |

#### Address Derivation

```
1. Generate secp256k1 key pair
2. Take uncompressed public key (64 bytes, no prefix)
3. Compute Keccak256 hash (32 bytes)
4. Take last 20 bytes
5. Apply EIP-55 checksum encoding
```

### Data Encoding

| Format | Usage |
|--------|-------|
| **RLP** | Transaction serialization |
| **Hex** | Addresses, transaction data (`0x` prefix) |
| **EIP-55** | Mixed-case checksum address encoding |
| **ABI** | Smart contract call data encoding |

---

## Address Types & Key Derivation

### Address Type

This system supports **EOA (Externally Owned Account)** addresses only. Smart contract wallets are not supported.

EOA addresses:

- Are 20 bytes (40 hex chars) with `0x` prefix
- Example: `0x742d35Cc6634C0532925a3b8BC9e7595f0bEb123`
- Validated via `common.IsHexAddress(addr)`

### HD Wallet Derivation Path

**Standard:** BIP44 with Ethereum coin type

```
m / 44' / 60' / account' / 0 / address_index
```

| Component | Value | Notes |
|-----------|-------|-------|
| **Purpose** | 44' | BIP44 |
| **Coin Type** | 60' | SLIP-0044 Ethereum |
| **Account** | See table below | Non-hardened for deposit/payment/stored |
| **Change** | 0 | External chain only (no internal change addresses) |
| **Index** | 0..N | Sequential address generation |

### Account Types and Derivation Paths

| Account | Index | Hardened | Path | Use |
|---------|-------|----------|------|-----|
| **deposit** | 0 | No | `m/44'/60'/0'/0/i` | Aggregate client funds |
| **payment** | 1 | No | `m/44'/60'/1'/0/i` | Outgoing payments |
| **stored** | 2 | No | `m/44'/60'/2'/0/i` | Cold storage |
| **auth1** | 11 | Yes | `m/44'/60'/11''/0/i` | Authentication key 1 |
| **auth2** | 12 | Yes | `m/44'/60'/12''/0/i` | Authentication key 2 |

> **Note:** Deposit/payment/stored use non-hardened derivation so the Watch Wallet can derive child addresses from the extended public key (xpub) without private keys. Auth accounts use hardened derivation for enhanced security.

---

## Transaction Architecture

### Account Model

Ethereum uses an account-based model (not UTXO):

```
Account State = {
    nonce:    uint64    // Number of transactions sent
    balance:  *big.Int  // Wei balance
    codeHash: []byte    // Empty for EOA
    storage:  map       // Empty for EOA
}
```

**Key Concepts:**

- Each transaction from an address must use the next sequential nonce
- Transactions are identified by hash (not UTXO references)
- Balance changes are atomic (no partial spending)

### Nonce Management

Critical for transaction ordering and double-spend prevention:

```go
// Get pending nonce (includes in-flight transactions)
nonce, err = eth.GetTransactionCount(ctx, fromAddr, QuantityTagPending)

// For batch operations, increment additionalNonce per transaction
effectiveNonce = nonce + additionalNonce
```

**Rules:**

- Always use `pending` state to account for unconfirmed transactions
- Nonces must be sequential — gaps cause transactions to be stuck
- Increment `additionalNonce` when creating multiple transactions in the same batch

### Gas

Gas is the unit of computation cost on Ethereum:

```
Transaction Fee = Gas Used × Effective Gas Price

// Legacy (pre-EIP-1559):
Fee = gasLimit × gasPrice

// EIP-1559 (post-London):
Fee = gasLimit × (baseFee + priorityFee)
```

**Gas Estimation:**

```go
estimatedGas, err = eth.EstimateGas(ctx, msg)
```

---

## Transaction Types

### Legacy Transactions (Type 0)

Uses a single `gasPrice` field. Still supported on all networks.

```go
types.LegacyTx{
    Nonce:    nonce,
    GasPrice: gasPrice,
    Gas:      gasLimit,
    To:       &toAddr,
    Value:    amount,
    Data:     data,
}
```

### EIP-1559 Transactions (Type 2)

Introduced in the London hard fork (2021). Separates gas price into base fee (burned) + priority fee (to validator).

```go
types.DynamicFeeTx{
    ChainID:   chainID,
    Nonce:     nonce,
    GasTipCap: maxPriorityFeePerGas,  // tip to validator
    GasFeeCap: maxFeePerGas,           // max total fee willing to pay
    Gas:       gasLimit,
    To:        &toAddr,
    Value:     amount,
    Data:      data,
}
```

**Fee Calculation (EIP-1559):**

```go
maxPriorityFeePerGas = 2 Gwei  // configurable via config.Ethereum.MaxPriorityFeePerGas
maxFeePerGas = (baseFeePerGas * 2) + maxPriorityFeePerGas
```

**EIP-1559 Detection:**

```go
// Check if network supports EIP-1559 by presence of baseFeePerGas in latest block
supported = eth.SupportsEIP1559(ctx)
```

**Current Implementation Note:** The Watch Wallet currently creates Legacy transactions (`types.LegacyTx`). The EIP-1559 path (`CreateRawTransactionEIP1559`) exists in the infrastructure layer but is not yet used by the Watch use case.

---

## Signing Mechanism

### Signer Types

| Signer | Transaction Type | Description |
|--------|-----------------|-------------|
| `types.HomesteadSigner` | Legacy | Original signer (no replay protection) |
| `types.NewEIP155Signer(chainID)` | Legacy | Replay protection (EIP-155) |
| `types.NewLondonSigner(chainID)` | Legacy + EIP-1559 | Used in this system |

### Signing Flow

```go
// 1. Decode hex transaction
tx := new(types.Transaction)
rlp.DecodeBytes(rawBytes, tx)

// 2. Get private key from keystore
privateKey := keystore.GetPrivKey(address, password)

// 3. Determine chain ID
// netID 1 → chainID 1 (mainnet)
// other → chainID 4 (testnet/regtest)

// 4. Create London signer (supports both legacy and EIP-1559)
signer := types.NewLondonSigner(chainID)

// 5. Sign transaction
signedTx, err := types.SignTx(tx, signer, privateKey)

// 6. Verify sender address
sender, err := types.Sender(signer, signedTx)

// 7. Encode to hex for file output
```

### Key Storage

Private keys are stored in the **go-ethereum keystore** format on the local filesystem:

```
File: UTC--{timestamp}--{address}
Path: Configurable (default: ./data/keystore/)
Encryption: scrypt (StandardScryptN, StandardScryptP)
Password: Required to decrypt (set in configuration)
```

**Key Import (Keygen Wallet):**

```go
// Uses local filesystem keystore (compatible with Anvil)
// Does NOT use personal_importRawKey RPC
ks := keystore.NewKeyStore(keyDir, keystore.StandardScryptN, keystore.StandardScryptP)
account, err := ks.ImportECDSA(ecdsaPrivKey, password)
```

---

## ERC-20 Token Support

### Overview

ERC-20 token transfers are supported via the `erc20` infrastructure package.

**Key Difference from ETH transfers:**

| | ETH Transfer | ERC-20 Transfer |
|-|-------------|-----------------|
| **To address** | Recipient | Contract address |
| **Value** | Transfer amount | 0 (always) |
| **Data** | Empty | ABI-encoded `transfer()` call |
| **Fee payer** | Sender | Master address |

### Function Selector

ERC-20 `transfer(address,uint256)` is encoded as:

```go
// Function selector = first 4 bytes of Keccak256("transfer(address,uint256)")
selector := crypto.Keccak256([]byte("transfer(address,uint256)"))[:4]

// Call data = selector + padded address (32 bytes) + padded amount (32 bytes)
data = append(selector, paddedAddress...)
data = append(data, paddedAmount...)
```

### Limitations

- Requires prior `approve()` transaction for some flows
- Master address pays ETH gas fees (not token holder)
- No ERC-721 or ERC-1155 support

---

## Network & Confirmation Guidelines

### Supported Networks

| Network | Chain ID | Purpose | RPC Port |
|---------|----------|---------|----------|
| **Mainnet** | 1 | Production | 8545 |
| **Sepolia** | 11155111 | Public testnet | 8545 |
| **Holesky** | 17000 | Staking testnet | 8545 |
| **Anvil (local)** | 31337 | Local development | 8545 |

**Chain ID Mapping in this system:**

```go
// netID 1 → chainID 1 (mainnet)
// other → chainID 4 (treated as testnet; historical Rinkeby ID)
```

### Confirmation Guidelines

| Confirmations | Risk Level | Typical Use Case |
|---------------|------------|------------------|
| 0 (pending) | High | Development only |
| 1 | Medium | Low-value transfers |
| 6 | Low | Standard commerce |
| 12 | Very Low | Large transactions |
| 32+ | Negligible | High-value custody |

---

## Fee Management

### Fee Components (EIP-1559)

| Component | Description | Set By |
|-----------|-------------|--------|
| **Base Fee** | Burned by protocol, adjusts per block | Protocol |
| **Priority Fee** (tip) | Paid to validator | Sender (default: 2 Gwei) |
| **Max Fee** | Maximum total gas price willing to pay | Sender |

### Fee Sources

| Source | Method |
|--------|--------|
| **Gas Price (legacy)** | `eth_gasPrice` RPC |
| **Base Fee** | Latest block header (`baseFeePerGas`) |
| **Priority Fee** | Configurable (default 2 Gwei) |
| **Gas Estimate** | `eth_estimateGas` RPC |

### Fee Calculation per Transaction Type

**Deposit / Transfer (full balance sweep):**

```
transferAmount = balance - (gasPrice × estimatedGas)
```

**Payment (fixed amount):**

```
senderCost = amount + (gasPrice × estimatedGas)
recipientReceives = amount
```

---

## Wallet Implementation

### Wallet Roles

| Wallet | Role | Network |
|--------|------|---------|
| **Watch** | Create transactions, broadcast, monitor | Online |
| **Keygen** | Generate keys, import to keystore, sign (first sig) | Offline (air-gapped) |
| **Sign** | Additional signatures | Offline (air-gapped) |

### Keygen Wallet Operations

1. **`create seed`** — Generate BIP39 mnemonic and store encrypted seed
2. **`create hdkey`** — Derive BIP44 HD keys for specified account
3. **`import privkey`** — Import private keys into local keystore (encrypted)
4. **`export address`** — Export public addresses to file for Watch Wallet

### Watch Wallet Operations

1. **`import address`** — Import addresses from Keygen Wallet file
2. **`create deposit/payment/transfer`** — Build unsigned transactions
3. **`send`** — Broadcast signed transactions
4. **`monitor`** — Track transaction confirmation status

### Sign Wallet Operations

1. **`sign`** — Sign unsigned transaction files using keystore private keys

### Database Schema (ETH-specific tables)

| Table | Description |
|-------|-------------|
| `account_key` | Key generation tracking with status |
| `address` | Generated Ethereum addresses |
| `full_pubkey` | Extended public keys |
| `eth_tx` | Unsigned/signed transaction records |
| `eth_detail_tx` | Per-transaction detail (nonce, gas, amounts) |
| `payment_request` | Outbound payment requests |

### Address Status Lifecycle

```
AddrStatusHDKeyGenerated
    → (import privkey)
AddrStatusPrivKeyImported
```

### Transaction Status Lifecycle

```
TxTypeSent ──(confirmations >= threshold)──> TxTypeDone ──(notified)──> TxTypeNotified
```

---

## RPC & API Reference

### Key go-ethereum Methods Used

| Method | Description |
|--------|-------------|
| `ethclient.BalanceAt` | Get account balance at block |
| `ethclient.PendingNonceAt` | Get pending nonce |
| `ethclient.SendTransaction` | Broadcast raw transaction |
| `eth_gasPrice` | Get current gas price (legacy) |
| `eth_estimateGas` | Estimate gas for transaction |
| `eth_getTransactionCount` | Get nonce (use `pending` tag) |
| `eth_getTransactionReceipt` | Get receipt (for confirmation status) |
| `eth_getBlockByNumber` | Get block (for baseFeePerGas detection) |
| `eth_call` | Call contract function (read-only) |
| `net_version` | Get network/chain ID |

### Go Libraries

| Library | Version | Purpose |
|---------|---------|---------|
| **go-ethereum** | v1.17.0 | Ethereum client library |
| `ethclient` | — | HTTP/WS client |
| `core/types` | — | Transaction types (LegacyTx, DynamicFeeTx) |
| `crypto` | — | Keccak256, PubkeyToAddress |
| `accounts/keystore` | — | Encrypted key storage |
| `common` | — | HexToAddress, IsHexAddress |
| `rlp` | — | RLP encode/decode |
| **btcd/btcec/v2** | — | ECDSA private key operations |
| **btcd/btcutil/hdkeychain** | — | BIP44 HD key derivation |

### Application Port Interface

The `Ethereumer` interface (defined in `internal/application/ports/api/eth/interface.go`) provides:

- Balance operations: `GetTotalBalance`, `BalanceAt`
- Network: `Syncing`, `ProtocolVersion`, `NetVersion`
- Block: `BlockNumber`, `GetBlockByNumber`, `EnsureBlockNumber`
- Transaction: `CreateRawTransaction`, `CreateRawTransactionEIP1559`, `SendRawTransaction`
- Gas: `GasPrice`, `EstimateGas`
- Signing: `SignOnRawTransaction`
- Key management: `ToECDSA`, `GetPrivKey`, `ImportRawKey`

---

## Security Considerations

### Private Key Security

- **NEVER** log or expose private keys
- Private keys only exist on offline Keygen/Sign wallets
- Keys stored encrypted via scrypt in local keystore files
- Watch Wallet holds no private keys — only public addresses

### Keystore Security

```
// Encryption: scrypt with standard parameters
N = 1 << 18 (StandardScryptN)
P = 1       (StandardScryptP)
r = 8
```

> **Warning:** The keystore password is currently hardcoded in the codebase. This must be externalized to a secure configuration source before production use.

### Transaction Security

- Always validate destination address with `common.IsHexAddress(addr)`
- Verify nonce before signing to prevent replay
- Check sender address after signing (compare recovered address with expected)
- EIP-155 replay protection is enforced via `LondonSigner`

### Anvil / Local Development

- Anvil is detected by version string containing "Anvil"
- Key operations use local filesystem (not RPC) for Anvil compatibility
- Chain ID is forced to testnet value for non-mainnet environments

---

## Testing Resources

### Local Development with Anvil

Anvil (from Foundry) is the recommended local Ethereum node for development:

```bash
# Start Anvil with deterministic accounts
anvil --deterministic

# With specific chain ID
anvil --chain-id 1337

# Fork mainnet state
anvil --fork-url https://mainnet.infura.io/v3/<key>
```

### Block Explorers

| Network | Explorer |
|---------|----------|
| **Mainnet** | [etherscan.io](https://etherscan.io/) |
| **Sepolia** | [sepolia.etherscan.io](https://sepolia.etherscan.io/) |
| **Holesky** | [holesky.etherscan.io](https://holesky.etherscan.io/) |

### Testnet Faucets

| Network | Faucet |
|---------|--------|
| **Sepolia** | [sepoliafaucet.com](https://sepoliafaucet.com/) |
| **Holesky** | [holesky-faucet.pk910.de](https://holesky-faucet.pk910.de/) |

---

## Official References

### Ethereum Improvement Proposals (EIPs)

| EIP | Title | Relevance |
|-----|-------|-----------|
| [EIP-55](https://eips.ethereum.org/EIPS/eip-55) | Mixed-case checksum address encoding | Address format |
| [EIP-155](https://eips.ethereum.org/EIPS/eip-155) | Simple replay attack protection | Chain ID in signatures |
| [EIP-191](https://eips.ethereum.org/EIPS/eip-191) | Signed data standard | Message signing |
| [EIP-1559](https://eips.ethereum.org/EIPS/eip-1559) | Fee market change | Base fee + priority fee |
| [EIP-2718](https://eips.ethereum.org/EIPS/eip-2718) | Typed Transaction Envelope | Transaction type field |
| [EIP-2930](https://eips.ethereum.org/EIPS/eip-2930) | Access list transactions (Type 1) | Optional optimization |
| [EIP-4844](https://eips.ethereum.org/EIPS/eip-4844) | Shard blob transactions | L2 scaling |

### BIP Standards Used

| BIP | Title | Usage |
|-----|-------|-------|
| [BIP32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki) | HD Wallets | Key derivation |
| [BIP39](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) | Mnemonic Codes | Seed generation |
| [BIP44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki) | Multi-Account HD | Derivation path standard |
| [SLIP-0044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md) | Coin Types | Coin type 60 = Ethereum |

### Official Documentation

- [Ethereum Developer Documentation](https://ethereum.org/en/developers/docs/)
- [go-ethereum Documentation](https://geth.ethereum.org/docs/)
- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [ERC-20 Standard](https://eips.ethereum.org/EIPS/eip-20)

---

## ETH-Specific Flow Details

> The common 3-wallet setup, signing, and monitoring flows are defined in the chain-agnostic reference:
> [docs/transaction-flow.md](../../transaction-flow.md).
> This section describes Ethereum-specific concerns on top of that common flow.

### Single-Sig Flow (Ethereum)

Follows the [common single-sig flow](../../transaction-flow.md#single-sig-flow).

Ethereum-specific steps:

- Watch Wallet fetches the pending nonce via `eth_getTransactionCount` with `pending` tag
- Watch Wallet fetches gas price via `eth_gasPrice` and estimates gas via `eth_estimateGas`
- Watch Wallet serializes the transaction as RLP-encoded hex
- Keygen Wallet decodes the hex, signs with `types.NewLondonSigner(chainID)`, re-encodes to hex
- Watch Wallet sends via `eth_sendRawTransaction`
- Watch Wallet polls for receipt via `eth_getTransactionReceipt` to confirm finality

### Transaction Types vs Use Cases

| Use Case | Watch Action | Nonce Source | Amount Calculation |
|----------|-------------|-------------|-------------------|
| **Deposit** | Client address → deposit address | Pending nonce of client address | `balance - fee` |
| **Payment** | Payment address → external address | Pending nonce of payment address | Fixed amount from payment_request |
| **Transfer** | deposit/stored → other internal | Pending nonce of source address | `balance - fee` or fixed amount |

### Multisig

> Ethereum EOA does not natively support multisig at the protocol level (unlike Bitcoin's P2SH/P2WSH).
> This system's multisig support for Ethereum would require a smart contract (e.g., Gnosis Safe).
> The current implementation uses **single-sig EOA only**.

---

**Document Version:** 1.0
**Last Updated:** 2026-02-24
**Maintainer:** go-crypto-wallet team
