package watchrepo

import (
	"context"
	"database/sql"

	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/types"
)

// TxOutputRepositorySqlc is repository for btc_tx_output table using sqlc
type TxOutputRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewBTCTxOutputRepositorySqlc returns TxOutputRepositorySqlc object
func NewBTCTxOutputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger,
) *TxOutputRepositorySqlc {
	return &TxOutputRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
func (r *TxOutputRepositorySqlc) GetOne(id int64) (*models.BTCTXOutput, error) {
	ctx := context.Background()

	output, err := r.queries.GetBtcTxOutputByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetBtcTxOutputByID()")
	}

	return convertSqlcBtcTxOutputToModel(&output), nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *TxOutputRepositorySqlc) GetAllByTxID(id int64) ([]*models.BTCTXOutput, error) {
	ctx := context.Background()

	outputs, err := r.queries.GetBtcTxOutputsByTxID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call GetBtcTxOutputsByTxID()")
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

	_, err := r.queries.InsertBtcTxOutput(ctx, sqlcgen.InsertBtcTxOutputParams{
		TxID:          txItem.TXID,
		OutputAddress: txItem.OutputAddress,
		OutputAccount: txItem.OutputAccount,
		OutputAmount:  txItem.OutputAmount.String(),
		IsChange:      txItem.IsChange,
		UpdatedAt:     convertNullTimeToSqlNullTime(txItem.UpdatedAt),
	})
	if err != nil {
		return errors.Wrap(err, "failed to call InsertBtcTxOutput()")
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

func convertSqlcBtcTxOutputToModel(output *sqlcgen.BtcTxOutput) *models.BTCTXOutput {
	var amount types.Decimal
	_ = amount.UnmarshalText([]byte(output.OutputAmount))

	return &models.BTCTXOutput{
		ID:            output.ID,
		TXID:          output.TxID,
		OutputAddress: output.OutputAddress,
		OutputAccount: output.OutputAccount,
		OutputAmount:  amount,
		IsChange:      output.IsChange,
		UpdatedAt:     convertSqlNullTimeToNullTime(output.UpdatedAt),
	}
}
