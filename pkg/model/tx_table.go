package model

import (
	"fmt"
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/jmoiron/sqlx"
)

//enum.Actionに応じて、テーブルを切り替える
//TODO:パラメータに、enum.ActionTypeを追加して、tx_payment*.goをすべて排除したほうがシンプルかも。メンテは楽

//const (
//	tableNameReceipt = "tx_receipt"
//	tableNamePayment = "tx_payment"
//)

var txTableName = map[enum.ActionType]string{
	"receipt": "tx_receipt",
	"payment": "tx_payment",
}

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
//func (m *DB) TableNameReceipt() string {
//	return tableNameReceipt
//}

// getTxReceiptByID TxReceiptテーブルから該当するIDのレコードを返す
func (m *DB) getTxByID(tbl string, id int64) (*TxTable, error) {
	sql := "SELECT * FROM %s WHERE id=?"
	sql = fmt.Sprintf(sql, tbl)

	txReceipt := TxTable{}
	err := m.RDB.Get(&txReceipt, sql, id)

	return &txReceipt, err
}

// GetTxReceiptByID TxReceiptテーブルから該当するIDのレコードを返す
func (m *DB) GetTxByID(actionType enum.ActionType, id int64) (*TxTable, error) {
	return m.getTxByID(txTableName[actionType], id)
}

// getTxReceiptByUnsignedHex unsigned_hex_txをキーとしてレコード数を取得する
func (m *DB) getTxCountByUnsignedHex(tbl, hex string) (int64, error) {
	var count int64
	sql := "SELECT count(id) FROM %s WHERE unsigned_hex_tx=?"
	sql = fmt.Sprintf(sql, tbl)

	err := m.RDB.Get(&count, sql, hex)

	return count, err
}

// GetTxReceiptCountByUnsignedHex unsigned_hex_txをキーとしてレコード数を取得する
func (m *DB) GetTxCountByUnsignedHex(actionType enum.ActionType, hex string) (int64, error) {
	return m.getTxCountByUnsignedHex(txTableName[actionType], hex)
}

// getTxReceiptIDBySentHash sent_hash_txをキーとしてreceipt_idを取得する
func (m *DB) getTxIDBySentHash(tbl, hash string) (int64, error) {
	var receiptID int64
	sql := "SELECT id FROM %s WHERE sent_hash_tx=?"
	sql = fmt.Sprintf(sql, tbl)

	err := m.RDB.Get(&receiptID, sql, hash)

	return receiptID, err
}

// GetTxReceiptIDBySentHash sent_hash_txをキーとしてreceipt_idを取得する
func (m *DB) GetTxIDBySentHash(actionType enum.ActionType, hash string) (int64, error) {
	return m.getTxIDBySentHash(txTableName[actionType], hash)
}

func (m *DB) getSentTxHash(tbl string, txTypeValue uint8) ([]string, error) {
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
func (m *DB) GetSentTxHashByTxTypeSent(actionType enum.ActionType) ([]string, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeSent]
	return m.getSentTxHash(txTableName[actionType], txTypeValue)
}

// InsertTxReceiptForUnsigned TxReceiptテーブルに未署名トランザクションレコードを作成する
func (m *DB) insertTxForUnsigned(tbl string, txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {

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
func (m *DB) InsertTxForUnsigned(actionType enum.ActionType, txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.insertTxForUnsigned(txTableName[actionType], txReceipt, tx, isCommit)
}

// updsateTxReceiptForSent TxReceiptテーブルのsigned_hex_tx, sent_hash_txを更新する
func (m *DB) updateTxAfterSent(tbl string, txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
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
func (m *DB) UpdateTxAfterSent(actionType enum.ActionType, txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.updateTxAfterSent(txTableName[actionType], txReceipt, tx, isCommit)
}

// updateTxTypeOnTxReceiptByTxHash TxReceiptテーブルの該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (m *DB) updateTxTypeByTxHash(tbl string, hash string, txTypeValue uint8, tx *sqlx.Tx, isCommit bool) (int64, error) {

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
func (m *DB) UpdateTxTypeDoneByTxHash(actionType enum.ActionType, hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeDone]
	return m.updateTxTypeByTxHash(txTableName[actionType], hash, txTypeValue, tx, isCommit)
}

// UpdateTxReceipNotifiedByTxHash TxReceiptテーブルの該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (m *DB) UpdateTxTypeNotifiedByTxHash(actionType enum.ActionType, hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeNotified]
	return m.updateTxTypeByTxHash(txTableName[actionType], hash, txTypeValue, tx, isCommit)
}

// updateTxTypeOnTxReceiptByID TxReceiptテーブルの該当するIDのレコードのcurrnt_tx_typeを更新する
func (m *DB) updateTxTypeByID(tbl string, ID int64, txTypeValue uint8, tx *sqlx.Tx, isCommit bool) (int64, error) {

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
func (m *DB) UpdateTxTypeDoneByID(actionType enum.ActionType, ID int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeDone]
	return m.updateTxTypeByID(txTableName[actionType], ID, txTypeValue, tx, isCommit)
}

// UpdateTxReceipNotifiedByID TxReceiptテーブルの該当するIDのレコードのcurrnt_tx_typeを更新する
func (m *DB) UpdateTxTypeNotifiedByID(actionType enum.ActionType, ID int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeNotified]
	return m.updateTxTypeByID(txTableName[actionType], ID, txTypeValue, tx, isCommit)
}
