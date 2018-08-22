package model

import (
	"fmt"
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/jmoiron/sqlx"
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

	err := m.RDB.Get(&count, sql, hex)

	return count, err
}

// GetTxReceiptCountByUnsignedHex unsigned_hex_txをキーとしてレコード数を取得する
func (m *DB) GetTxReceiptCountByUnsignedHex(hex string) (int64, error) {
	return m.getTxReceiptCountByUnsignedHex(m.TableNameReceipt(), hex)
}

// getTxReceiptIDBySentHash sent_hash_txをキーとしてreceipt_idを取得する
func (m *DB) getTxReceiptIDBySentHash(tbl, hash string) (int64, error) {
	var receiptID int64
	sql := "SELECT id FROM %s WHERE sent_hash_tx=?"
	sql = fmt.Sprintf(sql, tbl)

	err := m.RDB.Get(&receiptID, sql, hash)

	return receiptID, err
}

// GetTxReceiptIDBySentHash sent_hash_txをキーとしてreceipt_idを取得する
func (m *DB) GetTxReceiptIDBySentHash(hash string) (int64, error) {
	return m.getTxReceiptIDBySentHash(m.TableNameReceipt(), hash)
}

func (m *DB) getSentTxHashOnTxReceipt(tbl string, txTypeValue uint8) ([]string, error) {
	var txHashs []string
	sql := "SELECT sent_hash_tx FROM %s WHERE current_tx_type=?"
	sql = fmt.Sprintf(sql, tbl)

	err := m.RDB.Select(&txHashs, sql, enum.TxTypeValue[enum.TxTypeSent])
	if err != nil {
		return nil, err
	}

	return txHashs, nil
}

// GetSentTxHashOnTxReceiptByTxTypeSent TxReceiptテーブルから送信済ステータスであるsent_hash_txの配列を返す
func (m *DB) GetSentTxHashOnTxReceiptByTxTypeSent() ([]string, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeSent]
	return m.getSentTxHashOnTxReceipt(m.TableNameReceipt(), txTypeValue)
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

// updateTxTypeOnTxReceiptByTxHash TxReceiptテーブルの該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (m *DB) updateTxTypeOnTxReceiptByTxHash(tbl string, hash string, txTypeValue uint8, tx *sqlx.Tx, isCommit bool) (int64, error) {

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	sql := `
UPDATE %s SET current_tx_type=? WHERE sent_hash_tx=?
`
	sql = fmt.Sprintf(sql, tbl)
	res, err := tx.Exec(sql, txTypeValue, hash)
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

// UpdateTxReceipDoneByTxHash TxReceiptテーブルの該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (m *DB) UpdateTxReceipDoneByTxHash(hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeDone]
	return m.updateTxTypeOnTxReceiptByTxHash(m.TableNameReceipt(), hash, txTypeValue, tx, isCommit)
}

// UpdateTxReceipNotifiedByTxHash TxReceiptテーブルの該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (m *DB) UpdateTxReceipNotifiedByTxHash(hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeNotified]
	return m.updateTxTypeOnTxReceiptByTxHash(m.TableNameReceipt(), hash, txTypeValue, tx, isCommit)
}

// updateTxTypeOnTxReceiptByID TxReceiptテーブルの該当するIDのレコードのcurrnt_tx_typeを更新する
func (m *DB) updateTxTypeOnTxReceiptByID(tbl string, ID int64, txTypeValue uint8, tx *sqlx.Tx, isCommit bool) (int64, error) {

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	sql := `
UPDATE %s SET current_tx_type=? WHERE id=?
`
	sql = fmt.Sprintf(sql, tbl)
	res, err := tx.Exec(sql, txTypeValue, ID)
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

// UpdateTxReceipDoneByID TxReceiptテーブルの該当するIDのレコードのcurrnt_tx_typeを更新する
func (m *DB) UpdateTxReceipDoneByID(ID int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeDone]
	return m.updateTxTypeOnTxReceiptByID(m.TableNameReceipt(), ID, txTypeValue, tx, isCommit)
}

// UpdateTxReceipNotifiedByID TxReceiptテーブルの該当するIDのレコードのcurrnt_tx_typeを更新する
func (m *DB) UpdateTxReceipNotifiedByID(ID int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeNotified]
	return m.updateTxTypeOnTxReceiptByID(m.TableNameReceipt(), ID, txTypeValue, tx, isCommit)
}
