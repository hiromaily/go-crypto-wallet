package watchsrv

import (
	"database/sql"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp/btc"
)

// create_tx.go is for common func among create transaction functionalities

// - check unspentTx for sender account
// - get utxo and create unsigned tx
// - fee is fluctuating, so outdated unsigned tx must not be used, re-create unsigned tx
// - after signed tx sent, utxo could not be retrieved by listUnspent() anymore

// create unsigned tx
// - [actionType:deposit] sender account: client,  receiver account: deposit
// - [actionType:payment] sender account: payment, receiver account: anonymous
// - [actionType:transfer] sender account: it depends on params, receiver account: it depends on params
// value of given requiredAmount
// - depositAction:  0. total amount in clients to deposit
// - transferAction: if 0, total amount of sender is sent to receiver
//                   if not 0, amount is sent from sender to receiver
// - paymentAction: 0. total amount in payment users

// TxCreate type
type TxCreate struct {
	btc             btcgrp.Bitcoiner
	logger          *zap.Logger
	dbConn          *sql.DB
	addrRepo        watchrepo.AddressRepositorier
	txRepo          watchrepo.BTCTxRepositorier
	txInputRepo     watchrepo.TxInputRepositorier
	txOutputRepo    watchrepo.TxOutputRepositorier
	payReqRepo      watchrepo.PaymentRequestRepositorier
	txFileRepo      tx.FileRepositorier
	depositReceiver account.AccountType
	paymentSender   account.AccountType
	wtype           wallet.WalletType
}

// NewTxCreate returns TxCreate object
func NewTxCreate(
	btc btcgrp.Bitcoiner,
	logger *zap.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.BTCTxRepositorier,
	txInputRepo watchrepo.TxInputRepositorier,
	txOutputRepo watchrepo.TxOutputRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	txFileRepo tx.FileRepositorier,
	depositReceiver account.AccountType,
	paymentSender account.AccountType,
	wtype wallet.WalletType,
) *TxCreate {
	return &TxCreate{
		btc:             btc,
		logger:          logger,
		dbConn:          dbConn,
		addrRepo:        addrRepo,
		txRepo:          txRepo,
		txInputRepo:     txInputRepo,
		txOutputRepo:    txOutputRepo,
		payReqRepo:      payReqRepo,
		txFileRepo:      txFileRepo,
		depositReceiver: depositReceiver,
		paymentSender:   paymentSender,
		wtype:           wtype,
	}
}

type parsedTx struct {
	txInputs       []btcjson.TransactionInput
	txRepoTxInputs []*models.BTCTXInput
	prevTxs        []btc.PrevTx
	addresses      []string // input, sender's address
}

// FIXME: receiver account covers fee, but is should be flexible
// TODO: what if `listtransactions` api is called to see result after this func
func (t *TxCreate) createTx(
	sender,
	receiver account.AccountType,
	targetAction action.ActionType,
	requiredAmount btcutil.Amount,
	adjustmentFee float64,
	paymentRequestIds []int64,
	userPayments []UserPayment,
) (string, string, error) {
	t.logger.Debug("createTx()",
		zap.String("sender_acount", sender.String()),
		zap.String("receiver_acount", receiver.String()),
		zap.String("target_action", targetAction.String()),
		zap.Any("required_amount", requiredAmount),
		zap.Float64("adjustmentFee", adjustmentFee))

	// get listUnspent
	unspentList, unspentAddrs, err := t.getUnspentList(sender)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call getUnspentList()")
	}
	if len(unspentList) == 0 {
		t.logger.Info("no listunspent")
		return "", "", nil
	}

	// parse listUnspent
	parsedTx, inputTotal, isDone := t.parseListUnspentTx(unspentList, requiredAmount)
	if len(parsedTx.txInputs) == 0 {
		t.logger.Info("no input tx in listUnspent")
		return "", "", nil
	}
	if !isDone {
		return "", "", errors.New("sender account can't meet amount to send")
	}
	if requiredAmount != 0 {
		t.logger.Debug("amount", zap.Any("expected_change", inputTotal-requiredAmount))
	}

	// create txOutputs
	var txPrevOutputs map[btcutil.Address]btcutil.Amount
	switch targetAction {
	case action.ActionTypeDeposit, action.ActionTypeTransfer:
		var isChange bool
		if requiredAmount != 0 {
			isChange = true
		}
		txPrevOutputs, err = t.createTxOutputs(receiver, requiredAmount, inputTotal, unspentAddrs[0], isChange)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call createTxOutputs()")
		}
	case action.ActionTypePayment:
		changeAddr := unspentAddrs[0] // this is actually sender's address because it's for change
		changeAmount := inputTotal - requiredAmount
		txPrevOutputs = t.createPaymentTxOutputs(userPayments, changeAddr, changeAmount)
		t.logger.Debug("before createPaymentOutputs()",
			zap.Any("change_addr", changeAddr),
			zap.Any("change_amount", changeAmount),
			zap.Int("len(txPrevOutputs)", len(txPrevOutputs)),
		)
	default:
		return "", "", errors.Errorf("invalid actionType: %s", targetAction)
	}

	// create raw transaction as temporary use
	//  - later calculate by tx size
	msgTx, err := t.btc.CreateRawTransaction(parsedTx.txInputs, txPrevOutputs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.CreateRawTransaction()")
	}

	// calculate fee and output total
	//  - adjust outputTotal by fee and re-run CreateRawTransaction
	//  - this logic would be different from payment
	outputTotal, fee, txOutputs, txRepoTxOutputs, err := t.calculateOutputTotal(sender, receiver, msgTx, adjustmentFee, inputTotal, txPrevOutputs)
	if err != nil {
		return "", "", err
	}

	// for debug
	//for addr, amt := range txOutputs {
	//	t.logger.Debug("txOutputs",
	//		zap.String("address_string", addr.String()),
	//		zap.String("address_encoded", addr.EncodeAddress()),
	//		zap.String("amount", amt.String()),
	//	)
	//}
	//for _, v := range txRepoTxOutputs {
	//	t.logger.Debug("txRepoTxOutputs",
	//		zap.String("output_account", v.OutputAccount),
	//		zap.String("output_address", v.OutputAddress),
	//		zap.String("output_amount", v.OutputAmount.String()),
	//	)
	//}

	// re call CreateRawTransaction
	msgTx, err = t.btc.CreateRawTransaction(parsedTx.txInputs, txOutputs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.CreateRawTransaction()")
	}

	// convert msgTx to hex
	hex, err := t.btc.ToHex(msgTx)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call btc.ToHex(msgTx)")
	}

	// insert to tx_table for unsigned tx
	//  - txID would be 0 if record is already existing then csv file is not created
	txID, err := t.insertTxTableForUnsigned(
		targetAction,
		hex,
		inputTotal,
		outputTotal,
		fee,
		parsedTx.txRepoTxInputs,
		txRepoTxOutputs,
		paymentRequestIds)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call insertTxTableForUnsigned()")
	}

	// serialize previous txs for multisig signature
	previousTxs := btc.PreviousTxs{
		SenderAccount: sender,
		PrevTxs:       parsedTx.prevTxs,
		Addrs:         parsedTx.addresses,
	}
	encodedAddrsPrevs, err := serial.EncodeToString(previousTxs)
	if err != nil {
		return "", "", errors.Wrap(err, "fail to call serial.EncodeToString()")
	}

	// generate tx file
	// TODO: how to recover when error occurred here
	// - inserted data in database must be deleted to generate hex file
	var generatedFileName string
	if txID != 0 {
		generatedFileName, err = t.generateHexFile(targetAction, hex, encodedAddrsPrevs, txID)
		if err != nil {
			return "", "", errors.Wrap(err, "fail to call generateHexFile()")
		}
	}

	t.logger.Debug("getUnspentList()",
		zap.Any("unspentList", unspentList),
		zap.Any("unspentAddrs", unspentAddrs),
		zap.Any("requiredAmount", requiredAmount),
		zap.Any("input_total", inputTotal),
		zap.Int("len(inputs)", len(parsedTx.txInputs)),
		zap.Int("len(txPrevOutputs)", len(txPrevOutputs)),
		zap.Int("len(txOutputs)", len(txOutputs)),
		zap.Int("len(txRepoTxOutputs)", len(txRepoTxOutputs)),
		zap.String("encodedAddrsPrevs", encodedAddrsPrevs),
	)

	return hex, generatedFileName, nil
}

// call API `unspentlist`
// this func returns no result, no error possibly, so caller should check both returned value
func (t *TxCreate) getUnspentList(accountType account.AccountType) ([]btc.ListUnspentResult, []string, error) {
	// unlock locked UnspentTransaction
	//if err := w.BTC.UnlockUnspent(); err != nil {
	//	return "", "", err
	//}

	// get listUnspent
	unspentList, err := t.btc.ListUnspentByAccount(accountType, t.btc.ConfirmationBlock())
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to call btc.ListUnspentByAccount()")
	}
	unspentAddrs := t.btc.GetUnspentListAddrs(unspentList, accountType)

	return unspentList, unspentAddrs, nil
}

// parse result of listUnspent
// retured *parsedTx could be nil
func (t *TxCreate) parseListUnspentTx(unspentList []btc.ListUnspentResult, amount btcutil.Amount) (*parsedTx, btcutil.Amount, bool) {
	var inputTotal btcutil.Amount
	txInputs := make([]btcjson.TransactionInput, 0, len(unspentList))
	txRepoTxInputs := make([]*models.BTCTXInput, 0, len(unspentList))
	prevTxs := make([]btc.PrevTx, 0, len(unspentList))
	addresses := make([]string, 0, len(unspentList))

	var isDone bool // if isDone is false, sender can't meet amount
	if amount == 0 {
		isDone = true
	}

	for _, tx := range unspentList {
		// Amount
		amt, err := btcutil.NewAmount(tx.Amount)
		if err != nil {
			// this error is not expected
			t.logger.Error(
				"fail to call btcutil.NewAmount() then skipped",
				zap.String("tx_id", tx.TxID),
				zap.Float64("tx_amount", tx.Amount),
				zap.Error(err))
			continue
		}
		inputTotal += amt

		txInputs = append(txInputs, btcjson.TransactionInput{
			Txid: tx.TxID,
			Vout: tx.Vout,
		})

		txRepoTxInputs = append(txRepoTxInputs, &models.BTCTXInput{
			TXID:               0,
			InputTxid:          tx.TxID,
			InputVout:          tx.Vout,
			InputAddress:       tx.Address,
			InputAccount:       tx.Label,
			InputAmount:        t.btc.FloatToDecimal(tx.Amount),
			InputConfirmations: uint64(tx.Confirmations),
		})

		// TODO: if sender is client account (non-multisig address), RedeemScript is blank
		prevTxs = append(prevTxs, btc.PrevTx{
			Txid:         tx.TxID,
			Vout:         tx.Vout,
			ScriptPubKey: tx.ScriptPubKey,
			RedeemScript: tx.RedeemScript, // required if target account is multsig address
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

// for ActionTypeDeposit, ActionTypeTransfer
func (t *TxCreate) createTxOutputs(
	reciver account.AccountType,
	requiredAmount btcutil.Amount,
	inputTotal btcutil.Amount,
	senderAddr string,
	isChange bool,
) (map[btcutil.Address]btcutil.Amount, error) {
	// get unallocated address for receiver
	// - deposit/transfer
	pubkeyTable, err := t.addrRepo.GetOneUnAllocated(reciver)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call pubkeyRepo.GetOneUnAllocated()")
	}
	receiverAddr := pubkeyTable.WalletAddress

	// create receiver txOutput
	receiverDecodedAddr, err := btcutil.DecodeAddress(receiverAddr, t.btc.GetChainConf())
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call btcutil.DecodeAddress(%s)", receiverAddr)
	}
	txPrevOutputs := make(map[btcutil.Address]btcutil.Amount)
	if isChange {
		txPrevOutputs[receiverDecodedAddr] = requiredAmount // satoshi
	} else {
		txPrevOutputs[receiverDecodedAddr] = inputTotal // satoshi
	}
	t.logger.Debug("receiver txOutput",
		zap.String("receiverAddr", receiverAddr),
		zap.Any("receivedAmount", txPrevOutputs[receiverDecodedAddr]))

	// if change is required
	if isChange {
		t.logger.Debug("change is required")
		senderDecodedAddr, err := btcutil.DecodeAddress(senderAddr, t.btc.GetChainConf())
		if err != nil {
			return nil, errors.Wrapf(err, "fail to call btcutil.DecodeAddress(%s)", receiverAddr)
		}

		// for change
		// TODO: what if (inputTotal - requiredAmount) = 0,
		//  fee can not be paid from txOutput for change
		txPrevOutputs[senderDecodedAddr] = inputTotal - requiredAmount

		t.logger.Debug("change(sender) txOutput",
			zap.String("senderAddr", senderAddr),
			zap.Any("inputTotal - requiredAmount", inputTotal-requiredAmount))

	}
	return txPrevOutputs, nil
}

func (t *TxCreate) calculateOutputTotal(
	sender account.AccountType,
	receiver account.AccountType,
	msgTx *wire.MsgTx,
	adjustmentFee float64,
	inputTotal btcutil.Amount,
	txPrevOutputs map[btcutil.Address]btcutil.Amount,
) (btcutil.Amount, btcutil.Amount, map[btcutil.Address]btcutil.Amount, []*models.BTCTXOutput, error) {
	// get fee
	fee, err := t.btc.GetFee(msgTx, adjustmentFee)
	if err != nil {
		return 0, 0, nil, nil, errors.Wrap(err, "fail to call btc.GetFee()")
	}
	var outputTotal btcutil.Amount
	txRepoOutputs := make([]*models.BTCTXOutput, 0, len(txPrevOutputs))

	// subtract fee from output transaction for change
	// FIXME: what if change is short, should re-run form the beginning with shortage-flag
	for addr, amt := range txPrevOutputs {
		if len(txPrevOutputs) == 1 {
			// no change
			txPrevOutputs[addr] -= fee
			txRepoOutputs = append(txRepoOutputs, &models.BTCTXOutput{
				TXID:          0,
				OutputAddress: addr.String(),
				OutputAccount: receiver.String(),
				OutputAmount:  t.btc.AmountToDecimal(amt - fee),
				IsChange:      false,
			})
			outputTotal += amt
			break
		}

		if acnt, _ := t.btc.GetAccount(addr.String()); acnt == sender.String() {
			t.logger.Debug("detect sender account in calculateOutputTotal")
			// address is used for change
			txPrevOutputs[addr] -= fee
			txRepoOutputs = append(txRepoOutputs, &models.BTCTXOutput{
				TXID:          0,
				OutputAddress: addr.String(),
				OutputAccount: sender.String(),
				OutputAmount:  t.btc.AmountToDecimal(amt - fee),
				IsChange:      true,
			})
		} else {
			txRepoOutputs = append(txRepoOutputs, &models.BTCTXOutput{
				TXID:          0,
				OutputAddress: addr.String(),
				OutputAccount: receiver.String(),
				OutputAmount:  t.btc.AmountToDecimal(amt),
				IsChange:      false,
			})
		}
		outputTotal += amt
	}
	t.logger.Debug("calculateOutputTotal",
		zap.Any("fee", fee),
		zap.Any("outputTotal (before fee adjustment)", outputTotal),
		zap.Any("outputTotal by (inputTotal - fee)", inputTotal-fee),
		zap.Any("outputTotal by (outputTotal - fee)", outputTotal-fee),
	)

	// total amount should be same
	outputTotal = inputTotal - fee

	if outputTotal <= 0 {
		t.logger.Debug(
			"inputTotal is short of coin to pay fee",
			zap.Any("amount of inputTotal", inputTotal),
			zap.Any("fee", fee))
		return 0, 0, nil, nil, errors.Wrapf(err, "inputTotal is short of coin to pay fee")
	}

	return outputTotal, fee, txPrevOutputs, txRepoOutputs, nil
}

func (t *TxCreate) insertTxTableForUnsigned(
	actionType action.ActionType,
	hex string,
	inputTotal,
	outputTotal,
	fee btcutil.Amount,
	txInputs []*models.BTCTXInput,
	txOutputs []*models.BTCTXOutput,
	paymentRequestIds []int64,
) (int64, error) {
	// skip if same hex is already stored
	count, err := t.txRepo.GetCountByUnsignedHex(actionType, hex)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call repo.Tx().GetCountByUnsignedHex()")
	}
	if count != 0 {
		// skip
		return 0, nil
	}

	// TxReceipt table
	txItem := &models.BTCTX{
		Action:            actionType.String(),
		UnsignedHexTX:     hex,
		TotalInputAmount:  t.btc.AmountToDecimal(inputTotal),
		TotalOutputAmount: t.btc.AmountToDecimal(outputTotal),
		Fee:               t.btc.AmountToDecimal(fee),
	}

	// start database transaction
	dtx, err := t.dbConn.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "fail to start transaction")
	}
	defer func() {
		if err != nil {
			dtx.Rollback()
		} else {
			dtx.Commit()
		}
	}()

	// tx := w.repo.MustBegin()
	txID, err := t.txRepo.InsertUnsignedTx(actionType, txItem)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call repo.Tx().InsertUnsignedTx()")
	}

	// TxReceiptInput table
	//  update txID
	for idx := range txInputs {
		txInputs[idx].TXID = txID
	}
	err = t.txInputRepo.InsertBulk(txInputs)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call txInRepo.InsertBulk()")
	}

	// TxReceiptOutput table
	//  update txID
	for idx := range txOutputs {
		txOutputs[idx].TXID = txID
	}
	err = t.txOutputRepo.InsertBulk(txOutputs)
	if err != nil {
		return 0, errors.Wrap(err, "fail to call repo.TxOutput().InsertBulk()")
	}

	// update payment_id in payment_request table for only action.ActionTypePayment
	if actionType == action.ActionTypePayment {
		_, err = t.payReqRepo.UpdatePaymentID(txID, paymentRequestIds)
		if err != nil {
			return 0, errors.Wrap(err, "fail to call repo.PayReq().UpdatePaymentID(txID, paymentRequestIds)")
		}
	}

	return txID, nil
}

// generateHexFile generate file for hex and encoded previous addresses
func (t *TxCreate) generateHexFile(actionType action.ActionType, hex, encodedAddrsPrevs string, id int64) (string, error) {
	var (
		generatedFileName string
		err               error
	)

	savedata := hex
	if encodedAddrsPrevs != "" {
		savedata = fmt.Sprintf("%s,%s", savedata, encodedAddrsPrevs)
	}

	// create file
	path := t.txFileRepo.CreateFilePath(actionType, tx.TxTypeUnsigned, id, 0)
	generatedFileName, err = t.txFileRepo.WriteFile(path, savedata)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.WriteFile()")
	}

	return generatedFileName, nil
}

// IsFoundTxIDAndVout finds out txID and vout from related txInputs
// nolint: unused
func (t *TxCreate) IsFoundTxIDAndVout(txID string, vout uint32, inputs []btcjson.TransactionInput) bool {
	for _, val := range inputs {
		if val.Txid == txID && val.Vout == vout {
			return true
		}
	}
	return false
}
