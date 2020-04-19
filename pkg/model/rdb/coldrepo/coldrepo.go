package coldrepo

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// ColdRepository is repository for keygen/sign wallet
type ColdRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewColdRepository returns ColdRepository
func NewColdRepository(db *sqlx.DB, logger *zap.Logger) *ColdRepository {
	return &ColdRepository{
		db:     db,
		logger: logger,
	}
}

// Close db connection
func (r *ColdRepository) Close() error {
	return r.db.Close()
}
