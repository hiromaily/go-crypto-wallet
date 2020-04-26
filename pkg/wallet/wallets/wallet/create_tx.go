package wallet

import (
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/serial"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
)

// create_tx.go is for common func among create transaction fuctionalites

type parsedTx struct {
	txInputs       []btcjson.TransactionInput
	txRepoTxInputs []*models.TXInput
	prevTxs        []btc.PrevTx
	addresses      []string //input, sender's address
}

// - check unspentTx for sender account
// - get utxo and create unsigned tx
// - fee is fluctuating, so outdated unsigned tx must not be used, re-create unsigned tx
// - after signed tx sent, utxo could not be retrieved by listUnspent()

// create unsigned tx
// - sender account: client, receiver account: receipt
// - if amount=0, all coin is sent
// FIXME: receiver account covers fee, but is should be flexible
// TODO: change functionality is not implemented yet
// TODO: after this func, what if `listtransactions` api is called to see result
func (w *Wallet) createTx(
	sender,
	receiver account.AccountType,
	targetAction action.ActionType,
	requiredAmount btcutil.Amount,
	adjustmentFee float64,
	paymentRequestIds []int64,
	userPayments []UserPayment) (string, string, error) {

	//amount
	// - receiptAction: it's 0. total amount in clients to receipt
	// - transferAction: if 0, total amount of sender to receiver
	//                   if not 0, amount is sent from sender to receiver
	// - paymentAction: total amount in payment users

	w.logger.Debug("createTx()",
		zap.String("sender_acount", sender.String()),
		zap.String("receiver_acount", receiver.String()),
		zap.String("target_action", targetAction.String()),
		zap.Any("required_amount", requiredAmount),
		zap.Float64("adjustmentFee", adjustmentFee))

	// 1. get listUnspent
	unspentList, unspentAddrs, err := w.getUnspentList(sender)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call getUnspentList()")
	}
	if len(unspentList) == 0 {
		w.logger.Info("no listunspent")
		return "", "", nil
	}
	w.logger.Debug("getUnspentList()",
		zap.Any("unspentList", unspentList),
		zap.Any("unspentAddrs", unspentAddrs))

	// 2. parse listUnspent
	parsedTx, inputTotal, isDone := w.parseListUnspentTx(unspentList, requiredAmount)
	if len(parsedTx.txInputs) == 0 {
		w.logger.Info("no input tx in listUnspent")
		return "", "", nil
	}
	if !isDone {
		return "", "", errors.New("sender account can't meet amount to send")
	}
	w.logger.Debug(
		"amount",
		zap.Any("requiredAmount", requiredAmount),
		zap.Any("input_total", inputTotal),
		zap.Int("len(inputs)", len(parsedTx.txInputs)),
	)
	if requiredAmount != 0 {
		w.logger.Debug("amount", zap.Any("expected_change", inputTotal-requiredAmount))
	}

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:         parsedTx.addresses,
		PrevTxs:       parsedTx.prevTxs,
		SenderAccount: sender,
	}

	// 3. create txOutputs
	var txPrevOutputs map[btcutil.Address]btcutil.Amount
	switch targetAction {
	case action.ActionTypeReceipt, action.ActionTypeTransfer:
		var isChange bool
		if requiredAmount != 0 {
			isChange = true
		}
		txPrevOutputs, err = w.createTxOutputs(receiver, requiredAmount, inputTotal, unspentAddrs[0], isChange)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call createTxOutputs()")
		}
	case action.ActionTypePayment:
		changeAddr := unspentAddrs[0] //this is actually sender's address because it's for change
		changeAmount := inputTotal - requiredAmount
		w.logger.Debug("before createPaymentOutputs()",
			zap.Any("change_addr", changeAddr),
			zap.Any("change_amount", changeAmount))
		txPrevOutputs = w.createPaymentOutputs(userPayments, changeAddr, changeAmount)
		w.logger.Debug("txPrevOutputs",
			zap.Int("len(txPrevOutputs)", len(txPrevOutputs)))

	default:
		return "", "", errors.Errorf("invalid actionType: %s", targetAction)
	}
	w.logger.Debug("txPrevOutputs",
		zap.Int("len(txPrevOutputs)", len(txPrevOutputs)))
	//"txPrevOutputsError":"json: unsupported type: map[btcutil.Address]btcutil.Amount"

	// create raw tx
	hex, fileName, err := w.createRawTx(
		targetAction,
		sender,
		receiver,
		adjustmentFee,
		parsedTx.txInputs,
		inputTotal,
		parsedTx.txRepoTxInputs,
		txPrevOutputs,
		&addrsPrevs,
		paymentRequestIds)

	//TODO: notify what unsigned tx is created
	// => in command pkg

	return hex, fileName, err
}

// call API `unspentlist`
// no result and no error is possible, so caller should check both returned value
func (w *Wallet) getUnspentList(accountType account.AccountType) ([]btc.ListUnspentResult, []string, error) {
	// unlock locked UnspentTransaction
	//if err := w.BTC.UnlockUnspent(); err != nil {
	//	return "", "", err
	//}

	// get listUnspent
	unspentList, err := w.btc.ListUnspentByAccount(accountType)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call btc.Client().ListUnspentByAccount()")
	}
	unspentAddrs := w.btc.GetUnspentListAddrs(unspentList, accountType)

	return unspentList, unspentAddrs, nil
}

// parse result of listUnspent
// retured *parsedTx could be nil
func (w *Wallet) parseListUnspentTx(unspentList []btc.ListUnspentResult, amount btcutil.Amount) (*parsedTx, btcutil.Amount, bool) {
	var inputTotal btcutil.Amount
	txInputs := make([]btcjson.TransactionInput, 0, len(unspentList))
	txRepoTxInputs := make([]*models.TXInput, 0, len(unspentList))
	prevTxs := make([]btc.PrevTx, 0, len(unspentList))
	addresses := make([]string, 0, len(unspentList))

	var isDone bool //if isDone is false, sender can't meet amount
	if amount == 0 {
		isDone = true
	}

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

		txRepoTxInputs = append(txRepoTxInputs, &models.TXInput{
			TXID:               0,
			InputTxid:          tx.TxID,
			InputVout:          tx.Vout,
			InputAddress:       tx.Address,
			InputAccount:       tx.Label,
			InputAmount:        w.btc.FloatToDecimal(tx.Amount),
			InputConfirmations: uint64(tx.Confirmations),
		})

		//TODO: if sender is client account, RedeemScript is blank
		prevTxs = append(prevTxs, btc.PrevTx{
			Txid:         tx.TxID,
			Vout:         tx.Vout,
			ScriptPubKey: tx.ScriptPubKey,
			RedeemScript: tx.RedeemScript, //required if target account is multsig address
			Amount:       tx.Amount,
		})

		addresses = append(addresses, tx.Address)

		// check total if amount is set as parameter
		if amount == 0 {
			continue
		}
		if inputTotal > amount {
			isDone = true
			break
		}
	}

	return &parsedTx{
		txInputs:       txInputs,
		txRepoTxInputs: txRepoTxInputs,
		prevTxs:        prevTxs,
		addresses:      addresses,
	}, inputTotal, isDone
}

//TODO: this code should be integrated with ActionTypePayment
func (w *Wallet) createTxOutputs(
	reciver account.AccountType,
	requiredAmount btcutil.Amount,
	inputTotal btcutil.Amount,
	senderAddr string,
	isChange bool) (map[btcutil.Address]btcutil.Amount, error) {

	// 1. get unallocated address for receiver
	// - receipt/transfer
	pubkeyTable, err := w.pubkeyRepo.GetOneUnAllocated(reciver)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call pubkeyRepo.GetOneUnAllocated()")
	}
	receiverAddr := pubkeyTable.WalletAddress

	// 2. create receiver txOutput
	receiverDecodedAddr, err := btcutil.DecodeAddress(receiverAddr, w.btc.GetChainConf())
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call btcutil.DecodeAddress(%s)", receiverAddr)
	}
	txPrevOutputs := make(map[btcutil.Address]btcutil.Amount)
	if isChange {
		txPrevOutputs[receiverDecodedAddr] = requiredAmount //satoshi
	} else {
		txPrevOutputs[receiverDecodedAddr] = inputTotal //satoshi
	}
	w.logger.Debug("receiver txOutput",
		zap.String("receiverAddr", receiverAddr),
		zap.Any("receivedAmount", txPrevOutputs[receiverDecodedAddr]))

	// 3. if change is required
	if isChange {
		w.logger.Debug("change is required")
		senderDecodedAddr, err := btcutil.DecodeAddress(senderAddr, w.btc.GetChainConf())
		if err != nil {
			return nil, errors.Wrapf(err, "fail to call btcutil.DecodeAddress(%s)", receiverAddr)
		}

		// for change
		// TODO: what if (inputTotal - requiredAmount) = 0,
		//  fee can not be paid from txOutput for change
		txPrevOutputs[senderDecodedAddr] = inputTotal - requiredAmount

		w.logger.Debug("change(sender) txOutput",
			zap.String("senderAddr", senderAddr),
			zap.Any("inputTotal - requiredAmount", inputTotal-requiredAmount))

	}
	return txPrevOutputs, nil
}

// createRawTx create raw tx
// - calculate fee
// - create raw tx
// - insert data to detabase
// - available from receipt/transfer action
//TODO: is_allocated should be updated to true when tx sent
func (w *Wallet) createRawTx(
	targetAction action.ActionType,
	sender account.AccountType,
	receiver account.AccountType,
	adjustmentFee float64,
	txInputs []btcjson.TransactionInput,
	inputTotal btcutil.Amount,
	txRepoTxInputs []*models.TXInput,
	txPrevOutputs map[btcutil.Address]btcutil.Amount,
	addrsPrevs *btc.AddrsPrevTxs,
	paymentRequestIds []int64) (string, string, error) {

	// 1. create raw transaction as temporary use
	// - receipt/transfer
	// - later calculate by tx size
	msgTx, err := w.btc.CreateRawTransaction(txInputs, txPrevOutputs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.CreateRawTransaction()")
	}

	// 2. calculate fee and output total
	// - receipt/transfer
	//  - adjust outputTotal by fee and re-run CreateRawTransaction
	//  - this logic would be different from payment
	//TODO: payment has different logic
	outputTotal, fee, txOutputs, txRepoTxOutputs, err := w.calculateOutputTotal(sender, receiver, msgTx, adjustmentFee, inputTotal, txPrevOutputs)
	if err != nil {
		return "", "", err
	}
	w.logger.Debug("txOutputs", zap.Int("len(txOutputs)", len(txOutputs)))
	w.logger.Debug("txRepoTxOutputs", zap.Int("len(txRepoTxOutputs)", len(txRepoTxOutputs)))
	//zap.Any("txOutputs", txOutputs),
	//zap.Any("txRepoTxOutputs", txRepoTxOutputs),

	//for debug
	//map[btcutil.Address]btcutil.Amount
	for addr, amt := range txOutputs {
		w.logger.Debug("txOutputs",
			zap.String("address_string", addr.String()),
			zap.String("address_encoded", addr.EncodeAddress()),
			zap.String("amount", amt.String()),
		)
	}
	for _, v := range txRepoTxOutputs {
		w.logger.Debug("txRepoTxOutputs",
			zap.String("output_account", v.OutputAccount),
			zap.String("output_address", v.OutputAddress),
			zap.String("output_amount", v.OutputAmount.String()),
		)
	}

	// 3. re call CreateRawTransaction
	// - receipt/transfer
	msgTx, err = w.btc.CreateRawTransaction(txInputs, txOutputs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.CreateRawTransaction()")
	}

	// 4. convert msgTx to hex
	// - receipt/transfer/payment
	hex, err := w.btc.ToHex(msgTx)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.ToHex(msgTx)")
	}

	// 5. insert to tx_table for unsigned tx
	// - receipt/transfer/payment
	//  - txID would be 0 if record is already existing then csv file is not created
	txID, err := w.insertTxTableForUnsigned(
		targetAction,
		hex,
		inputTotal,
		outputTotal,
		fee,
		txRepoTxInputs,
		txRepoTxOutputs,
		paymentRequestIds)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call insertTxTableForUnsigned()")
	}

	// 6. serialize previous txs for multisig signature
	// - receipt/transfer/payment
	encodedAddrsPrevs, err := serial.EncodeToString(*addrsPrevs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call serial.EncodeToString()")
	}
	w.logger.Debug("encodedAddrsPrevs", zap.String("encodedAddrsPrevs", encodedAddrsPrevs))

	// 7. generate tx file
	// - receipt/transfer/payment
	//TODO: how to recover when error occurred here
	// - inserted data in database must be deleted to generate hex file
	var generatedFileName string
	if txID != 0 {
		generatedFileName, err = w.generateHexFile(targetAction, hex, encodedAddrsPrevs, txID)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call generateHexFile()")
		}
	}

	return hex, generatedFileName, nil
}

func (w *Wallet) calculateOutputTotal(
	sender account.AccountType,
	receiver account.AccountType,
	msgTx *wire.MsgTx,
	adjustmentFee float64,
	inputTotal btcutil.Amount,
	txPrevOutputs map[btcutil.Address]btcutil.Amount,
) (btcutil.Amount, btcutil.Amount, map[btcutil.Address]btcutil.Amount, []*models.TXOutput, error) {

	// get fee
	fee, err := w.btc.GetFee(msgTx, adjustmentFee)
	if err != nil {
		return 0, 0, nil, nil, errors.Wrap(err, "fail to call btc.GetFee()")
	}
	var outputTotal btcutil.Amount
	txRepoOutputs := make([]*models.TXOutput, 0, len(txPrevOutputs))

	// subtract fee from output transaction for change
	// FIXME: what if change is short, should re-run form the beginning with shortage-flag
	for addr, amt := range txPrevOutputs {
		if len(txPrevOutputs) == 1 {
			//no change
			txPrevOutputs[addr] -= fee
			txRepoOutputs = append(txRepoOutputs, &models.TXOutput{
				TXID:          0,
				OutputAddress: addr.String(),
				OutputAccount: receiver.String(),
				OutputAmount:  w.btc.AmountToDecimal(amt - fee),
				IsChange:      false,
			})
			outputTotal += amt
			break
		}

		if acnt, _ := w.btc.GetAccount(addr.String()); acnt == sender.String() {
			w.logger.Debug("detect sender account in calculateOutputTotal")
			//chang address
			txPrevOutputs[addr] -= fee
			txRepoOutputs = append(txRepoOutputs, &models.TXOutput{
				TXID:          0,
				OutputAddress: addr.String(),
				OutputAccount: sender.String(),
				OutputAmount:  w.btc.AmountToDecimal(amt - fee),
				IsChange:      true,
			})
		} else {
			txRepoOutputs = append(txRepoOutputs, &models.TXOutput{
				TXID:          0,
				OutputAddress: addr.String(),
				OutputAccount: receiver.String(),
				OutputAmount:  w.btc.AmountToDecimal(amt),
				IsChange:      false,
			})
		}
		outputTotal += amt
	}
	w.logger.Debug("calculateOutputTotal",
		zap.Any("fee", fee),
		zap.Any("outputTotal (before fee adjustment)", outputTotal),
		zap.Any("outputTotal by (inputTotal - fee)", inputTotal-fee),
		zap.Any("outputTotal by (outputTotal - fee)", outputTotal-fee),
	)

	//result of total should be same
	outputTotal = inputTotal - fee //for no change type
	//outputTotal -= fee

	if outputTotal <= 0 {
		w.logger.Debug(
			"inputTotal is short of coin to pay fee",
			zap.Any("amount of inputTotal", inputTotal),
			zap.Any("fee", fee))
		return 0, 0, nil, nil, errors.Wrapf(err, "inputTotal is short of coin to pay fee")
	}

	return outputTotal, fee, txPrevOutputs, txRepoOutputs, nil
}

// - available from receipt/payment/transfer action
func (w *Wallet) insertTxTableForUnsigned(
	actionType action.ActionType,
	hex string,
	inputTotal,
	outputTotal,
	fee btcutil.Amount,
	txInputs []*models.TXInput,
	txOutputs []*models.TXOutput,
	paymentRequestIds []int64) (int64, error) {

	// 1. skip if same hex is already stored
	//count, err := w.repo.GetTxCountByUnsignedHex(actionType, hex)
	count, err := w.txRepo.GetCountByUnsignedHex(actionType, hex)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call txRepo.GetCount()")
	}
	if count != 0 {
		//skip
		return 0, nil
	}

	// 2.TxReceipt table //TODO: remove after replacement is done
	//txReceipt := walletrepo.TxTable{}
	//txReceipt.UnsignedHexTx = hex
	//txReceipt.TotalInputAmount = w.btc.AmountString(inputTotal)
	//txReceipt.TotalOutputAmount = w.btc.AmountString(outputTotal)
	//txReceipt.Fee = w.btc.AmountString(fee)
	//txReceipt.TxType = txType
	txItem := &models.TX{
		Action:            action.ActionTypePayment.String(),
		UnsignedHexTX:     hex,
		TotalInputAmount:  w.btc.AmountToDecimal(inputTotal),
		TotalOutputAmount: w.btc.AmountToDecimal(outputTotal),
		Fee:               w.btc.AmountToDecimal(fee),
	}

	// start db transaction //TODO: implement transaction
	//tx := w.repo.MustBegin()
	txID, err := w.txRepo.InsertUnsignedTx(actionType, txItem)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call txRepo.InsertUnsignedTx()")
	}

	// 3.TxReceiptInput table
	// update ReceiptID
	for idx := range txInputs {
		txInputs[idx].TXID = txID
	}
	err = w.txInRepo.InsertBulk(txInputs)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call txInRepo.InsertBulk()")
	}

	// 4.TxReceiptOutput table
	// update ReceiptID
	for idx := range txOutputs {
		txOutputs[idx].TXID = txID
	}
	//commit flag //TODO: transaction with defer
	//isCommit := true
	//if actionType == action.ActionTypePayment {
	//	isCommit = false
	//}
	err = w.txOutRepo.InsertBulk(txOutputs)
	if err != nil {
		return 0, errors.Wrap(err, "storager.InsertTxOutputForUnsigned()")
	}

	//TODO: not implemented yet
	// 5. address for receiver account should be updated `is_allocated=1`

	// 6. update payment_id in payment_request table for only action.ActionTypePayment
	if actionType == action.ActionTypePayment {
		_, err = w.payReqRepo.UpdatePaymentID(txID, paymentRequestIds) //TODO: transaction commit
		if err != nil {
			return 0, errors.Wrap(err, "storager.UpdatePaymentIDOnPaymentRequest()")
		}
	}

	return txID, nil
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
	path := w.txFileRepo.CreateFilePath(actionType, tx.TxTypeUnsigned, id, 0)
	generatedFileName, err = w.txFileRepo.WriteFile(path, savedata)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.WriteFile()")
	}

	return generatedFileName, nil
}

// IsFoundTxIDAndVout finds out txID and vout from related txInputs
// nolint: unused
func (w *Wallet) IsFoundTxIDAndVout(txID string, vout uint32, inputs []btcjson.TransactionInput) bool {
	for _, val := range inputs {
		if val.Txid == txID && val.Vout == vout {
			return true
		}
	}
	return false
}
