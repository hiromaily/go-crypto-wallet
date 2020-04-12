package wallet

import (
	"strings"

	"github.com/bookerzzz/grok"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
)

//CSVのpublickeyをDBにimportする、このとき、clientの場合はaccount無し
//importしたclientをBitcoin core APIを通じて、walletにimportする
// ImportPublicKeyForWatchWallet csvファイルからpublicアドレスをimportする for WatchOnyWallet
func (w *Wallet) ImportPublicKey(fileName string, accountType account.AccountType, isRescan bool) error {
	//ファイル読み込み
	pubKeys, err := key.ImportPubKey(fileName)
	if err != nil {
		return errors.Errorf("key.ImportPubKey() error: %s", err)
	}

	//[]AccountPublicKeyTable //import時にアカウント名をcoreのwalletに登録する
	acnt := string(accountType)
	//if accountType == enum.AccountTypeClient {
	//	account = ""
	//}

	var pubKeyData []walletrepo.AccountPublicKeyTable
	for _, key := range pubKeys {
		inner := strings.Split(key, ",")
		grok.Value(inner)
		var addr string
		if accountType == account.AccountTypeClient {
			addr = inner[1] //p2sh_segwit_address
		} else {
			addr = inner[3] //wallet_import_format
		}

		//Bitcoin core APIから`importaddress`をcallする
		//TODO:1台のPCで検証しているときなど、すでにimport済の場合はエラーが出る
		err := w.btc.ImportAddressWithLabel(addr, acnt, isRescan) //基本falseのはず
		if err != nil {
			//-4: The wallet already contains the private key for this address or script
			w.logger.Error(
				"btc.ImportAddressWithLabel(%s) error: %s",
				zap.String("address", addr),
				zap.Error(err))
			continue
		}

		pubKeyData = append(pubKeyData, walletrepo.AccountPublicKeyTable{
			WalletAddress: addr,
			Account:       acnt,
		})

		//watch only walletとして追加されているかチェックする
		w.checkImportedPublicAddress(addr)
	}

	//DBにInsert
	err = w.storager.InsertAccountPubKeyTable(accountType, pubKeyData, nil, true)
	if err != nil {
		return errors.Errorf("DB.InsertAccountPubKeyTable() error: %s", err)
		//TODO:これが失敗したら、どうやって、登録済みのデータを再度Insertするか？再度実行すればOKのはず
	}

	return nil
}

//checkImportedPublicAddress watch only walletとして追加されているかチェックする
func (w *Wallet) checkImportedPublicAddress(addr string) {
	if w.btc.Version() >= enum.BTCVer17 {
		w.checkImportedPublicAddressVer17(addr)
		return
	}

	//1.getaccount address(wallet_address)
	account, err := w.btc.GetAccount(addr)
	if err != nil {
		w.logger.Error(
			"w.btc.GetAccount()",
			zap.String("address", addr),
			zap.Error(err))
	}
	w.logger.Debug(
		"account is found",
		zap.String("account", account),
		zap.String("address", addr))

	//2.check full_public_key by validateaddress retrieving it
	res, err := w.btc.ValidateAddress(addr)
	if err != nil {
		w.logger.Error(
			"w.btc.ValidateAddress()",
			zap.String("address", addr),
			zap.Error(err))
	}
	grok.Value(res)
	//watch only walletを想定している
	if !res.IsWatchOnly {
		w.logger.Error("this address must be watch only wallet")
	}

}

//checkImportedPublicAddressVer17 watch only walletとして追加されているかチェックする (for bitcoin version 17)
func (w *Wallet) checkImportedPublicAddressVer17(addr string) {
	w.logger.Info("checkImportedPublicAddressVer17()")

	//getaddressinfo "address"
	addrInfo, err := w.btc.GetAddressInfo(addr)
	if err != nil {
		w.logger.Error(
			"w.btc.GetAddressInfo()",
			zap.String("address", addr),
			zap.Error(err))
	}
	w.logger.Debug("account is found",
		zap.String("account", addrInfo.Label),
		zap.String("address", addr))

	//watch only walletを想定している
	if !addrInfo.Iswatchonly {
		w.logger.Error("this address must be watch only wallet")
	}
}
