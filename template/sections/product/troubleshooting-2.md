## Troubleshooting

### Common Errors

#### Error: "PSBT is not fully signed"

**Symptom:**

```
Error: PSBT is not fully signed - cannot finalize incomplete PSBT
```

**Cause:** Trying to broadcast a PSBT that doesn't have all required signatures.

**Solution:**

```bash
# Check signature status
./keygen sign --file payment_12_unsigned_1_*.psbt --dry-run

# Ensure all required signatures are collected:
# - Single-sig: 1 signature (Keygen)
# - 2-of-2 multisig: 2 signatures (Keygen + Sign)
```

#### Error: "Invalid PSBT format"

**Symptom:**

```
Error: failed to validate PSBT: invalid PSBT format
```

**Cause:** PSBT file is corrupted or invalid base64.

**Solution:**

```bash
# Verify file integrity
sha256sum payment_12_unsigned_0_*.psbt

# Re-create transaction if file is corrupted
./watch create payment --fee 0.0002

# Do not edit PSBT files manually
```

#### Error: "Transaction already broadcast"

**Symptom:**

```
Error: transaction already sent
Transaction ID: (empty)
```

**Cause:** Transaction was already broadcast to the network.

**Solution:**

- This is **not an error** - the transaction was successfully sent previously
- Check blockchain explorer to confirm transaction status
- No action needed

#### Error: "Failed to read PSBT file"

**Symptom:**

```
Error: failed to read PSBT file: invalid file extension (expected .psbt)
```

**Cause:** File doesn't have `.psbt` extension.

**Solution:**

```bash
# Ensure file has correct extension
mv payment_12_signed_2_1534744600000000002 payment_12_signed_2_1534744600000000002.psbt
```

#### Error: "Missing private key"

**Symptom:**

```
Error: failed to sign PSBT: private key not found for address
```

**Cause:** Wallet doesn't have the required private key.

**Solution:**

```bash
# Verify key exists in database
# Keygen wallet:
./keygen list-keys --account client

# Sign wallet:
./sign list-keys --auth auth1

# If key is missing, import or regenerate keys
./keygen import-privkey --file keys.txt
```

### Debugging Tips

#### Inspect PSBT Contents

Use Bitcoin Core to inspect PSBT:

```bash
# Decode PSBT (Bitcoin Core)
bitcoin-cli decodepsbt "$(cat payment_12_unsigned_0_*.psbt)"
```

**Output shows:**

- Inputs and their UTXOs
- Outputs and amounts
- Fee
- Signatures status
- Missing information

#### Validate PSBT

```bash
# Validate PSBT format (Bitcoin Core)
bitcoin-cli analyzepsbt "$(cat payment_12_unsigned_0_*.psbt)"
```

**Output shows:**

- Number of inputs
- Required signatures per input
- Current signature count
- Missing signatures
- Estimated fee
- Next steps

#### Check Transaction Status

```bash
# Check transaction in database
sqlite3 data/db/btc_watch.db "SELECT * FROM btc_tx WHERE id = 12;"

# Check if transaction is on blockchain
bitcoin-cli getrawtransaction <txid> 1
```

---
