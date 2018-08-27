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
	Idx                uint32     `db:"idx"`
	UpdatedAt          *time.Time `db:"updated_at"`
}

// InsertAccountKeyClient account_key_clientテーブルにレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) InsertAccountKeyClient(accountKeyClients []AccountKeyClient, tx *sqlx.Tx, isCommit bool) error {

	sql := "INSERT INTO account_key_client (wallet_address, wallet_import_format, account, key_type, idx) "
	sql += "VALUES (:wallet_address, :wallet_import_format, :account, :key_type, :idx)"

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

//GetMaxClientIndex indexの最大値を返す
func (m *DB) GetMaxClientIndex() (int64, error) {
	sql := "SELECT MAX(idx) from account_key_client;"

	var idx int64
	err := m.RDB.Get(&idx, sql)

	return idx, err
}
