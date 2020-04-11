package wallets

//Cold wallet

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
)

// ImportPrivateKey 指定したAccountTypeに属するテーブルのis_imported_priv_keyがfalseのWIFをImportPrivKeyRescanする
// https://en.bitcoin.it/wiki/How_to_import_private_keys
// getaddressesbyaccount "" で内容を確認可能？？
func (w *Wallet) ImportPrivateKey(accountType account.AccountType) error {
	//TODO:remove it
	if w.wtype == WalletTypeWatchOnly {
		return errors.New("it's available on Coldwallet1, Coldwallet2")
	}

	//AccountType問わずimportは可能にしておく

	//DBから未登録のPrivateKey情報を取得する
	//WIFs, err := w.DB.GetNotImportedKeyWIF(accountType)
	accountKeyTable, err := w.storager.GetAllAccountKeyByKeyStatus(accountType, enum.KeyStatusGenerated) //key_status=0
	if err != nil {
		return errors.Errorf("key.GenerateSeed() error: %s", err)
	}

	if len(accountKeyTable) == 0 {
		w.logger.Info("No unimported Private Key")
		return nil
	}

	//var account string
	//if accountType != enum.AccountTypeClient {
	//	account = string(accountType)
	//}
	account := string(accountType)

	//bitcoin APIにて登録をする
	for _, record := range accountKeyTable {
		w.logger.Debugf("[%s] address: %s, WIF: %s", accountType, record.WalletAddress, record.WalletImportFormat)
		wif, err := btcutil.DecodeWIF(record.WalletImportFormat)
		if err != nil {
			//ここでエラーが出るのであれば生成ロジックが抜本的に問題があるので、return
			return errors.Errorf("WIF is invalid format. btcutil.DecodeWIF(%s) error: %v", record.WalletImportFormat, err)
		}
		//rescanはいらないはず
		w.logger.Debugf("BTC.ImportPrivKeyWithoutReScan(%s, %s)", wif, account)
		err = w.btc.ImportPrivKeyWithoutReScan(wif, account)
		//err = w.BTC.ImportPrivKeyWithoutReScan(wif, "")
		//err = w.BTC.ImportPrivKey(wif)
		if err != nil {
			//Bitcoin coreの状況によってエラーが返ることも想定する。よってエラー時はcontinue
			w.logger.Errorf("BTC.ImportPrivKeyWithoutReScan() error: %s", err)
			continue
		}
		//update DB
		//_, err = w.DB.UpdateIsImprotedPrivKey(accountType, record.WalletImportFormat, nil, true)
		_, err = w.storager.UpdateKeyStatusByWIF(accountType, enum.KeyStatusImportprivkey, record.WalletImportFormat, nil, true)
		if err != nil {
			w.logger.Errorf("BTC.UpdateKeyStatusByWIF(%s, %s, %s) error: %s", accountType, enum.KeyStatusImportprivkey, record.WalletImportFormat, err)
		}

		//アドレスがbitcoin core walletに登録されているかチェック
		w.checkImportedAddress(record.WalletAddress, record.P2shSegwitAddress, record.FullPublicKey)
	}

	return nil
}

// checkImportedAddress addresssをチェックする (for bitcoin version 16)
func (w *Wallet) checkImportedAddress(walletAddress, p2shSegwitAddress, fullPublicKey string) {
	if w.btc.Version() >= enum.BTCVer17 {
		w.checkImportedAddressVer17(walletAddress, p2shSegwitAddress, fullPublicKey)
		return
	}

	//1.getaccount address(wallet_address)
	account, err := w.btc.GetAccount(walletAddress)
	if err != nil {
		w.logger.Warnf("w.BTC.GetAccount(%s) error: %s", walletAddress, err)
		//for new version check
		w.checkImportedAddressVer17(walletAddress, p2shSegwitAddress, fullPublicKey)
		return
	}
	logger.Debugf("account[%s] is found by wallet_address:%s", account, walletAddress)

	//2.getaccount address(p2sh_segwit_address)
	account, err = w.btc.GetAccount(p2shSegwitAddress)
	if err != nil {
		w.logger.Errorf("w.BTC.GetAccount(%s) error: %s", p2shSegwitAddress, err)
	}
	w.logger.Debugf("account[%s] is found by p2sh_segwit_address:%s", account, p2shSegwitAddress)

	//3.check full_public_key by validateaddress retrieving it
	res, err := w.btc.ValidateAddress(p2shSegwitAddress)
	if err != nil {
		w.logger.Errorf("w.BTC.ValidateAddress(%s) error: %s", p2shSegwitAddress, err)
	}
	if res.PubKey != fullPublicKey {
		w.logger.Errorf("generating pubkey logic is wrong")
	}
}

// checkImportedAddress addresssをチェックする (for bitcoin version 17)
func (w *Wallet) checkImportedAddressVer17(walletAddress, p2shSegwitAddress, fullPublicKey string) {
	w.logger.Info("checkImportedAddressVer17()")

	//getaddressinfo "address"
	addrInfo, err := w.btc.GetAddressInfo(walletAddress)
	if err != nil {
		w.logger.Errorf("w.BTC.GetAddressInfo(%s) error: %s", walletAddress, err)
	}
	w.logger.Debugf("account[%s] is found by wallet_address:%s", addrInfo.Label, walletAddress)

	//2.getaccount address(p2sh_segwit_address)
	addrInfo, err = w.btc.GetAddressInfo(p2shSegwitAddress)
	if err != nil {
		w.logger.Errorf("w.BTC.GetAccount(%s) error: %s", p2shSegwitAddress, err)
	}
	w.logger.Debugf("account[%s] is found by p2sh_segwit_address:%s", addrInfo.Label, p2shSegwitAddress)

	//3.getaccount address(p2sh_segwit_address)
	addrInfo, err = w.btc.GetAddressInfo(p2shSegwitAddress)
	if err != nil {
		w.logger.Errorf("w.BTC.GetAccount(%s) error: %s", p2shSegwitAddress, err)
	}
	w.logger.Debugf("account[%s] is found by p2sh_segwit_address:%s", addrInfo.Label, p2shSegwitAddress)

	if addrInfo.Pubkey != fullPublicKey {
		w.logger.Errorf("generating pubkey logic is wrong")
	}
}
