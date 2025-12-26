package key

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160" //nolint:gosec

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
	bchaddr "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address/bch"
	xrpaddr "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// TODO: except client address, same address is used only once due to security
// - address could be traced easily
// - so any internal addresses should be disposable

// BIP44
// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#Purpose
// m / purpose' / coin_type' / account' / change / address_index

// e.g.  BIP44, Bitcoin Mainnet
// Client account  => m/44/0/0/0/0~xxxxx
// Receipt account => m/44/0/1/0/0~xxxxx
// Payment account => m/44/0/2/0/0~xxxxx

// PurposeType BIP44/BIP49, for now 44 is used as fixed value
type PurposeType uint32

// Uint32 converter
func (t PurposeType) Uint32() uint32 {
	return uint32(t)
}

// purpose depends on BIP, BIP44  is a constant set to `44`
const (
	PurposeTypeBIP44 PurposeType = 44 // BIP44
	PurposeTypeBIP49 PurposeType = 49 // BIP49
	PurposeTypeBIP84 PurposeType = 84 // BIP84
	PurposeTypeBIP86 PurposeType = 86 // BIP86
)

// CoinType creates a separate subtree for every cryptocoin
//
//	which come from `CoinType` in go-crypto-wallet/pkg/wallet/coin/types.go
type CoinType uint32

// Uint32 converter
func (t CoinType) Uint32() uint32 {
	return uint32(t)
}

// coin_type
// https://github.com/satoshilabs/slips/blob/master/slip-0044.md

// Account
// account come from `AccountType` in go-crypto-wallet/pkg/account/public_account.go

// ChangeType  external or internal use
type ChangeType uint32

// Uint32 converter
func (t ChangeType) Uint32() uint32 {
	return uint32(t)
}

// change_type
const (
	ChangeTypeExternal ChangeType = 0 // constant 0 is used for external chain
	ChangeTypeInternal ChangeType = 1 // constant 1 for internal chain (also known as change addresses)
)

// HDKey HD Wallet Key object
type HDKey struct {
	purpose      PurposeType
	coinType     domainCoin.CoinType
	coinTypeCode domainCoin.CoinTypeCode
	conf         *chaincfg.Params
}

// NewHDKey returns Key
func NewHDKey(
	purpose PurposeType, coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params,
) *HDKey {
	keyData := HDKey{
		purpose:      purpose,
		coinType:     domainCoin.GetCoinType(coinTypeCode, conf),
		coinTypeCode: coinTypeCode,
		conf:         conf,
	}

	return &keyData
}

// CreateKey create hd key
func (k *HDKey) CreateKey(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
	// create privateKey, publicKey by account level
	privKey, _, err := k.createKeyByAccount(seed, accountType)
	if err != nil {
		return nil, fmt.Errorf("fail to call createKeyByAccount(): %w", err)
	}
	// create keys by index and count
	return k.createKeysWithIndex(privKey, idxFrom, count)
}

// KeyType returns the key type this generator supports (implements Generator interface)
func (k *HDKey) KeyType() domainKey.KeyType {
	switch k.purpose {
	case PurposeTypeBIP44:
		return domainKey.KeyTypeBIP44
	case PurposeTypeBIP49:
		return domainKey.KeyTypeBIP49
	case PurposeTypeBIP84:
		return domainKey.KeyTypeBIP84
	case PurposeTypeBIP86:
		return domainKey.KeyTypeBIP86
	default:
		return domainKey.KeyTypeBIP44
	}
}

// SupportsAddressType checks if this generator supports the given address type (implements Generator interface)
func (k *HDKey) SupportsAddressType(addrType address.AddrType) bool {
	switch k.purpose {
	case PurposeTypeBIP44:
		return addrType == address.AddrTypeLegacy
	case PurposeTypeBIP49:
		return addrType == address.AddrTypeP2shSegwit
	case PurposeTypeBIP84:
		return addrType == address.AddrTypeBech32
	case PurposeTypeBIP86:
		return addrType == address.AddrTypeTaproot
	default:
		return false
	}
}

// GetDerivationPath returns the derivation path for the given account and index (implements Generator interface)
func (k *HDKey) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
	return fmt.Sprintf("m/%d'/%d'/%d'/0/%d",
		k.purpose.Uint32(),
		k.coinType.Uint32(),
		accountType.Uint32(),
		index)
}

// createKeyByAccount create privateKey, publicKey by account level
func (k *HDKey) createKeyByAccount(
	seed []byte, accountType domainAccount.AccountType,
) (*hdkeychain.ExtendedKey, *hdkeychain.ExtendedKey, error) {
	// Master
	masterKey, err := hdkeychain.NewMaster(seed, k.conf)
	if err != nil {
		return nil, nil, err
	}
	// Purpose
	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + k.purpose.Uint32())
	if err != nil {
		return nil, nil, err
	}
	// CoinType
	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + k.coinType.Uint32())
	if err != nil {
		return nil, nil, err
	}
	// Account
	logger.Debug(
		"create_key_by_account",
		"account_type", accountType.String(),
		"account_value", accountType.Uint32(),
	)
	accountPrivKey, err := coinType.Derive(hdkeychain.HardenedKeyStart + accountType.Uint32())
	if err != nil {
		return nil, nil, err
	}
	// Change
	// Index

	// get pubKey
	publicKey, err := accountPrivKey.Neuter()
	if err != nil {
		return nil, nil, err
	}

	// strPrivateKey := account.String()
	// strPublicKey := publicKey.String()
	return accountPrivKey, publicKey, nil
}

// createKeysWithIndex create keys by index and count
// e.g. - idxFrom:0,  count 10 => 0-9
//   - idxFrom:10, count 10 => 10-19
func (k *HDKey) createKeysWithIndex(
	accountPrivKey *hdkeychain.ExtendedKey, idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
	// accountPrivKey, err := hdkeychain.NewKeyFromString(accountPrivKey)

	// Change
	change, err := accountPrivKey.Derive(ChangeTypeExternal.Uint32())
	if err != nil {
		return nil, err
	}

	// Index
	walletKeys := make([]domainKey.WalletKey, count)
	for i := uint32(0); i < count; i++ {
		var loopErr error
		var child *hdkeychain.ExtendedKey
		child, loopErr = change.Derive(idxFrom + i)
		if loopErr != nil {
			return nil, loopErr
		}

		// privateKey
		var privateKey *btcec.PrivateKey
		privateKey, loopErr = child.ECPrivKey()
		if loopErr != nil {
			return nil, loopErr
		}

		switch k.coinTypeCode {
		case domainCoin.BTC, domainCoin.BCH:
			// WIF　(compressed: true) => bitcoin core expresses compressed address
			var wif *btcutil.WIF
			wif, loopErr = btcutil.NewWIF(privateKey, k.conf, true)
			if loopErr != nil {
				return nil, loopErr
			}

			var strP2PKHAddr, strP2SHSegWitAddr string
			var bech32Addr *btcutil.AddressWitnessPubKeyHash
			var redeemScript string
			strP2PKHAddr, strP2SHSegWitAddr, bech32Addr, redeemScript, loopErr = k.btcAddrs(wif, privateKey)
			if loopErr != nil {
				return nil, loopErr
			}

			// Generate Taproot address
			var taprootAddr *btcutil.AddressTaproot
			taprootAddr, loopErr = k.getTaprootAddr(privateKey)
			if loopErr != nil {
				return nil, loopErr
			}

			// address.String() is equal to address.EncodeAddress()
			walletKeys[i] = domainKey.WalletKey{
				WIF:            wif.String(),
				P2PKHAddr:      strP2PKHAddr,
				P2SHSegWitAddr: strP2SHSegWitAddr,
				Bech32Addr:     bech32Addr.EncodeAddress(),
				TaprootAddr:    taprootAddr.EncodeAddress(),
				FullPubKey:     getFullPubKey(privateKey, true),
				RedeemScript:   redeemScript,
			}

		case domainCoin.ETH:
			var ethAddr, ethPubKey, ethPrivKey string
			ethAddr, ethPubKey, ethPrivKey, loopErr = k.ethAddrs(privateKey)
			if loopErr != nil {
				return nil, loopErr
			}

			walletKeys[i] = domainKey.WalletKey{
				WIF:            ethPrivKey,
				P2PKHAddr:      ethAddr,
				P2SHSegWitAddr: "",
				Bech32Addr:     "",
				TaprootAddr:    "",
				FullPubKey:     ethPubKey,
				RedeemScript:   "",
			}
		case domainCoin.XRP:
			var xrpAddr, xrpPubKey, xrpPrivKey string
			xrpAddr, xrpPubKey, xrpPrivKey, loopErr = k.xrpAddrs(privateKey)
			if loopErr != nil {
				return nil, loopErr
			}

			// eth address is used as passphrase for generating key by API `wallet_propose`
			var ethAddr string
			ethAddr, _, _, loopErr = k.ethAddrs(privateKey)
			if loopErr != nil {
				return nil, loopErr
			}

			walletKeys[i] = domainKey.WalletKey{
				WIF:            xrpPrivKey,
				P2PKHAddr:      xrpAddr,
				P2SHSegWitAddr: ethAddr,
				Bech32Addr:     "",
				TaprootAddr:    "",
				FullPubKey:     xrpPubKey,
				RedeemScript:   "",
			}
		case domainCoin.LTC, domainCoin.ERC20, domainCoin.HYT:
			return nil, fmt.Errorf("coinType[%s] is not implemented yet", k.coinTypeCode.String())
		default:
			return nil, fmt.Errorf("coinType[%s] is not implemented yet", k.coinTypeCode.String())
		}
	}

	return walletKeys, nil
}

func (k *HDKey) btcAddrs(
	wif *btcutil.WIF, privKey *btcec.PrivateKey,
) (string, string, *btcutil.AddressWitnessPubKeyHash, string, error) {
	// P2SH address

	// get P2PKH address as string for BTC/BCH
	// - P2PKH Address, Pay To PubKey Hash
	// - if only BTC, this logic would be enough
	//  address, err := child.Address(conf)
	//  address.String()
	strP2PKHAddr, err := k.getP2PKHAddr(privKey)
	if err != nil {
		return "", "", nil, "", err
	}

	// P2SH-SegWit address
	strP2SHSegWitAddr, redeemScript, err := k.getP2SHSegWitAddr(privKey)
	if err != nil {
		return "", "", nil, "", err
	}

	// Bech32 address
	bech32Addr, err := k.getBech32Addr(wif)
	if err != nil {
		return "", "", nil, "", err
	}
	return strP2PKHAddr, strP2SHSegWitAddr, bech32Addr, redeemScript, nil
}

// https://goethereumbook.org/wallet-generate/
func (*HDKey) ethAddrs(privKey *btcec.PrivateKey) (string, string, string, error) {
	// private key
	ethPrivKey := privKey.ToECDSA()
	ethHexPrivKey := hexutil.Encode(crypto.FromECDSA(ethPrivKey))

	// pubkey, address
	ethPubkey := ethPrivKey.Public()
	pubkeyECDSA, ok := ethPubkey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", "", errors.New("fail to call cast pubkey to ecsda pubkey")
	}
	// pubkey
	ethHexPubKey := hexutil.Encode(crypto.FromECDSAPub(pubkeyECDSA))[4:]

	// address
	address := crypto.PubkeyToAddress(*pubkeyECDSA).Hex()

	return address, ethHexPubKey, ethHexPrivKey, nil
}

func (*HDKey) xrpAddrs(privKey *btcec.PrivateKey) (string, string, string, error) {
	// private key (same as ethereum for now)
	xrpPrivKey := privKey.ToECDSA()
	// xrpHexPrivKey := hexutil.Encode(crypto.FromECDSA(xrpPrivKey))
	xrpHexPrivKey, err := xrpaddr.NewAccountPrivateKey(crypto.FromECDSA(xrpPrivKey))
	if err != nil {
		return "", "", "", fmt.Errorf("fail to call xrpaddr.NewAccountPrivateKey(): %w", err)
	}

	serializedPubKey := privKey.PubKey().SerializeCompressed()
	pubKeyHash := xrpaddr.Sha256RipeMD160(serializedPubKey)
	if len(pubKeyHash) != ripemd160.Size {
		return "", "", "", errors.New("pubKeyHash must be 20 bytes")
	}
	// address
	address, err := xrpaddr.NewAccountID(pubKeyHash)
	if err != nil {
		return "", "", "", fmt.Errorf("fail to call rcrypto.NewAccountID(): %w", err)
	}
	// publicKey
	publicKey, err := xrpaddr.NewAccountPublicKey(pubKeyHash)
	if err != nil {
		return "", "", "", fmt.Errorf("fail to call rcrypto.NewAccountPublicKey(): %w", err)
	}

	return address.String(), publicKey.String(), xrpHexPrivKey.String(), nil
}

// get Address(P2PKH) as string for BTC/BCH
// P2PKH Address, Pay To PubKey Hash
// https://bitcoin.org/en/glossary/p2pkh-address
func (k *HDKey) getP2PKHAddr(privKey *btcec.PrivateKey) (string, error) {
	serializedPubKey := privKey.PubKey().SerializeCompressed()
	pkHash := btcutil.Hash160(serializedPubKey)

	// *btcutil.AddressPubKeyHash
	p2PKHAddr, err := btcutil.NewAddressPubKeyHash(pkHash, k.conf)
	if err != nil {
		return "", fmt.Errorf("fail to call btcutil.NewAddressPubKeyHash(): %w", err)
	}

	switch k.coinTypeCode {
	case domainCoin.BTC:
		return p2PKHAddr.String(), nil
	case domainCoin.BCH:
		return k.getP2PKHAddrBCH(p2PKHAddr)
	case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.ERC20, domainCoin.HYT:
		return "", fmt.Errorf("getP2pkhAddr() is not implemented for %s", k.coinTypeCode)
	default:
		return "", fmt.Errorf("getP2pkhAddr() is not implemented for %s", k.coinTypeCode)
	}
}

// getP2PKHAddrBCH get P2PKH Addr for BCH
func (k *HDKey) getP2PKHAddrBCH(p2PKHAddr *btcutil.AddressPubKeyHash) (string, error) {
	addrBCH, err := bchaddr.NewCashAddressPubKeyHash(p2PKHAddr.ScriptAddress(), k.conf)
	if err != nil {
		return "", fmt.Errorf("fail to call btcutil.NewAddressPubKeyHash(): %w", err)
	}

	// get prefix
	prefix, ok := bchaddr.Prefixes[k.conf.Name]
	if !ok {
		return "", fmt.Errorf("invalid BCH *chaincfg : %s", k.conf.Name)
	}
	return fmt.Sprintf("%s:%s", prefix, addrBCH.String()), nil
}

// getP2SHSegWitAddr get P2SH-SegWit address (P2SH nested SegWit) and redeemScript as string
//   - it's for only BTC
//   - Though BCH would not require it, just in case
//
// FIXME: getting RedeemScript is not fixed yet
//
//nolint:unparam // redeemScript (second return value) is not implemented yet, will be fixed in future
func (k *HDKey) getP2SHSegWitAddr(privKey *btcec.PrivateKey) (string, string, error) {
	// []byte
	pubKeyHash := btcutil.Hash160(privKey.PubKey().SerializeCompressed())
	segwitAddress, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, k.conf)
	if err != nil {
		return "", "", fmt.Errorf("fail to call btcutil.NewAddressWitnessPubKeyHash(): %w", err)
	}

	// FIXME: getting RedeemScript is not fixed yet
	// get redeemScript
	payToAddrScript, err := txscript.PayToAddrScript(segwitAddress)
	if err != nil {
		return "", "", fmt.Errorf("fail to call txscript.PayToAddrScript(): %w", err)
	}

	// value of payToAddrScript is equal to scriptPubKey, but it's not redeemScript
	// if call `getaddressinfo` API, result includes this value as scriptPubKey in embedded in p2sh_segwit_address
	// That's why payToAddrScript is not used as redeemScript
	// Redeem Script => Hash of RedeemScript => p2SH ScriptPubKey

	var strRedeemScript string // FIXME: not implemented yet
	switch k.coinTypeCode {
	case domainCoin.BTC:
		btcAddress, addrErr := btcutil.NewAddressScriptHash(payToAddrScript, k.conf)
		if addrErr != nil {
			return "", "", fmt.Errorf("fail to call btcutil.NewAddressScriptHash(): %w", addrErr)
		}
		return btcAddress.String(), strRedeemScript, nil
	case domainCoin.BCH:
		bchAddress, addrErr := bchaddr.NewCashAddressScriptHash(payToAddrScript, k.conf)
		if addrErr != nil {
			return "", "", fmt.Errorf("fail to call bchaddr.NewCashAddressScriptHash(): %w", addrErr)
		}
		return bchAddress.String(), strRedeemScript, nil
	case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.ERC20, domainCoin.HYT:
		return "", "", fmt.Errorf("getP2shSegwitAddr() is not implemented yet for %s", k.coinTypeCode)
	default:
		return "", "", fmt.Errorf("getP2shSegwitAddr() is not implemented yet for %s", k.coinTypeCode)
	}
}

// getBech32Addr returns bech32 address
func (k *HDKey) getBech32Addr(wif *btcutil.WIF) (*btcutil.AddressWitnessPubKeyHash, error) {
	witnessProg := btcutil.Hash160(wif.SerializePubKey())
	bech32Addr, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, k.conf)
	if err != nil {
		return nil, fmt.Errorf("fail to call NewAddressWitnessPubKeyHash(): %w", err)
	}
	return bech32Addr, nil
}

// getTaprootAddr returns a Taproot address (BIP86) for the given private key
// BIP86 uses key path spending without script path (no merkle root)
func (k *HDKey) getTaprootAddr(privKey *btcec.PrivateKey) (*btcutil.AddressTaproot, error) {
	// Get the internal public key
	internalPubKey := privKey.PubKey()

	// Compute the tweaked Taproot output key (BIP341) without script path
	taprootKey := txscript.ComputeTaprootKeyNoScript(internalPubKey)

	// Get the 32-byte x-only public key (Schnorr public key)
	witnessProg := taprootKey.SerializeCompressed()[1:] // Remove the parity byte

	// Create Taproot address
	taprootAddr, err := btcutil.NewAddressTaproot(witnessProg, k.conf)
	if err != nil {
		return nil, fmt.Errorf("fail to call NewAddressTaproot(): %w", err)
	}

	return taprootAddr, nil
}

// getFullPubKey returns full Public Key
func getFullPubKey(privKey *btcec.PrivateKey, isCompressed bool) string {
	var bPubKey []byte
	if isCompressed {
		// Compressed
		bPubKey = privKey.PubKey().SerializeCompressed()
	} else {
		// Uncompressed
		bPubKey = privKey.PubKey().SerializeUncompressed()
	}
	hexPubKey := hex.EncodeToString(bPubKey)
	return hexPubKey
}
