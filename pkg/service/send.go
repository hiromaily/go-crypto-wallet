package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/file"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/pkg/errors"
	"log"
	"time"
)

// coldwallet側で署名済みトランザクションを作成したものから、送金処理を行う

// SendFromFile 渡されたファイルから署名済transactionを読み取り、送信を行う
func (w *Wallet) SendFromFile(filePath string) (string, error) {
	//ファイル名から、tx_receipt_idを取得する
	//payment_5_unsigned_1534466246366489473
	txReceiptID, actionType, _, err := file.ParseFile(filePath, "signed")
	if err != nil {
		return "", err
	}

	//ファイルからhexを読み取る
	signedHex, err := file.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	log.Println(signedHex)

	//送信
	hash, err := w.BTC.SendTransactionByHex(signedHex)
	if err != nil {
		//TODO: これが失敗するのはどういうときか？
		//-26: 16: mandatory-script-verify-flag-failed (Operation not valid with the current stack size)
		//=> 署名が不十分だとこれが出るらしい
		return "", err
	}

	//DB更新
	err = w.updateHexForSentTx(txReceiptID, signedHex, hash.String(), actionType)
	if err != nil {
		//TODO:仮にここでエラーが出たとしても、送信したという事実に変わりはない
		return "", err
	}

	return hash.String(), nil
}

//
func (w *Wallet) updateHexForSentTx(txReceiptID int64, signedHex, sentTxID string, actionType enum.ActionType) error {
	//1.TxReceiptテーブル
	t := time.Now()
	txReceipt := model.TxTable{}
	txReceipt.ID = txReceiptID
	txReceipt.SignedHexTx = signedHex
	txReceipt.SentHashTx = sentTxID
	txReceipt.SentUpdatedAt = &t
	txReceipt.TxType = enum.TxTypeValue[enum.TxTypeSent] //3:未署名

	var (
		affectedNum int64
		err         error
	)

	//ActionTypeによって、処理を分ける
	if actionType == enum.ActionTypeReceipt {
		affectedNum, err = w.DB.UpdateTxReceiptForSent(
			&txReceipt, nil, true)
	} else if actionType == enum.ActionTypePayment {
		affectedNum, err = w.DB.UpdateTxPaymentForSent(
			&txReceipt, nil, true)
	}

	if err != nil {
		return errors.Errorf("DB.UpdateTxReceiptForSent(): error: %v", err)
	}
	if affectedNum == 0 {
		return errors.Errorf("DB.UpdateTxReceiptForSent(): tx_receipt table was not updated")
	}

	return nil
}
