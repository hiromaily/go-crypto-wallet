## Multisig & MuSig2

### Traditional Multisig (P2SH/P2WSH)

**M-of-N Redeem Script:**

```
<M> <PubKey1> <PubKey2> ... <PubKeyN> <N> OP_CHECKMULTISIG
```

### MuSig2 (Schnorr Signature Aggregation)

MuSig2 enables N-of-N multisig that appears as single-sig on-chain.

**Benefits:**

- 30-50% smaller transactions
- Maximum privacy (looks like single-sig)
- Lower fees
- No on-chain multisig indicator

**Critical Security: NONCE MANAGEMENT**

- **NEVER reuse nonces** - reusing leaks private key
- Generate fresh nonces for every transaction
- Delete nonces after signing

See [musig2/](../../../../docs/chains/btc/musig2/README.md) for detailed documentation.

---
