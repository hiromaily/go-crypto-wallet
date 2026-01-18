package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	repowatch "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/watch"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainPayment "github.com/hiromaily/go-crypto-wallet/internal/domain/payment"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
)

// PaymentRequestRepositorySqlc is repository for payment_request table using sqlc (SQLite)
type PaymentRequestRepositorySqlc struct {
	db           *sql.DB
	queries      *sqlcgen.Queries
	coinTypeCode domainCoin.CoinTypeCode
}

// NewPaymentRequestRepositorySqlc returns PaymentRequestRepositorySqlc object
func NewPaymentRequestRepositorySqlc(
	dbConn *sql.DB, coinTypeCode domainCoin.CoinTypeCode,
) *PaymentRequestRepositorySqlc {
	return &PaymentRequestRepositorySqlc{
		db:           dbConn,
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
		IsDone:          sqlcReq.IsDone != 0, // INTEGER → bool
	}

	if sqlcReq.PaymentID.Valid {
		req.PaymentID = &sqlcReq.PaymentID.Int64
	}

	// Parse TEXT timestamp
	if sqlcReq.UpdatedAt.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", sqlcReq.UpdatedAt.String)
		if err == nil {
			req.UpdatedAt = &t
		}
	}

	return req, nil
}

// convertFromPaymentRequest converts domain.PaymentRequest entity to sqlcgen.PaymentRequest
func convertFromPaymentRequest(req *domainPayment.PaymentRequest) *sqlcgen.PaymentRequest {
	sqlcReq := &sqlcgen.PaymentRequest{
		ID:              req.ID,
		Coin:            req.CoinTypeCode.String(),
		SenderAddress:   req.SenderAddress,
		SenderAccount:   req.SenderAccount,
		ReceiverAddress: req.ReceiverAddress,
		Amount:          req.Amount,
		IsDone:          boolToInt64(req.IsDone), // bool → INTEGER
	}

	if req.PaymentID != nil {
		sqlcReq.PaymentID = sql.NullInt64{Int64: *req.PaymentID, Valid: true}
	}

	// Convert time.Time to TEXT format
	if req.UpdatedAt != nil {
		sqlcReq.UpdatedAt = sql.NullString{
			String: req.UpdatedAt.Format("2006-01-02 15:04:05"),
			Valid:  true,
		}
	}

	return sqlcReq
}

// GetAll returns all records whose payment_id is null
func (r *PaymentRequestRepositorySqlc) GetAll() ([]*domainPayment.PaymentRequest, error) {
	ctx := context.Background()

	requests, err := r.queries.GetAllPaymentRequests(ctx, r.coinTypeCode.String())
	if err != nil {
		return nil, fmt.Errorf("failed to call GetAllPaymentRequests(): %w", err)
	}

	result := make([]*domainPayment.PaymentRequest, 0, len(requests))
	for i := range requests {
		domainReq, err := convertToPaymentRequest(&requests[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert payment request: %w", err)
		}
		result = append(result, domainReq)
	}

	return result, nil
}

// GetAllByPaymentID returns all records searched by payment_id
func (r *PaymentRequestRepositorySqlc) GetAllByPaymentID(paymentID int64) ([]*domainPayment.PaymentRequest, error) {
	ctx := context.Background()

	requests, err := r.queries.GetPaymentRequestsByPaymentID(ctx, sqlcgen.GetPaymentRequestsByPaymentIDParams{
		Coin:      r.coinTypeCode.String(),
		PaymentID: sql.NullInt64{Int64: paymentID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetPaymentRequestsByPaymentID(): %w", err)
	}

	result := make([]*domainPayment.PaymentRequest, 0, len(requests))
	for i := range requests {
		domainReq, err := convertToPaymentRequest(&requests[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert payment request: %w", err)
		}
		result = append(result, domainReq)
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
		IsDone:    1, // bool true → int64 1
		Coin:      r.coinTypeCode.String(),
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

	result, err := r.queries.DeleteAllPaymentRequests(ctx, r.coinTypeCode.String())
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
	tx persistence.Transaction,
) (repowatch.PaymentRequestRepositorier, error) {
	sqlTx := database.UnwrapSQLTx(tx)
	if sqlTx == nil {
		return nil, database.ErrUnsupportedTransaction
	}
	return &PaymentRequestRepositorySqlc{
		db:           r.db,
		queries:      r.queries.WithTx(sqlTx),
		coinTypeCode: r.coinTypeCode,
	}, nil
}
