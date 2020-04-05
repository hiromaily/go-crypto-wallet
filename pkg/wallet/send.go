package wallet

//Watch only wallet

import (
	"time"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/txfile"
)

// coldwallet側で署名済みトランザクションを作成したものから、送金処理を行う

// SendFromFile 渡されたファイルから署名済transactionを読み取り、送信を行う
func (w *Wallet) SendFromFile(filePath string) (string, error) {
	if w.Type != enum.WalletTypeWatchOnly {
		return "", errors.New("it's available on WatchOnlyWallet")
	}

	//ファイル名から、tx_receipt_idを取得する
	//payment_5_unsigned_1534466246366489473
	txReceiptID, actionType, _, err := txfile.ParseFile(filePath, []enum.TxType{enum.TxTypeSigned})
	if err != nil {
		return "", errors.Errorf("txfile.ParseFile() error: %s", err)
	}

	//ファイルからhexを読み取る
	signedHex, err := txfile.ReadFile(filePath)
	if err != nil {
		return "", errors.Errorf("txfile.ReadFile() error: %s", err)
	}

	//送信
	hash, err := w.BTC.SendTransactionByHex(signedHex)
	if err != nil {
		//TODO:本番環境ではBitcoinネットワークに取り込まれなくても、ここでエラーがでる？？その場合、手数料をあげて再トランザクションを生成するように促す必要がある
		//-26: 16: mandatory-script-verify-flag-failed (Operation not valid with the current stack size)
		//=> 署名が不十分だとこれが出るらしい
		logger.Error("This error implies new transsaction should be created from the beginning")
		return "", errors.Errorf("BTC.SendTransactionByHex() error: %s", err)
	}

	//DB更新 tx_receipt/tx_payment
	err = w.updateHexForSentTx(txReceiptID, signedHex, hash.String(), actionType)
	if err != nil {
		//TODO:仮にここでエラーが出たとしても、送信したという事実に変わりはない。ここのみを再度実行する仕組みが必要
		return "", errors.Errorf("w.updateHexForSentTx() error: %s", err)
	}

	//DB更新 account_pubkey_receiptのみ
	err = w.updateIsAllocatedForAccountPubkey(txReceiptID, actionType)
	if err != nil {
		//TODO:仮にここでエラーが出たとしても、送信したという事実に変わりはない。ここのみを再度実行する仕組みが必要
		return "", errors.Errorf("w.updateIsAllocatedForAccountPubkey() error: %s", err)
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

	affectedNum, err = w.DB.UpdateTxAfterSent(actionType, &txReceipt, nil, true)

	if err != nil {
		return errors.Errorf("DB.UpdateTxAfterSent(): error: %s", err)
	}
	if affectedNum == 0 {
		return errors.New("DB.UpdateTxAfterSent(): tx_receipt table was not updated")
	}

	return nil
}

func (w *Wallet) updateIsAllocatedForAccountPubkey(txReceiptID int64, actionType enum.ActionType) error {
	//tx_receiptの場合のみ
	if actionType == enum.ActionTypeReceipt {
		return nil
	}

	//1.tx_receipt_outputのreceipt_idに一致する1レコードのoutput_addressを取得
	txOutputs, err := w.DB.GetTxOutputByReceiptID(actionType, txReceiptID)
	if err != nil {
		return errors.Errorf("DB.GetTxOutputByReceiptID(): error: %s", err)
	}
	if len(txOutputs) == 0 {
		return errors.New("output tx could not be found in tx_receipt_output")
	}

	//2.account_pubkey_receiptのwallet_addressで検索し、is_allocatedがfalseであれば、trueに更新する
	//tx_paymentの場合、勝手に分散されていて、使用済かどうかは、Quoineから補充するタイミングで、更新する必要がある
	tm := time.Now()
	accountPublicKeyTable := make([]model.AccountPublicKeyTable, 1)
	accountPublicKeyTable[0].WalletAddress = txOutputs[0].OutputAddress
	accountPublicKeyTable[0].IsAllocated = true
	accountPublicKeyTable[0].UpdatedAt = &tm

	w.DB.UpdateIsAllocatedOnAccountPubKeyTable(account.AccountTypeReceipt, accountPublicKeyTable, nil, true)

	return nil
}
