package watch

import (
	"context"
	"database/sql"
	"fmt"

	portsPersistence "github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// TxRepositorySqlc is repository for tx table using sqlc
type TxRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewTxRepositorySqlc returns TxRepositorySqlc object
func NewTxRepositorySqlc(dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode) *TxRepositorySqlc {
	return &TxRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne returns one record by ID
func (r *TxRepositorySqlc) GetOne(id int64) (*sqlcgen.Tx, error) {
	ctx := context.Background()

	tx, err := r.queries.GetTxByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetTxByID(): %w", err)
	}

	return &tx, nil
}

// GetMaxID returns max id
func (r *TxRepositorySqlc) GetMaxID(actionType domainTx.ActionType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.GetMaxTxID(ctx, sqlcgen.GetMaxTxIDParams{
		Coin:   sqlcgen.TxCoin(r.coinTypeCode.String()),
		Action: sqlcgen.TxAction(actionType.String()),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call GetMaxTxID(): %w", err)
	}

	if result == nil {
		return 0, nil
	}

	// Convert interface{} to int64
	if maxID, ok := result.(int64); ok {
		return maxID, nil
	}

	return 0, nil
}

// InsertUnsignedTx inserts records
func (r *TxRepositorySqlc) InsertUnsignedTx(actionType domainTx.ActionType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.InsertTx(ctx, sqlcgen.InsertTxParams{
		Coin:   sqlcgen.TxCoin(r.coinTypeCode.String()),
		Action: sqlcgen.TxAction(actionType.String()),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call InsertTx(): %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get LastInsertId(): %w", err)
	}

	return id, nil
}

// Update updates by sqlcgen.Tx (entire update)
func (r *TxRepositorySqlc) Update(txItem *sqlcgen.Tx) (int64, error) {
	ctx := context.Background()

	err := r.queries.UpdateTx(ctx, sqlcgen.UpdateTxParams{
		Coin:      txItem.Coin,
		Action:    txItem.Action,
		UpdatedAt: txItem.UpdatedAt,
		ID:        txItem.ID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdateTx(): %w", err)
	}

	return 1, nil // sqlc Update doesn't return rows affected for :exec queries
}

// DeleteAll deletes all records
func (r *TxRepositorySqlc) DeleteAll() (int64, error) {
	ctx := context.Background()

	result, err := r.queries.DeleteAllTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to call DeleteAllTx(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// WithTx returns a new repository instance that uses the provided transaction
func (r *TxRepositorySqlc) WithTx(tx *sql.Tx) portsPersistence.TxRepositorier {
	return &TxRepositorySqlc{
		queries:      r.queries.WithTx(tx),
		coinTypeCode: r.coinTypeCode,
	}
}
