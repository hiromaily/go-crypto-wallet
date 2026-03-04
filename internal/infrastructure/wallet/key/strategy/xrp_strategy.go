package strategy

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160" //nolint:gosec

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	xrpaddr "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
)

// XRPKeyStrategy implements CoinKeyStrategy for Ripple (XRP).
//
// XRP uses account-based addresses with a unique encoding scheme.
// XRP addresses are Base58Check encoded with specific prefix bytes.
//
// Address Generation:
// 1. Private key → Compressed public key
// 2. SHA256(RIPEMD160(public key)) → 20-byte hash
// 3. Base58Check encode with XRP prefix → Address (e.g., rN7n7otQDd6FczFgLdlqtyMVrn3P8vKV1Y)
//
// Special Notes:
// - XRP addresses start with 'r' for mainnet
// - P2SHSegWitAddr field is reused to store Ethereum address (used as passphrase for API operations)
// - This is a workaround for generating deterministic keys via xrpl-grpc-server API
type XRPKeyStrategy struct{}

// Compile-time check that XRPKeyStrategy implements CoinKeyStrategy
var _ portsWallet.CoinKeyStrategy = (*XRPKeyStrategy)(nil)

// NewXRPKeyStrategy creates a new XRPKeyStrategy
func NewXRPKeyStrategy() *XRPKeyStrategy {
	return &XRPKeyStrategy{}
}

// CoinTypeCode returns the coin type code for Ripple
func (*XRPKeyStrategy) CoinTypeCode() domainCoin.CoinTypeCode {
	return domainCoin.XRP
}

// GenerateWalletKey generates XRP addresses from a private key.
// Note: P2SHSegWitAddr is used to store an Ethereum address, which serves as
// a passphrase for API operations with xrpl-grpc-server.
func (s *XRPKeyStrategy) GenerateWalletKey(privKey *btcec.PrivateKey) (*domainKey.WalletKey, error) {
	// Generate XRP address, public key, and private key
	xrpAddr, xrpPubKey, xrpPrivKey, err := s.generateXRPKeys(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate XRP keys: %w", err)
	}

	// Generate Ethereum address (used as passphrase for API `wallet_propose`)
	ethAddr, err := s.generateEthAddress(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ETH address for XRP: %w", err)
	}

	return &domainKey.WalletKey{
		WIF:            xrpPrivKey,
		P2PKHAddr:      xrpAddr,
		P2SHSegWitAddr: ethAddr, // ETH address used as passphrase
		Bech32Addr:     "",      // Not used for XRP
		TaprootAddr:    "",      // Not used for XRP
		FullPubKey:     xrpPubKey,
		RedeemScript:   "", // Not used for XRP
	}, nil
}

// generateXRPKeys generates XRP address, public key, and hex-encoded private key bytes.
// The private key is returned as a 32-byte hex-encoded string for storage in the WIF field,
// allowing repositories to recover raw bytes via hex.DecodeString.
func (*XRPKeyStrategy) generateXRPKeys(privKey *btcec.PrivateKey) (string, string, string, error) {
	// Convert btcec private key to ECDSA
	xrpPrivKey := privKey.ToECDSA()

	// Store raw private key bytes as hex in the WIF field
	privKeyHex := hex.EncodeToString(crypto.FromECDSA(xrpPrivKey))

	// Generate public key hash
	serializedPubKey := privKey.PubKey().SerializeCompressed()
	pubKeyHash := xrpaddr.Sha256RipeMD160(serializedPubKey)
	if len(pubKeyHash) != ripemd160.Size {
		return "", "", "", errors.New("pubKeyHash must be 20 bytes")
	}

	// Generate XRP address
	address, err := xrpaddr.NewAccountID(pubKeyHash)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to call xrpaddr.NewAccountID(): %w", err)
	}

	// Generate XRP public key
	publicKey, err := xrpaddr.NewAccountPublicKey(pubKeyHash)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to call xrpaddr.NewAccountPublicKey(): %w", err)
	}

	return address.String(), publicKey.String(), privKeyHex, nil
}

// generateEthAddress generates an Ethereum address from the same private key
// This is used as a passphrase for XRP API operations
func (*XRPKeyStrategy) generateEthAddress(privKey *btcec.PrivateKey) (string, error) {
	// Convert to ECDSA private key
	ethPrivKey := privKey.ToECDSA()

	// Get public key
	ethPubkey := ethPrivKey.Public()
	pubkeyECDSA, ok := ethPubkey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("failed to cast public key to ECDSA public key")
	}

	// Generate Ethereum address
	address := crypto.PubkeyToAddress(*pubkeyECDSA).Hex()

	return address, nil
}
