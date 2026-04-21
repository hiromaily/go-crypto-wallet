## 1. Library Support Validation

### 1.1 btcd PSBT Package

**Package**: `github.com/btcsuite/btcd/btcutil/psbt`
**Version**: v1.1.6 (via btcd v0.25.0)
**Status**: ✅ Full Support

#### Available Operations

| Operation | Function | Description |
|-----------|----------|-------------|
| **Create** | `New()`, `NewFromUnsignedTx()` | Create PSBT from scratch or unsigned tx |
| **Parse** | `NewFromRawBytes()` | Parse from base64 or binary |
| **Update** | `NewUpdater()` | Add inputs, outputs, scripts, metadata |
| **Sign** | `updater.Sign()` | Add partial signatures |
| **Finalize** | `Finalize()`, `MaybeFinalizeAll()` | Finalize inputs |
| **Extract** | `Extract()` | Extract signed transaction |
| **Serialize** | `Serialize()`, `B64Encode()` | Serialize to binary/base64 |
| **Validate** | `SanityCheck()`, `IsComplete()` | Validation functions |

#### Key Features

1. **Complete BIP174 Implementation**
   - All PSBT roles: Creator, Updater, Signer, Combiner, Finalizer, Extractor
   - Full compliance with BIP174 specification

2. **Address Type Support**
   - ✅ P2PKH (Legacy, 1...)
   - ✅ P2SH-SegWit (3...)
   - ✅ P2WPKH (Bech32, bc1q...)
   - ✅ P2TR (Taproot, bc1p...) with Schnorr signatures

3. **Offline Signing**
   - No network communication required
   - Perfect for air-gapped wallets (Keygen, Sign)
   - Private key operations done locally

4. **Multisig Support**
   - Partial signature handling
   - Signature aggregation
   - 2-of-2, 2-of-3, M-of-N multisig

#### Core Data Structures

```go
// Main PSBT container
type Packet struct {
    UnsignedTx *wire.MsgTx  // Unsigned transaction
    Inputs     []PInput     // Per-input metadata
    Outputs    []POutput    // Per-output metadata
    Unknowns   []Unknown    // Custom fields
}

// Per-input data (includes Taproot fields)
type PInput struct {
    NonWitnessUtxo         *wire.MsgTx          // Full prev tx (non-SegWit)
    WitnessUtxo            *wire.TxOut          // Prev output (SegWit)
    PartialSigs            []*PartialSig        // Signatures
    SighashType            txscript.SigHashType // Sighash type
    RedeemScript           []byte               // P2SH redeem script
    WitnessScript          []byte               // Witness script
    FinalScriptSig         []byte               // Final scriptSig
    FinalScriptWitness     []byte               // Final witness
    TaprootKeySpendSig     []byte               // Taproot key path sig
    TaprootScriptSpendSig  []*TaprootScriptSpendSig // Taproot script sigs
    TaprootInternalKey     []byte               // Taproot internal key
    TaprootMerkleRoot      []byte               // Taproot merkle root
    // ... BIP32 derivation fields
}
```

### 1.2 Bitcoin Core RPC PSBT Methods

**Minimum Version**: Bitcoin Core v0.17+ (PSBT support)
**Taproot Support**: Bitcoin Core v22.0+ (Schnorr signatures)
**Current Project**: Compatible (supports v17+)

#### Available RPC Methods

| Method | Purpose | Use Case |
|--------|---------|----------|
| `walletcreatefundedpsbt` | Create and fund PSBT | Watch wallet transaction creation |
| `walletprocesspsbt` | Sign PSBT with wallet keys | Watch wallet signing (optional) |
| `finalizepsbt` | Finalize completed PSBT | Watch wallet finalization |
| `combinepsbt` | Combine multiple PSBTs | Multisig signature combining |
| `converttopsbt` | Convert raw tx to PSBT | Legacy compatibility |

#### Method Details

**walletcreatefundedpsbt**

```json
// Parameters
{
  "inputs": [],            // Empty for auto-selection
  "outputs": {             // Recipient addresses and amounts
    "address": amount
  },
  "locktime": 0,
  "options": {
    "add_inputs": true,    // Auto-add inputs
    "changeAddress": "",   // Custom change address
    "feeRate": 0.0001,     // Custom fee rate
    "subtractFeeFromOutputs": []
  },
  "bip32derivs": true      // Include derivation paths
}

// Returns
{
  "psbt": "base64...",     // Base64-encoded PSBT
  "fee": 0.00001,          // Transaction fee (BTC)
  "changepos": 1           // Change output position
}
```

**walletprocesspsbt**

```json
// Parameters
{
  "psbt": "base64...",     // PSBT to sign
  "sign": true,            // Whether to sign
  "sighashtype": "ALL"     // Signature hash type
}

// Returns
{
  "psbt": "base64...",     // Signed PSBT
  "complete": true         // Fully signed?
}
```

**finalizepsbt**

```json
// Parameters
{
  "psbt": "base64...",     // PSBT to finalize
  "extract": true          // Extract transaction if complete
}

// Returns
{
  "hex": "...",            // Final transaction hex (if extracted)
  "complete": true         // All signatures present?
}
```

---
