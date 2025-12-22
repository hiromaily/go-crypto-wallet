package watchrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/types"
)

// BTCTxRepositorySqlc is repository for btc_tx table using sqlc
type BTCTxRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewBTCTxRepositorySqlc returns BTCTxRepositorySqlc object
func NewBTCTxRepositorySqlc(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger) *BTCTxRepositorySqlc {
	return &BTCTxRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne returns one record by ID
func (r *BTCTxRepositorySqlc) GetOne(id int64) (*models.BTCTX, error) {
	ctx := context.Background()

	btcTx, err := r.queries.GetBtcTxByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetBtcTxByID()")
	}

	return convertSqlcBtcTxToModel(&btcTx), nil
}

// GetCountByUnsignedHex returns count by hex string
func (r *BTCTxRepositorySqlc) GetCountByUnsignedHex(actionType action.ActionType, hex string) (int64, error) {
	ctx := context.Background()

	count, err := r.queries.GetBtcTxCountByUnsignedHex(ctx, sqlcgen.GetBtcTxCountByUnsignedHexParams{
		Coin:          sqlcgen.BtcTxCoin(r.coinTypeCode.String()),
		Action:        sqlcgen.BtcTxAction(actionType.String()),
		UnsignedHexTx: hex,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call GetBtcTxCountByUnsignedHex()")
	}

	return count, nil
}

// GetTxIDBySentHash returns txID by sentHashTx
func (r *BTCTxRepositorySqlc) GetTxIDBySentHash(actionType action.ActionType, hash string) (int64, error) {
	ctx := context.Background()

	id, err := r.queries.GetBtcTxIDBySentHash(ctx, sqlcgen.GetBtcTxIDBySentHashParams{
		Coin:       sqlcgen.BtcTxCoin(r.coinTypeCode.String()),
		Action:     sqlcgen.BtcTxAction(actionType.String()),
		SentHashTx: hash,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call GetBtcTxIDBySentHash()")
	}

	return id, nil
}

// GetSentHashTx returns list of sent_hash_tx by txType
func (r *BTCTxRepositorySqlc) GetSentHashTx(actionType action.ActionType, txType tx.TxType) ([]string, error) {
	ctx := context.Background()

	hashes, err := r.queries.GetBtcTxSentHashList(ctx, sqlcgen.GetBtcTxSentHashListParams{
		Coin:          sqlcgen.BtcTxCoin(r.coinTypeCode.String()),
		Action:        sqlcgen.BtcTxAction(actionType.String()),
		CurrentTxType: int8(txType.Int8()),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetBtcTxSentHashList()")
	}

	return hashes, nil
}

// InsertUnsignedTx inserts records
func (r *BTCTxRepositorySqlc) InsertUnsignedTx(actionType action.ActionType, txItem *models.BTCTX) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.InsertBtcTx(ctx, sqlcgen.InsertBtcTxParams{
		Coin:              sqlcgen.BtcTxCoin(r.coinTypeCode.String()),
		Action:            sqlcgen.BtcTxAction(actionType.String()),
		UnsignedHexTx:     txItem.UnsignedHexTX,
		SignedHexTx:       txItem.SignedHexTX,
		SentHashTx:        txItem.SentHashTX,
		TotalInputAmount:  txItem.TotalInputAmount.String(),
		TotalOutputAmount: txItem.TotalOutputAmount.String(),
		Fee:               txItem.Fee.String(),
		CurrentTxType:     txItem.CurrentTXType,
		UnsignedUpdatedAt: convertNullTimeToSqlNullTime(txItem.UnsignedUpdatedAt),
		SentUpdatedAt:     convertNullTimeToSqlNullTime(txItem.SentUpdatedAt),
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call InsertBtcTx()")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get LastInsertId()")
	}

	return id, nil
}

// Update updates by models.BTCTX (entire update)
func (r *BTCTxRepositorySqlc) Update(txItem *models.BTCTX) (int64, error) {
	ctx := context.Background()

	err := r.queries.UpdateBtcTx(ctx, sqlcgen.UpdateBtcTxParams{
		Coin:              sqlcgen.BtcTxCoin(txItem.Coin),
		Action:            sqlcgen.BtcTxAction(txItem.Action),
		UnsignedHexTx:     txItem.UnsignedHexTX,
		SignedHexTx:       txItem.SignedHexTX,
		SentHashTx:        txItem.SentHashTX,
		TotalInputAmount:  txItem.TotalInputAmount.String(),
		TotalOutputAmount: txItem.TotalOutputAmount.String(),
		Fee:               txItem.Fee.String(),
		CurrentTxType:     txItem.CurrentTXType,
		UnsignedUpdatedAt: convertNullTimeToSqlNullTime(txItem.UnsignedUpdatedAt),
		SentUpdatedAt:     convertNullTimeToSqlNullTime(txItem.SentUpdatedAt),
		ID:                txItem.ID,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateBtcTx()")
	}

	return 1, nil
}

// UpdateAfterTxSent updates when tx sent
func (r *BTCTxRepositorySqlc) UpdateAfterTxSent(
	txID int64,
	txType tx.TxType,
	signedHex,
	sentHashTx string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcTxAfterSent(ctx, sqlcgen.UpdateBtcTxAfterSentParams{
		CurrentTxType: txType.Int8(),
		SignedHexTx:   signedHex,
		SentHashTx:    sentHashTx,
		SentUpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:            txID,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateBtcTxAfterSent()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// UpdateTxType updates txType
func (r *BTCTxRepositorySqlc) UpdateTxType(id int64, txType tx.TxType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcTxType(ctx, sqlcgen.UpdateBtcTxTypeParams{
		CurrentTxType: txType.Int8(),
		ID:            id,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateBtcTxType()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// UpdateTxTypeBySentHashTx updates txType
func (r *BTCTxRepositorySqlc) UpdateTxTypeBySentHashTx(
	actionType action.ActionType, txType tx.TxType, sentHashTx string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcTxTypeBySentHash(ctx, sqlcgen.UpdateBtcTxTypeBySentHashParams{
		CurrentTxType: txType.Int8(),
		Coin:          sqlcgen.BtcTxCoin(r.coinTypeCode.String()),
		Action:        sqlcgen.BtcTxAction(actionType.String()),
		SentHashTx:    sentHashTx,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateBtcTxTypeBySentHash()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// DeleteAll deletes all records
func (r *BTCTxRepositorySqlc) DeleteAll() (int64, error) {
	ctx := context.Background()

	result, err := r.queries.DeleteAllBtcTx(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call DeleteAllBtcTx()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// Helper functions

func convertSqlcBtcTxToModel(btcTx *sqlcgen.BtcTx) *models.BTCTX {
	var totalInputAmount, totalOutputAmount, fee types.Decimal
	_ = totalInputAmount.UnmarshalText([]byte(btcTx.TotalInputAmount))
	_ = totalOutputAmount.UnmarshalText([]byte(btcTx.TotalOutputAmount))
	_ = fee.UnmarshalText([]byte(btcTx.Fee))

	return &models.BTCTX{
		ID:                btcTx.ID,
		Coin:              string(btcTx.Coin),
		Action:            string(btcTx.Action),
		UnsignedHexTX:     btcTx.UnsignedHexTx,
		SignedHexTX:       btcTx.SignedHexTx,
		SentHashTX:        btcTx.SentHashTx,
		TotalInputAmount:  totalInputAmount,
		TotalOutputAmount: totalOutputAmount,
		Fee:               fee,
		CurrentTXType:     btcTx.CurrentTxType,
		UnsignedUpdatedAt: convertSqlNullTimeToNullTime(btcTx.UnsignedUpdatedAt),
		SentUpdatedAt:     convertSqlNullTimeToNullTime(btcTx.SentUpdatedAt),
	}
}
