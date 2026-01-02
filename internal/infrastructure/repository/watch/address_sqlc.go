package watch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
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

// convertToAddress converts sqlcgen.Address to domain.Address entity
func convertToAddress(sqlcAddr *sqlcgen.Address) (*domainAddress.Address, error) {
	addr := &domainAddress.Address{
		ID:            sqlcAddr.ID,
		CoinTypeCode:  domainCoin.CoinTypeCode(sqlcAddr.Coin),
		AccountType:   domainAccount.AccountType(sqlcAddr.Account),
		WalletAddress: sqlcAddr.WalletAddress,
		IsAllocated:   sqlcAddr.IsAllocated,
	}

	if sqlcAddr.UpdatedAt.Valid {
		addr.UpdatedAt = &sqlcAddr.UpdatedAt.Time
	}

	return addr, nil
}

// convertFromAddress converts domain.Address entity to sqlcgen.Address
func convertFromAddress(addr *domainAddress.Address) *sqlcgen.Address {
	sqlcAddr := &sqlcgen.Address{
		ID:            addr.ID,
		Coin:          sqlcgen.AddressCoin(addr.CoinTypeCode.String()),
		Account:       sqlcgen.AddressAccount(addr.AccountType.String()),
		WalletAddress: addr.WalletAddress,
		IsAllocated:   addr.IsAllocated,
	}

	if addr.UpdatedAt != nil {
		sqlcAddr.UpdatedAt = sql.NullTime{Time: *addr.UpdatedAt, Valid: true}
	}

	return sqlcAddr
}

// GetAll returns all records by account
func (r *AddressRepositorySqlc) GetAll(accountType domainAccount.AccountType) ([]*domainAddress.Address, error) {
	ctx := context.Background()

	addresses, err := r.queries.GetAllAddresses(ctx, sqlcgen.GetAllAddressesParams{
		Coin:    sqlcgen.AddressCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllAddresses(): %w", err)
	}

	result := make([]*domainAddress.Address, 0, len(addresses))
	for i := range addresses {
		addr, err := convertToAddress(&addresses[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert address at index %d: %w", i, err)
		}
		result = append(result, addr)
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
func (r *AddressRepositorySqlc) GetOneUnAllocated(accountType domainAccount.AccountType) (*domainAddress.Address, error) {
	ctx := context.Background()

	addr, err := r.queries.GetOneUnallocatedAddress(ctx, sqlcgen.GetOneUnallocatedAddressParams{
		Coin:    sqlcgen.AddressCoin(r.coinTypeCode.String()),
		Account: sqlcgen.AddressAccount(accountType.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneUnallocatedAddress(): %w", err)
	}

	return convertToAddress(&addr)
}

// InsertBulk inserts multiple records
func (r *AddressRepositorySqlc) InsertBulk(ctx context.Context, items []*domainAddress.Address) error {
	for _, item := range items {
		sqlcAddr := convertFromAddress(item)
		_, err := r.queries.InsertAddress(ctx, sqlcgen.InsertAddressParams{
			Coin:          sqlcAddr.Coin,
			Account:       sqlcAddr.Account,
			WalletAddress: sqlcAddr.WalletAddress,
			IsAllocated:   sqlcAddr.IsAllocated,
			UpdatedAt:     sqlcAddr.UpdatedAt,
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
