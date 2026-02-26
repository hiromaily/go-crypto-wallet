package postgres

import (
	"context"
	"database/sql"
	"fmt"

	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/postgres/sqlcgen"
)

// TxInputRepositorySqlc is repository for btc_tx_input table using sqlc
type TxInputRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewBTCTxInputRepositorySqlc returns TxInputRepositorySqlc object
func NewBTCTxInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *TxInputRepositorySqlc {
	return &TxInputRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// convertToBTCTxInput converts sqlcgen.BtcTxInput to domain.BTCTxInput entity
func convertToBTCTxInput(sqlcInput *sqlcgen.BtcTxInput) (*domainBTC.BTCTxInput, error) {
	input := &domainBTC.BTCTxInput{
		ID:                 sqlcInput.ID,
		TxID:               sqlcInput.TxID,
		InputTxid:          sqlcInput.InputTxid,
		InputVout:          uint32(sqlcInput.InputVout),
		InputAddress:       sqlcInput.InputAddress,
		InputAccount:       sqlcInput.InputAccount,
		InputAmount:        sqlcInput.InputAmount,
		InputConfirmations: uint64(sqlcInput.InputConfirmations),
	}

	if sqlcInput.UpdatedAt.Valid {
		input.UpdatedAt = &sqlcInput.UpdatedAt.Time
	}

	return input, nil
}

// GetOne get one record by ID
func (r *TxInputRepositorySqlc) GetOne(id int64) (*domainBTC.BTCTxInput, error) {
	ctx := context.Background()

	input, err := r.queries.GetBtcTxInputByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxInputByID(): %w", err)
	}

	return convertToBTCTxInput(&input)
}

// GetAllByTxID returns all records searched by tx_id
func (r *TxInputRepositorySqlc) GetAllByTxID(id int64) ([]*domainBTC.BTCTxInput, error) {
	ctx := context.Background()

	inputs, err := r.queries.GetBtcTxInputsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxInputsByTxID(): %w", err)
	}

	result := make([]*domainBTC.BTCTxInput, 0, len(inputs))
	for i := range inputs {
		input, err := convertToBTCTxInput(&inputs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert input at index %d: %w", i, err)
		}
		result = append(result, input)
	}

	return result, nil
}

// Insert inserts one record
func (r *TxInputRepositorySqlc) Insert(txItem *domainBTC.BTCTxInput) error {
	ctx := context.Background()

	var updatedAt sql.NullTime
	if txItem.UpdatedAt != nil {
		updatedAt = sql.NullTime{Time: *txItem.UpdatedAt, Valid: true}
	}

	_, err := r.queries.InsertBtcTxInput(ctx, sqlcgen.InsertBtcTxInputParams{
		TxID:               txItem.TxID,
		InputTxid:          txItem.InputTxid,
		InputVout:          int32(txItem.InputVout),
		InputAddress:       txItem.InputAddress,
		InputAccount:       txItem.InputAccount,
		InputAmount:        txItem.InputAmount,
		InputConfirmations: int64(txItem.InputConfirmations),
		UpdatedAt:          updatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to call InsertBtcTxInput(): %w", err)
	}

	return nil
}

// InsertBulk inserts multiple records
func (r *TxInputRepositorySqlc) InsertBulk(txItems []*domainBTC.BTCTxInput) error {
	for _, item := range txItems {
		if err := r.Insert(item); err != nil {
			return err
		}
	}
	return nil
}
