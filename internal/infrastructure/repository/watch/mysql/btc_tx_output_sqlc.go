package mysql

import (
	"context"
	"database/sql"
	"fmt"

	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// TxOutputRepositorySqlc is repository for btc_tx_output table using sqlc
type TxOutputRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewBTCTxOutputRepositorySqlc returns TxOutputRepositorySqlc object
func NewBTCTxOutputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *TxOutputRepositorySqlc {
	return &TxOutputRepositorySqlc{
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
		IsChange:      sqlcOutput.IsChange,
	}

	if sqlcOutput.UpdatedAt.Valid {
		output.UpdatedAt = &sqlcOutput.UpdatedAt.Time
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
		IsChange:      output.IsChange,
	}

	if output.UpdatedAt != nil {
		sqlcOutput.UpdatedAt = sql.NullTime{Time: *output.UpdatedAt, Valid: true}
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

// GetAllByTxID returns all records searched by tx_id
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
			return nil, fmt.Errorf("failed to convert output at index %d: %w", i, err)
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
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}
