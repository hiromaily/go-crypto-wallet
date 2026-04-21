## 2. Recommended Architecture: Hybrid Approach

### 2.1 Approach Comparison

| Approach | Watch Wallet | Keygen Wallet | Sign Wallet | Pros | Cons |
|----------|--------------|---------------|-------------|------|------|
| **btcd Only** | btcd package | btcd package | btcd package | Consistent API, offline-first | No RPC convenience |
| **RPC Only** | Bitcoin Core RPC | Bitcoin Core RPC | Bitcoin Core RPC | Simpler code | Requires online wallets |
| **Hybrid** ✅ | Bitcoin Core RPC | btcd package | btcd package | Best of both | Two API surfaces |

**Selected: Hybrid Approach**

### 2.2 Rationale

1. **Watch Wallet (Online)**
   - Uses Bitcoin Core RPC for convenience
   - `walletcreatefundedpsbt` handles input selection and fee calculation
   - `finalizepsbt` prepares transaction for broadcast
   - Already connected to Bitcoin Core

2. **Keygen/Sign Wallets (Offline)**
   - Use btcd PSBT package for offline signing
   - No network communication required
   - Maintains air-gapped security
   - Full control over signing process

3. **Compatibility**
   - PSBTs are standardized (BIP174)
   - Watch wallet can create PSBT, offline wallets can sign
   - Final PSBT can be finalized by Watch wallet

### 2.3 Wallet-Specific Implementation

#### Watch Wallet (Online)

```go
// Transaction Creation
func (w *WatchWallet) CreatePSBT(inputs, outputs, options) (string, error) {
    // Use Bitcoin Core RPC
    result, err := w.rpcClient.WalletCreateFundedPSBT(inputs, outputs, 0, options, true)
    if err != nil {
        return "", err
    }
    return result.PSBT, nil // Base64-encoded PSBT
}

// Transaction Finalization
func (w *WatchWallet) FinalizePSBT(psbtBase64 string) (*wire.MsgTx, error) {
    // Use Bitcoin Core RPC
    result, err := w.rpcClient.FinalizePSBT(psbtBase64, true)
    if err != nil {
        return nil, err
    }
    if !result.Complete {
        return nil, errors.New("PSBT not fully signed")
    }
    // Parse hex to wire.MsgTx
    return w.HexToMsgTx(result.Hex)
}
```

#### Keygen Wallet (Offline)

```go
// PSBT Signing (First Signature)
func (k *KeygenWallet) SignPSBT(psbtBase64 string, wifs []string) (string, bool, error) {
    // Decode base64 to bytes
    psbtBytes, err := base64.StdEncoding.DecodeString(psbtBase64)
    if err != nil {
        return "", false, fmt.Errorf("failed to decode base64: %w", err)
    }

    // Parse PSBT using btcd package
    packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
    if err != nil {
        return "", false, fmt.Errorf("invalid PSBT: %w", err)
    }

    // Create updater
    updater, err := psbt.NewUpdater(packet)
    if err != nil {
        return "", false, err
    }

    // Sign each input
    for i := range packet.UnsignedTx.TxIn {
        // Get private key for this input
        privKey, err := k.getPrivateKeyForInput(wifs, i)
        if err != nil {
            continue
        }

        // Create signature
        sig, err := k.createSignature(packet, i, privKey)
        if err != nil {
            return "", false, err
        }

        // Add partial signature
        err = updater.Sign(i, sig, privKey.PubKey().SerializeCompressed(), nil, nil)
        if err != nil {
            return "", false, err
        }
    }

    // Serialize back to base64
    var buf bytes.Buffer
    err = packet.Serialize(&buf)
    if err != nil {
        return "", false, err
    }

    psbtBase64Out := base64.StdEncoding.EncodeToString(buf.Bytes())
    isComplete := packet.IsComplete()

    return psbtBase64Out, isComplete, nil
}
```

#### Sign Wallet (Offline)

```go
// PSBT Signing (Second Signature)
func (s *SignWallet) SignPSBT(psbtBase64 string, authWIF string) (string, bool, error) {
    // Decode base64 to bytes
    psbtBytes, err := base64.StdEncoding.DecodeString(psbtBase64)
    if err != nil {
        return "", false, fmt.Errorf("failed to decode base64: %w", err)
    }

    // Parse PSBT using btcd package
    packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
    if err != nil {
        return "", false, fmt.Errorf("invalid PSBT: %w", err)
    }

    // Verify PSBT already has partial signatures
    if !hasPartialSignatures(packet) {
        return "", false, errors.New("PSBT must have at least one signature from Keygen wallet")
    }
}

// hasPartialSignatures checks if PSBT has at least one partial signature
func hasPartialSignatures(packet *psbt.Packet) bool {
    for _, input := range packet.Inputs {
        if len(input.PartialSigs) > 0 {
            return true
        }
    }
    return false
}

    // Create updater
    updater, err := psbt.NewUpdater(packet)
    if err != nil {
        return "", false, err
    }

    // Get auth private key
    privKey, err := btcutil.DecodeWIF(authWIF)
    if err != nil {
        return "", false, err
    }

    // Add second signature to each input
    for i := range packet.UnsignedTx.TxIn {
        sig, err := s.createSignature(packet, i, privKey)
        if err != nil {
            continue
        }

        err = updater.Sign(i, sig, privKey.PubKey().SerializeCompressed(), nil, nil)
        if err != nil {
            return "", false, err
        }
    }

    // Serialize
    var buf bytes.Buffer
    err = packet.Serialize(&buf)
    if err != nil {
        return "", false, err
    }

    psbtBase64Out := base64.StdEncoding.EncodeToString(buf.Bytes())
    isComplete := packet.IsComplete()

    return psbtBase64Out, isComplete, nil
}
```

---
