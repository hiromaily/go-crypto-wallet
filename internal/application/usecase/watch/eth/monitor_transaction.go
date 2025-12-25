package eth

import (
	"context"
	"fmt"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ethereum"
	watchrepo "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type monitorTransactionUseCase struct {
	ethClient    ethereum.Ethereumer
	addrRepo     watchrepo.AddressRepositorier
	txDetailRepo watchrepo.EthDetailTxRepositorier
	confirmNum   uint64
}

// NewMonitorTransactionUseCase creates a new MonitorTransactionUseCase
func NewMonitorTransactionUseCase(
	ethClient ethereum.Ethereumer,
	addrRepo watchrepo.AddressRepositorier,
	txDetailRepo watchrepo.EthDetailTxRepositorier,
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
	// update tx_type for TxTypeSent
	err := u.updateStatusTxTypeSent(ctx)
	if err != nil {
		return fmt.Errorf("fail to call updateStatusTxTypeSent(): %w", err)
	}

	// update tx_type for TxTypeDone
	// - TODO: notification
	// for _, actionType := range types {
	//	err := u.updateStatusTxTypeDone(actionType)
	//	if err != nil {
	//		return fmt.Errorf("fail to call updateStatusTxTypeDone() ActionType: %s: %w", actionType, err)
	//	}
	//}
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

// update TxTypeSent to TxTypeDone if confirmation is 6 or more
func (u *monitorTransactionUseCase) updateStatusTxTypeSent(ctx context.Context) error {
	// get records whose status is TxTypeSent
	hashes, err := u.txDetailRepo.GetSentHashTx(domainTx.TxTypeSent)
	if err != nil {
		return fmt.Errorf("fail to call txDetailRepo.GetSentHashTx(TxTypeSent): %w", err)
	}

	// get hash in detail and check confirmation
	for _, sentHash := range hashes {
		// check confirmation
		var confirmNum uint64
		confirmNum, err = u.ethClient.GetConfirmation(ctx, sentHash)
		if err != nil {
			return fmt.Errorf("fail to call eth.GetConfirmation() sentHash: %s: %w", sentHash, err)
		}
		logger.Info("confirmation",
			"sentHash", sentHash,
			"confirmation num", confirmNum)
		if confirmNum < u.confirmNum {
			continue
		}
		// update status
		_, err = u.txDetailRepo.UpdateTxTypeBySentHashTx(domainTx.TxTypeDone, sentHash)
		if err != nil {
			logger.Warn("failed to call txDetailRepo.UpdateTxTypeBySentHashTx()",
				"error", err,
			)
		}
	}
	return nil
}
