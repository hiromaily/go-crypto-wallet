package key

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
)

//PurposeType BIP44は44固定
type PurposeType uint32

// purpose
const (
	PurposeTypeBIP44 PurposeType = 44 //BIP44
)

//CoinType コインの種類
type CoinType uint32

// coin_type
const (
	CoinTypeBitcoin CoinType = 0 //Bitcoin
	CoinTypeTestnet CoinType = 1 //Testnet
)

//AccountType 利用目的
type AccountType uint32

// account
const (
	AccountTypeClient   AccountType = 0 //ユーザーの入金受付用アドレス
	AccountTypeReceipt  AccountType = 1 //入金を受け付けるアドレス用
	AccountTypePayment  AccountType = 2 //出金時に支払いをするアドレス
	AccountTypeMultisig AccountType = 3 //マルチシグアドレスのための承認アドレス
)

//TODO:同じアドレスを使い回すと、アドレスから総額情報がバレて危険
//よって、内部利用のアドレスは毎回使い捨てにすること

//ChangeType 受け取り階層
type ChangeType uint32

// change
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
func CreateAccount(conf *chaincfg.Params, seed []byte, actType AccountType) (string, string, error) {

	//Master
	masterKey, err := hdkeychain.NewMaster(seed, conf)
	if err != nil {
		return "", "", err
	}
	//Purpose
	purpose, err := masterKey.Child(hdkeychain.HardenedKeyStart + uint32(PurposeTypeBIP44))
	if err != nil {
		return "", "", err
	}
	//CoinType
	coinType, err := purpose.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return "", "", err
	}
	//Account
	account, err := coinType.Child(hdkeychain.HardenedKeyStart + uint32(actType))
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
	max := idxFrom + count
	for i := uint32(idxFrom); i < max; i++ {
		child, err := change.Child(i)
		if err != nil {
			return nil, err
		}

		// privateKey
		privateKey, err := child.ECPrivKey()
		if err != nil {
			return nil, err
		}

		// WIF
		wif, err := btcutil.NewWIF(privateKey, conf, false)
		if err != nil {
			return nil, err
		}
		strPrivateKey := wif.String()

		// Address
		address, err := child.Address(conf)
		if err != nil {
			return nil, err
		}

		walletKeys[i] = WalletKey{WIF: strPrivateKey, Address: address.String()}
	}

	return walletKeys, nil
}
