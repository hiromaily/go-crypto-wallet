package model

import (
	"fmt"
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/jmoiron/sqlx"
)

//CREATE TABLE `tx_receipt_output` (
//`id`             BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
//`receipt_id`     BIGINT(20) UNSIGNED NOT NULL COMMENT'tx_receipt ID',
//`output_address` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'outputに利用されるaddress(受け取る人)',
//`output_account` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'outputに利用されるaccount(受け取る人)',
//`output_amount`  DECIMAL(26,10) NOT NULL COMMENT'outputに利用されるamount(入金金額)',
//`isChange`       BOOL DEFAULT false COMMENT'お釣り用のoutputであればtrue',
//`updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',

var txTableOutputName = map[enum.ActionType]string{
	"receipt": "tx_receipt_output",
	"payment": "tx_payment_output",
}

// TxOutput tx_receipt_output/tx_payment_outputテーブル
type TxOutput struct {
	ID            int64      `db:"id"`
	ReceiptID     int64      `db:"receipt_id"`
	OutputAddress string     `db:"output_address"`
	OutputAccount string     `db:"output_account"`
	OutputAmount  string     `db:"output_amount"`
	IsChange      bool       `db:"is_change"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

// getTxOutputByReceiptID 該当するIDのレコードを返す
func (m *DB) getTxOutputByReceiptID(tbl string, receiptID int64) ([]TxOutput, error) {
	sql := "SELECT * FROM %s WHERE receipt_id=?"
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var txReceiptOutputs []TxOutput
	err := m.RDB.Select(&txReceiptOutputs, sql, receiptID)

	return txReceiptOutputs, err
}

// GetTxOutputByReceiptID 該当するIDのレコードを返す
func (m *DB) GetTxOutputByReceiptID(actionType enum.ActionType, receiptID int64) ([]TxOutput, error) {
	return m.getTxOutputByReceiptID(txTableOutputName[actionType], receiptID)
}

// insertTxOutputForUnsigned 未署名トランザクションのoutputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) insertTxOutputForUnsigned(tbl string, txReceiptOutputs []TxOutput, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (receipt_id, output_address, output_account, output_amount, is_change) 
VALUES (:receipt_id,  :output_address, :output_account, :output_amount, :is_change)
`
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, txReceiptOutput := range txReceiptOutputs {
		_, err := tx.NamedExec(sql, txReceiptOutput)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if isCommit {
		tx.Commit()
	}

	return nil
}

// InsertTxOutputForUnsigned 未署名トランザクションのoutputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) InsertTxOutputForUnsigned(actionType enum.ActionType, txReceiptOutputs []TxOutput, tx *sqlx.Tx, isCommit bool) error {
	return m.insertTxOutputForUnsigned(txTableOutputName[actionType], txReceiptOutputs, tx, isCommit)
}
