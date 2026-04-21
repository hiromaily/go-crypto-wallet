## PSBT Infrastructure

### Bitcoin API Layer

Location: `internal/infrastructure/api/btc/btc/psbt.go`

#### Core PSBT Methods

**1. CreatePSBT**

```go
// CreatePSBT creates a PSBT from a wire.MsgTx and previous transaction data
func (b *Bitcoin) CreatePSBT(
    msgTx *wire.MsgTx,
    prevTxs []PrevTx,
) (string, error)
```

**Purpose:** Creates an unsigned PSBT from transaction inputs and outputs.

**Implementation:**

```go
// Create PSBT packet
packet, err := psbt.NewFromUnsignedTx(msgTx)

// Add witness UTXO information for SegWit/Taproot inputs
for i, input := range msgTx.TxIn {
    prevOut := prevTxs[i]
    packet.Inputs[i].WitnessUtxo = &wire.TxOut{
        Value:    prevOut.Amount,
        PkScript: prevOut.ScriptPubKey,
    }
    // Add additional metadata (derivation paths, etc.)
}

// Serialize to base64
return packet.B64Encode()
```

**2. SignPSBTWithKey**

```go
// SignPSBTWithKey signs a PSBT with provided private keys (offline)
func (b *Bitcoin) SignPSBTWithKey(
    psbtBase64 string,
    wifs []string,
) (string, bool, error)
```

**Purpose:** Signs PSBT inputs with provided WIF private keys (offline signing).

**Implementation:**

```go
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

// Parse private keys
privKeys := parseWIFs(wifs)

// Sign each input
for i := range packet.Inputs {
    // Determine signature type based on input type
    if isTaprootInput(packet.Inputs[i]) {
        // Schnorr signature (BIP340)
        sig, err := schnorr.Sign(privKey, sigHash)
    } else {
        // ECDSA signature
        sig, err := ecdsa.Sign(privKey, sigHash)
    }

    // Add signature to PSBT
    packet.Inputs[i].PartialSigs = append(
        packet.Inputs[i].PartialSigs,
        psbt.PartialSig{PubKey: pubKey, Signature: sig},
    )
}

// Check if fully signed
isSigned := isComplete(packet)

return packet.B64Encode(), isSigned, nil
```

**3. FinalizePSBT**

```go
// FinalizePSBT finalizes a fully signed PSBT
func (b *Bitcoin) FinalizePSBT(psbtBase64 string) (string, error)
```

**Purpose:** Combines signatures into final scriptSig/witness for broadcasting.

**Implementation:**

```go
// Decode base64 PSBT
psbtBytes, err := base64.StdEncoding.DecodeString(psbtBase64)
if err != nil {
    return "", fmt.Errorf("failed to decode PSBT: %w", err)
}

// Parse PSBT
packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
if err != nil {
    return "", fmt.Errorf("failed to parse PSBT: %w", err)
}

// Finalize each input
for i := range packet.Inputs {
    err := psbt.Finalize(packet, i)
    if err != nil {
        return "", fmt.Errorf("failed to finalize input %d: %w", i, err)
    }
}

return packet.B64Encode(), nil
```

**4. ExtractTransaction**

```go
// ExtractTransaction extracts the final transaction from a finalized PSBT
func (b *Bitcoin) ExtractTransaction(psbtBase64 string) (*wire.MsgTx, error)
```

**Purpose:** Extracts the final, broadcastable transaction from PSBT.

**Implementation:**

```go
// Decode base64 PSBT
psbtBytes, err := base64.StdEncoding.DecodeString(psbtBase64)
if err != nil {
    return nil, fmt.Errorf("failed to decode PSBT: %w", err)
}

// Parse PSBT
packet, err := psbt.NewFromRawBytes(bytes.NewReader(psbtBytes), false)
if err != nil {
    return nil, fmt.Errorf("failed to parse PSBT: %w", err)
}

// Extract transaction
tx, err := psbt.Extract(packet)
if err != nil {
    return nil, fmt.Errorf("failed to extract transaction: %w", err)
}

return tx, nil
```

#### Helper Methods

**5. IsPSBTComplete**

```go
// IsPSBTComplete checks if a PSBT has all required signatures
func (b *Bitcoin) IsPSBTComplete(psbtBase64 string) (bool, error)
```

**6. ParsePSBT**

```go
// ParsePSBT parses a PSBT and returns structured data
func (b *Bitcoin) ParsePSBT(psbtBase64 string) (*ParsedPSBT, error)
```

**7. ValidatePSBT**

```go
// ValidatePSBT validates PSBT format and structure
func (b *Bitcoin) ValidatePSBT(psbtBase64 string) error
```

### File Storage Layer

Location: `internal/infrastructure/storage/file/transaction.go`

#### PSBT File Operations

**1. WritePSBTFile**

```go
// WritePSBTFile writes a PSBT to a file with .psbt extension
func (r *TransactionFileRepository) WritePSBTFile(
    path string,
    psbtBase64 string,
) (string, error)
```

**Implementation:**

```go
// Validate PSBT format
if !isValidBase64(psbtBase64) {
    return "", errors.New("invalid PSBT base64 encoding")
}

// Add .psbt extension if missing
if !strings.HasSuffix(path, ".psbt") {
    path += ".psbt"
}

// Create parent directory if needed
os.MkdirAll(filepath.Dir(path), 0755)

// Write PSBT to file
err := os.WriteFile(path, []byte(psbtBase64), 0644)

return path, err
```

**2. ReadPSBTFile**

```go
// ReadPSBTFile reads a PSBT from a file
func (r *TransactionFileRepository) ReadPSBTFile(path string) (string, error)
```

**Implementation:**

```go
// Validate extension
if !strings.HasSuffix(strings.ToLower(path), ".psbt") {
    return "", fmt.Errorf("invalid PSBT file extension: %s", path)
}

// Security: Prevent path traversal
cleanPath := filepath.Clean(path)
if r.filePath != "" && !strings.HasPrefix(cleanPath, r.filePath) {
    return "", fmt.Errorf("path traversal attempt detected: %s", path)
}

// Read file
data, err := os.ReadFile(cleanPath)
if err != nil {
    return "", fmt.Errorf("failed to read PSBT file: %w", err)
}

return string(data), nil
```

---
