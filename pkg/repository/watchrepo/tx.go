package watchrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// TxRepositorier is TxRepository interface
type TxRepositorier interface {
	GetOne(id int64) (*models.TX, error)
	GetMaxID(actionType action.ActionType) (int64, error)
	InsertUnsignedTx(actionType action.ActionType) (int64, error)
	Update(txItem *models.TX) (int64, error)
	DeleteAll() (int64, error)
}

// TxRepository is repository for tx table
type TxRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewTxRepository returns TxRepository object
func NewTxRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *TxRepository {
	return &TxRepository{
		dbConn:       dbConn,
		tableName:    "tx",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne returns one record by ID
func (r *TxRepository) GetOne(id int64) (*models.TX, error) {
	ctx := context.Background()

	txItem, err := models.FindTX(ctx, r.dbConn, id) //unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindEthTX()")
	}
	return txItem, nil
}

// GetMaxID returns max id
func (r *TxRepository) GetMaxID(actionType action.ActionType) (int64, error) {
	ctx := context.Background()

	type Response struct {
		MaxCount int64 `boil:"max_count"`
	}
	var maxCount Response
	err := models.Txes(
		qm.Select("MAX(id) as max_count"),
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("action=?", actionType.String()),
	).Bind(ctx, r.dbConn, &maxCount)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call models.EthTxes().Bind()")
	}
	return maxCount.MaxCount, nil
}

// InsertUnsignedTx inserts records
func (r *TxRepository) InsertUnsignedTx(actionType action.ActionType) (int64, error) {
	//set coin
	txItem := &models.TX{
		Coin:   r.coinTypeCode.String(),
		Action: actionType.String(),
	}

	ctx := context.Background()
	if err := txItem.Insert(ctx, r.dbConn, boil.Infer()); err != nil {
		return 0, errors.Wrap(err, "failed to call txItem.Insert()")
	}
	id, err := r.GetMaxID(actionType)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update updates by models.Tx (entire update)
func (r *TxRepository) Update(txItem *models.TX) (int64, error) {
	ctx := context.Background()
	return txItem.Update(ctx, r.dbConn, boil.Infer())
}

// DeleteAll deletes all records
func (r *TxRepository) DeleteAll() (int64, error) {
	ctx := context.Background()
	txItems, _ := models.Txes().All(ctx, r.dbConn)
	return txItems.DeleteAll(ctx, r.dbConn)
}
