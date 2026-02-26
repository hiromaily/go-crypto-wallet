package mysql

import (
	"context"
	"database/sql"
	"fmt"

	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
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
		InputVout:          sqlcInput.InputVout,
		InputAddress:       sqlcInput.InputAddress,
		InputAccount:       sqlcInput.InputAccount,
		InputAmount:        sqlcInput.InputAmount,
		InputConfirmations: sqlcInput.InputConfirmations,
	}

	if sqlcInput.UpdatedAt.Valid {
		input.UpdatedAt = &sqlcInput.UpdatedAt.Time
	}

	return input, nil
}

// convertFromBTCTxInput converts domain.BTCTxInput entity to sqlcgen.BtcTxInput
func convertFromBTCTxInput(input *domainBTC.BTCTxInput) *sqlcgen.BtcTxInput {
	sqlcInput := &sqlcgen.BtcTxInput{
		ID:                 input.ID,
		TxID:               input.TxID,
		InputTxid:          input.InputTxid,
		InputVout:          input.InputVout,
		InputAddress:       input.InputAddress,
		InputAccount:       input.InputAccount,
		InputAmount:        input.InputAmount,
		InputConfirmations: input.InputConfirmations,
	}

	if input.UpdatedAt != nil {
		sqlcInput.UpdatedAt = sql.NullTime{Time: *input.UpdatedAt, Valid: true}
	}

	return sqlcInput
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

	sqlcInput := convertFromBTCTxInput(txItem)
	_, err := r.queries.InsertBtcTxInput(ctx, sqlcgen.InsertBtcTxInputParams{
		TxID:               sqlcInput.TxID,
		InputTxid:          sqlcInput.InputTxid,
		InputVout:          sqlcInput.InputVout,
		InputAddress:       sqlcInput.InputAddress,
		InputAccount:       sqlcInput.InputAccount,
		InputAmount:        sqlcInput.InputAmount,
		InputConfirmations: sqlcInput.InputConfirmations,
		UpdatedAt:          sqlcInput.UpdatedAt,
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
