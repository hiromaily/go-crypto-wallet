## Address Types & Key Derivation

### Address Type

This system supports **EOA (Externally Owned Account)** addresses for key generation and single-sig flows. **Safe (Gnosis Safe v1.4.1) smart contract wallets** are also supported for multisig flows — see [multisig.md](../../../../docs/chains/eth/multisig.md).

EOA addresses:

- Are 20 bytes (40 hex chars) with `0x` prefix
- Example: `0x742d35Cc6634C0532925a3b8BC9e7595f0bEb123`
- Validated via `common.IsHexAddress(addr)`

### HD Wallet Derivation Path

**Standard:** BIP44 with Ethereum coin type

```
m / 44' / 60' / account' / 0 / address_index
```

| Component | Value | Notes |
|-----------|-------|-------|
| **Purpose** | 44' | BIP44 |
| **Coin Type** | 60' | SLIP-0044 Ethereum |
| **Account** | See table below | Non-hardened for deposit/payment/stored |
| **Change** | 0 | External chain only (no internal change addresses) |
| **Index** | 0..N | Sequential address generation |

### Account Types and Derivation Paths

| Account | Index | Hardened | Path | Use |
|---------|-------|----------|------|-----|
| **deposit** | 0 | No | `m/44'/60'/0'/0/i` | Aggregate client funds |
| **payment** | 1 | No | `m/44'/60'/1'/0/i` | Outgoing payments |
| **stored** | 2 | No | `m/44'/60'/2'/0/i` | Cold storage |
| **auth1** | 11 | Yes | `m/44'/60'/11''/0/i` | Authentication key 1 |
| **auth2** | 12 | Yes | `m/44'/60'/12''/0/i` | Authentication key 2 |

> **Note:** Deposit/payment/stored use non-hardened derivation so the Watch Wallet can derive child addresses from the extended public key (xpub) without private keys. Auth accounts use hardened derivation for enhanced security.

---
