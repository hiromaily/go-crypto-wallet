package eth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	watchusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/ethtx"
	watchrepo "github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/serial"
)

type createTransactionUseCase struct {
	ethClient       ethereum.EtherTxCreator
	dbConn          *sql.DB
	addrRepo        watchrepo.AddressRepositorier
	txRepo          watchrepo.TxRepositorier
	txDetailRepo    watchrepo.EthDetailTxRepositorier
	payReqRepo      watchrepo.PaymentRequestRepositorier
	txFileRepo      file.TransactionFileRepositorier
	depositReceiver domainAccount.AccountType
	paymentSender   domainAccount.AccountType
}

// NewCreateTransactionUseCase creates a new CreateTransactionUseCase
func NewCreateTransactionUseCase(
	ethClient ethereum.EtherTxCreator,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txRepo watchrepo.TxRepositorier,
	txDetailRepo watchrepo.EthDetailTxRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	txFileRepo file.TransactionFileRepositorier,
	depositReceiver domainAccount.AccountType,
	paymentSender domainAccount.AccountType,
) watchusecase.CreateTransactionUseCase {
	return &createTransactionUseCase{
		ethClient:       ethClient,
		dbConn:          dbConn,
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

	var hex, fileName string
	var execErr error

	switch actionType {
	case domainTx.ActionTypeDeposit:
		hex, fileName, execErr = u.createDepositTx(ctx)
	case domainTx.ActionTypePayment:
		hex, fileName, execErr = u.createPaymentTx(ctx)
	case domainTx.ActionTypeTransfer:
		hex, fileName, execErr = u.createTransferTx(ctx, input.SenderAccount, input.ReceiverAccount, input.Amount)
	default:
		return watchusecase.CreateTransactionOutput{}, fmt.Errorf("unsupported action type: %s", input.ActionType)
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
func (u *createTransactionUseCase) createDepositTx(ctx context.Context) (string, string, error) {
	sender := domainAccount.AccountTypeClient
	receiver := u.depositReceiver
	targetAction := domainTx.ActionTypeDeposit
	logger.Debug("account",
		"sender", sender.String(),
		"receiver", receiver.String(),
	)

	userAmounts, err := u.getUserAmounts(ctx, sender)
	if err != nil {
		return "", "", err
	}
	if len(userAmounts) == 0 {
		logger.Info("no data")
		return "", "", nil
	}

	serializedTxs, txDetailItems, err := u.createDepositRawTransactions(ctx, sender, receiver, userAmounts)
	if err != nil {
		return "", "", err
	}
	if len(txDetailItems) == 0 {
		return "", "", nil
	}

	txID, err := u.updateDB(targetAction, txDetailItems, nil)
	logger.Debug("update result",
		"txID", txID,
		"error", err,
	)
	if err != nil {
		return "", "", err
	}

	// save transaction result to file
	var generatedFileName string
	if len(serializedTxs) != 0 {
		generatedFileName, err = u.generateHexFile(targetAction, sender, txID, serializedTxs)
		if err != nil {
			return "", "", fmt.Errorf("fail to call generateHexFile(): %w", err)
		}
	}

	return "", generatedFileName, nil
}

// createPaymentTx creates unsigned tx for user (anonymous addresses)
// sender: payment, receiver: addresses coming from payment_request table
// Note: only one address of sender should afford to send coin to all payment request users
func (u *createTransactionUseCase) createPaymentTx(ctx context.Context) (string, string, error) {
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
		return "", "", err
	}
	if len(userPayments) == 0 {
		logger.Debug("no data in userPayments")
		// no data
		return "", "", nil
	}

	// get sender address
	senderAddr, err := u.addrRepo.GetOneUnAllocated(sender)
	if err != nil {
		return "", "", fmt.Errorf("fail to call addrRepo.GetAll(domainAccount.AccountTypeClient): %w", err)
	}
	err = u.validateAmount(ctx, senderAddr, totalAmount)
	if err != nil {
		return "", "", err
	}

	// create raw transaction each address
	serializedTxs, txDetailItems, err := u.createPaymentRawTransactions(ctx, sender, receiver, userPayments, senderAddr)
	if err != nil {
		return "", "", err
	}
	if len(txDetailItems) == 0 {
		return "", "", nil
	}

	txID, err := u.updateDB(targetAction, txDetailItems, paymentRequestIds)
	if err != nil {
		return "", "", err
	}

	// save transaction result to file
	var generatedFileName string
	if len(serializedTxs) != 0 {
		generatedFileName, err = u.generateHexFile(targetAction, sender, txID, serializedTxs)
		if err != nil {
			return "", "", fmt.Errorf("fail to call generateHexFile(): %w", err)
		}
	}

	return "", generatedFileName, nil
}

// createTransferTx creates unsigned tx for transfer coin among internal accounts except client, authorization
// FIXME: for now, receiver account covers fee, but should be flexible
// - sender pays fee
// - any internal account should have only one address in Ethereum because no utxo
func (u *createTransactionUseCase) createTransferTx(
	ctx context.Context,
	sender, receiver domainAccount.AccountType,
	floatValue float64,
) (string, string, error) {
	targetAction := domainTx.ActionTypeTransfer

	// validation account
	if receiver == domainAccount.AccountTypeClient || receiver == domainAccount.AccountTypeAuthorization {
		return "", "", errors.New("invalid receiver account. client, authorization account is not allowed as receiver")
	}
	if sender == receiver {
		return "", "", errors.New("invalid account. sender and receiver is same")
	}

	// check sender's balance
	senderAddr, err := u.addrRepo.GetOneUnAllocated(sender)
	if err != nil {
		return "", "", fmt.Errorf("fail to call addrRepo.GetOneUnAllocated(sender): %w", err)
	}
	senderBalance, err := u.ethClient.GetBalance(ctx, senderAddr.WalletAddress, eth.QuantityTagLatest)
	if err != nil {
		return "", "", fmt.Errorf("fail to call eth.GetBalance(sender): %w", err)
	}

	if senderBalance.Uint64() == 0 {
		return "", "", errors.New("sender has no balance")
	}

	requiredValue := u.ethClient.FloatToBigInt(floatValue)
	if floatValue != 0 && (senderBalance.Uint64() <= requiredValue.Uint64()) {
		return "", "", errors.New("sender balance is insufficient to send")
	}
	logger.Debug("amount",
		"floatValue(Ether)", floatValue,
		"requiredValue(Ether)", requiredValue.Uint64(),
		"senderBalance", senderBalance.Uint64(),
	)

	// get receiver address
	receiverAddr, err := u.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return "", "", fmt.Errorf("fail to call addrRepo.GetOneUnAllocated(receiver): %w", err)
	}

	// call CreateRawTransaction
	rawTx, txDetailItem, err := u.ethClient.CreateRawTransaction(ctx,
		senderAddr.WalletAddress, receiverAddr.WalletAddress, requiredValue.Uint64(), 0)
	if err != nil {
		return "", "", fmt.Errorf(
			"fail to call eth.CreateRawTransaction(), sender address: %s: %w",
			senderAddr.WalletAddress, err)
	}

	rawTxHex := rawTx.TxHex
	logger.Debug("rawTxHex", "rawTxHex", rawTxHex)

	serializedTx, err := serial.EncodeToString(rawTx)
	if err != nil {
		return "", "", fmt.Errorf("fail to call serial.EncodeToString(rawTx): %w", err)
	}
	serializedTxs := []string{serializedTx}

	// create insert data for　eth_detail_tx
	txDetailItem.SenderAccount = sender.String()
	txDetailItem.ReceiverAccount = receiver.String()
	txDetailItems := []*models.EthDetailTX{txDetailItem}

	txID, err := u.updateDB(targetAction, txDetailItems, nil)
	if err != nil {
		return "", "", err
	}

	// save transaction result to file
	var generatedFileName string
	if len(serializedTxs) != 0 {
		generatedFileName, err = u.generateHexFile(targetAction, sender, txID, serializedTxs)
		if err != nil {
			return "", "", fmt.Errorf("fail to call generateHexFile(): %w", err)
		}
	}

	return "", generatedFileName, nil
}

// userPayment represents user's payment address and amount
type userPayment struct {
	senderAddr   string   // sender address for just checking
	receiverAddr string   // receiver address
	floatAmount  float64  // float amount (Ether)
	amount       *big.Int // amount (Wei)
}

func (u *createTransactionUseCase) getUserAmounts(
	ctx context.Context,
	sender domainAccount.AccountType,
) ([]eth.UserAmount, error) {
	// get addresses for client account
	addrs, err := u.addrRepo.GetAll(sender)
	if err != nil {
		return nil, fmt.Errorf("fail to call addrRepo.GetAll(domainAccount.AccountTypeClient): %w", err)
	}

	// target addresses
	var userAmounts []eth.UserAmount

	// address list for client
	for _, addr := range addrs {
		// TODO: if previous tx is not done, wrong amount is returned. how to manage it??
		var balance *big.Int
		balance, err = u.ethClient.GetBalance(ctx, addr.WalletAddress, eth.QuantityTagLatest)
		if err != nil {
			logger.Warn("fail to call .GetBalance()",
				"address", addr.WalletAddress,
				"error", err,
			)
		} else if balance.Uint64() != 0 {
			userAmounts = append(userAmounts, eth.UserAmount{Address: addr.WalletAddress, Amount: balance.Uint64()})
		}
	}

	return userAmounts, nil
}

func (u *createTransactionUseCase) createDepositRawTransactions(
	ctx context.Context,
	sender, receiver domainAccount.AccountType,
	userAmounts []eth.UserAmount,
) ([]string, []*models.EthDetailTX, error) {
	// get address for deposit account
	depositAddr, err := u.addrRepo.GetOneUnAllocated(receiver)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"fail to call addrRepo.GetOneUnAllocated(domainAccount.AccountTypeDeposit): %w", err,
		)
	}

	// create raw transaction each address
	serializedTxs := make([]string, 0, len(userAmounts))
	txDetailItems := make([]*models.EthDetailTX, 0, len(userAmounts))
	for _, val := range userAmounts {
		// call CreateRawTransaction
		var rawTx *ethtx.RawTx
		var txDetailItem *models.EthDetailTX
		rawTx, txDetailItem, err = u.ethClient.CreateRawTransaction(
			ctx, val.Address, depositAddr.WalletAddress, 0, 0)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"fail to call addrRepo.CreateRawTransaction(), sender address: %s: %w",
				val.Address, err)
		}

		rawTxHex := rawTx.TxHex
		logger.Debug("rawTxHex", "rawTxHex", rawTxHex)

		var serializedTx string
		serializedTx, err = serial.EncodeToString(rawTx)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to call serial.EncodeToString(rawTx): %w", err)
		}
		serializedTxs = append(serializedTxs, serializedTx)

		// create insert data for　eth_detail_tx
		txDetailItem.SenderAccount = sender.String()
		txDetailItem.ReceiverAccount = receiver.String()
		txDetailItems = append(txDetailItems, txDetailItem)
	}
	return serializedTxs, txDetailItems, nil
}

func (u *createTransactionUseCase) createUserPayment() ([]userPayment, *big.Int, []int64, error) {
	// get payment_request
	paymentRequests, err := u.payReqRepo.GetAll()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("fail to call repo.GetPaymentRequestAll(): %w", err)
	}
	if len(paymentRequests) == 0 {
		logger.Debug("no data in payment_request")
		return nil, nil, nil, nil
	}

	userPayments := make([]userPayment, len(paymentRequests))
	paymentRequestIds := make([]int64, len(paymentRequests))
	totalAmount := new(big.Int)

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
			return nil, nil, nil, errors.New("payment_request table includes invalid amount field")
		}
		userPayments[idx].floatAmount = amt

		// validate address
		if err = u.ethClient.ValidateAddr(userPayments[idx].receiverAddr); err != nil {
			// fatal error
			logger.Error("fail to call ValidationAddr",
				"address", userPayments[idx].receiverAddr,
				"error", err,
			)
			return nil, nil, nil, fmt.Errorf("fail to call eth.ValidationAddr(): %w", err)
		}

		// amount
		userPayments[idx].amount = u.ethClient.FloatToBigInt(userPayments[idx].floatAmount)
		totalAmount = new(big.Int).Add(totalAmount, userPayments[idx].amount)
	}

	return userPayments, totalAmount, paymentRequestIds, nil
}

func (u *createTransactionUseCase) validateAmount(
	ctx context.Context,
	senderAddr *models.Address,
	totalAmount *big.Int,
) error {
	// check sender's total balance
	senderBalance, err := u.ethClient.GetBalance(ctx, senderAddr.WalletAddress, eth.QuantityTagPending)
	if err != nil {
		return fmt.Errorf("fail to call eth.GetBalance(): %w", err)
	}

	if senderBalance.Uint64() <= totalAmount.Uint64() {
		return errors.New("sender balance is insufficient to send")
	}
	return nil
}

func (u *createTransactionUseCase) createPaymentRawTransactions(
	ctx context.Context,
	sender, receiver domainAccount.AccountType,
	userPayments []userPayment,
	senderAddr *models.Address,
) ([]string, []*models.EthDetailTX, error) {
	serializedTxs := make([]string, 0, len(userPayments))
	txDetailItems := make([]*models.EthDetailTX, 0, len(userPayments))
	additionalNonce := 0
	for _, userPayment := range userPayments {
		// call CreateRawTransaction
		rawTx, txDetailItem, err := u.ethClient.CreateRawTransaction(ctx,
			senderAddr.WalletAddress, userPayment.receiverAddr, userPayment.amount.Uint64(), additionalNonce)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"fail to call addrRepo.CreateRawTransaction(), sender address: %s: %w",
				senderAddr.WalletAddress, err)
		}
		additionalNonce++

		rawTxHex := rawTx.TxHex
		logger.Debug("rawTxHex", "rawTxHex", rawTxHex)

		serializedTx, err := serial.EncodeToString(rawTx)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to call serial.EncodeToString(rawTx): %w", err)
		}
		serializedTxs = append(serializedTxs, serializedTx)

		// create insert data for　eth_detail_tx
		txDetailItem.SenderAccount = sender.String()
		txDetailItem.ReceiverAccount = receiver.String()
		txDetailItems = append(txDetailItems, txDetailItem)
	}
	return serializedTxs, txDetailItems, nil
}

func (u *createTransactionUseCase) updateDB(
	targetAction domainTx.ActionType,
	txDetailItems []*models.EthDetailTX,
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

	// Insert eth_tx
	txID, err := u.txRepo.InsertUnsignedTx(targetAction)
	if err != nil {
		return 0, fmt.Errorf("fail to call txRepo.InsertUnsignedTx(): %w", err)
	}
	// Insert to eth_detail_tx
	for idx := range txDetailItems {
		txDetailItems[idx].TXID = txID
	}
	if err = u.txDetailRepo.InsertBulk(txDetailItems); err != nil {
		return 0, fmt.Errorf("fail to call txDetailRepo.InsertBulk(): %w", err)
	}

	if targetAction == domainTx.ActionTypePayment {
		_, err = u.payReqRepo.UpdatePaymentID(txID, paymentRequestIds)
		if err != nil {
			return 0, fmt.Errorf("fail to call repo.PayReq().UpdatePaymentID(txID, paymentRequestIds): %w", err)
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
		return "", fmt.Errorf("fail to call txFileRepo.WriteFile(): %w", err)
	}

	return generatedFileName, nil
}
