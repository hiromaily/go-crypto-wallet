package btc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	watchrepo "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type monitorTransactionUseCase struct {
	btcClient   bitcoin.Bitcoiner
	dbConn      *sql.DB
	txRepo      watchrepo.BTCTxRepositorier
	txInputRepo watchrepo.TxInputRepositorier
	payReqRepo  watchrepo.PaymentRequestRepositorier
}

// NewMonitorTransactionUseCase creates a new MonitorTransactionUseCase
func NewMonitorTransactionUseCase(
	btcClient bitcoin.Bitcoiner,
	dbConn *sql.DB,
	txRepo watchrepo.BTCTxRepositorier,
	txInputRepo watchrepo.TxInputRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
) watchusecase.MonitorTransactionUseCase {
	return &monitorTransactionUseCase{
		btcClient:   btcClient,
		dbConn:      dbConn,
		txRepo:      txRepo,
		txInputRepo: txInputRepo,
		payReqRepo:  payReqRepo,
	}
}

func (u *monitorTransactionUseCase) UpdateTxStatus(ctx context.Context) error {
	types := []domainTx.ActionType{
		domainTx.ActionTypeDeposit,
		domainTx.ActionTypePayment,
		domainTx.ActionTypeTransfer,
	}

	// 1. Update transactions from Sent → Done (when confirmations meet threshold)
	for _, actionType := range types {
		if err := u.updateStatusFromSentToDone(actionType); err != nil {
			return fmt.Errorf("failed to update status to done for %s: %w", actionType, err)
		}
	}

	// 2. Update transactions from Done → Notified (notify users and mark as notified)
	for _, actionType := range types {
		if err := u.updateStatusFromDoneToNotified(actionType); err != nil {
			return fmt.Errorf("failed to update status to notified for %s: %w", actionType, err)
		}
	}

	return nil
}

func (u *monitorTransactionUseCase) MonitorBalance(
	ctx context.Context,
	input watchusecase.MonitorBalanceInput,
) error {
	targetAccounts := []domainAccount.AccountType{
		domainAccount.AccountTypeClient,
		domainAccount.AccountTypeDeposit,
		domainAccount.AccountTypePayment,
		domainAccount.AccountTypeStored,
	}

	for _, account := range targetAccounts {
		balance, err := u.btcClient.GetBalanceByAccount(account, input.ConfirmationNum)
		if err != nil {
			return fmt.Errorf("failed to get balance for %s: %w", account, err)
		}

		logger.Info("account balance",
			"account", account.String(),
			"balance", balance.String(),
			"confirmations", input.ConfirmationNum)
	}

	return nil
}

// updateStatusFromSentToDone updates transactions from Sent to Done when confirmations are met
func (u *monitorTransactionUseCase) updateStatusFromSentToDone(actionType domainTx.ActionType) error {
	// Get transactions with Sent status
	hashes, err := u.txRepo.GetSentHashTx(actionType, domainTx.TxTypeSent)
	if err != nil {
		return fmt.Errorf("failed to get sent transactions: %w", err)
	}

	// Check confirmation for each transaction
	for _, hash := range hashes {
		isDone, err := u.checkTransactionConfirmation(hash, actionType)
		if err != nil {
			logger.Error("failed to check transaction confirmation",
				"action_type", actionType.String(),
				"hash", hash,
				"error", err)
			continue
		}

		if isDone {
			// Update status to Done
			_, err = u.txRepo.UpdateTxTypeBySentHashTx(actionType, domainTx.TxTypeDone, hash)
			if err != nil {
				return fmt.Errorf("failed to update tx to done status: %w", err)
			}
			logger.Info("transaction status updated to done",
				"action_type", actionType.String(),
				"hash", hash)
		}
	}

	return nil
}

// updateStatusFromDoneToNotified notifies users and updates status from Done to Notified
func (u *monitorTransactionUseCase) updateStatusFromDoneToNotified(actionType domainTx.ActionType) error {
	// Get transactions with Done status
	hashes, err := u.txRepo.GetSentHashTx(actionType, domainTx.TxTypeDone)
	if err != nil {
		return fmt.Errorf("failed to get done transactions: %w", err)
	}

	logger.Debug("checking done transactions",
		"action_type", actionType.String(),
		"count", len(hashes))

	// Notify for each transaction
	for _, hash := range hashes {
		txID, err := u.notifyTransactionDone(hash, actionType)
		if err != nil {
			logger.Error("failed to notify transaction done",
				"action_type", actionType.String(),
				"hash", hash,
				"error", err)
			continue
		}

		// Skip if already notified
		if txID == 0 {
			continue
		}

		// Update status to Notified
		if err := u.updateToNotifiedStatus(txID, actionType); err != nil {
			logger.Error("failed to update to notified status",
				"action_type", actionType.String(),
				"tx_id", txID,
				"error", err)
			continue
		}

		logger.Info("transaction notified",
			"action_type", actionType.String(),
			"tx_id", txID,
			"hash", hash)
	}

	return nil
}

// checkTransactionConfirmation checks if transaction has enough confirmations
func (u *monitorTransactionUseCase) checkTransactionConfirmation(
	hash string,
	actionType domainTx.ActionType,
) (bool, error) {
	// Get transaction details from Bitcoin network
	tx, err := u.btcClient.GetTransactionByTxID(hash)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction details: %w", err)
	}

	logger.Debug("transaction confirmation status",
		"action_type", actionType.String(),
		"hash", hash,
		"confirmations", tx.Confirmations,
		"required", u.btcClient.ConfirmationBlock())

	// Check if confirmations meet threshold
	if tx.Confirmations >= u.btcClient.ConfirmationBlock() {
		return true, nil
	}

	// Not enough confirmations yet
	logger.Info("waiting for more confirmations",
		"hash", hash,
		"current", tx.Confirmations,
		"required", u.btcClient.ConfirmationBlock())

	return false, nil
}

// notifyTransactionDone notifies relevant parties that transaction is confirmed
func (u *monitorTransactionUseCase) notifyTransactionDone(
	hash string,
	actionType domainTx.ActionType,
) (int64, error) {
	// Get transaction ID
	txID, err := u.txRepo.GetTxIDBySentHash(actionType, hash)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction ID: %w", err)
	}

	switch actionType {
	case domainTx.ActionTypeDeposit:
		return u.notifyDepositTransaction(txID)
	case domainTx.ActionTypePayment:
		return u.notifyPaymentTransaction(txID)
	case domainTx.ActionTypeTransfer:
		logger.Warn("transfer notification not implemented yet")
		return 0, errors.New("transfer transaction notification not implemented")
	default:
		return 0, fmt.Errorf("unknown action type: %s", actionType)
	}
}

// notifyDepositTransaction notifies about deposit transaction
func (u *monitorTransactionUseCase) notifyDepositTransaction(txID int64) (int64, error) {
	// Get transaction inputs
	txInputs, err := u.txInputRepo.GetAllByTxID(txID)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction inputs: %w", err)
	}

	if len(txInputs) == 0 {
		logger.Debug("no transaction inputs found", "tx_id", txID)
		return 0, nil
	}

	// Notify affected addresses (TODO: implement actual notification mechanism)
	for _, input := range txInputs {
		logger.Debug("deposit transaction input address",
			"tx_id", txID,
			"address", input.InputAddress)
		// TODO: Send notification to address owner
	}

	return txID, nil
}

// notifyPaymentTransaction notifies about payment transaction
func (u *monitorTransactionUseCase) notifyPaymentTransaction(txID int64) (int64, error) {
	// Get payment requests
	paymentRequests, err := u.payReqRepo.GetAllByPaymentID(txID)
	if err != nil {
		return 0, fmt.Errorf("failed to get payment requests: %w", err)
	}

	if len(paymentRequests) == 0 {
		logger.Debug("no payment requests found", "tx_id", txID)
		return 0, nil
	}

	// Notify payment recipients (TODO: implement actual notification mechanism)
	for _, req := range paymentRequests {
		logger.Debug("payment transaction recipient",
			"tx_id", txID,
			"sender_address", req.SenderAddress)
		// TODO: Send notification to payment recipient
	}

	return txID, nil
}

// updateToNotifiedStatus updates transaction status to Notified
func (u *monitorTransactionUseCase) updateToNotifiedStatus(txID int64, actionType domainTx.ActionType) error {
	switch actionType {
	case domainTx.ActionTypeDeposit:
		_, err := u.txRepo.UpdateTxType(txID, domainTx.TxTypeNotified)
		if err != nil {
			return fmt.Errorf("failed to update tx type to notified: %w", err)
		}

	case domainTx.ActionTypePayment:
		// Use database transaction for payment (updates both tx and payment_request tables)
		dtx, err := u.dbConn.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer func() {
			if err != nil {
				_ = dtx.Rollback()
			} else {
				_ = dtx.Commit()
			}
		}()

		// Update transaction type
		_, err = u.txRepo.UpdateTxType(txID, domainTx.TxTypeNotified)
		if err != nil {
			return fmt.Errorf("failed to update tx type to notified: %w", err)
		}

		// Mark payment request as done
		_, err = u.payReqRepo.UpdateIsDone(txID)
		if err != nil {
			return fmt.Errorf("failed to update payment request: %w", err)
		}

	case domainTx.ActionTypeTransfer:
		logger.Warn("transfer status update not implemented yet")
		return errors.New("transfer transaction status update not implemented")

	default:
		return fmt.Errorf("unknown action type: %s", actionType)
	}

	return nil
}
