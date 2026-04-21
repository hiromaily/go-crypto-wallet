## PSBT Basics

### PSBT File Format

#### File Naming Convention

PSBT files follow this naming pattern:

```
{actionType}_{txID}_{txType}_{signedCount}_{timestamp}.psbt
```

**Components:**

- `actionType`: Transaction type (`deposit`, `payment`, `transfer`)
- `txID`: Database transaction ID
- `txType`: Status (`unsigned`, `signed`)
- `signedCount`: Number of signatures collected (0, 1, 2, ...)
- `timestamp`: Unix timestamp in nanoseconds

**Examples:**

```
deposit_8_unsigned_0_1534744535097796209.psbt    # Unsigned deposit (0 signatures)
deposit_8_unsigned_1_1534744535097796210.psbt    # Partially signed (1 signature)
deposit_8_signed_2_1534744535097796211.psbt      # Fully signed (2 signatures)
```

#### File Content

PSBT files contain base64-encoded binary data:

```
cHNidP8BAHECAAAAAZt/TvyKa6hVH3n8FwUPKA...
```

**Do not edit PSBT files manually!** Use wallet commands to create and sign PSBTs.

### PSBT States

A PSBT progresses through these states:

```
1. Unsigned (0 signatures) ───> Watch Wallet creates
                    │
2. Partially Signed (1 sig) ──> Keygen Wallet signs
                    │
3. Partially Signed (2 sig) ──> Sign Wallet signs (multisig)
                    │
4. Fully Signed ──────────────> Watch Wallet finalizes & broadcasts
```

---
