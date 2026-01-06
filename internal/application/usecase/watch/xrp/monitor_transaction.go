package xrp

import (
	"context"
	"fmt"

	portsRipple "github.com/hiromaily/go-crypto-wallet/internal/application/ports/ripple"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	watchrepo "github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type monitorTransactionUseCase struct {
	rippler  portsRipple.Rippler
	addrRepo watchrepo.AddressRepositorier
}

// NewMonitorTransactionUseCase creates a new MonitorTransactionUseCase
func NewMonitorTransactionUseCase(
	rippler portsRipple.Rippler,
	addrRepo watchrepo.AddressRepositorier,
) watchusecase.MonitorTransactionUseCase {
	return &monitorTransactionUseCase{
		rippler:  rippler,
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
		total := u.rippler.GetTotalBalance(ctx, addrs)
		logger.Info("total balance",
			"account", acnt.String(),
			"balance", total)
	}

	return nil
}
