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

// Key Key object
type Key struct {
	purpose      PurposeType
	coinType     coin.CoinType
	coinTypeCode coin.CoinTypeCode
	conf         *chaincfg.Params
	logger       *zap.Logger
}

// NewKey returns Key
func NewKey(purpose PurposeType, coinTypeCode coin.CoinTypeCode, conf *chaincfg.Params, logger *zap.Logger) *Key {
	keyData := Key{
		purpose:      purpose,
		coinType:     coinTypeCode.CoinType(conf),
		coinTypeCode: coinTypeCode,
		conf:         conf,
		logger:       logger,
	}

	return &keyData
}

func (k Key) CreateKey(seed []byte, actType account.AccountType, idxFrom, count uint32) ([]WalletKey, error) {
	// create privateKey, publicKey by account level
	privKey, _, err := k.createKeyByAccount(seed, actType)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call createKeyByAccount()")
	}
	// create keys by index and count
	return k.createKeysWithIndex(privKey, idxFrom, count)
}

// createKeyByAccount create privateKey, publicKey by account level
func (k Key) createKeyByAccount(seed []byte, actType account.AccountType) (*hdkeychain.ExtendedKey, *hdkeychain.ExtendedKey, error) {

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
func (k Key) createKeysWithIndex(accountPrivKey *hdkeychain.ExtendedKey, idxFrom, count uint32) ([]WalletKey, error) {
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

		// full public Key
		//getFullPubKey(privateKey)

		// WIF　(compressed: true) => bitcoin core expresses compressed address
		wif, err := btcutil.NewWIF(privateKey, k.conf, true)
		if err != nil {
			return nil, err
		}

		// Address(P2PKH) BTC/BCH
		strAddr, err := k.addressString(privateKey)
		if err != nil {
			return nil, err
		}

		// p2sh-segwit
		p2shSegwit, redeemScript, err := k.getP2shSegwit(privateKey)
		if err != nil {
			return nil, err
		}
		if redeemScript != "" {
			k.logger.Debug("result of getP2shSegwit()", zap.String("redeemScript", redeemScript))
		}

		//address.String() とaddress.EncodeAddress()は結果として同じ
		walletKeys[i] = WalletKey{
			WIF:          wif.String(),
			Address:      strAddr, //[P2PKH]AddressPubKeyHash is an Address for a pay-to-pubkey-hash
			P2shSegwit:   p2shSegwit,
			FullPubKey:   getFullPubKey(privateKey, true),
			RedeemScript: redeemScript,
		}
		//idxFrom++
	}

	return walletKeys, nil
}

// BTC/BCHのaddress P2PKHを返す
//  address, err := child.Address(conf)
//  address.String() と同じ結果を返す
func (k Key) addressString(privKey *btcec.PrivateKey) (string, error) {
	serializedKey := privKey.PubKey().SerializeCompressed()
	pkHash := btcutil.Hash160(serializedKey)

	//*btcutil.AddressPubKeyHash
	addr, err := btcutil.NewAddressPubKeyHash(pkHash, k.conf)
	if err != nil {
		return "", errors.Errorf("btcutil.NewAddressPubKeyHash() error: %s", err)
	}

	if k.coinTypeCode == coin.BTC {
		//BTC
		return addr.String(), nil
	}

	//BCH *bchutil.CashAddressPubKeyHash
	//return bchutil.NewCashAddressPubKeyHash(pkHash, k.conf)
	addrBCH, err := bchutil.NewCashAddressPubKeyHash(addr.ScriptAddress(), k.conf)
	if err != nil {
		return "", errors.Errorf("btcutil.NewAddressPubKeyHash() error: %s", err)
	}

	//prefixを取得
	prefix, ok := bchutil.Prefixes[k.conf.Name]
	if !ok {
		return "", errors.New("[fatal error] chainConf *chaincfg.Params is wrong")
	}

	return fmt.Sprintf("%s:%s", prefix, addrBCH.String()), nil
}

// p2sh-segwitと、redeemScriptのstringとを返す
// BCHは利用する予定はないが、念の為
//func (k Key) getP2shSegwit(privKey *btcec.PrivateKey) (*btcutil.AddressScriptHash, error) {
func (k Key) getP2shSegwit(privKey *btcec.PrivateKey) (string, string, error) {
	// []byte
	publicKeyHash := btcutil.Hash160(privKey.PubKey().SerializeCompressed())
	segwitAddress, err := btcutil.NewAddressWitnessPubKeyHash(publicKeyHash, k.conf)
	if err != nil {
		return "", "", errors.Errorf("btcutil.NewAddressWitnessPubKeyHash() error: %s", err)
	}
	//logger.Debugf("segwitAddress: %s", segwitAddress)

	redeemScript, err := txscript.PayToAddrScript(segwitAddress)
	if err != nil {
		return "", "", errors.Errorf("txscript.PayToAddrScript() error: %s", err)
	}

	//この値はscriptPubKeyと一致するが、これはredeemScriptではない。
	//getaddressinfo APIでp2sh_segwit_addressをサーチすると、embedded側のscriptPubKeyと一致した
	//よって、この値は使えない。。。
	//Redeem Script => Hash of RedeemScript => p2SH ScriptPubKey
	//strRedeemScript := hex.EncodeToString(redeemScript)
	var strRedeemScript string //暫定

	//BTC
	if k.coinTypeCode == coin.BTC {
		address, err := btcutil.NewAddressScriptHash(redeemScript, k.conf)
		if err != nil {
			return "", "", errors.Errorf("btcutil.NewAddressScriptHash() error: %s", err)
		}
		//logger.Debugf("address.String() %s", address.String())

		return address.String(), strRedeemScript, nil
	}
	//BCH
	address, err := bchutil.NewCashAddressScriptHash(redeemScript, k.conf)
	if err != nil {
		return "", "", errors.Errorf("bchutil.NewCashAddressScriptHash() error: %s", err)
	}
	//logger.Debugf("address.String() %s", address.String())
	return address.String(), strRedeemScript, nil
}

// getPubKey fullのPublic Keyを返す
func getFullPubKey(privKey *btcec.PrivateKey, isCompressed bool) string {
	var bPubKey []byte
	if isCompressed {
		//Compressed
		bPubKey = privKey.PubKey().SerializeCompressed()
	} else {
		//Uncompressed
		bPubKey = privKey.PubKey().SerializeUncompressed()
	}
	//logger.Debugf("bPubKey hash: %s", btcutil.Hash160(bPubKey))

	hexPubKey := hex.EncodeToString(bPubKey)
	//logger.Debugf("hex.EncodeToString(bPubKey): %s", hexPubKey)

	//key *PublicKey
	//bHexPubKey, _ := hex.DecodeString(hexPubKey)
	//pubKey, _ := btcec.ParsePubKey(bHexPubKey, btcec.S256())

	return hexPubKey
}

// GetExtendedKey for only debug use
func GetExtendedKey(accountPrivateKey string) (*hdkeychain.ExtendedKey, error) {
	account, err := hdkeychain.NewKeyFromString(accountPrivateKey)
	if err != nil {
		return nil, err
	}
	// Change
	change, err := account.Child(ChangeTypeExternal.Uint32())
	if err != nil {
		return nil, err
	}

	child, err := change.Child(0)
	if err != nil {
		return nil, err
	}

	// extendedKey
	return child, nil
}
