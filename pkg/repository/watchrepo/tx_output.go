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

// TxOutputRepositorier is TxOutputRepository interface
type TxOutputRepositorier interface {
	GetOne(id int64) (*models.TXOutput, error)
	GetAllByTxID(id int64) ([]*models.TXOutput, error)
	Insert(txItem *models.TXOutput) error
	InsertBulk(txItems []*models.TXOutput) error
}

// TxOutputRepository is repository for tx_output table
type TxOutputRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewTxOutputRepository returns TxOutputRepository object
func NewTxOutputRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *TxOutputRepository {
	return &TxOutputRepository{
		dbConn:       dbConn,
		tableName:    "tx_output",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
func (r *TxOutputRepository) GetOne(id int64) (*models.TXOutput, error) {
	ctx := context.Background()

	txItem, err := models.FindTXOutput(ctx, r.dbConn, id) //unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindTXOutput()")
	}
	return txItem, nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *TxOutputRepository) GetAllByTxID(id int64) ([]*models.TXOutput, error) {
	ctx := context.Background()
	txItems, err := models.TXOutputs(
		qm.Where("tx_id=?", id),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.TXOutputs().All()")
	}

	return txItems, nil
}

// Insert inserts one record
func (r *TxOutputRepository) Insert(txItem *models.TXOutput) error {
	ctx := context.Background()
	return txItem.Insert(ctx, r.dbConn, boil.Infer())
}

// InsertBulk inserts multiple records
func (r *TxOutputRepository) InsertBulk(txItems []*models.TXOutput) error {
	//FIXME: when comuns includes booleans and both value is included in that comun,
	// this func causes error
	// Error 1136: Column count doesn't match value count at row 5
	//ctx := context.Background()
	//return models.TXOutputSlice(txItems).InsertAll(ctx, r.dbConn, boil.Infer())

	// this code is temporary
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}
