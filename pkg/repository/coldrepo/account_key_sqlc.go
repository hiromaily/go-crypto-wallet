package coldrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// AccountKeyRepositorySqlc is repository for account_key table using sqlc
type AccountKeyRepositorySqlc struct {
	queries      *sqlcgen.Queries
	dbConn       *sql.DB
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewAccountKeyRepositorySqlc returns AccountKeyRepositorySqlc object
func NewAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger,
) *AccountKeyRepositorySqlc {
	return &AccountKeyRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		dbConn:       dbConn,
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetMaxIndex returns max idx
func (r *AccountKeyRepositorySqlc) GetMaxIndex(accountType account.AccountType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.GetMaxAccountKeyIndex(ctx, sqlcgen.GetMaxAccountKeyIndexParams{
		Coin:    sqlcgen.AccountKeyCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AccountKeyAccount(accountType.String()),
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call GetMaxAccountKeyIndex()")
	}

	// Type assert interface{} to int64
	if maxIdx, ok := result.(int64); ok {
		return maxIdx, nil
	}

	return 0, nil
}

// GetOneMaxID returns one record by max id
func (r *AccountKeyRepositorySqlc) GetOneMaxID(accountType account.AccountType) (*models.AccountKey, error) {
	ctx := context.Background()

	accountKey, err := r.queries.GetOneAccountKeyByMaxID(ctx, sqlcgen.GetOneAccountKeyByMaxIDParams{
		Coin:    sqlcgen.AccountKeyCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AccountKeyAccount(accountType.String()),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetOneAccountKeyByMaxID()")
	}

	return convertSqlcAccountKeyToModel(&accountKey), nil
}

// GetAllAddrStatus returns all AccountKey by addr_status
func (r *AccountKeyRepositorySqlc) GetAllAddrStatus(
	accountType account.AccountType, addrStatus address.AddrStatus,
) ([]*models.AccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetAccountKeysByAddrStatus(ctx, sqlcgen.GetAccountKeysByAddrStatusParams{
		Coin:       sqlcgen.AccountKeyCoin(r.coinTypeCode.String()),
		Account:    sqlcgen.AccountKeyAccount(accountType.String()),
		AddrStatus: addrStatus.Int8(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetAccountKeysByAddrStatus()")
	}

	result := make([]*models.AccountKey, len(accountKeys))
	for i, accountKey := range accountKeys {
		result[i] = convertSqlcAccountKeyToModel(&accountKey)
	}

	return result, nil
}

// GetAllMultiAddr returns all AccountKey by multisig_address
func (r *AccountKeyRepositorySqlc) GetAllMultiAddr(
	accountType account.AccountType, addrs []string,
) ([]*models.AccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetAccountKeysByMultisigAddresses(ctx, sqlcgen.GetAccountKeysByMultisigAddressesParams{
		Coin:    sqlcgen.AccountKeyCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AccountKeyAccount(accountType.String()),
		Addrs:   addrs,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetAccountKeysByMultisigAddresses()")
	}

	result := make([]*models.AccountKey, len(accountKeys))
	for i, accountKey := range accountKeys {
		result[i] = convertSqlcAccountKeyToModel(&accountKey)
	}

	return result, nil
}

// InsertBulk inserts multiple records
func (r *AccountKeyRepositorySqlc) InsertBulk(items []*models.AccountKey) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertAccountKey(ctx, sqlcgen.InsertAccountKeyParams{
			Coin:               sqlcgen.AccountKeyCoin(item.Coin),
			Account:            sqlcgen.AccountKeyAccount(item.Account),
			P2pkhAddress:       item.P2PKHAddress,
			P2shSegwitAddress:  item.P2SHSegwitAddress,
			Bech32Address:      item.Bech32Address,
			FullPublicKey:      item.FullPublicKey,
			MultisigAddress:    item.MultisigAddress,
			RedeemScript:       item.RedeemScript,
			WalletImportFormat: item.WalletImportFormat,
			Idx:                item.Idx,
			AddrStatus:         item.AddrStatus,
		})
		if err != nil {
			return errors.Wrap(err, "failed to call InsertAccountKey()")
		}
	}

	return nil
}

// UpdateAddr updates address by P2SHSegWitAddr
func (r *AccountKeyRepositorySqlc) UpdateAddr(accountType account.AccountType, addr, keyAddress string) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateAccountKeyAddress(ctx, sqlcgen.UpdateAccountKeyAddressParams{
		P2pkhAddress:      addr,
		UpdatedAt:         sql.NullTime{Time: time.Now(), Valid: true},
		Coin:              sqlcgen.AccountKeyCoin(r.coinTypeCode.String()),
		Account:           sqlcgen.AccountKeyAccount(accountType.String()),
		P2shSegwitAddress: keyAddress,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateAccountKeyAddress()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// UpdateAddrStatus updates addr_status
func (r *AccountKeyRepositorySqlc) UpdateAddrStatus(
	accountType account.AccountType, addrStatus address.AddrStatus, strWIFs []string,
) (int64, error) {
	ctx := context.Background()
	var totalAffected int64

	// sqlc doesn't support IN clauses with variable arguments, so update one at a time
	for _, wif := range strWIFs {
		result, err := r.queries.UpdateAccountKeyAddrStatus(ctx, sqlcgen.UpdateAccountKeyAddrStatusParams{
			AddrStatus:         addrStatus.Int8(),
			UpdatedAt:          sql.NullTime{Time: time.Now(), Valid: true},
			Coin:               sqlcgen.AccountKeyCoin(r.coinTypeCode.String()),
			Account:            sqlcgen.AccountKeyAccount(accountType.String()),
			WalletImportFormat: wif,
		})
		if err != nil {
			return 0, errors.Wrap(err, "failed to call UpdateAccountKeyAddrStatus()")
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return 0, errors.Wrap(err, "failed to get RowsAffected()")
		}
		totalAffected += affected
	}

	return totalAffected, nil
}

// UpdateMultisigAddr updates multisig_address
func (r *AccountKeyRepositorySqlc) UpdateMultisigAddr(
	accountType account.AccountType, item *models.AccountKey,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateAccountKeyMultisigAddr(ctx, sqlcgen.UpdateAccountKeyMultisigAddrParams{
		MultisigAddress: item.MultisigAddress,
		RedeemScript:    item.RedeemScript,
		AddrStatus:      item.AddrStatus,
		UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		Coin:            sqlcgen.AccountKeyCoin(r.coinTypeCode.String()),
		Account:         sqlcgen.AccountKeyAccount(accountType.String()),
		FullPublicKey:   item.FullPublicKey,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateAccountKeyMultisigAddr()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// UpdateMultisigAddrs updates all multisig_address with transaction
func (r *AccountKeyRepositorySqlc) UpdateMultisigAddrs(
	accountType account.AccountType, items []*models.AccountKey,
) (int64, error) {
	ctx := context.Background()

	// transaction
	dtx, err := r.dbConn.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "failed to call db.Begin()")
	}
	defer func() {
		if err != nil {
			dtx.Rollback()
		} else {
			dtx.Commit()
		}
	}()

	qtx := r.queries.WithTx(dtx)
	var totalAffected int64

	for _, item := range items {
		result, err := qtx.UpdateAccountKeyMultisigAddr(ctx, sqlcgen.UpdateAccountKeyMultisigAddrParams{
			MultisigAddress: item.MultisigAddress,
			RedeemScript:    item.RedeemScript,
			AddrStatus:      item.AddrStatus,
			UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
			Coin:            sqlcgen.AccountKeyCoin(r.coinTypeCode.String()),
			Account:         sqlcgen.AccountKeyAccount(accountType.String()),
			FullPublicKey:   item.FullPublicKey,
		})
		if err != nil {
			return 0, errors.Wrap(err, "failed to call UpdateAccountKeyMultisigAddr()")
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return 0, errors.Wrap(err, "failed to get RowsAffected()")
		}
		totalAffected += affected
	}

	return totalAffected, nil
}

// Helper functions

func convertSqlcAccountKeyToModel(accountKey *sqlcgen.AccountKey) *models.AccountKey {
	return &models.AccountKey{
		ID:                accountKey.ID,
		Coin:              string(accountKey.Coin),
		Account:           string(accountKey.Account),
		P2PKHAddress:      accountKey.P2pkhAddress,
		P2SHSegwitAddress: accountKey.P2shSegwitAddress,
		Bech32Address:     accountKey.Bech32Address,
		FullPublicKey:     accountKey.FullPublicKey,
		MultisigAddress:   accountKey.MultisigAddress,
		RedeemScript:      accountKey.RedeemScript,
		WalletImportFormat: accountKey.WalletImportFormat,
		Idx:               accountKey.Idx,
		AddrStatus:        accountKey.AddrStatus,
		UpdatedAt:         convertSQLNullTimeToNullTime(accountKey.UpdatedAt),
	}
}

func convertSQLNullTimeToNullTime(t sql.NullTime) null.Time {
	if !t.Valid {
		return null.Time{}
	}
	return null.TimeFrom(t.Time)
}

func convertNullTimeToSQLNullTime(t null.Time) sql.NullTime {
	if !t.Valid {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t.Time, Valid: true}
}
