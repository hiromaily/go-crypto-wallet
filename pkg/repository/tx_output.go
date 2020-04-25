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

type TxOutputRepository interface {
	GetOne(id int64) (*models.TXOutput, error)
	GetAllByTxID(id int64) ([]*models.TXOutput, error)
	Insert(txItem *models.TXOutput) error
	InsertBulk(txItems []*models.TXOutput) error
}

type txOutputRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewTxOutputRepository returns TxOutputRepository interface
func NewTxOutputRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) TxOutputRepository {
	return &txOutputRepository{
		dbConn:       dbConn,
		tableName:    "tx_output",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
// - replaced from GetTxByID
func (r *txOutputRepository) GetOne(id int64) (*models.TXOutput, error) {
	ctx := context.Background()

	txItem, err := models.FindTXOutput(ctx, r.dbConn, id) //unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindTXOutput()")
	}
	return txItem, nil
}

// GetAllByTxID
// - replaced from GetTxOutputByReceiptID
func (r *txOutputRepository) GetAllByTxID(id int64) ([]*models.TXOutput, error) {
	ctx := context.Background()
	txItems, err := models.TXOutputs(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("tx_id=?", id),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.TXOutputs().All()")
	}

	return txItems, nil
}

// Insert inserts one record
func (r *txOutputRepository) Insert(txItem *models.TXOutput) error {
	ctx := context.Background()
	return txItem.Insert(ctx, r.dbConn, boil.Infer())
}

// Insert inserts multiple records
// - replaced from InsertTxOutputForUnsigned()
func (r *txOutputRepository) InsertBulk(txItems []*models.TXOutput) error {
	ctx := context.Background()
	return models.TXOutputSlice(txItems).InsertAll(ctx, r.dbConn, boil.Infer())
}
