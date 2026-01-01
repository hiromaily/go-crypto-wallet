package watch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// AddressRepositorySqlc is repository for address table using sqlc
type AddressRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAddressRepositorySqlc returns AddressRepositorySqlc object
func NewAddressRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *AddressRepositorySqlc {
	return &AddressRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetAll returns all records by account
func (r *AddressRepositorySqlc) GetAll(accountType domainAccount.AccountType) ([]*sqlcgen.Address, error) {
	ctx := context.Background()

	addresses, err := r.queries.GetAllAddresses(ctx, sqlcgen.GetAllAddressesParams{
		Coin:    sqlcgen.AddressCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllAddresses(): %w", err)
	}

	result := make([]*sqlcgen.Address, len(addresses))
	for i := range addresses {
		result[i] = &addresses[i]
	}

	return result, nil
}

// GetAllAddress returns all addresses by account
func (r *AddressRepositorySqlc) GetAllAddress(accountType domainAccount.AccountType) ([]string, error) {
	ctx := context.Background()

	addresses, err := r.queries.GetAllAddressStrings(ctx, sqlcgen.GetAllAddressStringsParams{
		Coin:    sqlcgen.AddressCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllAddressStrings(): %w", err)
	}

	return addresses, nil
}

// GetOneUnAllocated returns one records by is_allocated=false
func (r *AddressRepositorySqlc) GetOneUnAllocated(accountType domainAccount.AccountType) (*sqlcgen.Address, error) {
	ctx := context.Background()

	addr, err := r.queries.GetOneUnallocatedAddress(ctx, sqlcgen.GetOneUnallocatedAddressParams{
		Coin:    sqlcgen.AddressCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneUnallocatedAddress(): %w", err)
	}

	return &addr, nil
}

// InsertBulk inserts multiple records
func (r *AddressRepositorySqlc) InsertBulk(ctx context.Context, items []*sqlcgen.Address) error {
	for _, item := range items {
		_, err := r.queries.InsertAddress(ctx, sqlcgen.InsertAddressParams{
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

	result, err := r.queries.UpdateAddressIsAllocated(ctx, sqlcgen.UpdateAddressIsAllocatedParams{
		IsAllocated:   isAllocated,
		UpdatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
		Coin:          sqlcgen.AddressCoin(r.coinTypeCode.String()),
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
