package watch

import (
	"context"
	"database/sql"
	"fmt"

	portsPersistence "github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// PaymentRequestRepositorySqlc is repository for payment_request table using sqlc
type PaymentRequestRepositorySqlc struct {
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewPaymentRequestRepositorySqlc returns PaymentRequestRepositorySqlc object
func NewPaymentRequestRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *PaymentRequestRepositorySqlc {
	return &PaymentRequestRepositorySqlc{
		queries:      sqlcgen.New(dbConn),
		coinTypeCode: coinTypeCode,
	}
}

// GetAll returns all records whose payment_id is null
func (r *PaymentRequestRepositorySqlc) GetAll() ([]*sqlcgen.PaymentRequest, error) {
	ctx := context.Background()

	requests, err := r.queries.GetAllPaymentRequests(ctx, sqlcgen.PaymentRequestCoin(r.coinTypeCode.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllPaymentRequests(): %w", err)
	}

	result := make([]*sqlcgen.PaymentRequest, len(requests))
	for i := range requests {
		result[i] = &requests[i]
	}

	return result, nil
}

// GetAllByPaymentID returns all records searched by payment_id
func (r *PaymentRequestRepositorySqlc) GetAllByPaymentID(paymentID int64) ([]*sqlcgen.PaymentRequest, error) {
	ctx := context.Background()

	requests, err := r.queries.GetPaymentRequestsByPaymentID(ctx, sqlcgen.GetPaymentRequestsByPaymentIDParams{
		Coin:      sqlcgen.PaymentRequestCoin(r.coinTypeCode.String()),
		PaymentID: sql.NullInt64{Int64: paymentID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetPaymentRequestsByPaymentID(): %w", err)
	}

	result := make([]*sqlcgen.PaymentRequest, len(requests))
	for i := range requests {
		result[i] = &requests[i]
	}

	return result, nil
}

// InsertBulk inserts multiple records
func (r *PaymentRequestRepositorySqlc) InsertBulk(items []*sqlcgen.PaymentRequest) error {
	ctx := context.Background()

	for _, item := range items {
		_, err := r.queries.InsertPaymentRequest(ctx, sqlcgen.InsertPaymentRequestParams{
			Coin:            item.Coin,
			PaymentID:       item.PaymentID,
			SenderAddress:   item.SenderAddress,
			SenderAccount:   item.SenderAccount,
			ReceiverAddress: item.ReceiverAddress,
			Amount:          item.Amount,
			IsDone:          item.IsDone,
			UpdatedAt:       item.UpdatedAt,
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

// WithTx returns a new repository instance that uses the provided transaction
func (r *PaymentRequestRepositorySqlc) WithTx(tx *sql.Tx) portsPersistence.PaymentRequestRepositorier {
	return &PaymentRequestRepositorySqlc{
		queries:      r.queries.WithTx(tx),
		coinTypeCode: r.coinTypeCode,
	}
}
