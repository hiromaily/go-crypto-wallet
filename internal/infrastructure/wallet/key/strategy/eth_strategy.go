package strategy

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// ETHKeyStrategy implements CoinKeyStrategy for Ethereum (ETH).
//
// Ethereum uses account-based addresses derived from ECDSA public keys.
// Unlike Bitcoin's UTXO model, Ethereum addresses are permanent account identifiers.
//
// Address Generation:
// 1. Private key → ECDSA public key
// 2. Public key → Keccak256 hash → Last 20 bytes
// 3. Prefix with "0x" → Ethereum address (e.g., 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb)
//
// Ethereum addresses are case-insensitive but use EIP-55 checksumming for error detection.
type ETHKeyStrategy struct{}

// Compile-time check that ETHKeyStrategy implements CoinKeyStrategy
var _ portsWallet.CoinKeyStrategy = (*ETHKeyStrategy)(nil)

// NewETHKeyStrategy creates a new ETHKeyStrategy
func NewETHKeyStrategy() *ETHKeyStrategy {
	return &ETHKeyStrategy{}
}

// CoinTypeCode returns the coin type code for Ethereum
func (*ETHKeyStrategy) CoinTypeCode() domainCoin.CoinTypeCode {
	return domainCoin.ETH
}

// GenerateWalletKey generates Ethereum addresses from a private key
// https://goethereumbook.org/wallet-generate/
func (s *ETHKeyStrategy) GenerateWalletKey(privKey *btcec.PrivateKey) (*domainKey.WalletKey, error) {
	// Convert btcec private key to ECDSA private key
	ethPrivKey := privKey.ToECDSA()
	ethHexPrivKey := hexutil.Encode(crypto.FromECDSA(ethPrivKey))

	// Get public key and address
	ethAddr, ethPubKey, err := s.generateAddressAndPubKey(ethPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ETH address: %w", err)
	}

	// For Ethereum, we use WIF field to store the private key
	// P2PKHAddr stores the Ethereum address
	// Other address fields (P2SHSegWitAddr, Bech32Addr, TaprootAddr) are not used
	return &domainKey.WalletKey{
		WIF:            ethHexPrivKey,
		P2PKHAddr:      ethAddr,
		P2SHSegWitAddr: "", // Not used for Ethereum
		Bech32Addr:     "", // Not used for Ethereum
		TaprootAddr:    "", // Not used for Ethereum
		FullPubKey:     ethPubKey,
		RedeemScript:   "", // Not used for Ethereum
	}, nil
}

// generateAddressAndPubKey generates Ethereum address and public key from ECDSA private key
func (*ETHKeyStrategy) generateAddressAndPubKey(ethPrivKey *ecdsa.PrivateKey) (string, string, error) {
	// Get public key
	ethPubkey := ethPrivKey.Public()
	pubkeyECDSA, ok := ethPubkey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("failed to cast public key to ECDSA public key")
	}

	// Public key (without the 0x04 prefix)
	ethHexPubKey := hexutil.Encode(crypto.FromECDSAPub(pubkeyECDSA))[4:]

	// Address (checksummed according to EIP-55)
	address := crypto.PubkeyToAddress(*pubkeyECDSA).Hex()

	return address, ethHexPubKey, nil
}
