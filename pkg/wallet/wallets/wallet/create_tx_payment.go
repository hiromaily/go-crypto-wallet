package wallet

import (
	"github.com/btcsuite/btcd/btcjson"
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
// - sender account(payment) covers fee, but is should be flexible
func (w *Wallet) CreatePaymentTx(adjustmentFee float64) (string, string, error) {
	sender := account.AccountTypePayment
	//receiver := account.AccountTypeAnonymous
	targetAction := action.ActionTypePayment

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
	balance, err := w.btc.GetBalanceByAccount(account.AccountTypePayment)
	//balance, err := w.btc.GetReceivedByLabelAndMinConf(account.AccountTypePayment.String(), w.btc.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= userTotal {
		//balance is short
		return "", "", errors.New("balance for payment account is insufficient")
	}
	w.logger.Debug("payment balane and userTotal",
		zap.Any("balance", balance),
		zap.Any("userTotal", userTotal))

	//FIXME: how to commonalize code from here
	// create transfer transaction
	//return w.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee)

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
	w.logger.Debug(
		"amount",
		zap.Any("requiredAmount", 0),
		zap.Any("input_total", inputTotal),
		zap.Int("len(inputs)", len(parsedTx.txInputs)),
	)

	addrsPrevs := btc.AddrsPrevTxs{
		Addrs:         parsedTx.addresses,
		PrevTxs:       parsedTx.prevTxs,
		SenderAccount: sender,
	}

	// 3. create payment txOutputs from userPayment
	// get sender address for change
	//pubkeyTable, err := w.storager.GetOneUnAllocatedAccountPubKeyTable(sender)
	//if err != nil {
	//	return "", "", errors.Wrap(err, "fail to call storager.GetOneUnAllocatedAccountPubKeyTable()")
	//}

	//changeAddr := pubkeyTable.WalletAddress
	//FIXME: changAddr must be wrong...
	changeAddr := unspentAddrs[0] //this is actually sender's address because it's for change
	changeAmount := inputTotal - userTotal
	w.logger.Debug("before createPaymentOutputs()",
		zap.Any("change_addr", changeAddr),
		zap.Any("change_amount", changeAmount))

	txPrevOutputs := w.createPaymentOutputs(userPayments, changeAddr, changeAmount)
	w.logger.Debug("txPrevOutputs",
		zap.Int("len(txPrevOutputs)", len(txPrevOutputs)))

	//debug
	for addr, amt := range txPrevOutputs {
		w.logger.Debug("txPrevOutputs",
			zap.String("address_string", addr.String()),
			zap.String("address_encoded", addr.EncodeAddress()),
			zap.String("amount", amt.String()),
		)
	}

	// create raw tx
	return w.createPaymentRawTx(
		targetAction,
		sender,
		"",
		adjustmentFee,
		parsedTx.txInputs,
		inputTotal,
		parsedTx.txRepoTxInputs,
		txPrevOutputs,
		&addrsPrevs,
		paymentRequestIds)
}

// userPayments is given for receiverAddr
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

	return txOutputs
}

//TODO: logic is similar with createRawTx()
func (w *Wallet) createPaymentRawTx(
	targetAction action.ActionType,
	sender account.AccountType,
	receiver account.AccountType,
	adjustmentFee float64,
	txInputs []btcjson.TransactionInput,
	inputTotal btcutil.Amount,
	txRepoTxInputs []walletrepo.TxInput,
	txPrevOutputs map[btcutil.Address]btcutil.Amount,
	addrsPrevs *btc.AddrsPrevTxs,
	paymentRequestIds []int64) (string, string, error) {

	// 1. create raw transaction as temporary use
	msgTx, err := w.btc.CreateRawTransaction(txInputs, txPrevOutputs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.CreateRawTransactionWithOutput()")
	}

	// 2. calculate fee and output total
	//  - adjust outputTotal by fee and re-run CreateRawTransaction
	outputTotal, fee, txOutputs, txRepoTxOutputs, err := w.calculateOutputTotal(sender, receiver, msgTx, adjustmentFee, inputTotal, txPrevOutputs)
	if err != nil {
		return "", "", err
	}
	w.logger.Debug("txOutputs", zap.Int("len(txOutputs)", len(txOutputs)))
	w.logger.Debug("txRepoTxOutputs", zap.Int("len(txRepoTxOutputs)", len(txRepoTxOutputs)))

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
			zap.String("output_amount", v.OutputAmount),
		)
	}

	// 3. re call CreateRawTransactionWithOutput
	msgTx, err = w.btc.CreateRawTransaction(txInputs, txOutputs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.CreateRawTransactionWithOutput()")
	}

	// 4. convert msgTx to hex
	hex, err := w.btc.ToHex(msgTx)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.ToHex(msgTx)")
	}

	// 5. insert to tx_table for unsigned tx
	//  - txReceiptID would be 0 if record is already existing then csv file is not created
	txReceiptID, err := w.insertTxTableForUnsigned(
		targetAction,
		hex, inputTotal,
		outputTotal,
		fee,
		tx.TxTypeValue[tx.TxTypeUnsigned],
		txRepoTxInputs,
		txRepoTxOutputs,
		paymentRequestIds)

	if err != nil {
		return "", "", errors.Wrap(err, "fail to call insertTxTableForUnsigned()")
	}

	// 6. serialize previous txs for multisig signature
	encodedAddrsPrevs, err := serial.EncodeToString(*addrsPrevs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call serial.EncodeToString()")
	}
	w.logger.Debug(
		"encodedAddrsPrevs",
		zap.String("encodedAddrsPrevs", encodedAddrsPrevs))

	// 7. generate tx file
	var generatedFileName string
	if txReceiptID != 0 {
		generatedFileName, err = w.generateHexFile(targetAction, hex, encodedAddrsPrevs, txReceiptID)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call generateHexFile()")
		}
	}

	return hex, generatedFileName, nil
}

//func (w *Wallet) calculatePaymentOutputTotal(
//	sender account.AccountType,
//	msgTx *wire.MsgTx,
//	adjustmentFee float64,
//	inputTotal btcutil.Amount,
//	txPrevOutputs map[btcutil.Address]btcutil.Amount,
//) (btcutil.Amount, btcutil.Amount, map[btcutil.Address]btcutil.Amount, []walletrepo.TxOutput, error) {
//	//get fee
//	fee, err := w.btc.GetFee(msgTx, adjustmentFee)
//	if err != nil {
//		return 0, 0, nil, nil, errors.Wrap(err, "fail to call btc.GetFee()")
//	}
//
//	var outputTotal btcutil.Amount
//	txRepoOutputs := make([]walletrepo.TxOutput, 0, len(txPrevOutputs))
//
//	// subtract fee from output transaction for change
//	// FIXME: what if change is short, should re-run form the beginning with shortage-flag
//	for addr, amt := range txPrevOutputs {
//		if acnt, _ := w.btc.GetAccount(addr.String()); acnt == sender.String() {
//			//chang address
//			txPrevOutputs[addr] -= fee
//			txRepoOutputs = append(txRepoOutputs, walletrepo.TxOutput{
//				ReceiptID:     0,
//				OutputAddress: addr.String(),
//				OutputAccount: sender.String(),
//				OutputAmount:  w.btc.AmountString(amt - fee),
//				IsChange:      true,
//			})
//		} else {
//			txRepoOutputs = append(txRepoOutputs, walletrepo.TxOutput{
//				ReceiptID:     0,
//				OutputAddress: addr.String(),
//				OutputAccount: "",
//				OutputAmount:  w.btc.AmountString(amt),
//				IsChange:      false,
//			})
//		}
//		outputTotal += amt
//	}
//
//	//total
//	outputTotal -= fee
//	if outputTotal <= 0 {
//		w.logger.Debug(
//			"inputTotal is short of coin to pay fee",
//			zap.Any("amount", inputTotal),
//			zap.Any("len(inputs)", fee))
//		return 0, 0, nil, nil, errors.Wrapf(err, "inputTotal is short of coin to pay fee")
//	}
//	return outputTotal, fee, txPrevOutputs, txRepoOutputs, nil
//}
