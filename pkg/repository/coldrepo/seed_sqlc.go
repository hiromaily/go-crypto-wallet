package coldrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
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
func (r *SeedRepositorySqlc) GetOne() (*models.Seed, error) {
	ctx := context.Background()

	seed, err := r.queries.GetSeed(ctx, sqlcgen.SeedCoin(r.coinTypeCode.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to call GetSeed(): %w", err)
	}

	return convertSqlcSeedToModel(&seed), nil
}

// Insert inserts record
func (r *SeedRepositorySqlc) Insert(strSeed string) error {
	ctx := context.Background()

	_, err := r.queries.InsertSeed(ctx, sqlcgen.InsertSeedParams{
		Coin: sqlcgen.SeedCoin(r.coinTypeCode.String()),
		Seed: strSeed,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertSeed(): %w", err)
	}

	return nil
}

// Helper functions

func convertSqlcSeedToModel(seed *sqlcgen.Seed) *models.Seed {
	return &models.Seed{
		ID:        seed.ID,
		Coin:      string(seed.Coin),
		Seed:      seed.Seed,
		UpdatedAt: convertSQLNullTimeToNullTime(seed.UpdatedAt),
	}
}
