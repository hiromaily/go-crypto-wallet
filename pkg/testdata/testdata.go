package testdata

import (
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/api"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/icrowley/fake"
	"github.com/pkg/errors"
)

// CreateInitialTestData
func CreateInitialTestData(m *model.DB, btc *api.Bitcoin) error{
	//1. account_pubkey_clientにアカウント名を仮登録
	accountPubKeyTable, err := m.GetAllAccountPubKeyTable(enum.AccountTypeClient)
	if err != nil{
		return errors.Errorf("DB.GetAllAccountPubKeyTable() error: %s", err)
	}
	for idx , _ := range accountPubKeyTable {
		tm := time.Now()
		accountPubKeyTable[idx].Account = fake.FirstName()
		accountPubKeyTable[idx].UpdatedAt = &tm
	}

	//update
	tx := m.RDB.MustBegin()
	err = m.UpdateAccountOnAccountPubKeyTable(enum.AccountTypeClient, accountPubKeyTable, tx, false)
	if err != nil{
		return errors.Errorf("DB.UpdateAccountOnAccountPubKeyTable() error: %s", err)
	}

	//2. アドレスにaccount名を登録(bitcoin core経由)
	for _ , pubkey := range accountPubKeyTable {
		err = btc.SetAccount(pubkey.WalletAddress, pubkey.Account)
		if err != nil{
			return errors.Errorf("btc.SetAccount() error: %s", err)
		}
		//if btc.Version() >= enum.BTCVer17 {
		//	err = btc.SetLabel(pubkey.WalletAddress, pubkey.Account)
		//	if err != nil{
		//		return errors.Errorf("btc.SetLabel() error: %s", err)
		//	}
		//} else {
		//	err = btc.SetAccount(pubkey.WalletAddress, pubkey.Account)
		//	if err != nil{
		//		return errors.Errorf("btc.SetAccount() error: %s", err)
		//	}
		//}
	}

	//3. payment_requestテーブルに情報をInsert
	paymentRequests := []model.PaymentRequest{
		{
			AddressFrom: accountPubKeyTable[0].WalletAddress,
			AccountFrom: accountPubKeyTable[0].Account,
			AddressTo: accountPubKeyTable[5].WalletAddress,
			Amount: "0.1",
 			IsDone: false,
		},
		{
			AddressFrom: accountPubKeyTable[1].WalletAddress,
			AccountFrom: accountPubKeyTable[1].Account,
			AddressTo: accountPubKeyTable[6].WalletAddress,
			Amount: "0.2",
			IsDone: false,
		},
		{
			AddressFrom: accountPubKeyTable[2].WalletAddress,
			AccountFrom: accountPubKeyTable[2].Account,
			AddressTo: accountPubKeyTable[7].WalletAddress,
			Amount: "0.3",
			IsDone: false,
		},
		{
			AddressFrom: accountPubKeyTable[3].WalletAddress,
			AccountFrom: accountPubKeyTable[3].Account,
			AddressTo: accountPubKeyTable[8].WalletAddress,
			Amount: "0.4",
			IsDone: false,
		},
		{
			AddressFrom: accountPubKeyTable[4].WalletAddress,
			AccountFrom: accountPubKeyTable[4].Account,
			AddressTo: accountPubKeyTable[9].WalletAddress,
			Amount: "0.5",
			IsDone: false,
		},
	}
	//insert
	err = m.InsertPaymentRequest(paymentRequests, tx, true)
	if err != nil{
		return errors.Errorf("btc.InsertPaymentRequest() error: %s", err)
	}

	//INSERT INTO `payment_request` VALUES
	//(1,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2N33pRYgyuHn6K2xCrrq9dPzuW6ZAvFJfVz',0.1,false,now()),
	//(2,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2NFd6TEUgSpy8LvttBgVrLB6ZBA5X9BSUSz',0.2,false,now()),
	//(3,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV',0.25,false,now()),
	//(4,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV',0.3,false,now()),
	//(5,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2N7WsiDc4yK7PoUL9saGE5ZGsbRQ8R9NafS',0.4,false,now());

	return nil
}