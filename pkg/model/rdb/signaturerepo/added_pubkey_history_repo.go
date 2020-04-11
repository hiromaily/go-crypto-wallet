package signaturerepo

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
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

var addedPubkeyHistoryTableName = map[account.AccountType]string{
	account.AccountTypeReceipt: "added_pubkey_history_receipt",
	account.AccountTypePayment: "added_pubkey_history_payment",
	account.AccountTypeQuoine:  "added_pubkey_history_quoine",
	account.AccountTypeFee:     "added_pubkey_history_fee",
	account.AccountTypeStored:  "added_pubkey_history_stored",
}

//getAddedPubkeyHistoryTableByNoWalletMultisigAddress WalletMultisigAddressが発行されていないレコードを返す
func (r *SignatureRepository) getAddedPubkeyHistoryTableByNoWalletMultisigAddress(tbl string, accountType account.AccountType) ([]AddedPubkeyHistoryTable, error) {
	sql := "SELECT * FROM %s WHERE wallet_multisig_address = '';"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	var addedPubkeyHistoryTable []AddedPubkeyHistoryTable
	err := r.db.Select(&addedPubkeyHistoryTable, sql)
	if err != nil {
		return nil, err
	}

	return addedPubkeyHistoryTable, nil
}

//GetAddedPubkeyHistoryTableByNoWalletMultisigAddress WalletMultisigAddressが発行されていないレコードを返す
func (r *SignatureRepository) GetAddedPubkeyHistoryTableByNoWalletMultisigAddress(accountType account.AccountType) ([]AddedPubkeyHistoryTable, error) {
	return r.getAddedPubkeyHistoryTableByNoWalletMultisigAddress(addedPubkeyHistoryTableName[accountType], accountType)
}

//getAddedPubkeyHistoryTableByNoWalletMultisigAddress WalletMultisigAddressが発行済かつ、exportされていないレコードを返す
func (r *SignatureRepository) getAddedPubkeyHistoryTableByNotExported(tbl string, accountType account.AccountType) ([]AddedPubkeyHistoryTable, error) {
	sql := "SELECT * FROM %s WHERE wallet_multisig_address != '' AND is_exported=false;"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	var addedPubkeyHistoryTable []AddedPubkeyHistoryTable
	err := r.db.Select(&addedPubkeyHistoryTable, sql)
	if err != nil {
		return nil, err
	}

	return addedPubkeyHistoryTable, nil
}

//GetAddedPubkeyHistoryTableByNotExported WalletMultisigAddressが発行済かつ、exportされていないレコードを返す
func (r *SignatureRepository) GetAddedPubkeyHistoryTableByNotExported(accountType account.AccountType) ([]AddedPubkeyHistoryTable, error) {
	return r.getAddedPubkeyHistoryTableByNotExported(addedPubkeyHistoryTableName[accountType], accountType)
}

// insertAddedPubkeyHistoryTable added_pubkey_history_table(payment, receipt...)テーブルにレコードを作成する
//TODO:BulkInsertがやりたい
func (r *SignatureRepository) insertAddedPubkeyHistoryTable(tbl string, addedPubkeyHistoryTables []AddedPubkeyHistoryTable, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (full_public_key, auth_address1, auth_address2, wallet_multisig_address, redeem_script) 
VALUES (:full_public_key, :auth_address1, :auth_address2, :wallet_multisig_address, :redeem_script)
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = r.db.MustBegin()
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
func (r *SignatureRepository) InsertAddedPubkeyHistoryTable(accountType account.AccountType, addedPubkeyHistoryTables []AddedPubkeyHistoryTable, tx *sqlx.Tx, isCommit bool) error {
	return r.insertAddedPubkeyHistoryTable(addedPubkeyHistoryTableName[accountType], addedPubkeyHistoryTables, tx, isCommit)
}

// updateMultisigAddrOnAddedPubkeyHistoryTable added_pubkey_history_table(payment, receipt...)テーブルのmultisigアドレスを更新する
func (r *SignatureRepository) updateMultisigAddrOnAddedPubkeyHistoryTable(tbl, multiSigAddr, redeemScript, authAddr1, fullPublicKey string, tx *sqlx.Tx, isCommit bool) error {
	sql := `
UPDATE %s SET wallet_multisig_address=?, redeem_script=?, auth_address1=? WHERE full_public_key=? 
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = r.db.MustBegin()
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
func (r *SignatureRepository) UpdateMultisigAddrOnAddedPubkeyHistoryTable(accountType account.AccountType, multiSigAddr, redeemScript, authAddr1, fullPublicKey string, tx *sqlx.Tx, isCommit bool) error {
	return r.updateMultisigAddrOnAddedPubkeyHistoryTable(addedPubkeyHistoryTableName[accountType], multiSigAddr, redeemScript, authAddr1, fullPublicKey, tx, isCommit)
}

// updateIsExportedOnAddedPubkeyHistoryTable added_pubkey_history_table(payment, receipt...)テーブルのis_exportedを更新する
func (r *SignatureRepository) updateIsExportedOnAddedPubkeyHistoryTable(tbl string, ids []int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
	var sql string
	sql = "UPDATE %s SET is_exported=true WHERE id IN (?);"
	sql = fmt.Sprintf(sql, tbl)

	//In対応
	query, args, err := sqlx.In(sql, ids)
	if err != nil {
		return 0, errors.Errorf("sqlx.In() error: %v", err)
	}
	query = r.db.Rebind(query)
	//logger.Debugf("sql: %s", query)

	if tx == nil {
		tx = r.db.MustBegin()
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
func (r *SignatureRepository) UpdateIsExportedOnAddedPubkeyHistoryTable(accountType account.AccountType, ids []int64, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return r.updateIsExportedOnAddedPubkeyHistoryTable(addedPubkeyHistoryTableName[accountType], ids, tx, isCommit)
}
