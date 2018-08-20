package model

import (
	"github.com/jmoiron/sqlx"
)

//enum.Actionに応じて、テーブルを切り替える

const (
	tableNamePaymentOutput = "tx_payment_output"
)

// TableNamePaymentOutput tx_payment_outputテーブル名を返す
func (m *DB) TableNamePaymentOutput() string {
	return tableNamePaymentOutput
}

// GetTxPaymentOutputByReceiptID TxReceiptInputテーブルから該当するIDのレコードを返す
func (m *DB) GetTxPaymentOutputByReceiptID(receiptID int64) (*TxOutput, error) {
	return m.getTxReceiptOutputByReceiptID(m.TableNamePaymentOutput(), receiptID)
}

// InsertTxPaymentOutputForUnsigned TxReceiptOutputテーブルに未署名トランザクションのinputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) InsertTxPaymentOutputForUnsigned(txReceiptOutputs []TxOutput, tx *sqlx.Tx, isCommit bool) error {
	return m.insertTxReceiptOutputForUnsigned(m.TableNamePaymentOutput(), txReceiptOutputs, tx, isCommit)
}
