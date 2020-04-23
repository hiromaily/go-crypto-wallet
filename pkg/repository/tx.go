package repository

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"go.uber.org/zap"

	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

type txRepository struct {
	dbConn       *sql.DB
	tableName    string
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
}

type TxItem struct {
	ID         uint        `boil:"id"`
	MasterID   uint        `boil:"genre_master_id"`
	Code       string      `boil:"genre_code"`
	ParentCode null.String `boil:"parent_code"`
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
	sql := `SELECT * FROM $1 WHERE id=$2`
	ctx := context.Background()

	var txItem models.TX
	err := queries.Raw(sql, r.tableName, id).Bind(ctx, r.dbConn, &txItem)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to call queries.Raw() sql:%s", sql)
	}
	return &txItem, nil
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
func (r *txRepository) GetTxID(hash string) (uint64, error) {
	sql := `SELECT id FROM $1 WHERE sent_hash_tx=$2`
	ctx := context.Background()

	var txItem models.TX
	err := queries.Raw(sql, r.tableName, hash).Bind(ctx, r.dbConn, &txItem)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to call queries.Raw() sql:%s", sql)
	}
	return txItem.ID, nil
}

// GetSentHashTx
// - replaced from GetSentTxHashByTxTypeSent, GetSentTxHashByTxTypeDone
func (r *txRepository) GetSentHashTx(txType tx.TxType) ([]string, error) {
	sql := `SELECT sent_hash_tx FROM $1 WHERE current_tx_type=$2`
	ctx := context.Background()

	var txItems []string
	err := queries.Raw(sql, r.tableName, txType.String()).Bind(ctx, r.dbConn, &txItems)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to call queries.Raw() sql:%s", sql)
	}
	return txItems, nil
}

// InsertUnsignedTx
// - replaced from InsertTxForUnsigned()
func (r *txRepository) InsertUnsignedTx(txItem *models.TX) error {
	ctx := context.Background()
	if err := txItem.Insert(ctx, r.dbConn, boil.Infer()); err != nil {
		return errors.Wrap(err, "failed to call txItem.Insert()")
	}
	//TODO: how to get LastInsertId()??
	return nil
}

// UpdateTx
// - replaced from UpdateTxAfterSent
func (r *txRepository) UpdateTx(txItem *models.TX) (int64, error) {
	ctx := context.Background()
	return txItem.Update(ctx, r.dbConn, boil.Infer())
}
