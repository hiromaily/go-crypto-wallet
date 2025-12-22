package coldrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// AuthFullPubkeyRepositorySqlc is repository for auth_fullpubkey table using sqlc
type AuthFullPubkeyRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewAuthFullPubkeyRepositorySqlc returns AuthFullPubkeyRepositorySqlc object
func NewAuthFullPubkeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger,
) *AuthFullPubkeyRepositorySqlc {
	return &AuthFullPubkeyRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne returns one record by authType
func (r *AuthFullPubkeyRepositorySqlc) GetOne(authType account.AuthType) (*models.AuthFullpubkey, error) {
	ctx := context.Background()

	authPubkey, err := r.queries.GetAuthFullPubkey(ctx, sqlcgen.GetAuthFullPubkeyParams{
		Coin:        sqlcgen.AuthFullpubkeyCoin(r.coinTypeCode.String()),
		AuthAccount: authType.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetAuthFullPubkey()")
	}

	return convertSqlcAuthFullPubkeyToModel(&authPubkey), nil
}

// Insert inserts record
func (r *AuthFullPubkeyRepositorySqlc) Insert(authType account.AuthType, fullPubKey string) error {
	ctx := context.Background()

	_, err := r.queries.InsertAuthFullPubkey(ctx, sqlcgen.InsertAuthFullPubkeyParams{
		Coin:          sqlcgen.AuthFullpubkeyCoin(r.coinTypeCode.String()),
		AuthAccount:   authType.String(),
		FullPublicKey: fullPubKey,
	})
	if err != nil {
		return errors.Wrap(err, "failed to call InsertAuthFullPubkey()")
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *AuthFullPubkeyRepositorySqlc) InsertBulk(items []*models.AuthFullpubkey) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertAuthFullPubkey(ctx, sqlcgen.InsertAuthFullPubkeyParams{
			Coin:          sqlcgen.AuthFullpubkeyCoin(item.Coin),
			AuthAccount:   item.AuthAccount,
			FullPublicKey: item.FullPublicKey,
		})
		if err != nil {
			return errors.Wrap(err, "failed to call InsertAuthFullPubkey()")
		}
	}

	return nil
}

// Helper functions

func convertSqlcAuthFullPubkeyToModel(authPubkey *sqlcgen.AuthFullpubkey) *models.AuthFullpubkey {
	return &models.AuthFullpubkey{
		ID:            int16(authPubkey.ID),
		Coin:          string(authPubkey.Coin),
		AuthAccount:   authPubkey.AuthAccount,
		FullPublicKey: authPubkey.FullPublicKey,
		UpdatedAt:     convertSQLNullTimeToNullTime(authPubkey.UpdatedAt),
	}
}
