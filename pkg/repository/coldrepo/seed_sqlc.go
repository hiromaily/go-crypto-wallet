package coldrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// SeedRepositorySqlc is repository for seed table using sqlc
type SeedRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewSeedRepositorySqlc returns SeedRepositorySqlc object
func NewSeedRepositorySqlc(
	dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger,
) *SeedRepositorySqlc {
	return &SeedRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne returns one record
func (r *SeedRepositorySqlc) GetOne() (*models.Seed, error) {
	ctx := context.Background()

	seed, err := r.queries.GetSeed(ctx, sqlcgen.SeedCoin(r.coinTypeCode.String()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetSeed()")
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
		return errors.Wrap(err, "failed to call InsertSeed()")
	}

	return nil
}

// Helper functions

func convertSqlcSeedToModel(seed *sqlcgen.Seed) *models.Seed {
	return &models.Seed{
		ID:        int8(seed.ID),
		Coin:      string(seed.Coin),
		Seed:      seed.Seed,
		UpdatedAt: convertSQLNullTimeToNullTime(seed.UpdatedAt),
	}
}
