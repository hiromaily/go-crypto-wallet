package service

//Cold wallet

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/pkg/errors"
)

//1.Seedの生成+DBに登録
//2.Multisig Keyの生成+DBに登録(承認用は端末を分けて管理しないと意味がないかも)

//CreateMultiSig(addmultisigaddress)にwalletにmultisig用のprivate keyを登録する
//これのパラメータには、multisigしないと送金許可しないアドレス(receipt, payment)+承認用のアドレスをセット
//これによって、生成されたアドレスから送金する場合、パラメータにセットしたアドレスに紐づく秘密鍵が必要
//payment,receiptのアドレスは、実際には、addmultisigaddressによって生成されたアドレスに置き換えられる。

//3.Client Keyの生成+DBに登録
//4.Receipt Keyの生成 + Multisig対応 + DBに登録 (1日1Key消費するイメージ)
//5.Payment Keyの生成+ Multisig + DBに登録 (1日1Key消費するイメージ)

// GenerateSeed seedを生成する
func (w *Wallet) GenerateSeed() ([]byte, error) {

	bSeed, err := w.retrieveSeed()
	if err == nil {
		return bSeed, nil
	}

	// seed生成
	bSeed, err = key.GenerateSeed()
	if err != nil {
		return nil, errors.Errorf("key.GenerateSeed() error: %s", err)
	}
	strSeed := key.SeedToString(bSeed)

	// set default seed
	if w.Env == enum.EnvDev && w.Seed != "" {
		strSeed = w.Seed
	}

	// DBにseed情報を登録
	_, err = w.DB.InsertSeed(strSeed, nil, true)
	if err != nil {
		return nil, errors.Errorf("DB.InsertSeed() error: %s", err)
	}

	return bSeed, nil
}

func (w *Wallet) retrieveSeed() ([]byte, error) {
	// DBからseed情報を登録
	seed, err := w.DB.GetSeedOne()
	if err == nil && seed.Seed != "" {
		logger.Info("seed have already been generated")
		return key.SeedToByte(seed.Seed)
	}

	return nil, errors.Errorf("DB.GetSeedOne() error: %s", err)
}

// GenerateAccountKey AccountType属性のアカウントKeyを生成する
// TODO:AccountTypeAuthorizationのときは、レコードがある場合は追加できないようにしたほうがいい？？
func (w *Wallet) GenerateAccountKey(accountType enum.AccountType, coinType enum.CoinType, seed []byte, count uint32) ([]key.WalletKey, error) {
	//現在のindexを取得
	idx, err := w.DB.GetMaxIndexOnAccountKeyTable(accountType)
	if err != nil {
		idx = 0
	} else {
		idx++
	}
	logger.Infof("idx: %d", idx)

	return w.generateAccountKey(accountType, coinType, seed, uint32(idx), count)
}

// generateKey AccountType属性のアカウントKeyを生成する
func (w *Wallet) generateAccountKey(accountType enum.AccountType, coinType enum.CoinType, seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
	// HDウォレットのkeyを生成する
	walletKeys, err := w.generateAccountKeyData(accountType, coinType, seed, idxFrom, count)
	if err != nil {
		return nil, errors.Errorf("key.generateAccountKeyData(AccountTypeClient) error: %s", err)
	}

	// keyTypeを取得
	keyID, err := w.getKeyTypeByAccount(accountType)
	if err != nil {
		return nil, errors.Errorf("getKeyTypeByAccount(AccountTypeClient) error: %s", err)
	}

	// Account
	var account string
	if accountType != enum.AccountTypeClient {
		account = string(accountType)
	}

	// DBにClientAccountのKey情報を登録
	accountKeyClients := make([]model.AccountKeyTable, len(walletKeys))
	for idx, key := range walletKeys {
		accountKeyClients[idx] = model.AccountKeyTable{
			WalletAddress:         key.Address,
			P2shSegwitAddress:     key.P2shSegwit,
			FullPublicKey:         key.FullPubKey,
			WalletMultisigAddress: "",
			RedeemScript:          key.RedeemScript,
			WalletImportFormat:    key.WIF,
			Account:               account,
			KeyType:               keyID,
			Idx:                   idxFrom,
		}
		idxFrom++
	}
	err = w.DB.InsertAccountKeyTable(accountType, accountKeyClients, nil, true)
	if err != nil {
		return nil, errors.Errorf("DB.InsertAccountKeyTable() error: %s", err)
	}

	return walletKeys, err
}

// generateKeyData AccountType属性のアカウントKeyを生成する
func (w *Wallet) generateAccountKeyData(accountType enum.AccountType, coinType enum.CoinType, seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
	// Keyオブジェクト
	keyData := key.NewKey(coinType, w.BTC.GetChainConf())

	// key生成
	priv, _, err := keyData.CreateAccount(seed, accountType)
	if err != nil {
		return nil, errors.Errorf("key.CreateAccount() error: %s", err)
	}

	walletKeys, err := keyData.CreateKeysWithIndex(priv, idxFrom, count)
	if err != nil {
		return nil, errors.Errorf("key.CreateKeysWithIndex() error: %s", err)
	}

	return walletKeys, nil
}

func (w *Wallet) getKeyTypeByAccount(accountType enum.AccountType) (uint8, error) {
	//accountType:0
	//coin_typeを取得
	ct := key.CoinTypeBitcoin
	if w.BTC.GetChainConf().Name != string(enum.NetworkTypeMainNet) {
		ct = key.CoinTypeTestnet
	}
	keyType, err := w.DB.GetKeyTypeByCoinAndAccountType(ct, accountType)
	if err != nil {
		return 0, errors.Errorf("DB.GetKeyTypeByCoinAndAccountType() error: %s", err)
	}

	return keyType.ID, nil
}
