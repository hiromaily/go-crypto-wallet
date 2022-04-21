package watchsrv

import (
	"database/sql"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
)

// TxMonitor type
type TxMonitor struct {
	xrp          xrpgrp.Rippler
	logger       *zap.Logger
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier
	txDetailRepo watchrepo.XrpDetailTxRepositorier
	wtype        wallet.WalletType
}

// NewTxMonitor returns TxMonitor object
func NewTxMonitor(
	xrp xrpgrp.Rippler,
	logger *zap.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txDetailRepo watchrepo.XrpDetailTxRepositorier,
	wtype wallet.WalletType,
) *TxMonitor {
	return &TxMonitor{
		xrp:          xrp,
		logger:       logger,
		dbConn:       dbConn,
		addrRepo:     addrRepo,
		txDetailRepo: txDetailRepo,
		wtype:        wtype,
	}
}

// UpdateTxStatus update transaction status
// - no need for xrp
func (t *TxMonitor) UpdateTxStatus() error {
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
			return errors.Wrap(err, "fail to call addrRepo.GetAllAddress()")
		}
		total := t.xrp.GetTotalBalance(addrs)
		t.logger.Info("total balance",
			zap.String("account", acnt.String()),
			zap.Float64("balance", total))
	}

	return nil
}
