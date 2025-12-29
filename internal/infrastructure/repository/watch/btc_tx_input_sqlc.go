package watch

import (
	"context"
	"database/sql"
	"fmt"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
)

// TxInputRepositorySqlc is repository for btc_tx_input table using sqlc
type TxInputRepositorySqlc struct {
	queries      *sqlc.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewBTCTxInputRepositorySqlc returns TxInputRepositorySqlc object
func NewBTCTxInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *TxInputRepositorySqlc {
	return &TxInputRepositorySqlc{
		queries:      sqlc.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetOne get one record by ID
func (r *TxInputRepositorySqlc) GetOne(id int64) (*sqlc.BtcTxInput, error) {
	ctx := context.Background()

	input, err := r.queries.GetBtcTxInputByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxInputByID(): %w", err)
	}

	return &input, nil
}

// GetAllByTxID returns all records searched by tx_id
func (r *TxInputRepositorySqlc) GetAllByTxID(id int64) ([]*sqlc.BtcTxInput, error) {
	ctx := context.Background()

	inputs, err := r.queries.GetBtcTxInputsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxInputsByTxID(): %w", err)
	}

	result := make([]*sqlc.BtcTxInput, len(inputs))
	for i := range inputs {
		result[i] = &inputs[i]
	}

	return result, nil
}

// Insert inserts one record
func (r *TxInputRepositorySqlc) Insert(txItem *sqlc.BtcTxInput) error {
	ctx := context.Background()

	_, err := r.queries.InsertBtcTxInput(ctx, sqlc.InsertBtcTxInputParams{
		TxID:               txItem.TxID,
		InputTxid:          txItem.InputTxid,
		InputVout:          txItem.InputVout,
		InputAddress:       txItem.InputAddress,
		InputAccount:       txItem.InputAccount,
		InputAmount:        txItem.InputAmount,
		InputConfirmations: txItem.InputConfirmations,
		UpdatedAt:          txItem.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertBtcTxInput(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *TxInputRepositorySqlc) InsertBulk(txItems []*sqlc.BtcTxInput) error {
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}
