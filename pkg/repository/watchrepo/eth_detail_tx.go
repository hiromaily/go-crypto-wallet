package watchrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/zap"

	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// EthDetailTxRepositorier is EthDetailTxInputRepository interface
type EthDetailTxRepositorier interface {
	GetOne(id int64) (*models.EthDetailTX, error)
	GetAllByTxID(id int64) ([]*models.EthDetailTX, error)
	Insert(txItem *models.EthDetailTX) error
	InsertBulk(txItems []*models.EthDetailTX) error
}

// EthDetailTxInputRepository is repository for eth_detail_tx table
type EthDetailTxInputRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewEthDetailTxInputRepository returns EthDetailTxInputRepository object
func NewEthDetailTxInputRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *EthDetailTxInputRepository {
	return &EthDetailTxInputRepository{
		dbConn:       dbConn,
		tableName:    "eth_detail_tx",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
func (r *EthDetailTxInputRepository) GetOne(id int64) (*models.EthDetailTX, error) {
	ctx := context.Background()

	txItem, err := models.FindEthDetailTX(ctx, r.dbConn, id) //unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindTXInput()")
	}
	return txItem, nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *EthDetailTxInputRepository) GetAllByTxID(id int64) ([]*models.EthDetailTX, error) {
	ctx := context.Background()
	txItems, err := models.EthDetailTxes(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("tx_id=?", id),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.BTCTXInputs().All()")
	}

	return txItems, nil
}

// Insert inserts one record
func (r *EthDetailTxInputRepository) Insert(txItem *models.EthDetailTX) error {
	ctx := context.Background()
	return txItem.Insert(ctx, r.dbConn, boil.Infer())
}

// InsertBulk inserts multiple records
func (r *EthDetailTxInputRepository) InsertBulk(txItems []*models.EthDetailTX) error {
	ctx := context.Background()
	return models.EthDetailTXSlice(txItems).InsertAll(ctx, r.dbConn, boil.Infer())
}
