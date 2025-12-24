package watchsrv

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
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
	btc             bitcoin.Bitcoiner
	dbConn          *sql.DB
	addrRepo        watch.AddressRepositorier
	txRepo          watch.BTCTxRepositorier
	txInputRepo     watch.TxInputRepositorier
	txOutputRepo    watch.TxOutputRepositorier
	payReqRepo      watch.PaymentRequestRepositorier
	txFileRepo      file.TransactionFileRepositorier
	depositReceiver domainAccount.AccountType
	paymentSender   domainAccount.AccountType
	wtype           domainWallet.WalletType
}

// NewTxCreate returns TxCreate object
func NewTxCreate(
	btcAPI bitcoin.Bitcoiner,
	dbConn *sql.DB,
	addrRepo watch.AddressRepositorier,
	txRepo watch.BTCTxRepositorier,
	txInputRepo watch.TxInputRepositorier,
	txOutputRepo watch.TxOutputRepositorier,
	payReqRepo watch.PaymentRequestRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	depositReceiver domainAccount.AccountType,
	paymentSender domainAccount.AccountType,
	wtype domainWallet.WalletType,
) *TxCreate {
	return &TxCreate{
		btc:             btcAPI,
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
//
//nolint:gocyclo
func (t *TxCreate) createTx(
	sender,
	receiver domainAccount.AccountType,
	targetAction domainTx.ActionType,
	requiredAmount btcutil.Amount,
	adjustmentFee float64,
	paymentRequestIds []int64,
	userPayments []UserPayment,
) (string, string, error) {
	logger.Debug("createTx()",
		"sender_acount", sender.String(),
		"receiver_acount", receiver.String(),
		"target_action", targetAction.String(),
		"required_amount", requiredAmount,
		"adjustmentFee", adjustmentFee)

	// get listUnspent
	unspentList, unspentAddrs, err := t.getUnspentList(sender)
	if err != nil {
		return "", "", fmt.Errorf("fail to call getUnspentList(): %w", err)
	}
	if len(unspentList) == 0 {
		logger.Info("no listunspent")
		return "", "", nil
	}

	// parse listUnspent
	parsedTx, inputTotal, isDone := t.parseListUnspentTx(unspentList, requiredAmount)
	if len(parsedTx.txInputs) == 0 {
		logger.Info("no input tx in listUnspent")
		return "", "", nil
	}
	if !isDone {
		return "", "", errors.New("sender account can't meet amount to send")
	}
	if requiredAmount != 0 {
		logger.Debug("amount", "expected_change", inputTotal-requiredAmount)
	}

	// create txOutputs
	var txPrevOutputs map[btcutil.Address]btcutil.Amount
	switch targetAction {
	case domainTx.ActionTypeDeposit, domainTx.ActionTypeTransfer:
		var isChange bool
		if requiredAmount != 0 {
			isChange = true
		}
		txPrevOutputs, err = t.createTxOutputs(receiver, requiredAmount, inputTotal, unspentAddrs[0], isChange)
		if err != nil {
			return "", "", fmt.Errorf("fail to call createTxOutputs(): %w", err)
		}
	case domainTx.ActionTypePayment:
		changeAddr := unspentAddrs[0] // this is actually sender's address because it's for change
		changeAmount := inputTotal - requiredAmount
		txPrevOutputs = t.createPaymentTxOutputs(userPayments, changeAddr, changeAmount)
		logger.Debug("before createPaymentOutputs()",
			"change_addr", changeAddr,
			"change_amount", changeAmount,
			"len(txPrevOutputs)", len(txPrevOutputs),
		)
	default:
		return "", "", fmt.Errorf("invalid actionType: %s", targetAction)
	}

	// create raw transaction as temporary use
	//  - later calculate by tx size
	msgTx, err := t.btc.CreateRawTransaction(parsedTx.txInputs, txPrevOutputs)
	if err != nil {
		return "", "", fmt.Errorf("fail to call btc.CreateRawTransaction(): %w", err)
	}

	// calculate fee and output total
	//  - adjust outputTotal by fee and re-run CreateRawTransaction
	//  - this logic would be different from payment
	outputTotal, fee, txOutputs, txRepoTxOutputs, err := t.calculateOutputTotal(
		sender, receiver, msgTx, adjustmentFee, inputTotal, txPrevOutputs)
	if err != nil {
		return "", "", err
	}

	// for debug
	// for addr, amt := range txOutputs {
	//	t.logger.Debug("txOutputs",
	//		"address_string", addr.String(),
	//		"address_encoded", addr.EncodeAddress(),
	//		"amount", amt.String(),
	//	)
	//}
	// for _, v := range txRepoTxOutputs {
	//	t.logger.Debug("txRepoTxOutputs",
	//		"output_account", v.OutputAccount,
	//		"output_address", v.OutputAddress,
	//		"output_amount", v.OutputAmount.String(),
	//	)
	//}

	// re call CreateRawTransaction
	msgTx, err = t.btc.CreateRawTransaction(parsedTx.txInputs, txOutputs)
	if err != nil {
		return "", "", fmt.Errorf("fail to call btc.CreateRawTransaction(): %w", err)
	}

	// convert msgTx to hex
	hex, err := t.btc.ToHex(msgTx)
	if err != nil {
		return "", "", fmt.Errorf("fail to call btc.ToHex(msgTx): %w", err)
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
		return "", "", fmt.Errorf("fail to call insertTxTableForUnsigned(): %w", err)
	}

	// serialize previous txs for multisig signature
	previousTxs := btc.PreviousTxs{
		SenderAccount: sender,
		PrevTxs:       parsedTx.prevTxs,
		Addrs:         parsedTx.addresses,
	}
	encodedAddrsPrevs, err := serial.EncodeToString(previousTxs)
	if err != nil {
		return "", "", fmt.Errorf("fail to call serial.EncodeToString(): %w", err)
	}

	// generate tx file
	// TODO: how to recover when error occurred here
	// - inserted data in database must be deleted to generate hex file
	var generatedFileName string
	if txID != 0 {
		generatedFileName, err = t.generateHexFile(targetAction, hex, encodedAddrsPrevs, txID)
		if err != nil {
			return "", "", fmt.Errorf("fail to call generateHexFile(): %w", err)
		}
	}

	logger.Debug("getUnspentList()",
		"unspentList", unspentList,
		"unspentAddrs", unspentAddrs,
		"requiredAmount", requiredAmount,
		"input_total", inputTotal,
		"len(inputs)", len(parsedTx.txInputs),
		"len(txPrevOutputs)", len(txPrevOutputs),
		"len(txOutputs)", len(txOutputs),
		"len(txRepoTxOutputs)", len(txRepoTxOutputs),
		"encodedAddrsPrevs", encodedAddrsPrevs,
	)

	return hex, generatedFileName, nil
}

// call API `unspentlist`
// this func returns no result, no error possibly, so caller should check both returned value
func (t *TxCreate) getUnspentList(accountType domainAccount.AccountType) ([]btc.ListUnspentResult, []string, error) {
	// unlock locked UnspentTransaction
	// if err := w.BTC.UnlockUnspent(); err != nil {
	//	return "", "", err
	//}

	// get listUnspent
	unspentList, err := t.btc.ListUnspentByAccount(accountType, t.btc.ConfirmationBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call btc.ListUnspentByAccount(): %w", err)
	}
	unspentAddrs := t.btc.GetUnspentListAddrs(unspentList, accountType)

	return unspentList, unspentAddrs, nil
}

// parse result of listUnspent
// retured *parsedTx could be nil
func (t *TxCreate) parseListUnspentTx(
	unspentList []btc.ListUnspentResult, amount btcutil.Amount,
) (*parsedTx, btcutil.Amount, bool) {
	var inputTotal btcutil.Amount
	txInputs := make([]btcjson.TransactionInput, 0, len(unspentList))
	txRepoTxInputs := make([]*models.BTCTXInput, 0, len(unspentList))
	prevTxs := make([]btc.PrevTx, 0, len(unspentList))
	addresses := make([]string, 0, len(unspentList))

	var isDone bool // if isDone is false, sender can't meet amount
	if amount == 0 {
		isDone = true
	}

	for _, txItem := range unspentList {
		// Amount
		amt, err := btcutil.NewAmount(txItem.Amount)
		if err != nil {
			// this error is not expected
			logger.Error(
				"fail to call btcutil.NewAmount() then skipped",
				"tx_id", txItem.TxID,
				"tx_amount", txItem.Amount,
				"error", err)
			continue
		}
		inputTotal += amt

		txInputs = append(txInputs, btcjson.TransactionInput{
			Txid: txItem.TxID,
			Vout: txItem.Vout,
		})

		inputAmount, err := t.btc.FloatToDecimal(txItem.Amount)
		if err != nil {
			logger.Error("fail to convert input amount to decimal", "error", err)
			continue
		}
		txRepoTxInputs = append(txRepoTxInputs, &models.BTCTXInput{
			TXID:               0,
			InputTxid:          txItem.TxID,
			InputVout:          txItem.Vout,
			InputAddress:       txItem.Address,
			InputAccount:       txItem.Label,
			InputAmount:        inputAmount,
			InputConfirmations: uint64(txItem.Confirmations),
		})

		// TODO: if sender is client account (non-multisig address), RedeemScript is blank
		prevTxs = append(prevTxs, btc.PrevTx{
			Txid:         txItem.TxID,
			Vout:         txItem.Vout,
			ScriptPubKey: txItem.ScriptPubKey,
			RedeemScript: txItem.RedeemScript, // required if target account is multsig address
			Amount:       txItem.Amount,
		})

		addresses = append(addresses, txItem.Address)

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
	reciver domainAccount.AccountType,
	requiredAmount btcutil.Amount,
	inputTotal btcutil.Amount,
	senderAddr string,
	isChange bool,
) (map[btcutil.Address]btcutil.Amount, error) {
	// get unallocated address for receiver
	// - deposit/transfer
	pubkeyTable, err := t.addrRepo.GetOneUnAllocated(reciver)
	if err != nil {
		return nil, fmt.Errorf("fail to call pubkeyRepo.GetOneUnAllocated(): %w", err)
	}
	receiverAddr := pubkeyTable.WalletAddress

	// create receiver txOutput
	receiverDecodedAddr, err := btcutil.DecodeAddress(receiverAddr, t.btc.GetChainConf())
	if err != nil {
		return nil, fmt.Errorf("fail to call btcutil.DecodeAddress(%s): %w", receiverAddr, err)
	}
	txPrevOutputs := make(map[btcutil.Address]btcutil.Amount)
	if isChange {
		txPrevOutputs[receiverDecodedAddr] = requiredAmount // satoshi
	} else {
		txPrevOutputs[receiverDecodedAddr] = inputTotal // satoshi
	}
	logger.Debug("receiver txOutput",
		"receiverAddr", receiverAddr,
		"receivedAmount", txPrevOutputs[receiverDecodedAddr])

	// if change is required
	if isChange {
		logger.Debug("change is required")
		senderDecodedAddr, decodeErr := btcutil.DecodeAddress(senderAddr, t.btc.GetChainConf())
		if decodeErr != nil {
			return nil, fmt.Errorf("fail to call btcutil.DecodeAddress(%s): %w", receiverAddr, decodeErr)
		}

		// for change
		// TODO: what if (inputTotal - requiredAmount) = 0,
		//  fee can not be paid from txOutput for change
		txPrevOutputs[senderDecodedAddr] = inputTotal - requiredAmount

		logger.Debug("change(sender) txOutput",
			"senderAddr", senderAddr,
			"inputTotal - requiredAmount", inputTotal-requiredAmount)
	}
	return txPrevOutputs, nil
}

func (t *TxCreate) calculateOutputTotal(
	sender domainAccount.AccountType,
	receiver domainAccount.AccountType,
	msgTx *wire.MsgTx,
	adjustmentFee float64,
	inputTotal btcutil.Amount,
	txPrevOutputs map[btcutil.Address]btcutil.Amount,
) (btcutil.Amount, btcutil.Amount, map[btcutil.Address]btcutil.Amount, []*models.BTCTXOutput, error) {
	// get fee
	fee, err := t.btc.GetFee(msgTx, adjustmentFee)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("fail to call btc.GetFee(): %w", err)
	}
	var outputTotal btcutil.Amount
	txRepoOutputs := make([]*models.BTCTXOutput, 0, len(txPrevOutputs))

	// subtract fee from output transaction for change
	// FIXME: what if change is short, should re-run form the beginning with shortage-flag
	for addr, amt := range txPrevOutputs {
		if len(txPrevOutputs) == 1 {
			// no change
			txPrevOutputs[addr] -= fee
			outputAmount, err := t.btc.AmountToDecimal(amt - fee)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to convert output amount to decimal: %w", err)
			}
			txRepoOutputs = append(txRepoOutputs, &models.BTCTXOutput{
				TXID:          0,
				OutputAddress: addr.String(),
				OutputAccount: receiver.String(),
				OutputAmount:  outputAmount,
				IsChange:      false,
			})
			outputTotal += amt
			break
		}

		if acnt, _ := t.btc.GetAccount(addr.String()); acnt == sender.String() {
			logger.Debug("detect sender account in calculateOutputTotal")
			// address is used for change
			txPrevOutputs[addr] -= fee
			outputAmount, err := t.btc.AmountToDecimal(amt - fee)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to convert change amount to decimal: %w", err)
			}
			txRepoOutputs = append(txRepoOutputs, &models.BTCTXOutput{
				TXID:          0,
				OutputAddress: addr.String(),
				OutputAccount: sender.String(),
				OutputAmount:  outputAmount,
				IsChange:      true,
			})
		} else {
			outputAmount, err := t.btc.AmountToDecimal(amt)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to convert output amount to decimal: %w", err)
			}
			txRepoOutputs = append(txRepoOutputs, &models.BTCTXOutput{
				TXID:          0,
				OutputAddress: addr.String(),
				OutputAccount: receiver.String(),
				OutputAmount:  outputAmount,
				IsChange:      false,
			})
		}
		outputTotal += amt
	}
	logger.Debug("calculateOutputTotal",
		"fee", fee,
		"outputTotal (before fee adjustment)", outputTotal,
		"outputTotal by (inputTotal - fee)", inputTotal-fee,
		"outputTotal by (outputTotal - fee)", outputTotal-fee,
	)

	// total amount should be same
	outputTotal = inputTotal - fee

	if outputTotal <= 0 {
		logger.Debug(
			"inputTotal is short of coin to pay fee",
			"amount of inputTotal", inputTotal,
			"fee", fee)
		return 0, 0, nil, nil, fmt.Errorf("inputTotal is short of coin to pay fee: %w", err)
	}

	return outputTotal, fee, txPrevOutputs, txRepoOutputs, nil
}

func (t *TxCreate) insertTxTableForUnsigned(
	actionType domainTx.ActionType,
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
		return 0, fmt.Errorf("fail to call repo.Tx().GetCountByUnsignedHex(): %w", err)
	}
	if count != 0 {
		// skip
		return 0, nil
	}

	// TxReceipt table
	totalInputAmt, err := t.btc.AmountToDecimal(inputTotal)
	if err != nil {
		return 0, fmt.Errorf("fail to convert total input amount to decimal: %w", err)
	}
	totalOutputAmt, err := t.btc.AmountToDecimal(outputTotal)
	if err != nil {
		return 0, fmt.Errorf("fail to convert total output amount to decimal: %w", err)
	}
	feeAmt, err := t.btc.AmountToDecimal(fee)
	if err != nil {
		return 0, fmt.Errorf("fail to convert fee amount to decimal: %w", err)
	}
	txItem := &models.BTCTX{
		Action:            actionType.String(),
		UnsignedHexTX:     hex,
		TotalInputAmount:  totalInputAmt,
		TotalOutputAmount: totalOutputAmt,
		Fee:               feeAmt,
	}

	// start database transaction
	dtx, err := t.dbConn.Begin()
	if err != nil {
		return 0, fmt.Errorf("fail to start transaction: %w", err)
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
		return 0, fmt.Errorf("fail to call repo.Tx().InsertUnsignedTx(): %w", err)
	}

	// TxReceiptInput table
	//  update txID
	for idx := range txInputs {
		txInputs[idx].TXID = txID
	}
	err = t.txInputRepo.InsertBulk(txInputs)
	if err != nil {
		return 0, fmt.Errorf("fail to call txInRepo.InsertBulk(): %w", err)
	}

	// TxReceiptOutput table
	//  update txID
	for idx := range txOutputs {
		txOutputs[idx].TXID = txID
	}
	err = t.txOutputRepo.InsertBulk(txOutputs)
	if err != nil {
		return 0, fmt.Errorf("fail to call repo.TxOutput().InsertBulk(): %w", err)
	}

	// update payment_id in payment_request table for only domainTx.ActionTypePayment
	if actionType == domainTx.ActionTypePayment {
		_, err = t.payReqRepo.UpdatePaymentID(txID, paymentRequestIds)
		if err != nil {
			return 0, fmt.Errorf("fail to call repo.PayReq().UpdatePaymentID(txID, paymentRequestIds): %w", err)
		}
	}

	return txID, nil
}

// generateHexFile generate file for hex and encoded previous addresses
func (t *TxCreate) generateHexFile(
	actionType domainTx.ActionType, hex, encodedAddrsPrevs string, id int64,
) (string, error) {
	var (
		generatedFileName string
		err               error
	)

	savedata := hex
	if encodedAddrsPrevs != "" {
		savedata = fmt.Sprintf("%s,%s", savedata, encodedAddrsPrevs)
	}

	// create file
	path := t.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeUnsigned, id, 0)
	generatedFileName, err = t.txFileRepo.WriteFile(path, savedata)
	if err != nil {
		return "", fmt.Errorf("fail to call txFileRepo.WriteFile(): %w", err)
	}

	return generatedFileName, nil
}

// IsFoundTxIDAndVout finds out txID and vout from related txInputs
func (*TxCreate) IsFoundTxIDAndVout(txID string, vout uint32, inputs []btcjson.TransactionInput) bool {
	for _, val := range inputs {
		if val.Txid == txID && val.Vout == vout {
			return true
		}
	}
	return false
}
