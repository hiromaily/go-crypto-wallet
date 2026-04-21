# About [xrpscan/xrpl-go](https://github.com/xrpscan/xrpl-go)

> **Note**: This document covers **xrpscan/xrpl-go** specifically. For the overall library selection strategy (why we use both xrpscan and Peersyst), see [library-selection.md](../../../../docs/chains/xrp/library-selection.md).

[xrpscan/xrpl-go](https://github.com/xrpscan/xrpl-go) does *not* provide full offline transaction building/signing + submission without network access. It *can* do some offline work (keypairs, serialization, signing), but **most useful functionality still assumes a connection to an XRPL node (via WebSocket or RPC)** to query ledger state and submit transactions.

**In go-crypto-wallet**: We use xrpscan/xrpl-go primarily for **transaction submission** and **ledger queries** (online operations on the watch wallet).

---

## What `xrpl-go` *does* support

### Offline-capable parts

* **Keypair/Wallet generation**
  * The library provides utilities for keypairs/wallets, including seed-based derivation and address handling.
* **Serialization / binary codec**
  * You can serialize/deserialize XRPL transactions offline if you have all required fields.
* **Signing transactions offline**
  * The library supports signing. For example, the `wallet` package can sign transactions if all inputs (fields) are provided.

These parts can be used without a network connection **as long as you already know all the input values** (e.g., Sequence, Fee).

---

## What `xrpl-go` *cannot* do offline alone

* **Querying current account state**
  * You must connect to an XRPL node to get the `Sequence`, account balances, ledger index, etc.
* **Fetching reliable network parameters**
  * Fee recommendation, validated ledger index, last ledger sequence windows require node access.
* **Submitting transactions**
  * Submission and confirmation (reliable submission) require hitting a rippled node (WebSocket or RPC).

Even the XRPL documentation uses the library **to connect to the XRP Ledger network** in its examples.

---

## Summary (2026-ready)

| Feature                                 | Offline (No Node) | Requires Node |
|------------------------------------------|-------------------|---------------|
| Key & address generation                 | ✅                | —             |
| Serialization / Sign                     | ✅                | —             |
| Query account / ledger state             | —                 | ✅            |
| Submit transaction                       | —                 | ✅            |
| Reliable submission (confirm validation) | —                 | ✅            |

So **`xrpl-go` does not replace the need to talk to XRPL nodes** — it is a library for interacting *with* XRPL, not a full offline SDK that builds every required field by itself

---

## How xrpscan/xrpl-go Fits in go-crypto-wallet

### Current Usage (2026)

| Component | Uses xrpscan/xrpl-go For |
|-----------|--------------------------|
| **Watch Wallet** | ✅ Account queries (`GetAccountInfo`, `GetBalance`) |
| **Watch Wallet** | ✅ Transaction submission (`SubmitTransaction`) |
| **Watch Wallet** | ✅ Validation waiting (`WaitValidation`) |
| **Keygen/Sign Wallets** | ❌ Not used (offline wallets) |

### What We DON'T Use xrpscan For

| Feature | Why Not | Alternative |
|---------|---------|-------------|
| Transaction Signing | Limited wallet API | ✅ Peersyst/xrpl-go (`wallet.Sign()`) |
| Multi-Signature | No `Multisign()` method | ✅ Peersyst/xrpl-go (`wallet.Multisign()`) |
| Offline Operations | Requires node connection | ✅ Peersyst/xrpl-go (zero network) |

### Specialized Roles

```
┌─────────────────────────────────────┐
│         Watch Wallet (Online)        │
│  ┌───────────────────────────────┐  │
│  │    xrpscan/xrpl-go v0.2.11    │  │
│  │  • Query account info          │  │
│  │  • Submit transactions         │  │
│  │  • Wait for validation         │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│    Keygen/Sign Wallets (Offline)    │
│  ┌───────────────────────────────┐  │
│  │   Peersyst/xrpl-go v0.1.15    │  │
│  │  • Derive wallets from seed    │  │
│  │  • Sign transactions (offline) │  │
│  │  • Multi-signature support     │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

**For the complete rationale**, see [library-selection.md](../../../../docs/chains/xrp/library-selection.md).
