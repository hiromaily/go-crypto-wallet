package watchsrv

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
)

// TxMonitor type
type TxMonitor struct {
	xrp          xrpgrp.Rippler
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier
	txDetailRepo watchrepo.XrpDetailTxRepositorier
	wtype        wallet.WalletType
}

// NewTxMonitor returns TxMonitor object
func NewTxMonitor(
	xrp xrpgrp.Rippler,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txDetailRepo watchrepo.XrpDetailTxRepositorier,
	wtype wallet.WalletType,
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
	targetAccounts := []account.AccountType{
		account.AccountTypeClient,
		account.AccountTypeDeposit,
		account.AccountTypePayment,
		account.AccountTypeStored,
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
