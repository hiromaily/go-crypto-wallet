package model

import (
	"fmt"
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/jmoiron/sqlx"
)

// AddedPubkeyHistoryTable added_pubkey_history_receiptテーブル
type AddedPubkeyHistoryTable struct {
	ID                    int64      `db:"id"`
	WalletAddress         string     `db:"wallet_address"`
	AuthAddress1          string     `db:"auth_address1"`
	AuthAddress2          string     `db:"auth_address2"`
	WalletMultisigAddress string     `db:"wallet_multisig_address"`
	RedeemScript          string     `db:"redeem_script"`
	UpdatedAt             *time.Time `db:"updated_at"`
}

var addedPubkeyHistoryTableName = map[enum.AccountType]string{
	enum.AccountTypeReceipt: "added_pubkey_history_receipt",
	enum.AccountTypePayment: "added_pubkey_history_payment",
}

// insertAddedPubkeyHistoryTable added_pubkey_history_table(payment, receipt...)テーブルにレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) insertAddedPubkeyHistoryTable(tbl string, addedPubkeyHistoryTables []AddedPubkeyHistoryTable, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (wallet_address, auth_address1, auth_address2, wallet_multisig_address, redeem_script) 
VALUES (:wallet_address, :auth_address1, :auth_address2, :wallet_multisig_address, :redeem_script)
`
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	for _, addedPubkeyHistory := range addedPubkeyHistoryTables {
		_, err := tx.NamedExec(sql, addedPubkeyHistory)
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

// InsertAddedPubkeyHistoryTable added_pubkey_history_table(payment, receipt...)テーブルにレコードを作成する
func (m *DB) InsertAddedPubkeyHistoryTable(accountType enum.AccountType, addedPubkeyHistoryTables []AddedPubkeyHistoryTable, tx *sqlx.Tx, isCommit bool) error {
	return m.insertAddedPubkeyHistoryTable(addedPubkeyHistoryTableName[accountType], addedPubkeyHistoryTables, tx, isCommit)
}
