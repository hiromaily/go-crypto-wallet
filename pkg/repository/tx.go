package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/zap"

	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

type TxRepository interface {
	GetOne(id uint64) (*models.TX, error)
	GetCount(hex string) (int64, error)
	GetTxIDBySentHash(hash string) (uint64, error)
	GetSentHashTx(txType tx.TxType) ([]string, error)
	InsertUnsignedTx(txItem *models.TX) (uint64, error)
	Update(txItem *models.TX) (int64, error)
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
func (r *txRepository) GetOne(id uint64) (*models.TX, error) {
	//sql := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, r.tableName, id)
	ctx := context.Background()

	//var txItem models.TX
	//err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItem)
	txItem, err := models.FindTX(ctx, r.dbConn, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call models.FindTX()")
		//return nil, errors.Wrapf(err, "failed to call queries.Raw() sql:%s", sql)
	}
	return txItem, nil
}

// GetCount get count by hex string
// - replaced from GetTxCountByUnsignedHex
func (r *txRepository) GetCount(hex string) (int64, error) {
	ctx := context.Background()
	queries := []qm.QueryMod{
		qm.Where(models.TXColumns.UnsignedHexTX+"=?", hex),
	}
	count, err := models.Txes(queries...).Count(ctx, r.dbConn)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call models.Txes().Count()")
	}
	return count, nil
}

// GetOne get one record by ID
// - replaced from GetTxIDBySentHash
func (r *txRepository) GetTxIDBySentHash(hash string) (uint64, error) {
	sql := fmt.Sprintf(`SELECT id FROM %s WHERE sent_hash_tx="%s"`, r.tableName, hash)
	ctx := context.Background()

	var txItem models.TX
	err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItem)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to call queries.Raw() sql:%s", sql)
	}
	return txItem.ID, nil
}

//OK
func (r *txRepository) getTxIDByUnsignedHexTx(hex string) (uint64, error) {
	sql := fmt.Sprintf(`SELECT id FROM %s WHERE unsigned_hex_tx="%s"`, r.tableName, hex)
	ctx := context.Background()

	var txItem models.TX
	err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItem)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to call queries.Raw() sql:%s", sql)
	}
	return txItem.ID, nil
}

// GetSentHashTx
// - replaced from GetSentTxHashByTxTypeSent, GetSentTxHashByTxTypeDone
func (r *txRepository) GetSentHashTx(txType tx.TxType) ([]string, error) {
	sql := fmt.Sprintf(`SELECT sent_hash_tx FROM %s WHERE current_tx_type="%s"`, r.tableName, txType.String())
	ctx := context.Background()

	var txItems []string
	err := queries.Raw(sql).Bind(ctx, r.dbConn, &txItems)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to call queries.Raw() sql:%s", sql)
	}
	return txItems, nil
}

// InsertUnsignedTx
// - replaced from InsertTxForUnsigned()
func (r *txRepository) InsertUnsignedTx(txItem *models.TX) (uint64, error) {
	//set coin
	txItem.Coin = r.coinTypeCode.String()

	ctx := context.Background()
	if err := txItem.Insert(ctx, r.dbConn, boil.Infer()); err != nil {
		return 0, errors.Wrap(err, "failed to call txItem.Insert()")
	}
	//TODO: how to get LastInsertId() without implementation
	// search by txItem.UnsignedHexTX
	id, err := r.getTxIDByUnsignedHexTx(txItem.UnsignedHexTX)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateTx
// - replaced from UpdateTxAfterSent
func (r *txRepository) Update(txItem *models.TX) (int64, error) {
	ctx := context.Background()
	return txItem.Update(ctx, r.dbConn, boil.Infer())
}

// Delete
func (r *txRepository) DeleteAll() (int64, error) {
	ctx := context.Background()
	txItems, _ := models.Txes().All(ctx, r.dbConn)
	return txItems.DeleteAll(ctx, r.dbConn)
}
