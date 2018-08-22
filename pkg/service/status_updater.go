package service

import (
	"log"

	"github.com/bookerzzz/grok"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/pkg/errors"
)

// UpdateStatus tx_paymentテーブル/tx_receiptテーブルのcurrent_tx_typeが3(送信済)のものを監視し、statusをupdateする
func (w *Wallet) UpdateStatus() error {

	//1.tx_receipt
	receiptHashs, err := w.DB.GetSentTxHashOnTxReceiptByTxTypeSent()
	if err != nil {
		errors.Errorf("DB.GetSentTxHashOnTxReceipt() error: %v", err)
	}
	log.Println("[Debug] receiptHashs", receiptHashs)

	// hashの詳細を取得する
	err = w.checkTransaction(receiptHashs, enum.ActionTypeReceipt)
	if err != nil {
		return err
	}

	//2.tx_payment
	paymentHashs, err := w.DB.GetSentTxHashOnTxPaymentByTxTypeSent()
	if err != nil {
		errors.Errorf("DB.GetSentTxHashOnTxPayment() error: %v", err)
	}
	log.Println("[Debug] paymentHashs", paymentHashs)

	// hashの詳細を取得する
	return w.checkTransaction(paymentHashs, enum.ActionTypePayment)
}

// checkTransaction Bitcoin core APIでhashの状況をチェックし、もろもろ更新、通知を行う
func (w *Wallet) checkTransaction(hashs []string, actionType enum.ActionType) error {
	for _, hash := range hashs {
		//トランザクションの状態を取得
		tran, err := w.BTC.GetTransactionByTxID(hash)
		if err != nil {
			log.Printf("w.BTC.GetTransactionByTxID(): txID:%s, err:%v", hash, err)
			//TODO:実際に起きる場合はcanceledに更新したほうがいいか？
			continue
		}
		log.Println("[Debug]Transactions Confirmations")
		grok.Value(tran.Confirmations)

		//debug 終了したら消す
		//err = w.notifyUsers(hash, actionType)
		//if err != nil {
		//	return err
		//}

		//現在のconfirmationをチェック
		if tran.Confirmations >= int64(w.BTC.ConfirmationBlock()) {
			//指定にconfirmationに達したので、current_tx_typeをdoneに更新する
			if actionType == enum.ActionTypeReceipt {
				_, err = w.DB.UpdateTxReceipDoneByTxHash(hash, nil, true)
				if err != nil {
					return errors.Errorf("DB.UpdateTxReceipDoneByTxHash() error: %v", err)
				}
				//ユーザーに通知
				w.notifyUsers(hash, actionType)
			} else if actionType == enum.ActionTypePayment {
				_, err = w.DB.UpdateTxPaymentDoneByTxHash(hash, nil, true)
				if err != nil {
					return errors.Errorf("DB.UpdateTxPaymentDoneByTxHash() error: %v", err)
				}
				//ユーザーに通知
				w.notifyUsers(hash, actionType)
			}
		} else {
			//TODO:TestNet環境だと1000satoshiでもトランザクションが処理されてしまう
			//TODO:DBのsent_updated_atフィールドから一定時間立っても、指定したconfirmationに達しないものはキャンセルにして、
			//TODO:手数料を上げて再度トランザクションを作成する？？
			log.Println("TODO:一定時間を過ぎてもトランザクションが終了しないものは通知したほうがいいかもしれない。")
		}
	}
	return nil
}

func (w *Wallet) notifyUsers(hash string, actionType enum.ActionType) error {
	//[tx_receiptの場合]
	if actionType == enum.ActionTypeReceipt {
		log.Println("[Debug] notifyUsers() hash:", hash)
		// 1.hashからidを取得(tx_receipt/tx_payment)
		receiptID, err := w.DB.GetTxReceiptIDBySentHash(hash)
		if err != nil {
			return errors.Errorf("DB.GetTxReceiptIDBySentHash() error: %v", err)
		}
		log.Println("[Debug] notifyUsers() receiptID:", receiptID)

		// 2.tx_receipt_inputテーブルから該当のreceipt_idでレコードを取得
		txInputs, err := w.DB.GetTxReceiptInputByReceiptID(receiptID)
		if err != nil {
			return errors.Errorf("DB.GetTxReceiptInputByReceiptID(%d) error: %v", receiptID, err)
		}
		if len(txInputs) == 0 {
			log.Println("[Debug] notifyUsers() len(txInputs) == 0")
			return nil
		}

		// 3.取得したinput_addressesに対して、入金が終了したことを通知する
		// TODO:NatsのPublisherとして通知すればいいか？
		for _, input := range txInputs {
			log.Println("[Debug]input.InputAddress: ", input.InputAddress)
		}

		// 4.通知後はstatusをnotifiedに変更する
		_, err = w.DB.UpdateTxReceipNotifiedByTxHash(hash, nil, true)
		if err != nil {
			return errors.Errorf("DB.UpdateTxReceipNotifiedByTxHash() error: %v", err)
		}

	}
	return nil
}
