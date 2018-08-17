package service

import (
	"fmt"
	"log"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/hiromaily/go-bitcoin/pkg/file"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/pkg/errors"
)

// 入金チェックから、utxoを取得し、未署名トランザクションを作成する

// DetectReceivedCoin Wallet内アカウントに入金があれば、そこから、未署名のトランザクションを返す
// 古い未署名のトランザクションは変動するfeeの関係で、stackしていく(再度実行時は差分を抽出する)仕様にはしていない。
// 送信処理後には、unspent()でutxoとして取得できなくなるので、シーケンスで送信まで行うことを想定している
// - 未署名トランザクション作成(本機能)
// - 署名(オフライン)
// - 送信(オンライン)
func (w *Wallet) DetectReceivedCoin() (string, string, error) {
	log.Println("[DetectReceivedCoin]")
	//TODO:このロジックを連続で走らせた場合、現在処理中のものが、タイミングによってはまた取得できてしまうかもしれない??
	// => LockUnspent()

	// LockされたUnspentTransactionを解除する
	if err := w.BTC.UnlockAllUnspentTransaction(); err != nil {
		return "", "", err
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
	//TODO:ListUnspent内で、検索すべきBlock番号まで内部的に保持できてるっぽい
	// Watch only walletであれば、ListUnspentで実現可能
	unspentList, err := w.BTC.Client().ListUnspentMin(6)
	//FIXME: multisigのアドレスはこれで取得できないかも。。。それか、Bitcoin Coreの表示がおかしい。。。
	//FIXME: multisigからhokanに転送したので、multisigには残高がないことが正しい
	if err != nil {
		return "", "", errors.Errorf("ListUnspentMin(): error: %v", err)
	}
	log.Println("[Debug]List Unspent")
	grok.Value(unspentList) //Debug

	if len(unspentList) == 0 {
		return "", "", nil
	}

	var (
		total            btcutil.Amount
		inputs           []btcjson.TransactionInput
		txReceiptDetails []model.TxReceiptDetail
	)

	for _, tx := range unspentList {
		//TODO: spendableは実環境では使えない。
		// 6に満たない場合、まだ未確定であることを意味するはず => これはListUnspent()のパラメータで可能のためコメントアウト
		//https://bitcoin.stackexchange.com/questions/63198/why-outputs-spendable-and-solvable-are-false
		//if tx.Confirmations < w.BTC.ConfirmationBlock() {
		//	continue
		//}

		// Transaction詳細を取得 => これはこのタイミングでは不要
		// (アカウント名によって、はじく、もしくはユーザーのアドレスの取得に必要)
		//tran, err := w.BTC.GetTransactionByTxID(tx.TxID)
		//if err != nil {
		//	//このエラーは起こりえない
		//	log.Printf("[Error] w.BTC.GetTransactionByTxID(): txID:%s, err:%v", tx.TxID, err)
		//	continue
		//}
		//log.Println("[Debug]Transactions")
		//grok.Value(tran)

		//除外するアカウント
		//TODO:本番環境ではこの条件がかわる気がする
		if tx.Account == w.BTC.StoredAccountName() ||
			tx.Account == w.BTC.PaymentAccountName() || tx.Account == "" {
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
		if w.BTC.LockUnspent(tx) != nil {
			continue
		}

		//TODO:以下処理は不要だが、仕様がFIXするまでコメントアウトとして残しておく
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
		// txReceiptDetails
		txReceiptDetails = append(txReceiptDetails, model.TxReceiptDetail{
			ReceiptID:    0,
			InputTxid:    tx.TxID,
			InputVout:    tx.Vout,
			InputAddress: tx.Address,
			InputAccount: tx.Account,
			//InputAmount:        strconv.FormatFloat(tx.Amount, 'f',4, 64),
			InputAmount:        fmt.Sprintf("%f", tx.Amount),
			InputConfirmations: tx.Confirmations,
		})

	}
	log.Printf("[Debug]Total Coin to send:%d(Satoshi) before fee calculated, input length: %d", total, len(inputs))
	if len(inputs) == 0 {
		return "", "", nil
	}

	// 一連の処理を実行
	return w.createRawTransactionAndFee(total, inputs, txReceiptDetails)
}

// createRawTransactionAndFee feeの抽出からtransaction作成、DBへの必要情報保存など、もろもろこちらで行う
func (w *Wallet) createRawTransactionAndFee(total btcutil.Amount, inputs []btcjson.TransactionInput, txReceiptDetails []model.TxReceiptDetail) (string, string, error) {

	// 1.CreateRawTransaction(仮で作成し、この後サイズから手数料を算出する)
	log.Println("[Debug] w.BTC.StoredAddress() :", w.BTC.StoredAddress())
	msgTx, err := w.BTC.CreateRawTransaction(w.BTC.StoredAddress(), total, inputs)
	if err != nil {
		return "", "", errors.Errorf("CreateRawTransaction(): error: %v", err)
	}
	log.Printf("[Debug] CreateRawTransaction: %v\n", msgTx)
	//grok.Value(msgTx)

	// 2.fee算出
	fee, err := w.BTC.GetTransactionFee(msgTx)
	if err != nil {
		return "", "", errors.Errorf("GetTransactionFee(): error: %v", err)
	}
	log.Printf("[Debug]fee: %v", fee) //0.000208 BTC

	// 3.totalを調整し、再度RawTransactionを作成する
	total = total - fee
	if total <= 0 {
		return "", "", errors.Errorf("calculated fee must be wrong: fee:%v, error: %v", fee, err)
	}
	log.Printf("[Debug]Total Coin to send:%d(Satoshi) after fee calculated, input length: %d", total, len(inputs))

	// 4.再度 CreateRawTransaction
	msgTx, err = w.BTC.CreateRawTransaction(w.BTC.StoredAddress(), total, inputs)
	if err != nil {
		return "", "", errors.Errorf("CreateRawTransaction(): error: %v", err)
	}

	// 5.出力用にHexに変換する
	hex, err := w.BTC.ToHex(msgTx)
	if err != nil {
		return "", "", errors.Errorf("w.BTC.ToHex(msgTx): error: %v", err)
	}

	//TODO:以下処理は不要だが、仕様がFIXするまでコメントアウトとしてのこしておく
	//TODO:fundrawtransactionによる手数料の算出は全額送金においては機能しない
	//https://bitcoincore.org/en/doc/0.16.2/rpc/rawtransactions/fundrawtransaction/
	//res, err := w.BTC.FundRawTransaction(hex)
	//if err != nil {
	//	//FIXME:error: -4: Insufficient funds
	//	return "", errors.Errorf("w.BTC.FundRawTransaction(hex): error: %v", err)
	//}
	//[Debug]res.Hex
	//w.BTC.GetRawTransactionByHex(res.Hex)

	// 6. Databaseに必要な情報を保存
	txReceiptID, err := w.insertHexForUnsignedTx(hex, total+fee, fee, w.BTC.StoredAddress(), 1, txReceiptDetails)
	if err != nil {
		return "", "", errors.Errorf("insertHexOnDB(): error: %v", err)
	}

	// 7. GCSにトランザクションファイルを作成
	//TODO:本来、この戻り値をDumpして、GCSに保存、それをDLして、USBに入れてコールドウォレットに移動しなくてはいけない
	//TODO:Debug時はlocalに出力することとする。=> これはフラグで判別したほうがいいかもしれない
	var generatedFileName string
	if txReceiptID != 0 {
		//generatedFileName = file.WriteFileForUnsigned(txReceiptID, "inside/", hex)
		path := file.CreateFilePath(file.ActionReceipt, "unsigned", txReceiptID)
		generatedFileName, err = file.WriteFile(path, hex)
		if err != nil {
			return "", "", errors.Errorf("file.WriteFile(): error: %v", err)
		}
	}

	return hex, generatedFileName, nil
}

func (w *Wallet) storeHexOnGPS(hexTx string) {

}

//TODO:引数の数が多いのはGoにおいてはBad practice...
func (w *Wallet) insertHexForUnsignedTx(hex string, total, fee btcutil.Amount, addr string, txType int, txReceiptDetails []model.TxReceiptDetail) (int64, error) {
	//1.内容が同じだと、生成されるhexもまったく同じ為、同一のhexが合った場合は処理をskipする
	count, err := w.DB.GetTxReceiptByUnsignedHex(hex)
	if err != nil {
		return 0, errors.Errorf("DB.GetTxReceiptByUnsignedHex(): error: %v", err)
	}
	if count != 0 {
		//skip
		return 0, nil
	}

	//2.TxReceiptテーブル
	txReceipt := model.TxReceipt{}
	txReceipt.UnsignedHexTx = hex
	txReceipt.TotalAmount = w.BTC.AmountString(total)
	txReceipt.Fee = w.BTC.AmountString(fee)
	txReceipt.ReceiverAddress = addr
	txReceipt.TxType = 1 //未署名:TODO:Constとして定義しておく

	tx := w.DB.RDB.MustBegin()
	txReceiptID, err := w.DB.InsertTxReceiptForUnsigned(&txReceipt, tx, false)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxReceiptForUnsigned(): error: %v", err)
	}

	//3.TxReceiptDetailテーブル
	for idx, _ := range txReceiptDetails {
		txReceiptDetails[idx].ReceiptID = txReceiptID
	}

	err = w.DB.InsertTxReceiptDetailForUnsigned(txReceiptDetails, tx, true)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxReceiptDetailForUnsigned(): error: %v", err)
	}

	return txReceiptID, nil
}
