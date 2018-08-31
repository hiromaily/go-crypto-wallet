package model

import (
	"fmt"
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// AccountKeyTable account_key_clientテーブル
type AccountKeyTable struct {
	ID                    int64      `db:"id"`
	WalletAddress         string     `db:"wallet_address"`
	P2shSegwitAddress     string     `db:"p2sh_segwit_address"`
	FullPublicKey         string     `db:"full_public_key"`
	WalletMultisigAddress string     `db:"wallet_multisig_address"`
	RedeemScript          string     `db:"redeem_script"`
	WalletImportFormat    string     `db:"wallet_import_format"`
	Account               string     `db:"account"`
	KeyType               uint8      `db:"key_type"`
	Idx                   uint32     `db:"idx"`
	IsImprotedPrivKey     bool       `db:"is_imported_priv_key"`
	IsExprotedPubKey      bool       `db:"is_exported_pub_key"`
	UpdatedAt             *time.Time `db:"updated_at"`
}

var accountKeyTableName = map[enum.AccountType]string{
	enum.AccountTypeClient:        "account_key_client",
	enum.AccountTypeReceipt:       "account_key_receipt",
	enum.AccountTypePayment:       "account_key_payment",
	enum.AccountTypeAuthorization: "account_key_authorization",
}

// insertAccountKeyClient account_key_table(client, payment, receipt...)テーブルにレコードを作成する
//TODO:BulkInsertがやりたい
func (m *DB) insertAccountKeyTable(tbl string, accountKeyTables []AccountKeyTable, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (wallet_address, p2sh_segwit_address, full_public_key, wallet_multisig_address, redeem_script, wallet_import_format, account, key_type, idx) 
VALUES (:wallet_address, :p2sh_segwit_address,:full_public_key, :wallet_multisig_address, :redeem_script, :wallet_import_format, :account, :key_type, :idx)
`
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

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
func (m *DB) InsertAccountKeyTable(accountType enum.AccountType, accountKeyTables []AccountKeyTable, tx *sqlx.Tx, isCommit bool) error {
	return m.insertAccountKeyTable(accountKeyTableName[accountType], accountKeyTables, tx, isCommit)
}

// updateIsImprotedPrivKey is_imported_priv_keyをtrueに更新する
func (m *DB) updateIsImprotedPrivKey(tbl, strWIF string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	sql := `
UPDATE %s SET is_imported_priv_key=true WHERE wallet_import_format=? 
`
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = m.RDB.MustBegin()
	}

	res, err := tx.Exec(sql, strWIF)
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

// UpdateIsImprotedPrivKey is_imported_priv_keyをtrueに更新する
func (m *DB) UpdateIsImprotedPrivKey(accountType enum.AccountType, strWIF string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.updateIsImprotedPrivKey(accountKeyTableName[accountType], strWIF, tx, isCommit)
}

// updateIsExprotedPubKey is_exported_pub_keyをtrueに更新する
func (m *DB) updateIsExprotedPubKey(tbl string, accountType enum.AccountType, pubKeys []string, isMultisig bool, tx *sqlx.Tx, isCommit bool) (int64, error) {
	var sql string
	if accountType == enum.AccountTypeClient {
		sql = "UPDATE %s SET is_exported_pub_key=true WHERE wallet_address IN (?);"
	} else if accountType != enum.AccountTypeClient && isMultisig {
		sql = "UPDATE %s SET is_exported_pub_key=true WHERE wallet_multisig_address IN (?) ;"
	} else {
		logger.Info("is_exported_pub_key is not needed to update")
		return 0, nil
	}
	sql = fmt.Sprintf(sql, tbl)

	//In対応
	query, args, err := sqlx.In(sql, pubKeys)
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

// UpdateIsExprotedPubKey is_exported_pub_keyをtrueに更新する
func (m *DB) UpdateIsExprotedPubKey(accountType enum.AccountType, pubKeys []string, isMultisig bool, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return m.updateIsExprotedPubKey(accountKeyTableName[accountType], accountType, pubKeys, isMultisig, tx, isCommit)
}

//getMaxIndex indexの最大値を返す
func (m *DB) getMaxIndex(tbl string) (int64, error) {
	sql := "SELECT MAX(idx) from %s;"
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var idx int64
	err := m.RDB.Get(&idx, sql)

	return idx, err
}

//GetMaxIndex indexの最大値を返す
func (m *DB) GetMaxIndex(accountType enum.AccountType) (int64, error) {
	return m.getMaxIndex(accountKeyTableName[accountType])
}

//getNotImportedKeyWIF IsImprotedPrivKeyがfalseのレコードをすべて返す
func (m *DB) getNotImportedKeyWIF(tbl string) ([]string, error) {
	sql := "SELECT wallet_import_format FROM %s WHERE is_imported_priv_key=false;"
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var WIFs []string
	err := m.RDB.Select(&WIFs, sql)
	if err != nil {
		return nil, err
	}

	return WIFs, nil
}

//GetNotImportedKeyWIF IsImprotedPrivKeyがfalseのレコードをすべて返す
func (m *DB) GetNotImportedKeyWIF(accountType enum.AccountType) ([]string, error) {
	return m.getNotImportedKeyWIF(accountKeyTableName[accountType])
}

//getNotExportedPubKey IsExprotedPubKeyがfalseのレコードをすべて返す
func (m *DB) getPubkeyNotExportedPubKey(tbl string, accountType enum.AccountType, isMultisig bool) ([]string, error) {
	//wallet_address
	//wallet_multisig_address
	var sql string
	//if accountType == enum.AccountTypeClient {
	if !isMultisig {
		sql = "SELECT wallet_address FROM %s WHERE is_exported_pub_key=false;"
	} else {
		sql = "SELECT wallet_multisig_address FROM %s WHERE is_exported_pub_key=false;"
	}
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var pubKeys []string
	err := m.RDB.Select(&pubKeys, sql)
	if err != nil {
		return nil, err
	}

	return pubKeys, nil
}

//GetPubkeyNotExportedPubKey IsExprotedPubKeyがfalseのレコードをすべて返す
func (m *DB) GetPubkeyNotExportedPubKey(accountType enum.AccountType, isMultisig bool) ([]string, error) {
	return m.getPubkeyNotExportedPubKey(accountKeyTableName[accountType], accountType, isMultisig)
}

//getFullNotExportedPubKey IsExprotedPubKeyがfalseのレコードをすべて返す
func (m *DB) getFullPubkeyNotExportedPubKey(tbl string) ([]string, error) {

	var sql = "SELECT full_public_key FROM %s WHERE is_exported_pub_key=false;"
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var pubKeys []string
	err := m.RDB.Select(&pubKeys, sql)
	if err != nil {
		return nil, err
	}

	return pubKeys, nil
}

//GetFullPubkeyNotExportedPubKey IsExprotedPubKeyがfalseのレコードをすべて返す
func (m *DB) GetFullPubkeyNotExportedPubKey(accountType enum.AccountType) ([]string, error) {
	return m.getFullPubkeyNotExportedPubKey(accountKeyTableName[accountType])
}

//getAllNotExportedPubKey IsExprotedPubKeyがfalseのレコードをすべて返す
func (m *DB) getAllNotExportedPubKey(tbl string, accountType enum.AccountType) ([]AccountKeyTable, error) {
	sql := "SELECT * FROM %s WHERE is_exported_pub_key=false;"
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var accountKeyTable []AccountKeyTable
	err := m.RDB.Select(&accountKeyTable, sql)
	if err != nil {
		return nil, err
	}

	return accountKeyTable, nil
}

//GetAllNotExportedPubKey IsExprotedPubKeyがfalseのレコードをすべて返す
func (m *DB) GetAllNotExportedPubKey(accountType enum.AccountType) ([]AccountKeyTable, error) {
	return m.getAllNotExportedPubKey(accountKeyTableName[accountType], accountType)
}

//getOneByMaxID idが最大の1レコードを返す
func (m *DB) getOneByMaxID(tbl string, accountType enum.AccountType) (*AccountKeyTable, error) {
	sql := "SELECT * FROM %s ORDER BY ID DESC LIMIT 1;"
	sql = fmt.Sprintf(sql, tbl)
	logger.Debugf("sql: %s", sql)

	var accountKeyTable AccountKeyTable
	err := m.RDB.Get(&accountKeyTable, sql)
	if err != nil {
		return nil, err
	}

	return &accountKeyTable, nil
}

//GetOneByMaxID idが最大の1レコードを返す
func (m *DB) GetOneByMaxID(accountType enum.AccountType) (*AccountKeyTable, error) {
	return m.getOneByMaxID(accountKeyTableName[accountType], accountType)
}
