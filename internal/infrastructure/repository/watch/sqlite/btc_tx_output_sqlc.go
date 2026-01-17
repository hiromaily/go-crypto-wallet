package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
)

// TxOutputRepositorySqlc is repository for btc_tx_output table using sqlc (SQLite)
type TxOutputRepositorySqlc struct {
	db           *sql.DB
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewBTCTxOutputRepositorySqlc returns TxOutputRepositorySqlc object
func NewBTCTxOutputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *TxOutputRepositorySqlc {
	return &TxOutputRepositorySqlc{
		db:           dbConn,
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// convertToBtcTxOutput converts sqlcgen.BtcTxOutput to domain.BtcTxOutput entity
func convertToBtcTxOutput(sqlcOutput *sqlcgen.BtcTxOutput) (*domainBitcoin.BtcTxOutput, error) {
	output := &domainBitcoin.BtcTxOutput{
		ID:            sqlcOutput.ID,
		TxID:          sqlcOutput.TxID,
		OutputAddress: sqlcOutput.OutputAddress,
		OutputAccount: sqlcOutput.OutputAccount,
		OutputAmount:  sqlcOutput.OutputAmount,
		IsChange:      sqlcOutput.IsChange != 0,
	}

	// Parse TEXT timestamp
	if sqlcOutput.UpdatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", sqlcOutput.UpdatedAt.String)
		if err == nil {
			output.UpdatedAt = &t
		}
	}

	return output, nil
}

// convertFromBtcTxOutput converts domain.BtcTxOutput entity to sqlcgen.BtcTxOutput
func convertFromBtcTxOutput(output *domainBitcoin.BtcTxOutput) *sqlcgen.BtcTxOutput {
	sqlcOutput := &sqlcgen.BtcTxOutput{
		ID:            output.ID,
		TxID:          output.TxID,
		OutputAddress: output.OutputAddress,
		OutputAccount: output.OutputAccount,
		OutputAmount:  output.OutputAmount,
		IsChange:      boolToInt64(output.IsChange),
	}

	// Convert time.Time to TEXT format
	if output.UpdatedAt != nil {
		sqlcOutput.UpdatedAt = sql.NullString{
			String: output.UpdatedAt.Format("2006-01-02 15:04:05"),
			Valid:  true,
		}
	}

	return sqlcOutput
}

// GetOne get one record by ID
func (r *TxOutputRepositorySqlc) GetOne(id int64) (*domainBitcoin.BtcTxOutput, error) {
	ctx := context.Background()

	output, err := r.queries.GetBtcTxOutputByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxOutputByID(): %w", err)
	}

	return convertToBtcTxOutput(&output)
}

// GetAllByTxID returns all records by tx ID
func (r *TxOutputRepositorySqlc) GetAllByTxID(id int64) ([]*domainBitcoin.BtcTxOutput, error) {
	ctx := context.Background()

	outputs, err := r.queries.GetBtcTxOutputsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxOutputsByTxID(): %w", err)
	}

	result := make([]*domainBitcoin.BtcTxOutput, 0, len(outputs))
	for i := range outputs {
		output, err := convertToBtcTxOutput(&outputs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert output: %w", err)
		}
		result = append(result, output)
	}

	return result, nil
}

// Insert inserts one record
func (r *TxOutputRepositorySqlc) Insert(txItem *domainBitcoin.BtcTxOutput) error {
	ctx := context.Background()

	sqlcOutput := convertFromBtcTxOutput(txItem)
	_, err := r.queries.InsertBtcTxOutput(ctx, sqlcgen.InsertBtcTxOutputParams{
		TxID:          sqlcOutput.TxID,
		OutputAddress: sqlcOutput.OutputAddress,
		OutputAccount: sqlcOutput.OutputAccount,
		OutputAmount:  sqlcOutput.OutputAmount,
		IsChange:      sqlcOutput.IsChange,
		UpdatedAt:     sqlcOutput.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertBtcTxOutput(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *TxOutputRepositorySqlc) InsertBulk(txItems []*domainBitcoin.BtcTxOutput) error {
	ctx := context.Background()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := r.queries.WithTx(tx)

	for _, item := range txItems {
		sqlcOutput := convertFromBtcTxOutput(item)
		_, err := qtx.InsertBtcTxOutput(ctx, sqlcgen.InsertBtcTxOutputParams{
			TxID:          sqlcOutput.TxID,
			OutputAddress: sqlcOutput.OutputAddress,
			OutputAccount: sqlcOutput.OutputAccount,
			OutputAmount:  sqlcOutput.OutputAmount,
			IsChange:      sqlcOutput.IsChange,
			UpdatedAt:     sqlcOutput.UpdatedAt,
		})
		if err != nil {
			return fmt.Errorf("failed to insert output: %w", err)
		}
	}

	return tx.Commit()
}
