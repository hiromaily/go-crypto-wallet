package watch

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/quagmt/udecimal"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	models "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
)

// TxOutputRepositorySqlc is repository for btc_tx_output table using sqlc
type TxOutputRepositorySqlc struct {
	queries      *sqlc.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewBTCTxOutputRepositorySqlc returns TxOutputRepositorySqlc object
func NewBTCTxOutputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *TxOutputRepositorySqlc {
	return &TxOutputRepositorySqlc{
		queries:      sqlc.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne get one record by ID
func (r *TxOutputRepositorySqlc) GetOne(id int64) (*models.BTCTXOutput, error) {
	ctx := context.Background()

	output, err := r.queries.GetBtcTxOutputByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxOutputByID(): %w", err)
	}

	return convertSqlcBtcTxOutputToModel(&output), nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *TxOutputRepositorySqlc) GetAllByTxID(id int64) ([]*models.BTCTXOutput, error) {
	ctx := context.Background()

	outputs, err := r.queries.GetBtcTxOutputsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxOutputsByTxID(): %w", err)
	}

	result := make([]*models.BTCTXOutput, len(outputs))
	for i, output := range outputs {
		result[i] = convertSqlcBtcTxOutputToModel(&output)
	}

	return result, nil
}

// Insert inserts one record
func (r *TxOutputRepositorySqlc) Insert(txItem *models.BTCTXOutput) error {
	ctx := context.Background()

	_, err := r.queries.InsertBtcTxOutput(ctx, sqlc.InsertBtcTxOutputParams{
		TxID:          txItem.TXID,
		OutputAddress: txItem.OutputAddress,
		OutputAccount: txItem.OutputAccount,
		OutputAmount:  txItem.OutputAmount.String(),
		IsChange:      txItem.IsChange,
		UpdatedAt:     convertNullTimeToSQLNullTime(txItem.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertBtcTxOutput(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *TxOutputRepositorySqlc) InsertBulk(txItems []*models.BTCTXOutput) error {
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}

// Helper functions

func convertSqlcBtcTxOutputToModel(output *sqlc.BtcTxOutput) *models.BTCTXOutput {
	amount, _ := udecimal.Parse(output.OutputAmount)

	return &models.BTCTXOutput{
		ID:            output.ID,
		TXID:          output.TxID,
		OutputAddress: output.OutputAddress,
		OutputAccount: output.OutputAccount,
		OutputAmount:  amount,
		IsChange:      output.IsChange,
		UpdatedAt:     convertSQLNullTimeToNullTime(output.UpdatedAt),
	}
}
