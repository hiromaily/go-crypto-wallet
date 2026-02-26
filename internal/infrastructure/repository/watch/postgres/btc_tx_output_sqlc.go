package postgres

import (
	"context"
	"database/sql"
	"fmt"

	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/postgres/sqlcgen"
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

// convertToBTCTxOutput converts sqlcgen.BtcTxOutput to domain.BTCTxOutput entity
func convertToBTCTxOutput(sqlcOutput *sqlcgen.BtcTxOutput) (*domainBTC.BTCTxOutput, error) {
	output := &domainBTC.BTCTxOutput{
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

// GetOne get one record by ID
func (r *TxOutputRepositorySqlc) GetOne(id int64) (*domainBTC.BTCTxOutput, error) {
	ctx := context.Background()

	output, err := r.queries.GetBtcTxOutputByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxOutputByID(): %w", err)
	}

	return convertToBTCTxOutput(&output)
}

// GetAllByTxID returns all records searched by tx_id
func (r *TxOutputRepositorySqlc) GetAllByTxID(id int64) ([]*domainBTC.BTCTxOutput, error) {
	ctx := context.Background()

	outputs, err := r.queries.GetBtcTxOutputsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxOutputsByTxID(): %w", err)
	}

	result := make([]*domainBTC.BTCTxOutput, 0, len(outputs))
	for i := range outputs {
		output, err := convertToBTCTxOutput(&outputs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert output at index %d: %w", i, err)
		}
		result = append(result, output)
	}

	return result, nil
}

// Insert inserts one record
func (r *TxOutputRepositorySqlc) Insert(txItem *domainBTC.BTCTxOutput) error {
	ctx := context.Background()

	var updatedAt sql.NullTime
	if txItem.UpdatedAt != nil {
		updatedAt = sql.NullTime{Time: *txItem.UpdatedAt, Valid: true}
	}

	_, err := r.queries.InsertBtcTxOutput(ctx, sqlcgen.InsertBtcTxOutputParams{
		TxID:          txItem.TxID,
		OutputAddress: txItem.OutputAddress,
		OutputAccount: txItem.OutputAccount,
		OutputAmount:  txItem.OutputAmount,
		IsChange:      txItem.IsChange,
		UpdatedAt:     updatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertBtcTxOutput(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *TxOutputRepositorySqlc) InsertBulk(txItems []*domainBTC.BTCTxOutput) error {
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}
