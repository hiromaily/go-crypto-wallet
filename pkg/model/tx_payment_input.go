package model

import (
	"github.com/jmoiron/sqlx"
)

//enum.Actionに応じて、テーブルを切り替える

const (
	tableNamePaymentInput = "tx_payment_input"
)

// TxPaymentInput tx_receipt_inputテーブル(tx_payment_inputとしても利用)
type TxPaymentInput struct {
	TxReceiptInput
}

// TableNamePaymentDetail tx_payment_detailテーブル名を返す
func (m *DB) TableNamePaymentInput() string {
	return tableNamePaymentInput
}

// GetTxPaymentInputByReceiptID TxReceiptInputテーブルから該当するIDのレコードを返す
func (m *DB) GetTxPaymentInputByReceiptID(receiptID int64) (*TxReceiptInput, error) {
	return m.getTxReceiptInputByReceiptID(m.TableNamePaymentInput(), receiptID)
}

// InsertTxPaymentInputForUnsigned TxReceiptInputテーブルに未署名トランザクションのinputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) InsertTxPaymentInputForUnsigned(txReceiptDetails []TxReceiptInput, tx *sqlx.Tx, isCommit bool) error {
	return m.insertTxReceiptInputForUnsigned(m.TableNamePaymentInput(), txReceiptDetails, tx, isCommit)
}
