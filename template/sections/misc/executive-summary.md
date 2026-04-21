## Executive Summary

This document outlines the technical approach for implementing PSBT (Partially Signed Bitcoin Transaction, BIP174) support in go-crypto-wallet. The research validates that **both btcd library and Bitcoin Core RPC have full PSBT support**, enabling a hybrid approach that maintains offline wallet security while leveraging online wallet capabilities.

### Key Findings

✅ **btcd v0.25.0** has comprehensive PSBT support (`github.com/btcsuite/btcd/btcutil/psbt`)
✅ **Bitcoin Core RPC** provides PSBT methods for online wallets
✅ **Offline signing** fully supported for Keygen and Sign wallets
✅ **All address types** supported (P2PKH, P2SH, P2WPKH, P2TR)
✅ **No blockers** identified for implementation

### Recommended Approach: **Hybrid**

- **Watch Wallet** (online): Bitcoin Core RPC PSBT methods
- **Keygen/Sign Wallets** (offline): btcd PSBT package
- **File Format**: Base64-encoded PSBT with `.psbt` extension
- **Migration**: Clean break from CSV to PSBT (no backward compatibility)

---
