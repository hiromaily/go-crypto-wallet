package model

import (
	"time"

	"fmt"
	"github.com/jmoiron/sqlx"
)

//enum.Actionに応じて、テーブルを切り替える

const (
	tableNameReceiptInput = "tx_receipt_input"
)

// TxReceiptInput tx_receipt_inputテーブル
type TxInput struct {
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

// TableNameReceiptInput tx_receipt_inputテーブル名を返す
func (m *DB) TableNameReceiptInput() string {
	return tableNameReceiptInput
}

// getTxReceiptInputByReceiptID TxReceiptInputテーブルから該当するIDのレコードを返す
func (m *DB) getTxReceiptInputByReceiptID(tbl string, receiptID int64) (*TxInput, error) {
	sql := "SELECT * FROM %s WHERE receipt_id=$1"
	sql = fmt.Sprintf(sql, tbl)

	txReceiptInput := TxInput{}
	err := m.RDB.Select(&txReceiptInput, sql, receiptID)

	return &txReceiptInput, err
}

// GetTxReceiptInputByReceiptID TxReceiptInputテーブルから該当するIDのレコードを返す
func (m *DB) GetTxReceiptInputByReceiptID(receiptID int64) (*TxInput, error) {
	return m.getTxReceiptInputByReceiptID(m.TableNameReceiptInput(), receiptID)
}

// insertTxReceiptInputForUnsigned TxReceiptInputテーブルに未署名トランザクションのinputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) insertTxReceiptInputForUnsigned(tbl string, txReceiptInputs []TxInput, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (receipt_id, input_txid, input_vout, input_address, input_account, input_amount, input_confirmations) 
VALUES (:receipt_id, :input_txid, :input_vout, :input_address, :input_account, :input_amount, :input_confirmations)
`
	sql = fmt.Sprintf(sql, tbl)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, txReceiptInput := range txReceiptInputs {
		_, err := tx.NamedExec(sql, txReceiptInput)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if isCommit {
		tx.Commit()
	}

	return nil
}

// InsertTxReceiptInputForUnsigned TxReceiptInputテーブルに未署名トランザクションのinputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) InsertTxReceiptInputForUnsigned(txReceiptInputs []TxInput, tx *sqlx.Tx, isCommit bool) error {
	return m.insertTxReceiptInputForUnsigned(m.TableNameReceiptInput(), txReceiptInputs, tx, isCommit)
}
