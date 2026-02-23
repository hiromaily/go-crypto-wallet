package postgres

import (
	"context"
	"database/sql"
	"fmt"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/postgres/sqlcgen"
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

// convertToSeed converts sqlcgen.Seed to domain.Seed entity.
// SECURITY: Handles encrypted seed data - never log the seed field.
func convertToSeed(sqlcSeed *sqlcgen.Seed) (*domainKey.Seed, error) {
	coinTypeCode := domainCoin.CoinTypeCode(interfaceToString(sqlcSeed.Coin))
	if !domainCoin.IsCoinTypeCode(string(coinTypeCode)) {
		return nil, fmt.Errorf("invalid coin type code from database: %v", sqlcSeed.Coin)
	}

	seed := &domainKey.Seed{
		ID:           int8(sqlcSeed.ID),
		CoinTypeCode: coinTypeCode,
		Seed:         sqlcSeed.Seed, // Encrypted seed - NEVER log
	}

	if sqlcSeed.UpdatedAt.Valid {
		seed.UpdatedAt = &sqlcSeed.UpdatedAt.Time
	}

	return seed, nil
}

// GetOne returns one record
func (r *SeedRepositorySqlc) GetOne(ctx context.Context) (*domainKey.Seed, error) {
	sqlcSeed, err := r.queries.GetSeed(ctx, r.coinTypeCode.String())
	if err != nil {
		return nil, fmt.Errorf("failed to call GetSeed(): %w", err)
	}

	return convertToSeed(&sqlcSeed)
}

// Insert inserts record
func (r *SeedRepositorySqlc) Insert(ctx context.Context, strSeed string) error {
	_, err := r.queries.InsertSeed(ctx, sqlcgen.InsertSeedParams{
		Coin: r.coinTypeCode.String(),
		Seed: strSeed,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertSeed(): %w", err)
	}

	return nil
}
