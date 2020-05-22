package watchrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/zap"

	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// TxOutputRepositorier is TxOutputRepository interface
type TxOutputRepositorier interface {
	GetOne(id int64) (*models.BTCTXOutput, error)
	GetAllByTxID(id int64) ([]*models.BTCTXOutput, error)
	Insert(txItem *models.BTCTXOutput) error
	InsertBulk(txItems []*models.BTCTXOutput) error
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
		tableName:    "btc_tx_output",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
func (r *TxOutputRepository) GetOne(id int64) (*models.BTCTXOutput, error) {
	ctx := context.Background()

	txItem, err := models.FindBTCTXOutput(ctx, r.dbConn, id) //unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindTXOutput()")
	}
	return txItem, nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *TxOutputRepository) GetAllByTxID(id int64) ([]*models.BTCTXOutput, error) {
	ctx := context.Background()
	txItems, err := models.BTCTXOutputs(
		qm.Where("tx_id=?", id),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.BTCTXOutputs().All()")
	}

	return txItems, nil
}

// Insert inserts one record
func (r *TxOutputRepository) Insert(txItem *models.BTCTXOutput) error {
	ctx := context.Background()
	return txItem.Insert(ctx, r.dbConn, boil.Infer())
}

// InsertBulk inserts multiple records
func (r *TxOutputRepository) InsertBulk(txItems []*models.BTCTXOutput) error {
	//FIXME: when comuns includes booleans and both value is included in that comun,
	// this func causes error
	// Error 1136: Column count doesn't match value count at row 5
	//ctx := context.Background()
	//return models.BTCTXOutputSlice(txItems).InsertAll(ctx, r.dbConn, boil.Infer())

	// this code is temporary
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}
