// Package wallet defines interfaces for wallet key generation operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer.
package wallet

import (
	"github.com/btcsuite/btcd/chaincfg"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// Generator is the key generator interface for creating HD wallet keys.
// Implementations are responsible for generating keys according to specific BIP standards
// (BIP44, BIP49, BIP84, BIP86, etc.).
type Generator interface {
	// KeyType returns the key type this generator supports (e.g., BIP44, BIP84).
	KeyType() domainKey.KeyType

	// CreateKey creates keys based on the seed and account type.
	// Parameters:
	//   - seed: The HD wallet seed (BIP39 mnemonic-derived)
	//   - actType: The account type (client, deposit, payment, etc.)
	//   - idxFrom: Starting index for key derivation
	//   - count: Number of keys to generate
	// Returns a slice of generated wallet keys with addresses and metadata.
	CreateKey(seed []byte, actType domainAccount.AccountType, idxFrom, count uint32) ([]domainKey.WalletKey, error)

	// SupportsAddressType checks if this generator supports the given address type.
	// Different generators support different address formats (P2PKH, P2WPKH, P2TR, etc.).
	SupportsAddressType(addrType domainAddress.AddrType) bool

	// GetDerivationPath returns the BIP32 derivation path for the given account and index.
	// Format: m/purpose'/coin_type'/account'/change/index
	GetDerivationPath(accountType domainAccount.AccountType, index uint32) string
}

// GeneratorFactory creates a Generator based on key type.
// This factory pattern allows for runtime selection of the appropriate key generator
// implementation based on the desired BIP standard.
type GeneratorFactory interface {
	// CreateGenerator creates a key generator for the specified key type.
	// Parameters:
	//   - keyType: The BIP standard to use (BIP44, BIP49, BIP84, BIP86, MuSig2)
	//   - coinTypeCode: The cryptocurrency (BTC, BCH, ETH, XRP)
	//   - conf: Bitcoin network parameters (mainnet, testnet, regtest)
	// Returns a Generator implementation or an error if the key type is unsupported.
	CreateGenerator(
		keyType domainKey.KeyType,
		coinTypeCode domainCoin.CoinTypeCode,
		conf *chaincfg.Params,
	) (Generator, error)
}
