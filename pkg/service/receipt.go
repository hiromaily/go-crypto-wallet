package service

import (
	"fmt"
	"log"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
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
func (w *Wallet) DetectReceivedCoin(adjustmentFee float64) (string, string, error) {
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
	if err != nil {
		return "", "", errors.Errorf("ListUnspentMin(): error: %v", err)
	}
	log.Println("[Debug]List Unspent")
	grok.Value(unspentList) //Debug

	if len(unspentList) == 0 {
		return "", "", nil
	}

	var (
		inputs          []btcjson.TransactionInput
		inputTotal      btcutil.Amount
		txReceiptInputs []model.TxInput
	)

	for _, tx := range unspentList {
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
		inputTotal += amt //合計

		//lockunspentによって、該当トランザクションをロックして再度ListUnspent()で出力されることを防ぐ
		if w.BTC.LockUnspent(tx) != nil {
			continue
		}

		// inputs
		inputs = append(inputs, btcjson.TransactionInput{
			Txid: tx.TxID,
			Vout: tx.Vout,
		})
		// txReceiptInputs
		txReceiptInputs = append(txReceiptInputs, model.TxInput{
			ReceiptID:          0,
			InputTxid:          tx.TxID,
			InputVout:          tx.Vout,
			InputAddress:       tx.Address,
			InputAccount:       tx.Account,
			InputAmount:        fmt.Sprintf("%f", tx.Amount),
			InputConfirmations: tx.Confirmations,
		})

	}
	log.Printf("[Debug]Total Coin to send:%d(Satoshi) before fee calculated, input length: %d", inputTotal, len(inputs))
	if len(inputs) == 0 {
		return "", "", nil
	}

	// 一連の処理を実行
	return w.createRawTransactionAndFee(adjustmentFee, inputs, inputTotal, txReceiptInputs)
}

// createRawTransactionAndFee feeの抽出からtransaction作成、DBへの必要情報保存など、もろもろこちらで行う
func (w *Wallet) createRawTransactionAndFee(adjustmentFee float64, inputs []btcjson.TransactionInput, inputTotal btcutil.Amount, txReceiptInputs []model.TxInput) (string, string, error) {
	var outputTotal btcutil.Amount

	// 1.CreateRawTransaction(仮で作成し、この後サイズから手数料を算出する)
	log.Println("[Debug] w.BTC.StoredAddress() :", w.BTC.StoredAddress())
	msgTx, err := w.BTC.CreateRawTransaction(w.BTC.StoredAddress(), inputTotal, inputs)
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
	log.Printf("[Debug]fee: %v, %f", fee, adjustmentFee) //0.000208 BTC

	//FIXME:処理が受理されないトランザクションを作るために、意図的に1Satothiのfeeでトランザクションを作る
	//DEBUG: Relayfeeにより、最低でも1000Satoshi必要
	//fee = 1000



	// 3.手数料のために、totalを調整し、再度RawTransactionを作成する
	//このパートのみが、出金とロジックが異なる
	outputTotal = inputTotal - fee
	if outputTotal <= 0 {
		return "", "", errors.Errorf("calculated fee must be wrong: fee:%v, error: %v", fee, err)
	}
	log.Printf("[Debug]Total Coin to send:%d(Satoshi) after fee calculated, input length: %d", outputTotal, len(inputs))

	// 4.outputs作成
	txReceiptOutputs := []model.TxOutput{
		{
			ReceiptID:     0,
			OutputAddress: w.BTC.StoredAddress(),
			OutputAccount: w.BTC.StoredAccountName(),
			OutputAmount:  w.BTC.AmountString(outputTotal),
			IsChange:      false,
		},
	}

	// 5.再度 CreateRawTransaction
	//TODO:同一I/Fのために、outputsを作成して、CreateRawTransactionWithOutput()を呼び出したほうがいいかもしれない。
	msgTx, err = w.BTC.CreateRawTransaction(w.BTC.StoredAddress(), outputTotal, inputs)
	if err != nil {
		return "", "", errors.Errorf("CreateRawTransaction(): error: %v", err)
	}

	// 6.出力用にHexに変換する
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

	// 7. Databaseに必要な情報を保存
	//  txType //1.未署名
	txReceiptID, err := w.insertHexForUnsignedTxOnReceipt(hex, inputTotal, outputTotal, fee, enum.TxTypeValue[enum.TxTypeUnsigned], txReceiptInputs, txReceiptOutputs)
	if err != nil {
		return "", "", errors.Errorf("insertHexOnDB(): error: %v", err)
	}

	// 8. GCSにトランザクションファイルを作成
	//TODO:本来、この戻り値をDumpして、GCSに保存、それをDLして、USBに入れてコールドウォレットに移動しなくてはいけない
	//TODO:Debug時はlocalに出力することとする。=> これはフラグで判別したほうがいいかもしれない/Interface型にして対応してもいいかも
	var generatedFileName string
	if txReceiptID != 0 {
		path := file.CreateFilePath(enum.ActionTypeReceipt, enum.TxTypeUnsigned, txReceiptID)
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
func (w *Wallet) insertHexForUnsignedTxOnReceipt(hex string, inputTotal, outputTotal, fee btcutil.Amount, txType uint8,
	txReceiptInputs []model.TxInput, txReceiptOutputs []model.TxOutput) (int64, error) {

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
	txReceipt := model.TxTable{}
	txReceipt.UnsignedHexTx = hex
	txReceipt.TotalInputAmount = w.BTC.AmountString(inputTotal)
	txReceipt.TotalOutputAmount = w.BTC.AmountString(outputTotal)
	txReceipt.Fee = w.BTC.AmountString(fee)
	txReceipt.TxType = txType

	tx := w.DB.RDB.MustBegin()
	txReceiptID, err := w.DB.InsertTxReceiptForUnsigned(&txReceipt, tx, false)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxReceiptForUnsigned(): error: %v", err)
	}

	//3.TxReceiptInputテーブル
	//ReceiptIDの更新
	for idx := range txReceiptInputs {
		txReceiptInputs[idx].ReceiptID = txReceiptID
	}

	err = w.DB.InsertTxReceiptInputForUnsigned(txReceiptInputs, tx, false)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxReceiptInputForUnsigned(): error: %v", err)
	}

	//4.TxReceiptOutputテーブル
	//ReceiptIDの更新
	for idx := range txReceiptOutputs {
		txReceiptOutputs[idx].ReceiptID = txReceiptID
	}

	err = w.DB.InsertTxReceiptOutputForUnsigned(txReceiptOutputs, tx, true)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxReceiptOutputForUnsigned(): error: %v", err)
	}

	return txReceiptID, nil
}
