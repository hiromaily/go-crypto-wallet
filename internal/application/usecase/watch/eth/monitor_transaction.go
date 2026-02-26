package eth

import (
	"context"
	"fmt"
	"time"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// maxReceiptRetries is the maximum number of retry attempts for node connectivity failures.
const maxReceiptRetries = 3

// initialRetryDelay is the base delay for exponential backoff.
const initialRetryDelay = time.Second

type monitorTransactionUseCase struct {
	ethClient    apieth.EtherTxMonitor
	addrRepo     repowatch.AddressRepositorier
	txDetailRepo repowatch.ETHDetailTXRepositorier
	confirmNum   uint64
}

// NewMonitorTransactionUseCase creates a new MonitorTransactionUseCase
func NewMonitorTransactionUseCase(
	ethClient apieth.EtherTxMonitor,
	addrRepo repowatch.AddressRepositorier,
	txDetailRepo repowatch.ETHDetailTXRepositorier,
	confirmNum uint64,
) watchusecase.MonitorTransactionUseCase {
	return &monitorTransactionUseCase{
		ethClient:    ethClient,
		addrRepo:     addrRepo,
		txDetailRepo: txDetailRepo,
		confirmNum:   confirmNum,
	}
}

func (u *monitorTransactionUseCase) UpdateTxStatus(ctx context.Context) error {
	err := u.updateStatusTxTypeSent(ctx)
	if err != nil {
		return fmt.Errorf("fail to call updateStatusTxTypeSent(): %w", err)
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

	for _, acnt := range targetAccounts {
		addrs, err := u.addrRepo.GetAllAddress(acnt)
		if err != nil {
			return fmt.Errorf("fail to call addrRepo.GetAllAddress(): %w", err)
		}
		total, _ := u.ethClient.GetTotalBalance(ctx, addrs)
		logger.Info("total balance",
			"account", acnt.String(),
			"balance", total.Uint64())
	}

	return nil
}

// updateStatusTxTypeSent processes all TxTypeSent transactions and updates their status.
func (u *monitorTransactionUseCase) updateStatusTxTypeSent(ctx context.Context) error {
	hashes, err := u.txDetailRepo.GetSentHashTx(domainTx.TxTypeSent)
	if err != nil {
		return fmt.Errorf("fail to call txDetailRepo.GetSentHashTx(TxTypeSent): %w", err)
	}

	for _, sentHash := range hashes {
		if err := u.processTransaction(ctx, sentHash); err != nil {
			logger.Warn("failed to process transaction",
				"sentHash", sentHash,
				"error", err,
			)
		}
	}
	return nil
}

// processTransaction handles status check and update for a single transaction.
// It checks the receipt for on-chain success/failure, counts confirmations,
// and transitions the transaction through TxTypeSent → TxTypeDone (or TxTypeCancel).
func (u *monitorTransactionUseCase) processTransaction(ctx context.Context, sentHash string) error {
	// Get receipt with exponential backoff retry for node connectivity failures.
	receipt, err := u.getReceiptWithRetry(ctx, sentHash)
	if err != nil {
		return fmt.Errorf("fail to get transaction receipt for %s: %w", sentHash, err)
	}
	if receipt == nil {
		// Transaction not yet mined; will be rechecked in the next monitoring cycle.
		logger.Info("transaction not yet mined, skipping", "sentHash", sentHash)
		return nil
	}

	// EVM status 0 means the transaction was reverted or failed on-chain.
	if receipt.Status == 0 {
		logger.Warn("transaction failed or reverted on-chain",
			"sentHash", sentHash,
			"revertReason", receipt.RevertReason,
		)
		if _, err := u.txDetailRepo.UpdateTxTypeBySentHashTx(domainTx.TxTypeCancel, sentHash); err != nil {
			logger.Warn("failed to update tx type to canceled",
				"sentHash", sentHash,
				"error", err,
			)
		}
		return nil
	}

	// Check confirmation depth against the configured threshold.
	confirmNum, err := u.ethClient.GetConfirmation(ctx, sentHash)
	if err != nil {
		return fmt.Errorf("fail to call eth.GetConfirmation() sentHash=%s: %w", sentHash, err)
	}
	logger.Info("confirmation",
		"sentHash", sentHash,
		"confirmationNum", confirmNum,
	)
	if confirmNum < u.confirmNum {
		return nil
	}

	// Confirmations sufficient — mark as done.
	if _, err := u.txDetailRepo.UpdateTxTypeBySentHashTx(domainTx.TxTypeDone, sentHash); err != nil {
		logger.Warn("failed to update tx type to done",
			"sentHash", sentHash,
			"error", err,
		)
		return nil
	}

	// Free the sender address to prevent double-spending on future transactions.
	if _, err := u.addrRepo.UpdateIsAllocated(false, receipt.From); err != nil {
		logger.Warn("failed to update is_allocated",
			"address", receipt.From,
			"error", err,
		)
	}

	return nil
}

// getReceiptWithRetry retrieves a transaction receipt with exponential backoff retry
// for node connectivity failures. Logs each failure at WARN level with the retry count.
// Returns (nil, nil) when the transaction has not yet been mined.
// Returns (nil, error) when all retry attempts are exhausted.
func (u *monitorTransactionUseCase) getReceiptWithRetry(
	ctx context.Context, txHash string,
) (*domainETH.TransactionReceipt, error) {
	for attempt := 0; attempt < maxReceiptRetries; attempt++ {
		receipt, err := u.ethClient.GetTxReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil // nil means "not yet mined"
		}

		logger.Warn("Ethereum node connectivity failure, retrying",
			"txHash", txHash,
			"attempt", attempt+1,
			"maxRetries", maxReceiptRetries,
			"error", err,
		)

		if attempt+1 == maxReceiptRetries {
			return nil, fmt.Errorf("failed to get receipt after %d attempts: %w", maxReceiptRetries, err)
		}

		// Exponential backoff: 1s, 2s, 4s, ...
		backoff := initialRetryDelay * (1 << uint(attempt))
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
		}
	}
	return nil, nil
}
