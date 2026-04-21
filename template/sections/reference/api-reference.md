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
