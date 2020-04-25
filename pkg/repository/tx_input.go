package repository

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

type TxInputRepository interface {
	GetOne(id int64) (*models.TXInput, error)
	GetAllByTxID(id int64) ([]*models.TXInput, error)
	Insert(txItem *models.TXInput) error
	InsertBulk(txItems []*models.TXInput) error
}

type txInputRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

func NewTxInputRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) TxInputRepository {
	return &txInputRepository{
		dbConn:       dbConn,
		tableName:    "tx_input",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
// - replaced from GetTxByID
func (r *txInputRepository) GetOne(id int64) (*models.TXInput, error) {
	ctx := context.Background()

	txItem, err := models.FindTXInput(ctx, r.dbConn, id) //unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindTXInput()")
	}
	return txItem, nil
}

// GetAllByTxID
// - replaced from GetTxInputByReceiptID
func (r *txInputRepository) GetAllByTxID(id int64) ([]*models.TXInput, error) {
	ctx := context.Background()
	txItems, err := models.TXInputs(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("tx_id=?", id),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.TXInputs().All()")
	}

	return txItems, nil
}

// Insert inserts one record
func (r *txInputRepository) Insert(txItem *models.TXInput) error {
	ctx := context.Background()
	return txItem.Insert(ctx, r.dbConn, boil.Infer())
}

// Insert inserts multiple records
// - replaced from InsertTxInputForUnsigned()
func (r *txInputRepository) InsertBulk(txItems []*models.TXInput) error {
	ctx := context.Background()
	return models.TXInputSlice(txItems).InsertAll(ctx, r.dbConn, boil.Infer())
}
