package wallet

//Watch only wallet

import (
	"github.com/bookerzzz/grok"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
)

// UpdateStatus tx_paymentテーブル/tx_receiptテーブルのcurrent_tx_typeが3(送信済)のものを監視し、statusをupdateする
func (w *Wallet) UpdateStatus() error {
	//TODO:remove it
	if w.Type != WalletTypeWatchOnly {
		return errors.New("it's available on WalletTypeWatchOnly")
	}

	//tx_typeが`done`で処理が止まっているものがあるという前提で、処理を分ける

	types := []enum.ActionType{enum.ActionTypeReceipt, enum.ActionTypePayment, enum.ActionTypeTransfer}

	//1.ここでは送信済のみが対象
	for _, actionType := range types {
		err := w.updateStatusForTxTypeSent(actionType)
		if err != nil {
			return errors.Errorf("ActionType: %s, updateStatusForTxTypeSent() error: %s", actionType, err)
		}
	}

	//2.tx_typeがdoneのトランザクションに対して、通知
	for _, actionType := range types {
		err := w.updateStatusForTxTypeDone(actionType)
		if err != nil {
			return errors.Errorf("ActionType: %s, updateStatusForTxTypeDone() error: %s", actionType, err)
		}
	}

	return nil
}

// current_tx_type更新処理
func (w *Wallet) updateStatusForTxTypeSent(actionType enum.ActionType) error {
	// 送信済statusのものを取得
	hashes, err := w.DB.GetSentTxHashByTxTypeSent(actionType)
	if err != nil {
		return errors.Errorf("ActionType: %s, DB.GetSentTxHashByTxTypeSent() error: %s", actionType, err)
	}
	logger.Debugf("ActionType: %s, hashes: %v", actionType, hashes)

	// hashの詳細を取得し、confirmationをチェックし、指定数以上であれば、ux_typeを更新する
	for _, hash := range hashes {
		err = w.checkTransaction(hash, actionType)
		if err != nil {
			logger.Errorf("ActionType: %s, w.checkTransaction(%s, %s) error: %s", actionType, hash, actionType, err)
			continue
		}
	}
	return nil
}

func (w *Wallet) updateStatusForTxTypeDone(actionType enum.ActionType) error {
	// 更新
	hashes, err := w.DB.GetSentTxHashByTxTypeDone(actionType)
	if err != nil {
		return errors.Errorf("ActionType: %s, DB.GetSentTxHashByTxTypeDone() error: %s", actionType, err)
	}
	logger.Debugf("ActionType: %s, hashes: %v", actionType, hashes)

	// 通知する
	for _, hash := range hashes {
		//ユーザーに通知
		id, err := w.notifyUsers(hash, actionType)
		if err != nil {
			logger.Errorf("ActionType: %s, w.notifyUsers(%s, %s) error: %s", actionType, hash, actionType, err)
			continue
		}
		if id == 0 {
			continue
		}

		//通知が成功したので、tx_typeを更新する(別funcで対応)
		err = w.updateTxTypeNotified(id, hash, actionType)
		//仮にここがエラーになっても、通知は成功している。。。が、また処理が走ってしまう。。。
		if err != nil {
			logger.Errorf("ActionType: %s, w.updateTxTypeNotified(%s, %s) error: %s", actionType, hash, actionType, err)
			continue
		}
	}
	return nil
}

// checkTransaction Bitcoin core APIでhashの状況をチェックし、もろもろ更新、通知を行う
func (w *Wallet) checkTransaction(hash string, actionType enum.ActionType) error {
	//トランザクションの状態を取得
	tran, err := w.BTC.GetTransactionByTxID(hash)
	if err != nil {
		//logger.Errorf("ActionType: %s, w.BTC.GetTransactionByTxID(): txID:%s, err:%s", actionType, hash, err)
		//TODO:実際に起きる場合はcanceledに更新したほうがいいか？
		return errors.Errorf("ActionType: %s, w.BTC.GetTransactionByTxID(): txID:%s, err: %s", actionType, hash, err)
	}
	logger.Debugf("ActionType: %s, Transactions Confirmations", actionType)
	grok.Value(tran.Confirmations)

	//現在のconfirmationをチェック
	if tran.Confirmations >= int64(w.BTC.ConfirmationBlock()) {
		//指定にconfirmationに達したので、current_tx_typeをdoneに更新する
		_, err = w.DB.UpdateTxTypeDoneByTxHash(actionType, hash, nil, true)
		if err != nil {
			return errors.Errorf("ActionType: %s, DB.UpdateTxTypeDoneByTxHash() error: %s", actionType, err)
		}
	} else {
		//TODO:TestNet環境だと1000satoshiでもトランザクションが処理されてしまう
		//TODO:DBのsent_updated_atフィールドから一定時間立っても、指定したconfirmationに達しないものはキャンセルにして、
		//TODO:手数料を上げて再度トランザクションを作成する？？
		logger.Info("TODO:一定時間を過ぎてもトランザクションが終了しないものは通知したほうがいいかもしれない。")
	}

	return nil
}

// notifyUsers 入金/出金が終了したことを通知する
func (w *Wallet) notifyUsers(hash string, actionType enum.ActionType) (int64, error) {
	logger.Debugf("ActionType: %s, notifyUsers() hash: %s", actionType, hash)

	//id: receiptID/paymentID
	var (
		id  int64
		err error
	)

	//[tx_receiptの場合]
	if actionType == enum.ActionTypeReceipt {

		// 1.hashからidを取得(tx_receipt/tx_payment)
		id, err = w.DB.GetTxIDBySentHash(actionType, hash)
		if err != nil {
			return 0, errors.Errorf("ActionType: %s, DB.GetTxIDBySentHash() error: %s", actionType, err)
		}
		logger.Debug("notifyUsers() receiptID:", id)

		// 2.tx_receipt_inputテーブルから該当のreceipt_idでレコードを取得
		txInputs, err := w.DB.GetTxInputByReceiptID(enum.ActionTypeReceipt, id)
		if err != nil {
			return 0, errors.Errorf("ActionType: %s, DB.GetTxInputByReceiptID(%d) error: %s", actionType, id, err)
		}
		if len(txInputs) == 0 {
			logger.Debug("notifyUsers() len(txInputs) == 0")
			return 0, nil
		}

		// 3.取得したinput_addressesに対して、入金が終了したことを通知する
		// TODO:NatsのPublisherとして通知すればいいか？
		for _, input := range txInputs {
			logger.Debug("input.InputAddress: ", input.InputAddress)
		}

	} else if actionType == enum.ActionTypePayment {
		//出金の通知フローは異なる。inputsはstoredの内部アドレスになっているため、payment_requestテーブルから情報を取得しないといけない
		// 1.hashからidを取得(tx_receipt/tx_payment)
		id, err = w.DB.GetTxIDBySentHash(actionType, hash)

		if err != nil {
			return 0, errors.Errorf("ActionType: %s, DB.GetTxIDBySentHash() error: %s", actionType, err)
		}
		logger.Debugf("notifyUsers() paymentID: %d", id)

		// 2.payment_requestテーブルから該当のpayment_idでレコードを取得
		paymentUsers, err := w.DB.GetPaymentRequestByPaymentID(id)
		if err != nil {
			return 0, errors.Errorf("ActionType: %s, DB.GetPaymentRequestByPaymentID(%d) error: %s", actionType, id, err)
		}
		if len(paymentUsers) == 0 {
			logger.Debug("[Debug] notifyUsers() len(paymentUsers) == 0")
			return 0, nil
		}

		// 3.取得したinput_addressesに対して、入金が終了したことを通知する
		// TODO:NatsのPublisherとして通知すればいいか？
		for _, user := range paymentUsers {
			logger.Debugf("user.AddressFrom: %s", user.AddressFrom)
		}

	}

	return id, nil
}

//updateTxTypeNotified tx_typeを通知済に更新する
func (w *Wallet) updateTxTypeNotified(id int64, hash string, actionType enum.ActionType) error {
	//id: receiptID/paymentID

	if actionType == enum.ActionTypeReceipt {
		// 通知後はstatusをnotifiedに変更する
		_, err := w.DB.UpdateTxTypeNotifiedByID(actionType, id, nil, true)
		if err != nil {
			return errors.Errorf("ActionType: %s, DB.UpdateTxTypeNotifiedByID() error: %s", actionType, err)
		}

	} else if actionType == enum.ActionTypePayment {
		tx := w.DB.RDB.MustBegin()
		// 通知後はstatusをnotifiedに変更する
		_, err := w.DB.UpdateTxTypeNotifiedByID(actionType, id, tx, false)
		if err != nil {
			return errors.Errorf("ActionType: %s, DB.UpdateTxTypeNotifiedByID() error: %s", actionType, err)
		}

		// payment_requestテーブルのis_doneをtrueに更新する
		_, err = w.DB.UpdateIsDoneOnPaymentRequest(id, tx, true)
		if err != nil {
			return errors.Errorf("ActionType: %s, DB.UpdateIsDoneOnPaymentRequest() error: %s", actionType, err)
		}
	}

	return nil
}
