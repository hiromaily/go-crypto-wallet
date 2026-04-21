## File Management

### File Locations

#### Watch Wallet

```
data/tx/btc/
├── deposit_8_unsigned_0_*.psbt      # Created here
├── deposit_8_signed_1_*.psbt        # Receives from Keygen
├── payment_12_unsigned_0_*.psbt     # Created here
└── payment_12_signed_2_*.psbt       # Receives from Sign
```

#### Keygen Wallet

```
data/tx/btc/
├── deposit_8_unsigned_0_*.psbt      # Receives from Watch
├── deposit_8_signed_1_*.psbt        # Creates here
├── payment_12_unsigned_0_*.psbt     # Receives from Watch
└── payment_12_unsigned_1_*.psbt     # Creates here
```

#### Sign Wallet

```
data/tx/btc/
├── payment_12_unsigned_1_*.psbt     # Receives from Keygen
└── payment_12_signed_2_*.psbt       # Creates here
```

### File Transfer Best Practices

#### For Air-Gapped Systems

1. **Use USB Drives**

   ```bash
   # Mount USB
   mount /dev/sdb1 /media/usb

   # Copy PSBT
   cp data/tx/btc/payment_12_unsigned_0_*.psbt /media/usb/

   # Safely unmount
   umount /media/usb
   ```

2. **Use QR Codes** (for smaller transactions)

   ```bash
   # Generate QR code
   qrencode -o psbt.png < payment_12_unsigned_0_*.psbt

   # Scan and decode on offline system
   zbarimg psbt.png > payment_12_unsigned_0_*.psbt
   ```

3. **Use Optical Data Transfer** (most secure)
   - Print PSBT as QR code or text
   - Manually type or scan on offline system

#### Security Considerations

- ✅ **Virus scan USB drives** before use on offline systems
- ✅ **Use dedicated USB drives** for wallet operations only
- ✅ **Verify file integrity** with checksums (sha256sum)
- ❌ **Never** connect offline wallets to the internet
- ❌ **Never** use the same USB drive for other purposes

### File Cleanup

#### Automatic Cleanup

Wallets automatically manage PSBT files:

- Unsigned PSBTs are kept until signed
- Partially signed PSBTs are kept until fully signed
- Fully signed PSBTs are kept until broadcasted
- Broadcasted transaction PSBTs can be archived

#### Manual Cleanup

```bash
# Archive old PSBT files
mkdir -p data/tx/btc/archive/2025-01
mv data/tx/btc/*_signed_*_1735*.psbt data/tx/btc/archive/2025-01/

# Compress archives
tar -czf psbt-archive-2025-01.tar.gz data/tx/btc/archive/2025-01/
```

---
