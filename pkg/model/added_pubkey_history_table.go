package model

import (
	"fmt"
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AddedPubkeyHistoryTable added_pubkey_historyテーブル
type AddedPubkeyHistoryTable struct {
	ID                    int64      `db:"id"`
	FullPublicKey         string     `db:"full_public_key"`
	AuthAddress1          string     `db:"auth_address1"`
	AuthAddress2          string     `db:"auth_address2"`
	WalletMultisigAddress string     `db:"wallet_multisig_address"`
	RedeemScript          string     `db:"redeem_script"`
	IsExported            bool       `db:"is_exported"`
	UpdatedAt             *time.Time `db:"updated_at"`
}

var addedPubkeyHistoryTableName = map[enum.AccountType]string{
	enum.AccountTypeReceipt: "added_pubkey_history_receipt",
	enum.AccountTypePayment: "added_pubkey_history_payment",
	enum.AccountTypeQuoine:  "added_pubkey_history_quoine",
	enum.AccountTypeFee:     "added_pubkey_history_fee",
	enum.AccountTypeStored:  "added_pubkey_history_stored",
}

//getAddedPubkeyHistoryTableByNoWalletMultisigAddress WalletMultisigAddressが発行されていないレコードを返す
func (m *DB) getAddedPubkeyHistoryTableByNoWalletMultisigAddress(tbl string, accountType enum.AccountType) ([]AddedPubkeyHistoryTable, error) {
	sql := "SELECT * FROM %s WHERE wallet_multisig_address = '';"
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var addedPubkeyHistoryTable []AddedPubkeyHistoryTable
	err := m.RDB.Select(&addedPubkeyHistoryTable, sql)
	if err != nil {
		return nil, err
	}

	return addedPubkeyHistoryTable, nil
}

//GetAddedPubkeyHistoryTableByNoWalletMultisigAddress WalletMultisigAddressが発行されていないレコードを返す
func (m *DB) GetAddedPubkeyHistoryTableByNoWalletMultisigAddress(accountType enum.AccountType) ([]AddedPubkeyHistoryTable, error) {
	return m.getAddedPubkeyHistoryTableByNoWalletMultisigAddress(addedPubkeyHistoryTableName[accountType], accountType)
}

//getAddedPubkeyHistoryTableByNoWalletMultisigAddress WalletMultisigAddressが発行済かつ、exportされていないレコードを返す
func (m *DB) getAddedPubkeyHistoryTableByNotExported(tbl string, accountType enum.AccountType) ([]AddedPubkeyHistoryTable, error) {
	sql := "SELECT * FROM %s WHERE wallet_multisig_address != '' AND is_exported=false;"
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var addedPubkeyHistoryTable []AddedPubkeyHistoryTable
	err := m.RDB.Select(&addedPubkeyHistoryTable, sql)
	if err != nil {
		return nil, err
	}

	return addedPubkeyHistoryTable, nil
}

//GetAddedPubkeyHistoryTableByNotExported WalletMultisigAddressが発行済かつ、exportされていないレコードを返す
func (m *DB) GetAddedPubkeyHistoryTableByNotExported(accountType enum.AccountType) ([]AddedPubkeyHistoryTable, error) {
	return m.getAddedPubkeyHistoryTableByNotExported(addedPubkeyHistoryTableName[accountType], accountType)
}

// insertAddedPubkeyHistoryTable added_pubkey_history_table(payment, receipt...)テーブルにレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) insertAddedPubkeyHistoryTable(tbl string, addedPubkeyHistoryTables []AddedPubkeyHistoryTable, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (full_public_key, auth_address1, auth_address2, wallet_multisig_address, redeem_script) 
VALUES (:full_public_key, :auth_address1, :auth_address2, :wallet_multisig_address, :redeem_script)
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

// updateMultisigAddrOnAddedPubkeyHistoryTable added_pubkey_history_table(payment, receipt...)テーブルのmultisigアドレスを更新する
func (m *DB) updateMultisigAddrOnAddedPubkeyHistoryTable(tbl, multiSigAddr, redeemScript, authAddr1, fullPublicKey string, tx *sqlx.Tx, isCommit bool) error {
	sql := `
UPDATE %s SET wallet_multisig_address=?, redeem_script=?, auth_address1=? WHERE full_public_key=? 
`
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	_, err := tx.Exec(sql, multiSigAddr, redeemScript, authAddr1, fullPublicKey)
	if err != nil {
		tx.Rollback()
		return err
	}
	if isCommit {
		tx.Commit()
	}
	//affectedNum, _ := res.RowsAffected()

	return nil
}

// UpdateMultisigAddrOnAddedPubkeyHistoryTable added_pubkey_history_table(payment, receipt...)テーブルのmultisigアドレスを更新する
func (m *DB) UpdateMultisigAddrOnAddedPubkeyHistoryTable(accountType enum.AccountType, multiSigAddr, redeemScript, authAddr1, fullPublicKey string, tx *sqlx.Tx, isCommit bool) error {
	return m.updateMultisigAddrOnAddedPubkeyHistoryTable(addedPubkeyHistoryTableName[accountType], multiSigAddr, redeemScript, authAddr1, fullPublicKey, tx, isCommit)
}

// updateIsExportedOnAddedPubkeyHistoryTable added_pubkey_history_table(payment, receipt...)テーブルのis_exportedを更新する
func (m *DB) updateIsExportedOnAddedPubkeyHistoryTable(tbl string, ids []int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
	var sql string
	sql = "UPDATE %s SET is_exported=true WHERE id IN (?);"
	sql = fmt.Sprintf(sql, tbl)

	//In対応
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return 0, errors.Errorf("sqlx.In() error: %v", err)
	}
	query = m.RDB.Rebind(query)
	logger.Debugf("sql: %s", query)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	res, err := tx.Exec(query, args...)
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

// UpdateIsExportedOnAddedPubkeyHistoryTable added_pubkey_history_table(payment, receipt...)テーブルのis_exportedを更新する
func (m *DB) UpdateIsExportedOnAddedPubkeyHistoryTable(accountType enum.AccountType, ids []int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.updateIsExportedOnAddedPubkeyHistoryTable(addedPubkeyHistoryTableName[accountType], ids, tx, isCommit)
}
