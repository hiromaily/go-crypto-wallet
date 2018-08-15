package service

import (
	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"sort"
)

// UserPayment ユーザーの支払先アドレスと金額
type UserPayment struct {
	address string
	amount  float64
}

// debug用データ作成
func createDebugUserPayment() []UserPayment {
	//getnewaddress pay1
	//2N33pRYgyuHn6K2xCrrq9dPzuW6ZAvFJfVz

	//getnewaddress pay2
	//2NFd6TEUgSpy8LvttBgVrLB6ZBA5X9BSUSz

	//getnewaddress pay3
	//2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV

	userPayments := []UserPayment{
		{"2N33pRYgyuHn6K2xCrrq9dPzuW6ZAvFJfVz", 0.1},
		{"2NFd6TEUgSpy8LvttBgVrLB6ZBA5X9BSUSz", 0.2},
		{"2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV", 0.3},
	}

	//TODO:Addressが重複する場合、合算したほうがいいかも
	return userPayments
}

// CreateUnsignedTransactionForPayment 支払いのための未署名トランザクションを作成する
func (w *Wallet) CreateUnsignedTransactionForPayment() error {
	//DBから情報を取得、もしくはNatsからのリクエストをトリガーにするか、今はまだ未定
	//とりあえず、テストデータで実装
	userPayments := createDebugUserPayment()

	//2. Listunspent()にてpaymentアカウント用のutxoをすべて取得する
	//listunspent 6  9999999 [\"2N54KrNdyuAkqvvadqSencgpr9XJZnwFYKW\"]
	addr, err := w.Btc.DecodeAddress(w.Btc.PaymentAddress())
	if err != nil {
		return errors.Errorf("DecodeAddress(): error: %v", err)
	}
	unspentList, err := w.Btc.Client().ListUnspentMinMaxAddresses(6, 9999999, []btcutil.Address{addr})

	//3. 送金の金額と近しいutxoでtxを作成するため、ソートしておく
	//unspentList[0].Amount
	sort.Slice(unspentList, func(i, j int) bool {
		//small to big
		return unspentList[i].Amount < unspentList[j].Amount
	})

	grok.Value(userPayments)
	grok.Value(unspentList)

	//for _, userPayment := range userPayments {
	//	//
	//	log.Printf()
	//}

	return nil
}
