## Changelog

### Version 1.10 (2026-01-17)

- ✅ Fixed BCH Pattern 3 (3-of-3 Multisig) UTXO retrieval issue (Closes #423, PR #426)
- Resolved CashAddr format mismatch in `ListUnspentByAccount` address comparison
- Updated BCH Pattern Matrix to show all 3 patterns with correct E2E script paths
- Updated E2E Script Reference with proper BCH script locations
- Added implementation status details for BCH Pattern 3
- All BCH patterns now have dedicated E2E scripts in `scripts/operation/bch/e2e/`

### Version 1.9 (2026-01-17)

- Updated BCH Pattern Matrix with accurate pattern numbering (1, 2, 3)
- Added BCH limitations section (no SegWit, Taproot, Schnorr, MuSig2)
- Updated BCH Pattern 3 (3-of-3 Multisig) with detailed workflow and signing flow
- Added cross-reference link to BCH Technical Reference
- Renamed "BCH Pattern 2" to "BCH Pattern 3" to reflect correct pattern numbering

### Version 1.8 (2026-01-16)

- ✅ Pattern 11 (P2TR Tapscript M-of-N) framework implemented
- E2E script `e2e-p11-p2tr-tapscript.sh` created (Closes #381)
- Implements Tapscript Script Path spending framework with 2-of-3 threshold
- Uses BIP86 key derivation + BIP342 Tapscript semantics
- Script tree with Merkle proof and control block structure
- M × Schnorr signatures for Script Path spend
- Address format: `bcrt1p...` (62 chars, Bech32m encoding)
- ~50% smaller than P2WSH 2-of-3 multisig
- Enhanced privacy: unused script paths hidden in Merkle tree
- Note: Full Tapscript implementation pending (currently uses placeholder)

### Version 1.7 (2026-01-16)

- ✅ Pattern 9 (P2TR Taproot Single-sig) is now fully working
- E2E script `e2e-p9-p2tr-singlesig.sh` completed and verified (Closes #377)
- Fixed Taproot address derivation to use x-only public keys (32 bytes)
- Implements BIP86 key derivation for Taproot key path spending
- Uses Schnorr signatures (BIP340, 64 bytes)
- Most efficient single-sig transaction format
- Address format: `bcrt1p...` (62 chars, Bech32m encoding)

### Version 1.6 (2026-01-16)

- ✅ Pattern 8 (P2SH-P2WSH 3-of-3 Multisig) is now fully working
- E2E script `e2e-p8-p2sh-p2wsh-3of3.sh` completed and verified
- Fixed receiver address generation to use P2SH-SegWit format (Closes #374)
- P2SH-wrapped SegWit multisig with legacy compatibility (`2...` addresses in regtest)
- Implements BIP49 key derivation for 3-of-3 multisig
- All 3 signatures required (Keygen + Sign1 + Sign2)

### Version 1.5 (2026-01-16)

- ✅ Pattern 6 (P2WSH Native SegWit 2-of-3 Multisig) is now fully working
- E2E script `e2e-p6-p2wsh-2of3.sh` completed and verified
- Native SegWit multisig with Bech32 encoding (`bcrt1q...` 62-char addresses)
- Added native SegWit descriptor support (`wsh`, `wpkh`) in PSBT infrastructure
- Most efficient multisig format - no P2SH wrapper overhead
- Implements BIP84 key derivation and BIP67 sorted multisig keys

### Version 1.4 (2026-01-15)

- ✅ Pattern 5 (P2WPKH Native SegWit Single-sig) is now fully working
- E2E script `e2e-p5-p2wpkh-singlesig.sh` completed and verified
- Native SegWit with Bech32 encoding (`bcrt1q...` addresses)

### Version 1.3 (2026-01-15)

- ✅ Pattern 3 (P2SH-P2WPKH Single-sig) is now fully working
- E2E script `e2e-p3-p2sh-p2wpkh-singlesig.sh` completed and verified

### Version 1.2 (2026-01-15)

- ✅ Pattern 2 (P2PKH 2-of-3 Multisig) is now fully working
- Fixed key derivation path mismatch for multisig accounts (PR #357)
- Added detailed explanation of the fix and root cause
- Updated implementation status for 2-of-3 multisig

### Version 1.1 (2026-01-14)

- Initial comprehensive documentation of all patterns
