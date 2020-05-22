package watchrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// ETHTxRepositorier is EthTxRepository interface
type ETHTxRepositorier interface {
	GetOne(id int64) (*models.EthTX, error)
	GetMaxID(actionType action.ActionType) (int64, error)
	InsertUnsignedTx(actionType action.ActionType) (int64, error)
	Update(txItem *models.EthTX) (int64, error)
	DeleteAll() (int64, error)
}

// ETHTxRepository is repository for eth_tx table
type ETHTxRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewETHTxRepository returns ETHTxRepository object
func NewETHTxRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *ETHTxRepository {
	return &ETHTxRepository{
		dbConn:       dbConn,
		tableName:    "eth_tx",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne returns one record by ID
func (r *ETHTxRepository) GetOne(id int64) (*models.EthTX, error) {
	ctx := context.Background()

	txItem, err := models.FindEthTX(ctx, r.dbConn, id) //unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindEthTX()")
	}
	return txItem, nil
}

// GetMaxID returns max id
func (r *ETHTxRepository) GetMaxID(actionType action.ActionType) (int64, error) {
	ctx := context.Background()

	type Response struct {
		MaxCount int64 `boil:"max_count"`
	}
	var maxCount Response
	err := models.EthTxes(
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
func (r *ETHTxRepository) InsertUnsignedTx(actionType action.ActionType) (int64, error) {
	//set coin
	txItem := &models.EthTX{
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
func (r *ETHTxRepository) Update(txItem *models.EthTX) (int64, error) {
	ctx := context.Background()
	return txItem.Update(ctx, r.dbConn, boil.Infer())
}

// DeleteAll deletes all records
func (r *ETHTxRepository) DeleteAll() (int64, error) {
	ctx := context.Background()
	txItems, _ := models.EthTxes().All(ctx, r.dbConn)
	return txItems.DeleteAll(ctx, r.dbConn)
}
