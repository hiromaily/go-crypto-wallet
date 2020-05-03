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

	"github.com/hiromaily/go-bitcoin/pkg/action"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// TxRepositorier is TxRepository interface
type TxRepositorier interface {
	GetOne(id int64) (*models.TX, error)
	GetCountByUnsignedHex(actionType action.ActionType, hex string) (int64, error)
	GetTxIDBySentHash(actionType action.ActionType, hash string) (int64, error)
	GetSentHashTx(actionType action.ActionType, txType tx.TxType) ([]string, error)
	InsertUnsignedTx(actionType action.ActionType, txItem *models.TX) (int64, error)
	Update(txItem *models.TX) (int64, error)
	UpdateAfterTxSent(txID int64, txType tx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType tx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(actionType action.ActionType, txType tx.TxType, sentHashTx string) (int64, error)
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
	//sql := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, r.tableName, id)
	ctx := context.Background()

	//var txItem models.TX
	//err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItem)
	txItem, err := models.FindTX(ctx, r.dbConn, id) //unique
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindTX()")
		//return nil, errors.Wrapf(err, "failed to call queries.Raw() sql:%s", sql)
	}
	return txItem, nil
}

// GetCountByUnsignedHex returns count by hex string
func (r *TxRepository) GetCountByUnsignedHex(actionType action.ActionType, hex string) (int64, error) {
	ctx := context.Background()
	queries := []qm.QueryMod{
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("action=?", actionType.String()),
		qm.And("unsigned_hex_tx=?", hex),
	}
	count, err := models.Txes(queries...).Count(ctx, r.dbConn)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call models.Txes().Count()")
	}
	return count, nil
}

// GetTxIDBySentHash returns txID  by sentHashTx
func (r *TxRepository) GetTxIDBySentHash(actionType action.ActionType, hash string) (int64, error) {
	//sql := fmt.Sprintf(`SELECT id FROM %s WHERE sent_hash_tx="%s"`, r.tableName, hash)
	ctx := context.Background()

	//var txItem models.TX
	//err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItem)
	txItem, err := models.Txes(
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

func (r *TxRepository) getTxIDByUnsignedHexTx(actionType action.ActionType, hex string) (int64, error) {
	//sql := fmt.Sprintf(`SELECT id FROM %s WHERE unsigned_hex_tx="%s"`, r.tableName, hex)
	ctx := context.Background()

	//var txItem models.TX
	//err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItem)
	txItem, err := models.Txes(
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
func (r *TxRepository) GetSentHashTx(actionType action.ActionType, txType tx.TxType) ([]string, error) {
	//sql := fmt.Sprintf(`SELECT sent_hash_tx FROM %s WHERE current_tx_type="%s"`, r.tableName, txType.String())
	ctx := context.Background()

	//var txItems []string
	//err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItems)
	txItems, err := models.Txes(
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
func (r *TxRepository) InsertUnsignedTx(actionType action.ActionType, txItem *models.TX) (int64, error) {
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
func (r *TxRepository) Update(txItem *models.TX) (int64, error) {
	ctx := context.Background()
	return txItem.Update(ctx, r.dbConn, boil.Infer())
}

// UpdateAfterTxSent updates when tx sent
func (r *TxRepository) UpdateAfterTxSent(
	txID int64,
	txType tx.TxType,
	signedHex,
	sentHashTx string) (int64, error) {

	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.TXColumns.CurrentTXType: txType.Int8(),
		models.TXColumns.SignedHexTX:   signedHex,
		models.TXColumns.SentHashTX:    sentHashTx,
		models.TXColumns.SentUpdatedAt: null.TimeFrom(time.Now()),
	}
	return models.Txes(
		qm.Where("id=?", txID), //unique
	).UpdateAll(ctx, r.dbConn, updCols)
}

// UpdateTxType updates txType
func (r *TxRepository) UpdateTxType(id int64, txType tx.TxType) (int64, error) {
	//UPDATE %s SET current_tx_type=? WHERE id=?
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.TXColumns.CurrentTXType: txType.Int8(),
	}
	return models.Txes(
		qm.Where("id=?", id), //unique
	).UpdateAll(ctx, r.dbConn, updCols)
}

// UpdateTxTypeBySentHashTx updates txType
func (r *TxRepository) UpdateTxTypeBySentHashTx(actionType action.ActionType, txType tx.TxType, sentHashTx string) (int64, error) {
	//UPDATE %s SET current_tx_type=? WHERE sent_hash_tx=?
	ctx := context.Background()

	// Set updating columns
	updCols := map[string]interface{}{
		models.TXColumns.CurrentTXType: txType.Int8(),
	}
	return models.Txes(
		qm.Where("coin=?", r.coinTypeCode.String()),
		qm.And("action=?", actionType.String()),
		qm.And("sent_hash_tx=?", sentHashTx),
	).UpdateAll(ctx, r.dbConn, updCols)
}

// DeleteAll deletes all records
func (r *TxRepository) DeleteAll() (int64, error) {
	ctx := context.Background()
	txItems, _ := models.Txes().All(ctx, r.dbConn)
	return txItems.DeleteAll(ctx, r.dbConn)
}
