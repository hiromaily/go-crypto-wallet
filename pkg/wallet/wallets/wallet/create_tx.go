package wallet

import (
	"fmt"
	"github.com/bookerzzz/grok"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/serial"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
)

type parsedTx struct {
	txInputs       []btcjson.TransactionInput
	txRepoTxInputs []walletrepo.TxInput
	prevTxs        []btc.PrevTx
	addresses      []string
}

// - check unspentTx for sender account
// - get utxo and create unsigned tx
// - fee is fluctuating, so outdated unsigned tx must not be used, re-create unsigned tx
// - after signed tx sent, utxo could not be retrieved by listUnspent()

// create unsigned tx
// - sender account: client, receiver account: receipt
// - receiver account covers fee
func (w *Wallet) createTx(
	sender,
	receiver account.AccountType,
	targetAction action.ActionType,
	amount btcutil.Amount,
	adjustmentFee float64) (string, string, error) {

	// get listUnspent
	unspentList, err := w.getUnspentList(sender)
	if len(unspentList) == 0 {
		w.logger.Info("no listunspent")
		return "", "", nil
	}

	// parse listUnspent
	parsedTx, inputTotal := w.parseListUnspentTx(unspentList, amount)
	w.logger.Debug(
		"total coin to send (Satoshi) before fee calculated",
		zap.Any("input_amount", inputTotal),
		zap.Int("len(inputs)", len(parsedTx.txInputs)))
	if len(parsedTx.txInputs) == 0 {
		return "", "", nil
	}

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:         parsedTx.addresses,
		PrevTxs:       parsedTx.prevTxs,
		SenderAccount: sender,
	}

	//TODO: how this code can be integrated with CreateTransferTx ??

	// create raw tx
	hex, fileName, err := w.createRawTx(
		targetAction,
		receiver,
		adjustmentFee,
		parsedTx.txInputs,
		inputTotal,
		parsedTx.txRepoTxInputs,
		&addrsPrevs)

	//TODO: notify what unsigned tx is created
	// => in command pkg

	return hex, fileName, err
}

// call API `unspentlist`
// no result and no error is possible, so caller should check both returned value
func (w *Wallet) getUnspentList(accountType account.AccountType) ([]btc.ListUnspentResult, error) {
	// unlock locked UnspentTransaction
	//if err := w.BTC.UnlockAllUnspentTransaction(); err != nil {
	//	return "", "", err
	//}

	// get listUnspent
	unspentList, _, err := w.btc.ListUnspentByAccount(accountType)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call btc.Client().ListUnspent()")
	}
	grok.Value(unspentList)
	if len(unspentList) == 0 {
		w.logger.Info("no listunspent")
		return nil, nil
	}
	return unspentList, nil
}

// parse result of listUnspent
// retured *parsedTx could be nil
func (w *Wallet) parseListUnspentTx(unspentList []btc.ListUnspentResult, amount btcutil.Amount) (*parsedTx, btcutil.Amount) {
	var inputTotal btcutil.Amount
	txInputs := make([]btcjson.TransactionInput, 0, len(unspentList))
	txRepoTxInputs := make([]walletrepo.TxInput, 0, len(unspentList))
	prevTxs := make([]btc.PrevTx, 0, len(unspentList))
	addresses := make([]string, 0, len(unspentList))

	for _, tx := range unspentList {
		// Amount
		amt, err := btcutil.NewAmount(tx.Amount)
		if err != nil {
			//this error is not expected
			w.logger.Error(
				"btcutil.NewAmount()",
				zap.Float64("tx amount", tx.Amount),
				zap.Error(err))
			continue
		}
		inputTotal += amt

		txInputs = append(txInputs, btcjson.TransactionInput{
			Txid: tx.TxID,
			Vout: tx.Vout,
		})

		txRepoTxInputs = append(txRepoTxInputs, walletrepo.TxInput{
			ReceiptID:          0,
			InputTxid:          tx.TxID,
			InputVout:          tx.Vout,
			InputAddress:       tx.Address,
			InputAccount:       tx.Label,
			InputAmount:        fmt.Sprintf("%f", tx.Amount),
			InputConfirmations: tx.Confirmations,
		})

		prevTxs = append(prevTxs, btc.PrevTx{
			Txid:         tx.TxID,
			Vout:         tx.Vout,
			ScriptPubKey: tx.ScriptPubKey,
			RedeemScript: "", //required if target account is multsig address
			Amount:       tx.Amount,
		})

		addresses = append(addresses, tx.Address)

		// check total if amount is set as parameter
		if amount == 0 {
			continue
		}
		if inputTotal > amount {
			break
		}
	}

	return &parsedTx{
		txInputs:       txInputs,
		txRepoTxInputs: txRepoTxInputs,
		prevTxs:        prevTxs,
		addresses:      addresses,
	}, inputTotal
}

// createRawTx create raw tx
// - calculate fee
// - create raw tx
// - insert data to detabase
// - available from receipt/transfer action
//TODO: is_allocated should be updated to true when tx sent
func (w *Wallet) createRawTx(
	actionType action.ActionType,
	receiverAccountType account.AccountType,
	adjustmentFee float64,
	txInputs []btcjson.TransactionInput,
	inputTotal btcutil.Amount,
	txRepoTxInputs []walletrepo.TxInput,
	addrsPrevs *btc.AddrsPrevTxs) (string, string, error) {

	// 1. get unallocated address for receiver
	pubkeyTable, err := w.storager.GetOneUnAllocatedAccountPubKeyTable(receiverAccountType)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call storager.GetOneUnAllocatedAccountPubKeyTable()")
	}
	receiverAddr := pubkeyTable.WalletAddress
	//storedAccount := pubkeyTable.Account //used to OutputAccount before

	// 2. create raw transaction as temporary use
	//  - later calculate by tx size
	msgTx, err := w.btc.CreateRawTransaction(receiverAddr, inputTotal, txInputs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.CreateRawTransaction()")
	}

	// 3. calculate fee and output total
	//  - adjust outputTotal by fee and re-run CreateRawTransaction
	//  - this logic would be different from payment
	outputTotal, fee, err := w.calculateOutputTotal(msgTx, adjustmentFee, inputTotal)
	if err != nil {
		return "", "", err
	}
	w.logger.Debug(
		"total coin to send (Satoshi) after fee calculated",
		zap.Any("output_amount", outputTotal),
		zap.Int("len(inputs)", len(txInputs)))

	txRepoTxOutputs := []walletrepo.TxOutput{
		{
			ReceiptID:     0,
			OutputAddress: receiverAddr,
			OutputAccount: receiverAccountType.String(),
			OutputAmount:  w.btc.AmountString(outputTotal),
			IsChange:      false,
		},
	}

	// 4. re call CreateRawTransaction by output
	msgTx, err = w.btc.CreateRawTransaction(receiverAddr, outputTotal, txInputs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.CreateRawTransaction()")
	}

	// 5. convert msgTx to hex
	hex, err := w.btc.ToHex(msgTx)
	if err != nil {
		return "", "", errors.Errorf("BTC.ToHex(msgTx): error: %s", err)
	}

	// 6. insert to tx_table for unsigned tx
	//  - txReceiptID would be 0 if record is already existing then csv file is not created
	txReceiptID, err := w.insertTxTableForUnsigned(
		actionType,
		hex,
		inputTotal,
		outputTotal,
		fee,
		tx.TxTypeValue[tx.TxTypeUnsigned],
		txRepoTxInputs,
		txRepoTxOutputs,
		nil)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call insertTxTableForUnsigned()")
	}

	// 7. serialize previous txs for multisig signature
	encodedAddrsPrevs, err := serial.EncodeToString(*addrsPrevs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call serial.EncodeToString()")
	}
	w.logger.Debug("encodedAddrsPrevs", zap.String("encodedAddrsPrevs", encodedAddrsPrevs))

	// 8. generate tx file
	//TODO: how to recover when error occurred here
	// - inserted data in database must be deleted to generate hex file
	var generatedFileName string
	if txReceiptID != 0 {
		generatedFileName, err = w.generateHexFile(actionType, hex, encodedAddrsPrevs, txReceiptID)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call generateHexFile()")
		}
	}

	return hex, generatedFileName, nil
}

func (w *Wallet) calculateOutputTotal(msgTx *wire.MsgTx, adjustmentFee float64, inputTotal btcutil.Amount) (btcutil.Amount, btcutil.Amount, error) {
	var outputTotal btcutil.Amount
	fee, err := w.btc.GetFee(msgTx, adjustmentFee)
	outputTotal = inputTotal - fee
	if outputTotal <= 0 {
		w.logger.Debug(
			"inputTotal is short of coin to pay fee",
			zap.Any("amount", inputTotal),
			zap.Any("len(inputs)", fee))
		return 0, 0, errors.Wrapf(err, "inputTotal is short of coin to pay fee")
	}
	return outputTotal, fee, nil
}

// - available from receipt/payment/transfer action
func (w *Wallet) insertTxTableForUnsigned(
	actionType action.ActionType,
	hex string,
	inputTotal,
	outputTotal,
	fee btcutil.Amount,
	txType uint8,
	txInputs []walletrepo.TxInput,
	txOutputs []walletrepo.TxOutput,
	paymentRequestIds []int64) (int64, error) {

	// 1. skip if same hex is already stored
	count, err := w.storager.GetTxCountByUnsignedHex(actionType, hex)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call storager.GetTxCountByUnsignedHex()")
	}
	if count != 0 {
		//skip
		return 0, nil
	}

	// 2.TxReceipt table
	txReceipt := walletrepo.TxTable{}
	txReceipt.UnsignedHexTx = hex
	txReceipt.TotalInputAmount = w.btc.AmountString(inputTotal)
	txReceipt.TotalOutputAmount = w.btc.AmountString(outputTotal)
	txReceipt.Fee = w.btc.AmountString(fee)
	txReceipt.TxType = txType

	// start db transaction
	tx := w.storager.MustBegin()
	txReceiptID, err := w.storager.InsertTxForUnsigned(actionType, &txReceipt, tx, false)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call storager.InsertTxForUnsigned()")
	}

	// 3.TxReceiptInput table
	// update ReceiptID
	for idx := range txInputs {
		txInputs[idx].ReceiptID = txReceiptID
	}
	err = w.storager.InsertTxInputForUnsigned(actionType, txInputs, tx, false)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call storager.InsertTxInputForUnsigned()")
	}

	// 4.TxReceiptOutput table
	// update ReceiptID
	for idx := range txOutputs {
		txOutputs[idx].ReceiptID = txReceiptID
	}
	//commit flag
	isCommit := true
	if actionType == action.ActionTypePayment {
		isCommit = false
	}
	err = w.storager.InsertTxOutputForUnsigned(actionType, txOutputs, tx, isCommit)
	if err != nil {
		return 0, errors.Wrap(err, "storager.InsertTxOutputForUnsigned()")
	}

	//TODO: not implemented yet
	// 5. address for receiver account should be updated `is_allocated=1`

	// 6. update payment_id in payment_request table for only action.ActionTypePayment
	if actionType == action.ActionTypePayment {
		_, err = w.storager.UpdatePaymentIDOnPaymentRequest(txReceiptID, paymentRequestIds, tx, true)
		if err != nil {
			return 0, errors.Wrap(err, "storager.UpdatePaymentIDOnPaymentRequest()")
		}
	}

	return txReceiptID, nil
}

// generateHexFile generate file for hex and encoded previous addresses
// - available from receipt/payment/transfer action
func (w *Wallet) generateHexFile(actionType action.ActionType, hex, encodedAddrsPrevs string, id int64) (string, error) {
	var (
		generatedFileName string
		err               error
	)

	savedata := hex
	if encodedAddrsPrevs != "" {
		savedata = fmt.Sprintf("%s,%s", savedata, encodedAddrsPrevs)
	}

	// create file
	path := w.txFileRepo.CreateFilePath(actionType, tx.TxTypeUnsigned, id)
	generatedFileName, err = w.txFileRepo.WriteFile(path, savedata)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.WriteFile()")
	}

	return generatedFileName, nil
}
