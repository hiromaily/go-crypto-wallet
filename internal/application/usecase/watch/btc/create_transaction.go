// Package btc provides BTC-specific transaction creation use cases.
//
// WARNING: This file is for BTC (Bitcoin) ONLY.
// If cointype is BCH (Bitcoin Cash), DO NOT modify this file.
// Use the BCH-specific implementation instead:
//
//	internal/application/usecase/watch/bch/create_transaction.go
//
// Key differences:
//   - BTC: Uses PSBT (BIP174) format, outputs .psbt files
//   - BCH: Uses raw transaction hex format, outputs .hex files
package btc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// createTxBTCClient defines the minimal interface needed for BTC transaction creation.
// This follows the Interface Segregation Principle - depend only on methods actually used.
type createTxBTCClient interface {
	apibtc.ChainConfigProvider
	apibtc.AmountConverter
	apibtc.UTXOProvider
	apibtc.RawTransactionCreator
	apibtc.AddressOperator
	apibtc.BalanceChecker
	apibtc.PSBTCreator
}

type createTransactionUseCase struct {
	btcClient       createTxBTCClient
	dbConn          *sql.DB
	addrRepo        repowatch.AddressRepositorier
	txRepo          repowatch.BTCTxRepositorier
	txInputRepo     repowatch.TxInputRepositorier
	txOutputRepo    repowatch.TxOutputRepositorier
	payReqRepo      repowatch.PaymentRequestRepositorier
	txFileRepo      file.TransactionFileRepositorier
	depositReceiver domainAccount.AccountType
	paymentSender   domainAccount.AccountType
	walletType      domainWallet.WalletType
}

// NewCreateTransactionUseCase creates a new CreateTransactionUseCase
func NewCreateTransactionUseCase(
	btcClient createTxBTCClient,
	dbConn *sql.DB,
	addrRepo repowatch.AddressRepositorier,
	txRepo repowatch.BTCTxRepositorier,
	txInputRepo repowatch.TxInputRepositorier,
	txOutputRepo repowatch.TxOutputRepositorier,
	payReqRepo repowatch.PaymentRequestRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	depositReceiver domainAccount.AccountType,
	paymentSender domainAccount.AccountType,
	walletType domainWallet.WalletType,
) watchusecase.CreateTransactionUseCase {
	return &createTransactionUseCase{
		btcClient:       btcClient,
		dbConn:          dbConn,
		addrRepo:        addrRepo,
		txRepo:          txRepo,
		txInputRepo:     txInputRepo,
		txOutputRepo:    txOutputRepo,
		payReqRepo:      payReqRepo,
		txFileRepo:      txFileRepo,
		depositReceiver: depositReceiver,
		paymentSender:   paymentSender,
		walletType:      walletType,
	}
}

func (u *createTransactionUseCase) Execute(
	ctx context.Context,
	input watchusecase.CreateTransactionInput,
) (watchusecase.CreateTransactionOutput, error) {
	// Convert action type string to domain type
	actionType := domainTx.ActionType(input.ActionType)
	if !domainTx.ValidateActionType(input.ActionType) {
		return watchusecase.CreateTransactionOutput{}, fmt.Errorf("invalid action type: %s", input.ActionType)
	}

	var hex, fileName string
	var execErr error

	switch actionType {
	case domainTx.ActionTypeDeposit:
		hex, fileName, execErr = u.createDepositTx(input.AdjustmentFee)
	case domainTx.ActionTypePayment:
		hex, fileName, execErr = u.createPaymentTx(input.AdjustmentFee)
	case domainTx.ActionTypeTransfer:
		hex, fileName, execErr = u.createTransferTx(
			input.SenderAccount,
			input.ReceiverAccount,
			input.Amount,
			input.AdjustmentFee,
		)
	default:
		return watchusecase.CreateTransactionOutput{},
			fmt.Errorf("unsupported action type: %s", input.ActionType)
	}

	if execErr != nil {
		return watchusecase.CreateTransactionOutput{}, fmt.Errorf("failed to create transaction: %w", execErr)
	}

	return watchusecase.CreateTransactionOutput{
		TransactionHex: hex,
		FileName:       fileName,
	}, nil
}

// createDepositTx creates unsigned tx if client accounts have coins
// - sender: client, receiver: deposit
// - receiver account covers fee, but should be flexible
func (u *createTransactionUseCase) createDepositTx(adjustmentFee float64) (string, string, error) {
	sender := domainAccount.AccountTypeClient
	receiver := u.depositReceiver
	targetAction := domainTx.ActionTypeDeposit
	logger.Debug("account",
		"sender", sender.String(),
		"receiver", receiver.String(),
	)

	requiredAmount, err := u.btcClient.FloatToAmount(0)
	if err != nil {
		return "", "", err
	}

	// create deposit transaction
	return u.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, nil, nil)
}

// createPaymentTx creates unsigned tx for user (anonymous addresses)
// sender: payment, receiver: addresses coming from payment_request table
// - sender account (payment) covers fee, but should be flexible
func (u *createTransactionUseCase) createPaymentTx(adjustmentFee float64) (string, string, error) {
	sender := u.paymentSender
	receiver := domainAccount.AccountTypeAnonymous
	targetAction := domainTx.ActionTypePayment
	logger.Debug("account",
		"sender", sender.String(),
		"receiver", receiver.String(),
	)

	// get payment data from payment_request
	userPayments, paymentRequestIds, err := u.createUserPayment()
	if err != nil {
		return "", "", err
	}
	if len(userPayments) == 0 {
		logger.Debug("no data in userPayments")
		// no data
		return "", "", nil
	}

	// calculate total amount to send from payment_request
	var requiredAmount btcutil.Amount
	for _, val := range userPayments {
		requiredAmount += val.validAmount
	}

	// get balance for payment account
	balance, err := u.btcClient.GetBalanceByAccount(domainAccount.AccountTypePayment, u.btcClient.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= requiredAmount {
		// balance is short
		logger.Info("balance for payment account is insufficient",
			"payment_balance", balance.ToBTC(),
			"required_amount", requiredAmount.ToBTC(),
		)
		return "", "", nil
	}
	logger.Debug("payment balance and userTotal",
		"balance", balance,
		"userTotal", requiredAmount)

	// create payment transaction
	return u.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, paymentRequestIds, userPayments)
}

// createTransferTx creates unsigned tx for transfer coin among internal accounts except client, authorization
// FIXME: for now, receiver account covers fee, but should be flexible
func (u *createTransactionUseCase) createTransferTx(
	sender, receiver domainAccount.AccountType, floatAmount, adjustmentFee float64,
) (string, string, error) {
	targetAction := domainTx.ActionTypeTransfer

	// validation account
	if receiver == domainAccount.AccountTypeClient || receiver == domainAccount.AccountTypeAuthorization {
		return "", "", errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return "", "", errors.New("invalid account. sender and receiver is same")
	}

	// amount btcutil.Amount
	requiredAmount, err := u.btcClient.FloatToAmount(floatAmount)
	if err != nil {
		return "", "", err
	}

	// check balance for sender
	balance, err := u.btcClient.GetBalanceByAccount(sender, u.btcClient.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= requiredAmount {
		// balance is short
		return "", "", fmt.Errorf("account: %s balance is insufficient", sender)
	}

	// create transfer transaction
	return u.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, nil, nil)
}

type parsedTx struct {
	txInputs       []btcjson.TransactionInput
	txRepoTxInputs []*domainBitcoin.BTCTxInput
	prevTxs        []dtobtc.PreviousTx
	addresses      []string // input, sender's address
}

// userPayment represents user's payment address and amount
type userPayment struct {
	senderAddr   string          // sender address for just checking
	receiverAddr string          // receiver address
	validRecAddr btcutil.Address // decoded receiver address
	amount       float64         // amount
	validAmount  btcutil.Amount  // decoded amount
}

// FIXME: receiver account covers fee, but should be flexible
// TODO: what if `listtransactions` api is called to see result after this func
//
//nolint:gocyclo
func (u *createTransactionUseCase) createTx(
	sender,
	receiver domainAccount.AccountType,
	targetAction domainTx.ActionType,
	requiredAmount btcutil.Amount,
	adjustmentFee float64,
	paymentRequestIds []int64,
	userPayments []userPayment,
) (string, string, error) {
	logger.Debug("createTx()",
		"sender_account", sender.String(),
		"receiver_account", receiver.String(),
		"target_action", targetAction.String(),
		"required_amount", requiredAmount,
		"adjustmentFee", adjustmentFee)

	// get listUnspent
	unspentList, unspentAddrs, err := u.getUnspentList(sender)
	if err != nil {
		return "", "", fmt.Errorf("fail to call getUnspentList(): %w", err)
	}
	if len(unspentList) == 0 {
		logger.Info("no listunspent")
		return "", "", nil
	}

	// parse listUnspent
	parsedTx, inputTotal, isDone := u.parseListUnspentTx(unspentList, requiredAmount)
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
		txPrevOutputs, err = u.createTxOutputs(receiver, requiredAmount, inputTotal, unspentAddrs[0], isChange)
		if err != nil {
			return "", "", fmt.Errorf("fail to call createTxOutputs(): %w", err)
		}
	case domainTx.ActionTypePayment:
		changeAddr := unspentAddrs[0] // this is actually sender's address because it's for change
		changeAmount := inputTotal - requiredAmount
		txPrevOutputs = u.createPaymentTxOutputs(userPayments, changeAddr, changeAmount)
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
	msgTx, err := u.btcClient.CreateRawTransaction(parsedTx.txInputs, txPrevOutputs)
	if err != nil {
		return "", "", fmt.Errorf("fail to call btc.CreateRawTransaction(): %w", err)
	}

	// calculate fee and output total
	//  - adjust outputTotal by fee and re-run CreateRawTransaction
	//  - this logic would be different from payment
	// Pass the change address (unspentAddrs[0]) to identify which output should have fee subtracted
	changeAddr := unspentAddrs[0]
	outputTotal, fee, txOutputs, txRepoTxOutputs, err := u.calculateOutputTotal(
		sender, receiver, msgTx, adjustmentFee, inputTotal, txPrevOutputs, changeAddr)
	if err != nil {
		return "", "", err
	}

	// re call CreateRawTransaction
	msgTx, err = u.btcClient.CreateRawTransaction(parsedTx.txInputs, txOutputs)
	if err != nil {
		return "", "", fmt.Errorf("fail to call btc.CreateRawTransaction(): %w", err)
	}

	// convert msgTx to hex
	hex, err := u.btcClient.ToHex(msgTx)
	if err != nil {
		return "", "", fmt.Errorf("fail to call btc.ToHex(msgTx): %w", err)
	}

	// insert to tx_table for unsigned tx
	//  - txID would be 0 if record is already existing then csv file is not created
	txID, err := u.insertTxTableForUnsigned(
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

	// generate PSBT file
	// TODO: how to recover when error occurred here
	// - inserted data in database must be deleted to generate PSBT file
	var generatedFileName string
	if txID != 0 {
		generatedFileName, err = u.generatePSBTFile(targetAction, msgTx, parsedTx.prevTxs, sender, txID)
		if err != nil {
			return "", "", fmt.Errorf("fail to call generatePSBTFile(): %w", err)
		}
	}

	logger.Debug("createTx completed",
		"unspentList", unspentList,
		"unspentAddrs", unspentAddrs,
		"requiredAmount", requiredAmount,
		"input_total", inputTotal,
		"len(inputs)", len(parsedTx.txInputs),
		"len(prevTxs)", len(parsedTx.prevTxs),
		"len(txPrevOutputs)", len(txPrevOutputs),
		"len(txOutputs)", len(txOutputs),
		"len(txRepoTxOutputs)", len(txRepoTxOutputs),
		"sender_account", sender.String(),
		"generated_file", generatedFileName,
	)

	return hex, generatedFileName, nil
}

// call API `listunspent`
// this func returns no result, no error possibly, so caller should check both returned value
func (u *createTransactionUseCase) getUnspentList(
	accountType domainAccount.AccountType,
) ([]dtobtc.UnspentOutput, []string, error) {
	// get listUnspent
	unspentList, err := u.btcClient.ListUnspentByAccount(accountType, u.btcClient.ConfirmationBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call btc.ListUnspentByAccount(): %w", err)
	}
	unspentAddrs := u.btcClient.GetUnspentListAddrs(unspentList, accountType)

	return unspentList, unspentAddrs, nil
}

// parse result of listUnspent
// returned *parsedTx could be nil
func (u *createTransactionUseCase) parseListUnspentTx(
	unspentList []dtobtc.UnspentOutput, amount btcutil.Amount,
) (*parsedTx, btcutil.Amount, bool) {
	var inputTotal btcutil.Amount
	txInputs := make([]btcjson.TransactionInput, 0, len(unspentList))
	txRepoTxInputs := make([]*domainBitcoin.BTCTxInput, 0, len(unspentList))
	prevTxs := make([]dtobtc.PreviousTx, 0, len(unspentList))
	addresses := make([]string, 0, len(unspentList))

	var isDone bool // if isDone is false, sender can't meet amount
	if amount == 0 {
		isDone = true
	}

	for _, txItem := range unspentList {
		// Amount (already btcutil.Amount, no conversion needed)
		inputTotal += txItem.Amount

		txInputs = append(txInputs, btcjson.TransactionInput{
			Txid: txItem.TxID,
			Vout: txItem.Vout,
		})

		// Convert btcutil.Amount to float64 for decimal conversion
		inputAmount, err := u.btcClient.FloatToDecimal(txItem.Amount.ToBTC())
		if err != nil {
			logger.Error("fail to convert input amount to decimal", "error", err)
			continue
		}
		input, err := domainBitcoin.NewBTCTxInput(
			0, // TxID will be set later
			txItem.TxID,
			txItem.Vout,
			txItem.Address,
			txItem.Label,
			inputAmount.String(),
			uint64(txItem.Confirmations),
		)
		if err != nil {
			logger.Error("fail to create BtcTxInput", "error", err)
			continue
		}
		txRepoTxInputs = append(txRepoTxInputs, input)

		// TODO: if sender is client account (non-multisig address), RedeemScript is blank
		prevTxs = append(prevTxs, dtobtc.PreviousTx{
			TxID:          txItem.TxID,
			Vout:          txItem.Vout,
			ScriptPubKey:  txItem.ScriptPubKey,
			RedeemScript:  txItem.RedeemScript,  // required if target account is multisig address
			WitnessScript: txItem.WitnessScript, // for SegWit addresses
			Amount:        txItem.Amount,        // Keep as btcutil.Amount
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

// createTxOutputs creates transaction outputs for ActionTypeDeposit and ActionTypeTransfer.
// This function supports all Bitcoin address types including:
//   - P2PKH (Legacy, 1...)
//   - P2SH-SegWit (3...)
//   - P2WPKH (Bech32, bc1q...)
//   - P2TR (Taproot, bc1p...) - BIP86
//
// The address type is automatically detected by btcutil.DecodeAddress() based on the
// address format stored in the database. The underlying btcsuite library and Bitcoin Core
// RPC handle creating the appropriate scriptPubKey for each address type.
func (u *createTransactionUseCase) createTxOutputs(
	receiver domainAccount.AccountType,
	requiredAmount btcutil.Amount,
	inputTotal btcutil.Amount,
	senderAddr string,
	isChange bool,
) (map[btcutil.Address]btcutil.Amount, error) {
	// get unallocated address for receiver
	// - deposit/transfer
	pubkeyTable, err := u.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return nil, fmt.Errorf("fail to call pubkeyRepo.GetOneUnAllocated(): %w", err)
	}
	receiverAddr := pubkeyTable.WalletAddress

	// create receiver txOutput
	// DecodeAddress automatically recognizes all address types including Taproot (bc1p...)
	receiverDecodedAddr, err := btcutil.DecodeAddress(receiverAddr, u.btcClient.GetChainConf())
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
		senderDecodedAddr, decodeErr := btcutil.DecodeAddress(senderAddr, u.btcClient.GetChainConf())
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

// createPaymentTxOutputs creates transaction outputs for ActionTypePayment.
// This function supports all Bitcoin address types for payment receivers including:
//   - P2PKH (Legacy, 1...)
//   - P2SH-SegWit (3...)
//   - P2WPKH (Bech32, bc1q...)
//   - P2TR (Taproot, bc1p...) - BIP86
//
// The address type is automatically detected by btcutil.DecodeAddress() when processing
// user payment addresses from the payment_request table.
func (u *createTransactionUseCase) createPaymentTxOutputs(
	userPayments []userPayment, changeAddr string, changeAmount btcutil.Amount,
) map[btcutil.Address]btcutil.Amount {
	var (
		txOutputs = map[btcutil.Address]btcutil.Amount{}
		// if key of map is btcutil.Address which is interface type, uniqueness can't be found from map key
		// so this key is string
		tmpOutputs = map[string]btcutil.Amount{}
	)

	// create txOutput from userPayment
	for _, userPayment := range userPayments {
		tmpOutputs[userPayment.receiverAddr] += userPayment.validAmount
	}

	// add txOutput as change address and amount for change
	// TODO:
	// - what if user register for address which is same to payment address?
	//   Though it's impossible in real but systematically possible
	// - BIP44, hdwallet has `ChangeType`. ideally this address should be used
	tmpOutputs[changeAddr] += changeAmount

	// create txOutputs from tmpOutputs switching string address type to btcutil.Address
	for strAddr, amount := range tmpOutputs {
		addr, err := u.btcClient.DecodeAddress(strAddr)
		if err != nil {
			// this case is impossible because addresses are checked in advance
			logger.Error("fail to call DecodeAddress",
				"address", strAddr)
			continue
		}
		txOutputs[addr] = amount
	}

	return txOutputs
}

func (u *createTransactionUseCase) calculateOutputTotal(
	sender domainAccount.AccountType,
	receiver domainAccount.AccountType,
	msgTx *wire.MsgTx,
	adjustmentFee float64,
	inputTotal btcutil.Amount,
	txPrevOutputs map[btcutil.Address]btcutil.Amount,
	changeAddrStr string, // Address string of the change output (where fee should be subtracted from)
) (btcutil.Amount, btcutil.Amount, map[btcutil.Address]btcutil.Amount, []*domainBitcoin.BTCTxOutput, error) {
	// get fee
	fee, err := u.btcClient.GetFee(msgTx, adjustmentFee)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("fail to call btc.GetFee(): %w", err)
	}
	var outputTotal btcutil.Amount
	txRepoOutputs := make([]*domainBitcoin.BTCTxOutput, 0, len(txPrevOutputs))

	// subtract fee from output transaction for change
	// FIXME: what if change is short, should re-run from the beginning with shortage-flag
	for addr, amt := range txPrevOutputs {
		// Check if this address is the change address (by comparing address strings)
		// This works with descriptor-based wallets where GetAccount() may not work correctly
		isChangeAddr := addr.String() == changeAddrStr

		if len(txPrevOutputs) == 1 {
			// Only one output - subtract fee from it
			txPrevOutputs[addr] -= fee
			outputAmount, err := u.btcClient.AmountToDecimal(amt - fee)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to convert output amount to decimal: %w", err)
			}
			output, err := domainBitcoin.NewBTCTxOutput(
				0, // TxID will be set later
				addr.String(),
				receiver.String(),
				outputAmount.String(),
				false,
			)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to create BtcTxOutput: %w", err)
			}
			txRepoOutputs = append(txRepoOutputs, output)
			outputTotal += amt
			break
		}

		if isChangeAddr {
			logger.Debug("detect change address in calculateOutputTotal", "address", addr.String())
			// address is used for change - subtract fee from it
			txPrevOutputs[addr] -= fee
			outputAmount, err := u.btcClient.AmountToDecimal(amt - fee)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to convert change amount to decimal: %w", err)
			}
			output, err := domainBitcoin.NewBTCTxOutput(
				0, // TxID will be set later
				addr.String(),
				sender.String(),
				outputAmount.String(),
				true, // Mark as change output
			)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to create BtcTxOutput: %w", err)
			}
			txRepoOutputs = append(txRepoOutputs, output)
		} else {
			// Not change address - don't subtract fee
			outputAmount, err := u.btcClient.AmountToDecimal(amt)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to convert output amount to decimal: %w", err)
			}
			output, err := domainBitcoin.NewBTCTxOutput(
				0, // TxID will be set later
				addr.String(),
				receiver.String(),
				outputAmount.String(),
				false,
			)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to create BtcTxOutput: %w", err)
			}
			txRepoOutputs = append(txRepoOutputs, output)
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

func (u *createTransactionUseCase) insertTxTableForUnsigned(
	actionType domainTx.ActionType,
	hex string,
	inputTotal,
	outputTotal,
	fee btcutil.Amount,
	txInputs []*domainBitcoin.BTCTxInput,
	txOutputs []*domainBitcoin.BTCTxOutput,
	paymentRequestIds []int64,
) (int64, error) {
	// skip if same hex is already stored
	count, err := u.txRepo.GetCountByUnsignedHex(actionType, hex)
	if err != nil {
		return 0, fmt.Errorf("fail to call repo.Tx().GetCountByUnsignedHex(): %w", err)
	}
	if count != 0 {
		// skip
		return 0, nil
	}

	// TxReceipt table
	totalInputAmt, err := u.btcClient.AmountToDecimal(inputTotal)
	if err != nil {
		return 0, fmt.Errorf("fail to convert total input amount to decimal: %w", err)
	}
	totalOutputAmt, err := u.btcClient.AmountToDecimal(outputTotal)
	if err != nil {
		return 0, fmt.Errorf("fail to convert total output amount to decimal: %w", err)
	}
	feeAmt, err := u.btcClient.AmountToDecimal(fee)
	if err != nil {
		return 0, fmt.Errorf("fail to convert fee amount to decimal: %w", err)
	}
	txItem := domainBitcoin.NewBTCTransaction(
		u.btcClient.CoinTypeCode(),
		actionType,
		domainTx.TxTypeUnsigned,
	)
	txItem.SetUnsignedTx(hex, totalInputAmt.String(), totalOutputAmt.String(), feeAmt.String())

	// start database transaction
	dtx, err := u.dbConn.Begin()
	if err != nil {
		return 0, fmt.Errorf("fail to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = dtx.Rollback() // Error already being handled
		} else {
			_ = dtx.Commit() // Error already being handled
		}
	}()

	txID, err := u.txRepo.InsertUnsignedTx(actionType, txItem)
	if err != nil {
		return 0, fmt.Errorf("fail to call repo.Tx().InsertUnsignedTx(): %w", err)
	}

	// TxReceiptInput table
	//  update txID
	for idx := range txInputs {
		txInputs[idx].TxID = txID
	}
	err = u.txInputRepo.InsertBulk(txInputs)
	if err != nil {
		return 0, fmt.Errorf("fail to call txInRepo.InsertBulk(): %w", err)
	}

	// TxReceiptOutput table
	//  update txID
	for idx := range txOutputs {
		txOutputs[idx].TxID = txID
	}
	err = u.txOutputRepo.InsertBulk(txOutputs)
	if err != nil {
		return 0, fmt.Errorf("fail to call repo.TxOutput().InsertBulk(): %w", err)
	}

	// update payment_id in payment_request table for only domainTx.ActionTypePayment
	if actionType == domainTx.ActionTypePayment {
		_, err = u.payReqRepo.UpdatePaymentID(txID, paymentRequestIds)
		if err != nil {
			return 0, fmt.Errorf("fail to call repo.PayReq().UpdatePaymentID(txID, paymentRequestIds): %w", err)
		}
	}

	return txID, nil
}

// generatePSBTFile creates PSBT file for unsigned transaction with all required metadata.
// This replaces the legacy CSV-based generateHexFile method.
//
// The PSBT includes:
//   - All transaction inputs and outputs
//   - Previous output metadata (amounts, scriptPubKeys)
//   - Redeem scripts (for P2SH multisig)
//   - Witness scripts (for P2WSH multisig)
//   - Taproot metadata (for P2TR addresses)
//
// The generated PSBT file follows BIP174 format and is compatible with:
//   - Keygen wallet (offline signing)
//   - Sign wallet (multisig second signature)
//   - Hardware wallets (via BIP32 derivation paths)
func (u *createTransactionUseCase) generatePSBTFile(
	actionType domainTx.ActionType,
	msgTx *wire.MsgTx,
	prevTxs []dtobtc.PreviousTx,
	senderAccount domainAccount.AccountType,
	id int64,
) (string, error) {
	// Create PSBT from msgTx and previous outputs
	// This includes all necessary metadata for offline signing
	psbtBase64, err := u.btcClient.CreatePSBT(msgTx, prevTxs, senderAccount)
	if err != nil {
		return "", fmt.Errorf("fail to create PSBT: %w", err)
	}

	// Create file path with .psbt extension
	path := u.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeUnsigned, id, 0)

	// Write PSBT file
	generatedFileName, err := u.txFileRepo.WritePSBTFile(path, psbtBase64)
	if err != nil {
		return "", fmt.Errorf("fail to write PSBT file: %w", err)
	}

	logger.Debug("generated PSBT file",
		"action", actionType.String(),
		"txID", id,
		"fileName", generatedFileName,
		"inputs", len(prevTxs),
		"sender", senderAccount.String(),
	)

	return generatedFileName, nil
}

// createUserPayment gets payment data from payment_request table
func (u *createTransactionUseCase) createUserPayment() ([]userPayment, []int64, error) {
	// get payment_request
	paymentRequests, err := u.payReqRepo.GetAll()
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call repo.GetPaymentRequestAll(): %w", err)
	}
	if len(paymentRequests) == 0 {
		logger.Debug("no data in payment_request")
		return nil, nil, nil
	}

	userPayments := make([]userPayment, len(paymentRequests))
	paymentRequestIds := make([]int64, len(paymentRequests))

	// store `id` separately for key updating
	for idx, val := range paymentRequests {
		paymentRequestIds[idx] = val.ID

		userPayments[idx].senderAddr = val.SenderAddress
		userPayments[idx].receiverAddr = val.ReceiverAddress
		amt, parseErr := strconv.ParseFloat(val.Amount, 64)
		if parseErr != nil {
			// fatal error because table includes invalid data
			logger.Error("payment_request table includes invalid amount field")
			return nil, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].amount = amt

		// decode address
		// TODO: may be it's not necessary
		userPayments[idx].validRecAddr, err = u.btcClient.DecodeAddress(userPayments[idx].receiverAddr)
		if err != nil {
			// fatal error
			logger.Error("unexpected error occurred converting receiverAddr from string type to address type")
			return nil, nil, errors.New(
				"unexpected error occurred converting receiverAddr from string type to address type",
			)
		}

		// amount
		userPayments[idx].validAmount, err = u.btcClient.FloatToAmount(userPayments[idx].amount)
		if err != nil {
			// fatal error
			logger.Error("unexpected error occurred converting amount from float64 type to Amount type")
			return nil, nil, errors.New("unexpected error occurred converting amount from float64 type to Amount type")
		}
	}

	return userPayments, paymentRequestIds, nil
}
