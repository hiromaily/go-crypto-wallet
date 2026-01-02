package cold

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// ETHAccountKeyRepositorySqlc is repository for eth_account_key table using sqlc
type ETHAccountKeyRepositorySqlc struct {
	queries *sqlcgen.Queries
	dbConn  *sql.DB
}

// NewETHAccountKeyRepositorySqlc returns ETHAccountKeyRepositorySqlc object
func NewETHAccountKeyRepositorySqlc(dbConn *sql.DB) *ETHAccountKeyRepositorySqlc {
	return &ETHAccountKeyRepositorySqlc{
		queries: sqlcgen.New(dbConn),
		dbConn:  dbConn,
	}
}

// GetMaxIndex returns max idx
func (r *ETHAccountKeyRepositorySqlc) GetMaxIndex(accountType domainAccount.AccountType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.GetMaxEthAccountKeyIndex(ctx, sqlcgen.EthAccountKeyAccount(accountType.String()))
	if err != nil {
		return 0, fmt.Errorf("failed to call GetMaxEthAccountKeyIndex(): %w", err)
	}

	// Type assert interface{} to int64
	if maxIdx, ok := result.(int64); ok {
		return maxIdx, nil
	}

	return 0, nil
}

// GetOneMaxID returns one record by max id
func (r *ETHAccountKeyRepositorySqlc) GetOneMaxID(accountType domainAccount.AccountType,
) (*sqlcgen.EthAccountKey, error) {
	ctx := context.Background()

	accountKey, err := r.queries.GetOneEthAccountKeyByMaxID(ctx, sqlcgen.EthAccountKeyAccount(accountType.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneEthAccountKeyByMaxID(): %w", err)
	}

	return &accountKey, nil
}

// GetAllAddrStatus returns all EthAccountKey by addr_status
func (r *ETHAccountKeyRepositorySqlc) GetAllAddrStatus(
	accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
) ([]*sqlcgen.EthAccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetEthAccountKeysByAddrStatus(ctx, sqlcgen.GetEthAccountKeysByAddrStatusParams{
		Account:    sqlcgen.EthAccountKeyAccount(accountType.String()),
		AddrStatus: addrStatus.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetEthAccountKeysByAddrStatus(): %w", err)
	}

	result := make([]*sqlcgen.EthAccountKey, len(accountKeys))
	for i := range accountKeys {
		result[i] = &accountKeys[i]
	}

	return result, nil
}

// GetByAddress returns EthAccountKey by address
func (r *ETHAccountKeyRepositorySqlc) GetByAddress(addr string) (*sqlcgen.EthAccountKey, error) {
	ctx := context.Background()

	accountKey, err := r.queries.GetEthAccountKeyByAddress(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetEthAccountKeyByAddress(): %w", err)
	}

	return &accountKey, nil
}

// InsertBulk inserts multiple records
func (r *ETHAccountKeyRepositorySqlc) InsertBulk(items []*sqlcgen.EthAccountKey) error {
	ctx := context.Background()

	tx, err := r.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := r.queries.WithTx(tx)

	for _, item := range items {
		_, err := qtx.InsertEthAccountKey(ctx, sqlcgen.InsertEthAccountKeyParams{
			Account:       item.Account,
			Address:       item.Address,
			FullPublicKey: item.FullPublicKey,
			PrivateKey:    item.PrivateKey,
			Idx:           item.Idx,
			AddrStatus:    item.AddrStatus,
		})
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to call InsertEthAccountKey(): %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateAddrStatus updates addr_status
func (r *ETHAccountKeyRepositorySqlc) UpdateAddrStatus(
	accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus, privateKeys []string,
) (int64, error) {
	ctx := context.Background()

	tx, err := r.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := r.queries.WithTx(tx)

	var totalAffected int64
	// sqlc doesn't support IN clauses with variable arguments, so update one at a time
	for _, privateKey := range privateKeys {
		result, err := qtx.UpdateEthAccountKeyAddrStatus(ctx, sqlcgen.UpdateEthAccountKeyAddrStatusParams{
			AddrStatus: addrStatus.Int8(),
			UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
			Account:    sqlcgen.EthAccountKeyAccount(accountType.String()),
			PrivateKey: privateKey,
		})
		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("failed to call UpdateEthAccountKeyAddrStatus(): %w", err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
		}
		totalAffected += affected
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return totalAffected, nil
}
