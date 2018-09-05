package testdata

import (
	"github.com/hiromaily/go-bitcoin/pkg/model"
)

// CreateTestPaymentRequest
func CreateInitialTestData(m *model.DB){
	//INSERT INTO `payment_request` VALUES
	//(1,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2N33pRYgyuHn6K2xCrrq9dPzuW6ZAvFJfVz',0.1,false,now()),
	//(2,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2NFd6TEUgSpy8LvttBgVrLB6ZBA5X9BSUSz',0.2,false,now()),
	//(3,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV',0.25,false,now()),
	//(4,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV',0.3,false,now()),
	//(5,NULL,'2NFAtuEUzfhEqWgiKYEkSAXUYRutnH75Hkf','yasui','2N7WsiDc4yK7PoUL9saGE5ZGsbRQ8R9NafS',0.4,false,now());

	//1. account_pubkey_clientにアカウント名を仮登録
	//fake.FirstName()


	//2. アドレスにaccount名を登録(bitcoin core経由)

	//3. account_pubkey_clientから情報を取得

	//4. payment_requestテーブルに情報をInsert
	//var paymentRequests []model.PaymentRequest
	//m.InsertPaymentRequest()
}