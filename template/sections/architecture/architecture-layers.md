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
