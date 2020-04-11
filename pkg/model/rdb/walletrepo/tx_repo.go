package walletrepo

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
)

var txTableName = map[enum.ActionType]string{
	"receipt":  "tx_receipt",
	"payment":  "tx_payment",
	"transfer": "tx_transfer",
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

// getTxByID 該当するIDのレコードを返す
func (r *WalletRepository) getTxByID(tbl string, id int64) (*TxTable, error) {
	sql := "SELECT * FROM %s WHERE id=?"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	txReceipt := TxTable{}
	err := r.db.Get(&txReceipt, sql, id)

	return &txReceipt, err
}

// GetTxByID 該当するIDのレコードを返す
func (r *WalletRepository) GetTxByID(actionType enum.ActionType, id int64) (*TxTable, error) {
	return r.getTxByID(txTableName[actionType], id)
}

// getTxCountByUnsignedHex unsigned_hex_txをキーとしてレコード数を取得する
func (r *WalletRepository) getTxCountByUnsignedHex(tbl, hex string) (int64, error) {
	var count int64
	sql := "SELECT count(id) FROM %s WHERE unsigned_hex_tx=?"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	err := r.db.Get(&count, sql, hex)

	return count, err
}

// GetTxCountByUnsignedHex unsigned_hex_txをキーとしてレコード数を取得する
func (r *WalletRepository) GetTxCountByUnsignedHex(actionType enum.ActionType, hex string) (int64, error) {
	return r.getTxCountByUnsignedHex(txTableName[actionType], hex)
}

// getTxIDBySentHash sent_hash_txをキーとしてreceipt_idを取得する
func (r *WalletRepository) getTxIDBySentHash(tbl, hash string) (int64, error) {
	var receiptID int64
	sql := "SELECT id FROM %s WHERE sent_hash_tx=?"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	err := r.db.Get(&receiptID, sql, hash)

	return receiptID, err
}

// GetTxIDBySentHash sent_hash_txをキーとしてreceipt_idを取得する
func (r *WalletRepository) GetTxIDBySentHash(actionType enum.ActionType, hash string) (int64, error) {
	return r.getTxIDBySentHash(txTableName[actionType], hash)
}

// getSentTxHash
func (r *WalletRepository) getSentTxHash(tbl string, txTypeValue uint8) ([]string, error) {
	var txHashs []string
	sql := "SELECT sent_hash_tx FROM %s WHERE current_tx_type=?"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	err := r.db.Select(&txHashs, sql, txTypeValue)
	if err != nil {
		return nil, err
	}

	return txHashs, nil
}

// GetSentTxHashByTxTypeSent tx_typeが`sent`であるsent_hash_txの配列を返す
func (r *WalletRepository) GetSentTxHashByTxTypeSent(actionType enum.ActionType) ([]string, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeSent]
	return r.getSentTxHash(txTableName[actionType], txTypeValue)
}

// GetSentTxHashByTxTypeDone tx_typeが`done`のステータスであるsent_hash_txの配列を返す
func (r *WalletRepository) GetSentTxHashByTxTypeDone(actionType enum.ActionType) ([]string, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeDone]
	return r.getSentTxHash(txTableName[actionType], txTypeValue)
}

// insertTxForUnsigned 未署名トランザクションレコードを作成する
func (r *WalletRepository) insertTxForUnsigned(tbl string, txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {

	sql := `
INSERT INTO %s (unsigned_hex_tx, signed_hex_tx, sent_hash_tx, total_input_amount, total_output_amount, fee, current_tx_type) 
VALUES (:unsigned_hex_tx, :signed_hex_tx, :sent_hash_tx, :total_input_amount, :total_output_amount, :fee, :current_tx_type)
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = r.db.MustBegin()
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

// InsertTxForUnsigned 未署名トランザクションレコードを作成する
func (r *WalletRepository) InsertTxForUnsigned(actionType enum.ActionType, txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return r.insertTxForUnsigned(txTableName[actionType], txReceipt, tx, isCommit)
}

// updateTxAfterSent signed_hex_tx, sent_hash_txを更新する
func (r *WalletRepository) updateTxAfterSent(tbl string, txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
	if tx == nil {
		tx = r.db.MustBegin()
	}

	sql := `
UPDATE %s SET signed_hex_tx=:signed_hex_tx, sent_hash_tx=:sent_hash_tx, current_tx_type=:current_tx_type,
 sent_updated_at=:sent_updated_at
 WHERE id=:id
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

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

// UpdateTxAfterSent signed_hex_tx, sent_hash_txを更新する
func (r *WalletRepository) UpdateTxAfterSent(actionType enum.ActionType, txReceipt *TxTable, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return r.updateTxAfterSent(txTableName[actionType], txReceipt, tx, isCommit)
}

// updateTxTypeByTxHash 該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (r *WalletRepository) updateTxTypeByTxHash(tbl string, hash string, txTypeValue uint8, tx *sqlx.Tx, isCommit bool) (int64, error) {

	if tx == nil {
		tx = r.db.MustBegin()
	}

	sql := `
UPDATE %s SET current_tx_type=? WHERE sent_hash_tx=?
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

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

// UpdateTxTypeDoneByTxHash 該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (r *WalletRepository) UpdateTxTypeDoneByTxHash(actionType enum.ActionType, hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeDone]
	return r.updateTxTypeByTxHash(txTableName[actionType], hash, txTypeValue, tx, isCommit)
}

// UpdateTxTypeNotifiedByTxHash 該当するsent_hash_txのレコードのcurrnt_tx_typeを更新する
func (r *WalletRepository) UpdateTxTypeNotifiedByTxHash(actionType enum.ActionType, hash string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeNotified]
	return r.updateTxTypeByTxHash(txTableName[actionType], hash, txTypeValue, tx, isCommit)
}

// updateTxTypeByID 該当するIDのレコードのcurrnt_tx_typeを更新する
func (r *WalletRepository) updateTxTypeByID(tbl string, ID int64, txTypeValue uint8, tx *sqlx.Tx, isCommit bool) (int64, error) {

	if tx == nil {
		tx = r.db.MustBegin()
	}

	sql := `
UPDATE %s SET current_tx_type=? WHERE id=?
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

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

// UpdateTxTypeDoneByID 該当するIDのレコードのcurrnt_tx_typeを更新する
func (r *WalletRepository) UpdateTxTypeDoneByID(actionType enum.ActionType, ID int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeDone]
	return r.updateTxTypeByID(txTableName[actionType], ID, txTypeValue, tx, isCommit)
}

// UpdateTxTypeNotifiedByID 該当するIDのレコードのcurrnt_tx_typeを更新する
func (r *WalletRepository) UpdateTxTypeNotifiedByID(actionType enum.ActionType, ID int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
	txTypeValue := enum.TxTypeValue[enum.TxTypeNotified]
	return r.updateTxTypeByID(txTableName[actionType], ID, txTypeValue, tx, isCommit)
}
