## Overview

### What Changed?

The go-crypto-wallet system has migrated from a custom CSV-based transaction format to the industry-standard PSBT format (BIP 174).

**Key Changes:**

| Aspect | Before (CSV) | After (PSBT) |
|--------|-------------|--------------|
| **File Format** | Custom CSV | Base64 PSBT |
| **Extension** | Various | `.psbt` |
| **Encoding** | Plain text | Base64 |
| **Standard** | Custom | BIP 174 |
| **Compatibility** | go-crypto-wallet only | Bitcoin Core, hardware wallets |
| **Validation** | Basic | Comprehensive |
| **Metadata** | Limited | Rich (UTXO info, scripts) |

### Why Migrate?

**Benefits of PSBT:**

1. **Standardization**
   - Industry-standard format (BIP 174)
   - Compatible with Bitcoin Core, Electrum, hardware wallets
   - Better interoperability

2. **Security**
   - Structured data format prevents parsing errors
   - Built-in validation
   - Better error detection

3. **Functionality**
   - Richer transaction metadata
   - Support for complex scripts
   - Better debugging capabilities

4. **Future-Proofing**
   - Foundation for hardware wallet integration
   - Enables Taproot multisig features
   - Supports MuSig2 (future)

### Migration Scope

**Affected Components:**

- ✅ Watch Wallet - Transaction creation and broadcasting
- ✅ Keygen Wallet - Transaction signing (first signature)
- ✅ Sign Wallet - Transaction signing (second signature)
- ✅ Transaction file storage
- ✅ Transaction file naming convention

**Not Affected:**

- ❌ Database schema (no changes)
- ❌ Key generation process
- ❌ Address generation
- ❌ Configuration files
- ❌ Command-line interface (same commands)

---
