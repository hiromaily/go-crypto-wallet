package xrp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/bookerzzz/grok"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple/xrp"
	watchrepo "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

type createTransactionUseCase struct {
	rippler         ripple.Rippler
	dbConn          *sql.DB
	uuidHandler     uuid.UUIDHandler
	addrRepo        watchrepo.AddressRepositorier
	txRepo          watchrepo.TxRepositorier
	txDetailRepo    watchrepo.XrpDetailTxRepositorier
	payReqRepo      watchrepo.PaymentRequestRepositorier
	txFileRepo      file.TransactionFileRepositorier
	depositReceiver domainAccount.AccountType
	paymentSender   domainAccount.AccountType
}

// NewCreateTransactionUseCase creates a new CreateTransactionUseCase
func NewCreateTransactionUseCase(
	rippler ripple.Rippler,
	dbConn *sql.DB,
	uuidHandler uuid.UUIDHandler,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.TxRepositorier,
	txDetailRepo watchrepo.XrpDetailTxRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	depositReceiver domainAccount.AccountType,
	paymentSender domainAccount.AccountType,
) watchusecase.CreateTransactionUseCase {
	return &createTransactionUseCase{
		rippler:         rippler,
		dbConn:          dbConn,
		uuidHandler:     uuidHandler,
		addrRepo:        addrRepo,
		txRepo:          txRepo,
		txDetailRepo:    txDetailRepo,
		payReqRepo:      payReqRepo,
		txFileRepo:      txFileRepo,
		depositReceiver: depositReceiver,
		paymentSender:   paymentSender,
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

	var fileName string
	var execErr error

	switch actionType {
	case domainTx.ActionTypeDeposit:
		fileName, execErr = u.createDepositTx(ctx)
	case domainTx.ActionTypePayment:
		fileName, execErr = u.createPaymentTx(ctx)
	case domainTx.ActionTypeTransfer:
		fileName, execErr = u.createTransferTx(ctx, input.SenderAccount, input.ReceiverAccount, input.Amount)
	default:
		return watchusecase.CreateTransactionOutput{}, fmt.Errorf("unsupported action type: %s", input.ActionType)
	}

	if execErr != nil {
		return watchusecase.CreateTransactionOutput{}, fmt.Errorf("failed to create transaction: %w", execErr)
	}

	return watchusecase.CreateTransactionOutput{
		TransactionHex: "",
		FileName:       fileName,
	}, nil
}

// createDepositTx creates unsigned tx if client accounts have coins
// - sender: client, receiver: deposit
// - receiver account covers fee, but this should be flexible
func (u *createTransactionUseCase) createDepositTx(ctx context.Context) (string, error) {
	sender := domainAccount.AccountTypeClient
	receiver := u.depositReceiver
	targetAction := domainTx.ActionTypeDeposit
	logger.Debug("account",
		"sender", sender.String(),
		"receiver", receiver.String(),
	)

	userAmounts, err := u.getUserAmounts(ctx, sender)
	if err != nil {
		return "", err
	}
	if len(userAmounts) == 0 {
		logger.Info("no data")
		return "", nil
	}

	serializedTxs, txDetailItems, err := u.createDepositRawTransactions(ctx, sender, receiver, userAmounts)
	if err != nil {
		return "", err
	}
	if len(txDetailItems) == 0 {
		return "", nil
	}

	txID, err := u.updateDB(targetAction, txDetailItems, nil)
	if err != nil {
		return "", err
	}

	// save transaction result to file
	var generatedFileName string
	if len(serializedTxs) != 0 {
		generatedFileName, err = u.generateHexFile(targetAction, sender, txID, serializedTxs)
		if err != nil {
			return "", fmt.Errorf("fail to call generateHexFile(): %w", err)
		}
	}

	return generatedFileName, nil
}

// createPaymentTx creates unsigned tx for user (anonymous addresses)
// sender: payment, receiver: addresses coming from user_payment table
// - sender account (payment) covers fee, but this should be flexible
// Note:
// - to avoid complex logic to create raw transaction
// - only one address of sender should afford to send coin to all payment request users.
func (u *createTransactionUseCase) createPaymentTx(ctx context.Context) (string, error) {
	sender := u.paymentSender
	receiver := domainAccount.AccountTypeAnonymous
	targetAction := domainTx.ActionTypePayment
	logger.Debug("account",
		"sender", sender.String(),
		"receiver", receiver.String(),
	)

	// get payment data from payment_request
	userPayments, totalAmount, paymentRequestIds, err := u.createUserPayment()
	if err != nil {
		return "", err
	}
	if len(userPayments) == 0 {
		logger.Debug("no userPayments")
		// no data
		return "", nil
	}

	// check sender's total balance
	senderAddr, err := u.addrRepo.GetOneUnAllocated(sender)
	if err != nil {
		return "", fmt.Errorf("fail to call addrRepo.GetOneUnAllocated(): %w", err)
	}
	if err = u.validateAmount(ctx, senderAddr, totalAmount); err != nil {
		return "", nil
	}

	// create raw transaction for each address
	serializedTxs, txDetailItems := u.createPaymentRawTransactions(ctx, sender, receiver, userPayments, senderAddr)
	if len(txDetailItems) == 0 {
		return "", nil
	}

	txID, err := u.updateDB(targetAction, txDetailItems, paymentRequestIds)
	if err != nil {
		return "", err
	}

	// save transaction result to file
	var generatedFileName string
	if len(serializedTxs) != 0 {
		generatedFileName, err = u.generateHexFile(targetAction, sender, txID, serializedTxs)
		if err != nil {
			return "", fmt.Errorf("fail to call generateHexFile(): %w", err)
		}
	}

	return generatedFileName, nil
}

// createTransferTx creates unsigned tx for transfer coin among internal accounts except client, authorization
// FIXME: for now, receiver account covers fee, but this should be flexible
// - sender pays fee
// - any internal account should have only one address in XRP because no utxo
func (u *createTransactionUseCase) createTransferTx(
	ctx context.Context,
	sender, receiver domainAccount.AccountType,
	floatValue float64,
) (string, error) {
	targetAction := domainTx.ActionTypeTransfer

	// validation account
	if receiver == domainAccount.AccountTypeClient || receiver == domainAccount.AccountTypeAuthorization {
		return "", errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return "", errors.New("invalid account. sender and receiver is same")
	}

	// check sender's balance
	senderAddr, err := u.addrRepo.GetOneUnAllocated(sender)
	if err != nil {
		return "", fmt.Errorf("fail to call addrRepo.GetOneUnAllocated(sender): %w", err)
	}
	senderBalance, err := u.rippler.GetBalance(ctx, senderAddr.WalletAddress)
	if err != nil {
		return "", fmt.Errorf("fail to call rippler.GetBalance(): %w", err)
	}
	if senderBalance <= 20 {
		return "", errors.New("sender balance is insufficient to send")
	}
	if floatValue != 0 && senderBalance <= floatValue {
		return "", errors.New("sender balance is insufficient to send")
	}

	logger.Debug("amount",
		"floatValue", floatValue,
		"senderBalance", senderBalance,
	)

	// get receiver address
	receiverAddr, err := u.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return "", fmt.Errorf("fail to call addrRepo.GetOneUnAllocated(receiver): %w", err)
	}

	// call CreateRawTransaction
	instructions := &xrp.Instructions{
		MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
	}
	txJSON, rawTxString, err := u.rippler.CreateRawTransaction(
		ctx, senderAddr.WalletAddress, receiverAddr.WalletAddress, floatValue, instructions)
	if err != nil {
		return "", fmt.Errorf(
			"fail to call rippler.CreateRawTransaction(), sender address: %s: %w",
			senderAddr.WalletAddress, err)
	}
	logger.Debug("txJSON", "txJSON", txJSON)
	grok.Value(txJSON)

	// generate UUID to trace transaction because unsignedTx is not unique
	uid, err := u.uuidHandler.GenerateV7()
	if err != nil {
		return "", fmt.Errorf("fail to call uuidHandler.GenerateV7(): %w", err)
	}

	serializedTxs := []string{fmt.Sprintf("%s,%s", uid, rawTxString)}

	// create insert data for xrp_detail_tx
	txDetailItem := &models.XRPDetailTX{
		UUID:               uid.String(),
		CurrentTXType:      domainTx.TxTypeUnsigned.Int8(),
		SenderAccount:      sender.String(),
		SenderAddress:      senderAddr.WalletAddress,
		ReceiverAccount:    receiver.String(),
		ReceiverAddress:    receiverAddr.WalletAddress,
		Amount:             txJSON.Amount,
		XRPTXType:          txJSON.TransactionType,
		Fee:                txJSON.Fee,
		Flags:              txJSON.Flags,
		LastLedgerSequence: txJSON.LastLedgerSequence,
		Sequence:           txJSON.Sequence,
	}
	txDetailItems := []*models.XRPDetailTX{txDetailItem}

	txID, err := u.updateDB(targetAction, txDetailItems, nil)
	if err != nil {
		return "", err
	}

	// save transaction result to file
	var generatedFileName string
	if len(serializedTxs) != 0 {
		generatedFileName, err = u.generateHexFile(targetAction, sender, txID, serializedTxs)
		if err != nil {
			return "", fmt.Errorf("fail to call generateHexFile(): %w", err)
		}
	}

	return generatedFileName, nil
}

// getUserAmounts gets user amounts from addresses with balances
func (u *createTransactionUseCase) getUserAmounts(
	ctx context.Context,
	sender domainAccount.AccountType,
) ([]xrp.UserAmount, error) {
	// get addresses for sender account
	addrs, err := u.addrRepo.GetAll(sender)
	if err != nil {
		return nil, fmt.Errorf("fail to call addrRepo.GetAll(): %w", err)
	}

	// target addresses
	var userAmounts []xrp.UserAmount
	// address list for sender
	for _, addr := range addrs {
		// TODO: if previous tx is not done, wrong amount is returned. how to manage it??
		var balance float64
		balance, err = u.rippler.GetBalance(ctx, addr.WalletAddress)
		if err != nil {
			logger.Warn("fail to call rippler.GetBalance()",
				"address", addr.WalletAddress,
			)
		} else {
			logger.Debug("account_info",
				"address", addr.WalletAddress, "balance", balance)
			if balance != 0 {
				userAmounts = append(userAmounts, xrp.UserAmount{Address: addr.WalletAddress, Amount: balance})
			}
		}
	}
	return userAmounts, nil
}

// createDepositRawTransactions creates raw transactions for deposit
func (u *createTransactionUseCase) createDepositRawTransactions(
	ctx context.Context,
	sender, receiver domainAccount.AccountType,
	userAmounts []xrp.UserAmount,
) ([]string, []*models.XRPDetailTX, error) {
	// get address for deposit account
	depositAddr, err := u.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"fail to call addrRepo.GetOneUnAllocated(): %w", err,
		)
	}

	// create raw transaction for each address
	serializedTxs := make([]string, 0, len(userAmounts))
	txDetailItems := make([]*models.XRPDetailTX, 0, len(userAmounts))

	var sequence uint64
	for _, val := range userAmounts {
		// call CreateRawTransaction
		instructions := &xrp.Instructions{
			MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
		}
		if sequence != 0 {
			instructions.Sequence = sequence
		}
		var txJSON *xrp.TxInput
		var rawTxString string
		txJSON, rawTxString, err = u.rippler.CreateRawTransaction(
			ctx, val.Address, depositAddr.WalletAddress, 0, instructions)
		if err != nil {
			logger.Warn("fail to call rippler.CreateRawTransaction()", "error", err)
			continue
		}
		logger.Debug("txJSON", "txJSON", txJSON)
		grok.Value(txJSON)

		// sequence for next rawTransaction
		sequence = txJSON.Sequence + 1

		// generate UUID to trace transaction because unsignedTx is not unique
		uid, err := u.uuidHandler.GenerateV7()
		if err != nil {
			return nil, nil, fmt.Errorf("fail to call uuidHandler.GenerateV7(): %w", err)
		}

		serializedTxs = append(serializedTxs, fmt.Sprintf("%s,%s", uid, rawTxString))

		// create insert data for xrp_detail_tx
		txDetailItem := &models.XRPDetailTX{
			UUID:               uid.String(),
			CurrentTXType:      domainTx.TxTypeUnsigned.Int8(),
			SenderAccount:      sender.String(),
			SenderAddress:      val.Address,
			ReceiverAccount:    receiver.String(),
			ReceiverAddress:    depositAddr.WalletAddress,
			Amount:             txJSON.Amount,
			XRPTXType:          txJSON.TransactionType,
			Fee:                txJSON.Fee,
			Flags:              txJSON.Flags,
			LastLedgerSequence: txJSON.LastLedgerSequence,
			Sequence:           txJSON.Sequence,
		}
		txDetailItems = append(txDetailItems, txDetailItem)
	}

	return serializedTxs, txDetailItems, nil
}

// userPayment represents user's payment address and amount
type userPayment struct {
	senderAddr   string  // sender address for just checking
	receiverAddr string  // receiver address
	floatAmount  float64 // float amount (XRP)
}

// createUserPayment gets payment data from payment_request table
func (u *createTransactionUseCase) createUserPayment() ([]userPayment, float64, []int64, error) {
	// get payment_request
	paymentRequests, err := u.payReqRepo.GetAll()
	if err != nil {
		return nil, 0, nil, fmt.Errorf("fail to call payReqRepo.GetAll(): %w", err)
	}
	if len(paymentRequests) == 0 {
		logger.Debug("no data in payment_request")
		return nil, 0, nil, nil
	}

	userPayments := make([]userPayment, len(paymentRequests))
	paymentRequestIds := make([]int64, len(paymentRequests))
	var totalAmount float64

	// store `id` separately for key updating
	for idx, val := range paymentRequests {
		paymentRequestIds[idx] = val.ID

		userPayments[idx].senderAddr = val.SenderAddress
		userPayments[idx].receiverAddr = val.ReceiverAddress
		var amt float64
		amt, err = strconv.ParseFloat(val.Amount.String(), 64)
		if err != nil {
			// fatal error because table includes invalid data
			logger.Error("payment_request table includes invalid amount field")
			return nil, 0, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].floatAmount = amt

		// validate address
		if !xrp.ValidateAddress(userPayments[idx].receiverAddr) {
			// fatal error
			logger.Error("address is invalid",
				"address", userPayments[idx].receiverAddr,
				"error", err,
			)
			return nil, 0, nil, fmt.Errorf("address is invalid: %s: %w", userPayments[idx].receiverAddr, err)
		}

		// total amount
		totalAmount += amt
	}

	return userPayments, totalAmount, paymentRequestIds, nil
}

// validateAmount validates that sender has sufficient balance
func (u *createTransactionUseCase) validateAmount(
	ctx context.Context,
	senderAddr *models.Address,
	totalAmount float64,
) error {
	senderBalance, err := u.rippler.GetBalance(ctx, senderAddr.WalletAddress)
	if err != nil {
		return fmt.Errorf("fail to call rippler.GetBalance(): %w", err)
	}

	if senderBalance <= totalAmount {
		return errors.New("sender balance is insufficient to send")
	}
	return nil
}

// createPaymentRawTransactions creates raw transactions for payment
func (u *createTransactionUseCase) createPaymentRawTransactions(
	ctx context.Context,
	sender, receiver domainAccount.AccountType,
	userPayments []userPayment,
	senderAddr *models.Address,
) ([]string, []*models.XRPDetailTX) {
	serializedTxs := make([]string, 0, len(userPayments))
	txDetailItems := make([]*models.XRPDetailTX, 0, len(userPayments))
	var sequence uint64
	for _, userPayment := range userPayments {
		// call CreateRawTransaction
		instructions := &xrp.Instructions{
			MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
		}
		if sequence != 0 {
			instructions.Sequence = sequence
		}
		txJSON, rawTxString, err := u.rippler.CreateRawTransaction(
			ctx, senderAddr.WalletAddress, userPayment.receiverAddr, userPayment.floatAmount, instructions)
		if err != nil {
			// TODO: which is better to return err or continue?
			// return error in ethereum logic
			logger.Warn("fail to call rippler.CreateRawTransaction()", "error", err)
			continue
		}
		logger.Debug("txJSON", "txJSON", txJSON)
		grok.Value(txJSON)

		// sequence for next rawTransaction
		sequence = txJSON.Sequence + 1

		// generate UUID to trace transaction because unsignedTx is not unique
		uid, err := u.uuidHandler.GenerateV7()
		if err != nil {
			logger.Warn("fail to call uuidHandler.GenerateV7()", "error", err)
			continue
		}

		serializedTxs = append(serializedTxs, fmt.Sprintf("%s,%s", uid, rawTxString))

		// create insert data for xrp_detail_tx
		txDetailItem := &models.XRPDetailTX{
			UUID:               uid.String(),
			CurrentTXType:      domainTx.TxTypeUnsigned.Int8(),
			SenderAccount:      sender.String(),
			SenderAddress:      senderAddr.WalletAddress,
			ReceiverAccount:    receiver.String(),
			ReceiverAddress:    userPayment.receiverAddr,
			Amount:             txJSON.Amount,
			XRPTXType:          txJSON.TransactionType,
			Fee:                txJSON.Fee,
			Flags:              txJSON.Flags,
			LastLedgerSequence: txJSON.LastLedgerSequence,
			Sequence:           txJSON.Sequence,
		}
		txDetailItems = append(txDetailItems, txDetailItem)
	}
	return serializedTxs, txDetailItems
}

// updateDB updates database in a transaction
func (u *createTransactionUseCase) updateDB(
	targetAction domainTx.ActionType,
	txDetailItems []*models.XRPDetailTX,
	paymentRequestIds []int64,
) (int64, error) {
	// start transaction
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

	// Insert tx
	txID, err := u.txRepo.InsertUnsignedTx(targetAction)
	if err != nil {
		return 0, fmt.Errorf("fail to call txRepo.InsertUnsignedTx(): %w", err)
	}
	// Insert to xrp_detail_tx
	for idx := range txDetailItems {
		txDetailItems[idx].TXID = txID
	}
	if err = u.txDetailRepo.InsertBulk(txDetailItems); err != nil {
		return 0, fmt.Errorf("fail to call txDetailRepo.InsertBulk(): %w", err)
	}

	if targetAction == domainTx.ActionTypePayment {
		_, err = u.payReqRepo.UpdatePaymentID(txID, paymentRequestIds)
		if err != nil {
			return 0, fmt.Errorf("fail to call payReqRepo.UpdatePaymentID(): %w", err)
		}
	}
	return txID, nil
}

// generateHexFile generates file for hex txID and encoded previous addresses
func (u *createTransactionUseCase) generateHexFile(
	actionType domainTx.ActionType, senderAccount domainAccount.AccountType, txID int64, serializedTxs []string,
) (string, error) {
	// add senderAccount to first line
	serializedTxs = append([]string{senderAccount.String()}, serializedTxs...)

	// create file
	path := u.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeUnsigned, txID, 0)
	generatedFileName, err := u.txFileRepo.WriteFileSlice(path, serializedTxs)
	if err != nil {
		return "", fmt.Errorf("fail to call txFileRepo.WriteFileSlice(): %w", err)
	}

	return generatedFileName, nil
}
