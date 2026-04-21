## Signing Mechanisms

### ECDSA Signatures (Legacy/SegWit)

Used for P2PKH, P2SH, P2WPKH, and P2WSH transactions.

**Sighash Types:**

| Type | Value | Description |
|------|-------|-------------|
| SIGHASH_ALL | 0x01 | Sign all inputs and outputs (default) |
| SIGHASH_NONE | 0x02 | Sign all inputs, no outputs |
| SIGHASH_SINGLE | 0x03 | Sign all inputs, matching output only |
| SIGHASH_ANYONECANPAY | 0x80 | Modifier: sign only current input |

### Schnorr Signatures (Taproot)

Used for P2TR transactions. Introduced with Taproot (BIP340).

**Advantages:**

- Fixed 64-byte size (vs variable ECDSA)
- Linear: enables signature aggregation (MuSig2)
- Provably secure under standard assumptions
- Batch verification is faster

See [taproot/user-guide.md](../../../../docs/chains/btc/taproot/user-guide.md) for details.

---
