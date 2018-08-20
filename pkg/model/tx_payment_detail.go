package model

import (
	"github.com/jmoiron/sqlx"
)

//enum.Actionに応じて、テーブルを切り替える

const (
	tableNamePaymentDetail = "tx_payment_detail"
)

// TxPaymentDetail tx_receipt_detailテーブル(tx_payment_detailとしても利用)
type TxPaymentDetail struct {
	TxReceiptDetail
}

// TableNamePaymentDetail tx_payment_detailテーブル名を返す
func (m *DB) TableNamePaymentDetail() string {
	return tableNamePaymentDetail
}

// GetTxPaymentDetailByReceiptID TxReceiptDetailテーブルから該当するIDのレコードを返す
func (m *DB) GetTxPaymentDetailByReceiptID(receiptID int64) (*TxReceiptDetail, error) {
	return m.getTxReceiptDetailByReceiptID(m.TableNamePaymentDetail(), receiptID)
}

// InsertTxPaymentDetailForUnsigned TxReceiptDetailテーブルに未署名トランザクションのinputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) InsertTxPaymentDetailForUnsigned(txReceiptDetails []TxReceiptDetail, tx *sqlx.Tx, isCommit bool) error {
	return m.insertTxReceiptDetailForUnsigned(m.TableNamePaymentDetail(), txReceiptDetails, tx, isCommit)
}
