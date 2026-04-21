## Adding New Features

### Adding Support for New Address Type

**Example: Adding P2TR Multisig (Script Path)**

#### Step 1: Update Address Generation

Location: `internal/infrastructure/wallet/key/btc/hdwallet.go`

```go
// Add Taproot multisig address generation
func (h *HDWallet) GenerateTaprootMultisigAddress(
    pubKeys []*btcec.PublicKey,
    threshold int,
) (string, error) {
    // Create Taproot script tree
    script := createMultisigScript(pubKeys, threshold)
    taprootKey := txscript.ComputeTaprootOutputKey(internalKey, script)

    // Generate address
    address, err := btcutil.NewAddressTaproot(
        schnorr.SerializePubKey(taprootKey),
        h.chainConfig,
    )

    return address.EncodeAddress(), nil
}
```

#### Step 2: Update PSBT Creation

Location: `internal/infrastructure/api/btc/btc/psbt.go`

```go
// Update CreatePSBT to include Taproot witness data
func (b *Bitcoin) CreatePSBT(msgTx *wire.MsgTx, prevTxs []PrevTx) (string, error) {
    packet, err := psbt.NewFromUnsignedTx(msgTx)

    for i, input := range msgTx.TxIn {
        packet.Inputs[i].WitnessUtxo = &wire.TxOut{
            Value:    prevTxs[i].Amount,
            PkScript: prevTxs[i].ScriptPubKey,
        }

        // Add Taproot-specific data
        if isTaprootOutput(prevTxs[i].ScriptPubKey) {
            packet.Inputs[i].TaprootInternalKey = prevTxs[i].InternalKey
            packet.Inputs[i].TaprootScriptTree = prevTxs[i].ScriptTree
        }
    }

    return packet.B64Encode(), nil
}
```

#### Step 3: Update Signing Logic

Location: `internal/infrastructure/api/btc/btc/psbt.go`

```go
// Update SignPSBTWithKey for Taproot script path
func (b *Bitcoin) SignPSBTWithKey(psbtBase64 string, wifs []string) (string, bool, error) {
    // Decode base64 PSBT
    psbtBytes, err := base64.StdEncoding.DecodeString(psbtBase64)
    if err != nil {
        return "", false, fmt.Errorf("failed to decode PSBT: %w", err)
    }

    // Parse PSBT
    packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
    if err != nil {
        return "", false, fmt.Errorf("failed to parse PSBT: %w", err)
    }

    for i := range packet.Inputs {
        if isTaprootScriptPath(packet.Inputs[i]) {
            // Schnorr signature for Taproot script path
            sig, err := signTaprootScriptPath(privKey, packet, i)
        } else if isTaprootKeyPath(packet.Inputs[i]) {
            // Schnorr signature for Taproot key path
            sig, err := schnorr.Sign(privKey, sigHash)
        } else {
            // ECDSA for legacy/SegWit
            sig, err := ecdsa.Sign(privKey, sigHash)
        }

        packet.Inputs[i].PartialSigs = append(
            packet.Inputs[i].PartialSigs,
            psbt.PartialSig{PubKey: pubKey, Signature: sig},
        )
    }

    return packet.B64Encode(), isComplete(packet), nil
}
```

#### Step 4: Add Tests

```go
func TestSignPSBTWithKey_TaprootScriptPath(t *testing.T) {
    // Create Taproot multisig PSBT
    psbt := createTestTaprootMultisigPSBT(t)

    // Sign with key
    signed, isComplete, err := btc.SignPSBTWithKey(psbt, []string{testWIF})

    // Verify
    assert.NoError(t, err)
    assert.True(t, isComplete)
    assert.NotEmpty(t, signed)
}
```

---
