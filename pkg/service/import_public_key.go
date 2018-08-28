package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/csv"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/pkg/errors"
)

//Watch only wallet

//CSVのpublickeyをDBにimportする、このとき、clientの場合はaccount無し
//importしたclientをBitcoin core APIを通じて、walletにimportする

// ImportPublicKey csvファイルからpublicアドレスをimportする
func (w *Wallet) ImportPublicKey(fileName string, accountType enum.AccountType) error {
	//ファイル読み込み
	pubKeys, err := csv.ImportPubKey(fileName)
	if err != nil {
		return errors.Errorf("csv.ImportPubKey() error: %v", err)
	}

	//[]AccountPublicKeyTable
	account := string(accountType)
	if accountType == enum.AccountTypeClient {
		account = ""
	}

	var pubKeyData []model.AccountPublicKeyTable
	for _, key := range pubKeys {
		//Bitcoin core APIから`importaddress`をcallする
		//err := w.BTC.ImportAddressWithLabel(key, account, false)
		//if err != nil {
		//	logger.Errorf("BTC.ImportAddressWithLabel(%s) error: %v", key, err)
		//	continue
		//}

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
		//TODO:これが失敗したら、どうやって、登録済みのデータを再度Insertするか？
	}

	return nil
}
