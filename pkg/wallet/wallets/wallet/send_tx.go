package wallet

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
)

// SendTx send signed tx by keygen/sign walet
func (w *Wallet) SendTx(filePath string) (string, error) {

	// get tx_receipt_id from file name
	//payment_5_unsigned_1_1534466246366489473
	actionType, _, txID, _, err := w.txFileRepo.ValidateFilePath(filePath, tx.TxTypeSigned)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.ValidateFilePath()")
	}

	w.logger.Debug("send_tx", zap.String("action_type", actionType.String()))

	// read hex from file
	signedHex, err := w.txFileRepo.ReadFile(filePath)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.ReadFile()")
	}

	// send signed tx
	hash, err := w.btc.SendTransactionByHex(signedHex)
	if err != nil {
		// if signature is not completed
		//-26: 16: mandatory-script-verify-flag-failed (Operation not valid with the current stack size)
		return "", errors.Wrap(err, "fail to call btc.SendTransactionByHex()")
	}

	if hash == nil {
		//tx is already sent
		return "", nil
	}

	// update tx_table
	affectedNum, err := w.repo.Tx().UpdateAfterTxSent(txID, tx.TxTypeSent, signedHex, hash.String())
	if err != nil {
		//TODO: even if error occurred, tx is already sent. so db should be corrected manually
		w.logger.Warn("fail to call repo.Tx().UpdateAfterTxSent() but tx is already sent. So database should be updated manually",
			zap.Int64("tx_id", txID),
			zap.String("tx_type", tx.TxTypeSent.String()),
			zap.Int8("tx_type_vakue", tx.TxTypeSent.Int8()),
			zap.String("signed_hex_tx", signedHex),
			zap.String("sent_hash_tx", hash.String()),
		)
		return "", errors.Wrapf(err, "fail to call updateHexForSentTx(), but tx is sent. txID: %d", txID)
	}
	if affectedNum == 0 {
		w.logger.Info("no records to update tx_table",
			zap.Int64("tx_id", txID),
			zap.String("tx_type", tx.TxTypeSent.String()),
			zap.Int8("tx_type_vakue", tx.TxTypeSent.Int8()),
			zap.String("signed_hex_tx", signedHex),
			zap.String("sent_hash_tx", hash.String()),
		)
		return "", nil
	}

	// update account_pubkey_table
	if actionType != action.ActionTypePayment {
		//skip for that receiver address is anonymous
		err = w.updateIsAllocatedAccountPubkey(txID)
		if err != nil {
			//TODO: even if error occurred, tx is already sent. so db should be corrected manually
			return "", err
		}
	}
	return hash.String(), nil
}

func (w *Wallet) updateIsAllocatedAccountPubkey(txID int64) error {
	// get txOutputs by tx_id
	txOutputs, err := w.repo.TxOutput().GetAllByTxID(txID)
	if err != nil {
		return errors.Wrap(err, "fail to call repo.TxOutput().GetAllByTxID()")
	}
	if len(txOutputs) == 0 {
		return errors.New("output tx could not be found in tx_receipt_output")
	}

	_, err = w.repo.Pubkey().UpdateIsAllocated(true, txOutputs[0].OutputAddress)
	if err != nil {
		return errors.Wrap(err, "fail to call repo.Pubkey().UpdateIsAllocated()")
	}

	return nil
}
