package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database"
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

// convertToTransaction converts sqlcgen.Tx to domain.Transaction entity
func convertToTransaction(sqlcTx *sqlcgen.Tx) (*domainTx.Transaction, error) {
	tx := &domainTx.Transaction{
		ID:           sqlcTx.ID,
		CoinTypeCode: domainCoin.CoinTypeCode(sqlcTx.Coin),
		ActionType:   domainTx.ActionType(sqlcTx.Action),
	}

	if sqlcTx.UpdatedAt.Valid {
		tx.UpdatedAt = &sqlcTx.UpdatedAt.Time
	}

	return tx, nil
}

// convertFromTransaction converts domain.Transaction entity to sqlcgen.Tx
func convertFromTransaction(tx *domainTx.Transaction) *sqlcgen.Tx {
	sqlcTx := &sqlcgen.Tx{
		ID:     tx.ID,
		Coin:   sqlcgen.TxCoin(tx.CoinTypeCode.String()),
		Action: sqlcgen.TxAction(tx.ActionType.String()),
	}

	if tx.UpdatedAt != nil {
		sqlcTx.UpdatedAt = sql.NullTime{Time: *tx.UpdatedAt, Valid: true}
	}

	return sqlcTx
}

// GetOne returns one record by ID
func (r *TxRepositorySqlc) GetOne(id int64) (*domainTx.Transaction, error) {
	ctx := context.Background()

	tx, err := r.queries.GetTxByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetTxByID(): %w", err)
	}

	return convertToTransaction(&tx)
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

// Update updates by domain.Transaction (entire update)
func (r *TxRepositorySqlc) Update(txItem *domainTx.Transaction) (int64, error) {
	ctx := context.Background()

	sqlcTx := convertFromTransaction(txItem)
	err := r.queries.UpdateTx(ctx, sqlcgen.UpdateTxParams{
		Coin:      sqlcTx.Coin,
		Action:    sqlcTx.Action,
		UpdatedAt: sqlcTx.UpdatedAt,
		ID:        sqlcTx.ID,
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

// WithTransaction returns a new repository instance that uses the provided transaction
func (r *TxRepositorySqlc) WithTransaction(tx persistence.Transaction) repowatch.TxRepositorier {
	sqlTx := database.UnwrapSQLTx(tx)
	if sqlTx == nil {
		return r // Return the same repository if transaction cannot be unwrapped
	}
	return &TxRepositorySqlc{
		queries:      r.queries.WithTx(sqlTx),
		coinTypeCode: r.coinTypeCode,
	}
}
