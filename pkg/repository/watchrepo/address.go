package watchrepo

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

// AddressRepositorier is AddressRepository interface
type AddressRepositorier interface {
	GetAll(accountType account.AccountType) ([]*models.Address, error)
	GetAllAddress(accountType account.AccountType) ([]string, error)
	GetOneUnAllocated(accountType account.AccountType) (*models.Address, error)
	InsertBulk(items []*models.Address) error
	UpdateIsAllocated(isAllocated bool, Address string) (int64, error)
}

// AddressRepository is repository for address table
type AddressRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewAddressRepository returns AddressRepository object
func NewAddressRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *AddressRepository {
	return &AddressRepository{
		dbConn:       dbConn,
		tableName:    "address",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetAll returns all records by account
func (r *AddressRepository) GetAll(accountType account.AccountType) ([]*models.Address, error) {
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

// GetAllAddress returns all addresses by account
func (r *AddressRepository) GetAllAddress(accountType account.AccountType) ([]string, error) {
	ctx := context.Background()

	type Response struct {
		Address string `boil:"wallet_address"`
	}
	var response []*Response
	err := models.EthTxes(
		qm.Select("wallet_address"),
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
	).Bind(ctx, r.dbConn, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.EthTxes().Bind()")
	}

	// convert
	addrs := make([]string, 0, len(response))
	for _, val := range response {
		addrs = append(addrs, val.Address)
	}
	return addrs, nil
}

// GetOneUnAllocated returns one records by is_allocated=false
func (r *AddressRepository) GetOneUnAllocated(accountType account.AccountType) (*models.Address, error) {
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

// InsertBulk inserts multiple records
func (r *AddressRepository) InsertBulk(items []*models.Address) error {
	ctx := context.Background()
	return models.AddressSlice(items).InsertAll(ctx, r.dbConn, boil.Infer())
}

// UpdateIsAllocated updates is_allocated
func (r *AddressRepository) UpdateIsAllocated(isAllocated bool, Address string) (int64, error) {
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
