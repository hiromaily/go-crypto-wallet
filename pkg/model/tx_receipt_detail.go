package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type TxReceiptDetail struct {
	ID           int64  `db:"id"`
	ReceiptID    int64  `db:"receipt_id"`
	InputTxid    string `db:"input_txid"`
	InputVout    uint32 `db:"input_vout"`
	InputAddress string `db:"input_address"`
	InputAccount string `db:"input_account"`
	//TotalAmount   float64 `db:"total_amount"` //Float型はInsert後に誤差が生じる可能性がある
	//Fee           float64 `db:"fee"`
	InputAmount        string     `db:"input_amount"`
	InputConfirmations int64      `db:"input_confirmations"`
	UpdatedAt          *time.Time `db:"updated_at"`
}

// GetTxReceiptByID TxReceiptテーブルから該当するIDのレコードを返す
func (m *DB) GetTxReceiptDetailByReceiptID(receiptID int64) (*TxReceiptDetail, error) {
	txReceiptDetail := TxReceiptDetail{}
	err := m.DB.Select(&txReceiptDetail, "SELECT * FROM tx_receipt_detail WHERE receipt_id=$1", receiptID)

	return &txReceiptDetail, err
}

// InsertTxReceiptDetailForUnsigned TxReceiptDetailテーブルに未署名トランザクションのinputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) InsertTxReceiptDetailForUnsigned(txReceiptDetail []TxReceiptDetail, tx *sqlx.Tx) (int64, error) {
	if tx == nil {
		tx = m.DB.MustBegin()
	}

	sql := `
INSERT INTO tx_receipt (unsigned_hex_tx, signed_hex_tx, sent_hex_tx, total_amount, fee, receiver_address, current_tx_type) 
VALUES (:unsigned_hex_tx, :signed_hex_tx, :sent_hex_tx, :total_amount, :fee, :receiver_address, :current_tx_type)
`
	res, err := tx.NamedExec(sql, txReceiptDetail)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	id, _ := res.LastInsertId()
	return id, nil
}
