## Overview

### What is PSBT?

PSBT (Partially Signed Bitcoin Transaction) is a Bitcoin standard (BIP 174) that provides a structured format for unsigned and partially signed transactions. It allows multiple parties to collaborate on signing a transaction without exposing private keys.

### Benefits Over Legacy CSV Format

| Feature | Legacy CSV | PSBT |
|---------|-----------|------|
| **Format** | Custom CSV | Standardized (BIP 174) |
| **Encoding** | Plain text | Base64 |
| **Compatibility** | go-crypto-wallet only | Bitcoin Core, hardware wallets, other tools |
| **Metadata** | Limited | Rich (UTXO info, derivation paths, etc.) |
| **Security** | Basic | Enhanced with structured validation |
| **Error Handling** | Limited | Comprehensive |
| **File Extension** | Various | `.psbt` |

### When to Use PSBT

- ✅ **All new transactions** - PSBT is the recommended format
- ✅ **Multisig transactions** - Better tracking of signatures
- ✅ **Offline signing** - Secure air-gapped wallet operations
- ✅ **Hardware wallet integration** - Future compatibility
- ⚠️ **Legacy transactions** - CSV format still supported for backward compatibility

---
