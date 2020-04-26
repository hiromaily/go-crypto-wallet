package coldrepo

import (
	"database/sql"

	"go.uber.org/zap"
)

// SignRepository is sign wallet repository interface
type SignRepository interface {
	Close() error
	BeginTx() (*sql.Tx, error)
	Seed() SeedRepository
	AccountKey() AccountKeyRepository
	MultisigHistory() MultisigHistoryRepository
}

// signRepository is repository for sign wallet
type signRepository struct {
	db     *sql.DB
	logger *zap.Logger
	SeedRepository
	AccountKeyRepository
	MultisigHistoryRepository
}

// NewSignWalletRepository returns SignRepository
func NewSignWalletRepository(
	db *sql.DB,
	logger *zap.Logger,
	seedRepo SeedRepository,
	accountKeyRepo AccountKeyRepository,
	multisigRepo MultisigHistoryRepository) SignRepository {

	return &signRepository{
		db:                        db,
		logger:                    logger,
		SeedRepository:            seedRepo,
		AccountKeyRepository:      accountKeyRepo,
		MultisigHistoryRepository: multisigRepo,
	}
}

// Close close db connection
func (r *signRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// BeginTx starts transaction
func (r *signRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *signRepository) Seed() SeedRepository {
	return r.SeedRepository
}

func (r *signRepository) AccountKey() AccountKeyRepository {
	return r.AccountKeyRepository
}

func (r *signRepository) MultisigHistory() MultisigHistoryRepository {
	return r.MultisigHistoryRepository
}
