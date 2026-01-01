package watch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// BTCTxRepositorySqlc is repository for btc_tx table using sqlc
type BTCTxRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewBTCTxRepositorySqlc returns BTCTxRepositorySqlc object
func NewBTCTxRepositorySqlc(dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode) *BTCTxRepositorySqlc {
	return &BTCTxRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne returns one record by ID
func (r *BTCTxRepositorySqlc) GetOne(id int64) (*sqlcgen.BtcTx, error) {
	ctx := context.Background()

	btcTx, err := r.queries.GetBtcTxByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxByID(): %w", err)
	}

	return &btcTx, nil
}

// GetCountByUnsignedHex returns count by hex string
func (r *BTCTxRepositorySqlc) GetCountByUnsignedHex(actionType domainTx.ActionType, hex string) (int64, error) {
	ctx := context.Background()

	count, err := r.queries.GetBtcTxCountByUnsignedHex(ctx, sqlcgen.GetBtcTxCountByUnsignedHexParams{
		Coin:          sqlcgen.BtcTxCoin(r.coinTypeCode.String()),
		Action:        sqlcgen.BtcTxAction(actionType.String()),
		UnsignedHexTx: hex,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call GetBtcTxCountByUnsignedHex(): %w", err)
	}

	return count, nil
}

// GetTxIDBySentHash returns txID by sentHashTx
func (r *BTCTxRepositorySqlc) GetTxIDBySentHash(actionType domainTx.ActionType, hash string) (int64, error) {
	ctx := context.Background()

	id, err := r.queries.GetBtcTxIDBySentHash(ctx, sqlcgen.GetBtcTxIDBySentHashParams{
		Coin:       sqlcgen.BtcTxCoin(r.coinTypeCode.String()),
		Action:     sqlcgen.BtcTxAction(actionType.String()),
		SentHashTx: hash,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call GetBtcTxIDBySentHash(): %w", err)
	}

	return id, nil
}

// GetSentHashTx returns list of sent_hash_tx by txType
func (r *BTCTxRepositorySqlc) GetSentHashTx(actionType domainTx.ActionType, txType domainTx.TxType) ([]string, error) {
	ctx := context.Background()

	hashes, err := r.queries.GetBtcTxSentHashList(ctx, sqlcgen.GetBtcTxSentHashListParams{
		Coin:          sqlcgen.BtcTxCoin(r.coinTypeCode.String()),
		Action:        sqlcgen.BtcTxAction(actionType.String()),
		CurrentTxType: txType.Int8(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxSentHashList(): %w", err)
	}

	return hashes, nil
}

// InsertUnsignedTx inserts records
func (r *BTCTxRepositorySqlc) InsertUnsignedTx(actionType domainTx.ActionType, txItem *sqlcgen.BtcTx) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.InsertBtcTx(ctx, sqlcgen.InsertBtcTxParams{
		Coin:              sqlcgen.BtcTxCoin(r.coinTypeCode.String()),
		Action:            sqlcgen.BtcTxAction(actionType.String()),
		UnsignedHexTx:     txItem.UnsignedHexTx,
		SignedHexTx:       txItem.SignedHexTx,
		SentHashTx:        txItem.SentHashTx,
		TotalInputAmount:  txItem.TotalInputAmount,
		TotalOutputAmount: txItem.TotalOutputAmount,
		Fee:               txItem.Fee,
		CurrentTxType:     txItem.CurrentTxType,
		UnsignedUpdatedAt: txItem.UnsignedUpdatedAt,
		SentUpdatedAt:     txItem.SentUpdatedAt,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call InsertBtcTx(): %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get LastInsertId(): %w", err)
	}

	return id, nil
}

// Update updates by sqlc.BtcTx (entire update)
func (r *BTCTxRepositorySqlc) Update(txItem *sqlcgen.BtcTx) (int64, error) {
	ctx := context.Background()

	err := r.queries.UpdateBtcTx(ctx, sqlcgen.UpdateBtcTxParams{
		Coin:              txItem.Coin,
		Action:            txItem.Action,
		UnsignedHexTx:     txItem.UnsignedHexTx,
		SignedHexTx:       txItem.SignedHexTx,
		SentHashTx:        txItem.SentHashTx,
		TotalInputAmount:  txItem.TotalInputAmount,
		TotalOutputAmount: txItem.TotalOutputAmount,
		Fee:               txItem.Fee,
		CurrentTxType:     txItem.CurrentTxType,
		UnsignedUpdatedAt: txItem.UnsignedUpdatedAt,
		SentUpdatedAt:     txItem.SentUpdatedAt,
		ID:                txItem.ID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateBtcTx(): %w", err)
	}

	return 1, nil
}

// UpdateAfterTxSent updates when tx sent
func (r *BTCTxRepositorySqlc) UpdateAfterTxSent(
	txID int64,
	txType domainTx.TxType,
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
		return 0, fmt.Errorf("failed to call UpdateBtcTxAfterSent(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateTxType updates txType
func (r *BTCTxRepositorySqlc) UpdateTxType(id int64, txType domainTx.TxType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcTxType(ctx, sqlcgen.UpdateBtcTxTypeParams{
		CurrentTxType: txType.Int8(),
		ID:            id,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateBtcTxType(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// UpdateTxTypeBySentHashTx updates txType
func (r *BTCTxRepositorySqlc) UpdateTxTypeBySentHashTx(
	actionType domainTx.ActionType, txType domainTx.TxType, sentHashTx string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcTxTypeBySentHash(ctx, sqlcgen.UpdateBtcTxTypeBySentHashParams{
		CurrentTxType: txType.Int8(),
		Coin:          sqlcgen.BtcTxCoin(r.coinTypeCode.String()),
		Action:        sqlcgen.BtcTxAction(actionType.String()),
		SentHashTx:    sentHashTx,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateBtcTxTypeBySentHash(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// DeleteAll deletes all records
func (r *BTCTxRepositorySqlc) DeleteAll() (int64, error) {
	ctx := context.Background()

	result, err := r.queries.DeleteAllBtcTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to call DeleteAllBtcTx(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}
