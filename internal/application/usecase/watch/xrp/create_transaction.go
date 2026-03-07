package xrp

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bookerzzz/grok"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
	dbtx "github.com/hiromaily/go-crypto-wallet/pkg/db/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// userAmount holds an XRP address and its balance
type userAmount struct {
	Address string
	Amount  float64
}

// multisigTxEntry holds per-transaction data needed to build an XRPTransactionEntry
// for the multisig JSON file path.
type multisigTxEntry struct {
	uid           string
	txJSON        *dtoxrp.TxInput
	senderAddr    string
	senderAccType string
}

type createTransactionUseCase struct {
	accountInfo     apixrp.AccountInfoProvider
	txPreparer      apixrp.TransactionPreparer
	unitOfWork      dbtx.UnitOfWork
	uuidHandler     uuid.UUIDHandler
	addrRepo        repowatch.AddressRepositorier
	txRepo          repowatch.TxRepositorier
	txDetailRepo    repowatch.XRPDetailTXRepositorier
	payReqRepo      repowatch.PaymentRequestRepositorier
	txFileRepo      file.TransactionFileRepositorier
	depositReceiver domainAccount.AccountType
	paymentSender   domainAccount.AccountType
	network         string // "mainnet" or "testnet" (normalized for XRPTransactionFile)
}

// NewCreateTransactionUseCase creates a new CreateTransactionUseCase.
//
// This constructor follows the Interface Segregation Principle by depending on focused interfaces
// (AccountInfoProvider, TransactionPreparer) instead of the monolithic XRPer interface.
//
// Parameters:
//   - accountInfo: Provides account balance and information queries
//   - txPreparer: Prepares unsigned raw transactions
//   - network: XRP network type string (e.g. "mainnet", "testnet", "standalone"); used
//     when writing multisig JSON files. Any non-"mainnet" value is treated as "testnet".
//
// Note: Typically both interfaces are implemented by the same XRPer concrete type,
// but accepting them separately allows for better testability and clearer dependencies.
func NewCreateTransactionUseCase(
	accountInfo apixrp.AccountInfoProvider,
	txPreparer apixrp.TransactionPreparer,
	unitOfWork dbtx.UnitOfWork,
	uuidHandler uuid.UUIDHandler,
	addrRepo repowatch.AddressRepositorier,
	txRepo repowatch.TxRepositorier,
	txDetailRepo repowatch.XRPDetailTXRepositorier,
	payReqRepo repowatch.PaymentRequestRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	depositReceiver domainAccount.AccountType,
	paymentSender domainAccount.AccountType,
	network string,
) watchusecase.CreateTransactionUseCase {
	return &createTransactionUseCase{
		accountInfo:     accountInfo,
		txPreparer:      txPreparer,
		unitOfWork:      unitOfWork,
		uuidHandler:     uuidHandler,
		addrRepo:        addrRepo,
		txRepo:          txRepo,
		txDetailRepo:    txDetailRepo,
		payReqRepo:      payReqRepo,
		txFileRepo:      txFileRepo,
		depositReceiver: depositReceiver,
		paymentSender:   paymentSender,
		network:         network,
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
		fileName, execErr = u.createDepositTx(ctx, input.MultisigQuorum)
	case domainTx.ActionTypePayment:
		fileName, execErr = u.createPaymentTx(ctx, input.MultisigQuorum)
	case domainTx.ActionTypeTransfer:
		fileName, execErr = u.createTransferTx(
			ctx, input.SenderAccount, input.ReceiverAccount, input.Amount, input.MultisigQuorum)
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
func (u *createTransactionUseCase) createDepositTx(ctx context.Context, multisigQuorum uint32) (string, error) {
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

	serializedTxs, txDetailItems, msEntries, err := u.createDepositRawTransactions(
		ctx, sender, receiver, userAmounts, multisigQuorum)
	if err != nil {
		return "", err
	}
	if len(txDetailItems) == 0 {
		return "", nil
	}

	txID, err := u.updateDB(ctx, targetAction, txDetailItems, nil)
	if err != nil {
		return "", err
	}

	// save transaction result to file
	var generatedFileName string
	if multisigQuorum > 1 {
		if len(msEntries) != 0 {
			generatedFileName, err = u.generateMultisigJSONFile(targetAction, txID, multisigQuorum, msEntries)
			if err != nil {
				return "", fmt.Errorf("fail to call generateMultisigJSONFile(): %w", err)
			}
		}
	} else if len(serializedTxs) != 0 {
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
func (u *createTransactionUseCase) createPaymentTx(ctx context.Context, multisigQuorum uint32) (string, error) {
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
	senderAddr, err := u.getAndValidateAddress(sender, "sender")
	if err != nil {
		return "", err
	}
	if err = u.validateAmount(ctx, senderAddr, totalAmount); err != nil {
		return "", nil
	}

	// create raw transaction for each address
	serializedTxs, txDetailItems, msEntries := u.createPaymentRawTransactions(
		ctx, sender, receiver, userPayments, senderAddr, multisigQuorum)
	if len(txDetailItems) == 0 {
		return "", nil
	}

	txID, err := u.updateDB(ctx, targetAction, txDetailItems, paymentRequestIds)
	if err != nil {
		return "", err
	}

	// save transaction result to file
	var generatedFileName string
	if multisigQuorum > 1 {
		if len(msEntries) != 0 {
			generatedFileName, err = u.generateMultisigJSONFile(targetAction, txID, multisigQuorum, msEntries)
			if err != nil {
				return "", fmt.Errorf("fail to call generateMultisigJSONFile(): %w", err)
			}
		}
	} else if len(serializedTxs) != 0 {
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
// validateTransferAccounts validates that sender and receiver accounts are valid for transfer
func (*createTransactionUseCase) validateTransferAccounts(sender, receiver domainAccount.AccountType) error {
	if receiver == domainAccount.AccountTypeClient || receiver == domainAccount.AccountTypeAuthorization {
		return errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return errors.New("invalid account. sender and receiver is same")
	}
	return nil
}

// getAndValidateAddress retrieves an unallocated address and validates it's not nil
func (u *createTransactionUseCase) getAndValidateAddress(
	accountType domainAccount.AccountType,
	role string,
) (*domainAddress.Address, error) {
	addr, err := u.addrRepo.GetOneUnAllocated(accountType)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s address for account %s: %w", role, accountType.String(), err)
	}
	if addr == nil {
		return nil, fmt.Errorf("no unallocated address found for %s account %s", role, accountType.String())
	}
	return addr, nil
}

func (u *createTransactionUseCase) createTransferTx(
	ctx context.Context,
	sender, receiver domainAccount.AccountType,
	floatValue float64,
	multisigQuorum uint32,
) (string, error) {
	targetAction := domainTx.ActionTypeTransfer

	// validation account
	if err := u.validateTransferAccounts(sender, receiver); err != nil {
		return "", err
	}

	// get sender address and check balance
	senderAddr, err := u.getAndValidateAddress(sender, "sender")
	if err != nil {
		return "", err
	}

	senderBalance, err := u.accountInfo.GetBalance(ctx, senderAddr.WalletAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get balance for sender address %s: %w", senderAddr.WalletAddress, err)
	}
	if senderBalance <= 20 {
		return "", fmt.Errorf(
			"sender balance %.2f XRP is insufficient (minimum 20 XRP required for account reserve)",
			senderBalance,
		)
	}
	if floatValue != 0 && senderBalance <= floatValue {
		return "", fmt.Errorf("sender balance %.2f XRP is insufficient to send %.2f XRP", senderBalance, floatValue)
	}

	logger.Debug("amount",
		"floatValue", floatValue,
		"senderBalance", senderBalance,
	)

	// get receiver address
	receiverAddr, err := u.getAndValidateAddress(receiver, "receiver")
	if err != nil {
		return "", err
	}

	// call CreateRawTransaction
	instructions := &dtoxrp.Instructions{
		MaxLedgerVersionOffset: domainXRP.MaxLedgerVersionOffset,
	}
	txJSON, rawTxString, err := u.txPreparer.CreateRawTransaction(
		ctx, senderAddr.WalletAddress, receiverAddr.WalletAddress, floatValue, instructions)
	if err != nil {
		return "", fmt.Errorf(
			"failed to call txPreparer.CreateRawTransaction() for sender address %s: %w",
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
	txDetailItem, err := domainXRP.NewXRPDetailTx(
		0, // TxID will be set after insertion
		uid.String(),
		domainTx.TxTypeUnsigned,
		sender.String(),
		senderAddr.WalletAddress,
		receiver.String(),
		receiverAddr.WalletAddress,
		txJSON.Amount,
		txJSON.TransactionType,
		txJSON.Fee,
		txJSON.Flags,
		txJSON.LastLedgerSequence,
		txJSON.Sequence,
	)
	if err != nil {
		return "", fmt.Errorf("fail to create XRPDetailTx: %w", err)
	}
	txDetailItems := []*domainXRP.XRPDetailTx{txDetailItem}

	txID, err := u.updateDB(ctx, targetAction, txDetailItems, nil)
	if err != nil {
		return "", err
	}

	// save transaction result to file
	var generatedFileName string
	if multisigQuorum > 1 {
		msEntries := []multisigTxEntry{{
			uid:           uid.String(),
			txJSON:        txJSON,
			senderAddr:    senderAddr.WalletAddress,
			senderAccType: sender.String(),
		}}
		generatedFileName, err = u.generateMultisigJSONFile(targetAction, txID, multisigQuorum, msEntries)
		if err != nil {
			return "", fmt.Errorf("fail to call generateMultisigJSONFile(): %w", err)
		}
	} else if len(serializedTxs) != 0 {
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
) ([]userAmount, error) {
	// get addresses for sender account
	addrs, err := u.addrRepo.GetAll(sender)
	if err != nil {
		return nil, fmt.Errorf("fail to call addrRepo.GetAll(): %w", err)
	}

	// target addresses
	var userAmounts []userAmount
	// address list for sender
	for _, addr := range addrs {
		// TODO: if previous tx is not done, wrong amount is returned. how to manage it??
		var balance float64
		balance, err = u.accountInfo.GetBalance(ctx, addr.WalletAddress)
		if err != nil {
			logger.Warn("failed to call accountInfo.GetBalance()",
				"address", addr.WalletAddress,
			)
		} else {
			logger.Debug("account_info",
				"address", addr.WalletAddress, "balance", balance)
			if balance != 0 {
				userAmounts = append(userAmounts, userAmount{Address: addr.WalletAddress, Amount: balance})
			}
		}
	}
	return userAmounts, nil
}

// createDepositRawTransactions creates raw transactions for deposit.
// When multisigQuorum > 1, it also populates msEntries for the JSON file path.
func (u *createTransactionUseCase) createDepositRawTransactions(
	ctx context.Context,
	sender, receiver domainAccount.AccountType,
	userAmounts []userAmount,
	multisigQuorum uint32,
) ([]string, []*domainXRP.XRPDetailTx, []multisigTxEntry, error) {
	// get address for deposit account
	depositAddr, err := u.getAndValidateAddress(receiver, "deposit")
	if err != nil {
		return nil, nil, nil, err
	}

	// create raw transaction for each address
	serializedTxs := make([]string, 0, len(userAmounts))
	txDetailItems := make([]*domainXRP.XRPDetailTx, 0, len(userAmounts))
	var msEntries []multisigTxEntry

	var sequence uint64
	for _, val := range userAmounts {
		// call CreateRawTransaction
		instructions := &dtoxrp.Instructions{
			MaxLedgerVersionOffset: domainXRP.MaxLedgerVersionOffset,
		}
		if sequence != 0 {
			instructions.Sequence = sequence
		}
		var txJSON *dtoxrp.TxInput
		var rawTxString string
		txJSON, rawTxString, err = u.txPreparer.CreateRawTransaction(
			ctx, val.Address, depositAddr.WalletAddress, 0, instructions)
		if err != nil {
			logger.Warn("failed to call txPreparer.CreateRawTransaction()", "error", err)
			continue
		}
		logger.Debug("txJSON", "txJSON", txJSON)
		grok.Value(txJSON)

		// sequence for next rawTransaction
		sequence = txJSON.Sequence + 1

		// generate UUID to trace transaction because unsignedTx is not unique
		uid, err := u.uuidHandler.GenerateV7()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("fail to call uuidHandler.GenerateV7(): %w", err)
		}

		serializedTxs = append(serializedTxs, fmt.Sprintf("%s,%s", uid, rawTxString))

		if multisigQuorum > 1 {
			msEntries = append(msEntries, multisigTxEntry{
				uid:           uid.String(),
				txJSON:        txJSON,
				senderAddr:    val.Address,
				senderAccType: sender.String(),
			})
		}

		// create insert data for xrp_detail_tx
		txDetailItem, err := domainXRP.NewXRPDetailTx(
			0, // TxID will be set after insertion
			uid.String(),
			domainTx.TxTypeUnsigned,
			sender.String(),
			val.Address,
			receiver.String(),
			depositAddr.WalletAddress,
			txJSON.Amount,
			txJSON.TransactionType,
			txJSON.Fee,
			txJSON.Flags,
			txJSON.LastLedgerSequence,
			txJSON.Sequence,
		)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("fail to create XRPDetailTx: %w", err)
		}
		txDetailItems = append(txDetailItems, txDetailItem)
	}

	return serializedTxs, txDetailItems, msEntries, nil
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
		amt, err = strconv.ParseFloat(val.Amount, 64)
		if err != nil {
			// fatal error because table includes invalid data
			logger.Error("payment_request table includes invalid amount field")
			return nil, 0, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].floatAmount = amt

		// validate address
		if !xrpkg.ValidateAddress(userPayments[idx].receiverAddr) {
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
	senderAddr *domainAddress.Address,
	totalAmount float64,
) error {
	senderBalance, err := u.accountInfo.GetBalance(ctx, senderAddr.WalletAddress)
	if err != nil {
		return fmt.Errorf("failed to get balance for address %s: %w", senderAddr.WalletAddress, err)
	}

	if senderBalance <= totalAmount {
		return fmt.Errorf(
			"insufficient balance: sender has %.2f XRP, but %.2f XRP required",
			senderBalance, totalAmount,
		)
	}
	return nil
}

// createPaymentRawTransactions creates raw transactions for payment.
// When multisigQuorum > 1, it also populates msEntries for the JSON file path.
func (u *createTransactionUseCase) createPaymentRawTransactions(
	ctx context.Context,
	sender, receiver domainAccount.AccountType,
	userPayments []userPayment,
	senderAddr *domainAddress.Address,
	multisigQuorum uint32,
) ([]string, []*domainXRP.XRPDetailTx, []multisigTxEntry) {
	serializedTxs := make([]string, 0, len(userPayments))
	txDetailItems := make([]*domainXRP.XRPDetailTx, 0, len(userPayments))
	var msEntries []multisigTxEntry
	var sequence uint64
	for _, up := range userPayments {
		// call CreateRawTransaction
		instructions := &dtoxrp.Instructions{
			MaxLedgerVersionOffset: domainXRP.MaxLedgerVersionOffset,
		}
		if sequence != 0 {
			instructions.Sequence = sequence
		}
		txJSON, rawTxString, err := u.txPreparer.CreateRawTransaction(
			ctx, senderAddr.WalletAddress, up.receiverAddr, up.floatAmount, instructions)
		if err != nil {
			// TODO: which is better to return err or continue?
			// return error in ethereum logic
			logger.Warn("failed to call txPreparer.CreateRawTransaction()", "error", err)
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

		if multisigQuorum > 1 {
			msEntries = append(msEntries, multisigTxEntry{
				uid:           uid.String(),
				txJSON:        txJSON,
				senderAddr:    senderAddr.WalletAddress,
				senderAccType: sender.String(),
			})
		}

		// create insert data for xrp_detail_tx
		txDetailItem, err := domainXRP.NewXRPDetailTx(
			0, // TxID will be set after insertion
			uid.String(),
			domainTx.TxTypeUnsigned,
			sender.String(),
			senderAddr.WalletAddress,
			receiver.String(),
			up.receiverAddr,
			txJSON.Amount,
			txJSON.TransactionType,
			txJSON.Fee,
			txJSON.Flags,
			txJSON.LastLedgerSequence,
			txJSON.Sequence,
		)
		if err != nil {
			logger.Warn("fail to create XRPDetailTx", "error", err)
			continue
		}
		txDetailItems = append(txDetailItems, txDetailItem)
	}
	return serializedTxs, txDetailItems, msEntries
}

// updateDB updates database in a transaction
func (u *createTransactionUseCase) updateDB(
	ctx context.Context,
	targetAction domainTx.ActionType,
	txDetailItems []*domainXRP.XRPDetailTx,
	paymentRequestIds []int64,
) (int64, error) {
	// start transaction
	tx, err := u.unitOfWork.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("fail to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback() // Error already being handled
		} else {
			_ = tx.Commit() // Error already being handled
		}
	}()

	// Create transactional repositories that use the transaction
	txRepoWithTx, err := u.txRepo.WithTransaction(tx)
	if err != nil {
		return 0, fmt.Errorf("fail to create transactional txRepo: %w", err)
	}
	txDetailRepoWithTx, err := u.txDetailRepo.WithTransaction(tx)
	if err != nil {
		return 0, fmt.Errorf("fail to create transactional txDetailRepo: %w", err)
	}
	payReqRepoWithTx, err := u.payReqRepo.WithTransaction(tx)
	if err != nil {
		return 0, fmt.Errorf("fail to create transactional payReqRepo: %w", err)
	}

	// Insert tx
	txID, err := txRepoWithTx.InsertUnsignedTx(targetAction)
	if err != nil {
		return 0, fmt.Errorf("fail to call txRepo.InsertUnsignedTx(): %w", err)
	}
	// Insert to xrp_detail_tx
	for idx := range txDetailItems {
		txDetailItems[idx].TxID = txID
	}
	if err = txDetailRepoWithTx.InsertBulk(txDetailItems); err != nil {
		return 0, fmt.Errorf("fail to call txDetailRepo.InsertBulk(): %w", err)
	}

	if targetAction == domainTx.ActionTypePayment {
		_, err = payReqRepoWithTx.UpdatePaymentID(txID, paymentRequestIds)
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

// generateMultisigJSONFile generates an XRPTransactionFile JSON for multisig unsigned transactions.
// This is a parallel path to generateHexFile, invoked only when MultisigQuorum > 1.
func (u *createTransactionUseCase) generateMultisigJSONFile(
	actionType domainTx.ActionType,
	txID int64,
	quorum uint32,
	entries []multisigTxEntry,
) (string, error) {
	// Normalize network: only "mainnet" or "testnet" are valid in XRPTransactionFile
	network := u.network
	if network != string(xrpkg.NetworkTypeXRPMainNet) {
		network = string(xrpkg.NetworkTypeXRPTestNet)
	}

	txEntries := make([]dtoxrp.XRPTransactionEntry, len(entries))
	for i, e := range entries {
		txEntries[i] = dtoxrp.XRPTransactionEntry{
			UUID:               e.uid,
			UnsignedData:       *e.txJSON,
			SenderAccount:      e.senderAddr,
			SenderAccountType:  e.senderAccType,
			SignatureCount:     0,
			RequiredSignatures: int(quorum),
			SignedBlob:         nil,
			IsComplete:         false,
		}
	}

	txFile := &dtoxrp.XRPTransactionFile{
		Version:      "1.0.0",
		Chain:        "XRP",
		Network:      network,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		Transactions: txEntries,
	}

	path := u.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeUnsigned, txID, 0)
	generatedFileName, err := u.txFileRepo.WriteXRPJSONFile(path, txFile)
	if err != nil {
		return "", fmt.Errorf("fail to call txFileRepo.WriteXRPJSONFile(): %w", err)
	}

	return generatedFileName, nil
}
