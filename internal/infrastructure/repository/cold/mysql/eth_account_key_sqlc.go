package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainEth "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
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

// convertToETHAccountKey converts sqlcgen.EthAccountKey to domain.EthAccountKey entity.
// SECURITY: Handles private key data - never log the private key or account_extended_privkey fields.
func convertToETHAccountKey(sqlcKey *sqlcgen.EthAccountKey) (*domainEth.ETHAccountKey, error) {
	addrStatus, err := domainAddress.AddrStatusFromInt8(sqlcKey.AddrStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid addr status in database: %w", err)
	}

	key := &domainEth.ETHAccountKey{
		ID:            sqlcKey.ID,
		Account:       domainAccount.AccountType(sqlcKey.Account),
		Address:       sqlcKey.Address,
		FullPublicKey: sqlcKey.FullPublicKey,
		PrivateKey:    sqlcKey.PrivateKey, // Private key - NEVER log
		Idx:           sqlcKey.Idx,
		AddrStatus:    addrStatus,
	}

	if sqlcKey.AccountExtendedPrivkey.Valid {
		key.AccountExtendedPrivkey = sqlcKey.AccountExtendedPrivkey.String
	}

	if sqlcKey.UpdatedAt.Valid {
		key.UpdatedAt = &sqlcKey.UpdatedAt.Time
	}

	return key, nil
}

// convertFromETHAccountKey converts domain.EthAccountKey entity to sqlcgen.EthAccountKey.
func convertFromETHAccountKey(key *domainEth.ETHAccountKey) *sqlcgen.EthAccountKey {
	sqlcKey := &sqlcgen.EthAccountKey{
		ID:            key.ID,
		Account:       sqlcgen.EthAccountKeyAccount(key.Account.String()),
		Address:       key.Address,
		FullPublicKey: key.FullPublicKey,
		PrivateKey:    key.PrivateKey,
		Idx:           key.Idx,
		AddrStatus:    key.AddrStatus.Int8(),
	}

	if key.AccountExtendedPrivkey != "" {
		sqlcKey.AccountExtendedPrivkey = sql.NullString{String: key.AccountExtendedPrivkey, Valid: true}
	}

	if key.UpdatedAt != nil {
		sqlcKey.UpdatedAt = sql.NullTime{Time: *key.UpdatedAt, Valid: true}
	}

	return sqlcKey
}

// GetMaxIndex returns max idx
func (r *ETHAccountKeyRepositorySqlc) GetMaxIndex(
	ctx context.Context, accountType domainAccount.AccountType,
) (int64, error) {
	result, err := r.queries.GetMaxETHAccountKeyIndex(ctx, sqlcgen.EthAccountKeyAccount(accountType.String()))
	if err != nil {
		return 0, fmt.Errorf("failed to call GetMaxETHAccountKeyIndex(): %w", err)
	}

	// Type assert interface{} to int64
	if maxIdx, ok := result.(int64); ok {
		return maxIdx, nil
	}

	return 0, nil
}

// GetOneMaxID returns one record by max id
func (r *ETHAccountKeyRepositorySqlc) GetOneMaxID(accountType domainAccount.AccountType,
) (*domainEth.ETHAccountKey, error) {
	ctx := context.Background()

	accountKey, err := r.queries.GetOneETHAccountKeyByMaxID(ctx, sqlcgen.EthAccountKeyAccount(accountType.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to call GetOneETHAccountKeyByMaxID(): %w", err)
	}

	return convertToETHAccountKey(&accountKey)
}

// GetAllAddrStatus returns all EthAccountKey by addr_status
func (r *ETHAccountKeyRepositorySqlc) GetAllAddrStatus(
	accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
) ([]*domainEth.ETHAccountKey, error) {
	ctx := context.Background()

	accountKeys, err := r.queries.GetETHAccountKeysByAddrStatus(ctx, sqlcgen.GetETHAccountKeysByAddrStatusParams{
		Account:    sqlcgen.EthAccountKeyAccount(accountType.String()),
		AddrStatus: addrStatus.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHAccountKeysByAddrStatus(): %w", err)
	}

	result := make([]*domainEth.ETHAccountKey, 0, len(accountKeys))
	for i := range accountKeys {
		key, err := convertToETHAccountKey(&accountKeys[i])
		if err != nil {
			return nil, err
		}
		result = append(result, key)
	}

	return result, nil
}

// GetByAddress returns EthAccountKey by address
func (r *ETHAccountKeyRepositorySqlc) GetByAddress(addr string) (*domainEth.ETHAccountKey, error) {
	ctx := context.Background()

	accountKey, err := r.queries.GetETHAccountKeyByAddress(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetETHAccountKeyByAddress(): %w", err)
	}

	return convertToETHAccountKey(&accountKey)
}

// InsertBulk inserts multiple records
func (r *ETHAccountKeyRepositorySqlc) InsertBulk(items []*domainEth.ETHAccountKey) error {
	ctx := context.Background()

	tx, err := r.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	qtx := r.queries.WithTx(tx)

	for _, item := range items {
		sqlcItem := convertFromETHAccountKey(item)
		_, err := qtx.InsertETHAccountKey(ctx, sqlcgen.InsertETHAccountKeyParams{
			Account:                sqlcItem.Account,
			Address:                sqlcItem.Address,
			FullPublicKey:          sqlcItem.FullPublicKey,
			PrivateKey:             sqlcItem.PrivateKey,
			AccountExtendedPrivkey: sqlcItem.AccountExtendedPrivkey,
			Idx:                    sqlcItem.Idx,
			AddrStatus:             sqlcItem.AddrStatus,
		})
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to call InsertETHAccountKey(): %w", err)
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
		result, err := qtx.UpdateETHAccountKeyAddrStatus(ctx, sqlcgen.UpdateETHAccountKeyAddrStatusParams{
			AddrStatus: addrStatus.Int8(),
			UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
			Account:    sqlcgen.EthAccountKeyAccount(accountType.String()),
			PrivateKey: privateKey,
		})
		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("failed to call UpdateETHAccountKeyAddrStatus(): %w", err)
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
