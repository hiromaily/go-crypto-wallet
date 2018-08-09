package service

import (
	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"log"
)

// DetectReceivedCoin Wallet内アカウントに入金があれば、そこから、未署名のトランザクションを返す
func (w *Wallet) DetectReceivedCoin() (*wire.MsgTx, error) {
	log.Println("[DetectReceivedCoin]")
	//TODO:このロジックを連続で走らせた場合、現在処理中のものが、タイミングによってはまた取得できてしまうかもしれない
	// ので、そこを考慮しないといけない

	//TODO:ここはgoroutineで並列化されたタスク内で、ロックされたtxidを監視し、confirmationが6になったら、
	// 解除するようにしたほうがいいかも。暫定でここに設定
	if err := w.Btc.UnlockAllUnspentTransaction(); err != nil {
		return nil, err
	}

	//1. アカウント一覧からまとめて残高を取得
	//type ListUnspentResult struct {
	//	TxID          string  `json:"txid"`
	//	Vout          uint32  `json:"vout"`
	//	Address       string  `json:"address"`
	//	Account       string  `json:"account"`
	//	ScriptPubKey  string  `json:"scriptPubKey"`
	//	RedeemScript  string  `json:"redeemScript,omitempty"`
	//	Amount        float64 `json:"amount"`
	//	Confirmations int64   `json:"confirmations"`
	//	Spendable     bool    `json:"spendable"`
	//}
	//TODO:とりあえず、ListUnspentを使っているが、全ユーザーにGetUnspentByAddress()を使わないといけないかも
	// Watch only walletであれば、実現できるはず
	list, err := w.Btc.Client().ListUnspent()
	if err != nil {
		return nil, errors.Errorf("ListUnspent(): error: %v", err)
	}
	log.Printf("List Unspent: %v\n", list)
	grok.Value(list) //Debug

	if len(list) == 0 {
		return nil, nil
	}

	var total btcutil.Amount
	var inputs []btcjson.TransactionInput
	//CreateRawTransaction()は外で実行する
	//Loop内ではパラメータを作成するのみ
	for _, tx := range list {
		//FIXME: pendableは実環境では使えない。とりあえず、confirmation数でチェックにしておく
		// 6に満たない場合、まだ未確定であることを意味するはず
		//https://bitcoin.stackexchange.com/questions/63198/why-outputs-spendable-and-solvable-are-false
		if tx.Confirmations < ConfirmationBlockNum {
			//if tx.Spendable == false {
			continue
		}

		// Transaction詳細を取得(必要な情報があるかどうか不明)
		tran, err := w.Btc.GetTransactionByTxID(tx.TxID)
		if err != nil {
			//このエラーは起こりえない
			log.Printf("w.Btc.GetTransactionByTxID(): txID:%s, err:%v", tx.TxID, err)
			continue
		}
		//log.Printf("Transactions: %v\n", tran)
		//grok.Value(tran)

		//除外するアカウント(TODO:これは必要であれば外部定義すべき)
		if tran.Details[0].Account == "hokan" || tran.Details[0].Account == "" {
			continue
		}

		// Amount
		// Satoshiに変換しないといけない
		// 1Satoshi＝0.00000001BTC
		// Float型の計算は微妙なのでint64型に変換して計算する
		amt, err := btcutil.NewAmount(tx.Amount)
		if err != nil {
			//このエラーは起こりえない
			log.Printf("btcutil.NewAmount(): Amount:%d, err:%v", tx.Amount, err)
			continue
		}
		//TODO:全額は送金できないので、このタイミングで手数料を差し引かねばならないが、
		// どこで手数料を算出すべき？total算出後でもいい？
		// => fundrawtransaction()にて算出できるっぽい
		total += amt //合計

		//lockunspentによって、該当トランザクションをロックして再度ListUnspent()で出力されることを防ぐ
		if w.Btc.LockUnspent(tx) != nil {
			continue
		}
		//TODO:ここはまとめてLevelDBにPutしたほうが効率的かもしれない
		//このトランザクションIDはDBに保存が必要 => ここで保存されたIDはconfirmationのチェックに使われる
		err = w.Db.Put("unspent", tx.TxID+string(tx.Vout), nil)
		if err != nil {
			//このタイミングでエラーがおきるのであれば、設計ミス
			log.Printf("Error by w.Db.Put(unspent). This error should not occurred.:, error:%v", err)
			continue
		}

		// inputs
		inputs = append(inputs, btcjson.TransactionInput{
			Txid: tx.TxID,
			Vout: tx.Vout,
		})
	}
	log.Printf("[Debug]Total Coin to send:%d(Satoshi), input length: %d", total, len(inputs))
	if len(inputs) == 0 {
		return nil, nil
	}

	//TODO: このタイミングで手数料を算出して、totalから差し引く？？ => これは不要
	//total = 18500000 //桁を間違えた。。。0.185 BTC
	//total = 195000000

	// CreateRawTransaction
	//FIXME:おつりがあるのであれば、おつりのトランザクションも作らないといけないし、このfuncのインターフェースを見直す必要がある
	msgTx, err := w.Btc.CreateRawTransaction(HokanAddress, total, inputs) //hokanのアドレス
	if err != nil {
		return nil, errors.Errorf("CreateRawTransaction(): error: %v", err)
	}
	log.Printf("[Done]CreateRawTransaction: %v\n", msgTx)
	//grok.Value(msgTx)

	//TODO:fundrawtransactionによって手数料を算出したほうがいい。
	//https://bitcoincore.org/en/doc/0.16.2/rpc/rawtransactions/fundrawtransaction/

	//TODO:本来、この戻り値をDumpして、どっかに保存する必要があるはず、それをUSBに入れてコールドウォレットに移動しなくてはいけない

	return msgTx, nil
}
