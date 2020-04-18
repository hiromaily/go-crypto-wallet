package key

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cpacia/bchutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

//TODO: except client address, same address is used only once due to security
// - address could be traced easily
// - so any internal addresses should be disposable

//BIP44
//https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#Purpose
// m / purpose' / coin_type' / account' / change / address_index

//e.g.  BIP44, Bitcoin Mainnet
//Client account  => m/44/0/0/0/0~xxxxx
//Receipt account => m/44/0/1/0/0~xxxxx
//Payment account => m/44/0/2/0/0~xxxxx

// PurposeType BIP44/BIP49, for now 44 is used as fixed value
type PurposeType uint32

func (t PurposeType) Uint32() uint32 {
	return uint32(t)
}

// purpose depends on BIP, BIP44  is a constant set to `44`
const (
	PurposeTypeBIP44 PurposeType = 44 //BIP44
	PurposeTypeBIP49 PurposeType = 49 //BIP49
)

// CoinType creates a separate subtree for every cryptocoin
//  which come from `CoinType` in go-bitcoin/pkg/wallet/coin/types.go

type CoinType uint32

func (t CoinType) Uint32() uint32 {
	return uint32(t)
}

// coin_type
// https://github.com/satoshilabs/slips/blob/master/slip-0044.md

// Account
// account come from `AccountType` in go-bitcoin/pkg/account.go

//ChangeType  external or internal use
type ChangeType uint32

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
	coinType     coin.CoinType
	coinTypeCode coin.CoinTypeCode
	conf         *chaincfg.Params
	logger       *zap.Logger
}

// NewHDKey returns Key
func NewHDKey(purpose PurposeType, coinTypeCode coin.CoinTypeCode, conf *chaincfg.Params, logger *zap.Logger) *HDKey {
	keyData := HDKey{
		purpose:      purpose,
		coinType:     coinTypeCode.CoinType(conf),
		coinTypeCode: coinTypeCode,
		conf:         conf,
		logger:       logger,
	}

	return &keyData
}

func (k *HDKey) CreateKey(seed []byte, actType account.AccountType, idxFrom, count uint32) ([]WalletKey, error) {
	// create privateKey, publicKey by account level
	privKey, _, err := k.createKeyByAccount(seed, actType)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call createKeyByAccount()")
	}
	// create keys by index and count
	return k.createKeysWithIndex(privKey, idxFrom, count)
}

// createKeyByAccount create privateKey, publicKey by account level
func (k *HDKey) createKeyByAccount(seed []byte, actType account.AccountType) (*hdkeychain.ExtendedKey, *hdkeychain.ExtendedKey, error) {

	//Master
	masterKey, err := hdkeychain.NewMaster(seed, k.conf)
	if err != nil {
		return nil, nil, err
	}
	//Purpose
	purpose, err := masterKey.Child(hdkeychain.HardenedKeyStart + k.purpose.Uint32())
	if err != nil {
		return nil, nil, err
	}
	//CoinType
	coinType, err := purpose.Child(hdkeychain.HardenedKeyStart + k.coinType.Uint32())
	if err != nil {
		return nil, nil, err
	}
	//Account
	accountPrivKey, err := coinType.Child(hdkeychain.HardenedKeyStart + account.AccountTypeValue[actType])
	if err != nil {
		return nil, nil, err
	}
	//Change
	//Index

	// get pubKey
	publicKey, err := accountPrivKey.Neuter()
	if err != nil {
		return nil, nil, err
	}

	//strPrivateKey := account.String()
	//strPublicKey := publicKey.String()
	return accountPrivKey, publicKey, nil
}

// createKeysWithIndex create keys by index and count
// e.g. - idxFrom:0,  count 10 => 0-9
//      - idxFrom:10, count 10 => 10-19
func (k *HDKey) createKeysWithIndex(accountPrivKey *hdkeychain.ExtendedKey, idxFrom, count uint32) ([]WalletKey, error) {
	//accountPrivKey, err := hdkeychain.NewKeyFromString(accountPrivKey)

	// Change
	change, err := accountPrivKey.Child(ChangeTypeExternal.Uint32())
	if err != nil {
		return nil, err
	}

	// Index
	walletKeys := make([]WalletKey, count)
	for i := uint32(0); i < count; i++ {
		child, err := change.Child(idxFrom + i)
		if err != nil {
			return nil, err
		}

		// privateKey
		privateKey, err := child.ECPrivKey()
		if err != nil {
			return nil, err
		}

		// WIF　(compressed: true) => bitcoin core expresses compressed address
		wif, err := btcutil.NewWIF(privateKey, k.conf, true)
		if err != nil {
			return nil, err
		}

		// get P2PKH address as string for BTC/BCH
		//if only BTC, this logic would be enough
		//address, err := child.Address(conf)
		//address.String()
		strP2PKHAddr, err := k.getP2pkhAddr(privateKey)
		if err != nil {
			return nil, err
		}

		// p2sh-segwit address
		strP2shSegwit, redeemScript, err := k.getP2shSegwitAddr(privateKey)
		if err != nil {
			return nil, err
		}

		// address.String() is equal to address.EncodeAddress()
		walletKeys[i] = WalletKey{
			WIF:          wif.String(),
			Address:      strP2PKHAddr, //[P2PKH]AddressPubKeyHash is an Address for a pay-to-pubkey-hash
			P2shSegwit:   strP2shSegwit,
			FullPubKey:   getFullPubKey(privateKey, true),
			RedeemScript: redeemScript,
		}
	}

	return walletKeys, nil
}

// get Address(P2PKH) as string for BTC/BCH
func (k *HDKey) getP2pkhAddr(privKey *btcec.PrivateKey) (string, error) {
	serializedPubKey := privKey.PubKey().SerializeCompressed()
	pkHash := btcutil.Hash160(serializedPubKey)

	//*btcutil.AddressPubKeyHash
	p2PKHAddr, err := btcutil.NewAddressPubKeyHash(pkHash, k.conf)
	if err != nil {
		return "", errors.Wrapf(err, "fail to call btcutil.NewAddressPubKeyHash()")
	}

	switch k.coinTypeCode {
	case coin.BTC:
		return p2PKHAddr.String(), nil
	case coin.BCH:
		k.getP2PKHAddrBCH(p2PKHAddr)
	}
	return "", errors.Errorf("getP2pkhAddr() is not implemented for %s", k.coinTypeCode)
}

// BCH *bchutil.CashAddressPubKeyHash
func (k *HDKey) getP2PKHAddrBCH(p2PKHAddr *btcutil.AddressPubKeyHash) (string, error) {
	addrBCH, err := bchutil.NewCashAddressPubKeyHash(p2PKHAddr.ScriptAddress(), k.conf)
	if err != nil {
		return "", errors.Wrap(err, "fail to call btcutil.NewAddressPubKeyHash()")
	}

	// get prefix
	prefix, ok := bchutil.Prefixes[k.conf.Name]
	if !ok {
		return "", errors.Errorf("invalid BCH *chaincfg : %s", k.conf.Name)
	}
	return fmt.Sprintf("%s:%s", prefix, addrBCH.String()), nil
}

// FIXME: getting RedeemScript is not fixed yet
// get p2sh-segwit address and redeemScript as string
//  - it's for only BTC
//  - Though BCH would not require it, just in case
func (k *HDKey) getP2shSegwitAddr(privKey *btcec.PrivateKey) (string, string, error) {
	// []byte
	pubKeyHash := btcutil.Hash160(privKey.PubKey().SerializeCompressed())
	segwitAddress, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, k.conf)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btcutil.NewAddressWitnessPubKeyHash()")
	}

	//FIXME: getting RedeemScript is not fixed yet
	// get redeemScript
	payToAddrScript, err := txscript.PayToAddrScript(segwitAddress)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call txscript.PayToAddrScript()")
	}

	// value of payToAddrScript is equal to scriptPubKey, but it's not redeemScript
	// if call `getaddressinfo` API, result includes this value as scriptPubKey in embedded in p2sh_segwit_address
	// That's why payToAddrScript is not used as redeemScript
	// Redeem Script => Hash of RedeemScript => p2SH ScriptPubKey

	var strRedeemScript string //FIXME: not implemented yet
	switch k.coinTypeCode {
	case coin.BTC:
		address, err := btcutil.NewAddressScriptHash(payToAddrScript, k.conf)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call btcutil.NewAddressScriptHash()")
		}
		return address.String(), strRedeemScript, nil
	case coin.BCH:
		address, err := bchutil.NewCashAddressScriptHash(payToAddrScript, k.conf)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call bchutil.NewCashAddressScriptHash()")
		}
		return address.String(), strRedeemScript, nil
	}
	return "", "", errors.Errorf("getP2shSegwitAddr() is not implemented yet for %s", k.coinTypeCode)
}

// getPubKey returns full Public Key
func getFullPubKey(privKey *btcec.PrivateKey, isCompressed bool) string {
	var bPubKey []byte
	if isCompressed {
		//Compressed
		bPubKey = privKey.PubKey().SerializeCompressed()
	} else {
		//Uncompressed
		bPubKey = privKey.PubKey().SerializeUncompressed()
	}
	hexPubKey := hex.EncodeToString(bPubKey)
	return hexPubKey
}
