package cold

import (
	"context"
	"database/sql"
	"fmt"

	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/database/sqlc"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

// SeedRepositorySqlc is repository for seed table using sqlc
type SeedRepositorySqlc struct {
	queries      *sqlc.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewSeedRepositorySqlc returns SeedRepositorySqlc object
func NewSeedRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *SeedRepositorySqlc {
	return &SeedRepositorySqlc{
		queries:      sqlc.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne returns one record
func (r *SeedRepositorySqlc) GetOne() (*models.Seed, error) {
	ctx := context.Background()

	seed, err := r.queries.GetSeed(ctx, sqlc.SeedCoin(r.coinTypeCode.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to call GetSeed(): %w", err)
	}

	return convertSqlcSeedToModel(&seed), nil
}

// Insert inserts record
func (r *SeedRepositorySqlc) Insert(strSeed string) error {
	ctx := context.Background()

	_, err := r.queries.InsertSeed(ctx, sqlc.InsertSeedParams{
		Coin: sqlc.SeedCoin(r.coinTypeCode.String()),
		Seed: strSeed,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertSeed(): %w", err)
	}

	return nil
}

// Helper functions

func convertSqlcSeedToModel(seed *sqlc.Seed) *models.Seed {
	return &models.Seed{
		ID:        seed.ID,
		Coin:      string(seed.Coin),
		Seed:      seed.Seed,
		UpdatedAt: convertSQLNullTimeToNullTime(seed.UpdatedAt),
	}
}
