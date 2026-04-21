## Quick Reference

### Address Format Summary

| Type | Mainnet Prefix | Testnet Prefix | Example |
|------|----------------|----------------|---------|
| P2PKH (CashAddr) | `bitcoincash:q` | `bchtest:q` | `bitcoincash:qp3wjpa3tjlj042z2wv7hahsldgwhwy0rq9sywjpyy` |
| P2SH (CashAddr) | `bitcoincash:p` | `bchtest:p` | `bitcoincash:pr0662zpd7vr936d83f64u629v886aan7c77r3j5v5` |
| P2PKH (Legacy) | `1` | `m`/`n` | `1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2` |
| P2SH (Legacy) | `3` | `2` | `3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy` |

### Network Parameters

| Parameter | Mainnet | Testnet |
|-----------|---------|---------|
| Network Magic | `0xe8f3e1e3` | `0xf4f3e5f4` |
| P2P Port | 8333 | 18333 |
| RPC Port | 8332 | 18332 |
| Coin Type (SLIP44) | 145 | 1 |
| Address Prefix | `bitcoincash:` | `bchtest:` |

### Transaction Checklist

- [ ] Use CashAddr format for addresses
- [ ] Include SIGHASH_FORKID (0x40) in signatures
- [ ] Set appropriate fee (typically 1 sat/byte)
- [ ] Verify replay protection is active
- [ ] Wait for sufficient confirmations
