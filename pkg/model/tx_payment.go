package model

import (
	"github.com/jmoiron/sqlx"
)

//enum.Actionに応じて、テーブルを切り替える

const (
	tableNamePayment = "tx_payment"
)

// TxPayment tx_paymentテーブル
type TxPayment struct {
	TxReceipt
}

// TableNamePayment tx_paymentテーブル名を返す
func (m *DB) TableNamePayment() string {
	return tableNamePayment
}

// GetTxPaymentByID TxReceiptテーブルから該当するIDのレコードを返す
func (m *DB) GetTxPaymentByID(id int64) (*TxReceipt, error) {
	return m.getTxReceiptByID(m.TableNamePayment(), id)
}

// GetTxPaymentByUnsignedHex unsigned_hex_txをキーとしてレコードを取得する
func (m *DB) GetTxPaymentByUnsignedHex(hex string) (int64, error) {
	return m.getTxReceiptByUnsignedHex(m.TableNamePayment(), hex)
}

// InsertTxPaymentForUnsigned TxReceiptテーブルに未署名トランザクションレコードを作成する
func (m *DB) InsertTxPaymentForUnsigned(txReceipt *TxReceipt, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.insertTxReceiptForUnsigned(m.TableNamePayment(), txReceipt, tx, isCommit)
}

// UpdateTxPaymentForSent TxReceiptテーブルのsigned_hex_tx, sent_hash_txを更新する
func (m *DB) UpdateTxPaymentForSent(txReceipt *TxReceipt, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.updateTxReceiptForSent(m.TableNamePayment(), txReceipt, tx, isCommit)
}
