package xrp

import (
	"context"
	"fmt"

	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// xrpMonitorClient defines the interface for XRP operations needed by monitorTransactionUseCase.
// This follows the Interface Segregation Principle - depend only on what you need.
type xrpMonitorClient interface {
	apixrp.BalanceChecker
}

type monitorTransactionUseCase struct {
	xrper    xrpMonitorClient
	addrRepo repowatch.AddressRepositorier
}

// NewMonitorTransactionUseCase creates a new MonitorTransactionUseCase.
// The xrper parameter accepts any type that implements xrpMonitorClient (BalanceChecker).
// Typically, apixrp.XRPer is passed which implements all required methods.
func NewMonitorTransactionUseCase(
	xrper xrpMonitorClient,
	addrRepo repowatch.AddressRepositorier,
) watchusecase.MonitorTransactionUseCase {
	return &monitorTransactionUseCase{
		xrper:    xrper,
		addrRepo: addrRepo,
	}
}

// UpdateTxStatus updates transaction status
// Note: For XRP, UpdateTxStatus is a no-op (returns nil)
func (*monitorTransactionUseCase) UpdateTxStatus(ctx context.Context) error {
	// No need for XRP - transactions are validated immediately upon submission
	return nil
}

// MonitorBalance monitors balance across all account types
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
		total := u.xrper.GetTotalBalance(ctx, addrs)
		logger.Info("total balance",
			"account", acnt.String(),
			"balance", total)
	}

	return nil
}
