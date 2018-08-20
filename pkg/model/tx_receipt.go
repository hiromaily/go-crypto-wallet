package model

import (
	"time"

	"fmt"
	"github.com/jmoiron/sqlx"
)

//enum.Actionに応じて、テーブルを切り替える

const (
	tableNameReceipt = "tx_receipt"
)

// TxReceipt tx_receiptテーブル
type TxReceipt struct {
	ID            int64  `db:"id"`
	UnsignedHexTx string `db:"unsigned_hex_tx"`
	SignedHexTx   string `db:"signed_hex_tx"`
	SentHexTx     string `db:"sent_hash_tx"`
	//TotalAmount   float64 `db:"total_amount"` //Float型はInsert後に誤差が生じる可能性がある
	//Fee           float64 `db:"fee"`
	TotalAmount       string     `db:"total_amount"`
	Fee               string     `db:"fee"`
	ReceiverAddress   string     `db:"receiver_address"`
	TxType            uint8      `db:"current_tx_type"` //TODO: intからuint8にできないか？
	UnsignedUpdatedAt *time.Time `db:"unsigned_updated_at"`
	SignedUpdatedAt   *time.Time `db:"signed_updated_at"`
	SentUpdatedAt     *time.Time `db:"sent_updated_at"`
}

// TableNameReceipt tx_receiptテーブル名を返す
func (m *DB) TableNameReceipt() string {
	return tableNameReceipt
}

// getTxReceiptByID TxReceiptテーブルから該当するIDのレコードを返す
func (m *DB) getTxReceiptByID(tbl string, id int64) (*TxReceipt, error) {
	sql := "SELECT * FROM %s WHERE id=?"
	sql = fmt.Sprintf(sql, tbl)

	txReceipt := TxReceipt{}
	err := m.RDB.Get(&txReceipt, sql, id)

	return &txReceipt, err
}

// GetTxReceiptByID TxReceiptテーブルから該当するIDのレコードを返す
func (m *DB) GetTxReceiptByID(id int64) (*TxReceipt, error) {
	return m.getTxReceiptByID(m.TableNameReceipt(), id)
}

// getTxReceiptByUnsignedHex unsigned_hex_txをキーとしてレコードを取得する
func (m *DB) getTxReceiptByUnsignedHex(tbl, hex string) (int64, error) {
	var count int64
	sql := "SELECT count(id) FROM %s WHERE unsigned_hex_tx=?"
	sql = fmt.Sprintf(sql, tbl)

	//err := m.RDB.Get(&count, "SELECT count(id) FROM tx_receipt WHERE unsigned_hex_tx=?", hex)
	err := m.RDB.Get(&count, sql, hex)

	return count, err
}

// GetTxReceiptByUnsignedHex unsigned_hex_txをキーとしてレコードを取得する
func (m *DB) GetTxReceiptByUnsignedHex(hex string) (int64, error) {
	return m.getTxReceiptByUnsignedHex(m.TableNameReceipt(), hex)
}

// InsertTxReceiptForUnsigned TxReceiptテーブルに未署名トランザクションレコードを作成する
func (m *DB) insertTxReceiptForUnsigned(tbl string, txReceipt *TxReceipt, tx *sqlx.Tx, isCommit bool) (int64, error) {

	sql := `
INSERT INTO %s (unsigned_hex_tx, signed_hex_tx, sent_hash_tx, total_amount, fee, receiver_address, current_tx_type) 
VALUES (:unsigned_hex_tx, :signed_hex_tx, :sent_hash_tx, :total_amount, :fee, :receiver_address, :current_tx_type)
`
	sql = fmt.Sprintf(sql, tbl)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	res, err := tx.NamedExec(sql, txReceipt)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if isCommit {
		tx.Commit()
	}

	id, _ := res.LastInsertId()
	return id, nil
}

// InsertTxReceiptForUnsigned TxReceiptテーブルに未署名トランザクションレコードを作成する
func (m *DB) InsertTxReceiptForUnsigned(txReceipt *TxReceipt, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.insertTxReceiptForUnsigned(m.TableNameReceipt(), txReceipt, tx, isCommit)
}

// updsateTxReceiptForSent TxReceiptテーブルのsigned_hex_tx, sent_hash_txを更新する
func (m *DB) updateTxReceiptForSent(tbl string, txReceipt *TxReceipt, tx *sqlx.Tx, isCommit bool) (int64, error) {
	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	sql := `
UPDATE %s SET signed_hex_tx=:signed_hex_tx, sent_hash_tx=:sent_hash_tx, current_tx_type=:current_tx_type,
 sent_updated_at=:sent_updated_at
 WHERE id=:id
`
	sql = fmt.Sprintf(sql, tbl)

	res, err := tx.NamedExec(sql, txReceipt)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if isCommit {
		tx.Commit()
	}
	affectedNum, _ := res.RowsAffected()

	return affectedNum, nil
}

// UpdateTxReceiptForSent TxReceiptテーブルのsigned_hex_tx, sent_hash_txを更新する
func (m *DB) UpdateTxReceiptForSent(txReceipt *TxReceipt, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.updateTxReceiptForSent(m.TableNameReceipt(), txReceipt, tx, isCommit)
}
