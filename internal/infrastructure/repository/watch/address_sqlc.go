package watch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	sqlc "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
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
func (r *AddressRepositorySqlc) GetAll(accountType domainAccount.AccountType) ([]*sqlc.Address, error) {
	ctx := context.Background()

	addresses, err := r.queries.GetAllAddresses(ctx, sqlc.GetAllAddressesParams{
		Coin:    sqlc.AddressCoin(r.coinTypeCode.String()),
		Account: sqlc.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllAddresses(): %w", err)
	}

	result := make([]*sqlc.Address, len(addresses))
	for i := range addresses {
		result[i] = &addresses[i]
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
func (r *AddressRepositorySqlc) GetOneUnAllocated(accountType domainAccount.AccountType) (*sqlc.Address, error) {
	ctx := context.Background()

	addr, err := r.queries.GetOneUnallocatedAddress(ctx, sqlc.GetOneUnallocatedAddressParams{
		Coin:    sqlc.AddressCoin(r.coinTypeCode.String()),
		Account: sqlc.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneUnallocatedAddress(): %w", err)
	}

	return &addr, nil
}

// InsertBulk inserts multiple records
func (r *AddressRepositorySqlc) InsertBulk(ctx context.Context, items []*sqlc.Address) error {
	for _, item := range items {
		_, err := r.queries.InsertAddress(ctx, sqlc.InsertAddressParams{
			Coin:          item.Coin,
			Account:       item.Account,
			WalletAddress: item.WalletAddress,
			IsAllocated:   item.IsAllocated,
			UpdatedAt:     item.UpdatedAt,
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
