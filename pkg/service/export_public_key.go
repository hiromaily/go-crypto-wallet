package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/csv"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
)

func (w *Wallet) ExportPublicKey(accountType key.AccountType) error {
	//AccountType問わずexportは可能にしておく

	//DBから該当するpublic keyを取得
	pubKeys, err := w.DB.GetNotExportedPubKey(accountType)
	if err != nil {
		return errors.Errorf("key.GenerateSeed() error: %s", err)
	}

	//accountTypeから必要なファイルパスを取得
	at, err := w.DB.GetAccountTypeByID(accountType)
	if err != nil {
		return errors.Errorf("DB.GetAccountTypeByID() error: %s", err)
	}

	//CSVに書き出す
	fileName, err := csv.ExportPubKey(pubKeys, at.Type)
	if err != nil {
		return errors.Errorf("csv.ExportPubKey() error: %s", err)
	}
	logger.Infof("file name is %s", fileName)

	//DBの該当レコードをアップデート
	_, err = w.DB.UpdateIsExprotedPubKey(accountType, pubKeys, nil, true)
	if err != nil {
		return errors.Errorf("csv.UpdateIsExprotedPubKey() error: %s", err)
	}

	return nil
}
