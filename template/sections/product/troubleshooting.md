## Troubleshooting

### Common Errors

#### Error: "Nonce already used for this transaction"

**Symptom:**

```
Error: nonce reuse detected - cannot sign with same nonce
Transaction ID: 15
```

**Cause:** Attempting to reuse a nonce (critical security violation).

**Solution:**

- **Never reuse nonces** - this will leak your private key
- Generate fresh nonces for each transaction
- If nonces were already generated, create a new transaction

#### Error: "Not all nonces collected"

**Symptom:**

```
Error: cannot proceed to Round 2 - missing nonces
Expected: 3 nonces
Received: 2 nonces
```

**Cause:** Trying to sign before all nonces are generated.

**Solution:**

```bash
# Check PSBT status
./keygen musig2 status --file payment_15_unsigned_0_*.psbt

# Ensure all wallets have generated nonces:
# - Keygen wallet
# - Sign wallet 1
# - Sign wallet 2
```

#### Error: "Invalid partial signature"

**Symptom:**

```
Error: partial signature verification failed
Signer: auth1
```

**Cause:** Partial signature doesn't match expected format or is corrupted.

**Solution:**

```bash
# Verify PSBT integrity
sha256sum payment_15_unsigned_2_*.psbt

# Re-sign with the problematic wallet
./sign musig2 sign --file payment_15_unsigned_1_*.psbt
```

#### Error: "Signature aggregation failed"

**Symptom:**

```
Error: failed to aggregate signatures - verification failed
```

**Cause:** One or more partial signatures are invalid or missing.

**Solution:**

```bash
# Check all partial signatures are present
./watch musig2 status --file payment_15_unsigned_3_*.psbt

# Expected output should show:
# - Nonces: 3/3 ✓
# - Partial signatures: 3/3 ✓

# If any are missing, re-run signing for that wallet
```

#### Error: "MuSig2 not supported by Bitcoin Core"

**Symptom:**

```
Error: Bitcoin Core does not support Taproot/MuSig2
Current version: v21.0
```

**Cause:** Bitcoin Core version is too old.

**Solution:**

- Upgrade Bitcoin Core to v22.0 or higher
- Taproot and Schnorr signatures require Bitcoin Core 22.0+

### Debugging Tips

#### Inspect PSBT Status

```bash
# Check MuSig2 PSBT status
./keygen musig2 status --file payment_15_unsigned_0_*.psbt
```

**Output:**

```
PSBT Status:
  Transaction ID: 15
  Type: payment
  State: Round 1 - Nonce generation
  Nonces collected: 2/3
    ✓ Keygen: [66 bytes]
    ✓ Sign (auth1): [66 bytes]
    ✗ Sign (auth2): missing
  Partial signatures: 0/3
  Next step: Generate nonce for Sign (auth2)
```

#### Verify Nonce Uniqueness

```bash
# Check nonce uniqueness
./keygen musig2 verify-nonces --file payment_15_nonce_0_*.psbt
```

**Output:**

```
✓ All nonces are unique
✓ No nonce reuse detected
✓ Safe to proceed to Round 2
```

#### Check Signature Validity

```bash
# Verify partial signatures
./watch musig2 verify --file payment_15_unsigned_3_*.psbt
```

**Output:**

```
Partial Signature Verification:
  ✓ Keygen signature: valid
  ✓ Sign (auth1) signature: valid
  ✓ Sign (auth2) signature: valid
  ✓ All signatures compatible for aggregation
```

---
