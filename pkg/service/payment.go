package service

import (
	"log"
	"sort"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// ユーザーからの出金依頼事に出金する処理

// UserPayment ユーザーの支払先アドレスと金額
type UserPayment struct {
	senderAddr   string          //送信者のアドレス (履歴を追うためだけに保持)
	receiverAddr string          //受信者のアドレス
	validRecAddr btcutil.Address //受診者のアドレス(変換後)
	amount       float64         //送信金額
	validAmount  btcutil.Amount  //送金金額(変換後)
}

// debug用データ作成
func (w *Wallet) createDebugUserPayment() []UserPayment {
	//getnewaddress pay1
	//2N33pRYgyuHn6K2xCrrq9dPzuW6ZAvFJfVz

	//getnewaddress pay2
	//2NFd6TEUgSpy8LvttBgVrLB6ZBA5X9BSUSz

	//getnewaddress pay3
	//2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV

	//getnewaddress pay4
	//2N7WsiDc4yK7PoUL9saGE5ZGsbRQ8R9NafS

	//5レコードあるが、4種類の送信先
	userPayments := []UserPayment{
		{"", "2N33pRYgyuHn6K2xCrrq9dPzuW6ZAvFJfVz", nil, 0.1, 0},
		{"", "2NFd6TEUgSpy8LvttBgVrLB6ZBA5X9BSUSz", nil, 0.2, 0},
		{"", "2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV", nil, 0.25, 0},
		{"", "2MucBdUqkP5XqNFVTCj35H6WQPC5u2a2BKV", nil, 0.3, 0},
		{"", "2N7WsiDc4yK7PoUL9saGE5ZGsbRQ8R9NafS", nil, 0.4, 0},
	}
	//btcutil.Address型, btcutil.Amount型に変換,
	var err error
	for idx, val := range userPayments {
		//Address TODO:このタイミングでaddressは不要かもしれない
		userPayments[idx].validRecAddr, err = w.Btc.DecodeAddress(val.receiverAddr)
		if err != nil {
			//これは本来事前にチェックされるため、ありえないはず
			log.Printf("[Error] unexpected error converting string to address")
		}
		//grok.Value(userPayments[idx].validRecAddr)

		//Amount
		userPayments[idx].validAmount, err = w.Btc.FloatBitToAmount(val.amount)
		if err != nil {
			//これは本来事前にチェックされるため、ありえないはず
			log.Printf("[Error] unexpected error converting float64 to Amount")
		}
	}

	//Addressが重複する場合、合算したほうがいいかも => transactionのoutputを作成するタイミングで合算可能
	return userPayments
}

// isFoundTxIDAndVout 指定したtxIDとvoutに紐づくinputが既に存在しているか確認する
// TODO:utilityとして有用なので、transaction.goに移動するか
func (w *Wallet) isFoundTxIDAndVout(txID string, vout uint32, inputs []btcjson.TransactionInput) bool {
	for _, val := range inputs {
		if val.Txid == txID && val.Vout == vout {
			return true
		}
	}
	return false
}

// CreateUnsignedTransactionForPayment 支払いのための未署名トランザクションを作成する
func (w *Wallet) CreateUnsignedTransactionForPayment() error {
	//DBから情報を取得、もしくはNatsからのリクエストをトリガーにするか、今はまだ未定
	//とりあえず、テストデータで実装
	//実際には、こちらもamountでソートしておく
	userPayments := w.createDebugUserPayment()

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

	//grok.Value(userPayments)
	grok.Value(unspentList)

	var (
		inputs     []btcjson.TransactionInput
		tmpOutputs = map[string]btcutil.Amount{}
		outputs = map[btcutil.Address]btcutil.Amount{}
	)
	//送信ユーザー毎に処理
	//TODO:lockの処理は後ほど
	for _, userPayment := range userPayments {
		//送金金額をチェック
		//userPayment
		var isAllocated bool
		for idx, utxo := range unspentList {
			if utxo.Amount >= userPayment.amount {
				//利用可能 utxoを発見
				isAllocated = true

				//if !(usedTxID == utxo.TxID && usedTxVout == utxo.Vout){
				if !w.isFoundTxIDAndVout(utxo.TxID, utxo.Vout, inputs) {
					//新規利用のため、inputを作成
					inputs = append(inputs, btcjson.TransactionInput{
						Txid: utxo.TxID,
						Vout: utxo.Vout,
					})
				}
				//Outputを更新
				//log.Println("[Debug]", userPayment.validRecAddr)
				//FIXME:userPayment.validRecAddrはポインタとして保持しているので、値が同じでも異なるものとして認識されてしまう。。。
				//if _, ok := outputs[userPayment.validRecAddr]; ok {
				if _, ok := tmpOutputs[userPayment.receiverAddr]; ok {
					//加算する
					//outputs[userPayment.validRecAddr] += userPayment.validAmount
					tmpOutputs[userPayment.receiverAddr] += userPayment.validAmount
				} else {
					//新規送信先を作成
					//outputs[userPayment.validRecAddr] = userPayment.validAmount
					tmpOutputs[userPayment.receiverAddr] = userPayment.validAmount
				}
				//unspentListも更新する
				unspentList[idx].Amount -= userPayment.amount

				//
				break
			}
		}
		if !isAllocated {
			//送信に必要なutxoが見つからなかった。
			//こんな致命的なエラーは発生しないよう、運用せねばならない。
			log.Printf("[Error] unexpected error: proper utxo could not be found in our accout to send")
		}
	}
	//tmpOutputsをoutputsとして変換する
	for key, val := range tmpOutputs {
		addr, err = w.Btc.DecodeAddress(key)
		if err != nil {
			//これは本来事前にチェックされるため、ありえないはず
			log.Printf("[Error] unexpected error converting string to address")
		}
		outputs[addr] = val
	}
	//TODO:おつりを自分に送信せねばならない
	//inputsにあるものはすべて
	grok.Value(unspentList)


	//TODO:inputsをすべてlockする
	grok.Value(inputs)
	//grok.Value(tmpOutputs)
	grok.Value(outputs)

	return nil
}
