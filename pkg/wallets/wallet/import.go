package wallet

import (
	"strings"

	"github.com/bookerzzz/grok"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
)

//CSVのpublickeyをDBにimportする、このとき、clientの場合はaccount無し
//importしたclientをBitcoin core APIを通じて、walletにimportする
// ImportPublicKeyForWatchWallet csvファイルからpublicアドレスをimportする for WatchOnyWallet
func (w *Wallet) ImportPublicKey(fileName string, accountType account.AccountType, isRescan bool) error {
	//ファイル読み込み
	pubKeys, err := w.addrFileStorager.ImportPubKey(fileName)
	if err != nil {
		return errors.Errorf("key.ImportPubKey() error: %s", err)
	}

	//[]AccountPublicKeyTable //import時にアカウント名をcoreのwalletに登録する
	acnt := string(accountType)
	//if accountType == ctype.AccountTypeClient {
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
	addrInfo, err := w.btc.GetAddressInfo(addr)
	if err != nil {
		w.logger.Error(
			"w.btc.GetAddressInfo()",
			zap.String("address", addr),
			zap.Error(err))
	}
	w.logger.Debug("account is found",
		zap.String("account", addrInfo.GetLabelName()),
		zap.String("address", addr))

	//`watch only wallet` is expected
	if !addrInfo.Iswatchonly {
		w.logger.Error("this address must be watch only wallet")
	}

}
