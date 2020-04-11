package walletrepo

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
)

var txInputTableName = map[enum.ActionType]string{
	"receipt":  "tx_receipt_input",
	"payment":  "tx_payment_input",
	"transfer": "tx_transfer_input",
}

// TxInput tx_receipt_input/tx_payment_inputテーブル
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

// getTxInputByReceiptID 該当するIDのレコードを返す
func (r *WalletRepository) getTxInputByReceiptID(tbl string, receiptID int64) ([]TxInput, error) {
	sql := "SELECT * FROM %s WHERE receipt_id=?"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	var txReceiptInputs []TxInput
	err := r.db.Select(&txReceiptInputs, sql, receiptID)

	return txReceiptInputs, err
}

// GetTxInputByReceiptID 該当するIDのレコードを返す
func (r *WalletRepository) GetTxInputByReceiptID(actionType enum.ActionType, receiptID int64) ([]TxInput, error) {
	return r.getTxInputByReceiptID(txInputTableName[actionType], receiptID)
}

// insertTxInputForUnsigned 未署名トランザクションのinputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (r *WalletRepository) insertTxInputForUnsigned(tbl string, txReceiptInputs []TxInput, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (receipt_id, input_txid, input_vout, input_address, input_account, input_amount, input_confirmations) 
VALUES (:receipt_id, :input_txid, :input_vout, :input_address, :input_account, :input_amount, :input_confirmations)
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = r.db.MustBegin()
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

// InsertTxInputForUnsigned 未署名トランザクションのinputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (r *WalletRepository) InsertTxInputForUnsigned(actionType enum.ActionType, txReceiptInputs []TxInput, tx *sqlx.Tx, isCommit bool) error {
	return r.insertTxInputForUnsigned(txInputTableName[actionType], txReceiptInputs, tx, isCommit)
}
