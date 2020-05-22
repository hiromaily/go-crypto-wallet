package watchsrv

import (
	"database/sql"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
)

// TxMonitor type
type TxMonitor struct {
	eth          ethgrp.Ethereumer
	logger       *zap.Logger
	dbConn       *sql.DB
	txRepo       watchrepo.ETHTxRepositorier
	txDetailRepo watchrepo.EthDetailTxRepositorier
	payReqRepo   watchrepo.PaymentRequestRepositorier
	wtype        wallet.WalletType
}

// NewTxMonitor returns TxMonitor object
func NewTxMonitor(
	eth ethgrp.Ethereumer,
	logger *zap.Logger,
	dbConn *sql.DB,
	txRepo watchrepo.ETHTxRepositorier,
	txDetailRepo watchrepo.EthDetailTxRepositorier,
	payReqRepo watchrepo.PaymentRequestRepositorier,
	wtype wallet.WalletType) *TxMonitor {

	return &TxMonitor{
		eth:          eth,
		logger:       logger,
		dbConn:       dbConn,
		txRepo:       txRepo,
		txDetailRepo: txDetailRepo,
		payReqRepo:   payReqRepo,
		wtype:        wtype,
	}
}

// TODO: implementation

// UpdateTxStatus update transaction status
// - monitor transaction whose tx_type=3(TxTypeSent) in tx_payment/tx_deposit/tx_transfer
func (t *TxMonitor) UpdateTxStatus() error {
	//TODO: as possibility tx_type is not updated from `done`

	types := []action.ActionType{
		action.ActionTypeDeposit,
		action.ActionTypePayment,
		action.ActionTypeTransfer,
	}

	//1. update tx_type for TxTypeSent
	for _, actionType := range types {
		err := t.updateStatusTxTypeSent(actionType)
		if err != nil {
			return errors.Wrapf(err, "fail to call updateStatusTxTypeSent() ActionType: %s", actionType)
		}
	}

	//2. update tx_type for TxTypeDone
	// - TODO: notification
	for _, actionType := range types {
		err := t.updateStatusTxTypeDone(actionType)
		if err != nil {
			return errors.Wrapf(err, "fail to call updateStatusTxTypeDone() ActionType: %s", actionType)
		}
	}

	return nil
}

// update TxTypeSent to TxTypeDone if confirmation is 6 or more
func (t *TxMonitor) updateStatusTxTypeSent(actionType action.ActionType) error {
	//// get records whose status is TxTypeSent
	//hashes, err := t.txRepo.GetSentHashTx(actionType, tx.TxTypeSent)
	//if err != nil {
	//	return errors.Wrapf(err, "fail to call txRepo.GetSentHashTx(TxTypeSent) ActionType: %s", actionType)
	//}
	//
	//// get hash in detail and check confirmation
	//// update txType if confirmation is 6 or more (or configured number
	//for _, hash := range hashes {
	//	isDone, err := t.checkTxConfirmation(hash, actionType)
	//	if err != nil {
	//		t.logger.Error(
	//			"fail to call w.checkTransaction()",
	//			zap.String("actionType", actionType.String()),
	//			zap.String("hash", hash),
	//			zap.Error(err))
	//		continue
	//	}
	//	if isDone {
	//		// current confirmation meets 6 or more
	//		_, err = t.txRepo.UpdateTxTypeBySentHashTx(actionType, tx.TxTypeDone, hash)
	//		if err != nil {
	//			return errors.Wrapf(err, "fail to call repo.Tx().UpdateTxTypeBySentHashTx(tx.TxTypeDone) ActionType: %s", actionType)
	//		}
	//	}
	//}
	return nil
}

func (t *TxMonitor) updateStatusTxTypeDone(actionType action.ActionType) error {
	//// get records whose status is TxTypeDone
	//hashes, err := t.txRepo.GetSentHashTx(actionType, tx.TxTypeDone)
	//if err != nil {
	//	return errors.Wrapf(err, "fail to call txRepo.GetSentHashTx(TxTypeDone) ActionType: %s", actionType)
	//}
	//t.logger.Debug(
	//	"called repo.Tx().GetSentHashTx(TxTypeDone)",
	//	zap.String("actionType", actionType.String()),
	//	zap.Any("hashes", hashes))
	//
	//// notify tx get done
	//for _, hash := range hashes {
	//	txID, err := t.notifyTxDone(hash, actionType)
	//	if err != nil {
	//		t.logger.Error(
	//			"fail to call w.notifyUsers()",
	//			zap.String("actionType", actionType.String()),
	//			zap.String("hash", hash),
	//			zap.Error(err))
	//		continue
	//	}
	//	// update is already done
	//	if txID == 0 {
	//		continue
	//	}
	//
	//	// update tx_type to TxTypeNotified
	//	err = t.updateTxTypeNotified(txID, actionType)
	//	//TODO: even if update is failed, notification is done. so how to manage??
	//	if err != nil {
	//		t.logger.Error(
	//			"fail to call w.updateTxTypeNotified()",
	//			zap.String("actionType", actionType.String()),
	//			zap.String("hash", hash),
	//			zap.Error(err))
	//		continue
	//	}
	//}
	return nil
}

// checkTxConfirmation check confirmation for hash tx
//func (t *TxMonitor) checkTxConfirmation(hash string, actionType action.ActionType) (bool, error) {
//	// get tx in detail by RPC `gettransaction`
//	tran, err := t.btc.GetTransactionByTxID(hash)
//	if err != nil {
//		return false, errors.Wrapf(err, "fail to call btc.GetTransactionByTxID(): ActionType: %s, txID:%s", actionType, hash)
//	}
//	t.logger.Debug("confirmation detail",
//		zap.String("actionType", actionType.String()),
//		zap.Uint64("confirmation", tran.Confirmations))
//
//	// check current confirmation
//	if tran.Confirmations >= uint64(t.btc.ConfirmationBlock()) {
//		// current confirmation meets 6 or more
//		return true, nil
//	}
//
//	// not completed yet
//	//TODO: what if confirmation doesn't proceed for a long time after signed tx is sent
//	// - should it be canceled??
//	// - then raise fee and should unsigned tx be re-created again??
//	t.logger.Info("confirmation is not met yet",
//		zap.Uint64("want", t.btc.ConfirmationBlock()),
//		zap.Uint64("got", tran.Confirmations))
//
//	return false, nil
//}

// notifyTxDone notify tx is sent and met specific confirmation number
//func (t *TxMonitor) notifyTxDone(hash string, actionType action.ActionType) (int64, error) {
//
//	var (
//		txID int64
//		err  error
//	)
//
//	switch actionType {
//	case action.ActionTypeDeposit:
//		// 1. get txID from hash
//		txID, err = t.txRepo.GetTxIDBySentHash(actionType, hash)
//		if err != nil {
//			return 0, errors.Wrapf(err, "fail to call txRepo.GetTxIDBySentHash() ActionType: %s", actionType)
//		}
//
//		// 2. get txInputs
//		txInputs, err := t.txInputRepo.GetAllByTxID(txID)
//		if err != nil {
//			return 0, errors.Wrapf(err, "fail to call txInRepo.GetAllByTxID(%d) ActionType: %s", txID, actionType)
//		}
//		if len(txInputs) == 0 {
//			t.logger.Debug("txInputs is not found in tx_input table",
//				zap.Int64("tx_id", txID))
//			return 0, nil
//		}
//
//		// 3. notify to given input_addresses tx is done
//		// TODO:how to notify
//		for _, input := range txInputs {
//			t.logger.Debug("address in txInputs", zap.String("input.InputAddress", input.InputAddress))
//		}
//	case action.ActionTypePayment:
//		// 1. get txID from hash
//		txID, err = t.txRepo.GetTxIDBySentHash(actionType, hash)
//		if err != nil {
//			return 0, errors.Wrapf(err, "fail to call txRepo.GetTxIDBySentHash() ActionType: %s", actionType)
//		}
//
//		// 2. get info from payment_request table
//		paymentUsers, err := t.payReqRepo.GetAllByPaymentID(txID)
//		if err != nil {
//			return 0, errors.Wrapf(err, "fail to call repo.GetPaymentRequestByPaymentID(%d) ActionType: %s", txID, actionType)
//		}
//		if len(paymentUsers) == 0 {
//			t.logger.Debug("payment user is not found",
//				zap.Int64("tx_id", txID))
//			return 0, nil
//		}
//
//		// 3. notify to given input_addresses tx is done
//		// TODO:how to notify
//		for _, user := range paymentUsers {
//			t.logger.Debug("address in paymentUsers", zap.String("user.AddressFrom", user.SenderAddress))
//		}
//	case action.ActionTypeTransfer:
//		//TODO: not implemented yet
//		t.logger.Warn("action.ActionTypeTransfer is not implemented yet in notifyTxDone()")
//		return 0, errors.New("action.ActionTypeTransfer is not implemented yet in notifyTxDone()")
//	}
//
//	return txID, nil
//}

// update tx_type TxTypeNotified
//func (t *TxMonitor) updateTxTypeNotified(id int64, actionType action.ActionType) error {
//	switch actionType {
//	case action.ActionTypeDeposit:
//		_, err := t.txRepo.UpdateTxType(id, tx.TxTypeNotified)
//		if err != nil {
//			return errors.Wrapf(err, "fail to call repo.Tx().UpdateTxType(tx.TxTypeNotified) ActionType: %s", actionType)
//		}
//	case action.ActionTypePayment:
//		dtx, err := t.dbConn.Begin()
//		if err != nil {
//			return errors.Wrapf(err, "fail to start transaction")
//		}
//		defer func() {
//			if err != nil {
//				dtx.Rollback()
//			} else {
//				dtx.Commit()
//			}
//		}()
//		_, err = t.txRepo.UpdateTxType(id, tx.TxTypeNotified)
//		if err != nil {
//			return errors.Wrapf(err, "fail to call repo.Tx().UpdateTxType(tx.TxTypeNotified) ActionType: %s", actionType)
//		}
//
//		// update is_done=true in payment_request
//		_, err = t.payReqRepo.UpdateIsDone(id)
//		if err != nil {
//			return errors.Wrapf(err, "fail to call repo.UpdateIsDoneOnPaymentRequest() ActionType: %s", actionType)
//		}
//	case action.ActionTypeTransfer:
//		//TODO: not implemented yet, it could be same to action.ActionTypeDeposit
//		t.logger.Warn("action.ActionTypeTransfer is not implemented yet in updateTxTypeNotified()")
//		return errors.New("action.ActionTypeTransfer is not implemented yet in updateTxTypeNotified()")
//	}
//
//	return nil
//}
