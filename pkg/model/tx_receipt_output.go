package model

//
//import (
//	"fmt"
//	"github.com/jmoiron/sqlx"
//	"time"
//)
//
////CREATE TABLE `tx_receipt_output` (
////`id`             BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
////`receipt_id`     BIGINT(20) UNSIGNED NOT NULL COMMENT'tx_receipt ID',
////`output_address` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'outputに利用されるaddress(受け取る人)',
////`output_account` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'outputに利用されるaccount(受け取る人)',
////`output_amount`  DECIMAL(26,10) NOT NULL COMMENT'outputに利用されるamount(入金金額)',
////`isChange`       BOOL DEFAULT false COMMENT'お釣り用のoutputであればtrue',
////`updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',
//
//const (
//	tableNameReceiptOutput = "tx_receipt_output"
//)
//
//// TxOutput tx_receipt_output/tx_payment_outputテーブル
//type TxOutput struct {
//	ID            int64      `db:"id"`
//	ReceiptID     int64      `db:"receipt_id"`
//	OutputAddress string     `db:"output_address"`
//	OutputAccount string     `db:"output_account"`
//	OutputAmount  string     `db:"output_amount"`
//	IsChange      bool       `db:"is_change"`
//	UpdatedAt     *time.Time `db:"updated_at"`
//}
//
//// TableNameReceiptOutput tx_receipt_outputテーブル名を返す
//func (m *DB) TableNameReceiptOutput() string {
//	return tableNameReceiptOutput
//}
//
//// getTxReceiptOutputByReceiptID TxReceiptOutputテーブルから該当するIDのレコードを返す
//func (m *DB) getTxReceiptOutputByReceiptID(tbl string, receiptID int64) ([]TxOutput, error) {
//	sql := "SELECT * FROM %s WHERE receipt_id=?"
//	sql = fmt.Sprintf(sql, tbl)
//
//	txReceiptOutputs := []TxOutput{}
//	err := m.RDB.Select(&txReceiptOutputs, sql, receiptID)
//
//	return txReceiptOutputs, err
//}
//
//// GetTxReceiptOutputByReceiptID TxReceiptOutputテーブルから該当するIDのレコードを返す
//func (m *DB) GetTxReceiptOutputByReceiptID(receiptID int64) ([]TxOutput, error) {
//	return m.getTxReceiptOutputByReceiptID(m.TableNameReceiptOutput(), receiptID)
//}
//
//// insertTxReceiptOutputForUnsigned TxReceiptOutputテーブルに未署名トランザクションのoutputに使われたtxレコードを作成する
////TODO:BulkInsertがやりたい
//func (m *DB) insertTxReceiptOutputForUnsigned(tbl string, txReceiptOutputs []TxOutput, tx *sqlx.Tx, isCommit bool) error {
//
//	sql := `
//INSERT INTO %s (receipt_id, output_address, output_account, output_amount, is_change)
//VALUES (:receipt_id,  :output_address, :output_account, :output_amount, :is_change)
//`
//	sql = fmt.Sprintf(sql, tbl)
//
//	if tx == nil {
//		tx = m.RDB.MustBegin()
//	}
//
//	for _, txReceiptOutput := range txReceiptOutputs {
//		_, err := tx.NamedExec(sql, txReceiptOutput)
//		if err != nil {
//			tx.Rollback()
//			return err
//		}
//	}
//
//	if isCommit {
//		tx.Commit()
//	}
//
//	return nil
//}
//
//// InsertTxReceiptOutputForUnsigned TxReceiptOutputテーブルに未署名トランザクションのoutputに使われたtxレコードを作成する
////TODO:BulkInsertがやりたい
//func (m *DB) InsertTxReceiptOutputForUnsigned(txReceiptOutputs []TxOutput, tx *sqlx.Tx, isCommit bool) error {
//	return m.insertTxReceiptOutputForUnsigned(m.TableNameReceiptOutput(), txReceiptOutputs, tx, isCommit)
//}
