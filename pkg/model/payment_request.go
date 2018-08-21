package model

import (
	"github.com/jmoiron/sqlx"
	"time"
)

//PaymentRequest payment_requestテーブル
type PaymentRequest struct {
	ID          int        `db:"id"`
	AddressFrom string     `db:"address_from"`
	AccountFrom string     `db:"account_from"`
	AddressTo   string     `db:"address_to"`
	Amount      string     `db:"amount"`
	IsDone      bool       `db:"is_done"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

// GetPaymentRequest PaymentRequestテーブル全体を返す
func (m *DB) GetPaymentRequest() ([]PaymentRequest, error) {
	var paymentRequests []PaymentRequest
	err := m.RDB.Select(&paymentRequests, "SELECT * FROM payment_request WHERE is_done=false")

	return paymentRequests, err
}

// InsertPaymentRequest PaymentRequestテーブルに出金依頼レコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) InsertPaymentRequest(paymentRequests []PaymentRequest, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO payment_request (address_from, account_from, address_to, amount) 
VALUES (:address_from, :account_from, :address_to, :amount)
`

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, paymentRequest := range paymentRequests {
		_, err := tx.NamedExec(sql, paymentRequest)
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

// UpdatePaymentRequestForIsDone is_doneフィールドをtrueに更新する
//TODO:暫定で追加したのみ、実際の仕様に合わせて修正が必要
func (m *DB) UpdatePaymentRequestForIsDone(tx *sqlx.Tx, isCommit bool) (int64, error) {
	sql := `
UPDATE payment_request SET is_done=true WHERE is_done=false 
`

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	//res, err := tx.NamedExec(sql, nil)
	//res := tx.MustExec(sql)
	res, err := tx.Exec(sql)
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
