package service

//Cold wallet

import (
	"github.com/hiromaily/go-bitcoin/pkg/csv"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
)

//ExportPublicKey publicのアドレスをcsvとして出力する
//TODO:watch only walletにセットするアドレスは、clientの場合は、wallet_address, receipt/paymentの場合、`wallet_multisig_address`
func (w *Wallet) ExportPublicKey(accountType enum.AccountType, isMultisig bool) error {
	//AccountType問わずexportは可能にしておく

	//DBから該当するpublic keyを取得
	pubKeys, err := w.DB.GetNotExportedPubKey(accountType, isMultisig)
	if err != nil {
		return errors.Errorf("key.GenerateSeed() error: %s", err)
	}

	//accountTypeから必要なファイルパスを取得
	//at, err := w.DB.GetAccountTypeByID(accountType)
	//if err != nil {
	//	return errors.Errorf("DB.GetAccountTypeByID() error: %s", err)
	//}

	//CSVに書き出す
	fileName, err := csv.ExportPubKey(pubKeys, string(accountType))
	if err != nil {
		return errors.Errorf("csv.ExportPubKey() error: %s", err)
	}
	logger.Infof("file name is %s", fileName)

	if len(pubKeys) == 0 {
		logger.Info("no public key in table")
		return nil
	}

	//DBの該当レコードをアップデート
	_, err = w.DB.UpdateIsExprotedPubKey(accountType, pubKeys,isMultisig, nil, true)
	if err != nil {
		return errors.Errorf("csv.UpdateIsExprotedPubKey() error: %s", err)
	}

	return nil
}
