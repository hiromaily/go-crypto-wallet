## FAQ

### General Questions

**Q: Will the migration cause downtime?**

A: Yes, a brief maintenance window (2-4 hours) is required to complete pending transactions and deploy new binaries. However, the system will be offline only during binary deployment (15-30 minutes).

**Q: Do I need to regenerate keys or addresses?**

A: No, keys and addresses are not affected by PSBT migration. Only the transaction file format changes.

**Q: Can I use PSBT and CSV simultaneously?**

A: No, the PSBT-enabled version removes CSV support. You must complete all CSV transactions before migration.

**Q: What happens to existing transactions after migration?**

A: Completed transactions (already broadcast) are not affected. Only pending transactions need to be completed or converted.

### Technical Questions

**Q: Are PSBT files larger than CSV files?**

A: Yes, PSBT files contain more metadata (UTXO info, scripts, derivation paths). However, files are base64-encoded and compressed, so the size difference is minimal (typically 20-30% larger).

**Q: Can I inspect PSBT contents?**

A: Yes, use Bitcoin Core:

```bash
bitcoin-cli decodepsbt "$(cat transaction.psbt)"
```

**Q: Are PSBTs compatible with Bitcoin Core?**

A: Yes, PSBT is a Bitcoin standard (BIP 174) fully supported by Bitcoin Core v0.17+.

**Q: Can I convert CSV files to PSBT?**

A: Technically possible but complex. Recommended approach is to complete CSV transactions before migration.

**Q: Does PSBT support all address types?**

A: Yes, PSBT supports all Bitcoin address types: P2PKH, P2SH, P2WPKH, P2WSH, and P2TR (Taproot).

### Operational Questions

**Q: How do I transfer PSBT files between wallets?**

A: Same as CSV files - use USB drives, QR codes, or secure file transfer. PSBT files have `.psbt` extension.

**Q: Do commands change with PSBT?**

A: No, commands remain the same. The format change is transparent to operators:

```bash
# Same commands work with PSBT
./watch create deposit --fee 0.0001
./keygen sign --file transaction.psbt
./watch send --file transaction.psbt
```

**Q: How do I know if a PSBT is fully signed?**

A: Check the filename:

- `_unsigned_0_*.psbt` - No signatures
- `_unsigned_1_*.psbt` - Partially signed (1 signature)
- `_signed_2_*.psbt` - Fully signed (2 signatures)

**Q: What if I lose a PSBT file?**

A: PSBT files can be recreated from database transaction records. Contact support for recovery procedures.

### Troubleshooting

**Q: Error: "PSBT is not fully signed"**

A: The PSBT needs more signatures. For multisig:

1. Check current signature count in filename
2. Send to next signer (Keygen → Sign)
3. Verify all required signatures collected

**Q: Error: "Invalid PSBT format"**

A: PSBT file may be corrupted:

1. Verify file integrity (checksum)
2. Don't edit PSBT files manually
3. Re-create transaction if needed

**Q: Error: "Failed to read PSBT file"**

A: Check file extension:

```bash
# Must have .psbt extension
mv transaction transaction.psbt
```

**Q: Can I edit PSBT files?**

A: No, never edit PSBT files manually. They are base64-encoded binary data. Use wallet commands to create and sign PSBTs.

---
