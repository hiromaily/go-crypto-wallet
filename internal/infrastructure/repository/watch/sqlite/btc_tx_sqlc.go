package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
)

// BTCTxRepositorySqlc is repository for btc_tx table using sqlc (SQLite)
type BTCTxRepositorySqlc struct {
	db           *sql.DB
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewBTCTxRepositorySqlc returns BTCTxRepositorySqlc object
func NewBTCTxRepositorySqlc(dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode) *BTCTxRepositorySqlc {
	return &BTCTxRepositorySqlc{
		db:           dbConn,
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// convertToBtcTransaction converts sqlcgen.BtcTx to domain.BtcTransaction entity
func convertToBtcTransaction(sqlcTx *sqlcgen.BtcTx) (*domainBTC.BTCTransaction, error) {
	currentTxType, err := domainTx.TxTypeFromInt8(int8(sqlcTx.CurrentTxType))
	if err != nil {
		return nil, fmt.Errorf("invalid tx type in database: %w", err)
	}

	tx := &domainBTC.BTCTransaction{
		ID:                sqlcTx.ID,
		CoinTypeCode:      domainCoin.CoinTypeCode(sqlcTx.Coin),
		ActionType:        domainTx.ActionType(sqlcTx.Action),
		UnsignedHexTx:     sqlcTx.UnsignedHexTx,
		SignedHexTx:       sqlcTx.SignedHexTx,
		SentHashTx:        sqlcTx.SentHashTx,
		TotalInputAmount:  strconv.FormatFloat(sqlcTx.TotalInputAmount, 'f', -1, 64),
		TotalOutputAmount: strconv.FormatFloat(sqlcTx.TotalOutputAmount, 'f', -1, 64),
		Fee:               strconv.FormatFloat(sqlcTx.Fee, 'f', -1, 64),
		CurrentTxType:     currentTxType,
	}

	// Parse TEXT timestamps
	if sqlcTx.UnsignedUpdatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", sqlcTx.UnsignedUpdatedAt.String)
		if err == nil {
			tx.UnsignedUpdatedAt = &t
		}
	}
	if sqlcTx.SentUpdatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", sqlcTx.SentUpdatedAt.String)
		if err == nil {
			tx.SentUpdatedAt = &t
		}
	}

	return tx, nil
}

// convertFromBtcTransaction converts domain.BtcTransaction entity to sqlcgen.BtcTx
func convertFromBtcTransaction(tx *domainBTC.BTCTransaction) *sqlcgen.BtcTx {
	totalInputAmount, _ := strconv.ParseFloat(tx.TotalInputAmount, 64)
	totalOutputAmount, _ := strconv.ParseFloat(tx.TotalOutputAmount, 64)
	fee, _ := strconv.ParseFloat(tx.Fee, 64)
	sqlcTx := &sqlcgen.BtcTx{
		ID:                tx.ID,
		Coin:              tx.CoinTypeCode.String(),
		Action:            tx.ActionType.String(),
		UnsignedHexTx:     tx.UnsignedHexTx,
		SignedHexTx:       tx.SignedHexTx,
		SentHashTx:        tx.SentHashTx,
		TotalInputAmount:  totalInputAmount,
		TotalOutputAmount: totalOutputAmount,
		Fee:               fee,
		CurrentTxType:     int64(tx.CurrentTxType.Int8()),
	}

	// Convert time.Time to TEXT format
	if tx.UnsignedUpdatedAt != nil {
		sqlcTx.UnsignedUpdatedAt = sql.NullString{
			String: tx.UnsignedUpdatedAt.Format("2006-01-02 15:04:05"),
			Valid:  true,
		}
	}
	if tx.SentUpdatedAt != nil {
		sqlcTx.SentUpdatedAt = sql.NullString{
			String: tx.SentUpdatedAt.Format("2006-01-02 15:04:05"),
			Valid:  true,
		}
	}

	return sqlcTx
}

// GetOne returns one record by ID
func (r *BTCTxRepositorySqlc) GetOne(id int64) (*domainBTC.BTCTransaction, error) {
	ctx := context.Background()

	btcTx, err := r.queries.GetBtcTxByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxByID(): %w", err)
	}

	return convertToBtcTransaction(&btcTx)
}

// GetCountByUnsignedHex returns count by hex string
func (r *BTCTxRepositorySqlc) GetCountByUnsignedHex(actionType domainTx.ActionType, hex string) (int64, error) {
	ctx := context.Background()

	count, err := r.queries.GetBtcTxCountByUnsignedHex(ctx, sqlcgen.GetBtcTxCountByUnsignedHexParams{
		Coin:          r.coinTypeCode.String(),
		Action:        actionType.String(),
		UnsignedHexTx: hex,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call GetBtcTxCountByUnsignedHex(): %w", err)
	}

	return count, nil
}

// GetTxIDBySentHash returns tx ID by sent hash
func (r *BTCTxRepositorySqlc) GetTxIDBySentHash(actionType domainTx.ActionType, hash string) (int64, error) {
	ctx := context.Background()

	id, err := r.queries.GetBtcTxIDBySentHash(ctx, sqlcgen.GetBtcTxIDBySentHashParams{
		Coin:       r.coinTypeCode.String(),
		Action:     actionType.String(),
		SentHashTx: hash,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call GetBtcTxIDBySentHash(): %w", err)
	}

	return id, nil
}

// GetSentHashTx returns sent hash transactions
func (r *BTCTxRepositorySqlc) GetSentHashTx(actionType domainTx.ActionType, txType domainTx.TxType) ([]string, error) {
	ctx := context.Background()

	hashes, err := r.queries.GetBtcTxSentHashList(ctx, sqlcgen.GetBtcTxSentHashListParams{
		Coin:          r.coinTypeCode.String(),
		Action:        actionType.String(),
		CurrentTxType: int64(txType.Int8()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxSentHashList(): %w", err)
	}

	return hashes, nil
}

// InsertUnsignedTx inserts unsigned transaction
func (r *BTCTxRepositorySqlc) InsertUnsignedTx(
	actionType domainTx.ActionType, txItem *domainBTC.BTCTransaction,
) (int64, error) {
	ctx := context.Background()

	totalInputAmount, _ := strconv.ParseFloat(txItem.TotalInputAmount, 64)
	totalOutputAmount, _ := strconv.ParseFloat(txItem.TotalOutputAmount, 64)
	fee, _ := strconv.ParseFloat(txItem.Fee, 64)
	result, err := r.queries.InsertBtcTx(ctx, sqlcgen.InsertBtcTxParams{
		Coin:              r.coinTypeCode.String(),
		Action:            actionType.String(),
		UnsignedHexTx:     txItem.UnsignedHexTx,
		SignedHexTx:       txItem.SignedHexTx,
		SentHashTx:        txItem.SentHashTx,
		TotalInputAmount:  totalInputAmount,
		TotalOutputAmount: totalOutputAmount,
		Fee:               fee,
		CurrentTxType:     int64(txItem.CurrentTxType.Int8()),
		UnsignedUpdatedAt: sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call InsertBtcTx(): %w", err)
	}

	return result.LastInsertId()
}

// Update updates transaction
func (r *BTCTxRepositorySqlc) Update(txItem *domainBTC.BTCTransaction) (int64, error) {
	ctx := context.Background()

	sqlcTx := convertFromBtcTransaction(txItem)
	err := r.queries.UpdateBtcTx(ctx, sqlcgen.UpdateBtcTxParams{
		UnsignedHexTx:     sqlcTx.UnsignedHexTx,
		SignedHexTx:       sqlcTx.SignedHexTx,
		SentHashTx:        sqlcTx.SentHashTx,
		TotalInputAmount:  sqlcTx.TotalInputAmount,
		TotalOutputAmount: sqlcTx.TotalOutputAmount,
		Fee:               sqlcTx.Fee,
		CurrentTxType:     sqlcTx.CurrentTxType,
		ID:                txItem.ID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateBtcTx(): %w", err)
	}

	// UpdateBtcTx doesn't return RowsAffected in SQLite, return 1 for success
	return 1, nil
}

// UpdateAfterTxSent updates after transaction sent
func (r *BTCTxRepositorySqlc) UpdateAfterTxSent(
	txID int64, txType domainTx.TxType, signedHex, sentHashTx string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcTxAfterSent(ctx, sqlcgen.UpdateBtcTxAfterSentParams{
		SignedHexTx:   signedHex,
		SentHashTx:    sentHashTx,
		CurrentTxType: int64(txType.Int8()),
		SentUpdatedAt: sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
		ID:            txID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateBtcTxAfterSent(): %w", err)
	}

	return result.RowsAffected()
}

// UpdateTxType updates transaction type
func (r *BTCTxRepositorySqlc) UpdateTxType(id int64, txType domainTx.TxType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcTxType(ctx, sqlcgen.UpdateBtcTxTypeParams{
		CurrentTxType: int64(txType.Int8()),
		ID:            id,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateBtcTxType(): %w", err)
	}

	return result.RowsAffected()
}

// UpdateTxTypeBySentHashTx updates transaction type by sent hash
func (r *BTCTxRepositorySqlc) UpdateTxTypeBySentHashTx(
	actionType domainTx.ActionType, txType domainTx.TxType, sentHashTx string,
) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdateBtcTxTypeBySentHash(ctx, sqlcgen.UpdateBtcTxTypeBySentHashParams{
		CurrentTxType: int64(txType.Int8()),
		Coin:          r.coinTypeCode.String(),
		Action:        actionType.String(),
		SentHashTx:    sentHashTx,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateBtcTxTypeBySentHash(): %w", err)
	}

	return result.RowsAffected()
}

// DeleteAll deletes all records
func (r *BTCTxRepositorySqlc) DeleteAll() (int64, error) {
	ctx := context.Background()

	result, err := r.queries.DeleteAllBtcTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to call DeleteAllBtcTx(): %w", err)
	}

	return result.RowsAffected()
}
