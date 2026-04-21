# Overview

## What is MuSig2?

MuSig2 is a cryptographic protocol that enables multiple parties to create a **single aggregated Schnorr signature** that looks identical to a single-signature transaction on the blockchain. Unlike traditional multisig (P2SH, P2WSH) where multiple signatures are stored on-chain, MuSig2 aggregates multiple signatures into one, providing significant benefits.

## Benefits Over Traditional Multisig

| Feature | Traditional P2WSH Multisig | MuSig2 |
|---------|---------------------------|--------|
| **On-Chain Appearance** | Multiple signatures visible | Single signature (looks like single-sig) |
| **Transaction Size** | ~370-400 bytes (2-of-3) | ~200-250 bytes (30-50% smaller) |
| **Privacy** | Multisig is visible | Indistinguishable from single-sig |
| **Fees** | Higher (proportional to size) | 30-50% lower |
| **Signature Algorithm** | ECDSA | Schnorr (BIP340) |
| **Address Type** | P2WSH (bc1q...) | P2TR Taproot (bc1p...) |
| **Compatibility** | Older standard | Modern (Bitcoin Core 22.0+) |

## When to Use MuSig2

- ✅ **New multisig setups** - Best privacy and efficiency
- ✅ **High-volume operations** - Significant fee savings over time
- ✅ **Privacy-focused applications** - Transactions look like single-sig
- ✅ **Modern infrastructure** - Requires Bitcoin Core 22.0+ and Taproot support
- ⚠️ **Legacy multisig** - Traditional P2WSH still supported for backward compatibility

## How MuSig2 Works (Two-Round Protocol)

```
Round 1: Nonce Generation (Parallel)
┌─────────────────────────────────────────────────────┐
│  Keygen Wallet → Generate Nonce 1                   │
│  Sign Wallet 1 → Generate Nonce 2  (can run in     │
│  Sign Wallet 2 → Generate Nonce 3   parallel)      │
└─────────────────────────────────────────────────────┘
                        ↓
            Exchange nonces via PSBT files
                        ↓
Round 2: Signing (Sequential)
┌─────────────────────────────────────────────────────┐
│  Keygen Wallet → Create Partial Signature 1         │
│  Sign Wallet 1 → Create Partial Signature 2         │
│  Sign Wallet 2 → Create Partial Signature 3         │
└─────────────────────────────────────────────────────┘
                        ↓
            Collect partial signatures
                        ↓
Aggregation (Watch Wallet)
┌─────────────────────────────────────────────────────┐
│  Watch Wallet → Aggregate Partial Signatures        │
│              → Verify Final Signature               │
│              → Broadcast Transaction                │
└─────────────────────────────────────────────────────┘
```

**Key Security Feature:**

- Each wallet generates a **nonce** (random value) in Round 1
- Nonces must be **unique per transaction** and **never reused**
- Reusing nonces can leak private keys - this is critical!

---

# Wallet Type

This is explained for BTC/BCH for now.
There are mainly 3 wallets separately and these wallets are expected to be installed in each different devices.

## 1.Watch only wallet

- Only this wallet run online to access to BTC/BCH Nodes.
- Only pubkey address is stored. Private key is NOT stored for security reason. That's why this is called `watch only wallet`.
- Major functionalities are
  - creating unsigned transaction
  - sending signed transaction
  - monitoring transaction status.

## 2.Keygen wallet as cold wallet

- Key management functionalities for accounts.
- This wallet is expected to work offline.
- Major functionalities are
  - generating seed for accounts
  - generating keys based on `HD Wallet`
  - generating multisig addressed according to account setting
  - exporting pubkey addresses as csv file which is imported from `Watch only wallet`
  - signing on unsigned transaction as first sign. However, multisig addresses could not be completed by only this wallet.

## 3.Sign wallet as cold wallet (Auth wallet)

- The internal authorization operators would use this wallet to sign on unsigned transaction for multisig addresses.
- Each of operators would be given own authorization account and Sing wallet apps.
- This wallet is expected to work offline.
- Major functionalities are
  - generating seed for accounts for own auth account
  - generating keys based on `HD Wallet` for own auth account
  - exporting full-pubkey addresses as csv file which is imported from `Keygen wallet` to generate multisig address
  - signing on unsigned transaction as second or more signs for multisig addresses.

# Workflow diagram

## BTC

### 1. Generate keys

![generate keys](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/0_key%20generation%20diagram.png?raw=true)

### 2. Create unsigned transaction, Sign on unsigned tx, Send signed tx for non-multisig address

![create tx](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/1_Handle%20transactions%20for%20non-multisig%20address.png?raw=true)

### 3. Create unsigned transaction, Sign on unsigned tx, Send signed tx for multisig address

![create tx for multisig](https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/2_Handle%20transactions%20for%20multisig%20address.png?raw=true)

# Wallet Architecture

## Three Wallet Types

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Watch Wallet  │     │  Keygen Wallet  │     │   Sign Wallet   │
│    (Online)     │     │   (Offline)     │     │   (Offline)     │
├─────────────────┤     ├─────────────────┤     ├─────────────────┤
│ • Monitor txs   │     │ • Generate keys │     │ • Auth signing  │
│ • Create unsig  │     │ • Create multis │     │ • Second+ sign  │
│ • Send signed   │     │ • First sign    │     │ • Export pubkey │
│ • Import pubkey │     │ • Export pubkey │     │                 │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                       │                       │
        │    CSV/File Export    │    CSV/File Export    │
        └───────────────────────┴───────────────────────┘
```

## Security Model

1. **Keygen Wallet** (Offline): Generates HD wallet seeds and keys. Never connects to network.
2. **Sign Wallet** (Offline): Provides authorization signatures. Each operator has own instance.
3. **Watch Wallet** (Online): Only stores public keys. Cannot sign transactions.
