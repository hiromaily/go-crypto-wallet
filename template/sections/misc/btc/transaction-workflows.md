## Transaction Workflows

### Workflow 1: Payment Transaction with MuSig2 (3-of-3)

**Scenario**: Sending funds to external addresses using MuSig2 multisig

#### Step 1: Create MuSig2 Address (One-Time Setup)

```bash
# On Keygen wallet - create MuSig2 Taproot addresses
./keygen create musig2-address --account payment
```

**Output:**

```
✓ MuSig2 Taproot addresses created
Account: payment
Addresses created: 10
Address type: P2TR (Taproot)
Multisig type: MuSig2 3-of-3
Example address: tb1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxq...
```

#### Step 2: Create Unsigned Transaction (Watch Wallet)

```bash
# On the online Watch wallet system
./watch create payment --multisig-type musig2 --fee 0.0001
```

**Output:**

```
Created unsigned MuSig2 transaction: payment_15_unsigned_0_1735680000000000000.psbt
Transaction ID: 15
Inputs: 2
Outputs: 2 (1 recipient + 1 change)
Amount: 0.5 BTC
Fee: 0.0001 BTC
Multisig: MuSig2 3-of-3
Round 1: Ready for nonce generation
```

**File Generated:**

- `data/tx/btc/payment_15_unsigned_0_1735680000000000000.psbt`

#### Step 3: Round 1 - Generate Nonces (All Wallets, Parallel)

**Important**: Nonce generation can be done **in parallel** for all three wallets.

##### Keygen Wallet

```bash
# Transfer PSBT to Keygen wallet
# On offline Keygen wallet system
./keygen musig2 nonce --file /media/usb/payment_15_unsigned_0_*.psbt
```

**Output:**

```
✓ MuSig2 nonce generated successfully
Wallet: Keygen
Nonce size: 66 bytes
Output file: payment_15_unsigned_0_1735680000000000001.psbt
Status: Nonce added to PSBT
Next: Transfer to Sign wallets for their nonce generation
```

##### Sign Wallet 1

```bash
# Transfer updated PSBT to Sign wallet 1
# On offline Sign wallet #1 system
./sign musig2 nonce --file /media/usb/payment_15_unsigned_0_*.psbt
```

**Output:**

```
✓ MuSig2 nonce generated successfully
Wallet: Sign (auth1)
Nonce size: 66 bytes
Output file: payment_15_unsigned_0_1735680000000000002.psbt
Status: Nonce added to PSBT
Next: Transfer to Sign wallet #2 for nonce generation
```

##### Sign Wallet 2

```bash
# Transfer updated PSBT to Sign wallet 2
# On offline Sign wallet #2 system
./sign musig2 nonce --file /media/usb/payment_15_unsigned_0_*.psbt
```

**Output:**

```
✓ MuSig2 nonce generated successfully
Wallet: Sign (auth2)
Nonce size: 66 bytes
Output file: payment_15_nonce_0_1735680000000000003.psbt
Status: All nonces collected (3/3)
Next: Round 2 - Partial signature creation
```

**Note**: After Step 3, all three nonces are stored in the PSBT file.

#### Step 4: Round 2 - Create Partial Signatures (All Wallets, Sequential)

**Important**: Signing must be done **sequentially** after all nonces are collected.

##### Keygen Wallet Signs

```bash
# Transfer PSBT with all nonces to Keygen wallet
# On offline Keygen wallet system
./keygen musig2 sign --file /media/usb/payment_15_nonce_0_*.psbt
```

**Output:**

```
✓ MuSig2 partial signature created
Wallet: Keygen
Signature size: 32 bytes
Output file: payment_15_unsigned_1_1735680000000000004.psbt
Status: Partial signature 1/3
Next: Transfer to Sign wallet #1 for signing
```

##### Sign Wallet 1 Signs

```bash
# Transfer PSBT to Sign wallet 1
# On offline Sign wallet #1 system
./sign musig2 sign --file /media/usb/payment_15_unsigned_1_*.psbt
```

**Output:**

```
✓ MuSig2 partial signature created
Wallet: Sign (auth1)
Signature size: 32 bytes
Output file: payment_15_unsigned_2_1735680000000000005.psbt
Status: Partial signature 2/3
Next: Transfer to Sign wallet #2 for signing
```

##### Sign Wallet 2 Signs

```bash
# Transfer PSBT to Sign wallet 2
# On offline Sign wallet #2 system
./sign musig2 sign --file /media/usb/payment_15_unsigned_2_*.psbt
```

**Output:**

```
✓ MuSig2 partial signature created
Wallet: Sign (auth2)
Signature size: 32 bytes
Output file: payment_15_unsigned_3_1735680000000000006.psbt
Status: All partial signatures collected (3/3)
Next: Transfer to Watch wallet for signature aggregation
```

#### Step 5: Aggregate Signatures (Watch Wallet)

```bash
# Transfer PSBT with all partial signatures to Watch wallet
# On online Watch wallet system
./watch musig2 aggregate --file /media/usb/payment_15_unsigned_3_*.psbt
```

**Output:**

```
✓ MuSig2 signatures aggregated successfully
Final signature size: 64 bytes (Schnorr)
Verification: PASSED
Output file: payment_15_signed_3_1735680000000000007.psbt
Status: Ready for broadcasting
Transaction size: 215 bytes
Traditional multisig size: ~370 bytes
Size reduction: 41.9%
Fee reduction: 41.9%
```

#### Step 6: Broadcast Transaction (Watch Wallet)

```bash
# On online Watch wallet system
./watch send --file payment_15_signed_3_*.psbt
```

**Output:**

```
✓ Transaction broadcast successfully
Transaction hash: a1b2c3d4e5f6789...
Status: Sent
Confirmations: 0 (pending)
On-chain appearance: Single-signature transaction (MuSig2)
Privacy: Maximum (indistinguishable from single-sig)
```

---

### Workflow 2: Deposit Transaction (Single-Signature Taproot Spend)

**Scenario**: Receiving funds from users (client → deposit account) using a single-signature Taproot address.

For deposit transactions using single-signature Taproot addresses:

#### Step 1: Create Unsigned Transaction

```bash
# On Watch wallet
./watch create deposit --multisig-type musig2 --fee 0.0001
```

#### Step 2: Round 1 - Generate Nonce (Keygen Wallet Only)

```bash
# On Keygen wallet
./keygen musig2 nonce --file deposit_8_unsigned_0_*.psbt
```

#### Step 3: Round 2 - Sign (Keygen Wallet Only)

```bash
# On Keygen wallet
./keygen musig2 sign --file deposit_8_nonce_0_*.psbt
```

**Note**: For single-signature, only Keygen wallet participates.

#### Step 4: Finalize and Broadcast

```bash
# On Watch wallet
./watch musig2 aggregate --file deposit_8_unsigned_1_*.psbt
./watch send --file deposit_8_signed_1_*.psbt
```

---
