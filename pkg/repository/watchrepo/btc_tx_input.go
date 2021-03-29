package watchrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// TxInputRepositorier is TxInputRepository interface
type TxInputRepositorier interface {
	GetOne(id int64) (*models.BTCTXInput, error)
	GetAllByTxID(id int64) ([]*models.BTCTXInput, error)
	Insert(txItem *models.BTCTXInput) error
	InsertBulk(txItems []*models.BTCTXInput) error
}

// TxInputRepository is repository for tx_input table
type TxInputRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewBTCTxInputRepository returns TxInputRepository object
func NewBTCTxInputRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *TxInputRepository {
	return &TxInputRepository{
		dbConn:       dbConn,
		tableName:    "btc_tx_input",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
func (r *TxInputRepository) GetOne(id int64) (*models.BTCTXInput, error) {
	ctx := context.Background()

	txItem, err := models.FindBTCTXInput(ctx, r.dbConn, id) // unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindTXInput()")
	}
	return txItem, nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *TxInputRepository) GetAllByTxID(id int64) ([]*models.BTCTXInput, error) {
	ctx := context.Background()
	txItems, err := models.BTCTXInputs(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("tx_id=?", id),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.BTCTXInputs().All()")
	}

	return txItems, nil
}

// Insert inserts one record
func (r *TxInputRepository) Insert(txItem *models.BTCTXInput) error {
	ctx := context.Background()
	return txItem.Insert(ctx, r.dbConn, boil.Infer())
}

// InsertBulk inserts multiple records
func (r *TxInputRepository) InsertBulk(txItems []*models.BTCTXInput) error {
	ctx := context.Background()
	return models.BTCTXInputSlice(txItems).InsertAll(ctx, r.dbConn, boil.Infer())
}
