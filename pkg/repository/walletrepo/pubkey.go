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

type PubkeyRepository interface {
	//GetAll(accountType account.AccountType) ([]*models.Pubkey, error)
	GetOneUnAllocated(accountType account.AccountType) (*models.Pubkey, error)
	InsertBulk(items []*models.Pubkey) error
	UpdateIsAllocated(isAllocated bool, pubkey string) (int64, error)
}

type pubkeyRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewPubkeyRepository returns NewPubkeyRepository interface
func NewPubkeyRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) PubkeyRepository {
	return &pubkeyRepository{
		dbConn:       dbConn,
		tableName:    "pubkey",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetAll returns all records by account
//func (r *pubkeyRepository) GetAll(accountType account.AccountType) ([]*models.Pubkey, error) {
//	//sql := "SELECT * FROM %s;"
//	ctx := context.Background()
//
//	items, err := models.Pubkeys(
//		qm.Where("coin=?", r.coinTypeCode.String()),
//		qm.And("account=?", accountType.String()),
//	).All(ctx, r.dbConn)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to call models.Pubkeys().All()")
//	}
//	return items, nil
//}

// GetOneUnAllocated returns one records by is_allocated=false
// - replaced from GetOneUnAllocatedAccountPubKeyTable
func (r *pubkeyRepository) GetOneUnAllocated(accountType account.AccountType) (*models.Pubkey, error) {
	//sql := "SELECT * FROM %s WHERE is_allocated=false ORDER BY id LIMIT 1;"
	ctx := context.Background()

	item, err := models.Pubkeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("account=?", accountType.String()),
		qm.And("is_allocated=?", false),
	).One(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.Pubkeys().One()")
	}
	return item, nil
}

// Insert inserts multiple records
// - replaced from InsertAccountPubKeyTable()
func (r *pubkeyRepository) InsertBulk(items []*models.Pubkey) error {
	ctx := context.Background()
	return models.PubkeySlice(items).InsertAll(ctx, r.dbConn, boil.Infer())
}

// UpdateIsAllocated updates is_allocated
// - replaced from UpdateIsAllocatedOnAccountPubKeyTable
func (r *pubkeyRepository) UpdateIsAllocated(isAllocated bool, pubkey string) (int64, error) {
	//	sql := `UPDATE %s SET is_allocated=:is_allocated, updated_at=:updated_at
	//WHERE wallet_address=:wallet_address`
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.PubkeyColumns.IsAllocated: isAllocated,
		models.PubkeyColumns.UpdatedAt:   null.TimeFrom(time.Now()),
	}
	return models.Pubkeys(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("wallet_address=?", pubkey),
	).UpdateAll(ctx, r.dbConn, updCols)
}
