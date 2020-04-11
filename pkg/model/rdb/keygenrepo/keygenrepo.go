package keygenrepo

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type KeygenRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// KeygenRepository
func NewKeygenRepository(db *sqlx.DB, logger *zap.Logger) *KeygenRepository {
	return &KeygenRepository{
		db:     db,
		logger: logger,
	}
}

func (r *KeygenRepository) Close() error {
	return r.db.Close()
}
