package service

//Cold wallet

import (
	"github.com/btcsuite/btcutil"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
)

// ImportPrivateKey 指定したAccountTypeに属するテーブルのis_imported_priv_keyがfalseのWIFをImportPrivKeyRescanする
// https://en.bitcoin.it/wiki/How_to_import_private_keys
// TODO: import後は再起動が必要かもしれない
// getaddressesbyaccount "" で内容を確認可能？？
func (w *Wallet) ImportPrivateKey(accountType enum.AccountType) error {
	//AccountType問わずimportは可能にしておく

	//DBから未登録のPrivateKey情報を取得する
	//WIFs, err := w.DB.GetNotImportedKeyWIF(accountType)
	accountKeyTable, err := w.DB.GetAllNotImportedKey(accountType)
	if err != nil {
		return errors.Errorf("key.GenerateSeed() error: %s", err)
	}

	if len(accountKeyTable) == 0 {
		logger.Info("No unimported Private Key")
		return nil
	}

	var account string
	if accountType != enum.AccountTypeClient {
		account = string(accountType)
	}

	//bitcoin APIにて登録をする
	for _, record := range accountKeyTable {
		logger.Debugf("[%s] address: %s, WIF: %s", accountType, record.WalletAddress, record.WalletImportFormat)
		wif, err := btcutil.DecodeWIF(record.WalletImportFormat)
		if err != nil {
			//ここでエラーが出るのであれば生成ロジックが抜本的に問題があるので、return
			return errors.Errorf("WIF is invalid format. btcutil.DecodeWIF(%s) error: %v", record.WalletImportFormat, err)
		}
		//TODO:rescanはいらないはず
		logger.Debugf("BTC.ImportPrivKeyWithoutReScan(%s, %s)", wif, account)
		err = w.BTC.ImportPrivKeyWithoutReScan(wif, account)
		//err = w.BTC.ImportPrivKeyWithoutReScan(wif, "")
		//err = w.BTC.ImportPrivKey(wif)
		if err != nil {
			//Bitcoin coreの状況によってエラーが返ることも想定する。よってエラー時はcontinue
			logger.Errorf("BTC.ImportPrivKeyWithoutReScan() error: %v", err)
			continue
		}
		//update DB
		_, err = w.DB.UpdateIsImprotedPrivKey(accountType, record.WalletImportFormat, nil, true)
		if err != nil {
			logger.Errorf("BTC.UpdateIsImprotedPrivKey(%s, %s) error: %v", accountType, record.WalletImportFormat, err)
		}

		//TODO:getaddressesbyaccount "receipt"/"payment"で、登録されたアドレスが表示されるかチェック
		if accountType != enum.AccountTypeClient {
			//1.getaccount address(wallet_address)
			account, err := w.BTC.GetAccount(record.WalletAddress)
			if err != nil {
				logger.Errorf("w.BTC.GetAccount(%s) error: %v", record.WalletAddress, err)
			}
			logger.Debugf("account[%s] is found by wallet_address:%s", account, record.WalletAddress)

			//2.getaccount address(p2sh_segwit_address)
			account, err = w.BTC.GetAccount(record.P2shSegwitAddress)
			if err != nil {
				logger.Errorf("w.BTC.GetAccount(%s) error: %v", record.P2shSegwitAddress, err)
			}
			logger.Debugf("account[%s] is found by p2sh_segwit_address:%s", account, record.P2shSegwitAddress)

			//3.TODO check full_public_key by validateaddress retrieving it

		}
	}

	//for _, strWIF := range WIFs {
	//	logger.Debugf("[%s] strWIF: %s", accountType, strWIF)
	//	wif, err := btcutil.DecodeWIF(strWIF)
	//	if err != nil {
	//		//ここでエラーが出るのであれば生成ロジックが抜本的に問題があるので、return
	//		return errors.Errorf("WIF is invalid format. btcutil.DecodeWIF(%s) error: %v", strWIF, err)
	//	}
	//	//TODO:rescanはいらないはず
	//	logger.Debugf("BTC.ImportPrivKeyWithoutReScan(%s, %s)", wif, account)
	//	err = w.BTC.ImportPrivKeyWithoutReScan(wif, account)
	//	//err = w.BTC.ImportPrivKeyWithoutReScan(wif, "")
	//	//err = w.BTC.ImportPrivKey(wif)
	//	if err != nil {
	//		//Bitcoin coreの状況によってエラーが返ることも想定する。よってエラー時はcontinue
	//		logger.Errorf("BTC.ImportPrivKeyWithoutReScan() error: %v", err)
	//		continue
	//	}
	//	//update DB
	//	_, err = w.DB.UpdateIsImprotedPrivKey(accountType, strWIF, nil, true)
	//	if err != nil {
	//		logger.Errorf("BTC.UpdateIsImprotedPrivKey(%s, %s) error: %v", accountType, strWIF, err)
	//	}
	//}

	return nil
}
