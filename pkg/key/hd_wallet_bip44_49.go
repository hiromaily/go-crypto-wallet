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
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/pkg/errors"
)

//PurposeType BIP44は44固定
type PurposeType uint32

// purpose
const (
	PurposeTypeBIP44 PurposeType = 44 //BIP44
	PurposeTypeBIP49 PurposeType = 49 //BIP49
)

//TODO:同じアドレスを使い回すと、アドレスから総額情報がバレて危険
//よって、内部利用のアドレスは毎回使い捨てにすること

//ChangeType 受け取り階層
type ChangeType uint32

// change_type
const (
	ChangeTypeExternal ChangeType = 0 //外部送金者からの受け取り用 (ユーザー、集約用、マルチシグ)
	ChangeTypeInternal ChangeType = 1 //自身のトランザクションのおつり用 (出金時に使うトレード用アドレス) //TODO:これは使わないでいいかも
)

//CoinType コインの種類(
type CoinType uint32

// coin_type
const (
	CoinTypeBitcoin CoinType = 0 //Bitcoin
	CoinTypeTestnet CoinType = 1 //Testnet
)

//e.g. for Mainnet
//Client  => m/44/0/0/0/0~xxxxx
//Receipt => m/44/0/1
//Payment => m/44/0/2/0/0 => quoineから購入したものを受け取る用
//Payment => m/44/0/2/1/0 => 出金による支払いに利用、かつ、おつりも受け取る => TODO:ChangeTypeによってアドレスが変わってしまったら、どう運用するのか

// Key Keyオブジェクト
type Key struct {
	coinType enum.CoinType
	conf     *chaincfg.Params
}

// NewKey Keyオブジェクトを返す
func NewKey(coinType enum.CoinType, conf *chaincfg.Params) *Key {
	keyData := Key{
		coinType: coinType,
		conf:     conf,
	}

	return &keyData
}

// CreateAccount アカウント階層までのprivateKey及び publicKeyを生成する
func (k Key) CreateAccount(seed []byte, actType enum.AccountType) (string, string, error) {

	//Master
	masterKey, err := hdkeychain.NewMaster(seed, k.conf)
	if err != nil {
		return "", "", err
	}
	//Purpose
	purpose, err := masterKey.Child(hdkeychain.HardenedKeyStart + uint32(PurposeTypeBIP44))
	//purpose, err := masterKey.Child(hdkeychain.HardenedKeyStart + uint32(PurposeTypeBIP49))
	if err != nil {
		return "", "", err
	}
	//CoinType
	ct := uint32(CoinTypeBitcoin)
	if k.conf.Name != string(enum.NetworkTypeMainNet) {
		ct = uint32(CoinTypeTestnet)
	}
	coinType, err := purpose.Child(hdkeychain.HardenedKeyStart + ct)
	if err != nil {
		return "", "", err
	}
	//Account
	account, err := coinType.Child(hdkeychain.HardenedKeyStart + uint32(enum.AccountTypeValue[actType]))
	if err != nil {
		return "", "", err
	}
	//Change
	//Index

	publicKey, err := account.Neuter()
	if err != nil {
		return "", "", err
	}

	strPrivateKey := account.String()
	strPublicKey := publicKey.String()

	return strPrivateKey, strPublicKey, nil
}

// CreateKeysWithIndex 指定したindexに応じて複数のkeyを生成する
// e.g. [1] idxFrom:0,  count 10 => 0-9
//      [2] idxFrom:10, count 10 => 10-19
func (k Key) CreateKeysWithIndex(accountPrivateKey string, idxFrom, count uint32) ([]WalletKey, error) {
	account, err := hdkeychain.NewKeyFromString(accountPrivateKey)
	if err != nil {
		return nil, err
	}
	// Change
	change, err := account.Child(uint32(ChangeTypeExternal))
	if err != nil {
		return nil, err
	}

	// Index
	walletKeys := make([]WalletKey, count)
	//max := idxFrom + count
	//for i := uint32(idxFrom); i < max; i++ {
	for i := uint32(0); i < count; i++ {
		child, err := change.Child(idxFrom)
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

		// WIF　(compress: true) => bitcoin coreでは圧縮したアドレスを表示する
		wif, err := btcutil.NewWIF(privateKey, k.conf, true)
		if err != nil {
			return nil, err
		}
		strPrivateKey := wif.String()

		// Address(P2PKH) BTC/BCH
		//  btcutil.NewAddressPubKeyHash(pkHash, net)
		//if k.coinType == enum.BTC {
		//	address, err := child.Address(k.conf)
		//}
		strAddr, err := k.addressString(privateKey)
		if err != nil {
			return nil, err
		}

		// p2sh-segwit
		p2shSegwit, err := k.getP2shSegwit(privateKey)
		if err != nil {
			return nil, err
		}

		//address.String() とaddress.EncodeAddress()は結果として同じ
		walletKeys[i] = WalletKey{
			WIF:        strPrivateKey,
			Address:    strAddr, //[P2PKH]AddressPubKeyHash is an Address for a pay-to-pubkey-hash
			P2shSegwit: p2shSegwit,
			FullPubKey: getFullPubKey(privateKey, true),
		}

		idxFrom++
	}

	return walletKeys, nil
}

// GetExtendedKey for only debug use
func (k Key) GetExtendedKey(accountPrivateKey string) (*hdkeychain.ExtendedKey, error) {
	account, err := hdkeychain.NewKeyFromString(accountPrivateKey)
	if err != nil {
		return nil, err
	}
	// Change
	change, err := account.Child(uint32(ChangeTypeExternal))
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

	if k.coinType == enum.BTC {
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

// p2sh-segwitのstringを返す
// BCHは利用する予定はないが、念の為
//func (k Key) getP2shSegwit(privKey *btcec.PrivateKey) (*btcutil.AddressScriptHash, error) {
func (k Key) getP2shSegwit(privKey *btcec.PrivateKey) (string, error) {
	// []byte
	publicKeyHash := btcutil.Hash160(privKey.PubKey().SerializeCompressed())
	segwitAddress, err := btcutil.NewAddressWitnessPubKeyHash(publicKeyHash, k.conf)
	if err != nil {
		return "", errors.Errorf("btcutil.NewAddressWitnessPubKeyHash() error: %s", err)
	}
	//logger.Debugf("segwitAddress: %s", segwitAddress)

	redeemScript, err := txscript.PayToAddrScript(segwitAddress)
	if err != nil {
		return "", errors.Errorf("txscript.PayToAddrScript() error: %s", err)
	}
	//logger.Debugf("redeemScript: %s", redeemScript)

	//BTC
	if k.coinType == enum.BTC {
		address, err := btcutil.NewAddressScriptHash(redeemScript, k.conf)
		if err != nil {
			return "", errors.Errorf("btcutil.NewAddressScriptHash() error: %s", err)
		}
		//logger.Debugf("address.String() %s", address.String())
		return address.String(), nil
	}
	//BCH
	address, err := bchutil.NewCashAddressScriptHash(redeemScript, k.conf)
	if err != nil {
		return "", errors.Errorf("bchutil.NewCashAddressScriptHash() error: %s", err)
	}
	//logger.Debugf("address.String() %s", address.String())
	return address.String(), nil
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
