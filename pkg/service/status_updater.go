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
	receiptHashs, err := w.DB.GetSentTxHashOnTxReceipt()
	if err != nil {
		errors.Errorf("DB.GetSentTxHashOnTxReceipt() error: %v", err)
	}
	log.Println(receiptHashs)

	// hashの詳細を取得する
	w.checkTransaction(receiptHashs, enum.ActionTypeReceipt)

	//2.tx_payment
	paymentHashs, err := w.DB.GetSentTxHashOnTxPayment()
	if err != nil {
		errors.Errorf("DB.GetSentTxHashOnTxPayment() error: %v", err)
	}
	log.Println(paymentHashs)

	// hashの詳細を取得する
	w.checkTransaction(paymentHashs, enum.ActionTypePayment)

	return nil
}

func (w *Wallet) checkTransaction(hashs []string, actionType enum.ActionType) error {
	for _, hash := range hashs {
		tran, err := w.BTC.GetTransactionByTxID(hash)
		if err != nil {
			log.Printf("w.BTC.GetTransactionByTxID(): txID:%s, err:%v", hash, err)
			//TODO:実際に起きる場合はcanceledに更新したほうがいいか？
			continue
		}
		log.Println("[Debug]Transactions")
		grok.Value(tran.Confirmations)

		if tran.Confirmations >= int64(w.BTC.ConfirmationBlock()) {
			//指定にconfirmationに達したので、doneに更新する
			if actionType == enum.ActionTypeReceipt {
				_, err = w.DB.UpdateTxReceiptForDone(hash, nil, true)
				if err != nil {
					return errors.Errorf("DB.UpdateTxReceiptForDone() error: %v", err)
				}
			} else if actionType == enum.ActionTypePayment {
				_, err = w.DB.UpdateTxPaymentForDone(hash, nil, true)
				if err != nil {
					return errors.Errorf("DB.UpdateTxPaymentForDone() error: %v", err)
				}
			}
		}
	}
	return nil
}
