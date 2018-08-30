package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/pkg/errors"
	"strings"
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
		//Bitcoin core APIから`importaddress`をcallする
		err := w.BTC.ImportAddressWithLabel(key, account, false)
		if err != nil {
			logger.Errorf("BTC.ImportAddressWithLabel(%s) error: %v", key, err)
			continue
		}

		pubKeyData = append(pubKeyData, model.AccountPublicKeyTable{
			WalletAddress: key,
			Account:       account,
		})
	}

	//DBにInsert
	err = w.DB.InsertAccountPubKeyTable(accountType, pubKeyData, nil, true)
	//logger.Error(err)
	if err != nil {
		return errors.Errorf("csv.DB.InsertAccountPubKeyTable() error: %v", err)
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

	//ファイル読み込み
	pubKeys, err := key.ImportPubKey(fileName)
	if err != nil {
		return errors.Errorf("key.ImportPubKey() error: %v", err)
	}

	//added_pubkey_history_receiptテーブルにInsert
	addedPubkeyHistorys := make([]model.AddedPubkeyHistoryTable, len(pubKeys))
	for i, key := range pubKeys {
		//TODO:とりあえず、1カラムのデータを前提でコーディングしておく
		inner := strings.Split(key, ",")

		addedPubkeyHistorys[i] = model.AddedPubkeyHistoryTable{
			WalletAddress:         inner[0],
			AuthAddress1:          "",
			AuthAddress2:          "",
			WalletMultisigAddress: "",
			RedeemScript:          "",
		}
	}
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
