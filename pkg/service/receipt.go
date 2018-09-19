package service

//Watch only wallet

import (
	"fmt"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/hiromaily/go-bitcoin/pkg/api/btc"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/serial"
	"github.com/hiromaily/go-bitcoin/pkg/txfile"
	"github.com/pkg/errors"
)

// 入金チェックから、utxoを取得し、未署名トランザクションを作成する
// 古い未署名のトランザクションは変動するfeeの関係で、stackしていく(再度実行時は差分を抽出する)仕様にはしていない。
// 送信処理後には、unspent()でutxoとして取得できなくなるので、シーケンスで送信まで行うことを想定している
// - 未署名トランザクション作成(本機能)
// - 署名(オフライン)
// - 送信(オンライン)

// DetectReceivedCoin Wallet内アカウントに入金があれば、そこから、未署名のトランザクションを返す
func (w *Wallet) DetectReceivedCoin(adjustmentFee float64) (string, string, error) {
	//TODO:このロジックを連続で走らせた場合、現在処理中のものが、タイミングによってはまた取得できてしまうかもしれない??
	// => LockUnspent()

	// LockされたUnspentTransactionを解除する
	//if err := w.BTC.UnlockAllUnspentTransaction(); err != nil {
	//	return "", "", err
	//}

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

	// Watch only walletであれば、ListUnspentで実現可能
	//unspentList, err := w.BTC.Client().ListUnspentMin(w.BTC.ConfirmationBlock()) //6
	unspentList, err := w.BTC.ListUnspent()
	if err != nil {
		return "", "", errors.Errorf("BTC.Client().ListUnspent(): error: %s", err)
	}
	logger.Debug("List Unspent")
	grok.Value(unspentList) //Debug

	if len(unspentList) == 0 {
		return "", "", nil
	}

	var (
		inputs          []btcjson.TransactionInput
		inputTotal      btcutil.Amount
		txReceiptInputs []model.TxInput
		prevTxs         []btc.PrevTx
		addresses       []string
	)

	for _, tx := range unspentList {

		//除外するアカウント
		//TODO:本番環境ではこの条件がかわる気がする
		//if tx.Account == w.BTC.StoredAccountName() ||
		//	tx.Account == w.BTC.PaymentAccountName() || tx.Account == "" {
		//	continue
		//}
		if tx.Label == string(enum.AccountTypeReceipt) ||
			tx.Label == string(enum.AccountTypePayment) || tx.Label == "" {
			continue
		}

		// Amount
		amt, err := btcutil.NewAmount(tx.Amount)
		if err != nil {
			//このエラーは起こりえない
			logger.Errorf("btcutil.NewAmount(%f): error: %s", tx.Amount, err)
			continue
		}
		inputTotal += amt //合計

		//TODO:Ver17対応
		//lockunspentによって、該当トランザクションをロックして再度ListUnspent()で出力されることを防ぐ
		//if w.BTC.LockUnspent(tx) != nil {
		//	continue
		//}

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
			InputAccount:       tx.Label,
			InputAmount:        fmt.Sprintf("%f", tx.Amount),
			InputConfirmations: tx.Confirmations,
		})

		// prevTxs(walletでの署名でもversion17からは必要になる。。。fuck)
		prevTxs = append(prevTxs, btc.PrevTx{
			Txid:         tx.TxID,
			Vout:         tx.Vout,
			ScriptPubKey: tx.ScriptPubKey,
			RedeemScript: "", //TODO:redeemScriptはどうやって算出する??おそらく、txidから詳細を取得する必要がある?gettransactionはだめだった。。。
			Amount:       tx.Amount,
		})
		//tx.Address
		addresses = append(addresses, tx.Address)
	}
	logger.Debugf("total coin to send:%d(Satoshi) before fee calculated, input length: %d", inputTotal, len(inputs))
	if len(inputs) == 0 {
		return "", "", nil
	}

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:   addresses,
		PrevTxs: prevTxs,
	}

	// 一連の処理を実行
	hex, fileName, err := w.createRawTransactionAndFee(adjustmentFee, inputs, inputTotal, txReceiptInputs, &addrsPrevs)

	//TODO:Ver17対応
	// LockされたUnspentTransactionを解除する
	//if err := w.BTC.UnlockAllUnspentTransaction(); err != nil {
	//	return "", "", err
	//}

	return hex, fileName, err
}

// createRawTransactionAndFee feeの抽出からtransaction作成、DBへの必要情報保存など、もろもろこちらで行う
func (w *Wallet) createRawTransactionAndFee(adjustmentFee float64, inputs []btcjson.TransactionInput,
	inputTotal btcutil.Amount, txReceiptInputs []model.TxInput, addrsPrevs *btc.AddrsPrevTxs) (string, string, error) {

	var outputTotal btcutil.Amount

	//TODO:w.BTC.StoredAddress() この部分をDatabaseから取得しないといけない
	//TODO:送金時に、フラグ(is_allocated)をONにすることとする
	pubkeyTable, err := w.DB.GetOneUnAllocatedAccountPubKeyTable(enum.AccountTypeReceipt)
	if err != nil {
		return "", "", errors.Errorf("DB.GetOneUnAllocatedAccountPubKeyTable(): error: %s", err)
	}
	storedAddr := pubkeyTable.WalletAddress //change from w.BTC.StoredAddress()
	storedAccount := pubkeyTable.Account    //change from w.BTC.StoredAccountName()

	// 1.CreateRawTransaction(仮で作成し、この後サイズから手数料を算出する)
	msgTx, err := w.BTC.CreateRawTransaction(storedAddr, inputTotal, inputs)
	if err != nil {
		return "", "", errors.Errorf("BTC.CreateRawTransaction(): error: %s", err)
	}

	// 2.fee算出
	fee, err := w.BTC.GetFee(msgTx, adjustmentFee)

	// 3.手数料のために、totalを調整し、再度RawTransactionを作成する
	//このパートは、出金とロジックが異なる
	outputTotal = inputTotal - fee
	if outputTotal <= 0 {
		return "", "", errors.Errorf("calculated fee must be wrong: fee:%v, error: %s", fee, err)
	}
	logger.Debugf("Total Coin to send:%d(Satoshi) after fee calculated, input length: %d", outputTotal, len(inputs))

	// 4.outputs作成
	txReceiptOutputs := []model.TxOutput{
		{
			ReceiptID:     0,
			OutputAddress: storedAddr,
			OutputAccount: storedAccount,
			OutputAmount:  w.BTC.AmountString(outputTotal),
			IsChange:      false,
		},
	}

	// 5.再度 CreateRawTransaction
	msgTx, err = w.BTC.CreateRawTransaction(storedAddr, outputTotal, inputs)
	if err != nil {
		return "", "", errors.Errorf("BTC.CreateRawTransaction(): error: %s", err)
	}

	// 6.出力用にHexに変換する
	hex, err := w.BTC.ToHex(msgTx)
	if err != nil {
		return "", "", errors.Errorf("BTC.ToHex(msgTx): error: %s", err)
	}

	// 7. Databaseに必要な情報を保存
	txReceiptID, err := w.insertTxTableForUnsigned(enum.ActionTypeReceipt, hex, inputTotal, outputTotal, fee, enum.TxTypeValue[enum.TxTypeUnsigned], txReceiptInputs, txReceiptOutputs, nil)
	if err != nil {
		return "", "", errors.Errorf("insertTxTableForUnsigned(): error: %s", err)
	}

	// 8. serialize previous txs for multisig signature
	encodedAddrsPrevs, err := serial.EncodeToString(*addrsPrevs)
	if err != nil {
		return "", "", errors.Errorf("serial.EncodeToString(): error: %s", err)
	}
	logger.Debugf("encodedAddrsPrevs: %s", encodedAddrsPrevs)

	// 9. GCSにトランザクションファイルを作成
	//TODO:本来、この戻り値をDumpして、GCSに保存、それをDLして、USBに入れてコールドウォレットに移動しなくてはいけない
	//TODO:Debug時はlocalに出力することとする。=> これはフラグで判別したほうがいいかもしれない/Interface型にして対応してもいいかも
	var generatedFileName string
	if txReceiptID != 0 {
		generatedFileName, err = w.storeHex(hex, encodedAddrsPrevs, txReceiptID, enum.ActionTypeReceipt)
		if err != nil {
			return "", "", errors.Errorf("wallet.storeHex(): error: %s", err)
		}
	}

	// 10. 入金準備に入ったことをユーザーに通知
	// TODO:NatsのPublisherとして通知すればいいか？

	return hex, generatedFileName, nil
}

//TODO:引数の数が多いのはGoにおいてはBad practice...
//[共通(receipt/payment)]
func (w *Wallet) insertTxTableForUnsigned(actionType enum.ActionType, hex string, inputTotal, outputTotal, fee btcutil.Amount, txType uint8,
	txInputs []model.TxInput, txOutputs []model.TxOutput, paymentRequestIds []int64) (int64, error) {

	//1.内容が同じだと、生成されるhexもまったく同じ為、同一のhexが合った場合は処理をskipする
	//count, err := w.DB.GetTxReceiptCountByUnsignedHex(hex)
	count, err := w.DB.GetTxCountByUnsignedHex(actionType, hex)
	if err != nil {
		return 0, errors.Errorf("DB.GetTxCountByUnsignedHex(): error: %s", err)
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
	txReceiptID, err := w.DB.InsertTxForUnsigned(actionType, &txReceipt, tx, false)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxForUnsigned(): error: %s", err)
	}

	//3.TxReceiptInputテーブル
	//ReceiptIDの更新
	for idx := range txInputs {
		txInputs[idx].ReceiptID = txReceiptID
	}
	err = w.DB.InsertTxInputForUnsigned(actionType, txInputs, tx, false)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxInputForUnsigned(): error: %s", err)
	}

	//4.TxReceiptOutputテーブル
	//ReceiptIDの更新
	for idx := range txOutputs {
		txOutputs[idx].ReceiptID = txReceiptID
	}

	//commit flag
	isCommit := true
	if actionType == enum.ActionTypePayment {
		isCommit = false
	}

	err = w.DB.InsertTxOutputForUnsigned(actionType, txOutputs, tx, isCommit)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxOutputForUnsigned(): error: %s", err)
	}

	//5. payment_requestのpayment_idを更新する paymentRequestIds
	if actionType == enum.ActionTypePayment {
		//txReceiptID
		_, err = w.DB.UpdatePaymentIDOnPaymentRequest(txReceiptID, paymentRequestIds, tx, true)
		if err != nil {
			return 0, errors.Errorf("DB.UpdatePaymentIDOnPaymentRequest(): error: %s", err)
		}
	}

	return txReceiptID, nil
}

// storeHex　hex情報を保存し、ファイル名を返す
// [共通(receipt/payment)]
func (w *Wallet) storeHex(hex, encodedAddrsPrevs string, id int64, actionType enum.ActionType) (string, error) {
	var (
		generatedFileName string
		err               error
	)

	savedata := hex
	if encodedAddrsPrevs != "" {
		savedata = fmt.Sprintf("%s,%s", savedata, encodedAddrsPrevs)
	}

	//To File
	if w.Env == enum.EnvDev {
		path := txfile.CreateFilePath(actionType, enum.TxTypeUnsigned, id, true)
		generatedFileName, err = txfile.WriteFile(path, savedata)
		if err != nil {
			return "", errors.Errorf("txfile.WriteFile(): error: %s", err)
		}
	}

	//[WIP] GCS
	path := txfile.CreateFilePath(actionType, enum.TxTypeUnsigned, id, false)

	//GCS上に、Clientを作成(セッションの関係で都度作成する)
	_, err = w.GCS[actionType].WriteOnce(path, savedata)
	if err != nil {
		return "", errors.Errorf("storage.WriteOnce(): error: %s", err)
	}

	return generatedFileName, nil
}
