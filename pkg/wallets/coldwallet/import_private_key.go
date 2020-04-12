package coldwallet

//Cold wallet

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/keystatus"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

// ImportPrivateKey 指定したAccountTypeに属するテーブルのis_imported_priv_keyがfalseのWIFをImportPrivKeyRescanする
// https://en.bitcoin.it/wiki/How_to_import_private_keys
// getaddressesbyaccount "" で内容を確認可能？？
func (w *ColdWallet) ImportPrivateKey(accountType account.AccountType) error {
	//TODO:remove it
	if w.wtype == types.WalletTypeWatchOnly {
		return errors.New("it's available on Coldwallet1, Coldwallet2")
	}

	//AccountType問わずimportは可能にしておく

	//DBから未登録のPrivateKey情報を取得する
	//WIFs, err := w.DB.GetNotImportedKeyWIF(accountType)
	accountKeyTable, err := w.storager.GetAllAccountKeyByKeyStatus(accountType, keystatus.KeyStatusGenerated) //key_status=0
	if err != nil {
		return errors.Errorf("key.GenerateSeed() error: %s", err)
	}

	if len(accountKeyTable) == 0 {
		w.logger.Info("No unimported Private Key")
		return nil
	}

	//var account string
	//if accountType != ctype.AccountTypeClient {
	//	account = string(accountType)
	//}
	account := string(accountType)

	//bitcoin APIにて登録をする
	for _, record := range accountKeyTable {
		w.logger.Debug(
			"call GetAllAccountKeyByKeyStatus()",
			zap.String("accountType", accountType.String()),
			zap.String("record.WalletAddress", record.WalletAddress),
			zap.String("record.WalletImportFormat", record.WalletImportFormat))
		wif, err := btcutil.DecodeWIF(record.WalletImportFormat)
		if err != nil {
			//ここでエラーが出るのであれば生成ロジックが抜本的に問題があるので、return
			return errors.Errorf(
				"fail to call btcutil.DecodeWIF(%s). WIF is invalid format.  error: %v", record.WalletImportFormat, err)
		}
		//rescanはいらないはず
		err = w.btc.ImportPrivKeyWithoutReScan(wif, account)
		//err = w.BTC.ImportPrivKeyWithoutReScan(wif, "")
		//err = w.BTC.ImportPrivKey(wif)
		if err != nil {
			//Bitcoin coreの状況によってエラーが返ることも想定する。よってエラー時はcontinue
			w.logger.Error("fail to call btc.ImportPrivKeyWithoutReScan()", zap.Error(err))
			continue
		}
		//update DB
		//_, err = w.DB.UpdateIsImprotedPrivKey(accountType, record.WalletImportFormat, nil, true)
		_, err = w.storager.UpdateKeyStatusByWIF(accountType, keystatus.KeyStatusImportprivkey, record.WalletImportFormat, nil, true)
		if err != nil {
			w.logger.Error(
				"fail to call btc.UpdateKeyStatusByWIF()",
				zap.String("accountType", accountType.String()),
				zap.String("KeyStatus", keystatus.KeyStatusImportprivkey.String()),
				zap.String("record.WalletImportFormat", record.WalletImportFormat),
				zap.Error(err))
		}

		//アドレスがbitcoin core walletに登録されているかチェック
		w.checkImportedAddress(record.WalletAddress, record.P2shSegwitAddress, record.FullPublicKey)
	}

	return nil
}

// checkImportedAddress check imported address
func (w *ColdWallet) checkImportedAddress(walletAddress, p2shSegwitAddress, fullPublicKey string) {
	w.logger.Info("checkImportedAddress")
	//if w.btc.Version() >= ctype.BTCVer17 {
	//	w.checkImportedAddressVer17(walletAddress, p2shSegwitAddress, fullPublicKey)
	//	return
	//}
	//For now, version17 is expected
	//getaddressinfo "address"
	addrInfo, err := w.btc.GetAddressInfo(walletAddress)
	if err != nil {
		w.logger.Error(
			"fail to call btc.GetAddressInfo()",
			zap.String("walletAddress", walletAddress),
			zap.Error(err))
	}
	w.logger.Debug(
		"account is found by wallet_address",
		zap.String("account", addrInfo.Label),
		zap.String("walletAddress", walletAddress))

	//2.getaccount address(p2sh_segwit_address)
	addrInfo, err = w.btc.GetAddressInfo(p2shSegwitAddress)
	if err != nil {
		w.logger.Error(
			"fail to call btc.GetAccount()",
			zap.String("p2shSegwitAddress", p2shSegwitAddress),
			zap.Error(err))
	}
	w.logger.Debug(
		"account is found by p2sh_segwit_address",
		zap.String("acount", addrInfo.Label),
		zap.String("p2shSegwitAddress", p2shSegwitAddress))

	//3.getaccount address(p2sh_segwit_address)
	addrInfo, err = w.btc.GetAddressInfo(p2shSegwitAddress)
	if err != nil {
		w.logger.Error(
			"fail to call btc.GetAccount()",
			zap.String("p2shSegwitAddress", p2shSegwitAddress),
			zap.Error(err))
	}
	w.logger.Debug(
		"account is found by p2sh_segwit_address",
		zap.String("acount", addrInfo.Label),
		zap.String("p2shSegwitAddress", p2shSegwitAddress))

	if addrInfo.Pubkey != fullPublicKey {
		w.logger.Error("generating pubkey logic is wrong")
	}
}

// checkImportedAddress addresssをチェックする (for bitcoin version 17)
func (w *ColdWallet) checkImportedAddressVer19(walletAddress, p2shSegwitAddress, fullPublicKey string) {
	w.logger.Info("checkImportedAddressVer19()")

}
