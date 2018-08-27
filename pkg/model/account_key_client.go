package model

import (
	"github.com/jmoiron/sqlx"
	"time"
)

// AccountKeyClient account_key_clientテーブル
type AccountKeyClient struct {
	ID                 int64      `db:"id"`
	WalletAddress      string     `db:"wallet_address"`
	WalletImportFormat string     `db:"wallet_import_format"`
	Account            string     `db:"account"`
	KeyType            uint8      `db:"key_type"`
	Index              uint32     `db:"index"`
	UpdatedAt          *time.Time `db:"updated_at"`
}

// InsertAccountKeyClient account_key_clientテーブルにレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) InsertAccountKeyClient(accountKeyClients []AccountKeyClient, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO account_key_client (wallet_address, account_from, address_to, amount) 
VALUES (:address_from, :account_from, :address_to, :amount)
`

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, accountKeyClient := range accountKeyClients {
		_, err := tx.NamedExec(sql, accountKeyClient)
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
