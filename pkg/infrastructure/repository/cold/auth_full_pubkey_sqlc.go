package cold

import (
	"context"
	"database/sql"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/database/sqlc"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

// AuthFullPubkeyRepositorySqlc is repository for auth_fullpubkey table using sqlc
type AuthFullPubkeyRepositorySqlc struct {
	queries      *sqlc.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAuthFullPubkeyRepositorySqlc returns AuthFullPubkeyRepositorySqlc object
func NewAuthFullPubkeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *AuthFullPubkeyRepositorySqlc {
	return &AuthFullPubkeyRepositorySqlc{
		queries:      sqlc.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne returns one record by authType
func (r *AuthFullPubkeyRepositorySqlc) GetOne(authType domainAccount.AuthType) (*models.AuthFullpubkey, error) {
	ctx := context.Background()

	authPubkey, err := r.queries.GetAuthFullPubkey(ctx, sqlc.GetAuthFullPubkeyParams{
		Coin:        sqlc.AuthFullpubkeyCoin(r.coinTypeCode.String()),
		AuthAccount: authType.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAuthFullPubkey(): %w", err)
	}

	return convertSqlcAuthFullPubkeyToModel(&authPubkey), nil
}

// Insert inserts record
func (r *AuthFullPubkeyRepositorySqlc) Insert(authType domainAccount.AuthType, fullPubKey string) error {
	ctx := context.Background()

	_, err := r.queries.InsertAuthFullPubkey(ctx, sqlc.InsertAuthFullPubkeyParams{
		Coin:          sqlc.AuthFullpubkeyCoin(r.coinTypeCode.String()),
		AuthAccount:   authType.String(),
		FullPublicKey: fullPubKey,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertAuthFullPubkey(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *AuthFullPubkeyRepositorySqlc) InsertBulk(items []*models.AuthFullpubkey) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertAuthFullPubkey(ctx, sqlc.InsertAuthFullPubkeyParams{
			Coin:          sqlc.AuthFullpubkeyCoin(item.Coin),
			AuthAccount:   item.AuthAccount,
			FullPublicKey: item.FullPublicKey,
		})
		if err != nil {
			return fmt.Errorf("failed to call InsertAuthFullPubkey(): %w", err)
		}
	}

	return nil
}

// Helper functions

func convertSqlcAuthFullPubkeyToModel(authPubkey *sqlc.AuthFullpubkey) *models.AuthFullpubkey {
	return &models.AuthFullpubkey{
		ID:            authPubkey.ID,
		Coin:          string(authPubkey.Coin),
		AuthAccount:   authPubkey.AuthAccount,
		FullPublicKey: authPubkey.FullPublicKey,
		UpdatedAt:     convertSQLNullTimeToNullTime(authPubkey.UpdatedAt),
	}
}
