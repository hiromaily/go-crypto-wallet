package service

import (
	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	//"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"log"
)

// DetectReceivedCoin Wallet内アカウントに入金があれば、そこから、未署名のトランザクションを返す
//func (w *Wallet) DetectReceivedCoin() (*wire.MsgTx, error) {
func (w *Wallet) DetectReceivedCoin() (string, error) {
	log.Println("[DetectReceivedCoin]")
	//TODO:このロジックを連続で走らせた場合、現在処理中のものが、タイミングによってはまた取得できてしまうかもしれない
	// ので、そこを考慮しないといけない
	// => LockUnspent()

	//TODO:ここはgoroutineで並列化されたタスク内で、ロックされたtxidを監視し、confirmationが6になったら、
	// 解除するようにしたほうがいいかも。暫定でここに設定
	if err := w.Btc.UnlockAllUnspentTransaction(); err != nil {
		return "", err
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
		return "", errors.Errorf("ListUnspent(): error: %v", err)
	}
	log.Printf("[Debug]List Unspent: %v\n", list)
	grok.Value(list) //Debug

	if len(list) == 0 {
		return "", nil
	}

	var total btcutil.Amount
	var inputs []btcjson.TransactionInput
	for _, tx := range list {
		//FIXME: pendableは実環境では使えない。とりあえず、confirmation数でチェックにしておく
		// 6に満たない場合、まだ未確定であることを意味するはず => これはListUnspent()のパラメータで可能
		//https://bitcoin.stackexchange.com/questions/63198/why-outputs-spendable-and-solvable-are-false
		if tx.Confirmations < w.Btc.ConfirmationBlock() {
			//if tx.Spendable == false {
			continue
		}

		// Transaction詳細を取得(必要な情報があるかどうか不明)
		tran, err := w.Btc.GetTransactionByTxID(tx.TxID)
		if err != nil {
			//このエラーは起こりえない
			log.Printf("[Error] w.Btc.GetTransactionByTxID(): txID:%s, err:%v", tx.TxID, err)
			continue
		}
		//log.Printf("[Debug]Transactions: %v\n", tran)
		//grok.Value(tran)

		//除外するアカウント(TODO:これは必要であれば外部定義すべき)
		//=> これは本番環境では不要なはず
		if tran.Details[0].Account == "hokan" || tran.Details[0].Account == "" {
			continue
		}

		// Amount
		amt, err := btcutil.NewAmount(tx.Amount)
		if err != nil {
			//このエラーは起こりえない
			log.Printf("[Error] btcutil.NewAmount(): Amount:%f, err:%v", tx.Amount, err)
			continue
		}
		total += amt //合計

		//lockunspentによって、該当トランザクションをロックして再度ListUnspent()で出力されることを防ぐ
		if w.Btc.LockUnspent(tx) != nil {
			continue
		}

		//TODO:ここはまとめてLevelDBにPutしたほうが効率的かもしれない
		//TODO:ここで保存した情報が本当に必要か？監視すべき対象は何で、何を実現させるかはっきりさせる
		//このトランザクションIDはDBに保存が必要 => ここで保存されたIDはconfirmationのチェックに使われる
		//err = w.Db.Put("unspent", tx.TxID+string(tx.Vout), nil)
		//if err != nil {
		//	//このタイミングでエラーがおきるのであれば、設計ミス
		//	log.Printf("[Error] Error by w.Db.Put(unspent). This error should not occurred.:, error:%v", err)
		//	continue
		//}

		// inputs
		inputs = append(inputs, btcjson.TransactionInput{
			Txid: tx.TxID,
			Vout: tx.Vout,
		})
	}
	log.Printf("[Debug]Total Coin to send:%d(Satoshi) before fee calculated, input length: %d", total, len(inputs))
	if len(inputs) == 0 {
		return "", nil
	}

	// CreateRawTransaction(仮で作成し、この後サイズから手数料を算出する)
	//msgTx, err := w.Btc.CreateRawTransaction(w.Btc.StoreAddr(), total, inputs)
	//[For Only Debug]
	msgTx, err := w.Btc.CreateRawTransaction("muVSWToBoNWusjLCbxcQNBWTmPjioRLpaA", total, inputs)
	if err != nil {
		return "", errors.Errorf("CreateRawTransaction(): error: %v", err)
	}
	log.Printf("[Debug] CreateRawTransaction: %v\n", msgTx)
	//grok.Value(msgTx)

	//fee算出
	fee, err := w.Btc.GetTransactionFee(msgTx)
	if err != nil {
		return "", errors.Errorf("GetTransactionFee(): error: %v", err)
	}
	log.Printf("[Debug]fee: %v", fee) //0.000208 BTC

	//TODO: totalを調整し、再度RawTransactionを作成する
	total = total - fee
	if total <= 0 {
		return "", errors.Errorf("calculated fee must be wrong: fee:%v, error: %v", fee, err)
	}
	log.Printf("[Debug]Total Coin to send:%d(Satoshi) after fee calculated, input length: %d", total, len(inputs))

	msgTx, err = w.Btc.CreateRawTransaction(w.Btc.StoreAddr(), total, inputs)
	if err != nil {
		return "", errors.Errorf("CreateRawTransaction(): error: %v", err)
	}

	//FIXME:ロジックがおかしいかも
	hex, err := w.Btc.ToHex(msgTx)
	if err != nil {
		return "", errors.Errorf("w.Btc.ToHex(msgTx): error: %v", err)
	}

	//TODO:fundrawtransactionによる手数料の算出は全額送金においては機能しない
	//https://bitcoincore.org/en/doc/0.16.2/rpc/rawtransactions/fundrawtransaction/
	//res, err := w.Btc.FundRawTransaction(hex)
	//if err != nil {
	//	//FIXME:error: -4: Insufficient funds
	//	return "", errors.Errorf("w.Btc.FundRawTransaction(hex): error: %v", err)
	//}
	//[Debug]res.Hex
	//w.Btc.GetRawTransactionByHex(res.Hex)

	//TODO:本来、この戻り値をDumpして、GCSに保存、それをDLして、USBに入れてコールドウォレットに移動しなくてはいけない
	return hex, nil
}
