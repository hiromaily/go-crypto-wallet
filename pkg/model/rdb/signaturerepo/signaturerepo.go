package signaturerepo

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type SignatureRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// SignatureRepository
func NewSignatureRepository(db *sqlx.DB, logger *zap.Logger) *SignatureRepository {
	return &SignatureRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SignatureRepository) Close() error {
	return r.db.Close()
}
