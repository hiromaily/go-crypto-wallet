package sign

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
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

// Input/Output DTOs

// SignTransactionInput represents input for signing a transaction
type SignTransactionInput struct {
	FilePath string
}

// SignTransactionOutput represents output from signing a transaction
type SignTransactionOutput struct {
	SignedHex    string
	IsComplete   bool
	NextFilePath string
}

// ImportPrivateKeyInput represents input for importing private keys
type ImportPrivateKeyInput struct {
	AuthType domainAccount.AuthType // For BTC, this is authType instead of accountType
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
