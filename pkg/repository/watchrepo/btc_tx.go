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

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// BTCTxRepositorier is BTCTxRepository interface
type BTCTxRepositorier interface {
	GetOne(id int64) (*models.BTCTX, error)
	GetCountByUnsignedHex(actionType action.ActionType, hex string) (int64, error)
	GetTxIDBySentHash(actionType action.ActionType, hash string) (int64, error)
	GetSentHashTx(actionType action.ActionType, txType tx.TxType) ([]string, error)
	InsertUnsignedTx(actionType action.ActionType, txItem *models.BTCTX) (int64, error)
	Update(txItem *models.BTCTX) (int64, error)
	UpdateAfterTxSent(txID int64, txType tx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType tx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(actionType action.ActionType, txType tx.TxType, sentHashTx string) (int64, error)
	DeleteAll() (int64, error)
}

// BTCTxRepository is repository for tx table
type BTCTxRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

// NewBTCTxRepository returns BTCTxRepository object
func NewBTCTxRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) *BTCTxRepository {
	return &BTCTxRepository{
		dbConn:       dbConn,
		tableName:    "btc_tx",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne returns one record by ID
func (r *BTCTxRepository) GetOne(id int64) (*models.BTCTX, error) {
	//sql := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, r.tableName, id)
	ctx := context.Background()

	//var txItem models.BTCTX
	//err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItem)
	txItem, err := models.FindBTCTX(ctx, r.dbConn, id) //unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindTX()")
		//return nil, errors.Wrapf(err, "failed to call queries.Raw() sql:%s", sql)
	}
	return txItem, nil
}

// GetCountByUnsignedHex returns count by hex string
func (r *BTCTxRepository) GetCountByUnsignedHex(actionType action.ActionType, hex string) (int64, error) {
	ctx := context.Background()
	queries := []qm.QueryMod{
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("action=?", actionType.String()),
		qm.And("unsigned_hex_tx=?", hex),
	}
	count, err := models.BTCTxes(queries...).Count(ctx, r.dbConn)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call models.Txes().Count()")
	}
	return count, nil
}

// GetTxIDBySentHash returns txID  by sentHashTx
func (r *BTCTxRepository) GetTxIDBySentHash(actionType action.ActionType, hash string) (int64, error) {
	//sql := fmt.Sprintf(`SELECT id FROM %s WHERE sent_hash_tx="%s"`, r.tableName, hash)
	ctx := context.Background()

	//var txItem models.BTCTX
	//err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItem)
	txItem, err := models.BTCTxes(
		qm.Select("id"),
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("action=?", actionType.String()),
		qm.And("sent_hash_tx=?", hash),
	).One(ctx, r.dbConn)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call models.Txes().One()")
	}
	return txItem.ID, nil
}

func (r *BTCTxRepository) getTxIDByUnsignedHexTx(actionType action.ActionType, hex string) (int64, error) {
	//sql := fmt.Sprintf(`SELECT id FROM %s WHERE unsigned_hex_tx="%s"`, r.tableName, hex)
	ctx := context.Background()

	//var txItem models.BTCTX
	//err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItem)
	txItem, err := models.BTCTxes(
		qm.Select("id"),
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("action=?", actionType.String()),
		qm.And("unsigned_hex_tx=?", hex),
	).One(ctx, r.dbConn)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call models.Txes().One()")
	}
	return txItem.ID, nil
}

// GetSentHashTx returns list of sent_hash_tx by txType
func (r *BTCTxRepository) GetSentHashTx(actionType action.ActionType, txType tx.TxType) ([]string, error) {
	//sql := fmt.Sprintf(`SELECT sent_hash_tx FROM %s WHERE current_tx_type="%s"`, r.tableName, txType.String())
	ctx := context.Background()

	//var txItems []string
	//err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItems)
	txItems, err := models.BTCTxes(
		qm.Select("sent_hash_tx"),
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("action=?", actionType.String()),
		qm.And("current_tx_type=?", txType.Int8()),
	).All(ctx, r.dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.Txes().One()")
	}

	hashes := make([]string, 0, len(txItems))
	for _, txItem := range txItems {
		hashes = append(hashes, txItem.SentHashTX)
	}
	return hashes, nil
}

// InsertUnsignedTx inserts records
func (r *BTCTxRepository) InsertUnsignedTx(actionType action.ActionType, txItem *models.BTCTX) (int64, error) {
	//set coin
	txItem.Coin = r.coinTypeCode.String()

	ctx := context.Background()
	if err := txItem.Insert(ctx, r.dbConn, boil.Infer()); err != nil {
		return 0, errors.Wrap(err, "failed to call txItem.Insert()")
	}
	//TODO: how to get LastInsertId() without implementation
	// search by txItem.UnsignedHexTX
	id, err := r.getTxIDByUnsignedHexTx(actionType, txItem.UnsignedHexTX)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update updates by models.Tx (entire update)
func (r *BTCTxRepository) Update(txItem *models.BTCTX) (int64, error) {
	ctx := context.Background()
	return txItem.Update(ctx, r.dbConn, boil.Infer())
}

// UpdateAfterTxSent updates when tx sent
func (r *BTCTxRepository) UpdateAfterTxSent(
	txID int64,
	txType tx.TxType,
	signedHex,
	sentHashTx string) (int64, error) {

	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.BTCTXColumns.CurrentTXType: txType.Int8(),
		models.BTCTXColumns.SignedHexTX:   signedHex,
		models.BTCTXColumns.SentHashTX:    sentHashTx,
		models.BTCTXColumns.SentUpdatedAt: null.TimeFrom(time.Now()),
	}
	return models.BTCTxes(
		qm.Where("id=?", txID), //unique
	).UpdateAll(ctx, r.dbConn, updCols)
}

// UpdateTxType updates txType
func (r *BTCTxRepository) UpdateTxType(id int64, txType tx.TxType) (int64, error) {
	//UPDATE %s SET current_tx_type=? WHERE id=?
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.BTCTXColumns.CurrentTXType: txType.Int8(),
	}
	return models.BTCTxes(
		qm.Where("id=?", id), //unique
	).UpdateAll(ctx, r.dbConn, updCols)
}

// UpdateTxTypeBySentHashTx updates txType
func (r *BTCTxRepository) UpdateTxTypeBySentHashTx(actionType action.ActionType, txType tx.TxType, sentHashTx string) (int64, error) {
	//UPDATE %s SET current_tx_type=? WHERE sent_hash_tx=?
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.BTCTXColumns.CurrentTXType: txType.Int8(),
	}
	return models.BTCTxes(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("action=?", actionType.String()),
		qm.And("sent_hash_tx=?", sentHashTx),
	).UpdateAll(ctx, r.dbConn, updCols)
}

// DeleteAll deletes all records
func (r *BTCTxRepository) DeleteAll() (int64, error) {
	ctx := context.Background()
	txItems, _ := models.BTCTxes().All(ctx, r.dbConn)
	return txItems.DeleteAll(ctx, r.dbConn)
}
