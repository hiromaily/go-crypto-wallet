package keygen

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
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

// ExportFullPubkeyUseCase exports the account-level extended public key to a file (ETH only).
// The exported xpub allows the Watch wallet to derive and verify child addresses
// without holding private keys.
type ExportFullPubkeyUseCase interface {
	Export(ctx context.Context, input ExportFullPubkeyInput) (ExportFullPubkeyOutput, error)
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

// SignMultisigTransactionUseCase performs offline EIP-712 signing of a Safe multisig proposal.
// It reads the multisig file, verifies the safeTxHash, signs with the signer's private key,
// and writes the updated file. No network calls are made — suitable for air-gapped wallets.
// Used by both ETHKeygen and ETHSign wallet adapters.
type SignMultisigTransactionUseCase interface {
	Sign(ctx context.Context, input SignMultisigTransactionInput) (SignMultisigTransactionOutput, error)
}

// RunDKGUseCase runs the DKG ceremony for this node, saves the encrypted key shard, and
// returns the joint Ethereum address shared by all ceremony participants.
type RunDKGUseCase interface {
	Execute(ctx context.Context, input RunDKGInput) (RunDKGOutput, error)
}

// ServeMPCUseCase starts the MPC node gRPC daemon that accepts signing sessions from the
// Watch wallet coordinator. It blocks until ctx is cancelled (graceful shutdown).
type ServeMPCUseCase interface {
	Serve(ctx context.Context, input ServeMPCInput) error
}

// RunDKGInput carries the configuration for a DKG ceremony on this node.
type RunDKGInput struct {
	PartyID         string
	AllPartyIDs     []string
	Threshold       int
	ListenAddr      string // gRPC listen address for this node in P2P DKG mode (e.g., "localhost:9001")
	PeerAddrs       []string
	PreParamsPath   string
	ShardOutputPath string
	Passphrase      string
}

// RunDKGOutput carries the result of a successful DKG ceremony.
type RunDKGOutput struct {
	EthAddress string // Checksummed joint Ethereum address
	ShardPath  string // Filesystem path where the encrypted shard was saved
}

// ServeMPCInput carries the configuration for starting the MPC signing daemon.
type ServeMPCInput struct {
	ListenAddr  string
	ShardPath   string
	Passphrase  string
	PartyID     string
	AllPartyIDs []string
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

// ExportFullPubkeyInput represents input for exporting the account-level extended public key (ETH)
type ExportFullPubkeyInput struct {
	AccountType domainAccount.AccountType
}

// ExportFullPubkeyOutput represents output from exporting the account-level extended public key (ETH)
type ExportFullPubkeyOutput struct {
	FileName string
}

// ImportPrivateKeyInput represents input for importing private keys
type ImportPrivateKeyInput struct {
	AccountType domainAccount.AccountType
}

// CreateMultisigAddressInput represents input for creating multisig addresses
type CreateMultisigAddressInput struct {
	AccountType domainAccount.AccountType
	AddressType domainBTC.AddressType
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

// SignMultisigTransactionInput represents input for offline EIP-712 Safe multisig signing.
type SignMultisigTransactionInput struct {
	FilePath      string // path to the ETHMultisigTransactionFile JSON
	SignerAddress string // EIP-55 checksummed EOA address of this signer
}

// SignMultisigTransactionOutput represents the result of a successful multisig signing step.
type SignMultisigTransactionOutput struct {
	FilePath   string // path to the newly written file (counter incremented)
	IsComplete bool   // true when len(Signatures) >= Threshold after this signing
	SignCount  int    // number of signatures collected (including this one)
	Threshold  int    // threshold required for execution
}
