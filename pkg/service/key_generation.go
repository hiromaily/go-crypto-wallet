package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/pkg/errors"
)

//Seed生成は完全に分離したほうがいい
//1.Seedの生成+DBに登録
//2.Multisig Keyの生成+DBに登録(承認用は端末を分けて管理しないと意味がないかも)

//CreateMultiSig(addmultisigaddress)にwalletにmultisig用のprivate keyを登録する
//これのパラメータには、multisigしないと送金許可しないアドレス(receipt, payment)+承認用のアドレスをセット
//これによって、生成されたアドレスから、送金する場合、パラメータにセットしたアドレスに紐づく秘密鍵が必要
//payment,receiptのアドレスは、実際には、addmultisigaddressによって生成されたアドレスに置き換えられる。

//含まれるもの
//coldwallet1: client, receipt, payment
//coldwallet2: multisig address
//TODO: =>どこのマシンで、addmultisigaddressを行う？？

//https://www.slideshare.net/ssusere174e3/ss-33733512

//3.Client Keyの生成+DBに登録
//4.Receipt Keyの生成 + Multisig対応 + DBに登録 (1日1Key消費するイメージ)
//5.Payment Keyの生成+ Multisig + DBに登録 (1日1Key消費するイメージ)

// InitialKeyGeneration 初回生成用のfunc
func (w *Wallet) InitialKeyGeneration() error {

	// Seed
	seed, err := w.retrieveSeed()
	//seed, err := w.GenerateSeed()
	if err != nil {
		return errors.Errorf("retrieveSeed() error: %v", err)
	}

	//accounts := []key.AccountType{
	//	key.AccountTypeClient,
	//	key.AccountTypeReceipt,
	//	key.AccountTypePayment,
	//	key.AccountTypeMultisig,
	//}

	// TODO:初回の生成数をどのように調整し、決定するか？とりあえず固定
	// アカウント(Client)を生成
	_, err = w.GenerateAccountKey(key.AccountTypeClient, seed, 10000)
	if err != nil {
		return errors.Errorf("GenerateClientAccount() error: %v", err)
	}

	// アカウント(Receipt)を生成
	//_, err = w.GenerateReceiptAccount(seed, 0, 100)
	//if err != nil {
	//	return errors.Errorf("GenerateReceiptAccount() error: %v", err)
	//}

	// アカウント(Payment)を生成
	//_, err = w.GeneratePaymentAccount(seed, 0, 100)
	//if err != nil {
	//	return errors.Errorf("GeneratePaymentAccount() error: %v", err)
	//}

	// アカウント(Authorization)を生成
	//_, err = w.GenerateAuthorizationAccount(seed, 0, 3)
	//if err != nil {
	//	return errors.Errorf("GenerateMultisigAccount() error: %v", err)
	//}

	return nil
}

// GenerateSeed seedを生成する
func (w *Wallet) GenerateSeed() ([]byte, error) {

	//seed, err := w.DB.GetSeedOne()
	//if err == nil && seed.Seed != "" {
	//	logger.Info("seed have already been generated")
	//	return key.SeedToByte(seed.Seed)
	//}
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

	// DBにseed情報を登録
	_, err = w.DB.InsertSeed(strSeed, nil, true)
	if err != nil {
		return nil, errors.Errorf("key.InsertSeed() error: %s", err)
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

	return nil, errors.Errorf("DB.GetSeedOne() error: %v", err)
}

// GenerateAccountKey AccountType属性のアカウントKeyを生成する
func (w *Wallet) GenerateAccountKey(accountType key.AccountType, seed []byte, count uint32) ([]key.WalletKey, error) {
	//現在のindexを取得
	idx, err := w.DB.GetMaxIndex(accountType)
	if err != nil {
		//return nil, errors.Errorf("DB.GetMaxClientIndex() error: %s", err)
		idx = 0
	} else {
		idx++
	}
	logger.Infof("idx: %d", idx)

	return w.generateAccountKey(accountType, seed, uint32(idx), count)
}

// generateKey AccountType属性のアカウントKeyを生成する
func (w *Wallet) generateAccountKey(accountType key.AccountType, seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
	// HDウォレットのkeyを生成する
	walletKeys, err := w.generateAccountKeyData(accountType, seed, idxFrom, count)
	if err != nil {
		return nil, errors.Errorf("key.generateAccount(AccountTypeClient) error: %s", err)
	}

	// keyTypeを取得
	keyID, err := w.getKeyTypeByAccount(accountType)
	if err != nil {
		return nil, errors.Errorf("getKeyTypeByAccount(AccountTypeClient) error: %s", err)
	}

	// DBにClientAccountのKey情報を登録
	accountKeyClients := make([]model.AccountKeyTable, len(walletKeys))
	for idx, key := range walletKeys {
		accountKeyClients[idx] = model.AccountKeyTable{
			WalletAddress:      key.Address,
			WalletImportFormat: key.WIF,
			KeyType:            keyID,
			Idx:                idxFrom,
		}
		idxFrom++
	}
	err = w.DB.InsertAccountKeyTable(accountType, accountKeyClients, nil, true)
	if err != nil {
		return nil, errors.Errorf("DB.InsertAccountKeyClient() error: %s", err)
	}

	return walletKeys, err
}

// generateKeyData AccountType属性のアカウントKeyを生成する
func (w *Wallet) generateAccountKeyData(accountType key.AccountType, seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
	// key生成
	priv, _, err := key.CreateAccount(w.BTC.GetChainConf(), seed, accountType)
	if err != nil {
		return nil, errors.Errorf("key.CreateAccount() error: %s", err)
	}

	walletKeys, err := key.CreateKeysWithIndex(w.BTC.GetChainConf(), priv, idxFrom, count)
	if err != nil {
		return nil, errors.Errorf("key.CreateKeysWithIndex() error: %s", err)
	}

	return walletKeys, nil
}

// GenerateReceiptAccount Receiptアカウントを生成する
//func (w *Wallet) GenerateReceiptAccount(seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
//	walletKeys, err := w.generateAccount(seed, idxFrom, count, key.AccountTypeReceipt)
//	if err != nil {
//		return nil, errors.Errorf("key.generateAccount(AccountTypeReceipt) error: %s", err)
//	}
//
//	// TODO:DBにReceiptAccountのKey情報を登録
//
//	return walletKeys, err
//}

// GeneratePaymentAccount Paymentアカウントを生成する
//func (w *Wallet) GeneratePaymentAccount(seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
//	walletKeys, err := w.generateAccount(seed, idxFrom, count, key.AccountTypePayment)
//	if err != nil {
//		return nil, errors.Errorf("key.generateAccount(AccountTypePayment) error: %s", err)
//	}
//
//	// TODO:DBにPaymentAccountのKey情報を登録
//
//	return walletKeys, err
//}

// GenerateAuthorizationAccount Paymentアカウントを生成する
//func (w *Wallet) GenerateAuthorizationAccount(seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
//	walletKeys, err := w.generateAccount(seed, idxFrom, count, key.AccountTypeAuthorization)
//	if err != nil {
//		return nil, errors.Errorf("key.generateAccount(AccountTypeAuthorization) error: %s", err)
//	}
//
//	// TODO:DBにAuthorizationAccountのKey情報を登録
//
//	return walletKeys, err
//}

func (w *Wallet) getKeyTypeByAccount(accountType key.AccountType) (uint8, error) {
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
