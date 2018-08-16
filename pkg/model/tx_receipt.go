package model

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type TxReceipt struct {
	ID            int    `db:"id"`
	UnsignedHexTx string `db:"unsigned_hex_tx"`
	SignedHexTx   string `db:"signed_hex_tx"`
	SentHexTx     string `db:"sent_hex_tx"`
	//TotalAmount   float64 `db:"total_amount"` //Float型はInsert後に誤差が生じる可能性がある
	//Fee           float64 `db:"fee"`
	TotalAmount       string     `db:"total_amount"`
	Fee               string     `db:"fee"`
	ReceiverAddress   string     `db:"receiver_address"`
	TxType            int        `db:"current_tx_type"`
	UnsignedUpdatedAt *time.Time `db:"unsigned_updated_at"`
	SignedUpdatedAt   *time.Time `db:"unsigned_updated_at"`
	SentUpdatedAt     *time.Time `db:"sent_updated_at"`
}

// GetTxReceiptByID TxReceiptテーブルから該当するIDのレコードを返す
func (m *DB) GetTxReceiptByID(id int) (*TxReceipt, error) {
	txReceipt := TxReceipt{}
	err := m.DB.Get(&txReceipt, "SELECT * FROM tx_receipt WHERE id=$1", id)

	return &txReceipt, err
}

// InsertTxReceiptForUnsigned TxReceiptテーブルに未署名トランザクションレコードを作成する
func (m *DB) InsertTxReceiptForUnsigned(txReceipt *TxReceipt, tx *sqlx.Tx) (int64, error) {
	if tx == nil {
		tx = m.DB.MustBegin()
	}
	//sql := "INSERT INTO tx_receipt (unsigned_hex_tx, signed_hex_tx, sent_hex_tx, total_amount, fee, receiver_address, current_tx_type) VALUES ('$1', '$2', '$3', '$4', '$5', '$6', '$7')"
	//res := tx.MustExec(sql, txReceipt.UnsignedHexTx, "", "", txReceipt.TotalAmount, txReceipt.Fee, txReceipt.ReceiverAddress, txReceipt.TxType)
	//if _, err := res.LastInsertId(); err != nil {
	//	tx.Rollback()
	//	return err
	//}
	sql := "INSERT INTO tx_receipt (unsigned_hex_tx, signed_hex_tx, sent_hex_tx, total_amount, fee, receiver_address, current_tx_type) VALUES (:unsigned_hex_tx, :signed_hex_tx, :sent_hex_tx, :total_amount, :fee, :receiver_address, :current_tx_type)"
	res, err := tx.NamedExec(sql, txReceipt)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	id, _ := res.LastInsertId()
	return id, nil
}
