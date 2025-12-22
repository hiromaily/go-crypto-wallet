package watchrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
)

// AddressRepositorySqlc is repository for address table using sqlc
type AddressRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewAddressRepositorySqlc returns AddressRepositorySqlc object
func NewAddressRepositorySqlc(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger) *AddressRepositorySqlc {
	return &AddressRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetAll returns all records by account
func (r *AddressRepositorySqlc) GetAll(accountType account.AccountType) ([]*models.Address, error) {
	ctx := context.Background()

	addresses, err := r.queries.GetAllAddresses(ctx, sqlcgen.GetAllAddressesParams{
		Coin:    sqlcgen.AddressCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetAllAddresses()")
	}

	// Convert sqlc types to sqlboiler types for backward compatibility
	result := make([]*models.Address, len(addresses))
	for i, addr := range addresses {
		result[i] = convertSqlcAddressToModel(&addr)
	}

	return result, nil
}

// GetAllAddress returns all addresses by account
func (r *AddressRepositorySqlc) GetAllAddress(accountType account.AccountType) ([]string, error) {
	ctx := context.Background()

	addresses, err := r.queries.GetAllAddressStrings(ctx, sqlcgen.GetAllAddressStringsParams{
		Coin:    sqlcgen.AddressCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetAllAddressStrings()")
	}

	return addresses, nil
}

// GetOneUnAllocated returns one records by is_allocated=false
func (r *AddressRepositorySqlc) GetOneUnAllocated(accountType account.AccountType) (*models.Address, error) {
	ctx := context.Background()

	addr, err := r.queries.GetOneUnallocatedAddress(ctx, sqlcgen.GetOneUnallocatedAddressParams{
		Coin:    sqlcgen.AddressCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetOneUnallocatedAddress()")
	}

	return convertSqlcAddressToModel(&addr), nil
}

// InsertBulk inserts multiple records
func (r *AddressRepositorySqlc) InsertBulk(items []*models.Address) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertAddress(ctx, sqlcgen.InsertAddressParams{
			Coin:          sqlcgen.AddressCoin(item.Coin),
			Account:       sqlcgen.AddressAccount(item.Account),
			WalletAddress: item.WalletAddress,
			IsAllocated:   item.IsAllocated,
			UpdatedAt:     convertNullTimeToSqlNullTime(item.UpdatedAt),
		})
		if err != nil {
			return errors.Wrap(err, "failed to call InsertAddress()")
		}
	}

	return nil
}

// UpdateIsAllocated updates is_allocated
func (r *AddressRepositorySqlc) UpdateIsAllocated(isAllocated bool, address string) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateAddressIsAllocated(ctx, sqlcgen.UpdateAddressIsAllocatedParams{
		IsAllocated:   isAllocated,
		UpdatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
		Coin:          sqlcgen.AddressCoin(r.coinTypeCode.String()),
		WalletAddress: address,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateAddressIsAllocated()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// Helper functions for type conversion

func convertSqlcAddressToModel(addr *sqlcgen.Address) *models.Address {
	return &models.Address{
		ID:            addr.ID,
		Coin:          string(addr.Coin),
		Account:       string(addr.Account),
		WalletAddress: addr.WalletAddress,
		IsAllocated:   addr.IsAllocated,
		UpdatedAt:     convertSqlNullTimeToNullTime(addr.UpdatedAt),
	}
}

func convertSqlNullTimeToNullTime(t sql.NullTime) null.Time {
	if !t.Valid {
		return null.Time{}
	}
	return null.TimeFrom(t.Time)
}

func convertNullTimeToSqlNullTime(t null.Time) sql.NullTime {
	if !t.Valid {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t.Time, Valid: true}
}
