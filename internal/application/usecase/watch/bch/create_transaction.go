package bch

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/shared"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// createTxBCHClient defines the minimal interface needed for BCH transaction creation.
// This follows the Interface Segregation Principle - depend only on methods actually used.
// BCH does not use PSBT, so this interface excludes PSBT-related methods.
type createTxBCHClient interface {
	apibtc.ChainConfigProvider
	apibtc.AmountConverter
	apibtc.UTXOProvider
	apibtc.RawTransactionCreator
	apibtc.AddressOperator
	apibtc.BalanceChecker
}

// createTransactionUseCase implements BCH-specific transaction creation.
// Unlike BTC, BCH uses Raw Transaction Hex format instead of PSBT.
type createTransactionUseCase struct {
	bchClient       createTxBCHClient
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
	txDBHelper      *shared.TxDBHelper
}

// NewCreateBCHTransactionUseCase creates a new BCH CreateTransactionUseCase.
// BCH uses Raw Transaction Hex format instead of PSBT.
func NewCreateBCHTransactionUseCase(
	bchClient apibtc.BCHer,
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
	// Runtime type verification for debugging
	logger.Debug("BCH use case: received client type",
		"type", fmt.Sprintf("%T", bchClient),
		"coin", bchClient.CoinTypeCode())

	return &createTransactionUseCase{
		bchClient:       bchClient,
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
		txDBHelper: shared.NewTxDBHelper(
			dbConn, txRepo, txInputRepo, txOutputRepo, payReqRepo, bchClient,
		),
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
func (u *createTransactionUseCase) createDepositTx(adjustmentFee float64) (string, string, error) {
	sender := domainAccount.AccountTypeClient
	receiver := u.depositReceiver
	targetAction := domainTx.ActionTypeDeposit
	logger.Debug("account",
		"sender", sender.String(),
		"receiver", receiver.String(),
	)

	requiredAmount, err := u.bchClient.FloatToAmount(0)
	if err != nil {
		return "", "", err
	}

	return u.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, nil, nil)
}

// createPaymentTx creates unsigned tx for user (anonymous addresses)
func (u *createTransactionUseCase) createPaymentTx(adjustmentFee float64) (string, string, error) {
	sender := u.paymentSender
	receiver := domainAccount.AccountTypeAnonymous
	targetAction := domainTx.ActionTypePayment
	logger.Debug("account",
		"sender", sender.String(),
		"receiver", receiver.String(),
	)

	userPayments, paymentRequestIds, err := u.createUserPayment()
	if err != nil {
		return "", "", err
	}
	if len(userPayments) == 0 {
		logger.Debug("no data in userPayments")
		return "", "", nil
	}

	var requiredAmount btcutil.Amount
	for _, val := range userPayments {
		requiredAmount += val.ValidAmount
	}

	balance, err := u.bchClient.GetBalanceByAccount(domainAccount.AccountTypePayment, u.bchClient.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= requiredAmount {
		logger.Info("balance for payment account is insufficient",
			"payment_balance", balance.ToBTC(),
			"required_amount", requiredAmount.ToBTC(),
		)
		return "", "", nil
	}

	return u.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, paymentRequestIds, userPayments)
}

// createTransferTx creates unsigned tx for transfer coin among internal accounts
func (u *createTransactionUseCase) createTransferTx(
	sender, receiver domainAccount.AccountType, floatAmount, adjustmentFee float64,
) (string, string, error) {
	targetAction := domainTx.ActionTypeTransfer

	if receiver == domainAccount.AccountTypeClient || receiver == domainAccount.AccountTypeAuthorization {
		return "", "", errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return "", "", errors.New("invalid account. sender and receiver is same")
	}

	requiredAmount, err := u.bchClient.FloatToAmount(floatAmount)
	if err != nil {
		return "", "", err
	}

	balance, err := u.bchClient.GetBalanceByAccount(sender, u.bchClient.ConfirmationBlock())
	if err != nil {
		return "", "", err
	}
	if balance <= requiredAmount {
		return "", "", fmt.Errorf("account: %s balance is insufficient", sender)
	}

	return u.createTx(sender, receiver, targetAction, requiredAmount, adjustmentFee, nil, nil)
}

// createTx creates a BCH transaction using Raw Transaction Hex format (not PSBT).
//

func (u *createTransactionUseCase) createTx(
	sender,
	receiver domainAccount.AccountType,
	targetAction domainTx.ActionType,
	requiredAmount btcutil.Amount,
	adjustmentFee float64,
	paymentRequestIds []int64,
	userPayments []shared.UserPayment,
) (string, string, error) {
	logger.Debug("createTx()",
		"sender_account", sender.String(),
		"receiver_account", receiver.String(),
		"target_action", targetAction.String(),
		"required_amount", requiredAmount,
		"adjustmentFee", adjustmentFee)

	// Get listUnspent
	unspentList, unspentAddrs, err := u.getUnspentList(sender)
	if err != nil {
		return "", "", fmt.Errorf("fail to call getUnspentList(): %w", err)
	}
	if len(unspentList) == 0 {
		logger.Info("no listunspent")
		return "", "", nil
	}

	// Parse listUnspent using shared helper
	parsedTx, inputTotal, isDone := shared.ParseListUnspentTx(unspentList, requiredAmount, u.bchClient)
	if len(parsedTx.TxInputs) == 0 {
		logger.Info("no input tx in listUnspent")
		return "", "", nil
	}
	if !isDone {
		return "", "", errors.New("sender account can't meet amount to send")
	}

	// Create txOutputs
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
		changeAddr := unspentAddrs[0]
		changeAmount := inputTotal - requiredAmount
		txPrevOutputs = u.createPaymentTxOutputs(userPayments, changeAddr, changeAmount)
	default:
		return "", "", fmt.Errorf("invalid actionType: %s", targetAction)
	}

	// Create raw transaction
	msgTx, err := u.bchClient.CreateRawTransaction(parsedTx.TxInputs, txPrevOutputs)
	if err != nil {
		return "", "", fmt.Errorf("fail to call bch.CreateRawTransaction(): %w", err)
	}

	// Calculate fee and output total
	changeAddr := unspentAddrs[0]
	outputTotal, fee, txOutputs, txRepoTxOutputs, err := u.calculateOutputTotal(
		sender, receiver, msgTx, adjustmentFee, inputTotal, txPrevOutputs, changeAddr)
	if err != nil {
		return "", "", err
	}
	logger.Debug("transaction fee calculated",
		"input_total", inputTotal.ToBTC(),
		"output_total", outputTotal.ToBTC(),
		"fee", fee.ToBTC())

	// Re-create raw transaction with adjusted outputs
	msgTx, err = u.bchClient.CreateRawTransaction(parsedTx.TxInputs, txOutputs)
	if err != nil {
		return "", "", fmt.Errorf("fail to call bch.CreateRawTransaction(): %w", err)
	}

	// Convert msgTx to hex
	hex, err := u.bchClient.ToHex(msgTx)
	if err != nil {
		return "", "", fmt.Errorf("fail to call bch.ToHex(msgTx): %w", err)
	}

	// Insert to tx_table for unsigned tx using shared helper
	txID, err := u.txDBHelper.InsertTxTableForUnsigned(
		targetAction,
		hex,
		inputTotal,
		outputTotal,
		fee,
		parsedTx.TxRepoTxInputs,
		txRepoTxOutputs,
		paymentRequestIds)
	if err != nil {
		return "", "", fmt.Errorf("fail to call insertTxTableForUnsigned(): %w", err)
	}

	// Generate Raw TX Hex file (BCH uses hex, not PSBT)
	var generatedFileName string
	if txID != 0 {
		generatedFileName, err = u.generateRawTxHexFile(targetAction, hex, parsedTx.PrevTxs, txID)
		if err != nil {
			return "", "", fmt.Errorf("fail to call generateRawTxHexFile(): %w", err)
		}
	}

	logger.Debug("createTx completed",
		"unspentList", unspentList,
		"input_total", inputTotal,
		"len(inputs)", len(parsedTx.TxInputs),
		"len(prevTxs)", len(parsedTx.PrevTxs),
		"sender_account", sender.String(),
		"generated_file", generatedFileName,
	)

	return hex, generatedFileName, nil
}

// getUnspentList calls API `listunspent`
func (u *createTransactionUseCase) getUnspentList(
	accountType domainAccount.AccountType,
) ([]dtobtc.UnspentOutput, []string, error) {
	unspentList, err := u.bchClient.ListUnspentByAccount(accountType, u.bchClient.ConfirmationBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call bch.ListUnspentByAccount(): %w", err)
	}
	unspentAddrs := u.bchClient.GetUnspentListAddrs(unspentList, accountType)

	return unspentList, unspentAddrs, nil
}

// createTxOutputs creates transaction outputs for deposit and transfer
func (u *createTransactionUseCase) createTxOutputs(
	receiver domainAccount.AccountType,
	requiredAmount btcutil.Amount,
	inputTotal btcutil.Amount,
	senderAddr string,
	isChange bool,
) (map[btcutil.Address]btcutil.Amount, error) {
	pubkeyTable, err := u.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return nil, fmt.Errorf("fail to call pubkeyRepo.GetOneUnAllocated(): %w", err)
	}
	receiverAddr := pubkeyTable.WalletAddress

	receiverDecodedAddr, err := btcutil.DecodeAddress(receiverAddr, u.bchClient.GetChainConf())
	if err != nil {
		return nil, fmt.Errorf("fail to call btcutil.DecodeAddress(%s): %w", receiverAddr, err)
	}
	txPrevOutputs := make(map[btcutil.Address]btcutil.Amount)
	if isChange {
		txPrevOutputs[receiverDecodedAddr] = requiredAmount
	} else {
		txPrevOutputs[receiverDecodedAddr] = inputTotal
	}

	if isChange {
		senderDecodedAddr, decodeErr := btcutil.DecodeAddress(senderAddr, u.bchClient.GetChainConf())
		if decodeErr != nil {
			return nil, fmt.Errorf("fail to call btcutil.DecodeAddress(%s): %w", senderAddr, decodeErr)
		}
		txPrevOutputs[senderDecodedAddr] = inputTotal - requiredAmount
	}
	return txPrevOutputs, nil
}

// createPaymentTxOutputs creates transaction outputs for payment
func (u *createTransactionUseCase) createPaymentTxOutputs(
	userPayments []shared.UserPayment, changeAddr string, changeAmount btcutil.Amount,
) map[btcutil.Address]btcutil.Amount {
	var (
		txOutputs  = map[btcutil.Address]btcutil.Amount{}
		tmpOutputs = map[string]btcutil.Amount{}
	)

	for _, userPayment := range userPayments {
		tmpOutputs[userPayment.ReceiverAddr] += userPayment.ValidAmount
	}
	tmpOutputs[changeAddr] += changeAmount

	for strAddr, amount := range tmpOutputs {
		addr, err := u.bchClient.DecodeAddress(strAddr)
		if err != nil {
			logger.Error("fail to call DecodeAddress", "address", strAddr)
			continue
		}
		txOutputs[addr] = amount
	}

	return txOutputs
}

// calculateOutputTotal calculates output total and adjusts for fee
func (u *createTransactionUseCase) calculateOutputTotal(
	sender domainAccount.AccountType,
	receiver domainAccount.AccountType,
	msgTx *wire.MsgTx,
	adjustmentFee float64,
	inputTotal btcutil.Amount,
	txPrevOutputs map[btcutil.Address]btcutil.Amount,
	changeAddrStr string,
) (btcutil.Amount, btcutil.Amount, map[btcutil.Address]btcutil.Amount, []*domainBitcoin.BTCTxOutput, error) {
	logger.Debug("BCH use case: calculateOutputTotal called",
		"tx_size", msgTx.SerializeSize(),
		"adjustment_fee", adjustmentFee,
		"input_total", inputTotal.ToBTC())

	// Call GetFee through interface - Go should dispatch to BitcoinCash.GetFee override
	logger.Debug("calling GetFee",
		"tx_size", msgTx.SerializeSize(),
		"adjustment_fee", adjustmentFee)
	fee, err := u.bchClient.GetFee(msgTx, adjustmentFee)
	if err != nil {
		return 0, 0, nil, nil, fmt.Errorf("fail to call bch.GetFee(): %w", err)
	}
	logger.Debug("BCH fee calculated",
		"fee_satoshis", fee,
		"fee_btc", fee.ToBTC(),
		"tx_size", msgTx.SerializeSize())

	// BCH addresses may have network prefix (bchreg:, bitcoincash:, bchtest:)
	// Strip prefix for comparison since addr.String() doesn't include it
	changeAddrWithoutPrefix := strings.TrimPrefix(changeAddrStr, "bchreg:")
	changeAddrWithoutPrefix = strings.TrimPrefix(changeAddrWithoutPrefix, "bitcoincash:")
	changeAddrWithoutPrefix = strings.TrimPrefix(changeAddrWithoutPrefix, "bchtest:")

	var outputTotal btcutil.Amount
	txRepoOutputs := make([]*domainBitcoin.BTCTxOutput, 0, len(txPrevOutputs))

	for addr, amt := range txPrevOutputs {
		isChangeAddr := addr.String() == changeAddrWithoutPrefix

		if len(txPrevOutputs) == 1 {
			txPrevOutputs[addr] -= fee
			outputAmount, err := u.bchClient.AmountToDecimal(amt - fee)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to convert output amount to decimal: %w", err)
			}
			output, err := domainBitcoin.NewBTCTxOutput(
				0,
				addr.String(),
				receiver.String(),
				outputAmount.String(),
				false,
			)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to create BtcTxOutput: %w", err)
			}
			txRepoOutputs = append(txRepoOutputs, output)
			break
		}

		if isChangeAddr {
			txPrevOutputs[addr] -= fee
			outputAmount, err := u.bchClient.AmountToDecimal(amt - fee)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to convert change amount to decimal: %w", err)
			}
			output, err := domainBitcoin.NewBTCTxOutput(
				0,
				addr.String(),
				sender.String(),
				outputAmount.String(),
				true,
			)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to create BtcTxOutput: %w", err)
			}
			txRepoOutputs = append(txRepoOutputs, output)
		} else {
			outputAmount, err := u.bchClient.AmountToDecimal(amt)
			if err != nil {
				return 0, 0, nil, nil, fmt.Errorf("fail to convert output amount to decimal: %w", err)
			}
			output, err := domainBitcoin.NewBTCTxOutput(
				0,
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
	}

	outputTotal = inputTotal - fee

	if outputTotal <= 0 {
		return 0, 0, nil, nil, fmt.Errorf(
			"inputTotal is short of coin to pay fee: inputTotal=%d, fee=%d", inputTotal, fee)
	}

	return outputTotal, fee, txPrevOutputs, txRepoOutputs, nil
}

// generateRawTxHexFile creates a Raw Transaction Hex file for BCH unsigned transaction.
// BCH uses Raw TX Hex format instead of PSBT (which is BTC-only).
//
// The file contains:
//   - Transaction hex
//   - Previous transaction metadata (for signing reference)
func (u *createTransactionUseCase) generateRawTxHexFile(
	actionType domainTx.ActionType,
	txHex string,
	prevTxs []dtobtc.PreviousTx,
	id int64,
) (string, error) {
	// Create file path (use .hex extension for BCH raw transactions)
	// CreateFilePath returns path ending with underscore: payment_1_unsign_0_
	// WriteHexFile adds timestamp and .hex extension (similar to WritePSBTFile pattern)
	// Final format: payment_1_unsigned_0_1768816802572944000.hex
	basePath := u.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeUnsigned, id, 0)

	// Write transaction hex to file
	// For BCH, we store raw hex along with prevTx metadata
	content := u.formatRawTxContent(txHex, prevTxs)
	generatedFileName, err := u.txFileRepo.WriteHexFile(basePath, content)
	if err != nil {
		return "", fmt.Errorf("fail to write raw tx file: %w", err)
	}

	logger.Debug("generated BCH raw tx file",
		"action", actionType.String(),
		"txID", id,
		"fileName", generatedFileName,
		"inputs", len(prevTxs),
	)

	return generatedFileName, nil
}

// formatRawTxContent formats the raw transaction content for BCH.
// Format: txHex followed by prevTx metadata (JSON-like format for each prevTx)
func (*createTransactionUseCase) formatRawTxContent(txHex string, prevTxs []dtobtc.PreviousTx) string {
	// Simple format: hex on first line, then prevTx info
	var builder strings.Builder
	builder.WriteString(txHex)
	builder.WriteString("\n")

	for _, prev := range prevTxs {
		fmt.Fprintf(&builder, "prevtx:%s:%d:%s:%s:%d\n",
			prev.TxID,
			prev.Vout,
			prev.ScriptPubKey,
			prev.RedeemScript,
			prev.Amount,
		)
	}
	return builder.String()
}

// createUserPayment gets payment data from payment_request table
func (u *createTransactionUseCase) createUserPayment() ([]shared.UserPayment, []int64, error) {
	paymentRequests, err := u.payReqRepo.GetAll()
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call repo.GetPaymentRequestAll(): %w", err)
	}
	if len(paymentRequests) == 0 {
		logger.Debug("no data in payment_request")
		return nil, nil, nil
	}

	userPayments := make([]shared.UserPayment, len(paymentRequests))
	paymentRequestIds := make([]int64, len(paymentRequests))

	for idx, val := range paymentRequests {
		paymentRequestIds[idx] = val.ID

		userPayments[idx].SenderAddr = val.SenderAddress
		userPayments[idx].ReceiverAddr = val.ReceiverAddress
		amt, parseErr := strconv.ParseFloat(val.Amount, 64)
		if parseErr != nil {
			logger.Error("payment_request table includes invalid amount field")
			return nil, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].Amount = amt

		// Validate address format (decoded address stored in ValidRecAddr for potential future use)
		// Note: createPaymentTxOutputs re-decodes addresses to aggregate payments to same address
		userPayments[idx].ValidRecAddr, err = u.bchClient.DecodeAddress(userPayments[idx].ReceiverAddr)
		if err != nil {
			logger.Error("invalid receiver address format", "address", userPayments[idx].ReceiverAddr)
			return nil, nil, fmt.Errorf(
				"invalid receiver address format (%s): %w", userPayments[idx].ReceiverAddr, err)
		}

		userPayments[idx].ValidAmount, err = u.bchClient.FloatToAmount(userPayments[idx].Amount)
		if err != nil {
			logger.Error("unexpected error occurred converting amount from float64 type to Amount type")
			return nil, nil, errors.New("unexpected error occurred converting amount from float64 type to Amount type")
		}
	}

	return userPayments, paymentRequestIds, nil
}
