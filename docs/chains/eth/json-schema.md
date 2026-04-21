<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/eth/json-schema.tpl.md · Run `make docs` to regenerate.
-->

# Ethereum Raw Transaction JSON (AI-Agent Friendly Specification)

## 1. Overview

A raw transaction represents a signed operation that will be broadcast to an Ethereum node.

Typical lifecycle:

```
build transaction
→ sign transaction
→ serialize
→ sendRawTransaction
```

Supported transaction types are defined by:

* EIP-2718

Most applications should use **Type 2 (EIP-1559)**.

---

## 2. Type 2 Transaction JSON (Standard)

**Structure**

```json
{
  "type": "0x2",
  "chainId": "0x1",
  "nonce": "0x0",
  "maxPriorityFeePerGas": "0x3b9aca00",
  "maxFeePerGas": "0x77359400",
  "gas": "0x5208",
  "to": "0xRecipientAddress",
  "value": "0xde0b6b3a7640000",
  "data": "0x",
  "accessList": [],
  "v": "0x1",
  "r": "0x...",
  "s": "0x..."
}
```

---

## 3. Field Description

| Field                  | Required        | Description                           |
| ---------------------- | --------------- | ------------------------------------- |
| `type`                 | yes             | Transaction type (`0x2` for EIP-1559) |
| `chainId`              | yes             | Ethereum chain ID                     |
| `nonce`                | yes             | Sender transaction count              |
| `maxPriorityFeePerGas` | yes             | Tip for validator                     |
| `maxFeePerGas`         | yes             | Maximum total gas fee                 |
| `gas`                  | yes             | Gas limit                             |
| `to`                   | optional        | Recipient address                     |
| `value`                | optional        | ETH amount in wei                     |
| `data`                 | optional        | Contract call data                    |
| `accessList`           | optional        | Access list (EIP-2930)                |
| `v`                    | yes (signed tx) | Recovery id                           |
| `r`                    | yes (signed tx) | Signature component                   |
| `s`                    | yes (signed tx) | Signature component                   |

---

## 4. Special Cases

### ETH Transfer

```
data = "0x"
to = recipient
value > 0
```

Example

```json
{
  "type": "0x2",
  "chainId": "0x1",
  "nonce": "0x1",
  "maxPriorityFeePerGas": "0x59682f00",
  "maxFeePerGas": "0x59682f00",
  "gas": "0x5208",
  "to": "0xReceiverAddress",
  "value": "0xde0b6b3a7640000",
  "data": "0x",
  "accessList": []
}
```

---

### Smart Contract Call

```
to = contract address
data = ABI encoded function
value = optional
```

Example

```json
{
  "type": "0x2",
  "chainId": "0x1",
  "nonce": "0x2",
  "maxPriorityFeePerGas": "0x59682f00",
  "maxFeePerGas": "0x77359400",
  "gas": "0x186a0",
  "to": "0xContractAddress",
  "value": "0x0",
  "data": "0xa9059cbb000000000000000000000000...",
  "accessList": []
}
```

---

### Contract Deployment

```
to = null
data = contract bytecode
```

Example

```json
{
  "type": "0x2",
  "chainId": "0x1",
  "nonce": "0x3",
  "maxPriorityFeePerGas": "0x59682f00",
  "maxFeePerGas": "0x77359400",
  "gas": "0x2dc6c0",
  "to": null,
  "value": "0x0",
  "data": "0x6080604052348015610010...",
  "accessList": []
}
```

---

## 5. Signing Flow

Signing payload excludes signature fields.

Unsigned transaction:

```json
{
  "type": "0x2",
  "chainId": "0x1",
  "nonce": "0x1",
  "maxPriorityFeePerGas": "...",
  "maxFeePerGas": "...",
  "gas": "...",
  "to": "...",
  "value": "...",
  "data": "...",
  "accessList": []
}
```

Steps

```
1. encode transaction
2. hash transaction
3. sign with secp256k1
4. append r,s,v
5. broadcast via eth_sendRawTransaction
```

---

## 6. Blob Transaction (Type 3)

Defined by

* EIP-4844

Example structure

```json
{
  "type": "0x3",
  "chainId": "0x1",
  "nonce": "0x1",
  "maxPriorityFeePerGas": "...",
  "maxFeePerGas": "...",
  "maxFeePerBlobGas": "...",
  "gas": "...",
  "to": "...",
  "value": "0x0",
  "data": "0x",
  "blobVersionedHashes": [],
  "accessList": []
}
```

Used mainly by L2 rollups.

---

## 7. Typical Wallet Support

Minimal support required for Ethereum wallets:

```
Transaction Types
- Type 2 (EIP-1559)
- Type 3 (Blob)

Transaction Intents
- ETH transfer
- ERC20 transfer
- Contract call
- Contract deployment
```
