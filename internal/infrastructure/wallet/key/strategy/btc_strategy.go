package strategy

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// BTCKeyStrategy implements CoinKeyStrategy for Bitcoin (BTC).
// Bitcoin supports multiple address types:
// - P2PKH (Pay To PubKey Hash): Legacy addresses starting with '1'
// - P2SH-SegWit (Pay To Script Hash with SegWit): Nested SegWit addresses starting with '3'
// - Bech32 (P2WPKH): Native SegWit addresses starting with 'bc1' (or 'tb1' for testnet)
// - Taproot (P2TR): Taproot addresses starting with 'bc1p' (BIP86)
//
// All four address types are generated for maximum compatibility.
type BTCKeyStrategy struct {
	conf *chaincfg.Params
}

// Compile-time check that BTCKeyStrategy implements CoinKeyStrategy
var _ portsWallet.CoinKeyStrategy = (*BTCKeyStrategy)(nil)

// NewBTCKeyStrategy creates a new BTCKeyStrategy
func NewBTCKeyStrategy(conf *chaincfg.Params) *BTCKeyStrategy {
	return &BTCKeyStrategy{
		conf: conf,
	}
}

// CoinTypeCode returns the coin type code for Bitcoin
func (*BTCKeyStrategy) CoinTypeCode() domainCoin.CoinTypeCode {
	return domainCoin.BTC
}

// GenerateWalletKey generates all BTC address types from a private key
func (s *BTCKeyStrategy) GenerateWalletKey(privKey *btcec.PrivateKey) (*domainKey.WalletKey, error) {
	// WIF (Wallet Import Format) - compressed
	wif, err := btcutil.NewWIF(privKey, s.conf, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create WIF: %w", err)
	}

	// P2PKH address (Legacy)
	p2pkhAddr, err := s.getP2PKHAddr(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate P2PKH address: %w", err)
	}

	// P2SH-SegWit address (Nested SegWit)
	p2shSegWitAddr, redeemScript, err := s.getP2SHSegWitAddr(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate P2SH-SegWit address: %w", err)
	}

	// Bech32 address (Native SegWit)
	bech32Addr, err := s.getBech32Addr(wif)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Bech32 address: %w", err)
	}

	// Taproot address (BIP86)
	taprootAddr, err := s.getTaprootAddr(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Taproot address: %w", err)
	}

	// Full public key (compressed)
	fullPubKey := getFullPubKey(privKey, true)

	return &domainKey.WalletKey{
		WIF:            wif.String(),
		P2PKHAddr:      p2pkhAddr,
		P2SHSegWitAddr: p2shSegWitAddr,
		Bech32Addr:     bech32Addr.EncodeAddress(),
		TaprootAddr:    taprootAddr.EncodeAddress(),
		FullPubKey:     fullPubKey,
		RedeemScript:   redeemScript,
	}, nil
}

// getP2PKHAddr generates a P2PKH (Pay To PubKey Hash) address
// P2PKH is the original Bitcoin address format (addresses starting with '1')
// https://bitcoin.org/en/glossary/p2pkh-address
func (s *BTCKeyStrategy) getP2PKHAddr(privKey *btcec.PrivateKey) (string, error) {
	serializedPubKey := privKey.PubKey().SerializeCompressed()
	pkHash := btcutil.Hash160(serializedPubKey)

	p2pkhAddr, err := btcutil.NewAddressPubKeyHash(pkHash, s.conf)
	if err != nil {
		return "", fmt.Errorf("failed to call btcutil.NewAddressPubKeyHash(): %w", err)
	}

	return p2pkhAddr.String(), nil
}

// getP2SHSegWitAddr generates a P2SH-SegWit address (P2SH nested SegWit) and redeem script
// P2SH-SegWit addresses provide SegWit benefits while maintaining compatibility with older wallets
// The redeem script reveals the witness program when spending
func (s *BTCKeyStrategy) getP2SHSegWitAddr(privKey *btcec.PrivateKey) (string, string, error) {
	pubKeyHash := btcutil.Hash160(privKey.PubKey().SerializeCompressed())
	segwitAddress, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, s.conf)
	if err != nil {
		return "", "", fmt.Errorf("failed to call btcutil.NewAddressWitnessPubKeyHash(): %w", err)
	}

	// Get redeem script (witness program)
	payToAddrScript, err := txscript.PayToAddrScript(segwitAddress)
	if err != nil {
		return "", "", fmt.Errorf("failed to call txscript.PayToAddrScript(): %w", err)
	}

	// For P2SH-SegWit, the redeemScript IS the witness program
	// The witness program (OP_0 <hash>) is hashed to create the P2SH address
	strRedeemScript := hex.EncodeToString(payToAddrScript)

	btcAddress, err := btcutil.NewAddressScriptHash(payToAddrScript, s.conf)
	if err != nil {
		return "", "", fmt.Errorf("failed to call btcutil.NewAddressScriptHash(): %w", err)
	}

	return btcAddress.String(), strRedeemScript, nil
}

// getBech32Addr generates a Bech32 address (Native SegWit)
// Bech32 addresses are more efficient than P2SH-SegWit and have better error detection
func (s *BTCKeyStrategy) getBech32Addr(wif *btcutil.WIF) (*btcutil.AddressWitnessPubKeyHash, error) {
	witnessProg := btcutil.Hash160(wif.SerializePubKey())
	bech32Addr, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, s.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to call NewAddressWitnessPubKeyHash(): %w", err)
	}
	return bech32Addr, nil
}

// getTaprootAddr generates a Taproot address (BIP86) for the given private key
// BIP86 uses key path spending without script path (no merkle root)
// Taproot provides improved privacy and more efficient complex transactions
func (s *BTCKeyStrategy) getTaprootAddr(privKey *btcec.PrivateKey) (*btcutil.AddressTaproot, error) {
	// Get the internal public key
	internalPubKey := privKey.PubKey()

	// Compute the tweaked Taproot output key (BIP341) without script path
	taprootKey := txscript.ComputeTaprootKeyNoScript(internalPubKey)

	// Get the 32-byte x-only public key (Schnorr public key)
	witnessProg := taprootKey.SerializeCompressed()[1:] // Remove the parity byte

	// Create Taproot address
	taprootAddr, err := btcutil.NewAddressTaproot(witnessProg, s.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to call NewAddressTaproot(): %w", err)
	}

	return taprootAddr, nil
}
