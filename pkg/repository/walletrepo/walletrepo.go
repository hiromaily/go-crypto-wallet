package walletrepo

import (
	"database/sql"

	"go.uber.org/zap"
)

// WalletRepository is wallet repository interface
type WalletRepositorier interface {
	Close() error
	BeginTx() (*sql.Tx, error)
	Tx() TxRepository
	TxInput() TxInputRepository
	TxOutput() TxOutputRepository
	PayReq() PaymentRequestRepository
	Addr() AddressRepository
}

// walletRepository is repository for wallet
type WalletRepository struct {
	db     *sql.DB
	logger *zap.Logger
	TxRepository
	TxInputRepository
	TxOutputRepository
	PaymentRequestRepository
	AddressRepository
}

// NewWalletRepository returns WalletRepository
func NewWalletRepository(
	db *sql.DB,
	logger *zap.Logger,
	txRepo TxRepository,
	txInRepo TxInputRepository,
	txOutRepo TxOutputRepository,
	payReqRepo PaymentRequestRepository,
	addrRepo AddressRepository) *WalletRepository {

	return &WalletRepository{
		db:                       db,
		logger:                   logger,
		TxRepository:             txRepo,
		TxInputRepository:        txInRepo,
		TxOutputRepository:       txOutRepo,
		PaymentRequestRepository: payReqRepo,
		AddressRepository:        addrRepo,
	}
}

// Close close db connection
func (r *WalletRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// BeginTx starts transaction
func (r *WalletRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *WalletRepository) Tx() TxRepository {
	return r.TxRepository
}

func (r *WalletRepository) TxInput() TxInputRepository {
	return r.TxInputRepository
}

func (r *WalletRepository) TxOutput() TxOutputRepository {
	return r.TxOutputRepository
}

func (r *WalletRepository) PayReq() PaymentRequestRepository {
	return r.PaymentRequestRepository
}

func (r *WalletRepository) Addr() AddressRepository {
	return r.AddressRepository
}
