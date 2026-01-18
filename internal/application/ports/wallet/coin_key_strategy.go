// Package wallet defines interfaces for wallet key generation operations.
package wallet

import (
	"github.com/btcsuite/btcd/btcec/v2"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// CoinKeyStrategy defines the interface for coin-specific key generation logic.
// This interface enables separation of concerns between BIP derivation (handled by HDKey)
// and coin-specific address generation (handled by strategy implementations).
//
// Each cryptocurrency (BTC, BCH, ETH, XRP) has unique address generation requirements:
// - BTC: Supports P2PKH, P2SH-SegWit, Bech32, and Taproot addresses
// - BCH: Supports P2PKH (CashAddr format) only, no SegWit/Taproot
// - ETH: Account-based addresses from ECDSA public keys
// - XRP: Ripple-specific address encoding
//
// By implementing this strategy pattern, we:
// 1. Isolate coin-specific logic from BIP32 derivation
// 2. Prevent accidental cross-coin contamination (e.g., BTC logic affecting BCH)
// 3. Make it easier to add new cryptocurrencies
// 4. Improve testability with focused unit tests per coin
type CoinKeyStrategy interface {
	// GenerateWalletKey creates a WalletKey from a private key using coin-specific logic.
	// The private key has already been derived using BIP32 by the HDKey struct.
	// This method is responsible for:
	// - Generating all supported address types for this coin
	// - Encoding addresses in the correct format
	// - Populating WIF, public keys, and redeem scripts as needed
	//
	// Parameters:
	//   - privKey: BIP32-derived private key (from HDKey)
	//
	// Returns:
	//   - WalletKey populated with coin-specific addresses and metadata
	//   - Error if address generation fails
	GenerateWalletKey(privKey *btcec.PrivateKey) (*domainKey.WalletKey, error)

	// CoinTypeCode returns the coin type this strategy handles.
	// Used for logging and validation.
	CoinTypeCode() domainCoin.CoinTypeCode
}
