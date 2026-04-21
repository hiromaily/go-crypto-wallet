<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/eth/transaction-patterns.tpl.md · Run `make docs` to regenerate.
-->

# Ethereum Raw Transaction Patterns (2026)

## 1. Transaction Types (Protocol Level)

### Type 0 — Legacy Transaction

Standard before EIP-1559.

EIPs: EIP-155

Fields

```
nonce
gasPrice
gasLimit
to
value
data
v
r
s
```

Notes

* Uses `gasPrice`
* No fee market
* Mostly deprecated but still valid

---

### Type 1 — Access List Transaction

EIPs: EIP-2930

Fields

```
chainId
nonce
gasPrice
gasLimit
to
value
data
accessList
v
r
s
```

Notes

* Optional access list
* Gas optimization
* Rarely used in practice

---

### Type 2 — EIP-1559 Transaction (Standard)

EIPs: EIP-1559

Fields

```
chainId
nonce
maxPriorityFeePerGas
maxFeePerGas
gasLimit
to
value
data
accessList
v
r
s
```

Notes

* Current default transaction type
* Uses dynamic fee market
* Recommended for most applications

---

### Type 3 — Blob Transaction

EIPs: EIP-4844

Fields

```
chainId
nonce
maxPriorityFeePerGas
maxFeePerGas
gasLimit
to
value
data
accessList

maxFeePerBlobGas
blobVersionedHashes

v
r
s
```

External blob data

```
blobs
commitments
proofs
```

Notes

* Designed for rollup data availability
* Used by L2 networks such as

  * Optimism
  * Arbitrum
  * Base

---

### Future / Emerging

Example

* EIP-7702

Purpose

* Account abstraction
* Temporary code execution for EOAs

---

## 2. Transaction Intent Patterns (Application Level)

These patterns exist regardless of transaction type.

---

### ETH Transfer

Simple native token transfer.

Pattern

```
to = recipient
value > 0
data = empty
```

Used for

```
wallet transfer
exchange withdrawal
payment
```

---

### Smart Contract Call

Invoke contract method.

Pattern

```
to = contract address
data = ABI encoded function call
value = optional
```

Examples

```
DEX swap
NFT mint
DeFi interaction
```

---

### Contract Deployment

Deploy new smart contract.

Pattern

```
to = null
data = contract bytecode + constructor args
value = optional
```

---

### ERC-20 Token Transfer

Token transfer using contract.

Example contract standard

* ERC-20

Pattern

```
to = token contract
data = transfer(address,uint256)
value = 0
```

---

### ERC-721 / NFT Transfer

NFT transfer.

Standard

* ERC-721

Pattern

```
to = NFT contract
data = safeTransferFrom(...)
value = 0
```

---

## 3. Summary Matrix

| Category        | Pattern         | Typical Type |
| --------------- | --------------- | ------------ |
| ETH transfer    | send ETH        | Type 2       |
| Contract call   | DeFi / swap     | Type 2       |
| Contract deploy | deploy bytecode | Type 2       |
| ERC20 transfer  | token transfer  | Type 2       |
| NFT transfer    | ERC721 transfer | Type 2       |
| Rollup blob     | L2 data publish | Type 3       |

---

## 4. Minimal Patterns AI Agents Should Understand

If building wallet / signing infrastructure for Ethereum, the **minimal required patterns** are:

```
1. ETH transfer
2. ERC20 transfer
3. Contract call
4. Contract deployment
5. Blob transaction (L2)
```

Transaction type support

```
Type 2 (EIP-1559)  ← default
Type 3 (Blob)      ← L2 support
```
