package walletrepo

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
)

var txOutputTableName = map[enum.ActionType]string{
	"receipt":  "tx_receipt_output",
	"payment":  "tx_payment_output",
	"transfer": "tx_transfer_output",
}

// TxOutput tx_receipt_output/tx_payment_outputテーブル
type TxOutput struct {
	ID            int64      `db:"id"`
	ReceiptID     int64      `db:"receipt_id"`
	OutputAddress string     `db:"output_address"`
	OutputAccount string     `db:"output_account"`
	OutputAmount  string     `db:"output_amount"`
	IsChange      bool       `db:"is_change"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

// getTxOutputByReceiptID 該当するIDのレコードを返す
func (r *WalletRepository) getTxOutputByReceiptID(tbl string, receiptID int64) ([]TxOutput, error) {
	sql := "SELECT * FROM %s WHERE receipt_id=?"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	var txReceiptOutputs []TxOutput
	err := r.db.Select(&txReceiptOutputs, sql, receiptID)

	return txReceiptOutputs, err
}

// GetTxOutputByReceiptID 該当するIDのレコードを返す
func (r *WalletRepository) GetTxOutputByReceiptID(actionType enum.ActionType, receiptID int64) ([]TxOutput, error) {
	return r.getTxOutputByReceiptID(txOutputTableName[actionType], receiptID)
}

// insertTxOutputForUnsigned 未署名トランザクションのoutputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (r *WalletRepository) insertTxOutputForUnsigned(tbl string, txReceiptOutputs []TxOutput, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (receipt_id, output_address, output_account, output_amount, is_change) 
VALUES (:receipt_id,  :output_address, :output_account, :output_amount, :is_change)
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = r.db.MustBegin()
	}

	for _, txReceiptOutput := range txReceiptOutputs {
		_, err := tx.NamedExec(sql, txReceiptOutput)
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

// InsertTxOutputForUnsigned 未署名トランザクションのoutputに使われたtxレコードを作成する
//TODO:BulkInsertがやりたい
func (r *WalletRepository) InsertTxOutputForUnsigned(actionType enum.ActionType, txReceiptOutputs []TxOutput, tx *sqlx.Tx, isCommit bool) error {
	return r.insertTxOutputForUnsigned(txOutputTableName[actionType], txReceiptOutputs, tx, isCommit)
}
