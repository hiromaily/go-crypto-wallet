package cold

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// XRPAccountKeyRepositorySqlc is repository for xrp_account_key table using sqlc
type XRPAccountKeyRepositorySqlc struct {
	queries      *sqlc.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewXRPAccountKeyRepositorySqlc returns XRPAccountKeyRepositorySqlc object
func NewXRPAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *XRPAccountKeyRepositorySqlc {
	return &XRPAccountKeyRepositorySqlc{
		queries:      sqlc.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetAllAddrStatus returns all XRPAccountKey by addr_status
func (r *XRPAccountKeyRepositorySqlc) GetAllAddrStatus(
	ctx context.Context, accountType domainAccount.AccountType, addrStatus address.AddrStatus,
) ([]*sqlc.XrpAccountKey, error) {
	xrpKeys, err := r.queries.GetXRPAccountKeysByAddrStatus(ctx, sqlc.GetXRPAccountKeysByAddrStatusParams{
		Coin:       sqlc.XrpAccountKeyCoin(r.coinTypeCode.String()),
		Account:    sqlc.XrpAccountKeyAccount(accountType.String()),
		AddrStatus: addrStatus.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPAccountKeysByAddrStatus(): %w", err)
	}

	result := make([]*sqlc.XrpAccountKey, len(xrpKeys))
	for i := range xrpKeys {
		result[i] = &xrpKeys[i]
	}

	return result, nil
}

// GetSecret returns secret (master_seed)
func (r *XRPAccountKeyRepositorySqlc) GetSecret(
	ctx context.Context, accountType domainAccount.AccountType, addr string,
) (string, error) {
	secret, err := r.queries.GetXRPAccountKeySecret(ctx, sqlc.GetXRPAccountKeySecretParams{
		Coin:      sqlc.XrpAccountKeyCoin(r.coinTypeCode.String()),
		Account:   sqlc.XrpAccountKeyAccount(accountType.String()),
		AccountID: addr,
	})
	if err != nil {
		return "", fmt.Errorf("failed to call GetXRPAccountKeySecret(): %w", err)
	}

	return secret, nil
}

// InsertBulk inserts multiple records
func (r *XRPAccountKeyRepositorySqlc) InsertBulk(ctx context.Context, items []*sqlc.XrpAccountKey) error {
	for _, item := range items {
		_, err := r.queries.InsertXRPAccountKey(ctx, sqlc.InsertXRPAccountKeyParams{
			Coin:             item.Coin,
			Account:          item.Account,
			AccountID:        item.AccountID,
			KeyType:          item.KeyType,
			MasterKey:        item.MasterKey,
			MasterSeed:       item.MasterSeed,
			MasterSeedHex:    item.MasterSeedHex,
			PublicKey:        item.PublicKey,
			PublicKeyHex:     item.PublicKeyHex,
			IsRegularKeyPair: item.IsRegularKeyPair,
			AllocatedID:      item.AllocatedID,
			AddrStatus:       item.AddrStatus,
		})
		if err != nil {
			return fmt.Errorf("failed to call InsertXRPAccountKey(): %w", err)
		}
	}

	return nil
}

// UpdateAddrStatus updates addr_status
func (r *XRPAccountKeyRepositorySqlc) UpdateAddrStatus(
	ctx context.Context, accountType domainAccount.AccountType, addrStatus address.AddrStatus, accountIDs []string,
) (int64, error) {
	var totalAffected int64

	// Update one at a time since IN clause not supported for multiple updates
	for _, accountID := range accountIDs {
		result, err := r.queries.UpdateXRPAccountKeyAddrStatus(ctx, sqlc.UpdateXRPAccountKeyAddrStatusParams{
			AddrStatus: addrStatus.Int8(),
			UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
			Coin:       sqlc.XrpAccountKeyCoin(r.coinTypeCode.String()),
			Account:    sqlc.XrpAccountKeyAccount(accountType.String()),
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
