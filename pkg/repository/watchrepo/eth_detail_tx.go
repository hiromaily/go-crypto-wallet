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

	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// EthDetailTxRepositorier is EthDetailTxInputRepository interface
type EthDetailTxRepositorier interface {
	GetOne(id int64) (*models.EthDetailTX, error)
	GetAllByTxID(id int64) ([]*models.EthDetailTX, error)
	Insert(txItem *models.EthDetailTX) error
	InsertBulk(txItems []*models.EthDetailTX) error
	UpdateAfterTxSent(uuid string, txType tx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType tx.TxType) (int64, error)
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
		return nil, errors.Wrap(err, "failed to call models.EthDetailTxes().All()")
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

// UpdateAfterTxSent updates when tx sent
func (r *EthDetailTxInputRepository) UpdateAfterTxSent(
	uuid string,
	txType tx.TxType,
	signedHex,
	sentHashTx string) (int64, error) {

	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.EthDetailTXColumns.CurrentTXType: txType.Int8(),
		models.EthDetailTXColumns.SignedHexTX:   signedHex,
		models.EthDetailTXColumns.SentHashTX:    sentHashTx,
		models.EthDetailTXColumns.SentUpdatedAt: null.TimeFrom(time.Now()),
	}
	return models.EthDetailTxes(
		qm.Where("uuid=?", uuid), //unique
	).UpdateAll(ctx, r.dbConn, updCols)
}

// UpdateTxType updates txType
func (r *EthDetailTxInputRepository) UpdateTxType(id int64, txType tx.TxType) (int64, error) {
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.EthDetailTXColumns.CurrentTXType: txType.Int8(),
	}
	return models.EthTxes(
		qm.Where("id=?", id), //unique
	).UpdateAll(ctx, r.dbConn, updCols)
}
