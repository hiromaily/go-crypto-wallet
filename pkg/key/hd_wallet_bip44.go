package key

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
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

//e.g. for Mainnet
//Client  => m/44/0/0/0/0~xxxxx
//Receipt => m/44/0/1
//Payment => m/44/0/2/0/0 => quoineから購入したものを受け取る用
//Payment => m/44/0/2/1/0 => 出金による支払いに利用、かつ、おつりも受け取る => TODO:ChangeTypeによってアドレスが変わってしまったら、どう運用するのか

// CreateAccount アカウント階層までのprivateKey及び publicKeyを生成する
func CreateAccount(conf *chaincfg.Params, seed []byte, actType enum.AccountType) (string, string, error) {

	//Master
	masterKey, err := hdkeychain.NewMaster(seed, conf)
	if err != nil {
		return "", "", err
	}
	//Purpose
	purpose, err := masterKey.Child(hdkeychain.HardenedKeyStart + uint32(PurposeTypeBIP44))
	//purpose, err := masterKey.Child(hdkeychain.HardenedKeyStart + uint32(PurposeTypeBIP49))
	if err != nil {
		return "", "", err
	}
	//CoinType TODO:切り替えが必要
	ct := uint32(enum.CoinTypeBitcoin)
	if conf.Name != string(enum.NetworkTypeMainNet) {
		ct = uint32(enum.CoinTypeTestnet)
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
func CreateKeysWithIndex(conf *chaincfg.Params, accountPrivateKey string, idxFrom, count uint32) ([]WalletKey, error) {
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

		// WIF　(compress: false)
		//wif, err := btcutil.NewWIF(privateKey, conf, false)
		// WIF　(compress: false)
		wif, err := btcutil.NewWIF(privateKey, conf, true)
		if err != nil {
			return nil, err
		}
		strPrivateKey := wif.String()

		// Address
		address, err := child.Address(conf)
		if err != nil {
			return nil, err
		}

		// p2sh-segwit
		p2shSegwit, err := getP2shSegwit(privateKey, conf)
		if err != nil {
			return nil, err
		}

		//address.String() とaddress.EncodeAddress()は結果として同じ
		walletKeys[i] = WalletKey{WIF: strPrivateKey, Address: address.String(), P2shSegwit: p2shSegwit.String(), FullPubKey: getFullPubKey(privateKey)}

		idxFrom++
	}

	return walletKeys, nil
}

//For only Debug
func experimentalKey() {
	// [Debug] different way to generate address
	//serializedKey := privateKey.PubKey().SerializeCompressed()
	//pubKeyAddr, err := btcutil.NewAddressPubKey(serializedKey, conf)
	//log.Println("address.String()", address.String())       //mySBc7pWWXjBUmAtjBY3sCdgnPAvAzwCoA
	//log.Println("pubKeyAddr.String()", pubKeyAddr.String()) //022c70901aac621c4436c4cb1f2daa8b9a6ff2c9d707b3f2639319d902679e1dfd
	//log.Println("pubKeyAddr.AddressPubKeyHash().String()", pubKeyAddr.AddressPubKeyHash().String()) //mySBc7pWWXjBUmAtjBY3sCdgnPAvAzwCoA
	//log.Println("getFullPubKey(privateKey)", getFullPubKey(privateKey)) //pubKeyAddr.String()とは微妙に異なる。。
	//log.Println(" ")
}

// p2sh-segwitを返す
func getP2shSegwit(privKey *btcec.PrivateKey, conf *chaincfg.Params) (*btcutil.AddressScriptHash, error) {
	// []byte
	publicKeyHash := btcutil.Hash160(privKey.PubKey().SerializeCompressed())
	segwitAddress, err := btcutil.NewAddressWitnessPubKeyHash(publicKeyHash, conf)
	if err != nil {
		return nil, errors.Errorf("btcutil.NewAddressWitnessPubKeyHash() error: %s", err)
	}
	logger.Debugf("segwitAddress: %s", segwitAddress)

	redeemScript, err := txscript.PayToAddrScript(segwitAddress)
	if err != nil {
		return nil, errors.Errorf("txscript.PayToAddrScript() error: %s", err)
	}
	logger.Debugf("redeemScript: %s", redeemScript)

	address, err := btcutil.NewAddressScriptHash(redeemScript, conf)
	if err != nil {
		return nil, errors.Errorf("btcutil.NewAddressScriptHash() error: %s", err)
	}
	logger.Debugf("address.String() %s", address.String())

	return address, nil
}

// getPubKey fullのPublic Keyを返す
func getFullPubKey(privKey *btcec.PrivateKey) string {
	//bPubKey := privKey.PubKey().SerializeCompressed()
	bPubKey := privKey.PubKey().SerializeUncompressed()
	//logger.Debugf("bPubKey hash: %s", btcutil.Hash160(bPubKey))

	hexPubKey := hex.EncodeToString(bPubKey)
	//logger.Debugf("hex.EncodeToString(bPubKey): %s", hexPubKey)

	//key *PublicKey
	//bHexPubKey, _ := hex.DecodeString(hexPubKey)
	//pubKey, _ := btcec.ParsePubKey(bHexPubKey, btcec.S256())

	return hexPubKey
}
