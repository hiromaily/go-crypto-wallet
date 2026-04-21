## Transaction Workflows

### Workflow 1: Deposit Transaction (Single-Signature)

**Scenario**: Receiving funds from users (client → deposit account)

#### Step 1: Create Unsigned Transaction (Watch Wallet)

```bash
# On the online Watch wallet system
cd /path/to/go-crypto-wallet
./watch create deposit --fee 0.0001
```

**Output:**

```
Created unsigned transaction: deposit_8_unsigned_0_1534744535097796209.psbt
Transaction ID: 8
Inputs: 5
Outputs: 2
Amount: 1.5 BTC
Fee: 0.0001 BTC
```

**File Generated:**

- `data/tx/btc/deposit_8_unsigned_0_1534744535097796209.psbt`

#### Step 2: Transfer PSBT to Keygen Wallet

Transfer the unsigned PSBT file to the offline Keygen wallet system:

```bash
# Using USB drive or secure file transfer
cp data/tx/btc/deposit_8_unsigned_0_*.psbt /media/usb/
```

#### Step 3: Sign Transaction (Keygen Wallet)

```bash
# On the offline Keygen wallet system
cd /path/to/go-crypto-wallet
./keygen sign --file /media/usb/deposit_8_unsigned_0_1534744535097796209.psbt
```

**Output:**

```
Signed transaction successfully
Input file: deposit_8_unsigned_0_1534744535097796209.psbt
Output file: deposit_8_signed_1_1534744535097796210.psbt
Is fully signed: true
Transaction ready for broadcasting
```

**File Generated:**

- `data/tx/btc/deposit_8_signed_1_1534744535097796210.psbt`

#### Step 4: Transfer Signed PSBT Back to Watch Wallet

```bash
# Copy signed PSBT back to Watch wallet
cp data/tx/btc/deposit_8_signed_1_*.psbt /media/usb/
```

#### Step 5: Broadcast Transaction (Watch Wallet)

```bash
# On the online Watch wallet system
./watch send --file /media/usb/deposit_8_signed_1_1534744535097796210.psbt
```

**Output:**

```
Transaction broadcast successfully
Transaction hash: a1b2c3d4e5f6...
Status: Sent
```

---

### Workflow 2: Payment Transaction (Multisig 2-of-2)

**Scenario**: Sending funds to external addresses (payment account → external)

#### Step 1: Create Unsigned Transaction (Watch Wallet)

```bash
./watch create payment --fee 0.0002
```

**Output:**

```
Created unsigned transaction: payment_12_unsigned_0_1534744600000000000.psbt
Transaction ID: 12
Inputs: 3
Outputs: 2 (1 recipient + 1 change)
Amount: 0.5 BTC
Fee: 0.0002 BTC
Multisig: 2-of-2 (requires 2 signatures)
```

#### Step 2: First Signature (Keygen Wallet)

```bash
# Transfer to Keygen wallet
./keygen sign --file payment_12_unsigned_0_1534744600000000000.psbt
```

**Output:**

```
Signed transaction successfully
Output file: payment_12_unsigned_1_1534744600000000001.psbt
Is fully signed: false
Signatures: 1/2
Next: Transfer to Sign wallet for second signature
```

**Note:** Transaction is **not** fully signed yet (1 of 2 signatures).

#### Step 3: Second Signature (Sign Wallet)

```bash
# Transfer to Sign wallet
./sign sign --file payment_12_unsigned_1_1534744600000000001.psbt
```

**Output:**

```
Signed transaction successfully
Output file: payment_12_signed_2_1534744600000000002.psbt
Is fully signed: true
Signatures: 2/2
Transaction ready for broadcasting
```

#### Step 4: Broadcast Transaction (Watch Wallet)

```bash
./watch send --file payment_12_signed_2_1534744600000000002.psbt
```

**Output:**

```
Transaction broadcast successfully
Transaction hash: f6e5d4c3b2a1...
Status: Sent
```

---

### Workflow 3: Transfer Transaction (Multisig 2-of-2)

**Scenario**: Moving funds between internal accounts (stored → payment)

#### Step 1: Create Unsigned Transaction (Watch Wallet)

```bash
./watch create transfer --sender stored --receiver payment --amount 10.0 --fee 0.0003
```

**Output:**

```
Created unsigned transaction: transfer_15_unsigned_0_1534744700000000000.psbt
Transaction ID: 15
Sender: stored account
Receiver: payment account
Amount: 10.0 BTC
Fee: 0.0003 BTC
Multisig: 2-of-2
```

#### Steps 2-4: Same as Payment Workflow

Follow the same signing and broadcasting steps as the payment transaction:

1. Keygen wallet signs (first signature)
2. Sign wallet signs (second signature)
3. Watch wallet broadcasts (finalization)

---
