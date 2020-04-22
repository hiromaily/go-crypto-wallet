package wallet

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/action"
)

// UpdateTxStatus update transaction status
// - monitor transaction whose tx_type=3(TxTypeSent) in tx_payment/tx_receipt/tx_transfer
func (w *Wallet) UpdateTxStatus() error {
	//TODO: as possibility tx_type is not updated from `done`

	types := []action.ActionType{
		action.ActionTypeReceipt,
		action.ActionTypePayment,
		action.ActionTypeTransfer,
	}

	//1. update tx_type for TxTypeSent
	for _, actionType := range types {
		err := w.updateStatusForTxTypeSent(actionType)
		if err != nil {
			return errors.Wrapf(err, "fail to call updateStatusForTxTypeSent() ActionType: %s", actionType)
		}
	}

	//2. update tx_type for TxTypeDone
	// - TODO: notification
	for _, actionType := range types {
		err := w.updateStatusForTxTypeDone(actionType)
		if err != nil {
			return errors.Wrapf(err, "fail to call updateStatusForTxTypeDone() ActionType: %s", actionType)
		}
	}

	return nil
}

// update TxTypeSent to TxTypeDone if confirmation is 6 or more
func (w *Wallet) updateStatusForTxTypeSent(actionType action.ActionType) error {
	// get records whose status is TxTypeSent
	hashes, err := w.storager.GetSentTxHashByTxTypeSent(actionType)
	if err != nil {
		return errors.Wrapf(err, "fail to call storager.GetSentTxHashByTxTypeSent() ActionType: %s", actionType)
	}

	// get hash in detail and check confirmation
	// update txType if confirmation is 6 or more (or configured number
	for _, hash := range hashes {
		err = w.checkTxConfirmation(hash, actionType)
		if err != nil {
			w.logger.Error(
				"fail to call w.checkTransaction()",
				zap.String("actionType", actionType.String()),
				zap.String("hash", hash),
				zap.Error(err))
			continue
		}
	}
	return nil
}

func (w *Wallet) updateStatusForTxTypeDone(actionType action.ActionType) error {
	// get records whose status is TxTypeDone
	hashes, err := w.storager.GetSentTxHashByTxTypeDone(actionType)
	if err != nil {
		return errors.Wrapf(err, "fail to call storager.GetSentTxHashByTxTypeDone() ActionType: %s", actionType)
	}
	w.logger.Debug(
		"called storager.GetSentTxHashByTxTypeDone()",
		zap.String("actionType", actionType.String()),
		zap.Any("hashes", hashes))

	// notify tx get done
	for _, hash := range hashes {
		txID, err := w.notifyTxDone(hash, actionType)
		if err != nil {
			w.logger.Error(
				"fail to call w.notifyUsers()",
				zap.String("actionType", actionType.String()),
				zap.String("hash", hash),
				zap.Error(err))
			continue
		}
		// update is already done
		if txID == 0 {
			continue
		}

		// update tx_type to TxTypeNotified
		err = w.updateTxTypeNotified(txID, actionType)
		//TODO: even if update is failed, notification is done. so how to manage??
		if err != nil {
			w.logger.Error(
				"fail to call w.updateTxTypeNotified()",
				zap.String("actionType", actionType.String()),
				zap.String("hash", hash),
				zap.Error(err))
			continue
		}
	}
	return nil
}

// checkTxConfirmation check confirmation for hash tx
func (w *Wallet) checkTxConfirmation(hash string, actionType action.ActionType) error {
	// get tx in detail by RPC `gettransaction`
	tran, err := w.btc.GetTransactionByTxID(hash)
	if err != nil {
		return errors.Wrapf(err, "fail to call btc.GetTransactionByTxID(): ActionType: %s, txID:%s", actionType, hash)
	}
	w.logger.Debug("confirmation detail",
		zap.String("actionType", actionType.String()),
		zap.Uint64("confirmation", tran.Confirmations))

	// check current confirmation
	if tran.Confirmations >= uint64(w.btc.ConfirmationBlock()) {
		//current confirmation meet 6 or more
		_, err = w.storager.UpdateTxTypeDoneByTxHash(actionType, hash, nil, true)
		if err != nil {
			return errors.Wrapf(err, "fail to call storager.UpdateTxTypeDoneByTxHash() ActionType: %s", actionType)
		}
	} else {
		// not completed yet
		//TODO: what if confirmation doesn't proceed for a long time after signed tx is sent
		// - should it be canceled??
		// - then raise fee and should unsigned tx be re-created again??
		w.logger.Info("confirmation is not met yet",
			zap.Uint64("want", w.btc.ConfirmationBlock()),
			zap.Uint64("got", tran.Confirmations))
	}

	return nil
}

// notifyTxDone notify tx is sent and met specific confirmation number
func (w *Wallet) notifyTxDone(hash string, actionType action.ActionType) (int64, error) {

	var (
		txID int64
		err  error
	)

	switch actionType {
	case action.ActionTypeReceipt:
		// 1. get txID from hash
		txID, err = w.storager.GetTxIDBySentHash(actionType, hash)
		if err != nil {
			return 0, errors.Wrapf(err, "fail to call storager.GetTxIDBySentHash() ActionType: %s", actionType)
		}

		// 2. get txInputs
		txInputs, err := w.storager.GetTxInputByReceiptID(actionType, txID)
		if err != nil {
			return 0, errors.Wrapf(err, "fail to call storager.GetTxInputByReceiptID(%d) ActionType: %s", txID, actionType)
		}
		if len(txInputs) == 0 {
			w.logger.Debug("txInputs is not found in tx_input table",
				zap.Int64("tx_id", txID))
			return 0, nil
		}

		// 3. notify to given input_addresses tx is done
		// TODO:how to notify
		for _, input := range txInputs {
			w.logger.Debug("address in txInputs", zap.String("input.InputAddress", input.InputAddress))
		}
	case action.ActionTypePayment:
		// 1. get txID from hash
		txID, err = w.storager.GetTxIDBySentHash(actionType, hash)
		if err != nil {
			return 0, errors.Wrapf(err, "fail to call storager.GetTxIDBySentHash() ActionType: %s", actionType)
		}

		// 2. get info from payment_request table
		paymentUsers, err := w.storager.GetPaymentRequestByPaymentID(txID)
		if err != nil {
			return 0, errors.Wrapf(err, "fail to call storager.GetPaymentRequestByPaymentID(%d) ActionType: %s", txID, actionType)
		}
		if len(paymentUsers) == 0 {
			w.logger.Debug("payment user is not found",
				zap.Int64("tx_id", txID))
			return 0, nil
		}

		// 3. notify to given input_addresses tx is done
		// TODO:how to notify
		for _, user := range paymentUsers {
			w.logger.Debug("address in paymentUsers", zap.String("user.AddressFrom", user.AddressFrom))
		}
	case action.ActionTypeTransfer:
		//TODO: not implemented yet
		w.logger.Warn("action.ActionTypeTransfer is not implemented yet in notifyTxDone()")
		return 0, errors.New("action.ActionTypeTransfer is not implemented yet in notifyTxDone()")
	}

	return txID, nil
}

//upadte tx_type TxTypeNotified
func (w *Wallet) updateTxTypeNotified(id int64, actionType action.ActionType) error {
	switch actionType {
	case action.ActionTypeReceipt:
		_, err := w.storager.UpdateTxTypeNotifiedByID(actionType, id, nil, true)
		if err != nil {
			return errors.Wrapf(err, "fail to call storager.UpdateTxTypeNotifiedByID() ActionType: %s", actionType)
		}
	case action.ActionTypePayment:
		tx := w.storager.MustBegin()
		_, err := w.storager.UpdateTxTypeNotifiedByID(actionType, id, tx, false)
		if err != nil {
			return errors.Wrapf(err, "fail to call storager.UpdateTxTypeNotifiedByID() ActionType: %s", actionType)
		}

		// update is_done=true in payment_request
		_, err = w.storager.UpdateIsDoneOnPaymentRequest(id, tx, true)
		if err != nil {
			return errors.Wrapf(err, "fail to call storager.UpdateIsDoneOnPaymentRequest() ActionType: %s", actionType)
		}
	case action.ActionTypeTransfer:
		//TODO: not implemented yet
		w.logger.Warn("action.ActionTypeTransfer is not implemented yet in updateTxTypeNotified()")
		return errors.New("action.ActionTypeTransfer is not implemented yet in updateTxTypeNotified()")
	}

	return nil
}
