package model

//
//import (
//	"github.com/jmoiron/sqlx"
//)
//
////enum.Actionに応じて、テーブルを切り替える
//
//const (
//	tableNamePaymentInput = "tx_payment_input"
//)
//
//// TableNamePaymentInput tx_payment_inputテーブル名を返す
//func (m *DB) TableNamePaymentInput() string {
//	return tableNamePaymentInput
//}
//
//// GetTxPaymentInputByReceiptID TxReceiptInputテーブルから該当するIDのレコードを返す
//func (m *DB) GetTxPaymentInputByReceiptID(receiptID int64) ([]TxInput, error) {
//	return m.getTxReceiptInputByReceiptID(m.TableNamePaymentInput(), receiptID)
//}
//
//// InsertTxPaymentInputForUnsigned TxReceiptInputテーブルに未署名トランザクションのinputに使われたtxレコードを作成する
////TODO:BulkInsertがやりたい
//func (m *DB) InsertTxPaymentInputForUnsigned(txReceiptDetails []TxInput, tx *sqlx.Tx, isCommit bool) error {
//	return m.insertTxReceiptInputForUnsigned(m.TableNamePaymentInput(), txReceiptDetails, tx, isCommit)
//}
