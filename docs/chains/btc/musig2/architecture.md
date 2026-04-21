<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/btc/musig2/architecture.tpl.md · Run `make docs` to regenerate.
-->

# MuSig2 Architecture Documentation

This document describes the architecture and implementation of MuSig2 (Simple Two-Round Schnorr Multisignatures) support in the go-crypto-wallet project, following Clean Architecture principles.

## Table of Contents

1. [Overview](#overview)
2. [Architecture Layers](#architecture-layers)
3. [Component Interactions](#component-interactions)
4. [Data Flow](#data-flow)
5. [Security Architecture](#security-architecture)
6. [Database Schema](#database-schema)
7. [API Reference](#api-reference)
8. [Implementation Notes](#implementation-notes)

---

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

## Architecture Layers

### Domain Layer

**Location**: `internal/domain/musig2/`

**Purpose**: Pure business logic with zero infrastructure dependencies

**Components**:

- **Types**: MuSig2-specific value objects (NonceCommitment, PartialSignature, AggregatedKey, SigningSession)
- **Validators**: Business rule validation (nonce uniqueness, signer count, signature validation)
- **Errors**: Domain-specific errors for MuSig2 operations

**Key Characteristics**:

- No external dependencies (no database, no API clients)
- Immutable value objects
- Validation logic for business rules
- Domain errors with clear semantics

**Example Types**:

```go
// NonceCommitment represents a MuSig2 nonce commitment from a signer
type NonceCommitment struct {
    SignerID string
    Nonce    [66]byte  // Public nonce (66 bytes)
}

// PartialSignature represents a partial signature from a signer
type PartialSignature struct {
    SignerID  string
    Signature []byte  // Partial signature scalar (32 bytes)
}

// AggregatedKey represents an aggregated public key
type AggregatedKey struct {
    PublicKey   []byte  // Aggregated public key (33 bytes compressed)
    TweakApplied bool    // Whether Taproot tweak was applied
}

// SigningSession tracks the state of a MuSig2 signing session
type SigningSession struct {
    SessionID        string
    ParticipantCount int
    Nonces           []NonceCommitment
    PartialSigs      []PartialSignature
}
```

**Key Validators**:

```go
// ValidateSignerCount validates the number of signers for MuSig2
func ValidateSignerCount(count int) error

// ValidateNonceUniqueness ensures all nonces are unique
func ValidateNonceUniqueness(nonces []NonceCommitment) error

// ValidatePartialSignatures validates partial signatures
func ValidatePartialSignatures(sigs []PartialSignature, expected int) error

// ValidatePublicKeysForMuSig2 validates public keys for aggregation
func ValidatePublicKeysForMuSig2(pubKeys [][]byte) error

// ValidateAggregatedPublicKey validates aggregated public key format
func ValidateAggregatedPublicKey(aggKey []byte) error

// ValidateSigningSessionComplete validates session is ready for aggregation
func ValidateSigningSessionComplete(session *SigningSession) error
```

### Application Layer (Use Cases)

**Location**: `internal/application/usecase/`

**Purpose**: Orchestrate business logic by coordinating domain objects and infrastructure services

**Organization**:

```
internal/application/usecase/
├── keygen/
│   ├── interfaces.go              # Use case interfaces
│   └── btc/
│       ├── create_musig2_address.go
│       ├── musig2_nonce.go
│       └── musig2_sign.go
├── sign/
│   ├── interfaces.go
│   └── btc/
│       ├── musig2_nonce.go
│       └── musig2_sign.go
└── watch/
    ├── interfaces.go
    └── btc/
        └── musig2_aggregate.go
```

#### Keygen Wallet Use Cases

**1. CreateMuSig2AddressUseCase**

Creates MuSig2 Taproot addresses by aggregating public keys from all signers.

```go
type CreateMuSig2AddressUseCase interface {
    Create(ctx context.Context, input CreateMuSig2AddressInput) error
}

type CreateMuSig2AddressInput struct {
    AccountType account.AccountType
}
```

**Process**:

1. Validate account is multisig account
2. Retrieve public keys from all signers (auth_fullpubkey table)
3. Add account's own public key
4. Aggregate public keys using MuSig2Service
5. Apply Taproot tweak (BIP86)
6. Create Taproot address (P2TR)
7. Store address in account_key table

**2. GenerateMuSig2NonceUseCase (Round 1)**

Generates cryptographically secure nonces for a specific transaction.

```go
type GenerateMuSig2NonceUseCase interface {
    Generate(ctx context.Context, input GenerateMuSig2NonceInput) (*GenerateMuSig2NonceOutput, error)
}

type GenerateMuSig2NonceInput struct {
    TransactionID  int64
    PSBTData       []byte
    AccountType    account.AccountType
}

type GenerateMuSig2NonceOutput struct {
    Nonce     [66]byte  // Public nonce
    PSBTData  []byte    // Updated PSBT with nonce
}
```

**Process**:

1. Validate PSBT and transaction
2. Generate secure random nonce using MuSig2Service
3. Store nonce in PSBT proprietary field
4. Track nonce in database (for uniqueness enforcement)
5. Return updated PSBT

**3. MuSig2SignUseCase (Round 2)**

Creates partial signature after all nonces are collected.

```go
type MuSig2SignUseCase interface {
    Sign(ctx context.Context, input MuSig2SignInput) (*MuSig2SignOutput, error)
}

type MuSig2SignInput struct {
    PSBTData       []byte
    AccountType    account.AccountType
}

type MuSig2SignOutput struct {
    PartialSignature []byte
    PSBTData         []byte  // Updated PSBT with partial signature
}
```

**Process**:

1. Validate all nonces are present in PSBT
2. Extract nonces from PSBT
3. Retrieve private key for signing
4. Create partial signature using MuSig2Service
5. Store partial signature in PSBT
6. Return updated PSBT

#### Sign Wallet Use Cases

Sign wallet has similar use cases for nonce generation and signing:

**1. GenerateMuSig2NonceUseCase**: Same as Keygen, but uses auth key
**2. MuSig2SignUseCase**: Same as Keygen, but uses auth key

**Key Difference**: Sign wallets use `auth_account_key` table instead of `account_key` table.

#### Watch Wallet Use Cases

**AggregateMuSig2SignaturesUseCase**

Aggregates partial signatures from all signers into final signature.

```go
type AggregateMuSig2SignaturesUseCase interface {
    Aggregate(ctx context.Context, input AggregateMuSig2SignaturesInput) (*AggregateMuSig2SignaturesOutput, error)
}

type AggregateMuSig2SignaturesInput struct {
    PSBTData []byte
}

type AggregateMuSig2SignaturesOutput struct {
    FinalSignature [64]byte  // Schnorr signature
    PSBTData       []byte    // Finalized PSBT
}
```

**Process**:

1. Validate all partial signatures are present in PSBT
2. Extract partial signatures
3. Aggregate signatures using MuSig2Service
4. Verify aggregated signature
5. Finalize PSBT with final signature
6. Return finalized PSBT ready for broadcasting

### Infrastructure Layer

**Location**: `internal/infrastructure/`

**Purpose**: Implement interfaces defined by domain layer; handle external dependencies

#### MuSig2Service

**Location**: `internal/infrastructure/api/btc/btc/musig2.go`

**Purpose**: Wrapper around `github.com/btcsuite/btcd/btcec/v2/schnorr/musig2` library

**Key Methods**:

```go
type MuSig2Service struct {
    chainConfig *chaincfg.Params
}

// AggregatePublicKeys aggregates multiple public keys into a single key
func (s *MuSig2Service) AggregatePublicKeys(
    pubKeys []*btcec.PublicKey,
    applyTaprootTweak bool,
) (*btcec.PublicKey, error)

// GenerateNonce generates a secure nonce for MuSig2 signing
func (s *MuSig2Service) GenerateNonce(
    privateKey *btcec.PrivateKey,
    pubKeys []*btcec.PublicKey,
) ([66]byte, error)

// CreatePartialSignature creates a partial signature
func (s *MuSig2Service) CreatePartialSignature(
    privateKey *btcec.PrivateKey,
    pubKeys []*btcec.PublicKey,
    nonces [][66]byte,
    messageHash [32]byte,
) ([]byte, error)

// AggregateSignatures combines partial signatures into final signature
func (s *MuSig2Service) AggregateSignatures(
    pubKeys []*btcec.PublicKey,
    nonces [][66]byte,
    partialSigs [][]byte,
    messageHash [32]byte,
) (*schnorr.Signature, error)

// VerifyAggregatedSignature verifies the final aggregated signature
func (s *MuSig2Service) VerifyAggregatedSignature(
    aggregatedPubKey *btcec.PublicKey,
    messageHash [32]byte,
    signature *schnorr.Signature,
) bool
```

**Library Integration**:

The service uses `github.com/btcsuite/btcd/btcec/v2/schnorr/musig2` (v2.3.6):

```go
import (
    "github.com/btcsuite/btcd/btcec/v2"
    "github.com/btcsuite/btcd/btcec/v2/schnorr"
    "github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
)

// Example: Key aggregation
ctx, err := musig2.NewContext(
    privateKey,
    true,  // Sort keys
    musig2.WithKnownSigners(allPublicKeys),
    musig2.WithBip86TweakCtx(),  // Taproot tweak
)

// Example: Nonce generation
session, err := ctx.NewSession()
pubNonce := session.PublicNonce()  // [66]byte

// Example: Signing
partialSig, err := session.Sign(messageHash)

// Example: Aggregation
haveAll, err := session.CombineSig(partialSig)
finalSig := session.FinalSig()
```

#### Repository Layer

**AccountKeyRepository**

**Location**: `internal/infrastructure/repository/cold/account_key.go`

**Purpose**: Manage account keys including MuSig2 Taproot addresses

**Key Methods**:

```go
type AccountKeyRepositorier interface {
    GetAllAddrStatus(accountType account.AccountType, addrStatus address.AddrStatus) ([]*sqlc.AccountKey, error)
    UpdateMultisigAddr(accountType account.AccountType, item *sqlc.AccountKey) (int64, error)
    // ... other methods
}
```

**AuthFullPubkeyRepository**

**Location**: `internal/infrastructure/repository/cold/auth_fullpubkey.go`

**Purpose**: Manage full public keys from auth accounts (Sign wallets)

**Key Methods**:

```go
type AuthFullPubkeyRepositorier interface {
    GetOne(authType account.AuthType) (*sqlc.AuthFullpubkey, error)
    // ... other methods
}
```

#### File Storage

**Location**: `internal/infrastructure/storage/file/`

**Purpose**: Handle PSBT file operations

MuSig2 data is stored in PSBT proprietary fields:

- **Public Nonces**: 66 bytes per signer
- **Partial Signatures**: 32 bytes per signer
- **Aggregated Signature**: 64 bytes (final)

### Interface Adapters Layer

**Location**: `internal/interface-adapters/cli/`

**Purpose**: CLI commands that invoke use cases

**Commands**:

```bash
# Keygen Wallet
keygen create musig2-address --account payment
keygen musig2 nonce --file payment_15_unsigned_0.psbt
keygen musig2 sign --file payment_15_nonce_0.psbt

# Sign Wallet
sign musig2 nonce --file payment_15_unsigned_0.psbt
sign musig2 sign --file payment_15_unsigned_1.psbt

# Watch Wallet
watch musig2 aggregate --file payment_15_unsigned_3.psbt
watch send --file payment_15_signed_3.psbt
```

**Command Flow**:

```
CLI Command → Parse Args → Create Use Case → Execute → Handle Result → Display Output
```

---

## Component Interactions

### Address Creation Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. CLI Command                                              │
│    keygen create musig2-address --account payment           │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. CreateMuSig2AddressUseCase                               │
│    - Validate account is multisig                           │
│    - Get signer public keys                                 │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. AuthFullPubkeyRepository                                 │
│    - GetOne(auth1) → pubKey1                                │
│    - GetOne(auth2) → pubKey2                                │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. AccountKeyRepository                                     │
│    - GetAllAddrStatus(payment, PrivKeyImported)             │
│    → Returns account keys needing multisig addresses        │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 5. MuSig2Service.AggregatePublicKeys()                      │
│    Input: [pubKey1, pubKey2, accountPubKey]                 │
│    - Sort public keys                                       │
│    - Aggregate using MuSig2 protocol                        │
│    - Apply Taproot tweak (BIP86)                            │
│    Output: aggregatedPubKey                                 │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 6. Create Taproot Address                                   │
│    - Create P2TR address from aggregatedPubKey              │
│    - Format: bc1p... (mainnet) or tb1p... (testnet)        │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 7. AccountKeyRepository.UpdateMultisigAddr()                │
│    - Store P2TR address                                     │
│    - Update addr_status                                     │
└─────────────────────────────────────────────────────────────┘
```

### Transaction Signing Flow (Two-Round Protocol)

#### Round 1: Nonce Generation (Parallel)

```
┌─────────────────────────────────────────────────────────────┐
│ KEYGEN WALLET                                               │
│ 1. keygen musig2 nonce --file payment_15_unsigned_0.psbt   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. GenerateMuSig2NonceUseCase (Keygen)                     │
│    - Parse PSBT                                             │
│    - Validate transaction                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. MuSig2Service.GenerateNonce()                            │
│    - Generate secure random nonce                           │
│    - Create public nonce (66 bytes)                         │
│    Output: [66]byte nonce                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. Store in PSBT                                            │
│    - Add nonce to PSBT proprietary field                    │
│    - Keygen nonce: field_id = "musig2_nonce_keygen"        │
│    Output: payment_15_unsigned_0_...1.psbt                  │
└─────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────┐
│ SIGN WALLET 1 (Parallel)                                    │
│ 1. sign musig2 nonce --file payment_15_unsigned_0_...1.psbt│
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. GenerateMuSig2NonceUseCase (Sign)                       │
│    - Parse PSBT                                             │
│    - Validate transaction                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. MuSig2Service.GenerateNonce()                            │
│    - Generate secure random nonce                           │
│    Output: [66]byte nonce                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. Store in PSBT                                            │
│    - Add nonce to PSBT proprietary field                    │
│    - Sign nonce: field_id = "musig2_nonce_sign1"           │
│    Output: payment_15_unsigned_0_...2.psbt                  │
└─────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────┐
│ SIGN WALLET 2 (Parallel)                                    │
│ 1. sign musig2 nonce --file payment_15_unsigned_0_...2.psbt│
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. GenerateMuSig2NonceUseCase (Sign)                       │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. MuSig2Service.GenerateNonce()                            │
│    Output: [66]byte nonce                                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. Store in PSBT                                            │
│    - Sign nonce: field_id = "musig2_nonce_sign2"           │
│    Output: payment_15_nonce_0_...3.psbt                     │
│    Status: All nonces collected (3/3)                       │
└─────────────────────────────────────────────────────────────┘
```

#### Round 2: Signing (Sequential)

```
┌─────────────────────────────────────────────────────────────┐
│ KEYGEN WALLET                                               │
│ 1. keygen musig2 sign --file payment_15_nonce_0.psbt       │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. MuSig2SignUseCase (Keygen)                              │
│    - Validate all nonces present                            │
│    - Extract nonces from PSBT                               │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. AccountKeyRepository                                     │
│    - Get private key for signing                            │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. MuSig2Service.CreatePartialSignature()                  │
│    Input:                                                    │
│    - privateKey (keygen)                                    │
│    - allPublicKeys                                          │
│    - allNonces [nonce1, nonce2, nonce3]                     │
│    - messageHash (transaction hash)                         │
│    Process:                                                  │
│    - Create MuSig2 context with all signers                 │
│    - Register all nonces                                    │
│    - Sign message hash                                      │
│    Output: partialSignature (32 bytes)                      │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 5. Store in PSBT                                            │
│    - Add partial signature to PSBT proprietary field        │
│    - Field ID: "musig2_partialsig_keygen"                  │
│    Output: payment_15_unsigned_1.psbt                       │
└─────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────┐
│ SIGN WALLET 1 (Sequential after Keygen)                    │
│ 1. sign musig2 sign --file payment_15_unsigned_1.psbt      │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. MuSig2SignUseCase (Sign)                                │
│    - Validate all nonces present                            │
│    - Extract nonces from PSBT                               │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. AuthAccountKeyRepository                                 │
│    - Get auth private key for signing                       │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. MuSig2Service.CreatePartialSignature()                  │
│    Input:                                                    │
│    - privateKey (sign1/auth1)                               │
│    - allPublicKeys                                          │
│    - allNonces                                              │
│    - messageHash                                            │
│    Output: partialSignature (32 bytes)                      │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 5. Store in PSBT                                            │
│    - Field ID: "musig2_partialsig_sign1"                   │
│    Output: payment_15_unsigned_2.psbt                       │
└─────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────┐
│ SIGN WALLET 2 (Sequential after Sign1)                     │
│ 1. sign musig2 sign --file payment_15_unsigned_2.psbt      │
└────────────────────┬────────────────────────────────────────┘
                     │
                 [Same process as Sign1]
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 5. Store in PSBT                                            │
│    - Field ID: "musig2_partialsig_sign2"                   │
│    Output: payment_15_unsigned_3.psbt                       │
│    Status: All partial signatures collected (3/3)           │
└─────────────────────────────────────────────────────────────┘
```

#### Aggregation (Watch Wallet)

```
┌─────────────────────────────────────────────────────────────┐
│ WATCH WALLET                                                │
│ 1. watch musig2 aggregate --file payment_15_unsigned_3.psbt│
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 2. AggregateMuSig2SignaturesUseCase                        │
│    - Validate all partial signatures present                │
│    - Extract partial signatures from PSBT                   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 3. MuSig2Service.AggregateSignatures()                     │
│    Input:                                                    │
│    - allPublicKeys                                          │
│    - allNonces                                              │
│    - partialSignatures [sig1, sig2, sig3]                  │
│    - messageHash                                            │
│    Process:                                                  │
│    - Create MuSig2 context                                  │
│    - Register all nonces                                    │
│    - Combine all partial signatures                         │
│    - Generate final signature                               │
│    Output: finalSignature (64 bytes Schnorr)               │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 4. MuSig2Service.VerifyAggregatedSignature()               │
│    - Verify signature against aggregated public key        │
│    - Ensure signature is valid before broadcasting          │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 5. Finalize PSBT                                            │
│    - Add final signature to PSBT witness                    │
│    - Remove temporary MuSig2 fields                         │
│    Output: payment_15_signed_3.psbt                         │
│    Status: Ready for broadcasting                           │
└─────────────────────────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│ 6. Broadcast Transaction                                    │
│    watch send --file payment_15_signed_3.psbt              │
│    → Transaction broadcast to Bitcoin network               │
└─────────────────────────────────────────────────────────────┘
```

---

## Data Flow

### Nonce Data Flow

```
Round 1: Nonce Generation
=========================

Keygen Wallet:
  Private Key → MuSig2Service.GenerateNonce()
                       ↓
                [66]byte Public Nonce
                       ↓
              PSBT Proprietary Field
                       ↓
         payment_15_unsigned_0_...1.psbt

Sign Wallet 1:
  Private Key → MuSig2Service.GenerateNonce()
                       ↓
                [66]byte Public Nonce
                       ↓
              PSBT Proprietary Field
                       ↓
         payment_15_unsigned_0_...2.psbt

Sign Wallet 2:
  Private Key → MuSig2Service.GenerateNonce()
                       ↓
                [66]byte Public Nonce
                       ↓
              PSBT Proprietary Field
                       ↓
          payment_15_nonce_0_...3.psbt
          (All nonces collected)
```

### Signature Data Flow

```
Round 2: Partial Signature Creation
====================================

Each Signer:
  Private Key + All Nonces + Message Hash
              ↓
  MuSig2Service.CreatePartialSignature()
              ↓
  [32]byte Partial Signature
              ↓
  PSBT Proprietary Field
              ↓
  Updated PSBT file

Aggregation:
  All Partial Signatures + All Nonces + Message Hash
              ↓
  MuSig2Service.AggregateSignatures()
              ↓
  [64]byte Schnorr Signature
              ↓
  PSBT Witness Field
              ↓
  Finalized PSBT
              ↓
  Broadcast to Bitcoin Network
```

### PSBT Field Structure

MuSig2 uses PSBT proprietary fields (BIP174 extension):

```
Proprietary Field Format:
Key:   <identifier> <subtype> <keydata>
Value: <valuedata>

MuSig2 Nonces:
Key:   FC <musig2> <nonce> <signer_id>
Value: [66 bytes] public nonce

MuSig2 Partial Signatures:
Key:   FC <musig2> <psig> <signer_id>
Value: [32 bytes] partial signature

Example PSBT with MuSig2 data:
┌─────────────────────────────────────┐
│ PSBT Global Fields                  │
│ - Version                           │
│ - Transaction                       │
└─────────────────────────────────────┘
┌─────────────────────────────────────┐
│ PSBT Input #0                       │
│ - Non-witness UTXO                  │
│ - Witness UTXO                      │
│ - Derivation paths                  │
│ - Proprietary Fields:               │
│   * musig2_nonce_keygen: [66]byte  │
│   * musig2_nonce_sign1:  [66]byte  │
│   * musig2_nonce_sign2:  [66]byte  │
│   * musig2_psig_keygen:  [32]byte  │
│   * musig2_psig_sign1:   [32]byte  │
│   * musig2_psig_sign2:   [32]byte  │
└─────────────────────────────────────┘
┌─────────────────────────────────────┐
│ PSBT Output #0 (Recipient)          │
│ - Script                            │
│ - Amount                            │
└─────────────────────────────────────┘
┌─────────────────────────────────────┐
│ PSBT Output #1 (Change)             │
│ - Script                            │
│ - Amount                            │
└─────────────────────────────────────┘
```

---

## Security Architecture

### Nonce Security

**Critical Requirement**: Each nonce must be used exactly once. Nonce reuse will leak the private key!

#### Multi-Layer Protection

**1. Application Layer**:

```go
// Validate nonce uniqueness before use
func ValidateNonceUniqueness(nonces []NonceCommitment) error {
    seen := make(map[string]bool)
    for _, nonce := range nonces {
        key := string(nonce.Nonce[:])
        if seen[key] {
            return ErrDuplicateNonce
        }
        seen[key] = true
    }
    return nil
}
```

**2. Database Layer** (Future Enhancement):

```sql
CREATE TABLE musig2_nonces (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    signer_id VARCHAR(255) NOT NULL,
    nonce BINARY(66) NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_nonce (nonce),
    UNIQUE KEY unique_tx_signer (transaction_id, signer_id)
);
```

**3. Cryptographic Layer**:

- Nonces generated using secure random number generator
- Each session creates fresh nonces
- Nonces are automatically deleted after use in btcd library

#### Nonce Lifecycle

```
1. Generate → Fresh random nonce
2. Store    → PSBT proprietary field + DB tracking
3. Exchange → Via PSBT file transfers
4. Use      → Sign message once
5. Delete   → Remove after signing complete
```

### Key Aggregation Security

**Public Key Validation**:

```go
// Validate all public keys before aggregation
func ValidatePublicKeysForMuSig2(pubKeys [][]byte) error {
    // Check count (2-15 signers)
    if len(pubKeys) < 2 || len(pubKeys) > 15 {
        return ErrInvalidSignerCount
    }

    // Check each key
    for i, pubKey := range pubKeys {
        if len(pubKey) == 0 {
            return fmt.Errorf("empty public key at index %d", i)
        }
        // Validate key format (33 or 65 bytes)
        if len(pubKey) != 33 && len(pubKey) != 65 {
            return fmt.Errorf("invalid key length at index %d", i)
        }
    }

    // Check for duplicates
    seen := make(map[string]bool)
    for _, pubKey := range pubKeys {
        key := string(pubKey)
        if seen[key] {
            return ErrDuplicatePublicKey
        }
        seen[key] = true
    }

    return nil
}
```

**Taproot Tweak Application**:

```go
// Apply BIP86 Taproot tweak for key-path spending
aggregatedKey, err := musig2Service.AggregatePublicKeys(
    pubKeys,
    true, // applyTaprootTweak
)
```

### Signature Validation

**Partial Signature Verification**:

```go
// Each partial signature is validated before aggregation
func ValidatePartialSignatures(sigs []PartialSignature, expected int) error {
    if len(sigs) != expected {
        return fmt.Errorf("expected %d signatures, got %d", expected, len(sigs))
    }

    for i, sig := range sigs {
        if len(sig.Signature) != 32 {
            return fmt.Errorf("invalid signature length at index %d", i)
        }
        if sig.SignerID == "" {
            return fmt.Errorf("missing signer ID at index %d", i)
        }
    }

    return nil
}
```

**Final Signature Verification**:

```go
// Verify aggregated signature before broadcasting
isValid := musig2Service.VerifyAggregatedSignature(
    aggregatedPubKey,
    messageHash,
    finalSignature,
)
if !isValid {
    return errors.New("aggregated signature verification failed")
}
```

### Attack Vectors and Mitigations

#### 1. Nonce Reuse Attack

**Attack**: Reusing the same nonce for different messages allows an attacker to compute the private key.

**Mitigation**:

- Database unique constraint on nonces
- Application-level validation
- Automatic nonce deletion after use
- Session-based nonce management

#### 2. Rogue Key Attack

**Attack**: Attacker provides a crafted public key that allows them to control the aggregated key.

**Mitigation**:

- MuSig2 protocol includes built-in rogue key protection
- Key sorting ensures deterministic aggregation
- All signers must participate in signing

#### 3. Signature Forgery

**Attack**: Attempt to create valid signature without all required partial signatures.

**Mitigation**:

- Cryptographic proof prevents forgery without all partial signatures
- Verification step before broadcasting
- PSBT validation at each stage

#### 4. PSBT Tampering

**Attack**: Modify PSBT fields to change transaction or signatures.

**Mitigation**:

- PSBT format includes checksums
- Transaction hash validation at each step
- Offline signing prevents network attacks

---

## Database Schema

### Current Tables (Used by MuSig2)

#### account_key Table

Stores account keys including MuSig2 Taproot addresses:

```sql
CREATE TABLE account_key (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,
    account VARCHAR(255) NOT NULL,
    key_index INT NOT NULL,
    full_public_key TEXT NOT NULL,
    wallet_import_format TEXT,
    p2pkh_address VARCHAR(255),
    p2sh_segwit_address VARCHAR(255),
    bech32_address VARCHAR(255),
    multisig_address VARCHAR(255),     -- MuSig2 P2TR address stored here
    redeem_script TEXT,
    addr_status TINYINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_coin_account (coin, account),
    KEY idx_addr_status (addr_status)
);
```

**MuSig2 Usage**:

- `multisig_address`: Stores P2TR (Taproot) address (`bc1p...` or `tb1p...`)
- `addr_status`: Tracks address creation status
- `full_public_key`: Account's public key for aggregation

#### auth_fullpubkey Table

Stores public keys from Sign wallets (auth accounts):

```sql
CREATE TABLE auth_fullpubkey (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,
    auth_account VARCHAR(255) NOT NULL,
    full_public_key TEXT NOT NULL,
    p2pkh_address VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_coin_auth (coin, auth_account)
);
```

**MuSig2 Usage**:

- `full_public_key`: Sign wallet's public key for key aggregation
- `auth_account`: Identifies which Sign wallet (auth1, auth2, etc.)

#### auth_account_key Table

Stores Sign wallet's private keys for signing:

```sql
CREATE TABLE auth_account_key (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,
    auth_account VARCHAR(255) NOT NULL,
    full_public_key TEXT NOT NULL,
    wallet_import_format TEXT NOT NULL,
    p2pkh_address VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_coin_auth (coin, auth_account)
);
```

**MuSig2 Usage**:

- `wallet_import_format`: Private key for creating partial signatures
- Used by Sign wallets during Round 2 (signing)

### Future Tables (For Enhanced Nonce Management)

#### musig2_nonces Table (Proposed)

Track nonces to prevent reuse:

```sql
CREATE TABLE musig2_nonces (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,
    transaction_id BIGINT NOT NULL,
    signer_id VARCHAR(255) NOT NULL,
    nonce BINARY(66) NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    used_at TIMESTAMP NULL,
    UNIQUE KEY unique_nonce (nonce),
    UNIQUE KEY unique_tx_signer (transaction_id, signer_id),
    KEY idx_coin_tx (coin, transaction_id),
    KEY idx_used (used)
);
```

**Purpose**:

- Enforce nonce uniqueness at database level
- Track which nonces have been used
- Prevent nonce reuse attacks
- Audit trail for debugging

#### musig2_sessions Table (Proposed)

Track signing sessions:

```sql
CREATE TABLE musig2_sessions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL UNIQUE,
    coin VARCHAR(10) NOT NULL,
    transaction_id BIGINT NOT NULL,
    participant_count INT NOT NULL,
    nonces_collected INT DEFAULT 0,
    partial_sigs_collected INT DEFAULT 0,
    status ENUM('nonce_generation', 'signing', 'aggregating', 'completed', 'failed') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    KEY idx_coin_tx (coin, transaction_id),
    KEY idx_status (status)
);
```

**Purpose**:

- Track signing session progress
- Validate all nonces/signatures collected
- Monitor signing workflow state
- Debugging and audit trail

---

## API Reference

### Use Case Interfaces

#### CreateMuSig2AddressUseCase

```go
package keygenusecase

type CreateMuSig2AddressUseCase interface {
    Create(ctx context.Context, input CreateMuSig2AddressInput) error
}

type CreateMuSig2AddressInput struct {
    AccountType account.AccountType
}
```

**Example Usage**:

```go
useCase := container.NewKeygenCreateMuSig2AddressUseCase()
err := useCase.Create(ctx, keygenusecase.CreateMuSig2AddressInput{
    AccountType: domainAccount.AccountTypePayment,
})
```

#### GenerateMuSig2NonceUseCase

```go
package keygenusecase  // or signusecase

type GenerateMuSig2NonceUseCase interface {
    Generate(ctx context.Context, input GenerateMuSig2NonceInput) (*GenerateMuSig2NonceOutput, error)
}

type GenerateMuSig2NonceInput struct {
    TransactionID  int64
    PSBTData       []byte
    AccountType    account.AccountType
}

type GenerateMuSig2NonceOutput struct {
    Nonce     [66]byte
    PSBTData  []byte
}
```

**Example Usage**:

```go
useCase := container.NewKeygenGenerateMuSig2NonceUseCase()
output, err := useCase.Generate(ctx, keygenusecase.GenerateMuSig2NonceInput{
    TransactionID: 15,
    PSBTData:      psbtBytes,
    AccountType:   domainAccount.AccountTypePayment,
})
```

#### MuSig2SignUseCase

```go
package keygenusecase  // or signusecase

type MuSig2SignUseCase interface {
    Sign(ctx context.Context, input MuSig2SignInput) (*MuSig2SignOutput, error)
}

type MuSig2SignInput struct {
    PSBTData       []byte
    AccountType    account.AccountType
}

type MuSig2SignOutput struct {
    PartialSignature []byte
    PSBTData         []byte
}
```

**Example Usage**:

```go
useCase := container.NewKeygenMuSig2SignUseCase()
output, err := useCase.Sign(ctx, keygenusecase.MuSig2SignInput{
    PSBTData:    psbtBytesWithNonces,
    AccountType: domainAccount.AccountTypePayment,
})
```

#### AggregateMuSig2SignaturesUseCase

```go
package watchusecase

type AggregateMuSig2SignaturesUseCase interface {
    Aggregate(ctx context.Context, input AggregateMuSig2SignaturesInput) (*AggregateMuSig2SignaturesOutput, error)
}

type AggregateMuSig2SignaturesInput struct {
    PSBTData []byte
}

type AggregateMuSig2SignaturesOutput struct {
    FinalSignature [64]byte
    PSBTData       []byte
}
```

**Example Usage**:

```go
useCase := container.NewWatchAggregateMuSig2SignaturesUseCase()
output, err := useCase.Aggregate(ctx, watchusecase.AggregateMuSig2SignaturesInput{
    PSBTData: psbtBytesWithPartialSigs,
})
```

### Service Interfaces

#### MuSig2Service

```go
package btc

type MuSig2Service struct {
    chainConfig *chaincfg.Params
}

func NewMuSig2Service(chainConfig *chaincfg.Params) *MuSig2Service

// AggregatePublicKeys aggregates multiple public keys using MuSig2
func (s *MuSig2Service) AggregatePublicKeys(
    pubKeys []*btcec.PublicKey,
    applyTaprootTweak bool,
) (*btcec.PublicKey, error)

// GenerateNonce generates a secure nonce for MuSig2 signing
func (s *MuSig2Service) GenerateNonce(
    privateKey *btcec.PrivateKey,
    pubKeys []*btcec.PublicKey,
) ([66]byte, error)

// CreatePartialSignature creates a partial signature
func (s *MuSig2Service) CreatePartialSignature(
    privateKey *btcec.PrivateKey,
    pubKeys []*btcec.PublicKey,
    nonces [][66]byte,
    messageHash [32]byte,
) ([]byte, error)

// AggregateSignatures combines partial signatures
func (s *MuSig2Service) AggregateSignatures(
    pubKeys []*btcec.PublicKey,
    nonces [][66]byte,
    partialSigs [][]byte,
    messageHash [32]byte,
) (*schnorr.Signature, error)

// VerifyAggregatedSignature verifies the final signature
func (s *MuSig2Service) VerifyAggregatedSignature(
    aggregatedPubKey *btcec.PublicKey,
    messageHash [32]byte,
    signature *schnorr.Signature,
) bool
```

---

## Implementation Notes

### Library Integration: btcd/btcec/v2/schnorr/musig2

**Version**: v2.3.6
**Documentation**: <https://pkg.go.dev/github.com/btcsuite/btcd/btcec/v2/schnorr/musig2>
**License**: ISC (permissive)

#### Key Features Used

**1. MuSig2 Context Creation**:

```go
ctx, err := musig2.NewContext(
    privateKey,           // Signer's private key
    true,                 // Sort keys (MUST be same for all signers)
    musig2.WithKnownSigners(allPublicKeys),  // All participants
    musig2.WithBip86TweakCtx(),              // Taproot tweak for BIP86
)
```

**2. Session Management**:

```go
session, err := ctx.NewSession()
pubNonce := session.PublicNonce()  // [66]byte public nonce
```

**3. Nonce Registration**:

```go
haveAllNonces, err := session.RegisterPubNonce(otherNonce)
```

**4. Signing**:

```go
partialSig, err := session.Sign(messageHash)
```

**5. Aggregation**:

```go
haveAllSigs, err := session.CombineSig(partialSig)
finalSig := session.FinalSig()
```

**6. Verification**:

```go
isValid := finalSig.Verify(messageHash[:], aggregatedPubKey)
```

### PSBT Handling

**PSBT Library**: `github.com/btcsuite/btcd/psbt`

**MuSig2 Extension Fields**:

MuSig2 data is stored in PSBT proprietary fields (BIP174 allows custom fields):

```go
// Proprietary field format
type ProprietaryData struct {
    Identifier []byte  // "musig2"
    Subtype    []byte  // "nonce" or "psig"
    KeyData    []byte  // signer_id
    ValueData  []byte  // nonce or partial signature
}

// Example: Store nonce
func AddMuSig2NonceToPSBT(p *psbt.Packet, signerID string, nonce [66]byte) error {
    prop := psbt.ProprietaryData{
        Identifier: []byte("musig2"),
        Subtype:    []byte("nonce"),
        KeyData:    []byte(signerID),
        ValueData:  nonce[:],
    }
    p.Inputs[0].Unknowns = append(p.Inputs[0].Unknowns, &prop)
    return nil
}

// Example: Extract nonces
func ExtractMuSig2NoncesFromPSBT(p *psbt.Packet) ([][66]byte, error) {
    var nonces [][66]byte
    for _, unknown := range p.Inputs[0].Unknowns {
        if bytes.Equal(unknown.Identifier, []byte("musig2")) &&
           bytes.Equal(unknown.Subtype, []byte("nonce")) {
            var nonce [66]byte
            copy(nonce[:], unknown.ValueData)
            nonces = append(nonces, nonce)
        }
    }
    return nonces, nil
}
```

### Error Handling

**Error Wrapping**:

```go
if err != nil {
    return fmt.Errorf("failed to create MuSig2 context: %w", err)
}
```

**Domain Errors**:

```go
var (
    ErrDuplicateNonce       = errors.New("duplicate nonce detected")
    ErrInvalidSignerCount   = errors.New("invalid number of signers")
    ErrMissingPartialSig    = errors.New("missing partial signature")
    ErrSignatureVerifyFailed = errors.New("signature verification failed")
)
```

### Logging

**Structured Logging**:

```go
logger.Debug("create musig2 taproot address",
    "account_type", input.AccountType.String(),
)

logger.Info("MuSig2 signatures aggregated",
    "signer_count", len(partialSigs),
    "transaction_id", txID,
)

logger.Error("failed to aggregate signatures",
    "error", err,
    "partial_sig_count", len(partialSigs),
)
```

**Security**: Never log private keys, nonces (secret part), or sensitive data.

### Testing Strategy

**Unit Tests**:

- Test each use case in isolation
- Mock external dependencies
- Test error conditions

**Integration Tests**:

- Test complete signing workflows
- Test with real PSBT data
- Test nonce uniqueness enforcement

**Test Files**:

```
internal/application/usecase/keygen/btc/
├── create_musig2_address_test.go
├── musig2_nonce_test.go
└── musig2_sign_test.go

internal/application/usecase/sign/btc/
├── musig2_nonce_test.go
└── musig2_sign_test.go

internal/application/usecase/watch/btc/
└── musig2_aggregate_test.go

internal/infrastructure/api/btc/btc/
└── musig2_test.go
```

---

## Performance Considerations

### Transaction Size Comparison

| Type | Size | Reduction |
|------|------|-----------|
| Traditional 2-of-3 P2WSH | ~370 bytes | Baseline |
| MuSig2 3-of-3 P2TR | ~215 bytes | 41.9% |

**Cost Savings**:

- At 10 sat/vB: 1,550 sats saved per transaction
- At 50 sat/vB: 7,750 sats saved per transaction

### Parallel vs Sequential Operations

**Parallel** (Round 1 - Nonce Generation):

- All signers can generate nonces simultaneously
- No dependencies between nonce generations
- Reduces overall workflow time

**Sequential** (Round 2 - Signing):

- Must wait for all nonces before signing
- Each signer creates partial signature independently
- Can still be done in parallel after nonces collected

### Database Performance

**Nonce Uniqueness Check**:

- Index on nonce column for fast lookups
- Unique constraint prevents duplicates
- Consider cleanup of old nonces

---

## Migration from Traditional Multisig

### Compatibility

- **Backward Compatible**: Traditional P2WSH multisig continues to work
- **Opt-In**: Users choose when to migrate to MuSig2
- **No Breaking Changes**: Existing addresses and transactions unaffected

### Migration Path

1. **Create MuSig2 addresses** for new deposits
2. **Keep existing P2WSH addresses** for backward compatibility
3. **Gradually transition** funds to MuSig2 addresses
4. **Monitor savings** in transaction fees

---

## References

### Bitcoin Improvement Proposals (BIPs)

- [BIP 327: MuSig2](https://github.com/bitcoin/bips/blob/master/bip-0327.mediawiki)
- [BIP 340: Schnorr Signatures](https://github.com/bitcoin/bips/blob/master/bip-0340.mediawiki)
- [BIP 341: Taproot](https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki)
- [BIP 86: Key Derivation for Taproot](https://github.com/bitcoin/bips/blob/master/bip-0086.mediawiki)
- [BIP 174: PSBT Specification](https://github.com/bitcoin/bips/blob/master/bip-0174.mediawiki)

### Research Papers

- [MuSig2: Simple Two-Round Schnorr Multisignatures](https://eprint.iacr.org/2020/1261)
- [MuSig: A New Multisignature Standard](https://eprint.iacr.org/2018/068)

### Library Documentation

- [btcd/btcec/v2/schnorr/musig2](https://pkg.go.dev/github.com/btcsuite/btcd/btcec/v2/schnorr/musig2)
- [btcd/psbt](https://pkg.go.dev/github.com/btcsuite/btcd/psbt)

### Internal Documentation

- [MuSig2 User Guide](/chains/btc/musig2/user-guide)
- MuSig2 Library Selection (internal research document)
- MuSig2 Usage Guide (internal research document)

---

**Last Updated**: 2025-12-30
**Version**: 1.0
**Authors**: Development Team with Claude Code
