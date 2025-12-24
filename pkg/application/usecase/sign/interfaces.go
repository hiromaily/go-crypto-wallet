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
