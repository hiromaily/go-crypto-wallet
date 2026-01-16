# Bitcoin Address Types vs Address Formats - Reference Guide

**CRITICAL REFERENCE**: This document clarifies the distinction between address types and address formats to prevent common AI agent and developer confusion.

## Table of Contents

- [Core Concepts](#core-concepts)
- [Terminology Reference](#terminology-reference)
- [Configuration Values](#configuration-values)
- [Common Mistakes](#common-mistakes)
- [Quick Reference Table](#quick-reference-table)
- [Examples](#examples)

## Core Concepts

### Address Type (Semantic Concept)

**Address Type** refers to the **semantic purpose** and **script structure** of a Bitcoin address.

- Defines HOW bitcoins are locked (script type)
- Determines signature requirements
- Examples: P2PKH, P2SH, P2WPKH, P2WSH, P2TR

### Address Format (Encoding Method)

**Address Format** refers to the **encoding scheme** used to represent the address as a string.

- Defines HOW the address is displayed/encoded
- Does NOT change the underlying script type
- Examples: Base58, Bech32, Bech32m

### Critical Distinction

```
Address Type ≠ Address Format

P2TR (address type) USES bech32m (address format)
```

## Terminology Reference

### Address Types (Script Types)

| Address Type | Full Name | Description | BIP |
|--------------|-----------|-------------|-----|
| **P2PKH** | Pay-to-PubKey-Hash | Legacy single-sig | BIP13 |
| **P2SH** | Pay-to-Script-Hash | Legacy multisig wrapper | BIP16 |
| **P2WPKH** | Pay-to-Witness-PubKey-Hash | Native SegWit single-sig | BIP141, BIP84 |
| **P2WSH** | Pay-to-Witness-Script-Hash | Native SegWit multisig | BIP141, BIP84 |
| **P2SH-P2WPKH** | Nested SegWit single-sig | SegWit wrapped in P2SH | BIP141, BIP49 |
| **P2SH-P2WSH** | Nested SegWit multisig | SegWit multisig wrapped in P2SH | BIP141, BIP49 |
| **P2TR** | Pay-to-Taproot | Taproot (Schnorr) | BIP341, BIP86 |

### Address Formats (Encoding Schemes)

| Format | Description | Address Prefix | Used For | Checksum |
|--------|-------------|----------------|----------|----------|
| **Base58** | Legacy encoding | `1...` (P2PKH)<br>`3...` (P2SH) | P2PKH, P2SH | Base58Check |
| **Bech32** | SegWit encoding | `bc1q...` (mainnet)<br>`tb1q...` (testnet)<br>`bcrt1q...` (regtest) | P2WPKH, P2WSH | BCH error correction |
| **Bech32m** | Modified Bech32 | `bc1p...` (mainnet)<br>`tb1p...` (testnet)<br>`bcrt1p...` (regtest) | P2TR (Taproot) | Fixed checksum |

### Key Relationships

```
┌─────────────────────────────────────────────────────────────┐
│                      Address Types                          │
├──────────────┬──────────────┬───────────────────────────────┤
│   Legacy     │   SegWit     │         Taproot              │
├──────────────┼──────────────┼───────────────────────────────┤
│ P2PKH, P2SH  │ P2WPKH,      │          P2TR                │
│              │ P2WSH,       │                              │
│              │ P2SH-P2WPKH, │                              │
│              │ P2SH-P2WSH   │                              │
└──────────────┴──────────────┴───────────────────────────────┘
       │              │                    │
       ▼              ▼                    ▼
   Base58          Bech32              Bech32m
 (1... 3...)      (bc1q...)           (bc1p...)
```

## Configuration Values

### In `go-crypto-wallet` Configuration

#### `address_type` Field

This field specifies the **address type** (semantic concept), NOT the encoding format.

**Valid Values:**

| Config Value | Actual Address Type | Address Format | Address Prefix | BIP |
|--------------|---------------------|----------------|----------------|-----|
| `"legacy"` | P2PKH (single-sig)<br>P2SH (multisig) | Base58 | `1...` or `3...` | BIP44 |
| `"p2sh-segwit"` | P2SH-P2WPKH (single-sig)<br>P2SH-P2WSH (multisig) | Base58 | `3...` | BIP49 |
| `"bech32"` | P2WPKH (single-sig)<br>P2WSH (multisig) | Bech32 | `bc1q...` (62 chars for WSH) | BIP84 |
| `"taproot"` | P2TR | Bech32m | `bc1p...` | BIP86 |
| `"bech32m"` | ⚠️ **INVALID** - Use `"taproot"` instead | - | - | - |

#### ⚠️ Common Mistake

```yaml
# ❌ WRONG - "bech32m" is an encoding format, not an address type
address_type: "bech32m"

# ✅ CORRECT - "taproot" is the address type
address_type: "taproot"
```

#### `key_type` Field

This field specifies the BIP derivation path standard.

**Valid Values:**

| Config Value | BIP | Derivation Path | Address Type Compatibility |
|--------------|-----|-----------------|----------------------------|
| `"bip44"` | BIP44 | `m/44'/0'/account'/0/index` | `legacy` (P2PKH) |
| `"bip49"` | BIP49 | `m/49'/0'/account'/0/index` | `p2sh-segwit` (P2SH-P2WPKH) |
| `"bip84"` | BIP84 | `m/84'/0'/account'/0/index` | `bech32` (P2WPKH) |
| `"bip86"` | BIP86 | `m/86'/0'/account'/0/index` | `taproot` (P2TR) |
| `"musig2"` | Custom | Coin-level xpub export | `taproot` (P2TR with MuSig2) |

#### Automatic Derivation

In `go-crypto-wallet`, `key_type` is **automatically derived** from `address_type`:

```go
// internal/domain/address/types.go
func (a AddrType) ToKeyType() (key.KeyType, error) {
	switch a {
	case AddrTypeLegacy:
		return key.KeyTypeBIP44, nil
	case AddrTypeP2shSegwit:
		return key.KeyTypeBIP49, nil
	case AddrTypeBech32:
		return key.KeyTypeBIP84, nil
	case AddrTypeTaproot:
		return key.KeyTypeBIP86, nil
	case AddrTypeBCHCashAddr:
		// BCH uses BIP44 derivation path (same as legacy BTC)
		return key.KeyTypeBIP44, nil
	case AddrTypeETH:
		// Ethereum uses BIP44 derivation path
		return key.KeyTypeBIP44, nil
	default:
		return "", fmt.Errorf("unsupported address type for key derivation: %s", a)
	}
}
```

**Configuration Example:**

```yaml
# Only specify address_type - key_type is derived automatically
address_type: "taproot"  # Automatically uses BIP86 (key_type)

# You do NOT need to specify:
# key_type: "bip86"  # This is automatic!
```

## Common Mistakes

### Mistake 1: Confusing "bech32m" with "taproot"

**❌ WRONG:**
```yaml
address_type: "bech32m"  # This will fail - bech32m is a format, not a type
```

**✅ CORRECT:**
```yaml
address_type: "taproot"  # Taproot uses bech32m format automatically
```

**Explanation:**
- `"bech32m"` is an **encoding format** (how addresses are displayed)
- `"taproot"` is an **address type** (the script structure)
- Taproot addresses automatically use bech32m encoding
- The configuration expects address types, not encoding formats

### Mistake 2: Using "bech32" for Taproot

**❌ WRONG:**
```yaml
address_type: "bech32"  # This creates P2WPKH addresses (bc1q...), NOT Taproot!
```

**✅ CORRECT:**
```yaml
address_type: "taproot"  # This creates P2TR addresses (bc1p...)
```

**Explanation:**
- `"bech32"` refers to **Native SegWit** (P2WPKH/P2WSH) addresses
- Native SegWit uses the original **bech32** encoding
- Taproot uses the improved **bech32m** encoding
- These are different address types with different encodings

### Mistake 3: Mixing Address Prefixes

**Address Prefix Meanings:**

| Prefix | Network | Format | Address Type |
|--------|---------|--------|--------------|
| `bc1q...` | Mainnet | Bech32 | P2WPKH (20-byte hash) or P2WSH (32-byte hash) |
| `bc1p...` | Mainnet | Bech32m | P2TR (Taproot) |
| `tb1q...` | Testnet | Bech32 | P2WPKH or P2WSH |
| `tb1p...` | Testnet | Bech32m | P2TR |
| `bcrt1q...` | Regtest | Bech32 | P2WPKH or P2WSH |
| `bcrt1p...` | Regtest | Bech32m | P2TR |

**❌ WRONG Interpretation:**
> "I need a `bc1p...` address, so I'll set `address_type: "bech32m"`"

**✅ CORRECT Interpretation:**
> "I need a Taproot address, so I'll set `address_type: "taproot"`, which produces `bc1p...` addresses automatically"

### Mistake 4: Not Understanding the Relationship

```
┌──────────────────────────────────────────────────────┐
│ Configuration:  address_type: "taproot"             │
│                                                      │
│        ↓ (determines)                               │
│                                                      │
│ Address Type:   P2TR (Pay-to-Taproot)              │
│                                                      │
│        ↓ (automatically uses)                       │
│                                                      │
│ Address Format: bech32m                             │
│                                                      │
│        ↓ (results in)                               │
│                                                      │
│ Address Prefix: bc1p... (mainnet)                  │
│                 tb1p... (testnet)                   │
│                 bcrt1p... (regtest)                 │
└──────────────────────────────────────────────────────┘
```

## Quick Reference Table

### Address Type → Format → Prefix

| Config: `address_type` | Address Type | Format | Mainnet | Testnet | Regtest | Length |
|------------------------|--------------|--------|---------|---------|---------|--------|
| `"legacy"` | P2PKH | Base58 | `1...` | `m.../n...` | `m.../n...` | 26-35 |
| `"legacy"` (multisig) | P2SH | Base58 | `3...` | `2...` | `2...` | 26-35 |
| `"p2sh-segwit"` | P2SH-P2WPKH | Base58 | `3...` | `2...` | `2...` | 26-35 |
| `"p2sh-segwit"` (multisig) | P2SH-P2WSH | Base58 | `3...` | `2...` | `2...` | 26-35 |
| `"bech32"` | P2WPKH | Bech32 | `bc1q...` | `tb1q...` | `bcrt1q...` | 42 |
| `"bech32"` (multisig) | P2WSH | Bech32 | `bc1q...` | `tb1q...` | `bcrt1q...` | 62 |
| `"taproot"` | P2TR | Bech32m | `bc1p...` | `tb1p...` | `bcrt1p...` | 62 |

### E2E Script Configuration

**Pattern 9 (P2TR Taproot Single-sig):**

```bash
# Environment variable override
export WALLET_ADDRESS_TYPE="taproot"  # ✅ CORRECT

# NOT this:
export WALLET_ADDRESS_TYPE="bech32m"  # ❌ WRONG - will fail
export WALLET_ADDRESS_TYPE="bc1p"     # ❌ WRONG - that's a prefix, not a type
```

## Examples

### Example 1: Generating Taproot Addresses

**Configuration:**

```yaml
# config/wallet/btc_keygen.yaml
address_type: "taproot"  # This is all you need!
# key_type is automatically set to "bip86"
```

**Result:**

```bash
$ keygen --coin btc create hdkey --account payment --count 1
# Generates address with derivation path: m/86'/0'/1'/0/0
# Address format: bcrt1p... (regtest) or bc1p... (mainnet)
# Address encoding: bech32m (automatic)
```

### Example 2: Bitcoin Core RPC

**Generating New Address:**

```bash
# Correct: Specify the address type
bitcoin-cli -regtest -rpcwallet=watch getnewaddress "" bech32m
# Returns: bcrt1p... (Taproot address)

# Also correct: Use label "taproot"
bitcoin-cli -regtest -rpcwallet=watch getnewaddress "mylabel" bech32m

# IMPORTANT: The "bech32m" parameter here refers to the desired format,
# which Bitcoin Core interprets as "generate a Taproot address"
```

### Example 3: Identifying Address Type from Prefix

```python
def identify_address_type(address: str) -> str:
    """Identify Bitcoin address type from prefix."""

    if address.startswith('1'):
        return "P2PKH (Legacy)"
    elif address.startswith(('3', '2')):
        return "P2SH (Legacy multisig or Nested SegWit)"
    elif address.startswith(('bc1q', 'tb1q', 'bcrt1q')):
        # Check length: 42 = P2WPKH, 62 = P2WSH
        if len(address) == 42:
            return "P2WPKH (Native SegWit)"
        elif len(address) == 62:
            return "P2WSH (Native SegWit Multisig)"
    elif address.startswith(('bc1p', 'tb1p', 'bcrt1p')):
        return "P2TR (Taproot)"

    return "Unknown"

# Examples
print(identify_address_type("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"))
# Output: "P2PKH (Legacy)"

print(identify_address_type("bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4"))
# Output: "P2WPKH (Native SegWit)"

print(identify_address_type("bc1p5d7rjq7g6rdk2yhzks9smlaqtedr4dekq08ge8ztwac72sfr9rusxg3297"))
# Output: "P2TR (Taproot)"
```

## AI Agent Guidelines

### When Reviewing Code or Configuration

1. **Never suggest `address_type: "bech32m"`**
   - ✅ Use `address_type: "taproot"` instead
   - Bech32m is the encoding, not the type

2. **Verify Environment Variable Names**
   - ✅ `WALLET_ADDRESS_TYPE="taproot"`
   - ❌ `WALLET_ADDRESS_TYPE="bech32m"`

3. **Check Address Prefix Expectations**
   - `bc1q...` (42 chars) = P2WPKH (Native SegWit single-sig)
   - `bc1q...` (62 chars) = P2WSH (Native SegWit multisig)
   - `bc1p...` (62 chars) = P2TR (Taproot)

4. **Bitcoin Core RPC Calls**
   - `getnewaddress "" bech32m` → Creates Taproot address (`bc1p...`)
   - `getnewaddress "" bech32` → Creates Native SegWit address (`bc1q...`)

### Common Confusion Patterns

| AI Agent Says... | Correct Interpretation |
|------------------|------------------------|
| "Set address format to bech32m" | "Set address type to taproot" |
| "Use bech32m addresses" | "Use Taproot addresses (which use bech32m encoding)" |
| "Configure bech32m" | "Configure taproot (which automatically uses bech32m)" |

## References

### BIPs (Bitcoin Improvement Proposals)

- **BIP13**: Address Format for pay-to-script-hash
- **BIP16**: Pay to Script Hash (P2SH)
- **BIP44**: Multi-Account Hierarchy for Deterministic Wallets (Legacy)
- **BIP49**: Derivation scheme for P2WPKH-nested-in-P2SH (Nested SegWit)
- **BIP84**: Derivation scheme for P2WPKH (Native SegWit)
- **BIP86**: Key Derivation for Single Key P2TR Outputs (Taproot)
- **BIP141**: Segregated Witness (SegWit)
- **BIP173**: Base32 address format for native v0-16 witness outputs (Bech32)
- **BIP350**: Bech32m format for v1+ witness addresses (Taproot)
- **BIP340**: Schnorr Signatures for secp256k1
- **BIP341**: Taproot: SegWit version 1 spending rules

### Internal Code References

- `internal/domain/address/types.go`: Address type definitions and conversions
- `internal/domain/key/types.go`: Key type definitions (BIP44/49/84/86)
- `pkg/address/`: Address encoding/decoding utilities

### Related Documentation

- [TAPROOT_GUIDE.md](./TAPROOT_GUIDE.md): Comprehensive Taproot usage guide
- [e2e_transaction_patterns.md](./e2e_transaction_patterns.md): E2E test patterns for all address types
- [btc_bch_technical_guide.md](./btc_bch_technical_guide.md): Bitcoin and Bitcoin Cash technical details

## Summary

### Remember These Key Points

1. **Taproot** is an address type (P2TR)
2. **Bech32m** is an encoding format
3. In configuration, use `address_type: "taproot"`, NOT `"bech32m"`
4. Taproot addresses automatically use bech32m encoding
5. Address prefixes tell you the type:
   - `bc1q...` = Native SegWit (Bech32)
   - `bc1p...` = Taproot (Bech32m)

### Quick Decision Guide

**Need to create Taproot addresses?**

```yaml
# Step 1: Configure
address_type: "taproot"  # That's it!

# Step 2: Generate
$ keygen create hdkey --account payment

# Step 3: Verify
$ keygen export address --account payment
# Should see: bcrt1p... addresses
```

---

**Last Updated**: 2026-01-16
**Version**: 1.0
**Maintainer**: go-crypto-wallet team
