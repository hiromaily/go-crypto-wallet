## Transaction Architecture

### UTXO Model

Bitcoin uses the Unspent Transaction Output (UTXO) model:

```
UTXO = {
    txid:      32-byte transaction hash
    vout:      output index (uint32)
    value:     satoshi amount (int64)
    scriptPubKey: locking script
}
```

**Key Concepts:**

- Each transaction consumes UTXOs (inputs) and creates new UTXOs (outputs)
- Total inputs must equal outputs + transaction fee
- UTXOs can only be spent once (double-spend protection)

### Transaction Weight & Virtual Size

SegWit introduced weight units for fee calculation:

```
Weight = (Non-witness data × 4) + Witness data
Virtual Size (vBytes) = Weight ÷ 4

Fee = Virtual Size × Fee Rate (sat/vB)
```

**Typical Sizes:**

| Transaction Type | Weight | vBytes | Fee @ 10 sat/vB |
|------------------|--------|--------|-----------------|
| P2PKH (1-in, 2-out) | ~680 | ~170 | ~1,700 sats |
| P2WPKH (1-in, 2-out) | ~440 | ~110 | ~1,100 sats |
| P2TR (1-in, 2-out) | ~396 | ~99 | ~990 sats |
| 2-of-3 Multisig (P2WSH) | ~1,100 | ~275 | ~2,750 sats |
| 2-of-3 MuSig2 (P2TR) | ~560 | ~140 | ~1,400 sats |

**Reference:**

- [BIP141 - SegWit](https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki)
- [BIP144 - SegWit Peer Services](https://github.com/bitcoin/bips/blob/master/bip-0144.mediawiki)

---
