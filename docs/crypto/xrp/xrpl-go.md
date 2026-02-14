# About [xrpl-go](https://github.com/xrpscan/xrpl-go)

[xrpl-go](https://github.com/xrpscan/xrpl-go) does *not* provide full offline transaction building/signing + submission without network access. It *can* do some offline work (keypairs, serialization, signing), but **most useful functionality still assumes a connection to an XRPL node (via WebSocket or RPC)** to query ledger state and submit transactions.

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
