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
		userPayments[idx].validRecAddr, err = w.BTC.DecodeAddress(val.receiverAddr)
		if err != nil {
			//これは本来事前にチェックされるため、ありえないはず
			log.Printf("[Error] unexpected error converting string to address")
		}
		//grok.Value(userPayments[idx].validRecAddr)

		//Amount
		userPayments[idx].validAmount, err = w.BTC.FloatBitToAmount(val.amount)
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
func (w *Wallet) CreateUnsignedTransactionForPayment() (string, error) {
	//DBから情報を取得、もしくはNatsからのリクエストをトリガーにするか、今はまだ未定
	//とりあえず、テストデータで実装
	//実際には、こちらもamountでソートしておく=> ソートは不要
	userPayments := w.createDebugUserPayment()

	//2.送信合計金額を算出
	var userTotal btcutil.Amount
	for _, val := range userPayments {
		userTotal += val.validAmount
	}

	//3.保管用アドレスの残高を確認し、金額が不足しているのであればエラー
	//getbalance hokan 6
	balance, err := w.BTC.GetBalanceByAccountAndMinConf(w.BTC.PaymentAccountName(), w.BTC.ConfirmationBlock())
	if err != nil {
		return "", err
	}
	if balance <= userTotal {
		return "", errors.New("Stored account balance is insufficient")
	}

	//4. Listunspent()にてpaymentアカウント用のutxoをすべて取得する
	//listunspent 6  9999999 [\"2N54KrNdyuAkqvvadqSencgpr9XJZnwFYKW\"]
	addr, err := w.BTC.DecodeAddress(w.BTC.PaymentAddress())
	if err != nil {
		return "", errors.Errorf("DecodeAddress(): error: %v", err)
	}
	unspentList, err := w.BTC.Client().ListUnspentMinMaxAddresses(6, 9999999, []btcutil.Address{addr})
	if err != nil {
		//致命的なエラー
		return "", err
	}

	//5. 送金の金額と近しいutxoでtxを作成するため、ソートしておく => 小さなutxoから利用していくのに便利だが、MUSTではない
	sort.Slice(unspentList, func(i, j int) bool {
		//small to big
		return unspentList[i].Amount < unspentList[j].Amount
	})

	//grok.Value(userPayments)
	//grok.Value(unspentList)

	var (
		inputs     []btcjson.TransactionInput
		tmpOutputs = map[string]btcutil.Amount{}
		outputs    = map[btcutil.Address]btcutil.Amount{}
	)

	//6.合計金額を超えるまで、listunspentからinputsを作成する
	var inputTotal btcutil.Amount
	var isDone bool
	for _, utxo := range unspentList {
		//utxo.Amount
		//utxo.Vout
		//utxo.TxID
		amt, err := btcutil.NewAmount(utxo.Amount)
		if err != nil {
			log.Println("[Error] unexpected error:", err)
			continue
		}
		inputTotal += amt
		//新規利用のため、inputを作成
		inputs = append(inputs, btcjson.TransactionInput{
			Txid: utxo.TxID,
			Vout: utxo.Vout,
		})
		if inputTotal > userTotal {
			isDone = true
			break
		}
	}
	if !isDone {
		//これは金額が足りなかったことを意味するので致命的
		return "", errors.New("[fatal error] bitcoin is insufficient in trading account of ours")
	}
	//outputは別途こちらで作成
	for _, userPayment := range userPayments {
		if _, ok := tmpOutputs[userPayment.receiverAddr]; ok {
			//加算する
			tmpOutputs[userPayment.receiverAddr] += userPayment.validAmount
		} else {
			//新規送信先を作成
			tmpOutputs[userPayment.receiverAddr] = userPayment.validAmount
		}
	}

	//差分でお釣り用のoutputを作成する
	change := inputTotal - userTotal
	tmpOutputs[w.BTC.PaymentAddress()] = change

	//tmpOutputsをoutputsとして変換する
	for key, val := range tmpOutputs {
		addr, err = w.BTC.DecodeAddress(key)
		if err != nil {
			//これは本来事前にチェックされるため、ありえないはず
			log.Printf("[Error] unexpected error converting string to address")
			continue
		}
		outputs[addr] = val
	}

	//TODO:inputsをすべてlockする
	grok.Value(inputs)
	grok.Value(tmpOutputs)
	grok.Value(outputs)

	//rawtransactionを
	// 一連の処理を実行
	return w.createRawTransactionForPayment(inputs, outputs)
}

//TODO:receipt側のロジックとほぼ同じ為、まとめたほうがいい(手数料のロジックのみ異なる)
func (w *Wallet) createRawTransactionForPayment(inputs []btcjson.TransactionInput, outputs map[btcutil.Address]btcutil.Amount) (string, error) {
	// 1.CreateRawTransactionWithOutput(仮で作成し、この後サイズから手数料を算出する)
	msgTx, err := w.BTC.CreateRawTransactionWithOutput(inputs, outputs)
	if err != nil {
		return "", errors.Errorf("CreateRawTransactionWithOutput(): error: %v", err)
	}
	log.Printf("[Debug] CreateRawTransactionWithOutput: %v\n", msgTx)
	//grok.Value(msgTx)

	// 2.fee算出
	fee, err := w.BTC.GetTransactionFee(msgTx)
	if err != nil {
		return "", errors.Errorf("GetTransactionFee(): error: %v", err)
	}
	log.Printf("[Debug]fee: %v", fee) //0.001183 BTC

	// 3.TODO:お釣り用のoutputのトランザクションから、手数料を差し引かねばならい
	// FIXME: これが足りない場合がめんどくさい。。。これをどう回避すべきか
	for addr := range outputs {
		if addr.String() == w.BTC.PaymentAddress() {
			outputs[addr] -= fee
		}
	}
	grok.Value(outputs)

	//debug
	//return "", nil

	// 4.再度 CreateRawTransaction
	msgTx, err = w.BTC.CreateRawTransactionWithOutput(inputs, outputs)
	if err != nil {
		return "", errors.Errorf("CreateRawTransaction(): error: %v", err)
	}

	// 5.出力用にHexに変換する
	hex, err := w.BTC.ToHex(msgTx)
	if err != nil {
		return "", errors.Errorf("w.BTC.ToHex(msgTx): error: %v", err)
	}

	// 6. Databaseに必要な情報を保存
	//txReceiptID, err := w.insertHexForUnsignedTx(hex, total+fee, fee, w.BTC.StoredAddress(), 1, txReceiptDetails)
	//if err != nil {
	//	return "", "", errors.Errorf("insertHexOnDB(): error: %v", err)
	//}

	// 6. GCSにトランザクションファイルを作成
	//TODO:本来、この戻り値をDumpして、GCSに保存、それをDLして、USBに入れてコールドウォレットに移動しなくてはいけない
	//TODO:Debug時はlocalに出力することとする
	//FIXME:該当するtransactionのIDを設定せねばならない
	//var generatedFileName string
	//if txReceiptID != 0 {
	//	generatedFileName = file.WriteFileForUnsigned(txReceiptID, hex)
	//}

	//return hex, generatedFileName, nil
	return hex, nil
}

// CreateUnsignedTransactionForPaymentOld 支払いのための未署名トランザクションを作成する
// TODO:別ロジックとして考えるため、一旦退避
//func (w *Wallet) CreateUnsignedTransactionForPaymentOld() error {
//	//DBから情報を取得、もしくはNatsからのリクエストをトリガーにするか、今はまだ未定
//	//とりあえず、テストデータで実装
//	//実際には、こちらもamountでソートしておく
//	userPayments := w.createDebugUserPayment()
//
//	//2. Listunspent()にてpaymentアカウント用のutxoをすべて取得する
//	//listunspent 6  9999999 [\"2N54KrNdyuAkqvvadqSencgpr9XJZnwFYKW\"]
//	addr, err := w.BTC.DecodeAddress(w.BTC.PaymentAddress())
//	if err != nil {
//		return errors.Errorf("DecodeAddress(): error: %v", err)
//	}
//	unspentList, err := w.BTC.Client().ListUnspentMinMaxAddresses(6, 9999999, []btcutil.Address{addr})
//	if err != nil {
//		//致命的なエラー
//		return err
//	}
//
//	//3. 送金の金額と近しいutxoでtxを作成するため、ソートしておく
//	//unspentList[0].Amount
//	sort.Slice(unspentList, func(i, j int) bool {
//		//small to big
//		return unspentList[i].Amount < unspentList[j].Amount
//	})
//
//	//grok.Value(userPayments)
//	grok.Value(unspentList)
//
//	var (
//		inputs     []btcjson.TransactionInput
//		tmpOutputs = map[string]btcutil.Amount{}
//		outputs    = map[btcutil.Address]btcutil.Amount{}
//	)
//	//送信ユーザー毎に処理
//	for _, userPayment := range userPayments {
//		//送金金額をチェック
//		//userPayment
//		var isAllocated bool
//		for idx, utxo := range unspentList {
//			if utxo.Amount >= userPayment.amount {
//				//利用可能 utxoを発見
//				isAllocated = true
//
//				//if !(usedTxID == utxo.TxID && usedTxVout == utxo.Vout){
//				if !w.isFoundTxIDAndVout(utxo.TxID, utxo.Vout, inputs) {
//					//新規利用のため、inputを作成
//					inputs = append(inputs, btcjson.TransactionInput{
//						Txid: utxo.TxID,
//						Vout: utxo.Vout,
//					})
//				}
//				//Outputを更新
//				//log.Println("[Debug]", userPayment.validRecAddr)
//				//FIXME:userPayment.validRecAddrはポインタとして保持しているので、値が同じでも異なるものとして認識されてしまう。。。
//				//if _, ok := outputs[userPayment.validRecAddr]; ok {
//				if _, ok := tmpOutputs[userPayment.receiverAddr]; ok {
//					//加算する
//					//outputs[userPayment.validRecAddr] += userPayment.validAmount
//					tmpOutputs[userPayment.receiverAddr] += userPayment.validAmount
//				} else {
//					//新規送信先を作成
//					//outputs[userPayment.validRecAddr] = userPayment.validAmount
//					tmpOutputs[userPayment.receiverAddr] = userPayment.validAmount
//				}
//				//unspentListも更新する
//				unspentList[idx].Amount -= userPayment.amount
//
//				//
//				break
//			}
//		}
//		if !isAllocated {
//			//送信に必要なutxoが見つからなかった。
//			//こんな致命的なエラーは発生しないよう、運用せねばならない。
//			log.Printf("[Error] unexpected error: proper utxo could not be found in our accout to send")
//		}
//	}
//	//tmpOutputsをoutputsとして変換する
//	for key, val := range tmpOutputs {
//		addr, err = w.BTC.DecodeAddress(key)
//		if err != nil {
//			//これは本来事前にチェックされるため、ありえないはず
//			log.Printf("[Error] unexpected error converting string to address")
//		}
//		outputs[addr] = val
//	}
//	//TODO:おつりを自分に送信せねばならない
//	//inputsにあるものはすべて
//	grok.Value(unspentList)
//
//	//TODO:inputsをすべてlockする
//	grok.Value(inputs)
//	//grok.Value(tmpOutputs)
//	grok.Value(outputs)
//
//	return nil
//}
