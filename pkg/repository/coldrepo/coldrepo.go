package coldrepo

import (
	"database/sql"

	"go.uber.org/zap"
)

// ColdRepository is cold wallet repository interface
type ColdRepository interface {
	Close() error
	BeginTx() (*sql.Tx, error)
	Seed() SeedRepository
	AccountKey() AccountKeyRepository
	MultisigHistory() MultisigHistoryRepository
}

// coldRepository is repository for cold wallet
type coldRepository struct {
	db     *sql.DB
	logger *zap.Logger
	SeedRepository
	AccountKeyRepository
	MultisigHistoryRepository
}

// NewColdWalletRepository returns ColdRepository
// nolint:golint
func NewColdWalletRepository(
	db *sql.DB,
	logger *zap.Logger,
	seedRepo SeedRepository,
	accountKeyRepo AccountKeyRepository,
	multisigRepo MultisigHistoryRepository) *coldRepository {

	return &coldRepository{
		db:                        db,
		logger:                    logger,
		SeedRepository:            seedRepo,
		AccountKeyRepository:      accountKeyRepo,
		MultisigHistoryRepository: multisigRepo,
	}
}

// Close close db connection
func (r *coldRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// BeginTx starts transaction
func (r *coldRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *coldRepository) Seed() SeedRepository {
	return r.SeedRepository
}

func (r *coldRepository) AccountKey() AccountKeyRepository {
	return r.AccountKeyRepository
}

func (r *coldRepository) MultisigHistory() MultisigHistoryRepository {
	return r.MultisigHistoryRepository
}
