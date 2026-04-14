### Schema Design

#### Watch Schema (`watch`)

**Purpose**: Manages online wallet operations including address tracking, transaction monitoring, and payment requests.

**Tables**:

| Table | Description |
|-------|-------------|
| `address` | Wallet addresses for all coin/account types (btc, bch, eth, xrp, hyt) |
| `btc_tx` | Bitcoin/BCH transaction records |
| `btc_tx_input` | Bitcoin transaction inputs (UTXOs) |
| `btc_tx_output` | Bitcoin transaction outputs |
| `tx` | Generic transaction records (ETH, XRP, HYT) |
| `eth_detail_tx` | Ethereum transaction details |
| `xrp_detail_tx` | XRP transaction details |
| `xrp_pending_multisig` | Pending XRP multi-signature transactions awaiting signatures |
| `xrp_multisig_signature` | Collected signatures for XRP multi-signature transactions |
| `payment_request` | Payment request queue |

**Access Pattern**: High read/write - monitors blockchain, creates transactions

#### Keygen Schema (`keygen`)

**Purpose**: Stores key generation data for offline key generation wallet.

**Tables**:

| Table | Description |
|-------|-------------|
| `seed` | Encrypted seed phrases (all coins) |
| `btc_account_key` | BTC/BCH keys with multiple address formats (P2PKH, P2SH-SegWit, Bech32, Taproot) |
| `eth_account_key` | Ethereum keys (address, public key, private key) |
| `xrp_account_key` | XRP-specific account keys |
| `auth_fullpubkey` | Full public keys for multisig authentication (with BIP32 extended key support) |
| `xrp_signer_list` | XRP SignerList configuration for multi-signature accounts |
| `xrp_signer_entry` | Individual signer entries within an XRP signer list |
| `xrp_regular_key` | XRP regular key assignments for enhanced security |
| `musig2_nonces` | MuSig2 nonce commitments (shared with sign schema) |

**Access Pattern**: Write-heavy during key generation, read-only during export

**Security**: This schema contains sensitive key material - should be in offline/cold storage in production

#### Sign Schema (`sign` / `sign2`)

**Purpose**: Stores signing wallet data for offline transaction signing. `sign2` uses the same schema for a second signing wallet instance.

**Tables**:

| Table | Description |
|-------|-------------|
| `seed` | Encrypted seed phrases for signing wallet (BTC/BCH only) |
| `auth_account_key` | Authentication account keys for signing (with multiple address formats) |
| `musig2_nonces` | MuSig2 nonce commitments (shared with keygen schema) |

**Access Pattern**: Read-only during signing operations

**Security**: This schema contains sensitive signing keys - should be in offline/cold storage in production

#### Key Type Support

The `btc_account_key` and `auth_account_key` tables support multiple key types:

| Key Type | BIP Standard | Address Format |
|----------|-------------|----------------|
| `bip44` | BIP-44 | P2PKH (legacy) |
| `bip49` | BIP-49 | P2SH-SegWit |
| `bip84` | BIP-84 | Bech32 (native SegWit) |
| `bip86` | BIP-86 | Taproot |
| `musig2` | MuSig2 | Multi-signature |

#### Type Differences Across Dialects

| Feature | PostgreSQL | MySQL | SQLite |
|---------|-----------|-------|--------|
| Auto-increment | `identity` | `auto_increment` | `INTEGER PRIMARY KEY` |
| Numeric | `numeric(26,10)` | `decimal(26,10)` | `TEXT` |
| Enums | Named types | Inline `enum()` | `CHECK` constraints |
| Timestamps | `timestamptz` | `datetime` | `TEXT` |
| Boolean | `boolean` | `bool` (tinyint) | `INTEGER` |
| Binary | `bytea` | `blob` | `BLOB` |

---
