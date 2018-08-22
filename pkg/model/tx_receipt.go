package model

import (
	"time"

	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mf-financial/cayenne_wallet/pkg/enum"
)

//enum.Actionに応じて、テーブルを切り替える

const (
	tableNameReceipt = "tx_receipt"
)

// TxTable tx_receipt/tx_paymentテーブル
type TxTable struct {
	ID                int64      `db:"id"`
	UnsignedHexTx     string     `db:"unsigned_hex_tx"`
	SignedHexTx       string     `db:"signed_hex_tx"`
	SentHashTx        string     `db:"sent_hash_tx"`
	TotalInputAmount  string     `db:"total_input_amount"`  //inputの合計
	TotalOutputAmount string     `db:"total_output_amount"` //outputの合計(input-feeがこの金額になるはず)
	Fee               string     `db:"fee"`
	TxType            uint8      `db:"current_tx_type"`
	UnsignedUpdatedAt *time.Time `db:"unsigned_updated_at"`
	SignedUpdatedAt   *time.Time `db:"signed_updated_at"`
	SentUpdatedAt     *time.Time `db:"sent_updated_at"`
}

// TableNameReceipt tx_receiptテーブル名を返す
func (m *DB) TableNameReceipt() string {
	return tableNameReceipt
}

// getTxReceiptByID TxReceiptテーブルから該当するIDのレコードを返す
func (m *DB) getTxReceiptByID(tbl string, id int64) (*TxTable, error) {
	sql := "SELECT * FROM %s WHERE id=?"
	sql = fmt.Sprintf(sql, tbl)

	txReceipt := TxTable{}
	err := m.RDB.Get(&txReceipt, sql, id)

	return &txReceipt, err
}

// GetTxReceiptByID TxReceiptテーブルから該当するIDのレコードを返す
func (m *DB) GetTxReceiptByID(id int64) (*TxTable, error) {
	return m.getTxReceiptByID(m.TableNameReceipt(), id)
}

// getTxReceiptByUnsignedHex unsigned_hex_txをキーとしてレコード数を取得する
func (m *DB) getTxReceiptCountByUnsignedHex(tbl, hex string) (int64, error) {
	var count int64
	sql := "SELECT count(id) FROM %s WHERE unsigned_hex_tx=?"
	sql = fmt.Sprintf(sql, tbl)

	//err := m.RDB.Get(&count, "SELECT count(id) FROM tx_receipt WHERE unsigned_hex_tx=?", hex)
	err := m.RDB.Get(&count, sql, hex)

	return count, err
}

// GetTxReceiptByUnsignedHex unsigned_hex_txをキーとしてレコード数を取得する
func (m *DB) GetTxReceiptCountByUnsignedHex(hex string) (int64, error) {
	return m.getTxReceiptCountByUnsignedHex(m.TableNameReceipt(), hex)
}

func (m *DB) getSentTxHashOnTxReceipt(tbl string) ([]string, error) {
	var txHashs []string
	sql := "SELECT sent_hash_tx FROM %s WHERE current_tx_type=?"
	sql = fmt.Sprintf(sql, tbl)

	err := m.RDB.Select(&txHashs, sql, enum.TxTypeValue[enum.TxTypeSent])
	if err != nil {
		return nil, err
	}

	return txHashs, nil
}

// GetSentTxHashOnTxReceipt TxReceiptテーブルから送信済ステータスであるsent_hash_txの配列を返す
func (m *DB) GetSentTxHashOnTxReceipt() ([]string, error) {
	return m.getSentTxHashOnTxReceipt(m.TableNameReceipt())
}

// InsertTxReceiptForUnsigned TxReceiptテーブルに未署名トランザクションレコードを作成する
func (m *DB) insertTxReceiptForUnsigned(tbl string, txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {

	sql := `
INSERT INTO %s (unsigned_hex_tx, signed_hex_tx, sent_hash_tx, total_input_amount, total_output_amount, fee, current_tx_type) 
VALUES (:unsigned_hex_tx, :signed_hex_tx, :sent_hash_tx, :total_input_amount, :total_output_amount, :fee, :current_tx_type)
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
func (m *DB) InsertTxReceiptForUnsigned(txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.insertTxReceiptForUnsigned(m.TableNameReceipt(), txReceipt, tx, isCommit)
}

// updsateTxReceiptForSent TxReceiptテーブルのsigned_hex_tx, sent_hash_txを更新する
func (m *DB) updateTxReceiptForSent(tbl string, txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
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
func (m *DB) UpdateTxReceiptForSent(txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.updateTxReceiptForSent(m.TableNameReceipt(), txReceipt, tx, isCommit)
}

// updsateTxReceiptForDone TxReceiptテーブルの該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (m *DB) updateTxReceiptForDone(tbl string, hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	sql := `
UPDATE %s SET current_tx_type=? WHERE sent_hash_tx=?
`
	sql = fmt.Sprintf(sql, tbl)
	res, err := tx.Exec(sql, enum.TxTypeValue[enum.TxTypeDone], hash)
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

// UpdateTxReceiptForDone TxReceiptテーブルの該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (m *DB) UpdateTxReceiptForDone(hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.updateTxReceiptForDone(m.TableNameReceipt(), hash, tx, isCommit)
}
