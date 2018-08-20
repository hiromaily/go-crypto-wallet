package model

import (
	"time"

	"fmt"
	"github.com/jmoiron/sqlx"
)

//enum.Actionに応じて、テーブルを切り替える

const (
	tableNameReceiptDetail = "tx_receipt_detail"
	tableNamePaymentDetail = "tx_payment_detail"
)

// TxReceiptDetail tx_receipt_detailテーブル(tx_payment_detailとしても利用)
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

// TableNameReceiptDetail tx_receipt_detailテーブル名を返す
func (m *DB) TableNameReceiptDetail() string {
	return tableNameReceiptDetail
}

// TableNamePaymentDetail tx_payment_detailテーブル名を返す
func (m *DB) TableNamePaymentDetail() string {
	return tableNamePaymentDetail
}

// getTxReceiptDetailByReceiptID TxReceiptDetailテーブルから該当するIDのレコードを返す
func (m *DB) getTxReceiptDetailByReceiptID(tbl string, receiptID int64) (*TxReceiptDetail, error) {
	sql := "SELECT * FROM %s WHERE receipt_id=$1"
	sql = fmt.Sprintf(sql, tbl)

	txReceiptDetail := TxReceiptDetail{}
	err := m.RDB.Select(&txReceiptDetail, sql, receiptID)

	return &txReceiptDetail, err
}

// GetTxReceiptDetailByReceiptID TxReceiptDetailテーブルから該当するIDのレコードを返す
func (m *DB) GetTxReceiptDetailByReceiptID(tbl string, receiptID int64) (*TxReceiptDetail, error) {
	return m.getTxReceiptDetailByReceiptID(m.TableNameReceipt(), receiptID)
}

// insertTxReceiptDetailForUnsigned TxReceiptDetailテーブルに未署名トランザクションのinputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) insertTxReceiptDetailForUnsigned(tbl string, txReceiptDetails []TxReceiptDetail, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (receipt_id, input_txid, input_vout, input_address, input_account, input_amount, input_confirmations) 
VALUES (:receipt_id, :input_txid, :input_vout, :input_address, :input_account, :input_amount, :input_confirmations)
`
	sql = fmt.Sprintf(sql, tbl)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, txReceiptDetail := range txReceiptDetails {
		_, err := tx.NamedExec(sql, txReceiptDetail)
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

// InsertTxReceiptDetailForUnsigned TxReceiptDetailテーブルに未署名トランザクションのinputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) InsertTxReceiptDetailForUnsigned(tbl string, txReceiptDetails []TxReceiptDetail, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (receipt_id, input_txid, input_vout, input_address, input_account, input_amount, input_confirmations) 
VALUES (:receipt_id, :input_txid, :input_vout, :input_address, :input_account, :input_amount, :input_confirmations)
`
	sql = fmt.Sprintf(sql, tbl)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, txReceiptDetail := range txReceiptDetails {
		_, err := tx.NamedExec(sql, txReceiptDetail)
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
