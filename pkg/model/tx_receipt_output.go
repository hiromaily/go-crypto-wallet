package model

import "time"

//CREATE TABLE `tx_receipt_output` (
//`id`             BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT'ID',
//`receipt_id`     BIGINT(20) UNSIGNED NOT NULL COMMENT'tx_receipt ID',
//`output_address` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'outputに利用されるaddress(受け取る人)',
//`output_account` VARCHAR(255) COLLATE utf8_unicode_ci NOT NULL COMMENT'outputに利用されるaccount(受け取る人)',
//`output_amount`  DECIMAL(26,10) NOT NULL COMMENT'outputに利用されるamount(入金金額)',
//`isChange`       BOOL DEFAULT false COMMENT'お釣り用のoutputであればtrue',
//`updated_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT'更新日時',

const (
	tableNameReceiptOutput = "tx_receipt_output"
)

// TxReceiptOutput tx_receipt_outputテーブル
type TxReceiptOutput struct {
	ID            int64      `db:"id"`
	ReceiptID     int64      `db:"receipt_id"`
	OutputAddress string     `db:"output_address"`
	OutputAccount string     `db:"output_account"`
	OutputAmount  string     `db:"output_amount"`
	IsChange      bool       `db:"is_change"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

// TableNameReceiptOutput tx_receipt_outputテーブル名を返す
func (m *DB) TableNameReceiptOutput() string {
	return tableNameReceiptOutput
}
