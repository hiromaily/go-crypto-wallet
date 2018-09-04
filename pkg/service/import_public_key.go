package service

import (
	"github.com/bookerzzz/grok"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/pkg/errors"
	"strings"
	"time"
)

//Watch only wallet
//Cold wallet2

//CSVのpublickeyをDBにimportする、このとき、clientの場合はaccount無し
//importしたclientをBitcoin core APIを通じて、walletにimportする

// ImportPublicKeyForWatchWallet csvファイルからpublicアドレスをimportする for WatchOnyWallet
func (w *Wallet) ImportPublicKeyForWatchWallet(fileName string, accountType enum.AccountType) error {
	//ファイル読み込み
	pubKeys, err := key.ImportPubKey(fileName)
	if err != nil {
		return errors.Errorf("key.ImportPubKey() error: %v", err)
	}

	//[]AccountPublicKeyTable
	account := string(accountType)
	if accountType == enum.AccountTypeClient {
		account = ""
	}

	var pubKeyData []model.AccountPublicKeyTable
	for _, key := range pubKeys {
		inner := strings.Split(key, ",")
		grok.Value(inner)
		var addr string
		if accountType == enum.AccountTypeClient {
			addr = inner[1] //p2sh_segwit_address
		} else {
			addr = inner[3] //wallet_import_format
		}

		//Bitcoin core APIから`importaddress`をcallする
		//TODO:1台のPCで検証しているときなど、すでにimport済の場合はエラーが出る
		err := w.BTC.ImportAddressWithLabel(addr, account, false)
		if err != nil {
			//-4: The wallet already contains the private key for this address or script
			logger.Errorf("BTC.ImportAddressWithLabel(%s) error: %s", addr, err)
			continue
		}

		pubKeyData = append(pubKeyData, model.AccountPublicKeyTable{
			WalletAddress: addr,
			Account:       account,
		})

		//watch only walletとして追加されているかチェックする
		w.checkImportedPublicAddress(addr)
	}

	//DBにInsert
	err = w.DB.InsertAccountPubKeyTable(accountType, pubKeyData, nil, true)
	if err != nil {
		return errors.Errorf("DB.InsertAccountPubKeyTable() error: %v", err)
		//TODO:これが失敗したら、どうやって、登録済みのデータを再度Insertするか？再度実行すればOKのはず
	}

	return nil
}

// ImportPublicKeyForColdWallet2 csvファイルからpublicアドレスをimportする for Cold Wallet2
func (w *Wallet) ImportPublicKeyForColdWallet2(fileName string, accountType enum.AccountType) error {
	//accountチェック
	if accountType != enum.AccountTypeReceipt && accountType != enum.AccountTypePayment {
		logger.Info("AccountType should be AccountTypeReceipt or AccountTypePayment")
		return nil
	}
	//TODO:ImportするファイルのaccountTypeもチェックしたほうがBetter
	//e.g. ./data/pubkey/receipt
	tmp := strings.Split(strings.Split(fileName, "_")[0], "/")
	if tmp[len(tmp)-1] != string(accountType) {
		return errors.Errorf("mismatching between accountType(%s) and file prefix [%s]", accountType, tmp[0])
	}

	//ファイル読み込み(full public key)
	pubKeys, err := key.ImportPubKey(fileName)
	if err != nil {
		return errors.Errorf("key.ImportPubKey() error: %v", err)
	}

	//added_pubkey_history_receiptテーブルにInsert
	addedPubkeyHistorys := make([]model.AddedPubkeyHistoryTable, len(pubKeys))
	for i, key := range pubKeys {
		//TODO:とりあえず、1カラムのデータを前提でコーディングしておく
		inner := strings.Split(key, ",")
		//tmpData := []string{
		//	record.WalletAddress,
		//	record.P2shSegwitAddress,
		//	record.FullPublicKey,
		//	record.WalletMultisigAddress,
		//	record.Account,
		//	strconv.Itoa(int(record.KeyType)),
		//	strconv.Itoa(int(record.Idx)),
		//}

		//TODO:ここでは、FullPublicKeyをセットする必要がある
		addedPubkeyHistorys[i] = model.AddedPubkeyHistoryTable{
			FullPublicKey:         inner[2],
			AuthAddress1:          "",
			AuthAddress2:          "",
			WalletMultisigAddress: "",
			RedeemScript:          "",
		}
	}
	//TODO:Upsertに変えたほうがいいか？Insert済の場合、エラーが出る
	err = w.DB.InsertAddedPubkeyHistoryTable(accountType, addedPubkeyHistorys, nil, true)
	if err != nil {
		return errors.Errorf("DB.InsertAccountKeyClient() error: %s", err)
	}

	// DBにClientAccountのKey情報を登録 (CSVの読み込み)
	//accountKeyClients := make([]model.AccountKeyTable, len(pubKeys))
	//for i, key := range pubKeys {
	//	inner := strings.Split(key, ",")
	//	if len(inner) != 4 {
	//		return errors.New("exported file should be changed as right specification")
	//	}
	//	kt, err := strconv.Atoi(inner[2])
	//	if err != nil {
	//		return errors.New("exported file should be changed as right specification")
	//	}
	//	idx, err := strconv.Atoi(inner[3])
	//	if err != nil {
	//		return errors.New("exported file should be changed as right specification")
	//	}
	//
	//	accountKeyClients[i] = model.AccountKeyTable{
	//		WalletAddress: inner[0],
	//		Account:       inner[1],
	//		KeyType:       uint8(kt),
	//		Idx:           uint32(idx),
	//	}
	//}
	//err = w.DB.InsertAccountKeyTable(accountType, accountKeyClients, nil, true)
	//if err != nil {
	//	return errors.Errorf("DB.InsertAccountKeyClient() error: %s", err)
	//}

	return nil
}

// ImportMultisigAddrForColdWallet1 coldwallet2でexportされたmultisigアドレス情報をimportする for Cold Wallet1
func (w *Wallet) ImportMultisigAddrForColdWallet1(fileName string, accountType enum.AccountType) error {
	if accountType != enum.AccountTypeReceipt && accountType != enum.AccountTypePayment {
		logger.Info("AccountType should be AccountTypeReceipt or AccountTypePayment")
		return nil
	}
	//TODO:ImportするファイルのaccountTypeもチェックしたほうがBetter
	//e.g. ./data/pubkey/receipt
	tmp := strings.Split(strings.Split(fileName, "_")[0], "/")
	if tmp[len(tmp)-1] != string(accountType) {
		return errors.Errorf("mismatching between accountType(%s) and file prefix [%s]", accountType, tmp[0])
	}

	//ファイル読み込み(full public key)
	pubKeys, err := key.ImportPubKey(fileName)
	if err != nil {
		return errors.Errorf("key.ImportPubKey() error: %v", err)
	}

	//added_pubkey_history_receiptテーブルにInsert
	accountKeyTable := make([]model.AccountKeyTable, len(pubKeys))

	tm := time.Now()
	for i, key := range pubKeys {
		//TODO:とりあえず、1カラムのデータを前提でコーディングしておく
		inner := strings.Split(key, ",")
		//csvファイル
		//tmpData := []string{
		//	record.FullPublicKey,
		//	record.AuthAddress1,
		//	record.AuthAddress2,
		//	record.WalletMultisigAddress,
		//	record.RedeemScript,
		//}

		//Upsertをかけるには情報が不足しているので、とりあえず1行ずつUpdateする
		accountKeyTable[i] = model.AccountKeyTable{
			FullPublicKey:         inner[0],
			WalletMultisigAddress: inner[3],
			RedeemScript:          inner[4],
			KeyStatus:             enum.KeyStatusValue[enum.KeyStatusMultiAddressImported],
			UpdatedAt:             &tm,
		}
	}
	//Update
	err = w.DB.UpdateMultisigAddrOnAccountKeyTableByFullPubKey(accountType, accountKeyTable, nil, true)
	if err != nil {
		return errors.Errorf("DB.UpdateMultisigAddrOnAccountKeyTable() error: %s", err)
	}

	return nil
}

//checkImportedPublicAddress watch only walletとして追加されているかチェックする
func (w *Wallet) checkImportedPublicAddress(addr string) {
	if w.BTC.Version() >= 170000 {
		w.checkImportedPublicAddressVer17(addr)
		return
	}

	//1.getaccount address(wallet_address)
	account, err := w.BTC.GetAccount(addr)
	if err != nil {
		logger.Errorf("w.BTC.GetAccount(%s) error: %s", addr, err)
	}
	logger.Debugf("account[%s] is found by wallet_address:%s", account, addr)

	//2.check full_public_key by validateaddress retrieving it
	res, err := w.BTC.ValidateAddress(addr)
	if err != nil {
		logger.Errorf("w.BTC.ValidateAddress(%s) error: %s", addr, err)
	}
	grok.Value(res)
	//watch only walletを想定している
	if !res.IsWatchOnly {
		logger.Errorf("this address must be watch only wallet")
	}

}

//checkImportedPublicAddressVer17 watch only walletとして追加されているかチェックする (for bitcoin version 17)
func (w *Wallet) checkImportedPublicAddressVer17(addr string) {
	logger.Info("checkImportedPublicAddressVer17()")

	//getaddressinfo "address"
	addrInfo, err := w.BTC.GetAddressInfo(addr)
	if err != nil {
		logger.Errorf("w.BTC.GetAddressInfo(%s) error: %s", addr, err)
	}
	logger.Debugf("account[%s] is found by wallet_address:%s", addrInfo.Label, addr)

	//watch only walletを想定している
	if !addrInfo.Iswatchonly {
		logger.Errorf("this address must be watch only wallet")
	}
}
