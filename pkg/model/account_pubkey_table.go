package model

import (
	"fmt"
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/jmoiron/sqlx"
)

// AccountPublicKeyTable account_key_clientテーブル
type AccountPublicKeyTable struct {
	ID            int64      `db:"id"`
	WalletAddress string     `db:"wallet_address"`
	Account       string     `db:"account"`
	IsAllocated   bool       `db:"is_allocated"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

var accountPubKeyTableName = map[enum.AccountType]string{
	enum.AccountTypeClient:  "account_pubkey_client",
	enum.AccountTypeReceipt: "account_pubkey_receipt",
	enum.AccountTypePayment: "account_pubkey_payment",
}

//getAllAccountPubKeyTable
func (m *DB) getAllAccountPubKeyTable(tbl string) ([]AccountPublicKeyTable, error) {
	sql := "SELECT * FROM %s;"
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var accountKeyTable []AccountPublicKeyTable
	err := m.RDB.Select(&accountKeyTable, sql)
	if err != nil {
		return nil, err
	}

	return accountKeyTable, nil
}

// GetAllAccountPubKeyTable account_pubkey_table(client, payment, receipt...)テーブルから全レコードを取得
func (m *DB) GetAllAccountPubKeyTable(accountType enum.AccountType) ([]AccountPublicKeyTable, error) {
	return m.getAllAccountPubKeyTable(accountPubKeyTableName[accountType])
}

//getOneUnAllocatedAccountPubKeyTable account_pubkey_table(client, payment, receipt...)テーブルからis_allocated=falseの1レコードを取得
func (m *DB) getOneUnAllocatedAccountPubKeyTable(tbl string) (*AccountPublicKeyTable, error) {
	sql := "SELECT * FROM %s WHERE is_allocated=false ORDER BY id LIMIT 1;"
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var accountKeyTable AccountPublicKeyTable
	err := m.RDB.Get(&accountKeyTable, sql)
	if err != nil {
		return nil, err
	}

	return &accountKeyTable, nil
}

// GetOneUnAllocatedAccountPubKeyTable account_pubkey_table(client, payment, receipt...)テーブルからis_allocated=falseの1レコードを取得
func (m *DB) GetOneUnAllocatedAccountPubKeyTable(accountType enum.AccountType) (*AccountPublicKeyTable, error) {
	return m.getOneUnAllocatedAccountPubKeyTable(accountPubKeyTableName[accountType])
}

// insertAccountPubKeyTable account_key_table(client, payment, receipt...)テーブルにレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) insertAccountPubKeyTable(tbl string, accountPubKeyTables []AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (wallet_address, account) 
VALUES (:wallet_address, :account)
`
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, accountPubKeyTable := range accountPubKeyTables {
		_, err := tx.NamedExec(sql, accountPubKeyTable)
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

// InsertAccountPubKeyTable account_pubkey_table(client, payment, receipt...)テーブルにレコードを作成する
func (m *DB) InsertAccountPubKeyTable(accountType enum.AccountType, accountPubKeyTables []AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error {
	return m.insertAccountPubKeyTable(accountPubKeyTableName[accountType], accountPubKeyTables, tx, isCommit)
}

// updateAccountOnAccountPubKeyTable Accountを更新する
func (m *DB) updateAccountOnAccountPubKeyTable(tbl string, accountKeyTable []AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error {
	sql := `
UPDATE %s SET account=:account, updated_at=:updated_at 
WHERE id=:id
`
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, accountKey := range accountKeyTable {
		_, err := tx.NamedExec(sql, accountKey)
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

// UpdateAccountOnAccountPubKeyTable Accountを更新する
func (m *DB) UpdateAccountOnAccountPubKeyTable(accountType enum.AccountType, accountKeyTable []AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error {
	return m.updateAccountOnAccountPubKeyTable(accountPubKeyTableName[accountType], accountKeyTable, tx, isCommit)
}

// updateIsAllocatedOnAccountPubKeyTable IsAllocatedを更新する
func (m *DB) updateIsAllocatedOnAccountPubKeyTable(tbl string, accountKeyTable []AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error {
	sql := `
UPDATE %s SET is_allocated=:is_allocated, updated_at=:updated_at 
WHERE wallet_address=:wallet_address
`
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, accountKey := range accountKeyTable {
		_, err := tx.NamedExec(sql, accountKey)
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

// UpdateIsAllocatedOnAccountPubKeyTable IsAllocatedを更新する
func (m *DB) UpdateIsAllocatedOnAccountPubKeyTable(accountType enum.AccountType, accountKeyTable []AccountPublicKeyTable, tx *sqlx.Tx, isCommit bool) error {
	return m.updateIsAllocatedOnAccountPubKeyTable(accountPubKeyTableName[accountType], accountKeyTable, tx, isCommit)
}
