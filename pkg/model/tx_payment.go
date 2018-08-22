package model

import (
	"github.com/jmoiron/sqlx"
)

//enum.Actionに応じて、テーブルを切り替える

const (
	tableNamePayment = "tx_payment"
)

// TableNamePayment tx_paymentテーブル名を返す
func (m *DB) TableNamePayment() string {
	return tableNamePayment
}

// GetTxPaymentByID TxReceiptテーブルから該当するIDのレコードを返す
func (m *DB) GetTxPaymentByID(id int64) (*TxTable, error) {
	return m.getTxReceiptByID(m.TableNamePayment(), id)
}

// GetTxPaymentByUnsignedHex unsigned_hex_txをキーとしてレコードを取得する
func (m *DB) GetTxPaymentCountByUnsignedHex(hex string) (int64, error) {
	return m.getTxReceiptCountByUnsignedHex(m.TableNamePayment(), hex)
}

// GetSentTxHashOnTxPayment TxPaymentテーブルから送信済ステータスであるsent_hash_txの配列を返す
func (m *DB) GetSentTxHashOnTxPayment() ([]string, error) {
	return m.getSentTxHashOnTxReceipt(m.TableNamePayment())
}

// InsertTxPaymentForUnsigned TxReceiptテーブルに未署名トランザクションレコードを作成する
func (m *DB) InsertTxPaymentForUnsigned(txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.insertTxReceiptForUnsigned(m.TableNamePayment(), txReceipt, tx, isCommit)
}

// UpdateTxPaymentForSent TxReceiptテーブルのsigned_hex_tx, sent_hash_txを更新する
func (m *DB) UpdateTxPaymentForSent(txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.updateTxReceiptForSent(m.TableNamePayment(), txReceipt, tx, isCommit)
}

// UpdateTxPaymentForDone TxReceiptテーブルの該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (m *DB) UpdateTxPaymentForDone(hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.updateTxReceiptForDone(m.TableNamePayment(), hash, tx, isCommit)
}
