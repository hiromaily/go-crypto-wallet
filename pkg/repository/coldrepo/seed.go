package coldrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// SeedRepositorier is SeedRepository interface
type SeedRepositorier interface {
	GetOne() (*models.Seed, error)
	Insert(strSeed string) error
}

// SeedRepository is repository for seed table
type SeedRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewSeedRepository returns SeedRepository interface
func NewSeedRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *SeedRepository {
	return &SeedRepository{
		dbConn:       dbConn,
		tableName:    "seed",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne returns one record
func (r *SeedRepository) GetOne() (*models.Seed, error) {
	ctx := context.Background()

	item, err := models.Seeds(
		qm.Where("coin=?", r.coinTypeCode.String()),
	).One(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.Seeds().One()")
	}
	return item, nil
}

// Insert inserts record
func (r *SeedRepository) Insert(strSeed string) error {
	//set coin
	item := &models.Seed{
		Coin: r.coinTypeCode.String(),
		Seed: strSeed,
	}

	ctx := context.Background()
	if err := item.Insert(ctx, r.dbConn, boil.Infer()); err != nil {
		return errors.Wrap(err, "failed to call item.Insert()")
	}

	return nil
}
