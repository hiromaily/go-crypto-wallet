package sign

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// SignTransactionUseCase signs transactions
type SignTransactionUseCase interface {
	Sign(ctx context.Context, input SignTransactionInput) (SignTransactionOutput, error)
}

// ImportPrivateKeyUseCase imports private keys (BTC only)
type ImportPrivateKeyUseCase interface {
	Import(ctx context.Context, input ImportPrivateKeyInput) error
}

// ExportFullPubkeyUseCase exports full public keys (BTC only)
type ExportFullPubkeyUseCase interface {
	Export(ctx context.Context) (ExportFullPubkeyOutput, error)
}

// GenerateSeedUseCase generates a new seed for auth accounts
type GenerateSeedUseCase interface {
	Generate(ctx context.Context) (GenerateSeedOutput, error)
}

// StoreSeedUseCase stores a seed for auth accounts
type StoreSeedUseCase interface {
	Store(ctx context.Context, input StoreSeedInput) (StoreSeedOutput, error)
}

// GenerateAuthKeyUseCase generates HD keys for auth accounts
type GenerateAuthKeyUseCase interface {
	Generate(ctx context.Context, input GenerateAuthKeyInput) (GenerateAuthKeyOutput, error)
}

// GenerateMuSig2NonceUseCase generates MuSig2 nonces for Round 1 (BTC only)
type GenerateMuSig2NonceUseCase interface {
	Generate(ctx context.Context, input GenerateMuSig2NonceInput) (GenerateMuSig2NonceOutput, error)
}

// MuSig2SignUseCase creates MuSig2 partial signatures for Round 2 (BTC only)
type MuSig2SignUseCase interface {
	Sign(ctx context.Context, input MuSig2SignInput) (MuSig2SignOutput, error)
}

// Input/Output DTOs

// SignTransactionInput represents input for signing a transaction
type SignTransactionInput struct {
	FilePath string
}

// SignTransactionOutput represents output from signing a transaction
type SignTransactionOutput struct {
	SignedData   string // Signed transaction data (hex for legacy, base64 PSBT for BTC)
	IsComplete   bool
	NextFilePath string
}

// ImportPrivateKeyInput represents input for importing private keys
// Note: AuthType is not needed here as it's already configured in the use case during construction
type ImportPrivateKeyInput struct {
	// Empty struct - AuthType is handled by the use case factory
}

// ExportFullPubkeyOutput represents output from exporting full public keys
type ExportFullPubkeyOutput struct {
	FileName string
}

// GenerateSeedOutput represents output from generating a seed
type GenerateSeedOutput struct {
	Seed []byte
}

// StoreSeedInput represents input for storing a seed
type StoreSeedInput struct {
	Seed string
}

// StoreSeedOutput represents output from storing a seed
type StoreSeedOutput struct {
	Seed []byte
}

// GenerateAuthKeyInput represents input for generating auth keys
type GenerateAuthKeyInput struct {
	AuthType domainAccount.AuthType
	Seed     []byte
	Count    uint32
}

// GenerateAuthKeyOutput represents output from generating auth keys
type GenerateAuthKeyOutput struct {
	GeneratedCount int
}

// GenerateMuSig2NonceInput represents input for generating MuSig2 nonces
type GenerateMuSig2NonceInput struct {
	TransactionID string
	AuthType      domainAccount.AuthType
}

// GenerateMuSig2NonceOutput represents output from generating MuSig2 nonces
type GenerateMuSig2NonceOutput struct {
	PublicNonce [66]byte
	SignerID    string
}

// MuSig2SignInput represents input for MuSig2 partial signature creation
type MuSig2SignInput struct {
	TransactionID    string
	AuthType         domainAccount.AuthType
	MessageHash      [32]byte
	AggregatedNonces [][66]byte // Public nonces from all signers
}

// MuSig2SignOutput represents output from MuSig2 partial signature creation
type MuSig2SignOutput struct {
	PartialSignature [32]byte
	SignerID         string
}
