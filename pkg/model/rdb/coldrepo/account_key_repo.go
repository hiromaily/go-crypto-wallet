package coldrepo

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/keystatus"
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
	Idx                   uint32     `db:"idx"`
	KeyStatus             uint8      `db:"key_status"`
	UpdatedAt             *time.Time `db:"updated_at"`
}

var accountKeyTableName = map[account.AccountType]string{
	account.AccountTypeClient:        "account_key_client",
	account.AccountTypeReceipt:       "account_key_receipt",
	account.AccountTypePayment:       "account_key_payment",
	account.AccountTypeQuoine:        "account_key_quoine",
	account.AccountTypeFee:           "account_key_fee",
	account.AccountTypeStored:        "account_key_stored",
	account.AccountTypeAuthorization: "account_key_authorization",
}

//getMaxIndexOnAccountKeyTable indexの最大値を返す
func (r *ColdRepository) getMaxIndexOnAccountKeyTable(tbl string) (int64, error) {
	sql := "SELECT MAX(idx) from %s;"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	var idx int64
	err := r.db.Get(&idx, sql)

	return idx, err
}

//GetMaxIndexOnAccountKeyTable indexの最大値を返す
func (r *ColdRepository) GetMaxIndexOnAccountKeyTable(accountType account.AccountType) (int64, error) {
	return r.getMaxIndexOnAccountKeyTable(accountKeyTableName[accountType])
}

//getOneByMaxIDOnAccountKeyTable idが最大の1レコードを返す
func (r *ColdRepository) getOneByMaxIDOnAccountKeyTable(tbl string, accountType account.AccountType) (*AccountKeyTable, error) {
	sql := "SELECT * FROM %s ORDER BY ID DESC LIMIT 1;"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	var accountKeyTable AccountKeyTable
	err := r.db.Get(&accountKeyTable, sql)
	if err != nil {
		return nil, err
	}

	return &accountKeyTable, nil
}

//GetOneByMaxIDOnAccountKeyTable idが最大の1レコードを返す
func (r *ColdRepository) GetOneByMaxIDOnAccountKeyTable(accountType account.AccountType) (*AccountKeyTable, error) {
	return r.getOneByMaxIDOnAccountKeyTable(accountKeyTableName[accountType], accountType)
}

// getAllAccountKeyByKeyStatus 指定したkeyStatusのレコードをすべて返す
func (r *ColdRepository) getAllAccountKeyByKeyStatus(tbl string, keyStatus keystatus.KeyStatus) ([]AccountKeyTable, error) {
	//sql := "SELECT * FROM %s WHERE is_imported_priv_key=false;"
	sql := "SELECT * FROM %s WHERE key_status=?;"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	var accountKeyTable []AccountKeyTable
	err := r.db.Select(&accountKeyTable, sql, keystatus.KeyStatusValue[keyStatus])
	if err != nil {
		return nil, err
	}

	return accountKeyTable, nil
}

// GetAllAccountKeyByKeyStatus 指定したkeyStatusのレコードをすべて返す
func (r *ColdRepository) GetAllAccountKeyByKeyStatus(accountType account.AccountType, keyStatus keystatus.KeyStatus) ([]AccountKeyTable, error) {
	return r.getAllAccountKeyByKeyStatus(accountKeyTableName[accountType], keyStatus)
}

func (r *ColdRepository) getAllAccountKeyByMultiAddrs(tbl string, addrs []string) ([]AccountKeyTable, error) {
	sql := "SELECT * FROM %s WHERE wallet_multisig_address IN (?);"
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	//In対応
	query, args, err := sqlx.In(sql, addrs)
	if err != nil {
		return nil, errors.Errorf("sqlx.In() error: %v", err)
	}
	query = r.db.Rebind(query)
	//logger.Debugf("sql: %s", query)

	var accountKeyTable []AccountKeyTable
	err = r.db.Select(&accountKeyTable, query, args...)
	if err != nil {
		return nil, err
	}

	return accountKeyTable, nil
}

// GetAllAccountKeyByMultiAddrs WIPをmultiAddressから取得する
func (r *ColdRepository) GetAllAccountKeyByMultiAddrs(accountType account.AccountType, addrs []string) ([]AccountKeyTable, error) {
	return r.getAllAccountKeyByMultiAddrs(accountKeyTableName[accountType], addrs)
}

// insertAccountKeyClient account_key_table(client, payment, receipt...)テーブルにレコードを作成する
//TODO:BulkInsertがやりたい
func (r *ColdRepository) insertAccountKeyTable(tbl string, accountKeyTables []AccountKeyTable, tx *sqlx.Tx, isCommit bool) error {

	sql := `
INSERT INTO %s (wallet_address, p2sh_segwit_address, full_public_key, wallet_multisig_address, redeem_script, wallet_import_format, account, idx) 
VALUES (:wallet_address, :p2sh_segwit_address,:full_public_key, :wallet_multisig_address, :redeem_script, :wallet_import_format, :account, :idx)
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = r.db.MustBegin()
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
func (r *ColdRepository) InsertAccountKeyTable(accountType account.AccountType, accountKeyTables []AccountKeyTable, tx *sqlx.Tx, isCommit bool) error {
	return r.insertAccountKeyTable(accountKeyTableName[accountType], accountKeyTables, tx, isCommit)
}

// updateKeyStatus key_statusを更新する
func (r *ColdRepository) updateKeyStatusByWIF(tbl string, keyStatus keystatus.KeyStatus, strWIF string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	sql := `
UPDATE %s SET key_status=? WHERE wallet_import_format=?
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = r.db.MustBegin()
	}

	res, err := tx.Exec(sql, keystatus.KeyStatusValue[keyStatus], strWIF)
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

// UpdateKeyStatusByWIF key_statusを更新する
func (r *ColdRepository) UpdateKeyStatusByWIF(accountType account.AccountType, keyStatus keystatus.KeyStatus, strWIF string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return r.updateKeyStatusByWIF(accountKeyTableName[accountType], keyStatus, strWIF, tx, isCommit)
}

// updateKeyStatusByWIFs key_statusを更新する
func (r *ColdRepository) updateKeyStatusByWIFs(tbl string, keyStatus keystatus.KeyStatus, wifs []string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	var sql string
	sql = "UPDATE %s SET key_status=%d WHERE wallet_import_format IN (?);"
	sql = fmt.Sprintf(sql, tbl, keystatus.KeyStatusValue[keyStatus])

	//In対応
	query, args, err := sqlx.In(sql, wifs)
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

// UpdateKeyStatusByWIFs key_statusを更新する
func (r *ColdRepository) UpdateKeyStatusByWIFs(accountType account.AccountType, keyStatus keystatus.KeyStatus, wifs []string, tx *sqlx.Tx, isCommit bool) (int64, error) {
	return r.updateKeyStatusByWIFs(accountKeyTableName[accountType], keyStatus, wifs, tx, isCommit)
}

// updateMultisigAddrOnAccountKeyTableByFullPubKey wallet_multisig_addressを更新する
func (r *ColdRepository) updateMultisigAddrOnAccountKeyTableByFullPubKey(tbl string, accountKeyTable []AccountKeyTable, tx *sqlx.Tx, isCommit bool) error {
	sql := `
UPDATE %s SET wallet_multisig_address=:wallet_multisig_address, redeem_script=:redeem_script, key_status=:key_status, updated_at=:updated_at 
WHERE full_public_key=:full_public_key
`
	sql = fmt.Sprintf(sql, tbl)
	//logger.Debugf("sql: %s", sql)

	if tx == nil {
		tx = r.db.MustBegin()
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

// UpdateMultisigAddrOnAccountKeyTableByFullPubKey wallet_multisig_addressを更新する
func (r *ColdRepository) UpdateMultisigAddrOnAccountKeyTableByFullPubKey(accountType account.AccountType, accountKeyTable []AccountKeyTable, tx *sqlx.Tx, isCommit bool) error {
	return r.updateMultisigAddrOnAccountKeyTableByFullPubKey(accountKeyTableName[accountType], accountKeyTable, tx, isCommit)
}

// GetRedeedScriptByAddress 与えられたmultiSigアドレスから、RedeemScriptを取得する
func GetRedeedScriptByAddress(accountKeys []AccountKeyTable, addr string) string {
	for _, val := range accountKeys {
		if val.WalletMultisigAddress == addr {
			return val.RedeemScript
		}
	}
	return ""
}
