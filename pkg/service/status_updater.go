package service

import (
	"log"

	"github.com/pkg/errors"
	"github.com/bookerzzz/grok"
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
	w.checkTransaction(receiptHashs)


	//2.tx_payment
	paymentHashs, err := w.DB.GetSentTxHashOnTxPayment()
	if err != nil {
		errors.Errorf("DB.GetSentTxHashOnTxPayment() error: %v", err)
	}
	log.Println(paymentHashs)

	// hashの詳細を取得する
	w.checkTransaction(receiptHashs)

	return nil
}

func (w *Wallet) checkTransaction(hashs []string) error {
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
			log.Println("OK!")
		}
	}
	return nil
}