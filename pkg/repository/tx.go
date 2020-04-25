package repository

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

type TxRepository interface {
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

type txRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

func NewTxRepository(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger *zap.Logger) TxRepository {
	return &txRepository{
		dbConn:       dbConn,
		tableName:    "tx",
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
// - replaced from GetTxByID
func (r *txRepository) GetOne(id int64) (*models.TX, error) {
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

// GetCount get count by hex string
// - replaced from GetTxCountByUnsignedHex
func (r *txRepository) GetCountByUnsignedHex(actionType action.ActionType, hex string) (int64, error) {
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

// GetOne get one record by ID
// - replaced from GetTxIDBySentHash
func (r *txRepository) GetTxIDBySentHash(actionType action.ActionType, hash string) (int64, error) {
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

func (r *txRepository) getTxIDByUnsignedHexTx(actionType action.ActionType, hex string) (int64, error) {
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

// GetSentHashTx
// - replaced from GetSentTxHashByTxTypeSent, GetSentTxHashByTxTypeDone
func (r *txRepository) GetSentHashTx(actionType action.ActionType, txType tx.TxType) ([]string, error) {
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

// InsertUnsignedTx
// - replaced from InsertTxForUnsigned()
func (r *txRepository) InsertUnsignedTx(actionType action.ActionType, txItem *models.TX) (int64, error) {
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

// Update
func (r *txRepository) Update(txItem *models.TX) (int64, error) {
	ctx := context.Background()
	return txItem.Update(ctx, r.dbConn, boil.Infer())
}

// UpdateAfterTxSent
// - replaced from UpdateTxAfterSent
func (r *txRepository) UpdateAfterTxSent(
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

// UpdateTxType
// - replaced from UpdateTxTypeNotifiedByID
func (r *txRepository) UpdateTxType(id int64, txType tx.TxType) (int64, error) {
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

// UpdateTxTypeBySentHashTx
// - replaced from UpdateTxTypeDoneByTxHash
func (r *txRepository) UpdateTxTypeBySentHashTx(actionType action.ActionType, txType tx.TxType, sentHashTx string) (int64, error) {
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

// Delete
func (r *txRepository) DeleteAll() (int64, error) {
	ctx := context.Background()
	txItems, _ := models.Txes().All(ctx, r.dbConn)
	return txItems.DeleteAll(ctx, r.dbConn)
}
