package xrp

import (
	"context"
	"database/sql"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// TxMonitor type
type TxMonitor struct {
	xrp          ripple.Rippler
	dbConn       *sql.DB
	addrRepo     watch.AddressRepositorier
	txDetailRepo watch.XrpDetailTxRepositorier
	wtype        domainWallet.WalletType
}

// NewTxMonitor returns TxMonitor object
func NewTxMonitor(
	xrp ripple.Rippler,
	dbConn *sql.DB,
	addrRepo watch.AddressRepositorier,
	txDetailRepo watch.XrpDetailTxRepositorier,
	wtype domainWallet.WalletType,
) *TxMonitor {
	return &TxMonitor{
		xrp:          xrp,
		dbConn:       dbConn,
		addrRepo:     addrRepo,
		txDetailRepo: txDetailRepo,
		wtype:        wtype,
	}
}

// UpdateTxStatus update transaction status
// - no need for xrp
func (*TxMonitor) UpdateTxStatus() error {
	return nil
}

// MonitorBalance monitors balance
func (t *TxMonitor) MonitorBalance(_ uint64) error {
	targetAccounts := []domainAccount.AccountType{
		domainAccount.AccountTypeClient,
		domainAccount.AccountTypeDeposit,
		domainAccount.AccountTypePayment,
		domainAccount.AccountTypeStored,
	}

	for _, acnt := range targetAccounts {
		addrs, err := t.addrRepo.GetAllAddress(acnt)
		if err != nil {
			return fmt.Errorf("fail to call addrRepo.GetAllAddress(): %w", err)
		}
		total := t.xrp.GetTotalBalance(context.TODO(), addrs)
		logger.Info("total balance",
			"account", acnt.String(),
			"balance", total)
	}

	return nil
}
