package watchsrv

import (
	"database/sql"
	"errors"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

// TxMonitor type
type TxMonitor struct {
	btc         btcgrp.Bitcoiner
	dbConn      *sql.DB
	txRepo      watchrepo.BTCTxRepositorier
	txInputRepo watchrepo.TxInputRepositorier
	payReqRepo  watchrepo.PaymentRequestRepositorier
	wtype       domainWallet.WalletType
}

// NewTxMonitor returns TxMonitor object
func NewTxMonitor(
	btc btcgrp.Bitcoiner,
	dbConn *sql.DB,
	txRepo watchrepo.BTCTxRepositorier,
	txInputRepo watchrepo.TxInputRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	wtype domainWallet.WalletType,
) *TxMonitor {
	return &TxMonitor{
		btc:         btc,
		dbConn:      dbConn,
		txRepo:      txRepo,
		txInputRepo: txInputRepo,
		payReqRepo:  payReqRepo,
		wtype:       wtype,
	}
}

// UpdateTxStatus update transaction status
// - monitor transaction whose tx_type=3(TxTypeSent) in tx_payment/tx_deposit/tx_transfer
func (t *TxMonitor) UpdateTxStatus() error {
	// TODO: as possibility tx_type is not updated from `done`

	types := []domainTx.ActionType{
		domainTx.ActionTypeDeposit,
		domainTx.ActionTypePayment,
		domainTx.ActionTypeTransfer,
	}

	// 1. update tx_type for TxTypeSent
	for _, actionType := range types {
		err := t.updateStatusTxTypeSent(actionType)
		if err != nil {
			return fmt.Errorf("fail to call updateStatusTxTypeSent() ActionType: %s: %w", actionType, err)
		}
	}

	// 2. update tx_type for TxTypeDone
	// - TODO: notification
	for _, actionType := range types {
		err := t.updateStatusTxTypeDone(actionType)
		if err != nil {
			return fmt.Errorf("fail to call updateStatusTxTypeDone() ActionType: %s: %w", actionType, err)
		}
	}

	return nil
}

// update TxTypeSent to TxTypeDone if confirmation is 6 or more
func (t *TxMonitor) updateStatusTxTypeSent(actionType domainTx.ActionType) error {
	// get records whose status is TxTypeSent
	hashes, err := t.txRepo.GetSentHashTx(actionType, domainTx.TxTypeSent)
	if err != nil {
		return fmt.Errorf("fail to call txRepo.GetSentHashTx(TxTypeSent) ActionType: %s: %w", actionType, err)
	}

	// get hash in detail and check confirmation
	// update txType if confirmation is 6 or more (or configured number
	for _, hash := range hashes {
		isDone, checkErr := t.checkTxConfirmation(hash, actionType)
		if checkErr != nil {
			logger.Error(
				"fail to call w.checkTransaction()",
				"actionType", actionType.String(),
				"hash", hash,
				"error", checkErr)
			continue
		}
		if isDone {
			// current confirmation meets 6 or more
			_, err = t.txRepo.UpdateTxTypeBySentHashTx(actionType, domainTx.TxTypeDone, hash)
			if err != nil {
				return fmt.Errorf(
					"fail to call repo.Tx().UpdateTxTypeBySentHashTx(domainTx.TxTypeDone) ActionType: %s: %w",
					actionType, err)
			}
		}
	}
	return nil
}

func (t *TxMonitor) updateStatusTxTypeDone(actionType domainTx.ActionType) error {
	// get records whose status is TxTypeDone
	hashes, err := t.txRepo.GetSentHashTx(actionType, domainTx.TxTypeDone)
	if err != nil {
		return fmt.Errorf("fail to call txRepo.GetSentHashTx(TxTypeDone) ActionType: %s: %w", actionType, err)
	}
	logger.Debug(
		"called repo.Tx().GetSentHashTx(TxTypeDone)",
		"actionType", actionType.String(),
		"hashes", hashes)

	// notify tx get done
	for _, hash := range hashes {
		txID, notifyErr := t.notifyTxDone(hash, actionType)
		if notifyErr != nil {
			logger.Error(
				"fail to call w.notifyUsers()",
				"actionType", actionType.String(),
				"hash", hash,
				"error", notifyErr)
			continue
		}
		// update is already done
		if txID == 0 {
			continue
		}

		// update tx_type to TxTypeNotified
		err = t.updateTxTypeNotified(txID, actionType)
		// TODO: even if update is failed, notification is done. so how to manage??
		if err != nil {
			logger.Error(
				"fail to call w.updateTxTypeNotified()",
				"actionType", actionType.String(),
				"hash", hash,
				"error", err)
			continue
		}
	}
	return nil
}

// checkTxConfirmation check confirmation for hash tx
func (t *TxMonitor) checkTxConfirmation(hash string, actionType domainTx.ActionType) (bool, error) {
	// get tx in detail by RPC `gettransaction`
	tran, err := t.btc.GetTransactionByTxID(hash)
	if err != nil {
		return false, fmt.Errorf(
			"fail to call btc.GetTransactionByTxID(): ActionType: %s, txID:%s: %w",
			actionType, hash, err)
	}
	logger.Debug("confirmation detail",
		"actionType", actionType.String(),
		"confirmation", tran.Confirmations)

	// check current confirmation
	if tran.Confirmations >= t.btc.ConfirmationBlock() {
		// current confirmation meets 6 or more
		return true, nil
	}

	// not completed yet
	// TODO: what if confirmation doesn't proceed for a long time after signed tx is sent
	// - should it be canceled??
	// - then raise fee and should unsigned tx be re-created again??
	logger.Info("confirmation is not met yet",
		"want", t.btc.ConfirmationBlock(),
		"got", tran.Confirmations)

	return false, nil
}

// notifyTxDone notify tx is sent and met specific confirmation number
func (t *TxMonitor) notifyTxDone(hash string, actionType domainTx.ActionType) (int64, error) {
	var (
		txID int64
		err  error
	)

	switch actionType {
	case domainTx.ActionTypeDeposit:
		// 1. get txID from hash
		txID, err = t.txRepo.GetTxIDBySentHash(actionType, hash)
		if err != nil {
			return 0, fmt.Errorf("fail to call txRepo.GetTxIDBySentHash() ActionType: %s: %w", actionType, err)
		}

		// 2. get txInputs
		var txInputs []*models.BTCTXInput
		txInputs, err = t.txInputRepo.GetAllByTxID(txID)
		if err != nil {
			return 0, fmt.Errorf("fail to call txInRepo.GetAllByTxID(%d) ActionType: %s: %w", txID, actionType, err)
		}
		if len(txInputs) == 0 {
			logger.Debug("txInputs is not found in tx_input table",
				"tx_id", txID)
			return 0, nil
		}

		// 3. notify to given input_addresses tx is done
		// TODO:how to notify
		for _, input := range txInputs {
			logger.Debug("address in txInputs", "input.InputAddress", input.InputAddress)
		}
	case domainTx.ActionTypePayment:
		// 1. get txID from hash
		txID, err = t.txRepo.GetTxIDBySentHash(actionType, hash)
		if err != nil {
			return 0, fmt.Errorf("fail to call txRepo.GetTxIDBySentHash() ActionType: %s: %w", actionType, err)
		}

		// 2. get info from payment_request table
		var paymentUsers []*models.PaymentRequest
		paymentUsers, err = t.payReqRepo.GetAllByPaymentID(txID)
		if err != nil {
			return 0, fmt.Errorf(
				"fail to call repo.GetPaymentRequestByPaymentID(%d) ActionType: %s: %w",
				txID, actionType, err)
		}
		if len(paymentUsers) == 0 {
			logger.Debug("payment user is not found",
				"tx_id", txID)
			return 0, nil
		}

		// 3. notify to given input_addresses tx is done
		// TODO:how to notify
		for _, user := range paymentUsers {
			logger.Debug("address in paymentUsers", "user.AddressFrom", user.SenderAddress)
		}
	case domainTx.ActionTypeTransfer:
		// TODO: not implemented yet
		logger.Warn("domainTx.ActionTypeTransfer is not implemented yet in notifyTxDone()")
		return 0, errors.New("domainTx.ActionTypeTransfer is not implemented yet in notifyTxDone()")
	}

	return txID, nil
}

// update tx_type TxTypeNotified
func (t *TxMonitor) updateTxTypeNotified(id int64, actionType domainTx.ActionType) error {
	switch actionType {
	case domainTx.ActionTypeDeposit:
		_, err := t.txRepo.UpdateTxType(id, domainTx.TxTypeNotified)
		if err != nil {
			return fmt.Errorf(
				"fail to call repo.Tx().UpdateTxType(domainTx.TxTypeNotified) ActionType: %s: %w",
				actionType, err)
		}
	case domainTx.ActionTypePayment:
		dtx, err := t.dbConn.Begin()
		if err != nil {
			return fmt.Errorf("fail to start transaction: %w", err)
		}
		defer func() {
			if err != nil {
				dtx.Rollback()
			} else {
				dtx.Commit()
			}
		}()
		_, err = t.txRepo.UpdateTxType(id, domainTx.TxTypeNotified)
		if err != nil {
			return fmt.Errorf(
				"fail to call repo.Tx().UpdateTxType(domainTx.TxTypeNotified) ActionType: %s: %w",
				actionType, err)
		}

		// update is_done=true in payment_request
		_, err = t.payReqRepo.UpdateIsDone(id)
		if err != nil {
			return fmt.Errorf("fail to call repo.UpdateIsDoneOnPaymentRequest() ActionType: %s: %w", actionType, err)
		}
	case domainTx.ActionTypeTransfer:
		// TODO: not implemented yet, it could be same to domainTx.ActionTypeDeposit
		logger.Warn("domainTx.ActionTypeTransfer is not implemented yet in updateTxTypeNotified()")
		return errors.New("domainTx.ActionTypeTransfer is not implemented yet in updateTxTypeNotified()")
	}

	return nil
}

// MonitorBalance monitors balances
func (t *TxMonitor) MonitorBalance(confirmationNum uint64) error {
	targetAccounts := []domainAccount.AccountType{
		domainAccount.AccountTypeClient,
		domainAccount.AccountTypeDeposit,
		domainAccount.AccountTypePayment,
		domainAccount.AccountTypeStored,
	}

	for _, acnt := range targetAccounts {
		total, err := t.btc.GetBalanceByAccount(acnt, confirmationNum)
		if err != nil {
			return fmt.Errorf("fail to call btc.GetBalanceByAccount() confirmation: %d: %w", confirmationNum, err)
		}
		logger.Info("total balance",
			"account", acnt.String(),
			"balance", total.String(),
		)
	}

	return nil
}
