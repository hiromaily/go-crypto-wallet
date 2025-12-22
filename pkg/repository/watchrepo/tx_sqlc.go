package watchrepo

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// TxRepositorySqlc is repository for tx table using sqlc
type TxRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewTxRepositorySqlc returns TxRepositorySqlc object
func NewTxRepositorySqlc(dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger) *TxRepositorySqlc {
	return &TxRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne returns one record by ID
func (r *TxRepositorySqlc) GetOne(id int64) (*models.TX, error) {
	ctx := context.Background()

	tx, err := r.queries.GetTxByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetTxByID()")
	}

	return convertSqlcTxToModel(&tx), nil
}

// GetMaxID returns max id
func (r *TxRepositorySqlc) GetMaxID(actionType action.ActionType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.GetMaxTxID(ctx, sqlcgen.GetMaxTxIDParams{
		Coin:   sqlcgen.TxCoin(r.coinTypeCode.String()),
		Action: sqlcgen.TxAction(actionType.String()),
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call GetMaxTxID()")
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
func (r *TxRepositorySqlc) InsertUnsignedTx(actionType action.ActionType) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.InsertTx(ctx, sqlcgen.InsertTxParams{
		Coin:   sqlcgen.TxCoin(r.coinTypeCode.String()),
		Action: sqlcgen.TxAction(actionType.String()),
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call InsertTx()")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get LastInsertId()")
	}

	return id, nil
}

// Update updates by models.Tx (entire update)
func (r *TxRepositorySqlc) Update(txItem *models.TX) (int64, error) {
	ctx := context.Background()

	err := r.queries.UpdateTx(ctx, sqlcgen.UpdateTxParams{
		Coin:      sqlcgen.TxCoin(txItem.Coin),
		Action:    sqlcgen.TxAction(txItem.Action),
		UpdatedAt: convertNullTimeToSQLNullTime(txItem.UpdatedAt),
		ID:        txItem.ID,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to call UpdateTx()")
	}

	return 1, nil // sqlc Update doesn't return rows affected for :exec queries
}

// DeleteAll deletes all records
func (r *TxRepositorySqlc) DeleteAll() (int64, error) {
	ctx := context.Background()

	result, err := r.queries.DeleteAllTx(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to call DeleteAllTx()")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "failed to get RowsAffected()")
	}

	return rowsAffected, nil
}

// Helper functions

func convertSqlcTxToModel(tx *sqlcgen.Tx) *models.TX {
	return &models.TX{
		ID:        tx.ID,
		Coin:      string(tx.Coin),
		Action:    string(tx.Action),
		UpdatedAt: convertSQLNullTimeToNullTime(tx.UpdatedAt),
	}
}
