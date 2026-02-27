package mysql

import (
	"context"
	"database/sql"
	"fmt"

	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainPayment "github.com/hiromaily/go-crypto-wallet/internal/domain/payment"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
	dbtx "github.com/hiromaily/go-crypto-wallet/pkg/db/tx"
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

// convertToPaymentRequest converts sqlcgen.PaymentRequest to domain.PaymentRequest entity
func convertToPaymentRequest(sqlcReq *sqlcgen.PaymentRequest) (*domainPayment.PaymentRequest, error) {
	req := &domainPayment.PaymentRequest{
		ID:              sqlcReq.ID,
		CoinTypeCode:    domainCoin.CoinTypeCode(sqlcReq.Coin),
		SenderAddress:   sqlcReq.SenderAddress,
		SenderAccount:   sqlcReq.SenderAccount,
		ReceiverAddress: sqlcReq.ReceiverAddress,
		Amount:          sqlcReq.Amount,
		IsDone:          sqlcReq.IsDone,
	}

	if sqlcReq.PaymentID.Valid {
		req.PaymentID = &sqlcReq.PaymentID.Int64
	}
	if sqlcReq.UpdatedAt.Valid {
		req.UpdatedAt = &sqlcReq.UpdatedAt.Time
	}

	return req, nil
}

// convertFromPaymentRequest converts domain.PaymentRequest entity to sqlcgen.PaymentRequest
func convertFromPaymentRequest(req *domainPayment.PaymentRequest) *sqlcgen.PaymentRequest {
	sqlcReq := &sqlcgen.PaymentRequest{
		ID:              req.ID,
		Coin:            sqlcgen.PaymentRequestCoin(req.CoinTypeCode.String()),
		SenderAddress:   req.SenderAddress,
		SenderAccount:   req.SenderAccount,
		ReceiverAddress: req.ReceiverAddress,
		Amount:          req.Amount,
		IsDone:          req.IsDone,
	}

	if req.PaymentID != nil {
		sqlcReq.PaymentID = sql.NullInt64{Int64: *req.PaymentID, Valid: true}
	}
	if req.UpdatedAt != nil {
		sqlcReq.UpdatedAt = sql.NullTime{Time: *req.UpdatedAt, Valid: true}
	}

	return sqlcReq
}

// GetAll returns all records whose payment_id is null
func (r *PaymentRequestRepositorySqlc) GetAll() ([]*domainPayment.PaymentRequest, error) {
	ctx := context.Background()

	requests, err := r.queries.GetAllPaymentRequests(ctx, sqlcgen.PaymentRequestCoin(r.coinTypeCode.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllPaymentRequests(): %w", err)
	}

	result := make([]*domainPayment.PaymentRequest, len(requests))
	for i := range requests {
		domainReq, err := convertToPaymentRequest(&requests[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert payment request: %w", err)
		}
		result[i] = domainReq
	}

	return result, nil
}

// GetAllByPaymentID returns all records searched by payment_id
func (r *PaymentRequestRepositorySqlc) GetAllByPaymentID(paymentID int64) ([]*domainPayment.PaymentRequest, error) {
	ctx := context.Background()

	requests, err := r.queries.GetPaymentRequestsByPaymentID(ctx, sqlcgen.GetPaymentRequestsByPaymentIDParams{
		Coin:      sqlcgen.PaymentRequestCoin(r.coinTypeCode.String()),
		PaymentID: sql.NullInt64{Int64: paymentID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetPaymentRequestsByPaymentID(): %w", err)
	}

	result := make([]*domainPayment.PaymentRequest, len(requests))
	for i := range requests {
		domainReq, err := convertToPaymentRequest(&requests[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert payment request: %w", err)
		}
		result[i] = domainReq
	}

	return result, nil
}

// InsertBulk inserts multiple records
func (r *PaymentRequestRepositorySqlc) InsertBulk(items []*domainPayment.PaymentRequest) error {
	ctx := context.Background()

	for _, item := range items {
		sqlcReq := convertFromPaymentRequest(item)
		_, err := r.queries.InsertPaymentRequest(ctx, sqlcgen.InsertPaymentRequestParams{
			Coin:            sqlcReq.Coin,
			PaymentID:       sqlcReq.PaymentID,
			SenderAddress:   sqlcReq.SenderAddress,
			SenderAccount:   sqlcReq.SenderAccount,
			ReceiverAddress: sqlcReq.ReceiverAddress,
			Amount:          sqlcReq.Amount,
			IsDone:          sqlcReq.IsDone,
			UpdatedAt:       sqlcReq.UpdatedAt,
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

// WithTransaction returns a new repository instance that uses the provided transaction
func (r *PaymentRequestRepositorySqlc) WithTransaction(
	tx dbtx.Transaction,
) (repowatch.PaymentRequestRepositorier, error) {
	sqlTx := dbtx.UnwrapSQLTx(tx)
	if sqlTx == nil {
		return nil, dbtx.ErrUnsupportedTransaction
	}
	return &PaymentRequestRepositorySqlc{
		queries:      r.queries.WithTx(sqlTx),
		coinTypeCode: r.coinTypeCode,
	}, nil
}
