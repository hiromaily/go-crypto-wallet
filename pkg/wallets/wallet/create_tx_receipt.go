package wallet

import (
	"fmt"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
	"go.uber.org/zap"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/serial"
	"github.com/hiromaily/go-bitcoin/pkg/txfile"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api/btc"
)

// 入金チェックから、utxoを取得し、未署名トランザクションを作成する
// 古い未署名のトランザクションは変動するfeeの関係で、stackしていく(再度実行時は差分を抽出する)仕様にはしていない。
// 送信処理後には、unspent()でutxoとして取得できなくなるので、シーケンスで送信まで行うことを想定している
// - 未署名トランザクション作成(本機能)
// - 署名(オフライン)
// - 送信(オンライン)

// DetectReceivedCoin Wallet内アカウントに入金があれば、そこから、未署名のトランザクションを返す
func (w *Wallet) DetectReceivedCoin(adjustmentFee float64) (string, string, error) {
	//TODO:remove it

	// LockされたUnspentTransactionを解除する
	//if err := w.BTC.UnlockAllUnspentTransaction(); err != nil {
	//	return "", "", err
	//}

	// Watch only walletであれば、ListUnspentで実現可能
	//unspentList, err := w.BTC.ListUnspent()
	unspentList, _, err := w.btc.ListUnspentByAccount(account.AccountTypeClient)

	if err != nil {
		return "", "", errors.Errorf("BTC.Client().ListUnspent(): error: %s", err)
	}
	w.logger.Debug("List Unspent")
	grok.Value(unspentList) //Debug

	if len(unspentList) == 0 {
		w.logger.Info("no listunspent")
		return "", "", nil
	}

	var (
		inputs          []btcjson.TransactionInput
		inputTotal      btcutil.Amount
		txReceiptInputs []walletrepo.TxInput
		prevTxs         []btc.PrevTx
		addresses       []string
	)

	for _, tx := range unspentList {

		//除外するアカウント
		//TODO:本番環境ではこの条件がかわる気がする=>はじめからclientアカウントの情報を取得しておく
		//if tx.Label == string(enum.AccountTypeReceipt) ||
		//	tx.Label == string(enum.AccountTypePayment) || tx.Label == "" {
		//	continue
		//}

		// Amount
		amt, err := btcutil.NewAmount(tx.Amount)
		if err != nil {
			//このエラーは起こりえない
			w.logger.Error(
				"btcutil.NewAmount()",
				zap.Float64("tx amount", tx.Amount),
				zap.Error(err))
			continue
		}
		inputTotal += amt //合計

		//TODO:Ver17対応が必要
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
		txReceiptInputs = append(txReceiptInputs, walletrepo.TxInput{
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
			RedeemScript: "", //multisigではない場合は、不要
			Amount:       tx.Amount,
		})

		//tx.Address
		addresses = append(addresses, tx.Address)
	}
	w.logger.Debug(
		"total coin to send (Satoshi) before fee calculated",
		zap.Any("amount", inputTotal),
		zap.Int("len(inputs)", len(inputs)))
	if len(inputs) == 0 {
		return "", "", nil
	}

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:         addresses,
		PrevTxs:       prevTxs,
		SenderAccount: account.AccountTypeClient,
	}

	// 一連の処理を実行
	hex, fileName, err := w.createRawTransactionAndFee(enum.ActionTypeReceipt, account.AccountTypeReceipt, adjustmentFee,
		inputs, inputTotal, txReceiptInputs, &addrsPrevs)

	//TODO:Ver17対応が必要
	// LockされたUnspentTransactionを解除する
	//if err := w.BTC.UnlockAllUnspentTransaction(); err != nil {
	//	return "", "", err
	//}

	return hex, fileName, err
}

// createRawTransactionAndFee feeの抽出からtransaction作成、DBへの必要情報保存など、もろもろこちらで行う
// receipt/transfer共通
func (w *Wallet) createRawTransactionAndFee(actionType enum.ActionType, accountType account.AccountType,
	adjustmentFee float64, inputs []btcjson.TransactionInput, inputTotal btcutil.Amount,
	txReceiptInputs []walletrepo.TxInput, addrsPrevs *btc.AddrsPrevTxs) (string, string, error) {

	var outputTotal btcutil.Amount

	//TODO:送金時に、フラグ(is_allocated)をONにすることとする(ここで設定するAccountは受信者)
	pubkeyTable, err := w.storager.GetOneUnAllocatedAccountPubKeyTable(accountType)
	if err != nil {
		return "", "", errors.Errorf("DB.GetOneUnAllocatedAccountPubKeyTable(): error: %s", err)
	}
	storedAddr := pubkeyTable.WalletAddress //change from w.BTC.StoredAddress()
	storedAccount := pubkeyTable.Account    //change from w.BTC.StoredAccountName()

	// 1.CreateRawTransaction(仮で作成し、この後サイズから手数料を算出する)
	msgTx, err := w.btc.CreateRawTransaction(storedAddr, inputTotal, inputs)
	if err != nil {
		return "", "", errors.Errorf("BTC.CreateRawTransaction(): error: %s", err)
	}

	// 2.fee算出
	fee, err := w.btc.GetFee(msgTx, adjustmentFee)

	// 3.手数料のために、totalを調整し、再度RawTransactionを作成する
	//このパートは、出金とロジックが異なる
	outputTotal = inputTotal - fee
	if outputTotal <= 0 {
		return "", "", errors.Errorf("calculated fee must be wrong: fee:%v, error: %s", fee, err)
	}
	w.logger.Debug(
		"Total Coin to send:%d(Satoshi) after fee calculated, input length: %d",
		zap.Any("amount",outputTotal),
		zap.Int("len(inputs)", len(inputs)))

	// 4.outputs作成
	txReceiptOutputs := []walletrepo.TxOutput{
		{
			ReceiptID:     0,
			OutputAddress: storedAddr,
			OutputAccount: storedAccount,
			OutputAmount:  w.btc.AmountString(outputTotal),
			IsChange:      false,
		},
	}

	// 5.再度 CreateRawTransaction
	msgTx, err = w.btc.CreateRawTransaction(storedAddr, outputTotal, inputs)
	if err != nil {
		return "", "", errors.Errorf("BTC.CreateRawTransaction(): error: %s", err)
	}

	// 6.出力用にHexに変換する
	hex, err := w.btc.ToHex(msgTx)
	if err != nil {
		return "", "", errors.Errorf("BTC.ToHex(msgTx): error: %s", err)
	}

	// 7. Databaseに必要な情報を保存
	txReceiptID, err := w.insertTxTableForUnsigned(actionType, hex, inputTotal, outputTotal, fee, enum.TxTypeValue[enum.TxTypeUnsigned], txReceiptInputs, txReceiptOutputs, nil)
	if err != nil {
		return "", "", errors.Errorf("insertTxTableForUnsigned(): error: %s", err)
	}

	// 8. serialize previous txs for multisig signature
	encodedAddrsPrevs, err := serial.EncodeToString(*addrsPrevs)
	if err != nil {
		return "", "", errors.Errorf("serial.EncodeToString(): error: %s", err)
	}
	w.logger.Debug("encodedAddrsPrevs", zap.String("encodedAddrsPrevs", encodedAddrsPrevs))

	// 9. GCSにトランザクションファイルを作成
	//TODO:本来、この戻り値をDumpして、GCSに保存、それをDLして、USBに入れてコールドウォレットに移動しなくてはいけない
	//TODO:Debug時はlocalに出力することとする。=> これはフラグで判別したほうがいいかもしれない/Interface型にして対応してもいいかも
	//TODO:ここでエラーが起きると、再生成するために、insertされたデータを削除しないといけない。DB更新より先にファイルを作成したほうがいい？？
	var generatedFileName string
	if txReceiptID != 0 {
		generatedFileName, err = w.storeHex(hex, encodedAddrsPrevs, txReceiptID, actionType)
		if err != nil {
			return "", "", errors.Errorf("wallet.storeHex(): error: %s", err)
		}
	}

	// 10. 入金準備に入ったことをユーザーに通知
	// TODO:NatsのPublisherとして通知すればいいか？

	return hex, generatedFileName, nil
}

//TODO:引数の数が多いのはGoにおいてはBad practice...
//[共通(receipt/payment/transfer)]
func (w *Wallet) insertTxTableForUnsigned(actionType enum.ActionType, hex string, inputTotal, outputTotal, fee btcutil.Amount, txType uint8,
	txInputs []walletrepo.TxInput, txOutputs []walletrepo.TxOutput, paymentRequestIds []int64) (int64, error) {

	//1.内容が同じだと、生成されるhexもまったく同じ為、同一のhexが合った場合は処理をskipする
	count, err := w.storager.GetTxCountByUnsignedHex(actionType, hex)
	if err != nil {
		return 0, errors.Errorf("DB.GetTxCountByUnsignedHex(): error: %s", err)
	}
	if count != 0 {
		//skip
		return 0, nil
	}

	//2.TxReceiptテーブル
	txReceipt := walletrepo.TxTable{}
	txReceipt.UnsignedHexTx = hex
	txReceipt.TotalInputAmount = w.btc.AmountString(inputTotal)
	txReceipt.TotalOutputAmount = w.btc.AmountString(outputTotal)
	txReceipt.Fee = w.btc.AmountString(fee)
	txReceipt.TxType = txType

	tx := w.storager.MustBegin()
	txReceiptID, err := w.storager.InsertTxForUnsigned(actionType, &txReceipt, tx, false)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxForUnsigned(): error: %s", err)
	}

	//3.TxReceiptInputテーブル
	//ReceiptIDの更新
	for idx := range txInputs {
		txInputs[idx].ReceiptID = txReceiptID
	}
	err = w.storager.InsertTxInputForUnsigned(actionType, txInputs, tx, false)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxInputForUnsigned(): error: %s", err)
	}

	//4.TxReceiptOutputテーブル
	//ReceiptIDの更新
	for idx := range txOutputs {
		txOutputs[idx].ReceiptID = txReceiptID
	}

	//commit flag
	//paymentのみ、後続の処理が存在する(payment_requestテーブル)
	isCommit := true
	if actionType == enum.ActionTypePayment {
		isCommit = false
	}

	err = w.storager.InsertTxOutputForUnsigned(actionType, txOutputs, tx, isCommit)
	if err != nil {
		return 0, errors.Errorf("DB.InsertTxOutputForUnsigned(): error: %s", err)
	}

	//TODO:未着手
	//5.Toに指定されたaccount_pubkey_receiptなどの使用されたwalletのis_allocatedを1に更新する

	//6. payment_requestのpayment_idを更新する paymentRequestIds
	if actionType == enum.ActionTypePayment {
		//txReceiptID
		_, err = w.storager.UpdatePaymentIDOnPaymentRequest(txReceiptID, paymentRequestIds, tx, true)
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
	path := txfile.CreateFilePath(actionType, enum.TxTypeUnsigned, id, true)
	generatedFileName, err = txfile.WriteFile(path, savedata)
	if err != nil {
		return "", errors.Errorf("txfile.WriteFile(): error: %s", err)
	}

	return generatedFileName, nil
}
