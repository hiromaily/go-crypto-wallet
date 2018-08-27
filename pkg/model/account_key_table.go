package model

import (
	"fmt"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/jmoiron/sqlx"
	"time"
)

// AccountKeyTable account_key_clientテーブル
type AccountKeyTable struct {
	ID                 int64      `db:"id"`
	WalletAddress      string     `db:"wallet_address"`
	WalletImportFormat string     `db:"wallet_import_format"`
	Account            string     `db:"account"`
	KeyType            uint8      `db:"key_type"`
	Idx                uint32     `db:"idx"`
	UpdatedAt          *time.Time `db:"updated_at"`
}

var accountKeyTableName = map[key.AccountType]string{
	0: "account_key_client",
	1: "account_key_receipt",
	2: "account_key_payment",
	3: "account_key_authorization",
}

// insertAccountKeyClient account_key_table(client, payment, receipt...)テーブルにレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) insertAccountKeyTable(tbl string, accountKeyTables []AccountKeyTable, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (wallet_address, wallet_import_format, account, key_type, idx) 
VALUES (:wallet_address, :wallet_import_format, :account, :key_type, :idx)
`
	sql = fmt.Sprintf(sql, tbl)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, accountKeyClient := range accountKeyTables {
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

// InsertAccountKeyTable account_key_table(client, payment, receipt...)テーブルにレコードを作成する
func (m *DB) InsertAccountKeyTable(accountType key.AccountType, accountKeyTables []AccountKeyTable, tx *sqlx.Tx, isCommit bool) error {
	return m.insertAccountKeyTable(accountKeyTableName[accountType], accountKeyTables, tx, isCommit)
}

//getMaxIndex indexの最大値を返す
func (m *DB) getMaxIndex(tbl string) (int64, error) {
	sql := "SELECT MAX(idx) from %s;"
	sql = fmt.Sprintf(sql, tbl)

	var idx int64
	err := m.RDB.Get(&idx, sql)

	return idx, err
}

//GetMaxIndex indexの最大値を返す
func (m *DB) GetMaxIndex(accountType key.AccountType) (int64, error) {
	return m.getMaxIndex(accountKeyTableName[accountType])
}
