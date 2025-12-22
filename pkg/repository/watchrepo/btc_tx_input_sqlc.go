package watchrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ericlagergren/decimal"

	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// TxInputRepositorySqlc is repository for btc_tx_input table using sqlc
type TxInputRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewBTCTxInputRepositorySqlc returns TxInputRepositorySqlc object
func NewBTCTxInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger,
) *TxInputRepositorySqlc {
	return &TxInputRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetOne get one record by ID
func (r *TxInputRepositorySqlc) GetOne(id int64) (*models.BTCTXInput, error) {
	ctx := context.Background()

	input, err := r.queries.GetBtcTxInputByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxInputByID(): %w", err)
	}

	return convertSqlcBtcTxInputToModel(&input), nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *TxInputRepositorySqlc) GetAllByTxID(id int64) ([]*models.BTCTXInput, error) {
	ctx := context.Background()

	inputs, err := r.queries.GetBtcTxInputsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxInputsByTxID(): %w", err)
	}

	result := make([]*models.BTCTXInput, len(inputs))
	for i, input := range inputs {
		result[i] = convertSqlcBtcTxInputToModel(&input)
	}

	return result, nil
}

// Insert inserts one record
func (r *TxInputRepositorySqlc) Insert(txItem *models.BTCTXInput) error {
	ctx := context.Background()

	_, err := r.queries.InsertBtcTxInput(ctx, sqlcgen.InsertBtcTxInputParams{
		TxID:               txItem.TXID,
		InputTxid:          txItem.InputTxid,
		InputVout:          txItem.InputVout,
		InputAddress:       txItem.InputAddress,
		InputAccount:       txItem.InputAccount,
		InputAmount:        txItem.InputAmount.String(),
		InputConfirmations: txItem.InputConfirmations,
		UpdatedAt:          convertNullTimeToSQLNullTime(txItem.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertBtcTxInput(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *TxInputRepositorySqlc) InsertBulk(txItems []*models.BTCTXInput) error {
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}

// Helper functions

func convertSqlcBtcTxInputToModel(input *sqlcgen.BtcTxInput) *models.BTCTXInput {
	amount := new(decimal.Big)
	_, _ = amount.SetString(input.InputAmount)

	return &models.BTCTXInput{
		ID:                 input.ID,
		TXID:               input.TxID,
		InputTxid:          input.InputTxid,
		InputVout:          input.InputVout,
		InputAddress:       input.InputAddress,
		InputAccount:       input.InputAccount,
		InputAmount:        amount,
		InputConfirmations: input.InputConfirmations,
		UpdatedAt:          convertSQLNullTimeToNullTime(input.UpdatedAt),
	}
}
