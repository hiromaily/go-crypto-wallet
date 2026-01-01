package cold

import (
	"context"
	"database/sql"
	"fmt"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// SeedRepositorySqlc is repository for seed table using sqlc
type SeedRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewSeedRepositorySqlc returns SeedRepositorySqlc object
func NewSeedRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *SeedRepositorySqlc {
	return &SeedRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne returns one record
func (r *SeedRepositorySqlc) GetOne(ctx context.Context) (*sqlcgen.Seed, error) {
	seed, err := r.queries.GetSeed(ctx, sqlcgen.SeedCoin(r.coinTypeCode.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to call GetSeed(): %w", err)
	}

	return &seed, nil
}

// Insert inserts record
func (r *SeedRepositorySqlc) Insert(ctx context.Context, strSeed string) error {
	_, err := r.queries.InsertSeed(ctx, sqlcgen.InsertSeedParams{
		Coin: sqlcgen.SeedCoin(r.coinTypeCode.String()),
		Seed: strSeed,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertSeed(): %w", err)
	}

	return nil
}
