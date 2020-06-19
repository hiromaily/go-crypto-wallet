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

// XrpDetailTxRepositorier is XrpDetailTxInputRepository interface
type XrpDetailTxRepositorier interface {
	GetOne(id int64) (*models.XRPDetailTX, error)
	GetAllByTxID(id int64) ([]*models.XRPDetailTX, error)
	GetSentHashTx(txType tx.TxType) ([]string, error)
	Insert(txItem *models.XRPDetailTX) error
	InsertBulk(txItems []*models.XRPDetailTX) error
	UpdateAfterTxSent(uuid string, txType tx.TxType, signedTxID, signedTxBlob string, earlistLedgerVersion uint64) (int64, error)
	UpdateTxType(id int64, txType tx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType tx.TxType, sentHashTx string) (int64, error)
}

// XrpDetailTxInputRepository is repository for Xrp_detail_tx table
type XrpDetailTxInputRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewXrpDetailTxInputRepository returns XrpDetailTxInputRepository object
func NewXrpDetailTxInputRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *XrpDetailTxInputRepository {
	return &XrpDetailTxInputRepository{
		dbConn:       dbConn,
		tableName:    "xrp_detail_tx",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
func (r *XrpDetailTxInputRepository) GetOne(id int64) (*models.XRPDetailTX, error) {
	ctx := context.Background()

	txItem, err := models.FindXRPDetailTX(ctx, r.dbConn, id) //unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindTXInput()")
	}
	return txItem, nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *XrpDetailTxInputRepository) GetAllByTxID(id int64) ([]*models.XRPDetailTX, error) {
	ctx := context.Background()
	txItems, err := models.XRPDetailTxes(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("tx_id=?", id),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.XRPDetailTxes().All()")
	}

	return txItems, nil
}

// GetSentHashTx returns list of sent_hash_tx by txType
// - tx_blob is equivalent to sent_hash_tx in XRP
func (r *XrpDetailTxInputRepository) GetSentHashTx(txType tx.TxType) ([]string, error) {
	ctx := context.Background()

	txItems, err := models.XRPDetailTxes(
		qm.Select("tx_blob"),
		qm.InnerJoin("tx on tx.id=xrp_detail_tx.tx_id"),
		qm.Where("tx.coin=?", r.coinTypeCode.String()),
		qm.And("xrp_detail_tx.current_tx_type=?", txType.Int8()),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.XRPDetailTxes().All()")
	}

	hashes := make([]string, 0, len(txItems))
	for _, txItem := range txItems {
		hashes = append(hashes, txItem.TXBlob)
	}
	return hashes, nil
}

// Insert inserts one record
func (r *XrpDetailTxInputRepository) Insert(txItem *models.XRPDetailTX) error {
	ctx := context.Background()
	return txItem.Insert(ctx, r.dbConn, boil.Infer())
}

// InsertBulk inserts multiple records
func (r *XrpDetailTxInputRepository) InsertBulk(txItems []*models.XRPDetailTX) error {
	ctx := context.Background()
	return models.XRPDetailTXSlice(txItems).InsertAll(ctx, r.dbConn, boil.Infer())
}

// UpdateAfterTxSent updates when tx sent
// - xxx is equivalent to SignedHexTX
// - tx_blob is equivalent to sent_hash_tx in XRP
func (r *XrpDetailTxInputRepository) UpdateAfterTxSent(
	uuid string,
	txType tx.TxType,
	signedTxID,
	TxBlob string,
	earlistLedgerVersion uint64) (int64, error) {

	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.XRPDetailTXColumns.CurrentTXType:         txType.Int8(),
		models.XRPDetailTXColumns.SignedTXID:            signedTxID,
		models.XRPDetailTXColumns.TXBlob:                TxBlob,
		models.XRPDetailTXColumns.EarliestLedgerVersion: earlistLedgerVersion,
		models.XRPDetailTXColumns.SentUpdatedAt:         null.TimeFrom(time.Now()),
	}
	return models.XRPDetailTxes(
		qm.Where("uuid=?", uuid), //unique
	).UpdateAll(ctx, r.dbConn, updCols)
}

// UpdateTxType updates txType
func (r *XrpDetailTxInputRepository) UpdateTxType(id int64, txType tx.TxType) (int64, error) {
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.XRPDetailTXColumns.CurrentTXType: txType.Int8(),
	}
	return models.XRPDetailTxes(
		qm.Where("id=?", id), //unique
	).UpdateAll(ctx, r.dbConn, updCols)
}

// UpdateTxTypeBySentHashTx updates txType
func (r *XrpDetailTxInputRepository) UpdateTxTypeBySentHashTx(txType tx.TxType, sentHashTx string) (int64, error) {
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.XRPDetailTXColumns.CurrentTXType: txType.Int8(),
	}
	return models.XRPDetailTxes(
		qm.Where("sent_tx_blob=?", sentHashTx),
	).UpdateAll(ctx, r.dbConn, updCols)
}
