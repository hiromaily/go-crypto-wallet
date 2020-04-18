package wallet

import (
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

// CreatePaymentTx create unsigned tx for user(anonymous addresses)
// sender: payment, receiver: addresses coming from user_payment table
func (w *Wallet) CreatePaymentTx(adjustmentFee float64) (string, string, error) {
	sender := account.AccountTypePayment
	//receiver := account.AccountTypeAnonymous
	//targetAction := action.ActionTypePayment

	// get payment data from payment_request
	userPayments, paymentRequestIds, err := w.createUserPayment()
	if err != nil {
		return "", "", err
	}
	if len(userPayments) == 0 {
		w.logger.Debug("no data in userPayments")
		// no data
		return "", "", nil
	}

	// calculate total amount to send from payment_request
	var userTotal btcutil.Amount
	for _, val := range userPayments {
		userTotal += val.validAmount
	}

	// get balance for payment account
	balance, err := w.btc.GetReceivedByLabelAndMinConf(account.AccountTypePayment.String(), w.btc.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= userTotal {
		//balance is short
		return "", "", errors.New("balance for payment account is insufficient")
	}
	w.logger.Debug("balane_userTotal",
		zap.Any("balance", balance),
		zap.Any("userTotal", userTotal))

	//FIXME: how to commonalize code from here

	// get listUnspent
	unspentList, unspentAddrs, err := w.getUnspentList(sender)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call getUnspentList()")
	}
	if len(unspentList) == 0 {
		w.logger.Info("no listunspent")
		return "", "", nil
	}

	// parse listUnspent
	parsedTx, inputTotal, isDone := w.parseListUnspentTx(unspentList, userTotal)
	w.logger.Debug(
		"total coin to send (Satoshi) before fee calculated",
		zap.Any("input_amount", inputTotal),
		zap.Int("len(inputs)", len(parsedTx.txInputs)))
	if len(parsedTx.txInputs) == 0 {
		return "", "", nil
	}
	if !isDone {
		return "", "", errors.New("sender account can't meet amount to send")
	}

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:         parsedTx.addresses,
		PrevTxs:       parsedTx.prevTxs,
		SenderAccount: sender,
	}

	// payment logic exclusively
	// create payment txOutputs
	changeAddr := unspentAddrs[0].String() //this is actually sender's address because it's for change
	changeAmount := inputTotal - userTotal
	txOutputs := w.createPaymentOutputs(userPayments, changeAddr, changeAmount)

	// create raw tx
	return w.createRawPaymentTx(
		adjustmentFee,
		parsedTx.txInputs,
		inputTotal,
		parsedTx.txRepoTxInputs,
		txOutputs,
		&addrsPrevs,
		paymentRequestIds)
}

func (w *Wallet) createPaymentOutputs(userPayments []UserPayment, changeAddr string, changeAmount btcutil.Amount) map[btcutil.Address]btcutil.Amount {
	var (
		txOutputs = map[btcutil.Address]btcutil.Amount{}
		//if key of map is btcutil.Address which is interface type, uniqueness can't be found from map key
		// so this key is string
		tmpOutputs = map[string]btcutil.Amount{}
	)

	// create txOutput from userPayment
	for _, userPayment := range userPayments {
		if _, ok := tmpOutputs[userPayment.receiverAddr]; ok {
			// if there are multiple receiver same address in user_payment,
			//  sum these
			tmpOutputs[userPayment.receiverAddr] += userPayment.validAmount
		} else {
			// new receiver address
			tmpOutputs[userPayment.receiverAddr] = userPayment.validAmount
		}
	}

	// add txOutput as change address and amount for change
	//TODO:
	// - what if user register for address which is same to payment address?
	//   Though it's impossible in real but systematically possible
	// - BIP44, hdwallet has `ChangeType`. ideally this address should be used
	if _, ok := tmpOutputs[changeAddr]; ok {
		// in real, this is impossible
		tmpOutputs[changeAddr] += changeAmount
	} else {
		// of course, change address is new in txOutputs
		tmpOutputs[changeAddr] = changeAmount
	}

	// create txOutputs from tmpOutputs switching string address type to btcutil.Address
	for strAddr, amount := range tmpOutputs {
		addr, err := w.btc.DecodeAddress(strAddr)
		if err != nil {
			// this case is impossible because addresses are checked in advance
			w.logger.Error("fail to call DecodeAddress",
				zap.String("address", strAddr))
			continue
		}
		txOutputs[addr] = amount
	}

	//Debug
	//grok.Value(tmpOutputs)
	grok.Value(txOutputs)

	return txOutputs
}

//TODO: logic is similar with createRawTx()
func (w *Wallet) createRawPaymentTx(
	adjustmentFee float64,
	txInputs []btcjson.TransactionInput,
	inputTotal btcutil.Amount,
	txPaymentInputs []walletrepo.TxInput,
	txPrevOutputs map[btcutil.Address]btcutil.Amount,
	addrsPrevs *btc.AddrsPrevTxs,
	paymentRequestIds []int64) (string, string, error) {

	// 1. get unallocated address for receiver (done already)

	// 2. create raw transaction as temporary use
	msgTx, err := w.btc.CreateRawTransactionWithOutput(txInputs, txPrevOutputs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.CreateRawTransactionWithOutput()")
	}

	// 3. calculate fee and output total
	// - receipt/transfer
	//  - adjust outputTotal by fee and re-run CreateRawTransaction
	outputTotal, fee, txOutputs, txPaymentOutputs, err := w.calculatePaymentOutputTotal(msgTx, adjustmentFee, inputTotal, txPrevOutputs)
	if err != nil {
		return "", "", err
	}
	w.logger.Debug(
		"total coin to send (Satoshi) after fee calculated",
		zap.Any("output_amount", outputTotal),
		zap.Int("len(inputs)", len(txInputs)))

	// 4. re call CreateRawTransactionWithOutput
	msgTx, err = w.btc.CreateRawTransactionWithOutput(txInputs, txOutputs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.CreateRawTransactionWithOutput()")
	}

	// 5. convert msgTx to hex
	hex, err := w.btc.ToHex(msgTx)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.ToHex(msgTx)")
	}

	// 6. insert to tx_table for unsigned tx
	//  - txReceiptID would be 0 if record is already existing then csv file is not created
	txReceiptID, err := w.insertTxTableForUnsigned(
		action.ActionTypePayment,
		hex, inputTotal,
		outputTotal,
		fee,
		tx.TxTypeValue[tx.TxTypeUnsigned],
		txPaymentInputs,
		txPaymentOutputs,
		paymentRequestIds)

	if err != nil {
		return "", "", errors.Wrap(err, "fail to call insertTxTableForUnsigned()")
	}

	// 7. serialize previous txs for multisig signature
	encodedAddrsPrevs, err := serial.EncodeToString(*addrsPrevs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call serial.EncodeToString()")
	}
	w.logger.Debug(
		"encodedAddrsPrevs",
		zap.String("encodedAddrsPrevs", encodedAddrsPrevs))

	// 8. generate tx file
	var generatedFileName string
	if txReceiptID != 0 {
		generatedFileName, err = w.generateHexFile(action.ActionTypePayment, hex, encodedAddrsPrevs, txReceiptID)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call generateHexFile()")
		}
	}

	return hex, generatedFileName, nil
}

func (w *Wallet) calculatePaymentOutputTotal(
	msgTx *wire.MsgTx,
	adjustmentFee float64,
	inputTotal btcutil.Amount,
	txOutputs map[btcutil.Address]btcutil.Amount,
) (btcutil.Amount, btcutil.Amount, map[btcutil.Address]btcutil.Amount, []walletrepo.TxOutput, error) {

	fee, err := w.btc.GetFee(msgTx, adjustmentFee)
	if err != nil {
		return 0, 0, nil, nil, errors.Wrap(err, "fail to call btc.GetFee()")
	}

	var outputTotal btcutil.Amount
	txPaymentOutputs := make([]walletrepo.TxOutput, 0, len(txOutputs))

	// subtract fee from output transaction for change
	// FIXME: what if change is short, should re-run form the beginning with shortage-flag
	for addr, amt := range txOutputs {
		if acnt, _ := w.btc.GetAccount(addr.String()); acnt == string(account.AccountTypePayment) {
			//chang address
			txOutputs[addr] -= fee
			txPaymentOutputs = append(txPaymentOutputs, walletrepo.TxOutput{
				ReceiptID:     0,
				OutputAddress: addr.String(),
				OutputAccount: string(account.AccountTypePayment),
				OutputAmount:  w.btc.AmountString(amt - fee),
				IsChange:      true,
			})
		} else {
			txPaymentOutputs = append(txPaymentOutputs, walletrepo.TxOutput{
				ReceiptID:     0,
				OutputAddress: addr.String(),
				OutputAccount: "",
				OutputAmount:  w.btc.AmountString(amt),
				IsChange:      false,
			})
		}
		outputTotal += amt
	}
	//total
	outputTotal -= fee
	if outputTotal <= 0 {
		w.logger.Debug(
			"inputTotal is short of coin to pay fee",
			zap.Any("amount", inputTotal),
			zap.Any("len(inputs)", fee))
		return 0, 0, nil, nil, errors.Wrapf(err, "inputTotal is short of coin to pay fee")
	}
	return outputTotal, fee, txOutputs, txPaymentOutputs, nil
}
