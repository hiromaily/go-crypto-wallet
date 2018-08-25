package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/pkg/errors"
)

// InitialKeyGeneration
func (w *Wallet) InitialKeyGeneration() error {
	// Seed
	seed, err := w.GenerateSeed()
	if err != nil {
		return errors.Errorf("GenerateSeed() error: %v", err)
	}

	//accounts := []key.AccountType{
	//	key.AccountTypeClient,
	//	key.AccountTypeReceipt,
	//	key.AccountTypePayment,
	//	key.AccountTypeMultisig,
	//}

	// TODO:初回の生成数をどのように調整し、決定するか？とりあえず固定
	// アカウント(Client)を生成
	_, err = w.GenerateClientAccount(seed, 0, 10000)
	if err != nil {
		return errors.Errorf("GenerateClientAccount() error: %v", err)
	}

	// アカウント(Receipt)を生成
	_, err = w.GenerateReceiptAccount(seed, 0, 100)
	if err != nil {
		return errors.Errorf("GenerateReceiptAccount() error: %v", err)
	}

	// アカウント(Payment)を生成
	_, err = w.GeneratePaymentAccount(seed, 0, 100)
	if err != nil {
		return errors.Errorf("GeneratePaymentAccount() error: %v", err)
	}

	// アカウント(Multisig)を生成
	_, err = w.GenerateMultisigAccount(seed, 0, 3)
	if err != nil {
		return errors.Errorf("GenerateMultisigAccount() error: %v", err)
	}

	return nil
}

// GenerateSeed seedを生成する
func (w *Wallet) GenerateSeed() ([]byte, error) {

	// 生成
	bSeed, err := key.GenerateSeed()
	if err != nil {
		return nil, errors.Errorf("key.GenerateSeed() error: %s", err)
	}

	// TODO:DBにseed情報を登録

	return bSeed, nil
}

func (w *Wallet) retrieveSeed() ([]byte, error) {
	// TODO:DBからseed情報を登録

	return nil, nil
}

// generateAccount アカウントを生成する
func (w *Wallet) generateAccount(seed []byte, idxFrom, count uint32, account key.AccountType) ([]key.WalletKey, error) {
	// アカウント(Client)を生成
	priv, _, err := key.CreateAccount(w.BTC.GetChainConf(), seed, account)
	if err != nil {
		return nil, errors.Errorf("key.CreateAccount() error: %s", err)
	}

	walletKeys, err := key.CreateKeysWithIndex(w.BTC.GetChainConf(), priv, idxFrom, count)
	if err != nil {
		return nil, errors.Errorf("key.CreateKeysWithIndex() error: %s", err)
	}

	return walletKeys, nil
}

// GenerateClientAccount Clientアカウントを生成する
func (w *Wallet) GenerateClientAccount(seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
	walletKeys, err := w.generateAccount(seed, idxFrom, count, key.AccountTypeClient)
	if err != nil {
		return nil, errors.Errorf("key.generateAccount(AccountTypeClient) error: %s", err)
	}

	// TODO:DBにClientAccountのKey情報を登録

	return walletKeys, err
}

// GenerateReceiptAccount Receiptアカウントを生成する
func (w *Wallet) GenerateReceiptAccount(seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
	walletKeys, err := w.generateAccount(seed, idxFrom, count, key.AccountTypeReceipt)
	if err != nil {
		return nil, errors.Errorf("key.generateAccount(AccountTypeReceipt) error: %s", err)
	}

	// TODO:DBにReceiptAccountのKey情報を登録

	return walletKeys, err
}

// GeneratePaymentAccount Paymentアカウントを生成する
func (w *Wallet) GeneratePaymentAccount(seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
	walletKeys, err := w.generateAccount(seed, idxFrom, count, key.AccountTypePayment)
	if err != nil {
		return nil, errors.Errorf("key.generateAccount(AccountTypePayment) error: %s", err)
	}

	// TODO:DBにPaymentAccountのKey情報を登録

	return walletKeys, err
}

// GenerateMultisigAccount Paymentアカウントを生成する
func (w *Wallet) GenerateMultisigAccount(seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
	walletKeys, err := w.generateAccount(seed, idxFrom, count, key.AccountTypeMultisig)
	if err != nil {
		return nil, errors.Errorf("key.generateAccount(AccountTypeMultisig) error: %s", err)
	}

	// TODO:DBにMultisigAccountのKey情報を登録

	return walletKeys, err
}
