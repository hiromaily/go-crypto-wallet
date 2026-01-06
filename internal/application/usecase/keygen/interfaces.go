package keygen

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
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

// CreateMuSig2AddressUseCase creates MuSig2 Taproot addresses (BTC only)
type CreateMuSig2AddressUseCase interface {
	Create(ctx context.Context, input CreateMuSig2AddressInput) error
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

// GenerateMuSig2NonceUseCase generates MuSig2 nonces for Round 1 (BTC only)
type GenerateMuSig2NonceUseCase interface {
	Generate(ctx context.Context, input GenerateMuSig2NonceInput) (GenerateMuSig2NonceOutput, error)
}

// MuSig2SignUseCase creates MuSig2 partial signatures for Round 2 (BTC only)
type MuSig2SignUseCase interface {
	Sign(ctx context.Context, input MuSig2SignInput) (MuSig2SignOutput, error)
}

// GenerateDescriptorUseCase generates descriptors for an account (single-sig or multisig)
type GenerateDescriptorUseCase interface {
	Generate(ctx context.Context, input GenerateDescriptorInput) (GenerateDescriptorOutput, error)
}

// ExportDescriptorUseCase exports descriptors to file
type ExportDescriptorUseCase interface {
	Export(ctx context.Context, input ExportDescriptorInput) (ExportDescriptorOutput, error)
}

// DescriptorFormat specifies output format
type DescriptorFormat string

// Descriptor format constants
const (
	DescriptorFormatText        DescriptorFormat = "text"
	DescriptorFormatJSON        DescriptorFormat = "json"
	DescriptorFormatBitcoinCore DescriptorFormat = "bitcoin-core"
)

// GenerateDescriptorInput represents the input for generating a descriptor.
type GenerateDescriptorInput struct {
	AccountType  domainAccount.AccountType
	AddressType  domainAddress.AddrType
	IsChange     bool
	RequiredSigs int // Optional for multisig; 0 selects the minimal required-sigs config
}

type GenerateDescriptorOutput struct {
	Descriptor  string
	AccountType domainAccount.AccountType
	AddressType domainAddress.AddrType
	IsMultisig  bool
}

// ExportDescriptorInput contains export parameters
type ExportDescriptorInput struct {
	AccountType   domainAccount.AccountType
	OutputPath    string
	Format        DescriptorFormat
	IncludeChange bool
}

// ExportDescriptorOutput contains export result
type ExportDescriptorOutput struct {
	FilePath string
}

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
	AddressType domainBitcoin.AddressType
}

// CreateMuSig2AddressInput represents input for creating MuSig2 Taproot addresses
type CreateMuSig2AddressInput struct {
	AccountType domainAccount.AccountType
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

// GenerateMuSig2NonceInput represents input for generating MuSig2 nonces
type GenerateMuSig2NonceInput struct {
	TransactionID string
	SignerID      string
}

// GenerateMuSig2NonceOutput represents output from generating MuSig2 nonces
type GenerateMuSig2NonceOutput struct {
	PublicNonce [66]byte
	SignerID    string
}

// MuSig2SignInput represents input for MuSig2 partial signature creation
type MuSig2SignInput struct {
	TransactionID    string
	SignerID         string
	MessageHash      [32]byte
	AggregatedNonces [][66]byte // Public nonces from all signers
}

// MuSig2SignOutput represents output from MuSig2 partial signature creation
type MuSig2SignOutput struct {
	PartialSignature [32]byte
	SignerID         string
}
