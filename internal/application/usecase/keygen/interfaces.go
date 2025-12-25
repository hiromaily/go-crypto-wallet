package keygen

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
)

// GenerateHDWalletUseCase generates HD wallet keys
type GenerateHDWalletUseCase interface {
	Generate(ctx context.Context, input GenerateHDWalletInput) (GenerateHDWalletOutput, error)
}

// GenerateSeedUseCase generates and stores BIP39 seeds
type GenerateSeedUseCase interface {
	Generate(ctx context.Context) (GenerateSeedOutput, error)
	Store(ctx context.Context, input StoreSeedInput) (StoreSeedOutput, error)
}

// ExportAddressUseCase exports addresses to files
type ExportAddressUseCase interface {
	Export(ctx context.Context, input ExportAddressInput) (ExportAddressOutput, error)
}

// ImportPrivateKeyUseCase imports private keys
type ImportPrivateKeyUseCase interface {
	Import(ctx context.Context, input ImportPrivateKeyInput) error
}

// CreateMultisigAddressUseCase creates multisig addresses (BTC only)
type CreateMultisigAddressUseCase interface {
	Create(ctx context.Context, input CreateMultisigAddressInput) error
}

// ImportFullPubkeyUseCase imports full public keys from other signers (BTC only)
type ImportFullPubkeyUseCase interface {
	Import(ctx context.Context, input ImportFullPubkeyInput) error
}

// GenerateKeyUseCase generates keys (XRP only)
type GenerateKeyUseCase interface {
	Generate(ctx context.Context, input GenerateKeyInput) error
}

// SignTransactionUseCase signs unsigned transactions (first signature for multisig)
type SignTransactionUseCase interface {
	Sign(ctx context.Context, input SignTransactionInput) (SignTransactionOutput, error)
}

// Input/Output DTOs

// GenerateHDWalletInput represents input for generating HD wallet keys
type GenerateHDWalletInput struct {
	AccountType domainAccount.AccountType
	Seed        []byte
	Count       uint32
}

// GenerateHDWalletOutput represents output from generating HD wallet keys
type GenerateHDWalletOutput struct {
	GeneratedCount int
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

// ExportAddressInput represents input for exporting addresses
type ExportAddressInput struct {
	AccountType domainAccount.AccountType
}

// ExportAddressOutput represents output from exporting addresses
type ExportAddressOutput struct {
	FileName string
}

// ImportPrivateKeyInput represents input for importing private keys
type ImportPrivateKeyInput struct {
	AccountType domainAccount.AccountType
}

// CreateMultisigAddressInput represents input for creating multisig addresses
type CreateMultisigAddressInput struct {
	AccountType domainAccount.AccountType
	AddressType address.AddrType
}

// ImportFullPubkeyInput represents input for importing full public keys
type ImportFullPubkeyInput struct {
	FileName string
}

// GenerateKeyInput represents input for generating keys (XRP)
type GenerateKeyInput struct {
	AccountType domainAccount.AccountType
	IsKeyPair   bool
	WalletKeys  any // []domainKey.WalletKey - using any to avoid import cycle
}

// SignTransactionInput represents input for signing transactions
type SignTransactionInput struct {
	FilePath string
}

// SignTransactionOutput represents output from signing transactions
type SignTransactionOutput struct {
	FilePath      string
	IsDone        bool
	SignedCount   int
	UnsignedCount int
}
