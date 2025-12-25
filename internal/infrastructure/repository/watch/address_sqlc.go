package watch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/guregu/null/v6"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	models "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
)

// AddressRepositorySqlc is repository for address table using sqlc
type AddressRepositorySqlc struct {
	queries      *sqlc.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAddressRepositorySqlc returns AddressRepositorySqlc object
func NewAddressRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *AddressRepositorySqlc {
	return &AddressRepositorySqlc{
		queries:      sqlc.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetAll returns all records by account
func (r *AddressRepositorySqlc) GetAll(accountType domainAccount.AccountType) ([]*models.Address, error) {
	ctx := context.Background()

	addresses, err := r.queries.GetAllAddresses(ctx, sqlc.GetAllAddressesParams{
		Coin:    sqlc.AddressCoin(r.coinTypeCode.String()),
		Account: sqlc.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllAddresses(): %w", err)
	}

	// Convert sqlc types to sqlboiler types for backward compatibility
	result := make([]*models.Address, len(addresses))
	for i, addr := range addresses {
		result[i] = convertSqlcAddressToModel(&addr)
	}

	return result, nil
}

// GetAllAddress returns all addresses by account
func (r *AddressRepositorySqlc) GetAllAddress(accountType domainAccount.AccountType) ([]string, error) {
	ctx := context.Background()

	addresses, err := r.queries.GetAllAddressStrings(ctx, sqlc.GetAllAddressStringsParams{
		Coin:    sqlc.AddressCoin(r.coinTypeCode.String()),
		Account: sqlc.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllAddressStrings(): %w", err)
	}

	return addresses, nil
}

// GetOneUnAllocated returns one records by is_allocated=false
func (r *AddressRepositorySqlc) GetOneUnAllocated(accountType domainAccount.AccountType) (*models.Address, error) {
	ctx := context.Background()

	addr, err := r.queries.GetOneUnallocatedAddress(ctx, sqlc.GetOneUnallocatedAddressParams{
		Coin:    sqlc.AddressCoin(r.coinTypeCode.String()),
		Account: sqlc.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneUnallocatedAddress(): %w", err)
	}

	return convertSqlcAddressToModel(&addr), nil
}

// InsertBulk inserts multiple records
func (r *AddressRepositorySqlc) InsertBulk(items []*models.Address) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertAddress(ctx, sqlc.InsertAddressParams{
			Coin:          sqlc.AddressCoin(item.Coin),
			Account:       sqlc.AddressAccount(item.Account),
			WalletAddress: item.WalletAddress,
			IsAllocated:   item.IsAllocated,
			UpdatedAt:     convertNullTimeToSQLNullTime(item.UpdatedAt),
		})
		if err != nil {
			return fmt.Errorf("failed to call InsertAddress(): %w", err)
		}
	}

	return nil
}

// UpdateIsAllocated updates is_allocated
func (r *AddressRepositorySqlc) UpdateIsAllocated(isAllocated bool, address string) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateAddressIsAllocated(ctx, sqlc.UpdateAddressIsAllocatedParams{
		IsAllocated:   isAllocated,
		UpdatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
		Coin:          sqlc.AddressCoin(r.coinTypeCode.String()),
		WalletAddress: address,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateAddressIsAllocated(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// Helper functions for type conversion

func convertSqlcAddressToModel(addr *sqlc.Address) *models.Address {
	return &models.Address{
		ID:            addr.ID,
		Coin:          string(addr.Coin),
		Account:       string(addr.Account),
		WalletAddress: addr.WalletAddress,
		IsAllocated:   addr.IsAllocated,
		UpdatedAt:     convertSQLNullTimeToNullTime(addr.UpdatedAt),
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
