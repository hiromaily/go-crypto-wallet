package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type TxReceiptDetail struct {
	ID                 int64      `db:"id"`
	ReceiptID          int64      `db:"receipt_id"`
	InputTxid          string     `db:"input_txid"`
	InputVout          uint32     `db:"input_vout"`
	InputAddress       string     `db:"input_address"`
	InputAccount       string     `db:"input_account"`
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
func (m *DB) InsertTxReceiptDetailForUnsigned(txReceiptDetails []TxReceiptDetail, tx *sqlx.Tx) error {
	if tx == nil {
		tx = m.DB.MustBegin()
	}

	sql := `
INSERT INTO tx_receipt_detail (receipt_id, input_txid, input_vout, input_address, input_account, input_amount, input_confirmations) 
VALUES (:receipt_id, :input_txid, :input_vout, :input_address, :input_account, :input_amount, :input_confirmations)
`
	for _, txReceiptDetail := range txReceiptDetails {
		_, err := tx.NamedExec(sql, txReceiptDetail)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
