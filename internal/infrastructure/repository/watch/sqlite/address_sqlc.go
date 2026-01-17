package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
)

// AddressRepositorySqlc is repository for address table using sqlc (SQLite)
type AddressRepositorySqlc struct {
	db           *sql.DB
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAddressRepositorySqlc returns AddressRepositorySqlc object
func NewAddressRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *AddressRepositorySqlc {
	return &AddressRepositorySqlc{
		db:           dbConn,
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
		IsAllocated:   sqlcAddr.IsAllocated != 0, // SQLite uses INTEGER for boolean
	}

	// Parse SQLite TEXT timestamp to time.Time
	if sqlcAddr.UpdatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", sqlcAddr.UpdatedAt.String)
		if err == nil {
			addr.UpdatedAt = &t
		}
	}

	return addr, nil
}

// convertFromAddress converts domain.Address entity to sqlcgen.Address
func convertFromAddress(addr *domainAddress.Address) *sqlcgen.Address {
	sqlcAddr := &sqlcgen.Address{
		ID:            addr.ID,
		Coin:          addr.CoinTypeCode.String(),
		Account:       addr.AccountType.String(),
		WalletAddress: addr.WalletAddress,
		IsAllocated:   boolToInt64(addr.IsAllocated),
	}

	// Convert time.Time to SQLite TEXT format
	if addr.UpdatedAt != nil {
		sqlcAddr.UpdatedAt = sql.NullString{
			String: addr.UpdatedAt.Format("2006-01-02 15:04:05"),
			Valid:  true,
		}
	}

	return sqlcAddr
}

// GetAll returns all records by account
func (r *AddressRepositorySqlc) GetAll(accountType domainAccount.AccountType) ([]*domainAddress.Address, error) {
	ctx := context.Background()

	addresses, err := r.queries.GetAllAddresses(ctx, sqlcgen.GetAllAddressesParams{
		Coin:    r.coinTypeCode.String(),
		Account: accountType.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllAddresses(): %w", err)
	}

	result := make([]*domainAddress.Address, 0, len(addresses))
	for i := range addresses {
		addr, err := convertToAddress(&addresses[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert address: %w", err)
		}
		result = append(result, addr)
	}

	return result, nil
}

// GetAllAddress returns all wallet addresses by account
func (r *AddressRepositorySqlc) GetAllAddress(accountType domainAccount.AccountType) ([]string, error) {
	ctx := context.Background()

	addresses, err := r.queries.GetAllAddressStrings(ctx, sqlcgen.GetAllAddressStringsParams{
		Coin:    r.coinTypeCode.String(),
		Account: accountType.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllAddressStrings(): %w", err)
	}

	return addresses, nil
}

// GetOneUnAllocated returns one unallocated address
func (r *AddressRepositorySqlc) GetOneUnAllocated(
	accountType domainAccount.AccountType,
) (*domainAddress.Address, error) {
	ctx := context.Background()

	addr, err := r.queries.GetOneUnallocatedAddress(ctx, sqlcgen.GetOneUnallocatedAddressParams{
		Coin:    r.coinTypeCode.String(),
		Account: accountType.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneUnallocatedAddress(): %w", err)
	}

	return convertToAddress(&addr)
}

// InsertBulk inserts multiple addresses
func (r *AddressRepositorySqlc) InsertBulk(ctx context.Context, items []*domainAddress.Address) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // Rollback is safe to call even after Commit
	}()

	qtx := r.queries.WithTx(tx)

	for _, item := range items {
		sqlcAddr := convertFromAddress(item)
		_, err := qtx.InsertAddress(ctx, sqlcgen.InsertAddressParams{
			Coin:          sqlcAddr.Coin,
			Account:       sqlcAddr.Account,
			WalletAddress: sqlcAddr.WalletAddress,
			IsAllocated:   sqlcAddr.IsAllocated,
			UpdatedAt:     sqlcAddr.UpdatedAt,
		})
		if err != nil {
			return fmt.Errorf("failed to insert address: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateIsAllocated updates is_allocated flag
func (r *AddressRepositorySqlc) UpdateIsAllocated(isAllocated bool, address string) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateAddressIsAllocated(ctx, sqlcgen.UpdateAddressIsAllocatedParams{
		IsAllocated:   boolToInt64(isAllocated),
		UpdatedAt:     sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
		Coin:          r.coinTypeCode.String(),
		WalletAddress: address,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to update is_allocated: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
