package cold

import (
	"context"
	"database/sql"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// AuthFullPubkeyRepositorySqlc is repository for auth_fullpubkey table using sqlc
type AuthFullPubkeyRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewAuthFullPubkeyRepositorySqlc returns AuthFullPubkeyRepositorySqlc object
func NewAuthFullPubkeyRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *AuthFullPubkeyRepositorySqlc {
	return &AuthFullPubkeyRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne returns one record by authType
func (r *AuthFullPubkeyRepositorySqlc) GetOne(authType domainAccount.AuthType) (*sqlcgen.AuthFullpubkey, error) {
	ctx := context.Background()

	authPubkey, err := r.queries.GetAuthFullPubkey(ctx, sqlcgen.GetAuthFullPubkeyParams{
		Coin:        sqlcgen.AuthFullpubkeyCoin(r.coinTypeCode.String()),
		AuthAccount: authType.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAuthFullPubkey(): %w", err)
	}

	return &authPubkey, nil
}

// Insert inserts record
func (r *AuthFullPubkeyRepositorySqlc) Insert(authType domainAccount.AuthType, fullPubKey string) error {
	ctx := context.Background()

	_, err := r.queries.InsertAuthFullPubkey(ctx, sqlcgen.InsertAuthFullPubkeyParams{
		Coin:          sqlcgen.AuthFullpubkeyCoin(r.coinTypeCode.String()),
		AuthAccount:   authType.String(),
		FullPublicKey: fullPubKey,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertAuthFullPubkey(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *AuthFullPubkeyRepositorySqlc) InsertBulk(items []*sqlcgen.AuthFullpubkey) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertAuthFullPubkey(ctx, sqlcgen.InsertAuthFullPubkeyParams{
			Coin:          item.Coin,
			AuthAccount:   item.AuthAccount,
			FullPublicKey: item.FullPublicKey,
		})
		if err != nil {
			return fmt.Errorf("failed to call InsertAuthFullPubkey(): %w", err)
		}
	}

	return nil
}
