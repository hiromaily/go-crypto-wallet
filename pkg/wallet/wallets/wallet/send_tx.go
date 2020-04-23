package wallet

import (
	"time"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
)

// SendTx send signed tx by keygen/sign walet
func (w *Wallet) SendTx(filePath string) (string, error) {

	// get tx_receipt_id from file name
	//payment_5_unsigned_1_1534466246366489473
	actionType, _, txReceiptID, _, err := w.txFileRepo.ValidateFilePath(filePath, tx.TxTypeSigned)
	if err != nil {
		return "", errors.Wrap(err, "fail to call txFileRepo.ValidateFilePath()")
	}

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

	// update tx_table
	err = w.updateHexForSentTx(txReceiptID, signedHex, hash.String(), actionType)
	if err != nil {
		//TODO: even if error occurred, tx is already sent. so db should be corrected manually
		return "", errors.Wrap(err, "fail to call updateHexForSentTx(), but tx is sent")
	}

	// update account_pubkey_table
	err = w.updateIsAllocatedForAccountPubkey(txReceiptID, actionType)
	if err != nil {
		//TODO: even if error occurred, tx is already sent. so db should be corrected manually
		return "", errors.Wrap(err, "fail to call updateIsAllocatedForAccountPubkey()")
	}

	return hash.String(), nil
}

func (w *Wallet) updateHexForSentTx(txReceiptID int64, signedHex, sentTxID string, actionType action.ActionType) error {
	// 1.TxReceipt table
	t := time.Now()
	txReceipt := walletrepo.TxTable{}
	txReceipt.ID = txReceiptID
	txReceipt.SignedHexTx = signedHex
	txReceipt.SentHashTx = sentTxID
	txReceipt.SentUpdatedAt = &t
	txReceipt.TxType = tx.TxTypeValue[tx.TxTypeSent] //3:未署名

	var (
		affectedNum int64
		err         error
	)

	affectedNum, err = w.repo.UpdateTxAfterSent(actionType, &txReceipt, nil, true)

	if err != nil {
		return errors.Wrap(err, "fail to call storager.UpdateTxAfterSent()")
	}
	if affectedNum == 0 {
		return errors.New("tx_table was not updated by storager.UpdateTxAfterSent()")
	}

	return nil
}

func (w *Wallet) updateIsAllocatedForAccountPubkey(txReceiptID int64, actionType action.ActionType) error {
	if actionType == action.ActionTypeReceipt {
		return nil
	}

	// get txOutputs from .tx_receipt_output by receipt_id
	txOutputs, err := w.repo.GetTxOutputByReceiptID(actionType, txReceiptID)
	if err != nil {
		return errors.Wrap(err, "fail to call storager.GetTxOutputByReceiptID()")
	}
	if len(txOutputs) == 0 {
		return errors.New("output tx could not be found in tx_receipt_output")
	}

	// search record from account_pubkey_receipt by wallet_address
	// if is_allocated=`false`, update it `true`
	tm := time.Now()
	accountPublicKeyTable := make([]walletrepo.AccountPublicKeyTable, 1)
	accountPublicKeyTable[0].WalletAddress = txOutputs[0].OutputAddress
	accountPublicKeyTable[0].IsAllocated = true
	accountPublicKeyTable[0].UpdatedAt = &tm

	err = w.repo.UpdateIsAllocatedOnAccountPubKeyTable(account.AccountTypeReceipt, accountPublicKeyTable, nil, true)
	if err != nil {
		return errors.Wrap(err, "fail to call storager.UpdateIsAllocatedOnAccountPubKeyTable()")
	}

	return nil
}
