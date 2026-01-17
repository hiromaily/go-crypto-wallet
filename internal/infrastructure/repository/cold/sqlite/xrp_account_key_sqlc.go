package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
)

// XRPAccountKeyRepositorySqlc is repository for xrp_account_key table using sqlc
type XRPAccountKeyRepositorySqlc struct {
	queries      *sqlcgen.Queries
	dbConn       *sql.DB
	coinTypeCode domainCoin.CoinTypeCode
}

// NewXRPAccountKeyRepositorySqlc returns XRPAccountKeyRepositorySqlc object
func NewXRPAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *XRPAccountKeyRepositorySqlc {
	return &XRPAccountKeyRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		dbConn:       dbConn,
		coinTypeCode: coinTypeCode,
	}
}

// convertToXRPAccountKey converts sqlcgen.XrpAccountKey to domain.XrpAccountKey entity.
// SECURITY: Handles master seed data - never log the master seed fields.
func convertToXRPAccountKey(sqlcKey *sqlcgen.XrpAccountKey) (*domainXrp.XRPAccountKey, error) {
	addrStatus, err := domainAddress.AddrStatusFromInt8(int8(sqlcKey.AddrStatus))
	if err != nil {
		return nil, fmt.Errorf("invalid addr status in database: %w", err)
	}

	key := &domainXrp.XRPAccountKey{
		ID:               sqlcKey.ID,
		CoinTypeCode:     domainCoin.CoinTypeCode(sqlcKey.Coin),
		Account:          domainAccount.AccountType(sqlcKey.Account),
		AccountID:        sqlcKey.AccountID,
		KeyType:          domainXrp.XRPKeyType(sqlcKey.KeyType),
		MasterKey:        sqlcKey.MasterKey,
		MasterSeed:       sqlcKey.MasterSeed,    // Master seed - NEVER log
		MasterSeedHex:    sqlcKey.MasterSeedHex, // Master seed hex - NEVER log
		PublicKey:        sqlcKey.PublicKey,
		PublicKeyHex:     sqlcKey.PublicKeyHex,
		IsRegularKeyPair: sqlcKey.IsRegularKeyPair == 1,
		AllocatedID:      sqlcKey.AllocatedID,
		AddrStatus:       addrStatus,
	}

	if sqlcKey.UpdatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", sqlcKey.UpdatedAt.String)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp in database: %w", err)
		}
		key.UpdatedAt = &t
	}

	return key, nil
}

// convertFromXRPAccountKey converts domain.XrpAccountKey entity to sqlcgen.XrpAccountKey.
func convertFromXRPAccountKey(key *domainXrp.XRPAccountKey) *sqlcgen.XrpAccountKey {
	sqlcKey := &sqlcgen.XrpAccountKey{
		ID:               key.ID,
		Coin:             key.CoinTypeCode.String(),
		Account:          key.Account.String(),
		AccountID:        key.AccountID,
		KeyType:          int64(key.KeyType),
		MasterKey:        key.MasterKey,
		MasterSeed:       key.MasterSeed,
		MasterSeedHex:    key.MasterSeedHex,
		PublicKey:        key.PublicKey,
		PublicKeyHex:     key.PublicKeyHex,
		IsRegularKeyPair: boolToInt64(key.IsRegularKeyPair),
		AllocatedID:      key.AllocatedID,
		AddrStatus:       int64(key.AddrStatus.Int8()),
	}

	if key.UpdatedAt != nil {
		sqlcKey.UpdatedAt = sql.NullString{String: key.UpdatedAt.Format("2006-01-02 15:04:05"), Valid: true}
	}

	return sqlcKey
}

// GetAllAddrStatus returns all XRPAccountKey by addr_status
func (r *XRPAccountKeyRepositorySqlc) GetAllAddrStatus(
	ctx context.Context, accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
) ([]*domainXrp.XRPAccountKey, error) {
	xrpKeys, err := r.queries.GetXRPAccountKeysByAddrStatus(ctx, sqlcgen.GetXRPAccountKeysByAddrStatusParams{
		Coin:       r.coinTypeCode.String(),
		Account:    accountType.String(),
		AddrStatus: int64(addrStatus.Int8()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPAccountKeysByAddrStatus(): %w", err)
	}

	result := make([]*domainXrp.XRPAccountKey, 0, len(xrpKeys))
	for i := range xrpKeys {
		key, err := convertToXRPAccountKey(&xrpKeys[i])
		if err != nil {
			return nil, err
		}
		result = append(result, key)
	}

	return result, nil
}

// GetSecret returns secret (master_seed)
func (r *XRPAccountKeyRepositorySqlc) GetSecret(
	ctx context.Context, accountType domainAccount.AccountType, addr string,
) (string, error) {
	secret, err := r.queries.GetXRPAccountKeySecret(ctx, sqlcgen.GetXRPAccountKeySecretParams{
		Coin:      r.coinTypeCode.String(),
		Account:   accountType.String(),
		AccountID: addr,
	})
	if err != nil {
		return "", fmt.Errorf("failed to call GetXRPAccountKeySecret(): %w", err)
	}

	return secret, nil
}

// InsertBulk inserts multiple records with transaction atomicity
func (r *XRPAccountKeyRepositorySqlc) InsertBulk(ctx context.Context, items []*domainXrp.XRPAccountKey) error {
	// Begin transaction to ensure atomicity of bulk insert
	tx, err := r.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // Safe to call even if already committed
	}()

	qtx := r.queries.WithTx(tx)

	for _, item := range items {
		sqlcItem := convertFromXRPAccountKey(item)
		_, err := qtx.InsertXRPAccountKey(ctx, sqlcgen.InsertXRPAccountKeyParams{
			Coin:             sqlcItem.Coin,
			Account:          sqlcItem.Account,
			AccountID:        sqlcItem.AccountID,
			KeyType:          sqlcItem.KeyType,
			MasterKey:        sqlcItem.MasterKey,
			MasterSeed:       sqlcItem.MasterSeed,
			MasterSeedHex:    sqlcItem.MasterSeedHex,
			PublicKey:        sqlcItem.PublicKey,
			PublicKeyHex:     sqlcItem.PublicKeyHex,
			IsRegularKeyPair: sqlcItem.IsRegularKeyPair,
			AllocatedID:      sqlcItem.AllocatedID,
			AddrStatus:       sqlcItem.AddrStatus,
		})
		if err != nil {
			return fmt.Errorf("failed to call InsertXRPAccountKey(): %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateAddrStatus updates addr_status with transaction atomicity
func (r *XRPAccountKeyRepositorySqlc) UpdateAddrStatus(
	ctx context.Context,
	accountType domainAccount.AccountType,
	addrStatus domainAddress.AddrStatus,
	accountIDs []string,
) (int64, error) {
	// Begin transaction to ensure atomicity of bulk update
	tx, err := r.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // Safe to call even if already committed
	}()

	qtx := r.queries.WithTx(tx)
	var totalAffected int64

	// Update one at a time since IN clause not supported for multiple updates
	for _, accountID := range accountIDs {
		result, err := qtx.UpdateXRPAccountKeyAddrStatus(ctx, sqlcgen.UpdateXRPAccountKeyAddrStatusParams{
			AddrStatus: int64(addrStatus.Int8()),
			UpdatedAt:  sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
			Coin:       r.coinTypeCode.String(),
			Account:    accountType.String(),
			AccountID:  accountID,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to call UpdateXRPAccountKeyAddrStatus(): %w", err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
		}
		totalAffected += affected
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return totalAffected, nil
}
