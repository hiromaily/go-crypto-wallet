package coldrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// XRPAccountKeyRepositorySqlc is repository for xrp_account_key table using sqlc
type XRPAccountKeyRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewXRPAccountKeyRepositorySqlc returns XRPAccountKeyRepositorySqlc object
func NewXRPAccountKeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger,
) *XRPAccountKeyRepositorySqlc {
	return &XRPAccountKeyRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetAllAddrStatus returns all XRPAccountKey by addr_status
func (r *XRPAccountKeyRepositorySqlc) GetAllAddrStatus(
	accountType account.AccountType, addrStatus address.AddrStatus,
) ([]*models.XRPAccountKey, error) {
	ctx := context.Background()

	xrpKeys, err := r.queries.GetXRPAccountKeysByAddrStatus(ctx, sqlcgen.GetXRPAccountKeysByAddrStatusParams{
		Coin:       sqlcgen.XrpAccountKeyCoin(r.coinTypeCode.String()),
		Account:    sqlcgen.XrpAccountKeyAccount(accountType.String()),
		AddrStatus: addrStatus.Int8(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetXRPAccountKeysByAddrStatus()")
	}

	result := make([]*models.XRPAccountKey, len(xrpKeys))
	for i, xrpKey := range xrpKeys {
		result[i] = convertSqlcXRPAccountKeyToModel(&xrpKey)
	}

	return result, nil
}

// GetSecret returns secret (master_seed)
func (r *XRPAccountKeyRepositorySqlc) GetSecret(accountType account.AccountType, addr string) (string, error) {
	ctx := context.Background()

	secret, err := r.queries.GetXRPAccountKeySecret(ctx, sqlcgen.GetXRPAccountKeySecretParams{
		Coin:      sqlcgen.XrpAccountKeyCoin(r.coinTypeCode.String()),
		Account:   sqlcgen.XrpAccountKeyAccount(accountType.String()),
		AccountID: addr,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to call GetXRPAccountKeySecret()")
	}

	return secret, nil
}

// InsertBulk inserts multiple records
func (r *XRPAccountKeyRepositorySqlc) InsertBulk(items []*models.XRPAccountKey) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertXRPAccountKey(ctx, sqlcgen.InsertXRPAccountKeyParams{
			Coin:              sqlcgen.XrpAccountKeyCoin(item.Coin),
			Account:           sqlcgen.XrpAccountKeyAccount(item.Account),
			AccountID:         item.AccountID,
			KeyType:           item.KeyType,
			MasterKey:         item.MasterKey,
			MasterSeed:        item.MasterSeed,
			MasterSeedHex:     item.MasterSeedHex,
			PublicKey:         item.PublicKey,
			PublicKeyHex:      item.PublicKeyHex,
			IsRegularKeyPair:  item.IsRegularKeyPair,
			AllocatedID:       item.AllocatedID,
			AddrStatus:        item.AddrStatus,
		})
		if err != nil {
			return errors.Wrap(err, "failed to call InsertXRPAccountKey()")
		}
	}

	return nil
}

// UpdateAddrStatus updates addr_status
func (r *XRPAccountKeyRepositorySqlc) UpdateAddrStatus(
	accountType account.AccountType, addrStatus address.AddrStatus, accountIDs []string,
) (int64, error) {
	ctx := context.Background()
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
			return 0, errors.Wrap(err, "failed to call UpdateXRPAccountKeyAddrStatus()")
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return 0, errors.Wrap(err, "failed to get RowsAffected()")
		}
		totalAffected += affected
	}

	return totalAffected, nil
}

// Helper functions

func convertSqlcXRPAccountKeyToModel(xrpKey *sqlcgen.XrpAccountKey) *models.XRPAccountKey {
	return &models.XRPAccountKey{
		ID:               xrpKey.ID,
		Coin:             string(xrpKey.Coin),
		Account:          string(xrpKey.Account),
		AccountID:        xrpKey.AccountID,
		KeyType:          xrpKey.KeyType,
		MasterKey:        xrpKey.MasterKey,
		MasterSeed:       xrpKey.MasterSeed,
		MasterSeedHex:    xrpKey.MasterSeedHex,
		PublicKey:        xrpKey.PublicKey,
		PublicKeyHex:     xrpKey.PublicKeyHex,
		IsRegularKeyPair: xrpKey.IsRegularKeyPair,
		AllocatedID:      xrpKey.AllocatedID,
		AddrStatus:       xrpKey.AddrStatus,
		UpdatedAt:        convertSQLNullTimeToNullTime(xrpKey.UpdatedAt),
	}
}
