package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// XRPAccountKeyRepositorySqlc is repository for xrp_account_key table using sqlc
type XRPAccountKeyRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewXRPAccountKeyRepositorySqlc returns XRPAccountKeyRepositorySqlc object
func NewXRPAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *XRPAccountKeyRepositorySqlc {
	return &XRPAccountKeyRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// convertToXRPAccountKey converts sqlcgen.XrpAccountKey to domain.XrpAccountKey entity.
// SECURITY: Handles master seed data - never log the master seed fields.
func convertToXRPAccountKey(sqlcKey *sqlcgen.XrpAccountKey) (*domainXRP.XRPAccountKey, error) {
	addrStatus, err := domainAddress.AddrStatusFromInt8(sqlcKey.AddrStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid addr status in database: %w", err)
	}

	key := &domainXRP.XRPAccountKey{
		ID:               sqlcKey.ID,
		CoinTypeCode:     domainCoin.CoinTypeCode(sqlcKey.Coin),
		Account:          domainAccount.AccountType(sqlcKey.Account),
		AccountID:        sqlcKey.AccountID,
		KeyType:          domainXRP.XRPKeyType(sqlcKey.KeyType),
		MasterKey:        sqlcKey.MasterKey,
		MasterSeed:       sqlcKey.MasterSeed,    // Master seed - NEVER log
		MasterSeedHex:    sqlcKey.MasterSeedHex, // Master seed hex - NEVER log
		PublicKey:        sqlcKey.PublicKey,
		PublicKeyHex:     sqlcKey.PublicKeyHex,
		IsRegularKeyPair: sqlcKey.IsRegularKeyPair,
		AllocatedID:      sqlcKey.AllocatedID,
		AddrStatus:       addrStatus,
	}

	if sqlcKey.UpdatedAt.Valid {
		key.UpdatedAt = &sqlcKey.UpdatedAt.Time
	}

	return key, nil
}

// convertFromXRPAccountKey converts domain.XrpAccountKey entity to sqlcgen.XrpAccountKey.
func convertFromXRPAccountKey(key *domainXRP.XRPAccountKey) *sqlcgen.XrpAccountKey {
	sqlcKey := &sqlcgen.XrpAccountKey{
		ID:               key.ID,
		Coin:             sqlcgen.XrpAccountKeyCoin(key.CoinTypeCode.String()),
		Account:          sqlcgen.XrpAccountKeyAccount(key.Account.String()),
		AccountID:        key.AccountID,
		KeyType:          int8(key.KeyType),
		MasterKey:        key.MasterKey,
		MasterSeed:       key.MasterSeed,
		MasterSeedHex:    key.MasterSeedHex,
		PublicKey:        key.PublicKey,
		PublicKeyHex:     key.PublicKeyHex,
		IsRegularKeyPair: key.IsRegularKeyPair,
		AllocatedID:      key.AllocatedID,
		AddrStatus:       key.AddrStatus.Int8(),
	}

	if key.UpdatedAt != nil {
		sqlcKey.UpdatedAt = sql.NullTime{Time: *key.UpdatedAt, Valid: true}
	}

	return sqlcKey
}

// GetAllAddrStatus returns all XRPAccountKey by addr_status
func (r *XRPAccountKeyRepositorySqlc) GetAllAddrStatus(
	ctx context.Context, accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
) ([]*domainXRP.XRPAccountKey, error) {
	xrpKeys, err := r.queries.GetXRPAccountKeysByAddrStatus(ctx, sqlcgen.GetXRPAccountKeysByAddrStatusParams{
		Coin:       sqlcgen.XrpAccountKeyCoin(r.coinTypeCode.String()),
		Account:    sqlcgen.XrpAccountKeyAccount(accountType.String()),
		AddrStatus: addrStatus.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPAccountKeysByAddrStatus(): %w", err)
	}

	result := make([]*domainXRP.XRPAccountKey, 0, len(xrpKeys))
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
		Coin:      sqlcgen.XrpAccountKeyCoin(r.coinTypeCode.String()),
		Account:   sqlcgen.XrpAccountKeyAccount(accountType.String()),
		AccountID: addr,
	})
	if err != nil {
		return "", fmt.Errorf("failed to call GetXRPAccountKeySecret(): %w", err)
	}

	return secret, nil
}

// InsertBulk inserts multiple records
func (r *XRPAccountKeyRepositorySqlc) InsertBulk(ctx context.Context, items []*domainXRP.XRPAccountKey) error {
	for _, item := range items {
		sqlcItem := convertFromXRPAccountKey(item)
		_, err := r.queries.InsertXRPAccountKey(ctx, sqlcgen.InsertXRPAccountKeyParams{
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

	return nil
}

// UpdateAddrStatus updates addr_status
func (r *XRPAccountKeyRepositorySqlc) UpdateAddrStatus(
	ctx context.Context,
	accountType domainAccount.AccountType,
	addrStatus domainAddress.AddrStatus,
	accountIDs []string,
) (int64, error) {
	var totalAffected int64

	// Update one at a time since IN clause not supported for multiple updates
	for _, accountID := range accountIDs {
		result, err := r.queries.UpdateXRPAccountKeyAddrStatus(ctx, sqlcgen.UpdateXRPAccountKeyAddrStatusParams{
			AddrStatus: addrStatus.Int8(),
			UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
			Coin:       sqlcgen.XrpAccountKeyCoin(r.coinTypeCode.String()),
			Account:    sqlcgen.XrpAccountKeyAccount(accountType.String()),
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

	return totalAffected, nil
}
