package service

import (
	"github.com/bookerzzz/grok"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
)

// UpdateStatus tx_paymentテーブル/tx_receiptテーブルのcurrent_tx_typeが3(送信済)のものを監視し、statusをupdateする
func (w *Wallet) UpdateStatus() error {
	//TODO:もしこのタイミングで、tx_typeが`done`で処理が止まっているものがあるケースはどうする？

	types := []enum.ActionType{enum.ActionTypeReceipt, enum.ActionTypePayment}
	for _, actionType := range types {
		hashes, err := w.DB.GetSentTxHashByTxTypeSent(actionType)
		if err != nil {
			//TODO:continueしたほうがいいか？
			return errors.Errorf("ActionType: %s, DB.GetSentTxHashByTxTypeSent() error: %v", actionType, err)
		}
		logger.Debug("hashes: ", hashes)

		// hashの詳細を取得する
		err = w.checkTransaction(hashes, actionType)
		if err != nil {
			//TODO:continueしたほうがいいか？
			return err
		}
	}

	return nil

	//tx_receipt/tx_payment
	//receiptHashs, err := w.DB.GetSentTxHashOnTxReceiptByTxTypeSent()
	//receiptHashs, err := w.DB.GetSentTxHashByTxTypeSent(enum.ActionTypeReceipt)
	//if err != nil {
	//	return errors.Errorf("DB.GetSentTxHashOnTxReceipt() error: %v", err)
	//}
	//logger.Debug("receiptHashes: ", receiptHashs)
	//
	//// hashの詳細を取得する
	//err = w.checkTransaction(receiptHashs, enum.ActionTypeReceipt)
	//if err != nil {
	//	return err
	//}

	//2.tx_payment(まとめる)
	//paymentHashs, err := w.DB.GetSentTxHashOnTxPaymentByTxTypeSent()
	//if err != nil {
	//	return errors.Errorf("DB.GetSentTxHashOnTxPayment() error: %v", err)
	//}
	//logger.Debug("paymentHashes: ", paymentHashs)
	//
	//// hashの詳細を取得する
	//return w.checkTransaction(paymentHashs, enum.ActionTypePayment)
}

// checkTransaction Bitcoin core APIでhashの状況をチェックし、もろもろ更新、通知を行う
func (w *Wallet) checkTransaction(hashs []string, actionType enum.ActionType) error {
	for _, hash := range hashs {
		//トランザクションの状態を取得
		tran, err := w.BTC.GetTransactionByTxID(hash)
		if err != nil {
			logger.Errorf("ActionType: %s, w.BTC.GetTransactionByTxID(): txID:%s, err:%v", actionType, hash, err)
			//TODO:実際に起きる場合はcanceledに更新したほうがいいか？
			continue
		}
		logger.Debugf("ActionType: %s, Transactions Confirmations", actionType)
		grok.Value(tran.Confirmations)

		//[debug] 終了したら消す
		//err = w.notifyUsers(hash, actionType)
		//if err != nil {
		//	return err
		//}
		//return nil

		//現在のconfirmationをチェック
		if tran.Confirmations >= int64(w.BTC.ConfirmationBlock()) {
			//指定にconfirmationに達したので、current_tx_typeをdoneに更新する
			_, err = w.DB.UpdateTxTypeDoneByTxHash(actionType, hash, nil, true)
			if err != nil {
				return errors.Errorf("ActionType: %s, DB.UpdateTxDoneByTxHash() error: %v", actionType, err)
			}
			//ユーザーに通知
			err = w.notifyUsers(hash, actionType)
			//TODO:errはどう処理すべき？他にもhashがあるので、continueしてもいいか？
			if err != nil {
				logger.Errorf("ActionType: %s, w.notifyUsers(%s, %s) error:%v",actionType, hash, actionType, err)
			}

			//if actionType == enum.ActionTypeReceipt {
			//	_, err = w.DB.UpdateTxReceipDoneByTxHash(hash, nil, true)
			//	if err != nil {
			//		return errors.Errorf("DB.UpdateTxReceipDoneByTxHash() error: %v", err)
			//	}
			//	//ユーザーに通知
			//	err = w.notifyUsers(hash, actionType)
			//	//TODO:errはどう処理すべき？他にもhashがあるので、continueしてもいいか？
			//	if err != nil {
			//		logger.Errorf("w.notifyUsers(%s, %s) error:%v", hash, actionType, err)
			//	}
			//} else if actionType == enum.ActionTypePayment {
			//	_, err = w.DB.UpdateTxPaymentDoneByTxHash(hash, nil, true)
			//	if err != nil {
			//		return errors.Errorf("DB.UpdateTxPaymentDoneByTxHash() error: %v", err)
			//	}
			//	//ユーザーに通知
			//	err = w.notifyUsers(hash, actionType)
			//	//TODO:errはどう処理すべき？他にもhashがあるので、continueしてもいいか？
			//	if err != nil {
			//		logger.Errorf("w.notifyUsers(%s, %s) error:%v", hash, actionType, err)
			//	}
			//}
		} else {
			//TODO:TestNet環境だと1000satoshiでもトランザクションが処理されてしまう
			//TODO:DBのsent_updated_atフィールドから一定時間立っても、指定したconfirmationに達しないものはキャンセルにして、
			//TODO:手数料を上げて再度トランザクションを作成する？？
			logger.Info("TODO:一定時間を過ぎてもトランザクションが終了しないものは通知したほうがいいかもしれない。")
		}
	}

	return nil
}

func (w *Wallet) notifyUsers(hash string, actionType enum.ActionType) error {

	logger.Debugf("notifyUsers() hash: %s", hash)

	//[tx_receiptの場合]
	if actionType == enum.ActionTypeReceipt {

		// 1.hashからidを取得(tx_receipt/tx_payment)
		receiptID, err := w.DB.GetTxIDBySentHash(actionType, hash)
		if err != nil {
			return errors.Errorf("ActionType: %s, DB.GetTxIDBySentHash() error: %v", actionType, err)
		}
		logger.Debug("notifyUsers() receiptID:", receiptID)

		// 2.tx_receipt_inputテーブルから該当のreceipt_idでレコードを取得
		txInputs, err := w.DB.GetTxReceiptInputByReceiptID(receiptID)
		if err != nil {
			return errors.Errorf("DB.GetTxReceiptInputByReceiptID(%d) error: %v", receiptID, err)
		}
		if len(txInputs) == 0 {
			logger.Debug("notifyUsers() len(txInputs) == 0")
			return nil
		}

		// 3.取得したinput_addressesに対して、入金が終了したことを通知する
		// TODO:NatsのPublisherとして通知すればいいか？
		for _, input := range txInputs {
			logger.Debug("input.InputAddress: ", input.InputAddress)
		}

		// 4.通知後はstatusをnotifiedに変更する
		_, err = w.DB.UpdateTxTypeNotifiedByID(actionType, receiptID, nil, true)
		if err != nil {
			return errors.Errorf("DB.UpdateTxReceipNotifiedByID() error: %v", err)
		}

	} else if actionType == enum.ActionTypePayment {
		//TODO:出金の通知フローは異なる。。。inputsはstoredの内部アドレスになっているため、payment_requestテーブルから情報を取得しないといけない
		// 1.hashからidを取得(tx_receipt/tx_payment)
		//paymentID, err := w.DB.GetTxPaymentIDBySentHash(hash)
		paymentID, err := w.DB.GetTxIDBySentHash(actionType, hash)

		if err != nil {
			return errors.Errorf("DB.GetTxPaymentIDBySentHash() error: %v", err)
		}
		logger.Debugf("notifyUsers() paymentID: %d", paymentID)

		// 2.payment_requestテーブルから該当のpayment_idでレコードを取得
		//txInputs, err := w.DB.GetTxReceiptInputByReceiptID(paymentID)
		paymentUsers, err := w.DB.GetPaymentRequestByPaymentID(paymentID)
		if err != nil {
			return errors.Errorf("DB.GetPaymentRequestByPaymentID(%d) error: %v", paymentID, err)
		}
		if len(paymentUsers) == 0 {
			logger.Debug("[Debug] notifyUsers() len(paymentUsers) == 0")
			return nil
		}

		// 3.取得したinput_addressesに対して、入金が終了したことを通知する
		// TODO:NatsのPublisherとして通知すればいいか？
		for _, user := range paymentUsers {
			logger.Debugf("user.AddressFrom: %s", user.AddressFrom)
		}

		tx := w.DB.RDB.MustBegin()
		// 4.通知後はstatusをnotifiedに変更する
		//_, err = w.DB.UpdateTxPaymentNotifiedByID(paymentID, tx, false)
		_, err = w.DB.UpdateTxTypeNotifiedByID(actionType, paymentID, tx, false)
		if err != nil {
			return errors.Errorf("DB.UpdateTxPaymentNotifiedByID() error: %v", err)
		}

		// 5.payment_requestテーブルのis_doneをtrueに更新する
		w.DB.UpdateIsDoneOnPaymentRequest(paymentID, tx, true)

	}

	return nil
}
