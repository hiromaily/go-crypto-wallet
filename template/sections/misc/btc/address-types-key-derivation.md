## Address Types & Key Derivation

### Address Types Supported

| Type | BIP | Prefix (Mainnet) | Prefix (Testnet) | Description |
|------|-----|------------------|------------------|-------------|
| **P2PKH** | BIP44 | `1` | `m`/`n` | Legacy Pay-to-Public-Key-Hash |
| **P2SH** | BIP16 | `3` | `2` | Pay-to-Script-Hash |
| **P2SH-P2WPKH** | BIP49 | `3` | `2` | SegWit wrapped in P2SH |
| **P2WPKH** | BIP84 | `bc1q` | `tb1q` | Native SegWit |
| **P2WSH** | BIP141 | `bc1q` | `tb1q` | SegWit Script Hash |
| **P2TR** | BIP86 | `bc1p` | `tb1p` | Taproot (recommended) |

See [overview/address-types.md](../../../../docs/chains/btc/overview/address-types.md) for detailed comparison.

### HD Wallet Derivation Paths

| Standard | Path | Address Type |
|----------|------|--------------|
| **BIP44** | `m/44'/0'/account'/change/index` | P2PKH (Legacy) |
| **BIP49** | `m/49'/0'/account'/change/index` | P2SH-P2WPKH |
| **BIP84** | `m/84'/0'/account'/change/index` | P2WPKH (Native SegWit) |
| **BIP86** | `m/86'/0'/account'/change/index` | P2TR (Taproot) |

**Coin Types:**

- `0'` = Bitcoin Mainnet
- `1'` = Bitcoin Testnet/Signet

**References:**

- [BIP32 - HD Wallets](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)
- [BIP39 - Mnemonic Codes](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)
- [BIP43 - Purpose Field](https://github.com/bitcoin/bips/blob/master/bip-0043.mediawiki)
- [BIP44 - Multi-Account Hierarchy](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)

### ScriptPubKey Formats

```
P2PKH:      OP_DUP OP_HASH160 <20-byte pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
P2SH:       OP_HASH160 <20-byte scriptHash> OP_EQUAL
P2WPKH:     0x00 <20-byte pubKeyHash>
P2WSH:      0x00 <32-byte witnessScriptHash>
P2TR:       0x51 <32-byte x-only pubKey>
```

---
