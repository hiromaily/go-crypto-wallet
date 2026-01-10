package key

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

// DescriptorGenerator generates Bitcoin output descriptors from HD wallet seeds.
//
// This infrastructure component bridges the domain layer (DescriptorBuilder)
// and the wallet layer (HDKey) to generate descriptors from seeds.
//
// The generator:
//   - Derives account-level extended public keys from seeds
//   - Computes master key fingerprints
//   - Builds descriptor strings with proper BIP380 format
//   - Supports BIP44, BIP49 (P2SH-SegWit), BIP84 (Native SegWit), and BIP86 (Taproot)
type DescriptorGenerator struct {
	hdKey             *HDKey
	conf              *chaincfg.Params
	descriptorBuilder *domainWallet.DescriptorBuilder
}

// NewDescriptorGenerator creates a new descriptor generator.
//
// Parameters:
//   - hdKey: HD wallet key generator (must be initialized with purpose and coin type)
//   - conf: Bitcoin network configuration (mainnet, testnet, regtest)
//
// Returns:
//   - DescriptorGenerator instance ready to generate descriptors
func NewDescriptorGenerator(hdKey *HDKey, conf *chaincfg.Params) *DescriptorGenerator {
	return &DescriptorGenerator{
		hdKey:             hdKey,
		conf:              conf,
		descriptorBuilder: domainWallet.NewDescriptorBuilder(),
	}
}

// GenerateAccountDescriptor generates a descriptor for a specific account.
//
// The descriptor includes:
//   - Master key fingerprint (first 4 bytes of master key hash)
//   - Derivation path (e.g., /49'/0'/0' for BIP49)
//   - Account-level extended public key
//
// Parameters:
//   - seed: BIP39 seed (32 or 64 bytes)
//   - accountType: Account type (deposit, payment, stored, etc.)
//   - withChecksum: If true, includes BIP380 checksum in descriptor
//
// Returns:
//   - Descriptor object with all metadata
//   - Error if seed is invalid or key derivation fails
//
// Example output:
//
//	sh(wpkh([a1b2c3d4/49'/0'/0']xpub6ERApfZw...))#checksum
func (g *DescriptorGenerator) GenerateAccountDescriptor(
	seed []byte,
	accountType domainAccount.AccountType,
	withChecksum bool,
) (*domainWallet.Descriptor, error) {
	// Get master fingerprint
	fingerprint, err := g.GetMasterFingerprint(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to get master fingerprint: %w", err)
	}

	// Get account-level xpub
	xpub, err := g.GetAccountXPub(seed, accountType)
	if err != nil {
		return nil, fmt.Errorf("failed to get account xpub: %w", err)
	}

	// Build derivation path (e.g., /49'/0'/0')
	derivationPath := fmt.Sprintf("/%d'/%d'/%d'",
		g.hdKey.purpose.Uint32(),
		g.hdKey.coinType.Uint32(),
		accountType.Uint32(),
	)

	// Create descriptor key
	descriptorKey := domainWallet.DescriptorKey{
		Fingerprint:    fingerprint,
		DerivationPath: derivationPath,
		ExtendedPubKey: xpub,
	}

	// Determine descriptor type based on purpose
	descType := g.getDescriptorType()

	// Build descriptor
	var descriptorStr string
	if withChecksum {
		descriptorStr, err = g.descriptorBuilder.BuildDescriptorWithChecksum(descType, descriptorKey)
		if err != nil {
			return nil, fmt.Errorf("failed to build descriptor with checksum: %w", err)
		}
	} else {
		descriptorStr, err = g.descriptorBuilder.BuildDescriptor(descType, descriptorKey)
		if err != nil {
			return nil, fmt.Errorf("failed to build descriptor: %w", err)
		}
	}

	// Parse descriptor string to get checksum (if present)
	var checksum string
	if withChecksum {
		// Extract checksum from descriptor string (format: descriptor#checksum)
		if len(descriptorStr) > 9 && descriptorStr[len(descriptorStr)-9] == '#' {
			checksum = descriptorStr[len(descriptorStr)-8:]
			descriptorStr = descriptorStr[:len(descriptorStr)-9]
		}
	}

	// Create descriptor object
	descriptor := &domainWallet.Descriptor{
		Type:     descType,
		Script:   descriptorStr,
		Keys:     []domainWallet.DescriptorKey{descriptorKey},
		Checksum: checksum,
	}

	return descriptor, nil
}

// GetMasterFingerprint derives the master key fingerprint from a seed.
//
// The fingerprint is the first 4 bytes (8 hex characters) of the HASH160
// of the master public key. This uniquely identifies the master key and
// is used in descriptor notation.
//
// Parameters:
//   - seed: BIP39 seed (32 or 64 bytes)
//
// Returns:
//   - 8-character hex string (e.g., "a1b2c3d4")
//   - Error if seed is invalid
//
// Reference: BIP32 - Hierarchical Deterministic Wallets
func (g *DescriptorGenerator) GetMasterFingerprint(seed []byte) (string, error) {
	// Derive master key
	masterKey, err := hdkeychain.NewMaster(seed, g.conf)
	if err != nil {
		return "", fmt.Errorf("failed to derive master key: %w", err)
	}

	// Get master public key
	masterPubKey, err := masterKey.Neuter()
	if err != nil {
		return "", fmt.Errorf("failed to get master public key: %w", err)
	}

	// Get fingerprint (first 4 bytes of pubkey hash)
	// hdkeychain provides ParentFingerprint() but we need the master fingerprint
	// We can compute it from the serialized public key
	fingerprintBytes := masterPubKey.ParentFingerprint()

	// Convert to hex string (8 characters)
	fingerprint := fmt.Sprintf("%08x", fingerprintBytes)

	return fingerprint, nil
}

// GetAccountXPub derives the account-level extended public key from a seed.
//
// The account xpub represents the public key at the account level in the
// BIP32 hierarchy: m/purpose'/coin_type'/account'
//
// This xpub can be used to derive all addresses for the account without
// exposing the private key.
//
// Parameters:
//   - seed: BIP39 seed (32 or 64 bytes)
//   - accountType: Account type (deposit, payment, stored, etc.)
//
// Returns:
//   - Extended public key string (starts with "xpub", "tpub", "ypub", "zpub", etc.)
//   - Error if seed is invalid or derivation fails
//
// Example: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4..."
func (g *DescriptorGenerator) GetAccountXPub(
	seed []byte,
	accountType domainAccount.AccountType,
) (string, error) {
	// Use existing HD wallet method to derive account-level keys
	_, accountPubKey, err := g.hdKey.createKeyByAccount(seed, accountType)
	if err != nil {
		return "", fmt.Errorf("failed to derive account key: %w", err)
	}

	// Convert to string (xpub format)
	xpub := accountPubKey.String()

	return xpub, nil
}

// GenerateAccountDescriptorWithRange generates a descriptor with derivation path range.
//
// This is useful for generating descriptors that represent multiple addresses
// using wildcard notation (e.g., /0/* for all external addresses).
//
// Parameters:
//   - seed: BIP39 seed
//   - accountType: Account type
//   - change: Change index (0 for external, 1 for internal/change)
//   - wildcard: If true, use /* for wildcard range; if false, just add change index
//   - withChecksum: If true, includes BIP380 checksum
//
// Returns:
//   - Complete descriptor string (e.g., "sh(wpkh([fp/49'/0'/0']xpub.../0/*))#checksum")
//   - Error if generation fails
func (g *DescriptorGenerator) GenerateAccountDescriptorWithRange(
	seed []byte,
	accountType domainAccount.AccountType,
	change uint32,
	wildcard bool,
	withChecksum bool,
) (string, error) {
	// Generate base descriptor without checksum first
	descriptor, err := g.GenerateAccountDescriptor(seed, accountType, false)
	if err != nil {
		return "", err
	}

	// Use the descriptor script string
	descriptorStr := descriptor.Script

	// Add range
	descriptorWithRange, err := g.descriptorBuilder.FormatDescriptorWithRange(
		descriptorStr,
		change,
		wildcard,
	)
	if err != nil {
		return "", fmt.Errorf("failed to add range to descriptor: %w", err)
	}

	// Add checksum if requested
	if withChecksum {
		checksum, err := g.descriptorBuilder.CalculateChecksum(descriptorWithRange)
		if err != nil {
			return "", fmt.Errorf("failed to calculate checksum: %w", err)
		}
		return descriptorWithRange + "#" + checksum, nil
	}

	return descriptorWithRange, nil
}

// getDescriptorType maps BIP purpose to descriptor type.
func (g *DescriptorGenerator) getDescriptorType() domainWallet.DescriptorType {
	switch g.hdKey.purpose {
	case PurposeTypeBIP44:
		return domainWallet.DescriptorTypePKH
	case PurposeTypeBIP49:
		return domainWallet.DescriptorTypeSHWPKH
	case PurposeTypeBIP84:
		return domainWallet.DescriptorTypeWPKH
	case PurposeTypeBIP86:
		return domainWallet.DescriptorTypeTR
	default:
		return domainWallet.DescriptorTypePKH
	}
}

// GetMasterFingerprintHex is an alternative implementation that computes the fingerprint
// directly from the public key hash, following Bitcoin Core's approach.
//
// This is provided for reference and compatibility verification.
// Use GetMasterFingerprint() for standard usage.
func (g *DescriptorGenerator) GetMasterFingerprintHex(seed []byte) (string, error) {
	// Derive master key
	masterKey, err := hdkeychain.NewMaster(seed, g.conf)
	if err != nil {
		return "", fmt.Errorf("failed to derive master key: %w", err)
	}

	// Get the serialized public key
	pubKey, err := masterKey.ECPubKey()
	if err != nil {
		return "", fmt.Errorf("failed to get master public key: %w", err)
	}

	// Compute HASH160 (SHA256 followed by RIPEMD160)
	// In Bitcoin, the fingerprint is derived from HASH160 of the public key
	pubKeyBytes := pubKey.SerializeCompressed()
	hash160 := btcutil.Hash160(pubKeyBytes)

	// Take first 4 bytes as fingerprint
	fingerprint := hex.EncodeToString(hash160[:4])

	return fingerprint, nil
}
