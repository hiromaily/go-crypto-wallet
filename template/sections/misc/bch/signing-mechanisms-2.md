## Signing Mechanisms

### ECDSA Signatures

BCH uses ECDSA (secp256k1) for all signatures. Unlike Bitcoin, BCH does NOT support Schnorr signatures.

**Signature Format (DER Encoded):**

```
0x30 [total-length]
  0x02 [r-length] [r]
  0x02 [s-length] [s]
[sighash-type]  (includes FORKID: 0x41)
```

**Size:** 71-73 bytes (variable due to DER encoding)

### BCH Signature Digest Algorithm

Post-fork BCH uses a modified signature digest algorithm (similar to BIP143) for replay protection:

```
Signature Digest:
1. nVersion (4 bytes)
2. hashPrevouts (32 bytes)
3. hashSequence (32 bytes)
4. outpoint (36 bytes)
5. scriptCode (variable)
6. value (8 bytes)
7. nSequence (4 bytes)
8. hashOutputs (32 bytes)
9. nLocktime (4 bytes)
10. sighash type (4 bytes) - includes FORKID
```

---
