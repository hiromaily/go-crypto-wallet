package watchrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/null/v8"

	"github.com/hiromaily/go-crypto-wallet/pkg/db/rdb/sqlcgen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// PaymentRequestRepositorySqlc is repository for payment_request table using sqlc
type PaymentRequestRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode coin.CoinTypeCode
	logger       logger.Logger
}

// NewPaymentRequestRepositorySqlc returns PaymentRequestRepositorySqlc object
func NewPaymentRequestRepositorySqlc(
	dbConn *sql.DB, coinTypeCode coin.CoinTypeCode, logger logger.Logger,
) *PaymentRequestRepositorySqlc {
	return &PaymentRequestRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
		logger:       logger,
	}
}

// GetAll returns all records whose payment_id is null
func (r *PaymentRequestRepositorySqlc) GetAll() ([]*models.PaymentRequest, error) {
	ctx := context.Background()

	requests, err := r.queries.GetAllPaymentRequests(ctx, sqlcgen.PaymentRequestCoin(r.coinTypeCode.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllPaymentRequests(): %w", err)
	}

	// Convert sqlc types to sqlboiler types
	result := make([]*models.PaymentRequest, len(requests))
	for i, req := range requests {
		result[i] = convertSqlcPaymentRequestToModel(&req)
	}

	return result, nil
}

// GetAllByPaymentID returns all records searched by payment_id
func (r *PaymentRequestRepositorySqlc) GetAllByPaymentID(paymentID int64) ([]*models.PaymentRequest, error) {
	ctx := context.Background()

	requests, err := r.queries.GetPaymentRequestsByPaymentID(ctx, sqlcgen.GetPaymentRequestsByPaymentIDParams{
		Coin:      sqlcgen.PaymentRequestCoin(r.coinTypeCode.String()),
		PaymentID: sql.NullInt64{Int64: paymentID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetPaymentRequestsByPaymentID(): %w", err)
	}

	// Convert sqlc types to sqlboiler types
	result := make([]*models.PaymentRequest, len(requests))
	for i, req := range requests {
		result[i] = convertSqlcPaymentRequestToModel(&req)
	}

	return result, nil
}

// InsertBulk inserts multiple records
func (r *PaymentRequestRepositorySqlc) InsertBulk(items []*models.PaymentRequest) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertPaymentRequest(ctx, sqlcgen.InsertPaymentRequestParams{
			Coin:            sqlcgen.PaymentRequestCoin(item.Coin),
			PaymentID:       convertNullInt64ToSQLNullInt64(item.PaymentID),
			SenderAddress:   item.SenderAddress,
			SenderAccount:   item.SenderAccount,
			ReceiverAddress: item.ReceiverAddress,
			Amount:          item.Amount.String(),
			IsDone:          item.IsDone,
			UpdatedAt:       convertNullTimeToSQLNullTime(item.UpdatedAt),
		})
		if err != nil {
			return fmt.Errorf("failed to call InsertPaymentRequest(): %w", err)
		}
	}

	return nil
}

// UpdatePaymentID updates payment_id for multiple IDs
func (r *PaymentRequestRepositorySqlc) UpdatePaymentID(paymentID int64, ids []int64) (int64, error) {
	ctx := context.Background()
	var totalAffected int64

	// sqlc doesn't support IN clauses with variable arguments,
	// so we update one at a time
	for _, id := range ids {
		result, err := r.queries.UpdatePaymentRequestPaymentID(ctx, sqlcgen.UpdatePaymentRequestPaymentIDParams{
			PaymentID: sql.NullInt64{Int64: paymentID, Valid: true},
			ID:        id,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to call UpdatePaymentRequestPaymentID(): %w", err)
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
		}
		totalAffected += affected
	}

	return totalAffected, nil
}

// UpdateIsDone updates isDone
func (r *PaymentRequestRepositorySqlc) UpdateIsDone(paymentID int64) (int64, error) {
	ctx := context.Background()

	result, err := r.queries.UpdatePaymentRequestIsDone(ctx, sqlcgen.UpdatePaymentRequestIsDoneParams{
		IsDone:    true,
		Coin:      sqlcgen.PaymentRequestCoin(r.coinTypeCode.String()),
		PaymentID: sql.NullInt64{Int64: paymentID, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to call UpdatePaymentRequestIsDone(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// DeleteAll deletes all records
func (r *PaymentRequestRepositorySqlc) DeleteAll() (int64, error) {
	ctx := context.Background()

	result, err := r.queries.DeleteAllPaymentRequests(ctx, sqlcgen.PaymentRequestCoin(r.coinTypeCode.String()))
	if err != nil {
		return 0, fmt.Errorf("failed to call DeleteAllPaymentRequests(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}

// Helper functions

func convertSqlcPaymentRequestToModel(req *sqlcgen.PaymentRequest) *models.PaymentRequest {
	amount := new(decimal.Big)
	_, _ = amount.SetString(req.Amount)

	return &models.PaymentRequest{
		ID:              req.ID,
		Coin:            string(req.Coin),
		PaymentID:       convertSQLNullInt64ToNullInt64(req.PaymentID),
		SenderAddress:   req.SenderAddress,
		SenderAccount:   req.SenderAccount,
		ReceiverAddress: req.ReceiverAddress,
		Amount:          amount,
		IsDone:          req.IsDone,
		UpdatedAt:       convertSQLNullTimeToNullTime(req.UpdatedAt),
	}
}

func convertNullInt64ToSQLNullInt64(n null.Int64) sql.NullInt64 {
	if !n.Valid {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: n.Int64, Valid: true}
}

func convertSQLNullInt64ToNullInt64(n sql.NullInt64) null.Int64 {
	if !n.Valid {
		return null.Int64{}
	}
	return null.Int64From(n.Int64)
}
