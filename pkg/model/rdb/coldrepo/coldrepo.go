package coldrepo

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ColdRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// KeygenRepository
func NewColdRepository(db *sqlx.DB, logger *zap.Logger) *ColdRepository {
	return &ColdRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ColdRepository) Close() error {
	return r.db.Close()
}
