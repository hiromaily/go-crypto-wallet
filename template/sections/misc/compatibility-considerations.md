## Compatibility Considerations

### Database Schema Compatibility

#### Backward Compatibility

The schema maintains backward compatibility:

```sql
-- Old queries still work
SELECT * FROM account_key WHERE p2wsh_address IS NOT NULL;

-- New queries for MuSig2
SELECT * FROM account_key WHERE taproot_address IS NOT NULL;

-- Both types
SELECT * FROM account_key WHERE
    p2wsh_address IS NOT NULL OR taproot_address IS NOT NULL;
```

#### Schema Evolution

| Version | P2WSH Support | Taproot Support | MuSig2 Support |
|---------|---------------|-----------------|----------------|
| v1.0    | ✅ p2wsh_address | ❌ | ❌ |
| v2.0    | ✅ p2wsh_address | ✅ taproot_address | ❌ |
| v3.0    | ✅ p2wsh_address | ✅ taproot_address | ✅ (via taproot_address + musig2_nonces table) |

### Wallet Configuration Compatibility

#### Keygen Wallet

```toml
# v2.x configuration (legacy)
[multisig]
require_num = 2
pubkey_num = 3

# v3.x configuration (MuSig2)
[multisig]
require_num = 2
pubkey_num = 3
use_musig2 = false  # Default: false (backward compat)

[taproot]
enabled = true      # Required for MuSig2
```

#### Sign Wallet

Same configuration structure as Keygen wallet.

#### Watch Wallet

```toml
# Additional Bitcoin Core settings for Taproot
[bitcoin]
host = "127.0.0.1:8332"
# Minimum version: 22.0 for Taproot support
# Recommended: 25.0+
```

### External System Compatibility

#### Block Explorers

- **P2WSH**: Supported by all explorers
- **MuSig2 (P2TR)**: Supported by explorers with Taproot support
  - ✅ blockstream.info
  - ✅ mempool.space
  - ✅ blockchain.com (since late 2021)
  - ⚠️ Older explorers may not recognize bc1p... addresses

#### Wallets and Exchanges

- **Sending to MuSig2 addresses**: Most modern wallets/exchanges support sending to bc1p... addresses (as of 2023)
- **Receiving from MuSig2 addresses**: All wallets/exchanges accept from bc1p... addresses (looks like normal payment)

**Compatibility Matrix**:

| Service Type | P2WSH (bc1q...) | MuSig2 (bc1p...) | Notes |
|--------------|-----------------|------------------|-------|
| **Bitcoin Core 22.0+** | ✅ | ✅ | Full support |
| **Bitcoin Core 0.21.x** | ✅ | ❌ | No Taproot |
| **Hardware Wallets** | ✅ | ✅ (most) | Ledger: Firmware 2.1.0+, Trezor: Firmware 2.4.2+ |
| **Major Exchanges** | ✅ | ✅ (most) | Verify with specific exchange |
| **Mobile Wallets** | ✅ | ✅ (varies) | Check wallet documentation |

#### API Compatibility

If you expose APIs for wallet integration:

```javascript
// Old API (P2WSH)
POST /api/v1/address/create
{
  "account": "deposit",
  "type": "p2wsh"  // explicit type
}

// New API (MuSig2)
POST /api/v1/address/create
{
  "account": "deposit",
  "type": "musig2"  // or "taproot"
}

// Backward-compatible default
POST /api/v1/address/create
{
  "account": "deposit"
  // type: "p2wsh" by default for backward compat
}
```

---
