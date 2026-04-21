## File Management

### File Locations

#### Watch Wallet

```
data/tx/btc/
├── payment_15_unsigned_0_*.psbt     # Created here (unsigned)
├── payment_15_nonce_0_*.psbt        # Receives after Round 1
├── payment_15_unsigned_3_*.psbt     # Receives after Round 2
└── payment_15_signed_3_*.psbt       # Creates after aggregation
```

#### Keygen Wallet

```
data/tx/btc/
├── payment_15_unsigned_0_...0.psbt     # Receives from Watch
├── payment_15_unsigned_0_...1.psbt     # Creates after nonce gen
├── payment_15_nonce_0_*.psbt           # Receives from Sign wallets
└── payment_15_unsigned_1_*.psbt        # Creates after signing
```

#### Sign Wallet 1

```
data/tx/btc/
├── payment_15_unsigned_0_...1.psbt     # Receives from Keygen
├── payment_15_unsigned_0_...2.psbt     # Creates after nonce gen
├── payment_15_unsigned_1_*.psbt        # Receives from Keygen
└── payment_15_unsigned_2_*.psbt        # Creates after signing
```

#### Sign Wallet 2

```
data/tx/btc/
├── payment_15_unsigned_0_...2.psbt     # Receives from Sign 1
├── payment_15_unsigned_0_...3.psbt     # Creates after nonce gen
├── payment_15_unsigned_2_*.psbt        # Receives from Sign 1
└── payment_15_unsigned_3_*.psbt        # Creates after signing
```

### File Transfer Best Practices

#### For Air-Gapped Systems

1. **Use USB Drives**

   ```bash
   # Mount USB
   mount /dev/sdb1 /media/usb

   # Copy PSBT
   cp data/tx/btc/payment_15_unsigned_0_*.psbt /media/usb/

   # Safely unmount
   umount /media/usb
   ```

2. **Verify File Integrity**

   ```bash
   # Generate checksum
   sha256sum payment_15_unsigned_0_*.psbt > checksum.txt

   # Verify on destination
   sha256sum -c checksum.txt
   ```

3. **Track PSBT State**

   ```bash
   # Check PSBT status
   ./keygen musig2 status --file payment_15_unsigned_0_*.psbt
   ```

#### Security Considerations

- ✅ **Virus scan USB drives** before use on offline systems
- ✅ **Use dedicated USB drives** for wallet operations only
- ✅ **Verify file checksums** after transfer
- ✅ **Track PSBT state** to avoid confusion
- ❌ **Never** connect offline wallets to the internet
- ❌ **Never** reuse nonces (critical security issue)

### Nonce Management

#### Nonce Lifecycle

```
1. Generate → Store in PSBT → Collect all nonces → Sign → Delete
```

**Critical Rules:**

- Each transaction needs fresh nonces
- Nonces are stored securely during Round 1
- Nonces are deleted after signing (Round 2)
- **Never reuse nonces** (will leak private key)

#### Nonce Storage

Nonces are stored in:

- **PSBT proprietary fields** (during transaction flow)
- **Wallet database** (temporary, for Round 2 signing)

**Cleanup**: Nonces are automatically deleted after signing.

---
