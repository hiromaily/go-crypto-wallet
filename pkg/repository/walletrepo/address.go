package walletrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// AddressRepository is repository for Address table
type AddressRepository interface {
	GetAll(accountType account.AccountType) ([]*models.Address, error)
	GetOneUnAllocated(accountType account.AccountType) (*models.Address, error)
	InsertBulk(items []*models.Address) error
	UpdateIsAllocated(isAllocated bool, Address string) (int64, error)
}

type addressRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewAddressRepository returns NewAddressRepository interface
func NewAddressRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) AddressRepository {
	return &addressRepository{
		dbConn:       dbConn,
		tableName:    "address",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetAll returns all records by account
func (r *addressRepository) GetAll(accountType account.AccountType) ([]*models.Address, error) {
	//sql := "SELECT * FROM %s WHERE account=%s;"
	ctx := context.Background()

	items, err := models.Addresses(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.Addresss().All()")
	}
	return items, nil
}

// GetOneUnAllocated returns one records by is_allocated=false
// - replaced from GetOneUnAllocatedAccountAddressTable
func (r *addressRepository) GetOneUnAllocated(accountType account.AccountType) (*models.Address, error) {
	//sql := "SELECT * FROM %s WHERE is_allocated=false ORDER BY id LIMIT 1;"
	ctx := context.Background()

	item, err := models.Addresses(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.And("is_allocated=?", false),
	).One(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.Addresss().One()")
	}
	return item, nil
}

// Insert inserts multiple records
// - replaced from InsertAccountAddressTable()
func (r *addressRepository) InsertBulk(items []*models.Address) error {
	ctx := context.Background()
	return models.AddressSlice(items).InsertAll(ctx, r.dbConn, boil.Infer())
}

// UpdateIsAllocated updates is_allocated
// - replaced from UpdateIsAllocatedOnAccountAddressTable
func (r *addressRepository) UpdateIsAllocated(isAllocated bool, Address string) (int64, error) {
	//	sql := `UPDATE %s SET is_allocated=:is_allocated, updated_at=:updated_at
	//WHERE wallet_address=:wallet_address`
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.AddressColumns.IsAllocated: isAllocated,
		models.AddressColumns.UpdatedAt:   null.TimeFrom(time.Now()),
	}
	return models.Addresses(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("wallet_address=?", Address),
	).UpdateAll(ctx, r.dbConn, updCols)
}
