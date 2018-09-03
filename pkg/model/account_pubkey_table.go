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
	UpdatedAt     *time.Time `db:"updated_at"`
}

var accountPubKeyTableName = map[enum.AccountType]string{
	enum.AccountTypeClient:  "account_pubkey_client",
	enum.AccountTypeReceipt: "account_pubkey_receipt",
	enum.AccountTypePayment: "account_pubkey_payment",
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
