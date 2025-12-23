package watchsrv

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

// TxMonitor type
type TxMonitor struct {
	eth          ethgrp.Ethereumer
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier
	txDetailRepo watchrepo.EthDetailTxRepositorier
	confirmNum   uint64
	wtype        wallet.WalletType
}

// NewTxMonitor returns TxMonitor object
func NewTxMonitor(
	eth ethgrp.Ethereumer,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	txDetailRepo watchrepo.EthDetailTxRepositorier,
	confirmNum uint64,
	wtype wallet.WalletType,
) *TxMonitor {
	return &TxMonitor{
		eth:          eth,
		dbConn:       dbConn,
		addrRepo:     addrRepo,
		txDetailRepo: txDetailRepo,
		confirmNum:   confirmNum,
		wtype:        wtype,
	}
}

// UpdateTxStatus update transaction status
// - monitor transaction whose tx_type=3(TxTypeSent) in tx_payment/tx_deposit/tx_transfer
func (t *TxMonitor) UpdateTxStatus() error {
	// update tx_type for TxTypeSent
	err := t.updateStatusTxTypeSent()
	if err != nil {
		return fmt.Errorf("fail to call updateStatusTxTypeSent(): %w", err)
	}

	// update tx_type for TxTypeDone
	// - TODO: notification
	// for _, actionType := range types {
	//	err := t.updateStatusTxTypeDone(actionType)
	//	if err != nil {
	//		return fmt.Errorf("fail to call updateStatusTxTypeDone() ActionType: %s: %w", actionType, err)
	//	}
	//}
	return nil
}

// update TxTypeSent to TxTypeDone if confirmation is 6 or more
func (t *TxMonitor) updateStatusTxTypeSent() error {
	// get records whose status is TxTypeSent
	hashes, err := t.txDetailRepo.GetSentHashTx(tx.TxTypeSent)
	if err != nil {
		return fmt.Errorf("fail to call txDetailRepo.GetSentHashTx(TxTypeSent): %w", err)
	}

	// get hash in detail and check confirmation
	for _, sentHash := range hashes {
		// check confirmation
		var confirmNum uint64
		confirmNum, err = t.eth.GetConfirmation(context.TODO(), sentHash)
		if err != nil {
			return fmt.Errorf("fail to call eth.GetConfirmation() sentHash: %s: %w", sentHash, err)
		}
		logger.Info("confirmation",
			"sentHash", sentHash,
			"confirmation num", confirmNum)
		if confirmNum < t.confirmNum {
			continue
		}
		// update status
		_, err = t.txDetailRepo.UpdateTxTypeBySentHashTx(tx.TxTypeDone, sentHash)
		if err != nil {
			logger.Warn("failed to call txDetailRepo.UpdateTxTypeBySentHashTx()",
				"error", err,
			)
		}
	}
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
		total, _ := t.eth.GetTotalBalance(context.TODO(), addrs)
		logger.Info("total balance",
			"account", acnt.String(),
			"balance", total.Uint64())
	}

	return nil
}
