## Transaction Architecture

### UTXO Model

BCH uses the same UTXO (Unspent Transaction Output) model as Bitcoin:

```
UTXO = {
    txid:      32-byte transaction hash
    vout:      output index (uint32)
    value:     satoshi amount (int64)
    scriptPubKey: locking script
}
```

### Transaction Structure

BCH uses legacy transaction format (no SegWit):

```
+--------------------+
| Version (4 bytes)  |  Usually 1 or 2
+--------------------+
| Input Count        |  VarInt
+--------------------+
| Inputs[]           |
|   - prevTxHash     |  32 bytes
|   - prevVout       |  4 bytes
|   - scriptSigLen   |  VarInt
|   - scriptSig      |  variable (contains signature)
|   - sequence       |  4 bytes
+--------------------+
| Output Count       |  VarInt
+--------------------+
| Outputs[]          |
|   - value          |  8 bytes (satoshis)
|   - scriptPubKeyLen|  VarInt
|   - scriptPubKey   |  variable
+--------------------+
| Locktime (4 bytes) |
+--------------------+
```

### Transaction Size

Since BCH doesn't support SegWit, transaction sizes are larger than BTC SegWit transactions:

| Transaction Type | Size | Fee @ 1 sat/byte |
|------------------|------|------------------|
| P2PKH (1-in, 2-out) | ~226 bytes | ~226 sats |
| P2PKH (2-in, 2-out) | ~374 bytes | ~374 sats |
| 2-of-3 Multisig (1-in, 2-out) | ~370-400 bytes | ~400 sats |

### Sighash Types

BCH uses the same sighash types as legacy Bitcoin, with an additional fork ID:

| Type | Value | Description |
|------|-------|-------------|
| SIGHASH_ALL | 0x01 | Sign all inputs and outputs |
| SIGHASH_NONE | 0x02 | Sign all inputs, no outputs |
| SIGHASH_SINGLE | 0x03 | Sign all inputs, matching output only |
| SIGHASH_ANYONECANPAY | 0x80 | Modifier: sign only current input |
| SIGHASH_FORKID | 0x40 | BCH fork ID flag (added post-fork) |

**BCH Sighash Example:**

```
SIGHASH_ALL + SIGHASH_FORKID = 0x41
```

### Replay Protection

BCH implements replay protection using SIGHASH_FORKID to prevent transactions from being valid on both BCH and BTC chains.

---
