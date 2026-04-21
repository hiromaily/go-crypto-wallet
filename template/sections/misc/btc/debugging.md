## Debugging

### Debugging PSBT Issues

#### Enable Debug Logging

```go
// In code
logger.SetLevel(logger.DebugLevel)

// Or via config
[logger]
level = "debug"
```

#### Inspect PSBT with Bitcoin Core

```bash
# Decode PSBT
bitcoin-cli decodepsbt "$(cat transaction.psbt)"

# Analyze PSBT
bitcoin-cli analyzepsbt "$(cat transaction.psbt)"
```

**Output shows:**

- Inputs and their metadata
- Outputs
- Current signatures
- Missing signatures
- Fee estimation

#### Common Issues and Solutions

**Issue 1: "PSBT missing witness UTXO"**

**Solution:** Ensure witness UTXO data is added in CreatePSBT:

```go
packet.Inputs[i].WitnessUtxo = &wire.TxOut{
    Value:    prevTxs[i].Amount,
    PkScript: prevTxs[i].ScriptPubKey,
}
```

**Issue 2: "Invalid signature"**

**Solution:** Verify correct signature algorithm:

- Taproot → Schnorr (BIP340)
- Legacy/SegWit → ECDSA

**Issue 3: "PSBT not finalizing"**

**Solution:** Check all required signatures present:

```go
isComplete, err := btc.IsPSBTComplete(psbtBase64)
if !isComplete {
    // Add missing signatures
}
```

### Debugging Tools

#### PSBT Inspector Script

```bash
#!/bin/bash
# inspect_psbt.sh

PSBT_FILE="$1"

if [ ! -f "$PSBT_FILE" ]; then
    echo "Usage: $0 <psbt_file>"
    exit 1
fi

echo "=== PSBT Analysis ==="
bitcoin-cli analyzepsbt "$(cat "$PSBT_FILE")"

echo ""
echo "=== PSBT Decode ==="
bitcoin-cli decodepsbt "$(cat "$PSBT_FILE")"
```

---
