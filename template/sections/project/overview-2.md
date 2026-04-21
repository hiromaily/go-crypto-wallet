## Overview

### What is MuSig2?

MuSig2 is a two-round Schnorr multisignature protocol (BIP327) that enables multiple parties to create a single aggregated signature that is indistinguishable from a standard single-signature transaction on the blockchain. This provides:

- **Smaller transactions**: 30-50% size reduction compared to traditional P2WSH multisig
- **Lower fees**: Proportional to transaction size reduction
- **Better privacy**: Multisig transactions look like single-sig on-chain
- **Schnorr signatures**: Uses BIP340 Schnorr signatures via Taproot (P2TR)

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│   Interface Adapters (CLI Commands)                     │
│   - keygen create musig2-address                        │
│   - keygen musig2 nonce                                 │
│   - keygen musig2 sign                                  │
│   - sign musig2 nonce                                   │
│   - sign musig2 sign                                    │
│   - watch musig2 aggregate                              │
└────────────────────┬────────────────────────────────────┘
                     │ depends on
┌────────────────────▼────────────────────────────────────┐
│   Application Layer (Use Cases)                         │
│   Keygen:  CreateMuSig2AddressUseCase                   │
│            GenerateMuSig2NonceUseCase                   │
│            MuSig2SignUseCase                            │
│   Sign:    GenerateMuSig2NonceUseCase                   │
│            MuSig2SignUseCase                            │
│   Watch:   AggregateMuSig2SignaturesUseCase             │
└────────────────────┬────────────────────────────────────┘
                     │ depends on
┌────────────────────▼────────────────────────────────────┐
│   Domain Layer (Business Logic)                         │
│   - MuSig2 Types (domain/musig2/)                       │
│   - Validators                                          │
│   - Business Rules                                      │
└─────────────────────────────────────────────────────────┘
                     ▲ implements
┌────────────────────┴────────────────────────────────────┐
│   Infrastructure Layer (External Dependencies)          │
│   - MuSig2Service (btcd/btcec/v2/schnorr/musig2)       │
│   - AccountKeyRepository (MySQL)                        │
│   - AuthFullPubkeyRepository (MySQL)                    │
│   - FileStorage (PSBT files)                            │
└─────────────────────────────────────────────────────────┘
```

### Design Principles

1. **Clean Architecture**: Strict layer separation with dependency inversion
2. **Security First**: Nonce uniqueness enforced at multiple levels
3. **Type Safety**: Domain types for all MuSig2 operations
4. **Testability**: All components have clear interfaces
5. **Offline Support**: Keygen and Sign wallets work completely offline
6. **PSBT Integration**: MuSig2 data stored in PSBT proprietary fields

---
