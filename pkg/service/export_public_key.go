package service

//Cold wallet

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
)

//ExportPublicKey publicのアドレスをcsvとして出力する
//TODO:watch only walletにセットするアドレスは、clientの場合は、wallet_address, receipt/paymentの場合、`wallet_multisig_address`
func (w *Wallet) ExportPublicKey(accountType enum.AccountType, isMultisig bool) error {
	//AccountType問わずexportは可能にしておくか？はじめ厳しくしておき、徐々に設定を緩める方向で

	//DBから該当するpublic keyを取得
	pubKeys, err := w.DB.GetPubkeyNotExportedPubKey(accountType, isMultisig)
	if err != nil {
		return errors.Errorf("key.GetPubkeyNotExportedPubKey() error: %s", err)
	}

	if len(pubKeys) == 0 {
		logger.Info("no public key in table")
		return nil
	}

	//CSVに書き出す
	fileName, err := key.ExportPubKey(pubKeys, string(accountType))
	if err != nil {
		return errors.Errorf("key.ExportPubKey() error: %s", err)
	}
	logger.Infof("file name is %s", fileName)

	//DBの該当レコードをアップデート
	_, err = w.DB.UpdateIsExprotedPubKey(accountType, pubKeys, isMultisig, nil, true)
	if err != nil {
		return errors.Errorf("DB.UpdateIsExprotedPubKey() error: %s", err)
	}

	return nil
}

//ExportFullPublicKey full publicのアドレスをcsvとして出力する
func (w *Wallet) ExportFullPublicKey(accountType enum.AccountType) error {
	//AccountType問わずexportは可能にしておくか？はじめ厳しくしておき、徐々に設定を緩める方向で

	//DBから該当するpublic keyを取得
	pubKeys, err := w.DB.GetFullPubkeyNotExportedPubKey(accountType)
	if err != nil {
		return errors.Errorf("key.GetPubkeyNotExportedPubKey() error: %s", err)
	}

	if len(pubKeys) == 0 {
		logger.Info("no full public key in table")
		return nil
	}

	//CSVに書き出す
	fileName, err := key.ExportPubKey(pubKeys, string(accountType))
	if err != nil {
		return errors.Errorf("key.ExportPubKey() error: %s", err)
	}
	logger.Infof("file name is %s", fileName)

	return nil
}

//ExportAllKeyTable publicのアドレスをcsvとして出力する
//TODO:一旦使用中止
func (w *Wallet) ExportAllKeyTable(accountType enum.AccountType) error {
	//AccountType問わずexportは可能にしておく

	//DBから該当するpublic keyを取得
	accountKeyTable, err := w.DB.GetAllNotExportedPubKey(accountType)
	if err != nil {
		return errors.Errorf("key.GetAllNotExportedPubKey() error: %s", err)
	}

	if len(accountKeyTable) == 0 {
		logger.Info("no public key in table")
		return nil
	}

	//CSVに書き出す
	fileName, err := key.ExportAccountKeyTable(accountKeyTable, string(accountType))
	if err != nil {
		return errors.Errorf("key.ExportAccountKeyTable() error: %s", err)
	}
	logger.Infof("file name is %s", fileName)

	return nil
}
