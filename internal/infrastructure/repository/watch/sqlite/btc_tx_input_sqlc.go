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

// TxInputRepositorySqlc is repository for btc_tx_input table using sqlc (SQLite)
type TxInputRepositorySqlc struct {
	db           *sql.DB
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewBTCTxInputRepositorySqlc returns TxInputRepositorySqlc object
func NewBTCTxInputRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *TxInputRepositorySqlc {
	return &TxInputRepositorySqlc{
		db:           dbConn,
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// convertToBtcTxInput converts sqlcgen.BtcTxInput to domain.BtcTxInput entity
func convertToBtcTxInput(sqlcInput *sqlcgen.BtcTxInput) (*domainBitcoin.BtcTxInput, error) {
	input := &domainBitcoin.BtcTxInput{
		ID:                 sqlcInput.ID,
		TxID:               sqlcInput.TxID,
		InputTxid:          sqlcInput.InputTxid,
		InputVout:          uint32(sqlcInput.InputVout),
		InputAddress:       sqlcInput.InputAddress,
		InputAccount:       sqlcInput.InputAccount,
		InputAmount:        sqlcInput.InputAmount,
		InputConfirmations: uint64(sqlcInput.InputConfirmations),
	}

	// Parse TEXT timestamp
	if sqlcInput.UpdatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", sqlcInput.UpdatedAt.String)
		if err == nil {
			input.UpdatedAt = &t
		}
	}

	return input, nil
}

// convertFromBtcTxInput converts domain.BtcTxInput entity to sqlcgen.BtcTxInput
func convertFromBtcTxInput(input *domainBitcoin.BtcTxInput) *sqlcgen.BtcTxInput {
	sqlcInput := &sqlcgen.BtcTxInput{
		ID:                 input.ID,
		TxID:               input.TxID,
		InputTxid:          input.InputTxid,
		InputVout:          int64(input.InputVout),
		InputAddress:       input.InputAddress,
		InputAccount:       input.InputAccount,
		InputAmount:        input.InputAmount,
		InputConfirmations: int64(input.InputConfirmations),
	}

	// Convert time.Time to TEXT format
	if input.UpdatedAt != nil {
		sqlcInput.UpdatedAt = sql.NullString{
			String: input.UpdatedAt.Format("2006-01-02 15:04:05"),
			Valid:  true,
		}
	}

	return sqlcInput
}

// GetOne get one record by ID
func (r *TxInputRepositorySqlc) GetOne(id int64) (*domainBitcoin.BtcTxInput, error) {
	ctx := context.Background()

	input, err := r.queries.GetBtcTxInputByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxInputByID(): %w", err)
	}

	return convertToBtcTxInput(&input)
}

// GetAllByTxID returns all records by tx ID
func (r *TxInputRepositorySqlc) GetAllByTxID(id int64) ([]*domainBitcoin.BtcTxInput, error) {
	ctx := context.Background()

	inputs, err := r.queries.GetBtcTxInputsByTxID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetBtcTxInputsByTxID(): %w", err)
	}

	result := make([]*domainBitcoin.BtcTxInput, 0, len(inputs))
	for i := range inputs {
		input, err := convertToBtcTxInput(&inputs[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert input: %w", err)
		}
		result = append(result, input)
	}

	return result, nil
}

// Insert inserts one record
func (r *TxInputRepositorySqlc) Insert(txItem *domainBitcoin.BtcTxInput) error {
	ctx := context.Background()

	sqlcInput := convertFromBtcTxInput(txItem)
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
func (r *TxInputRepositorySqlc) InsertBulk(txItems []*domainBitcoin.BtcTxInput) error {
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
		sqlcInput := convertFromBtcTxInput(item)
		_, err := qtx.InsertBtcTxInput(ctx, sqlcgen.InsertBtcTxInputParams{
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
			return fmt.Errorf("failed to insert input: %w", err)
		}
	}

	return tx.Commit()
}
