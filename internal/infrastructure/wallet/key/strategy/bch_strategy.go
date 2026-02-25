package strategy

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	bchutil "github.com/hiromaily/go-crypto-wallet/pkg/chains/bch"
)

// BCHKeyStrategy implements CoinKeyStrategy for Bitcoin Cash (BCH).
//
// CRITICAL: BCH vs BTC Feature Differences
// Bitcoin Cash diverged from Bitcoin in 2017 and does NOT support SegWit or Taproot:
// - ✅ P2PKH (CashAddr format): The ONLY address type BCH supports
// - ❌ SegWit (P2SH-SegWit, Bech32): Not supported - BCH removed SegWit
// - ❌ Taproot: Not supported - BCH doesn't have Schnorr/Taproot
// - ❌ MuSig2: Not available without Schnorr signatures
//
// Address Format: CashAddr (prefix:payload)
// - Mainnet: bitcoincash:qpm2qsznhks23z7629mms6s4cwef74vcwvy22gdx6a
// - Testnet: bchtest:qpm2qsznhks23z7629mms6s4cwef74vcwvhanqgjxu
//
// This isolation ensures BCH address generation cannot be contaminated by BTC's
// SegWit/Taproot logic, which would result in invalid or unsupported addresses.
type BCHKeyStrategy struct {
	conf *chaincfg.Params
}

// Compile-time check that BCHKeyStrategy implements CoinKeyStrategy
var _ portsWallet.CoinKeyStrategy = (*BCHKeyStrategy)(nil)

// NewBCHKeyStrategy creates a new BCHKeyStrategy
func NewBCHKeyStrategy(conf *chaincfg.Params) *BCHKeyStrategy {
	return &BCHKeyStrategy{
		conf: conf,
	}
}

// CoinTypeCode returns the coin type code for Bitcoin Cash
func (*BCHKeyStrategy) CoinTypeCode() domainCoin.CoinTypeCode {
	return domainCoin.BCH
}

// GenerateWalletKey generates BCH addresses from a private key.
// BCH only supports P2PKH addresses in CashAddr format.
// All SegWit and Taproot fields are explicitly set to empty strings.
func (s *BCHKeyStrategy) GenerateWalletKey(privKey *btcec.PrivateKey) (*domainKey.WalletKey, error) {
	// WIF (Wallet Import Format) - compressed
	wif, err := btcutil.NewWIF(privKey, s.conf, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create WIF: %w", err)
	}

	// P2PKH address in CashAddr format (the ONLY address type BCH supports)
	p2pkhAddr, err := s.getP2PKHAddr(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate P2PKH CashAddr: %w", err)
	}

	// Full public key (compressed)
	fullPubKey := getFullPubKey(privKey, true)

	// IMPORTANT: BCH does NOT support SegWit or Taproot
	// These fields are explicitly set to empty to prevent confusion
	return &domainKey.WalletKey{
		WIF:            wif.String(),
		P2PKHAddr:      p2pkhAddr,
		P2SHSegWitAddr: "", // BCH does not support SegWit
		Bech32Addr:     "", // BCH does not support Bech32
		TaprootAddr:    "", // BCH does not support Taproot
		FullPubKey:     fullPubKey,
		RedeemScript:   "", // Not used for P2PKH
	}, nil
}

// getP2PKHAddr generates a P2PKH address in CashAddr format for BCH
// CashAddr format: <prefix>:<payload>
// Example: bitcoincash:qpm2qsznhks23z7629mms6s4cwef74vcwvy22gdx6a
func (s *BCHKeyStrategy) getP2PKHAddr(privKey *btcec.PrivateKey) (string, error) {
	serializedPubKey := privKey.PubKey().SerializeCompressed()
	pkHash := btcutil.Hash160(serializedPubKey)

	// First, create a standard BTC-style P2PKH address
	p2pkhAddr, err := btcutil.NewAddressPubKeyHash(pkHash, s.conf)
	if err != nil {
		return "", fmt.Errorf("failed to call btcutil.NewAddressPubKeyHash(): %w", err)
	}

	// Convert to BCH CashAddr format
	return s.convertToCashAddr(p2pkhAddr)
}

// convertToCashAddr converts a BTC-style P2PKH address to BCH CashAddr format
func (s *BCHKeyStrategy) convertToCashAddr(p2pkhAddr *btcutil.AddressPubKeyHash) (string, error) {
	addrBCH, err := bchutil.NewCashAddressPubKeyHash(p2pkhAddr.ScriptAddress(), s.conf)
	if err != nil {
		return "", fmt.Errorf("failed to call bchutil.NewCashAddressPubKeyHash(): %w", err)
	}

	// Get prefix (bitcoincash for mainnet, bchtest for testnet)
	prefix, ok := bchutil.Prefixes[s.conf.Name]
	if !ok {
		return "", fmt.Errorf("invalid BCH network config: %s", s.conf.Name)
	}

	return fmt.Sprintf("%s:%s", prefix, addrBCH.String()), nil
}
