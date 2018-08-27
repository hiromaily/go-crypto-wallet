package service

import (
	"github.com/btcsuite/btcutil"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
)

// ImportPrivateKey 指定したAccountTypeに属するテーブルのis_imported_priv_keyがfalseのWIFをImportPrivKeyRescanする
func (w *Wallet) ImportPrivateKey(accountType key.AccountType) error {
	//DBから未登録のPrivateKey情報を取得する
	WIFs, err := w.DB.GetNotImportedKeyWIF(accountType)
	if err != nil {
		return errors.Errorf("key.GenerateSeed() error: %s", err)
	}

	//bitcoin APIにて登録をする
	for _, strWIF := range WIFs {
		wif, err := btcutil.DecodeWIF(strWIF)
		if err != nil {
			//ここでエラーが出るのであれば生成ロジックが抜本的に問題があるので、return
			return errors.Errorf("WIF is invalid format. btcutil.DecodeWIF(%s) error: %v", strWIF, err)
		}
		err = w.BTC.ImportPrivKeyWithoutReScan(wif, "")
		if err != nil {
			//Bitcoin coreの状況によってエラーが返ることも想定する。よってエラー時はcontinue
			logger.Errorf("BTC.ImportPrivKeyWithoutReScan() error: %v", err)
			continue
		}
		//update DB
		_, err = w.DB.UpdateIsImprotedPrivKey(accountType, strWIF, nil, true)
		if err != nil {
			logger.Errorf("BTC.UpdateIsImprotedPrivKey(%s, %s) error: %v", accountType, strWIF, err)
		}
	}

	return nil
}
